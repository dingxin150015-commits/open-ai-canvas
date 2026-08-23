# 千问/百炼模型集成总结报告

## 完成状态

✅ **已完成：系统性整理所有模型的 API 格式**  
✅ **已完成：创建技术文档** → [`.claude/qianwen-models-api-reference.md`](qianwen-models-api-reference.md)  
⏳ **待完成：沉淀千问技能**（视需要而定）

---

## 核心发现：DashScope 图片生成的关键问题

### 当前错误根因

**错误信息：**
```json
{
  "code": "InvalidParameter",
  "message": "Input should be a valid list: input.messages.0.content"
}
```

**问题诊断：**

当前代码（[`provider_dashscope.go:22-36`](../backend/internal/service/provider_dashscope.go:22-36)）发送的格式：

```go
messages = append(messages, map[string]interface{}{
    "role":    "user",
    "content": input.Prompt,  // ❌ 字符串
})
```

实际发送给 API：
```json
{
  "messages": [{
    "role": "user",
    "content": "一只小猫"  // ❌ 字符串
  }]
}
```

**DashScope API 期望的格式：**

```json
{
  "messages": [{
    "role": "user",
    "content": [              // ✅ 必须是数组
      {"text": "一只小猫"}   // ✅ 对象包装
    ]
  }]
}
```

### 修复方案

**修改位置：** `backend/internal/service/provider_dashscope.go` 函数 `runDashScopeImageTask`

**修改要点：**
1. `content` 从字符串改为对象数组
2. 文本包装在 `{"text": "..."}` 对象中
3. 参考图包装在 `{"image": "data:image/..."}` 对象中

**完整修复代码见：** [`.claude/qianwen-models-api-reference.md` 第三章](qianwen-models-api-reference.md#31-当前-provider_dashscopego-的问题)

---

## 模型 API 格式对比表

### 图片生成模型

| 模型 | 调用方式 | 最大分辨率 | content格式 | 特殊参数 |
|------|---------|-----------|------------|---------|
| `wan2.7-image-pro` | 同步+异步 | 4096×4096 | **对象数组** | `thinking_mode`, `enable_sequential`, `color_palette` |
| `wan2.7-image` | 同步+异步 | 2048×2048 | **对象数组** | 同上 |
| `qwen-image-3.0-pro` | 仅同步 | 2048×2048 | **对象数组** | `prompt_extend`, `prompt_extend_mode` |
| `qwen-image-3.0` | 仅同步 | 2048×2048 | **对象数组** | `prompt_extend` |
| `qwen-image-2.0-pro` | 仅同步 | 2048×2048 | **对象数组** | - |
| `qwen-image-2.0` | 仅同步 | 2048×2048 | **对象数组** | - |

**关键统一点：**
- ✅ 所有模型的 `content` 都必须是 **对象数组**
- ✅ 文本使用 `{"text": "..."}` 格式
- ✅ 图片使用 `{"image": "data:image/..."}` 格式

### 视频生成模型

| 模型 | 类型 | 调用方式 | 状态 |
|------|------|---------|------|
| `happyhorse-1.1-t2v` | 文生视频 | **仅异步** | 📝 格式待查 |
| `happyhorse-1.1-i2v` | 图生视频 | **仅异步** | 📝 格式待查 |
| `happyhorse-1.1-r2v` | 参考生视频 | **仅异步** | 📝 格式待查 |
| `wan2.7-t2v-2026-06-12` | 文生视频 | **仅异步** | 📝 格式待查 |
| `wan2.7-i2v-2026-04-25` | 图生视频 | **仅异步** | 📝 格式待查 |
| `wan2.7-r2v-2026-06-12` | 参考生视频 | **仅异步** | 📝 格式待查 |

**注意：** 
- 视频模型**必须**设置 `X-DashScope-Async: enable` header
- 需要轮询 `/api/v1/tasks/{task_id}` 获取结果
- 推测使用相同的 messages 格式（对象数组）

---

## 三协议方案架构

项目需要支持的三种协议：

### 1. OpenAI 兼容协议（文本对话）

**用途：** 文本生成、对话  
**端点：** `https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions`  
**格式：** 标准 OpenAI messages，`content` 可以是字符串或数组  
**实现：** 已有（通过 OpenAI 插件）

### 2. DashScope 图片协议

**用途：** 图片生成与编辑  
**端点：** `https://dashscope.aliyuncs.com/api/v1/services/aigc/multimodal-generation/generation`  
**格式：** DashScope messages，`content` **必须是对象数组**  
**实现：** 部分完成，需修复 content 格式

### 3. Gemini 协议

**用途：** 多模态生成  
**端点：** Google Gemini API  
**格式：** 使用 `parts` 数组  
**实现：** 通过 OpenAI 插件的 Gemini 格式支持

---

## 实施优先级

### 🔴 优先级1：修复图片生成（立即执行）

**任务：**
1. 修改 `provider_dashscope.go` 的 content 构建逻辑
2. 测试 `qwen-image-3.0-pro` 同步调用
3. 测试 `wan2.7-image-pro` 异步调用

**验证步骤：**
1. Docker 重新构建 backend
2. 管理后台配置 DashScope 渠道
3. 选择模型并点击"测试模型"
4. 确认返回图片URL而不是错误

### 🟡 优先级2：实现视频生成（功能扩展）

**任务：**
1. 访问视频生成的官方文档，确认确切格式
2. 创建 `runDashScopeVideoTask` 函数
3. 实现异步轮询机制
4. 添加 `ChannelInterfaceDashScopeVideo` 协议

**前置条件：**
- 图片生成已验证成功

### 🟢 优先级3：完善模型发现（优化体验）

**任务：**
1. 在 `bailian/discovery.go` 中补充视频模型定义
2. 前端添加视频协议选项
3. 完善错误处理和用户提示

---

## 已创建的文档

1. **API 格式参考文档**  
   文件：`.claude/qianwen-models-api-reference.md`  
   内容：完整的请求/响应格式、参数说明、错误诊断

2. **本总结报告**  
   文件：`.claude/QIANWEN_INTEGRATION_SUMMARY.md`  
   内容：核心发现、对比表格、实施建议

---

## 待确认事项

### 需要访问官方文档

- [ ] 视频生成的确切端点路径
- [ ] 视频模型的参数列表（duration, resolution, fps等）
- [ ] 参考生视频（r2v）的具体输入格式
- [ ] HappyHorse 和 Wan 视频模型的差异

### 需要实际测试

- [ ] 图片生成修复后是否成功
- [ ] 4K分辨率是否真的可用（wan2.7-image-pro）
- [ ] 异步轮询的合理间隔
- [ ] 图片/视频URL的实际有效期

---

## 下一步行动建议

**立即执行（你的决策）：**

**选项A：** 立即修复图片生成
- 修改 `provider_dashscope.go`
- 测试验证
- 提交修复

**选项B：** 先获取视频文档，再统一实现
- 继续访问官方文档获取视频API格式
- 一次性实现图片+视频
- 提交完整功能

**选项C：** 分阶段实现
- 本次只修复图片生成并验证
- 下次单独实现视频生成

---

**报告生成时间：** 2026-01-20  
**文档状态：** ✅ 已完成系统性整理  
**技能沉淀：** 可根据需要创建 Claude 技能文件
