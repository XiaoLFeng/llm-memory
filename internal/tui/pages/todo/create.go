package todo

import (
	"fmt"
	"strings"
	"time"

	"github.com/XiaoLFeng/llm-memory/internal/models/dto"
	"github.com/XiaoLFeng/llm-memory/internal/models/entity"
	"github.com/XiaoLFeng/llm-memory/internal/tui/components"
	"github.com/XiaoLFeng/llm-memory/internal/tui/core"
	"github.com/XiaoLFeng/llm-memory/internal/tui/layout"
	"github.com/XiaoLFeng/llm-memory/internal/tui/theme"
	"github.com/XiaoLFeng/llm-memory/startup"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type createSuccessMsg struct{}

type CreatePage struct {
	bs     *startup.Bootstrap
	frame  *layout.Frame
	pop    func(core.PageID) tea.Cmd
	width  int
	height int

	// 表单字段
	codeInput        *components.Input // Code 输入框
	planSelect       *components.Select
	titleInput       *components.Input
	descriptionInput *components.TextArea
	prioritySelect   *components.Select
	dueDateInput     *components.Input
	tagsInput        *components.Input

	// Plan 列表缓存
	planCodes []string

	// 预设 Plan（从 Plan 详情页创建时使用）
	presetPlanCode  string
	presetPlanTitle string

	focusIndex int
	maxFocus   int

	submitting bool
	err        error
}

func NewCreatePage(bs *startup.Bootstrap, pop func(core.PageID) tea.Cmd) *CreatePage {
	return newCreatePageInternal(bs, pop, "", "")
}

// NewCreatePageWithPlan 创建带预设 Plan 的 Todo 创建页
func NewCreatePageWithPlan(bs *startup.Bootstrap, pop func(core.PageID) tea.Cmd, planCode, planTitle string) *CreatePage {
	return newCreatePageInternal(bs, pop, planCode, planTitle)
}

// newCreatePageInternal 内部构造函数
func newCreatePageInternal(bs *startup.Bootstrap, pop func(core.PageID) tea.Cmd, presetPlanCode, presetPlanTitle string) *CreatePage {
	// 初始化表单组件
	codeInput := components.NewInput("标识码", "小写字母+连字符，如: my-todo", true)

	// 获取可用的 Plan 列表
	ctx := bs.Context()
	plans, _ := bs.PlanService.ListPlans(ctx, bs.CurrentScope)
	planOptions := make([]components.SelectOption, 0, len(plans))
	planCodes := make([]string, 0, len(plans))
	selectedPlanIndex := 0
	for _, plan := range plans {
		// 只显示活跃的计划（非已完成/已取消）
		if plan.Status != entity.PlanStatusCompleted && plan.Status != entity.PlanStatusCancelled {
			// 如果有预设 PlanCode，记录索引
			if presetPlanCode != "" && plan.Code == presetPlanCode {
				selectedPlanIndex = len(planOptions)
			}
			planOptions = append(planOptions, components.SelectOption{
				Label: fmt.Sprintf("[%s] %s", plan.Code, plan.Title),
				Value: plan.Code,
			})
			planCodes = append(planCodes, plan.Code)
		}
	}
	// 如果没有可用计划，添加一个提示选项
	if len(planOptions) == 0 {
		planOptions = append(planOptions, components.SelectOption{
			Label: "（暂无可用计划，请先创建计划）",
			Value: "",
		})
	}
	planSelect := components.NewSelect("所属计划", planOptions)
	// 如果有预设 PlanCode，设置选中索引
	if presetPlanCode != "" && selectedPlanIndex < len(planOptions) {
		planSelect.SetSelectedIndex(selectedPlanIndex)
	}

	titleInput := components.NewInput("标题", "待办事项标题", true)
	descriptionInput := components.NewTextArea("描述", "详细描述（可选）", false)
	descriptionInput.SetHeight(4)

	prioritySelect := components.NewSelect("优先级", []components.SelectOption{
		{Label: "低", Value: int(entity.ToDoPriorityLow)},
		{Label: "中", Value: int(entity.ToDoPriorityMedium)},
		{Label: "高", Value: int(entity.ToDoPriorityHigh)},
		{Label: "紧急", Value: int(entity.ToDoPriorityUrgent)},
	})
	prioritySelect.SetSelectedIndex(1) // 默认中等优先级

	dueDateInput := components.NewInput("截止日期", "YYYY-MM-DD（可选）", false)
	tagsInput := components.NewInput("标签", "逗号分隔（可选）", false)

	return &CreatePage{
		bs:               bs,
		frame:            layout.NewFrame(80, 24),
		pop:              pop,
		codeInput:        codeInput,
		planSelect:       planSelect,
		titleInput:       titleInput,
		descriptionInput: descriptionInput,
		prioritySelect:   prioritySelect,
		dueDateInput:     dueDateInput,
		tagsInput:        tagsInput,
		planCodes:        planCodes,
		presetPlanCode:   presetPlanCode,
		presetPlanTitle:  presetPlanTitle,
		maxFocus:         6, // 7个字段，索引0-6
	}
}

func (p *CreatePage) Init() tea.Cmd {
	return p.codeInput.Focus()
}

func (p *CreatePage) Resize(w, h int) {
	p.width, p.height = w, h
	p.frame.Resize(w, h)

	// 设置表单组件宽度
	formWidth := 60
	if w < 70 {
		formWidth = w - 10
	}
	p.codeInput.SetWidth(formWidth)
	p.planSelect.SetWidth(formWidth)
	p.titleInput.SetWidth(formWidth)
	p.descriptionInput.SetWidth(formWidth)
	p.prioritySelect.SetWidth(formWidth)
	p.dueDateInput.SetWidth(formWidth)
	p.tagsInput.SetWidth(formWidth)
}

func (p *CreatePage) Update(msg tea.Msg) (core.Page, tea.Cmd) {
	switch v := msg.(type) {
	case tea.KeyMsg:
		if p.submitting {
			return p, nil
		}

		switch v.String() {
		case "ctrl+s":
			// 提交表单
			return p, p.submit()
		case "tab", "down":
			// 下一个字段
			p.nextField()
		case "shift+tab", "up":
			// 上一个字段
			p.prevField()
		default:
			// 分发到当前聚焦的组件
			return p, p.updateFocused(msg)
		}

	case createSuccessMsg:
		// 创建成功，返回列表
		if p.pop != nil {
			return p, p.pop(core.PageTodo)
		}
	}

	return p, nil
}

func (p *CreatePage) View() string {
	cw, _ := p.frame.ContentSize()
	cardW := layout.FitCardWidth(cw)

	if p.submitting {
		return components.LoadingState(theme.IconTodo+" 创建待办", "正在保存...", cardW)
	}

	if p.err != nil {
		errCard := lipgloss.JoinVertical(lipgloss.Left,
			theme.FormError.Render("错误: "+p.err.Error()),
			"",
			p.renderForm(cardW-6),
		)
		return components.Card(theme.IconTodo+" 创建待办", errCard, cardW)
	}

	body := p.renderForm(cardW - 6)
	return components.Card(theme.IconTodo+" 创建待办", body, cardW)
}

func (p *CreatePage) renderForm(width int) string {
	parts := []string{
		p.codeInput.View(),
		p.planSelect.View(),
		p.titleInput.View(),
		p.descriptionInput.View(),
		p.prioritySelect.View(),
		p.dueDateInput.View(),
		p.tagsInput.View(),
		"",
		theme.FormHint.Render("提示: Tab切换字段 · Ctrl+S保存 · Esc返回"),
	}

	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}

func (p *CreatePage) Meta() core.Meta {
	return core.Meta{
		Title:      "创建待办",
		Breadcrumb: "待办管理 > 创建",
		Extra:      "填写待办信息",
		Keys: []components.KeyHint{
			{Key: "Ctrl+S", Desc: "保存"},
			{Key: "Tab", Desc: "下一字段"},
			{Key: "Esc", Desc: "返回"},
		},
	}
}

// nextField 聚焦下一个字段
func (p *CreatePage) nextField() {
	p.blurCurrent()
	p.focusIndex = (p.focusIndex + 1) % (p.maxFocus + 1)
	p.focusCurrent()
}

// prevField 聚焦上一个字段
func (p *CreatePage) prevField() {
	p.blurCurrent()
	p.focusIndex--
	if p.focusIndex < 0 {
		p.focusIndex = p.maxFocus
	}
	p.focusCurrent()
}

// focusCurrent 聚焦当前字段
func (p *CreatePage) focusCurrent() tea.Cmd {
	switch p.focusIndex {
	case 0:
		return p.codeInput.Focus()
	case 1:
		return p.planSelect.Focus()
	case 2:
		return p.titleInput.Focus()
	case 3:
		return p.descriptionInput.Focus()
	case 4:
		return p.prioritySelect.Focus()
	case 5:
		return p.dueDateInput.Focus()
	case 6:
		return p.tagsInput.Focus()
	}
	return nil
}

// blurCurrent 取消当前字段焦点
func (p *CreatePage) blurCurrent() {
	switch p.focusIndex {
	case 0:
		p.codeInput.Blur()
	case 1:
		p.planSelect.Blur()
	case 2:
		p.titleInput.Blur()
	case 3:
		p.descriptionInput.Blur()
	case 4:
		p.prioritySelect.Blur()
	case 5:
		p.dueDateInput.Blur()
	case 6:
		p.tagsInput.Blur()
	}
}

// updateFocused 更新当前聚焦的组件
func (p *CreatePage) updateFocused(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch p.focusIndex {
	case 0:
		p.codeInput, cmd = p.codeInput.Update(msg)
	case 1:
		p.planSelect, cmd = p.planSelect.Update(msg)
	case 2:
		p.titleInput, cmd = p.titleInput.Update(msg)
	case 3:
		p.descriptionInput, cmd = p.descriptionInput.Update(msg)
	case 4:
		p.prioritySelect, cmd = p.prioritySelect.Update(msg)
	case 5:
		p.dueDateInput, cmd = p.dueDateInput.Update(msg)
	case 6:
		p.tagsInput, cmd = p.tagsInput.Update(msg)
	}
	return cmd
}

// submit 提交表单
func (p *CreatePage) submit() tea.Cmd {
	return func() tea.Msg {
		p.submitting = true
		p.err = nil

		// 验证必填字段
		if err := p.codeInput.Validate(); err != nil {
			p.err = fmt.Errorf("标识码不能为空")
			p.submitting = false
			return nil
		}

		// 验证所属计划
		planCode, ok := p.planSelect.Value().(string)
		if !ok || planCode == "" {
			p.err = fmt.Errorf("请选择所属计划（如没有可用计划，请先创建计划）")
			p.submitting = false
			return nil
		}

		if err := p.titleInput.Validate(); err != nil {
			p.err = fmt.Errorf("标题不能为空")
			p.submitting = false
			return nil
		}

		// 解析截止日期
		var dueDate *time.Time
		dueDateStr := strings.TrimSpace(p.dueDateInput.Value())
		if dueDateStr != "" {
			parsed, err := time.Parse("2006-01-02", dueDateStr)
			if err != nil {
				p.err = fmt.Errorf("日期格式错误，请使用 YYYY-MM-DD 格式")
				p.submitting = false
				return nil
			}
			dueDate = &parsed
		}

		// 解析标签
		var tags []string
		tagsStr := strings.TrimSpace(p.tagsInput.Value())
		if tagsStr != "" {
			rawTags := strings.Split(tagsStr, ",")
			for _, tag := range rawTags {
				trimmed := strings.TrimSpace(tag)
				if trimmed != "" {
					tags = append(tags, trimmed)
				}
			}
		}

		// 构建创建请求
		createDTO := &dto.ToDoCreateDTO{
			Code:        p.codeInput.Value(),
			PlanCode:    planCode,
			Title:       p.titleInput.Value(),
			Description: p.descriptionInput.Value(),
			Priority:    p.prioritySelect.Value().(int),
			DueDate:     dueDate,
			Tags:        tags,
		}

		// 调用服务创建待办
		ctx := p.bs.Context()
		_, err := p.bs.ToDoService.CreateToDo(ctx, createDTO, p.bs.CurrentScope)
		if err != nil {
			p.err = err
			p.submitting = false
			return nil
		}

		return createSuccessMsg{}
	}
}
