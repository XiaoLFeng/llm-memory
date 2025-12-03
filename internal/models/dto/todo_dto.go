package dto

import "time"

// ToDoCreateDTO 创建待办请求
type ToDoCreateDTO struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Priority    int        `json:"priority"` // 1-4，默认 2
	DueDate     *time.Time `json:"due_date"`
	Tags        []string   `json:"tags"`
	Global      bool       `json:"global"` // true=全局；false=当前路径(私有/组内)
}

// ToDoUpdateDTO 更新待办请求
type ToDoUpdateDTO struct {
	ID          int64      `json:"id"`
	Title       *string    `json:"title,omitempty"`
	Description *string    `json:"description,omitempty"`
	Priority    *int       `json:"priority,omitempty"`
	Status      *int       `json:"status,omitempty"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	Tags        *[]string  `json:"tags,omitempty"`
}

// ToDoResponseDTO 待办响应
type ToDoResponseDTO struct {
	ID          int64      `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Priority    int        `json:"priority"`
	PriorityStr string     `json:"priority_str"` // 低/中/高/紧急
	Status      int        `json:"status"`
	StatusStr   string     `json:"status_str"` // 待处理/进行中/已完成/已取消
	DueDate     *time.Time `json:"due_date"`
	CompletedAt *time.Time `json:"completed_at"`
	Tags        []string   `json:"tags"`
	Scope       string     `json:"scope"`
	IsOverdue   bool       `json:"is_overdue"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// ToDoListDTO 待办列表项
type ToDoListDTO struct {
	ID          int64      `json:"id"`
	Title       string     `json:"title"`
	Priority    int        `json:"priority"`
	PriorityStr string     `json:"priority_str"`
	Status      int        `json:"status"`
	StatusStr   string     `json:"status_str"`
	DueDate     *time.Time `json:"due_date"`
	Scope       string     `json:"scope"`
	IsOverdue   bool       `json:"is_overdue"`
}

// ToDoBatchCreateDTO 批量创建待办请求
type ToDoBatchCreateDTO struct {
	Items []ToDoCreateDTO `json:"items"` // 最多 100 个
}

// ToDoBatchUpdateDTO 批量更新待办请求
type ToDoBatchUpdateDTO struct {
	Items []ToDoUpdateDTO `json:"items"` // 最多 100 个
}

// ToDoBatchCompleteDTO 批量完成待办请求
type ToDoBatchCompleteDTO struct {
	IDs []int64 `json:"ids"` // 最多 100 个
}

// ToDoBatchDeleteDTO 批量删除待办请求
type ToDoBatchDeleteDTO struct {
	IDs []int64 `json:"ids"` // 最多 100 个
}

// ToDoBatchProgressDTO 批量更新进度请求（按状态批量更新）
type ToDoBatchProgressDTO struct {
	IDs    []int64 `json:"ids"`    // 最多 100 个
	Status int     `json:"status"` // 0-3
}

// ToDoBatchResultDTO 批量操作结果
type ToDoBatchResultDTO struct {
	Total     int      `json:"total"`     // 总数
	Succeeded int      `json:"succeeded"` // 成功数
	Failed    int      `json:"failed"`    // 失败数
	Errors    []string `json:"errors"`    // 错误信息列表
}

// 批量操作限制常量
const (
	MaxBatchSize = 100 // 单次批量操作最大数量
)
