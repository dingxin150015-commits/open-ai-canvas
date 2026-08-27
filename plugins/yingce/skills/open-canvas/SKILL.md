---
name: open-canvas
description: 打开影策网页画布并自动连接本地 Canvas Agent。用户要求打开、启动、进入、使用影策或画布时使用。
---

# 打开影策画布

当用户要求打开、启动、进入或使用影策时，不要把 URL 交给用户手动复制，不要通过浏览器点击“新建画布”。优先快速拉起本地画布和本地 Canvas Agent，然后打开只带 `mode` 的画布 URL。网页使用浏览器本地不可导出密钥与 Local Runtime 完成签名握手；禁止读取、打印或把 Agent master token 放入 URL、日志、剪贴板或浏览器存储。

## 默认打开方式

- 新建画布：`<画布网页地址>/canvas?mode=new`
- 最近画布：`<画布网页地址>/canvas?mode=recent`
- 自己选择：`<画布网页地址>/canvas?mode=choose`

默认打开新建本地画布；只有用户明确要求线上地址、最近画布或自己选择时，才改用对应模式。

## 工作流

1. 如果当前仓库是影策项目，优先使用当前仓库的 `web/` 前端。
2. 先检查网页端口归属：如果 `3000`、`3001` 等端口已被占用，必须用当前操作系统的端口/进程工具或服务输出确认监听进程的工作目录属于当前仓库的 `web/`，不能只因为端口存在就当成本地画布。
3. 如果已有当前仓库的 Vite dev 服务，复用它并记录精确 Origin，例如 `http://localhost:3001`。
4. 如果没有当前仓库的服务，在 `web/` 下运行 `bun run dev -- --host 127.0.0.1 --port <端口>`；默认使用 `3000`，被其他项目占用时选择空闲端口。不要使用 Next 命令，不要执行构建或测试。
5. 检查 `127.0.0.1:17371`：只有 `/runtime/info` 返回 `runtime=framefield-local-runtime` 才能复用；若是其他程序占用，不得终止它或冒充影策 Runtime，应停下说明冲突。网页默认只连接这个固定回环端点。
6. 启动 Canvas Agent 时，把第 3/4 步的 localhost 和 127.0.0.1 Origin 作为 `FRAMEFIELD_TRUSTED_WEB_ORIGINS` 的精确逗号分隔值，再运行 `npx -y @ddcat666/open-ai-canvas-agent`。该变量不是画布 URL，也不能使用旧 `CANVAS_URL`。
7. 不读取 `~/.infinite-canvas/canvas-agent.json` 的 token，不请求用户复制 token；`/config` 只可用于确认非敏感状态，不能作为凭据入口。
8. 直接打开最终 URL：`<真实画布地址>/canvas?mode=new`。禁止附加 `agentUrl`、`agentToken` 或其他凭据参数；网页会主动移除这些旧参数。
9. 画布网页会自动新建具体画布、打开本机 Agent 面板并通过签名挑战连接 Local Runtime；不要用浏览器点击新建画布。
10. 打开后再使用 `canvas_get_state` 检查画布是否已经连接；如果尚未连接，先检查 `/runtime/info` 的 `originTrusted`、精确 Origin 和 `FRAMEFIELD_TRUSTED_WEB_ORIGINS`，不要读取 token，也不要改用 Chrome 或线上站点，除非用户明确要求。

## 用户只安装插件时

- 如果当前工作区不是影策源码仓库，优先提示用户先打开或启动影策网页，再连接本地 Agent。
- 可以使用线上画布地址或用户给出的本地地址作为 `<画布网页地址>`，但该精确 Origin 必须由用户信任并进入 `FRAMEFIELD_TRUSTED_WEB_ORIGINS`；最终 URL 仍只带 `mode`。
- 不要假设用户已经安装本仓库依赖；插件的 MCP 会通过 `npx -y @ddcat666/open-ai-canvas-agent mcp` 使用已发布的 Canvas Agent。

不要要求用户手动填写 token 或复制 JSON。外部状态不明确、端口属于其他程序或 Edge 未登录时应停下说明，禁止切换到 Chrome。
