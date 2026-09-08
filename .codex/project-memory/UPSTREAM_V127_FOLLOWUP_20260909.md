# v1.2.7 阶段8复盘整改与部署记录

## 当前状态

- 运行状态：`runtime_followup_healthy_edge_pending`。修复版源码提交 `ae01c5d667a979c835da03bca21bc4d418382e4e` 已部署；Web/Backend均healthy、RestartCount=0，显示版本`v1.2.7+dingxin.1`。
- 当前运行镜像：Backend `sha256:03b068b2782b4533688d574b12f273e2dc7ea7b016cd3e9fc2d40ff782b3b012`，Web `sha256:f736a584ed204db1087cbc606d57c660754f150d3d270725686c2b1d102d36d6`。
- 本文件将形成纯文档提交，因此候选HEAD会比运行源码多一个记忆提交；业务源码和运行构建均固定在`ae01c5d6`。未推送修复提交，远程候选仍为`aff2c8e5`；stable/main未变。

## 已整改事项

1. D盘三个审计子目录中的15个数据库/密钥副本已迁入C盘私有备份集；逐文件SHA256一致，ACL仅当前用户、SYSTEM、Administrators，宽泛主体0。D盘只保留脱敏JSON和构建日志。
2. `ImportAdminChannelModels`在选择导入时同时处理新模型与安全补齐：只补未启用、未定价、无能力配置的已有记录，并以`updated_at`乐观条件保护并发管理员写入；第二次同步零更新。UI显示新增/补齐结果，已有配置不覆盖。
3. Backend生产镜像和Linux CGO测试层均安装ffmpeg/ffprobe；运行容器已确认两者可执行，MPEG-4回填专项不再因缺依赖跳过。
4. Web quality workflow改为`bun run test:all`，动态发现`web/test`与`web/src`下全部测试文件；当前173文件全部进入门禁。
5. 重新从verified Schema6恢复点创建全新无网络克隆卷，以与真实环境一致的`ENABLE_PROVIDER_PLUGINS=true`和并发参数连续启动两次。验证完成后临时卷已删除。

## 验证证据

- Backend Linux CGO+ffmpeg全量：第二轮全部包通过。第一轮仅既有资源删除并发测试出现一次SQLite表锁，原样重跑未复现；没有修改业务删除逻辑或放宽断言。
- 专项：模型目录首次安全补齐`updated=1`、第二次零更新；MPEG-4播放回填终态测试通过且未跳过。
- Web：TypeScript通过；1619项主集合+5项隔离注册表测试，合计1624通过/0失败；生产构建通过；本轮文件Prettier检查通过。
- 配置一致克隆：Schema6→8；两次启动healthy。业务计数、八组关键字段摘要、迁移记录、密钥/标记摘要前后一致；完整性ok、外键0、重复活动逻辑code0；2个密文均可解密，失败0。
- 新部署前恢复点`BKP-20260909-005845-V127-FOLLOWUP`为verified/protected：完整数据、DB/WAL/SHM、匹配设置密钥、迁移标记、plugin registry和源码bundle；Backend暂停3.84秒后恢复healthy，逐文件ACL宽泛主体0。
- 真实部署前后业务计数、关键字段摘要、迁移记录、密钥/标记一致；2个密文均可解密。部署后Backend暂停1.25秒完成审计并恢复healthy，最近错误日志匹配0。

## Edge与后续边界

尚未执行Microsoft Edge验收。必须继续检查F6、首页/创作入口、画布顶底栏在侧栏/Agent/DevTools组合下的响应式布局、明暗主题与网格、自定义外观、Runtime有界探测、报价400和偶发加载失败。模型目录真实拉取会访问上游并写数据库，不纳入本轮无外部写入的视觉验收；其源码/事务行为已由离线服务测试覆盖。

通过用户Edge验收前，不晋升stable，不操作Host Updater，不调用Provider或付费生成。
