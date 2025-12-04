package tools

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/XiaoLFeng/llm-memory/startup"
)

// GroupAddPathInput group_add_path 工具输入
type GroupAddPathInput struct {
	GroupName string `json:"group_name" jsonschema:"要添加路径的组名称"`
	Path      string `json:"path,omitempty" jsonschema:"要添加的路径，留空则添加当前工作目录"`
}

// validateGroupOperationPermission 验证组操作权限
func validateGroupOperationPermission(ctx context.Context, bs *startup.Bootstrap, groupName string) error {
	// 1. 验证组是否存在
	group, err := bs.GroupService.GetGroupByName(ctx, groupName)
	if err != nil {
		return fmt.Errorf("组不存在: %s", groupName)
	}

	// 2. 获取当前路径
	currentPath := bs.CurrentScope.CurrentPath
	if currentPath == "" {
		return fmt.Errorf("当前不在任何项目路径中")
	}

	// 3. 验证当前路径是否已经在其他组中
	existingGroup, err := bs.GroupService.GetGroupByPath(ctx, currentPath)
	if err == nil && existingGroup.ID != group.ID {
		return fmt.Errorf("当前路径已在组 '%s' 中，无法重复加入", existingGroup.Name)
	}

	return nil
}

// RegisterGroupTools 注册组管理工具
func RegisterGroupTools(server *mcp.Server, bs *startup.Bootstrap) {
	// group_add_path - 添加路径到组（增加权限验证）
	mcp.AddTool(server, &mcp.Tool{
		Name:        "group_add_path",
		Description: `将当前路径添加到指定组。注意：只能操作当前路径，不能操作其他路径。如果当前路径已在其他组中，会先移除再加入新组。`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, input GroupAddPathInput) (*mcp.CallToolResult, any, error) {
		// 权限验证
		if err := validateGroupOperationPermission(ctx, bs, input.GroupName); err != nil {
			return NewErrorResult(err.Error()), nil, nil
		}

		// 获取组
		group, err := bs.GroupService.GetGroupByName(ctx, input.GroupName)
		if err != nil {
			return NewErrorResult(fmt.Sprintf("找不到组 '%s': %v", input.GroupName, err)), nil, nil
		}

		// 确定要添加的路径（只能是当前路径）
		pathToAdd := input.Path
		if pathToAdd != "" {
			// 如果用户指定了路径，验证是否与当前路径一致
			currentPath := bs.CurrentScope.CurrentPath
			if pathToAdd != currentPath {
				return NewErrorResult(fmt.Sprintf("只能操作当前路径。当前路径: %s，指定路径: %s", currentPath, pathToAdd)), nil, nil
			}
		} else {
			// 使用当前路径
			pathToAdd = bs.CurrentScope.CurrentPath
			if pathToAdd == "" {
				return NewErrorResult("当前不在任何项目路径中"), nil, nil
			}
		}

		// 检查是否已经在该组中
		for _, existingPath := range group.Paths {
			if existingPath.GetPath() == pathToAdd {
				return NewTextResult(fmt.Sprintf("当前路径已在组 '%s' 中", input.GroupName)), nil, nil
			}
		}

		// 添加路径
		if err := bs.GroupService.AddPath(ctx, group.ID, pathToAdd); err != nil {
			return NewErrorResult(err.Error()), nil, nil
		}

		return NewTextResult(fmt.Sprintf("已将当前路径 '%s' 添加到组 '%s'", pathToAdd, input.GroupName)), nil, nil
	})
}
