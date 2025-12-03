package memory

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

// CreateModel 记忆创建模型
// 呀~ 创建新记忆的表单！✨
type CreateModel struct {
	bs            *startup.Bootstrap
	focusIndex    int
	titleInput    textinput.Model
	contentArea   textarea.Model
	categoryInput textinput.Model
	tagsInput     textinput.Model
	global        bool
	frame         *components.Frame
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
	ti.Width = 60
	ti.PromptStyle = lipgloss.NewStyle().Foreground(styles.Primary)
	ti.TextStyle = lipgloss.NewStyle().Foreground(styles.Text)

	// 内容输入框
	ta := textarea.New()
	ta.Placeholder = "记忆内容..."
	ta.SetWidth(60)
	ta.SetHeight(8)
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.FocusedStyle.Prompt = lipgloss.NewStyle().Foreground(styles.Primary)
	ta.FocusedStyle.Text = lipgloss.NewStyle().Foreground(styles.Text)
	ta.BlurredStyle.Prompt = lipgloss.NewStyle().Foreground(styles.Overlay0)
	ta.BlurredStyle.Text = lipgloss.NewStyle().Foreground(styles.Subtext0)

	// 分类输入框
	ci := textinput.New()
	ci.Placeholder = "分类（可选）"
	ci.CharLimit = 50
	ci.Width = 60
	ci.PromptStyle = lipgloss.NewStyle().Foreground(styles.Primary)
	ci.TextStyle = lipgloss.NewStyle().Foreground(styles.Text)

	// 标签输入框
	tgi := textinput.New()
	tgi.Placeholder = "标签，用逗号分隔（可选）"
	tgi.CharLimit = 200
	tgi.Width = 60
	tgi.PromptStyle = lipgloss.NewStyle().Foreground(styles.Primary)
	tgi.TextStyle = lipgloss.NewStyle().Foreground(styles.Text)

	return &CreateModel{
		bs:            bs,
		titleInput:    ti,
		contentArea:   ta,
		categoryInput: ci,
		tagsInput:     tgi,
		frame:         components.NewFrame(80, 24),
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

		case "g":
			// 切换全局/私有
			m.global = !m.global
			target := "当前路径/组内"
			if m.global {
				target = "全局"
			}
			return m, common.ShowToast(fmt.Sprintf("已切换为 %s", target), common.ToastInfo)
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.frame.SetSize(msg.Width, msg.Height)

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

		// 使用 DTO 创建记忆
		createDTO := &dto.MemoryCreateDTO{
			Title:    title,
			Content:  content,
			Category: category,
			Tags:     tags,
			Priority: 2,
			Global:   m.global,
		}
		_, err := m.bs.MemoryService.CreateMemory(context.Background(), createDTO, m.bs.CurrentScope)
		if err != nil {
			return memoriesErrorMsg{err: err}
		}

		return memoryCreatedMsg{}
	}
}

// View 渲染界面
func (m *CreateModel) View() string {
	// 计算表单卡片宽度
	// 构建表单内容
	var formContent strings.Builder

	// 标题字段
	labelStyle := lipgloss.NewStyle().
		Foreground(styles.Text).
		Bold(true).
		MarginBottom(1)

	formContent.WriteString(labelStyle.Render("标题"))
	formContent.WriteString("\n")
	formContent.WriteString(m.titleInput.View())
	formContent.WriteString("\n\n")

	// 内容字段
	formContent.WriteString(labelStyle.Render("内容"))
	formContent.WriteString("\n")
	formContent.WriteString(m.contentArea.View())
	formContent.WriteString("\n\n")

	// 分类字段
	formContent.WriteString(labelStyle.Render("分类"))
	formContent.WriteString("\n")
	formContent.WriteString(m.categoryInput.View())
	formContent.WriteString("\n\n")

	// 标签字段
	formContent.WriteString(labelStyle.Render("标签"))
	formContent.WriteString("\n")
	formContent.WriteString(m.tagsInput.View())
	formContent.WriteString("\n\n")

	// 作用域切换
	scopeLabel := styles.IconGlobe + " 作用域 (按 g 切换)"
	scopeValue := "当前路径/组内"
	if m.global {
		scopeValue = "全局"
	}
	formContent.WriteString(labelStyle.Render(scopeLabel))
	formContent.WriteString("\n")
	formContent.WriteString(lipgloss.NewStyle().Foreground(styles.Accent).Render(scopeValue))
	formContent.WriteString("\n")

	// 错误信息
	if m.err != nil {
		formContent.WriteString("\n")
		errorStyle := lipgloss.NewStyle().Foreground(styles.Error)
		formContent.WriteString(errorStyle.Render("错误: " + m.err.Error()))
	}

	// 将表单包装在卡片中
	centeredContent := components.RenderCard(m.frame, "创建新记忆", formContent.String(), 56, 72, 4, lipgloss.Top)

	// 快捷键
	keys := []string{
		styles.StatusKeyStyle.Render("Tab") + " " + styles.StatusValueStyle.Render("切换"),
		styles.StatusKeyStyle.Render("Ctrl+S") + " " + styles.StatusValueStyle.Render("保存"),
		styles.StatusKeyStyle.Render("esc") + " " + styles.StatusValueStyle.Render("取消"),
	}

	return m.frame.Render("记忆管理 > 创建记忆", centeredContent, keys, "")
}
