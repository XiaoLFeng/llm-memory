package theme

import "github.com/charmbracelet/lipgloss"

// 常用间距
const (
	PadX = 2
	PadY = 1
)

// Header/Footer/Body 样式
var (
	Header = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Border).
		Background(Mantle).
		Foreground(Text).
		Padding(0, 1)

	Footer = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Border).
		Background(Mantle).
		Foreground(Subtext0).
		Padding(0, 1)

	Content = lipgloss.NewStyle().
		Background(Base).
		Padding(PadY, PadX)
)

// 文本
var (
	Title     = lipgloss.NewStyle().Foreground(Primary).Bold(true)
	Subtitle  = lipgloss.NewStyle().Foreground(Accent)
	TextMain  = lipgloss.NewStyle().Foreground(Text)
	TextDim   = lipgloss.NewStyle().Foreground(Subtext0)
	MutedText = lipgloss.NewStyle().Foreground(Muted)
)

// 状态键提示
var (
	KeyStyle   = lipgloss.NewStyle().Foreground(Primary).Background(Mantle).Bold(true)
	ValueStyle = lipgloss.NewStyle().Foreground(Text).Background(Mantle)
)

// 卡片
var (
	Card = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Border).
		Background(Surface0).
		Padding(1, 2)

	CardStrong = Card.Copy().BorderForeground(BorderStrong)
)

// 选中行
var (
	Row         = lipgloss.NewStyle().Foreground(Text)
	RowSelected = lipgloss.NewStyle().
			Foreground(Text).
			Background(Surface1)
)
