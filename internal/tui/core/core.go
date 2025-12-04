package core

import (
	"github.com/XiaoLFeng/llm-memory/internal/tui/components"
	tea "github.com/charmbracelet/bubbletea"
)

// PageID 路由枚举
type PageID string

const (
	PageMenu   PageID = "menu"
	PageMemory PageID = "memory"
	PagePlan   PageID = "plan"
	PageTodo   PageID = "todo"
	PageGroup  PageID = "group"

	// CRUD 页面
	PageMemoryCreate PageID = "memory_create"
	PageMemoryEdit   PageID = "memory_edit"
	PagePlanCreate   PageID = "plan_create"
	PagePlanEdit     PageID = "plan_edit"
	PageTodoCreate   PageID = "todo_create"
	PageTodoEdit     PageID = "todo_edit"
	PageGroupCreate  PageID = "group_create"
	PageGroupEdit    PageID = "group_edit"

	// 帮助页面
	PageHelp PageID = "help"
)

// Page 页面接口
type Page interface {
	Init() tea.Cmd
	Update(msg tea.Msg) (Page, tea.Cmd)
	View() string
	Meta() Meta
	Resize(w, h int)
}

type Meta struct {
	Title      string
	Breadcrumb string
	Extra      string
	Keys       []components.KeyHint
}
