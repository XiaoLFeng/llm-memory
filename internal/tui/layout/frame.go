package layout

import (
	"strings"

	"github.com/XiaoLFeng/llm-memory/internal/tui/theme"
	"github.com/charmbracelet/lipgloss"
)

// Frame 统一框架：Header / Body / Footer
type Frame struct {
	width  int
	height int
}

func NewFrame(w, h int) *Frame {
	return &Frame{width: w, height: h}
}

func (f *Frame) Resize(w, h int) { f.width, f.height = w, h }

func (f *Frame) ContentSize() (int, int) {
	contentH := f.height - 8 // header+footer+padding
	if contentH < 5 {
		contentH = 5
	}
	contentW := f.width - 2
	if contentW < 40 {
		contentW = 40
	}
	return contentW, contentH
}

// Render 将 body 放入框架
func (f *Frame) Render(breadcrumb, extra, body string, keys []string) string {
	contentW, contentH := f.ContentSize()

	header := renderHeader(breadcrumb, extra, contentW)
	main := renderBody(body, contentW, contentH)
	footer := renderFooter(keys, contentW)

	return lipgloss.JoinVertical(lipgloss.Left, header, main, footer)
}

func renderHeader(breadcrumb, extra string, width int) string {
	left := theme.Title.Render(theme.IconLogo + " LLM-Memory")
	mid := theme.Subtitle.Render(" ｜ " + breadcrumb)
	right := ""
	if extra != "" {
		right = theme.MutedText.Render(extra)
	}

	gap := width - lipgloss.Width(left) - lipgloss.Width(mid) - lipgloss.Width(right) - 2
	if gap < 0 {
		gap = 0
	}
	content := left + mid + strings.Repeat(" ", gap) + right

	return theme.Header.Copy().Width(width).Render(content)
}

func renderBody(body string, width, height int) string {
	return lipgloss.NewStyle().
		Width(width).
		Height(height).
		Background(theme.Base).
		Padding(1, 1).
		Render(body)
}

func renderFooter(keys []string, width int) string {
	if len(keys) == 0 {
		return theme.Footer.Copy().Width(width).Render("")
	}

	var items []string
	for _, k := range keys {
		items = append(items, theme.KeyStyle.Render(k))
	}
	line := strings.Join(items, "  ")
	return theme.Footer.Copy().Width(width).Render(line)
}

// FitCardWidth 计算卡片总宽度
func FitCardWidth(contentWidth int) int {
	if contentWidth > 96 {
		return 96
	}
	if contentWidth < 48 {
		return 48
	}
	return contentWidth - 4
}
