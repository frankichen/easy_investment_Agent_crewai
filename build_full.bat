@echo off
REM ============================================================================
REM A股智能分析系统 - 一键完整构建脚本
REM ============================================================================
REM 自动构建 Python 引擎 + Go UI + 创建发布包
REM ============================================================================

echo.
echo ========================================
echo  A股智能分析系统 - 一键完整构建
echo ========================================
echo.

REM 检查是否已设置虚拟环境
if not exist "stock_analysis_a_stock\venv\Scripts\python.exe" (
    echo [!] 检测到虚拟环境未设置
    echo [!] 正在自动设置虚拟环境...
    echo.
    cd stock_analysis_a_stock
    call setup_env.bat
    if %ERRORLEVEL% NEQ 0 (
        echo ERROR: 虚拟环境设置失败
        cd ..
        pause
        exit /b 1
    )
    cd ..
    echo.
    echo 虚拟环境设置完成！
    echo.
)

REM 调用 PowerShell 脚本执行构建
echo 开始完整构建流程...
echo.
powershell -ExecutionPolicy Bypass -File "%~dp0build_full.ps1"

if %ERRORLEVEL% NEQ 0 (
    echo.
    echo ========================================
    echo  构建失败！
    echo ========================================
    pause
    exit /b 1
)

echo.
echo ========================================
echo  构建成功完成！
echo ========================================
echo.
pause
