package group

import (
	"context"
	"fmt"
	"strings"

	"github.com/XiaoLFeng/llm-memory/internal/tui/common"
	"github.com/XiaoLFeng/llm-memory/internal/tui/styles"
	"github.com/XiaoLFeng/llm-memory/startup"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// CreateModel 组创建模型
// 呀~ 创建新组的表单！📝
type CreateModel struct {
	bs         *startup.Bootstrap
	focusIndex int
	nameInput  textinput.Model
	descArea   textarea.Model
	width      int
	height     int
	err        error
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
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("👥 创建新组"))
	b.WriteString("\n\n")

	// 名称
	b.WriteString(styles.LabelStyle.Render("名称"))
	b.WriteString("\n")
	b.WriteString(m.nameInput.View())
	b.WriteString("\n\n")

	// 描述
	b.WriteString(styles.LabelStyle.Render("描述"))
	b.WriteString("\n")
	b.WriteString(m.descArea.View())
	b.WriteString("\n\n")

	// 错误信息
	if m.err != nil {
		b.WriteString(styles.ErrorStyle.Render("错误: " + m.err.Error()))
		b.WriteString("\n\n")
	}

	// 帮助
	b.WriteString(styles.HelpStyle.Render("Tab 切换 | Ctrl+S 保存 | Esc 返回"))

	return b.String()
}
