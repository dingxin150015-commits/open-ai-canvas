# 固定 v1.6.6 增量合并、验证与部署执行记录

更新时间：2026-10-12 23:14 CST。当前是启动事实，不是完成报告；每个关键动作后追加证据和时间。此任务获用户明确要求：合入最新上游、解决冲突、全量验证后更新本地部署，并非强制推送到自己的 GitHub fork。用户此前的 Edge 手工验收、恢复点和旧 A 保留规则仍须按当前风险执行；真实付费模型调用没有获准。

## 已确认基线

- 官方 GitHub `releases/latest` 于 2026-10-12 返回 `v1.6.6`、`draft=false`、`prerelease=false`、发布时间 2026-10-11 16:50:56（API 显示格式），`git ls-remote official` 的 `refs/tags/v1.6.6^{}` 为 `27c8066e6da9f0d4754dc0e5c40097049fe730cf`，标签对象为 `8cd7da6dd0d89179f11390112bc2893c43c8ee5b`。移动 `official/main=e70fe1b9` 不作为目标。网页搜索索引只显示 v1.6.4，已被实时 API+Git 标签取代。
- v1.6.0 A″ 位于 `.local/worktrees/upstream-v160`，`HEAD=7419e2b3db6aa08e47a208529c0f27173ccf39c3`、`MERGE_HEAD=42f14cd73eec7c2ded98eb15b6a4f5aa2f3f3b4a`，既有暂存树 `8b6d037a84b4eb6f62e7728e94dd1cfbff3f3c71`；未合并、未暂存、未跟踪均 0。它尚未获 Edge 手验，不能作为稳定已部署 A。原 A、A′、A″ 不覆盖/删除。
- 生产容器 `open-ai-canvas-backend-1`/`open-ai-canvas-web-1` 均 healthy、restart 0，当前分别为镜像 `sha256:d67ee794...`、`sha256:9bfc9c67...`，Compose 来源 `.local/worktrees/dashscope-plugins/docker-compose.local.yml`。后端 ready 报 `v1.2.7+dingxin.1.catalog1 / ad3efa10`、Schema 9/9、活动 Worker 0；生产卷 `open-ai-canvas_backend-data` 挂载 `/data`。这不是 v1.6.0 候选，也不等于 v1.6.6 部署条件满足。
- 用户 fork 的 `origin` 为 `https://github.com/dingxin150015-commits/open-ai-canvas.git`，官方 remote 是只 fetch 的 `official`。第一次 `origin ls-remote` 遇到 Windows Schannel TLS 握手失败；尚未核实目标远端分支状态，不更改全局代理或 Git TLS 设置来绕过。

## 本轮待执行门禁

1. 在 A″ 只暂存新记忆后记录最终树，使用不移动旧分支的本地 WIP 双父检查点固定未验收来源；新建隔离 A‴，从固定 v1.6.6 标签增量合并。先记录 merge-tree 预判与真实冲突，再按模块合成和审计自动重叠；不在 A″ 原工作树叠加 merge。
2. 验证 Web/Backend Linux CGO/Canvas Agent、锁与生成插件包、Compose、镜像、无端口 HTTP；对生产 Schema9 和 A″ Schema45 分别用隔离副本演练到候选目标，并校验重复迁移、完整性、外键、脱敏计数与匹配密钥。不得在此步骤写生产卷。
3. 生产发布前重新确认当前服务、运行任务、卷和 GitHub 目标；制作新的私有 `verified/protected` 完整恢复点，保留旧镜像标签和旧 A。先用户 Microsoft Edge 隔离预览验收；若测试失败或不满足安全门禁，停在候选，不宣称“测试无误”。
4. 只有前述门禁通过，才以保留卷、明确 Compose 工作目录升级本地部署，并验证 build-info、ready、Schema、业务计数/密钥合同和 Edge；然后正式提交、非强制推送到核实的用户 fork 分支、核对远端 SHA/CI。绝不 force，绝不删除旧工作树/数据/卷。若 Edge 手验缺失，报告并等待用户，不替其声称通过。
