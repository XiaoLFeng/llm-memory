package types

import (
	"time"

	"github.com/asdine/storm/v3"
)

// Priority 表示 Todo 的优先级
type Priority int

const (
	TodoPriorityLow    Priority = iota + 1 // 低优先级
	TodoPriorityMedium                     // 中优先级
	TodoPriorityHigh                       // 高优先级
	TodoPriorityUrgent                     // 紧急优先级
)

// String 将 Priority 转换为字符串表示
func (p Priority) String() string {
	switch p {
	case TodoPriorityLow:
		return "low"
	case TodoPriorityMedium:
		return "medium"
	case TodoPriorityHigh:
		return "high"
	case TodoPriorityUrgent:
		return "urgent"
	default:
		return "unknown"
	}
}

// TodoStatus 表示 Todo 的状态
type TodoStatus int

const (
	TodoStatusPending    TodoStatus = iota // 待处理状态
	TodoStatusInProgress                   // 进行中状态
	TodoStatusCompleted                    // 已完成状态
	TodoStatusCancelled                    // 已取消状态
)

// TodoStatusToString 将 TodoStatus 转换为字符串表示
func (s TodoStatus) String() string {
	switch s {
	case TodoStatusPending:
		return "pending"
	case TodoStatusInProgress:
		return "in_progress"
	case TodoStatusCompleted:
		return "completed"
	case TodoStatusCancelled:
		return "cancelled"
	default:
		return "unknown"
	}
}

// Todo 表示一个待办事项实体
type Todo struct {
	ID          int        `storm:"id,increment"` // 主键，自增
	GroupID     int        `storm:"index"`        // 所属组ID（0=Global）
	Path        string     `storm:"index"`        // 精确路径（Personal作用域）
	Title       string     `storm:"index"`        // 标题，带索引
	Description string     `storm:""`             // 描述
	Priority    Priority   `storm:"index"`        // 优先级，带索引
	Status      TodoStatus `storm:"index"`        // 状态，带索引
	DueDate     *time.Time `storm:""`             // 截止日期
	Tags        []string   `storm:""`             // 标签
	CreatedAt   time.Time  `storm:"index"`        // 创建时间，带索引
	UpdatedAt   time.Time  `storm:""`             // 更新时间
	CompletedAt *time.Time `storm:""`             // 完成时间
}

// NewTodo 创建一个新的 Todo 实例
// 嘿嘿~ 现在支持设置作用域啦！💖
func NewTodo(title, description string, priority Priority, groupID int, path string) *Todo {
	now := time.Now()
	return &Todo{
		GroupID:     groupID,
		Path:        path,
		Title:       title,
		Description: description,
		Priority:    priority,
		Status:      TodoStatusPending,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// NewGlobalTodo 创建全局待办实例
func NewGlobalTodo(title, description string, priority Priority) *Todo {
	return NewTodo(title, description, priority, GlobalGroupID, "")
}

// NewPersonalTodo 创建 Personal 作用域的待办实例
func NewPersonalTodo(title, description string, priority Priority, path string) *Todo {
	return NewTodo(title, description, priority, GlobalGroupID, path)
}

// NewGroupTodo 创建 Group 作用域的待办实例
func NewGroupTodo(title, description string, priority Priority, groupID int) *Todo {
	return NewTodo(title, description, priority, groupID, "")
}

// IsGlobal 检查待办是否为全局待办
func (t *Todo) IsGlobal() bool {
	return t.GroupID == GlobalGroupID && t.Path == ""
}

// IsPersonal 检查待办是否为 Personal 作用域
func (t *Todo) IsPersonal() bool {
	return t.Path != ""
}

// IsGroup 检查待办是否为 Group 作用域
func (t *Todo) IsGroup() bool {
	return t.GroupID != GlobalGroupID && t.Path == ""
}

// GetScope 获取待办的作用域类型
func (t *Todo) GetScope() Scope {
	if t.Path != "" {
		return ScopePersonal
	}
	if t.GroupID != GlobalGroupID {
		return ScopeGroup
	}
	return ScopeGlobal
}

// MarkAsCompleted 标记 Todo 为已完成
func (t *Todo) MarkAsCompleted() {
	now := time.Now()
	t.Status = TodoStatusCompleted
	t.CompletedAt = &now
	t.UpdatedAt = now
}

// MarkAsInProgress 标记 Todo 为进行中
func (t *Todo) MarkAsInProgress() {
	t.Status = TodoStatusInProgress
	t.UpdatedAt = time.Now()
}

// MarkAsCancelled 标记 Todo 为已取消
func (t *Todo) MarkAsCancelled() {
	t.Status = TodoStatusCancelled
	t.UpdatedAt = time.Now()
}

// IsOverdue 检查 Todo 是否已过期
func (t *Todo) IsOverdue() bool {
	if t.DueDate == nil {
		return false
	}
	return time.Now().After(*t.DueDate) && t.Status != TodoStatusCompleted
}

// AddTag 添加标签
func (t *Todo) AddTag(tag string) {
	for _, existingTag := range t.Tags {
		if existingTag == tag {
			return // 标签已存在
		}
	}
	t.Tags = append(t.Tags, tag)
	t.UpdatedAt = time.Now()
}

// RemoveTag 移除标签
func (t *Todo) RemoveTag(tag string) {
	for i, existingTag := range t.Tags {
		if existingTag == tag {
			t.Tags = append(t.Tags[:i], t.Tags[i+1:]...)
			t.UpdatedAt = time.Now()
			return
		}
	}
}

// Validate 验证 Todo 数据的有效性
func (t *Todo) Validate() error {
	if t.Title == "" {
		return storm.ErrAlreadyExists
	}
	return nil
}
