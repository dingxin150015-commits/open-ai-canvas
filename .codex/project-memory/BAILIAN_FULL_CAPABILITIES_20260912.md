# 百炼七模型完整能力补齐

状态：`deployed_preverified_edge_pending`。最终运行代码 5889345a 已部署，非浏览器预验证通过；停在用户 Microsoft Edge 人工验收前。未进行真实模型生成。

后续增量见 `CANVAS_NATIVE_BATCH_20260912.md`：画布原生单请求多图现已补齐。Web 更新为 bailian3 / cd0b4340，Backend 保持本记录的 bailian2 / 5889345a；本记录原“画布仅独立任务”限制及前后端版本一致检查已被该增量替代。

## 授权与基线

- 用户批准修复审计中的全部使用阻断、高级参数、媒体/引用和计费缺口，最后停在 Edge 人工验证前。
- 用户确认计费采用推荐方案：保留现有自主零售价和历史订单，新增可选实际用量计费，不自动切换模型策略。
- 用户另行确认将未配置时的默认图片任务超时从 8 分钟提高到 15 分钟；显式已保存策略及其他任务类型不变。只读生产检查确认无活动任务、Schema 8、未配置图片超时覆盖。
- 正确工作树为 `.local/worktrees/dashscope-plugins`，分支 `codex/dashscope-system-plugins-20260911`，起点 13ff9e42，集成基线 d9fec0f4。旧根目录 b6254b41 不作为开发或部署基线。
- 不执行 Edge、不调用真实 Provider、不 push、不启用模型或改写现有价格。外部模型顾问没有调用，使用工程编排技能的 local-only 路径。

## 实现范围与证据

- 七模型：qwen-image-3.0-pro、happyhorse-1.1-{i2v,r2v,t2v}、wan2.7-{i2v,r2v,t2v}，保留已登记快照。
- 创作页附件按图片/视频/音频及视觉总量分别判断；支持 T2V 音频、首尾帧加音频、五视觉加参考音频。素材库选择也检查数量。
- 图生视频空提示词：修复 Create、Canvas、后端 admission 和 worker 的独立阻断，保持其他模型必填规则。
- 首尾帧按模型区别显示；旧失效帧允许用户清除，创作页删除/替换素材维护帧 ID。
- 高级参数进入创作草稿和 Canvas metadata；负向提示词、种子、扩写、Qwen 思考/扩写模式/水印、请求方式；显式 false 不丢失。
- Qwen 原生同步、原生异步和 OpenAI 兼容仍通过同一个声明式插件；创作页可选原生单请求多图，默认独立任务行为不变。
- 媒体引用按类型分别计数，HappyHorse R2V 转为原生 Image 标记；Wan R2V 首帧不占 reference_image 编号。续写视频增加尺寸与比例校验。
- 修正错误的长音频拒绝规则和对应测试；新测试不只对照旧 Provider，还覆盖官网合同。
- 新计费模式 `per_second_usage` 为百炼视频渠道规格的显式选项；保留 `per_second`/`fixed_request`，逻辑模型可通过渠道价格策略使用。预估后按可信终态 usage 结算并退款；缺用量或实际费用超出预扣时转人工核对，不自动负余额补扣。旧 Token 策略不变。
- Schema 9 `provider_duration_metering` 为 API 日志增加三种媒体计量时长，为订单增加实际计量时长，不修改旧迁移校验和。

## 验证进度

- 超预扣人工结算使用已有管理员确认和必填核对依据入口，并原子检查余额足够；自动任务不得调用复核路径。
- 逻辑模型使用内部 ID。运行时从有效线路导出公共 generationContract，前端据此开放专属控件并绑定请求合同；路由拒绝跨合同降级。不同合同混合线路不冒充单一合同，需拆分模型使用专属能力。不暴露渠道 ID、密钥或线路关系，不改写既有模型数据。

- 第一轮干净类型检查发现 optional input 与 selectedModel 引用问题，已修正；第二轮类型通过。
- 首批前端专项 34 项通过；后续修改仍需最终全量重验。
- 最终 `.local/bailian-backend-release.log`：Linux CGO `go test -count=1 -run . ./...` 全量通过，service 254.539 秒。
- 最终 `.local/bailian-web-release.log`：干净类型、34 专项、1632 全量（1627 + 隔离 5）、生产构建通过。
- 测试镜像与最终工作树逐文件 SHA-256 比对：Backend/插件 31 文件、Web 25 文件完全一致；默认图片超时 15 分钟和空提示词媒体准入均已包含。
- 上述为主体检查点。追加收口最终证据：`.local/bailian-backend-contract-final.log` 全量通过，service 328.882 秒；`.local/bailian-web-final-verified.log` 类型、35 专项、1633 全量（1628 + 隔离 5）及构建通过。Backend/插件 32 文件、Web 26 文件与测试镜像逐字节一致。
- 全量中发现并修正 OpenAI 自动尺寸 nil 替换失败；旧前端“提示词必填/音频必须短于视频”断言按新合同更新后全量通过。
- 以上早期部署状态已被下方最终部署记录替代，旧 Schema 8 / 8dcdfd5d 仅作为恢复基线。

## 最终部署与数据证据

- 代码提交：5889345a8f2f861ecd0a9bd8bdf67751fc354275；前置主体 b9cb3534。未 push，稳定分支未晋升。
- Backend 与 Web 均为 v1.2.7+dingxin.1.bailian2 / 5889345a；Backend 容器 7f9ad8088cc0，Web 容器 b0c46e44efef，均 healthy / restart 0。
- Backend 镜像 sha256:f5edb2742ece697a024df713e333ba8b236579dd28bed6012de20128d698f5a2；Web 镜像 sha256:f53d523c579c6d9831262a5808375e3e65b2bfe53db68ab916c0306b5c9a19b2。
- 健康入口返回 ready=true、Schema 9/9、活动任务 0；前端 build-info.json 版本和提交与 Backend 一致。未登录后台 GET 返回 HTTP 401，严重启动错误筛选为 0。
- 原运行环境逐项一致；真实卷仍为 open-ai-canvas_backend-data。升级前后模型、价格、逻辑模型/版本/线路、订单、账户和流水摘要不变，数据库完整性 ok、外键违规 0、密钥解密失败 0。只追加批准的 Schema 9。
- 新恢复点：BKP-20260912-203916-BAILIAN-COMPLETE-PREDEPLOY / BKP-20260912-204448-BAILIAN-COMPLETE-POSTDEPLOY，均 verified/protected。私有位置与镜像身份见 Git 忽略的 .local/deployment-complete-state.json；所有旧恢复点保留。
- 两个官方插件均启用、版本 3.1.0。运行包与源码 SHA-256 相同：image a3a8def33990ee3ee45b78238cb7d9049cf5ccf9076aaba1a4feb90f7180692e；video 624c9543eb81a9dc751704f6ecf8a4ca1b4174f6dd76630a4d787651ada4b0fa。
- 断网克隆双启动和数据摘要核对通过；仅删除了本轮已停止的专用克隆容器和测试卷，私有克隆快照及正式恢复点保留，可重建测试环境。
- 运行候选从逐文件验证的测试镜像构建；Backend 复用原运行基线和未修改的原生支付插件，仅替换已验证服务器、迁移器及两个百炼包。

## 用户 Edge 清单

1. 使用已登录 Microsoft Edge 硬刷新。核对版本 v1.2.7+dingxin.1.bailian2 / 5889345a；插件中心两包为 3.1.0。
2. 后台七模型：Qwen 选择 DashScope Qwen Image；六个视频选择 DashScope Wan / HappyHorse Video。旧 Wanx/Wan 协议不能替代。
3. 原启用状态、价格和订单必须保持。若模型仍停用或未定价，前台不可见属于正常门禁；由用户确认价格和启用配置，代理没有代为启用或设价。
4. 创作页和画布检查高级参数及刷新恢复：seed=0、显式 false、负向提示词、扩写；Qwen 的思考、Direct/Agent（图像编辑禁用 Agent）、水印及三种请求方式。
5. 使用已有且获准使用的素材检查输入入口：Wan T2V 音频；Wan I2V 首帧/首尾帧加音频、续写；Wan R2V 五视觉加参考音频；两个 I2V 空提示词。此轮不点击发送，不上传未经确认的真实云素材。
6. HappyHorse I2V 只显示首帧，HappyHorse R2V 无首尾帧，Wan R2V 只有首帧。失效旧设置可清除；删除/替换参考图不保留悬空帧 ID。
7. Qwen 单请求多图已扩展至创作页和画布，默认每张独立任务不变。画布逐输出映射到子节点，单图重试只请求一张；详见后续 CANVAS_NATIVE_BATCH_20260912。单请求模式按所选按次价格，不等同于供应商按图成本。
8. 逻辑模型专属能力来自共同 generationContract；不同合同混合线路需拆分，不做静默跨合同降级。更换合同后旧专属参数有清除提示。
9. 默认图片任务超时 15 分钟；已有显式超时不覆盖。实际时长计费为可选项，原策略不自动切换；超预扣需管理员明确核对，补扣时余额必须足够。
10. 若字段丢失、引用错位、报价异常、控制台报错或版本不一致，停止保存/生成并反馈 Edge 证据。本轮不要点击“测试模型”或真实生成；账号权限、额度、费用和生成质量仍须另行授权验证。

## 最终门禁

源码、包一致性、类型、全量回归、生产构建、克隆迁移及保留卷部署预验证已完成。当前只等待用户 Edge 人工确认；真实生成、上游账号权限和生成效果不记为已验证。
