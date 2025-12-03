package plan

import (
	"context"
	"fmt"
	"strings"

	"github.com/XiaoLFeng/llm-memory/internal/models/dto"
	"github.com/XiaoLFeng/llm-memory/internal/tui/common"
	"github.com/XiaoLFeng/llm-memory/internal/tui/components"
	"github.com/XiaoLFeng/llm-memory/internal/tui/styles"
	"github.com/XiaoLFeng/llm-memory/startup"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// CreateModel 计划创建模型
// 呀~ 创建新计划的表单！
type CreateModel struct {
	bs         *startup.Bootstrap
	focusIndex int
	titleInput textinput.Model
	descArea   textarea.Model
	width      int
	height     int
	err        error
	frame      *components.Frame
}

// NewCreateModel 创建计划创建模型
func NewCreateModel(bs *startup.Bootstrap) *CreateModel {
	// 标题输入框
	ti := textinput.New()
	ti.Placeholder = "计划标题"
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 50

	// 描述输入框
	ta := textarea.New()
	ta.Placeholder = "计划描述（可选）..."
	ta.SetWidth(50)
	ta.SetHeight(6)

	return &CreateModel{
		bs:         bs,
		titleInput: ti,
		descArea:   ta,
		frame:      components.NewFrame(80, 24),
	}
}

// Title 返回页面标题
func (m *CreateModel) Title() string {
	return "创建计划"
}

// ShortHelp 返回快捷键帮助
func (m *CreateModel) ShortHelp() []key.Binding {
	return []key.Binding{common.KeyTab, common.KeyEnter, common.KeyBack}
}

// Init 初始化
func (m *CreateModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update 处理输入
func (m *CreateModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, common.Back()

		case "tab", "shift+tab":
			// 切换焦点
			if msg.String() == "tab" {
				m.focusIndex = (m.focusIndex + 1) % 2
			} else {
				m.focusIndex = (m.focusIndex - 1 + 2) % 2
			}
			m.updateFocus()

		case "ctrl+s":
			// 保存
			return m, m.save()
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.frame.SetSize(msg.Width, msg.Height)

	case planCreatedMsg:
		return m, tea.Batch(
			common.ShowToast("计划创建成功！", common.ToastSuccess),
			common.Back(),
		)

	case plansErrorMsg:
		m.err = msg.err
	}

	// 更新当前聚焦的输入框
	cmd := m.updateInputs(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// updateFocus 更新焦点状态
func (m *CreateModel) updateFocus() {
	m.titleInput.Blur()
	m.descArea.Blur()

	switch m.focusIndex {
	case 0:
		m.titleInput.Focus()
	case 1:
		m.descArea.Focus()
	}
}

// updateInputs 更新输入框
func (m *CreateModel) updateInputs(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	switch m.focusIndex {
	case 0:
		m.titleInput, cmd = m.titleInput.Update(msg)
	case 1:
		m.descArea, cmd = m.descArea.Update(msg)
	}

	return cmd
}

type planCreatedMsg struct{}

// save 保存计划
func (m *CreateModel) save() tea.Cmd {
	return func() tea.Msg {
		title := strings.TrimSpace(m.titleInput.Value())
		if title == "" {
			return plansErrorMsg{err: fmt.Errorf("标题不能为空")}
		}

		description := strings.TrimSpace(m.descArea.Value())

		// 使用 DTO 创建计划
		createDTO := &dto.PlanCreateDTO{
			Title:       title,
			Description: description,
			Scope:       "global",
		}
		_, err := m.bs.PlanService.CreatePlan(context.Background(), createDTO, nil)
		if err != nil {
			return plansErrorMsg{err: err}
		}

		return planCreatedMsg{}
	}
}

// View 渲染界面
func (m *CreateModel) View() string {
	// 构建表单内容
	var formParts []string

	// 标题输入
	titleLabel := lipgloss.NewStyle().
		Foreground(styles.Subtext1).
		Bold(true).
		Render("标题")

	titleInput := m.titleInput.View()
	if m.focusIndex == 0 {
		titleInput = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(styles.Primary).
			Padding(0, 1).
			Render(titleInput)
	}
	formParts = append(formParts, titleLabel+"\n"+titleInput)

	// 描述输入
	descLabel := lipgloss.NewStyle().
		Foreground(styles.Subtext1).
		Bold(true).
		Render("描述")

	descArea := m.descArea.View()
	if m.focusIndex == 1 {
		descArea = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(styles.Primary).
			Padding(0, 1).
			Render(descArea)
	}
	formParts = append(formParts, descLabel+"\n"+descArea)

	// 提示信息
	hint := lipgloss.NewStyle().
		Foreground(styles.Overlay1).
		Italic(true).
		Render(styles.IconBulb + " 提示：按 tab 切换输入框，ctrl+s 保存")
	formParts = append(formParts, hint)

	// 错误信息
	if m.err != nil {
		errorBox := components.CardError("错误", m.err.Error(), 60)
		formParts = append(formParts, errorBox)
	}

	formContent := strings.Join(formParts, "\n\n")

	// 用卡片包装表单
	cardContent := components.Card(styles.IconEdit+" 创建新计划", formContent, m.frame.GetContentWidth()-4)

	// 居中显示
	content := lipgloss.Place(
		m.frame.GetContentWidth(),
		m.frame.GetContentHeight(),
		lipgloss.Center,
		lipgloss.Center,
		cardContent,
	)

	keys := []string{
		"tab 切换",
		"ctrl+s 保存",
		"esc 取消",
	}

	return m.frame.Render("计划管理 > 创建计划", content, keys, "")
}
