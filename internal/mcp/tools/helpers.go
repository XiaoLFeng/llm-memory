package tools

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/XiaoLFeng/llm-memory/pkg/types"
)

// NewTextResult 创建文本结果
func NewTextResult(text string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: text},
		},
	}
}

// NewErrorResult 创建错误结果
func NewErrorResult(errMsg string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		IsError: true,
		Content: []mcp.Content{
			&mcp.TextContent{Text: errMsg},
		},
	}
}

// getScopeTagWithContext 根据 PathID 和作用域上下文返回中文作用域标签
// 优先匹配组路径，其次个人路径（无 Global 支持）
func getScopeTagWithContext(pathID int64, scopeCtx *types.ScopeContext) string {
	if pathID == 0 {
		return "[未知]"
	}
	if scopeCtx != nil {
		for _, gid := range scopeCtx.GroupPathIDs {
			if pathID == gid {
				return "[小组]"
			}
		}
		if scopeCtx.PathID == pathID {
			return "[私有]"
		}
	}
	// 缺省认为是私有
	return "[私有]"
}

// getScopeTagWithGlobal 根据 Global 标记和 PathID 返回中文作用域标签（供 Memory 使用）
// 优先匹配 Global，其次组路径，最后个人路径
func getScopeTagWithGlobal(global bool, pathID int64, scopeCtx *types.ScopeContext) string {
	if global {
		return "[全局]"
	}
	if pathID == 0 {
		return "[私有]"
	}
	if scopeCtx != nil {
		for _, gid := range scopeCtx.GroupPathIDs {
			if pathID == gid {
				return "[小组]"
			}
		}
		if scopeCtx.PathID == pathID {
			return "[私有]"
		}
	}
	// 缺省认为是私有
	return "[私有]"
}
