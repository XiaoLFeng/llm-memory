package todo

import (
	"context"
	"os"

	"github.com/XiaoLFeng/llm-memory/internal/cli"
	"github.com/XiaoLFeng/llm-memory/internal/cli/handlers"
	"github.com/XiaoLFeng/llm-memory/startup"
	"github.com/spf13/cobra"
)

var todoCompleteCode string

// todoCompleteCmd 完成待办
var todoCompleteCmd = &cobra.Command{
	Use:   "complete",
	Short: "完成待办事项",
	Long:  `将指定待办事项标记为已完成~ 🎉`,
	Run: func(cmd *cobra.Command, args []string) {
		if todoCompleteCode == "" {
			cli.PrintError("请使用 --code 参数指定有效的待办标识码")
			os.Exit(1)
		}

		bs := startup.New(
			startup.WithSignalHandler(false),
		).MustInitialize(context.Background())
		defer bs.Shutdown()

		handler := handlers.NewTodoHandler(bs)
		if err := handler.Complete(bs.Context(), todoCompleteCode); err != nil {
			cli.PrintError(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	todoCompleteCmd.Flags().StringVarP(&todoCompleteCode, "code", "c", "", "待办标识码（必填）")
	_ = todoCompleteCmd.MarkFlagRequired("code")

	todoCmd.AddCommand(todoCompleteCmd)
}
