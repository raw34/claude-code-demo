# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a user management system built with a frontend-backend separated architecture using Vue3, Go, and MySQL, all orchestrated with Docker Compose.

## Repository Structure

```
claude-code-demo/
├── frontend/              # Vue3 SPA application
│   ├── src/              # Source code
│   │   ├── api/          # API client
│   │   ├── views/        # Page components
│   │   ├── stores/       # Pinia stores
│   │   └── router/       # Vue Router config
│   ├── Dockerfile        # Frontend container config
│   └── nginx.conf        # Nginx configuration
├── backend/              # Go API service
│   ├── cmd/server/       # Application entry point
│   ├── internal/         # Internal packages
│   │   ├── config/       # Configuration
│   │   ├── handlers/     # HTTP handlers
│   │   ├── middleware/   # Middleware functions
│   │   ├── models/       # Data models
│   │   ├── repository/   # Data access layer
│   │   └── service/      # Business logic
│   ├── migrations/       # Database migrations
│   └── Dockerfile        # Backend container config
├── docs/                 # Documentation
│   └── architecture.md   # System architecture design
├── docker-compose.yml    # Container orchestration
├── Makefile             # Build automation
└── .env.example         # Environment variables template
```

## System Architecture

The system uses a multi-tier architecture:
- **Frontend**: Vue3 + TypeScript SPA served by Nginx on port 80
- **Backend**: Go REST API service on port 8080 (internal)
- **Database**: MySQL 8.0 for data persistence
- **Cache**: Redis for session management and caching

All services communicate through a Docker bridge network `user-net`.

## Development Commands

### Quick Start
```bash
# Copy environment variables
cp .env.example .env

# Build and start all services
make build
make up

# View logs
make logs

# Stop services
make down
```

### Frontend Development
```bash
cd frontend
yarn install        # Install dependencies
yarn dev           # Start dev server (port 5173)
yarn build         # Build for production
yarn lint          # Run linter
yarn type-check    # Check TypeScript types
```

### Backend Development
```bash
cd backend
go mod download     # Install dependencies
go run cmd/server/main.go  # Run server
go test ./...       # Run tests
go build -o server cmd/server/main.go  # Build binary
```

### Database Operations
```bash
# Access MySQL container
docker exec -it user-db mysql -uroot -p

# Run migrations manually
docker exec -i user-db mysql -uroot -p user_db < backend/migrations/001_init.sql
```

## API Endpoints

Base URL: `/api/v1`

### Authentication
- `POST /auth/register` - User registration
- `POST /auth/login` - User login
- `POST /auth/logout` - User logout (requires auth)
- `POST /auth/refresh` - Refresh access token

### User Management (requires authentication)
- `GET /users` - List users with pagination
- `GET /users/:id` - Get user details
- `PUT /users/:id` - Update user
- `DELETE /users/:id` - Delete user
- `GET /users/profile` - Get current user profile
- `PUT /users/profile` - Update current user profile

## Key Technologies

### Frontend
- **Vue 3**: Composition API with setup syntax
- **TypeScript**: Type safety and better IDE support
- **Pinia**: State management
- **Vue Router**: Client-side routing
- **Element Plus**: UI component library
- **Axios**: HTTP client with interceptors
- **Vite**: Build tool and dev server
- **Yarn**: Package management

### Backend
- **Gin**: Web framework
- **GORM**: ORM for database operations
- **JWT**: Authentication with access/refresh tokens
- **bcrypt**: Password hashing
- **godotenv**: Environment variable management
- **go-redis**: Redis client for session management

### Infrastructure
- **Docker Compose**: Container orchestration
- **Nginx**: Static file serving and API proxy
- **MySQL 8.0**: Relational database
- **Redis**: Session store and caching

## Security Features

- JWT-based authentication with 15-minute access tokens
- Refresh tokens with 7-day expiry stored in MySQL
- Session management with Redis for immediate revocation
- Password hashing with bcrypt
- CORS configuration for cross-origin requests
- Input validation on all endpoints
- SQL injection protection via GORM
- Environment-based configuration
- Redis password protection

## Development Tips

1. **API Testing**: Use tools like Postman or curl to test API endpoints
2. **Database Access**: MySQL is not exposed to host by default for security
3. **Hot Reload**: Frontend supports hot reload in development mode
4. **Logs**: Use `docker-compose logs -f [service]` to view specific service logs
5. **Environment Variables**: Never commit `.env` file, use `.env.example` as template