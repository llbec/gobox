package polygon

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// ------------------------
// ✅ 构造函数
// ------------------------
func NewPolygon(apiKey, symbol string) *Polygon {
	if symbol == "" {
		symbol = "LHSW"
	}
	return &Polygon{
		ApiKey:  apiKey,
		Symbol:  strings.ToUpper(symbol),
		BaseURL: "https://api.polygon.io",
		Client:  &http.Client{Timeout: 10 * time.Second},
	}
}

// ------------------------
// ✅ 网络请求层
// ------------------------
func (r *Polygon) fetchPriceData(maxRetries int) ([]byte, error) {
	url := fmt.Sprintf("%s/v2/aggs/ticker/%s/prev?apiKey=%s", r.BaseURL, r.Symbol, r.ApiKey)

	for attempt := 1; attempt <= maxRetries; attempt++ {
		resp, err := r.Client.Get(url)
		if err != nil {
			if attempt < maxRetries {
				log.Printf("请求失败（%d/%d），10秒后重试...\n", attempt, maxRetries)
				time.Sleep(10 * time.Second)
				continue
			}
			return nil, fmt.Errorf("网络错误: %v", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("读取响应失败: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("HTTP错误: %d, 内容: %s", resp.StatusCode, string(body))
		}

		log.Printf("成功获取数据（%d/%d）\nresponse:\n%s\n", attempt, maxRetries, string(body))

		return body, nil
	}
	return nil, fmt.Errorf("多次请求失败")
}

// ------------------------
// ✅ 数据解析层
// ------------------------
func (r *Polygon) parsePriceData(body []byte) (*PriceInfo, error) {
	var data PolygonResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("JSON解析失败: %v", err)
	}

	if strings.ToUpper(data.Status) != "OK" {
		return nil, fmt.Errorf("API响应错误: %s", data.Error)
	}
	if len(data.Results) == 0 {
		return nil, fmt.Errorf("未找到价格数据")
	}

	bar := data.Results[0]
	ts := int64(bar.Timestamp)

	priceInfo := &PriceInfo{
		Symbol:            r.Symbol,
		Price:             bar.Close,
		Volume:            bar.Volume,
		Timestamp:         ts,
		TimestampReadable: time.UnixMilli(ts).Format("2006-01-02 15:04:05"),
		UpdatedAt:         time.Now().Format("2006-01-02 15:04:05"),
	}
	return priceInfo, nil
}

// ------------------------
// ✅ 业务逻辑层
// ------------------------
func (r *Polygon) GetRealtimePrice(maxRetries int) (*PriceInfo, error) {
	body, err := r.fetchPriceData(maxRetries)
	if err != nil {
		return nil, err
	}

	priceInfo, err := r.parsePriceData(body)
	if err != nil {
		return nil, err
	}

	r.LastPrice = priceInfo.Price
	r.LastUpdated = priceInfo.UpdatedAt

	/*fileName := fmt.Sprintf("%s_realtime_price.json", r.Symbol)
	fileBytes, _ := json.MarshalIndent(priceInfo, "", "  ")
	if err := os.WriteFile(fileName, fileBytes, 0644); err != nil {
		return nil, fmt.Errorf("写入文件失败: %v", err)
	}*/

	return priceInfo, nil
}
