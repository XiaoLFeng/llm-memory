package utils

import (
	"github.com/XiaoLFeng/llm-memory/internal/tui/theme"
)

// RenderDetailSection 渲染详情页的区块
// 参数：
// - icon: Nerd Font 图标（如 📝、📄）
// - title: 区块标题
// - content: 区块内容（支持多行）
// - width: 可用宽度
// 返回：格式化后的行数组
func RenderDetailSection(icon, title, content string, width int) []string {
	// 空值检查（优雅处理）
	if content == "" {
		return []string{}
	}

	var lines []string

	// === 小标题行 ===
	// 使用 theme.Subtitle 样式 + 图标
	titleLine := theme.Subtitle.Render(icon + " " + title)
	lines = append(lines, titleLine)
	lines = append(lines, "") // 标题下方留空行

	// === 内容行 ===
	// 自动换行处理
	contentLines := WrapText(content, width)

	// 使用 theme.TextMain 样式渲染每一行
	for _, line := range contentLines {
		if line == "" {
			// 保留空行（段落分隔）
			lines = append(lines, "")
		} else {
			// 渲染文本内容
			lines = append(lines, theme.TextMain.Render(line))
		}
	}

	return lines
}
