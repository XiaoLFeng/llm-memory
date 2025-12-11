package app

import (
	"github.com/XiaoLFeng/llm-memory/internal/tui/components"
	"github.com/XiaoLFeng/llm-memory/internal/tui/core"
	"github.com/XiaoLFeng/llm-memory/internal/tui/layout"
	"github.com/XiaoLFeng/llm-memory/internal/tui/pages/group"
	"github.com/XiaoLFeng/llm-memory/internal/tui/pages/help"
	"github.com/XiaoLFeng/llm-memory/internal/tui/pages/memory"
	"github.com/XiaoLFeng/llm-memory/internal/tui/pages/menu"
	"github.com/XiaoLFeng/llm-memory/internal/tui/pages/plan"
	"github.com/XiaoLFeng/llm-memory/internal/tui/pages/todo"
	"github.com/XiaoLFeng/llm-memory/startup"
	tea "github.com/charmbracelet/bubbletea"
)

// 导航消息（支持携带数据）
type navMsg struct {
	to   core.PageID
	data interface{}
}

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
		m.page = m.makePageWithData(v.to, v.data)
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
	return m.makePageWithData(id, nil)
}

// create page with data
func (m *AppModel) makePageWithData(id core.PageID, data interface{}) core.Page {
	switch id {
	case core.PageMemory:
		return memory.NewListPage(m.bs, m.navigate, m.navigateWithData)
	case core.PageMemoryCreate:
		return memory.NewCreatePage(m.bs, m.navigate)
	case core.PageMemoryEdit:
		if memoryID, ok := data.(int64); ok {
			return memory.NewEditPage(m.bs, memoryID, m.navigate)
		}
		return memory.NewListPage(m.bs, m.navigate, m.navigateWithData)
	case core.PagePlan:
		return plan.NewListPage(m.bs, m.navigate, m.navigateWithData)
	case core.PagePlanCreate:
		return plan.NewCreatePage(m.bs, m.navigate)
	case core.PagePlanEdit:
		if planID, ok := data.(int64); ok {
			return plan.NewEditPage(m.bs, planID, m.navigate)
		}
		return plan.NewListPage(m.bs, m.navigate, m.navigateWithData)
	case core.PageTodo:
		return todo.NewListPage(m.bs, m.navigate, m.navigateWithData)
	case core.PageTodoCreate:
		// 支持从 Plan 详情页创建 Todo（传递 TodoCreateContext）
		if ctx, ok := data.(*plan.TodoCreateContext); ok && ctx != nil {
			return todo.NewCreatePageWithPlan(m.bs, m.navigate, ctx.PlanCode, ctx.PlanTitle)
		}
		return todo.NewCreatePage(m.bs, m.navigate)
	case core.PageTodoEdit:
		if todoID, ok := data.(int64); ok {
			return todo.NewEditPage(m.bs, todoID, m.navigate)
		}
		return todo.NewListPage(m.bs, m.navigate, m.navigateWithData)
	case core.PageGroup:
		return group.NewListPage(m.bs, m.navigate, m.navigateWithData)
	case core.PageGroupCreate:
		return group.NewCreatePage(m.bs, m.navigate)
	case core.PageGroupEdit:
		return group.NewEditPage(m.bs, m.navigate, data)
	case core.PageHelp:
		return help.NewPage(m.navigate)
	default:
		return menu.NewPage(m.navigate)
	}
}

// navigate helper
func (m *AppModel) navigate(id core.PageID) tea.Cmd {
	return func() tea.Msg { return navMsg{to: id} }
}

// navigateWithData helper
func (m *AppModel) navigateWithData(id core.PageID, data interface{}) tea.Cmd {
	return func() tea.Msg { return navMsg{to: id, data: data} }
}

// helper convert keys
func componentsToStrings(keys []components.KeyHint) []string {
	var out []string
	for _, k := range keys {
		out = append(out, components.JoinKeys([]components.KeyHint{k}))
	}
	return out
}
