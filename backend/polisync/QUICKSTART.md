# PoliSync Quick Start

Get PoliSync running in 5 minutes.

## Prerequisites

- .NET 8 SDK ([download](https://dotnet.microsoft.com/download/dotnet/8.0))
- Docker Desktop (for dependencies)
- PowerShell 7+ or Bash

## Step 1: Start Dependencies

```bash
# Start PostgreSQL, Redis, Kafka
cd backend/inscore
docker-compose up -d postgres redis kafka

# Run database migrations (Go side)
./scripts/run_migrations.sh
```

## Step 2: Generate PII Encryption Key

```bash
# Generate 32-byte key
openssl rand -base64 32

# Set environment variable
export PII_ENCRYPTION_KEY="<generated_key>"
```

## Step 3: Configure Database

```bash
# Set database password
export DB_PASSWORD="your_postgres_password"
export REDIS_PASSWORD="your_redis_password"
```

## Step 4: Build PoliSync

```bash
cd backend/polisync

# Windows
.\build.ps1

# Linux/Mac
chmod +x build.ps1
pwsh build.ps1
```

## Step 5: Run

```bash
cd src/PoliSync.ApiHost
dotnet run
```

## Step 6: Verify

```bash
# Health check
curl http://localhost:50121/health

# Should return: {"status":"Healthy"}
```

## Available Endpoints

| Service | gRPC | HTTP | Status |
|---------|------|------|--------|
| Product | 50120 | 50121 | 🚧 In Progress |
| Quote | 50130 | 50131 | 🚧 In Progress |
| Order | 50140 | 50141 | 🚧 In Progress |
| Commission | 50150 | 50151 | 🚧 In Progress |
| Policy | 50160 | 50161 | 🚧 In Progress |
| Underwriting | 50170 | 50171 | 🚧 In Progress |
| Claim | 50210 | 50211 | 🚧 In Progress |

## Docker Compose (Alternative)

```bash
cd backend/polisync

# Set environment variables in .env file
cat > .env << EOF
DB_PASSWORD=your_password
REDIS_PASSWORD=your_password
PII_ENCRYPTION_KEY=$(openssl rand -base64 32)
EOF

# Start everything
docker-compose up --build
```

## Testing gRPC Services

### Using grpcurl

```bash
# Install grpcurl
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

# List services
grpcurl -plaintext localhost:50120 list

# Call a method (once implemented)
grpcurl -plaintext -d '{"product_code":"MOTOR-001"}' \
  localhost:50120 products.v1.ProductService/GetProduct
```

### Using BloomRPC

1. Download [BloomRPC](https://github.com/bloomrpc/bloomrpc)
2. Import proto files from `api/protos/`
3. Connect to `localhost:50120`
4. Test RPCs interactively

## Development Workflow

### 1. Make Changes

Edit code in `src/PoliSync.*/`

### 2. Hot Reload

```bash
# Run with hot reload
dotnet watch run --project src/PoliSync.ApiHost
```

### 3. Run Tests

```bash
# All tests
dotnet test

# Specific project
dotnet test tests/PoliSync.Products.Tests
```

### 4. Check Logs

```bash
# Tail logs
tail -f logs/polisync-*.log

# Or use Docker logs
docker logs -f polisync
```

## Common Issues

### Port Already in Use

```bash
# Find process using port
netstat -ano | findstr :50120  # Windows
lsof -i :50120                 # Linux/Mac

# Kill process
taskkill /PID <pid> /F         # Windows
kill -9 <pid>                  # Linux/Mac
```

### Database Connection Failed

```bash
# Check PostgreSQL is running
docker ps | grep postgres

# Test connection
psql -h localhost -U app -d insuretech

# Check schema exists
\dn
# Should show: insurance_schema, commission_schema
```

### Kafka Connection Failed

```bash
# Check Kafka is running
docker ps | grep kafka

# List topics
docker exec -it kafka kafka-topics --bootstrap-server localhost:9092 --list
```

### Build Errors

```bash
# Clean everything
./build.ps1 -Clean

# Restore packages
dotnet restore

# Rebuild
dotnet build
```

## Next Steps

1. **Read** [POLISYNC_REFERENCE.md](POLISYNC_REFERENCE.md) for status, state machines & implementation guide
2. **Review** [ARCHITECTURE.md](ARCHITECTURE.md) for system diagrams
3. **Start** insurance-service (Go :50115) before running PoliSync
4. **Deploy** using Docker Compose or Kubernetes

## Getting Help

- Check logs in `logs/polisync-*.log`
- Review health check: `http://localhost:50121/health`
- Inspect metrics: `http://localhost:9090/metrics`
- Read documentation in `README.md`

## Useful Commands

```bash
# Build
dotnet build

# Run
dotnet run --project src/PoliSync.ApiHost

# Test
dotnet test

# Clean
dotnet clean

# Format code
dotnet format

# Publish
dotnet publish -c Release -o ./publish

# Docker build
docker build -t polisync:latest .

# Docker run
docker run -p 50120-50211:50120-50211 polisync:latest
```

## Environment Variables Reference

| Variable | Description | Example |
|----------|-------------|---------|
| `DB_PASSWORD` | PostgreSQL password | `secure_password` |
| `REDIS_PASSWORD` | Redis password | `redis_password` |
| `PII_ENCRYPTION_KEY` | 32-byte base64 key | `base64_encoded_key` |
| `ASPNETCORE_ENVIRONMENT` | Environment | `Development` / `Production` |
| `Kafka__BootstrapServers` | Kafka brokers | `localhost:9092` |

## Production Checklist

- [ ] Set strong passwords for DB and Redis
- [ ] Generate secure PII encryption key
- [ ] Configure JWT public key path
- [ ] Set up TLS certificates for gRPC
- [ ] Configure log retention
- [ ] Set up monitoring (Prometheus + Grafana)
- [ ] Configure distributed tracing
- [ ] Set up backup strategy
- [ ] Configure auto-scaling (HPA)
- [ ] Set resource limits (CPU/Memory)
- [ ] Enable audit logging
- [ ] Configure rate limiting
- [ ] Set up alerting rules

Happy coding! 🚀
