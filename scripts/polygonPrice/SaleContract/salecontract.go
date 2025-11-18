package SaleContract

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// NewSaleContract - 创建实例：rpcURL 合约地址 私钥(0x 前缀可有可无)
func NewSaleContract(rpcURL, contractAddr, hexPrivKey string) (*SaleContract, error) {
	cli, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("dial rpc failed: %w", err)
	}

	// parse private key (允许带或不带 0x)
	hex := strings.TrimPrefix(hexPrivKey, "0x")
	priv, err := crypto.HexToECDSA(hex)
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %w", err)
	}

	// parse abi
	parsed, err := abi.JSON(strings.NewReader(minimalABI))
	if err != nil {
		return nil, fmt.Errorf("parse abi failed: %w", err)
	}

	// 获取 chain id（让签名和 gas 策略适配当前链）
	ctx := context.Background()
	chainID, err := cli.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf("get chain id failed: %w", err)
	}

	return &SaleContract{
		client:  cli,
		address: common.HexToAddress(contractAddr),
		privKey: priv,
		abi:     parsed,
		chainID: chainID,
	}, nil
}

// ReadSaleData - 查询 price 和 saleEndTimestamp
func (c *SaleContract) ReadSaleData(ctx context.Context) (price *big.Int, saleEnd *big.Int, err error) {
	// pack call data for price()
	dataPrice, err := c.abi.Pack("price")
	if err != nil {
		return nil, nil, fmt.Errorf("abi pack price failed: %w", err)
	}
	msg := ethereum.CallMsg{To: &c.address, Data: dataPrice}
	outPrice, err := c.client.CallContract(ctx, msg, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("call price failed: %w", err)
	}
	// unpack returns []interface{}
	resPrice, err := c.abi.Unpack("price", outPrice)
	if err != nil {
		return nil, nil, fmt.Errorf("abi unpack price failed: %w", err)
	}
	if len(resPrice) == 0 {
		return nil, nil, fmt.Errorf("price: empty result")
	}
	price = resPrice[0].(*big.Int)

	// pack call data for saleEndTimestamp()
	dataEnd, err := c.abi.Pack("saleEndTimestamp")
	if err != nil {
		return nil, nil, fmt.Errorf("abi pack saleEndTimestamp failed: %w", err)
	}
	msg2 := ethereum.CallMsg{To: &c.address, Data: dataEnd}
	outEnd, err := c.client.CallContract(ctx, msg2, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("call saleEndTimestamp failed: %w", err)
	}
	resEnd, err := c.abi.Unpack("saleEndTimestamp", outEnd)
	if err != nil {
		return nil, nil, fmt.Errorf("abi unpack saleEndTimestamp failed: %w", err)
	}
	if len(resEnd) == 0 {
		return nil, nil, fmt.Errorf("saleEndTimestamp: empty result")
	}
	saleEnd = resEnd[0].(*big.Int)

	return price, saleEnd, nil
}

// ConfigureSale - 调用 configureSale(_price, _salePeriod)
// priceWei: *big.Int 单位 wei
// expireTime: 到期时间（本地 time），salePeriod = expireTime - now （秒）
func (c *SaleContract) ConfigureSale(ctx context.Context, priceWei *big.Int, expireTime time.Time) (txHash string, err error) {
	now := time.Now()
	if expireTime.Before(now) {
		return "", fmt.Errorf("expireTime must be after now")
	}
	period := big.NewInt(int64(expireTime.Sub(now).Seconds()))
	if period.Sign() <= 0 {
		return "", fmt.Errorf("salePeriod must be positive")
	}

	// 打包数据
	payload, err := c.abi.Pack("configureSale", priceWei, period)
	if err != nil {
		return "", fmt.Errorf("abi pack failed: %w", err)
	}

	// 构造交易：需要 nonce, gasPrice, gasLimit
	fromAddr := crypto.PubkeyToAddress(c.privKey.PublicKey)
	nonce, err := c.client.PendingNonceAt(ctx, fromAddr)
	if err != nil {
		return "", fmt.Errorf("get nonce failed: %w", err)
	}

	// 建议 gas price
	gasPrice, err := c.client.SuggestGasPrice(ctx)
	if err != nil {
		return "", fmt.Errorf("suggest gas price failed: %w", err)
	}

	// 估算 gas limit（CallMsg 中包含 To, From, Data）
	callMsg := ethereum.CallMsg{
		From:     fromAddr,
		To:       &c.address,
		GasPrice: gasPrice,
		Value:    big.NewInt(0),
		Data:     payload,
	}
	gasLimit, err := c.client.EstimateGas(ctx, callMsg)
	if err != nil {
		// 如果估算失败，使用一个保守默认值
		gasLimit = uint64(300000) // fallback
	}

	// 创建交易（value 为0）
	tx := types.NewTransaction(nonce, c.address, big.NewInt(0), gasLimit, gasPrice, payload)

	// 签名
	signed, err := types.SignTx(tx, types.NewEIP155Signer(c.chainID), c.privKey)
	if err != nil {
		return "", fmt.Errorf("sign tx failed: %w", err)
	}

	// 发送
	if err := c.client.SendTransaction(ctx, signed); err != nil {
		return "", fmt.Errorf("send tx failed: %w", err)
	}

	return signed.Hash().Hex(), nil
}

// Helper: 打印 ABI（调试用）
func (c *SaleContract) ABIJSON() string {
	b, _ := json.MarshalIndent(c.abi, "", "  ")
	return string(b)
}
