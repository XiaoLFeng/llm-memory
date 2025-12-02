package group

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/XiaoLFeng/llm-memory/internal/tui/common"
	"github.com/XiaoLFeng/llm-memory/internal/tui/components"
	"github.com/XiaoLFeng/llm-memory/internal/tui/styles"
	"github.com/XiaoLFeng/llm-memory/internal/tui/utils"
	"github.com/XiaoLFeng/llm-memory/pkg/types"
	"github.com/XiaoLFeng/llm-memory/startup"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// DetailModel 组详情模型
// 嘿嘿~ 查看组的详细信息和管理路径！📋
type DetailModel struct {
	bs            *startup.Bootstrap
	groupID       int
	group         *types.Group
	selectedIndex int
	width         int
	height        int
	loading       bool
	err           error
}

// NewDetailModel 创建组详情模型
func NewDetailModel(bs *startup.Bootstrap, groupID int) *DetailModel {
	return &DetailModel{
		bs:      bs,
		groupID: groupID,
		loading: true,
	}
}

// Title 返回页面标题
func (m *DetailModel) Title() string {
	if m.group != nil {
		return "组: " + m.group.Name
	}
	return "组详情"
}

// ShortHelp 返回快捷键帮助
func (m *DetailModel) ShortHelp() []key.Binding {
	return []key.Binding{common.KeyUp, common.KeyDown, common.KeyDelete, common.KeyBack}
}

// Init 初始化
func (m *DetailModel) Init() tea.Cmd {
	return m.loadGroup()
}

// loadGroup 加载组详情
func (m *DetailModel) loadGroup() tea.Cmd {
	return func() tea.Msg {
		group, err := m.bs.GroupService.GetGroup(context.Background(), m.groupID)
		if err != nil {
			return groupDetailErrorMsg{err}
		}
		return groupDetailLoadedMsg{group}
	}
}

type groupDetailLoadedMsg struct {
	group *types.Group
}

type groupDetailErrorMsg struct {
	err error
}

type pathAddedMsg struct{}

type pathRemovedMsg struct {
	path string
}

// Update 处理输入
func (m *DetailModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, common.KeyBack):
			return m, common.Back()

		case key.Matches(msg, common.KeyUp):
			if m.group != nil && m.selectedIndex > 0 {
				m.selectedIndex--
			}

		case key.Matches(msg, common.KeyDown):
			if m.group != nil && m.selectedIndex < len(m.group.Paths)-1 {
				m.selectedIndex++
			}

		case msg.String() == "a":
			// 添加当前路径
			return m, m.addCurrentPath()

		case key.Matches(msg, common.KeyDelete):
			// 删除选中的路径
			if m.group != nil && len(m.group.Paths) > 0 && m.selectedIndex < len(m.group.Paths) {
				path := m.group.Paths[m.selectedIndex]
				return m, common.ShowConfirm(
					"移除路径",
					fmt.Sprintf("确定要从组中移除路径「%s」吗？", path),
					m.removePath(path),
					nil,
				)
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case groupDetailLoadedMsg:
		m.loading = false
		m.group = msg.group
		m.selectedIndex = 0

	case groupDetailErrorMsg:
		m.loading = false
		m.err = msg.err

	case pathAddedMsg:
		cmds = append(cmds, m.loadGroup())
		cmds = append(cmds, common.ShowToast("路径已添加", common.ToastSuccess))

	case pathRemovedMsg:
		cmds = append(cmds, m.loadGroup())
		cmds = append(cmds, common.ShowToast("路径已移除", common.ToastSuccess))

	case common.RefreshMsg:
		m.loading = true
		cmds = append(cmds, m.loadGroup())
	}

	return m, tea.Batch(cmds...)
}

// addCurrentPath 添加当前工作目录到组
func (m *DetailModel) addCurrentPath() tea.Cmd {
	return func() tea.Msg {
		pwd, err := os.Getwd()
		if err != nil {
			return groupDetailErrorMsg{fmt.Errorf("无法获取当前目录: %v", err)}
		}

		err = m.bs.GroupService.AddPath(context.Background(), m.groupID, pwd)
		if err != nil {
			return groupDetailErrorMsg{err}
		}

		return pathAddedMsg{}
	}
}

// removePath 从组中移除路径
func (m *DetailModel) removePath(path string) tea.Cmd {
	return func() tea.Msg {
		err := m.bs.GroupService.RemovePath(context.Background(), m.groupID, path)
		if err != nil {
			return groupDetailErrorMsg{err}
		}
		return pathRemovedMsg{path}
	}
}

// View 渲染界面
func (m *DetailModel) View() string {
	if m.loading {
		loadingStyle := lipgloss.NewStyle().
			Foreground(styles.Info).
			Align(lipgloss.Center)
		content := loadingStyle.Render("加载中...")
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
	}

	if m.err != nil {
		errorStyle := lipgloss.NewStyle().
			Foreground(styles.Error).
			Align(lipgloss.Center)
		content := errorStyle.Render("错误: " + m.err.Error())
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
	}

	if m.group == nil {
		errorStyle := lipgloss.NewStyle().
			Foreground(styles.Error).
			Align(lipgloss.Center)
		content := errorStyle.Render("组不存在")
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
	}

	// 计算卡片宽度
	cardWidth := m.width - 4
	if cardWidth < 60 {
		cardWidth = 60
	}

	// 基本信息卡片
	var basicInfo strings.Builder
	basicInfo.WriteString(components.InfoRow("ID", fmt.Sprintf("%d", m.group.ID)))
	basicInfo.WriteString("\n")
	basicInfo.WriteString(components.InfoRow("名称", m.group.Name))
	basicInfo.WriteString("\n")
	if m.group.Description != "" {
		basicInfo.WriteString(components.InfoRow("描述", m.group.Description))
		basicInfo.WriteString("\n")
	}
	basicInfo.WriteString(components.InfoRow("创建时间", utils.FormatRelativeTime(m.group.CreatedAt)))

	basicCard := components.Card("📋 基本信息", basicInfo.String(), cardWidth)

	// 路径列表卡片
	var pathsList strings.Builder
	if len(m.group.Paths) == 0 {
		pathsList.WriteString(lipgloss.NewStyle().
			Foreground(styles.Subtext0).
			Render("暂无关联路径~ 按 a 添加当前目录"))
	} else {
		for i, path := range m.group.Paths {
			var line string
			if i == m.selectedIndex {
				indicator := lipgloss.NewStyle().Foreground(styles.Primary).Render("▸ ")
				pathStyle := lipgloss.NewStyle().Foreground(styles.Primary).Bold(true)
				line = indicator + pathStyle.Render(path)
			} else {
				line = "  " + lipgloss.NewStyle().Foreground(styles.Text).Render(path)
			}
			pathsList.WriteString(line)
			if i < len(m.group.Paths)-1 {
				pathsList.WriteString("\n")
			}
		}
	}

	pathsTitle := fmt.Sprintf("📂 关联路径 %s",
		lipgloss.NewStyle().Foreground(styles.Subtext0).Render(fmt.Sprintf("(%d)", len(m.group.Paths))))
	pathsCard := components.Card(pathsTitle, pathsList.String(), cardWidth)

	// 组合卡片
	cards := lipgloss.JoinVertical(lipgloss.Left, basicCard, "", pathsCard)

	// 状态栏
	keys := []string{
		lipgloss.NewStyle().Foreground(styles.Primary).Render("↑/↓") + " 选择路径",
		lipgloss.NewStyle().Foreground(styles.Primary).Render("a") + " 添加当前目录",
		lipgloss.NewStyle().Foreground(styles.Primary).Render("d") + " 移除路径",
		lipgloss.NewStyle().Foreground(styles.Primary).Render("esc") + " 返回",
	}
	statusBar := components.RenderKeysOnly(keys, m.width)

	// 组合视图
	contentHeight := m.height - 3
	centeredCards := lipgloss.Place(m.width, contentHeight, lipgloss.Center, lipgloss.Center, cards)

	return lipgloss.JoinVertical(lipgloss.Left, centeredCards, statusBar)
}
