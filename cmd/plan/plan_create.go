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
	planTitle       string
	planDescription string
)

// planCreateCmd 创建新计划
// 嘿嘿~ 创建一个新计划！💫
var planCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "创建新计划",
	Long:  `创建一个新的计划~ ✨`,
	Run: func(cmd *cobra.Command, args []string) {
		if planTitle == "" {
			cli.PrintError("标题不能为空，请使用 --title 参数")
			os.Exit(1)
		}

		bs := startup.New(
			startup.WithSignalHandler(false),
		).MustInitialize(context.Background())
		defer bs.Shutdown()

		handler := handlers.NewPlanHandler(bs)
		if err := handler.Create(bs.Context(), planTitle, planDescription); err != nil {
			cli.PrintError(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	planCreateCmd.Flags().StringVarP(&planTitle, "title", "t", "", "计划标题（必填）")
	planCreateCmd.Flags().StringVarP(&planDescription, "description", "d", "", "计划描述")

	_ = planCreateCmd.MarkFlagRequired("title")

	planCmd.AddCommand(planCreateCmd)
}
