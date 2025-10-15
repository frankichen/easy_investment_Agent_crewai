@echo off
REM ============================================================================
REM 设置虚拟环境和安装依赖
REM ============================================================================

echo ====================================
echo 设置虚拟环境
echo ====================================

REM 检查Python是否可用
where python >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo ERROR: Python未找到，请先安装Python
    exit /b 1
)

echo [1/3] 创建虚拟环境...
if exist venv (
    echo 虚拟环境已存在，跳过创建
) else (
    python -m venv venv
    if %ERRORLEVEL% NEQ 0 (
        echo ERROR: 创建虚拟环境失败
        exit /b 1
    )
    echo ✓ 虚拟环境创建成功
)

echo.
echo [2/3] 激活虚拟环境并安装依赖...
echo 正在安装，请稍候...
call venv\Scripts\activate.bat
if %ERRORLEVEL% NEQ 0 (
    echo ERROR: 激活虚拟环境失败
    exit /b 1
)

pip install --upgrade pip
if %ERRORLEVEL% NEQ 0 (
    echo ERROR: 升级pip失败
    exit /b 1
)

pip install crewai crewai-tools akshare pyinstaller
if %ERRORLEVEL% NEQ 0 (
    echo ERROR: 安装依赖失败
    exit /b 1
)

echo.
echo ====================================
echo ✓ 环境设置完成！
echo ====================================
echo.
echo 下一步：运行 build_python_bundle_venv.bat 来构建
echo.
