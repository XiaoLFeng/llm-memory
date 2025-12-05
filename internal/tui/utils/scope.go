package utils

import "github.com/XiaoLFeng/llm-memory/startup"

// ScopeTag 根据 pathID 返回作用域标签（供 Todo/Plan 使用，无 Global 支持）
// 用于在 TUI 列表中显示数据的归属范围
func ScopeTag(pathID int64, bs *startup.Bootstrap) string {
	ctx := bs.CurrentScope
	if ctx != nil {
		for _, gid := range ctx.GroupPathIDs {
			if pathID == gid {
				return "[小组]"
			}
		}
		if ctx.PathID == pathID {
			return "[项目]"
		}
	}
	if pathID > 0 {
		return "[项目]"
	}
	return "[未知]"
}

// ScopeTagWithGlobal 根据 global 标记和 pathID 返回作用域标签（供 Memory 使用）
// 用于在 TUI 列表中显示数据的归属范围
func ScopeTagWithGlobal(global bool, pathID int64, bs *startup.Bootstrap) string {
	if global {
		return "[全局]"
	}
	ctx := bs.CurrentScope
	if ctx != nil {
		for _, gid := range ctx.GroupPathIDs {
			if pathID == gid {
				return "[小组]"
			}
		}
		if ctx.PathID == pathID {
			return "[项目]"
		}
	}
	if pathID > 0 {
		return "[项目]"
	}
	return "[未知]"
}

// ScopeFilter 作用域过滤状态
// 用于在 TUI 列表页面中切换显示不同作用域的数据
type ScopeFilter int

const (
	ScopeAll      ScopeFilter = iota // 全部（默认）
	ScopePersonal                    // 仅项目
	ScopeGroup                       // 仅小组
)

// String 返回用于 Service 层查询的作用域字符串
func (s ScopeFilter) String() string {
	switch s {
	case ScopePersonal:
		return "personal"
	case ScopeGroup:
		return "group"
	default:
		return "all"
	}
}

// Label 返回用于 UI 显示的作用域标签
func (s ScopeFilter) Label() string {
	switch s {
	case ScopePersonal:
		return "项目"
	case ScopeGroup:
		return "小组"
	default:
		return "全部"
	}
}

// Next 循环切换到下一个作用域
func (s ScopeFilter) Next() ScopeFilter {
	return (s + 1) % 3
}
