package polygon

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// ------------------------
// 自定义 Timestamp 类型（兼容 number / "number"）
// ------------------------
type Timestamp int64

// UnmarshalJSON 实现 JSON 反序列化逻辑
func (t *Timestamp) UnmarshalJSON(data []byte) error {
	// 1️⃣ 去掉可能的引号，比如 "1760126400000"
	var raw int64
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("cannot parse timestamp: %v", err)
	}

	// 2️⃣ 直接将毫秒级时间戳转为 time.Time
	*t = Timestamp(raw)
	return nil
}

// ToTime 将 Timestamp 转换为 time.Time
func (t Timestamp) ToTime() time.Time {
	return time.UnixMilli(int64(t))
}

// ------------------------
// ✅ 数据结构定义
// ------------------------
type PriceInfo struct {
	Symbol            string  `json:"symbol"`
	Price             float64 `json:"price"`
	Volume            float64 `json:"volume"`
	Timestamp         int64   `json:"timestamp"`
	TimestampReadable string  `json:"timestamp_readable,omitempty"`
	UpdatedAt         string  `json:"updated_at"`
}

// ------------------------
// ✅ Polygon 响应结构
// ------------------------
type PolygonResponse struct {
	Ticker string `json:"ticker"`
	//QueryCount   int    `json:"queryCount"`
	//ResultsCount int    `json:"resultsCount"`
	//Adjusted     bool   `json:"adjusted"`
	Results []struct {
		Ticker string  `json:"T"`
		Volume float64 `json:"v"`
		//VWAP      float64   `json:"vw"`
		//Open      float64   `json:"o"`
		Close float64 `json:"c"`
		//High      float64   `json:"h"`
		//Low       float64   `json:"l"`
		Timestamp Timestamp `json:"t"`
		//NumTrades int       `json:"n"`
	} `json:"results"`
	Status string `json:"status"`
	//RequestID string `json:"request_id"`
	//Count     int    `json:"count"`
	Error string `json:"error,omitempty"`
}

type Polygon struct {
	ApiKey      string
	Symbol      string
	BaseURL     string
	LastPrice   float64
	LastUpdated string
	Client      *http.Client
}
