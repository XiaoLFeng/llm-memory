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
