package components

import (
	"github.com/XiaoLFeng/llm-memory/internal/tui/styles"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Spinner 加载动画组件
// 嘿嘿~ 用于显示加载状态的动画！⏳
type Spinner struct {
	spinner spinner.Model
	message string
	visible bool
}

// NewSpinner 创建加载动画组件
func NewSpinner() *Spinner {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(styles.Primary)

	return &Spinner{
		spinner: s,
		message: "加载中...",
		visible: false,
	}
}

// NewSpinnerWithStyle 创建带自定义样式的加载动画
func NewSpinnerWithStyle(spinnerType spinner.Spinner, color lipgloss.Color) *Spinner {
	s := spinner.New()
	s.Spinner = spinnerType
	s.Style = lipgloss.NewStyle().Foreground(color)

	return &Spinner{
		spinner: s,
		message: "加载中...",
		visible: false,
	}
}

// Show 显示加载动画
func (s *Spinner) Show(message string) tea.Cmd {
	s.message = message
	s.visible = true
	return s.spinner.Tick
}

// Hide 隐藏加载动画
func (s *Spinner) Hide() {
	s.visible = false
}

// IsVisible 是否可见
func (s *Spinner) IsVisible() bool {
	return s.visible
}

// SetMessage 设置加载消息
func (s *Spinner) SetMessage(message string) {
	s.message = message
}

// Init 初始化
func (s *Spinner) Init() tea.Cmd {
	return nil
}

// Update 处理输入
func (s *Spinner) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !s.visible {
		return s, nil
	}

	var cmd tea.Cmd
	s.spinner, cmd = s.spinner.Update(msg)
	return s, cmd
}

// View 渲染界面
func (s *Spinner) View() string {
	if !s.visible {
		return ""
	}

	messageStyle := lipgloss.NewStyle().
		Foreground(styles.Subtext0)

	return s.spinner.View() + " " + messageStyle.Render(s.message)
}

// ViewInline 渲染为内联样式（不带消息）
func (s *Spinner) ViewInline() string {
	if !s.visible {
		return ""
	}
	return s.spinner.View()
}

// Tick 返回 tick 命令
func (s *Spinner) Tick() tea.Cmd {
	return s.spinner.Tick
}

// 预定义的 Spinner 样式
var (
	// SpinnerDot 点状加载动画
	SpinnerDot = spinner.Dot

	// SpinnerLine 线条加载动画
	SpinnerLine = spinner.Line

	// SpinnerMiniDot 小点加载动画
	SpinnerMiniDot = spinner.MiniDot

	// SpinnerJump 跳跃加载动画
	SpinnerJump = spinner.Jump

	// SpinnerPulse 脉冲加载动画
	SpinnerPulse = spinner.Pulse

	// SpinnerPoints 多点加载动画
	SpinnerPoints = spinner.Points

	// SpinnerGlobe 地球加载动画
	SpinnerGlobe = spinner.Globe

	// SpinnerMoon 月亮加载动画
	SpinnerMoon = spinner.Moon

	// SpinnerMonkey 猴子加载动画
	SpinnerMonkey = spinner.Monkey
)

// LoadingView 简单的加载视图（静态）
func LoadingView(message string) string {
	spinnerStyle := lipgloss.NewStyle().
		Foreground(styles.Primary)

	messageStyle := lipgloss.NewStyle().
		Foreground(styles.Subtext0)

	return spinnerStyle.Render("⏳") + " " + messageStyle.Render(message)
}

// LoadingCard 加载卡片（带边框）
func LoadingCard(message string, width int) string {
	content := LoadingView(message)
	return CardSimple(content, width)
}
