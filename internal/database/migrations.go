package database

import "gorm.io/gorm"

// RunMigrations 执行数据库迁移
// 全新架构：无需兼容历史表结构，交给 AutoMigrate 直接创建/更新
func RunMigrations(db *gorm.DB) error {
	return nil
}
