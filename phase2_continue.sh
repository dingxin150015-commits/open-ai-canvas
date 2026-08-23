#!/bin/bash
# 阶段 2 继续执行：加深 main + 统一获取 feature

set -e

echo "=========================================="
echo "  阶段 2 继续：加深历史 + 获取 feature"
echo "=========================================="
echo ""

REPO_DIR="/tmp/upstream-repo"
CHECK_INTERVAL=5
NO_CHANGE_TIMEOUT=600

# ========== 函数：监控下载 ==========
monitor_process() {
    local pid=$1
    local watch_dir="$REPO_DIR/.git/objects"
    local old_size=0
    local no_change_time=0
    local start_time=$(date +%s)

    echo ">>> 开始监控进程 (PID: $pid)"

    while kill -0 $pid 2>/dev/null; do
        if [ ! -d "$watch_dir" ]; then
            sleep $CHECK_INTERVAL
            continue
        fi

        local new_size=$(du -sb "$watch_dir" 2>/dev/null | cut -f1 || echo "0")
        local current_time=$(date +%s)
        local elapsed=$((current_time - start_time))

        if [ "$new_size" != "$old_size" ] && [ "$new_size" != "0" ]; then
            local change=$((new_size - old_size))
            local speed=$((change / CHECK_INTERVAL))
            echo "[$(date +%H:%M:%S)] 进度: $((new_size / 1024))KB (+$((change / 1024))KB), 速度: $((speed / 1024))KB/s, 已用: ${elapsed}s"
            no_change_time=0
            old_size=$new_size
        else
            no_change_time=$((no_change_time + CHECK_INTERVAL))
            if [ $no_change_time -le $NO_CHANGE_TIMEOUT ]; then
                echo "[$(date +%H:%M:%S)] 无增量: ${no_change_time}s / ${NO_CHANGE_TIMEOUT}s"
            fi

            if [ $no_change_time -ge $NO_CHANGE_TIMEOUT ]; then
                echo "⚠️  10分钟无增量，终止进程"
                kill $pid 2>/dev/null
                return 1
            fi
        fi

        sleep $CHECK_INTERVAL
    done

    wait $pid
    return $?
}

# ========== 第一部分：加深 main 历史 ==========

echo "=========================================="
echo "  第一部分：加深 main 分支历史"
echo "=========================================="
echo ""

cd "$REPO_DIR"

current_depth=$(git rev-list --count HEAD)
echo "当前深度: $current_depth 个提交"
echo ""

# 阶段 2-6: 逐步加深
for target_depth in 10 50 100 200; do
    current=$(git rev-list --count HEAD)
    
    if [ $current -ge $target_depth ]; then
        echo ">>> 跳过阶段 $target_depth（当前已有 $current 个提交）"
        continue
    fi
    
    deepen_by=$((target_depth - current))
    
    echo ""
    echo "=========================================="
    echo "  阶段 depth=$target_depth（加深 $deepen_by 个提交）"
    echo "=========================================="
    echo ""
    
    git fetch --deepen=$deepen_by 2>&1 &
    fetch_pid=$!
    
    if monitor_process $fetch_pid; then
        new_depth=$(git rev-list --count HEAD)
        echo "✅ 成功，当前深度: $new_depth"
    else
        echo "❌ 失败，继续下一阶段"
    fi
done

# 阶段 final: 完整历史
echo ""
echo "=========================================="
echo "  阶段 FINAL：获取完整历史（--unshallow）"
echo "=========================================="
echo ""

if [ -f .git/shallow ]; then
    git fetch --unshallow 2>&1 &
    fetch_pid=$!
    
    if monitor_process $fetch_pid; then
        final_depth=$(git rev-list --count HEAD)
        echo "✅ 完整历史获取成功！"
        echo "✅ 总提交数: $final_depth"
    else
        echo "⚠️  完整历史获取失败，但已有部分历史"
        partial_depth=$(git rev-list --count HEAD)
        echo "当前深度: $partial_depth"
    fi
else
    echo "⚠️  已经是完整克隆"
fi

echo ""
echo "=========================================="
echo "  main 分支历史加深完成"
echo "=========================================="

final_main_depth=$(git rev-list --count HEAD)
echo "最终深度: $final_main_depth 个提交"

# ========== 第二部分：获取 feature 分支 ==========

echo ""
echo "=========================================="
echo "  第二部分：获取 feature 分支（统一克隆）"
echo "=========================================="
echo ""

cd "$REPO_DIR"

echo ">>> 获取 feature 分支（depth=1）"
git fetch origin feature:feature 2>&1 &
fetch_pid=$!

if monitor_process $fetch_pid; then
    echo "✅ feature 分支获取成功"
    
    # 检查 feature 提交数
    if git rev-parse feature &>/dev/null; then
        feature_depth=$(git rev-list --count feature)
        echo "feature 分支深度: $feature_depth 个提交"
    fi
else
    echo "❌ feature 分支获取失败"
fi

# ========== 第三部分：更新本地引用 ==========

echo ""
echo "=========================================="
echo "  第三部分：更新本地引用"
echo "=========================================="
echo ""

cd /d/open-ai-canvas

echo ">>> 更新 upstream 远程仓库"
git fetch upstream 2>&1 | head -20

echo ""
echo "✅ 本地引用已更新"

# ========== 总结 ==========

echo ""
echo "=========================================="
echo "  阶段 2 继续执行完成"
echo "=========================================="
echo ""

cd "$REPO_DIR"

echo ">>> 最终状态"
echo ""
echo "main 分支:"
main_count=$(git rev-list --count main 2>/dev/null || echo "0")
echo "  - 提交数: $main_count"

if git rev-parse feature &>/dev/null 2>&1; then
    feature_count=$(git rev-list --count feature 2>/dev/null || echo "0")
    echo "feature 分支:"
    echo "  - 提交数: $feature_count"
else
    echo "feature 分支: 未获取"
fi

echo ""
echo "是否完整克隆:"
if [ -f .git/shallow ]; then
    echo "  - ❌ 浅克隆（部分历史）"
else
    echo "  - ✅ 完整克隆"
fi

echo ""
echo "下一步: 生成分析报告"

