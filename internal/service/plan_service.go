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
	planModel *models.PlanModel
}

// NewPlanService 创建新的计划服务实例
func NewPlanService(model *models.PlanModel) *PlanService {
	return &PlanService{
		planModel: model,
	}
}

// CreatePlan 创建新计划
// PathID 关联个人或小组路径
func (s *PlanService) CreatePlan(ctx context.Context, input *dto.PlanCreateDTO, scopeCtx *types.ScopeContext) (*entity.Plan, error) {
	// 参数验证 - code 格式验证
	if err := entity.ValidateCode(input.Code); err != nil {
		return nil, err
	}
	// 参数验证 - code 唯一性验证（活跃状态中）
	exists, err := s.planModel.ExistsActiveCode(ctx, input.Code, 0)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("活跃状态中已存在相同的 code，请使用不同的标识码")
	}
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

	// 解析作用域 -> PathID（必须指定路径）
	pathID := resolveDefaultPathID(scopeCtx)
	if pathID == 0 {
		return nil, errors.New("无法确定作用域，请先初始化 paths")
	}

	// 创建计划实例
	plan := &entity.Plan{
		Code:        input.Code,
		PathID:      pathID,
		Title:       strings.TrimSpace(input.Title),
		Description: strings.TrimSpace(input.Description),
		Content:     strings.TrimSpace(input.Content),
		Status:      entity.PlanStatusPending,
		Progress:    0,
	}

	// 保存到数据库
	if err := s.planModel.Create(ctx, plan); err != nil {
		return nil, err
	}

	return plan, nil
}

// UpdatePlan 更新计划
func (s *PlanService) UpdatePlan(ctx context.Context, input *dto.PlanUpdateDTO) error {
	// 参数验证
	if strings.TrimSpace(input.Code) == "" {
		return errors.New("计划 code 不能为空")
	}

	// 通过 code 获取现有计划
	plan, err := s.planModel.FindByCode(ctx, input.Code)
	if err != nil {
		return errors.New("计划不存在或已完成/取消")
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
	return s.planModel.Update(ctx, plan)
}

// DeletePlan 删除计划（通过 code）
func (s *PlanService) DeletePlan(ctx context.Context, code string) error {
	// 参数验证
	if strings.TrimSpace(code) == "" {
		return errors.New("无效的计划 code")
	}

	// 通过 code 获取计划
	plan, err := s.planModel.FindByCode(ctx, code)
	if err != nil {
		return errors.New("计划不存在或已完成/取消")
	}

	// 执行删除操作
	return s.planModel.Delete(ctx, plan.ID)
}

// DeletePlanByID 删除计划（通过 ID，TUI 内部使用）
func (s *PlanService) DeletePlanByID(ctx context.Context, id int64) error {
	// 参数验证
	if id == 0 {
		return errors.New("无效的计划ID")
	}

	// 验证计划是否存在
	_, err := s.planModel.FindByID(ctx, id)
	if err != nil {
		return errors.New("计划不存在")
	}

	// 执行删除操作
	return s.planModel.Delete(ctx, id)
}

// GetPlan 获取单个计划（通过 code）
func (s *PlanService) GetPlan(ctx context.Context, code string) (*entity.Plan, error) {
	// 参数验证
	if strings.TrimSpace(code) == "" {
		return nil, errors.New("无效的计划 code")
	}

	// 查询计划
	plan, err := s.planModel.FindByCode(ctx, code)
	if err != nil {
		return nil, errors.New("计划不存在或已完成/取消")
	}
	if plan == nil {
		return nil, errors.New("计划不存在")
	}

	return plan, nil
}

// GetPlanByID 获取单个计划（通过 ID，TUI 内部使用）
func (s *PlanService) GetPlanByID(ctx context.Context, id int64) (*entity.Plan, error) {
	// 参数验证
	if id == 0 {
		return nil, errors.New("无效的计划ID")
	}

	// 查询计划
	plan, err := s.planModel.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if plan == nil {
		return nil, errors.New("计划不存在")
	}

	return plan, nil
}

// ListPlans 获取所有计划列表（需要提供作用域上下文）
// 使用 PathOnlyFilter 进行过滤
func (s *PlanService) ListPlans(ctx context.Context, scopeCtx *types.ScopeContext) ([]entity.Plan, error) {
	filter := buildPathOnlyFilter("all", scopeCtx)
	plans, err := s.planModel.FindByPathOnlyFilter(ctx, filter)
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
// 使用 PathOnlyFilter 进行过滤（无 Global 支持）
func (s *PlanService) ListPlansByScope(ctx context.Context, scope string, scopeCtx *types.ScopeContext) ([]entity.Plan, error) {
	filter := buildPathOnlyFilter(scope, scopeCtx)
	plans, err := s.planModel.FindByPathOnlyFilter(ctx, filter)
	if err != nil {
		return nil, err
	}

	if plans == nil {
		return make([]entity.Plan, 0), nil
	}

	return plans, nil
}

// ListByStatus 根据状态获取计划列表
func (s *PlanService) ListByStatus(ctx context.Context, status entity.PlanStatus, scopeCtx *types.ScopeContext) ([]entity.Plan, error) {
	// 验证状态值是否有效
	if !isValidPlanStatus(status) {
		return nil, errors.New("无效的计划状态")
	}

	filter := buildPathOnlyFilter("all", scopeCtx)
	plans, err := s.planModel.FindByStatus(ctx, status, filter)
	if err != nil {
		return nil, err
	}

	// 如果没有计划，返回空切片而不是nil
	if plans == nil {
		return make([]entity.Plan, 0), nil
	}

	return plans, nil
}

// StartPlan 开始计划（通过 code）
func (s *PlanService) StartPlan(ctx context.Context, code string) error {
	// 参数验证
	if strings.TrimSpace(code) == "" {
		return errors.New("无效的计划 code")
	}

	// 获取计划
	plan, err := s.planModel.FindByCode(ctx, code)
	if err != nil {
		return errors.New("计划不存在或已完成/取消")
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
	return s.planModel.Update(ctx, plan)
}

// StartPlanByID 开始计划（通过 ID，TUI 内部使用）
func (s *PlanService) StartPlanByID(ctx context.Context, id int64) error {
	// 参数验证
	if id == 0 {
		return errors.New("无效的计划ID")
	}

	// 获取计划
	plan, err := s.planModel.FindByID(ctx, id)
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
	return s.planModel.Update(ctx, plan)
}

// CompletePlan 完成计划（通过 code）
func (s *PlanService) CompletePlan(ctx context.Context, code string) error {
	// 参数验证
	if strings.TrimSpace(code) == "" {
		return errors.New("无效的计划 code")
	}

	// 获取计划
	plan, err := s.planModel.FindByCode(ctx, code)
	if err != nil {
		return errors.New("计划不存在或已完成/取消")
	}

	// 验证状态转换是否合法
	if plan.Status == entity.PlanStatusCancelled {
		return errors.New("已取消的计划无法标记为完成")
	}

	// 执行完成
	plan.Complete()

	// 保存更新
	return s.planModel.Update(ctx, plan)
}

// CompletePlanByID 完成计划（通过 ID，TUI 内部使用）
func (s *PlanService) CompletePlanByID(ctx context.Context, id int64) error {
	// 参数验证
	if id == 0 {
		return errors.New("无效的计划ID")
	}

	// 获取计划
	plan, err := s.planModel.FindByID(ctx, id)
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
	return s.planModel.Update(ctx, plan)
}

// CancelPlan 取消计划（通过 code）
func (s *PlanService) CancelPlan(ctx context.Context, code string) error {
	// 参数验证
	if strings.TrimSpace(code) == "" {
		return errors.New("无效的计划 code")
	}

	// 获取计划
	plan, err := s.planModel.FindByCode(ctx, code)
	if err != nil {
		return errors.New("计划不存在或已完成/取消")
	}

	// 验证状态
	if plan.Status == entity.PlanStatusCompleted {
		return errors.New("已完成的计划无法取消")
	}

	// 执行取消
	plan.Cancel()

	// 保存更新
	return s.planModel.Update(ctx, plan)
}

// CancelPlanByID 取消计划（通过 ID，TUI 内部使用）
func (s *PlanService) CancelPlanByID(ctx context.Context, id int64) error {
	if id == 0 {
		return errors.New("无效的计划ID")
	}

	plan, err := s.planModel.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if plan.Status == entity.PlanStatusCompleted {
		return errors.New("已完成的计划无法取消")
	}

	plan.Cancel()
	return s.planModel.Update(ctx, plan)
}

// UpdateProgress 更新计划进度（通过 code）
func (s *PlanService) UpdateProgress(ctx context.Context, code string, progress int) error {
	// 参数验证
	if strings.TrimSpace(code) == "" {
		return errors.New("无效的计划 code")
	}
	if progress < 0 || progress > 100 {
		return errors.New("进度值必须在0-100之间")
	}

	// 获取计划
	plan, err := s.planModel.FindByCode(ctx, code)
	if err != nil {
		return errors.New("计划不存在或已完成/取消")
	}

	// 验证状态 - 已取消的计划不能更新进度
	if plan.Status == entity.PlanStatusCancelled {
		return errors.New("已取消的计划无法更新进度")
	}

	// 使用 Plan 类型的 UpdateProgress 方法
	plan.UpdateProgress(progress)

	// 保存更新
	return s.planModel.Update(ctx, plan)
}

// UpdateProgressByID 更新计划进度（通过 ID，TUI 内部使用）
func (s *PlanService) UpdateProgressByID(ctx context.Context, id int64, progress int) error {
	// 参数验证
	if id == 0 {
		return errors.New("无效的计划ID")
	}
	if progress < 0 || progress > 100 {
		return errors.New("进度值必须在0-100之间")
	}

	// 获取计划
	plan, err := s.planModel.FindByID(ctx, id)
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
	return s.planModel.Update(ctx, plan)
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
	_, err := s.planModel.FindByID(ctx, planID)
	if err != nil {
		return nil, errors.New("计划不存在")
	}

	return s.planModel.AddSubTask(ctx, planID, strings.TrimSpace(title), strings.TrimSpace(description))
}

// UpdateSubTask 更新子任务
func (s *PlanService) UpdateSubTask(ctx context.Context, input *dto.SubTaskUpdateDTO) error {
	if input.ID == 0 {
		return errors.New("无效的子任务ID")
	}

	subTask, err := s.planModel.GetSubTask(ctx, input.ID)
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

	return s.planModel.UpdateSubTask(ctx, subTask)
}

// DeleteSubTask 删除子任务
func (s *PlanService) DeleteSubTask(ctx context.Context, subTaskID int64) error {
	if subTaskID == 0 {
		return errors.New("无效的子任务ID")
	}
	return s.planModel.DeleteSubTask(ctx, subTaskID)
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
// 使用 PathID 判断作用域（无 Global 支持）
func ToPlanResponseDTO(plan *entity.Plan, scopeCtx *types.ScopeContext) *dto.PlanResponseDTO {
	if plan == nil {
		return nil
	}

	scope := types.GetScopeForDisplayNoGlobal(plan.PathID, scopeCtx)

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
		Code:        plan.Code,
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
		Code:        plan.Code,
		Title:       plan.Title,
		Description: plan.Description,
		Status:      string(plan.Status),
		Progress:    plan.Progress,
	}
}
