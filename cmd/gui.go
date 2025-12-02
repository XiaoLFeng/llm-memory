package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/XiaoLFeng/llm-memory/internal/gui"
	"github.com/XiaoLFeng/llm-memory/startup"
	"github.com/spf13/cobra"
)

// guiCmd 是 gui 子命令
// 呀~ 启动图形管理界面！(´∀｀)💖
var guiCmd = &cobra.Command{
	Use:   "gui",
	Short: "启动图形管理界面",
	Long:  `启动 LLM-Memory 的图形管理界面，可以进行记忆、计划和 TODO 的管理操作~ ✨`,
	Run: func(cmd *cobra.Command, args []string) {
		runGUI()
	},
}

func init() {
	rootCmd.AddCommand(guiCmd)
}

// runGUI 运行图形界面
func runGUI() {
	// 使用 startup 包统一初始化
	bs := startup.New(
		startup.WithSignalHandler(true),
	).MustInitialize(context.Background())
	defer bs.Shutdown()

	// 启动 GUI
	if err := gui.Run(bs); err != nil {
		fmt.Printf("GUI 运行出错: %v\n", err)
		os.Exit(1)
	}
}
