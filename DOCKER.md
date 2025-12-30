# Docker Development Guide

This guide covers running SmoothWeb with Docker and Docker Compose for development.

## Prerequisites

- Docker (20.10+)
- Docker Compose (2.0+)

## Quick Start

1. **Copy environment files:**
   ```bash
   cp backend/.env.example backend/.env
   cp frontend/.env.example frontend/.env
   ```

2. **Start all services:**
   ```bash
   docker-compose up --build
   ```

3. **Access the application:**
   - Backend API: http://localhost:8080
   - Frontend Dev: http://localhost:5173

## Project Structure

```
.
├── backend/
│   ├── Dockerfile          # Backend container image
│   ├── .dockerignore       # Files to exclude from build context
│   └── .env.example       # Backend environment variables
├── frontend/
│   ├── Dockerfile.dev      # Frontend dev container
│   ├── .dockerignore       # Files to exclude from build context
│   └── .env.example       # Frontend environment variables
└── docker-compose.yml       # Multi-service orchestration
```

## Services

### Backend Service

- **Image:** Built from `./backend/Dockerfile`
- **Ports:** 8080:8080 (configurable via `BACKEND_PORT`)
- **Volumes:**
  - `backend_data`: Persists SQLite database
  - `backend_logs`: Persistent log storage
- **Environment Variables:**
  - `SERVER_PORT=8080`: API server port
  - `GIN_MODE=release`: Gin mode (debug/release)
  - `DB_PATH=/app/data/smoothweb.db`: Database file path
  - `JWT_SECRET`: JWT signing key (change in production!)
  - `JWT_EXPIRATION=24h`: Token expiration time
  - `CORS_ORIGINS=http://localhost:5173`: Allowed CORS origins
  - `LOG_LEVEL=info`: Logging level

### Frontend Service

- **Image:** Built from `./frontend/Dockerfile.dev`
- **Ports:** 5173:5173 (configurable via `FRONTEND_PORT`)
- **Volumes:**
  - `./frontend:/app`: Source code mount for hot-reload
  - `/app/node_modules`: Named volume for dependency isolation
- **Environment Variables:**
  - `VITE_API_URL=http://localhost:8080/api/v1`: Backend API URL

### Networks

- **smoothweb-network**: Bridge network for service communication

## Development Workflow

### Hot Reload

- **Backend:** Requires rebuild with `docker-compose up --build backend`
- **Frontend:** Automatic via Vite dev server with source code mount

### Database Persistence

Database is persisted in the `backend_data` Docker volume:
- Survives container restarts
- Survives `docker-compose down`
- To reset database: `docker-compose down -v`

### Logs

Backend logs are persisted in the `backend_logs` Docker volume:
- Access via `docker-compose logs backend`
- Persistent across container restarts
- To clear logs: `docker-compose down -v`

## Common Commands

### Start Services
```bash
# Build and start all services
docker-compose up --build

# Start in background with logs
docker-compose up -d --build

# Start specific service
docker-compose up backend
docker-compose up frontend
```

### Stop Services
```bash
# Stop all services (keeps volumes)
docker-compose down

# Stop all services and remove volumes (resets database)
docker-compose down -v
```

### View Logs
```bash
# Follow all service logs
docker-compose logs -f

# Follow specific service logs
docker-compose logs -f backend
docker-compose logs -f frontend

# View last 100 lines
docker-compose logs --tail=100
```

### Rebuild Services
```bash
# Rebuild backend after code changes
docker-compose up --build backend

# Rebuild frontend after dependency changes
docker-compose up --build frontend

# Force rebuild without cache
docker-compose build --no-cache
```

### Run Commands in Container
```bash
# Access backend shell
docker-compose exec backend sh

# Run tests in backend
docker-compose exec backend go test ./...

# Access frontend shell
docker-compose exec frontend sh

# Install npm package in frontend
docker-compose exec frontend npm install <package>
```

## Configuration

### Custom Ports

Override default ports in docker-compose.yml or via environment:
```bash
BACKEND_PORT=9000 FRONTEND_PORT=3000 docker-compose up
```

### Custom Environment Variables

Override any environment variable:
```bash
JWT_SECRET=custom-secret docker-compose up
```

### Production Considerations

For production deployment:
1. **Change JWT_SECRET** to a secure random value
2. **Update CORS_ORIGINS** to include production domain
3. **Set GIN_MODE** to `release` (already default)
4. **Use backend/production.Dockerfile** (not yet created)
5. **Disable hot-reload mounts** in frontend Dockerfile
6. **Remove debug environment variables**

## Troubleshooting

### Port Already in Use
```bash
# Check what's using port 8080
sudo lsof -i :8080

# Change backend port
BACKEND_PORT=9000 docker-compose up
```

### Database Locked
```bash
# Stop all services first
docker-compose down

# Then start again
docker-compose up
```

### Frontend Not Connecting to Backend
1. Check backend is running: `docker-compose ps`
2. Check backend health: `docker-compose logs backend`
3. Verify `VITE_API_URL` matches backend port
4. Check CORS settings in backend environment

### Build Context Issues
```bash
# Verify .dockerignore is working
docker-compose config

# Clean build cache
docker-compose build --no-cache
```

## Cleaning Up

### Remove All Containers
```bash
docker-compose down
```

### Remove Containers and Volumes
```bash
docker-compose down -v
```

### Remove Containers, Volumes, and Images
```bash
docker-compose down -v --rmi all
```

### Prune Docker System
```bash
docker system prune -a
```

## Development Tips

1. **Backend Changes:**
   - Go code changes require rebuild: `docker-compose up --build backend`
   - Database migrations should be run manually inside container

2. **Frontend Changes:**
   - Hot reload is enabled automatically
   - Component changes reflect immediately in browser
   - New dependencies require container rebuild

3. **Database Access:**
   ```bash
   docker-compose exec backend sqlite3 /app/data/smoothweb.db
   ```

4. **Viewing Logs:**
   - All logs also saved to `backend_logs` volume
   - Use `docker-compose logs -f` for real-time monitoring
