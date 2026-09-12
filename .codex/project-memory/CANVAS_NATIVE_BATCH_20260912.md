# 画布 Qwen 原生单请求多图

状态：`deployed_preverified_edge_pending`。用户要求补齐创作页与画布的单请求多图差异；本次为前端增量，继续停在 Microsoft Edge 人工门禁前。

## 当前已确认

- 正确工作树 `.local/worktrees/dashscope-plugins`，分支 `codex/dashscope-system-plugins-20260911`；在 d2fe747f / 5889345a 交付上追加代码提交 cd0b434060ef8ca9a967ab2c9cc1ce824a7c043b，未切换旧根分支、未 push。
- 创作页和画布均有显式「单请求多图」，默认关闭。适用 Qwen Image 3.0 Pro 及声明相同 generationContract 的逻辑模型；不对工作流或不同合同静默生效。
- 画布开启后提交一次 count=N，生成 N 个子节点；节点持久保存 taskOutputIndex / taskSubmissionNodeId。结果按序号附着、缺失输出报错、被删除的目标不回退覆盖兄弟节点。
- 使用现有统一任务物化和持久副作用确认；共享任务的输出顺序消费，部分输出失败不吞掉其余成功结果，不自动重新付费生成。
- 刷新恢复按任务和节点组合去重，可根据提交锚点找到尚未保存 taskId 的同批次子图。单张重试请求 count=1，按原输出序号隔离重试幂等键；接收新任务后清除旧输出序号，并保留生成参数。
- 报价配置带入节点高级选项：单请求按一次报价，默认拆分按请求数报价；模型单价与策略不改写。

## 验证证据

- `.local/canvas-native-batch-final.log`：干净 TypeScript、84 专项（含本次 12 项）、1645 全量（1640 + 隔离 5）、生产构建通过；断网 Docker 执行，无真实 Provider。
- 9 个修改源码/测试文件与验证镜像 SHA-256 一致。测试镜像 `open-ai-canvas-web-test:canvas-native-batch`；发布镜像从此镜像构建，只增加版本标识。
- 新用例覆盖一次请求/三个独立输出、默认三次请求、单图、部分失败、整体失败、取消、逻辑合同与工作流隔离、JSON 恢复、删除目标、缺失输出、单图重试幂等、报价配置及物化输出幂等。
- 复查发现节点报价配置缺少 generationOptions，已补齐并加测试。首次新增报价测试缺少有效模型渠道，已修正夹具，未放松生产模型选择规则。
- 部署脚本曾在容器替换前触发 PowerShell `-p` 参数绑定歧义；确认旧容器/镜像标签未变后改为数组参数传递，再执行通过。

## 最终部署

- Web：v1.2.7+dingxin.1.bailian3 / cd0b4340，容器 0e8cec3338e7d2ca8c92478a7a250811646ec6174c0c4db3c837081f3bacbf54，镜像 sha256:44380b76e9aaaab801832e577716639d144ee04fdc6608e5403208d2381fba20，healthy。
- Backend：仍是 v1.2.7+dingxin.1.bailian2 / 5889345a，容器 7f9ad8088cc0599b4723e1b7a3d50734b37e992571c34697793bf7974f77a301，未重建。Schema 9/9、插件 3.1.0 和数据卷不变；本轮没有数据库写操作。
- 前后端版本不同是本轮仅前端升级的预期状态，不按“版本必须完全相同”误判。Web 环境、端口、重启策略与原容器一致，Web 无数据挂载。
- 旧 Web 镜像 f53d523c579c6d9831262a5808375e3e65b2bfe53db68ab916c0306b5c9a19b2 保留为 `open-ai-canvas-web:pre-canvas-native-batch-20260912`。回退仅替换 Web，不回退 Schema 或 Backend；旧数据库恢复点全部保留。
- 部署证据 `.local/canvas-native-batch-deployment.json`，发布构建日志 `.local/canvas-native-batch-release.log`。

## Edge 人工门禁

1. 已登录 Edge 硬刷新，确认前端 bailian3 / cd0b4340；后端仍为 bailian2 / 5889345a。
2. Qwen 画布配置节点和图片节点的高级参数均可切换「单请求多图」，默认关闭；切换后报价按一次/多次对应变化。
3. 修改张数、开关和种子后刷新，确认恢复；同步/异步/OpenAI 兼容模式入口仍在。
4. 本轮不要点击生成或模型测试。真实生成后的一任务多结果、展开/主图选择、页面刷新及单图重试须在另行费用授权后验收。
5. 发现字段丢失、报价异常、控制台错误或版本不符时停止保存/生成并反馈。未执行 Edge、真实生成、模型启用、价格修改、远程推送。

本记录替代 BAILIAN_FULL_CAPABILITIES_20260912 中“Canvas 多图仅独立任务”的限制；其他七模型能力及人工/付费门禁保持。
