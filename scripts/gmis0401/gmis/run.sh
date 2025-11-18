#!/bin/bash
set -e

IMAGE_NAME="gmis0401-helper"
CONTAINER_NAME="gmisHelper"

# 获取当前目录路径
BASE_DIR="$(pwd)"

# 宿主机可执行文件和配置文件路径
HOST_EXEC_PATH="$BASE_DIR/gmis0401-linux-amd64"
HOST_CONFIG_PATH="$BASE_DIR/config.yml"

# 检查执行文件和配置文件是否存在
if [ ! -f "$HOST_EXEC_PATH" ]; then
  echo "❌ 找不到执行文件：$HOST_EXEC_PATH"
  exit 1
fi

if [ ! -f "$HOST_CONFIG_PATH" ]; then
  echo "❌ 找不到配置文件：$HOST_CONFIG_PATH"
  exit 1
fi

# 确保可执行文件有权限
chmod +x "$HOST_EXEC_PATH"

# 删除旧容器（如果存在）
if [ "$(docker ps -aq -f name=$CONTAINER_NAME)" ]; then
  echo "������ 删除旧容器 $CONTAINER_NAME ..."
  docker rm -f $CONTAINER_NAME
fi

# 启动容器并挂载执行文件和配置文件
echo "������ 启动容器 $CONTAINER_NAME ..."
docker run -d \
  --name $CONTAINER_NAME \
  -v "$HOST_EXEC_PATH":/app/gmis0401-linux-amd64 \
  -v "$HOST_CONFIG_PATH":/app/config.yml \
  --network host \
  --restart always \
  $IMAGE_NAME \
  /app/gmis0401-linux-amd64

echo "✅ 容器启动完成！使用 'docker logs -f $CONTAINER_NAME' 查看运行日志。"

