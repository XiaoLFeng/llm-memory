package plan

import (
	"context"
	"os"

	"github.com/XiaoLFeng/llm-memory/internal/cli"
	"github.com/XiaoLFeng/llm-memory/internal/cli/handlers"
	"github.com/XiaoLFeng/llm-memory/startup"
	"github.com/spf13/cobra"
)

var planGetCode string

// planGetCmd 获取计划详情
var planGetCmd = &cobra.Command{
	Use:   "get",
	Short: "获取计划详情",
	Long:  `获取指定计划的详细信息，包括标题、描述、进度、内容等~ 📋`,
	Run: func(cmd *cobra.Command, args []string) {
		if planGetCode == "" {
			cli.PrintError("标识码不能为空，请使用 --code 参数")
			os.Exit(1)
		}

		bs := startup.New(
			startup.WithSignalHandler(false),
		).MustInitialize(context.Background())
		defer bs.Shutdown()

		handler := handlers.NewPlanHandler(bs)
		if err := handler.Get(bs.Context(), planGetCode); err != nil {
			cli.PrintError(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	planGetCmd.Flags().StringVarP(&planGetCode, "code", "c", "", "计划标识码（必填）")
	_ = planGetCmd.MarkFlagRequired("code")
	planCmd.AddCommand(planGetCmd)
}
