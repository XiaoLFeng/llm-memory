package plan

import (
	"context"
	"os"

	"github.com/XiaoLFeng/llm-memory/internal/cli"
	"github.com/XiaoLFeng/llm-memory/internal/cli/handlers"
	"github.com/XiaoLFeng/llm-memory/startup"
	"github.com/spf13/cobra"
)

var planDeleteCode string

// planDeleteCmd 删除计划
var planDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "删除计划",
	Long:  `删除指定标识码的计划~ 🗑️`,
	Run: func(cmd *cobra.Command, args []string) {
		if planDeleteCode == "" {
			cli.PrintError("请使用 --code 参数指定有效的计划标识码")
			os.Exit(1)
		}

		bs := startup.New(
			startup.WithSignalHandler(false),
		).MustInitialize(context.Background())
		defer bs.Shutdown()

		handler := handlers.NewPlanHandler(bs)
		if err := handler.Delete(bs.Context(), planDeleteCode); err != nil {
			cli.PrintError(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	planDeleteCmd.Flags().StringVarP(&planDeleteCode, "code", "c", "", "计划标识码（必填）")
	_ = planDeleteCmd.MarkFlagRequired("code")

	planCmd.AddCommand(planDeleteCmd)
}
