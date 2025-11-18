package SaleContract

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// SaleContract 封装
type SaleContract struct {
	client  *ethclient.Client
	address common.Address
	abi     abi.ABI
	privKey *ecdsa.PrivateKey
	chainID *big.Int
}

// Minimal ABI：只包含 configureSale, price(), saleEndTimestamp()
const minimalABI = `[
  {"inputs":[{"internalType":"uint256","name":"_price","type":"uint256"},{"internalType":"uint256","name":"_salePeriod","type":"uint256"}],"name":"configureSale","outputs":[],"stateMutability":"nonpayable","type":"function"},
  {"inputs":[],"name":"price","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},
  {"inputs":[],"name":"saleEndTimestamp","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"}
]`
