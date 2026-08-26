# 阶段 9 发布交付报告

## 交付结论

阶段 1–8 的百炼模型目录、支持状态、Wan 2.7、HappyHorse 1.1、Wan 3.0、精确 SKU 计费、阿里云 OSS 物化和 Edge 真实验证已经完成。阶段 9 已将这些变更整理为可审阅的本地 Git 提交，并保留明确的未完成问题清单。

本报告不把“主线已完成”扩大解释为“整个仓库所有历史问题均已修复”。用户将继续手工测试其他功能；任何新发现另行诊断和授权。

## 当前运行基线

- Backend 与 Web 均运行最终候选镜像并保持 healthy、RestartCount=0。
- 数据卷未删除、未替换；禁止使用 `docker compose down -v`。
- 运行目录为 319 个渠道模型：11 Ready、308 Planned。
- Wan 3.0 模型和临时价格档已恢复停用、未定价、单价 0。
- 保留 1 个成功真实任务、1 个 settled 账务、1 个 ready 阿里云 OSS Resource、1 个 confirmed Asset 和 1 个成功 API 调用审计。
- 用户余额为 98 积分，预留为 0。

## 真实验证证据

- 模型：`wan3.0-video`。
- API 模式：DashScope 原生异步创建和任务轮询。
- 创建请求：1 次；HTTP 200；任务 attempts=1；没有创建重试。
- 实际参数：adaptive、480P、2 秒、audio=false、watermark=false、seed 省略。
- 账务：按秒 1 积分，实际结算 2 积分。
- 结果：MP4，805,526 bytes，阿里云 OSS 乌兰察布 Endpoint，Resource ready。
- 前端：Asset confirmed，使用 `resource:` 持久引用；资源 Range 读取返回 206。
- Edge：视频正常播放，播放器显示 0:02；用户确认无声、无水印。
- 偏差：刷新后画幅从原计划 16:9 回到 adaptive；未补发第二次任务。

## 测试与构建

- Backend：隔离 Linux CGO `go test -count=1 ./...` 通过；测试不挂载运行数据、备份或密钥。
- Web：主套件 449/449、跨 Runtime 1/1、TypeScript 和容器生产构建通过，无 Rolldown panic。
- Canvas Agent：TypeScript 构建通过；阶段 4 曾有 286/286 历史通过记录，但阶段 9 当前 Windows 复验在 `dreamina-cli-process.test.ts` 稳定出现 4 个进程树清理失败并使测试进程不退出。失败涉及 oversized output、取消清理、早期 receipt 清理和 progress timeout 后的 `taskkill`/child close 语义；不能写成当前全绿，需 Linux CI 复核或单独修复。
- 后续阶段 10 已确认上述 4 个失败是 Codex 受限进程无权执行 `taskkill` 的测试环境假阴性；主机权限专项 12/12、全量 287/287 和构建通过。阶段 9 此处保留为当时历史结论，不再代表当前状态。
- Manifest：80 个百炼官方补充模型，11 Ready；生成前后哈希稳定。
- Git 静态门禁：`gofmt -l backend` 为空，`git diff --check` 通过；最终密钥/本机路径扫描见阶段 9 执行记录。
- 远程 GitHub Actions 尚未触发；本地通过不能写成远程 CI 已通过。

## 数据保护与回滚

- 阶段 0、阶段 6、阶段 7 和阶段 8 的 verified/protected 恢复点已登记。
- 阶段 8 调用前恢复点：`BKP-20260826-185248-STAGE8-PRECALL`。
- 阶段 8 调用后恢复点：`BKP-20260826-201444-STAGE8-POSTCALL`。
- 公开仓库只保存脱敏 ID、哈希和验证状态；本机精确路径位于 Git 忽略的 `.local/backups/*/backup-pointer.md`。
- 数据库恢复不能撤销供应商费用，也不会自动删除或恢复 OSS 对象。
- 旧 Backend/Web 镜像保留本地回滚标签；回滚镜像不能修复共享卷权限问题。

## 历史问题完成度审计

| 问题 | 当前结论 | 证据或剩余边界 |
| --- | --- | --- |
| GitHub“无法获取 PR 状态” | 已解释，非仓库代码故障 | 受限进程看不到 Windows Keyring；主机环境认证正常，当前没有对应 PR。 |
| 数据库无法读取确认 | 已解决 | verified 快照、独立恢复、同 UID/GID 只读监控和完整性检查均通过。 |
| 阿里云 OSS 配置是否生效 | 已修复并真实验证 | 真实视频 Resource ready、Edge 可播放、Asset 使用 `resource:`。 |
| 百炼拉取只有约 240 个且没有视频 | 已修复 | 上游 241 + 官方补充 80；最终 319，其中视频 78。 |
| 重复拉取持续“补齐 239” | 已修复 | 空目标字段不再写回；重复拉取新增 0、补齐 0，目录摘要不变。 |
| Ready/Planned 状态和非 Ready 误启用 | 已修复 | Backend 强制门禁、后台只读状态、普通用户目录过滤和测试均通过。 |
| Catalog/Admission 能力分裂 | 已修复（系统渠道） | 二者共用 `effectiveChannelModelCapability`；缺失/损坏/不完整 fail closed。 |
| Wan 3.0 Adapter | 已修复并真实验证 | 原生端点、轮询、计费、OSS、播放完整通过。 |
| Wan 2.7 / HappyHorse Adapter | 代码和 mock 已完成，真实账号调用待验证 | 模型专属能力和 payload 测试通过；没有额外付费调用授权。 |
| DashScope 原生路径错误 | 已修复 | 保留配置的账号/区域 host，只替换 `/compatible-mode/v1` 为原生 `/api/v1` 路径。 |
| 结果 URL 泄露 API Key/自定义头 | 已修复 | 结果下载不再转发模型凭据，并执行出站 URL 校验。 |
| `generation-*` 素材不进入 OSS | 已修复并真实验证 | 识别 image/media 存储源，上传失败 fail closed；真实 Resource/Asset 成功。 |
| `priceConfigured=false` 仍可能路由 | 已修复可用性边界 | 目录、路由和价格档同时要求 enabled、PriceConfigured 和正价格。 |
| 系统渠道精确 SKU 计费 | 已修复并真实验证 | operation/vquality/videoSeconds 精确匹配，真实任务结算正确。 |
| Rolldown worker panic | 已修复 | 移除 Bun 主进程内 Vite server，增加 panic guard，最终全量无 panic。 |
| Windows CGO SQLite 测试失败 | 已修复 | 独立 Linux CGO test target 和 CI 工具链检查，不改生产驱动。 |
| 未使用 modernc SQLite 依赖漂移 | 已修复 | go.mod/go.sum 回到 go-sqlite3 生产合同。 |
| Create 缺少生成声音/水印 | 已修复 | 能力驱动输出控制、摘要、449+1 测试和真实无声/无水印请求。 |
| Qwen 图片参数/UI 旧问题 | 部分解决，未完成 | 修正真实 Create 组件和 missing/empty/wildcard 语义；Qwen 图片仍 Planned，Provider 参数合同和 Edge 真实验证未完成。 |
| DashScope 图片通用能力夸大 | 未完成但已隔离 | 图片模型保持 Planned，不可定价/启用；晋升 Ready 前需模型专属序列化与测试。 |
| 通用视频 fallback 夸大能力 | 部分解决 | Ready Wan/HappyHorse 已脱离通用 fallback；其他 Planned/手工模型仍需逐模型收敛。 |
| 任务错误响应和结构化日志 | 部分解决 | `failService` 已安全投影未分类错误；`POST /tasks` 等部分路径仍把 service error 交给 `fail`，全链路日志未完成。 |
| 系统渠道 `/model-catalog/quote` | 未修复 | feature off 分支仍返回 501 TODO。 |
| 插件备用端口/Vite/Origin | 未修复 | skill 仍误写 Next、使用无效 `CANVAS_URL`，未对接 `FRAMEFIELD_TRUSTED_WEB_ORIGINS`。 |
| 插件 90 秒超时 | 未修复 | 低于 Canvas Agent 35 分钟生成续接合同。 |
| 插件 `agentToken` 查询参数 | 未修复 | 文档/skill 仍采用已过时且可能泄露的 URL token 协议。 |
| Compose 插件开关透传 | 未修复 | `ENABLE_PROVIDER_PLUGINS` 只在 local Compose 透传。 |
| Server CORS 默认 `*` | 未修复 | `docker-compose.server.yml` 仍默认允许所有 Origin。 |
| 文档站构建 | 未修复，用户已暂缓 | `docs/package.json`、`source.config.ts` 和多篇链接页面仍缺失。 |
| 远程 GitHub CI | 待验证 | workflow 已补 Backend/Web/Canvas Agent 门禁，但尚未 push/PR 触发。 |
| Canvas Agent Windows 进程清理测试 | 当前失败，未修复 | 4 个 `dreamina-cli-process` 用例失败且遗留子进程；构建通过。按用户要求本轮只记录，不改生产进程代码。 |
| SQLite 卷诊断权限 | 已恢复，自动防护未完成 | DB/WAL/SHM 已为服务 UID/GID 100:101、0660；以后诊断容器必须同 UID/GID。 |
| Create 刷新后画幅回默认值 | 未修复 | 真实调用从 16:9 回到 adaptive；需要单独产品状态设计。 |
| 音频/水印控件视觉像普通方框 | 可用但交互表达仍可优化 | 整块按钮可点击且摘要准确；可后续改为更明确的 Switch/Checkbox。 |
| `LogicalModelPriceSKU` 死结构 | 未修复 | 类型仍存在但未进入实际 schema/repository/service 主链。 |
| `.claude/settings.local.json` 宽泛权限 | 未修复，本地治理项 | `.claude/` 已 Git 忽略；应由用户在 Claude 环境单独收敛。 |

## 不纳入本次发布的后续项

- Qwen 图片模型 Ready 化和真实生成。
- Wan 2.7、HappyHorse 真实付费回归。
- 商业价格矩阵、批量启用模型和运营策略。
- 插件协议、CORS、文档站和系统 quote 的单独修复。
- push、创建 PR 和远程 CI；这些是外部写入，需要独立授权。

## Git 提交

阶段 9 从 `main@0b202a1` 起按职责拆分为以下本地提交：

- `64015db feat(backend): 百炼模型目录 - 增加支持状态与视频适配`
- `62cf5d8 feat(web): 百炼模型管理 - 增加支持状态与生成控制`
- `f97adde test(build): 发布门禁 - 固定工具链与隔离验证`
- 文档、项目记忆和百炼技能由本报告所在的最终文档提交固化。

本阶段没有 push、没有创建 PR，也没有触发远程 GitHub Actions。
