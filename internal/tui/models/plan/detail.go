package plan

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

// DetailModel 计划详情模型
// 嘿嘿~ 查看计划的详细内容！📋
type DetailModel struct {
	bs       *startup.Bootstrap
	id       int
	plan     *types.Plan
	viewport viewport.Model
	ready    bool
	width    int
	height   int
	loading  bool
	err      error
}

// NewDetailModel 创建计划详情模型
func NewDetailModel(bs *startup.Bootstrap, id int) *DetailModel {
	return &DetailModel{
		bs:      bs,
		id:      id,
		loading: true,
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
	plan *types.Plan
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
			if m.plan != nil && m.plan.Status == types.PlanStatusPending {
				return m, m.startPlan()
			}

		case msg.String() == "f":
			// 完成计划
			if m.plan != nil && m.plan.Status == types.PlanStatusInProgress {
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
		if !m.ready {
			m.viewport = viewport.New(msg.Width-4, msg.Height-10)
			m.viewport.YPosition = 0
			m.ready = true
		} else {
			m.viewport.Width = msg.Width - 4
			m.viewport.Height = msg.Height - 10
		}
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

	var b strings.Builder

	// 标题
	b.WriteString(styles.SubtitleStyle.Render("标题"))
	b.WriteString("\n")
	b.WriteString(m.plan.Title)
	b.WriteString("\n\n")

	// 状态
	b.WriteString(styles.SubtitleStyle.Render("状态"))
	b.WriteString("\n")
	b.WriteString(utils.FormatStatusIcon(string(m.plan.Status)) + " " + utils.FormatStatus(string(m.plan.Status)))
	b.WriteString("\n\n")

	// 进度
	b.WriteString(styles.SubtitleStyle.Render("进度"))
	b.WriteString("\n")
	b.WriteString(utils.FormatProgress(m.plan.Progress, 20))
	b.WriteString("\n\n")

	// 描述
	if m.plan.Description != "" {
		b.WriteString(styles.SubtitleStyle.Render("描述"))
		b.WriteString("\n")
		b.WriteString(m.plan.Description)
		b.WriteString("\n\n")
	}

	// 开始时间
	b.WriteString(styles.SubtitleStyle.Render("开始时间"))
	b.WriteString("\n")
	b.WriteString(utils.FormatTimePtr(m.plan.StartDate))
	b.WriteString("\n\n")

	// 结束时间
	b.WriteString(styles.SubtitleStyle.Render("结束时间"))
	b.WriteString("\n")
	b.WriteString(utils.FormatTimePtr(m.plan.EndDate))
	b.WriteString("\n\n")

	// 创建时间
	b.WriteString(styles.SubtitleStyle.Render("创建时间"))
	b.WriteString("\n")
	b.WriteString(utils.FormatTime(m.plan.CreatedAt))
	b.WriteString("\n\n")

	// 子任务
	if len(m.plan.SubTasks) > 0 {
		b.WriteString(styles.SubtitleStyle.Render("子任务"))
		b.WriteString("\n")
		for _, task := range m.plan.SubTasks {
			b.WriteString(utils.FormatStatusIcon(string(task.Status)) + " " + task.Title)
			b.WriteString("\n")
		}
	}

	return b.String()
}

// View 渲染界面
func (m *DetailModel) View() string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("📋 计划详情"))
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
	b.WriteString(styles.HelpStyle.Render("↑/↓ 滚动 | s 开始 | f 完成 | p 进度 | esc 返回"))

	return b.String()
}
