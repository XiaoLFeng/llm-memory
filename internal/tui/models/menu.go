package models

import (
	"strings"

	"github.com/XiaoLFeng/llm-memory/internal/tui/common"
	"github.com/XiaoLFeng/llm-memory/internal/tui/components"
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
		{Title: "记忆管理", Description: "查看和管理记忆内容", Icon: styles.IconBrain, Page: common.PageMemoryList},
		{Title: "计划管理", Description: "管理你的计划", Icon: styles.IconTasks, Page: common.PagePlanList},
		{Title: "待办管理", Description: "管理待办事项", Icon: styles.IconTodo, Page: common.PageTodoList},
		{Title: "组管理", Description: "管理路径组，组内共享数据", Icon: styles.IconUsers, Page: common.PageGroupList},
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
	// 计算合适的宽度
	contentWidth := m.width - 4
	if contentWidth > 70 {
		contentWidth = 70
	}
	if contentWidth < 40 {
		contentWidth = 40
	}

	// Logo 区域
	logoStyle := lipgloss.NewStyle().
		Foreground(styles.Primary).
		Bold(true).
		Align(lipgloss.Center).
		Width(contentWidth)

	subtitleStyle := lipgloss.NewStyle().
		Foreground(styles.Subtext0).
		Align(lipgloss.Center).
		Width(contentWidth)

	logo := logoStyle.Render(styles.IconLogo + " LLM-Memory")
	subtitle := subtitleStyle.Render("AI 记忆管理系统 v1.0")

	// Logo 卡片
	logoContent := logo + "\n" + subtitle
	logoCard := components.CardSimple(logoContent, contentWidth)

	// 菜单项
	var menuItems strings.Builder
	for i, item := range m.items {
		var itemLine string

		// 选中指示器
		indicator := "  "
		if i == m.selected {
			indicator = "▸ "
		}

		// 图标和标题
		iconStyle := lipgloss.NewStyle().Foreground(styles.Primary)
		titleStyle := lipgloss.NewStyle().Foreground(styles.Text)
		if i == m.selected {
			titleStyle = titleStyle.Bold(true).Foreground(styles.Primary)
		}

		itemLine = indicator + iconStyle.Render(item.Icon) + "  " + titleStyle.Render(item.Title)

		// 描述（仅选中项显示）
		if i == m.selected {
			descStyle := lipgloss.NewStyle().
				Foreground(styles.Subtext0).
				MarginLeft(5)
			itemLine += "\n" + descStyle.Render(item.Description)
		}

		if i > 0 {
			menuItems.WriteString("\n")
		}
		menuItems.WriteString(itemLine)
		if i == m.selected {
			menuItems.WriteString("\n")
		}
	}

	// 菜单卡片
	menuCard := components.Card("功能菜单", menuItems.String(), contentWidth)

	// 快捷键提示
	keys := []key.Binding{
		key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "向上")),
		key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "向下")),
		key.NewBinding(key.WithKeys("enter"), key.WithHelp("Enter", "确认")),
		key.NewBinding(key.WithKeys("q"), key.WithHelp("q", "退出")),
	}

	var keyStrs []string
	for _, k := range keys {
		keyStr := styles.StatusKeyStyle.Render(k.Help().Key) + " " +
			styles.StatusValueStyle.Render(k.Help().Desc)
		keyStrs = append(keyStrs, keyStr)
	}
	statusBar := components.RenderKeysOnly(keyStrs, contentWidth)

	// 组合所有内容
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		logoCard,
		"",
		menuCard,
		"",
		statusBar,
	)

	// 居中显示
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}

// SetSize 设置窗口大小
func (m *MenuModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}
