#!/bin/bash
# Nginx Setup Script for Trendico Platform
# Version: 2.0

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
NGINX_DIR="$(dirname "$SCRIPT_DIR")"
INSTALL_DIR="/etc/nginx"

echo "=========================================="
echo "Trendico Nginx Configuration Setup"
echo "=========================================="
echo ""

# Check if running as root
if [ "$EUID" -ne 0 ]; then 
    echo "❌ Please run as root (sudo)"
    exit 1
fi

echo "📋 Pre-flight checks..."

# Backup existing configuration
if [ -d "$INSTALL_DIR" ]; then
    BACKUP_DIR="/etc/nginx.backup.$(date +%Y%m%d-%H%M%S)"
    echo "📦 Backing up existing configuration to $BACKUP_DIR"
    cp -r "$INSTALL_DIR" "$BACKUP_DIR"
fi

# Create directory structure
echo "📁 Creating directory structure..."
mkdir -p "$INSTALL_DIR"/{conf.d,snippets,upstreams,cache,maps,stream.d,includes}
mkdir -p "$INSTALL_DIR"/{sites-available,sites-enabled}
mkdir -p "$INSTALL_DIR"/ssl/{certs,ocsp}

# Create cache directories
echo "💾 Creating cache directories..."
mkdir -p /var/cache/nginx/{static,api,microcache,proxy_temp}
chown -R nginx:nginx /var/cache/nginx

# Create error pages directory
echo "📄 Creating error pages directory..."
mkdir -p /var/www/html
cp "$NGINX_DIR"/error-pages/*.html /var/www/html/ 2>/dev/null || true

# Copy configuration files
echo "📝 Copying configuration files..."
cp "$NGINX_DIR"/nginx.conf "$INSTALL_DIR"/
cp "$NGINX_DIR"/conf.d/*.conf "$INSTALL_DIR"/conf.d/
cp "$NGINX_DIR"/snippets/*.conf "$INSTALL_DIR"/snippets/
cp "$NGINX_DIR"/upstreams/*.conf "$INSTALL_DIR"/upstreams/
cp "$NGINX_DIR"/cache/*.conf "$INSTALL_DIR"/cache/
cp "$NGINX_DIR"/maps/*.conf "$INSTALL_DIR"/maps/ 2>/dev/null || true
cp "$NGINX_DIR"/sites-available/*.conf "$INSTALL_DIR"/sites-available/

# Enable sites
echo "🔗 Enabling sites..."
cd "$INSTALL_DIR"/sites-enabled

# Enable main shop
if [ ! -L trendyco.com.bd.conf ]; then
    ln -s ../sites-available/trendyco.com.bd.conf .
    echo "  ✅ Enabled trendyco.com.bd"
fi

# Enable portal
if [ ! -L portal.trendyco.com.bd.conf ]; then
    ln -s ../sites-available/portal.trendyco.com.bd.conf .
    echo "  ✅ Enabled portal.trendyco.com.bd"
fi

# Enable MTA-STS
if [ ! -L mta-sts.trendyco.com.bd.conf ]; then
    ln -s ../sites-available/mta-sts.trendyco.com.bd.conf .
    echo "  ✅ Enabled mta-sts.trendyco.com.bd"
fi

# Test configuration
echo ""
echo "🧪 Testing nginx configuration..."
nginx -t

if [ $? -eq 0 ]; then
    echo ""
    echo "✅ Configuration test passed!"
    echo ""
    echo "🔄 Reloading nginx..."
    nginx -s reload || systemctl reload nginx
    
    echo ""
    echo "=========================================="
    echo "✅ Setup completed successfully!"
    echo "=========================================="
    echo ""
    echo "Next steps:"
    echo "1. Verify all sites are accessible"
    echo "2. Check logs: tail -f /var/log/nginx/error.log"
    echo "3. Monitor cache: du -sh /var/cache/nginx/*"
    echo ""
else
    echo ""
    echo "❌ Configuration test failed!"
    echo "Please review the errors above and fix them."
    exit 1
fi
