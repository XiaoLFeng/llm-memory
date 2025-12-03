package components

import (
	"strings"

	"github.com/XiaoLFeng/llm-memory/internal/tui/styles"
	"github.com/charmbracelet/lipgloss"
)

// Position 浮动位置
type Position int

const (
	// TopCenter 顶部居中
	TopCenter Position = iota
	// Center 屏幕中央
	Center
	// BottomCenter 底部居中
	BottomCenter
)

// PlaceOverlay 将浮动内容放置在基础内容上
func PlaceOverlay(base, overlay string, width, height int, pos Position) string {
	baseLines := strings.Split(base, "\n")
	overlayLines := strings.Split(overlay, "\n")

	overlayWidth := lipgloss.Width(overlay)
	overlayHeight := len(overlayLines)

	// 计算 overlay 的起始位置
	var startY int
	switch pos {
	case TopCenter:
		startY = 2
	case Center:
		startY = (height - overlayHeight) / 2
	case BottomCenter:
		startY = height - overlayHeight - 2
	}

	startX := (width - overlayWidth) / 2
	if startX < 0 {
		startX = 0
	}
	if startY < 0 {
		startY = 0
	}

	// 确保 base 有足够的行数
	for len(baseLines) < height {
		baseLines = append(baseLines, strings.Repeat(" ", width))
	}

	// 创建结果
	result := make([]string, len(baseLines))
	copy(result, baseLines)

	// 叠加 overlay
	for i, overlayLine := range overlayLines {
		lineY := startY + i
		if lineY >= 0 && lineY < len(result) {
			// 获取当前行
			baseLine := result[lineY]
			baseRunes := []rune(baseLine)

			// 确保行足够宽
			for len(baseRunes) < width {
				baseRunes = append(baseRunes, ' ')
			}

			// 创建新行
			newLine := string(baseRunes[:startX])
			newLine += overlayLine
			endX := startX + lipgloss.Width(overlayLine)
			if endX < len(baseRunes) {
				newLine += string(baseRunes[endX:])
			}

			result[lineY] = newLine
		}
	}

	return strings.Join(result, "\n")
}

// PlaceOverlayWithDim 带半透明遮罩的浮动
func PlaceOverlayWithDim(base, overlay string, width, height int, pos Position) string {
	// 先将背景变暗
	dimmedBase := dimContent(base)
	// 再放置浮动内容
	return PlaceOverlay(dimmedBase, overlay, width, height, pos)
}

// dimContent 将内容变暗（添加遮罩效果）
func dimContent(content string) string {
	dimStyle := lipgloss.NewStyle().
		Foreground(styles.Overlay0)

	lines := strings.Split(content, "\n")
	for i, line := range lines {
		// 移除原有样式并应用变暗样式
		lines[i] = dimStyle.Render(stripAnsi(line))
	}
	return strings.Join(lines, "\n")
}

// stripAnsi 移除 ANSI 转义序列（简化版本）
func stripAnsi(s string) string {
	var result strings.Builder
	inEscape := false

	for _, r := range s {
		if r == '\x1b' {
			inEscape = true
			continue
		}
		if inEscape {
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
				inEscape = false
			}
			continue
		}
		result.WriteRune(r)
	}

	return result.String()
}

// CenterBox 创建居中的盒子
func CenterBox(content string, width, height int, borderColor lipgloss.Color) string {
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(1, 2).
		Align(lipgloss.Center)

	box := boxStyle.Render(content)

	// 在空白背景上居中
	containerStyle := lipgloss.NewStyle().
		Width(width).
		Height(height).
		Align(lipgloss.Center, lipgloss.Center)

	return containerStyle.Render(box)
}
