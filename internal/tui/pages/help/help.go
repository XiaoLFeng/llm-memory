package help

import (
	"github.com/XiaoLFeng/llm-memory/internal/tui/components"
	"github.com/XiaoLFeng/llm-memory/internal/tui/core"
	"github.com/XiaoLFeng/llm-memory/internal/tui/layout"
	"github.com/XiaoLFeng/llm-memory/internal/tui/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Page 帮助页面
type Page struct {
	frame  *layout.Frame
	width  int
	height int
	scroll int
	push   func(core.PageID) tea.Cmd
}

// NewPage 创建帮助页面
func NewPage(push func(core.PageID) tea.Cmd) *Page {
	return &Page{
		frame: layout.NewFrame(80, 24),
		width: 80,
		push:  push,
	}
}

// Init 初始化
func (p *Page) Init() tea.Cmd {
	return nil
}

// Update 更新
func (p *Page) Update(msg tea.Msg) (core.Page, tea.Cmd) {
	switch v := msg.(type) {
	case tea.KeyMsg:
		switch v.String() {
		case "up", "k":
			if p.scroll > 0 {
				p.scroll--
			}
		case "down", "j":
			p.scroll++
		case "home":
			p.scroll = 0
		}
	}
	return p, nil
}

// View 渲染
func (p *Page) View() string {
	contentWidth, _ := p.frame.ContentSize()
	cardWidth := layout.FitCardWidth(contentWidth)

	content := p.buildHelpContent(cardWidth)
	return lipgloss.NewStyle().Width(contentWidth).Render(content)
}

// Meta 返回页面元数据
func (p *Page) Meta() core.Meta {
	return core.Meta{
		Title:      "帮助",
		Breadcrumb: theme.IconHelp + " 帮助",
		Extra:      "",
		Keys: []components.KeyHint{
			{Key: "↑/↓", Desc: "滚动"},
			{Key: "Esc", Desc: "返回"},
		},
	}
}

// Resize 调整大小
func (p *Page) Resize(w, h int) {
	p.width = w
	p.height = h
	p.frame.Resize(w, h)
}

// buildHelpContent 构建帮助内容
func (p *Page) buildHelpContent(width int) string {
	titleStyle := theme.Title.Copy().MarginBottom(1)
	sectionStyle := theme.Subtitle.Copy().Bold(true).MarginTop(1).MarginBottom(1)
	keyStyle := theme.KeyStyle.Copy().Width(12)
	descStyle := theme.TextMain

	sections := []string{
		titleStyle.Render("🌊 LLM-Memory TUI 帮助"),
		"",
		sectionStyle.Render("全局快捷键"),
		renderKeyRow(keyStyle, descStyle, "Ctrl+C / q", "退出程序"),
		renderKeyRow(keyStyle, descStyle, "Esc", "返回上一页"),
		renderKeyRow(keyStyle, descStyle, "?", "打开帮助"),
		"",
		sectionStyle.Render("列表页快捷键"),
		renderKeyRow(keyStyle, descStyle, "↑ / k", "向上移动"),
		renderKeyRow(keyStyle, descStyle, "↓ / j", "向下移动"),
		renderKeyRow(keyStyle, descStyle, "Enter", "查看详情 / 切换视图"),
		renderKeyRow(keyStyle, descStyle, "Tab", "切换作用域过滤"),
		renderKeyRow(keyStyle, descStyle, "c", "创建新项"),
		renderKeyRow(keyStyle, descStyle, "e", "编辑选中项"),
		renderKeyRow(keyStyle, descStyle, "d", "删除选中项"),
		renderKeyRow(keyStyle, descStyle, "r", "刷新列表"),
		"",
		sectionStyle.Render("表单页快捷键"),
		renderKeyRow(keyStyle, descStyle, "Tab / ↓", "下一个字段"),
		renderKeyRow(keyStyle, descStyle, "Shift+Tab / ↑", "上一个字段"),
		renderKeyRow(keyStyle, descStyle, "← / →", "切换选项（选择器）"),
		renderKeyRow(keyStyle, descStyle, "Ctrl+S", "保存"),
		renderKeyRow(keyStyle, descStyle, "Esc", "取消并返回"),
		"",
		sectionStyle.Render("删除确认"),
		renderKeyRow(keyStyle, descStyle, "y / Y / Enter", "确认删除"),
		renderKeyRow(keyStyle, descStyle, "n / N / Esc", "取消删除"),
		"",
		sectionStyle.Render("作用域说明"),
		theme.TextDim.Render("  [全局] - 全局可见，所有路径都可访问"),
		theme.TextDim.Render("  [私有] - 仅当前路径可见"),
		theme.TextDim.Render("  [小组] - 组内所有路径可见"),
		"",
		sectionStyle.Render("优先级说明"),
		theme.TextDim.Render("  P1 低   - 低优先级，不紧急"),
		theme.TextDim.Render("  P2 中   - 中等优先级"),
		theme.TextDim.Render("  P3 高   - 高优先级，需要关注"),
		theme.TextDim.Render("  P4 紧急 - 紧急任务，立即处理"),
	}

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// renderKeyRow 渲染快捷键行
func renderKeyRow(keyStyle, descStyle lipgloss.Style, key, desc string) string {
	return lipgloss.JoinHorizontal(lipgloss.Left,
		keyStyle.Render(key),
		descStyle.Render(desc),
	)
}
