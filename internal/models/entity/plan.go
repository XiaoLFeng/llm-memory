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
// 纯关联模式：PathID=0 表示 Global，PathID>0 关联 PersonalPath
type Plan struct {
	ID          int64      `gorm:"primaryKey"`                             // 雪花算法生成
	PathID      int64      `gorm:"index;default:0;comment:路径ID（0=Global）"` // 关联 PersonalPath.ID，0 表示 Global
	Title       string     `gorm:"index;size:255;not null;comment:标题"`
	Description string     `gorm:"type:text;comment:简要描述（摘要）"`
	Content     string     `gorm:"type:text;comment:详细内容（新增字段）"` // 新增：详细内容
	Status      PlanStatus `gorm:"index;size:20;default:'pending'"`
	StartDate   *time.Time `gorm:"comment:开始日期"`
	EndDate     *time.Time `gorm:"comment:结束日期"`
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

// IsGlobal 检查是否为全局计划
func (p *Plan) IsGlobal() bool {
	return p.PathID == 0
}

// IsPersonal 检查是否为 Personal 作用域
// 纯关联模式下，PathID > 0 表示关联某个路径
func (p *Plan) IsPersonal() bool {
	return p.PathID > 0
}

// GetScope 获取作用域类型字符串
// 注意：纯关联模式下只有 personal 和 global，group 通过 join 查询实现
func (p *Plan) GetScope() string {
	if p.PathID > 0 {
		return "personal"
	}
	return "global"
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
		now := time.Now()
		p.EndDate = &now
	} else {
		p.Status = PlanStatusInProgress
	}
}

// Start 开始计划
func (p *Plan) Start() {
	now := time.Now()
	p.StartDate = &now
	p.Status = PlanStatusInProgress
	if p.Progress == 0 {
		p.Progress = 1
	}
}

// Complete 完成计划
func (p *Plan) Complete() {
	now := time.Now()
	p.EndDate = &now
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
