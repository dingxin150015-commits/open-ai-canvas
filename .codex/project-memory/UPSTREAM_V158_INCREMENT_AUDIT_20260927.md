# 官方 v1.5.8 增量冲突审计（只做模拟，未合并）

## 固定证据与边界

- 已确认：2026-09-27 `git ls-remote official refs/heads/main refs/tags/v1.5.*` 显示官方 `main` 与 `v1.5.8^{}` 同为 `258422757c98467f24381219e90f4809394ec226`。官方 `v1.5.6^{}` 为 `473d535ec161fd21001e1ad904d3b316e2d45c39`，且是前者祖先。不要把移动的 `official/main` 当实施目标。
- 已确认：有效候选工作树为 `.local/worktrees/upstream-v156`，分支 `codex/upstream-v1.5.6-20260921`，`HEAD=6b16d25e`、`MERGE_HEAD=473d535e`，当前 `git diff --name-only --diff-filter=U` 为空；仍未提交、未部署，Backend/Web 全量测试尚有失败。根工作树不是此候选。
- 审计把当前索引连同 9 个未暂存文件放入一次性临时索引，生成只用于模拟的完整候选树 `c2018dd72b10a6b2bc9f6e3fcf4f5492c3112c22`；临时索引已删除，实际索引与工作树未被暂存或改写。`git merge-tree --write-tree --messages --name-only --merge-base=473d535e c2018dd7 25842275` 生成预览树 `1fdc73cbdf10f7d672ab690aed71437d54513d85`，退出码 1。
- 已确认：官方 `v1.5.6..v1.5.8` 改动 830 个路径；本地候选相对 v1.5.6 改动 674 个路径，其中 157 个双方重叠。模拟得到 59 个显式冲突，均报 `CONFLICT (content)`；剩余 98 个重叠路径由 Git 自动合并，尚未完成语义审计。此计数不是“功能已兼容”或“可部署”。
- 已确认：本轮仅获取固定上游 Git 对象、运行只读差异/merge-tree 模拟并更新本记忆；未执行实际 merge、commit、push、构建、数据库操作或容器变更。

## 冲突解决合同

所有显式冲突默认按职责合成，不能对整文件盲用 ours/theirs。以下分组覆盖全部 59 路径；每批先读双方块、调用点和测试，再作最小合成。`needs-human` 指实现时发现两边产品目标真正互斥，不能靠编译通过代替决策。

| 批次 | 文件（显式冲突） | 策略和关键风险 |
| --- | --- | --- |
| 1 数据与协议（5） | `backend/internal/database/migrations.go`、`schema.go`；`backend/internal/model/models_channel.go`；`backend/internal/protocol/manifest.go`；`backend/internal/service/aliases_types.go` | synthesize，高风险。官方迁移已占 30–40，本地 30/31 分别为实际时长计费和 Comfy Bridge；不能覆盖官方编号或改已应用校验和。仅在确认本地候选迁移未部署后，把本地新增迁移追加到 40 之后；旧生产 Schema9 已记录时长迁移的谱系须把官方 9–40 顺移且跳过重复时长迁移，逐项验证无缺号/重放/校验和漂移。官方 SMS 插件贡献校验与本地 v3/v4 百炼 Manifest 均要保留；模型 tags 与本地 Eligibility 同时存在。 |
| 2 后端执行与错误（4） | `backend/internal/app/provider_http_client.go`、`admin_storage_delete_test.go`、`resource_delete_test.go`；`backend/internal/handler/response.go` | synthesize，中高风险。保留本地安全的 requestId 提取和上游 request-kind/错误 details 合同，统一响应语义；测试夹具可采用上游 WAL 库，但不得放宽 SSRF 或丢失本地数据保护断言。 |
| 3 前端协议与持久化（5） | `web/src/services/api/request.ts`、`web/src/services/canvas-generation-consumer.ts`、`web/src/lib/generation-error.ts`、`web/src/lib/canvas/canvas-generation-task-sync.ts`、`web/src/types/canvas.ts` | synthesize，高风险。错误 envelope 要同时保留上游 details/retry 和本地 requestId/错误分类；生成结果必须兼顾上游资源恢复/CAS 与本地提交后持久化、节点精确绑定，禁止把已删输出回退给同任务兄弟节点；Canvas 类型保留 Comfy provider 和上游产出模型快照。 |
| 4 画布节点与媒体（7） | `web/src/components/canvas/canvas-frame-node.tsx`、`canvas-node-content.tsx`、`canvas-node-prompt-panel.tsx`、`canvas-node-search-modal.tsx`、`canvas-node-toolbar.tsx`、`canvas-version-compare-modal.tsx`、`storyboard-assets-cell.tsx` | synthesize，中高风险。上游稳定媒体地址/缓存、字幕、产出模型标记、列表生成与节点卡片要进入；本地无提示词视频合同、Comfy/人像入口及资源权限不得丢。源地址直连与浏览器缓存的取舍需按签名和媒体安全合同验证。 |
| 5 画布外观与工作区（7） | `web/src/components/canvas/canvas-appearance-controls.tsx`、`canvas-short-drama-entry.tsx`、`canvas-toolbar.tsx`、`canvas-workspace-overlays.tsx`、`director/canvas-director-workbench.tsx`；`web/src/lib/canvas-theme.ts`；`web/src/styles/globals.css` | synthesize，中高风险。采用上游新入口/悬浮侧栏/3D 导演 UI，不覆盖本地窄屏底栏、Comfy 连接创建、导演台 F6 交互和深色网格可见性。上游黑底/黑网格默认与本地浅灰点网格互斥，需 Edge 视觉验收，不可仅选一侧。CSS 格式改动不能掩盖语义差异。 |
| 6 画布页面编排（4） | `web/src/pages/canvas/canvas-project-status-dialogs.tsx`、`project.tsx`、`use-canvas-director.ts`、`use-canvas-node-editor.ts` | synthesize，高风险。上游批量设置和 3D 导演输出接入；本地资源保存失败留页、画布工作模式、人像能力和视频预览保留。Director 应按新输出模型逐字段合成，不复制两个独立状态机。 |
| 7 创作、管理和设置（10） | `web/src/pages/create/index.tsx`；`web/src/pages/admin/channels/channels-page.tsx`、`components/channel-model-editor.tsx`、`components/channel-model-manager.tsx`、`settings/system-performance-page.tsx`；`web/src/pages/auth/auth-scene.tsx`；`web/src/pages/settings/channel-settings-pane.tsx`、`index.tsx`；`web/src/pages/tasks/index.tsx`；`web/src/components/layout/user-oss-settings-form.tsx` | synthesize，高风险。创作页 10 块涉及本地百炼参数/本机渠道和上游附件限制、媒体物化；管理员保留本地 planned/ready 只读及实际时长价格，同时接上游标签、重定价和列表状态。任务页保留本机即梦刷新/删除与上游媒体恢复，OSS 直连须维持安全提示。 |
| 8 测试与专题文档（7） | `web/test/canvas-appearance.test.ts`、`canvas-media-performance.test.ts`、`canvas-visual-contrast.test.ts`、`generation-storage-consistency.test.ts`；`docs/content/docs/backend/backend-database.mdx`、`overview/features.mdx`、`progress/pending-test.mdx` | synthesize，中风险。先确定生产合同再改断言；数据库文档覆盖官方 40 及本地后续迁移两谱系，测试不能为了消红而删除边界。待测试文件必须继续声明未通过项。 |
| 9 构建、部署与发布（10） | `.github/workflows/publish-images.yml`、`quality.yml`；`Dockerfile`、`backend/Dockerfile`、`VERSION`；`docker-compose.deploy.yml`、`dev.yml`、`server.yml`、`docker-compose.yml`；`web/vite.config.ts` | synthesize，中高风险。候选完成后再设 `v1.5.8+dingxin.1`，保留本地 CGO 测试 target、Bun 固定版本、插件开关与生产 CORS fail-closed；吸收官方可重用 commit 镜像、构建时间和 `CANVAS_PUBLIC_BASE_URL`。不能用上游 Server `CORS=*` 覆盖本地强校验。 |

非显式冲突：官方修改了多个 `.yingce-plugin` 归档及 `web/bun.lock`、`backend/go.mod/go.sum`。这些不手工改二进制或锁文件；最终按官方生成脚本与锁定依赖重新生成/校验。官方新 `backend/builtin/` Markdown 技能包、Agent 工具和资源访问策略应进行独立语义审计。尤其 98 个自动合并重叠路径中 `backend/internal/app/protocol_plugins.go`、`provider_protocol.go`、`task_worker.go`、`handler/routes.go`、`web/src/pages/canvas/use-canvas-generation.ts`、`web/src/services/api/generation-task.ts`、`resources.ts`、`web/src/pages/create/creation-workspace.tsx` 等不能因无标记而跳过。

## 后续执行顺序与门禁

1. 先冻结现有 v1.5.6 候选的完整快照与 Backend/Web 已知失败清单；只修阻断后续归因的基线错误，不为旧目标重复跑最终全量。当前 `MERGE_HEAD` 未关闭，禁止直接在原工作树继续叠加 v1.5.8 merge；真正实施前须明确如何结束现有 merge 或以经批准的隔离快照建立新候选。
2. 获准实施后，在新的隔离候选中只用固定 `v1.5.8@25842275`，先复核相同 merge-tree 计数，再依上表 1→9 处理。每批结束检查无冲突标记、调用方与最小专项测试；有产品目标互斥则停为 needs-human。
3. 数据库必须分别用生产旧 Schema9 谱系和官方 v1.5.6/v1.5.8 谱系的脱敏克隆验证升级、二次启动、历史数据/费用/加密字段不变，并保持后台使用中的真实卷不受影响。
4. 完成后再跑 Backend Linux CGO 全量、Web 全量/类型/生产构建、Canvas Agent 测试/构建、插件打包校验、Compose 配置检查和候选镜像；远程 CI、真实模型/OSS/Comfy 运行及用户 Edge 均为后续独立门禁。未通过前 NO-GO：不得发布/部署。
