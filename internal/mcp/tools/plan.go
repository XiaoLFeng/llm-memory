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

// PlanListInput plan_list 工具输入
type PlanListInput struct {
	Scope string `json:"scope,omitempty" jsonschema:"作用域过滤(personal/group/all)，默认all显示全部"`
}

// PlanCreateInput plan_create 工具输入
type PlanCreateInput struct {
	Code        string `json:"code" jsonschema:"计划唯一标识码"`
	Title       string `json:"title" jsonschema:"计划标题，简洁描述计划目标"`
	Description string `json:"description" jsonschema:"计划的详细描述，包含具体步骤和目标"`
	Content     string `json:"content" jsonschema:"计划的详细内容，支持 Markdown 格式"`
	Scope       string `json:"scope,omitempty" jsonschema:"查询筛选仍可用的作用域 personal/group/all"`
}

// PlanGetInput plan_get 工具输入
type PlanGetInput struct {
	Code string `json:"code" jsonschema:"要获取的计划code"`
}

// PlanUpdateInput plan_update 工具输入
type PlanUpdateInput struct {
	Code        string  `json:"code" jsonschema:"要更新的计划code"`
	Title       *string `json:"title,omitempty" jsonschema:"新标题（可选）"`
	Description *string `json:"description,omitempty" jsonschema:"新描述（可选）"`
	Content     *string `json:"content,omitempty" jsonschema:"新内容（可选），支持 Markdown 格式"`
	Progress    *int    `json:"progress,omitempty" jsonschema:"完成进度 0-100（可选），系统自动调整状态：0=待开始，1-99=进行中，100=已完成"`
}

// RegisterPlanTools 注册计划管理工具
func RegisterPlanTools(server *mcp.Server, bs *startup.Bootstrap) {
	// plan_list - 列出所有计划
	mcp.AddTool(server, &mcp.Tool{
		Name: "plan_list",
		Description: `列出所有计划及进度状态。scope参数说明（安全隔离）：
  - personal: 仅当前路径的私有数据
  - group: 仅当前小组的数据（需已加入小组）
  - all/省略: 当前路径 + 小组数据（默认，权限隔离）`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, input PlanListInput) (*mcp.CallToolResult, any, error) {
		// 构建作用域上下文
		scopeCtx := buildScopeContext(input.Scope, bs)

		plans, err := bs.PlanService.ListPlansByScope(ctx, input.Scope, scopeCtx)
		if err != nil {
			return NewErrorResult(err.Error()), nil, nil
		}
		if len(plans) == 0 {
			return NewTextResult("暂无计划"), nil, nil
		}
		result := "计划列表:\n"
		for _, p := range plans {
			status := getPlanStatusText(p.Status)
			scopeTag := getScopeTagWithContext(p.PathID, bs.CurrentScope)
			result += fmt.Sprintf("- [%s] %s (%s, 进度: %d%%) %s\n", p.Code, p.Title, status, p.Progress, scopeTag)
		}
		return NewTextResult(result), nil, nil
	})

	// plan_create - 创建新计划
	mcp.AddTool(server, &mcp.Tool{
		Name:        "plan_create",
		Description: `创建计划，用于"需要跟踪进度的多步骤目标"。必填: title、description、content(Markdown)。短动作请用 todo_create；长期事实请用 memory_create。scope 参数仅用于列表筛选。`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, input PlanCreateInput) (*mcp.CallToolResult, any, error) {
		// 构建创建 DTO
		createDTO := &dto.PlanCreateDTO{
			Code:        input.Code,
			Title:       input.Title,
			Description: input.Description,
			Content:     input.Content,
		}

		// 构建作用域上下文
		scopeCtx := buildScopeContext(input.Scope, bs)

		plan, err := bs.PlanService.CreatePlan(ctx, createDTO, scopeCtx)
		if err != nil {
			return NewErrorResult(err.Error()), nil, nil
		}
		scopeTag := getScopeTagWithContext(plan.PathID, bs.CurrentScope)
		return NewTextResult(fmt.Sprintf("计划创建成功! Code: %s, 标题: %s %s", plan.Code, plan.Title, scopeTag)), nil, nil
	})

	// plan_get - 获取计划详情
	mcp.AddTool(server, &mcp.Tool{
		Name:        "plan_get",
		Description: `获取指定code计划的完整详情，包括标题、描述、内容、进度、子任务等。`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, input PlanGetInput) (*mcp.CallToolResult, any, error) {
		plan, err := bs.PlanService.GetPlan(ctx, input.Code)
		if err != nil {
			return NewErrorResult(err.Error()), nil, nil
		}

		scopeTag := getScopeTagWithContext(plan.PathID, bs.CurrentScope)

		var sb strings.Builder
		sb.WriteString("计划详情:\n")
		sb.WriteString(fmt.Sprintf("Code: %s\n", plan.Code))
		sb.WriteString(fmt.Sprintf("标题: %s\n", plan.Title))
		sb.WriteString(fmt.Sprintf("状态: %s\n", getPlanStatusText(plan.Status)))
		sb.WriteString(fmt.Sprintf("进度: %d%%\n", plan.Progress))
		sb.WriteString(fmt.Sprintf("作用域: %s\n", scopeTag))
		sb.WriteString(fmt.Sprintf("创建时间: %s\n", plan.CreatedAt.Format("2006-01-02 15:04:05")))
		sb.WriteString(fmt.Sprintf("更新时间: %s\n", plan.UpdatedAt.Format("2006-01-02 15:04:05")))
		sb.WriteString(fmt.Sprintf("\n描述:\n%s\n", plan.Description))
		sb.WriteString(fmt.Sprintf("\n内容:\n%s", plan.Content))

		// 如果有子任务，也显示出来
		if len(plan.SubTasks) > 0 {
			sb.WriteString("\n\n子任务:\n")
			for _, st := range plan.SubTasks {
				stStatus := getPlanStatusText(st.Status)
				sb.WriteString(fmt.Sprintf("  - [%d] %s (%s, %d%%)\n", st.ID, st.Title, stStatus, st.Progress))
			}
		}

		return NewTextResult(sb.String()), nil, nil
	})

	// plan_update - 更新计划
	mcp.AddTool(server, &mcp.Tool{
		Name:        "plan_update",
		Description: `更新计划，只更新提供的字段（title/description/content/progress）；至少提供一个字段，否则返回错误。progress: 0=待开始，1-99=进行中，100=已完成（状态自动调整）。`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, input PlanUpdateInput) (*mcp.CallToolResult, any, error) {
		// 检查是否有更新（至少一个字段）
		if input.Title == nil && input.Description == nil &&
			input.Content == nil && input.Progress == nil {
			return NewErrorResult("至少提供一个要更新的字段"), nil, nil
		}

		// 构建更新 DTO
		updateDTO := &dto.PlanUpdateDTO{
			Code:        input.Code,
			Title:       input.Title,
			Description: input.Description,
			Content:     input.Content,
			Progress:    input.Progress,
		}

		// 执行更新
		if err := bs.PlanService.UpdatePlan(ctx, updateDTO); err != nil {
			return NewErrorResult(err.Error()), nil, nil
		}

		// 构建响应消息
		var parts []string
		if input.Title != nil {
			parts = append(parts, "标题")
		}
		if input.Description != nil {
			parts = append(parts, "描述")
		}
		if input.Content != nil {
			parts = append(parts, "内容")
		}
		if input.Progress != nil {
			parts = append(parts, fmt.Sprintf("进度(%d%%)", *input.Progress))
		}

		return NewTextResult(fmt.Sprintf("计划 %s 更新成功: %s", input.Code, strings.Join(parts, "、"))), nil, nil
	})
}

// getPlanStatusText 获取计划状态文本
func getPlanStatusText(status entity.PlanStatus) string {
	switch status {
	case entity.PlanStatusPending:
		return "待开始"
	case entity.PlanStatusInProgress:
		return "进行中"
	case entity.PlanStatusCompleted:
		return "已完成"
	case entity.PlanStatusCancelled:
		return "已取消"
	default:
		return "未知"
	}
}
