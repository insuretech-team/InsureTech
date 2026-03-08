# Nginx Dockerfile for Trendico
# Production-ready nginx with modular configuration

# ============================================
# Stage 1: Validation Stage
# ============================================
FROM nginx:1.27-alpine AS validate

LABEL maintainer="Trendico DevOps"
LABEL description="Nginx configuration validation stage"

# Install required packages for validation
RUN apk add --no-cache bash curl

# Create required directories
RUN mkdir -p /etc/nginx/conf.d \
    && mkdir -p /etc/nginx/snippets \
    && mkdir -p /etc/nginx/upstreams \
    && mkdir -p /etc/nginx/cache \
    && mkdir -p /etc/nginx/maps \
    && mkdir -p /etc/nginx/sites-available \
    && mkdir -p /etc/nginx/sites-enabled \
    && mkdir -p /etc/nginx/stream.d \
    && mkdir -p /etc/nginx/ssl/certs \
    && mkdir -p /etc/nginx/ssl/ocsp \
    && mkdir -p /var/cache/nginx/static \
    && mkdir -p /var/cache/nginx/api \
    && mkdir -p /var/cache/nginx/microcache \
    && mkdir -p /var/www/html \
    && mkdir -p /var/www/certbot

# Remove default config
RUN rm -f /etc/nginx/conf.d/default.conf

# Copy all nginx configuration files and directories comprehensively
# This ensures sites-available, sites-enabled, and all other subdirectories are included
COPY infra/nginx/mime.types /etc/nginx/mime.types
COPY infra/nginx/nginx.conf /etc/nginx/nginx.conf

# Copy all configuration subdirectories
COPY infra/nginx/conf.d/ /etc/nginx/conf.d/
COPY infra/nginx/snippets/ /etc/nginx/snippets/
COPY infra/nginx/upstreams/ /etc/nginx/upstreams/
COPY infra/nginx/cache/ /etc/nginx/cache/
COPY infra/nginx/maps/ /etc/nginx/maps/
COPY infra/nginx/error-pages/ /var/www/html/

# Copy sites-available AND sites-enabled to ensure all site configurations are included
# Both directories are copied to preserve any enabled site configurations
COPY infra/nginx/sites-available/ /etc/nginx/sites-available/
COPY infra/nginx/sites-enabled/ /etc/nginx/sites-enabled/

# Verify all required site configurations are present
RUN echo "=== Verifying nginx site configurations ===" && \
    ls -la /etc/nginx/sites-available/ && \
    echo "=== Sites enabled ===" && \
    ls -la /etc/nginx/sites-enabled/ && \
    echo "=== Checking for MTA-STS configuration ===" && \
    if [ -f /etc/nginx/sites-enabled/mta-sts.trendyco.com.bd.conf ]; then \
        echo "✓ MTA-STS configuration found in sites-enabled"; \
    else \
        echo "⚠ WARNING: MTA-STS configuration NOT found in sites-enabled"; \
    fi && \
    if [ -f /etc/nginx/sites-available/mta-sts.trendyco.com.bd.conf ]; then \
        echo "✓ MTA-STS configuration found in sites-available"; \
    else \
        echo "⚠ WARNING: MTA-STS configuration NOT found in sites-available"; \
    fi

# Copy scripts directory (if needed for maintenance/debugging)
COPY infra/nginx/scripts/ /usr/local/bin/nginx-scripts/

# Validate nginx configuration syntax (skip upstream DNS resolution during build)
# Temporarily move files that reference upstreams, validate base config, then restore
RUN echo "=== Validating Nginx Base Configuration ===" \
    && mv /etc/nginx/upstreams /etc/nginx/upstreams.bak \
    && mv /etc/nginx/sites-enabled /etc/nginx/sites-enabled.bak \
    && mkdir -p /etc/nginx/upstreams \
    && mkdir -p /etc/nginx/sites-enabled \
    && nginx -t \
    && rm -rf /etc/nginx/upstreams /etc/nginx/sites-enabled \
    && mv /etc/nginx/upstreams.bak /etc/nginx/upstreams \
    && mv /etc/nginx/sites-enabled.bak /etc/nginx/sites-enabled \
    && echo "=== Base Configuration Validated ===" \
    && echo "Note: Full validation with upstreams will happen at container startup"

# ============================================
# Stage 2: Production Stage
# ============================================
FROM nginx:1.27-alpine

LABEL maintainer="Trendico DevOps"
LABEL description="Nginx reverse proxy with modular configuration"
LABEL version="2.0"

# Install required packages
RUN apk add --no-cache \
    bash \
    curl \
    tzdata \
    ca-certificates

# Set timezone
ENV TZ=Asia/Dhaka
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# Create required directories (keep default nginx files)
RUN mkdir -p /etc/nginx/conf.d \
    && mkdir -p /etc/nginx/snippets \
    && mkdir -p /etc/nginx/upstreams \
    && mkdir -p /etc/nginx/cache \
    && mkdir -p /etc/nginx/maps \
    && mkdir -p /etc/nginx/sites-available \
    && mkdir -p /etc/nginx/sites-enabled \
    && mkdir -p /etc/nginx/stream.d \
    && mkdir -p /etc/nginx/ssl/certs \
    && mkdir -p /etc/nginx/ssl/ocsp \
    && mkdir -p /var/cache/nginx/static \
    && mkdir -p /var/cache/nginx/api \
    && mkdir -p /var/cache/nginx/microcache \
    && mkdir -p /var/www/html \
    && mkdir -p /var/www/certbot

# Remove only the default config files we'll replace
RUN rm -f /etc/nginx/conf.d/default.conf

# Copy all validated configuration from validate stage
COPY --from=validate /etc/nginx/mime.types /etc/nginx/mime.types
COPY --from=validate /etc/nginx/nginx.conf /etc/nginx/nginx.conf

# Copy all configuration subdirectories from validate stage
COPY --from=validate /etc/nginx/conf.d/ /etc/nginx/conf.d/
COPY --from=validate /etc/nginx/snippets/ /etc/nginx/snippets/
COPY --from=validate /etc/nginx/upstreams/ /etc/nginx/upstreams/
COPY --from=validate /etc/nginx/cache/ /etc/nginx/cache/
COPY --from=validate /etc/nginx/maps/ /etc/nginx/maps/
COPY --from=validate /etc/nginx/sites-available/ /etc/nginx/sites-available/
COPY --from=validate /etc/nginx/sites-enabled/ /etc/nginx/sites-enabled/
COPY --from=validate /var/www/html/ /var/www/html/
COPY --from=validate /usr/local/bin/nginx-scripts/ /usr/local/bin/nginx-scripts/

# Verify configurations were copied correctly
RUN echo "=== Production stage: Verifying nginx configurations ===" && \
    echo "Sites available:" && ls -la /etc/nginx/sites-available/ && \
    echo "Sites enabled:" && ls -la /etc/nginx/sites-enabled/ && \
    echo "=== Checking for required site configurations ===" && \
    for site in trendyco.com.bd portal.trendyco.com.bd mta-sts.trendyco.com.bd; do \
        if [ -f "/etc/nginx/sites-enabled/${site}.conf" ]; then \
            echo "✓ ${site}.conf found in sites-enabled"; \
        else \
            echo "⚠ WARNING: ${site}.conf NOT found in sites-enabled"; \
        fi; \
    done

# Generate self-signed SSL certificate for default server
RUN apk add --no-cache openssl \
    && openssl req -x509 -nodes -days 3650 -newkey rsa:2048 \
    -keyout /etc/nginx/ssl/certs/default.key \
    -out /etc/nginx/ssl/certs/default.crt \
    -subj "/C=BD/ST=Dhaka/L=Dhaka/O=Trendico/CN=default.local" \
    && chmod 644 /etc/nginx/ssl/certs/default.crt \
    && chmod 600 /etc/nginx/ssl/certs/default.key

# Set proper permissions
RUN chown -R nginx:nginx /var/cache/nginx \
    && chown -R nginx:nginx /var/www/html \
    && chmod -R 755 /var/cache/nginx \
    && chmod -R 755 /var/www/html

# Note: Configuration validation with upstream resolution happens at container startup
# when actual services are available in the network

# Expose ports
EXPOSE 80 443

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=40s --retries=3 \
    CMD curl -f http://localhost/health || exit 1

# Start nginx in foreground
CMD ["nginx", "-g", "daemon off;"]
