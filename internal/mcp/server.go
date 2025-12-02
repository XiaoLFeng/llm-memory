package mcp

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/XiaoLFeng/llm-memory/internal/mcp/tools"
	"github.com/XiaoLFeng/llm-memory/startup"
)

// Server MCP 服务器
// 嘿嘿~ 使用官方 SDK 实现的 MCP 服务器！(´∀｀)💖
type Server struct {
	bs     *startup.Bootstrap
	server *mcp.Server
}

// NewServer 创建新的 MCP 服务器
// 呀~ 初始化服务器并注册所有工具！✨
func NewServer(bs *startup.Bootstrap) *Server {
	// 创建 MCP 服务器
	mcpServer := mcp.NewServer(&mcp.Implementation{
		Name:    "llm-memory",
		Version: "0.0.1",
	}, nil)

	s := &Server{
		bs:     bs,
		server: mcpServer,
	}

	// 注册所有工具
	s.registerTools()

	return s
}

// Run 运行 MCP 服务器
// 嘿嘿~ 通过 stdio 传输运行服务！🚀
func (s *Server) Run() error {
	return s.server.Run(context.Background(), &mcp.StdioTransport{})
}

// registerTools 注册所有 MCP 工具
// 嘿嘿~ 使用 tools 包中的注册函数！✨
func (s *Server) registerTools() {
	// 记忆管理工具
	tools.RegisterMemoryTools(s.server, s.bs)
	// 计划管理工具
	tools.RegisterPlanTools(s.server, s.bs)
	// TODO 管理工具
	tools.RegisterTodoTools(s.server, s.bs)
}
