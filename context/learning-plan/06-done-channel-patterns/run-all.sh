#!/bin/bash

echo "======================================"
echo "  Done Channel 模式实战示例"
echo "======================================"
echo ""

GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

run_example() {
    local file=$1
    local name=$2
    
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${GREEN}运行: $name${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""
    
    go run "$file"
    
    echo ""
    echo -e "${YELLOW}按 Enter 继续下一个示例...${NC}"
    read
    echo ""
}

echo "这些示例展示了 Done Channel 在实际项目中的应用"
echo ""
echo "按 Enter 开始..."
read

run_example "01-basic-done.go" "示例1：基础应用（初始化等待）"
run_example "02-multiple-workers.go" "示例2：多 Worker 协作"
run_example "03-server-graceful-shutdown.go" "示例3：服务优雅关闭 ⭐"
run_example "04-real-world-crawler.go" "示例4：并发爬虫系统"
run_example "05-producer-consumer.go" "示例5：生产者消费者模式"

echo -e "${GREEN}======================================"
echo "  所有示例运行完成！"
echo "======================================${NC}"
echo ""
echo "💡 提示："
echo "  - 查看 README.md 了解详细说明"
echo "  - 尝试修改代码中的参数观察不同行为"
echo "  - 对比这些模式与你自己的代码"

