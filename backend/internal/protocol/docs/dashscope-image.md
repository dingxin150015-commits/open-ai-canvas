# {{NAME}}

DashScope 图片是影策连接阿里云百炼原生图片生成接口的宿主内置协议。Registry 负责协议发现、能力分类和管理后台展示；真实请求由模型专属 Provider 执行，不能把 Qwen Image 与其他万相图片模型压缩为同一套通用参数。渠道 Base URL 应保留账号实际使用的区域或 MaaS 主机，影策只在协议边界替换兼容路径，不擅自切换账号、区域或公共域名。

## 接口

{{OPERATIONS}}

当前同步接口使用 `POST /api/v1/services/aigc/multimodal-generation/generation`。请求使用 Bearer API Key，但密钥只由后端渠道读取，不进入浏览器 URL、任务正文或日志。成功响应必须包含可识别图片结果；签名结果 URL 有时效，进入影策后应尽快物化为受控资源。

## 模型

- `qwen-image-3.0-pro` 使用严格的单轮 multimodal message：文生图是一段 text，图片编辑是 1–3 段 image 后跟一段 text。
- Qwen Image 3.0 的 `auto` 表示省略上游 `size`；比例和内部 `宽x高` 只在 Adapter 边界转换为官方 `宽*高`。
- 其他 DashScope 图片模型不能自动继承 Qwen 的输出数量、参考图、thinking、prompt extend、质量、透明背景或 mask 合同。
- 官方目录发现只表示模型存在；只有 supportStatus=ready、持久能力完整并有精确序列化测试的模型才能执行。

## 参数

{{PARAMETERS}}

影策能力 JSON 是 UI、任务准入和 Provider 的共同真相。缺失、损坏或不完整配置会 fail closed。参考图必须满足模型声明的 MIME、大小和数量限制；不支持的质量、透明背景、蒙版或输出格式必须明确拒绝，禁止静默忽略。上游返回的签名 URL 属于不同信任域，下载时不得转发模型 API Key、自定义租户头或其他渠道凭据。

## 官方

- 阿里云百炼模型服务与图片生成 API 文档：<https://help.aliyun.com/zh/model-studio/image-generation-api-reference>
- Qwen 文生图和图片编辑的精确字段应以账号区域对应的当前官方 API 页面为准。
- 模型价格、免费额度和区域可用性会变化，不在协议元数据中固化；任务执行仍由影策显式价格档和授权门禁控制。

{{CONTRACT}}

## 当前实现边界

本协议的 execution 为 `host:dashscope-image`，不会进入声明式 JSON 运行器。这样可以保留模型级消息顺序、尺寸转换、参考图限制、同步结果解析和安全下载逻辑。新增模型必须先补官方证据、能力 Manifest、Adapter 和拒绝测试，再从 planned 晋升 ready。
