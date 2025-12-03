package service

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/XiaoLFeng/llm-memory/internal/models"
	"github.com/XiaoLFeng/llm-memory/internal/models/dto"
	"github.com/XiaoLFeng/llm-memory/internal/models/entity"
	"github.com/XiaoLFeng/llm-memory/pkg/types"
)

// GroupService 组服务层
// 嘿嘿~ 用于管理 Group 的业务逻辑！📦
type GroupService struct {
	model *models.GroupModel
}

// NewGroupService 创建新的组服务实例
// 呀~ 构造函数来啦！(´∀｀)
func NewGroupService(model *models.GroupModel) *GroupService {
	return &GroupService{
		model: model,
	}
}

// CreateGroup 创建新组
// 嘿嘿~ 创建一个新的组来管理多个路径！💖
func (s *GroupService) CreateGroup(ctx context.Context, name, description string) (*entity.Group, error) {
	// 验证组名
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("组名称不能为空哦~ 📝")
	}

	// 检查组名是否已存在
	existing, _ := s.model.FindByName(ctx, name)
	if existing != nil {
		return nil, errors.New("组名称已存在哦~ 💫")
	}

	// 创建组
	group := &entity.Group{
		Name:        name,
		Description: strings.TrimSpace(description),
	}
	if err := s.model.Create(ctx, group); err != nil {
		return nil, err
	}

	return group, nil
}

// UpdateGroup 更新组信息
func (s *GroupService) UpdateGroup(ctx context.Context, id uint, name, description *string) error {
	if id == 0 {
		return errors.New("组ID必须大于 0 哦~ 🎮")
	}

	// 获取现有组
	group, err := s.model.FindByID(ctx, id)
	if err != nil {
		return errors.New("组不存在哦~ 💫")
	}

	// 应用更新
	if name != nil {
		trimmedName := strings.TrimSpace(*name)
		if trimmedName == "" {
			return errors.New("组名称不能为空哦~ 📝")
		}
		// 检查新名称是否被其他组使用
		existing, _ := s.model.FindByName(ctx, trimmedName)
		if existing != nil && existing.ID != id {
			return errors.New("组名称已被使用哦~ 💫")
		}
		group.Name = trimmedName
	}
	if description != nil {
		group.Description = strings.TrimSpace(*description)
	}

	return s.model.Update(ctx, group)
}

// DeleteGroup 删除组
// 注意：这不会删除关联的数据，只是解除路径关联
func (s *GroupService) DeleteGroup(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("组ID必须大于 0 哦~ 🎮")
	}
	return s.model.Delete(ctx, id)
}

// GetGroup 获取组详情
func (s *GroupService) GetGroup(ctx context.Context, id uint) (*entity.Group, error) {
	if id == 0 {
		return nil, errors.New("组ID必须大于 0 哦~ 🎮")
	}
	return s.model.FindByID(ctx, id)
}

// GetGroupByName 根据名称获取组
func (s *GroupService) GetGroupByName(ctx context.Context, name string) (*entity.Group, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("组名称不能为空哦~ 📝")
	}
	return s.model.FindByName(ctx, name)
}

// ListGroups 列出所有组
func (s *GroupService) ListGroups(ctx context.Context) ([]entity.Group, error) {
	return s.model.FindAll(ctx)
}

// AddCurrentPath 将当前工作目录添加到组
// 这是最常用的添加路径方法~ ✨
func (s *GroupService) AddCurrentPath(ctx context.Context, groupID uint) error {
	pwd, err := os.Getwd()
	if err != nil {
		return errors.New("无法获取当前工作目录: " + err.Error())
	}
	return s.AddPath(ctx, groupID, pwd)
}

// AddPath 添加指定路径到组
func (s *GroupService) AddPath(ctx context.Context, groupID uint, path string) error {
	if groupID == 0 {
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

	return s.model.AddPath(ctx, groupID, absPath)
}

// AddPathByName 根据组名添加路径
// 嘿嘿~ 通过组名添加路径更方便！💖
func (s *GroupService) AddPathByName(ctx context.Context, groupName string, path string) error {
	groupName = strings.TrimSpace(groupName)
	if groupName == "" {
		return errors.New("组名称不能为空哦~ 📝")
	}

	// 查找组
	group, err := s.model.FindByName(ctx, groupName)
	if err != nil {
		return errors.New("组不存在: " + groupName)
	}

	// 如果路径为空，使用当前目录
	if path == "" {
		pwd, err := os.Getwd()
		if err != nil {
			return errors.New("无法获取当前工作目录: " + err.Error())
		}
		path = pwd
	}

	return s.AddPath(ctx, group.ID, path)
}

// RemovePath 从组中移除路径
func (s *GroupService) RemovePath(ctx context.Context, groupID uint, path string) error {
	if groupID == 0 {
		return errors.New("组ID必须大于 0 哦~ 🎮")
	}

	// 规范化路径
	absPath, err := filepath.Abs(path)
	if err != nil {
		absPath = path
	}

	return s.model.RemovePath(ctx, groupID, absPath)
}

// RemovePathByName 根据组名移除路径
func (s *GroupService) RemovePathByName(ctx context.Context, groupName string, path string) error {
	groupName = strings.TrimSpace(groupName)
	if groupName == "" {
		return errors.New("组名称不能为空哦~ 📝")
	}

	// 查找组
	group, err := s.model.FindByName(ctx, groupName)
	if err != nil {
		return errors.New("组不存在: " + groupName)
	}

	return s.RemovePath(ctx, group.ID, path)
}

// GetGroupByPath 根据路径获取所属组
func (s *GroupService) GetGroupByPath(ctx context.Context, path string) (*entity.Group, error) {
	if path == "" {
		return nil, errors.New("路径不能为空哦~ 📝")
	}

	// 规范化路径
	absPath, err := filepath.Abs(path)
	if err != nil {
		absPath = path
	}

	return s.model.FindByPath(ctx, absPath)
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
	group, err := s.model.FindByPath(ctx, absPath)
	if err != nil {
		// 查找失败，使用默认作用域
		return scope, nil
	}

	if group != nil {
		// 找到了组，设置组信息
		scope.GroupID = int(group.ID)
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

// GetScopeInfo 获取当前作用域信息 DTO
// 呀~ 返回作用域的详细信息！✨
func (s *GroupService) GetScopeInfo(ctx context.Context) (*dto.ScopeInfoDTO, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return &dto.ScopeInfoDTO{
			CurrentPath: "",
			IsInGroup:   false,
		}, nil
	}

	// 规范化路径
	absPath, err := filepath.Abs(pwd)
	if err != nil {
		absPath = pwd
	}

	info := &dto.ScopeInfoDTO{
		CurrentPath: absPath,
		IsInGroup:   false,
	}

	// 查找所属组
	group, err := s.model.FindByPath(ctx, absPath)
	if err == nil && group != nil {
		info.GroupID = group.ID
		info.GroupName = group.Name
		info.IsInGroup = true
	}

	return info, nil
}

// ToGroupResponseDTO 将 Group entity 转换为 ResponseDTO
// 嘿嘿~ 数据转换小助手！💖
func ToGroupResponseDTO(group *entity.Group) *dto.GroupResponseDTO {
	if group == nil {
		return nil
	}

	paths := make([]string, 0, len(group.Paths))
	for _, p := range group.Paths {
		paths = append(paths, p.Path)
	}

	return &dto.GroupResponseDTO{
		ID:          group.ID,
		Name:        group.Name,
		Description: group.Description,
		Paths:       paths,
		PathCount:   len(paths),
		CreatedAt:   group.CreatedAt,
		UpdatedAt:   group.UpdatedAt,
	}
}

// ToGroupListDTO 将 Group entity 转换为 ListDTO
func ToGroupListDTO(group *entity.Group) *dto.GroupListDTO {
	if group == nil {
		return nil
	}

	return &dto.GroupListDTO{
		ID:          group.ID,
		Name:        group.Name,
		Description: group.Description,
		PathCount:   len(group.Paths),
	}
}
