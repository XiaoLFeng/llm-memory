package plan

import (
	"context"
	"os"

	"github.com/XiaoLFeng/llm-memory/internal/cli"
	"github.com/XiaoLFeng/llm-memory/internal/cli/handlers"
	"github.com/XiaoLFeng/llm-memory/startup"
	"github.com/spf13/cobra"
)

var planStartCode string

// planStartCmd 开始计划
var planStartCmd = &cobra.Command{
	Use:   "start",
	Short: "开始计划",
	Long:  `将指定计划标记为进行中~ 🚀`,
	Run: func(cmd *cobra.Command, args []string) {
		if planStartCode == "" {
			cli.PrintError("请使用 --code 参数指定有效的计划标识码")
			os.Exit(1)
		}

		bs := startup.New(
			startup.WithSignalHandler(false),
		).MustInitialize(context.Background())
		defer bs.Shutdown()

		handler := handlers.NewPlanHandler(bs)
		if err := handler.Start(bs.Context(), planStartCode); err != nil {
			cli.PrintError(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	planStartCmd.Flags().StringVarP(&planStartCode, "code", "c", "", "计划标识码（必填）")
	_ = planStartCmd.MarkFlagRequired("code")

	planCmd.AddCommand(planStartCmd)
}
