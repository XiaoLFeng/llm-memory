package models

import (
	"strings"

	"github.com/XiaoLFeng/llm-memory/internal/tui/common"
	"github.com/XiaoLFeng/llm-memory/internal/tui/styles"
	"github.com/XiaoLFeng/llm-memory/startup"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// MenuItem 菜单项
type MenuItem struct {
	Title       string
	Description string
	Icon        string
	Page        common.PageType
}

// MenuModel 主菜单模型
// 嘿嘿~ 这是主菜单界面！💖
type MenuModel struct {
	bs       *startup.Bootstrap
	items    []MenuItem
	selected int
	width    int
	height   int
}

// NewMenuModel 创建主菜单模型
func NewMenuModel(bs *startup.Bootstrap) *MenuModel {
	items := []MenuItem{
		{Title: "记忆管理", Description: "查看和管理记忆内容", Icon: "📝", Page: common.PageMemoryList},
		{Title: "计划管理", Description: "管理你的计划", Icon: "📋", Page: common.PagePlanList},
		{Title: "待办管理", Description: "管理待办事项", Icon: "✅", Page: common.PageTodoList},
		{Title: "组管理", Description: "管理路径组，组内共享数据", Icon: "👥", Page: common.PageGroupList},
	}

	return &MenuModel{
		bs:       bs,
		items:    items,
		selected: 0,
		width:    80,
		height:   24,
	}
}

// Title 返回页面标题
func (m *MenuModel) Title() string {
	return "主菜单"
}

// ShortHelp 返回快捷键帮助
func (m *MenuModel) ShortHelp() []key.Binding {
	return []key.Binding{common.KeyUp, common.KeyDown, common.KeyEnter, common.KeyQuit}
}

// Init 初始化
func (m *MenuModel) Init() tea.Cmd {
	return nil
}

// Update 处理输入
func (m *MenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, common.KeyQuit):
			return m, tea.Quit

		case key.Matches(msg, common.KeyUp):
			if m.selected > 0 {
				m.selected--
			} else {
				m.selected = len(m.items) - 1
			}

		case key.Matches(msg, common.KeyDown):
			if m.selected < len(m.items)-1 {
				m.selected++
			} else {
				m.selected = 0
			}

		case key.Matches(msg, common.KeyEnter):
			if m.selected < len(m.items) {
				return m, common.Navigate(m.items[m.selected].Page)
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

// View 渲染界面
func (m *MenuModel) View() string {
	var b strings.Builder

	// 标题
	title := styles.TitleStyle.Render("🧠 LLM-Memory 管理系统")
	b.WriteString(title)
	b.WriteString("\n\n")

	// 菜单项
	for i, item := range m.items {
		var line string
		itemText := item.Icon + " " + item.Title
		if i == m.selected {
			line = styles.SelectedStyle.Render("> " + itemText)
		} else {
			line = styles.NormalStyle.Render("  " + itemText)
		}
		b.WriteString(line)
		b.WriteString("\n")
	}

	// 退出选项
	b.WriteString("\n")
	exitText := "🚪 退出"
	if m.selected == len(m.items) {
		b.WriteString(styles.SelectedStyle.Render("> " + exitText))
	} else {
		b.WriteString(styles.NormalStyle.Render("  " + exitText))
	}
	b.WriteString("\n")

	// 当前选中项的描述
	b.WriteString("\n")
	if m.selected < len(m.items) {
		desc := styles.DescStyle.Render(m.items[m.selected].Description)
		b.WriteString(desc)
	} else {
		desc := styles.DescStyle.Render("退出程序")
		b.WriteString(desc)
	}

	// 帮助信息
	help := styles.HelpStyle.Render("↑/↓ 选择 | Enter 确认 | q 退出")
	b.WriteString("\n\n")
	b.WriteString(help)

	return b.String()
}

// SetSize 设置窗口大小
func (m *MenuModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// containerStyle 容器样式
var containerStyle = lipgloss.NewStyle().
	Padding(1, 2)
