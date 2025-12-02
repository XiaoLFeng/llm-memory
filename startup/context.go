package startup

import (
	"context"
	"errors"
	"sync"
	"time"
)

// ErrShutdownTimeout 关闭超时错误
var ErrShutdownTimeout = errors.New("关闭超时")

// AppContext 应用级别的 Context 封装
// 嘿嘿~ 这是整个应用的上下文管理器！(´∀｀)
type AppContext struct {
	ctx        context.Context
	cancel     context.CancelFunc
	wg         sync.WaitGroup
	shutdownCh chan struct{}
	mu         sync.RWMutex
	closed     bool
}

// NewAppContext 创建新的应用 Context
// 支持超时设置和优雅关闭~ ✨
func NewAppContext(parent context.Context) *AppContext {
	ctx, cancel := context.WithCancel(parent)
	return &AppContext{
		ctx:        ctx,
		cancel:     cancel,
		shutdownCh: make(chan struct{}),
	}
}

// Context 获取底层 context.Context
// 可以传递给 Service 和 Repository 层~ 💖
func (ac *AppContext) Context() context.Context {
	return ac.ctx
}

// Done 返回关闭信号通道
func (ac *AppContext) Done() <-chan struct{} {
	return ac.ctx.Done()
}

// Shutdown 触发优雅关闭
// 呀~ 会等待所有任务完成后再关闭哦！🎮
func (ac *AppContext) Shutdown(timeout time.Duration) error {
	ac.mu.Lock()
	if ac.closed {
		ac.mu.Unlock()
		return nil
	}
	ac.closed = true
	ac.mu.Unlock()

	// 发送关闭信号
	close(ac.shutdownCh)

	// 取消 context
	ac.cancel()

	// 等待所有任务完成（带超时）
	done := make(chan struct{})
	go func() {
		ac.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-time.After(timeout):
		return ErrShutdownTimeout
	}
}

// Go 启动一个受管理的 goroutine
// 嘿嘿~ 这样可以追踪所有 goroutine！💫
func (ac *AppContext) Go(fn func(ctx context.Context)) {
	ac.wg.Add(1)
	go func() {
		defer ac.wg.Done()
		fn(ac.ctx)
	}()
}

// ShutdownCh 获取关闭信号通道
func (ac *AppContext) ShutdownCh() <-chan struct{} {
	return ac.shutdownCh
}

// IsClosed 检查是否已关闭
func (ac *AppContext) IsClosed() bool {
	ac.mu.RLock()
	defer ac.mu.RUnlock()
	return ac.closed
}
