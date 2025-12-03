package group

import (
	"context"
	"fmt"

	"github.com/XiaoLFeng/llm-memory/startup"
	"github.com/spf13/cobra"
)

// groupDeleteCmd 删除组的命令
// 呀~ 删除一个组（不会删除组内的数据）！💨
var groupDeleteCmd = &cobra.Command{
	Use:   "delete <name>",
	Short: "删除组",
	Long: `删除指定的组~ ✨

注意：删除组不会删除组内的记忆、待办和计划，只是解除路径关联。

示例：
  llm-memory group delete my-project`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		groupName := args[0]

		// 初始化 Bootstrap
		boot := startup.New()
		if err := boot.Initialize(context.Background()); err != nil {
			fmt.Printf("初始化失败: %v\n", err)
			return
		}
		defer func() { _ = boot.Shutdown() }()

		// 获取组
		group, err := boot.GroupService.GetGroupByName(boot.Context(), groupName)
		if err != nil {
			fmt.Printf("找不到组 '%s': %v\n", groupName, err)
			return
		}

		// 删除组
		if err := boot.GroupService.DeleteGroup(boot.Context(), group.ID); err != nil {
			fmt.Printf("删除组失败: %v\n", err)
			return
		}

		fmt.Printf("✨ 组 '%s' 已删除！\n", groupName)
	},
}

func init() {
	groupCmd.AddCommand(groupDeleteCmd)
}
