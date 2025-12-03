package group

import (
	"context"
	"fmt"
	"os"

	"github.com/XiaoLFeng/llm-memory/startup"
	"github.com/spf13/cobra"
)

// groupRemovePathCmd 从组中移除路径的命令
// 呀~ 把路径从组中移除！🗑️
var groupRemovePathCmd = &cobra.Command{
	Use:   "remove-path <group-name> [path]",
	Short: "从组中移除路径",
	Long: `从组中移除当前目录或指定路径~ ✨

如果不指定路径，则默认移除当前工作目录。

示例：
  llm-memory group remove-path my-project           # 移除当前目录
  llm-memory group remove-path my-project /path/to  # 移除指定路径`,
	Args: cobra.RangeArgs(1, 2),
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

		// 确定要移除的路径
		var pathToRemove string
		if len(args) > 1 {
			pathToRemove = args[1]
		} else {
			// 使用当前目录
			pwd, err := os.Getwd()
			if err != nil {
				fmt.Printf("无法获取当前目录: %v\n", err)
				return
			}
			pathToRemove = pwd
		}

		// 从组中移除路径
		if err := boot.GroupService.RemovePath(boot.Context(), group.ID, pathToRemove); err != nil {
			fmt.Printf("移除路径失败: %v\n", err)
			return
		}

		fmt.Printf("✨ 已从组 '%s' 中移除路径 '%s'！\n", groupName, pathToRemove)
	},
}

func init() {
	groupCmd.AddCommand(groupRemovePathCmd)
}
