# Test script for Leaf project using gotestsum (PowerShell)

Write-Host "🌱 Leaf Test Suite" -ForegroundColor Green
Write-Host ""

# Check if gotestsum is installed
$gotestsumExists = Get-Command gotestsum -ErrorAction SilentlyContinue

if (-not $gotestsumExists) {
    Write-Host "⚠️  gotestsum is not installed" -ForegroundColor Yellow
    Write-Host "Installing gotestsum..."
    go install gotest.tools/gotestsum@latest
    Write-Host "✓ gotestsum installed" -ForegroundColor Green
    Write-Host ""
}

# Run tests with gotestsum
Write-Host "Running tests..." -ForegroundColor Green
gotestsum --format pkgname-and-test-fails -- -cover -v ./...

Write-Host ""
Write-Host "✓ Tests completed" -ForegroundColor Green
