package tui

import (
	"github.com/XiaoLFeng/llm-memory/internal/tui/app"
	"github.com/XiaoLFeng/llm-memory/startup"
	tea "github.com/charmbracelet/bubbletea"
)

// Run 启动新 TUI
func Run(bs *startup.Bootstrap) error {
	prog := tea.NewProgram(app.New(bs), tea.WithAltScreen())
	_, err := prog.Run()
	return err
}
