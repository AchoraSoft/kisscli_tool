# uninstall.ps1
$BINARY_NAME = "kissc"
$INSTALL_DIR = "$env:USERPROFILE\bin"
$INSTALL_PATH = "$INSTALL_DIR\$BINARY_NAME.exe"

if (Test-Path $INSTALL_PATH) {
    Write-Host "Removing KISSC..."
    Remove-Item $INSTALL_PATH -Force
    
    if ($?) {
        Write-Host "KISSC successfully uninstalled" -ForegroundColor Green
    } else {
        Write-Host "Failed to uninstall KISSC" -ForegroundColor Red
        exit 1
    }
} else {
    Write-Host "KISSC is not installed at $INSTALL_PATH"
}