package plan

import (
	"context"
	"os"

	"github.com/XiaoLFeng/llm-memory/internal/cli"
	"github.com/XiaoLFeng/llm-memory/internal/cli/handlers"
	"github.com/XiaoLFeng/llm-memory/startup"
	"github.com/spf13/cobra"
)

var planCompleteCode string

// planCompleteCmd 完成计划
var planCompleteCmd = &cobra.Command{
	Use:   "complete",
	Short: "完成计划",
	Long:  `将指定计划标记为已完成~ 🎉`,
	Run: func(cmd *cobra.Command, args []string) {
		if planCompleteCode == "" {
			cli.PrintError("请使用 --code 参数指定有效的计划标识码")
			os.Exit(1)
		}

		bs := startup.New(
			startup.WithSignalHandler(false),
		).MustInitialize(context.Background())
		defer bs.Shutdown()

		handler := handlers.NewPlanHandler(bs)
		if err := handler.Complete(bs.Context(), planCompleteCode); err != nil {
			cli.PrintError(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	planCompleteCmd.Flags().StringVarP(&planCompleteCode, "code", "c", "", "计划标识码（必填）")
	_ = planCompleteCmd.MarkFlagRequired("code")

	planCmd.AddCommand(planCompleteCmd)
}
