package memory

import (
	"github.com/XiaoLFeng/llm-memory/internal/tui/components"
	"github.com/XiaoLFeng/llm-memory/internal/tui/core"
	"github.com/XiaoLFeng/llm-memory/internal/tui/layout"
	"github.com/XiaoLFeng/llm-memory/internal/tui/theme"
	"github.com/XiaoLFeng/llm-memory/startup"
	tea "github.com/charmbracelet/bubbletea"
)

// ListPage 简化版：展示占位与快捷键
type ListPage struct {
	bs     *startup.Bootstrap
	frame  *layout.Frame
	width  int
	height int
}

func NewListPage(bs *startup.Bootstrap, _ func(core.PageID) tea.Cmd) *ListPage {
	return &ListPage{
		bs:    bs,
		frame: layout.NewFrame(80, 24),
		width: 80, height: 24,
	}
}

func (p *ListPage) Init() tea.Cmd { return nil }

func (p *ListPage) Resize(w, h int) {
	p.width, p.height = w, h
	p.frame.Resize(w, h)
}

func (p *ListPage) Update(msg tea.Msg) (core.Page, tea.Cmd) {
	return p, nil
}

func (p *ListPage) View() string {
	cw, _ := p.frame.ContentSize()
	card := components.EmptyState(theme.IconMemory+" 记忆列表", "尚未接入数据服务，按 Enter 可接入实现。", layout.FitCardWidth(cw))
	return card
}

func (p *ListPage) Meta() core.Meta {
	return core.Meta{
		Title:      "记忆列表",
		Breadcrumb: "记忆管理 > 列表",
		Extra:      "占位模式",
		Keys: []components.KeyHint{
			{Key: "Enter", Desc: "（预留）打开详情"},
			{Key: "c", Desc: "新建"},
			{Key: "Esc", Desc: "返回"},
		},
	}
}
