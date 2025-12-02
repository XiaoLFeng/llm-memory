package todo

import (
	"context"
	"strings"

	"github.com/XiaoLFeng/llm-memory/internal/tui/common"
	"github.com/XiaoLFeng/llm-memory/internal/tui/styles"
	"github.com/XiaoLFeng/llm-memory/internal/tui/utils"
	"github.com/XiaoLFeng/llm-memory/pkg/types"
	"github.com/XiaoLFeng/llm-memory/startup"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// TodayModel 今日待办模型
// 嘿嘿~ 展示今日待办的列表！📅
type TodayModel struct {
	bs      *startup.Bootstrap
	list    list.Model
	todos   []types.Todo
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
		todos, err := m.bs.TodoService.ListToday(context.Background())
		if err != nil {
			return todosErrorMsg{err}
		}
		return todayTodosLoadedMsg{todos}
	}
}

type todayTodosLoadedMsg struct {
	todos []types.Todo
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
				if item.todo.Status == types.TodoStatusPending {
					return m, m.startTodo(item.todo.ID)
				}
			}

		case msg.String() == "f":
			// 完成待办
			if item, ok := m.list.SelectedItem().(todoItem); ok {
				if item.todo.Status == types.TodoStatusInProgress {
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
func (m *TodayModel) startTodo(id int) tea.Cmd {
	return func() tea.Msg {
		err := m.bs.TodoService.StartTodo(context.Background(), id)
		if err != nil {
			return todosErrorMsg{err}
		}
		return todoStartedMsg{id}
	}
}

// completeTodo 完成待办
func (m *TodayModel) completeTodo(id int) tea.Cmd {
	return func() tea.Msg {
		err := m.bs.TodoService.CompleteTodo(context.Background(), id)
		if err != nil {
			return todosErrorMsg{err}
		}
		return todoCompletedMsg{id}
	}
}

// View 渲染界面
func (m *TodayModel) View() string {
	var b strings.Builder

	if m.loading {
		b.WriteString(styles.InfoStyle.Render("加载中..."))
		return b.String()
	}

	if m.err != nil {
		b.WriteString(styles.ErrorStyle.Render("错误: " + m.err.Error()))
		return b.String()
	}

	if len(m.todos) == 0 {
		b.WriteString(styles.TitleStyle.Render("📅 今日待办"))
		b.WriteString("\n\n")
		b.WriteString(styles.MutedStyle.Render("今日暂无待办事项~ 🎉"))
		b.WriteString("\n\n")
		b.WriteString(styles.HelpStyle.Render("esc 返回"))
		return b.String()
	}

	b.WriteString(m.list.View())
	b.WriteString("\n")
	b.WriteString(styles.HelpStyle.Render("↑/↓ 选择 | enter 查看 | s 开始 | f 完成 | esc 返回"))

	return b.String()
}

// 引入 utils 进行格式化
var _ = utils.FormatTime
