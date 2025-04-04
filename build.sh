#!/bin/bash

# Initialize Go module if not already done
if [ ! -f go.mod ]; then
    echo "Initializing Go module..."
    go mod init github.com/your-username/kissc
fi

# Download dependencies
echo "Downloading dependencies..."
go mod tidy

VERSION="1.0.2"
BINARY_NAME="kissc"
PLATFORMS=("linux/amd64" "darwin/amd64" "darwin/arm64" "windows/amd64")

# Create builds directory if it doesn't exist
mkdir -p builds

echo "Building binaries..."
for platform in "${PLATFORMS[@]}"
do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    output_name="builds/${BINARY_NAME}-${VERSION}-${GOOS}-${GOARCH}"
    
    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
    fi

    echo "Building for $GOOS/$GOARCH..."
    env GOOS=$GOOS GOARCH=$GOARCH go build -o $output_name .
    
    if [ $? -eq 0 ]; then
        echo "Successfully built $output_name"
    else
        echo "Error building for $GOOS/$GOARCH"
        exit 1
    fi
done

echo "Build complete. Binaries are in the builds/ directory."