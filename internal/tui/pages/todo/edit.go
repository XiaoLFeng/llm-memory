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

type editLoadMsg struct {
	todo *entity.ToDo
	err  error
}

type editSuccessMsg struct{}

type EditPage struct {
	bs     *startup.Bootstrap
	frame  *layout.Frame
	pop    func(core.PageID) tea.Cmd
	todoID int64
	width  int
	height int

	// 加载状态
	loading bool
	loadErr error
	todo    *entity.ToDo

	// 表单字段
	titleInput       *components.Input
	descriptionInput *components.TextArea
	prioritySelect   *components.Select
	statusSelect     *components.Select
	dueDateInput     *components.Input
	tagsInput        *components.Input

	focusIndex int
	maxFocus   int

	submitting bool
	err        error
}

func NewEditPage(bs *startup.Bootstrap, todoID int64, pop func(core.PageID) tea.Cmd) *EditPage {
	// 初始化表单组件
	titleInput := components.NewInput("标题", "待办事项标题", true)
	descriptionInput := components.NewTextArea("描述", "详细描述（可选）", false)
	descriptionInput.SetHeight(4)

	prioritySelect := components.NewSelect("优先级", []components.SelectOption{
		{Label: "低", Value: int(entity.ToDoPriorityLow)},
		{Label: "中", Value: int(entity.ToDoPriorityMedium)},
		{Label: "高", Value: int(entity.ToDoPriorityHigh)},
		{Label: "紧急", Value: int(entity.ToDoPriorityUrgent)},
	})

	statusSelect := components.NewSelect("状态", []components.SelectOption{
		{Label: "待处理", Value: int(entity.ToDoStatusPending)},
		{Label: "进行中", Value: int(entity.ToDoStatusInProgress)},
		{Label: "已完成", Value: int(entity.ToDoStatusCompleted)},
		{Label: "已取消", Value: int(entity.ToDoStatusCancelled)},
	})

	dueDateInput := components.NewInput("截止日期", "YYYY-MM-DD（可选）", false)
	tagsInput := components.NewInput("标签", "逗号分隔（可选）", false)

	return &EditPage{
		bs:               bs,
		frame:            layout.NewFrame(80, 24),
		pop:              pop,
		todoID:           todoID,
		loading:          true,
		titleInput:       titleInput,
		descriptionInput: descriptionInput,
		prioritySelect:   prioritySelect,
		statusSelect:     statusSelect,
		dueDateInput:     dueDateInput,
		tagsInput:        tagsInput,
		maxFocus:         5, // 6个字段，索引0-5
	}
}

func (p *EditPage) Init() tea.Cmd {
	return p.load()
}

func (p *EditPage) load() tea.Cmd {
	return func() tea.Msg {
		ctx := p.bs.Context()
		todo, err := p.bs.ToDoService.GetToDoByID(ctx, p.todoID)
		if err != nil {
			return editLoadMsg{err: err}
		}
		return editLoadMsg{todo: todo}
	}
}

func (p *EditPage) Resize(w, h int) {
	p.width, p.height = w, h
	p.frame.Resize(w, h)

	// 设置表单组件宽度
	formWidth := 60
	if w < 70 {
		formWidth = w - 10
	}
	p.titleInput.SetWidth(formWidth)
	p.descriptionInput.SetWidth(formWidth)
	p.prioritySelect.SetWidth(formWidth)
	p.statusSelect.SetWidth(formWidth)
	p.dueDateInput.SetWidth(formWidth)
	p.tagsInput.SetWidth(formWidth)
}

func (p *EditPage) Update(msg tea.Msg) (core.Page, tea.Cmd) {
	switch v := msg.(type) {
	case editLoadMsg:
		p.loading = false
		if v.err != nil {
			p.loadErr = v.err
			return p, nil
		}
		p.todo = v.todo
		p.populateForm()
		return p, p.titleInput.Focus()

	case tea.KeyMsg:
		if p.loading || p.submitting {
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

	case editSuccessMsg:
		// 更新成功，返回列表
		if p.pop != nil {
			return p, p.pop(core.PageTodo)
		}
	}

	return p, nil
}

func (p *EditPage) View() string {
	cw, _ := p.frame.ContentSize()
	cardW := layout.FitCardWidth(cw)

	if p.loading {
		return components.LoadingState(theme.IconTodo+" 编辑待办", "加载中...", cardW)
	}

	if p.loadErr != nil {
		return components.ErrorState(theme.IconTodo+" 编辑待办", p.loadErr.Error(), cardW)
	}

	if p.submitting {
		return components.LoadingState(theme.IconTodo+" 编辑待办", "正在保存...", cardW)
	}

	if p.err != nil {
		errCard := lipgloss.JoinVertical(lipgloss.Left,
			theme.FormError.Render("错误: "+p.err.Error()),
			"",
			p.renderForm(cardW-6),
		)
		return components.Card(theme.IconTodo+" 编辑待办", errCard, cardW)
	}

	body := p.renderForm(cardW - 6)
	return components.Card(theme.IconTodo+" 编辑待办", body, cardW)
}

func (p *EditPage) renderForm(width int) string {
	parts := []string{
		p.titleInput.View(),
		p.descriptionInput.View(),
		p.prioritySelect.View(),
		p.statusSelect.View(),
		p.dueDateInput.View(),
		p.tagsInput.View(),
		"",
		theme.FormHint.Render("提示: Tab切换字段 · Ctrl+S保存 · Esc返回"),
	}

	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}

func (p *EditPage) Meta() core.Meta {
	return core.Meta{
		Title:      "编辑待办",
		Breadcrumb: "待办管理 > 编辑",
		Extra:      "修改待办信息",
		Keys: []components.KeyHint{
			{Key: "Ctrl+S", Desc: "保存"},
			{Key: "Tab", Desc: "下一字段"},
			{Key: "Esc", Desc: "返回"},
		},
	}
}

// populateForm 填充表单数据
func (p *EditPage) populateForm() {
	if p.todo == nil {
		return
	}

	// 填充标题和描述
	p.titleInput.SetValue(p.todo.Title)
	p.descriptionInput.SetValue(p.todo.Description)

	// 设置优先级选项
	p.prioritySelect.SetSelectedIndex(int(p.todo.Priority) - 1)

	// 设置状态选项
	p.statusSelect.SetSelectedIndex(int(p.todo.Status))

	// 填充截止日期
	if p.todo.DueDate != nil {
		p.dueDateInput.SetValue(p.todo.DueDate.Format("2006-01-02"))
	}

	// 填充标签
	if len(p.todo.Tags) > 0 {
		tags := make([]string, len(p.todo.Tags))
		for i, tag := range p.todo.Tags {
			tags[i] = tag.Tag
		}
		p.tagsInput.SetValue(strings.Join(tags, ", "))
	}
}

// nextField 聚焦下一个字段
func (p *EditPage) nextField() {
	p.blurCurrent()
	p.focusIndex = (p.focusIndex + 1) % (p.maxFocus + 1)
	p.focusCurrent()
}

// prevField 聚焦上一个字段
func (p *EditPage) prevField() {
	p.blurCurrent()
	p.focusIndex--
	if p.focusIndex < 0 {
		p.focusIndex = p.maxFocus
	}
	p.focusCurrent()
}

// focusCurrent 聚焦当前字段
func (p *EditPage) focusCurrent() tea.Cmd {
	switch p.focusIndex {
	case 0:
		return p.titleInput.Focus()
	case 1:
		return p.descriptionInput.Focus()
	case 2:
		return p.prioritySelect.Focus()
	case 3:
		return p.statusSelect.Focus()
	case 4:
		return p.dueDateInput.Focus()
	case 5:
		return p.tagsInput.Focus()
	}
	return nil
}

// blurCurrent 取消当前字段焦点
func (p *EditPage) blurCurrent() {
	switch p.focusIndex {
	case 0:
		p.titleInput.Blur()
	case 1:
		p.descriptionInput.Blur()
	case 2:
		p.prioritySelect.Blur()
	case 3:
		p.statusSelect.Blur()
	case 4:
		p.dueDateInput.Blur()
	case 5:
		p.tagsInput.Blur()
	}
}

// updateFocused 更新当前聚焦的组件
func (p *EditPage) updateFocused(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch p.focusIndex {
	case 0:
		p.titleInput, cmd = p.titleInput.Update(msg)
	case 1:
		p.descriptionInput, cmd = p.descriptionInput.Update(msg)
	case 2:
		p.prioritySelect, cmd = p.prioritySelect.Update(msg)
	case 3:
		p.statusSelect, cmd = p.statusSelect.Update(msg)
	case 4:
		p.dueDateInput, cmd = p.dueDateInput.Update(msg)
	case 5:
		p.tagsInput, cmd = p.tagsInput.Update(msg)
	}
	return cmd
}

// submit 提交表单
func (p *EditPage) submit() tea.Cmd {
	return func() tea.Msg {
		p.submitting = true
		p.err = nil

		// 验证必填字段
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

		// 构建更新请求
		title := p.titleInput.Value()
		description := p.descriptionInput.Value()
		priority := p.prioritySelect.Value().(int)
		status := p.statusSelect.Value().(int)

		// 先通过ID获取当前待办的Code
		ctx := p.bs.Context()
		todo, err := p.bs.ToDoService.GetToDoByID(ctx, p.todoID)
		if err != nil {
			return editLoadMsg{err: err}
		}

		updateDTO := &dto.ToDoUpdateDTO{
			Code:        todo.Code,
			Title:       &title,
			Description: &description,
			Priority:    &priority,
			Status:      &status,
			DueDate:     dueDate,
			Tags:        &tags,
		}

		// 调用服务更新待办
		err = p.bs.ToDoService.UpdateToDo(ctx, updateDTO)
		if err != nil {
			p.err = err
			p.submitting = false
			return nil
		}

		return editSuccessMsg{}
	}
}
