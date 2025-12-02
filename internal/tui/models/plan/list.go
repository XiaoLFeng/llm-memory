package plan

import (
	"context"
	"fmt"
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

// planItem 计划列表项
type planItem struct {
	plan types.Plan
}

func (i planItem) Title() string {
	return fmt.Sprintf("%d. %s %s", i.plan.ID, utils.FormatStatusIcon(string(i.plan.Status)), i.plan.Title)
}

func (i planItem) Description() string {
	return fmt.Sprintf("%s | %s", utils.FormatStatus(string(i.plan.Status)), utils.FormatProgress(i.plan.Progress, 10))
}

func (i planItem) FilterValue() string {
	return i.plan.Title
}

// ListModel 计划列表模型
// 嘿嘿~ 展示所有计划的列表！📋
type ListModel struct {
	bs      *startup.Bootstrap
	list    list.Model
	plans   []types.Plan
	width   int
	height  int
	loading bool
	err     error
}

// NewListModel 创建计划列表模型
func NewListModel(bs *startup.Bootstrap) *ListModel {
	// 创建列表
	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = styles.ListSelectedStyle
	delegate.Styles.SelectedDesc = styles.ListDescStyle
	delegate.Styles.NormalTitle = styles.ListItemStyle
	delegate.Styles.NormalDesc = styles.ListDescStyle

	l := list.New([]list.Item{}, delegate, 80, 20)
	l.Title = "📋 计划列表"
	l.SetShowHelp(false)
	l.SetFilteringEnabled(true)
	l.Styles.Title = styles.ListTitleStyle

	return &ListModel{
		bs:      bs,
		list:    l,
		loading: true,
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
	plans []types.Plan
}

type plansErrorMsg struct {
	err error
}

type planDeletedMsg struct {
	id int
}

type planStartedMsg struct {
	id int
}

type planCompletedMsg struct {
	id int
}

// Update 处理输入
func (m *ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// 如果正在过滤，让列表处理
		if m.list.FilterState() == list.Filtering {
			break
		}

		switch {
		case key.Matches(msg, common.KeyBack):
			return m, common.Back()

		case key.Matches(msg, common.KeyCreate):
			return m, common.Navigate(common.PagePlanCreate)

		case key.Matches(msg, common.KeyEnter):
			if item, ok := m.list.SelectedItem().(planItem); ok {
				return m, common.Navigate(common.PagePlanDetail, map[string]any{"id": item.plan.ID})
			}

		case key.Matches(msg, common.KeyDelete):
			if item, ok := m.list.SelectedItem().(planItem); ok {
				return m, common.ShowConfirm(
					"删除计划",
					fmt.Sprintf("确定要删除计划「%s」吗？", item.plan.Title),
					m.deletePlan(item.plan.ID),
					nil,
				)
			}

		case msg.String() == "s":
			// 开始计划
			if item, ok := m.list.SelectedItem().(planItem); ok {
				if item.plan.Status == types.PlanStatusPending {
					return m, m.startPlan(item.plan.ID)
				}
			}

		case msg.String() == "f":
			// 完成计划
			if item, ok := m.list.SelectedItem().(planItem); ok {
				if item.plan.Status == types.PlanStatusInProgress {
					return m, m.completePlan(item.plan.ID)
				}
			}

		case msg.String() == "p":
			// 更新进度
			if item, ok := m.list.SelectedItem().(planItem); ok {
				return m, common.Navigate(common.PagePlanProgress, map[string]any{
					"id":       item.plan.ID,
					"progress": item.plan.Progress,
				})
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetSize(msg.Width-4, msg.Height-8)

	case plansLoadedMsg:
		m.loading = false
		m.plans = msg.plans
		items := make([]list.Item, len(msg.plans))
		for i, plan := range msg.plans {
			items[i] = planItem{plan: plan}
		}
		m.list.SetItems(items)

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

	// 更新列表
	newList, cmd := m.list.Update(msg)
	m.list = newList
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// deletePlan 删除计划
func (m *ListModel) deletePlan(id int) tea.Cmd {
	return func() tea.Msg {
		err := m.bs.PlanService.DeletePlan(context.Background(), id)
		if err != nil {
			return plansErrorMsg{err}
		}
		return planDeletedMsg{id}
	}
}

// startPlan 开始计划
func (m *ListModel) startPlan(id int) tea.Cmd {
	return func() tea.Msg {
		err := m.bs.PlanService.StartPlan(context.Background(), id)
		if err != nil {
			return plansErrorMsg{err}
		}
		return planStartedMsg{id}
	}
}

// completePlan 完成计划
func (m *ListModel) completePlan(id int) tea.Cmd {
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
	var b strings.Builder

	if m.loading {
		b.WriteString(styles.InfoStyle.Render("加载中..."))
		return b.String()
	}

	if m.err != nil {
		b.WriteString(styles.ErrorStyle.Render("错误: " + m.err.Error()))
		return b.String()
	}

	if len(m.plans) == 0 {
		b.WriteString(styles.TitleStyle.Render("📋 计划列表"))
		b.WriteString("\n\n")
		b.WriteString(styles.MutedStyle.Render("暂无计划~ 按 c 创建新计划"))
		b.WriteString("\n\n")
		b.WriteString(styles.HelpStyle.Render("c 新建 | esc 返回"))
		return b.String()
	}

	b.WriteString(m.list.View())
	b.WriteString("\n")
	b.WriteString(styles.HelpStyle.Render("↑/↓ 选择 | enter 查看 | c 新建 | s 开始 | f 完成 | p 进度 | d 删除 | esc 返回"))

	return b.String()
}
