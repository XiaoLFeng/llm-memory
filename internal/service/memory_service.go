package service

import (
	"context"
	"errors"
	"strings"

	"github.com/XiaoLFeng/llm-memory/internal/database"
	"github.com/XiaoLFeng/llm-memory/pkg/types"
)

// MemoryService 记忆服务结构体
// 嘿嘿~ 这是处理记忆业务逻辑的服务层哦！💖
// 负责验证、处理和协调各种记忆操作~ ✨
type MemoryService struct {
	repo database.MemoryRepository
}

// NewMemoryService 创建新的记忆服务实例
// 呀~ 构造函数来啦！接收一个 MemoryRepository 实例~ (´∀｀)
func NewMemoryService(repo database.MemoryRepository) *MemoryService {
	return &MemoryService{
		repo: repo,
	}
}

// CreateMemory 创建新的记忆
// 嘿嘿~ 创建记忆前会先验证数据的完整性呢！💫
// 参数验证通过后才会调用仓储层创建~ 🎯
func (s *MemoryService) CreateMemory(ctx context.Context, title, content, category string, tags []string, priority int, groupID int, path string) (*types.Memory, error) {
	// 验证标题不能为空
	if strings.TrimSpace(title) == "" {
		return nil, errors.New("标题不能为空哦~ 📝")
	}

	// 验证内容不能为空
	if strings.TrimSpace(content) == "" {
		return nil, errors.New("内容不能为空哦~ 📖")
	}

	// 验证分类不能为空
	if strings.TrimSpace(category) == "" {
		return nil, errors.New("分类不能为空哦~ 🏷️")
	}

	// 验证优先级范围（1-4）
	// 呀~ 优先级必须在合法范围内呢！✨
	if priority < types.PriorityLow || priority > types.PriorityUrgent {
		return nil, errors.New("优先级必须在 1-4 之间哦~ 🎮")
	}

	// 创建记忆实例
	// 嗯嗯！使用 types 包的构造函数，优雅地初始化~ 💖
	memory := types.NewMemory(title, content, category, tags, priority, groupID, path)

	// 保存到数据库
	err := s.repo.Create(ctx, memory)
	if err != nil {
		return nil, err
	}

	return memory, nil
}

// CreateGlobalMemory 创建全局记忆
// 便捷方法，创建 Global 作用域的记忆~ 🌐
func (s *MemoryService) CreateGlobalMemory(ctx context.Context, title, content, category string, tags []string, priority int) (*types.Memory, error) {
	return s.CreateMemory(ctx, title, content, category, tags, priority, types.GlobalGroupID, "")
}

// CreatePersonalMemory 创建 Personal 作用域的记忆
// 便捷方法，创建属于特定路径的记忆~ 📍
func (s *MemoryService) CreatePersonalMemory(ctx context.Context, title, content, category string, tags []string, priority int, path string) (*types.Memory, error) {
	return s.CreateMemory(ctx, title, content, category, tags, priority, types.GlobalGroupID, path)
}

// CreateGroupMemory 创建 Group 作用域的记忆
// 便捷方法，创建属于特定组的记忆~ 👥
func (s *MemoryService) CreateGroupMemory(ctx context.Context, title, content, category string, tags []string, priority int, groupID int) (*types.Memory, error) {
	return s.CreateMemory(ctx, title, content, category, tags, priority, groupID, "")
}

// UpdateMemory 更新记忆
// 呀~ 更新前会验证记忆的完整性，确保数据有效！✨
func (s *MemoryService) UpdateMemory(ctx context.Context, memory *types.Memory) error {
	// 验证记忆实例不能为空
	if memory == nil {
		return errors.New("记忆实例不能为空哦~ 💫")
	}

	// 验证ID必须大于0
	if memory.ID <= 0 {
		return errors.New("记忆ID必须大于 0 哦~ 🎮")
	}

	// 验证标题不能为空
	if strings.TrimSpace(memory.Title) == "" {
		return errors.New("标题不能为空哦~ 📝")
	}

	// 验证内容不能为空
	if strings.TrimSpace(memory.Content) == "" {
		return errors.New("内容不能为空哦~ 📖")
	}

	// 验证分类不能为空
	if strings.TrimSpace(memory.Category) == "" {
		return errors.New("分类不能为空哦~ 🏷️")
	}

	// 验证优先级范围
	if memory.Priority < types.PriorityLow || memory.Priority > types.PriorityUrgent {
		return errors.New("优先级必须在 1-4 之间哦~ 🎮")
	}

	// 检查记忆是否存在
	// 嘿嘿~ 更新前要先确认记忆存在呢！💖
	existingMemory, err := s.repo.FindByID(ctx, memory.ID)
	if err != nil {
		return errors.New("记忆不存在，无法更新哦~ 🔍")
	}

	if existingMemory == nil {
		return errors.New("记忆不存在，无法更新哦~ 🔍")
	}

	// 执行更新操作
	return s.repo.Update(ctx, memory)
}

// DeleteMemory 删除记忆
// 呀~ 删除前会验证ID和记忆是否存在！💨
func (s *MemoryService) DeleteMemory(ctx context.Context, id int) error {
	// 验证ID必须大于0
	if id <= 0 {
		return errors.New("记忆ID必须大于 0 哦~ 🎮")
	}

	// 检查记忆是否存在
	// 嗯嗯！删除不存在的记忆可不行呢~ 🔍
	existingMemory, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New("记忆不存在，无法删除哦~ 💫")
	}

	if existingMemory == nil {
		return errors.New("记忆不存在，无法删除哦~ 💫")
	}

	// 执行删除操作
	return s.repo.Delete(ctx, id)
}

// GetMemory 获取单个记忆
// 嘿嘿~ 根据ID精准查找记忆！就像寻宝一样~ 🏴‍☠️
func (s *MemoryService) GetMemory(ctx context.Context, id int) (*types.Memory, error) {
	// 验证ID必须大于0
	if id <= 0 {
		return nil, errors.New("记忆ID必须大于 0 哦~ 🎮")
	}

	// 从仓储层获取记忆
	return s.repo.FindByID(ctx, id)
}

// ListMemories 列出所有记忆
// 呀~ 获取所有记忆列表，就像打开记忆相册一样！📖
func (s *MemoryService) ListMemories(ctx context.Context) ([]types.Memory, error) {
	return s.repo.FindAll(ctx)
}

// ListByCategory 根据分类列出记忆
// 嗯嗯！按分类筛选记忆，让记忆管理更有条理~ 🏷️
func (s *MemoryService) ListByCategory(ctx context.Context, category string) ([]types.Memory, error) {
	// 验证分类不能为空
	if strings.TrimSpace(category) == "" {
		return nil, errors.New("分类名称不能为空哦~ 📝")
	}

	return s.repo.FindByCategory(ctx, category)
}

// SearchMemories 搜索记忆
// 嘿嘿~ 智能搜索功能！根据关键词在标题和内容中查找~ 🔍
func (s *MemoryService) SearchMemories(ctx context.Context, keyword string) ([]types.Memory, error) {
	// 验证关键词不能为空
	if strings.TrimSpace(keyword) == "" {
		return nil, errors.New("搜索关键词不能为空哦~ 🎯")
	}

	return s.repo.Search(ctx, keyword)
}

// ArchiveMemory 归档记忆
// 呀~ 将记忆标记为已归档状态！💼
// 归档后的记忆不会显示在常规列表中~ ✨
func (s *MemoryService) ArchiveMemory(ctx context.Context, id int) error {
	// 验证ID必须大于0
	if id <= 0 {
		return errors.New("记忆ID必须大于 0 哦~ 🎮")
	}

	// 获取记忆实例
	memory, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New("记忆不存在，无法归档哦~ 💫")
	}

	if memory == nil {
		return errors.New("记忆不存在，无法归档哦~ 💫")
	}

	// 检查是否已经归档
	// 嘿嘿~ 避免重复归档呢！✨
	if memory.IsArchived {
		return errors.New("记忆已经归档过了哦~ 📦")
	}

	// 设置为归档状态
	memory.IsArchived = true

	// 更新到数据库
	// 嗯嗯！使用 BeforeUpdate 自动更新时间戳~ 💖
	return s.repo.Update(ctx, memory)
}

// ListMemoriesByScope 根据作用域列出记忆
// 嘿嘿~ 支持 Personal/Group/Global 三层作用域过滤！💖
func (s *MemoryService) ListMemoriesByScope(ctx context.Context, scope *types.ScopeContext) ([]types.Memory, error) {
	return s.repo.FindByScope(ctx, scope)
}

// SearchMemoriesByScope 根据作用域搜索记忆
// 在指定作用域内搜索关键词~ 🔍
func (s *MemoryService) SearchMemoriesByScope(ctx context.Context, scope *types.ScopeContext, keyword string) ([]types.Memory, error) {
	// 验证关键词不能为空
	if strings.TrimSpace(keyword) == "" {
		return nil, errors.New("搜索关键词不能为空哦~ 🎯")
	}

	return s.repo.SearchByScope(ctx, scope, keyword)
}
