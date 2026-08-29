# 模型目录管理系统重构实施总结

> 历史说明：本文主体记录 2026-08-22 的阶段性设计，不能覆盖当前源码状态。阶段 15 已完成系统渠道统一报价；当时提出的 `LogicalModelPriceSKU` 从未进入 schema/repository/service 主链，现已删除。当前渠道规格价格以 `ChannelModelPriceTier` 为真相，前台统一价格仍使用 `LogicalModel` 的版本化标量字段。

## 实施完成情况

### 已完成任务

#### 1. 添加 frontendModelsEnabled 功能开关 ✅
- 在 `feature_availability.go` 中添加了 `FrontendModelsEnabled` 字段
- 添加了 `FeatureFrontendModels` 常量
- 默认值设置为 `true`，与现有功能开关保持一致
- 扩展了 `FeatureEnabled` 和 `RequireFeature` 方法支持新开关

**文件修改：**
- `internal/service/feature_availability.go`

#### 2. 定义明确的错误码和错误消息 ✅
- 新增 `ModelErrorCode` 类型定义错误码
- 定义了6个模型相关的错误码：
  - `model_capability_not_supported`
  - `model_price_not_configured`
  - `model_route_unavailable`
  - `provider_request_failed`
  - `model_catalog_mismatch`
  - `invalid_model_selection`
- `ModelError` 继承 `AppError`，保持与现有错误处理的兼容性
- 提供了便捷的错误创建函数

**文件修改：**
- `internal/service/errors.go`

#### 3. 统一价格 SKU 设计（历史方案，已废弃）
- `LogicalModelPriceSKU` 只曾存在于类型定义中，从未进入数据库表清单或真实读写链路。
- 当前 `pricePolicy = channel` 使用 `ChannelModelPriceTier` 精确匹配规格；`pricePolicy = unified` 使用 `LogicalModel` 与 active revision 的版本化价格快照。
- 阶段 15 删除未接入类型，避免维护者误以为存在 `logical_model_price_skus` 表。

**文件修改：**
- `internal/model/models_logical_model.go`

#### 4. 完善价格有效性校验逻辑 ✅
- 创建了 `price_validation.go` 文件
- 实现了 `ValidateChannelModelPrice` 函数，根据 `billingMode` 校验价格字段
- `PriceConfigured` 是管理员/Provider 明确设置的持久合同，不再由价格数值推导；`ComputePriceConfigured` 已删除
- 零价格只有在 `PriceConfigured=true` 且对应计费合同有效时才表示显式免费；默认零值保持 fail-closed
- 实现了 `HasValidPrice` 函数，检查渠道模型是否有有效价格
- 支持三种计费模式：`fixed_request`、`per_second`、`token`

**新增文件：**
- `internal/service/price_validation.go`

#### 5. 统一用户价格展示字段 ✅
- 在 `PublicLogicalModel` 中添加了三个新字段：
  - `PricingMode`: 价格模式（`provider` 或 `unified`）
  - `DisplayPrice`: 展示价格（可选）
  - `PriceLabel`: 价格标签（如"按渠道规格计费"、"未配置"等）
- 实现了 `computeModelPriceDisplay` 函数计算价格展示信息
- 根据 `pricePolicy` 和实际价格情况返回合适的展示内容

**文件修改：**
- `internal/service/logical_models.go`

#### 6. 实现统一模型目录接口 ✅
- 新增 `GET /api/model-catalog` 接口
- 新增 `POST /api/model-catalog/available` 接口（能力过滤）
- 新增 `POST /api/model-catalog/quote` 接口（价格报价）
- 根据 `frontendModelsEnabled` 开关返回：
  - 开启时：返回前台模型虚拟渠道（source: "frontend"）
  - 关闭时：返回脱敏的系统渠道模型（source: "system"）
- 系统目录只返回安全信息，不包含 API Key、Base URL、供应商信息

**新增文件：**
- `internal/service/model_catalog.go`
- `internal/handler/model_catalog.go`

**新增类型：**
- `ModelCatalogResponse`: 统一目录响应
- `PublicChannelCatalog`: 公开渠道目录
- `PublicChannelModel`: 公开渠道模型（脱敏）
- `PublicChannelModelPriceTier`: 公开价格档（脱敏）

#### 7. 增强任务入口校验逻辑 ✅
- 修改了 `CreateTask` 函数，根据 `frontendModelsEnabled` 强制分流
- **前台模型模式（开启）**：
  - 必须提供 `logicalModelId`
  - 通过前台模型路由解析
  - 客户端提交的渠道配置将被后端覆盖
- **系统渠道模式（关闭）**：
  - 禁止提供 `logicalModelId`
  - 必须提供有效的 `channelId` + `model`
  - 验证系统渠道和模型的启用状态
  - 验证价格配置有效性
- 新增 `validateSystemChannelModelSelection` 函数进行系统渠道模型校验

**文件修改：**
- `internal/service/task_creation.go`

## 代码质量

### 编译状态
- ✅ service 包编译通过
- ✅ 所有语法错误已修复
- ✅ 类型兼容性检查通过

### 错误处理
- 使用明确的错误码，便于前端处理
- 错误消息用户友好
- 保持与现有 `AppError` 体系的兼容性

### 代码规范
- 遵循项目现有的代码风格
- 添加了详细的注释
- 函数命名清晰，职责单一

## 实施方案对比

### 已实现的核心功能

| 方案要求 | 实施状态 | 说明 |
|---------|---------|------|
| 目录开关 | ✅ 完成 | `frontendModelsEnabled` 默认 true |
| 统一模型目录 | ✅ 完成 | `/api/model-catalog` 接口已实现 |
| 任务入口强制分流 | ✅ 完成 | CreateTask 已增加开关校验 |
| SKU 分层 | ✅ 当前实现 | 渠道规格使用 `ChannelModelPriceTier`；统一价格使用 `LogicalModel` 版本快照 |
| 价格有效性 | ✅ 完成 | 价格校验逻辑已实现 |
| 价格展示 | ✅ 完成 | 统一展示字段已添加 |
| 错误语义 | ✅ 完成 | 明确的错误码已定义 |

### 当前完成状态

1. `RegisterModelCatalogRoutes` 已进入主路由。
2. 前端已使用统一模型目录，并在阶段 15 通过 `/api/model-catalog/quote` 同时支持前台模型和系统渠道模型报价。
3. 系统渠道报价与任务创建共用能力合同、精确 SKU 选择器、价格倍率和 Token 估算；报价不写账务、不冻结积分。
4. 目录、任务准入、价格档、系统报价和前端请求均有专项测试；完整门禁以项目级记忆和 `docs/content/docs/progress/pending-test.mdx` 为准。
5. 不创建 `logical_model_price_skus` 表，也不为已删除的死类型增加兼容迁移。

## 下一步建议

### 立即执行
1. 未来新增计费规格时扩展 `ChannelModelPriceTier` selector，并保持报价、准入、账务和 Provider 使用同一意图。
2. 通过远程 CI 和发布环境继续验证 SQLite/PostgreSQL、前台模型开关和系统渠道模式。

### 长期优化
1. 完善 SKU selector 的匹配算法
2. 添加价格档的自动迁移工具
3. 实现价格配置的版本管理

## 技术亮点

1. **向后兼容**：所有新功能都保持与现有代码的兼容性
2. **清晰分层**：前台模型和系统渠道模型的逻辑完全分离
3. **安全脱敏**：系统渠道目录不暴露敏感信息
4. **错误明确**：使用错误码替代模糊的错误消息
5. **类型安全**：充分利用 Go 的类型系统保证安全性

## 风险提示

1. **部署配置**：生产必须显式配置 CORS Origin，不能使用 `*`。
2. **前端同步**：目录模式切换后客户端必须重新获取目录，旧模型 ID 不得跨模式报价。
3. **性能影响**：统一目录和报价会读取能力及价格档，需要继续关注 PostgreSQL 和大目录下的查询性能。

## 文件清单

### 新增文件
- `internal/service/model_catalog.go`
- `internal/service/price_validation.go`
- `internal/handler/model_catalog.go`

### 修改文件
- `internal/service/feature_availability.go`
- `internal/service/errors.go`
- `internal/service/logical_models.go`
- `internal/service/task_creation.go`
- `internal/model/models_logical_model.go`

---

**实施日期**: 2026-08-22
**代码审查**: 建议进行 Code Review
**部署建议**: 建议先在测试环境验证后再上线
