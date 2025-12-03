package plan

import (
	"context"
	"os"

	"github.com/XiaoLFeng/llm-memory/internal/cli"
	"github.com/XiaoLFeng/llm-memory/internal/cli/handlers"
	"github.com/XiaoLFeng/llm-memory/startup"
	"github.com/spf13/cobra"
)

var (
	planProgressID    int
	planProgressValue int
)

// planProgressCmd 更新计划进度
var planProgressCmd = &cobra.Command{
	Use:   "progress",
	Short: "更新计划进度",
	Long:  `更新指定计划的完成进度（0-100）~ 📊`,
	Run: func(cmd *cobra.Command, args []string) {
		if planProgressID <= 0 {
			cli.PrintError("请使用 --id 参数指定有效的计划ID")
			os.Exit(1)
		}
		if planProgressValue < 0 || planProgressValue > 100 {
			cli.PrintError("进度必须在 0-100 之间")
			os.Exit(1)
		}

		bs := startup.New(
			startup.WithSignalHandler(false),
		).MustInitialize(context.Background())
		defer bs.Shutdown()

		handler := handlers.NewPlanHandler(bs)
		if err := handler.UpdateProgress(bs.Context(), int64(planProgressID), planProgressValue); err != nil {
			cli.PrintError(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	planProgressCmd.Flags().IntVarP(&planProgressID, "id", "i", 0, "计划ID（必填）")
	planProgressCmd.Flags().IntVarP(&planProgressValue, "progress", "p", 0, "进度值 0-100（必填）")

	_ = planProgressCmd.MarkFlagRequired("id")
	_ = planProgressCmd.MarkFlagRequired("progress")

	planCmd.AddCommand(planProgressCmd)
}
