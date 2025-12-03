package components

import (
	"strings"

	"github.com/XiaoLFeng/llm-memory/internal/tui/styles"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// StatusBar 状态栏组件
type StatusBar struct {
	breadcrumb string
	keys       []key.Binding
	extra      string
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

// SetExtra 设置额外信息
func (s *StatusBar) SetExtra(extra string) {
	s.extra = extra
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
	// 快捷键提示
	var keyStrs []string
	for _, k := range s.keys {
		keyStr := styles.StatusKeyStyle.Render(k.Help().Key) + " " +
			styles.StatusValueStyle.Render(k.Help().Desc)
		keyStrs = append(keyStrs, keyStr)
	}
	keysStr := strings.Join(keyStrs, "  │  ")

	// 状态栏样式 - 带边框
	statusStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Border).
		Foreground(styles.Subtext0).
		Width(s.width-2).
		Padding(0, 1)

	return statusStyle.Render(keysStr)
}

// ViewWithBreadcrumb 带面包屑的渲染
func (s *StatusBar) ViewWithBreadcrumb() string {
	// 面包屑
	breadcrumb := styles.StatusKeyStyle.Render(s.breadcrumb)

	// 快捷键提示
	var keyStrs []string
	for _, k := range s.keys {
		keyStr := styles.StatusKeyStyle.Render(k.Help().Key) + " " +
			styles.StatusValueStyle.Render(k.Help().Desc)
		keyStrs = append(keyStrs, keyStr)
	}
	keysStr := strings.Join(keyStrs, "  │  ")

	// 计算间距
	left := breadcrumb
	right := keysStr

	gap := s.width - lipgloss.Width(left) - lipgloss.Width(right) - 6
	if gap < 0 {
		gap = 0
	}

	content := left + strings.Repeat(" ", gap) + right

	// 状态栏样式 - 带边框
	statusStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Border).
		Foreground(styles.Subtext0).
		Width(s.width-2).
		Padding(0, 1)

	return statusStyle.Render(content)
}

// RenderKeysOnly 只渲染快捷键（用于状态栏）
func RenderKeysOnly(keys []string, width int) string {
	keysStr := strings.Join(keys, "  │  ")

	statusStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Border).
		Foreground(styles.Subtext0).
		Width(width-2).
		Padding(0, 1)

	return statusStyle.Render(keysStr)
}
