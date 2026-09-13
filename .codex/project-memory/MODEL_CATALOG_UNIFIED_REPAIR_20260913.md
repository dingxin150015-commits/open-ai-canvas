# 系统模型目录统一修复交付

状态：`deployed_preverified_edge_pending`。用户批准统一方案后已完成实现、全量回归、数据副本预演与部署，停在用户 Microsoft Edge 人工验收前。没有真实模型调用或自动启用/定价。

## 基线与版本

- 正确工作树 `.local/worktrees/dashscope-plugins`，分支 `codex/dashscope-system-plugins-20260911`，本轮起点81896044，继承固定集成基线d9fec0f4；旧根分支b6254b41不用于构建或发布。
- 主体提交ad3efa10ee0854e65fe67916d1e530b478c6ad04；Gemini文本插件防丢失修订cc7d0b63dd2d044701b2606bae5b70d846dbadab。未push、未稳定分支晋升。
- 前后端主程序均为v1.2.7+dingxin.1.catalog1 / ad3efa10。收尾按插件独立版本部署：Gemini文本4.0.1，其余本轮六个协议包4.0.0，包合同均v4；收尾未改变主程序二进制或前端。
- 最新Backend容器b7497517f291a9448dfba9ac4bbddccd59b7498b10e79c4f8bde9b5d5262a938，镜像sha256:d67ee794a880f01d436de970b6016db6bbd37096c2edbc50de36230bf8551081。
- Web容器1479e3582b95e1fbf30f7a62b08e43e8614df53b51915eb0c0452a66cca7fae2，镜像sha256:9bfc9c671a45a79e09790987579bca4d9a7c08ecb3f8ac67f077be288a7d8bf3。

## 已实现

1. `evaluateChannelModel`区分可编辑、配置完整、协议支持、可启用与模型级可执行资格，返回原因码、说明和协议版本。目录标记不再锁死表单，前端列表与筛选使用运行资格。
2. 新发现条目默认unknown而不是planned；手工新增不再凭来源直接获得ready。既有planned未配置记录无需删建，可原地完善配置；已有价格、启用状态、模型/路由ID和引用不批量修改。
3. 保存停用草稿强制enabled=false；真实保存、报价、系统渠道解析、逻辑路由、计费和任务入口采用公共资格判定。移除Repository两处ready过滤，仍保留启用查询和service强校验。
4. 管理员`POST /api/admin/channels/:id/models/validate`及“检查配置（不调用上游）”复用配置/价格校验和资格评估，不发供应商请求、不创建模型/订单/任务、不分配ID。模型测试仍是独立真实调用入口。
5. v4在`contributes.providers[]`中增加可选modelContract及modelCatalog，不另建运行时。exact模型集合同时用于配置和BuildCreate/BuildAgent；compatible协议允许配置新ID但不声称账号可用。未使用新字段的旧包沿用原协议合同。
6. 百炼补充目录从当前正式协议包注册视图投影；静态目录仅作生成资料，不再由service另外注册始终可用的bailian发现实例。安装/卸载/启停使路由缓存失效，不自动改模型开关或价格。
7. Gemini解析displayName、supportedGenerationMethods，按nextPageToken同源翻页；100页、10000条、2分钟上限及重复游标检测，部分错误不导入部分结果。页面100条/页与此无关。歧义输出类型由管理员选择，不从generateContent或名字推断完整能力。
8. Gemini普通文本队列可使用安装的Create合同；文本插件4.0.1明确拒绝普通非Agent路径上的媒体附件，避免静默丢失。多模态Agent使用独立Agent合同，不受此纯文本限制。缺参考素材的测试提示不再冒称凭据/端点已验证。

## 验证证据与发现

- 初次后端全量发现10项失败：旧planned不可编辑/删除断言及缺失合法持久能力配置的测试夹具。按新合同更新断言和夹具，未放宽生产校验；定位项专项通过。
- `.local/catalog-repair-backend-final.log`：Linux CGO `go test -count=1 -run . ./...`全量通过，service244.801秒。
- 收尾`.local/catalog-repair-backend-guard-final.log`：后端再次全量通过，service153.684秒；Gemini普通文本本机模拟成功，图片/视频/音频三个负例确认供应商请求数不增加。
- `.local/catalog-repair-web-final.log`：干净TypeScript、40专项、1645全量（1640+隔离5）、生产构建通过。
- 后端/插件65文件、前端5文件与测试镜像逐文件SHA-256一致；收尾5个变更文件另行核验。最终7个运行协议包全部与已验证源码匹配；主程序二进制与主体发布镜像逐字节相同。
- 真数据断网克隆各完成双启动；Schema保持9/9、80表、完整性ok、外键违规0、解密失败0。模型、价格、逻辑模型/修订/线路、任务、订单、账户、流水、资源身份等关键摘要部署前后相同。
- 前后端healthy、restart0；活动任务0；后端严重启动错误筛选0。未登录管理GET及新增配置检查POST均HTTP401。
- 实际数据卷仍为open-ai-canvas_backend-data，后端环境、前端环境及端口保持。没有新增Schema版本，不直接改生产support_status。

## 恢复点与运行边界

- 主体：catalog-predeploy-20260913-143843、catalog-postdeploy-20260913-144131；收尾：catalog-predeploy-20260913-145155、catalog-postdeploy-20260913-145410，均完成独立关键文件比对、私有ACL和数据库审计。
- DB/WAL/SHM（存在时）、设置密钥、存储迁移标记、插件registry分别校验。完整私有位置见Git忽略的`.local/catalog-repair-deployment.json`和`.local/catalog-guard-deployment.json`；旧备份全部保留。
- 专用克隆open-ai-canvas-catalog-preflight-20260913及open-ai-canvas-catalog-guard-preflight-20260913均已停止且network=none，保留测试卷和私有克隆快照；未删除生产卷。
- 回退必须匹配主程序、协议包与插件registry/安装状态。旧v3宿主不能直接沿用v4包；不可只换旧镜像就宣称恢复完成。若交付后产生新配置/订单，须先评估增量数据，不盲目覆盖为旧备份。本轮未做生产回退演练。

## 用户 Edge 验收

1. 在已登录Edge硬刷新，核对主程序catalog1/ad3efa10；插件中心Gemini文本4.0.1，Gemini图片/Veo、两个百炼、OpenAI Chat Completions/Responses为4.0.0。
2. 打开原Gemini渠道模型列表，已有条目应能“编辑”，不需要删除重建；选择正确的文本/图片/视频能力和协议，检查字段及已保存配置恢复。
3. 使用“检查配置（不调用上游）”，检查缺协议、缺价格等原因可读；不要点击“模型测试”或生成。
4. 如需保存不完整配置，可使用“保存停用草稿”，保存后不得进入前台可执行目录。正式启用必须明确配置有效价格，零价格也必须显式确认。
5. 核对用户已启用Qwen的价格、开关、高级参数、逻辑线路及画布单请求多图设置保持；以前的Qwen人工反馈是历史证据，不自动代替本次新版本验收。
6. 新配置Gemini时不要把TTS/Embedding等全部当文本，缺少执行协议时保持草稿。68条planned百炼补充模型没有被本修复批量实现或放开。
7. 若需要重新拉取完整目录，由用户确认后操作；本轮没有使用真实渠道Key再拉取Google模型，因此“该账号是否确实还有第二页”仍未实测。
8. 字段丢失、报价异常、协议错误或版本不符时停止保存/启用并反馈。真实模型权限、额度、生成效果、计费及媒体恢复均需另行授权验收。

本记录替代此前Gemini“未修复/待授权”的诊断状态，保留初始诊断作为历史。项目记忆同步不包含真实API Key、Cookie或账号数据正文。
