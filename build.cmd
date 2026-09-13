@echo off
setlocal
cd /d "%~dp0"

set "APP=gah.exe"
set "VERSION=unknown"
for /f "usebackq delims=" %%V in ("GAHversion") do set "VERSION=%%V"

echo ============================================
echo  GoAccountHub one-click build
echo  version: %VERSION%
echo ============================================
echo.

where go >nul 2>nul
if errorlevel 1 (
    echo [ERROR] "go" was not found in PATH.
    exit /b 1
)
where npm >nul 2>nul
if errorlevel 1 (
    echo [ERROR] "npm" was not found in PATH.
    exit /b 1
)

pushd "GAHFrontend"
if not exist "node_modules" (
    echo [1/3] Installing frontend dependencies ...
    call npm ci
    if errorlevel 1 set "FAILED=npm ci"
) else (
    echo [1/3] Frontend dependencies already installed, skipping npm ci
)
if not defined FAILED (
    echo [2/3] Building frontend ...
    call npm run build
    if errorlevel 1 set "FAILED=npm run build"
)
popd
if defined FAILED (
    echo.
    echo [ERROR] %FAILED% failed.
    exit /b 1
)

if not exist "GAHFrontend\dist\index.html" (
    echo.
    echo [ERROR] GAHFrontend\dist\index.html is missing, the frontend was not built.
    exit /b 1
)

echo [3/3] Building backend, embedding GAHFrontend\dist ...
go build -trimpath -ldflags "-s -w" -o "%APP%" .
if errorlevel 1 (
    echo.
    echo [ERROR] go build failed.
    exit /b 1
)

for %%F in ("%APP%") do set "SIZE=%%~zF"
echo.
echo [OK] %APP% built, version %VERSION%, %SIZE% bytes
echo      Run: %APP% start
exit /b 0
