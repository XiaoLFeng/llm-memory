package group

import (
	"fmt"
	"strings"

	"github.com/XiaoLFeng/llm-memory/internal/tui/components"
	"github.com/XiaoLFeng/llm-memory/internal/tui/core"
	"github.com/XiaoLFeng/llm-memory/internal/tui/layout"
	"github.com/XiaoLFeng/llm-memory/internal/tui/theme"
	"github.com/XiaoLFeng/llm-memory/internal/tui/utils"
	"github.com/XiaoLFeng/llm-memory/startup"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type loadMsg struct {
	items []groupItem
	err   error
}

type groupItem struct {
	ID          int64
	Name        string
	Description string
	PathCount   int
}

type ListPage struct {
	bs               *startup.Bootstrap
	frame            *layout.Frame
	navigateWithData func(core.PageID, interface{}) tea.Cmd
	loading          bool
	err              error
	items            []groupItem
	cursor           int
	showing          bool
	confirmDelete    bool
	deleteTarget     string
}

func NewListPage(bs *startup.Bootstrap, _ func(core.PageID) tea.Cmd, navigateWithData func(core.PageID, interface{}) tea.Cmd) *ListPage {
	return &ListPage{
		bs:               bs,
		frame:            layout.NewFrame(80, 24),
		navigateWithData: navigateWithData,
		loading:          true,
	}
}

func (p *ListPage) Init() tea.Cmd { return p.load() }

func (p *ListPage) load() tea.Cmd {
	return func() tea.Msg {
		ctx := p.bs.Context()
		groups, err := p.bs.GroupService.ListGroups(ctx)
		if err != nil {
			return loadMsg{err: err}
		}
		items := make([]groupItem, 0, len(groups))
		for _, g := range groups {
			items = append(items, groupItem{
				ID:          g.ID,
				Name:        g.Name,
				Description: g.Description,
				PathCount:   len(g.Paths),
			})
		}
		return loadMsg{items: items}
	}
}

func (p *ListPage) Resize(w, h int) { p.frame.Resize(w, h) }

func (p *ListPage) Update(msg tea.Msg) (core.Page, tea.Cmd) {
	switch v := msg.(type) {
	case tea.KeyMsg:
		// 删除确认对话框激活状态
		if p.confirmDelete {
			switch v.String() {
			case "y", "Y":
				return p, p.doDelete()
			case "n", "N", "esc":
				p.confirmDelete = false
				p.deleteTarget = ""
			}
			return p, nil
		}

		// 正常按键处理
		if v.String() == "r" {
			p.loading = true
			p.err = nil
			return p, p.load()
		}
		switch v.String() {
		case "up", "k":
			if p.cursor > 0 {
				p.cursor--
			}
		case "down", "j":
			if p.cursor < len(p.items)-1 {
				p.cursor++
			}
		case "enter":
			p.showing = !p.showing
		case "esc":
			p.showing = false
		case "c":
			if p.navigateWithData != nil {
				return p, p.navigateWithData(core.PageGroupCreate, nil)
			}
		case "e":
			if len(p.items) > 0 && p.navigateWithData != nil {
				item := p.items[p.cursor]
				return p, p.navigateWithData(core.PageGroupEdit, item.ID)
			}
		case "d":
			if len(p.items) > 0 {
				p.confirmDelete = true
				p.deleteTarget = p.items[p.cursor].Name
			}
		case "?":
			// 切换帮助显示（暂未实现）
		}
	case loadMsg:
		p.loading = false
		p.err = v.err
		if v.err == nil {
			p.items = v.items
			if p.cursor >= len(p.items) {
				p.cursor = len(p.items) - 1
			}
		}
	case deleteResult:
		p.confirmDelete = false
		p.deleteTarget = ""
		if v.err != nil {
			p.err = v.err
		} else {
			// 删除成功，重新加载列表
			p.loading = true
			return p, p.load()
		}
	}
	return p, nil
}

func (p *ListPage) View() string {
	cw, _ := p.frame.ContentSize()
	cardW := layout.FitCardWidth(cw)

	// 删除确认对话框覆盖层
	if p.confirmDelete {
		return components.DeleteConfirmDialog(p.deleteTarget, cardW)
	}

	switch {
	case p.loading:
		return components.LoadingState(theme.IconGroup+" 组列表", "加载组信息中...", cardW)
	case p.err != nil:
		return components.ErrorState(theme.IconGroup+" 组列表", p.err.Error(), cardW)
	case len(p.items) == 0:
		return components.EmptyState(theme.IconGroup+" 组列表", "暂无组，按 c 创建吧~", cardW)
	default:
		if p.showing {
			body := p.renderDetail(cardW - 6)
			return components.Card(theme.IconGroup+" 组详情", body, cardW)
		}
		body := p.renderList(cardW - 6)
		return components.Card(theme.IconGroup+" 组列表", body, cardW)
	}
}

func (p *ListPage) renderList(width int) string {
	var b strings.Builder
	max := len(p.items)
	if max > 20 {
		max = 20
	}
	for i := 0; i < max; i++ {
		g := p.items[i]
		desc := g.Description
		if desc == "" {
			desc = "暂无描述"
		}
		line := fmt.Sprintf("%s · 路径 %d · %s", g.Name, g.PathCount, desc)
		if utils.LipWidth(line) > width {
			line = utils.Truncate(line, width)
		}
		if i == p.cursor {
			line = lipgloss.NewStyle().Foreground(theme.Info).Render("▶ " + line)
		} else {
			line = "  " + line
		}
		b.WriteString(line)
		if i != max-1 {
			b.WriteRune('\n')
		}
	}
	return b.String()
}

func (p *ListPage) renderDetail(width int) string {
	if len(p.items) == 0 {
		return "暂无数据"
	}
	g := p.items[p.cursor]
	desc := g.Description
	if desc == "" {
		desc = "暂无描述"
	}
	lines := []string{
		fmt.Sprintf("名称: %s", g.Name),
		fmt.Sprintf("路径数量: %d", g.PathCount),
		fmt.Sprintf("描述: %s", desc),
	}
	for i, l := range lines {
		if utils.LipWidth(l) > width {
			lines[i] = utils.Truncate(l, width)
		}
	}
	return strings.Join(lines, "\n")
}

func (p *ListPage) Meta() core.Meta {
	return core.Meta{
		Title:      "组列表",
		Breadcrumb: "组管理 > 列表",
		Extra:      "r 刷新",
		Keys: []components.KeyHint{
			{Key: "Enter", Desc: "切换详情/列表"},
			{Key: "c", Desc: "新建"},
			{Key: "e", Desc: "编辑"},
			{Key: "d", Desc: "删除"},
			{Key: "r", Desc: "刷新"},
			{Key: "?", Desc: "帮助"},
			{Key: "Esc", Desc: "返回"},
			{Key: "↑/↓", Desc: "移动"},
		},
	}
}

// deleteResult 删除结果消息
type deleteResult struct {
	err error
}

// doDelete 执行删除操作
func (p *ListPage) doDelete() tea.Cmd {
	return func() tea.Msg {
		if len(p.items) == 0 {
			return deleteResult{err: nil}
		}
		item := p.items[p.cursor]
		ctx := p.bs.Context()
		err := p.bs.GroupService.DeleteGroup(ctx, item.ID)
		return deleteResult{err: err}
	}
}
