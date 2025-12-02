package styles

import "github.com/charmbracelet/lipgloss"

// 颜色定义
// 嘿嘿~ 采用 Tailwind CSS 配色方案，统一整个 TUI 的视觉风格！💖
var (
	// 主色调 - 紫色系
	PrimaryColor   = lipgloss.Color("#A78BFA")
	SecondaryColor = lipgloss.Color("#7C3AED")

	// 状态色
	SuccessColor = lipgloss.Color("#22C55E")
	ErrorColor   = lipgloss.Color("#EF4444")
	WarningColor = lipgloss.Color("#F59E0B")
	InfoColor    = lipgloss.Color("#3B82F6")

	// 文字色
	TextColor      = lipgloss.Color("#E2E8F0")
	MutedColor     = lipgloss.Color("#64748B")
	HighlightColor = lipgloss.Color("#93C5FD")

	// 背景色
	BgColor      = lipgloss.Color("#1E1E2E")
	BgLightColor = lipgloss.Color("#313244")
	BorderColor  = lipgloss.Color("#45475A")
)

// 通用样式
// 呀~ 这些样式可以在整个 TUI 中复用！✨
var (
	// 标题样式
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(PrimaryColor).
			MarginBottom(1)

	// 副标题样式
	SubtitleStyle = lipgloss.NewStyle().
			Foreground(HighlightColor).
			MarginBottom(1)

	// 选中项样式
	SelectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(SecondaryColor).
			Padding(0, 1)

	// 普通项样式
	NormalStyle = lipgloss.NewStyle().
			Foreground(HighlightColor).
			Padding(0, 1)

	// 帮助文本样式
	HelpStyle = lipgloss.NewStyle().
			Foreground(MutedColor)

	// 成功样式
	SuccessStyle = lipgloss.NewStyle().
			Foreground(SuccessColor)

	// 错误样式
	ErrorStyle = lipgloss.NewStyle().
			Foreground(ErrorColor)

	// 警告样式
	WarningStyle = lipgloss.NewStyle().
			Foreground(WarningColor)

	// 信息样式
	InfoStyle = lipgloss.NewStyle().
			Foreground(InfoColor)

	// 静默文本样式
	MutedStyle = lipgloss.NewStyle().
			Foreground(MutedColor)

	// 描述文本样式
	DescStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#94A3B8")).
			Italic(true)
)

// 表单样式
var (
	// 聚焦状态的输入框
	FocusedInputStyle = lipgloss.NewStyle().
				Foreground(PrimaryColor).
				BorderForeground(PrimaryColor).
				BorderStyle(lipgloss.RoundedBorder()).
				Padding(0, 1)

	// 未聚焦状态的输入框
	BlurredInputStyle = lipgloss.NewStyle().
				Foreground(MutedColor).
				BorderForeground(BorderColor).
				BorderStyle(lipgloss.RoundedBorder()).
				Padding(0, 1)

	// 标签样式
	LabelStyle = lipgloss.NewStyle().
			Foreground(TextColor).
			Bold(true)

	// 占位符样式
	PlaceholderStyle = lipgloss.NewStyle().
				Foreground(MutedColor)
)

// 列表样式
var (
	// 列表标题样式
	ListTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(PrimaryColor).
			Padding(0, 1).
			MarginBottom(1)

	// 列表项样式
	ListItemStyle = lipgloss.NewStyle().
			Foreground(TextColor).
			Padding(0, 1)

	// 列表选中项样式
	ListSelectedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFFFF")).
				Background(SecondaryColor).
				Padding(0, 1)

	// 列表项描述样式
	ListDescStyle = lipgloss.NewStyle().
			Foreground(MutedColor).
			Padding(0, 1)
)

// 状态栏样式
var (
	StatusBarStyle = lipgloss.NewStyle().
			Foreground(TextColor).
			Background(BgLightColor).
			Padding(0, 1)

	StatusKeyStyle = lipgloss.NewStyle().
			Foreground(PrimaryColor).
			Bold(true)

	StatusValueStyle = lipgloss.NewStyle().
				Foreground(TextColor)
)

// 对话框样式
var (
	DialogStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(PrimaryColor).
			Padding(1, 2)

	DialogTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(PrimaryColor).
				MarginBottom(1)

	DialogButtonStyle = lipgloss.NewStyle().
				Foreground(TextColor).
				Background(BgLightColor).
				Padding(0, 2).
				MarginRight(1)

	DialogButtonActiveStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFFFF")).
				Background(SecondaryColor).
				Padding(0, 2).
				MarginRight(1)
)

// 优先级样式
// 嘿嘿~ 不同优先级用不同颜色标记！🎨
func PriorityStyle(priority int) lipgloss.Style {
	switch priority {
	case 1: // 低
		return lipgloss.NewStyle().Foreground(MutedColor)
	case 2: // 中
		return lipgloss.NewStyle().Foreground(InfoColor)
	case 3: // 高
		return lipgloss.NewStyle().Foreground(WarningColor)
	case 4: // 紧急
		return lipgloss.NewStyle().Foreground(ErrorColor).Bold(true)
	default:
		return lipgloss.NewStyle().Foreground(TextColor)
	}
}

// 状态样式
func StatusStyle(status string) lipgloss.Style {
	switch status {
	case "pending", "待开始":
		return lipgloss.NewStyle().Foreground(MutedColor)
	case "in_progress", "进行中":
		return lipgloss.NewStyle().Foreground(InfoColor)
	case "completed", "已完成":
		return lipgloss.NewStyle().Foreground(SuccessColor)
	case "cancelled", "已取消":
		return lipgloss.NewStyle().Foreground(ErrorColor)
	default:
		return lipgloss.NewStyle().Foreground(TextColor)
	}
}
