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

// PlanHandler 计划命令处理器
type PlanHandler struct {
	bs *startup.Bootstrap
}

// NewPlanHandler 创建计划处理器
func NewPlanHandler(bs *startup.Bootstrap) *PlanHandler {
	return &PlanHandler{bs: bs}
}

// List 列出所有计划
func (h *PlanHandler) List(ctx context.Context) error {
	// 使用 ListPlansByScope 确保权限隔离：全局 + 当前路径相关
	plans, err := h.bs.PlanService.ListPlansByScope(ctx, "all", h.bs.CurrentScope)
	if err != nil {
		return err
	}

	if len(plans) == 0 {
		cli.PrintInfo("暂无计划~ 快创建一个吧！")
		return nil
	}

	cli.PrintTitle(cli.IconPlan + " 计划列表")
	table := output.NewTable("标识码", "标题", "状态", "进度")
	for _, p := range plans {
		table.AddRow(
			p.Code,
			p.Title,
			getPlanStatusText(p.Status),
			fmt.Sprintf("%d%%", p.Progress),
		)
	}
	table.Print()

	return nil
}

// Create 创建计划
func (h *PlanHandler) Create(ctx context.Context, code, title, description string, global bool) error {
	createDTO := &dto.PlanCreateDTO{
		Code:        code,
		Title:       title,
		Description: description,
		Global:      global,
	}
	plan, err := h.bs.PlanService.CreatePlan(ctx, createDTO, h.bs.CurrentScope)
	if err != nil {
		return err
	}

	cli.PrintSuccess(fmt.Sprintf("计划创建成功！标识码: %s, 标题: %s", plan.Code, plan.Title))
	return nil
}

// UpdateProgress 更新计划进度
func (h *PlanHandler) UpdateProgress(ctx context.Context, code string, progress int) error {
	if err := h.bs.PlanService.UpdateProgress(ctx, code, progress); err != nil {
		return err
	}

	cli.PrintSuccess(fmt.Sprintf("计划 %s 进度已更新为 %d%%", code, progress))
	return nil
}

// Start 开始计划
func (h *PlanHandler) Start(ctx context.Context, code string) error {
	if err := h.bs.PlanService.StartPlan(ctx, code); err != nil {
		return err
	}

	cli.PrintSuccess(fmt.Sprintf("计划 %s 已开始", code))
	return nil
}

// Complete 完成计划
func (h *PlanHandler) Complete(ctx context.Context, code string) error {
	if err := h.bs.PlanService.CompletePlan(ctx, code); err != nil {
		return err
	}

	cli.PrintSuccess(fmt.Sprintf("计划 %s 已完成", code))
	return nil
}

// Delete 删除计划
func (h *PlanHandler) Delete(ctx context.Context, code string) error {
	if err := h.bs.PlanService.DeletePlan(ctx, code); err != nil {
		return err
	}

	cli.PrintSuccess(fmt.Sprintf("计划 %s 已删除", code))
	return nil
}

// Get 获取计划详情
func (h *PlanHandler) Get(ctx context.Context, code string) error {
	plan, err := h.bs.PlanService.GetPlan(ctx, code)
	if err != nil {
		return err
	}

	cli.PrintTitle(cli.IconClipboard + " 计划详情")
	fmt.Printf("标识码:   %s\n", plan.Code)
	fmt.Printf("标题:     %s\n", plan.Title)
	fmt.Printf("状态:     %s\n", getPlanStatusText(plan.Status))
	fmt.Printf("进度:     %d%%\n", plan.Progress)
	fmt.Printf("创建时间: %s\n", plan.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("更新时间: %s\n", plan.UpdatedAt.Format("2006-01-02 15:04:05"))
	if plan.Description != "" {
		fmt.Println("\n描述:")
		fmt.Println(plan.Description)
	}
	if plan.Content != "" {
		fmt.Println("\n内容:")
		fmt.Println(plan.Content)
	}

	return nil
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
