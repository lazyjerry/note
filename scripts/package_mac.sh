#!/usr/bin/env bash
set -euo pipefail

# Mac 筆記本應用程式一鍵打包腳本（Fyne）
# 此腳本負責完整的 macOS 應用程式打包流程
#
# 使用方法：
# ./scripts/package_mac.sh [--clean] [版本號]
#
# 執行流程：
# 1. 檢查必要工具（Go, Xcode Command Line Tools, Fyne CLI）
# 2. 生成應用程式資源（圖示、字體等）
# 3. 整理依賴並編譯應用程式
# 4. 使用 Fyne 打包工具創建 .app 包
# 5. 移動到發布目錄並可選擇自動開啟

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

# 應用程式設定（可由環境變數覆蓋）
APP_NAME_DEFAULT="Mac筆記本"
APP_ID_DEFAULT="com.notebook.mac-notebook-app"
APP_NAME="${APP_NAME:-$APP_NAME_DEFAULT}"
APP_ID="${APP_ID:-$APP_ID_DEFAULT}"
APP_VERSION="${APP_VERSION:-1.0.0}"
APP_BUILD="${APP_BUILD:-1}"
OPEN_AFTER_BUILD="${OPEN_AFTER_BUILD:-1}"

DIST_DIR="$ROOT_DIR/dist"
mkdir -p "$DIST_DIR"

info() { echo -e "\033[1;34m[INFO]\033[0m $*"; }
warn() { echo -e "\033[1;33m[WARN]\033[0m $*"; }
err()  { echo -e "\033[1;31m[ERR ]\033[0m $*"; }

# 檢查 Go
if ! command -v go >/dev/null 2>&1; then
  err "未找到 go，請先安裝（建議：brew install go）"
  exit 1
fi

# 檢查 Xcode Command Line Tools（Fyne 在 macOS 需要）
if ! xcode-select -p >/dev/null 2>&1; then
  warn "未偵測到 Xcode Command Line Tools，嘗試安裝..."
  xcode-select --install || true
  err "請完成 Command Line Tools 安裝後再重試。"
  exit 1
fi

# 檢查或安裝 fyne CLI
if ! command -v fyne >/dev/null 2>&1; then
  info "未找到 fyne CLI，嘗試安裝..."
  GO_BIN="$(go env GOPATH)/bin"
  mkdir -p "$GO_BIN"
  GOBIN="$GO_BIN" go install fyne.io/tools/cmd/fyne@latest
  export PATH="$GO_BIN:$PATH"
  if ! command -v fyne >/dev/null 2>&1; then
    err "安裝 fyne CLI 失敗，請手動執行：go install fyne.io/tools/cmd/fyne@latest"
    exit 1
  fi
fi

# 生成應用程式資源
info "生成應用程式資源（圖示、字體等）..."
if [[ -f "scripts/generate_resources.go" ]]; then
  go run scripts/generate_resources.go || warn "資源生成失敗，繼續打包流程"
else
  warn "未找到資源生成腳本，跳過資源生成"
fi

# 整理依賴
info "整理依賴..."
go mod tidy

# 可選：清理
if [[ "${1:-}" == "--clean" ]]; then
  info "清理舊的產物..."
  rm -rf "$DIST_DIR/${APP_NAME}.app" "$ROOT_DIR/${APP_NAME}.app" "$ROOT_DIR"/*.app || true
fi

# 打包應用程式
info "開始打包 macOS 應用程式..."
PKG_APP=""
if [[ -f "$ROOT_DIR/FyneApp.toml" ]]; then
  info "偵測到 FyneApp.toml，採用設定檔打包"
  fyne package -os darwin
  PKG_APP="$ROOT_DIR/${APP_NAME}.app"
else
  info "未找到 FyneApp.toml，使用預設資訊打包"
  fyne package -os darwin -name "$APP_NAME" -appID "$APP_ID" -appBuild "$APP_BUILD" -appVersion "$APP_VERSION"
  PKG_APP="$ROOT_DIR/${APP_NAME}.app"
fi

# 偵測產物（若檔名有差異，自動找第一個 .app）
if [[ ! -d "$PKG_APP" ]]; then
  CANDIDATE="$(ls -1d "$ROOT_DIR"/*.app 2>/dev/null | head -n1 || true)"
  if [[ -n "${CANDIDATE}" ]]; then
    PKG_APP="$CANDIDATE"
  else
    err "找不到已產生的 .app！"
    exit 1
  fi
fi

# 移動到 dist
APP_BASENAME="$(basename "$PKG_APP")"
mv -f "$PKG_APP" "$DIST_DIR/" 2>/dev/null || { cp -R "$PKG_APP" "$DIST_DIR/"; rm -rf "$PKG_APP"; }
OUT_APP="$DIST_DIR/$APP_BASENAME"

info "打包完成：$OUT_APP"

# 自動開啟（可用 OPEN_AFTER_BUILD=0 關閉）
if [[ "$OPEN_AFTER_BUILD" == "1" ]]; then
  info "開啟應用程式..."
  open "$OUT_APP"
fi

info "完成。"


