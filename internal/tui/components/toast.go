package components

import (
	"time"

	"github.com/XiaoLFeng/llm-memory/internal/tui/styles"
	tea "github.com/charmbracelet/bubbletea"
)

// ToastType 提示消息类型
type ToastType int

const (
	ToastSuccess ToastType = iota
	ToastError
	ToastWarning
	ToastInfo
)

// Toast 提示消息组件
// 嘿嘿~ 用于显示操作反馈的短暂提示！💬
type Toast struct {
	message   string
	toastType ToastType
	visible   bool
	duration  time.Duration
}

// NewToast 创建 Toast 组件
func NewToast() *Toast {
	return &Toast{
		duration: 3 * time.Second,
	}
}

// Show 显示提示消息
func (t *Toast) Show(message string, toastType ToastType) {
	t.message = message
	t.toastType = toastType
	t.visible = true
}

// Hide 隐藏提示消息
func (t *Toast) Hide() {
	t.visible = false
}

// IsVisible 是否可见
func (t *Toast) IsVisible() bool {
	return t.visible
}

// HideAfter 延迟隐藏
func (t *Toast) HideAfter() tea.Cmd {
	return tea.Tick(t.duration, func(time.Time) tea.Msg {
		return hideToastMsg{}
	})
}

type hideToastMsg struct{}

// Init 初始化
func (t *Toast) Init() tea.Cmd {
	return nil
}

// Update 处理输入
func (t *Toast) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case hideToastMsg:
		t.Hide()
	}
	return t, nil
}

// View 渲染界面
func (t *Toast) View() string {
	if !t.visible {
		return ""
	}

	var style = styles.InfoStyle
	var icon string

	switch t.toastType {
	case ToastSuccess:
		style = styles.SuccessStyle
		icon = "✓ "
	case ToastError:
		style = styles.ErrorStyle
		icon = "✗ "
	case ToastWarning:
		style = styles.WarningStyle
		icon = "⚠ "
	case ToastInfo:
		style = styles.InfoStyle
		icon = "ℹ "
	}

	return style.Render(icon + t.message)
}
