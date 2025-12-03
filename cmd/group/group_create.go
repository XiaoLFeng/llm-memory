package group

import (
	"context"
	"fmt"

	"github.com/XiaoLFeng/llm-memory/startup"
	"github.com/spf13/cobra"
)

// groupCreateCmd 创建新组的命令
// 嘿嘿~ 创建一个新的组来管理多个路径！💖
var groupCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "创建新组",
	Long: `创建一个新的组，组可以包含多个路径，组内共享记忆、待办和计划~ ✨

示例：
  llm-memory group create my-project
  llm-memory group create my-project --desc "我的项目描述"`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		description, _ := cmd.Flags().GetString("desc")

		// 初始化 Bootstrap
		boot := startup.New()
		if err := boot.Initialize(context.Background()); err != nil {
			fmt.Printf("初始化失败: %v\n", err)
			return
		}
		defer func() { _ = boot.Shutdown() }()

		// 创建组
		group, err := boot.GroupService.CreateGroup(boot.Context(), name, description)
		if err != nil {
			fmt.Printf("创建组失败: %v\n", err)
			return
		}

		fmt.Printf("✨ 组 '%s' 创建成功！(ID: %d)\n", group.Name, group.ID)
		if description != "" {
			fmt.Printf("   描述: %s\n", description)
		}
	},
}

func init() {
	groupCmd.AddCommand(groupCreateCmd)
	groupCreateCmd.Flags().String("desc", "", "组的描述信息")
}
