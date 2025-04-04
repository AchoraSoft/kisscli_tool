#!/bin/bash

# Configuration
VERSION="1.0.2"
REPO="AchoraSoft/kisscli_tool"
BINARY_NAME="kissc"
INSTALL_DIR="/usr/local/bin"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Detect OS and Architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# Normalize architecture
case $ARCH in
    x86_64) ARCH="amd64" ;;
    arm64|aarch64) ARCH="arm64" ;;
    *) echo -e "${RED}Unsupported architecture: $ARCH${NC}"; exit 1 ;;
esac

# Windows detection (for WSL)
if [[ "$OS" == "mingw"* || "$OS" == "msys"* || "$OS" == "cygwin"* ]]; then
    echo -e "${RED}Please use the Windows installer (install.ps1)${NC}"
    exit 1
fi

# Download URL
URL="https://github.com/$REPO/releases/download/v$VERSION/$BINARY_NAME-$VERSION-$OS-$ARCH"

# Download and install
echo -e "${GREEN}Installing KISSC v$VERSION for $OS-$ARCH...${NC}"
echo -e "Downloading from: $URL"

if ! curl -fSL "$URL" -o "$BINARY_NAME"; then
    echo -e "${RED}Failed to download KISSC${NC}"
    exit 1
fi

chmod +x "$BINARY_NAME"

# Move to install directory
echo -e "Installing to $INSTALL_DIR (may require sudo)"
sudo mv "$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"

# Verify installation
if command -v "$BINARY_NAME" >/dev/null 2>&1; then
    echo -e "${GREEN}KISSC installed successfully!${NC}"
    echo -e "Try running: ${YELLOW}$BINARY_NAME <your_project_name>${NC}"
else
    echo -e "${RED}Installation failed - $BINARY_NAME not found in PATH${NC}"
    exit 1
fi