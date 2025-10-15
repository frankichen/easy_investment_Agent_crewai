# 🚀 一键构建指南

## 快速开始

### 方法 1: 一键构建（推荐）⭐

直接双击运行：
```
build_full.bat
```

这会自动完成：
- ✅ 检查并设置虚拟环境（如果需要）
- ✅ 构建 Python 分析引擎
- ✅ 构建 Go UI 程序
- ✅ 创建完整的发布包

### 方法 2: PowerShell 构建

在 PowerShell 中运行：
```powershell
.\build_full.ps1
```

## 📋 前提条件

### 必需工具
- ✅ **Python 3.12+** - [下载地址](https://www.python.org/downloads/)
- ✅ **Go 1.24+** - [下载地址](https://go.dev/dl/)

### 首次使用
如果是第一次构建，脚本会自动：
1. 创建 Python 虚拟环境
2. 安装所有必需的依赖：
   - crewai
   - crewai-tools
   - akshare
   - pyinstaller

## 🔧 构建流程

### 完整流程说明

```
[1/5] 检查构建环境
      ↓ 验证 Python 和 Go 是否安装
      
[2/5] 构建 Python 分析引擎
      ↓ 使用 PyInstaller 打包 Python 代码
      ↓ 生成独立可执行文件
      
[3/5] 准备打包资源
      ↓ 复制 Python 引擎到 UI 目录
      
[4/5] 构建 Go UI 程序
      ↓ 编译 Go 代码为 Windows 可执行文件
      
[5/5] 创建发布包
      ↓ 打包所有文件到 release 目录
      ✓ 完成！
```

### 构建时间

- 首次构建（含依赖安装）: **5-10 分钟**
- 后续构建: **2-5 分钟**

## 📦 构建输出

### 目录结构

构建完成后，会在 `release` 目录生成：

```
release/
└── A股智能分析系统/
    ├── A股智能分析系统.exe      # Go UI 程序（主程序）
    ├── 启动.bat                  # 快速启动脚本
    ├── README.txt                # 使用说明
    ├── .env.example              # 环境变量配置模板
    └── python_bundle/            # Python 分析引擎
        └── stock_analysis_engine/
            ├── stock_analysis_engine.exe  # Python 引擎
            └── _internal/                 # Python 运行时和依赖
                ├── python313.dll
                ├── python3.dll
                └── ... (其他依赖)
```

### 发布包大小

- 总大小: **~500MB - 800MB**
- 包含完整的 Python 运行时和所有依赖

## 🐛 常见问题

### 1. "Python 未找到"

**解决方案**:
- 安装 Python 3.12+
- 确保 Python 已添加到系统 PATH
- 重启终端/命令提示符

### 2. "Go 未找到"

**解决方案**:
- 安装 Go 1.24+
- 确保 Go 已添加到系统 PATH
- 重启终端/命令提示符

### 3. "虚拟环境不存在"

**解决方案**:
```cmd
cd stock_analysis_a_stock
setup_env.bat
cd ..
build_full.bat
```

### 4. PyInstaller 构建很慢

这是正常的，PyInstaller 需要分析和打包大量依赖。
- 首次构建: 5-10 分钟
- 后续增量构建会更快

### 5. "UPX 警告"或压缩错误

已通过在 spec 文件中设置 `upx=False` 解决。

### 6. "DLL 加载失败"

已修复：
- ✅ 禁用 UPX 压缩
- ✅ 设置正确的工作目录

### 7. "ModuleNotFoundError"

已修复：
- ✅ 添加所有必需模块到 `hiddenimports`
- ✅ 使用绝对导入

## 🔍 手动步骤（高级用户）

如果需要分步构建：

### 步骤 1: 设置环境（仅首次）
```cmd
cd stock_analysis_a_stock
setup_env.bat
```

### 步骤 2: 构建 Python 引擎
```powershell
cd stock_analysis_a_stock
powershell -ExecutionPolicy Bypass -File build_python_bundle.ps1
```

### 步骤 3: 构建完整程序
```cmd
cd ..
build_windows_exe.bat
```

## ✅ 验证构建

### 测试 Python 引擎
```cmd
cd stock_analysis_a_stock
test_engine.bat
```

### 测试完整程序
```cmd
cd release\A股智能分析系统
.\启动.bat
```

然后访问 http://localhost:8080

## 📝 修改记录

### v2.0 - 一键构建版本
- ✅ 创建 `build_full.ps1` - 完整 PowerShell 构建脚本
- ✅ 创建 `build_full.bat` - 一键批处理脚本
- ✅ 自动检测并设置虚拟环境
- ✅ 集成所有构建步骤
- ✅ 自动创建发布包

### v1.0 - 修复版本
- ✅ 修复 DLL 加载问题
- ✅ 修复模块导入错误
- ✅ 改进虚拟环境检测
- ✅ 禁用 UPX 压缩

## 🎯 下一步

构建完成后：

1. **测试程序**
   ```cmd
   cd release\A股智能分析系统
   .\启动.bat
   ```

2. **配置 API**
   - 复制 `.env.example` 为 `.env`
   - 填入您的 AI API 密钥

3. **分发程序**
   - 将 `release\A股智能分析系统` 整个文件夹打包
   - 分发给用户使用

4. **创建安装包（可选）**
   - 使用 Inno Setup 或 NSIS 创建安装程序
   - 或者直接压缩为 ZIP 文件分发

## 💡 提示

- 🔧 修改代码后，只需重新运行 `build_full.bat`
- 📦 发布包是独立的，可以在任何 Windows 10/11 系统运行
- 🚀 程序启动后会自动打开浏览器
- 🔑 API 密钥配置在 `.env` 文件中

## 🆘 获取帮助

如果遇到问题：
1. 查看 `BUILD_FIX_SUMMARY.md` - 问题修复总结
2. 查看 `TROUBLESHOOTING.md` - 故障排查指南
3. 提交 Issue 到 GitHub 仓库

祝构建顺利！🎉
