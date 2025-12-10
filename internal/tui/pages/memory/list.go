package memory

import (
	"fmt"
	"strings"
	"time"

	"github.com/XiaoLFeng/llm-memory/internal/tui/components"
	"github.com/XiaoLFeng/llm-memory/internal/tui/core"
	"github.com/XiaoLFeng/llm-memory/internal/tui/layout"
	"github.com/XiaoLFeng/llm-memory/internal/tui/theme"
	"github.com/XiaoLFeng/llm-memory/internal/tui/utils"
	"github.com/XiaoLFeng/llm-memory/startup"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type (
	loadMsg struct {
		items []typesMemory
		err   error
	}
	deleteSuccessMsg struct{}
	deleteErrorMsg   struct{ err error }
)

// typesMemory 只包含 TUI 展示需要的字段，避免直接耦合 entity
type typesMemory struct {
	ID        int64
	Title     string
	Content   string // 记忆内容
	Category  string
	Priority  int
	Global    bool
	PathID    int64
	Tags      []string
	CreatedAt time.Time
}

type ListPage struct {
	bs               *startup.Bootstrap
	frame            *layout.Frame
	width            int
	height           int
	loading          bool
	err              error
	items            []typesMemory
	cursor           int
	showing          bool              // true 展示详情，false 展示列表
	scopeFilter      utils.ScopeFilter // 作用域过滤状态
	detailViewport   viewport.Model    // 详情页滚动视图
	push             func(core.PageID) tea.Cmd
	pushWithData     func(core.PageID, interface{}) tea.Cmd
	confirmDelete    bool  // 是否在删除确认模式
	deleteTarget     int64 // 要删除的 ID
	deleteProcessing bool  // 是否正在处理删除
	deleteYesActive  bool  // true=选中确认，false=选中取消
}

func NewListPage(bs *startup.Bootstrap, push func(core.PageID) tea.Cmd, pushWithData func(core.PageID, interface{}) tea.Cmd) *ListPage {
	// 初始化 viewport（初始尺寸，后续动态调整）
	vp := viewport.New(60, 10)
	vp.Style = lipgloss.NewStyle()

	return &ListPage{
		bs:             bs,
		frame:          layout.NewFrame(80, 24),
		width:          80,
		height:         24,
		loading:        true,
		detailViewport: vp,
		push:           push,
		pushWithData:   pushWithData,
	}
}

func (p *ListPage) Init() tea.Cmd {
	return p.load()
}

func (p *ListPage) load() tea.Cmd {
	return func() tea.Msg {
		ctx := p.bs.Context()
		scopeStr := p.scopeFilter.String()
		memories, err := p.bs.MemoryService.ListMemoriesByScope(ctx, scopeStr, p.bs.CurrentScope)
		if err != nil {
			return loadMsg{err: err}
		}
		items := make([]typesMemory, 0, len(memories))
		for _, m := range memories {
			items = append(items, typesMemory{
				ID:        m.ID,
				Title:     m.Title,
				Content:   m.Content,
				Category:  m.Category,
				Priority:  m.Priority,
				Global:    m.Global,
				PathID:    m.PathID,
				Tags:      m.GetTagStrings(),
				CreatedAt: m.CreatedAt,
			})
		}
		return loadMsg{items: items}
	}
}

func (p *ListPage) Resize(w, h int) {
	p.width, p.height = w, h
	p.frame.Resize(w, h)
}

func (p *ListPage) Update(msg tea.Msg) (core.Page, tea.Cmd) {
	switch v := msg.(type) {
	case tea.KeyMsg:
		// 删除确认模式处理
		if p.confirmDelete {
			switch v.String() {
			case "left", "h", "right", "l":
				p.deleteYesActive = !p.deleteYesActive
				return p, nil
			case "y", "Y":
				p.confirmDelete = false
				p.deleteProcessing = true
				return p, p.doDelete()
			case "n", "N", "esc":
				p.confirmDelete = false
				p.deleteTarget = 0
				return p, nil
			case "enter":
				if p.deleteYesActive {
					p.confirmDelete = false
					p.deleteProcessing = true
					return p, p.doDelete()
				} else {
					p.confirmDelete = false
					p.deleteTarget = 0
					return p, nil
				}
			}
			return p, nil
		}

		// 详情页模式：处理滚动
		if p.showing {
			switch v.String() {
			case "esc", "q":
				p.showing = false
				return p, nil
			case "up", "k":
				p.detailViewport.LineUp(1)
				return p, nil
			case "down", "j":
				p.detailViewport.LineDown(1)
				return p, nil
			case "pgup":
				p.detailViewport.HalfViewUp()
				return p, nil
			case "pgdown":
				p.detailViewport.HalfViewDown()
				return p, nil
			case "home":
				p.detailViewport.GotoTop()
				return p, nil
			case "end":
				p.detailViewport.GotoBottom()
				return p, nil
			}
			return p, nil
		}

		// 列表模式
		switch v.String() {
		case "tab":
			p.scopeFilter = p.scopeFilter.Next()
			p.loading = true
			p.cursor = 0
			return p, p.load()
		case "r":
			p.loading = true
			p.err = nil
			return p, p.load()
		case "up", "k":
			if p.cursor > 0 {
				p.cursor--
			}
		case "down", "j":
			if p.cursor < len(p.items)-1 {
				p.cursor++
			}
		case "enter":
			if len(p.items) > 0 {
				p.showing = !p.showing
				// 进入详情页时重置滚动位置
				if p.showing {
					p.detailViewport.GotoTop()
				}
			}
		case "esc":
			p.showing = false
		case "c":
			return p, p.push(core.PageMemoryCreate)
		case "e":
			if len(p.items) > 0 {
				id := p.items[p.cursor].ID
				return p, p.pushWithData(core.PageMemoryEdit, id)
			}
		case "d":
			if len(p.items) > 0 {
				p.deleteTarget = p.items[p.cursor].ID
				p.confirmDelete = true
				p.deleteYesActive = false // 默认选中"取消"，更安全
			}
		case "?":
			return p, p.push(core.PageHelp)
		}
	case loadMsg:
		p.loading = false
		p.err = v.err
		if v.err == nil {
			p.items = v.items
			if p.cursor >= len(p.items) {
				p.cursor = len(p.items) - 1
			}
			if p.cursor < 0 {
				p.cursor = 0
			}
		}
	case deleteSuccessMsg:
		p.deleteProcessing = false
		p.deleteTarget = 0
		p.loading = true
		return p, p.load()
	case deleteErrorMsg:
		p.deleteProcessing = false
		p.deleteTarget = 0
		p.err = v.err
	case tea.WindowSizeMsg:
		// 动态调整 viewport 尺寸
		if p.showing {
			// 直接使用终端尺寸，减去详情视图自身的 header/footer
			const detailOverhead = 4
			p.detailViewport.Width = v.Width - 4
			p.detailViewport.Height = v.Height - detailOverhead
		}
	}
	return p, nil
}

func (p *ListPage) View() string {
	cw, _ := p.frame.ContentSize()
	cardWidth := layout.FitCardWidth(cw)
	scopeLabel := p.scopeFilter.Label()
	titleWithScope := fmt.Sprintf("%s 记忆列表 [%s]", theme.IconMemory, scopeLabel)

	// 删除确认对话框
	if p.confirmDelete {
		var itemName string
		if p.cursor < len(p.items) {
			itemName = p.items[p.cursor].Title
		}
		return components.ConfirmDialogWithButtons("确认删除",
			fmt.Sprintf("确定要删除「%s」吗？\n此操作不可撤销。", itemName),
			cardWidth, p.deleteYesActive)
	}

	switch {
	case p.loading || p.deleteProcessing:
		msg := "努力加载中..."
		if p.deleteProcessing {
			msg = "正在删除..."
		}
		return components.LoadingState(titleWithScope, msg, cardWidth)
	case p.err != nil:
		return components.ErrorState(titleWithScope, p.err.Error(), cardWidth)
	case len(p.items) == 0:
		return components.EmptyState(titleWithScope, "暂无记忆，按 c 创建一条吧~", cardWidth)
	default:
		if p.showing {
			// === 使用 viewport 渲染详情页 ===
			// 直接使用终端尺寸，减去详情视图自身的 header/footer
			// title(1) + 空行(1) + 空行(1) + scrollHint(1) = 4行
			const detailOverhead = 4

			viewportWidth := p.width - 4
			viewportHeight := p.height - detailOverhead

			p.detailViewport.Width = viewportWidth
			p.detailViewport.Height = viewportHeight

			// 生成详情内容并设置到 viewport
			detailContent := p.renderDetail(p.detailViewport.Width)
			p.detailViewport.SetContent(detailContent)

			// 滚动进度指示器
			scrollPercent := p.detailViewport.ScrollPercent() * 100
			scrollInfo := fmt.Sprintf("%.0f%%", scrollPercent)
			scrollHint := theme.TextDim.Render(fmt.Sprintf(
				"滚动: %s | ↑/↓ j/k PgUp/PgDn Home/End | Esc 返回", scrollInfo))

			// 组合视图
			title := theme.Title.Render(theme.IconMemory + " 记忆详情")
			viewportView := p.detailViewport.View()

			return lipgloss.JoinVertical(lipgloss.Left,
				title,
				"",
				viewportView,
				"",
				scrollHint,
			)
		}
		body := p.renderList(cardWidth - 6)
		return components.Card(titleWithScope, body, cardWidth)
	}
}

func (p *ListPage) renderList(width int) string {
	var b strings.Builder
	max := len(p.items)
	if max > 20 {
		max = 20
	}
	for i := 0; i < max; i++ {
		m := p.items[i]
		scope := utils.ScopeTagWithGlobal(m.Global, m.PathID, p.bs)
		tagStr := ""
		if len(m.Tags) > 0 {
			tagStr = " #" + strings.Join(m.Tags, " #")
		}
		line := fmt.Sprintf("%s %s · %s · P%d · %s%s",
			scope, m.Title, m.Category, m.Priority,
			m.CreatedAt.Format("01-02 15:04"), tagStr)
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

	m := p.items[p.cursor]
	scope := utils.ScopeTagWithGlobal(m.Global, m.PathID, p.bs)

	var lines []string

	// === 区块 1：标题 ===
	titleLine := theme.FormLabel.Bold(true).Render("标题: ") + theme.TextMain.Render(m.Title)
	lines = append(lines, titleLine)
	lines = append(lines, "")

	// === 区块 2：元数据 ===
	metaStyle := theme.TextDim
	tagStr := "无"
	if len(m.Tags) > 0 {
		tagStr = strings.Join(m.Tags, ", ")
	}
	lines = append(lines, metaStyle.Render(fmt.Sprintf(
		"分类: %s | 优先级: P%d | 作用域: %s",
		m.Category, m.Priority, scope)))
	lines = append(lines, metaStyle.Render(fmt.Sprintf(
		"标签: %s", tagStr)))
	lines = append(lines, metaStyle.Render(fmt.Sprintf(
		"创建时间: %s", m.CreatedAt.Format("2006-01-02 15:04:05"))))

	// === 分隔线 ===
	lines = append(lines, "")
	separatorLine := lipgloss.NewStyle().
		Foreground(theme.Border).
		Render(strings.Repeat("─", width))
	lines = append(lines, separatorLine)

	// === 区块 3：内容 ===
	if m.Content != "" {
		lines = append(lines, "")
		contentLines := utils.RenderDetailSection("📄", "内容", m.Content, width)
		lines = append(lines, contentLines...)
	}

	return strings.Join(lines, "\n")
}

func (p *ListPage) Meta() core.Meta {
	// 详情页模式
	if p.showing {
		return core.Meta{
			Title:      "记忆详情",
			Breadcrumb: "记忆管理 > 详情",
			Keys: []components.KeyHint{
				{Key: "↑/↓ j/k", Desc: "滚动"},
				{Key: "PgUp/PgDn", Desc: "翻页"},
				{Key: "Home/End", Desc: "首/尾"},
				{Key: "Esc", Desc: "返回列表"},
			},
		}
	}

	// 列表模式
	return core.Meta{
		Title:      "记忆列表",
		Breadcrumb: "记忆管理 > 列表",
		Extra:      fmt.Sprintf("[%s] Tab切换 r刷新", p.scopeFilter.Label()),
		Keys: []components.KeyHint{
			{Key: "Tab", Desc: "切换作用域"},
			{Key: "Enter", Desc: "详情"},
			{Key: "c", Desc: "新建"},
			{Key: "e", Desc: "编辑"},
			{Key: "d", Desc: "删除"},
			{Key: "r", Desc: "刷新"},
			{Key: "Esc", Desc: "返回"},
			{Key: "↑/↓", Desc: "移动"},
		},
	}
}

// doDelete 执行删除操作
func (p *ListPage) doDelete() tea.Cmd {
	return func() tea.Msg {
		ctx := p.bs.Context()
		if err := p.bs.MemoryService.DeleteMemoryByID(ctx, p.deleteTarget); err != nil {
			return deleteErrorMsg{err: err}
		}
		return deleteSuccessMsg{}
	}
}
