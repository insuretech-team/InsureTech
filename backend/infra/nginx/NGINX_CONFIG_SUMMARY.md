# InsureTech Nginx Configuration Summary

## Overview
The nginx configuration for InsureTech is located at `E:\Projects\InsureTech\backend\infra\nginx` with a modular architecture supporting multiple insurance portals and APIs.

## Directory Structure

```
nginx/
├── nginx.conf                          # Main configuration (modular architecture)
├── Dockerfile                          # Multi-stage Docker build with validation
├── docker-compose.nginx.yml            # Production docker-compose
├── docker-compose.test.yml             # Local testing docker-compose
├── mime.types                          # MIME type definitions
├── README.md                           # Documentation
├── .gitignore                          # Git ignore rules

├── conf.d/                             # Global HTTP configurations
│   ├── 00-global.conf                 # Basic HTTP settings, client limits
│   ├── 01-logging.conf                # Log formats (main, extended, JSON, security)
│   ├── 02-performance.conf            # Performance tuning, buffer settings
│   ├── 03-security.conf               # Rate limiting, DDoS protection, exploit blocking
│   └── 04-compression.conf            # Gzip & Brotli compression

├── snippets/                           # Reusable configuration snippets
│   ├── proxy-params.conf              # Standard proxy headers
│   ├── security-headers.conf          # Security headers (CSP, X-Frame, etc.)
│   ├── ssl-params.conf                # SSL/TLS parameters (TLS 1.2/1.3)
│   └── websocket.conf                 # WebSocket upgrade support

├── upstreams/                          # Backend server definitions
│   └── gateway.conf                   # API Gateway upstream

├── cache/                              # Cache configurations
│   ├── cache-zones.conf               # Three-tier cache (static, api, microcache)
│   └── cache-bypass.conf              # Cache bypass rules (auth, POST, etc.)

├── maps/                               # Map definitions
│   └── bot-detection.conf             # Good/bad bot classification

├── sites-available/                    # Virtual host configurations
│   ├── b2b_portal.conf                # B2B Portal (b2b.labaidinsuretech.com)
│   ├── coming_soon.conf               # Coming soon placeholder
│   ├── default.conf                   # Default fallback server
│   ├── insuretech-api.conf            # API endpoint (api.labaidinsuretech.com)
│   ├── labaidinsuretech.com.conf      # Main website
│   └── system_portal.conf             # System/admin portal

├── sites-enabled/                      # Symlinks to enabled sites (empty - enable manually)

├── error-pages/                        # Error page templates
│   ├── 404.html
│   └── 502.html

├── stream.d/                           # TCP/UDP configurations (optional)
│   └── .gitkeep

├── scripts/                            # Maintenance and deployment scripts
│   ├── clear-cache.sh
│   ├── create-dev-branch.ps1
│   ├── enable-sites.ps1
│   ├── enable-sites.sh
│   ├── pre-deployment-check.ps1
│   ├── setup.sh
│   ├── test-config.sh
│   └── validate-config.ps1

└── static/                             # Logo and static assets
    ├── Insuretech-Logo.png
    └── logo_transparent.png
```

## Key Configuration Details

### 1. Main Configuration (`nginx.conf`)
- **User**: `www-data`
- **Worker Processes**: `auto` (scales with CPU cores)
- **Worker Connections**: 4096 per worker
- **Worker File Limit**: 65535
- **Event Model**: `epoll` with `multi_accept`

### 2. Global Settings (`conf.d/00-global.conf`)
- **Max Body Size**: 20MB (for file uploads)
- **Sendfile**: Enabled with 512k chunks
- **Keepalive**: 65s timeout, 1000 requests per connection
- **DNS Resolver**: Google's public DNS (8.8.8.8, 8.8.4.4)

### 3. Logging (`conf.d/01-logging.conf`)
Four log formats available:
- **main**: Standard access log
- **extended**: Includes timing info (request, upstream, cache status)
- **json_combined**: Structured JSON logging
- **security**: Malicious activity logging

Main access log: `/var/log/nginx/access.log` with 32k buffer, 5s flush

### 4. Performance (`conf.d/02-performance.conf`)
- **File Cache**: 10,000 files max, 30s inactive timeout
- **Proxy Buffers**: 256 × 8k
- **FastCGI Buffers**: 256 × 16k
- **Timeouts**: Connect 10s, send 60s, read 60s
- **HTTP/1.1** for upstream keepalive

### 5. Security (`conf.d/03-security.conf`)

**Rate Limiting Zones**:
- `general`: 10 req/s
- `api`: 30 req/s
- `login`: 5 req/min
- `strict`: 1 req/s

**Connection Limits**:
- Per IP: 50 concurrent connections
- Per server: 500 total connections

**Exploit Blocking**:
- Git/SVN/Mercurial exposure (`.git`, `.svn`, `.hg`)
- Environment files (`.env*`)
- PHP exploits (`.php`, `phpinfo`, `eval-stdin`)
- Path traversal (`../`, `%2e%2e`)
- Shell injection commands
- Framework exploits (actuator, swagger, GraphQL, WordPress)
- SQL injection patterns
- XSS patterns
- File inclusion attempts

**Bad Bot Blocking**:
- Empty user agents
- Scan tools (nmap, sqlmap, nikto, masscan, nessus)
- Aggressive crawlers (Bytespider, AhrefsBot, SemrushBot)
- Suspicious patterns (Python HTTP clients, Go HTTP clients)
- Shellshock attempts

**Bad Methods**: TRACE, TRACK, PROPFIND, MKCOL, COPY, MOVE, LOCK, UNLOCK

**Host Blocking**:
- Direct IP access (forces domain-based access)
- All IPv4 and IPv6 literal addresses

### 6. Compression (`conf.d/04-compression.conf`)
- **Gzip**: Level 6, min 256 bytes
- **Types**: JSON, JavaScript, CSS, fonts, SVG, HTML, XML, etc.
- **Brotli**: Available (commented out, uncomment if module available)

### 7. Caching Strategy (`cache/`)

**Three-Tier Cache**:

1. **Static Cache** (`/var/cache/nginx/static`)
   - Size: 1GB
   - Inactive: 7 days
   - Use case: Assets, images

2. **API Cache** (`/var/cache/nginx/api`)
   - Size: 500MB
   - Inactive: 1 hour
   - Use case: API responses

3. **Microcache** (`/var/cache/nginx/microcache`)
   - Size: 100MB
   - Inactive: 10 minutes
   - Use case: Dynamic content

**Cache Bypass Conditions**:
- Non-GET/HEAD requests (POST, PUT, DELETE, PATCH)
- Requests with Authorization header
- Requests with session/token/auth cookies
- Admin, cart, checkout, account URIs

### 8. SSL/TLS (`snippets/ssl-params.conf`)
- **Protocols**: TLSv1.2, TLSv1.3 only
- **Session Cache**: 50MB shared cache, 1-day timeout
- **Session Tickets**: Disabled
- **OCSP Stapling**: Disabled (Firefox compatibility)

### 9. Security Headers (`snippets/security-headers.conf`)
- **X-Frame-Options**: SAMEORIGIN (prevents clickjacking)
- **X-Content-Type-Options**: nosniff (prevents MIME sniffing)
- **X-XSS-Protection**: 1; mode=block
- **Referrer-Policy**: strict-origin-when-cross-origin
- **Permissions-Policy**: Blocks geolocation, microphone, camera, payment
- **CSP**: Allows self, Cloudflare Insights, Google Fonts, CDN
- **HSTS**: Set per HTTPS server (63,072,000 seconds = 2 years)

### 10. Virtual Hosts (`sites-available/`)

#### a) Main Website (`labaidinsuretech.com.conf`)
- HTTP: ACME challenge only, redirects to HTTPS
- HTTPS: Proxies to backend (upstream defined)
- Cache aggressive for `/assets/`, `/static/`

#### b) B2B Portal (`b2b_portal.conf`)
- Domain: `b2b.labaidinsuretech.com`
- Proxies to `127.0.0.1:3000` (Next.js)
- WebSocket support for HMR
- Static assets cached aggressively

#### c) API (`insuretech-api.conf`)
- Domain: `api.labaidinsuretech.com`
- Proxies to API gateway
- No caching for `/auth/**` endpoints
- Rate limiting: 30 req/s

#### d) System Portal (`system_portal.conf`)
- Domain: `system.labaidinsuretech.com`
- Admin/system access

#### e) Coming Soon (`coming_soon.conf`)
- Placeholder for under-maintenance sites

#### f) Default (`default.conf`)
- Fallback for undefined hosts
- Returns 444 (close connection)

### 11. WebSocket Support (`snippets/websocket.conf`)
- HTTP/1.1 upgrade
- Upgrade and Connection headers set
- 7-day timeouts for long-lived connections
- Buffering disabled

### 12. Bot Detection (`maps/bot-detection.conf`)

**Good Bots** (allowed):
- Googlebot, Bingbot, DuckDuckBot, Baiduspider, YandexBot
- Facebook, LinkedIn, Twitter crawlers

**Bot Rate Limiting**: 1 req/s for identified bots

## Docker Setup

### Production (`docker-compose.nginx.yml`)
- **Image**: nginx:1.25-alpine
- **Container**: trendico_nginx
- **Ports**: 80, 443
- **Volumes**: All config directories, SSL certs, cache storage
- **Dependencies**: gateway, trendyco, trendfront services
- **Healthcheck**: HTTP GET to `/health` endpoint

### Testing (`docker-compose.test.yml`)
- **Image**: nginx:1.25 (Debian, not Alpine)
- **Ports**: 8880 (HTTP), 8443 (HTTPS)
- **Self-Signed Certificates**: Generated for all domains
- **Entrypoint**: Generates certs, enables all sites, validates config

## Dockerfile Stages

### Stage 1: Validation
- Creates all required directories
- Copies and validates configuration
- Moves upstreams/sites-enabled aside for base validation
- Checks for MTA-STS configuration

### Stage 2: Production
- Copies validated configs from Stage 1
- Generates default self-signed certificate
- Sets proper permissions (nginx:nginx)
- Exposes ports 80, 443
- Healthcheck via curl to `/health`

## Environment Files

**Search Results**: No `.env*` files found in `E:\Projects\InsureTech\` root

The configuration does NOT use environment variables for settings — all values are hardcoded in configuration files.

## Security Highlights

1. **Comprehensive Exploit Blocking**: Blocks common attack vectors including SQL injection, XSS, path traversal, shell injection
2. **Rate Limiting**: Different limits for general (10/s), API (30/s), login (5/min), and strict (1/s) zones
3. **DDoS Protection**: Connection timeouts, aggressive user agent blocking
4. **No Version Disclosure**: `server_tokens off`
5. **No Direct IP Access**: Blocks requests to server IP
6. **HTTPS Enforcement**: Automatic HTTP→HTTPS redirect
7. **Bot Management**: Differentiates good bots (search engines) from bad (scanners)

## Performance Features

1. **Multi-Tier Caching**: Static (1GB), API (500MB), Microcache (100MB)
2. **Cache Locking**: Prevents thundering herd on cache misses
3. **Background Updates**: Cache can be updated in background
4. **Sendfile**: Zero-copy file transmission
5. **Gzip Compression**: Level 6, min 256 bytes
6. **File Descriptor Caching**: 10,000 files cached

## Deployment Notes

1. **Sites-Enabled**: Currently empty — sites must be symlinked from `sites-available/`
2. **SSL Certificates**: Expected in `/etc/letsencrypt/live/{domain}/` (Certbot managed)
3. **Cache Directories**: Must be created with proper permissions:
   - `/var/cache/nginx/static`
   - `/var/cache/nginx/api`
   - `/var/cache/nginx/microcache`
4. **Upstream Resolution**: Deferred to container startup (DNS lookup timing)
5. **Certbot Integration**: Automatically adds HSTS header for HTTPS

## Configuration Validation

```bash
# Test syntax (on server)
sudo nginx -t

# Reload after changes
sudo nginx -s reload
# or
sudo systemctl reload nginx

# Docker validation
docker compose -f docker-compose.test.yml up
```

## Customization Points

1. **Rate Limits**: Adjust in `conf.d/03-security.conf`
2. **Cache Settings**: Modify `cache/cache-zones.conf`
3. **Security Rules**: Update `conf.d/03-security.conf` and `snippets/security-headers.conf`
4. **Upstream Servers**: Edit `upstreams/gateway.conf`
5. **Log Formats**: Change in `conf.d/01-logging.conf`

## Known Issues/Considerations

1. **OCSP Stapling**: Disabled for Firefox compatibility
2. **Content-Security-Policy**: Partially truncated in file but visible structure is complete
3. **MTA-STS**: Dockerfile checks for MTA-STS config but file not found in sites-available
4. **Sites-Enabled**: All sites must be manually enabled via symlinks
