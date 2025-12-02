package handlers

import (
	"context"
	"fmt"
	"strings"

	"github.com/XiaoLFeng/llm-memory/internal/cli"
	"github.com/XiaoLFeng/llm-memory/internal/cli/output"
	"github.com/XiaoLFeng/llm-memory/pkg/types"
	"github.com/XiaoLFeng/llm-memory/startup"
)

// MemoryHandler 记忆命令处理器
// 嘿嘿~ 处理所有记忆相关的 CLI 命令！(´∀｀)💖
type MemoryHandler struct {
	bs *startup.Bootstrap
}

// NewMemoryHandler 创建记忆处理器
func NewMemoryHandler(bs *startup.Bootstrap) *MemoryHandler {
	return &MemoryHandler{bs: bs}
}

// List 列出所有记忆
// 呀~ 展示所有记忆条目！✨
func (h *MemoryHandler) List(ctx context.Context) error {
	memories, err := h.bs.MemoryService.ListMemories(ctx)
	if err != nil {
		return err
	}

	if len(memories) == 0 {
		cli.PrintInfo("暂无记忆~ 快创建一条吧！")
		return nil
	}

	cli.PrintTitle("📚 记忆列表")
	table := output.NewTable("ID", "标题", "分类", "创建时间")
	for _, m := range memories {
		table.AddRow(
			fmt.Sprintf("%d", m.ID),
			m.Title,
			m.Category,
			m.CreatedAt.Format("2006-01-02 15:04"),
		)
	}
	table.Print()

	return nil
}

// Create 创建记忆
// 嘿嘿~ 创建新的记忆条目！💫
func (h *MemoryHandler) Create(ctx context.Context, title, content, category string, tags []string) error {
	if category == "" {
		category = "默认"
	}

	memory, err := h.bs.MemoryService.CreateMemory(ctx, title, content, category, tags, 2, types.GlobalGroupID, "")
	if err != nil {
		return err
	}

	cli.PrintSuccess(fmt.Sprintf("记忆创建成功！ID: %d, 标题: %s", memory.ID, memory.Title))
	return nil
}

// Search 搜索记忆
// 呀~ 根据关键词搜索记忆！🔍
func (h *MemoryHandler) Search(ctx context.Context, keyword string) error {
	memories, err := h.bs.MemoryService.SearchMemories(ctx, keyword)
	if err != nil {
		return err
	}

	if len(memories) == 0 {
		cli.PrintInfo(fmt.Sprintf("未找到包含 \"%s\" 的记忆~", keyword))
		return nil
	}

	cli.PrintTitle(fmt.Sprintf("🔍 搜索结果 (%d 条)", len(memories)))
	table := output.NewTable("ID", "标题", "分类")
	for _, m := range memories {
		table.AddRow(
			fmt.Sprintf("%d", m.ID),
			m.Title,
			m.Category,
		)
	}
	table.Print()

	return nil
}

// Delete 删除记忆
func (h *MemoryHandler) Delete(ctx context.Context, id int) error {
	if err := h.bs.MemoryService.DeleteMemory(ctx, id); err != nil {
		return err
	}

	cli.PrintSuccess(fmt.Sprintf("记忆 %d 已删除", id))
	return nil
}

// Get 获取单个记忆详情
// 嗯嗯！查看记忆的详细内容！📝
func (h *MemoryHandler) Get(ctx context.Context, id int) error {
	memory, err := h.bs.MemoryService.GetMemory(ctx, id)
	if err != nil {
		return err
	}

	cli.PrintTitle("📝 记忆详情")
	fmt.Printf("ID:       %d\n", memory.ID)
	fmt.Printf("标题:     %s\n", memory.Title)
	fmt.Printf("分类:     %s\n", memory.Category)
	if len(memory.Tags) > 0 {
		fmt.Printf("标签:     %s\n", strings.Join(memory.Tags, ", "))
	}
	fmt.Printf("优先级:   %d\n", memory.Priority)
	fmt.Printf("创建时间: %s\n", memory.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("更新时间: %s\n", memory.UpdatedAt.Format("2006-01-02 15:04:05"))
	fmt.Println("\n内容:")
	fmt.Println(memory.Content)

	return nil
}
