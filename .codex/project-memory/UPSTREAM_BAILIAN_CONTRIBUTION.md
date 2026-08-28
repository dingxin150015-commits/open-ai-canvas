# 上游百炼能力贡献专题记忆

更新时间：2026-08-28

## 1. 目标和当前结论

本专题记录从本地百炼能力开发、上游差异研究、贡献排序、设计 Issue 提交到维护者决策等待的完整链路。它用于上下文压缩或换代理后继续工作，不能替代提交前对官方最新 `main`、Issue 和官方百炼文档的实时复核。

当前结论：官方 Issue [#332](https://github.com/ddcat-ai/open-ai-canvas/issues/332) 已提交，但维护者尚未回复。决策状态为 `pending`。在五项架构决策确认前，支持状态 Schema、百炼 Adapter 和版本化 Manifest PR 全部为 **NO-GO**。

## 2. 本地资产与上游差异

### 本地已形成的可贡献资产

- 80 项版本化百炼补充目录，以及 `ready/planned/unsupported/deprecated` 支持状态；
- 非 Ready 在启用、定价、测试、逻辑线路和任务准入边界的 Backend 强制门禁；
- Qwen Image 3.0 Pro 原生同步 Adapter、模型专属能力和 mock/fixture 测试；
- Wan 3.0 All-in-One、Wan 2.7、HappyHorse 1.1 的模型专属能力、序列化、异步状态与拒绝测试；
- 目录同步保留完整元数据、乐观并发、安全补齐和重复拉取幂等；
- 精确 SKU、OSS 结果物化和一次获授权的 Wan 3.0 真实运行证据。

### 研究时确认的官方现状

- 2026-08-27 只读快照为官方 `main@5289bef4`；该哈希会漂移，任何新分支创建前必须重新 fetch/API 核验。
- 官方百炼 discovery 已列 Qwen Image、Wan 2.7、HappyHorse 等静态目录项，但没有明确区分“官方已发现”和“项目可执行”；当时未包含 Wan 3.0。
- 官方 `feature` 当时落后 main 414 个提交且独有 0，不适合作为贡献基线。
- 本地 `main` 包含长期产品化开发和大规模历史合并，不能直接推送成上游 PR。

## 3. 用户确认的贡献顺序

用户把原贡献顺序调整为：

1. **第三批：百炼能力贡献**；
2. **第一批：低耦合工程贡献**，包括 Provider 插件 Compose 开关、Codex 插件本机连接/超时、系统渠道统一报价、CI 门禁；
3. **第二批：产品与安全贡献**，包括 Create 声音/水印、Create 参数持久化、数据库日志安全和生产 CORS。

本专题只约束第一项“百炼能力贡献”。后两批仍需在各自提交前重新对照官方最新代码，避免上游已经实现或架构发生变化。

## 4. 为什么先提设计 Issue

本地百炼实现规模大、跨数据库、后台、目录、任务准入、Provider、计费和测试。如果直接提交大 PR，会把以下尚未对齐的公共合同强加给上游：

1. 支持状态存放位置；
2. 非 Ready 模型在后台如何显示；
3. 人工协议配置能否视为 Ready；
4. Adapter 属于 Provider、Protocol 还是新的 Executor Registry；
5. 官方目录采用生成式 Manifest 还是手写 Go 列表。

因此先提交中文设计 Issue，只陈述问题、证据、建议和可选方案，不提前创建 Schema 或 Adapter 分支。

## 5. Issue #332 提交记录

- 仓库：`ddcat-ai/open-ai-canvas`；
- 标题：`[Feature] 百炼模型目录需要区分“官方已发现”和“项目可执行”`；
- URL：https://github.com/ddcat-ai/open-ai-canvas/issues/332
- 创建时间：2026-08-27；
- 作者：当前已认证 GitHub 账号；
- 正文包含：问题、风险、官方代码证据、四态支持状态、Ready 十项条件、数据兼容、版本化 Manifest、五项架构问题、六步 PR 拆分、非范围和测试建议；
- 未创建分支、PR、push 或 CI；
- 尝试添加 `enhancement` 标签被 GitHub 以权限不足拒绝。Issue 无标签是权限边界，不影响正文有效性；禁止重复尝试或声称标签已添加。

## 6. 等待维护者确认的五项决策

### 决策 1：支持状态真相位置

选项：仅 `provider.Model`、仅 `ChannelModel`、或 Provider 给目录默认值而 ChannelModel 保存管理员有效状态。

本地推荐第三种：Provider 表达官方目录事实，ChannelModel 表达当前系统的持久状态。

### 决策 2：非 Ready 后台表现

选项：完全隐藏、只读展示且禁止定价/启用、或允许管理员强制覆盖。

本地推荐只读展示，既保留官方目录资产，也不制造可执行假象。

### 决策 3：人工 Ready

选项：完整人工协议和能力可视为 `ready/manual`，或只有内置 Adapter 才可 Ready。

本地推荐允许 `ready/manual`，但仍需通过协议、能力、价格和任务准入校验。

### 决策 4：Adapter 位置

选项：`backend/internal/provider/bailian`、`backend/internal/protocol`、或新的 Provider Executor Registry。

本地推荐优先复用官方当前 Provider/Protocol 注册体系；声明式协议表达不了媒体顺序、角色和跨媒体限制时，再增加模型级 Hook，避免继续扩张通用 service switch。

### 决策 5：Manifest 维护方式

选项：脚本生成的 JSON/Go Manifest、手工 Go 列表、或只维护 Ready 模型。

本地推荐版本化生成 Manifest，提交生成器、公开官方链接和确定性测试，不提交完整官方文档镜像。

## 7. 维护者回复的有效判定

### 有效确认

同时满足：

1. 回复者 `author_association` 为 `OWNER`、`MEMBER` 或 `COLLABORATOR`；
2. 明确逐项选择五个架构方案，或明确表示接受 Issue 中的推荐方案；
3. 若要求替代架构，描述足以确定首个 PR 的边界。

### 不能单独视为确认

- Issue 仍为 open 或被关闭；
- 添加标签、里程碑、指派或项目看板；
- 点赞、订阅或其他 reaction；
- 普通 `CONTRIBUTOR`/`NONE` 用户的回复；
- “欢迎 PR”“可以试试”“看起来不错”等笼统意见；
- 其他 Issue/PR 的自动交叉引用；
- 没有正文决策的提交或分支引用。

### 状态记录枚举

- `pending`：没有有效维护者回复；
- `partial`：只回答部分问题或笼统欢迎 PR；
- `confirmed`：五项均明确，可进入首个 PR 设计；
- `alternative_requested`：维护者要求另一架构，先修订方案；
- `rejected`：方向被拒绝，不继续提交；
- `superseded`：维护者用其他 Issue/PR/实现替代本方案。

任何未决项都继续保持对应代码范围的 NO-GO。

## 8. 2026-08-28 实时状态与误信号排除

通过 GitHub API 读取 Issue、comments 和 timeline，确认：

- state：open；
- comments：0；
- labels：0；
- assignees：0；
- milestone：无；
- `updated_at=2026-08-27T09:13:58Z`，与创建时间相同；
- 没有维护者对五项决策作答。

时间线出现 PR #335 的 `cross-referenced`。进一步核对后确认：

- PR 主题是“3D 导演台 - 补齐非 P0 分层交互与回归保护”；
- 作者 `echoD886` 的身份是 `CONTRIBUTOR`；
- PR 无 Issue 评论、无 Review；
- 异常重复正文偶然包含 `#332`，触发 GitHub 自动引用；
- 与百炼目录、支持状态或 Adapter 完全无关。

结论：这不是维护者回复，也不是关注或批准信号；Issue 决策状态保持 `pending`。

## 9. 跟进节奏和外部写入边界

1. Issue 提交后先等待 2–3 个工作日；
2. 3–5 个工作日仍无回复时，可以准备一条简洁评论，把五项问题压缩成编号选择；
3. 评论属于 GitHub 外部写入，必须再次获得用户明确授权；
4. 只发送一次结构化跟进，不高频催办、不 @ 多名维护者；
5. 超过 7 天仍无回复，不自行假定推荐方案获批，可转而先贡献不依赖该架构的低耦合 PR，或等待用户决定。

## 10. 获确认后的 PR 拆分

每个 PR 都从**当时官方最新 main** 新建干净 `contrib/*` 分支，不直接使用本地 `main` 或旧 Fork feature：

1. 支持状态、历史回填、后台只读行为和强制门禁；
2. Qwen Image 3.0 Pro 原生同步生成/编辑；
3. Wan 3.0 All-in-One；
4. Wan 2.7 T2V/I2V/R2V；
5. HappyHorse 1.1 T2V/I2V/R2V；
6. 版本化百炼 Manifest 和幂等同步。

不得把本地约 5,000 行以上的百炼实现整体作为一个 PR。每个模型 PR 只携带该模型的能力、Serializer、创建/查询/终态解析、拒绝路径和测试。

## 11. 上游测试与安全门禁

### 支持状态基础 PR

- SQLite/PostgreSQL 迁移；
- 历史完整配置与人工配置回填；
- 非 Ready 保存、定价、启用、测试、逻辑线路和任务准入拒绝；
- 后台只读展示；
- 管理员并发保存保护；
- 重复目录同步新增 0、更新 0；
- Deprecated/退休模型不自动恢复；
- Backend 全量、Web 测试/类型检查/生产构建。

### 单模型 Adapter PR

- 精确模型 ID、operation、同步/异步生命周期；
- 能力与 Provider payload 一致；
- 合法序列化 fixture；
- 非法参数、媒体类型/数量/大小/组合拒绝；
- 不静默钳制、替换、丢弃或扩展参数；
- 临时 URL、SSRF 和结果物化边界；
- Catalog、Admission、计费和 Provider 共用同一有效合同。

默认只运行 fixture/mock，不执行真实付费生成。真实百炼调用必须重新声明账号/区域、模型、端点、参数、预计费用和重试次数，并取得用户单次授权；本地历史 Wan 3.0 授权不可复用。

## 12. 可复用经验与踩坑

- 目录发现不是执行能力；Ready 必须是端到端可证明合同。
- 先设计 Issue、后小 PR，能避免在数据库和 Provider 层形成维护者不接受的公共结构。
- 贡献材料全部使用中文，代码符号和必要协议字段保持原名，便于中文背景维护者审阅。
- GitHub 标签写权限不足不影响 Issue 内容；权限错误应记录一次后停止重试。
- Issue 状态检查必须读取 comments 和 timeline，并打开交叉引用来源；仅看页面顶部或通知可能误判。
- 维护者沉默不是同意；普通贡献者互动也不能替代仓库决策者。
- 上游分支和官方代码会快速变化；历史 commit 只作为当时证据，开始每个 PR 前必须刷新并重新做差异审计。
- 本地官方文档镜像是开发真源，但不能整体提交上游；PR 使用公开链接和最小元数据。
- 外部评论、push、PR、CI 和真实模型调用都是独立授权边界，不能从“允许提交 Issue”自动扩大权限。
