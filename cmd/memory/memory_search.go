package memory

import (
	"context"
	"os"

	"github.com/XiaoLFeng/llm-memory/internal/cli"
	"github.com/XiaoLFeng/llm-memory/internal/cli/handlers"
	"github.com/XiaoLFeng/llm-memory/startup"
	"github.com/spf13/cobra"
)

var memorySearchKeyword string

// memorySearchCmd 搜索记忆
// 呀~ 根据关键词搜索记忆！🔍
var memorySearchCmd = &cobra.Command{
	Use:   "search",
	Short: "搜索记忆",
	Long:  `根据关键词搜索记忆条目~ 🔍`,
	Run: func(cmd *cobra.Command, args []string) {
		if memorySearchKeyword == "" {
			cli.PrintError("请使用 --keyword 参数指定搜索关键词")
			os.Exit(1)
		}

		bs := startup.New(
			startup.WithSignalHandler(false),
		).MustInitialize(context.Background())
		defer bs.Shutdown()

		handler := handlers.NewMemoryHandler(bs)
		if err := handler.Search(bs.Context(), memorySearchKeyword); err != nil {
			cli.PrintError(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	memorySearchCmd.Flags().StringVarP(&memorySearchKeyword, "keyword", "k", "", "搜索关键词（必填）")
	_ = memorySearchCmd.MarkFlagRequired("keyword")

	memoryCmd.AddCommand(memorySearchCmd)
}
