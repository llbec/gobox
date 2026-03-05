package drage

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type DecodeInput struct {
	HexString string `json:"hexString" jsonschema:"十六进制字符串，包含要解码的数据"`
}

type DecodeOutput struct {
	Result string `json:"result" jsonschema:"解码结果信息，包含检测项数量和每个检测项的详细信息"`
}

func Decode(ctx context.Context, req *mcp.CallToolRequest, input DecodeInput) (
	*mcp.CallToolResult,
	DecodeOutput,
	error,
) {
	bytes, err := hexStringToBytes(input.HexString)
	if err != nil {
		return nil, DecodeOutput{}, fmt.Errorf("解析十六进制字符串失败: %w", err)
	}

	result := decode(bytes)
	return nil, DecodeOutput{Result: result}, nil
}
