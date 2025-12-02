package database

import (
	"context"

	"github.com/XiaoLFeng/llm-memory/pkg/types"
)

// MemoryRepository 记忆仓储接口
// 嘿嘿~ 这是记忆管理的核心接口呢！💖
type MemoryRepository interface {
	// Create 创建新的记忆
	Create(ctx context.Context, memory *types.Memory) error

	// Update 更新现有记忆
	Update(ctx context.Context, memory *types.Memory) error

	// Delete 删除指定ID的记忆
	Delete(ctx context.Context, id int) error

	// FindByID 根据ID查找记忆
	FindByID(ctx context.Context, id int) (*types.Memory, error)

	// FindAll 查找所有记忆
	FindAll(ctx context.Context) ([]types.Memory, error)

	// FindByCategory 根据分类查找记忆
	FindByCategory(ctx context.Context, category string) ([]types.Memory, error)

	// Search 根据关键词搜索记忆
	Search(ctx context.Context, keyword string) ([]types.Memory, error)
}

// PlanRepository 计划仓储接口
// 呀~ 用于管理所有的计划和任务哦！📋
type PlanRepository interface {
	// Create 创建新的计划
	Create(ctx context.Context, plan *types.Plan) error

	// Update 更新现有计划
	Update(ctx context.Context, plan *types.Plan) error

	// Delete 删除指定ID的计划
	Delete(ctx context.Context, id int) error

	// FindByID 根据ID查找计划
	FindByID(ctx context.Context, id int) (*types.Plan, error)

	// FindAll 查找所有计划
	FindAll(ctx context.Context) ([]types.Plan, error)

	// FindByStatus 根据状态查找计划
	FindByStatus(ctx context.Context, status types.PlanStatus) ([]types.Plan, error)
}

// TodoRepository 待办事项仓储接口
// 嗯嗯！管理所有待办任务的接口~ ✨
type TodoRepository interface {
	// Create 创建新的待办事项
	Create(ctx context.Context, todo *types.Todo) error

	// Update 更新现有待办事项
	Update(ctx context.Context, todo *types.Todo) error

	// Delete 删除指定ID的待办事项
	Delete(ctx context.Context, id int) error

	// FindByID 根据ID查找待办事项
	FindByID(ctx context.Context, id int) (*types.Todo, error)

	// FindAll 查找所有待办事项
	FindAll(ctx context.Context) ([]types.Todo, error)

	// FindByStatus 根据状态查找待办事项
	FindByStatus(ctx context.Context, status types.TodoStatus) ([]types.Todo, error)

	// FindToday 查找今天的待办事项
	FindToday(ctx context.Context) ([]types.Todo, error)
}
