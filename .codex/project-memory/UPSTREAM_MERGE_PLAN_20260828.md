# 2026-08-28 官方上游增量合并执行计划

## 状态

- 当前阶段：阶段 0-2 已完成，等待用户批准阶段 3。
- 阶段 3 未获批准前，不创建集成分支、不执行 merge。
- 本计划不授权 push、PR、CI、tag、Release、付费调用、数据库恢复或卷删除。

## 固定起点

- 本地起始分支：`main`。
- 本地起始提交：`405d25aec95b307144e5a3e84170bebb948ebe31`。
- 阶段 0-2 完成后，项目记忆与备份登记以一个纯文档提交固化；阶段 3 应从届时干净的当前 `main` 创建集成分支。功能差异和 29 个冲突仍以 `405d25a` 对 `ab89c05` 的代码关系为分析基线。
- 当前应用版本：`v1.1.4`，由根 `VERSION` 在 Vite 构建时注入。
- 阶段 2 实时固定官方快照：`main@ab89c05362394623d00c95a4899be03e5243e75a`。它比旧分析的 `913cf4b` 多 1 个提交；后续阶段开始前仍需确认该 SHA 未被用户要求刷新。
- 用户 fork：`https://github.com/dingxin150015-commits/open-ai-canvas`。计划作为后续统一维护入口；当前不 push、不创建 PR。

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

- 只有用户在阶段 0-2 汇报后明确批准，才创建集成分支并引入官方 SHA。
- 后续每个阶段独立汇报、独立批准。

## 永久 NO-GO

- `git reset --hard`、宽范围 checkout、批量 ours/theirs。
- 在当前 `main` 直接 pull/merge。
- `docker compose down -v`、删除或替换数据卷。
- 未验证恢复点前迁移、重建或部署。
- 将大型 fork 整体合入。
- 未重新授权的真实模型、OSS 删除、GitHub 写入、代理或账户变更。
