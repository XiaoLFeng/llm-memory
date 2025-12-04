package components

import (
	"github.com/XiaoLFeng/llm-memory/internal/tui/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SelectOption 选择器选项
type SelectOption struct {
	Label string
	Value interface{}
}

// Select 选择器组件
type Select struct {
	label   string
	options []SelectOption
	cursor  int
	focused bool
	width   int
}

// NewSelect 创建选择器
func NewSelect(label string, options []SelectOption) *Select {
	return &Select{
		label:   label,
		options: options,
		cursor:  0,
		width:   40,
	}
}

// Init 初始化
func (s *Select) Init() tea.Cmd {
	return nil
}

// Update 更新
func (s *Select) Update(msg tea.Msg) (*Select, tea.Cmd) {
	if !s.focused {
		return s, nil
	}

	switch v := msg.(type) {
	case tea.KeyMsg:
		switch v.String() {
		case "left", "h":
			if s.cursor > 0 {
				s.cursor--
			}
		case "right", "l":
			if s.cursor < len(s.options)-1 {
				s.cursor++
			}
		}
	}
	return s, nil
}

// View 渲染
func (s *Select) View() string {
	// 标签
	labelStr := theme.FormLabel.Render(s.label)

	// 选项
	var optionViews []string
	for i, opt := range s.options {
		style := theme.SelectOption
		prefix := "  "
		if i == s.cursor {
			if s.focused {
				style = theme.SelectOptionSelected
				prefix = theme.SelectCursor.Render()
			} else {
				style = theme.SelectOptionSelected.Copy().Background(theme.Surface0)
			}
		}
		optionViews = append(optionViews, prefix+style.Render(opt.Label))
	}

	optionsBox := lipgloss.JoinHorizontal(lipgloss.Left, optionViews...)

	// 边框
	boxStyle := theme.FormInput
	if s.focused {
		boxStyle = theme.FormInputFocused
	}
	box := boxStyle.Width(s.width).Render(optionsBox)

	return lipgloss.JoinVertical(lipgloss.Left, labelStr, box)
}

// Focus 获取焦点
func (s *Select) Focus() tea.Cmd {
	s.focused = true
	return nil
}

// Blur 失去焦点
func (s *Select) Blur() {
	s.focused = false
}

// Value 获取当前选中值
func (s *Select) Value() interface{} {
	if s.cursor < len(s.options) {
		return s.options[s.cursor].Value
	}
	return nil
}

// SelectedIndex 获取选中索引
func (s *Select) SelectedIndex() int {
	return s.cursor
}

// SetSelectedIndex 设置选中索引
func (s *Select) SetSelectedIndex(idx int) {
	if idx >= 0 && idx < len(s.options) {
		s.cursor = idx
	}
}

// SetWidth 设置宽度
func (s *Select) SetWidth(w int) {
	s.width = w
}

// IsFocused 是否聚焦
func (s *Select) IsFocused() bool {
	return s.focused
}

// SelectedLabel 获取选中标签
func (s *Select) SelectedLabel() string {
	if s.cursor < len(s.options) {
		return s.options[s.cursor].Label
	}
	return ""
}
