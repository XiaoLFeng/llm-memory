package output

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// 表格样式定义
var (
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#A78BFA"))

	cellStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E2E8F0"))

	separatorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#475569"))
)

// Table 表格结构
// 嘿嘿~ 用于渲染漂亮的表格输出！(´∀｀)💖
type Table struct {
	Headers []string
	Rows    [][]string
}

// NewTable 创建新表格
func NewTable(headers ...string) *Table {
	return &Table{
		Headers: headers,
		Rows:    make([][]string, 0),
	}
}

// AddRow 添加一行数据
func (t *Table) AddRow(cells ...string) {
	t.Rows = append(t.Rows, cells)
}

// Render 渲染表格
// 呀~ 生成漂亮的表格字符串！✨
func (t *Table) Render() string {
	if len(t.Headers) == 0 {
		return ""
	}

	// 计算每列最大宽度
	colWidths := make([]int, len(t.Headers))
	for i, h := range t.Headers {
		colWidths[i] = runeWidth(h)
	}
	for _, row := range t.Rows {
		for i, cell := range row {
			if i < len(colWidths) {
				w := runeWidth(cell)
				if w > colWidths[i] {
					colWidths[i] = w
				}
			}
		}
	}

	var sb strings.Builder

	// 渲染表头
	sb.WriteString(renderRow(t.Headers, colWidths, headerStyle))
	sb.WriteString("\n")
	sb.WriteString(renderSeparator(colWidths))
	sb.WriteString("\n")

	// 渲染数据行
	for _, row := range t.Rows {
		sb.WriteString(renderRow(row, colWidths, cellStyle))
		sb.WriteString("\n")
	}

	return sb.String()
}

// String 实现 Stringer 接口
func (t *Table) String() string {
	return t.Render()
}

// Print 直接打印表格
func (t *Table) Print() {
	fmt.Print(t.Render())
}

// renderRow 渲染一行
func renderRow(cells []string, widths []int, style lipgloss.Style) string {
	var parts []string
	for i, cell := range cells {
		if i < len(widths) {
			// 计算填充空格数
			padLen := widths[i] - runeWidth(cell)
			if padLen < 0 {
				padLen = 0
			}
			paddedCell := cell + strings.Repeat(" ", padLen)
			parts = append(parts, style.Render(paddedCell))
		}
	}
	return strings.Join(parts, separatorStyle.Render(" │ "))
}

// renderSeparator 渲染分隔线
func renderSeparator(widths []int) string {
	var parts []string
	for _, w := range widths {
		parts = append(parts, strings.Repeat("─", w))
	}
	return separatorStyle.Render(strings.Join(parts, "─┼─"))
}

// runeWidth 计算字符串显示宽度（考虑中文字符）
func runeWidth(s string) int {
	width := 0
	for _, r := range s {
		if r >= 0x4E00 && r <= 0x9FFF || // CJK统一汉字
			r >= 0x3400 && r <= 0x4DBF || // CJK扩展A
			r >= 0xFF00 && r <= 0xFFEF { // 全角字符
			width += 2
		} else {
			width += 1
		}
	}
	return width
}
