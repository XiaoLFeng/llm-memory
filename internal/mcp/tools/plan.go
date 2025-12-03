package tools

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/XiaoLFeng/llm-memory/internal/models/dto"
	"github.com/XiaoLFeng/llm-memory/internal/models/entity"
	"github.com/XiaoLFeng/llm-memory/startup"
)

// PlanListInput plan_list 工具输入
type PlanListInput struct {
	Scope string `json:"scope,omitempty" jsonschema:"作用域过滤(personal/group/global/all)，默认all显示全部"`
}

// PlanCreateInput plan_create 工具输入
type PlanCreateInput struct {
	Title       string `json:"title" jsonschema:"计划标题，简洁描述计划目标"`
	Description string `json:"description,omitempty" jsonschema:"计划的详细描述，包含具体步骤和目标"`
	Content     string `json:"content,omitempty" jsonschema:"计划的详细内容（新增），支持 Markdown 格式"`
	Scope       string `json:"scope,omitempty" jsonschema:"保存到哪个作用域(personal/group/global)，默认global"`
}

// PlanUpdateProgressInput plan_update_progress 工具输入
type PlanUpdateProgressInput struct {
	ID       int64 `json:"id" jsonschema:"要更新的计划ID"`
	Progress int   `json:"progress" jsonschema:"完成进度(0-100)，系统会自动更新状态"`
}

// RegisterPlanTools 注册计划管理工具
func RegisterPlanTools(server *mcp.Server, bs *startup.Bootstrap) {
	// plan_list - 列出所有计划
	mcp.AddTool(server, &mcp.Tool{
		Name:        "plan_list",
		Description: `列出所有计划及进度状态。scope: personal/group/global/all，默认all=使用当前作用域集合。未指定也会落在 currentScope（无则 global）。`,
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
			scopeTag := getScopeTagFromPathID(p.PathID)
			result += fmt.Sprintf("- [%d] %s (%s, 进度: %d%%) %s\n", p.ID, p.Title, status, p.Progress, scopeTag)
		}
		return NewTextResult(result), nil, nil
	})

	// plan_create - 创建新计划
	mcp.AddTool(server, &mcp.Tool{
		Name:        "plan_create",
		Description: `创建计划，用于“需要跟踪进度的多步骤目标”。必填: title。可选: description、content(Markdown)、scope。短平快的单一步行动请用 todo_create；长期事实/偏好请用 memory_create。未指定 scope 默认使用 currentScope（缺省则 global）。`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, input PlanCreateInput) (*mcp.CallToolResult, any, error) {
		// 构建创建 DTO
		createDTO := &dto.PlanCreateDTO{
			Title:       input.Title,
			Description: input.Description,
			Content:     input.Content,
			Scope:       input.Scope,
		}

		// 构建作用域上下文
		scopeCtx := buildScopeContext(input.Scope, bs)

		plan, err := bs.PlanService.CreatePlan(ctx, createDTO, scopeCtx)
		if err != nil {
			return NewErrorResult(err.Error()), nil, nil
		}
		scopeTag := getScopeTagFromPathID(plan.PathID)
		return NewTextResult(fmt.Sprintf("计划创建成功! ID: %d, 标题: %s %s", plan.ID, plan.Title, scopeTag)), nil, nil
	})

	// plan_update_progress - 更新计划进度
	mcp.AddTool(server, &mcp.Tool{
		Name:        "plan_update_progress",
		Description: `更新计划进度(0-100)，状态自动调整：0=待开始，1-99=进行中，100=已完成。`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, input PlanUpdateProgressInput) (*mcp.CallToolResult, any, error) {
		if err := bs.PlanService.UpdateProgress(ctx, input.ID, input.Progress); err != nil {
			return NewErrorResult(err.Error()), nil, nil
		}
		return NewTextResult(fmt.Sprintf("计划 %d 进度已更新为 %d%%", input.ID, input.Progress)), nil, nil
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
