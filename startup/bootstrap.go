package startup

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/XiaoLFeng/llm-memory/internal/app"
	"github.com/XiaoLFeng/llm-memory/internal/database"
	"github.com/XiaoLFeng/llm-memory/internal/repository"
	"github.com/XiaoLFeng/llm-memory/internal/service"
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
	db *database.DB

	// Service 层（公开，供外部使用）
	MemoryService *service.MemoryService
	PlanService   *service.PlanService
	TodoService   *service.TodoService

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
// 顺序：Context -> Config -> Database -> Repository -> Service
func (b *Bootstrap) Initialize(ctx context.Context) error {
	if b.initialized {
		return ErrAlreadyInitialized
	}

	// 1. 创建应用级 Context
	b.appCtx = NewAppContext(ctx)

	// 2. 加载配置
	config, err := b.loadConfig()
	if err != nil {
		return fmt.Errorf("加载配置失败: %w", err)
	}
	b.config = config

	// 3. 初始化数据库
	db, err := database.Open(config.DBPath)
	if err != nil {
		return fmt.Errorf("初始化数据库失败: %w", err)
	}
	b.db = db

	// 4. 创建 Repository 实例
	memoryRepo := repository.NewMemoryRepo(db)
	planRepo := repository.NewPlanRepo(db)
	todoRepo := repository.NewTodoRepo(db)

	// 5. 创建 Service 实例
	b.MemoryService = service.NewMemoryService(memoryRepo)
	b.PlanService = service.NewPlanService(planRepo)
	b.TodoService = service.NewTodoService(todoRepo)

	// 6. 启动信号处理
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

// DB 获取数据库实例
func (b *Bootstrap) DB() *database.DB {
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

	// 关闭数据库
	if b.db != nil {
		if err := b.db.Close(); err != nil {
			return fmt.Errorf("关闭数据库失败: %w", err)
		}
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
