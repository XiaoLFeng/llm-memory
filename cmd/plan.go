package cmd

import (
	"github.com/spf13/cobra"
)

// planCmd 是 plan 父命令
// 嘿嘿~ 计划管理命令组！📋
var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "计划管理命令",
	Long:  `管理 LLM-Memory 中的计划，包括创建、更新进度、完成等~ 📋`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(planCmd)
}
