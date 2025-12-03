package styles

import "github.com/charmbracelet/lipgloss"

// Teal/Cyan 青绿色配色方案
// 嘿嘿~ 这是现代化的青绿色配色定义！清新又优雅！(´∀｀)💖
var (
	// 基础色 - 背景层次 (Slate 深色系)
	Base     = lipgloss.Color("#0f172a") // Slate-900 最深背景
	Mantle   = lipgloss.Color("#1e293b") // Slate-800 侧边栏/卡片底
	Crust    = lipgloss.Color("#0c1222") // 最暗边缘
	Surface0 = lipgloss.Color("#334155") // Slate-700 卡片背景
	Surface1 = lipgloss.Color("#475569") // Slate-600 悬浮/选中背景
	Surface2 = lipgloss.Color("#64748b") // Slate-500 更亮的表面

	// 主色调 - 青绿色系 (Teal)
	Primary       = lipgloss.Color("#2dd4bf") // Teal-400 主色
	PrimaryDim    = lipgloss.Color("#0d9488") // Teal-600 暗色
	PrimaryBright = lipgloss.Color("#5eead4") // Teal-300 亮色

	// 强调色
	Accent   = lipgloss.Color("#22d3ee") // Cyan-400 链接/强调
	Lavender = lipgloss.Color("#67e8f9") // Cyan-300 次要强调
	Teal     = lipgloss.Color("#14b8a6") // Teal-500 特殊强调
	Emerald  = lipgloss.Color("#10b981") // Emerald-500 翠绿色

	// 语义色
	Success = lipgloss.Color("#4ade80") // Green-400 成功
	Warning = lipgloss.Color("#fbbf24") // Amber-400 警告
	Error   = lipgloss.Color("#f87171") // Red-400 错误
	Info    = lipgloss.Color("#38bdf8") // Sky-400 信息

	// 文字色
	Text     = lipgloss.Color("#e2e8f0") // Slate-200 主文字
	Subtext1 = lipgloss.Color("#cbd5e1") // Slate-300 次要文字
	Subtext0 = lipgloss.Color("#94a3b8") // Slate-400 更暗文字
	Overlay2 = lipgloss.Color("#64748b") // Slate-500 占位符
	Overlay1 = lipgloss.Color("#475569") // Slate-600 禁用文字
	Overlay0 = lipgloss.Color("#334155") // Slate-700 最暗文字

	// 边框色
	Border       = lipgloss.Color("#334155") // Slate-700 默认边框
	BorderFocus  = lipgloss.Color("#2dd4bf") // Teal-400 聚焦边框
	BorderSubtle = lipgloss.Color("#475569") // Slate-600 微妙边框
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
