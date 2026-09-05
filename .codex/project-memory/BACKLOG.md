# 待办与项目计划

更新时间：2026-09-04

## P0：接管与数据安全

- [x] 固定 `v1.2.5@f8e87bc` 并完成阶段 0–2：历史检查点与新 `origin/codex/dingxin-stable` 均为 `44e8c10`；完整对象审计为 66 提交/506 文件、89 重叠路径、模拟 47 显式冲突文件/149 区块。官方 `main@7ec5173` 的 6 个后续提交保持在外。
- [x] 在 `codex/upstream-v1.2.5-20260904` 完成阶段 3–4：真实 49 未合并路径/139 文本区块已按五组解决，未解决 0；当前 `HEAD=39c7e62`、`MERGE_HEAD=f8e87bc`、状态 `merge_open_resolved_unvalidated`。
- [ ] 获得新批准后进入阶段 5；先审计最终 staged bytes，再执行格式、Backend Linux CGO、迁移副本双启动、Web/Agent/Docs/Compose/Windows 脚本门禁。阶段 5 通过前不得创建 merge commit。
- [ ] 阶段 5 第一项必须移除 Create 对已删除 `backendModelRuntimeRequired` 的依赖，按官方 v1.2.5 统一走 `runBackendGenerationTask`，保留本地元数据/reasoning/任务恢复合同并补专项复核；修复前 TypeScript 预期无法通过。
- [ ] 在阶段 8/9 前为 Host Updater增加用户 fork Release/镜像源与版本策略，或显式保持不可用；当前默认 `ddcat-ai/open-ai-canvas` 不能用于自维护版本在线更新。
- [ ] v1.2.5 完成后再实施本地扩展边界收口，并在阶段 0–9 全部完成后单独评估/安装 Astravia Windows；当前不安装、不向其开放生产数据库、密钥或稳定 checkout。
- [x] 固定 `4ba9694` 已固化为本地 merge commit `d0def5f`，阶段 C通过；提交后评审记忆提交 `7687053` 已fast-forward推送fork并独立复核一致。本次最终检查点记忆继续以小型docs提交同步，完成后该集成分支作为已验证固定4ba远程检查点。没有force、未更新`origin/main`、未触发远程CI或部署。
- [ ] 后续独立事项：决定是否为fork集成分支建立可手工触发或PR触发的远程质量门禁；审计官方 `4ba9694..b7348ab` 22提交/216文件新尾差；另行规划阶段16部署到新本地提交的备份、迁移、镜像和Edge验收。
- [x] 路线B已完成升级前恢复点、回滚镜像标签、候选构建、克隆数据迁移/二次启动幂等和真实保留卷升级；非浏览器健康、日志、数据库和版本检查通过。
- [x] 2026-09-04 用户使用已登录Microsoft Edge完成人工只读冒烟并回传11张截图；版本/登录、稳定首页、导航、项目/章节、空资产分镜、素材/任务/Skill、后台及Console/Network基础健康通过，状态为 `runtime_candidate_edge_smoke_passed_limited`。
- [ ] 另行授权并准备真实项目资产后，补测分镜资产绑定/解绑、`@` mention、历史角色图/当前声音、音频引用、长名称溢出、刷新恢复和多镜头任务恢复；这不是本次空资产冒烟已覆盖范围。
- [ ] Provider真实生成、精确计费、OSS新写入、模型/价格保存、存储连接测试、Skill同步和删除操作继续保持独立门禁；不得因Edge只读冒烟通过而自动执行。

- [x] 对当前 SQLite 主库和 WAL 做一致性备份，并验证完整性、独立恢复副本和 OSS 密文解密。
- [x] 记录当前镜像 ID、容器、卷、数据目录和关键文件哈希。
- [ ] 收敛 `.claude/settings.local.json` 的宽泛授权。
- [x] 使用 Git binary patch 和 44 文件 ZIP 保护当前 9 个已修改文件和 35 个未跟踪文件，并验证 CRC、路径集合和 SHA-256。
- [x] 建立 `BACKUP_REGISTRY.md`，登记阶段 0 的非敏感基线、私有数据库快照、密钥、哈希、ACL、验证和保留策略。
- [x] 在主机执行环境复核 `gh auth status`/`gh api user` 并查询明确仓库；确认受限进程只是无法读取 Keyring，当前没有对应 PR。远程命名整理仍需在未来 Git 操作阶段单独评估。

## P1：Qwen 与能力合同

- [x] 建立唯一 `effectiveChannelModelCapability` 持久边界，由系统渠道 Catalog 与 Admission 共用；精确 SKU 计费和 Ready Provider 使用同一持久模型合同。
- [x] 明确定义并测试 missing、explicit-empty、wildcard 三态；显式空数组不回填，`*` 只开启自定义输入。
- [x] 按 Qwen、Wan、HappyHorse 具体模型建立保守能力矩阵：Qwen Image 3.0 Pro、Wan 2.7、HappyHorse 1.1、Wan 3.0 已完成；其他图片模型继续 Planned。
- [ ] 移除执行器不支持或会被静默改写的图片/视频参数：Ready 视频模型已完成；Planned/手工图片与通用 fallback 尚未全量收敛。
- [ ] 统一 Create、Canvas、Admin 和逻辑模型的尺寸派生纯函数。
- [x] Create 参数按用户、模型和生成方式持久化；刷新保持画幅/分辨率/时长/声音/水印，明确 Switch 与请求 payload 一致，并完成 Edge 验证。
- [x] 在管理后台和运行目录验证 Qwen Pro 的 auto、比例/像素、自定义尺寸、3 参考图、6 输出及不支持项；普通 Create 真实生成设置仍需模型定价/启用后的单独授权验证。
- [x] 删除 Qwen 临时 DEBUG、完整提示词/请求体/签名 URL/API Key 片段日志和未使用导入；Provider 只保留结构化安全状态。
- [x] 为官方发现目录增加 ready/planned/unsupported/deprecated 支持状态，并由 Backend 强制限制非 Ready 状态不可定价、启用、测试、删除或路由。
- [x] 重构 `FetchAdminChannelModels`，保留完整 CatalogItem 元数据并为现有待配置记录提供安全补齐；阶段 7 已在运行数据库完成拉取并验证真实幂等。
- [x] 完成 Wan 2.7、HappyHorse 1.1、Wan 3.0 的模型专属能力合同并晋升对应 Ready 模型；Prime、旧版和视频编辑继续由精确合同与账号证据控制。
- [x] 为 Wan 3.0 实现独立 All-in-One Adapter，支持 T2V、首帧、首尾帧、参考图片/视频/音频；file/link 保持未开放。

## P1：测试与可观察性

- [x] 复核 Canvas Agent Windows 进程树清理：确认 Codex 受限进程的 `taskkill` 被拒绝导致假阴性；测试夹具显式探测能力，主机真实专项 12/12、全量 287/287 和 TypeScript 构建通过。远程 Linux CI 仍随未来 PR 验证。
- [x] 用 Qwen 3.0 专属官方尺寸合同替代旧 38/39 项通用列表，并覆盖存量 Planned→Ready 安全补齐、字段缺失/空数组/`*` 与非法默认值边界。
- [x] 增加 Qwen `auto` 省略、`x -> *`、比例映射、1–6 输出、1–3 图编辑、参数依赖和拒绝路径完整序列化测试。
- [x] 为 Ready 视频模型增加 480P/时长/音频/水印、跨媒体数量、输入相关时长和非法组合拒绝测试。
- [x] 为系统渠道持久能力边界、目录/准入 fail-closed、Ready Provider payload、精确 SKU 账务和真实幂等增加测试。
- [x] 为任务失败建立前端、handler、service、worker、provider、资源物化全链路结构化日志，并以 request/task/provider/resource ID 关联。
- [x] 修复错误响应不记录日志或可能暴露内部错误的问题；统一稳定错误分类、retryable、requestId 和前端诊断编号。
- [x] 收敛生产 GORM 日志：忽略正常 `record not found`，保留真实错误元数据并整体隐藏 SQL 文本，覆盖 Raw SQL 预插值边界。

## P2：插件、部署与 CI

- [x] 修复插件备用端口、Vite 启动命令和可信 Origin 配置。
- [x] 让插件超时与 Agent 35 分钟生成续接合同一致，并保留 60 秒余量。
- [x] 用签名 Local Runtime 握手替代 `agentToken` 查询参数旧协议；Web 继续主动清除旧参数。
- [x] 修复不同 Compose 对 `ENABLE_PROVIDER_PLUGINS` 的透传差异，五个 Backend 入口默认均为 false。
- [ ] 在可调用 Codex CLI 的环境重新安装 `yingce@yingce-local`，新建线程验证插件 skill/MCP 加载和 Microsoft Edge 自动打开；当前 WindowsApps ACL 阻止 `codex plugin list`，不得伪报安装态通过。
- [x] 收紧 Server CORS：生产 Server Compose 必须显式提供非空 Origin，不再默认 `*`。
- [x] 在 CI 中增加 Canvas Agent 测试/构建，并为 Web 增加生产构建；远程 Actions 仍需提交/PR 后验证。
- [x] 消除 Bun test 内嵌 Vite/Rolldown 的退出竞争，增加 panic guard，并用 `.bun-version` 统一本机/CI；远程 Actions 仍需提交/PR 后验证。
- [ ] 补齐文档站缺页、路径和工具表。

## P2：百炼技能与协议资产

- [x] 逐篇索引本地百炼官方文档并保留 SHA-256、行数、标题和分类。
- [x] 建立全局 `bailian-model-api` 技能，按文本、图片、视频、语音、向量和平台 API 分流。
- [x] 建立逐模型同步/异步、端点、输入、限制、轮询、输出和临时 URL 的检索与证据合同。
- [x] 完成技能结构验证、项目/全局副本哈希核对和典型文本、图片、视频、语音、向量检索测试。
- [ ] 官方文档更新时重新生成索引并审阅差异。
- [ ] 跟进官方 Issue [#332](https://github.com/ddcat-ai/open-ai-canvas/issues/332) 的维护者回复；2026-08-28 实时复核仍为 open/评论 0/无标签/无指派，五项架构决策确认前，不提交支持状态 Schema、百炼 Adapter 或版本化 Manifest PR。
- [ ] Issue 提交后先等待 2–3 个工作日；若仍无维护者回复，再由用户明确授权后发送一次简洁的五项决策确认评论。禁止自动催办、重复评论或把点赞、标签、指派、交叉引用、关闭 Issue、普通贡献者回复和笼统“欢迎 PR”当成完整批准。
- [ ] 收到回复后核对回复者的 `author_association`；只有 `OWNER`、`MEMBER` 或 `COLLABORATOR` 的明确回答才作为维护者决策。逐项记录为 `confirmed/partial/alternative_requested/rejected/superseded`，部分回答继续保持未决项 NO-GO。
- [ ] Issue #332 获认可后，按“支持状态 → Qwen Image 3.0 Pro → Wan 3.0 → Wan 2.7 → HappyHorse 1.1 → 版本化 Manifest”顺序从官方最新 main 建立独立贡献分支。

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

## 阶段 16 本地交付状态

- [x] 最终 Backend、Web、Canvas Agent、插件和 Compose 全量门禁。
- [x] 最终数据完整性、OSS 加密、容器、镜像、数据卷和日志核对。
- [x] Microsoft Edge 首页、管理后台和创作页人工只读冒烟。
- [x] 创建 `BKP-20260827-160505-STAGE16-FINAL` 并登记最终本地镜像标签。
- [x] 生成 `STAGE16_FINAL_HANDOFF.md` 并清理项目记忆中的已修复旧结论。
- [ ] 集成分支 push 已完成；PR、远程 Actions、更新 `origin/main`、Git tag 和 Release 仍等待用户另行授权。

## 阶段 1 后续门禁

- [x] 在隔离 Linux CGO Backend target 中运行 SQLite 迁移和 Backend 全量测试；远程 GitHub Actions 仍待后续提交/PR 后确认。
- [x] 阶段 2 完成 Wan 3.0 All-in-One Adapter 与序列化/拒绝测试，标准版已从 `planned` 晋升为 `ready`；Prime 仍以精确官方合同和账号可用性证据为门禁。
- [x] 阶段 3 完成 Wan 2.7/HappyHorse 1.1 模型专属合同并晋升对应目录项。
- [x] 阶段 4 完成 Backend CGO、Web 和 Canvas Agent 全量测试与构建；阶段 6 以后才允许备份后的数据库迁移、容器重建和 Edge 运行验证。
- [x] 阶段 5 完成发布候选审计、能力/价格 fail-closed 收敛、敏感日志清理、独立候选镜像和迁移/回滚方案。
- [x] 阶段 6 创建 pre/post-deploy verified/protected 备份，保留卷逐个更新 Backend/Web，完成 Schema、数据计数、OSS 密文和运行哈希核对。
- [x] 阶段 7 已使用登录 Microsoft Edge 验证模型拉取、80 项补充目录、11 Ready/308 Planned、管理后台支持状态、重复拉取真实幂等，以及普通用户创作台不暴露未定价且停用模型；未点击生成。
- [x] 阶段 8 已完成一次授权内 Wan 3.0 真实调用、精确 SKU 账务、阿里云 OSS 物化、Resource/Asset、Edge 播放、无声/无水印、模型安全恢复和调用后 verified/protected 备份。实际画幅为 adaptive，未补发第二次任务。

## v1.2.5 阶段 5 后续

- [x] Create Backend-only 主链、重叠 Web 语义修复、官方公告图片自引用缺陷修复。
- [x] Web 类型/默认全套/补充专项/生产构建、Backend Linux CGO 全量、Compose 与 PowerShell 静态门禁。
- [x] 独立候选镜像、verified/protected 快照、两次无网络克隆迁移和计数/密钥验证。
- [x] Canvas Agent 以 Node 22 + Bun 1.4 Linux 隔离环境完成 328 项测试和构建；修正 unref heartbeat 测试的有界等待，0 fail、0 cancelled。
- [x] 阶段 5 最终 staged-bytes、冲突标记、敏感路径和版本复核通过；阶段 6 merge commit 仍需用户批准。
- [ ] Host Updater 完成 fork 仓库、镜像命名空间、版本策略和祖先校验前保持不可用。
- [ ] 最终 Microsoft Edge 验收对比本地画布底板/柔和网格与官方 v1.2.5 底板/80% 网格；由用户选择最终默认，当前实现不得被表述为永久产品决定。
- [x] 阶段 6 创建本地双父 merge commit `b6cff79`，固定纳入官方 `v1.2.5@f8e87bc`。
- [ ] 阶段 7 仅将候选分支推送到 fork 并独立核对远程 SHA；不得更新 `origin/main`。
