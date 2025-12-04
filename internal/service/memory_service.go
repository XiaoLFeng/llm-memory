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

	// 解析作用域 -> PathID（global=true 存全局；否则使用当前路径）
	pathID := int64(0)
	if !input.Global {
		pathID = resolveDefaultPathID(scopeCtx)
		if pathID == 0 {
			return nil, errors.New("无法确定私有/小组作用域，请先初始化 paths 或指定全局模式")
		}
	}

	// 创建记忆实例
	memory := &entity.Memory{
		Global:   input.Global,
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
	return s.model.FindByFilter(ctx, models.DefaultVisibilityFilter())
}

// ListMemoriesByScope 根据作用域列出记忆
// scope 参数: personal/group/global/all
func (s *MemoryService) ListMemoriesByScope(ctx context.Context, scope string, scopeCtx *types.ScopeContext) ([]entity.Memory, error) {
	filter := buildVisibilityFilter(scope, scopeCtx)
	return s.model.FindByFilter(ctx, filter)
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

	return s.model.SearchByFilter(ctx, keyword, models.DefaultVisibilityFilter())
}

// SearchMemoriesByScope 根据作用域搜索记忆
// 纯关联模式：使用 PathID 和 GroupPathIDs 进行查询
func (s *MemoryService) SearchMemoriesByScope(ctx context.Context, keyword string, scope string, scopeCtx *types.ScopeContext) ([]entity.Memory, error) {
	// 验证关键词不能为空
	if strings.TrimSpace(keyword) == "" {
		return nil, errors.New("搜索关键词不能为空")
	}

	filter := buildVisibilityFilter(scope, scopeCtx)
	return s.model.SearchByFilter(ctx, keyword, filter)
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

	scope := types.GetScopeForDisplayWithGlobal(memory.Global, memory.PathID, scopeCtx)

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
