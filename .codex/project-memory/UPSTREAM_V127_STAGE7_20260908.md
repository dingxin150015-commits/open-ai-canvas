# v1.2.7阶段7执行记录

## 授权与边界

- 用户明确批准阶段7，并决定Director F6与其他视觉项统一留到阶段8，由用户手工Microsoft Edge验收。F6仍是待验收项。
- 本阶段只允许推送 `origin/codex/upstream-v1.2.7-20260908`、独立核对远程SHA并检查远程CI状态。
- 不更新 `origin/codex/dingxin-stable`、`origin/main`、GitHub默认分支，不创建tag/Release/PR，不触发部署或工作流手工运行。
- 当前候选起点 `a23cad205663bf3fa864164dcafbbb3544608e36`，源码检查点 `8808750199549523cc699087db10ca6be828901f`。推送前会固化本决定，最终SHA以后续记录为准。

## 执行中

- 已核对候选工作树干净，origin为 `https://github.com/dingxin150015-commits/open-ai-canvas.git`，官方remote push保持禁用。
- 后续先核对GitHub稳定分支和候选目标，再执行dry-run与非强制push；结果完成后增量更新本文件。
