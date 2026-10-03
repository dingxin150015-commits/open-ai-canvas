# v1.5.8 A′ 执行记录（进行中，NO-GO）

## 已确认检查点

- 2026-09-27 15:27:36 CST：仅本地双父 WIP 提交 `56e9b124d7cfda3c4cb2d8a38f9da477ce14f8f5` 固化未提交 v1.5.6 候选；父提交为本地 `6b16d25e` 与官方 `473d535e`。旧 A `codex/dashscope-system-plugins-20260911` 未改，原 v1.5.6 合并工作树未改。
- 2026-09-27 下午：新隔离 A′ `.local/worktrees/upstream-v158`、分支 `codex/upstream-v1.5.8-20260927` 从 WIP 检查点建立；对固定官方 `v1.5.8@25842275` 执行 `git merge --no-commit --no-ff`，实际显式文本冲突 59。`MERGE_HEAD` 打开，未正式提交、推送、部署。
- 截至 2026-09-27 17:22 CST：显式未解决冲突从 59 降至 38。已按职责合成数据库/协议/后端错误 9 个文件，前端 API/持久化 5 个文件和画布媒体节点 7 个文件；所有这些文件已暂存为冲突解决，不含冲突标记。
- 数据库：官方 Schema 30–40 保留，本地时长迁移顺延至 41、Comfy Bridge 至 42；旧生产 Schema9 历史时长迁移仍以原校验和验证，官方 9–40 迁移在该谱系顺移一位且避免时长重复执行。`schema_test.go` 增加两谱系编号连续性断言；尚未在 Linux CGO 或数据克隆上运行实际迁移，不可声称安全通过。
- 插件：Manifest 同时接收本地 v3/v4 声明式协议与上游 SMS 贡献校验。后端响应保留 `requestId/errorCode/errorCategory/retryable`，同时加入上游 `details`；本地权限/SSRF 边界未放宽。
- 后端初次 `go test ./internal/database ./internal/protocol ./internal/handler ./internal/app -run '^$'` 暴露自动合并的 `model_catalog.go` 失效变量及 `resource.go` 缺少日志导入；两处按双方调用链修复后，同一仅编译命令通过。此命令没有运行测试用例，Linux CGO 全量仍待执行。
- 前端生成持久化冲突保留已提交副作用不被后续普通写入失败误判为未提交，同时纳入上游对未落盘编辑继续普通保存的判断。媒体节点采用上游稳定预览与产出模型信息，保留本地可选空视频提示词、资产标签与人像入口；尚未完成前端类型/运行验证。

## 当前未完成与门禁

- 2026-09-27 晚间：A′ `web/` 执行 `bun install --frozen-lockfile` 成功（Bun 1.4.0，1440 包，仅隔离工作树）；`bun run typecheck` 完整通过；创作配置、画布外观/媒体和生成持久化 5 文件专项为 56/56 通过。专项测试发现画布下载仍用已移除的 `saveAs`，已改接上游 `downloadBrowserMedia`，真实浏览器下载仍待 Edge。
- 完整 `bun run test:all` 执行 2236 个用例/296 文件，2190 通过、46 失败、5 个错误，耗时约 515 秒，退出码 1。失败包括：画布单请求多图输出映射和任务消费、本机渠道权限投影、创作设置用户隔离、API 客户端合同、导演台入口/底栏/深色主题等静态断言、镜像发布脚本、旧本机 Agent 教程、旧任务会话合同。必须逐项分辨真实回归与过时测试；不能宣称全量通过或通过删除测试处理。
- 随后复测关键资产：`bun test test/canvas-native-image-batch.test.ts` 12/12 通过（旧夹具补上资源访问响应与完整模型能力）；本机渠道权限、运行时投影和新任务合同三个文件 12/12 通过。统一 `resolveModelRequestConfig` 运行时精确投影本机许可，测试改用当前 `http` 的共享 Axios adapter；旧 `/sessions` 测试按用户决策改为 `/tasks`，未恢复旧 API。
- A′ 前端生产构建 `bun run build` 已成功（Bun 1.4.0，包含 Comfy Bridge 生成步骤、TypeScript 与 Vite，15,427 模块，约 4 分 42 秒）。该结果只验证构建，不代表运行 UI/Edge 或全量测试通过。
- 后续全量失败归因：`canvas-native-image-batch.test.ts` 旧夹具未模拟官方新增 `/resources/access` 合同，补共享客户端 adapter 和完整模型能力后 12/12 通过；本机渠道三文件测试改当前 `http` 合同并把许可投影收敛到 `resolveModelRequestConfig`，12/12 通过。
- 旧网页 Agent 模块已按用户决策移除，三个遗留测试文件仍 import 其面板/hook。仅退役与已删除旧入口直接绑定的 7+5+6 个用例及相应导入；保留其余任务/素材/持久化测试。新的 Cloud Agent 专项覆盖须在最终审计中核对，不能把退役测试算通过。
- 旧 A 的本机即梦任务中心职责在 v1.5.6 WIP 中被自动合并漏掉，导致 11 个专项失败；已按当前 `http` 客户端恢复任务列表合并/项目过滤、查询/等待/手动刷新、只删除本机记录、提交不确定禁止重试及安全日志，并保留上游媒体恢复。`bun test test/generation-task-local.test.ts` 现 49/49 通过；仍需完整回归。
- 2026-09-28 后续：旧网页 Agent 已删除模块导致 3 个测试文件导入错误。只退役与旧 `/sessions` 面板/hook 直接绑定的测试块，保留其余测试。任务物化文件按当前资源访问合同补假 adapter 后 41/41 通过；`generation-task-local.test.ts` 在恢复本机即梦调用链后 49/49 通过。`user-scoped-generation-persistence.test.ts` 复核了账号切换迟到响应、真实批量删除 API、增量画布最新读取及显式资源修复+版本 CAS 夹具；当前相关专项已逐项通过，尚待全文件与全量重跑。
- `web/src/services/user-data-sync.ts` 在素材分页请求返回后增加 session epoch 与用户 scope 复核，拒绝旧账号迟到结果；无浏览器全局时重置定时器改用 `globalThis.clearTimeout`。这是本轮发现的真实跨账号边界修复，须纳入完整回归。
- Backend `go test ./... -run '^$'` 全仓仅编译通过；真正 Linux CGO 全量未执行。Docker Desktop Engine 当前不可连接；启动 Docker 可能自动启动已有容器，已向用户异步请求单独授权，在收到答复前不得启动。Windows CGO=0 不可替代 SQLite/Linux 完整测试。
- 截至 2026-09-27 20:25 CST：59 个显式 Git 未合并路径均已解决并暂存，`git diff --name-only --diff-filter=U` 与未合并状态计数均为 0。画布相关 14 个 TS/TSX 文件及 `globals.css` 完成语法解析；创作页、后台页面与相关测试文件逐组完成 Babel 语法解析。尚不能以此替代完整类型/运行测试。
- 自动合并复核已修复 `backend/internal/app/model_catalog.go` 的失效 `normalized/err` 引用、`backend/internal/app/resource.go` 的缺失日志导入、`web/src/lib/canvas/canvas-appearance.ts` 深色自定义网格黑底黑点问题，以及创作页上游默认音频/水印与本地显式开关的报价/提交分裂；新增显式选项一致性测试。
- 上游资源分发合同采用浏览器直连 OSS/CDN，七牛设置说明不再承诺旧浏览器媒体代理；本机即梦任务更新/删除、Comfy Bridge、百炼 v4 Manifest 与实际时长计费继续保留。`VERSION` 候选为 `v1.5.8+dingxin.1`，未发布。
- A′ 当前仍处于 `MERGE_HEAD=25842275` 的无提交状态。Windows Bun 单测尝试先被 PowerShell `bun.ps1` 执行策略挡住，改用 `bun.cmd` 后因新工作树无 `react` 依赖而停在导入阶段：0 pass/1 error，不是业务测试通过；后续在隔离依赖环境重跑。
- 历史检查点（17:22 CST）：当时还剩 38 个显式冲突；现已清零。98 个双方重叠但 Git 自动合并的路径仍需语义复核。
- Backend Linux CGO 全量、Web 全量/类型/生产构建、Canvas Agent 测试/构建、插件包再生校验、两谱系脱敏克隆迁移和用户 Edge 手工验收均未完成。不得将当前 A′ 标为稳定，不得部署、推送或删除旧 A。
- `git diff --cached --check` 全树会报告官方新内置技能 Markdown 的历史尾随空格；已解决的目标源码路径专项 `--check` 无误。这是上游新增内容的格式门禁，不等同于本次冲突块错误，最终需按项目质量规则处理或明确分类。

## 2026-09-28 后续实测与修正（本节覆盖上方过时的“未完成”描述）

- Web `bun.cmd run typecheck` 通过；`bun.cmd run test:all` 最终主组 2374/2374、隔离 slot 组 5/5 通过。首轮残余两项导入错误分别由 Windows 文件 URL `.pathname` 和隔离 `canvas-agent/node_modules` 缺失造成；前者改 `fileURLToPath`，后者在隔离工作树补齐依赖，未恢复旧网页 Agent 模块。Canvas Agent 在主机权限下 `bun.cmd install --frozen-lockfile` 校验 156 个安装/228 个包、无变更；`npm.cmd test` 328/328 通过，`npm.cmd run build` 通过。Web `bun.cmd run build` 最终通过，Vite 转换 15,427 模块、用时约 4 分 57 秒；输出包含 Windows amd64、Linux amd64/arm64 三个 Comfy Bridge 文件。源码构建不等于 Docker 镜像构建。
- Windows 主机 `go env CGO_ENABLED=0` 且无 GCC；`go test ./...` 运行失败，失败包报告 `go-sqlite3 requires cgo to work`，不能记作后端业务回归。`go test ./... -run '^$'` 全仓编译通过；不依赖 SQLite 的 `go test ./internal/protocol ./internal/provider/... ./payment-plugins` 通过。Linux CGO 全量与两谱系克隆迁移仍是发布阻断。
- Docker Desktop Engine 管道不存在。启动它可能自动启动已有容器，已向用户异步请求单独确认；未获答复前未启动、未构建镜像、未碰真实卷。
- 静态 Compose 逐份解析发现 4 处自动合并引入的 `environment` 键缩进错误，已修复。`docker compose --env-file .env.example -f ... config --quiet` 的默认、本地、开发编排通过；Server、Deploy 及 Deploy+Build 用仅当前进程的占位密码、数据库地址、镜像 digest 和 HTTPS Origin 通过。未执行 `up/down`。
- `Dockerfile` 增加独立 Go 构建阶段，交叉编译 Windows amd64、Linux amd64/arm64 的 Comfy Bridge 下载产物并复制到 Vite 构建上下文；宿主生成二进制从 `.gitignore`/`.dockerignore` 排除。`web/test/comfy-bridge-release.test.ts` 通过；Docker Web 镜像仍待实际构建验证。
- 两个百炼 v4 插件原归档与其 `manifest.json`、README、接口文档逐字节一致；发现旧 `scripts/build-bailian-plugins.ps1` 会先调用 v2 生成器，可能降级 v4。已让生成器拒绝两个 v4 ID，脚本先验证 v4 Manifest、再归档现有源文件；实际重新打包成功，两个新归档仍均为 3 项、逐字节一致、`yingce.plugin/v4`。未触碰插件安装态或生产注册表。
- 官方内置技能 5 个 Markdown 文件中原有尾随空格已清理；仅改变行尾空格，不改技能正文。最终暂存后的 `git diff --cached --check` 为 0，显式未合并路径为 0。源码/配置目录不存在行首冲突标记；`.claude/memory/git-conflict-resolution-strategy.md` 中的三行标记是历史示例，不是本轮冲突。
- 自动合并重叠路径再次由 Git 双方差集计算为 98；与最初 merge-tree 预览树比较，其中 9 个在本轮已手工补修（`.gitignore`、模型目录、资源日志、响应测试、v2 生成器、画布网格、任务中心、跨用户同步和模型配置），其余 89 个仍与 Git 自动结果一致。已审阅高风险协议插件、Provider、任务 Worker、路由、画布生成、创作/任务 API 的关键本地差异并运行上述专项/全量；受 Linux CGO/Edge 缺失限制，不能把静态复核写成完整运行验收。
- 上游自动合并把 `AGENTS.md` 中本机 Canvas Agent、Codex 插件及其验证说明丢掉，现已恢复；同时把已不存在的 `docs` Next.js 构建描述改为当前 MDX 内容目录。上游 Director E2E 脚本仍显式使用 Chrome，不符合本项目 Microsoft Edge 专用约束，因此本轮没有运行，F6 仍保留隔离 Edge 门禁。
- CI 辅助检查：`gofmt -l .` 为 0；Web `bun.cmd run lint` 通过；Git Bash `plugin-packages/verify-payment-packages.sh --all linux amd64` 返回 0。`bun.cmd scripts/check-changed-formatting.mjs` 在未提交 A′ 上返回 1，原因是脚本按 `HEAD=56e9b124` 的父提交枚举历史 v1.5.6 WIP 中数百个文件，再用当前工作区内容检查，其中含大量旧格式文件；它不是当前未提交 merge 的准确提交差异门禁。正式 A′ 形成提交前须对实际拟提交差异重新运行并处理剩余格式失败，不能记为已通过。
- 对固定官方 `25842275` 的文件集做差：A′ 工作树相对官方没有整文件删除，另有 301 个本地新增路径；这只能排除“官方文件整体遗漏”，不能证明被修改的上游调用链语义正确。

## 2026-09-28 Docker Desktop 启动后的隔离验证（进行中）

- 用户明确告知已启动 Docker Desktop。只读核实 Linux Engine 29.1.3 可用；现有 `open-ai-canvas` Compose 仍指向旧 A `.local/worktrees/dashscope-plugins/docker-compose.local.yml`，生产后端镜像 `sha256:d67ee794`、重启计数 0、挂载 `open-ai-canvas_backend-data:/data`。本轮未对该 Compose、容器或真实卷执行 up/down/restart/write；n8n 和 Firecrawl 也保持原状。
- 首次独立 `docker build --target backend-test` 在真实 Linux CGO 下运行约 13 分钟后失败：管理渠道和 Custom Relay 的旧 loopback 测试夹具触发新桌面本机渠道安全门禁；NewAPI 旧协议名不再可选；Gemini 普通 Create 被旧 guard 拦截；超时长计费人工复核曾可透支；`internal/app` 的默认 10 分钟测试超时发生在大量 SQLite AutoMigrate 的 CPU 路径。没有把这些失败归为 Windows CGO stub 或忽略。
- 已把后端 Docker 测试阶段拆成可复用 `backend-test-runtime`，把单包/全量 Go 测试超时增为 30 分钟；CI `quality.yml` 合并重复 env，显式保留 CGO=1 与依赖检测。独立镜像 `open-ai-canvas-backend:v158-a-prime-test-runtime` 不覆盖生产 tag。
- 针对第一批失败：测试改用显式桌面本机渠道许可或测试专用公网地址+本机 RoundTripper，未放宽生产 SSRF；过时 `newapi-video`/`custom-video` 用例改为当前 `newapi-channel-1/2`；Gemini 普通 Create 与 Agent 各要求已安装相应声明式合同；人工复核时长账单现在检查可用余额。Linux CGO 定向复验 `internal/app`、`internal/handler`、`internal/repository` 全部通过。
- 为不依赖当前宿主机目录 bind mount，创建无挂载/无宿主端口的独立测试容器与网络；隔离 PostgreSQL 16 使用 tmpfs、无持久卷，临时测试容器内安装 redis-server。依赖齐全的首轮全量进一步暴露了公开目录旧夹具无协议、Cloud Agent mock 使用旧 localhost 放行、`ModelError` 缺 Unwrap、调用日志回显上游原文和旧视频 Token 断言等；已做最小修复且相关 `internal/app` 定向测试通过。完整 Linux CGO 全量复跑仍在执行，不能记为通过。
- 已从 2026-09-13 受保护且验证过的 `catalog-clone` 备份（只读）创建新独立卷 `open-ai-canvas-v158-a-prime-legacy-clone-20260928`。备份内 DB、WAL/SHM、`.settings-key`、迁移标记、插件登记与资源目录均存在；复制后 DB 与设置密钥逐字节比对一致。首次复制命令因 PowerShell 提前展开 shell `$()` 失败、目标经只读确认仍为空后才重试成功。未写源备份或生产卷；该新卷尚未运行迁移/完整性验证。
- 依赖齐全的第一次全量还有 8 个 `internal/app` 失败族：公开目录旧夹具遗漏协议/有效价格、Cloud Agent mock 仍依赖 legacy localhost 许可、ModelError 未 Unwrap、日志泄漏上游原文、旧视频 Token 断言等；更新为新安全合同或修复真实错误后，相关定向测试通过。第二次全量再次暴露若干旧本地测试与上游新插件/价格合同的差异，以及任务错误可操作文案被完全泛化的回归。已恢复损坏能力 JSON 的目录隔离和仅有效价格档可用，任务日志保留泛化脱敏但允许白名单类型经 `providerErrorDetail` 过滤后的提示，机器判断使用内部原文而非公开消息；相应定向测试通过。完整 `internal/app` 包与最终全仓正在后续验证，不能提前记为通过。
- 最终合成保留“无活动价格档时已显式配置的标量价格可用；只要存在价格档就绝不回退标量”的本地既有合同，同时保持损坏能力 JSON 失败关闭。媒体恢复发现新桌面回环能力虽用于任务请求，同源结果下载仍走旧通用 SSRF 客户端；改为仅当后端桌面能力启用、渠道显式 `allowLocalChannel` 且同源精确回环时使用拒绝重定向的严格下载策略，缺任一条件不发起请求。专项通过，文档已同步 `task-runtime-stability.mdx` 与 `local-development.mdx`。
- 依赖齐全的 Linux CGO 最终 `go test -count=1 -timeout 30m -json ./...` 全包退出码 0。隔离 PostgreSQL 16、测试容器内 Redis 可执行文件及 ffmpeg 均在场；另逐个以 `-v` 核实 CI 要求的 5 个 PostgreSQL/Redis 用例均实际 RUN/PASS：`TestPostgresAssetIDMigration`、`TestRepairAssetBytesPostgres`、`TestCanvasRevision/postgres`、`TestRuntimeCoordinatorRedisLeasePreservesLongestTTL`、`TestClearAdminRuntimeCacheUsesRedisWhitelist`。未授权的 `TestDDCatLiveImageRoundTrip` 及平台不适用的 `TestRPCProviderRejectsIncompatibleFallbackBinary` 仍按各自条件跳过，不计为付费实测通过。下一步为候选镜像与旧谱系克隆迁移。
- 隔离旧谱系克隆卷经只读 `PRAGMA integrity_check=ok`、`foreign_key_check` 0 项；迁移前 `Schema 9`、81 表，脱敏基线计数 users/channels/channel_models/tasks/resources/assets/canvas_projects/billing_orders = 1/2/370/3/5/5/2/3。未查询任何用户正文。该卷尚未迁移。

## 2026-09-29 隔离迁移、镜像与收尾门禁（进行中）

- Backend 候选镜像 `open-ai-canvas-backend:v158-a-prime-candidate-20260928` 成功构建，ID `sha256:f8b73d259347...`；Web 候选镜像 `open-ai-canvas-web:v158-a-prime-candidate-20260929` 成功构建，ID `sha256:623419094756...`。两者都带 `v1.5.8+dingxin.1` 和明确的未提交来源标识 `tree-f85b2ab8c4d7`，不是正式 commit；未覆盖生产 `:local` 标签。Web 镜像内 Windows x64 Bridge 为 `MZ`、Linux x64/ARM64 为 ELF，机器码分别 `3e 00` / `b7 00`，三个文件均在 `/opt/canvas-release`，`build-info.json` 与候选树标识一致。
- 旧生产谱系：候选 `migrate-schema status` 在新克隆返回 `current=9 expected=42 ready=false`，执行 `up` 后 42/42，第二次 `up` 与 `verify` 均 42/42。迁移后 `integrity_check=ok`、`foreign_key_check` 0 项，103 表，users/channels/channel_models/tasks/resources/assets/canvas_projects/billing_orders 仍为 1/2/370/3/5/5/2/3；`.settings-key`、存储迁移标记、`plugin_registry.json` 与只读备份逐字节一致。首次用只读卷打开迁移后 SQLite 因 WAL 侧文件已清理而无法创建共享内存，改为仅对隔离卷提供可写目录、数据库仍以 `sqlite3 -readonly` 查询后通过；这不是数据损坏。
- 官方谱系：从干净固定官方 `v1.5.8@25842275` 临时 worktree 构建仅包含迁移二进制的独立镜像，在新卷生成官方 40/40，并逐字节复制到另一升级卷；候选二进制识别 40/42、升级 42/42、重复 up 和 verify 均通过。官方源/升级副本分别为 98/100 表，业务关键表皆 0，双方 integrity_check=ok、外键违规 0。该谱系验证真实官方 40→候选 42 的结构和幂等，但使用空业务数据；真实业务保持由旧 Schema9 克隆覆盖。
- 官方临时 worktree `D:\open-ai-canvas\.local\worktrees\official-v158-schema-20260929` 已在核实 clean、HEAD 固定 SHA、路径限定后用非强制 `git worktree remove` 移除；可从官方 SHA 重新恢复。旧 A、A′、原 v1.5.6 候选和上述隔离卷未删除。生产卷及其 Compose 仍未改变。
- 候选 Backend/Web 各在无数据卷、无宿主端口的隔离 Docker 网络中短暂启动：Backend `/api/health/ready` 为 HTTP 200、`ready=true`、Schema42/42、活动任务0、Go1.25.14、`tree-f85b2ab8c4d7`；Web `/build-info.json` 与后端身份一致，Nginx `/api/health/ready` 代理成功，三种 Comfy Bridge 下载 URL 的 HTTP HEAD 均成功。没有用户登录态或真实业务调用。结束后只移除本轮 5 个无生产挂载临时容器与空闲测试网络；3 个隔离迁移卷和候选镜像仍保留。现有影策生产 Backend/Web 镜像分别仍为 `sha256:d67ee794`/`sha256:9bfc9c67`、RestartCount 0，Compose 仍指向旧 A。
- 镜像构建后暂存树相对镜像源码树 `f85b2ab8c4d7bfd2ddd3d6284a1bb42e3b116b7a` 只改变项目记忆与 `docs/content/docs`；Backend/Web/Agent/插件、Dockerfile 和 CI 源码未变。写入本记忆会再次改变完整树 SHA，接管时请以 `git write-tree` 重新计算，不能把历史树值当当前 commit。`git diff --cached --check`、`gofmt -l backend`、未合并路径与业务源码行首冲突标记均为 0；旧 A 工作树干净，A′ 无未暂存/未跟踪文件。镜像构建身份标注的是未提交源码树，不应冒充正式 commit。
- 最终再以共同基点和原 WIP/固定官方重新计算：双方重叠 157 路径、显式冲突 59、自动合并 98；其中 23 路径在 A′ 相对 Git 自动预览树有定向修复或格式变化、其余 75 与自动结果一致。高风险的协议/任务/错误、仓储计费、创作和画布状态已通过相应源码/全量/镜像门禁；自动结果的视觉与真实用户交互仍保留 Edge 门禁，不能因无冲突或测试通过伪称用户已验收。
- 2026-09-29 只读远端复核：`official/main=f9b8c5a43e9c76d21656ec5b514ca7fe26c85e5a`，最新正式 v1.5 系列标签仍是 `v1.5.8^{}`=`258422757c98467f24381219e90f4809394ec226`。本轮已批准的固定主版本目标不变；移动 main 未标记尾差留下一轮独立审计，不在打开的 merge 上叠加。
- 为使最终合并提交的格式门禁只检真实改动，`web/scripts/check-changed-formatting.mjs` 现在识别第二父提交：若文件与第二父完全相同则跳过，真正合成/新增仍检查。新增合并夹具证明上游原样未格式化文件被跳过，而合成文件退化仍失败；专项 10/10 通过。未提交工作树用相同算法只读模拟：382 个 Web 改动候选中 262 个与官方第二父相同、48 个为第一父旧格式、72 个应检查；发现的 24 个真实候选格式问题只对精确文件运行 Prettier（约 234 行增、106 行删的机械排版）。复跑后模拟 72/72 格式通过，Web `bun.cmd run typecheck` 通过，`bun.cmd run test:all` 主组 2375/2375 + 隔离组 5/5，`bun.cmd run build` 通过（15,427 模块、约 2 分 52 秒）。正式提交后仍须运行原始 CI 格式命令，当前模拟不能替代真实提交门禁。
- GitHub 只读查询：设计 Issue #332 已由 `ddcat666` 在 2026-09-07 以 `closed/completed` 关闭，评论 0；不能由此推断五项架构选择。2026-09-28 `git ls-remote` 显示 `v1.5.8^{}` 仍为固定目标 `25842275`，移动的 `official/main` 已到 `f9b8c5a4`；本轮不在打开的 merge 中追逐新尾差，后续 PR/下一轮合并重新审计。中文贡献方案见 `docs/plans/bailian-upstream-pr-v158.md`。

## 2026-09-29 11:37 CST 回环地址隔离预览（用户已批准）

- 启动前只读确认候选 Backend/Web 镜像 ID 为 `sha256:f8b73d259347...`/`sha256:623419094756...`，`3001` 无监听，拟用预览容器/网络/卷名不存在；原生产 Web 仍发布 3000。Docker 访问需沙箱批准，获得本轮批准后只在精确隔离对象上操作。
- 新建网络 `open-ai-canvas-v158-edge-preview-net-20260929`、空卷 `open-ai-canvas-v158-edge-preview-data-20260929`；启动 `open-ai-canvas-v158-preview-backend-20260929`（网络内 `backend` 别名，挂载新卷 `/data`，无宿主端口、自动迁移、回环 Origin、注册测试账号可用、私网上游拒绝、Provider 插件开关关闭）与 `open-ai-canvas-v158-preview-web-20260929`（只发布 `127.0.0.1:3001:3000`）。未连接/挂载生产 Compose 或 `open-ai-canvas_backend-data`。
- 实测预览 `/build-info.json`、`/api/health/ready` 与 `/` 都是 HTTP 200；Web/Backend 构建身份一致，为 `v1.5.8+dingxin.1/tree-f85b2ab8c4d7`，Backend `ready=true`、Schema42/42、活动任务0，容器均 healthy。生产 Backend/Web 仍分别使用镜像 `sha256:d67ee794`/`sha256:9bfc9c67`，均 running/healthy、RestartCount 0。预览保留运行等待用户 Edge 手验；HTTP 检查不证明登录态或 UI 正确。
- 上游百炼贡献提案已把给合并大模型的省 token 操作卡放在第一部分、维护者自然语言说明放在第二部分；未向上游发评论/PR，Issue #332 未给架构答复的门禁仍保留。文档/记忆更新不改变已经构建镜像的应用源码，最终提交仍需真实 SHA 重建与 CI 验证。
