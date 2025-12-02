package components

import (
	"strings"

	"github.com/XiaoLFeng/llm-memory/internal/tui/styles"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// Help 帮助面板组件
// 呀~ 显示所有可用的快捷键！❓
type Help struct {
	keys    []key.Binding
	visible bool
	width   int
}

// NewHelp 创建帮助面板组件
func NewHelp() *Help {
	return &Help{
		width: 60,
	}
}

// SetKeys 设置快捷键
func (h *Help) SetKeys(keys []key.Binding) {
	h.keys = keys
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
		if h.visible && msg.String() == "?" || msg.String() == "esc" {
			h.Hide()
		}
	case tea.WindowSizeMsg:
		h.width = msg.Width
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
	title := styles.DialogTitleStyle.Render("📖 快捷键帮助")
	b.WriteString(title)
	b.WriteString("\n\n")

	// 快捷键列表
	for _, k := range h.keys {
		keyStr := styles.StatusKeyStyle.Render(k.Help().Key)
		descStr := styles.MutedStyle.Render(k.Help().Desc)
		b.WriteString("  " + keyStr + "  " + descStr + "\n")
	}

	b.WriteString("\n")
	b.WriteString(styles.MutedStyle.Render("按 ? 或 ESC 关闭"))

	return styles.DialogStyle.Render(b.String())
}

// ShortHelp 获取简短帮助
func (h *Help) ShortHelp() string {
	var parts []string
	for _, k := range h.keys {
		parts = append(parts, k.Help().Key+" "+k.Help().Desc)
	}
	return styles.HelpStyle.Render(strings.Join(parts, " | "))
}
