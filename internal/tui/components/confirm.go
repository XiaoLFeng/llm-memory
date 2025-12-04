package components

import (
	"github.com/XiaoLFeng/llm-memory/internal/tui/theme"
	"github.com/charmbracelet/lipgloss"
)

// ConfirmDialog 渲染确认对话框（静态函数）
func ConfirmDialog(title, message, hint string, width int) string {
	if width < 40 {
		width = 40
	}

	titleStr := theme.ConfirmTitle.Render(theme.IconWarning + " " + title)
	msgStr := theme.ConfirmMessage.Width(width - 6).Render(message)
	hintStr := theme.ConfirmHint.Render(hint)

	content := lipgloss.JoinVertical(lipgloss.Left, titleStr, msgStr, hintStr)

	return theme.ConfirmBox.Width(width).Render(content)
}

// ConfirmDialogWithButtons 带按钮的确认对话框
func ConfirmDialogWithButtons(title, message string, width int, yesSelected bool) string {
	if width < 40 {
		width = 40
	}

	titleStr := theme.ConfirmTitle.Render(theme.IconWarning + " " + title)
	msgStr := theme.ConfirmMessage.Width(width - 6).Render(message)

	// 按钮
	yesStyle := theme.ConfirmButton
	noStyle := theme.ConfirmButton
	if yesSelected {
		yesStyle = theme.ConfirmButtonActive
	} else {
		noStyle = theme.ConfirmButtonActive
	}

	buttons := lipgloss.JoinHorizontal(lipgloss.Left,
		yesStyle.Render("[Y] 确认"),
		noStyle.Render("[N] 取消"),
	)

	content := lipgloss.JoinVertical(lipgloss.Center, titleStr, msgStr, buttons)

	return theme.ConfirmBox.Width(width).Render(content)
}

// DeleteConfirmDialog 删除确认对话框
func DeleteConfirmDialog(itemName string, width int) string {
	title := "确认删除"
	message := "确定要删除「" + itemName + "」吗？\n此操作不可撤销。"
	hint := "[Y] 确认删除  [N/Esc] 取消"

	return ConfirmDialog(title, message, hint, width)
}
