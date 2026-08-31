# 官方主线尾差审计：`4f07daa..c8b60ce`

## 状态与边界

- 审计日期：2026-08-30。
- 当前本地/用户 fork 集成分支：`codex/upstream-20260828-ab89c05@98154647b3458f51909799c333b6088dcb258636`，本地与 `origin` ahead/behind 为 0/0，工作树在审计前干净。
- 已验证并合并的官方固定基线：`4f07daae9ec3b4e4cb0a8cd35a6e0c1a4b593f29`。
- 本轮两次只读确认官方 `main`，最终固定审计候选为 `c8b60cea078f8a7093d6186f89febb6c3a10d9b4`；已 fetch 到本地只读 `official/main`，审计结束前未漂移。
- 用户只批准独立尾差审计；不授权真实 merge、冲突编辑、merge commit、push、PR、远程 CI、部署、数据库/卷操作、技能自动同步或 Provider 调用。

## 增量规模

- 共同基点：`4f07daa`。
- 官方新增：4 个提交、164 个文件、26,364 行新增、2,696 行删除。
- 当前集成分支相对共同基点改动 229 个文件；双方重叠 28 个文件。
- `HEAD...official/main`：本地独有 46、官方独有 4。
- 只读 `git merge-tree --write-tree --messages --name-only HEAD official/main` 得到 17 个显式冲突文件、48 个冲突区块；工作树未进入合并状态。

## 四个新增提交与功能

### `2f6832f`：短剧工作台、技能运行与任务恢复

- 短剧项目增加封面资源、章节字数、镜头/素材索引、分镜生成恢复、任务上下文、项目角色声音样本和资产引用。
- 技能从单一 Markdown 扩展为版本化包：新增 `skill_versions`、`skill_files`，支持 Markdown/ZIP/GitHub 安装、文件浏览/搜索、自动更新、Backend 包存储和 Canvas Agent 原生 skill 文件输入。
- Canvas Agent 支持最多 8 个技能、每包 512 文件、单文件 8MB、总计 20MB；保留目录结构并校验路径/base64。
- Provider 轮询响应可写回真实生成进度；图片/视频连接上游阶段不再显示统一假 35%。

### `0893741`：长篇项目接口拆分和分页

- 项目详情拆为 core、chapters、assets、workbench 等分页/按需接口，避免一次返回长正文和全部素材。
- 新增章节字数回填与复合索引，优化项目、章节、镜头、任务、画布和素材查询。

### `90da729`：管理后台桌面信息架构与视觉规范

- 重构后台导航、页面框架、表格、筛选、设置页和运营页；新增 `admin-ui.css` 及设计/视觉回归文档。
- `admin-ui.css` 约 8,770 行，静态审计发现 39 行十六进制颜色、14 行 `rgb/rgba`、16 个直接 `.ant-*` 选择器和 91 个 `!important`，与本地“三层 token、限定第三方覆盖、避免全局 Ant 补丁”的合同存在明显语义冲突，不能原样整体接受。

### `c8b60ce`：后台配置运营与积分日志

- 完善公告、资源、功能开放、邮件、第三方、绘图、方舟和日志页面；API 调用日志区分用户积分账单与上游估算成本。
- 新增后台 UI 回归测试，并通过 `pretest` 自动运行。

## 数据和启动副作用

- 新表：`skill_versions`、`skill_files`；`skills`、`user_skill_states` 增加版本、来源、哈希、同步和安装字段。
- `projects` 增加 `cover_resource_id`；`project_units` 增加 `word_count`；多张短剧/任务/画布表增加复合索引。
- `MigrateSchema` 增加章节字数回填；必须与本地 `backfillChannelModelCatalogState` 和价格选择器迁移按明确顺序共存。
- Server 启动新增 `EnsureSkillPackages`，会将旧技能物化为 ZIP/版本记录；Worker 每 6 小时检查开启自动更新的 GitHub 技能。部署前必须备份、使用隔离数据库验证幂等/回滚，并明确网络与磁盘写入边界。

## 17 个显式冲突及解决合同

| 文件                                  | 区块 | 解决合同                                                                                                                                                                                                                             |
| ------------------------------------- | ---: | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `README.md`                           |    1 | 保留当前格式化表格与全部既有成员，追加 Rou；图片单独验证，不覆盖赞助商信息。                                                                                                                                                         |
| `backend/internal/database/schema.go` |    2 | `Models()` 取并集；先执行本地 Catalog 状态回填，再执行章节字数回填，再执行价格选择器迁移。两个 helper 均保留并补顺序/幂等测试。                                                                                                      |
| `canvas-agent/test/skills.test.ts`    |    2 | 保留本地 Windows `path.basename`/清理结构和 24K 限制测试，加入官方多文件包、reference 文件、版本和非法路径/base64/20MB/512 文件边界。                                                                                                |
| `backend-database.mdx`                |    1 | 表格取并集：补 `projects`、`project_units`、技能版本表及本地镜头/产物/工作流边界。                                                                                                                                                   |
| `pending-test.mdx`                    |    1 | 保留本地历史和 NO-GO，再加入素材密度、技能包、项目恢复与后台 UI 清单；二级标题去重，自动化通过不升级为运行验收。                                                                                                                     |
| `canvas-node-generation.ts`           |    1 | 采用官方从真实 Resource MIME/ObjectKey 推导声音样本扩展名的异步实现；保留本地资源 URL/storage key、音频容量和角色声音提示合同。                                                                                                      |
| `app-top-nav.tsx`                     |    1 | 接入官方 TopBar Extension，但保留 hideChrome、移动遮罩、折叠侧栏和创建/空间工作台 class；修复嵌套缩进并做响应式验证。                                                                                                                |
| `admin-route-pages.tsx`               |    6 | 采用官方发布门禁、运营动作和任务优先说明；保留当前路由、Suspense、后台引用上下文及本地安全组件。                                                                                                                                     |
| `channels-page.tsx`                   |    2 | 采用官方可访问筛选、凭证状态和确认式行操作；保留本地 URL 状态、完整地址展示和删除/停用安全语义。                                                                                                                                     |
| `admin-announcements-panel.tsx`       |    5 | 以官方安全预览/发布状态机为主，保留 `AnnouncementContent`、配图草稿、置顶、重新发布与焦点恢复；消除重复 preview/lifecycle 结构，不能静默丢字段。                                                                                     |
| `storage-resources-panel.tsx`         |    2 | 采用官方 scoped filter wrapper 与 aria；保留 URL 筛选、用户 ID、批量引用校验/事务删除和危险操作确认。                                                                                                                                |
| `create/index.tsx`                    |    2 | 取并集：保留本地 `creationReferenceMetadata`、设置持久化、声音/水印和无画幅清空；加入官方 `skillRuntime` 与完整技能文件提交。                                                                                                        |
| `projects/detail.tsx`                 |    3 | 采用官方 core/units 分页与 TopBar Extension；保留本地画布创建后的远端关联失败降级、归档提示、路由纠正和用户归属门禁。                                                                                                                |
| `workflow-production-workbench.tsx`   |   14 | 最大风险。加入官方技能选择、任务恢复、真实进度、镜头删除/解绑、资产 mention、分页和三阶段状态；保留本地生成配置、持久任务生命周期、提示词字段、资产版本、失败可重试和 fail-closed 语义。必须逐区块重构，不能选择整文件 ours/theirs。 |
| `workflow-shared.tsx`                 |    1 | 扩展 `ArtifactStatus` 接收任务状态，保留产物 ready/failed/stale 显示；任务成功但产物未同步必须明确显示，不伪造已生成。                                                                                                               |
| `workflow-stage-views.tsx`            |    1 | 采用官方分页统计数据，保留现有语义和可读 JSX；主要是格式冲突，仍需真实空态/加载态验证。                                                                                                                                              |
| `workflow.css`                        |    3 | 合并任务状态、时间线媒体与 1280/1050/760 断点；去重重复规则，保持 token/reduced-motion，禁止把官方大范围 raw color/`!important` 扩散到全局。                                                                                         |

## 11 个自动合并但必须语义复核的重叠文件

1. `backend/cmd/server/main.go`：`EnsureSkillPackages` 是启动写路径，必须隔离验证幂等和失败回滚。
2. `models_channel.go` / `repository/finance.go`：积分日志只按 BillingOrder ID 批量读取，并再次校验 `order.UserID == log.UserID`；保持积分与上游成本分离。
3. `provider.go` / `task_worker.go`：只从成功 create/poll 响应读取 0–100 进度，不保存响应正文；进度写入失败不得让供应商成功请求失败。
4. `service.go`：GitHub 技能自动同步 Worker 只能服务显式 `autoUpdate`，需要关闭/并发/超时证据。
5. `web/package.json`：官方 `pretest` 在 panic guard 外运行；应把 `admin-ui-regressions.test.ts` 纳入本地 `test:main`，不要依赖绕过 guard 的 pretest。
6. `web/src/pages/canvas/project.tsx`：技能/任务上下文必须继续走统一 Backend 任务与用户 scope。
7. `web/src/services/api/projects.ts`：分页接口必须保持 envelope 解包、AbortSignal、用户归属和旧长篇项目安全降级。
8. `web/src/styles/globals.css`：只保留真正全局 token/布局入口；后台具体规则放 scoped `admin-ui.css`，不得新增全局 Ant 补丁。
9. `generation-task-local.test.ts`：保留 Dreamina/local 既有安全合同并加入章节/镜头恢复上下文，不降低原测试覆盖。

## 技能包与安全审计

- 正面证据：上传/解压限制为 20MB、512 文件、单文件 8MB；拒绝软链接、特殊文件、重复/越界/`.git` 路径；GitHub 只接受 `https://github.com`，下载固定 commit 的 codeload ZIP，HTTP client 30 秒超时；原始文件响应带 `nosniff`、sandbox CSP、private/no-store。
- Canvas Agent 再次校验路径、base64、单文件/总量，并要求存在 `SKILL.md`；技能临时目录按 turn 清理。
- 待解决：官方把 Local Runtime JSON body 上限从 30MB 提高到 64MB。20MB 包的 base64 加元数据仍可落在约 30MB 内，建议保留 30MB，除非用测试证明合法最大包确实需要更高上限；不能无依据扩大所有 JSON 路由攻击面。
- 自动 GitHub 同步和 `EnsureSkillPackages` 都会产生网络/磁盘/数据库副作用；源码合并不授权运行迁移或启用自动更新。

## 管理后台设计冲突

- 官方后台重构功能完整，但 `admin-ui.css` 规模和直接覆盖方式与本项目 AGENTS/UI 设计系统约束冲突。
- 推荐先保留官方组件结构和交互，再做一次 scoped CSS 收口：颜色进入 Primitive/Semantic/Component token，所有 Ant 覆盖限定在 `.admin-shell` 或具体组件，合并重复选择器，删除无依据 `!important` 和 raw color。
- 后台改动必须用 Microsoft Edge 验证主要路由、浅色/深色、1280/1440/1920、筛选、Drawer/Modal、键盘焦点和危险操作确认；静态测试不能替代视觉验收。

## 推荐后续阶段

### 尾差阶段 A：固定 SHA 并打开 no-commit 合并

- 开始前再次确认官方实时 SHA；本轮只合并已审计的 `c8b60ce`，新尾差另行记录。
- 确认本地与 fork 同步、工作树干净、`VERSION=v1.1.4` 和恢复点可用。
- 执行 `git merge --no-commit --no-ff c8b60ce`；仅核对真实冲突仍为 17 文件/48 区块，变化即停止。

### 尾差阶段 B：分层解决

1. README、Schema、Canvas Agent skills 测试与两份文档。
2. 技能包 Backend/Agent、安全上限和测试入口。
3. 任务进度、积分日志、分页 API 和 Create/Canvas 任务上下文。
4. 短剧工作台 5 个文件逐区块重构。
5. 管理后台组件结构后做 scoped CSS/token 收口；不能原样接受 8,770 行样式。

### 尾差阶段 C：验证

- Backend：技能包上传/ZIP/GitHub/版本/幂等/自动同步，Schema/字数回填，项目分页/归属，任务恢复/进度，积分日志；隔离 Linux CGO 全量。
- Web：把 admin regression 纳入 panic guard；技能、短剧、分页、恢复、后台安全和现有 Qwen/Wan/HappyHorse/Create 合同；全量、跨 Runtime、TypeScript、最终格式化后生产构建。
- Canvas Agent：多文件 skill、路径/base64/大小、临时目录清理、30MB body 边界；主机全量和构建。
- 文档/配置：gofmt、Prettier、MDX/JSON、冲突标记、Compose；新增 `admin-ui.css` 做选择器/token/重复规则静态审计。
- Browser/运行：源码门禁后另行批准开发预览或候选部署，使用已登录 Microsoft Edge；不回退 Chrome。运行数据库迁移、自动同步和真实技能 GitHub 下载需另行授权。

## 当前结论

- `c8b60ce` 尾差功能价值高，但复杂度显著高于前一轮：17 个冲突文件、48 个区块，且包含启动迁移、网络 Worker、技能执行边界、短剧工作台大重构和后台设计系统冲突。
- 建议继续分阶段处理；当前仅完成审计，不能宣称本地/用户 fork 已跟上官方 `c8b60ce`。
- 官方尾差未修改根 `VERSION`，本地继续保持 `v1.1.4`。

## 尾差阶段 A 实际结果

- 用户批准后，开始前再次确认官方实时 `main=c8b60ce`，并复核既有恢复点六项关键哈希一致。
- 已执行固定 SHA 的 `git merge --no-commit --no-ff c8b60ce`。真实冲突恰好为上文 17 个文件、48 个区块；模拟结论无漂移。
- 当前处于未提交合并状态，尚未解决任何冲突。尾差阶段 B、merge commit 和 push 分别需要后续批准。

## 尾差阶段 B 实际结果

- 2026-08-30，用户批准尾差阶段 B。17 个冲突文件/48 区块全部解决，未解决索引项为 0；关键安全合同按本审计方案保留。
- Local Runtime 保持 30MB；admin regression 进入 panic guard；Schema 三类回填/迁移共存；技能、Create、短剧、任务恢复、后台发布安全和 scoped CSS 完成取并集。
- Web TypeScript、后台/技能/短剧/恢复等专项、Canvas Agent skills/Runtime 12/12、Backend 纯逻辑专项通过。完整数据库和跨栈验证仍未执行。
- 官方结束时已前进到 `4ba9694`；其相对 `c8b60ce` 的 8 个提交中多项直接修复本轮短剧工作台和资产/视频协议，未纳入当前合并。建议将其作为阶段 C 前的独立审计决策。
