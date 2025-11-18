package main

import (
	"log"
	"os"
	"polygonprice/polygon"
)

type Config struct {
	APIKey          string `yaml:"api_key"`          // Polygon API Key，用于访问行情数据
	Symbol          string `yaml:"symbol"`           // 需要查询的代币或股票符号，例如 "ETH" 或 "AAPL"
	RPCUrl          string `yaml:"rpc_url"`          // 区块链节点 RPC 地址，用于调用合约接口
	PrivateKey      string `yaml:"private_key"`      // 调用合约的私钥（需具备 updaterole 权限）
	ContractAddress string `yaml:"contract_address"` // 价格设置合约地址
	Date            string `yaml:"expire_date"`      // 价格有效截止时间（到秒），例如 "2025-10-12 23:59:59"
}

// ------------------------
// ✅ 主程序
// ------------------------
func main() {
	apiKey := os.Getenv("POLYGON_API_KEY")
	if apiKey == "" {
		apiKey = "FhYkL9lATs7NoHnzKDSpwsd_YESfyp74"
	}

	rlinkFi := polygon.NewPolygon(apiKey, "LHSW")

	log.Println("获取15分钟延迟的实时价格...")
	priceData, err := rlinkFi.GetRealtimePrice(3)
	if err != nil {
		log.Println("错误:", err)
		os.Exit(1)
	}

	log.Printf("✅ 成功获取 %s 价格: %.2f USD (%s)\n", priceData.Symbol, priceData.Price, priceData.UpdatedAt)
	//fmt.Printf("📁 数据已保存至 %s_realtime_price.json\n", priceData.Symbol)
}
