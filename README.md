# Mac 加密 Notebook

一個功能完整的安全 Markdown 筆記應用程式，使用 Go 和 Fyne 建立，專為 macOS 設計，支援加密、自動保存和效能優化。

## 🌟 功能特色

### 核心功能

- **Markdown 編輯與即時預覽** - 支援 GitHub Flavored Markdown，包含表格、任務列表、程式碼高亮
- **檔案和資料夾管理** - 完整的檔案系統操作，支援建立、刪除、重新命名、移動和複製
- **密碼和生物識別加密** - AES-256 加密，支援 Touch ID/Face ID 生物識別驗證
- **自動保存功能** - 智慧自動保存，支援自訂間隔和加密檔案處理
- **拖放支援** - 支援檔案和資料夾的拖放操作

### 進階功能

- **效能監控** - 即時記憶體使用監控、背景任務追蹤和效能優化
- **錯誤處理系統** - 統一的錯誤處理、本地化錯誤訊息和日誌記錄
- **通知系統** - 完整的用戶通知和保存狀態指示器
- **主題支援** - 支援淺色和深色主題切換
- **大檔案處理** - 分塊處理大檔案，優化記憶體使用

## 📁 專案結構

```
mac-notebook-app/
├── main.go                     # 應用程式入口點
├── go.mod                      # Go 模組定義
├── go.sum                      # 相依套件版本鎖定
├── CHANGELOG.md                # 變更日誌
├── README.md                   # 專案說明文件
├── internal/                   # 內部應用程式程式碼
│   ├── models/                 # 資料模型和結構體
│   │   ├── note.go            # 筆記模型
│   │   ├── settings.go        # 設定模型
│   │   ├── file_info.go       # 檔案資訊模型
│   │   └── errors.go          # 錯誤定義
│   ├── services/              # 業務邏輯服務（15+ 個服務）
│   │   ├── interfaces.go      # 服務介面定義
│   │   ├── editor_service.go  # 編輯器服務
│   │   ├── file_manager_service.go # 檔案管理服務
│   │   ├── encryption_service.go   # 加密服務
│   │   ├── auto_save_service.go    # 自動保存服務
│   │   ├── performance_service.go  # 效能監控服務
│   │   ├── error_service.go        # 錯誤處理服務
│   │   ├── notification_service.go # 通知服務
│   │   └── ...                     # 其他服務
│   └── repositories/          # 資料存取層
│       ├── interfaces.go      # 儲存庫介面
│       └── file_repository.go # 檔案儲存庫實作
└── ui/                        # 使用者介面元件（10+ 個 UI 元件）
    ├── main_window.go         # 主視窗實作
    ├── editor.go              # 編輯器元件
    ├── preview.go             # 預覽元件
    ├── file_tree.go           # 檔案樹元件
    ├── settings_dialog.go     # 設定對話框
    └── ...                    # 其他 UI 元件
```

## 🔧 系統需求

### 最低需求

- **作業系統**: macOS 10.15 (Catalina) 或更新版本
- **Go 版本**: Go 1.21 或更新版本
- **記憶體**: 最少 512MB RAM
- **儲存空間**: 最少 100MB 可用空間

### 建議需求

- **作業系統**: macOS 12.0 (Monterey) 或更新版本
- **記憶體**: 1GB RAM 或更多
- **儲存空間**: 500MB 可用空間

## 📦 相依套件

```go
require (
    fyne.io/fyne/v2 v2.4.0+        // 跨平台 GUI 框架
    github.com/yuin/goldmark v1.6.0+ // Markdown 解析器
    github.com/fsnotify/fsnotify v1.7.0+ // 檔案系統通知
    golang.org/x/crypto v0.17.0+    // 加密函數
    github.com/google/uuid v1.4.0+  // UUID 生成
)
```

## 🚀 快速開始

### 1. 環境準備

```bash
# 確認 Go 版本
go version  # 需要 Go 1.21+

# 安裝 Xcode Command Line Tools（macOS 必需）
xcode-select --install
```

### 2. 下載專案

```bash
# 克隆專案
git clone https://github.com/your-username/mac-notebook-app.git
cd mac-notebook-app

# 下載相依套件
go mod download
```

### 3. 開發模式運行

```bash
# 運行應用程式
go run main.go

# 或者使用 Fyne 工具
go install fyne.io/fyne/v2/cmd/fyne@latest
fyne run -os darwin
```

### 4. 執行測試

```bash
# 執行所有測試
go test ./...

# 執行特定套件測試
go test ./internal/services

# 執行效能測試
go test -bench=. ./internal/services

# 執行整合測試
go test -v ./internal/services -run TestEndToEnd

# 執行長時間穩定性測試
go test -v ./internal/services -run TestLongRunning
```

## 📱 打包和部署

### 方法一：使用 Fyne 工具打包（推薦）

```bash
# 安裝 Fyne 打包工具
go install fyne.io/fyne/v2/cmd/fyne@latest

# 打包為 macOS 應用程式
fyne package -os darwin -name "Mac Notebook" -icon icon.png

# 打包並指定應用程式資訊
fyne package -os darwin \
  -name "Mac Notebook" \
  -icon icon.png \
  -appID "com.yourcompany.macnotebook" \
  -appVersion "0.14.0" \
  -appBuild "1"
```

### 方法二：手動建置

```bash
# 建置可執行檔
CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -o mac-notebook-amd64 main.go
CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -o mac-notebook-arm64 main.go

# 建立通用二進位檔（Universal Binary）
lipo -create -output mac-notebook mac-notebook-amd64 mac-notebook-arm64

# 建立應用程式包結構
mkdir -p "Mac Notebook.app/Contents/MacOS"
mkdir -p "Mac Notebook.app/Contents/Resources"

# 複製執行檔
cp mac-notebook "Mac Notebook.app/Contents/MacOS/"

# 建立 Info.plist
cat > "Mac Notebook.app/Contents/Info.plist" << EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleExecutable</key>
    <string>mac-notebook</string>
    <key>CFBundleIdentifier</key>
    <string>com.yourcompany.macnotebook</string>
    <key>CFBundleName</key>
    <string>Mac Notebook</string>
    <key>CFBundleVersion</key>
    <string>0.14.0</string>
    <key>CFBundleShortVersionString</key>
    <string>0.14.0</string>
    <key>CFBundlePackageType</key>
    <string>APPL</string>
    <key>LSMinimumSystemVersion</key>
    <string>10.15</string>
    <key>NSHighResolutionCapable</key>
    <true/>
</dict>
</plist>
EOF
```

### 方法三：建立安裝程式

```bash
# 使用 create-dmg 建立 DMG 安裝檔
brew install create-dmg

# 建立 DMG
create-dmg \
  --volname "Mac Notebook Installer" \
  --volicon "icon.icns" \
  --window-pos 200 120 \
  --window-size 600 300 \
  --icon-size 100 \
  --icon "Mac Notebook.app" 175 120 \
  --hide-extension "Mac Notebook.app" \
  --app-drop-link 425 120 \
  "Mac Notebook v0.14.0.dmg" \
  "Mac Notebook.app"
```

### 程式碼簽名和公證（發布用）

```bash
# 程式碼簽名
codesign --force --options runtime --sign "Developer ID Application: Your Name" "Mac Notebook.app"

# 建立簽名的 DMG
codesign --force --sign "Developer ID Application: Your Name" "Mac Notebook v0.14.0.dmg"

# 上傳到 Apple 進行公證
xcrun notarytool submit "Mac Notebook v0.14.0.dmg" \
  --apple-id "your-apple-id@example.com" \
  --password "app-specific-password" \
  --team-id "YOUR_TEAM_ID" \
  --wait

# 裝訂公證票據
xcrun stapler staple "Mac Notebook v0.14.0.dmg"
```

## 🧪 測試覆蓋率

專案包含完整的測試套件：

- **單元測試**: 150+ 個測試函數
- **整合測試**: 6 個端到端測試場景
- **效能測試**: 10+ 個基準測試
- **穩定性測試**: 長時間運行和記憶體洩漏檢測
- **測試覆蓋率**: 85%+ 程式碼覆蓋率

```bash
# 檢查測試覆蓋率
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

## 📊 專案統計

- **總檔案數**: 45+ 個 Go 檔案
- **程式碼行數**: 15,000+ 行（包含註解）
- **測試檔案數量**: 20+ 個測試檔案
- **服務模組**: 15+ 個業務邏輯服務
- **UI 元件**: 10+ 個使用者介面元件
- **架構完成度**: 95% - 核心功能完成，準備發布

## 🔒 安全性功能

- **AES-256 加密**: 軍用級別的檔案加密
- **生物識別驗證**: Touch ID/Face ID 支援
- **密碼強度驗證**: 強制使用強密碼
- **安全檔案處理**: 防止目錄遍歷攻擊
- **記憶體安全**: 自動清理敏感資料

## 🚀 效能特色

- **記憶體優化**: 智慧記憶體管理和垃圾回收
- **大檔案支援**: 分塊處理 5MB+ 檔案
- **快取系統**: LRU 快取策略，最多快取 100 個筆記
- **背景監控**: 即時效能指標和資源使用監控
- **並發安全**: 執行緒安全的並發操作

## 🛠️ 開發工具

```bash
# 程式碼格式化
go fmt ./...

# 程式碼檢查
go vet ./...

# 安全性掃描
go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
gosec ./...

# 相依套件更新
go mod tidy
go mod verify
```

## 📝 開發狀態

✅ **已完成功能**:

- 核心筆記編輯和預覽
- 檔案系統管理
- 加密和安全功能
- 自動保存系統
- 效能監控和優化
- 錯誤處理和通知
- 完整的測試套件

🚧 **進行中**:

- 應用程式打包和部署優化
- 使用者介面最終調整

## 🤝 貢獻指南

1. Fork 專案
2. 建立功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交變更 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 開啟 Pull Request

## 📄 授權條款

本專案採用 MIT 授權條款 - 詳見 [LICENSE](LICENSE) 檔案。

## 📞 支援和回饋

如有問題或建議，請：

- 開啟 [GitHub Issue](https://github.com/your-username/mac-notebook-app/issues)
- 發送郵件至：support@yourcompany.com
- 查看 [CHANGELOG.md](CHANGELOG.md) 了解最新變更

---

**Mac Notebook v0.14.0** - 一個安全、高效的 Markdown 筆記應用程式 🚀
