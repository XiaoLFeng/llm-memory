package components

import (
	"time"

	"github.com/XiaoLFeng/llm-memory/internal/tui/styles"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
type Toast struct {
	message   string
	toastType ToastType
	visible   bool
	duration  time.Duration
	width     int
	height    int
}

// NewToast 创建 Toast 组件
func NewToast() *Toast {
	return &Toast{
		duration: 3 * time.Second,
		width:    80,
		height:   24,
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

// SetSize 设置窗口大小
func (t *Toast) SetSize(width, height int) {
	t.width = width
	t.height = height
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

	var borderColor lipgloss.Color
	var icon string

	switch t.toastType {
	case ToastSuccess:
		borderColor = styles.Success
		icon = "+"
	case ToastError:
		borderColor = styles.Error
		icon = "x"
	case ToastWarning:
		borderColor = styles.Warning
		icon = "!"
	case ToastInfo:
		borderColor = styles.Info
		icon = "i"
	}

	// 创建 Toast 样式 - 带边框
	toastStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Foreground(styles.Text).
		Padding(0, 2)

	// 图标样式
	iconStyle := lipgloss.NewStyle().
		Foreground(borderColor).
		Bold(true)

	content := iconStyle.Render(icon) + "  " + t.message

	return toastStyle.Render(content)
}

// RenderOverlay 渲染为浮动层（用于居中显示）
func (t *Toast) RenderOverlay(base string) string {
	if !t.visible {
		return base
	}

	toastView := t.View()
	return PlaceOverlay(base, toastView, t.width, t.height, TopCenter)
}
