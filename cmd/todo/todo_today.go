package todo

import (
	"context"
	"os"

	"github.com/XiaoLFeng/llm-memory/internal/cli"
	"github.com/XiaoLFeng/llm-memory/internal/cli/handlers"
	"github.com/XiaoLFeng/llm-memory/startup"
	"github.com/spf13/cobra"
)

// todoTodayCmd 获取今日待办
// 嘿嘿~ 查看今天要做的事！📅
var todoTodayCmd = &cobra.Command{
	Use:   "today",
	Short: "获取今日待办事项",
	Long:  `获取今天的待办事项列表~ 📅`,
	Run: func(cmd *cobra.Command, args []string) {
		bs := startup.New(
			startup.WithSignalHandler(false),
		).MustInitialize(context.Background())
		defer bs.Shutdown()

		handler := handlers.NewTodoHandler(bs)
		if err := handler.Today(bs.Context()); err != nil {
			cli.PrintError(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	todoCmd.AddCommand(todoTodayCmd)
}
