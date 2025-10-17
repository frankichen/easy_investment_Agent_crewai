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

# Add tiktoken data files - fixes "Unknown encoding cl100k_base" error
datas += collect_data_files('tiktoken')
datas += collect_data_files('tiktoken_ext')

# Add litellm data files - fixes anthropic_tokenizer.json and other tokenizer files
datas += collect_data_files('litellm', includes=['**/*.json'])

# Add akshare data files - fixes calendar.json and other data files
datas += collect_data_files('akshare', includes=['**/*.json', '**/*.csv', '**/*.txt'])

# Add config files (agents.yaml, tasks.yaml) - must be in a_stock_analysis package
config_dir = os.path.join(a_stock_analysis_path, 'config')
if os.path.exists(config_dir):
    datas.append((config_dir, 'a_stock_analysis/config'))

# Add env.example as reference
env_example = os.path.join(a_stock_analysis_path, 'env.example')
if os.path.exists(env_example):
    datas.append((env_example, '.'))

# Collect metadata for packages that need it
datas += copy_metadata('crewai')
datas += copy_metadata('akshare')
datas += copy_metadata('pandas')
datas += copy_metadata('tiktoken')
datas += copy_metadata('litellm')

# Add metadata for jaraco packages to fix pkg_resources issues
try:
    datas += copy_metadata('jaraco.text')
except Exception:
    pass  # jaraco.text might not always be needed
try:
    datas += copy_metadata('jaraco.functools')
except Exception:
    pass

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
    
    # tiktoken and encoding support
    'tiktoken',
    'tiktoken.core',
    'tiktoken.registry',
    'tiktoken_ext',
    'tiktoken_ext.openai_public',
    
    # litellm and its dependencies
    'litellm',
    'litellm.litellm_core_utils',
    'litellm.litellm_core_utils.tokenizers',
    
    # pkg_resources and setuptools dependencies
    # These are needed by PyInstaller runtime hooks
    'pkg_resources',
    'pkg_resources.extern',
    'jaraco',
    'jaraco.text',
    'jaraco.functools',
    'jaraco.context',
    
    # Application specific modules - tools package
    'a_stock_analysis.tools',
    'a_stock_analysis.tools.a_stock_data_tool',
    'a_stock_analysis.tools.financial_tool',
    'a_stock_analysis.tools.market_sentiment_tool',
    'a_stock_analysis.tools.calculator_tool',
    'a_stock_analysis.streaming_output',
    'a_stock_analysis.config_loader',
]

# Additional hidden imports from submodules
hiddenimports += collect_submodules('crewai')
hiddenimports += collect_submodules('akshare')

# Collect jaraco submodules to fix pkg_resources runtime hook issues
try:
    hiddenimports += collect_submodules('jaraco.text')
except Exception:
    pass  # jaraco.text might not be installed
try:
    hiddenimports += collect_submodules('jaraco.functools')
except Exception:
    pass

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
        # Note: Don't exclude IPython, jupyter as they may be needed by pyvis/crewai
        # Note: Don't exclude matplotlib as it may be needed by some analysis tools
        'tkinter',
        'PyQt5',
        'PyQt6',
        'PySide2',
        'PySide6',
        'pytest',
        # Note: Don't exclude setuptools, pip, wheel as they may be needed by pkg_resources
        # 'setuptools',
        # 'pip',
        # 'wheel',
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
