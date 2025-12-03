package plan

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/XiaoLFeng/llm-memory/internal/tui/common"
	"github.com/XiaoLFeng/llm-memory/internal/tui/styles"
	"github.com/XiaoLFeng/llm-memory/internal/tui/utils"
	"github.com/XiaoLFeng/llm-memory/startup"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
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
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render(styles.IconChart + " 更新进度"))
	b.WriteString("\n\n")

	// 当前进度
	b.WriteString(styles.SubtitleStyle.Render("当前进度"))
	b.WriteString("\n")
	progress, _ := strconv.Atoi(m.input.Value())
	b.WriteString(utils.FormatProgress(progress, 30))
	b.WriteString("\n\n")

	// 输入框
	b.WriteString(styles.LabelStyle.Render("新进度 (0-100)"))
	b.WriteString("\n")
	b.WriteString(m.input.View())
	b.WriteString(" %")
	b.WriteString("\n\n")

	// 错误信息
	if m.err != nil {
		b.WriteString(styles.ErrorStyle.Render("错误: " + m.err.Error()))
		b.WriteString("\n\n")
	}

	// 帮助信息
	b.WriteString(styles.HelpStyle.Render("↑/↓ 调整 | enter 保存 | esc 取消"))

	return b.String()
}
