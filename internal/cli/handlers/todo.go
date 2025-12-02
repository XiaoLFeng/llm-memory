package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/XiaoLFeng/llm-memory/internal/cli"
	"github.com/XiaoLFeng/llm-memory/internal/cli/output"
	"github.com/XiaoLFeng/llm-memory/pkg/types"
	"github.com/XiaoLFeng/llm-memory/startup"
)

// TodoHandler TODO 命令处理器
// 嘿嘿~ 处理所有待办相关的 CLI 命令！✅
type TodoHandler struct {
	bs *startup.Bootstrap
}

// NewTodoHandler 创建 TODO 处理器
func NewTodoHandler(bs *startup.Bootstrap) *TodoHandler {
	return &TodoHandler{bs: bs}
}

// List 列出所有待办
// 呀~ 展示所有待办事项！✨
func (h *TodoHandler) List(ctx context.Context) error {
	todos, err := h.bs.TodoService.ListTodos(ctx)
	if err != nil {
		return err
	}

	if len(todos) == 0 {
		cli.PrintInfo("暂无待办事项~")
		return nil
	}

	cli.PrintTitle("📝 待办事项列表")
	table := output.NewTable("ID", "标题", "状态", "优先级")
	for _, t := range todos {
		table.AddRow(
			fmt.Sprintf("%d", t.ID),
			t.Title,
			getTodoStatusText(t.Status),
			getPriorityText(t.Priority),
		)
	}
	table.Print()

	return nil
}

// Today 获取今日待办
// 嘿嘿~ 查看今天要做的事！📅
func (h *TodoHandler) Today(ctx context.Context) error {
	todos, err := h.bs.TodoService.ListToday(ctx)
	if err != nil {
		return err
	}

	if len(todos) == 0 {
		cli.PrintInfo("今日暂无待办事项~ 🎉")
		return nil
	}

	cli.PrintTitle(fmt.Sprintf("📅 今日待办 (%s)", time.Now().Format("2006-01-02")))
	table := output.NewTable("ID", "标题", "状态", "优先级")
	for _, t := range todos {
		table.AddRow(
			fmt.Sprintf("%d", t.ID),
			t.Title,
			getTodoStatusText(t.Status),
			getPriorityText(t.Priority),
		)
	}
	table.Print()

	return nil
}

// Create 创建待办
// 呀~ 创建新的待办事项！💫
func (h *TodoHandler) Create(ctx context.Context, title, description string, priority int) error {
	p := types.Priority(priority)
	if p == 0 {
		p = types.TodoPriorityMedium
	}

	todo, err := h.bs.TodoService.CreateTodo(ctx, title, description, p, nil)
	if err != nil {
		return err
	}

	cli.PrintSuccess(fmt.Sprintf("待办创建成功！ID: %d, 标题: %s", todo.ID, todo.Title))
	return nil
}

// Complete 完成待办
func (h *TodoHandler) Complete(ctx context.Context, id int) error {
	if err := h.bs.TodoService.CompleteTodo(ctx, id); err != nil {
		return err
	}

	cli.PrintSuccess(fmt.Sprintf("待办 %d 已完成", id))
	return nil
}

// Start 开始待办
func (h *TodoHandler) Start(ctx context.Context, id int) error {
	if err := h.bs.TodoService.StartTodo(ctx, id); err != nil {
		return err
	}

	cli.PrintSuccess(fmt.Sprintf("待办 %d 已开始", id))
	return nil
}

// Delete 删除待办
func (h *TodoHandler) Delete(ctx context.Context, id int) error {
	if err := h.bs.TodoService.DeleteTodo(ctx, id); err != nil {
		return err
	}

	cli.PrintSuccess(fmt.Sprintf("待办 %d 已删除", id))
	return nil
}

// Get 获取待办详情
// 嗯嗯！查看待办的详细信息！📝
func (h *TodoHandler) Get(ctx context.Context, id int) error {
	todo, err := h.bs.TodoService.GetTodo(ctx, id)
	if err != nil {
		return err
	}

	cli.PrintTitle("✅ 待办详情")
	fmt.Printf("ID:       %d\n", todo.ID)
	fmt.Printf("标题:     %s\n", todo.Title)
	fmt.Printf("状态:     %s\n", getTodoStatusText(todo.Status))
	fmt.Printf("优先级:   %s\n", getPriorityText(todo.Priority))
	if todo.DueDate != nil {
		fmt.Printf("截止日期: %s\n", todo.DueDate.Format("2006-01-02"))
	}
	if todo.CompletedAt != nil {
		fmt.Printf("完成时间: %s\n", todo.CompletedAt.Format("2006-01-02 15:04:05"))
	}
	fmt.Printf("创建时间: %s\n", todo.CreatedAt.Format("2006-01-02 15:04:05"))
	if todo.Description != "" {
		fmt.Println("\n描述:")
		fmt.Println(todo.Description)
	}

	return nil
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
