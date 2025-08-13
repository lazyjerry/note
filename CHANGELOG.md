# 變更日誌 (Changelog)

本檔案記錄 Mac Notebook App 專案的所有重要變更和開發進度。

格式基於 [Keep a Changelog](https://keepachangelog.com/zh-TW/1.0.0/)，
並且本專案遵循 [語義化版本](https://semver.org/lang/zh-TW/)。

## [未發布] - 開發中

### 🚧 進行中

- 準備開始 Task 4: 實作加密功能

### 📋 待辦事項

- [ ] 4. 實作加密功能
- [ ] 5. 建立編輯器服務
- [ ] 6. 實作自動保存系統
- [ ] 7. 建立 Fyne UI 基礎架構
- [ ] 8. 實作編輯器 UI 元件
- [ ] 9. 實作檔案操作 UI
- [ ] 10. 建立加密 UI 元件
- [ ] 11. 實作設定管理 UI
- [ ] 12. 整合所有服務到 UI
- [ ] 13. 實作錯誤處理和用戶回饋
- [ ] 14. 效能優化和測試
- [ ] 15. 應用程式打包和部署

---

## [0.3.0] - 2024-01-XX - 檔案系統操作完成

### ✅ 新增功能

- **Task 3.1: 實作 FileRepository**

  - 建立 `LocalFileRepository` 實作檔案儲存庫介面
  - 實作基礎檔案操作：`ReadFile`, `WriteFile`, `FileExists`, `DeleteFile`
  - 實作目錄操作：`CreateDirectory`, `ListDirectory`, `WalkDirectory`
  - 實作 Markdown 檔案特殊處理：`ReadMarkdownFile`, `WriteMarkdownFile`, `IsMarkdownFile`
  - 新增檔案路徑安全性驗證，防止目錄遍歷攻擊
  - 新增完整的錯誤處理機制

- **Task 3.2: 實作 FileManagerService**
  - 建立 `LocalFileManagerService` 實作檔案管理服務介面
  - 實作檔案和資料夾 CRUD 操作：`ListFiles`, `CreateDirectory`, `DeleteFile`, `RenameFile`, `MoveFile`
  - 實作檔案樹狀結構遍歷：`GetFileTree` 與 `FileTreeNode` 結構
  - 實作檔案搜尋功能：`SearchFiles` 支援模式匹配和遞迴搜尋
  - 實作檔案複製功能：`CopyFile` 支援檔案和目錄複製
  - 實作目錄大小計算：`GetDirectorySize` 遞迴計算目錄總大小
  - 新增檔案排序功能（目錄優先，按名稱排序）

### 🧪 測試改進

- 新增 `file_repository_test.go` 包含 8 個測試函數
  - 測試建構函數、檔案操作、目錄操作、Markdown 操作
  - 測試路徑驗證、目錄遍歷、並發存取
  - 包含效能測試 (Benchmark)
- 新增 `file_manager_service_test.go` 包含 10 個測試函數
  - 測試所有 CRUD 操作、搜尋、複製、移動功能
  - 測試檔案樹生成、目錄大小計算、路徑驗證
  - 涵蓋正常情況、邊界條件和錯誤處理

### 📝 文件更新

- 所有新增程式碼都包含詳細的繁體中文註解
- 函數說明包含參數、回傳值和執行流程描述
- 複雜邏輯都有步驟化的執行流程說明

### 🔧 技術改進

- 實作模組化設計，清楚分離儲存庫層和服務層
- 使用依賴注入提高程式碼可測試性
- 統一使用 `AppError` 進行結構化錯誤處理
- 實作安全的檔案路徑處理，防止安全漏洞

---

## [0.2.0] - 2024-01-XX - 資料模型和驗證完成

### ✅ 新增功能

- **Task 2.1: 創建 Note 資料模型**

  - 建立 `Note` 結構體，包含 ID、標題、內容、檔案路徑等欄位
  - 實作筆記建立函數 `NewNote` 與唯一 ID 生成
  - 實作筆記內容更新：`UpdateContent` 方法
  - 實作筆記保存標記：`MarkSaved` 方法
  - 實作筆記修改狀態檢查：`IsModified` 方法
  - 實作筆記資料驗證：`Validate` 方法
  - 實作加密設定管理：`SetEncryption`, `RemoveEncryption` 方法
  - 實作實用功能：`GetWordCount`, `Clone` 方法

- **Task 2.2: 實作 Settings 資料模型**
  - 建立 `Settings` 結構體，包含加密、自動保存、主題等設定
  - 實作預設設定建立：`NewDefaultSettings` 函數
  - 實作設定驗證：`Validate` 方法
  - 實作設定更新方法：`UpdateEncryption`, `UpdateAutoSaveInterval`, `UpdateTheme` 等
  - 實作生物識別設定：`ToggleBiometric`, `SetBiometric` 方法
  - 實作設定複製和比較：`Clone`, `IsDefault` 方法
  - 實作檔案保存和載入：`SaveToFile`, `LoadFromFile` 方法
  - 實作支援清單取得：`GetSupportedEncryptionAlgorithms`, `GetSupportedThemes`

### 🧪 測試改進

- 新增 `note_test.go` 包含 11 個測試函數
  - 涵蓋筆記建立、內容更新、保存標記、修改狀態檢查
  - 測試資料驗證、加密設定、字數統計、複製功能
  - 測試 ID 生成和隨機字串生成功能
- 新增 `settings_test.go` 包含 14 個測試函數
  - 涵蓋預設設定建立、設定驗證、各種更新方法
  - 測試生物識別切換、設定複製、預設狀態檢查
  - 測試檔案保存載入、錯誤處理、支援清單功能

### 📝 文件更新

- 所有資料模型都包含完整的繁體中文註解
- 每個方法都有詳細的參數說明和執行流程描述
- 測試函數都有清楚的測試目的說明

### 🔧 技術改進

- 實作結構化的錯誤處理系統 (`AppError`, `ValidationError`)
- 使用時間戳追蹤筆記的建立、修改和保存狀態
- 實作設定的檔案持久化機制
- 加入資料驗證確保資料完整性

---

## [0.1.0] - 2024-01-XX - 專案基礎架構完成

### ✅ 新增功能

- **Task 1: 建立專案結構和核心介面**

  - 初始化 Go 模組 (`go.mod`) 和依賴管理
  - 建立專案目錄結構：
    - `internal/models/` - 資料模型
    - `internal/services/` - 業務邏輯服務
    - `internal/repositories/` - 資料存取層
    - `ui/` - 使用者介面
  - 定義核心介面：
    - `FileRepository` - 檔案操作介面
    - `SettingsRepository` - 設定持久化介面
    - `EncryptionRepository` - 加密金鑰管理介面
    - `EditorService` - 筆記編輯服務介面
    - `FileManagerService` - 檔案管理服務介面
    - `EncryptionService` - 加密服務介面
    - `AutoSaveService` - 自動保存服務介面
    - `SettingsService` - 設定管理服務介面

- **Task 1.5: 為所有程式碼添加繁體中文註解**
  - 所有介面和結構體都有完整的繁體中文說明
  - 每個方法都包含參數、回傳值和用途的詳細描述
  - 註解格式符合 Go 語言文件慣例

### 🔧 技術設定

- 設定 Go 1.21 作為最低版本要求
- 新增核心依賴：
  - `fyne.io/fyne/v2 v2.4.0` - GUI 框架
  - `github.com/yuin/goldmark v1.6.0` - Markdown 解析
  - `github.com/fsnotify/fsnotify v1.7.0` - 檔案系統監控
  - `golang.org/x/crypto v0.17.0` - 加密功能
- 建立基本的 `main.go` 檔案與 Fyne 應用程式架構

### 📝 文件建立

- 建立 `README.md` 專案說明文件
- 建立 `.gitignore` 檔案
- 設定專案授權 (`LICENSE`)

---

## 專案統計

### 📊 程式碼統計 (截至 v0.3.0)

- **總檔案數**: 15+ 個 Go 檔案
- **程式碼行數**: 2000+ 行 (包含註解)
- **測試檔案**: 4 個
- **測試函數**: 35+ 個
- **測試覆蓋率**: 高 (包含正常情況、邊界條件、錯誤處理)

### 🏗️ 架構完成度

- ✅ **資料層 (Models)**: 100% 完成
- ✅ **儲存庫層 (Repositories)**: 檔案操作 100% 完成
- ✅ **服務層 (Services)**: 檔案管理 100% 完成
- 🚧 **服務層 (Services)**: 加密服務 0% 完成
- 🚧 **使用者介面層 (UI)**: 0% 完成

### 🎯 里程碑達成

- [x] **里程碑 1**: 專案基礎架構 (v0.1.0)
- [x] **里程碑 2**: 核心資料模型 (v0.2.0)
- [x] **里程碑 3**: 檔案系統操作 (v0.3.0)
- [ ] **里程碑 4**: 加密和安全功能
- [ ] **里程碑 5**: 編輯器核心功能
- [ ] **里程碑 6**: 使用者介面
- [ ] **里程碑 7**: 完整應用程式

---

## 貢獻指南

### 🔄 版本控制

- 使用 [語義化版本](https://semver.org/lang/zh-TW/) 進行版本管理
- 主要版本 (Major): 不相容的 API 變更
- 次要版本 (Minor): 向下相容的功能新增
- 修訂版本 (Patch): 向下相容的問題修正

### 📝 變更日誌格式

- **新增功能** (Added): 新功能
- **變更** (Changed): 現有功能的變更
- **棄用** (Deprecated): 即將移除的功能
- **移除** (Removed): 已移除的功能
- **修正** (Fixed): 錯誤修正
- **安全性** (Security): 安全性相關變更

### 🧪 品質標準

- 所有新功能都必須包含單元測試
- 程式碼覆蓋率應保持在高水準
- 所有程式碼都必須包含繁體中文註解
- 遵循 Go 語言最佳實踐和慣例
