package entity

import (
	"time"
)

// Group 组实体（数据表结构）
// 用于管理多个路径共享数据的组实体
type Group struct {
	ID          int64     `gorm:"primaryKey"` // 雪花算法生成
	Name        string    `gorm:"uniqueIndex;size:100;not null;comment:组名称"`
	Description string    `gorm:"type:text;comment:组描述"`
	CreatedAt   time.Time `gorm:"index;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`

	// 关联：路径列表
	Paths []GroupPath `gorm:"foreignKey:GroupID;constraint:OnDelete:CASCADE"`
}

// TableName 指定表名
func (Group) TableName() string {
	return "groups"
}

// GroupPath 组路径映射表
// 存储组和路径的关联关系（纯关联模式）
type GroupPath struct {
	ID             int64 `gorm:"primaryKey"`                                 // 雪花算法生成
	GroupID        int64 `gorm:"index;not null"`                             // 关联组ID
	PersonalPathID int64 `gorm:"uniqueIndex;not null;comment:关联的路径ID（全局唯一）"` // 关联 PersonalPath.ID

	// 关联：PersonalPath（用于预加载获取路径字符串）
	PersonalPath PersonalPath `gorm:"foreignKey:PersonalPathID"`
}

// TableName 指定表名
func (GroupPath) TableName() string {
	return "group_paths"
}

// GetPath 获取路径字符串（需要预加载 PersonalPath）
func (gp *GroupPath) GetPath() string {
	return gp.PersonalPath.Path
}

// GetPathIDs 获取路径 ID 列表
// 嘿嘿~ 纯关联模式下，返回 PersonalPath 的 ID 列表！💖
func (g *Group) GetPathIDs() []int64 {
	ids := make([]int64, len(g.Paths))
	for i, p := range g.Paths {
		ids[i] = p.PersonalPathID
	}
	return ids
}

// ContainsPathID 检查是否包含路径 ID
func (g *Group) ContainsPathID(pathID int64) bool {
	for _, p := range g.Paths {
		if p.PersonalPathID == pathID {
			return true
		}
	}
	return false
}

// PathCount 获取路径数量
func (g *Group) PathCount() int {
	return len(g.Paths)
}
