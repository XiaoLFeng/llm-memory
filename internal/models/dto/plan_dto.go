package dto

import "time"

// PlanCreateDTO 创建计划请求
type PlanCreateDTO struct {
	Code        string `json:"code"` // 人类可读的唯一标识码（必填）
	Title       string `json:"title"`
	Description string `json:"description"` // 摘要
	Content     string `json:"content"`     // 详细内容
}

// PlanUpdateDTO 更新计划请求
type PlanUpdateDTO struct {
	Code        string  `json:"code"` // 通过 code 定位计划
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Content     *string `json:"content,omitempty"`
	Progress    *int    `json:"progress,omitempty"`
}

// PlanProgressDTO 更新计划进度请求
type PlanProgressDTO struct {
	Code     string `json:"code"`     // 通过 code 定位计划
	Progress int    `json:"progress"` // 0-100
}

// PlanResponseDTO 计划响应
type PlanResponseDTO struct {
	ID          int64         `json:"id"`
	Code        string        `json:"code"` // 人类可读的唯一标识码
	Title       string        `json:"title"`
	Description string        `json:"description"` // 摘要
	Content     string        `json:"content"`     // 详细内容
	Status      string        `json:"status"`
	StatusStr   string        `json:"status_str"` // 状态显示文本
	Progress    int           `json:"progress"`
	Todos       []ToDoListDTO `json:"todos"` // 关联的待办事项列表
	Scope       string        `json:"scope"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

// PlanListDTO 计划列表项
type PlanListDTO struct {
	ID          int64  `json:"id"`
	Code        string `json:"code"` // 人类可读的唯一标识码
	Title       string `json:"title"`
	Description string `json:"description"` // 摘要用于列表展示
	Status      string `json:"status"`
	Progress    int    `json:"progress"`
	TodoCount   int    `json:"todo_count"` // 待办数量
	Scope       string `json:"scope"`
}
