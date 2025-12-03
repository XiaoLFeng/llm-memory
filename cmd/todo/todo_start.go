package todo

import (
	"context"
	"os"

	"github.com/XiaoLFeng/llm-memory/internal/cli"
	"github.com/XiaoLFeng/llm-memory/internal/cli/handlers"
	"github.com/XiaoLFeng/llm-memory/startup"
	"github.com/spf13/cobra"
)

var todoStartID int

// todoStartCmd 开始待办
var todoStartCmd = &cobra.Command{
	Use:   "start",
	Short: "开始待办事项",
	Long:  `将指定待办事项标记为进行中~ 🚀`,
	Run: func(cmd *cobra.Command, args []string) {
		if todoStartID <= 0 {
			cli.PrintError("请使用 --id 参数指定有效的待办ID")
			os.Exit(1)
		}

		bs := startup.New(
			startup.WithSignalHandler(false),
		).MustInitialize(context.Background())
		defer bs.Shutdown()

		handler := handlers.NewTodoHandler(bs)
		if err := handler.Start(bs.Context(), int64(todoStartID)); err != nil {
			cli.PrintError(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	todoStartCmd.Flags().IntVarP(&todoStartID, "id", "i", 0, "待办ID（必填）")
	_ = todoStartCmd.MarkFlagRequired("id")

	todoCmd.AddCommand(todoStartCmd)
}
