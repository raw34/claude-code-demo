# 用户管理系统架构设计

## 1. 系统架构图描述

```
┌────────────────────────────────────────────────────────────────────────────┐
│                           Docker Network (user-net)                         │
├────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌─────────────────┐     ┌─────────────────┐     ┌──────────────┐         │
│  │   nginx:80      │────▶│   go-api:8080   │────▶│   mysql:3306 │         │
│  │                 │     │                 │     │              │         │
│  │ Vue3+TS SPA    │     │  User API       │     │   Database   │         │
│  │  /app/dist      │     │  REST Service   │     │   user_db    │         │
│  └─────────────────┘     └─────────────────┘     └──────────────┘         │
│        frontend                 │ backend              database            │
│                                │                                           │
│                                ▼                                           │
│                          ┌─────────────────┐                              │
│                          │   redis:6379    │                              │
│                          │                 │                              │
│                          │ Session Store   │                              │
│                          │   Cache         │                              │
│                          └─────────────────┘                              │
│                                redis                                        │
└────────────────────────────────────────────────────────────────────────────┘
                                ▲
                                │
                          Host: localhost:80
```

### 容器说明：
- **frontend**: Nginx容器，托管Vue3+TypeScript构建后的静态文件，监听80端口
- **backend**: Go API服务容器，提供RESTful接口，监听8080端口  
- **database**: MySQL容器，存储用户数据，监听3306端口
- **redis**: Redis容器，管理用户登录态和session，监听6379端口

### 通信流程：
1. 用户访问 `http://localhost:80`
2. Nginx提供Vue3应用的静态文件
3. Vue3应用通过 `/api` 路径向后端发起请求
4. Nginx将 `/api` 请求代理到 `backend:8080`
5. Go API服务验证请求：
   - 检查Redis中的Session状态
   - 验证JWT Token有效性
6. Go API服务处理请求并访问MySQL数据库
7. 响应按原路返回给用户

### 登录态管理：
- 用户登录后，后端生成JWT Token并在Redis中创建Session
- 每个请求都会验证Redis中的Session状态
- Token过期或用户登出时，清除Redis中的Session
- 支持Token刷新机制，延长用户登录状态

## 2. RESTful API 接口设计

### 基础路径：`/api/v1`

| 方法 | 路径 | 描述 | 请求体 | 响应 |
|------|------|------|--------|------|
| POST | `/auth/register` | 用户注册 | `{username, email, password}` | `{id, username, email, created_at}` |
| POST | `/auth/login` | 用户登录 | `{email, password}` | `{token, user}` |
| POST | `/auth/logout` | 用户登出 | - | `{message}` |
| POST | `/auth/refresh` | 刷新令牌 | `{refresh_token}` | `{token, refresh_token}` |
| GET | `/users` | 获取用户列表 | - | `{users[], total, page, limit}` |
| GET | `/users/:id` | 获取用户详情 | - | `{id, username, email, created_at, updated_at}` |
| PUT | `/users/:id` | 更新用户信息 | `{username?, email?, password?}` | `{user}` |
| DELETE | `/users/:id` | 删除用户 | - | `{message}` |
| GET | `/users/profile` | 获取当前用户信息 | - | `{user}` |
| PUT | `/users/profile` | 更新当前用户信息 | `{username?, email?, password?}` | `{user}` |

### 请求/响应示例：

#### 注册用户
```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "username": "johndoe",
  "email": "john@example.com",
  "password": "securePassword123"
}

Response: 201 Created
{
  "id": 1,
  "username": "johndoe",
  "email": "john@example.com",
  "created_at": "2024-01-15T10:00:00Z"
}
```

#### 用户登录
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "securePassword123"
}

Response: 200 OK
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": 1,
    "username": "johndoe",
    "email": "john@example.com"
  }
}
```

### 认证方式
- 使用JWT Bearer Token认证
- Token放在请求头：`Authorization: Bearer <token>`
- Token有效期：15分钟
- Refresh Token有效期：7天
- Session状态存储在Redis中，支持跨服务器的分布式部署
- 每次请求都会验证Redis中的Session有效性

## 3. 数据库表结构设计

### users 表
```sql
CREATE TABLE `users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(50) NOT NULL UNIQUE,
  `email` varchar(100) NOT NULL UNIQUE,
  `password_hash` varchar(255) NOT NULL,
  `is_active` boolean DEFAULT true,
  `created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_email` (`email`),
  KEY `idx_username` (`username`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### refresh_tokens 表
```sql
CREATE TABLE `refresh_tokens` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `token` varchar(255) NOT NULL UNIQUE,
  `expires_at` timestamp NOT NULL,
  `created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_token` (`token`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_expires_at` (`expires_at`),
  CONSTRAINT `fk_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### user_sessions 表（可选，用于跟踪用户会话）
```sql
CREATE TABLE `user_sessions` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `ip_address` varchar(45),
  `user_agent` varchar(255),
  `last_activity` timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_last_activity` (`last_activity`),
  CONSTRAINT `fk_session_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

## 4. 容器间网络通信方案

### Docker网络配置
- 使用自定义bridge网络 `user-net`
- 所有容器连接到同一网络，可通过容器名互相访问
- 容器间通信使用内部端口，不需要暴露到宿主机

### 网络安全策略
1. **前端容器**：只暴露80端口到宿主机
2. **后端容器**：不暴露端口到宿主机，仅在内部网络访问
3. **数据库容器**：不暴露端口到宿主机，仅接受后端容器连接
4. **Redis容器**：不暴露端口到宿主机，仅接受后端容器连接，使用密码认证

### 服务发现
- 使用Docker内置DNS，通过容器名访问：
  - frontend → backend: `http://backend:8080`
  - backend → database: `mysql://database:3306`
  - backend → redis: `redis://redis:6379`

### Nginx反向代理配置
```nginx
server {
    listen 80;
    server_name localhost;
    
    # Vue应用静态文件
    location / {
        root /usr/share/nginx/html;
        try_files $uri $uri/ /index.html;
    }
    
    # API代理
    location /api {
        proxy_pass http://backend:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### 环境变量配置
```yaml
frontend:
  - VITE_API_BASE_URL=/api/v1
  
backend:
  - DB_HOST=database
  - DB_PORT=3306
  - DB_NAME=user_db
  - REDIS_HOST=redis
  - REDIS_PORT=6379
  - REDIS_PASSWORD=${REDIS_PASSWORD}
  - JWT_SECRET=${JWT_SECRET}
  
database:
  - MYSQL_DATABASE=user_db
  - MYSQL_ROOT_PASSWORD=${DB_ROOT_PASSWORD}
  
redis:
  - REDIS_PASSWORD=${REDIS_PASSWORD}
```