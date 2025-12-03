package group

import (
	"context"
	"fmt"
	"strings"

	"github.com/XiaoLFeng/llm-memory/internal/models/entity"
	"github.com/XiaoLFeng/llm-memory/internal/tui/common"
	"github.com/XiaoLFeng/llm-memory/internal/tui/components"
	"github.com/XiaoLFeng/llm-memory/internal/tui/layout"
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
	return fmt.Sprintf("%s %d 个路径 | %s", styles.IconFolder, pathCount, utils.FormatRelativeTime(i.group.CreatedAt))
}

func (i groupItem) FilterValue() string {
	return i.group.Name
}

// ListModel 组列表模型
type ListModel struct {
	bs            *startup.Bootstrap
	groups        []entity.Group
	selectedIndex int
	width         int
	height        int
	frame         *components.Frame
	loading       bool
	err           error
}

// NewListModel 创建组列表模型
func NewListModel(bs *startup.Bootstrap) *ListModel {
	return &ListModel{
		bs:            bs,
		loading:       true,
		selectedIndex: 0,
		frame:         components.NewFrame(80, 24),
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
	return tea.Batch(m.loadGroups(), common.StartAutoRefresh())
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
	id int64
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
		m.frame.SetSize(msg.Width, msg.Height)

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

	case common.AutoRefreshMsg:
		// 自动刷新：静默加载数据
		cmds = append(cmds, m.loadGroups())
	}

	return m, tea.Batch(cmds...)
}

// deleteGroup 删除组
func (m *ListModel) deleteGroup(id int64) tea.Cmd {
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
		content := lipgloss.NewStyle().
			Foreground(styles.Info).
			Render("加载中...")
		return layout.ListPage(
			m.frame,
			"组管理 > 组列表",
			styles.IconUsers+" 组管理",
			content,
			[]string{},
			"",
		)
	}

	if m.err != nil {
		content := lipgloss.NewStyle().
			Foreground(styles.Error).
			Render("错误: " + m.err.Error())
		return layout.ListPage(
			m.frame,
			"组管理 > 组列表",
			styles.IconUsers+" 组管理",
			content,
			[]string{},
			"",
		)
	}

	if len(m.groups) == 0 {
		content := lipgloss.NewStyle().
			Foreground(styles.Subtext0).
			Render("暂无组~ 按 c 创建新组")

		keys := []string{
			styles.StatusKeyStyle.Render("c") + " " + styles.StatusValueStyle.Render("新建"),
			styles.StatusKeyStyle.Render("esc") + " " + styles.StatusValueStyle.Render("返回"),
		}

		return layout.ListPage(
			m.frame,
			"组管理 > 组列表",
			styles.IconUsers+" 组管理",
			content,
			keys,
			"",
		)
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

	// 状态栏
	keys := []string{
		styles.StatusKeyStyle.Render("↑/↓") + " " + styles.StatusValueStyle.Render("选择"),
		styles.StatusKeyStyle.Render("enter") + " " + styles.StatusValueStyle.Render("查看"),
		styles.StatusKeyStyle.Render("c") + " " + styles.StatusValueStyle.Render("新建"),
		styles.StatusKeyStyle.Render("d") + " " + styles.StatusValueStyle.Render("删除"),
		styles.StatusKeyStyle.Render("esc") + " " + styles.StatusValueStyle.Render("返回"),
	}

	extra := fmt.Sprintf("共 %d 个组", len(m.groups))

	return layout.ListPage(
		m.frame,
		"组管理 > 组列表",
		styles.IconUsers+" 组管理",
		listItems.String(),
		keys,
		extra,
	)
}
