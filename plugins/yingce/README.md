# 影策 Codex 插件

这个插件把影策的本地 Canvas Agent MCP 打包给 Codex app 使用，让 Codex 能打开本地画布、读取当前节点、创建内容并触发生成流程。

## 安装

> 影策尚未上架 Codex 公共插件目录，直接搜索不会显示。请从本仓库自带的 marketplace 安装。

### AI 自动安装

把下面这段发给 Codex：

```text
请从 https://github.com/ddcat-ai/open-ai-canvas.git 安装影策 Codex 插件。
请 clone 仓库到 ~/plugins/open-ai-canvas，确认 .agents/plugins/marketplace.json 和
plugins/yingce/.codex-plugin/plugin.json 都存在。然后运行
codex plugin marketplace add ~/plugins/open-ai-canvas，
再运行 codex plugin add yingce@yingce-local。
安装后请校验插件，并告诉我是否需要开启一个新对话来加载新技能和 MCP 工具。
```

### 手动安装

如果本机还没有仓库，先 clone：

```bash
mkdir -p ~/plugins
git clone https://github.com/ddcat-ai/open-ai-canvas.git ~/plugins/open-ai-canvas
```

注册仓库 marketplace 并安装插件；如果使用已有仓库，请把路径替换为仓库的绝对路径：

```bash
codex plugin marketplace add ~/plugins/open-ai-canvas
codex plugin add yingce@yingce-local
```

安装后建议开启一个新的 Codex 对话，让新的 skill 和 MCP 工具完整加载。

### 本仓库开发调试

如果你就在影策仓库中调试插件，可以直接添加当前仓库。建议使用仓库绝对路径，避免 Codex 从其他工作目录解析失败：

```bash
cd /path/to/open-ai-canvas
codex plugin marketplace add "$(pwd)"
codex plugin add yingce@yingce-local
```

## 使用

1. 新建 Codex 线程后说“打开影策”。
2. 插件会确认当前仓库的本地画布服务是否已运行；端口被占用时会检查进程归属，不会把其他项目的 `3000` 当作影策。
3. 确认或启动后，插件会打开只带 `mode` 的新建画布 URL；网页使用浏览器本地密钥与 Local Runtime 完成签名握手，不把 Agent token 放入 URL。
4. 画布打开后，让 Codex 读取或操作当前画布。

常用提示：

```text
打开影策
读取当前画布并总结节点结构
根据选中节点创建一组生图提示词
```

## 工作机制

插件默认通过以下命令启动 MCP。MCP 进程只注册工具并连接已有的 Local Runtime；`open-canvas` 技能负责检查并启动网页和本地 Agent，不能把 MCP 进程误当成 HTTP Runtime：

```bash
npx -y @ddcat666/open-ai-canvas-agent mcp
```

外部 MCP 工具超时为 2160 秒，覆盖 Canvas Agent 的 35 分钟生成续接窗口并保留一分钟响应/清理余量。修改本地插件源后需要重新执行 `codex plugin add yingce@yingce-local` 并新建对话，已运行线程不会热加载新的 skill 或 MCP 配置。

## 手动排查

优先本地启动画布：

```bash
cd web
bun install
bun run dev
```

然后启动本地 Agent。必须把实际 Vite Origin 作为精确可信来源传入；这里的变量只声明允许握手的网页 Origin，不包含 token：

```powershell
$env:FRAMEFIELD_TRUSTED_WEB_ORIGINS = "http://localhost:3000,http://127.0.0.1:3000"
npx -y @ddcat666/open-ai-canvas-agent
```

如果 `3000` 被当前仓库以外的进程占用，使用 Vite 的端口参数启动空闲端口，例如：

```powershell
bun run dev -- --host 127.0.0.1 --port 3001
$env:FRAMEFIELD_TRUSTED_WEB_ORIGINS = "http://localhost:3001,http://127.0.0.1:3001"
npx -y @ddcat666/open-ai-canvas-agent
```

手动排查时访问 `http://127.0.0.1:17371/runtime/info`，确认返回 `framefield-local-runtime` 且当前网页 Origin 已获信任；`/config` 只返回非敏感状态，不提供 token。随后打开 `<画布网页地址>/canvas?mode=new`。旧的 `agentUrl`、`agentToken` 查询参数会被网页主动移除并拒绝，不能再用于连接。

## 技能库

插件的 `skills/` 是可按需加载的 Codex 技能库，不需要把整套规则复制到每次对话里。安装插件后，建议新建一个 Codex 线程；Codex 会根据请求自动发现并加载对应技能：

- `canvas-context`：先读语义化画布上下文、选区、连接和资源就绪状态。
- `canvas-editing`：写入前校验真实节点 id，批量操作后复核结果。
- `asset-aware-generation`：复用已有角色、场景、道具、风格和媒体资源创建生成流程。

安装/更新流程：

```bash
cd /path/to/open-ai-canvas
codex plugin marketplace add "$(pwd)"
codex plugin add yingce@yingce-local
# 更新插件后开启一个新的 Codex 对话
```

验证技能和 MCP 是否已加载，可以在新对话中直接说：

```text
读取当前画布上下文，列出可用媒体资源，并说明哪些资源可以作为生图参考。
```

如果回答没有调用 `canvas_get_context` / `canvas_get_resources`，先确认当前对话是安装插件后新建的，并检查 `codex mcp list` 中的 `yingce`。
