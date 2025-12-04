package entity

import (
	"time"
)

// Memory 记忆实体（数据表结构）
// 记忆条目，用于持久化存储重要信息
// 纯关联模式：PathID=0 表示 Global，PathID>0 关联 PersonalPath
type Memory struct {
	ID         int64     `gorm:"primaryKey"`                               // 雪花算法生成
	Global     bool      `gorm:"index;default:false;comment:是否全局可见"`       // true=全局，false=私有/小组
	PathID     int64     `gorm:"index;default:0;comment:关联路径ID(0=无绑定/全局)"` // 关联 Path.ID，0 表示未绑定
	Title      string    `gorm:"index;size:255;not null;comment:标题"`
	Content    string    `gorm:"type:text;not null;comment:内容"`
	Category   string    `gorm:"index;size:100;default:'默认';comment:分类"`
	Priority   int       `gorm:"default:1;comment:优先级 1-4"`
	IsArchived bool      `gorm:"index;default:false;comment:是否归档"`
	CreatedAt  time.Time `gorm:"index;autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`

	// 关联：标签
	Tags []MemoryTag `gorm:"foreignKey:MemoryID;constraint:OnDelete:CASCADE"`
}

// TableName 指定表名
func (Memory) TableName() string {
	return "memories"
}

// MemoryTag 记忆标签关联表
// 存储记忆的标签关联
type MemoryTag struct {
	ID       int64  `gorm:"primaryKey"`     // 雪花算法生成
	MemoryID int64  `gorm:"index;not null"` // 关联记忆ID
	Tag      string `gorm:"index;size:100;not null"`
}

// TableName 指定表名
func (MemoryTag) TableName() string {
	return "memory_tags"
}

// MemoryPriority 记忆优先级常量
// 统一的优先级定义
const (
	MemoryPriorityLow    = 1 // 低优先级
	MemoryPriorityMedium = 2 // 中优先级
	MemoryPriorityHigh   = 3 // 高优先级
	MemoryPriorityUrgent = 4 // 紧急优先级
)

// IsGlobal 检查是否为全局记忆
func (m *Memory) IsGlobal() bool {
	return m.Global
}

// IsPersonal 检查是否为 Personal 作用域
func (m *Memory) IsPersonal() bool {
	return !m.Global && m.PathID > 0
}

// GetScope 获取作用域类型字符串
func (m *Memory) GetScope() string {
	if m.Global {
		return "global"
	}
	if m.PathID > 0 {
		return "personal"
	}
	return "unknown"
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
