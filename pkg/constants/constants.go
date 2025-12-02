package constants

import (
	"os"
	"path/filepath"
	"runtime"
)

// ========================================
// 🎮 应用基础信息常量
// ========================================

// AppName 应用名称
const AppName = "llm-memory"

// AppVersion 应用版本
const AppVersion = "0.1.0"

// DefaultDBPath 默认数据库路径
var DefaultDBPath = getDefaultDBPath()

// getDefaultDBPath 获取默认数据库路径 (嘿~ 这个函数很智能哦！)
func getDefaultDBPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// 如果获取不到用户目录，使用当前目录
		return "./data.db"
	}

	// 根据操作系统构建数据库路径
	var dbDir string
	switch runtime.GOOS {
	case "windows":
		dbDir = filepath.Join(homeDir, "AppData", "Local", AppName)
	default:
		// Unix-like systems (Linux, macOS, etc.)
		dbDir = filepath.Join(homeDir, "."+AppName)
	}

	return filepath.Join(dbDir, "data.db")
}

// ========================================
// 📋 菜单相关常量
// ========================================

// MainMenuOptionCount 主菜单选项数量
const MainMenuOptionCount = 6

// MainMenuItems 主菜单项目
var MainMenuItems = []string{
	"💾 保存记忆",
	"🔍 搜索记忆",
	"📚 浏览记忆",
	"🗑️ 删除记忆",
	"📊 统计信息",
	"👋 退出程序",
}

// ModuleNames 各功能模块名称
const (
	ModuleSave   = "save"   // 保存模块
	ModuleSearch = "search" // 搜索模块
	ModuleBrowse = "browse" // 浏览模块
	ModuleDelete = "delete" // 删除模块
	ModuleStats  = "stats"  // 统计模块
	ModuleExit   = "exit"   // 退出模块
)

// ========================================
// 📊 数据相关常量
// ========================================

// DefaultPageSize 默认分页大小
const DefaultPageSize = 10

// MaxPageSize 最大分页大小
const MaxPageSize = 50

// MinPageSize 最小分页大小
const MinPageSize = 5

// MaxTitleLength 最大标题长度
const MaxTitleLength = 100

// MaxContentLength 最大内容长度
const MaxContentLength = 10000

// MaxTagLength 标签最大长度
const MaxTagLength = 50

// MaxTagCount 最大标签数量
const MaxTagCount = 10

// ========================================
// 🎨 UI 相关常量
// ========================================

// DefaultWidth 默认界面宽度
const DefaultWidth = 80

// DefaultHeight 默认界面高度
const DefaultHeight = 24

// MinWidth 最小界面宽度
const MinWidth = 60

// MinHeight 最小界面高度
const MinHeight = 20

// ========================================
// ⏱️ 时间相关常量
// ========================================

// DefaultDateFormat 默认日期格式
const DefaultDateFormat = "2006-01-02 15:04:05"

// ShortDateFormat 短日期格式
const ShortDateFormat = "2006-01-02"

// TimeFormatOnly 仅时间格式
const TimeFormatOnly = "15:04:05"

// ========================================
// 🔧 配置相关常量
// ========================================

// ConfigFileName 配置文件名
const ConfigFileName = "config.json"

// LogFileName 日志文件名
const LogFileName = "app.log"

// BackupDirName 备份目录名
const BackupDirName = "backups"

// ========================================
// 🚀 性能相关常量
// ========================================

// DefaultTimeout 默认超时时间 (秒)
const DefaultTimeout = 30

// MaxSearchResults 最大搜索结果数量
const MaxSearchResults = 100

// DatabaseConnectionRetries 数据库连接重试次数
const DatabaseConnectionRetries = 3

// ========================================
// 🎯 错误码常量
// ========================================

// ErrSuccess 操作成功
const ErrSuccess = 0

// ErrGeneral 一般错误
const ErrGeneral = 1

// ErrDBError 数据库错误
const ErrDBError = 2

// ErrInvalidInput 输入无效
const ErrInvalidInput = 3

// ErrNotFound 未找到
const ErrNotFound = 4

// ErrPermission 权限错误
const ErrPermission = 5

// ========================================
// 🌟 其他常量
// ========================================

// Author 作者信息 (嘿嘿~ 就是我啦！)
const Author = "XiaoLFeng"

// Description 应用描述
const Description = "一个轻量级的本地记忆管理工具，使用 BubbleTea 构建 TUI 界面"

// GitHubRepo GitHub 仓库地址
const GitHubRepo = "https://github.com/XiaoLFeng/llm-memory"
