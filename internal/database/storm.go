package database

import (
	"os"
	"path/filepath"

	"github.com/XiaoLFeng/llm-memory/pkg/types"
	"github.com/asdine/storm/v3"
)

// DB 数据库包装结构体
// 嘿嘿~ 封装 storm.DB 让使用更优雅~ ✨
type DB struct {
	*storm.DB
}

// Open 打开数据库连接
// 自动创建目录和初始化索引哦~ 💖
func Open(path string) (*DB, error) {
	// 确保目录存在
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	// 打开数据库
	db, err := storm.Open(path)
	if err != nil {
		return nil, err
	}

	// 初始化索引
	if err := initIndexes(db); err != nil {
		db.Close()
		return nil, err
	}

	return &DB{DB: db}, nil
}

// Close 关闭数据库连接
// 记得关闭数据库释放资源呢~ (´∀｀)
func (db *DB) Close() error {
	return db.DB.Close()
}

// initIndexes 初始化所有实体的索引
// 为 Memory、Plan、Todo 三个实体创建索引~ 🎯
func initIndexes(db *storm.DB) error {
	// 初始化 Memory 索引
	if err := db.Init(&types.Memory{}); err != nil {
		return err
	}

	// 初始化 MemoryCategory 索引
	if err := db.Init(&types.MemoryCategory{}); err != nil {
		return err
	}

	// 初始化 Plan 索引
	if err := db.Init(&types.Plan{}); err != nil {
		return err
	}

	// 初始化 Todo 索引
	if err := db.Init(&types.Todo{}); err != nil {
		return err
	}

	return nil
}
