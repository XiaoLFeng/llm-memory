package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/XiaoLFeng/llm-memory/pkg/types"
	"github.com/XiaoLFeng/llm-memory/startup"
	"github.com/spf13/cobra"
)

// groupCurrentCmd 显示当前作用域的命令
// 嘿嘿~ 看看当前目录的作用域是什么！🔍
var groupCurrentCmd = &cobra.Command{
	Use:   "current",
	Short: "显示当前作用域",
	Long: `显示当前工作目录所属的作用域信息~ ✨

包括：当前路径 (Personal)、所属组 (Group)、全局 (Global)

示例：
  llm-memory group current`,
	Run: func(cmd *cobra.Command, args []string) {
		// 初始化 Bootstrap
		boot := startup.New()
		if err := boot.Initialize(context.Background()); err != nil {
			fmt.Printf("初始化失败: %v\n", err)
			return
		}
		defer func() { _ = boot.Shutdown() }()

		// 获取当前目录
		pwd, err := os.Getwd()
		if err != nil {
			fmt.Printf("无法获取当前目录: %v\n", err)
			return
		}

		// 获取当前作用域
		scope := boot.CurrentScope
		if scope == nil {
			scope = types.NewGlobalOnlyScope()
		}

		fmt.Println("🔍 当前作用域信息:")
		fmt.Println("─────────────────────────────────────")
		fmt.Printf("📍 当前路径: %s\n", pwd)

		// Personal 作用域
		if scope.IncludePersonal {
			fmt.Printf("👤 Personal: ✅ 启用 (精确匹配当前路径)\n")
		} else {
			fmt.Println("👤 Personal: ❌ 未启用")
		}

		// Group 作用域
		if scope.GroupID != types.GlobalGroupID {
			fmt.Printf("👥 Group: ✅ 启用 (组: %s, ID: %d)\n", scope.GroupName, scope.GroupID)
		} else {
			fmt.Println("👥 Group: ❌ 当前路径不属于任何组")
		}

		// Global 作用域
		if scope.IncludeGlobal {
			fmt.Println("🌐 Global: ✅ 启用")
		} else {
			fmt.Println("🌐 Global: ❌ 未启用")
		}

		fmt.Println("─────────────────────────────────────")

		// 提示信息
		if scope.GroupID == types.GlobalGroupID {
			fmt.Println("\n💡 提示: 使用 'llm-memory group add-path <group-name>' 将当前目录添加到组")
		}
	},
}

func init() {
	groupCmd.AddCommand(groupCurrentCmd)
}
