@echo off
REM 快速测试打包的 Python 引擎
echo ====================================
echo 测试 Python 分析引擎
echo ====================================
echo.

set "ENGINE_PATH=dist\stock_analysis_engine\stock_analysis_engine.exe"

if not exist "%ENGINE_PATH%" (
    echo ERROR: 找不到引擎文件: %ENGINE_PATH%
    echo 请先运行构建脚本
    exit /b 1
)

echo 测试命令: %ENGINE_PATH% --company "测试公司" --code "000001" --market "SZ"
echo.
echo ====================================
echo 开始测试...
echo ====================================
echo.

cd dist\stock_analysis_engine
stock_analysis_engine.exe --company "测试公司" --code "000001" --market "SZ"

echo.
echo ====================================
echo 测试完成，退出码: %ERRORLEVEL%
echo ====================================
