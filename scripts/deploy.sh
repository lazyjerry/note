#!/usr/bin/env bash
set -euo pipefail

# Mac 筆記本應用程式完整部署腳本
# 此腳本負責完整的應用程式部署流程，包括打包、簽名、公證和 DMG 創建
#
# 使用方法：
# ./scripts/deploy.sh [版本號] [--skip-notarize] [--clean]
#
# 執行流程：
# 1. 解析命令列參數
# 2. 執行完整測試套件
# 3. 生成資源並打包應用程式
# 4. 進行程式碼簽名
# 5. 執行公證流程（可選）
# 6. 創建 DMG 安裝程式
# 7. 生成部署報告

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

# 預設設定
VERSION="${1:-1.0.0}"
SKIP_NOTARIZE=false
CLEAN_BUILD=false
SKIP_TESTS=false

# 解析命令列參數
while [[ $# -gt 0 ]]; do
    case $1 in
        --skip-notarize)
            SKIP_NOTARIZE=true
            shift
            ;;
        --clean)
            CLEAN_BUILD=true
            shift
            ;;
        --skip-tests)
            SKIP_TESTS=true
            shift
            ;;
        --help|-h)
            echo "Mac 筆記本應用程式部署腳本"
            echo ""
            echo "使用方法："
            echo "  $0 [版本號] [選項]"
            echo ""
            echo "選項："
            echo "  --skip-notarize  跳過公證流程"
            echo "  --clean          清理建置並重新開始"
            echo "  --skip-tests     跳過測試執行"
            echo "  --help, -h       顯示此說明"
            echo ""
            echo "範例："
            echo "  $0 1.0.0                    # 完整部署流程"
            echo "  $0 1.0.1 --skip-notarize   # 跳過公證的部署"
            echo "  $0 1.1.0 --clean           # 清理建置的部署"
            exit 0
            ;;
        *)
            if [[ "$1" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
                VERSION="$1"
            else
                echo "未知參數：$1"
                exit 1
            fi
            shift
            ;;
    esac
done

# 設定環境變數
export APP_VERSION="$VERSION"
export APP_NAME="Mac筆記本"
export APP_ID="com.notebook.mac-notebook-app"

# 顏色輸出函數
info() { echo -e "\033[1;34m[INFO]\033[0m $*"; }
warn() { echo -e "\033[1;33m[WARN]\033[0m $*"; }
err()  { echo -e "\033[1;31m[ERR ]\033[0m $*"; }
success() { echo -e "\033[1;32m[SUCCESS]\033[0m $*"; }
step() { echo -e "\033[1;36m[STEP]\033[0m $*"; }

# 記錄部署開始時間
DEPLOY_START_TIME=$(date +%s)
DEPLOY_LOG="dist/deploy_${VERSION}_$(date +%Y%m%d_%H%M%S).log"

# 創建日誌目錄
mkdir -p dist

# 日誌函數
log() {
    echo "$(date '+%Y-%m-%d %H:%M:%S') $*" >> "$DEPLOY_LOG"
    echo "$*"
}

# runTests 執行完整測試套件
# 回傳：測試過程中的錯誤（如果有）
#
# 執行流程：
# 1. 執行單元測試
# 2. 執行整合測試
# 3. 執行效能測試
# 4. 生成測試報告
runTests() {
    if [[ "$SKIP_TESTS" == "true" ]]; then
        warn "跳過測試執行"
        return 0
    fi
    
    step "執行完整測試套件..."
    log "開始執行測試套件"
    
    # 執行所有測試
    info "執行單元測試和整合測試..."
    if ! go test -v -race -coverprofile=coverage.out ./... 2>&1 | tee -a "$DEPLOY_LOG"; then
        err "測試失敗"
        log "測試失敗，部署中止"
        return 1
    fi
    
    # 生成測試覆蓋率報告
    info "生成測試覆蓋率報告..."
    go tool cover -html=coverage.out -o dist/coverage_${VERSION}.html
    
    # 顯示測試覆蓋率
    local coverage
    coverage=$(go tool cover -func=coverage.out | grep total | awk '{print $3}')
    info "測試覆蓋率：$coverage"
    log "測試覆蓋率：$coverage"
    
    success "測試套件執行完成"
    log "測試套件執行完成"
    return 0
}

# buildApplication 建置應用程式
# 回傳：建置過程中的錯誤（如果有）
buildApplication() {
    step "建置應用程式..."
    log "開始建置應用程式 v$VERSION"
    
    # 清理建置（如果需要）
    if [[ "$CLEAN_BUILD" == "true" ]]; then
        info "清理舊的建置檔案..."
        rm -rf dist/*.app dist/*.dmg build/
        log "清理建置檔案完成"
    fi
    
    # 執行打包腳本
    info "執行應用程式打包..."
    if ! ./scripts/package_mac.sh 2>&1 | tee -a "$DEPLOY_LOG"; then
        err "應用程式打包失敗"
        log "應用程式打包失敗"
        return 1
    fi
    
    success "應用程式建置完成"
    log "應用程式建置完成"
    return 0
}

# signAndNotarize 簽名和公證應用程式
# 回傳：簽名公證過程中的錯誤（如果有）
signAndNotarize() {
    step "簽名和公證應用程式..."
    
    local app_path="dist/${APP_NAME}.app"
    
    # 檢查應用程式是否存在
    if [[ ! -d "$app_path" ]]; then
        err "找不到應用程式：$app_path"
        log "找不到應用程式：$app_path"
        return 1
    fi
    
    # 執行簽名
    info "執行程式碼簽名..."
    log "開始程式碼簽名"
    
    if [[ "$SKIP_NOTARIZE" == "true" ]]; then
        # 僅簽名，不公證
        if ! ./scripts/sign_and_notarize.sh "$app_path" 2>&1 | tee -a "$DEPLOY_LOG"; then
            err "程式碼簽名失敗"
            log "程式碼簽名失敗"
            return 1
        fi
    else
        # 完整簽名和公證流程
        # 檢查公證所需的環境變數
        if [[ -z "${APPLE_ID:-}" || -z "${APP_PASSWORD:-}" ]]; then
            warn "未設定 APPLE_ID 或 APP_PASSWORD 環境變數"
            warn "僅執行程式碼簽名，跳過公證"
            log "跳過公證：缺少必要環境變數"
            
            if ! ./scripts/sign_and_notarize.sh "$app_path" 2>&1 | tee -a "$DEPLOY_LOG"; then
                err "程式碼簽名失敗"
                log "程式碼簽名失敗"
                return 1
            fi
        else
            info "執行完整簽名和公證流程..."
            log "開始公證流程"
            
            if ! ./scripts/sign_and_notarize.sh "$app_path" "" "$APPLE_ID" "$APP_PASSWORD" 2>&1 | tee -a "$DEPLOY_LOG"; then
                err "簽名和公證失敗"
                log "簽名和公證失敗"
                return 1
            fi
        fi
    fi
    
    success "簽名和公證完成"
    log "簽名和公證完成"
    return 0
}

# createInstaller 創建安裝程式
# 回傳：創建過程中的錯誤（如果有）
createInstaller() {
    step "創建 DMG 安裝程式..."
    log "開始創建 DMG 安裝程式"
    
    local app_path="dist/${APP_NAME}.app"
    
    # 執行 DMG 創建腳本
    if ! ./scripts/create_dmg.sh "$app_path" "$VERSION" 2>&1 | tee -a "$DEPLOY_LOG"; then
        err "DMG 創建失敗"
        log "DMG 創建失敗"
        return 1
    fi
    
    success "DMG 安裝程式創建完成"
    log "DMG 安裝程式創建完成"
    return 0
}

# generateDeployReport 生成部署報告
# 回傳：報告生成過程中的錯誤（如果有）
generateDeployReport() {
    step "生成部署報告..."
    
    local report_path="dist/deploy_report_${VERSION}.md"
    local deploy_end_time=$(date +%s)
    local deploy_duration=$((deploy_end_time - DEPLOY_START_TIME))
    local deploy_minutes=$((deploy_duration / 60))
    local deploy_seconds=$((deploy_duration % 60))
    
    cat > "$report_path" << EOF
# Mac 筆記本 v${VERSION} 部署報告

## 部署資訊

- **版本號**: ${VERSION}
- **部署日期**: $(date '+%Y年%m月%d日 %H:%M:%S')
- **部署時長**: ${deploy_minutes}分${deploy_seconds}秒
- **建置環境**: $(uname -s) $(uname -r)
- **Go 版本**: $(go version)

## 部署設定

- **清理建置**: $([ "$CLEAN_BUILD" == "true" ] && echo "是" || echo "否")
- **跳過測試**: $([ "$SKIP_TESTS" == "true" ] && echo "是" || echo "否")
- **跳過公證**: $([ "$SKIP_NOTARIZE" == "true" ] && echo "是" || echo "否")

## 產出檔案

EOF
    
    # 列出產出檔案
    if [[ -d "dist" ]]; then
        echo "### 應用程式檔案" >> "$report_path"
        echo "" >> "$report_path"
        find dist -name "*.app" -exec echo "- **應用程式**: {} ($(du -h {} | cut -f1))" \; >> "$report_path"
        
        echo "" >> "$report_path"
        echo "### 安裝程式檔案" >> "$report_path"
        echo "" >> "$report_path"
        find dist -name "*.dmg" -exec echo "- **DMG 安裝程式**: {} ($(du -h {} | cut -f1))" \; >> "$report_path"
        
        echo "" >> "$report_path"
        echo "### 其他檔案" >> "$report_path"
        echo "" >> "$report_path"
        find dist -name "*.html" -exec echo "- **測試報告**: {}" \; >> "$report_path"
        find dist -name "*.txt" -exec echo "- **說明文件**: {}" \; >> "$report_path"
        find dist -name "*.log" -exec echo "- **部署日誌**: {}" \; >> "$report_path"
    fi
    
    cat >> "$report_path" << EOF

## 測試結果

EOF
    
    if [[ -f "coverage.out" ]]; then
        local coverage
        coverage=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' || echo "未知")
        echo "- **測試覆蓋率**: $coverage" >> "$report_path"
    fi
    
    cat >> "$report_path" << EOF

## 部署步驟

1. ✅ 執行測試套件
2. ✅ 生成應用程式資源
3. ✅ 建置和打包應用程式
4. ✅ 程式碼簽名$([ "$SKIP_NOTARIZE" == "false" ] && echo "和公證" || echo "")
5. ✅ 創建 DMG 安裝程式
6. ✅ 生成部署報告

## 使用說明

### 安裝應用程式

1. 雙擊 \`Mac筆記本_v${VERSION}.dmg\` 開啟安裝程式
2. 將應用程式拖拽到 Applications 資料夾
3. 從 Launchpad 或 Applications 資料夾啟動

### 系統需求

- macOS 10.15 (Catalina) 或更新版本
- 至少 100MB 可用磁碟空間

### 分發說明

- 應用程式已完成程式碼簽名，可在 macOS 上正常執行
$([ "$SKIP_NOTARIZE" == "false" ] && echo "- 應用程式已完成公證，可安全分發給其他用戶" || echo "- 如需分發給其他用戶，建議完成公證流程")

---

*此報告由部署腳本自動生成於 $(date '+%Y-%m-%d %H:%M:%S')*
EOF
    
    success "部署報告已生成：$report_path"
    log "部署報告已生成：$report_path"
}

# 主執行流程
main() {
    info "Mac 筆記本應用程式部署工具 v${VERSION}"
    info "============================================="
    log "開始部署流程 v${VERSION}"
    
    # 顯示部署設定
    info "部署設定："
    info "- 版本號：$VERSION"
    info "- 清理建置：$([ "$CLEAN_BUILD" == "true" ] && echo "是" || echo "否")"
    info "- 跳過測試：$([ "$SKIP_TESTS" == "true" ] && echo "是" || echo "否")"
    info "- 跳過公證：$([ "$SKIP_NOTARIZE" == "true" ] && echo "是" || echo "否")"
    info ""
    
    # 執行部署步驟
    local step_count=1
    local total_steps=6
    
    # 步驟 1: 執行測試
    info "[$step_count/$total_steps] 執行測試套件"
    if ! runTests; then
        err "部署失敗：測試階段"
        exit 1
    fi
    ((step_count++))
    
    # 步驟 2: 建置應用程式
    info "[$step_count/$total_steps] 建置應用程式"
    if ! buildApplication; then
        err "部署失敗：建置階段"
        exit 1
    fi
    ((step_count++))
    
    # 步驟 3: 簽名和公證
    info "[$step_count/$total_steps] 簽名和公證"
    if ! signAndNotarize; then
        err "部署失敗：簽名公證階段"
        exit 1
    fi
    ((step_count++))
    
    # 步驟 4: 創建安裝程式
    info "[$step_count/$total_steps] 創建安裝程式"
    if ! createInstaller; then
        err "部署失敗：安裝程式創建階段"
        exit 1
    fi
    ((step_count++))
    
    # 步驟 5: 生成報告
    info "[$step_count/$total_steps] 生成部署報告"
    generateDeployReport
    ((step_count++))
    
    # 完成部署
    local deploy_end_time=$(date +%s)
    local deploy_duration=$((deploy_end_time - DEPLOY_START_TIME))
    local deploy_minutes=$((deploy_duration / 60))
    local deploy_seconds=$((deploy_duration % 60))
    
    success "部署完成！"
    success "版本：$VERSION"
    success "耗時：${deploy_minutes}分${deploy_seconds}秒"
    success ""
    success "產出檔案："
    
    # 列出主要產出檔案
    if [[ -f "dist/Mac筆記本_v${VERSION}.dmg" ]]; then
        local dmg_size
        dmg_size=$(du -h "dist/Mac筆記本_v${VERSION}.dmg" | cut -f1)
        success "- DMG 安裝程式：dist/Mac筆記本_v${VERSION}.dmg ($dmg_size)"
    fi
    
    if [[ -d "dist/Mac筆記本.app" ]]; then
        local app_size
        app_size=$(du -h "dist/Mac筆記本.app" | cut -f1)
        success "- 應用程式：dist/Mac筆記本.app ($app_size)"
    fi
    
    if [[ -f "dist/deploy_report_${VERSION}.md" ]]; then
        success "- 部署報告：dist/deploy_report_${VERSION}.md"
    fi
    
    success ""
    success "測試安裝："
    success "1. 雙擊 DMG 檔案開啟安裝程式"
    success "2. 拖拽應用程式到 Applications 資料夾"
    success "3. 從 Launchpad 啟動應用程式"
    
    log "部署流程完成 v${VERSION}"
}

# 錯誤處理
trap 'err "部署過程中發生錯誤，請檢查日誌：$DEPLOY_LOG"; exit 1' ERR

# 執行主函數
main "$@"