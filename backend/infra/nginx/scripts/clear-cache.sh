#!/bin/bash
# Clear Nginx Cache
# Clears all or specific cache zones

CACHE_DIR="/var/cache/nginx"

echo "=========================================="
echo "Nginx Cache Management"
echo "=========================================="
echo ""

# Check if running as root
if [ "$EUID" -ne 0 ]; then 
    echo "❌ Please run as root (sudo)"
    exit 1
fi

# Show current cache sizes
echo "📊 Current cache sizes:"
du -sh $CACHE_DIR/* 2>/dev/null || echo "No cache directories found"
echo ""

# Ask what to clear
echo "What would you like to clear?"
echo "1) All caches"
echo "2) Static cache only"
echo "3) API cache only"
echo "4) Microcache only"
echo "5) Cancel"
echo ""
read -p "Enter choice [1-5]: " choice

case $choice in
    1)
        echo "🗑️  Clearing all caches..."
        rm -rf $CACHE_DIR/static/*
        rm -rf $CACHE_DIR/api/*
        rm -rf $CACHE_DIR/microcache/*
        echo "✅ All caches cleared"
        ;;
    2)
        echo "🗑️  Clearing static cache..."
        rm -rf $CACHE_DIR/static/*
        echo "✅ Static cache cleared"
        ;;
    3)
        echo "🗑️  Clearing API cache..."
        rm -rf $CACHE_DIR/api/*
        echo "✅ API cache cleared"
        ;;
    4)
        echo "🗑️  Clearing microcache..."
        rm -rf $CACHE_DIR/microcache/*
        echo "✅ Microcache cleared"
        ;;
    5)
        echo "❌ Cancelled"
        exit 0
        ;;
    *)
        echo "❌ Invalid choice"
        exit 1
        ;;
esac

echo ""
echo "🔄 Reloading nginx..."
nginx -s reload

echo ""
echo "✅ Cache cleared and nginx reloaded"
