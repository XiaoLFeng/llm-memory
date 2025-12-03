package components

import (
	"github.com/XiaoLFeng/llm-memory/internal/tui/common"
	"github.com/XiaoLFeng/llm-memory/internal/tui/styles"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Confirm 确认对话框组件
type Confirm struct {
	title     string
	message   string
	visible   bool
	selected  int // 0: 取消, 1: 确认
	onConfirm tea.Cmd
	onCancel  tea.Cmd
	width     int
	height    int
}

// NewConfirm 创建确认对话框组件
func NewConfirm() *Confirm {
	return &Confirm{
		width:  80,
		height: 24,
	}
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

// SetSize 设置窗口大小
func (c *Confirm) SetSize(width, height int) {
	c.width = width
	c.height = height
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
				return common.ConfirmResultMsg{Confirmed: c.selected == 1}
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("esc", "n"))):
			c.visible = false
			return c, func() tea.Msg {
				return common.ConfirmResultMsg{Confirmed: false}
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("y"))):
			c.visible = false
			return c, func() tea.Msg {
				return common.ConfirmResultMsg{Confirmed: true}
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

	// 图标
	iconStyle := lipgloss.NewStyle().
		Foreground(styles.Warning).
		Bold(true)
	icon := iconStyle.Render("!")

	// 标题样式
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.Text).
		MarginBottom(1)
	title := titleStyle.Render(icon + " " + c.title)

	// 消息样式
	messageStyle := lipgloss.NewStyle().
		Foreground(styles.Subtext0).
		MarginBottom(1)
	message := messageStyle.Render(c.message)

	// 按钮样式
	btnNormalStyle := lipgloss.NewStyle().
		Foreground(styles.Text).
		Background(styles.Surface0).
		Padding(0, 3).
		MarginRight(2)

	btnActiveStyle := lipgloss.NewStyle().
		Foreground(styles.Text).
		Background(styles.Primary).
		Padding(0, 3).
		MarginRight(2)

	btnDangerStyle := lipgloss.NewStyle().
		Foreground(styles.Text).
		Background(styles.Error).
		Padding(0, 3)

	// 渲染按钮
	var cancelBtn, confirmBtn string
	if c.selected == 0 {
		cancelBtn = btnActiveStyle.Render("取消")
		confirmBtn = btnDangerStyle.Copy().Background(styles.Surface1).Render("确认")
	} else {
		cancelBtn = btnNormalStyle.Render("取消")
		confirmBtn = btnDangerStyle.Render("确认")
	}

	// 按钮行居中
	buttonsStyle := lipgloss.NewStyle().
		MarginTop(1).
		Align(lipgloss.Center)
	buttons := buttonsStyle.Render(cancelBtn + confirmBtn)

	// 组合对话框内容
	contentStyle := lipgloss.NewStyle().
		Align(lipgloss.Center)
	content := contentStyle.Render(title + "\n\n" + message + "\n" + buttons)

	// 对话框外框
	dialogStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Warning).
		Padding(1, 3).
		Width(50)

	return dialogStyle.Render(content)
}

// RenderOverlay 渲染为浮动层（用于居中显示）
func (c *Confirm) RenderOverlay(base string) string {
	if !c.visible {
		return base
	}

	confirmView := c.View()
	return PlaceOverlay(base, confirmView, c.width, c.height, Center)
}
