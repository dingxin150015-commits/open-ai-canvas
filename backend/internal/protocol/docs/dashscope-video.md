# {{NAME}}

DashScope 视频现通过 `dashscope-video.yingce-plugin` 独立声明式插件加载。Wan 3.0、Wan 2.7 与 HappyHorse 按精确 modelKey 声明各自的时长、分辨率、媒体角色和请求规则，不使用静默截断参数的通用 Wan 模板。本文件保留协议迁移说明；当前完整执行合同以插件包的 Manifest 与文档为准。

## 接口

{{OPERATIONS}}

创建请求使用 `POST /api/v1/services/aigc/video-generation/video-synthesis` 并携带 `X-DashScope-Async: enable`；轮询使用 `GET /api/v1/tasks/{task_id}`。影策从已配置 Base URL 中保留实际账号和区域主机，仅移除 `/compatible-mode/v1` 或旧 `/api/v1` 路径后拼接原生端点，禁止擅自改到另一区域。

## 模型

- `wan3.0-video` 支持 2–30 秒或 `-1` 智能时长，并按精确能力处理 adaptive、分辨率、声音、水印及参考音视频。
- Wan 2.7 T2V/I2V/R2V 主线和已登记快照使用不同参考结构；R2V 的图片与视频合计上限、含参考视频时输出上限、驱动音频长度都必须跨层校验。
- HappyHorse 1.1 T2V、I2V、R2V 分别限制为无参考、单首帧、1–9 张参考图；旧版和视频编辑模型不能继承 1.1 Ready 合同。
- 官方 discovery 中出现模型不等于影策已经具备执行器；planned、unsupported、deprecated 都不能启用、定价或进入用户路由。

## 参数

{{PARAMETERS}}

视频参考素材要求上游在任务存续期间可访问的 HTTPS URL，因此元数据声明 `RequiresPublicMediaURLs=true`。影策可以为归属当前用户且状态 ready 的平台资源签发短时地址；Data URL、Cookie 保护地址、未完成资源和私网地址不得直接提交。首帧、尾帧、普通参考图、参考视频和参考音频具有不同角色，必须按模型合同排序并拒绝非法混用。

## 官方

- 阿里云百炼视频生成 API 文档：<https://help.aliyun.com/zh/model-studio/video-generation-api-reference>
- 创建、查询、模型能力和价格应以账号区域对应的当前官方页面为准。
- 返回 task ID 和媒体 URL 通常有有效期；完成后应立即进入影策资源物化。真实费用以供应商账单为准，影策价格档只用于任务预授权和内部结算。

{{CONTRACT}}

## 当前实现边界

本协议的 execution 为 `declarative`，使用 v3 包合同。插件声明创建、轮询、状态和结果映射；宿主统一处理任务恢复、出站校验、匿名结果下载、账务与资源持久化。结果下载不转发渠道 Authorization，即使结果地址与上游同源。原 Go Provider 仅保留为测试对照；资源持久化失败仍必须向上返回真实失败。
