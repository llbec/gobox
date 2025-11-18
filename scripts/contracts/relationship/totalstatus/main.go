package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/viper"
)

type MemberInfo struct {
	Member    common.Address
	Referrer  common.Address
	Timestamp *big.Int
	Buys      *big.Int
	Payments  *big.Int
}

func main() {
	// 1. 加载配置文件
	viper.SetConfigFile("config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}

	rpcURL := viper.GetString("rpc_url")
	saleAddr := common.HexToAddress(viper.GetString("contracts.sale"))
	relationAddr := common.HexToAddress(viper.GetString("contracts.relationship"))

	// 2. 连接 RPC
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatalf("无法连接 RPC: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	// 获取 chain ID
	chainID, err := client.NetworkID(ctx)
	if err != nil {
		log.Fatalf("获取 chain ID 失败: %v", err)
	}

	// 简单映射 chain name
	var chainName string
	switch chainID.Int64() {
	case 1:
		chainName = "Ethereum Mainnet"
	case 56:
		chainName = "BSC Mainnet"
	case 137:
		chainName = "Polygon Mainnet"
	default:
		chainName = fmt.Sprintf("Chain (%d)", chainID.Int64())
	}
	fmt.Printf("✅ 已连接到: %v\n", chainName)

	// 3. 定义 ABI（简化版，只包含需要的函数）
	saleABI := `[{"inputs":[],"name":"totalSold","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"}]`
	relationABI := `[
		{"inputs":[],"name":"getTotalInfo","outputs":[
			{"internalType":"uint256","name":"totalMembers","type":"uint256"},
			{"internalType":"uint256","name":"totalBuys","type":"uint256"},
			{"internalType":"uint256","name":"totalPayments","type":"uint256"}
		],"stateMutability":"view","type":"function"},
		{"inputs":[],"name":"totalSupply","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},
		{"inputs":[{"internalType":"uint256","name":"id","type":"uint256"}],"name":"ownerOf","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},
		{"inputs":[{"internalType":"uint256","name":"id","type":"uint256"}],"name":"getGroupMembers","outputs":[{"components":[
			{"internalType":"address","name":"member","type":"address"},
			{"internalType":"address","name":"referrer","type":"address"},
			{"internalType":"uint256","name":"timestamp","type":"uint256"},
			{"internalType":"uint256","name":"buys","type":"uint256"},
			{"internalType":"uint256","name":"payments","type":"uint256"}
		],"internalType":"struct MemberInfo[]","name":"","type":"tuple[]"}],"stateMutability":"view","type":"function"}
	]`

	saleParsed, err := abi.JSON(strings.NewReader(saleABI))
	if err != nil {
		log.Fatalf("Sale ABI 解析失败: %v", err)
	}
	relationParsed, err := abi.JSON(strings.NewReader(relationABI))
	if err != nil {
		log.Fatalf("Relationship ABI 解析失败: %v", err)
	}

	// 4. 调用 totalSold
	totalSoldData, _ := saleParsed.Pack("totalSold")
	soldBytes, err := client.CallContract(ctx, callMsg(saleAddr, totalSoldData), nil)
	if err != nil {
		log.Fatalf("调用 totalSold 失败: %v", err)
	}
	var totalSold *big.Int
	_ = saleParsed.UnpackIntoInterface(&totalSold, "totalSold", soldBytes)
	fmt.Println("Total Sold:", formatBig(totalSold, 18))

	// 5. 调用 getTotalInfo
	data, _ := relationParsed.Pack("getTotalInfo")
	output, err := client.CallContract(ctx, callMsg(relationAddr, data), nil)
	if err != nil {
		log.Fatalf("调用 getTotalInfo 失败: %v", err)
	}
	var totalMembers, totalBuys, totalPayments *big.Int
	err = relationParsed.UnpackIntoInterface(&[]interface{}{&totalMembers, &totalBuys, &totalPayments}, "getTotalInfo", output)
	if err != nil {
		log.Fatalf("解包 getTotalInfo 失败: %v", err)
	}
	fmt.Println("Total Members:", totalMembers)
	fmt.Println("Total Buys:", formatBig(totalBuys, 18))
	fmt.Println("Total Payments:", formatBig(totalPayments, 18))

	// 6. 调用 totalSupply
	data, _ = relationParsed.Pack("totalSupply")
	out, err := client.CallContract(ctx, callMsg(relationAddr, data), nil)
	if err != nil {
		log.Fatalf("调用 totalSupply 失败: %v", err)
	}
	var totalSupply *big.Int
	_ = relationParsed.UnpackIntoInterface(&totalSupply, "totalSupply", out)
	fmt.Println("Total Groups:", totalSupply)

	// 7. 循环获取组成员
	for i := big.NewInt(1); i.Cmp(totalSupply) <= 0; i.Add(i, big.NewInt(1)) {
		// ownerOf
		data, _ = relationParsed.Pack("ownerOf", i)
		out, err = client.CallContract(ctx, callMsg(relationAddr, data), nil)
		if err != nil {
			log.Fatalf("调用 ownerOf 失败: %v", err)
		}
		var owner common.Address
		_ = relationParsed.UnpackIntoInterface(&owner, "ownerOf", out)
		fmt.Printf("\nGroup %d owner: %s\n", i, owner.Hex())

		// getGroupMembers
		data, _ = relationParsed.Pack("getGroupMembers", i)
		out, err = client.CallContract(ctx, callMsg(relationAddr, data), nil)
		if err != nil {
			log.Fatalf("调用 getGroupMembers 失败: %v", err)
		}

		var members []MemberInfo
		if err := relationParsed.UnpackIntoInterface(&members, "getGroupMembers", out); err != nil {
			log.Fatalf("解包成员失败: %v", err)
		}

		fmt.Println("Members count:", len(members))
		totalBuys := big.NewInt(0)
		totalPayments := big.NewInt(0)
		for _, m := range members {
			//t := time.Unix(m.Timestamp.Int64(), 0)
			fmt.Printf("\t%s  Buys: %12s  Payment: %12s referrer:%s\n",
				m.Member.Hex(),
				formatBig(m.Buys, 18),
				formatBig(m.Payments, 18),
				m.Referrer.Hex(),
				//t,
			)
			totalBuys.Add(totalBuys, m.Buys)
			totalPayments.Add(totalPayments, m.Payments)
		}
		fmt.Println("\tTotal Buys:", formatBig(totalBuys, 18))
		fmt.Println("\tTotal Payments:", formatBig(totalPayments, 18))
	}

	fmt.Println("程序执行完毕，按回车键退出...")
	fmt.Scanln() // 等待用户输入（按回车）
}

func callMsg(to common.Address, data []byte) ethereum.CallMsg {
	return ethereum.CallMsg{
		To:   &to,
		Data: data,
	}
}

func formatBig(num *big.Int, decimals int) string {
	if num == nil {
		return "0.000000"
	}
	f := new(big.Float).SetInt(num)
	denom := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil))
	val := new(big.Float).Quo(f, denom)
	return fmt.Sprintf("%.6f", val) // 固定保留 6 位小数
}
