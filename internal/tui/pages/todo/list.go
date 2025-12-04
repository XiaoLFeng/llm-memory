package todo

import (
	"fmt"
	"strings"
	"time"

	"github.com/XiaoLFeng/llm-memory/internal/models/entity"
	"github.com/XiaoLFeng/llm-memory/internal/tui/components"
	"github.com/XiaoLFeng/llm-memory/internal/tui/core"
	"github.com/XiaoLFeng/llm-memory/internal/tui/layout"
	"github.com/XiaoLFeng/llm-memory/internal/tui/theme"
	"github.com/XiaoLFeng/llm-memory/startup"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type loadMsg struct {
	items []todoItem
	err   error
}

type todoItem struct {
	Title       string
	Priority    entity.ToDoPriority
	Status      entity.ToDoStatus
	Global      bool
	PathID      int64
	DueDate     *time.Time
	CompletedAt *time.Time
}

type ListPage struct {
	bs      *startup.Bootstrap
	frame   *layout.Frame
	loading bool
	err     error
	items   []todoItem
	cursor  int
	showing bool
}

func NewListPage(bs *startup.Bootstrap, _ func(core.PageID) tea.Cmd) *ListPage {
	return &ListPage{
		bs:      bs,
		frame:   layout.NewFrame(80, 24),
		loading: true,
	}
}

func (p *ListPage) Init() tea.Cmd { return p.load() }

func (p *ListPage) load() tea.Cmd {
	return func() tea.Msg {
		ctx := p.bs.Context()
		todos, err := p.bs.ToDoService.ListToDos(ctx)
		if err != nil {
			return loadMsg{err: err}
		}
		items := make([]todoItem, 0, len(todos))
		for _, t := range todos {
			items = append(items, todoItem{
				Title:       t.Title,
				Priority:    t.Priority,
				Status:      t.Status,
				Global:      t.Global,
				PathID:      t.PathID,
				DueDate:     t.DueDate,
				CompletedAt: t.CompletedAt,
			})
		}
		return loadMsg{items: items}
	}
}

func (p *ListPage) Resize(w, h int) { p.frame.Resize(w, h) }

func (p *ListPage) Update(msg tea.Msg) (core.Page, tea.Cmd) {
	switch v := msg.(type) {
	case tea.KeyMsg:
		if v.String() == "r" {
			p.loading = true
			p.err = nil
			return p, p.load()
		}
		switch v.String() {
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
	switch {
	case p.loading:
		return components.LoadingState(theme.IconTodo+" 待办列表", "加载待办中...", cardW)
	case p.err != nil:
		return components.ErrorState(theme.IconTodo+" 待办列表", p.err.Error(), cardW)
	case len(p.items) == 0:
		return components.EmptyState(theme.IconTodo+" 待办列表", "暂无待办，按 c 创建吧~", cardW)
	default:
		if p.showing {
			body := p.renderDetail(cardW - 6)
			return components.Card(theme.IconTodo+" 待办详情", body, cardW)
		}
		body := p.renderList(cardW - 6)
		return components.Card(theme.IconTodo+" 待办列表", body, cardW)
	}
}

func (p *ListPage) renderList(width int) string {
	var b strings.Builder
	max := len(p.items)
	if max > 20 {
		max = 20
	}
	for i := 0; i < max; i++ {
		t := p.items[i]
		scope := scopeTag(t.Global, t.PathID, p.bs)
		status := statusText(t.Status)
		priority := priorityText(t.Priority)
		due := ""
		if t.DueDate != nil {
			due = " · 截止 " + t.DueDate.Format("01-02")
		}
		line := fmt.Sprintf("%s %s · %s · %s%s",
			scope, t.Title, priority, status, due)
		if lipWidth(line) > width {
			line = truncate(line, width)
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
	t := p.items[p.cursor]
	scope := scopeTag(t.Global, t.PathID, p.bs)
	due := "无"
	if t.DueDate != nil {
		due = t.DueDate.Format("2006-01-02 15:04")
	}
	comp := "未完成"
	if t.CompletedAt != nil {
		comp = t.CompletedAt.Format("2006-01-02 15:04")
	}
	lines := []string{
		fmt.Sprintf("标题: %s", t.Title),
		fmt.Sprintf("优先级: %s", priorityText(t.Priority)),
		fmt.Sprintf("状态: %s", statusText(t.Status)),
		fmt.Sprintf("作用域: %s", scope),
		fmt.Sprintf("截止: %s", due),
		fmt.Sprintf("完成时间: %s", comp),
	}
	for i, l := range lines {
		if lipWidth(l) > width {
			lines[i] = truncate(l, width)
		}
	}
	return strings.Join(lines, "\n")
}

func (p *ListPage) Meta() core.Meta {
	return core.Meta{
		Title:      "待办列表",
		Breadcrumb: "待办管理 > 列表",
		Extra:      "r 刷新",
		Keys: []components.KeyHint{
			{Key: "Enter", Desc: "切换详情/列表"},
			{Key: "c", Desc: "新建"},
			{Key: "r", Desc: "刷新"},
			{Key: "Esc", Desc: "返回"},
			{Key: "↑/↓", Desc: "移动"},
		},
	}
}

func scopeTag(global bool, pathID int64, bs *startup.Bootstrap) string {
	if global {
		return "[全局]"
	}
	ctx := bs.CurrentScope
	if ctx != nil {
		for _, gid := range ctx.GroupPathIDs {
			if pathID == gid {
				return "[小组]"
			}
		}
		if ctx.PathID == pathID {
			return "[私有]"
		}
	}
	if pathID > 0 {
		return "[私有]"
	}
	return "[未知]"
}

func statusText(status entity.ToDoStatus) string {
	switch status {
	case entity.ToDoStatusCompleted:
		return "已完成"
	case entity.ToDoStatusInProgress:
		return "进行中"
	case entity.ToDoStatusCancelled:
		return "已取消"
	default:
		return "待处理"
	}
}

func priorityText(p entity.ToDoPriority) string {
	switch p {
	case entity.ToDoPriorityUrgent:
		return "紧急"
	case entity.ToDoPriorityHigh:
		return "高"
	case entity.ToDoPriorityMedium:
		return "中"
	default:
		return "低"
	}
}

func lipWidth(s string) int { return len([]rune(s)) }

func truncate(s string, limit int) string {
	runes := []rune(s)
	if len(runes) <= limit {
		return s
	}
	if limit <= 1 {
		return string(runes[:limit])
	}
	return string(runes[:limit-1]) + "…"
}
