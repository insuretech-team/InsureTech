#!/bin/bash
# InsureTech Server Setup Script
# Run this as root to install all required software for the InsureTech Backend

set -e

echo "==============================================="
echo " Starting InsureTech Server Provisioning Phase"
echo "==============================================="

# Update the system
echo "-> Updating system packages..."
apt update && apt upgrade -y

# Install prerequisite packages
echo "-> Installing prerequisites (curl, wget, git, build-essential)..."
apt install -y curl wget git build-essential unzip software-properties-common

# Install Docker & Docker Compose
echo "-> Installing Docker..."
if ! command -v docker &> /dev/null; then
    # Add Docker's official GPG key:
    apt-get install -y ca-certificates
    install -m 0755 -d /etc/apt/keyrings
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
    chmod a+r /etc/apt/keyrings/docker.asc

    # Add the repository to Apt sources:
    echo \
      "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \
      $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
      tee /etc/apt/sources.list.d/docker.list > /dev/null
    apt-get update

    apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
else
    echo "Docker is already installed."
fi

# Add our user to the docker group
echo "-> Adding 'insureadmin' to the docker group..."
usermod -aG docker insureadmin || true

# Install Go 1.24 (Required by Gateway, Auth, etc)
echo "-> Installing Go 1.24..."
if ! command -v go &> /dev/null; then
    wget https://go.dev/dl/go1.24.0.linux-amd64.tar.gz
    rm -rf /usr/local/go
    tar -C /usr/local -xzf go1.24.0.linux-amd64.tar.gz
    rm go1.24.0.linux-amd64.tar.gz
    
    # Add to system profile
    echo "export PATH=$PATH:/usr/local/go/bin" >> /etc/profile
    echo "export PATH=$PATH:/usr/local/go/bin:/home/insureadmin/go/bin" >> /home/insureadmin/.bashrc
else
    echo "Go is already installed."
fi

# Install Node.js 20 (Required by Payment & Ticketing service)
echo "-> Installing Node.js 20..."
if ! command -v node &> /dev/null; then
    curl -fsSL https://deb.nodesource.com/setup_20.x | bash -
    apt-get install -y nodejs
else
    echo "Node.js is already installed."
fi

# Install Python 3.11 & pip (Required for AI and OCR)
echo "-> Installing Python..."
apt install -y python3-pip python3-venv

echo "==============================================="
echo " Provisioning Complete!"
echo " It's recommended to log out and log back in as 'insureadmin' to apply Docker group changes."
echo " Use 'docker compose' instead of 'docker-compose'."
echo "==============================================="
