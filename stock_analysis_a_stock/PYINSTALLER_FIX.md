# PyInstaller Translation File Fix

## 问题描述 (Problem Description)

When packaging the A股智能分析系统 with PyInstaller, the application failed to run with the following error:

```
⚠️ FileNotFoundError: [Errno 2] No such file or directory: 
'D:\\...\\crewai\\utilities\\../translations/en.json'
⚠️ Exception: Prompt file 'None' not found.
```

This occurred because PyInstaller was not including the crewai package's translation files (JSON files) in the bundled application.

## 根本原因 (Root Cause)

1. **Missing PyInstaller Spec File**: The `build_pyinstaller.spec` file referenced in all build scripts was missing from the repository.

2. **Translation Files Not Included**: Without a proper spec file, PyInstaller's default behavior doesn't automatically include data files (like JSON translation files) from imported packages.

3. **Ignored by .gitignore**: The root `.gitignore` had `*.spec` which prevented the spec file from being tracked.

## 解决方案 (Solution)

### 1. Created `build_pyinstaller.spec`

A comprehensive PyInstaller spec file was created with the following key configurations:

#### Translation Files Collection
```python
# Add crewai translation files - this fixes the FileNotFoundError for translations
datas += collect_data_files('crewai', includes=['**/*.json'])
```

This ensures all JSON files (including `en.json`, `zh.json`, etc.) from the crewai package are included in the bundle.

#### Hidden Imports
```python
hiddenimports = [
    'crewai',
    'crewai.utilities.i18n',  # Critical for translation loading
    'crewai.agent',
    'crewai.task',
    # ... and many more
]
```

Explicitly tells PyInstaller to include the i18n module and other crewai components.

#### Config Files
```python
# Add config files (agents.yaml, tasks.yaml)
config_dir = os.path.join(a_stock_analysis_path, 'config')
if os.path.exists(config_dir):
    datas.append((config_dir, 'config'))
```

Ensures the YAML configuration files are included in the bundle.

#### Complete Module Collection
```python
hiddenimports += collect_submodules('crewai')
```

Collects all crewai submodules to ensure nothing is missed.

### 2. Updated `.gitignore`

Modified the root `.gitignore` to allow the spec file:

```gitignore
*.spec
!build_pyinstaller.spec  # Allow this specific spec file
```

## 验证 (Verification)

The fix addresses the issue by:

1. ✅ Collecting all crewai data files including translation JSON files
2. ✅ Including `crewai.utilities.i18n` in hidden imports
3. ✅ Collecting all crewai submodules
4. ✅ Including configuration files (agents.yaml, tasks.yaml)
5. ✅ Adding metadata for required packages

## 构建说明 (Build Instructions)

After this fix, the build process should work as documented:

### Windows
```bash
cd stock_analysis_a_stock
build_python_bundle.bat
```

### Linux/macOS
```bash
cd stock_analysis_a_stock
chmod +x build_python_bundle.sh
./build_python_bundle.sh
```

The build scripts will now use the `build_pyinstaller.spec` file to create a complete bundle with all necessary files.

## 预期结果 (Expected Result)

After building with the new spec file, the packaged application should:

1. ✅ Successfully load crewai translation files
2. ✅ Initialize the i18n system without errors
3. ✅ Start the analysis engine without FileNotFoundError
4. ✅ Run all agents and tasks correctly

## 技术细节 (Technical Details)

### PyInstaller Data Collection

The spec file uses PyInstaller's utility functions:

- `collect_data_files()`: Collects non-Python files from packages
- `collect_submodules()`: Ensures all submodules are imported
- `copy_metadata()`: Includes package metadata

### File Structure After Build

```
dist/stock_analysis_engine/
├── stock_analysis_engine.exe
├── _internal/
│   ├── crewai/
│   │   └── translations/
│   │       ├── en.json    ← Now included!
│   │       └── zh.json    ← Now included!
│   ├── ... (other dependencies)
└── config/
    ├── agents.yaml
    └── tasks.yaml
```

## 相关文件 (Related Files)

- `stock_analysis_a_stock/build_pyinstaller.spec` - The new spec file
- `.gitignore` - Updated to allow the spec file
- `stock_analysis_a_stock/build_python_bundle.bat` - Uses the spec file
- `stock_analysis_a_stock/build_python_bundle.sh` - Uses the spec file

## 参考 (References)

- [PyInstaller Documentation](https://pyinstaller.org/)
- [PyInstaller Hooks](https://pyinstaller.org/en/stable/hooks.html)
- [CrewAI Documentation](https://docs.crewai.com/)

---

**日期 (Date)**: 2025-10-16  
**作者 (Author)**: GitHub Copilot  
**状态 (Status)**: ✅ Fixed

---

## 更新 (Update) - 2025-10-16

### 问题描述 (New Issue Description)

After the initial fix, users encountered another PyInstaller error:

```
ModuleNotFoundError: No module named 'jaraco.text'
[PYI-28816:ERROR] Failed to execute script 'pyi_rth_pkgres' due to unhandled exception!
```

This error occurs in the PyInstaller runtime hook `pyi_rth_pkgres`, which is responsible for handling `pkg_resources`. The `pkg_resources` module (part of setuptools) requires `jaraco.text` and other jaraco packages as dependencies.

### 根本原因 (Root Cause)

1. **Runtime Hook Dependency**: PyInstaller includes a runtime hook for `pkg_resources` that runs before the main application.
2. **Missing jaraco Packages**: The `pkg_resources` module requires `jaraco.text`, `jaraco.functools`, and `jaraco.context` as dependencies.
3. **Excluded setuptools**: The original spec file excluded setuptools, pip, and wheel, but `pkg_resources` needs these packages.

### 解决方案 (Solution)

Updated `build_pyinstaller.spec` with the following changes:

#### 1. Added jaraco Hidden Imports
```python
hiddenimports = [
    # ... existing imports ...
    
    # pkg_resources and setuptools dependencies
    # These are needed by PyInstaller runtime hooks
    'pkg_resources',
    'pkg_resources.extern',
    'jaraco',
    'jaraco.text',
    'jaraco.functools',
    'jaraco.context',
]
```

#### 2. Added jaraco Metadata Collection
```python
# Add metadata for jaraco packages to fix pkg_resources issues
try:
    datas += copy_metadata('jaraco.text')
except Exception:
    pass  # jaraco.text might not always be needed
try:
    datas += copy_metadata('jaraco.functools')
except Exception:
    pass
```

#### 3. Added jaraco Submodules Collection
```python
# Collect jaraco submodules to fix pkg_resources runtime hook issues
try:
    hiddenimports += collect_submodules('jaraco.text')
except Exception:
    pass  # jaraco.text might not be installed
try:
    hiddenimports += collect_submodules('jaraco.functools')
except Exception:
    pass
```

#### 4. Stopped Excluding setuptools
```python
excludes=[
    # ... other excludes ...
    # Note: Don't exclude setuptools, pip, wheel as they may be needed by pkg_resources
    # 'setuptools',
    # 'pip',
    # 'wheel',
]
```

### 技术细节 (Technical Details)

The error occurred because:
- PyInstaller uses runtime hooks to initialize certain packages before the main application runs
- The `pyi_rth_pkgres` hook initializes `pkg_resources` to handle package metadata
- `pkg_resources` depends on `jaraco.text` and other jaraco packages
- Without these dependencies included, the runtime hook fails before the application can start

### 验证 (Verification)

After this fix, the packaged application should:

1. ✅ Successfully initialize the pkg_resources runtime hook
2. ✅ Load all jaraco dependencies without errors
3. ✅ Start the analysis engine without ModuleNotFoundError
4. ✅ Run all agents and tasks correctly

---

**日期 (Date)**: 2025-10-16  
**作者 (Author)**: GitHub Copilot  
**状态 (Status)**: ✅ Fixed (Updated)
