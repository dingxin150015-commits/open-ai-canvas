---
name: protocol-registration-checklist
description: 新增协议时必须同步更新的 6 个位置
metadata:
  type: reference
---

每次新增协议时，必须在以下 6 个位置注册，缺一不可。

## 注册位置

### 后端（Go）- 5 个位置

#### 1. 常量定义

**文件**：`backend/internal/model/models.go`

**位置**：`ChannelInterfaceType` 常量声明区域

**格式**：
```go
ChannelInterfaceXXXYyy ChannelInterfaceType = "xxx-yyy"
```

**命名规则**：
- 常量名：大驼峰（ChannelInterfaceXXXYyy）
- 协议值：kebab-case（"xxx-yyy"）
- 保持字母顺序或按功能分组

**示例**：
```go
const (
    // ... 其他协议
    ChannelInterfaceDashScopeImage        ChannelInterfaceType = "dashscope-image"
    ChannelInterfaceDashScopeVideo        ChannelInterfaceType = "dashscope-video"
    ChannelInterfaceGeminiImage           ChannelInterfaceType = "gemini-image"
    ChannelInterfaceMiniMaxVideo          ChannelInterfaceType = "minimax-video"
    // ...
)
```

---

#### 2. allowed 表注册

**文件**：`backend/internal/service/provider.go`

**位置**：`allowed` map 定义（约 2800 行）

**格式**：
```go
allowed := map[string]map[string]bool{
    "image": {"xxx-image": true, ...},
    "video": {"xxx-video": true, ...},
    "audio": {"xxx-audio": true, ...},
    "text":  {"chat-completion": true, ...},
}
```

**作用**：
- 验证协议是否允许用于该能力类型
- 防止错误配置（如将视频协议用于图片生成）

**示例**：
```go
allowed := map[string]map[string]bool{
    "image": {
        "openai-image": true, 
        "dashscope-image": true,  // 新增
        "gemini-image": true,     // 新增
    },
    "video": {
        "newapi": true, 
        "dashscope-video": true,  // 新增
        "minimax-video": true,    // 新增
    },
}
```

---

#### 3. 路由分支

**文件**：`backend/internal/service/provider.go`

**位置**：`runImageTask` 或 `runVideoTask` 等函数内

**格式**：
```go
if input.Config.InterfaceType == string(model.ChannelInterfaceXXX) {
    return runXXXTask(ctx, input)
}
```

**作用**：
- 根据协议类型路由到对应的处理函数
- 实际执行协议逻辑的入口

**示例**：
```go
func runImageTask(ctx context.Context, input canvasGenerationInput) (map[string]interface{}, error) {
    if input.Config.InterfaceType == string(model.ChannelInterfaceGeminiImage) {
        return runGeminiImageTask(ctx, input)
    }
    if input.Config.InterfaceType == string(model.ChannelInterfaceDashScopeImage) {
        return runDashScopeImageTask(ctx, input)  // 新增
    }
    // ... 其他协议
}

func runVideoTask(ctx context.Context, input canvasGenerationInput) (map[string]interface{}, error) {
    if input.Config.InterfaceType == string(model.ChannelInterfaceMiniMaxVideo) {
        return runMiniMaxVideoTask(ctx, input)
    }
    if input.Config.InterfaceType == string(model.ChannelInterfaceDashScopeVideo) {
        return runDashScopeVideoTask(ctx, input)  // 新增
    }
    // ... 其他协议
}
```

---

#### 4. 协议验证

**文件**：`backend/internal/service/admin.go`

**位置**：`validChannelInterfaceType` 函数（约 825 行）

**格式**：
```go
case model.ChannelInterfaceXXX, ..., model.ChannelInterfaceYYY:
    return true
```

**作用**：
- 管理后台保存渠道时验证协议名称合法性
- 防止无效协议进入数据库

**示例**：
```go
func validChannelInterfaceType(value model.ChannelInterfaceType) bool {
    switch value {
    case model.ChannelInterfaceChatCompletion,
         model.ChannelInterfaceOpenAIImage,
         model.ChannelInterfaceDashScopeImage,   // 新增
         model.ChannelInterfaceDashScopeVideo,   // 新增
         model.ChannelInterfaceGeminiImage,      // 新增
         model.ChannelInterfaceMiniMaxVideo:     // 新增
        return true
    default:
        return false
    }
}
```

**注意**：
- 所有协议必须加入同一个 `case` 语句
- 不要创建多个 case（会降低可读性）

---

#### 5. 能力映射

**文件**：`backend/internal/service/channel_models.go`

**位置**：`capabilityForProtocol` 函数（约 852 行）

**格式**：
```go
case model.ChannelInterfaceXXX:
    return "image"  // 或 "video", "audio", "text"
```

**作用**：
- 将协议名映射到能力类型
- 用于模型配置、路由选择

**示例**：
```go
func capabilityForProtocol(protocol model.ChannelInterfaceType) string {
    switch protocol {
    case model.ChannelInterfaceOpenAIImage,
         model.ChannelInterfaceGrokImage,
         model.ChannelInterfaceDashScopeImage,   // 新增
         model.ChannelInterfaceGeminiImage:      // 新增
        return "image"
    
    case model.ChannelInterfaceNewAPIVideo,
         model.ChannelInterfaceDashScopeVideo,   // 新增
         model.ChannelInterfaceMiniMaxVideo:     // 新增
        return "video"
    
    case model.ChannelInterfaceOpenAIAudio:
        return "audio"
    
    case model.ChannelInterfaceChatCompletion:
        return "text"
    
    default:
        return ""
    }
}
```

---

### 前端（TypeScript）- 1 个位置

#### 6. 类型定义和配置

**文件**：`web/src/lib/model-protocols.ts`

**位置 A**：`ModelProtocol` 类型联合（约第 7 行）

**格式**：
```typescript
export type ModelProtocol =
    | "xxx-yyy"
    | ...
```

**位置 B**：`MODEL_PROTOCOLS` 配置数组（约第 30 行）

**格式**：
```typescript
export const MODEL_PROTOCOLS = [
    { 
        value: "xxx-yyy", 
        label: "XXX 显示名", 
        capability: "image",  // 或 "video", "audio"
        create: "POST /api/endpoint",
        contentType: "application/json",
        media: "参考素材说明"
    },
    ...
]
```

**作用**：
- 前端下拉框选项
- 协议配置信息展示
- TypeScript 类型安全

**示例**：
```typescript
export type ModelProtocol =
    | "chat-completion"
    | "openai-image"
    | "dashscope-image"    // 新增
    | "dashscope-video"    // 新增
    | "gemini-image"       // 新增
    | "minimax-video";     // 新增

export const MODEL_PROTOCOLS = [
    { 
        value: "dashscope-image", 
        label: "DashScope 图片", 
        capability: "image", 
        create: "POST /api/v1/services/aigc/multimodal-generation/generation", 
        contentType: "application/json (同步模式)", 
        media: "文生图与参考图，content 多模态数组格式" 
    },
    { 
        value: "dashscope-video", 
        label: "DashScope 视频", 
        capability: "video", 
        create: "POST /services/aigc/video-generation/video-synthesis", 
        poll: "GET /tasks/{task_id}", 
        contentType: "application/json (异步模式)", 
        media: "文生视频、图生视频、参考生视频，media 多模态数组" 
    },
    // ...
]
```

---

## 验证方法

### 自动验证脚本

```bash
#!/bin/bash
# 验证协议注册完整性

PROTOCOL_NAME="dashscope-image"
PROTOCOL_CONSTANT="DashScopeImage"

echo "=== 协议注册完整性检查 ==="
echo "协议名: $PROTOCOL_NAME"
echo

# 1. 常量定义
echo "1. 常量定义 (models.go):"
grep "$PROTOCOL_CONSTANT" backend/internal/model/models.go && echo "  ✅" || echo "  ❌ 缺失"

# 2. allowed 表
echo "2. allowed 表 (provider.go):"
grep "\"$PROTOCOL_NAME\"" backend/internal/service/provider.go | grep -q "true" && echo "  ✅" || echo "  ❌ 缺失"

# 3. 路由分支
echo "3. 路由分支 (provider.go):"
grep -q "ChannelInterface$PROTOCOL_CONSTANT" backend/internal/service/provider.go && echo "  ✅" || echo "  ❌ 缺失"

# 4. 协议验证
echo "4. 协议验证 (admin.go):"
grep -q "$PROTOCOL_CONSTANT" backend/internal/service/admin.go && echo "  ✅" || echo "  ❌ 缺失"

# 5. 能力映射
echo "5. 能力映射 (channel_models.go):"
grep -q "$PROTOCOL_CONSTANT" backend/internal/service/channel_models.go && echo "  ✅" || echo "  ❌ 缺失"

# 6. TS 类型
echo "6. TS 类型 (model-protocols.ts):"
grep -q "\"$PROTOCOL_NAME\"" web/src/lib/model-protocols.ts && echo "  ✅" || echo "  ❌ 缺失"

echo
echo "=== 编译验证 ==="
cd backend && go build ./cmd/server && echo "Go 编译: ✅" || echo "Go 编译: ❌"
cd ../web && npm run typecheck && echo "TS 类型: ✅" || echo "TS 类型: ❌"
```

---

### 手动验证清单

新增协议后，逐项检查：

- [ ] **models.go** - 常量定义存在
- [ ] **provider.go** - allowed 表包含协议
- [ ] **provider.go** - 路由分支调用正确函数
- [ ] **admin.go** - validChannelInterfaceType 包含协议
- [ ] **channel_models.go** - capabilityForProtocol 映射正确
- [ ] **model-protocols.ts** - 类型联合包含协议
- [ ] **model-protocols.ts** - 配置数组包含条目
- [ ] **Go 编译通过** - `go build ./cmd/server`
- [ ] **TS 类型通过** - `npm run typecheck`
- [ ] **实现文件存在** - `provider_xxx.go`

---

## 常见错误

### 错误 1: 遗漏 admin.go 验证

**现象**：
- 编译通过
- 运行时管理后台保存渠道报错"非法协议"

**原因**：
- `validChannelInterfaceType` 未包含新协议
- 保存时被拦截

**修复**：
```go
// admin.go:825
case ..., model.ChannelInterfaceNewProtocol:  // 添加这里
    return true
```

---

### 错误 2: 遗漏 channel_models.go 映射

**现象**：
- 编译通过
- 模型配置时无法识别协议能力

**原因**：
- `capabilityForProtocol` 返回空字符串
- 无法解析协议属于哪种能力

**修复**：
```go
// channel_models.go:852
case ..., model.ChannelInterfaceNewProtocol:  // 添加这里
    return "image"  // 或 "video" 等
```

---

### 错误 3: TS 类型定义不完整

**现象**：
- 前端下拉框有选项
- TypeScript 编译报错："Type 'xxx' is not assignable"

**原因**：
- 只添加了 `MODEL_PROTOCOLS` 配置
- 忘记添加 `ModelProtocol` 类型联合

**修复**：
```typescript
// model-protocols.ts
export type ModelProtocol =
    | "existing-protocol"
    | "new-protocol"  // 添加这里
    | ...
```

---

### 错误 4: allowed 表分类错误

**现象**：
- 编译通过
- 运行时报错"协议不支持该能力类型"

**原因**：
- 将 image 协议加到 video 分类
- 或将 video 协议加到 image 分类

**修复**：
```go
allowed := map[string]map[string]bool{
    "image": {"image-protocol": true},  // 确保匹配
    "video": {"video-protocol": true},  // 确保匹配
}
```

---

## 新增协议模板

### 步骤 1: 创建实现文件

```go
// backend/internal/service/provider_xxx.go
package service

func runXXXTask(ctx context.Context, input canvasGenerationInput) (map[string]interface{}, error) {
    // 实现协议逻辑
    return result, nil
}
```

---

### 步骤 2: 注册到 6 个位置

使用本文档的"注册位置"章节作为清单，逐个完成。

---

### 步骤 3: 验证

```bash
# 1. 编译验证
cd backend && go build ./cmd/server
cd ../web && npm run typecheck

# 2. 运行验证脚本
./verify_protocol.sh

# 3. 功能测试
# - 管理后台创建渠道
# - 配置模型使用新协议
# - 测试生成功能
```

---

## 协议命名规范

### 后端（Go）

**常量名**：大驼峰
```go
ChannelInterfaceDashScopeImage  // ✅
ChannelInterfacedashscopeimage  // ❌ 小写
ChannelInterface_DashScope      // ❌ 下划线
```

**协议值**：kebab-case
```go
"dashscope-image"  // ✅
"DashScopeImage"   // ❌ 大驼峰
"dashscope_image"  // ❌ 下划线
```

---

### 前端（TypeScript）

**类型值**：与后端协议值一致
```typescript
"dashscope-image"  // ✅ 与后端一致
"dashScopeImage"   // ❌ 驼峰
```

**显示名**：用户友好的中文或英文
```typescript
label: "DashScope 图片"     // ✅ 中文
label: "DashScope Images"   // ✅ 英文
label: "dashscope-image"    // ❌ 技术名称
```

---

## 相关链接

- [[large-scale-merge-best-practices]] - 大规模合并时的检查清单
- [[git-conflict-resolution-strategy]] - 冲突解决策略
- [[api-development-best-practices]] - API 开发最佳实践

---

**Why**: 协议定义分散在多个文件，遗漏任何一处会导致功能失效。系统性检查是唯一可靠的方法。

**How to apply**: 
- 每次新增协议时，打印本文档作为检查清单
- 逐个位置验证，完成后打勾
- 使用验证脚本自动检查
- 大规模合并后，用此清单验证所有协议完整性
