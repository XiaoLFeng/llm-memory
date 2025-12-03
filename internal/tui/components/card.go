package components

import (
	"github.com/XiaoLFeng/llm-memory/internal/tui/theme"
	"github.com/charmbracelet/lipgloss"
)

// Card 渲染卡片，width 为总宽度
func Card(title, content string, width int) string {
	if width < 32 {
		width = 32
	}
	titleLine := theme.Subtitle.Bold(true).Render(title)
	body := lipgloss.NewStyle().
		Width(width - 4).
		Render(content)

	return theme.Card.Copy().
		Width(width).
		Render(titleLine + "\n" + body)
}
