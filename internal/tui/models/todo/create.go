package todo

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/XiaoLFeng/llm-memory/internal/tui/common"
	"github.com/XiaoLFeng/llm-memory/internal/tui/styles"
	"github.com/XiaoLFeng/llm-memory/pkg/types"
	"github.com/XiaoLFeng/llm-memory/startup"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// CreateModel 待办创建模型
// 呀~ 创建新待办的表单！📝
type CreateModel struct {
	bs            *startup.Bootstrap
	focusIndex    int
	titleInput    textinput.Model
	descArea      textarea.Model
	priorityInput textinput.Model
	width         int
	height        int
	err           error
}

// NewCreateModel 创建待办创建模型
func NewCreateModel(bs *startup.Bootstrap) *CreateModel {
	// 标题输入框
	ti := textinput.New()
	ti.Placeholder = "待办标题"
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 50

	// 描述输入框
	ta := textarea.New()
	ta.Placeholder = "待办描述（可选）..."
	ta.SetWidth(50)
	ta.SetHeight(4)

	// 优先级输入框
	pi := textinput.New()
	pi.Placeholder = "1-4"
	pi.CharLimit = 1
	pi.Width = 10
	pi.SetValue("2")

	return &CreateModel{
		bs:            bs,
		titleInput:    ti,
		descArea:      ta,
		priorityInput: pi,
	}
}

// Title 返回页面标题
func (m *CreateModel) Title() string {
	return "创建待办"
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
				m.focusIndex = (m.focusIndex + 1) % 3
			} else {
				m.focusIndex = (m.focusIndex - 1 + 3) % 3
			}
			m.updateFocus()

		case "ctrl+s":
			// 保存
			return m, m.save()
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case todoCreatedMsg:
		return m, tea.Batch(
			common.ShowToast("待办创建成功！", common.ToastSuccess),
			common.Back(),
		)

	case todosErrorMsg:
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
	m.priorityInput.Blur()

	switch m.focusIndex {
	case 0:
		m.titleInput.Focus()
	case 1:
		m.descArea.Focus()
	case 2:
		m.priorityInput.Focus()
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
	case 2:
		m.priorityInput, cmd = m.priorityInput.Update(msg)
	}

	return cmd
}

type todoCreatedMsg struct{}

// save 保存待办
func (m *CreateModel) save() tea.Cmd {
	return func() tea.Msg {
		title := strings.TrimSpace(m.titleInput.Value())
		if title == "" {
			return todosErrorMsg{err: fmt.Errorf("标题不能为空")}
		}

		description := strings.TrimSpace(m.descArea.Value())

		priorityStr := strings.TrimSpace(m.priorityInput.Value())
		priority := 2
		if priorityStr != "" {
			p, err := strconv.Atoi(priorityStr)
			if err != nil || p < 1 || p > 4 {
				return todosErrorMsg{err: fmt.Errorf("优先级必须是 1-4 之间的数字")}
			}
			priority = p
		}

		_, err := m.bs.TodoService.CreateTodo(context.Background(), title, description, types.Priority(priority), nil)
		if err != nil {
			return todosErrorMsg{err: err}
		}

		return todoCreatedMsg{}
	}
}

// View 渲染界面
func (m *CreateModel) View() string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("📝 创建新待办"))
	b.WriteString("\n\n")

	// 标题
	b.WriteString(styles.LabelStyle.Render("标题"))
	b.WriteString("\n")
	b.WriteString(m.titleInput.View())
	b.WriteString("\n\n")

	// 描述
	b.WriteString(styles.LabelStyle.Render("描述"))
	b.WriteString("\n")
	b.WriteString(m.descArea.View())
	b.WriteString("\n\n")

	// 优先级
	b.WriteString(styles.LabelStyle.Render("优先级 (1低/2中/3高/4紧急)"))
	b.WriteString("\n")
	b.WriteString(m.priorityInput.View())
	b.WriteString("\n\n")

	// 错误信息
	if m.err != nil {
		b.WriteString(styles.ErrorStyle.Render("错误: " + m.err.Error()))
		b.WriteString("\n\n")
	}

	// 帮助信息
	b.WriteString(styles.HelpStyle.Render("tab 切换 | ctrl+s 保存 | esc 取消"))

	return b.String()
}
