package memory

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

// memoryItem 记忆列表项
type memoryItem struct {
	memory types.Memory
}

func (i memoryItem) Title() string {
	return fmt.Sprintf("%d. %s", i.memory.ID, i.memory.Title)
}

func (i memoryItem) Description() string {
	return fmt.Sprintf("[%s] %s", i.memory.Category, utils.FormatRelativeTime(i.memory.CreatedAt))
}

func (i memoryItem) FilterValue() string {
	return i.memory.Title
}

// ListModel 记忆列表模型
// 嘿嘿~ 展示所有记忆的列表！📚
type ListModel struct {
	bs       *startup.Bootstrap
	list     list.Model
	memories []types.Memory
	width    int
	height   int
	loading  bool
	err      error
}

// NewListModel 创建记忆列表模型
func NewListModel(bs *startup.Bootstrap) *ListModel {
	// 创建列表
	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = styles.ListSelectedStyle
	delegate.Styles.SelectedDesc = styles.ListDescStyle
	delegate.Styles.NormalTitle = styles.ListItemStyle
	delegate.Styles.NormalDesc = styles.ListDescStyle

	l := list.New([]list.Item{}, delegate, 80, 20)
	l.Title = "📚 记忆列表"
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
	return "记忆列表"
}

// ShortHelp 返回快捷键帮助
func (m *ListModel) ShortHelp() []key.Binding {
	return []key.Binding{
		common.KeyUp, common.KeyDown, common.KeyEnter,
		common.KeyCreate, common.KeyDelete, common.KeySearch, common.KeyBack,
	}
}

// Init 初始化
func (m *ListModel) Init() tea.Cmd {
	return m.loadMemories()
}

// loadMemories 加载记忆列表
func (m *ListModel) loadMemories() tea.Cmd {
	return func() tea.Msg {
		memories, err := m.bs.MemoryService.ListMemories(context.Background())
		if err != nil {
			return memoriesErrorMsg{err}
		}
		return memoriesLoadedMsg{memories}
	}
}

type memoriesLoadedMsg struct {
	memories []types.Memory
}

type memoriesErrorMsg struct {
	err error
}

type memoryDeletedMsg struct {
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
			return m, common.Navigate(common.PageMemoryCreate)

		case key.Matches(msg, common.KeySearch):
			return m, common.Navigate(common.PageMemorySearch)

		case key.Matches(msg, common.KeyEnter):
			if item, ok := m.list.SelectedItem().(memoryItem); ok {
				return m, common.Navigate(common.PageMemoryDetail, map[string]any{"id": item.memory.ID})
			}

		case key.Matches(msg, common.KeyDelete):
			if item, ok := m.list.SelectedItem().(memoryItem); ok {
				return m, common.ShowConfirm(
					"删除记忆",
					fmt.Sprintf("确定要删除记忆「%s」吗？", item.memory.Title),
					m.deleteMemory(item.memory.ID),
					nil,
				)
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetSize(msg.Width-4, msg.Height-8)

	case memoriesLoadedMsg:
		m.loading = false
		m.memories = msg.memories
		items := make([]list.Item, len(msg.memories))
		for i, memory := range msg.memories {
			items[i] = memoryItem{memory: memory}
		}
		m.list.SetItems(items)

	case memoriesErrorMsg:
		m.loading = false
		m.err = msg.err

	case memoryDeletedMsg:
		cmds = append(cmds, m.loadMemories())
		cmds = append(cmds, common.ShowToast("记忆已删除", common.ToastSuccess))

	case common.RefreshMsg:
		m.loading = true
		cmds = append(cmds, m.loadMemories())
	}

	// 更新列表
	newList, cmd := m.list.Update(msg)
	m.list = newList
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// deleteMemory 删除记忆
func (m *ListModel) deleteMemory(id int) tea.Cmd {
	return func() tea.Msg {
		err := m.bs.MemoryService.DeleteMemory(context.Background(), id)
		if err != nil {
			return memoriesErrorMsg{err}
		}
		return memoryDeletedMsg{id}
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

	if len(m.memories) == 0 {
		b.WriteString(styles.TitleStyle.Render("📚 记忆列表"))
		b.WriteString("\n\n")
		b.WriteString(styles.MutedStyle.Render("暂无记忆~ 按 c 创建新记忆"))
		b.WriteString("\n\n")
		b.WriteString(styles.HelpStyle.Render("c 新建 | / 搜索 | esc 返回"))
		return b.String()
	}

	b.WriteString(m.list.View())
	b.WriteString("\n")
	b.WriteString(styles.HelpStyle.Render("↑/↓ 选择 | enter 查看 | c 新建 | d 删除 | / 搜索 | esc 返回"))

	return b.String()
}
