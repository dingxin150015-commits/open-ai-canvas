# 固定 v1.6.6 增量合并、验证与部署执行记录

## 2026-10-15 23:51 CST：远端两正式分支到位，临时传输ref清理

- 用户fork候选与稳定分支已先后非强制到`b740ee0489a2570a52f2e5b7fbdb12c42b76261e`，`origin/main`仍11931d00，本地`codex/dingxin-stable`原子快进同SHA；旧A`codex/dashscope-system-plugins-20260911@6b16d25e`保留。查询GitHub提交b740的PR工作流运行数组为0；质量workflow当前只配置`pull_request`及`workflow_call`，本次分支push没触发CI，不能记通过。
- 技术路径：先候选d9→d6→b9→6b小包；官方v1.5.6在单独临时分支按第一父153提交分段到473d，第152原提交新增58.6MB支付归档/多平台二进制时HTTP408，以忽略区备用Git索引和四个临时合成提交精确分组传原blob，合成树与原官方提交树SHA相同后回推未经改写的原提交。官方随后v1.5.8的110、v1.6.0的47、v1.6.6的100个第一父提交依次分段到固定SHA；候选再56e9→7419→4323→2a→00351→b740，稳定d9→b740。全程未force、未修改系统/Clash全局代理或Git历史。
- 正式候选/稳定均能到官方v1.6.6头27c8066e及所有大blob后，精确删除本任务两条远端传输临时分支，及本地合成支付临时分支；原官方提交对象和正式代码仍可从正式分支到达。其他旧工作树/分支与生产恢复材料均不清理。
- 当前在写最终项目记忆/交接，形成后继`docs`提交后还需再非强制推候选/稳定同新SHA，核远端及CI触发边界。生产镜像BuildInfo仍指代码提交00351，文档后继提交不改应用源码。Edge登录态与付费供应商仍未验。

## 2026-10-15 23:18 CST：v1.5.6传输完毕，候选快进到WIP检查点

- 旧A/用户稳定共同已有官方v1.2.7@371bc68a。临时官方传输ref从此逐步非强制推送v1.5.6的153个第一父提交，远端达到473d535e。第50段GitHub仅对上游自带backend/server约62.67MB警告并接受；第140段批量408，改单提交达151。第152原官方commit b35c66b2新增8份支付产物合计58.6MB，单提交HTTP408：额外临时`codex/transfer-payment-v156-20261015`从原151父提交经备用索引/commit-tree分四组上传原blob，最终临时树与原152树SHA完全一致，再把未经改写的原152及官方153提交快进到原官方传输ref。临时支付ref最终需清理。
- 候选远端从用户旧A6b16非强制快进到v1.5.6 WIP检查点56e9b124（第二父官方473d已在远端）。稳定仍d9，正式目标b740尚未远端发布。其后v1.5.8、v1.6.0、v1.6.6官方头是25842275/42f14cd7/27c8066e，从473d第一父分别余110/157/257个提交；继续采用同一临时官方ref的10提交批次，若遇大blob单提交408则先按原对象分组暂存、回推原SHA，绝不rewrite/force。直到候选/稳定都确认b740后精确删除本轮临时分支。

## 2026-10-15 19:35 CST：GitHub传输分段处理

- 用户fork`origin`远端候选此前不存在、稳定d9fec0f4、main11931d00。直接push最终b740两次HTTP408且远端未建候选；这不是合并/数据/测试失败，未改全局代理/路由或使用force。
- 已以非强制小包把`codex/upstream-v1.6.6-20261012`按d9→d6fed7→b9cb35→6b16d25e推进，远端逐段确认。`56e9b124`单次WIP检查点携带1757文件/官方v1.5.6第二父历史，再次HTTP408，远端仍6b16。
- 共同祖先`371bc68a`是旧A与官方v1.5.6`473d535e`的merge-base。建立短期`codex/transfer-official-v156-20261015@371bc68a`，计划对153个官方第一父提交分小段快进后再推WIP。候选/稳定目标仍b740ee04；完成后只有在最终分支均验证后才删除临时传输ref。旧A及生产继续保留，未声称GitHub发布完成。

## 2026-10-15 14:00 CST：最终生产二次发布核验完成

- 增量修复提交`00351b4f`从上游AutoDL2.3.0视频归档合入旧2.0.0音频Provider和`indextts2-v1`，同步v2 Manifest/README/interface内嵌合同/zip；新增`build-packages.sh --provider-only`精确生成入口。实际归档加载、音色参考请求、任务创建与WAV轮询单测以及Backend Linux CGO最终全仓通过。Web、Agent和PG/Redis原专项按应用代码无变更继续有效，无真实付费模型调用。
- 最终Backend/Web镜像按代码提交`00351b4f3c922b5acf73e6d17a029a049ef01af9`构建，Backend内插件包SHA与仓库一致。二次切换前Schema51完整备份`v166-audio-before-20261015-132809`相较此前真实恢复证明源没有关键数据变化；保留卷重建backend/web后两者healthy、build提交一致、Schema51/51/任务0。插件109、AutoDL2.3.1+dingxin.1视频＋音频二Provider和IndexTTS2工作流存在，其余旧插件合同与启用状态不变。
- 私有post-audio审计：完整性ok/FK0、解密失败0，关键表计数/旧字段摘要/密钥hash/marker/迁移51项同新鲜备份；旧2a Web静态清单1615项保留，旧v1.2.7代表性JS/CSS仍200且非SPA回退。原v1.2.7和首轮2a镜像、备份、克隆、旧A均未删；当前待交接文档commit/非强制push与远端复核。
- Microsoft Edge未出现在当前电脑控制工具的可用浏览器列表；本轮无登录态页面视觉/交互验收，未回退到其他浏览器。手工指南`docs/plans/edge-v166-production-20261015.md`；F6和真实收费调用继续待验证。

## 2026-10-15 13:10 CST：AutoDL音频合包与Schema51恢复验证通过

- 上游2.3.0归档只含视频Provider及15视频工作流，旧2.0.0归档额外含`autodl-comfyui-audio`/`indextts2-v1`，仓库源manifest最初是1.0.3而归档是2.3.0。以官方2.3.0为基底，精确保留旧2.0.0音频合同，生成2.3.1+dingxin.1源码及归档；用仓库`embed-documentation.mjs`、新增`build-packages.sh --provider-only autodl-comfyui`同步内嵌接口定义，避免源码/文档/zip再次分叉。
- `TestAutoDLPluginArtifact`补验实际包加载、IndexTTS文本/参考音频请求与异步创建/轮询WAV返回。Linux CGO最新全仓`go test -count=1 -timeout 30m ./...`最终退出0。一次全仓中间失败有明确环境原因：为zip生成临时安装的Node24触发原本因无Node跳过的Agent E2E，但隔离测试容器未安装undici；恢复原测试容器工具状态后全仓通过。生产镜像内Node22依赖未因此更改，也未使用真实付费Provider。
- 生产首轮2a/Schema51仍健康、活动任务0。二次切换前私有恢复点`v166-audio-before-20261015-120957`包含完整数据及Web静态资产，完整性ok/FK0/加密检查失败0；旧2a镜像在无网络`open-ai-canvas-v166-audio-restore-proof-20261015`从该备份实际启动到Schema51/ready且业务摘要、密钥、插件、迁移历史与备份一致。旧v1.2.7/Schema9恢复点与旧镜像标签也保留。尚未二次发布或push。

## 2026-10-15 11:40 CST：生产首轮晋升与打包回归

- 用户明确要求继续生产晋升/推送。正式双父提交`2a71a5f8`（父本地WIP4323c759、官方27c8066e），工作树曾清洁；其应用路径与此前2491+5Web、Backend全仓、Agent328测试树相同。新鲜私有Schema9完整备份`v166-predeploy-20261015-100745`，与先前已迁移/恢复证明的备份在关键业务、密钥、marker、插件、迁移摘要完全一致；任务0、完整性ok。
- 新Web卷预置旧971个静态资产，后端14个原运行变量与新Compose一致；源目录切到A‴，生产backend/web均升级到2a镜像、ready及Schema51/51。私有postdeploy审计已记录完整性ok/FK0/解密失败0/任务0，关键业务旧字段摘要和密钥一致，旧静态JS/CSS可达。
- 最后插件合约检查发现旧AutoDL 2.0.0包含`autodl-comfyui-audio`及`indextts2-v1`，上游2.3.0归档仅保留视频；本地源manifest仍含音频但原版本1.0.3，源/归档失配。数据库无AutoDL协议模型行，但潜在能力受损。已以官方2.3.0归档为基底合入旧2.0.0音频块，生成2.3.1+dingxin.1源码、README/docs及归档，添加归档实际加载+创建请求测试。旧归档在忽略的验证目录可恢复；当前定向复测尚在进行。正式修复commit、二次发布和push未完成，不能把首轮2a当作最终版。

## 2026-10-14 23:41 CST：收尾与根接管入口同步

最终源树与应用路径比对无差异，只有记忆/验收说明变化；61显式冲突清零，最终2491+5全部通过。3001最终身份/Schema核验通过，上一候选静态chunk仍返回200/真实类型。私有恢复点后增审计材料的ACL再次验证，仅授权SID。PG/debug/旧Schema9恢复测试容器停止且保留，生产与预览继续运行。根旧分支只增量更新README/CURRENT_STATE/BACKLOG接管指针，既有内容保留，未作为构建源。跨保护目录多文件patch曾部分写入后失败，造成最新记忆段重复；已按相同完整段机械去重，并经主机授权的native apply_patch完成根入口，未重复更改应用源码。无正式commit/push/晋升/付费生成/卷删除，待Edge。

## 2026-10-14 23:18 CST：候选验证完成，交付 Edge 门禁

- 最终Web `bun scripts/run-all-tests.mjs --timeout 60000`：2491/2491主套件＋5/5独立注册表，0 fail；最后一轮1645.48秒，45项持久化专项/类型/lint通过。补了“并发删除连接时回滚不生成undefined记录”保护和双情形用例。`web-full-release.log`、`web-persistence-release.log`、`web-typecheck-release.log`、`web-lint-release.log`为最终证据。
- Backend Linux CGO全仓、实际独立PG16/Redis（不连接生产）用例通过；Agent328/构建通过。未跑Chrome浏览器runner、真实供应商/扣费，不宣称全部端到端。
- 六谱系：生产9、legacy45/local45/canonical45、官方43、固定官方v1.6.6原生50，均独立副本up/up/verify→51，旧迁移身份、业务摘要、key/marker/plugin状态保留。Windows tar对官方中文路径报错，改同一Git归档在隔离Linux解包/编译；官方50夹具本来没有key/registry，审计明确null，不制造密钥、不掩盖有加密记录时的解密失败。
- 私有完整恢复点 `v166-predeploy-20261014-183917` 在NTFS受限ACL；旧镜像保护标签 `pre-v166-20261014`。真实恢复至无网络 `open-ai-canvas-v166-restore9-proof-20261014` 后Schema9/9/healthy/任务0，关键计数/digest/key/plugin全部一致，恢复点登记verified/protected。生产3000没有替换镜像或卷。
- 候选启动业务摘要不变；内置技能33→71、旧33版本与文件记录保留，registry84→109；仅旧三包随上游升版本，启用状态无变化，七个本地关键插件审计摘要一致。其他运行验证与数据变化在handoff明确区分。
- 最终镜像源树 `cbf3ebd381531e887c295931b45a408af8cec501`，构建标识 `candidate-cbf3ebd38153`，版本v1.6.6+dingxin.1。Backend镜像15cc1a2c、Web0d06c840；Compose仅更新独立预览，Schema51/51/healthy/活动0。7路由和入口8静态资源HTTP通过，不等于浏览器UI验收。
- Preview第一次只读DB因root离线迁移属主，精确副本卷chown为候选app解决；Docker internal网络未发布宿主端口，只给Web加frontend网络，backend仍无默认外网路由。PowerShell单次HTTP用-NoProxy，系统/Clash配置未改。
- 交付 `HANDOFF_V166_20261014.md` 与Edge清单。待用户手验确认后继续已获授权的正式发布/push，需新鲜恢复点和运行复核。无正式commit/force/push/生产晋升/数据删除/真实付费调用。

## 2026-10-14 17:43 CST：全量结果与语义回归修复

- Backend Linux CGO `go test -count=1 -timeout 30m ./...` 全仓退出0；Agent `npm test` 328/328、`npm run build`通过。Web 首轮构建、typecheck、lint通过；首轮完整测试2467/23，修复后重跑2486/4，尚非发布通过。
- 本地原生单请求多图和上游展开堆栈冲突：已有 `taskOutputIndex` 的结果节点不再重复展开、不再覆写为主输出；12 项多图及10项版本/入口专项通过。生成连接尝试增加来源快照，回滚只撤销本次生成字段并保留后续普通编辑，44项持久化专项通过。
- 上游 `yingce:`/IndexedDB名称更名会使旧浏览器缓存不可见；精确保留本地 `infinite-canvas:`/数据库名（实现标识，非产品名称），继续使用上游新同步/恢复算法。不添加平行存储兼容层。48项物化/刷新专项通过，刷新节流使用单调时钟。
- 上游新持久日志及删除冲突恢复是有意语义变化，已同步旧夹具：保留本地未确认工作作为冲突副本，禁止其自动重建云端已删除画布；本地成功写入不能清掉远端冲突状态。最后四项为静态运行器/存储命名断言及跨套件残留资产夹具，正在修复；不跳过断言。
- 尚未完成真实副本Schema迁移、PG/Redis依赖测试、最终镜像、Edge、生产晋升或push。旧服务/卷、回退工作树保留，无付费模型调用。

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

## 2026-10-13 约 00:00 CST：隔离 A‴ 已建立并开始固定合并

- A″ 的全部暂存内容连同本记录启动段写成树 `04de369256a308373eb6b8036dccc61d459672b8`，`git commit-tree` 生成仅本地 WIP 双父 `4323c7592471ba04481f696b1052f55f183ff05c`（父 `7419e2b3` 和官方 v1.6.0 `42f14cd7`）；A″ 原 HEAD/MERGE_HEAD 未动，未 push。
- 新隔离工作树 `C:\Users\DingXin\.codex\worktrees\upstream-v166\open-ai-canvas` 创建在该 WIP SHA，分支 `codex/upstream-v1.6.6-20261012`。`git fetch` 后标签 peeled SHA 与远端一致：`27c8066e6da9f0d4754dc0e5c40097049fe730cf`。相对 v1.6.0，官方 1,606 路径，本地 732 路径，重叠 195。
- `git merge-tree` 预判冲突；实际 `git merge --no-commit --no-ff v1.6.6` 的 `MERGE_HEAD` 为标签对象 `8cd7da6d`（应进一步核实剥离到 `27c8066e`），`git diff --name-only --diff-filter=U` 为 61 路径。包含一个上游删除 vs 本地修改的画布版本比较弹窗，其余以文本冲突为主；不能简单整文件 ours/theirs。尚未解决或测试，不得部署/push。

## 2026-10-13 01:31 CST：冲突分批边界

- 逐文件计数确定 61 个未合并路径；`web/src/components/canvas/canvas-version-compare-modal.tsx` 与 `web/src/components/canvas/director/canvas-director-workbench.tsx` 是 2 个 modify/delete（无文本标记），`web/src/pages/canvas/project.tsx` 有 27 个文本冲突块。其余模块多为 1–4 块，但不能只按块数决定风险。
- 处理顺序：基础 `VERSION`/CI/依赖；Backend 模型、数据库、协议和资源/任务；Web 模型合同、状态及路由；画布、导演/预演、创作/分镜；对应测试；最后专题文档与项目记忆。每批复核调用点及双方阶段内容，允许 synthesize、保留一侧或标 needs-human；未理解互斥产品目标不自动赌。
- 当前没有运行 v1.6.6 编译、迁移、镜像或 Edge。生产 v1.2.7/Schema9 是实测现状，不可从旧 v1.6.0 A″ 的 Schema45 克隆结果推断实际部署安全。

## 2026-10-13 03:26 CST：基础/Backend 第一批合成中

- `VERSION` 合成为 `v1.6.6+dingxin.1`；质量流程与 `web/package.json` 同时保留本地完整测试/Director E2E 和上游预演 E2E。当前机器遵守 Edge-only，不在此处运行 Chrome 脚本。首三文件 JSON 语法检查通过，仍未跑全部 CI。
- 上游 Go module 改为 `yingce/backend`；逐项确认后，对 26 个本地 Go 源文件的旧 `infinite-canvas/backend/` 导入做精确机械替换，生成二进制 `backend/server` 未触碰。此项是独立的自动合并语义风险，必须用 Linux CGO 编译/测试复验。
- Backend 已合成数据库连接（上游 PostgreSQL 连接前配置 + 本地脱敏 GORM logger）、渠道结构/管理视图（上游 `ProxyURL` + 本地 `AllowLocalChannel`）、调用链（本机回环安全直连优先，普通渠道代理）、资源/结果下载和视频能力约束。上游 Bailian 补充目录 20 个 ID 与本地受版本管理的 80 项 manifest 比较，20 个均已包含；保留本地逐模型支持状态与文档合同，不降级为泛化插件目录。
- 数据库新候选暂定 Schema 51：canonical 44 本地时长、官方 45–50、51 本地 Comfy；历史本地 45 已是 Comfy 时，保留该记录并把官方上传预留迁移移到 51；更早本地 v9 与 v41/42 谱系保留旧名称/校验和，官方 46–50 再接 51 上传预留。新添计划测试，但 Windows `go test -run '^TestMigrationPlan' ./internal/database` 因 `CGO_ENABLED=0` 的 go-sqlite3 stub 全部无法运行，属于环境阻断而非通过/失败结论；独立 Linux CGO 测试镜像构建中，尚未验证这些计划或真实克隆，**生产迁移仍 NO-GO**。
- 当前仅暂存已解决的冲突文件，仍有约 43 个未合并路径；主要为画布、预演台、模型/页面/样式与文档。未正式提交、推送、生产部署或真实付费调用。

## 2026-10-13 09:47 CST：Backend 专项回归与画布底层收口

- Linux CGO 无卷容器 `go test -run '^TestMigrationPlan' -count=1 ./internal/database` 通过；Windows 同一组因 CGO stub 未运行，不能把后者计入代码失败。Backend 全仓首轮未通过，仅 `internal/app` 两项（图片质量/分辨率档位测试夹具缺 BillingMode；排队任务测试误把应脱敏的持久化错误要求成原始异常）及 `internal/handler` 一项（登录身份防枚举测试误把每次随机 requestId 当作相等字段）。修正的是三个测试夹具/断言，未放宽实际价格校验或错误脱敏。三项在隔离 Linux CGO 调试容器中定向复测 `app`、`handler` 均通过；全仓最终重跑仍待。
- 前端模型合同、主题和入口：上游图片价格选择器与 `1.5k` 能力接入，保留本地百炼 video generationContract；用户切换使用上游先恢复缓存再开放工作区，保留本地账号清缓存；主题保留旧持久 key 和画布临时覆盖，首屏读取旧/新 key；本地首页独立、上游 `/inspirations` 加入导航、`/create` 保留上游游客能力；管理员插件中心继续有本地权限旁路。以上尚未运行 Web 类型/浏览器验证。
- 画布底层：新 `CanvasSettingsPopover` 替换旧视频专用 Portal（旧文件可由 Git WIP 恢复）；上游 `CanvasViewport`/预演 Dock 名称、缩略图和 MediaPackage 导出采纳，导出仍支持本地绘图；本地 Comfy Bridge workflow 配置、单任务多图节点输出序号与上游逐层提取/素材保存前合成检查已合成。未解路径降至 19，剩项目画布 27 块、导演→预演数据/组件删除与其他页面、样式、测试和文档。不可把这一阶段称为 UI 通过。

## 2026-10-13 13:53 CST：用户选择上游版本历史移除行为

- 用户对“官方 v1.6.6 删除本地画布版本历史/预览/恢复 UI 及调用链”的明确回复是：**采用上游删除方案**。因此在 A‴ 中清除旧 UI/调用方和先前暂存保留的旧比较弹窗；历史代码仍可从本地 WIP `4323c759` 恢复，旧持久版本数据/备份不得因本决策删除。最终交付须把这是有意接受的功能变化说清楚，而非误报为无影响。
- 原导演台→上游预演台是另一项独立数据边界：已增加纯转换 helper，将有效旧 `directorScenes`/节点关联映射到 `previsScenes`，保留原字段并对不完整结构拒绝强转；上游新预演台已有保存失败可留在界面的关闭流程。相关单测待 Web 依赖/类型与 Edge 验证，不能先宣称旧场景运行验收完成。

## 2026-10-14 15:21 CST：冲突清零与全量启动

- A‴ 61 个显式冲突已全部解开/暂存，应用路径无残留标记。项目画布 27 块按用户版本历史决策采用上游新布局/预演调用，同时保留本地人像筛选/Modal、百炼配置模型和参数；旧版本比较弹窗已移除，旧版本数据/备份不删。提示词编辑器采用上游引用组件并补回百炼专属参数面板；资产标签编辑保留本地读写规则。
- 样式仅按冲突精确合成；上游 `globals-*.css` 分片目前未被实际 import，没有据此盲删主 CSS。新视口的本地网格强度只乘一次，回归路径同步为 `canvas-viewport.tsx`。文档保留上游新功能并明确候选 Schema51 和历史迁移身份。
- Windows 沙箱查询 `.git` 曾报 not a git repository，升级为只读主机权限后确认工作树登记/分支/对象完整，未修复或重建 Git 元数据。Docker Engine 当时未运行，按已授权测试/部署目标启动 Docker Desktop；旧生产容器随其恢复，当前尚未替换镜像。
- Web `bun install --frozen-lockfile` 1440 包成功，Agent 冻结安装 152 包成功，无锁改动；Web typecheck/lint 通过。七文件专项 75 pass/4 fail：三项 Windows worker 超出默认5秒（非断言失败），一项旧 infinite-canvas 文件名静态断言已修正。新场景转换2项与预演接线34项已通过。测试运行器增加参数透传，以 `bun scripts/run-all-tests.mjs --timeout 60000` 跑完整清单，断言未跳过；Agent npm test/build、Backend Linux CGO 全仓同时运行，尚未有最终退出码。
