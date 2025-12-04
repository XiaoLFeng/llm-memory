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

// EditPage 计划编辑页面
type EditPage struct {
	bs       *startup.Bootstrap
	frame    *layout.Frame
	pop      func(core.PageID) tea.Cmd
	width    int
	height   int
	focusIdx int // 当前聚焦的表单字段索引
	planID   int64

	// 表单字段
	titleInput      *components.Input
	descriptionArea *components.TextArea
	contentArea     *components.TextArea
	progressInput   *components.Input

	// 状态
	loading    bool
	submitting bool
	err        error
}

type editLoadMsg struct {
	title       string
	description string
	content     string
	progress    int
	err         error
}

type editResultMsg struct {
	success bool
	err     error
}

// NewEditPage 创建计划编辑页面
func NewEditPage(bs *startup.Bootstrap, planID int64, pop func(core.PageID) tea.Cmd) *EditPage {
	// 初始化表单组件
	titleInput := components.NewInput("标题", "请输入计划标题", true)
	descriptionArea := components.NewTextArea("描述", "请输入计划描述（摘要）", true)
	contentArea := components.NewTextArea("详细内容", "请输入详细的计划内容，支持 Markdown 格式", true)
	progressInput := components.NewInput("进度 (0-100)", "输入 0-100 的数字", false)

	// 设置文本域高度
	descriptionArea.SetHeight(3)
	contentArea.SetHeight(8)

	return &EditPage{
		bs:              bs,
		frame:           layout.NewFrame(80, 24),
		pop:             pop,
		width:           80,
		height:          24,
		planID:          planID,
		titleInput:      titleInput,
		descriptionArea: descriptionArea,
		contentArea:     contentArea,
		progressInput:   progressInput,
		loading:         true,
	}
}

func (p *EditPage) Init() tea.Cmd {
	return p.loadPlan()
}

func (p *EditPage) Resize(w, h int) {
	p.width, p.height = w, h
	p.frame.Resize(w, h)

	// 调整组件宽度
	cw, _ := p.frame.ContentSize()
	fieldWidth := cw - 10
	if fieldWidth > 80 {
		fieldWidth = 80
	}

	p.titleInput.SetWidth(fieldWidth)
	p.descriptionArea.SetWidth(fieldWidth)
	p.contentArea.SetWidth(fieldWidth)
	p.progressInput.SetWidth(fieldWidth)
}

func (p *EditPage) Update(msg tea.Msg) (core.Page, tea.Cmd) {
	var cmd tea.Cmd

	switch v := msg.(type) {
	case tea.KeyMsg:
		// 加载中或提交中禁用输入
		if p.loading || p.submitting {
			if v.String() == "esc" && p.pop != nil {
				return p, p.pop(core.PagePlan)
			}
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

	case editLoadMsg:
		p.loading = false
		if v.err != nil {
			p.err = v.err
		} else {
			// 填充表单
			p.titleInput.SetValue(v.title)
			p.descriptionArea.SetValue(v.description)
			p.contentArea.SetValue(v.content)
			p.progressInput.SetValue(fmt.Sprintf("%d", v.progress))
			// 聚焦第一个字段
			return p, p.titleInput.Focus()
		}
		return p, nil

	case editResultMsg:
		p.submitting = false
		if v.success {
			// 更新成功，返回列表页
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
		p.titleInput, cmd = p.titleInput.Update(msg)
	case 1:
		p.descriptionArea, cmd = p.descriptionArea.Update(msg)
	case 2:
		p.contentArea, cmd = p.contentArea.Update(msg)
	case 3:
		p.progressInput, cmd = p.progressInput.Update(msg)
	}

	return p, cmd
}

func (p *EditPage) View() string {
	cw, _ := p.frame.ContentSize()
	cardW := layout.FitCardWidth(cw)

	// 加载中状态
	if p.loading {
		return components.LoadingState(theme.IconEdit+" 编辑计划", "加载计划信息中...", cardW)
	}

	// 错误状态
	if p.err != nil && !p.submitting {
		return components.ErrorState(theme.IconEdit+" 编辑计划", p.err.Error(), cardW)
	}

	var body strings.Builder

	// 渲染表单
	body.WriteString(p.titleInput.View())
	body.WriteString("\n\n")
	body.WriteString(p.descriptionArea.View())
	body.WriteString("\n\n")
	body.WriteString(p.contentArea.View())
	body.WriteString("\n\n")
	body.WriteString(p.progressInput.View())
	body.WriteString("\n\n")

	// 提示信息
	if p.submitting {
		body.WriteString(theme.TextDim.Render("正在保存更改..."))
	} else if p.err != nil {
		body.WriteString(theme.FormError.Render(fmt.Sprintf("错误: %s", p.err.Error())))
	} else {
		body.WriteString(theme.FormHint.Render("Ctrl+S 保存 | Tab/Shift+Tab 切换字段 | Esc 取消"))
	}

	title := theme.IconEdit + " 编辑计划"
	return components.Card(title, body.String(), cardW)
}

func (p *EditPage) Meta() core.Meta {
	return core.Meta{
		Title:      "编辑计划",
		Breadcrumb: "计划管理 > 编辑",
		Extra:      "修改表单后按 Ctrl+S 保存",
		Keys: []components.KeyHint{
			{Key: "Ctrl+S", Desc: "保存"},
			{Key: "Tab", Desc: "下一字段"},
			{Key: "Shift+Tab", Desc: "上一字段"},
			{Key: "Esc", Desc: "取消"},
		},
	}
}

// loadPlan 加载计划数据
func (p *EditPage) loadPlan() tea.Cmd {
	return func() tea.Msg {
		ctx := p.bs.Context()
		plan, err := p.bs.PlanService.GetPlanByID(ctx, p.planID)
		if err != nil {
			return editLoadMsg{err: err}
		}

		return editLoadMsg{
			title:       plan.Title,
			description: plan.Description,
			content:     plan.Content,
			progress:    plan.Progress,
		}
	}
}

// nextField 切换到下一个字段
func (p *EditPage) nextField() tea.Cmd {
	p.blurCurrent()
	p.focusIdx = (p.focusIdx + 1) % 4
	return p.focusCurrent()
}

// prevField 切换到上一个字段
func (p *EditPage) prevField() tea.Cmd {
	p.blurCurrent()
	p.focusIdx = (p.focusIdx - 1 + 4) % 4
	return p.focusCurrent()
}

// focusCurrent 聚焦当前字段
func (p *EditPage) focusCurrent() tea.Cmd {
	switch p.focusIdx {
	case 0:
		return p.titleInput.Focus()
	case 1:
		return p.descriptionArea.Focus()
	case 2:
		return p.contentArea.Focus()
	case 3:
		return p.progressInput.Focus()
	}
	return nil
}

// blurCurrent 取消当前字段聚焦
func (p *EditPage) blurCurrent() {
	switch p.focusIdx {
	case 0:
		p.titleInput.Blur()
	case 1:
		p.descriptionArea.Blur()
	case 2:
		p.contentArea.Blur()
	case 3:
		p.progressInput.Blur()
	}
}

// submit 提交表单
func (p *EditPage) submit() tea.Cmd {
	// 验证表单
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
		title := strings.TrimSpace(p.titleInput.Value())
		description := strings.TrimSpace(p.descriptionArea.Value())
		content := strings.TrimSpace(p.contentArea.Value())
		progressStr := strings.TrimSpace(p.progressInput.Value())

		// 先通过ID获取当前计划的Code
		plan, err := p.bs.PlanService.GetPlanByID(ctx, p.planID)
		if err != nil {
			return editResultMsg{success: false, err: err}
		}

		// 构建更新 DTO
		input := &dto.PlanUpdateDTO{
			Code: plan.Code,
		}

		// 只更新非空字段
		if title != "" {
			input.Title = &title
		}
		if description != "" {
			input.Description = &description
		}
		if content != "" {
			input.Content = &content
		}
		if progressStr != "" {
			progress := 0
			if _, err := fmt.Sscanf(progressStr, "%d", &progress); err == nil {
				if progress >= 0 && progress <= 100 {
					input.Progress = &progress
				}
			}
		}

		// 更新计划
		if err := p.bs.PlanService.UpdatePlan(ctx, input); err != nil {
			return editResultMsg{success: false, err: err}
		}

		return editResultMsg{success: true}
	}
}
