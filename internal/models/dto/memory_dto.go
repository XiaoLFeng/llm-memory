package dto

import "time"

// MemoryCreateDTO 创建记忆请求
type MemoryCreateDTO struct {
	Title    string   `json:"title"`
	Content  string   `json:"content"`
	Category string   `json:"category"`
	Tags     []string `json:"tags"`
	Priority int      `json:"priority"`
	Global   bool     `json:"global"` // true=全局；false=当前路径(私有/组内)
}

// MemoryUpdateDTO 更新记忆请求
type MemoryUpdateDTO struct {
	ID       int64     `json:"id"`
	Title    *string   `json:"title,omitempty"`
	Content  *string   `json:"content,omitempty"`
	Category *string   `json:"category,omitempty"`
	Tags     *[]string `json:"tags,omitempty"`
	Priority *int      `json:"priority,omitempty"`
}

// MemoryResponseDTO 记忆响应
type MemoryResponseDTO struct {
	ID         int64     `json:"id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	Category   string    `json:"category"`
	Tags       []string  `json:"tags"`
	Priority   int       `json:"priority"`
	Scope      string    `json:"scope"` // Personal/Group/Global
	IsArchived bool      `json:"is_archived"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// MemoryListDTO 记忆列表项
type MemoryListDTO struct {
	ID         int64  `json:"id"`
	Title      string `json:"title"`
	Category   string `json:"category"`
	Priority   int    `json:"priority"`
	IsArchived bool   `json:"is_archived"`
	Scope      string `json:"scope"`
}

// MemorySearchDTO 记忆搜索请求
type MemorySearchDTO struct {
	Keyword string `json:"keyword"`
	Scope   string `json:"scope"` // personal/group/global/all
}
