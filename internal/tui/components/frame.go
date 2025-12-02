package components

import (
	"strings"

	"github.com/XiaoLFeng/llm-memory/internal/tui/styles"
	"github.com/charmbracelet/lipgloss"
)

// Frame 页面框架组件
// 嘿嘿~ 这是统一的页面框架，所有页面都用它！💖
type Frame struct {
	width  int
	height int
}

// NewFrame 创建页面框架
func NewFrame(width, height int) *Frame {
	return &Frame{
		width:  width,
		height: height,
	}
}

// SetSize 设置尺寸
func (f *Frame) SetSize(width, height int) {
	f.width = width
	f.height = height
}

// GetContentHeight 获取内容区域高度
func (f *Frame) GetContentHeight() int {
	// 标题栏高度 3 + 状态栏高度 3 + 边距 2
	return f.height - 8
}

// GetContentWidth 获取内容区域宽度
func (f *Frame) GetContentWidth() int {
	// 边距 4
	return f.width - 4
}

// Render 渲染完整框架
// breadcrumb: 面包屑导航，如 "记忆管理 > 记忆列表"
// content: 主要内容
// keys: 快捷键列表，如 []string{"↑/↓ 移动", "Enter 确认", "esc 返回"}
// extra: 额外信息，显示在标题栏右侧
func (f *Frame) Render(breadcrumb, content string, keys []string, extra string) string {
	// 计算宽度
	contentWidth := f.width - 2
	if contentWidth < 40 {
		contentWidth = 40
	}

	// 标题栏
	header := f.renderHeader(breadcrumb, extra, contentWidth)

	// 内容区域
	mainContent := f.renderContent(content, contentWidth)

	// 状态栏
	footer := f.renderFooter(keys, contentWidth)

	return lipgloss.JoinVertical(lipgloss.Left, header, mainContent, footer)
}

// RenderSimple 渲染简单框架（无状态栏）
func (f *Frame) RenderSimple(breadcrumb, content string) string {
	contentWidth := f.width - 2
	if contentWidth < 40 {
		contentWidth = 40
	}

	header := f.renderHeader(breadcrumb, "", contentWidth)
	mainContent := f.renderContent(content, contentWidth)

	return lipgloss.JoinVertical(lipgloss.Left, header, mainContent)
}

// renderHeader 渲染标题栏
func (f *Frame) renderHeader(breadcrumb, extra string, width int) string {
	// Logo
	logo := styles.LogoStyle.Render(styles.LogoText)

	// 面包屑
	nav := styles.SubtitleStyle.Render(breadcrumb)

	// 左侧内容
	left := logo + styles.Separator + nav

	// 右侧额外信息
	right := ""
	if extra != "" {
		right = styles.MutedStyle.Render(extra)
	}

	// 计算填充
	leftWidth := lipgloss.Width(left)
	rightWidth := lipgloss.Width(right)
	padding := width - leftWidth - rightWidth - 4
	if padding < 0 {
		padding = 0
	}

	content := left + strings.Repeat(" ", padding) + right

	headerStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Border).
		Foreground(styles.Text).
		Width(width).
		Padding(0, 1)

	return headerStyle.Render(content)
}

// renderContent 渲染内容区域
func (f *Frame) renderContent(content string, width int) string {
	contentHeight := f.GetContentHeight()
	if contentHeight < 5 {
		contentHeight = 5
	}

	contentStyle := lipgloss.NewStyle().
		Width(width).
		Height(contentHeight).
		Padding(1, 1)

	return contentStyle.Render(content)
}

// renderFooter 渲染状态栏
func (f *Frame) renderFooter(keys []string, width int) string {
	if len(keys) == 0 {
		return ""
	}

	// 用分隔符连接快捷键
	keysStr := strings.Join(keys, "  │  ")

	footerStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Border).
		Foreground(styles.Subtext0).
		Width(width).
		Padding(0, 1)

	return footerStyle.Render(keysStr)
}

// RenderWithCard 渲染带卡片包装的框架
func (f *Frame) RenderWithCard(breadcrumb, cardTitle, content string, keys []string, extra string) string {
	// 将内容用卡片包装
	cardContent := Card(cardTitle, content, f.GetContentWidth()-4)
	return f.Render(breadcrumb, cardContent, keys, extra)
}
