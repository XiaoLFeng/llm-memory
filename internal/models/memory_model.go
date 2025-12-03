package models

import (
	"context"
	"strings"

	"github.com/XiaoLFeng/llm-memory/internal/database"
	"github.com/XiaoLFeng/llm-memory/internal/models/entity"
	"gorm.io/gorm"
)

// MemoryModel 记忆数据访问层
type MemoryModel struct {
	db *gorm.DB
}

// NewMemoryModel 创建 MemoryModel 实例
func NewMemoryModel(db *gorm.DB) *MemoryModel {
	return &MemoryModel{db: db}
}

// Create 创建记忆
func (m *MemoryModel) Create(ctx context.Context, memory *entity.Memory) error {
	memory.ID = database.GenerateID()
	return m.db.WithContext(ctx).Create(memory).Error
}

// Update 更新记忆
func (m *MemoryModel) Update(ctx context.Context, memory *entity.Memory) error {
	return m.db.WithContext(ctx).Save(memory).Error
}

// Delete 删除记忆（硬删除）
func (m *MemoryModel) Delete(ctx context.Context, id int64) error {
	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 先删除关联的标签
		if err := tx.Where("memory_id = ?", id).Unscoped().Delete(&entity.MemoryTag{}).Error; err != nil {
			return err
		}
		// 硬删除记忆本身
		return tx.Unscoped().Delete(&entity.Memory{}, id).Error
	})
}

// FindByID 根据 ID 查找记忆
func (m *MemoryModel) FindByID(ctx context.Context, id int64) (*entity.Memory, error) {
	var memory entity.Memory
	err := m.db.WithContext(ctx).Preload("Tags").First(&memory, id).Error
	if err != nil {
		return nil, err
	}
	return &memory, nil
}

// FindAll 查找所有记忆
func (m *MemoryModel) FindAll(ctx context.Context) ([]entity.Memory, error) {
	var memories []entity.Memory
	err := m.db.WithContext(ctx).Preload("Tags").Order("created_at DESC").Find(&memories).Error
	return memories, err
}

// FindByCategory 根据分类查找记忆
func (m *MemoryModel) FindByCategory(ctx context.Context, category string) ([]entity.Memory, error) {
	var memories []entity.Memory
	err := m.db.WithContext(ctx).Preload("Tags").Where("category = ?", category).Order("created_at DESC").Find(&memories).Error
	return memories, err
}

// FindByScope 根据作用域查找记忆
// 纯关联模式：基于 PathID 进行查询
// pathID: 当前路径的 PathID（0 表示无路径）
// groupPathIDs: 组内所有路径 ID 列表（空切片表示无组）
// includeGlobal: 是否包含全局记忆
func (m *MemoryModel) FindByScope(ctx context.Context, pathID int64, groupPathIDs []int64, includeGlobal bool) ([]entity.Memory, error) {
	var memories []entity.Memory
	query := m.db.WithContext(ctx).Preload("Tags")

	// 构建作用域条件
	var conditions []string
	var args []interface{}

	// Personal: 当前路径
	if pathID > 0 {
		conditions = append(conditions, "(path_id = ?)")
		args = append(args, pathID)
	}

	// Group: 组内所有路径（排除已包含的 pathID 避免重复）
	if len(groupPathIDs) > 0 {
		conditions = append(conditions, "(path_id IN ?)")
		args = append(args, groupPathIDs)
	}

	// Global: PathID = 0
	if includeGlobal {
		conditions = append(conditions, "(path_id = 0)")
	}

	if len(conditions) > 0 {
		query = query.Where(strings.Join(conditions, " OR "), args...)
	} else {
		// 无条件时默认返回全局数据（path_id = 0）以避免调用方因 scopeCtx 为空而取不到数据
		query = query.Where("path_id = 0")
	}

	err := query.Order("created_at DESC").Find(&memories).Error
	return memories, err
}

// Search 搜索记忆（在标题和内容中搜索）
func (m *MemoryModel) Search(ctx context.Context, keyword string) ([]entity.Memory, error) {
	var memories []entity.Memory
	pattern := "%" + keyword + "%"
	err := m.db.WithContext(ctx).Preload("Tags").
		Where("title LIKE ? OR content LIKE ?", pattern, pattern).
		Order("created_at DESC").
		Find(&memories).Error
	return memories, err
}

// SearchByScope 在指定作用域内搜索记忆
// 纯关联模式：基于 PathID 进行查询
func (m *MemoryModel) SearchByScope(ctx context.Context, keyword string, pathID int64, groupPathIDs []int64, includeGlobal bool) ([]entity.Memory, error) {
	var memories []entity.Memory
	pattern := "%" + keyword + "%"
	query := m.db.WithContext(ctx).Preload("Tags").
		Where("title LIKE ? OR content LIKE ?", pattern, pattern)

	// 构建作用域条件
	var conditions []string
	var args []interface{}

	// Personal: 当前路径
	if pathID > 0 {
		conditions = append(conditions, "(path_id = ?)")
		args = append(args, pathID)
	}

	// Group: 组内所有路径
	if len(groupPathIDs) > 0 {
		conditions = append(conditions, "(path_id IN ?)")
		args = append(args, groupPathIDs)
	}

	// Global: PathID = 0
	if includeGlobal {
		conditions = append(conditions, "(path_id = 0)")
	}

	if len(conditions) > 0 {
		query = query.Where(strings.Join(conditions, " OR "), args...)
	} else {
		return memories, nil
	}

	err := query.Order("created_at DESC").Find(&memories).Error
	return memories, err
}

// Archive 归档记忆
func (m *MemoryModel) Archive(ctx context.Context, id int64) error {
	return m.db.WithContext(ctx).Model(&entity.Memory{}).Where("id = ?", id).Update("is_archived", true).Error
}

// Unarchive 取消归档记忆
func (m *MemoryModel) Unarchive(ctx context.Context, id int64) error {
	return m.db.WithContext(ctx).Model(&entity.Memory{}).Where("id = ?", id).Update("is_archived", false).Error
}

// UpdateTags 更新记忆标签
// 先删除旧标签再添加新标签
func (m *MemoryModel) UpdateTags(ctx context.Context, memoryID int64, tags []string) error {
	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 删除旧标签
		if err := tx.Where("memory_id = ?", memoryID).Delete(&entity.MemoryTag{}).Error; err != nil {
			return err
		}
		// 添加新标签
		for _, tag := range tags {
			memoryTag := entity.MemoryTag{
				ID:       database.GenerateID(),
				MemoryID: memoryID,
				Tag:      tag,
			}
			if err := tx.Create(&memoryTag).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// Count 获取记忆总数
func (m *MemoryModel) Count(ctx context.Context) (int64, error) {
	var count int64
	err := m.db.WithContext(ctx).Model(&entity.Memory{}).Count(&count).Error
	return count, err
}
