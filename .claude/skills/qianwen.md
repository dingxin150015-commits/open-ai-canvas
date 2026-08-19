---
name: qianwen
description: 查询千问/百炼模型的API格式和使用方法
---

# 千问/百炼模型技能

快速查询千问AI平台（百炼/DashScope）的模型API格式。

## 使用方法

调用此技能：`/qianwen [查询内容]`

示例：
- `/qianwen 图片生成格式`
- `/qianwen wan2.7-image-pro参数`
- `/qianwen 视频生成异步调用`
- `/qianwen happyhorse-1.1-t2v`

## 核心知识点

### 统一端点

| 类型 | 端点 | 异步Header |
|------|------|-----------|
| 图片 | `/api/v1/services/aigc/multimodal-generation/generation` | 可选：`X-DashScope-Async: enable` |
| 视频 | `/api/v1/services/aigc/video-generation/video-synthesis` | **必须**：`X-DashScope-Async: enable` |

Base URL: `https://dashscope.aliyuncs.com`

### 图片生成核心格式

**⚠️ 关键：content 必须是对象数组，不能是字符串！**

```json
{
  "model": "qwen-image-3.0-pro",
  "input": {
    "messages": [
      {
        "role": "user",
        "content": [
          {"text": "描述文本"},
          {"image": "data:image/png;base64,..."}  // 可选
        ]
      }
    ]
  },
  "parameters": {
    "size": "1024*1024"
  }
}
```

### 视频生成核心格式

**⚠️ 所有视频模型必须异步调用！**

```json
{
  "model": "happyhorse-1.1-t2v",
  "input": {
    "prompt": "描述文本",
    "media": [  // I2V/R2V时需要
      {"type": "first_frame", "url": "https://..."}
    ]
  },
  "parameters": {
    "duration": 5,
    "ratio": "16:9",
    "SR": 720  // HappyHorse用SR，Wan用resolution
  }
}
```

### 图片模型速查

| 模型 | 调用 | 最大分辨率 | 特殊能力 |
|------|------|-----------|---------|
| wan2.7-image-pro | 同步+异步 | 4096×4096 | 组图(12张)、色彩控制、思考模式 |
| wan2.7-image | 同步+异步 | 2048×2048 | 思考模式 |
| qwen-image-3.0-pro | 仅同步 | 2048×2048 | 智能改写、文字渲染 |
| qwen-image-3.0 | 仅同步 | 2048×2048 | 智能改写 |
| qwen-image-2.0-pro | 仅同步 | 2048×2048 | - |
| qwen-image-2.0 | 仅同步 | 2048×2048 | - |

### 本项目使用说明

**已实现功能：**
- ✅ DashScope 图片生成（同步模式）
- ✅ 支持全部6个图片模型
- ✅ 自动下载临时图片并转换为 data URL
- ✅ 正常显示、下载和保存到OSS

**使用方法：**
1. 在管理后台配置百炼系统渠道的 API Key
2. 启用对应的图片生成模型（协议选择"DashScope 图片"）
3. 在创作页面选择模型生成图片

**已知限制：**
- ❌ 不支持蒙版编辑功能
- ⚠️ 图片 URL 有效期24小时（已自动处理）
- ⚠️ 使用同步模式，响应较慢时请耐心等待

**技术细节：**
- 协议类型：`dashscope-image`
- API 端点：`/api/v1/services/aigc/multimodal-generation/generation`
- Content-Type：`application/json`
- 返回格式：`{"mode": "image", "images": [{"dataUrl": "..."}]}`

### 视频模型速查

| 模型 | 类型 | 时长 | 角色引用 |
|------|------|------|---------|
| happyhorse-1.1-t2v | 文生视频 | 1-15秒 | - |
| happyhorse-1.1-i2v | 图生视频 | 1-15秒 | - |
| happyhorse-1.1-r2v | 参考生视频 | 1-15秒 | `character1`, `character2` |
| wan2.7-t2v-2026-06-12 | 文生视频 | 2-15秒 | - |
| wan2.7-i2v-2026-04-25 | 图生视频 | 2-15秒 | - |
| wan2.7-r2v-2026-06-12 | 参考生视频 | 2-15秒 | `Video 1`, `Image 1` |

### 异步任务轮询

```bash
# 1. 提交任务
POST /api/v1/services/aigc/video-generation/video-synthesis
Header: X-DashScope-Async: enable
→ 返回 task_id

# 2. 轮询状态
GET /api/v1/tasks/{task_id}
→ task_status: PENDING/RUNNING/SUCCEEDED/FAILED
```

**轮询策略：**
- 前30秒：每3秒一次
- 30秒后：每10秒一次
- 超时：10分钟

### 关键参数

**图片生成（Wan系列专属）：**
- `thinking_mode`: bool - 思考模式（默认true）
- `enable_sequential`: bool - 组图生成（1-12张）
- `color_palette`: array - 色彩控制（3-10色，总和100%）

**图片生成（Qwen系列专属）：**
- `prompt_extend`: bool - 智能改写（默认true，+3-5秒）
- `prompt_extend_mode`: string - "direct"或"agent"（仅3.0系列）

**视频生成（通用）：**
- `duration`: int - 时长（秒）
- `ratio`: string - "16:9", "9:16", "1:1", "4:3", "3:4"
- `SR`: int - 720或1080（HappyHorse）
- `resolution`: string - "720P"或"1080P"（Wan）

### 常见错误

1. **"Input should be a valid list: input.messages.0.content"**
   - ❌ `content: "文本"`
   - ✅ `content: [{"text": "文本"}]`

2. **"Model not exist"**
   - 检查模型名拼写
   - 确认模型在支持列表中

3. **视频生成失败**
   - 检查是否设置 `X-DashScope-Async: enable`
   - 检查media数组的type字段是否正确

### 图片URL有效期

⚠️ **24小时** - 必须立即下载并转换为data URL或上传到自己的OSS

### 角色引用规则

**HappyHorse（R2V）：**
```json
"media": [
  {"type": "reference_image", "url": "..."},  // character1
  {"type": "reference_image", "url": "..."}   // character2
]
"prompt": "character1 和 character2 在花园散步"
```

**Wan（R2V）：**
```json
"media": [
  {"type": "reference_video", "url": "..."},  // Video 1
  {"type": "reference_image", "url": "..."}   // Image 1
]
"prompt": "Video 1 和 Image 1 在咖啡厅交谈"
```

---

## 完整文档

详见：`.claude/qianwen-models-api-reference-complete.md`

## 官方文档

https://platform.qianwenai.com/docs
