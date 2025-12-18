# YOLOSINT Build Script
# Ensures binary is rebuilt and warns if unchanged

$ErrorActionPreference = "Stop"
$ProjectRoot = Split-Path -Parent (Split-Path -Parent $MyInvocation.MyCommand.Path)
$BinaryPath = Join-Path $ProjectRoot "yolosint.exe"
$HashFile = Join-Path $ProjectRoot ".build_hash"

Set-Location $ProjectRoot

# Get hash of existing binary (if exists)
$OldHash = $null
if (Test-Path $BinaryPath) {
    $OldHash = (Get-FileHash $BinaryPath -Algorithm SHA256).Hash
}

Write-Host "Running go mod tidy..." -ForegroundColor Cyan
go mod tidy
if ($LASTEXITCODE -ne 0) {
    Write-Host "ERROR: go mod tidy failed" -ForegroundColor Red
    exit 1
}

Write-Host "Building yolosint.exe..." -ForegroundColor Cyan
go build -o yolosint.exe ./cmd/yolosint
if ($LASTEXITCODE -ne 0) {
    Write-Host "ERROR: Build failed" -ForegroundColor Red
    exit 1
}

# Get hash of new binary
$NewHash = (Get-FileHash $BinaryPath -Algorithm SHA256).Hash

# Compare hashes
if ($OldHash -eq $NewHash) {
    Write-Host ""
    Write-Host "WARNING: Binary unchanged! No code changes detected." -ForegroundColor Yellow
    Write-Host "Hash: $NewHash" -ForegroundColor Yellow
    Write-Host ""
} else {
    Write-Host ""
    Write-Host "SUCCESS: Binary rebuilt with changes." -ForegroundColor Green
    Write-Host "Old: $OldHash" -ForegroundColor DarkGray
    Write-Host "New: $NewHash" -ForegroundColor Green
    Write-Host ""
}

# Save hash for reference
$NewHash | Out-File -FilePath $HashFile -Encoding UTF8

Write-Host "Binary ready: $BinaryPath" -ForegroundColor Cyan
Write-Host "Launching yolosint.exe..." -ForegroundColor Cyan
Write-Host ""

# Run the binary
& $BinaryPath

