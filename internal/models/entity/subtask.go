package entity

import (
	"time"
)

// SubTask 子任务实体（独立表）
// 嘿嘿~ 每个大计划都需要拆分成小任务来管理哦~ ✨
type SubTask struct {
	ID          uint       `gorm:"primaryKey;autoIncrement"`
	PlanID      uint       `gorm:"index;not null;comment:所属计划ID"`
	Title       string     `gorm:"size:255;not null"`
	Description string     `gorm:"type:text"`
	Status      PlanStatus `gorm:"size:20;default:'pending'"`
	Progress    int        `gorm:"default:0"`
	SortOrder   int        `gorm:"default:0;comment:排序顺序"`
	CreatedAt   time.Time  `gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime"`
}

// TableName 指定表名
func (SubTask) TableName() string {
	return "sub_tasks"
}

// IsCompleted 检查子任务是否已完成
func (s *SubTask) IsCompleted() bool {
	return s.Status == PlanStatusCompleted
}

// IsInProgress 检查子任务是否正在进行中
func (s *SubTask) IsInProgress() bool {
	return s.Status == PlanStatusInProgress
}

// UpdateProgress 更新子任务进度（自动调整状态）
func (s *SubTask) UpdateProgress(progress int) {
	if progress < 0 {
		progress = 0
	}
	if progress > 100 {
		progress = 100
	}

	s.Progress = progress

	// 根据进度自动更新状态
	if progress == 0 {
		s.Status = PlanStatusPending
	} else if progress == 100 {
		s.Status = PlanStatusCompleted
	} else {
		s.Status = PlanStatusInProgress
	}
}

// Complete 完成子任务
func (s *SubTask) Complete() {
	s.Status = PlanStatusCompleted
	s.Progress = 100
}

// Cancel 取消子任务
func (s *SubTask) Cancel() {
	s.Status = PlanStatusCancelled
}
