package todo

import (
	"context"
	"os"

	"github.com/XiaoLFeng/llm-memory/internal/cli"
	"github.com/XiaoLFeng/llm-memory/internal/cli/handlers"
	"github.com/XiaoLFeng/llm-memory/startup"
	"github.com/spf13/cobra"
)

var todoCancelCode string

// todoCancelCmd 取消待办
var todoCancelCmd = &cobra.Command{
	Use:   "cancel",
	Short: "取消待办",
	Long:  `将指定待办标记为已取消状态~ ❌`,
	Run: func(cmd *cobra.Command, args []string) {
		if todoCancelCode == "" {
			cli.PrintError("标识码不能为空，请使用 --code 参数")
			os.Exit(1)
		}

		bs := startup.New(
			startup.WithSignalHandler(false),
		).MustInitialize(context.Background())
		defer bs.Shutdown()

		handler := handlers.NewTodoHandler(bs)
		if err := handler.Cancel(bs.Context(), todoCancelCode); err != nil {
			cli.PrintError(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	todoCancelCmd.Flags().StringVarP(&todoCancelCode, "code", "c", "", "待办标识码（必填）")
	_ = todoCancelCmd.MarkFlagRequired("code")
	todoCmd.AddCommand(todoCancelCmd)
}
