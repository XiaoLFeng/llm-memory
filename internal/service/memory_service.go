package service

import (
	"context"
	"errors"
	"strings"

	"github.com/XiaoLFeng/llm-memory/internal/models"
	"github.com/XiaoLFeng/llm-memory/internal/models/dto"
	"github.com/XiaoLFeng/llm-memory/internal/models/entity"
	"github.com/XiaoLFeng/llm-memory/pkg/types"
)

// MemoryService 记忆服务结构体
// 嘿嘿~ 这是处理记忆业务逻辑的服务层哦！💖
// 负责验证、处理和协调各种记忆操作~ ✨
type MemoryService struct {
	model *models.MemoryModel
}

// NewMemoryService 创建新的记忆服务实例
// 呀~ 构造函数来啦！接收一个 MemoryModel 实例~ (´∀｀)
func NewMemoryService(model *models.MemoryModel) *MemoryService {
	return &MemoryService{
		model: model,
	}
}

// resolveDefaultScope 解析默认作用域
// 嘿嘿~ group 优先，无组则 personal！💖
func resolveDefaultScope(scopeCtx *types.ScopeContext) (uint, string) {
	// 1. 如果在组内，使用 group 作用域
	if scopeCtx != nil && scopeCtx.GroupID > 0 {
		return uint(scopeCtx.GroupID), ""
	}
	// 2. 否则使用 personal 作用域（当前目录）
	if scopeCtx != nil && scopeCtx.CurrentPath != "" {
		return 0, scopeCtx.CurrentPath
	}
	// 3. 最后回退到 global
	return 0, ""
}

// parseScope 解析 scope 参数
// 呀~ 支持 personal/group/global 三种作用域过滤！✨
func parseScope(scope string, scopeCtx *types.ScopeContext) (uint, string, bool) {
	switch strings.ToLower(scope) {
	case "personal":
		if scopeCtx != nil && scopeCtx.CurrentPath != "" {
			return 0, scopeCtx.CurrentPath, false
		}
		return 0, "", false
	case "group":
		if scopeCtx != nil && scopeCtx.GroupID > 0 {
			return uint(scopeCtx.GroupID), "", false
		}
		return 0, "", false
	case "global":
		return 0, "", true // groupID=0, path="", includeGlobal=true 代表只要全局
	case "all", "":
		// all 或不指定则显示所有可见数据
		return 0, "", true
	default:
		return 0, "", true
	}
}

// CreateMemory 创建新的记忆
// 嘿嘿~ 创建记忆前会先验证数据的完整性呢！💫
// scope 参数: personal/group/global，留空则使用默认作用域（group > personal）
func (s *MemoryService) CreateMemory(ctx context.Context, input *dto.MemoryCreateDTO, scopeCtx *types.ScopeContext) (*entity.Memory, error) {
	// 验证标题不能为空
	if strings.TrimSpace(input.Title) == "" {
		return nil, errors.New("标题不能为空哦~ 📝")
	}

	// 验证内容不能为空
	if strings.TrimSpace(input.Content) == "" {
		return nil, errors.New("内容不能为空哦~ 📖")
	}

	// 默认分类
	category := strings.TrimSpace(input.Category)
	if category == "" {
		category = "默认"
	}

	// 默认优先级
	priority := input.Priority
	if priority < 1 || priority > 4 {
		priority = 1
	}

	// 解析作用域
	var groupID uint
	var path string

	scope := strings.ToLower(input.Scope)
	switch scope {
	case "personal":
		if scopeCtx != nil && scopeCtx.CurrentPath != "" {
			path = scopeCtx.CurrentPath
		}
	case "group":
		if scopeCtx != nil && scopeCtx.GroupID > 0 {
			groupID = uint(scopeCtx.GroupID)
		}
	case "global":
		// groupID 和 path 都为空即为 global
	default:
		// 默认：group 优先，然后 personal
		groupID, path = resolveDefaultScope(scopeCtx)
	}

	// 创建记忆实例
	memory := &entity.Memory{
		GroupID:  groupID,
		Path:     path,
		Title:    strings.TrimSpace(input.Title),
		Content:  strings.TrimSpace(input.Content),
		Category: category,
		Priority: priority,
	}

	// 保存到数据库
	if err := s.model.Create(ctx, memory); err != nil {
		return nil, err
	}

	// 更新标签
	if len(input.Tags) > 0 {
		if err := s.model.UpdateTags(ctx, memory.ID, input.Tags); err != nil {
			return nil, err
		}
		// 重新获取以包含标签
		memory, _ = s.model.FindByID(ctx, memory.ID)
	}

	return memory, nil
}

// UpdateMemory 更新记忆
// 呀~ 更新前会验证记忆的完整性，确保数据有效！✨
func (s *MemoryService) UpdateMemory(ctx context.Context, input *dto.MemoryUpdateDTO) error {
	// 验证ID必须大于0
	if input.ID == 0 {
		return errors.New("记忆ID必须大于 0 哦~ 🎮")
	}

	// 获取现有记忆
	memory, err := s.model.FindByID(ctx, input.ID)
	if err != nil {
		return errors.New("记忆不存在，无法更新哦~ 🔍")
	}

	// 应用更新
	if input.Title != nil {
		title := strings.TrimSpace(*input.Title)
		if title == "" {
			return errors.New("标题不能为空哦~ 📝")
		}
		memory.Title = title
	}
	if input.Content != nil {
		content := strings.TrimSpace(*input.Content)
		if content == "" {
			return errors.New("内容不能为空哦~ 📖")
		}
		memory.Content = content
	}
	if input.Category != nil {
		category := strings.TrimSpace(*input.Category)
		if category == "" {
			return errors.New("分类不能为空哦~ 🏷️")
		}
		memory.Category = category
	}
	if input.Priority != nil {
		priority := *input.Priority
		if priority < 1 || priority > 4 {
			return errors.New("优先级必须在 1-4 之间哦~ 🎮")
		}
		memory.Priority = priority
	}

	// 执行更新操作
	if err := s.model.Update(ctx, memory); err != nil {
		return err
	}

	// 更新标签（如果提供）
	if input.Tags != nil {
		if err := s.model.UpdateTags(ctx, memory.ID, *input.Tags); err != nil {
			return err
		}
	}

	return nil
}

// DeleteMemory 删除记忆
// 呀~ 删除前会验证ID和记忆是否存在！💨
func (s *MemoryService) DeleteMemory(ctx context.Context, id uint) error {
	// 验证ID必须大于0
	if id == 0 {
		return errors.New("记忆ID必须大于 0 哦~ 🎮")
	}

	// 检查记忆是否存在
	_, err := s.model.FindByID(ctx, id)
	if err != nil {
		return errors.New("记忆不存在，无法删除哦~ 💫")
	}

	// 执行删除操作
	return s.model.Delete(ctx, id)
}

// GetMemory 获取单个记忆
// 嘿嘿~ 根据ID精准查找记忆！就像寻宝一样~ 🏴‍☠️
func (s *MemoryService) GetMemory(ctx context.Context, id uint) (*entity.Memory, error) {
	// 验证ID必须大于0
	if id == 0 {
		return nil, errors.New("记忆ID必须大于 0 哦~ 🎮")
	}

	// 从模型层获取记忆
	return s.model.FindByID(ctx, id)
}

// ListMemories 列出所有记忆
// 呀~ 获取所有记忆列表，就像打开记忆相册一样！📖
func (s *MemoryService) ListMemories(ctx context.Context) ([]entity.Memory, error) {
	return s.model.FindAll(ctx)
}

// ListMemoriesByScope 根据作用域列出记忆
// 嘿嘿~ 支持 Personal/Group/Global 三层作用域过滤！💖
// scope 参数: personal/group/global/all
func (s *MemoryService) ListMemoriesByScope(ctx context.Context, scope string, scopeCtx *types.ScopeContext) ([]entity.Memory, error) {
	var groupID uint
	var path string
	var includeGlobal bool

	switch strings.ToLower(scope) {
	case "personal":
		if scopeCtx != nil && scopeCtx.CurrentPath != "" {
			path = scopeCtx.CurrentPath
		}
		includeGlobal = false
	case "group":
		if scopeCtx != nil && scopeCtx.GroupID > 0 {
			groupID = uint(scopeCtx.GroupID)
		}
		includeGlobal = false
	case "global":
		includeGlobal = true
		// groupID=0, path="" 代表只要全局数据
	case "all", "":
		// all 或不指定则显示所有可见数据
		if scopeCtx != nil {
			if scopeCtx.CurrentPath != "" {
				path = scopeCtx.CurrentPath
			}
			if scopeCtx.GroupID > 0 {
				groupID = uint(scopeCtx.GroupID)
			}
		}
		includeGlobal = true
	default:
		includeGlobal = true
	}

	return s.model.FindByScope(ctx, groupID, path, includeGlobal)
}

// ListByCategory 根据分类列出记忆
// 嗯嗯！按分类筛选记忆，让记忆管理更有条理~ 🏷️
func (s *MemoryService) ListByCategory(ctx context.Context, category string) ([]entity.Memory, error) {
	// 验证分类不能为空
	if strings.TrimSpace(category) == "" {
		return nil, errors.New("分类名称不能为空哦~ 📝")
	}

	return s.model.FindByCategory(ctx, category)
}

// SearchMemories 搜索记忆
// 嘿嘿~ 智能搜索功能！根据关键词在标题和内容中查找~ 🔍
func (s *MemoryService) SearchMemories(ctx context.Context, keyword string) ([]entity.Memory, error) {
	// 验证关键词不能为空
	if strings.TrimSpace(keyword) == "" {
		return nil, errors.New("搜索关键词不能为空哦~ 🎯")
	}

	return s.model.Search(ctx, keyword)
}

// SearchMemoriesByScope 根据作用域搜索记忆
// 在指定作用域内搜索关键词~ 🔍
func (s *MemoryService) SearchMemoriesByScope(ctx context.Context, keyword string, scope string, scopeCtx *types.ScopeContext) ([]entity.Memory, error) {
	// 验证关键词不能为空
	if strings.TrimSpace(keyword) == "" {
		return nil, errors.New("搜索关键词不能为空哦~ 🎯")
	}

	var groupID uint
	var path string
	var includeGlobal bool

	switch strings.ToLower(scope) {
	case "personal":
		if scopeCtx != nil && scopeCtx.CurrentPath != "" {
			path = scopeCtx.CurrentPath
		}
		includeGlobal = false
	case "group":
		if scopeCtx != nil && scopeCtx.GroupID > 0 {
			groupID = uint(scopeCtx.GroupID)
		}
		includeGlobal = false
	case "global":
		includeGlobal = true
	case "all", "":
		if scopeCtx != nil {
			if scopeCtx.CurrentPath != "" {
				path = scopeCtx.CurrentPath
			}
			if scopeCtx.GroupID > 0 {
				groupID = uint(scopeCtx.GroupID)
			}
		}
		includeGlobal = true
	default:
		includeGlobal = true
	}

	return s.model.SearchByScope(ctx, keyword, groupID, path, includeGlobal)
}

// ArchiveMemory 归档记忆
// 呀~ 将记忆标记为已归档状态！💼
// 归档后的记忆不会显示在常规列表中~ ✨
func (s *MemoryService) ArchiveMemory(ctx context.Context, id uint) error {
	// 验证ID必须大于0
	if id == 0 {
		return errors.New("记忆ID必须大于 0 哦~ 🎮")
	}

	// 获取记忆实例
	memory, err := s.model.FindByID(ctx, id)
	if err != nil {
		return errors.New("记忆不存在，无法归档哦~ 💫")
	}

	// 检查是否已经归档
	if memory.IsArchived {
		return errors.New("记忆已经归档过了哦~ 📦")
	}

	// 执行归档
	return s.model.Archive(ctx, id)
}

// UnarchiveMemory 取消归档记忆
func (s *MemoryService) UnarchiveMemory(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("记忆ID必须大于 0 哦~ 🎮")
	}

	memory, err := s.model.FindByID(ctx, id)
	if err != nil {
		return errors.New("记忆不存在哦~ 💫")
	}

	if !memory.IsArchived {
		return errors.New("记忆未归档哦~ 📦")
	}

	return s.model.Unarchive(ctx, id)
}

// ToMemoryResponseDTO 将 Memory entity 转换为 ResponseDTO
// 嘿嘿~ 数据转换小助手！💖
func ToMemoryResponseDTO(memory *entity.Memory, currentPath string) *dto.MemoryResponseDTO {
	if memory == nil {
		return nil
	}

	tags := make([]string, 0, len(memory.Tags))
	for _, t := range memory.Tags {
		tags = append(tags, t.Tag)
	}

	// 判断作用域
	var scope types.Scope
	if memory.Path != "" {
		scope = types.ScopePersonal
	} else if memory.GroupID > 0 {
		scope = types.ScopeGroup
	} else {
		scope = types.ScopeGlobal
	}

	return &dto.MemoryResponseDTO{
		ID:         memory.ID,
		Title:      memory.Title,
		Content:    memory.Content,
		Category:   memory.Category,
		Tags:       tags,
		Priority:   memory.Priority,
		IsArchived: memory.IsArchived,
		Scope:      string(scope),
		CreatedAt:  memory.CreatedAt,
		UpdatedAt:  memory.UpdatedAt,
	}
}

// ToMemoryListDTO 将 Memory entity 转换为 ListDTO
func ToMemoryListDTO(memory *entity.Memory) *dto.MemoryListDTO {
	if memory == nil {
		return nil
	}

	return &dto.MemoryListDTO{
		ID:         memory.ID,
		Title:      memory.Title,
		Category:   memory.Category,
		Priority:   memory.Priority,
		IsArchived: memory.IsArchived,
	}
}
