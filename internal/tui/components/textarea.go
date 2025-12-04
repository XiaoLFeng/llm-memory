package components

import (
	"github.com/XiaoLFeng/llm-memory/internal/tui/theme"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TextArea 多行文本输入组件
type TextArea struct {
	model    textarea.Model
	label    string
	required bool
	focused  bool
	err      error
	width    int
	height   int
}

// NewTextArea 创建多行输入框
func NewTextArea(label, placeholder string, required bool) *TextArea {
	ta := textarea.New()
	ta.Placeholder = placeholder
	ta.CharLimit = 4096
	ta.SetHeight(5)
	ta.ShowLineNumbers = false

	return &TextArea{
		model:    ta,
		label:    label,
		required: required,
		width:    60,
		height:   5,
	}
}

// Init 初始化
func (t *TextArea) Init() tea.Cmd {
	return nil
}

// Update 更新
func (t *TextArea) Update(msg tea.Msg) (*TextArea, tea.Cmd) {
	var cmd tea.Cmd
	t.model, cmd = t.model.Update(msg)
	return t, cmd
}

// View 渲染
func (t *TextArea) View() string {
	// 标签
	labelStr := theme.FormLabel.Render(t.label)
	if t.required {
		labelStr += theme.FormLabelRequired.Render()
	}

	// 输入框样式
	inputStyle := theme.FormInput
	if t.focused {
		inputStyle = theme.FormInputFocused
	}

	// 设置文本域尺寸
	t.model.SetWidth(t.width - 4)
	t.model.SetHeight(t.height)

	// 渲染文本域
	inputBox := inputStyle.Width(t.width).Render(t.model.View())

	// 错误提示
	errStr := ""
	if t.err != nil {
		errStr = theme.FormError.Render(t.err.Error())
	}

	return lipgloss.JoinVertical(lipgloss.Left, labelStr, inputBox, errStr)
}

// Focus 获取焦点
func (t *TextArea) Focus() tea.Cmd {
	t.focused = true
	return t.model.Focus()
}

// Blur 失去焦点
func (t *TextArea) Blur() {
	t.focused = false
	t.model.Blur()
}

// Value 获取值
func (t *TextArea) Value() string {
	return t.model.Value()
}

// SetValue 设置值
func (t *TextArea) SetValue(v string) {
	t.model.SetValue(v)
}

// SetWidth 设置宽度
func (t *TextArea) SetWidth(w int) {
	t.width = w
	t.model.SetWidth(w - 4)
}

// SetHeight 设置高度
func (t *TextArea) SetHeight(h int) {
	t.height = h
	t.model.SetHeight(h)
}

// SetError 设置错误
func (t *TextArea) SetError(err error) {
	t.err = err
}

// Validate 验证
func (t *TextArea) Validate() error {
	if t.required && t.model.Value() == "" {
		return &ValidationError{Field: t.label, Message: "不能为空"}
	}
	return nil
}

// IsFocused 是否聚焦
func (t *TextArea) IsFocused() bool {
	return t.focused
}
