package memory

import (
	"context"
	"strings"

	"github.com/XiaoLFeng/llm-memory/internal/models/entity"
	"github.com/XiaoLFeng/llm-memory/internal/tui/common"
	"github.com/XiaoLFeng/llm-memory/internal/tui/components"
	"github.com/XiaoLFeng/llm-memory/internal/tui/styles"
	"github.com/XiaoLFeng/llm-memory/internal/tui/utils"
	"github.com/XiaoLFeng/llm-memory/startup"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// DetailModel 记忆详情模型
type DetailModel struct {
	bs       *startup.Bootstrap
	id       int64
	memory   *entity.Memory
	viewport viewport.Model
	frame    *components.Frame
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
		id:      int64(id),
		frame:   components.NewFrame(80, 24),
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
	memory *entity.Memory
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
		m.frame.SetSize(msg.Width, msg.Height)

		// 直接使用 frame 的内容尺寸，不再额外减法
		contentHeight := m.frame.GetContentHeight()
		contentWidth := m.frame.GetContentWidth()

		if !m.ready {
			m.viewport = viewport.New(contentWidth, contentHeight)
			m.viewport.YPosition = 0
			m.ready = true
		} else {
			m.viewport.Width = contentWidth
			m.viewport.Height = contentHeight
		}
		// 无论数据是否已加载，都尝试更新内容
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

	// 直接使用 viewport 宽度，减去卡片边框和内边距
	cardWidth := m.viewport.Width - 4
	if cardWidth < 40 {
		cardWidth = 40
	}

	// 基本信息卡片
	var basicInfo strings.Builder
	basicInfo.WriteString(components.InfoRow("标题", m.memory.Title))
	basicInfo.WriteString("\n")
	basicInfo.WriteString(components.InfoRow("分类", components.CategoryBadge(m.memory.Category)))
	basicInfo.WriteString("\n")
	basicInfo.WriteString(components.InfoRow("优先级", components.PriorityBadge(m.memory.Priority)))
	basicInfo.WriteString("\n")
	basicInfo.WriteString(components.InfoRow("作用域", components.ScopeBadgeFromPathID(m.memory.PathID)))
	basicInfo.WriteString("\n")
	if len(m.memory.Tags) > 0 {
		// 转换 []entity.MemoryTag 为 []string
		tags := make([]string, len(m.memory.Tags))
		for i, t := range m.memory.Tags {
			tags[i] = t.Tag
		}
		basicInfo.WriteString(components.InfoRow("标签", components.TagsBadge(tags)))
		basicInfo.WriteString("\n")
	}
	basicInfo.WriteString(components.InfoRow("创建时间", utils.FormatTime(m.memory.CreatedAt)))
	basicInfo.WriteString("\n")
	basicInfo.WriteString(components.InfoRow("更新时间", utils.FormatTime(m.memory.UpdatedAt)))

	basicCard := components.NestedCard("基本信息", basicInfo.String(), cardWidth)

	// 内容卡片
	contentStyle := lipgloss.NewStyle().
		Foreground(styles.Text)
	contentCard := components.NestedCard("记忆内容", contentStyle.Render(m.memory.Content), cardWidth)

	// 组合所有卡片
	return lipgloss.JoinVertical(
		lipgloss.Left,
		basicCard,
		"",
		contentCard,
	)
}

// View 渲染界面
func (m *DetailModel) View() string {
	// 加载中
	if m.loading {
		loadingContent := lipgloss.NewStyle().
			Foreground(styles.Info).
			Render("加载中...")
		return m.frame.Render("记忆管理 > 记忆详情", loadingContent, []string{}, "")
	}

	// 错误
	if m.err != nil {
		errorContent := lipgloss.NewStyle().
			Foreground(styles.Error).
			Render("错误: " + m.err.Error())
		return m.frame.Render("记忆管理 > 记忆详情", errorContent, []string{}, "")
	}

	// 内容
	var content string
	if m.ready {
		content = m.viewport.View()
	}

	// 快捷键
	keys := []string{
		styles.StatusKeyStyle.Render("↑/↓") + " " + styles.StatusValueStyle.Render("滚动"),
		styles.StatusKeyStyle.Render("esc") + " " + styles.StatusValueStyle.Render("返回"),
	}

	// 面包屑
	breadcrumb := "记忆管理 > 记忆详情"
	if m.memory != nil {
		breadcrumb = "记忆管理 > " + m.memory.Title
	}

	return m.frame.Render(breadcrumb, content, keys, "")
}
