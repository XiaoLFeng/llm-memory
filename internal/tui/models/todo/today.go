package todo

import (
	"context"
	"strings"

	"github.com/XiaoLFeng/llm-memory/internal/models/entity"
	"github.com/XiaoLFeng/llm-memory/internal/tui/common"
	"github.com/XiaoLFeng/llm-memory/internal/tui/components"
	"github.com/XiaoLFeng/llm-memory/internal/tui/styles"
	"github.com/XiaoLFeng/llm-memory/internal/tui/utils"
	"github.com/XiaoLFeng/llm-memory/startup"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TodayModel 今日待办模型
// 嘿嘿~ 展示今日待办的列表！📅
type TodayModel struct {
	bs      *startup.Bootstrap
	list    list.Model
	todos   []entity.ToDo
	width   int
	height  int
	loading bool
	err     error
}

// NewTodayModel 创建今日待办模型
func NewTodayModel(bs *startup.Bootstrap) *TodayModel {
	// 创建列表
	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = styles.ListSelectedStyle
	delegate.Styles.SelectedDesc = styles.ListDescStyle
	delegate.Styles.NormalTitle = styles.ListItemStyle
	delegate.Styles.NormalDesc = styles.ListDescStyle

	l := list.New([]list.Item{}, delegate, 80, 20)
	l.Title = "📅 今日待办"
	l.SetShowHelp(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = styles.ListTitleStyle

	return &TodayModel{
		bs:      bs,
		list:    l,
		loading: true,
	}
}

// Title 返回页面标题
func (m *TodayModel) Title() string {
	return "今日待办"
}

// ShortHelp 返回快捷键帮助
func (m *TodayModel) ShortHelp() []key.Binding {
	return []key.Binding{
		common.KeyUp, common.KeyDown, common.KeyEnter, common.KeyBack,
	}
}

// Init 初始化
func (m *TodayModel) Init() tea.Cmd {
	return m.loadTodayTodos()
}

// loadTodayTodos 加载今日待办列表
func (m *TodayModel) loadTodayTodos() tea.Cmd {
	return func() tea.Msg {
		todos, err := m.bs.ToDoService.ListToday(context.Background())
		if err != nil {
			return todosErrorMsg{err}
		}
		return todayTodosLoadedMsg{todos}
	}
}

type todayTodosLoadedMsg struct {
	todos []entity.ToDo
}

// Update 处理输入
func (m *TodayModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, common.KeyBack):
			return m, common.Back()

		case key.Matches(msg, common.KeyEnter):
			if item, ok := m.list.SelectedItem().(todoItem); ok {
				return m, common.Navigate(common.PageTodoDetail, map[string]any{"id": item.todo.ID})
			}

		case msg.String() == "s":
			// 开始待办
			if item, ok := m.list.SelectedItem().(todoItem); ok {
				if item.todo.Status == entity.ToDoStatusPending {
					return m, m.startTodo(item.todo.ID)
				}
			}

		case msg.String() == "f":
			// 完成待办
			if item, ok := m.list.SelectedItem().(todoItem); ok {
				if item.todo.Status == entity.ToDoStatusInProgress {
					return m, m.completeTodo(item.todo.ID)
				}
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetSize(msg.Width-4, msg.Height-8)

	case todayTodosLoadedMsg:
		m.loading = false
		m.todos = msg.todos
		items := make([]list.Item, len(msg.todos))
		for i, todo := range msg.todos {
			items[i] = todoItem{todo: todo}
		}
		m.list.SetItems(items)

	case todosErrorMsg:
		m.loading = false
		m.err = msg.err

	case todoStartedMsg:
		cmds = append(cmds, m.loadTodayTodos())
		cmds = append(cmds, common.ShowToast("待办已开始", common.ToastSuccess))

	case todoCompletedMsg:
		cmds = append(cmds, m.loadTodayTodos())
		cmds = append(cmds, common.ShowToast("待办已完成", common.ToastSuccess))

	case common.RefreshMsg:
		m.loading = true
		cmds = append(cmds, m.loadTodayTodos())
	}

	// 更新列表
	newList, cmd := m.list.Update(msg)
	m.list = newList
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// startTodo 开始待办
func (m *TodayModel) startTodo(id uint) tea.Cmd {
	return func() tea.Msg {
		err := m.bs.ToDoService.StartToDo(context.Background(), id)
		if err != nil {
			return todosErrorMsg{err}
		}
		return todoStartedMsg{id}
	}
}

// completeTodo 完成待办
func (m *TodayModel) completeTodo(id uint) tea.Cmd {
	return func() tea.Msg {
		err := m.bs.ToDoService.CompleteToDo(context.Background(), id)
		if err != nil {
			return todosErrorMsg{err}
		}
		return todoCompletedMsg{id}
	}
}

// View 渲染界面
func (m *TodayModel) View() string {
	var content string

	if m.loading {
		content = styles.InfoStyle.Render("加载中...")
		if m.width > 0 && m.height > 0 {
			return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
		}
		return content
	}

	if m.err != nil {
		content = styles.ErrorStyle.Render("错误: " + m.err.Error())
		if m.width > 0 && m.height > 0 {
			return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
		}
		return content
	}

	if len(m.todos) == 0 {
		var b strings.Builder
		emptyContent := strings.Join([]string{
			styles.MutedStyle.Render("今日暂无待办事项~"),
			"",
			styles.SuccessStyle.Render("🎉"),
		}, "\n")

		cardContent := components.Card("📅 今日待办", emptyContent, m.width-4)
		b.WriteString(cardContent)
		b.WriteString("\n\n")

		keys := []string{
			styles.StatusKeyStyle.Render("esc") + " 返回",
		}
		b.WriteString(components.RenderKeysOnly(keys, m.width))

		content = b.String()
		if m.width > 0 && m.height > 0 {
			return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
		}
		return content
	}

	// 渲染列表内容
	var b strings.Builder
	listContent := m.renderList()
	cardContent := components.Card("📅 今日待办", listContent, m.width-4)
	b.WriteString(cardContent)
	b.WriteString("\n\n")

	// 底部快捷键状态栏
	keys := []string{
		styles.StatusKeyStyle.Render("↑/↓") + " 选择",
		styles.StatusKeyStyle.Render("enter") + " 查看",
		styles.StatusKeyStyle.Render("s") + " 开始",
		styles.StatusKeyStyle.Render("f") + " 完成",
		styles.StatusKeyStyle.Render("esc") + " 返回",
	}
	b.WriteString(components.RenderKeysOnly(keys, m.width))

	content = b.String()
	if m.width > 0 && m.height > 0 {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
	}
	return content
}

// renderList 渲染待办列表
func (m *TodayModel) renderList() string {
	var b strings.Builder
	selected := m.list.Index()

	for i, item := range m.list.Items() {
		if todoItem, ok := item.(todoItem); ok {
			line := m.renderTodoItem(todoItem.todo, i == selected)
			b.WriteString(line)
			if i < len(m.list.Items())-1 {
				b.WriteString("\n")
			}
		}
	}

	return b.String()
}

// renderTodoItem 渲染单个待办项
func (m *TodayModel) renderTodoItem(todo entity.ToDo, selected bool) string {
	// 指示器
	indicator := " "
	if selected {
		indicator = lipgloss.NewStyle().Foreground(styles.Primary).Render("▸")
	}

	// 状态图标
	statusIcon := components.StatusBadgeSimple(todo.Status.String())

	// 标题
	titleStyle := lipgloss.NewStyle().Foreground(styles.Text)
	if selected {
		titleStyle = titleStyle.Bold(true)
	}
	title := titleStyle.Render(todo.Title)

	// 优先级徽章
	priority := components.PriorityBadge(int(todo.Priority))

	// 截止时间
	dueDate := ""
	if todo.DueDate != nil {
		dueDate = components.TimeBadge(utils.FormatDate(*todo.DueDate))
	}

	// 组合行
	parts := []string{indicator, statusIcon, title, priority}
	if dueDate != "" {
		parts = append(parts, dueDate)
	}

	line := strings.Join(parts, " ")

	// 选中项带背景
	if selected {
		return lipgloss.NewStyle().
			Background(styles.Surface1).
			Foreground(styles.Text).
			Width(m.width-8).
			Padding(0, 1).
			Render(line)
	}

	return lipgloss.NewStyle().
		Width(m.width-8).
		Padding(0, 1).
		Render(line)
}

// 引入 utils 进行格式化
var _ = utils.FormatTime
