# 环境变量配置指南

## 配置文件说明

所有环境变量配置文件都在项目根目录，采用统一管理：

```
项目根目录/
├── .env                 # 实际使用的配置文件（gitignore）
├── .env.example         # Docker环境模板
└── .env.local.example   # 本地开发模板
```

## 使用方式

### Docker部署
```bash
cp .env.example .env
docker-compose up -d
```

### 本地开发
```bash
# 1. 启动MySQL和Redis容器
docker-compose up -d database redis

# 2. 配置环境变量
cp .env.local.example .env

# 3. 启动后端（会自动读取根目录的.env）
cd backend
go run cmd/server/main.go

# 4. 启动前端
cd frontend
yarn dev
```

## 主要区别

两种环境的唯一区别是主机地址：
- **Docker环境**: DB_HOST=database, REDIS_HOST=redis
- **本地开发**: DB_HOST=localhost, REDIS_HOST=localhost

其他所有配置保持一致，便于管理和部署。