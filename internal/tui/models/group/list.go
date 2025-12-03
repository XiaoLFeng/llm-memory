package group

import (
	"context"
	"fmt"
	"strings"

	"github.com/XiaoLFeng/llm-memory/internal/models/entity"
	"github.com/XiaoLFeng/llm-memory/internal/tui/common"
	"github.com/XiaoLFeng/llm-memory/internal/tui/components"
	"github.com/XiaoLFeng/llm-memory/internal/tui/styles"
	"github.com/XiaoLFeng/llm-memory/internal/tui/utils"
	"github.com/XiaoLFeng/llm-memory/startup"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// groupItem 组列表项
type groupItem struct {
	group entity.Group
}

func (i groupItem) Title() string {
	return fmt.Sprintf("%d. %s", i.group.ID, i.group.Name)
}

func (i groupItem) Description() string {
	pathCount := len(i.group.Paths)
	return fmt.Sprintf("📂 %d 个路径 | %s", pathCount, utils.FormatRelativeTime(i.group.CreatedAt))
}

func (i groupItem) FilterValue() string {
	return i.group.Name
}

// ListModel 组列表模型
// 嘿嘿~ 展示所有组的列表！👥
type ListModel struct {
	bs            *startup.Bootstrap
	groups        []entity.Group
	selectedIndex int
	width         int
	height        int
	loading       bool
	err           error
}

// NewListModel 创建组列表模型
func NewListModel(bs *startup.Bootstrap) *ListModel {
	return &ListModel{
		bs:            bs,
		loading:       true,
		selectedIndex: 0,
	}
}

// Title 返回页面标题
func (m *ListModel) Title() string {
	return "组管理"
}

// ShortHelp 返回快捷键帮助
func (m *ListModel) ShortHelp() []key.Binding {
	return []key.Binding{
		common.KeyUp, common.KeyDown, common.KeyEnter,
		common.KeyCreate, common.KeyDelete, common.KeyBack,
	}
}

// Init 初始化
func (m *ListModel) Init() tea.Cmd {
	return m.loadGroups()
}

// loadGroups 加载组列表
func (m *ListModel) loadGroups() tea.Cmd {
	return func() tea.Msg {
		groups, err := m.bs.GroupService.ListGroups(context.Background())
		if err != nil {
			return groupsErrorMsg{err}
		}
		return groupsLoadedMsg{groups}
	}
}

type groupsLoadedMsg struct {
	groups []entity.Group
}

type groupsErrorMsg struct {
	err error
}

type groupDeletedMsg struct {
	id uint
}

// Update 处理输入
func (m *ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, common.KeyBack):
			return m, common.Back()

		case key.Matches(msg, common.KeyCreate):
			return m, common.Navigate(common.PageGroupCreate)

		case key.Matches(msg, common.KeyUp):
			if m.selectedIndex > 0 {
				m.selectedIndex--
			}

		case key.Matches(msg, common.KeyDown):
			if m.selectedIndex < len(m.groups)-1 {
				m.selectedIndex++
			}

		case key.Matches(msg, common.KeyEnter):
			if len(m.groups) > 0 && m.selectedIndex < len(m.groups) {
				return m, common.Navigate(common.PageGroupDetail, map[string]any{"id": m.groups[m.selectedIndex].ID})
			}

		case key.Matches(msg, common.KeyDelete):
			if len(m.groups) > 0 && m.selectedIndex < len(m.groups) {
				group := m.groups[m.selectedIndex]
				return m, common.ShowConfirm(
					"删除组",
					fmt.Sprintf("确定要删除组「%s」吗？", group.Name),
					m.deleteGroup(group.ID),
					nil,
				)
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case groupsLoadedMsg:
		m.loading = false
		m.groups = msg.groups
		m.selectedIndex = 0

	case groupsErrorMsg:
		m.loading = false
		m.err = msg.err

	case groupDeletedMsg:
		cmds = append(cmds, m.loadGroups())
		cmds = append(cmds, common.ShowToast("组已删除", common.ToastSuccess))

	case common.RefreshMsg:
		m.loading = true
		cmds = append(cmds, m.loadGroups())
	}

	return m, tea.Batch(cmds...)
}

// deleteGroup 删除组
func (m *ListModel) deleteGroup(id uint) tea.Cmd {
	return func() tea.Msg {
		err := m.bs.GroupService.DeleteGroup(context.Background(), id)
		if err != nil {
			return groupsErrorMsg{err}
		}
		return groupDeletedMsg{id}
	}
}

// View 渲染界面
func (m *ListModel) View() string {
	if m.loading {
		loadingStyle := lipgloss.NewStyle().
			Foreground(styles.Info).
			Align(lipgloss.Center)
		content := loadingStyle.Render("加载中...")
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
	}

	if m.err != nil {
		errorStyle := lipgloss.NewStyle().
			Foreground(styles.Error).
			Align(lipgloss.Center)
		content := errorStyle.Render("错误: " + m.err.Error())
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
	}

	if len(m.groups) == 0 {
		emptyStyle := lipgloss.NewStyle().
			Foreground(styles.Subtext0).
			Align(lipgloss.Center)
		content := emptyStyle.Render("暂无组~ 按 c 创建新组")

		// 状态栏
		keys := []string{
			lipgloss.NewStyle().Foreground(styles.Primary).Render("c") + " 新建",
			lipgloss.NewStyle().Foreground(styles.Primary).Render("esc") + " 返回",
		}
		statusBar := components.RenderKeysOnly(keys, m.width)

		view := lipgloss.Place(m.width, m.height-3, lipgloss.Center, lipgloss.Center, content)
		return lipgloss.JoinVertical(lipgloss.Left, view, statusBar)
	}

	// 构建列表内容
	var listItems strings.Builder
	for i, group := range m.groups {
		var line string

		// 构建一行内容
		indicator := "  "
		if i == m.selectedIndex {
			indicator = lipgloss.NewStyle().Foreground(styles.Primary).Render("▸ ")
		} else {
			indicator = "  "
		}

		// 组名和ID
		nameStyle := lipgloss.NewStyle().Foreground(styles.Text).Bold(true)
		if i == m.selectedIndex {
			nameStyle = nameStyle.Foreground(styles.Primary)
		}
		name := nameStyle.Render(fmt.Sprintf("%d. %s", group.ID, group.Name))

		// 路径数量
		pathCount := len(group.Paths)
		countBadge := lipgloss.NewStyle().
			Foreground(styles.Subtext0).
			Render(fmt.Sprintf("(%d 个路径)", pathCount))

		// 时间
		timeStr := lipgloss.NewStyle().
			Foreground(styles.Overlay1).
			Italic(true).
			Render(utils.FormatRelativeTime(group.CreatedAt))

		line = fmt.Sprintf("%s%s %s │ %s", indicator, name, countBadge, timeStr)

		listItems.WriteString(line)
		if i < len(m.groups)-1 {
			listItems.WriteString("\n")
		}
	}

	// 计算卡片宽度
	cardWidth := m.width - 4
	if cardWidth < 60 {
		cardWidth = 60
	}

	// 使用卡片包装列表
	titleWithCount := fmt.Sprintf("👥 组管理 %s",
		lipgloss.NewStyle().Foreground(styles.Subtext0).Render(fmt.Sprintf("(%d)", len(m.groups))))
	card := components.Card(titleWithCount, listItems.String(), cardWidth)

	// 状态栏
	keys := []string{
		lipgloss.NewStyle().Foreground(styles.Primary).Render("↑/↓") + " 选择",
		lipgloss.NewStyle().Foreground(styles.Primary).Render("enter") + " 查看",
		lipgloss.NewStyle().Foreground(styles.Primary).Render("c") + " 新建",
		lipgloss.NewStyle().Foreground(styles.Primary).Render("d") + " 删除",
		lipgloss.NewStyle().Foreground(styles.Primary).Render("esc") + " 返回",
	}
	statusBar := components.RenderKeysOnly(keys, m.width)

	// 组合视图
	contentHeight := m.height - 3
	centeredCard := lipgloss.Place(m.width, contentHeight, lipgloss.Center, lipgloss.Center, card)

	return lipgloss.JoinVertical(lipgloss.Left, centeredCard, statusBar)
}
