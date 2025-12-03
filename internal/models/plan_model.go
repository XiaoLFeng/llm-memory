package models

import (
	"context"
	"strings"

	"github.com/XiaoLFeng/llm-memory/internal/models/entity"
	"gorm.io/gorm"
)

// PlanModel 计划数据访问层
// 嘿嘿~ 这是计划的数据访问模型！💖
type PlanModel struct {
	db *gorm.DB
}

// NewPlanModel 创建 PlanModel 实例
func NewPlanModel(db *gorm.DB) *PlanModel {
	return &PlanModel{db: db}
}

// Create 创建计划
func (m *PlanModel) Create(ctx context.Context, plan *entity.Plan) error {
	return m.db.WithContext(ctx).Create(plan).Error
}

// Update 更新计划
func (m *PlanModel) Update(ctx context.Context, plan *entity.Plan) error {
	return m.db.WithContext(ctx).Save(plan).Error
}

// Delete 删除计划（软删除）
func (m *PlanModel) Delete(ctx context.Context, id uint) error {
	return m.db.WithContext(ctx).Delete(&entity.Plan{}, id).Error
}

// FindByID 根据 ID 查找计划
func (m *PlanModel) FindByID(ctx context.Context, id uint) (*entity.Plan, error) {
	var plan entity.Plan
	err := m.db.WithContext(ctx).Preload("SubTasks", func(db *gorm.DB) *gorm.DB {
		return db.Order("sort_order ASC")
	}).First(&plan, id).Error
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

// FindAll 查找所有计划
func (m *PlanModel) FindAll(ctx context.Context) ([]entity.Plan, error) {
	var plans []entity.Plan
	err := m.db.WithContext(ctx).Preload("SubTasks", func(db *gorm.DB) *gorm.DB {
		return db.Order("sort_order ASC")
	}).Order("created_at DESC").Find(&plans).Error
	return plans, err
}

// FindByStatus 根据状态查找计划
func (m *PlanModel) FindByStatus(ctx context.Context, status entity.PlanStatus) ([]entity.Plan, error) {
	var plans []entity.Plan
	err := m.db.WithContext(ctx).Preload("SubTasks", func(db *gorm.DB) *gorm.DB {
		return db.Order("sort_order ASC")
	}).Where("status = ?", status).Order("created_at DESC").Find(&plans).Error
	return plans, err
}

// FindByScope 根据作用域查找计划
// 呀~ 支持 Personal/Group/Global 三层作用域过滤！✨
func (m *PlanModel) FindByScope(ctx context.Context, groupID uint, path string, includeGlobal bool) ([]entity.Plan, error) {
	var plans []entity.Plan
	query := m.db.WithContext(ctx).Preload("SubTasks", func(db *gorm.DB) *gorm.DB {
		return db.Order("sort_order ASC")
	})

	// 构建作用域条件
	var conditions []string
	var args []interface{}

	if path != "" {
		conditions = append(conditions, "(path = ?)")
		args = append(args, path)
	}
	if groupID > 0 {
		conditions = append(conditions, "(group_id = ? AND path = '')")
		args = append(args, groupID)
	}
	if includeGlobal {
		conditions = append(conditions, "(group_id = 0 AND path = '')")
	}

	if len(conditions) > 0 {
		query = query.Where(strings.Join(conditions, " OR "), args...)
	}

	err := query.Order("created_at DESC").Find(&plans).Error
	return plans, err
}

// UpdateProgress 更新计划进度
func (m *PlanModel) UpdateProgress(ctx context.Context, id uint, progress int) error {
	plan, err := m.FindByID(ctx, id)
	if err != nil {
		return err
	}
	plan.UpdateProgress(progress)
	return m.Update(ctx, plan)
}

// Start 开始计划
func (m *PlanModel) Start(ctx context.Context, id uint) error {
	plan, err := m.FindByID(ctx, id)
	if err != nil {
		return err
	}
	plan.Start()
	return m.Update(ctx, plan)
}

// Complete 完成计划
func (m *PlanModel) Complete(ctx context.Context, id uint) error {
	plan, err := m.FindByID(ctx, id)
	if err != nil {
		return err
	}
	plan.Complete()
	return m.Update(ctx, plan)
}

// Cancel 取消计划
func (m *PlanModel) Cancel(ctx context.Context, id uint) error {
	plan, err := m.FindByID(ctx, id)
	if err != nil {
		return err
	}
	plan.Cancel()
	return m.Update(ctx, plan)
}

// AddSubTask 添加子任务
func (m *PlanModel) AddSubTask(ctx context.Context, planID uint, title, description string) (*entity.SubTask, error) {
	// 获取当前最大排序值
	var maxOrder int
	m.db.WithContext(ctx).Model(&entity.SubTask{}).Where("plan_id = ?", planID).Select("COALESCE(MAX(sort_order), 0)").Scan(&maxOrder)

	subTask := &entity.SubTask{
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

// DeleteSubTask 删除子任务
func (m *PlanModel) DeleteSubTask(ctx context.Context, subTaskID uint) error {
	return m.db.WithContext(ctx).Delete(&entity.SubTask{}, subTaskID).Error
}

// GetSubTask 获取子任务
func (m *PlanModel) GetSubTask(ctx context.Context, subTaskID uint) (*entity.SubTask, error) {
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
