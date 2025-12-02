package types

import (
	"time"
)

// Memory 记忆实体结构体
type Memory struct {
	ID         int       `json:"id" storm:"id,increment"`   // 主键，自增
	Title      string    `json:"title" storm:"index"`       // 标题，带索引
	Content    string    `json:"content"`                   // 内容
	Category   string    `json:"category" storm:"index"`    // 分类，带索引
	Tags       []string  `json:"tags"`                      // 标签
	Priority   int       `json:"priority"`                  // 优先级
	CreatedAt  time.Time `json:"created_at" storm:"index"`  // 创建时间，带索引
	UpdatedAt  time.Time `json:"updated_at"`                // 更新时间
	IsArchived bool      `json:"is_archived" storm:"index"` // 是否归档，带索引
}

// MemoryCategory 记忆分类结构体
type MemoryCategory struct {
	ID          int       `json:"id" storm:"id,increment"` // 主键，自增
	Name        string    `json:"name" storm:"unique"`     // 分类名称，唯一
	Description string    `json:"description"`             // 分类描述
	Color       string    `json:"color"`                   // 分类颜色
	CreatedAt   time.Time `json:"created_at"`              // 创建时间
	UpdatedAt   time.Time `json:"updated_at"`              // 更新时间
}

// 优先级常量定义
const (
	PriorityLow    = 1 // 低优先级
	PriorityNormal = 2 // 普通优先级
	PriorityHigh   = 3 // 高优先级
	PriorityUrgent = 4 // 紧急优先级
)

// 优先级名称映射
var PriorityNames = map[int]string{
	PriorityLow:    "低优先级",
	PriorityNormal: "普通优先级",
	PriorityHigh:   "高优先级",
	PriorityUrgent: "紧急优先级",
}

// GetPriorityName 获取优先级名称
func (m *Memory) GetPriorityName() string {
	if name, exists := PriorityNames[m.Priority]; exists {
		return name
	}
	return "未知优先级"
}

// NewMemory 创建新的记忆实例
func NewMemory(title, content, category string, tags []string, priority int) *Memory {
	now := time.Now()
	return &Memory{
		Title:      title,
		Content:    content,
		Category:   category,
		Tags:       tags,
		Priority:   priority,
		CreatedAt:  now,
		UpdatedAt:  now,
		IsArchived: false,
	}
}

// NewMemoryCategory 创建新的记忆分类实例
func NewMemoryCategory(name, description, color string) *MemoryCategory {
	now := time.Now()
	return &MemoryCategory{
		Name:        name,
		Description: description,
		Color:       color,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// BeforeUpdate 在更新前自动设置更新时间
func (m *Memory) BeforeUpdate() error {
	m.UpdatedAt = time.Now()
	return nil
}

// BeforeUpdate 在更新前自动设置更新时间
func (mc *MemoryCategory) BeforeUpdate() error {
	mc.UpdatedAt = time.Now()
	return nil
}
