package todo

import (
	"context"
	"os"

	"github.com/XiaoLFeng/llm-memory/internal/cli"
	"github.com/XiaoLFeng/llm-memory/internal/cli/handlers"
	"github.com/XiaoLFeng/llm-memory/startup"
	"github.com/spf13/cobra"
)

// todoFinalCmd 删除所有待办
var todoFinalCmd = &cobra.Command{
	Use:   "final",
	Short: "删除当前作用域的所有待办",
	Long:  `删除当前作用域内的所有待办事项（不可恢复）~ 🗑️`,
	Run: func(cmd *cobra.Command, args []string) {
		bs := startup.New(
			startup.WithSignalHandler(false),
		).MustInitialize(context.Background())
		defer bs.Shutdown()

		handler := handlers.NewTodoHandler(bs)
		if err := handler.Final(bs.Context()); err != nil {
			cli.PrintError(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	todoCmd.AddCommand(todoFinalCmd)
}
