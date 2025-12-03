package models

import (
	"context"
	"strings"

	"github.com/XiaoLFeng/llm-memory/internal/models/entity"
	"gorm.io/gorm"
)

// MemoryModel 记忆数据访问层
// 嘿嘿~ 这是记忆的数据访问模型！💖
type MemoryModel struct {
	db *gorm.DB
}

// NewMemoryModel 创建 MemoryModel 实例
func NewMemoryModel(db *gorm.DB) *MemoryModel {
	return &MemoryModel{db: db}
}

// Create 创建记忆
func (m *MemoryModel) Create(ctx context.Context, memory *entity.Memory) error {
	return m.db.WithContext(ctx).Create(memory).Error
}

// Update 更新记忆
func (m *MemoryModel) Update(ctx context.Context, memory *entity.Memory) error {
	return m.db.WithContext(ctx).Save(memory).Error
}

// Delete 删除记忆（软删除）
func (m *MemoryModel) Delete(ctx context.Context, id uint) error {
	return m.db.WithContext(ctx).Delete(&entity.Memory{}, id).Error
}

// FindByID 根据 ID 查找记忆
func (m *MemoryModel) FindByID(ctx context.Context, id uint) (*entity.Memory, error) {
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
// 呀~ 支持 Personal/Group/Global 三层作用域过滤！✨
func (m *MemoryModel) FindByScope(ctx context.Context, groupID uint, path string, includeGlobal bool) ([]entity.Memory, error) {
	var memories []entity.Memory
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
func (m *MemoryModel) SearchByScope(ctx context.Context, keyword string, groupID uint, path string, includeGlobal bool) ([]entity.Memory, error) {
	var memories []entity.Memory
	pattern := "%" + keyword + "%"
	query := m.db.WithContext(ctx).Preload("Tags").
		Where("title LIKE ? OR content LIKE ?", pattern, pattern)

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

	err := query.Order("created_at DESC").Find(&memories).Error
	return memories, err
}

// Archive 归档记忆
func (m *MemoryModel) Archive(ctx context.Context, id uint) error {
	return m.db.WithContext(ctx).Model(&entity.Memory{}).Where("id = ?", id).Update("is_archived", true).Error
}

// Unarchive 取消归档记忆
func (m *MemoryModel) Unarchive(ctx context.Context, id uint) error {
	return m.db.WithContext(ctx).Model(&entity.Memory{}).Where("id = ?", id).Update("is_archived", false).Error
}

// UpdateTags 更新记忆标签
// 嘿嘿~ 先删除旧标签再添加新标签！💖
func (m *MemoryModel) UpdateTags(ctx context.Context, memoryID uint, tags []string) error {
	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 删除旧标签
		if err := tx.Where("memory_id = ?", memoryID).Delete(&entity.MemoryTag{}).Error; err != nil {
			return err
		}
		// 添加新标签
		for _, tag := range tags {
			memoryTag := entity.MemoryTag{
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
