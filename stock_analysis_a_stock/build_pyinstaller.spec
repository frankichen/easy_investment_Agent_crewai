# -*- mode: python ; coding: utf-8 -*-
"""
PyInstaller spec file for A股智能分析系统
This spec file configures how PyInstaller packages the Python application.
"""

import sys
import os
from PyInstaller.utils.hooks import collect_data_files, collect_submodules, copy_metadata

# Get the absolute path to the source directory
src_path = os.path.join(os.getcwd(), 'src')
a_stock_analysis_path = os.path.join(src_path, 'a_stock_analysis')

# Collect data files from packages
datas = []

# Add crewai translation files - this fixes the FileNotFoundError for translations
datas += collect_data_files('crewai', includes=['**/*.json'])

# Add config files (agents.yaml, tasks.yaml)
config_dir = os.path.join(a_stock_analysis_path, 'config')
if os.path.exists(config_dir):
    datas.append((config_dir, 'config'))

# Add env.example as reference
env_example = os.path.join(a_stock_analysis_path, 'env.example')
if os.path.exists(env_example):
    datas.append((env_example, '.'))

# Collect metadata for packages that need it
datas += copy_metadata('crewai')
datas += copy_metadata('akshare')
datas += copy_metadata('pandas')

# Hidden imports - modules that PyInstaller might miss
hiddenimports = [
    # Core crewai modules
    'crewai',
    'crewai.agent',
    'crewai.task',
    'crewai.crew',
    'crewai.process',
    'crewai.tools',
    'crewai.utilities',
    'crewai.utilities.i18n',
    'crewai.project',
    'crewai.events',
    'crewai.llm',
    
    # Crewai tools
    'crewai_tools',
    
    # AKShare and its dependencies
    'akshare',
    
    # Data processing
    'pandas',
    'pandas._libs',
    'pandas._libs.tslibs',
    'pandas._libs.tslibs.timedeltas',
    'pandas._libs.tslibs.np_datetime',
    'pandas._libs.tslibs.nattype',
    'pandas._libs.tslibs.timestamps',
    'numpy',
    'numpy.core',
    'numpy.core._methods',
    'numpy.lib.format',
    
    # HTTP and web
    'requests',
    'urllib3',
    'html2text',
    
    # Environment and configuration
    'dotenv',
    
    # YAML support
    'yaml',
    
    # JSON support
    'json',
    
    # Other potential dependencies
    'pydantic',
    'pydantic_core',
]

# Additional hidden imports from submodules
hiddenimports += collect_submodules('crewai')
hiddenimports += collect_submodules('akshare')

# Binaries - none explicitly needed
binaries = []

# Analysis configuration
a = Analysis(
    ['src/a_stock_analysis/cli_entry.py'],
    pathex=[src_path],
    binaries=binaries,
    datas=datas,
    hiddenimports=hiddenimports,
    hookspath=[],
    hooksconfig={},
    runtime_hooks=[],
    excludes=[
        # Exclude unnecessary packages to reduce size
        'matplotlib',
        'tkinter',
        'PyQt5',
        'PyQt6',
        'PySide2',
        'PySide6',
        'IPython',
        'jupyter',
        'notebook',
        'pytest',
        'setuptools',
        'pip',
        'wheel',
    ],
    noarchive=False,
    optimize=0,
)

# PYZ (Python ZIP archive)
pyz = PYZ(a.pure)

# EXE configuration
exe = EXE(
    pyz,
    a.scripts,
    [],
    exclude_binaries=True,
    name='stock_analysis_engine',
    debug=False,
    bootloader_ignore_signals=False,
    strip=False,
    upx=True,  # Use UPX compression if available
    console=True,  # Keep console window to see output
    disable_windowed_traceback=False,
    argv_emulation=False,
    target_arch=None,
    codesign_identity=None,
    entitlements_file=None,
)

# COLLECT - gather all files into dist folder
coll = COLLECT(
    exe,
    a.binaries,
    a.datas,
    strip=False,
    upx=True,
    upx_exclude=[],
    name='stock_analysis_engine',
)
