package styles

import "github.com/charmbracelet/lipgloss"

// 边框样式
// 嘿嘿~ 这些是各种边框配置！✨
var (
	// 圆角边框 - 主要使用
	RoundedBorder = lipgloss.RoundedBorder()

	// 普通边框 - 次要使用
	NormalBorder = lipgloss.NormalBorder()

	// 双线边框 - 强调使用
	DoubleBorder = lipgloss.DoubleBorder()

	// 粗边框 - 标题使用
	ThickBorder = lipgloss.ThickBorder()
)

// 卡片样式
// 呀~ 这些是各种卡片容器样式！💖
var (
	// 基础卡片 - 带圆角边框
	CardStyle = lipgloss.NewStyle().
			Border(RoundedBorder).
			BorderForeground(Border).
			Padding(1, 2)

	// 聚焦卡片 - 紫色边框
	CardFocusedStyle = lipgloss.NewStyle().
				Border(RoundedBorder).
				BorderForeground(Primary).
				Padding(1, 2)

	// 成功卡片 - 绿色边框
	CardSuccessStyle = lipgloss.NewStyle().
				Border(RoundedBorder).
				BorderForeground(Success).
				Padding(1, 2)

	// 错误卡片 - 红色边框
	CardErrorStyle = lipgloss.NewStyle().
			Border(RoundedBorder).
			BorderForeground(Error).
			Padding(1, 2)

	// 警告卡片 - 黄色边框
	CardWarningStyle = lipgloss.NewStyle().
				Border(RoundedBorder).
				BorderForeground(Warning).
				Padding(1, 2)

	// 信息卡片 - 蓝色边框
	CardInfoStyle = lipgloss.NewStyle().
			Border(RoundedBorder).
			BorderForeground(Info).
			Padding(1, 2)
)

// 页面框架样式
// 嘿嘿~ 这是统一的页面框架样式！🎮
var (
	// 标题栏样式
	HeaderStyle = lipgloss.NewStyle().
			Border(RoundedBorder).
			BorderForeground(Border).
			Foreground(Text).
			Padding(0, 1)

	// 内容区样式
	ContentStyle = lipgloss.NewStyle().
			Padding(1, 2)

	// 状态栏样式
	FooterStyle = lipgloss.NewStyle().
			Border(RoundedBorder).
			BorderForeground(Border).
			Foreground(Subtext0).
			Padding(0, 1)
)

// Logo 和品牌样式
var (
	// Logo 样式
	LogoStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Primary)

	// 品牌文字
	BrandStyle = lipgloss.NewStyle().
			Foreground(Subtext1).
			Italic(true)

	// Logo 文本
	LogoText = "🧠 LLM-Memory"

	// 分隔符
	Separator = " ┃ "
)

// 列表项样式 - 现代化
var (
	// 列表项图标
	ListItemIcon = "📝"

	// 选中指示器
	SelectedIndicator = "▸"

	// 未选中占位
	UnselectedIndicator = " "

	// 列表项主标题样式
	ListItemTitleStyle = lipgloss.NewStyle().
				Foreground(Text).
				Bold(true)

	// 列表项主标题样式（选中）
	ListItemTitleSelectedStyle = lipgloss.NewStyle().
					Foreground(Text).
					Bold(true).
					Background(Surface1)

	// 列表项描述样式
	ListItemDescStyle = lipgloss.NewStyle().
				Foreground(Subtext0)

	// 列表项描述样式（选中）
	ListItemDescSelectedStyle = lipgloss.NewStyle().
					Foreground(Subtext1).
					Background(Surface1)

	// 列表项元信息样式
	ListItemMetaStyle = lipgloss.NewStyle().
				Foreground(Overlay0)

	// 元信息分隔符
	MetaSeparator = " │ "
)

// 徽章样式
var (
	// 作用域徽章样式
	BadgeGlobalStyle = lipgloss.NewStyle().
				Foreground(Accent).
				Bold(true)

	BadgeGroupStyle = lipgloss.NewStyle().
			Foreground(Teal).
			Bold(true)

	BadgePersonalStyle = lipgloss.NewStyle().
				Foreground(Pink).
				Bold(true)

	// 优先级徽章样式
	BadgeLowStyle = lipgloss.NewStyle().
			Foreground(Overlay0)

	BadgeMediumStyle = lipgloss.NewStyle().
				Foreground(Accent)

	BadgeHighStyle = lipgloss.NewStyle().
			Foreground(Warning).
			Bold(true)

	BadgeUrgentStyle = lipgloss.NewStyle().
				Foreground(Error).
				Bold(true)
)

// 时间戳样式
var (
	TimeStyle = lipgloss.NewStyle().
		Foreground(Overlay1).
		Italic(true)
)

// 标签样式
var (
	TagStyle = lipgloss.NewStyle().
		Foreground(Lavender)
)

// 内容区域嵌套卡片样式
var (
	// 嵌套卡片 - 用于详情页的信息分组
	NestedCardStyle = lipgloss.NewStyle().
			Border(RoundedBorder).
			BorderForeground(BorderSubtle).
			Padding(0, 1)

	// 嵌套卡片标题
	NestedCardTitleStyle = lipgloss.NewStyle().
				Foreground(Subtext1).
				Bold(true)
)
