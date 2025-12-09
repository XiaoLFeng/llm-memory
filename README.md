# llm-memory

<div align="center">
  <img src="https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go" alt="Go Version">
  <img src="https://img.shields.io/badge/License-MIT-green.svg" alt="License">
  <img src="https://img.shields.io/badge/Platform-macOS%20%7C%20Linux%20%7C%20Windows-lightgrey" alt="Platform">
</div>

大模型统一记忆管理工具，提供记忆（Memory）、计划（Plan）、待办（Todo）的全生命周期管理，支持 MCP 协议与 TUI 界面。

## ✨ 特性

- 🧠 **记忆管理**：支持全局/项目/小组三级作用域的记忆存储
- 📋 **计划管理**：多步骤计划创建与进度跟踪（支持子任务）
- ✅ **待办管理**：优先级、标签、截止日期的完整任务管理
- 👥 **小组协作**：基于路径的小组隔离与共享
- 🔌 **MCP 协议**：完整的 Model Context Protocol 服务端实现
- 🎨 **现代 TUI**：基于 Bubble Tea 的青绿色主题终端界面
- 🗄️ **纯 Go 实现**：无 CGO 依赖，使用 glebarez/sqlite（纯 Go SQLite 驱动）
- 🌐 **跨平台**：支持 macOS、Linux、Windows

## 📦 安装

### Homebrew（macOS & Linux）

```bash
# 添加 Tap（首次安装）
brew tap XiaoLFeng/tap

# 安装
brew install XiaoLFeng/tap/llm-memory

# 验证安装
llm-memory --version
```

### 一键安装脚本

**Unix/Linux/macOS:**

```bash
# 安装最新版本
curl -fsSL https://raw.githubusercontent.com/XiaoLFeng/llm-memory/master/scripts/install.sh | bash

# 或指定版本
curl -fsSL https://raw.githubusercontent.com/XiaoLFeng/llm-memory/master/scripts/install.sh | bash -s v0.1.0
```

**Windows (PowerShell):**

```powershell
# 安装最新版本
iwr -useb https://raw.githubusercontent.com/XiaoLFeng/llm-memory/master/scripts/install.ps1 | iex

# 或指定版本
& ([scriptblock]::Create((iwr -useb https://raw.githubusercontent.com/XiaoLFeng/llm-memory/master/scripts/install.ps1))) -Version v0.1.0
```

**卸载:**

```bash
# Unix/Linux/macOS
curl -fsSL https://raw.githubusercontent.com/XiaoLFeng/llm-memory/master/scripts/uninstall.sh | bash

# Windows (PowerShell)
iwr -useb https://raw.githubusercontent.com/XiaoLFeng/llm-memory/master/scripts/uninstall.ps1 | iex
```

### 从源码编译

```bash
# 克隆仓库
git clone https://github.com/XiaoLFeng/llm-memory.git
cd llm-memory

# 编译（需要 Go 1.23+）
go build -o llm-memory

# 移动到 PATH（可选）
sudo mv llm-memory /usr/local/bin/
```

### 预编译二进制

从 [Releases](https://github.com/XiaoLFeng/llm-memory/releases) 页面下载对应平台的二进制文件：

- macOS: `llm-memory-darwin-amd64` / `llm-memory-darwin-arm64`
- Linux: `llm-memory-linux-amd64` / `llm-memory-linux-arm64`
- Windows: `llm-memory-windows-amd64.exe`

## 🚀 快速开始

### 启动 MCP 服务

```bash
llm-memory mcp
```

### 启动 TUI 界面

```bash
llm-memory tui
```

### CLI 命令示例

```bash
# 记忆管理
llm-memory memory create --title "API 密钥" --content "sk-xxx" --global
llm-memory memory list
llm-memory memory search "API"

# 计划管理
llm-memory plan create --title "重构项目" --description "模块化架构"
llm-memory plan list
llm-memory plan update <code> --progress 50

# 待办管理
llm-memory todo create --title "修复 Bug #123" --priority 4
llm-memory todo list --scope personal
llm-memory todo complete <code>

# 小组管理
llm-memory group create --name "开发组"
llm-memory group add --group "开发组" --path /path/to/project
```

## 📖 架构设计

```
cmd/           -> Cobra CLI 命令入口
startup/       -> 启动引导器（初始化顺序管理）
internal/
  ├── service/ -> 业务逻辑层（MemoryService、PlanService 等）
  ├── models/  -> 数据层
  │   ├── entity/  -> GORM 实体（数据库表）
  │   ├── dto/     -> 数据传输对象
  │   └── *_model.go -> 数据访问对象（DAO）
  ├── mcp/     -> MCP 协议实现
  ├── tui/     -> Bubble Tea TUI（青绿色主题）
  ├── cli/     -> CLI 处理器与输出格式化
  └── database/-> SQLite + 雪花 ID
pkg/types/     -> 共享类型定义（Scope、ScopeContext）
```

### 关键设计模式

- **纯关联模式**：通过 `Global` + `PathID` 字段实现作用域隔离
- **可见性过滤器**：统一的 `VisibilityFilter` 处理权限查询
- **雪花 ID**：分布式唯一 ID 生成（非自增）
- **WAL 模式**：SQLite 写前日志模式，支持并发读写

## 🛠️ 开发指南

### 构建

```bash
go build -o llm-memory
```

### 测试

```bash
go test ./...
```

### 数据库

- 位置：`~/.llm-memory/llm-memory.db`
- 驱动：`github.com/glebarez/sqlite`（纯 Go 实现）
- 模式：WAL（Write-Ahead Logging）

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 许可

MIT License - 详见 [LICENSE](LICENSE)

## 👤 作者

筱锋 (xiao_lfeng)

## 🔗 相关链接

- [MCP 协议规范](https://modelcontextprotocol.io/)
- [Bubble Tea 框架](https://github.com/charmbracelet/bubbletea)
- [GORM 文档](https://gorm.io/)
