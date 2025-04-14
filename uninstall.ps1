# uninstall.ps1
$BINARY_NAME = "tofo"
$INSTALL_DIR = "$env:USERPROFILE\bin"
$INSTALL_PATH = "$INSTALL_DIR\$BINARY_NAME.exe"

if (Test-Path $INSTALL_PATH) {
    Write-Host "Removing TOFO..."
    Remove-Item $INSTALL_PATH -Force
    
    if ($?) {
        Write-Host "TOFO successfully uninstalled" -ForegroundColor Green
    } else {
        Write-Host "Failed to uninstall TOFO" -ForegroundColor Red
        exit 1
    }
} else {
    Write-Host "TOFO is not installed at $INSTALL_PATH"
}