# 用户管理系统

基于 Vue3 + TypeScript + Go + MySQL + Redis 的前后端分离用户管理系统，使用 Docker Compose 进行容器编排。

## 系统特性

- 🎨 **现代化前端**：Vue3 + TypeScript + Element Plus 构建的响应式UI
- 🚀 **高性能后端**：Go + Gin 框架的 RESTful API
- 🔐 **双重认证**：JWT Token + Redis Session 的安全认证机制
- 🐳 **容器化部署**：Docker Compose 一键部署所有服务
- 📦 **模块化设计**：清晰的分层架构，易于扩展和维护
- 🔄 **即时登出**：基于 Redis 的 Session 管理，支持即时失效

## 快速开始

### 前置要求

- Docker 和 Docker Compose
- Node.js 18+ 和 Yarn 1.x（本地开发）
- Go 1.21+（本地开发）
- Git

### 安装步骤

1. 克隆项目
```bash
git clone https://github.com/yourusername/claude-code-demo.git
cd claude-code-demo
```

2. 配置环境变量
```bash
cp .env.example .env
# 编辑 .env 文件，设置以下关键配置：
# - DB_ROOT_PASSWORD: MySQL root密码
# - DB_PASSWORD: 应用数据库密码
# - REDIS_PASSWORD: Redis密码
# - JWT_SECRET: JWT签名密钥
```

3. 使用 Docker Compose 启动所有服务
```bash
make build
make up
```

4. 访问应用
- 打开浏览器访问 http://localhost
- 使用注册功能创建新用户
- 登录后可以管理用户信息

### 停止服务
```bash
make down
```

## 项目结构

```
.
├── frontend/              # 前端 Vue3 + TypeScript 应用
│   ├── src/              
│   │   ├── api/          # API 接口封装
│   │   ├── views/        # 页面组件
│   │   ├── stores/       # Pinia 状态管理
│   │   ├── router/       # Vue Router 路由配置
│   │   └── types/        # TypeScript 类型定义
│   ├── Dockerfile        
│   ├── nginx.conf        # Nginx 配置
│   └── package.json      
├── backend/              # 后端 Go API 服务
│   ├── cmd/server/       # 应用入口
│   ├── internal/         
│   │   ├── config/       # 配置管理
│   │   ├── handlers/     # HTTP 处理器
│   │   ├── middleware/   # 中间件
│   │   ├── models/       # 数据模型
│   │   ├── repository/   # 数据访问层
│   │   └── service/      # 业务逻辑层
│   ├── migrations/       # 数据库迁移脚本
│   ├── Dockerfile        
│   └── go.mod           
├── docs/                 # 项目文档
│   └── architecture.md   # 详细架构设计
├── docker-compose.yml    # Docker Compose 配置
├── Makefile             # 自动化脚本
├── .env.example         # 环境变量模板
└── CLAUDE.md            # Claude Code 开发指南
```

## 功能列表

- ✅ 用户注册/登录/登出
- ✅ JWT Token 认证（15分钟过期）
- ✅ Refresh Token 机制（7天过期）
- ✅ Redis Session 管理
- ✅ 用户列表查看（支持分页）
- ✅ 用户信息编辑
- ✅ 用户删除（不能删除自己）
- ✅ 个人资料管理
- ✅ 密码修改
- ✅ 健康检查端点

## 技术栈

### 前端
- **Vue 3.4**: Composition API + `<script setup>` 语法
- **TypeScript 5.3**: 类型安全，更好的开发体验
- **Vue Router 4**: 客户端路由管理
- **Pinia**: 新一代状态管理
- **Element Plus**: 企业级 UI 组件库
- **Axios**: HTTP 客户端，支持请求/响应拦截
- **Vite 5**: 极速的开发服务器和构建工具
- **Yarn**: 快速、可靠的包管理器

### 后端
- **Go 1.21**: 高性能的编译型语言
- **Gin 1.8**: 轻量级 Web 框架
- **GORM**: 功能强大的 ORM 库
- **JWT-go**: JWT 认证实现
- **go-redis/v9**: Redis 客户端
- **bcrypt**: 密码加密
- **MySQL 8.0**: 关系型数据库

### 基础设施
- **Docker**: 容器化技术
- **Docker Compose**: 多容器编排
- **Nginx**: 静态文件服务和反向代理
- **Redis 7**: 高性能缓存和 Session 存储

## 开发指南

### 本地开发环境设置

#### 1. 数据库和 Redis
```bash
# 只启动数据库和 Redis 服务
docker-compose up -d database redis

# 验证服务状态
docker-compose ps

# 测试MySQL连接
mysql -h localhost -P 3306 -u userapp -puserpassword user_db

# 测试Redis连接
redis-cli -h localhost -p 6379 -a redispassword ping
```

#### 2. 后端开发
```bash
cd backend

# 复制本地开发环境配置
cp .env.local .env

# 安装依赖
go mod download

# 运行服务（使用本地环境变量）
go run cmd/server/main.go

# 或者使用环境变量文件
source .env && go run cmd/server/main.go

# 运行测试
go test ./...

# 构建二进制文件
go build -o server cmd/server/main.go
```

#### 3. 前端开发
```bash
cd frontend

# 安装依赖
yarn install

# 启动开发服务器（支持热重载）
yarn dev

# 构建生产版本
yarn build

# 类型检查
yarn type-check

# 代码格式化
yarn format

# ESLint 检查
yarn lint
```

### 常用命令

```bash
# Docker 操作
make build    # 构建所有镜像
make up       # 启动所有服务
make down     # 停止所有服务
make logs     # 查看服务日志
make clean    # 清理容器和数据卷

# 开发调试
docker-compose logs -f backend    # 查看后端日志
docker-compose logs -f frontend   # 查看前端日志
docker exec -it user-db mysql -u root -p    # 连接数据库
docker exec -it user-redis redis-cli    # 连接 Redis
```

## API 文档

详细的 API 文档请参考 [docs/architecture.md](docs/architecture.md)

### 主要接口

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| POST | `/api/v1/auth/register` | 用户注册 | 否 |
| POST | `/api/v1/auth/login` | 用户登录 | 否 |
| POST | `/api/v1/auth/logout` | 用户登出 | 是 |
| POST | `/api/v1/auth/refresh` | 刷新令牌 | 否 |
| GET | `/api/v1/users` | 获取用户列表 | 是 |
| GET | `/api/v1/users/:id` | 获取用户详情 | 是 |
| PUT | `/api/v1/users/:id` | 更新用户信息 | 是 |
| DELETE | `/api/v1/users/:id` | 删除用户 | 是 |
| GET | `/api/v1/users/profile` | 获取当前用户信息 | 是 |
| PUT | `/api/v1/users/profile` | 更新当前用户信息 | 是 |

## 安全特性

- **密码安全**: 使用 bcrypt 加密存储，不可逆
- **双重认证**: JWT Token + Redis Session 双重验证
- **即时失效**: 登出即刻清除 Redis Session
- **跨域保护**: CORS 配置限制访问来源
- **SQL注入防护**: 使用 GORM 参数化查询
- **敏感信息保护**: 环境变量管理敏感配置
- **密码强度**: 前后端双重验证密码复杂度
- **会话管理**: 支持查看和管理所有登录会话

## 部署建议

### 生产环境配置

1. **环境变量**
   - 使用强密码和复杂的 JWT Secret
   - 不要使用示例中的默认值
   - 使用密钥管理服务存储敏感信息

2. **网络安全**
   - 使用 HTTPS（配置 SSL 证书）
   - 配置防火墙规则
   - 限制数据库和 Redis 的访问

3. **性能优化**
   - 启用 Redis 持久化
   - 配置数据库连接池
   - 使用 CDN 加速静态资源

4. **监控告警**
   - 配置日志收集
   - 设置性能监控
   - 配置异常告警

## 本地开发与Docker服务连接

### 快速启动本地开发环境

```bash
# 方式1：使用提供的脚本
./scripts/local-dev.sh

# 方式2：使用Makefile
make dev-services  # 启动MySQL和Redis
make dev-backend   # 在新终端运行后端
make dev-frontend  # 在新终端运行前端
```

### 连接配置

Docker Compose已配置端口映射，允许本地开发连接：

- **MySQL**: `localhost:3306`
  - 用户名: `userapp`
  - 密码: `userpassword`
  - 数据库: `user_db`

- **Redis**: `localhost:6379`
  - 密码: `redispassword`

### 后端本地开发配置

后端项目包含 `.env.local` 文件，已预配置本地开发环境变量：

```bash
cd backend
cp .env.local .env  # 使用本地配置
go run cmd/server/main.go
```

## 故障排查

### 常见问题

1. **前端依赖安装失败**
   ```bash
   # 清理缓存重新安装
   rm -rf node_modules yarn.lock
   yarn install
   ```

2. **后端依赖问题**
   ```bash
   # 清理模块缓存
   go clean -modcache
   go mod download
   ```

3. **Docker 容器启动失败**
   ```bash
   # 查看详细日志
   docker-compose logs [service-name]
   
   # 重新构建镜像
   docker-compose build --no-cache [service-name]
   ```

4. **数据库连接失败**
   - 检查环境变量配置
   - 确认数据库服务健康状态
   - 验证网络连通性
   
5. **本地开发数据库权限问题**
   ```bash
   # 如果遇到用户权限错误，清理数据卷重新初始化
   docker-compose down -v
   docker-compose up -d database redis
   # 等待30秒让MySQL完全初始化
   sleep 30
   ```

## 贡献指南

1. Fork 本仓库
2. 创建你的特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交你的改动 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启一个 Pull Request

### 代码规范

- 前端：遵循 Vue 3 风格指南和 TypeScript 最佳实践
- 后端：遵循 Go 官方代码规范
- 提交信息：使用语义化提交规范

## 许可证

MIT License

## 作者

Randy

## 致谢

- Vue.js 团队提供的优秀前端框架
- Go 团队提供的高性能编程语言
- 所有开源项目贡献者

---

🤖 使用 [Claude Code](https://claude.ai/code) 生成和维护

Co-Authored-By: Claude <noreply@anthropic.com>