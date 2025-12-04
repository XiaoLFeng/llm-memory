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
	todoCode        string
	todoTitle       string
	todoDescription string
	todoPriority    int
	todoGlobal      bool
)

// todoCreateCmd 创建待办
// 呀~ 创建新的待办事项！💫
var todoCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "创建新待办事项",
	Long:  `创建一个新的待办事项~ ✨`,
	Run: func(cmd *cobra.Command, args []string) {
		if todoCode == "" {
			cli.PrintError("标识码不能为空，请使用 --code 参数")
			os.Exit(1)
		}
		if todoTitle == "" {
			cli.PrintError("标题不能为空，请使用 --title 参数")
			os.Exit(1)
		}

		bs := startup.New(
			startup.WithSignalHandler(false),
		).MustInitialize(context.Background())
		defer bs.Shutdown()

		handler := handlers.NewTodoHandler(bs)
		if err := handler.Create(bs.Context(), todoCode, todoTitle, todoDescription, todoPriority, todoGlobal); err != nil {
			cli.PrintError(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	todoCreateCmd.Flags().StringVarP(&todoCode, "code", "c", "", "待办标识码（必填）")
	todoCreateCmd.Flags().StringVarP(&todoTitle, "title", "t", "", "待办标题（必填）")
	todoCreateCmd.Flags().StringVarP(&todoDescription, "description", "d", "", "待办描述")
	todoCreateCmd.Flags().IntVarP(&todoPriority, "priority", "p", 2, "优先级：1低/2中/3高/4紧急")
	todoCreateCmd.Flags().BoolVar(&todoGlobal, "global", false, "将待办保存为全局（默认当前路径/组内可见）")

	_ = todoCreateCmd.MarkFlagRequired("code")
	_ = todoCreateCmd.MarkFlagRequired("title")

	todoCmd.AddCommand(todoCreateCmd)
}
