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

// TUI 最大尺寸限制常量
const (
	MaxWidth  = 120 // 最大宽度（列数）
	MaxHeight = 40  // 最大高度（行数）
)

// AppModel 根应用模型
type AppModel struct {
	bs              *startup.Bootstrap
	pageStack       []common.Page       // 页面栈
	currentPage     common.Page         // 当前页面
	toast           *components.Toast   // 提示消息
	confirm         *components.Confirm // 确认对话框
	width           int                 // 终端实际宽度
	height          int                 // 终端实际高度
	effectiveWidth  int                 // 有效宽度（受限后）
	effectiveHeight int                 // 有效高度（受限后）
	quitting        bool
}

// NewAppModel 创建根应用模型
func NewAppModel(bs *startup.Bootstrap) *AppModel {
	menu := NewMenuModel(bs)
	return &AppModel{
		bs:              bs,
		pageStack:       []common.Page{},
		currentPage:     menu,
		toast:           components.NewToast(),
		confirm:         components.NewConfirm(),
		width:           80,
		height:          24,
		effectiveWidth:  80,
		effectiveHeight: 24,
	}
}

// Init 初始化
func (m *AppModel) Init() tea.Cmd {
	return tea.Batch(m.currentPage.Init(), tea.WindowSize())
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
		// 计算有效尺寸（受最大限制）
		m.effectiveWidth = msg.Width
		m.effectiveHeight = msg.Height
		if m.effectiveWidth > MaxWidth {
			m.effectiveWidth = MaxWidth
		}
		if m.effectiveHeight > MaxHeight {
			m.effectiveHeight = MaxHeight
		}
		// 传递限制后的尺寸给子页面
		effectiveMsg := tea.WindowSizeMsg{
			Width:  m.effectiveWidth,
			Height: m.effectiveHeight,
		}
		newPage, cmd := m.currentPage.Update(effectiveMsg)
		m.currentPage = newPage.(common.Page)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
		return m, tea.Batch(cmds...)

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

	case common.AutoRefreshMsg:
		// 自动刷新消息，转发给当前页面处理
		// 同时启动下一轮自动刷新
		cmds = append(cmds, common.StartAutoRefresh())
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
		return quitStyle.Render("再见~ "+styles.IconWave) + "\n"
	}

	// 设置组件尺寸（使用有效尺寸）
	m.toast.SetSize(m.effectiveWidth, m.effectiveHeight)
	m.confirm.SetSize(m.effectiveWidth, m.effectiveHeight)

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

	// 如果终端尺寸大于有效尺寸，居中显示
	if m.width > m.effectiveWidth || m.height > m.effectiveHeight {
		mainView = lipgloss.Place(
			m.width,
			m.height,
			lipgloss.Center,
			lipgloss.Center,
			mainView,
		)
	}

	return mainView
}

// createPage 创建页面
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
		return plan.NewDetailModel(m.bs, int64(id))
	case common.PagePlanProgress:
		id := getIntParam(params, "id")
		progress := getIntParam(params, "progress")
		return plan.NewProgressModel(m.bs, int64(id), progress)
	case common.PageTodoList:
		return todo.NewListModel(m.bs)
	case common.PageTodoCreate:
		return todo.NewCreateModel(m.bs)
	case common.PageTodoDetail:
		id := getIntParam(params, "id")
		return todo.NewDetailModel(m.bs, int64(id))
	case common.PageGroupList:
		return group.NewListModel(m.bs)
	case common.PageGroupCreate:
		return group.NewCreateModel(m.bs)
	case common.PageGroupDetail:
		id := getIntParam(params, "id")
		return group.NewDetailModel(m.bs, int64(id))
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
