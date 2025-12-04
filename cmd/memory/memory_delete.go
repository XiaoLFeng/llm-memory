package memory

import (
	"context"
	"os"

	"github.com/XiaoLFeng/llm-memory/internal/cli"
	"github.com/XiaoLFeng/llm-memory/internal/cli/handlers"
	"github.com/XiaoLFeng/llm-memory/startup"
	"github.com/spf13/cobra"
)

var memoryDeleteCode string

// memoryDeleteCmd 删除记忆
var memoryDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "删除记忆",
	Long:  `删除指定标识码的记忆条目~ 🗑️`,
	Run: func(cmd *cobra.Command, args []string) {
		if memoryDeleteCode == "" {
			cli.PrintError("请使用 --code 参数指定有效的记忆标识码")
			os.Exit(1)
		}

		bs := startup.New(
			startup.WithSignalHandler(false),
		).MustInitialize(context.Background())
		defer bs.Shutdown()

		handler := handlers.NewMemoryHandler(bs)
		if err := handler.Delete(bs.Context(), memoryDeleteCode); err != nil {
			cli.PrintError(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	memoryDeleteCmd.Flags().StringVarP(&memoryDeleteCode, "code", "c", "", "记忆标识码（必填）")
	_ = memoryDeleteCmd.MarkFlagRequired("code")

	memoryCmd.AddCommand(memoryDeleteCmd)
}
