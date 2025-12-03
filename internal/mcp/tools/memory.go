package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/XiaoLFeng/llm-memory/internal/models/dto"
	"github.com/XiaoLFeng/llm-memory/pkg/types"
	"github.com/XiaoLFeng/llm-memory/startup"
)

// MemoryListInput memory_list 工具输入
type MemoryListInput struct {
	Scope string `json:"scope,omitempty" jsonschema:"作用域过滤(personal/group/global/all)，默认all显示全部"`
}

// MemoryCreateInput memory_create 工具输入
type MemoryCreateInput struct {
	Title    string   `json:"title" jsonschema:"记忆标题，简洁概括内容"`
	Content  string   `json:"content" jsonschema:"记忆的详细内容，支持多行文本"`
	Category string   `json:"category,omitempty" jsonschema:"记忆分类，如：用户偏好、技术文档。默认为'默认'"`
	Tags     []string `json:"tags,omitempty" jsonschema:"标签列表，用于细粒度分类和搜索"`
	Scope    string   `json:"scope,omitempty" jsonschema:"保存到哪个作用域(personal/group/global)，默认global"`
}

// MemoryDeleteInput memory_delete 工具输入
type MemoryDeleteInput struct {
	ID int64 `json:"id" jsonschema:"要删除的记忆ID"`
}

// MemorySearchInput memory_search 工具输入
type MemorySearchInput struct {
	Keyword string `json:"keyword" jsonschema:"搜索关键词，在标题和内容中模糊匹配"`
	Scope   string `json:"scope,omitempty" jsonschema:"作用域过滤(personal/group/global/all)，默认all显示全部"`
}

// MemoryGetInput memory_get 工具输入
type MemoryGetInput struct {
	ID int64 `json:"id" jsonschema:"要获取的记忆ID"`
}

// MemoryUpdateInput memory_update 工具输入
type MemoryUpdateInput struct {
	ID       int64    `json:"id" jsonschema:"要更新的记忆ID"`
	Title    string   `json:"title,omitempty" jsonschema:"新标题（可选）"`
	Content  string   `json:"content,omitempty" jsonschema:"新内容（可选）"`
	Category string   `json:"category,omitempty" jsonschema:"新分类（可选）"`
	Tags     []string `json:"tags,omitempty" jsonschema:"新标签列表（可选）"`
	Priority int      `json:"priority,omitempty" jsonschema:"新优先级 1-4（可选）"`
}

// RegisterMemoryTools 注册记忆管理工具
func RegisterMemoryTools(server *mcp.Server, bs *startup.Bootstrap) {
	// memory_list - 列出所有记忆
	mcp.AddTool(server, &mcp.Tool{
		Name:        "memory_list",
		Description: `列出可见记忆。scope参数: personal/group/global/all(默认)。`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, input MemoryListInput) (*mcp.CallToolResult, any, error) {
		// 构建作用域上下文
		scopeCtx := buildScopeContext(input.Scope, bs)

		memories, err := bs.MemoryService.ListMemoriesByScope(ctx, input.Scope, scopeCtx)
		if err != nil {
			return NewErrorResult(err.Error()), nil, nil
		}
		if len(memories) == 0 {
			return NewTextResult("暂无记忆"), nil, nil
		}
		result := "记忆列表:\n"
		for _, m := range memories {
			scopeTag := getScopeTag(m.GroupID, m.Path)
			result += fmt.Sprintf("- [%d] %s (分类: %s) %s\n", m.ID, m.Title, m.Category, scopeTag)
		}
		return NewTextResult(result), nil, nil
	})

	// memory_create - 创建新记忆
	mcp.AddTool(server, &mcp.Tool{
		Name:        "memory_create",
		Description: `创建记忆条目。必填: title(标题)、content(内容)。可选: category(分类)、tags(标签列表)、scope(作用域personal/group/global，默认global)。`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, input MemoryCreateInput) (*mcp.CallToolResult, any, error) {
		// 构建创建 DTO
		createDTO := &dto.MemoryCreateDTO{
			Title:    input.Title,
			Content:  input.Content,
			Category: input.Category,
			Tags:     input.Tags,
			Priority: 1, // 默认优先级
			Scope:    input.Scope,
		}

		// 构建作用域上下文
		scopeCtx := buildScopeContext(input.Scope, bs)

		memory, err := bs.MemoryService.CreateMemory(ctx, createDTO, scopeCtx)
		if err != nil {
			return NewErrorResult(err.Error()), nil, nil
		}
		scopeTag := getScopeTag(memory.GroupID, memory.Path)
		return NewTextResult(fmt.Sprintf("记忆创建成功! ID: %d, 标题: %s %s", memory.ID, memory.Title, scopeTag)), nil, nil
	})

	// memory_delete - 删除记忆
	mcp.AddTool(server, &mcp.Tool{
		Name:        "memory_delete",
		Description: `删除指定ID的记忆，不可恢复。`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, input MemoryDeleteInput) (*mcp.CallToolResult, any, error) {
		if err := bs.MemoryService.DeleteMemory(ctx, input.ID); err != nil {
			return NewErrorResult(err.Error()), nil, nil
		}
		return NewTextResult(fmt.Sprintf("记忆 %d 已删除", input.ID)), nil, nil
	})

	// memory_search - 搜索记忆
	mcp.AddTool(server, &mcp.Tool{
		Name:        "memory_search",
		Description: `搜索记忆，在标题和内容中模糊匹配keyword。scope参数: personal/group/global/all(默认)。`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, input MemorySearchInput) (*mcp.CallToolResult, any, error) {
		// 构建作用域上下文
		scopeCtx := buildScopeContext(input.Scope, bs)

		memories, err := bs.MemoryService.SearchMemoriesByScope(ctx, input.Keyword, input.Scope, scopeCtx)
		if err != nil {
			return NewErrorResult(err.Error()), nil, nil
		}
		if len(memories) == 0 {
			return NewTextResult("未找到匹配的记忆"), nil, nil
		}
		result := fmt.Sprintf("搜索结果 (%d 条):\n", len(memories))
		for _, m := range memories {
			scopeTag := getScopeTag(m.GroupID, m.Path)
			result += fmt.Sprintf("- [%d] %s %s\n", m.ID, m.Title, scopeTag)
		}
		return NewTextResult(result), nil, nil
	})

	// memory_get - 获取记忆详情
	mcp.AddTool(server, &mcp.Tool{
		Name:        "memory_get",
		Description: `获取指定ID记忆的完整详情，包括内容、分类、标签等。`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, input MemoryGetInput) (*mcp.CallToolResult, any, error) {
		memory, err := bs.MemoryService.GetMemory(ctx, input.ID)
		if err != nil {
			return NewErrorResult(err.Error()), nil, nil
		}

		// 转换标签
		tags := make([]string, 0, len(memory.Tags))
		for _, t := range memory.Tags {
			tags = append(tags, t.Tag)
		}

		scopeTag := getScopeTag(memory.GroupID, memory.Path)
		result := fmt.Sprintf(`记忆详情:
ID: %d
标题: %s
分类: %s
优先级: %d
标签: %v
作用域: %s
创建时间: %s
更新时间: %s

内容:
%s`,
			memory.ID,
			memory.Title,
			memory.Category,
			memory.Priority,
			tags,
			scopeTag,
			memory.CreatedAt.Format("2006-01-02 15:04:05"),
			memory.UpdatedAt.Format("2006-01-02 15:04:05"),
			memory.Content,
		)
		return NewTextResult(result), nil, nil
	})

	// memory_update - 更新记忆
	mcp.AddTool(server, &mcp.Tool{
		Name:        "memory_update",
		Description: `更新记忆，只更新提供的字段。可选: title、content、category、tags、priority(1-4)。`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, input MemoryUpdateInput) (*mcp.CallToolResult, any, error) {
		// 构建更新 DTO
		updateDTO := &dto.MemoryUpdateDTO{
			ID: input.ID,
		}

		// 只设置提供了的字段
		if input.Title != "" {
			updateDTO.Title = &input.Title
		}
		if input.Content != "" {
			updateDTO.Content = &input.Content
		}
		if input.Category != "" {
			updateDTO.Category = &input.Category
		}
		if len(input.Tags) > 0 {
			updateDTO.Tags = &input.Tags
		}
		if input.Priority > 0 && input.Priority <= 4 {
			updateDTO.Priority = &input.Priority
		}

		// 检查是否有更新
		if updateDTO.Title == nil && updateDTO.Content == nil && updateDTO.Category == nil && updateDTO.Tags == nil && updateDTO.Priority == nil {
			return NewErrorResult("没有提供要更新的字段"), nil, nil
		}

		// 执行更新
		if err := bs.MemoryService.UpdateMemory(ctx, updateDTO); err != nil {
			return NewErrorResult(err.Error()), nil, nil
		}

		return NewTextResult(fmt.Sprintf("记忆 %d 更新成功", input.ID)), nil, nil
	})
}

// buildScopeContext 根据 scope 字符串构建 ScopeContext
func buildScopeContext(scope string, bs *startup.Bootstrap) *types.ScopeContext {
	// 获取当前工作目录和作用域上下文
	currentScope := bs.CurrentScope
	if currentScope == nil {
		currentScope = types.NewGlobalOnlyScope()
	}

	switch strings.ToLower(scope) {
	case "personal":
		return &types.ScopeContext{
			CurrentPath:     currentScope.CurrentPath,
			GroupID:         types.GlobalGroupID,
			IncludePersonal: true,
			IncludeGroup:    false,
			IncludeGlobal:   false,
		}
	case "group":
		return &types.ScopeContext{
			CurrentPath:     currentScope.CurrentPath,
			GroupID:         currentScope.GroupID,
			GroupName:       currentScope.GroupName,
			IncludePersonal: false,
			IncludeGroup:    true,
			IncludeGlobal:   false,
		}
	case "global":
		return &types.ScopeContext{
			CurrentPath:     currentScope.CurrentPath,
			GroupID:         types.GlobalGroupID,
			IncludePersonal: false,
			IncludeGroup:    false,
			IncludeGlobal:   true,
		}
	default: // "all" 或空字符串
		return currentScope
	}
}

// getScopeTag 获取作用域标签
func getScopeTag(groupID int64, path string) string {
	if path != "" {
		return "[Personal]"
	}
	if groupID > 0 {
		return "[Group]"
	}
	return "[Global]"
}

// tagsToStringSlice 将 MemoryTag 切片转换为字符串切片
func tagsToStringSlice(tags interface{}) []string {
	result := []string{}
	// 处理不同类型的 tags
	switch t := tags.(type) {
	case []string:
		return t
	default:
		return result
	}
}
