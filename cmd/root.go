package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var Version = "0.0.1"

// rootCmd 是应用的根命令
// 呀~ 这是所有子命令的入口点！(´∀｀)💖
var rootCmd = &cobra.Command{
	Use:   "llm-memory",
	Short: "LLM-Memory - 大模型统一记忆系统",
	Long: `LLM-Memory 是一个为大模型设计的统一记忆管理系统。

嘿嘿~ 支持记忆管理、计划管理和 TODO 管理功能！
可以通过 GUI 界面操作，也可以作为 MCP 服务运行~ ✨`,
	// 如果直接运行 llm-memory 不带子命令，显示帮助
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

// Execute 执行根命令
// 这是程序的入口点~ 🚀
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// 添加版本标志
	rootCmd.Version = Version
	rootCmd.SetVersionTemplate("LLM-Memory 版本: {{.Version}}\n")
}
