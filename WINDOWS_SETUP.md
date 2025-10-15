# Windows 编译说明

## 问题解决：虚拟环境设置

如果你遇到 "No poetry or virtual environment detected" 错误，按照以下步骤操作：

### 方法 1：使用新的脚本（推荐）

1. **设置虚拟环境**
   ```cmd
   cd stock_analysis_a_stock
   setup_env.bat
   ```
   这会创建虚拟环境并安装所有依赖。

2. **构建 Python Bundle**
   ```cmd
   build_python_bundle_venv.bat
   ```

3. **构建完整程序**
   ```cmd
   cd ..
   build_windows_exe.bat
   ```

### 方法 2：手动步骤

1. **进入项目目录**
   ```cmd
   cd stock_analysis_a_stock
   ```

2. **创建虚拟环境**
   ```cmd
   python -m venv venv
   ```

3. **激活虚拟环境**
   ```cmd
   venv\Scripts\activate.bat
   ```

4. **安装依赖**
   ```cmd
   pip install --upgrade pip
   pip install crewai crewai-tools akshare pyinstaller
   ```

5. **构建**
   ```cmd
   pyinstaller build_pyinstaller.spec --clean
   ```

## 常见问题

### "系统找不到指定的路径"
- 确保你在正确的目录（`stock_analysis_a_stock`）中
- 使用 `cd` 命令切换到正确的目录

### "VIRTUAL_ENV 未定义"
- 需要先激活虚拟环境：`venv\Scripts\activate.bat`
- 或者使用新的 `setup_env.bat` 脚本

## 一键构建（从根目录）

如果已经设置好虚拟环境，可以直接运行：
```cmd
build_windows_exe.bat
```

这会自动构建 Python Bundle 和 Go UI。
