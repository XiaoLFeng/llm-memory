package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/XiaoLFeng/llm-memory/internal/models/dto"
	"github.com/XiaoLFeng/llm-memory/internal/models/entity"
	"github.com/XiaoLFeng/llm-memory/startup"
)

// TodoListInput todo_list 工具输入
type TodoListInput struct {
	Scope string `json:"scope,omitempty" jsonschema:"作用域过滤 personal group all 默认all显示全部"`
}

// TodoCreateItem 批量创建的待办项
type TodoCreateItem struct {
	Code        string `json:"code" jsonschema:"待办唯一标识码"`
	PlanCode    string `json:"plan_code" jsonschema:"所属计划的标识码（必填）"`
	Title       string `json:"title" jsonschema:"待办标题 简洁描述任务"`
	Description string `json:"description,omitempty" jsonschema:"待办的详细描述"`
	Priority    int    `json:"priority,omitempty" jsonschema:"优先级 1低2中3高4紧急 默认2"`
}

// TodoBatchCreateInput todo_batch_create 工具输入
type TodoBatchCreateInput struct {
	Items []TodoCreateItem `json:"items" jsonschema:"待办事项列表最多100个，每个待办必须指定plan_code"`
}

// TodoBatchOperationInput 批量操作输入（完成/取消/删除）
type TodoBatchOperationInput struct {
	Codes []string `json:"codes" jsonschema:"待办事项代码列表最多100个"`
}

// TodoBatchOperationResult 批量操作结果
type TodoBatchOperationResult struct {
	SuccessCount int                `json:"success_count"`
	FailCount    int                `json:"fail_count"`
	Failures     []TodoBatchFailure `json:"failures,omitempty"`
}

type TodoBatchFailure struct {
	Code  string `json:"code"`
	Error string `json:"error"`
}

// TodoBatchUpdateInput todo_batch_update 工具输入
type TodoBatchUpdateInput struct {
	Items []TodoUpdateItem `json:"items" jsonschema:"待办事项更新列表最多100个"`
}

type TodoUpdateItem struct {
	Code        string `json:"code" jsonschema:"要更新的待办事项代码"`
	Title       string `json:"title,omitempty" jsonschema:"新的待办标题"`
	Description string `json:"description,omitempty" jsonschema:"新的待办描述"`
	Priority    int    `json:"priority,omitempty" jsonschema:"新的优先级 1低2中3高4紧急"`
	Status      int    `json:"status,omitempty" jsonschema:"新的状态 0待处理1进行中2已完成3已取消"`
}

// TodoBatchStartInput todo_batch_start 工具输入
type TodoBatchStartInput struct {
	Codes []string `json:"codes" jsonschema:"要开始的待办事项代码列表最多100个"`
}

// TodoFinalInput todo_final 工具输入
type TodoFinalInput struct {
	Scope string `json:"scope,omitempty" jsonschema:"作用域过滤 personal group all 默认all显示全部"`
}

// validateTodoOwnership 验证待办所有权权限
func validateTodoOwnership(ctx context.Context, bs *startup.Bootstrap, code string) (*entity.ToDo, error) {
	todo, err := bs.ToDoService.GetToDo(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("待办不存在: %s", code)
	}

	// 权限检查：只能操作自己作用域内的待办
	scope := bs.CurrentScope
	if scope == nil {
		return nil, fmt.Errorf("当前作用域为空")
	}

	// 检查个人权限
	if todo.PathID > 0 {
		if scope.IncludePersonal && todo.PathID == scope.PathID {
			return todo, nil
		}
		// 检查组权限
		if scope.IncludeGroup {
			for _, groupPathID := range scope.GroupPathIDs {
				if todo.PathID == groupPathID {
					return todo, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("无权限操作待办: %s", code)
}

// validateBatchSize 验证批量操作大小
func validateBatchSize(count int) error {
	if count > 100 {
		return fmt.Errorf("批量操作最多支持100个项目")
	}
	if count == 0 {
		return fmt.Errorf("批量操作至少需要1个项目")
	}
	return nil
}

// formatBatchResult 格式化批量操作结果为混合模式
func formatBatchResult(result *TodoBatchOperationResult, operation string) string {
	if result.FailCount == 0 {
		return fmt.Sprintf("✅ %s成功! 共处理 %d 个待办事项", operation, result.SuccessCount)
	}
	if result.SuccessCount == 0 {
		return fmt.Sprintf("❌ %s失败! 所有 %d 个待办事项都无法处理:\n%s",
			operation, result.FailCount, formatFailures(result.Failures))
	}
	return fmt.Sprintf("⚠️ %s部分完成! 成功 %d 个，失败 %d 个:\n%s",
		operation, result.SuccessCount, result.FailCount, formatFailures(result.Failures))
}

// formatFailures 格式化失败信息
func formatFailures(failures []TodoBatchFailure) string {
	var sb strings.Builder
	for _, failure := range failures {
		sb.WriteString(fmt.Sprintf("  • %s: %s\n", failure.Code, failure.Error))
	}
	return sb.String()
}

// RegisterTodoTools 注册 TODO 管理工具
func RegisterTodoTools(server *mcp.Server, bs *startup.Bootstrap) {
	// todo_list - 列出所有待办
	mcp.AddTool(server, &mcp.Tool{
		Name: "todo_list",
		Description: `列出所有待办及状态。每个 Todo 都归属于一个 Plan。
scope参数说明（安全隔离）：
  - personal: 仅当前路径的项目数据
  - group: 仅当前小组的数据（需已加入小组）
  - all/省略: 当前路径 + 小组数据（默认，权限隔离）`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, input TodoListInput) (*mcp.CallToolResult, any, error) {
		// 构建作用域上下文
		scopeCtx := getScopeContext(bs)

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
			scopeTag := getScopeTagWithContext(t.PathID, bs.CurrentScope)
			// 获取 Plan Code
			planCode, _ := bs.ToDoService.GetPlanCodeByTodoID(ctx, t.ID)
			// 格式: title - description (如果有描述)
			titlePart := t.Title
			if t.Description != "" {
				desc := t.Description
				if len(desc) > 60 {
					desc = desc[:57] + "..."
				}
				titlePart = t.Title + " - " + desc
			}
			result += fmt.Sprintf("- [%s] %s (计划:%s, %s, %s) %s\n", t.Code, titlePart, planCode, status, priority, scopeTag)
		}
		return NewTextResult(result), nil, nil
	})

	// todo_batch_create - 批量创建待办
	mcp.AddTool(server, &mcp.Tool{
		Name: "todo_batch_create",
		Description: `批量创建待办事项，提高AI处理效率。支持最多100个待办项的批量创建。
重要：每个待办必须指定 plan_code（所属计划的标识码），Todo 必须归属于一个 Plan。
返回混合模式结果：成功显示统计，失败显示详情。`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, input TodoBatchCreateInput) (*mcp.CallToolResult, any, error) {
		// 验证批量大小
		if err := validateBatchSize(len(input.Items)); err != nil {
			return NewErrorResult(err.Error()), nil, nil
		}

		result := &TodoBatchOperationResult{}
		scopeCtx := getScopeContext(bs)

		// 批量创建
		for _, item := range input.Items {
			// 验证 plan_code 必填
			if strings.TrimSpace(item.PlanCode) == "" {
				result.FailCount++
				result.Failures = append(result.Failures, TodoBatchFailure{
					Code:  item.Code,
					Error: "plan_code 是必填项，Todo 必须归属于一个 Plan",
				})
				continue
			}

			// 默认优先级
			priority := item.Priority
			if priority == 0 {
				priority = 2 // 默认中等优先级
			}

			createDTO := &dto.ToDoCreateDTO{
				Code:        item.Code,
				PlanCode:    item.PlanCode,
				Title:       item.Title,
				Description: item.Description,
				Priority:    priority,
			}

			_, err := bs.ToDoService.CreateToDo(ctx, createDTO, scopeCtx)
			if err != nil {
				result.FailCount++
				result.Failures = append(result.Failures, TodoBatchFailure{
					Code:  item.Code,
					Error: err.Error(),
				})
			} else {
				result.SuccessCount++
			}
		}

		response := formatBatchResult(result, "批量创建")
		return NewTextResult(response), result, nil
	})

	// todo_batch_complete - 批量完成待办
	mcp.AddTool(server, &mcp.Tool{
		Name:        "todo_batch_complete",
		Description: `批量标记待办事项为已完成。支持最多100个待办的批量完成操作。返回混合模式结果。`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, input TodoBatchOperationInput) (*mcp.CallToolResult, any, error) {
		// 验证批量大小
		if err := validateBatchSize(len(input.Codes)); err != nil {
			return NewErrorResult(err.Error()), nil, nil
		}

		result := &TodoBatchOperationResult{}

		// 批量完成
		for _, code := range input.Codes {
			// 权限验证
			_, err := validateTodoOwnership(ctx, bs, code)
			if err != nil {
				result.FailCount++
				result.Failures = append(result.Failures, TodoBatchFailure{
					Code:  code,
					Error: err.Error(),
				})
				continue
			}

			// 标记完成
			if err := bs.ToDoService.CompleteToDo(ctx, code); err != nil {
				result.FailCount++
				result.Failures = append(result.Failures, TodoBatchFailure{
					Code:  code,
					Error: err.Error(),
				})
			} else {
				result.SuccessCount++
			}
		}

		response := formatBatchResult(result, "批量完成")
		return NewTextResult(response), result, nil
	})

	// todo_batch_cancel - 批量取消待办
	mcp.AddTool(server, &mcp.Tool{
		Name:        "todo_batch_cancel",
		Description: `批量标记待办事项为已取消。支持最多100个待办的批量取消操作。返回混合模式结果。`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, input TodoBatchOperationInput) (*mcp.CallToolResult, any, error) {
		// 验证批量大小
		if err := validateBatchSize(len(input.Codes)); err != nil {
			return NewErrorResult(err.Error()), nil, nil
		}

		result := &TodoBatchOperationResult{}

		// 批量取消
		for _, code := range input.Codes {
			// 权限验证
			_, err := validateTodoOwnership(ctx, bs, code)
			if err != nil {
				result.FailCount++
				result.Failures = append(result.Failures, TodoBatchFailure{
					Code:  code,
					Error: err.Error(),
				})
				continue
			}

			// 标记取消
			if err := bs.ToDoService.CancelToDo(ctx, code); err != nil {
				result.FailCount++
				result.Failures = append(result.Failures, TodoBatchFailure{
					Code:  code,
					Error: err.Error(),
				})
			} else {
				result.SuccessCount++
			}
		}

		response := formatBatchResult(result, "批量取消")
		return NewTextResult(response), result, nil
	})

	// todo_batch_start - 批量开始待办
	mcp.AddTool(server, &mcp.Tool{
		Name:        "todo_batch_start",
		Description: `批量标记待办事项为进行中状态。支持最多100个待办的批量开始操作。返回混合模式结果。`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, input TodoBatchStartInput) (*mcp.CallToolResult, any, error) {
		// 验证批量大小
		if err := validateBatchSize(len(input.Codes)); err != nil {
			return NewErrorResult(err.Error()), nil, nil
		}

		result := &TodoBatchOperationResult{}

		// 批量开始
		for _, code := range input.Codes {
			// 权限验证
			_, err := validateTodoOwnership(ctx, bs, code)
			if err != nil {
				result.FailCount++
				result.Failures = append(result.Failures, TodoBatchFailure{
					Code:  code,
					Error: err.Error(),
				})
				continue
			}

			// 标记开始
			if err := bs.ToDoService.StartToDo(ctx, code); err != nil {
				result.FailCount++
				result.Failures = append(result.Failures, TodoBatchFailure{
					Code:  code,
					Error: err.Error(),
				})
			} else {
				result.SuccessCount++
			}
		}

		response := formatBatchResult(result, "批量开始")
		return NewTextResult(response), result, nil
	})

	// todo_batch_update - 批量更新待办
	mcp.AddTool(server, &mcp.Tool{
		Name:        "todo_batch_update",
		Description: `批量更新待办事项的标题、描述、优先级或状态。支持最多100个待办的批量更新。返回混合模式结果。`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, input TodoBatchUpdateInput) (*mcp.CallToolResult, any, error) {
		// 验证批量大小
		if err := validateBatchSize(len(input.Items)); err != nil {
			return NewErrorResult(err.Error()), nil, nil
		}

		result := &TodoBatchOperationResult{}

		// 批量更新
		for _, item := range input.Items {
			// 权限验证
			_, err := validateTodoOwnership(ctx, bs, item.Code)
			if err != nil {
				result.FailCount++
				result.Failures = append(result.Failures, TodoBatchFailure{
					Code:  item.Code,
					Error: err.Error(),
				})
				continue
			}

			// 检查是否有更新内容
			hasUpdates := false
			updateDTO := &dto.ToDoUpdateDTO{
				Code: item.Code,
			}

			if item.Title != "" {
				updateDTO.Title = &item.Title
				hasUpdates = true
			}
			if item.Description != "" {
				updateDTO.Description = &item.Description
				hasUpdates = true
			}
			if item.Priority > 0 && item.Priority <= 4 {
				updateDTO.Priority = &item.Priority
				hasUpdates = true
			}
			if item.Status >= 0 && item.Status <= 3 {
				updateDTO.Status = &item.Status
				hasUpdates = true
			}

			if !hasUpdates {
				result.FailCount++
				result.Failures = append(result.Failures, TodoBatchFailure{
					Code:  item.Code,
					Error: "至少需要提供一个要更新的字段",
				})
				continue
			}

			// 更新待办
			if err := bs.ToDoService.UpdateToDo(ctx, updateDTO); err != nil {
				result.FailCount++
				result.Failures = append(result.Failures, TodoBatchFailure{
					Code:  item.Code,
					Error: err.Error(),
				})
			} else {
				result.SuccessCount++
			}
		}

		response := formatBatchResult(result, "批量更新")
		return NewTextResult(response), result, nil
	})

	// todo_final - 删除所有待办
	mcp.AddTool(server, &mcp.Tool{
		Name: "todo_final",
		Description: `删除当前作用域内的所有待办事项。这是一个清理工具，会直接删除指定作用域内的所有待办（不可恢复）。
删除逻辑：
  - 未加入小组：删除当前路径的项目待办
  - 已加入小组：删除小组内所有路径的待办`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, input TodoFinalInput) (*mcp.CallToolResult, any, error) {
		// 构建作用域上下文
		scopeCtx := getScopeContext(bs)

		// 删除所有待办
		deletedCount, err := bs.ToDoService.DeleteAllByScope(ctx, input.Scope, scopeCtx)
		if err != nil {
			return NewErrorResult(err.Error()), nil, nil
		}

		if deletedCount == 0 {
			return NewTextResult("当前作用域内没有待办事项需要删除"), nil, nil
		}

		return NewTextResult(fmt.Sprintf("已删除 %d 个待办事项", deletedCount)), nil, nil
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
