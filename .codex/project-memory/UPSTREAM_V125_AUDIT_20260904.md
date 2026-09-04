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

- 审计完成后可按既有批准打开固定 no-commit/no-ff 合并并解决冲突。
- 阶段 4 完成后必须停止；未运行阶段 5 前，任何 resolved index 都是 `unvalidated`，禁止 merge commit、候选 push、stable 晋升、部署或 Provider 调用。
