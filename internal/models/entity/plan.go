package entity

import (
	"time"
)

// PlanStatus 计划状态类型
// 用字符串类型让状态更清晰
type PlanStatus string

// 计划状态常量定义
// 这些状态涵盖了完整的计划生命周期
const (
	PlanStatusPending    PlanStatus = "pending"     // 待开始状态
	PlanStatusInProgress PlanStatus = "in_progress" // 进行中状态
	PlanStatusCompleted  PlanStatus = "completed"   // 已完成状态
	PlanStatusCancelled  PlanStatus = "cancelled"   // 已取消状态
)

// Plan 计划实体（数据表结构）
// 用于跟踪长期目标和复杂任务的计划实体
// PathID 关联 PersonalPath 或 Group 中的路径
type Plan struct {
	ID          int64      `gorm:"primaryKey"`                                 // 雪花算法生成
	Code        string     `gorm:"index;size:100;not null;comment:人类可读的唯一标识码"` // 外部查询标识，活跃状态唯一
	PathID      int64      `gorm:"index;not null;comment:路径ID（关联个人或小组路径）"`     // 关联 Path.ID
	Title       string     `gorm:"index;size:255;not null;comment:标题"`
	Description string     `gorm:"type:text;not null;comment:简要描述（摘要）"`
	Content     string     `gorm:"type:text;not null;comment:详细内容"`
	Status      PlanStatus `gorm:"index;size:20;default:'pending'"`
	Progress    int        `gorm:"default:0;comment:进度 0-100"`
	CreatedAt   time.Time  `gorm:"index;autoCreateTime"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime"`

	// 关联：子任务（独立存储，不再 inline）
	SubTasks []SubTask `gorm:"foreignKey:PlanID;constraint:OnDelete:CASCADE"`
}

// TableName 指定表名
func (Plan) TableName() string {
	return "plans"
}

// IsPersonal 检查是否为 Personal 作用域
// PathID > 0 表示关联某个路径（个人或小组）
func (p *Plan) IsPersonal() bool {
	return p.PathID > 0
}

// GetScope 获取作用域类型字符串
// 注意：group 通过上下文判断，此处只能返回 personal 或 unknown
func (p *Plan) GetScope() string {
	if p.PathID > 0 {
		return "personal"
	}
	return "unknown"
}

// IsCompleted 检查计划是否已完成
func (p *Plan) IsCompleted() bool {
	return p.Status == PlanStatusCompleted
}

// IsInProgress 检查计划是否正在进行中
func (p *Plan) IsInProgress() bool {
	return p.Status == PlanStatusInProgress
}

// UpdateProgress 更新计划进度（自动调整状态）
// 智能进度管理，还会更新状态
func (p *Plan) UpdateProgress(progress int) {
	if progress < 0 {
		progress = 0
	}
	if progress > 100 {
		progress = 100
	}

	p.Progress = progress

	// 根据进度自动更新状态
	if progress == 0 {
		p.Status = PlanStatusPending
	} else if progress == 100 {
		p.Status = PlanStatusCompleted
	} else {
		p.Status = PlanStatusInProgress
	}
}

// Start 开始计划
func (p *Plan) Start() {
	p.Status = PlanStatusInProgress
	if p.Progress == 0 {
		p.Progress = 1
	}
}

// Complete 完成计划
func (p *Plan) Complete() {
	p.Status = PlanStatusCompleted
	p.Progress = 100
}

// Cancel 取消计划
func (p *Plan) Cancel() {
	p.Status = PlanStatusCancelled
}

// CalculateProgress 根据子任务计算总进度
// 智能计算整体进度，让计划管理更准确
func (p *Plan) CalculateProgress() {
	if len(p.SubTasks) == 0 {
		return
	}

	totalProgress := 0
	for _, subTask := range p.SubTasks {
		totalProgress += subTask.Progress
	}

	averageProgress := totalProgress / len(p.SubTasks)
	p.UpdateProgress(averageProgress)
}
