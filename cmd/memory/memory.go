package memory

import (
	"github.com/XiaoLFeng/llm-memory/cmd"
	"github.com/spf13/cobra"
)

// memoryCmd 是 memory 父命令
// 嘿嘿~ 记忆管理命令组！📚
var memoryCmd = &cobra.Command{
	Use:   "memory",
	Short: "记忆管理命令",
	Long:  `管理 LLM-Memory 中的记忆条目，包括创建、查看、搜索和删除~ ✨`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(memoryCmd)
}
