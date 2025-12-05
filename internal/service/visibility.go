package service

import (
	"strings"

	"github.com/XiaoLFeng/llm-memory/internal/models"
	"github.com/XiaoLFeng/llm-memory/pkg/types"
)

// resolveDefaultPathID 获取默认的非全局路径ID
func resolveDefaultPathID(scopeCtx *types.ScopeContext) int64 {
	if scopeCtx != nil && scopeCtx.PathID > 0 {
		return scopeCtx.PathID
	}
	return 0
}

// buildVisibilityFilter 将 scope 字符串解析为统一的查询过滤器（供 Memory 使用）
// 默认（空字符串）返回"全局 + 所有非全局"结果，避免遗漏数据
func buildVisibilityFilter(scope string, scopeCtx *types.ScopeContext) models.VisibilityFilter {
	filter := models.DefaultVisibilityFilter()

	switch strings.ToLower(scope) {
	case "global":
		filter.IncludeNonGlobal = false
	case "personal":
		filter.IncludeGlobal = false
		pathID := resolveDefaultPathID(scopeCtx)
		if pathID > 0 {
			filter.PathIDs = []int64{pathID}
		} else {
			filter.IncludeNonGlobal = false
		}
	case "group":
		filter.IncludeGlobal = false
		if scopeCtx != nil && len(scopeCtx.GroupPathIDs) > 0 {
			filter.PathIDs = scopeCtx.GroupPathIDs
		} else {
			filter.IncludeNonGlobal = false
		}
	case "all", "":
		// 默认：全局 + 当前路径相关（使用 scopeCtx，安全优先）
		// 嘿嘿~ 这样就不会泄露其他路径的私有数据啦！(´∀｀)b
		filter.IncludeGlobal = true
		filter.IncludeNonGlobal = true
		if scopeCtx != nil {
			filter.PathIDs = models.MergePathIDs(scopeCtx.PathID, scopeCtx.GroupPathIDs)
		}
		// 如果没有有效的 PathID，只显示全局数据
		if len(filter.PathIDs) == 0 {
			filter.IncludeNonGlobal = false
		}
	default:
		// 未知 scope：只显示全局（安全优先）
		filter.IncludeGlobal = true
		filter.IncludeNonGlobal = false
	}

	return filter
}

// buildPathOnlyFilter 将 scope 字符串解析为路径过滤器（供 Todo/Plan 使用，无 Global 支持）
// 嘿嘿~ Todo 和 Plan 不需要全局，只用路径过滤就够啦！＼(^o^)／
func buildPathOnlyFilter(scope string, scopeCtx *types.ScopeContext) models.PathOnlyVisibilityFilter {
	filter := models.DefaultPathOnlyFilter()

	switch strings.ToLower(scope) {
	case "personal":
		pathID := resolveDefaultPathID(scopeCtx)
		if pathID > 0 {
			filter.PathIDs = []int64{pathID}
		}
	case "group":
		if scopeCtx != nil && len(scopeCtx.GroupPathIDs) > 0 {
			filter.PathIDs = scopeCtx.GroupPathIDs
		}
	case "all", "":
		// 默认：当前路径 + 组路径（无全局）
		if scopeCtx != nil {
			filter.PathIDs = models.MergePathIDs(scopeCtx.PathID, scopeCtx.GroupPathIDs)
		}
	}

	return filter
}
