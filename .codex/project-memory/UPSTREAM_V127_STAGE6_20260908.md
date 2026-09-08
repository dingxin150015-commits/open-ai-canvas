# v1.2.7阶段6执行记录

## 收尾状态：自动化源码门禁通过，人工E2E门禁未闭合

源码提交 `8808750199549523cc699087db10ca6be828901f` 已固化；本文件随独立记忆提交交接。当前 `v127_stage6_automated_pass_manual_e2e_pending`，不能写为阶段6无条件完成或发布就绪。没有推送、部署、挂载真实数据库或操作浏览器。

| 自动化检查 | 结果及证据 |
| --- | --- |
| Backend Linux CGO全量 | 通过；stage6-backend-second.log。依赖URL校验的用例需要DNS，测试数据仍为独立临时目录 |
| Web完整测试集合 | 1624通过/0失败（1619+独立注册表5）；包含test与src同置测试，stage6-web-final.log |
| Web干净类型/生产构建 | 通过；最终构建28.43秒，使用incremental=false |
| Canvas Agent | Linux322通过/0失败；6个Windows专用用例由本机补跑6通过/0跳过，候选与该Agent源码一致；TypeScript构建通过 |
| Compose | 六文件+deploy/build叠加共7组合解析通过，仅假占位变量；没有启动服务 |
| Windows脚本 | 两个PowerShell脚本AST解析通过；未执行启动脚本 |
| 格式/结构 | 变更Web格式、Go格式、JSON与workflow YAML、diff检查通过；80归档目录和80Manifest校验通过 |
| 版本与身份 | 前端build-info.json与后端ELF静态字符串均为v1.2.7+dingxin.1 / 8808750199549523cc699087db10ca6be828901f；未运行服务器二进制 |

日志与只读审计工具均位于主工作区 `.local/audits/v127-20260908-110746/`。Windows专项日志stage6-agent-windows.log；Agent全量stage6-agent-third.log；身份读回stage6-identity-readback.log。Docker目录bind不可用后仅用审计代码构建上下文静态读ELF，没有数据挂载。

### 自动化证据覆盖的功能族

| 功能/保护合同 | 代码与测试证据 |
| --- | --- |
| Schema6双谱系→8、活动code索引、资源播放字段 | backend/internal/database及其全量测试 |
| Ready/Planned、目录补齐、批量原子删除、精确SKU | channel_models/model_router/finance及v127_merge_contract_test.go、model-pricing测试 |
| Qwen/Wan/HappyHorse、xAI拒绝、媒体URL/OSS/归属 | provider系列、protocol、resource测试，Backend全量 |
| 任务租约、恢复/取消、旧结果拒绝、日志脱敏 | task_worker/runtime_policy/provider_error及相关全量测试 |
| Create Backend-only、元数据/声音水印、用户scope | Create/生成持久化与跨Runtime测试 |
| 导航、主题、自定义网格、实际宽度控件 | v127-ui-replay/workspace-theme/canvas-appearance测试；像素验收仍待用户 |
| 官方编辑器、命令/撤销、插件slot、媒体交互 | editor-*、timeline-*、canvas-*全测试集合；slot空初态隔离执行 |
| Agent协议、30MB/敏感信息边界、插件引用文件 | Agent全量及Windows专项、80Manifest读取校验 |
| 发布与部署边界 | 标签合法性、构建身份、Compose解析；生产卷迁移/真实启动不在本轮 |

这些是源码/单元/构建证据，不等于每项功能的真实账号、付费Provider或生产数据验收。Host Updater仍不得用于fork自动更新，Compose中官方默认镜像源也不代表可直接部署本地fork。

### 构建标识

- 显示版本：`v1.2.7+dingxin.1`；Docker可用标签：`1.2.7-dingxin.1`。
- Web测试构建manifest：`sha256:d6bb6738ff6a1069dcc8d755c2f0bcb5d78a505eca5a9e1f11359ea4d325e342`。
- Backend身份验证构建manifest：`sha256:435217453934b83ed81d65d8df259649ca941907417ae8e4f9cf95e7e3b61ecd`。
- Agent测试构建manifest：`sha256:ae2b8d1987bd4128afc7df3accb0b90dfee718f9807a69ad2dadeb909fb958b0`。
- 以上为隔离验证产物，未替换现有容器，不是生产迁移或用户视觉验收通过证明。

以下为执行过程记录。

用户仅批准阶段6最终源码验证/候选固化，基点61ce2d3d。未授权push、部署、真实数据库或浏览器操作；进行中，不能视为阶段6已通过。

- 全Web目录首次1611通过/11失败，发现默认链外的过时控件断言及editor-slot-registry跨文件静态注册污染。保留全部测试并将空初态注册表测试隔离进程；其他断言按实际上游控件和阶段5合同更新，不恢复旧实现。
- Backend断网测试因既有公网域名DNS校验失败（含任务恢复前置URL检查）而失败；改回隔离Linux但允许DNS，不接触真实数据/Provider。Agent首轮缺跨目录fixture文件，补齐Web、插件与Compose源文件后再验。
- 官方日志归档本轮可读取：固定运行34138313829为51通过/7失败。A13/B9/C9/D9/E10/F8涉及复现台/api/public/appearance 502，A13另有favicon404；F6为选择留在导演台后弹窗未消失。后者未浏览器复验，不猜测或标为修复。
- 发现发布workflow直接把VERSION含+值用于Docker标签，已映射为1.2.7-dingxin.1并防止fork构建被自动标为无后缀版本。Web补齐构建commit注入及build-info.json，便于后续与Backend身份核对。
- 显示版本的build metadata在SemVer比较中被忽略，Host Updater仍不能据此区分dingxin补丁，继续保持不启用；不改变其仓库/镜像/祖先控制边界。

## 自动化门禁进展

- Backend隔离Linux CGO全量已通过，日志stage6-backend-second.log；包含目录/精确SKU/安全归属/任务租约/资源/迁移与版本比较回归。
- Web全测试集合同时发现test与src下同置测试，不遗漏旧素材测试；注册表测试单独进程。最新1624通过（1619+5），0失败，干净类型与生产构建通过，最终构建身份复核进行中。
- Canvas Agent原命令全量：Linux322通过、0失败、6个Windows专用用例跳过；在与候选代码一致的主工作区复用现有依赖，Windows专项6通过、0失败、0跳过。Agent TypeScript构建通过，未启动真实Agent。
- Compose六文件及deploy+build叠加共7组合解析通过；缺CORS的负向解析按预期失败后补仅解析用假Origin，没有放宽生产校验。PowerShell两脚本仅AST解析通过，未运行启动脚本。
- 80个插件归档可读取且80个Manifest JSON解析通过，无越界条目；改变的JSON解析和diff检查通过。最终源码提交8808750199549523cc699087db10ca6be828901f已固化，但这不是F6/视觉验收放行。

## 明确保留门禁

官方F6确认框交互尚未由浏览器复验，本轮不操作用户Edge或改用Chrome。其余6项网络失败已增加仅测试模式启用的精确fixture，并通过非浏览器测试，但也不宣称整个E2E已重跑通过。用户随后明确同意将F6与其他视觉项统一留到阶段8，由用户手工Microsoft Edge验收，并批准阶段7。因此F6不再阻挡候选分支推送，但仍是阶段8的必测项，不能标为已通过。
