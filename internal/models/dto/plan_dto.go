package dto

import "time"

// PlanCreateDTO 创建计划请求
type PlanCreateDTO struct {
	Title       string `json:"title"`
	Description string `json:"description"` // 摘要
	Content     string `json:"content"`     // 详细内容
	Scope       string `json:"scope"`       // personal/group/global（默认 group）
}

// PlanUpdateDTO 更新计划请求
type PlanUpdateDTO struct {
	ID          int64      `json:"id"`
	Title       *string    `json:"title,omitempty"`
	Description *string    `json:"description,omitempty"`
	Content     *string    `json:"content,omitempty"`
	Status      *string    `json:"status,omitempty"`
	Progress    *int       `json:"progress,omitempty"`
	StartDate   *time.Time `json:"start_date,omitempty"`
	EndDate     *time.Time `json:"end_date,omitempty"`
}

// PlanProgressDTO 更新计划进度请求
type PlanProgressDTO struct {
	ID       int64 `json:"id"`
	Progress int   `json:"progress"` // 0-100
}

// PlanResponseDTO 计划响应
type PlanResponseDTO struct {
	ID          int64        `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description"` // 摘要
	Content     string       `json:"content"`     // 详细内容
	Status      string       `json:"status"`
	StatusStr   string       `json:"status_str"` // 状态显示文本
	Progress    int          `json:"progress"`
	StartDate   *time.Time   `json:"start_date"`
	EndDate     *time.Time   `json:"end_date"`
	SubTasks    []SubTaskDTO `json:"sub_tasks"`
	Scope       string       `json:"scope"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

// PlanListDTO 计划列表项
type PlanListDTO struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"` // 摘要用于列表展示
	Status      string `json:"status"`
	Progress    int    `json:"progress"`
	Scope       string `json:"scope"`
}

// SubTaskDTO 子任务 DTO
type SubTaskDTO struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Progress    int       `json:"progress"`
	SortOrder   int       `json:"sort_order"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// SubTaskCreateDTO 创建子任务请求
type SubTaskCreateDTO struct {
	PlanID      int64  `json:"plan_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

// SubTaskUpdateDTO 更新子任务请求
type SubTaskUpdateDTO struct {
	ID          int64   `json:"id"`
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Status      *string `json:"status,omitempty"`
	Progress    *int    `json:"progress,omitempty"`
}
