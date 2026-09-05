# 官方 `v1.2.5` 固定增量审计

## 固定对象与分支

- 稳定基线：`codex/dingxin-stable@44e8c10205d0e306031ccd12014cdd3f44f67651`，已推送并独立复核为 `origin/codex/dingxin-stable` 同一 SHA。
- 历史检查点：`origin/codex/upstream-20260828-ab89c05` 已从 `e02cee1` fast-forward 到同一 `44e8c10`，未 force、未改 `origin/main@11931d0`。
- 集成分支：`codex/upstream-v1.2.5-20260904`，从稳定基线创建，尚未 push。
- 官方标签对象：`2a45d0cc3eb25424907e1f29a2aae9a657db7ece`；peeled commit：`f8e87bcc4ce3e6f7eae7a89dc8b9231801116072`。
- 标签没有可验证 Git signature；本轮以 HTTPS 官方 remote、固定标签对象和 peeled commit 三者一致作为来源证据，不把它描述为签名发布。
- 实时 `official/main@7ec51732f7b736f9250fd93b558a1a98f6652a71` 的 6 个后续提交不纳入本轮。

## 真实增量规模

- 共同基点：`4ba9694f459808cec8fa57b55f2db767d9f412b8`。
- 当前本地独有 54 个提交，官方目标独有 66 个提交。
- 官方完整 Git 对象差异：506 个文件、`+38,322/-3,646`；本地差异：303 个文件、`+83,085/-18,138`。
- 双方重叠路径 89 个；merge-tree 模拟得到 47 个显式冲突文件、149 个冲突区块，另有 40 个 changed-both 自动合并路径需要语义复核。
- 早先 GitHub Compare API 返回 300 个文件、`+16,725/-657`，是 API 文件列表上限造成的截断估计；本记录以固定标签完整 fetch 后的本地 Git 对象统计替代旧值。

## 官方功能族

- 在线更新：后台检查、Host Updater、强制备份、独立迁移、健康稳定窗口和失败回退；Windows 仅包含编译兼容，本地 fork 更新源尚未适配。
- 数据/存储：schema v2、`resources.upload_key` 独立迁移和唯一索引、COS/S3/CDN、资源上传幂等、回收站、历史位置和删除一致性。
- 模型/任务：声明式协议、官方插件包、能力保护、任务恢复、真实用量结算、视频参数/提示词约束和供应商错误识别。
- 项目/画布：五类资产、历史和生成副本、大画布性能、媒体预览/播放器、资产候选/引用/同步、空白起点和外观设置。
- 创作/素材：参考内容重构、模型价格展示、批量上传、自定义分类、分页、网格密度和历史缺失字段兼容。
- 平台：密码找回、平台外观、微信/支付宝支付插件、AI 审美批改插件、文档和发布工作流。

## 47 个显式冲突文件

### 根目录、部署和文档

- `.gitignore`
- `CHANGELOG.md`
- `README.md`
- `VERSION`
- `backend/Dockerfile`
- `backend/cmd/server/main.go`
- `docker-compose.deploy.yml`
- `docs/content/docs/backend/backend-database.mdx`
- `docs/content/docs/progress/pending-test.mdx`
- `web/package.json`

### Backend 合同

- `backend/internal/service/model_capability.go`
- `backend/internal/service/resource.go`

### Web 资产、模型和画布

- `web/src/components/assets/asset-library-picker-modal.tsx`
- `web/src/components/canvas/asset-picker-modal.tsx`
- `web/src/components/canvas/canvas-assistant-panel.tsx`
- `web/src/components/canvas/canvas-node-generation.ts`
- `web/src/components/canvas/canvas-project-asset-modal.tsx`
- `web/src/components/canvas/canvas-resource-mention-textarea.tsx`
- `web/src/components/canvas/canvas-toolbar.tsx`
- `web/src/components/canvas/infinite-canvas.tsx`
- `web/src/components/layout/workspace-top-bar.tsx`
- `web/src/components/model-picker.tsx`
- `web/src/components/video-settings-panel.tsx`
- `web/src/lib/canvas-theme.ts`
- `web/src/lib/canvas/canvas-resource-references.ts`
- `web/src/lib/model-pricing.ts`
- `web/src/pages/admin/plugins/plugins-page.tsx`
- `web/src/pages/assets/index.tsx`
- `web/src/pages/canvas/project.tsx`
- `web/src/pages/canvas/shared.tsx`
- `web/src/pages/canvas/use-canvas-media-tools.ts`
- `web/src/pages/create/index.tsx`
- `web/src/pages/projects/detail/assets.tsx`
- `web/src/pages/projects/detail/chapters.tsx`
- `web/src/pages/projects/detail/workflow-production-workbench.tsx`
- `web/src/pages/projects/detail/workflow-shot-references.ts`
- `web/src/pages/projects/detail/workflow-stage-views.tsx`
- `web/src/pages/projects/detail/workflow.css`
- `web/src/pages/projects/index.tsx`
- `web/src/services/api/projects.ts`
- `web/src/services/user-data-sync.ts`
- `web/src/styles/globals.css`

### Web 测试

- `web/test/canvas-model-policy.test.ts`
- `web/test/canvas-node-registry.test.ts`
- `web/test/create-library-button.test.ts`
- `web/test/generation-task-local.test.ts`
- `web/test/project-chapter-skill-runtime.test.ts`

## 阶段 4 解决合同

1. 文档/版本/测试：保留本地 `Unreleased` 和完整待测状态，加入官方 v1.2.3–v1.2.5 历史；最终 fork 版本先保持待阶段 6 决策，不在未验证索引中冒充官方发布；测试按功能并集修复夹具。
2. 部署：保留 `backend-test` Linux CGO target、生产 CORS fail-closed、Provider 插件开关和本地 Compose；接入 buildinfo/migrate-schema/Host Updater，但 fork 更新源适配前保持更新写路径不可用。
3. 数据/资源：迁移必须把本地 Catalog/价格/Skill/项目字段与官方 schema v2、upload_key、支付/回收站表取并集；资源读取/删除/上传继续校验用户归属、引用和物理失败。
4. Provider/计费/任务：保留精确 Bailian Adapter、Ready/Planned、完整 intent/price tier、xAI fail-closed、脱敏日志、SSRF 和 OSS 物化；吸收声明式协议、任务恢复、真实 Token 补扣和新插件骨架。
5. Web：把官方素材分类/分页/批量上传、画布历史/副本/性能、媒体预览、外观、工作流修复移植到当前本地架构；保留本地 Create 设置、资产 mention、任务恢复、模型能力、安全错误和稳定首页；不整文件选边。

## 当前停止条件

- 阶段 3 已从本地审计提交 `39c7e62` 打开固定 no-commit/no-ff 合并；`MERGE_HEAD=f8e87bcc4ce3e6f7eae7a89dc8b9231801116072`。
- 真实 merge 得到 49 个未合并路径：47 个文本冲突和 `plugin-packages/autodl-comfyui/manifest.json`、`web/src/pages/home/index.tsx` 两个 modify/delete 决策；工作树冲突标记为 139 个。差异来自 modify/delete 不生成文本标记以及 merge-tree 对相邻区块的计数方式。
- 阶段 4 已按五组逐区块解决并加入 index，未解决路径为 0。AutoDL 源 Manifest 作为本地扩展源保留；稳定首页按既有产品决定保留；官方打包插件、素材/回收站/搜索、平台外观、画布媒体/历史/外观、支付和更新基础同时纳入。
- 当前根版本为 `v1.2.5+dingxin.1`，明确区别于官方原版；是否作为最终构建标识仍由阶段 6 前复核。
- 阶段 4 结构检查：冲突标记 0、未解决索引 0、非暂存文件 0；冲突 Go 文件已执行 gofmt，`web/package.json` 与保留的 AutoDL Manifest 可解析，cached diff check 通过。稳定首页组件保留并恢复 `/home` 路由；官方 `/` 创作入口保持不变。以上不是阶段 5 类型、测试、构建或运行验证。
- 2026-09-05 只读吸收审计确认：417 个本地未触碰的官方变更路径与 `f8e87bc` blob 逐字节一致；229 个官方新增路径全部存在。在线更新、upload_key迁移、素材分类/回收站/批量上传、密码找回、平台外观、支付、AI审美、画布历史/外观/空间索引、视频布尔能力和80个官方插件包的结构入口均存在。
- 已确认一个阶段 5 阻断项：官方已把 Create 文本生成统一为 `runBackendGenerationTask`，并删除 `backendModelRuntimeRequired`；当前 `web/src/pages/create/index.tsx` 仍导入该已不存在的导出，并保留 `requestImageQuestion` 旧分流。必须在阶段 5 采用官方 Backend-only 生命周期，同时保留本地任务元数据、reasoning、Prompt Cache所需的其他路径，然后运行类型和任务恢复测试。
- 有意不原样吸收的官方差异：保留 `/home` 稳定首页和AutoDL源 Manifest；保留本地画布底色/网格 token；章节角色再提取只对待确认候选去重，而不以已确认角色永久阻止再提取。这些是本地产品策略，不是遗漏，但需阶段 5 专项锁定。
- 官方 Host Updater代码已吸收，但默认仓库仍为 `ddcat-ai/open-ai-canvas`；在更新源、镜像仓库和版本比较适配用户 fork 前只能视为结构存在、不可用于本地自维护版本更新。
- 当前状态是 `merge_open_resolved_unvalidated`。阶段 4 完成后必须停止；阶段 5 尚未运行，禁止 merge commit、候选 push、stable 晋升、部署、真实数据迁移或 Provider 调用。

## 阶段 5：源码、构建与克隆迁移门禁（2026-09-05）

- Create 已采用官方 Backend-only 任务生命周期，同时保留本地 `metadata`、`reasoning`、`clientOperationId`、任务恢复/取消和 Dreamina 专用路径。
- 修复真实合并回归：画布节点动作 Context 恢复稳定 memo，节点拖动首帧 `dragPreview` 同步关闭浮层；保留 `/home`、AutoDL 源 Manifest、本地主题 token、章节角色再提取和 Provider 脱敏合同。
- 修复官方 v1.2.5 公告图片自引用缺陷：仅豁免正在原子替换的当前公告和正在丢弃的自身草稿，其他业务引用继续 fail closed。
- Backend 测试镜像复制全部插件包；项目封面引用测试补齐快照表夹具，字符串错误码测试同步本地脱敏合同。
- Web TypeScript、默认全套测试（主套件 1158/1158、画布 33/33、跨 Runtime 1/1）、补充专项和生产构建通过；构建转换 13,543 模块，无 Rolldown panic。直接并发 `bun test test` 会因共享 `window/localforage` 竞态产生假失败，不能替代项目串行分组。
- Backend 隔离 Linux CGO `go test -count=1 ./...` 全部通过；Canvas Agent 在 Node 22 + Bun 1.4 Linux 隔离环境完成 328 项（322 pass、6 个 Windows 专项 skip、0 fail、0 cancelled）及 TypeScript 构建。为避免 unref heartbeat 令测试宿主提前退出，测试增加了有界 renewal-start watchdog；Windows 本机长驻进程树用例受 `taskkill` 权限限制，不作为最终宿主。
- Compose 六种组合仅解析通过；2 个 PowerShell 文件 AST 解析通过。Host Updater 在 fork 源适配前仍不可执行。
- 候选镜像：Backend `sha256:6b04464b6d03246db988cbf5b79906fcca20bddb5a26c7033eea63e269a51ca4`，Web `sha256:36dc2302122eb17be3e4c14de059978d872a2c14ad52f2280770333129d50118`。
- `BKP-20260905-143709-STAGE5-V125` 已 `verified/protected`：预迁移 `ok`/外键0/72表；克隆卷连续两次无网络启动 healthy，迁移后 `ok`/外键0/80表，业务计数和 OSS 密钥合同不变。
- 真实 Backend 保持旧镜像 `sha256:8182cc0a...` 且 `running/healthy`、Paused=false、RestartCount=0；真实卷和 Web 均未替换。没有 merge commit、push、PR、Edge 控制或 Provider 调用。
- 最终 staged-bytes 复核通过：517 个 staged 文件、unmerged 0、unstaged 0、新增冲突标记 0、敏感文件路径 0、高置信 Secret 新增命中 0、`git diff --cached --check` 通过；仓库既有 `.claude` 文档中的 3 行冲突示例不属于本次新增。
- 当前状态：`stage5_complete_ready_for_stage6_approval`。阶段 6 merge commit 仍需用户明确批准；不得自动 push、部署、控制 Edge 或调用 Provider。

## 阶段 6–7 授权与画布颜色门禁

- 2026-09-05 用户批准阶段 6 本地 merge commit 与阶段 7 fork 候选分支推送连续执行；授权不包含 `origin/main`、稳定分支晋升、部署、Edge 控制或 Provider 调用。
- 当前候选继续保留本地 `workspaceBackground` 和柔和网格 token，但该选择不是永久封死。最终 Microsoft Edge 视觉验收必须将画布颜色列为独立验收项，比较本地方案与官方 v1.2.5 的 `#f0f0f0/#000000 + 80% 网格` 方向。
- 若最终选择官方方案，应仅调整画布语义 token/默认外观来源与对应测试，不回退官方已吸收的自定义外观、项目持久化和账号级默认功能；不得在部署前凭静态测试代替视觉决定。

### 阶段 6 结果

- 已创建本地 merge commit `b6cff796d81361b1e98fd8a83f3f340d340a63e1`，提交说明为 `chore(upstream): 官方版本 - 集成 v1.2.5`。
- 两个父提交精确为 `39c7e626cdf58c7f901cdfaa7f138041e44d1459` 和 `f8e87bcc4ce3e6f7eae7a89dc8b9231801116072`；官方固定提交是新 HEAD 的祖先。
- merge commit 后工作树干净，阶段 7 只允许推送 `origin/codex/upstream-v1.2.5-20260904`，不允许 force、更新 `origin/main`、稳定分支晋升或部署。

### 阶段 7 结果

- 在项目记忆提交 `dced038188ca2e400fe6e7a5704cd0be9a3b7eef` 上执行 dry-run，结果仅为新建 `origin/codex/upstream-v1.2.5-20260904`；随后非 force 推送成功并建立同名 upstream tracking。
- 独立 `git ls-remote` 确认远程候选 SHA 与本地 `dced038...` 一致；`origin/main` 仍为 `11931d0085ccfec3952a4a7e30bd482173c51492`，未变化。
- `gh run list` 对候选分支返回空数组：此次 branch push 没有产生可见 GitHub Actions 运行，不能宣称远程 CI 通过；本地 Stage 5 完整门禁仍是当前验证依据。
- 本节作为阶段 7 最终记忆检查点再 fast-forward 推送一次；仍不创建 PR、不晋升 `codex/dingxin-stable`、不更新 `origin/main`、不部署。
