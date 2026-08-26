# 阶段 5 发布候选登记

更新时间：2026-08-26

## 候选身份

- Git 基线：`main@0b202a1`；当前候选仍是未提交工作树，不能用单一 commit 表示。
- 百炼目录：版本 `2026-08-25`，80 项；11 个 Ready 视频、67 个 Planned 视频、2 个 Planned 图片。
- Manifest SHA-256：`D52B4560F25C72C55C5566C6DB8980B60441E04CC2E648A87E8858DE6E0EFB62`；重新生成前后完全一致。
- Backend 候选镜像：`open-ai-canvas-backend:stage5-20260826`，ID `sha256:0e5744fdfad8767320dadf66a3fb510dd3decbdd335f071b2e64d1dd43200f91`。
- Web 候选镜像：`open-ai-canvas-web:stage5-20260826`，ID `sha256:0a18498c9634c82fc47d327da8c67ef8e6b62e2e3506e62070c1a53e7e7b0bc6`。
- 候选镜像没有替换 `:local` 标签，也没有重建当前运行容器。

## 阶段 5 收敛项

- 修复 Bun test 内嵌 Vite/Rolldown 的退出竞争，增加 panic guard，并统一 Bun 1.4.0。
- 删除 Qwen 图片能力与 Canvas 设置面板的临时 DEBUG 日志。
- 明确图片尺寸三态：缺失继承默认；显式空列表不回填；`*` 只开启自定义输入，不制造预设。
- Catalog 与系统任务准入共用持久能力配置；缺失、损坏或不完整时 fail closed。
- 旧记录只有能力、协议和有效能力 JSON 全部存在时才回填为 `ready/manual`；其余为 `planned/legacy`。
- 路由和公开目录都要求 `price_configured=true`，残留正价格的草稿不会被误判为可用。
- DashScope 视频日志不再输出负向提示词、媒体 URL 前缀、API Key 片段、完整 URL 或完整响应 JSON。
- 变更与新增候选源码进行了敏感凭据模式扫描，命中 0。

## 验证证据

- Web：446/446 主套件、1/1 跨 Runtime；TypeScript 通过；Vite 生产构建通过且无 Rolldown panic。
- Canvas Agent：使用仓库固定的 Bun 1.4.0 路径，286/286 全量测试和 TypeScript 构建通过。
- Rolldown 专项：直接导入测试连续 10 次，0 失败、0 panic；合成 guard fixture 验证 clean=0、panic=86。
- Backend：隔离 Linux CGO `go test -count=1 ./...` 全部通过；Database、Repository、Handler、Provider、Service 均通过。
- Backend 候选：临时、无挂载容器启动并通过 `/api/health`，退出后自动删除。
- Web 候选：静态入口存在；提供临时 `backend` hosts 后 `nginx -t` 通过，容器退出后自动删除。
- 当前运行容器仍为 `open-ai-canvas-web-1@583a381f2cc3` 和 `open-ai-canvas-backend-1@b8415a78e587`，均使用旧 `:local` 镜像且保持 healthy。

## 阶段 6 迁移顺序

1. 按 `BACKUP_REGISTRY.md` 新建并验证 pre-deploy SQLite/WAL/SHM/`.settings-key` 备份，登记 ACL、哈希、完整性、外键和恢复结果。
2. 给当前运行的两个旧镜像 ID 增加明确的 pre-stage6 回滚标签；记录容器 ID、镜像 ID、卷和 Compose 配置哈希。
3. 先将 Backend 候选切换为 `:local`，只重新创建 Backend，不执行 `down`，绝不删除卷；等待 schema migration 和 health 完成。
4. 只读核对新增列、旧记录 ready/planned 回填数量、用户/渠道/模型/价格/资源/任务计数，以及 OSS 密文仍可解密。
5. 再将 Web 候选切换为 `:local`，只重新创建 Web；核对 HTML 与主要 chunk 哈希。
6. 阶段 7 才使用已登录 Microsoft Edge 进行管理后台拉取、分类/状态/幂等和无费用任务前置验证。

## 回滚边界

- 代码回滚：把 pre-stage6 镜像重新标记为 `:local` 并逐个重新创建服务；新增数据库列为加法变化，旧代码可忽略。
- 数据回滚：若 backfill 或目录同步结果错误，必须经用户再次批准后从阶段 6 pre-deploy 备份恢复数据库和匹配 `.settings-key`；不得以删除卷代替恢复。
- 不使用 `docker compose down -v`、不清理当前卷、不在失败时盲目重复迁移。

## NO-GO

- 当前运行环境尚未应用候选代码和 schema，不能声称管理后台已出现 80 项补充目录或 11 个 Ready 模型。
- 阶段 6 前的新备份与恢复验证是硬门禁；阶段 0 备份不能替代部署前新快照。
- 当前候选未提交、未创建 PR、远程 GitHub Actions 未运行；如需 Git/PR，应另行授权。
- Edge、真实百炼、OSS 上传和任何计费生成均未执行；真实调用继续要求单独确认账号、区域、模型、参数、费用和重试次数。

## 阶段 6 执行结果

- 候选镜像已按计划逐个切换为 `:local` 并重新创建 Backend/Web；两个服务 healthy。
- 数据卷保持 `open-ai-canvas_backend-data`，旧镜像回滚标签仍存在。
- pre/post-deploy 备份均已在 `BACKUP_REGISTRY.md` 登记为 `verified/protected`。
- 迁移后 241 条旧模型为 `planned/legacy`，等待阶段 7 Edge 拉取官方补充目录。
