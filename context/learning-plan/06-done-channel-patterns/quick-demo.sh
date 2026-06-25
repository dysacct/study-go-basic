#!/bin/bash

# 快速演示脚本 - 无需手动按 Enter，自动运行所有示例

echo "======================================"
echo "  Done Channel 模式快速演示"
echo "======================================"
echo ""

GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

demos=(
    "01-basic-done.go:基础应用"
    "02-multiple-workers.go:多Worker协作"
    "03-server-graceful-shutdown.go:优雅关闭"
    "04-real-world-crawler.go:并发爬虫"
    "05-producer-consumer.go:生产者消费者"
)

for demo in "${demos[@]}"; do
    IFS=':' read -r file name <<< "$demo"
    
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${GREEN}演示: $name${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""
    
    timeout 10s go run "$file" 2>&1 || true
    
    echo ""
    sleep 1
done

echo -e "${GREEN}======================================"
echo "  演示完成！"
echo "======================================${NC}"
echo ""
echo "💡 提示："
echo "  - 使用 ./run-all.sh 可以逐步查看每个示例"
echo "  - 查看 YOUR-CODE-ANALYSIS.md 了解如何改进你的代码"

