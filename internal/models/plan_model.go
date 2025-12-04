package models

import (
	"context"

	"github.com/XiaoLFeng/llm-memory/internal/database"
	"github.com/XiaoLFeng/llm-memory/internal/models/entity"
	"gorm.io/gorm"
)

// PlanModel 计划数据访问层
type PlanModel struct {
	db *gorm.DB
}

// NewPlanModel 创建 PlanModel 实例
func NewPlanModel(db *gorm.DB) *PlanModel {
	return &PlanModel{db: db}
}

// Create 创建计划
func (m *PlanModel) Create(ctx context.Context, plan *entity.Plan) error {
	plan.ID = database.GenerateID()
	return m.db.WithContext(ctx).Create(plan).Error
}

// Update 更新计划
func (m *PlanModel) Update(ctx context.Context, plan *entity.Plan) error {
	return m.db.WithContext(ctx).Save(plan).Error
}

// Delete 删除计划（硬删除）
func (m *PlanModel) Delete(ctx context.Context, id int64) error {
	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 先删除关联的子任务
		if err := tx.Where("plan_id = ?", id).Unscoped().Delete(&entity.SubTask{}).Error; err != nil {
			return err
		}
		// 硬删除计划本身
		return tx.Unscoped().Delete(&entity.Plan{}, id).Error
	})
}

// FindByID 根据 ID 查找计划
func (m *PlanModel) FindByID(ctx context.Context, id int64) (*entity.Plan, error) {
	var plan entity.Plan
	err := m.db.WithContext(ctx).Preload("SubTasks", func(db *gorm.DB) *gorm.DB {
		return db.Order("sort_order ASC")
	}).First(&plan, id).Error
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

// FindByCode 根据 code 查找活跃的计划
// 只查询 pending 和 in_progress 状态的计划（排除 completed 和 cancelled）
func (m *PlanModel) FindByCode(ctx context.Context, code string) (*entity.Plan, error) {
	var plan entity.Plan
	err := m.db.WithContext(ctx).
		Preload("SubTasks", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC")
		}).
		Where("code = ?", code).
		Where("status NOT IN ?", []entity.PlanStatus{entity.PlanStatusCompleted, entity.PlanStatusCancelled}).
		First(&plan).Error
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

// ExistsActiveCode 检查活跃记录中是否存在指定 code
// 只检查活跃状态（status NOT IN completed/cancelled）
// excludeID: 如果 > 0，则排除该 ID（用于更新时检查）
func (m *PlanModel) ExistsActiveCode(ctx context.Context, code string, excludeID int64) (bool, error) {
	var count int64
	query := m.db.WithContext(ctx).
		Model(&entity.Plan{}).
		Where("code = ?", code).
		Where("status NOT IN ?", []entity.PlanStatus{entity.PlanStatusCompleted, entity.PlanStatusCancelled})

	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}

	err := query.Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// FindAll 查找所有计划
func (m *PlanModel) FindAll(ctx context.Context) ([]entity.Plan, error) {
	return m.FindByFilter(ctx, DefaultVisibilityFilter())
}

// FindByStatus 根据状态查找计划
func (m *PlanModel) FindByStatus(ctx context.Context, status entity.PlanStatus) ([]entity.Plan, error) {
	filter := DefaultVisibilityFilter()
	var plans []entity.Plan
	err := applyVisibilityFilter(m.db.WithContext(ctx).Preload("SubTasks", func(db *gorm.DB) *gorm.DB {
		return db.Order("sort_order ASC")
	}), filter).Where("status = ?", status).Order("created_at DESC").Find(&plans).Error
	return plans, err
}

// FindByScope 根据作用域查找计划
// 纯关联模式：基于 PathID 进行查询
// pathID: 当前路径的 PathID（0 表示无路径）
// groupPathIDs: 组内所有路径 ID 列表（空切片表示无组）
// includeGlobal: 是否包含全局计划
// 默认排除已完成和已取消的计划
func (m *PlanModel) FindByScope(ctx context.Context, pathID int64, groupPathIDs []int64, includeGlobal bool) ([]entity.Plan, error) {
	filter := VisibilityFilter{
		IncludeGlobal:    includeGlobal,
		IncludeNonGlobal: true,
		PathIDs:          MergePathIDs(pathID, groupPathIDs),
	}
	return m.FindByFilter(ctx, filter)
}

// FindByFilter 统一过滤器查询计划
// 默认排除已完成和已取消的计划
func (m *PlanModel) FindByFilter(ctx context.Context, filter VisibilityFilter) ([]entity.Plan, error) {
	var plans []entity.Plan
	err := applyVisibilityFilter(m.db.WithContext(ctx).Preload("SubTasks", func(db *gorm.DB) *gorm.DB {
		return db.Order("sort_order ASC")
	}), filter).
		Where("status NOT IN ?", []entity.PlanStatus{entity.PlanStatusCompleted, entity.PlanStatusCancelled}).
		Order("created_at DESC").
		Find(&plans).Error
	return plans, err
}

// UpdateProgress 更新计划进度
func (m *PlanModel) UpdateProgress(ctx context.Context, id int64, progress int) error {
	plan, err := m.FindByID(ctx, id)
	if err != nil {
		return err
	}
	plan.UpdateProgress(progress)
	return m.Update(ctx, plan)
}

// Start 开始计划
func (m *PlanModel) Start(ctx context.Context, id int64) error {
	plan, err := m.FindByID(ctx, id)
	if err != nil {
		return err
	}
	plan.Start()
	return m.Update(ctx, plan)
}

// Complete 完成计划
func (m *PlanModel) Complete(ctx context.Context, id int64) error {
	plan, err := m.FindByID(ctx, id)
	if err != nil {
		return err
	}
	plan.Complete()
	return m.Update(ctx, plan)
}

// Cancel 取消计划
func (m *PlanModel) Cancel(ctx context.Context, id int64) error {
	plan, err := m.FindByID(ctx, id)
	if err != nil {
		return err
	}
	plan.Cancel()
	return m.Update(ctx, plan)
}

// AddSubTask 添加子任务
func (m *PlanModel) AddSubTask(ctx context.Context, planID int64, title, description string) (*entity.SubTask, error) {
	// 获取当前最大排序值
	var maxOrder int
	m.db.WithContext(ctx).Model(&entity.SubTask{}).Where("plan_id = ?", planID).Select("COALESCE(MAX(sort_order), 0)").Scan(&maxOrder)

	subTask := &entity.SubTask{
		ID:          database.GenerateID(),
		PlanID:      planID,
		Title:       title,
		Description: description,
		Status:      entity.PlanStatusPending,
		Progress:    0,
		SortOrder:   maxOrder + 1,
	}

	if err := m.db.WithContext(ctx).Create(subTask).Error; err != nil {
		return nil, err
	}
	return subTask, nil
}

// UpdateSubTask 更新子任务
func (m *PlanModel) UpdateSubTask(ctx context.Context, subTask *entity.SubTask) error {
	return m.db.WithContext(ctx).Save(subTask).Error
}

// DeleteSubTask 删除子任务（硬删除）
func (m *PlanModel) DeleteSubTask(ctx context.Context, subTaskID int64) error {
	return m.db.WithContext(ctx).Unscoped().Delete(&entity.SubTask{}, subTaskID).Error
}

// GetSubTask 获取子任务
func (m *PlanModel) GetSubTask(ctx context.Context, subTaskID int64) (*entity.SubTask, error) {
	var subTask entity.SubTask
	err := m.db.WithContext(ctx).First(&subTask, subTaskID).Error
	if err != nil {
		return nil, err
	}
	return &subTask, nil
}

// Count 获取计划总数
func (m *PlanModel) Count(ctx context.Context) (int64, error) {
	var count int64
	err := m.db.WithContext(ctx).Model(&entity.Plan{}).Count(&count).Error
	return count, err
}
