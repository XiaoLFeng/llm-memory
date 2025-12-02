package models

import (
	"github.com/XiaoLFeng/llm-memory/internal/tui/common"
	"github.com/XiaoLFeng/llm-memory/internal/tui/components"
	"github.com/XiaoLFeng/llm-memory/internal/tui/models/group"
	"github.com/XiaoLFeng/llm-memory/internal/tui/models/memory"
	"github.com/XiaoLFeng/llm-memory/internal/tui/models/plan"
	"github.com/XiaoLFeng/llm-memory/internal/tui/models/todo"
	"github.com/XiaoLFeng/llm-memory/internal/tui/styles"
	"github.com/XiaoLFeng/llm-memory/startup"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// AppModel 根应用模型
// 嘿嘿~ 这是整个 TUI 的根模型，管理页面栈和全局状态！💖
type AppModel struct {
	bs          *startup.Bootstrap
	pageStack   []common.Page       // 页面栈
	currentPage common.Page         // 当前页面
	toast       *components.Toast   // 提示消息
	confirm     *components.Confirm // 确认对话框
	width       int
	height      int
	quitting    bool
}

// NewAppModel 创建根应用模型
func NewAppModel(bs *startup.Bootstrap) *AppModel {
	menu := NewMenuModel(bs)
	return &AppModel{
		bs:          bs,
		pageStack:   []common.Page{},
		currentPage: menu,
		toast:       components.NewToast(),
		confirm:     components.NewConfirm(),
		width:       80,
		height:      24,
	}
}

// Init 初始化
func (m *AppModel) Init() tea.Cmd {
	return m.currentPage.Init()
}

// Update 处理输入
func (m *AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// 如果确认对话框正在显示，优先处理
		if m.confirm.IsVisible() {
			newConfirm, cmd := m.confirm.Update(msg)
			m.confirm = newConfirm.(*components.Confirm)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
			return m, tea.Batch(cmds...)
		}

		// 全局退出快捷键（仅在主菜单时生效）
		if msg.String() == "ctrl+c" {
			m.quitting = true
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case common.NavigateMsg:
		// 导航到新页面
		m.pageStack = append(m.pageStack, m.currentPage)
		m.currentPage = m.createPage(msg.Page, msg.Params)
		return m, m.currentPage.Init()

	case common.BackMsg:
		// 返回上一页
		if len(m.pageStack) > 0 {
			m.currentPage = m.pageStack[len(m.pageStack)-1]
			m.pageStack = m.pageStack[:len(m.pageStack)-1]
			return m, common.Refresh()
		}

	case common.RefreshMsg:
		// 刷新当前页面
		return m, m.currentPage.Init()

	case common.ToastMsg:
		// 显示提示消息
		m.toast.Show(msg.Message, components.ToastType(msg.Type))
		cmds = append(cmds, m.toast.HideAfter())

	case common.CloseToastMsg:
		// 关闭提示消息
		m.toast.Hide()

	case common.ConfirmMsg:
		// 显示确认对话框
		m.confirm.Show(msg.Title, msg.Message, msg.OnConfirm, msg.OnCancel)

	case common.ConfirmResultMsg:
		// 确认对话框结果
		if msg.Confirmed {
			if cmd := m.confirm.GetOnConfirm(); cmd != nil {
				cmds = append(cmds, cmd)
			}
		} else {
			if cmd := m.confirm.GetOnCancel(); cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
		m.confirm.Hide()

	case common.CloseConfirmMsg:
		// 关闭确认对话框
		m.confirm.Hide()
	}

	// 更新当前页面
	if !m.confirm.IsVisible() {
		newPage, cmd := m.currentPage.Update(msg)
		m.currentPage = newPage.(common.Page)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	// 更新 Toast
	newToast, cmd := m.toast.Update(msg)
	m.toast = newToast.(*components.Toast)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// View 渲染界面
func (m *AppModel) View() string {
	if m.quitting {
		quitStyle := lipgloss.NewStyle().
			Foreground(styles.Primary).
			Bold(true)
		return quitStyle.Render("再见~ 👋") + "\n"
	}

	// 设置组件尺寸
	m.toast.SetSize(m.width, m.height)
	m.confirm.SetSize(m.width, m.height)

	// 渲染当前页面
	content := m.currentPage.View()

	// 主视图
	mainView := content

	// 如果有 Toast，叠加显示（使用 Overlay 居中）
	if m.toast.IsVisible() {
		mainView = m.toast.RenderOverlay(mainView)
	}

	// 如果有确认对话框，叠加显示（使用 Overlay 居中）
	if m.confirm.IsVisible() {
		mainView = m.confirm.RenderOverlay(mainView)
	}

	return mainView
}

// createPage 创建页面
// 呀~ 根据页面类型创建对应的页面模型！✨
func (m *AppModel) createPage(pageType common.PageType, params map[string]any) common.Page {
	switch pageType {
	case common.PageMainMenu:
		return NewMenuModel(m.bs)
	case common.PageMemoryList:
		return memory.NewListModel(m.bs)
	case common.PageMemoryCreate:
		return memory.NewCreateModel(m.bs)
	case common.PageMemoryDetail:
		id := getIntParam(params, "id")
		return memory.NewDetailModel(m.bs, id)
	case common.PageMemorySearch:
		return memory.NewSearchModel(m.bs)
	case common.PagePlanList:
		return plan.NewListModel(m.bs)
	case common.PagePlanCreate:
		return plan.NewCreateModel(m.bs)
	case common.PagePlanDetail:
		id := getIntParam(params, "id")
		return plan.NewDetailModel(m.bs, id)
	case common.PagePlanProgress:
		id := getIntParam(params, "id")
		progress := getIntParam(params, "progress")
		return plan.NewProgressModel(m.bs, id, progress)
	case common.PageTodoList:
		return todo.NewListModel(m.bs)
	case common.PageTodoToday:
		return todo.NewTodayModel(m.bs)
	case common.PageTodoCreate:
		return todo.NewCreateModel(m.bs)
	case common.PageTodoDetail:
		id := getIntParam(params, "id")
		return todo.NewDetailModel(m.bs, id)
	case common.PageGroupList:
		return group.NewListModel(m.bs)
	case common.PageGroupCreate:
		return group.NewCreateModel(m.bs)
	case common.PageGroupDetail:
		id := getIntParam(params, "id")
		return group.NewDetailModel(m.bs, id)
	default:
		return NewMenuModel(m.bs)
	}
}

// getIntParam 从参数中获取整数
func getIntParam(params map[string]any, key string) int {
	if params == nil {
		return 0
	}
	if v, ok := params[key]; ok {
		if i, ok := v.(int); ok {
			return i
		}
	}
	return 0
}
