package memory

import (
	"context"
	"os"
	"strings"

	"github.com/XiaoLFeng/llm-memory/internal/cli"
	"github.com/XiaoLFeng/llm-memory/internal/cli/handlers"
	"github.com/XiaoLFeng/llm-memory/startup"
	"github.com/spf13/cobra"
)

var (
	updateCode     string
	updateTitle    string
	updateContent  string
	updateCategory string
	updateTags     string
	updatePriority int
)

// memoryUpdateCmd 更新记忆
var memoryUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "更新记忆",
	Long:  `更新已有记忆的标题、内容、分类、标签或优先级~ ✨`,
	Run: func(cmd *cobra.Command, args []string) {
		if updateCode == "" {
			cli.PrintError("标识码不能为空，请使用 --code 参数")
			os.Exit(1)
		}

		// 检查是否至少提供一个更新字段
		hasTitle := cmd.Flags().Changed("title")
		hasContent := cmd.Flags().Changed("content")
		hasCategory := cmd.Flags().Changed("category")
		hasTags := cmd.Flags().Changed("tags")
		hasPriority := cmd.Flags().Changed("priority")

		if !hasTitle && !hasContent && !hasCategory && !hasTags && !hasPriority {
			cli.PrintError("至少需要提供一个更新字段（--title, --content, --category, --tags, --priority）")
			os.Exit(1)
		}

		bs := startup.New(
			startup.WithSignalHandler(false),
		).MustInitialize(context.Background())
		defer bs.Shutdown()

		// 构建更新参数
		var title, content, category *string
		var tags *[]string
		var priority *int

		if hasTitle {
			title = &updateTitle
		}
		if hasContent {
			content = &updateContent
		}
		if hasCategory {
			category = &updateCategory
		}
		if hasTags {
			tagList := strings.Split(updateTags, ",")
			for i := range tagList {
				tagList[i] = strings.TrimSpace(tagList[i])
			}
			tags = &tagList
		}
		if hasPriority {
			if updatePriority < 1 || updatePriority > 4 {
				cli.PrintError("优先级必须在 1-4 之间")
				os.Exit(1)
			}
			priority = &updatePriority
		}

		handler := handlers.NewMemoryHandler(bs)
		if err := handler.Update(bs.Context(), updateCode, title, content, category, tags, priority); err != nil {
			cli.PrintError(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	memoryUpdateCmd.Flags().StringVarP(&updateCode, "code", "c", "", "记忆标识码（必填）")
	memoryUpdateCmd.Flags().StringVarP(&updateTitle, "title", "t", "", "新标题")
	memoryUpdateCmd.Flags().StringVar(&updateContent, "content", "", "新内容")
	memoryUpdateCmd.Flags().StringVarP(&updateCategory, "category", "C", "", "新分类")
	memoryUpdateCmd.Flags().StringVar(&updateTags, "tags", "", "新标签（逗号分隔）")
	memoryUpdateCmd.Flags().IntVarP(&updatePriority, "priority", "p", 0, "新优先级 1-4")

	_ = memoryUpdateCmd.MarkFlagRequired("code")

	memoryCmd.AddCommand(memoryUpdateCmd)
}
