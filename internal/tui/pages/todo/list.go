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
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type (
	loadMsg struct {
		items []todoItem
		err   error
	}
	deleteSuccessMsg    struct{}              // 删除成功消息
	deleteErrorMsg      struct{ err error }   // 删除失败消息
	finalDeleteMsg      struct{ count int64 } // 批量删除成功消息
	finalDeleteErrorMsg struct{ err error }   // 批量删除失败消息
)

type todoItem struct {
	ID          int64 // 添加ID字段用于编辑和删除
	Title       string
	Description string // 待办描述
	Priority    entity.ToDoPriority
	Status      entity.ToDoStatus
	PathID      int64
	DueDate     *time.Time
	CompletedAt *time.Time
}

type ListPage struct {
	bs               *startup.Bootstrap
	frame            *layout.Frame
	loading          bool
	err              error
	items            []todoItem
	cursor           int
	showing          bool
	scopeFilter      utils.ScopeFilter // 作用域过滤状态
	detailViewport   viewport.Model    // 详情页滚动视图
	push             func(core.PageID) tea.Cmd
	pushWithData     func(core.PageID, interface{}) tea.Cmd
	confirmDelete    bool  // 是否在删除确认模式
	deleteTarget     int64 // 要删除的 ID
	deleteProcessing bool  // 是否正在处理删除
	deleteYesActive  bool  // true=选中确认，false=选中取消
	confirmFinal     bool  // 是否在批量删除确认模式
	finalProcessing  bool  // 是否正在处理批量删除
	finalYesActive   bool  // true=选中确认，false=选中取消
}

func NewListPage(bs *startup.Bootstrap, push func(core.PageID) tea.Cmd, pushWithData func(core.PageID, interface{}) tea.Cmd) *ListPage {
	// 初始化 viewport（初始尺寸，后续动态调整）
	vp := viewport.New(60, 10)
	vp.Style = lipgloss.NewStyle()

	return &ListPage{
		bs:             bs,
		frame:          layout.NewFrame(80, 24),
		loading:        true,
		detailViewport: vp,
		push:           push,
		pushWithData:   pushWithData,
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
				Description: t.Description, // 添加描述字段
				Priority:    t.Priority,
				Status:      t.Status,
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
		// 批量删除确认模式
		if p.confirmFinal {
			switch v.String() {
			case "left", "h":
				p.finalYesActive = !p.finalYesActive
				return p, nil
			case "right", "l":
				p.finalYesActive = !p.finalYesActive
				return p, nil
			case "enter":
				if p.finalYesActive {
					p.confirmFinal = false
					p.finalProcessing = true
					return p, p.doFinalDelete()
				} else {
					p.confirmFinal = false
					return p, nil
				}
			case "esc", "n", "N":
				p.confirmFinal = false
				return p, nil
			}
			return p, nil
		}

		// 删除确认模式
		if p.confirmDelete {
			switch v.String() {
			case "y", "Y", "enter":
				p.confirmDelete = false
				p.deleteProcessing = true
				return p, p.doDelete()
			case "n", "N", "esc":
				p.confirmDelete = false
				p.deleteTarget = 0
				return p, nil
			}
			return p, nil
		}

		// 详情页模式：处理滚动
		if p.showing {
			switch v.String() {
			case "esc", "q":
				p.showing = false
				return p, nil
			case "up", "k":
				p.detailViewport.LineUp(1)
				return p, nil
			case "down", "j":
				p.detailViewport.LineDown(1)
				return p, nil
			case "pgup":
				p.detailViewport.HalfViewUp()
				return p, nil
			case "pgdown":
				p.detailViewport.HalfViewDown()
				return p, nil
			case "home":
				p.detailViewport.GotoTop()
				return p, nil
			case "end":
				p.detailViewport.GotoBottom()
				return p, nil
			}
			return p, nil
		}

		// 列表模式
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
			if len(p.items) > 0 {
				p.showing = !p.showing
				// 进入详情页时重置滚动位置
				if p.showing {
					p.detailViewport.GotoTop()
				}
			}
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
		case "f", "F":
			// 批量删除所有待办（final）
			if len(p.items) > 0 {
				p.confirmFinal = true
				p.finalYesActive = false
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
	case deleteSuccessMsg:
		p.deleteProcessing = false
		p.deleteTarget = 0
		p.loading = true
		return p, p.load()
	case deleteErrorMsg:
		p.deleteProcessing = false
		p.deleteTarget = 0
		p.err = v.err
	case finalDeleteMsg:
		p.finalProcessing = false
		p.loading = true
		p.cursor = 0
		return p, p.load()
	case finalDeleteErrorMsg:
		p.finalProcessing = false
		p.err = v.err
	case tea.WindowSizeMsg:
		// 动态调整 viewport 尺寸
		if p.showing {
			const headerHeight = 4 // 标题 + 空行
			const footerHeight = 3 // 空行 + 操作提示
			p.detailViewport.Width = v.Width - 4
			p.detailViewport.Height = v.Height - headerHeight - footerHeight
		}
	}
	return p, nil
}

func (p *ListPage) View() string {
	cw, _ := p.frame.ContentSize()
	cardW := layout.FitCardWidth(cw)
	scopeLabel := p.scopeFilter.Label()
	titleWithScope := fmt.Sprintf("%s 待办列表 [%s]", theme.IconTodo, scopeLabel)

	// 批量删除确认对话框
	if p.confirmFinal {
		totalCount := len(p.items)
		title := "批量删除确认"
		message := fmt.Sprintf("确定要删除当前作用域的所有 %d 个待办吗？此操作不可恢复！", totalCount)
		dialog := components.ConfirmDialogWithButtons(title, message, cardW, p.finalYesActive)
		return dialog
	}

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
			// === 使用 viewport 渲染详情页 ===
			// 动态计算并设置 viewport 尺寸
			cw, ch := p.frame.ContentSize()
			const headerHeight = 4 // 标题 + 空行
			const footerHeight = 3 // 空行 + 操作提示

			viewportWidth := cw - 4
			viewportHeight := ch - headerHeight - footerHeight

			p.detailViewport.Width = viewportWidth
			p.detailViewport.Height = viewportHeight

			// 生成详情内容并设置到 viewport
			detailContent := p.renderDetail(p.detailViewport.Width)
			p.detailViewport.SetContent(detailContent)

			// 滚动进度指示器
			scrollPercent := p.detailViewport.ScrollPercent() * 100
			scrollInfo := fmt.Sprintf("%.0f%%", scrollPercent)
			scrollHint := theme.TextDim.Render(fmt.Sprintf(
				"滚动: %s | ↑/↓ j/k PgUp/PgDn Home/End | Esc 返回", scrollInfo))

			// 组合视图
			title := theme.Title.Render(theme.IconTodo + " 待办详情")
			viewportView := p.detailViewport.View()

			return lipgloss.JoinVertical(lipgloss.Left,
				title,
				"",
				viewportView,
				"",
				scrollHint,
			)
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
		scope := utils.ScopeTag(t.PathID, p.bs)
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
	scope := utils.ScopeTag(t.PathID, p.bs)

	var lines []string

	// === 区块 1：标题 ===
	titleLine := theme.FormLabel.Bold(true).Render("标题: ") + theme.TextMain.Render(t.Title)
	lines = append(lines, titleLine)
	lines = append(lines, "")

	// === 区块 2：元数据 ===
	metaStyle := theme.TextDim
	due := "无"
	if t.DueDate != nil {
		due = t.DueDate.Format("2006-01-02 15:04")
	}
	comp := "未完成"
	if t.CompletedAt != nil {
		comp = t.CompletedAt.Format("2006-01-02 15:04")
	}
	lines = append(lines, metaStyle.Render(fmt.Sprintf(
		"优先级: %s | 状态: %s | 作用域: %s",
		priorityText(t.Priority), statusText(t.Status), scope)))
	lines = append(lines, metaStyle.Render(fmt.Sprintf(
		"截止时间: %s | 完成时间: %s", due, comp)))

	// === 分隔线 ===
	lines = append(lines, "")
	separatorLine := lipgloss.NewStyle().
		Foreground(theme.Border).
		Render(strings.Repeat("─", width))
	lines = append(lines, separatorLine)

	// === 区块 3：描述 ===
	if t.Description != "" {
		lines = append(lines, "")
		descLines := utils.RenderDetailSection("📝", "描述", t.Description, width)
		lines = append(lines, descLines...)
	}

	return strings.Join(lines, "\n")
}

func (p *ListPage) Meta() core.Meta {
	// 详情页模式
	if p.showing {
		return core.Meta{
			Title:      "待办详情",
			Breadcrumb: "待办管理 > 详情",
			Keys: []components.KeyHint{
				{Key: "↑/↓ j/k", Desc: "滚动"},
				{Key: "PgUp/PgDn", Desc: "翻页"},
				{Key: "Home/End", Desc: "首/尾"},
				{Key: "Esc", Desc: "返回列表"},
			},
		}
	}

	// 列表模式
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
			{Key: "f", Desc: "清空全部"},
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
			return deleteErrorMsg{err: err} // ✅ 返回错误消息
		}
		return deleteSuccessMsg{} // ✅ 返回成功消息
	}
}

// doFinalDelete 执行批量删除操作
func (p *ListPage) doFinalDelete() tea.Cmd {
	return func() tea.Msg {
		ctx := p.bs.Context()
		scope := p.scopeFilter.String()
		deletedCount, err := p.bs.ToDoService.DeleteAllByScope(ctx, scope, p.bs.CurrentScope)
		if err != nil {
			return finalDeleteErrorMsg{err: err}
		}
		return finalDeleteMsg{count: deletedCount}
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
