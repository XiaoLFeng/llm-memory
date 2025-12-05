package utils

import "strings"

// LipWidth 计算字符串宽度（处理中文字符）
// 使用 rune 长度而非字节长度，确保 Unicode 字符正确计数
func LipWidth(s string) int {
	return len([]rune(s))
}

// Truncate 截断字符串并添加省略号
// 当字符串超过指定限制时，截断并在末尾添加 "…"
func Truncate(s string, limit int) string {
	runes := []rune(s)
	if len(runes) <= limit {
		return s
	}
	if limit <= 1 {
		return string(runes[:limit])
	}
	return string(runes[:limit-1]) + "…"
}

// WrapText 将长文本按指定宽度自动换行
// 特点：
// 1. 保留原始段落分隔（空行）
// 2. 在空格/标点处优雅断行
// 3. 支持 Unicode 字符（中文等）
func WrapText(text string, width int) []string {
	// 空值检查（优雅处理）
	if text == "" {
		return []string{}
	}

	// 防御性编程：默认宽度
	if width <= 0 {
		width = 60
	}

	// 按换行符分割段落
	paragraphs := strings.Split(text, "\n")
	var result []string

	for _, para := range paragraphs {
		// 保留空行（Markdown 段落分隔）
		if strings.TrimSpace(para) == "" {
			result = append(result, "")
			continue
		}

		// 按宽度切分
		wrapped := wrapLine(para, width)
		result = append(result, wrapped...)
	}

	return result
}

// wrapLine 将单行文本按宽度切分（私有函数）
// 特性：
// - 基于 rune 计算长度（支持 Unicode）
// - 尝试在空格/标点处断行（优雅换行）
// - 避免单词中间断开
func wrapLine(line string, width int) []string {
	runes := []rune(line)

	// 短于宽度，直接返回
	if len(runes) <= width {
		return []string{line}
	}

	var lines []string
	start := 0

	for start < len(runes) {
		// 计算本行的结束位置
		end := start + width
		if end > len(runes) {
			end = len(runes)
		}

		// 尝试在空格/标点处断行（优雅换行）
		if end < len(runes) {
			// 向前查找最近的断点（最多回退 20 个字符）
			bestBreak := end
			maxLookback := 20
			if end-start < maxLookback {
				maxLookback = end - start
			}

			for i := end; i > start && i > end-maxLookback; i-- {
				r := runes[i]
				// 在空格或中文标点处断行
				if r == ' ' || r == '、' || r == '，' || r == '。' ||
					r == '；' || r == '：' || r == '！' || r == '？' {
					bestBreak = i
					break
				}
			}

			// 如果找到更好的断点，使用它
			if bestBreak != end {
				end = bestBreak + 1 // 包含标点符号
			}
		}

		// 添加本行
		lines = append(lines, string(runes[start:end]))
		start = end

		// 跳过前导空格
		for start < len(runes) && runes[start] == ' ' {
			start++
		}
	}

	return lines
}
