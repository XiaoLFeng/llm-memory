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
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// DetailModel 待办详情模型
// 嘿嘿~ 查看待办的详细内容！✅
type DetailModel struct {
	bs       *startup.Bootstrap
	id       uint
	todo     *entity.ToDo
	viewport viewport.Model
	frame    *components.Frame // 添加 Frame 支持
	ready    bool
	width    int
	height   int
	loading  bool
	err      error
}

// NewDetailModel 创建待办详情模型
func NewDetailModel(bs *startup.Bootstrap, id uint) *DetailModel {
	return &DetailModel{
		bs:      bs,
		id:      id,
		frame:   components.NewFrame(80, 24), // 初始化 Frame
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
		todo, err := m.bs.ToDoService.GetToDo(context.Background(), m.id)
		if err != nil {
			return todosErrorMsg{err}
		}
		return todoLoadedMsg{todo}
	}
}

type todoLoadedMsg struct {
	todo *entity.ToDo
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
			if m.todo != nil && m.todo.Status == entity.ToDoStatusPending {
				return m, m.startTodo()
			}

		case msg.String() == "f":
			// 完成待办
			if m.todo != nil && m.todo.Status == entity.ToDoStatusInProgress {
				return m, m.completeTodo()
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.frame.SetSize(msg.Width, msg.Height)

		// 使用 frame 的内容尺寸
		contentWidth := m.frame.GetContentWidth()
		contentHeight := m.frame.GetContentHeight()

		if !m.ready {
			m.viewport = viewport.New(contentWidth, contentHeight)
			m.viewport.YPosition = 0
			m.ready = true
		} else {
			m.viewport.Width = contentWidth
			m.viewport.Height = contentHeight
		}
		// 无论数据是否已加载，都尝试更新内容
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
		err := m.bs.ToDoService.StartToDo(context.Background(), m.id)
		if err != nil {
			return todosErrorMsg{err}
		}
		return todoStartedMsg{m.id}
	}
}

// completeTodo 完成待办
func (m *DetailModel) completeTodo() tea.Cmd {
	return func() tea.Msg {
		err := m.bs.ToDoService.CompleteToDo(context.Background(), m.id)
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

	// 直接使用 viewport 宽度，减去卡片边框和内边距
	cardWidth := m.viewport.Width - 4
	if cardWidth < 40 {
		cardWidth = 40
	}

	var sections []string

	// 基本信息卡片
	basicInfo := m.renderBasicInfo()
	sections = append(sections, components.NestedCard("基本信息", basicInfo, cardWidth))

	// 详细信息卡片
	detailInfo := m.renderDetailInfo()
	sections = append(sections, components.NestedCard("详细信息", detailInfo, cardWidth))

	// 时间信息卡片
	timeInfo := m.renderTimeInfo()
	sections = append(sections, components.NestedCard("时间信息", timeInfo, cardWidth))

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
	if len(m.todo.Tags) > 0 {
		// 转换 []entity.ToDoTag 为 []string
		tags := make([]string, len(m.todo.Tags))
		for i, t := range m.todo.Tags {
			tags[i] = t.Tag
		}
		tagsBadge := components.TagsBadge(tags)
		lines = append(lines, components.InfoRow("标签", tagsBadge))
	} else {
		lines = append(lines, components.InfoRow("标签", lipgloss.NewStyle().Foreground(styles.Overlay0).Render("-")))
	}

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
// 嘿嘿~ 现在使用统一的 Frame 渲染！
func (m *DetailModel) View() string {
	breadcrumb := "待办管理 > 待办详情"
	if m.todo != nil {
		breadcrumb = "待办管理 > " + m.todo.Title
	}

	// 加载中
	if m.loading {
		content := lipgloss.Place(
			m.frame.GetContentWidth(),
			m.frame.GetContentHeight(),
			lipgloss.Center,
			lipgloss.Center,
			components.CardInfo("", "加载中...", 40),
		)
		keys := []string{"esc 返回"}
		return m.frame.Render(breadcrumb, content, keys, "")
	}

	// 错误显示
	if m.err != nil {
		content := lipgloss.Place(
			m.frame.GetContentWidth(),
			m.frame.GetContentHeight(),
			lipgloss.Center,
			lipgloss.Center,
			components.CardError("错误", m.err.Error(), 60),
		)
		keys := []string{"esc 返回"}
		return m.frame.Render(breadcrumb, content, keys, "")
	}

	// 正常显示
	content := ""
	if m.ready {
		content = m.viewport.View()
	}

	keys := []string{
		"↑/↓ 滚动",
		"s 开始",
		"f 完成",
		"esc 返回",
	}

	return m.frame.Render(breadcrumb, content, keys, "")
}
