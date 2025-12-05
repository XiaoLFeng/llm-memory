# 示例 2: 修复登录 Bug

## 场景描述

**用户需求**：用户报告无法登录，需要紧急修复。

**项目特点**：
- 复杂度：低
- 预计时长：2 小时
- 涉及模块：单一（认证逻辑）
- 依赖关系：无

## 步骤 1：需求分析

### 任务评估
- **复杂度**：低（单一问题排查）
- **时间跨度**：短期（预计 2 小时）
- **依赖关系**：无
- **知识积累**：不需要（简单 Bug 修复）

### 决策
❌ 不创建 Plan：单一任务
✅ 创建 Todo：1 个紧急任务
❌ 不创建 Memory：简单修复，无需记录

## 步骤 2：设计方案

### Todo 设计

```markdown
📋 **工作流建议方案**

---

## Todos: (共 1 个任务)

### 1️⃣ [Priority 4 🔴 紧急] 修复登录失败 Bug
- **Code:** `todo-fix-login-bug`
- **描述:** 排查并修复用户无法登录的问题，涉及认证逻辑
- **原因:** 影响用户使用，需要紧急处理

---

**是否确认创建？**
```

## 步骤 3：用户确认

用户："确认，尽快修复！"

## 步骤 4：执行命令

```bash
./llm-memory todo create \
  --code "todo-fix-login-bug" \
  --title "修复登录失败 Bug" \
  --description "排查并修复用户无法登录的问题，涉及认证逻辑" \
  --priority 4
```

## 步骤 5：执行结果

```
✅ 待办创建成功！标识码: todo-fix-login-bug, 标题: 修复登录失败 Bug

任务已创建！开始排查问题吧~ (｡･ω･｡)ﾉ゛
```

## 任务执行流程

### 1. 开始任务

```bash
./llm-memory todo start --code "todo-fix-login-bug"
# 输出：待办 todo-fix-login-bug 已开始
```

### 2. 排查和修复

```bash
# 查看日志，发现问题...
# 修改代码，测试修复...
```

### 3. 完成任务

```bash
./llm-memory todo complete --code "todo-fix-login-bug"
# 输出：待办 todo-fix-login-bug 已完成
```

### 4. 验证结果

```bash
# 查看待办列表
./llm-memory todo list

# 应该看到 todo-fix-login-bug 状态为"已完成"
```

## 关键学习点

### 1. 优先级判断
- **Bug 修复通常是 Priority 4（紧急）**
- 特别是影响核心功能的 Bug
- 用户无法使用系统 = 最高优先级

### 2. 简单任务处理
- 不需要 Plan（避免过度设计）
- 不需要 Memory（除非发现有价值的知识点）
- 快速创建、快速完成

### 3. 何时升级为 Plan
如果修复过程中发现需要重构多个模块，可以考虑升级：

```bash
# 1. 先完成当前的紧急修复
./llm-memory todo complete --code "todo-fix-login-bug"

# 2. 创建 Memory 记录发现的问题
./llm-memory memory create \
  --code "mem-auth-issue-found" \
  --title "认证系统架构问题" \
  --content "在修复登录 Bug 时发现..." \
  --category "技术债务"

# 3. 创建新的 Plan 进行系统性重构
#（参考示例 1）
```

## 扩展场景：发现需要深入重构

如果在修复过程中发现认证系统存在架构问题：

### 场景 1: 发现安全漏洞

```bash
# 完成紧急修复
./llm-memory todo complete --code "todo-fix-login-bug"

# 创建新的紧急 Todo
./llm-memory todo create \
  --code "todo-fix-security-vuln" \
  --title "修复认证安全漏洞" \
  --description "发现 SQL 注入风险，需要立即修复" \
  --priority 4

# 记录安全问题
./llm-memory memory create \
  --code "mem-security-finding" \
  --title "认证系统安全漏洞记录" \
  --content "## 漏洞描述\n\n## 修复方案\n\n## 预防措施" \
  --category "安全" \
  --tags "安全,漏洞,认证"
```

### 场景 2: 发现需要重构

```bash
# 完成紧急修复
./llm-memory todo complete --code "todo-fix-login-bug"

# 创建重构 Plan
./llm-memory plan create \
  --code "plan-auth-refactor" \
  --title "认证系统重构" \
  --description "解决架构问题，提升可维护性" \
  --content "# 重构计划\n\n..."

# 创建重构任务（参考示例 1）
./llm-memory todo batch-create --json '[...]'
```

## 总结

这个示例展示了：

✅ **简单任务的快速处理流程**
✅ **紧急优先级的正确判断**
✅ **何时需要升级为复杂工作流**
✅ **如何在问题排查中积累知识**

记住：不是所有任务都需要复杂的工作流！简单的 Bug 修复，一个 Todo 就够了~ (´∀｀)💖

---

**返回**: [示例索引](./README.md) | [上一个示例](./auth-system-refactor.md)
