package tools

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/XiaoLFeng/llm-memory/pkg/types"
	"github.com/XiaoLFeng/llm-memory/startup"
)

// PlanListInput plan_list 工具输入
type PlanListInput struct{}

// PlanCreateInput plan_create 工具输入
type PlanCreateInput struct {
	Title       string `json:"title" jsonschema:"description=计划的标题，简洁描述计划目标，例如：'完成用户认证模块开发'、'学习Go语言基础'"`
	Description string `json:"description,omitempty" jsonschema:"description=计划的详细描述，包含具体步骤、目标、注意事项等，例如：'实现JWT认证，包括登录、注册、token刷新等功能'"`
}

// PlanUpdateProgressInput plan_update_progress 工具输入
type PlanUpdateProgressInput struct {
	ID       int `json:"id" jsonschema:"description=要更新的计划ID，可通过 plan_list 获取"`
	Progress int `json:"progress" jsonschema:"description=计划的完成进度，范围0-100的整数。0表示未开始，100表示已完成。系统会根据进度自动更新计划状态"`
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
- 已取消：计划被取消`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, input PlanListInput) (*mcp.CallToolResult, any, error) {
		plans, err := bs.PlanService.ListPlans(ctx)
		if err != nil {
			return NewErrorResult(err.Error()), nil, nil
		}
		if len(plans) == 0 {
			return NewTextResult("暂无计划"), nil, nil
		}
		result := "计划列表:\n"
		for _, p := range plans {
			status := getStatusText(p.Status)
			result += fmt.Sprintf("- [%d] %s (%s, 进度: %d%%)\n", p.ID, p.Title, status, p.Progress)
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
- 描述中包含具体的里程碑或检查点
- 创建后可通过 plan_update_progress 更新进度

示例：
- 标题："重构用户模块"，描述："1. 分析现有代码 2. 设计新架构 3. 实现核心功能 4. 编写测试 5. 代码审查"`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, input PlanCreateInput) (*mcp.CallToolResult, any, error) {
		plan, err := bs.PlanService.CreatePlan(ctx, input.Title, input.Description)
		if err != nil {
			return NewErrorResult(err.Error()), nil, nil
		}
		return NewTextResult(fmt.Sprintf("计划创建成功! ID: %d, 标题: %s", plan.ID, plan.Title)), nil, nil
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

// getStatusText 获取计划状态文本
func getStatusText(status types.PlanStatus) string {
	switch status {
	case types.PlanStatusPending:
		return "待开始"
	case types.PlanStatusInProgress:
		return "进行中"
	case types.PlanStatusCompleted:
		return "已完成"
	case types.PlanStatusCancelled:
		return "已取消"
	default:
		return "未知"
	}
}
