package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/XiaoLFeng/llm-memory/internal/tui"
	"github.com/XiaoLFeng/llm-memory/startup"
	"github.com/spf13/cobra"
)

// tuiCmd 是 tui 子命令
// 呀~ 启动终端用户界面！(´∀｀)💖
var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "启动终端用户界面",
	Long:  `启动 LLM-Memory 的终端用户界面，可以进行记忆、计划和待办的管理操作~ ✨`,
	Run: func(cmd *cobra.Command, args []string) {
		runTUI()
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}

// runTUI 运行终端用户界面
func runTUI() {
	// 使用 startup 包统一初始化
	bs := startup.New(
		startup.WithSignalHandler(true),
	).MustInitialize(context.Background())
	defer bs.Shutdown()

	// 启动 TUI
	if err := tui.Run(bs); err != nil {
		fmt.Printf("TUI 运行出错: %v\n", err)
		os.Exit(1)
	}
}
