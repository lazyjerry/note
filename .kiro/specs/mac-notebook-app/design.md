# 設計文檔

## 概述

Mac 筆記本應用是一個基於 Golang 和 Fyne 框架開發的桌面應用程式。該應用採用模組化架構，支援 Markdown 編輯、檔案管理、多層級加密保護和智慧自動保存功能。應用程式設計重點在於用戶體驗、資料安全性和跨平台相容性。

## 架構

### 整體架構

```
┌─────────────────────────────────────────┐
│                UI Layer                 │
│  (Fyne Widgets & Custom Components)     │
├─────────────────────────────────────────┤
│              Service Layer              │
│  ┌─────────────┬─────────────────────┐  │
│  │   Editor    │    File Manager     │  │
│  │   Service   │      Service        │  │
│  ├─────────────┼─────────────────────┤  │
│  │  Encryption │    Auto Save        │  │
│  │   Service   │     Service         │  │
│  └─────────────┴─────────────────────┘  │
├─────────────────────────────────────────┤
│             Repository Layer            │
│  ┌─────────────┬─────────────────────┐  │
│  │    File     │     Settings        │  │
│  │ Repository  │    Repository       │  │
│  └─────────────┴─────────────────────┘  │
├─────────────────────────────────────────┤
│              Data Layer                 │
│  ┌─────────────┬─────────────────────┐  │
│  │ File System │    Keychain/        │  │
│  │   Storage   │  Biometric Auth     │  │
│  └─────────────┴─────────────────────┘  │
└─────────────────────────────────────────┘
```

### 技術棧

- **UI 框架**: Fyne v2.4+
- **程式語言**: Go 1.21+
- **加密庫**: crypto/aes, golang.org/x/crypto/chacha20poly1305
- **Markdown 處理**: github.com/yuin/goldmark
- **生物驗證**: LocalAuthentication (透過 CGO 調用 macOS API)
- **檔案監控**: github.com/fsnotify/fsnotify

## 元件與介面

### 1. 核心資料模型

#### Note 模型

```go
type Note struct {
    ID          string    `json:"id"`
    Title       string    `json:"title"`
    Content     string    `json:"content"`
    FilePath    string    `json:"file_path"`
    IsEncrypted bool      `json:"is_encrypted"`
    EncryptionType string `json:"encryption_type"` // "password", "biometric", "both"
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    LastSaved   time.Time `json:"last_saved"`
}
```

#### Settings 模型

```go
type Settings struct {
    DefaultEncryption   string `json:"default_encryption"`   // "aes256", "chacha20"
    AutoSaveInterval    int    `json:"auto_save_interval"`    // minutes
    DefaultSaveLocation string `json:"default_save_location"`
    BiometricEnabled    bool   `json:"biometric_enabled"`
    Theme              string `json:"theme"`                 // "light", "dark", "auto"
}
```

### 2. 服務介面

#### EditorService

```go
type EditorService interface {
    CreateNote(title, content string) (*Note, error)
    OpenNote(filePath string) (*Note, error)
    SaveNote(note *Note) error
    UpdateContent(noteID, content string) error
    PreviewMarkdown(content string) string
}
```

#### FileManagerService

```go
type FileManagerService interface {
    ListFiles(directory string) ([]FileInfo, error)
    CreateDirectory(path string) error
    DeleteFile(path string) error
    RenameFile(oldPath, newPath string) error
    MoveFile(sourcePath, destPath string) error
}
```

#### EncryptionService

```go
type EncryptionService interface {
    EncryptContent(content, password string, algorithm string) ([]byte, error)
    DecryptContent(encryptedData []byte, password string, algorithm string) (string, error)
    SetupBiometricAuth(noteID string) error
    AuthenticateWithBiometric(noteID string) (bool, error)
    ValidatePassword(password string) bool
}
```

#### AutoSaveService

```go
type AutoSaveService interface {
    StartAutoSave(note *Note, interval time.Duration)
    StopAutoSave(noteID string)
    SaveNow(noteID string) error
    GetSaveStatus(noteID string) SaveStatus
}
```

### 3. UI 元件架構

#### 主視窗結構

```
MainWindow
├── MenuBar
│   ├── File Menu (New, Open, Save, Settings)
│   ├── Edit Menu (Undo, Redo, Find)
│   └── View Menu (Theme, Preview)
├── ToolBar
│   ├── New Note Button
│   ├── Save Button
│   ├── Encryption Toggle
│   └── Preview Toggle
├── Content Area (HSplit)
│   ├── Left Panel (VSplit)
│   │   ├── File Tree (30%)
│   │   └── Note List (70%)
│   └── Right Panel
│       ├── Editor Area (70%)
│       └── Preview Area (30%)
└── Status Bar
    ├── Save Status
    ├── Encryption Status
    └── Word Count
```

## 資料模型

### 檔案結構

```
~/Documents/NotebookApp/
├── notes/
│   ├── folder1/
│   │   ├── note1.md
│   │   └── note2.md.enc
│   └── note3.md
├── .notebook/
│   ├── settings.json
│   ├── index.db
│   └── keys/
│       ├── biometric_keys.json
│       └── password_hashes.json
└── backups/
    └── auto_backup_YYYYMMDD/
```

### 加密檔案格式

```json
{
	"version": "1.0",
	"algorithm": "aes256",
	"auth_type": "password",
	"salt": "base64_encoded_salt",
	"iv": "base64_encoded_iv",
	"data": "base64_encoded_encrypted_content",
	"checksum": "sha256_hash"
}
```

## 錯誤處理

### 錯誤類型定義

```go
type AppError struct {
    Code    string `json:"code"`
    Message string `json:"message"`
    Details string `json:"details,omitempty"`
}

const (
    ErrFileNotFound     = "FILE_NOT_FOUND"
    ErrInvalidPassword  = "INVALID_PASSWORD"
    ErrEncryptionFailed = "ENCRYPTION_FAILED"
    ErrBiometricFailed  = "BIOMETRIC_FAILED"
    ErrSaveFailed       = "SAVE_FAILED"
    ErrPermissionDenied = "PERMISSION_DENIED"
)
```

### 錯誤處理策略

1. **檔案操作錯誤**

   - 顯示用戶友好的錯誤訊息
   - 提供重試機制
   - 記錄詳細錯誤日誌

2. **加密錯誤**

   - 密碼錯誤：提供重新輸入機會（最多 3 次）
   - 生物驗證失敗：回退到密碼驗證
   - 加密演算法錯誤：使用預設演算法

3. **自動保存錯誤**
   - 顯示保存失敗通知
   - 提供手動保存選項
   - 暫停自動保存直到問題解決

## 測試策略

### 單元測試

- **服務層測試**: 每個服務介面的完整測試覆蓋
- **加密功能測試**: 各種加密演算法的加解密測試
- **檔案操作測試**: 檔案 CRUD 操作的測試
- **資料模型測試**: 資料驗證和序列化測試

### 整合測試

- **UI 整合測試**: 使用 Fyne 的測試框架測試 UI 互動
- **檔案系統整合測試**: 實際檔案操作的端到端測試
- **加密整合測試**: 完整的加密工作流程測試

### 效能測試

- **大檔案處理測試**: 測試大型 Markdown 檔案的處理效能
- **自動保存效能測試**: 測試自動保存對系統效能的影響
- **記憶體使用測試**: 長時間運行的記憶體洩漏測試

### 安全測試

- **加密強度測試**: 驗證加密演算法的實作正確性
- **密碼安全測試**: 測試密碼儲存和驗證的安全性
- **生物驗證測試**: 測試 macOS 生物驗證整合的安全性

### 用戶體驗測試

- **響應性測試**: 確保 UI 操作的即時響應
- **錯誤處理測試**: 驗證錯誤訊息的清晰度和有用性
- **工作流程測試**: 測試完整的用戶使用場景
