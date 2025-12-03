package entity

import (
	"time"
)

// PersonalPath 个人路径索引实体
// 记录用户访问过的路径，用于自动创建 Personal 作用域索引
type PersonalPath struct {
	ID        int64     `gorm:"primaryKey"`                     // 雪花算法生成
	Path      string    `gorm:"uniqueIndex;size:1024;not null"` // 路径（全局唯一）
	LastVisit time.Time `gorm:"index"`                          // 最后访问时间
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// TableName 指定表名
func (PersonalPath) TableName() string {
	return "personal_paths"
}

// Touch 更新最后访问时间
func (p *PersonalPath) Touch() {
	p.LastVisit = time.Now()
}
