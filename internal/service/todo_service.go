package service

import (
	"context"
	"errors"
	"time"

	"github.com/XiaoLFeng/llm-memory/internal/database"
	"github.com/XiaoLFeng/llm-memory/pkg/types"
)

// TodoService 待办事项服务
// 嘿嘿~ 这是处理待办事项业务逻辑的服务层哦！💖
type TodoService struct {
	repo database.TodoRepository
}

// NewTodoService 创建新的待办事项服务实例
// 呀~ 构造函数来啦！(´∀｀)
func NewTodoService(repo database.TodoRepository) *TodoService {
	return &TodoService{
		repo: repo,
	}
}

// CreateTodo 创建新的待办事项
// 嘿嘿~ 创建待办前会先验证数据的完整性呢！💫
// 参数：
//   - ctx: 上下文
//   - title: 标题
//   - description: 描述
//   - priority: 优先级
//   - dueDate: 截止日期（可选）
//   - groupID: 组ID（0=Global）
//   - path: 路径（Personal 作用域）
//
// 返回：
//   - 创建的待办事项
//   - 错误信息（如果有的话）
func (s *TodoService) CreateTodo(ctx context.Context, title, description string, priority types.Priority, dueDate *time.Time, groupID int, path string) (*types.Todo, error) {
	// 验证标题不能为空
	if title == "" {
		return nil, errors.New("标题不能为空哦~ 📝")
	}

	// 验证优先级的有效性
	if priority < types.TodoPriorityLow || priority > types.TodoPriorityUrgent {
		return nil, errors.New("无效的优先级哦~ 🎮")
	}

	// 创建新的待办事项实例
	// 嗯嗯！使用 types 包的构造函数，优雅地初始化~ 💖
	todo := types.NewTodo(title, description, priority, groupID, path)
	todo.DueDate = dueDate

	// 保存到数据库
	if err := s.repo.Create(ctx, todo); err != nil {
		return nil, err
	}

	return todo, nil
}

// CreateGlobalTodo 创建全局待办事项
// 便捷方法，创建 Global 作用域的待办~ 🌐
func (s *TodoService) CreateGlobalTodo(ctx context.Context, title, description string, priority types.Priority, dueDate *time.Time) (*types.Todo, error) {
	return s.CreateTodo(ctx, title, description, priority, dueDate, types.GlobalGroupID, "")
}

// CreatePersonalTodo 创建 Personal 作用域的待办事项
// 便捷方法，创建属于特定路径的待办~ 📍
func (s *TodoService) CreatePersonalTodo(ctx context.Context, title, description string, priority types.Priority, dueDate *time.Time, path string) (*types.Todo, error) {
	return s.CreateTodo(ctx, title, description, priority, dueDate, types.GlobalGroupID, path)
}

// CreateGroupTodo 创建 Group 作用域的待办事项
// 便捷方法，创建属于特定组的待办~ 👥
func (s *TodoService) CreateGroupTodo(ctx context.Context, title, description string, priority types.Priority, dueDate *time.Time, groupID int) (*types.Todo, error) {
	return s.CreateTodo(ctx, title, description, priority, dueDate, groupID, "")
}

// UpdateTodo 更新待办事项
// 参数：
//   - ctx: 上下文
//   - todo: 要更新的待办事项
//
// 返回：
//   - 错误信息（如果有的话）
func (s *TodoService) UpdateTodo(ctx context.Context, todo *types.Todo) error {
	// 验证待办事项不能为空
	if todo == nil {
		return errors.New("待办事项不能为空")
	}

	// 验证标题不能为空
	if todo.Title == "" {
		return errors.New("标题不能为空")
	}

	// 验证优先级的有效性
	if todo.Priority < types.TodoPriorityLow || todo.Priority > types.TodoPriorityUrgent {
		return errors.New("无效的优先级")
	}

	// 更新时间戳
	todo.UpdatedAt = time.Now()

	// 保存更新
	return s.repo.Update(ctx, todo)
}

// DeleteTodo 删除待办事项
// 参数：
//   - ctx: 上下文
//   - id: 待办事项ID
//
// 返回：
//   - 错误信息（如果有的话）
func (s *TodoService) DeleteTodo(ctx context.Context, id int) error {
	// 验证ID有效性
	if id <= 0 {
		return errors.New("无效的待办事项ID")
	}

	// 检查待办事项是否存在
	if _, err := s.repo.FindByID(ctx, id); err != nil {
		return errors.New("待办事项不存在")
	}

	// 执行删除
	return s.repo.Delete(ctx, id)
}

// GetTodo 获取指定ID的待办事项
// 参数：
//   - ctx: 上下文
//   - id: 待办事项ID
//
// 返回：
//   - 待办事项
//   - 错误信息（如果有的话）
func (s *TodoService) GetTodo(ctx context.Context, id int) (*types.Todo, error) {
	// 验证ID有效性
	if id <= 0 {
		return nil, errors.New("无效的待办事项ID")
	}

	return s.repo.FindByID(ctx, id)
}

// ListTodos 获取所有待办事项
// 返回：
//   - 待办事项列表
//   - 错误信息（如果有的话）
func (s *TodoService) ListTodos(ctx context.Context) ([]types.Todo, error) {
	return s.repo.FindAll(ctx)
}

// ListByStatus 根据状态获取待办事项列表
// 参数：
//   - ctx: 上下文
//   - status: 待办事项状态
//
// 返回：
//   - 待办事项列表
//   - 错误信息（如果有的话）
func (s *TodoService) ListByStatus(ctx context.Context, status types.TodoStatus) ([]types.Todo, error) {
	// 验证状态的有效性
	if status < types.TodoStatusPending || status > types.TodoStatusCancelled {
		return nil, errors.New("无效的状态")
	}

	return s.repo.FindByStatus(ctx, status)
}

// ListToday 获取今天的待办事项
// 返回：
//   - 今天的待办事项列表
//   - 错误信息（如果有的话）
func (s *TodoService) ListToday(ctx context.Context) ([]types.Todo, error) {
	return s.repo.FindToday(ctx)
}

// CompleteTodo 标记待办事项为已完成
// 参数：
//   - ctx: 上下文
//   - id: 待办事项ID
//
// 返回：
//   - 错误信息（如果有的话）
func (s *TodoService) CompleteTodo(ctx context.Context, id int) error {
	// 验证ID有效性
	if id <= 0 {
		return errors.New("无效的待办事项ID")
	}

	// 获取待办事项
	todo, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// 检查状态是否已经完成
	if todo.Status == types.TodoStatusCompleted {
		return errors.New("待办事项已经完成")
	}

	// 标记为已完成
	todo.MarkAsCompleted()

	// 保存更新
	return s.repo.Update(ctx, todo)
}

// StartTodo 标记待办事项为进行中
// 参数：
//   - ctx: 上下文
//   - id: 待办事项ID
//
// 返回：
//   - 错误信息（如果有的话）
func (s *TodoService) StartTodo(ctx context.Context, id int) error {
	// 验证ID有效性
	if id <= 0 {
		return errors.New("无效的待办事项ID")
	}

	// 获取待办事项
	todo, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// 检查状态是否已经完成或取消
	if todo.Status == types.TodoStatusCompleted {
		return errors.New("已完成的待办事项无法开始")
	}
	if todo.Status == types.TodoStatusCancelled {
		return errors.New("已取消的待办事项无法开始")
	}

	// 标记为进行中
	todo.MarkAsInProgress()

	// 保存更新
	return s.repo.Update(ctx, todo)
}

// ListTodosByScope 根据作用域列出待办事项
// 嘿嘿~ 支持 Personal/Group/Global 三层作用域过滤！💖
func (s *TodoService) ListTodosByScope(ctx context.Context, scope *types.ScopeContext) ([]types.Todo, error) {
	return s.repo.FindByScope(ctx, scope)
}

// ListTodayByScope 根据作用域获取今天的待办事项
// 在指定作用域内查找今天截止的任务~ ⏰
func (s *TodoService) ListTodayByScope(ctx context.Context, scope *types.ScopeContext) ([]types.Todo, error) {
	return s.repo.FindTodayByScope(ctx, scope)
}
