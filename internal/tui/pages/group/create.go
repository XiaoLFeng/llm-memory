package group

import (
	"fmt"

	"github.com/XiaoLFeng/llm-memory/internal/tui/components"
	"github.com/XiaoLFeng/llm-memory/internal/tui/core"
	"github.com/XiaoLFeng/llm-memory/internal/tui/layout"
	"github.com/XiaoLFeng/llm-memory/internal/tui/theme"
	"github.com/XiaoLFeng/llm-memory/startup"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// createResult 创建结果消息
type createResult struct {
	err error
}

// CreatePage 组创建页面
type CreatePage struct {
	bs       *startup.Bootstrap
	frame    *layout.Frame
	navigate func(core.PageID) tea.Cmd

	// 表单字段
	nameInput *components.Input
	descArea  *components.TextArea

	// 状态
	focusIndex int
	submitting bool
	err        error
}

// NewCreatePage 创建新建组页面
func NewCreatePage(bs *startup.Bootstrap, navigate func(core.PageID) tea.Cmd) *CreatePage {
	return &CreatePage{
		bs:         bs,
		frame:      layout.NewFrame(80, 24),
		navigate:   navigate,
		nameInput:  components.NewInput("组名称", "输入组名称", true),
		descArea:   components.NewTextArea("描述", "输入组描述（可选）", false),
		focusIndex: 0,
	}
}

func (p *CreatePage) Init() tea.Cmd {
	// 初始化时聚焦第一个输入框
	return p.nameInput.Focus()
}

func (p *CreatePage) Resize(w, h int) {
	p.frame.Resize(w, h)
}

func (p *CreatePage) Update(msg tea.Msg) (core.Page, tea.Cmd) {
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

	case createResult:
		p.submitting = false
		if v.err != nil {
			p.err = v.err
		} else {
			// 创建成功，返回列表页
			if p.navigate != nil {
				return p, p.navigate(core.PageGroup)
			}
		}
	}

	return p, nil
}

func (p *CreatePage) View() string {
	cw, _ := p.frame.ContentSize()
	cardW := layout.FitCardWidth(cw)

	// 如果正在提交
	if p.submitting {
		return components.LoadingState(theme.IconGroup+" 新建组", "正在创建组...", cardW)
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

	return components.Card(theme.IconGroup+" 新建组", body, cardW)
}

func (p *CreatePage) Meta() core.Meta {
	return core.Meta{
		Title:      "新建组",
		Breadcrumb: "组管理 > 新建",
		Extra:      "Ctrl+S 保存",
		Keys: []components.KeyHint{
			{Key: "Tab", Desc: "切换字段"},
			{Key: "Ctrl+S", Desc: "保存"},
			{Key: "Esc", Desc: "取消"},
		},
	}
}

// updateFocus 更新组件焦点状态
func (p *CreatePage) updateFocus() tea.Cmd {
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
func (p *CreatePage) submit() tea.Cmd {
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

	// 异步创建
	return func() tea.Msg {
		ctx := p.bs.Context()
		name := p.nameInput.Value()
		desc := p.descArea.Value()

		_, err := p.bs.GroupService.CreateGroup(ctx, name, desc)
		return createResult{err: err}
	}
}
