package memory

import (
	"fmt"
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
	loadMemoryMsg struct {
		title    string
		content  string
		category string
		tags     []string
		priority int
		global   bool
		err      error
	}
	updateSuccessMsg struct{}
	updateErrorMsg   struct{ err error }
)

type EditPage struct {
	bs       *startup.Bootstrap
	frame    *layout.Frame
	width    int
	height   int
	pop      func(core.PageID) tea.Cmd
	memoryID int64
	focusIdx int
	loading  bool
	saving   bool
	err      error

	// 表单字段
	inputTitle     *components.Input
	textContent    *components.TextArea
	inputCategory  *components.Input
	inputTags      *components.Input
	selectPriority *components.Select
	selectGlobal   *components.Select
}

func NewEditPage(bs *startup.Bootstrap, memoryID int64, pop func(core.PageID) tea.Cmd) *EditPage {
	return &EditPage{
		bs:       bs,
		frame:    layout.NewFrame(80, 24),
		width:    80,
		height:   24,
		pop:      pop,
		memoryID: memoryID,
		loading:  true,

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

func (p *EditPage) Init() tea.Cmd {
	return p.loadMemory()
}

func (p *EditPage) Resize(w, h int) {
	p.width, p.height = w, h
	p.frame.Resize(w, h)
}

func (p *EditPage) Update(msg tea.Msg) (core.Page, tea.Cmd) {
	switch v := msg.(type) {
	case tea.KeyMsg:
		if p.loading || p.saving {
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

	case loadMemoryMsg:
		p.loading = false
		if v.err != nil {
			p.err = v.err
		} else {
			// 填充表单
			p.inputTitle.SetValue(v.title)
			p.textContent.SetValue(v.content)
			p.inputCategory.SetValue(v.category)
			if len(v.tags) > 0 {
				p.inputTags.SetValue(strings.Join(v.tags, ", "))
			}
			p.selectPriority.SetSelectedIndex(v.priority - 1)
			if v.global {
				p.selectGlobal.SetSelectedIndex(1)
			} else {
				p.selectGlobal.SetSelectedIndex(0)
			}
			return p, p.inputTitle.Focus()
		}

	case updateSuccessMsg:
		p.saving = false
		// 保存成功后返回列表
		return p, tea.Quit

	case updateErrorMsg:
		p.saving = false
		p.err = v.err
	}

	// 更新当前聚焦的字段
	return p, p.updateFocusedField(msg)
}

func (p *EditPage) View() string {
	cw, _ := p.frame.ContentSize()
	cardWidth := layout.FitCardWidth(cw)

	if p.loading {
		return components.LoadingState(theme.IconMemory+" 编辑记忆", "正在加载...", cardWidth)
	}

	if p.saving {
		return components.LoadingState(theme.IconMemory+" 编辑记忆", "正在保存...", cardWidth)
	}

	// 加载失败
	if p.err != nil && !p.loading && !p.saving {
		return components.ErrorState(theme.IconMemory+" 编辑记忆", p.err.Error(), cardWidth)
	}

	// 设置所有组件宽度
	formWidth := cardWidth - 8
	p.inputTitle.SetWidth(formWidth)
	p.textContent.SetWidth(formWidth)
	p.inputCategory.SetWidth(formWidth)
	p.inputTags.SetWidth(formWidth)
	p.selectPriority.SetWidth(formWidth)
	p.selectGlobal.SetWidth(formWidth)

	// 表单内容
	var formParts []string
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
	return components.Card(theme.IconMemory+" 编辑记忆", body, cardWidth)
}

func (p *EditPage) Meta() core.Meta {
	return core.Meta{
		Title:      "编辑记忆",
		Breadcrumb: "记忆管理 > 编辑",
		Extra:      fmt.Sprintf("ID: %d", p.memoryID),
		Keys: []components.KeyHint{
			{Key: "Ctrl+S", Desc: "保存"},
			{Key: "Tab", Desc: "下一项"},
			{Key: "Shift+Tab", Desc: "上一项"},
			{Key: "Esc", Desc: "取消"},
		},
	}
}

// nextField 切换到下一个字段
func (p *EditPage) nextField() tea.Cmd {
	p.blurAll()
	p.focusIdx = (p.focusIdx + 1) % 6
	return p.focusCurrent()
}

// prevField 切换到上一个字段
func (p *EditPage) prevField() tea.Cmd {
	p.blurAll()
	p.focusIdx = (p.focusIdx - 1 + 6) % 6
	return p.focusCurrent()
}

// blurAll 取消所有字段焦点
func (p *EditPage) blurAll() {
	p.inputTitle.Blur()
	p.textContent.Blur()
	p.inputCategory.Blur()
	p.inputTags.Blur()
	p.selectPriority.Blur()
	p.selectGlobal.Blur()
}

// focusCurrent 聚焦当前字段
func (p *EditPage) focusCurrent() tea.Cmd {
	switch p.focusIdx {
	case 0:
		return p.inputTitle.Focus()
	case 1:
		return p.textContent.Focus()
	case 2:
		return p.inputCategory.Focus()
	case 3:
		return p.inputTags.Focus()
	case 4:
		return p.selectPriority.Focus()
	case 5:
		return p.selectGlobal.Focus()
	}
	return nil
}

// updateFocusedField 更新当前聚焦的字段
func (p *EditPage) updateFocusedField(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch p.focusIdx {
	case 0:
		_, cmd = p.inputTitle.Update(msg)
	case 1:
		_, cmd = p.textContent.Update(msg)
	case 2:
		_, cmd = p.inputCategory.Update(msg)
	case 3:
		_, cmd = p.inputTags.Update(msg)
	case 4:
		_, cmd = p.selectPriority.Update(msg)
	case 5:
		_, cmd = p.selectGlobal.Update(msg)
	}
	return cmd
}

// loadMemory 加载记忆数据
func (p *EditPage) loadMemory() tea.Cmd {
	return func() tea.Msg {
		ctx := p.bs.Context()
		memory, err := p.bs.MemoryService.GetMemoryByID(ctx, p.memoryID)
		if err != nil {
			return loadMemoryMsg{err: err}
		}

		return loadMemoryMsg{
			title:    memory.Title,
			content:  memory.Content,
			category: memory.Category,
			tags:     memory.GetTagStrings(),
			priority: memory.Priority,
			global:   memory.Global,
		}
	}
}

// save 保存记忆
func (p *EditPage) save() tea.Cmd {
	// 验证表单
	if err := p.inputTitle.Validate(); err != nil {
		p.inputTitle.SetError(err)
		return nil
	}
	if err := p.textContent.Validate(); err != nil {
		p.textContent.SetError(err)
		return nil
	}

	// 清除错误
	p.inputTitle.SetError(nil)
	p.textContent.SetError(nil)
	p.err = nil

	// 解析标签
	tags := parseTags(p.inputTags.Value())

	// 获取优先级
	priority := p.selectPriority.Value().(int)

	// 分类
	category := strings.TrimSpace(p.inputCategory.Value())
	if category == "" {
		category = "默认"
	}

	// 准备更新数据
	title := p.inputTitle.Value()
	content := p.textContent.Value()

	p.saving = true
	return func() tea.Msg {
		ctx := p.bs.Context()

		// 先通过ID获取当前记忆的Code
		memory, err := p.bs.MemoryService.GetMemoryByID(ctx, p.memoryID)
		if err != nil {
			return updateErrorMsg{err: err}
		}

		input := &dto.MemoryUpdateDTO{
			Code:     memory.Code,
			Title:    &title,
			Content:  &content,
			Category: &category,
			Tags:     &tags,
			Priority: &priority,
		}

		if err := p.bs.MemoryService.UpdateMemory(ctx, input); err != nil {
			return updateErrorMsg{err: err}
		}

		return updateSuccessMsg{}
	}
}
