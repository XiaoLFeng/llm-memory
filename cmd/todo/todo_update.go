package todo

import (
	"context"
	"os"

	"github.com/XiaoLFeng/llm-memory/internal/cli"
	"github.com/XiaoLFeng/llm-memory/internal/cli/handlers"
	"github.com/XiaoLFeng/llm-memory/startup"
	"github.com/spf13/cobra"
)

var (
	todoUpdateCode        string
	todoUpdateTitle       string
	todoUpdateDescription string
	todoUpdatePriority    int
	todoUpdateStatus      int
)

// todoUpdateCmd 更新待办
var todoUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "更新待办",
	Long:  `更新已有待办的标题、描述、优先级或状态~ 📝`,
	Run: func(cmd *cobra.Command, args []string) {
		if todoUpdateCode == "" {
			cli.PrintError("标识码不能为空，请使用 --code 参数")
			os.Exit(1)
		}

		// 检查是否至少提供一个更新字段
		hasTitle := cmd.Flags().Changed("title")
		hasDescription := cmd.Flags().Changed("description")
		hasPriority := cmd.Flags().Changed("priority")
		hasStatus := cmd.Flags().Changed("status")

		if !hasTitle && !hasDescription && !hasPriority && !hasStatus {
			cli.PrintError("至少需要提供一个更新字段（--title, --description, --priority, --status）")
			os.Exit(1)
		}

		bs := startup.New(
			startup.WithSignalHandler(false),
		).MustInitialize(context.Background())
		defer bs.Shutdown()

		// 构建更新参数
		var title, description *string
		var priority, status *int

		if hasTitle {
			title = &todoUpdateTitle
		}
		if hasDescription {
			description = &todoUpdateDescription
		}
		if hasPriority {
			if todoUpdatePriority < 1 || todoUpdatePriority > 4 {
				cli.PrintError("优先级必须在 1-4 之间")
				os.Exit(1)
			}
			priority = &todoUpdatePriority
		}
		if hasStatus {
			if todoUpdateStatus < 0 || todoUpdateStatus > 3 {
				cli.PrintError("状态必须在 0-3 之间（0待处理/1进行中/2已完成/3已取消）")
				os.Exit(1)
			}
			status = &todoUpdateStatus
		}

		handler := handlers.NewTodoHandler(bs)
		if err := handler.Update(bs.Context(), todoUpdateCode, title, description, priority, status); err != nil {
			cli.PrintError(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	todoUpdateCmd.Flags().StringVarP(&todoUpdateCode, "code", "c", "", "待办标识码（必填）")
	todoUpdateCmd.Flags().StringVarP(&todoUpdateTitle, "title", "t", "", "新标题")
	todoUpdateCmd.Flags().StringVarP(&todoUpdateDescription, "description", "d", "", "新描述")
	todoUpdateCmd.Flags().IntVarP(&todoUpdatePriority, "priority", "p", 0, "新优先级 1-4")
	todoUpdateCmd.Flags().IntVarP(&todoUpdateStatus, "status", "s", -1, "新状态 0-3（0待处理/1进行中/2已完成/3已取消）")

	_ = todoUpdateCmd.MarkFlagRequired("code")

	todoCmd.AddCommand(todoUpdateCmd)
}
