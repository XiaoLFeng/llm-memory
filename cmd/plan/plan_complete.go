package plan

import (
	"context"
	"os"

	"github.com/XiaoLFeng/llm-memory/internal/cli"
	"github.com/XiaoLFeng/llm-memory/internal/cli/handlers"
	"github.com/XiaoLFeng/llm-memory/startup"
	"github.com/spf13/cobra"
)

var planCompleteID int

// planCompleteCmd 完成计划
// 呀~ 标记计划为已完成！🎉
var planCompleteCmd = &cobra.Command{
	Use:   "complete",
	Short: "完成计划",
	Long:  `将指定计划标记为已完成~ 🎉`,
	Run: func(cmd *cobra.Command, args []string) {
		if planCompleteID <= 0 {
			cli.PrintError("请使用 --id 参数指定有效的计划ID")
			os.Exit(1)
		}

		bs := startup.New(
			startup.WithSignalHandler(false),
		).MustInitialize(context.Background())
		defer bs.Shutdown()

		handler := handlers.NewPlanHandler(bs)
		if err := handler.Complete(bs.Context(), uint(planCompleteID)); err != nil {
			cli.PrintError(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	planCompleteCmd.Flags().IntVarP(&planCompleteID, "id", "i", 0, "计划ID（必填）")
	_ = planCompleteCmd.MarkFlagRequired("id")

	planCmd.AddCommand(planCompleteCmd)
}
