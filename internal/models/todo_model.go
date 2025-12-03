package models

import (
	"context"
	"fmt"
	"strings"
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

// Delete 删除待办（软删除）
func (m *ToDoModel) Delete(ctx context.Context, id int64) error {
	return m.db.WithContext(ctx).Delete(&entity.ToDo{}, id).Error
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

// FindAll 查找所有待办
func (m *ToDoModel) FindAll(ctx context.Context) ([]entity.ToDo, error) {
	var todos []entity.ToDo
	err := m.db.WithContext(ctx).Preload("Tags").Order("priority DESC, created_at DESC").Find(&todos).Error
	return todos, err
}

// FindByStatus 根据状态查找待办
func (m *ToDoModel) FindByStatus(ctx context.Context, status entity.ToDoStatus) ([]entity.ToDo, error) {
	var todos []entity.ToDo
	err := m.db.WithContext(ctx).Preload("Tags").Where("status = ?", status).Order("priority DESC, created_at DESC").Find(&todos).Error
	return todos, err
}

// FindByScope 根据作用域查找待办
// 支持 Personal/Group/Global 三层作用域过滤
func (m *ToDoModel) FindByScope(ctx context.Context, groupID int64, path string, includeGlobal bool) ([]entity.ToDo, error) {
	var todos []entity.ToDo
	query := m.db.WithContext(ctx).Preload("Tags")

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

	err := query.Order("priority DESC, created_at DESC").Find(&todos).Error
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
			todo, err := m.FindByID(ctx, update.ID)
			if err != nil {
				result.Failed++
				result.Errors = append(result.Errors,
					fmt.Sprintf("第 %d 项（ID=%d）不存在", i+1, update.ID))
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
					fmt.Sprintf("第 %d 项（ID=%d）更新失败: %s", i+1, update.ID, err.Error()))
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
		deleteResult := tx.Delete(&entity.ToDo{}, ids)
		if deleteResult.Error != nil {
			result.Failed = len(ids)
			result.Errors = append(result.Errors, deleteResult.Error.Error())
		} else {
			result.Succeeded = int(deleteResult.RowsAffected)
			result.Failed = len(ids) - result.Succeeded
			if result.Failed > 0 {
				result.Errors = append(result.Errors,
					fmt.Sprintf("%d 个待办不存在或已删除", result.Failed))
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
