package layout

import (
	"github.com/XiaoLFeng/llm-memory/internal/tui/components"
	"github.com/charmbracelet/lipgloss"
)

// 统一的卡片宽度计算配置
const (
	defaultMargin   = 4  // 内容左右预留
	defaultMinWidth = 48 // 列表/表单卡片最小宽度（含边框与 padding）
	defaultMaxWidth = 80 // 列表/表单卡片最大宽度
)

// DefaultCardWidth 根据 Frame 内容宽度计算推荐卡片宽度
func DefaultCardWidth(frame *components.Frame) int {
	return components.FitCardWidth(frame.GetContentWidth(), defaultMargin, defaultMinWidth, defaultMaxWidth)
}

// CardInnerWidth 由卡片总宽度换算出内容可用宽度（去掉边框和左右 padding）
func CardInnerWidth(cardWidth int) int {
	inner := cardWidth - 6 // 2 边框 + 4 padding
	if inner < 10 {
		inner = 10
	}
	return inner
}

// ListPage 统一的列表页布局：顶部面包屑 + 居中卡片列表 + 底部快捷键栏
func ListPage(frame *components.Frame, breadcrumb, title, listContent string, keys []string, extra string) string {
	cardWidth := DefaultCardWidth(frame)
	card := components.Card(title, listContent, cardWidth)

	body := lipgloss.Place(
		frame.GetContentWidth(),
		frame.GetContentHeight(),
		lipgloss.Center,
		lipgloss.Top,
		card,
	)

	return frame.Render(breadcrumb, body, keys, extra)
}

// FormPage 统一的表单页布局：顶部面包屑 + 表单卡片 + 底部快捷键栏
func FormPage(frame *components.Frame, breadcrumb, title, formContent string, keys []string, extra string) string {
	cardWidth := DefaultCardWidth(frame)
	card := components.Card(title, formContent, cardWidth)

	body := lipgloss.Place(
		frame.GetContentWidth(),
		frame.GetContentHeight(),
		lipgloss.Center,
		lipgloss.Top,
		card,
	)

	return frame.Render(breadcrumb, body, keys, extra)
}

// DetailPage 统一的详情页布局：将已有内容直接交给 Frame 包裹
// 详情页内部通常已经处理了 viewport/换行，不需要再强制卡片。
func DetailPage(frame *components.Frame, breadcrumb, content string, keys []string, extra string) string {
	return frame.Render(breadcrumb, content, keys, extra)
}

// MenuPage 菜单页布局：内容垂直居中，后续可统一背景/边距策略
func MenuPage(width, height int, content string) string {
	return lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}
