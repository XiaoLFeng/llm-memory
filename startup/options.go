package startup

import "time"

// Options 启动选项
// 呀~ 配置启动行为的各种选项！✨
type Options struct {
	// ConfigPath 配置文件路径（可选，为空时使用默认路径）
	ConfigPath string

	// ShutdownTimeout 优雅关闭超时时间
	ShutdownTimeout time.Duration

	// EnableSignalHandler 是否启用信号处理
	EnableSignalHandler bool

	// Debug 调试模式
	Debug bool
}

// DefaultOptions 返回默认选项
// 嘿嘿~ 合理的默认值让使用更简单！(´∀｀)
func DefaultOptions() *Options {
	return &Options{
		ConfigPath:          "",
		ShutdownTimeout:     30 * time.Second,
		EnableSignalHandler: true,
		Debug:               false,
	}
}

// Option 选项函数类型
type Option func(*Options)

// WithConfigPath 设置配置文件路径
func WithConfigPath(path string) Option {
	return func(o *Options) {
		o.ConfigPath = path
	}
}

// WithShutdownTimeout 设置关闭超时时间
func WithShutdownTimeout(timeout time.Duration) Option {
	return func(o *Options) {
		o.ShutdownTimeout = timeout
	}
}

// WithSignalHandler 启用/禁用信号处理
func WithSignalHandler(enabled bool) Option {
	return func(o *Options) {
		o.EnableSignalHandler = enabled
	}
}

// WithDebug 设置调试模式
func WithDebug(debug bool) Option {
	return func(o *Options) {
		o.Debug = debug
	}
}
