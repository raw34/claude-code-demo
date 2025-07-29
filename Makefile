.PHONY: help build up down logs clean test dev-services dev-backend dev-frontend

# 默认目标
help:
	@echo "Available commands:"
	@echo "  make build         - Build all Docker images"
	@echo "  make up            - Start all services"
	@echo "  make down          - Stop all services"
	@echo "  make logs          - View logs"
	@echo "  make clean         - Clean up containers and volumes"
	@echo "  make test          - Run tests"
	@echo "  make dev-services  - Start MySQL and Redis for local development"
	@echo "  make dev-backend   - Run backend locally (requires dev-services)"
	@echo "  make dev-frontend  - Run frontend locally"

# 构建所有镜像
build:
	docker-compose build

# 启动所有服务
up:
	docker-compose up -d

# 停止所有服务
down:
	docker-compose down

# 查看日志
logs:
	docker-compose logs -f

# 清理容器和卷
clean:
	docker-compose down -v
	docker system prune -f

# 运行测试
test:
	cd backend && go test ./...

# 启动开发环境服务（MySQL和Redis）
dev-services:
	@echo "Starting MySQL and Redis for local development..."
	docker-compose up -d database redis
	@echo "Waiting for services to be ready..."
	@sleep 5
	@echo "Services started. MySQL on localhost:3306, Redis on localhost:6379"

# 本地运行后端
dev-backend:
	@echo "Starting backend development server..."
	cd backend && cp .env.local .env 2>/dev/null || true && go run cmd/server/main.go

# 本地运行前端
dev-frontend:
	@echo "Starting frontend development server..."
	cd frontend && yarn dev