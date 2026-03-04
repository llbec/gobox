package main

import (
	"gobox/tools/act"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func createMcpServer() *mcp.Server {
	server := mcp.NewServer(&mcp.Implementation{Name: "cdcs-squadron", Version: "v1.0.0"}, nil)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "set_config",
		Description: "设置目标路径和备份路径配置",
	}, act.SetConfig)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "archive_files",
		Description: "归档指定时间之前的文件",
	}, act.Archive)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "filter_files",
		Description: "筛选符合条件的文件",
	}, act.FilterFiles)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "resend_files",
		Description: "重传指定文件",
	}, act.Resend)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "modify_files",
		Description: "修改文件的患者ID并重新发送",
	}, act.Modify)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "rollback_operation",
		Description: "回滚指定的操作",
	}, act.Rollback)

	return server
}
