# 官方修复尾差独立审计：`c8b60ce..4ba9694`

## 状态与边界

- 审计日期：2026-08-30。
- 固定区间：`c8b60cea078f8a7093d6186f89febb6c3a10d9b4..4ba9694f459808cec8fa57b55f2db767d9f412b8`。
- 线上 `official/main` 在审计时正是 `4ba9694`；固定目标只 fetch 到独立只读引用 `official/audit-4ba9694`，本地 `official/main` 仍锁定 `c8b60ce`。
- 当前合并现场保持 `HEAD=3510997`、`MERGE_HEAD=c8b60ce`、170 个暂存路径、未解决索引项 0、非暂存业务改动 0、`VERSION=v1.1.4`。没有 merge commit、push、PR、远程 CI、部署、数据库/卷、技能同步或 Provider 调用。

## 规模、依赖与功能

- 8 个线性提交、42 个文件、2,398 行新增、414 行删除。
- `a52376d`：项目默认图片/视频模型、工作台能力归一、逻辑模型报价和规格选择。
- `b5f8b66`：生成任务改为后台可靠提交，允许不同镜头并行；提交镜头修订 ID、角色图声资产和任务产物元数据。
- `37eb2b7`：为图片参考增加 ID/Role，统一首帧、尾帧与 `reference_image` 在内置协议、Backend Provider 和前端直连 Provider 中的映射。
- `f2fe163`：工作台同时返回历史绑定资产版本与当前角色声音，生成时使用历史视觉版本、当前声音样本和镜头台词。
- `83781d9`：分镜导入/工作台自动补可编辑资产 mention，模型能力检查不再在前端静默丢弃角色声音样本。
- `b0fa90b`：单独调整浅色画布背景。
- `8353f3c`：把首页改为创作仪表盘，新增项目/画布/素材/任务统计、图表和最近项目表。
- `4ba9694`：移除首页欢迎 kicker，并修复工作台资产名导致的横向滚动；后者属于前 5 个短剧修复的 UI 收口，前者依赖首页重构。

前 5 个提交相互依赖，应作为一个协议/工作台修复组审阅；`b0fa90b` 独立；`8353f3c` 与 `4ba9694` 的首页部分是独立 UI 组。若暂缓首页，可只移植 `4ba9694` 的 `workflow.css` 部分，不能整提交 cherry-pick。

## 对当前已解决索引的精确模拟

- 使用当前 index 的 `git write-tree` 生成树 `eafbe070280d0688cb496060e4760c278be1c5e6`，再创建无引用审计提交 `aaa2e08f3723a4880b81ef5872dda7641725cc67`；没有移动分支或索引。
- 与 `4ba9694` 的 `git merge-tree --write-tree` 结果为 3 个显式冲突文件、15 个冲突区块：
  - `docs/content/docs/backend/backend-database.mdx`：1 区块；取现有完整领域表与新增默认模型字段的并集。
  - `web/package.json`：1 区块；保留 panic guard、`test:admin-ui`、`test:main`、跨 Runtime 拆分，并把新增测试显式加入受保护入口。
  - `web/src/pages/projects/detail/workflow-production-workbench.tsx`：13 区块；必须保留本地删除/解绑、任务恢复、预览历史、字段布局和 fail-closed 合同，同时吸收默认模型/报价、并行提交、历史资产版本、声音/台词和 mention 修复。禁止整文件 ours/theirs。
- 双方共有 11 个重叠文件；除上述 3 个外，8 个自动合并文件仍需语义复核：`builtin.go`、`builtin_test.go`、`provider.go`、`provider_test.go`、`pending-test.mdx`、`workflow.css`、`projects.ts`、`globals.css`。

## 独立代码审计发现

### 阻断项 1：xAI/起始帧协议对“只有尾帧”的解释不一致

- `backend/internal/protocol/builtin.go` 的 xAI Adapter 把单个 `last_frame` 放进通用 `frameImages`，随后作为 `image` 起始图发送；错误信息声称不支持尾帧，但代码没有拒绝该输入。
- `backend/internal/service/provider.go` 的类型化 xAI 路径只检查 `videoStartFrameNodeId`；只有 `videoEndFrameNodeId` 时会走 `reference_images`，把尾帧当角色/风格参考图。
- `web/src/services/api/video-provider-openai.ts` 又会正确拒绝 xAI 尾帧。因此同一用户输入在前端直连、Backend 类型化 Provider 和协议 Adapter 三条路径产生三种结果。
- 影响：可能在付费视频请求中把尾帧误当首帧或普通参考图，输出语义错误且难以从任务状态发现。
- 解决：建立单一 `validate/resolveVideoImageRoles` 纯合同；xAI 对任何 `last_frame` fail-closed，只允许 0/1 个 `first_frame`，`reference_to_video` 才允许 `reference_images`。Backend、协议 Adapter、前端共用等价测试，新增“仅尾帧”“首帧+普通参考”“尾帧+普通参考”用例。

### 阻断项 2：首页任务卡没有遵守功能开关，失败时伪装为零数据

- 首页无条件渲染指向 `/tasks` 的“本周任务”卡；当 `taskCenterEnabled=false` 时查询被禁用，但卡仍显示 0，点击后进入受 `RequireFeature` 拦截的路由。
- 查询使用 `listGenerationTasks(300).catch(() => [])`，权限、网络或 Backend 错误都会显示为“本周还没有生成任务”，把“不可读取”误报为“确实为零”。
- 影响：功能开放边界和运营统计不可信；用户无法区分没有任务与读取失败。
- 解决：把 `taskCenterEnabled` 传入仪表盘并隐藏/禁用任务卡；保留 query error，显示可重试的“统计暂不可用”，不得吞错为 0。

### 重要项 3：首页“本周任务”不是完整周统计

- 仪表盘最多读取最近 300 个任务，再从中筛选 7 天；高频账号一周超过 300 个任务时会无提示少计。
- 解决：使用 Backend 聚合接口或按时间游标分页到 7 天边界；若暂不实现，文案明确“最近 300 条中的本周任务”，并显示截断提示。

### 重要项 4：浅色画布修复绕过设计 token

- `b0fa90b` 在组件内硬编码 `#e6e6e6`，与 Primitive → Semantic → Component token 合同不一致。
- 解决：新增/复用浅色 canvas background token，由主题表提供；组件继续只读取 `theme.canvas.background`，并补明暗主题静态/Edge 视觉验证。

## 正面证据与验证缺口

- 固定 4ba 快照前端 `tsc --noEmit` 退出码 0。
- Backend `internal/protocol` 通过；新增 4 个不依赖 SQLite 的 Provider 角色测试通过。
- Windows `internal/service` 全包只因 `CGO_ENABLED=0` 的 go-sqlite3 stub 失败，不能视为业务失败，也不能视为通过。默认模型持久化、历史视觉/当前声音和修订产物登记测试仍必须在隔离 Linux CGO 中运行。
- 本轮没有 Bun 可执行文件，因此新增 Web Bun 测试尚未独立运行；源码测试存在不等于已执行。
- 没有 Microsoft Edge 视觉验收；首页、浅色画布、工作台并行状态和横向滚动仍是待验证项。

## 推荐决策与后续阶段

当前结论为 **NO-GO：不要直接创建当前 `c8b60ce` merge commit，也不要原样 merge/cherry-pick 全部 8 提交**。

建议在阶段 C 前新增“修复尾差阶段 B2”，经用户批准后：

1. 保持当前未提交 `c8b60ce` 合并，固定 `4ba9694`，以 `--no-commit` 方式引入修复尾差。
2. 逐区块解决 3 文件/15 区块，语义复核 8 个自动合并重叠文件。
3. 前 5 个提交的功能取并集，但先修复 xAI 尾帧三路径不一致；保留本地精确计价、百炼、30MB Runtime、panic guard、镜头删除/解绑和任务恢复合同。
4. `b0fa90b` 改为 token 化后纳入。
5. 首页仪表盘建议暂缓；先纳入 `4ba9694` 的工作台 CSS 横向滚动修复。若要纳入首页，必须先修功能开关、错误态和 300 条截断统计。
6. 完成 Linux CGO Backend、新增 Web/Bun、全量 Web/Agent/构建和最终 staged-bytes 门禁后，再单独决定 merge commit；push 仍是后续独立授权。

## 修复尾差阶段 B2 实际结果

- 2026-08-31，用户批准阶段 B2。开始前线上 `official/main` 已前进到 `fb089b2d2125e9ef1fcd2440f7153b8cce5f8157`；本阶段继续只处理已审计固定 `4ba9694`，新尾差未 fetch/merge。
- 复核私有恢复点 14 个文件：登记数据库、WAL、`.settings-key`、迁移标记、独立恢复库和 Git bundle 哈希一致；逐文件 ACL 没有当前用户、SYSTEM、Administrators 之外的条目。
- 当前已解决 `c8b60ce` 索引先固化为安全树 `1972ff0dcaec2fc23ca31028d80ae6e78519d041`、双父安全提交 `82758b208856c3d1a9697c5e35e49ae68453c25f`，保存在 `refs/codex/safety/b2-pre-4ba9694`。随后结束旧未提交合并，从干净 `HEAD` 标准打开 `git merge --no-commit --no-ff 4ba9694`。
- 直接重开会重新出现 18 个文件/52 区块；使用安全提交与 `4ba9694` 的三方自动合并树恢复既有人工解决后，精确收敛到审计预测的 3 文件/15 区块，并全部逐区块解决。当前 `MERGE_HEAD=4ba9694`、未解决索引项 0、非暂存改动 0、`VERSION=v1.1.4`。
- 前五个短剧/协议修复组已纳入；xAI 在协议 Adapter 与 Backend 类型化 Provider 中对尾帧 fail-closed，R2V 忽略陈旧首尾帧元数据，新增协议/Service/Web 回归断言。工作台保留本地删除、解绑、任务恢复、预览与 fail-closed 合同，同时接入默认模型、统一报价、并行提交、历史视觉版本、当前声音、台词和 mention。
- 浅色画布背景改为 `canvas.workspaceBackground` 组件 token，`InfiniteCanvas` 不再硬编码颜色。首页仪表盘按方案暂缓，稳定首页与全局样式只做 Prettier 机械格式化；只纳入工作台资产名横向滚动 CSS。
- 验证：Web TypeScript 退出码 0；尾差 28 个受支持文本文件 Prettier 通过；xAI 协议测试和 xAI Service 专项通过；JSON、gofmt、冲突标记、cached diff 和结构合同通过。Windows Go 系统缓存 ACL 失败后改用项目隔离 GOCACHE，首次编译回传丢失但缓存完成，随后精确复跑通过。
- Web Bun 专项未取得结果：主机无 Bun；已有 Bun Docker 镜像可启动，但 Docker Desktop 在当前执行隔离中无法看到 D 盘 bind source，两次运行均为 0 个测试。没有把 0 测试记为通过。Linux CGO Backend、Web Bun 全量、Canvas Agent、生产构建和 Edge 视觉留阶段 C。
- 没有 merge commit、push、PR、远程 CI、部署、运行数据库/卷、技能同步或 Provider 调用。下一步等待用户批准阶段 C。

## 阶段 C 完整验证与发布可行性评估

- 2026-08-31，用户批准阶段 C，并要求完成后只评估 merge commit 和 push 可行性，不执行二者。
- Backend：使用正式 `backend/Dockerfile` 的 `backend-test` target 构建隔离 Linux CGO 镜像，`go test -count=1 ./...` 全部包通过；测试镜像 manifest list 为 `sha256:4b57312dbd2f1be77bc11c51eaad9573b9fada7727014b39ac62ea979ebb6993`，未挂载运行数据、备份、密钥或数据卷。
- Web：首次完整主套件 1119/1122，3 个失败均为测试合同漂移：xAI 测试按错误签名未传到 options，两个源码断言依赖单行空白/旧两参数调用。只修正测试夹具和空白鲁棒性后，专项 23/23、主套件 1122/1122、admin regression 10/10、跨 Runtime 1/1、TypeScript 全部通过。
- Web 生产构建：使用与正式 `web-build` 等价、但排除本机 `node_modules` 的临时 Dockerfile；三平台 Comfy Bridge 通过，Vite 转换 13,469 个模块并完成生产构建。仅有既有大 chunk 和插件耗时警告；未导出、部署或替换运行镜像。
- Canvas Agent：最终格式化字节在主机权限下全量 328/328、0 失败，`npm run build` 通过；Windows 无 symlink 权限的拒绝子项继续由 Linux CI 保留。
- 最终格式/配置：145 个暂存受支持文本全部 Prettier；暂存 Go 无 gofmt 漂移；JSON 可解析；`pending-test.mdx` 无重复二级标题；default/dev/local/server/build/deploy/deploy+build 七组 Compose 均通过 `config --quiet`。deploy 首轮因必填 `DATABASE_URL` 缺失按设计拒绝，使用仅供解析的虚拟 URL 后通过，没有启动服务。
- 运行边界：现有 Backend/Web 仍为 `:local` 镜像，running/healthy、RestartCount=0、OOM=false；没有候选部署或 Edge 候选 UI 验收。Edge 运行验收属于后续候选部署门禁，不把源码/构建通过写成浏览器通过。
- 最终远程复核：`origin/codex/upstream-20260828-ab89c05=9815464`，`origin/main=11931d0`，远程检查点是当前 `HEAD=3510997` 的祖先。官方实时 `main=e124a9ea0459287957c309f8a859d415a3084b03`，所以 `4ba9694..e124a9e` 是新的未审计尾差，继续保持在外。
- merge commit 可行性：**可行，但需单独授权**。当前固定 `MERGE_HEAD=4ba9694`、未解决/非暂存为 0，完整源码门禁通过；提交后必须核对双父为 `3510997` 与 `4ba9694`、最终树与阶段 C 暂存树一致、`VERSION=v1.1.4`，且不得声称跟上实时官方 `e124a9e`。
- push 可行性：**当前未提交状态不可推送为完整合并；在 merge commit 单独批准并成功后，推送到用户 fork 的同名集成分支技术上可 fast-forward，无需 force，但仍需单独授权和推送前 SHA 复核**。不更新 `origin/main`，不自动创建 PR/CI/部署。

## 本地 merge commit 与后续三方评审

- 用户单独授权后，已创建 merge commit `d0def5f807e7d8c410103c85bd5ded43506b39b3`；双父 `3510997` / `4ba9694`，tree `49580c569687e0ef30b47cce0f5d1aa8e9aa9ff7`。没有 push。
- fork推送评审：远程集成分支 `9815464` 是新提交祖先，ahead/behind为14/0，dry-run成功；分支公开且未保护。push到该分支可fast-forward，但不会触发现有仅监听main push的quality/publish工作流，因此远程检查点不等于远程CI通过。
- 公开内容预检：待推送194个路径中无 `.env`、`.settings-key`、数据库、私钥/证书扩展名；新增行高置信 AWS/GitHub/OpenAI风格Key、Bearer和私钥头命中均为0。该启发式不能替代GitHub Secret Scanning，但未发现明确阻断项。
- 当前部署评审：运行镜像仍对应阶段16代码 `f9317fa`，不是 `d0def5f`。两者相差90提交、659文件、+128,714/-21,305；同为 `v1.1.4`，版本字符串无法区分。候选新增数据库/启动迁移、skills包、项目分页/分镜/资产/后台和协议功能，未来部署必须另做备份、迁移演练和Edge验收。
- 官方实时评审：`official/audit-current=b7348ab`、版本 `v1.2.3.1`；相对固定 `4ba9694` 有22提交/216文件。新本地与官方实时的共同基点仍为 `4ba9694`，双方独有48/22提交；未来模拟冲突29文件/52区块。不得把当前本地提交或部署描述为跟上官方实时主干。
