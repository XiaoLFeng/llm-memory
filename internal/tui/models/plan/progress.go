package plan

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/XiaoLFeng/llm-memory/internal/tui/common"
	"github.com/XiaoLFeng/llm-memory/internal/tui/components"
	"github.com/XiaoLFeng/llm-memory/internal/tui/styles"
	"github.com/XiaoLFeng/llm-memory/internal/tui/utils"
	"github.com/XiaoLFeng/llm-memory/startup"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ProgressModel 计划进度更新模型
type ProgressModel struct {
	bs       *startup.Bootstrap
	id       int64
	progress int
	input    textinput.Model
	width    int
	height   int
	err      error
	frame    *components.Frame
}

// NewProgressModel 创建计划进度更新模型
func NewProgressModel(bs *startup.Bootstrap, id int64, progress int) *ProgressModel {
	ti := textinput.New()
	ti.Placeholder = "0-100"
	ti.Focus()
	ti.CharLimit = 3
	ti.Width = 10
	ti.SetValue(strconv.Itoa(progress))

	return &ProgressModel{
		bs:       bs,
		id:       id,
		progress: progress,
		input:    ti,
		frame:    components.NewFrame(80, 24),
	}
}

// Title 返回页面标题
func (m *ProgressModel) Title() string {
	return "更新进度"
}

// ShortHelp 返回快捷键帮助
func (m *ProgressModel) ShortHelp() []key.Binding {
	return []key.Binding{common.KeyEnter, common.KeyBack}
}

// Init 初始化
func (m *ProgressModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update 处理输入
func (m *ProgressModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, common.Back()

		case "enter", "ctrl+s":
			// 保存
			return m, m.save()

		case "up":
			// 增加进度
			m.adjustProgress(10)

		case "down":
			// 减少进度
			m.adjustProgress(-10)
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.frame.SetSize(msg.Width, msg.Height)

	case progressUpdatedMsg:
		return m, tea.Batch(
			common.ShowToast("进度已更新", common.ToastSuccess),
			common.Back(),
		)

	case plansErrorMsg:
		m.err = msg.err
	}

	// 更新输入框
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// adjustProgress 调整进度
func (m *ProgressModel) adjustProgress(delta int) {
	progress, _ := strconv.Atoi(m.input.Value())
	progress += delta
	if progress < 0 {
		progress = 0
	}
	if progress > 100 {
		progress = 100
	}
	m.input.SetValue(strconv.Itoa(progress))
}

type progressUpdatedMsg struct{}

// save 保存进度
func (m *ProgressModel) save() tea.Cmd {
	return func() tea.Msg {
		progressStr := strings.TrimSpace(m.input.Value())
		progress, err := strconv.Atoi(progressStr)
		if err != nil {
			return plansErrorMsg{err: fmt.Errorf("请输入有效的数字")}
		}

		if progress < 0 || progress > 100 {
			return plansErrorMsg{err: fmt.Errorf("进度必须在 0-100 之间")}
		}

		err = m.bs.PlanService.UpdateProgress(context.Background(), m.id, progress)
		if err != nil {
			return plansErrorMsg{err: err}
		}

		return progressUpdatedMsg{}
	}
}

// View 渲染界面
func (m *ProgressModel) View() string {
	var formContent strings.Builder

	// 当前进度
	formContent.WriteString(styles.LabelStyle.Render("当前进度"))
	formContent.WriteString("\n")
	progress, _ := strconv.Atoi(m.input.Value())
	formContent.WriteString(utils.FormatProgress(progress, 30))
	formContent.WriteString("\n\n")

	// 输入框
	formContent.WriteString(styles.LabelStyle.Render("新进度 (0-100)"))
	formContent.WriteString("\n")
	formContent.WriteString(m.input.View())
	formContent.WriteString(" %")
	formContent.WriteString("\n")

	// 错误信息
	if m.err != nil {
		formContent.WriteString("\n")
		formContent.WriteString(styles.ErrorStyle.Render("错误: " + m.err.Error()))
	}

	// 使用卡片包装表单
	cardWidth := m.frame.GetContentWidth() - 4
	if cardWidth > 60 {
		cardWidth = 60
	}
	cardContent := components.Card(styles.IconChart+" 更新进度", formContent.String(), cardWidth)

	// 居中显示
	centeredContent := lipgloss.Place(
		m.frame.GetContentWidth(),
		m.frame.GetContentHeight(),
		lipgloss.Center,
		lipgloss.Center,
		cardContent,
	)

	// 快捷键
	keys := []string{
		styles.StatusKeyStyle.Render("↑/↓") + " " + styles.StatusValueStyle.Render("调整"),
		styles.StatusKeyStyle.Render("Enter") + " " + styles.StatusValueStyle.Render("保存"),
		styles.StatusKeyStyle.Render("esc") + " " + styles.StatusValueStyle.Render("取消"),
	}

	return m.frame.Render("计划管理 > "+styles.IconChart+" 更新进度", centeredContent, keys, "")
}
