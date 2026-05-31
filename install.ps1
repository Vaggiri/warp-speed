<#
.SYNOPSIS
Installs warp-speed CLI to the user's system so it can be run globally.

.DESCRIPTION
This script creates a .warp-speed\bin directory in the user's profile,
copies the warp-speed.exe binary into it, and updates the User PATH environment
variable permanently.

.EXAMPLE
.\install.ps1
#>

$ErrorActionPreference = 'Stop'

$binaryName = "warp-speed.exe"
$installDir = Join-Path -Path $env:USERPROFILE -ChildPath ".warp-speed\bin"
$binaryPath = Join-Path -Path $installDir -ChildPath $binaryName

Write-Host "[INFO] Starting installation of warp-speed..." -ForegroundColor Cyan

# 1. Check if binary exists in current directory
if (-Not (Test-Path -Path ".\$binaryName")) {
    Write-Host "[ERROR] Could not find $binaryName in the current directory." -ForegroundColor Red
    Write-Host "Please run 'go build -o $binaryName' first." -ForegroundColor Yellow
    exit 1
}

# 2. Create the installation directory if it doesn't exist
if (-Not (Test-Path -Path $installDir)) {
    Write-Host "[INFO] Creating directory: $installDir" -ForegroundColor Blue
    New-Item -ItemType Directory -Force -Path $installDir | Out-Null
}

# 3. Copy/Move the binary
Write-Host "[INFO] Copying binary to install directory..." -ForegroundColor Blue
Copy-Item -Path ".\$binaryName" -Destination $binaryPath -Force

# 4. Update the User PATH Environment Variable
$userPath = [Environment]::GetEnvironmentVariable("PATH", [EnvironmentVariableTarget]::User)

if ($userPath -like "*$installDir*") {
    Write-Host "[OK] The installation directory is already in your PATH." -ForegroundColor Green
} else {
    Write-Host "[INFO] Adding $installDir to your User PATH..." -ForegroundColor Blue
    $newPath = "$userPath;$installDir"
    [Environment]::SetEnvironmentVariable("PATH", $newPath, [EnvironmentVariableTarget]::User)
    
    Write-Host "[WARN] PATH updated! You will need to restart your terminal for the changes to take effect." -ForegroundColor Yellow
}

Write-Host "[SUCCESS] warp-speed has been successfully installed!" -ForegroundColor Green
Write-Host "Try running 'warp-speed' in a new terminal window." -ForegroundColor White
