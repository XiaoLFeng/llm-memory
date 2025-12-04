package memory

import (
	"strings"

	"github.com/XiaoLFeng/llm-memory/internal/models/dto"
	"github.com/XiaoLFeng/llm-memory/internal/tui/components"
	"github.com/XiaoLFeng/llm-memory/internal/tui/core"
	"github.com/XiaoLFeng/llm-memory/internal/tui/layout"
	"github.com/XiaoLFeng/llm-memory/internal/tui/theme"
	"github.com/XiaoLFeng/llm-memory/startup"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type (
	createSuccessMsg struct{}
	createErrorMsg   struct{ err error }
)

type CreatePage struct {
	bs       *startup.Bootstrap
	frame    *layout.Frame
	width    int
	height   int
	pop      func(core.PageID) tea.Cmd
	focusIdx int
	saving   bool
	err      error

	// 表单字段
	inputCode      *components.Input
	inputTitle     *components.Input
	textContent    *components.TextArea
	inputCategory  *components.Input
	inputTags      *components.Input
	selectPriority *components.Select
	selectGlobal   *components.Select
}

func NewCreatePage(bs *startup.Bootstrap, pop func(core.PageID) tea.Cmd) *CreatePage {
	return &CreatePage{
		bs:     bs,
		frame:  layout.NewFrame(80, 24),
		width:  80,
		height: 24,
		pop:    pop,

		inputCode:     components.NewInput("标识码", "小写字母+连字符，如: my-memory", true),
		inputTitle:    components.NewInput("标题", "请输入记忆标题", true),
		textContent:   components.NewTextArea("内容", "请输入记忆内容", true),
		inputCategory: components.NewInput("分类", "默认", false),
		inputTags:     components.NewInput("标签", "多个标签用逗号分隔", false),
		selectPriority: components.NewSelect("优先级", []components.SelectOption{
			{Label: "1-低", Value: 1},
			{Label: "2-中", Value: 2},
			{Label: "3-高", Value: 3},
			{Label: "4-紧急", Value: 4},
		}),
		selectGlobal: components.NewSelect("作用域", []components.SelectOption{
			{Label: "私有", Value: false},
			{Label: "全局", Value: true},
		}),
	}
}

func (p *CreatePage) Init() tea.Cmd {
	p.inputCategory.SetValue("默认")
	p.selectPriority.SetSelectedIndex(1) // 默认中优先级
	return p.inputCode.Focus()
}

func (p *CreatePage) Resize(w, h int) {
	p.width, p.height = w, h
	p.frame.Resize(w, h)
}

func (p *CreatePage) Update(msg tea.Msg) (core.Page, tea.Cmd) {
	switch v := msg.(type) {
	case tea.KeyMsg:
		if p.saving {
			return p, nil
		}

		switch v.String() {
		case "ctrl+s":
			return p, p.save()
		case "esc":
			return p, tea.Quit
		case "tab", "down":
			return p, p.nextField()
		case "shift+tab", "up":
			return p, p.prevField()
		}

	case createSuccessMsg:
		p.saving = false
		// 保存成功后返回列表
		return p, tea.Quit
	case createErrorMsg:
		p.saving = false
		p.err = v.err
	}

	// 更新当前聚焦的字段
	return p, p.updateFocusedField(msg)
}

func (p *CreatePage) View() string {
	cw, _ := p.frame.ContentSize()
	cardWidth := layout.FitCardWidth(cw)

	if p.saving {
		return components.LoadingState(theme.IconMemory+" 创建记忆", "正在保存...", cardWidth)
	}

	// 设置所有组件宽度
	formWidth := cardWidth - 8
	p.inputCode.SetWidth(formWidth)
	p.inputTitle.SetWidth(formWidth)
	p.textContent.SetWidth(formWidth)
	p.inputCategory.SetWidth(formWidth)
	p.inputTags.SetWidth(formWidth)
	p.selectPriority.SetWidth(formWidth)
	p.selectGlobal.SetWidth(formWidth)

	// 表单内容
	var formParts []string
	formParts = append(formParts, p.inputCode.View())
	formParts = append(formParts, p.inputTitle.View())
	formParts = append(formParts, p.textContent.View())
	formParts = append(formParts, p.inputCategory.View())
	formParts = append(formParts, p.inputTags.View())
	formParts = append(formParts, p.selectPriority.View())
	formParts = append(formParts, p.selectGlobal.View())

	// 错误提示
	if p.err != nil {
		errMsg := theme.FormError.Render("错误: " + p.err.Error())
		formParts = append(formParts, errMsg)
	}

	// 提示信息
	hint := theme.FormHint.Render("Ctrl+S 保存 | Tab/↓ 下一项 | Shift+Tab/↑ 上一项 | Esc 取消")
	formParts = append(formParts, hint)

	body := lipgloss.JoinVertical(lipgloss.Left, formParts...)
	return components.Card(theme.IconMemory+" 创建记忆", body, cardWidth)
}

func (p *CreatePage) Meta() core.Meta {
	return core.Meta{
		Title:      "创建记忆",
		Breadcrumb: "记忆管理 > 创建",
		Extra:      "",
		Keys: []components.KeyHint{
			{Key: "Ctrl+S", Desc: "保存"},
			{Key: "Tab", Desc: "下一项"},
			{Key: "Shift+Tab", Desc: "上一项"},
			{Key: "Esc", Desc: "取消"},
		},
	}
}

// nextField 切换到下一个字段
func (p *CreatePage) nextField() tea.Cmd {
	p.blurAll()
	p.focusIdx = (p.focusIdx + 1) % 7
	return p.focusCurrent()
}

// prevField 切换到上一个字段
func (p *CreatePage) prevField() tea.Cmd {
	p.blurAll()
	p.focusIdx = (p.focusIdx - 1 + 7) % 7
	return p.focusCurrent()
}

// blurAll 取消所有字段焦点
func (p *CreatePage) blurAll() {
	p.inputCode.Blur()
	p.inputTitle.Blur()
	p.textContent.Blur()
	p.inputCategory.Blur()
	p.inputTags.Blur()
	p.selectPriority.Blur()
	p.selectGlobal.Blur()
}

// focusCurrent 聚焦当前字段
func (p *CreatePage) focusCurrent() tea.Cmd {
	switch p.focusIdx {
	case 0:
		return p.inputCode.Focus()
	case 1:
		return p.inputTitle.Focus()
	case 2:
		return p.textContent.Focus()
	case 3:
		return p.inputCategory.Focus()
	case 4:
		return p.inputTags.Focus()
	case 5:
		return p.selectPriority.Focus()
	case 6:
		return p.selectGlobal.Focus()
	}
	return nil
}

// updateFocusedField 更新当前聚焦的字段
func (p *CreatePage) updateFocusedField(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch p.focusIdx {
	case 0:
		_, cmd = p.inputCode.Update(msg)
	case 1:
		_, cmd = p.inputTitle.Update(msg)
	case 2:
		_, cmd = p.textContent.Update(msg)
	case 3:
		_, cmd = p.inputCategory.Update(msg)
	case 4:
		_, cmd = p.inputTags.Update(msg)
	case 5:
		_, cmd = p.selectPriority.Update(msg)
	case 6:
		_, cmd = p.selectGlobal.Update(msg)
	}
	return cmd
}

// save 保存记忆
func (p *CreatePage) save() tea.Cmd {
	// 验证表单
	if err := p.inputCode.Validate(); err != nil {
		p.inputCode.SetError(err)
		return nil
	}
	if err := p.inputTitle.Validate(); err != nil {
		p.inputTitle.SetError(err)
		return nil
	}
	if err := p.textContent.Validate(); err != nil {
		p.textContent.SetError(err)
		return nil
	}

	// 清除错误
	p.inputCode.SetError(nil)
	p.inputTitle.SetError(nil)
	p.textContent.SetError(nil)
	p.err = nil

	// 解析标签
	tags := parseTags(p.inputTags.Value())

	// 获取优先级
	priority := p.selectPriority.Value().(int)

	// 获取全局标志
	global := p.selectGlobal.Value().(bool)

	// 分类默认值
	category := strings.TrimSpace(p.inputCategory.Value())
	if category == "" {
		category = "默认"
	}

	p.saving = true
	return func() tea.Msg {
		ctx := p.bs.Context()
		input := &dto.MemoryCreateDTO{
			Code:     p.inputCode.Value(),
			Title:    p.inputTitle.Value(),
			Content:  p.textContent.Value(),
			Category: category,
			Tags:     tags,
			Priority: priority,
			Global:   global,
		}

		if _, err := p.bs.MemoryService.CreateMemory(ctx, input, p.bs.CurrentScope); err != nil {
			return createErrorMsg{err: err}
		}

		return createSuccessMsg{}
	}
}

// parseTags 解析标签字符串
func parseTags(tagsStr string) []string {
	if tagsStr == "" {
		return nil
	}
	parts := strings.Split(tagsStr, ",")
	var tags []string
	for _, part := range parts {
		if tag := strings.TrimSpace(part); tag != "" {
			tags = append(tags, tag)
		}
	}
	return tags
}
