# 百炼声明式插件完整迁移

最终状态（2026-09-12）：`bailian_declarative_deployed_preverified_edge_pending`。正式 Backend 已部署，非浏览器预验证完成；下一步仅由用户手工 Microsoft Edge 确认。

## 基线与授权

- 基线为已晋升的 `d9fec0f4b6f2cec48ca11c8f69bd5c4e7c8b2f48`，修复分支 `codex/dashscope-system-plugins-20260911`，独立工作树 `.local/worktrees/dashscope-plugins`。
- 用户明确批准完整迁移、必要的通用引擎扩展、系统性预验证和部署，完成后停在 Microsoft Edge 门禁前。没有批准真实付费生成或远程发布。
- 主工作区仍在旧 `b6254b41`；此前宿主直接注册补丁及其测试不能作为当前基线的发布证据。受控宿主插件过渡实现已从新分支移除。

## 当前实现

- 新增正式 `dashscope-image` / `dashscope-video` 插件，使用统一官方包 bootstrap、Registry、目录及启停流程，未放宽上传包 host 绑定限制。
- 模型线协议迁入 `plugin-packages/bailian-specs.mjs`；生产不再调用旧 DashScope Provider。旧文件改为 `_legacy_test.go`，仅做原始请求对照。
- 引擎增加纯变量、严格数值/布尔解析、通用尺寸及媒体规则、帧选择、数组与字符汇总、区域路径后缀处理、严格响应状态和匿名结果下载。
- 插件整包 Provider 注册具有原子性；停用影响目录、报价和新任务准入。已知新模型不允许保存为旧 Wanx/Wan 协议。
- 当前产品能力画像保留 UI/报价/准入职责；完全插件化产品能力目录属于后续独立范围，不能把线协议迁移等同于删除全部模型知识。

## 预验证进度

- 已完成 12 个 Ready 模型及关键操作的请求对照夹具、非法输入、启停/重启、保存读回、同步/异步本地 HTTP 模拟与恢复不重复创建测试。
- 首轮问题：重复变量导致包加载失败；缺省空字符串被 coalesce 变为 nil；视频续写素材顺序不同；通用下载同源带鉴权。均已定位并修正，加入对照或安全断言。
- 首轮 Linux CGO 全量只有既有 `TestDeleteAdminResourcesSupportsPartialSuccessAndBatchedReferences` 出现 SQLite 表锁；本次新增测试通过。最终代码仍须重新全量验证，不提前宣称最终通过。
- 第二轮全量的唯一失败为既有 `TestExpiredArchivedAssetCleanupRespectsCanvasReferencesAndUsesDeletionOutbox` 的同类内存 SQLite 表锁。已将两处资源删除测试夹具限定为单连接，保留全部断言及生产连接池不变；后续全量通过。

## 最终验证与部署

- 源码提交 `d6fed73033ebc4bd8cd508f45e0b16f5d8f862e5`；收敛无关格式后的部署提交 `8dcdfd5dd1fda78b0d000cfd98d143c73fc4b6f3`。修复分支干净检查后构建，未推送或晋升稳定分支。
- 最终 `scripts/test-backend-cgo.ps1` 在相同提交上运行 `go test -count=1 ./...` 全量通过，service 200.199 秒；完整日志在修复树 `.local/backend-8dcdfd5d-full.log`，没有 FAIL、panic 或构建 ERROR。
- 两包使用 `yingce.plugin/v3` / `3.0.0`；新语义不会被旧引擎静默忽略。后台视图仍使用既有归一化展示合同。
- 镜像 `open-ai-canvas-backend:bailian-plugins-20260911` / `sha256:8102a7c11dd4d30897c86b6a9e7245ce7e533031809f46dff9f3df17792f3609`；运行容器 `45b8cc449d3c`，版本 `v1.2.7+dingxin.1.bailian`，build commit 为 `8dcdfd5d`，healthy / restart 0 / OOM false / schema 8/8。
- 候选通过断网克隆卷双启动，模型、价格、任务、资源、系统设置、密文可解密性和迁移摘要均保持；正式保留卷部署前后也完成相同比对，完整性 ok、外键违规 0、密文解密失败 0、80 张表。
- 真卷仍为 `open-ai-canvas_backend-data`；运行环境参数逐项一致；Web 容器 `f965eb451aef` 未替换。入口 HTTP 200；未登录协议目录请求 HTTP 401；启动错误筛选 0。
- 运行注册表显示两个官方插件启用且版本为 3.0.0，原插件启停状态保持。模型记录、价格、启用状态均未改写；已保存错误旧协议的记录不会被静默迁移，需在 Edge 中核对。
- 镜像内与源码包 SHA-256 一致：image `6475add3a9f6fd12d32a038523baf9ab9092c998eafb8ee658ffd3272848e430`；video `63fb3532c8091282d6bfb8ec1c021fa50f69d5cc73bfec371809769eae08863c`。
- 新恢复点：`BKP-20260912-172928-BAILIAN-PREDEPLOY` 与 `BKP-20260912-173253-BAILIAN-POSTDEPLOY`，均 verified/protected；私有位置见本机 `.local/deployment-state.json`。部署前镜像由 `open-ai-canvas-backend:pre-bailian-20260912-172928` 保护。原 2026-09-11 恢复点也保留。
- 专用断网克隆容器和卷 `open-ai-canvas-bailian-preflight-20260911` 已在核对身份后删除；可从私有备份重建，真实卷与备份未删除。
- 旧主工作区的三个宿主注入代码改动已归档后撤销；入口记忆已修正为正确工作树和基线，历史记录明确标注被替代。

## Edge 人工门禁

1. 使用已有登录态 Edge 刷新页面，打开 `/admin/plugins`，确认两个百炼协议插件为官方来源、版本 3.0.0，且没有加载错误。
2. 打开 `/admin/channels` 的百炼模型编辑器：Qwen Pro 应选择 `DashScope Qwen Image`，路径为 `/api/v1/services/aigc/multimodal-generation/generation`；视频应选择 `DashScope Wan / HappyHorse Video`，创建与任务查询路径正确。
3. 确认原模型价格和启用状态保持，已有参数未被重置；旧 Wanx/Wan 选项不能作为上述 Ready 模型的替代协议。
4. 本批只做 UI/配置入口确认，不点击“测试模型”或生成。真实上游权限、额度和生成质量仍需独立授权验证。

本次没有操作 Edge、发起真实 Provider 生成、推送 GitHub 或触发远程 CI。以上为源码/离线夹具/运行数据与非浏览器入口验证，不等于 Edge 验收或真实付费生成通过。
