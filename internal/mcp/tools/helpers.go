package tools

import "github.com/modelcontextprotocol/go-sdk/mcp"

// NewTextResult 创建文本结果
// 嘿嘿~ 封装一下官方 SDK 的结果创建！(´∀｀)
func NewTextResult(text string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: text},
		},
	}
}

// NewErrorResult 创建错误结果
// 呀~ 出错时返回这个！💫
func NewErrorResult(errMsg string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		IsError: true,
		Content: []mcp.Content{
			&mcp.TextContent{Text: errMsg},
		},
	}
}
