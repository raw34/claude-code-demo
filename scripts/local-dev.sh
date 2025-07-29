#!/bin/bash

# 本地开发环境启动脚本

echo "🚀 Starting local development environment..."

# 检查Docker是否运行
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker first."
    exit 1
fi

# 启动MySQL和Redis
echo "📦 Starting MySQL and Redis containers..."
docker-compose up -d database redis

# 等待服务就绪
echo "⏳ Waiting for services to be ready..."
sleep 5

# 测试MySQL连接
echo "🔍 Testing MySQL connection..."
if mysql -h localhost -P 3306 -u userapp -puserpassword -e "SELECT 1" user_db > /dev/null 2>&1; then
    echo "✅ MySQL is ready on localhost:3306"
else
    echo "⚠️  MySQL connection failed. It might need more time to start."
fi

# 测试Redis连接
echo "🔍 Testing Redis connection..."
if redis-cli -h localhost -p 6379 -a redispassword ping > /dev/null 2>&1; then
    echo "✅ Redis is ready on localhost:6379"
else
    echo "⚠️  Redis connection failed. It might need more time to start."
fi

echo ""
echo "📝 Connection details for local development:"
echo "   MySQL:"
echo "     Host: localhost"
echo "     Port: 3306"
echo "     User: userapp"
echo "     Password: userpassword"
echo "     Database: user_db"
echo ""
echo "   Redis:"
echo "     Host: localhost"
echo "     Port: 6379"
echo "     Password: redispassword"
echo ""
echo "🎯 Next steps:"
echo "   1. Backend: cd backend && cp .env.local .env && go run cmd/server/main.go"
echo "   2. Frontend: cd frontend && yarn dev"
echo ""
echo "✨ Happy coding!"