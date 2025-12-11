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
// 注意：关联的 Todo 会通过 GORM 的 CASCADE 约束自动删除
func (m *PlanModel) Delete(ctx context.Context, id int64) error {
	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 先删除关联的 Todo 标签
		var todoIDs []int64
		if err := tx.Model(&entity.ToDo{}).Where("plan_id = ?", id).Pluck("id", &todoIDs).Error; err != nil {
			return err
		}
		if len(todoIDs) > 0 {
			if err := tx.Where("to_do_id IN ?", todoIDs).Delete(&entity.ToDoTag{}).Error; err != nil {
				return err
			}
		}
		// 删除关联的 Todo
		if err := tx.Where("plan_id = ?", id).Unscoped().Delete(&entity.ToDo{}).Error; err != nil {
			return err
		}
		// 硬删除计划本身
		return tx.Unscoped().Delete(&entity.Plan{}, id).Error
	})
}

// FindByID 根据 ID 查找计划
func (m *PlanModel) FindByID(ctx context.Context, id int64) (*entity.Plan, error) {
	var plan entity.Plan
	err := m.db.WithContext(ctx).Preload("Todos", func(db *gorm.DB) *gorm.DB {
		return db.Preload("Tags").Order("sort_order ASC, priority DESC")
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
		Preload("Todos", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Tags").Order("sort_order ASC, priority DESC")
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

// FindAll 查找所有计划（需要提供路径过滤器）
func (m *PlanModel) FindAll(ctx context.Context, filter PathOnlyVisibilityFilter) ([]entity.Plan, error) {
	return m.FindByPathOnlyFilter(ctx, filter)
}

// FindByStatus 根据状态查找计划
func (m *PlanModel) FindByStatus(ctx context.Context, status entity.PlanStatus, filter PathOnlyVisibilityFilter) ([]entity.Plan, error) {
	var plans []entity.Plan
	err := ApplyPathOnlyFilter(m.db.WithContext(ctx).Preload("Todos", func(db *gorm.DB) *gorm.DB {
		return db.Preload("Tags").Order("sort_order ASC, priority DESC")
	}), filter).Where("status = ?", status).Order("created_at DESC").Find(&plans).Error
	return plans, err
}

// FindByScope 根据作用域查找计划（无 Global 支持）
// pathIDs: 路径 ID 列表（个人路径 + 组内路径）
// 默认排除已完成和已取消的计划
func (m *PlanModel) FindByScope(ctx context.Context, pathIDs []int64) ([]entity.Plan, error) {
	filter := PathOnlyVisibilityFilter{
		PathIDs: pathIDs,
	}
	return m.FindByPathOnlyFilter(ctx, filter)
}

// FindByFilter 根据统一过滤器查询计划（兼容旧接口，供 Memory 使用）
// 默认排除已完成和已取消的计划
func (m *PlanModel) FindByFilter(ctx context.Context, filter VisibilityFilter) ([]entity.Plan, error) {
	var plans []entity.Plan
	err := applyVisibilityFilter(m.db.WithContext(ctx).Preload("Todos", func(db *gorm.DB) *gorm.DB {
		return db.Preload("Tags").Order("sort_order ASC, priority DESC")
	}), filter).
		Where("status NOT IN ?", []entity.PlanStatus{entity.PlanStatusCompleted, entity.PlanStatusCancelled}).
		Order("created_at DESC").
		Find(&plans).Error
	return plans, err
}

// FindByPathOnlyFilter 根据路径过滤器查询计划（供 Plan 使用）
// 默认排除已完成和已取消的计划（Plan 保持隐藏机制）
func (m *PlanModel) FindByPathOnlyFilter(ctx context.Context, filter PathOnlyVisibilityFilter) ([]entity.Plan, error) {
	var plans []entity.Plan
	err := ApplyPathOnlyFilter(m.db.WithContext(ctx).Preload("Todos", func(db *gorm.DB) *gorm.DB {
		return db.Preload("Tags").Order("sort_order ASC, priority DESC")
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

// Count 获取计划总数
func (m *PlanModel) Count(ctx context.Context) (int64, error) {
	var count int64
	err := m.db.WithContext(ctx).Model(&entity.Plan{}).Count(&count).Error
	return count, err
}
