---
name: git-conflict-resolution-strategy
description: 不同类型冲突的最佳解决策略
metadata:
  type: reference
---

根据冲突类型选择正确的解决策略。

## 策略矩阵

| 冲突类型 | 推荐策略 | 工具 | 预估时间 | 风险 |
|---|---|---|---|---|
| 简单文本冲突（1-2处） | Git 自动合并 | `git merge` | 快（5-10分钟） | 低 |
| 多区域独立冲突 | 手动 Edit | `Edit` 工具 | 中（30-60分钟） | 低 |
| 结构性重构冲突 | 逐文件对比 | `diff` + `Edit` | 慢（1-3小时） | 中 |
| 依赖关系变更 | 理解新架构 | `Read` + 分析 | 慢（2-4小时） | 高 |
| 批量文件冲突 | 逐个处理 | `Edit` 逐个 | 很慢（3-8小时） | 中 |

---

## 详细策略

### 策略 1: Git 自动合并

**适用场景**：
- 独立区域的简单修改
- 没有逻辑冲突
- 冲突标记清晰（1-2 处）

**操作步骤**：
```bash
# 1. 开始合并
git merge upstream/main

# 2. 查看冲突
git status
git diff --name-only --diff-filter=U

# 3. 编辑冲突文件
# 保留需要的部分，删除冲突标记
vim conflicted_file.go

# 4. 标记已解决
git add conflicted_file.go

# 5. 完成合并
git commit
```

**何时使用**：
- ✅ 文档文件（README.md, CHANGELOG.md）
- ✅ 配置文件（简单的键值对修改）
- ✅ 独立功能的新增代码

**何时避免**：
- ❌ 核心逻辑文件有多处冲突
- ❌ 架构重构导致的冲突
- ❌ 不确定如何合并的复杂逻辑

---

### 策略 2: 手动 Edit（推荐首选）

**适用场景**：
- 多个独立的冲突区域
- 需要精确控制合并结果
- 需要审计合并逻辑

**操作步骤**：
```bash
# 1. 查看冲突内容
git diff upstream/main -- file.go

# 2. 决定基础版本（通常选择架构更新的一方）
git checkout --theirs file.go  # 或 --ours

# 3. Read 当前文件状态
Read file.go

# 4. 用 Edit 精确插入另一方的逻辑
Edit file.go "precise_old_string" "merged_new_string"

# 5. 立即编译验证
go build ./cmd/server

# 6. 标记已解决
git add file.go
```

**优点**：
- 精确控制每一处修改
- 可审计（每个 Edit 都有明确的 old/new）
- 易于回滚（git reset 到某次 Edit 前）
- 避免脚本引入的微妙错误

**案例**：provider.go 的 4 个冲突区域
```bash
# 接受上游完整内容
git checkout --theirs provider.go

# 精确插入 DashScope 逻辑
Edit provider.go "old_1" "new_1"  # 插入 requirePublicURL 条件
Edit provider.go "old_2" "new_2"  # 插入 allowed 表条目
Edit provider.go "old_3" "new_3"  # 插入图片路由分支
Edit provider.go "old_4" "new_4"  # 插入视频路由分支

# 验证
go build ./cmd/server
```

---

### 策略 3: 逐文件对比

**适用场景**：
- 不确定哪个版本更好
- 结构性重构，难以直接合并
- 需要理解双方的意图

**操作步骤**：
```bash
# 1. 保存双方版本
git show HEAD:path/to/file.go > /tmp/ours.go
git show upstream/main:path/to/file.go > /tmp/theirs.go

# 2. 对比差异
diff -u /tmp/ours.go /tmp/theirs.go | less

# 3. 理解双方的修改意图
# - 哪些是 bug 修复？
# - 哪些是功能增强？
# - 哪些是架构优化？

# 4. 决定合并策略
# - 接受某一方为基础
# - 手动合并关键部分

# 5. 实施合并（通常用 Edit）
git checkout --theirs file.go
Edit file.go ...
```

**何时使用**：
- ✅ 核心架构文件（service.go 从 2118 → 363 行）
- ✅ 接口定义发生变化
- ✅ 双方都有重要修改

---

### 策略 4: 理解新架构

**适用场景**：
- 上游进行了架构级重构
- API 签名变更
- 依赖关系重组

**操作步骤**：
```bash
# 1. 识别架构变更
git log upstream/main --oneline | grep -i "refactor\|restructure\|rename"

# 2. 阅读关键提交
git show <commit_hash>

# 3. 理解新架构
# - 哪些方法被移动了？（Service → Repository）
# - 哪些方法被重命名了？（providerResourceURL → directResourceURL）
# - 新增了哪些抽象层？（Provider interface, Registry）

# 4. 修改本地代码适配新架构
Edit file.go "s.OldMethod(...)" "s.repo.NewMethod(...)"

# 5. 编译验证
go build
```

**案例**：provider_ark_private_assets.go
```go
// 旧架构（上游 bug）
s.ResourceForUser(&model.User{ID: userID}, resourceID)
s.providerResourceURL(resource, expiry)

// 新架构（正确）
s.repo.ResourceForUser(userID, resourceID)
s.directResourceURL(resource, expiry)
```

---

## 禁止的操作

### ❌ 批量 `git checkout`

**错误示例**：
```bash
# 危险：覆盖所有已合并的文件
for f in *.go; do
    if [[ ! "$f" =~ dashscope ]]; then
        git checkout upstream/main -- "$f"
    fi
done
```

**后果**：
- 覆盖合并提交中已正确解决的改动
- 编译可能通过，但功能静默失效
- 极难发现（需要逐行对比合并提交）

**正确做法**：
```bash
# 只更新确实有问题的特定文件
git checkout upstream/main -- specific_broken_file.go

# 或使用 stash 保护
git stash push -m "保护合并结果"
git checkout upstream/main -- file.go
git stash pop  # 如果有冲突，手动解决
```

---

### ❌ 盲目接受 ours/theirs

**错误示例**：
```bash
# 危险：不看内容直接接受
git checkout --ours .
# 或
git checkout --theirs .
```

**后果**：
- 丢失另一方的功能
- 可能丢失重要的 bug 修复
- 可能丢失架构改进

**正确做法**：
```bash
# 逐个文件决策
git checkout --theirs path/to/file1.go  # 这个文件用上游
Edit path/to/file1.go ...               # 插入本地逻辑

git checkout --ours path/to/file2.go    # 这个文件保留本地
Edit path/to/file2.go ...               # 插入上游改进
```

---

### ❌ 基于记忆编辑代码

**错误示例**：
```bash
# 危险：不 Read，直接凭记忆 Edit
Edit file.go "我记得的内容" "新内容"
# 结果：String to replace not found
```

**后果**：
- Edit 失败，浪费一次工具调用
- 可能编辑错位置（记忆不准确）
- 在合并过程中文件状态不断变化

**正确做法**：
```bash
# 1. 先 Read 确认当前状态
Read file.go

# 2. 找到精确的 old_string（复制粘贴）

# 3. 再 Edit
Edit file.go "精确的old_string" "new_string"
```

---

### ❌ 不编译就继续

**错误示例**：
```bash
Edit file1.go ...
Edit file2.go ...
Edit file3.go ...
# 最后才编译
go build  # 发现一堆错误，不知道哪个 Edit 引入的
```

**后果**：
- 问题累积，难以定位
- 可能在错误的方向上继续
- 浪费时间

**正确做法**：
```bash
# 每个文件 Edit 后立即验证
Edit file1.go ...
go build ./cmd/server  # ✅ 通过

Edit file2.go ...
go build ./cmd/server  # ❌ 失败，立即修复

Edit file3.go ...
go build ./cmd/server  # ✅ 通过
```

---

## 验证步骤

每解决一个冲突后，按顺序执行：

### 1. 标记已解决
```bash
git add <file>
```

### 2. 编译验证
```bash
# Go
cd backend && go build ./cmd/server

# TypeScript
cd web && npm run build
```

### 3. 静态检查
```bash
# Go
go vet ./...

# TypeScript
npm run typecheck
```

### 4. 测试验证（可选）
```bash
# Go
go test ./...

# 前端
npm test
```

### 5. 功能检查
- 协议注册完整性 → [[protocol-registration-checklist]]
- 调用链完整性（grep 检查函数引用）
- 本地功能未被覆盖

---

## 决策流程图

```
遇到冲突
  ↓
是否为简单文本冲突？
  ├─ 是 → Git 自动合并
  └─ 否 → 继续判断
      ↓
是否为多区域独立冲突？
  ├─ 是 → 手动 Edit（首选）
  └─ 否 → 继续判断
      ↓
是否为架构重构冲突？
  ├─ 是 → 逐文件对比 + 理解新架构
  └─ 否 → 继续判断
      ↓
是否为批量文件冲突？
  ├─ 是 → 逐个处理（禁用批量操作）
  └─ 否 → 寻求帮助
```

---

## 常见场景示例

### 场景 1: 协议常量冲突

**冲突**：
```
<<<<<<< HEAD
    ChannelInterfaceDashScopeImage = "dashscope-image"
=======
    ChannelInterfaceGeminiImage = "gemini-image"
>>>>>>> upstream/main
```

**策略**：Git 自动合并（保留双方）
```go
    ChannelInterfaceDashScopeImage = "dashscope-image"
    ChannelInterfaceGeminiImage = "gemini-image"
```

---

### 场景 2: 函数签名变更

**冲突**：本地调用了被重命名的方法

**策略**：理解新架构 + 手动 Edit
```bash
# 1. 查找旧方法的新位置
git log upstream/main -p -S "ResourceForUser"

# 2. 理解变更
# Service.ResourceForUser → Repository.ResourceForUser

# 3. 修改调用
Edit file.go "s.ResourceForUser(...)" "s.repo.ResourceForUser(...)"
```

---

### 场景 3: service.go 重构冲突

**冲突**：本地 2118 行，上游 363 行（职责分离）

**策略**：逐文件对比 + 接受上游
```bash
# 1. 对比理解
diff <(git show HEAD:service.go) <(git show upstream/main:service.go)

# 2. 确认本地修改是否在其他文件
grep -r "本地特有函数" backend/internal/service/

# 3. 决策：接受上游（架构更优）
git checkout --theirs service.go

# 4. 如果本地有独立功能，单独保留
# （本例中 DashScope 在独立文件，无需特殊处理）
```

---

## 相关链接

- [[large-scale-merge-best-practices]] - 大规模合并最佳实践
- [[protocol-registration-checklist]] - 协议注册检查清单
- [[api-development-best-practices]] - API 开发最佳实践

---

**Why**: 不同冲突需要不同策略，选对方法事半功倍。

**How to apply**: 
- 遇到冲突时先判断类型
- 根据矩阵选择策略
- 每步验证，不累积问题
- 禁用高风险操作（批量 checkout、盲目 ours/theirs）
