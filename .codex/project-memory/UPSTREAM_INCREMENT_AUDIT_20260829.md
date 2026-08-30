# 官方主线新增量审计：`ab89c05..4f07daa`

## 状态与授权边界

- 审计日期：2026-08-29。
- 当前本地分支：`codex/upstream-20260828-ab89c05`；审计起点 `HEAD=c15a3aeca1e47821300bb9388bc62ccdd00afc12`。
- 当前分支已经通过 merge commit `d04c4d1d968ad2c53de919985596d61c8d841f5f` 包含官方 `ab89c05362394623d00c95a4899be03e5243e75a`。
- 本轮两次只读确认官方 `main`，最终固定候选为 `4f07daae9ec3b4e4cb0a8cd35a6e0c1a4b593f29`；审计结束前再次确认未漂移。
- 用户本轮批准的是新增量审计，不包括真实 merge、冲突编辑、merge commit、push、PR、远程 CI、部署、数据库/卷操作或 Provider 调用。

## 增量规模

- 共同基点：`ab89c05362394623d00c95a4899be03e5243e75a`。
- 官方新增：14 个提交、107 个文件、6,843 行新增、661 行删除。
- 当前本地相对共同基点改动 195 个文件；官方新增量与本地重叠 22 个文件。
- `git rev-list --left-right --count HEAD...official/main`：本地独有 42、官方独有 14。
- `git merge-tree --write-tree --messages --name-only HEAD official/main` 只读模拟得到 6 个显式文本冲突；工作树未进入合并状态。

## 新增功能清单

1. `115e228`：README 新增赞助商展示和两张图片。
2. `10c57ac`：视频首帧/当前帧/尾帧与最多 30 个时间点的批量取帧；视频片段截取与重新生成解耦。
3. `8535658`：复制、参数变体和批量生成结果保留唯一命名，避免下载重名。
4. `a8b5f11`：视频画幅能力可为空；按画幅和分辨率推导只读宽高，不支持画幅时隐藏尺寸设置。
5. `d66c7ad`：短剧分镜工作流 v2、镜头版本/产物/生产任务上下文、项目工作台，以及系统渠道默认价格档。
6. `d84b2b6`：已有图片节点显式引用进入图生图；图片加音频按图生视频匹配，纯音频才按音频生视频。
7. `bf37e45`、`0f274bb`：公告置顶、管理员标识、配图草稿、自动回收和引用安全。
8. `b0d77db`、`ae0fdda`：管理后台资源查询/统计/预览，以及引用校验、共享对象保护和事务化批量删除。
9. `f58481c`：忽略 Playwright CLI 本地文件。
10. `80d2aa6`：更新两张赞助商/品牌图片。
11. `aa804d4`：本地公开资源直链追加扩展名；AutoDL 工作流 Manifest 参数对齐。
12. `4f07daa`：公告 Markdown 预览后发布。

## 六个显式冲突及解决合同

| 文件                                          | 风险                                                                                             | 解决合同                                                                                                                                                                                                                                                                                                                                            |
| --------------------------------------------- | ------------------------------------------------------------------------------------------------ | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `backend/internal/service/finance.go`         | 高。错误选择价格档会导致错误预扣、错误免费、线路切换失败或账单与请求规格不一致。                 | 保留本地“服务端按完整 intent 重新选择价格档、缺失即 fail-closed”的行为；采用官方 intent-aware `newBillingOrderWithPriceTier` 签名。直接系统模型使用服务端选出的 tier ID 并同时传入 intent；逻辑模型和线路切换使用已选路由的 `PriceTier.ID`，不信任客户端或旧任务中的 `config.priceTierId`。补默认档、精确档、零价显式配置、线路切换和过期 ID 测试。 |
| `docs/content/docs/progress/pending-test.mdx` | 低运行风险、高审计风险。任一侧覆盖都会丢失本地阶段历史或官方新增手工验收项。                     | 保留本地历史/NO-GO 前言，再追加官方短剧、取帧、命名、默认价格、图生图、存储管理等章节；按二级标题去重，不把自动化通过升级成浏览器/生产验证。                                                                                                                                                                                                        |
| `web/package.json`                            | 中。错误解决会丢测试、绕过 panic guard，形成假绿。                                               | 保留本地 `test` → panic guard → `test:main`/`test:cross-runtime` 拆分；把官方新增 `canvas-node-generation-mentions.test.ts`、`canvas-video-frame.test.ts` 纳入 `test:main`，并保留 `creation-settings-store`、`qwen-image-capability` 和跨 Runtime 独立门禁。官方没有 `bun.lock` 变化。                                                             |
| `web/src/components/video-settings-panel.tsx` | 高。可能重新显示不支持的尺寸、丢失智能时长，或让价格档与实际请求不一致。                         | 采用官方 `ratios.length > 0` 才显示尺寸、画幅+分辨率推导只读宽高；保留本地 `-1` 智能时长、智能时长仅 fixed-request 可定价、声音/水印与模型能力驱动。删除任意宽高编辑入口，避免向只声明分辨率的工作流发送旧全局 size。                                                                                                                               |
| `web/src/lib/model-capabilities.ts`           | 高。能力漂移会放大 Ready 模型参数、破坏显式空数组或丢失插件工作流能力。                          | 取并集：保留 Qwen/Wan/HappyHorse 精确能力、视觉总数、智能时长、显式空数组和 `*` 三态；接入官方 `pluginWorkflowCapabilityConfig`、`resolveVideoRatioValue` 和“画幅可为空”。`normalizeVideoValue` 同时保留智能时长并使用 `resolveVideoRatioValue`。                                                                                                   |
| `web/src/pages/create/index.tsx`              | 高。可能丢失按用户/模型/操作的设置持久化、声音/水印授权摘要，或向不支持画幅的模型发送陈旧 size。 | 保留本地设置持久化与声音/水印 Switch/摘要；采用官方画幅可为空和 nullish 语义。视频请求只有 `videoProfile.ratios.length > 0` 时才写 `size`，避免空字符串或旧全局比例进入 payload；分辨率仍按能力归一化。                                                                                                                                             |

## 自动合并但必须语义复核的区域

1. **系统价格与支持状态**：`channel-model-manager.tsx` 自动合并后同时保留本地 `SupportStatus` 过滤和官方默认/规格价格表单。必须验证 `PriceConfigured=true` 才允许零价、关闭“可供用户使用”后不参与匹配，以及 Planned/Unsupported 模型仍不能定价或启用。
2. **Backend 视频能力**：官方允许 `Ratios=[]`，本地精确 DashScope 模型仍保留各自比例合同。必须验证只声明分辨率的插件工作流不会被默认画幅补全，Ready 模型也不会意外失去比例。
3. **图片+音频路由**：官方把图片/角色+音频归入 `image_to_video`，再由参考音频数量约束筛选模型。必须保留本地 `MatchCapability` 的音频上限和 Wan/HappyHorse 精确能力，防止路由到不支持音频的 I2V 模型。
4. **公开资源直链**：自动合并同时保留本地脱敏日志和官方扩展名路径。必须验证签名仍只绑定资源 ID/过期时间、旧无文件名路由兼容、Range/Content-Type、HTTPS/私网 S3 门禁和不转发 Provider 凭据。
5. **AutoDL 插件包**：官方同时修改 Manifest JSON、二进制 `.yingce-plugin` 和协议制品测试。必须验证包内外 manifest 一致、哈希/版本和宿主兼容 shim，不能只看 JSON。
6. **短剧工作流 v2**：新增 `shot_revisions`、`shot_artifacts`、`production_task_links`，并扩展 `shots`。必须在隔离数据库验证 AutoMigrate、SQLite→PostgreSQL 迁移、旧 v1 实例保留、用户/项目/画布/任务归属和局部 stale 传播。
7. **公告与资源管理**：新增 `announcement_image_drafts`、后台跨用户资源读取和批量删除。虽然 service 有管理员校验、引用快照、共享对象计数、事务和 outbox，仍需验证公告/素材/画布/项目/工作流引用取并集，物理删除失败不删除业务记录或丢任务。
8. **Markdown 公告**：复核预览与最终渲染使用同一安全 Markdown 边界，禁止原始 HTML/XSS、危险 URL 和管理员预览/用户显示不一致。

## 推荐实际合并阶段

### 阶段 12A：固定 SHA 并打开本地 no-commit 合并

- 开始前再次 `ls-remote`；即使官方继续前进，本轮只合并已审计的固定 SHA `4f07daa`，新尾差另行审计。
- 确认工作树干净、`VERSION=v1.1.4`、恢复点仍可用。
- 执行 `git merge --no-commit --no-ff 4f07daa`；只核对冲突是否仍为上述 6 文件，若集合变化立即停止。

### 阶段 12B：分层解决冲突

1. 先处理 `web/package.json` 和待测文档，锁定测试入口与审计清单。
2. 处理 `finance.go`，先锁定服务端价格选择合同，再运行 Backend 价格/路由专项测试。
3. 处理 `model-capabilities.ts` 与 `video-settings-panel.tsx`，锁定能力和画幅/智能时长边界。
4. 最后处理 `create/index.tsx`，接入持久化、输出开关和“无画幅不发 size”。
5. 逐个复核上述 8 类自动合并语义区域；禁止批量 ours/theirs。

### 阶段 12C：验证门禁

- Backend：gofmt；价格档/线路切换、模型能力、资源直链、AutoDL 制品、短剧 v2、公告草稿、后台资源安全删除专项；随后隔离 Linux CGO `go test ./...`。
- Web：Prettier、TypeScript；新增取帧/命名/模型匹配/价格表单/工作区状态测试；本地 Qwen/Wan/HappyHorse、设置持久化、声音/水印和跨 Runtime 测试；最后 panic guard 全量与生产构建。
- Canvas Agent/插件：插件包合同和 Canvas Agent 全量/构建，确认 AutoDL 包与宿主兼容。
- 文档/配置：MDX 格式和章节去重、`git diff --check`、冲突标记扫描；Compose 本轮无文件变化，可只做关键配置解析而不部署。
- 数据：只在隔离 SQLite/PostgreSQL 测试库验证四张新表与旧工作流迁移；不操作运行卷。
- 浏览器和运行环境：源码门禁完成后另行申请部署/Edge 验证授权；本阶段不能把测试通过写成运行 UI 已验证。

### 阶段 12D：提交与远程门禁

- 全部门禁通过后停在本地 merge commit 决策；merge commit 需要用户再次批准。
- push 到 `origin/codex/upstream-20260828-ab89c05`、PR、远程 CI、更新 `origin/main`、tag/Release 均为后续独立授权。

## 当前结论

- 新增量可合并，但不是低风险三文件同步；建议按 6 个显式冲突、8 类语义复核和完整跨栈验证执行。
- 当前框架只明确跟到 `ab89c05`，尚未 follow `4f07daa`；不能宣称已跟上官方实时主干。
- 官方新增量没有修改根 `VERSION`，本地仍应保持用户确认的 `v1.1.4`，除非另行做版本发布决策。

## 阶段 12A 实际结果

- 用户批准后，开始前再次确认官方实时 `main` 仍为 `4f07daa`；恢复点六项关键文件 SHA-256 与登记值一致。
- 已执行固定 SHA 的 `git merge --no-commit --no-ff 4f07daa`。真实冲突恰好为上文 6 个文件，共 8 个冲突区块；模拟结论无漂移。
- 当前处于未提交合并状态，尚未解决任何冲突。阶段 12B、merge commit 和 push 分别需要后续批准。

## 阶段 12B 实际结果

- 2026-08-30，用户批准阶段 12B。六个冲突文件的八个区块均按本文件合同手工取并集，未解决索引项归零，没有创建 merge commit。
- 计费增加 intent 与显式 tier ID 一致性校验；Create 对无画幅模型清空旧 `size`；测试入口保留 panic guard 并加入官方新增取帧和引用测试。
- Web 相关 9 文件 69/69、TypeScript、Backend 纯逻辑计费/能力和 AutoDL 制品测试通过。宿主 SQLite 测试仅因既知 CGO stub 无法执行，数据库与完整跨栈验证保留到阶段 12C。
- 八类自动合并语义区域已完成源码级复核，但运行数据库迁移、完整 Backend/Web/Canvas Agent、生产构建和浏览器验证尚未执行，不能升级为完整验证通过。
- 阶段结束复核发现官方实时 `main=2f6832f`，相对固定合并目标 `4f07daa` 新增 1 个巨型提交、113 文件、7,016 行新增和 723 行删除。该尾差未 fetch/merge，不影响当前冲突解决结论，但意味着当前分支即使完成阶段 12C 也只能声明跟到固定 `4f07daa`。

## 阶段 12C 实际结果

- 2026-08-30，用户批准阶段 12C。隔离 Linux CGO Backend 全部包、Web 1082/1082 + 跨 Runtime 1/1、TypeScript、13,493 模块生产构建、Canvas Agent 主机 327/327 与构建全部通过。
- Prettier 首次发现 36 个合并后格式漂移并机械修正；修正后 Prettier、TypeScript、相关 Web 69/69、AutoDL 制品、gofmt、Git diff/冲突标记、文档标题、Compose 和 13,493 模块生产构建复验通过；最终构建耗时 7 分 32 秒，最终暂存字节已有构建证据。
- 两张赞助商 PNG 均可解析和查看：`fruivision.png` 为 793×801，`xmzm.png` 为 5270×3536；README 均以 160 像素宽度展示。未进行浏览器渲染验收。
- 结束时官方 `main` 已为 `0893741`；`4f07daa..0893741` 是 2 个未审计提交、122 文件、8,431 行新增和 863 行删除。本轮固定合并和全部验证仍只覆盖到 `4f07daa`。
- 当前停在本地 merge commit 决策门禁；提交、push、部署和新尾差合并均未执行。

## 阶段 12D 实际结果

- 2026-08-30，用户批准阶段 12D；创建本地 merge commit `e2326856b1cc3a854ea11fcd54c5c3e2866c8ac3`。
- 父提交为本地验证基线 `ae58ecf` 与官方固定目标 `4f07daa`，提交后合并状态结束、工作树干净、`VERSION=v1.1.4`。
- 官方实时 `main` 仍为 `0893741`，新尾差未 fetch 或合并。本阶段没有 push、PR、远程 CI、部署或运行数据操作。
