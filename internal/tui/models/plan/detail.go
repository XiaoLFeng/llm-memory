package plan

import (
	"context"
	"strings"

	"github.com/XiaoLFeng/llm-memory/internal/models/entity"
	"github.com/XiaoLFeng/llm-memory/internal/tui/common"
	"github.com/XiaoLFeng/llm-memory/internal/tui/components"
	"github.com/XiaoLFeng/llm-memory/internal/tui/utils"
	"github.com/XiaoLFeng/llm-memory/startup"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// DetailModel 计划详情模型
// 嘿嘿~ 查看计划的详细内容！📋
type DetailModel struct {
	bs       *startup.Bootstrap
	id       uint
	plan     *entity.Plan
	viewport viewport.Model
	ready    bool
	width    int
	height   int
	loading  bool
	err      error
	frame    *components.Frame
}

// NewDetailModel 创建计划详情模型
func NewDetailModel(bs *startup.Bootstrap, id uint) *DetailModel {
	return &DetailModel{
		bs:      bs,
		id:      id,
		loading: true,
		frame:   components.NewFrame(80, 24),
	}
}

// Title 返回页面标题
func (m *DetailModel) Title() string {
	if m.plan != nil {
		return m.plan.Title
	}
	return "计划详情"
}

// ShortHelp 返回快捷键帮助
func (m *DetailModel) ShortHelp() []key.Binding {
	return []key.Binding{common.KeyUp, common.KeyDown, common.KeyBack}
}

// Init 初始化
func (m *DetailModel) Init() tea.Cmd {
	return m.loadPlan()
}

// loadPlan 加载计划详情
func (m *DetailModel) loadPlan() tea.Cmd {
	return func() tea.Msg {
		plan, err := m.bs.PlanService.GetPlan(context.Background(), m.id)
		if err != nil {
			return plansErrorMsg{err}
		}
		return planLoadedMsg{plan}
	}
}

type planLoadedMsg struct {
	plan *entity.Plan
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
			// 开始计划
			if m.plan != nil && m.plan.Status == entity.PlanStatusPending {
				return m, m.startPlan()
			}

		case msg.String() == "f":
			// 完成计划
			if m.plan != nil && m.plan.Status == entity.PlanStatusInProgress {
				return m, m.completePlan()
			}

		case msg.String() == "p":
			// 更新进度
			if m.plan != nil {
				return m, common.Navigate(common.PagePlanProgress, map[string]any{
					"id":       m.plan.ID,
					"progress": m.plan.Progress,
				})
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.frame.SetSize(msg.Width, msg.Height)

		// 使用 frame 的内容尺寸，不再使用硬编码的减法
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
		if m.plan != nil {
			m.viewport.SetContent(m.renderContent())
		}

	case planLoadedMsg:
		m.loading = false
		m.plan = msg.plan
		if m.ready {
			m.viewport.SetContent(m.renderContent())
		}

	case planStartedMsg:
		m.loading = true
		cmds = append(cmds, m.loadPlan())
		cmds = append(cmds, common.ShowToast("计划已开始", common.ToastSuccess))

	case planCompletedMsg:
		m.loading = true
		cmds = append(cmds, m.loadPlan())
		cmds = append(cmds, common.ShowToast("计划已完成", common.ToastSuccess))

	case plansErrorMsg:
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

// startPlan 开始计划
func (m *DetailModel) startPlan() tea.Cmd {
	return func() tea.Msg {
		err := m.bs.PlanService.StartPlan(context.Background(), m.id)
		if err != nil {
			return plansErrorMsg{err}
		}
		return planStartedMsg{m.id}
	}
}

// completePlan 完成计划
func (m *DetailModel) completePlan() tea.Cmd {
	return func() tea.Msg {
		err := m.bs.PlanService.CompletePlan(context.Background(), m.id)
		if err != nil {
			return plansErrorMsg{err}
		}
		return planCompletedMsg{m.id}
	}
}

// renderContent 渲染内容
func (m *DetailModel) renderContent() string {
	if m.plan == nil {
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
	sections = append(sections, components.NestedCard("📝 基本信息", basicInfo, cardWidth))

	// 进度信息卡片
	progressInfo := m.renderProgressInfo()
	sections = append(sections, components.NestedCard("📊 进度信息", progressInfo, cardWidth))

	// 时间信息卡片
	timeInfo := m.renderTimeInfo()
	sections = append(sections, components.NestedCard("⏰ 时间信息", timeInfo, cardWidth))

	// 描述卡片（如果有）
	if m.plan.Description != "" {
		sections = append(sections, components.NestedCard("📄 描述", m.plan.Description, cardWidth))
	}

	// 子任务列表（如果有）
	if len(m.plan.SubTasks) > 0 {
		subTasksInfo := m.renderSubTasks()
		sections = append(sections, components.NestedCard("✓ 子任务", subTasksInfo, cardWidth))
	}

	return strings.Join(sections, "\n\n")
}

// renderBasicInfo 渲染基本信息
func (m *DetailModel) renderBasicInfo() string {
	var lines []string

	// 标题
	lines = append(lines, components.InfoRow("标题", m.plan.Title))

	// 状态
	status := components.StatusBadge(string(m.plan.Status))
	lines = append(lines, components.InfoRow("状态", status))

	// 作用域
	scope := components.ScopeBadgeFromGroupIDPath(m.plan.GroupID, m.plan.Path)
	lines = append(lines, components.InfoRow("作用域", scope))

	return strings.Join(lines, "\n")
}

// renderProgressInfo 渲染进度信息
func (m *DetailModel) renderProgressInfo() string {
	var lines []string

	// 进度徽章
	progressBadge := components.ProgressBadge(m.plan.Progress)
	lines = append(lines, components.InfoRow("进度", progressBadge))

	// 进度条
	progressBar := utils.FormatProgress(m.plan.Progress, 30)
	lines = append(lines, progressBar)

	return strings.Join(lines, "\n")
}

// renderTimeInfo 渲染时间信息
func (m *DetailModel) renderTimeInfo() string {
	var lines []string

	lines = append(lines, components.InfoRow("创建时间", utils.FormatTime(m.plan.CreatedAt)))
	lines = append(lines, components.InfoRow("开始时间", utils.FormatTimePtr(m.plan.StartDate)))
	lines = append(lines, components.InfoRow("结束时间", utils.FormatTimePtr(m.plan.EndDate)))

	return strings.Join(lines, "\n")
}

// renderSubTasks 渲染子任务列表
func (m *DetailModel) renderSubTasks() string {
	var lines []string

	for _, task := range m.plan.SubTasks {
		statusIcon := components.StatusBadgeSimple(string(task.Status))
		taskLine := statusIcon + " " + task.Title
		lines = append(lines, taskLine)
	}

	return strings.Join(lines, "\n")
}

// View 渲染界面
func (m *DetailModel) View() string {
	breadcrumb := "计划管理 > 计划详情"
	if m.plan != nil {
		breadcrumb = "计划管理 > " + m.plan.Title
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
		"p 进度",
		"esc 返回",
	}

	return m.frame.Render(breadcrumb, content, keys, "")
}
