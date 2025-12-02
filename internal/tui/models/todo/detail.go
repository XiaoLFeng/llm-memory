package todo

import (
	"context"
	"strings"

	"github.com/XiaoLFeng/llm-memory/internal/tui/common"
	"github.com/XiaoLFeng/llm-memory/internal/tui/components"
	"github.com/XiaoLFeng/llm-memory/internal/tui/styles"
	"github.com/XiaoLFeng/llm-memory/internal/tui/utils"
	"github.com/XiaoLFeng/llm-memory/pkg/types"
	"github.com/XiaoLFeng/llm-memory/startup"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

	var sections []string

	// 基本信息卡片
	basicInfo := m.renderBasicInfo()
	sections = append(sections, components.NestedCard("基本信息", basicInfo, m.width-12))

	// 详细信息卡片
	detailInfo := m.renderDetailInfo()
	sections = append(sections, components.NestedCard("详细信息", detailInfo, m.width-12))

	// 时间信息卡片
	timeInfo := m.renderTimeInfo()
	sections = append(sections, components.NestedCard("时间信息", timeInfo, m.width-12))

	return strings.Join(sections, "\n\n")
}

// renderBasicInfo 渲染基本信息
func (m *DetailModel) renderBasicInfo() string {
	var lines []string

	// 标题
	lines = append(lines, components.InfoRow("标题", m.todo.Title))

	// 状态
	statusBadge := components.StatusBadge(m.todo.Status.String())
	lines = append(lines, components.InfoRow("状态", statusBadge))

	// 优先级
	priorityBadge := components.PriorityBadge(int(m.todo.Priority))
	lines = append(lines, components.InfoRow("优先级", priorityBadge))

	// 作用域
	scopeBadge := components.ScopeBadgeFromGroupIDPath(m.todo.GroupID, m.todo.Path)
	lines = append(lines, components.InfoRow("作用域", scopeBadge))

	return strings.Join(lines, "\n")
}

// renderDetailInfo 渲染详细信息
func (m *DetailModel) renderDetailInfo() string {
	var lines []string

	// 描述
	description := m.todo.Description
	if description == "" {
		description = lipgloss.NewStyle().Foreground(styles.Overlay0).Render("-")
	}
	lines = append(lines, components.InfoRow("描述", description))

	// 标签
	tags := utils.JoinTags(m.todo.Tags)
	if len(m.todo.Tags) > 0 {
		tags = components.TagsBadge(m.todo.Tags)
	}
	lines = append(lines, components.InfoRow("标签", tags))

	return strings.Join(lines, "\n")
}

// renderTimeInfo 渲染时间信息
func (m *DetailModel) renderTimeInfo() string {
	var lines []string

	// 截止日期
	dueDate := utils.FormatDatePtr(m.todo.DueDate)
	if m.todo.DueDate != nil {
		dueDate = components.TimeBadge(dueDate)
	}
	lines = append(lines, components.InfoRow("截止日期", dueDate))

	// 创建时间
	createdAt := components.TimeBadge(utils.FormatTime(m.todo.CreatedAt))
	lines = append(lines, components.InfoRow("创建时间", createdAt))

	// 完成时间
	if m.todo.CompletedAt != nil {
		completedAt := components.TimeBadge(utils.FormatTime(*m.todo.CompletedAt))
		lines = append(lines, components.InfoRow("完成时间", completedAt))
	}

	return strings.Join(lines, "\n")
}

// View 渲染界面
func (m *DetailModel) View() string {
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

	if !m.ready {
		return ""
	}

	var b strings.Builder

	// 使用卡片包装内容
	viewportContent := m.viewport.View()
	cardContent := components.Card("✅ 待办详情", viewportContent, m.width-4)
	b.WriteString(cardContent)
	b.WriteString("\n\n")

	// 底部快捷键状态栏
	keys := []string{
		styles.StatusKeyStyle.Render("↑/↓") + " 滚动",
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
