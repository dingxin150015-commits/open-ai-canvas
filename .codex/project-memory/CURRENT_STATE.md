# 当前状态

更新时间：2026-08-27

## Git 与远程

- 当前分支：`main`。
- 阶段 9 发布固化从 `0b202a1` 开始；Backend、Web、构建测试已分别形成 `64015db`、`62cf5d8`、`f97adde` 三个本地提交，文档与工程记忆由阶段 9 最终文档提交承载。
- 这些提交仅存在于本地 `main`，尚未 push、未创建 PR、未触发远程 GitHub Actions；本地验证结果不能替代远程 CI。
- 阶段 0 基线记录时工作树包含 9 个已修改文件和 35 个未跟踪文件，主要是 Qwen 能力补全、前端 fallback、调试日志、项目记忆、百炼技能和交接报告。
- `origin` 指向官方仓库 `ddcat-ai/open-ai-canvas`；`myfork` 指向用户 Fork；`upstream` 当前指向本机临时仓库，而非实时 GitHub。
- 用户在 2026-08-25 提供截图，确认交互式 PowerShell 中 `gh auth status` 已登录 `<authenticated-github-account>`，使用 HTTPS，具备 `gist/read:org/repo/workflow` scope。
- 2026-08-26 已复核“无法获取拉取请求状态”：默认受限进程不能读取 Windows Keyring，因此 `gh auth status` 会把凭据显示为无效；同一工作区在主机权限下 `gh auth status` 与 `gh api user` 均成功，账号仍为 `<authenticated-github-account>`。这是执行隔离，不是 GitHub 认证再次损坏，禁止为此重写或导出 Token。
- 官方仓库和用户 Fork 当前均无匹配的开放 PR，`ddcat-ai/open-ai-canvas` 上以 `<authenticated-github-account>:main` 为 head 的历史 PR 查询也为空。因此“无法获取 PR 状态”当前无需项目修复；后续查询必须在主机权限下显式传 `--repo`，无 PR 应显示为“未创建”，不能显示为认证失败。

## 运行与数据

- 项目全程使用 Microsoft Edge 进行开发联调；用户已确认 Edge 当前完成登录。后续管理后台和 UI 检查必须显式选择 Edge，禁止使用 Chrome 或自动回退其他浏览器。
- 2026-08-25 阶段 0 已完成一致性数据库备份与恢复验证。敏感快照位置：`<private-backup-root>\stage0-20260825-160625`；ACL 仅允许 `<current-windows-user>`、SYSTEM、Administrators。
- D 盘为 exFAT，不支持可靠 NTFS ACL；项目 `.local/backups/stage0-20260825-160625` 只保留 Git/容器非敏感基线和敏感备份位置指针，不保存数据库或设置密钥。
- 数据库快照 `integrity_check=ok`、外键违规 0、60 张表，独立恢复副本验证通过；配套 `.settings-key` 可成功解密 OSS 密文，未输出明文。
- 最近只读检查时，本地 Web 和 Backend 容器均为 healthy；这只证明健康入口可用。
- 当前 Web 容器曾返回 `index-ppkVtJ7Q.js` 和 `create-V_EnNwcp.js`，已不同于旧交接报告中的分块名。
- SQLite 主库、WAL 和存储迁移前快照存在；历史上发生过 `docker-compose down -v` 删除卷事故，并从 2026-08-17 备份恢复。
- 阶段 0 脱敏计数：1 个用户、1 个会话、1 个系统渠道、241 个渠道模型、0 个逻辑模型、0 个资源、0 个任务、0 个素材、0 个画布项目。
- 百炼渠道 241 个模型全部停用，241 个能力或协议待配置，视频能力模型为 0。
- 平台 OSS 配置已启用、字段完整、Secret 已加密保存且可用备份密钥解密；资源表为 0，因此仍无真实 OSS 上传成功证据。
- 当前数据库文件曾被设置为 `666`。这解决了短期写权限问题，但不符合生产最小权限要求；修复前必须先确认容器 UID/GID、卷权限和可恢复备份。

## Qwen 图片 UI

- `/create` 使用 `web/src/pages/create/index.tsx` 中的 `GenerationSettingsMenu`。
- `ImageSettingsPanel` 主要用于 Canvas 弹层。旧调试把日志加在错误组件上，因此“Create 页面无日志”不能证明缓存或组件未加载。
- 阶段 12 后，`qwen-image-3.0-pro` 使用独立 Qwen 3.0 能力：默认 `auto`、官方推荐比例/像素、自定义尺寸、最多 3 参考图/6 输出，不宣称质量、透明背景、mask、response/output format。

## 能力合同

- 系统渠道 Catalog 与任务 Admission 已共用 `effectiveChannelModelCapability` 持久能力边界；缺失、损坏或不完整配置 fail closed。
- Ready 的 Wan 2.7、HappyHorse 1.1 和 Wan 3.0 使用模型专属能力、Adapter 与 Provider payload 测试，不再依赖夸大的通用视频默认。
- `qwen-image-3.0-pro` 已晋升 Ready，使用模型专属同步 Provider；其他 DashScope 图片仍保持 Planned/通用隔离，不能继承 Qwen 合同。
- Planned/手工模型的通用能力 fallback 仍需谨慎；只有完成 `provider + protocol + modelKey` 合同与测试后才可晋升 Ready。

## 其他已确认问题

- `values` 缺失、显式空数组和通配符 `*` 的语义尚未统一。
- 任务失败响应路径可能不记录日志，因此“后端无日志”不能证明请求未发送。
- 插件备用端口流程把 Vite 写成 Next，且使用源码不读取的环境变量。
- 插件 MCP 超时为 90 秒，低于 Agent 长任务续接时长。
- 插件仍使用已被前端拒绝的 `agentUrl/agentToken` 查询参数旧协议。
- 文档站缺页、Compose 变量透传不一致和 Server CORS 默认值过宽；Canvas Agent CI 门禁已在阶段 4 补入工作流，但远程尚未触发。
- `.claude/settings.local.json` 历史上存在大量宽泛写入和破坏性授权，不应继承为 Codex 执行权限。

## 当前验证状态

- Backend 隔离 Linux CGO 全量测试、Web 449/449 + 跨 Runtime 1/1、TypeScript和生产构建均通过。
- 阶段 10 已推翻“Canvas Agent Windows 生产进程清理失败”的旧判断：4 个失败由 Codex 受限测试进程调用 `taskkill` 时被系统拒绝访问造成；同一源码在主机权限下真实 `taskkill /T /F` 专项 12/12、Canvas Agent 全量 287/287 和 TypeScript 构建均通过。
- 进程测试现会探测精确 Windows 进程树终止能力：受限环境用注入的直属测试子进程终止器继续验证取消、超时、输出上限和 receipt 语义，并把精确树能力明确标记为 skip；正常主机和 CI 仍执行真实生产终止器。
- 运行数据库已迁移并完成 319 条目录拉取；重复拉取真实幂等。Backend/Web 当前运行最终候选且 healthy，数据卷保持原位。
- Microsoft Edge 已验证管理后台支持状态、普通用户门禁、Create 音频/水印控制和 Wan 3.0 结果播放。
- 一次授权内 Wan 3.0 真实任务 succeeded，精确 SKU 账务 settled，阿里云 OSS Resource ready、Asset confirmed；随后模型恢复停用/未定价。
- 2026-08-25 的 Chrome/内置浏览器失效登录结果仍属于错误浏览器证据，后续只使用 Edge。
- 后续任何真实模型调用仍必须重新确认账号、模型、参数、额度、费用和重试次数；阶段 8 授权不可复用。

## 阶段 10 前空间清理

- 已删除可再生成的 `.local/cache` 大部分缓存、`web/dist`、`canvas-agent/dist` 和 `backend/server.exe`，按清理前精确计数至少释放 775,473,371 bytes（约 739.55 MiB）。
- `node_modules` 是阶段 10–16 持续开发依赖，未删除；阶段 0/6/7/8 的 verified/protected 恢复点和运行数据卷均未触碰。
- `.local/cache/stage1-go-mod` 初次清理因 Windows 文件属性和访问限制延迟完成；后续 `go clean -modcache` 已清除它，`.local/cache` 当前不存在，阶段 10 没有未删除或未停止的已确认目标。

## 阶段 11：Create 设置持久化与明确开关

- Create 设置现在按“用户作用域 + 用户选择的模型 + 生成方式”写入 IndexedDB；文生视频、图生视频、文生图和图片编辑不会互相覆盖设置。
- 页面不再在刷新或模型重算时无条件覆盖用户比例；设置草稿读取完成前禁止提交，读取失败才安全回退模型默认值。
- 生成声音和添加水印使用真实 Ant Design Switch；页面摘要、消息历史和任务请求显式使用当前 Switch 值，不再由旧全局默认静默决定。
- 阶段 11 验证：专项 13/13、Web 主套件 453/453、跨 Runtime 1/1、TypeScript 和生产构建通过；用户在已登录 Microsoft Edge 确认 16:9、480P、2 秒及无声/无水印刷新后保持，Switch 开关与摘要同步，未提交生成任务。
- 运行 Web 为 `sha256:b7cec94c5ce9c42a78a00a0129915dbd7f1ef8790cf480554d5bd793b60120f7`，容器 `25816ac62175` healthy、RestartCount=0；Backend 容器和镜像未变。旧 Web 镜像保留为 `open-ai-canvas-web:pre-stage11-20260826-2350`。

## 阶段 12：Qwen Image 3.0 Pro Ready

- 官方依据：`api-reference/image-generation/qwen-text-to-image.md` 与 `qwen-image-editing.md`；采用原生同步 `POST /api/v1/services/aigc/multimodal-generation/generation`。
- 文生图严格发送一条 user message/一个 text；图片编辑严格发送 1–3 个 image 后跟一个 text。`auto` 省略 `size`，内部 `x`/比例只在 Adapter 转成上游 `宽*高`。
- Provider 支持并校验 `n=1-6`、负面提示词 500 字符、prompt extend/direct-agent 依赖、thinking、watermark、seed；拒绝 mask、质量、透明背景、超量/超大/非法格式参考图和非法尺寸，不静默钳制或丢参。
- Catalog 版本 `2026-08-27`，80 项，SHA-256 `28D060E90D47B9DA943D6908330BCB82F352A95D47353FE0F10322A84F2B03B1`；Qwen Pro 从 Planned 晋升 Ready，其他图片模型保持 Planned。
- 隔离 Backend CGO 全量、Web 456/456、跨 Runtime 1/1、TypeScript 和生产构建通过。Microsoft Edge 首次拉取显示上游 242、官方 80、新增 1、补齐 69；第二次新增 0/补齐 0。
- 运行数据库为 320 模型：Ready 12、Planned 308；Qwen 能力 JSON 558 bytes，停用、未定价。未执行测试模型或真实图片生成，未产生费用。
- 当前 Backend `3b7b1e526a42` / `sha256:d0926f28a48cc4c6331860cf8006af040cb974caecc23a874eb180f81868b206`，Web `aac0f955c51d` / `sha256:3bf43cf9c603531f14c70d99cce0d7a4d8a93bc324d0f5c75d8c7a0c1420edfa`，均 healthy、RestartCount=0。

## 阶段 13：统一错误链路与安全日志

- Backend 业务失败响应保留既有数字 `code`，新增稳定的 `errorCode`、`errorCategory`、`retryable` 和 `requestId`；`X-Request-ID` 同时写入响应头。合法客户端请求编号会保留，含空格或可疑内容的编号会被替换。
- `ModelError` 现在可被 `errors.As` 正确投影；能力、价格、路由和 Provider 错误不再被误判为匿名 500。未分类 4xx/5xx、panic 和 JSON 解析错误只返回安全中文，不公开内部路径、SQL、网络错误或请求正文。
- 任务创建、重试、取消、Worker、Provider 和资源物化日志使用 request/task/provider/resource 关联键、分类、状态、耗时和错误类型；提示词、文本、content/input/output、签名 URL、API Key、Provider 原始 message 和非结构化任务错误均被隐藏。
- Web `ApiError` 读取 Backend 错误元数据；生成错误仅在请求编号符合安全字符合同后显示“诊断编号”，不能把伪造头或敏感字符串带入用户界面。
- 当前门禁：Backend 隔离 Linux CGO 全量测试通过；Web 主套件 458/458、跨 Runtime 1/1、TypeScript 与生产构建通过，11,022 个模块且无 Rolldown panic。
- 无费用运行验证：未认证任务请求返回 401/authentication；畸形登录 JSON 返回安全 400/validation；合法 request ID 端到端保留，非法 ID 被替换，Backend 日志不包含测试敏感标记。
- 当前 Backend `975ca3fa982b` / `sha256:71deebc35ad6ce81771c798dbcb45bdf68aa8d03c0092d11363993c1b4683910`，Web `5c5a8a6dd877` / `sha256:58a93110a7f7337f874d6863e417999f83b8c6949d573fccc2c3f4f343853e65`，均 healthy、RestartCount=0、OOMKilled=false。
- 本阶段没有 Schema/目录/价格变化，没有真实模型调用、OSS 上传或费用；远程 GitHub CI 仍未触发。GORM 的既有 `record not found` SQL 调试输出仍需在后续日志治理中单独收敛，不能视为阶段 13 的业务错误响应泄露。

## 百炼官方文档与全局技能

- 本地官方镜像 `<project-root>\百炼千问文档\百炼千问文档` 已逐篇索引：473 个文本文件、223,797 行。
- 每篇文档的标题、行数、模式标签和 SHA-256 位于 `.codex/skills/bailian-model-api/references/source-index.md`；机器可读目录位于 `source-catalog.json`。
- 项目内可版本控制的技能源：`.codex/skills/bailian-model-api`。
- 已安装的全局副本：`<codex-home>\skills\bailian-model-api`。
- 2026-08-25 校验结果：技能结构有效；项目源与全局副本 12/12 个文件 SHA-256 完全一致；文本、图片、视频、语音和向量检索脚本均返回官方文档结果。
- 技能只沉淀协议知识和检索流程，不授权真实 API 调用；任何可能计费的调用仍需用户单独批准。

## 百炼模型发现目录：阶段 1

- 2026-08-25 已在当前未提交工作树实现阶段 1：官方版本化 Manifest、`ready/planned/unsupported/deprecated` 支持状态、完整 CatalogItem 同步、既有未配置记录安全补齐和管理后台只读状态。
- Manifest 由 `scripts/build_bailian_catalog.py` 从本地官方文档生成，当前版本 `2026-08-25`，共 80 项；包含 78 个视频目录项和 2 个图片目录项，`wan3.0-video` 与 `wan3.0-video-prime` 均已进入发现目录。
- 阶段 1 中 80 项全部保持 `planned`。Backend 已在启用查询、保存/定价、模型测试、删除和逻辑模型线路处阻止非 `ready` 模型；不会制造“目录可见即执行器可用”的假合同。
- `FetchAdminChannelModels` 已改用完整 `FetchChannelModelCatalog`，保留 capability、protocol、provider model key、endpoint、operation、来源版本与文档证据。同步只补齐停用、未定价、无能力配置记录，使用 `updated_at` 条件防止并发管理员保存被覆盖；相同目录重复拉取不产生更新。
- 本阶段没有迁移当前数据库、没有重建/重启容器、没有操作 Edge、没有调用百炼或 OSS。当前运行中的 241 条模型记录仍保持阶段 0 状态，只有未来获批迁移/重建后才会应用新字段与目录。
- 验证：百炼/Service 合同 Go 测试通过；Database 包编译通过；前端专项测试 2/2 通过；`bun run build` 成功（11,020 模块，存在既有大 chunk 警告）；Manifest 重新生成前后 SHA-256 均为 `5276FC8B5304536B8FA2456CA63A3CD982D4BC03ABDC1D6323AB46E7035164E8`。
- 已解决：Windows 宿主机仍为 `CGO_ENABLED=0`，但新增隔离 Linux `backend-test` target，使用 `build-base + CGO_ENABLED=1` 完成当前源码 Backend 全量测试；SQLite 迁移、Database、Repository、Handler、Service 均通过。

## Wan 3.0：阶段 2

- 2026-08-25 已在当前未提交工作树实现 `wan3.0-video` 独立 All-in-One Adapter；标准版从 `planned` 晋升为 `ready`，`wan3.0-video-prime` 继续保持 `planned`。
- Adapter 使用官方原生异步合同：`POST /api/v1/services/aigc/video-generation/video-synthesis`、`X-DashScope-Async: enable`，复用现有 `GET /api/v1/tasks/{task_id}` 轮询和结果物化链。
- 已支持文生视频、单首帧、首尾帧、参考图片、参考视频和参考音频。官方 `file/link` 输入尚未进入产品素材合同，未列入 Ready 操作。
- 参数合同为 480P/720P/1080P、`adaptive/16:9/4:3/1:1/3:4/9:16`、2–30 秒或 `-1` 智能时长、音频开关、水印和可选 seed。默认值为 1080P、adaptive、5 秒、有声、无水印。
- Adapter 明确拒绝非法时长/分辨率/比例/seed、超过 20000 字符的提示词、帧与参考媒体混用、缺失帧、未物化 `asset://`、超量/超大/超时素材和固定时长下超过 30 秒的视频预算；不会静默钳制或截断。
- Wan 3.0 Ready 目录项会自动携带模型专属能力 JSON；存量未配置记录可在后续获批拉取时安全补齐，已配置记录仍不会被覆盖。
- 当前只完成 mock/离线验证；没有迁移运行数据库、重建容器、操作 Edge 或调用真实百炼。真实账号/区域/额度/计费与 OSS 结果仍未验证。
- 阶段 2 验证：Backend Manifest、目录补齐、能力投影、序列化、拒绝、终态和 mock 提交专项测试通过；Web 专项测试 6/6 通过；最终 `bun run build` 成功（11,020 模块，存在既有大 chunk 警告）；Manifest 重生成前后 SHA-256 均为 `44264D4A4C8B13697864380C638844DBC224E99003C9DAB74E5C10D96D86BDA1`。

## CGO 测试门禁与阶段 3

- `backend/Dockerfile` 已拆出共享 build base 和独立 `backend-test` target；`scripts/test-backend-cgo.ps1` 在 Windows 通过 Docker 运行 Linux CGO 全量测试，不安装全局 GCC、不挂载运行数据。
- `.github/workflows/quality.yml` 已显式设置 `CGO_ENABLED=1`，验证/按需安装 GCC，并禁用 Go test cache；远程 GitHub Actions 尚未实际触发，不能称 CI 已通过。
- 首次完整测试构建证明所有 Backend 包通过；阶段 3 最终源码在 gofmt 后再次全量测试，测试层耗时约 54 秒，Database/Repository/Handler/Service 等全部通过。测试前后现有容器 ID、镜像和状态完全一致。
- `.dockerignore` 排除 `.claude/.codex`、本地百炼镜像和 `web/dist`，Backend 测试上下文从约 56.6 MB 降至约 255 KB，同时不影响根 Dockerfile 使用 Web 源码。
- 阶段 3 为 Wan 2.7 T2V/I2V/R2V 和 HappyHorse 1.1 T2V/I2V/R2V 建立独立 Adapter 与模型能力。当前 Manifest 有 11 个 Ready 视频条目、67 个 Planned 视频条目、2 个 Planned 图片条目。
- Wan 2.7 Ready 覆盖 T2V 主线及两个快照、I2V 主线及快照、R2V 主线及快照；支持精确的音频 URL、帧角色、驱动音频、视频续写、1-5 视觉参考和参考音色合同。
- HappyHorse Ready 只覆盖 1.1 三模型：T2V、单首帧 I2V、1-9 图 R2V；1.0、视频编辑和其他模型继续 Planned。
- 新增跨层约束：图片+视频总视觉素材数、含参考视频时最大输出时长、驱动音频不得长于输出。Catalog 配置、逻辑模型投影、前端选择/提交、Backend admission 和 Adapter 均使用这些字段。
- 阶段 3 验证：Backend 专项与隔离 CGO 全量测试通过；Web 专项测试 9/9 和 TypeScript 检查通过；最终 `bun run build` 成功（11,020 模块，既有大 chunk 警告）；Manifest 重生成 SHA-256 为 `D52B4560F25C72C55C5566C6DB8980B60441E04CC2E648A87E8858DE6E0EFB62`。
- 为满足现有 CI 的全仓 `gofmt -l .` 门禁，另对 5 个历史 Go 文件执行了纯机械 gofmt；随后 `gofmt -l backend` 为空且隔离 CGO 全量测试再次通过。
- 本阶段仍未迁移运行数据库、重建服务容器、操作 Edge、调用真实百炼或验证 OSS 产物。

## 阶段 4：全量质量门禁

- 2026-08-26 Backend 最终源码通过隔离 Linux CGO 全量测试；测试镜像 manifest 为 `sha256:91be000c0997916ca4269a67e80830ae1812372274c3884cd66f1b917b61659d`。测试 target 未挂载运行数据库、备份、密钥或数据卷。
- Web 全量测试通过：主套件 445/445、跨 Runtime 套件 1/1；TypeScript 与 Vite 生产构建通过，11,020 个模块。原 Bun/Rolldown worker panic 已修复：`local-channel-runtime-projection.node.test.mjs` 不再在 Bun test 主进程内启动 Vite，而是直接导入业务模块；专项连续 10 次 0 失败/0 panic，最终全量测试无 panic，生产构建约 3 分 35 秒通过。
- 新增 `web/scripts/run-with-panic-guard.mjs`，即使子测试错误地以 0 退出，只要输出 Rolldown/Rust ModuleLoader panic 就以 86 失败；合成 clean/panic fixture 分别验证退出码 0/86。仓库根 `.bun-version` 固定 Bun 1.4.0，GitHub Actions 两个 Bun job 改为读取该唯一版本源。
- Canvas Agent 最终全量测试 286/286 通过，随后 `npm run build` 通过。Windows 本机没有创建文件符号链接的权限，因此该单一安全子项明确交给 Linux CI；本机仍执行并通过遍历、根逃逸、硬链接、超限、文件头、TOCTOU 与路径交换检查。
- 修复了 Windows/exFAT 下小于一次文件锁往返的测试租约、依赖环境定时器先后顺序的脆弱断言，以及原子替换瞬间 `ENOENT` 的轮询测试。生产默认 30 秒租约未改动。
- 已撤销未被源码引用的 `modernc sqlite` 依赖族及伴随版本漂移，`backend/go.mod`、`backend/go.sum` 精确回到 Git 基线；生产继续使用 `go-sqlite3 + CGO`。
- CI 已增加 Web 生产构建和独立 Canvas Agent 测试/构建 job。远程 GitHub Actions 尚未触发，不能称远程 CI 已通过。
- 文档目录仍没有 `package.json`、`source.config.ts` 或完整页面集合，README/AGENTS 中声明的 Next/Fumadocs 构建当前不可执行；这是既有 P2 缺口，不在本阶段伪造“文档构建通过”。
- 阶段 4 仍未迁移数据库、重建运行服务、操作 Edge、调用真实百炼/OSS 或产生模型费用。

## 阶段 5：发布候选与迁移准备

- 发布前审计已收敛旧 Qwen DEBUG、图片尺寸三态、Catalog/Admission 能力分裂、未配置价格误路由和 DashScope 日志敏感信息风险；对应 Web/Backend 测试均已补充。
- 系统渠道目录和任务准入现共用持久能力配置；缺失、损坏或不完整时 fail closed。旧记录只有能力、协议和有效能力 JSON 同时存在才回填 `ready/manual`。
- Manifest 重生成仍为 80 项、11 Ready/67 Planned 视频和 2 Planned 图片，SHA-256 仍为 `D52B4560F25C72C55C5566C6DB8980B60441E04CC2E648A87E8858DE6E0EFB62`。
- Web 最终门禁为主套件 446/446、跨 Runtime 1/1、TypeScript 与生产构建通过且无 panic；Backend 隔离 Linux CGO 全量测试通过。
- Canvas Agent 使用统一后的 Bun 1.4.0 入口再次完成 286/286 全量测试与 TypeScript 构建。
- 已构建独立候选镜像：Backend `sha256:0e5744fdfad8767320dadf66a3fb510dd3decbdd335f071b2e64d1dd43200f91`，Web `sha256:0a18498c9634c82fc47d327da8c67ef8e6b62e2e3506e62070c1a53e7e7b0bc6`。临时候选容器健康/配置检查通过。
- 当前运行容器仍为旧 `:local` 镜像，ID `583a381f2cc3`（Web）和 `b8415a78e587`（Backend），保持 healthy；数据库和卷未触碰。
- 阶段 6 前必须按 `STAGE5_RELEASE_CANDIDATE.md` 创建新的 pre-deploy verified 备份；当前仍不能声称运行后台已应用新目录。

## 阶段 6：备份、迁移与运行镜像更新

- 已创建并登记两份私有恢复点：`BKP-20260826-130230-STAGE6-PREDEPLOY` 和 `BKP-20260826-131009-STAGE6-POSTDEPLOY`，均为 `verified/protected`；主快照、独立恢复副本、外键和 OSS 密文解密验证全部通过。
- 旧镜像由 `open-ai-canvas-backend:pre-stage6-20260826-130230` 和 `open-ai-canvas-web:pre-stage6-20260826-130230` 标签保护；数据卷未删除、未替换。
- 当前 Backend 容器 `5693d185fab6` 运行候选镜像 `sha256:0e5744fdfad8767320dadf66a3fb510dd3decbdd335f071b2e64d1dd43200f91`；Web 容器 `f278caf9beaa` 运行 `sha256:0a18498c9634c82fc47d327da8c67ef8e6b62e2e3506e62070c1a53e7e7b0bc6`。两者 healthy、Paused=false、RestartCount=0、OOMKilled=false。
- Schema 迁移后仍为 60 张表，完整性 `ok`、外键违规 0；ChannelModel 目录列全部存在。241 条旧记录全部为 `planned/legacy`，Provider Key、操作 JSON 和文档 JSON 缺失均为 0。
- 用户、会话、渠道、模型、资源、任务、素材和画布项目计数与部署前一致；OSS 1 条配置/1 个 Secret 仍全部加密且可解密。
- Web 首页响应与容器 `index.html` SHA-256 一致：`29B1107BDC54C935A69AF9549A17BD8297513E4C6C01452D20E594326326FE6C`；Nginx 代理 `/api/health` 返回业务 code 0。
- 本阶段没有点击模型拉取，因此运行数据库仍是 241 条 legacy/planned；80 项官方补充目录和 11 个 Ready 模型要到阶段 7 通过已登录 Microsoft Edge 执行无费用拉取后才能进入数据库。
- 未进行真实百炼生成或 OSS 上传；这些仍需要后续独立授权与验证。

## 阶段 7：运行目录拉取与 Edge 管理后台验证

- 用户在已登录 Microsoft Edge 中完成百炼模型拉取。运行数据库现有 319 条渠道模型：`ready=11`、`planned=308`；官方补充 80 项，其中图片 2、视频 78。
- 11 个 Ready 均为已完成 Adapter 与合同测试的 Wan 2.7、HappyHorse 1.1 和 `wan3.0-video`；全部保持未定价、停用。非 Ready 误启用、误定价和 Ready 合同缺失均为 0。
- Edge 管理后台已直接验证：支持状态筛选“可用”共 11 条；Planned 条目只显示“查看”；搜索 `wan3.0` 得到 `wan3.0-video`，能力为视频、协议为 DashScope 视频、支持状态为可用、价格未配置、状态停用。
- 首轮重复拉取暴露真实幂等缺陷：上游 239 条没有能力/协议时，同步仍提交 `capability=""`、`protocol=""` 更新，导致 toast 长期显示“补齐 239”并改变 `updated_at`。
- `channel_model_catalog_sync.go` 现只在目标 capability/protocol 非空时补齐；新增 `TestUpstreamCatalogEnrichmentIsIdempotent` 和按字段聚合的安全更新计数。修复后重复拉取显示“新增 0，补齐 0”，目录摘要前后均为 `997BF93019C1BA6F1E34191BDFCBAF0AA76ABFFF3627370FD9F3F4908C97B7C4`。
- 当前 Backend 容器 `5293ca038048` 使用镜像 `sha256:6faf8c852dfcab4289afc2f87ea0deb2e2c0ff7927a010a9fa2a8fe40c35023e`；Web 容器仍为 `f278caf9beaa`。两者 healthy，数据卷仍为 `open-ai-canvas_backend-data`。
- 已创建并登记拉取后恢复点 `BKP-20260826-171800-STAGE7-CATALOG`；数据库、独立恢复副本、外键、OSS 密文、目录合同、ACL 和哈希全部验证通过。
- 普通用户侧门禁已由 Edge 截图确认：视频生成模型选择器为空并提示“当前没有可用模型，请联系管理员或检查模型配置”；因此尚未定价且停用的 `wan3.0-video` 未被误开放。用户未点击生成，任务表仍为 0。
- 新增 `TestPublicSystemChannelCatalogExcludesDisabledAndMarksUnpricedUnavailable`，锁定停用模型不进入普通用户目录、启用但未定价模型仍为 unavailable。隔离 Linux CGO Backend 全量测试通过，测试镜像 manifest 为 `sha256:a18a21f61b6ee41d73844fef2969aff4ed7bc906dd578c0cfcb0003ddd78cf5e`。
- 阶段 7 已正式收口；下一阶段涉及真实百炼生成和 OSS 结果物化，可能产生费用，必须先取得新的明确授权。
- 未进行真实百炼生成或 OSS 上传；任何可能计费的阶段 8 调用仍需单独批准。

## 阶段 8：真实调用前审计与候选部署

- 官方当前合同：`wan3.0-video` 使用原生异步 `POST /api/v1/services/aigc/video-generation/video-synthesis`，轮询 `GET /api/v1/tasks/{task_id}`；最短 2 秒，最低 480P，结果 URL 和 task ID 24 小时有效。官方模型市场当前显示北京 480P 限时价 0.21 元/秒、原价 0.30 元/秒；账单仍以账号实际结算为准。
- 脱敏运行配置确认：百炼渠道是北京 MaaS 主机的 `/compatible-mode/v1`，渠道启用且 API Key 存在；`wan3.0-video` Ready，但仍未定价、停用。平台 OSS 已启用，提供商为阿里云、Region 为 `oss-cn-wulanchabu`，Endpoint/Bucket/AccessKey/加密 Secret 均存在。
- 前置审计发现并修复三项真实阻断：原生视频路径曾错误追加在 `/compatible-mode/v1` 后；结果 OSS URL 下载曾携带模型 API Key/自定义头；`generation-image/video/audio:*` 未被远端同步识别，导致生成素材不能进入平台 OSS。
- 原生端点现在保留账号/区域绑定主机，只替换协议路径；结果下载验证安全 URL 且不再转发渠道凭据；生成素材按 image/media 存储源读取，上传失败会 fail closed 并保留重试，不再把本地 blob 地址伪装成服务端成功。
- 验证：Backend 隔离 Linux CGO 全量测试通过；Web 主套件 448/448、跨 Runtime 1/1，通过 TypeScript 和生产构建；候选 Backend/Web 临时容器 healthy。
- 当前已部署候选：Backend `c33b097bd7a2` / `sha256:5b236e5e...`，Web `b27398fc1af4` / `sha256:fa6f2eaa...`，均 healthy、RestartCount=0；旧镜像由 `pre-stage8-preflight-20260826` 标签保护，数据卷未改变。
- 已创建 `BKP-20260826-185248-STAGE8-PRECALL` 并验证完整性、独立恢复、OSS 密文、ACL 与哈希。调用前仍为任务 0、资源 0、素材 0、启用模型 0、定价模型 0。
- 2026-08-26 用户已明确批准一次 `wan3.0-video` 真实调用：北京 MaaS、480P、16:9、2 秒、无声、无水印、seed=`20260826`、创建重试 0，供应商费用上限 0.60 元。授权不包含第二次创建、其他模型/参数或提高费用；执行和结果验证进行中。
- 调用前又修复系统渠道精确 SKU 计费：任务账务现在使用生成方式、分辨率和时长匹配价格档，不再用空规格查价。隔离 Linux CGO Backend 全量测试通过；当前 Backend 为 `f97220ae4273` / `sha256:0d3b2403dbcf106f8dd61dcb3cd88692e2858c17a591b2aeec8d3aca6a08cdca`，healthy、RestartCount=0。
- 工具链事故：首次运行卷监控容器时虽以 SQLite readonly 打开，但容器用户为 root，导致 WAL/SHM 变为 `root:root 0644`，Backend `app` 用户启动回填时报 readonly 并重启。付费任务尚未创建。已停止重启循环，仅将 DB/WAL/SHM 恢复为 UID/GID `100:101`、`0660`，回滚容器恢复 healthy；随后以相同 UID/GID 监控验证任务/资源/素材/账务仍为 0，文件权限保持不变，再成功部署精确计费候选。
- 新门禁：任何直接挂载运行 SQLite 卷的诊断容器必须使用 Backend 相同 UID/GID；`:ro` 卷和 SQLite readonly 连接也可能创建/修改 SHM，不能把它们等同于文件系统零写入。
- 用户已进一步明确批准本次调用省略 seed；其余合同仍锁定为 `wan3.0-video`、480P、16:9、2 秒、无声、无水印、创建重试 0、供应商费用上限 0.60 元。当前进入 Edge 临时价格配置与单次提交阶段。
- Edge 已保存临时价格配置，数据库复核通过：模型 enabled/priceConfigured=true；唯一价格档 `PTIER_000001` 精确 selector 为 `text_to_video + 480p + 2秒`，按秒 1 积分，价格档启用。任务/资源/素材/账务/API 调用仍全部为 0，余额 100、预留 0，Backend healthy。
- 真实请求前 Edge 截图发现 Create 的 `GenerationSettingsMenu` 只显示画幅/分辨率/时长，没有渲染能力中已声明的生成音频和水印；全局旧默认会发送音频=true，若直接提交将违反无声授权。
- 已在真实 `/create` 菜单增加能力驱动的“生成声音/添加水印”输出开关和摘要，不支持的模型不显示；专项、TypeScript、Web 449/449、跨 Runtime 1/1 和容器生产构建通过。当前 Web `7b71d615ca98` / `sha256:3900fbb49fceda5da019be3991b818d36c2e4f6831c688f4d1adaaeab1e7daa2`，Backend `f97220ae4273`，均 healthy、RestartCount=0。真实任务仍未提交。
- 用户在看到“无声·无水印”摘要后提交了唯一任务；实际画幅在刷新后回到 `adaptive`，不是已批准的 16:9，但 480P/2 秒和费用上限未变化。未授权且不会补发第二次任务。
- 唯一真实任务 `a434…d0ab` 已 succeeded，attempts=1；API create 仅 1 条且 HTTP 200。供应商实际请求为 `adaptive`、480P、2 秒、audio=false、watermark=false、无 seed。
- 内部账务已 settled：数量 2 秒、实际 2 积分、预留归零、余额 98。阿里云 OSS Resource `f6da…e3b0` 为 ready、video/mp4、805,526 bytes、Endpoint `oss-cn-wulanchabu.aliyuncs.com`；Asset 已 confirmed 并使用 `resource:` 存储键。服务日志记录资源文件 Range 读取 206，证明浏览器已通过资源接口读取对象。等待用户最终确认实际播放、无声和无水印。
- Edge 最终视觉验证通过：视频可正常播放，播放器显示 0:02；用户确认无声音、无水印，截图无可见水印。下一步仅恢复模型停用/未定价并创建调用后恢复点，不删除任务、账务、Resource、Asset 或 OSS 对象。
- 阶段 8 已完成：Wan 3.0 模型和唯一临时价格档已恢复 disabled/unpriced、单价 0；成功任务、settled 账务、ready OSS Resource 和 confirmed Asset 保留。调用后 `BKP-20260826-201444-STAGE8-POSTCALL` 已完成完整性、外键、独立恢复、OSS 密文、运行审计、ACL 和哈希验证。
