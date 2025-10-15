#!/bin/bash
# ============================================================================
# 验证构建文件格式的脚本
# ============================================================================
# 使用此脚本检查你构建的可执行文件是否为正确的格式
# ============================================================================

echo ""
echo "========================================"
echo " 构建文件格式验证工具"
echo "========================================"
echo ""

# 检查发布目录是否存在
if [ ! -d "release/A股智能分析系统" ]; then
    echo "❌ 错误: release/A股智能分析系统 目录不存在"
    echo "请先运行构建脚本"
    exit 1
fi

# 检查主程序
echo "[1] 检查 Go UI 程序..."
ui_exe="release/A股智能分析系统/A股智能分析系统.exe"
if [ -f "$ui_exe" ]; then
    file_type=$(file "$ui_exe")
    echo "文件: $ui_exe"
    echo "类型: $file_type"
    
    if echo "$file_type" | grep -q "PE32.*Windows"; then
        echo "✅ 正确: Windows PE 可执行文件"
    else
        echo "❌ 错误: 不是有效的 Windows 可执行文件"
        echo "   请在 Windows 系统上重新构建"
    fi
else
    echo "❌ 错误: 主程序不存在"
fi

echo ""
echo "[2] 检查 Python 引擎..."

# 检查 .exe 版本
py_exe="release/A股智能分析系统/python_bundle/stock_analysis_engine/stock_analysis_engine.exe"
if [ -f "$py_exe" ]; then
    file_type=$(file "$py_exe")
    echo "文件: $py_exe"
    echo "类型: $file_type"
    
    if echo "$file_type" | grep -q "PE32.*Windows"; then
        echo "✅ 正确: Windows PE 可执行文件"
    else
        echo "❌ 错误: 不是有效的 Windows 可执行文件"
        echo "   这个文件虽然有 .exe 扩展名，但是 Linux 格式"
        echo "   请在 Windows 系统上重新构建"
    fi
else
    echo "⚠️  警告: stock_analysis_engine.exe 不存在"
fi

# 检查无扩展名版本（Linux格式）
py_linux="release/A股智能分析系统/python_bundle/stock_analysis_engine/stock_analysis_engine"
if [ -f "$py_linux" ]; then
    file_type=$(file "$py_linux")
    echo ""
    echo "文件: $py_linux"
    echo "类型: $file_type"
    
    if echo "$file_type" | grep -q "ELF.*executable"; then
        echo "⚠️  这是 Linux ELF 可执行文件，无法在 Windows 上运行"
        echo "   请在 Windows 系统上重新构建"
    fi
fi

echo ""
echo "========================================"
echo " 验证完成"
echo "========================================"
echo ""
echo "如果发现错误，请按以下步骤修复："
echo "1. 在 Windows 10/11 系统上安装 Python 3.12+ 和 Go 1.24+"
echo "2. 运行: build_windows_exe.bat"
echo "3. 如果没有 Windows 环境，使用虚拟机或 CI/CD"
echo ""
echo "详细说明请参考: CROSS_PLATFORM_BUILD.md"
echo ""
