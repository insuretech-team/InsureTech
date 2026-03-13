#!/bin/bash
# One-time setup: grant insureadmin passwordless sudo on remote
# Usage: bash scripts/setup_sudoers.sh
# Will prompt for insureadmin sudo password ONCE.

read -rsp "insureadmin sudo password: " PASS
echo ""

ssh -o StrictHostKeyChecking=no insureadmin@146.190.97.242 \
    "echo '$PASS' | sudo -S bash -c \
    'echo \"insureadmin ALL=(ALL) NOPASSWD:ALL\" > /etc/sudoers.d/insureadmin && \
    chmod 440 /etc/sudoers.d/insureadmin && \
    echo OK: passwordless sudo granted'"
