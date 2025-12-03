package dto

import "time"

// PlanCreateDTO 创建计划请求
// 嘿嘿~ 用于创建新计划的数据传输对象！💖
type PlanCreateDTO struct {
	Title       string `json:"title"`
	Description string `json:"description"` // 摘要
	Content     string `json:"content"`     // 详细内容（新增）
	Scope       string `json:"scope"`       // personal/group/global（默认 group）
}

// PlanUpdateDTO 更新计划请求
// 呀~ 用于更新已有计划的数据传输对象！✨
type PlanUpdateDTO struct {
	ID          uint       `json:"id"`
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
	ID       uint `json:"id"`
	Progress int  `json:"progress"` // 0-100
}

// PlanResponseDTO 计划响应
// 嘿嘿~ 用于返回计划详情的数据传输对象！💖
type PlanResponseDTO struct {
	ID          uint         `json:"id"`
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
// 呀~ 用于列表展示的简化计划数据！✨
type PlanListDTO struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"` // 摘要用于列表展示
	Status      string `json:"status"`
	Progress    int    `json:"progress"`
	Scope       string `json:"scope"`
}

// SubTaskDTO 子任务 DTO
type SubTaskDTO struct {
	ID          uint      `json:"id"`
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
	PlanID      uint   `json:"plan_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

// SubTaskUpdateDTO 更新子任务请求
type SubTaskUpdateDTO struct {
	ID          uint    `json:"id"`
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Status      *string `json:"status,omitempty"`
	Progress    *int    `json:"progress,omitempty"`
}
