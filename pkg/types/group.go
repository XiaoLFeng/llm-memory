package types

import (
	"time"
)

// Group 组实体结构体
// 嘿嘿~ 用于管理多个路径共享同一套记忆、计划、待办！📦
type Group struct {
	ID          int       `json:"id"`          // 主键，自增
	Name        string    `json:"name"`        // 组名称，唯一
	Description string    `json:"description"` // 组描述
	Paths       []string  `json:"paths"`       // 关联的路径列表
	CreatedAt   time.Time `json:"created_at"`  // 创建时间
	UpdatedAt   time.Time `json:"updated_at"`  // 更新时间
}

// GroupPathMapping 路径到组的映射
// 呀~ 用于快速查找路径属于哪个组呢！🔍
type GroupPathMapping struct {
	ID      int    `json:"id"`       // 主键
	Path    string `json:"path"`     // 路径（唯一索引）
	GroupID int    `json:"group_id"` // 所属组ID
}

// NewGroup 创建新的组实例
// 💖 构造函数模式，让创建组更优雅~
func NewGroup(name, description string) *Group {
	now := time.Now()
	return &Group{
		Name:        name,
		Description: description,
		Paths:       make([]string, 0),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// AddPath 添加路径到组
// 返回 false 表示路径已存在
func (g *Group) AddPath(path string) bool {
	// 检查路径是否已存在
	for _, p := range g.Paths {
		if p == path {
			return false
		}
	}
	g.Paths = append(g.Paths, path)
	g.UpdatedAt = time.Now()
	return true
}

// RemovePath 从组中移除路径
// 返回 false 表示路径不存在
func (g *Group) RemovePath(path string) bool {
	for i, p := range g.Paths {
		if p == path {
			g.Paths = append(g.Paths[:i], g.Paths[i+1:]...)
			g.UpdatedAt = time.Now()
			return true
		}
	}
	return false
}

// ContainsPath 检查组是否包含指定路径
func (g *Group) ContainsPath(path string) bool {
	for _, p := range g.Paths {
		if p == path {
			return true
		}
	}
	return false
}

// PathCount 返回组中的路径数量
func (g *Group) PathCount() int {
	return len(g.Paths)
}

// BeforeUpdate 在更新前自动设置更新时间
func (g *Group) BeforeUpdate() error {
	g.UpdatedAt = time.Now()
	return nil
}

// NewGroupPathMapping 创建新的路径映射
func NewGroupPathMapping(path string, groupID int) *GroupPathMapping {
	return &GroupPathMapping{
		Path:    path,
		GroupID: groupID,
	}
}
