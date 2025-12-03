package models

import (
	"context"
	"path/filepath"
	"time"

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

// Delete 删除组（硬删除）
func (m *GroupModel) Delete(ctx context.Context, id int64) error {
	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 先删除路径映射
		if err := tx.Where("group_id = ?", id).Unscoped().Delete(&entity.GroupPath{}).Error; err != nil {
			return err
		}
		// 硬删除组
		return tx.Unscoped().Delete(&entity.Group{}, id).Error
	})
}

// FindByID 根据 ID 查找组
func (m *GroupModel) FindByID(ctx context.Context, id int64) (*entity.Group, error) {
	var group entity.Group
	err := m.db.WithContext(ctx).Preload("Paths.PersonalPath").First(&group, id).Error
	if err != nil {
		return nil, err
	}
	return &group, nil
}

// FindByName 根据名称查找组
func (m *GroupModel) FindByName(ctx context.Context, name string) (*entity.Group, error) {
	var group entity.Group
	err := m.db.WithContext(ctx).Preload("Paths.PersonalPath").Where("name = ?", name).First(&group).Error
	if err != nil {
		return nil, err
	}
	return &group, nil
}

// FindByPath 根据路径查找所属组
// 纯关联模式：通过 PersonalPath -> GroupPath 查找
func (m *GroupModel) FindByPath(ctx context.Context, path string) (*entity.Group, error) {
	// 规范化路径
	absPath, err := filepath.Abs(path)
	if err != nil {
		absPath = path
	}

	// 先找 PersonalPath
	var personalPath entity.PersonalPath
	err = m.db.WithContext(ctx).Where("path = ?", absPath).First(&personalPath).Error
	if err != nil {
		return nil, err
	}

	// 再找 GroupPath 关联
	var groupPath entity.GroupPath
	err = m.db.WithContext(ctx).Where("personal_path_id = ?", personalPath.ID).First(&groupPath).Error
	if err != nil {
		return nil, err
	}

	return m.FindByID(ctx, groupPath.GroupID)
}

// FindAll 查找所有组
func (m *GroupModel) FindAll(ctx context.Context) ([]entity.Group, error) {
	var groups []entity.Group
	err := m.db.WithContext(ctx).Preload("Paths.PersonalPath").Order("created_at DESC").Find(&groups).Error
	return groups, err
}

// AddPath 添加路径到组
// 纯关联模式：先确保 PersonalPath 存在，再创建 GroupPath 关联
func (m *GroupModel) AddPath(ctx context.Context, groupID int64, path string) error {
	// 规范化路径
	absPath, err := filepath.Abs(path)
	if err != nil {
		absPath = path
	}

	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. 确保 PersonalPath 存在
		var personalPath entity.PersonalPath
		err := tx.Where("path = ?", absPath).First(&personalPath).Error
		if err == gorm.ErrRecordNotFound {
			// 不存在，创建新的 PersonalPath
			personalPath = entity.PersonalPath{
				ID:        database.GenerateID(),
				Path:      absPath,
				LastVisit: time.Now(),
			}
			if err := tx.Create(&personalPath).Error; err != nil {
				return err
			}
		} else if err != nil {
			return err
		}

		// 2. 检查 GroupPath 关联是否已存在
		var existingGroupPath entity.GroupPath
		err = tx.Where("personal_path_id = ?", personalPath.ID).First(&existingGroupPath).Error
		if err == nil {
			// 关联已存在
			if existingGroupPath.GroupID == groupID {
				return nil // 已经在当前组
			}
			return gorm.ErrDuplicatedKey // 被其他组占用
		}
		if err != gorm.ErrRecordNotFound {
			return err
		}

		// 3. 创建 GroupPath 关联
		groupPath := entity.GroupPath{
			ID:             database.GenerateID(),
			GroupID:        groupID,
			PersonalPathID: personalPath.ID,
		}
		return tx.Create(&groupPath).Error
	})
}

// RemovePath 从组移除路径
// 纯关联模式：通过 PersonalPath 查找关联并删除
func (m *GroupModel) RemovePath(ctx context.Context, groupID int64, path string) error {
	// 规范化路径
	absPath, err := filepath.Abs(path)
	if err != nil {
		absPath = path
	}

	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 先找 PersonalPath
		var personalPath entity.PersonalPath
		if err := tx.Where("path = ?", absPath).First(&personalPath).Error; err != nil {
			return err
		}

		// 删除 GroupPath 关联
		return tx.Where("group_id = ? AND personal_path_id = ?", groupID, personalPath.ID).
			Delete(&entity.GroupPath{}).Error
	})
}

// PathExists 检查路径是否已被任何组占用
// 纯关联模式：检查 GroupPath 关联
func (m *GroupModel) PathExists(ctx context.Context, path string) (bool, error) {
	// 规范化路径
	absPath, err := filepath.Abs(path)
	if err != nil {
		absPath = path
	}

	// 先找 PersonalPath
	var personalPath entity.PersonalPath
	err = m.db.WithContext(ctx).Where("path = ?", absPath).First(&personalPath).Error
	if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	// 检查 GroupPath 关联
	var count int64
	err = m.db.WithContext(ctx).Model(&entity.GroupPath{}).Where("personal_path_id = ?", personalPath.ID).Count(&count).Error
	return count > 0, err
}

// GetGroupIDByPath 获取路径所属的组 ID
// 纯关联模式：通过 PersonalPath -> GroupPath 查找
func (m *GroupModel) GetGroupIDByPath(ctx context.Context, path string) (int64, error) {
	// 规范化路径
	absPath, err := filepath.Abs(path)
	if err != nil {
		absPath = path
	}

	// 先找 PersonalPath
	var personalPath entity.PersonalPath
	err = m.db.WithContext(ctx).Where("path = ?", absPath).First(&personalPath).Error
	if err != nil {
		return 0, err
	}

	// 再找 GroupPath 关联
	var groupPath entity.GroupPath
	err = m.db.WithContext(ctx).Where("personal_path_id = ?", personalPath.ID).First(&groupPath).Error
	if err != nil {
		return 0, err
	}
	return groupPath.GroupID, nil
}

// GetPathIDByPath 根据路径字符串获取 PersonalPath ID
// 嘿嘿~ 纯关联模式下的辅助方法！💖
func (m *GroupModel) GetPathIDByPath(ctx context.Context, path string) (int64, error) {
	// 规范化路径
	absPath, err := filepath.Abs(path)
	if err != nil {
		absPath = path
	}

	var personalPath entity.PersonalPath
	err = m.db.WithContext(ctx).Where("path = ?", absPath).First(&personalPath).Error
	if err != nil {
		return 0, err
	}
	return personalPath.ID, nil
}

// GetPathIDsByGroupID 获取组下所有路径 ID
// 嘿嘿~ 用于 Scope 查询时获取组内所有路径！💖
func (m *GroupModel) GetPathIDsByGroupID(ctx context.Context, groupID int64) ([]int64, error) {
	var groupPaths []entity.GroupPath
	err := m.db.WithContext(ctx).Where("group_id = ?", groupID).Find(&groupPaths).Error
	if err != nil {
		return nil, err
	}

	pathIDs := make([]int64, len(groupPaths))
	for i, gp := range groupPaths {
		pathIDs[i] = gp.PersonalPathID
	}
	return pathIDs, nil
}

// GetPathStringsByGroupID 获取组下所有路径字符串
// 嘿嘿~ 用于展示时获取实际路径！💖
func (m *GroupModel) GetPathStringsByGroupID(ctx context.Context, groupID int64) ([]string, error) {
	pathIDs, err := m.GetPathIDsByGroupID(ctx, groupID)
	if err != nil {
		return nil, err
	}

	if len(pathIDs) == 0 {
		return []string{}, nil
	}

	var personalPaths []entity.PersonalPath
	err = m.db.WithContext(ctx).Where("id IN ?", pathIDs).Find(&personalPaths).Error
	if err != nil {
		return nil, err
	}

	paths := make([]string, len(personalPaths))
	for i, pp := range personalPaths {
		paths[i] = pp.Path
	}
	return paths, nil
}

// Count 获取组总数
func (m *GroupModel) Count(ctx context.Context) (int64, error) {
	var count int64
	err := m.db.WithContext(ctx).Model(&entity.Group{}).Count(&count).Error
	return count, err
}
