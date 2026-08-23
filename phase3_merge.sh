#!/bin/bash
# 阶段 3: 合并 upstream/main

set -e

echo "=========================================="
echo "  阶段 3: 合并 upstream/main"
echo "=========================================="
echo ""

# ========== 准备工作 ==========

echo ">>> 检查当前状态"
git status --short

current_branch=$(git branch --show-current)
echo "当前分支: $current_branch"

if [ "$current_branch" != "main" ]; then
    echo "⚠️  当前不在 main 分支"
    echo "切换到 main 分支..."
    git checkout main
fi

echo ""

# ========== 创建合并分支 ==========

echo "=========================================="
echo "  步骤 1: 创建合并分支"
echo "=========================================="
echo ""

merge_branch="merge/upstream-integration-$(date +%Y%m%d-%H%M%S)"
echo "创建分支: $merge_branch"

git checkout -b "$merge_branch"
echo "✅ 合并分支已创建"

echo ""

# ========== 开始合并 ==========

echo "=========================================="
echo "  步骤 2: 开始合并 upstream/main"
echo "=========================================="
echo ""

echo ">>> 合并信息"
echo "源: upstream/main (455 个新提交)"
echo "目标: $merge_branch"
echo "策略: --no-ff (保留合并历史)"
echo ""

echo ">>> 执行合并..."
git merge upstream/main --no-ff -m "merge: 合并上游 main 分支 (455 个新提交)

包含以下主要功能:
- 模型规格价格系统 (SKU)
- 方舟私域素材支持
- 画布系统大幅优化
- 插件系统重命名
- 文本生成能力
- 技能系统
- 计费系统优化
- 任务系统治理

上游提交: 455 个
文件变更: 420 个
代码行: +47,693

Co-Authored-By: Claude Opus 5 (1M context) <noreply@anthropic.com>
" 2>&1 || {
    merge_result=$?
    echo ""
    echo "=========================================="
    echo "  合并遇到冲突（预期行为）"
    echo "=========================================="
    echo ""
    
    echo ">>> 检查冲突文件"
    git status --short | grep "^UU\|^AA\|^DD\|^AU\|^UA\|^DU\|^UD"
    
    echo ""
    echo ">>> 冲突文件数量"
    conflict_count=$(git status --short | grep "^UU\|^AA\|^DD\|^AU\|^UA\|^DU\|^UD" | wc -l)
    echo "冲突文件: $conflict_count 个"
    
    echo ""
    echo "=========================================="
    echo "  合并暂停，等待手动解决冲突"
    echo "=========================================="
    echo ""
    echo "当前分支: $merge_branch"
    echo ""
    echo "下一步操作:"
    echo "1. 查看冲突文件列表: git status"
    echo "2. 逐个解决冲突"
    echo "3. 标记已解决: git add <file>"
    echo "4. 继续合并: git commit"
    echo ""
    
    exit 0
}

# 如果没有冲突（不太可能）
echo ""
echo "=========================================="
echo "  合并成功（无冲突）"
echo "=========================================="
echo ""
echo "当前分支: $merge_branch"
echo ""

