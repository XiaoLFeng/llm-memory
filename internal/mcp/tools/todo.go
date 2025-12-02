package tools

import (
	"context"
	"fmt"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/XiaoLFeng/llm-memory/pkg/types"
	"github.com/XiaoLFeng/llm-memory/startup"
)

// TodoListInput todo_list 工具输入
type TodoListInput struct{}

// TodoCreateInput todo_create 工具输入
type TodoCreateInput struct {
	Title       string `json:"title" jsonschema:"待办标题，简洁描述任务"`
	Description string `json:"description,omitempty" jsonschema:"待办的详细描述"`
	Priority    int    `json:"priority,omitempty" jsonschema:"优先级(1低/2中/3高/4紧急)，默认2"`
}

// TodoCompleteInput todo_complete 工具输入
type TodoCompleteInput struct {
	ID int `json:"id" jsonschema:"要完成的待办事项ID"`
}

// TodoTodayInput todo_today 工具输入
type TodoTodayInput struct{}

// RegisterTodoTools 注册 TODO 管理工具
// 嗯嗯！待办事项相关的 MCP 工具都在这里！🎮
func RegisterTodoTools(server *mcp.Server, bs *startup.Bootstrap) {
	// todo_list - 列出所有待办
	mcp.AddTool(server, &mcp.Tool{
		Name: "todo_list",
		Description: `列出用户的所有待办事项，包含状态和优先级信息。

使用场景：
- 查看所有待办事项的整体情况
- 了解各任务的状态和优先级
- 获取待办ID用于标记完成

返回信息：待办ID、标题、状态、优先级

状态说明：
- 待处理：新创建，尚未开始
- 进行中：已开始处理
- 已完成：任务完成
- 已取消：任务取消

优先级说明：
- 低：可延后处理的任务
- 中：正常优先级（默认）
- 高：需要优先处理
- 紧急：需要立即处理`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, input TodoListInput) (*mcp.CallToolResult, any, error) {
		todos, err := bs.TodoService.ListTodos(ctx)
		if err != nil {
			return NewErrorResult(err.Error()), nil, nil
		}
		if len(todos) == 0 {
			return NewTextResult("暂无待办事项"), nil, nil
		}
		result := "待办事项列表:\n"
		for _, t := range todos {
			status := getTodoStatusText(t.Status)
			priority := getPriorityText(t.Priority)
			result += fmt.Sprintf("- [%d] %s (%s, %s)\n", t.ID, t.Title, status, priority)
		}
		return NewTextResult(result), nil, nil
	})

	// todo_create - 创建待办
	mcp.AddTool(server, &mcp.Tool{
		Name: "todo_create",
		Description: `创建一个新的待办事项，用于记录需要完成的短期任务。

使用场景：
- 用户提出需要完成的具体任务
- 记录会话中发现的待处理事项
- 分解复杂任务为多个待办

待办事项 vs 计划：
- 待办事项：短期、具体、一次性完成的任务（如"修复Bug"、"回复邮件"）
- 计划：长期、复杂、需要跟踪进度的目标（如"完成项目重构"）

优先级选择指南：
- 1（低）：不紧急且不重要，可以延后
- 2（中）：正常任务，按顺序处理（默认）
- 3（高）：重要任务，需要优先安排
- 4（紧急）：紧急任务，需要立即处理

示例：
- 标题："修复用户登录失败问题"，优先级：4（紧急）
- 标题："更新项目文档"，优先级：2（中）`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, input TodoCreateInput) (*mcp.CallToolResult, any, error) {
		priority := types.Priority(input.Priority)
		if priority == 0 {
			priority = types.TodoPriorityMedium
		}
		todo, err := bs.TodoService.CreateTodo(ctx, input.Title, input.Description, priority, nil)
		if err != nil {
			return NewErrorResult(err.Error()), nil, nil
		}
		return NewTextResult(fmt.Sprintf("待办事项创建成功! ID: %d, 标题: %s", todo.ID, todo.Title)), nil, nil
	})

	// todo_complete - 完成待办
	mcp.AddTool(server, &mcp.Tool{
		Name: "todo_complete",
		Description: `将指定的待办事项标记为已完成。

使用场景：
- 用户确认任务已完成
- 任务目标已达成
- 需要关闭某个待办事项

注意事项：
- 已完成的待办无法再次标记为完成
- 已取消的待办无法标记为完成
- 完成后会记录完成时间

建议：在用户明确表示任务完成后使用此工具，可以先通过 todo_list 或 todo_today 确认待办ID`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, input TodoCompleteInput) (*mcp.CallToolResult, any, error) {
		if err := bs.TodoService.CompleteTodo(ctx, input.ID); err != nil {
			return NewErrorResult(err.Error()), nil, nil
		}
		return NewTextResult(fmt.Sprintf("待办事项 %d 已标记为完成", input.ID)), nil, nil
	})

	// todo_today - 获取今日待办
	mcp.AddTool(server, &mcp.Tool{
		Name: "todo_today",
		Description: `获取今日的待办事项列表，帮助用户聚焦当天任务。

使用场景：
- 每日工作开始时查看今天的任务
- 用户询问"今天有什么任务"
- 快速了解当天需要处理的事项

返回信息：待办ID、标题、状态

使用建议：
- 每天开始工作时先查看今日待办
- 根据优先级安排处理顺序
- 完成后及时使用 todo_complete 标记

提示：如果需要查看所有待办（不仅是今天的），请使用 todo_list`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, input TodoTodayInput) (*mcp.CallToolResult, any, error) {
		todos, err := bs.TodoService.ListToday(ctx)
		if err != nil {
			return NewErrorResult(err.Error()), nil, nil
		}
		if len(todos) == 0 {
			return NewTextResult("今日暂无待办事项"), nil, nil
		}
		result := fmt.Sprintf("今日待办事项 (%s):\n", time.Now().Format("2006-01-02"))
		for _, t := range todos {
			status := getTodoStatusText(t.Status)
			result += fmt.Sprintf("- [%d] %s (%s)\n", t.ID, t.Title, status)
		}
		return NewTextResult(result), nil, nil
	})
}

// getTodoStatusText 获取待办状态文本
func getTodoStatusText(status types.TodoStatus) string {
	switch status {
	case types.TodoStatusPending:
		return "待处理"
	case types.TodoStatusInProgress:
		return "进行中"
	case types.TodoStatusCompleted:
		return "已完成"
	case types.TodoStatusCancelled:
		return "已取消"
	default:
		return "未知"
	}
}

// getPriorityText 获取优先级文本
func getPriorityText(priority types.Priority) string {
	switch priority {
	case types.TodoPriorityLow:
		return "低"
	case types.TodoPriorityMedium:
		return "中"
	case types.TodoPriorityHigh:
		return "高"
	case types.TodoPriorityUrgent:
		return "紧急"
	default:
		return "未知"
	}
}
