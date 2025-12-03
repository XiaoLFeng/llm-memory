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
	plans, err := h.bs.PlanService.ListPlans(ctx)
	if err != nil {
		return err
	}

	if len(plans) == 0 {
		cli.PrintInfo("暂无计划~ 快创建一个吧！")
		return nil
	}

	cli.PrintTitle(cli.IconPlan + " 计划列表")
	table := output.NewTable("ID", "标题", "状态", "进度")
	for _, p := range plans {
		table.AddRow(
			fmt.Sprintf("%d", p.ID),
			p.Title,
			getPlanStatusText(p.Status),
			fmt.Sprintf("%d%%", p.Progress),
		)
	}
	table.Print()

	return nil
}

// Create 创建计划
func (h *PlanHandler) Create(ctx context.Context, title, description string) error {
	createDTO := &dto.PlanCreateDTO{
		Title:       title,
		Description: description,
		Scope:       "global",
	}
	plan, err := h.bs.PlanService.CreatePlan(ctx, createDTO, nil)
	if err != nil {
		return err
	}

	cli.PrintSuccess(fmt.Sprintf("计划创建成功！ID: %d, 标题: %s", plan.ID, plan.Title))
	return nil
}

// UpdateProgress 更新计划进度
func (h *PlanHandler) UpdateProgress(ctx context.Context, id int64, progress int) error {
	if err := h.bs.PlanService.UpdateProgress(ctx, id, progress); err != nil {
		return err
	}

	cli.PrintSuccess(fmt.Sprintf("计划 %d 进度已更新为 %d%%", id, progress))
	return nil
}

// Start 开始计划
func (h *PlanHandler) Start(ctx context.Context, id int64) error {
	if err := h.bs.PlanService.StartPlan(ctx, id); err != nil {
		return err
	}

	cli.PrintSuccess(fmt.Sprintf("计划 %d 已开始", id))
	return nil
}

// Complete 完成计划
func (h *PlanHandler) Complete(ctx context.Context, id int64) error {
	if err := h.bs.PlanService.CompletePlan(ctx, id); err != nil {
		return err
	}

	cli.PrintSuccess(fmt.Sprintf("计划 %d 已完成", id))
	return nil
}

// Delete 删除计划
func (h *PlanHandler) Delete(ctx context.Context, id int64) error {
	if err := h.bs.PlanService.DeletePlan(ctx, id); err != nil {
		return err
	}

	cli.PrintSuccess(fmt.Sprintf("计划 %d 已删除", id))
	return nil
}

// Get 获取计划详情
func (h *PlanHandler) Get(ctx context.Context, id int64) error {
	plan, err := h.bs.PlanService.GetPlan(ctx, id)
	if err != nil {
		return err
	}

	cli.PrintTitle(cli.IconClipboard + " 计划详情")
	fmt.Printf("ID:       %d\n", plan.ID)
	fmt.Printf("标题:     %s\n", plan.Title)
	fmt.Printf("状态:     %s\n", getPlanStatusText(plan.Status))
	fmt.Printf("进度:     %d%%\n", plan.Progress)
	if plan.StartDate != nil {
		fmt.Printf("开始时间: %s\n", plan.StartDate.Format("2006-01-02 15:04:05"))
	}
	if plan.EndDate != nil {
		fmt.Printf("结束时间: %s\n", plan.EndDate.Format("2006-01-02 15:04:05"))
	}
	fmt.Printf("创建时间: %s\n", plan.CreatedAt.Format("2006-01-02 15:04:05"))
	if plan.Description != "" {
		fmt.Println("\n描述:")
		fmt.Println(plan.Description)
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
