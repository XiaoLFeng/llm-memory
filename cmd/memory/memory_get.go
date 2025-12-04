package memory

import (
	"context"
	"os"

	"github.com/XiaoLFeng/llm-memory/internal/cli"
	"github.com/XiaoLFeng/llm-memory/internal/cli/handlers"
	"github.com/XiaoLFeng/llm-memory/startup"
	"github.com/spf13/cobra"
)

var memoryGetCode string

// memoryGetCmd 获取记忆详情
var memoryGetCmd = &cobra.Command{
	Use:   "get",
	Short: "获取记忆详情",
	Long:  `获取指定标识码的记忆详细信息~ 📝`,
	Run: func(cmd *cobra.Command, args []string) {
		if memoryGetCode == "" {
			cli.PrintError("请使用 --code 参数指定有效的记忆标识码")
			os.Exit(1)
		}

		bs := startup.New(
			startup.WithSignalHandler(false),
		).MustInitialize(context.Background())
		defer bs.Shutdown()

		handler := handlers.NewMemoryHandler(bs)
		if err := handler.Get(bs.Context(), memoryGetCode); err != nil {
			cli.PrintError(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	memoryGetCmd.Flags().StringVarP(&memoryGetCode, "code", "c", "", "记忆标识码（必填）")
	_ = memoryGetCmd.MarkFlagRequired("code")

	memoryCmd.AddCommand(memoryGetCmd)
}
