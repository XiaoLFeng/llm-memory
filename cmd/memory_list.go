package cmd

import (
	"context"
	"os"

	"github.com/XiaoLFeng/llm-memory/internal/cli"
	"github.com/XiaoLFeng/llm-memory/internal/cli/handlers"
	"github.com/XiaoLFeng/llm-memory/startup"
	"github.com/spf13/cobra"
)

// memoryListCmd 列出所有记忆
// 呀~ 查看所有记忆条目！✨
var memoryListCmd = &cobra.Command{
	Use:   "list",
	Short: "列出所有记忆",
	Long:  `列出系统中保存的所有记忆条目~ 📚`,
	Run: func(cmd *cobra.Command, args []string) {
		bs := startup.New(
			startup.WithSignalHandler(false),
		).MustInitialize(context.Background())
		defer bs.Shutdown()

		handler := handlers.NewMemoryHandler(bs)
		if err := handler.List(bs.Context()); err != nil {
			cli.PrintError(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	memoryCmd.AddCommand(memoryListCmd)
}
