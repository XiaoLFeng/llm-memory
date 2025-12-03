package group

import (
	"github.com/XiaoLFeng/llm-memory/cmd"
	"github.com/spf13/cobra"
)

// groupCmd 是 group 父命令
// 嘿嘿~ 组管理命令组！用于管理多路径关联的组~ 👥
var groupCmd = &cobra.Command{
	Use:   "group",
	Short: "组管理命令",
	Long: `管理 LLM-Memory 中的组，组可以包含多个路径，组内共享记忆、待办和计划~ ✨

示例：
  # 创建新组
  llm-memory group create my-project --desc "我的项目"

  # 列出所有组
  llm-memory group list

  # 将当前目录添加到组
  llm-memory group add-path my-project

  # 显示当前作用域
  llm-memory group current`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(groupCmd)
}
