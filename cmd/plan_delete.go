package cmd

import (
	"context"
	"os"

	"github.com/XiaoLFeng/llm-memory/internal/cli"
	"github.com/XiaoLFeng/llm-memory/internal/cli/handlers"
	"github.com/XiaoLFeng/llm-memory/startup"
	"github.com/spf13/cobra"
)

var planDeleteID int

// planDeleteCmd 删除计划
// 嘿嘿~ 删除指定的计划！🗑️
var planDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "删除计划",
	Long:  `删除指定ID的计划~ 🗑️`,
	Run: func(cmd *cobra.Command, args []string) {
		if planDeleteID <= 0 {
			cli.PrintError("请使用 --id 参数指定有效的计划ID")
			os.Exit(1)
		}

		bs := startup.New(
			startup.WithSignalHandler(false),
		).MustInitialize(context.Background())
		defer bs.Shutdown()

		handler := handlers.NewPlanHandler(bs)
		if err := handler.Delete(bs.Context(), planDeleteID); err != nil {
			cli.PrintError(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	planDeleteCmd.Flags().IntVarP(&planDeleteID, "id", "i", 0, "计划ID（必填）")
	_ = planDeleteCmd.MarkFlagRequired("id")

	planCmd.AddCommand(planDeleteCmd)
}
