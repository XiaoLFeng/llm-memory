package memory

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

// CreateModel 记忆创建模型
// 呀~ 创建新记忆的表单！✨
type CreateModel struct {
	bs            *startup.Bootstrap
	focusIndex    int
	titleInput    textinput.Model
	contentArea   textarea.Model
	categoryInput textinput.Model
	tagsInput     textinput.Model
	width         int
	height        int
	err           error
}

// NewCreateModel 创建记忆创建模型
func NewCreateModel(bs *startup.Bootstrap) *CreateModel {
	// 标题输入框
	ti := textinput.New()
	ti.Placeholder = "记忆标题"
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 50

	// 内容输入框
	ta := textarea.New()
	ta.Placeholder = "记忆内容..."
	ta.SetWidth(50)
	ta.SetHeight(6)

	// 分类输入框
	ci := textinput.New()
	ci.Placeholder = "分类（可选）"
	ci.CharLimit = 50
	ci.Width = 50

	// 标签输入框
	tgi := textinput.New()
	tgi.Placeholder = "标签，用逗号分隔（可选）"
	tgi.CharLimit = 200
	tgi.Width = 50

	return &CreateModel{
		bs:            bs,
		titleInput:    ti,
		contentArea:   ta,
		categoryInput: ci,
		tagsInput:     tgi,
	}
}

// Title 返回页面标题
func (m *CreateModel) Title() string {
	return "创建记忆"
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
				m.focusIndex = (m.focusIndex + 1) % 4
			} else {
				m.focusIndex = (m.focusIndex - 1 + 4) % 4
			}
			m.updateFocus()

		case "ctrl+s":
			// 保存
			return m, m.save()
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case memoryCreatedMsg:
		return m, tea.Batch(
			common.ShowToast("记忆创建成功！", common.ToastSuccess),
			common.Back(),
		)

	case memoriesErrorMsg:
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
	m.contentArea.Blur()
	m.categoryInput.Blur()
	m.tagsInput.Blur()

	switch m.focusIndex {
	case 0:
		m.titleInput.Focus()
	case 1:
		m.contentArea.Focus()
	case 2:
		m.categoryInput.Focus()
	case 3:
		m.tagsInput.Focus()
	}
}

// updateInputs 更新输入框
func (m *CreateModel) updateInputs(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	switch m.focusIndex {
	case 0:
		m.titleInput, cmd = m.titleInput.Update(msg)
	case 1:
		m.contentArea, cmd = m.contentArea.Update(msg)
	case 2:
		m.categoryInput, cmd = m.categoryInput.Update(msg)
	case 3:
		m.tagsInput, cmd = m.tagsInput.Update(msg)
	}

	return cmd
}

type memoryCreatedMsg struct{}

// save 保存记忆
func (m *CreateModel) save() tea.Cmd {
	return func() tea.Msg {
		title := strings.TrimSpace(m.titleInput.Value())
		if title == "" {
			return memoriesErrorMsg{err: fmt.Errorf("标题不能为空")}
		}

		content := strings.TrimSpace(m.contentArea.Value())
		if content == "" {
			return memoriesErrorMsg{err: fmt.Errorf("内容不能为空")}
		}

		category := strings.TrimSpace(m.categoryInput.Value())
		if category == "" {
			category = "默认"
		}

		var tags []string
		tagsStr := strings.TrimSpace(m.tagsInput.Value())
		if tagsStr != "" {
			for _, tag := range strings.Split(tagsStr, ",") {
				tag = strings.TrimSpace(tag)
				if tag != "" {
					tags = append(tags, tag)
				}
			}
		}

		_, err := m.bs.MemoryService.CreateMemory(context.Background(), title, content, category, tags, 2)
		if err != nil {
			return memoriesErrorMsg{err: err}
		}

		return memoryCreatedMsg{}
	}
}

// View 渲染界面
func (m *CreateModel) View() string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("✨ 创建新记忆"))
	b.WriteString("\n\n")

	// 标题
	b.WriteString(styles.LabelStyle.Render("标题"))
	b.WriteString("\n")
	b.WriteString(m.titleInput.View())
	b.WriteString("\n\n")

	// 内容
	b.WriteString(styles.LabelStyle.Render("内容"))
	b.WriteString("\n")
	b.WriteString(m.contentArea.View())
	b.WriteString("\n\n")

	// 分类
	b.WriteString(styles.LabelStyle.Render("分类"))
	b.WriteString("\n")
	b.WriteString(m.categoryInput.View())
	b.WriteString("\n\n")

	// 标签
	b.WriteString(styles.LabelStyle.Render("标签"))
	b.WriteString("\n")
	b.WriteString(m.tagsInput.View())
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
