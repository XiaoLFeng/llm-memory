package styles

import "github.com/charmbracelet/lipgloss"

// Catppuccin Mocha 配色方案
// 嘿嘿~ 这是现代化的配色定义！采用 Catppuccin Mocha 风格！(´∀｀)💖
var (
	// 基础色 - 背景层次
	Base     = lipgloss.Color("#1E1E2E") // 最深背景
	Mantle   = lipgloss.Color("#181825") // 侧边栏/卡片底
	Crust    = lipgloss.Color("#11111B") // 最暗边缘
	Surface0 = lipgloss.Color("#313244") // 卡片背景
	Surface1 = lipgloss.Color("#45475A") // 悬浮/选中背景
	Surface2 = lipgloss.Color("#585B70") // 更亮的表面

	// 主色调 - 紫色系（保持品牌一致）
	Primary       = lipgloss.Color("#CBA6F7") // Mauve - 主紫色
	PrimaryDim    = lipgloss.Color("#A78BFA") // 暗紫
	PrimaryBright = lipgloss.Color("#DDB6FF") // 亮紫

	// 强调色
	Accent   = lipgloss.Color("#89B4FA") // Blue - 链接/强调
	Lavender = lipgloss.Color("#B4BEFE") // 薰衣草 - 次要强调
	Teal     = lipgloss.Color("#94E2D5") // Teal - 特殊强调
	Pink     = lipgloss.Color("#F5C2E7") // Pink - 装饰色

	// 语义色
	Success = lipgloss.Color("#A6E3A1") // Green - 成功
	Warning = lipgloss.Color("#F9E2AF") // Yellow - 警告
	Error   = lipgloss.Color("#F38BA8") // Red/Pink - 错误
	Info    = lipgloss.Color("#89DCEB") // Sky - 信息

	// 文字色
	Text     = lipgloss.Color("#CDD6F4") // 主文字
	Subtext1 = lipgloss.Color("#BAC2DE") // 次要文字
	Subtext0 = lipgloss.Color("#A6ADC8") // 更暗文字
	Overlay2 = lipgloss.Color("#9399B2") // 占位符
	Overlay1 = lipgloss.Color("#7F849C") // 禁用文字
	Overlay0 = lipgloss.Color("#6C7086") // 最暗文字

	// 边框色
	Border       = lipgloss.Color("#45475A") // 默认边框
	BorderFocus  = lipgloss.Color("#CBA6F7") // 聚焦边框
	BorderSubtle = lipgloss.Color("#313244") // 微妙边框
)

// 优先级颜色映射
// 呀~ 不同优先级用不同颜色标记！🎨
var PriorityColors = map[int]lipgloss.Color{
	1: Overlay0, // 低 - 灰色
	2: Accent,   // 中 - 蓝色
	3: Warning,  // 高 - 橙色
	4: Error,    // 紧急 - 红色
}

// 状态颜色映射
var StatusColors = map[string]lipgloss.Color{
	"pending":     Overlay0, // 待处理 - 灰色
	"in_progress": Info,     // 进行中 - 蓝色
	"completed":   Success,  // 已完成 - 绿色
	"cancelled":   Error,    // 已取消 - 红色
}
