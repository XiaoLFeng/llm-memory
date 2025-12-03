package models

import (
	"context"
	"path/filepath"

	"github.com/XiaoLFeng/llm-memory/internal/database"
	"github.com/XiaoLFeng/llm-memory/internal/models/entity"
	"gorm.io/gorm"
)

// GroupModel 组数据访问层
type GroupModel struct {
	db *gorm.DB
}

// NewGroupModel 创建 GroupModel 实例
func NewGroupModel(db *gorm.DB) *GroupModel {
	return &GroupModel{db: db}
}

// Create 创建组
func (m *GroupModel) Create(ctx context.Context, group *entity.Group) error {
	group.ID = database.GenerateID()
	return m.db.WithContext(ctx).Create(group).Error
}

// Update 更新组
func (m *GroupModel) Update(ctx context.Context, group *entity.Group) error {
	return m.db.WithContext(ctx).Save(group).Error
}

// Delete 删除组（软删除）
func (m *GroupModel) Delete(ctx context.Context, id int64) error {
	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 先删除路径映射
		if err := tx.Where("group_id = ?", id).Delete(&entity.GroupPath{}).Error; err != nil {
			return err
		}
		// 再删除组
		return tx.Delete(&entity.Group{}, id).Error
	})
}

// FindByID 根据 ID 查找组
func (m *GroupModel) FindByID(ctx context.Context, id int64) (*entity.Group, error) {
	var group entity.Group
	err := m.db.WithContext(ctx).Preload("Paths").First(&group, id).Error
	if err != nil {
		return nil, err
	}
	return &group, nil
}

// FindByName 根据名称查找组
func (m *GroupModel) FindByName(ctx context.Context, name string) (*entity.Group, error) {
	var group entity.Group
	err := m.db.WithContext(ctx).Preload("Paths").Where("name = ?", name).First(&group).Error
	if err != nil {
		return nil, err
	}
	return &group, nil
}

// FindByPath 根据路径查找所属组
func (m *GroupModel) FindByPath(ctx context.Context, path string) (*entity.Group, error) {
	// 规范化路径
	absPath, err := filepath.Abs(path)
	if err != nil {
		absPath = path
	}

	var groupPath entity.GroupPath
	err = m.db.WithContext(ctx).Where("path = ?", absPath).First(&groupPath).Error
	if err != nil {
		return nil, err
	}

	return m.FindByID(ctx, groupPath.GroupID)
}

// FindAll 查找所有组
func (m *GroupModel) FindAll(ctx context.Context) ([]entity.Group, error) {
	var groups []entity.Group
	err := m.db.WithContext(ctx).Preload("Paths").Order("created_at DESC").Find(&groups).Error
	return groups, err
}

// AddPath 添加路径到组
func (m *GroupModel) AddPath(ctx context.Context, groupID int64, path string) error {
	// 规范化路径
	absPath, err := filepath.Abs(path)
	if err != nil {
		absPath = path
	}

	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 检查路径是否已存在
		var existing entity.GroupPath
		err := tx.Where("path = ?", absPath).First(&existing).Error
		if err == nil {
			// 路径已存在
			if existing.GroupID == groupID {
				return nil // 已经在当前组，无需操作
			}
			return gorm.ErrDuplicatedKey // 被其他组占用
		}
		if err != gorm.ErrRecordNotFound {
			return err
		}

		// 创建路径映射
		groupPath := entity.GroupPath{
			ID:      database.GenerateID(),
			GroupID: groupID,
			Path:    absPath,
		}
		return tx.Create(&groupPath).Error
	})
}

// RemovePath 从组移除路径
func (m *GroupModel) RemovePath(ctx context.Context, groupID int64, path string) error {
	// 规范化路径
	absPath, err := filepath.Abs(path)
	if err != nil {
		absPath = path
	}

	return m.db.WithContext(ctx).
		Where("group_id = ? AND path = ?", groupID, absPath).
		Delete(&entity.GroupPath{}).Error
}

// PathExists 检查路径是否已被任何组占用
func (m *GroupModel) PathExists(ctx context.Context, path string) (bool, error) {
	// 规范化路径
	absPath, err := filepath.Abs(path)
	if err != nil {
		absPath = path
	}

	var count int64
	err = m.db.WithContext(ctx).Model(&entity.GroupPath{}).Where("path = ?", absPath).Count(&count).Error
	return count > 0, err
}

// GetGroupIDByPath 获取路径所属的组 ID
func (m *GroupModel) GetGroupIDByPath(ctx context.Context, path string) (int64, error) {
	// 规范化路径
	absPath, err := filepath.Abs(path)
	if err != nil {
		absPath = path
	}

	var groupPath entity.GroupPath
	err = m.db.WithContext(ctx).Where("path = ?", absPath).First(&groupPath).Error
	if err != nil {
		return 0, err
	}
	return groupPath.GroupID, nil
}

// Count 获取组总数
func (m *GroupModel) Count(ctx context.Context) (int64, error) {
	var count int64
	err := m.db.WithContext(ctx).Model(&entity.Group{}).Count(&count).Error
	return count, err
}
