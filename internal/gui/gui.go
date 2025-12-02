package gui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/XiaoLFeng/llm-memory/startup"
)

// 样式定义
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#A78BFA")).
			MarginBottom(1)

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#7C3AED")).
			Padding(0, 1)

	normalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#93C5FD")).
			Padding(0, 1)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#64748B")).
			MarginTop(1)
)

// MenuItem 菜单项
type MenuItem struct {
	Title       string
	Description string
}

// Model 主界面模型
// 嘿嘿~ 这是一个简单的菜单界面！(´∀｀)💖
type Model struct {
	bs       *startup.Bootstrap
	items    []MenuItem
	selected int
	width    int
	height   int
	quitting bool
}

// NewModel 创建新的 GUI 模型
func NewModel(bs *startup.Bootstrap) Model {
	items := []MenuItem{
		{Title: "📝 记忆管理", Description: "查看和管理记忆内容"},
		{Title: "📋 计划管理", Description: "管理你的计划"},
		{Title: "✅ TODO 管理", Description: "管理待办事项"},
		{Title: "🚪 退出", Description: "退出程序"},
	}

	return Model{
		bs:       bs,
		items:    items,
		selected: 0,
		width:    80,
		height:   24,
	}
}

// Init 初始化
func (m Model) Init() tea.Cmd {
	return nil
}

// Update 处理输入
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "up", "k":
			if m.selected > 0 {
				m.selected--
			} else {
				m.selected = len(m.items) - 1
			}

		case "down", "j":
			if m.selected < len(m.items)-1 {
				m.selected++
			} else {
				m.selected = 0
			}

		case "enter":
			// 处理选择
			switch m.selected {
			case 0:
				// 记忆管理 - TODO: 实现子菜单
				return m, nil
			case 1:
				// 计划管理 - TODO: 实现子菜单
				return m, nil
			case 2:
				// TODO 管理 - TODO: 实现子菜单
				return m, nil
			case 3:
				// 退出
				m.quitting = true
				return m, tea.Quit
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

// View 渲染界面
func (m Model) View() string {
	if m.quitting {
		return "再见~ 👋\n"
	}

	var b strings.Builder

	// 标题
	title := titleStyle.Render("🧠 LLM-Memory 管理系统")
	b.WriteString(title)
	b.WriteString("\n\n")

	// 菜单项
	for i, item := range m.items {
		var line string
		if i == m.selected {
			line = selectedStyle.Render(fmt.Sprintf("> %s", item.Title))
		} else {
			line = normalStyle.Render(fmt.Sprintf("  %s", item.Title))
		}
		b.WriteString(line)
		b.WriteString("\n")
	}

	// 当前选中项的描述
	b.WriteString("\n")
	if m.selected < len(m.items) {
		desc := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#94A3B8")).
			Italic(true).
			Render(m.items[m.selected].Description)
		b.WriteString(desc)
	}

	// 帮助信息
	help := helpStyle.Render("↑/↓ 选择 | Enter 确认 | q 退出")
	b.WriteString("\n\n")
	b.WriteString(help)

	return b.String()
}

// Run 运行 GUI
// 呀~ 启动简单的管理界面！✨
func Run(bs *startup.Bootstrap) error {
	p := tea.NewProgram(NewModel(bs), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
