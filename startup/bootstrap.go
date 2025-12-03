package startup

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/XiaoLFeng/llm-memory/internal/app"
	"github.com/XiaoLFeng/llm-memory/internal/database"
	"github.com/XiaoLFeng/llm-memory/internal/models"
	"github.com/XiaoLFeng/llm-memory/internal/models/entity"
	"github.com/XiaoLFeng/llm-memory/internal/service"
	"github.com/XiaoLFeng/llm-memory/pkg/types"
	"gorm.io/gorm"
)

// 错误定义
var (
	ErrAlreadyInitialized = errors.New("bootstrap 已经初始化")
	ErrNotInitialized     = errors.New("bootstrap 尚未初始化")
)

// Bootstrap 应用启动器
// 嘿嘿~ 这是统一的应用启动入口！(´∀｀)💖
// 负责初始化配置、数据库、服务等所有组件~
type Bootstrap struct {
	// Context 管理
	appCtx *AppContext

	// 配置
	config  *app.Config
	options *Options

	// 数据库
	db *gorm.DB

	// Service 层（公开，供外部使用）
	MemoryService *service.MemoryService
	PlanService   *service.PlanService
	ToDoService   *service.ToDoService  // 注意：类型名使用 ToDo
	GroupService  *service.GroupService // 组服务

	// 当前作用域上下文
	// 嘿嘿~ 启动时自动解析当前目录的作用域！✨
	CurrentScope *types.ScopeContext

	// 信号处理
	signalHandler *SignalHandler

	// 状态
	initialized bool
}

// New 创建新的 Bootstrap 实例
// 呀~ 只是创建实例，还没有初始化哦！✨
func New(opts ...Option) *Bootstrap {
	options := DefaultOptions()
	for _, opt := range opts {
		opt(options)
	}

	return &Bootstrap{
		options: options,
	}
}

// Initialize 初始化应用
// 嘿嘿~ 按照正确的顺序初始化所有组件！💫
// 顺序：Snowflake -> Context -> Config -> Database -> Model -> Service
func (b *Bootstrap) Initialize(ctx context.Context) error {
	if b.initialized {
		return ErrAlreadyInitialized
	}

	// 0. 初始化雪花算法
	// 嘿嘿~ 节点 ID 基于机器 MAC 地址或 hostname 自动生成！✨
	if err := database.InitSnowflake(); err != nil {
		return fmt.Errorf("初始化雪花算法失败: %w", err)
	}

	// 1. 创建应用级 Context
	b.appCtx = NewAppContext(ctx)

	// 2. 加载配置
	config, err := b.loadConfig()
	if err != nil {
		return fmt.Errorf("加载配置失败: %w", err)
	}
	b.config = config

	// 3. 初始化 GORM 数据库
	// 嘿嘿~ 使用 SQLite + WAL 模式支持并发读写！💖
	gormDB, err := database.OpenSQLite(&database.SQLiteConfig{
		DBPath: config.DBPath,
		Debug:  config.Debug,
	})
	if err != nil {
		return fmt.Errorf("初始化数据库失败: %w", err)
	}
	b.db = gormDB

	// 4. 自动迁移表结构
	// 呀~ 确保数据库表结构是最新的！✨
	if err := database.AutoMigrateSQLite(gormDB,
		&entity.Memory{},
		&entity.MemoryTag{},
		&entity.Plan{},
		&entity.SubTask{},
		&entity.ToDo{},
		&entity.ToDoTag{},
		&entity.Group{},
		&entity.GroupPath{},
		&entity.PersonalPath{},
	); err != nil {
		return fmt.Errorf("迁移数据库表结构失败: %w", err)
	}

	// 5. 创建 Model 实例
	memoryModel := models.NewMemoryModel(gormDB)
	planModel := models.NewPlanModel(gormDB)
	todoModel := models.NewToDoModel(gormDB)
	groupModel := models.NewGroupModel(gormDB)

	// 6. 创建 Service 实例
	b.MemoryService = service.NewMemoryService(memoryModel)
	b.PlanService = service.NewPlanService(planModel)
	b.ToDoService = service.NewToDoService(todoModel)
	b.GroupService = service.NewGroupService(groupModel)

	// 7. 解析当前作用域
	// 嘿嘿~ 启动时自动获取当前目录的作用域上下文！💖
	scope, err := b.GroupService.GetCurrentScope(b.appCtx.Context())
	if err != nil {
		// 如果解析失败，使用仅包含 Global 的作用域
		scope = types.NewGlobalOnlyScope()
	}
	b.CurrentScope = scope

	// 8. 启动信号处理
	if b.options.EnableSignalHandler {
		b.signalHandler = NewSignalHandler()
		b.signalHandler.Start(func(sig os.Signal) {
			fmt.Printf("\n收到信号 %v，正在优雅关闭...\n", sig)
			_ = b.Shutdown()
		})
	}

	b.initialized = true
	return nil
}

// loadConfig 加载配置
func (b *Bootstrap) loadConfig() (*app.Config, error) {
	// TODO: 支持从指定路径加载配置
	return app.LoadConfig()
}

// Context 获取应用级 Context
// 呀~ 可以传递给 Service 和其他组件使用！💖
func (b *Bootstrap) Context() context.Context {
	if b.appCtx == nil {
		return context.Background()
	}
	return b.appCtx.Context()
}

// AppContext 获取 AppContext 实例
func (b *Bootstrap) AppContext() *AppContext {
	return b.appCtx
}

// Config 获取配置
func (b *Bootstrap) Config() *app.Config {
	return b.config
}

// DB 获取 GORM 数据库实例
// 嘿嘿~ 现在使用 GORM 管理数据库连接！💖
func (b *Bootstrap) DB() *gorm.DB {
	return b.db
}

// Shutdown 优雅关闭
// 嘿嘿~ 按照逆序关闭所有组件！✨
func (b *Bootstrap) Shutdown() error {
	if !b.initialized {
		return ErrNotInitialized
	}

	// 停止信号处理
	if b.signalHandler != nil {
		b.signalHandler.Stop()
	}

	// 关闭 AppContext（会等待所有 goroutine）
	if b.appCtx != nil {
		if err := b.appCtx.Shutdown(b.options.ShutdownTimeout); err != nil {
			fmt.Printf("等待任务完成超时: %v\n", err)
		}
	}

	// 关闭数据库连接
	if err := database.CloseSQLite(); err != nil {
		fmt.Printf("关闭数据库连接失败: %v\n", err)
	}

	b.initialized = false
	return nil
}

// MustInitialize 初始化应用（失败时退出）
// 嘿嘿~ 简化启动代码，失败直接退出！💫
func (b *Bootstrap) MustInitialize(ctx context.Context) *Bootstrap {
	if err := b.Initialize(ctx); err != nil {
		fmt.Printf("初始化应用失败: %v\n", err)
		os.Exit(1)
	}
	return b
}

// IsInitialized 检查是否已初始化
func (b *Bootstrap) IsInitialized() bool {
	return b.initialized
}
