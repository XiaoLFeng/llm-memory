package todo

import (
	"context"
	"os"
	"strings"

	"github.com/XiaoLFeng/llm-memory/internal/cli"
	"github.com/XiaoLFeng/llm-memory/internal/cli/handlers"
	"github.com/XiaoLFeng/llm-memory/startup"
	"github.com/spf13/cobra"
)

var todoBatchDeleteCodes string

// todoBatchDeleteCmd 批量删除待办
var todoBatchDeleteCmd = &cobra.Command{
	Use:   "batch-delete",
	Short: "批量删除待办事项",
	Long:  `批量删除多个待办事项~ 🗑️`,
	Example: `  # 批量删除多个待办
  llm-memory todo batch-delete --codes "todo-1,todo-2,todo-3"

  # 批量删除（使用空格分隔也会自动处理）
  llm-memory todo batch-delete --codes "todo-1, todo-2, todo-3"`,
	Run: func(cmd *cobra.Command, args []string) {
		if todoBatchDeleteCodes == "" {
			cli.PrintError("请使用 --codes 参数指定待办标识码列表（逗号分隔）")
			os.Exit(1)
		}

		// 解析 codes
		codes := splitCodes(todoBatchDeleteCodes)
		if len(codes) == 0 {
			cli.PrintError("未提供有效的待办标识码")
			os.Exit(1)
		}

		bs := startup.New(
			startup.WithSignalHandler(false),
		).MustInitialize(context.Background())
		defer bs.Shutdown()

		handler := handlers.NewTodoHandler(bs)
		if err := handler.BatchDelete(bs.Context(), codes); err != nil {
			cli.PrintError(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	todoBatchDeleteCmd.Flags().StringVar(&todoBatchDeleteCodes, "codes", "", "待办标识码列表（逗号分隔，必填）")
	_ = todoBatchDeleteCmd.MarkFlagRequired("codes")

	todoCmd.AddCommand(todoBatchDeleteCmd)
}

// splitCodes 解析逗号分隔的 codes 字符串
func splitCodes(codesStr string) []string {
	if codesStr == "" {
		return nil
	}

	parts := strings.Split(codesStr, ",")
	codes := make([]string, 0, len(parts))

	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			codes = append(codes, trimmed)
		}
	}

	return codes
}
