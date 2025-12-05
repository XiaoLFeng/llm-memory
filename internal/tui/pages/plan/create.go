package plan

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
)

// CreatePage 计划创建页面
type CreatePage struct {
	bs       *startup.Bootstrap
	frame    *layout.Frame
	pop      func(core.PageID) tea.Cmd
	width    int
	height   int
	focusIdx int // 当前聚焦的表单字段索引

	// 表单字段
	codeInput       *components.Input
	titleInput      *components.Input
	descriptionArea *components.TextArea
	contentArea     *components.TextArea
	scopeSelect     *components.Select

	// 状态
	submitting bool
	err        error
}

type createResultMsg struct {
	success bool
	err     error
}

// NewCreatePage 创建计划创建页面
func NewCreatePage(bs *startup.Bootstrap, pop func(core.PageID) tea.Cmd) *CreatePage {
	// 初始化表单组件
	codeInput := components.NewInput("标识码", "小写字母+连字符，如: my-plan", true)
	titleInput := components.NewInput("标题", "请输入计划标题", true)
	descriptionArea := components.NewTextArea("描述", "请输入计划描述（摘要）", true)
	contentArea := components.NewTextArea("详细内容", "请输入详细的计划内容，支持 Markdown 格式", true)
	scopeSelect := components.NewSelect("作用域", []components.SelectOption{
		{Label: "项目/组内", Value: false},
		{Label: "全局可见", Value: true},
	})

	// 设置文本域高度
	descriptionArea.SetHeight(3)
	contentArea.SetHeight(8)

	return &CreatePage{
		bs:              bs,
		frame:           layout.NewFrame(80, 24),
		pop:             pop,
		width:           80,
		height:          24,
		codeInput:       codeInput,
		titleInput:      titleInput,
		descriptionArea: descriptionArea,
		contentArea:     contentArea,
		scopeSelect:     scopeSelect,
	}
}

func (p *CreatePage) Init() tea.Cmd {
	// 聚焦第一个字段
	return p.codeInput.Focus()
}

func (p *CreatePage) Resize(w, h int) {
	p.width, p.height = w, h
	p.frame.Resize(w, h)

	// 调整组件宽度
	cw, _ := p.frame.ContentSize()
	fieldWidth := cw - 10
	if fieldWidth > 80 {
		fieldWidth = 80
	}

	p.codeInput.SetWidth(fieldWidth)
	p.titleInput.SetWidth(fieldWidth)
	p.descriptionArea.SetWidth(fieldWidth)
	p.contentArea.SetWidth(fieldWidth)
	p.scopeSelect.SetWidth(fieldWidth)
}

func (p *CreatePage) Update(msg tea.Msg) (core.Page, tea.Cmd) {
	var cmd tea.Cmd

	switch v := msg.(type) {
	case tea.KeyMsg:
		// 提交中禁用输入
		if p.submitting {
			return p, nil
		}

		switch v.String() {
		case "esc":
			if p.pop != nil {
				return p, p.pop(core.PagePlan)
			}
		case "ctrl+s":
			return p, p.submit()
		case "tab", "down":
			return p, p.nextField()
		case "shift+tab", "up":
			return p, p.prevField()
		}

	case createResultMsg:
		p.submitting = false
		if v.success {
			// 创建成功，返回列表页
			if p.pop != nil {
				return p, p.pop(core.PagePlan)
			}
		} else {
			p.err = v.err
		}
		return p, nil
	}

	// 更新当前聚焦的组件
	switch p.focusIdx {
	case 0:
		p.codeInput, cmd = p.codeInput.Update(msg)
	case 1:
		p.titleInput, cmd = p.titleInput.Update(msg)
	case 2:
		p.descriptionArea, cmd = p.descriptionArea.Update(msg)
	case 3:
		p.contentArea, cmd = p.contentArea.Update(msg)
	case 4:
		p.scopeSelect, cmd = p.scopeSelect.Update(msg)
	}

	return p, cmd
}

func (p *CreatePage) View() string {
	cw, _ := p.frame.ContentSize()
	cardW := layout.FitCardWidth(cw)

	var body strings.Builder

	// 渲染表单
	body.WriteString(p.codeInput.View())
	body.WriteString("\n\n")
	body.WriteString(p.titleInput.View())
	body.WriteString("\n\n")
	body.WriteString(p.descriptionArea.View())
	body.WriteString("\n\n")
	body.WriteString(p.contentArea.View())
	body.WriteString("\n\n")
	body.WriteString(p.scopeSelect.View())
	body.WriteString("\n\n")

	// 提示信息
	if p.submitting {
		body.WriteString(theme.TextDim.Render("正在创建计划..."))
	} else if p.err != nil {
		body.WriteString(theme.FormError.Render(fmt.Sprintf("错误: %s", p.err.Error())))
	} else {
		body.WriteString(theme.FormHint.Render("Ctrl+S 保存 | Tab/Shift+Tab 切换字段 | Esc 取消"))
	}

	title := theme.IconCreate + " 创建计划"
	return components.Card(title, body.String(), cardW)
}

func (p *CreatePage) Meta() core.Meta {
	return core.Meta{
		Title:      "创建计划",
		Breadcrumb: "计划管理 > 创建",
		Extra:      "填写表单后按 Ctrl+S 保存",
		Keys: []components.KeyHint{
			{Key: "Ctrl+S", Desc: "保存"},
			{Key: "Tab", Desc: "下一字段"},
			{Key: "Shift+Tab", Desc: "上一字段"},
			{Key: "Esc", Desc: "取消"},
		},
	}
}

// nextField 切换到下一个字段
func (p *CreatePage) nextField() tea.Cmd {
	p.blurCurrent()
	p.focusIdx = (p.focusIdx + 1) % 5
	return p.focusCurrent()
}

// prevField 切换到上一个字段
func (p *CreatePage) prevField() tea.Cmd {
	p.blurCurrent()
	p.focusIdx = (p.focusIdx - 1 + 5) % 5
	return p.focusCurrent()
}

// focusCurrent 聚焦当前字段
func (p *CreatePage) focusCurrent() tea.Cmd {
	switch p.focusIdx {
	case 0:
		return p.codeInput.Focus()
	case 1:
		return p.titleInput.Focus()
	case 2:
		return p.descriptionArea.Focus()
	case 3:
		return p.contentArea.Focus()
	case 4:
		return p.scopeSelect.Focus()
	}
	return nil
}

// blurCurrent 取消当前字段聚焦
func (p *CreatePage) blurCurrent() {
	switch p.focusIdx {
	case 0:
		p.codeInput.Blur()
	case 1:
		p.titleInput.Blur()
	case 2:
		p.descriptionArea.Blur()
	case 3:
		p.contentArea.Blur()
	case 4:
		p.scopeSelect.Blur()
	}
}

// submit 提交表单
func (p *CreatePage) submit() tea.Cmd {
	// 验证表单
	if err := p.codeInput.Validate(); err != nil {
		p.err = err
		return nil
	}
	if err := p.titleInput.Validate(); err != nil {
		p.err = err
		return nil
	}
	if err := p.descriptionArea.Validate(); err != nil {
		p.err = err
		return nil
	}
	if err := p.contentArea.Validate(); err != nil {
		p.err = err
		return nil
	}

	p.submitting = true
	p.err = nil

	return func() tea.Msg {
		ctx := p.bs.Context()

		// 获取表单值
		code := strings.TrimSpace(p.codeInput.Value())
		title := strings.TrimSpace(p.titleInput.Value())
		description := strings.TrimSpace(p.descriptionArea.Value())
		content := strings.TrimSpace(p.contentArea.Value())

		// 创建计划
		input := &dto.PlanCreateDTO{
			Code:        code,
			Title:       title,
			Description: description,
			Content:     content,
		}

		_, err := p.bs.PlanService.CreatePlan(ctx, input, p.bs.CurrentScope)
		if err != nil {
			return createResultMsg{success: false, err: err}
		}

		return createResultMsg{success: true}
	}
}
