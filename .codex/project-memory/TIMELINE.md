# 开发历史时间线

## 2026-08-17：模型渠道插件体系

- `97223bc`：管理后台系统渠道支持 API 格式选择。
- 建立 Provider 插件、系统渠道和前端协议选择相关基础。

## 2026-08-19：DashScope 图片验证阶段

- `52b2600`：移除 DashScope 图片调试日志。
- 历史材料记录过图片生成、OSS 结果和六模型尝试；这是当时分支的历史证据，合并后的当前主分支尚未重新完成真实调用回归。

## 2026-08-20：多参考图

- `11931d0`：DashScope 支持多张参考图生成。
- 形成多参考图格式、上传、模型限制和适配器转换经验。

## 2026-08-23：DashScope 图片与视频实现

- `a903485`：完整实现 DashScope 图片和视频协议。
- 覆盖同步图片、异步视频、任务轮询、媒体输入和模型发现等方向。
- 历史文档中的接口样例与当前 Provider 已发生漂移，不能直接复制旧请求。

## 2026-08-23：上游历史获取与大规模合并

- 通过本机临时上游仓库获取较完整历史。
- `c4eabc5`：合并上游 main，共 455 个新提交。
- `c30a685`：修复合并后的编译问题。
- 合并阶段曾使用批量 checkout，静默覆盖至少 4 处本地功能，成为后续重要教训。

## 2026-08-24：补回被覆盖功能并完成主分支集成

- `e95050f`：修复合并遗留缺陷并补回被覆盖的本地改动。
- `73faab2`、`ae65991`：补充合并报告和经验总结。
- `0b202a1`：合并进入当前 main。
- 历史报告声称 Go/前端构建和部分测试通过，但仍存在失败项，且没有完成当前版本的 DashScope、管理后台、UI 和端到端回归。

## 2026-08-24 至 2026-08-25：Qwen 模型目录与 UI 调试

- 新增 DashScope 图片/视频默认能力、空配置补全、前端 fallback 和大量临时日志。
- 数据库/API 历史快照显示 Qwen 图片模型曾返回约 38/39 项尺寸配置。
- 多份报告先后把格式、allowCustom、硬编码、缓存和组件加载判断为根因，但结论互相矛盾。
- 最终确认 `/create` 的真实设置 UI 未被直接调试；旧日志位于 Canvas 设置组件，缓存根因未成立。

## 2026-08-25：数据卷事故与恢复

- 历史执行 `docker-compose down -v`，删除用户账号、密钥、模型配置和生成内容所在卷。
- 使用 2026-08-17 的 SQLite 备份恢复并重新创建容器。
- 恢复只能证明备份可被加载，不能证明事故前全部数据已恢复。

## 2026-08-25：Codex 接管审计

- 完整审阅 Web、Backend、Canvas Agent、插件、根文档、`.claude`、官方百炼文档镜像和当前差异。
- 将旧报告降级为历史证据，建立当前状态、风险、经验和 NO-GO 计划。
- 用户确认交互式 PowerShell 中 GitHub CLI 认证正常；工具进程的无效 token 结果归因于环境凭据差异，待后续执行时在目标进程复核。

## 2026-08-25：建立项目记忆与百炼全局技能

- 在 `AGENTS.md` 增加项目记忆强制入口和维护规则。
- 新建 `.codex/project-memory/`，沉淀当前状态、时间线、经验教训、待办和维护规范。
- 完整读取并索引本地百炼官方镜像 473 个文本文件、223,797 行，为每篇记录 SHA-256。
- 新建并验证 `bailian-model-api` 技能，按文本、图片、视频、语音、向量、平台 API 和通用安全协议渐进加载。
- 将技能从项目源安装到全局 Codex skills，并验证 12 个文件哈希一致。

## 2026-08-25：锁定浏览器联调环境

- 用户确认本项目开发过程全程使用 Microsoft Edge，且 Edge 已完成登录。
- 后续管理后台检查、登录态验证和 UI 联调只使用 Edge，不使用 Chrome，也不自动回退其他浏览器。
- 此前 Chrome/内置浏览器显示的本地登录失效只代表错误浏览器环境，不作为项目实际登录状态证据。

## 2026-08-25：阶段 0 基线保护完成

- 记录 `main@0b202a1`、9 个已修改文件、35 个未跟踪文件，以及 Backend/Web 容器、镜像、Compose 和数据卷基线。
- 生成 Git binary patch、44 文件当前状态 ZIP、逐文件 SHA-256；ZIP CRC 和路径集合验证通过。
- 短暂停止 Backend 调度复制 SQLite/WAL/SHM/`.settings-key`/迁移标记，随后立即恢复，Backend 保持 running/healthy。
- 使用 SQLite Online Backup 生成一致性快照：完整性 `ok`、外键违规 0、60 张表，独立恢复副本验证通过。
- 使用匹配 `.settings-key` 在内存中成功解密 OSS 密文，只验证成功和长度，未输出明文。
- 发现 D 盘为 exFAT、无法可靠收紧 ACL；敏感快照和密钥转移到 C 盘私有 Codex 备份目录，仅允许 <current-windows-user>、SYSTEM、Administrators。
- D 盘原始数据库/WAL/SHM/密钥及快照副本已删除，只保留非敏感基线和安全位置指针。

## 2026-08-25：阶段 1 百炼官方发现目录

- 从本地官方百炼文档生成 80 项版本化 Manifest，并显式纳入 `wan3.0-video` 和 `wan3.0-video-prime`。
- 建立 `ready/planned/unsupported/deprecated` 支持状态；本阶段所有官方项为 `planned`，等待模型专属 Adapter 和合同测试后再晋升。
- 管理后台拉取从字符串列表升级为完整 CatalogItem，同步保留能力、协议、上游模型 ID、操作、来源版本和文档证据。
- 增加安全补齐与幂等边界：不覆盖已启用、已定价或已有能力配置的记录，并以更新时间保护并发写入。
- Backend 和管理 UI 同时落实非 Ready 只读门禁；拉取结果新增上游/官方补充、能力分类、新增/补齐等汇总。
- 只执行离线测试与构建，不迁移运行数据库、不重建容器、不进行 Edge 或真实模型验证。

## 2026-08-25：阶段 2 Wan 3.0 All-in-One Adapter

- 按官方 `wan30-video/create-task.md` 与 `query-result.md` 建立模型专属原生 DashScope 合同。
- 标准版 `wan3.0-video` 晋升 Ready；Prime 因本地精确合同和账号运行证据不足继续 Planned。
- 支持 T2V、首帧、首尾帧、参考图片/视频/音频，以及智能时长、adaptive 比例、有声、水印和 seed。
- 增加严格拒绝与 mock 提交测试，并移除 DashScope 视频日志中的完整提示词、Metadata、请求体和临时视频 URL 输出。
- 只执行源码、mock 和前端构建验证；没有运行环境变更或真实付费调用。

## 2026-08-26：CGO 测试门禁与阶段 3

- 为 Backend 增加隔离 Linux CGO 测试 target、Windows 脚本和显式 CI 工具链检查，关闭阶段 1 的 SQLite 迁移测试缺口。
- 在不挂载任何运行数据的测试镜像中两次完成 Backend 全量测试；现有服务容器前后完全一致。
- 分别实现 Wan 2.7 T2V/I2V/R2V 与 HappyHorse 1.1 T2V/I2V/R2V 模型专属 Adapter 和能力合同。
- 7 个 Wan 2.7 主线/快照和 3 个 HappyHorse 1.1 模型晋升 Ready；旧版、视频编辑和未适配模型保持 Planned。
- 引入跨媒体总数、输入相关时长和音频截断保护，清除旧通用 Adapter 的静默截断路径。
- 完成 Backend 专项/全量测试、Web 9 项专项测试、TypeScript 和生产构建；未进行运行数据库、容器或真实 API 变更。
- 修复 5 个历史 Go 文件的小规模格式差异，使 Backend 全仓 gofmt 门禁可通过；格式化后再次完成隔离 CGO 全量测试。

## 2026-08-26：GitHub 状态复核与阶段 4

- 证明交互主机中的 GitHub CLI 认证有效，受限进程报错来自 Windows Keyring 不可见；官方仓库、用户 Fork 和指定 head 当前均没有 PR，因此无需重置认证。
- 完成 Backend 隔离 CGO、Web 全量测试/生产构建、Canvas Agent 286 条全量测试/构建。
- 修正一组仅在 Windows/exFAT 高 I/O 负载下出现的 Agent 测试时序假设；未改变生产 30 秒租约和并发业务语义。
- 撤销未使用的 modernc SQLite 依赖漂移，保持生产 go-sqlite3 + CGO。
- GitHub Actions 增加 Web build 与独立 Canvas Agent test/build 门禁；尚未通过提交或 PR 触发远程执行。
- 确认文档站构建配置从 HEAD 和 `origin/main` 都不存在，保留为后续 P2 工程任务，不虚报文档构建成功。

## 2026-08-26：Rolldown 测试退出竞争修复

- 确认唯一在 Bun 主测试进程内启动 Vite 的 `local-channel-runtime-projection.node.test.mjs` 是 Rolldown worker 生命周期进入测试进程的路径。
- 改为 Bun 直接导入业务模块，移除测试内 Vite SSR server、临时 cache 和异步 close；专项 10 次、Web 445+1 全量、TypeScript 与生产构建均通过且无 panic。
- 增加 panic guard，将 Rolldown worker panic 从“退出码 0 的告警”升级为明确失败；合成 fixture 验证 guard 有效。
- 用根 `.bun-version` 将本机和 CI 统一为 Bun 1.4.0；未升级 Vite/Rolldown 依赖，因为当前结构性修复已消除触发条件，版本升级保留为独立变更。
- 用户决定文档站暂缓，并授权随后进入阶段 5。

## 2026-08-26：阶段 5 发布候选审计

- 清理旧 Qwen UI/能力调试日志，并用测试锁定 missing、explicit-empty、wildcard 三态。
- 建立系统渠道持久能力唯一边界，Catalog 与任务 admission 对缺失、损坏和不完整配置一致 fail closed。
- 收紧旧数据回填和价格可用语义，防止无能力 JSON 或 `price_configured=false` 的记录进入可用目录/路由。
- 删除 DashScope 视频日志中的提示词、媒体 URL、API Key 片段、完整请求 URL 与完整响应输出。
- Manifest、Web 447 条、Canvas Agent 286 条、Backend CGO 全量门禁通过；构建 Backend/Web 独立候选镜像并完成无数据临时容器 smoke。
- 运行中的旧容器、数据库和卷保持不变；完整阶段 6 顺序与回滚边界登记在 `STAGE5_RELEASE_CANDIDATE.md`。

## 2026-08-26：阶段 6 备份与部署

- 在 C 盘私有目录创建 pre-deploy verified/protected 快照，短暂 pause Backend 复制 SQLite/WAL/SHM/密钥后立即恢复 healthy；主/恢复副本与 OSS 解密均通过。
- 给旧 Backend/Web 镜像增加 pre-stage6 回滚标签，随后使用 `--no-build --no-deps --force-recreate` 逐个更新服务，未执行 `down`，数据卷保持原位。
- Backend 自动完成加法 Schema 迁移；241 条旧记录安全回填为 `planned/legacy`，业务计数不变，OSS 密文仍可解密。
- Web/Backend 新容器均 healthy、无重启/OOM/关键错误；首页原始响应与镜像内 index 哈希一致。
- 创建并验证 post-deploy 恢复点，登记迁移后数据库和新容器基线。
- 按阶段门禁未操作 Edge、未点击拉取模型、未调用真实百炼或 OSS 上传。

## 2026-08-26：阶段 7 目录拉取、幂等修复与管理后台验证

- 用户在已登录 Microsoft Edge 中执行模型拉取，运行目录从 241 条扩展到 319 条：官方补充 80 项，11 Ready、308 Planned。
- 管理后台 UI 验证了 Planned 只读、Ready 筛选恰为 11 条，以及 `wan3.0-video` 的视频能力、DashScope 视频协议、未配置价格和停用状态。
- 第三次拉取仍显示“补齐 239”，数据库摘要也发生变化；专项失败测试证明空目标 capability/protocol 被误当作补齐字段。
- 收紧同步条件并重新构建/保留卷更新 Backend。修复后的重复拉取显示“新增 0，补齐 0”，目录业务字段及 `updated_at` 摘要保持不变。
- 创建 `BKP-20260826-171800-STAGE7-CATALOG`，验证 60 张表、319 条目录、外键、独立恢复副本、OSS 密文、ACL 和恢复材料。
- 普通用户创作台 Edge 截图显示视频模型选择器为空并提示没有可用模型，证明停用且未定价的 Ready 模型没有被误开放；未点击生成。
- 增加公开系统目录回归测试，锁定停用模型被服务端省略、启用但未定价模型仍标记 unavailable；隔离 Linux CGO Backend 全量测试通过。
- 阶段 7 正式完成；未进行真实百炼生成或 OSS 调用。

## 2026-08-26：阶段 8 真实调用前收敛

- 读取 Wan 3.0 官方创建/查询合同和当前模型市场价格，并脱敏核对北京 MaaS 渠道、Wan 3.0 状态、积分余额与乌兰察布阿里云 OSS 配置。
- 在真实调用前发现兼容路径错误、结果下载跨主机转发凭据、generation 素材不进入 OSS 三项阻断；完成修复和回归测试。
- Backend CGO 全量、Web 448+1 全量和生产构建通过；构建并部署阶段 8 Backend/Web 候选，旧镜像打回滚标签，数据卷保持原位。
- 创建 `BKP-20260826-185248-STAGE8-PRECALL`。此时任务/资源/素材仍为 0，模型仍停用且未定价，未发生真实调用或费用。
- 用户随后明确批准一次 `wan3.0-video`、480P、16:9、2 秒真实调用，创建重试 0、供应商费用上限 0.60 元；第二次创建及参数扩展仍未授权。
- 调用前发现系统渠道账务没有按任务规格匹配精确价格档；修复为按 operation/vquality/videoSeconds 选择 SKU，并通过 Backend CGO 全量测试。
- 首次运行卷监控容器以 root 打开 SQLite，意外把 WAL/SHM 权限改为 root:root 0644，Backend 随即出现 readonly 重启。自动回滚没有解决共享卷权限；精确恢复 DB/WAL/SHM 为 app 的 100:101/0660 后服务恢复，数据计数仍为 0/0/0/0。
- 监控改为 Backend 相同 UID/GID，复核权限不再漂移；精确计费候选随后部署 healthy。真实任务仍未创建。
- 用户明确批准本次 UI 调用省略 seed，其他模型、参数、费用上限和创建重试合同不变。
- 用户在 Edge 保存 Wan 3.0 临时价格档；数据库确认只有 `text_to_video/480p/2s` 一个已启用价格档，模型临时启用，尚未创建任务或预留积分。
- Edge 预提交截图发现 Create 缺少音频/水印配置，且旧全局默认为有声；为避免违反无声授权，暂停提交并补齐能力驱动输出开关。
- Web 449+1、TypeScript 和容器构建通过，保留回滚标签部署新 Web；Backend/价格档/零任务基线未改变。
- 用户提交唯一 Wan 3.0 任务。刷新后画幅实际为 adaptive 而非 16:9；未重试。Provider create 一次成功，轮询后任务 succeeded。
- 请求参数审计为 adaptive/480P/2s/audio=false/watermark=false；账务结算 2 积分，OSS 生成一个 ready MP4 资源并创建 confirmed Asset，浏览器 Range 读取返回 206。
- 用户在 Edge 播放生成结果，确认 2 秒视频正常播放、无声、无水印；真实 Wan 3.0 与 OSS 视觉门禁通过。
- 用户将 Wan 3.0 模型和临时价格档恢复停用/未定价；数据库确认审计证据完整保留。创建并验证 `BKP-20260826-201444-STAGE8-POSTCALL`，阶段 8 完成。

## 2026-08-26：阶段 9 发布固化与最终交付

- 从 `main@0b202a1` 开始，将阶段 1–8 的实现拆分为 Backend `64015db`、Web `62cf5d8`、构建测试 `f97adde` 和最终文档记忆提交。
- Backend 隔离 Linux CGO 全量、Web 449+1、TypeScript、生产构建和静态门禁通过；官方补充目录稳定为 80 项、11 Ready。
- Canvas Agent TypeScript 构建通过，但当前 Windows 复验有 4 个进程树清理用例失败并挂住测试进程；阶段 4 的 286/286 降级为历史证据，未把本次候选写成全绿。
- 全量复核旧问题并按已修复、部分修复、未修复、待真实验证分类，详见 `STAGE9_RELEASE_HANDOFF.md`。
- 本阶段只做本地 Git 固化；没有 push、没有创建 PR、没有触发远程 CI，也没有新增真实模型调用或数据库变更。

## 2026-08-26：空间清理与阶段 10 Canvas Agent 稳定性

- 先按备份登记审计空间；保留阶段 0/6/7/8 的 verified/protected 恢复点、运行数据卷和阶段 10–16 所需 `node_modules`。
- 删除可再生成的 Go/索引缓存、Web/Canvas Agent 构建目录和旧 Backend 可执行文件，按清理前精确计数至少释放 775,473,371 bytes；旧 `stage1-go-mod` 因 Windows 文件属性和访问限制仍有残余，连续三种安全路径失败后停止强删。
- 受限进程专项稳定复现 4 个清理失败，并证明内部 `taskkill` 返回“拒绝访问”；主机权限下同一原始实现 11/11 全通过，确认不是生产 `taskkill /T /F + child close` 缺陷。
- 测试夹具增加精确终止能力探测：受限环境继续用注入终止器验证业务语义并明确 skip 精确树能力；主机环境真实终止专项 12/12。
- 全量测试首次 286/287，唯一失败是在原子替换窗口读取 `runtime.json` 得到 `ENOENT`；按既有规则仅重试瞬时 `ENOENT`，其他错误继续 fail closed。专项连续 5 次通过，最终全量 287/287、TypeScript 构建通过。
- 本阶段没有修改生产进程管理代码、没有运行容器、操作数据库、打开浏览器或调用外部模型。

## 2026-08-27：阶段 11 Create 设置持久化与输出开关

- 确认根因是 Create 的图片/视频 Effect 在刷新和模型重算时总把比例、时长、清晰度重置为模型默认值；声音和水印又依赖全局配置，无法按模型和生成方式隔离。
- 新增用户作用域 IndexedDB 草稿，以用户选择模型和生成方式为键；并发写入串行合并，损坏文档 fail closed，设置未加载完成前禁止提交。
- Create 的声音、水印改为真实 Ant Design Switch；摘要、消息历史和任务 requestConfig 使用同一显式状态。
- 专项 13/13、Web 主套件 453/453、跨 Runtime 1/1、TypeScript 和生产构建通过；构建 11,022 个模块，无 Rolldown panic。
- 旧 Web 镜像标记 `open-ai-canvas-web:pre-stage11-20260826-2350`，只重建 Web；Backend ID/镜像不变，双方 healthy、RestartCount=0，首页和 API health 均为 200。
- 用户在已登录 Microsoft Edge 确认 16:9、480P、2 秒和无声/无水印刷新后保持，声音/水印 Switch 与摘要同步；未点击发送、没有模型调用或费用。

## 2026-08-27：阶段 12 Qwen Image 3.0 Pro Ready

- 逐行读取 Qwen 同步文生图和图片编辑官方 OpenAPI，锁定单轮 messages、文生图单 text、编辑 1–3 image + 单 text、`auto` 省略、`宽*高`、`n=1-6` 和 24 小时 URL 合同。
- 新增 Qwen 3.0 专属同步请求构造器和保守能力；覆盖尺寸/比例、参考图格式与大小、prompt extend/thinking、负面提示词、watermark、seed 和严格拒绝，不把合同扩散到其他 DashScope 图片模型。
- `qwen-image-3.0-pro` 从 Planned 晋升 Ready；Catalog 版本更新为 `2026-08-27`，80 项，重生成哈希两次一致。
- Backend 隔离 CGO 全量、Web 456+1、TypeScript 和 11,022 模块生产构建通过；首次 Backend 门禁仅暴露一个局部 Go 短声明语法错误，修正后全量通过。
- 创建 `BKP-20260827-104744-STAGE12-PREDEPLOY`，完整性/外键/OSS 密文和 ACL 通过；给旧 Backend/Web 镜像加阶段 12 回滚标签后保留卷逐个更新服务。
- Edge 首轮拉取发现上游从 241 增至 242：总目录 320，新增 1、补齐 69；Ready 筛选 12，Qwen 显示图片/可用/DashScope 图片/停用/未定价，能力为 3 参考图、10MB、6 输出、auto、自定义，所有不支持项关闭。
- 第二次拉取新增 0/补齐 0。创建 `BKP-20260827-113656-STAGE12-POSTDEPLOY`；数据库 320、Ready 12、Planned 308、Qwen 能力 JSON 558 bytes，invalid 计数全 0。
- 本阶段未启用或定价 Qwen，未点击测试模型或生成图片，没有供应商调用和费用。

## 2026-08-27：阶段 13 统一错误链路与安全日志

- 确认 `ModelError` 嵌入 `*AppError` 但未实现 `Unwrap`，以及多条任务路由仍把 service error 交给旧 `fail`；修复后能力、价格、路由、Provider 与普通 HTTP 错误都进入统一安全投影。
- 增加请求编号中间件、响应头和稳定错误元数据；Gin 访问日志、handler、任务、Worker、Provider 与资源物化以关联 ID、状态、耗时和错误类型串联。
- API 调用与任务日志新增 prompt/text/content/input/output 等敏感正文过滤；Provider 原始 message、签名 URL、内部错误和非结构化任务失败不再进入响应、持久任务日志或进程日志。
- Backend 隔离 Linux CGO 全量测试通过；Web 458+1、TypeScript 和 11,022 模块生产构建通过，无 Rolldown panic。
- 创建并验证 `BKP-20260827-125039-STAGE13-PREDEPLOY`；保留阶段 12 镜像标签后逐个更新 Backend/Web，数据卷未删除或替换。
- 本地无费用验证返回稳定 401/400 元数据，合法 request ID 保留、非法 ID 替换，测试敏感标记未出现在 Backend 日志；没有创建任务、模型调用、OSS 上传或费用。
- 代码、测试、文档和工程记忆由本地提交 `f39512a` 固化；未 push、未创建 PR、未触发远程 CI。

## 2026-08-27：阶段 14 影策 Codex 插件与部署入口合同

- 对照当前 Web/Canvas Agent 签名 Local Runtime 实现，确认插件文档仍在指导用户把 master token 放入 URL，并错误使用 Next、`CANVAS_URL` 和 `/config` 读取 token；这些都是已经失效且不安全的历史合同。
- 将插件深链收敛为只带 `mode`，备用端口改为 Vite，可信网页来源使用精确 `FRAMEFIELD_TRUSTED_WEB_ORIGINS`；MCP 超时从 90 秒提高到 2160 秒。
- root/local/dev/deploy/server 五个 Backend Compose 入口统一透传 `ENABLE_PROVIDER_PLUGINS`，默认 false；五份 Compose config 校验通过。
- 插件合同 4/4、Web 本机签名连接 17/17、Canvas Agent 主机全量 291/291、TypeScript、官方插件/skill 校验全部通过。
- 受限全量测试因 `taskkill` ACL 每分钟遗留子进程；精确终止仅本次测试树后在主机权限复验通过，最终测试进程为 0，运行 Backend/Web 未变化。
- 未执行插件 cachebuster/reinstall/new-thread 验证；`codex plugin list` 在 WindowsApps ACL 下无法启动。未重建容器、操作数据库或调用外部模型。
- 插件、Compose、测试、文档和工程记忆由本地提交 `8946917` 固化；未 push、未创建 PR、未触发远程 CI。

## 2026-08-27：阶段 15 统一系统报价与部署安全收口

- 将统一目录报价从仅前台逻辑模型扩展到系统渠道模型；系统模式按 channel model ID、持久能力和精确 SKU selector 构造只读账单快照，前端模型选择器改用统一端点。
- 删除从未进入数据库表清单和读写主链的 `LogicalModelPriceSKU`，校正旧 Backend 阶段总结，明确当前 `ChannelModelPriceTier`/LogicalModel 价格边界。
- Server Compose 的 CORS 从默认 `*` 改为显式必填；空覆盖配置失败、明确 Origin 配置通过。
- GORM 官方 `ParameterizedQueries` 对 Raw SQL 错误路径仍可能看到预插值值；最终用包装 logger 整体隐藏 SQL 文本，同时忽略正常 not-found 并保留真实数据库错误诊断。
- Backend 前三次门禁依次暴露错误构造函数命名、GORM Raw 边界和测试空表夹具遗漏；逐项修正后最终全量通过，没有修改运行数据库。
- Web 459+1、TypeScript 和生产构建通过；创建 `BKP-20260827-150444-STAGE15-PREDEPLOY`，保留旧镜像后逐个更新 Backend/Web，数据卷未替换。
- 无费用运行验证：不存在用户登录返回安全 401；报价路由鉴权返回 401；Backend 日志敏感标记、`record not found` 和 SELECT SQL 计数均为 0。未启用/定价模型或产生外部调用。
- 代码、测试、部署文档和工程记忆由本地提交 `f9317fa` 固化；未 push、未创建 PR、未触发远程 CI。

## 2026-08-27：阶段 16 最终本地交付

- 从干净的 `main@2b54c2e` 启动最终审计；gofmt、diff、私有路径和凭据启发式扫描通过。
- Backend 隔离 Linux CGO 全量、Web 459+1、TypeScript、11,022 模块构建、Canvas Agent Windows 主机 292/292 与构建全部通过；插件/skill、五份 Compose 和 Server CORS 门禁通过。
- Computer Use 两次因无法可靠确定 Edge URL 而安全停止，没有操作 Chrome；用户随后提供 Microsoft Edge 首页、管理后台、创作页三张截图，补齐人工只读验收。
- 创建并验证 `BKP-20260827-160505-STAGE16-FINAL`，数据库完整性 `ok`、外键 0、60 表、OSS 密文可恢复；给当前 Backend/Web 镜像增加阶段 16 最终本地标签。
- Web 镜像内 `index.html` 与容器 HTTP 首页哈希一致；容器 healthy、RestartCount=0、OOMKilled=false，数据卷未替换，安全日志计数正常。
- 生成 `STAGE16_FINAL_HANDOFF.md`，保留插件安装态、文档站、远程 CI、付费模型和 Planned 模型等明确后续门禁；按用户要求停在远程 push/PR 前。
- 最终交付报告、恢复点、Edge 人工证据和记忆清理由本地提交 `058f64b` 固化；没有远程写入。

## 2026-08-27：上游百炼能力贡献设计 Issue

- 只读复核官方 `main@5289bef4`、官方 `feature`、用户 Fork 和活跃 Fork；确认官方 feature 落后 main 414 个提交且独有 0，本地 main 不能整体作为 PR head。
- 官方百炼 discovery 已列出 Qwen Image、Wan 2.7、HappyHorse 等模型，但目录发现与项目可执行状态仍无显式分层；Wan 3.0 尚未进入 discovery。
- 用户授权后，在 `ddcat-ai/open-ai-canvas` 创建中文 Issue [#332](https://github.com/ddcat-ai/open-ai-canvas/issues/332)，正文包含支持状态、Ready 条件、数据兼容、Manifest、五项架构决策和六步 PR 拆分。
- Issue 创建后核对为 open、作者正确、正文关键章节齐全、评论 0、无指派/里程碑。账号无官方仓库标签写权限，`enhancement` 补加被拒绝；未创建分支、PR、push 或 CI。

## 2026-08-28：Issue #332 维护者决策状态复核

- 使用 `gh api` 只读检查 Issue 本体、评论和时间线；Issue 仍为 open，评论 0、标签 0、无指派/里程碑，`updated_at=2026-08-27T09:13:58Z` 与创建时间相同，五项架构决策均未获维护者回答。
- 时间线唯一新增事件是 PR #335 的自动 `cross-referenced`。进一步读取 PR 本体、Issue comments 和 reviews，确认它属于 3D 导演台改进，作者 `echoD886` 的身份为 `CONTRIBUTOR`，无评论/Review，异常重复正文偶然出现 `#332`；该事件与百炼无关。
- 决策状态保持 `pending`，原 NO-GO 不变：不创建支持状态 Schema、百炼 Adapter 或版本化 Manifest PR，不用标签、指派、点赞、关闭或无关交叉引用替代明确架构答复。
- 约定后续判定：有效维护者回复需来自 `OWNER/MEMBER/COLLABORATOR`，并逐项回答五个问题或明确接受推荐方案；笼统“欢迎 PR”只算部分认可。提交后等待 2–3 个工作日，后续评论必须重新取得用户授权。

## 2026-08-28：官方上游增量合并阶段 0-1

- 用户批准按阶段执行合并；阶段 0-2 可连续，阶段 3 起逐阶段批准。建立 `UPSTREAM_MERGE_PLAN_20260828.md` 和会话级持久备注。
- 冻结 `main@405d25a`、`VERSION=v1.1.4`、现有 remote 和运行镜像；Backend/Web 均 healthy、RestartCount=0、OOMKilled=false。
- 创建 `BKP-20260828-164620-UPSTREAM-PREMERGE`：Git bundle、当前工作树 overlay、SQLite/WAL/SHM、`.settings-key`、迁移标记和独立序列化恢复副本均完成哈希与恢复验证。
- 当前数据已从阶段 16 继续增长到任务 3、资源 5、素材 5、画布项目 1；320 模型仍为 Ready 12、Planned 308。
- 首次副本校验发现 Docker 复制文件未正确继承 Windows 可读 ACL；仅对新备份目录逐文件授权当前用户、SYSTEM、Administrators，随后双副本验证通过。live 卷未变。
- 给当前 Backend/Web 镜像增加 `pre-upstream-merge-20260828-164620` 标签；阶段 1 完成，可进入阶段 2。

## 2026-08-28：官方上游增量合并阶段 2

- 将远程规范为用户 fork `origin`、官方 `official`、本机旧快照 `legacy-upstream-snapshot`；官方、旧快照及其他人的 fork 禁止 push，默认 push remote 为用户 `origin`。
- 只读 fetch 官方 main/feature/tags 和用户 fork。官方 main 从旧分析的 `913cf4b` 前进到 `ab89c05`；用户 fork main 仍为 `11931d0`，未 push。
- 重新计算共同基点 `70a6640`、本地独有 38、官方独有 49、双方改动 180/504 文件、重叠 54；新 merge-tree 仍为同一 29 个显式冲突文件。
- 官方 feature 独有 0、落后 main 427；无需合并。官方 `VERSION` 仍是 `v1.2.0-preview.1`，但 main Git 描述为 `v1.2.1-21-gab89c05`，版本文件列入语义冲突。
- 阶段 0-2 完成；没有 merge/rebase/切换分支/push。阶段 3 等待用户批准。

## 2026-08-28：官方上游增量合并阶段 3 启动

- 用户确认其 fork 是集成分支未来的线上归属，并批准阶段 3 按推荐方案执行；本阶段仍不 push。
- 将 `CURRENT_STATE.md` 中阶段 2 前的旧 remote 描述明确降级为历史，当前真相保持 `origin`=用户 fork、`official`=官方、`legacy-upstream-snapshot`=本机旧快照。
- 阶段 3 计划从阶段 0-2 纯文档基线 `2c584f0` 创建 `codex/upstream-20260828-ab89c05`，以 no-commit/no-ff 引入固定官方 `ab89c05`；只确认冲突，阶段 4 获批前不解决业务冲突。

## 2026-08-29：官方上游增量合并阶段 10 第一提交与停止门禁

- 用户批准“两步可审计合并”。先修正项目记忆、`BACKLOG.md` 和历史价格说明，保持 `VERSION=v1.1.4`，再创建本地 merge commit `d04c4d1d968ad2c53de919985596d61c8d841f5f`；父提交为本地 `03604df` 与固定官方快照 `ab89c05`。
- 提交后重新只读查询官方实时 `main`，发现已从原审计目标 `115e228` 前进到 `80d2aa6`。GitHub Compare 显示新增 11 个提交、96 个文件、6,195 行新增和 526 行删除。
- 新增范围包含视频批量取帧、生成结果命名、视频设置、短剧分镜工作流、媒体引用/模型匹配、公告、存储管理、定价与品牌素材，不再是原预计的 README 加两张赞助商图片。
- 按预设门禁停止，没有 fetch/merge 第二步，没有 push、PR、远程 CI、部署、数据库/卷操作或 Provider 调用。后续需以 `80d2aa6` 为新固定 SHA 重新审计并获得用户批准。

## 2026-08-29：官方上游增量合并阶段 11 新增量审计

- 用户批准新增量审计；开始时官方实时 `main` 已从 `80d2aa6` 前进到 `4f07daa`，因此以 `4f07daae9ec3b4e4cb0a8cd35a6e0c1a4b593f29` 为固定候选并 fetch 到只读跟踪引用。结束前再次确认官方 SHA 未漂移。
- 相对已合并的 `ab89c05`，官方新增 14 个提交、107 个文件、6,843 行新增和 661 行删除；本地/官方新增量重叠 22 个文件。
- 只读 `merge-tree` 模拟得到 6 个显式冲突，位于 Backend 计费、待测文档、Web 测试入口、视频设置、能力模型和 Create 页面；另列出 8 类自动合并但必须语义复核的安全/数据区域。
- 新功能包含视频批量取帧、生成命名、能力驱动尺寸、短剧工作流 v2、默认价格、图生图引用、公告配图/Markdown、后台存储管理、资源直链扩展名和 AutoDL 参数对齐。
- 生成 `UPSTREAM_INCREMENT_AUDIT_20260829.md`，明确逐文件解决合同、分层执行步骤和验证门禁。本阶段没有真实 merge、源代码冲突编辑、push、部署或运行数据操作。

## 2026-08-29：官方上游增量合并阶段 12A 打开固定 SHA 合并

- 用户批准阶段 12A。开始前确认官方实时 `main` 仍为固定 `4f07daa`，工作树干净、`VERSION=v1.1.4`；私有上游合并前恢复点的数据库、WAL、独立恢复库、设置密钥、迁移标记和 Git bundle 六项哈希全部复核一致。
- 执行 `git merge --no-commit --no-ff 4f07daa`，当前 `HEAD=ae58ecf`、`MERGE_HEAD=4f07daa`，未创建 merge commit。
- 真实冲突与阶段 11 模拟完全一致：6 个文件、8 个冲突区块；107 个官方增量路径进入索引，无意外冲突新增或缺失。
- 本阶段未处理冲突、运行源码测试、push、部署或操作运行数据；停在阶段 12B 批准门禁。

## 2026-08-30：官方上游增量合并阶段 12B 冲突解决

- 用户批准阶段 12B；按“测试/文档 → Backend 计费 → 能力/视频设置 → Create”的顺序手工解决 6 个文件、8 个区块，未使用批量 ours/theirs，未解决索引项归零。
- Backend 计费改为服务端按完整 intent 选择价格档，旧/客户端 tier ID 不得覆盖真实规格；逻辑模型和线路切换同样使用服务端选出的 tier。补充旧 ID 拒绝和请求规格重选测试。
- Web 保留百炼精确能力、智能时长、设置持久化、声音/水印和 panic guard，同时吸收插件工作流能力、画幅可为空、只读尺寸推导、视频取帧和图片+音频模型匹配；无画幅模型不显示或沿用旧 size。
- Web 相关 9 文件 69/69、TypeScript、Backend 纯逻辑专项和 AutoDL 制品测试通过。宿主数据库测试只有已知 `CGO_ENABLED=0` stub 失败，隔离 Linux CGO 全量门禁留到阶段 12C。
- 结束复核发现官方实时 `main` 已前进到 `2f6832f`；相对固定 `4f07daa` 是 1 个涉及 113 文件、7,016 行新增和 723 行删除的短剧工作台/技能运行/任务恢复大提交。本轮未 fetch 或合并，未来需独立审计。
- 没有 merge commit、push、部署、数据库/卷或 Provider 操作；停在阶段 12C 批准门禁。

## 2026-08-30：官方上游增量合并阶段 12C 完整验证

- 用户批准阶段 12C。隔离 Linux CGO Backend 全部包通过；Web 1082/1082、跨 Runtime 1/1、TypeScript 和 13,493 模块生产构建通过；Canvas Agent 主机 327/327 和 TypeScript 构建通过。
- Prettier 首次发现 36 个官方新增/修改文件格式漂移，机械修正后全部受支持暂存文本、TypeScript、相关 Web 69/69、AutoDL 制品专项和最终 13,493 模块生产构建复验通过；Go 格式、cached diff、冲突标记和待测文档标题检查通过。
- 六份 Compose 与 deploy+build 叠加配置解析通过；生产 Server 缺密码时按设计拒绝，使用仅供解析的虚拟值后通过。未启动、重建或部署服务。
- 两张赞助商图片可解析并完成本地视觉检查；README 使用固定 160 像素宽度。没有 Edge 浏览器渲染验收。
- 结束时官方主线又前进到 `0893741`；相对固定 `4f07daa` 新增 2 个提交、122 文件。本轮未 fetch 或合并尾差，只验证固定 `4f07daa`。
- 阶段 12C 完成但没有 merge commit、push、部署、数据库/卷或 Provider 操作；停在本地提交批准门禁。

## 2026-08-30：官方上游增量合并阶段 12D 本地提交

- 用户批准阶段 12D；将完整验证后的固定 `4f07daa` 合并固化为本地 merge commit `e2326856b1cc3a854ea11fcd54c5c3e2866c8ac3`。
- 双父提交为 `ae58ecf` 与 `4f07daa`，提交后合并状态结束、工作树干净、`VERSION=v1.1.4`。
- 官方实时 `main` 仍为 `0893741`，相对固定目标的 2 个提交/122 文件未 fetch 或合并。
- 本阶段未 push、创建 PR、触发远程 CI、部署、操作数据库/卷或调用 Provider；远程写入继续等待独立批准。

## 2026-08-30：用户 fork 集成分支远程检查点

- 用户授权执行上一条建议中的第一项：只把已验证集成分支推送到用户 fork，不更新 `origin/main`，不处理官方新尾差。
- 新建远程 `origin/codex/upstream-20260828-ab89c05`，首次 SHA 为 `265c997`，本地分支设置为跟踪该远程分支。
- 复核 `origin/main` 仍为 `11931d0`，同名 head 没有 PR。没有部署、数据操作或 Provider 调用。

## 2026-08-30：官方 `c8b60ce` 尾差独立审计

- 用户同意下一步后只执行尾差审计。官方实时 `main` 已从 `0893741` 前进到 `c8b60ce`，fetch 到只读 `official/main`，审计结束前未漂移。
- 相对已集成的 `4f07daa`，官方新增 4 个提交、164 个文件、26,364 行新增和 2,696 行删除；双方重叠 28 文件。
- merge-tree 模拟得到 17 个冲突文件、48 个区块。风险集中在技能包/自动更新、启动迁移、短剧工作台、任务恢复、管理后台大规模重构和设计系统。
- 新增 `UPSTREAM_TAIL_AUDIT_20260830.md`，记录逐文件解决合同、技能安全、CSS 收口和分阶段验证方案。没有真实 merge、业务代码修改、远程写入或运行数据操作。

## 2026-08-30：官方尾差阶段 A 打开固定 SHA 合并

- 用户批准尾差阶段 A；开始前确认官方实时 `main` 仍为 `c8b60ce`，工作树干净、`VERSION=v1.1.4`，既有恢复点六项哈希一致。
- 执行 `git merge --no-commit --no-ff c8b60ce`，当前 `HEAD=3510997`、`MERGE_HEAD=c8b60ce`，未创建 merge commit。
- 真实冲突与审计完全一致：17 个文件、48 个区块；164 个尾差路径进入索引，无意外冲突新增或缺失。
- 本阶段未处理冲突、运行源码测试、push、部署、技能同步或操作运行数据；停在尾差阶段 B 批准门禁。

## 2026-08-30：官方尾差阶段 B 冲突解决

- 用户批准尾差阶段 B；解决固定 `c8b60ce` 的 17 个冲突文件、48 个区块，未解决索引项归零。
- 保留本地迁移、计价、百炼/Create、30MB Runtime 和 panic guard 合同；接入官方技能包、项目分页、真实进度、短剧恢复和后台重构，并收口模型编辑器 raw color。
- Web TypeScript及后台/技能/短剧/恢复专项、Canvas Agent 12/12、Backend 纯逻辑专项通过；完整跨栈门禁留到阶段 C。
- 官方实时主线又前进到 `4ba9694`，相对固定目标新增 8 个直接修复短剧生成/资产/视频协议的提交。本轮未 fetch/merge；停在“先审计新修复尾差或验证当前固定合并”的用户决策门禁。

## 2026-08-30：`c8b60ce..4ba9694` 八个修复提交独立审计

- 用户要求先独立审计。线上 `official/main` 确认为固定 `4ba9694`；只 fetch 到 `official/audit-4ba9694`，没有更新本地 `official/main` 或当前合并现场。
- 8 个提交共 42 文件、+2,398/-414。前 5 个形成默认模型/计价、并行生成、历史角色图/当前声音、视频角色协议和 mention 的耦合修复组；后 3 个是浅色画布与首页仪表盘/UI 收口。
- 使用当前已解决 index 构造无引用审计提交，merge-tree 得到 11 个重叠文件、3 个显式冲突文件/15 区块；工作台单文件 13 区块。
- 审计发现 xAI 仅尾帧在前端、Backend 类型化 Provider 和协议 Adapter 中分别被拒绝、当普通参考、当起始图；首页任务卡未遵守功能开关，查询失败伪装为零数据，且最多 300 条不能保证完整周统计；浅色画布使用硬编码颜色。
- 固定快照 TypeScript 退出码 0，Backend 协议包和 4 个新增纯 Provider 角色测试通过；SQLite 服务测试仍需 Linux CGO，Web Bun 测试和 Edge 视觉验收未执行。
- 结论为 NO-GO：不先创建当前 `c8b60ce` merge commit，不原样接受全部 8 提交。建议下一步单独批准修复尾差阶段 B2；完整证据和解决合同见 `UPSTREAM_FIX_TAIL_AUDIT_20260830.md`。

## 2026-08-31：修复尾差阶段 B2

- 用户批准 B2。官方实时主线已前进到 `fb089b2`，本阶段只处理固定 `4ba9694`，新尾差保持在外。
- 复核恢复点 14 个文件哈希和逐文件 ACL 后，将当前 `c8b60ce` 已解决索引固化为安全树/双父安全提交 `82758b2`，再从干净 `HEAD` 标准打开 `4ba9694` no-commit 合并。
- 直接重开出现 18 文件/52 区块；恢复安全树三方结果后收敛为审计预测的 3 文件/15 区块并逐区块解决。当前 `HEAD=3510997`、`MERGE_HEAD=4ba9694`、未解决 0、非暂存 0、`VERSION=v1.1.4`。
- 前五个短剧/协议修复组已纳入，并修复 xAI 尾帧三路径不一致；浅色画布使用组件 token；首页仪表盘暂缓，只保留工作台横向滚动修复。
- TypeScript、尾差 Prettier、xAI 协议/Service 专项、gofmt、JSON、冲突标记和 diff check 通过。Web Bun 因主机无 Bun 且 Docker bind source 不可见而 0 测试，未记为通过；Linux CGO/全量/构建/Edge 留阶段 C。
- 没有 merge commit、push、PR、CI、部署、数据库/卷、技能同步或 Provider 调用；停在阶段 C 批准门禁。

## 2026-08-31：固定 `4ba9694` 阶段 C 完整验证

- 用户批准阶段 C，并限定完成后只评估 merge commit/push 可行性，不执行。
- Backend 正式 test target 在隔离 Linux CGO 下全部包通过；测试镜像 manifest list 为 `sha256:4b57312dbd2f1be77bc11c51eaad9573b9fada7727014b39ac62ea979ebb6993`。
- Web 首轮主套件 1119/1122，3 个失败均是测试夹具/空白锚点漂移；只修测试后专项 23/23，最终 admin 10/10、主套件 1122/1122、跨 Runtime 1/1、TypeScript通过。
- 三平台 Comfy Bridge 和 Vite 13,469 模块生产构建通过；仅有既有大 chunk/插件耗时警告。Canvas Agent 最终格式化字节全量 328/328，build通过。
- 145 个暂存文本 Prettier、gofmt、JSON、文档标题和七组 Compose 通过；deploy 缺必填 URL 首轮按设计拒绝，虚拟解析值下通过。运行容器保持旧 `:local`、healthy、restart 0、OOM false，没有部署或 Edge 候选验收。
- 远程 fork 集成分支仍为 `9815464`，是当前 HEAD祖先；`origin/main=11931d0`。官方实时主线已到 `e124a9e`，新尾差未审计、未纳入。
- 结论：固定 `4ba9694` 本地 merge commit 技术可行但需新授权；当前未提交状态不可推送。提交后同名 fork 集成分支可 fast-forward 推送但仍需再授权；不 force、不更新 `origin/main`、不声称跟上实时官方。

## 2026-08-31：固定 `4ba9694` 本地 merge commit 授权

- 用户单独批准创建本地 merge commit；push 仍只做评审分析，不执行。提交前必须再次验证 `HEAD=3510997`、`MERGE_HEAD=4ba9694`、最终暂存树、未解决/非暂存 0 和 `VERSION=v1.1.4`。
- 本记录随合并提交固化；实际 merge commit SHA 在提交后由 Git 历史和会话记忆登记。推送、PR、CI、部署、运行数据和官方新尾差继续保持独立门禁。

## 2026-08-31：固定 `4ba9694` 本地 merge commit 与只读发布评审

- 创建本地 merge commit `d0def5f807e7d8c410103c85bd5ded43506b39b3`，双父精确为 `3510997`、`4ba9694`，tree为阶段 C最终暂存树 `49580c569687e0ef30b47cce0f5d1aa8e9aa9ff7`；提交后 `MERGE_HEAD` 消失、工作树干净、`VERSION=v1.1.4`。
- 未 push。fork公开、账号具备admin/push，同名集成分支未保护；远程仍为 `9815464`，`origin/main=11931d0`。`git push --dry-run` 确认可 fast-forward到 `d0def5f`，无需 force，将新增14提交、194路径、约 +53,429/-15,401。
- Actions已启用，但 `quality.yml` 和 `publish-images.yml` 的 push过滤都只含 `main`；推送集成分支不会自动运行远程质量或发布镜像。历史含官方首页仪表盘提交，但最终树按决策恢复稳定首页，不能以祖先提交存在推断功能启用。
- 当前运行容器仍是阶段16：Backend `f94f555b`、Web `31417063`，创建于8月27日，运行代码来自 `f9317fa`；Web入口哈希仍为 `bf8376e...`，healthy/restart0/OOMfalse。它比新本地提交少90提交、659个变更文件。
- 官方实时固定为 `b7348ab43cd174354724d3e6f9d88d105479c4e1`、`VERSION=v1.2.3.1`。相对本地固定基点 `4ba9694` 新增22提交、216文件、+10,391/-955，包含在线更新/备份回退、密码找回、画布媒体与资产、后台/渠道/存储修复和版本发布。
- 新本地提交与官方实时分别有48/22个独有提交，直接树差461文件；只读未来 merge-tree预测29个冲突文件、52区块。当前部署、本地提交、官方实时三者均不是可互换的同一版本。

## 2026-08-31：用户 fork 固定 `4ba9694` 远程检查点

- 先提交5份提交后项目记忆为 `768705338b525af63ba8f9fb1babfcd7b92670cd`。推送前重新读取远程仍为 `9815464`，tracking ref一致且为本地祖先；dry-run通过。
- 实际push将 `origin/codex/upstream-20260828-ab89c05` 从 `9815464` fast-forward到 `7687053`，没有force、没有更新`origin/main`。随后独立 `ls-remote` 返回同一SHA，本地/远程 ahead/behind为0/0，工作树干净。
- 本记录作为最终检查点记忆提交继续fast-forward同步；最终远程必须与包含本记录的当前HEAD精确一致。非main push未触发quality/publish工作流，不能声称远程CI通过。

## 2026-09-01：路线B保留卷升级至固定 `4ba9694` 本地候选

- 新增本地部署标识提交 `c3f5479`，根版本改为 `v1.1.4+local.4ba.e02cee1`，未推送fork。
- 创建 `BKP-20260901-114427-ROUTE-B-PREDEPLOY`：完整数据、原始DB/WAL/SHM/密钥/迁移标记、独立恢复库、Git bundle和回滚镜像标签均验证；完整性ok、外键0、60表、业务计数1/1/320/3/5/5/1，逐文件ACL无额外主体。
- 构建Backend `8182cc0a`、Web `522109d7`候选；Web 13,469模块并包含新版本标识。
- 禁网克隆卷第一次迁移60→72表、生成33个规范化skill包；第二次启动表数和全部业务计数稳定，迁移错误0、密钥不变。演练容器/卷随后删除。
- 在旧镜像回滚标签保护下，仅执行local Compose `up -d --no-build --force-recreate`，保留真实 `open-ai-canvas_backend-data`；新容器healthy/restart0/OOMfalse。
- HTTP健康、Web版本、日志和升级后真实卷快照通过；数据库完整性ok、外键0、72表，关键Schema/计数与演练一致。当前停在用户Microsoft Edge人工验收前，无真实模型/OSS外部写入。
