package cmd

import (
	"context"
	"os"

	"github.com/XiaoLFeng/llm-memory/internal/cli"
	"github.com/XiaoLFeng/llm-memory/internal/cli/handlers"
	"github.com/XiaoLFeng/llm-memory/startup"
	"github.com/spf13/cobra"
)

// todoListCmd 列出所有待办
// 呀~ 查看所有待办事项！✨
var todoListCmd = &cobra.Command{
	Use:   "list",
	Short: "列出所有待办事项",
	Long:  `列出系统中的所有待办事项~ 📝`,
	Run: func(cmd *cobra.Command, args []string) {
		bs := startup.New(
			startup.WithSignalHandler(false),
		).MustInitialize(context.Background())
		defer bs.Shutdown()

		handler := handlers.NewTodoHandler(bs)
		if err := handler.List(bs.Context()); err != nil {
			cli.PrintError(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	todoCmd.AddCommand(todoListCmd)
}
