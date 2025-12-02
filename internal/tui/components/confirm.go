package components

import (
	"github.com/XiaoLFeng/llm-memory/internal/tui/styles"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// Confirm 确认对话框组件
// 呀~ 用于危险操作前的确认！⚠️
type Confirm struct {
	title     string
	message   string
	visible   bool
	selected  int // 0: 取消, 1: 确认
	onConfirm tea.Cmd
	onCancel  tea.Cmd
}

// NewConfirm 创建确认对话框组件
func NewConfirm() *Confirm {
	return &Confirm{}
}

// Show 显示确认对话框
func (c *Confirm) Show(title, message string, onConfirm, onCancel tea.Cmd) {
	c.title = title
	c.message = message
	c.visible = true
	c.selected = 0 // 默认选中取消
	c.onConfirm = onConfirm
	c.onCancel = onCancel
}

// Hide 隐藏确认对话框
func (c *Confirm) Hide() {
	c.visible = false
}

// IsVisible 是否可见
func (c *Confirm) IsVisible() bool {
	return c.visible
}

// GetOnConfirm 获取确认回调
func (c *Confirm) GetOnConfirm() tea.Cmd {
	return c.onConfirm
}

// GetOnCancel 获取取消回调
func (c *Confirm) GetOnCancel() tea.Cmd {
	return c.onCancel
}

// Init 初始化
func (c *Confirm) Init() tea.Cmd {
	return nil
}

// confirmResultMsg 确认结果消息
type confirmResultMsg struct {
	confirmed bool
}

// Update 处理输入
func (c *Confirm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !c.visible {
		return c, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("left", "h"))):
			c.selected = 0
		case key.Matches(msg, key.NewBinding(key.WithKeys("right", "l"))):
			c.selected = 1
		case key.Matches(msg, key.NewBinding(key.WithKeys("tab"))):
			c.selected = (c.selected + 1) % 2
		case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
			c.visible = false
			return c, func() tea.Msg {
				return confirmResultMsg{confirmed: c.selected == 1}
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("esc", "n"))):
			c.visible = false
			return c, func() tea.Msg {
				return confirmResultMsg{confirmed: false}
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("y"))):
			c.visible = false
			return c, func() tea.Msg {
				return confirmResultMsg{confirmed: true}
			}
		}
	}

	return c, nil
}

// View 渲染界面
func (c *Confirm) View() string {
	if !c.visible {
		return ""
	}

	// 标题
	title := styles.DialogTitleStyle.Render(c.title)

	// 消息
	message := styles.MutedStyle.Render(c.message)

	// 按钮
	cancelBtn := styles.DialogButtonStyle.Render("取消")
	confirmBtn := styles.DialogButtonStyle.Render("确认")

	if c.selected == 0 {
		cancelBtn = styles.DialogButtonActiveStyle.Render("取消")
	} else {
		confirmBtn = styles.DialogButtonActiveStyle.Render("确认")
	}

	buttons := cancelBtn + " " + confirmBtn

	// 组合对话框
	content := title + "\n\n" + message + "\n\n" + buttons

	return styles.DialogStyle.Render(content)
}
