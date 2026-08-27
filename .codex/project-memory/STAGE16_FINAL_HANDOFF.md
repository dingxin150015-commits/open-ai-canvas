# 阶段 16 最终本地交付报告

## 结论

影策阶段 1–16 的本地开发、测试、数据保护、运行部署和交付固化已完成。当前状态是“本地发布候选通过”，不是“远程发布完成”：没有 push、PR、远程 GitHub Actions、Git tag 或 Release。

## 最终源码与运行基线

- 分支：`main`。
- 阶段 16 执行前 Git：`2b54c2e`；最终运行代码来自 `f9317fa`，后续阶段 16 提交只更新交付文档和工程记忆。
- Backend 容器：`e9261b627b37`，镜像 `sha256:f94f555bbdf74bda2c4b36f04ed83c7dae428045a25eab8698ed51bdfad99d9d`。
- Web 容器：`6a66c52b0882`，镜像 `sha256:3141706386beba3bbe4dedae77f0c4420d4ea20782266fb33f55a149a13fff83`。
- 两个容器均 `running/healthy`、RestartCount=0、OOMKilled=false。
- 数据卷：`open-ai-canvas_backend-data:/data`，未删除、未替换。
- 最终本地镜像标签：
  - `open-ai-canvas-backend:stage16-final-20260827-160505`
  - `open-ai-canvas-web:stage16-final-20260827-160505`
- Web 镜像内 `index.html` 与容器 HTTP 首页 SHA-256 均为 `bf8376eab6bc3727c9fdc2552a79d201dddb90ff3f643e90538c06c17b3c1c3b`。

## 最终质量门禁

| 单元 | 最终结果 |
| --- | --- |
| Go 格式 | `gofmt -l backend` 0 文件 |
| Git 差异 | `git diff --check` 通过 |
| Backend | 隔离 Linux CGO `go test -count=1 ./...` 全部包通过 |
| Web 主套件 | 459/459 |
| Web 跨 Runtime | 1/1 |
| Web TypeScript | 通过 |
| Web 生产构建 | 11,022 模块，通过，无 Rolldown panic |
| Canvas Agent | Windows 主机权限 292/292 |
| Canvas Agent 构建 | TypeScript 通过 |
| Canvas Agent 测试进程 | 0 个残留 |
| 影策插件 | manifest validator 通过；`open-canvas` skill validator 通过 |
| Compose | root/local/dev/deploy/server 五份 `config --no-interpolate` 通过 |
| Server CORS | 空值拒绝；显式 Origin 通过 |
| 凭据启发式扫描 | 无 AWS/GitHub/Bearer 命中；8 个 `sk-` 命中均为 23–28 字符代码/测试片段，不是可用项目密钥 |

Windows 本机没有创建文件符号链接的权限；符号链接拒绝的单一子项继续由未来 Linux CI 保留。遍历、根逃逸、硬链接、超限、文件头、TOCTOU 和路径交换均在本机测试覆盖。

## 数据和 OSS

最终恢复点：`BKP-20260827-160505-STAGE16-FINAL`。

- 数据库完整性：`ok`。
- 外键违规：0。
- 表：60。
- 用户/会话/系统渠道：1/1/1。
- 渠道模型：320，其中 Ready 12、Planned 308。
- 任务/资源/素材：1/1/1。
- 逻辑模型/画布项目：0/0。
- OSS 配置记录和 Secret：1/1；Secret 为 `enc:v1`，匹配 `.settings-key` 可解密且非空，未输出明文。
- 最终备份 ACL 仅当前用户、SYSTEM、Administrators。

恢复必须同时使用数据库、WAL/SHM、匹配 `.settings-key` 和迁移标记，并需要用户即时批准。禁止用 `docker compose down -v` 代替恢复。数据库恢复不能撤销历史供应商费用，也不会自动删除 OSS 对象。

## 运行与 Microsoft Edge 验收

- `/api/health` 返回 `{code:0,status:ok}`。
- 最近 30 分钟 Backend 日志中：`record not found` 0、SQL 脱敏错误 0、阶段测试敏感标记 0。
- 用户以已登录 Microsoft Edge 提供首页、管理后台数据概览、创作页三张截图：
  - 首页导航、项目/画布入口正常；
  - 管理后台显示 1 活跃用户、1 上游请求、1 生成任务、100% 成功率、队列 0；
  - 创作页正常加载历史视频结果和 Wan 3.0 参数摘要，积分余额 98；
  - 未点击生成、未改变模型或价格配置。
- Computer Use 两次因不能可靠确认 Edge URL 而安全终止；没有改用 Chrome。最终 Edge 证据来自用户人工截图。

## 阶段 1–15 已交付能力摘要

- 百炼官方目录 Manifest、Ready/Planned 状态和真实幂等同步。
- Wan 3.0、Wan 2.7、HappyHorse 1.1 模型专属能力与 Provider Adapter。
- Qwen Image 3.0 Pro 专属同步图片合同和 Ready 目录。
- 精确渠道 SKU 计费、系统/前台统一报价和安全路由。
- 一次授权内 Wan 3.0 真实生成、2 积分结算、阿里云 OSS Resource/Asset 和 Edge 播放验证。
- Create 参数按用户/模型/生成方式持久化，声音和水印使用明确 Switch。
- 统一 requestId、错误分类、前端诊断编号和安全结构化日志。
- GORM not-found/SQL 文本收口、Server CORS fail closed。
- Canvas Agent Windows 进程树、签名 Local Runtime 和插件无 Token 深链合同。

## 明确保留的后续事项

1. 在可调用 Codex CLI 的环境重新安装 `yingce@yingce-local`，新建线程验证插件 skill/MCP 与 Microsoft Edge 自动打开。
2. 文档站缺页、Next/Fumadocs 配置和链接修复；用户已决定暂缓。
3. push、PR、远程 GitHub Actions、tag 和 Release；均需单独授权。
4. Qwen Image 3.0 Pro 真实付费生成与 OSS 物化；需单独授权模型、参数、费用和重试次数。
5. Wan 2.7、HappyHorse 的真实账号付费回归。
6. Planned/手工模型逐模型 Adapter；不得因目录存在就开放执行。
7. Create、Canvas、Admin 和逻辑模型尺寸派生纯函数进一步统一。
8. `.claude/settings.local.json` 宽泛权限由用户在 Claude 环境单独收敛。
9. 百炼官方文档更新后重新生成索引并审阅差异。

## 回滚与 NO-GO

- 代码回滚优先使用已标记的本地镜像；不要先恢复数据库。
- 只有确认数据迁移或目录数据错误，且用户再次批准后，才恢复 verified/protected 数据快照。
- 任一完整性失败、外键违规、容器重启、镜像不一致、日志泄密、测试失败或数据计数异常，都阻止远程发布。
- 真实模型调用、OSS 删除、数据库恢复、全局代理改变和远程 Git 写入始终需要新的明确授权。

## 远程发布门禁

阶段 16 按用户批准停在本地边界。下一步若要进入远程发布，必须先确定：

1. 目标 remote（优先用户 Fork，而不是误用本机 `upstream` 临时仓库）；
2. 目标分支和是否允许 push；
3. PR 的 base/head；
4. 是否允许触发 GitHub Actions；
5. 远程 CI 失败时的修复和回滚权限。

在这些事项获得批准前，不执行远程写入。
