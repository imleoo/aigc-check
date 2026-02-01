#!/bin/bash

# AIGC-Check 构建和导出脚本

echo "开始构建 Docker 镜像..."

# 1. 检查 Docker 是否运行
if ! docker info > /dev/null 2>&1; then
    echo "错误: Docker 未运行，请先启动 Docker。"
    exit 1
fi

# 2. 构建镜像
# 使用 --no-cache 确保获取最新依赖
# 标记为 aigc-check:latest
echo "正在构建镜像 aigc-check:latest ..."
docker build --no-cache -t aigc-check:latest .

if [ $? -ne 0 ]; then
    echo "构建失败！请检查 Dockerfile 和网络连接。"
    exit 1
fi

echo "构建成功！"

# 3. 导出镜像
EXPORT_FILE="aigc-check-image.tar"
echo "正在导出镜像到 $EXPORT_FILE ..."
docker save -o $EXPORT_FILE aigc-check:latest

if [ $? -ne 0 ]; then
    echo "导出失败！"
    exit 1
fi

echo "镜像已导出为 $EXPORT_FILE"
echo ""
echo "=== 部署说明 ==="
echo "1. 将 $EXPORT_FILE 上传到服务器"
echo "2. 在服务器上运行: docker load -i $EXPORT_FILE"
echo "3. 启动容器:"
echo "   docker run -d -p 8080:8080 \\"
echo "     -e GEMINI_API_KEY=your_key_here \\"
echo "     --name aigc-check aigc-check:latest"
echo ""
echo "注意: 配置文件中的 data 目录位于容器内的 /app/data"
echo "如果需要持久化数据，请添加挂载: -v ./data:/app/data"
