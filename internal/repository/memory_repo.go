package repository

import (
	"context"
	"errors"

	"github.com/XiaoLFeng/llm-memory/internal/database"
	"github.com/XiaoLFeng/llm-memory/pkg/types"
)

// MemoryRepo 记忆仓储实现
// 嘿嘿~ 这是 MemoryRepository 接口的具体实现哦！💖
// 使用 storm 数据库进行优雅的 CRUD 操作~ ✨
type MemoryRepo struct {
	db *database.DB
}

// NewMemoryRepo 创建新的记忆仓储实例
// 呀~ 构造函数来啦！接收一个 DB 实例~ (´∀｀)
func NewMemoryRepo(db *database.DB) *MemoryRepo {
	return &MemoryRepo{
		db: db,
	}
}

// Create 创建新的记忆
// 使用 db.Save 方法优雅地保存到数据库~ 🎯
func (r *MemoryRepo) Create(ctx context.Context, memory *types.Memory) error {
	if memory == nil {
		return errors.New("memory 不能为空哦~ 💫")
	}
	// TODO: 未来可以使用 ctx 实现超时控制
	return r.db.Save(memory)
}

// Update 更新现有记忆
// 自动更新 UpdatedAt 字段，然后使用 db.Update~ ✨
func (r *MemoryRepo) Update(ctx context.Context, memory *types.Memory) error {
	if memory == nil {
		return errors.New("memory 不能为空哦~ 💫")
	}

	// 自动设置更新时间，嘿嘿~ 这样数据更完整呢！
	memory.BeforeUpdate()

	return r.db.Update(memory)
}

// Delete 删除指定ID的记忆
// 使用 db.DeleteStruct 方法优雅地删除~ 💨
func (r *MemoryRepo) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("ID 必须大于 0 哦~ 🎮")
	}

	memory := &types.Memory{ID: id}
	return r.db.DeleteStruct(memory)
}

// FindByID 根据ID查找记忆
// 使用 db.One 方法精准查找！就像寻宝一样~ 🏴‍☠️
func (r *MemoryRepo) FindByID(ctx context.Context, id int) (*types.Memory, error) {
	if id <= 0 {
		return nil, errors.New("ID 必须大于 0 哦~ 🎮")
	}

	var memory types.Memory
	err := r.db.One("ID", id, &memory)
	if err != nil {
		return nil, err
	}

	return &memory, nil
}

// FindAll 查找所有记忆
// 使用 db.All 方法获取所有记忆，就像打开记忆相册一样~ 📖
func (r *MemoryRepo) FindAll(ctx context.Context) ([]types.Memory, error) {
	var memories []types.Memory
	err := r.db.All(&memories)
	if err != nil {
		return nil, err
	}

	return memories, nil
}

// FindByCategory 根据分类查找记忆
// 使用 db.Find 方法按分类筛选，整理记忆就靠它了！🏷️
func (r *MemoryRepo) FindByCategory(ctx context.Context, category string) ([]types.Memory, error) {
	if category == "" {
		return nil, errors.New("分类名称不能为空哦~ 📝")
	}

	var memories []types.Memory
	err := r.db.Find("Category", category, &memories)
	if err != nil {
		return nil, err
	}

	return memories, nil
}

// Search 根据关键词搜索记忆
// 呀~ 这是个智能搜索功能！在标题和内容中查找关键词~ 🔍
// 支持模糊匹配，让记忆检索更方便呢！💫
func (r *MemoryRepo) Search(ctx context.Context, keyword string) ([]types.Memory, error) {
	if keyword == "" {
		return nil, errors.New("搜索关键词不能为空哦~ 🎯")
	}

	var memories []types.Memory
	var allMemories []types.Memory

	// 先获取所有记忆
	err := r.db.All(&allMemories)
	if err != nil {
		return nil, err
	}

	// 过滤包含关键词的记忆
	for _, memory := range allMemories {
		// 在标题中搜索关键词
		titleMatch := contains(memory.Title, keyword)

		// 在内容中搜索关键词
		contentMatch := contains(memory.Content, keyword)

		// 如果标题或内容包含关键词，就添加到结果中
		// 嘿嘿~ 这样就可以从多个地方找到记忆啦！✨
		if titleMatch || contentMatch {
			memories = append(memories, memory)
		}
	}

	return memories, nil
}

// contains 辅助函数：检查字符串是否包含关键词
// 呀~ 简单的字符串匹配，就像玩文字游戏一样！🎮
func contains(text, keyword string) bool {
	textRunes := []rune(text)
	keywordRunes := []rune(keyword)
	keywordLen := len(keywordRunes)

	if keywordLen == 0 || len(textRunes) < keywordLen {
		return false
	}

	for i := 0; i <= len(textRunes)-keywordLen; i++ {
		match := true
		for j := 0; j < keywordLen; j++ {
			if textRunes[i+j] != keywordRunes[j] {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}

	return false
}
