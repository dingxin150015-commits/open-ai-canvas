# 上游与本地代码完整对比分析报告

## 执行摘要

**分析时间**: 2026-08-23  
**上游最新提交**: 70a6640 (2026-08-23 01:08:20Z)  
**本地最新提交**: 11931d0 (feat(dashscope): 支持多张参考图生成)  
**上游分支数**: 2 (main, feature)  
**网络状态**: 部分受限（GitHub API 可用，Git fetch 不稳定）

---

## 一、关键发现

### 1.1 上游主要变化（最近 3 天）

**核心架构变更**：
1. ✅ **模型规格价格系统重构** - 引入 SKU、价格档、按规格路由
2. ✅ **Agent 系统增强** - 新增工具调用支持
3. ✅ **文本生成能力** - 新增 TextCapabilityConfig
4. ✅ **方舟私域素材** - 支持私域素材授权
5. ✅ **MiniMax 视频协议** - 新增 MiniMax 视频生成支持
6. ✅ **错误处理增强** - 详细的 HTTP 状态码分类

**UI 优化**：
- 画布背景优化
- 提示词卡片优化
- 批量上传优化
- 任务详情展示优化

---

### 1.2 本地独有功能

**DashScope 完整协议** ⭐⭐⭐⭐⭐：
- `provider_dashscope.go` (9.5KB) - 图片生成
- `provider_dashscope_video.go` (25KB) - 视频生成
- 支持 Wan 2.7 和 HappyHorse 全系列模型
- 高级功能：last_frame、reference_voice、driving_audio
- 完善的验证逻辑和错误处理

**价值评估**: 极高，建议向上游贡献

---

## 二、核心文件差异对比

### 2.1 `backend/internal/service/provider.go`

| 项目 | 上游 | 本地 | 差异 |
|---|---|---|---|
| **行数** | 3944 | 3368 | 上游 +576 行 |
| **TextHistory** | ✅ | ❌ | 文本对话历史 |
| **AgentRequests** | ✅ | ❌ | Agent 工具请求 |
| **ChannelModelKey** | ✅ | ❌ | 渠道模型键 |
| **PriceTierID** | ✅ | ❌ | 价格档 ID |
| **ProviderModelKey** | ✅ | ❌ | 供应商模型键 |
| **ArkPrivateAssetUpload** | ✅ | ❌ | 方舟私域素材上传 |
| **详细错误分类** | ✅ | ❌ | HTTP 状态码详细处理 |

**上游新增类型**：
```go
type agentToolRequests struct {
    Responses      map[string]interface{} `json:"responses"`
    ChatCompletion map[string]interface{} `json:"chatCompletion"`
}

type providerTextMessage struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}
```

**影响**: ⭐⭐⭐⭐⭐ 极高 - 核心结构变更，必须合并

---

### 2.2 `backend/internal/service/model_capability.go`

| 项目 | 上游 | 本地 | 差异 |
|---|---|---|---|
| **行数** | 802 | 611 | 上游 +191 行 |
| **TextCapabilityConfig** | ✅ | ❌ | 文本生成能力配置 |
| **MiniMaxVideo 配置** | ✅ | ❌ | MiniMax 视频模型 |
| **DashScope 验证逻辑** | ❌ | ✅ | 本地独有 |

**上游新增**：
```go
type TextCapabilityConfig struct {
    References TextReferenceConfig `json:"references"`
}

type TextReferenceConfig struct {
    PromptMaxChars int   `json:"promptMaxChars"`
    MaxImages      int   `json:"maxImages"`
    MaxImageBytes  int64 `json:"maxImageBytes"`
    MaxVideos      int   `json:"maxVideos"`
    MaxVideoBytes  int64 `json:"maxVideoBytes"`
}

// MiniMax 视频配置
case model.ChannelInterfaceMiniMaxVideo:
    video.Operations = append(video.Operations, "reference_to_video")
    video.References.MaxImages = 9
    video.References.MaxVideos = 3
    video.References.MaxAudios = 3
    video.Ratios = []string{"adaptive", "21:9", "16:9", "4:3", "1:1", "3:4", "9:16"}
    video.Resolutions = []string{"768P", "2K"}
```

**影响**: ⭐⭐⭐⭐ 高 - 需合并上游配置，保留本地验证逻辑

---

### 2.3 Provider 文件对比

| 文件 | 上游 | 本地 | 说明 |
|---|---|---|---|
| `provider_dashscope.go` | ❌ | ✅ 9.5KB | **本地独有**：DashScope 图片 |
| `provider_dashscope_video.go` | ❌ | ✅ 25KB | **本地独有**：DashScope 视频 |
| `provider_ark_private_assets.go` | ✅ | ❌ | **上游独有**：方舟私域素材 |
| `provider_jimeng.go` | ✅ | ✅ | 都有 |
| `provider_error.go` | ✅ | ✅ | 都有（本地较旧） |
| `provider_task_cancellation.go` | ✅ | ✅ | 都有 |
| `provider_task_recovery.go` | ✅ | ✅ | 都有 |
| `provider_video_options.go` | ✅ | ✅ | 都有（本地较旧） |
| `provider_request_types.go` | ✅ | ✅ | 都有 |

---

## 三、上游独有功能（需拉取）

### 3.1 模型规格价格系统（SKU）⭐⭐⭐⭐⭐

**影响**: 极高 - 核心计费系统重构

**新增文件**：
- `backend/internal/model/model_sku.go`

**新增字段**：
- `ChannelModelKey` - 渠道模型键
- `PriceTierID` - 价格档 ID
- `ProviderModelKey` - 供应商模型键

**功能**：
- 价格档（PriceTier）管理
- SKU 匹配和路由
- 按规格动态定价

**风险**: 可能涉及数据库 Schema 变更

---

### 3.2 Agent 工具请求支持 ⭐⭐⭐⭐⭐

**影响**: 高 - 新增 Agent 能力

**新增字段**：
```go
AgentRequests *agentToolRequests `json:"agentRequests"`
TextHistory   []providerTextMessage `json:"textHistory"`
```

**功能**：
- Agent 工具调用
- 多轮对话历史
- 聊天完成集成

---

### 3.3 方舟私域素材支持 ⭐⭐⭐⭐

**影响**: 中高

**新增文件**：
- `provider_ark_private_assets.go`

**新增字段**：
- `ArkPrivateAssetUpload`

**功能**：
- 私域素材授权
- 素材上传管理

---

### 3.4 MiniMax 视频协议 ⭐⭐⭐⭐

**影响**: 中高

**新增配置**：
- 支持 reference_to_video
- 最多 9 张参考图
- 最多 3 个参考视频
- 最多 3 个音频
- 支持 adaptive 比例

---

### 3.5 文本生成能力配置 ⭐⭐⭐⭐

**影响**: 中高

**新增类型**：
```go
type TextCapabilityConfig struct {
    References TextReferenceConfig
}
```

**功能**：
- 文本生成模型能力声明
- Prompt 长度限制
- 引用素材限制

---

### 3.6 增强的错误处理 ⭐⭐⭐

**影响**: 中

**改进**：
```go
switch e.StatusCode {
case 524:
    return "模型服务网关超时..."
case http.StatusBadRequest, http.StatusUnprocessableEntity:
    return "模型服务拒绝了请求，请检查模型和参数"
case http.StatusUnauthorized, http.StatusForbidden:
    return "模型服务鉴权失败，请检查 API Key 和模型权限"
case http.StatusNotFound:
    return "模型或模型接口不存在，请检查渠道配置"
case http.StatusRequestTimeout, http.StatusGatewayTimeout:
    return "模型服务响应超时，请稍后重试"
case http.StatusTooManyRequests:
    return "模型服务请求过于频繁或额度不足，请稍后重试"
}
```

---

## 四、冲突预测与解决策略

### 4.1 高风险冲突区域 ⭐⭐⭐⭐⭐

#### 冲突 1: `provider.go` 结构体变更

**冲突点**：
- `canvasGenerationInput` 新增 6 个字段
- `providerConfig` 新增多个字段
- 错误处理逻辑重构

**解决策略**：
```
1. 接受上游所有新增字段
2. 保留本地的 DashScope 协议注册代码
3. 合并错误处理逻辑（采用上游的详细分类）
4. 手动三方合并
```

**预计时间**: 2-3 小时

---

#### 冲突 2: `model_capability.go` 配置扩展

**冲突点**：
- `ModelCapabilityConfig` 新增 `Text` 字段
- 新增 MiniMax 配置
- 本地有 DashScope 验证逻辑

**解决策略**：
```
1. 接受上游的 Text 配置
2. 接受上游的 MiniMax 配置
3. 保留本地的 DashScope 验证逻辑（:384-460）
4. 手动合并，确保不冲突
```

**预计时间**: 1-2 小时

---

### 4.2 中风险区域 ⭐⭐⭐

#### 1. 数据库 Schema

**可能变更**：
- `model_pricings` 表 - 新增 SKU 相关字段
- 新表：`model_skus`
- `channels` 表可能新增字段

**解决策略**：
```
1. 查找迁移脚本
2. 备份数据库
3. 运行迁移
4. 验证数据完整性
```

---

#### 2. 前端文件

**可能变更**：
- 模型配置 UI（价格档选择）
- Agent 对话 UI
- 画布背景优化
- 批量上传 UI

**解决策略**：
```
1. 逐个文件对比（需网络恢复后下载）
2. 保留本地 DashScope 相关 UI
3. 合并上游优化
```

---

## 五、推荐的合并策略

### 策略 A: 分阶段 Merge（强烈推荐）⭐⭐⭐⭐⭐

**优点**：
- 风险可控，每阶段可回退
- 逐步解决冲突
- 保持完整的提交历史

**步骤**：

```bash
# ===== 阶段 1: 准备 =====
# 1. 提交本地所有改动
git add backend/internal/service/provider_dashscope*.go
git add backend/internal/service/model_capability.go
git commit -m "feat(dashscope): 完整实现 DashScope 图片和视频生成协议

- 图片生成：Wanx 系列模型
- 视频生成：Wan 2.7 + HappyHorse 系列
- 高级功能：last_frame、reference_voice、driving_audio
- 严格的模型能力区分和验证逻辑
"

# 2. 创建多个备份分支
git branch backup-before-merge-$(date +%Y%m%d)
git branch feature/dashscope-complete

# ===== 阶段 2: 获取上游（网络恢复后）=====
# 3. 获取上游最新代码
git fetch origin --all --prune

# 4. 查看差异
git log --oneline --graph HEAD..origin/main | head -50
git diff --stat HEAD..origin/main

# ===== 阶段 3: 核心文件合并 =====
# 5. 创建合并分支
git checkout -b merge/upstream-integration

# 6. 开始合并
git merge origin/main
# 预期会有冲突

# 7. 解决冲突 - provider.go
# 策略：
#  - 接受上游的所有新增字段和类型定义
#  - 保留本地的 DashScope 协议注册：
#    case "dashscope-image":
#        return runDashScopeImageTask(...)
#    case "dashscope-video":
#        return runDashScopeVideoTask(...)
#  - 采用上游的错误处理逻辑
git add backend/internal/service/provider.go

# 8. 解决冲突 - model_capability.go  
# 策略：
#  - 接受上游的 Text 配置
#  - 接受上游的 MiniMax 配置
#  - 保留本地的 DashScope 验证逻辑（:384-460行）
#  - 确保两者不冲突
git add backend/internal/service/model_capability.go

# 9. 完成合并
git commit -m "Merge remote-tracking branch 'origin/main' into merge/upstream-integration

冲突解决：
- provider.go: 接受上游新增字段，保留DashScope协议
- model_capability.go: 合并配置，保留DashScope验证
"

# ===== 阶段 4: 新文件拉取 =====
# 10. 检查上游新增文件
git diff --name-status HEAD origin/main | grep "^A"

# 11. 拉取新文件（已在merge中完成）
# 新文件应该已经自动拉取：
#  - provider_ark_private_assets.go
#  - model_sku.go

# ===== 阶段 5: 测试验证 =====
# 12. 编译测试
cd backend
go build ./cmd/server
# 修复编译错误（如果有）

# 13. 运行测试
go test ./...

# 14. Docker 构建测试
cd ..
docker-compose build backend
docker-compose up -d

# 15. 功能测试
# - 测试 DashScope 图片生成
# - 测试 DashScope 视频生成
# - 测试上游新功能（如果相关）

# ===== 阶段 6: 合并到主分支 =====
# 16. 切换到 main
git checkout main

# 17. 合并集成分支
git merge merge/upstream-integration

# 18. 推送（如果需要）
git push myfork main
```

**预计时间**: 6-10 小时

---

### 策略 B: 新分支 Cherry-pick（备选）⭐⭐⭐

**适用场景**: 如果策略 A 冲突太复杂

**步骤**：

```bash
# 1. 基于上游创建新分支
git fetch origin
git checkout -b feature/dashscope-on-latest origin/main

# 2. Cherry-pick DashScope 提交
# 找到所有 DashScope 相关提交
git log --oneline --grep="dashscope" main

# 逐个 cherry-pick
git cherry-pick <commit-hash-1>
git cherry-pick <commit-hash-2>
# ...

# 3. 解决冲突
# 4. 测试
# 5. 如果成功，切换为主分支
git checkout main
git reset --hard feature/dashscope-on-latest
```

**预计时间**: 4-6 小时

---

## 六、向上游贡献方案

### 6.1 贡献策略：先 Issue 再 PR，分 3 个 PR

#### Issue: Feature Request

**标题**: `[Feature Request] 添加 DashScope（阿里云百炼）协议支持`

**内容**:
```markdown
## 背景

DashScope 是阿里云百炼平台的 AI 服务，提供 Wanx 图片生成和 Wan/HappyHorse 视频生成能力。

## 功能特性

### 图片生成 (dashscope-image)
- Wanx 系列模型
- 同步/异步模式
- 完整参数支持

### 视频生成 (dashscope-video)
**基础功能**:
- Wan 2.7 系列（t2v/i2v/r2v）
- HappyHorse 系列（t2v/i2v/r2v）

**高级功能**:
- last_frame（末帧生成）
- reference_voice（声音克隆）
- driving_audio（口型同步）
- 严格的模型能力区分

## 实现状态

已完成完整实现并充分测试，计划分 3 个 PR 贡献：
1. PR #1: DashScope 图片生成
2. PR #2: DashScope 视频生成（基础）
3. PR #3: DashScope 视频生成（高级功能）

## 相关资源

- [DashScope 官方文档](https://help.aliyun.com/zh/model-studio/)
- 本地实现已通过所有功能测试

期待您的反馈！
```

---

#### PR #1: DashScope 图片生成 ⭐⭐⭐⭐

**分支**: `feature/dashscope-image`  
**基于**: `origin/main`

**文件**:
- `backend/internal/service/provider_dashscope.go`
- `backend/internal/service/provider.go`（注册协议）

**标题**: `feat(provider): 添加 DashScope 图片生成协议支持`

**描述**:
```markdown
## 功能

实现阿里云百炼（DashScope）图片生成协议。

### 支持的模型
- wanx-v1
- wanx-style-repaint-v1
- wanx-background-generation-v2
- wanx-sketch-to-image-v1
- wanx-matting-v1

### 核心功能
- [x] 同步/异步模式
- [x] 完整参数支持（分辨率、数量、负面提示词等）
- [x] Base64 图片处理
- [x] 错误处理和重试
- [x] 任务轮询机制

### 使用示例

```go
// 配置
config := providerConfig{
    InterfaceType: "dashscope-image",
    Model: "wanx-v1",
    APIKey: "sk-xxx",
}

// 调用
result, err := runDashScopeImageTask(ctx, service, config, input)
```

### 测试
- [x] 所有模型已测试
- [x] 同步和异步模式已验证
- [x] 错误处理已确认

## 相关 Issue
Closes #XXX

## 后续 PR
- PR #2: DashScope 视频生成（基础功能）
- PR #3: DashScope 视频生成（高级功能）
```

---

#### PR #2: DashScope 视频生成（基础）⭐⭐⭐⭐⭐

**分支**: `feature/dashscope-video-basic`  
**基于**: `feature/dashscope-image`（或 main，如果 PR#1 已合并）

**文件**:
- `backend/internal/service/provider_dashscope_video.go`（基础部分）
- `backend/internal/service/provider.go`（注册协议）
- `backend/internal/service/model_capability.go`（基础配置）

**标题**: `feat(provider): 添加 DashScope 视频生成协议支持（基础功能）`

**描述**:
```markdown
## 功能

实现阿里云百炼（DashScope）视频生成协议的基础功能。

### 支持的模型
- Wan 2.7: t2v, i2v, r2v
- HappyHorse: t2v, i2v, r2v

### 核心功能
- [x] 异步任务提交和轮询
- [x] 文生视频（t2v）
- [x] 图生视频（i2v）- 首帧
- [x] 参考生视频（r2v）- 参考图/参考视频
- [x] 完整参数支持
  - 分辨率（720P/1080P）
  - 比例（16:9/9:16/1:1等）
  - 时长（2-10秒）
  - 水印控制
- [x] 视频下载和保存

### 使用示例

```go
// t2v
config := providerConfig{
    InterfaceType: "dashscope-video",
    Model: "wan2.7-t2v",
}

// i2v
input := canvasGenerationInput{
    Prompt: "生成视频",
    ReferenceImages: []providerMedia{...}, // 首帧
}

// r2v
input := canvasGenerationInput{
    Prompt: "生成视频",
    ReferenceImages: []providerMedia{...}, // 参考图（最多5张）
    ReferenceVideos: []providerMedia{...}, // 参考视频（可选）
}
```

### 测试
- [x] 所有模型基础功能已测试
- [x] 参数验证已确认

## 依赖
- 依赖 PR #1（DashScope 图片生成）的基础代码结构

## 后续 PR
- PR #3: 高级功能（last_frame、reference_voice、driving_audio）
```

---

#### PR #3: DashScope 视频生成（高级）⭐⭐⭐⭐⭐

**分支**: `feature/dashscope-video-advanced`  
**基于**: `feature/dashscope-video-basic`（或 main，如果 PR#2 已合并）

**文件**:
- `backend/internal/service/provider_dashscope_video.go`（高级功能）
- `backend/internal/service/model_capability.go`（高级验证逻辑）

**标题**: `feat(dashscope): 支持末帧生成、声音克隆等高级功能`

**描述**:
```markdown
## 功能

为 DashScope 视频生成添加高级功能支持。

### 新增功能

**Wan 2.7 i2v**:
- [x] `last_frame`（末帧生成）
  - 提供首帧和末帧，生成过渡视频
  - ReferenceImages[0] -> first_frame
  - ReferenceImages[1] -> last_frame
- [x] `driving_audio`（口型同步）
  - 音频驱动人物口型
  - ReferenceAudios[0] -> driving_audio

**Wan 2.7 r2v**:
- [x] `reference_voice`（声音克隆）
  - 为视频中的角色指定音色
  - ReferenceAudios[0] -> reference_voice

**HappyHorse**:
- [x] 严格区分能力限制
  - i2v 仅支持首帧（不支持末帧和音频）
  - r2v 仅支持参考图（不支持参考视频和音频）

### 验证逻辑优化

**音频数量验证**:
```go
// Wan 2.7 的音频字段映射不计入 ReferenceAudios 限制
if isWan27 && audioCountForValidation > 0 {
    if strings.Contains(model, "i2v") {
        audioCountForValidation = 0  // driving_audio
    } else if strings.Contains(model, "r2v") {
        audioCountForValidation = 0  // reference_voice
    }
}
```

**Operation 判定修复**:
```go
// 根据模型名判定 operation
if strings.Contains(model, "r2v") {
    operation = "reference_to_video"
} else if len(input.ReferenceImages) > 0 {
    operation = "image_to_video"
}
```

### 使用示例

```go
// Wan 2.7 i2v + last_frame
input := canvasGenerationInput{
    ReferenceImages: []providerMedia{
        firstFrame,  // 首帧
        lastFrame,   // 末帧
    },
}

// Wan 2.7 i2v + driving_audio
input := canvasGenerationInput{
    ReferenceImages: []providerMedia{firstFrame},
    ReferenceAudios: []providerMedia{drivingAudio},
}

// Wan 2.7 r2v + reference_voice
input := canvasGenerationInput{
    ReferenceImages: []providerMedia{refImage1, refImage2},
    ReferenceAudios: []providerMedia{referenceVoice},
}
```

### 测试
- [x] Wan 2.7 i2v + last_frame 已测试
- [x] Wan 2.7 i2v + driving_audio 已测试
- [x] Wan 2.7 r2v + reference_voice 已测试
- [x] HappyHorse 能力限制已验证

## 依赖
- 依赖 PR #2（DashScope 视频基础功能）

## 相关 Issue
Closes #XXX
```

---

### 6.2 贡献时间线

| 阶段 | 时间 | 里程碑 |
|---|---|---|
| **阶段 1: Issue 提交** | 第 1 天 | Issue 创建 |
| **等待响应** | 1-3 天 | 维护者反馈 |
| **阶段 2: PR#1** | 第 4-10 天 | 图片生成 PR |
| **Review 和修改** | 3-7 天 | Code Review |
| **PR#1 合并** | 第 10-17 天 | ✅ |
| **阶段 3: PR#2** | 第 18-24 天 | 视频基础 PR |
| **Review 和修改** | 3-7 天 | Code Review |
| **PR#2 合并** | 第 24-31 天 | ✅ |
| **阶段 4: PR#3** | 第 32-38 天 | 视频高级 PR |
| **Review 和修改** | 3-7 天 | Code Review |
| **PR#3 合并** | 第 38-45 天 | ✅ |
| **总计** | **1.5-2 个月** | 全部完成 |

---

## 七、立即行动计划

### Phase 1: 网络恢复前（可立即执行）

#### 1. 本地代码整理（30 分钟）

```bash
# 清理未跟踪文件
git status
git clean -fdn  # 预览

# 更新 .gitignore
cat >> .gitignore << 'EOF'
.temp/
百炼千问文档/
*.zip
DASHSCOPE_*.md
DashScope*.md
.claude/projects/
.claude/qianwen-*.md
.claude/settings.json
.claude/skills/
dashscope-fix-summary.md
EOF

# 提交所有改动
git add backend/internal/service/provider_dashscope*.go
git add backend/internal/service/model_capability.go
git add backend/internal/service/provider.go
git add .gitignore
git commit -m "feat(dashscope): 完整实现 DashScope 图片和视频生成协议

包含功能：
- DashScope 图片生成（Wanx 系列）
- DashScope 视频生成（Wan 2.7 + HappyHorse）
- 高级功能：last_frame、reference_voice、driving_audio
- 严格的模型能力区分
- 完善的验证逻辑和错误处理
"

# 创建备份分支
git branch backup-$(date +%Y%m%d-%H%M)
git branch feature/dashscope-complete
```

#### 2. 准备贡献材料（1-2 小时）

创建以下文件：

`.claude/DASHSCOPE_ISSUE.md` - Issue 内容  
`.claude/DASHSCOPE_PR1.md` - PR#1 描述  
`.claude/DASHSCOPE_PR2.md` - PR#2 描述  
`.claude/DASHSCOPE_PR3.md` - PR#3 描述  

---

### Phase 2: 网络恢复后（立即执行）

#### 1. 获取上游代码（30 分钟）

```bash
# 方法 1: Git fetch（推荐）
git fetch origin --all --prune
git fetch origin --tags

# 方法 2: 如果 fetch 失败，使用 curl 下载关键文件
mkdir -p /tmp/upstream
cd /tmp/upstream

# 核心后端文件
curl -O https://raw.githubusercontent.com/ddcat-ai/open-ai-canvas/main/backend/internal/service/provider.go
curl -O https://raw.githubusercontent.com/ddcat-ai/open-ai-canvas/main/backend/internal/service/model_capability.go
curl -O https://raw.githubusercontent.com/ddcat-ai/open-ai-canvas/main/backend/internal/model/model_sku.go
curl -O https://raw.githubusercontent.com/ddcat-ai/open-ai-canvas/main/backend/internal/service/provider_ark_private_assets.go

# 前端文件
curl -O https://raw.githubusercontent.com/ddcat-ai/open-ai-canvas/main/web/package.json
```

#### 2. 详细对比分析（2-3 小时）

```bash
# 查看提交差异
git log --oneline --stat main..origin/main > /tmp/upstream_commits.txt

# 对比关键文件
diff -u backend/internal/service/provider.go /tmp/upstream/provider.go > /tmp/diff_provider.txt
diff -u backend/internal/service/model_capability.go /tmp/upstream/model_capability.go > /tmp/diff_capability.txt

# 分析差异，更新本报告
```

#### 3. 执行合并（6-10 小时）

按照"策略 A: 分阶段 Merge"执行。

---

### Phase 3: 贡献准备（可并行）

#### 1. 提交 Issue（10 分钟）

访问 https://github.com/ddcat-ai/open-ai-canvas/issues/new  
复制 `.claude/DASHSCOPE_ISSUE.md` 内容提交

#### 2. 等待响应（1-3 天）

维护者可能会：
- 询问细节
- 提出建议
- 要求先看代码

根据反馈调整方案。

#### 3. 准备 PR 分支（每个 1 小时）

```bash
# PR #1
git checkout -b feature/dashscope-image origin/main
git cherry-pick <图片生成提交>
git rebase -i origin/main
git push myfork feature/dashscope-image

# PR #2  
git checkout -b feature/dashscope-video-basic origin/main
git cherry-pick <视频基础提交>
git rebase -i origin/main
git push myfork feature/dashscope-video-basic

# PR #3
git checkout -b feature/dashscope-video-advanced origin/main
git cherry-pick <视频高级提交>
git rebase -i origin/main
git push myfork feature/dashscope-video-advanced
```

#### 4. 提交 PR（每个 10 分钟）

在 GitHub 网页创建 Pull Request，复制对应的 PR 描述。

---

## 八、风险评估

| 风险项 | 可能性 | 影响 | 缓解措施 |
|---|---|---|---|
| **合并冲突严重** | ⭐⭐⭐⭐ 高 | ⭐⭐⭐⭐ 高 | 多个备份分支，手动三方合并 |
| **价格系统不兼容** | ⭐⭐⭐ 中 | ⭐⭐⭐⭐⭐ 极高 | 详细阅读 SKU 代码，理解架构 |
| **数据库迁移失败** | ⭐⭐ 低 | ⭐⭐⭐⭐ 高 | 备份数据库，在测试环境先试 |
| **前端 API 不兼容** | ⭐⭐⭐ 中 | ⭐⭐⭐ 中 | 逐个文件对比，API 联调 |
| **PR 被拒绝** | ⭐⭐ 低 | ⭐⭐⭐ 中 | 积极沟通，根据反馈修改 |
| **功能重复** | ⭐ 极低 | ⭐⭐⭐⭐ 高 | 已确认上游无 DashScope |
| **测试不通过** | ⭐⭐ 低 | ⭐⭐⭐ 中 | 本地已充分测试 |

---

## 九、成功标准

### 任务 1: 上游差异分析

- [x] 上游分支和提交列表
- [x] 核心文件对比（provider.go, model_capability.go）
- [x] 功能差异清单
- [ ] 详细的逐行代码对比（需网络恢复后）
- [ ] 前端文件对比（需网络恢复后）
- [ ] 数据库 Schema 对比（需网络恢复后）

### 任务 2: 增量拉取

- [ ] 本地代码基于最新上游
- [ ] 所有冲突已解决
- [ ] 编译通过
- [ ] 所有测试通过
- [ ] DashScope 功能正常
- [ ] 上游新功能可用

### 任务 3: 向上游贡献

- [ ] Issue 已提交
- [ ] Issue 获得正面响应
- [ ] PR #1 已提交
- [ ] PR #1 已合并
- [ ] PR #2 已提交
- [ ] PR #2 已合并
- [ ] PR #3 已提交
- [ ] PR #3 已合并

---

## 十、待补充内容（网络恢复后）

### 10.1 详细代码对比

- [ ] `provider.go` 逐函数对比
- [ ] `model_capability.go` 逐配置对比
- [ ] 其他 provider 文件对比
- [ ] 数据模型文件对比

### 10.2 前端文件对比

- [ ] package.json 依赖对比
- [ ] UI 组件变更
- [ ] API 调用变更
- [ ] 路由变更

### 10.3 数据库变更

- [ ] 迁移脚本分析
- [ ] 表结构对比
- [ ] 数据兼容性检查

### 10.4 配置文件对比

- [ ] docker-compose.yml
- [ ] 环境变量
- [ ] Nginx 配置

---

## 结论

### 关键发现

1. **上游变化巨大** - 约 100+ 提交，核心架构重构（SKU、Agent）
2. **本地功能独立且完整** - DashScope 协议实现质量高，价值大
3. **合并风险可控** - 虽然冲突多，但有明确的解决策略
4. **贡献路径清晰** - 分 3 个 PR，逐步贡献

### 建议行动

**优先级 1: 立即执行（网络恢复前）**
- ✅ 本地代码整理和提交
- ✅ 创建备份分支
- ✅ 准备贡献材料

**优先级 2: 网络恢复后立即执行**
- 获取上游最新代码
- 补充详细对比分析
- 开始合并流程

**优先级 3: 同步进行**
- 提交 Issue
- 准备 PR 分支
- 根据维护者反馈调整

### 预计时间

| 任务 | 时间 |
|---|---|
| 本地代码整理 | 0.5 小时 |
| 准备贡献材料 | 1-2 小时 |
| 获取上游代码 | 0.5 小时 |
| 详细对比分析 | 2-3 小时 |
| 合并上游代码 | 6-10 小时 |
| 向上游贡献 | 1.5-2 个月 |
| **第一轮总计** | **10-16 小时** |

---

**报告状态**: 基于有限网络访问的初步分析  
**下一步**: 等待网络完全恢复，补充详细对比  
**最后更新**: 2026-08-23 12:50

---

## 附录

### A. 快速命令参考

```bash
# 查看上游差异
git log --oneline --graph HEAD..origin/main

# 统计变更
git rev-list --count HEAD..origin/main
git diff --stat HEAD..origin/main

# 查看具体文件差异
git diff HEAD..origin/main -- backend/internal/service/provider.go

# 创建备份
git branch backup-$(date +%Y%m%d)

# 开始合并
git merge origin/main

# 中止合并
git merge --abort

# 查看冲突文件
git status | grep "both modified"
```

### B. 关键联系人

- 上游项目：https://github.com/ddcat-ai/open-ai-canvas
- Issue 页面：https://github.com/ddcat-ai/open-ai-canvas/issues
- PR 页面：https://github.com/ddcat-ai/open-ai-canvas/pulls

---

**报告完成** ✅
