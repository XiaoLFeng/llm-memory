package plan

import (
	"github.com/XiaoLFeng/llm-memory/internal/tui/components"
	"github.com/XiaoLFeng/llm-memory/internal/tui/core"
	"github.com/XiaoLFeng/llm-memory/internal/tui/layout"
	"github.com/XiaoLFeng/llm-memory/internal/tui/theme"
	"github.com/XiaoLFeng/llm-memory/startup"
	tea "github.com/charmbracelet/bubbletea"
)

type ListPage struct {
	bs    *startup.Bootstrap
	frame *layout.Frame
}

func NewListPage(bs *startup.Bootstrap, _ func(core.PageID) tea.Cmd) *ListPage {
	return &ListPage{
		bs:    bs,
		frame: layout.NewFrame(80, 24),
	}
}

func (p *ListPage) Init() tea.Cmd                           { return nil }
func (p *ListPage) Resize(w, h int)                         { p.frame.Resize(w, h) }
func (p *ListPage) Update(msg tea.Msg) (core.Page, tea.Cmd) { return p, nil }

func (p *ListPage) View() string {
	cw, _ := p.frame.ContentSize()
	return components.EmptyState(theme.IconPlan+" 计划列表", "计划页面待接入，保持期待~", layout.FitCardWidth(cw))
}

func (p *ListPage) Meta() core.Meta {
	return core.Meta{
		Title:      "计划列表",
		Breadcrumb: "计划管理 > 列表",
		Extra:      "占位模式",
		Keys: []components.KeyHint{
			{Key: "Enter", Desc: "查看计划"},
			{Key: "c", Desc: "新建计划"},
			{Key: "Esc", Desc: "返回"},
		},
	}
}
