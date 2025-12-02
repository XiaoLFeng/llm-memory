package components

import (
	"strings"

	"github.com/XiaoLFeng/llm-memory/internal/tui/styles"
	"github.com/charmbracelet/lipgloss"
)

// Card 创建卡片容器
// 嘿嘿~ 这是通用的卡片组件，带标题和边框！💖
func Card(title, content string, width int) string {
	if width < 20 {
		width = 20
	}

	// 计算标题行
	titleLine := createTitleLine(title, width-4)

	cardStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Border).
		Width(width).
		Padding(1, 2)

	innerContent := titleLine + "\n" + content
	return cardStyle.Render(innerContent)
}

// CardWithColor 带自定义边框颜色的卡片
func CardWithColor(title, content string, width int, borderColor lipgloss.Color) string {
	if width < 20 {
		width = 20
	}

	titleLine := createTitleLine(title, width-4)

	cardStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Width(width).
		Padding(1, 2)

	innerContent := titleLine + "\n" + content
	return cardStyle.Render(innerContent)
}

// CardSimple 简单卡片（无标题）
func CardSimple(content string, width int) string {
	if width < 20 {
		width = 20
	}

	cardStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Border).
		Width(width).
		Padding(1, 2)

	return cardStyle.Render(content)
}

// CardFocused 聚焦状态的卡片
func CardFocused(title, content string, width int) string {
	return CardWithColor(title, content, width, styles.Primary)
}

// CardSuccess 成功状态的卡片
func CardSuccess(title, content string, width int) string {
	return CardWithColor(title, content, width, styles.Success)
}

// CardError 错误状态的卡片
func CardError(title, content string, width int) string {
	return CardWithColor(title, content, width, styles.Error)
}

// CardWarning 警告状态的卡片
func CardWarning(title, content string, width int) string {
	return CardWithColor(title, content, width, styles.Warning)
}

// CardInfo 信息状态的卡片
func CardInfo(title, content string, width int) string {
	return CardWithColor(title, content, width, styles.Info)
}

// NestedCard 嵌套卡片（用于详情页的信息分组）
func NestedCard(title, content string, width int) string {
	if width < 20 {
		width = 20
	}

	titleLine := createTitleLine(title, width-4)

	nestedStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.BorderSubtle).
		Width(width).
		Padding(0, 1)

	innerContent := titleLine + "\n" + content
	return nestedStyle.Render(innerContent)
}

// createTitleLine 创建标题行（带装饰线）
func createTitleLine(title string, width int) string {
	if title == "" {
		return ""
	}

	titleStyle := lipgloss.NewStyle().
		Foreground(styles.Subtext1).
		Bold(true)

	styledTitle := titleStyle.Render(title)
	titleWidth := lipgloss.Width(styledTitle)

	// 计算右侧需要的破折号数量
	dashCount := width - titleWidth - 3
	if dashCount < 0 {
		dashCount = 0
	}

	lineStyle := lipgloss.NewStyle().Foreground(styles.BorderSubtle)
	dashes := lineStyle.Render(strings.Repeat("─", dashCount))

	return "─ " + styledTitle + " " + dashes
}

// InfoRow 信息行（用于详情页的键值对显示）
func InfoRow(label, value string) string {
	labelStyle := lipgloss.NewStyle().
		Foreground(styles.Subtext0).
		Width(12)

	valueStyle := lipgloss.NewStyle().
		Foreground(styles.Text)

	return labelStyle.Render(label+":") + " " + valueStyle.Render(value)
}

// InfoGrid 信息网格（多列显示）
func InfoGrid(items [][]string, colWidth int) string {
	if len(items) == 0 {
		return ""
	}

	var rows []string
	for _, row := range items {
		var cols []string
		for _, item := range row {
			colStyle := lipgloss.NewStyle().Width(colWidth)
			cols = append(cols, colStyle.Render(item))
		}
		rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, cols...))
	}

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}
