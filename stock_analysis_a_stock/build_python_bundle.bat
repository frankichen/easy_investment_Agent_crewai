@echo off
REM Build script to bundle Python application using PyInstaller

echo ====================================
echo Building Python Bundle with PyInstaller
echo ====================================

REM Check if poetry is available
where poetry >nul 2>nul
if %ERRORLEVEL% EQU 0 (
    echo Using Poetry environment...
    poetry run pip install pyinstaller
    poetry run pyinstaller build_pyinstaller.spec --clean
) else (
    REM Check if virtual environment exists
    if exist venv\Scripts\activate.bat (
        echo Using virtual environment...
        call venv\Scripts\activate.bat
        pip install pyinstaller
        pyinstaller build_pyinstaller.spec --clean
    ) else if defined VIRTUAL_ENV (
        echo Using active virtual environment...
        pip install pyinstaller
        pyinstaller build_pyinstaller.spec --clean
    ) else (
        echo ERROR: No poetry or virtual environment detected!
        echo Please install dependencies first:
        echo   运行: setup_env.bat
        echo   然后运行: build_python_bundle_venv.bat
        echo   OR
        echo   poetry install --no-root
        exit /b 1
    )
)

if %ERRORLEVEL% EQU 0 (
    echo.
    echo ====================================
    echo Python bundle created successfully!
    echo Location: dist\stock_analysis_engine\
    echo ====================================
) else (
    echo.
    echo ====================================
    echo ERROR: Failed to create Python bundle
    echo ====================================
    exit /b 1
)
