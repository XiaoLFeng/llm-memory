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
	"github.com/XiaoLFeng/llm-memory/internal/tui/utils"
	"github.com/XiaoLFeng/llm-memory/startup"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type loadMsg struct {
	items []todoItem
	err   error
}

type todoItem struct {
	ID          int64 // 添加ID字段用于编辑和删除
	Title       string
	Priority    entity.ToDoPriority
	Status      entity.ToDoStatus
	Global      bool
	PathID      int64
	DueDate     *time.Time
	CompletedAt *time.Time
}

type ListPage struct {
	bs              *startup.Bootstrap
	frame           *layout.Frame
	loading         bool
	err             error
	items           []todoItem
	cursor          int
	showing         bool
	scopeFilter     utils.ScopeFilter // 作用域过滤状态
	push            func(core.PageID) tea.Cmd
	pushWithData    func(core.PageID, interface{}) tea.Cmd
	confirmDelete   bool
	deleteTarget    int64
	deleteYesActive bool // true=选中确认，false=选中取消
}

func NewListPage(bs *startup.Bootstrap, push func(core.PageID) tea.Cmd, pushWithData func(core.PageID, interface{}) tea.Cmd) *ListPage {
	return &ListPage{
		bs:           bs,
		frame:        layout.NewFrame(80, 24),
		loading:      true,
		push:         push,
		pushWithData: pushWithData,
	}
}

func (p *ListPage) Init() tea.Cmd { return p.load() }

func (p *ListPage) load() tea.Cmd {
	return func() tea.Msg {
		ctx := p.bs.Context()
		scopeStr := p.scopeFilter.String()
		todos, err := p.bs.ToDoService.ListToDosByScope(ctx, scopeStr, p.bs.CurrentScope)
		if err != nil {
			return loadMsg{err: err}
		}
		items := make([]todoItem, 0, len(todos))
		for _, t := range todos {
			items = append(items, todoItem{
				ID:          t.ID, // 添加ID
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
		// 删除确认模式
		if p.confirmDelete {
			switch v.String() {
			case "left", "right", "h", "l":
				p.deleteYesActive = !p.deleteYesActive
			case "y", "enter":
				if p.deleteYesActive {
					return p, p.doDelete()
				}
				p.confirmDelete = false
				p.deleteTarget = 0
			case "n", "esc":
				p.confirmDelete = false
				p.deleteTarget = 0
			}
			return p, nil
		}

		// 普通模式
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
		case "c":
			// 创建待办
			if p.push != nil {
				return p, p.push(core.PageTodoCreate)
			}
		case "e":
			// 编辑待办
			if len(p.items) > 0 && p.pushWithData != nil {
				todoID := p.items[p.cursor].ID
				return p, p.pushWithData(core.PageTodoEdit, todoID)
			}
		case "d":
			// 删除待办
			if len(p.items) > 0 {
				p.confirmDelete = true
				p.deleteTarget = p.items[p.cursor].ID
				p.deleteYesActive = false
			}
		case "?":
			// 查看帮助（可选实现）
			if p.push != nil {
				return p, p.push(core.PageHelp)
			}
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
	titleWithScope := fmt.Sprintf("%s 待办列表 [%s]", theme.IconTodo, scopeLabel)

	// 删除确认对话框
	if p.confirmDelete && len(p.items) > 0 {
		itemName := p.items[p.cursor].Title
		dialog := components.DeleteConfirmDialog(itemName, cardW)
		return dialog
	}

	switch {
	case p.loading:
		return components.LoadingState(titleWithScope, "加载待办中...", cardW)
	case p.err != nil:
		return components.ErrorState(titleWithScope, p.err.Error(), cardW)
	case len(p.items) == 0:
		return components.EmptyState(titleWithScope, "暂无待办，按 c 创建吧~", cardW)
	default:
		if p.showing {
			body := p.renderDetail(cardW - 6)
			return components.Card(theme.IconTodo+" 待办详情", body, cardW)
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
		t := p.items[i]
		scope := utils.ScopeTag(t.Global, t.PathID, p.bs)
		status := statusText(t.Status)
		priority := priorityText(t.Priority)
		due := ""
		if t.DueDate != nil {
			due = " · 截止 " + t.DueDate.Format("01-02")
		}
		line := fmt.Sprintf("%s %s · %s · %s%s",
			scope, t.Title, priority, status, due)
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
	t := p.items[p.cursor]
	scope := utils.ScopeTag(t.Global, t.PathID, p.bs)
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
		if utils.LipWidth(l) > width {
			lines[i] = utils.Truncate(l, width)
		}
	}
	return strings.Join(lines, "\n")
}

func (p *ListPage) Meta() core.Meta {
	return core.Meta{
		Title:      "待办列表",
		Breadcrumb: "待办管理 > 列表",
		Extra:      fmt.Sprintf("[%s] Tab切换 r刷新", p.scopeFilter.Label()),
		Keys: []components.KeyHint{
			{Key: "Tab", Desc: "切换作用域"},
			{Key: "Enter", Desc: "详情"},
			{Key: "c", Desc: "新建"},
			{Key: "e", Desc: "编辑"},
			{Key: "d", Desc: "删除"},
			{Key: "r", Desc: "刷新"},
			{Key: "Esc", Desc: "返回"},
			{Key: "↑/↓", Desc: "移动"},
		},
	}
}

// doDelete 执行删除操作
func (p *ListPage) doDelete() tea.Cmd {
	return func() tea.Msg {
		ctx := p.bs.Context()
		if err := p.bs.ToDoService.DeleteToDoByID(ctx, p.deleteTarget); err != nil {
			p.err = err
			p.confirmDelete = false
			p.deleteTarget = 0
			return nil
		}
		// 删除成功，重新加载
		p.confirmDelete = false
		p.deleteTarget = 0
		return p.load()()
	}
}

// statusText 将待办状态转换为中文显示
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

// priorityText 将待办优先级转换为中文显示
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
