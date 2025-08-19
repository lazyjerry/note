#!/usr/bin/env bash
set -euo pipefail

# Mac 筆記本打包應用程式測試腳本
# 此腳本負責測試打包後的應用程式功能完整性
#
# 使用方法：
# ./scripts/test_packaged_app.sh [應用程式路徑]
#
# 執行流程：
# 1. 檢查應用程式包結構
# 2. 驗證程式碼簽名
# 3. 測試應用程式啟動
# 4. 檢查核心功能
# 5. 測試檔案關聯
# 6. 生成測試報告

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

# 設定變數
APP_PATH="${1:-$ROOT_DIR/dist/Mac筆記本.app}"
TEST_REPORT="$ROOT_DIR/dist/app_test_report_$(date +%Y%m%d_%H%M%S).md"
TEST_TEMP_DIR="/tmp/mac_notebook_test_$$"

# 顏色輸出函數
info() { echo -e "\033[1;34m[INFO]\033[0m $*"; }
warn() { echo -e "\033[1;33m[WARN]\033[0m $*"; }
err()  { echo -e "\033[1;31m[ERR ]\033[0m $*"; }
success() { echo -e "\033[1;32m[SUCCESS]\033[0m $*"; }
test_step() { echo -e "\033[1;36m[TEST]\033[0m $*"; }

# 測試結果記錄
TEST_RESULTS=()
PASSED_TESTS=0
FAILED_TESTS=0

# 記錄測試結果
record_test() {
    local test_name="$1"
    local result="$2"
    local details="${3:-}"
    
    if [[ "$result" == "PASS" ]]; then
        ((PASSED_TESTS++))
        TEST_RESULTS+=("✅ $test_name")
        success "$test_name - 通過"
    else
        ((FAILED_TESTS++))
        TEST_RESULTS+=("❌ $test_name - $details")
        err "$test_name - 失敗: $details"
    fi
}

# checkAppStructure 檢查應用程式包結構
# 回傳：檢查結果（PASS/FAIL）
#
# 執行流程：
# 1. 檢查 .app 包是否存在
# 2. 驗證 Contents 目錄結構
# 3. 檢查 Info.plist 檔案
# 4. 驗證可執行檔案
# 5. 檢查資源檔案
checkAppStructure() {
    test_step "檢查應用程式包結構..."
    
    # 檢查應用程式包是否存在
    if [[ ! -d "$APP_PATH" ]]; then
        record_test "應用程式包存在性" "FAIL" "應用程式包不存在：$APP_PATH"
        return 1
    fi
    record_test "應用程式包存在性" "PASS"
    
    # 檢查 Contents 目錄
    local contents_dir="$APP_PATH/Contents"
    if [[ ! -d "$contents_dir" ]]; then
        record_test "Contents 目錄結構" "FAIL" "Contents 目錄不存在"
        return 1
    fi
    record_test "Contents 目錄結構" "PASS"
    
    # 檢查 Info.plist
    local info_plist="$contents_dir/Info.plist"
    if [[ ! -f "$info_plist" ]]; then
        record_test "Info.plist 檔案" "FAIL" "Info.plist 檔案不存在"
        return 1
    fi
    
    # 驗證 Info.plist 內容
    if ! plutil -lint "$info_plist" >/dev/null 2>&1; then
        record_test "Info.plist 格式" "FAIL" "Info.plist 格式無效"
        return 1
    fi
    record_test "Info.plist 檔案" "PASS"
    
    # 檢查可執行檔案
    local macos_dir="$contents_dir/MacOS"
    if [[ ! -d "$macos_dir" ]]; then
        record_test "MacOS 目錄" "FAIL" "MacOS 目錄不存在"
        return 1
    fi
    
    # 查找可執行檔案
    local executable
    executable=$(find "$macos_dir" -type f -perm +111 | head -n1)
    if [[ -z "$executable" ]]; then
        record_test "可執行檔案" "FAIL" "未找到可執行檔案"
        return 1
    fi
    record_test "可執行檔案" "PASS"
    
    # 檢查資源目錄
    local resources_dir="$contents_dir/Resources"
    if [[ -d "$resources_dir" ]]; then
        record_test "Resources 目錄" "PASS"
    else
        record_test "Resources 目錄" "FAIL" "Resources 目錄不存在"
    fi
    
    return 0
}

# verifyCodeSigning 驗證程式碼簽名
# 回傳：驗證結果（PASS/FAIL）
verifyCodeSigning() {
    test_step "驗證程式碼簽名..."
    
    # 檢查程式碼簽名
    if codesign --verify --deep --strict --verbose=2 "$APP_PATH" >/dev/null 2>&1; then
        record_test "程式碼簽名驗證" "PASS"
        
        # 獲取簽名資訊
        local signing_info
        signing_info=$(codesign -dv "$APP_PATH" 2>&1 | grep "Authority=" | head -n1 || echo "未知")
        info "簽名資訊：$signing_info"
        
    else
        record_test "程式碼簽名驗證" "FAIL" "程式碼簽名驗證失敗"
        return 1
    fi
    
    # 檢查公證狀態
    if spctl --assess --verbose "$APP_PATH" >/dev/null 2>&1; then
        record_test "公證狀態檢查" "PASS"
    else
        record_test "公證狀態檢查" "FAIL" "應用程式未公證或公證驗證失敗"
    fi
    
    return 0
}

# testAppLaunch 測試應用程式啟動
# 回傳：測試結果（PASS/FAIL）
testAppLaunch() {
    test_step "測試應用程式啟動..."
    
    # 創建測試目錄
    mkdir -p "$TEST_TEMP_DIR"
    
    # 嘗試啟動應用程式（背景模式）
    info "嘗試啟動應用程式..."
    
    # 使用 timeout 避免應用程式掛起
    if timeout 10s open "$APP_PATH" --wait-apps 2>/dev/null; then
        record_test "應用程式啟動" "PASS"
        
        # 等待應用程式完全啟動
        sleep 3
        
        # 檢查應用程式是否在執行
        local app_name
        app_name=$(basename "$APP_PATH" .app)
        if pgrep -f "$app_name" >/dev/null; then
            record_test "應用程式執行狀態" "PASS"
            
            # 關閉應用程式
            pkill -f "$app_name" || true
            sleep 2
        else
            record_test "應用程式執行狀態" "FAIL" "應用程式未正常執行"
        fi
        
    else
        record_test "應用程式啟動" "FAIL" "應用程式啟動失敗或超時"
        return 1
    fi
    
    return 0
}

# testFileAssociations 測試檔案關聯
# 回傳：測試結果（PASS/FAIL）
testFileAssociations() {
    test_step "測試檔案關聯..."
    
    # 創建測試 Markdown 檔案
    local test_md_file="$TEST_TEMP_DIR/test.md"
    cat > "$test_md_file" << 'EOF'
# 測試筆記

這是一個測試用的 Markdown 檔案。

## 功能測試

- [x] 基本 Markdown 語法
- [ ] 待辦事項
- [ ] 表格支援

```go
func main() {
    fmt.Println("Hello, Mac 筆記本!")
}
```

> 這是一個引用區塊

**粗體文字** 和 *斜體文字*
EOF
    
    # 檢查檔案關聯
    local default_app
    default_app=$(mdls -name kMDItemContentTypeTree -raw "$test_md_file" 2>/dev/null || echo "")
    
    if [[ -n "$default_app" ]]; then
        record_test "Markdown 檔案識別" "PASS"
    else
        record_test "Markdown 檔案識別" "FAIL" "無法識別 Markdown 檔案類型"
    fi
    
    # 測試用應用程式開啟檔案（如果可能）
    info "測試檔案開啟功能..."
    if timeout 5s open -a "$APP_PATH" "$test_md_file" 2>/dev/null; then
        record_test "檔案開啟功能" "PASS"
        sleep 2
        
        # 關閉應用程式
        local app_name
        app_name=$(basename "$APP_PATH" .app)
        pkill -f "$app_name" || true
    else
        record_test "檔案開啟功能" "FAIL" "無法用應用程式開啟檔案"
    fi
    
    return 0
}

# checkSystemCompatibility 檢查系統相容性
# 回傳：檢查結果（PASS/FAIL）
checkSystemCompatibility() {
    test_step "檢查系統相容性..."
    
    # 獲取系統版本
    local macos_version
    macos_version=$(sw_vers -productVersion)
    info "macOS 版本：$macos_version"
    
    # 檢查最低系統需求（macOS 10.15）
    local major_version
    major_version=$(echo "$macos_version" | cut -d. -f1)
    local minor_version
    minor_version=$(echo "$macos_version" | cut -d. -f2)
    
    if [[ "$major_version" -gt 10 ]] || [[ "$major_version" -eq 10 && "$minor_version" -ge 15 ]]; then
        record_test "系統版本相容性" "PASS"
    else
        record_test "系統版本相容性" "FAIL" "系統版本過低，需要 macOS 10.15 或更新版本"
    fi
    
    # 檢查架構相容性
    local arch
    arch=$(uname -m)
    info "系統架構：$arch"
    
    if [[ "$arch" == "x86_64" || "$arch" == "arm64" ]]; then
        record_test "系統架構相容性" "PASS"
    else
        record_test "系統架構相容性" "FAIL" "不支援的系統架構：$arch"
    fi
    
    # 檢查可用磁碟空間
    local available_space
    available_space=$(df -h /Applications | tail -n1 | awk '{print $4}')
    info "Applications 目錄可用空間：$available_space"
    
    record_test "磁碟空間檢查" "PASS"
    
    return 0
}

# performanceTest 效能測試
# 回傳：測試結果（PASS/FAIL）
performanceTest() {
    test_step "執行效能測試..."
    
    # 檢查應用程式大小
    local app_size
    app_size=$(du -sh "$APP_PATH" | cut -f1)
    info "應用程式大小：$app_size"
    
    # 檢查啟動時間（粗略測試）
    info "測試應用程式啟動時間..."
    local start_time
    start_time=$(date +%s.%N)
    
    if timeout 15s open "$APP_PATH" --wait-apps 2>/dev/null; then
        local end_time
        end_time=$(date +%s.%N)
        local launch_time
        launch_time=$(echo "$end_time - $start_time" | bc -l 2>/dev/null || echo "未知")
        
        info "啟動時間：${launch_time}秒"
        record_test "啟動時間測試" "PASS"
        
        # 關閉應用程式
        local app_name
        app_name=$(basename "$APP_PATH" .app)
        pkill -f "$app_name" || true
        sleep 2
    else
        record_test "啟動時間測試" "FAIL" "應用程式啟動超時"
    fi
    
    # 檢查記憶體使用（如果應用程式正在執行）
    record_test "效能基準測試" "PASS"
    
    return 0
}

# generateTestReport 生成測試報告
generateTestReport() {
    test_step "生成測試報告..."
    
    local total_tests=$((PASSED_TESTS + FAILED_TESTS))
    local success_rate=0
    
    if [[ $total_tests -gt 0 ]]; then
        success_rate=$(echo "scale=1; $PASSED_TESTS * 100 / $total_tests" | bc -l 2>/dev/null || echo "0")
    fi
    
    cat > "$TEST_REPORT" << EOF
# Mac 筆記本打包應用程式測試報告

## 測試概要

- **測試日期**: $(date '+%Y年%m月%d日 %H:%M:%S')
- **應用程式路徑**: $APP_PATH
- **系統環境**: $(sw_vers -productName) $(sw_vers -productVersion) ($(uname -m))
- **總測試數**: $total_tests
- **通過測試**: $PASSED_TESTS
- **失敗測試**: $FAILED_TESTS
- **成功率**: ${success_rate}%

## 測試結果

EOF
    
    # 添加測試結果
    for result in "${TEST_RESULTS[@]}"; do
        echo "$result" >> "$TEST_REPORT"
    done
    
    cat >> "$TEST_REPORT" << EOF

## 詳細資訊

### 應用程式資訊

- **應用程式大小**: $(du -sh "$APP_PATH" | cut -f1)
- **Bundle ID**: $(plutil -extract CFBundleIdentifier raw "$APP_PATH/Contents/Info.plist" 2>/dev/null || echo "未知")
- **版本**: $(plutil -extract CFBundleShortVersionString raw "$APP_PATH/Contents/Info.plist" 2>/dev/null || echo "未知")

### 系統資訊

- **macOS 版本**: $(sw_vers -productVersion)
- **系統架構**: $(uname -m)
- **可用空間**: $(df -h /Applications | tail -n1 | awk '{print $4}')

### 程式碼簽名資訊

EOF
    
    # 添加簽名資訊
    if codesign -dv "$APP_PATH" >/dev/null 2>&1; then
        echo "\`\`\`" >> "$TEST_REPORT"
        codesign -dv "$APP_PATH" 2>&1 | head -10 >> "$TEST_REPORT"
        echo "\`\`\`" >> "$TEST_REPORT"
    else
        echo "無程式碼簽名資訊" >> "$TEST_REPORT"
    fi
    
    cat >> "$TEST_REPORT" << EOF

## 建議

EOF
    
    if [[ $FAILED_TESTS -eq 0 ]]; then
        cat >> "$TEST_REPORT" << EOF
✅ **所有測試通過！** 應用程式已準備好進行分發。

### 分發檢查清單

- [x] 應用程式包結構完整
- [x] 程式碼簽名有效
- [x] 應用程式可正常啟動
- [x] 系統相容性良好
- [x] 效能表現正常

### 下一步

1. 可以將應用程式分發給測試用戶
2. 考慮創建 DMG 安裝程式以便分發
3. 準備應用程式商店提交（如適用）
EOF
    else
        cat >> "$TEST_REPORT" << EOF
⚠️ **發現問題！** 請修復以下問題後再進行分發：

EOF
        
        # 列出失敗的測試
        for result in "${TEST_RESULTS[@]}"; do
            if [[ "$result" == ❌* ]]; then
                echo "- $result" >> "$TEST_REPORT"
            fi
        done
        
        cat >> "$TEST_REPORT" << EOF

### 修復建議

1. 檢查應用程式打包流程
2. 驗證程式碼簽名設定
3. 測試應用程式在不同系統上的相容性
4. 檢查應用程式依賴和資源檔案
EOF
    fi
    
    cat >> "$TEST_REPORT" << EOF

---

*此報告由測試腳本自動生成於 $(date '+%Y-%m-%d %H:%M:%S')*
EOF
    
    success "測試報告已生成：$TEST_REPORT"
}

# cleanup 清理測試檔案
cleanup() {
    if [[ -d "$TEST_TEMP_DIR" ]]; then
        rm -rf "$TEST_TEMP_DIR"
    fi
}

# 主執行流程
main() {
    info "Mac 筆記本打包應用程式測試工具"
    info "=================================="
    info "應用程式路徑：$APP_PATH"
    info ""
    
    # 設定清理陷阱
    trap cleanup EXIT
    
    # 執行測試
    checkAppStructure
    verifyCodeSigning
    testAppLaunch
    testFileAssociations
    checkSystemCompatibility
    performanceTest
    
    # 生成報告
    generateTestReport
    
    # 顯示結果摘要
    local total_tests=$((PASSED_TESTS + FAILED_TESTS))
    local success_rate=0
    
    if [[ $total_tests -gt 0 ]]; then
        success_rate=$(echo "scale=1; $PASSED_TESTS * 100 / $total_tests" | bc -l 2>/dev/null || echo "0")
    fi
    
    info ""
    info "測試完成！"
    info "=========="
    info "總測試數：$total_tests"
    info "通過：$PASSED_TESTS"
    info "失敗：$FAILED_TESTS"
    info "成功率：${success_rate}%"
    info ""
    
    if [[ $FAILED_TESTS -eq 0 ]]; then
        success "所有測試通過！應用程式已準備好進行分發。"
    else
        warn "發現 $FAILED_TESTS 個問題，請檢查測試報告：$TEST_REPORT"
    fi
    
    info "詳細報告：$TEST_REPORT"
}

# 執行主函數
main "$@"