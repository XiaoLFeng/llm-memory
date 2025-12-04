package group

import (
	"fmt"

	"github.com/XiaoLFeng/llm-memory/internal/models/entity"
	"github.com/XiaoLFeng/llm-memory/internal/tui/components"
	"github.com/XiaoLFeng/llm-memory/internal/tui/core"
	"github.com/XiaoLFeng/llm-memory/internal/tui/layout"
	"github.com/XiaoLFeng/llm-memory/internal/tui/theme"
	"github.com/XiaoLFeng/llm-memory/startup"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// loadGroupMsg 加载组消息
type loadGroupMsg struct {
	group *entity.Group
	err   error
}

// updateResult 更新结果消息
type updateResult struct {
	err error
}

// EditPage 组编辑页面
type EditPage struct {
	bs       *startup.Bootstrap
	frame    *layout.Frame
	navigate func(core.PageID) tea.Cmd

	// 数据
	groupID int64
	group   *entity.Group

	// 表单字段
	nameInput *components.Input
	descArea  *components.TextArea

	// 状态
	loading    bool
	focusIndex int
	submitting bool
	err        error
}

// NewEditPage 创建编辑组页面
func NewEditPage(bs *startup.Bootstrap, navigate func(core.PageID) tea.Cmd, data interface{}) *EditPage {
	var groupID int64
	if id, ok := data.(int64); ok {
		groupID = id
	}

	return &EditPage{
		bs:         bs,
		frame:      layout.NewFrame(80, 24),
		navigate:   navigate,
		groupID:    groupID,
		nameInput:  components.NewInput("组名称", "输入组名称", true),
		descArea:   components.NewTextArea("描述", "输入组描述（可选）", false),
		focusIndex: 0,
		loading:    true,
	}
}

func (p *EditPage) Init() tea.Cmd {
	return p.loadGroup()
}

func (p *EditPage) Resize(w, h int) {
	p.frame.Resize(w, h)
}

func (p *EditPage) Update(msg tea.Msg) (core.Page, tea.Cmd) {
	switch v := msg.(type) {
	case tea.KeyMsg:
		// 提交中不接受其他按键
		if p.submitting {
			return p, nil
		}

		switch v.String() {
		case "esc":
			// 返回列表页
			if p.navigate != nil {
				return p, p.navigate(core.PageGroup)
			}
		case "tab", "shift+tab":
			// 切换焦点
			if v.String() == "tab" {
				p.focusIndex = (p.focusIndex + 1) % 2
			} else {
				p.focusIndex = (p.focusIndex - 1 + 2) % 2
			}
			return p, p.updateFocus()
		case "ctrl+s":
			// 提交表单
			return p, p.submit()
		}

		// 将按键传递给当前焦点的组件
		var cmd tea.Cmd
		if p.focusIndex == 0 {
			p.nameInput, cmd = p.nameInput.Update(msg)
		} else {
			p.descArea, cmd = p.descArea.Update(msg)
		}
		return p, cmd

	case loadGroupMsg:
		p.loading = false
		if v.err != nil {
			p.err = v.err
		} else {
			p.group = v.group
			// 填充表单
			p.nameInput.SetValue(v.group.Name)
			p.descArea.SetValue(v.group.Description)
			// 聚焦第一个字段
			return p, p.nameInput.Focus()
		}

	case updateResult:
		p.submitting = false
		if v.err != nil {
			p.err = v.err
		} else {
			// 更新成功，返回列表页
			if p.navigate != nil {
				return p, p.navigate(core.PageGroup)
			}
		}
	}

	return p, nil
}

func (p *EditPage) View() string {
	cw, _ := p.frame.ContentSize()
	cardW := layout.FitCardWidth(cw)

	// 加载中
	if p.loading {
		return components.LoadingState(theme.IconGroup+" 编辑组", "加载组信息中...", cardW)
	}

	// 加载失败
	if p.err != nil && p.group == nil {
		return components.ErrorState(theme.IconGroup+" 编辑组", p.err.Error(), cardW)
	}

	// 如果正在提交
	if p.submitting {
		return components.LoadingState(theme.IconGroup+" 编辑组", "正在保存更改...", cardW)
	}

	// 表单内容
	formWidth := cardW - 10
	p.nameInput.SetWidth(formWidth)
	p.descArea.SetWidth(formWidth)
	p.descArea.SetHeight(5)

	// 渲染表单
	form := lipgloss.JoinVertical(lipgloss.Left,
		p.nameInput.View(),
		"",
		p.descArea.View(),
	)

	// 提示信息
	hint := theme.FormHint.Render("Tab 切换字段 · Ctrl+S 保存 · Esc 返回")

	// 错误提示
	errMsg := ""
	if p.err != nil {
		errMsg = theme.FormError.Render(fmt.Sprintf("错误: %s", p.err.Error()))
	}

	body := lipgloss.JoinVertical(lipgloss.Left, form, "", hint, errMsg)

	title := fmt.Sprintf("%s 编辑组 - %s", theme.IconGroup, p.group.Name)
	return components.Card(title, body, cardW)
}

func (p *EditPage) Meta() core.Meta {
	title := "编辑组"
	if p.group != nil {
		title = fmt.Sprintf("编辑组 - %s", p.group.Name)
	}
	return core.Meta{
		Title:      title,
		Breadcrumb: "组管理 > 编辑",
		Extra:      "Ctrl+S 保存",
		Keys: []components.KeyHint{
			{Key: "Tab", Desc: "切换字段"},
			{Key: "Ctrl+S", Desc: "保存"},
			{Key: "Esc", Desc: "取消"},
		},
	}
}

// loadGroup 加载组数据
func (p *EditPage) loadGroup() tea.Cmd {
	return func() tea.Msg {
		ctx := p.bs.Context()
		group, err := p.bs.GroupService.GetGroup(ctx, p.groupID)
		return loadGroupMsg{group: group, err: err}
	}
}

// updateFocus 更新组件焦点状态
func (p *EditPage) updateFocus() tea.Cmd {
	// 先清除所有焦点
	p.nameInput.Blur()
	p.descArea.Blur()

	// 聚焦当前字段
	if p.focusIndex == 0 {
		return p.nameInput.Focus()
	}
	return p.descArea.Focus()
}

// submit 提交表单
func (p *EditPage) submit() tea.Cmd {
	// 验证表单
	if err := p.nameInput.Validate(); err != nil {
		p.nameInput.SetError(err)
		return nil
	}

	// 清除之前的错误
	p.nameInput.SetError(nil)
	p.descArea.SetError(nil)
	p.err = nil

	// 设置提交状态
	p.submitting = true

	// 异步更新
	return func() tea.Msg {
		ctx := p.bs.Context()
		name := p.nameInput.Value()
		desc := p.descArea.Value()

		err := p.bs.GroupService.UpdateGroup(ctx, p.groupID, &name, &desc)
		return updateResult{err: err}
	}
}
