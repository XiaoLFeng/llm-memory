package types

// Scope 作用域类型
// 嘿嘿~ 用于区分数据的可见范围呢！🎯
type Scope string

// 作用域常量定义
const (
	ScopePersonal Scope = "personal" // 当前目录专属
	ScopeGroup    Scope = "group"    // 组作用域
	ScopeGlobal   Scope = "global"   // 全局作用域
)

// String 将 Scope 转换为字符串
func (s Scope) String() string {
	return string(s)
}

// IsValid 检查作用域是否有效
func (s Scope) IsValid() bool {
	switch s {
	case ScopePersonal, ScopeGroup, ScopeGlobal:
		return true
	default:
		return false
	}
}

// GlobalGroupID 全局作用域的特殊 GroupID
// 值为 0，表示不属于任何特定组，全局可见
const GlobalGroupID = 0

// ScopeContext 作用域上下文
// 用于在请求链路中传递当前作用域信息
type ScopeContext struct {
	CurrentPath     string // 当前工作目录
	GroupID         int64  // 所属组 ID（0 表示无组）
	GroupName       string // 组名称（方便显示）
	IncludePersonal bool   // 查询时是否包含 Personal 数据
	IncludeGroup    bool   // 查询时是否包含 Group 数据
	IncludeGlobal   bool   // 查询时是否包含 Global 数据
}

// NewScopeContext 创建默认的作用域上下文
// 嘿嘿~ 默认显示所有作用域的数据呢！✨
func NewScopeContext(currentPath string) *ScopeContext {
	return &ScopeContext{
		CurrentPath:     currentPath,
		GroupID:         GlobalGroupID,
		GroupName:       "",
		IncludePersonal: true,
		IncludeGroup:    true,
		IncludeGlobal:   true,
	}
}

// NewGlobalOnlyScope 创建只包含全局数据的作用域
func NewGlobalOnlyScope() *ScopeContext {
	return &ScopeContext{
		CurrentPath:     "",
		GroupID:         GlobalGroupID,
		GroupName:       "",
		IncludePersonal: false,
		IncludeGroup:    false,
		IncludeGlobal:   true,
	}
}

// NewPersonalOnlyScope 创建只包含当前目录数据的作用域
func NewPersonalOnlyScope(currentPath string) *ScopeContext {
	return &ScopeContext{
		CurrentPath:     currentPath,
		GroupID:         GlobalGroupID,
		GroupName:       "",
		IncludePersonal: true,
		IncludeGroup:    false,
		IncludeGlobal:   false,
	}
}

// NewGroupOnlyScope 创建只包含组数据的作用域
func NewGroupOnlyScope(groupID int64, groupName string) *ScopeContext {
	return &ScopeContext{
		CurrentPath:     "",
		GroupID:         groupID,
		GroupName:       groupName,
		IncludePersonal: false,
		IncludeGroup:    true,
		IncludeGlobal:   false,
	}
}

// WithGroup 设置组信息
func (sc *ScopeContext) WithGroup(groupID int64, groupName string) *ScopeContext {
	sc.GroupID = groupID
	sc.GroupName = groupName
	return sc
}

// SetPersonalOnly 设置只显示 Personal 数据
func (sc *ScopeContext) SetPersonalOnly() *ScopeContext {
	sc.IncludePersonal = true
	sc.IncludeGroup = false
	sc.IncludeGlobal = false
	return sc
}

// SetGroupOnly 设置只显示 Group 数据
func (sc *ScopeContext) SetGroupOnly() *ScopeContext {
	sc.IncludePersonal = false
	sc.IncludeGroup = true
	sc.IncludeGlobal = false
	return sc
}

// SetGlobalOnly 设置只显示 Global 数据
func (sc *ScopeContext) SetGlobalOnly() *ScopeContext {
	sc.IncludePersonal = false
	sc.IncludeGroup = false
	sc.IncludeGlobal = true
	return sc
}

// HasGroup 检查是否有关联的组
func (sc *ScopeContext) HasGroup() bool {
	return sc.GroupID != GlobalGroupID
}

// GetScope 根据数据的 GroupID 和 Path 判断其作用域
func GetScope(groupID int64, path, currentPath string) Scope {
	if groupID == GlobalGroupID && path == "" {
		return ScopeGlobal
	}
	if path != "" && path == currentPath {
		return ScopePersonal
	}
	if groupID != GlobalGroupID {
		return ScopeGroup
	}
	return ScopeGlobal
}
