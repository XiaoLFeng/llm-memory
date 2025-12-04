package todo

import (
	"context"
	"os"

	"github.com/XiaoLFeng/llm-memory/internal/cli"
	"github.com/XiaoLFeng/llm-memory/internal/cli/handlers"
	"github.com/XiaoLFeng/llm-memory/startup"
	"github.com/spf13/cobra"
)

var todoDeleteCode string

// todoDeleteCmd 删除待办
var todoDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "删除待办事项",
	Long:  `删除指定标识码的待办事项~ 🗑️`,
	Run: func(cmd *cobra.Command, args []string) {
		if todoDeleteCode == "" {
			cli.PrintError("请使用 --code 参数指定有效的待办标识码")
			os.Exit(1)
		}

		bs := startup.New(
			startup.WithSignalHandler(false),
		).MustInitialize(context.Background())
		defer bs.Shutdown()

		handler := handlers.NewTodoHandler(bs)
		if err := handler.Delete(bs.Context(), todoDeleteCode); err != nil {
			cli.PrintError(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	todoDeleteCmd.Flags().StringVarP(&todoDeleteCode, "code", "c", "", "待办标识码（必填）")
	_ = todoDeleteCmd.MarkFlagRequired("code")

	todoCmd.AddCommand(todoDeleteCmd)
}
