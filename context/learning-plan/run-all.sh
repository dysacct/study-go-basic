#!/bin/bash

# Channel、Context 和 Select 学习示例运行脚本

echo "======================================"
echo "  Go 并发学习示例自动运行脚本"
echo "======================================"
echo ""

# 颜色定义
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

run_example() {
    local dir=$1
    local name=$2
    
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${GREEN}运行: $name${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""
    
    cd "$dir" && go run main.go
    
    echo ""
    echo "按 Enter 继续下一个示例..."
    read
    echo ""
    cd - > /dev/null
}

# 运行所有示例
run_example "01-channel-basics" "第1章：Channel 基础"
run_example "02-select-basics" "第2章：Select 多路复用"
run_example "03-context-basics" "第3章：Context 上下文"
run_example "04-combined-usage" "第4章：综合应用"
run_example "05-advanced-patterns" "第5章：高级模式"

echo ""
echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}恭喜！基础学习完成！${NC}"
echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo "是否继续学习实战模式？(y/n)"
read answer

if [ "$answer" == "y" ] || [ "$answer" == "Y" ]; then
    echo ""
    echo -e "${BLUE}开始第4阶段：Done Channel 实战模式...${NC}"
    echo ""
    cd 06-done-channel-patterns
    ./run-all.sh
    cd ..
fi

echo ""
echo -e "${GREEN}======================================"
echo "  所有示例运行完成！"
echo "======================================${NC}"
