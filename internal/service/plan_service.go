package service

import (
	"context"
	"errors"
	"time"

	"github.com/XiaoLFeng/llm-memory/internal/database"
	"github.com/XiaoLFeng/llm-memory/pkg/types"
)

// PlanService 计划服务层结构体
// 嘿嘿~ 这是计划管理的业务逻辑核心呢！📋✨
type PlanService struct {
	repo database.PlanRepository
}

// NewPlanService 创建新的计划服务实例
// 构造函数模式，优雅地初始化服务~ 💖
func NewPlanService(repo database.PlanRepository) *PlanService {
	return &PlanService{
		repo: repo,
	}
}

// CreatePlan 创建新计划
// 业务逻辑：验证参数并创建计划~ (´∀｀)
func (s *PlanService) CreatePlan(ctx context.Context, title, description string) (*types.Plan, error) {
	// 参数验证 - 标题不能为空哦！
	if title == "" {
		return nil, errors.New("计划标题不能为空")
	}

	// 使用 types 包的构造函数创建计划
	plan := types.NewPlan(title, description)

	// 保存到数据库~ ✨
	if err := s.repo.Create(ctx, plan); err != nil {
		return nil, err
	}

	return plan, nil
}

// UpdatePlan 更新计划
// 业务逻辑：验证并更新计划信息~ 🎮
func (s *PlanService) UpdatePlan(ctx context.Context, plan *types.Plan) error {
	// 参数验证
	if plan == nil {
		return errors.New("计划对象不能为空")
	}
	if plan.ID == 0 {
		return errors.New("计划ID不能为0")
	}
	if plan.Title == "" {
		return errors.New("计划标题不能为空")
	}

	// 验证计划是否存在
	existingPlan, err := s.repo.FindByID(ctx, plan.ID)
	if err != nil {
		return errors.New("计划不存在")
	}
	if existingPlan == nil {
		return errors.New("计划不存在")
	}

	// 更新时间戳
	plan.UpdatedAt = time.Now()

	// 执行更新操作~ ＼(^o^)／
	return s.repo.Update(ctx, plan)
}

// DeletePlan 删除计划
// 业务逻辑：验证并删除指定计划~ 🍫
func (s *PlanService) DeletePlan(ctx context.Context, id int) error {
	// 参数验证
	if id <= 0 {
		return errors.New("无效的计划ID")
	}

	// 验证计划是否存在
	existingPlan, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New("计划不存在")
	}
	if existingPlan == nil {
		return errors.New("计划不存在")
	}

	// 执行删除操作~ (´∀｀)
	return s.repo.Delete(ctx, id)
}

// GetPlan 获取单个计划
// 业务逻辑：根据ID查询计划~ 🎯
func (s *PlanService) GetPlan(ctx context.Context, id int) (*types.Plan, error) {
	// 参数验证
	if id <= 0 {
		return nil, errors.New("无效的计划ID")
	}

	// 查询计划
	plan, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if plan == nil {
		return nil, errors.New("计划不存在")
	}

	return plan, nil
}

// ListPlans 获取所有计划列表
// 业务逻辑：查询全部计划~ ＼(^o^)／
func (s *PlanService) ListPlans(ctx context.Context) ([]types.Plan, error) {
	plans, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	// 如果没有计划，返回空切片而不是nil
	if plans == nil {
		return make([]types.Plan, 0), nil
	}

	return plans, nil
}

// ListByStatus 根据状态获取计划列表
// 业务逻辑：按状态筛选计划~ 💖
func (s *PlanService) ListByStatus(ctx context.Context, status types.PlanStatus) ([]types.Plan, error) {
	// 验证状态值是否有效
	if !isValidStatus(status) {
		return nil, errors.New("无效的计划状态")
	}

	plans, err := s.repo.FindByStatus(ctx, status)
	if err != nil {
		return nil, err
	}

	// 如果没有计划，返回空切片而不是nil
	if plans == nil {
		return make([]types.Plan, 0), nil
	}

	return plans, nil
}

// StartPlan 开始计划
// 业务逻辑：将计划状态改为进行中~ ✨
func (s *PlanService) StartPlan(ctx context.Context, id int) error {
	// 参数验证
	if id <= 0 {
		return errors.New("无效的计划ID")
	}

	// 获取计划
	plan, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if plan == nil {
		return errors.New("计划不存在")
	}

	// 验证状态转换是否合法
	if plan.Status == types.PlanStatusCompleted {
		return errors.New("已完成的计划无法重新开始")
	}
	if plan.Status == types.PlanStatusCancelled {
		return errors.New("已取消的计划无法开始")
	}

	// 更新状态和开始时间
	plan.Status = types.PlanStatusInProgress
	now := time.Now()
	plan.StartDate = &now
	plan.UpdatedAt = now

	// 如果进度为0，更新为1表示已开始
	if plan.Progress == 0 {
		plan.Progress = 1
	}

	// 保存更新~ 🎮
	return s.repo.Update(ctx, plan)
}

// CompletePlan 完成计划
// 业务逻辑：将计划状态改为已完成~ (´∀｀)
func (s *PlanService) CompletePlan(ctx context.Context, id int) error {
	// 参数验证
	if id <= 0 {
		return errors.New("无效的计划ID")
	}

	// 获取计划
	plan, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if plan == nil {
		return errors.New("计划不存在")
	}

	// 验证状态转换是否合法
	if plan.Status == types.PlanStatusCancelled {
		return errors.New("已取消的计划无法标记为完成")
	}

	// 更新状态、进度和结束时间
	plan.Status = types.PlanStatusCompleted
	plan.Progress = 100
	now := time.Now()
	plan.EndDate = &now
	plan.UpdatedAt = now

	// 保存更新~ 🍫
	return s.repo.Update(ctx, plan)
}

// UpdateProgress 更新计划进度
// 业务逻辑：更新进度并自动调整状态~ ＼(^o^)／
func (s *PlanService) UpdateProgress(ctx context.Context, id int, progress int) error {
	// 参数验证
	if id <= 0 {
		return errors.New("无效的计划ID")
	}
	if progress < 0 || progress > 100 {
		return errors.New("进度值必须在0-100之间")
	}

	// 获取计划
	plan, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if plan == nil {
		return errors.New("计划不存在")
	}

	// 验证状态 - 已取消的计划不能更新进度
	if plan.Status == types.PlanStatusCancelled {
		return errors.New("已取消的计划无法更新进度")
	}

	// 使用 Plan 类型的 UpdateProgress 方法
	// 这个方法会自动根据进度更新状态~ 💖
	plan.UpdateProgress(progress)

	// 如果进度达到100%，设置结束时间
	if progress == 100 {
		now := time.Now()
		plan.EndDate = &now
	}

	// 保存更新~ ✨
	return s.repo.Update(ctx, plan)
}

// isValidStatus 验证计划状态是否有效
// 辅助函数：检查状态值是否在允许的范围内~ 🎯
func isValidStatus(status types.PlanStatus) bool {
	validStatuses := []types.PlanStatus{
		types.PlanStatusPending,
		types.PlanStatusInProgress,
		types.PlanStatusCompleted,
		types.PlanStatusCancelled,
	}

	for _, validStatus := range validStatuses {
		if status == validStatus {
			return true
		}
	}

	return false
}
