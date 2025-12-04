package models

import (
	"strings"

	"gorm.io/gorm"
)

// VisibilityFilter 用于统一控制全局/非全局数据的查询范围
// 嘿嘿~ 一个过滤器走天下，CLI/MCP/TUI 都能共用！(´∀｀)b
type VisibilityFilter struct {
	IncludeGlobal    bool    // 是否包含全局数据
	IncludeNonGlobal bool    // 是否包含非全局数据（私有/小组）
	PathIDs          []int64 // 限定非全局数据所属的路径ID列表，空则不过滤
}

// DefaultVisibilityFilter 返回默认过滤器：全局 + 所有非全局
func DefaultVisibilityFilter() VisibilityFilter {
	return VisibilityFilter{
		IncludeGlobal:    true,
		IncludeNonGlobal: true,
		PathIDs:          nil,
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
		} else {
			conditions = append(conditions, "global = 0")
		}
	}

	// 如果没有任何条件，返回空结果以避免全表扫描
	if len(conditions) == 0 {
		return db.Where("1 = 0")
	}

	return db.Where(strings.Join(conditions, " OR "), args...)
}

// mergePathIDs 将单个 pathID 与路径列表去重合并
func mergePathIDs(pathID int64, pathIDs []int64) []int64 {
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
