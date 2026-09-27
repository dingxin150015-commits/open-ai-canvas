# v1.5.6 无提交合并接管记录

## 已确认边界

- 用户选择方案 A：保留本机 Canvas Agent、17371 Runtime、影策 Codex 插件、本机渠道和人像能力；将它们适配到上游新结构。该选择取代 2026-09-22 的 `needs-human` 产品门禁。
- 2026-09-24 用户进一步选择：旧 `/sessions` 影视 Agent 会话改接上游新 Agent；本机 Comfy Bridge 设备管理/轮询/结果接口保留并迁移。该决定不授权丢弃旧会话记录、伪造新 Agent 运行结果、删除 Bridge 设备记录或部署到真实数据卷。
- 独立合并树：`codex/upstream-v1.5.6-20260921`，`HEAD=6b16d25ee8bcfbba40ebf384f89f122fdf2ad90a`，`MERGE_HEAD=473d535ec161fd21001e1ad904d3b316e2d45c39`（官方正式 v1.5.6）。起点为本地稳定锚点 `d9fec0f4` 加后续百炼插件增量；根目录旧分支不是实现来源。
- 2026-09-23 开始接管时 `git diff --name-only --diff-filter=U` 为 197 项。已合成 Batch 2–5 的部分源码；Batch 0/1 生成物、Batch 7–30 尚未收口。`backend/internal/handler/routes.go`、`.github/workflows/quality.yml`、`web/package.json` 仍未合并。Batch 6 的本机渠道迁移正在进行，前端、Agent 与插件等大量冲突仍在。
- 合并尚未提交；不允许 `commit`、`force`、push 或部署。运行容器与真实数据卷不是本合并树的验证环境。

## 已审计的高风险合同

- 本地 Schema 9 已用于 `provider_duration_metering`，上游 Schema 9 已用于 `channel_presentation` 且推进至 29。`backend/internal/database/migrations.go` 暂按双谱系合成为 Schema 30，必须在隔离 Linux CGO 环境分别验证本地 Schema9 与官方 Schema29 快照的升级、二次启动和摘要保持；当前仅有 Windows 编译门禁，不能称迁移通过。
- 上游 `backend/internal/service` 已拆入 `internal/app`、`platform` 等包。本地目录资格、百炼能力、实际时长、任务日志和本机渠道需迁到对应新包；不得复活整套旧 service。
- 上游删除本机 Agent/渠道/人像路径，但用户已决定保留。只允许精确 loopback、安全出站、签名握手等本地安全合同；Host Updater 仍不得指向官方镜像覆盖 Fork。
- Git 自动合并的重叠路径仍需语义复核。消除 `<<<<<<<` 只证明文本冲突已解决，不证明本地和上游行为都保留。

## 当前下一步

### 2026-09-23 增量进度（未完成，不能发布）

- Batch 6 本机渠道迁入 `internal/app` 并通过 `go test ./internal/app -run '^$'` 编译门禁；真实 SQLite/任务行为尚未验证。
- 根据用户方案 A，精确恢复上游已删除的 `canvas-agent/`（87 文件）与 `plugins/yingce/`（10 文件），并恢复被删的本机 Agent、人像 Web 文件与相邻测试；未覆盖上游仍存在的同名文件。
- `.github/workflows/quality.yml` 保留双方 CI 意图；`web/package.json` 补回 `build:bridge`；模型协议 Token 计费采用上游“所有视频能力可选”规则，保留本地 Ark 协议辅助函数；另合成元数据探测、画布导出、人像复制与错误诊断等低耦合文件。`bun test test/model-protocols.test.ts` 为 2/2 通过。
- 本次更新时 `git diff --name-only --diff-filter=U` 仍有 148 文件，集中于插件目录生成器、Web 画布与服务调用链；任何 `git diff --check` 全局检查此时都会因尚未解决的冲突标记失败，这不是验证通过。Legacy `/sessions` 路由与上游 Agent 架构互斥，是后续高风险依赖门禁。
- 仍无 commit、push、部署。不得把此工作树当成已完成镜像。

### 后续验证补记（2026-09-24）

- Canvas Agent 锁定依赖安装成功，`npm test` 运行了会话、即梦调度与安全用例，已显示大量通过；随后超过 90 秒无新输出，人工中断，退出码 1，不能称整套测试通过。一个 Windows `taskkill` 用例提示受限测试进程下跳过。需隔离定位停顿用例。
- Web 工作树尚无完整依赖；`bun install --frozen-lockfile` 即使设置工作区 `BUN_TMPDIR` 与缓存目录、甚至提升沙箱权限，仍报 `EINVAL accessing temporary directory`，独立 Web 测试因缺少 `zustand` 无法启动。这是工具链门禁，不能写成源码测试失败或通过。
- API 请求层保留上游 `reason`/`http` 与本地分类、requestId、retryable 合同；局部 Web 合并后未合并文件数降至 134。上游已抽取的画布提示词写入 helper 取代旧局部副本，保留人像、ComfyUI 与画布主题局部功能。
- 再处理画布/布局/媒体低耦合冲突后未合并文件数降至 121：`canvas-project-domain` 保留本地模式存储与上游生成契约；`model-capabilities` 保留本地 DashScope 与上游 Agnes；`infinite-canvas` 吸收上游只读开关和活动主题；导演台模型入库采用上游本地/远端状态提示；ComfyUI Bridge 内置注册和人像、Ark 引用保留。
- `workspace-top-bar-extension` 使用上游 mount/slot 结构，项目详情改用上游页内 Header 以避免旧 `useWorkspaceTopBarExtension` 悬空调用。后续仍需核对顶栏、导航与项目页视觉行为。
- 插件目录生成器以上游 v1.5.6 为底，但显式跳过五个本地 v4 Gemini/OpenAI 包，防止 v2 再生覆盖 v4 Manifest。五个同名包的本地 v4 源 Manifest 与文档保留；将上游 Gemini/Responses 过滤 system 消息、Gemini 图片质量档映射补入本地 v4，重嵌文档后五包仍为 `yingce.plugin/v4`。对应二进制包尚待最终从 Manifest 重新打包，不能使用旧归档验证新映射。
- 五个 v4 包归档随后已从各自 Manifest/README/interface 文档重新打包并逐项验证源/ZIP 三文件字节一致；无手工二进制合并。该项仅证明包内一致性，尚未证明运行时执行或镜像部署。未合并文件数降至 105。
- `go test ./internal/protocol` 首次发现 Gemini Image 本地 v4 包仍写 `candidateCount`，违反上游明确合同；已删除该映射、重嵌文档并重包，复测 `go test ./internal/protocol` 通过。此为协议单元测试，不代表付费生成或完整应用联调。
- 2026-09-24 后续：导航保留本地独立首页 `/` 和创作入口 `/create`，同时采用上游短剧 Agent/自由画布标签及资源组；`router.tsx` 采用上游 CanvasRefreshShell 加载壳。客户端根初始化保留本机即梦目录 bootstrap；账户菜单采用上游可复用卡片。未合并文件数 85。
- 上游删除了历史 `backend/IMPLEMENTATION_SUMMARY.md`，其中自述旧设计且引用已迁出的 service，已按上游删除；Git HEAD 可恢复，该操作不涉及运行数据。
- 后续吸收上游注册限流冷却、任务视频缩略图、插件上传防重复提交、画布只读连线/上下文菜单以及提示词写入 helper；本地 Canvas Asset Tray 保留窄屏 `compact` 与侧栏偏移，上游“素材+数量”按钮取代旧纯数字按钮。此时未合并文件 81。
- Web 锁定依赖安装仍报 Bun 临时目录 EPERM/EINVAL；三种工作区临时目录与提升沙箱权限均未成功，停止重复尝试。曾尝试在隔离工作树为已有根目录依赖创建 Junction，系统拒绝创建，已恢复原来的空 `node_modules` 目录，没有留下链接或修改根依赖。完整 Web 测试/构建仍阻塞。
- 报价底层按上游新 `ModelQuoteRequest`（`logicalModelID` 或 `channelId+modelKey`）合并，但保留本地在浏览器阻止已知无价格规格反复 400 的 `modelQuoteValidationError`。模型选择器及分镜工作台已改用新 `quoteModel`，仍显示本地明确的报价失败原因，分镜工作台保留删除/保存/生成按钮。画布节点提示词面板仍未解决，其旧 `quoteModelCatalog` 调用必须随新合同替换。此时未合并文件数 68。
- 项目章节资产提取采用上游角色/场景/道具统一链，删除了页面快照去重方式；项目库采用上游资源缓存与显式删除流程。旧会话 `/sessions` 是否保留已向用户发起非阻塞确认，未获回复前不处理该高风险路由。
- 发现仅看 `diff-filter=U` 会漏掉已暂存的自动删除：当时还有 62 个 D。按方案 A 精确恢复本机 Runtime/即梦、桌面渠道、人像节点及其前端服务、相邻测试和独立插件，不复活已拆包的旧 service 或模拟转写。余下暂存删除为 28，仍须逐项审计 Comfy Bridge、Session、历史工件等产品边界，不能简单当作上游已吸收。
- `canvas-node-content.tsx` 合成上游文件上传/批量表格/媒体转换/稳定资源文件地址和本地人像排查节点；显式冲突现为 62，仍无完整 Web 构建。恢复本地文件本身不证明调用链已连通。
- 截至 2026-09-24 本轮停止点：`HEAD=6b16d25e`、`MERGE_HEAD=473d535e`，显式未合并仍 58 文件，暂存删除 28（混合已淘汰旧 service 和待判断 Comfy/Session 等），`git diff --cached --check` 已通过；全局构建不可运行，因为冲突标记尚在。无 commit、push、部署。不能宣布合并完成。
- 必须先得到产品边界答复：旧 `/sessions` 影视 Agent 会话 API 已被上游删除，但本地任务页仍调用；是保留并迁入新 app/repository/model，还是让旧页面改接上游新 Agent 并移除旧 API。该决定涉及数据模型和原任务历史，不得用假接口/空实现蒙混。
- 下一阶段需按依赖处理 `web/src/services/api/task-center.ts`、`web/src/pages/tasks/index.tsx` 与 backend `routes.go`；Comfy Bridge 旧后端同样未迁移，前端旧设置/Comfy 包的保留范围需要复核。随后完成画布页/生成链路、管理员计费、样式与测试，再做双谱系 Schema 隔离验证。Bun Web 依赖安装故障未解决，不能用缺依赖的 0 测试宣称通过。

### 2026-09-24 用户决策后的实施进度

- 旧 `/sessions` handler 冲突已按“改接上游新 Agent”解除；未删实际数据库表和业务数据。Web 任务页及旧画布助手仍待改接，不能视为完成。
- `web/src/services/api/task-center.ts` 已去掉失效的旧会话 DTO，任务页已停止提交 `/sessions`，提供“前往画布使用 Agent”入口；任务页其余 UI 冲突已合成并保留本机即梦任务刷新/删除按钮。旧 `canvas-assistant-panel.tsx` 仍通过画布页大型冲突块可达，需切到上游 `CanvasCloudAgentPanel` 后才能删除旧轮询 helper。
- Comfy Bridge 旧 handler/model/repository 已精确恢复，旧 service 状态机迁入 `internal/app/comfy_bridge.go`，在 `handler/api.go` 注册；工作流插件目录、provider `bridgeId`、执行、取消前退款/已领取待核对路径已接回。Schema 从 30 提升至 31，两个来源谱系都在末尾追加 `comfy_bridge_runtime`，跨库迁移表清单也补入设备和请求表。
- Windows Go 编译门禁：`go test ./internal/app ./internal/handler ./internal/database ./cmd/migrate-sqlite-postgres -run '^$'` 通过。隔离 Linux 镜像中 `TestComfyBridgeRequestSurvivesServiceRestart` 与 `TestMigrationPlan...` 通过；这是专项测试，不等于所有安全/取消/生成验收。
- 隔离 Linux `go test ./internal/database` 全包通过，覆盖实际迁移；Comfy Bridge 新增 Token/跨用户撤销测试已通过 Windows 编译，尚待 Linux CGO 执行。`TestTaskInputUsesWorkflowProvider`、`TestComfyBridgeWorkflowRequiresDeviceAndWorkflow`、`TestBundledWorkflowPluginsControlAvailability` 已在 Windows 测试运行通过。
- 新增 Token/跨用户撤销测试首次 Linux 运行失败，是测试凭据短于真实 32 字符门禁，并非认证实现失败；测试夹具已改为合规长度并保留另一条错误 Token 查找路径，等待重跑确认。真实令牌未写入记忆。
- 夹具修正后隔离 Linux `go test -count=1 -run TestComfyBridge ./internal/app` 通过，覆盖持久队列重启、正确/错误 Token、跨用户撤销和撤销后失效；仍未覆盖真实 ComfyUI 生成与费用取消端到端。
- 全量 Linux CGO 暴露上游拆包后遗留的 20 个本地 `internal/service/*test.go`；它们已机械迁到 `internal/app`，未删除测试。但 `go test ./internal/app -run '^$'` 仍因旧私有函数/旧返回值合同编译失败（首批涉及 `providerRequestIDFromHeader`、`runProtocolAdapterTaskWithTiming`、`sanitizeChannelModel` 等），须逐个迁接真实当前 API 或经覆盖对照后退役；不得加生产兼容桩只为让旧测试通过。此时不能宣称 Backend 全量通过。
- 后续已将上述 20 个本地测试全部迁入 `internal/app`，恢复安全的上游 `X-Request-Id` 响应头关联记录（真实生产路径使用，非测试桩），并把旧测试调用改到当前 `runProtocolAdapterTask`、双返回 `sanitizeChannelModel` 与 `resolveSystemChannelModelSelection`。`go test ./internal/app -run '^$'` 现已通过编译；各测试用例及 Backend 全量仍待执行。
- 2026-09-25 Windows 工作区缓存下 `go test ./... -run '^$'` 全仓 Go 编译门禁通过（只编译、不执行测试）。不等同于 Linux CGO 全量、Schema 克隆或真实 ComfyUI 生成验收。
- 隔离 `docker build --target backend-test` 全量失败：旧 `internal/service` 测试仍引用拆包前私有类型；另有 app 管理渠道测试被本机地址 SSRF 拒绝等失败。全量输出还包含大量 GORM 预期/非预期日志，须逐包收敛，不能称 Backend 全量通过。Docker 未运行生产 Compose、未挂载真实卷。
- 最新状态：显式未合并 57，暂存删除 23，`git diff --cached --check` 通过；仍是 merge-open、无 commit/push/部署。
- Web 画布基础冲突继续减少至 110 个；仍未进行完整 Web 类型检查，因依赖安装门禁与其他冲突标记尚在。当前工具链尝试不能代表业务验收。
- 2026-09-25 后续专项：后端 `go test ./... -run '^$'` 已在工作区缓存下全包编译通过；隔离 Linux CGO 的 `TestComfyBridge` 与 `internal/database` 全包已通过。用户明确选择旧影视会话 Agent 改接上游 Cloud Agent，Comfy Bridge 则保留并迁移，均不允许假兼容接口。
- v4 百炼协议回归测试暴露两个真实下载/参数问题：声明式请求需要传递允许列表内的 `rawOptions`、`generationMetadata`、`referenceCounts`，匿名媒体引用下载不得携带 Provider 凭证。已修复并通过 `go test ./internal/app -run 'TestBailianOpenAIHeaderRequestID|TestBailianQwenAllTransportsMockLifecycle' -count=1`。此为 mock 专项测试，不等于真实付费生成。
- 另一组百炼 mock 暴露 `RequestSpec.BasePathSuffixes` 在拆包后未被 URL 拼接器执行，导致原生端点误落到 `/compatible-mode` 下；已按 Manifest 明确声明移除末尾兼容路径，保留自定义代理前缀，并通过同步创建/匿名下载、区域代理路径与视频轮询恢复三个专项测试。全量 Windows app 测试中大量 SQLite 用例因 `CGO_ENABLED=0` stub 失败，不能解释成业务回归；后续需隔离 Linux CGO 执行。
- 最新 `git status` 仍有 55 个未合并路径，其中核心 `web/src/pages/canvas/project.tsx` 仍有 31 个冲突块。画布页当前同时含旧 `CanvasAssistantPanel` 与上游 `CanvasCloudAgentPanel` 的冲突分支，尚未完成切换；不能构建、提交或部署。
- 生成任务服务合成本机即梦前台/批次生命周期与上游 Cloud Agent 的准备任务 API；本机即梦仍不允许后台提交，百炼视频合同校验与上游元数据追踪均保留。`web/src/services/api/generation-task.ts` 已无标记、Babel TypeScript 语法解析通过，但未进行完整类型/运行测试。
- 项目 API 的五处冲突按上游 `http`（仍复用唯一 `apiClient`）及当前 handler 采用 camelCase 查询参数解决；`web/src/services/api/projects.ts` 语法解析通过。显式未合并路径降至 53。画布主文件仅合并了主题 scope 和导入，仍有大量冲突。
- 数据库专题文档已合并官方 Cloud Agent/画布版本/渠道成本与本地目录支持状态、旧 Schema9 时长谱系及 Comfy Bridge，标为当前 Schema31；未改已发布迁移校验和。显式未合并路径降至 52。
- `use-config-store.ts` 恢复 Comfy Bridge 的类型、默认值和持久配置归一化，同时恢复本机即梦目录注入与本机模型值编解码，保留上游渠道排序等字段。Babel TypeScript 语法解析与暂存差异检查通过，完整 Web 类型/运行验证尚无。显式未合并路径降至 51。
- 隐藏自动合并复核发现 `lib/plugins/builtin/workflows.ts` 保留了 Comfy Manifest 却丢失 `COMFYUI_PLUGIN_ID` 定义与 provider 联合类型，已修复。`generation-workflow-execution.ts` 恢复 Comfy Bridge 的安全选择和 `bridgeId`，`generation-task.ts` 将 Bridge 独立配置送后端且不继承普通渠道凭据。设置首页恢复本机即梦、ComfyUI 配置入口并保留上游 Agent 记忆页；个人渠道页恢复 localhost 精确放行控件和 Comfy 入口。相关 TS/TSX 文件语法解析通过，显式未合并路径降至 49。
- Web `bun install --frozen-lockfile` 仍不可用，但对既有根依赖设置 `NODE_PATH=D:\open-ai-canvas\web\node_modules` 后，Bun 能启动候选工作树目标测试。`bun test test/generation-task-local.test.ts` 首次阻断是测试引用上游已删 `use-canvas-agent-operations` 源文件，需迁到 Cloud Agent/operation-history 合同，不应复活旧 hook 当兼容桩。该测试目前 0 pass/1 fail，不能记为生成任务回归已通过。
- `NODE_PATH` 对 ESM 的 `react` 解析仍不充分；候选空 `web/node_modules` 与现有根依赖之间的 Junction 和提权符号链接均被该文件系统拒绝，尝试后已恢复候选空目录、未改根依赖。停止重复链接尝试；完整 Web 测试/类型/构建仍是门禁。
- 画布主文件 `project.tsx` 的显式冲突已合并完并通过 Babel TSX 语法解析：采用上游 Cloud Agent 面板与画布版本历史/覆盖层布局；保留本地肖像排查状态同步、入口、弹窗和候选入画布功能。旧 `/sessions` 的 `canvas-assistant-panel.tsx` 与 `canvas-agent-session.ts` 已精确删除，独立 `canvas-agent/` 进程和 Codex 插件未删除。`CanvasNodeActionContext` 补回肖像排查动作。画布工作模式切换移入顶栏更多菜单，导演台入口同在更多菜单，避免拉长底栏；`canvas-toolbar.tsx` 用上游覆盖层定位。以上只证明静态合并与语法，不是 UI 验收。
- 此时显式未合并路径 46，另有较多无冲突标记却丢失本地符号的自动合并文件；例如工作流插件 ID、`generation-workflow-execution` 已修。候选 TS 7 CLI 在所有冲突文件未清前只报告冲突标记语法错误，无法形成完整类型证明。
- 随后恢复 `types/canvas.ts` 的 Comfy 工作流 ID 与肖像排查节点状态，`canvas-workflow.ts` 重新识别 Comfy provider，`canvas-project-generation.ts` 保留 Comfy 参数映射并把配置实际放回返回值（原自动合并遗漏）。上游 Canvas Toolbar 覆盖层、素材页未同步本地回退与文件上传文本类型已合成；后台价格档沿用上游面板，同时保留仅百炼视频可选的实际时长计费与时长/参考图匹配条件。
- 创作页采用上游拆分后的 `CreationAgentEntry`/`creation-workspace.tsx` UI，移除旧单文件 2K 行组件副本；本地百炼高级参数、首尾帧选择、生成声音/水印控件已迁到新 Composer，`CreationSettings` 保留相应数据，空提示词视频在符合百炼合同且有参考素材时仍可提交。TSX 语法解析及差异检查通过，未做运行验证。显式未合并路径 39。
- 后续画布节点提示词面板改用上游新 `quoteModel`/Modal 调整，同时保留本地无价规格的浏览器预校验、可选空视频提示词和生成专属参数。配置节点恢复 ComfyUI 选择、校验与执行配置；节点工具条/覆盖层采用上游新菜单、批量表/转换与分割选择器。画布上传和媒体工具 hook 的自动合并曾留下重复状态机或语法错误，核对上游完整实现已覆盖本地同名能力与时间线入轨后，按文件采用上游完整实现，保留 `ComfyUI` 和肖像节点的差异在节点操作/连线控制器中精确补回。`use-canvas-project-lifecycle` 采用上游 CAS/历史恢复逻辑。静态语法检查通过，显式未合并路径降至 30；尚无完整 Web 类型检查或 Edge 运行证明。
- 2026-09-25 隔离 `docker build --target backend-test`（Linux CGO）全量未通过：`internal/app` 多个测试失败且全包测试在 10 分钟默认超时中止；包括旧管理渠道测试被当前私网 SSRF 门禁拒绝、Gemini 普通文本路径、模型价格/目录可用性、错误分类、声明式协议 MiniMax 任务 ID、Cloud Agent mock 的本机地址，以及 `internal/repository` 超额时长账单测试。不能用 Windows CGO=0 的 SQLite stub 失败解释这些 Linux 失败；需逐个分辨陈旧测试夹具与真实合并回归，绝不通过开放私网来让测试变绿。Docker 构建只使用候选源码和临时构建层，没有部署/挂生产卷。
- 随后按上游新字段合同合并技能库四个页面和 `skill-runtime`；旧 Web 本地 Agent 面板、聊天 UI 与面板壳没有当前应用调用点且依赖已移除 hook，按用户选择的新 Agent 路线删除；独立 `canvas-agent/` 与 Codex 插件仍在。后台渠道页保留上游排序/复制并补回严格本机渠道开关和 API 格式，未放宽 SSRF。
- `globals.css` 60 处冲突采用上游新创作、素材、画布布局，保留仍被导演台/React Aria 使用的滚动条、滑块和弹层动画；删除原误保留的全局管理端覆盖，因为上游已在 `pages/admin/theme/admin-tokens.css`/`admin-chrome.css` 隔离主题。`admin-ui.css` 价格档颜色换成既有 token 并保留实际时长计费控件。CSS PostCSS 解析通过，未做 Edge 视觉验收。
- 2026-09-25 当前所有显式 Git 未合并路径已清零，`git diff --cached --check` 通过；只有 `.claude/memory/git-conflict-resolution-strategy.md` 的冲突标记教学示例会被 `git grep` 命中，并非未解合并。Web 旧测试中的失效 Agent API 正在清理/迁移，Docker `web-build` 验证运行中。仍无 commit/push/deploy，不能宣称集成完成。
- 隔离 Docker `web-build` 通过（15,275 模块）；完整 `bun run typecheck` 首次暴露隐藏自动合并错误：浮动 Dock Switch 溢出菜单被当命令、旧素材宽松修复函数残留、桌面本机渠道 feature flag 类型缺失、画布旧 Agent/旧分割器引用、Comfy 工作流测试台缺 provider、钱包旧 API 调用、旧 Agent store 的已删类型。逐项修正后重建并复跑完整 TypeScript 类型检查通过。该 Web 镜像仅为本地验证，未部署。
- 冲突相关 Web 专项首轮 101 项中 96 通过、5 失败：一个旧钱包页菜单遗留已删；两个测试断言仍指向旧静态类名/API 合同；模型报价跨能力误发请求已补 fail-closed；创作轨道 CSS 选择器测试改为实际作用域。修正后待复跑，不能将首轮记录为通过。

1. 完成 Batch 6 本机渠道适配及 Backend 编译，复核之前 Batch 2–5 的隐藏自动合并错误。
2. 按 Batch 7–30 逐模块解决所有冲突；生成包只从最终 Manifest/生成器再生，不手工合成 ZIP。
3. 在隔离环境执行 Backend Linux CGO、Web 类型/测试/构建、Canvas Agent 测试/构建、插件包/Compose/文档验证；任何失败只修相关回归。
4. 完成最终 `git status`、`git diff --check`、未合并路径和冲突标记、功能吸收清单检查后向用户报告。保持 merge-open，不提交。
