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
// 负责验证、处理和协调各种记忆操作
type MemoryService struct {
	model *models.MemoryModel
}

// NewMemoryService 创建新的记忆服务实例
func NewMemoryService(model *models.MemoryModel) *MemoryService {
	return &MemoryService{
		model: model,
	}
}

// resolveDefaultScope 解析默认作用域
// 纯关联模式：返回 PathID（0 表示 Global）
func resolveDefaultScope(scopeCtx *types.ScopeContext) int64 {
	// 1. 如果在组内，使用 personal 作用域（当前目录）
	if scopeCtx != nil && scopeCtx.PathID > 0 {
		return scopeCtx.PathID
	}
	// 2. 回退到 global
	return 0
}

// parseScope 解析 scope 参数
// 纯关联模式：返回 pathID, groupPathIDs, includeGlobal
func parseScope(scope string, scopeCtx *types.ScopeContext) (int64, []int64, bool) {
	switch strings.ToLower(scope) {
	case "personal":
		if scopeCtx != nil && scopeCtx.PathID > 0 {
			return scopeCtx.PathID, nil, false
		}
		return 0, nil, false
	case "group":
		if scopeCtx != nil && len(scopeCtx.GroupPathIDs) > 0 {
			return 0, scopeCtx.GroupPathIDs, false
		}
		return 0, nil, false
	case "global":
		return 0, nil, true
	case "all", "":
		// all 或不指定则显示所有可见数据
		var pathID int64
		var groupPathIDs []int64
		if scopeCtx != nil {
			pathID = scopeCtx.PathID
			groupPathIDs = scopeCtx.GroupPathIDs
		}
		return pathID, groupPathIDs, true
	default:
		return 0, nil, true
	}
}

// CreateMemory 创建新的记忆
// scope 参数: personal/group/global，留空则使用默认作用域
// 纯关联模式：数据存储时只使用 PathID
func (s *MemoryService) CreateMemory(ctx context.Context, input *dto.MemoryCreateDTO, scopeCtx *types.ScopeContext) (*entity.Memory, error) {
	// 验证标题不能为空
	if strings.TrimSpace(input.Title) == "" {
		return nil, errors.New("标题不能为空")
	}

	// 验证内容不能为空
	if strings.TrimSpace(input.Content) == "" {
		return nil, errors.New("内容不能为空")
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

	// 解析作用域 -> PathID
	var pathID int64

	scope := strings.ToLower(input.Scope)
	switch scope {
	case "personal":
		if scopeCtx != nil && scopeCtx.PathID > 0 {
			pathID = scopeCtx.PathID
		}
	case "group":
		// group 作用域下，仍然保存到当前路径（数据属于当前路径，只是查询时会关联组）
		if scopeCtx != nil && scopeCtx.PathID > 0 {
			pathID = scopeCtx.PathID
		}
	case "global":
		// pathID = 0 即为 global
		pathID = 0
	default:
		// 默认使用当前路径
		pathID = resolveDefaultScope(scopeCtx)
	}

	// 创建记忆实例
	memory := &entity.Memory{
		PathID:   pathID,
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
func (s *MemoryService) UpdateMemory(ctx context.Context, input *dto.MemoryUpdateDTO) error {
	// 验证ID必须大于0
	if input.ID == 0 {
		return errors.New("记忆ID必须大于 0")
	}

	// 获取现有记忆
	memory, err := s.model.FindByID(ctx, input.ID)
	if err != nil {
		return errors.New("记忆不存在，无法更新")
	}

	// 应用更新
	if input.Title != nil {
		title := strings.TrimSpace(*input.Title)
		if title == "" {
			return errors.New("标题不能为空")
		}
		memory.Title = title
	}
	if input.Content != nil {
		content := strings.TrimSpace(*input.Content)
		if content == "" {
			return errors.New("内容不能为空")
		}
		memory.Content = content
	}
	if input.Category != nil {
		category := strings.TrimSpace(*input.Category)
		if category == "" {
			return errors.New("分类不能为空")
		}
		memory.Category = category
	}
	if input.Priority != nil {
		priority := *input.Priority
		if priority < 1 || priority > 4 {
			return errors.New("优先级必须在 1-4 之间")
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
func (s *MemoryService) DeleteMemory(ctx context.Context, id int64) error {
	// 验证ID必须大于0
	if id == 0 {
		return errors.New("记忆ID必须大于 0")
	}

	// 检查记忆是否存在
	_, err := s.model.FindByID(ctx, id)
	if err != nil {
		return errors.New("记忆不存在，无法删除")
	}

	// 执行删除操作
	return s.model.Delete(ctx, id)
}

// GetMemory 获取单个记忆
func (s *MemoryService) GetMemory(ctx context.Context, id int64) (*entity.Memory, error) {
	// 验证ID必须大于0
	if id == 0 {
		return nil, errors.New("记忆ID必须大于 0")
	}

	// 从模型层获取记忆
	return s.model.FindByID(ctx, id)
}

// ListMemories 列出所有记忆
func (s *MemoryService) ListMemories(ctx context.Context) ([]entity.Memory, error) {
	return s.model.FindAll(ctx)
}

// ListMemoriesByScope 根据作用域列出记忆
// scope 参数: personal/group/global/all
// 纯关联模式：使用 PathID 和 GroupPathIDs 进行查询
func (s *MemoryService) ListMemoriesByScope(ctx context.Context, scope string, scopeCtx *types.ScopeContext) ([]entity.Memory, error) {
	pathID, groupPathIDs, includeGlobal := parseScope(scope, scopeCtx)
	return s.model.FindByScope(ctx, pathID, groupPathIDs, includeGlobal)
}

// ListByCategory 根据分类列出记忆
func (s *MemoryService) ListByCategory(ctx context.Context, category string) ([]entity.Memory, error) {
	// 验证分类不能为空
	if strings.TrimSpace(category) == "" {
		return nil, errors.New("分类名称不能为空")
	}

	return s.model.FindByCategory(ctx, category)
}

// SearchMemories 搜索记忆
func (s *MemoryService) SearchMemories(ctx context.Context, keyword string) ([]entity.Memory, error) {
	// 验证关键词不能为空
	if strings.TrimSpace(keyword) == "" {
		return nil, errors.New("搜索关键词不能为空")
	}

	return s.model.Search(ctx, keyword)
}

// SearchMemoriesByScope 根据作用域搜索记忆
// 纯关联模式：使用 PathID 和 GroupPathIDs 进行查询
func (s *MemoryService) SearchMemoriesByScope(ctx context.Context, keyword string, scope string, scopeCtx *types.ScopeContext) ([]entity.Memory, error) {
	// 验证关键词不能为空
	if strings.TrimSpace(keyword) == "" {
		return nil, errors.New("搜索关键词不能为空")
	}

	pathID, groupPathIDs, includeGlobal := parseScope(scope, scopeCtx)
	return s.model.SearchByScope(ctx, keyword, pathID, groupPathIDs, includeGlobal)
}

// ArchiveMemory 归档记忆
func (s *MemoryService) ArchiveMemory(ctx context.Context, id int64) error {
	// 验证ID必须大于0
	if id == 0 {
		return errors.New("记忆ID必须大于 0")
	}

	// 获取记忆实例
	memory, err := s.model.FindByID(ctx, id)
	if err != nil {
		return errors.New("记忆不存在，无法归档")
	}

	// 检查是否已经归档
	if memory.IsArchived {
		return errors.New("记忆已经归档过了")
	}

	// 执行归档
	return s.model.Archive(ctx, id)
}

// UnarchiveMemory 取消归档记忆
func (s *MemoryService) UnarchiveMemory(ctx context.Context, id int64) error {
	if id == 0 {
		return errors.New("记忆ID必须大于 0")
	}

	memory, err := s.model.FindByID(ctx, id)
	if err != nil {
		return errors.New("记忆不存在")
	}

	if !memory.IsArchived {
		return errors.New("记忆未归档")
	}

	return s.model.Unarchive(ctx, id)
}

// ToMemoryResponseDTO 将 Memory entity 转换为 ResponseDTO
// 纯关联模式：使用 PathID 判断作用域
func ToMemoryResponseDTO(memory *entity.Memory, scopeCtx *types.ScopeContext) *dto.MemoryResponseDTO {
	if memory == nil {
		return nil
	}

	tags := make([]string, 0, len(memory.Tags))
	for _, t := range memory.Tags {
		tags = append(tags, t.Tag)
	}

	// 使用 PathID 判断作用域
	var scope types.Scope
	if scopeCtx != nil {
		scope = types.GetScopeForDisplay(memory.PathID, scopeCtx.PathID, scopeCtx.GroupPathIDs)
	} else {
		scope = types.GetScope(memory.PathID)
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
