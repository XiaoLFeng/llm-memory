package models

import (
	"context"
	"fmt"
	"time"

	"github.com/XiaoLFeng/llm-memory/internal/database"
	"github.com/XiaoLFeng/llm-memory/internal/models/dto"
	"github.com/XiaoLFeng/llm-memory/internal/models/entity"
	"gorm.io/gorm"
)

// ToDoModel 待办数据访问层
type ToDoModel struct {
	db *gorm.DB
}

// NewToDoModel 创建 ToDoModel 实例
func NewToDoModel(db *gorm.DB) *ToDoModel {
	return &ToDoModel{db: db}
}

// Create 创建待办
func (m *ToDoModel) Create(ctx context.Context, todo *entity.ToDo) error {
	todo.ID = database.GenerateID()
	return m.db.WithContext(ctx).Create(todo).Error
}

// Update 更新待办
func (m *ToDoModel) Update(ctx context.Context, todo *entity.ToDo) error {
	return m.db.WithContext(ctx).Save(todo).Error
}

// Delete 删除待办（硬删除）
func (m *ToDoModel) Delete(ctx context.Context, id int64) error {
	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 先删除关联的标签
		if err := tx.Where("to_do_id = ?", id).Unscoped().Delete(&entity.ToDoTag{}).Error; err != nil {
			return err
		}
		// 硬删除待办本身
		return tx.Unscoped().Delete(&entity.ToDo{}, id).Error
	})
}

// FindByID 根据 ID 查找待办
func (m *ToDoModel) FindByID(ctx context.Context, id int64) (*entity.ToDo, error) {
	var todo entity.ToDo
	err := m.db.WithContext(ctx).Preload("Tags").First(&todo, id).Error
	if err != nil {
		return nil, err
	}
	return &todo, nil
}

// FindByCode 根据 code 查找活跃的待办
// 只查询活跃状态（待处理和进行中），排除已完成和已取消的记录
func (m *ToDoModel) FindByCode(ctx context.Context, code string) (*entity.ToDo, error) {
	var todo entity.ToDo
	err := m.db.WithContext(ctx).
		Preload("Tags").
		Where("code = ? AND status NOT IN ?", code, []entity.ToDoStatus{entity.ToDoStatusCompleted, entity.ToDoStatusCancelled}).
		First(&todo).Error
	if err != nil {
		return nil, err
	}
	return &todo, nil
}

// ExistsActiveCode 检查活跃记录中是否存在指定 code
// 只检查活跃状态（待处理和进行中）
// excludeID: 如果 > 0，则排除该 ID（用于更新时检查重复）
func (m *ToDoModel) ExistsActiveCode(ctx context.Context, code string, excludeID int64) (bool, error) {
	var count int64
	query := m.db.WithContext(ctx).Model(&entity.ToDo{}).
		Where("code = ? AND status NOT IN ?", code, []entity.ToDoStatus{entity.ToDoStatusCompleted, entity.ToDoStatusCancelled})

	// 如果提供了 excludeID，则排除该 ID
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}

	err := query.Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// FindAll 查找所有待办
func (m *ToDoModel) FindAll(ctx context.Context) ([]entity.ToDo, error) {
	return m.FindByFilter(ctx, DefaultVisibilityFilter())
}

// FindByStatus 根据状态查找待办
func (m *ToDoModel) FindByStatus(ctx context.Context, status entity.ToDoStatus) ([]entity.ToDo, error) {
	filter := DefaultVisibilityFilter()
	var todos []entity.ToDo
	err := applyVisibilityFilter(m.db.WithContext(ctx).Preload("Tags"), filter).
		Where("status = ?", status).
		Order("priority DESC, created_at DESC").
		Find(&todos).Error
	return todos, err
}

// FindByScope 根据作用域查找待办
// 纯关联模式：基于 PathID 进行查询
// pathID: 当前路径的 PathID（0 表示无路径）
// groupPathIDs: 组内所有路径 ID 列表（空切片表示无组）
// includeGlobal: 是否包含全局待办
// 默认排除已完成和已取消的记录
func (m *ToDoModel) FindByScope(ctx context.Context, pathID int64, groupPathIDs []int64, includeGlobal bool) ([]entity.ToDo, error) {
	filter := VisibilityFilter{
		IncludeGlobal:    includeGlobal,
		IncludeNonGlobal: true,
		PathIDs:          mergePathIDs(pathID, groupPathIDs),
	}
	return m.FindByFilter(ctx, filter)
}

// FindByFilter 根据统一过滤器查询待办
// 默认排除已完成和已取消的记录
func (m *ToDoModel) FindByFilter(ctx context.Context, filter VisibilityFilter) ([]entity.ToDo, error) {
	var todos []entity.ToDo
	err := applyVisibilityFilter(m.db.WithContext(ctx).Preload("Tags"), filter).
		Where("status NOT IN ?", []entity.ToDoStatus{entity.ToDoStatusCompleted, entity.ToDoStatusCancelled}).
		Order("priority DESC, created_at DESC").
		Find(&todos).Error
	return todos, err
}

// Complete 完成待办
func (m *ToDoModel) Complete(ctx context.Context, id int64) error {
	now := time.Now()
	return m.db.WithContext(ctx).Model(&entity.ToDo{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":       entity.ToDoStatusCompleted,
		"completed_at": now,
	}).Error
}

// Start 开始待办
func (m *ToDoModel) Start(ctx context.Context, id int64) error {
	return m.db.WithContext(ctx).Model(&entity.ToDo{}).Where("id = ?", id).Update("status", entity.ToDoStatusInProgress).Error
}

// Cancel 取消待办
func (m *ToDoModel) Cancel(ctx context.Context, id int64) error {
	return m.db.WithContext(ctx).Model(&entity.ToDo{}).Where("id = ?", id).Update("status", entity.ToDoStatusCancelled).Error
}

// UpdateTags 更新待办标签
// 删除旧标签并添加新标签
func (m *ToDoModel) UpdateTags(ctx context.Context, todoID int64, tags []string) error {
	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 删除旧标签
		if err := tx.Where("to_do_id = ?", todoID).Delete(&entity.ToDoTag{}).Error; err != nil {
			return err
		}
		// 添加新标签
		for _, tag := range tags {
			todoTag := entity.ToDoTag{
				ID:     database.GenerateID(),
				ToDoID: todoID,
				Tag:    tag,
			}
			if err := tx.Create(&todoTag).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// BatchCreate 批量创建待办
func (m *ToDoModel) BatchCreate(ctx context.Context, todos []entity.ToDo) (*dto.ToDoBatchResultDTO, error) {
	result := &dto.ToDoBatchResultDTO{
		Total:  len(todos),
		Errors: make([]string, 0),
	}

	// 使用事务批量插入
	err := m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for i := range todos {
			todos[i].ID = database.GenerateID()
			if err := tx.Create(&todos[i]).Error; err != nil {
				result.Failed++
				result.Errors = append(result.Errors,
					fmt.Sprintf("第 %d 项创建失败: %s", i+1, err.Error()))
			} else {
				result.Succeeded++
			}
		}
		return nil
	})

	return result, err
}

// BatchUpdate 批量更新待办
func (m *ToDoModel) BatchUpdate(ctx context.Context, updates []dto.ToDoUpdateDTO) (*dto.ToDoBatchResultDTO, error) {
	result := &dto.ToDoBatchResultDTO{
		Total:  len(updates),
		Errors: make([]string, 0),
	}

	err := m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for i, update := range updates {
			todo, err := m.FindByCode(ctx, update.Code)
			if err != nil {
				result.Failed++
				result.Errors = append(result.Errors,
					fmt.Sprintf("第 %d 项（Code=%s）不存在", i+1, update.Code))
				continue
			}

			// 应用更新
			if update.Title != nil {
				todo.Title = *update.Title
			}
			if update.Description != nil {
				todo.Description = *update.Description
			}
			if update.Priority != nil {
				todo.Priority = entity.ToDoPriority(*update.Priority)
			}
			if update.Status != nil {
				todo.Status = entity.ToDoStatus(*update.Status)
				if todo.Status == entity.ToDoStatusCompleted {
					now := time.Now()
					todo.CompletedAt = &now
				}
			}
			if update.DueDate != nil {
				todo.DueDate = update.DueDate
			}

			if err := tx.Save(todo).Error; err != nil {
				result.Failed++
				result.Errors = append(result.Errors,
					fmt.Sprintf("第 %d 项（Code=%s）更新失败: %s", i+1, update.Code, err.Error()))
			} else {
				result.Succeeded++
			}
		}
		return nil
	})

	return result, err
}

// BatchComplete 批量完成待办
func (m *ToDoModel) BatchComplete(ctx context.Context, ids []int64) (*dto.ToDoBatchResultDTO, error) {
	result := &dto.ToDoBatchResultDTO{
		Total:  len(ids),
		Errors: make([]string, 0),
	}

	now := time.Now()

	err := m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, id := range ids {
			updateResult := tx.Model(&entity.ToDo{}).
				Where("id = ? AND status != ?", id, entity.ToDoStatusCompleted).
				Updates(map[string]interface{}{
					"status":       entity.ToDoStatusCompleted,
					"completed_at": now,
				})

			if updateResult.Error != nil {
				result.Failed++
				result.Errors = append(result.Errors,
					fmt.Sprintf("ID=%d 完成失败: %s", id, updateResult.Error.Error()))
			} else if updateResult.RowsAffected == 0 {
				result.Failed++
				result.Errors = append(result.Errors,
					fmt.Sprintf("ID=%d 不存在或已完成", id))
			} else {
				result.Succeeded++
			}
		}
		return nil
	})

	return result, err
}

// BatchDelete 批量删除待办
func (m *ToDoModel) BatchDelete(ctx context.Context, ids []int64) (*dto.ToDoBatchResultDTO, error) {
	result := &dto.ToDoBatchResultDTO{
		Total:  len(ids),
		Errors: make([]string, 0),
	}

	err := m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, id := range ids {
			deleteResult := tx.Where("id = ?", id).Delete(&entity.ToDo{})
			if deleteResult.Error != nil {
				result.Failed++
				result.Errors = append(result.Errors,
					fmt.Sprintf("ID=%d 删除失败: %s", id, deleteResult.Error.Error()))
			} else if deleteResult.RowsAffected == 0 {
				result.Failed++
				result.Errors = append(result.Errors,
					fmt.Sprintf("ID=%d 不存在", id))
			} else {
				result.Succeeded++
			}
		}
		return nil
	})

	return result, err
}

// BatchUpdateStatus 批量更新状态
func (m *ToDoModel) BatchUpdateStatus(ctx context.Context, ids []int64, status entity.ToDoStatus) (*dto.ToDoBatchResultDTO, error) {
	result := &dto.ToDoBatchResultDTO{
		Total:  len(ids),
		Errors: make([]string, 0),
	}

	updates := map[string]interface{}{
		"status": status,
	}
	if status == entity.ToDoStatusCompleted {
		updates["completed_at"] = time.Now()
	}

	err := m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		updateResult := tx.Model(&entity.ToDo{}).Where("id IN ?", ids).Updates(updates)
		if updateResult.Error != nil {
			result.Failed = len(ids)
			result.Errors = append(result.Errors, updateResult.Error.Error())
		} else {
			result.Succeeded = int(updateResult.RowsAffected)
			result.Failed = len(ids) - result.Succeeded
		}
		return nil
	})

	return result, err
}

// Count 获取待办总数
func (m *ToDoModel) Count(ctx context.Context) (int64, error) {
	var count int64
	err := m.db.WithContext(ctx).Model(&entity.ToDo{}).Count(&count).Error
	return count, err
}
