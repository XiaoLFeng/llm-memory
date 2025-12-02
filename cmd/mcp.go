package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/XiaoLFeng/llm-memory/internal/mcp"
	"github.com/XiaoLFeng/llm-memory/startup"
	"github.com/spf13/cobra"
)

// mcpCmd 是 mcp 子命令
// 呀~ 启动 MCP 服务！(´∀｀)💖
var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "启动 MCP 服务",
	Long: `启动 LLM-Memory 的 MCP (Model Context Protocol) 服务。

MCP 服务支持以下功能：
  - 记忆管理：增删改查记忆内容
  - 计划管理：创建、更新、查询计划
  - TODO 管理：管理待办事项

嘿嘿~ AI 模型可以通过 MCP 协议与此服务通信！✨`,
	Run: func(cmd *cobra.Command, args []string) {
		runMCP()
	},
}

func init() {
	rootCmd.AddCommand(mcpCmd)
}

// runMCP 运行 MCP 服务
// 嘿嘿~ 使用 startup 包统一初始化！✨
func runMCP() {
	// 使用 startup 包统一初始化
	bs := startup.New(
		startup.WithSignalHandler(true),
	).MustInitialize(context.Background())
	defer bs.Shutdown()

	// 启动 MCP 服务
	server := mcp.NewServer(bs)
	if err := server.Run(); err != nil {
		fmt.Printf("MCP 服务运行出错: %v\n", err)
		os.Exit(1)
	}
}
