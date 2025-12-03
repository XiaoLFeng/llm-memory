package types

import (
	"time"
)

// Plan 计划实体结构体 - 用于管理项目计划和任务
// 嘿嘿~ 这是一个完整的计划管理结构呢！📋
type Plan struct {
	ID          int        `json:"id"`          // 主键，自增
	GroupID     int        `json:"group_id"`    // 所属组ID（0=Global）
	Path        string     `json:"path"`        // 精确路径（Personal作用域）
	Title       string     `json:"title"`       // 标题
	Description string     `json:"description"` // 描述，详细内容
	Status      PlanStatus `json:"status"`      // 状态
	StartDate   *time.Time `json:"start_date"`  // 开始日期，可为空
	EndDate     *time.Time `json:"end_date"`    // 结束日期，可为空
	Progress    int        `json:"progress"`    // 进度 0-100，表示完成百分比
	SubTasks    []SubTask  `json:"sub_tasks"`   // 子任务列表
	CreatedAt   time.Time  `json:"created_at"`  // 创建时间
	UpdatedAt   time.Time  `json:"updated_at"`  // 更新时间
}

// PlanStatus 计划状态类型
// 呀~ 用字符串类型让状态更清晰呢！🎯
type PlanStatus string

// 计划状态常量定义
// 嗯嗯！这些状态涵盖了完整的计划生命周期~
const (
	PlanStatusPending    PlanStatus = "pending"     // 待开始状态
	PlanStatusInProgress PlanStatus = "in_progress" // 进行中状态
	PlanStatusCompleted  PlanStatus = "completed"   // 已完成状态
	PlanStatusCancelled  PlanStatus = "cancelled"   // 已取消状态
)

// SubTask 子任务结构体
// 每个大计划都需要拆分成小任务来管理哦~ ✨
type SubTask struct {
	ID          int        `json:"id"`          // 子任务ID，自增
	Title       string     `json:"title"`       // 子任务标题
	Description string     `json:"description"` // 子任务描述
	Status      PlanStatus `json:"status"`      // 子任务状态
	Progress    int        `json:"progress"`    // 子任务进度 0-100
	CreatedAt   time.Time  `json:"created_at"`  // 创建时间
	UpdatedAt   time.Time  `json:"updated_at"`  // 更新时间
}

// NewPlan 创建新的计划实例
// 💖 构造函数模式，现在支持设置作用域啦~
func NewPlan(title, description string, groupID int, path string) *Plan {
	now := time.Now()
	return &Plan{
		GroupID:     groupID,
		Path:        path,
		Title:       title,
		Description: description,
		Status:      PlanStatusPending,
		Progress:    0,
		SubTasks:    make([]SubTask, 0),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// NewGlobalPlan 创建全局计划实例
func NewGlobalPlan(title, description string) *Plan {
	return NewPlan(title, description, GlobalGroupID, "")
}

// NewPersonalPlan 创建 Personal 作用域的计划实例
func NewPersonalPlan(title, description string, path string) *Plan {
	return NewPlan(title, description, GlobalGroupID, path)
}

// NewGroupPlan 创建 Group 作用域的计划实例
func NewGroupPlan(title, description string, groupID int) *Plan {
	return NewPlan(title, description, groupID, "")
}

// IsGlobal 检查计划是否为全局计划
func (p *Plan) IsGlobal() bool {
	return p.GroupID == GlobalGroupID && p.Path == ""
}

// IsPersonal 检查计划是否为 Personal 作用域
func (p *Plan) IsPersonal() bool {
	return p.Path != ""
}

// IsGroup 检查计划是否为 Group 作用域
func (p *Plan) IsGroup() bool {
	return p.GroupID != GlobalGroupID && p.Path == ""
}

// GetScope 获取计划的作用域类型
func (p *Plan) GetScope() Scope {
	if p.Path != "" {
		return ScopePersonal
	}
	if p.GroupID != GlobalGroupID {
		return ScopeGroup
	}
	return ScopeGlobal
}

// NewSubTask 创建新的子任务实例
// 为计划添加可爱的小任务~ 🍫
func NewSubTask(title, description string) SubTask {
	now := time.Now()
	return SubTask{
		Title:       title,
		Description: description,
		Status:      PlanStatusPending,
		Progress:    0,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// IsCompleted 检查计划是否已完成
// 方便的判断方法~ ＼(^o^)／
func (p *Plan) IsCompleted() bool {
	return p.Status == PlanStatusCompleted
}

// IsInProgress 检查计划是否正在进行中
// 查看计划状态的小帮手~
func (p *Plan) IsInProgress() bool {
	return p.Status == PlanStatusInProgress
}

// UpdateProgress 更新计划进度
// 智能进度管理，还会更新状态哦！🎮
func (p *Plan) UpdateProgress(progress int) {
	if progress < 0 {
		progress = 0
	}
	if progress > 100 {
		progress = 100
	}

	p.Progress = progress
	p.UpdatedAt = time.Now()

	// 根据进度自动更新状态
	if progress == 0 {
		p.Status = PlanStatusPending
	} else if progress == 100 {
		p.Status = PlanStatusCompleted
	} else {
		p.Status = PlanStatusInProgress
	}
}

// AddSubTask 添加子任务
// 为计划添加新的小任务~ ✨
func (p *Plan) AddSubTask(title, description string) {
	subTask := NewSubTask(title, description)
	p.SubTasks = append(p.SubTasks, subTask)
	p.UpdatedAt = time.Now()
}

// CalculateProgress 根据子任务计算总进度
// 智能计算整体进度，让计划管理更准确~
func (p *Plan) CalculateProgress() {
	if len(p.SubTasks) == 0 {
		p.UpdateProgress(0)
		return
	}

	totalProgress := 0
	for _, subTask := range p.SubTasks {
		totalProgress += subTask.Progress
	}

	averageProgress := totalProgress / len(p.SubTasks)
	p.UpdateProgress(averageProgress)
}
