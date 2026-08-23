# 阶段 3 冲突系统性分析报告

**生成时间**: 2026-08-23 16:10  
**合并分支**: merge/upstream-integration-20260823-160503  
**冲突文件**: 12 个

---

## 一、冲突全局概览

### 1.1 冲突分类

| 类型 | 数量 | 文件 |
|---|---|---|
| **协议注册冲突** | 3 | provider.go, models.go, model-protocols.ts |
| **新增文件冲突** | 3 | bailian/discovery.go, interface.go, registry.go |
| **服务层冲突** | 4 | admin.go, channel_models.go, service.go, channel_model_catalog_plugin.go |
| **前端冲突** | 1 | wallet.ts |
| **文档冲突** | 1 | CHANGELOG.md |

### 1.2 冲突根源分析

**核心矛盾**：

1. **协议系统并行开发**
   - 本地: 实现 DashScope (dashscope-image, dashscope-video)
   - 上游: 实现 MiniMax (minimax-video) 和 Gemini (gemini-image)
   - 冲突点: 协议注册表

2. **架构重构**
   - 上游: Provider 抽象层重构 (interface.go, registry.go)
   - 上游: 百炼模型发现机制 (bailian/discovery.go)
   - 本地: 基于旧架构实现 DashScope

3. **服务层独立演进**
   - 本地: DashScope 能力配置
   - 上游: SKU 价格系统、插件重命名

---

## 二、逐个文件详细分析

### 2.1 核心冲突：provider.go

#### 冲突位置 1: requirePublicURL 逻辑

**本地版本**:
```go
// DashScope 视频接口要求 media[].url 为公网 HTTP(S) 地址
requirePublicURL := input.Config.InterfaceType == "newapi-channel-1" || 
                   input.Config.InterfaceType == "newapi-channel-2" || 
                   input.Config.InterfaceType == string(model.ChannelInterfaceVolcengineArkVideo) || 
                   input.Config.InterfaceType == string(model.ChannelInterfaceDashScopeVideo)
```

**上游版本**:
```go
// 简化版，没有 DashScope
requirePublicURL := input.Config.InterfaceType == "newapi-channel-1" || 
                   input.Config.InterfaceType == "newapi-channel-2" || 
                   input.Config.InterfaceType == string(model.ChannelInterfaceVolcengineArkVideo)
```

**冲突性质**: 功能添加冲突

**影响**: DashScope 视频参考素材 URL 处理

---

#### 冲突位置 2: 协议注册表

**本地版本**:
```go
allowed := map[string]map[string]bool{
    "image": {
        "openai-image": true, 
        "grok-image": true, 
        "volcengine-ark-image": true, 
        "volcengine-jimeng-image": true, 
        "dashscope-image": true  // 本地新增
    },
    "video": {
        "newapi": true, 
        "newapi-channel-1": true, 
        "newapi-channel-2": true, 
        "xai-video": true, 
        "volcengine-ark-video": true, 
        "volcengine-jimeng-video": true, 
        "gemini-veo": true, 
        "novita-video": true, 
        "dashscope-video": true  // 本地新增
    },
}
```

**上游版本**:
```go
allowed := map[string]map[string]bool{
    "image": {
        "openai-image": true, 
        "grok-image": true, 
        "volcengine-ark-image": true, 
        "volcengine-jimeng-image": true, 
        "gemini-image": true  // 上游新增
    },
    "video": {
        "newapi": true, 
        "newapi-channel-1": true, 
        "newapi-channel-2": true, 
        "xai-video": true, 
        "volcengine-ark-video": true, 
        "volcengine-jimeng-video": true, 
        "gemini-veo": true, 
        "novita-video": true, 
        "minimax-video": true  // 上游新增
    },
}
```

**冲突性质**: 协议列表合并冲突

**影响**: 缺少任一方都会导致对应功能不可用

---

### 2.2 核心冲突：models.go

#### 冲突位置: ChannelInterfaceType 常量

**本地版本**:
```go
const (
    // ... 其他协议 ...
    ChannelInterfaceNovitaVideo           ChannelInterfaceType = "novita-video"
    ChannelInterfaceDashScopeImage        ChannelInterfaceType = "dashscope-image"  // 本地新增
    ChannelInterfaceDashScopeVideo        ChannelInterfaceType = "dashscope-video"  // 本地新增
)
```

**上游版本**:
```go
const (
    // ... 其他协议 ...
    ChannelInterfaceNovitaVideo           ChannelInterfaceType = "novita-video"
    ChannelInterfaceMiniMaxVideo          ChannelInterfaceType = "minimax-video"  // 上游新增
)
```

**冲突性质**: 常量定义位置冲突

**影响**: 协议常量缺失会导致编译错误

---

### 2.3 核心冲突：model-protocols.ts (前端)

#### 冲突位置 1: 类型定义

**本地版本**:
```typescript
export type ModelProtocol = 
    | "openai-image"
    | "grok-image"
    | "volcengine-ark-image"
    | "volcengine-jimeng-image"
    | "dashscope-image"    // 本地新增
    | "dashscope-video"    // 本地新增
    | "openai-audio"
    // ...
```

**上游版本**:
```typescript
export type ModelProtocol = 
    | "openai-image"
    | "grok-image"
    | "volcengine-ark-image"
    | "volcengine-jimeng-image"
    | "gemini-image"       // 上游新增
    | "openai-audio"
    // ...
```

**冲突性质**: 类型联合冲突

---

#### 冲突位置 2: 协议配置数组

**本地版本**:
```typescript
export const MODEL_PROTOCOLS: ModelProtocolInfo[] = [
    // ...
    { value: "volcengine-jimeng-image", label: "...", ... },
    { value: "dashscope-image", label: "DashScope Wanx", ... },     // 本地新增
    { value: "dashscope-video", label: "DashScope Wan/HappyHorse", ... },  // 本地新增
    // ...
]
```

**上游版本**:
```typescript
export const MODEL_PROTOCOLS: ModelProtocolInfo[] = [
    // ...
    { value: "volcengine-jimeng-image", label: "...", ... },
    { value: "gemini-image", label: "Gemini Imagen", ... },  // 上游新增
    // ...
]
```

**冲突性质**: 数组插入位置冲突

---

### 2.4 架构冲突：新增文件 (AA 状态)

#### 文件 1: backend/internal/provider/interface.go

**状态**: AA (both added)

**本地版本**: 不存在此文件（基于旧架构）

**上游版本**: 新建 Provider 抽象接口
```go
package provider

type Provider interface {
    Name() string
    Type() string
    FetchModels(ctx context.Context, config Config) ([]Model, error)
}
```

**冲突性质**: 架构重构 - 本地未感知新架构

---

#### 文件 2: backend/internal/provider/registry.go

**状态**: AA (both added)

**上游新增**: Provider 注册表
```go
package provider

var registry = make(map[string]Provider)

func Register(name string, provider Provider) {
    registry[name] = provider
}
```

**冲突性质**: 架构重构 - 本地未使用注册机制

---

#### 文件 3: backend/internal/provider/bailian/discovery.go

**状态**: AA (both added)

**上游新增**: 百炼模型自动发现
```go
package bailian

func DiscoverModels(ctx context.Context, config Config) ([]Model, error) {
    // 从百炼 API 获取模型列表
    // 支持自动发现模型能力
}
```

**冲突性质**: 新功能 - 本地不涉及

---

### 2.5 服务层冲突

#### 文件 1: service.go

**冲突点**:
1. CreateTask 函数签名变化（上游重构）
2. provider_dashscope.go 和 provider_dashscope_video.go 的 init() 调用

**本地版本**: 
```go
// 在文件末尾调用
func init() {
    // 注册 DashScope 协议处理函数
}
```

**上游版本**: 
- 可能移除或重构了 init 机制
- 改用 registry.Register()

---

#### 文件 2: admin.go

**冲突点**: 模型能力配置 API

**本地**: DashScope 能力字段
**上游**: SKU 价格系统字段

---

#### 文件 3: channel_models.go

**冲突点**: 渠道模型管理逻辑

**本地**: DashScope 渠道配置
**上游**: 插件系统重命名（yingce）

---

#### 文件 4: channel_model_catalog_plugin.go

**状态**: AA (both added)

**上游新增**: 插件目录管理
**本地**: 不存在此文件

---

### 2.6 其他冲突

#### wallet.ts (前端)

**冲突点**: 钱包 API 接口

**本地**: 可能有自定义修改
**上游**: API 结构调整

---

#### CHANGELOG.md

**冲突点**: 变更日志条目

**本地**: DashScope 相关条目
**上游**: 其他功能条目

**影响**: 低（文档冲突）

---

## 三、冲突影响评估

### 3.1 按影响程度分类

#### 🔴 高影响（必须正确解决）

1. **provider.go** - 协议注册表
   - 缺失任一方会导致功能不可用
   - 影响: DashScope, MiniMax, Gemini 全部功能

2. **models.go** - 常量定义
   - 缺失会导致编译错误
   - 影响: 后端编译

3. **model-protocols.ts** - 类型定义
   - 缺失会导致前端类型错误
   - 影响: 前端编译和运行时

---

#### 🟡 中影响（需要仔细处理）

4. **service.go** - 服务层逻辑
   - 需要确保 DashScope 初始化正确
   - 影响: DashScope 功能启动

5. **admin.go** - 管理接口
   - 需要合并能力配置字段
   - 影响: 管理后台功能

6. **channel_models.go** - 渠道管理
   - 需要合并渠道配置逻辑
   - 影响: 渠道配置功能

---

#### 🟢 低影响（相对安全）

7-9. **新增文件** (interface.go, registry.go, discovery.go)
   - 上游新架构，本地不冲突
   - 保留上游版本即可

10. **channel_model_catalog_plugin.go**
    - 上游新功能
    - 保留上游版本即可

11. **wallet.ts**
    - API 调整
    - 需要检查兼容性

12. **CHANGELOG.md**
    - 文档冲突
    - 手动合并即可

---

### 3.2 风险矩阵

| 文件 | 影响程度 | 复杂度 | 风险等级 |
|---|---|---|---|
| provider.go | 🔴 高 | 🔴 高 | **🔴 极高** |
| models.go | 🔴 高 | 🟡 中 | **🔴 高** |
| model-protocols.ts | 🔴 高 | 🟡 中 | **🔴 高** |
| service.go | 🟡 中 | 🔴 高 | **🟡 高** |
| admin.go | 🟡 中 | 🟡 中 | **🟡 中** |
| channel_models.go | 🟡 中 | 🟡 中 | **🟡 中** |
| 新增文件 (3个) | 🟢 低 | 🟢 低 | **🟢 低** |
| channel_model_catalog_plugin.go | 🟢 低 | 🟢 低 | **🟢 低** |
| wallet.ts | 🟢 低 | 🟡 中 | **🟢 低** |
| CHANGELOG.md | 🟢 低 | 🟢 低 | **🟢 低** |

---

## 四、系统性问题识别

### 4.1 架构不一致

**问题**: 
- 本地基于旧架构实现 DashScope
- 上游已重构为新 Provider 架构

**影响**:
- DashScope 未使用新的 registry.Register()
- 可能与未来架构不兼容

**建议**: 
- 先按旧架构合并（保证功能）
- 后续重构为新架构（技术债）

---

### 4.2 协议注册分散

**问题**: 
- 协议定义散布在多个文件
- provider.go, models.go, model-protocols.ts 都需要修改

**影响**:
- 容易遗漏某个位置
- 维护成本高

**建议**: 
- 建立检查清单，确保所有位置都更新

---

### 4.3 测试覆盖

**问题**: 
- 合并后需要全面测试
- 特别是协议注册相关功能

**建议**: 
- 合并后立即测试所有协议
- DashScope, MiniMax, Gemini 都要验证

---

## 五、冲突解决策略

### 5.1 总体原则

1. **保留双方功能**: DashScope + 上游新增（MiniMax, Gemini）
2. **优先上游架构**: 新增文件使用上游版本
3. **合并注册表**: 协议列表包含所有协议
4. **保守处理**: 不确定时保留双方逻辑

---

### 5.2 分组解决策略

#### 第一组：协议注册（优先级 🔴）

**文件**: provider.go, models.go, model-protocols.ts

**策略**: 
1. 合并协议列表（包含所有协议）
2. 保留双方的特殊处理逻辑

**顺序**: 
1. models.go（定义常量）
2. provider.go（注册协议）
3. model-protocols.ts（前端类型）

---

#### 第二组：新增文件（优先级 🟢）

**文件**: interface.go, registry.go, discovery.go, channel_model_catalog_plugin.go

**策略**: 
- 完全接受上游版本
- 这些是新架构，本地没有冲突

---

#### 第三组：服务层（优先级 🟡）

**文件**: service.go, admin.go, channel_models.go

**策略**: 
1. 仔细审查上游变更
2. 确保 DashScope 逻辑完整
3. 合并双方的功能字段

---

#### 第四组：低风险（优先级 🟢）

**文件**: wallet.ts, CHANGELOG.md

**策略**: 
- wallet.ts: 检查 API 兼容性，优先上游
- CHANGELOG.md: 手动合并双方条目

---

## 六、详细解决方案

### 6.1 provider.go 解决方案

#### 冲突 1: requirePublicURL

**解决方案**: 合并双方逻辑

```go
requirePublicURL := input.Config.InterfaceType == "newapi-channel-1" || 
                   input.Config.InterfaceType == "newapi-channel-2" || 
                   input.Config.InterfaceType == string(model.ChannelInterfaceVolcengineArkVideo) ||
                   input.Config.InterfaceType == string(model.ChannelInterfaceDashScopeVideo)  // 保留本地
```

**风险**: 低
**验证**: 测试 DashScope 视频参考素材

---

#### 冲突 2: 协议注册表

**解决方案**: 合并所有协议

```go
allowed := map[string]map[string]bool{
    "image": {
        "openai-image": true, 
        "grok-image": true, 
        "volcengine-ark-image": true, 
        "volcengine-jimeng-image": true, 
        "dashscope-image": true,  // 本地
        "gemini-image": true,     // 上游
    },
    "video": {
        "newapi": true, 
        "newapi-channel-1": true, 
        "newapi-channel-2": true, 
        "xai-video": true, 
        "volcengine-ark-video": true, 
        "volcengine-jimeng-video": true, 
        "gemini-veo": true, 
        "novita-video": true, 
        "dashscope-video": true,  // 本地
        "minimax-video": true,    // 上游
    },
    "audio": {"openai-audio": true, "async-audio": true},
}
```

**风险**: 低（只要包含所有协议）
**验证**: 测试所有协议的图片/视频生成

---

### 6.2 models.go 解决方案

**解决方案**: 合并常量定义

```go
const (
    // ... 其他协议 ...
    ChannelInterfaceNovitaVideo           ChannelInterfaceType = "novita-video"
    ChannelInterfaceDashScopeImage        ChannelInterfaceType = "dashscope-image"  // 本地
    ChannelInterfaceDashScopeVideo        ChannelInterfaceType = "dashscope-video"  // 本地
    ChannelInterfaceMiniMaxVideo          ChannelInterfaceType = "minimax-video"    // 上游
)
```

**风险**: 低（常量定义）
**验证**: 编译检查

---

### 6.3 model-protocols.ts 解决方案

#### 类型定义

**解决方案**: 合并类型联合

```typescript
export type ModelProtocol = 
    | "openai-image"
    | "grok-image"
    | "volcengine-ark-image"
    | "volcengine-jimeng-image"
    | "dashscope-image"    // 本地
    | "dashscope-video"    // 本地
    | "gemini-image"       // 上游
    | "openai-audio"
    // ...
```

---

#### 协议配置数组

**解决方案**: 合并配置项（按字母顺序）

```typescript
export const MODEL_PROTOCOLS: ModelProtocolInfo[] = [
    // ... 前面的协议 ...
    { value: "dashscope-image", label: "DashScope Wanx", capability: "image", ... },
    { value: "dashscope-video", label: "DashScope Wan/HappyHorse", capability: "video", ... },
    { value: "gemini-image", label: "Gemini Imagen", capability: "image", ... },
    // ... 后面的协议 ...
]
```

**风险**: 低
**验证**: 前端编译 + UI 检查

---

### 6.4 service.go 解决方案

**策略**: 仔细审查上游变更

**步骤**:
1. 接受上游的函数签名变更
2. 确保 provider_dashscope.go 的 init() 调用仍然有效
3. 如果上游移除了 init 机制，改用新的注册方式

**风险**: 中（需要理解架构变化）
**验证**: 启动服务，检查 DashScope 是否正常注册

---

### 6.5 新增文件解决方案

**文件**: 
- interface.go
- registry.go
- discovery.go
- channel_model_catalog_plugin.go

**解决方案**: 完全接受上游版本

```bash
git checkout --theirs backend/internal/provider/interface.go
git checkout --theirs backend/internal/provider/registry.go
git checkout --theirs backend/internal/provider/bailian/discovery.go
git checkout --theirs backend/internal/service/channel_model_catalog_plugin.go
git add <files>
```

**风险**: 低（新文件，无冲突）

---

### 6.6 admin.go 解决方案

**策略**: 合并能力配置字段

**步骤**:
1. 接受上游的结构变更
2. 确保 DashScope 能力字段仍在
3. 合并上游的 SKU 相关字段

**风险**: 中
**验证**: 管理后台能力配置功能

---

### 6.7 channel_models.go 解决方案

**策略**: 合并双方逻辑

**步骤**:
1. 接受上游的插件重命名
2. 保留 DashScope 渠道配置逻辑

**风险**: 中
**验证**: 渠道配置功能

---

### 6.8 wallet.ts 解决方案

**策略**: 优先上游，检查兼容性

**步骤**:
1. 接受上游的 API 变更
2. 检查是否影响现有功能
3. 如有必要，调整调用代码

**风险**: 低
**验证**: 钱包功能测试

---

### 6.9 CHANGELOG.md 解决方案

**策略**: 手动合并

**步骤**:
1. 保留双方的条目
2. 按时间顺序排列
3. 格式统一

**风险**: 极低（文档）

---

## 七、执行计划

### 7.1 推荐顺序

**阶段 1: 低风险文件（快速处理）** ⏱️ 5 分钟

1. interface.go - `git checkout --theirs`
2. registry.go - `git checkout --theirs`
3. discovery.go - `git checkout --theirs`
4. channel_model_catalog_plugin.go - `git checkout --theirs`
5. CHANGELOG.md - 手动合并

---

**阶段 2: 协议注册（核心）** ⏱️ 20 分钟

6. models.go - 合并常量定义
7. provider.go - 合并协议列表
8. model-protocols.ts - 合并类型定义

---

**阶段 3: 服务层（仔细处理）** ⏱️ 30 分钟

9. service.go - 审查架构变更
10. admin.go - 合并能力字段
11. channel_models.go - 合并渠道逻辑

---

**阶段 4: 前端其他（最后）** ⏱️ 10 分钟

12. wallet.ts - 检查兼容性

---

**总预计时间**: 1-1.5 小时

---

### 7.2 验证清单

#### 编译验证
- [ ] Go 后端编译通过
- [ ] TypeScript 前端编译通过
- [ ] 无类型错误

#### 功能验证
- [ ] DashScope 图片生成
- [ ] DashScope 视频生成
- [ ] MiniMax 视频生成（上游新增）
- [ ] Gemini 图片生成（上游新增）
- [ ] 管理后台能力配置
- [ ] 渠道配置功能
- [ ] 钱包功能

#### 架构验证
- [ ] DashScope 协议正确注册
- [ ] provider_dashscope.go 初始化成功
- [ ] provider_dashscope_video.go 初始化成功

---

## 八、风险预警

### 8.1 高风险操作

⚠️ **provider.go 协议列表**
- 风险: 遗漏任一协议会导致功能不可用
- 缓解: 使用清单检查所有协议

⚠️ **service.go 架构变更**
- 风险: DashScope 可能无法正确初始化
- 缓解: 测试启动日志，确认注册成功

---

### 8.2 常见错误

❌ **只保留一方的协议列表**
```go
// 错误：只有 DashScope，缺少 MiniMax
"video": {"dashscope-video": true}
```

✅ **正确：合并所有协议**
```go
// 正确：包含所有协议
"video": {"dashscope-video": true, "minimax-video": true, ...}
```

---

❌ **忽略架构变更**
```go
// 错误：假设 init() 仍然有效
func init() {
    // 可能已不被调用
}
```

✅ **正确：验证注册机制**
```go
// 检查日志确认注册成功
// 或改用新的 registry.Register()
```

---

### 8.3 回滚方案

如果合并失败：

```bash
# 1. 中止合并
git merge --abort

# 2. 返回原分支
git checkout main

# 3. 删除合并分支
git branch -D merge/upstream-integration-*

# 4. 检查备份
git checkout backup-before-upstream-merge-20260823-133037
```

---

## 九、后续建议

### 9.1 技术债处理

**问题**: DashScope 基于旧架构

**建议**: 合并后创建任务
```
TODO: 重构 DashScope 为新 Provider 架构
- 实现 provider.Interface
- 使用 registry.Register()
- 支持模型自动发现
```

---

### 9.2 协议管理优化

**问题**: 协议定义分散

**建议**: 建立协议注册检查清单
```
新增协议时必须更新：
- [ ] backend/internal/model/models.go (常量)
- [ ] backend/internal/service/provider.go (注册表)
- [ ] web/src/lib/model-protocols.ts (类型+配置)
- [ ] 协议处理文件 (provider_*.go)
```

---

### 9.3 自动化测试

**建议**: 增加协议注册测试
```go
func TestAllProtocolsRegistered(t *testing.T) {
    // 检查所有定义的协议都已注册
    // 防止遗漏
}
```

---

## 十、总结

### 10.1 冲突特征

- **数量**: 12 个文件
- **类型**: 功能添加 + 架构重构
- **复杂度**: 中高（需要仔细处理）
- **风险**: 可控（有明确解决方案）

---

### 10.2 关键点

1. ✅ **协议列表必须完整**: 包含所有协议
2. ✅ **新增文件接受上游**: 新架构文件
3. ✅ **服务层仔细审查**: 确保 DashScope 正常
4. ✅ **全面测试验证**: 所有协议都要测试

---

### 10.3 预期结果

**合并后**:
- ✅ DashScope 功能完整（本地）
- ✅ MiniMax 视频可用（上游）
- ✅ Gemini 图片可用（上游）
- ✅ 新 Provider 架构就位（上游）
- ✅ SKU 价格系统可用（上游）

---

**报告完成**  
**下一步**: 等待您审阅后，开始执行解决方案

---

## 附录：快速参考

### A. 冲突文件清单

```
1. CHANGELOG.md                                          [低风险]
2. backend/internal/model/models.go                      [高风险]
3. backend/internal/provider/bailian/discovery.go        [低风险-新增]
4. backend/internal/provider/interface.go                [低风险-新增]
5. backend/internal/provider/registry.go                 [低风险-新增]
6. backend/internal/service/admin.go                     [中风险]
7. backend/internal/service/channel_model_catalog_plugin.go  [低风险-新增]
8. backend/internal/service/channel_models.go            [中风险]
9. backend/internal/service/provider.go                  [高风险]
10. backend/internal/service/service.go                  [中风险]
11. web/src/lib/model-protocols.ts                       [高风险]
12. web/src/services/api/wallet.ts                       [低风险]
```

### B. 验证命令

```bash
# 编译检查
cd backend && go build
cd web && npm run build

# 启动检查
docker-compose up -d
docker-compose logs | grep -i dashscope

# 测试 DashScope
curl -X POST http://localhost:3000/api/generate \
  -d '{"model":"wanx-v1","interface":"dashscope-image",...}'
```

### C. 关键代码位置

```
协议注册:
- backend/internal/model/models.go:74-79
- backend/internal/service/provider.go:310-320
- web/src/lib/model-protocols.ts:5-15

DashScope 实现:
- backend/internal/service/provider_dashscope.go
- backend/internal/service/provider_dashscope_video.go

新架构:
- backend/internal/provider/interface.go
- backend/internal/provider/registry.go
```
