package database

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 嘿嘿~ 这是全新的 SQLite + GORM 数据库模块！(/)
// 用来替代旧的 Storm 数据库，支持并发访问哦~ 💖

var (
	gormDB   *gorm.DB
	gormOnce sync.Once
	gormErr  error
)

// SQLiteConfig SQLite 数据库配置
type SQLiteConfig struct {
	DBPath string // 数据库文件路径
	Debug  bool   // 是否开启调试模式
}

// OpenSQLite 打开 SQLite 数据库连接（单例模式）
func OpenSQLite(cfg *SQLiteConfig) (*gorm.DB, error) {
	gormOnce.Do(func() {
		// 确保目录存在
		dir := filepath.Dir(cfg.DBPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			gormErr = err
			return
		}

		// 配置 GORM Logger
		logLevel := logger.Silent
		if cfg.Debug {
			logLevel = logger.Info
		}

		// 打开 SQLite 连接
		// 启用 WAL 模式支持并发读写，设置忙等待超时
		dsn := cfg.DBPath + "?_journal_mode=WAL&_busy_timeout=5000&_foreign_keys=ON"
		conn, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logLevel),
		})
		if err != nil {
			gormErr = err
			return
		}

		// 获取底层 sql.DB 设置连接池
		// SQLite 推荐单连接，避免锁争用
		sqlDB, err := conn.DB()
		if err != nil {
			gormErr = err
			return
		}
		sqlDB.SetMaxOpenConns(1)
		sqlDB.SetMaxIdleConns(1)

		gormDB = conn
	})

	if gormErr != nil {
		return nil, gormErr
	}

	return gormDB, nil
}

// GetGormDB 获取 GORM 数据库实例
// 嘿嘿~ 方便在其他地方获取数据库连接！💖
func GetGormDB() *gorm.DB {
	return gormDB
}

// CloseSQLite 关闭 SQLite 数据库连接
// 优雅地关闭数据库连接~ (^o^)/
func CloseSQLite() error {
	if gormDB != nil {
		sqlDB, err := gormDB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

// ResetSQLiteConnection 重置数据库连接（仅用于测试）
// 呀~ 这个方法只在测试时使用哦！
func ResetSQLiteConnection() {
	gormDB = nil
	gormOnce = sync.Once{}
	gormErr = nil
}

// AutoMigrateSQLite 自动迁移数据库表结构
// 嘿嘿~ 自动创建和更新表结构，方便管理！✨
func AutoMigrateSQLite(db *gorm.DB, models ...interface{}) error {
	return db.AutoMigrate(models...)
}
