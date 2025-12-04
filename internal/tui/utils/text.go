package utils

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
