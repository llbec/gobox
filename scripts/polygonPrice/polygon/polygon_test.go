package polygon

import (
	"fmt"
	"testing"
)

func Test_FetchPriceData(t *testing.T) {
	rlinkFi := NewPolygon("FhYkL9lATs7NoHnzKDSpwsd_YESfyp74", "LHSW")
	maxRetries := 3

	body, err := rlinkFi.fetchPriceData(maxRetries)
	if err != nil {
		t.Errorf("fetchPriceData failed: %v", err)
	}

	if len(body) == 0 {
		t.Errorf("fetchPriceData returned empty body")
	}
	fmt.Println(string(body))
}

func Test_ParsePriceData(t *testing.T) {
	rlinkFi := NewPolygon("FhYkL9lATs7NoHnzKDSpwsd_YESfyp74", "LHSW")

	//{"ticker":"LHSW","queryCount":1,"resultsCount":1,"adjusted":true,"results":[{"T":"LHSW","v":282884,"vw":2.122,"o":2.11,"c":2.07,"h":2.12,"l":2.01,"t":1760126400000,"n":299}],"status":"OK","request_id":"8dd428b70abeb1941e5742782eadeb79","count":1}

	body := "{\"ticker\":\"LHSW\",\"queryCount\":1,\"resultsCount\":1,\"adjusted\":true,\"results\":[{\"T\":\"LHSW\",\"v\":282884,\"vw\":2.122,\"o\":2.11,\"c\":2.07,\"h\":2.12,\"l\":2.01,\"t\":1760126400000,\"n\":299}],\"status\":\"OK\",\"request_id\":\"8dd428b70abeb1941e5742782eadeb79\",\"count\":1}"

	priceInfo, err := rlinkFi.parsePriceData([]byte(body))
	if err != nil {
		t.Errorf("parsePriceData failed: %v", err)
	}

	if priceInfo == nil {
		t.Errorf("parsePriceData returned nil priceInfo")
	}
	fmt.Println(priceInfo)
}
