package handlers

import (
	"context"
	"fmt"
	"strings"

	"github.com/XiaoLFeng/llm-memory/internal/cli"
	"github.com/XiaoLFeng/llm-memory/internal/cli/output"
	"github.com/XiaoLFeng/llm-memory/internal/models/dto"
	"github.com/XiaoLFeng/llm-memory/internal/models/entity"
	"github.com/XiaoLFeng/llm-memory/startup"
)

// MemoryHandler 记忆命令处理器
type MemoryHandler struct {
	bs *startup.Bootstrap
}

// NewMemoryHandler 创建记忆处理器
func NewMemoryHandler(bs *startup.Bootstrap) *MemoryHandler {
	return &MemoryHandler{bs: bs}
}

// List 列出所有记忆
func (h *MemoryHandler) List(ctx context.Context) error {
	// 使用 ListMemoriesByScope 确保权限隔离：全局 + 当前路径相关
	memories, err := h.bs.MemoryService.ListMemoriesByScope(ctx, "all", h.bs.CurrentScope)
	if err != nil {
		return err
	}

	if len(memories) == 0 {
		cli.PrintInfo("暂无记忆~ 快创建一条吧！")
		return nil
	}

	cli.PrintTitle(cli.IconMemory + " 记忆列表")
	table := output.NewTable("标识码", "标题", "分类", "创建时间")
	for _, m := range memories {
		table.AddRow(
			m.Code,
			m.Title,
			m.Category,
			m.CreatedAt.Format("2006-01-02 15:04"),
		)
	}
	table.Print()

	return nil
}

// Create 创建记忆
func (h *MemoryHandler) Create(ctx context.Context, code, title, content, category string, tags []string, global bool) error {
	if category == "" {
		category = "默认"
	}

	createDTO := &dto.MemoryCreateDTO{
		Code:     code,
		Title:    title,
		Content:  content,
		Category: category,
		Tags:     tags,
		Priority: 2,
		Global:   global,
	}
	memory, err := h.bs.MemoryService.CreateMemory(ctx, createDTO, h.bs.CurrentScope)
	if err != nil {
		return err
	}

	cli.PrintSuccess(fmt.Sprintf("记忆创建成功！标识码: %s, 标题: %s", memory.Code, memory.Title))
	return nil
}

// Search 搜索记忆
func (h *MemoryHandler) Search(ctx context.Context, keyword string) error {
	// 使用 SearchMemoriesByScope 确保权限隔离：全局 + 当前路径相关
	memories, err := h.bs.MemoryService.SearchMemoriesByScope(ctx, keyword, "all", h.bs.CurrentScope)
	if err != nil {
		return err
	}

	if len(memories) == 0 {
		cli.PrintInfo(fmt.Sprintf("未找到包含 \"%s\" 的记忆~", keyword))
		return nil
	}

	cli.PrintTitle(fmt.Sprintf("%s 搜索结果 (%d 条)", cli.IconSearch, len(memories)))
	table := output.NewTable("标识码", "标题", "分类")
	for _, m := range memories {
		table.AddRow(
			m.Code,
			m.Title,
			m.Category,
		)
	}
	table.Print()

	return nil
}

// Delete 删除记忆
func (h *MemoryHandler) Delete(ctx context.Context, code string) error {
	if err := h.bs.MemoryService.DeleteMemory(ctx, code); err != nil {
		return err
	}

	cli.PrintSuccess(fmt.Sprintf("记忆 %s 已删除", code))
	return nil
}

// Get 获取单个记忆详情
func (h *MemoryHandler) Get(ctx context.Context, code string) error {
	memory, err := h.bs.MemoryService.GetMemory(ctx, code)
	if err != nil {
		return err
	}

	cli.PrintTitle(cli.IconEdit + " 记忆详情")
	fmt.Printf("标识码:   %s\n", memory.Code)
	fmt.Printf("标题:     %s\n", memory.Title)
	fmt.Printf("分类:     %s\n", memory.Category)
	if len(memory.Tags) > 0 {
		tags := make([]string, len(memory.Tags))
		for i := range memory.Tags {
			var tag entity.MemoryTag = memory.Tags[i]
			tags[i] = tag.Tag
		}
		fmt.Printf("标签:     %s\n", strings.Join(tags, ", "))
	}
	fmt.Printf("优先级:   %d\n", memory.Priority)
	fmt.Printf("创建时间: %s\n", memory.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("更新时间: %s\n", memory.UpdatedAt.Format("2006-01-02 15:04:05"))
	fmt.Println("\n内容:")
	fmt.Println(memory.Content)

	return nil
}

// Update 更新记忆
func (h *MemoryHandler) Update(ctx context.Context, code string, title, content, category *string, tags *[]string, priority *int) error {
	updateDTO := &dto.MemoryUpdateDTO{
		Code:     code,
		Title:    title,
		Content:  content,
		Category: category,
		Tags:     tags,
		Priority: priority,
	}

	if err := h.bs.MemoryService.UpdateMemory(ctx, updateDTO); err != nil {
		return err
	}

	// 构建更新信息
	var updated []string
	if title != nil {
		updated = append(updated, "标题")
	}
	if content != nil {
		updated = append(updated, "内容")
	}
	if category != nil {
		updated = append(updated, "分类")
	}
	if tags != nil {
		updated = append(updated, "标签")
	}
	if priority != nil {
		updated = append(updated, "优先级")
	}

	cli.PrintSuccess(fmt.Sprintf("记忆 %s 更新成功！更新字段: %s", code, strings.Join(updated, ", ")))
	return nil
}
