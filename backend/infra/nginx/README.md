# Trendico Nginx Configuration

## Modular Architecture

This directory contains the modular nginx configuration for the Trendico platform.

### Directory Structure

```
nginx/
├── nginx.conf                  # Main configuration file
├── conf.d/                     # Global HTTP configurations
│   ├── 00-global.conf         # Basic HTTP settings
│   ├── 01-logging.conf        # Log formats and output
│   ├── 02-performance.conf    # Performance tuning
│   ├── 03-security.conf       # Security settings
│   └── 04-compression.conf    # Gzip compression
├── snippets/                   # Reusable configuration snippets
│   ├── ssl-params.conf        # SSL/TLS parameters
│   ├── proxy-params.conf      # Proxy headers
│   ├── security-headers.conf  # Security headers
│   └── websocket.conf         # WebSocket support
├── upstreams/                  # Backend server definitions
│   ├── gateway.conf           # API Gateway upstream
│   ├── trendyco.conf          # Trendyco app upstream
│   └── trendfront.conf        # Trendfront app upstream
├── cache/                      # Cache configurations
│   ├── cache-zones.conf       # Cache zone definitions
│   └── cache-bypass.conf      # Cache bypass rules
├── sites-available/            # Virtual host configurations
│   ├── trendyco.com.bd.conf
│   ├── portal.trendyco.com.bd.conf
│   ├── mta-sts.trendyco.com.bd.conf
│   └── default.conf
├── sites-enabled/              # Symlinks to enabled sites
├── maps/                       # Map definitions
│   └── bot-detection.conf
├── stream.d/                   # TCP/UDP configurations
└── ssl/                        # SSL certificates
    ├── certs/
    └── ocsp/
```

## Deployment

### 1. Copy Configuration to Nginx

```bash
# On production server
sudo cp -r /path/to/trendico/infra/nginx/* /etc/nginx/

# Create cache directories
sudo mkdir -p /var/cache/nginx/{static,api,microcache}
sudo chown -R nginx:nginx /var/cache/nginx

# Create error pages directory
sudo mkdir -p /var/www/html
sudo cp error-pages/*.html /var/www/html/
```

### 2. Enable Sites

```bash
cd /etc/nginx/sites-enabled

# Enable main shop
sudo ln -s ../sites-available/trendyco.com.bd.conf .

# Enable portal
sudo ln -s ../sites-available/portal.trendyco.com.bd.conf .

# Enable MTA-STS
sudo ln -s ../sites-available/mta-sts.trendyco.com.bd.conf .
```

### 3. Test Configuration

```bash
sudo nginx -t
```

### 4. Reload Nginx

```bash
sudo nginx -s reload
# or
sudo systemctl reload nginx
```

## Docker Deployment

For Docker environments, mount this directory:

```yaml
services:
  nginx:
    image: nginx:latest
    volumes:
      - ./infra/nginx/nginx.conf:/etc/nginx/nginx.conf:ro
      - ./infra/nginx/conf.d:/etc/nginx/conf.d:ro
      - ./infra/nginx/snippets:/etc/nginx/snippets:ro
      - ./infra/nginx/upstreams:/etc/nginx/upstreams:ro
      - ./infra/nginx/cache:/etc/nginx/cache:ro
      - ./infra/nginx/maps:/etc/nginx/maps:ro
      - ./infra/nginx/sites-enabled:/etc/nginx/sites-enabled:ro
      - nginx_cache:/var/cache/nginx
```

## Customization

### Adding New Upstream Servers

Edit the appropriate upstream file in `upstreams/`:

```nginx
upstream trendyco_app {
    least_conn;
    server trendyco:3000 weight=1;
    server trendyco:3002 weight=1;  # Add new instance
}
```

### Adjusting Cache Settings

Modify `cache/cache-zones.conf` to adjust cache sizes and TTLs.

### Security Tuning

Edit `conf.d/03-security.conf` to adjust rate limits based on traffic patterns.

## Monitoring

Check logs:
```bash
tail -f /var/log/nginx/access.log
tail -f /var/log/nginx/error.log
tail -f /var/log/nginx/trendyco-access.log
tail -f /var/log/nginx/portal-access.log
```

## Troubleshooting

### Configuration Test Failed

```bash
sudo nginx -t
# Review error messages and fix configuration
```

### Backend Connection Refused

```bash
# Check backend services are running
docker ps
# or
curl http://localhost:3000
```

### Cache Issues

```bash
# Clear all caches
sudo rm -rf /var/cache/nginx/*
sudo nginx -s reload
```

## Performance Tuning

The configuration is optimized for:
- 4 CPU cores (worker_processes auto)
- 4096 connections per worker
- Multi-tier caching strategy
- HTTP/2 support
- Compression enabled

Adjust `conf.d/00-global.conf` and `conf.d/02-performance.conf` based on your server specifications.
