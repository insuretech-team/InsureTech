#!/bin/bash

# Load environment variables from .env file and run database test
echo "=== PoliSync Database Live Test ==="
echo "Loading environment variables from .env file..."
echo ""

# Load .env file
ENV_FILE="../../../.env"
if [ -f "$ENV_FILE" ]; then
    export $(grep -v '^#' $ENV_FILE | xargs)
    echo "✅ Environment variables loaded"
else
    echo "❌ .env file not found at: $ENV_FILE"
    exit 1
fi

# Display connection info (without password)
echo ""
echo "Database Connection:"
echo "  Host: $PGHOST"
echo "  Port: $PGPORT"
echo "  Database: $PGDATABASE"
echo "  User: $PGUSER"
echo "  SSL Mode: $PGSSLMODE"
echo ""

# Build and run the test
echo "Building test project..."
dotnet build --configuration Release

if [ $? -eq 0 ]; then
    echo ""
    echo "Running database tests..."
    echo ""
    dotnet run --configuration Release --no-build
else
    echo ""
    echo "❌ Build failed!"
    exit 1
fi
