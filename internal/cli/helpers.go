package cli

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// 样式定义
// 嘿嘿~ 使用 Charmbracelet lipgloss 美化输出！💖
var (
	// TitleStyle 标题样式
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#A78BFA"))

	// SuccessStyle 成功消息样式
	SuccessStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#22C55E"))

	// ErrorStyle 错误消息样式
	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#EF4444"))

	// WarningStyle 警告消息样式
	WarningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F59E0B"))

	// InfoStyle 信息样式
	InfoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#3B82F6"))

	// MutedStyle 次要文本样式
	MutedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#64748B"))
)

// PrintSuccess 打印成功消息
// 呀~ 操作成功时使用这个！✨
func PrintSuccess(message string) {
	fmt.Println(SuccessStyle.Render("✓ " + message))
}

// PrintError 打印错误消息
func PrintError(message string) {
	fmt.Println(ErrorStyle.Render("✗ " + message))
}

// PrintWarning 打印警告消息
func PrintWarning(message string) {
	fmt.Println(WarningStyle.Render("⚠ " + message))
}

// PrintInfo 打印信息消息
func PrintInfo(message string) {
	fmt.Println(InfoStyle.Render("ℹ " + message))
}

// PrintTitle 打印标题
func PrintTitle(title string) {
	fmt.Println(TitleStyle.Render(title))
}
