package menu

import (
	"strings"

	"github.com/XiaoLFeng/llm-memory/internal/tui/components"
	"github.com/XiaoLFeng/llm-memory/internal/tui/core"
	"github.com/XiaoLFeng/llm-memory/internal/tui/layout"
	"github.com/XiaoLFeng/llm-memory/internal/tui/theme"
	tea "github.com/charmbracelet/bubbletea"
)

type Page struct {
	width  int
	height int
	frame  *layout.Frame
	cursor int
	items  []item
	push   func(core.PageID) tea.Cmd
}

type item struct {
	title string
	desc  string
	id    core.PageID
	icon  string
}

func NewPage(push func(core.PageID) tea.Cmd) *Page {
	return &Page{
		width:  80,
		height: 24,
		frame:  layout.NewFrame(80, 24),
		push:   push,
		items: []item{
			{title: "记忆管理", desc: "管理和搜索你的记忆", id: core.PageMemory, icon: theme.IconMemory},
			{title: "计划管理", desc: "规划与跟踪计划", id: core.PagePlan, icon: theme.IconPlan},
			{title: "待办管理", desc: "日常待办与执行", id: core.PageTodo, icon: theme.IconTodo},
			{title: "组管理", desc: "路径组与共享", id: core.PageGroup, icon: theme.IconGroup},
		},
	}
}

func (p *Page) Init() tea.Cmd { return nil }

func (p *Page) Resize(w, h int) {
	p.width, p.height = w, h
	p.frame.Resize(w, h)
}

func (p *Page) Update(msg tea.Msg) (core.Page, tea.Cmd) {
	switch m := msg.(type) {
	case tea.KeyMsg:
		switch m.String() {
		case "up", "k":
			if p.cursor > 0 {
				p.cursor--
			}
		case "down", "j":
			if p.cursor < len(p.items)-1 {
				p.cursor++
			}
		case "enter":
			return p, p.push(p.items[p.cursor].id)
		}
	}
	return p, nil
}

func (p *Page) View() string {
	contentW, _ := p.frame.ContentSize()
	cw := layout.FitCardWidth(contentW)
	return p.renderList(cw)
}

func (p *Page) Meta() core.Meta {
	return core.Meta{
		Title:      "主菜单",
		Breadcrumb: "主菜单",
		Extra:      "",
		Keys: []components.KeyHint{
			{Key: "↑/↓", Desc: "选择"},
			{Key: "Enter", Desc: "确认"},
		},
	}
}

// renderList
func (p *Page) renderList(cardWidth int) string {
	var lines []string
	for i, it := range p.items {
		cursor := "  "
		style := theme.TextMain
		descStyle := theme.TextDim
		if i == p.cursor {
			cursor = "▸ "
			style = theme.TextMain.Copy().Bold(true)
			descStyle = theme.TextDim.Copy().Foreground(theme.Primary)
		}
		title := style.Render(it.icon + " " + it.title)
		desc := descStyle.Render(it.desc)
		lines = append(lines, cursor+title+"\n  "+desc)
	}
	body := strings.Join(lines, "\n\n")
	return components.Card("选择功能", body, cardWidth)
}
