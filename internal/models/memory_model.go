package models

import (
	"context"

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

// FindByCode 根据 code 查找记忆（排除已归档）
func (m *MemoryModel) FindByCode(ctx context.Context, code string) (*entity.Memory, error) {
	var memory entity.Memory
	err := m.db.WithContext(ctx).
		Preload("Tags").
		Where("code = ? AND is_archived = ?", code, false).
		First(&memory).Error
	if err != nil {
		return nil, err
	}
	return &memory, nil
}

// ExistsCode 检查 code 是否已存在（用于创建/更新时校验唯一性）
func (m *MemoryModel) ExistsCode(ctx context.Context, code string, excludeID int64) (bool, error) {
	var count int64
	query := m.db.WithContext(ctx).Model(&entity.Memory{}).Where("code = ?", code)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

// FindAll 查找所有记忆
func (m *MemoryModel) FindAll(ctx context.Context) ([]entity.Memory, error) {
	return m.FindByFilter(ctx, DefaultVisibilityFilter())
}

// FindByCategory 根据分类查找记忆
func (m *MemoryModel) FindByCategory(ctx context.Context, category string) ([]entity.Memory, error) {
	filter := DefaultVisibilityFilter()
	var memories []entity.Memory
	err := applyVisibilityFilter(m.db.WithContext(ctx).Preload("Tags"), filter).
		Where("category = ? AND is_archived = ?", category, false).
		Order("created_at DESC").
		Find(&memories).Error
	return memories, err
}

// FindByScope 兼容旧接口：根据 PathID / GroupPathIDs 过滤
func (m *MemoryModel) FindByScope(ctx context.Context, pathID int64, groupPathIDs []int64, includeGlobal bool) ([]entity.Memory, error) {
	filter := VisibilityFilter{
		IncludeGlobal:    includeGlobal,
		IncludeNonGlobal: true,
		PathIDs:          MergePathIDs(pathID, groupPathIDs),
	}
	return m.FindByFilter(ctx, filter)
}

// FindByFilter 根据统一过滤器查询记忆
func (m *MemoryModel) FindByFilter(ctx context.Context, filter VisibilityFilter) ([]entity.Memory, error) {
	var memories []entity.Memory
	err := applyVisibilityFilter(m.db.WithContext(ctx).Preload("Tags"), filter).
		Where("is_archived = ?", false).
		Order("created_at DESC").
		Find(&memories).Error
	return memories, err
}

// Search 搜索记忆（在标题和内容中搜索）
func (m *MemoryModel) Search(ctx context.Context, keyword string) ([]entity.Memory, error) {
	filter := DefaultVisibilityFilter()
	return m.SearchByFilter(ctx, keyword, filter)
}

// SearchByScope 在指定作用域内搜索记忆（兼容旧接口）
func (m *MemoryModel) SearchByScope(ctx context.Context, keyword string, pathID int64, groupPathIDs []int64, includeGlobal bool) ([]entity.Memory, error) {
	filter := VisibilityFilter{
		IncludeGlobal:    includeGlobal,
		IncludeNonGlobal: true,
		PathIDs:          MergePathIDs(pathID, groupPathIDs),
	}
	return m.SearchByFilter(ctx, keyword, filter)
}

// SearchByFilter 统一过滤器搜索
func (m *MemoryModel) SearchByFilter(ctx context.Context, keyword string, filter VisibilityFilter) ([]entity.Memory, error) {
	var memories []entity.Memory
	pattern := "%" + keyword + "%"
	err := applyVisibilityFilter(
		m.db.WithContext(ctx).Preload("Tags").
			Where("title LIKE ? OR content LIKE ?", pattern, pattern),
		filter,
	).Where("is_archived = ?", false).
		Order("created_at DESC").Find(&memories).Error
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
