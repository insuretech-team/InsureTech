#!/bin/bash
# Generate Authn JWT Keys for InsureTech
# Outputs to backend/inscore/secrets

set -e

OUTPUT_DIR="backend/inscore/secrets"
KEY_SIZE=${1:-2048}

if [ "$KEY_SIZE" -lt 2048 ]; then
    echo "Error: KeySize must be >= 2048." >&2
    exit 1
fi

mkdir -p "$OUTPUT_DIR"

PRIVATE_KEY="$OUTPUT_DIR/jwt_rsa_private.pem"
PUBLIC_KEY="$OUTPUT_DIR/jwt_rsa_public.pem"

echo "Generating RSA private key ($KEY_SIZE bits)..."
openssl genrsa -out "$PRIVATE_KEY" "$KEY_SIZE"

echo "Generating RSA public key..."
openssl rsa -in "$PRIVATE_KEY" -pubout -out "$PUBLIC_KEY"

echo ""
echo "Generated:"
echo "  $PRIVATE_KEY"
echo "  $PUBLIC_KEY"
