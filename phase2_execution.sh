#!/bin/bash
# 阶段 2 完整执行脚本 - 简化版
# 分阶段克隆 + 监控 + Fork 分析

echo "=========================================="
echo "  阶段 2: 获取上游代码和分析 Fork"
echo "=========================================="
echo ""

# 配置
TARGET_DIR="/tmp/upstream-repo"
CHECK_INTERVAL=5
NO_CHANGE_TIMEOUT=600

# 清理旧的失败克隆
if [ -d "$TARGET_DIR" ]; then
    echo ">>> 清理旧目录"
    rm -rf "$TARGET_DIR"
fi

# 尝试克隆（depth=1）
echo "=========================================="
echo "  开始克隆（depth=1）"
echo "=========================================="
echo ""

gh repo clone ddcat-ai/open-ai-canvas "$TARGET_DIR" -- --depth=1 &
CLONE_PID=$!

echo "克隆进程 PID: $CLONE_PID"
echo "开始监控..."
echo ""

# 监控循环
old_size=0
no_change_time=0
start_time=$(date +%s)

while kill -0 $CLONE_PID 2>/dev/null; do
    if [ -d "$TARGET_DIR/.git/objects" ]; then
        new_size=$(du -sb "$TARGET_DIR/.git/objects" 2>/dev/null | cut -f1 || echo "0")
        current_time=$(date +%s)
        elapsed=$((current_time - start_time))
        
        if [ "$new_size" != "$old_size" ] && [ "$new_size" != "0" ]; then
            change=$((new_size - old_size))
            speed=$((change / CHECK_INTERVAL))
            echo "[$(date +%H:%M:%S)] 进度: $((new_size / 1024))KB (+$((change / 1024))KB), 速度: $((speed / 1024))KB/s"
            no_change_time=0
            old_size=$new_size
        else
            no_change_time=$((no_change_time + CHECK_INTERVAL))
            echo "[$(date +%H:%M:%S)] 无增量: ${no_change_time}s / ${NO_CHANGE_TIMEOUT}s"
            
            if [ $no_change_time -ge $NO_CHANGE_TIMEOUT ]; then
                echo "⚠️  10分钟无增量，终止克隆"
                kill $CLONE_PID 2>/dev/null
                break
            fi
        fi
    fi
    
    sleep $CHECK_INTERVAL
done

wait $CLONE_PID 2>/dev/null
CLONE_RESULT=$?

echo ""
if [ $CLONE_RESULT -eq 0 ] && [ -d "$TARGET_DIR/.git" ]; then
    cd "$TARGET_DIR"
    if git log --oneline -1 &>/dev/null; then
        echo "✅ 克隆成功"
    else
        echo "❌ 克隆失败：仓库不完整"
        exit 1
    fi
else
    echo "❌ 克隆失败"
    exit 1
fi

echo ""
echo "=========================================="
echo "  配置远程仓库"
echo "=========================================="

cd /d/open-ai-canvas

if git remote | grep -q "^upstream$"; then
    git remote remove upstream
fi

git remote add upstream "$TARGET_DIR"
git fetch upstream

echo "✅ upstream 配置完成"

echo ""
echo "=========================================="
echo "  分析 Fork 分支"
echo "=========================================="
echo ""

# 添加 Fork 远程仓库
for fork in "daQzi" "mdsxbm"; do
    echo ">>> 添加 Fork: $fork"
    if git remote | grep -q "^$fork$"; then
        git remote remove $fork
    fi
    git remote add $fork "https://github.com/$fork/open-ai-canvas.git"
done

# 获取 Fork 分支
echo ""
echo ">>> 获取 daQzi/codex/branding"
git fetch daQzi codex/branding 2>&1 | head -10

echo ""
echo ">>> 获取 mdsxbm/trae/agent-YUI35U"  
git fetch mdsxbm trae/agent-YUI35U 2>&1 | head -10

echo ""
echo "✅ Fork 分支获取完成"

