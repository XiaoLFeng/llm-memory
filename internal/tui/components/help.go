package components

import (
	"strings"

	"github.com/XiaoLFeng/llm-memory/internal/tui/styles"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// KeyGroup 快捷键分组
type KeyGroup struct {
	Title string
	Keys  []key.Binding
}

// Help 帮助面板组件
// 呀~ 现代化的帮助面板，支持分组显示快捷键！❓
type Help struct {
	keys      []key.Binding
	keyGroups []KeyGroup
	visible   bool
	width     int
	height    int
}

// NewHelp 创建帮助面板组件
func NewHelp() *Help {
	return &Help{
		width:  60,
		height: 24,
	}
}

// SetKeys 设置快捷键
func (h *Help) SetKeys(keys []key.Binding) {
	h.keys = keys
}

// SetKeyGroups 设置快捷键分组
func (h *Help) SetKeyGroups(groups []KeyGroup) {
	h.keyGroups = groups
}

// SetSize 设置尺寸
func (h *Help) SetSize(width, height int) {
	h.width = width
	h.height = height
}

// Toggle 切换显示状态
func (h *Help) Toggle() {
	h.visible = !h.visible
}

// Show 显示帮助面板
func (h *Help) Show() {
	h.visible = true
}

// Hide 隐藏帮助面板
func (h *Help) Hide() {
	h.visible = false
}

// IsVisible 是否可见
func (h *Help) IsVisible() bool {
	return h.visible
}

// SetWidth 设置宽度
func (h *Help) SetWidth(width int) {
	h.width = width
}

// Init 初始化
func (h *Help) Init() tea.Cmd {
	return nil
}

// Update 处理输入
func (h *Help) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if h.visible && (msg.String() == "?" || msg.String() == "esc") {
			h.Hide()
		}
	case tea.WindowSizeMsg:
		h.width = msg.Width
		h.height = msg.Height
	}
	return h, nil
}

// View 渲染界面
func (h *Help) View() string {
	if !h.visible {
		return ""
	}

	var b strings.Builder

	// 标题
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.Primary).
		MarginBottom(1)
	title := titleStyle.Render("📖 快捷键帮助")
	b.WriteString(title)
	b.WriteString("\n\n")

	// 如果有分组，按分组显示
	if len(h.keyGroups) > 0 {
		for i, group := range h.keyGroups {
			if i > 0 {
				b.WriteString("\n")
			}
			b.WriteString(h.renderKeyGroup(group))
		}
	} else {
		// 否则直接显示列表
		b.WriteString(h.renderKeyList(h.keys))
	}

	b.WriteString("\n")
	footerStyle := lipgloss.NewStyle().
		Foreground(styles.Overlay0).
		Italic(true)
	b.WriteString(footerStyle.Render("按 ? 或 Esc 关闭"))

	// 对话框样式
	dialogStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Primary).
		Padding(1, 3).
		Width(50)

	return dialogStyle.Render(b.String())
}

// renderKeyGroup 渲染快捷键分组
func (h *Help) renderKeyGroup(group KeyGroup) string {
	var b strings.Builder

	// 分组标题
	groupTitleStyle := lipgloss.NewStyle().
		Foreground(styles.Subtext1).
		Bold(true)
	b.WriteString(groupTitleStyle.Render(group.Title))
	b.WriteString("\n")

	// 分隔线
	lineStyle := lipgloss.NewStyle().
		Foreground(styles.BorderSubtle)
	b.WriteString(lineStyle.Render(strings.Repeat("─", 15)))
	b.WriteString("\n")

	// 快捷键列表
	b.WriteString(h.renderKeyList(group.Keys))

	return b.String()
}

// renderKeyList 渲染快捷键列表
func (h *Help) renderKeyList(keys []key.Binding) string {
	var b strings.Builder

	keyStyle := lipgloss.NewStyle().
		Foreground(styles.Primary).
		Width(10)

	descStyle := lipgloss.NewStyle().
		Foreground(styles.Subtext0)

	for _, k := range keys {
		keyStr := keyStyle.Render(k.Help().Key)
		descStr := descStyle.Render(k.Help().Desc)
		b.WriteString(keyStr + descStr + "\n")
	}

	return b.String()
}

// RenderOverlay 渲染为浮动层
func (h *Help) RenderOverlay(base string) string {
	if !h.visible {
		return base
	}

	helpView := h.View()
	return PlaceOverlay(base, helpView, h.width, h.height, Center)
}

// ShortHelp 获取简短帮助
func (h *Help) ShortHelp() string {
	var parts []string
	for _, k := range h.keys {
		parts = append(parts, k.Help().Key+" "+k.Help().Desc)
	}
	return styles.HelpStyle.Render(strings.Join(parts, " | "))
}

// DefaultKeyGroups 默认快捷键分组
func DefaultKeyGroups() []KeyGroup {
	return []KeyGroup{
		{
			Title: "导航",
			Keys: []key.Binding{
				key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "向上")),
				key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "向下")),
				key.NewBinding(key.WithKeys("enter"), key.WithHelp("Enter", "确认")),
				key.NewBinding(key.WithKeys("esc"), key.WithHelp("Esc", "返回")),
			},
		},
		{
			Title: "操作",
			Keys: []key.Binding{
				key.NewBinding(key.WithKeys("c", "n"), key.WithHelp("c/n", "新建")),
				key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "删除")),
				key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "搜索")),
			},
		},
		{
			Title: "全局",
			Keys: []key.Binding{
				key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "帮助")),
				key.NewBinding(key.WithKeys("q"), key.WithHelp("q", "退出")),
			},
		},
	}
}
