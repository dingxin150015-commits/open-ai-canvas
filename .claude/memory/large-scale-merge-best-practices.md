---
name: large-scale-merge-best-practices
description: 合并大量上游提交（300+）的经验和教训
metadata:
  type: feedback
---

从本次合并 455 个上游提交的经验中学到的关键教训。

## 核心原则

### 1. 永远先 Read，再 Edit

**问题**：在没有 Read 的情况下，基于 Git diff 或记忆重构字符串。

**教训**：
- Git diff 显示的是"变化"，不是"完整文本"
- 上下文可能包含微妙差异（如"即梦图片" vs "即梦官方图片"）
- Edit 工具需要精确匹配，任何微小差异都会导致失败

**正确做法**：
```bash
# 1. 先 Read 文件
Read file_path

# 2. 找到精确的 old_string（复制粘贴，不要手打）

# 3. 再 Edit
Edit file_path old_string new_string
```

**Why**: 避免 "String to replace not found" 错误，确保编辑的是真实存在的内容。

---

### 2. 禁用批量 checkout 已合并文件

**问题**：执行 `git checkout upstream/main -- *.go` 覆盖了已在合并提交中正确解决的改动。

**危险性**：
- 编译可能仍然通过
- 功能会静默失效（如协议验证、能力映射）
- 极难发现（需要逐行对比合并提交）

**正确做法**：
```bash
# ❌ 错误：批量 checkout
git checkout upstream/main -- *.go

# ✅ 正确：逐个检查
# 1. 列出需要更新的文件（基于编译错误）
go build ./cmd/server 2>&1 | grep "\.go:"

# 2. 只更新确实有 API 不兼容的文件
git checkout upstream/main -- path/to/specific/file.go

# 3. 或使用 stash 保护已解决的内容
git stash push -m "保护合并结果"
git checkout upstream/main -- specific_file.go
git stash pop
```

**Why**: 合并提交已经包含了双方的正确逻辑，批量覆盖会丢失本地的关键改动。

---

### 3. 每次修改后立即编译

**问题**：等到"全部改完"再编译，导致问题累积。

**教训**：
- 小错误会掩盖大错误
- 难以定位是哪次修改引入的问题
- 浪费时间在错误的方向上

**正确做法**：
```bash
# 每次 Edit 后立即验证
Edit file.go
go build ./cmd/server

# 发现问题立即修复
# 不要继续下一个文件
```

**Why**: 及早发现、及早修复，降低调试成本。

**链接**: [[docker-verify-deployment-before-testing]]

---

### 4. 复杂冲突用手动 Edit，不用脚本

**问题**：Python/AWK 脚本在多区域、复杂结构的冲突时不可靠。

**脚本的局限**：
- 正则难以匹配所有变体
- 难以处理嵌套结构
- 难以调试（黑盒）
- 可能引入微妙错误（空格、缩进）

**正确做法**：
```bash
# ❌ 错误：用脚本处理复杂冲突
python merge_script.py provider.go

# ✅ 正确：手动 Edit
# 1. 接受上游完整内容
git checkout --theirs provider.go

# 2. 用 Edit 工具精确插入本地逻辑
Edit provider.go "old_string" "new_string"  # 第 1 处
Edit provider.go "old_string" "new_string"  # 第 2 处
...

# 3. 编译验证
go build
```

**策略选择**：

| 冲突类型 | 推荐方法 |
|---|---|
| 简单文本冲突（1-2 处） | Git 自动合并 |
| 多区域独立冲突 | 手动 Edit（首选） |
| 结构性重构冲突 | 逐文件对比 + 手动 Edit |
| 批量文件冲突 | **禁用脚本**，逐个处理 |

**Why**: 手动 Edit 精确、可审计、可验证，复杂场景下更可靠。

---

### 5. 记录所有决策

**做法**：
```markdown
## 冲突文件: provider.go

### 决策
- 接受上游的新字段（AgentRequests, TextHistory）
- 接受上游的错误处理逻辑（更详细）
- 保留本地的 DashScope 协议注册

### 理由
- 上游新字段是架构升级，必须接受
- 上游错误处理有详细分类，优于本地简单版本
- DashScope 是本地独有功能，上游无此协议

### 方法
1. `git checkout --theirs provider.go`
2. Edit 插入 DashScope 分支（4 处）
3. 编译验证通过
```

**Why**: 便于审计、便于未来合并、便于团队理解决策。

---

## 检查清单

合并完成后必须全部通过：

### 冲突解决完整性
- [ ] 无冲突标记 (`grep -r "<<<<<<< HEAD" backend/ web/`)
- [ ] 无 .orig/.rej 文件 (`find . -name "*.orig" -o -name "*.rej"`)
- [ ] Git status 显示工作区干净

### 编译验证
- [ ] Go 编译通过 (`go build ./cmd/server`)
- [ ] Go vet 干净 (`go vet ./...`)
- [ ] TypeScript 编译通过 (`npm run build`)
- [ ] 前端测试运行 (`npm test` 或 `bun test`)

### 功能完整性
- [ ] 协议注册完整（检查 6 个位置）→ [[protocol-registration-checklist]]
- [ ] 本地功能未被覆盖（对比合并前 HEAD）
- [ ] 调用链完整（grep 检查函数引用）
- [ ] 实现文件存在且完整

### 文档
- [ ] 更新 CHANGELOG.md（如果有本地功能）
- [ ] 记录合并决策（哪些接受上游，哪些保留本地）
- [ ] 记录已知技术债

---

## 常见错误模式

### 错误 1: "编译通过 = 没问题"

**现实**：编译通过只意味着语法正确，不代表功能完整。

**案例**：
- 协议常量存在，但 validChannelInterfaceType 中未注册 → 管理后台保存失败
- 协议常量存在，但 capabilityForProtocol 中未映射 → 模型配置失效
- APIFormat 字段被覆盖，编译通过，但前端下拉框失效

**防范**：
- 使用检查清单系统性验证
- 运行功能测试，不只是编译测试

---

### 错误 2: "快速完成"比"正确完成"重要

**现实**：修复静默失效的功能会花费更多时间。

**案例**：
- 批量 checkout 节省 30 分钟，但发现和修复被覆盖功能花费 2 小时
- 用脚本节省 1 小时，但调试脚本错误花费 3 小时

**防范**：
- 慢即是快，正确的慢速方法快于错误的快速方法
- 每个步骤都验证，不要累积问题

---

### 错误 3: 依赖记忆和推测

**现实**：合并过程中文件状态不断变化，记忆不可靠。

**案例**：
- "我记得这个字段是这样的" → Edit 失败
- "应该没问题" → 运行时发现功能失效
- "肯定是这个版本" → 实际是另一个版本

**防范**：
- 永远基于当前文件状态（Read）
- 不假设，用命令验证（grep, git show）
- 记录决策而非依赖记忆

---

## 时间预估

| 规模 | 冲突文件 | 预估时间 | 实际案例 |
|---|---|---|---|
| 小型合并 | < 5 个 | 1-2 小时 | - |
| 中型合并 | 5-15 个 | 4-8 小时 | - |
| 大型合并 | 15+ 个 | 8-16 小时 | 本次 12 个冲突，7 小时 |

**影响因素**：
- 冲突复杂度（简单 vs 结构性重构）
- 熟悉度（第一次 vs 有经验）
- 验证深度（编译 vs 编译+测试+功能验证）

---

## 相关链接

- [[git-conflict-resolution-strategy]] - 冲突解决策略矩阵
- [[protocol-registration-checklist]] - 协议注册检查清单
- [[api-development-best-practices]] - API 开发最佳实践
- [[docker-verify-deployment-before-testing]] - 部署验证纪律

---

**Why**: 大规模合并极易引入静默失效，必须系统性验证。功能完整性比速度更重要。

**How to apply**: 
- 在每次大规模合并（100+ 提交）时使用此检查清单
- 不要假设"编译通过 = 没问题"
- 优先使用手动 Edit 而非自动化脚本
- 记录所有决策和理由
