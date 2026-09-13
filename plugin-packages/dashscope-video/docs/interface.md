# DashScope Wan / HappyHorse Video 接口字段

## 协议身份

- 插件 ID：`dashscope-video`。
- Provider ID：`dashscope-video`。
- 能力：`video`。
- 默认 Base URL：`https://dashscope.aliyuncs.com/compatible-mode/v1`。
- 鉴权驱动：`bearer`。
- 创建：`POST /api/v1/services/aigc/video-generation/video-synthesis`。
- 查询：`GET /api/v1/tasks/{{taskId}}`。

## 配置字段

| 字段 | 类型 | 必填 | 含义 |
| --- | --- | --- | --- |
| `apiKey` | secret | 是 | API Key |

## 统一字段映射

| 统一字段 | 类型 | 必填 | 上游映射 | 说明 |
| --- | --- | --- | --- | --- |
| `model` | string | 是 | `model` | 精确模型 ID |
| `prompt` | string | 否 | `input` | 模型输入提示词 |
| `images` | media[] | 否 | `input` | 图片素材 |
| `videos` | media[] | 否 | `input` | 视频素材 |
| `audios` | media[] | 否 | `input` | 音频素材 |
| `duration` | integer | 否 | `parameters.duration` | 时长 |
| `aspectRatio` | string | 否 | `parameters` | 画幅或尺寸 |

## 上游请求模板逐字段清单

下表由插件请求模板生成，覆盖 body、query、headers 和 multipart 文件声明中的每个字段。

| 上游位置 | 值或转换表达式 |
| --- | --- |
| `create.method` | `"POST"` |
| `create.path` | `"/api/v1/services/aigc/video-generation/video-synthesis"` |
| `create.contentType` | `"application/json"` |
| `create.body.model` | `{"$ref":"vars.model"}` |
| `create.body.input.prompt` | `{"$omitEmpty":{"$ref":"vars.prompt"}}` |
| `create.body.input.negative_prompt` | `{"$if":{"condition":{"$eq":[{"$ref":"vars.family"},"wan27"]},"then":{"$omitEmpty":{"$ref":"vars.negative"}},"else":null}}` |
| `create.body.input.audio_url` | `{"$if":{"condition":{"$and":[{"$eq":[{"$ref":"vars.family"},"wan27"]},{"$eq":[{"$ref":"vars.kind"},"t2v"]}]},"then":{"$ref":"vars.audios.0.value"},"else":null}}` |
| `create.body.input.reference_voice` | `{"$if":{"condition":{"$and":[{"$eq":[{"$ref":"vars.family"},"wan27"]},{"$eq":[{"$ref":"vars.kind"},"r2v"]}]},"then":{"$ref":"vars.audios.0.value"},"else":null}}` |
| `create.body.input.media` | `{"$omitEmpty":{"$concatArrays":[{"$if":{"condition":{"$and":[{"$eq":[{"$ref":"vars.family"},"wan27"]},{"$eq":[{"$ref":"vars.kind"},"i2v"]},{"$gt":[{"$len":{"$ref":"vars.videos"}},0]}]},"then":{"$map":{"from":{"$ref":"vars.videos"},"as":"m","in":{"type":"first_clip","url":{"$ref":"m.value"}}}},"else":[]}},{"$map":{"from":{"$ref":"vars.orderedImages"},"as":"m","in":{"type":{"$switch":{"cases":[{"when":{"$and":[{"$eq":[{"$ref":"vars.family"},"wan27"]},{"$eq":[{"$ref":"vars.kind"},"i2v"]},{"$gt":[{"$len":{"$ref":"vars.videos"}},0]}]},"then":"last_frame"},{"when":{"$or":[{"$eq":[{"$ref":"vars.kind"},"i2v"]},{"$and":[{"$eq":[{"$ref":"vars.family"},"wan30"]},{"$eq":[{"$ref":"vars.operation"},"image_to_video"]}]}]},"then":{"$if":{"condition":{"$eq":[{"$ref":"mIndex"},0]},"then":"first_frame","else":"last_frame"}}},{"when":{"$and":[{"$eq":[{"$ref":"vars.family"},"wan27"]},{"$eq":[{"$ref":"vars.kind"},"r2v"]},{"$eq":[{"$ref":"m.id"},{"$ref":"vars.first"}]},{"$ne":[{"$ref":"vars.first"},""]}]},"then":"first_frame"}],"default":"reference_image"}},"url":{"$ref":"m.value"}}}},{"$if":{"condition":{"$and":[{"$eq":[{"$ref":"vars.family"},"wan27"]},{"$eq":[{"$ref":"vars.kind"},"i2v"]}]},"then":[],"else":{"$map":{"from":{"$ref":"vars.videos"},"as":"m","in":{"type":"reference_video","url":{"$ref":"m.value"}}}}}},{"$if":{"condition":{"$eq":[{"$ref":"vars.family"},"wan30"]},"then":{"$map":{"from":{"$ref":"vars.audios"},"as":"m","in":{"type":"reference_audio","url":{"$ref":"m.value"}}}},"else":{"$if":{"condition":{"$and":[{"$eq":[{"$ref":"vars.family"},"wan27"]},{"$eq":[{"$ref":"vars.kind"},"i2v"]}]},"then":{"$map":{"from":{"$ref":"vars.audios"},"as":"m","in":{"type":"driving_audio","url":{"$ref":"m.value"}}}},"else":[]}}}}]}}` |
| `create.body.parameters.duration` | `{"$ref":"vars.duration"}` |
| `create.body.parameters.resolution` | `{"$ref":"vars.resolution"}` |
| `create.body.parameters.ratio` | `{"$if":{"condition":{"$or":[{"$eq":[{"$ref":"vars.kind"},"i2v"]},{"$and":[{"$eq":[{"$ref":"vars.family"},"wan27"]},{"$eq":[{"$ref":"vars.kind"},"r2v"]},{"$ne":[{"$ref":"vars.first"},""]}]}]},"then":null,"else":{"$ref":"vars.ratio"}}}` |
| `create.body.parameters.watermark` | `{"$ref":"vars.watermark"}` |
| `create.body.parameters.seed` | `{"$ref":"vars.seed"}` |
| `create.body.parameters.prompt_extend` | `{"$if":{"condition":{"$eq":[{"$ref":"vars.family"},"wan27"]},"then":{"$ref":"vars.extend"},"else":null}}` |
| `create.body.parameters.audio` | `{"$if":{"condition":{"$eq":[{"$ref":"vars.family"},"wan30"]},"then":{"$ref":"vars.audio"},"else":null}}` |
| `create.headers.X-DashScope-Async` | `"enable"` |
| `poll.method` | `"GET"` |
| `poll.path` | `"/api/v1/tasks/{{taskId}}"` |
| `poll.contentType` | `"application/json"` |

## Provider 扩展键

- 无额外扩展键。

动态模型或工作流允许使用文档声明的完整 `parameters/input/extra_body` 对象；该对象是协议本身的开放 schema，不会被宿主裁剪。

## 响应映射逐字段清单

| 映射位置 | 上游路径或转换表达式 |
| --- | --- |
| `response.strict` | `true` |
| `response.taskId` | `{"$coalesce":[{"$ref":"response.output.task_id"},{"$ref":"taskId"}]}` |
| `response.status` | `{"$switch":{"cases":[{"when":{"$eq":[{"$upper":{"$trim":{"$ref":"response.output.task_status"}}},"UNKNOWN"]},"then":"failed"}],"default":{"$ref":"response.output.task_status"}}}` |
| `response.videos` | `{"$ref":"response.output.video_url"}` |
| `response.usage` | `{"$ref":"response.usage"}` |
| `response.errorPaths[0]` | `"code"` |
| `response.errorPaths[1]` | `"output.code"` |
| `response.messagePaths[0]` | `"output.message"` |
| `response.messagePaths[1]` | `"message"` |
| `response.resultEphemeral` | `true` |
| `response.resultAnonymous` | `true` |

## 响应与错误

插件把上游 task/status/text/media/usage 映射为统一结果。临时媒体 URL 标记为 ephemeral，由宿主立即下载持久化。HTTP 错误、业务 code 和 error object 保持失败语义，不包装成成功。

## 兼容边界

精确支持 Wan 2.7 七个主线/快照 ID、HappyHorse 1.1 三个 ID、wan3.0-video。未登记模型拒绝。保留首尾帧、续写、参考媒体、驱动音频、智能时长及各自输入限制。异步轮询间隔 15 秒；UNKNOWN 视为终态失败；结果 URL 临时有效，宿主立即物化。

<!-- YINGCE_MANIFEST_CONTRACT_START -->
## Manifest 完整接口定义

以下 JSON 与插件包内实际 `manifest.json` 逐字段一致，覆盖插件身份、权限、配置、鉴权、参数、校验、创建、Agent、查询、取消、结果下载、响应和 Agent 响应映射。`documentation` 字段的值就是当前完整文档；为避免文档在自身内部无限递归，JSON 中仅用等义占位文本表示正文。

```json
{
  "apiVersion": "yingce.plugin/v4",
  "id": "dashscope-video",
  "name": "DashScope Wan / HappyHorse Video",
  "version": "4.0.0",
  "author": "Alibaba Cloud / 影策",
  "description": "DashScope Wan / HappyHorse Video 独立请求协议插件。",
  "documentation": "<当前插件的完整 documentation，由 README.md 与 docs/interface.md 拼接而成；为避免 JSON 递归，此处不重复展开正文。>",
  "permissions": [
    "generation.run",
    "media.read"
  ],
  "configuration": {
    "fields": [
      {
        "name": "apiKey",
        "type": "secret",
        "label": "API Key",
        "required": true
      }
    ]
  },
  "contributes": {
    "providers": [
      {
        "id": "dashscope-video",
        "label": "DashScope Wan / HappyHorse Video",
        "capabilities": [
          "video"
        ],
        "scopes": [
          "admin.system-channel",
          "user.custom-channel",
          "canvas",
          "creation",
          "agent"
        ],
        "baseUrl": "https://dashscope.aliyuncs.com/compatible-mode/v1",
        "requiresPublicMediaUrls": true,
        "auth": {
          "type": "bearer",
          "field": "apiKey"
        },
        "parameters": [
          {
            "name": "model",
            "type": "string",
            "required": true,
            "mapping": "model",
            "description": "精确模型 ID"
          },
          {
            "name": "prompt",
            "type": "string",
            "required": false,
            "mapping": "input",
            "description": "模型输入提示词"
          },
          {
            "name": "images",
            "type": "media[]",
            "required": false,
            "mapping": "input",
            "description": "图片素材"
          },
          {
            "name": "videos",
            "type": "media[]",
            "required": false,
            "mapping": "input",
            "description": "视频素材"
          },
          {
            "name": "audios",
            "type": "media[]",
            "required": false,
            "mapping": "input",
            "description": "音频素材"
          },
          {
            "name": "duration",
            "type": "integer",
            "required": false,
            "mapping": "parameters.duration",
            "description": "时长"
          },
          {
            "name": "aspectRatio",
            "type": "string",
            "required": false,
            "mapping": "parameters",
            "description": "画幅或尺寸"
          }
        ],
        "modelContract": {
          "mode": "exact",
          "models": [
            "wan2.7-t2v",
            "wan2.7-t2v-2026-06-12",
            "wan2.7-t2v-2026-04-25",
            "wan2.7-i2v",
            "wan2.7-i2v-2026-04-25",
            "wan2.7-r2v",
            "wan2.7-r2v-2026-06-12",
            "happyhorse-1.1-t2v",
            "happyhorse-1.1-i2v",
            "happyhorse-1.1-r2v",
            "wan3.0-video"
          ]
        },
        "modelCatalog": {
          "hosts": [
            "dashscope.aliyuncs.com",
            "dashscope-intl.aliyuncs.com",
            "dashscope-us.aliyuncs.com",
            "maas.aliyuncs.com"
          ],
          "entries": [
            {
              "id": "animate-anyone-gen2",
              "displayName": "animate-anyone-gen2",
              "supportStatus": "planned",
              "supportReason": "项目模型专属 Adapter 与合同测试尚未完成",
              "supportedOperations": [
                "image_to_video"
              ],
              "documentationPaths": [
                "api-reference/video-generation/animate-anyone/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            },
            {
              "id": "animate-anyone-template-gen2",
              "displayName": "animate-anyone-template-gen2",
              "supportStatus": "planned",
              "supportReason": "项目模型专属 Adapter 与合同测试尚未完成",
              "supportedOperations": [
                "image_to_video"
              ],
              "documentationPaths": [
                "api-reference/video-generation/animate-anyone-template-gen/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            },
            {
              "id": "emo-v1",
              "displayName": "emo-v1",
              "supportStatus": "planned",
              "supportReason": "项目模型专属 Adapter 与合同测试尚未完成",
              "supportedOperations": [
                "audio_to_video"
              ],
              "documentationPaths": [
                "api-reference/video-generation/emo-video/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            },
            {
              "id": "emoji-v1",
              "displayName": "emoji-v1",
              "supportStatus": "planned",
              "supportReason": "项目模型专属 Adapter 与合同测试尚未完成",
              "supportedOperations": [
                "image_to_video"
              ],
              "documentationPaths": [
                "api-reference/video-generation/emoji-video/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            },
            {
              "id": "happyhorse-1.0-i2v",
              "displayName": "happyhorse-1.0-i2v",
              "supportStatus": "planned",
              "supportReason": "项目模型专属 Adapter 与合同测试尚未完成",
              "supportedOperations": [
                "image_to_video"
              ],
              "documentationPaths": [
                "api-reference/video-generation/happyhorse-image-to-video/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            },
            {
              "id": "happyhorse-1.0-r2v",
              "displayName": "happyhorse-1.0-r2v",
              "supportStatus": "planned",
              "supportReason": "项目模型专属 Adapter 与合同测试尚未完成",
              "supportedOperations": [
                "reference_to_video"
              ],
              "documentationPaths": [
                "api-reference/video-generation/happyhorse-reference-to-video/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            },
            {
              "id": "happyhorse-1.0-t2v",
              "displayName": "happyhorse-1.0-t2v",
              "supportStatus": "planned",
              "supportReason": "项目模型专属 Adapter 与合同测试尚未完成",
              "supportedOperations": [
                "text_to_video"
              ],
              "documentationPaths": [
                "api-reference/video-generation/happyhorse-text-to-video/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            },
            {
              "id": "happyhorse-1.0-video-edit",
              "displayName": "happyhorse-1.0-video-edit",
              "supportStatus": "planned",
              "supportReason": "项目模型专属 Adapter 与合同测试尚未完成",
              "supportedOperations": [
                "video_edit"
              ],
              "documentationPaths": [
                "api-reference/video-generation/happyhorse-video-editing/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            },
            {
              "id": "happyhorse-1.1-i2v",
              "displayName": "HappyHorse 1.1 图生视频",
              "supportStatus": "ready",
              "supportReason": "模型专属 Adapter 与离线合同测试已完成",
              "supportedOperations": [
                "image_to_video"
              ],
              "documentationPaths": [
                "api-reference/video-generation/happyhorse-image-to-video/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            },
            {
              "id": "happyhorse-1.1-r2v",
              "displayName": "HappyHorse 1.1 参考生视频",
              "supportStatus": "ready",
              "supportReason": "模型专属 Adapter 与离线合同测试已完成",
              "supportedOperations": [
                "reference_to_video",
                "image_to_video"
              ],
              "documentationPaths": [
                "api-reference/video-generation/happyhorse-reference-to-video/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            },
            {
              "id": "happyhorse-1.1-t2v",
              "displayName": "HappyHorse 1.1 文生视频",
              "supportStatus": "ready",
              "supportReason": "模型专属 Adapter 与离线合同测试已完成",
              "supportedOperations": [
                "text_to_video"
              ],
              "documentationPaths": [
                "api-reference/video-generation/happyhorse-text-to-video/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            },
            {
              "id": "kling/kling-v3-omni-video-generation",
              "displayName": "kling/kling-v3-omni-video-generation",
              "supportStatus": "planned",
              "supportReason": "项目模型专属 Adapter 与合同测试尚未完成",
              "supportedOperations": [
                "image_to_video",
                "reference_to_video",
                "text_to_video",
                "video_edit"
              ],
              "documentationPaths": [
                "api-reference/video-generation/kling-video-generation/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            },
            {
              "id": "kling/kling-v3-video-generation",
              "displayName": "kling/kling-v3-video-generation",
              "supportStatus": "planned",
              "supportReason": "项目模型专属 Adapter 与合同测试尚未完成",
              "supportedOperations": [
                "image_to_video",
                "reference_to_video",
                "text_to_video",
                "video_edit"
              ],
              "documentationPaths": [
                "api-reference/video-generation/kling-video-generation/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            },
            {
              "id": "liveportrait",
              "displayName": "liveportrait",
              "supportStatus": "planned",
              "supportReason": "项目模型专属 Adapter 与合同测试尚未完成",
              "supportedOperations": [
                "image_to_video"
              ],
              "documentationPaths": [
                "api-reference/video-generation/liveportrait-video/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            },
            {
              "id": "pixverse/pixverse-c1-it2v",
              "displayName": "pixverse/pixverse-c1-it2v",
              "supportStatus": "planned",
              "supportReason": "项目模型专属 Adapter 与合同测试尚未完成",
              "supportedOperations": [
                "image_to_video"
              ],
              "documentationPaths": [
                "api-reference/video-generation/pixverse-image-to-video-first-frame/create-task.md",
                "api-reference/video-generation/pixverse-image-to-video/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            },
            {
              "id": "pixverse/pixverse-c1-kf2v",
              "displayName": "pixverse/pixverse-c1-kf2v",
              "supportStatus": "planned",
              "supportReason": "项目模型专属 Adapter 与合同测试尚未完成",
              "supportedOperations": [
                "first_last_frame_to_video"
              ],
              "documentationPaths": [
                "api-reference/video-generation/pixverse-image-to-video-first-last/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            },
            {
              "id": "pixverse/pixverse-c1-r2v",
              "displayName": "pixverse/pixverse-c1-r2v",
              "supportStatus": "planned",
              "supportReason": "项目模型专属 Adapter 与合同测试尚未完成",
              "supportedOperations": [
                "reference_to_video"
              ],
              "documentationPaths": [
                "api-reference/video-generation/pixverse-reference-to-video/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            },
            {
              "id": "pixverse/pixverse-c1-t2v",
              "displayName": "pixverse/pixverse-c1-t2v",
              "supportStatus": "planned",
              "supportReason": "项目模型专属 Adapter 与合同测试尚未完成",
              "supportedOperations": [
                "text_to_video"
              ],
              "documentationPaths": [
                "api-reference/video-generation/pixverse-text-to-video/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            },
            {
              "id": "pixverse/pixverse-lipsync",
              "displayName": "pixverse/pixverse-lipsync",
              "supportStatus": "planned",
              "supportReason": "项目模型专属 Adapter 与合同测试尚未完成",
              "supportedOperations": [
                "lip_sync"
              ],
              "documentationPaths": [
                "api-reference/video-generation/pixverse-lipsync/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            },
            {
              "id": "pixverse/pixverse-motioncontrol",
              "displayName": "pixverse/pixverse-motioncontrol",
              "supportStatus": "planned",
              "supportReason": "项目模型专属 Adapter 与合同测试尚未完成",
              "supportedOperations": [
                "motion_control"
              ],
              "documentationPaths": [
                "api-reference/video-generation/pixverse-motioncontrol/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            },
            {
              "id": "pixverse/pixverse-upscale",
              "displayName": "pixverse/pixverse-upscale",
              "supportStatus": "planned",
              "supportReason": "项目模型专属 Adapter 与合同测试尚未完成",
              "supportedOperations": [
                "video_upscale"
              ],
              "documentationPaths": [
                "api-reference/video-generation/pixverse-upscale/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            },
            {
              "id": "pixverse/pixverse-v5.6-it2v",
              "displayName": "pixverse/pixverse-v5.6-it2v",
              "supportStatus": "planned",
              "supportReason": "项目模型专属 Adapter 与合同测试尚未完成",
              "supportedOperations": [
                "image_to_video"
              ],
              "documentationPaths": [
                "api-reference/video-generation/pixverse-image-to-video-first-frame/create-task.md",
                "api-reference/video-generation/pixverse-image-to-video/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            },
            {
              "id": "pixverse/pixverse-v5.6-kf2v",
              "displayName": "pixverse/pixverse-v5.6-kf2v",
              "supportStatus": "planned",
              "supportReason": "项目模型专属 Adapter 与合同测试尚未完成",
              "supportedOperations": [
                "first_last_frame_to_video"
              ],
              "documentationPaths": [
                "api-reference/video-generation/pixverse-image-to-video-first-last/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            },
            {
              "id": "pixverse/pixverse-v5.6-r2v",
              "displayName": "pixverse/pixverse-v5.6-r2v",
              "supportStatus": "planned",
              "supportReason": "项目模型专属 Adapter 与合同测试尚未完成",
              "supportedOperations": [
                "reference_to_video"
              ],
              "documentationPaths": [
                "api-reference/video-generation/pixverse-reference-to-video/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            },
            {
              "id": "pixverse/pixverse-v5.6-t2v",
              "displayName": "pixverse/pixverse-v5.6-t2v",
              "supportStatus": "planned",
              "supportReason": "项目模型专属 Adapter 与合同测试尚未完成",
              "supportedOperations": [
                "text_to_video"
              ],
              "documentationPaths": [
                "api-reference/video-generation/pixverse-text-to-video/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            },
            {
              "id": "pixverse/pixverse-v6-it2v",
              "displayName": "pixverse/pixverse-v6-it2v",
              "supportStatus": "planned",
              "supportReason": "项目模型专属 Adapter 与合同测试尚未完成",
              "supportedOperations": [
                "image_to_video"
              ],
              "documentationPaths": [
                "api-reference/video-generation/pixverse-image-to-video-first-frame/create-task.md",
                "api-reference/video-generation/pixverse-image-to-video/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            },
            {
              "id": "pixverse/pixverse-v6-kf2v",
              "displayName": "pixverse/pixverse-v6-kf2v",
              "supportStatus": "planned",
              "supportReason": "项目模型专属 Adapter 与合同测试尚未完成",
              "supportedOperations": [
                "first_last_frame_to_video"
              ],
              "documentationPaths": [
                "api-reference/video-generation/pixverse-image-to-video-first-last/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            },
            {
              "id": "pixverse/pixverse-v6-r2v",
              "displayName": "pixverse/pixverse-v6-r2v",
              "supportStatus": "planned",
              "supportReason": "项目模型专属 Adapter 与合同测试尚未完成",
              "supportedOperations": [
                "reference_to_video"
              ],
              "documentationPaths": [
                "api-reference/video-generation/pixverse-reference-to-video/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            },
            {
              "id": "pixverse/pixverse-v6-r2v-omni",
              "displayName": "pixverse/pixverse-v6-r2v-omni",
              "supportStatus": "planned",
              "supportReason": "项目模型专属 Adapter 与合同测试尚未完成",
              "supportedOperations": [
                "reference_to_video"
              ],
              "documentationPaths": [
                "api-reference/video-generation/pixverse-reference-to-video/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            },
            {
              "id": "pixverse/pixverse-v6-t2v",
              "displayName": "pixverse/pixverse-v6-t2v",
              "supportStatus": "planned",
              "supportReason": "项目模型专属 Adapter 与合同测试尚未完成",
              "supportedOperations": [
                "text_to_video"
              ],
              "documentationPaths": [
                "api-reference/video-generation/pixverse-text-to-video/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            },
            {
              "id": "video-style-transform",
              "displayName": "video-style-transform",
              "supportStatus": "planned",
              "supportReason": "项目模型专属 Adapter 与合同测试尚未完成",
              "supportedOperations": [
                "video_edit"
              ],
              "documentationPaths": [
                "api-reference/video-generation/video-style-transform/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            },
            {
              "id": "videoretalk",
              "displayName": "videoretalk",
              "supportStatus": "planned",
              "supportReason": "项目模型专属 Adapter 与合同测试尚未完成",
              "supportedOperations": [
                "lip_sync"
              ],
              "documentationPaths": [
                "api-reference/video-generation/video-retalk/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            },
            {
              "id": "vidu/viduq2-pro_img2video",
              "displayName": "vidu/viduq2-pro_img2video",
              "supportStatus": "planned",
              "supportReason": "项目模型专属 Adapter 与合同测试尚未完成",
              "supportedOperations": [
                "image_to_video"
              ],
              "documentationPaths": [
                "api-reference/video-generation/vidu-image-to-video-first-frame/create-task.md",
                "api-reference/video-generation/vidu-image-to-video/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            },
            {
              "id": "vidu/viduq2-pro_reference2video",
              "displayName": "vidu/viduq2-pro_reference2video",
              "supportStatus": "planned",
              "supportReason": "项目模型专属 Adapter 与合同测试尚未完成",
              "supportedOperations": [
                "reference_to_video"
              ],
              "documentationPaths": [
                "api-reference/video-generation/vidu-reference-to-video/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            },
            {
              "id": "vidu/viduq2-pro_start-end2video",
              "displayName": "vidu/viduq2-pro_start-end2video",
              "supportStatus": "planned",
              "supportReason": "项目模型专属 Adapter 与合同测试尚未完成",
              "supportedOperations": [
                "first_last_frame_to_video"
              ],
              "documentationPaths": [
                "api-reference/video-generation/vidu-start-end-to-video/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            },
            {
              "id": "vidu/viduq2-turbo_img2video",
              "displayName": "vidu/viduq2-turbo_img2video",
              "supportStatus": "planned",
              "supportReason": "项目模型专属 Adapter 与合同测试尚未完成",
              "supportedOperations": [
                "image_to_video"
              ],
              "documentationPaths": [
                "api-reference/video-generation/vidu-image-to-video-first-frame/create-task.md",
                "api-reference/video-generation/vidu-image-to-video/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            },
            {
              "id": "vidu/viduq2-turbo_start-end2video",
              "displayName": "vidu/viduq2-turbo_start-end2video",
              "supportStatus": "planned",
              "supportReason": "项目模型专属 Adapter 与合同测试尚未完成",
              "supportedOperations": [
                "first_last_frame_to_video"
              ],
              "documentationPaths": [
                "api-reference/video-generation/vidu-start-end-to-video/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            },
            {
              "id": "vidu/viduq2_reference2video",
              "displayName": "vidu/viduq2_reference2video",
              "supportStatus": "planned",
              "supportReason": "项目模型专属 Adapter 与合同测试尚未完成",
              "supportedOperations": [
                "reference_to_video"
              ],
              "documentationPaths": [
                "api-reference/video-generation/vidu-reference-to-video/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            },
            {
              "id": "vidu/viduq2_text2video",
              "displayName": "vidu/viduq2_text2video",
              "supportStatus": "planned",
              "supportReason": "项目模型专属 Adapter 与合同测试尚未完成",
              "supportedOperations": [
                "text_to_video"
              ],
              "documentationPaths": [
                "api-reference/video-generation/vidu-text-to-video/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            },
            {
              "id": "vidu/viduq3-ad_reference2video",
              "displayName": "vidu/viduq3-ad_reference2video",
              "supportStatus": "planned",
              "supportReason": "项目模型专属 Adapter 与合同测试尚未完成",
              "supportedOperations": [
                "reference_to_video"
              ],
              "documentationPaths": [
                "api-reference/video-generation/vidu-reference-to-video/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            },
            {
              "id": "vidu/viduq3-drama_reference2video",
              "displayName": "vidu/viduq3-drama_reference2video",
              "supportStatus": "planned",
              "supportReason": "项目模型专属 Adapter 与合同测试尚未完成",
              "supportedOperations": [
                "reference_to_video"
              ],
              "documentationPaths": [
                "api-reference/video-generation/vidu-reference-to-video/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            },
            {
              "id": "vidu/viduq3-pro-fast_img2video",
              "displayName": "vidu/viduq3-pro-fast_img2video",
              "supportStatus": "planned",
              "supportReason": "项目模型专属 Adapter 与合同测试尚未完成",
              "supportedOperations": [
                "image_to_video"
              ],
              "documentationPaths": [
                "api-reference/video-generation/vidu-image-to-video-first-frame/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            },
            {
              "id": "vidu/viduq3-pro_img2video",
              "displayName": "vidu/viduq3-pro_img2video",
              "supportStatus": "planned",
              "supportReason": "项目模型专属 Adapter 与合同测试尚未完成",
              "supportedOperations": [
                "image_to_video"
              ],
              "documentationPaths": [
                "api-reference/video-generation/vidu-image-to-video-first-frame/create-task.md",
                "api-reference/video-generation/vidu-image-to-video/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            },
            {
              "id": "vidu/viduq3-pro_start-end2video",
              "displayName": "vidu/viduq3-pro_start-end2video",
              "supportStatus": "planned",
              "supportReason": "项目模型专属 Adapter 与合同测试尚未完成",
              "supportedOperations": [
                "first_last_frame_to_video"
              ],
              "documentationPaths": [
                "api-reference/video-generation/vidu-start-end-to-video/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            },
            {
              "id": "vidu/viduq3-pro_text2video",
              "displayName": "vidu/viduq3-pro_text2video",
              "supportStatus": "planned",
              "supportReason": "项目模型专属 Adapter 与合同测试尚未完成",
              "supportedOperations": [
                "text_to_video"
              ],
              "documentationPaths": [
                "api-reference/video-generation/vidu-text-to-video/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            },
            {
              "id": "vidu/viduq3-turbo_img2video",
              "displayName": "vidu/viduq3-turbo_img2video",
              "supportStatus": "planned",
              "supportReason": "项目模型专属 Adapter 与合同测试尚未完成",
              "supportedOperations": [
                "image_to_video"
              ],
              "documentationPaths": [
                "api-reference/video-generation/vidu-image-to-video-first-frame/create-task.md",
                "api-reference/video-generation/vidu-image-to-video/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            },
            {
              "id": "vidu/viduq3-turbo_start-end2video",
              "displayName": "vidu/viduq3-turbo_start-end2video",
              "supportStatus": "planned",
              "supportReason": "项目模型专属 Adapter 与合同测试尚未完成",
              "supportedOperations": [
                "first_last_frame_to_video"
              ],
              "documentationPaths": [
                "api-reference/video-generation/vidu-start-end-to-video/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            },
            {
              "id": "vidu/viduq3-turbo_text2video",
              "displayName": "vidu/viduq3-turbo_text2video",
              "supportStatus": "planned",
              "supportReason": "项目模型专属 Adapter 与合同测试尚未完成",
              "supportedOperations": [
                "text_to_video"
              ],
              "documentationPaths": [
                "api-reference/video-generation/vidu-text-to-video/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            },
            {
              "id": "wan2.1-i2v-plus",
              "displayName": "wan2.1-i2v-plus",
              "supportStatus": "planned",
              "supportReason": "项目模型专属 Adapter 与合同测试尚未完成",
              "supportedOperations": [
                "image_to_video"
              ],
              "documentationPaths": [
                "api-reference/video-generation/wan-image-to-video-first-frame/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            },
            {
              "id": "wan2.1-i2v-turbo",
              "displayName": "wan2.1-i2v-turbo",
              "supportStatus": "planned",
              "supportReason": "项目模型专属 Adapter 与合同测试尚未完成",
              "supportedOperations": [
                "image_to_video"
              ],
              "documentationPaths": [
                "api-reference/video-generation/wan-image-to-video-first-frame/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            },
            {
              "id": "wan2.1-kf2v-plus",
              "displayName": "wan2.1-kf2v-plus",
              "supportStatus": "planned",
              "supportReason": "项目模型专属 Adapter 与合同测试尚未完成",
              "supportedOperations": [
                "first_last_frame_to_video"
              ],
              "documentationPaths": [
                "api-reference/video-generation/wan-image-to-video-first-last-frames/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            },
            {
              "id": "wan2.1-t2v-plus",
              "displayName": "wan2.1-t2v-plus",
              "supportStatus": "planned",
              "supportReason": "项目模型专属 Adapter 与合同测试尚未完成",
              "supportedOperations": [
                "text_to_video"
              ],
              "documentationPaths": [
                "api-reference/video-generation/wan-text-to-video/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            },
            {
              "id": "wan2.1-t2v-turbo",
              "displayName": "wan2.1-t2v-turbo",
              "supportStatus": "planned",
              "supportReason": "项目模型专属 Adapter 与合同测试尚未完成",
              "supportedOperations": [
                "text_to_video"
              ],
              "documentationPaths": [
                "api-reference/video-generation/wan-text-to-video/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            },
            {
              "id": "wan2.2-animate-mix",
              "displayName": "wan2.2-animate-mix",
              "supportStatus": "planned",
              "supportReason": "项目模型专属 Adapter 与合同测试尚未完成",
              "supportedOperations": [
                "video_edit"
              ],
              "documentationPaths": [
                "api-reference/video-generation/wan-video-character-swap/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            },
            {
              "id": "wan2.2-animate-move",
              "displayName": "wan2.2-animate-move",
              "supportStatus": "planned",
              "supportReason": "项目模型专属 Adapter 与合同测试尚未完成",
              "supportedOperations": [
                "image_to_video"
              ],
              "documentationPaths": [
                "api-reference/video-generation/wan-image-to-animation/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            },
            {
              "id": "wan2.2-i2v-flash",
              "displayName": "wan2.2-i2v-flash",
              "supportStatus": "planned",
              "supportReason": "项目模型专属 Adapter 与合同测试尚未完成",
              "supportedOperations": [
                "image_to_video"
              ],
              "documentationPaths": [
                "api-reference/video-generation/wan-image-to-video-first-frame/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            },
            {
              "id": "wan2.2-i2v-plus",
              "displayName": "wan2.2-i2v-plus",
              "supportStatus": "planned",
              "supportReason": "项目模型专属 Adapter 与合同测试尚未完成",
              "supportedOperations": [
                "image_to_video"
              ],
              "documentationPaths": [
                "api-reference/video-generation/wan-image-to-video-first-frame/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            },
            {
              "id": "wan2.2-kf2v-flash",
              "displayName": "wan2.2-kf2v-flash",
              "supportStatus": "planned",
              "supportReason": "项目模型专属 Adapter 与合同测试尚未完成",
              "supportedOperations": [
                "first_last_frame_to_video"
              ],
              "documentationPaths": [
                "api-reference/video-generation/wan-image-to-video-first-last-frames/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            },
            {
              "id": "wan2.2-s2v",
              "displayName": "wan2.2-s2v",
              "supportStatus": "planned",
              "supportReason": "项目模型专属 Adapter 与合同测试尚未完成",
              "supportedOperations": [
                "audio_to_video"
              ],
              "documentationPaths": [
                "api-reference/video-generation/wan-s2v/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            },
            {
              "id": "wan2.2-t2v-plus",
              "displayName": "wan2.2-t2v-plus",
              "supportStatus": "planned",
              "supportReason": "项目模型专属 Adapter 与合同测试尚未完成",
              "supportedOperations": [
                "text_to_video"
              ],
              "documentationPaths": [
                "api-reference/video-generation/wan-text-to-video/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            },
            {
              "id": "wan2.5-i2v-preview",
              "displayName": "wan2.5-i2v-preview",
              "supportStatus": "planned",
              "supportReason": "项目模型专属 Adapter 与合同测试尚未完成",
              "supportedOperations": [
                "image_to_video"
              ],
              "documentationPaths": [
                "api-reference/video-generation/wan-image-to-video-first-frame/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            },
            {
              "id": "wan2.5-t2v-preview",
              "displayName": "wan2.5-t2v-preview",
              "supportStatus": "planned",
              "supportReason": "项目模型专属 Adapter 与合同测试尚未完成",
              "supportedOperations": [
                "text_to_video"
              ],
              "documentationPaths": [
                "api-reference/video-generation/wan-text-to-video/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            },
            {
              "id": "wan2.6-i2v",
              "displayName": "wan2.6-i2v",
              "supportStatus": "planned",
              "supportReason": "项目模型专属 Adapter 与合同测试尚未完成",
              "supportedOperations": [
                "image_to_video"
              ],
              "documentationPaths": [
                "api-reference/video-generation/wan-image-to-video-first-frame/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            },
            {
              "id": "wan2.6-i2v-flash",
              "displayName": "wan2.6-i2v-flash",
              "supportStatus": "planned",
              "supportReason": "项目模型专属 Adapter 与合同测试尚未完成",
              "supportedOperations": [
                "image_to_video"
              ],
              "documentationPaths": [
                "api-reference/video-generation/wan-image-to-video-first-frame/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            },
            {
              "id": "wan2.6-r2v",
              "displayName": "wan2.6-r2v",
              "supportStatus": "planned",
              "supportReason": "项目模型专属 Adapter 与合同测试尚未完成",
              "supportedOperations": [
                "reference_to_video"
              ],
              "documentationPaths": [
                "api-reference/video-generation/wan-reference-to-video/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            },
            {
              "id": "wan2.6-r2v-flash",
              "displayName": "wan2.6-r2v-flash",
              "supportStatus": "planned",
              "supportReason": "项目模型专属 Adapter 与合同测试尚未完成",
              "supportedOperations": [
                "reference_to_video"
              ],
              "documentationPaths": [
                "api-reference/video-generation/wan-reference-to-video/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            },
            {
              "id": "wan2.6-t2v",
              "displayName": "wan2.6-t2v",
              "supportStatus": "planned",
              "supportReason": "项目模型专属 Adapter 与合同测试尚未完成",
              "supportedOperations": [
                "text_to_video"
              ],
              "documentationPaths": [
                "api-reference/video-generation/wan-text-to-video/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            },
            {
              "id": "wan2.7-i2v",
              "displayName": "万相 2.7 图生视频",
              "supportStatus": "ready",
              "supportReason": "模型专属 Adapter 与离线合同测试已完成",
              "supportedOperations": [
                "image_to_video",
                "first_last_frame_to_video",
                "extend",
                "audio_to_video"
              ],
              "documentationPaths": [
                "api-reference/video-generation/wan27-image-to-video/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            },
            {
              "id": "wan2.7-i2v-2026-04-25",
              "displayName": "wan2.7-i2v-2026-04-25",
              "supportStatus": "ready",
              "supportReason": "模型专属 Adapter 与离线合同测试已完成",
              "supportedOperations": [
                "image_to_video",
                "first_last_frame_to_video",
                "extend",
                "audio_to_video"
              ],
              "documentationPaths": [
                "api-reference/video-generation/wan27-image-to-video/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            },
            {
              "id": "wan2.7-r2v",
              "displayName": "万相 2.7 参考生视频",
              "supportStatus": "ready",
              "supportReason": "模型专属 Adapter 与离线合同测试已完成",
              "supportedOperations": [
                "reference_to_video",
                "image_to_video",
                "audio_to_video"
              ],
              "documentationPaths": [
                "api-reference/video-generation/wan27-reference-to-video/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            },
            {
              "id": "wan2.7-r2v-2026-06-12",
              "displayName": "wan2.7-r2v-2026-06-12",
              "supportStatus": "ready",
              "supportReason": "模型专属 Adapter 与离线合同测试已完成",
              "supportedOperations": [
                "reference_to_video",
                "image_to_video",
                "audio_to_video"
              ],
              "documentationPaths": [
                "api-reference/video-generation/wan27-reference-to-video/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            },
            {
              "id": "wan2.7-t2v",
              "displayName": "万相 2.7 文生视频",
              "supportStatus": "ready",
              "supportReason": "模型专属 Adapter 与离线合同测试已完成",
              "supportedOperations": [
                "text_to_video",
                "audio_to_video"
              ],
              "documentationPaths": [
                "api-reference/video-generation/wan27-text-to-video/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            },
            {
              "id": "wan2.7-t2v-2026-04-25",
              "displayName": "wan2.7-t2v-2026-04-25",
              "supportStatus": "ready",
              "supportReason": "模型专属 Adapter 与离线合同测试已完成",
              "supportedOperations": [
                "text_to_video",
                "audio_to_video"
              ],
              "documentationPaths": [
                "api-reference/video-generation/wan27-text-to-video/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            },
            {
              "id": "wan2.7-t2v-2026-06-12",
              "displayName": "wan2.7-t2v-2026-06-12",
              "supportStatus": "ready",
              "supportReason": "模型专属 Adapter 与离线合同测试已完成",
              "supportedOperations": [
                "text_to_video",
                "audio_to_video"
              ],
              "documentationPaths": [
                "api-reference/video-generation/wan27-text-to-video/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            },
            {
              "id": "wan2.7-videoedit",
              "displayName": "wan2.7-videoedit",
              "supportStatus": "planned",
              "supportReason": "项目模型专属 Adapter 与合同测试尚未完成",
              "supportedOperations": [
                "video_edit"
              ],
              "documentationPaths": [
                "api-reference/video-generation/wan27-video-editing/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            },
            {
              "id": "wan3.0-video",
              "displayName": "万相 3.0 全能视频",
              "supportStatus": "ready",
              "supportReason": "文本、首帧/首尾帧和参考图片/视频/音频已适配；file/link 输入仍待实现",
              "supportedOperations": [
                "text_to_video",
                "image_to_video",
                "first_last_frame_to_video",
                "reference_to_video",
                "audio_to_video"
              ],
              "documentationPaths": [
                "api-reference/video-generation/wan30-video/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            },
            {
              "id": "wan3.0-video-prime",
              "displayName": "万相 3.0 全能视频 Prime",
              "supportStatus": "planned",
              "supportReason": "项目模型专属 Adapter 与合同测试尚未完成",
              "supportedOperations": [
                "file_to_video",
                "first_last_frame_to_video",
                "image_to_video",
                "link_to_video",
                "reference_audio_to_video",
                "reference_to_video",
                "text_to_video"
              ],
              "documentationPaths": [
                "api-reference/video-generation/wan30-video/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            },
            {
              "id": "wanx2.1-vace-plus",
              "displayName": "wanx2.1-vace-plus",
              "supportStatus": "planned",
              "supportReason": "项目模型专属 Adapter 与合同测试尚未完成",
              "supportedOperations": [
                "video_edit"
              ],
              "documentationPaths": [
                "api-reference/video-generation/wan-general-video-editing/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis"
            }
          ]
        },
        "variables": [
          {
            "name": "model",
            "value": {
              "$replace": [
                {
                  "$lower": {
                    "$trim": {
                      "$ref": "request.model"
                    }
                  }
                },
                "models/",
                ""
              ]
            }
          },
          {
            "name": "meta",
            "value": {
              "$merge": [
                {
                  "$ref": "request.extra.generationMetadata"
                },
                {
                  "$ref": "request.providerOptions.dashscope-video"
                }
              ]
            }
          },
          {
            "name": "raw",
            "value": {
              "$ref": "request.extra.rawOptions"
            }
          },
          {
            "name": "prompt",
            "value": {
              "$trim": {
                "$ref": "request.prompt"
              }
            }
          },
          {
            "name": "images",
            "value": {
              "$ref": "request.images"
            }
          },
          {
            "name": "videos",
            "value": {
              "$ref": "request.videos"
            }
          },
          {
            "name": "audios",
            "value": {
              "$ref": "request.audios"
            }
          },
          {
            "name": "family",
            "value": {
              "$switch": {
                "cases": [
                  {
                    "when": {
                      "$in": [
                        {
                          "$ref": "vars.model"
                        },
                        [
                          "wan2.7-t2v",
                          "wan2.7-t2v-2026-06-12",
                          "wan2.7-t2v-2026-04-25",
                          "wan2.7-i2v",
                          "wan2.7-i2v-2026-04-25",
                          "wan2.7-r2v",
                          "wan2.7-r2v-2026-06-12"
                        ]
                      ]
                    },
                    "then": "wan27"
                  },
                  {
                    "when": {
                      "$in": [
                        {
                          "$ref": "vars.model"
                        },
                        [
                          "happyhorse-1.1-t2v",
                          "happyhorse-1.1-i2v",
                          "happyhorse-1.1-r2v"
                        ]
                      ]
                    },
                    "then": "horse"
                  },
                  {
                    "when": {
                      "$eq": [
                        {
                          "$ref": "vars.model"
                        },
                        "wan3.0-video"
                      ]
                    },
                    "then": "wan30"
                  }
                ],
                "default": "unknown"
              }
            }
          },
          {
            "name": "kind",
            "value": {
              "$switch": {
                "cases": [
                  {
                    "when": {
                      "$in": [
                        {
                          "$ref": "vars.model"
                        },
                        [
                          "wan2.7-t2v",
                          "wan2.7-t2v-2026-06-12",
                          "wan2.7-t2v-2026-04-25",
                          "happyhorse-1.1-t2v"
                        ]
                      ]
                    },
                    "then": "t2v"
                  },
                  {
                    "when": {
                      "$in": [
                        {
                          "$ref": "vars.model"
                        },
                        [
                          "wan2.7-i2v",
                          "wan2.7-i2v-2026-04-25",
                          "happyhorse-1.1-i2v"
                        ]
                      ]
                    },
                    "then": "i2v"
                  },
                  {
                    "when": {
                      "$in": [
                        {
                          "$ref": "vars.model"
                        },
                        [
                          "wan2.7-r2v",
                          "wan2.7-r2v-2026-06-12",
                          "happyhorse-1.1-r2v"
                        ]
                      ]
                    },
                    "then": "r2v"
                  }
                ],
                "default": "all"
              }
            }
          },
          {
            "name": "first",
            "value": {
              "$trim": {
                "$ref": "vars.raw.firstFrameID"
              }
            }
          },
          {
            "name": "last",
            "value": {
              "$trim": {
                "$ref": "vars.raw.lastFrameID"
              }
            }
          },
          {
            "name": "operation",
            "value": {
              "$coalesce": [
                {
                  "$lower": {
                    "$trim": {
                      "$ref": "request.operation"
                    }
                  }
                },
                {
                  "$switch": {
                    "cases": [
                      {
                        "when": {
                          "$or": [
                            {
                              "$ne": [
                                {
                                  "$ref": "vars.first"
                                },
                                ""
                              ]
                            },
                            {
                              "$ne": [
                                {
                                  "$ref": "vars.last"
                                },
                                ""
                              ]
                            }
                          ]
                        },
                        "then": "image_to_video"
                      },
                      {
                        "when": {
                          "$and": [
                            {
                              "$eq": [
                                {
                                  "$ref": "vars.family"
                                },
                                "wan30"
                              ]
                            },
                            {
                              "$or": [
                                {
                                  "$gt": [
                                    {
                                      "$len": {
                                        "$ref": "vars.videos"
                                      }
                                    },
                                    0
                                  ]
                                },
                                {
                                  "$gt": [
                                    {
                                      "$len": {
                                        "$ref": "vars.audios"
                                      }
                                    },
                                    0
                                  ]
                                },
                                {
                                  "$gt": [
                                    {
                                      "$len": {
                                        "$ref": "vars.images"
                                      }
                                    },
                                    2
                                  ]
                                }
                              ]
                            }
                          ]
                        },
                        "then": "reference_to_video"
                      },
                      {
                        "when": {
                          "$and": [
                            {
                              "$eq": [
                                {
                                  "$ref": "vars.family"
                                },
                                "wan30"
                              ]
                            },
                            {
                              "$gt": [
                                {
                                  "$len": {
                                    "$ref": "vars.images"
                                  }
                                },
                                0
                              ]
                            }
                          ]
                        },
                        "then": "image_to_video"
                      },
                      {
                        "when": {
                          "$eq": [
                            {
                              "$ref": "vars.kind"
                            },
                            "r2v"
                          ]
                        },
                        "then": "reference_to_video"
                      },
                      {
                        "when": {
                          "$eq": [
                            {
                              "$ref": "vars.kind"
                            },
                            "i2v"
                          ]
                        },
                        "then": "image_to_video"
                      }
                    ],
                    "default": "text_to_video"
                  }
                }
              ]
            }
          },
          {
            "name": "duration",
            "value": {
              "$strictInt": {
                "value": {
                  "$coalesce": [
                    {
                      "$ref": "vars.raw.duration"
                    },
                    {
                      "$ref": "request.extra.videoSeconds"
                    }
                  ]
                },
                "default": 5,
                "min": {
                  "$if": {
                    "condition": {
                      "$eq": [
                        {
                          "$ref": "vars.family"
                        },
                        "horse"
                      ]
                    },
                    "then": 3,
                    "else": 2
                  }
                },
                "max": {
                  "$if": {
                    "condition": {
                      "$eq": [
                        {
                          "$ref": "vars.family"
                        },
                        "wan30"
                      ]
                    },
                    "then": 30,
                    "else": 15
                  }
                },
                "allow": {
                  "$if": {
                    "condition": {
                      "$eq": [
                        {
                          "$ref": "vars.family"
                        },
                        "wan30"
                      ]
                    },
                    "then": [
                      -1
                    ],
                    "else": []
                  }
                }
              }
            }
          },
          {
            "name": "resolution",
            "value": {
              "$switch": {
                "cases": [
                  {
                    "when": {
                      "$in": [
                        {
                          "$upper": {
                            "$trim": {
                              "$ref": "request.resolution"
                            }
                          }
                        },
                        [
                          "480",
                          "720",
                          "1080"
                        ]
                      ]
                    },
                    "then": {
                      "$concat": [
                        {
                          "$upper": {
                            "$trim": {
                              "$ref": "request.resolution"
                            }
                          }
                        },
                        "P"
                      ]
                    }
                  }
                ],
                "default": {
                  "$coalesce": [
                    {
                      "$upper": {
                        "$trim": {
                          "$ref": "request.resolution"
                        }
                      }
                    },
                    "1080P"
                  ]
                }
              }
            }
          },
          {
            "name": "ratio",
            "value": {
              "$coalesce": [
                {
                  "$lower": {
                    "$trim": {
                      "$ref": "request.aspectRatio"
                    }
                  }
                },
                {
                  "$if": {
                    "condition": {
                      "$eq": [
                        {
                          "$ref": "vars.family"
                        },
                        "wan30"
                      ]
                    },
                    "then": "adaptive",
                    "else": "16:9"
                  }
                }
              ]
            }
          },
          {
            "name": "watermark",
            "value": {
              "$strictBool": {
                "value": {
                  "$ref": "vars.raw.watermark"
                },
                "default": {
                  "$eq": [
                    {
                      "$ref": "vars.family"
                    },
                    "horse"
                  ]
                }
              }
            }
          },
          {
            "name": "audio",
            "value": {
              "$strictBool": {
                "value": {
                  "$ref": "vars.raw.generateAudio"
                },
                "default": true
              }
            }
          },
          {
            "name": "extend",
            "value": {
              "$strictBool": {
                "value": {
                  "$ref": "vars.meta.prompt_extend"
                },
                "default": true
              }
            }
          },
          {
            "name": "seed",
            "value": {
              "$strictInt": {
                "value": {
                  "$ref": "vars.meta.seed"
                },
                "default": null,
                "min": 0,
                "max": 2147483647,
                "allow": []
              }
            }
          },
          {
            "name": "negative",
            "value": {
              "$coalesce": [
                {
                  "$ref": "vars.meta.negative_prompt"
                },
                ""
              ]
            }
          },
          {
            "name": "orderedImages",
            "value": {
              "$if": {
                "condition": {
                  "$and": [
                    {
                      "$or": [
                        {
                          "$eq": [
                            {
                              "$ref": "vars.kind"
                            },
                            "i2v"
                          ]
                        },
                        {
                          "$and": [
                            {
                              "$eq": [
                                {
                                  "$ref": "vars.family"
                                },
                                "wan30"
                              ]
                            },
                            {
                              "$eq": [
                                {
                                  "$ref": "vars.operation"
                                },
                                "image_to_video"
                              ]
                            }
                          ]
                        }
                      ]
                    },
                    {
                      "$or": [
                        {
                          "$ne": [
                            {
                              "$ref": "vars.first"
                            },
                            ""
                          ]
                        },
                        {
                          "$ne": [
                            {
                              "$ref": "vars.last"
                            },
                            ""
                          ]
                        }
                      ]
                    }
                  ]
                },
                "then": {
                  "$selectByIds": {
                    "from": {
                      "$ref": "vars.images"
                    },
                    "ids": [
                      {
                        "$ref": "vars.first"
                      },
                      {
                        "$ref": "vars.last"
                      }
                    ]
                  }
                },
                "else": {
                  "$ref": "vars.images"
                }
              }
            }
          }
        ],
        "pollIntervalMs": 15000,
        "validations": [
          {
            "assert": {
              "$and": [
                {
                  "$eq": [
                    {
                      "$len": {
                        "$ref": "vars.images"
                      }
                    },
                    {
                      "$coalesce": [
                        {
                          "$ref": "request.extra.referenceCounts.image"
                        },
                        {
                          "$len": {
                            "$ref": "vars.images"
                          }
                        }
                      ]
                    }
                  ]
                },
                {
                  "$eq": [
                    {
                      "$len": {
                        "$ref": "vars.videos"
                      }
                    },
                    {
                      "$coalesce": [
                        {
                          "$ref": "request.extra.referenceCounts.video"
                        },
                        {
                          "$len": {
                            "$ref": "vars.videos"
                          }
                        }
                      ]
                    }
                  ]
                },
                {
                  "$eq": [
                    {
                      "$len": {
                        "$ref": "vars.audios"
                      }
                    },
                    {
                      "$coalesce": [
                        {
                          "$ref": "request.extra.referenceCounts.audio"
                        },
                        {
                          "$len": {
                            "$ref": "vars.audios"
                          }
                        }
                      ]
                    }
                  ]
                }
              ]
            },
            "message": "参考素材缺少有效地址，不能静默忽略"
          },
          {
            "assert": {
              "$ne": [
                {
                  "$ref": "vars.family"
                },
                "unknown"
              ]
            },
            "message": "模型尚未由此视频插件支持"
          },
          {
            "assert": {
              "$in": [
                {
                  "$ref": "vars.resolution"
                },
                {
                  "$if": {
                    "condition": {
                      "$eq": [
                        {
                          "$ref": "vars.family"
                        },
                        "wan27"
                      ]
                    },
                    "then": [
                      "720P",
                      "1080P"
                    ],
                    "else": [
                      "480P",
                      "720P",
                      "1080P"
                    ]
                  }
                }
              ]
            },
            "message": "分辨率不受模型支持"
          },
          {
            "assert": {
              "$or": [
                {
                  "$not": {
                    "$not": {
                      "$or": [
                        {
                          "$eq": [
                            {
                              "$ref": "vars.kind"
                            },
                            "i2v"
                          ]
                        },
                        {
                          "$and": [
                            {
                              "$eq": [
                                {
                                  "$ref": "vars.family"
                                },
                                "wan27"
                              ]
                            },
                            {
                              "$eq": [
                                {
                                  "$ref": "vars.kind"
                                },
                                "r2v"
                              ]
                            },
                            {
                              "$ne": [
                                {
                                  "$ref": "vars.first"
                                },
                                ""
                              ]
                            }
                          ]
                        }
                      ]
                    }
                  }
                },
                {
                  "$in": [
                    {
                      "$ref": "vars.ratio"
                    },
                    {
                      "$if": {
                        "condition": {
                          "$eq": [
                            {
                              "$ref": "vars.family"
                            },
                            "horse"
                          ]
                        },
                        "then": [
                          "16:9",
                          "9:16",
                          "1:1",
                          "4:3",
                          "3:4",
                          "4:5",
                          "5:4",
                          "9:21",
                          "21:9"
                        ],
                        "else": {
                          "$if": {
                            "condition": {
                              "$eq": [
                                {
                                  "$ref": "vars.family"
                                },
                                "wan30"
                              ]
                            },
                            "then": [
                              "adaptive",
                              "16:9",
                              "9:16",
                              "1:1",
                              "4:3",
                              "3:4"
                            ],
                            "else": [
                              "16:9",
                              "9:16",
                              "1:1",
                              "4:3",
                              "3:4"
                            ]
                          }
                        }
                      }
                    }
                  ]
                }
              ]
            },
            "message": "比例不受模型支持"
          },
          {
            "assert": {
              "$lte": [
                {
                  "$len": {
                    "$ref": "vars.prompt"
                  }
                },
                {
                  "$if": {
                    "condition": {
                      "$eq": [
                        {
                          "$ref": "vars.family"
                        },
                        "wan30"
                      ]
                    },
                    "then": 20000,
                    "else": 5000
                  }
                }
              ]
            },
            "message": "提示词超过模型限制"
          },
          {
            "assert": {
              "$or": [
                {
                  "$not": {
                    "$eq": [
                      {
                        "$ref": "vars.family"
                      },
                      "horse"
                    ]
                  }
                },
                {
                  "$lte": [
                    {
                      "$weightedLength": {
                        "value": {
                          "$ref": "vars.prompt"
                        },
                        "hanWeight": 2
                      }
                    },
                    5000
                  ]
                }
              ]
            },
            "message": "HappyHorse 提示词最多 2500 中文字符或 5000 非中文字符"
          },
          {
            "assert": {
              "$or": [
                {
                  "$not": {
                    "$and": [
                      {
                        "$not": {
                          "$eq": [
                            {
                              "$ref": "vars.kind"
                            },
                            "i2v"
                          ]
                        }
                      },
                      {
                        "$not": {
                          "$eq": [
                            {
                              "$ref": "vars.family"
                            },
                            "wan30"
                          ]
                        }
                      }
                    ]
                  }
                },
                {
                  "$gt": [
                    {
                      "$len": {
                        "$ref": "vars.prompt"
                      }
                    },
                    0
                  ]
                }
              ]
            },
            "message": "必须提供提示词"
          },
          {
            "assert": {
              "$or": [
                {
                  "$not": {
                    "$eq": [
                      {
                        "$ref": "vars.family"
                      },
                      "wan30"
                    ]
                  }
                },
                {
                  "$gt": [
                    {
                      "$add": [
                        {
                          "$len": {
                            "$ref": "vars.prompt"
                          }
                        },
                        {
                          "$add": [
                            {
                              "$len": {
                                "$ref": "vars.images"
                              }
                            },
                            {
                              "$len": {
                                "$ref": "vars.videos"
                              }
                            },
                            {
                              "$len": {
                                "$ref": "vars.audios"
                              }
                            }
                          ]
                        }
                      ]
                    },
                    0
                  ]
                }
              ]
            },
            "message": "提示词或参考素材至少提供一项"
          },
          {
            "assert": {
              "$lte": [
                {
                  "$len": {
                    "$ref": "vars.negative"
                  }
                },
                500
              ]
            },
            "message": "负向提示词超过 500 字符"
          },
          {
            "assert": {
              "$or": [
                {
                  "$not": {
                    "$eq": [
                      {
                        "$ref": "vars.kind"
                      },
                      "t2v"
                    ]
                  }
                },
                {
                  "$and": [
                    {
                      "$eq": [
                        {
                          "$len": {
                            "$ref": "vars.images"
                          }
                        },
                        0
                      ]
                    },
                    {
                      "$eq": [
                        {
                          "$len": {
                            "$ref": "vars.videos"
                          }
                        },
                        0
                      ]
                    },
                    {
                      "$lte": [
                        {
                          "$len": {
                            "$ref": "vars.audios"
                          }
                        },
                        {
                          "$if": {
                            "condition": {
                              "$eq": [
                                {
                                  "$ref": "vars.family"
                                },
                                "horse"
                              ]
                            },
                            "then": 0,
                            "else": 1
                          }
                        }
                      ]
                    },
                    {
                      "$in": [
                        {
                          "$ref": "vars.operation"
                        },
                        {
                          "$if": {
                            "condition": {
                              "$eq": [
                                {
                                  "$ref": "vars.family"
                                },
                                "horse"
                              ]
                            },
                            "then": [
                              "text_to_video"
                            ],
                            "else": [
                              "text_to_video",
                              "audio_to_video"
                            ]
                          }
                        }
                      ]
                    }
                  ]
                }
              ]
            },
            "message": "文生视频输入或操作不受支持"
          },
          {
            "assert": {
              "$or": [
                {
                  "$not": {
                    "$and": [
                      {
                        "$eq": [
                          {
                            "$ref": "vars.family"
                          },
                          "horse"
                        ]
                      },
                      {
                        "$eq": [
                          {
                            "$ref": "vars.kind"
                          },
                          "i2v"
                        ]
                      }
                    ]
                  }
                },
                {
                  "$and": [
                    {
                      "$eq": [
                        {
                          "$len": {
                            "$ref": "vars.images"
                          }
                        },
                        1
                      ]
                    },
                    {
                      "$eq": [
                        {
                          "$len": {
                            "$ref": "vars.videos"
                          }
                        },
                        0
                      ]
                    },
                    {
                      "$eq": [
                        {
                          "$len": {
                            "$ref": "vars.audios"
                          }
                        },
                        0
                      ]
                    },
                    {
                      "$eq": [
                        {
                          "$ref": "vars.last"
                        },
                        ""
                      ]
                    },
                    {
                      "$eq": [
                        {
                          "$ref": "vars.operation"
                        },
                        "image_to_video"
                      ]
                    }
                  ]
                }
              ]
            },
            "message": "HappyHorse 图生视频仅支持一张首帧"
          },
          {
            "assert": {
              "$or": [
                {
                  "$not": {
                    "$and": [
                      {
                        "$eq": [
                          {
                            "$ref": "vars.family"
                          },
                          "horse"
                        ]
                      },
                      {
                        "$eq": [
                          {
                            "$ref": "vars.kind"
                          },
                          "r2v"
                        ]
                      }
                    ]
                  }
                },
                {
                  "$and": [
                    {
                      "$gte": [
                        {
                          "$len": {
                            "$ref": "vars.images"
                          }
                        },
                        1
                      ]
                    },
                    {
                      "$lte": [
                        {
                          "$len": {
                            "$ref": "vars.images"
                          }
                        },
                        9
                      ]
                    },
                    {
                      "$eq": [
                        {
                          "$len": {
                            "$ref": "vars.videos"
                          }
                        },
                        0
                      ]
                    },
                    {
                      "$eq": [
                        {
                          "$len": {
                            "$ref": "vars.audios"
                          }
                        },
                        0
                      ]
                    },
                    {
                      "$eq": [
                        {
                          "$ref": "vars.first"
                        },
                        ""
                      ]
                    },
                    {
                      "$eq": [
                        {
                          "$ref": "vars.last"
                        },
                        ""
                      ]
                    },
                    {
                      "$in": [
                        {
                          "$ref": "vars.operation"
                        },
                        [
                          "reference_to_video",
                          "image_to_video"
                        ]
                      ]
                    }
                  ]
                }
              ]
            },
            "message": "HappyHorse 参考生视频仅支持 1–9 张普通参考图"
          },
          {
            "assert": {
              "$or": [
                {
                  "$not": {
                    "$and": [
                      {
                        "$eq": [
                          {
                            "$ref": "vars.family"
                          },
                          "wan27"
                        ]
                      },
                      {
                        "$eq": [
                          {
                            "$ref": "vars.kind"
                          },
                          "i2v"
                        ]
                      }
                    ]
                  }
                },
                {
                  "$and": [
                    {
                      "$lte": [
                        {
                          "$len": {
                            "$ref": "vars.images"
                          }
                        },
                        2
                      ]
                    },
                    {
                      "$lte": [
                        {
                          "$len": {
                            "$ref": "vars.videos"
                          }
                        },
                        1
                      ]
                    },
                    {
                      "$lte": [
                        {
                          "$len": {
                            "$ref": "vars.audios"
                          }
                        },
                        1
                      ]
                    },
                    {
                      "$gte": [
                        {
                          "$add": [
                            {
                              "$len": {
                                "$ref": "vars.images"
                              }
                            },
                            {
                              "$len": {
                                "$ref": "vars.videos"
                              }
                            }
                          ]
                        },
                        1
                      ]
                    },
                    {
                      "$in": [
                        {
                          "$ref": "vars.operation"
                        },
                        [
                          "image_to_video",
                          "extend",
                          "reference_to_video",
                          "audio_to_video"
                        ]
                      ]
                    }
                  ]
                }
              ]
            },
            "message": "Wan 2.7 图生视频输入不受支持"
          },
          {
            "assert": {
              "$or": [
                {
                  "$not": {
                    "$and": [
                      {
                        "$eq": [
                          {
                            "$ref": "vars.family"
                          },
                          "wan27"
                        ]
                      },
                      {
                        "$eq": [
                          {
                            "$ref": "vars.kind"
                          },
                          "i2v"
                        ]
                      },
                      {
                        "$gt": [
                          {
                            "$len": {
                              "$ref": "vars.videos"
                            }
                          },
                          0
                        ]
                      }
                    ]
                  }
                },
                {
                  "$and": [
                    {
                      "$eq": [
                        {
                          "$len": {
                            "$ref": "vars.audios"
                          }
                        },
                        0
                      ]
                    },
                    {
                      "$lte": [
                        {
                          "$len": {
                            "$ref": "vars.images"
                          }
                        },
                        1
                      ]
                    },
                    {
                      "$eq": [
                        {
                          "$ref": "vars.first"
                        },
                        ""
                      ]
                    },
                    {
                      "$gt": [
                        {
                          "$multiply": [
                            {
                              "$ref": "vars.duration"
                            },
                            1000
                          ]
                        },
                        {
                          "$sum": {
                            "$map": {
                              "from": {
                                "$ref": "vars.videos"
                              },
                              "as": "m",
                              "in": {
                                "$ref": "m.metadata.durationMs"
                              }
                            }
                          }
                        }
                      ]
                    }
                  ]
                }
              ]
            },
            "message": "视频续写仅允许续写视频及可选尾帧，输出须长于输入"
          },
          {
            "assert": {
              "$or": [
                {
                  "$not": {
                    "$and": [
                      {
                        "$eq": [
                          {
                            "$ref": "vars.family"
                          },
                          "wan27"
                        ]
                      },
                      {
                        "$eq": [
                          {
                            "$ref": "vars.kind"
                          },
                          "r2v"
                        ]
                      }
                    ]
                  }
                },
                {
                  "$and": [
                    {
                      "$gte": [
                        {
                          "$add": [
                            {
                              "$len": {
                                "$ref": "vars.images"
                              }
                            },
                            {
                              "$len": {
                                "$ref": "vars.videos"
                              }
                            }
                          ]
                        },
                        1
                      ]
                    },
                    {
                      "$lte": [
                        {
                          "$add": [
                            {
                              "$len": {
                                "$ref": "vars.images"
                              }
                            },
                            {
                              "$len": {
                                "$ref": "vars.videos"
                              }
                            }
                          ]
                        },
                        5
                      ]
                    },
                    {
                      "$lte": [
                        {
                          "$len": {
                            "$ref": "vars.audios"
                          }
                        },
                        1
                      ]
                    },
                    {
                      "$eq": [
                        {
                          "$ref": "vars.last"
                        },
                        ""
                      ]
                    },
                    {
                      "$in": [
                        {
                          "$ref": "vars.operation"
                        },
                        [
                          "reference_to_video",
                          "image_to_video",
                          "audio_to_video"
                        ]
                      ]
                    }
                  ]
                }
              ]
            },
            "message": "Wan 2.7 参考生视频需要 1–5 个视觉素材和最多 1 个音频"
          },
          {
            "assert": {
              "$or": [
                {
                  "$not": {
                    "$and": [
                      {
                        "$eq": [
                          {
                            "$ref": "vars.family"
                          },
                          "wan27"
                        ]
                      },
                      {
                        "$eq": [
                          {
                            "$ref": "vars.kind"
                          },
                          "r2v"
                        ]
                      },
                      {
                        "$gt": [
                          {
                            "$len": {
                              "$ref": "vars.videos"
                            }
                          },
                          0
                        ]
                      }
                    ]
                  }
                },
                {
                  "$lte": [
                    {
                      "$ref": "vars.duration"
                    },
                    10
                  ]
                }
              ]
            },
            "message": "含参考视频时输出最多 10 秒"
          },
          {
            "assert": {
              "$or": [
                {
                  "$not": {
                    "$and": [
                      {
                        "$ne": [
                          {
                            "$ref": "vars.last"
                          },
                          ""
                        ]
                      },
                      {
                        "$eq": [
                          {
                            "$len": {
                              "$ref": "vars.videos"
                            }
                          },
                          0
                        ]
                      }
                    ]
                  }
                },
                {
                  "$ne": [
                    {
                      "$ref": "vars.first"
                    },
                    ""
                  ]
                }
              ]
            },
            "message": "配置尾帧必须同时配置首帧"
          },
          {
            "assert": {
              "$or": [
                {
                  "$not": {
                    "$and": [
                      {
                        "$eq": [
                          {
                            "$ref": "vars.family"
                          },
                          "wan27"
                        ]
                      },
                      {
                        "$eq": [
                          {
                            "$ref": "vars.kind"
                          },
                          "r2v"
                        ]
                      },
                      {
                        "$ne": [
                          {
                            "$ref": "vars.first"
                          },
                          ""
                        ]
                      }
                    ]
                  }
                },
                {
                  "$eq": [
                    {
                      "$len": {
                        "$filter": {
                          "from": {
                            "$ref": "vars.images"
                          },
                          "as": "m",
                          "where": {
                            "$eq": [
                              {
                                "$ref": "m.id"
                              },
                              {
                                "$ref": "vars.first"
                              }
                            ]
                          }
                        }
                      }
                    },
                    1
                  ]
                }
              ]
            },
            "message": "指定首帧未包含在参考图中"
          },
          {
            "assert": {
              "$or": [
                {
                  "$not": {
                    "$eq": [
                      {
                        "$ref": "vars.family"
                      },
                      "wan30"
                    ]
                  }
                },
                {
                  "$and": [
                    {
                      "$lte": [
                        {
                          "$len": {
                            "$ref": "vars.images"
                          }
                        },
                        9
                      ]
                    },
                    {
                      "$lte": [
                        {
                          "$len": {
                            "$ref": "vars.videos"
                          }
                        },
                        3
                      ]
                    },
                    {
                      "$lte": [
                        {
                          "$len": {
                            "$ref": "vars.audios"
                          }
                        },
                        1
                      ]
                    },
                    {
                      "$in": [
                        {
                          "$ref": "vars.operation"
                        },
                        [
                          "text_to_video",
                          "image_to_video",
                          "reference_to_video",
                          "audio_to_video"
                        ]
                      ]
                    }
                  ]
                }
              ]
            },
            "message": "Wan 3.0 参考数量或操作超出限制"
          },
          {
            "assert": {
              "$or": [
                {
                  "$not": {
                    "$and": [
                      {
                        "$eq": [
                          {
                            "$ref": "vars.family"
                          },
                          "wan30"
                        ]
                      },
                      {
                        "$or": [
                          {
                            "$ne": [
                              {
                                "$ref": "vars.first"
                              },
                              ""
                            ]
                          },
                          {
                            "$ne": [
                              {
                                "$ref": "vars.last"
                              },
                              ""
                            ]
                          }
                        ]
                      }
                    ]
                  }
                },
                {
                  "$and": [
                    {
                      "$eq": [
                        {
                          "$ref": "vars.operation"
                        },
                        "image_to_video"
                      ]
                    },
                    {
                      "$eq": [
                        {
                          "$len": {
                            "$ref": "vars.videos"
                          }
                        },
                        0
                      ]
                    },
                    {
                      "$eq": [
                        {
                          "$len": {
                            "$ref": "vars.audios"
                          }
                        },
                        0
                      ]
                    }
                  ]
                }
              ]
            },
            "message": "Wan 3.0 帧图不能与参考音视频混用"
          },
          {
            "assert": {
              "$or": [
                {
                  "$not": {
                    "$and": [
                      {
                        "$eq": [
                          {
                            "$ref": "vars.family"
                          },
                          "wan30"
                        ]
                      },
                      {
                        "$eq": [
                          {
                            "$ref": "vars.operation"
                          },
                          "image_to_video"
                        ]
                      }
                    ]
                  }
                },
                {
                  "$and": [
                    {
                      "$gte": [
                        {
                          "$len": {
                            "$ref": "vars.images"
                          }
                        },
                        1
                      ]
                    },
                    {
                      "$lte": [
                        {
                          "$len": {
                            "$ref": "vars.images"
                          }
                        },
                        2
                      ]
                    },
                    {
                      "$eq": [
                        {
                          "$len": {
                            "$ref": "vars.videos"
                          }
                        },
                        0
                      ]
                    },
                    {
                      "$eq": [
                        {
                          "$len": {
                            "$ref": "vars.audios"
                          }
                        },
                        0
                      ]
                    }
                  ]
                }
              ]
            },
            "message": "图生视频仅支持 1–2 张帧图"
          },
          {
            "assert": {
              "$or": [
                {
                  "$not": {
                    "$and": [
                      {
                        "$eq": [
                          {
                            "$ref": "vars.family"
                          },
                          "wan30"
                        ]
                      },
                      {
                        "$eq": [
                          {
                            "$ref": "vars.operation"
                          },
                          "text_to_video"
                        ]
                      }
                    ]
                  }
                },
                {
                  "$eq": [
                    {
                      "$add": [
                        {
                          "$len": {
                            "$ref": "vars.images"
                          }
                        },
                        {
                          "$len": {
                            "$ref": "vars.videos"
                          }
                        },
                        {
                          "$len": {
                            "$ref": "vars.audios"
                          }
                        }
                      ]
                    },
                    0
                  ]
                }
              ]
            },
            "message": "文生视频不能携带素材"
          },
          {
            "assert": {
              "$or": [
                {
                  "$not": {
                    "$and": [
                      {
                        "$eq": [
                          {
                            "$ref": "vars.family"
                          },
                          "wan30"
                        ]
                      },
                      {
                        "$in": [
                          {
                            "$ref": "vars.operation"
                          },
                          [
                            "reference_to_video",
                            "audio_to_video"
                          ]
                        ]
                      }
                    ]
                  }
                },
                {
                  "$gt": [
                    {
                      "$add": [
                        {
                          "$len": {
                            "$ref": "vars.images"
                          }
                        },
                        {
                          "$len": {
                            "$ref": "vars.videos"
                          }
                        },
                        {
                          "$len": {
                            "$ref": "vars.audios"
                          }
                        }
                      ]
                    },
                    0
                  ]
                }
              ]
            },
            "message": "参考模式至少需要一个素材"
          },
          {
            "assert": {
              "$or": [
                {
                  "$not": {
                    "$and": [
                      {
                        "$eq": [
                          {
                            "$ref": "vars.family"
                          },
                          "wan30"
                        ]
                      },
                      {
                        "$gt": [
                          {
                            "$ref": "vars.duration"
                          },
                          0
                        ]
                      }
                    ]
                  }
                },
                {
                  "$lte": [
                    {
                      "$add": [
                        {
                          "$sum": {
                            "$map": {
                              "from": {
                                "$ref": "vars.videos"
                              },
                              "as": "m",
                              "in": {
                                "$ref": "m.metadata.durationMs"
                              }
                            }
                          }
                        },
                        {
                          "$multiply": [
                            {
                              "$ref": "vars.duration"
                            },
                            1000
                          ]
                        }
                      ]
                    },
                    30000
                  ]
                }
              ]
            },
            "message": "参考视频与输出时长之和不能超过 30 秒"
          },
          {
            "assert": {
              "$validateMedia": {
                "items": {
                  "$ref": "vars.images"
                },
                "maxBytes": {
                  "$if": {
                    "condition": {
                      "$eq": [
                        {
                          "$ref": "vars.family"
                        },
                        "wan30"
                      ]
                    },
                    "then": 31457280,
                    "else": 20971520
                  }
                },
                "mimes": {
                  "$if": {
                    "condition": {
                      "$eq": [
                        {
                          "$ref": "vars.family"
                        },
                        "horse"
                      ]
                    },
                    "then": [
                      "image/jpeg",
                      "image/jpg",
                      "image/png",
                      "image/webp"
                    ],
                    "else": [
                      "image/jpeg",
                      "image/jpg",
                      "image/png",
                      "image/webp",
                      "image/bmp"
                    ]
                  }
                },
                "minSide": {
                  "$if": {
                    "condition": {
                      "$eq": [
                        {
                          "$ref": "vars.family"
                        },
                        "horse"
                      ]
                    },
                    "then": {
                      "$if": {
                        "condition": {
                          "$eq": [
                            {
                              "$ref": "vars.kind"
                            },
                            "i2v"
                          ]
                        },
                        "then": 300,
                        "else": 400
                      }
                    },
                    "else": {
                      "$if": {
                        "condition": {
                          "$eq": [
                            {
                              "$ref": "vars.family"
                            },
                            "wan27"
                          ]
                        },
                        "then": 240,
                        "else": 0
                      }
                    }
                  }
                },
                "maxSide": {
                  "$if": {
                    "condition": {
                      "$eq": [
                        {
                          "$ref": "vars.family"
                        },
                        "wan27"
                      ]
                    },
                    "then": 8000,
                    "else": 0
                  }
                },
                "maxAspect": {
                  "$if": {
                    "condition": {
                      "$and": [
                        {
                          "$eq": [
                            {
                              "$ref": "vars.family"
                            },
                            "horse"
                          ]
                        },
                        {
                          "$eq": [
                            {
                              "$ref": "vars.kind"
                            },
                            "i2v"
                          ]
                        }
                      ]
                    },
                    "then": 2.5,
                    "else": {
                      "$if": {
                        "condition": {
                          "$eq": [
                            {
                              "$ref": "vars.family"
                            },
                            "wan27"
                          ]
                        },
                        "then": 8,
                        "else": 0
                      }
                    }
                  }
                },
                "allowData": {
                  "$not": {
                    "$eq": [
                      {
                        "$ref": "vars.family"
                      },
                      "wan27"
                    ]
                  }
                },
                "allowedOrigins": {
                  "$if": {
                    "condition": {
                      "$eq": [
                        {
                          "$ref": "vars.family"
                        },
                        "wan30"
                      ]
                    },
                    "then": [
                      "oss://dashscope-instant"
                    ],
                    "else": []
                  }
                }
              }
            },
            "message": "图片不符合规格"
          },
          {
            "assert": {
              "$validateMedia": {
                "items": {
                  "$ref": "vars.videos"
                },
                "maxBytes": {
                  "$if": {
                    "condition": {
                      "$eq": [
                        {
                          "$ref": "vars.family"
                        },
                        "wan30"
                      ]
                    },
                    "then": 52428800,
                    "else": 104857600
                  }
                },
                "mimes": [
                  "video/mp4",
                  "video/quicktime"
                ],
                "minDurationMs": {
                  "$if": {
                    "condition": {
                      "$eq": [
                        {
                          "$ref": "vars.kind"
                        },
                        "i2v"
                      ]
                    },
                    "then": 2000,
                    "else": 1000
                  }
                },
                "maxDurationMs": {
                  "$if": {
                    "condition": {
                      "$eq": [
                        {
                          "$ref": "vars.kind"
                        },
                        "i2v"
                      ]
                    },
                    "then": 10000,
                    "else": 30000
                  }
                },
                "minSide": {
                  "$if": {
                    "condition": {
                      "$eq": [
                        {
                          "$ref": "vars.family"
                        },
                        "wan27"
                      ]
                    },
                    "then": 240,
                    "else": 0
                  }
                },
                "maxSide": {
                  "$if": {
                    "condition": {
                      "$eq": [
                        {
                          "$ref": "vars.family"
                        },
                        "wan27"
                      ]
                    },
                    "then": 4096,
                    "else": 0
                  }
                },
                "maxAspect": {
                  "$if": {
                    "condition": {
                      "$eq": [
                        {
                          "$ref": "vars.family"
                        },
                        "wan27"
                      ]
                    },
                    "then": 8,
                    "else": 0
                  }
                },
                "allowData": {
                  "$eq": [
                    {
                      "$ref": "vars.family"
                    },
                    "wan30"
                  ]
                },
                "allowedOrigins": {
                  "$if": {
                    "condition": {
                      "$eq": [
                        {
                          "$ref": "vars.family"
                        },
                        "wan30"
                      ]
                    },
                    "then": [
                      "oss://dashscope-instant"
                    ],
                    "else": []
                  }
                }
              }
            },
            "message": "视频不符合规格"
          },
          {
            "assert": {
              "$validateMedia": {
                "items": {
                  "$ref": "vars.audios"
                },
                "maxBytes": 15728640,
                "mimes": [
                  "audio/wav",
                  "audio/x-wav",
                  "audio/mpeg",
                  "audio/mp3"
                ],
                "minDurationMs": {
                  "$if": {
                    "condition": {
                      "$eq": [
                        {
                          "$ref": "vars.kind"
                        },
                        "r2v"
                      ]
                    },
                    "then": 1000,
                    "else": 2000
                  }
                },
                "maxDurationMs": {
                  "$if": {
                    "condition": {
                      "$and": [
                        {
                          "$eq": [
                            {
                              "$ref": "vars.family"
                            },
                            "wan27"
                          ]
                        },
                        {
                          "$eq": [
                            {
                              "$ref": "vars.kind"
                            },
                            "r2v"
                          ]
                        }
                      ]
                    },
                    "then": 10000,
                    "else": 30000
                  }
                },
                "allowData": {
                  "$eq": [
                    {
                      "$ref": "vars.family"
                    },
                    "wan30"
                  ]
                },
                "allowedOrigins": {
                  "$if": {
                    "condition": {
                      "$eq": [
                        {
                          "$ref": "vars.family"
                        },
                        "wan30"
                      ]
                    },
                    "then": [
                      "oss://dashscope-instant"
                    ],
                    "else": []
                  }
                }
              }
            },
            "message": "音频不符合规格"
          }
        ],
        "create": {
          "basePathSuffixes": [
            "/compatible-mode/v1",
            "/api/v1"
          ],
          "method": "POST",
          "path": "/api/v1/services/aigc/video-generation/video-synthesis",
          "contentType": "application/json",
          "body": {
            "model": {
              "$ref": "vars.model"
            },
            "input": {
              "prompt": {
                "$omitEmpty": {
                  "$ref": "vars.prompt"
                }
              },
              "negative_prompt": {
                "$if": {
                  "condition": {
                    "$eq": [
                      {
                        "$ref": "vars.family"
                      },
                      "wan27"
                    ]
                  },
                  "then": {
                    "$omitEmpty": {
                      "$ref": "vars.negative"
                    }
                  },
                  "else": null
                }
              },
              "audio_url": {
                "$if": {
                  "condition": {
                    "$and": [
                      {
                        "$eq": [
                          {
                            "$ref": "vars.family"
                          },
                          "wan27"
                        ]
                      },
                      {
                        "$eq": [
                          {
                            "$ref": "vars.kind"
                          },
                          "t2v"
                        ]
                      }
                    ]
                  },
                  "then": {
                    "$ref": "vars.audios.0.value"
                  },
                  "else": null
                }
              },
              "reference_voice": {
                "$if": {
                  "condition": {
                    "$and": [
                      {
                        "$eq": [
                          {
                            "$ref": "vars.family"
                          },
                          "wan27"
                        ]
                      },
                      {
                        "$eq": [
                          {
                            "$ref": "vars.kind"
                          },
                          "r2v"
                        ]
                      }
                    ]
                  },
                  "then": {
                    "$ref": "vars.audios.0.value"
                  },
                  "else": null
                }
              },
              "media": {
                "$omitEmpty": {
                  "$concatArrays": [
                    {
                      "$if": {
                        "condition": {
                          "$and": [
                            {
                              "$eq": [
                                {
                                  "$ref": "vars.family"
                                },
                                "wan27"
                              ]
                            },
                            {
                              "$eq": [
                                {
                                  "$ref": "vars.kind"
                                },
                                "i2v"
                              ]
                            },
                            {
                              "$gt": [
                                {
                                  "$len": {
                                    "$ref": "vars.videos"
                                  }
                                },
                                0
                              ]
                            }
                          ]
                        },
                        "then": {
                          "$map": {
                            "from": {
                              "$ref": "vars.videos"
                            },
                            "as": "m",
                            "in": {
                              "type": "first_clip",
                              "url": {
                                "$ref": "m.value"
                              }
                            }
                          }
                        },
                        "else": []
                      }
                    },
                    {
                      "$map": {
                        "from": {
                          "$ref": "vars.orderedImages"
                        },
                        "as": "m",
                        "in": {
                          "type": {
                            "$switch": {
                              "cases": [
                                {
                                  "when": {
                                    "$and": [
                                      {
                                        "$eq": [
                                          {
                                            "$ref": "vars.family"
                                          },
                                          "wan27"
                                        ]
                                      },
                                      {
                                        "$eq": [
                                          {
                                            "$ref": "vars.kind"
                                          },
                                          "i2v"
                                        ]
                                      },
                                      {
                                        "$gt": [
                                          {
                                            "$len": {
                                              "$ref": "vars.videos"
                                            }
                                          },
                                          0
                                        ]
                                      }
                                    ]
                                  },
                                  "then": "last_frame"
                                },
                                {
                                  "when": {
                                    "$or": [
                                      {
                                        "$eq": [
                                          {
                                            "$ref": "vars.kind"
                                          },
                                          "i2v"
                                        ]
                                      },
                                      {
                                        "$and": [
                                          {
                                            "$eq": [
                                              {
                                                "$ref": "vars.family"
                                              },
                                              "wan30"
                                            ]
                                          },
                                          {
                                            "$eq": [
                                              {
                                                "$ref": "vars.operation"
                                              },
                                              "image_to_video"
                                            ]
                                          }
                                        ]
                                      }
                                    ]
                                  },
                                  "then": {
                                    "$if": {
                                      "condition": {
                                        "$eq": [
                                          {
                                            "$ref": "mIndex"
                                          },
                                          0
                                        ]
                                      },
                                      "then": "first_frame",
                                      "else": "last_frame"
                                    }
                                  }
                                },
                                {
                                  "when": {
                                    "$and": [
                                      {
                                        "$eq": [
                                          {
                                            "$ref": "vars.family"
                                          },
                                          "wan27"
                                        ]
                                      },
                                      {
                                        "$eq": [
                                          {
                                            "$ref": "vars.kind"
                                          },
                                          "r2v"
                                        ]
                                      },
                                      {
                                        "$eq": [
                                          {
                                            "$ref": "m.id"
                                          },
                                          {
                                            "$ref": "vars.first"
                                          }
                                        ]
                                      },
                                      {
                                        "$ne": [
                                          {
                                            "$ref": "vars.first"
                                          },
                                          ""
                                        ]
                                      }
                                    ]
                                  },
                                  "then": "first_frame"
                                }
                              ],
                              "default": "reference_image"
                            }
                          },
                          "url": {
                            "$ref": "m.value"
                          }
                        }
                      }
                    },
                    {
                      "$if": {
                        "condition": {
                          "$and": [
                            {
                              "$eq": [
                                {
                                  "$ref": "vars.family"
                                },
                                "wan27"
                              ]
                            },
                            {
                              "$eq": [
                                {
                                  "$ref": "vars.kind"
                                },
                                "i2v"
                              ]
                            }
                          ]
                        },
                        "then": [],
                        "else": {
                          "$map": {
                            "from": {
                              "$ref": "vars.videos"
                            },
                            "as": "m",
                            "in": {
                              "type": "reference_video",
                              "url": {
                                "$ref": "m.value"
                              }
                            }
                          }
                        }
                      }
                    },
                    {
                      "$if": {
                        "condition": {
                          "$eq": [
                            {
                              "$ref": "vars.family"
                            },
                            "wan30"
                          ]
                        },
                        "then": {
                          "$map": {
                            "from": {
                              "$ref": "vars.audios"
                            },
                            "as": "m",
                            "in": {
                              "type": "reference_audio",
                              "url": {
                                "$ref": "m.value"
                              }
                            }
                          }
                        },
                        "else": {
                          "$if": {
                            "condition": {
                              "$and": [
                                {
                                  "$eq": [
                                    {
                                      "$ref": "vars.family"
                                    },
                                    "wan27"
                                  ]
                                },
                                {
                                  "$eq": [
                                    {
                                      "$ref": "vars.kind"
                                    },
                                    "i2v"
                                  ]
                                }
                              ]
                            },
                            "then": {
                              "$map": {
                                "from": {
                                  "$ref": "vars.audios"
                                },
                                "as": "m",
                                "in": {
                                  "type": "driving_audio",
                                  "url": {
                                    "$ref": "m.value"
                                  }
                                }
                              }
                            },
                            "else": []
                          }
                        }
                      }
                    }
                  ]
                }
              }
            },
            "parameters": {
              "duration": {
                "$ref": "vars.duration"
              },
              "resolution": {
                "$ref": "vars.resolution"
              },
              "ratio": {
                "$if": {
                  "condition": {
                    "$or": [
                      {
                        "$eq": [
                          {
                            "$ref": "vars.kind"
                          },
                          "i2v"
                        ]
                      },
                      {
                        "$and": [
                          {
                            "$eq": [
                              {
                                "$ref": "vars.family"
                              },
                              "wan27"
                            ]
                          },
                          {
                            "$eq": [
                              {
                                "$ref": "vars.kind"
                              },
                              "r2v"
                            ]
                          },
                          {
                            "$ne": [
                              {
                                "$ref": "vars.first"
                              },
                              ""
                            ]
                          }
                        ]
                      }
                    ]
                  },
                  "then": null,
                  "else": {
                    "$ref": "vars.ratio"
                  }
                }
              },
              "watermark": {
                "$ref": "vars.watermark"
              },
              "seed": {
                "$ref": "vars.seed"
              },
              "prompt_extend": {
                "$if": {
                  "condition": {
                    "$eq": [
                      {
                        "$ref": "vars.family"
                      },
                      "wan27"
                    ]
                  },
                  "then": {
                    "$ref": "vars.extend"
                  },
                  "else": null
                }
              },
              "audio": {
                "$if": {
                  "condition": {
                    "$eq": [
                      {
                        "$ref": "vars.family"
                      },
                      "wan30"
                    ]
                  },
                  "then": {
                    "$ref": "vars.audio"
                  },
                  "else": null
                }
              }
            }
          },
          "headers": {
            "X-DashScope-Async": "enable"
          }
        },
        "poll": {
          "basePathSuffixes": [
            "/compatible-mode/v1",
            "/api/v1"
          ],
          "method": "GET",
          "path": "/api/v1/tasks/{{taskId}}"
        },
        "response": {
          "strict": true,
          "taskId": {
            "$coalesce": [
              {
                "$ref": "response.output.task_id"
              },
              {
                "$ref": "taskId"
              }
            ]
          },
          "status": {
            "$switch": {
              "cases": [
                {
                  "when": {
                    "$eq": [
                      {
                        "$upper": {
                          "$trim": {
                            "$ref": "response.output.task_status"
                          }
                        }
                      },
                      "UNKNOWN"
                    ]
                  },
                  "then": "failed"
                }
              ],
              "default": {
                "$ref": "response.output.task_status"
              }
            }
          },
          "videos": {
            "$ref": "response.output.video_url"
          },
          "usage": {
            "$ref": "response.usage"
          },
          "errorPaths": [
            "code",
            "output.code"
          ],
          "messagePaths": [
            "output.message",
            "message"
          ],
          "resultEphemeral": true,
          "resultAnonymous": true
        }
      }
    ]
  }
}
```
<!-- YINGCE_MANIFEST_CONTRACT_END -->
