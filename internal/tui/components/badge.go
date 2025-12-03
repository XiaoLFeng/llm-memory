package components

import (
	"github.com/XiaoLFeng/llm-memory/internal/tui/styles"
	"github.com/charmbracelet/lipgloss"
)

// Badge 标签徽章组件
// 嘿嘿~ 用于显示状态、优先级、作用域等标签！🏷️

// ScopeBadge 作用域徽章
func ScopeBadge(scope string) string {
	var style lipgloss.Style
	var text string

	switch scope {
	case "global", "Global":
		style = lipgloss.NewStyle().
			Foreground(styles.Accent).
			Bold(true)
		text = "[Global]"
	case "group", "Group":
		style = lipgloss.NewStyle().
			Foreground(styles.Teal).
			Bold(true)
		text = "[Group]"
	case "personal", "Personal":
		style = lipgloss.NewStyle().
			Foreground(styles.Pink).
			Bold(true)
		text = "[Personal]"
	default:
		style = lipgloss.NewStyle().Foreground(styles.Overlay0)
		text = "[Unknown]"
	}

	return style.Render(text)
}

// ScopeBadgeFromGroupIDPath 根据 GroupID 和 Path 生成作用域徽章
func ScopeBadgeFromGroupIDPath(groupID uint, path string) string {
	if path != "" {
		return ScopeBadge("Personal")
	}
	if groupID != 0 {
		return ScopeBadge("Group")
	}
	return ScopeBadge("Global")
}

// PriorityBadge 优先级徽章
func PriorityBadge(priority int) string {
	var style lipgloss.Style
	var text string

	switch priority {
	case 1:
		style = lipgloss.NewStyle().Foreground(styles.Overlay0)
		text = "⬇️低"
	case 2:
		style = lipgloss.NewStyle().Foreground(styles.Accent)
		text = "➡️中"
	case 3:
		style = lipgloss.NewStyle().
			Foreground(styles.Warning).
			Bold(true)
		text = "⬆️高"
	case 4:
		style = lipgloss.NewStyle().
			Foreground(styles.Error).
			Bold(true)
		text = "🔥紧急"
	default:
		style = lipgloss.NewStyle().Foreground(styles.Overlay0)
		text = "未知"
	}

	return style.Render(text)
}

// PriorityBadgeSimple 简单优先级徽章（仅图标）
func PriorityBadgeSimple(priority int) string {
	switch priority {
	case 1:
		return "⬇️"
	case 2:
		return "➡️"
	case 3:
		return "⬆️"
	case 4:
		return "🔥"
	default:
		return "•"
	}
}

// StatusBadge 状态徽章
func StatusBadge(status string) string {
	var style lipgloss.Style
	var text string

	switch status {
	case "pending", "待开始", "待处理":
		style = lipgloss.NewStyle().Foreground(styles.Overlay0)
		text = "⏳待处理"
	case "in_progress", "进行中":
		style = lipgloss.NewStyle().Foreground(styles.Info)
		text = "🔄进行中"
	case "completed", "已完成":
		style = lipgloss.NewStyle().Foreground(styles.Success)
		text = "✅已完成"
	case "cancelled", "已取消":
		style = lipgloss.NewStyle().Foreground(styles.Error)
		text = "❌已取消"
	default:
		style = lipgloss.NewStyle().Foreground(styles.Overlay0)
		text = "未知"
	}

	return style.Render(text)
}

// StatusBadgeSimple 简单状态徽章（仅图标）
func StatusBadgeSimple(status string) string {
	switch status {
	case "pending", "待开始", "待处理":
		return "⏳"
	case "in_progress", "进行中":
		return "🔄"
	case "completed", "已完成":
		return "✅"
	case "cancelled", "已取消":
		return "❌"
	default:
		return "•"
	}
}

// CategoryBadge 分类徽章
func CategoryBadge(category string) string {
	style := lipgloss.NewStyle().
		Foreground(styles.Lavender)
	return style.Render(category)
}

// TagBadge 标签徽章
func TagBadge(tag string) string {
	style := lipgloss.NewStyle().
		Foreground(styles.Lavender)
	return style.Render("#" + tag)
}

// TagsBadge 多标签徽章
func TagsBadge(tags []string) string {
	if len(tags) == 0 {
		return ""
	}

	style := lipgloss.NewStyle().
		Foreground(styles.Lavender)

	result := ""
	for i, tag := range tags {
		if i > 0 {
			result += " "
		}
		result += style.Render("#" + tag)
	}
	return result
}

// ProgressBadge 进度徽章
func ProgressBadge(progress int) string {
	var style lipgloss.Style

	if progress == 0 {
		style = lipgloss.NewStyle().Foreground(styles.Overlay0)
	} else if progress < 50 {
		style = lipgloss.NewStyle().Foreground(styles.Warning)
	} else if progress < 100 {
		style = lipgloss.NewStyle().Foreground(styles.Info)
	} else {
		style = lipgloss.NewStyle().Foreground(styles.Success)
	}

	return style.Render(string(rune('0'+progress/10)) + string(rune('0'+progress%10)) + "%")
}

// CountBadge 计数徽章
func CountBadge(count int) string {
	style := lipgloss.NewStyle().
		Foreground(styles.Subtext0)
	return style.Render("(" + string(rune('0'+count/10)) + string(rune('0'+count%10)) + ")")
}

// TimeBadge 时间徽章
func TimeBadge(timeStr string) string {
	style := lipgloss.NewStyle().
		Foreground(styles.Overlay1).
		Italic(true)
	return style.Render(timeStr)
}
