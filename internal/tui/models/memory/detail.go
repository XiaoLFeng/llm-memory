package memory

import (
	"context"
	"strings"

	"github.com/XiaoLFeng/llm-memory/internal/tui/common"
	"github.com/XiaoLFeng/llm-memory/internal/tui/styles"
	"github.com/XiaoLFeng/llm-memory/internal/tui/utils"
	"github.com/XiaoLFeng/llm-memory/pkg/types"
	"github.com/XiaoLFeng/llm-memory/startup"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

// DetailModel 记忆详情模型
// 嘿嘿~ 查看记忆的详细内容！📝
type DetailModel struct {
	bs       *startup.Bootstrap
	id       int
	memory   *types.Memory
	viewport viewport.Model
	ready    bool
	width    int
	height   int
	loading  bool
	err      error
}

// NewDetailModel 创建记忆详情模型
func NewDetailModel(bs *startup.Bootstrap, id int) *DetailModel {
	return &DetailModel{
		bs:      bs,
		id:      id,
		loading: true,
	}
}

// Title 返回页面标题
func (m *DetailModel) Title() string {
	if m.memory != nil {
		return m.memory.Title
	}
	return "记忆详情"
}

// ShortHelp 返回快捷键帮助
func (m *DetailModel) ShortHelp() []key.Binding {
	return []key.Binding{common.KeyUp, common.KeyDown, common.KeyBack}
}

// Init 初始化
func (m *DetailModel) Init() tea.Cmd {
	return m.loadMemory()
}

// loadMemory 加载记忆详情
func (m *DetailModel) loadMemory() tea.Cmd {
	return func() tea.Msg {
		memory, err := m.bs.MemoryService.GetMemory(context.Background(), m.id)
		if err != nil {
			return memoriesErrorMsg{err}
		}
		return memoryLoadedMsg{memory}
	}
}

type memoryLoadedMsg struct {
	memory *types.Memory
}

// Update 处理输入
func (m *DetailModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, common.KeyBack):
			return m, common.Back()
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if !m.ready {
			m.viewport = viewport.New(msg.Width-4, msg.Height-10)
			m.viewport.YPosition = 0
			m.ready = true
		} else {
			m.viewport.Width = msg.Width - 4
			m.viewport.Height = msg.Height - 10
		}
		if m.memory != nil {
			m.viewport.SetContent(m.renderContent())
		}

	case memoryLoadedMsg:
		m.loading = false
		m.memory = msg.memory
		if m.ready {
			m.viewport.SetContent(m.renderContent())
		}

	case memoriesErrorMsg:
		m.loading = false
		m.err = msg.err
	}

	// 更新 viewport
	if m.ready {
		newViewport, cmd := m.viewport.Update(msg)
		m.viewport = newViewport
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// renderContent 渲染内容
func (m *DetailModel) renderContent() string {
	if m.memory == nil {
		return ""
	}

	var b strings.Builder

	// 标题
	b.WriteString(styles.SubtitleStyle.Render("标题"))
	b.WriteString("\n")
	b.WriteString(m.memory.Title)
	b.WriteString("\n\n")

	// 分类
	b.WriteString(styles.SubtitleStyle.Render("分类"))
	b.WriteString("\n")
	b.WriteString(m.memory.Category)
	b.WriteString("\n\n")

	// 优先级
	b.WriteString(styles.SubtitleStyle.Render("优先级"))
	b.WriteString("\n")
	b.WriteString(utils.FormatPriorityIcon(m.memory.Priority) + " " + utils.FormatPriority(m.memory.Priority))
	b.WriteString("\n\n")

	// 标签
	b.WriteString(styles.SubtitleStyle.Render("标签"))
	b.WriteString("\n")
	b.WriteString(utils.JoinTags(m.memory.Tags))
	b.WriteString("\n\n")

	// 创建时间
	b.WriteString(styles.SubtitleStyle.Render("创建时间"))
	b.WriteString("\n")
	b.WriteString(utils.FormatTime(m.memory.CreatedAt))
	b.WriteString("\n\n")

	// 更新时间
	b.WriteString(styles.SubtitleStyle.Render("更新时间"))
	b.WriteString("\n")
	b.WriteString(utils.FormatTime(m.memory.UpdatedAt))
	b.WriteString("\n\n")

	// 内容
	b.WriteString(styles.SubtitleStyle.Render("内容"))
	b.WriteString("\n")
	b.WriteString(m.memory.Content)

	return b.String()
}

// View 渲染界面
func (m *DetailModel) View() string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("📝 记忆详情"))
	b.WriteString("\n\n")

	if m.loading {
		b.WriteString(styles.InfoStyle.Render("加载中..."))
		return b.String()
	}

	if m.err != nil {
		b.WriteString(styles.ErrorStyle.Render("错误: " + m.err.Error()))
		return b.String()
	}

	if m.ready {
		b.WriteString(m.viewport.View())
	}

	b.WriteString("\n\n")
	b.WriteString(styles.HelpStyle.Render("↑/↓ 滚动 | esc 返回"))

	return b.String()
}
