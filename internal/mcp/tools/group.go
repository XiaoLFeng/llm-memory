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
func RegisterGroupTools(server *mcp.Server, bs *startup.Bootstrap) {
	// group_list - 列出所有组
	mcp.AddTool(server, &mcp.Tool{
		Name:        "group_list",
		Description: `列出所有组及关联路径。`,
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
				result += fmt.Sprintf("    %s\n", p.GetPath())
			}
		}
		return NewTextResult(result), nil, nil
	})

	// group_create - 创建组
	mcp.AddTool(server, &mcp.Tool{
		Name:        "group_create",
		Description: `创建组。必填: name。可选: description。`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, input GroupCreateInput) (*mcp.CallToolResult, any, error) {
		group, err := bs.GroupService.CreateGroup(ctx, input.Name, input.Description)
		if err != nil {
			return NewErrorResult(err.Error()), nil, nil
		}
		return NewTextResult(fmt.Sprintf("组创建成功! ID: %d, 名称: %s", group.ID, group.Name)), nil, nil
	})

	// group_add_path - 添加路径到组
	mcp.AddTool(server, &mcp.Tool{
		Name:        "group_add_path",
		Description: `添加路径到组，path留空则用当前目录。`,
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
		Name:        "group_remove_path",
		Description: `从组移除指定路径。`,
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
		Name:        "group_delete",
		Description: `删除组，只解除关联不删数据。`,
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
		Name:        "group_current",
		Description: `获取当前目录的作用域信息(personal/group)。`,
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
		result += fmt.Sprintf("当前路径: %s\n", pwd)

		if scope.IncludePersonal {
			result += "Personal: 已启用\n"
		} else {
			result += "Personal: 未启用\n"
		}

		if scope.GroupID != types.GlobalGroupID {
			result += fmt.Sprintf("Group: %s (ID: %d)\n", scope.GroupName, scope.GroupID)
		} else {
			result += "Group: 不属于任何组\n"
		}

		if scope.IncludeGlobal {
			result += "Global: 已启用\n"
		} else {
			result += "Global: 未启用\n"
		}

		return NewTextResult(result), nil, nil
	})
}
