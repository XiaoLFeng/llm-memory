package handlers

import (
	"context"
	"fmt"

	"github.com/XiaoLFeng/llm-memory/internal/cli"
	"github.com/XiaoLFeng/llm-memory/internal/cli/output"
	"github.com/XiaoLFeng/llm-memory/internal/models/dto"
	"github.com/XiaoLFeng/llm-memory/internal/models/entity"
	"github.com/XiaoLFeng/llm-memory/startup"
)

// TodoHandler TODO 命令处理器
type TodoHandler struct {
	bs *startup.Bootstrap
}

// NewTodoHandler 创建 TODO 处理器
func NewTodoHandler(bs *startup.Bootstrap) *TodoHandler {
	return &TodoHandler{bs: bs}
}

// List 列出所有待办
func (h *TodoHandler) List(ctx context.Context) error {
	// 使用 ListToDosByScope 确保权限隔离：全局 + 当前路径相关
	todos, err := h.bs.ToDoService.ListToDosByScope(ctx, "all", h.bs.CurrentScope)
	if err != nil {
		return err
	}

	if len(todos) == 0 {
		cli.PrintInfo("暂无待办事项~")
		return nil
	}

	cli.PrintTitle(cli.IconTodo + " 待办事项列表")
	table := output.NewTable("标识码", "标题", "状态", "优先级")
	for _, t := range todos {
		table.AddRow(
			t.Code,
			t.Title,
			getToDoStatusText(t.Status),
			getToDoStatusPriorityText(t.Priority),
		)
	}
	table.Print()

	return nil
}

// Create 创建待办
func (h *TodoHandler) Create(ctx context.Context, code, title, description string, priority int, global bool) error {
	if priority == 0 {
		priority = int(entity.ToDoPriorityMedium)
	}

	createDTO := &dto.ToDoCreateDTO{
		Code:        code,
		Title:       title,
		Description: description,
		Priority:    priority,
	}

	todo, err := h.bs.ToDoService.CreateToDo(ctx, createDTO, h.bs.CurrentScope)
	if err != nil {
		return err
	}

	cli.PrintSuccess(fmt.Sprintf("待办创建成功！标识码: %s, 标题: %s", todo.Code, todo.Title))
	return nil
}

// Complete 完成待办
func (h *TodoHandler) Complete(ctx context.Context, code string) error {
	if err := h.bs.ToDoService.CompleteToDo(ctx, code); err != nil {
		return err
	}

	cli.PrintSuccess(fmt.Sprintf("待办 %s 已完成", code))
	return nil
}

// Start 开始待办
func (h *TodoHandler) Start(ctx context.Context, code string) error {
	if err := h.bs.ToDoService.StartToDo(ctx, code); err != nil {
		return err
	}

	cli.PrintSuccess(fmt.Sprintf("待办 %s 已开始", code))
	return nil
}

// Delete 删除待办
func (h *TodoHandler) Delete(ctx context.Context, code string) error {
	if err := h.bs.ToDoService.DeleteToDo(ctx, code); err != nil {
		return err
	}

	cli.PrintSuccess(fmt.Sprintf("待办 %s 已删除", code))
	return nil
}

// Final 删除当前作用域的所有待办（清理功能）
func (h *TodoHandler) Final(ctx context.Context) error {
	deletedCount, err := h.bs.ToDoService.DeleteAllByScope(ctx, "all", h.bs.CurrentScope)
	if err != nil {
		return err
	}

	if deletedCount == 0 {
		cli.PrintInfo("当前作用域内没有待办事项需要删除")
		return nil
	}

	cli.PrintSuccess(fmt.Sprintf("已删除 %d 个待办事项", deletedCount))
	return nil
}

// Get 获取待办详情
func (h *TodoHandler) Get(ctx context.Context, code string) error {
	todo, err := h.bs.ToDoService.GetToDo(ctx, code)
	if err != nil {
		return err
	}

	cli.PrintTitle(cli.IconCheck + " 待办详情")
	fmt.Printf("标识码:   %s\n", todo.Code)
	fmt.Printf("标题:     %s\n", todo.Title)
	fmt.Printf("状态:     %s\n", getToDoStatusText(todo.Status))
	fmt.Printf("优先级:   %s\n", getToDoStatusPriorityText(todo.Priority))
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

// getToDoStatusPriorityText 获取优先级文本
func getToDoStatusPriorityText(priority entity.ToDoPriority) string {
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

// BatchCreate 批量创建待办
func (h *TodoHandler) BatchCreate(ctx context.Context, items []dto.ToDoCreateDTO) error {
	if err := h.validateBatchItems(len(items)); err != nil {
		return err
	}

	batchDTO := &dto.ToDoBatchCreateDTO{Items: items}
	result, err := h.bs.ToDoService.BatchCreateToDos(ctx, batchDTO, h.bs.CurrentScope)
	if err != nil {
		return err
	}

	h.printBatchResult(result, "批量创建")
	return nil
}

// BatchComplete 批量完成待办
func (h *TodoHandler) BatchComplete(ctx context.Context, codes []string) error {
	if err := h.validateBatchCodes(codes); err != nil {
		return err
	}

	batchDTO := &dto.ToDoBatchCompleteDTO{Codes: codes}
	result, err := h.bs.ToDoService.BatchCompleteToDos(ctx, batchDTO)
	if err != nil {
		return err
	}

	h.printBatchResult(result, "批量完成")
	return nil
}

// BatchStart 批量开始待办
func (h *TodoHandler) BatchStart(ctx context.Context, codes []string) error {
	if err := h.validateBatchCodes(codes); err != nil {
		return err
	}

	progressDTO := &dto.ToDoBatchProgressDTO{
		Codes:  codes,
		Status: int(entity.ToDoStatusInProgress),
	}
	result, err := h.bs.ToDoService.BatchUpdateProgress(ctx, progressDTO)
	if err != nil {
		return err
	}

	h.printBatchResult(result, "批量开始")
	return nil
}

// BatchCancel 批量取消待办
func (h *TodoHandler) BatchCancel(ctx context.Context, codes []string) error {
	if err := h.validateBatchCodes(codes); err != nil {
		return err
	}

	progressDTO := &dto.ToDoBatchProgressDTO{
		Codes:  codes,
		Status: int(entity.ToDoStatusCancelled),
	}
	result, err := h.bs.ToDoService.BatchUpdateProgress(ctx, progressDTO)
	if err != nil {
		return err
	}

	h.printBatchResult(result, "批量取消")
	return nil
}

// BatchDelete 批量删除待办
func (h *TodoHandler) BatchDelete(ctx context.Context, codes []string) error {
	if err := h.validateBatchCodes(codes); err != nil {
		return err
	}

	batchDTO := &dto.ToDoBatchDeleteDTO{Codes: codes}
	result, err := h.bs.ToDoService.BatchDeleteToDos(ctx, batchDTO)
	if err != nil {
		return err
	}

	h.printBatchResult(result, "批量删除")
	return nil
}

// BatchUpdate 批量更新待办
func (h *TodoHandler) BatchUpdate(ctx context.Context, items []dto.ToDoUpdateDTO) error {
	if err := h.validateBatchItems(len(items)); err != nil {
		return err
	}

	batchDTO := &dto.ToDoBatchUpdateDTO{Items: items}
	result, err := h.bs.ToDoService.BatchUpdateToDos(ctx, batchDTO)
	if err != nil {
		return err
	}

	h.printBatchResult(result, "批量更新")
	return nil
}

// validateBatchCodes 验证批量 codes
func (h *TodoHandler) validateBatchCodes(codes []string) error {
	if len(codes) == 0 {
		return fmt.Errorf("批量操作至少需要1个项目")
	}
	if len(codes) > 100 {
		return fmt.Errorf("批量操作最多支持100个项目")
	}
	return nil
}

// validateBatchItems 验证批量项目数量
func (h *TodoHandler) validateBatchItems(count int) error {
	if count == 0 {
		return fmt.Errorf("批量操作至少需要1个项目")
	}
	if count > 100 {
		return fmt.Errorf("批量操作最多支持100个项目")
	}
	return nil
}

// printBatchResult 打印批量操作结果（混合模式）
func (h *TodoHandler) printBatchResult(result *dto.ToDoBatchResultDTO, operation string) {
	if result.Failed == 0 {
		// 全部成功
		cli.PrintSuccess(fmt.Sprintf("%s成功! 共处理 %d 个待办事项", operation, result.Succeeded))
		return
	}

	if result.Succeeded == 0 {
		// 全部失败
		cli.PrintError(fmt.Sprintf("%s失败! 所有 %d 个待办事项都无法处理:", operation, result.Failed))
		for _, errMsg := range result.Errors {
			fmt.Printf("  • %s\n", errMsg)
		}
		return
	}

	// 部分成功
	cli.PrintWarning(fmt.Sprintf("%s部分完成! 成功 %d 个，失败 %d 个:", operation, result.Succeeded, result.Failed))
	for _, errMsg := range result.Errors {
		fmt.Printf("  • %s\n", errMsg)
	}
}
