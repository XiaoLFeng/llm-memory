package todo

import (
	"context"
	"os"

	"github.com/XiaoLFeng/llm-memory/internal/cli"
	"github.com/XiaoLFeng/llm-memory/internal/cli/handlers"
	"github.com/XiaoLFeng/llm-memory/startup"
	"github.com/spf13/cobra"
)

var todoGetCode string

// todoGetCmd 获取待办详情
var todoGetCmd = &cobra.Command{
	Use:   "get",
	Short: "获取待办详情",
	Long:  `获取指定待办的详细信息，包括标题、描述、状态、优先级等~ 📋`,
	Run: func(cmd *cobra.Command, args []string) {
		if todoGetCode == "" {
			cli.PrintError("标识码不能为空，请使用 --code 参数")
			os.Exit(1)
		}

		bs := startup.New(
			startup.WithSignalHandler(false),
		).MustInitialize(context.Background())
		defer bs.Shutdown()

		handler := handlers.NewTodoHandler(bs)
		if err := handler.Get(bs.Context(), todoGetCode); err != nil {
			cli.PrintError(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	todoGetCmd.Flags().StringVarP(&todoGetCode, "code", "c", "", "待办标识码（必填）")
	_ = todoGetCmd.MarkFlagRequired("code")
	todoCmd.AddCommand(todoGetCmd)
}
