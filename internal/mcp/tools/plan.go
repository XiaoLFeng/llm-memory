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
	ID       uint `json:"id" jsonschema:"要更新的计划ID"`
	Progress int  `json:"progress" jsonschema:"完成进度(0-100)，系统会自动更新状态"`
}

// RegisterPlanTools 注册计划管理工具
// 呀~ 计划相关的 MCP 工具都在这里！✨
func RegisterPlanTools(server *mcp.Server, bs *startup.Bootstrap) {
	// plan_list - 列出所有计划
	mcp.AddTool(server, &mcp.Tool{
		Name: "plan_list",
		Description: `列出用户的所有计划，包含计划状态和进度信息。

使用场景：
- 查看当前所有计划的整体情况
- 了解各计划的进度状态
- 获取计划ID用于更新进度

返回信息：计划ID、标题、状态（待开始/进行中/已完成/已取消）、进度百分比

计划状态说明：
- 待开始：进度为0，尚未启动
- 进行中：进度在1-99之间
- 已完成：进度达到100
- 已取消：计划被取消

作用域说明：
- personal: 只显示当前目录的计划
- group: 只显示当前组的计划
- global: 只显示全局计划
- all: 显示所有可见计划（默认）`,
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
			scopeTag := getScopeTag(p.GroupID, p.Path)
			result += fmt.Sprintf("- [%d] %s (%s, 进度: %d%%) %s\n", p.ID, p.Title, status, p.Progress, scopeTag)
		}
		return NewTextResult(result), nil, nil
	})

	// plan_create - 创建新计划
	mcp.AddTool(server, &mcp.Tool{
		Name: "plan_create",
		Description: `创建一个新的计划，用于跟踪长期目标或复杂任务的完成情况。

使用场景：
- 用户提出需要分阶段完成的任务
- 跟踪学习计划、项目开发进度
- 管理需要持续关注的目标

计划 vs 待办事项：
- 计划：适合长期、复杂、需要跟踪进度的任务
- 待办事项：适合短期、简单、一次性完成的任务

最佳实践：
- 标题应明确表达计划目标
- description 用于摘要展示
- content 用于详细内容（支持 Markdown）
- 创建后可通过 plan_update_progress 更新进度

示例：
- 标题："重构用户模块"
- 描述："分析、设计、实现、测试、审查"
- 内容："## 阶段1: 分析\n- 分析现有代码\n- 识别问题点\n\n## 阶段2: 设计\n..."

作用域说明：
- personal: 保存到当前目录（只在此目录可见）
- group: 保存到当前组（组内所有路径可见）
- global: 保存为全局（任何地方可见，默认）`,
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
		scopeTag := getScopeTag(plan.GroupID, plan.Path)
		return NewTextResult(fmt.Sprintf("计划创建成功! ID: %d, 标题: %s %s", plan.ID, plan.Title, scopeTag)), nil, nil
	})

	// plan_update_progress - 更新计划进度
	mcp.AddTool(server, &mcp.Tool{
		Name: "plan_update_progress",
		Description: `更新指定计划的完成进度。

使用场景：
- 计划有新进展时更新进度
- 用户汇报任务完成情况
- 里程碑达成时记录进度

进度值说明：
- 0：未开始
- 1-99：进行中（会自动将状态设为"进行中"）
- 100：已完成（会自动将状态设为"已完成"）

注意事项：
- 进度只能是0-100的整数
- 已取消的计划无法更新进度
- 进度更新后状态会自动调整

建议：根据计划描述中的里程碑来估算进度，例如5个步骤完成3个约为60%`,
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
