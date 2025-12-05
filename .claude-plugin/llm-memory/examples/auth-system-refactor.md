# 示例 1: 用户认证系统重构

## 场景描述

**用户需求**：重构现有的用户认证系统，采用 JWT 机制，支持 refresh token。

**项目特点**：
- 复杂度：高
- 预计时长：1 周
- 涉及模块：数据库、后端 API、中间件、测试
- 依赖关系：有明确的任务依赖

## 步骤 1：需求分析

### 任务评估
- **复杂度**：高（涉及数据库、API、前端集成）
- **时间跨度**：长期（预计 1 周）
- **依赖关系**：有（数据库设计 → API 实现 → 前端集成）
- **知识积累**：需要（设计决策、安全最佳实践）

### 决策
✅ 创建 Plan：多步骤长期项目
✅ 创建 Todos：5 个任务，按依赖关系排序
✅ 创建 Memory：记录架构设计决策

## 步骤 2：设计方案

### Plan 设计

```markdown
📋 **工作流建议方案**

---

## Plan: 用户认证系统重构

**Code:** `plan-auth-refactor`
**描述:** 采用 JWT 机制，支持 refresh token，提升安全性
**详细内容:**
```
# 用户认证系统重构实施计划

## 阶段 1: 数据库设计 (Day 1-2)
- 设计 users 表结构
- 设计 refresh_tokens 表结构
- 添加必要的索引和约束
- 编写数据库迁移脚本

## 阶段 2: JWT 核心实现 (Day 2-3)
- 实现 JWT 生成逻辑
- 实现 JWT 验证逻辑
- 实现 refresh token 机制
- 配置过期时间和密钥管理

## 阶段 3: API 端点开发 (Day 3-4)
- POST /api/auth/register - 用户注册
- POST /api/auth/login - 用户登录
- POST /api/auth/refresh - 刷新令牌
- POST /api/auth/logout - 登出

## 阶段 4: 中间件和安全 (Day 4-5)
- 实现 JWT 验证中间件
- 添加到受保护路由
- 实现登录失败限流
- CSRF 保护

## 阶段 5: 测试和验证 (Day 5-7)
- 单元测试（覆盖率 > 80%）
- 集成测试
- 安全测试
- 性能测试
```

---

## Todos: (共 5 个任务)

### 1️⃣ [Priority 3 🟠 高] 设计认证数据库架构
- **Code:** `todo-design-auth-schema`
- **描述:** 设计 users、sessions、tokens 等表结构，支持 JWT 和 refresh token
- **原因:** 基础架构，必须优先完成

### 2️⃣ [Priority 4 🔴 紧急] 实现 JWT 令牌机制
- **Code:** `todo-implement-jwt`
- **描述:** 实现 JWT 生成、验证、刷新逻辑
- **原因:** 核心功能，阻塞其他开发

### 3️⃣ [Priority 3 🟠 高] 开发登录和注册 API
- **Code:** `todo-auth-api-endpoints`
- **描述:** POST /login, POST /register, POST /refresh 等端点

### 4️⃣ [Priority 2 🟡 中] 添加认证中间件
- **Code:** `todo-auth-middleware`
- **描述:** 实现 JWT 验证中间件，保护受限路由

### 5️⃣ [Priority 2 🟡 中] 编写认证单元测试
- **Code:** `todo-auth-unit-tests`
- **描述:** 测试覆盖率达到 80% 以上

---

## Memory: 认证系统设计决策

**Code:** `mem-auth-system-design`
**标题:** 用户认证系统设计决策
**分类:** 架构设计
**标签:** 认证, JWT, 安全, 数据库
**内容:**

```markdown
# 用户认证系统设计决策

## 技术选型

### 选择 JWT 的原因
1. **无状态**：不需要服务器端 session 存储
2. **易扩展**：支持水平扩展和微服务架构
3. **跨域支持**：天然支持跨域认证
4. **标准化**：行业标准，库支持完善

### Refresh Token 机制
- **Access Token**: 短期（15分钟），存储少量信息
- **Refresh Token**: 长期（7天），用于刷新 access token
- **安全性**: refresh token 存储在 httpOnly cookie，防止 XSS

## 数据库设计

### users 表
- id (primary key)
- username (unique)
- email (unique)
- password_hash (bcrypt)
- created_at, updated_at

### refresh_tokens 表
- id (primary key)
- user_id (foreign key)
- token_hash (SHA256)
- expires_at
- created_at

## 安全考虑
1. 密码使用 bcrypt 哈希（cost=12）
2. 防止暴力破解：登录失败限流
3. HTTPS only
4. CSRF 保护（双重 cookie 提交）

## 参考资料
- [JWT Best Practices](https://jwt.io/introduction/)
- [OWASP Authentication Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Authentication_Cheat_Sheet.html)
```

---

**是否确认创建？**
```

## 步骤 3：用户确认

用户回复："确认，看起来很完善！"

## 步骤 4：执行命令

### 使用批量创建（推荐）

```bash
# 一次性创建所有 Todos
./llm-memory todo batch-create --json '[
  {"code":"todo-design-auth-schema","title":"设计认证数据库架构","description":"设计 users、sessions、tokens 等表结构，支持 JWT 和 refresh token","priority":3},
  {"code":"todo-implement-jwt","title":"实现 JWT 令牌机制","description":"实现 JWT 生成、验证、刷新逻辑","priority":4},
  {"code":"todo-auth-api-endpoints","title":"开发登录和注册 API","description":"POST /login, POST /register, POST /refresh 等端点","priority":3},
  {"code":"todo-auth-middleware","title":"添加认证中间件","description":"实现 JWT 验证中间件，保护受限路由","priority":2},
  {"code":"todo-auth-unit-tests","title":"编写认证单元测试","description":"测试覆盖率达到 80% 以上","priority":2}
]'

# 创建 Plan
./llm-memory plan create \
  --code "plan-auth-refactor" \
  --title "用户认证系统重构" \
  --description "采用 JWT 机制，支持 refresh token，提升安全性" \
  --content "# 用户认证系统重构实施计划

## 阶段 1: 数据库设计 (Day 1-2)
- 设计 users 表结构
- 设计 refresh_tokens 表结构
- 添加必要的索引和约束
- 编写数据库迁移脚本

## 阶段 2: JWT 核心实现 (Day 2-3)
- 实现 JWT 生成逻辑
- 实现 JWT 验证逻辑
- 实现 refresh token 机制
- 配置过期时间和密钥管理

## 阶段 3: API 端点开发 (Day 3-4)
- POST /api/auth/register
- POST /api/auth/login
- POST /api/auth/refresh
- POST /api/auth/logout

## 阶段 4: 中间件和安全 (Day 4-5)
- 实现 JWT 验证中间件
- 添加到受保护路由
- 实现登录失败限流
- CSRF 保护

## 阶段 5: 测试和验证 (Day 5-7)
- 单元测试
- 集成测试
- 安全测试
- 性能测试"

# 创建 Memory
./llm-memory memory create \
  --code "mem-auth-system-design" \
  --title "用户认证系统设计决策" \
  --content "<详细内容见上面的 Memory 内容部分>" \
  --category "架构设计" \
  --tags "认证,JWT,安全,数据库"
```

### 或使用单个创建（传统方式）

```bash
# 逐个创建 Todos
./llm-memory todo create --code "todo-design-auth-schema" --title "设计认证数据库架构" --description "..." --priority 3
./llm-memory todo create --code "todo-implement-jwt" --title "实现 JWT 令牌机制" --description "..." --priority 4
# ... 其他 3 个
```

## 步骤 5：执行结果

```
✅ 批量创建成功! 共处理 5 个待办事项
✅ 计划创建成功！标识码: plan-auth-refactor, 标题: 用户认证系统重构
✅ 记忆创建成功！标识码: mem-auth-system-design

工作流创建完成！开始执行第一个任务吧~ (｡･ω･｡)ﾉ゛
```

## 进度跟踪示例

### 第 1 天：完成数据库设计

```bash
# 开始任务
./llm-memory todo start --code "todo-design-auth-schema"

# ... 工作中 ...

# 完成任务
./llm-memory todo complete --code "todo-design-auth-schema"

# 更新 Plan 进度
./llm-memory plan progress --code "plan-auth-refactor" --progress 20
```

### 第 3 天：完成 JWT 实现

```bash
./llm-memory todo start --code "todo-implement-jwt"
# ... 完成后 ...
./llm-memory todo complete --code "todo-implement-jwt"
./llm-memory plan progress --code "plan-auth-refactor" --progress 50
```

### 第 7 天：项目完成

```bash
# 批量完成剩余任务
./llm-memory todo batch-complete --codes "todo-auth-middleware,todo-auth-unit-tests"

# 完成计划
./llm-memory plan complete --code "plan-auth-refactor"
```

## 关键学习点

### 1. 优先级判断
- **JWT 实现设为 Priority 4（紧急）**：因为它阻塞其他开发
- **测试设为 Priority 2（中）**：因为不阻塞主流程
- **中间件设为 Priority 2（中）**：可以在API完成后进行

### 2. Memory 使用
- 记录架构决策的"为什么"，而不仅仅是"是什么"
- 包含参考链接和安全考虑
- 使用 Markdown 格式化，增强可读性

### 3. Plan 管理
- 定期更新进度百分比
- Content 字段使用 Markdown 格式化
- 划分清晰的阶段和时间线

### 4. 批量操作的优势
- **效率提升**：一次性创建 5 个待办，而不是执行 5 次命令
- **原子性**：所有待办在一个操作中创建
- **错误处理**：部分失败不影响其他项目

## 完整命令清单

<details>
<summary>点击展开所有命令</summary>

```bash
# Plan 创建
./llm-memory plan create \
  --code "plan-auth-refactor" \
  --title "用户认证系统重构" \
  --description "采用 JWT 机制，支持 refresh token，提升安全性" \
  --content "<见上面的 Plan 内容>"

# Todos 批量创建
./llm-memory todo batch-create --json '[
  {"code":"todo-design-auth-schema","title":"设计认证数据库架构","description":"设计 users、sessions、tokens 等表结构，支持 JWT 和 refresh token","priority":3},
  {"code":"todo-implement-jwt","title":"实现 JWT 令牌机制","description":"实现 JWT 生成、验证、刷新逻辑","priority":4},
  {"code":"todo-auth-api-endpoints","title":"开发登录和注册 API","description":"POST /login, POST /register, POST /refresh 等端点","priority":3},
  {"code":"todo-auth-middleware","title":"添加认证中间件","description":"实现 JWT 验证中间件，保护受限路由","priority":2},
  {"code":"todo-auth-unit-tests","title":"编写认证单元测试","description":"测试覆盖率达到 80% 以上","priority":2}
]'

# Memory 创建
./llm-memory memory create \
  --code "mem-auth-system-design" \
  --title "用户认证系统设计决策" \
  --content "<详细的设计文档>" \
  --category "架构设计" \
  --tags "认证,JWT,安全,数据库"

# 开始第一个任务
./llm-memory todo start --code "todo-design-auth-schema"

# 完成任务并更新进度
./llm-memory todo complete --code "todo-design-auth-schema"
./llm-memory plan progress --code "plan-auth-refactor" --progress 20

# 最终完成
./llm-memory plan complete --code "plan-auth-refactor"
```

</details>

## 扩展：使用 JSON 文件

对于复杂的批量操作，推荐使用 JSON 文件：

```bash
# 创建 todos.json
cat > todos.json <<EOF
[
  {
    "code": "todo-design-auth-schema",
    "title": "设计认证数据库架构",
    "description": "设计 users、sessions、tokens 等表结构",
    "priority": 3
  },
  {
    "code": "todo-implement-jwt",
    "title": "实现 JWT 令牌机制",
    "description": "实现 JWT 生成、验证、刷新逻辑",
    "priority": 4
  }
]
EOF

# 使用文件批量创建
./llm-memory todo batch-create --json-file ./todos.json
```

---

**返回**: [示例索引](./README.md) | [下一个示例](./fix-login-bug.md)
