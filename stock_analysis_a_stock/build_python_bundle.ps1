# Build Python Bundle with PyInstaller using PowerShell
Write-Host "====================================" -ForegroundColor Cyan
Write-Host "Building Python Bundle with PyInstaller" -ForegroundColor Cyan
Write-Host "====================================" -ForegroundColor Cyan
Write-Host ""

# Get the current directory
$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
Set-Location $scriptDir

# Check if virtual environment exists
$venvPython = Join-Path $scriptDir "venv\Scripts\python.exe"
$venvPyInstaller = Join-Path $scriptDir "venv\Scripts\pyinstaller.exe"

if (Test-Path $venvPython) {
    Write-Host "Using virtual environment..." -ForegroundColor Green
    
    # Install pyinstaller if not already installed
    if (-not (Test-Path $venvPyInstaller)) {
        Write-Host "Installing PyInstaller..." -ForegroundColor Yellow
        & $venvPython -m pip install pyinstaller
    }
    
    # Run PyInstaller
    Write-Host "Running PyInstaller..." -ForegroundColor Yellow
    & $venvPython -m PyInstaller build_pyinstaller.spec --clean
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host ""
        Write-Host "====================================" -ForegroundColor Green
        Write-Host "Python bundle created successfully!" -ForegroundColor Green
        Write-Host "Location: dist\stock_analysis_engine\" -ForegroundColor Green
        Write-Host "====================================" -ForegroundColor Green
    } else {
        Write-Host ""
        Write-Host "====================================" -ForegroundColor Red
        Write-Host "ERROR: Failed to create Python bundle" -ForegroundColor Red
        Write-Host "====================================" -ForegroundColor Red
        exit 1
    }
} else {
    Write-Host "ERROR: Virtual environment not found!" -ForegroundColor Red
    Write-Host "Please run: .\setup_env.bat" -ForegroundColor Yellow
    exit 1
}
