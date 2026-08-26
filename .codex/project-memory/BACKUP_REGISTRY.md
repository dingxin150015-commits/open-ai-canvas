# 影策备份登记簿

本文件是 `<project-root>` 备份、快照和恢复材料的权威登记。登记只保存脱敏路径、组成、哈希、验证状态和维护规则，不保存任何 Secret、Token、Cookie、API Key 或用户数据正文。公开版本使用 `<project-root>`、`<private-backup-root>` 等占位符；本机精确恢复位置保存在 Git 忽略的 `.local/backups/*/backup-pointer.md`。

## 状态定义

- `verified`：文件存在，关键哈希一致，数据库完整性和恢复检查通过。
- `protected`：访问权限已收紧并复核。
- `superseded`：已有更新备份替代，但尚未获准删除。
- `invalid`：缺失、哈希不一致或恢复验证失败，不可用于恢复。

## 保留规则

1. 备份目录使用 `stageN-YYYYMMDD-HHMMSS`，不可覆盖旧目录。
2. 数据库快照必须与匹配的 `.settings-key` 同时登记；缺少任一项即不可恢复加密设置。
3. D 盘为 exFAT，不保存数据库、WAL、密钥或其他敏感恢复材料。
4. 敏感备份只存入支持 ACL 的用户私有目录，并逐文件验证权限主体。
5. 新备份完成后先验证，再将旧备份标记 `superseded`；未经用户明确批准不得删除。
6. 至少保留阶段 0 基线，直到全部开发阶段完成、当前环境稳定运行且存在更新的 verified 备份。
7. 每次数据库 Schema、加密密钥、存储设置、生产部署或数据迁移前后都要新增登记。

---

## BKP-20260825-160625-STAGE0

### 基本信息

| 字段 | 值 |
| --- | --- |
| ID | `BKP-20260825-160625-STAGE0` |
| 创建时间 | 2026-08-25 16:06:25 +08:00 |
| 阶段 | 阶段 0：接管基线保护 |
| Git 基线 | `main@0b202a12aa71c0827548c54deaeaf65455d21686` |
| 状态 | `verified`, `protected` |
| 敏感级别 | 高：包含用户数据库和设置加密密钥 |
| 删除授权 | 未授权；禁止删除 |

### 非敏感工作树基线

位置：

`<project-root>\.local\backups\stage0-20260825-160625`

主要组成：

| 文件 | 用途 | SHA-256 |
| --- | --- | --- |
| `tracked-changes.patch` | 9 个已修改文件的 Git binary patch | `ED335AC3AFBBA1FE3BE47AE5A3C05B65406A6FCEB2549855508999A5839BB37A` |
| `workspace-current-files.zip` | 阶段开始时 44 个已修改/未跟踪文件 | `8148C94C2F99446B6D2A7F2CE7240DB07D07228D240DABD21E20E6EAEE5C83A1` |
| `project-memory-post-stage0.zip` | 阶段 0 完成后的 6 个项目记忆文件 | `86BB84D3B509D6DAAE8F0918907C4746C241F365023E8F423D83E765477C9041` |
| `project-memory-post-registry.zip` | 建立备份登记簿后的 7 个项目记忆文件 | `D85D5EEAE77B6B76AEC5621680BCD0322A7976F49CDE6046F7C7537B8F956B12` |
| `workspace-baseline.json` | 文件清单、大小和逐文件哈希 | 见备份目录 |
| `container-baseline.json` | Backend/Web 容器、镜像、Compose、卷基线 | 见备份目录 |
| `stage0-report.md` | 阶段 0 完整报告 | 见备份目录 |

验证：

- 工作树 ZIP：44/44 条目；CRC 通过；路径集合与基线一致。
- 项目记忆 ZIP：6/6 条目；CRC 通过。
- 登记后项目记忆 ZIP：7/7 条目；CRC 通过。
- `git diff --check`：通过。
- D 盘目录不含数据库、WAL、SHM 或 `.settings-key`。

### 敏感恢复材料

位置：

`<private-backup-root>\stage0-20260825-160625`

主要组成：

| 文件 | 大小 | SHA-256 |
| --- | ---: | --- |
| `database/open_ai_canvas-stage0.db` | 2,445,312 bytes | `30A5296BBF05C571ADC70925F9F7D96A422C3BA2F77E2A39FD4342E726273A36` |
| `restore-material/.settings-key` | 32 bytes | `E4B4B107C3060A18D65AD3AAD4E329D063A924A71CCF23EAB15B9ABF8120A6D5` |
| `restore-material/.storage-v2-migrated` | 20 bytes | `0435B398CEC770046D79B55E9C2E6AB06D236A7759FBB44503266E1DFEA0C8F3` |
| `database-verification.json` | 脱敏验证报告 | `9FDC2D2D6A7A4667E0A903C6B7705E0C2FC96088FCD8B047EF9B54A5377BF1C2` |

ACL 复核：

- `<current-windows-user>`：FullControl，非继承；
- `NT AUTHORITY\SYSTEM`：FullControl，非继承；
- `BUILTIN\Administrators`：FullControl，非继承；
- 无 Everyone/Users 访问项。

### 数据库与恢复验证

| 检查 | 结果 |
| --- | --- |
| `PRAGMA integrity_check` | `ok` |
| 外键违规 | 0 |
| 表数量 | 60 |
| 独立恢复副本完整性 | `ok` |
| 独立恢复副本外键违规 | 0 |
| `.settings-key` 与容器源哈希 | 一致 |
| OSS 密文解密 | 成功；仅验证成功和非空，未输出明文 |
| Backend 备份后状态 | `running`, `Paused=false`, `healthy` |

脱敏计数：1 个用户、1 个会话、1 个系统渠道、241 个渠道模型、0 个逻辑模型、0 个资源、0 个任务、0 个素材、0 个画布项目。

### 恢复依赖与边界

恢复必须同时使用：

1. `database/open_ai_canvas-stage0.db`；
2. `restore-material/.settings-key`；
3. `restore-material/.storage-v2-migrated`；
4. `container-baseline.json` 中记录的运行边界。

恢复不是当前授权事项。任何恢复前必须：

- 再次备份当时 live 数据；
- 停止写入并确认精确目标卷；
- 验证容器 UID/GID 和文件权限；
- 明确回滚步骤；
- 获得用户对具体恢复目标的即时批准。

### 后续维护

- 阶段 1 只修改代码和测试，不改变 live 数据库；阶段结束不自动替代本备份。
- 首次数据库 Schema 迁移前必须创建新的 pre-migration verified 备份并追加登记。
- 阶段 6 更新当前环境前必须创建新的 pre-deploy verified 备份并追加登记。
- 阶段 0 备份至少保留到阶段 8 完成并存在更新的稳定备份。

---

## BKP-20260826-130230-STAGE6-PREDEPLOY

### 基本信息

| 字段 | 值 |
| --- | --- |
| ID | `BKP-20260826-130230-STAGE6-PREDEPLOY` |
| 创建时间 | 2026-08-26 13:02:30 +08:00 |
| 阶段 | 阶段 6：候选部署前保护 |
| Git 基线 | `main@0b202a12aa71c0827548c54deaeaf65455d21686` + 未提交阶段 1-5 候选工作树 |
| 状态 | `verified`, `protected` |
| 敏感级别 | 高：包含用户数据库和设置加密密钥 |
| 删除授权 | 未授权；禁止删除 |

### 敏感恢复材料

位置：

`<private-backup-root>\stage6-20260826-130230`

| 文件 | 大小 | SHA-256 |
| --- | ---: | --- |
| `database/open_ai_canvas-stage6-predeploy.db` | 2,445,312 bytes | `C59211E5F8511654EED2ABC270AB13228F87DB7F4F033A3F1016CC3B00F20938` |
| `restore-material/.settings-key` | 32 bytes | `E4B4B107C3060A18D65AD3AAD4E329D063A924A71CCF23EAB15B9ABF8120A6D5` |
| `restore-material/.storage-v2-migrated` | 20 bytes | `0435B398CEC770046D79B55E9C2E6AB06D236A7759FBB44503266E1DFEA0C8F3` |
| `database-verification.json` | 1,136 bytes | `41AB5F19524C52977F005FB634D3002EE6970D038727E7D77C21F33AAECA9374` |
| `container-baseline.json` | 1,396 bytes | `D5BE1651AE20D3B5B66CBE2F857E445859E0E0C65A811EF35618627D63587F8D` |
| `checksums.json` | 854 bytes | `7BAF6396C2C8F986B0ADF4FBB1F24BCF298119405A407BAF3FBE002FDF6C67E7` |

ACL 复核：目录关闭继承；仅 `<current-windows-user>`、SYSTEM、Administrators 具有 FullControl；Everyone、Users、Authenticated Users 命中 0。

### 数据库与恢复验证

| 检查 | 结果 |
| --- | --- |
| 主快照 `PRAGMA integrity_check` | `ok` |
| 主快照外键违规 | 0 |
| 主快照表数量 | 60 |
| 独立恢复副本完整性 | `ok` |
| 独立恢复副本外键违规 | 0 |
| OSS 配置记录 | 1 |
| OSS 加密 Secret | 1；全部 `enc:v1` |
| OSS 密文解密 | 成功；仅验证非空，未输出明文 |
| Backend 复制后状态 | `running`, `Paused=false`, `healthy` |

脱敏计数：1 个用户、1 个登录会话、1 个系统渠道、241 个渠道模型、0 个逻辑模型、0 个资源、0 个任务、0 个素材、0 个画布项目。

### 部署前容器与候选

- 旧 Backend：容器 `b8415a78e587`，镜像 `sha256:5a2a02e15f28de8f531f4990f46a06ff4bb5ef27bbcfe92256a8834e6f936fd8`。
- 旧 Web：容器 `583a381f2cc3`，镜像 `sha256:cc1d8dae42ec3cf2cf7d39598278ea2c35c3d45afb61c24098db2dccdbf4f37a`。
- 数据卷：`open-ai-canvas_backend-data:/data`。
- Backend 候选：`sha256:0e5744fdfad8767320dadf66a3fb510dd3decbdd335f071b2e64d1dd43200f91`。
- Web 候选：`sha256:0a18498c9634c82fc47d327da8c67ef8e6b62e2e3506e62070c1a53e7e7b0bc6`。

### 恢复边界

恢复必须同时使用数据库快照、匹配 `.settings-key` 和迁移标记。恢复当前未授权；若阶段 6 迁移失败，先使用旧镜像做代码回滚。只有数据回填本身错误且用户再次批准时，才允许停止写入并恢复本备份。禁止删除卷。

---

## BKP-20260826-131009-STAGE6-POSTDEPLOY

### 基本信息

| 字段 | 值 |
| --- | --- |
| ID | `BKP-20260826-131009-STAGE6-POSTDEPLOY` |
| 创建时间 | 2026-08-26 13:10:09 +08:00 |
| 阶段 | 阶段 6：候选部署与 Schema 迁移后保护 |
| Git 基线 | `main@0b202a12aa71c0827548c54deaeaf65455d21686` + 未提交阶段 1-6 候选工作树 |
| 状态 | `verified`, `protected` |
| 敏感级别 | 高：包含迁移后用户数据库和设置加密密钥 |
| 删除授权 | 未授权；禁止删除 |

### 敏感恢复材料

位置：

`<private-backup-root>\stage6-postdeploy-20260826-131009`

| 文件 | 大小 | SHA-256 |
| --- | ---: | --- |
| `database/open_ai_canvas-stage6-postdeploy.db` | 2,555,904 bytes | `05DAE42E0FA56F9100034A34943DF4CCC8247781945B85299951C81D1AA25F72` |
| `restore-material/.settings-key` | 32 bytes | `E4B4B107C3060A18D65AD3AAD4E329D063A924A71CCF23EAB15B9ABF8120A6D5` |
| `restore-material/.storage-v2-migrated` | 20 bytes | `0435B398CEC770046D79B55E9C2E6AB06D236A7759FBB44503266E1DFEA0C8F3` |
| `database-verification.json` | 1,776 bytes | `820B0CD777418D423849BD3CBFA704519D62E2DF38FE66BA98EBC43F69744BB0` |
| `container-baseline.json` | 1,180 bytes | `266369CC9BEAACE8E3CEC036AD60D0B1F16E40DC5DE612966EA80C89B6DCCE86` |
| `checksums.json` | 855 bytes | `971E336907B0C8558CCAD9B46D826CE649D4149BF8A6F21FCF4F7D5C6CC0601A` |

ACL 复核：目录关闭继承；仅 `<current-windows-user>`、SYSTEM、Administrators 具有 FullControl；宽泛主体命中 0。

### 迁移后数据库验证

| 检查 | 结果 |
| --- | --- |
| 主快照 `PRAGMA integrity_check` | `ok` |
| 主快照外键违规 | 0 |
| 主快照表数量 | 60 |
| 独立恢复副本完整性/外键 | `ok` / 0 |
| ChannelModel 目录新增列 | 全部存在 |
| `support_status` | `planned=241` |
| `catalog_source` | `legacy=241` |
| 缺失 Provider Model Key | 0 |
| 缺失操作 JSON / 文档 JSON | 0 / 0 |
| OSS 配置记录 / Secret | 1 / 1 |
| OSS Secret 加密与解密 | 全部 `enc:v1`；可解密且非空，未输出明文 |

业务计数保持为：1 用户、1 会话、1 系统渠道、241 渠道模型、0 逻辑模型、0 资源、0 任务、0 素材、0 画布项目。

### 运行环境

- Backend：容器 `5693d185fab6`，镜像 `sha256:0e5744fdfad8767320dadf66a3fb510dd3decbdd335f071b2e64d1dd43200f91`，healthy。
- Web：容器 `f278caf9beaa`，镜像 `sha256:0a18498c9634c82fc47d327da8c67ef8e6b62e2e3506e62070c1a53e7e7b0bc6`，healthy。
- 数据卷仍为 `open-ai-canvas_backend-data:/data`，未删除、未替换。
- 旧镜像仍由 `pre-stage6-20260826-130230` 标签保护，可执行代码回滚。
- Web 首页原始响应与容器 `index.html` SHA-256 一致：`29B1107BDC54C935A69AF9549A17BD8297513E4C6C01452D20E594326326FE6C`。

### 恢复边界

本备份是阶段 6 完成后的当前稳定恢复点。恢复仍需用户即时批准；不得删除卷或只恢复数据库而遗漏匹配 `.settings-key`。阶段 7 Edge 验证和阶段 8 真实调用如改变数据，需要继续追加备份登记。

---

## BKP-20260826-171800-STAGE7-CATALOG

### 基本信息

| 字段 | 值 |
| --- | --- |
| ID | `BKP-20260826-171800-STAGE7-CATALOG` |
| 创建时间 | 2026-08-26 17:18 +08:00；容器基线采集于 17:20 |
| 阶段 | 阶段 7：百炼目录拉取、幂等修复与管理后台验证后保护 |
| Git 基线 | `main@0b202a12aa71c0827548c54deaeaf65455d21686` + 未提交阶段 1-7 候选工作树 |
| 状态 | `verified`, `protected` |
| 敏感级别 | 高：包含用户数据库和设置加密密钥 |
| 删除授权 | 未授权；禁止删除 |

### 敏感恢复材料

位置：

`<private-backup-root>\stage7-catalog-20260826-171800`

| 文件 | 大小 | SHA-256 |
| --- | ---: | --- |
| `database/open_ai_canvas-stage7-catalog.db` | 2,650,112 bytes | `8A34E57608E949F975A7DB3E089B3BCF13E7BA01D4C4C02C4C7CED9877BE014F` |
| `restore-material/.settings-key` | 32 bytes | `E4B4B107C3060A18D65AD3AAD4E329D063A924A71CCF23EAB15B9ABF8120A6D5` |
| `restore-material/.storage-v2-migrated` | 20 bytes | `0435B398CEC770046D79B55E9C2E6AB06D236A7759FBB44503266E1DFEA0C8F3` |
| `database-verification.json` | 877 bytes | `BA841D5BA8CF85FD154D363CF79EE78F7A9B833EE1A35A2BDF7A78AD02F1C4D2` |
| `catalog-verification.json` | 1,419 bytes | `2784A8031AFC89F9067159363B37FA641C5BDA2CEBC2B4EAE7DD6878B4724259` |
| `container-baseline.json` | 462 bytes | `B099B1110D9B8ADCACA69B96781CD468A54AA6668D033142476F68064B7B30A4` |
| `checksums.json` | 964 bytes | `A33D13F58C5DE4EDDC252ADC2AF3CCC1AD18B280A46B140F20F132C88AA34FF7` |

ACL 复核：目录关闭继承；仅 `<current-windows-user>`、SYSTEM、Administrators 具有 FullControl；宽泛主体命中 0。

### 数据库、目录与恢复验证

| 检查 | 结果 |
| --- | --- |
| 主快照 `PRAGMA integrity_check` | `ok` |
| 主快照外键违规 | 0 |
| 主快照表数量 | 60 |
| 独立恢复副本完整性/外键 | `ok` / 0 |
| OSS 配置记录 / Secret | 1 / 1 |
| OSS Secret 加密与解密 | 全部 `enc:v1`；可解密且非空，未输出明文 |
| 渠道模型总数 | 319 |
| 支持状态 | `ready=11`, `planned=308` |
| 官方补充 | 80：图片 2、视频 78 |
| 目录来源 | `upstream=239`, `bailian-official-docs=78`, `upstream+bailian-official-docs=2` |
| 已启用 / 已定价 | 0 / 0 |
| Ready 合同缺失 / 非 Ready 误启用或误定价 | 全部 0 |
| 上游 239 条保守目录字段偏差 | 全部 0 |
| 幂等摘要 | `997BF93019C1BA6F1E34191BDFCBAF0AA76ABFFF3627370FD9F3F4908C97B7C4` |

脱敏业务计数：1 个用户、1 个登录会话、1 个系统渠道、319 个渠道模型、0 个逻辑模型、0 个资源、0 个任务、0 个素材、0 个画布项目。

### 运行环境与阶段边界

- Backend：容器 `5293ca038048`，镜像 `sha256:6faf8c852dfcab4289afc2f87ea0deb2e2c0ff7927a010a9fa2a8fe40c35023e`，备份后恢复 `running/healthy`、Paused=false、RestartCount=0。
- Web：容器 `f278caf9beaa`，镜像 `sha256:0a18498c9634c82fc47d327da8c67ef8e6b62e2e3506e62070c1a53e7e7b0bc6`。
- 数据卷仍为 `open-ai-canvas_backend-data:/data`，未删除、未替换。
- 管理后台 Edge 已验证 11 个 Ready、Planned 只读、`wan3.0-video` 可用支持状态但未定价且停用；备份登记后用户继续以 Edge 验证普通创作台模型选择器为空，未点击生成。该只读检查未改变备份对应的数据状态。
- 本备份不包含真实百炼生成或 OSS 上传结果；两者仍未授权。

### 恢复边界

这是阶段 7 目录拉取后的恢复点。恢复必须同时使用数据库快照、匹配 `.settings-key` 和迁移标记，并需用户即时批准。禁止删除卷；普通代码回滚优先使用已保护镜像，只有数据目录本身需要回退且再次获批时才恢复本快照。

---

## BKP-20260826-185248-STAGE8-PRECALL

### 基本信息

| 字段 | 值 |
| --- | --- |
| ID | `BKP-20260826-185248-STAGE8-PRECALL` |
| 创建时间 | 2026-08-26 18:52:48 +08:00 |
| 阶段 | 阶段 8：真实 Wan 3.0 调用前保护 |
| Git 基线 | `main@0b202a12aa71c0827548c54deaeaf65455d21686` + 未提交阶段 1-8 候选工作树 |
| 状态 | `verified`, `protected` |
| 敏感级别 | 高：包含用户数据库和设置加密密钥 |
| 删除授权 | 未授权；禁止删除 |

### 敏感恢复材料

位置：

`<private-backup-root>\stage8-precall-20260826-185248`

| 文件 | 大小 | SHA-256 |
| --- | ---: | --- |
| `database/open_ai_canvas-stage8-precall.db` | 2,650,112 bytes | `18297EABC37C28E1B37F58CFAE09299BB151C7A4DA053310D2FEC2C7DC787F49` |
| `restore-material/.settings-key` | 32 bytes | `E4B4B107C3060A18D65AD3AAD4E329D063A924A71CCF23EAB15B9ABF8120A6D5` |
| `restore-material/.storage-v2-migrated` | 20 bytes | `0435B398CEC770046D79B55E9C2E6AB06D236A7759FBB44503266E1DFEA0C8F3` |
| `database-verification.json` | 877 bytes | `BA841D5BA8CF85FD154D363CF79EE78F7A9B833EE1A35A2BDF7A78AD02F1C4D2` |
| `catalog-verification.json` | 1,419 bytes | `2784A8031AFC89F9067159363B37FA641C5BDA2CEBC2B4EAE7DD6878B4724259` |
| `container-baseline.json` | 796 bytes | `974155D462D7168AF856E57BB7A30D4BC3DEA6905BEE00ABCA40388D0C7BD0D6` |
| `checksums.json` | 964 bytes | `F8224A255BEA823B5510522A004DCF1033254B3C0163A401AB1BB24ADF915D86` |

ACL 复核：目录关闭继承；仅 `<current-windows-user>`、SYSTEM、Administrators 具有 FullControl；宽泛主体命中 0。

### 数据库与恢复验证

| 检查 | 结果 |
| --- | --- |
| 主快照 `PRAGMA integrity_check` | `ok` |
| 主快照外键违规 | 0 |
| 主快照表数量 | 60 |
| 独立恢复副本完整性/外键 | `ok` / 0 |
| OSS 配置记录 / Secret | 1 / 1；全部加密、可解密，未输出明文 |
| 渠道目录 | 319：`ready=11`, `planned=308` |
| 目录摘要 | `997BF93019C1BA6F1E34191BDFCBAF0AA76ABFFF3627370FD9F3F4908C97B7C4` |
| 已启用 / 已定价模型 | 0 / 0 |
| 任务 / 资源 / 素材 | 0 / 0 / 0 |

### 运行环境与调用授权

- Backend：容器 `c33b097bd7a2`，镜像 `sha256:5b236e5e8b6980b8d9eb43e901927753f1299c46de03db9a9d162b6e5c9505bc`，healthy、Paused=false、RestartCount=0。
- Web：容器 `b27398fc1af4`，镜像 `sha256:fa6f2eaa11e8e30ef8fe4b232fdc2831a141077a97968fe7464e0873e58df5f6`，healthy、RestartCount=0。
- 数据卷仍为 `open-ai-canvas_backend-data:/data`，未删除、未替换。
- 阶段 8 候选已修复 OpenAI 兼容 Base URL 到原生 DashScope 路径、结果下载跨主机凭据泄露和 `generation-*` 素材未进入 OSS 同步的问题。
- 创建本备份时尚未配置 Wan 3.0 价格、未启用模型、未提交真实任务；`paidCallAuthorized=false`。
- 本备份创建后用户批准了精确调用；随后发生并已修复 SQLite WAL/SHM owner 漂移事故。修复后脱敏计数仍与本备份一致，未创建任务/资源/素材/账务；当前精确计费 Backend 镜像另由 `stage8-paid-ready-20260826` 标签保护。

### 恢复边界

这是阶段 8 真实调用前的权威恢复点。真实调用可能产生模型费用、任务、账务、资源、素材和 OSS 对象；数据库恢复不能撤销供应商账单或自动删除 OSS 对象。任何恢复、对象删除或卷操作仍需用户即时批准。

---

## BKP-20260826-201444-STAGE8-POSTCALL

### 基本信息

| 字段 | 值 |
| --- | --- |
| ID | `BKP-20260826-201444-STAGE8-POSTCALL` |
| 创建时间 | 2026-08-26 20:14:44 +08:00 |
| 阶段 | 阶段 8：真实 Wan 3.0 与 OSS 验证后保护 |
| Git 基线 | `main@0b202a12aa71c0827548c54deaeaf65455d21686` + 未提交阶段 1-8 最终工作树 |
| 状态 | `verified`, `protected` |
| 敏感级别 | 高：包含用户数据库、调用审计和设置加密密钥 |
| 删除授权 | 未授权；禁止删除 |

### 敏感恢复材料

位置：

`<private-backup-root>\stage8-postcall-20260826-201444`

| 文件 | 大小 | SHA-256 |
| --- | ---: | --- |
| `database/open_ai_canvas-stage8-postcall.db` | 2,650,112 bytes | `D235C46B020F326F2047E647ECE47030C39A9D7D4DFBDEC25A4534BDCF93098A` |
| `restore-material/.settings-key` | 32 bytes | `E4B4B107C3060A18D65AD3AAD4E329D063A924A71CCF23EAB15B9ABF8120A6D5` |
| `restore-material/.storage-v2-migrated` | 20 bytes | `0435B398CEC770046D79B55E9C2E6AB06D236A7759FBB44503266E1DFEA0C8F3` |
| `database-verification.json` | 877 bytes | `4EF553619ABE5D5653A3AA6F08C0C1D8CE9D228FD8F14EAFD8EA794358EF4763` |
| `catalog-verification.json` | 1,419 bytes | `A18986B94A00D60628329D0699A0BF3B6AE2CA931E577F3AAFBCED7F7F241C5C` |
| `runtime-verification.json` | 3,797 bytes | `386EF4AD3541B440242C8ACF5BE123B75139461BA6DA49E6B1F551121860CEE4` |
| `container-baseline.json` | 840 bytes | `0B0FF542E45D57814F99FF8515CCC210C3138C69084A1C5F10F4EA4DCB12F3BA` |
| `checksums.json` | 1,120 bytes | `0E8527EC7305304DB2018E9FB5FA18FF8F01CB3095F2808B1ACA8778703E6AE4` |

ACL 复核：目录关闭继承；仅 `<current-windows-user>`、SYSTEM、Administrators 具有 FullControl；宽泛主体命中 0。

### 数据库、调用和恢复验证

| 检查 | 结果 |
| --- | --- |
| 主快照 `PRAGMA integrity_check` | `ok` |
| 主快照外键违规 | 0 |
| 主快照表数量 | 60 |
| 独立恢复副本完整性/外键 | `ok` / 0 |
| OSS 配置记录 / Secret | 1 / 1；全部加密、可解密，未输出明文 |
| 渠道目录 | 319：`ready=11`, `planned=308` |
| Wan 3.0 安全状态 | 模型 disabled/unpriced；唯一临时价格档 disabled/unpriced、单价 0 |
| 真实创建请求 | 1；HTTP 200；任务 succeeded；attempts=1 |
| 实际 Provider 参数 | `adaptive`, 480P, 2 秒, audio=false, watermark=false, seed 省略 |
| 内部账务 | settled；2 秒；实际 2 积分；预留 0；余额 98 |
| OSS Resource | 1；ready；Aliyun；video/mp4；805,526 bytes；乌兰察布 Endpoint |
| Asset | 1；confirmed；`resource:` 存储键 |
| Edge 视觉验证 | 正常播放；0:02；无声；无水印 |

### 运行环境与事故记录

- Backend：容器 `f97220ae4273`，镜像 `sha256:0d3b2403dbcf106f8dd61dcb3cd88692e2858c17a591b2aeec8d3aca6a08cdca`，healthy、Paused=false、RestartCount=0。
- Web：容器 `7b71d615ca98`，镜像 `sha256:3900fbb49fceda5da019be3991b818d36c2e4f6831c688f4d1adaaeab1e7daa2`，healthy、RestartCount=0。
- 数据卷仍为 `open-ai-canvas_backend-data:/data`，未删除、未替换。
- 阶段中 root 监控容器曾把 WAL/SHM 改为 root:root 0644，导致 Backend readonly 重启；付费任务前已恢复 DB/WAL/SHM 为 100:101/0660，并把后续监控固定为相同 UID/GID。事故未产生任务或费用，数据计数保持不变。
- Create 初始缺少音频/水印入口；已补齐能力驱动输出控制，并通过 Web 449+1、TypeScript 和容器构建验证。

### 恢复边界

这是阶段 8 完成后的权威数据库恢复点。恢复数据库不能撤销已经发生的百炼供应商费用，也不会自动删除或恢复阿里云 OSS 对象；OSS 对象与数据库资源记录必须按同一审计范围处理。任何恢复、对象删除、价格重新启用或卷操作仍需用户即时批准。
