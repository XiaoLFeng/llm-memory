package plan

import (
	"context"
	"fmt"
	"strings"

	"github.com/XiaoLFeng/llm-memory/internal/models/entity"
	"github.com/XiaoLFeng/llm-memory/internal/tui/common"
	"github.com/XiaoLFeng/llm-memory/internal/tui/components"
	"github.com/XiaoLFeng/llm-memory/internal/tui/styles"
	"github.com/XiaoLFeng/llm-memory/internal/tui/utils"
	"github.com/XiaoLFeng/llm-memory/startup"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ListModel 计划列表模型
// 嘿嘿~ 展示所有计划的列表！📋
type ListModel struct {
	bs          *startup.Bootstrap
	plans       []entity.Plan
	cursor      int
	width       int
	height      int
	loading     bool
	err         error
	frame       *components.Frame
	scrollStart int
}

// NewListModel 创建计划列表模型
func NewListModel(bs *startup.Bootstrap) *ListModel {
	return &ListModel{
		bs:      bs,
		loading: true,
		frame:   components.NewFrame(80, 24),
	}
}

// Title 返回页面标题
func (m *ListModel) Title() string {
	return "计划列表"
}

// ShortHelp 返回快捷键帮助
func (m *ListModel) ShortHelp() []key.Binding {
	return []key.Binding{
		common.KeyUp, common.KeyDown, common.KeyEnter,
		common.KeyCreate, common.KeyDelete, common.KeyBack,
	}
}

// Init 初始化
func (m *ListModel) Init() tea.Cmd {
	return m.loadPlans()
}

// loadPlans 加载计划列表
func (m *ListModel) loadPlans() tea.Cmd {
	return func() tea.Msg {
		plans, err := m.bs.PlanService.ListPlans(context.Background())
		if err != nil {
			return plansErrorMsg{err}
		}
		return plansLoadedMsg{plans}
	}
}

type plansLoadedMsg struct {
	plans []entity.Plan
}

type plansErrorMsg struct {
	err error
}

type planDeletedMsg struct {
	id uint
}

type planStartedMsg struct {
	id uint
}

type planCompletedMsg struct {
	id uint
}

// Update 处理输入
func (m *ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, common.KeyBack):
			return m, common.Back()

		case key.Matches(msg, common.KeyCreate):
			return m, common.Navigate(common.PagePlanCreate)

		case key.Matches(msg, common.KeyUp):
			if m.cursor > 0 {
				m.cursor--
				m.updateScroll()
			}

		case key.Matches(msg, common.KeyDown):
			if m.cursor < len(m.plans)-1 {
				m.cursor++
				m.updateScroll()
			}

		case key.Matches(msg, common.KeyEnter):
			if len(m.plans) > 0 && m.cursor < len(m.plans) {
				return m, common.Navigate(common.PagePlanDetail, map[string]any{"id": m.plans[m.cursor].ID})
			}

		case key.Matches(msg, common.KeyDelete):
			if len(m.plans) > 0 && m.cursor < len(m.plans) {
				plan := m.plans[m.cursor]
				return m, common.ShowConfirm(
					"删除计划",
					fmt.Sprintf("确定要删除计划「%s」吗？", plan.Title),
					m.deletePlan(plan.ID),
					nil,
				)
			}

		case msg.String() == "s":
			// 开始计划
			if len(m.plans) > 0 && m.cursor < len(m.plans) {
				plan := m.plans[m.cursor]
				if plan.Status == entity.PlanStatusPending {
					return m, m.startPlan(plan.ID)
				}
			}

		case msg.String() == "f":
			// 完成计划
			if len(m.plans) > 0 && m.cursor < len(m.plans) {
				plan := m.plans[m.cursor]
				if plan.Status == entity.PlanStatusInProgress {
					return m, m.completePlan(plan.ID)
				}
			}

		case msg.String() == "p":
			// 更新进度
			if len(m.plans) > 0 && m.cursor < len(m.plans) {
				plan := m.plans[m.cursor]
				return m, common.Navigate(common.PagePlanProgress, map[string]any{
					"id":       plan.ID,
					"progress": plan.Progress,
				})
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.frame.SetSize(msg.Width, msg.Height)

	case plansLoadedMsg:
		m.loading = false
		m.plans = msg.plans
		// 确保光标不越界
		if m.cursor >= len(m.plans) {
			m.cursor = len(m.plans) - 1
		}
		if m.cursor < 0 {
			m.cursor = 0
		}

	case plansErrorMsg:
		m.loading = false
		m.err = msg.err

	case planDeletedMsg:
		cmds = append(cmds, m.loadPlans())
		cmds = append(cmds, common.ShowToast("计划已删除", common.ToastSuccess))

	case planStartedMsg:
		cmds = append(cmds, m.loadPlans())
		cmds = append(cmds, common.ShowToast("计划已开始", common.ToastSuccess))

	case planCompletedMsg:
		cmds = append(cmds, m.loadPlans())
		cmds = append(cmds, common.ShowToast("计划已完成", common.ToastSuccess))

	case common.RefreshMsg:
		m.loading = true
		cmds = append(cmds, m.loadPlans())
	}

	return m, tea.Batch(cmds...)
}

// updateScroll 更新滚动位置
func (m *ListModel) updateScroll() {
	visibleLines := m.frame.GetContentHeight() / 3 // 每个条目大约占 3 行
	if m.cursor < m.scrollStart {
		m.scrollStart = m.cursor
	}
	if m.cursor >= m.scrollStart+visibleLines {
		m.scrollStart = m.cursor - visibleLines + 1
	}
}

// deletePlan 删除计划
func (m *ListModel) deletePlan(id uint) tea.Cmd {
	return func() tea.Msg {
		err := m.bs.PlanService.DeletePlan(context.Background(), id)
		if err != nil {
			return plansErrorMsg{err}
		}
		return planDeletedMsg{id}
	}
}

// startPlan 开始计划
func (m *ListModel) startPlan(id uint) tea.Cmd {
	return func() tea.Msg {
		err := m.bs.PlanService.StartPlan(context.Background(), id)
		if err != nil {
			return plansErrorMsg{err}
		}
		return planStartedMsg{id}
	}
}

// completePlan 完成计划
func (m *ListModel) completePlan(id uint) tea.Cmd {
	return func() tea.Msg {
		err := m.bs.PlanService.CompletePlan(context.Background(), id)
		if err != nil {
			return plansErrorMsg{err}
		}
		return planCompletedMsg{id}
	}
}

// View 渲染界面
func (m *ListModel) View() string {
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
		return m.frame.Render("计划管理 > 计划列表", content, keys, "")
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
		return m.frame.Render("计划管理 > 计划列表", content, keys, "")
	}

	// 空列表
	if len(m.plans) == 0 {
		emptyText := lipgloss.NewStyle().
			Foreground(styles.Subtext0).
			Render("暂无计划~ 按 c 创建新计划吧！")
		content := lipgloss.Place(
			m.frame.GetContentWidth(),
			m.frame.GetContentHeight(),
			lipgloss.Center,
			lipgloss.Center,
			components.Card("📋 计划列表", emptyText, 60),
		)
		keys := []string{"c 新建", "esc 返回"}
		return m.frame.Render("计划管理 > 计划列表", content, keys, "")
	}

	// 渲染列表项
	var listItems []string
	visibleLines := m.frame.GetContentHeight() / 3
	endIdx := m.scrollStart + visibleLines
	if endIdx > len(m.plans) {
		endIdx = len(m.plans)
	}

	for i := m.scrollStart; i < endIdx; i++ {
		plan := m.plans[i]
		listItems = append(listItems, m.renderPlanItem(plan, i == m.cursor))
	}

	listContent := strings.Join(listItems, "\n")

	// 统计信息
	extra := fmt.Sprintf("共 %d 个计划", len(m.plans))

	// 包装在卡片中
	cardContent := components.Card("📋 计划列表", listContent, m.frame.GetContentWidth()-4)

	content := lipgloss.NewStyle().
		Width(m.frame.GetContentWidth()).
		Render(cardContent)

	keys := []string{
		"↑/↓ 选择",
		"enter 查看",
		"c 新建",
		"s 开始",
		"f 完成",
		"p 进度",
		"d 删除",
		"esc 返回",
	}

	return m.frame.Render("计划管理 > 计划列表", content, keys, extra)
}

// renderPlanItem 渲染计划列表项
func (m *ListModel) renderPlanItem(plan entity.Plan, selected bool) string {
	// 指示器
	indicator := "  "
	if selected {
		indicator = lipgloss.NewStyle().
			Foreground(styles.Primary).
			Bold(true).
			Render("▸ ")
	}

	// 标题 + 作用域
	titleStyle := lipgloss.NewStyle().
		Foreground(styles.Text).
		Bold(selected)
	if selected {
		titleStyle = titleStyle.Foreground(styles.Primary)
	}

	title := titleStyle.Render(plan.Title)
	scope := components.ScopeBadgeFromGroupIDPath(plan.GroupID, plan.Path)

	// 状态 + 进度
	status := components.StatusBadge(string(plan.Status))
	progress := components.ProgressBadge(plan.Progress)

	// 第一行：指示器 + 标题 + 作用域
	line1 := indicator + title + " " + scope

	// 第二行：状态 + 进度条
	progressBar := utils.FormatProgress(plan.Progress, 20)
	line2 := "   " + status + " " + progress + " " + progressBar

	return line1 + "\n" + line2
}
