package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/XiaoLFeng/llm-memory/startup"
	"github.com/spf13/cobra"
)

// groupAddPathCmd 将路径添加到组的命令
// 嘿嘿~ 把当前目录或指定路径加入组！📍
var groupAddPathCmd = &cobra.Command{
	Use:   "add-path <group-name> [path]",
	Short: "将路径添加到组",
	Long: `将当前目录或指定路径添加到组中~ ✨

如果不指定路径，则默认添加当前工作目录。

示例：
  llm-memory group add-path my-project           # 添加当前目录
  llm-memory group add-path my-project /path/to  # 添加指定路径`,
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

		// 确定要添加的路径
		var pathToAdd string
		if len(args) > 1 {
			pathToAdd = args[1]
		} else {
			// 使用当前目录
			pwd, err := os.Getwd()
			if err != nil {
				fmt.Printf("无法获取当前目录: %v\n", err)
				return
			}
			pathToAdd = pwd
		}

		// 添加路径到组
		if err := boot.GroupService.AddPath(boot.Context(), group.ID, pathToAdd); err != nil {
			fmt.Printf("添加路径失败: %v\n", err)
			return
		}

		fmt.Printf("✨ 已将路径 '%s' 添加到组 '%s'！\n", pathToAdd, groupName)
	},
}

func init() {
	groupCmd.AddCommand(groupAddPathCmd)
}
