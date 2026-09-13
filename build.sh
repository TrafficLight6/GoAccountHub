#!/usr/bin/env sh
set -eu

cd "$(dirname "$0")"

APP="gah"
case "$(uname -s)" in
    MINGW*|MSYS*|CYGWIN*) APP="gah.exe" ;;
esac

VERSION="unknown"
if [ -f GAHversion ]; then
    VERSION="$(tr -d '\r\n' < GAHversion)"
fi

echo "============================================"
echo " GoAccountHub one-click build"
echo " version: ${VERSION}"
echo "============================================"
echo

if ! command -v go >/dev/null 2>&1; then
    echo "[ERROR] \"go\" was not found in PATH."
    exit 1
fi
if ! command -v npm >/dev/null 2>&1; then
    echo "[ERROR] \"npm\" was not found in PATH."
    exit 1
fi

cd GAHFrontend
if [ ! -d node_modules ]; then
    echo "[1/3] Installing frontend dependencies ..."
    npm ci
else
    echo "[1/3] Frontend dependencies already installed, skipping npm ci"
fi

echo "[2/3] Building frontend ..."
npm run build
cd ..

if [ ! -f GAHFrontend/dist/index.html ]; then
    echo
    echo "[ERROR] GAHFrontend/dist/index.html is missing, the frontend was not built."
    exit 1
fi

echo "[3/3] Building backend, embedding GAHFrontend/dist ..."
go build -trimpath -ldflags "-s -w" -o "${APP}" .

SIZE="$(wc -c < "${APP}" | tr -d ' ')"
echo
echo "[OK] ${APP} built, version ${VERSION}, ${SIZE} bytes"
echo "     Run: ./${APP} start"
