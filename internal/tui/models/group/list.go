package group

import (
	"context"
	"fmt"
	"strings"

	"github.com/XiaoLFeng/llm-memory/internal/tui/common"
	"github.com/XiaoLFeng/llm-memory/internal/tui/styles"
	"github.com/XiaoLFeng/llm-memory/internal/tui/utils"
	"github.com/XiaoLFeng/llm-memory/pkg/types"
	"github.com/XiaoLFeng/llm-memory/startup"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// groupItem 组列表项
type groupItem struct {
	group types.Group
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
	bs      *startup.Bootstrap
	list    list.Model
	groups  []types.Group
	width   int
	height  int
	loading bool
	err     error
}

// NewListModel 创建组列表模型
func NewListModel(bs *startup.Bootstrap) *ListModel {
	// 创建列表
	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = styles.ListSelectedStyle
	delegate.Styles.SelectedDesc = styles.ListDescStyle
	delegate.Styles.NormalTitle = styles.ListItemStyle
	delegate.Styles.NormalDesc = styles.ListDescStyle

	l := list.New([]list.Item{}, delegate, 80, 20)
	l.Title = "👥 组管理"
	l.SetShowHelp(false)
	l.SetFilteringEnabled(true)
	l.Styles.Title = styles.ListTitleStyle

	return &ListModel{
		bs:      bs,
		list:    l,
		loading: true,
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
	groups []types.Group
}

type groupsErrorMsg struct {
	err error
}

type groupDeletedMsg struct {
	id int
}

// Update 处理输入
func (m *ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// 如果正在过滤，让列表处理
		if m.list.FilterState() == list.Filtering {
			break
		}

		switch {
		case key.Matches(msg, common.KeyBack):
			return m, common.Back()

		case key.Matches(msg, common.KeyCreate):
			return m, common.Navigate(common.PageGroupCreate)

		case key.Matches(msg, common.KeyEnter):
			if item, ok := m.list.SelectedItem().(groupItem); ok {
				return m, common.Navigate(common.PageGroupDetail, map[string]any{"id": item.group.ID})
			}

		case key.Matches(msg, common.KeyDelete):
			if item, ok := m.list.SelectedItem().(groupItem); ok {
				return m, common.ShowConfirm(
					"删除组",
					fmt.Sprintf("确定要删除组「%s」吗？", item.group.Name),
					m.deleteGroup(item.group.ID),
					nil,
				)
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetSize(msg.Width-4, msg.Height-8)

	case groupsLoadedMsg:
		m.loading = false
		m.groups = msg.groups
		items := make([]list.Item, len(msg.groups))
		for i, group := range msg.groups {
			items[i] = groupItem{group: group}
		}
		m.list.SetItems(items)

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

	// 更新列表
	newList, cmd := m.list.Update(msg)
	m.list = newList
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// deleteGroup 删除组
func (m *ListModel) deleteGroup(id int) tea.Cmd {
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
	var b strings.Builder

	if m.loading {
		b.WriteString(styles.InfoStyle.Render("加载中..."))
		return b.String()
	}

	if m.err != nil {
		b.WriteString(styles.ErrorStyle.Render("错误: " + m.err.Error()))
		return b.String()
	}

	if len(m.groups) == 0 {
		b.WriteString(styles.TitleStyle.Render("👥 组管理"))
		b.WriteString("\n\n")
		b.WriteString(styles.MutedStyle.Render("暂无组~ 按 c 创建新组"))
		b.WriteString("\n\n")
		b.WriteString(styles.HelpStyle.Render("c 新建 | esc 返回"))
		return b.String()
	}

	b.WriteString(m.list.View())
	b.WriteString("\n")
	b.WriteString(styles.HelpStyle.Render("↑/↓ 选择 | enter 查看 | c 新建 | d 删除 | esc 返回"))

	return b.String()
}
