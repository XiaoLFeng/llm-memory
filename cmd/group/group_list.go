package group

import (
	"context"
	"fmt"

	"github.com/XiaoLFeng/llm-memory/startup"
	"github.com/spf13/cobra"
)

// groupListCmd 列出所有组的命令
// 呀~ 看看我们有哪些组吧！📋
var groupListCmd = &cobra.Command{
	Use:   "list",
	Short: "列出所有组",
	Long:  `列出所有已创建的组及其包含的路径~ ✨`,
	Run: func(cmd *cobra.Command, args []string) {
		// 初始化 Bootstrap
		boot := startup.New()
		if err := boot.Initialize(context.Background()); err != nil {
			fmt.Printf("初始化失败: %v\n", err)
			return
		}
		defer func() { _ = boot.Shutdown() }()

		// 获取所有组
		groups, err := boot.GroupService.ListGroups(boot.Context())
		if err != nil {
			fmt.Printf("获取组列表失败: %v\n", err)
			return
		}

		if len(groups) == 0 {
			fmt.Println(iconInbox + " 暂无任何组，使用 'llm-memory group create <name>' 创建一个吧~")
			return
		}

		fmt.Println(iconPackage + " 组列表:")
		fmt.Println("─────────────────────────────────────")
		for _, group := range groups {
			fmt.Printf("\n"+iconTag+"  [%d] %s\n", group.ID, group.Name)
			if group.Description != "" {
				fmt.Printf("   "+iconEdit+" 描述: %s\n", group.Description)
			}
			if len(group.Paths) > 0 {
				fmt.Printf("   "+iconFolder+" 路径 (%d):\n", len(group.Paths))
				for _, path := range group.Paths {
					fmt.Printf("      - %s\n", path.GetPath())
				}
			} else {
				fmt.Println("   " + iconFolder + " 暂无关联路径")
			}
		}
		fmt.Println("\n─────────────────────────────────────")
		fmt.Printf("共 %d 个组\n", len(groups))
	},
}

func init() {
	groupCmd.AddCommand(groupListCmd)
}
