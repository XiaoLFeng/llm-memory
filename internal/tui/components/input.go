package components

import (
	"github.com/XiaoLFeng/llm-memory/internal/tui/theme"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Input 单行文本输入组件
type Input struct {
	model    textinput.Model
	label    string
	required bool
	focused  bool
	err      error
	width    int
}

// NewInput 创建输入框
func NewInput(label, placeholder string, required bool) *Input {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.CharLimit = 256

	return &Input{
		model:    ti,
		label:    label,
		required: required,
		width:    40,
	}
}

// Init 初始化
func (i *Input) Init() tea.Cmd {
	return nil
}

// Update 更新
func (i *Input) Update(msg tea.Msg) (*Input, tea.Cmd) {
	var cmd tea.Cmd
	i.model, cmd = i.model.Update(msg)
	return i, cmd
}

// View 渲染
func (i *Input) View() string {
	// 标签
	labelStr := theme.FormLabel.Render(i.label)
	if i.required {
		labelStr += theme.FormLabelRequired.Render()
	}

	// 输入框样式
	inputStyle := theme.FormInput
	if i.focused {
		inputStyle = theme.FormInputFocused
	}

	// 设置输入框宽度
	i.model.Width = i.width - 4

	// 渲染输入框
	inputBox := inputStyle.Width(i.width).Render(i.model.View())

	// 错误提示
	errStr := ""
	if i.err != nil {
		errStr = theme.FormError.Render(i.err.Error())
	}

	return lipgloss.JoinVertical(lipgloss.Left, labelStr, inputBox, errStr)
}

// Focus 获取焦点
func (i *Input) Focus() tea.Cmd {
	i.focused = true
	return i.model.Focus()
}

// Blur 失去焦点
func (i *Input) Blur() {
	i.focused = false
	i.model.Blur()
}

// Value 获取值
func (i *Input) Value() string {
	return i.model.Value()
}

// SetValue 设置值
func (i *Input) SetValue(v string) {
	i.model.SetValue(v)
}

// SetWidth 设置宽度
func (i *Input) SetWidth(w int) {
	i.width = w
	i.model.Width = w - 4
}

// SetError 设置错误
func (i *Input) SetError(err error) {
	i.err = err
}

// Validate 验证
func (i *Input) Validate() error {
	if i.required && i.model.Value() == "" {
		return &ValidationError{Field: i.label, Message: "不能为空"}
	}
	return nil
}

// IsFocused 是否聚焦
func (i *Input) IsFocused() bool {
	return i.focused
}

// ValidationError 验证错误
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Field + e.Message
}
