package app

import (
	"github.com/XiaoLFeng/llm-memory/internal/tui/components"
	"github.com/XiaoLFeng/llm-memory/internal/tui/core"
	"github.com/XiaoLFeng/llm-memory/internal/tui/layout"
	"github.com/XiaoLFeng/llm-memory/internal/tui/pages/group"
	"github.com/XiaoLFeng/llm-memory/internal/tui/pages/memory"
	"github.com/XiaoLFeng/llm-memory/internal/tui/pages/menu"
	"github.com/XiaoLFeng/llm-memory/internal/tui/pages/plan"
	"github.com/XiaoLFeng/llm-memory/internal/tui/pages/todo"
	"github.com/XiaoLFeng/llm-memory/startup"
	tea "github.com/charmbracelet/bubbletea"
)

// 导航消息
type navMsg struct{ to core.PageID }

// AppModel 根模型
type AppModel struct {
	bs     *startup.Bootstrap
	frame  *layout.Frame
	width  int
	height int

	page  core.Page
	stack []core.Page
}

func New(bs *startup.Bootstrap) *AppModel {
	m := &AppModel{
		bs:     bs,
		frame:  layout.NewFrame(80, 24),
		width:  80,
		height: 24,
	}
	m.page = m.makePage(core.PageMenu)
	return m
}

// Init
func (m *AppModel) Init() tea.Cmd {
	return m.page.Init()
}

// Update
func (m *AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch v := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = v.Width, v.Height
		m.frame.Resize(v.Width, v.Height)
		m.page.Resize(v.Width, v.Height)

	case navMsg:
		// push current
		m.stack = append(m.stack, m.page)
		m.page = m.makePage(v.to)
		return m, m.page.Init()

	case tea.KeyMsg:
		switch v.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "esc":
			if len(m.stack) > 0 {
				m.page = m.stack[len(m.stack)-1]
				m.stack = m.stack[:len(m.stack)-1]
				m.page.Resize(m.width, m.height)
				return m, nil
			}
		}
	}

	newPage, cmd := m.page.Update(msg)
	m.page = newPage
	return m, cmd
}

// View
func (m *AppModel) View() string {
	meta := m.page.Meta()
	keys := componentsToStrings(meta.Keys)
	return m.frame.Render(meta.Breadcrumb, meta.Extra, m.page.View(), keys)
}

// create page
func (m *AppModel) makePage(id core.PageID) core.Page {
	switch id {
	case core.PageMemory:
		return memory.NewListPage(m.bs, m.navigate)
	case core.PagePlan:
		return plan.NewListPage(m.bs, m.navigate)
	case core.PageTodo:
		return todo.NewListPage(m.bs, m.navigate)
	case core.PageGroup:
		return group.NewListPage(m.bs, m.navigate)
	default:
		return menu.NewPage(m.navigate)
	}
}

// navigate helper
func (m *AppModel) navigate(id core.PageID) tea.Cmd {
	return func() tea.Msg { return navMsg{to: id} }
}

// helper convert keys
func componentsToStrings(keys []components.KeyHint) []string {
	var out []string
	for _, k := range keys {
		out = append(out, components.JoinKeys([]components.KeyHint{k}))
	}
	return out
}
