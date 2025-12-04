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
	Code     string   `json:"code" jsonschema:"记忆唯一标识码"`
	Title    string   `json:"title" jsonschema:"记忆标题，简洁概括内容"`
	Content  string   `json:"content" jsonschema:"记忆的详细内容，支持多行文本"`
	Category string   `json:"category,omitempty" jsonschema:"记忆分类，如：用户偏好、技术文档。默认为'默认'"`
	Tags     []string `json:"tags,omitempty" jsonschema:"标签列表，用于细粒度分类和搜索"`
	Global   bool     `json:"global,omitempty" jsonschema:"是否写入全局（true 全局；false/省略 当前路径/组内）"`
	Scope    string   `json:"scope,omitempty" jsonschema:"查询筛选仍可用的作用域 personal/group/global/all"`
}

// MemoryDeleteInput memory_delete 工具输入
type MemoryDeleteInput struct {
	Code string `json:"code" jsonschema:"要删除的记忆code"`
}

// MemorySearchInput memory_search 工具输入
type MemorySearchInput struct {
	Keyword string `json:"keyword" jsonschema:"搜索关键词，在标题和内容中模糊匹配"`
	Scope   string `json:"scope,omitempty" jsonschema:"作用域过滤(personal/group/global/all)，默认all显示全部"`
}

// MemoryGetInput memory_get 工具输入
type MemoryGetInput struct {
	Code string `json:"code" jsonschema:"要获取的记忆code"`
}

// MemoryUpdateInput memory_update 工具输入
type MemoryUpdateInput struct {
	Code     string   `json:"code" jsonschema:"要更新的记忆code"`
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
		Name: "memory_list",
		Description: `列出可见记忆。scope参数说明（安全隔离）：
  - personal: 仅当前路径的私有数据
  - group: 仅当前小组的数据（需已加入小组）
  - global: 仅全局可见数据
  - all/省略: 全局 + 当前路径相关（默认，权限隔离）`,
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
			scopeTag := getScopeTagWithContext(m.Global, m.PathID, bs.CurrentScope)
			result += fmt.Sprintf("- [%s] %s (分类: %s) %s\n", m.Code, m.Title, m.Category, scopeTag)
		}
		return NewTextResult(result), nil, nil
	})

	// memory_create - 创建新记忆
	mcp.AddTool(server, &mcp.Tool{
		Name:        "memory_create",
		Description: `创建记忆条目，适合长期事实、偏好、上下文片段。必填: title、content。可选: category、tags、global。global=true 存入全局；省略/false 存当前路径(私有，若在组内则组可见)。短任务请用 todo_create，需要进度跟踪的多步骤目标请用 plan_create。scope 参数仅用于列表筛选。`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, input MemoryCreateInput) (*mcp.CallToolResult, any, error) {
		// 构建创建 DTO
		createDTO := &dto.MemoryCreateDTO{
			Code:     input.Code,
			Title:    input.Title,
			Content:  input.Content,
			Category: input.Category,
			Tags:     input.Tags,
			Priority: 1, // 默认优先级
			Global:   input.Global,
		}

		// 构建作用域上下文
		scopeCtx := buildScopeContext(input.Scope, bs)

		memory, err := bs.MemoryService.CreateMemory(ctx, createDTO, scopeCtx)
		if err != nil {
			return NewErrorResult(err.Error()), nil, nil
		}
		scopeTag := getScopeTagWithContext(memory.Global, memory.PathID, bs.CurrentScope)
		return NewTextResult(fmt.Sprintf("记忆创建成功! Code: %s, 标题: %s %s", memory.Code, memory.Title, scopeTag)), nil, nil
	})

	// memory_delete - 删除记忆
	mcp.AddTool(server, &mcp.Tool{
		Name:        "memory_delete",
		Description: `删除指定code的记忆，不可恢复。`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, input MemoryDeleteInput) (*mcp.CallToolResult, any, error) {
		if err := bs.MemoryService.DeleteMemory(ctx, input.Code); err != nil {
			return NewErrorResult(err.Error()), nil, nil
		}
		return NewTextResult(fmt.Sprintf("记忆 %s 已删除", input.Code)), nil, nil
	})

	// memory_search - 搜索记忆
	mcp.AddTool(server, &mcp.Tool{
		Name:        "memory_search",
		Description: `搜索记忆（标题与内容模糊匹配 keyword）。scope: personal/group/global/all；默认不填=全部（全局+私有+小组）。`,
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
			scopeTag := getScopeTagWithContext(m.Global, m.PathID, bs.CurrentScope)
			result += fmt.Sprintf("- [%s] %s %s\n", m.Code, m.Title, scopeTag)
		}
		return NewTextResult(result), nil, nil
	})

	// memory_get - 获取记忆详情
	mcp.AddTool(server, &mcp.Tool{
		Name:        "memory_get",
		Description: `获取指定code记忆的完整详情，包括内容、分类、标签等。`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, input MemoryGetInput) (*mcp.CallToolResult, any, error) {
		memory, err := bs.MemoryService.GetMemory(ctx, input.Code)
		if err != nil {
			return NewErrorResult(err.Error()), nil, nil
		}

		// 转换标签
		tags := make([]string, 0, len(memory.Tags))
		for _, t := range memory.Tags {
			tags = append(tags, t.Tag)
		}

		scopeTag := getScopeTagWithContext(memory.Global, memory.PathID, bs.CurrentScope)
		var sb strings.Builder
		_, _ = fmt.Fprintf(&sb, "记忆详情:\n")
		_, _ = fmt.Fprintf(&sb, "Code: %s\n", memory.Code)
		_, _ = fmt.Fprintf(&sb, "标题: %s\n", memory.Title)
		_, _ = fmt.Fprintf(&sb, "分类: %s\n", memory.Category)
		_, _ = fmt.Fprintf(&sb, "优先级: %d\n", memory.Priority)
		_, _ = fmt.Fprintf(&sb, "标签: %v\n", tags)
		_, _ = fmt.Fprintf(&sb, "作用域: %s\n", scopeTag)
		_, _ = fmt.Fprintf(&sb, "创建时间: %s\n", memory.CreatedAt.Format("2006-01-02 15:04:05"))
		_, _ = fmt.Fprintf(&sb, "更新时间: %s\n", memory.UpdatedAt.Format("2006-01-02 15:04:05"))
		_, _ = fmt.Fprintf(&sb, "\n内容:\n%s", memory.Content)
		result := sb.String()
		return NewTextResult(result), nil, nil
	})

	// memory_update - 更新记忆
	mcp.AddTool(server, &mcp.Tool{
		Name:        "memory_update",
		Description: `更新记忆，只更新提供的字段（title/content/category/tags/priority1-4）；至少提供一个字段，否则返回错误。`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, input MemoryUpdateInput) (*mcp.CallToolResult, any, error) {
		// 构建更新 DTO
		updateDTO := &dto.MemoryUpdateDTO{
			Code: input.Code,
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

		return NewTextResult(fmt.Sprintf("记忆 %s 更新成功", input.Code)), nil, nil
	})
}

// buildScopeContext 根据 scope 构建 ScopeContext（保持简单：默认返回当前上下文）
func buildScopeContext(_ string, bs *startup.Bootstrap) *types.ScopeContext {
	if bs.CurrentScope == nil {
		return types.NewScopeContext("")
	}
	return bs.CurrentScope
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
