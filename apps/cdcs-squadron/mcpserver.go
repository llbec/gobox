package cdcs_squadron

import "github.com/modelcontextprotocol/go-sdk/mcp"

func createMcpServer() *mcp.Server {
	server := mcp.NewServer(&mcp.Implementation{Name: "cdcs-squadron", Version: "v1.0.0"}, nil)
	return server
}
