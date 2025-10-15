#!/bin/bash
# ============================================================================
# A股智能分析系统 - Linux单体可执行文件构建脚本
# ============================================================================
# 本脚本将Python核心和Go UI打包成Linux可执行文件包
# 必须在Linux系统上运行以生成正确的Linux可执行文件
# ============================================================================

set -e

echo ""
echo "========================================"
echo " A股智能分析系统 - Linux版本构建"
echo "========================================"
echo ""

# 检查必要工具
echo "[1/5] 检查构建环境..."
if ! command -v python3 &> /dev/null && ! command -v python &> /dev/null; then
    echo "ERROR: Python未找到，请先安装Python 3.12+"
    exit 1
fi

if ! command -v go &> /dev/null; then
    echo "ERROR: Go未找到，请先安装Go 1.24+"
    exit 1
fi

# 检查操作系统
if [ "$(uname)" != "Linux" ]; then
    echo "⚠️  警告: 您似乎不在Linux系统上"
    echo "   当前系统: $(uname)"
    echo "   此脚本应该在Linux上运行以生成Linux可执行文件"
    echo ""
    read -p "按 Enter 继续，或 Ctrl+C 取消..." dummy
fi

echo "✓ Python和Go环境已就绪"

# 步骤1: 构建Python Bundle
echo ""
echo "[2/5] 构建Python分析引擎..."
cd stock_analysis_a_stock
chmod +x build_python_bundle.sh
./build_python_bundle.sh
cd ..
echo "✓ Python引擎构建完成"

# 步骤2: 复制Python Bundle到UI目录
echo ""
echo "[3/5] 准备打包资源..."
mkdir -p ui/python_bundle
cp -r stock_analysis_a_stock/dist/stock_analysis_engine ui/python_bundle/
echo "✓ Python引擎已复制到UI目录"

# 步骤3: 构建Go UI (Linux版本)
echo ""
echo "[4/5] 构建Go UI程序 (Linux)..."
cd ui
go build -ldflags="-s -w" -o stock-analysis-ui-linux
cd ..
echo "✓ Go UI程序构建完成"

# 步骤4: 创建发布包
echo ""
echo "[5/5] 创建发布包..."
mkdir -p "release/A股智能分析系统-Linux"

# 复制可执行文件
cp "ui/stock-analysis-ui-linux" "release/A股智能分析系统-Linux/a_stock_analysis"
chmod +x "release/A股智能分析系统-Linux/a_stock_analysis"

# 复制Python引擎
cp -r "ui/python_bundle" "release/A股智能分析系统-Linux/"

# 复制环境变量模板
if [ -f "stock_analysis_a_stock/src/a_stock_analysis/env.example" ]; then
    cp "stock_analysis_a_stock/src/a_stock_analysis/env.example" "release/A股智能分析系统-Linux/.env.example"
fi

# 创建启动脚本
cat > "release/A股智能分析系统-Linux/启动.sh" << 'EOF'
#!/bin/bash
echo "正在启动A股智能分析系统..."
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$DIR"
./a_stock_analysis
EOF
chmod +x "release/A股智能分析系统-Linux/启动.sh"

# 创建README
cat > "release/A股智能分析系统-Linux/README.txt" << 'EOF'
# A股智能分析系统 (Linux版本)

使用说明:
1. 运行 ./启动.sh 或 ./a_stock_analysis 启动程序
2. 浏览器会自动打开，如果没有，请手动访问 http://localhost:8080
3. 在界面中输入股票信息开始分析

配置说明:
- 如需配置API密钥，请复制.env.example为.env并编辑
- 环境变量文件应放在程序同目录下

注意事项:
- 首次运行可能需要几秒钟加载Python引擎
- 确保系统已安装必要的系统库
- 某些Linux发行版可能需要额外的依赖
EOF

echo "✓ 发布包已创建"

# 完成
echo ""
echo "========================================"
echo " 构建完成！"
echo "========================================"
echo ""
echo "发布包位置: release/A股智能分析系统-Linux/"
echo ""
echo "文件列表:"
ls -lh "release/A股智能分析系统-Linux/"
echo ""
echo "测试运行:"
echo "  cd release/A股智能分析系统-Linux/"
echo "  ./启动.sh"
echo ""
