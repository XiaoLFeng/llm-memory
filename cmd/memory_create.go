package cmd

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
	memoryTitle    string
	memoryContent  string
	memoryCategory string
	memoryTags     string
)

// memoryCreateCmd 创建新记忆
// 嘿嘿~ 创建一条新的记忆！💫
var memoryCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "创建新记忆",
	Long:  `创建一条新的记忆条目~ ✨`,
	Run: func(cmd *cobra.Command, args []string) {
		if memoryTitle == "" {
			cli.PrintError("标题不能为空，请使用 --title 参数")
			os.Exit(1)
		}
		if memoryContent == "" {
			cli.PrintError("内容不能为空，请使用 --content 参数")
			os.Exit(1)
		}

		bs := startup.New(
			startup.WithSignalHandler(false),
		).MustInitialize(context.Background())
		defer bs.Shutdown()

		// 处理标签
		var tags []string
		if memoryTags != "" {
			tags = strings.Split(memoryTags, ",")
			for i := range tags {
				tags[i] = strings.TrimSpace(tags[i])
			}
		}

		handler := handlers.NewMemoryHandler(bs)
		if err := handler.Create(bs.Context(), memoryTitle, memoryContent, memoryCategory, tags); err != nil {
			cli.PrintError(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	memoryCreateCmd.Flags().StringVarP(&memoryTitle, "title", "t", "", "记忆标题（必填）")
	memoryCreateCmd.Flags().StringVarP(&memoryContent, "content", "c", "", "记忆内容（必填）")
	memoryCreateCmd.Flags().StringVarP(&memoryCategory, "category", "C", "默认", "记忆分类")
	memoryCreateCmd.Flags().StringVar(&memoryTags, "tags", "", "标签（逗号分隔）")

	_ = memoryCreateCmd.MarkFlagRequired("title")
	_ = memoryCreateCmd.MarkFlagRequired("content")

	memoryCmd.AddCommand(memoryCreateCmd)
}
