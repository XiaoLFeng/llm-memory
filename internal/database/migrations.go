package database

import (
	"gorm.io/gorm"
)

// RunMigrations 执行数据库迁移
// 嘿嘿~ 在 AutoMigrate 之前运行，处理表重命名等特殊迁移！💖
func RunMigrations(db *gorm.DB) error {
	// 执行表重命名迁移
	if err := renameToDoTables(db); err != nil {
		return err
	}
	return nil
}

// renameToDoTables 重命名 ToDo 相关表
// 呀~ 把 to_dos 改成 todos，to_do_tags 改成 todo_tags！✨
func renameToDoTables(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// 检查旧表 to_dos 是否存在
		var todoCount int64
		tx.Raw("SELECT count(*) FROM sqlite_master WHERE type='table' AND name='to_dos'").Scan(&todoCount)

		if todoCount > 0 {
			// 重命名 to_dos -> todos
			if err := tx.Exec("ALTER TABLE to_dos RENAME TO todos").Error; err != nil {
				return err
			}
		}

		// 检查旧表 to_do_tags 是否存在
		var tagCount int64
		tx.Raw("SELECT count(*) FROM sqlite_master WHERE type='table' AND name='to_do_tags'").Scan(&tagCount)

		if tagCount > 0 {
			// 重命名 to_do_tags -> todo_tags
			if err := tx.Exec("ALTER TABLE to_do_tags RENAME TO todo_tags").Error; err != nil {
				return err
			}
		}

		return nil
	})
}
