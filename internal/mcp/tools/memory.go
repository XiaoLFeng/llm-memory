package tools

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/XiaoLFeng/llm-memory/startup"
)

// MemoryListInput memory_list 工具输入
type MemoryListInput struct{}

// MemoryCreateInput memory_create 工具输入
type MemoryCreateInput struct {
	Title    string   `json:"title" jsonschema:"description=记忆的标题，简洁明了地概括记忆内容，例如：'用户偏好设置'、'项目架构说明'"`
	Content  string   `json:"content" jsonschema:"description=记忆的详细内容，可以是任意文本信息，支持多行文本，例如：用户的具体偏好、技术方案细节等"`
	Category string   `json:"category,omitempty" jsonschema:"description=记忆的分类，用于组织和筛选记忆，例如：'用户偏好'、'技术文档'、'会议记录'。如不指定则默认为'默认'"`
	Tags     []string `json:"tags,omitempty" jsonschema:"description=记忆的标签列表，用于更细粒度的分类和搜索，例如：['重要', 'Go语言', '架构设计']"`
}

// MemoryDeleteInput memory_delete 工具输入
type MemoryDeleteInput struct {
	ID int `json:"id" jsonschema:"description=要删除的记忆ID，可通过 memory_list 或 memory_search 获取"`
}

// MemorySearchInput memory_search 工具输入
type MemorySearchInput struct {
	Keyword string `json:"keyword" jsonschema:"description=搜索关键词，将在记忆的标题和内容中进行模糊匹配，支持中英文"`
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

注意：如果记忆数量较多，建议使用 memory_search 进行精确查找`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, input MemoryListInput) (*mcp.CallToolResult, any, error) {
		memories, err := bs.MemoryService.ListMemories(ctx)
		if err != nil {
			return NewErrorResult(err.Error()), nil, nil
		}
		if len(memories) == 0 {
			return NewTextResult("暂无记忆"), nil, nil
		}
		result := "记忆列表:\n"
		for _, m := range memories {
			result += fmt.Sprintf("- [%d] %s (分类: %s)\n", m.ID, m.Title, m.Category)
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
- 标题："项目数据库设计"，分类："技术文档"，标签：["数据库", "MySQL"]`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, input MemoryCreateInput) (*mcp.CallToolResult, any, error) {
		category := input.Category
		if category == "" {
			category = "默认"
		}
		memory, err := bs.MemoryService.CreateMemory(ctx, input.Title, input.Content, category, input.Tags, 1)
		if err != nil {
			return NewErrorResult(err.Error()), nil, nil
		}
		return NewTextResult(fmt.Sprintf("记忆创建成功! ID: %d, 标题: %s", memory.ID, memory.Title)), nil, nil
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

建议：在执行任务前，先搜索是否有相关的记忆可以参考，这样可以提供更个性化的服务`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, input MemorySearchInput) (*mcp.CallToolResult, any, error) {
		memories, err := bs.MemoryService.SearchMemories(ctx, input.Keyword)
		if err != nil {
			return NewErrorResult(err.Error()), nil, nil
		}
		if len(memories) == 0 {
			return NewTextResult("未找到匹配的记忆"), nil, nil
		}
		result := fmt.Sprintf("搜索结果 (%d 条):\n", len(memories))
		for _, m := range memories {
			result += fmt.Sprintf("- [%d] %s\n", m.ID, m.Title)
		}
		return NewTextResult(result), nil, nil
	})
}
