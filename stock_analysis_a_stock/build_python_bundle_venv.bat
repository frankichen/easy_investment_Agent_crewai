@echo off
REM Build script to bundle Python application using PyInstaller with venv

echo ====================================
echo Building Python Bundle with PyInstaller
echo ====================================

REM Check if virtual environment exists
if not exist venv\Scripts\activate.bat (
    echo ERROR: 虚拟环境不存在！
    echo 请先运行: setup_env.bat
    exit /b 1
)

echo 使用虚拟环境构建...
call venv\Scripts\activate.bat
if %ERRORLEVEL% NEQ 0 (
    echo ERROR: 激活虚拟环境失败
    exit /b 1
)

pyinstaller build_pyinstaller.spec --clean

if %ERRORLEVEL% EQU 0 (
    echo.
    echo ====================================
    echo Python bundle 创建成功！
    echo 位置: dist\stock_analysis_engine\
    echo ====================================
) else (
    echo.
    echo ====================================
    echo ERROR: 创建Python bundle失败
    echo ====================================
    exit /b 1
)
