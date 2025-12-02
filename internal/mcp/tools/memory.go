package tools

import (
	"context"
	"fmt"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"

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
	ID int `json:"id" jsonschema:"要删除的记忆ID"`
}

// MemorySearchInput memory_search 工具输入
type MemorySearchInput struct {
	Keyword string `json:"keyword" jsonschema:"搜索关键词，在标题和内容中模糊匹配"`
	Scope   string `json:"scope,omitempty" jsonschema:"作用域过滤(personal/group/global/all)，默认all显示全部"`
}

// RegisterMemoryTools 注册记忆管理工具
// 嘿嘿~ 记忆相关的 MCP 工具都在这里！(´∀｀)💖
func RegisterMemoryTools(server *mcp.Server, bs *startup.Bootstrap) {
	// memory_list - 列出所有记忆
	mcp.AddTool(server, &mcp.Tool{
		Name: "memory_list",
		Description: `列出用户存储的所有记忆条目。

使用场景：
- 查看当前已保存的所有记忆
- 在创建新记忆前检查是否已存在类似内容
- 获取记忆ID用于后续的删除或更新操作

返回信息：记忆ID、标题、分类

注意：如果记忆数量较多，建议使用 memory_search 进行精确查找

作用域说明：
- personal: 只显示当前目录的记忆
- group: 只显示当前组的记忆
- global: 只显示全局记忆
- all: 显示所有可见记忆（默认）`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, input MemoryListInput) (*mcp.CallToolResult, any, error) {
		// 构建作用域上下文
		scope := buildScopeContext(input.Scope, bs)

		memories, err := bs.MemoryService.ListMemoriesByScope(ctx, scope)
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
		Name: "memory_create",
		Description: `创建一条新的记忆条目，用于持久化存储重要信息。

使用场景：
- 保存用户的偏好设置（如编程语言偏好、代码风格等）
- 记录项目相关的重要信息（架构决策、技术选型等）
- 存储需要跨会话记住的任何信息

最佳实践：
- 标题应简洁明了，便于后续搜索
- 内容应详细完整，包含所有相关信息
- 合理使用分类和标签，便于组织管理

示例：
- 标题："用户编程偏好"，分类："用户偏好"，标签：["编程", "偏好"]
- 标题："项目数据库设计"，分类："技术文档"，标签：["数据库", "MySQL"]

作用域说明：
- personal: 保存到当前目录（只在此目录可见）
- group: 保存到当前组（组内所有路径可见）
- global: 保存为全局（任何地方可见，默认）`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, input MemoryCreateInput) (*mcp.CallToolResult, any, error) {
		category := input.Category
		if category == "" {
			category = "默认"
		}

		// 根据 scope 确定 groupID 和 path
		groupID, path := resolveScopeForCreate(input.Scope, bs)

		memory, err := bs.MemoryService.CreateMemory(ctx, input.Title, input.Content, category, input.Tags, 1, groupID, path)
		if err != nil {
			return NewErrorResult(err.Error()), nil, nil
		}
		scopeTag := getScopeTag(groupID, path)
		return NewTextResult(fmt.Sprintf("记忆创建成功! ID: %d, 标题: %s %s", memory.ID, memory.Title, scopeTag)), nil, nil
	})

	// memory_delete - 删除记忆
	mcp.AddTool(server, &mcp.Tool{
		Name: "memory_delete",
		Description: `删除指定ID的记忆条目。

使用场景：
- 删除过时或不再需要的记忆
- 清理错误创建的记忆条目
- 用户明确要求删除某条记忆

注意事项：
- 删除操作不可恢复，请确认后再执行
- 需要先通过 memory_list 或 memory_search 获取正确的记忆ID
- 如果不确定要删除哪条记忆，建议先查看记忆列表`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, input MemoryDeleteInput) (*mcp.CallToolResult, any, error) {
		if err := bs.MemoryService.DeleteMemory(ctx, input.ID); err != nil {
			return NewErrorResult(err.Error()), nil, nil
		}
		return NewTextResult(fmt.Sprintf("记忆 %d 已删除", input.ID)), nil, nil
	})

	// memory_search - 搜索记忆
	mcp.AddTool(server, &mcp.Tool{
		Name: "memory_search",
		Description: `根据关键词搜索记忆，在标题和内容中进行模糊匹配。

使用场景：
- 快速查找特定主题的记忆
- 在回答用户问题前检索相关背景信息
- 查找与当前任务相关的历史记录

搜索技巧：
- 使用具体的关键词获得更精确的结果
- 可以搜索标题或内容中的任意文本
- 支持中英文关键词

建议：在执行任务前，先搜索是否有相关的记忆可以参考，这样可以提供更个性化的服务

作用域说明：
- personal: 只搜索当前目录的记忆
- group: 只搜索当前组的记忆
- global: 只搜索全局记忆
- all: 搜索所有可见记忆（默认）`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, input MemorySearchInput) (*mcp.CallToolResult, any, error) {
		// 构建作用域上下文
		scope := buildScopeContext(input.Scope, bs)

		memories, err := bs.MemoryService.SearchMemoriesByScope(ctx, scope, input.Keyword)
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
}

// buildScopeContext 根据 scope 字符串构建 ScopeContext
// 嘿嘿~ 这是通用的作用域构建辅助函数！✨
func buildScopeContext(scope string, bs *startup.Bootstrap) *types.ScopeContext {
	// 获取当前工作目录和作用域上下文
	currentScope := bs.CurrentScope
	if currentScope == nil {
		currentScope = types.NewGlobalOnlyScope()
	}

	switch scope {
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

// resolveScopeForCreate 解析创建时的作用域
// 返回 groupID 和 path
func resolveScopeForCreate(scope string, bs *startup.Bootstrap) (int, string) {
	currentScope := bs.CurrentScope
	if currentScope == nil {
		return types.GlobalGroupID, ""
	}

	switch scope {
	case "personal":
		pwd, _ := os.Getwd()
		return types.GlobalGroupID, pwd
	case "group":
		if currentScope.GroupID != types.GlobalGroupID {
			return currentScope.GroupID, ""
		}
		// 如果不属于任何组，回退到 global
		return types.GlobalGroupID, ""
	default: // "global" 或空字符串
		return types.GlobalGroupID, ""
	}
}

// getScopeTag 获取作用域标签
func getScopeTag(groupID int, path string) string {
	if path != "" {
		return "[Personal]"
	}
	if groupID != types.GlobalGroupID {
		return "[Group]"
	}
	return "[Global]"
}
