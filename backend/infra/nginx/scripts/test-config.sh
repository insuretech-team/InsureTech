#!/bin/bash
# Test Nginx Configuration
# Validates all configurations and checks connectivity

set -e

echo "=========================================="
echo "Nginx Configuration Test"
echo "=========================================="
echo ""

# Test nginx syntax
echo "1️⃣ Testing nginx configuration syntax..."
nginx -t
echo "   ✅ Syntax check passed"
echo ""

# Check if nginx is running
echo "2️⃣ Checking nginx status..."
if systemctl is-active --quiet nginx; then
    echo "   ✅ Nginx is running"
else
    echo "   ⚠️  Nginx is not running"
fi
echo ""

# Test upstream connectivity
echo "3️⃣ Testing upstream connectivity..."

# Gateway
if curl -sf http://localhost:8080/health > /dev/null 2>&1; then
    echo "   ✅ Gateway (8080) is reachable"
else
    echo "   ❌ Gateway (8080) is NOT reachable"
fi

# Trendyco
if curl -sf http://localhost:3000 > /dev/null 2>&1; then
    echo "   ✅ Trendyco (3000) is reachable"
else
    echo "   ❌ Trendyco (3000) is NOT reachable"
fi

# Trendfront
if curl -sf http://localhost:3001 > /dev/null 2>&1; then
    echo "   ✅ Trendfront (3001) is reachable"
else
    echo "   ❌ Trendfront (3001) is NOT reachable"
fi

echo ""

# Check SSL certificates
echo "4️⃣ Checking SSL certificates..."
if [ -f /etc/letsencrypt/live/trendyco.com.bd/fullchain.pem ]; then
    EXPIRY=$(openssl x509 -enddate -noout -in /etc/letsencrypt/live/trendyco.com.bd/fullchain.pem | cut -d= -f2)
    echo "   ✅ Certificate found"
    echo "   📅 Expires: $EXPIRY"
else
    echo "   ⚠️  Certificate not found (OK for local dev)"
fi

echo ""
echo "=========================================="
echo "Test completed"
echo "=========================================="
