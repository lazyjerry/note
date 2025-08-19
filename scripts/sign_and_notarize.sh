#!/usr/bin/env bash
set -euo pipefail

# Mac 筆記本應用程式簽名和公證腳本
# 此腳本負責對 macOS 應用程式進行程式碼簽名和公證
#
# 使用方法：
# ./scripts/sign_and_notarize.sh [應用程式路徑] [開發者ID] [Apple ID] [應用程式密碼]
#
# 執行流程：
# 1. 檢查必要的簽名工具和憑證
# 2. 對應用程式進行程式碼簽名
# 3. 創建公證用的 ZIP 檔案
# 4. 上傳到 Apple 進行公證
# 5. 等待公證完成並裝訂票據
# 6. 驗證最終的應用程式

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

# 設定變數
APP_PATH="${1:-$ROOT_DIR/dist/Mac筆記本.app}"
DEVELOPER_ID="${2:-}"
APPLE_ID="${3:-}"
APP_PASSWORD="${4:-}"
BUNDLE_ID="com.notebook.mac-notebook-app"

# 顏色輸出函數
info() { echo -e "\033[1;34m[INFO]\033[0m $*"; }
warn() { echo -e "\033[1;33m[WARN]\033[0m $*"; }
err()  { echo -e "\033[1;31m[ERR ]\033[0m $*"; }
success() { echo -e "\033[1;32m[SUCCESS]\033[0m $*"; }

# checkRequirements 檢查簽名和公證所需的工具和設定
# 回傳：如果缺少必要工具或設定則退出
#
# 執行流程：
# 1. 檢查 Xcode Command Line Tools
# 2. 檢查開發者憑證
# 3. 檢查 Apple ID 和應用程式密碼
checkRequirements() {
    info "檢查簽名和公證要求..."
    
    # 檢查 codesign 工具
    if ! command -v codesign >/dev/null 2>&1; then
        err "未找到 codesign 工具，請安裝 Xcode Command Line Tools"
        exit 1
    fi
    
    # 檢查 xcrun 工具
    if ! command -v xcrun >/dev/null 2>&1; then
        err "未找到 xcrun 工具，請安裝 Xcode Command Line Tools"
        exit 1
    fi
    
    # 檢查應用程式是否存在
    if [[ ! -d "$APP_PATH" ]]; then
        err "應用程式不存在：$APP_PATH"
        err "請先執行打包腳本：./scripts/package_mac.sh"
        exit 1
    fi
    
    # 檢查開發者 ID
    if [[ -z "$DEVELOPER_ID" ]]; then
        warn "未提供開發者 ID，嘗試自動偵測..."
        DEVELOPER_ID=$(security find-identity -v -p codesigning | grep "Developer ID Application" | head -n1 | sed 's/.*"\(.*\)".*/\1/' || true)
        if [[ -z "$DEVELOPER_ID" ]]; then
            err "未找到開發者憑證，請確保已安裝有效的 Developer ID Application 憑證"
            err "或手動指定開發者 ID：./scripts/sign_and_notarize.sh [APP_PATH] [DEVELOPER_ID]"
            exit 1
        fi
        info "自動偵測到開發者 ID：$DEVELOPER_ID"
    fi
    
    # 檢查公證設定（可選）
    if [[ -n "$APPLE_ID" && -n "$APP_PASSWORD" ]]; then
        info "將進行完整的簽名和公證流程"
    else
        warn "未提供 Apple ID 或應用程式密碼，僅進行程式碼簽名"
        warn "如需公證，請提供：./scripts/sign_and_notarize.sh [APP_PATH] [DEVELOPER_ID] [APPLE_ID] [APP_PASSWORD]"
    fi
}

# signApplication 對應用程式進行程式碼簽名
# 參數：
#   - app_path: 應用程式路徑
#   - developer_id: 開發者 ID
# 回傳：簽名過程中的錯誤（如果有）
#
# 執行流程：
# 1. 清理擴展屬性
# 2. 簽名應用程式內的所有可執行檔案
# 3. 簽名主應用程式包
# 4. 驗證簽名結果
signApplication() {
    local app_path="$1"
    local developer_id="$2"
    
    info "開始對應用程式進行程式碼簽名..."
    info "應用程式路徑：$app_path"
    info "開發者 ID：$developer_id"
    
    # 清理擴展屬性
    info "清理擴展屬性..."
    xattr -cr "$app_path" || warn "清理擴展屬性時出現警告"
    
    # 簽名應用程式內的所有可執行檔案
    info "簽名應用程式內的可執行檔案..."
    find "$app_path" -type f -perm +111 -exec codesign --force --verify --verbose --sign "$developer_id" {} \; || {
        err "簽名可執行檔案失敗"
        return 1
    }
    
    # 簽名主應用程式包
    info "簽名主應用程式包..."
    codesign --force --verify --verbose --sign "$developer_id" "$app_path" || {
        err "簽名主應用程式包失敗"
        return 1
    }
    
    # 驗證簽名
    info "驗證程式碼簽名..."
    codesign --verify --deep --strict --verbose=2 "$app_path" || {
        err "簽名驗證失敗"
        return 1
    }
    
    success "程式碼簽名完成"
    return 0
}

# notarizeApplication 對應用程式進行公證
# 參數：
#   - app_path: 應用程式路徑
#   - apple_id: Apple ID
#   - app_password: 應用程式專用密碼
# 回傳：公證過程中的錯誤（如果有）
#
# 執行流程：
# 1. 創建公證用的 ZIP 檔案
# 2. 上傳到 Apple 進行公證
# 3. 等待公證完成
# 4. 裝訂公證票據到應用程式
notarizeApplication() {
    local app_path="$1"
    local apple_id="$2"
    local app_password="$3"
    
    info "開始公證流程..."
    
    # 創建公證用的 ZIP 檔案
    local zip_path="${app_path%.*}.zip"
    info "創建公證用的 ZIP 檔案：$zip_path"
    
    # 移除舊的 ZIP 檔案
    rm -f "$zip_path"
    
    # 創建 ZIP 檔案
    ditto -c -k --keepParent "$app_path" "$zip_path" || {
        err "創建 ZIP 檔案失敗"
        return 1
    }
    
    # 上傳進行公證
    info "上傳到 Apple 進行公證..."
    local submit_result
    submit_result=$(xcrun notarytool submit "$zip_path" --apple-id "$apple_id" --password "$app_password" --team-id "$DEVELOPER_ID" --wait 2>&1) || {
        err "公證上傳失敗"
        err "$submit_result"
        return 1
    }
    
    info "公證上傳結果："
    echo "$submit_result"
    
    # 檢查公證是否成功
    if echo "$submit_result" | grep -q "status: Accepted"; then
        success "公證成功"
        
        # 裝訂公證票據
        info "裝訂公證票據到應用程式..."
        xcrun stapler staple "$app_path" || {
            warn "裝訂公證票據失敗，但公證已完成"
        }
        
        # 驗證裝訂結果
        info "驗證公證票據..."
        xcrun stapler validate "$app_path" || {
            warn "公證票據驗證失敗"
        }
        
    else
        err "公證失敗"
        return 1
    fi
    
    # 清理 ZIP 檔案
    rm -f "$zip_path"
    
    return 0
}

# verifyFinalApplication 驗證最終的應用程式
# 參數：
#   - app_path: 應用程式路徑
# 回傳：驗證過程中的錯誤（如果有）
#
# 執行流程：
# 1. 驗證程式碼簽名
# 2. 檢查公證狀態
# 3. 測試應用程式啟動
verifyFinalApplication() {
    local app_path="$1"
    
    info "驗證最終應用程式..."
    
    # 驗證程式碼簽名
    info "驗證程式碼簽名..."
    codesign --verify --deep --strict --verbose=2 "$app_path" || {
        err "程式碼簽名驗證失敗"
        return 1
    }
    
    # 檢查公證狀態
    info "檢查公證狀態..."
    spctl --assess --verbose "$app_path" || {
        warn "公證狀態檢查失敗，應用程式可能未公證"
    }
    
    # 顯示應用程式資訊
    info "應用程式資訊："
    codesign -dv --verbose=4 "$app_path" 2>&1 | head -20
    
    success "應用程式驗證完成"
    return 0
}

# 主執行流程
main() {
    info "Mac 筆記本應用程式簽名和公證工具"
    info "======================================="
    
    # 檢查要求
    checkRequirements
    
    # 執行程式碼簽名
    if ! signApplication "$APP_PATH" "$DEVELOPER_ID"; then
        err "程式碼簽名失敗"
        exit 1
    fi
    
    # 執行公證（如果提供了必要資訊）
    if [[ -n "$APPLE_ID" && -n "$APP_PASSWORD" ]]; then
        if ! notarizeApplication "$APP_PATH" "$APPLE_ID" "$APP_PASSWORD"; then
            err "公證失敗"
            exit 1
        fi
    else
        warn "跳過公證流程（未提供 Apple ID 或應用程式密碼）"
    fi
    
    # 驗證最終應用程式
    if ! verifyFinalApplication "$APP_PATH"; then
        err "最終驗證失敗"
        exit 1
    fi
    
    success "簽名和公證流程完成！"
    success "應用程式位置：$APP_PATH"
    
    # 顯示使用說明
    info ""
    info "使用說明："
    info "1. 應用程式已完成簽名，可以在 macOS 上正常執行"
    if [[ -n "$APPLE_ID" && -n "$APP_PASSWORD" ]]; then
        info "2. 應用程式已完成公證，可以安全分發給其他用戶"
    else
        info "2. 如需分發給其他用戶，建議完成公證流程"
    fi
    info "3. 可以使用 'open \"$APP_PATH\"' 測試應用程式"
}

# 執行主函數
main "$@"