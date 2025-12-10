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
	planUpdateCode        string
	planUpdateTitle       string
	planUpdateDescription string
	planUpdateContent     string
	planUpdateProgress    int
)

// planUpdateCmd 更新计划
var planUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "更新计划",
	Long:  `更新已有计划的标题、描述、内容或进度~ 📋`,
	Run: func(cmd *cobra.Command, args []string) {
		if planUpdateCode == "" {
			cli.PrintError("标识码不能为空，请使用 --code 参数")
			os.Exit(1)
		}

		// 检查是否至少提供一个更新字段
		hasTitle := cmd.Flags().Changed("title")
		hasDescription := cmd.Flags().Changed("description")
		hasContent := cmd.Flags().Changed("content")
		hasProgress := cmd.Flags().Changed("progress")

		if !hasTitle && !hasDescription && !hasContent && !hasProgress {
			cli.PrintError("至少需要提供一个更新字段（--title, --description, --content, --progress）")
			os.Exit(1)
		}

		bs := startup.New(
			startup.WithSignalHandler(false),
		).MustInitialize(context.Background())
		defer bs.Shutdown()

		// 构建更新参数
		var title, description, content *string
		var progress *int

		if hasTitle {
			title = &planUpdateTitle
		}
		if hasDescription {
			description = &planUpdateDescription
		}
		if hasContent {
			content = &planUpdateContent
		}
		if hasProgress {
			if planUpdateProgress < 0 || planUpdateProgress > 100 {
				cli.PrintError("进度必须在 0-100 之间")
				os.Exit(1)
			}
			progress = &planUpdateProgress
		}

		handler := handlers.NewPlanHandler(bs)
		if err := handler.Update(bs.Context(), planUpdateCode, title, description, content, progress); err != nil {
			cli.PrintError(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	planUpdateCmd.Flags().StringVarP(&planUpdateCode, "code", "c", "", "计划标识码（必填）")
	planUpdateCmd.Flags().StringVarP(&planUpdateTitle, "title", "t", "", "新标题")
	planUpdateCmd.Flags().StringVarP(&planUpdateDescription, "description", "d", "", "新描述")
	planUpdateCmd.Flags().StringVar(&planUpdateContent, "content", "", "新内容")
	planUpdateCmd.Flags().IntVarP(&planUpdateProgress, "progress", "p", -1, "新进度 0-100")

	_ = planUpdateCmd.MarkFlagRequired("code")

	planCmd.AddCommand(planUpdateCmd)
}
