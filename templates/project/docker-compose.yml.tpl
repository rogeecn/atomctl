version: '3.8'

services:
  # PostgreSQL 数据库
  postgres:
    image: postgres:15-alpine
    container_name: {{.ProjectName}}-postgres
    environment:
      POSTGRES_DB: {{.ProjectName}}
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./database/init:/docker-entrypoint-initdb.d
    networks:
      - {{.ProjectName}}-network
    restart: unless-stopped
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 5s
      timeout: 5s
      retries: 5

  # Redis 缓存
  redis:
    image: redis:7-alpine
    container_name: {{.ProjectName}}-redis
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    networks:
      - {{.ProjectName}}-network
    restart: unless-stopped
    command: redis-server --appendonly yes
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 5s
      timeout: 3s
      retries: 5

  # 应用服务
  app:
    build:
      context: .
      dockerfile: Dockerfile.dev
    container_name: {{.ProjectName}}-app
    environment:
      - APP_MODE=development
      - HTTP_PORT=8080
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_NAME={{.ProjectName}}
      - DB_USER=postgres
      - DB_PASSWORD=password
      - REDIS_HOST=redis
      - REDIS_PORT=6379
      - REDIS_DB=0
      - JWT_SECRET_KEY=your-secret-key-here
      - HASHIDS_SALT=your-salt-here
    ports:
      - "8080:8080"
      - "9090:9090"  # 监控端口
    volumes:
      - .:/app
      - /app/vendor
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
    networks:
      - {{.ProjectName}}-network
    restart: unless-stopped

  # PgAdmin 数据库管理工具
  pgadmin:
    image: dpage/pgadmin4:latest
    container_name: {{.ProjectName}}-pgadmin
    environment:
      PGADMIN_DEFAULT_EMAIL: admin@example.com
      PGADMIN_DEFAULT_PASSWORD: admin
    ports:
      - "8081:80"
    volumes:
      - pgadmin_data:/var/lib/pgadmin
    networks:
      - {{.ProjectName}}-network
    restart: unless-stopped
    profiles:
      - tools

  # Redis 管理工具
  redis-commander:
    image: rediscommander/redis-commander:latest
    container_name: {{.ProjectName}}-redis-commander
    environment:
      REDIS_HOSTS: local:redis:6379
    ports:
      - "8082:8081"
    networks:
      - {{.ProjectName}}-network
    restart: unless-stopped
    profiles:
      - tools

volumes:
  postgres_data:
    driver: local
  redis_data:
    driver: local
  pgadmin_data:
    driver: local

networks:
  {{.ProjectName}}-network:
    driver: bridge