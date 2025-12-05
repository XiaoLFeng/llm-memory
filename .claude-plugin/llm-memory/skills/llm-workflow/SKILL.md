---
name: llm-workflow
description: |
  智能工作流管理助手 - 为 llm-memory 项目设计，用于管理计划(Plan)、待办(Todo)和知识库(Memory)。

  **何时调用此 Skill：**
  - 用户说"帮我规划一个任务"、"创建计划"、"制定方案"
  - 用户要求"拆解任务"、"创建待办"、"添加 TODO"
  - 用户需要"记录信息"、"保存知识"、"创建笔记"
  - 用户提到"跟踪进度"、"更新状态"、"标记完成"
  - 用户询问"工作流"、"项目管理"、"任务管理"

  **工作方式：**
  1. 分析用户需求，设计 Plan + Todos + Memory 的完整方案
  2. 展示建议给用户确认（包括 code 命名、优先级判断）
  3. 用户确认后，执行 llm-memory CLI 命令创建
  4. 自动跟踪进度，提醒更新状态
  5. 在关键节点创建 Memory 记录知识
---

# LLM-Memory 智能工作流管理 Skill

嘿嘿~ 欢迎使用 LLM-Memory 工作流管理助手！(´∀｀)💖

## 背景介绍

llm-memory 是一个为大模型设计的统一记忆管理系统，提供三大核心功能：

1. **Plan（计划）**：管理复杂的长期目标和多步骤任务
2. **Todo（待办）**：跟踪短期任务和具体行动项
3. **Memory（记忆）**：构建项目知识库，记录关键信息

所有内容默认存储在**项目级别**（与当前工作目录关联），使用 `--global` 标志才会存储为全局可见。

## ⭐ Code 格式规则（必读！）

这是最容易出错的地方，必须严格遵守！

### 格式要求

```
✅ 规则：
- 全小写字母（a-z）
- 可包含连字符（-）
- 开头和末尾必须是字母
- 最少 3 个字符

✅ 有效示例：
plan-api-redesign
todo-fix-auth-bug
mem-deployment-notes
user-auth-system
api-v2

❌ 无效示例：
test_plan_001    ❌ 含下划线
Task-001         ❌ 含大写字母
-my-task         ❌ 开头不是字母
task-            ❌ 末尾不是字母
ab               ❌ 少于 3 个字符
myTask           ❌ 驼峰命名
```

### 推荐命名模式

```bash
# Plan 命名
plan-<简短描述>
例：plan-user-system, plan-database-migration

# Todo 命名
todo-<动作>-<对象>
例：todo-fix-login-bug, todo-add-api-docs

# Memory 命名
mem-<主题>
例：mem-api-design, mem-security-notes
```

## 📊 优先级判断规则

自动为 Todo 判断优先级（1-4）：

### Priority 4 - 紧急 🔴

```
触发条件：
✅ Bug 修复、安全漏洞、系统故障
✅ 阻塞其他工作的任务
✅ 生产环境问题
✅ 截止时间在 24 小时内

示例：
- 修复用户无法登录的 bug
- 处理数据库连接失败
- 紧急安全补丁
```

### Priority 3 - 高 🟠

```
触发条件：
✅ 重要功能开发
✅ 影响用户体验的问题
✅ 关键里程碑任务
✅ 截止时间在 3 天内

示例：
- 实现核心 API 端点
- 重要页面的 UI 优化
- 集成第三方服务
```

### Priority 2 - 中 🟡（默认）

```
触发条件：
✅ 常规开发任务
✅ 功能优化改进
✅ 无明确截止时间
✅ 不确定的优先级

示例：
- 添加单元测试
- 代码重构
- 文档更新
```

### Priority 1 - 低 🟢

```
触发条件：
✅ 可选的改进
✅ 技术债务清理
✅ 长期优化计划
✅ 学习和探索任务

示例：
- 性能优化探索
- 代码注释补充
- 依赖包升级
```

## 🔄 混合模式交互流程

这是此 Skill 的核心工作方式！必须严格遵循：

### 步骤 1：分析需求

当用户提出需求时，首先进行分析：

```
需要考虑：
1. 任务复杂度：是单一任务还是多步骤项目？
2. 时间跨度：短期（<1天）还是长期（>3天）？
3. 依赖关系：任务之间有无依赖？
4. 知识积累：是否需要记录关键信息？

决策逻辑：
- 单一简单任务 → 只创建 Todo
- 多步骤任务 → 创建 Plan + Todos
- 需要记录知识 → 添加 Memory
```

### 步骤 2：提出建议方案

向用户展示格式化的建议，等待确认：

```markdown
📋 **工作流建议方案**

---

## Plan: <计划标题>

**Code:** `plan-xxx-xxx`
**描述:** <简要说明计划目标>
**详细内容:**
```
<Markdown 格式的实施步骤>
- 步骤 1：...
- 步骤 2：...
- 步骤 3：...
```

---

## Todos: (共 N 个任务)

### 1️⃣ [Priority 4 🔴 紧急] <任务标题>
- **Code:** `todo-xxx-xxx`
- **描述:** <任务详情>
- **原因:** <为什么是紧急>

### 2️⃣ [Priority 3 🟠 高] <任务标题>
- **Code:** `todo-yyy-yyy`
- **描述:** <任务详情>

### 3️⃣ [Priority 2 🟡 中] <任务标题>
- **Code:** `todo-zzz-zzz`
- **描述:** <任务详情>

---

## Memory: (可选)

**Code:** `mem-xxx-xxx`
**标题:** <知识点标题>
**分类:** <分类名称>
**标签:** tag1, tag2, tag3
**内容:**
```
<要记录的知识内容>
```

---

**是否确认创建？**
如需调整（修改优先级、code、增删任务），请告诉我~ (｡･ω･｡)ﾉ゛
```

### 步骤 3：等待用户确认

用户可以：
- 直接确认 → 执行步骤 4
- 修改建议 → 调整后重新展示
- 取消操作 → 结束流程

### 步骤 4：执行 CLI 命令

⚠️ **重要提示**：当前 `plan create` 命令存在 BUG，无法使用！请跳过 Plan 创建。

```bash
# ❌ Plan 创建（当前不可用）
# ./main plan create --code "plan-xxx" --title "标题" --description "描述"
# 原因：CLI 缺少 --content 参数，导致创建失败

# ✅ 创建 Todos（按优先级排序）
./main todo create \
  --code "todo-task-one" \
  --title "第一个任务" \
  --description "任务详细说明" \
  --priority 4

./main todo create \
  --code "todo-task-two" \
  --title "第二个任务" \
  --priority 3

# ✅ 创建 Memory（如有需要）
./main memory create \
  --code "mem-knowledge-point" \
  --title "关键知识点" \
  --content "详细的知识内容..." \
  --category "架构设计" \
  --tags "tag1,tag2,tag3"
```

### 步骤 5：确认结果

检查每个命令的输出：

```
✅ 成功输出示例：
"待办创建成功！标识码: todo-xxx, 标题: xxx"
"记忆创建成功！标识码: mem-xxx"

❌ 错误输出示例：
"code 格式错误: 全小写字母，可含连字符，开头末尾必须是字母"
"活跃状态中已存在相同的 code"
```

如有错误，立即向用户报告并提供解决方案。

## 📝 CLI 命令完整清单

### Todo 命令（当前主要使用）

```bash
# 创建待办
./main todo create \
  --code <code> \
  --title <title> \
  [--description <desc>] \
  [--priority 1-4] \
  [--global]

# 列出所有待办
./main todo list

# 开始待办（状态 → in_progress）
./main todo start --code <code>

# 完成待办（状态 → completed）
./main todo complete --code <code>

# 删除待办
./main todo delete --code <code>

# 标记所有待办为完成
./main todo final
```

### Todo 批量命令（高效处理）

```bash
# 批量创建待办（JSON 格式）
./main todo batch-create --json '[
  {"code":"t1","title":"任务1","priority":3},
  {"code":"t2","title":"任务2","description":"详情"}
]'

# 或使用 JSON 文件
./main todo batch-create --json-file ./todos.json

# 批量开始待办
./main todo batch-start --codes "t1,t2,t3"

# 批量完成待办
./main todo batch-complete --codes "t1,t2"

# 批量取消待办
./main todo batch-cancel --codes "t1,t2"

# 批量删除待办
./main todo batch-delete --codes "t1,t2"

# 批量更新待办
./main todo batch-update --json '[
  {"code":"t1","title":"新标题","priority":4},
  {"code":"t2","status":2}
]'

# 批量操作特点：
# ✅ 支持最多 100 个项目
# ✅ 混合模式输出（全成功/全失败/部分成功）
# ✅ 详细的错误信息
```

### Memory 命令

```bash
# 创建记忆
./main memory create \
  --code <code> \
  --title <title> \
  --content <content> \
  [--category <category>] \
  [--tags <tag1,tag2,tag3>] \
  [--global]

# 列出所有记忆
./main memory list

# 搜索记忆
./main memory search --keyword <keyword>

# 获取记忆详情
./main memory get --code <code>

# 删除记忆
./main memory delete --code <code>
```

### Plan 命令

```bash
# 创建计划（支持 Markdown 格式的详细内容）
./main plan create \
  --code <code> \
  --title <title> \
  --description <desc> \
  --content <markdown-content> \
  [--global]

# 列出所有计划
./main plan list

# 开始计划
./main plan start --code <code>

# 更新进度（0-100）
./main plan progress --code <code> --progress <value>

# 完成计划
./main plan complete --code <code>

# 删除计划
./main plan delete --code <code>
```

## 💡 Memory 使用指南

### 何时创建 Memory

**必须创建的场景：**

```
1. 架构决策
   - 为什么选择某个技术方案
   - 设计模式的选择理由
   - API 设计约定

2. 问题解决方案
   - 复杂 Bug 的排查过程
   - 性能优化的关键点
   - 踩坑经验和解决办法

3. 代码示例
   - 常用代码片段
   - API 使用示例
   - 配置模板

4. 项目规范
   - 代码风格约定
   - 命名规范
   - Git 工作流
```

**可选创建的场景：**

```
1. 调试发现
   - 有价值的调试技巧
   - 工具使用心得

2. 第三方库
   - 库的使用技巧
   - 常见问题解决

3. 环境配置
   - 开发环境设置
   - 部署流程说明
```

### Memory 内容格式建议

使用 Markdown 格式编写，增强可读性：

```markdown
# <知识点标题>

## 背景
<为什么需要这个知识点>

## 核心内容
<详细说明>

## 代码示例
```<language>
<code here>
```

## 参考链接
- [Link 1](URL)
- [Link 2](URL)

## 注意事项
- 注意点 1
- 注意点 2
```

## 🎯 完整示例场景

我们提供了两个完整的示例场景，展示如何在真实项目中使用工作流管理：

### 示例索引

#### [示例 1: 用户认证系统重构](./examples/auth-system-refactor.md)
- **难度**: 高
- **特点**: 复杂的长期项目，展示 Plan + Todos + Memory 完整流程
- **学习要点**: 任务拆解、优先级判断、架构决策记录

#### [示例 2: 修复登录 Bug](./examples/fix-login-bug.md)
- **难度**: 低
- **特点**: 简单的单任务示例，快速处理
- **学习要点**: 紧急任务处理、简单流程

查看 [examples/README.md](./examples/README.md) 了解更多示例和使用指南。

---

## 🔧 进度跟踪

完成任务后，及时更新状态：

```bash
# 开始任务
./main todo start --code "todo-design-auth-schema"
# 输出：待办 todo-design-auth-schema 已开始

# 完成任务
./main todo complete --code "todo-design-auth-schema"
# 输出：待办 todo-design-auth-schema 已完成

# 查看所有待办
./main todo list
# 输出：待办列表（表格形式）
```

### 进度跟踪最佳实践

```
1. 开始任务前：
   ✅ 执行 start 命令标记为进行中
   ✅ 确认没有依赖任务未完成

2. 任务完成后：
   ✅ 立即执行 complete 命令
   ✅ 如果有新发现，创建 Memory 记录
   ✅ 如果有新任务，创建新 Todo

3. 定期检查：
   ✅ 每天执行 list 查看未完成任务
   ✅ 重新评估优先级
   ✅ 识别阻塞任务
```

## ❌ 错误处理

### 常见错误 1：Code 格式错误

```
错误信息：
code 格式错误: 全小写字母，可含连字符，开头末尾必须是字母

原因：
使用了大写字母、下划线、数字开头等不符合规则的字符

解决方案：
1. 检查 code 是否全小写
2. 将下划线（_）替换为连字符（-）
3. 确保开头和末尾是字母
4. 确保长度 >= 3

示例修正：
❌ test_plan_001 → ✅ test-plan-one
❌ Task-001      → ✅ task-one
❌ -my-task      → ✅ my-task
```

### 常见错误 2：Code 重复

```
错误信息：
活跃状态中已存在相同的 code，请使用不同的标识码

原因：
当前已有相同 code 的活跃（非完成/取消）记录

解决方案：
1. 执行 list 命令查看已存在的 code
2. 使用更具体的命名
3. 或者先完成/删除旧记录

示例修正：
❌ todo-fix-bug → ✅ todo-fix-login-bug
❌ mem-notes    → ✅ mem-api-design-notes
```

### 常见错误 3：Plan 创建失败

```
错误信息：
计划内容不能为空

原因：
CLI 存在 BUG，无法传递 content 参数

解决方案：
当前版本暂不支持 Plan 创建，请：
1. 只使用 Todo 管理任务
2. 或等待 CLI 修复后使用

预计修复：
需要在以下文件中添加 --content 参数：
- cmd/plan/plan_create.go
- internal/cli/handlers/plan.go
```

### 常见错误 4：找不到记录

```
错误信息：
计划不存在或已完成/取消
待办不存在
记忆不存在

原因：
1. Code 拼写错误
2. 作用域不匹配（全局 vs 项目）
3. 记录已删除

解决方案：
1. 执行 list 命令确认 code
2. 检查当前工作目录
3. 如需全局记录，使用 --global 创建
```

## ✅ 最佳实践

### 1. Code 命名规范

```
✅ 推荐做法：
- 使用语义化名称：plan-api-redesign
- 简洁但清晰：todo-fix-login
- 包含动作词：todo-add-tests, todo-update-docs
- 体现层级关系：todo-auth-jwt-impl, todo-auth-middleware

❌ 避免做法：
- 过长：plan-this-is-a-very-long-description-of-the-task
- 过短：abc, tsk
- 数字编号：task001, todo-123
- 无意义：temp, test, xxx
```

### 2. 优先级设置

```
✅ 推荐做法：
- 基于影响范围和紧急程度
- 阻塞性任务设为高优先级
- 不确定时使用默认（中）
- 定期重新评估

❌ 避免做法：
- 所有任务都是紧急（失去意义）
- 从不更新优先级
- 忽略依赖关系
```

### 3. Plan vs Todo 选择

```
使用 Plan（当 CLI 修复后）：
✅ 多步骤复杂任务（>3步）
✅ 长期目标（>3天）
✅ 需要跟踪整体进度
✅ 有子任务依赖关系

使用 Todo：
✅ 单一任务
✅ 短期行动项（<1天）
✅ 独立任务（无依赖）
✅ 简单的提醒事项

当前建议：
⚠️ 由于 Plan 创建不可用，暂时只使用 Todo
```

### 4. Memory 内容组织

```
✅ 推荐做法：
- 使用 Markdown 格式
- 包含代码示例
- 添加参考链接
- 使用合适的分类和标签
- 内容结构清晰

❌ 避免做法：
- 纯文本无格式
- 重复记录相似内容
- 缺少上下文信息
- 标签过多或过少
```

### 5. 工作流程

```
完整流程：
1. 需求分析 → 设计方案
2. 用户确认 → 创建记录
3. 执行任务 → 更新状态
4. 完成任务 → 记录知识
5. 回顾总结 → 改进流程

日常习惯：
- 每天开始：检查 todo list
- 任务开始：执行 start 命令
- 任务完成：执行 complete 命令
- 新发现：创建 memory 记录
- 每周回顾：评估进度和优先级
```

## 🚨 当前已知限制

### ~~限制 1：Plan 创建不可用~~  ✅ 已修复

~~**问题**：CLI 的 `plan create` 命令无法成功执行~~

**状态**：已修复！CLI 已支持 `--content` 参数，Plan 功能完全可用。

**使用方式**：
```bash
./main plan create \
  --code "my-plan" \
  --title "计划标题" \
  --description "计划描述" \
  --content "# 详细内容\n\n- 步骤1\n- 步骤2"
```

### ~~限制 2：CLI 批量操作不支持~~ ✅ 已支持

~~**问题**：只能逐个创建 Todo~~

**状态**：已支持！新增 6 个批量命令（batch-create/complete/start/cancel/delete/update）。

**详见**："Todo 批量命令" 部分

### 限制 3：默认项目级别

**行为**：所有命令默认作用于当前项目路径

**影响**：
- 切换目录后看不到之前的记录
- 需要使用 `--global` 创建全局记录

**建议**：
- 在项目根目录执行命令
- 重要的通用知识使用 `--global`

### 限制 4：SubTask 管理

**问题**：CLI 尚未支持 Plan 的子任务管理

**影响**：无法通过 CLI 添加/更新子任务

**临时方案**：使用 TUI 或直接编辑数据库

### 限制 5：Tags 搜索

**问题**：Memory 搜索不支持按标签过滤

**影响**：只能通过关键词搜索

**临时方案**：在 content 中包含关键词

## 📚 快速参考

### 常用命令速查

```bash
# 创建待办（最常用）
./main todo create --code <code> --title <title> --priority <1-4>

# 查看所有待办
./main todo list

# 开始任务
./main todo start --code <code>

# 完成任务
./main todo complete --code <code>

# 创建记忆
./main memory create --code <code> --title <title> --content <content>

# 搜索记忆
./main memory search --keyword <keyword>
```

### Code 格式速查

```
格式：全小写 + 连字符
规则：字母开头和结尾，>=3 字符
示例：my-task, plan-api-redesign
```

### 优先级速查

```
1 = 低   🟢 可选改进
2 = 中   🟡 常规任务（默认）
3 = 高   🟠 重要功能
4 = 紧急 🔴 Bug/阻塞
```