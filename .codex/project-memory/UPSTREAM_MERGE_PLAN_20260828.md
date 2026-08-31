# 2026-08-28 官方上游增量合并执行计划

## 状态

- 当前阶段：尾差阶段 B 已完成；停在 `4ba9694` 新尾差审计或固定 `c8b60ce` 完整阶段 C 验证的用户决策门禁。
- 阶段 4-8 已解决全部 29 个文本冲突，当前 `git diff --name-only --diff-filter=U` 为空。
- 本次批准仅包括先固化固定快照 `ab89c05`，再重新只读确认官方 `main` 后同步已审计的轻量提交 `115e228`；若官方已出现新的未审计代码增量，则停止并重新分析。
- push 仍需后续单独远程写入批准。
- 本计划不授权 push、PR、CI、tag、Release、付费调用、数据库恢复或卷删除。

## 固定起点

- 本地起始分支：`main`。
- 本地起始提交：`405d25aec95b307144e5a3e84170bebb948ebe31`。
- 阶段 0-2 完成后，项目记忆与备份登记以一个纯文档提交固化；阶段 3 应从届时干净的当前 `main` 创建集成分支。功能差异和 29 个冲突仍以 `405d25a` 对 `ab89c05` 的代码关系为分析基线。
- 当前应用版本：`v1.1.4`，由根 `VERSION` 在 Vite 构建时注入。
- 阶段 2 实时固定官方快照：`main@ab89c05362394623d00c95a4899be03e5243e75a`。它比旧分析的 `913cf4b` 多 1 个提交；后续阶段开始前仍需确认该 SHA 未被用户要求刷新。
- 用户 fork：`https://github.com/dingxin150015-commits/open-ai-canvas`。计划作为后续统一维护入口；当前不 push、不创建 PR。
- 用户已进一步确认：该 fork 是本次集成分支未来的线上归属。阶段 3 仍只做本地合并；待冲突解决和验证达到后续门禁并获得单独远程写入批准后，再推送同名集成分支，稳定 `origin/main` 不接收未验证代码。

## 核心策略

1. 保留本地经过阶段 1-16 验证的历史和运行合同，在独立集成分支引入官方主线。
2. 官方新协议/插件/工作流/存储/诊断架构作为新骨架；本地百炼支持状态、精确 Adapter、精确 SKU、Create 明确参数、安全日志和 OSS 物化必须保留。
3. 29 个已知文本冲突逐文件解决；自动合并文件继续做语义审查。
4. `ready` 必须绑定精确执行合同；官方 discovery 不自动升级可执行状态。
5. 免费价格必须显式配置；默认零值继续 fail closed。
6. Provider 大冲突按模块重构，不逐行混拼；精确 modelKey 优先于协议通用 fallback。
7. 数据迁移只先在恢复副本预演；live 更新前重新建立 pre-deploy verified/protected 恢复点。
8. UI 管理和人工验收统一使用已登录 Microsoft Edge。

## 阶段门禁

### 阶段 0：冻结基线

- 记录 Git HEAD、状态、分支、remote、标签、版本和哈希。
- 记录运行容器、镜像、卷、健康、重启/OOM、Compose 配置和端口。
- 记录数据库完整性、外键、表数、脱敏业务计数、模型支持状态和文件权限。
- 确认无其他写入任务；发现并发变更则重新冻结。

### 阶段 1：恢复点

- 源码：bundle/refs/status/remote/版本/哈希。
- 数据：SQLite 一致性快照、WAL/SHM、匹配 `.settings-key`、迁移标记和独立恢复副本。
- 验证：integrity、foreign key、表数、业务计数、OSS 密文可解密、逐文件哈希、ACL。
- 镜像：给当前 Backend/Web 增加明确 pre-upstream-merge 回滚标签。
- C 盘存敏感材料，D 盘只保存脱敏指针；禁止删除或覆盖旧备份。

### 阶段 2：远程与固定 SHA

- `official` 指向 `ddcat-ai/open-ai-canvas`。
- `origin` 计划指向用户 fork；现有主分支跟踪关系在 push 授权前不自动改成可写发布路径。
- 本机临时 `upstream` 改名为 `legacy-upstream-snapshot`。
- 获取官方 main/feature/tags 和用户 fork refs，记录实际 SHA；不 push、不 prune 用户数据。

阶段 2 已确认：共同基点 `70a6640`，本地独有 38、官方独有 49；双方分别改动 180/504 个文件，重叠 54 个；`merge-tree` 显式冲突仍为 29 个文件。官方 `feature` 独有 0、落后 main 427。用户 fork `origin/main=11931d0`，本地 main 领先 479 个提交且尚未 push。

### 阶段 3 以后

- 阶段 3 已获批准：从干净本地 main 创建 `codex/upstream-20260828-ab89c05`，以 `--no-commit --no-ff` 引入固定官方 `ab89c05`，只核对真实冲突，不解决业务冲突、不 push。
- 阶段 4 及之后每个阶段独立汇报、独立批准。

### 阶段 7 已确认结果

- 分支仍为 `codex/upstream-20260828-ab89c05`，`HEAD=03604df`，`MERGE_HEAD=ab89c05`；无 merge commit、无 push。
- Web 冲突已全部解决，动态协议以 Backend Registry 为事实源；显式免费价格、系统渠道统一报价、Create 参考轨道/输出开关/按用户模型操作持久化、本地素材缺失 fail-closed 和 Logo 的 Vite/Bun 双运行时边界均保留。
- `bun run typecheck` 通过；拆分定向测试共 41 项通过（核心 22、渠道目录 14、任务中心 2、用户同步 3），0 项失败。
- 大测试文件合并运行曾长时间无输出；拆分到相关用例后通过。`bun install --frozen-lockfile --offline` 因长时间无进展被终止，未联网、未改锁文件；当前依赖仍足以通过类型检查和上述测试，完整安装/全量测试留在后续总门禁复核。
- Backend/Web 运行容器保持原镜像、healthy、RestartCount=0；未重建、未部署、未操作数据卷。

### 阶段 8 已确认结果

- CI 工作流取并集：保留 Backend CGO、Web format/typecheck/test/build、Canvas Agent test/build，并接入官方 Director P0 Chrome E2E；本地未运行 Chrome，也未触发远程 Actions。
- `.gitignore` 保留本地缓存、敏感材料和临时报告边界，并接纳官方 `.superdesign` 忽略项。
- `CHANGELOG.md` 合并双方 `Unreleased`，移除未经本地发布的 `v1.2.0-preview.1` 分段；根 `VERSION` 恢复并保持 `v1.1.4`。
- 数据库文档合并支持状态、显式免费价格、插件状态、ComfyUI Bridge 和 S3/删除 Outbox；待测文档合并本地阶段历史与官方 33 类新增功能，形成 39 个无重复二级章节。
- 插件 README 保留签名 Local Runtime 与 URL Token 拒绝合同，同时接纳 `canvas-context`、`canvas-editing`、`asset-aware-generation` 技能说明。删除了官方旧的 URL Token 手动排查方案。
- CI YAML 解析和必需 job/step 检查通过；Prettier、冲突标记、cached diff check 通过。插件/技能合同首次因 Windows 路径分隔符断言失败，改用 `path.basename` 后 8/8 通过。
- Backend/Web 运行容器保持原 `open-ai-canvas-*:local` 镜像、healthy、RestartCount=0；未重建、未部署、未操作数据卷。
- 当前没有未解决文本冲突，但仍没有 merge commit。下一阶段建议先执行完整 Backend/Web/Canvas Agent/文档与插件门禁，再由用户单独批准是否提交和推送到其 fork 集成分支。

### 阶段 9 已确认结果

- Backend 隔离 Linux CGO 全量首次暴露 5 个合并后门禁缺口：测试镜像未复制 AutoDL 插件制品；两个日志测试仍期待原始上游/预检错误；三个 SKU 测试夹具缺少完整显式价格合同。修复测试镜像输入、保留脱敏断言，并补齐 `PriceConfigured + BillingMode` 后，第三次全量通过。
- Web 依赖按 `bun.lock` 恢复且锁文件未漂移。主套件首次 1061/1062，仅节点注册表测试错误地禁止插件额外节点；改为验证所有内置节点恰好有定义且全表 type 唯一后，专项 19/19、全量 1062/1062 和跨 Runtime 1/1 通过。
- Web TypeScript 通过。生产构建首次因沙箱拒绝 Go 缓存失败；主机权限完成 Windows、Linux AMD64、Linux ARM64 三个 Comfy Bridge 制品后，使用预构建校验完成 13,482 模块 Vite 生产构建，构建耗时 8 分 36 秒。大 chunk 和插件耗时提示为警告，不是失败。
- Canvas Agent 依赖按锁文件恢复且锁文件未漂移；主机权限全量 327/327 通过，包含 Windows 真实进程树终止、Dreamina 调度/恢复/围栏、Local Runtime、插件与技能合同；`tsc -p tsconfig.json` 通过。Windows 无符号链接权限的拒绝子项仍由 Linux CI 保留。
- CI/Changelog/MDX/插件 README/相邻测试 Prettier 通过；gofmt、冲突标记、工作树与暂存区 diff check 通过。五份独立 Compose 与 deploy+build 叠加配置共六组只读解析通过。
- `docs/` 仍没有独立 `package.json`，因此本阶段只有 MDX 格式和结构验证，不能宣称文档站构建通过。
- Backend/Web 运行容器保持 `open-ai-canvas-*:local`、healthy、RestartCount=0、OOMKilled=false；测试镜像和构建产物没有部署，数据卷未操作。
- 当前 `HEAD=03604df`、`MERGE_HEAD=ab89c05`、`VERSION=v1.1.4`、未解决冲突 0；尚未创建 merge commit、push、PR 或触发远程 CI。用户现已单独批准创建本地 merge commit；推送到 `origin/codex/upstream-20260828-ab89c05` 仍需再单独批准。

### 阶段 10：两步可审计合并（第一步完成，第二步按门禁停止）

1. 修正并暂存项目记忆、`BACKLOG.md` 和历史价格说明，根 `VERSION` 保持 `v1.1.4`。
2. 为已经完成全量验证的固定官方快照 `ab89c05` 创建本地 merge commit。
3. 使用只读远程查询重新确认官方 `main`；仅当它仍为已审计的 `115e228` 时，fetch 并创建第二个小型 merge commit。
4. 第二步预期只包含 `README.md` 与两张赞助商图片；逐文件检查后再提交。
5. 最后记录两个 merge commit 的精确 SHA、父提交、官方包含关系和剩余门禁；本阶段结束后停下汇报。

永久边界不变：不 push、不创建 PR、不触发远程 CI、不部署、不操作数据库/卷、不调用 Provider。

### 阶段 10 执行结果

- 第一个本地 merge commit 已创建：`d04c4d1d968ad2c53de919985596d61c8d841f5f`；父提交精确为 `03604df78cc8beffd1d5c7a721d635774ce2e8bc` 与 `ab89c05362394623d00c95a4899be03e5243e75a`。
- 提交后工作树干净；根 `VERSION` 保持 `v1.1.4`，没有 push、PR、远程 CI、部署、数据库/卷或 Provider 操作。
- 2026-08-29 重新读取官方 `main` 得到 `80d2aa688278a7cb9e2222842b24fcbf2f24de59`，不再是已审计的 `115e228129d895a56de9ad0999195f48a5605d7f`。
- GitHub Compare 确认 `115e228..80d2aa6` 为线性向前的 11 个提交，涉及 96 个文件、6,195 行新增和 526 行删除；包含视频批量取帧、生成结果命名、视频设置、短剧分镜工作流、媒体引用/模型匹配、公告、存储管理、定价与品牌素材等代码和数据结构变化，不属于“仅 README 与两张图片”的小同步。
- 因此没有 fetch/merge 第二步，也没有把 `115e228` 单独纳入后伪称已跟上实时主干。后续应以 `80d2aa6` 为新固定 SHA，重新做增量、语义冲突和验证范围分析，经用户批准后再继续。

### 阶段 11 新增量审计结果

- 用户批准审计后，官方实时 `main` 已继续前进并最终固定为 `4f07daae9ec3b4e4cb0a8cd35a6e0c1a4b593f29`；审计结束前再次只读确认未漂移。
- `ab89c05..4f07daa` 共 14 个提交、107 个文件、6,843 行新增和 661 行删除；当前本地与官方新增量重叠 22 个文件。
- 只读 `merge-tree` 模拟得到 6 个显式冲突：`finance.go`、`pending-test.mdx`、`web/package.json`、`video-settings-panel.tsx`、`model-capabilities.ts`、`create/index.tsx`。
- 本阶段未执行真实 merge、冲突编辑、merge commit、push、PR、远程 CI、部署、数据库/卷操作或 Provider 调用。
- 功能分组、逐文件解决合同、8 类语义复核和阶段 12 验证门禁见 `UPSTREAM_INCREMENT_AUDIT_20260829.md`。

### 阶段 12A 执行结果

- 开始前再次只读确认官方实时 `main` 仍为固定候选 `4f07daae9ec3b4e4cb0a8cd35a6e0c1a4b593f29`。
- 工作树干净、`VERSION=v1.1.4`、无进行中合并；私有恢复点 `BKP-20260828-164620-UPSTREAM-PREMERGE` 目录仍存在，数据库、WAL、独立恢复库、`.settings-key`、迁移标记和 Git bundle 六项 SHA-256 均与登记值一致。
- 已执行 `git merge --no-commit --no-ff 4f07daa`。当前 `HEAD=ae58ecf`、`MERGE_HEAD=4f07daa`，没有 merge commit。
- 真实冲突集合与模拟完全一致：6 个文件、8 个冲突区块，无意外新增或缺失；107 个官方增量路径进入合并索引，状态计数为新增 26、修改 75、未合并 6。
- 本阶段没有解决冲突、运行源码测试、push、PR、远程 CI、部署、数据库/卷操作或 Provider 调用。下一步必须等待用户批准阶段 12B。

### 阶段 12B 执行结果

- 六个显式冲突、八个冲突区块均已逐文件手工解决，未使用批量 ours/theirs；`git diff --name-only --diff-filter=U` 为空。
- `finance.go` 保留服务端按完整 `ModelRequestIntent` 重新选择价格档，直接任务、逻辑模型和线路切换均使用服务端选出的 tier ID；显式旧 ID 与 intent 不一致时 fail-closed。新增两个计费防回归测试。
- `model-capabilities.ts` 保留 Qwen/Wan/HappyHorse、视觉总数、智能时长和图片尺寸三态，同时接入插件工作流能力与画幅可为空；视频设置取官方只读尺寸推导并保留智能时长/声音/水印。
- Create 保留按用户/模型/操作的设置持久化与输出开关；无画幅模型不显示画幅，并把旧全局 `size` 清为空值，模型兼容检查也不再提交该选项。`web/package.json` 保留 panic guard/跨 Runtime 拆分并纳入两个官方新增测试。
- 自动合并区域已做源码语义复核：SupportStatus 与默认/规格价格并存；图片+音频继续受参考音频上限约束；资源新旧匿名路由、签名、Range 和脱敏日志并存；AutoDL 包/Manifest 合同、短剧归属、公告引用删除和安全 Markdown 边界均保留。
- 验证：Web 9 个相关测试文件 69/69、`bun run typecheck` 通过；Backend 不依赖 SQLite 的计费/能力专项通过；AutoDL 插件制品专项通过。
- Backend 宿主专项首次受 Go 缓存 ACL 阻断，主机权限启动后只有 `CGO_ENABLED=0` 的 `go-sqlite3` stub 失败；这属于已知环境限制，不是业务断言失败。数据库、短剧、公告、资源删除和完整 Backend 仍必须在阶段 12C 使用隔离 Linux CGO 门禁验证。
- 当前没有 merge commit、push、PR、远程 CI、部署、运行数据库/卷或 Provider 操作；下一步必须等待阶段 12C 批准。
- 结束时官方实时 `main` 已从固定候选 `4f07daa` 前进到 `2f6832f`。GitHub Compare 显示仅 1 个提交，但涉及短剧工作台、技能运行、任务恢复、Backend/Web/Canvas Agent 共 113 文件、7,016 行新增和 723 行删除；本轮没有 fetch/merge 该尾差。当前阶段 12C 仍只验证已解决的固定 `4f07daa` 合并，`2f6832f` 必须另行审计和批准。

### 阶段 12C 执行结果

- 隔离 Linux CGO Backend `go test -count=1 ./...` 全部包通过；测试镜像 manifest 为 `sha256:867a29ee02893ab284625580f15adcf734d81a5fad4619676f7dbc45cd355c03`，未挂载运行数据、备份、密钥或数据卷。
- Web panic guard 主套件 1082/1082、跨 Runtime 1/1、TypeScript 和生产构建通过；Vite 转换 13,493 个模块，构建耗时 8 分 36 秒，仅有大 chunk 与插件耗时警告。
- Canvas Agent 在主机权限下 327/327 通过，覆盖真实 Windows 进程树终止、Dreamina 围栏/调度/恢复、Local Runtime、插件和技能合同；Windows 无符号链接权限项按既定边界由 Linux CI 保留。`npm run build` 通过。
- Prettier 首次检查发现 36 个官方新增/修改文本文件格式漂移；仅做机械格式化后，全部受支持暂存文本文件检查通过，Web TypeScript、相关 9 文件 69/69、AutoDL 制品专项和 13,493 模块生产构建复验通过；最终构建耗时 7 分 32 秒。全部 staged Go 文件 gofmt 通过。
- 六份 Compose 单独配置和 deploy+build 叠加配置解析通过。`server` 首次在缺失必填数据库密码时按设计 fail-closed；使用仅供 `config --quiet` 的虚拟值后全部通过，未启动服务或连接数据库。
- `git diff --cached --check`、冲突标记扫描、待测文档二级标题去重和 JSON/MDX 格式通过；两张赞助商 PNG 可读取并完成视觉检查，README 使用固定宽度引用。
- Web 生产构建在 Prettier 机械修正前首次通过；修正后又完整复跑并通过 13,493 模块生产构建，最终暂存字节已获得构建证据。
- 结束时官方实时 `main=0893741`。相对固定目标 `4f07daa` 新增 2 个提交、122 文件、8,431 行新增和 863 行删除；本轮未 fetch/merge，未来必须独立审计。
- 当前没有 merge commit、push、PR、远程 CI、部署、运行数据库/卷或 Provider 操作。下一步只允许在用户新批准后创建本地 merge commit；push 继续是独立门禁。

### 阶段 12D 执行结果

- 用户批准后创建本地 merge commit `e2326856b1cc3a854ea11fcd54c5c3e2866c8ac3`，提交说明为 `chore(upstream): 官方主线 - 合并 4f07daa 新增量`。
- 双父提交精确为 `ae58ecfe1ebadcfd8e915559081983422516e58b` 与 `4f07daae9ec3b4e4cb0a8cd35a6e0c1a4b593f29`；`git merge-base --is-ancestor 4f07daa HEAD` 通过。
- 提交后 `.git/MERGE_HEAD` 不存在、工作树干净、根 `VERSION` 仍为 `v1.1.4`。
- 官方实时 `main` 复核仍为 `08937417ec8e11c20c77a7effeaf026787577650`；未 fetch/merge 该尾差。
- 没有 push、PR、远程 CI、部署、运行数据库/卷或 Provider 操作。下一步若推送集成分支到用户 fork，必须获得新的远程写入批准；是否先审计 `0893741` 也需用户决策。

### 用户 fork 远程检查点

- 用户明确授权“下一步”后，按既定顺序仅推送已验证集成分支，没有更新 `origin/main` 或创建 PR。
- `git push --set-upstream origin codex/upstream-20260828-ab89c05` 成功；首次远程 SHA 为 `265c9970840b8be755c14dda23e697b83bc1cfa0`，与当时本地 HEAD 一致，本地分支开始跟踪 `origin/codex/upstream-20260828-ab89c05`。
- 只读复核 `origin/main=11931d0085ccfec3952a4a7e30bd482173c51492`，同名 head 的 PR 查询返回空数组。
- 本项目记忆提交会在同一次已授权远程检查点操作中继续同步到该集成分支；最终远程 SHA 以 `git ls-remote` 实时结果为准。
- 未创建 PR、未触发部署、未操作运行数据、未调用 Provider，也未 fetch/merge 官方 `0893741` 尾差。

### `c8b60ce` 尾差审计

- 用户同意下一步后，只执行独立尾差审计。官方实时 `main` 从旧记忆的 `0893741` 前进到 `c8b60cea078f8a7093d6186f89febb6c3a10d9b4`，审计结束前再次确认未漂移。
- `4f07daa..c8b60ce` 共 4 个提交、164 个文件、26,364 行新增和 2,696 行删除；当前分支与尾差重叠 28 个文件。
- 只读 merge-tree 模拟得到 17 个显式冲突文件、48 个冲突区块；另有 11 个自动合并但必须语义复核的文件。
- 核心风险是技能包/自动 GitHub 同步与启动迁移、短剧工作台大重构、任务恢复/真实进度、积分日志，以及约 8,770 行后台 CSS 与本地设计系统合同冲突。
- 完整功能、逐文件解决合同、安全审计和后续阶段见 `UPSTREAM_TAIL_AUDIT_20260830.md`。
- 本阶段没有真实 merge、业务代码修改、push、PR、远程 CI、部署、数据库/卷、技能同步或 Provider 操作。

### 尾差阶段 A 执行结果

- 开始前复核本地工作树干净、`VERSION=v1.1.4`，本地 `HEAD=3510997` 比 fork 检查点领先 1 个纯审计文档提交；官方实时 `main` 三次确认均为固定 `c8b60ce`。
- 私有恢复点目录仍存在，数据库、WAL、独立恢复库、`.settings-key`、迁移标记和 Git bundle 六项 SHA-256 与登记值一致。
- 已执行 `git merge --no-commit --no-ff c8b60ce`。当前 `HEAD=3510997`、`MERGE_HEAD=c8b60ce`，没有 merge commit。
- 真实冲突集合与 merge-tree 审计完全一致：17 个文件、48 个区块，无新增或缺失；164 个官方尾差路径进入索引，状态为新增 36、修改 111、未合并 17。
- 本阶段没有解决冲突、运行源码测试、push、PR、远程 CI、部署、数据库/卷、技能同步或 Provider 操作。下一步必须等待尾差阶段 B 批准。

### 尾差阶段 B 执行结果

- 17 个显式冲突文件、48 个区块均已解决，未解决索引项归零。除公告编辑器和分页项目详情使用经过逐段审阅的官方完整重构作为基线外，短剧 workbench、Schema、技能、Create、后台和 CSS 均逐区块取并集；未做批量 ours/theirs。
- Schema 同时保留本地 Catalog 回填、官方章节字数回填和价格选择器迁移；Canvas Agent skills 测试取并集；声音样本使用真实 MIME/扩展名；Create 保留设置/输出/引用合同并接入 Skill Runtime。
- 短剧 workbench 接入分页资产、技能、任务恢复、真实进度、镜头删除/解绑和三阶段状态，同时保留本地镜头画面/动态提示词字段、Backend 任务生命周期和 fail-closed 生成门禁。
- Local Runtime 保持 30MB 签名 JSON 上限；后台回归测试改为经 panic guard 运行；后台模型编辑器 raw color 改为 admin token，并新增 CSS scope/panic guard 回归测试。
- 验证：Web TypeScript 通过；11 文件专项首轮 92/93，唯一失败是删除按钮测试对白空白敏感，修正后该文件 2/2、短剧关键 9/9 复验通过；后台回归 10/10。Canvas Agent skills/Runtime 12/12，Backend 技能包解析/Provider 进度纯逻辑专项通过。完整 SQLite/Backend/Web/Canvas Agent/构建仍留阶段 C。
- 格式、gofmt、cached diff、冲突标记和待测文档标题检查通过；当前 165 个业务/官方路径暂存，项目记忆另行暂存。
- 结束复核发现官方实时 `main=4ba9694`。`c8b60ce..4ba9694` 为 8 个提交、42 文件、+2,398/-414，包含多项直接修复分镜模型计价、并行生成、角色图声资产、视频参考协议和镜头台词；本轮未 fetch/merge。
- 没有 merge commit、push、PR、远程 CI、部署、数据库/卷、技能同步或 Provider 操作。建议用户先决定是否审计 `4ba9694` 修复尾差，再进入阶段 C。

### 修复尾差独立审计与阶段 B2 建议

- `c8b60ce..4ba9694` 已完成独立审计；固定目标通过独立 `official/audit-4ba9694` 引用获取，本地 `official/main` 和当前 `MERGE_HEAD=c8b60ce` 未改变。
- 当前已解决索引与该尾差有 11 个重叠文件；merge-tree 预测 3 个显式冲突文件、15 个区块，其中工作台 13 个。
- 前 5 个短剧/协议修复应作为一个耦合组在阶段 C 前纳入，但先修复 xAI 尾帧在前端、Backend Provider 和协议 Adapter 的三路径不一致。
- 浅色画布修复必须 token 化。首页仪表盘默认暂缓；若纳入，先修 `taskCenterEnabled`、查询错误态和 300 条统计截断。`4ba9694` 的工作台横向滚动 CSS 可与首页部分拆开纳入。
- 下一步定义为修复尾差阶段 B2：只在用户明确批准后打开固定 `4ba9694` 的 no-commit 增量、解决 3 文件/15 区块并做必要修正；完成后再进入完整阶段 C。merge commit、push、PR、CI、部署和 Provider 调用仍为独立门禁。

### 修复尾差阶段 B2 执行结果

- 用户批准后，先以安全树和双父安全提交 `82758b2` 保护当前 `c8b60ce` 已解决索引，再标准重开固定 `4ba9694` no-commit 合并。直接冲突 18 文件/52 区块，经安全三方结果恢复后精确收敛为 3 文件/15 区块并逐区块解决。
- 当前 `HEAD=3510997`、`MERGE_HEAD=4ba9694`、未解决索引项 0、非暂存 0、`VERSION=v1.1.4`。安全引用 `refs/codex/safety/b2-pre-4ba9694` 暂时保留到后续 merge commit 门禁完成。
- xAI 尾帧在协议 Adapter、Backend 类型化 Provider 和前端直连统一 fail-closed；R2V 忽略陈旧帧元数据。浅色画布背景 token 化；首页仪表盘暂缓，工作台横向滚动 CSS 保留。
- B2 专项验证通过 TypeScript、尾差 Prettier、xAI 协议/Service、gofmt/JSON/marker/diff。Web Bun 和 Linux CGO/全量/构建/Edge 未完成，转入阶段 C。
- 官方实时 `main=fb089b2`，相对固定目标的新尾差未审计、未 fetch/merge。没有 commit、push、PR、CI、部署、数据或 Provider 操作。

### 固定 `4ba9694` 阶段 C 执行结果

- Backend Linux CGO `go test -count=1 ./...` 全部通过；Web 最终 admin 10/10、主套件 1122/1122、跨 Runtime 1/1、TypeScript 和 13,469 模块生产构建通过；Canvas Agent 最终字节 328/328、build通过。
- 首轮 Web 3 个失败均为测试合同漂移，只修正 xAI 测试签名和两个空白/新调用合同断言；专项23/23后完整套件通过。测试镜像补齐 Agent锁文件依赖后跨 Runtime通过。
- 145 暂存文本 Prettier、暂存 Go gofmt、JSON、文档标题和七组 Compose解析通过；运行容器/数据未改变，没有 Edge 候选验收。
- 当前固定合并满足创建本地 merge commit的源码门禁，但提交仍需单独批准。提交时必须验证双父 `3510997` / `4ba9694`、最终树和版本；不得包含或宣称已包含官方实时 `e124a9e`。
- 当前未提交状态不可作为完整合并推送。merge commit完成后，用户 fork 同名集成分支 `9815464` 可被 fast-forward，但 push、PR、CI、`origin/main` 和部署继续是独立门禁。

### 固定 `4ba9694` 本地 merge commit 授权

- 用户已单独授权本轮创建本地 merge commit；提交说明使用项目约定格式，双父必须为提交前 `HEAD=3510997` 与固定官方 `4ba9694`，提交树必须等于最终暂存树。
- 本轮授权不包含 push。提交完成后只读复核 fork、运行部署和官方实时主干并评估推送可行性；不更新 `origin/main`，不处理官方新尾差。

### 固定 `4ba9694` 本地 merge commit 执行与推送评审

- 本地 merge commit已创建：`d0def5f807e7d8c410103c85bd5ded43506b39b3`，双父 `3510997` / `4ba9694`，tree `49580c569687e0ef30b47cce0f5d1aa8e9aa9ff7`。未 push。
- fork集成分支 `9815464` 到新提交可 fast-forward，dry-run通过；实际push仍需新授权。非main push不会触发当前quality/publish工作流，若需要远程质量证据必须另行设计PR或手工工作流门禁。
- 运行部署仍是阶段16 `f9317fa` 镜像，不应直接把新提交视为已部署；数据库/启动迁移、Edge和回滚门禁继续独立。
- 官方实时已到 `b7348ab` / `v1.2.3.1`；新尾差未合并，未来审计基线为22提交、216文件、模拟29冲突文件/52区块。

### 用户 fork 固定 `4ba9694` 检查点执行结果

- 5份提交后记忆先固化为 `7687053`。推送门禁确认远程 `9815464` 未漂移、可fast-forward且dry-run通过。
- 已将同名集成分支fast-forward到 `7687053`；独立远程复核一致，ahead/behind 0/0。没有force、`origin/main`、PR、CI或部署动作。
- 最终检查点记忆由本次小型docs提交继续同步；完成后再次核对远程与本地HEAD一致并保持工作树干净。

## 永久 NO-GO

- `git reset --hard`、宽范围 checkout、批量 ours/theirs。
- 在当前 `main` 直接 pull/merge。
- `docker compose down -v`、删除或替换数据卷。
- 未验证恢复点前迁移、重建或部署。
- 将大型 fork 整体合入。
- 未重新授权的真实模型、OSS 删除、GitHub 写入、代理或账户变更。
