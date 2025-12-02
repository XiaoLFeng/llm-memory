package cmd

import (
	"context"
	"os"

	"github.com/XiaoLFeng/llm-memory/internal/cli"
	"github.com/XiaoLFeng/llm-memory/internal/cli/handlers"
	"github.com/XiaoLFeng/llm-memory/startup"
	"github.com/spf13/cobra"
)

var planStartID int

// planStartCmd 开始计划
// 嘿嘿~ 开始执行计划！🚀
var planStartCmd = &cobra.Command{
	Use:   "start",
	Short: "开始计划",
	Long:  `将指定计划标记为进行中~ 🚀`,
	Run: func(cmd *cobra.Command, args []string) {
		if planStartID <= 0 {
			cli.PrintError("请使用 --id 参数指定有效的计划ID")
			os.Exit(1)
		}

		bs := startup.New(
			startup.WithSignalHandler(false),
		).MustInitialize(context.Background())
		defer bs.Shutdown()

		handler := handlers.NewPlanHandler(bs)
		if err := handler.Start(bs.Context(), planStartID); err != nil {
			cli.PrintError(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	planStartCmd.Flags().IntVarP(&planStartID, "id", "i", 0, "计划ID（必填）")
	_ = planStartCmd.MarkFlagRequired("id")

	planCmd.AddCommand(planStartCmd)
}
