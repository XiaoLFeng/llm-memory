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
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type loadMsg struct {
	items []planItem
	err   error
}

type planItem struct {
	ID          int64
	Code        string
	Title       string
	Description string // 计划摘要
	Content     string // 详细内容
	Status      string
	Progress    int
	PathID      int64
	CreatedAt   time.Time
	TodoCount   int        // 待办数量
	Todos       []todoItem // 关联的待办列表
}

// todoItem 计划详情中的待办项
type todoItem struct {
	ID       int64
	Code     string
	Title    string
	Status   entity.ToDoStatus
	Priority entity.ToDoPriority
}

type ListPage struct {
	bs              *startup.Bootstrap
	frame           *layout.Frame
	width           int
	height          int
	loading         bool
	err             error
	items           []planItem
	cursor          int
	showing         bool
	scopeFilter     utils.ScopeFilter // 作用域过滤状态
	detailViewport  viewport.Model    // 详情页滚动视图
	push            func(core.PageID) tea.Cmd
	pushWithData    func(core.PageID, interface{}) tea.Cmd
	confirmDelete   bool
	deleteTarget    int64
	deleteYesActive bool // true=选中确认，false=选中取消

	// Todo 交互相关
	todoMode          bool  // 是否处于 Todo 操作模式
	todoCursor        int   // Todo 列表游标
	todoConfirmDelete bool  // Todo 删除确认模式
	todoDeleteTarget  int64 // 要删除的 Todo ID
	todoYesActive     bool  // Todo 删除确认按钮状态
}

func NewListPage(bs *startup.Bootstrap, push func(core.PageID) tea.Cmd, pushWithData func(core.PageID, interface{}) tea.Cmd) *ListPage {
	// 初始化 viewport（初始尺寸，后续动态调整）
	vp := viewport.New(60, 10)
	vp.Style = lipgloss.NewStyle()

	return &ListPage{
		bs:              bs,
		frame:           layout.NewFrame(80, 24),
		width:           80,
		height:          24,
		loading:         true,
		detailViewport:  vp,
		push:            push,
		pushWithData:    pushWithData,
		deleteYesActive: true,
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
			// 获取关联的 Todos
			todos, _ := p.bs.ToDoService.ListToDosByPlanCode(ctx, pl.Code)
			todoItems := make([]todoItem, 0, len(todos))
			for _, t := range todos {
				todoItems = append(todoItems, todoItem{
					ID:       t.ID,
					Code:     t.Code,
					Title:    t.Title,
					Status:   t.Status,
					Priority: t.Priority,
				})
			}

			items = append(items, planItem{
				ID:          pl.ID,
				Code:        pl.Code,
				Title:       pl.Title,
				Description: pl.Description,
				Content:     pl.Content,
				Status:      string(pl.Status),
				Progress:    pl.Progress,
				PathID:      pl.PathID,
				CreatedAt:   pl.CreatedAt,
				TodoCount:   len(todos),
				Todos:       todoItems,
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
		// 删除确认模式
		if p.confirmDelete {
			switch v.String() {
			case "left", "h", "right", "l":
				p.deleteYesActive = !p.deleteYesActive
			case "y", "Y":
				return p, p.doDelete()
			case "n", "N", "esc":
				p.confirmDelete = false
				p.deleteTarget = 0
			case "enter":
				if p.deleteYesActive {
					return p, p.doDelete()
				} else {
					p.confirmDelete = false
					p.deleteTarget = 0
				}
			}
			return p, nil
		}

		// Todo 删除确认模式
		if p.todoConfirmDelete {
			switch v.String() {
			case "left", "h", "right", "l":
				p.todoYesActive = !p.todoYesActive
			case "y", "Y":
				return p, p.doDeleteTodo()
			case "n", "N", "esc":
				p.todoConfirmDelete = false
				p.todoDeleteTarget = 0
			case "enter":
				if p.todoYesActive {
					return p, p.doDeleteTodo()
				} else {
					p.todoConfirmDelete = false
					p.todoDeleteTarget = 0
				}
			}
			return p, nil
		}

		// 详情页模式
		if p.showing {
			// Todo 操作模式
			if p.todoMode {
				switch v.String() {
				case "tab":
					// 退出 Todo 模式
					p.todoMode = false
					return p, nil
				case "esc":
					// 退出 Todo 模式
					p.todoMode = false
					return p, nil
				case "up", "k":
					// 移动 Todo 游标
					if p.todoCursor > 0 {
						p.todoCursor--
					}
					return p, nil
				case "down", "j":
					// 移动 Todo 游标
					if len(p.items) > 0 && len(p.items[p.cursor].Todos) > 0 && p.todoCursor < len(p.items[p.cursor].Todos)-1 {
						p.todoCursor++
					}
					return p, nil
				case "n":
					// 创建新 Todo
					if p.pushWithData != nil && len(p.items) > 0 {
						return p, p.pushWithData(core.PageTodoCreate, &TodoCreateContext{
							PlanCode:  p.items[p.cursor].Code,
							PlanTitle: p.items[p.cursor].Title,
						})
					}
					return p, nil
				case "e":
					// 编辑选中的 Todo
					if p.pushWithData != nil && len(p.items) > 0 && len(p.items[p.cursor].Todos) > 0 {
						todoID := p.items[p.cursor].Todos[p.todoCursor].ID
						return p, p.pushWithData(core.PageTodoEdit, todoID)
					}
					return p, nil
				case "d":
					// 删除选中的 Todo
					if len(p.items) > 0 && len(p.items[p.cursor].Todos) > 0 {
						p.todoConfirmDelete = true
						p.todoDeleteTarget = p.items[p.cursor].Todos[p.todoCursor].ID
						p.todoYesActive = true
					}
					return p, nil
				case "s":
					// 开始 Todo
					return p, p.startTodo()
				case "c":
					// 完成 Todo
					return p, p.completeTodo()
				case "x":
					// 取消 Todo
					return p, p.cancelTodo()
				case "K":
					// 上移排序
					return p, p.moveTodoUp()
				case "J":
					// 下移排序
					return p, p.moveTodoDown()
				}
				return p, nil
			}

			// 详情页只读模式
			switch v.String() {
			case "esc", "q":
				p.showing = false
				p.todoMode = false
				p.todoCursor = 0
				return p, nil
			case "tab":
				// 进入 Todo 模式（如果有 Todo）
				if len(p.items) > 0 && len(p.items[p.cursor].Todos) > 0 {
					p.todoMode = true
					p.todoCursor = 0
				}
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
			case "n":
				// 在只读模式下也支持创建 Todo
				if p.pushWithData != nil && len(p.items) > 0 {
					return p, p.pushWithData(core.PageTodoCreate, &TodoCreateContext{
						PlanCode:  p.items[p.cursor].Code,
						PlanTitle: p.items[p.cursor].Title,
					})
				}
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
			if !p.showing && p.push != nil {
				return p, p.push(core.PagePlanCreate)
			}
		case "e":
			if !p.showing && len(p.items) > 0 && p.pushWithData != nil {
				return p, p.pushWithData(core.PagePlanEdit, p.items[p.cursor].ID)
			}
		case "d":
			if !p.showing && len(p.items) > 0 {
				p.confirmDelete = true
				p.deleteTarget = p.items[p.cursor].ID
				p.deleteYesActive = true
			}
		case "?":
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
			if p.cursor < 0 {
				p.cursor = 0
			}
		}
	case tea.WindowSizeMsg:
		// 动态调整 viewport 尺寸
		if p.showing {
			// 直接使用终端尺寸，减去详情视图自身的 header/footer
			const detailOverhead = 4
			p.detailViewport.Width = v.Width - 4
			p.detailViewport.Height = v.Height - detailOverhead
		}
	}
	return p, nil
}

func (p *ListPage) View() string {
	cw, _ := p.frame.ContentSize()
	cardW := layout.FitCardWidth(cw)
	scopeLabel := p.scopeFilter.Label()
	titleWithScope := fmt.Sprintf("%s 计划列表 [%s]", theme.IconPlan, scopeLabel)

	// 删除确认模式
	if p.confirmDelete {
		var itemName string
		for _, item := range p.items {
			if item.ID == p.deleteTarget {
				itemName = item.Title
				break
			}
		}
		return components.ConfirmDialogWithButtons("确认删除",
			fmt.Sprintf("确定要删除计划「%s」吗？\n此操作不可撤销。", itemName),
			cardW, p.deleteYesActive)
	}

	// Todo 删除确认模式
	if p.todoConfirmDelete {
		var todoTitle string
		if len(p.items) > 0 && len(p.items[p.cursor].Todos) > 0 {
			for _, todo := range p.items[p.cursor].Todos {
				if todo.ID == p.todoDeleteTarget {
					todoTitle = todo.Title
					break
				}
			}
		}
		return components.ConfirmDialogWithButtons("确认删除待办",
			fmt.Sprintf("确定要删除待办「%s」吗？\n此操作不可撤销。", todoTitle),
			cardW, p.todoYesActive)
	}

	switch {
	case p.loading:
		return components.LoadingState(titleWithScope, "加载计划中...", cardW)
	case p.err != nil:
		return components.ErrorState(titleWithScope, p.err.Error(), cardW)
	case len(p.items) == 0:
		return components.EmptyState(titleWithScope, "暂无计划，按 c 创建吧~", cardW)
	default:
		if p.showing {
			// === 使用 viewport 渲染详情页 ===
			// 直接使用终端尺寸，减去详情视图自身的 header/footer
			// title(1) + 空行(1) + 空行(1) + scrollHint(1) = 4行
			const detailOverhead = 4

			viewportWidth := p.width - 4
			viewportHeight := p.height - detailOverhead

			p.detailViewport.Width = viewportWidth
			p.detailViewport.Height = viewportHeight

			// 生成详情内容并设置到 viewport
			detailContent := p.renderDetail(p.detailViewport.Width)
			p.detailViewport.SetContent(detailContent)

			// 滚动进度指示器和模式提示
			scrollPercent := p.detailViewport.ScrollPercent() * 100
			scrollInfo := fmt.Sprintf("%.0f%%", scrollPercent)
			var scrollHint string
			if p.todoMode {
				scrollHint = theme.TextDim.Render(fmt.Sprintf(
					"[Todo 模式] %s | n新建 e编辑 d删除 s开始 c完成 x取消 J/K排序 | Tab/Esc 退出", scrollInfo))
			} else {
				scrollHint = theme.TextDim.Render(fmt.Sprintf(
					"滚动: %s | ↑/↓ j/k PgUp/PgDn Home/End | n新建Todo | Tab Todo模式 | Esc 返回", scrollInfo))
			}

			// 组合视图
			titleText := theme.IconPlan + " 计划详情"
			if p.todoMode {
				titleText = theme.IconTodo + " 计划详情 - Todo 模式"
			}
			title := theme.Title.Render(titleText)
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
		pl := p.items[i]
		scope := utils.ScopeTag(pl.PathID, p.bs)
		status := statusText(pl.Status, pl.Progress)
		todoCountStr := fmt.Sprintf("(%d)", pl.TodoCount)
		line := fmt.Sprintf("%s [%s] %s%s · %s · %d%%",
			scope, pl.Code, pl.Title, todoCountStr, status, pl.Progress)
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
	scope := utils.ScopeTag(pl.PathID, p.bs)

	var lines []string

	// === 区块 1：标题 ===
	titleLine := theme.FormLabel.Bold(true).Render("标题: ") + theme.TextMain.Render(pl.Title)
	lines = append(lines, titleLine)
	lines = append(lines, "")

	// === 区块 2：元数据 ===
	metaStyle := theme.TextDim
	lines = append(lines, metaStyle.Render(fmt.Sprintf(
		"标识码: %s | 状态: %s | 进度: %d%% | 作用域: %s",
		pl.Code, statusText(pl.Status, pl.Progress), pl.Progress, scope)))
	lines = append(lines, metaStyle.Render(fmt.Sprintf(
		"待办数量: %d | 创建时间: %s", pl.TodoCount, pl.CreatedAt.Format("2006-01-02 15:04:05"))))

	// === 分隔线 ===
	lines = append(lines, "")
	separatorLine := lipgloss.NewStyle().
		Foreground(theme.Border).
		Render(strings.Repeat("─", width))
	lines = append(lines, separatorLine)

	// === 区块 3：描述 ===
	if pl.Description != "" {
		lines = append(lines, "")
		descLines := utils.RenderDetailSection("📝", "描述", pl.Description, width)
		lines = append(lines, descLines...)
	}

	// === 区块 4：详细内容 ===
	if pl.Content != "" {
		lines = append(lines, "")
		lines = append(lines, "")
		contentLines := utils.RenderDetailSection("📄", "详细内容", pl.Content, width)
		lines = append(lines, contentLines...)
	}

	// === 区块 5：待办列表 ===
	if len(pl.Todos) > 0 {
		lines = append(lines, "")
		lines = append(lines, "")
		lines = append(lines, separatorLine)
		lines = append(lines, "")
		todoModeHint := ""
		if p.todoMode {
			todoModeHint = " (Tab 退出选择模式)"
		} else {
			todoModeHint = " (Tab 进入选择模式)"
		}
		todoHeader := theme.FormLabel.Bold(true).Render("📋 待办事项列表" + todoModeHint)
		lines = append(lines, todoHeader)
		lines = append(lines, "")
		for i, t := range pl.Todos {
			// 状态图标
			statusIcon := getStatusIcon(t.Status)
			// 优先级图标
			priorityIcon := getPriorityIcon(t.Priority)
			// 格式化行
			todoLine := fmt.Sprintf("  %s %s [%s] %s (%s, %s)",
				statusIcon, priorityIcon, t.Code, t.Title,
				todoStatusText(t.Status), todoPriorityText(t.Priority))

			// Todo 模式下高亮选中项
			if p.todoMode && i == p.todoCursor {
				todoLine = lipgloss.NewStyle().
					Foreground(theme.Primary).
					Bold(true).
					Render("▶" + todoLine[1:])
			}
			lines = append(lines, todoLine)
		}
	} else {
		// 无待办时显示提示
		lines = append(lines, "")
		lines = append(lines, "")
		lines = append(lines, separatorLine)
		lines = append(lines, "")
		todoHeader := theme.FormLabel.Bold(true).Render("📋 待办事项列表")
		lines = append(lines, todoHeader)
		lines = append(lines, "")
		lines = append(lines, theme.TextDim.Render("  暂无待办事项，按 n 创建新待办"))
	}

	return strings.Join(lines, "\n")
}

func (p *ListPage) Meta() core.Meta {
	// 详情页 + Todo 模式
	if p.showing && p.todoMode {
		return core.Meta{
			Title:      "计划详情 - Todo 模式",
			Breadcrumb: "计划管理 > 详情 > Todo",
			Keys: []components.KeyHint{
				{Key: "n", Desc: "新建 Todo"},
				{Key: "e", Desc: "编辑"},
				{Key: "d", Desc: "删除"},
				{Key: "s/c/x", Desc: "开始/完成/取消"},
				{Key: "J/K", Desc: "调整排序"},
				{Key: "↑/↓", Desc: "选择"},
				{Key: "Tab/Esc", Desc: "退出 Todo 模式"},
			},
		}
	}

	// 详情页只读模式
	if p.showing {
		return core.Meta{
			Title:      "计划详情",
			Breadcrumb: "计划管理 > 详情",
			Keys: []components.KeyHint{
				{Key: "↑/↓ j/k", Desc: "滚动"},
				{Key: "PgUp/PgDn", Desc: "翻页"},
				{Key: "n", Desc: "新建 Todo"},
				{Key: "Tab", Desc: "Todo 模式"},
				{Key: "Esc", Desc: "返回列表"},
			},
		}
	}

	// 列表模式
	return core.Meta{
		Title:      "计划列表",
		Breadcrumb: "计划管理 > 列表",
		Extra:      fmt.Sprintf("[%s] Tab切换 r刷新", p.scopeFilter.Label()),
		Keys: []components.KeyHint{
			{Key: "Tab", Desc: "切换作用域"},
			{Key: "Enter", Desc: "详情"},
			{Key: "c", Desc: "新建计划"},
			{Key: "e", Desc: "编辑"},
			{Key: "d", Desc: "删除"},
			{Key: "r", Desc: "刷新"},
			{Key: "?", Desc: "帮助"},
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

// todoStatusText 将待办状态转换为中文显示
func todoStatusText(status entity.ToDoStatus) string {
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

// todoPriorityText 将待办优先级转换为中文显示
func todoPriorityText(priority entity.ToDoPriority) string {
	switch priority {
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

// doDelete 执行删除操作
func (p *ListPage) doDelete() tea.Cmd {
	return func() tea.Msg {
		ctx := p.bs.Context()
		err := p.bs.PlanService.DeletePlanByID(ctx, p.deleteTarget)
		p.confirmDelete = false
		p.deleteTarget = 0
		if err != nil {
			return loadMsg{err: err}
		}
		// 重新加载列表
		scopeStr := p.scopeFilter.String()
		plans, err := p.bs.PlanService.ListPlansByScope(ctx, scopeStr, p.bs.CurrentScope)
		if err != nil {
			return loadMsg{err: err}
		}
		items := make([]planItem, 0, len(plans))
		for _, pl := range plans {
			// 获取关联的 Todos
			todos, _ := p.bs.ToDoService.ListToDosByPlanCode(ctx, pl.Code)
			todoItems := make([]todoItem, 0, len(todos))
			for _, t := range todos {
				todoItems = append(todoItems, todoItem{
					ID:       t.ID,
					Code:     t.Code,
					Title:    t.Title,
					Status:   t.Status,
					Priority: t.Priority,
				})
			}

			items = append(items, planItem{
				ID:          pl.ID,
				Code:        pl.Code,
				Title:       pl.Title,
				Description: pl.Description,
				Content:     pl.Content,
				Status:      string(pl.Status),
				Progress:    pl.Progress,
				PathID:      pl.PathID,
				CreatedAt:   pl.CreatedAt,
				TodoCount:   len(todos),
				Todos:       todoItems,
			})
		}
		return loadMsg{items: items}
	}
}

// TodoCreateContext 从 Plan 详情页传递到 Todo 创建页的上下文
type TodoCreateContext struct {
	PlanCode  string
	PlanTitle string
}

// getStatusIcon 获取状态图标
func getStatusIcon(status entity.ToDoStatus) string {
	switch status {
	case entity.ToDoStatusCompleted:
		return "✅"
	case entity.ToDoStatusInProgress:
		return "🔄"
	case entity.ToDoStatusCancelled:
		return "❌"
	default:
		return "⬜"
	}
}

// getPriorityIcon 获取优先级图标
func getPriorityIcon(priority entity.ToDoPriority) string {
	switch priority {
	case entity.ToDoPriorityUrgent:
		return "🔴"
	case entity.ToDoPriorityHigh:
		return "🟠"
	case entity.ToDoPriorityMedium:
		return "🟡"
	default:
		return "🟢"
	}
}

// doDeleteTodo 执行删除 Todo 操作
func (p *ListPage) doDeleteTodo() tea.Cmd {
	return func() tea.Msg {
		ctx := p.bs.Context()
		err := p.bs.ToDoService.DeleteToDoByID(ctx, p.todoDeleteTarget)
		p.todoConfirmDelete = false
		p.todoDeleteTarget = 0
		if err != nil {
			return loadMsg{err: err}
		}
		// 重新加载
		return p.load()()
	}
}

// startTodo 开始选中的 Todo
func (p *ListPage) startTodo() tea.Cmd {
	return func() tea.Msg {
		if len(p.items) == 0 || len(p.items[p.cursor].Todos) == 0 {
			return nil
		}
		code := p.items[p.cursor].Todos[p.todoCursor].Code
		ctx := p.bs.Context()
		if err := p.bs.ToDoService.StartToDo(ctx, code); err != nil {
			return loadMsg{err: err}
		}
		return p.load()()
	}
}

// completeTodo 完成选中的 Todo
func (p *ListPage) completeTodo() tea.Cmd {
	return func() tea.Msg {
		if len(p.items) == 0 || len(p.items[p.cursor].Todos) == 0 {
			return nil
		}
		code := p.items[p.cursor].Todos[p.todoCursor].Code
		ctx := p.bs.Context()
		if err := p.bs.ToDoService.CompleteToDo(ctx, code); err != nil {
			return loadMsg{err: err}
		}
		return p.load()()
	}
}

// cancelTodo 取消选中的 Todo
func (p *ListPage) cancelTodo() tea.Cmd {
	return func() tea.Msg {
		if len(p.items) == 0 || len(p.items[p.cursor].Todos) == 0 {
			return nil
		}
		code := p.items[p.cursor].Todos[p.todoCursor].Code
		ctx := p.bs.Context()
		if err := p.bs.ToDoService.CancelToDo(ctx, code); err != nil {
			return loadMsg{err: err}
		}
		return p.load()()
	}
}

// moveTodoUp 上移 Todo 排序
func (p *ListPage) moveTodoUp() tea.Cmd {
	return func() tea.Msg {
		if len(p.items) == 0 || len(p.items[p.cursor].Todos) < 2 {
			return nil
		}
		if p.todoCursor <= 0 {
			return nil
		}
		ctx := p.bs.Context()
		currentTodo := p.items[p.cursor].Todos[p.todoCursor]
		prevTodo := p.items[p.cursor].Todos[p.todoCursor-1]
		if err := p.bs.ToDoService.SwapTodoOrder(ctx, currentTodo.ID, prevTodo.ID); err != nil {
			return loadMsg{err: err}
		}
		p.todoCursor--
		return p.load()()
	}
}

// moveTodoDown 下移 Todo 排序
func (p *ListPage) moveTodoDown() tea.Cmd {
	return func() tea.Msg {
		if len(p.items) == 0 || len(p.items[p.cursor].Todos) < 2 {
			return nil
		}
		if p.todoCursor >= len(p.items[p.cursor].Todos)-1 {
			return nil
		}
		ctx := p.bs.Context()
		currentTodo := p.items[p.cursor].Todos[p.todoCursor]
		nextTodo := p.items[p.cursor].Todos[p.todoCursor+1]
		if err := p.bs.ToDoService.SwapTodoOrder(ctx, currentTodo.ID, nextTodo.ID); err != nil {
			return loadMsg{err: err}
		}
		p.todoCursor++
		return p.load()()
	}
}
