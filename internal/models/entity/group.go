package entity

import (
	"time"

	"gorm.io/gorm"
)

// Group 组实体（数据表结构）
// 用于管理多个路径共享数据的组实体
type Group struct {
	ID          int64          `gorm:"primaryKey"` // 雪花算法生成
	Name        string         `gorm:"uniqueIndex;size:100;not null;comment:组名称"`
	Description string         `gorm:"type:text;comment:组描述"`
	CreatedAt   time.Time      `gorm:"index;autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"` // 软删除支持

	// 关联：路径列表
	Paths []GroupPath `gorm:"foreignKey:GroupID;constraint:OnDelete:CASCADE"`
}

// TableName 指定表名
func (Group) TableName() string {
	return "groups"
}

// GroupPath 组路径映射表
// 存储组和路径的关联关系
type GroupPath struct {
	ID      int64  `gorm:"primaryKey"`     // 雪花算法生成
	GroupID int64  `gorm:"index;not null"` // 关联组ID
	Path    string `gorm:"uniqueIndex;size:1024;not null;comment:路径（全局唯一）"`
}

// TableName 指定表名
func (GroupPath) TableName() string {
	return "group_paths"
}

// GetPathStrings 获取路径字符串列表
func (g *Group) GetPathStrings() []string {
	paths := make([]string, len(g.Paths))
	for i, p := range g.Paths {
		paths[i] = p.Path
	}
	return paths
}

// AddPath 添加路径（返回 false 表示已存在）
func (g *Group) AddPath(path string) bool {
	for _, p := range g.Paths {
		if p.Path == path {
			return false
		}
	}
	g.Paths = append(g.Paths, GroupPath{Path: path})
	return true
}

// RemovePath 移除路径（返回 false 表示不存在）
func (g *Group) RemovePath(path string) bool {
	for i, p := range g.Paths {
		if p.Path == path {
			g.Paths = append(g.Paths[:i], g.Paths[i+1:]...)
			return true
		}
	}
	return false
}

// ContainsPath 检查是否包含路径
func (g *Group) ContainsPath(path string) bool {
	for _, p := range g.Paths {
		if p.Path == path {
			return true
		}
	}
	return false
}

// PathCount 获取路径数量
func (g *Group) PathCount() int {
	return len(g.Paths)
}
