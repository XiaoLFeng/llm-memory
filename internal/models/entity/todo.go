package entity

import (
	"time"

	"gorm.io/gorm"
)

// ToDoStatus 待办状态类型
// 用整数类型方便数据库存储
type ToDoStatus int

// 待办状态常量定义
const (
	ToDoStatusPending    ToDoStatus = iota // 待处理状态
	ToDoStatusInProgress                   // 进行中状态
	ToDoStatusCompleted                    // 已完成状态
	ToDoStatusCancelled                    // 已取消状态
)

// String 将 ToDoStatus 转换为字符串表示
func (s ToDoStatus) String() string {
	switch s {
	case ToDoStatusPending:
		return "pending"
	case ToDoStatusInProgress:
		return "in_progress"
	case ToDoStatusCompleted:
		return "completed"
	case ToDoStatusCancelled:
		return "cancelled"
	default:
		return "unknown"
	}
}

// ToDoStatusFromString 从字符串转换为 ToDoStatus
func ToDoStatusFromString(s string) ToDoStatus {
	switch s {
	case "pending":
		return ToDoStatusPending
	case "in_progress":
		return ToDoStatusInProgress
	case "completed":
		return ToDoStatusCompleted
	case "cancelled":
		return ToDoStatusCancelled
	default:
		return ToDoStatusPending
	}
}

// ToDoPriority 待办优先级类型
// 优先级从 1 开始
type ToDoPriority int

// 待办优先级常量定义
const (
	ToDoPriorityLow    ToDoPriority = iota + 1 // 低优先级
	ToDoPriorityMedium                         // 中优先级
	ToDoPriorityHigh                           // 高优先级
	ToDoPriorityUrgent                         // 紧急优先级
)

// String 将 ToDoPriority 转换为字符串表示
func (p ToDoPriority) String() string {
	switch p {
	case ToDoPriorityLow:
		return "low"
	case ToDoPriorityMedium:
		return "medium"
	case ToDoPriorityHigh:
		return "high"
	case ToDoPriorityUrgent:
		return "urgent"
	default:
		return "unknown"
	}
}

// ToDoPriorityFromString 从字符串转换为 ToDoPriority
func ToDoPriorityFromString(s string) ToDoPriority {
	switch s {
	case "low":
		return ToDoPriorityLow
	case "medium":
		return ToDoPriorityMedium
	case "high":
		return ToDoPriorityHigh
	case "urgent":
		return ToDoPriorityUrgent
	default:
		return ToDoPriorityMedium
	}
}

// ToDo 待办事项实体（数据表结构）
// 注意：类型名使用 ToDo（不是 Todo），避免 IDE 命名规范问题
// 用于管理短期任务的待办实体
type ToDo struct {
	ID          int64          `gorm:"primaryKey"`                              // 雪花算法生成
	GroupID     int64          `gorm:"index;default:0;comment:所属组ID（0=Global）"` // 关联组ID
	Path        string         `gorm:"index;size:1024;comment:精确路径（Personal作用域）"`
	Title       string         `gorm:"index;size:255;not null;comment:标题"`
	Description string         `gorm:"type:text;comment:描述"`
	Priority    ToDoPriority   `gorm:"index;default:2;comment:优先级 1-4"`
	Status      ToDoStatus     `gorm:"index;default:0;comment:状态"`
	DueDate     *time.Time     `gorm:"index;comment:截止日期"`
	CompletedAt *time.Time     `gorm:"comment:完成时间"`
	CreatedAt   time.Time      `gorm:"index;autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"` // 软删除支持

	// 关联：标签
	Tags []ToDoTag `gorm:"foreignKey:ToDoID;constraint:OnDelete:CASCADE"`
}

// TableName 指定表名
func (ToDo) TableName() string {
	return "to_dos"
}

// ToDoTag 待办标签关联表
// 存储待办的标签关联
type ToDoTag struct {
	ID     int64  `gorm:"primaryKey"`     // 雪花算法生成
	ToDoID int64  `gorm:"index;not null"` // 关联待办ID
	Tag    string `gorm:"index;size:100;not null"`
}

// TableName 指定表名
func (ToDoTag) TableName() string {
	return "to_do_tags"
}

// IsGlobal 检查是否为全局待办
func (t *ToDo) IsGlobal() bool {
	return t.GroupID == 0 && t.Path == ""
}

// IsPersonal 检查是否为 Personal 作用域
func (t *ToDo) IsPersonal() bool {
	return t.Path != ""
}

// IsGroup 检查是否为 Group 作用域
func (t *ToDo) IsGroup() bool {
	return t.GroupID != 0 && t.Path == ""
}

// GetScope 获取作用域类型字符串
func (t *ToDo) GetScope() string {
	if t.Path != "" {
		return "personal"
	}
	if t.GroupID != 0 {
		return "group"
	}
	return "global"
}

// GetTagStrings 获取标签字符串列表
func (t *ToDo) GetTagStrings() []string {
	tags := make([]string, len(t.Tags))
	for i, tag := range t.Tags {
		tags[i] = tag.Tag
	}
	return tags
}

// SetTags 设置标签（从字符串列表）
func (t *ToDo) SetTags(tags []string) {
	t.Tags = make([]ToDoTag, len(tags))
	for i, tag := range tags {
		t.Tags[i] = ToDoTag{Tag: tag}
	}
}

// MarkAsCompleted 标记为已完成
func (t *ToDo) MarkAsCompleted() {
	now := time.Now()
	t.Status = ToDoStatusCompleted
	t.CompletedAt = &now
}

// MarkAsInProgress 标记为进行中
func (t *ToDo) MarkAsInProgress() {
	t.Status = ToDoStatusInProgress
}

// MarkAsCancelled 标记为已取消
func (t *ToDo) MarkAsCancelled() {
	t.Status = ToDoStatusCancelled
}

// IsOverdue 检查是否已过期
func (t *ToDo) IsOverdue() bool {
	if t.DueDate == nil {
		return false
	}
	return time.Now().After(*t.DueDate) && t.Status != ToDoStatusCompleted
}

// IsCompleted 检查是否已完成
func (t *ToDo) IsCompleted() bool {
	return t.Status == ToDoStatusCompleted
}
