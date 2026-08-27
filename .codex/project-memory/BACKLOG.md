# 待办与项目计划

更新时间：2026-08-26

## P0：接管与数据安全

- [x] 对当前 SQLite 主库和 WAL 做一致性备份，并验证完整性、独立恢复副本和 OSS 密文解密。
- [x] 记录当前镜像 ID、容器、卷、数据目录和关键文件哈希。
- [ ] 收敛 `.claude/settings.local.json` 的宽泛授权。
- [x] 使用 Git binary patch 和 44 文件 ZIP 保护当前 9 个已修改文件和 35 个未跟踪文件，并验证 CRC、路径集合和 SHA-256。
- [x] 建立 `BACKUP_REGISTRY.md`，登记阶段 0 的非敏感基线、私有数据库快照、密钥、哈希、ACL、验证和保留策略。
- [x] 在主机执行环境复核 `gh auth status`/`gh api user` 并查询明确仓库；确认受限进程只是无法读取 Keyring，当前没有对应 PR。远程命名整理仍需在未来 Git 操作阶段单独评估。

## P1：Qwen 与能力合同

- [x] 建立唯一 `effectiveChannelModelCapability` 持久边界，由系统渠道 Catalog 与 Admission 共用；精确 SKU 计费和 Ready Provider 使用同一持久模型合同。
- [x] 明确定义并测试 missing、explicit-empty、wildcard 三态；显式空数组不回填，`*` 只开启自定义输入。
- [ ] 按 Qwen、Wan、HappyHorse 具体模型建立保守能力矩阵：Wan 2.7、HappyHorse 1.1、Wan 3.0 已完成；Qwen 图片仍 Planned。
- [ ] 移除执行器不支持或会被静默改写的图片/视频参数：Ready 视频模型已完成；Planned/手工图片与通用 fallback 尚未全量收敛。
- [ ] 统一 Create、Canvas、Admin 和逻辑模型的尺寸派生纯函数。
- [x] Create 参数按用户、模型和生成方式持久化；刷新保持画幅/分辨率/时长/声音/水印，明确 Switch 与请求 payload 一致，并完成 Edge 验证。
- [ ] 在真实 `/create` 设置菜单上验证 Qwen 比例、分辨率、质量和自定义尺寸。
- [x] 删除 Qwen 临时 DEBUG、完整提示词/请求体/签名 URL/API Key 片段日志和未使用导入；Provider 只保留结构化安全状态。
- [x] 为官方发现目录增加 ready/planned/unsupported/deprecated 支持状态，并由 Backend 强制限制非 Ready 状态不可定价、启用、测试、删除或路由。
- [x] 重构 `FetchAdminChannelModels`，保留完整 CatalogItem 元数据并为现有待配置记录提供安全补齐；阶段 7 已在运行数据库完成拉取并验证真实幂等。
- [x] 完成 Wan 2.7、HappyHorse 1.1、Wan 3.0 的模型专属能力合同并晋升对应 Ready 模型；Prime、旧版和视频编辑继续由精确合同与账号证据控制。
- [x] 为 Wan 3.0 实现独立 All-in-One Adapter，支持 T2V、首帧、首尾帧、参考图片/视频/音频；file/link 保持未开放。

## P1：测试与可观察性

- [x] 复核 Canvas Agent Windows 进程树清理：确认 Codex 受限进程的 `taskkill` 被拒绝导致假阴性；测试夹具显式探测能力，主机真实专项 12/12、全量 287/287 和 TypeScript 构建通过。远程 Linux CI 仍随未来 PR 验证。
- [ ] 增加 Qwen 38/39 项、同名 SKU 和真实存量配置测试；字段缺失、空数组、`*` 和非法默认值已有局部覆盖。
- [ ] 增加 Qwen `auto` 省略、`x -> *`、比例映射和 1–6 输出完整序列化测试；当前只有尺寸规范化局部覆盖。
- [x] 为 Ready 视频模型增加 480P/时长/音频/水印、跨媒体数量、输入相关时长和非法组合拒绝测试。
- [x] 为系统渠道持久能力边界、目录/准入 fail-closed、Ready Provider payload、精确 SKU 账务和真实幂等增加测试。
- [ ] 为任务失败建立前端、handler、service、worker、provider 全链路结构化日志。
- [ ] 修复错误响应不记录日志或可能暴露内部错误的问题。

## P2：插件、部署与 CI

- [ ] 修复插件备用端口、Vite 启动命令和可信 Origin 配置。
- [ ] 让插件超时与 Agent 长任务续接合同一致。
- [ ] 用当前本地运行时握手替代 `agentToken` 查询参数旧协议。
- [ ] 修复不同 Compose 对 `ENABLE_PROVIDER_PLUGINS` 的透传差异。
- [ ] 收紧 Server CORS 默认值。
- [x] 在 CI 中增加 Canvas Agent 测试/构建，并为 Web 增加生产构建；远程 Actions 仍需提交/PR 后验证。
- [x] 消除 Bun test 内嵌 Vite/Rolldown 的退出竞争，增加 panic guard，并用 `.bun-version` 统一本机/CI；远程 Actions 仍需提交/PR 后验证。
- [ ] 补齐文档站缺页、路径和工具表。

## P2：百炼技能与协议资产

- [x] 逐篇索引本地百炼官方文档并保留 SHA-256、行数、标题和分类。
- [x] 建立全局 `bailian-model-api` 技能，按文本、图片、视频、语音、向量和平台 API 分流。
- [x] 建立逐模型同步/异步、端点、输入、限制、轮询、输出和临时 URL 的检索与证据合同。
- [x] 完成技能结构验证、项目/全局副本哈希核对和典型文本、图片、视频、语音、向量检索测试。
- [ ] 官方文档更新时重新生成索引并审阅差异。

## 执行门禁

任何实施阶段都按以下顺序：

1. 最小专项单元测试；
2. Backend 全量测试；
3. Web 测试和构建；
4. Canvas Agent 测试和构建；
5. 数据备份；
6. 保留卷重建并重新创建容器；
7. 源码、镜像、容器和浏览器哈希核对；
8. Mock/无费用浏览器端到端验证；
9. 用户另行批准后才进行真实付费模型调用。

## 阶段 1 后续门禁

- [x] 在隔离 Linux CGO Backend target 中运行 SQLite 迁移和 Backend 全量测试；远程 GitHub Actions 仍待后续提交/PR 后确认。
- [x] 阶段 2 完成 Wan 3.0 All-in-One Adapter 与序列化/拒绝测试，标准版已从 `planned` 晋升为 `ready`；Prime 仍以精确官方合同和账号可用性证据为门禁。
- [x] 阶段 3 完成 Wan 2.7/HappyHorse 1.1 模型专属合同并晋升对应目录项。
- [x] 阶段 4 完成 Backend CGO、Web 和 Canvas Agent 全量测试与构建；阶段 6 以后才允许备份后的数据库迁移、容器重建和 Edge 运行验证。
- [x] 阶段 5 完成发布候选审计、能力/价格 fail-closed 收敛、敏感日志清理、独立候选镜像和迁移/回滚方案。
- [x] 阶段 6 创建 pre/post-deploy verified/protected 备份，保留卷逐个更新 Backend/Web，完成 Schema、数据计数、OSS 密文和运行哈希核对。
- [x] 阶段 7 已使用登录 Microsoft Edge 验证模型拉取、80 项补充目录、11 Ready/308 Planned、管理后台支持状态、重复拉取真实幂等，以及普通用户创作台不暴露未定价且停用模型；未点击生成。
- [x] 阶段 8 已完成一次授权内 Wan 3.0 真实调用、精确 SKU 账务、阿里云 OSS 物化、Resource/Asset、Edge 播放、无声/无水印、模型安全恢复和调用后 verified/protected 备份。实际画幅为 adaptive，未补发第二次任务。
