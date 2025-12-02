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
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// SearchModel 记忆搜索模型
// 呀~ 搜索记忆的界面！🔍
type SearchModel struct {
	bs        *startup.Bootstrap
	input     textinput.Model
	list      list.Model
	results   []types.Memory
	searching bool
	width     int
	height    int
	err       error
}

// NewSearchModel 创建记忆搜索模型
func NewSearchModel(bs *startup.Bootstrap) *SearchModel {
	// 搜索输入框
	ti := textinput.New()
	ti.Placeholder = "输入关键词搜索..."
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 50

	// 结果列表
	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = styles.ListSelectedStyle
	delegate.Styles.SelectedDesc = styles.ListDescStyle
	delegate.Styles.NormalTitle = styles.ListItemStyle
	delegate.Styles.NormalDesc = styles.ListDescStyle

	l := list.New([]list.Item{}, delegate, 80, 15)
	l.Title = "🔍 搜索结果"
	l.SetShowHelp(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = styles.ListTitleStyle

	return &SearchModel{
		bs:    bs,
		input: ti,
		list:  l,
	}
}

// Title 返回页面标题
func (m *SearchModel) Title() string {
	return "搜索记忆"
}

// ShortHelp 返回快捷键帮助
func (m *SearchModel) ShortHelp() []key.Binding {
	return []key.Binding{common.KeyEnter, common.KeyBack}
}

// Init 初始化
func (m *SearchModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update 处理输入
func (m *SearchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, common.KeyBack):
			if m.input.Focused() {
				return m, common.Back()
			}
			m.input.Focus()
			return m, nil

		case key.Matches(msg, common.KeyEnter):
			if m.input.Focused() && m.input.Value() != "" {
				// 执行搜索
				m.searching = true
				return m, m.search(m.input.Value())
			}
			// 查看详情
			if item, ok := m.list.SelectedItem().(memoryItem); ok {
				return m, common.Navigate(common.PageMemoryDetail, map[string]any{"id": item.memory.ID})
			}

		case key.Matches(msg, common.KeyDown):
			if m.input.Focused() && len(m.results) > 0 {
				m.input.Blur()
			}

		case key.Matches(msg, common.KeyUp):
			if !m.input.Focused() && m.list.Index() == 0 {
				m.input.Focus()
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.input.Width = msg.Width - 10
		m.list.SetSize(msg.Width-4, msg.Height-12)

	case searchResultsMsg:
		m.searching = false
		m.results = msg.memories
		items := make([]list.Item, len(msg.memories))
		for i, memory := range msg.memories {
			items[i] = memoryItem{memory: memory}
		}
		m.list.SetItems(items)
		if len(items) > 0 {
			m.input.Blur()
		}

	case memoriesErrorMsg:
		m.searching = false
		m.err = msg.err
	}

	// 更新输入框
	if m.input.Focused() {
		newInput, cmd := m.input.Update(msg)
		m.input = newInput
		cmds = append(cmds, cmd)
	} else {
		// 更新列表
		newList, cmd := m.list.Update(msg)
		m.list = newList
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

type searchResultsMsg struct {
	memories []types.Memory
}

// search 搜索记忆
func (m *SearchModel) search(keyword string) tea.Cmd {
	return func() tea.Msg {
		memories, err := m.bs.MemoryService.SearchMemories(context.Background(), keyword)
		if err != nil {
			return memoriesErrorMsg{err}
		}
		return searchResultsMsg{memories}
	}
}

// View 渲染界面
func (m *SearchModel) View() string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("🔍 搜索记忆"))
	b.WriteString("\n\n")

	// 搜索框
	b.WriteString(styles.LabelStyle.Render("关键词"))
	b.WriteString("\n")
	b.WriteString(m.input.View())
	b.WriteString("\n\n")

	// 错误信息
	if m.err != nil {
		b.WriteString(styles.ErrorStyle.Render("错误: " + m.err.Error()))
		b.WriteString("\n\n")
	}

	// 搜索中
	if m.searching {
		b.WriteString(styles.InfoStyle.Render("搜索中..."))
		b.WriteString("\n\n")
	}

	// 搜索结果
	if len(m.results) > 0 {
		b.WriteString(styles.SubtitleStyle.Render(fmt.Sprintf("找到 %d 条结果", len(m.results))))
		b.WriteString("\n\n")
		b.WriteString(m.list.View())
	} else if !m.searching && m.input.Value() != "" && m.results != nil {
		b.WriteString(styles.MutedStyle.Render("未找到匹配的记忆"))
	}

	b.WriteString("\n")
	b.WriteString(styles.HelpStyle.Render("enter 搜索/查看 | ↑/↓ 选择 | esc 返回"))

	return b.String()
}

// 引入 utils 进行格式化
var _ = utils.FormatTime
