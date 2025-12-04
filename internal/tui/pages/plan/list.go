package plan

import (
	"fmt"
	"strings"
	"time"

	"github.com/XiaoLFeng/llm-memory/internal/models/entity"
	"github.com/XiaoLFeng/llm-memory/internal/tui/components"
	"github.com/XiaoLFeng/llm-memory/internal/tui/core"
	"github.com/XiaoLFeng/llm-memory/internal/tui/layout"
	"github.com/XiaoLFeng/llm-memory/internal/tui/theme"
	"github.com/XiaoLFeng/llm-memory/internal/tui/utils"
	"github.com/XiaoLFeng/llm-memory/startup"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type loadMsg struct {
	items []planItem
	err   error
}

type planItem struct {
	Title     string
	Status    string
	Progress  int
	Global    bool
	PathID    int64
	CreatedAt time.Time
}

type ListPage struct {
	bs          *startup.Bootstrap
	frame       *layout.Frame
	loading     bool
	err         error
	items       []planItem
	cursor      int
	showing     bool
	scopeFilter utils.ScopeFilter // 作用域过滤状态
}

func NewListPage(bs *startup.Bootstrap, _ func(core.PageID) tea.Cmd) *ListPage {
	return &ListPage{
		bs:      bs,
		frame:   layout.NewFrame(80, 24),
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
		plans, err := p.bs.PlanService.ListPlansByScope(ctx, scopeStr, p.bs.CurrentScope)
		if err != nil {
			return loadMsg{err: err}
		}
		items := make([]planItem, 0, len(plans))
		for _, pl := range plans {
			items = append(items, planItem{
				Title:     pl.Title,
				Status:    string(pl.Status),
				Progress:  pl.Progress,
				Global:    pl.Global,
				PathID:    pl.PathID,
				CreatedAt: pl.CreatedAt,
			})
		}
		return loadMsg{items: items}
	}
}

func (p *ListPage) Resize(w, h int) { p.frame.Resize(w, h) }

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
	cardW := layout.FitCardWidth(cw)
	scopeLabel := p.scopeFilter.Label()
	titleWithScope := fmt.Sprintf("%s 计划列表 [%s]", theme.IconPlan, scopeLabel)

	switch {
	case p.loading:
		return components.LoadingState(titleWithScope, "加载计划中...", cardW)
	case p.err != nil:
		return components.ErrorState(titleWithScope, p.err.Error(), cardW)
	case len(p.items) == 0:
		return components.EmptyState(titleWithScope, "暂无计划，按 c 创建吧~", cardW)
	default:
		if p.showing {
			body := p.renderDetail(cardW - 6)
			return components.Card(theme.IconPlan+" 计划详情", body, cardW)
		}
		body := p.renderList(cardW - 6)
		return components.Card(titleWithScope, body, cardW)
	}
}

func (p *ListPage) renderList(width int) string {
	var b strings.Builder
	max := len(p.items)
	if max > 20 {
		max = 20
	}
	for i := 0; i < max; i++ {
		pl := p.items[i]
		scope := utils.ScopeTag(pl.Global, pl.PathID, p.bs)
		status := statusText(pl.Status, pl.Progress)
		line := fmt.Sprintf("%s %s · %s · %d%% · %s",
			scope, pl.Title, status, pl.Progress, pl.CreatedAt.Format("01-02 15:04"))
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
	pl := p.items[p.cursor]
	scope := utils.ScopeTag(pl.Global, pl.PathID, p.bs)
	lines := []string{
		fmt.Sprintf("标题: %s", pl.Title),
		fmt.Sprintf("状态: %s", statusText(pl.Status, pl.Progress)),
		fmt.Sprintf("进度: %d%%", pl.Progress),
		fmt.Sprintf("作用域: %s", scope),
		fmt.Sprintf("创建时间: %s", pl.CreatedAt.Format("2006-01-02 15:04:05")),
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
		Title:      "计划列表",
		Breadcrumb: "计划管理 > 列表",
		Extra:      fmt.Sprintf("[%s] Tab切换 r刷新", p.scopeFilter.Label()),
		Keys: []components.KeyHint{
			{Key: "Tab", Desc: "切换作用域"},
			{Key: "Enter", Desc: "详情"},
			{Key: "c", Desc: "新建计划"},
			{Key: "r", Desc: "刷新"},
			{Key: "Esc", Desc: "返回"},
			{Key: "↑/↓", Desc: "移动"},
		},
	}
}

// statusText 将计划状态转换为中文显示
func statusText(status string, progress int) string {
	switch entity.PlanStatus(status) {
	case entity.PlanStatusCompleted:
		return "已完成"
	case entity.PlanStatusInProgress:
		return "进行中"
	case entity.PlanStatusCancelled:
		return "已取消"
	default:
		if progress > 0 {
			return "进行中"
		}
		return "待开始"
	}
}
