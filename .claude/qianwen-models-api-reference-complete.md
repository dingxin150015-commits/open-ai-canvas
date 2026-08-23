# 百炼/千问模型完整 API 参考文档

> 基于官方文档完整版（2026-08-19）

## 文档说明

本文档基于**百炼千问官方文档压缩包**的完整内容整理，涵盖所有图片生成和视频生成模型的准确API格式。

**官方文档地址：** https://platform.qianwenai.com/docs  
**文档版本：** 2026-08-19  
**目标模型：**
- 图片生成：`wan2.7-image-pro`, `qwen-image-3.0-pro` 等
- 视频生成：HappyHorse 系列、Wan 系列

---

## 核心发现总结

### 1. 统一的异步任务架构

**所有视频生成模型都使用相同的异步任务模式：**

```
1. 创建任务 → 返回 task_id
2. 轮询状态 → GET /api/v1/tasks/{task_id}
3. 获取结果 → task_status = "SUCCEEDED"
```

### 2. 统一的端点

| 类型 | 端点 | Header要求 |
|------|------|-----------|
| 图片生成 | `/api/v1/services/aigc/multimodal-generation/generation` | 同步：无；异步：`X-DashScope-Async: enable` |
| 视频生成 | `/api/v1/services/aigc/video-generation/video-synthesis` | **必须** `X-DashScope-Async: enable` |

### 3. 请求体格式的统一规律

**图片生成**使用 `messages` 数组（多模态消息格式）：
```json
{
  "model": "qwen-image-3.0-pro",
  "input": {
    "messages": [
      {
        "role": "user",
        "content": [
          {"text": "描述"},
          {"image": "data:..."}  // 可选
        ]
      }
    ]
  }
}
```

**视频生成**使用 `prompt` + `media` 数组（专用格式）：
```json
{
  "model": "happyhorse-1.1-t2v",
  "input": {
    "prompt": "描述",
    "media": [  // 可选，用于i2v/r2v
      {
        "type": "first_frame",
        "url": "https://..."
      }
    ]
  },
  "parameters": {
    "duration": 5,
    "ratio": "16:9",
    "SR": 720
  }
}
```

---

## 一、图片生成模型（已确认）

### 1.1 模型完整列表

| 模型 ID | 调用方式 | 最大分辨率 | 特殊能力 |
|---------|---------|-----------|---------|
| `wan2.7-image-pro` | 同步+异步 | 4096×4096 | 组图(12张)、色彩控制、思考模式 |
| `wan2.7-image` | 同步+异步 | 2048×2048 | 思考模式、组图 |
| `qwen-image-3.0-pro` | 仅同步 | 2048×2048 | 智能改写、文字渲染 |
| `qwen-image-3.0` | 仅同步 | 2048×2048 | 智能改写 |
| `qwen-image-2.0-pro` | 仅同步 | 2048×2048 | - |
| `qwen-image-2.0` | 仅同步 | 2048×2048 | - |

### 1.2 请求格式（同步）

**端点：**
```
POST https://dashscope.aliyuncs.com/api/v1/services/aigc/multimodal-generation/generation
Authorization: Bearer $DASHSCOPE_API_KEY
Content-Type: application/json
```

**完整请求体（包含所有可能的参数）：**

```json
{
  "model": "qwen-image-3.0-pro",
  "input": {
    "messages": [
      {
        "role": "system",  // 可选
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
            "text": "一只可爱的橘色小猫"
          },
          {
            "image": "data:image/png;base64,iVBORw0KGgo..."  // 可选：参考图
          }
        ]
      }
    ]
  },
  "parameters": {
    // === 通用参数 ===
    "size": "1024*1024",           // 格式："宽*高" 或 "1K"/"2K"/"4K"
    "n": 1,                        // 生成数量
    
    // === Qwen系列专属 ===
    "prompt_extend": true,         // 智能改写（默认true，增加3-5秒）
    "prompt_extend_mode": "direct", // "direct" 或 "agent"（仅qwen-image-3.0系列）
    
    // === Wan系列专属 ===
    "thinking_mode": true,         // 思考模式（默认true）
    "enable_sequential": false,    // 组图生成（1-12张）
    "color_palette": [             // 色彩控制（3-10种颜色）
      {"hex": "#C2D1E6", "ratio": "23.51%"},
      {"hex": "#CDD8E9", "ratio": "20.13%"}
      // ... 总和必须为100%
    ]
  }
}
```

**⚠️ 关键点：**
1. ✅ `content` **必须是数组**，不能是字符串
2. ✅ 文本包装在 `{"text": "..."}` 中
3. ✅ 图片包装在 `{"image": "data:..."}` 中
4. ✅ 可以有多个 content 元素（文本+多张图片）

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
    "output_image_count": 1
  },
  "request_id": "xxx"
}
```

**⚠️ 图片URL有效期：24小时**

### 1.4 异步调用（Wan系列）

**步骤1：创建任务**
```http
POST /api/v1/services/aigc/image-generation/generation
Authorization: Bearer $DASHSCOPE_API_KEY
X-DashScope-Async: enable
Content-Type: application/json

{
  "model": "wan2.7-image-pro",
  "input": {...},
  "parameters": {...}
}
```

**响应：**
```json
{
  "output": {
    "task_id": "77093787-xxx",
    "task_status": "PENDING"
  }
}
```

**步骤2：轮询状态**
```http
GET /api/v1/tasks/77093787-xxx
Authorization: Bearer $DASHSCOPE_API_KEY
```

**响应（成功）：**
```json
{
  "output": {
    "task_id": "77093787-xxx",
    "task_status": "SUCCEEDED",
    "submit_time": "2026-03-31 23:04:46.166",
    "end_time": "2026-03-31 23:05:11.664",
    "choices": [
      {
        "message": {
          "content": [
            {"image": "https://...", "type": "image"}
          ]
        }
      }
    ]
  }
}
```

---

## 二、视频生成模型（完整确认）

### 2.1 模型完整列表

#### HappyHorse 系列

| 模型 ID | 类型 | 最大时长 | 分辨率 | 特点 |
|---------|------|---------|--------|------|
| `happyhorse-1.1-t2v` | 文生视频 | 15秒 | 720P/1080P | 物理真实、运动流畅 |
| `happyhorse-1.1-i2v` | 图生视频（首帧） | 15秒 | 720P/1080P | 自动跟随首帧宽高比 |
| `happyhorse-1.1-r2v` | 参考生视频 | 15秒 | 720P/1080P | 多角色融合（最多5个） |

#### Wan 系列

| 模型 ID | 类型 | 最大时长 | 分辨率 | 特点 |
|---------|------|---------|--------|------|
| `wan2.7-t2v-2026-06-12` | 文生视频 | 2-15秒 | 720P/1080P | 音画同步、多镜头 |
| `wan2.7-i2v-2026-04-25` | 图生视频 | 2-15秒 | 720P/1080P | 首帧/尾帧/双帧模式 |
| `wan2.7-r2v-2026-06-12` | 参考生视频 | 2-15秒 | 720P/1080P | 多角色、声音克隆 |

**所有视频模型都必须使用异步调用。**

### 2.2 统一的异步调用流程

#### 步骤1：创建任务

**端点（所有视频模型统一）：**
```
POST https://dashscope.aliyuncs.com/api/v1/services/aigc/video-generation/video-synthesis
Authorization: Bearer $DASHSCOPE_API_KEY
X-DashScope-Async: enable
Content-Type: application/json
```

#### 步骤2：轮询任务状态

**端点（通用）：**
```
GET https://dashscope.aliyuncs.com/api/v1/tasks/{task_id}
Authorization: Bearer $DASHSCOPE_API_KEY
```

**任务状态枚举：**
- `PENDING` - 排队中
- `RUNNING` - 生成中
- `SUCCEEDED` - 成功
- `FAILED` - 失败

### 2.3 HappyHorse 文生视频（T2V）

**模型：** `happyhorse-1.1-t2v`

**请求体：**
```json
{
  "model": "happyhorse-1.1-t2v",
  "input": {
    "prompt": "一只小猫在草地上奔跑，阳光明媚"
  },
  "parameters": {
    "duration": 5,        // 整数，1-15秒
    "ratio": "16:9",      // "16:9", "9:16", "1:1", "4:3", "3:4"
    "SR": 720,            // 720 或 1080
    "seed": 42            // 可选，固定随机种子
  }
}
```

**响应（创建任务）：**
```json
{
  "output": {
    "task_id": "video-task-xxx",
    "task_status": "PENDING"
  },
  "request_id": "xxx"
}
```

**响应（任务完成）：**
```json
{
  "output": {
    "task_id": "video-task-xxx",
    "task_status": "SUCCEEDED",
    "submit_time": "2026-01-20 10:00:00",
    "end_time": "2026-01-20 10:02:30",
    "video_url": [
      "https://dashscope-result-xxx.oss-cn-xxx.aliyuncs.com/xxx.mp4?Expires=xxx"
    ]
  },
  "usage": {
    "output_video_duration": 5,
    "video_count": 1,
    "ratio": "16:9",
    "SR": 720
  },
  "request_id": "xxx"
}
```

### 2.4 HappyHorse 图生视频（I2V）

**模型：** `happyhorse-1.1-i2v`

**请求体：**
```json
{
  "model": "happyhorse-1.1-i2v",
  "input": {
    "prompt": "让这张图片中的小猫动起来",  // 可选
    "media": [
      {
        "type": "first_frame",
        "url": "https://example.com/cat.jpg"
      }
    ]
  },
  "parameters": {
    "duration": 5,
    "SR": 1080,
    "seed": 42  // 可选
  }
}
```

**说明：**
- 输出视频宽高比自动跟随输入首帧图像
- 不需要指定 `ratio` 参数

### 2.5 HappyHorse 参考生视频（R2V）

**模型：** `happyhorse-1.1-r2v`

**请求体：**
```json
{
  "model": "happyhorse-1.1-r2v",
  "input": {
    "prompt": "character1 和 character2 在花园里散步聊天",
    "media": [
      {
        "type": "reference_image",
        "url": "https://example.com/person1.jpg"  // character1
      },
      {
        "type": "reference_image",
        "url": "https://example.com/person2.jpg"  // character2
      }
    ]
  },
  "parameters": {
    "duration": 10,
    "ratio": "16:9",
    "SR": 1080
  }
}
```

**关键点：**
- 最多支持 **5个角色**
- 在 prompt 中使用 `character1`, `character2` 等引用
- `media` 数组的顺序对应 character 编号

### 2.6 Wan 文生视频（T2V）

**模型：** `wan2.7-t2v-2026-06-12`

**请求体：**
```json
{
  "model": "wan2.7-t2v-2026-06-12",
  "input": {
    "prompt": "一只橘猫在阳光下慵懒地伸懒腰",
    "audio_url": "https://example.com/bgm.mp3"  // 可选：自定义音频
  },
  "parameters": {
    "duration": 8,           // 2-15秒，整数
    "ratio": "16:9",         // 可选："16:9"等
    "resolution": "1080P",   // "720P" 或 "1080P"
    "prompt_rewrite": true,  // 可选：提示词改写
    "watermark": false       // 可选：添加水印
  }
}
```

**响应格式：** 与 HappyHorse 相同

### 2.7 Wan 图生视频（I2V）

**模型：** `wan2.7-i2v-2026-04-25`

**请求体（首帧模式）：**
```json
{
  "model": "wan2.7-i2v-2026-04-25",
  "input": {
    "prompt": "让画面动起来，小猫开始奔跑",
    "media": [
      {
        "type": "first_frame",
        "url": "https://example.com/cat.jpg"
      }
    ]
  },
  "parameters": {
    "duration": 5,
    "resolution": "1080P"
  }
}
```

**请求体（尾帧模式）：**
```json
{
  "model": "wan2.7-i2v-2026-04-25",
  "input": {
    "prompt": "从起始状态过渡到结束状态",
    "media": [
      {
        "type": "last_frame",
        "url": "https://example.com/end.jpg"
      }
    ]
  },
  "parameters": {
    "duration": 5,
    "resolution": "720P"
  }
}
```

**请求体（首尾帧模式）：**
```json
{
  "model": "wan2.7-i2v-2026-04-25",
  "input": {
    "prompt": "从白天过渡到夜晚",
    "media": [
      {
        "type": "first_frame",
        "url": "https://example.com/day.jpg"
      },
      {
        "type": "last_frame",
        "url": "https://example.com/night.jpg"
      }
    ]
  },
  "parameters": {
    "duration": 10,
    "resolution": "1080P"
  }
}
```

### 2.8 Wan 参考生视频（R2V）

**模型：** `wan2.7-r2v-2026-06-12`

**请求体：**
```json
{
  "model": "wan2.7-r2v-2026-06-12",
  "input": {
    "prompt": "Video 1 和 Image 1 在咖啡厅里愉快交谈",
    "media": [
      {
        "type": "reference_video",
        "url": "https://example.com/person1.mp4"  // Video 1
      },
      {
        "type": "reference_image",
        "url": "https://example.com/person2.jpg"  // Image 1
      }
    ],
    "reference_voice": "https://example.com/voice.mp3"  // 可选：声音克隆
  },
  "parameters": {
    "duration": 12,
    "ratio": "16:9",
    "resolution": "1080P"
  }
}
```

**Wan 的角色引用规则：**
- 视频类型使用：`Video 1`, `Video 2`, ...
- 图片类型使用：`Image 1`, `Image 2`, ...
- **视频和图片分开计数**

---

## 三、与当前项目的对比与修复

### 3.1 图片生成的问题（已诊断）

**当前错误代码（provider_dashscope.go:22-36）：**

```go
messages = append(messages, map[string]interface{}{
    "role":    "user",
    "content": input.Prompt,  // ❌ 字符串
})
```

**正确代码：**

```go
// 构建 content 数组
content := []map[string]interface{}{}

// 合并 system prompt 和 user prompt
promptText := input.Prompt
if systemPrompt := strings.TrimSpace(input.Config.SystemPrompt); systemPrompt != "" {
    promptText = systemPrompt + "\n\n" + promptText
}

// 添加文本
content = append(content, map[string]interface{}{
    "text": promptText,  // ✅ 包装在对象中
})

// 添加参考图（如果有）
if len(input.ReferenceImages) > 0 {
    if len(input.ReferenceImages) > 1 {
        return nil, errors.New("DashScope 图片协议当前只支持 1 张参考图")
    }
    raw, _, err := mediaBytes(input.ReferenceImages[0])
    if err != nil {
        return nil, fmt.Errorf("读取参考图失败：%w", err)
    }
    content = append(content, map[string]interface{}{
        "image": "data:image/png;base64," + base64.StdEncoding.EncodeToString(raw),
    })
}

// 构建 messages
messages = append(messages, map[string]interface{}{
    "role":    "user",
    "content": content,  // ✅ 数组
})
```

### 3.2 视频生成的实现方案

**新增函数：** `runDashScopeVideoTask`

**位置：** `backend/internal/service/provider_dashscope.go`

**实现要点：**

1. **判断视频类型**（T2V/I2V/R2V）
2. **构建 media 数组**（根据 input.ReferenceImages/Videos）
3. **异步提交任务**
4. **轮询任务状态**（建议间隔：前30秒每3秒，之后每10秒）
5. **下载视频URL并转换为data URL**

**伪代码框架：**

```go
func runDashScopeVideoTask(ctx context.Context, input canvasGenerationInput) (map[string]interface{}, error) {
    // 1. 构建请求体
    body := map[string]interface{}{
        "model": input.Config.Model,
        "input": map[string]interface{}{
            "prompt": input.Prompt,
        },
        "parameters": map[string]interface{}{
            "duration": parseDuration(input.Config.Duration),
            "ratio": input.Config.AspectRatio,
            "SR": parseResolution(input.Config.Resolution),
        },
    }
    
    // 2. 处理 media（I2V/R2V）
    if len(input.ReferenceImages) > 0 || len(input.ReferenceVideos) > 0 {
        media := buildMediaArray(input)
        body["input"].(map[string]interface{})["media"] = media
    }
    
    // 3. 异步提交
    var taskResp struct {
        Output struct {
            TaskID string `json:"task_id"`
        } `json:"output"`
    }
    err := postDashScopeVideoAsync(ctx, input.Config, body, &taskResp)
    if err != nil {
        return nil, err
    }
    
    // 4. 轮询状态
    result, err := pollDashScopeTask(ctx, input.Config, taskResp.Output.TaskID)
    if err != nil {
        return nil, err
    }
    
    // 5. 下载视频
    videoURL := result.Output.VideoURL[0]
    data, mimeType, err := getExternalBinary(ctx, videoURL)
    if err != nil {
        return nil, err
    }
    
    return map[string]interface{}{
        "mode": "video",
        "videos": []string{dataURL(mimeType, data)},
    }, nil
}
```

### 3.3 项目需要支持的协议

| 协议名称 | 用途 | 端点 | 格式 |
|---------|------|------|------|
| `chat-completion` | 文本对话 | `/compatible-mode/v1/chat/completions` | OpenAI标准 |
| `dashscope-image` | 图片生成 | `/api/v1/services/aigc/multimodal-generation/generation` | messages数组 |
| `dashscope-video` | 视频生成 | `/api/v1/services/aigc/video-generation/video-synthesis` | prompt+media数组 |

---

## 四、关键参数对比表

### 4.1 图片生成参数

| 参数 | 类型 | 说明 | Wan系列 | Qwen系列 |
|------|------|------|---------|----------|
| `size` | string | 分辨率（"宽*高"或"1K"/"2K"/"4K"） | ✅ | ✅ |
| `n` | int | 生成数量 | ✅ | ✅ |
| `prompt_extend` | bool | 智能改写 | ❌ | ✅ |
| `prompt_extend_mode` | string | 改写模式（"direct"/"agent"） | ❌ | qwen-3.0系列 |
| `thinking_mode` | bool | 思考模式 | ✅ | ❌ |
| `enable_sequential` | bool | 组图生成 | ✅ | ❌ |
| `color_palette` | array | 色彩控制 | ✅ | ❌ |

### 4.2 视频生成参数

| 参数 | 类型 | 说明 | HappyHorse | Wan |
|------|------|------|------------|-----|
| `duration` | int | 时长（秒） | 1-15 | 2-15 |
| `ratio` | string | 宽高比 | ✅ | ✅ |
| `SR` | int | 分辨率档位 | 720/1080 | - |
| `resolution` | string | 分辨率 | - | "720P"/"1080P" |
| `seed` | int | 随机种子 | ✅ | ✅ |
| `prompt_rewrite` | bool | 提示词改写 | ❌ | ✅ |
| `watermark` | bool | 水印 | ❌ | ✅ |
| `audio_url` | string | 自定义音频 | ❌ | ✅ |
| `reference_voice` | string | 声音克隆 | ❌ | ✅ (R2V) |

---

## 五、实施优先级

### 🔴 优先级1：修复图片生成（立即）

**任务清单：**
- [x] 诊断错误根因（content必须是数组）
- [ ] 修改 `runDashScopeImageTask` 函数
- [ ] 测试同步调用（qwen-image-3.0-pro）
- [ ] 测试异步调用（wan2.7-image-pro）
- [ ] 提交代码

**预计时间：** 30分钟

### 🟡 优先级2：实现视频生成（功能扩展）

**任务清单：**
- [ ] 创建 `runDashScopeVideoTask` 函数
- [ ] 实现异步轮询机制（`pollDashScopeTask`）
- [ ] 实现 media 数组构建（`buildMediaArray`）
- [ ] 添加协议常量 `ChannelInterfaceDashScopeVideo`
- [ ] 前端添加视频协议选项
- [ ] 测试 T2V/I2V/R2V 三种模式

**预计时间：** 2-3小时

### 🟢 优先级3：完善模型发现

**任务清单：**
- [ ] 在 `bailian/discovery.go` 中添加所有视频模型定义
- [ ] 实现视频模型的参数配置
- [ ] 添加视频时长、分辨率选项
- [ ] 完善错误处理

**预计时间：** 1小时

---

## 六、代码实现示例

### 6.1 图片生成修复（完整代码）

```go
func runDashScopeImageTask(ctx context.Context, input canvasGenerationInput) (map[string]interface{}, error) {
    if input.Mask != nil {
        return nil, errors.New("DashScope 图片协议不支持蒙版编辑")
    }

    // 构建 content 数组
    content := []map[string]interface{}{}

    // 合并 system prompt
    promptText := input.Prompt
    if systemPrompt := strings.TrimSpace(input.Config.SystemPrompt); systemPrompt != "" {
        promptText = systemPrompt + "\n\n" + promptText
    }

    // 添加文本
    content = append(content, map[string]interface{}{
        "text": promptText,
    })

    // 添加参考图
    if len(input.ReferenceImages) > 0 {
        if len(input.ReferenceImages) > 1 {
            return nil, errors.New("DashScope 图片协议当前只支持 1 张参考图")
        }
        raw, _, err := mediaBytes(input.ReferenceImages[0])
        if err != nil {
            return nil, fmt.Errorf("读取参考图失败：%w", err)
        }
        content = append(content, map[string]interface{}{
            "image": "data:image/png;base64," + base64.StdEncoding.EncodeToString(raw),
        })
    }

    // 构建请求体
    body := map[string]interface{}{
        "model": input.Config.Model,
        "input": map[string]interface{}{
            "messages": []map[string]interface{}{
                {
                    "role":    "user",
                    "content": content,
                },
            },
        },
        "parameters": map[string]interface{}{},
    }

    // 设置参数
    if size := normalizeDashScopeImageSize(input.Config.Size); size != "" {
        body["parameters"].(map[string]interface{})["size"] = size
    }

    // 同步调用
    var response dashScopeResponse
    if err := postDashScopeJSONSync(ctx, input.Config, "/api/v1/services/aigc/multimodal-generation/generation", body, &response); err != nil {
        return nil, err
    }

    // 检查状态
    if response.Output.TaskStatus == "FAILED" || len(response.Output.Choices) == 0 {
        return nil, errors.New("DashScope 图片生成失败")
    }

    // 下载图片
    images, err := dashScopeImageDataURLs(ctx, input.Config, response)
    if err != nil {
        return nil, err
    }

    return map[string]interface{}{"mode": "image", "images": images}, nil
}
```

### 6.2 视频生成实现（框架代码）

```go
func runDashScopeVideoTask(ctx context.Context, input canvasGenerationInput) (map[string]interface{}, error) {
    // 构建请求体
    body := map[string]interface{}{
        "model": input.Config.Model,
        "input": map[string]interface{}{
            "prompt": input.Prompt,
        },
        "parameters": map[string]interface{}{
            "duration": getDuration(input.Config),
            "ratio":    getAspectRatio(input.Config),
        },
    }

    // 添加分辨率参数（HappyHorse用SR，Wan用resolution）
    if strings.HasPrefix(input.Config.Model, "happyhorse") {
        body["parameters"].(map[string]interface{})["SR"] = getResolution(input.Config)
    } else {
        body["parameters"].(map[string]interface{})["resolution"] = getResolutionString(input.Config)
    }

    // 构建 media 数组（I2V/R2V）
    if len(input.ReferenceImages) > 0 || len(input.ReferenceVideos) > 0 {
        media, err := buildVideoMediaArray(input)
        if err != nil {
            return nil, err
        }
        body["input"].(map[string]interface{})["media"] = media
    }

    // 异步提交
    taskID, err := submitDashScopeVideoTask(ctx, input.Config, body)
    if err != nil {
        return nil, err
    }

    // 轮询任务
    result, err := pollDashScopeVideoTask(ctx, input.Config, taskID)
    if err != nil {
        return nil, err
    }

    // 下载视频
    videoURL := result.Output.VideoURL[0]
    data, mimeType, err := getExternalBinary(withProviderRequestKind(ctx, "download"), videoURL)
    if err != nil {
        return nil, fmt.Errorf("下载视频失败：%w", err)
    }

    return map[string]interface{}{
        "mode":   "video",
        "videos": []string{dataURL(mimeType, data)},
    }, nil
}

func buildVideoMediaArray(input canvasGenerationInput) ([]map[string]interface{}, error) {
    media := []map[string]interface{}{}

    // 判断类型：首帧、尾帧、参考
    if isImageToVideo(input) {
        // I2V：first_frame
        url, err := uploadOrGetURL(input.ReferenceImages[0])
        if err != nil {
            return nil, err
        }
        media = append(media, map[string]interface{}{
            "type": "first_frame",
            "url":  url,
        })
    } else if isReferenceToVideo(input) {
        // R2V：reference_image / reference_video
        for _, img := range input.ReferenceImages {
            url, err := uploadOrGetURL(img)
            if err != nil {
                return nil, err
            }
            media = append(media, map[string]interface{}{
                "type": "reference_image",
                "url":  url,
            })
        }
        for _, vid := range input.ReferenceVideos {
            url, err := uploadOrGetURL(vid)
            if err != nil {
                return nil, err
            }
            media = append(media, map[string]interface{}{
                "type": "reference_video",
                "url":  url,
            })
        }
    }

    return media, nil
}

func pollDashScopeVideoTask(ctx context.Context, config providerConfig, taskID string) (*dashScopeVideoResult, error) {
    maxAttempts := 60 // 最多轮询10分钟
    interval := 3 * time.Second

    for i := 0; i < maxAttempts; i++ {
        select {
        case <-ctx.Done():
            return nil, ctx.Err()
        case <-time.After(interval):
            var result dashScopeVideoResult
            err := getDashScopeJSON(ctx, config, "/api/v1/tasks/"+taskID, &result)
            if err != nil {
                return nil, err
            }

            switch result.Output.TaskStatus {
            case "SUCCEEDED":
                return &result, nil
            case "FAILED":
                return nil, fmt.Errorf("任务失败：%s", result.Output.Message)
            case "PENDING", "RUNNING":
                // 动态调整间隔
                if i > 10 {
                    interval = 10 * time.Second
                }
                continue
            default:
                return nil, fmt.Errorf("未知任务状态：%s", result.Output.TaskStatus)
            }
        }
    }

    return nil, errors.New("任务超时")
}
```

---

## 七、API使用建议

### 7.1 错误处理

**常见错误码：**
- `InvalidParameter` - 参数错误（检查模型名、参数格式）
- `Throttling` - 限流（实施退避重试）
- `Unauthorized` - 认证失败（检查API Key）
- `DataInspectionFailed` - 内容审核失败（修改prompt）

**重试策略：**
```go
func retryWithBackoff(fn func() error, maxRetries int) error {
    for i := 0; i < maxRetries; i++ {
        err := fn()
        if err == nil {
            return nil
        }
        if !isRetryable(err) {
            return err
        }
        time.Sleep(time.Duration(math.Pow(2, float64(i))) * time.Second)
    }
    return errors.New("max retries exceeded")
}
```

### 7.2 性能优化

**图片URL处理：**
- ✅ 立即下载并转换为 data URL
- ✅ 或上传到自己的 OSS
- ❌ 不要直接使用临时URL（24小时后失效）

**视频轮询策略：**
- 前30秒：每3秒轮询一次
- 30秒后：每10秒轮询一次
- 超时时间：10分钟

**并发控制：**
- 图片生成：建议并发数 ≤ 5
- 视频生成：建议并发数 ≤ 2（生成时间长）

---

## 八、文档更新记录

| 版本 | 日期 | 更新内容 |
|------|------|---------|
| v2.0 | 2026-01-20 | 基于官方文档完整版重写，补充所有视频模型详细格式 |
| v1.0 | 2026-01-20 | 初版，基于浏览器快照和部分文档 |

---

**文档完整度：** ✅ 100%（基于官方2026-08-19版本）  
**代码实现状态：** ⏳ 图片修复待执行，视频功能待开发  
**下一步：** 执行优先级1任务（修复图片生成）
