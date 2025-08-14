# 變更日誌 (Changelog)

本檔案記錄 Mac Notebook App 專案的所有重要變更和開發進度。

格式基於 [Keep a Changelog](https://keepachangelog.com/zh-TW/1.0.0/)，
並且本專案遵循 [語義化版本](https://semver.org/lang/zh-TW/)。

## [未發布] - 開發中

### 🚧 進行中

- 準備開始 Task 8: 實作編輯器 UI 元件

### ✅ 新增功能

- **Task 7.1: 創建主視窗結構**

  - 建立完整的 `MainWindow` 結構體，包含所有主要 UI 元件
  - 實作主視窗建立功能：`NewMainWindow` 建立視窗並設定基本屬性
  - 實作選單欄建立：`createMenuBar` 包含檔案、編輯、檢視、說明選單
  - 實作工具欄建立：`createToolBar` 包含常用功能的快速存取按鈕
  - 實作狀態欄建立：`createStatusBar` 顯示保存狀態、加密狀態和字數統計
  - 實作主要內容區域：`createContentArea` 建立左右分割的面板佈局
  - 實作完整佈局組合：`assembleMainLayout` 組合所有元件到主視窗
  - 實作狀態更新功能：`UpdateSaveStatus`, `UpdateEncryptionStatus`, `UpdateWordCount`
  - 實作視窗實例存取：`GetWindow` 提供視窗實例給其他元件使用
  - 設定視窗基本屬性：標題、大小 (1200x800)、居中顯示、關閉攔截
  - 實作左右面板分割比例：左側 30%，右側 70%
  - 為所有功能預留 TODO 標記，準備後續任務實作

- **Task 7.2: 實作檔案樹狀視圖**

  - 建立完整的 `FileTreeWidget` 元件，支援檔案和目錄的樹狀顯示
  - 實作檔案樹節點結構：`FileNode` 包含路徑、名稱、類型和展開狀態
  - 實作檔案結構載入：`loadFileStructure`, `loadDirectoryChildren` 遞迴載入目錄內容
  - 實作樹狀元件回調：`getChildUIDs`, `isBranch`, `createNodeWidget`, `updateNodeWidget`
  - 實作節點選擇處理：`handleNodeSelection` 支援檔案和目錄選擇回調
  - 實作分支操作：`handleBranchOpened`, `handleBranchClosed` 管理目錄展開狀態
  - 實作檔案樹刷新：`Refresh` 重新載入檔案結構並更新顯示
  - 實作路徑展開：`ExpandPath` 自動展開指定路徑的所有父目錄
  - 實作回調函數設定：`SetOnFileSelect`, `SetOnFileOpen`, `SetOnDirectoryOpen`
  - 實作檔案類型識別：根據副檔名顯示不同圖示（資料夾、Markdown、文字檔案）
  - 整合到主視窗左側面板：建立檔案瀏覽器區域包含標題和控制按鈕
  - 建立示例檔案樹：展示基本的樹狀結構和節點選擇功能

- **Task 6.1: 建立自動保存服務**

  - 建立 `AutoSaveService` 介面的完整實作 (`AutoSaveServiceImpl`)
  - 實作定時保存邏輯：`StartAutoSave` 啟動自動保存定時器
  - 實作保存狀態追蹤：`SaveStatus` 結構體記錄保存狀態和統計
  - 實作立即保存功能：`SaveNow` 支援手動觸發保存
  - 實作自動保存管理：`StopAutoSave` 停止定時器並清理資源
  - 實作狀態查詢功能：`GetSaveStatus`, `GetAllSaveStatuses` 取得保存狀態
  - 實作服務關閉功能：`Shutdown` 安全關閉所有定時器
  - 實作執行緒安全的並發存取保護：使用 `sync.RWMutex`
  - 實作智慧保存邏輯：只保存已修改的筆記，跳過未修改筆記
  - 實作定時器重新排程：`rescheduleTimer` 自動重設下次保存時間
  - 實作保存錯誤處理和狀態更新
  - 實作筆記快取管理：避免重複載入筆記實例

- **Task 6.2: 整合加密檔案自動保存**

  - 增強 `AutoSaveService` 支援加密檔案的背景保存
  - 實作可配置的保存間隔：`StartAutoSaveWithSettings` 使用設定服務的間隔
  - 實作加密檔案特殊處理：`getAutoSaveInterval` 為加密檔案增加額外延遲
  - 實作保存失敗的重試機制：`saveNoteWithRetry` 對加密檔案進行最多 3 次重試
  - 實作動態間隔更新：`UpdateAutoSaveInterval` 支援運行時調整保存頻率
  - 實作加密檔案統計：`GetEncryptedFileCount` 追蹤加密檔案數量
  - 實作延遲配置：`SetEncryptedBackoff` 動態調整加密檔案的額外延遲
  - 實作設定服務整合：支援從 `SettingsService` 載入使用者配置的間隔
  - 實作錯誤分類：區分一般保存錯誤和加密檔案特定錯誤
  - 實作向後相容：`NewAutoSaveServiceWithDefaults` 支援無設定服務的使用
  - 實作智慧延遲：加密檔案自動增加 30 秒延遲以減少加密操作頻率

- **Task 5.1: 實作 Markdown 編輯器核心**
  - 建立 `EditorService` 介面的完整實作 (`editorService`)
  - 整合 goldmark Markdown 解析器，支援 GitHub Flavored Markdown (GFM)
  - 實作筆記建立功能：`CreateNote` 生成唯一 ID 和時間戳
  - 實作筆記開啟功能：`OpenNote` 支援一般和加密檔案識別
  - 實作筆記保存功能：`SaveNote` 自動生成檔案路徑和副檔名
  - 實作內容更新功能：`UpdateContent` 即時更新筆記內容
  - 實作 Markdown 預覽功能：`PreviewMarkdown` 轉換為 HTML
  - 實作活躍筆記管理：`GetActiveNote`, `CloseNote`, `GetActiveNotes`
  - 實作檔案名稱清理：`sanitizeFileName` 移除不合法字元
  - 支援表格、刪除線、任務列表等 Markdown 擴展功能
  - 支援自動標題 ID 生成和 XHTML 相容性

### 🧪 測試改進

- 新增 `main_window_test.go` 包含 5 個測試函數

  - 測試主視窗建立和初始化：`TestNewMainWindow`
  - 測試 UI 元件初始化：`TestMainWindowUIComponents` 驗證所有元件正確建立
  - 測試狀態更新功能：`TestMainWindowStatusUpdates` 測試保存、加密、字數狀態更新
  - 測試視窗實例存取：`TestMainWindowGetWindow` 驗證視窗實例回傳
  - 測試分割比例設定：`TestMainWindowSplitRatio` 驗證左右面板分割比例
  - 使用 Fyne 測試框架進行 UI 元件測試
  - 涵蓋正常情況、元件驗證和狀態管理測試

- 新增 `file_tree_test.go` 包含 9 個測試函數

  - 測試檔案樹元件建立：`TestNewFileTreeWidget` 驗證元件正確初始化
  - 測試檔案結構載入：`TestFileTreeLoadFileStructure` 驗證檔案和目錄正確載入
  - 測試子節點 ID 取得：`TestFileTreeGetChildUIDs` 驗證樹狀結構正確建立
  - 測試分支識別：`TestFileTreeIsBranch` 驗證目錄和檔案正確區分
  - 測試節點選擇：`TestFileTreeNodeSelection` 驗證選擇回調正確觸發
  - 測試分支操作：`TestFileTreeBranchOperations` 驗證目錄展開和收合狀態
  - 測試檔案樹刷新：`TestFileTreeRefresh` 驗證重新載入功能
  - 測試路徑展開：`TestFileTreeExpandPath` 驗證父目錄自動展開
  - 測試真實檔案系統整合：`TestFileTreeWithRealFileSystem` 使用臨時目錄測試
  - 實作完整的模擬檔案管理服務：`mockFileManagerService` 支援測試檔案結構
  - 涵蓋正常情況、錯誤處理、邊界條件和真實檔案系統整合測試

- 新增 `auto_save_service_test.go` 包含 15+ 個測試函數

  - 測試自動保存服務建立和初始化：`TestNewAutoSaveService`
  - 測試自動保存啟動和停止：`TestStartAutoSave`, `TestStopAutoSave`
  - 測試立即保存功能：`TestSaveNow`, `TestSaveNowWithNonExistentNote`
  - 測試並發保存防護：`TestSaveNowWithSaveInProgress`
  - 測試自動保存觸發：`TestAutoSaveWithModifiedNote`, `TestAutoSaveWithUnmodifiedNote`
  - 測試狀態查詢功能：`TestGetSaveStatus`, `TestGetAllSaveStatuses`
  - 測試服務關閉：`TestShutdown`
  - 測試錯誤處理：`TestSaveErrorHandling`
  - 測試並發自動保存：`TestConcurrentAutoSave`
  - 測試定時器重新排程：`TestRescheduleTimer`
  - 實作完整的模擬編輯器服務：`MockEditorService` 支援保存延遲和錯誤模擬
  - 包含執行緒安全性測試和效能測試
  - 新增加密檔案專用測試：`TestEncryptedFileAutoSave`, `TestEncryptedFileRetryMechanism`
  - 新增設定服務整合測試：`TestStartAutoSaveWithSettings`, `TestAutoSaveWithSettingsLoadError`
  - 新增動態配置測試：`TestUpdateAutoSaveInterval`, `TestSetEncryptedBackoff`
  - 新增統計功能測試：`TestGetEncryptedFileCount`
  - 實作完整的模擬設定服務：`MockSettingsService` 支援設定載入和錯誤模擬
  - 包含加密檔案重試機制的詳細測試和錯誤處理驗證

- 新增 `editor_service_test.go` 包含 15+ 個測試函數
  - 測試編輯器服務建立和初始化
  - 測試筆記建立、開啟、保存、更新功能
  - 測試加密檔案識別和處理
  - 測試 Markdown 預覽的各種語法轉換
  - 測試表格、任務列表等擴展功能
  - 測試檔案名稱清理和安全性
  - 測試活躍筆記管理功能
  - 測試錯誤處理：不存在檔案、空筆記、無效輸入
  - 包含模擬檔案系統 (`mockFileRepository`) 用於隔離測試

### 📝 文件更新

- 所有主視窗 UI 程式碼都包含詳細的繁體中文註解
- 每個 UI 建立函數都有完整的參數、回傳值和執行流程說明
- 選單欄、工具欄、狀態欄建立流程包含詳細的元件說明
- 佈局組合邏輯包含清楚的結構和比例設定說明
- 測試程式碼包含清楚的測試目的和驗證邏輯說明
- 更新 `main.go` 使用新的 MainWindow 結構，移除舊的佈局程式碼

- 所有檔案樹元件程式碼都包含詳細的繁體中文註解
- 每個檔案樹方法都有完整的參數、回傳值和執行流程說明
- 樹狀元件回調函數包含詳細的實作邏輯說明
- 檔案載入和節點管理流程包含清楚的步驟說明
- 測試程式碼包含清楚的測試目的和驗證邏輯說明
- 整合主視窗左側面板，建立完整的檔案瀏覽器介面

### 🔧 技術改進

- 實作標準桌面應用程式佈局架構：選單欄、工具欄、內容區域、狀態欄
- 使用 Fyne 框架建立跨平台 GUI 應用程式
- 實作模組化 UI 元件設計，每個元件獨立建立和管理
- 實作響應式佈局：水平分割容器支援動態調整面板大小
- 實作狀態管理系統：即時更新保存、加密和字數統計狀態
- 實作完整的選單系統：檔案、編輯、檢視、說明選單結構
- 實作工具欄快速存取：常用功能的圖示按鈕
- 實作視窗生命週期管理：關閉攔截和清理機制
- 支援中文介面：所有選單和標籤都使用繁體中文
- 實作 UI 測試架構：使用 Fyne 測試框架進行元件測試

- 實作自訂 Fyne 元件：FileTreeWidget 繼承 BaseWidget 建立可重用元件
- 實作樹狀資料結構：FileNode 支援父子關係和狀態管理
- 實作動態檔案載入：按需載入目錄內容，提高大型檔案系統的效能
- 實作檔案節點快取：使用 map 快速查找和存取檔案節點
- 實作回調函數架構：支援檔案選擇、開啟和目錄操作的事件處理
- 實作檔案類型識別：根據副檔名和屬性顯示適當的圖示
- 實作路徑管理：支援絕對路徑和相對路徑的正確處理
- 實作樹狀元件整合：完整實作 Fyne Tree 元件的所有回調介面
- 支援檔案系統監控：為未來的即時更新功能預留架構
- 實作模組化設計：檔案樹元件可獨立使用和測試

### 📋 待辦事項

- [ ] 8. 實作編輯器 UI 元件
- [ ] 9. 實作檔案操作 UI
- [ ] 10. 建立加密 UI 元件
- [ ] 11. 實作設定管理 UI
- [ ] 12. 整合所有服務到 UI
- [ ] 13. 實作錯誤處理和用戶回饋
- [ ] 14. 效能優化和測試
- [ ] 15. 應用程式打包和部署

---

## [0.7.0] - 2024-01-XX - Fyne UI 基礎架構完成

### ✅ 新增功能

- **Task 7.1: 創建主視窗結構**

  - 建立完整的 `MainWindow` 結構體，包含所有主要 UI 元件
  - 實作主視窗建立功能：`NewMainWindow` 建立視窗並設定基本屬性
  - 實作選單欄建立：`createMenuBar` 包含檔案、編輯、檢視、說明選單
  - 實作工具欄建立：`createToolBar` 包含常用功能的快速存取按鈕
  - 實作狀態欄建立：`createStatusBar` 顯示保存狀態、加密狀態和字數統計
  - 實作主要內容區域：`createContentArea` 建立左右分割的面板佈局
  - 實作完整佈局組合：`assembleMainLayout` 組合所有元件到主視窗
  - 實作狀態更新功能：`UpdateSaveStatus`, `UpdateEncryptionStatus`, `UpdateWordCount`
  - 實作視窗實例存取：`GetWindow` 提供視窗實例給其他元件使用
  - 設定視窗基本屬性：標題、大小 (1200x800)、居中顯示、關閉攔截
  - 實作左右面板分割比例：左側 30%，右側 70%
  - 為所有功能預留 TODO 標記，準備後續任務實作

- **Task 7.2: 實作檔案樹狀視圖**

  - 建立完整的 `FileTreeWidget` 元件，支援檔案和目錄的樹狀顯示
  - 實作檔案樹節點結構：`FileNode` 包含路徑、名稱、類型和展開狀態
  - 實作檔案結構載入：`loadFileStructure`, `loadDirectoryChildren` 遞迴載入目錄內容
  - 實作樹狀元件回調：`getChildUIDs`, `isBranch`, `createNodeWidget`, `updateNodeWidget`
  - 實作節點選擇處理：`handleNodeSelection` 支援檔案和目錄選擇回調
  - 實作分支操作：`handleBranchOpened`, `handleBranchClosed` 管理目錄展開狀態
  - 實作檔案樹刷新：`Refresh` 重新載入檔案結構並更新顯示
  - 實作路徑展開：`ExpandPath` 自動展開指定路徑的所有父目錄
  - 實作回調函數設定：`SetOnFileSelect`, `SetOnFileOpen`, `SetOnDirectoryOpen`
  - 實作檔案類型識別：根據副檔名顯示不同圖示（資料夾、Markdown、文字檔案）
  - 整合到主視窗左側面板：建立檔案瀏覽器區域包含標題和控制按鈕
  - 建立示例檔案樹：展示基本的樹狀結構和節點選擇功能

### 🧪 測試改進

- 新增 `main_window_test.go` 包含 5 個測試函數

  - 測試主視窗建立和初始化：`TestNewMainWindow`
  - 測試 UI 元件初始化：`TestMainWindowUIComponents` 驗證所有元件正確建立
  - 測試狀態更新功能：`TestMainWindowStatusUpdates` 測試保存、加密、字數狀態更新
  - 測試視窗實例存取：`TestMainWindowGetWindow` 驗證視窗實例回傳
  - 測試分割比例設定：`TestMainWindowSplitRatio` 驗證左右面板分割比例
  - 使用 Fyne 測試框架進行 UI 元件測試
  - 涵蓋正常情況、元件驗證和狀態管理測試

- 新增 `file_tree_test.go` 包含 9 個測試函數

  - 測試檔案樹元件建立：`TestNewFileTreeWidget` 驗證元件正確初始化
  - 測試檔案結構載入：`TestFileTreeLoadFileStructure` 驗證檔案和目錄正確載入
  - 測試子節點 ID 取得：`TestFileTreeGetChildUIDs` 驗證樹狀結構正確建立
  - 測試分支識別：`TestFileTreeIsBranch` 驗證目錄和檔案正確區分
  - 測試節點選擇：`TestFileTreeNodeSelection` 驗證選擇回調正確觸發
  - 測試分支操作：`TestFileTreeBranchOperations` 驗證目錄展開和收合狀態
  - 測試檔案樹刷新：`TestFileTreeRefresh` 驗證重新載入功能
  - 測試路徑展開：`TestFileTreeExpandPath` 驗證父目錄自動展開
  - 測試真實檔案系統整合：`TestFileTreeWithRealFileSystem` 使用臨時目錄測試
  - 實作完整的模擬檔案管理服務：`mockFileManagerService` 支援測試檔案結構
  - 涵蓋正常情況、錯誤處理、邊界條件和真實檔案系統整合測試

### 📝 文件更新

- 所有主視窗 UI 程式碼都包含詳細的繁體中文註解
- 每個 UI 建立函數都有完整的參數、回傳值和執行流程說明
- 選單欄、工具欄、狀態欄建立流程包含詳細的元件說明
- 佈局組合邏輯包含清楚的結構和比例設定說明
- 測試程式碼包含清楚的測試目的和驗證邏輯說明
- 更新 `main.go` 使用新的 MainWindow 結構，移除舊的佈局程式碼

- 所有檔案樹元件程式碼都包含詳細的繁體中文註解
- 每個檔案樹方法都有完整的參數、回傳值和執行流程說明
- 樹狀元件回調函數包含詳細的實作邏輯說明
- 檔案載入和節點管理流程包含清楚的步驟說明
- 測試程式碼包含清楚的測試目的和驗證邏輯說明
- 整合主視窗左側面板，建立完整的檔案瀏覽器介面

### 🔧 技術改進

- 實作標準桌面應用程式佈局架構：選單欄、工具欄、內容區域、狀態欄
- 使用 Fyne 框架建立跨平台 GUI 應用程式
- 實作模組化 UI 元件設計，每個元件獨立建立和管理
- 實作響應式佈局：水平分割容器支援動態調整面板大小
- 實作狀態管理系統：即時更新保存、加密和字數統計狀態
- 實作完整的選單系統：檔案、編輯、檢視、說明選單結構
- 實作工具欄快速存取：常用功能的圖示按鈕
- 實作視窗生命週期管理：關閉攔截和清理機制
- 支援中文介面：所有選單和標籤都使用繁體中文
- 實作 UI 測試架構：使用 Fyne 測試框架進行元件測試

- 實作自訂 Fyne 元件：FileTreeWidget 繼承 BaseWidget 建立可重用元件
- 實作樹狀資料結構：FileNode 支援父子關係和狀態管理
- 實作動態檔案載入：按需載入目錄內容，提高大型檔案系統的效能
- 實作檔案節點快取：使用 map 快速查找和存取檔案節點
- 實作回調函數架構：支援檔案選擇、開啟和目錄操作的事件處理
- 實作檔案類型識別：根據副檔名和屬性顯示適當的圖示
- 實作路徑管理：支援絕對路徑和相對路徑的正確處理
- 實作樹狀元件整合：完整實作 Fyne Tree 元件的所有回調介面
- 支援檔案系統監控：為未來的即時更新功能預留架構
- 實作模組化設計：檔案樹元件可獨立使用和測試

### 📋 待辦事項

- [ ] 8. 實作編輯器 UI 元件
- [ ] 9. 實作檔案操作 UI
- [ ] 10. 建立加密 UI 元件
- [ ] 11. 實作設定管理 UI
- [ ] 12. 整合所有服務到 UI
- [ ] 13. 實作錯誤處理和用戶回饋
- [ ] 14. 效能優化和測試
- [ ] 15. 應用程式打包和部署

---

## [0.6.0] - 2024-01-XX - 自動保存系統完成

### ✅ 新增功能

- 所有編輯器服務程式碼都包含詳細的繁體中文註解
- 每個編輯器函數都有完整的參數、回傳值和執行流程說明
- Markdown 解析器配置包含詳細的功能說明
- 測試程式碼包含清楚的測試目的和驗證邏輯說明

- **Task 5.2: 整合加密功能到編輯器**
  - 增強 `EditorService` 整合加密、密碼和生物識別服務
  - 實作加密檔案的開啟功能：`OpenNote` 支援自動解密流程
  - 實作加密檔案的保存功能：`SaveNote` 支援自動加密流程
  - 實作加密狀態管理：`EnableEncryption`, `DisableEncryption`
  - 實作密碼驗證整合：`DecryptWithPassword`, `EncryptWithPassword`
  - 實作生物識別驗證整合：支援 Touch ID/Face ID 解密
  - 實作加密狀態查詢：`IsEncrypted`, `GetEncryptionType`
  - 實作安全的加密檔案處理：`decryptFileContent`, `encryptFileContent`
  - 支援多種加密演算法：AES-256 和 ChaCha20-Poly1305
  - 實作加密檔案的自動識別（.enc 副檔名）

### 🧪 測試改進

- 更新 `editor_service_test.go` 包含 20+ 個測試函數
  - 新增加密功能測試：`TestEnableEncryption`, `TestDisableEncryption`
  - 新增加密狀態查詢測試：`TestIsEncrypted`, `TestGetEncryptionType`
  - 新增加密檔案開啟測試：`TestOpenEncryptedNote`
  - 更新所有現有測試以支援新的服務架構
  - 實作完整的模擬服務：`mockEncryptionService`, `mockPasswordService`, `mockBiometricService`
  - 測試加密檔案的錯誤處理：密碼驗證要求、解密失敗等
  - 測試加密狀態的正確管理和更新

### 📝 文件更新

- 所有加密整合程式碼都包含詳細的繁體中文註解
- 每個加密相關函數都有完整的參數、回傳值和執行流程說明
- 加密檔案處理流程包含詳細的安全性考量說明
- 生物識別整合包含回退機制的詳細說明

### 🔧 技術改進

- 使用 goldmark 作為高效能 Markdown 解析器
- 支援 GitHub Flavored Markdown 標準
- 實作活躍筆記快取機制提高效能
- 實作安全的檔案名稱處理防止檔案系統錯誤
- 使用 UUID 生成唯一筆記識別碼
- 實作完整的錯誤處理和狀態管理
- 整合多層級安全驗證：密碼 + 生物識別
- 實作加密檔案的透明處理：自動加密/解密
- 支援加密狀態的動態切換
- 實作安全的憑證管理架構（為未來 Keychain 整合做準備）

### 📋 待辦事項

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

## [0.6.0] - 2024-01-XX - 自動保存系統完成

### ✅ 新增功能

- **Task 6.1: 建立自動保存服務**

  - 建立 `AutoSaveService` 介面的完整實作 (`AutoSaveServiceImpl`)
  - 實作定時保存邏輯：`StartAutoSave` 啟動自動保存定時器
  - 實作保存狀態追蹤：`SaveStatus` 結構體記錄保存狀態和統計
  - 實作立即保存功能：`SaveNow` 支援手動觸發保存
  - 實作自動保存管理：`StopAutoSave` 停止定時器並清理資源
  - 實作狀態查詢功能：`GetSaveStatus`, `GetAllSaveStatuses` 取得保存狀態
  - 實作服務關閉功能：`Shutdown` 安全關閉所有定時器
  - 實作執行緒安全的並發存取保護：使用 `sync.RWMutex`
  - 實作智慧保存邏輯：只保存已修改的筆記，跳過未修改筆記
  - 實作定時器重新排程：`rescheduleTimer` 自動重設下次保存時間
  - 實作保存錯誤處理和狀態更新
  - 實作筆記快取管理：避免重複載入筆記實例

- **Task 6.2: 整合加密檔案自動保存**
  - 增強 `AutoSaveService` 支援加密檔案的背景保存
  - 實作可配置的保存間隔：`StartAutoSaveWithSettings` 使用設定服務的間隔
  - 實作加密檔案特殊處理：`getAutoSaveInterval` 為加密檔案增加額外延遲
  - 實作保存失敗的重試機制：`saveNoteWithRetry` 對加密檔案進行最多 3 次重試
  - 實作動態間隔更新：`UpdateAutoSaveInterval` 支援運行時調整保存頻率
  - 實作加密檔案統計：`GetEncryptedFileCount` 追蹤加密檔案數量
  - 實作延遲配置：`SetEncryptedBackoff` 動態調整加密檔案的額外延遲
  - 實作設定服務整合：支援從 `SettingsService` 載入使用者配置的間隔
  - 實作錯誤分類：區分一般保存錯誤和加密檔案特定錯誤
  - 實作向後相容：`NewAutoSaveServiceWithDefaults` 支援無設定服務的使用
  - 實作智慧延遲：加密檔案自動增加 30 秒延遲以減少加密操作頻率

### 🧪 測試改進

- 新增 `auto_save_service_test.go` 包含 20+ 個測試函數
  - 測試自動保存服務建立和初始化：`TestNewAutoSaveService`
  - 測試自動保存啟動和停止：`TestStartAutoSave`, `TestStopAutoSave`
  - 測試立即保存功能：`TestSaveNow`, `TestSaveNowWithNonExistentNote`
  - 測試並發保存防護：`TestSaveNowWithSaveInProgress`
  - 測試自動保存觸發：`TestAutoSaveWithModifiedNote`, `TestAutoSaveWithUnmodifiedNote`
  - 測試狀態查詢功能：`TestGetSaveStatus`, `TestGetAllSaveStatuses`
  - 測試服務關閉：`TestShutdown`
  - 測試錯誤處理：`TestSaveErrorHandling`
  - 測試並發自動保存：`TestConcurrentAutoSave`
  - 測試定時器重新排程：`TestRescheduleTimer`
  - 新增加密檔案專用測試：`TestEncryptedFileAutoSave`, `TestEncryptedFileRetryMechanism`
  - 新增設定服務整合測試：`TestStartAutoSaveWithSettings`, `TestAutoSaveWithSettingsLoadError`
  - 新增動態配置測試：`TestUpdateAutoSaveInterval`, `TestSetEncryptedBackoff`
  - 新增統計功能測試：`TestGetEncryptedFileCount`
  - 實作完整的模擬編輯器服務：`MockEditorService` 支援保存延遲和錯誤模擬
  - 實作完整的模擬設定服務：`MockSettingsService` 支援設定載入和錯誤模擬
  - 包含執行緒安全性測試、效能測試和加密檔案重試機制的詳細測試

### 📝 文件更新

- 所有自動保存服務程式碼都包含詳細的繁體中文註解
- 每個自動保存函數都有完整的參數、回傳值和執行流程說明
- 加密檔案處理邏輯包含詳細的重試機制和錯誤處理說明
- 設定服務整合包含回退機制的詳細說明

### 🔧 技術改進

- 實作高效能的定時器管理系統
- 使用執行緒安全的並發存取保護
- 實作智慧保存邏輯，避免不必要的保存操作
- 支援動態配置和運行時調整
- 實作完整的錯誤處理和重試機制
- 整合設定服務，支援使用者自訂保存間隔
- 針對加密檔案優化保存頻率，減少加密操作負載
- 實作資源管理和服務生命週期控制

### 📋 待辦事項

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

## [0.5.0] - 2024-01-XX - 編輯器服務完成

### ✅ 新增功能

- **Task 5.1: 實作 Markdown 編輯器核心**

  - 建立 `EditorService` 介面的完整實作 (`editorService`)
  - 整合 goldmark Markdown 解析器，支援 GitHub Flavored Markdown (GFM)
  - 實作筆記建立功能：`CreateNote` 生成唯一 ID 和時間戳
  - 實作筆記開啟功能：`OpenNote` 支援一般和加密檔案識別
  - 實作筆記保存功能：`SaveNote` 自動生成檔案路徑和副檔名
  - 實作內容更新功能：`UpdateContent` 即時更新筆記內容
  - 實作 Markdown 預覽功能：`PreviewMarkdown` 轉換為 HTML
  - 實作活躍筆記管理：`GetActiveNote`, `CloseNote`, `GetActiveNotes`
  - 實作檔案名稱清理：`sanitizeFileName` 移除不合法字元
  - 支援表格、刪除線、任務列表等 Markdown 擴展功能
  - 支援自動標題 ID 生成和 XHTML 相容性

- **Task 5.2: 整合加密功能到編輯器**
  - 增強 `EditorService` 整合加密、密碼和生物識別服務
  - 實作加密檔案的開啟功能：`OpenNote` 支援自動解密流程
  - 實作加密檔案的保存功能：`SaveNote` 支援自動加密流程
  - 實作加密狀態管理：`EnableEncryption`, `DisableEncryption`
  - 實作密碼驗證整合：`DecryptWithPassword`, `EncryptWithPassword`
  - 實作生物識別驗證整合：支援 Touch ID/Face ID 解密
  - 實作加密狀態查詢：`IsEncrypted`, `GetEncryptionType`
  - 實作安全的加密檔案處理：`decryptFileContent`, `encryptFileContent`
  - 支援多種加密演算法：AES-256 和 ChaCha20-Poly1305
  - 實作加密檔案的自動識別（.enc 副檔名）

### 🧪 測試改進

- 新增 `editor_service_test.go` 包含 20+ 個測試函數
  - 測試編輯器服務建立和初始化
  - 測試筆記建立、開啟、保存、更新功能
  - 測試加密檔案識別和處理
  - 測試 Markdown 預覽的各種語法轉換
  - 測試表格、任務列表等擴展功能
  - 測試檔案名稱清理和安全性
  - 測試活躍筆記管理功能
  - 測試錯誤處理：不存在檔案、空筆記、無效輸入
  - 新增加密功能測試：`TestEnableEncryption`, `TestDisableEncryption`
  - 新增加密狀態查詢測試：`TestIsEncrypted`, `TestGetEncryptionType`
  - 新增加密檔案開啟測試：`TestOpenEncryptedNote`
  - 實作完整的模擬服務：`mockEncryptionService`, `mockPasswordService`, `mockBiometricService`
  - 包含模擬檔案系統 (`mockFileRepository`) 用於隔離測試

### 📝 文件更新

- 所有編輯器服務程式碼都包含詳細的繁體中文註解
- 每個編輯器函數都有完整的參數、回傳值和執行流程說明
- Markdown 解析器配置包含詳細的功能說明
- 測試程式碼包含清楚的測試目的和驗證邏輯說明
- 所有加密整合程式碼都包含詳細的繁體中文註解
- 每個加密相關函數都有完整的參數、回傳值和執行流程說明
- 加密檔案處理流程包含詳細的安全性考量說明
- 生物識別整合包含回退機制的詳細說明

### 🔧 技術改進

- 使用 goldmark 作為高效能 Markdown 解析器
- 支援 GitHub Flavored Markdown 標準
- 實作活躍筆記快取機制提高效能
- 實作安全的檔案名稱處理防止檔案系統錯誤
- 使用 UUID 生成唯一筆記識別碼
- 實作完整的錯誤處理和狀態管理
- 整合多層級安全驗證：密碼 + 生物識別
- 實作加密檔案的透明處理：自動加密/解密
- 支援加密狀態的動態切換
- 實作安全的憑證管理架構（為未來 Keychain 整合做準備）

### 📋 待辦事項

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

## [0.4.0] - 2024-01-XX - 加密功能完成

### ✅ 新增功能

- **Task 4.1: 建立加密服務基礎架構**

  - 實作 `EncryptionService` 介面的完整實作 (`encryptionService`)
  - 實作 AES-256-GCM 加密演算法：`encryptWithAES`, `decryptWithAES`
  - 實作 ChaCha20-Poly1305 加密演算法：`encryptWithChaCha20`, `decryptWithChaCha20`
  - 實作 PBKDF2 金鑰衍生函數，使用 SHA-256 和 100,000 次迭代
  - 實作加密資料結構 `EncryptedData` 包含版本、演算法、鹽值、隨機數、資料和校驗和
  - 實作密碼強度驗證 `ValidatePassword`，要求大小寫字母、數字和特殊字元
  - 實作完整的錯誤處理和資料完整性驗證
  - 為生物識別驗證預留介面（將在 Task 4.3 實作）

- **Task 4.2: 實作密碼驗證系統**

  - 建立 `PasswordService` 介面和 `passwordService` 實作
  - 實作安全的密碼雜湊：`HashPassword` 使用 PBKDF2-SHA256 和隨機鹽值
  - 實作密碼驗證：`VerifyPassword` 使用安全比較函數防止時序攻擊
  - 實作密碼強度檢查：`CheckPasswordStrength` 評估密碼強度並提供改進建議
  - 實作密碼重試機制：最多 3 次嘗試，失敗後鎖定 5 分鐘
  - 實作重試狀態管理：`RecordFailedAttempt`, `IsLocked`, `ResetRetryCount`, `GetRetryInfo`
  - 實作常見弱密碼檢測和重複字元檢查
  - 實作執行緒安全的並發存取保護

- **Task 4.3: 整合 macOS 生物驗證**
  - 建立 `BiometricService` 介面和平台特定實作
  - 實作 macOS 生物識別驗證：使用 CGO 調用 LocalAuthentication API
  - 支援 Touch ID 和 Face ID 自動檢測和驗證
  - 實作生物識別可用性檢查：`IsAvailable`, `getBiometricType`
  - 實作生物識別驗證流程：`Authenticate`, `performBiometricAuthentication`
  - 實作筆記級別的生物識別管理：`SetupForNote`, `RemoveForNote`, `IsEnabledForNote`
  - 實作驗證失敗的回退機制：自動回退到密碼驗證
  - 實作跨平台支援：macOS 完整功能，其他平台回退實作
  - 整合到 `EncryptionService`：更新生物識別相關方法

### 🧪 測試改進

- 新增 `encryption_service_test.go` 包含 15+ 個測試函數
  - 測試 AES-256 和 ChaCha20-Poly1305 加密解密功能
  - 測試密碼強度驗證的各種情況
  - 測試錯誤處理：錯誤密碼、無效資料、跨演算法解密
  - 測試邊界條件：空內容、長內容、特殊字元
  - 包含效能測試 (Benchmark) 評估加密解密效能
  - 測試生物識別驗證的未實作狀態
- 新增 `password_service_test.go` 包含 12+ 個測試函數
  - 測試密碼雜湊和驗證的正確性和安全性
  - 測試密碼強度檢查的各種情況和建議生成
  - 測試重試機制：失敗記錄、鎖定狀態、重置功能
  - 測試錯誤處理：無效輸入、邊界條件
  - 測試並發存取的執行緒安全性
  - 包含效能測試評估雜湊和驗證效能
- 新增 `biometric_service_test.go` 包含 12+ 個整合測試函數
  - 測試生物識別可用性檢查和類型識別
  - 測試跨平台行為：macOS 功能測試，其他平台回退測試
  - 測試筆記級別的生物識別管理功能
  - 測試錯誤處理：無效輸入、不可用狀態、未啟用筆記
  - 測試與 EncryptionService 的整合
  - 測試並發存取的執行緒安全性
  - 包含效能測試評估生物識別操作效能

### 📝 文件更新

- 所有加密和生物識別服務程式碼都包含詳細的繁體中文註解
- 每個加密和驗證函數都有完整的參數、回傳值和執行流程說明
- 加密演算法和生物識別實作包含安全性考量和最佳實踐說明
- CGO 程式碼包含 Objective-C 和 C 語言的詳細註解

### 🔧 技術改進

- 使用業界標準的加密演算法和參數設定
- 實作安全的隨機數生成和鹽值管理
- 加入資料完整性校驗防止篡改
- 支援多種加密演算法的靈活切換
- 實作跨平台生物識別驗證架構
- 使用 CGO 安全調用 macOS 系統 API
- 實作完整的錯誤處理和回退機制
- 支援多種生物識別類型（Touch ID、Face ID）

### 📋 待辦事項

- [ ] 5. 建立編輯器服務
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

### 📊 程式碼統計 (截至 v0.7.0)

- **總檔案數**: 22+ 個 Go 檔案
- **程式碼行數**: 5000+ 行 (包含註解)
- **測試檔案**: 9 個
- **測試函數**: 85+ 個
- **測試覆蓋率**: 高 (包含正常情況、邊界條件、錯誤處理、並發測試、UI 測試)

### 🏗️ 架構完成度

- ✅ **資料層 (Models)**: 100% 完成
- ✅ **儲存庫層 (Repositories)**: 檔案操作 100% 完成
- ✅ **服務層 (Services)**: 檔案管理 100% 完成
- ✅ **服務層 (Services)**: 加密服務 100% 完成
- ✅ **服務層 (Services)**: 編輯器服務 100% 完成
- ✅ **服務層 (Services)**: 自動保存服務 100% 完成
- 🚧 **使用者介面層 (UI)**: 40% 完成 (主視窗架構和檔案樹)

### 🎯 里程碑達成

- [x] **里程碑 1**: 專案基礎架構 (v0.1.0)
- [x] **里程碑 2**: 核心資料模型 (v0.2.0)
- [x] **里程碑 3**: 檔案系統操作 (v0.3.0)
- [x] **里程碑 4**: 加密和安全功能 (v0.4.0)
- [x] **里程碑 5**: 編輯器核心功能 (v0.5.0)
- [x] **里程碑 6**: 自動保存系統 (v0.6.0)
- [x] **里程碑 7**: UI 基礎架構 (v0.7.0)
- [ ] **里程碑 8**: 完整使用者介面
- [ ] **里程碑 9**: 完整應用程式

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
