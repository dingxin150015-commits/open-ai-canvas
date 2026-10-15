# 待办与项目计划

## 当前执行顺序（2026-10-14 23:18 CST）

1. 自动化、实际PG/Redis、六谱系迁移、候选镜像、隔离预览与旧Schema9恢复已通过；最后Web2491+5全部通过。旧未解/全量失败条目只作为历史。
2. 用户在 http://localhost:3001 按 `docs/plans/edge-v166-acceptance-20261014.md` 验收；F6仅提供可选预览故障注入，未执行则仍test_only_pending。
3. 人工确认后，重核生产/远端和空闲状态、新鲜私有恢复点，正式commit/同SHA镜像、保留真实卷晋升、运行/数据/旧浏览器缓存复核，非强制push用户fork并核对CI。发布权限已有，不重复征求；不能省略手验与数据门禁。
4. 最后更新交接和根入口。旧A/候选/工作树/卷/镜像先保留，未经具体删除批准不清理；上游PR仍草案，贡献时重新核官方主线差集。

2026-10-14 17:43 CST：61显式冲突已清零，Backend Linux全仓、Agent328/构建与Web类型/lint/首轮构建通过；Web最新全量2486/4，最后四项收口中。真实迁移、PG/Redis依赖、最终候选镜像、Edge、恢复点/生产部署/推送未完成。下面旧未解状态为历史，按执行记录最新段继续。

2026-10-14 15:21 CST 当前进度：61 个显式冲突已全部解开；Web 类型/lint 与旧场景/新预演专项通过。全量测试、部署配置语义、真实迁移副本和生产恢复点继续执行，未部署/push。以下启动复选框按最新执行记录解释。

## v1.6.6 A‴ 当前门禁（更新：2026-10-13 01:31 CST）

- [x] 2026-10-13 13:53 CST 用户选择采用上游删除画布版本历史/预览/恢复 UI；实现需移除旧调用链而保留旧数据/备份，并在交付中披露行为变化。

- [x] 官方正式 v1.6.6 标签及剥离 SHA 核实，旧 A″ 未验收树封存为仅本地 WIP `4323c759`；新隔离 A‴ 已开始无提交合并，旧 A/A′/A″ 和生产均未动。
- [ ] 依次解决 61 个显式冲突（2 个删除/修改），并审计 195 个双方重叠路径中的自动合并；优先版本/CI、Backend 数据与协议、Web 模型和状态、画布/预演、测试/文档，不能仅因无冲突标记放行。
- [ ] v1.6.6 全量 Web/Backend Linux CGO/Agent、插件包、静态配置、候选镜像及无端口 HTTP；旧 Schema9 和 A″ Schema45 的隔离克隆迁移与关键数据/密钥保持。
- [ ] 私有完整生产恢复点与旧镜像保护、用户 Edge 隔离预览验收、保留卷本地部署验证、正式提交、非强制推送用户 fork/核对远端 SHA 和 CI。任一门禁失败即停，不做真实付费模型调用、不 force、不删除旧分支或卷。

## v1.6.0 A″ 增量合并门禁（更新：2026-10-04 04:09 CST）

- [x] 用户选择从 v1.5.8 未验收 A′ 制作仅本地 WIP 检查点，再在隔离 A″ 合入最新正式 v1.6.0；保留旧 A、原 A′、生产容器/卷，不作正式提交、推送或晋升。
- [x] 只读确认官方 GitHub `v1.6.0` Release 非草稿/非预发布，远端标签 `42f14cd7`；A′ 原 merge 无未合并路径、无未暂存/未跟踪项。
- [x] 固化 A′ 精确暂存树为仅本地 WIP 双父检查点 `7419e2b3`，fetch 固定 v1.6.0 后建立 A″ `.local/worktrees/upstream-v160`；原 A′/旧 A 保留。
- [x] v1.6.0 实际 40 个文本冲突已全部暂存解决；55 个自动合并重叠路径按模块、调用链和全量回归复核，高风险的重复提示词 Modal 已修复并补回归。95 个路径是总重叠（40 显式+55 自动），不把“无冲突标记”当业务验收。
- [x] Backend Linux CGO 全仓、五项 PostgreSQL/Redis 依赖用例、Web 全量 2,406+5/类型/lint/构建、Canvas Agent 328/构建、两个百炼 v4 包逐字节核对、旧9/本地42/官方43 三谱系隔离 Schema45 迁移、重复 up/verify、完整性/外键与脱敏计数保持通过。
- [x] 四套 Compose 静态配置、独立 Backend/Web 候选镜像及无宿主端口 HTTP 冒烟通过；隔离测试容器/网络已清理。源码测试和数据库副本不能替代登录态 Edge。
- [x] Microsoft Edge 手工验收方案写入 `HANDOFF_V160_20261004.md`；未启动 v1.6.0 宿主端口预览，须另核对/批准隔离预览后交付实际 URL。
- [ ] 用户 Edge 手验、真实付费/设备测试（另授权）、正式合并提交后的 CI 格式门禁、非强制远端推送、生产升级前恢复点与晋升均未完成；不替换生产、不删除旧 A 或 A′。

以下 v1.5.8 待办保留为上一候选历史；不要把它的测试结果自动赋予 v1.6.0。

## 当前 A′ 发布前门禁（2026-09-29 CST）

- [x] 固定官方 `v1.5.8@25842275`；在隔离 A′ 解决 59 个显式冲突，未提交/推送/部署，旧 A 保留。
- [x] Web 全量 2380、类型与源码构建通过；Canvas Agent 328 项、构建与 frozen-lockfile 校验通过；百炼 v4 包源/归档逐字节一致；默认、本地、开发与占位环境下生产 Compose 解析通过。
- [x] 对 98 个自动合并重叠路径复核文件集与高风险调用链；23 路径在合并后作定向修复/格式收口，其余 75 与 Git 自动结果一致。完整 Web/Agent/Linux CGO、迁移、镜像和静态 Git 门禁通过；这不替代画布/导演台的用户 Edge 视觉与交互验收。
- [ ] 正式提交前按最终差异重跑 `web/scripts/check-changed-formatting.mjs`；当前 WIP HEAD 的父提交不是最终 A′ 基线，其报出的历史格式失败不能直接视作候选已通过或真实新增问题。上游 Chrome Director E2E 因 Edge-only 约束未运行，另用隔离 Edge 夹具验证。
- [x] 用户已启动 Docker Desktop；隔离 Linux CGO 后端全量通过，CI 五项依赖测试均实际运行通过；旧 Schema 9 克隆复制及迁移前 integrity/foreign key/计数检查通过。
- [x] 旧 Schema 9 与固定官方 Schema 40 隔离卷均升级到 42，重复 up/verify、完整性/外键和脱敏计数通过；Backend/Web 独立候选镜像及隔离 HTTP ready/版本/三平台 Bridge 下载冒烟通过。生产卷、现有 Compose 与 `:local` 标签未变。
- [x] 最终候选暂存索引无冲突、无空白错误、无未暂存/未跟踪；旧 A clean，构建后新增差异仅项目记忆/专题文档。隔离测试容器和网络已清理，3 个迁移验证卷与候选镜像保留。
- [x] 2026-09-29 用户批准隔离预览：`http://127.0.0.1:3001` 仅绑定回环地址，独立空卷与专用 Docker 网络；Web/Backend healthy，HTTP build-info 与 ready 200、Schema42/42，生产容器镜像/健康/重启计数未变。预览保持运行供用户 Edge 手验，未使用生产账号、数据或密钥。
- [ ] 实际合并提交形成后需重跑真实 CI 格式门禁和远端检查。用户 Edge 手验仍是晋升前 NO-GO。
- [ ] 用户在已登录 Microsoft Edge 访问隔离预览并用新建测试账号手验；不得把生产 `localhost:3000` 的 Cookie、账号或素材当作预览数据。待用户确认通过并明确批准后再正式提交、非强制推送自己的 `origin`，核对远端 SHA/CI。旧 A 仍保留，删除另获精确批准。
- [ ] 真实百炼/即梦/Comfy 生成和其他付费调用另走单次费用授权；Director F6 只在隔离测试环境构造确定性保存失败。

以下为历史计划快照；旧标题/复选框仅代表当时状态，不是 A′ 的待办或完成证明。

## 历史：固定 v1.5.8 增量审计（实施前）

- [x] 只读确认官方 `main=v1.5.8@25842275`，以 v1.5.6 共同基线和当前完整候选树模拟得到 59 个显式冲突；方案见 `UPSTREAM_V158_INCREMENT_AUDIT_20260927.md`。
- [ ] 先固定当前 v1.5.6 候选快照和 Backend/Web 已知失败，处理阻断归因的基线错误；最终全量与脱敏双谱系克隆在 v1.5.8 新候选上验收。`MERGE_HEAD` 打开期间不在原工作树叠加第二次 merge。
- [ ] 获准实施后从隔离候选分批合成 59 个冲突，并审计 98 个自动合并重叠路径；重点门禁为迁移编号、v4 插件协议、生成持久化、画布/导演台和生产 CORS。
- [ ] 最终全量验证、插件归档再生、候选镜像与用户 Edge 门禁通过前，不发布/部署。

## 历史：v1.5.6 合并

- [x] 用户选择方案 A：保留并迁移本机 Agent、Runtime、影策 Codex 插件、本机渠道和人像能力。
- [x] 用户决定旧 `/sessions` 改接上游新 Agent；Comfy Bridge 后端接口保留并迁移（2026-09-24 明确授权）。
- [x] 旧网页 `/sessions` 影视会话入口改接上游 Cloud Agent；删除失效 Web 面板/轮询 helper，独立 Canvas Agent/Codex 插件仍保留。旧数据归档或只读展示仍需单独设计，不能假装自动迁入新 Agent。
- [x] Comfy Bridge handler、状态机、仓储、模型和 Schema31 迁入上游拆包；Linux CGO 专项覆盖 Token、跨用户撤销及重启恢复。真实设备生成、取消费用和克隆数据仍待验证。
- [x] 全部显式 Git 冲突清零；无提交、无推送、无部署。此项不是语义合并或发布门禁通过。
- [ ] 双谱系 Schema 9/29→31 在隔离 Linux CGO/克隆数据上验证两类来源、二次启动、历史业务摘要保持。
- [ ] 复核自动合并重叠路径并修复 Linux CGO 全量失败、Web 类型/测试/构建、插件包及部署配置；保持无提交，不触碰运行数据。

## Gemini 系统模型目录阻断（实现部署完成，待 Edge）

- [x] 用户批准统一方案；ad3efa10完成资格/插件合同/分页/草稿/检查配置，cc7d0b63补普通文本媒体防丢失。
- [x] 最终Backend全量、Web1645全量/40专项/类型/构建、70文件哈希与7运行包核对通过；两次断网克隆双启动、私有备份及部署数据保持通过。
- [ ] 用户按MODEL_CATALOG_UNIFIED_REPAIR_20260913.md完成新版本Edge配置验收；真实调用/账号权限/额度/计费仍未放行。

完整方案：MODEL_CATALOG_UNIFIED_REPAIR_PLAN_20260913.md。公共模型资格评估、上游插件合同扩展、发现分页与现有百炼合同统一一起设计，不仅补Gemini白名单；68条planned百炼补充模型不因目录修复被视为完成执行支持。

- [x] 2026-09-13 定位通用 planned + 缺 Gemini 目录补全 + 前后端只读闭环；见 GEMINI_CATALOG_DIAGNOSIS_20260913.md。
- [x] 补齐官方目录字段/分页、统一执行合同映射、已有候选受校验配置流程及跨厂商回归；保留模型合同与价格/启用门禁，不直接全量改ready。

## 历史：画布单请求多图

- [x] 补齐两种画布编辑入口、一次提交多结果、输出序号与物化、刷新恢复、单图重试幂等和报价。
- [x] 84 专项、1645 全量、类型/构建及 9 文件哈希验证；仅前端 cd0b4340 部署，后端与数据保持。
- [x] 2026-09-13 用户确认截图所列 Edge 无费用配置检查暂无问题，并自行配置启用 qwen-image-3.0-pro。
- [ ] 真实多图生成、生成中刷新恢复、单图重试和实际扣费需另行费用授权；其他六个视频模型不因本次反馈视为已验收。

## 历史：七模型能力补齐

- [x] 用户批准使用阻断、高级参数、媒体引用和可选实际时长计费；保留既有价格与订单。
- [x] 完成全部实现与系统性回归：Backend 全量、Web 1633 全量/类型/构建、58 个修改源码文件与测试镜像逐字节一致。
- [x] 完成 Schema 9 克隆双启动、保留卷部署及数据/环境/运行身份复核；新恢复点均 verified/protected。
- [x] 已准备完整交付清单，停止自动执行，不进入浏览器或真实生成。
- [ ] 用户使用已登录 Microsoft Edge 完成人工确认。

## 历史：百炼完整插件迁移

- [x] 在正确 v1.2.7 基线建立独立修复分支，移除过渡方案，接入正式声明式插件和通用引擎扩展。
- [x] 最终提交 8dcdfd5d 的 Backend Linux CGO 全量、插件源码与包一致性、私有备份、离线克隆双启动、保留卷部署前后数据验证全部通过。
- [ ] 用户在已登录 Microsoft Edge 进行最终协议选择器与插件中心人工验收；真实模型生成不在本次预验证授权内。

## 历史：阶段9完成，F6转为测试环境门禁

- [x] 用户已在Microsoft Edge完成当前生产环境可安全观察项复验，暂未发现问题。
- [x] 候选与稳定分支经dry-run后非强制快进到37bb2375；`origin/main`、默认分支设置、PR、Release和Host Updater未改动。
- [ ] 在隔离测试环境用确定性保存失败夹具完成Director F6；生产环境不制造网络、凭据或数据故障。
- [ ] 只读审计官方`371bc68..83d9622`尾差；通过新阶段批准后再决定增量合并。
- [ ] 主线稳定后另行评估Astravia次级工作台，不接入生产数据库、密钥或稳定checkout。

## 历史：Edge首轮失败项待复验

- [x] 普通页面取消自动Runtime探测；已知无价格规格不再发送报价；暗色点网格增强；宽/窄画布增加导演台入口。
- [x] 专项34、Web完整1624、类型、格式和生产构建通过；Web修复镜像已部署，Backend/数据不变。
- [x] 用户Edge清空DevTools后复验首页/创作/项目/画布/分镜；当前可安全观察部分暂未发现问题。
- [x] 项目/自由画布暗色点网格及导演台入口等可观察项已复验；F6另转隔离测试环境。

## 历史：阶段8复盘整改完成

- [x] 修复模型勾选导入与存量安全补齐断链；新增、补齐、保护配置与二次零更新已有测试证据。
- [x] 生产/测试镜像补ffmpeg/ffprobe；CI切换完整测试发现；配置一致克隆双启动和真实卷部署摘要通过。
- [x] D盘敏感审计副本迁移到私有ACL目录；新部署前恢复点verified/protected。
- [x] 用户使用登录态Microsoft Edge完成生产环境可安全观察的视觉/Runtime/报价/加载稳定性验收；F6单列测试环境。
- [x] 同步远程候选并完成阶段9稳定晋升；远程候选与stable首次对齐功能检查点37bb2375。

## 历史：阶段7已完成

- [x] 非强制推送fork候选并独立核对SHA；stable/main未变。远程CI按分支策略未触发，不记为通过。
- [ ] 阶段8需新批准：新恢复点、克隆数据Schema6→8双启动、候选部署，并停在用户Edge验收；F6与全部视觉项统一由用户手工验证。
- [ ] 若需远程CI，单独决定PR或安全workflow dispatch；不得从本次push授权推导。

## 历史：阶段6保留人工门禁

- [x] 自动化源码验证与构建身份固化：源码88087501，Web1624通过、Backend全量、Agent328跨平台覆盖、Compose/脚本/Manifest检查通过。
- [x] 用户决定Director F6与视觉项统一留到阶段8手工Edge验收；它不阻挡阶段7候选推送，但仍是阶段8必测项。
- [ ] 推送/部署/真实迁移/Provider调用仍未授权。默认官方镜像/Host Updater不得覆盖fork能力。

## 历史：阶段5源码已完成，等待阶段6

- [x] 三个UI提交选择性适配为8082e9e2/2aebd53e/d10c5e32；默认1278测试/干净类型/生产构建通过。
- [ ] 阶段6最终组合源码/吸收清单/安全/Agent/Compose等门禁，经批准后执行。
- [ ] 之后候选push、数据副本迁移、部署、用户Edge和稳定晋升继续分批批准；现有视觉失败及偶发加载失败不能提前关闭。

## 历史：v1.2.7阶段3–4检查点

- [x] 固定371bc68合入stable5c1203e的独立候选；63路径/184区块处理，Backend全量、Web类型/1255测试/构建通过。
- [ ] 阶段5获批后选择性重放e027aa4/f5fadb7/bcb6c6d；不整包恢复失败UI2布局。
- [ ] 阶段6最终验证、阶段7候选push、阶段8副本迁移/部署与用户Edge、阶段9稳定晋升按顺序另行批准。

以下旧版本清单为历史，当前状态以 `UPSTREAM_V127_STAGE34_20260908.md` 为准。

更新时间：2026-09-04

## P0：接管与数据安全

- [x] 固定 `v1.2.5@f8e87bc` 并完成阶段 0–2：历史检查点与新 `origin/codex/dingxin-stable` 均为 `44e8c10`；完整对象审计为 66 提交/506 文件、89 重叠路径、模拟 47 显式冲突文件/149 区块。官方 `main@7ec5173` 的 6 个后续提交保持在外。
- [x] 在 `codex/upstream-v1.2.5-20260904` 完成阶段 3–4：真实 49 未合并路径/139 文本区块已按五组解决，未解决 0；当前 `HEAD=39c7e62`、`MERGE_HEAD=f8e87bc`、状态 `merge_open_resolved_unvalidated`。
- [ ] 获得新批准后进入阶段 5；先审计最终 staged bytes，再执行格式、Backend Linux CGO、迁移副本双启动、Web/Agent/Docs/Compose/Windows 脚本门禁。阶段 5 通过前不得创建 merge commit。
- [ ] 阶段 5 第一项必须移除 Create 对已删除 `backendModelRuntimeRequired` 的依赖，按官方 v1.2.5 统一走 `runBackendGenerationTask`，保留本地元数据/reasoning/任务恢复合同并补专项复核；修复前 TypeScript 预期无法通过。
- [ ] 在阶段 8/9 前为 Host Updater增加用户 fork Release/镜像源与版本策略，或显式保持不可用；当前默认 `ddcat-ai/open-ai-canvas` 不能用于自维护版本在线更新。
- [ ] v1.2.5 完成后再实施本地扩展边界收口，并在阶段 0–9 全部完成后单独评估/安装 Astravia Windows；当前不安装、不向其开放生产数据库、密钥或稳定 checkout。
- [x] 固定 `4ba9694` 已固化为本地 merge commit `d0def5f`，阶段 C通过；提交后评审记忆提交 `7687053` 已fast-forward推送fork并独立复核一致。本次最终检查点记忆继续以小型docs提交同步，完成后该集成分支作为已验证固定4ba远程检查点。没有force、未更新`origin/main`、未触发远程CI或部署。
- [ ] 后续独立事项：决定是否为fork集成分支建立可手工触发或PR触发的远程质量门禁；审计官方 `4ba9694..b7348ab` 22提交/216文件新尾差；另行规划阶段16部署到新本地提交的备份、迁移、镜像和Edge验收。
- [x] 路线B已完成升级前恢复点、回滚镜像标签、候选构建、克隆数据迁移/二次启动幂等和真实保留卷升级；非浏览器健康、日志、数据库和版本检查通过。
- [x] 2026-09-04 用户使用已登录Microsoft Edge完成人工只读冒烟并回传11张截图；版本/登录、稳定首页、导航、项目/章节、空资产分镜、素材/任务/Skill、后台及Console/Network基础健康通过，状态为 `runtime_candidate_edge_smoke_passed_limited`。
- [ ] 另行授权并准备真实项目资产后，补测分镜资产绑定/解绑、`@` mention、历史角色图/当前声音、音频引用、长名称溢出、刷新恢复和多镜头任务恢复；这不是本次空资产冒烟已覆盖范围。
- [ ] Provider真实生成、精确计费、OSS新写入、模型/价格保存、存储连接测试、Skill同步和删除操作继续保持独立门禁；不得因Edge只读冒烟通过而自动执行。

- [x] 对当前 SQLite 主库和 WAL 做一致性备份，并验证完整性、独立恢复副本和 OSS 密文解密。
- [x] 记录当前镜像 ID、容器、卷、数据目录和关键文件哈希。
- [ ] 收敛 `.claude/settings.local.json` 的宽泛授权。
- [x] 使用 Git binary patch 和 44 文件 ZIP 保护当前 9 个已修改文件和 35 个未跟踪文件，并验证 CRC、路径集合和 SHA-256。
- [x] 建立 `BACKUP_REGISTRY.md`，登记阶段 0 的非敏感基线、私有数据库快照、密钥、哈希、ACL、验证和保留策略。
- [x] 在主机执行环境复核 `gh auth status`/`gh api user` 并查询明确仓库；确认受限进程只是无法读取 Keyring，当前没有对应 PR。远程命名整理仍需在未来 Git 操作阶段单独评估。

## P1：Qwen 与能力合同

- [x] 建立唯一 `effectiveChannelModelCapability` 持久边界，由系统渠道 Catalog 与 Admission 共用；精确 SKU 计费和 Ready Provider 使用同一持久模型合同。
- [x] 明确定义并测试 missing、explicit-empty、wildcard 三态；显式空数组不回填，`*` 只开启自定义输入。
- [x] 按 Qwen、Wan、HappyHorse 具体模型建立保守能力矩阵：Qwen Image 3.0 Pro、Wan 2.7、HappyHorse 1.1、Wan 3.0 已完成；其他图片模型继续 Planned。
- [ ] 移除执行器不支持或会被静默改写的图片/视频参数：Ready 视频模型已完成；Planned/手工图片与通用 fallback 尚未全量收敛。
- [ ] 统一 Create、Canvas、Admin 和逻辑模型的尺寸派生纯函数。
- [x] Create 参数按用户、模型和生成方式持久化；刷新保持画幅/分辨率/时长/声音/水印，明确 Switch 与请求 payload 一致，并完成 Edge 验证。
- [x] 在管理后台和运行目录验证 Qwen Pro 的 auto、比例/像素、自定义尺寸、3 参考图、6 输出及不支持项；普通 Create 真实生成设置仍需模型定价/启用后的单独授权验证。
- [x] 删除 Qwen 临时 DEBUG、完整提示词/请求体/签名 URL/API Key 片段日志和未使用导入；Provider 只保留结构化安全状态。
- [x] 为官方发现目录增加 ready/planned/unsupported/deprecated 支持状态，并由 Backend 强制限制非 Ready 状态不可定价、启用、测试、删除或路由。
- [x] 重构 `FetchAdminChannelModels`，保留完整 CatalogItem 元数据并为现有待配置记录提供安全补齐；阶段 7 已在运行数据库完成拉取并验证真实幂等。
- [x] 完成 Wan 2.7、HappyHorse 1.1、Wan 3.0 的模型专属能力合同并晋升对应 Ready 模型；Prime、旧版和视频编辑继续由精确合同与账号证据控制。
- [x] 为 Wan 3.0 实现独立 All-in-One Adapter，支持 T2V、首帧、首尾帧、参考图片/视频/音频；file/link 保持未开放。

## P1：测试与可观察性

- [x] 复核 Canvas Agent Windows 进程树清理：确认 Codex 受限进程的 `taskkill` 被拒绝导致假阴性；测试夹具显式探测能力，主机真实专项 12/12、全量 287/287 和 TypeScript 构建通过。远程 Linux CI 仍随未来 PR 验证。
- [x] 用 Qwen 3.0 专属官方尺寸合同替代旧 38/39 项通用列表，并覆盖存量 Planned→Ready 安全补齐、字段缺失/空数组/`*` 与非法默认值边界。
- [x] 增加 Qwen `auto` 省略、`x -> *`、比例映射、1–6 输出、1–3 图编辑、参数依赖和拒绝路径完整序列化测试。
- [x] 为 Ready 视频模型增加 480P/时长/音频/水印、跨媒体数量、输入相关时长和非法组合拒绝测试。
- [x] 为系统渠道持久能力边界、目录/准入 fail-closed、Ready Provider payload、精确 SKU 账务和真实幂等增加测试。
- [x] 为任务失败建立前端、handler、service、worker、provider、资源物化全链路结构化日志，并以 request/task/provider/resource ID 关联。
- [x] 修复错误响应不记录日志或可能暴露内部错误的问题；统一稳定错误分类、retryable、requestId 和前端诊断编号。
- [x] 收敛生产 GORM 日志：忽略正常 `record not found`，保留真实错误元数据并整体隐藏 SQL 文本，覆盖 Raw SQL 预插值边界。

## P2：插件、部署与 CI

- [x] 修复插件备用端口、Vite 启动命令和可信 Origin 配置。
- [x] 让插件超时与 Agent 35 分钟生成续接合同一致，并保留 60 秒余量。
- [x] 用签名 Local Runtime 握手替代 `agentToken` 查询参数旧协议；Web 继续主动清除旧参数。
- [x] 修复不同 Compose 对 `ENABLE_PROVIDER_PLUGINS` 的透传差异，五个 Backend 入口默认均为 false。
- [ ] 在可调用 Codex CLI 的环境重新安装 `yingce@yingce-local`，新建线程验证插件 skill/MCP 加载和 Microsoft Edge 自动打开；当前 WindowsApps ACL 阻止 `codex plugin list`，不得伪报安装态通过。
- [x] 收紧 Server CORS：生产 Server Compose 必须显式提供非空 Origin，不再默认 `*`。
- [x] 在 CI 中增加 Canvas Agent 测试/构建，并为 Web 增加生产构建；远程 Actions 仍需提交/PR 后验证。
- [x] 消除 Bun test 内嵌 Vite/Rolldown 的退出竞争，增加 panic guard，并用 `.bun-version` 统一本机/CI；远程 Actions 仍需提交/PR 后验证。
- [ ] 补齐文档站缺页、路径和工具表。

## P2：百炼技能与协议资产

- [x] 逐篇索引本地百炼官方文档并保留 SHA-256、行数、标题和分类。
- [x] 建立全局 `bailian-model-api` 技能，按文本、图片、视频、语音、向量和平台 API 分流。
- [x] 建立逐模型同步/异步、端点、输入、限制、轮询、输出和临时 URL 的检索与证据合同。
- [x] 完成技能结构验证、项目/全局副本哈希核对和典型文本、图片、视频、语音、向量检索测试。
- [ ] 官方文档更新时重新生成索引并审阅差异。
- [ ] 跟进官方 Issue [#332](https://github.com/ddcat-ai/open-ai-canvas/issues/332) 的维护者回复；2026-08-28 实时复核仍为 open/评论 0/无标签/无指派，五项架构决策确认前，不提交支持状态 Schema、百炼 Adapter 或版本化 Manifest PR。
- [ ] Issue 提交后先等待 2–3 个工作日；若仍无维护者回复，再由用户明确授权后发送一次简洁的五项决策确认评论。禁止自动催办、重复评论或把点赞、标签、指派、交叉引用、关闭 Issue、普通贡献者回复和笼统“欢迎 PR”当成完整批准。
- [ ] 收到回复后核对回复者的 `author_association`；只有 `OWNER`、`MEMBER` 或 `COLLABORATOR` 的明确回答才作为维护者决策。逐项记录为 `confirmed/partial/alternative_requested/rejected/superseded`，部分回答继续保持未决项 NO-GO。
- [ ] Issue #332 获认可后，按“支持状态 → Qwen Image 3.0 Pro → Wan 3.0 → Wan 2.7 → HappyHorse 1.1 → 版本化 Manifest”顺序从官方最新 main 建立独立贡献分支。

## 执行门禁

任何实施阶段都按以下顺序：

1. 最小专项单元测试；
2. Backend 全量测试；
3. Web 测试和构建；
4. Canvas Agent 测试和构建；
5. 数据备份；
6. 保留卷重建并重新创建容器；
7. 源码、镜像、容器和浏览器哈希核对；
8. Mock/无费用浏览器端到端验证；
9. 用户另行批准后才进行真实付费模型调用。

## 阶段 16 本地交付状态

- [x] 最终 Backend、Web、Canvas Agent、插件和 Compose 全量门禁。
- [x] 最终数据完整性、OSS 加密、容器、镜像、数据卷和日志核对。
- [x] Microsoft Edge 首页、管理后台和创作页人工只读冒烟。
- [x] 创建 `BKP-20260827-160505-STAGE16-FINAL` 并登记最终本地镜像标签。
- [x] 生成 `STAGE16_FINAL_HANDOFF.md` 并清理项目记忆中的已修复旧结论。
- [ ] 集成分支 push 已完成；PR、远程 Actions、更新 `origin/main`、Git tag 和 Release 仍等待用户另行授权。

## 阶段 1 后续门禁

- [x] 在隔离 Linux CGO Backend target 中运行 SQLite 迁移和 Backend 全量测试；远程 GitHub Actions 仍待后续提交/PR 后确认。
- [x] 阶段 2 完成 Wan 3.0 All-in-One Adapter 与序列化/拒绝测试，标准版已从 `planned` 晋升为 `ready`；Prime 仍以精确官方合同和账号可用性证据为门禁。
- [x] 阶段 3 完成 Wan 2.7/HappyHorse 1.1 模型专属合同并晋升对应目录项。
- [x] 阶段 4 完成 Backend CGO、Web 和 Canvas Agent 全量测试与构建；阶段 6 以后才允许备份后的数据库迁移、容器重建和 Edge 运行验证。
- [x] 阶段 5 完成发布候选审计、能力/价格 fail-closed 收敛、敏感日志清理、独立候选镜像和迁移/回滚方案。
- [x] 阶段 6 创建 pre/post-deploy verified/protected 备份，保留卷逐个更新 Backend/Web，完成 Schema、数据计数、OSS 密文和运行哈希核对。
- [x] 阶段 7 已使用登录 Microsoft Edge 验证模型拉取、80 项补充目录、11 Ready/308 Planned、管理后台支持状态、重复拉取真实幂等，以及普通用户创作台不暴露未定价且停用模型；未点击生成。
- [x] 阶段 8 已完成一次授权内 Wan 3.0 真实调用、精确 SKU 账务、阿里云 OSS 物化、Resource/Asset、Edge 播放、无声/无水印、模型安全恢复和调用后 verified/protected 备份。实际画幅为 adaptive，未补发第二次任务。

## v1.2.5 阶段 5 后续

- [x] Create Backend-only 主链、重叠 Web 语义修复、官方公告图片自引用缺陷修复。
- [x] Web 类型/默认全套/补充专项/生产构建、Backend Linux CGO 全量、Compose 与 PowerShell 静态门禁。
- [x] 独立候选镜像、verified/protected 快照、两次无网络克隆迁移和计数/密钥验证。
- [x] Canvas Agent 以 Node 22 + Bun 1.4 Linux 隔离环境完成 328 项测试和构建；修正 unref heartbeat 测试的有界等待，0 fail、0 cancelled。
- [x] 阶段 5 最终 staged-bytes、冲突标记、敏感路径和版本复核通过；阶段 6 merge commit 仍需用户批准。
- [ ] Host Updater 完成 fork 仓库、镜像命名空间、版本策略和祖先校验前保持不可用。
- [ ] 最终 Microsoft Edge 验收对比本地画布底板/柔和网格与官方 v1.2.5 底板/80% 网格；由用户选择最终默认，当前实现不得被表述为永久产品决定。
- [x] 阶段 6 创建本地双父 merge commit `b6cff79`，固定纳入官方 `v1.2.5@f8e87bc`。
- [x] 阶段 7 已将候选分支非 force 推送到 fork 并独立核对 `dced038...`；`origin/main` 未变化，未触发可见 Actions，不能称远程 CI 通过。
- [x] 阶段 8 将 `codex/dingxin-stable` 非 force fast-forward 到 v1.2.5 候选；候选/稳定一致，`origin/main` 未变化。
- [x] 阶段 9 创建 verified/protected 恢复点、完成新克隆卷双启动、旧镜像保护、真实保留卷 recreate、HTTP/版本/schema/数据库验证。
- [ ] 用户使用登录态 Microsoft Edge 完成阶段 9 最终验收；重点对比当前本地画布底色/柔和网格与官方 v1.2.5 底色/80%网格，再决定默认方案。
