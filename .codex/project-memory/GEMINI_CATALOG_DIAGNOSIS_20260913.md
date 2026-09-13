# Gemini 系统渠道拉取后只读：历史诊断

当前已由ad3efa10/cc7d0b63完成实现、全量预验证及部署，停在用户Edge门禁。权威交付见MODEL_CATALOG_UNIFIED_REPAIR_20260913.md；下方“未修复/未批准执行”为当时诊断阶段记录。

完整优化方案已收敛到 MODEL_CATALOG_UNIFIED_REPAIR_PLAN_20260913.md：不只补Gemini或新增状态枚举，而是分离发现/配置/执行/实测，统一手工和导入路径，并通过现有插件合同承载模型支持；其他百炼模型风险与实施阶段见该文档。当前仍仅设计，未执行。

2026-09-13 用户报告新增 Gemini 系统渠道、拉取全部模型后全部停用，编辑操作不可用。截图显示能力/协议待配置、planned、未定价、查看按钮，共50条。

## 源码已确认

- 有效树 `.local/worktrees/dashscope-plugins`，HEAD 81896044；已有前一轮用户验收文档修改保留。本轮未独立读取生产渠道凭据/数据库或操作浏览器。
- `backend/internal/service/channel_model_catalog.go:197` 通用发现全部设置 planned；原生 Gemini 目录字段 supportedGenerationMethods / displayName 未建模，解析器仍主要读取 model_type / supported_endpoint_types / display_name，catalogProtocol 没有 Gemini 原生映射。
- `channel_model_catalog_plugin.go:16` 的发现补全仅注册 bailian.New()，没有 Gemini 补全，因此 Gemini 条目停留在未知能力/协议和 planned。
- `channel_model_catalog_sync.go:49` 新导入模型 disabled/unpriced 是安全默认；真正的配置阻断是 planned。前端 channel-model-support.ts 把 planned 设为只读，editor 同时 disabled/inert；后端 requireReadyChannelModel 在保存/测试/批量更新也拒绝，不能仅解除前端禁用。
- `channel_models.go:399` 同渠道重复键拒绝新增，不能依靠同名手工新增修复已有只读项；不可建议删除记录规避。
- 代码历史 64015db3 / 62cf5d83 引入通用 planned 与管理只读门禁，先于本轮 cd0b4340 画布多图增量。属于此前目录治理接入缺口，不是画布开关直接造成。
- 存在 Gemini Image/Veo 执行代码，文本还有声明式插件门禁，不能由“已有部分执行器”推断所有 Gemini 模型均可执行。

## 官方合同及待核实项

官方 https://ai.google.dev/api/models 本轮已读取：原生 models.list 使用 displayName / supportedGenerationMethods，默认每页50条，通过 nextPageToken / pageToken 翻页。当前拉取只发一次 GET，不解析 nextPageToken；截图50条可能仅第一页。未读取该渠道实际响应，因此尚不能确认账号确实还有下一页，也未确认其上游是否为Google原生或第三方转发。

## 建议与授权边界

### 来源复核与分页澄清（2026-09-13）

- 对比实际采用的官方 v1.2.7 固定提交371bc68与当前分支：官方 channel_models.go 导入默认停用/未定价，但 manager 显示“编辑”，无本地 requireReadyChannelModel 限制。本地64015db3（2026-08-26）加入通用 planned 和后端只读校验，62cf5d83加入前端只读；两提交均不是371bc68祖先。因此“无法编辑启用”的直接原因是本地开发回归，而非官方原版已有同样锁死行为。
- 官方371bc68的目录解析同样缺少Gemini原生字段完整映射和nextPageToken分页。这些是上游遗留不足；不能与本地新增的编辑阻断混为一谈。此对比针对已采用固定上游版本，未声称最新远程HEAD也如此。
- 用户把页面显示改为100条/页。当前manager.tsx:242通过filteredItems.slice本地切页，底部total=filteredItems.length；此设置不改变后端向Google发出的目录请求。用户看到了已导入的全部50条，不能据此证明上游仅50条。是否实际漏页仍待原始响应nextPageToken证据；分页不足与模型只读是独立问题。
- 当前通用路径全部置planned，后续无目录补全的其他新渠道（包括仅返回ID的OpenAI兼容目录）也可能同样被锁；只补Gemini静态白名单不能根治。

### 建议的系统性修复范围（未批准执行）

1. 分离发现事实、配置完成度、执行支持、管理员启用/定价和账号真实可用性。缺元数据属于待配置，不等于明确planned/unsupported/deprecated。
2. 复用现有协议/插件Registry作为执行能力依据，按需要扩展插件目录解析/模型匹配元数据；不新增平行Gemini宿主注册表。目录GET与分页仍走统一安全出站。
3. 已有可执行协议允许候选受校验配置；后端验证精确能力/操作、插件可用性及任务路径后才放行启用和路由。未知自定义ID可以选择已有兼容协议；没有适配器仍明确阻断，不以名字或generateContent一个方法推断全部输出能力。
4. Gemini解析原生字段，按nextPageToken翻页，去重并限制页数/总量/循环令牌；中途失败不得声称已拉取全部，不将缺页视为上游删除模型。
5. 原地补全已有未配置候选，不删除重建、不更改模型ID/价格/启用/引用；已人工配置项保留，明确计划支持/不支持项不被批量转ready。
6. 覆盖Gemini分页与50/100显示分离、通用OpenAI未知ID候选配置、没有适配器/插件停用、重复拉取、已配置保护、百炼与Qwen能力回归；修复后停Edge人工门禁，真实调用单独授权。

- 沿统一协议/插件能力注册机制补 Gemini 目录解析、分页和精确可执行合同映射；未知/不支持模型不直接批量改 ready。
- 区分“已发现、等待管理员配置”与“项目确实尚未支持”；为已有候选提供受后端合同校验的配置/转化流程。保留显式价格和启用门禁。
- 重新拉取只补全未配置候选，保留已配置模型、价格、启用状态和引用；补 Gemini/OpenAI/百炼互不回归测试。
- 本轮为诊断，未实施修复、未改模型/价格、未测试模型或发起付费调用；仅按仓库约定更新未提交的项目记忆。
