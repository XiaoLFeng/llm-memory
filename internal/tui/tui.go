package tui

import (
	"github.com/XiaoLFeng/llm-memory/internal/tui/models"
	"github.com/XiaoLFeng/llm-memory/startup"
	tea "github.com/charmbracelet/bubbletea"
)

// Run 运行 TUI
// 呀~ 启动终端用户界面！✨
func Run(bs *startup.Bootstrap) error {
	app := models.NewAppModel(bs)
	p := tea.NewProgram(app, tea.WithAltScreen())
	_, err := p.Run()
	return err
}
