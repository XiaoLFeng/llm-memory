package tools

import (
	"context"
	"fmt"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/XiaoLFeng/llm-memory/pkg/types"
	"github.com/XiaoLFeng/llm-memory/startup"
)

// GroupListInput group_list 工具输入
type GroupListInput struct{}

// GroupCreateInput group_create 工具输入
type GroupCreateInput struct {
	Name        string `json:"name" jsonschema:"组名称，用于标识组"`
	Description string `json:"description,omitempty" jsonschema:"组的描述信息"`
}

// GroupAddPathInput group_add_path 工具输入
type GroupAddPathInput struct {
	GroupName string `json:"group_name" jsonschema:"要添加路径的组名称"`
	Path      string `json:"path,omitempty" jsonschema:"要添加的路径，留空则添加当前工作目录"`
}

// GroupRemovePathInput group_remove_path 工具输入
type GroupRemovePathInput struct {
	GroupName string `json:"group_name" jsonschema:"要移除路径的组名称"`
	Path      string `json:"path" jsonschema:"要移除的路径"`
}

// GroupDeleteInput group_delete 工具输入
type GroupDeleteInput struct {
	Name string `json:"name" jsonschema:"要删除的组名称"`
}

// GroupCurrentInput group_current 工具输入
type GroupCurrentInput struct{}

// RegisterGroupTools 注册组管理工具
// 嘿嘿~ 组管理相关的 MCP 工具都在这里！👥
func RegisterGroupTools(server *mcp.Server, bs *startup.Bootstrap) {
	// group_list - 列出所有组
	mcp.AddTool(server, &mcp.Tool{
		Name: "group_list",
		Description: `列出所有已创建的组及其包含的路径。

使用场景：
- 查看当前有哪些组
- 了解各组包含的路径
- 获取组名称用于其他操作

返回信息：组ID、名称、描述、路径列表`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, input GroupListInput) (*mcp.CallToolResult, any, error) {
		groups, err := bs.GroupService.ListGroups(ctx)
		if err != nil {
			return NewErrorResult(err.Error()), nil, nil
		}
		if len(groups) == 0 {
			return NewTextResult("暂无组"), nil, nil
		}
		result := "组列表:\n"
		for _, g := range groups {
			result += fmt.Sprintf("- [%d] %s", g.ID, g.Name)
			if g.Description != "" {
				result += fmt.Sprintf(" (%s)", g.Description)
			}
			result += fmt.Sprintf(" - %d 个路径\n", len(g.Paths))
			for _, p := range g.Paths {
				result += fmt.Sprintf("    📂 %s\n", p.Path)
			}
		}
		return NewTextResult(result), nil, nil
	})

	// group_create - 创建组
	mcp.AddTool(server, &mcp.Tool{
		Name: "group_create",
		Description: `创建一个新的组，用于管理多个路径的共享数据。

使用场景：
- 用户想要在多个目录之间共享记忆/待办/计划
- 项目有多个子目录需要统一管理
- 团队协作时需要共享信息

示例：
- 创建 "my-project" 组来管理前后端两个目录`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, input GroupCreateInput) (*mcp.CallToolResult, any, error) {
		group, err := bs.GroupService.CreateGroup(ctx, input.Name, input.Description)
		if err != nil {
			return NewErrorResult(err.Error()), nil, nil
		}
		return NewTextResult(fmt.Sprintf("组创建成功! ID: %d, 名称: %s", group.ID, group.Name)), nil, nil
	})

	// group_add_path - 添加路径到组
	mcp.AddTool(server, &mcp.Tool{
		Name: "group_add_path",
		Description: `将路径添加到指定组中。

使用场景：
- 将当前目录添加到组
- 将指定路径添加到组

注意：如果不指定路径，则添加当前工作目录`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, input GroupAddPathInput) (*mcp.CallToolResult, any, error) {
		// 获取组
		group, err := bs.GroupService.GetGroupByName(ctx, input.GroupName)
		if err != nil {
			return NewErrorResult(fmt.Sprintf("找不到组 '%s': %v", input.GroupName, err)), nil, nil
		}

		// 确定要添加的路径
		pathToAdd := input.Path
		if pathToAdd == "" {
			pwd, err := os.Getwd()
			if err != nil {
				return NewErrorResult(fmt.Sprintf("无法获取当前目录: %v", err)), nil, nil
			}
			pathToAdd = pwd
		}

		// 添加路径
		if err := bs.GroupService.AddPath(ctx, group.ID, pathToAdd); err != nil {
			return NewErrorResult(err.Error()), nil, nil
		}

		return NewTextResult(fmt.Sprintf("已将路径 '%s' 添加到组 '%s'", pathToAdd, input.GroupName)), nil, nil
	})

	// group_remove_path - 从组中移除路径
	mcp.AddTool(server, &mcp.Tool{
		Name: "group_remove_path",
		Description: `从指定组中移除路径。

使用场景：
- 某个目录不再需要与组共享数据
- 整理组的路径列表`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, input GroupRemovePathInput) (*mcp.CallToolResult, any, error) {
		// 获取组
		group, err := bs.GroupService.GetGroupByName(ctx, input.GroupName)
		if err != nil {
			return NewErrorResult(fmt.Sprintf("找不到组 '%s': %v", input.GroupName, err)), nil, nil
		}

		// 移除路径
		if err := bs.GroupService.RemovePath(ctx, group.ID, input.Path); err != nil {
			return NewErrorResult(err.Error()), nil, nil
		}

		return NewTextResult(fmt.Sprintf("已从组 '%s' 中移除路径 '%s'", input.GroupName, input.Path)), nil, nil
	})

	// group_delete - 删除组
	mcp.AddTool(server, &mcp.Tool{
		Name: "group_delete",
		Description: `删除指定的组。

注意：删除组不会删除组内的记忆、待办和计划数据，只是解除路径关联。`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, input GroupDeleteInput) (*mcp.CallToolResult, any, error) {
		// 获取组
		group, err := bs.GroupService.GetGroupByName(ctx, input.Name)
		if err != nil {
			return NewErrorResult(fmt.Sprintf("找不到组 '%s': %v", input.Name, err)), nil, nil
		}

		// 删除组
		if err := bs.GroupService.DeleteGroup(ctx, group.ID); err != nil {
			return NewErrorResult(err.Error()), nil, nil
		}

		return NewTextResult(fmt.Sprintf("组 '%s' 已删除", input.Name)), nil, nil
	})

	// group_current - 获取当前作用域
	mcp.AddTool(server, &mcp.Tool{
		Name: "group_current",
		Description: `获取当前工作目录的作用域信息。

返回信息：
- 当前路径 (Personal)
- 所属组 (Group)，如果有的话
- 全局 (Global) 状态

这可以帮助了解当前目录属于哪个组，以及会看到哪些数据。`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, input GroupCurrentInput) (*mcp.CallToolResult, any, error) {
		// 获取当前目录
		pwd, err := os.Getwd()
		if err != nil {
			return NewErrorResult(fmt.Sprintf("无法获取当前目录: %v", err)), nil, nil
		}

		// 获取当前作用域
		scope := bs.CurrentScope
		if scope == nil {
			scope = types.NewGlobalOnlyScope()
		}

		result := "当前作用域信息:\n"
		result += fmt.Sprintf("📍 当前路径: %s\n", pwd)

		if scope.IncludePersonal {
			result += "👤 Personal: ✅ 启用\n"
		} else {
			result += "👤 Personal: ❌ 未启用\n"
		}

		if scope.GroupID != types.GlobalGroupID {
			result += fmt.Sprintf("👥 Group: ✅ %s (ID: %d)\n", scope.GroupName, scope.GroupID)
		} else {
			result += "👥 Group: ❌ 不属于任何组\n"
		}

		if scope.IncludeGlobal {
			result += "🌐 Global: ✅ 启用\n"
		} else {
			result += "🌐 Global: ❌ 未启用\n"
		}

		return NewTextResult(result), nil, nil
	})
}
