package todo

import (
	"context"
	"os"

	"github.com/XiaoLFeng/llm-memory/internal/cli"
	"github.com/XiaoLFeng/llm-memory/internal/cli/handlers"
	"github.com/XiaoLFeng/llm-memory/startup"
	"github.com/spf13/cobra"
)

var todoStartCode string

// todoStartCmd 开始待办
var todoStartCmd = &cobra.Command{
	Use:   "start",
	Short: "开始待办事项",
	Long:  `将指定待办事项标记为进行中~ 🚀`,
	Run: func(cmd *cobra.Command, args []string) {
		if todoStartCode == "" {
			cli.PrintError("请使用 --code 参数指定有效的待办标识码")
			os.Exit(1)
		}

		bs := startup.New(
			startup.WithSignalHandler(false),
		).MustInitialize(context.Background())
		defer bs.Shutdown()

		handler := handlers.NewTodoHandler(bs)
		if err := handler.Start(bs.Context(), todoStartCode); err != nil {
			cli.PrintError(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	todoStartCmd.Flags().StringVarP(&todoStartCode, "code", "c", "", "待办标识码（必填）")
	_ = todoStartCmd.MarkFlagRequired("code")

	todoCmd.AddCommand(todoStartCmd)
}
