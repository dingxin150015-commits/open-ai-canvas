#!/bin/bash
# 冲突解决执行脚本

set -e

LOG_FILE="/tmp/conflict_resolution.log"
echo "=== 冲突解决日志 ===" > "$LOG_FILE"
echo "开始时间: $(date)" >> "$LOG_FILE"
echo "" >> "$LOG_FILE"

log() {
    echo "$1" | tee -a "$LOG_FILE"
}

log "=========================================="
log "  阶段 1: 处理低风险文件（新增文件）"
log "=========================================="
log ""

# 1. interface.go - 接受上游
log ">>> 1/12 处理 interface.go"
if [ -f backend/internal/provider/interface.go ]; then
    git checkout --theirs backend/internal/provider/interface.go
    git add backend/internal/provider/interface.go
    log "✅ interface.go - 已接受上游版本"
else
    log "⚠️  interface.go - 文件不存在"
fi

# 2. registry.go - 接受上游
log ">>> 2/12 处理 registry.go"
if [ -f backend/internal/provider/registry.go ]; then
    git checkout --theirs backend/internal/provider/registry.go
    git add backend/internal/provider/registry.go
    log "✅ registry.go - 已接受上游版本"
else
    log "⚠️  registry.go - 文件不存在"
fi

# 3. discovery.go - 接受上游
log ">>> 3/12 处理 discovery.go"
if [ -f backend/internal/provider/bailian/discovery.go ]; then
    git checkout --theirs backend/internal/provider/bailian/discovery.go
    git add backend/internal/provider/bailian/discovery.go
    log "✅ discovery.go - 已接受上游版本"
else
    log "⚠️  discovery.go - 文件不存在"
fi

# 4. channel_model_catalog_plugin.go - 接受上游
log ">>> 4/12 处理 channel_model_catalog_plugin.go"
if [ -f backend/internal/service/channel_model_catalog_plugin.go ]; then
    git checkout --theirs backend/internal/service/channel_model_catalog_plugin.go
    git add backend/internal/service/channel_model_catalog_plugin.go
    log "✅ channel_model_catalog_plugin.go - 已接受上游版本"
else
    log "⚠️  channel_model_catalog_plugin.go - 文件不存在"
fi

log ""
log "✅ 阶段 1 完成"
log ""

log "=========================================="
log "  阶段 1 完成，准备进入阶段 2"
log "=========================================="
log ""
log "下一步: 手动处理协议注册文件"
log "  - models.go"
log "  - provider.go"  
log "  - model-protocols.ts"
log ""

