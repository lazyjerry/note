#!/usr/bin/env bash
set -euo pipefail

# Mac 筆記本應用程式 DMG 創建腳本
# 此腳本負責創建美觀的 macOS 安裝程式 DMG 檔案
#
# 使用方法：
# ./scripts/create_dmg.sh [應用程式路徑] [版本號]
#
# 執行流程：
# 1. 檢查必要工具和應用程式
# 2. 創建臨時 DMG 映像檔
# 3. 設定 DMG 外觀和佈局
# 4. 添加應用程式和快捷方式
# 5. 壓縮並最佳化 DMG
# 6. 清理臨時檔案

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

# 設定變數
APP_PATH="${1:-$ROOT_DIR/dist/Mac筆記本.app}"
VERSION="${2:-1.0.0}"
APP_NAME="Mac筆記本"
DMG_NAME="${APP_NAME}_v${VERSION}"
DIST_DIR="$ROOT_DIR/dist"
TEMP_DMG="$DIST_DIR/${DMG_NAME}_temp.dmg"
FINAL_DMG="$DIST_DIR/${DMG_NAME}.dmg"
VOLUME_NAME="$APP_NAME $VERSION"
DMG_SIZE="100m"

# 顏色輸出函數
info() { echo -e "\033[1;34m[INFO]\033[0m $*"; }
warn() { echo -e "\033[1;33m[WARN]\033[0m $*"; }
err()  { echo -e "\033[1;31m[ERR ]\033[0m $*"; }
success() { echo -e "\033[1;32m[SUCCESS]\033[0m $*"; }

# checkRequirements 檢查 DMG 創建所需的工具和檔案
# 回傳：如果缺少必要工具或檔案則退出
#
# 執行流程：
# 1. 檢查 macOS 系統工具
# 2. 檢查應用程式檔案
# 3. 創建輸出目錄
checkRequirements() {
    info "檢查 DMG 創建要求..."
    
    # 檢查系統工具
    for tool in hdiutil osascript; do
        if ! command -v "$tool" >/dev/null 2>&1; then
            err "未找到必要工具：$tool"
            exit 1
        fi
    done
    
    # 檢查應用程式是否存在
    if [[ ! -d "$APP_PATH" ]]; then
        err "應用程式不存在：$APP_PATH"
        err "請先執行打包腳本：./scripts/package_mac.sh"
        exit 1
    fi
    
    # 創建輸出目錄
    mkdir -p "$DIST_DIR"
    
    # 清理舊的 DMG 檔案
    rm -f "$TEMP_DMG" "$FINAL_DMG"
    
    info "檢查完成，準備創建 DMG"
}

# createTempDMG 創建臨時 DMG 映像檔
# 回傳：創建過程中的錯誤（如果有）
#
# 執行流程：
# 1. 創建空白的 DMG 映像檔
# 2. 掛載 DMG 到系統
# 3. 複製應用程式到 DMG
# 4. 創建應用程式資料夾快捷方式
createTempDMG() {
    info "創建臨時 DMG 映像檔..."
    
    # 創建空白 DMG
    hdiutil create -srcfolder /tmp -volname "$VOLUME_NAME" -fs HFS+ -fsargs "-c c=64,a=16,e=16" -format UDRW -size "$DMG_SIZE" "$TEMP_DMG" || {
        err "創建臨時 DMG 失敗"
        return 1
    }
    
    # 掛載 DMG
    info "掛載 DMG 映像檔..."
    local mount_result
    mount_result=$(hdiutil attach -readwrite -noverify -noautoopen "$TEMP_DMG" | egrep '^/dev/' | sed 1q | awk '{print $1}')
    
    if [[ -z "$mount_result" ]]; then
        err "掛載 DMG 失敗"
        return 1
    fi
    
    local volume_path="/Volumes/$VOLUME_NAME"
    
    # 等待掛載完成
    sleep 2
    
    # 複製應用程式
    info "複製應用程式到 DMG..."
    cp -R "$APP_PATH" "$volume_path/" || {
        err "複製應用程式失敗"
        hdiutil detach "$mount_result" >/dev/null 2>&1 || true
        return 1
    }
    
    # 創建應用程式資料夾快捷方式
    info "創建應用程式資料夾快捷方式..."
    ln -s /Applications "$volume_path/Applications" || {
        warn "創建應用程式快捷方式失敗"
    }
    
    # 卸載 DMG
    info "卸載臨時 DMG..."
    hdiutil detach "$mount_result" >/dev/null || {
        warn "卸載 DMG 時出現警告"
    }
    
    success "臨時 DMG 創建完成"
    return 0
}

# customizeDMGAppearance 自訂 DMG 外觀
# 回傳：自訂過程中的錯誤（如果有）
#
# 執行流程：
# 1. 重新掛載 DMG
# 2. 設定視窗大小和位置
# 3. 設定圖示位置和大小
# 4. 設定背景圖片（如果有）
# 5. 隱藏工具列和狀態列
customizeDMGAppearance() {
    info "自訂 DMG 外觀..."
    
    # 重新掛載 DMG
    local mount_result
    mount_result=$(hdiutil attach -readwrite -noverify -noautoopen "$TEMP_DMG" | egrep '^/dev/' | sed 1q | awk '{print $1}')
    
    if [[ -z "$mount_result" ]]; then
        err "重新掛載 DMG 失敗"
        return 1
    fi
    
    local volume_path="/Volumes/$VOLUME_NAME"
    
    # 等待掛載完成
    sleep 2
    
    # 使用 AppleScript 設定 Finder 視窗外觀
    info "設定 Finder 視窗外觀..."
    osascript <<EOF
tell application "Finder"
    tell disk "$VOLUME_NAME"
        open
        set current view of container window to icon view
        set toolbar visible of container window to false
        set statusbar visible of container window to false
        set the bounds of container window to {400, 100, 920, 420}
        set viewOptions to the icon view options of container window
        set arrangement of viewOptions to not arranged
        set icon size of viewOptions to 72
        set position of item "$APP_NAME.app" of container window to {160, 205}
        set position of item "Applications" of container window to {360, 205}
        close
        open
        update without registering applications
        delay 2
    end tell
end tell
EOF
    
    # 同步檔案系統
    sync
    
    # 卸載 DMG
    info "卸載 DMG..."
    hdiutil detach "$mount_result" >/dev/null || {
        warn "卸載 DMG 時出現警告"
    }
    
    success "DMG 外觀設定完成"
    return 0
}

# finalizeDMG 完成 DMG 創建
# 回傳：完成過程中的錯誤（如果有）
#
# 執行流程：
# 1. 壓縮 DMG 映像檔
# 2. 最佳化檔案大小
# 3. 驗證 DMG 完整性
# 4. 清理臨時檔案
finalizeDMG() {
    info "完成 DMG 創建..."
    
    # 壓縮並最佳化 DMG
    info "壓縮和最佳化 DMG..."
    hdiutil convert "$TEMP_DMG" -format UDZO -imagekey zlib-level=9 -o "$FINAL_DMG" || {
        err "DMG 壓縮失敗"
        return 1
    }
    
    # 驗證 DMG
    info "驗證 DMG 完整性..."
    hdiutil verify "$FINAL_DMG" || {
        warn "DMG 驗證失敗"
    }
    
    # 清理臨時檔案
    info "清理臨時檔案..."
    rm -f "$TEMP_DMG"
    
    # 顯示檔案資訊
    local dmg_size
    dmg_size=$(du -h "$FINAL_DMG" | cut -f1)
    
    success "DMG 創建完成！"
    success "檔案位置：$FINAL_DMG"
    success "檔案大小：$dmg_size"
    
    return 0
}

# createReadme 創建 README 檔案
# 回傳：創建過程中的錯誤（如果有）
createReadme() {
    local readme_path="$DIST_DIR/README_安裝說明.txt"
    
    info "創建安裝說明檔案..."
    
    cat > "$readme_path" << EOF
Mac 筆記本 v$VERSION - 安裝說明
================================

感謝您下載 Mac 筆記本！

安裝步驟：
1. 雙擊 ${DMG_NAME}.dmg 檔案開啟安裝程式
2. 將 Mac筆記本.app 拖拽到 Applications 資料夾
3. 從 Launchpad 或 Applications 資料夾啟動應用程式

系統需求：
- macOS 10.15 (Catalina) 或更新版本
- 至少 100MB 可用磁碟空間

功能特色：
- Markdown 筆記編輯和即時預覽
- 檔案和資料夾管理
- 密碼和生物識別加密保護
- 自動保存功能
- 繁體中文輸入法優化
- 智慧編輯和匯出功能

如果遇到問題：
1. 確保您的 macOS 版本符合系統需求
2. 如果出現「無法開啟，因為它來自身分不明的開發者」錯誤：
   - 右鍵點擊應用程式，選擇「開啟」
   - 或在系統偏好設定 > 安全性與隱私權中允許執行
3. 如需技術支援，請訪問：https://github.com/mac-notebook-app

版本資訊：
- 版本：$VERSION
- 建置日期：$(date '+%Y年%m月%d日')
- 相容性：macOS 10.15+

© 2024 Mac筆記本開發團隊
EOF
    
    success "安裝說明檔案已創建：$readme_path"
}

# 主執行流程
main() {
    info "Mac 筆記本 DMG 創建工具"
    info "========================="
    info "應用程式：$APP_PATH"
    info "版本：$VERSION"
    info "輸出：$FINAL_DMG"
    info ""
    
    # 檢查要求
    checkRequirements
    
    # 創建臨時 DMG
    if ! createTempDMG; then
        err "創建臨時 DMG 失敗"
        exit 1
    fi
    
    # 自訂外觀
    if ! customizeDMGAppearance; then
        warn "DMG 外觀設定失敗，但繼續完成創建"
    fi
    
    # 完成 DMG
    if ! finalizeDMG; then
        err "完成 DMG 創建失敗"
        exit 1
    fi
    
    # 創建說明檔案
    createReadme
    
    success "DMG 創建流程完成！"
    info ""
    info "產出檔案："
    info "- DMG 安裝程式：$FINAL_DMG"
    info "- 安裝說明：$DIST_DIR/README_安裝說明.txt"
    info ""
    info "測試安裝："
    info "1. 雙擊 DMG 檔案開啟"
    info "2. 拖拽應用程式到 Applications 資料夾"
    info "3. 從 Launchpad 啟動應用程式"
}

# 執行主函數
main "$@"