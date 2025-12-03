package components

import (
	"github.com/XiaoLFeng/llm-memory/internal/tui/theme"
	"github.com/charmbracelet/lipgloss"
)

// EmptyState 标准空态
func EmptyState(title, hint string, width int) string {
	content := theme.TextDim.Render(hint)
	return Card(title, content, width)
}

// ErrorState 错误态
func ErrorState(title, msg string, width int) string {
	content := lipgloss.NewStyle().Foreground(theme.Error).Render(msg)
	return Card(title, content, width)
}

// LoadingState 加载态
func LoadingState(title, msg string, width int) string {
	content := lipgloss.NewStyle().Foreground(theme.Info).Render(msg)
	return Card(title, content, width)
}
