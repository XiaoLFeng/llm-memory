package utils

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"
)

// FormatTime 格式化时间
// 嘿嘿~ 将时间格式化为友好的显示格式！⏰
func FormatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04")
}

// FormatDate 格式化日期
func FormatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// FormatTimePtr 格式化时间指针
func FormatTimePtr(t *time.Time) string {
	if t == nil {
		return "-"
	}
	return FormatTime(*t)
}

// FormatDatePtr 格式化日期指针
func FormatDatePtr(t *time.Time) string {
	if t == nil {
		return "-"
	}
	return FormatDate(*t)
}

// FormatRelativeTime 格式化相对时间
// 呀~ 显示"刚刚"、"5分钟前"这样的友好格式！✨
func FormatRelativeTime(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	switch {
	case diff < time.Minute:
		return "刚刚"
	case diff < time.Hour:
		return fmt.Sprintf("%d分钟前", int(diff.Minutes()))
	case diff < 24*time.Hour:
		return fmt.Sprintf("%d小时前", int(diff.Hours()))
	case diff < 7*24*time.Hour:
		return fmt.Sprintf("%d天前", int(diff.Hours()/24))
	default:
		return FormatDate(t)
	}
}

// FormatProgress 格式化进度条
// 嘿嘿~ 用方块字符显示进度！📊
func FormatProgress(progress int, width int) string {
	if width <= 0 {
		width = 10
	}
	if progress < 0 {
		progress = 0
	}
	if progress > 100 {
		progress = 100
	}

	filled := width * progress / 100
	empty := width - filled

	return fmt.Sprintf("[%s%s] %3d%%",
		strings.Repeat("█", filled),
		strings.Repeat("░", empty),
		progress,
	)
}

// FormatPriority 格式化优先级
func FormatPriority(priority int) string {
	switch priority {
	case 1:
		return "低"
	case 2:
		return "中"
	case 3:
		return "高"
	case 4:
		return "紧急"
	default:
		return "未知"
	}
}

// FormatPriorityIcon 格式化优先级图标
func FormatPriorityIcon(priority int) string {
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
		return "❓"
	}
}

// FormatStatus 格式化状态
func FormatStatus(status string) string {
	switch status {
	case "pending":
		return "待开始"
	case "in_progress":
		return "进行中"
	case "completed":
		return "已完成"
	case "cancelled":
		return "已取消"
	default:
		return status
	}
}

// FormatStatusIcon 格式化状态图标
func FormatStatusIcon(status string) string {
	switch status {
	case "pending":
		return "⏳"
	case "in_progress":
		return "🔄"
	case "completed":
		return "✅"
	case "cancelled":
		return "❌"
	default:
		return "❓"
	}
}

// FormatTodoStatus 格式化待办状态
func FormatTodoStatus(status int) string {
	switch status {
	case 0:
		return "待处理"
	case 1:
		return "进行中"
	case 2:
		return "已完成"
	case 3:
		return "已取消"
	default:
		return "未知"
	}
}

// FormatTodoStatusIcon 格式化待办状态图标
func FormatTodoStatusIcon(status int) string {
	switch status {
	case 0:
		return "📋"
	case 1:
		return "🔄"
	case 2:
		return "✅"
	case 3:
		return "❌"
	default:
		return "❓"
	}
}

// Truncate 截断字符串
// 呀~ 如果字符串太长就截断并加省略号！📏
func Truncate(s string, maxLen int) string {
	if maxLen <= 0 {
		return ""
	}
	if utf8.RuneCountInString(s) <= maxLen {
		return s
	}
	runes := []rune(s)
	if maxLen <= 3 {
		return string(runes[:maxLen])
	}
	return string(runes[:maxLen-3]) + "..."
}

// PadRight 右填充字符串
func PadRight(s string, width int) string {
	runeCount := utf8.RuneCountInString(s)
	if runeCount >= width {
		return s
	}
	return s + strings.Repeat(" ", width-runeCount)
}

// PadLeft 左填充字符串
func PadLeft(s string, width int) string {
	runeCount := utf8.RuneCountInString(s)
	if runeCount >= width {
		return s
	}
	return strings.Repeat(" ", width-runeCount) + s
}

// Center 居中字符串
func Center(s string, width int) string {
	runeCount := utf8.RuneCountInString(s)
	if runeCount >= width {
		return s
	}
	left := (width - runeCount) / 2
	right := width - runeCount - left
	return strings.Repeat(" ", left) + s + strings.Repeat(" ", right)
}

// WrapText 自动换行
// 嘿嘿~ 将长文本按指定宽度自动换行！📝
func WrapText(text string, width int) string {
	if width <= 0 {
		return text
	}

	var result strings.Builder
	lines := strings.Split(text, "\n")

	for i, line := range lines {
		if i > 0 {
			result.WriteString("\n")
		}

		runes := []rune(line)
		for len(runes) > width {
			result.WriteString(string(runes[:width]))
			result.WriteString("\n")
			runes = runes[width:]
		}
		result.WriteString(string(runes))
	}

	return result.String()
}

// JoinTags 连接标签
func JoinTags(tags []string) string {
	if len(tags) == 0 {
		return "-"
	}
	return strings.Join(tags, ", ")
}

// RuneWidth 计算字符串显示宽度（考虑中文字符）
func RuneWidth(s string) int {
	width := 0
	for _, r := range s {
		if r >= 0x4E00 && r <= 0x9FFF || // CJK 统一汉字
			r >= 0x3000 && r <= 0x303F || // CJK 标点
			r >= 0xFF00 && r <= 0xFFEF { // 全角字符
			width += 2
		} else {
			width += 1
		}
	}
	return width
}
