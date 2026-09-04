# 官方 `v1.2.5` 增量集成与长期双轨维护计划

## 用户批准与当前边界

- 2026-09-04 用户批准以官方稳定 Release `v1.2.5` 为本轮目标，不跟随移动的 `official/main`。
- 固定标签对象：`2a45d0cc3eb25424907e1f29a2aae9a657db7ece`；peeled commit：`f8e87bcc4ce3e6f7eae7a89dc8b9231801116072`。
- 实时官方 `main@7ec51732f7b736f9250fd93b558a1a98f6652a71` 比目标多出的 6 个提交保持在本轮之外，后续单独审计。
- 用户批准阶段 0–4 连续执行并在阶段 4 后停止：阶段 0 固定目标和维护模型；阶段 1 保存当前 fork 检查点并建立长期稳定分支；阶段 2 固定目标、完整审计；阶段 3 在独立分支打开 no-commit/no-ff 合并；阶段 4 手工语义解决冲突。
- 阶段 5 全量验证、阶段 6 merge commit、阶段 7 fork 候选推送、阶段 8 稳定分支晋升、阶段 9 本地部署仍需阶段 4 汇报后的新批准。
- Astravia Windows 仅在阶段 0–9 全部完成后，按隔离 worktree、低权限、无生产密钥/数据库的边界作为辅助工作台和第二审阅面；当前不安装、不启动。

## 长期维护模型：双轨来源、单树交付

- `official/main` 和官方 Release refs：官方只读来源；不允许 push。
- `origin/codex/dingxin-stable`：用户 fork 的长期稳定交付分支，只接受完成源码、迁移、部署和 Edge 门禁的候选。
- `codex/upstream-v<version>-<date>`：每次固定官方 Release 的临时集成分支。
- `codex/local/<feature>`：本地自研功能分支，优先放在独立 Provider/Adapter/Manifest/component/hook/plugin/skill/Compose overlay 中。
- `codex/hotfix/upstream-carry-<issue>`：上游未发布但生产必须修复的临时补丁；记录上游来源和删除条件，上游正式修复到达后移除。
- 最终只维护一棵可编译、可部署代码树；不复制 `official-copy`/`local-copy` 两份完整项目。

## 已知规模和审计重点

- 相对已集成共同基点 `4ba9694`，固定标签完整 fetch 后确认官方 `v1.2.5` 有 66 个提交、506 个文件、`+38,322/-3,646`。早先 300 文件数字来自 GitHub Compare API 文件列表上限，已被本地 Git 对象统计替代。
- 当前本地与官方目标有 89 个重叠路径；merge-tree 模拟为 47 个显式冲突文件、149 个区块，另有 40 个 changed-both 自动合并路径需要语义复核。
- 高风险区域：Backend schema/model/provider/task/resource/finance/channel model；协议 Registry 和能力；Web 资产/mention/节点生成/工具栏；Docker/Compose/VERSION/CHANGELOG；插件 Manifest、文档和测试入口。
- 先完成 v1.2.5 集成，再单独做“本地扩展边界收口”，避免把上游大增量和本地重构混为一个变量。

## 必须保留的本地合同

- Qwen/Wan/HappyHorse 专属能力、Ready/Planned 状态和统一有效能力边界。
- 完整 `ModelRequestIntent`、服务端精确价格档和 `PriceConfigured` fail-closed。
- xAI 尾帧及首帧/角色引用混合输入 fail-closed。
- Provider/资源/任务日志脱敏、SSRF/签名 URL 和资源归属安全。
- OSS 结果物化、素材引用检查、30MB Local Runtime 边界。
- Create 参数持久化、声音/水印明确开关、用户作用域缓存和稳定首页决定。
- Microsoft Edge 是浏览器验收唯一来源。

## 应吸收的官方修复族

- `resources.upload_key` 独立迁移和用户维度唯一索引。
- 声音、水印、视频模型能力默认值、完整提示词长度拒绝和工作流产物修复。
- 声明式模型协议、任务恢复、运行时协议校验、模型价格管理。
- 画布历史/生成副本、素材分类/分页/批量上传、远程媒体兼容和删除引用检查。
- Host Updater Windows 编译兼容；实际在线更新仍不得在本地 fork 上启用，直到更新源、版本比较和回滚合同改为用户 fork。

## 阶段 0–9

1. 阶段 0：固定目标、记录双轨所有权和授权边界。
2. 阶段 1：审计并 fast-forward 推送当前稳定检查点；独立复核远程 SHA；建立 `codex/dingxin-stable`。
3. 阶段 2：fetch 固定 `v1.2.5`，计算 merge-base、独有提交、显式/语义冲突和迁移顺序，形成审计记录；不 merge。
4. 阶段 3：从稳定分支创建 `codex/upstream-v1.2.5-20260904`，执行 `git merge --no-commit --no-ff <fixed-sha>`，只核对真实冲突。
5. 阶段 4：按文档/版本、部署、数据库/资源、Provider/计费/任务、Web 画布五组逐区块语义取并集；不批量 ours/theirs；结束时允许未验证的 resolved index，但不创建 merge commit。
6. 阶段 5：格式、Backend Linux CGO、迁移双启动、Web/Agent/Docs/Compose/Windows 脚本及最终 staged-bytes 全量门禁。
7. 阶段 6：采用明确的 fork 构建标识，复核 `MERGE_HEAD` 后创建本地 merge commit。
8. 阶段 7：dry-run 后只推送 fork 候选分支，不 force、不更新 `origin/main`，独立核对 SHA 和真实 CI 状态。
9. 阶段 8：候选通过后经单独批准晋升 `origin/codex/dingxin-stable`。
10. 阶段 9：新恢复点、无网络克隆卷双启动、候选镜像、保留卷 recreate 和用户 Microsoft Edge 验收；不得执行 `down -v`。

## NO-GO

- 不把移动的 `official/main` 当作本轮目标，不在阶段中途追入新尾差。
- 不在 `main` 或现有历史检查点直接打开合并，不 force push，不改写历史。
- 不用整文件 ours/theirs 代替语义合并。
- 阶段 4 完成不代表代码已验证，不允许创建 merge commit、推送候选、部署、迁移真实卷或调用 Provider。
- 官方在线更新功能在适配用户 fork Release/channel 前保持禁用或只读，不能用官方镜像覆盖本地扩展。
