package service

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/XiaoLFeng/llm-memory/internal/database"
	"github.com/XiaoLFeng/llm-memory/pkg/types"
)

// GroupService 组服务层
// 嘿嘿~ 用于管理 Group 的业务逻辑！📦
type GroupService struct {
	repo database.GroupRepository
}

// NewGroupService 创建新的组服务实例
// 呀~ 构造函数来啦！(´∀｀)
func NewGroupService(repo database.GroupRepository) *GroupService {
	return &GroupService{
		repo: repo,
	}
}

// CreateGroup 创建新组
// 嘿嘿~ 创建一个新的组来管理多个路径！💖
func (s *GroupService) CreateGroup(ctx context.Context, name, description string) (*types.Group, error) {
	// 验证组名
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("组名称不能为空哦~ 📝")
	}

	// 检查组名是否已存在
	existing, _ := s.repo.FindByName(ctx, name)
	if existing != nil {
		return nil, errors.New("组名称已存在哦~ 💫")
	}

	// 创建组
	group := types.NewGroup(name, strings.TrimSpace(description))
	if err := s.repo.Create(ctx, group); err != nil {
		return nil, err
	}

	return group, nil
}

// UpdateGroup 更新组信息
func (s *GroupService) UpdateGroup(ctx context.Context, group *types.Group) error {
	if group == nil {
		return errors.New("组对象不能为空哦~ 💫")
	}
	if group.ID <= 0 {
		return errors.New("组ID必须大于 0 哦~ 🎮")
	}
	return s.repo.Update(ctx, group)
}

// DeleteGroup 删除组
// 注意：这不会删除关联的数据，只是解除路径关联
func (s *GroupService) DeleteGroup(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("组ID必须大于 0 哦~ 🎮")
	}
	return s.repo.Delete(ctx, id)
}

// GetGroup 获取组详情
func (s *GroupService) GetGroup(ctx context.Context, id int) (*types.Group, error) {
	if id <= 0 {
		return nil, errors.New("组ID必须大于 0 哦~ 🎮")
	}
	return s.repo.FindByID(ctx, id)
}

// GetGroupByName 根据名称获取组
func (s *GroupService) GetGroupByName(ctx context.Context, name string) (*types.Group, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("组名称不能为空哦~ 📝")
	}
	return s.repo.FindByName(ctx, name)
}

// ListGroups 列出所有组
func (s *GroupService) ListGroups(ctx context.Context) ([]types.Group, error) {
	return s.repo.FindAll(ctx)
}

// AddCurrentPath 将当前工作目录添加到组
// 这是最常用的添加路径方法~ ✨
func (s *GroupService) AddCurrentPath(ctx context.Context, groupID int) error {
	pwd, err := os.Getwd()
	if err != nil {
		return errors.New("无法获取当前工作目录: " + err.Error())
	}
	return s.AddPath(ctx, groupID, pwd)
}

// AddPath 添加指定路径到组
func (s *GroupService) AddPath(ctx context.Context, groupID int, path string) error {
	if groupID <= 0 {
		return errors.New("组ID必须大于 0 哦~ 🎮")
	}

	// 规范化路径
	absPath, err := filepath.Abs(path)
	if err != nil {
		return errors.New("无效的路径: " + err.Error())
	}

	// 检查路径是否存在
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return errors.New("路径不存在: " + absPath)
	}

	return s.repo.AddPath(ctx, groupID, absPath)
}

// RemovePath 从组中移除路径
func (s *GroupService) RemovePath(ctx context.Context, groupID int, path string) error {
	if groupID <= 0 {
		return errors.New("组ID必须大于 0 哦~ 🎮")
	}

	// 规范化路径
	absPath, err := filepath.Abs(path)
	if err != nil {
		absPath = path
	}

	return s.repo.RemovePath(ctx, groupID, absPath)
}

// GetGroupByPath 根据路径获取所属组
func (s *GroupService) GetGroupByPath(ctx context.Context, path string) (*types.Group, error) {
	if path == "" {
		return nil, errors.New("路径不能为空哦~ 📝")
	}

	// 规范化路径
	absPath, err := filepath.Abs(path)
	if err != nil {
		absPath = path
	}

	return s.repo.FindByPath(ctx, absPath)
}

// ResolveScope 解析当前作用域
// 这是核心方法！根据 pwd 确定当前的 ScopeContext~ 🎯
func (s *GroupService) ResolveScope(ctx context.Context, pwd string) (*types.ScopeContext, error) {
	// 规范化路径
	absPath, err := filepath.Abs(pwd)
	if err != nil {
		absPath = pwd
	}

	// 创建默认的作用域上下文
	scope := types.NewScopeContext(absPath)

	// 查找路径所属的组
	group, err := s.repo.FindByPath(ctx, absPath)
	if err != nil {
		// 查找失败，使用默认作用域
		return scope, nil
	}

	if group != nil {
		// 找到了组，设置组信息
		scope.GroupID = group.ID
		scope.GroupName = group.Name
	}

	return scope, nil
}

// GetCurrentScope 获取当前工作目录的作用域
// 便捷方法，自动获取 pwd~ ✨
func (s *GroupService) GetCurrentScope(ctx context.Context) (*types.ScopeContext, error) {
	pwd, err := os.Getwd()
	if err != nil {
		// 无法获取 pwd，返回只有 Global 的作用域
		return types.NewGlobalOnlyScope(), nil
	}
	return s.ResolveScope(ctx, pwd)
}
