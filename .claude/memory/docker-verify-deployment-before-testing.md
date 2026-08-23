---
name: docker-verify-deployment-before-testing
description: 改代码后必须验证容器跑的是新二进制，禁用 docker-compose restart
metadata:
  type: feedback
---

修改代码并重建镜像后，**必须先验证运行中的容器确实加载了新代码，再开始任何测试或诊断**。

**核心纪律：**

1. **禁用 `docker-compose restart <service>`** — 它只重启现有容器，**不会**换用新构建的镜像。必须用 `docker-compose up -d <service>`（会 Recreate 容器）。

2. **重启后强制验证二进制**，任选一种：
   - 扫描新增标记：`docker exec <container> sh -c "strings /path/to/binary | grep -c '<新日志标记>'"` — 必须非 0
   - 对比镜像 ID：`docker inspect <container> --format '{{.Image}}'` 与 `docker images <image>:tag --format '{{.ID}}'` 必须一致
   - 看重启输出：出现 `Recreated` 才是换了镜像；只有 `Restarting/Started` 说明用的还是旧容器

3. **验证通过前，任何"日志没输出""功能没生效"的观察都不可作为证据。**

**Why:** 2026-08-22 排查 DashScope 视频生成失败时，我加了大量诊断日志后用 `docker-compose restart backend` 重启，容器实际仍跑 10 小时前的旧二进制（镜像 ID `094312be` vs 新构建 `a1d4b402`，`strings | grep -c 'DashScope Video'` 返回 0）。因"日志一条都没有"而错误推断"函数从未被调用"，进而怀疑路由/配置解析，整条推理链全部作废，浪费了一轮完整的构建+测试+分析周期。

**How to apply:**
- 改代码 → `docker-compose build` → `docker-compose up -d` → **验证二进制** → 才开始测试
- 观察到"新代码像是没生效"时，第一反应先查部署，而不是先怀疑业务逻辑
- 同理适用于前端：容器里是打包产物，改了源码不重建镜像等于没改
- 参见 [[logging-must-cover-full-call-chain]]、[[api-development-best-practices]]
