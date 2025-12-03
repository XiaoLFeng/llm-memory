package todo

import (
	"context"
	"os"

	"github.com/XiaoLFeng/llm-memory/internal/cli"
	"github.com/XiaoLFeng/llm-memory/internal/cli/handlers"
	"github.com/XiaoLFeng/llm-memory/startup"
	"github.com/spf13/cobra"
)

var todoDeleteID int

// todoDeleteCmd 删除待办
var todoDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "删除待办事项",
	Long:  `删除指定ID的待办事项~ 🗑️`,
	Run: func(cmd *cobra.Command, args []string) {
		if todoDeleteID <= 0 {
			cli.PrintError("请使用 --id 参数指定有效的待办ID")
			os.Exit(1)
		}

		bs := startup.New(
			startup.WithSignalHandler(false),
		).MustInitialize(context.Background())
		defer bs.Shutdown()

		handler := handlers.NewTodoHandler(bs)
		if err := handler.Delete(bs.Context(), int64(todoDeleteID)); err != nil {
			cli.PrintError(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	todoDeleteCmd.Flags().IntVarP(&todoDeleteID, "id", "i", 0, "待办ID（必填）")
	_ = todoDeleteCmd.MarkFlagRequired("id")

	todoCmd.AddCommand(todoDeleteCmd)
}
