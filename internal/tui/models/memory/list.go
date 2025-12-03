package memory

import (
	"context"
	"fmt"
	"strings"

	"github.com/XiaoLFeng/llm-memory/internal/models/entity"
	"github.com/XiaoLFeng/llm-memory/internal/tui/common"
	"github.com/XiaoLFeng/llm-memory/internal/tui/components"
	"github.com/XiaoLFeng/llm-memory/internal/tui/styles"
	"github.com/XiaoLFeng/llm-memory/internal/tui/utils"
	"github.com/XiaoLFeng/llm-memory/startup"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ListModel 记忆列表模型
// 嘿嘿~ 展示所有记忆的列表！📚
type ListModel struct {
	bs       *startup.Bootstrap
	memories []entity.Memory
	selected int
	frame    *components.Frame
	width    int
	height   int
	loading  bool
	err      error
}

// NewListModel 创建记忆列表模型
func NewListModel(bs *startup.Bootstrap) *ListModel {
	return &ListModel{
		bs:      bs,
		frame:   components.NewFrame(80, 24),
		loading: true,
	}
}

// Title 返回页面标题
func (m *ListModel) Title() string {
	return "记忆列表"
}

// ShortHelp 返回快捷键帮助
func (m *ListModel) ShortHelp() []key.Binding {
	return []key.Binding{
		common.KeyUp, common.KeyDown, common.KeyEnter,
		common.KeyCreate, common.KeyDelete, common.KeySearch, common.KeyBack,
	}
}

// Init 初始化
func (m *ListModel) Init() tea.Cmd {
	return m.loadMemories()
}

// loadMemories 加载记忆列表
func (m *ListModel) loadMemories() tea.Cmd {
	return func() tea.Msg {
		memories, err := m.bs.MemoryService.ListMemories(context.Background())
		if err != nil {
			return memoriesErrorMsg{err}
		}
		return memoriesLoadedMsg{memories}
	}
}

type memoriesLoadedMsg struct {
	memories []entity.Memory
}

type memoriesErrorMsg struct {
	err error
}

type memoryDeletedMsg struct {
	id uint
}

// Update 处理输入
func (m *ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, common.KeyBack):
			return m, common.Back()

		case key.Matches(msg, common.KeyCreate):
			return m, common.Navigate(common.PageMemoryCreate)

		case key.Matches(msg, common.KeySearch):
			return m, common.Navigate(common.PageMemorySearch)

		case key.Matches(msg, common.KeyUp):
			if m.selected > 0 {
				m.selected--
			}

		case key.Matches(msg, common.KeyDown):
			if m.selected < len(m.memories)-1 {
				m.selected++
			}

		case key.Matches(msg, common.KeyEnter):
			if len(m.memories) > 0 && m.selected < len(m.memories) {
				return m, common.Navigate(common.PageMemoryDetail, map[string]any{"id": m.memories[m.selected].ID})
			}

		case key.Matches(msg, common.KeyDelete):
			if len(m.memories) > 0 && m.selected < len(m.memories) {
				return m, common.ShowConfirm(
					"删除记忆",
					fmt.Sprintf("确定要删除记忆「%s」吗？", m.memories[m.selected].Title),
					m.deleteMemory(m.memories[m.selected].ID),
					nil,
				)
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.frame.SetSize(msg.Width, msg.Height)

	case memoriesLoadedMsg:
		m.loading = false
		m.memories = msg.memories
		// 确保选中项在范围内
		if m.selected >= len(m.memories) {
			m.selected = len(m.memories) - 1
		}
		if m.selected < 0 {
			m.selected = 0
		}

	case memoriesErrorMsg:
		m.loading = false
		m.err = msg.err

	case memoryDeletedMsg:
		cmds = append(cmds, m.loadMemories())
		cmds = append(cmds, common.ShowToast("记忆已删除", common.ToastSuccess))

	case common.RefreshMsg:
		m.loading = true
		cmds = append(cmds, m.loadMemories())
	}

	return m, tea.Batch(cmds...)
}

// deleteMemory 删除记忆
func (m *ListModel) deleteMemory(id uint) tea.Cmd {
	return func() tea.Msg {
		err := m.bs.MemoryService.DeleteMemory(context.Background(), id)
		if err != nil {
			return memoriesErrorMsg{err}
		}
		return memoryDeletedMsg{id}
	}
}

// View 渲染界面
func (m *ListModel) View() string {
	// 加载中
	if m.loading {
		loadingContent := lipgloss.NewStyle().
			Foreground(styles.Info).
			Render("加载中...")
		return m.frame.Render("记忆管理 > 记忆列表", loadingContent, []string{}, "")
	}

	// 错误
	if m.err != nil {
		errorContent := lipgloss.NewStyle().
			Foreground(styles.Error).
			Render("错误: " + m.err.Error())
		return m.frame.Render("记忆管理 > 记忆列表", errorContent, []string{}, "")
	}

	// 空列表
	if len(m.memories) == 0 {
		emptyContent := lipgloss.NewStyle().
			Foreground(styles.Overlay0).
			Render("暂无记忆~ 按 c 创建新记忆")
		keys := []string{
			styles.StatusKeyStyle.Render("c") + " " + styles.StatusValueStyle.Render("新建"),
			styles.StatusKeyStyle.Render("/") + " " + styles.StatusValueStyle.Render("搜索"),
			styles.StatusKeyStyle.Render("esc") + " " + styles.StatusValueStyle.Render("返回"),
		}
		return m.frame.Render("记忆管理 > 记忆列表", emptyContent, keys, "")
	}

	// 渲染列表
	var listItems strings.Builder
	for i, memory := range m.memories {
		// 选中指示器
		indicator := "  "
		if i == m.selected {
			indicator = lipgloss.NewStyle().Foreground(styles.Primary).Render("▸ ")
		} else {
			indicator = "  "
		}

		// 标题样式
		titleStyle := styles.ListItemTitleStyle
		if i == m.selected {
			titleStyle = styles.ListItemTitleSelectedStyle
		}

		// 构建元信息
		var meta []string

		// 作用域徽章
		scopeBadge := components.ScopeBadgeFromGroupIDPath(memory.GroupID, memory.Path)
		meta = append(meta, scopeBadge)

		// 分类
		categoryBadge := components.CategoryBadge(memory.Category)
		meta = append(meta, categoryBadge)

		// 优先级
		priorityBadge := components.PriorityBadgeSimple(memory.Priority)
		meta = append(meta, priorityBadge)

		// 标签
		if len(memory.Tags) > 0 {
			// 转换 []entity.MemoryTag 为 []string
			tags := make([]string, len(memory.Tags))
			for i, t := range memory.Tags {
				tags[i] = t.Tag
			}
			tagsBadge := components.TagsBadge(tags)
			meta = append(meta, tagsBadge)
		}

		// 时间
		timeStr := utils.FormatRelativeTime(memory.CreatedAt)
		timeBadge := components.TimeBadge(timeStr)
		meta = append(meta, timeBadge)

		metaStr := strings.Join(meta, styles.MetaSeparator)

		// 描述样式
		descStyle := styles.ListItemDescStyle
		if i == m.selected {
			descStyle = styles.ListItemDescSelectedStyle
		}

		// 渲染列表项
		title := fmt.Sprintf("%s%s", indicator, titleStyle.Render(memory.Title))
		desc := "    " + descStyle.Render(metaStr)

		listItems.WriteString(title)
		listItems.WriteString("\n")
		listItems.WriteString(desc)

		if i < len(m.memories)-1 {
			listItems.WriteString("\n\n")
		}
	}

	// 快捷键
	keys := []string{
		styles.StatusKeyStyle.Render("↑/↓") + " " + styles.StatusValueStyle.Render("选择"),
		styles.StatusKeyStyle.Render("Enter") + " " + styles.StatusValueStyle.Render("查看"),
		styles.StatusKeyStyle.Render("c") + " " + styles.StatusValueStyle.Render("新建"),
		styles.StatusKeyStyle.Render("d") + " " + styles.StatusValueStyle.Render("删除"),
		styles.StatusKeyStyle.Render("/") + " " + styles.StatusValueStyle.Render("搜索"),
		styles.StatusKeyStyle.Render("esc") + " " + styles.StatusValueStyle.Render("返回"),
	}

	// 额外信息：总数
	extra := fmt.Sprintf("共 %d 条", len(m.memories))

	return m.frame.Render("记忆管理 > 记忆列表", listItems.String(), keys, extra)
}
