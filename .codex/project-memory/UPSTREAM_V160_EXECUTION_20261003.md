# 固定 v1.6.0 A″ 增量合并执行记录

更新时间：2026-10-03 11:45 CST。此文件是进行中流水；每次重要状态变化追加时间、命令或 SHA，并将过时的“尚未”标为历史，不用计划代替完成证据。

## 用户授权与固定目标

- 用户选择：从尚未完成 Edge 验收的 v1.5.8 A′ 已暂存合并结果制作**仅本地、明确标记 WIP** 的检查点；保留原 A′ 和旧 A，从检查点新建隔离 A″，增量合入官方最新稳定版本，逐项解冲突、验证并交付 Edge 手验指导。
- 2026-10-03 GitHub Release API 返回 `v1.6.0`、`draft=false`、`prerelease=false`；`git ls-remote official refs/tags/v1.6.0` 为 `42f14cd73eec7c2ded98eb15b6a4f5aa2f3f3b4a`。`official/main=24cba7650df53ef148eb434e960295371b717a53` 是移动主线，排除在固定目标外。
- 启动前 v1.5.8 A′ `.local/worktrees/upstream-v158`：`HEAD=56e9b124d7cfda3c4cb2d8a38f9da477ce14f8f5`、`MERGE_HEAD=258422757c98467f24381219e90f4809394ec226`，未合并路径 0、`git diff --cached --check` 0、未暂存/未跟踪 0。暂存路径 874 个，属于上轮完整候选。旧 A `.local/worktrees/dashscope-plugins@6b16d25e` 干净。
- 本轮不授权正式提交/推送 GitHub、生产部署、真实付费生成或删除任何旧分支/工作树/数据卷。现有 `127.0.0.1:3001` v1.5.8 预览和生产运行边界需在容器操作前重新核实，不得把上一轮健康结果当本轮现状。

## 待记录的里程碑

1. A′ WIP 树、双父 SHA、A″ 分支与工作树路径。
2. 固定标签 fetch/对象确认、增量文件集与显式/自动重叠冲突数量、分批策略。
3. 每批冲突决策、必要调用点、生成物处理与专项验证。
4. 全量门禁、双谱系隔离数据验证、候选镜像/预览（仅在必要且用户授权范围内）、未覆盖项。
5. 用户 Microsoft Edge 手工验收步骤、停止条件、晋升/远端/清理门禁。

## 2026-10-03 12:07 CST：WIP 固化与真实合并启动

- 在 v1.5.8 A′ 中只暂存启动记忆，`git diff --cached --check` 通过、未暂存/未跟踪/未合并路径均为 0；`git write-tree=24bbbc465cc0edd3cd40d84832344360e9fb4266`。`git commit-tree` 产生**仅本地 WIP 检查点** `7419e2b3db6aa08e47a208529c0f27173ccf39c3`，双父为 `56e9b124` 和固定官方 v1.5.8 `25842275`。A′ 原 `HEAD`/`MERGE_HEAD` 保持不变，不是正式晋升提交。
- 官方固定 `v1.6.0` 标签抓取后，本地对象类型为 commit，精确 SHA `42f14cd73eec7c2ded98eb15b6a4f5aa2f3f3b4a`；与 WIP 的 merge-base 是 v1.5.8 `25842275`。新分支 `codex/upstream-v1.6.0-20261003` 位于 `D:\open-ai-canvas\.local\worktrees\upstream-v160`，创建时 clean、HEAD 为 `7419e2b3`。旧 A/A′ 工作树未改动。
- `git diff --name-only 25842275..v1.6.0` 为上游 1,233 路径；`25842275..7419e2b3` 本地 709 路径；路径交集 95。`git merge-tree --write-tree --messages --name-only HEAD v1.6.0` 预判 40 个 content 冲突；随后在 A″ `git merge --no-commit --no-ff v1.6.0` 实际也是 40 个未合并路径，`HEAD=7419e2b3`、`MERGE_HEAD=42f14cd7`。未正式提交、推送或部署。
- 冲突分组：构建/版本 2（`.github/workflows/quality.yml`、`VERSION`）；Backend 启动/数据库/协议/业务/仓储 17（包括 `migrations.go`、`manifest.go`、任务/价格/资源）；Web 画布/模型/管理/服务 19；文档 2。实际冲突均为文本，但不能因同为 textual 就自动选边。先处理合约/迁移与低耦合叶子，后处理任务/画布/跨层调用，最后测试与文档；插件归档和自动合并重叠路径另做语义审计。
- `v1.5.8` 旧候选的 2,380 Web 测试、Linux CGO 与 Schema42 克隆结果只属于上一版本，不能记为 v1.6.0 A″ 验证。用户 Edge、正式提交/推送、生产替换、付费调用仍 NO-GO。

## 2026-10-03 19:11 CST：Backend 第一轮冲突合成（暂存，待全量）

- 40 个显式冲突中的 Backend 17 个与版本/CI 2 个已移除标记并暂存，当前仍有 Web 19、文档 2 个未合并。`VERSION` 合成 `v1.6.0+dingxin.1`；CI 同时保留 CGO/30 分钟依赖用例门禁和上游 Pi Agent npm 安装。
- 上游把 `admin.go`、`analytics.go`、`finance.go`、`channel_models.go`、`model_capability.go`、`model_router.go`、`provider_protocol.go`、`resource.go` 拆为细分文件。逐项对照函数集合后，删旧文件中的重复函数块，把本地差异迁到新文件：渠道目录补全/并发保护及资格过滤、Qwen/Wan/HappyHorse 能力、实际视频时长账单、匿名结果下载、严格任务错误与资源错误脱敏；同时保留上游音频实际时长、Pi Agent、资源运行策略、RequestAware 音频流解析。不是一律 ours/theirs。
- 迁移计划以官方 v1.6.0 的 41–43 为新库规范尾段，再接本地 44–45；已存在的本地 v1.5.8 41/42 或旧 v9 记录走各自验证后的替代顺序，历史 name/checksum 不重写。目标 Schema45，新增三谱系计划测试；**真实隔离克隆升级尚未运行，不可宣称迁移通过**。
- 协议首轮编译发现上游新 `manifest_adapter.go`、`manifest_request.go` 与本地 v4 方法重复；去重后在上游请求文件补变量与 base-path-suffix，在音频流解析入口保留本地严格响应检查。`GOCACHE` 指到工作区后 `go test -count=1 -timeout 10m ./internal/protocol` 在 Windows 通过；第一次默认 Go cache 被沙箱拒绝，随后数据库专项因 Windows `CGO_ENABLED=0` 的 go-sqlite3 stub 失败，需 Linux CGO 隔离重跑。
- Backend `go test -run '^$' ./internal/app ./internal/repository ./cmd/server` 正在使用工作区 Go cache 编译检查；待结果记录。生成物归档尚未重打包/核对；Backend 所有自动合并路径仍需语义检查，不能因编译通过视为业务验收。

## 2026-10-03 23:10 CST：40 个显式冲突清零，验证仍进行中

- Backend 编译检查实际通过：工作区 `GOCACHE` 下 `go test -run '^$' -timeout 10m ./internal/app ./internal/repository ./cmd/server` 三包均成功，运行测试数为 0；协议包 `go test -count=1 -timeout 10m ./internal/protocol` 真测通过。Windows 数据库专项由于 SQLite CGO stub 失败，必须在 Linux CGO 隔离环境补跑，不是数据库代码结论。
- Web 19 个显式冲突已分模块处理：工具/底栏与缩略图保留上游变体，项目画布同时保留本地肖像节点和上游角色资产规则；目录价格支持本地百炼实际视频时长与上游音频按秒；生成元数据把本地显式参数与上游默认/命名空间合成。提示词面板、导演台、素材库、模型能力和同步存储的大块冲突是上游拆文件，逐函数比较后去掉重复旧块，只把本地独有的生成选项、即梦模型编码与本地媒体存储识别迁进新文件。
- 文档两冲突已解决：数据库专题写候选 Schema45 的官方/本地/旧生产三谱系，不宣称迁移已运行；待测试页保留历史 A′ 与官方 v1.6.0 项，并新增 A″ 未完成门禁。`git diff --name-only --diff-filter=U` 实测为 0，所有 40 个显式冲突清零；这不等于自动合并语义或全量验证通过。
- A″ `web/node_modules` 起初不存在，已开始 `bun install --frozen-lockfile`；截至本时间该进程仍运行且占用 CPU，尚无退出状态。Web 类型、测试和构建尚未开始。全量 Backend/Linux CGO、Canvas Agent、三谱系隔离迁移、插件包、候选镜像/HTTP 与 Edge 均未通过；禁止晋升、推送或生产替换。

## 2026-10-04 00:42 CST：全量 Web 与 Linux 数据库专项结果

- Web `bun install --frozen-lockfile` 最终安装 1,440 包且锁文件未更改；`bun run typecheck`、`bun run lint` 均通过。12 文件专项初轮 167/168，唯一失败是即梦任务阶段标签已从 `canvas-node-content.tsx` 拆到 `canvas-node-status-content.tsx` 的旧静态断言；补检查真实新接线后该文件 49/49 通过。全量初轮 2,403 通过、3 失败、1 导入错误：两个静态断言（固定 v1.5.8 版本、模板选择器换行）已同步；同时发现自动合并让“放大提示词”弹窗重复渲染，保留可调整大小且携带模型合同的一份并新增只渲染一份断言。余下一项是隔离工作树 Canvas Agent 缺少 `zod`；其 `bun install --frozen-lockfile` 在沙箱临时目录 EPERM 后获批准重试并完成。最终 `bun run test:all` 主组 2,406/2,406、隔离组 5/5，退出码 0。Web 构建初次被默认 Go cache 权限拒绝，已用工作区 `GOCACHE` 重跑，**目前仍在执行**。
- Docker 只读核对：生产 Backend/Web 仍以旧 `:local` 镜像运行、均 healthy；2026-09-29 的 v1.5.8 `:3001` 预览容器当前不在运行列表。未替换生产、未挂载真实卷。
- 新工作树 Windows 路径在 Docker Desktop bind mount 映射下被 daemon 判为不存在；改从 A″ 发送构建上下文，独立构建 `open-ai-canvas-backend:v160-a-double-prime-test-runtime`，未覆盖生产标签。该 Linux CGO 镜像里 `go test -count=1 -timeout 10m ./internal/database` 真测通过；第一次 `sh -lc` 重置 PATH 找不到 Go，改用 `/usr/local/go/bin/go` 直接运行后通过。
- Backend 首轮无网络全仓 Go 测试未通过：`internal/assets` 的公网源站判定需要 DNS，`internal/app` 至少有同类出站域名夹具失败；`internal/hostupdate` 测试镜像缺根层 Compose 夹具。未放松 SSRF。测试专用 Dockerfile 阶段已补复制 `docker-compose.deploy.yml`，独立测试镜像重建中。App 其余失败需在隔离容器中用 JSON 失败摘要逐项归因；全量不能标绿。SQLite/官方/本地三谱系真实克隆还未执行。

## 2026-10-04 02:20 CST：Web/Agent 构建闭合、Backend 七项回归定向修复

- Web 冻结依赖安装后，`bun run test:all` 最终 2,406+5 全过；`bun run typecheck`、`bun run lint` 通过；`bun run build` 在工作区 `GOCACHE` 下完成三平台 Bridge 构建、TS 与 Vite 15,467 模块，约 4 分 48 秒，退出码 0。首轮默认 Go cache 被沙箱拒绝属环境故障，重跑结果才是通过证据。构建有大 chunk 与已有动态/静态同模块导入警告，未声称性能已验收。
- Canvas Agent `bun install --frozen-lockfile` 在沙箱临时目录 EPERM 后以主机权限核实 156 安装/228 包、无锁改动；`npm test` 首轮 `tsx` 初始化因沙箱 `os.userInfo` ENOMEM 未跑用例，主机权限重跑 328/328 通过；`npm run build` 通过。未操作本机实际 Runtime 账号或付费渠道。
- 独立 Backend Linux CGO 测试镜像 `open-ai-canvas-backend:v160-a-double-prime-test-runtime` 已构建；Docker Desktop 对新工作树 bind mount 报源路径不存在，因此使用客户端 build context 和无宿主端口、无卷调试容器，不绕过生产数据边界。`internal/database`、`internal/protocol` 真测通过；无网络首轮资源/Host Updater 失败分别由公网 DNS 测试条件和镜像缺根 Compose 夹具解释。测试阶段补 `COPY docker-compose.deploy.yml /src/` 后重建，隔离 bridge 网络下 `internal/assets`、`internal/hostupdate` 均通过。
- App 真测先定位七项：配置检查分配 PTIER ID、图片请求规格与陈旧档位 ID、隐藏错误标记、豆包语音资源拒绝分类、网络拒绝文案、音频多结果恢复、七牛配置。将价格档 ID 分配仅留在正式保存路径；订单从请求意图重选档位而保留上游音频时长数量；错误只输出固定安全类别与不含地址的 `connection refused`；语音业务错误输出固定可行动文案，不回传上游资源 ID。音频测试夹具使用新版显式 desktop loopback 能力而不启用旧全私网开关。七项定向真测现 7/7 通过，七牛单跑也通过；Backend 全仓第二轮正在独立容器中执行，尚未宣布全量通过。
- 当前没有生产替换、正式提交/push、三谱系真实克隆或用户 Edge 验收。2026-10-04 只读 Docker 显示生产 Backend/Web 仍为旧 `:local` 且 healthy；历史 v1.5.8 `:3001` 预览容器当前不在运行列表。

## 2026-10-04 03:28 CST：三谱系隔离迁移真实通过，全仓终轮仍在跑

- 只读确认来源测试卷未被运行容器挂载：`open-ai-canvas-catalog-preflight-20260913`（旧 Schema9）、`open-ai-canvas-v158-a-prime-legacy-clone-20260928`（本地 Schema42）。分别只读复制到新卷 `open-ai-canvas-v160-legacy9-clone-20261004` 和 `open-ai-canvas-v160-local42-clone-20261004`，复制时 `diff -qr` 逐字节一致。未接触 `open-ai-canvas_backend-data`。候选 `migrate-schema status` 升级前分别报告 9/45 和 42/45。
- 官方正式 `ghcr.io/ddcat-ai/open-ai-canvas-backend:1.6.0` 拉取 digest `sha256:0312f1185338656185a68ac64eea312ea86649fe11b31c8a336155eb570a9a1f`，仅在新空卷 `open-ai-canvas-v160-official43-clone-20261004` 生成并验证官方 43/43；只读复制到升级卷 `open-ai-canvas-v160-official43-upgrade-20261004`，复制时逐字节一致。该谱系为空业务结构源，不能替代有业务的旧谱系。
- 三个新升级卷分别执行候选 `migrate-schema up`、第二次 `up`、`verify`，每次均为 `current=45 expected=45 ready=true`。SQLite `integrity_check=ok`、`foreign_key_check` 违规 0、迁移记录各 45。旧 9 与本地 42 两份业务计数均为 users/channels/channel_models/tasks/resources/assets/projects/orders = `1/2/370/3/5/5/2/3`，与来源已登记计数保持；官方空业务升级卷这些表均为 0。旧 9/本地 42 的 `.settings-key` 与 `plugin_registry.json` 和各自只读源卷 `cmp` 字节相同。专用 `open-ai-canvas-v160-migrate-gocache-20261004` 仅为 Go 编译缓存，不是数据恢复点。
- Backend 全仓第二轮只剩一条 Cloud Agent 画布修订 CAS 竞争，单跑可过但重复 10 次偶发。原因是本地旧测试手工推进运行时，上游 v1.6.0 新审批 waiter 同时推进同一画布；**只在旧测试辅助函数里关闭该 waiter**，保留生产 CAS/异步实现，相关三个测试连续五轮通过，首个画布补丁测试十轮通过。最终 Backend 全仓在无生产挂载的调试容器重跑并把日志保存在容器 `/tmp`，**截至本时间仍运行，结果待查**。
- Web 2,406+5 全量、typecheck、lint、生产构建及 Canvas Agent 328/328、构建此前已通过；不能由这些结果推断 Edge 交互、付费模型或真实设备成功。A″ 仍为打开的无提交 merge，旧 A/A′ 和生产服务未改。

## 2026-10-04 03:44 CST：Backend 全仓与依赖集成绿灯

- Linux CGO 完整命令 `/usr/local/go/bin/go test -count=1 -timeout 30m ./...` 在独立无卷测试容器退出码 0；`internal/app` 约 304.5 秒，其他所有包均 `ok` 或无测试文件。日志仅保存在临时容器 `/tmp`，未写生产数据；此前丢失一次交互句柄的结果未计作证据，重新运行并读取明确包摘要后才认定通过。
- 依赖专项另在专用无宿主端口 PostgreSQL 16 测试容器与调试容器内安装的临时 Redis 8.8 可执行文件下运行：`TestPostgresAssetIDMigration`、`TestRepairAssetBytesPostgres`、`TestCanvasRevision/postgres`、`TestRuntimeCoordinatorRedisLeasePreservesLongestTTL`、`TestClearAdminRuntimeCacheUsesRedisWhitelist` 均实际 RUN/PASS，非静默 SKIP；`CANVAS_TEST_POSTGRES_DSN` 只指向专用空测试库。生产数据库和 n8n PostgreSQL 均未连接。
- 两个本地关键 v4 插件 `dashscope-image`、`dashscope-video` 的归档分别仅含 Manifest/README/接口文档 3 条目；逐条 SHA-256 与源目录字节相同，`apiVersion=yingce.plugin/v4`。未跑会重写全部支付包的生成脚本；Backend 候选镜像构建中的 payment-package-smoke 阶段已通过，完整镜像仍在构建。
- 截至本时间 A″ staged tree 暂存快照 `917414bc47833307908f2fdbb85fda7efa9a943d`，随后文档记忆写入会改变完整树 SHA；候选 Backend 构建标识临时使用 `tree-917414bc4783`，不是正式 commit。`git diff --cached --check` 剩 4 个上游原样文件的尾随空格（1 个 Pi runtime 源、3 个内置技能卡片），`git diff 7419e2b3..v1.6.0 --check` 已确认它们来自官方增量；未为这次合并大范围格式化官方内容。正式提交后仍需真实 CI 格式门禁。

## 2026-10-04 03:52 CST：候选镜像独立构建进度

- Backend 完整候选镜像 `open-ai-canvas-backend:v160-a-double-prime-candidate-20261004` 已构建，`docker image inspect` 确认 ID `sha256:ac4f54cfd33a05041d8df7d3d4eb092987e7d1e3d33b502709c40b48f47acdb0`；payment-package-smoke 阶段通过。构建标识是暂存树前缀 `tree-917414bc4783`，不是正式 Git commit。
- Web 候选标签 `open-ai-canvas-web:v160-a-double-prime-candidate-20261004` 构建前经 inspect 确认为不存在；现在独立构建中，Docker 已完成三平台 Comfy Bridge 阶段并在 Vite 阶段转换 15,467 模块，**尚无最终退出码，不得标记已通过**。该独立标签不覆盖生产 `:local`。
- 生产/隔离 HTTP 冒烟、最终 Compose 配置核对、Edge 验收与正式提交/推送仍未完成；生产服务及数据卷未被本轮镜像构建替换。

## 2026-10-04 04:09 CST：候选镜像与无宿主端口 HTTP 冒烟通过

- Web 完整候选镜像 `open-ai-canvas-web:v160-a-double-prime-candidate-20261004` 构建退出码 0，inspect ID `sha256:c4217b79b5774587a9c18a5f6c1c41136f7e382f2a54488a21c6c2238ee3124d`；Backend ID `sha256:ac4f54cfd33a05041d8df7d3d4eb092987e7d1e3d33b502709c40b48f47acdb0`。两者的正式版本参数 `v1.6.0+dingxin.1`，构建标识 `tree-917414bc4783` 不是正式 Git commit。Vite 有大 chunk 和静/动态同模块导入非阻断警告，浏览器性能尚未验收。
- 独立 `open-ai-canvas-v160-smoke-net-20261004` 网络中运行新前后端，无宿主端口和挂载。容器均 healthy；从 Web 容器读取 `/build-info.json` 为 `v1.6.0+dingxin.1`、`tree-917414bc4783`，经 Nginx `/api/health/ready` 报同标识、`ready=true`、`schema.current=45`/`expected=45`，首页 HTML 可取。不是登录态/付费/Edge 验证。
- `docker compose config --quiet` 对 `docker-compose.dev.yml`、`docker-compose.local.yml` 直接通过；部署及 deploy+build 首轮因缺强制 `POSTGRES_PASSWORD` 等环境变量按预期 fail-close，随后在一次性 PowerShell 进程提供**仅解析用占位值**后均通过。无服务启动、无真实密钥输出。
- 精确 inspect 四个本轮测试容器的镜像/挂载后，停止并删除两只烟测容器、App 调试容器和 PostgreSQL 测试容器，再删除两条专用网络；没有删除候选镜像、迁移克隆卷、PostgreSQL 匿名测试卷或任何生产资源。临时容器的 `/tmp/v160-backend-full.log` 随容器删除，已在前述日志中记录全仓结果。
- 新交接 `HANDOFF_V160_20261004.md` 记录 Edge 手验顺序。v1.6.0 未开放回环端口，因此此刻无可用 Edge 预览地址；不复用 v1.5.8 的 `:3001`。正式 commit/push/生产替换/旧 A 清理/付费仍 NO-GO。

## 2026-10-04 04:18 CST：最终 Git/旧工作树/生产只读复核

- A″ `HEAD=7419e2b3`、`MERGE_HEAD=42f14cd7`，暂存 1,246 路径，未合并/未暂存/未跟踪均 0。此时 `git write-tree=3f76000e0542df754311a268aee285e612b354af`；写入本段并暂存后完整树 SHA 会改变，后续须现场重算。应用路径无冲突标记；`.claude` 中 Git 教程的字面示例不属于未解冲突。
- `git diff --cached --check` 有 4 项由官方 v1.6.0 原样引入的尾随空格（Pi runtime 一项、内置卡片三项）。为了保持镜像与已验证源码一致，没有在候选镜像构建后编辑这些上游字节；不把此检查说成通过，正式提交与远端 CI 格式门禁仍待完成。
- 旧 A 工作树 clean；原 A′ 仍 `HEAD=56e9b124`、`MERGE_HEAD=25842275`，未暂存/未合并均 0，未被本轮晋升或删除。生产 `open-ai-canvas-backend-1`/`open-ai-canvas-web-1` 均继续运行 `:local` 且 healthy，Web 仍占用宿主 3000；没有本轮临时容器留在运行列表。项目记忆及待测试专题已更新至本轮候选结果。
