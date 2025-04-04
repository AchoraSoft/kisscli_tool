#!/bin/bash

set -euo pipefail

# Configuration
VERSION="1.0.0"
BINARY_NAME="kissc"
REPO_URL="https://github.com/AchoraSoft/kisscli_tool"
INSTALL_DIR="/usr/local/bin"
CDN_BASE="https://your-cdn.example.com/kissc"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print error and exit
error_exit() {
    echo -e "${RED}Error: $1${NC}" >&2
    exit 1
}

# Function to check for dependencies
check_deps() {
    # Check for curl
    if ! command -v curl &> /dev/null; then
        error_exit "curl is required but not installed. Please install curl first."
    fi
    
    # Check for Deno (but don't fail - we'll install it)
    if ! command -v deno &> /dev/null; then
        echo -e "${YELLOW}Deno not found. It will be installed automatically.${NC}"
    fi
}

# Function to install Deno
install_deno() {
    echo -e "${GREEN}Installing Deno...${NC}"
    
    if [[ "$OS" == "windows" ]]; then
        # Windows installation
        powershell -Command "irm https://deno.land/install.ps1 | iex"
        export PATH="$HOME/.deno/bin:$PATH"
    else
        # Unix-like installation
        curl -fsSL https://deno.land/x/install/install.sh | sh
        export PATH="$HOME/.deno/bin:$PATH"
    fi
    
    # Verify installation
    if ! command -v deno &> /dev/null; then
        echo -e "${YELLOW}Deno installed but not in PATH.${NC}"
        echo -e "Please add this to your shell configuration:"
        echo -e "export PATH=\"\$HOME/.deno/bin:\$PATH\""
        echo -e "Then run this script again."
        exit 1
    fi
    
    echo -e "${GREEN}Deno installed successfully!${NC}"
}

# Detect OS and Architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# Normalize architecture names
case $ARCH in
    x86_64)
        ARCH="amd64"
        ;;
    arm64|aarch64)
        ARCH="arm64"
        ;;
    *)
        error_exit "Unsupported architecture: $ARCH"
        ;;
esac

# Windows detection
if [[ "$OS" == "mingw"* || "$OS" == "msys"* || "$OS" == "cygwin"* ]]; then
    OS="windows"
    EXT=".exe"
    INSTALL_DIR="$HOME/bin"
else
    EXT=""
fi

# Check dependencies
check_deps

# Install Deno if needed
if ! command -v deno &> /dev/null; then
    install_deno
fi

# Download URL
DOWNLOAD_URL="$CDN_BASE/v${VERSION}/${BINARY_NAME}-${VERSION}-${OS}-${ARCH}${EXT}"

# Download and install
echo -e "${GREEN}Installing KISSC v${VERSION} for ${OS}-${ARCH}...${NC}"

# Create install directory if it doesn't exist
mkdir -p "$INSTALL_DIR"

# Download the binary
if ! curl -fSL "$DOWNLOAD_URL" -o "${INSTALL_DIR}/${BINARY_NAME}${EXT}"; then
    error_exit "Failed to download KISSC. Please try again."
fi

# Set executable permissions
chmod +x "${INSTALL_DIR}/${BINARY_NAME}${EXT}"

# Verify installation
if "${INSTALL_DIR}/${BINARY_NAME}${EXT}" --version &> /dev/null; then
    echo -e "${GREEN}KISSC installed successfully!${NC}"
else
    error_exit "Installation verification failed."
fi

# Print success message
echo -e "\n${GREEN}🚀 KISSC is ready to use!${NC}"
echo -e "Try creating your first project:"
echo -e "  ${YELLOW}kissc create my-project${NC}"
echo -e "\nFor more options:"
echo -e "  ${YELLOW}kissc --help${NC}"