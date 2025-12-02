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
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

// DetailModel 待办详情模型
// 嘿嘿~ 查看待办的详细内容！✅
type DetailModel struct {
	bs       *startup.Bootstrap
	id       int
	todo     *types.Todo
	viewport viewport.Model
	ready    bool
	width    int
	height   int
	loading  bool
	err      error
}

// NewDetailModel 创建待办详情模型
func NewDetailModel(bs *startup.Bootstrap, id int) *DetailModel {
	return &DetailModel{
		bs:      bs,
		id:      id,
		loading: true,
	}
}

// Title 返回页面标题
func (m *DetailModel) Title() string {
	if m.todo != nil {
		return m.todo.Title
	}
	return "待办详情"
}

// ShortHelp 返回快捷键帮助
func (m *DetailModel) ShortHelp() []key.Binding {
	return []key.Binding{common.KeyUp, common.KeyDown, common.KeyBack}
}

// Init 初始化
func (m *DetailModel) Init() tea.Cmd {
	return m.loadTodo()
}

// loadTodo 加载待办详情
func (m *DetailModel) loadTodo() tea.Cmd {
	return func() tea.Msg {
		todo, err := m.bs.TodoService.GetTodo(context.Background(), m.id)
		if err != nil {
			return todosErrorMsg{err}
		}
		return todoLoadedMsg{todo}
	}
}

type todoLoadedMsg struct {
	todo *types.Todo
}

// Update 处理输入
func (m *DetailModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, common.KeyBack):
			return m, common.Back()

		case msg.String() == "s":
			// 开始待办
			if m.todo != nil && m.todo.Status == types.TodoStatusPending {
				return m, m.startTodo()
			}

		case msg.String() == "f":
			// 完成待办
			if m.todo != nil && m.todo.Status == types.TodoStatusInProgress {
				return m, m.completeTodo()
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if !m.ready {
			m.viewport = viewport.New(msg.Width-4, msg.Height-10)
			m.viewport.YPosition = 0
			m.ready = true
		} else {
			m.viewport.Width = msg.Width - 4
			m.viewport.Height = msg.Height - 10
		}
		if m.todo != nil {
			m.viewport.SetContent(m.renderContent())
		}

	case todoLoadedMsg:
		m.loading = false
		m.todo = msg.todo
		if m.ready {
			m.viewport.SetContent(m.renderContent())
		}

	case todoStartedMsg:
		m.loading = true
		cmds = append(cmds, m.loadTodo())
		cmds = append(cmds, common.ShowToast("待办已开始", common.ToastSuccess))

	case todoCompletedMsg:
		m.loading = true
		cmds = append(cmds, m.loadTodo())
		cmds = append(cmds, common.ShowToast("待办已完成", common.ToastSuccess))

	case todosErrorMsg:
		m.loading = false
		m.err = msg.err
	}

	// 更新 viewport
	if m.ready {
		newViewport, cmd := m.viewport.Update(msg)
		m.viewport = newViewport
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// startTodo 开始待办
func (m *DetailModel) startTodo() tea.Cmd {
	return func() tea.Msg {
		err := m.bs.TodoService.StartTodo(context.Background(), m.id)
		if err != nil {
			return todosErrorMsg{err}
		}
		return todoStartedMsg{m.id}
	}
}

// completeTodo 完成待办
func (m *DetailModel) completeTodo() tea.Cmd {
	return func() tea.Msg {
		err := m.bs.TodoService.CompleteTodo(context.Background(), m.id)
		if err != nil {
			return todosErrorMsg{err}
		}
		return todoCompletedMsg{m.id}
	}
}

// renderContent 渲染内容
func (m *DetailModel) renderContent() string {
	if m.todo == nil {
		return ""
	}

	var b strings.Builder

	// 标题
	b.WriteString(styles.SubtitleStyle.Render("标题"))
	b.WriteString("\n")
	b.WriteString(m.todo.Title)
	b.WriteString("\n\n")

	// 状态
	b.WriteString(styles.SubtitleStyle.Render("状态"))
	b.WriteString("\n")
	b.WriteString(utils.FormatTodoStatusIcon(int(m.todo.Status)) + " " + utils.FormatTodoStatus(int(m.todo.Status)))
	b.WriteString("\n\n")

	// 优先级
	b.WriteString(styles.SubtitleStyle.Render("优先级"))
	b.WriteString("\n")
	b.WriteString(utils.FormatPriorityIcon(int(m.todo.Priority)) + " " + utils.FormatPriority(int(m.todo.Priority)))
	b.WriteString("\n\n")

	// 描述
	if m.todo.Description != "" {
		b.WriteString(styles.SubtitleStyle.Render("描述"))
		b.WriteString("\n")
		b.WriteString(m.todo.Description)
		b.WriteString("\n\n")
	}

	// 截止日期
	b.WriteString(styles.SubtitleStyle.Render("截止日期"))
	b.WriteString("\n")
	b.WriteString(utils.FormatDatePtr(m.todo.DueDate))
	b.WriteString("\n\n")

	// 标签
	b.WriteString(styles.SubtitleStyle.Render("标签"))
	b.WriteString("\n")
	b.WriteString(utils.JoinTags(m.todo.Tags))
	b.WriteString("\n\n")

	// 创建时间
	b.WriteString(styles.SubtitleStyle.Render("创建时间"))
	b.WriteString("\n")
	b.WriteString(utils.FormatTime(m.todo.CreatedAt))
	b.WriteString("\n\n")

	// 完成时间
	if m.todo.CompletedAt != nil {
		b.WriteString(styles.SubtitleStyle.Render("完成时间"))
		b.WriteString("\n")
		b.WriteString(utils.FormatTime(*m.todo.CompletedAt))
	}

	return b.String()
}

// View 渲染界面
func (m *DetailModel) View() string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("✅ 待办详情"))
	b.WriteString("\n\n")

	if m.loading {
		b.WriteString(styles.InfoStyle.Render("加载中..."))
		return b.String()
	}

	if m.err != nil {
		b.WriteString(styles.ErrorStyle.Render("错误: " + m.err.Error()))
		return b.String()
	}

	if m.ready {
		b.WriteString(m.viewport.View())
	}

	b.WriteString("\n\n")
	b.WriteString(styles.HelpStyle.Render("↑/↓ 滚动 | s 开始 | f 完成 | esc 返回"))

	return b.String()
}
