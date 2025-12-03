package app

import (
	"github.com/XiaoLFeng/llm-memory/internal/tui/components"
	tea "github.com/charmbracelet/bubbletea"
)

// Meta 描述页面的辅助信息
type Meta struct {
	Title      string
	Breadcrumb string
	Extra      string
	Keys       []components.KeyHint
}

// Page 页面接口
type Page interface {
	Init() tea.Cmd
	Update(msg tea.Msg) (Page, tea.Cmd)
	View() string
	Meta() Meta
	Resize(w, h int)
}
