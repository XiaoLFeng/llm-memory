package common

import (
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// PageType 页面类型枚举
// 呀~ 定义所有可能的页面类型！📄
type PageType int

const (
	PageMainMenu PageType = iota
	PageMemoryList
	PageMemoryCreate
	PageMemoryDetail
	PageMemorySearch
	PagePlanList
	PagePlanCreate
	PagePlanDetail
	PagePlanProgress
	PageTodoList
	PageTodoCreate
	PageTodoDetail
	PageGroupList   // 组列表
	PageGroupCreate // 创建组
	PageGroupDetail // 组详情
)

// ToastType 提示消息类型
type ToastType int

const (
	ToastSuccess ToastType = iota
	ToastError
	ToastWarning
	ToastInfo
)

// Page 页面接口
// 嘿嘿~ 所有页面都要实现这个接口！✨
type Page interface {
	tea.Model
	Title() string
	ShortHelp() []key.Binding
}

// 消息类型定义
// 这些消息用于页面间通信~

// NavigateMsg 导航消息
// 用于跳转到指定页面
type NavigateMsg struct {
	Page   PageType
	Params map[string]any
}

// BackMsg 返回消息
// 用于返回上一页
type BackMsg struct{}

// RefreshMsg 刷新消息
// 用于刷新当前页面数据
type RefreshMsg struct{}

// ToastMsg 提示消息
// 用于显示操作反馈
type ToastMsg struct {
	Message string
	Type    ToastType
}

// ConfirmMsg 确认对话框消息
// 用于危险操作前的确认
type ConfirmMsg struct {
	Title     string
	Message   string
	OnConfirm tea.Cmd
	OnCancel  tea.Cmd
}

// ConfirmResultMsg 确认结果消息
type ConfirmResultMsg struct {
	Confirmed bool
}

// CloseConfirmMsg 关闭确认对话框消息
type CloseConfirmMsg struct{}

// CloseToastMsg 关闭提示消息
type CloseToastMsg struct{}

// AutoRefreshMsg 自动刷新消息
type AutoRefreshMsg struct{}

// AutoRefreshInterval 自动刷新间隔 (30秒)
const AutoRefreshInterval = 30 * time.Second

// WindowSizeMsg 窗口大小消息
type WindowSizeMsg struct {
	Width  int
	Height int
}

// 通用快捷键定义
var (
	KeyQuit = key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "退出"),
	)

	KeyBack = key.NewBinding(
		key.WithKeys("esc", "backspace"),
		key.WithHelp("esc", "返回"),
	)

	KeyEnter = key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "确认"),
	)

	KeyUp = key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "上移"),
	)

	KeyDown = key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "下移"),
	)

	KeyCreate = key.NewBinding(
		key.WithKeys("c", "n"),
		key.WithHelp("c/n", "新建"),
	)

	KeyDelete = key.NewBinding(
		key.WithKeys("d", "delete"),
		key.WithHelp("d", "删除"),
	)

	KeySearch = key.NewBinding(
		key.WithKeys("/", "s"),
		key.WithHelp("/", "搜索"),
	)

	KeyHelp = key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "帮助"),
	)

	KeyTab = key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "下一项"),
	)

	KeyShiftTab = key.NewBinding(
		key.WithKeys("shift+tab"),
		key.WithHelp("shift+tab", "上一项"),
	)
)

// Navigate 创建导航命令
func Navigate(page PageType, params ...map[string]any) tea.Cmd {
	return func() tea.Msg {
		p := make(map[string]any)
		if len(params) > 0 {
			p = params[0]
		}
		return NavigateMsg{Page: page, Params: p}
	}
}

// Back 创建返回命令
func Back() tea.Cmd {
	return func() tea.Msg {
		return BackMsg{}
	}
}

// Refresh 创建刷新命令
func Refresh() tea.Cmd {
	return func() tea.Msg {
		return RefreshMsg{}
	}
}

// ShowToast 创建显示提示消息命令
func ShowToast(message string, toastType ToastType) tea.Cmd {
	return func() tea.Msg {
		return ToastMsg{Message: message, Type: toastType}
	}
}

// ShowConfirm 创建显示确认对话框命令
func ShowConfirm(title, message string, onConfirm, onCancel tea.Cmd) tea.Cmd {
	return func() tea.Msg {
		return ConfirmMsg{
			Title:     title,
			Message:   message,
			OnConfirm: onConfirm,
			OnCancel:  onCancel,
		}
	}
}

// StartAutoRefresh 启动自动刷新计时器
func StartAutoRefresh() tea.Cmd {
	return tea.Tick(AutoRefreshInterval, func(t time.Time) tea.Msg {
		return AutoRefreshMsg{}
	})
}

// 额外的快捷键定义
var (
	KeyEdit = key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "编辑"),
	)

	KeyToggle = key.NewBinding(
		key.WithKeys("space"),
		key.WithHelp("space", "切换"),
	)

	KeyConfirm = key.NewBinding(
		key.WithKeys("y"),
		key.WithHelp("y", "确认"),
	)

	KeyCancel = key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "取消"),
	)

	KeyRefresh = key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "刷新"),
	)

	KeyFilter = key.NewBinding(
		key.WithKeys("f"),
		key.WithHelp("f", "筛选"),
	)

	KeySort = key.NewBinding(
		key.WithKeys("o"),
		key.WithHelp("o", "排序"),
	)

	KeySave = key.NewBinding(
		key.WithKeys("ctrl+s"),
		key.WithHelp("ctrl+s", "保存"),
	)
)
