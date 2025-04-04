# install.ps1
$VERSION = "1.0.0"
$REPO = "AchoraSoft/kisscli_tool"
$BINARY_NAME = "kissc"
$INSTALL_DIR = "$HOME\bin"

# Create install directory if it doesn't exist
if (-not (Test-Path $INSTALL_DIR)) {
    New-Item -ItemType Directory -Path $INSTALL_DIR | Out-Null
}

# Download URL
$URL = "https://github.com/$REPO/releases/download/v$VERSION/${BINARY_NAME}-${VERSION}-windows-amd64.exe"

# Download and install
Write-Host "Installing KISSC v$VERSION for Windows..."
Write-Host "Downloading from: $URL"

try {
    Invoke-WebRequest -Uri $URL -OutFile "$INSTALL_DIR\$BINARY_NAME.exe"
    $env:PATH += ";$INSTALL_DIR"
    [System.Environment]::SetEnvironmentVariable("PATH", $env:PATH, [System.EnvironmentVariableTarget]::User)
    
    Write-Host "KISSC installed successfully!" -ForegroundColor Green
    Write-Host "Try running: $BINARY_NAME --help" -ForegroundColor Yellow
}
catch {
    Write-Host "Installation failed: $_" -ForegroundColor Red
    exit 1
}