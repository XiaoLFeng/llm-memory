package types

import (
	"time"
)

// Plan 计划实体结构体 - 用于管理项目计划和任务
// 嘿嘿~ 这是一个完整的计划管理结构呢！📋
type Plan struct {
	ID          int        `storm:"id,increment"` // 主键，自增
	Title       string     `storm:"index"`        // 标题，带索引以便快速查询
	Description string     `storm:""`             // 描述，详细内容
	Status      PlanStatus `storm:"index"`        // 状态，带索引用于状态筛选
	StartDate   *time.Time `storm:""`             // 开始日期，可为空
	EndDate     *time.Time `storm:""`             // 结束日期，可为空
	Progress    int        `storm:""`             // 进度 0-100，表示完成百分比
	SubTasks    []SubTask  `storm:"inline"`       // 子任务列表，使用inline存储
	CreatedAt   time.Time  `storm:"index"`        // 创建时间，带索引
	UpdatedAt   time.Time  `storm:"index"`        // 更新时间，带索引
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
	ID          int        `storm:"id,increment"` // 子任务ID，自增
	Title       string     `storm:""`             // 子任务标题
	Description string     `storm:""`             // 子任务描述
	Status      PlanStatus `storm:""`             // 子任务状态
	Progress    int        `storm:""`             // 子任务进度 0-100
	CreatedAt   time.Time  `storm:""`             // 创建时间
	UpdatedAt   time.Time  `storm:""`             // 更新时间
}

// NewPlan 创建新的计划实例
// 💖 构造函数模式，让创建计划更优雅~
func NewPlan(title, description string) *Plan {
	now := time.Now()
	return &Plan{
		Title:       title,
		Description: description,
		Status:      PlanStatusPending,
		Progress:    0,
		SubTasks:    make([]SubTask, 0),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
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
