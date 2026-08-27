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
