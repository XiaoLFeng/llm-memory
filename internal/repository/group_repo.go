package repository

import (
	"context"
	"errors"
	"path/filepath"

	"github.com/XiaoLFeng/llm-memory/internal/database"
	"github.com/XiaoLFeng/llm-memory/pkg/types"
	"github.com/asdine/storm/v3"
)

// GroupRepo 组仓储实现
// 嘿嘿~ 这是 Group 管理的核心仓储实现呢！📦
type GroupRepo struct {
	dbPath string
}

// NewGroupRepo 创建新的组仓储实例
// 呀~ 构造函数来啦！现在接收 dbPath 字符串~ (´∀｀)
func NewGroupRepo(dbPath string) *GroupRepo {
	return &GroupRepo{
		dbPath: dbPath,
	}
}

// Create 创建新的组
// 使用 db.Save 方法优���地保存到数据库~ 🎯
func (r *GroupRepo) Create(ctx context.Context, group *types.Group) error {
	if group == nil {
		return errors.New("group 不能为空哦~ 💫")
	}

	return database.OpenWithActionNoReturn(r.dbPath, func(db *database.DB) error {
		return db.Save(group)
	})
}

// Update 更新现有组
// 自动更新 UpdatedAt 字段，然后使用 db.Update~ ✨
func (r *GroupRepo) Update(ctx context.Context, group *types.Group) error {
	if group == nil {
		return errors.New("group 不能为空哦~ 💫")
	}

	// 自动设置更新时间
	_ = group.BeforeUpdate()

	return database.OpenWithActionNoReturn(r.dbPath, func(db *database.DB) error {
		return db.Update(group)
	})
}

// Delete 删除指定ID的组
// 同时删除相关的路径映射~ 💨
func (r *GroupRepo) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("ID 必须大于 0 哦~ 🎮")
	}

	return database.OpenWithActionNoReturn(r.dbPath, func(db *database.DB) error {
		// 先删除路径映射
		var mappings []types.GroupPathMapping
		err := db.Find("GroupID", id, &mappings)
		if err == nil {
			for _, mapping := range mappings {
				_ = db.DeleteStruct(&mapping)
			}
		}

		// 再删除组
		group := &types.Group{ID: id}
		return db.DeleteStruct(group)
	})
}

// FindByID 根据ID查找组
// 使用 db.One 方法精准查找！🏴‍☠️
func (r *GroupRepo) FindByID(ctx context.Context, id int) (*types.Group, error) {
	if id <= 0 {
		return nil, errors.New("ID 必须大于 0 哦~ 🎮")
	}

	return database.OpenWithAction(r.dbPath, func(db *database.DB) (*types.Group, error) {
		var group types.Group
		err := db.One("ID", id, &group)
		if err != nil {
			return nil, err
		}
		return &group, nil
	})
}

// FindByName 根据名称查找组
// 组名是唯一的，所以可以精准查找~ 🎯
func (r *GroupRepo) FindByName(ctx context.Context, name string) (*types.Group, error) {
	if name == "" {
		return nil, errors.New("组名称不能为空哦~ 📝")
	}

	return database.OpenWithAction(r.dbPath, func(db *database.DB) (*types.Group, error) {
		var group types.Group
		err := db.One("Name", name, &group)
		if err != nil {
			return nil, err
		}
		return &group, nil
	})
}

// FindByPath 根据路径查找所属组
// 先从路径映射表快速查找~ 🔍
func (r *GroupRepo) FindByPath(ctx context.Context, path string) (*types.Group, error) {
	if path == "" {
		return nil, errors.New("路径不能为空哦~ 📝")
	}

	// 规范化路径
	absPath, err := filepath.Abs(path)
	if err != nil {
		absPath = path
	}

	return database.OpenWithAction(r.dbPath, func(db *database.DB) (*types.Group, error) {
		// 从路径映射表查找
		var mapping types.GroupPathMapping
		err := db.One("Path", absPath, &mapping)
		if err != nil {
			if errors.Is(err, storm.ErrNotFound) {
				return nil, nil // 没找到，返回 nil
			}
			return nil, err
		}

		// 根据映射的 GroupID 查找组
		var group types.Group
		err = db.One("ID", mapping.GroupID, &group)
		if err != nil {
			return nil, err
		}
		return &group, nil
	})
}

// FindAll 查找所有组
// 使用 db.All 方法获取所有组~ 📖
func (r *GroupRepo) FindAll(ctx context.Context) ([]types.Group, error) {
	return database.OpenWithAction(r.dbPath, func(db *database.DB) ([]types.Group, error) {
		var groups []types.Group
		err := db.All(&groups)
		if err != nil {
			return nil, err
		}
		return groups, nil
	})
}

// AddPath 添加路径到组
// 同时更新组的 Paths 列表和路径映射表~ ✨
func (r *GroupRepo) AddPath(ctx context.Context, groupID int, path string) error {
	if groupID <= 0 {
		return errors.New("GroupID 必须大于 0 哦~ 🎮")
	}
	if path == "" {
		return errors.New("路径不能为空哦~ 📝")
	}

	// 规范化路径
	absPath, err := filepath.Abs(path)
	if err != nil {
		absPath = path
	}

	return database.OpenWithActionNoReturn(r.dbPath, func(db *database.DB) error {
		// 检查路径是否已被其他组占用
		var existingMapping types.GroupPathMapping
		err := db.One("Path", absPath, &existingMapping)
		if err == nil && existingMapping.GroupID != groupID {
			var existingGroup types.Group
			if err := db.One("ID", existingMapping.GroupID, &existingGroup); err == nil {
				return errors.New("该路径已属于其他组: " + existingGroup.Name)
			}
		}

		// 获取组
		var group types.Group
		err = db.One("ID", groupID, &group)
		if err != nil {
			return err
		}

		// 添加路径到组
		if !group.AddPath(absPath) {
			return errors.New("路径已存在于组中~ 📝")
		}

		// 更新组
		if err := db.Update(&group); err != nil {
			return err
		}

		// 创建路径映射
		mapping := types.NewGroupPathMapping(absPath, groupID)
		return db.Save(mapping)
	})
}

// RemovePath 从组中移除路径
// 同时删除路径映射~ 💨
func (r *GroupRepo) RemovePath(ctx context.Context, groupID int, path string) error {
	if groupID <= 0 {
		return errors.New("GroupID 必须大于 0 哦~ 🎮")
	}
	if path == "" {
		return errors.New("路径不能为空哦~ 📝")
	}

	// 规范化路径
	absPath, err := filepath.Abs(path)
	if err != nil {
		absPath = path
	}

	return database.OpenWithActionNoReturn(r.dbPath, func(db *database.DB) error {
		// 获取组
		var group types.Group
		err := db.One("ID", groupID, &group)
		if err != nil {
			return err
		}

		// 从组中移除路径
		if !group.RemovePath(absPath) {
			return errors.New("路径不存在于组中~ 📝")
		}

		// 更新组
		if err := db.Update(&group); err != nil {
			return err
		}

		// 删除路径映射
		var mapping types.GroupPathMapping
		err = db.One("Path", absPath, &mapping)
		if err == nil {
			_ = db.DeleteStruct(&mapping)
		}

		return nil
	})
}
