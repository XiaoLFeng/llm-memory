package memory

import (
	"fmt"
	"strings"
	"time"

	"github.com/XiaoLFeng/llm-memory/internal/tui/components"
	"github.com/XiaoLFeng/llm-memory/internal/tui/core"
	"github.com/XiaoLFeng/llm-memory/internal/tui/layout"
	"github.com/XiaoLFeng/llm-memory/internal/tui/theme"
	"github.com/XiaoLFeng/llm-memory/internal/tui/utils"
	"github.com/XiaoLFeng/llm-memory/startup"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type (
	loadMsg struct {
		items []typesMemory
		err   error
	}
)

// typesMemory 只包含 TUI 展示需要的字段，避免直接耦合 entity
type typesMemory struct {
	Title     string
	Category  string
	Priority  int
	Global    bool
	PathID    int64
	Tags      []string
	CreatedAt time.Time
}

type ListPage struct {
	bs          *startup.Bootstrap
	frame       *layout.Frame
	width       int
	height      int
	loading     bool
	err         error
	items       []typesMemory
	cursor      int
	showing     bool              // true 展示详情，false 展示列表
	scopeFilter utils.ScopeFilter // 作用域过滤状态
}

func NewListPage(bs *startup.Bootstrap, _ func(core.PageID) tea.Cmd) *ListPage {
	return &ListPage{
		bs:      bs,
		frame:   layout.NewFrame(80, 24),
		width:   80,
		height:  24,
		loading: true,
	}
}

func (p *ListPage) Init() tea.Cmd {
	return p.load()
}

func (p *ListPage) load() tea.Cmd {
	return func() tea.Msg {
		ctx := p.bs.Context()
		scopeStr := p.scopeFilter.String()
		memories, err := p.bs.MemoryService.ListMemoriesByScope(ctx, scopeStr, p.bs.CurrentScope)
		if err != nil {
			return loadMsg{err: err}
		}
		items := make([]typesMemory, 0, len(memories))
		for _, m := range memories {
			items = append(items, typesMemory{
				Title:     m.Title,
				Category:  m.Category,
				Priority:  m.Priority,
				Global:    m.Global,
				PathID:    m.PathID,
				Tags:      m.GetTagStrings(),
				CreatedAt: m.CreatedAt,
			})
		}
		return loadMsg{items: items}
	}
}

func (p *ListPage) Resize(w, h int) {
	p.width, p.height = w, h
	p.frame.Resize(w, h)
}

func (p *ListPage) Update(msg tea.Msg) (core.Page, tea.Cmd) {
	switch v := msg.(type) {
	case tea.KeyMsg:
		switch v.String() {
		case "tab":
			p.scopeFilter = p.scopeFilter.Next()
			p.loading = true
			p.cursor = 0
			return p, p.load()
		case "r":
			p.loading = true
			p.err = nil
			return p, p.load()
		case "up", "k":
			if p.cursor > 0 {
				p.cursor--
			}
		case "down", "j":
			if p.cursor < len(p.items)-1 {
				p.cursor++
			}
		case "enter":
			p.showing = !p.showing
		case "esc":
			p.showing = false
		}
	case loadMsg:
		p.loading = false
		p.err = v.err
		if v.err == nil {
			p.items = v.items
			if p.cursor >= len(p.items) {
				p.cursor = len(p.items) - 1
			}
		}
	}
	return p, nil
}

func (p *ListPage) View() string {
	cw, _ := p.frame.ContentSize()
	cardWidth := layout.FitCardWidth(cw)
	scopeLabel := p.scopeFilter.Label()
	titleWithScope := fmt.Sprintf("%s 记忆列表 [%s]", theme.IconMemory, scopeLabel)

	switch {
	case p.loading:
		return components.LoadingState(titleWithScope, "努力加载中...", cardWidth)
	case p.err != nil:
		return components.ErrorState(titleWithScope, p.err.Error(), cardWidth)
	case len(p.items) == 0:
		return components.EmptyState(titleWithScope, "暂无记忆，按 c 创建一条吧~", cardWidth)
	default:
		if p.showing {
			body := p.renderDetail(cardWidth - 6)
			return components.Card(theme.IconMemory+" 记忆详情", body, cardWidth)
		}
		body := p.renderList(cardWidth - 6)
		return components.Card(titleWithScope, body, cardWidth)
	}
}

func (p *ListPage) renderList(width int) string {
	var b strings.Builder
	max := len(p.items)
	if max > 20 {
		max = 20
	}
	for i := 0; i < max; i++ {
		m := p.items[i]
		scope := utils.ScopeTag(m.Global, m.PathID, p.bs)
		tagStr := ""
		if len(m.Tags) > 0 {
			tagStr = " #" + strings.Join(m.Tags, " #")
		}
		line := fmt.Sprintf("%s %s · %s · P%d · %s%s",
			scope, m.Title, m.Category, m.Priority,
			m.CreatedAt.Format("01-02 15:04"), tagStr)
		if utils.LipWidth(line) > width {
			line = utils.Truncate(line, width)
		}
		if i == p.cursor {
			line = lipgloss.NewStyle().Foreground(theme.Info).Render("▶ " + line)
		} else {
			line = "  " + line
		}
		b.WriteString(line)
		if i != max-1 {
			b.WriteRune('\n')
		}
	}
	return b.String()
}

func (p *ListPage) renderDetail(width int) string {
	if len(p.items) == 0 {
		return "暂无数据"
	}
	m := p.items[p.cursor]
	scope := utils.ScopeTag(m.Global, m.PathID, p.bs)
	tagStr := "无"
	if len(m.Tags) > 0 {
		tagStr = strings.Join(m.Tags, ", ")
	}
	lines := []string{
		fmt.Sprintf("标题: %s", m.Title),
		fmt.Sprintf("分类: %s", m.Category),
		fmt.Sprintf("优先级: P%d", m.Priority),
		fmt.Sprintf("作用域: %s", scope),
		fmt.Sprintf("标签: %s", tagStr),
		fmt.Sprintf("创建时间: %s", m.CreatedAt.Format("2006-01-02 15:04:05")),
	}
	for i, l := range lines {
		if utils.LipWidth(l) > width {
			lines[i] = utils.Truncate(l, width)
		}
	}
	return strings.Join(lines, "\n")
}

func (p *ListPage) Meta() core.Meta {
	return core.Meta{
		Title:      "记忆列表",
		Breadcrumb: "记忆管理 > 列表",
		Extra:      fmt.Sprintf("[%s] Tab切换 r刷新", p.scopeFilter.Label()),
		Keys: []components.KeyHint{
			{Key: "Tab", Desc: "切换作用域"},
			{Key: "Enter", Desc: "详情"},
			{Key: "c", Desc: "新建"},
			{Key: "r", Desc: "刷新"},
			{Key: "Esc", Desc: "返回"},
			{Key: "↑/↓", Desc: "移动"},
		},
	}
}
