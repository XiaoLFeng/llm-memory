package startup

import (
	"os"
	"os/signal"
	"syscall"
)

// SignalHandler 信号处理器
// 呀~ 处理系统信号，实现优雅关闭！🎮
type SignalHandler struct {
	sigCh  chan os.Signal
	doneCh chan struct{}
}

// NewSignalHandler 创建信号处理器
func NewSignalHandler() *SignalHandler {
	return &SignalHandler{
		sigCh:  make(chan os.Signal, 1),
		doneCh: make(chan struct{}),
	}
}

// Start 开始监听信号
// 嘿嘿~ 监听 SIGINT 和 SIGTERM 信号！💖
func (sh *SignalHandler) Start(onSignal func(os.Signal)) {
	signal.Notify(sh.sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		select {
		case sig := <-sh.sigCh:
			onSignal(sig)
		case <-sh.doneCh:
			return
		}
	}()
}

// Stop 停止监听
func (sh *SignalHandler) Stop() {
	signal.Stop(sh.sigCh)
	close(sh.doneCh)
}
