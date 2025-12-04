package service

import (
	"context"
	"errors"
	"strings"

	"github.com/XiaoLFeng/llm-memory/internal/models"
	"github.com/XiaoLFeng/llm-memory/internal/models/dto"
	"github.com/XiaoLFeng/llm-memory/internal/models/entity"
	"github.com/XiaoLFeng/llm-memory/pkg/types"
)

// PlanService 计划服务层结构体
type PlanService struct {
	model *models.PlanModel
}

// NewPlanService 创建新的计划服务实例
func NewPlanService(model *models.PlanModel) *PlanService {
	return &PlanService{
		model: model,
	}
}

// CreatePlan 创建新计划
// 纯关联模式：数据存储时只使用 PathID
func (s *PlanService) CreatePlan(ctx context.Context, input *dto.PlanCreateDTO, scopeCtx *types.ScopeContext) (*entity.Plan, error) {
	// 参数验证 - 标题不能为空
	if strings.TrimSpace(input.Title) == "" {
		return nil, errors.New("计划标题不能为空")
	}
	// 参数验证 - 描述不能为空
	if strings.TrimSpace(input.Description) == "" {
		return nil, errors.New("计划描述不能为空")
	}
	// 参数验证 - 内容不能为空
	if strings.TrimSpace(input.Content) == "" {
		return nil, errors.New("计划内容不能为空")
	}

	// 解析作用域 -> PathID（global=true 则存储到全局，否则使用当前路径）
	pathID := int64(0)
	if !input.Global {
		pathID = resolveDefaultPathID(scopeCtx)
		if pathID == 0 {
			return nil, errors.New("无法确定私有/小组作用域，请先初始化 paths 或选择全局")
		}
	}

	// 创建计划实例
	plan := &entity.Plan{
		Global:      input.Global,
		PathID:      pathID,
		Title:       strings.TrimSpace(input.Title),
		Description: strings.TrimSpace(input.Description),
		Content:     strings.TrimSpace(input.Content),
		Status:      entity.PlanStatusPending,
		Progress:    0,
	}

	// 保存到数据库
	if err := s.model.Create(ctx, plan); err != nil {
		return nil, err
	}

	return plan, nil
}

// UpdatePlan 更新计划
func (s *PlanService) UpdatePlan(ctx context.Context, input *dto.PlanUpdateDTO) error {
	// 参数验证
	if input.ID == 0 {
		return errors.New("计划ID不能为0")
	}

	// 获取现有计划
	plan, err := s.model.FindByID(ctx, input.ID)
	if err != nil {
		return errors.New("计划不存在")
	}

	// 验证状态 - 已取消的计划不能更新
	if plan.Status == entity.PlanStatusCancelled {
		return errors.New("已取消的计划无法更新")
	}

	// 应用更新
	if input.Title != nil {
		title := strings.TrimSpace(*input.Title)
		if title == "" {
			return errors.New("计划标题不能为空")
		}
		plan.Title = title
	}
	if input.Description != nil {
		desc := strings.TrimSpace(*input.Description)
		if desc == "" {
			return errors.New("计划描述不能为空")
		}
		plan.Description = desc
	}
	if input.Content != nil {
		content := strings.TrimSpace(*input.Content)
		if content == "" {
			return errors.New("计划内容不能为空")
		}
		plan.Content = content
	}
	if input.Progress != nil {
		progress := *input.Progress
		if progress < 0 || progress > 100 {
			return errors.New("进度值必须在0-100之间")
		}
		plan.UpdateProgress(progress)
	}

	// 执行更新操作
	return s.model.Update(ctx, plan)
}

// DeletePlan 删除计划
func (s *PlanService) DeletePlan(ctx context.Context, id int64) error {
	// 参数验证
	if id == 0 {
		return errors.New("无效的计划ID")
	}

	// 验证计划是否存在
	_, err := s.model.FindByID(ctx, id)
	if err != nil {
		return errors.New("计划不存在")
	}

	// 执行删除操作
	return s.model.Delete(ctx, id)
}

// GetPlan 获取单个计划
func (s *PlanService) GetPlan(ctx context.Context, id int64) (*entity.Plan, error) {
	// 参数验证
	if id == 0 {
		return nil, errors.New("无效的计划ID")
	}

	// 查询计划
	plan, err := s.model.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if plan == nil {
		return nil, errors.New("计划不存在")
	}

	return plan, nil
}

// ListPlans 获取所有计划列表
func (s *PlanService) ListPlans(ctx context.Context) ([]entity.Plan, error) {
	plans, err := s.model.FindByFilter(ctx, models.DefaultVisibilityFilter())
	if err != nil {
		return nil, err
	}

	// 如果没有计划，返回空切片而不是nil
	if plans == nil {
		return make([]entity.Plan, 0), nil
	}

	return plans, nil
}

// ListPlansByScope 根据作用域列出计划
// 纯关联模式：使用 PathID 和 GroupPathIDs 进行查询
func (s *PlanService) ListPlansByScope(ctx context.Context, scope string, scopeCtx *types.ScopeContext) ([]entity.Plan, error) {
	filter := buildVisibilityFilter(scope, scopeCtx)
	plans, err := s.model.FindByFilter(ctx, filter)
	if err != nil {
		return nil, err
	}

	if plans == nil {
		return make([]entity.Plan, 0), nil
	}

	return plans, nil
}

// ListByStatus 根据状态获取计划列表
func (s *PlanService) ListByStatus(ctx context.Context, status entity.PlanStatus) ([]entity.Plan, error) {
	// 验证状态值是否有效
	if !isValidPlanStatus(status) {
		return nil, errors.New("无效的计划状态")
	}

	plans, err := s.model.FindByStatus(ctx, status)
	if err != nil {
		return nil, err
	}

	// 如果没有计划，返回空切片而不是nil
	if plans == nil {
		return make([]entity.Plan, 0), nil
	}

	return plans, nil
}

// StartPlan 开始计划
func (s *PlanService) StartPlan(ctx context.Context, id int64) error {
	// 参数验证
	if id == 0 {
		return errors.New("无效的计划ID")
	}

	// 获取计划
	plan, err := s.model.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// 验证状态转换是否合法
	if plan.Status == entity.PlanStatusCompleted {
		return errors.New("已完成的计划无法重新开始")
	}
	if plan.Status == entity.PlanStatusCancelled {
		return errors.New("已取消的计划无法开始")
	}

	// 执行开始
	plan.Start()

	// 保存更新
	return s.model.Update(ctx, plan)
}

// CompletePlan 完成计划
func (s *PlanService) CompletePlan(ctx context.Context, id int64) error {
	// 参数验证
	if id == 0 {
		return errors.New("无效的计划ID")
	}

	// 获取计划
	plan, err := s.model.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// 验证状态转换是否合法
	if plan.Status == entity.PlanStatusCancelled {
		return errors.New("已取消的计划无法标记为完成")
	}

	// 执行完成
	plan.Complete()

	// 保存更新
	return s.model.Update(ctx, plan)
}

// CancelPlan 取消计划
func (s *PlanService) CancelPlan(ctx context.Context, id int64) error {
	if id == 0 {
		return errors.New("无效的计划ID")
	}

	plan, err := s.model.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if plan.Status == entity.PlanStatusCompleted {
		return errors.New("已完成的计划无法取消")
	}

	plan.Cancel()
	return s.model.Update(ctx, plan)
}

// UpdateProgress 更新计划进度
func (s *PlanService) UpdateProgress(ctx context.Context, id int64, progress int) error {
	// 参数验证
	if id == 0 {
		return errors.New("无效的计划ID")
	}
	if progress < 0 || progress > 100 {
		return errors.New("进度值必须在0-100之间")
	}

	// 获取计划
	plan, err := s.model.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// 验证状态 - 已取消的计划不能更新进度
	if plan.Status == entity.PlanStatusCancelled {
		return errors.New("已取消的计划无法更新进度")
	}

	// 使用 Plan 类型的 UpdateProgress 方法
	plan.UpdateProgress(progress)

	// 保存更新
	return s.model.Update(ctx, plan)
}

// AddSubTask 添加子任务
func (s *PlanService) AddSubTask(ctx context.Context, planID int64, title, description string) (*entity.SubTask, error) {
	if planID == 0 {
		return nil, errors.New("无效的计划ID")
	}
	if strings.TrimSpace(title) == "" {
		return nil, errors.New("子任务标题不能为空")
	}

	// 验证计划存在
	_, err := s.model.FindByID(ctx, planID)
	if err != nil {
		return nil, errors.New("计划不存在")
	}

	return s.model.AddSubTask(ctx, planID, strings.TrimSpace(title), strings.TrimSpace(description))
}

// UpdateSubTask 更新子任务
func (s *PlanService) UpdateSubTask(ctx context.Context, input *dto.SubTaskUpdateDTO) error {
	if input.ID == 0 {
		return errors.New("无效的子任务ID")
	}

	subTask, err := s.model.GetSubTask(ctx, input.ID)
	if err != nil {
		return errors.New("子任务不存在")
	}

	if input.Title != nil {
		title := strings.TrimSpace(*input.Title)
		if title == "" {
			return errors.New("子任务标题不能为空")
		}
		subTask.Title = title
	}
	if input.Description != nil {
		subTask.Description = strings.TrimSpace(*input.Description)
	}
	if input.Status != nil {
		subTask.Status = entity.PlanStatus(*input.Status)
	}
	if input.Progress != nil {
		progress := *input.Progress
		if progress < 0 || progress > 100 {
			return errors.New("进度值必须在0-100之间")
		}
		subTask.Progress = progress
	}

	return s.model.UpdateSubTask(ctx, subTask)
}

// DeleteSubTask 删除子任务
func (s *PlanService) DeleteSubTask(ctx context.Context, subTaskID int64) error {
	if subTaskID == 0 {
		return errors.New("无效的子任务ID")
	}
	return s.model.DeleteSubTask(ctx, subTaskID)
}

// isValidPlanStatus 验证计划状态是否有效
func isValidPlanStatus(status entity.PlanStatus) bool {
	validStatuses := []entity.PlanStatus{
		entity.PlanStatusPending,
		entity.PlanStatusInProgress,
		entity.PlanStatusCompleted,
		entity.PlanStatusCancelled,
	}

	for _, validStatus := range validStatuses {
		if status == validStatus {
			return true
		}
	}

	return false
}

// ToPlanResponseDTO 将 Plan entity 转换为 ResponseDTO
// 纯关联模式：使用 PathID 判断作用域
func ToPlanResponseDTO(plan *entity.Plan, scopeCtx *types.ScopeContext) *dto.PlanResponseDTO {
	if plan == nil {
		return nil
	}

	scope := types.GetScopeForDisplayWithGlobal(plan.Global, plan.PathID, scopeCtx)

	// 转换子任务
	subTasks := make([]dto.SubTaskDTO, 0, len(plan.SubTasks))
	for _, st := range plan.SubTasks {
		subTasks = append(subTasks, dto.SubTaskDTO{
			ID:          st.ID,
			Title:       st.Title,
			Description: st.Description,
			Status:      string(st.Status),
			Progress:    st.Progress,
			SortOrder:   st.SortOrder,
			CreatedAt:   st.CreatedAt,
			UpdatedAt:   st.UpdatedAt,
		})
	}

	return &dto.PlanResponseDTO{
		ID:          plan.ID,
		Title:       plan.Title,
		Description: plan.Description,
		Content:     plan.Content,
		Status:      string(plan.Status),
		Progress:    plan.Progress,
		SubTasks:    subTasks,
		Scope:       string(scope),
		CreatedAt:   plan.CreatedAt,
		UpdatedAt:   plan.UpdatedAt,
	}
}

// ToPlanListDTO 将 Plan entity 转换为 ListDTO
func ToPlanListDTO(plan *entity.Plan) *dto.PlanListDTO {
	if plan == nil {
		return nil
	}

	return &dto.PlanListDTO{
		ID:          plan.ID,
		Title:       plan.Title,
		Description: plan.Description,
		Status:      string(plan.Status),
		Progress:    plan.Progress,
	}
}
