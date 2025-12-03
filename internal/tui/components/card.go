package components

import (
	"strings"

	"github.com/XiaoLFeng/llm-memory/internal/tui/styles"
	"github.com/charmbracelet/lipgloss"
)

// FitCardWidth 根据容器宽度计算卡片合适宽度（含边框与 padding）
// margin 代表左右总预留空间（例如 Place 时两侧 2 列留白可设为 4）
// min/max 为最终卡片总宽度的上下限（含边框与 padding）
func FitCardWidth(containerWidth, margin, min, max int) int {
	if margin < 0 {
		margin = 0
	}
	usable := containerWidth - margin
	if max > 0 && usable > max {
		usable = max
	}
	if usable < min {
		usable = min
	}
	if usable < 20 {
		usable = 20
	}
	return usable
}

// RenderCard 渲染卡片并居中/定高放置
func RenderCard(frame *Frame, title, content string, minWidth, maxWidth, margin int, vPos lipgloss.Position) string {
	width := FitCardWidth(frame.GetContentWidth(), margin, minWidth, maxWidth)
	card := Card(title, content, width)
	return lipgloss.Place(
		frame.GetContentWidth(),
		frame.GetContentHeight(),
		lipgloss.Center,
		vPos,
		card,
	)
}

// Card 创建卡片容器
// width 参数代表卡片的期望总宽度（包含边框和 padding）
func Card(title, content string, width int) string {
	if width < 20 {
		width = 20
	}

	// 边框占 2 (左 1 + 右 1)
	// Padding 占 4 (左 2 + 右 2)
	// 所以内容区域 = width - 6
	contentWidth := width - 6
	if contentWidth < 10 {
		contentWidth = 10
	}

	// 计算标题行
	titleLine := createTitleLine(title, contentWidth)

	cardStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Border).
		Width(contentWidth).
		Padding(1, 2)

	innerContent := titleLine + "\n" + content
	return cardStyle.Render(innerContent)
}

// CardWithColor 带自定义边框颜色的卡片
// width 参数代表卡片的期望总宽度（包含边框和 padding）
func CardWithColor(title, content string, width int, borderColor lipgloss.Color) string {
	if width < 20 {
		width = 20
	}

	// 边框占 2，Padding 占 4
	contentWidth := width - 6
	if contentWidth < 10 {
		contentWidth = 10
	}

	titleLine := createTitleLine(title, contentWidth)

	cardStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Width(contentWidth).
		Padding(1, 2)

	innerContent := titleLine + "\n" + content
	return cardStyle.Render(innerContent)
}

// CardSimple 简单卡片（无标题）
// width 参数代表卡片的期望总宽度（包含边框和 padding）
func CardSimple(content string, width int) string {
	if width < 20 {
		width = 20
	}

	// 边框占 2，Padding 占 4
	contentWidth := width - 6
	if contentWidth < 10 {
		contentWidth = 10
	}

	cardStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Border).
		Width(contentWidth).
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
// width 参数代表卡片的期望总宽度（包含边框和 padding）
func NestedCard(title, content string, width int) string {
	if width < 20 {
		width = 20
	}

	// 边框占 2，Padding 占 2 (左右各1)
	contentWidth := width - 4
	if contentWidth < 10 {
		contentWidth = 10
	}

	titleLine := createTitleLine(title, contentWidth)

	nestedStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.BorderSubtle). // 使用提高对比度的边框色
		Width(contentWidth).
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

	// 如果宽度不够显示完整标题，截断标题
	if dashCount == 0 && titleWidth > width-3 {
		maxTitleLen := width - 6 // 留出 "─ " + " ─" 的空间
		if maxTitleLen > 3 {
			// 截断标题（考虑中文字符）
			runes := []rune(title)
			if len(runes) > maxTitleLen-3 {
				title = string(runes[:maxTitleLen-3]) + "..."
				styledTitle = titleStyle.Render(title)
				titleWidth = lipgloss.Width(styledTitle)
				dashCount = width - titleWidth - 3
				if dashCount < 0 {
					dashCount = 0
				}
			}
		}
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
