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

// 表单样式
var (
	FormLabel = lipgloss.NewStyle().
			Foreground(Primary).
			Bold(true).
			MarginBottom(1)

	FormLabelRequired = lipgloss.NewStyle().
				Foreground(Error).
				SetString(" *")

	FormInput = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Border).
			Padding(0, 1)

	FormInputFocused = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(Primary).
				Padding(0, 1)

	FormError = lipgloss.NewStyle().
			Foreground(Error).
			MarginTop(1)

	FormHint = lipgloss.NewStyle().
			Foreground(Muted).
			Italic(true)
)

// 选择器样式
var (
	SelectOption = lipgloss.NewStyle().
			Foreground(Text).
			Padding(0, 1)

	SelectOptionSelected = lipgloss.NewStyle().
				Foreground(Primary).
				Background(Surface1).
				Padding(0, 1).
				Bold(true)

	SelectCursor = lipgloss.NewStyle().
			Foreground(Primary).
			SetString("▸ ")
)

// 确认对话框样式
var (
	ConfirmBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Warning).
			Background(Surface0).
			Padding(1, 2)

	ConfirmTitle = lipgloss.NewStyle().
			Foreground(Warning).
			Bold(true).
			MarginBottom(1)

	ConfirmMessage = lipgloss.NewStyle().
			Foreground(Text).
			MarginBottom(1)

	ConfirmHint = lipgloss.NewStyle().
			Foreground(Subtext0)

	ConfirmButton = lipgloss.NewStyle().
			Foreground(Text).
			Background(Surface1).
			Padding(0, 2).
			MarginRight(1)

	ConfirmButtonActive = lipgloss.NewStyle().
				Foreground(Base).
				Background(Primary).
				Padding(0, 2).
				MarginRight(1).
				Bold(true)
)
