package tools

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/XiaoLFeng/llm-memory/internal/models/dto"
	"github.com/XiaoLFeng/llm-memory/internal/models/entity"
	"github.com/XiaoLFeng/llm-memory/startup"
)

// TodoListInput todo_list 工具输入
type TodoListInput struct {
	Scope string `json:"scope,omitempty" jsonschema:"作用域过滤(personal/group/global/all)，默认all显示全部"`
}

// TodoCreateInput todo_create 工具输入
type TodoCreateInput struct {
	Title       string `json:"title" jsonschema:"待办标题，简洁描述任务"`
	Description string `json:"description,omitempty" jsonschema:"待办的详细描述"`
	Priority    int    `json:"priority,omitempty" jsonschema:"优先级(1低/2中/3高/4紧急)，默认2"`
	Scope       string `json:"scope,omitempty" jsonschema:"保存到哪个作用域(personal/group/global)，默认global"`
}

// TodoCompleteInput todo_complete 工具输入
type TodoCompleteInput struct {
	ID int64 `json:"id" jsonschema:"要完成的待办事项ID"`
}

// RegisterTodoTools 注册 TODO 管理工具
func RegisterTodoTools(server *mcp.Server, bs *startup.Bootstrap) {
	// todo_list - 列出所有待办
	mcp.AddTool(server, &mcp.Tool{
		Name:        "todo_list",
		Description: `列出所有待办及状态。scope: personal/group/global/all，默认all=使用当前作用域集合；未指定也会落在 currentScope（无则 global）。`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, input TodoListInput) (*mcp.CallToolResult, any, error) {
		// 构建作用域上下文
		scopeCtx := buildScopeContext(input.Scope, bs)

		todos, err := bs.ToDoService.ListToDosByScope(ctx, input.Scope, scopeCtx)
		if err != nil {
			return NewErrorResult(err.Error()), nil, nil
		}
		if len(todos) == 0 {
			return NewTextResult("暂无待办事项"), nil, nil
		}
		result := "待办事项列表:\n"
		for _, t := range todos {
			status := getToDoStatusText(t.Status)
			priority := getToDoPriorityText(t.Priority)
			scopeTag := getScopeTagFromPathID(t.PathID)
			result += fmt.Sprintf("- [%d] %s (%s, %s) %s\n", t.ID, t.Title, status, priority, scopeTag)
		}
		return NewTextResult(result), nil, nil
	})

	// todo_create - 创建待办
	mcp.AddTool(server, &mcp.Tool{
		Name:        "todo_create",
		Description: `创建待办，适合可立即执行或短周期的单一步行动。必填: title。可选: description、priority(1低/2中/3高/4紧急，默认2)、scope。若是多步骤需跟踪的目标请用 plan_create；长期背景/事实请用 memory_create。未指定 scope 默认使用 currentScope（缺省则 global）。`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, input TodoCreateInput) (*mcp.CallToolResult, any, error) {
		// 默认优先级
		priority := input.Priority
		if priority == 0 {
			priority = 2 // 默认中等优先级
		}

		// 构建创建 DTO
		createDTO := &dto.ToDoCreateDTO{
			Title:       input.Title,
			Description: input.Description,
			Priority:    priority,
			Scope:       input.Scope,
		}

		// 构建作用域上下文
		scopeCtx := buildScopeContext(input.Scope, bs)

		todo, err := bs.ToDoService.CreateToDo(ctx, createDTO, scopeCtx)
		if err != nil {
			return NewErrorResult(err.Error()), nil, nil
		}
		scopeTag := getScopeTagFromPathID(todo.PathID)
		return NewTextResult(fmt.Sprintf("待办事项创建成功! ID: %d, 标题: %s %s", todo.ID, todo.Title, scopeTag)), nil, nil
	})

	// todo_complete - 完成待办
	mcp.AddTool(server, &mcp.Tool{
		Name:        "todo_complete",
		Description: `标记待办为已完成。`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, input TodoCompleteInput) (*mcp.CallToolResult, any, error) {
		if err := bs.ToDoService.CompleteToDo(ctx, input.ID); err != nil {
			return NewErrorResult(err.Error()), nil, nil
		}
		return NewTextResult(fmt.Sprintf("待办事项 %d 已标记为完成", input.ID)), nil, nil
	})
}

// getToDoStatusText 获取待办状态文本
func getToDoStatusText(status entity.ToDoStatus) string {
	switch status {
	case entity.ToDoStatusPending:
		return "待处理"
	case entity.ToDoStatusInProgress:
		return "进行中"
	case entity.ToDoStatusCompleted:
		return "已完成"
	case entity.ToDoStatusCancelled:
		return "已取消"
	default:
		return "未知"
	}
}

// getToDoPriorityText 获取优先级文本
func getToDoPriorityText(priority entity.ToDoPriority) string {
	switch priority {
	case entity.ToDoPriorityLow:
		return "低"
	case entity.ToDoPriorityMedium:
		return "中"
	case entity.ToDoPriorityHigh:
		return "高"
	case entity.ToDoPriorityUrgent:
		return "紧急"
	default:
		return "未知"
	}
}
