package cmd

import (
	"github.com/spf13/cobra"
)

// todoCmd 是 todo 父命令
// 嘿嘿~ 待办事项管理命令组！✅
var todoCmd = &cobra.Command{
	Use:   "todo",
	Short: "待办事项管理命令",
	Long:  `管理 LLM-Memory 中的待办事项~ ✅`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(todoCmd)
}
