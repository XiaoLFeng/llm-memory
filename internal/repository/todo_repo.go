package repository

import (
	"context"
	"time"

	"github.com/XiaoLFeng/llm-memory/internal/database"
	"github.com/XiaoLFeng/llm-memory/pkg/types"
)

// TodoRepo 待办事项仓储实现
// 嘿嘿~ 这是管理待办事项的核心仓储实现呢！💖
type TodoRepo struct {
	db *database.DB
}

// NewTodoRepo 创建新的待办事项仓储实例
// 呀~ 通过依赖注入的方式创建仓储对象~ ✨
func NewTodoRepo(db *database.DB) *TodoRepo {
	return &TodoRepo{
		db: db,
	}
}

// Create 创建新的待办事项
// 保存新的待办事项到数据库~ 🎯
func (r *TodoRepo) Create(ctx context.Context, todo *types.Todo) error {
	return r.db.Save(todo)
}

// Update 更新现有待办事项
// 更新待办事项的信息哦~ ✏️
func (r *TodoRepo) Update(ctx context.Context, todo *types.Todo) error {
	// 更新 UpdatedAt 时间戳
	todo.UpdatedAt = time.Now()
	return r.db.Update(todo)
}

// Delete 删除指定ID的待办事项
// 从数据库中移除待办事项~ 🗑️
func (r *TodoRepo) Delete(ctx context.Context, id int) error {
	return r.db.DeleteStruct(&types.Todo{ID: id})
}

// FindByID 根据ID查找待办事项
// 嗯嗯！通过ID精确查找待办事项~ 🔍
func (r *TodoRepo) FindByID(ctx context.Context, id int) (*types.Todo, error) {
	var todo types.Todo
	err := r.db.One("ID", id, &todo)
	if err != nil {
		return nil, err
	}
	return &todo, nil
}

// FindAll 查找所有待办事项
// 获取全部待办事项列表~ 📋
func (r *TodoRepo) FindAll(ctx context.Context) ([]types.Todo, error) {
	var todos []types.Todo
	err := r.db.All(&todos)
	if err != nil {
		return nil, err
	}
	return todos, nil
}

// FindByStatus 根据状态查找待办事项
// 按照状态筛选待办事项~ 🎨
func (r *TodoRepo) FindByStatus(ctx context.Context, status types.TodoStatus) ([]types.Todo, error) {
	var todos []types.Todo
	err := r.db.Find("Status", status, &todos)
	if err != nil {
		return nil, err
	}
	return todos, nil
}

// FindToday 查找今天截止的待办事项
// 呀~ 找出今天需要完成的任务！⏰
func (r *TodoRepo) FindToday(ctx context.Context) ([]types.Todo, error) {
	var todos []types.Todo

	// 获取所有待办事项
	err := r.db.All(&todos)
	if err != nil {
		return nil, err
	}

	// 获取今天的开始和结束时间
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, now.Location())

	// 筛选今天截止的待办事项
	var todayTodos []types.Todo
	for _, todo := range todos {
		if todo.DueDate != nil {
			// 检查截止日期是否在今天范围内
			if !todo.DueDate.Before(startOfDay) && !todo.DueDate.After(endOfDay) {
				todayTodos = append(todayTodos, todo)
			}
		}
	}

	return todayTodos, nil
}

// FindByScope 根据作用域查找待办事项
// 嘿嘿~ 支持 Personal/Group/Global 三层作用域过滤！💖
func (r *TodoRepo) FindByScope(ctx context.Context, scope *types.ScopeContext) ([]types.Todo, error) {
	if scope == nil {
		// 没有作用域限制，返回所有
		return r.FindAll(ctx)
	}

	var allTodos []types.Todo
	err := r.db.All(&allTodos)
	if err != nil {
		return nil, err
	}

	var result []types.Todo
	for _, todo := range allTodos {
		if r.matchScope(todo, scope) {
			result = append(result, todo)
		}
	}

	return result, nil
}

// FindTodayByScope 根据作用域查找今天的待办事项
// 在指定作用域内查找今天截止的任务~ ⏰
func (r *TodoRepo) FindTodayByScope(ctx context.Context, scope *types.ScopeContext) ([]types.Todo, error) {
	// 先按作用域过滤
	todos, err := r.FindByScope(ctx, scope)
	if err != nil {
		return nil, err
	}

	// 获取今天的开始和结束时间
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, now.Location())

	// 筛选今天截止的待办事项
	var todayTodos []types.Todo
	for _, todo := range todos {
		if todo.DueDate != nil {
			if !todo.DueDate.Before(startOfDay) && !todo.DueDate.After(endOfDay) {
				todayTodos = append(todayTodos, todo)
			}
		}
	}

	return todayTodos, nil
}

// matchScope 检查待办是否匹配作用域
// 核心过滤逻辑~ ✨
func (r *TodoRepo) matchScope(todo types.Todo, scope *types.ScopeContext) bool {
	// 检查 Global
	if scope.IncludeGlobal && todo.IsGlobal() {
		return true
	}

	// 检查 Personal（精确路径匹配）
	if scope.IncludePersonal && todo.Path != "" && todo.Path == scope.CurrentPath {
		return true
	}

	// 检查 Group
	if scope.IncludeGroup && scope.GroupID != types.GlobalGroupID && todo.GroupID == scope.GroupID {
		return true
	}

	return false
}
