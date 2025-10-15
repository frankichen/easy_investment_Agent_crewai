# ============================================================================
# A股智能分析系统 - Windows 完整构建脚本 (PowerShell)
# ============================================================================
# 本脚本将自动完成：
# 1. 检查构建环境
# 2. 构建 Python 分析引擎 (PyInstaller)
# 3. 构建 Go UI 程序
# 4. 创建完整的发布包
# ============================================================================

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host " A股智能分析系统 - 完整构建" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# 获取脚本根目录
$rootDir = Split-Path -Parent $MyInvocation.MyCommand.Path

# ============================================================================
# 步骤 1: 检查构建环境
# ============================================================================
Write-Host "[1/5] 检查构建环境..." -ForegroundColor Yellow

# 检查 Python
$pythonCmd = Get-Command python -ErrorAction SilentlyContinue
if (-not $pythonCmd) {
    Write-Host "ERROR: Python 未找到，请先安装 Python 3.12+" -ForegroundColor Red
    exit 1
}

# 检查 Go
$goCmd = Get-Command go -ErrorAction SilentlyContinue
if (-not $goCmd) {
    Write-Host "ERROR: Go 未找到，请先安装 Go 1.24+" -ForegroundColor Red
    exit 1
}

Write-Host "✓ Python 和 Go 环境已就绪" -ForegroundColor Green

# ============================================================================
# 步骤 2: 构建 Python 分析引擎
# ============================================================================
Write-Host ""
Write-Host "[2/5] 构建 Python 分析引擎..." -ForegroundColor Yellow

$stockAnalysisDir = Join-Path $rootDir "stock_analysis_a_stock"
Set-Location $stockAnalysisDir

# 检查虚拟环境
$venvPython = Join-Path $stockAnalysisDir "venv\Scripts\python.exe"
if (-not (Test-Path $venvPython)) {
    Write-Host "ERROR: 虚拟环境不存在！" -ForegroundColor Red
    Write-Host "请先运行: cd stock_analysis_a_stock && setup_env.bat" -ForegroundColor Yellow
    exit 1
}

Write-Host "  → 使用虚拟环境: $venvPython" -ForegroundColor Gray

# 运行 PyInstaller
Write-Host "  → 运行 PyInstaller (这可能需要几分钟)..." -ForegroundColor Gray
& $venvPython -m PyInstaller build_pyinstaller.spec --clean

if ($LASTEXITCODE -ne 0) {
    Write-Host "ERROR: Python 引擎构建失败" -ForegroundColor Red
    Set-Location $rootDir
    exit 1
}

# 检查构建输出
$enginePath = Join-Path $stockAnalysisDir "dist\stock_analysis_engine\stock_analysis_engine.exe"
if (-not (Test-Path $enginePath)) {
    Write-Host "ERROR: Python 引擎可执行文件未找到" -ForegroundColor Red
    Set-Location $rootDir
    exit 1
}

Write-Host "✓ Python 引擎构建完成" -ForegroundColor Green
Set-Location $rootDir

# ============================================================================
# 步骤 3: 准备打包资源
# ============================================================================
Write-Host ""
Write-Host "[3/5] 准备打包资源..." -ForegroundColor Yellow

$uiDir = Join-Path $rootDir "ui"
$uiPythonBundleDir = Join-Path $uiDir "python_bundle"

# 创建目录
if (-not (Test-Path $uiPythonBundleDir)) {
    New-Item -ItemType Directory -Path $uiPythonBundleDir -Force | Out-Null
}

# 复制 Python 引擎
$sourceEngine = Join-Path $stockAnalysisDir "dist\stock_analysis_engine"
$targetEngine = Join-Path $uiPythonBundleDir "stock_analysis_engine"

Write-Host "  → 复制 Python 引擎到 UI 目录..." -ForegroundColor Gray
if (Test-Path $targetEngine) {
    Remove-Item -Recurse -Force $targetEngine
}
Copy-Item -Recurse -Force $sourceEngine $targetEngine

if ($LASTEXITCODE -ne 0) {
    Write-Host "ERROR: 复制 Python 引擎失败" -ForegroundColor Red
    exit 1
}

Write-Host "✓ Python 引擎已复制到 UI 目录" -ForegroundColor Green

# ============================================================================
# 步骤 4: 构建 Go UI
# ============================================================================
Write-Host ""
Write-Host "[4/5] 构建 Go UI 程序..." -ForegroundColor Yellow

Set-Location $uiDir

Write-Host "  → 编译 Go 程序..." -ForegroundColor Gray
go build -ldflags="-s -w" -o stock-analysis-ui-windows-amd64.exe

if ($LASTEXITCODE -ne 0) {
    Write-Host "ERROR: Go 程序构建失败" -ForegroundColor Red
    Set-Location $rootDir
    exit 1
}

Write-Host "✓ Go UI 程序构建完成" -ForegroundColor Green
Set-Location $rootDir

# ============================================================================
# 步骤 5: 创建发布包
# ============================================================================
Write-Host ""
Write-Host "[5/5] 创建发布包..." -ForegroundColor Yellow

$releaseDir = Join-Path $rootDir "release"
$releaseAppDir = Join-Path $releaseDir "A股智能分析系统"

# 创建发布目录
if (-not (Test-Path $releaseDir)) {
    New-Item -ItemType Directory -Path $releaseDir -Force | Out-Null
}
if (Test-Path $releaseAppDir) {
    Write-Host "  → 清理旧的发布包..." -ForegroundColor Gray
    Remove-Item -Recurse -Force $releaseAppDir
}
New-Item -ItemType Directory -Path $releaseAppDir -Force | Out-Null

# 复制 Go UI 可执行文件
Write-Host "  → 复制 Go UI 可执行文件..." -ForegroundColor Gray
$goExe = Join-Path $uiDir "stock-analysis-ui-windows-amd64.exe"
Copy-Item $goExe (Join-Path $releaseAppDir "A股智能分析系统.exe")

# 复制 Python 引擎
Write-Host "  → 复制 Python 引擎..." -ForegroundColor Gray
$releasePythonBundle = Join-Path $releaseAppDir "python_bundle"
Copy-Item -Recurse $uiPythonBundleDir $releasePythonBundle

# 复制环境变量模板
$envExample = Join-Path $stockAnalysisDir "src\a_stock_analysis\env.example"
if (Test-Path $envExample) {
    Write-Host "  → 复制环境变量模板..." -ForegroundColor Gray
    Copy-Item $envExample (Join-Path $releaseAppDir ".env.example")
}

# 创建启动脚本
Write-Host "  → 创建启动脚本..." -ForegroundColor Gray
$startupScript = @"
@echo off
echo 正在启动A股智能分析系统...
start "" "A股智能分析系统.exe"
"@
$startupScript | Out-File -FilePath (Join-Path $releaseAppDir "启动.bat") -Encoding ASCII

# 创建 README
Write-Host "  → 创建 README..." -ForegroundColor Gray
$readme = @"
# A股智能分析系统

## 使用说明

1. 双击"启动.bat"或"A股智能分析系统.exe"启动程序
2. 浏览器会自动打开，如果没有，请手动访问 http://localhost:8080
3. 在界面中输入股票信息开始分析

## 配置说明

- 如需配置 API 密钥，请复制 .env.example 为 .env 并编辑
- 环境变量文件应放在程序同目录下
- 支持 OpenAI、Anthropic (Claude) 等 AI 模型

## 注意事项

- 首次运行可能需要几秒钟加载 Python 引擎
- 确保系统已安装 Microsoft Visual C++ Redistributable
  下载地址: https://aka.ms/vs/17/release/vc_redist.x64.exe
- Windows Defender 可能会提示，选择"仍要运行"即可
- 需要网络连接来访问 AI API 和获取股票数据

## 系统要求

- Windows 10/11 (64位)
- 8GB+ RAM 推荐
- 网络连接

## 技术支持

如遇问题，请访问项目主页或提交 Issue

构建时间: $(Get-Date -Format "yyyy-MM-dd HH:mm:ss")
"@
$readme | Out-File -FilePath (Join-Path $releaseAppDir "README.txt") -Encoding UTF8

Write-Host "✓ 发布包已创建" -ForegroundColor Green

# ============================================================================
# 完成
# ============================================================================
Write-Host ""
Write-Host "========================================" -ForegroundColor Green
Write-Host " 构建完成！" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Green
Write-Host ""
Write-Host "发布包位置: " -NoNewline
Write-Host $releaseAppDir -ForegroundColor Cyan
Write-Host ""
Write-Host "文件列表:" -ForegroundColor Yellow
Get-ChildItem $releaseAppDir | Format-Table Name, Length -AutoSize
Write-Host ""
Write-Host "您可以将整个 'A股智能分析系统' 文件夹分发给用户使用" -ForegroundColor Green
Write-Host ""

# 询问是否打开发布目录
$response = Read-Host "是否打开发布目录? (Y/N)"
if ($response -eq 'Y' -or $response -eq 'y') {
    Start-Process explorer.exe -ArgumentList $releaseAppDir
}
