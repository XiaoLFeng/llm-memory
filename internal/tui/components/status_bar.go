package components

import (
	"github.com/XiaoLFeng/llm-memory/internal/tui/styles"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// StatusBar 状态栏组件
// 嘿嘿~ 显示当前位置和快捷键提示！📍
type StatusBar struct {
	breadcrumb string
	keys       []key.Binding
	width      int
}

// NewStatusBar 创建状态栏组件
func NewStatusBar() *StatusBar {
	return &StatusBar{
		width: 80,
	}
}

// SetBreadcrumb 设置面包屑导航
func (s *StatusBar) SetBreadcrumb(breadcrumb string) {
	s.breadcrumb = breadcrumb
}

// SetKeys 设置快捷键
func (s *StatusBar) SetKeys(keys []key.Binding) {
	s.keys = keys
}

// SetWidth 设置宽度
func (s *StatusBar) SetWidth(width int) {
	s.width = width
}

// Init 初始化
func (s *StatusBar) Init() tea.Cmd {
	return nil
}

// Update 处理输入
func (s *StatusBar) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width = msg.Width
	}
	return s, nil
}

// View 渲染界面
func (s *StatusBar) View() string {
	// 面包屑
	breadcrumb := styles.StatusKeyStyle.Render(s.breadcrumb)

	// 快捷键提示
	var keysStr string
	for i, k := range s.keys {
		if i > 0 {
			keysStr += " | "
		}
		keysStr += styles.StatusKeyStyle.Render(k.Help().Key) + " " +
			styles.StatusValueStyle.Render(k.Help().Desc)
	}

	// 组合状态栏
	left := breadcrumb
	right := keysStr

	// 计算间距
	gap := s.width - lipgloss.Width(left) - lipgloss.Width(right)
	if gap < 0 {
		gap = 0
	}

	return styles.StatusBarStyle.
		Width(s.width).
		Render(left + lipgloss.NewStyle().Width(gap).Render("") + right)
}
