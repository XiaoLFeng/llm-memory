package theme

import "github.com/charmbracelet/lipgloss"

// Blue + Teal 清爽配色
var (
	// 基础底色
	Base     = lipgloss.Color("#0b1220") // 深底
	Mantle   = lipgloss.Color("#0f172a") // header/footer 背景
	Surface0 = lipgloss.Color("#122036") // 卡片底
	Surface1 = lipgloss.Color("#16304a") // 选中/悬浮
	Surface2 = lipgloss.Color("#1c3d5f") // 强调底

	// 主色
	Primary       = lipgloss.Color("#22d3ee") // Cyan-400（青）
	PrimaryDim    = lipgloss.Color("#0ea5e9") // Sky-500
	PrimaryBright = lipgloss.Color("#67e8f9")

	// 辅色（蓝）
	Accent     = lipgloss.Color("#38bdf8") // Sky-400
	AccentDeep = lipgloss.Color("#2563eb") // Blue-600

	// 语义
	Success = lipgloss.Color("#34d399")
	Warning = lipgloss.Color("#fbbf24")
	Error   = lipgloss.Color("#f87171")
	Info    = lipgloss.Color("#60a5fa")

	// 文本
	Text     = lipgloss.Color("#e2e8f0")
	Subtext1 = lipgloss.Color("#cbd5e1")
	Subtext0 = lipgloss.Color("#94a3b8")
	Muted    = lipgloss.Color("#64748b")

	// 边框
	Border       = lipgloss.Color("#1f2a3d")
	BorderStrong = lipgloss.Color("#2dd4bf")
)
