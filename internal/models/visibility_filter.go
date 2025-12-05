package models

import (
	"strings"

	"gorm.io/gorm"
)

// VisibilityFilter 用于统一控制全局/非全局数据的查询范围（供 Memory 使用）
// 嘿嘿~ 一个过滤器走天下，CLI/MCP/TUI 都能共用！(´∀｀)b
type VisibilityFilter struct {
	IncludeGlobal    bool    // 是否包含全局数据
	IncludeNonGlobal bool    // 是否包含非全局数据（私有/小组）
	PathIDs          []int64 // 限定非全局数据所属的路径ID列表，空则不过滤
}

// PathOnlyVisibilityFilter 用于 Todo/Plan 的查询范围（无 Global 支持）
// 呀~ Todo 和 Plan 不需要全局作用域，只用路径就够啦！＼(^o^)／
type PathOnlyVisibilityFilter struct {
	PathIDs []int64 // 限定数据所属的路径ID列表
}

// DefaultVisibilityFilter 返回默认过滤器：仅全局数据（安全默认值）
// 呀~ 默认只显示全局数据，私有数据需要明确指定作用域才能看到！(´∀｀)b
func DefaultVisibilityFilter() VisibilityFilter {
	return VisibilityFilter{
		IncludeGlobal:    true,
		IncludeNonGlobal: false, // 默认不显示非全局数据（安全优先）
		PathIDs:          nil,
	}
}

// DefaultPathOnlyFilter 返回默认的路径过滤器（空路径列表，安全默认值）
func DefaultPathOnlyFilter() PathOnlyVisibilityFilter {
	return PathOnlyVisibilityFilter{
		PathIDs: nil,
	}
}

// applyVisibilityFilter 将过滤条件应用到查询上
// 统一的查询拼装逻辑，避免各模型重复造轮子~
func applyVisibilityFilter(db *gorm.DB, filter VisibilityFilter) *gorm.DB {
	conditions := make([]string, 0, 2)
	args := make([]interface{}, 0, 1)

	if filter.IncludeGlobal {
		conditions = append(conditions, "global = 1")
	}

	if filter.IncludeNonGlobal {
		if len(filter.PathIDs) > 0 {
			conditions = append(conditions, "(global = 0 AND path_id IN ?)")
			args = append(args, filter.PathIDs)
		}
		// PathIDs 为空时，不添加非全局条件（只返回全局数据，安全优先）
	}

	// 如果没有任何条件，返回空结果以避免全表扫描
	if len(conditions) == 0 {
		return db.Where("1 = 0")
	}

	return db.Where(strings.Join(conditions, " OR "), args...)
}

// MergePathIDs 将单个 pathID 与路径列表去重合并（公开函数）
func MergePathIDs(pathID int64, pathIDs []int64) []int64 {
	pathSet := make(map[int64]struct{})
	if pathID > 0 {
		pathSet[pathID] = struct{}{}
	}
	for _, id := range pathIDs {
		if id > 0 {
			pathSet[id] = struct{}{}
		}
	}
	if len(pathSet) == 0 {
		return nil
	}
	result := make([]int64, 0, len(pathSet))
	for id := range pathSet {
		result = append(result, id)
	}
	return result
}

// ApplyPathOnlyFilter 将路径过滤条件应用到查询上（供 Todo/Plan 使用）
// 嘿嘿~ 没有 Global 的简化版过滤器！(´∀｀)b
func ApplyPathOnlyFilter(db *gorm.DB, filter PathOnlyVisibilityFilter) *gorm.DB {
	if len(filter.PathIDs) == 0 {
		// 没有有效路径时返回空结果（安全优先）
		return db.Where("1 = 0")
	}
	return db.Where("path_id IN ?", filter.PathIDs)
}
