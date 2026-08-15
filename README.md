# go-score-service

## 项目说明
一个轻量级的用户积分管理服务，基于内存存储，提供用户注册、积分增减、积分查询和排行榜功能。支持并发操作，适用于小规模用户积分场景。

## 标准命令
go build ./...     # 编译
go test ./...      # 测试
go run ./cmd       # 启动

## Docker 命令
docker build -t go-score-service .                      # 构建镜像
docker run -it go-score-service:latest                  # 运行容器
docker build --platform linux/arm64 -t go-score-service .  # 构建 arm64
