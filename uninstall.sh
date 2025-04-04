#!/bin/bash

BINARY_NAME="kissc"
INSTALL_DIR="/usr/local/bin"
INSTALL_PATH="$INSTALL_DIR/$BINARY_NAME"

if [ -f "$INSTALL_PATH" ]; then
    echo "Removing KISSC..."
    sudo rm "$INSTALL_PATH"
    
    if [ $? -eq 0 ]; then
        echo "KISSC successfully uninstalled"
    else
        echo "Failed to uninstall KISSC"
        exit 1
    fi
else
    echo "KISSC is not installed at $INSTALL_PATH"
fi