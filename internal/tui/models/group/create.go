package group

import (
	"context"
	"fmt"
	"strings"

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

// CreateModel 组创建模型
type CreateModel struct {
	bs         *startup.Bootstrap
	focusIndex int
	nameInput  textinput.Model
	descArea   textarea.Model
	width      int
	height     int
	err        error
	frame      *components.Frame
}

// NewCreateModel 创建组创建模型
func NewCreateModel(bs *startup.Bootstrap) *CreateModel {
	// 名称输入框
	ni := textinput.New()
	ni.Placeholder = "组名称"
	ni.Focus()
	ni.CharLimit = 50
	ni.Width = 50

	// 描述输入框
	ta := textarea.New()
	ta.Placeholder = "组描述（可选）..."
	ta.SetWidth(50)
	ta.SetHeight(4)

	return &CreateModel{
		bs:        bs,
		nameInput: ni,
		descArea:  ta,
		frame:     components.NewFrame(80, 24),
	}
}

// Title 返回页面标题
func (m *CreateModel) Title() string {
	return "创建组"
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

	case groupCreatedMsg:
		return m, tea.Batch(
			common.ShowToast("组创建成功！", common.ToastSuccess),
			common.Back(),
		)

	case groupsErrorMsg:
		m.err = msg.err
	}

	// 更新当前聚焦的输入框
	cmd := m.updateInputs(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// updateFocus 更新焦点状态
func (m *CreateModel) updateFocus() {
	m.nameInput.Blur()
	m.descArea.Blur()

	switch m.focusIndex {
	case 0:
		m.nameInput.Focus()
	case 1:
		m.descArea.Focus()
	}
}

// updateInputs 更新输入框
func (m *CreateModel) updateInputs(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	switch m.focusIndex {
	case 0:
		m.nameInput, cmd = m.nameInput.Update(msg)
	case 1:
		m.descArea, cmd = m.descArea.Update(msg)
	}

	return cmd
}

type groupCreatedMsg struct{}

// save 保存组
func (m *CreateModel) save() tea.Cmd {
	return func() tea.Msg {
		name := strings.TrimSpace(m.nameInput.Value())
		if name == "" {
			return groupsErrorMsg{err: fmt.Errorf("组名称不能为空")}
		}

		description := strings.TrimSpace(m.descArea.Value())

		_, err := m.bs.GroupService.CreateGroup(context.Background(), name, description)
		if err != nil {
			return groupsErrorMsg{err: err}
		}

		return groupCreatedMsg{}
	}
}

// View 渲染界面
func (m *CreateModel) View() string {
	// 计算卡片宽度
	cardWidth := m.frame.GetContentWidth() - 4
	if cardWidth > 70 {
		cardWidth = 70
	}
	if cardWidth < 60 {
		cardWidth = 60
	}

	// 表单内容
	var formContent strings.Builder

	// 名称输入
	labelStyle := lipgloss.NewStyle().
		Foreground(styles.Subtext0).
		Bold(true)

	formContent.WriteString(labelStyle.Render("名称"))
	formContent.WriteString("\n")

	// 输入框样式
	inputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Border).
		Width(cardWidth-8).
		Padding(0, 1)

	if m.focusIndex == 0 {
		inputStyle = inputStyle.BorderForeground(styles.Primary)
	}

	formContent.WriteString(inputStyle.Render(m.nameInput.View()))
	formContent.WriteString("\n\n")

	// 描述输入
	formContent.WriteString(labelStyle.Render("描述（可选）"))
	formContent.WriteString("\n")

	descStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Border).
		Width(cardWidth-8).
		Padding(0, 1)

	if m.focusIndex == 1 {
		descStyle = descStyle.BorderForeground(styles.Primary)
	}

	formContent.WriteString(descStyle.Render(m.descArea.View()))

	// 错误信息
	if m.err != nil {
		formContent.WriteString("\n\n")
		errorStyle := lipgloss.NewStyle().Foreground(styles.Error)
		formContent.WriteString(errorStyle.Render("错误: " + m.err.Error()))
	}

	// 使用卡片包装表单
	card := components.Card(styles.IconUsers+" 创建新组", formContent.String(), cardWidth)

	// 居中显示
	centeredContent := lipgloss.Place(
		m.frame.GetContentWidth(),
		m.frame.GetContentHeight(),
		lipgloss.Center,
		lipgloss.Center,
		card,
	)

	// 快捷键
	keys := []string{
		styles.StatusKeyStyle.Render("Tab") + " " + styles.StatusValueStyle.Render("切换"),
		styles.StatusKeyStyle.Render("Ctrl+S") + " " + styles.StatusValueStyle.Render("保存"),
		styles.StatusKeyStyle.Render("esc") + " " + styles.StatusValueStyle.Render("返回"),
	}

	return m.frame.Render("组管理 > "+styles.IconUsers+" 创建组", centeredContent, keys, "")
}
