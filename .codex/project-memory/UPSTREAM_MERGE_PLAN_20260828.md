# 2026-08-28 官方上游增量合并执行计划

## 状态

- 当前阶段：阶段 11 新增量审计已完成；等待用户批准以固定 `4f07daa` 执行 no-commit 合并和冲突处理。
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

## 永久 NO-GO

- `git reset --hard`、宽范围 checkout、批量 ours/theirs。
- 在当前 `main` 直接 pull/merge。
- `docker compose down -v`、删除或替换数据卷。
- 未验证恢复点前迁移、重建或部署。
- 将大型 fork 整体合入。
- 未重新授权的真实模型、OSS 删除、GitHub 写入、代理或账户变更。
