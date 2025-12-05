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

var todoBatchCompleteCodes string

// todoBatchCompleteCmd 批量完成待办
var todoBatchCompleteCmd = &cobra.Command{
	Use:   "batch-complete",
	Short: "批量完成待办事项",
	Long:  `批量将多个待办事项标记为已完成~ 🎉`,
	Example: `  # 批量完成多个待办
  llm-memory todo batch-complete --codes "todo-1,todo-2,todo-3"

  # 批量完成（使用空格分隔也会自动处理）
  llm-memory todo batch-complete --codes "todo-1, todo-2, todo-3"`,
	Run: func(cmd *cobra.Command, args []string) {
		if todoBatchCompleteCodes == "" {
			cli.PrintError("请使用 --codes 参数指定待办标识码列表（逗号分隔）")
			os.Exit(1)
		}

		// 解析 codes
		codes := parseCodes(todoBatchCompleteCodes)
		if len(codes) == 0 {
			cli.PrintError("未提供有效的待办标识码")
			os.Exit(1)
		}

		bs := startup.New(
			startup.WithSignalHandler(false),
		).MustInitialize(context.Background())
		defer bs.Shutdown()

		handler := handlers.NewTodoHandler(bs)
		if err := handler.BatchComplete(bs.Context(), codes); err != nil {
			cli.PrintError(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	todoBatchCompleteCmd.Flags().StringVar(&todoBatchCompleteCodes, "codes", "", "待办标识码列表（逗号分隔，必填）")
	_ = todoBatchCompleteCmd.MarkFlagRequired("codes")

	todoCmd.AddCommand(todoBatchCompleteCmd)
}

// parseCodes 解析逗号分隔的 codes 字符串
func parseCodes(codesStr string) []string {
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
