package entity

import (
	"time"

	"gorm.io/gorm"
)

// Memory 记忆实体（数据表结构）
// 嘿嘿~ 这是用于持久化存储的记忆实体！💖
type Memory struct {
	ID         uint           `gorm:"primaryKey;autoIncrement"`
	GroupID    uint           `gorm:"index;default:0;comment:所属组ID（0=Global）"`
	Path       string         `gorm:"index;size:1024;comment:精确路径（Personal作用域）"`
	Title      string         `gorm:"index;size:255;not null;comment:标题"`
	Content    string         `gorm:"type:text;not null;comment:内容"`
	Category   string         `gorm:"index;size:100;default:'默认';comment:分类"`
	Priority   int            `gorm:"default:1;comment:优先级 1-4"`
	IsArchived bool           `gorm:"index;default:false;comment:是否归档"`
	CreatedAt  time.Time      `gorm:"index;autoCreateTime"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime"`
	DeletedAt  gorm.DeletedAt `gorm:"index"` // 软删除支持

	// 关联：标签
	Tags []MemoryTag `gorm:"foreignKey:MemoryID;constraint:OnDelete:CASCADE"`
}

// TableName 指定表名
func (Memory) TableName() string {
	return "memories"
}

// MemoryTag 记忆标签关联表
// 呀~ 用于存储记忆的标签关联！✨
type MemoryTag struct {
	ID       uint   `gorm:"primaryKey;autoIncrement"`
	MemoryID uint   `gorm:"index;not null"`
	Tag      string `gorm:"index;size:100;not null"`
}

// TableName 指定表名
func (MemoryTag) TableName() string {
	return "memory_tags"
}

// MemoryPriority 记忆优先级常量
// 嘿嘿~ 统一的优先级定义！🎮
const (
	MemoryPriorityLow    = 1 // 低优先级
	MemoryPriorityMedium = 2 // 中优先级
	MemoryPriorityHigh   = 3 // 高优先级
	MemoryPriorityUrgent = 4 // 紧急优先级
)

// IsGlobal 检查是否为全局记忆
func (m *Memory) IsGlobal() bool {
	return m.GroupID == 0 && m.Path == ""
}

// IsPersonal 检查是否为 Personal 作用域
func (m *Memory) IsPersonal() bool {
	return m.Path != ""
}

// IsGroup 检查是否为 Group 作用域
func (m *Memory) IsGroup() bool {
	return m.GroupID != 0 && m.Path == ""
}

// GetScope 获取作用域类型字符串
func (m *Memory) GetScope() string {
	if m.Path != "" {
		return "personal"
	}
	if m.GroupID != 0 {
		return "group"
	}
	return "global"
}

// GetTagStrings 获取标签字符串列表
func (m *Memory) GetTagStrings() []string {
	tags := make([]string, len(m.Tags))
	for i, tag := range m.Tags {
		tags[i] = tag.Tag
	}
	return tags
}

// SetTags 设置标签（从字符串列表）
func (m *Memory) SetTags(tags []string) {
	m.Tags = make([]MemoryTag, len(tags))
	for i, tag := range tags {
		m.Tags[i] = MemoryTag{Tag: tag}
	}
}
