# Deployment Guide

Guide for deploying Notes API to production.

## Pre-Deployment Checklist

- [ ] All tests passing locally
- [ ] Environment variables configured
- [ ] Database migrations prepared
- [ ] SSL certificates ready (for HTTPS)
- [ ] Backup plan in place
- [ ] Monitoring configured
- [ ] Logging centralized (optional)

## Deployment Methods

### Method 1: Direct Binary Deployment

#### 1. Build Binary

```bash
go build -ldflags "-s -w" -o notes-api ./cmd/api/main.go
```

Flags:
- `-s` - Strip symbol table
- `-w` - Strip DWARF table
- Result: Smaller binary (~8-10 MB)

#### 2. Transfer to Server

```bash
scp notes-api user@server:/opt/notes-api/
```

#### 3. Setup Systemd Service

Create `/etc/systemd/system/notes-api.service`:

```ini
[Unit]
Description=Notes API Service
After=network.target mysql.service
Wants=network.target

[Service]
Type=simple
User=notesapi
WorkingDirectory=/opt/notes-api
EnvironmentFile=/opt/notes-api/.env
ExecStart=/opt/notes-api/notes-api
Restart=on-failure
RestartSec=10

# Security
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=yes

[Install]
WantedBy=multi-user.target
```

#### 4. Enable and Start Service

```bash
sudo systemctl enable notes-api
sudo systemctl start notes-api
sudo systemctl status notes-api
```

#### 5. View Logs

```bash
sudo journalctl -u notes-api -f
```

---

### Method 2: Docker Deployment

#### 1. Build Docker Image

Create `Dockerfile`:

```dockerfile
# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags "-s -w" \
    -o notes-api ./cmd/api/main.go

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Install runtime dependencies
RUN apk --no-cache add ca-certificates

# Copy binary from builder
COPY --from=builder /app/notes-api .

# Expose port
EXPOSE 3000

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:3000/health || exit 1

# Run
CMD ["./notes-api"]
```

#### 2. Build Image

```bash
docker build -t notes-api:latest .
```

#### 3. Push to Registry

```bash
docker tag notes-api:latest registry.example.com/notes-api:latest
docker push registry.example.com/notes-api:latest
```

#### 4. Run with Docker Compose

Create `docker-compose.prod.yml`:

```yaml
version: '3.9'

services:
  db:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: ${DB_PASS}
      MYSQL_DATABASE: ${DB_NAME}
    volumes:
      - mysql_data:/var/lib/mysql
    ports:
      - "3306:3306"
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost"]
      timeout: 3s
      retries: 5

  api:
    image: registry.example.com/notes-api:latest
    environment:
      DB_HOST: db
      DB_USER: ${DB_USER}
      DB_PASS: ${DB_PASS}
      DB_NAME: ${DB_NAME}
      DB_PORT: 3306
      APP_PORT: 3000
      ENV: production
      JWT_SECRET: ${JWT_SECRET}
      CORS_ALLOWED_ORIGINS: ${CORS_ALLOWED_ORIGINS}
    ports:
      - "3000:3000"
    depends_on:
      db:
        condition: service_healthy
    healthcheck:
      test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost:3000/health"]
      interval: 30s
      timeout: 3s
      retries: 3
    restart: unless-stopped

volumes:
  mysql_data:
```

#### 5. Deploy

```bash
docker-compose -f docker-compose.prod.yml up -d
```

---

### Method 3: Kubernetes Deployment

#### 1. Create Deployment Manifest

Create `k8s/deployment.yaml`:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: notes-api
  labels:
    app: notes-api
spec:
  replicas: 3
  selector:
    matchLabels:
      app: notes-api
  template:
    metadata:
      labels:
        app: notes-api
    spec:
      containers:
      - name: notes-api
        image: registry.example.com/notes-api:latest
        ports:
        - containerPort: 3000
        env:
        - name: DB_HOST
          valueFrom:
            configMapKeyRef:
              name: notes-api-config
              key: db-host
        - name: DB_USER
          valueFrom:
            secretKeyRef:
              name: notes-api-secret
              key: db-user
        - name: DB_PASS
          valueFrom:
            secretKeyRef:
              name: notes-api-secret
              key: db-pass
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: notes-api-secret
              key: jwt-secret
        livenessProbe:
          httpGet:
            path: /health
            port: 3000
          initialDelaySeconds: 10
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: 3000
          initialDelaySeconds: 5
          periodSeconds: 5
        resources:
          requests:
            cpu: 100m
            memory: 128Mi
          limits:
            cpu: 500m
            memory: 512Mi
```

#### 2. Create Service

Create `k8s/service.yaml`:

```yaml
apiVersion: v1
kind: Service
metadata:
  name: notes-api-service
spec:
  selector:
    app: notes-api
  ports:
  - protocol: TCP
    port: 80
    targetPort: 3000
  type: LoadBalancer
```

#### 3. Deploy to Kubernetes

```bash
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml
```

---

## Database Migration

### Using golang-migrate

```bash
# Install
go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Create connection
migrate -path migrations -database "mysql://user:pass@tcp(host:3306)/dbname" up
```

### Manual Migration

```bash
# Connect to database
mysql -h host -u user -p database_name

# Run SQL files
source migrations/000001_create_users_table.up.sql;
source migrations/000002_create_notes_table.up.sql;
```

## SSL/TLS Setup

### Using Let's Encrypt with Certbot

```bash
sudo certbot certonly --standalone -d api.example.com
```

### Configure Nginx as Reverse Proxy

Create `/etc/nginx/sites-available/notes-api`:

```nginx
server {
    listen 80;
    server_name api.example.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name api.example.com;

    ssl_certificate /etc/letsencrypt/live/api.example.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/api.example.com/privkey.pem;

    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;

    location / {
        proxy_pass http://localhost:3000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

Enable site:
```bash
sudo ln -s /etc/nginx/sites-available/notes-api /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl restart nginx
```

## Environment Configuration

### Production `.env` Example

```env
# Database
DB_HOST=db.example.com
DB_USER=produser
DB_PASS=strong_password_here
DB_NAME=notes_api_prod
DB_PORT=3306

# API
APP_PORT=3000
ENV=production

# JWT (minimum 32 characters, use strong random string)
JWT_SECRET=your_very_long_random_secret_key_minimum_32_chars

# CORS
CORS_ALLOWED_ORIGINS=https://app.example.com,https://www.example.com
```

### Security Best Practices

✅ **Do:**
- Use strong passwords (32+ chars)
- Store secrets in environment variables (not in code)
- Use HTTPS in production
- Enable CORS only for trusted origins
- Keep dependencies updated
- Regular backups

❌ **Don't:**
- Commit `.env` file to repository
- Use default credentials
- Expose sensitive logs
- Enable debug mode in production
- Run with excessive permissions

## Monitoring and Logging

### Health Check Endpoint

```bash
curl https://api.example.com/health
```

Should return:
```json
{
  "status": "OK"
}
```

### Application Metrics (Optional)

Consider adding metrics collection:
- Request count
- Response times
- Error rates
- Database connection pool status

Tools:
- Prometheus
- Grafana
- New Relic
- DataDog

### Centralized Logging (Optional)

Configure structured logging aggregation:
- ELK Stack (Elasticsearch, Logstash, Kibana)
- Loki + Grafana
- CloudWatch
- Splunk

Current setup uses Zerolog with JSON output for easy parsing.

## Backup Strategy

### Database Backup

```bash
# Daily backup
0 2 * * * mysqldump -h host -u user -p password database > /backup/notes_api_$(date +\%Y\%m\%d).sql

# Verify backup
mysql -h host -u user -p database < /backup/notes_api_backup.sql
```

### Backup Storage

- Off-site backup location
- Encrypted backups
- Regular restoration tests
- 30-day retention minimum

## Rollback Plan

### Before Deployment

1. Tag current version:
   ```bash
   git tag -a v1.0.0 -m "Production release"
   ```

2. Create backup of binary:
   ```bash
   cp notes-api notes-api.backup.$(date +%Y%m%d)
   ```

### Rollback Procedure

```bash
# Stop current service
sudo systemctl stop notes-api

# Restore previous binary
cp notes-api.backup.20240101 notes-api

# Start service
sudo systemctl start notes-api

# Verify
sudo systemctl status notes-api
```

## Performance Optimization

### Database Connection Pooling

Already configured in `pkg/database/mysql.go`:
- Max idle connections: 10
- Max open connections: 100
- Connection lifetime: 1 hour

### Caching (Future Enhancement)

Consider implementing:
- Redis for session cache
- In-memory cache for frequently accessed data
- HTTP caching headers

### Load Balancing

For multiple instances:
- Use Nginx or HAProxy
- Database connection pooling
- Session storage in database or Redis

## Security Headers

Already included in middleware:
```
Strict-Transport-Security: max-age=31536000
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
X-XSS-Protection: 1; mode=block
```

## Troubleshooting

### Application Won't Start

Check logs:
```bash
sudo journalctl -u notes-api -n 50
```

Common issues:
- Database connection failed
- Port already in use
- Missing environment variables
- File permission issues

### Database Connection Errors

```bash
# Test connection
mysql -h $DB_HOST -u $DB_USER -p$DB_PASS -e "SELECT 1"
```

### High Memory Usage

- Check for memory leaks
- Verify connection pool settings
- Monitor active goroutines

## Disaster Recovery

### Full System Restore

1. Restore database from backup
2. Deploy latest application binary
3. Verify all services are running
4. Test critical endpoints
5. Monitor for errors

## Support and Maintenance

- Monitor application health
- Review logs regularly
- Update dependencies quarterly
- Test disaster recovery annually
- Maintain documentation
