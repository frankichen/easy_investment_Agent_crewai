@echo off
REM ============================================================================
REM 验证构建文件格式的脚本 (Windows版本)
REM ============================================================================
REM 使用此脚本检查你构建的可执行文件是否为正确的格式
REM ============================================================================

echo.
echo ========================================
echo  构建文件格式验证工具
echo ========================================
echo.

REM 检查发布目录是否存在
if not exist "release\A股智能分析系统\" (
    echo ❌ 错误: release\A股智能分析系统 目录不存在
    echo 请先运行构建脚本
    exit /b 1
)

REM 检查主程序
echo [1] 检查 Go UI 程序...
set "ui_exe=release\A股智能分析系统\A股智能分析系统.exe"
if exist "%ui_exe%" (
    echo 文件: %ui_exe%
    echo ✅ 主程序存在
    
    REM 尝试获取文件版本信息
    powershell -Command "Get-Item '%ui_exe%' | Select-Object -ExpandProperty VersionInfo | Format-List FileDescription, ProductVersion"
) else (
    echo ❌ 错误: 主程序不存在
)

echo.
echo [2] 检查 Python 引擎...

REM 检查 .exe 版本
set "py_exe=release\A股智能分析系统\python_bundle\stock_analysis_engine\stock_analysis_engine.exe"
if exist "%py_exe%" (
    echo 文件: %py_exe%
    
    REM 检查文件大小
    for %%A in ("%py_exe%") do set size=%%~zA
    echo 文件大小: %size% bytes
    
    REM 检查 PE 头部 (MZ 标记)
    powershell -Command "$bytes = Get-Content '%py_exe%' -Encoding Byte -TotalCount 2; if ($bytes[0] -eq 77 -and $bytes[1] -eq 90) { Write-Host '✅ 正确: 有效的 Windows PE 可执行文件' -ForegroundColor Green } else { Write-Host '❌ 错误: 不是有效的 Windows 可执行文件' -ForegroundColor Red }"
) else (
    echo ⚠️  警告: stock_analysis_engine.exe 不存在
    
    REM 检查无扩展名版本
    set "py_noext=release\A股智能分析系统\python_bundle\stock_analysis_engine\stock_analysis_engine"
    if exist "%py_noext%" (
        echo.
        echo ❌ 发现文件: %py_noext% (无 .exe 扩展名^)
        echo    这可能是在 Linux 上编译的，无法在 Windows 上运行
        echo    请在 Windows 系统上重新构建
    )
)

REM 检查关键 DLL 文件
echo.
echo [3] 检查依赖文件...
set "internal_dir=release\A股智能分析系统\python_bundle\stock_analysis_engine\_internal"
if exist "%internal_dir%" (
    echo ✅ _internal 目录存在
    
    REM 统计文件数量
    for /f %%A in ('dir /b /a-d "%internal_dir%" 2^>nul ^| find /c /v ""') do set file_count=%%A
    echo    包含 %file_count% 个文件
) else (
    echo ❌ _internal 目录不存在
)

echo.
echo ========================================
echo  验证完成
echo ========================================
echo.
echo 测试建议:
echo 1. 尝试运行: release\A股智能分析系统\启动.bat
echo 2. 或直接运行: release\A股智能分析系统\A股智能分析系统.exe
echo 3. 如果失败，查看 TROUBLESHOOTING.md
echo.
pause
