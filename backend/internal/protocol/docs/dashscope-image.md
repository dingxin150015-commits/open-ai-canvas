# {{NAME}}

DashScope 图片现通过 `dashscope-image.yingce-plugin` 独立声明式插件加载。插件负责精确模型的请求校验、消息和响应映射，宿主负责鉴权、出站安全与资源持久化。本文件保留协议迁移说明；当前完整执行合同以插件包的 Manifest 与文档为准。渠道 Base URL 保留实际账号区域或 MaaS 主机，只替换兼容路径后缀。

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

本协议的 execution 为 `declarative`，使用 v3 包合同和通用表达式引擎。消息顺序、尺寸转换、参考图限制和响应映射由插件声明；匿名结果下载由宿主执行。原 Go Provider 仅保留为测试对照，不进入生产调用链。新增模型必须补充精确合同与成功、拒绝夹具，再晋升 Ready。
