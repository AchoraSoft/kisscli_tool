#!/bin/bash

BINARY_NAME="tofo"
INSTALL_DIR="/usr/local/bin"
INSTALL_PATH="$INSTALL_DIR/$BINARY_NAME"

if [ -f "$INSTALL_PATH" ]; then
    echo "Removing TOFO..."
    sudo rm "$INSTALL_PATH"
    
    if [ $? -eq 0 ]; then
        echo "TOFO successfully uninstalled"
    else
        echo "Failed to uninstall TOFO"
        exit 1
    fi
else
    echo "TOFO is not installed at $INSTALL_PATH"
fi