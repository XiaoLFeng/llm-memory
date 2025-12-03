package styles

import "github.com/charmbracelet/lipgloss"

// 颜色定义 - 兼容旧代码
// 嘿嘿~ 使用 Catppuccin Mocha 配色方案！💖
// 注意：这些是兼容旧代码的别名，新代码请使用 colors.go 中的定义
var (
	// 主色调 - 紫色系（映射到新配色）
	PrimaryColor   = Primary
	SecondaryColor = PrimaryDim

	// 状态色（映射到新配色）
	SuccessColor = Success
	ErrorColor   = Error
	WarningColor = Warning
	InfoColor    = Info

	// 文字色（映射到新配色）
	TextColor      = Text
	MutedColor     = Overlay0
	HighlightColor = Accent

	// 背景色（映射到新配色）
	BgColor      = Base
	BgLightColor = Surface0
	BorderColor  = Border
)

// 通用样式
// 呀~ 这些样式可以在整个 TUI 中复用！✨
var (
	// 标题样式 - 使用新配色
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Primary).
			MarginBottom(1)

	// 副标题样式
	SubtitleStyle = lipgloss.NewStyle().
			Foreground(Accent).
			MarginBottom(1)

	// 选中项样式 - 带背景
	SelectedStyle = lipgloss.NewStyle().
			Foreground(Text).
			Background(Surface1).
			Padding(0, 1)

	// 普通项样式
	NormalStyle = lipgloss.NewStyle().
			Foreground(Accent).
			Padding(0, 1)

	// 帮助文本样式
	HelpStyle = lipgloss.NewStyle().
			Foreground(Overlay0)

	// 成功样式
	SuccessStyle = lipgloss.NewStyle().
			Foreground(Success)

	// 错误样式
	ErrorStyle = lipgloss.NewStyle().
			Foreground(Error)

	// 警告样式
	WarningStyle = lipgloss.NewStyle().
			Foreground(Warning)

	// 信息样式
	InfoStyle = lipgloss.NewStyle().
			Foreground(Info)

	// 静默文本样式
	MutedStyle = lipgloss.NewStyle().
			Foreground(Overlay0)

	// 描述文本样式
	DescStyle = lipgloss.NewStyle().
			Foreground(Subtext0).
			Italic(true)
)

// 表单样式
var (
	// 聚焦状态的输入框 - 使用新配色
	FocusedInputStyle = lipgloss.NewStyle().
				Foreground(Text).
				BorderForeground(Primary).
				BorderStyle(lipgloss.RoundedBorder()).
				Padding(0, 1)

	// 未聚焦状态的输入框
	BlurredInputStyle = lipgloss.NewStyle().
				Foreground(Subtext0).
				BorderForeground(Border).
				BorderStyle(lipgloss.RoundedBorder()).
				Padding(0, 1)

	// 标签样式
	LabelStyle = lipgloss.NewStyle().
			Foreground(Text).
			Bold(true)

	// 占位符样式
	PlaceholderStyle = lipgloss.NewStyle().
				Foreground(Overlay0)
)

// 列表样式
var (
	// 列表标题样式 - 使用新配色
	ListTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Primary).
			Padding(0, 1).
			MarginBottom(1)

	// 列表项样式
	ListItemStyle = lipgloss.NewStyle().
			Foreground(Text).
			Padding(0, 1)

	// 列表选中项样式 - 带背景
	ListSelectedStyle = lipgloss.NewStyle().
				Foreground(Text).
				Background(Surface1).
				Padding(0, 1)

	// 列表项描述样式
	ListDescStyle = lipgloss.NewStyle().
			Foreground(Overlay0).
			Padding(0, 1)
)

// 状态栏样式
var (
	StatusBarStyle = lipgloss.NewStyle().
			Foreground(Subtext0).
			Background(Mantle).
			Padding(0, 1)

	StatusKeyStyle = lipgloss.NewStyle().
			Foreground(Primary).
			Background(Mantle).
			Bold(true)

	StatusValueStyle = lipgloss.NewStyle().
				Foreground(Text).
				Background(Mantle)
)

// 对话框样式
var (
	DialogStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Primary).
			Padding(1, 2)

	DialogTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(Primary).
				MarginBottom(1)

	DialogButtonStyle = lipgloss.NewStyle().
				Foreground(Text).
				Background(Surface0).
				Padding(0, 2).
				MarginRight(1)

	DialogButtonActiveStyle = lipgloss.NewStyle().
				Foreground(Text).
				Background(Primary).
				Padding(0, 2).
				MarginRight(1)
)

// 优先级样式
// 嘿嘿~ 不同优先级用不同颜色标记！🎨
func PriorityStyle(priority int) lipgloss.Style {
	switch priority {
	case 1: // 低
		return lipgloss.NewStyle().Foreground(Overlay0)
	case 2: // 中
		return lipgloss.NewStyle().Foreground(Accent)
	case 3: // 高
		return lipgloss.NewStyle().Foreground(Warning).Bold(true)
	case 4: // 紧急
		return lipgloss.NewStyle().Foreground(Error).Bold(true)
	default:
		return lipgloss.NewStyle().Foreground(Text)
	}
}

// 状态样式
func StatusStyle(status string) lipgloss.Style {
	switch status {
	case "pending", "待开始", "待处理":
		return lipgloss.NewStyle().Foreground(Overlay0)
	case "in_progress", "进行中":
		return lipgloss.NewStyle().Foreground(Info)
	case "completed", "已完成":
		return lipgloss.NewStyle().Foreground(Success)
	case "cancelled", "已取消":
		return lipgloss.NewStyle().Foreground(Error)
	default:
		return lipgloss.NewStyle().Foreground(Text)
	}
}
