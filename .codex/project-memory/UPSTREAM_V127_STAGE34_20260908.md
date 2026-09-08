# v1.2.7官方集成检查点（阶段3–4）

## 范围与固定输入

用户批准阶段3–4：从stable创建独立候选，合并固定官方v1.2.7，解决冲突并通过对应源码门禁后形成本地双父merge检查点。本文件随该检查点提交；不代表最终版本验收。

- 候选：`codex/upstream-v1.2.7-20260908`，独立worktree。
- 第一父基点：`5c1203eaf4876b697a4d89236f4979c03b4e34f5`。
- 官方第二父：`371bc68a5951c54f793cdfd37a2beb4d6d04a576`；tag对象 `fc85af2022a202ac8272d42e2f5fc93526f6fc9e`，未签名。
- 原UI留存 `codex/ui-v125-followup@b6254b4` 和稳定分支未移动；未重放e027aa4/f5fadb7/bcb6c6d。
- VERSION暂定 `v1.2.7+dingxin.1`，尚未发布或部署。

## 合并处理

真实63冲突文件/184区块与阶段1模拟一致。保留本地能力/Ready状态/精确SKU/错误诊断/OSS归属等合同，接入官方新模型编辑器、批量导入删除、尺寸预设、任务租约/恢复、资源播放副本、编辑器与懒加载UI。

复杂前端先保存原始冲突，在临时审计目录对base/ours/theirs统一格式后三方合并，消除格式噪声；其余语义冲突逐项处理。没有整目录/整文件盲选一方。

- 新导入复用完整CatalogItem，保留支持状态和能力；新编辑器与批量操作不绕过非Ready只读边界。
- Provider新媒体policy保留DashScope公网URL；任务租约与本地诊断日志合并；恢复分类在内部读取精确结构化错误，日志/对外消息继续脱敏。
- Create采用官方懒加载/会话结构，保留Backend-only、metadata/reasoning/clientOperationId、恢复取消、按模型草稿、settingsReady、声音水印；官方偏好为缺少模型草稿时的回退。
- 插件生成器重复注册和上传界面重复片段已收敛；新的控件类型/事件接口同步。
- 公告引用排除采用官方统一接口，仍只排除当前公告/自身草稿；最终官方软禁用用户行为保留，不恢复中间版硬清库。

## 已执行验证

- 隔离Linux CGO Backend全量通过，包括数据库迁移单元测试、Handler、repository、protocol、service及新增批量删除原子拒绝、导入状态与媒体policy回归。未读取真实数据库。
- 前端Node24/Bun1.4.0隔离Linux环境：干净TypeScript通过；默认测试链 `27+4+16+1174+33+1=1255 pass / 0 fail`；生产构建通过（39.15秒）。原有panic guard保留。
- 跨Runtime测试引用Agent源码，其zod与Web锁文件均为3.25.76且完整性相同；仅在测试镜像暴露这一相同依赖，未启动真实Agent。
- 修改的JSON解析、插件生成脚本语法、diff空白及冲突标记检查通过。
- 测试镜像：Backend `open-ai-canvas-backend-test:v127-stage4`，manifest list `sha256:1b19239054cd4c9627d0a51d8723b2d19e7746d7b09dcd503116f0b34d86087b`；Web `open-ai-canvas-web-test:v127-stage4`，manifest list `sha256:94a4f2c5934d7970a634d39dc251919edfc5e7fbe9f14e06d0e01bf726332cb2`。均非部署镜像切换。
- 完整原始日志/格式归一前后材料在主工作区 `.local/audits/v127-20260908-110746/`，主要为backend-second.log和web-fifth.log。

## 已定位的验证问题

第一轮Backend失败源于官方报价fixture缺少本地必需BillingMode，以及脱敏错误消息影响任务恢复分类；已修正并回归通过。旧Create/媒体源码断言随官方组件拆分调整，保留能力和播放边界。

Windows依赖安装持续高CPU后仅停止本任务安装进程，保留文件并转Linux验证。旧宿主机tsbuildinfo复制进镜像曾产生大面积过期类型诊断；断网同镜像禁用增量缓存后只剩两处真实Switch遗漏。已修正并将tsbuildinfo排除出Docker上下文；最终验证使用无增量缓存模式。

## 尚未执行与下一门禁

- 阶段5三个UI提交仍待选择性重放/跳过；失败UI2布局不得整包重放。现有布局/Runtime/报价复验问题不因此关闭。
- 阶段6最终全量/吸收清单、阶段7fork候选push、阶段8数据副本双启动及部署、用户Edge、阶段9稳定晋升均需后续批准。
- 官方Director P0 E2E失败根因仍待定位；未操作浏览器，不把源码通过视为Edge通过。
- Schema6→8两种历史谱系必须在真实数据副本确认；当前仅单元测试通过。未挂载、暂停或迁移生产卷，未启用Host Updater。
- 运行环境仍为原Web UI2和v1.2.5 Backend，未推送、未部署、未执行真实Provider/支付/OSS写入。
