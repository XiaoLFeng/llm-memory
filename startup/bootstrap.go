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
	"github.com/XiaoLFeng/llm-memory/pkg/types"
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

	// 数据库路径（不再持有长连接）
	dbPath string

	// Service 层（公开，供外部使用）
	MemoryService *service.MemoryService
	PlanService   *service.PlanService
	TodoService   *service.TodoService
	GroupService  *service.GroupService // 新增：组服务

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
	b.dbPath = config.DBPath

	// 3. 初始化数据库（仅用于确保索引创建，立即关闭）
	// 嘿嘿~ 每次操作都会自己打开关闭连接，这里只是确保表结构！💖
	db, err := database.Open(config.DBPath)
	if err != nil {
		return fmt.Errorf("初始化数据库失败: %w", err)
	}
	db.Close() // 立即关闭，不保持长连接

	// 4. 创建 Repository 实例（传入 dbPath）
	memoryRepo := repository.NewMemoryRepo(config.DBPath)
	planRepo := repository.NewPlanRepo(config.DBPath)
	todoRepo := repository.NewTodoRepo(config.DBPath)
	groupRepo := repository.NewGroupRepo(config.DBPath) // 新增：组仓储

	// 5. 创建 Service 实例
	b.MemoryService = service.NewMemoryService(memoryRepo)
	b.PlanService = service.NewPlanService(planRepo)
	b.TodoService = service.NewTodoService(todoRepo)
	b.GroupService = service.NewGroupService(groupRepo) // 新增：组服务

	// 6. 解析当前作用域
	// 嘿嘿~ 启动时自动获取当前目录的作用域上下文！💖
	scope, err := b.GroupService.GetCurrentScope(b.appCtx.Context())
	if err != nil {
		// 如果解析失败，使用仅包含 Global 的作用域
		scope = types.NewGlobalOnlyScope()
	}
	b.CurrentScope = scope

	// 7. 启动信号处理
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

// DBPath 获取数据库路径
// 嘿嘿~ 现在不再持有长连接，只提供路径！💖
func (b *Bootstrap) DBPath() string {
	return b.dbPath
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

	// 数据库连接现在由每次操作自己管理，不需要在这里关闭

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
