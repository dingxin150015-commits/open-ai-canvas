# 百炼声明式插件完整迁移

## 基线与授权

- 基线为已晋升的 `d9fec0f4b6f2cec48ca11c8f69bd5c4e7c8b2f48`，修复分支 `codex/dashscope-system-plugins-20260911`，独立工作树 `.local/worktrees/dashscope-plugins`。
- 用户明确批准完整迁移、必要的通用引擎扩展、系统性预验证和部署，完成后停在 Microsoft Edge 门禁前。没有批准真实付费生成或远程发布。
- 主工作区仍在旧 `b6254b41`；此前宿主直接注册补丁及其测试不能作为当前基线的发布证据。受控宿主插件过渡实现已从新分支移除。

## 当前实现

- 新增正式 `dashscope-image` / `dashscope-video` 插件，使用统一官方包 bootstrap、Registry、目录及启停流程，未放宽上传包 host 绑定限制。
- 模型线协议迁入 `plugin-packages/bailian-specs.mjs`；生产不再调用旧 DashScope Provider。旧文件改为 `_legacy_test.go`，仅做原始请求对照。
- 引擎增加纯变量、严格数值/布尔解析、通用尺寸及媒体规则、帧选择、数组与字符汇总、区域路径后缀处理、严格响应状态和匿名结果下载。
- 插件整包 Provider 注册具有原子性；停用影响目录、报价和新任务准入。已知新模型不允许保存为旧 Wanx/Wan 协议。
- 当前产品能力画像保留 UI/报价/准入职责；完全插件化产品能力目录属于后续独立范围，不能把线协议迁移等同于删除全部模型知识。

## 预验证进度

- 已完成 12 个 Ready 模型及关键操作的请求对照夹具、非法输入、启停/重启、保存读回、同步/异步本地 HTTP 模拟与恢复不重复创建测试。
- 首轮问题：重复变量导致包加载失败；缺省空字符串被 coalesce 变为 nil；视频续写素材顺序不同；通用下载同源带鉴权。均已定位并修正，加入对照或安全断言。
- 首轮 Linux CGO 全量只有既有 `TestDeleteAdminResourcesSupportsPartialSuccessAndBatchedReferences` 出现 SQLite 表锁；本次新增测试通过。最终代码仍须重新全量验证，不提前宣称最终通过。
- 部署、克隆数据校验、运行镜像及 Edge 状态将在完成后追加。
