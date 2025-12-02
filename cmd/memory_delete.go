package cmd

import (
	"context"
	"os"

	"github.com/XiaoLFeng/llm-memory/internal/cli"
	"github.com/XiaoLFeng/llm-memory/internal/cli/handlers"
	"github.com/XiaoLFeng/llm-memory/startup"
	"github.com/spf13/cobra"
)

var memoryDeleteID int

// memoryDeleteCmd 删除记忆
// 嘿嘿~ 删除指定的记忆！🗑️
var memoryDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "删除记忆",
	Long:  `删除指定ID的记忆条目~ 🗑️`,
	Run: func(cmd *cobra.Command, args []string) {
		if memoryDeleteID <= 0 {
			cli.PrintError("请使用 --id 参数指定有效的记忆ID")
			os.Exit(1)
		}

		bs := startup.New(
			startup.WithSignalHandler(false),
		).MustInitialize(context.Background())
		defer bs.Shutdown()

		handler := handlers.NewMemoryHandler(bs)
		if err := handler.Delete(bs.Context(), memoryDeleteID); err != nil {
			cli.PrintError(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	memoryDeleteCmd.Flags().IntVarP(&memoryDeleteID, "id", "i", 0, "记忆ID（必填）")
	_ = memoryDeleteCmd.MarkFlagRequired("id")

	memoryCmd.AddCommand(memoryDeleteCmd)
}
