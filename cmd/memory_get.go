package cmd

import (
	"context"
	"os"

	"github.com/XiaoLFeng/llm-memory/internal/cli"
	"github.com/XiaoLFeng/llm-memory/internal/cli/handlers"
	"github.com/XiaoLFeng/llm-memory/startup"
	"github.com/spf13/cobra"
)

var memoryGetID int

// memoryGetCmd 获取记忆详情
// 呀~ 查看记忆的详细内容！📝
var memoryGetCmd = &cobra.Command{
	Use:   "get",
	Short: "获取记忆详情",
	Long:  `获取指定ID的记忆详细信息~ 📝`,
	Run: func(cmd *cobra.Command, args []string) {
		if memoryGetID <= 0 {
			cli.PrintError("请使用 --id 参数指定有效的记忆ID")
			os.Exit(1)
		}

		bs := startup.New(
			startup.WithSignalHandler(false),
		).MustInitialize(context.Background())
		defer bs.Shutdown()

		handler := handlers.NewMemoryHandler(bs)
		if err := handler.Get(bs.Context(), memoryGetID); err != nil {
			cli.PrintError(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	memoryGetCmd.Flags().IntVarP(&memoryGetID, "id", "i", 0, "记忆ID（必填）")
	_ = memoryGetCmd.MarkFlagRequired("id")

	memoryCmd.AddCommand(memoryGetCmd)
}
