# 影策项目级记忆入口

最新记录：`UPSTREAM_V127_STAGE6_20260908.md`。自动化源码门禁通过，但Director F6人工E2E门禁未闭合；当前不是无条件发布就绪。

最新执行记录：`UPSTREAM_V127_STAGE5_20260908.md`。阶段5选择性UI适配及源码门禁完成，停在阶段6批准前；未部署或完成Edge验收。

最新候选源码记录：`UPSTREAM_V127_STAGE34_20260908.md`。阶段3–4源码门禁通过；后续UI适配、发布、部署和Edge尚未完成。

本目录保存 `<project-root>` 的持久工程记忆，用于跨任务、跨智能体和上下文压缩后的连续开发。

## 开始任务时

1. 始终先读 `CURRENT_STATE.md`。
2. 涉及历史、合并或旧方案时读 `TIMELINE.md`。
3. 涉及修复、部署、数据、协议或安全时读 `LESSONS.md`。
4. 涉及继续开发或排期时读 `BACKLOG.md`。
5. 涉及备份、迁移、恢复或部署时读 `BACKUP_REGISTRY.md`。
6. 更新记忆前遵循 `MAINTENANCE.md`。

## 证据优先级

从高到低：

1. 当前代码、Git 状态和当前配置；
2. 本轮运行测试、浏览器证据、API 响应和日志；
3. 可复现的历史测试输出和 Git 提交；
4. 正式专题文档；
5. `.claude/`、交接报告、旧分析和会话摘要。

后一级不能覆盖前一级。`healthy`、编译成功、日志存在、构建完成或历史“已验证”都不能自动升级为当前业务可用。

## 当前接管结论

- 项目已经完成一次大规模上游合并，并保留了 DashScope 图片、视频和多参考图等本地资产。
- 百炼目录、支持状态、Wan 2.7、HappyHorse 1.1、Wan 3.0、精确 SKU 计费和 OSS 结果物化已经完成阶段 1–8 验证；阶段 9 已完成本地 Git 发布固化，尚未 push、创建 PR 或触发远程 CI。
- Qwen Image 3.0 Pro 已在阶段 12 完成专属能力、严格同步 Provider、Ready 目录和 Microsoft Edge 无费用验证；模型保持停用、未定价，尚无真实付费生成，因此不能宣称账号运行调用已验证。
- 阶段 13 已完成统一错误分类、requestId、前端诊断编号和任务/Provider/资源安全关联日志；本地全量测试、构建、保留卷部署及无费用 401/400 运行验证通过，远程 CI 尚未触发。
- 阶段 14 已修复影策 Codex 插件的 Vite/Origin、URL Token、长任务超时和 Compose Provider 插件开关合同；源码与主机全量验证通过，插件安装态新线程仍待验证。
- 阶段 15 已完成系统渠道统一报价、生产 CORS fail-closed、GORM SQL 日志脱敏和未接入价格死结构清理；本地全量、部署与无费用日志验证通过。
- 阶段 16 已完成最终本地全量门禁、恢复点、镜像标签、Microsoft Edge 人工冒烟和交付报告；当前停在远程 push/PR/CI 授权门禁。
- 运行容器、GitHub 凭据、数据库内容和远程上游仍必须在每次操作前重新核验，不能沿用旧快照。
- 数据、外部 API、GitHub 写入、付费生成和全局代理变更均需单独授权。

## 文件职责

- `CURRENT_STATE.md`：当前代码、Git、运行环境和已确认问题。
- `TIMELINE.md`：功能、合并、事故和修复的时间线。
- `LESSONS.md`：踩坑、根因、解决策略和不可违反的边界。
- `BACKLOG.md`：待办、项目计划、执行顺序和验证门禁。
- `BACKUP_REGISTRY.md`：数据库、密钥、工作树和部署基线备份的权威登记。
- `STAGE5_RELEASE_CANDIDATE.md`：阶段 5 候选镜像、验证、迁移顺序、回滚和 NO-GO 清单。
- `STAGE9_RELEASE_HANDOFF.md`：阶段 9 发布结果、真实验证、历史问题完成度和剩余门禁。
- `STAGE16_FINAL_HANDOFF.md`：阶段 16 最终本地发布候选、测试、数据、镜像、回滚和远程发布门禁。
- `UPSTREAM_BAILIAN_CONTRIBUTION.md`：上游百炼贡献的背景、Issue #332、五项架构决策、PR 拆分、测试门禁和维护者回复判定规则。
- `UPSTREAM_MERGE_PLAN_20260828.md`：官方主线增量合并的用户授权、阶段门禁、恢复边界、remote 规范和冲突处理顺序。
- `UPSTREAM_INCREMENT_AUDIT_20260829.md`：`ab89c05..4f07daa` 新增量、六个显式冲突、语义冲突、解决合同和下一阶段验证门禁。
- `UPSTREAM_TAIL_AUDIT_20260830.md`：`4f07daa..c8b60ce` 尾差、17 个显式冲突、技能包/短剧/后台设计风险、逐文件解决合同和验证门禁。
- `UPSTREAM_FIX_TAIL_AUDIT_20260830.md`：`c8b60ce..4ba9694` 八个修复提交、当前已解决索引的 3 文件/15 区块模拟、协议/首页阻断项和阶段 B2 建议。
- `ROUTE_B_DEPLOYMENT_20260901.md`：固定 `4ba9694` 本地候选的构建标识、恢复点、克隆迁移、保留卷升级结果和 Microsoft Edge 人工验收清单。
- `UPSTREAM_V125_PLAN_20260904.md`：固定官方 `v1.2.5`、双轨来源/单树交付模型、阶段 0–9、冲突保护合同和 Astravia 延后边界。
- `UPSTREAM_V125_AUDIT_20260904.md`：`4ba9694..v1.2.5` 完整 506 文件增量、47 个显式冲突/149 区块、功能族和阶段 4 解决合同。
- `MAINTENANCE.md`：记忆更新规则。
