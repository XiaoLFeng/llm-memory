package dto

import "time"

// GroupCreateDTO 创建组请求
// 嘿嘿~ 用于创建新组的数据传输对象！💖
type GroupCreateDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// GroupUpdateDTO 更新组请求
type GroupUpdateDTO struct {
	ID          int64   `json:"id"`
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

// GroupAddPathDTO 添加路径到组请求
type GroupAddPathDTO struct {
	GroupName string `json:"group_name"`
	Path      string `json:"path"` // 留空则添加当前目录
}

// GroupRemovePathDTO 从组移除路径请求
type GroupRemovePathDTO struct {
	GroupName string `json:"group_name"`
	Path      string `json:"path"`
}

// GroupResponseDTO 组响应
type GroupResponseDTO struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Paths       []string  `json:"paths"`
	PathCount   int       `json:"path_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// GroupListDTO 组列表项
type GroupListDTO struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	PathCount   int    `json:"path_count"`
}

// ScopeInfoDTO 当前作用域信息
type ScopeInfoDTO struct {
	CurrentPath string `json:"current_path"`
	GroupID     int64  `json:"group_id"`
	GroupName   string `json:"group_name"`
	IsInGroup   bool   `json:"is_in_group"`
}
