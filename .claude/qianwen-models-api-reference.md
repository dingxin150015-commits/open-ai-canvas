# 百炼/千问模型 API 格式参考文档

## 文档说明

本文档整理了千问AI平台（百炼/DashScope）的图片生成和视频生成模型的API调用格式，为项目 `open-ai-canvas` 正确集成这些模型提供参考。

**官方文档地址：** https://platform.qianwenai.com/docs

**创建时间：** 2026-01-20  
**目标模型：**
- 图片生成：`wan2.7-image-pro`, `qwen-image-3.0-pro`, `qwen-image-2.0-pro`, `wan2.7-image`, `qwen-image-3.0`, `qwen-image-2.0`
- 视频生成：`happyhorse-1.1-t2v`, `wan2.7-t2v-2026-06-12` (文生视频)
- 视频生成：`happyhorse-1.1-i2v`, `wan2.7-i2v-2026-04-25` (图生视频)
- 视频生成：`happyhorse-1.1-r2v`, `wan2.7-r2v-2026-06-12` (参考生视频)

---

## 核心要点总结

### 1. 统一端点架构

千问AI平台提供两种API调用方式：

#### **方式A：OpenAI兼容模式（推荐用于文本对话）**
- Base URL: `https://dashscope.aliyuncs.com/compatible-mode/v1`
- 完全兼容 OpenAI SDK
- 适用于文本生成、对话等场景

#### **方式B：DashScope原生模式（图片/视频生成）**
- Base URL: `https://dashscope.aliyuncs.com`
- 使用 DashScope 特有的消息格式
- 必须用于多模态生成（图片、视频）

### 2. 认证方式

所有请求统一使用 **Bearer Token** 认证：

```http
Authorization: Bearer $DASHSCOPE_API_KEY
```

API Key 格式：`sk-ws-xxxxxxxx`

### 3. 同步 vs 异步

| 模型类型 | 调用方式 | 标识方法 |
|---------|---------|---------|
| 图片生成 | 同步或异步 | 异步需设置 `X-DashScope-Async: enable` Header |
| 视频生成 | **仅异步** | 必须设置 `X-DashScope-Async: enable` |

**异步任务轮询端点：**
```
GET https://dashscope.aliyuncs.com/api/v1/tasks/{task_id}
Authorization: Bearer $DASHSCOPE_API_KEY
```

---

## 一、图片生成模型

### 1.1 模型列表

| 模型 ID | 名称 | 调用方式 | 最大分辨率 | 特点 |
|---------|------|----------|-----------|------|
| `wan2.7-image-pro` | 万相2.7 Pro | 同步+异步 | 4096×4096 | 组图、超高分辨率、色彩控制 |
| `wan2.7-image` | 万相2.7 | 同步+异步 | 2048×2048 | 平衡版本 |
| `qwen-image-3.0-pro` | 千问图像3.0 Pro | **仅同步** | 2048×2048 | 文字渲染、智能改写 |
| `qwen-image-3.0` | 千问图像3.0 | **仅同步** | 2048×2048 | 平衡版本 |
| `qwen-image-2.0-pro` | 千问图像2.0 Pro | **仅同步** | 2048×2048 | 旧版本 |
| `qwen-image-2.0` | 千问图像2.0 | **仅同步** | 2048×2048 | 旧版本 |

### 1.2 请求格式（同步调用）

**端点：**
```
POST https://dashscope.aliyuncs.com/api/v1/services/aigc/multimodal-generation/generation
```

**请求体格式（关键变更）：**

```json
{
  "model": "qwen-image-3.0-pro",
  "input": {
    "messages": [
      {
        "role": "user",
        "content": [
          {
            "text": "一只可爱的橘色小猫"
          }
        ]
      }
    ]
  },
  "parameters": {
    "size": "1024*1024",
    "prompt_extend": true
  }
}
```

**⚠️ 关键要点：**
1. ✅ **`content` 必须是数组**，不能是字符串
2. ✅ 文本内容包装在 `{"text": "..."}` 对象中
3. ✅ 参考图包装在 `{"image": "data:image/..."}` 对象中
4. ✅ System prompt 可以合并到 user 的 text 中，或作为单独的 system message

**完整示例（带参考图）：**

```json
{
  "model": "wan2.7-image-pro",
  "input": {
    "messages": [
      {
        "role": "system",
        "content": [
          {
            "text": "你是专业的图像生成助手"
          }
        ]
      },
      {
        "role": "user",
        "content": [
          {
            "text": "根据这张图片生成相似风格的作品"
          },
          {
            "image": "data:image/png;base64,iVBORw0KGgo..."
          }
        ]
      }
    ]
  },
  "parameters": {
    "size": "2048*2048",
    "thinking_mode": true,
    "enable_sequential": false
  }
}
```

### 1.3 响应格式（同步）

```json
{
  "output": {
    "choices": [
      {
        "finish_reason": "stop",
        "message": {
          "role": "assistant",
          "content": [
            {
              "image": "https://dashscope-result-sz.oss-cn-shenzhen.aliyuncs.com/xxx.png?Expires=xxx"
            }
          ]
        }
      }
    ]
  },
  "usage": {
    "output_height": 1024,
    "output_width": 1024,
    "input_image_count": 0,
    "output_image_count": 1
  },
  "request_id": "571ae02f-5c9d-436c-83c2-f221e6df0xxx"
}
```

**⚠️ 图片URL有效期：24小时**  
必须立即下载并转换为 data URL 或保存到 OSS。

### 1.4 异步调用（万相模型）

**步骤1：创建任务**

```http
POST https://dashscope.aliyuncs.com/api/v1/services/aigc/image-generation/generation
Authorization: Bearer $DASHSCOPE_API_KEY
X-DashScope-Async: enable
Content-Type: application/json

{
  "model": "wan2.7-image-pro",
  "input": {
    "messages": [
      {
        "role": "user",
        "content": [
          {
            "text": "一位年轻女性，自然随性的自拍风格"
          }
        ]
      }
    ]
  },
  "parameters": {
    "size": "2K",
    "n": 1
  }
}
```

**响应（立即返回）：**

```json
{
  "output": {
    "task_id": "77093787-a217-4c29-9cd4-ca7b5ac86xxx",
    "task_status": "PENDING"
  },
  "request_id": "4fb3050f-de57-4a24-84ff-e37ee5xxxxxx"
}
```

**步骤2：轮询任务状态**

```http
GET https://dashscope.aliyuncs.com/api/v1/tasks/77093787-a217-4c29-9cd4-ca7b5ac86xxx
Authorization: Bearer $DASHSCOPE_API_KEY
```

**响应（成功）：**

```json
{
  "output": {
    "task_id": "77093787-a217-4c29-9cd4-ca7b5ac86xxx",
    "task_status": "SUCCEEDED",
    "submit_time": "2026-03-31 23:04:46.166",
    "end_time": "2026-03-31 23:05:11.664",
    "choices": [
      {
        "finish_reason": "stop",
        "message": {
          "role": "assistant",
          "content": [
            {
              "image": "https://dashscope-result-bj.oss-cn-beijing.aliyuncs.com/xxxxxx.png?Expires=xxxxxx",
              "type": "image"
            }
          ]
        }
      }
    ]
  },
  "usage": {
    "image_count": 1,
    "size": "2048*2048"
  }
}
```

**任务状态枚举：**
- `PENDING` - 排队中
- `RUNNING` - 生成中
- `SUCCEEDED` - 成功
- `FAILED` - 失败

### 1.5 关键参数说明

| 参数 | 类型 | 说明 | 适用模型 |
|------|------|------|----------|
| `size` | string | 格式：`"宽*高"` 或 `"1K"`/`"2K"`/`"4K"` | 所有 |
| `prompt_extend` | boolean | 智能改写提示词（默认true，增加3-5秒） | qwen-image系列 |
| `thinking_mode` | boolean | 思考模式，提升质量（默认true） | wan系列 |
| `enable_sequential` | boolean | 组图生成（1-12张连贯图） | wan系列 |
| `color_palette` | array | 自定义颜色方案（3-10种颜色） | wan系列 |
| `n` | integer | 生成数量 | 所有 |

---

## 二、视频生成模型

### 2.1 模型列表

| 模型 ID | 类型 | 调用方式 | 最大时长 | 最大分辨率 |
|---------|------|----------|----------|-----------|
| `happyhorse-1.1-t2v` | 文生视频 | **仅异步** | 待查 | 待查 |
| `happyhorse-1.1-i2v` | 图生视频 | **仅异步** | 待查 | 待查 |
| `happyhorse-1.1-r2v` | 参考生视频 | **仅异步** | 待查 | 待查 |
| `wan2.7-t2v-2026-06-12` | 文生视频 | **仅异步** | 待查 | 待查 |
| `wan2.7-i2v-2026-04-25` | 图生视频 | **仅异步** | 待查 | 待查 |
| `wan2.7-r2v-2026-06-12` | 参考生视频 | **仅异步** | 待查 | 待查 |

### 2.2 请求格式（待补充）

**⚠️ 注意：视频生成模型必须使用异步调用**

**推测端点：**
```
POST https://dashscope.aliyuncs.com/api/v1/services/aigc/video-generation/generation
X-DashScope-Async: enable
```

**推测请求格式（文生视频）：**

```json
{
  "model": "wan2.7-t2v-2026-06-12",
  "input": {
    "messages": [
      {
        "role": "user",
        "content": [
          {
            "text": "一只小猫在草地上奔跑"
          }
        ]
      }
    ]
  },
  "parameters": {
    "duration": "5s",
    "resolution": "1280*720"
  }
}
```

**推测请求格式（图生视频）：**

```json
{
  "model": "wan2.7-i2v-2026-04-25",
  "input": {
    "messages": [
      {
        "role": "user",
        "content": [
          {
            "text": "让这张图片动起来"
          },
          {
            "image": "data:image/png;base64,..."
          }
        ]
      }
    ]
  },
  "parameters": {
    "duration": "5s"
  }
}
```

### 2.3 响应格式（待补充）

**推测异步任务响应：**

```json
{
  "output": {
    "task_id": "video-task-xxx",
    "task_status": "PENDING"
  }
}
```

**推测完成后的结果：**

```json
{
  "output": {
    "task_id": "video-task-xxx",
    "task_status": "SUCCEEDED",
    "choices": [
      {
        "message": {
          "role": "assistant",
          "content": [
            {
              "video": "https://dashscope-result-xxx.oss-cn-xxx.aliyuncs.com/xxx.mp4?Expires=xxx",
              "type": "video"
            }
          ]
        }
      }
    ]
  },
  "usage": {
    "video_duration": "5s",
    "resolution": "1280*720"
  }
}
```

---

## 三、与当前项目的对比

### 3.1 当前 provider_dashscope.go 的问题

**错误代码（第22-36行）：**

```go
messages := []map[string]interface{}{}

if systemPrompt := strings.TrimSpace(input.Config.SystemPrompt); systemPrompt != "" {
    messages = append(messages, map[string]interface{}{
        "role":    "system",
        "content": systemPrompt,  // ❌ 错误：content 是字符串
    })
}

messages = append(messages, map[string]interface{}{
    "role":    "user",
    "content": input.Prompt,  // ❌ 错误：content 是字符串
})
```

**正确代码：**

```go
messages := []map[string]interface{}{}

// 构建 content 数组
content := []map[string]interface{}{}

// 合并 system prompt 和 user prompt
promptText := input.Prompt
if systemPrompt := strings.TrimSpace(input.Config.SystemPrompt); systemPrompt != "" {
    promptText = systemPrompt + "\n\n" + promptText
}

content = append(content, map[string]interface{}{
    "text": promptText,  // ✅ 正确：包装在 text 对象中
})

// 处理参考图
if len(input.ReferenceImages) > 0 {
    raw, _, err := mediaBytes(input.ReferenceImages[0])
    if err != nil {
        return nil, fmt.Errorf("读取参考图失败：%w", err)
    }
    content = append(content, map[string]interface{}{
        "image": "data:image/png;base64," + base64.StdEncoding.EncodeToString(raw),
    })
}

messages = append(messages, map[string]interface{}{
    "role":    "user",
    "content": content,  // ✅ 正确：content 是数组
})
```

### 3.2 项目需要支持的协议

| 协议名称 | 用途 | 端点 | 同步/异步 |
|---------|------|------|----------|
| `chat-completion` | 文本对话（OpenAI兼容） | `/compatible-mode/v1/chat/completions` | 同步+流式 |
| `dashscope-image` | 图片生成（DashScope原生） | `/api/v1/services/aigc/multimodal-generation/generation` | 同步或异步 |
| `dashscope-video` | 视频生成（待实现） | `/api/v1/services/aigc/video-generation/generation` | 仅异步 |

---

## 四、实施建议

### 4.1 修复图片生成

**优先级：🔴 高（阻塞当前测试）**

1. 修改 `runDashScopeImageTask` 中的 messages 构建逻辑
2. 确保 `content` 始终是对象数组
3. 测试同步调用是否成功

### 4.2 实现视频生成

**优先级：🟡 中（功能扩展）**

1. 创建新函数 `runDashScopeVideoTask`
2. 实现异步任务提交和轮询
3. 处理视频URL下载和转换
4. 添加协议常量 `ChannelInterfaceDashScopeVideo`

### 4.3 完善模型发现

**优先级：🟢 低（优化体验）**

1. 在 `backend/internal/provider/bailian/discovery.go` 中补充视频模型定义
2. 确保插件系统正确返回所有支持的模型
3. 前端 `model-protocols.ts` 添加视频协议选项

---

## 五、待确认事项

### 5.1 需要访问官方文档确认

- [ ] 视频生成的确切端点路径
- [ ] 视频模型的参数格式（duration, resolution等）
- [ ] 参考生视频（r2v）的输入格式
- [ ] 视频任务的平均生成时间

### 5.2 需要实际测试验证

- [ ] 图片生成的修复是否成功
- [ ] 异步轮询的最佳间隔时间
- [ ] 视频URL的有效期
- [ ] 4K图片生成的token消耗

---

## 六、相关资源

- **官方文档首页：** https://platform.qianwenai.com/docs
- **图片生成API：** https://help.aliyun.com/zh/model-studio/qwen-image-api
- **视频生成API：** （待查找确切链接）
- **异步任务管理：** https://platform.qianwenai.com/docs/developer-guides/run-and-scale/async-task-management
- **错误码文档：** https://help.aliyun.com/zh/model-studio/error-code

---

**文档版本：** v1.0  
**最后更新：** 2026-01-20  
**维护者：** Claude (open-ai-canvas 项目)
