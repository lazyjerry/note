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

#### 主視窗結構 (仿 macOS 備忘錄三欄式布局)

```
MainWindow (1200x800)
├── MenuBar (macOS 原生選單)
│   ├── File Menu
│   │   ├── New Note (⌘N)
│   │   ├── New Folder (⇧⌘N)
│   │   ├── Import... (⌘I)
│   │   ├── Export... (⌘E)
│   │   └── Settings... (⌘,)
│   ├── Edit Menu
│   │   ├── Undo (⌘Z)
│   │   ├── Redo (⇧⌘Z)
│   │   ├── Find (⌘F)
│   │   ├── Find and Replace (⌥⌘F)
│   │   └── Select All (⌘A)
│   ├── Format Menu
│   │   ├── Bold (⌘B)
│   │   ├── Italic (⌘I)
│   │   ├── Strikethrough
│   │   ├── Code Block
│   │   └── Insert Link (⌘K)
│   └── View Menu
│       ├── Show/Hide Folders (⌘1)
│       ├── Show/Hide Note List (⌘2)
│       ├── Show/Hide Preview (⌘3)
│       ├── Toggle Dark Mode (⌘D)
│       └── Enter Full Screen (⌃⌘F)
│
├── Toolbar (隱藏式，可選顯示)
│   ├── Back/Forward Navigation
│   ├── New Note Button
│   ├── Delete Button
│   ├── Share Button
│   └── Search Field
│
├── Content Area (三欄式布局)
│   ├── Sidebar (220px, 可調整)
│   │   ├── Header
│   │   │   ├── "資料夾" 標題
│   │   │   └── New Folder Button (+)
│   │   ├── Folder List (Tree View)
│   │   │   ├── 📁 所有筆記 (All Notes)
│   │   │   ├── 📁 最近使用 (Recently Used)
│   │   │   ├── 📁 我的最愛 (Favorites) ⭐
│   │   │   ├── 📁 已刪除 (Trash) 🗑️
│   │   │   ├── ─────────────────
│   │   │   ├── 📁 工作 (Work)
│   │   │   ├── 📁 個人 (Personal)
│   │   │   ├── 📁 專案 (Projects)
│   │   │   └── 📁 想法 (Ideas)
│   │   └── Footer
│   │       ├── Storage Usage Indicator
│   │       └── Sync Status
│   │
│   ├── Note List Panel (280px, 可調整)
│   │   ├── Header
│   │   │   ├── Folder Name / Search Results
│   │   │   ├── Sort Options (日期/標題/修改時間)
│   │   │   └── View Options (列表/網格)
│   │   ├── Search Bar
│   │   │   ├── 🔍 Search Field
│   │   │   └── Filter Options
│   │   ├── Note List (Scrollable)
│   │   │   ├── Note Item 1
│   │   │   │   ├── 📝 Note Title
│   │   │   │   ├── Preview Text (2-3 lines)
│   │   │   │   ├── Last Modified Date
│   │   │   │   └── Status Icons (🔒 encrypted, ⭐ favorite)
│   │   │   ├── Note Item 2
│   │   │   └── ...
│   │   └── Footer
│   │       ├── Note Count
│   │       └── New Note Button (+)
│   │
│   └── Editor Panel (剩餘空間, 最小 400px)
│       ├── Header (可隱藏)
│       │   ├── Note Title Field
│       │   ├── Action Buttons
│       │   │   ├── 🔒 Encryption Toggle
│       │   │   ├── ⭐ Favorite Toggle
│       │   │   ├── 📤 Share Button
│       │   │   ├── 🗑️ Delete Button
│       │   │   └── ⋯ More Options
│       │   └── Format Toolbar (可隱藏)
│       │       ├── Bold, Italic, Strikethrough
│       │       ├── Headers (H1-H6)
│       │       ├── Lists (Bullet, Numbered)
│       │       ├── Code Block, Quote
│       │       └── Link, Image
│       ├── Content Area (可分割)
│       │   ├── Editor View (Markdown)
│       │   │   ├── Line Numbers (可選)
│       │   │   ├── Syntax Highlighting
│       │   │   ├── Live Word Count
│       │   │   └── Auto-completion
│       │   └── Preview View (可選分割或全屏)
│       │       ├── Rendered HTML
│       │       ├── Table of Contents (可選)
│       │       └── Export Options
│       └── Footer
│           ├── Status Indicators
│           │   ├── 💾 Auto-save Status
│           │   ├── 🔒 Encryption Status
│           │   ├── 📊 Word/Character Count
│           │   └── 🌐 Sync Status
│           └── View Controls
│               ├── Zoom Level (50%-200%)
│               ├── 📝 Edit Mode Toggle
│               ├── 👁️ Preview Mode Toggle
│               ├── ⚡ Split View Toggle
│               └── 🎯 Focus Mode Toggle
│
└── Floating Elements
    ├── Search Overlay (⌘F)
    │   ├── Search Field
    │   ├── Replace Field
    │   ├── Match Count
    │   └── Navigation Buttons
    ├── Quick Actions Palette (⌘P)
    │   ├── Recent Commands
    │   ├── File Operations
    │   └── Settings Shortcuts
    └── Notifications
        ├── Auto-save Notifications
        ├── Sync Status Updates
        └── Error Messages
```

#### UI 設計原則 (遵循 macOS Human Interface Guidelines)

1. **視覺層次**

   - 使用 macOS 原生的視覺分層
   - 適當的陰影和邊框
   - 一致的間距和對齊

2. **顏色系統**

   - 支援淺色/深色模式自動切換
   - 使用 macOS 系統顏色
   - 高對比度支援

3. **字體系統**

   - 主要文字：SF Pro Text
   - 程式碼：SF Mono
   - 標題：SF Pro Display
   - 支援動態字體大小

4. **互動設計**

   - 原生 macOS 手勢支援
   - 鍵盤快捷鍵遵循 macOS 慣例
   - 拖放操作支援
   - 右鍵選單整合

5. **響應式布局**
   - 最小視窗大小：800x600
   - 面板可調整大小
   - 記住用戶偏好設定
   - 全螢幕模式支援

#### 特殊 UI 元件

##### 1. 智慧搜尋欄

```
Search Bar
├── 🔍 Search Icon
├── Search Field (with placeholder)
├── Filter Dropdown
│   ├── 📝 Content
│   ├── 📁 Folder
│   ├── 🏷️ Tags
│   └── 📅 Date Range
└── Recent Searches
```

##### 2. 筆記預覽卡片

```
Note Card
├── Header
│   ├── Note Title
│   ├── Status Icons (🔒, ⭐, 📎)
│   └── Last Modified
├── Content Preview (3 lines)
├── Tags (if any)
└── Footer
    ├── Word Count
    └── Folder Location
```

##### 3. 加密狀態指示器

```
Encryption Indicator
├── 🔓 Unlocked (Green)
├── 🔒 Locked (Orange)
├── 🔐 Encrypted (Blue)
└── ⚠️ Error (Red)
```

##### 4. 自動保存狀態

```
Auto-save Status
├── 💾 Saved (Gray)
├── ⏳ Saving... (Blue, animated)
├── ✅ Auto-saved (Green, fade out)
└── ❌ Save failed (Red)
```

#### 視圖模式系統 (View Mode System)

##### 1. 編輯器視圖模式

```
View Mode Controller
├── Edit Only Mode (編輯模式)
│   ├── 隱藏預覽面板
│   ├── 編輯器佔滿右側空間
│   ├── 快捷鍵: ⌘1
│   └── 專注寫作體驗
├── Preview Only Mode (預覽模式)
│   ├── 隱藏編輯器面板
│   ├── 預覽佔滿右側空間
│   ├── 快捷鍵: ⌘2
│   └── 閱讀和檢視體驗
├── Split View Mode (分割視圖)
│   ├── 編輯器和預覽並排顯示
│   ├── 可調整分割比例
│   ├── 同步滾動支援
│   ├── 快捷鍵: ⌘3
│   └── 即時預覽體驗
└── Focus Mode (專注模式)
    ├── 隱藏側邊欄和筆記列表
    ├── 全螢幕編輯體驗
    ├── 快捷鍵: ⌃⌘F
    └── 無干擾寫作環境
```

##### 2. 工具列重新設計

```
Top Toolbar (頂部工具列)
├── Left Section
│   ├── View Mode Toggles
│   │   ├── 📝 Edit Button
│   │   ├── 👁️ Preview Button
│   │   └── ⚡ Split Button
│   └── Quick Actions
│       ├── 🔍 Search
│       └── ⚙️ Settings
├── Center Section
│   ├── Document Title
│   └── Breadcrumb Navigation
└── Right Section
    ├── Sync Status
    ├── Word Count
    └── User Avatar

Side Toolbar (側邊工具列)
├── File Operations
│   ├── 📄 New Note
│   ├── 📁 New Folder
│   ├── 📤 Import
│   └── 📥 Export
├── Format Tools
│   ├── 𝐁 Bold
│   ├── 𝐼 Italic
│   ├── 𝐔 Underline
│   ├── 🔗 Link
│   └── 🖼️ Image
└── Advanced Tools
    ├── 🔒 Encryption
    ├── ⭐ Favorite
    ├── 🏷️ Tags
    └── 📊 Statistics
```

##### 3. 響應式佈局適應

```
Window Size Adaptations:
├── Large (1200px+)
│   ├── 三欄完整顯示
│   ├── 所有工具列可見
│   └── 最佳使用體驗
├── Medium (800-1199px)
│   ├── 可收合側邊欄
│   ├── 簡化工具列
│   └── 保持核心功能
└── Small (600-799px)
    ├── 單欄顯示模式
    ├── 抽屜式導航
    └── 觸控友好介面
```

#### 繁體中文輸入優化設計

##### 1. 輸入法整合架構

```
Chinese Input System
├── Input Method Engine (輸入法引擎)
│   ├── Zhuyin (注音) Support
│   ├── Pinyin (拼音) Support
│   ├── Cangjie (倉頡) Support
│   └── Quick (速成) Support
├── Candidate Window (候選字視窗)
│   ├── Native macOS Style
│   ├── Customizable Position
│   ├── Font Size Adaptation
│   └── Dark Mode Support
├── Composition Display (組字顯示)
│   ├── Inline Composition
│   ├── Underline Styling
│   ├── Tone Mark Display
│   └── Real-time Feedback
└── Text Rendering (文字渲染)
    ├── CJK Font Optimization
    ├── Character Spacing
    ├── Line Height Adjustment
    └── Unicode Normalization
```

##### 2. 中文編輯體驗優化

```
Chinese Editing Features:
├── Smart Input
│   ├── Auto Punctuation
│   ├── Smart Quotes (「」『』)
│   ├── Number Conversion (1→一)
│   └── Date Format (2024/1/1→2024年1月1日)
├── Text Selection
│   ├── Word Boundary Detection
│   ├── Phrase Selection
│   ├── Double-click Word Selection
│   └── Triple-click Paragraph Selection
├── Typography
│   ├── Proper Line Breaking
│   ├── Punctuation Hanging
│   ├── Vertical Text Support (未來)
│   └── Traditional/Simplified Toggle
└── Search & Replace
    ├── Fuzzy Pinyin Search
    ├── Traditional/Simplified Match
    ├── Regex with CJK Support
    └── Tone-insensitive Search
```

##### 3. 字體和渲染系統

```css
/* 中文字體優化 */
.chinese-text {
	font-family: "PingFang TC", /* macOS 繁體中文主字體 */ "Hiragino Sans TC", /* 備用繁體中文字體 */ "Microsoft JhengHei", /* Windows 繁體中文字體 */ "Noto Sans CJK TC", /* 跨平台 CJK 字體 */ sans-serif;

	/* 中文字體渲染優化 */
	text-rendering: optimizeLegibility;
	-webkit-font-smoothing: antialiased;
	-moz-osx-font-smoothing: grayscale;

	/* 中文文字間距 */
	letter-spacing: 0.05em;
	word-spacing: 0.1em;

	/* 行高優化 */
	line-height: 1.7;
}

/* 程式碼中的中文註解 */
.code-chinese-comment {
	font-family: "SF Mono", "PingFang TC", "Hiragino Sans TC", monospace;
	color: var(--comment-color);
	font-style: normal; /* 中文註解不使用斜體 */
}
```

#### 動畫與轉場效果

##### 1. 面板切換動畫

- 側邊欄展開/收合：0.3s ease-in-out
- 筆記列表載入：淡入效果 0.2s
- 編輯器內容切換：交叉淡化 0.15s

##### 2. 狀態變化動畫

- 保存狀態：脈衝效果
- 加密狀態：圖示旋轉
- 搜尋結果：高亮顯示

##### 3. 互動回饋

- 按鈕點擊：輕微縮放效果
- 拖放操作：半透明預覽
- 選擇狀態：漸變背景色

#### 無障礙設計

##### 1. VoiceOver 支援

- 所有 UI 元素都有適當的標籤
- 鍵盤導航支援
- 螢幕閱讀器友好的內容結構

##### 2. 鍵盤操作

```
全域快捷鍵：
⌘N          新建筆記
⌘O          開啟檔案
⌘S          保存
⌘F          搜尋
⌘G          搜尋下一個
⇧⌘G         搜尋上一個
⌘W          關閉筆記
⌘Q          退出應用程式

導航快捷鍵：
⌘1          顯示/隱藏資料夾面板
⌘2          顯示/隱藏筆記列表
⌘3          顯示/隱藏預覽面板
⌘↑/↓        在筆記列表中導航
⌘←/→        在面板間切換焦點
Tab/⇧Tab    在元件間導航

編輯快捷鍵：
⌘B          粗體
⌘I          斜體
⌘U          底線
⌘K          插入連結
⌘⇧K         移除連結
⌘L          插入列表
⌘⇧L         插入編號列表
⌘E          插入程式碼
⌘⇧E         插入程式碼區塊
```

##### 3. 高對比度支援

- 遵循系統高對比度設定
- 自訂顏色主題選項
- 文字大小動態調整

#### 效能優化設計

##### 1. 虛擬化列表

- 筆記列表使用虛擬滾動
- 只渲染可見項目
- 智慧預載入機制

##### 2. 延遲載入

- 筆記內容按需載入
- 圖片和附件延遲載入
- 預覽內容快取機制

##### 3. 記憶體管理

- 自動釋放未使用的筆記
- 圖片記憶體快取限制
- 定期垃圾回收

## 資料模型

### 檔案結構

```
~/Documents/NotebookApp/
├── notes/                          # 筆記檔案主目錄
│   ├── 📁 工作/                    # 工作相關筆記
│   │   ├── 📝 會議記錄.md
│   │   ├── 📝 專案計劃.md
│   │   └── 🔒 機密文件.md.enc
│   ├── 📁 個人/                    # 個人筆記
│   │   ├── 📝 日記.md.enc
│   │   ├── 📝 想法收集.md
│   │   └── 📁 旅行/
│   │       ├── 📝 日本行程.md
│   │       └── 📝 美食清單.md
│   ├── 📁 專案/                    # 專案筆記
│   │   ├── 📁 App開發/
│   │   │   ├── 📝 需求分析.md
│   │   │   ├── 📝 技術選型.md
│   │   │   └── 📝 進度追蹤.md
│   │   └── 📁 學習筆記/
│   │       ├── 📝 Go語言.md
│   │       └── 📝 設計模式.md
│   └── 📝 快速筆記.md              # 臨時筆記
├── .notebook/                      # 應用程式資料目錄
│   ├── settings.json               # 用戶設定
│   ├── index.db                    # 筆記索引資料庫
│   ├── cache/                      # 快取目錄
│   │   ├── thumbnails/             # 筆記縮圖
│   │   ├── search_index/           # 搜尋索引
│   │   └── preview_cache/          # 預覽快取
│   ├── keys/                       # 加密金鑰管理
│   │   ├── biometric_keys.json     # 生物識別金鑰
│   │   ├── password_hashes.json    # 密碼雜湊
│   │   └── master_key.enc          # 主金鑰（加密）
│   ├── logs/                       # 日誌檔案
│   │   ├── app_2024-01-15.log
│   │   ├── error_2024-01-15.log
│   │   └── performance_2024-01-15.log
│   └── plugins/                    # 外掛目錄
│       ├── markdown_extensions/
│       └── export_formats/
├── backups/                        # 備份目錄
│   ├── auto_backup_20240115/       # 自動備份
│   ├── manual_backup_20240110/     # 手動備份
│   └── export_backup_20240105/     # 匯出備份
├── templates/                      # 範本目錄
│   ├── 📝 會議記錄範本.md
│   ├── 📝 日報範本.md
│   ├── 📝 專案計劃範本.md
│   └── 📝 學習筆記範本.md
└── attachments/                    # 附件目錄
    ├── images/                     # 圖片附件
    ├── documents/                  # 文件附件
    └── media/                      # 媒體附件
```

### 資料庫結構 (SQLite)

```sql
-- 筆記索引表
CREATE TABLE notes (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    file_path TEXT UNIQUE NOT NULL,
    folder_path TEXT,
    is_encrypted BOOLEAN DEFAULT FALSE,
    encryption_type TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    last_opened DATETIME,
    word_count INTEGER DEFAULT 0,
    character_count INTEGER DEFAULT 0,
    is_favorite BOOLEAN DEFAULT FALSE,
    tags TEXT, -- JSON array
    preview_text TEXT,
    checksum TEXT -- 檔案完整性檢查
);

-- 資料夾結構表
CREATE TABLE folders (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    parent_id TEXT,
    path TEXT UNIQUE NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    color TEXT, -- 資料夾顏色標籤
    icon TEXT,  -- 自訂圖示
    sort_order INTEGER DEFAULT 0,
    FOREIGN KEY (parent_id) REFERENCES folders(id)
);

-- 標籤表
CREATE TABLE tags (
    id TEXT PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    color TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    usage_count INTEGER DEFAULT 0
);

-- 筆記標籤關聯表
CREATE TABLE note_tags (
    note_id TEXT,
    tag_id TEXT,
    PRIMARY KEY (note_id, tag_id),
    FOREIGN KEY (note_id) REFERENCES notes(id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
);

-- 搜尋歷史表
CREATE TABLE search_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    query TEXT NOT NULL,
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
    result_count INTEGER DEFAULT 0
);

-- 最近開啟表
CREATE TABLE recent_notes (
    note_id TEXT,
    opened_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (note_id),
    FOREIGN KEY (note_id) REFERENCES notes(id) ON DELETE CASCADE
);

-- 應用程式設定表
CREATE TABLE app_settings (
    key TEXT PRIMARY KEY,
    value TEXT,
    type TEXT, -- 'string', 'number', 'boolean', 'json'
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### 主題系統設計

#### 1. 顏色主題 (遵循 macOS 設計語言)

##### 淺色主題 (Light Mode)

```css
/* 主要顏色 */
--primary-background: #FFFFFF
--secondary-background: #F5F5F7
--tertiary-background: #EFEFF4
--sidebar-background: #F7F7F7
--selected-background: #007AFF1A

/* 文字顏色 */
--primary-text: #1D1D1F
--secondary-text: #86868B
--tertiary-text: #C7C7CC
--link-text: #007AFF
--accent-text: #FF3B30

/* 邊框和分隔線 */
--border-color: #D1D1D6
--separator-color: #C6C6C8
--shadow-color: #00000010

/* 狀態顏色 */
--success-color: #30D158
--warning-color: #FF9500
--error-color: #FF3B30
--info-color: #007AFF
```

##### 深色主題 (Dark Mode)

```css
/* 主要顏色 */
--primary-background: #1C1C1E
--secondary-background: #2C2C2E
--tertiary-background: #3A3A3C
--sidebar-background: #242426
--selected-background: #0A84FF1A

/* 文字顏色 */
--primary-text: #FFFFFF
--secondary-text: #98989D
--tertiary-text: #48484A
--link-text: #0A84FF
--accent-text: #FF453A

/* 邊框和分隔線 */
--border-color: #38383A
--separator-color: #48484A
--shadow-color: #00000030

/* 狀態顏色 */
--success-color: #32D74B
--warning-color: #FF9F0A
--error-color: #FF453A
--info-color: #0A84FF
```

#### 2. 字體系統

```css
/* 主要字體 */
--font-family-primary: -apple-system, BlinkMacSystemFont, 'SF Pro Text', sans-serif
--font-family-mono: 'SF Mono', Monaco, 'Cascadia Code', monospace
--font-family-display: -apple-system, BlinkMacSystemFont, 'SF Pro Display', sans-serif

/* 字體大小 */
--font-size-xs: 11px      /* 狀態列、標籤 */
--font-size-sm: 12px      /* 次要資訊 */
--font-size-base: 13px    /* 正文 */
--font-size-lg: 15px      /* 標題 */
--font-size-xl: 17px      /* 大標題 */
--font-size-2xl: 22px     /* 主標題 */
--font-size-3xl: 28px     /* 特大標題 */

/* 行高 */
--line-height-tight: 1.2
--line-height-normal: 1.4
--line-height-relaxed: 1.6
--line-height-loose: 1.8
```

#### 3. 間距系統

```css
/* 間距單位 (基於 4px 網格) */
--spacing-xs: 2px
--spacing-sm: 4px
--spacing-base: 8px
--spacing-md: 12px
--spacing-lg: 16px
--spacing-xl: 20px
--spacing-2xl: 24px
--spacing-3xl: 32px
--spacing-4xl: 40px
--spacing-5xl: 48px

/* 圓角 */
--radius-sm: 4px
--radius-base: 6px
--radius-md: 8px
--radius-lg: 12px
--radius-xl: 16px
--radius-full: 9999px

/* 陰影 */
--shadow-sm: 0 1px 2px rgba(0, 0, 0, 0.05)
--shadow-base: 0 1px 3px rgba(0, 0, 0, 0.1), 0 1px 2px rgba(0, 0, 0, 0.06)
--shadow-md: 0 4px 6px rgba(0, 0, 0, 0.07), 0 2px 4px rgba(0, 0, 0, 0.06)
--shadow-lg: 0 10px 15px rgba(0, 0, 0, 0.1), 0 4px 6px rgba(0, 0, 0, 0.05)
--shadow-xl: 0 20px 25px rgba(0, 0, 0, 0.1), 0 10px 10px rgba(0, 0, 0, 0.04)
```

### 加密檔案格式

```json
{
	"version": "2.0",
	"metadata": {
		"created_at": "2024-01-15T10:30:00Z",
		"algorithm": "aes256-gcm",
		"auth_type": "biometric+password",
		"key_derivation": "pbkdf2",
		"iterations": 100000
	},
	"encryption": {
		"salt": "base64_encoded_salt",
		"iv": "base64_encoded_iv",
		"auth_tag": "base64_encoded_auth_tag",
		"data": "base64_encoded_encrypted_content"
	},
	"integrity": {
		"checksum": "sha256_hash",
		"signature": "digital_signature"
	},
	"backup": {
		"recovery_hint": "encrypted_password_hint",
		"backup_key_id": "uuid_for_backup_key"
	}
}
```

### 設定檔案格式

```json
{
	"version": "1.0",
	"app_settings": {
		"theme": "auto",
		"language": "zh-TW",
		"font_size": 13,
		"font_family": "SF Pro Text",
		"line_height": 1.4,
		"auto_save_interval": 30,
		"backup_enabled": true,
		"backup_interval": 1440
	},
	"ui_settings": {
		"sidebar_width": 220,
		"note_list_width": 280,
		"show_line_numbers": false,
		"show_word_count": true,
		"show_preview": true,
		"split_view_ratio": 0.5,
		"toolbar_visible": true,
		"status_bar_visible": true
	},
	"editor_settings": {
		"tab_size": 4,
		"word_wrap": true,
		"syntax_highlighting": true,
		"auto_completion": true,
		"spell_check": true,
		"markdown_extensions": ["tables", "strikethrough", "task_lists", "code_highlighting", "math_expressions"]
	},
	"security_settings": {
		"default_encryption": "aes256",
		"biometric_enabled": true,
		"auto_lock_timeout": 300,
		"password_complexity": "medium",
		"secure_delete": true
	},
	"sync_settings": {
		"provider": "icloud",
		"auto_sync": true,
		"conflict_resolution": "manual",
		"sync_attachments": true
	},
	"export_settings": {
		"default_format": "markdown",
		"include_metadata": false,
		"preserve_structure": true,
		"image_handling": "embed"
	}
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

### 單元測試 (Unit Tests)

#### 1. 服務層測試

```go
// 測試覆蓋率目標：90%+
func TestEditorService_CreateNote(t *testing.T)
func TestEditorService_SaveNote(t *testing.T)
func TestEditorService_EncryptNote(t *testing.T)
func TestFileManagerService_ListFiles(t *testing.T)
func TestEncryptionService_AES256(t *testing.T)
func TestAutoSaveService_StartStop(t *testing.T)
```

#### 2. 資料模型測試

- **Note 模型驗證**: 標題、內容、加密狀態驗證
- **Settings 模型測試**: 設定值範圍和格式驗證
- **FileInfo 模型測試**: 檔案屬性和路徑處理
- **序列化測試**: JSON 編碼/解碼正確性

#### 3. 加密功能測試

- **演算法測試**: AES-256, ChaCha20-Poly1305
- **金鑰派生測試**: PBKDF2, Argon2
- **完整性驗證**: HMAC, 數位簽章
- **效能測試**: 加密/解密速度基準

### 整合測試 (Integration Tests)

#### 1. UI 整合測試

```go
func TestMainWindow_ThreeColumnLayout(t *testing.T)
func TestNoteList_Selection(t *testing.T)
func TestEditor_MarkdownPreview(t *testing.T)
func TestSidebar_FolderNavigation(t *testing.T)
func TestSearch_RealTimeResults(t *testing.T)
```

#### 2. 檔案系統整合測試

- **跨平台路徑處理**: Windows, macOS, Linux
- **檔案監控測試**: fsnotify 整合
- **大檔案處理**: 100MB+ Markdown 檔案
- **並發存取測試**: 多執行緒檔案操作

#### 3. 加密整合測試

- **端到端加密流程**: 建立 → 加密 → 保存 → 讀取 → 解密
- **生物識別整合**: Touch ID/Face ID 模擬測試
- **密碼管理測試**: 密碼變更、重設流程
- **備份恢復測試**: 加密檔案的備份和恢復

### 效能測試 (Performance Tests)

#### 1. 基準測試

```go
func BenchmarkMarkdownParsing(b *testing.B)
func BenchmarkFileSearch(b *testing.B)
func BenchmarkEncryption(b *testing.B)
func BenchmarkAutoSave(b *testing.B)
```

#### 2. 負載測試

- **大量筆記處理**: 10,000+ 筆記載入
- **記憶體使用監控**: 長時間運行記憶體洩漏檢測
- **CPU 使用率測試**: 高負載下的響應性
- **磁碟 I/O 測試**: 頻繁讀寫操作的效能

#### 3. UI 響應性測試

- **60fps 滾動測試**: 筆記列表流暢滾動
- **即時搜尋效能**: 輸入延遲 < 100ms
- **預覽更新速度**: Markdown 渲染 < 50ms
- **啟動時間測試**: 冷啟動 < 2s, 熱啟動 < 0.5s

### 安全測試 (Security Tests)

#### 1. 加密安全測試

- **密碼學強度驗證**: NIST 標準合規性
- **金鑰管理測試**: 安全金鑰生成和儲存
- **側通道攻擊防護**: 時間攻擊、功耗分析
- **隨機數品質測試**: 熵源驗證

#### 2. 資料保護測試

- **記憶體清理測試**: 敏感資料自動清除
- **檔案權限測試**: 適當的檔案系統權限
- **備份安全測試**: 備份檔案的加密保護
- **日誌安全測試**: 敏感資訊不洩漏到日誌

#### 3. 身份驗證測試

- **生物識別測試**: Touch ID/Face ID 整合安全性
- **密碼強度測試**: 密碼複雜度要求
- **會話管理測試**: 自動鎖定和解鎖機制
- **多因素驗證**: 密碼+生物識別組合

### 用戶體驗測試 (UX Tests)

#### 1. 可用性測試

- **新用戶引導**: 首次使用體驗流程
- **工作流程測試**: 常見使用場景完整測試
- **錯誤恢復測試**: 用戶錯誤操作的恢復能力
- **快捷鍵測試**: 鍵盤操作的完整性

#### 2. 無障礙測試

- **VoiceOver 相容性**: 螢幕閱讀器支援
- **鍵盤導航**: 純鍵盤操作完整性
- **高對比度支援**: 視覺障礙用戶支援
- **字體縮放測試**: 動態字體大小調整

#### 3. 多語言測試

- **國際化測試**: 繁體中文、英文介面
- **文字長度測試**: 不同語言的 UI 適應性
- **RTL 語言支援**: 阿拉伯文、希伯來文（未來）
- **日期時間格式**: 本地化格式正確性

#### 4. 中文輸入法測試

- **注音輸入測試**: 候選字視窗顯示正確性
- **拼音輸入測試**: 智慧選字和詞組預測
- **倉頡輸入測試**: 複雜字符輸入準確性
- **混合輸入測試**: 中英文混合輸入流暢性
- **特殊符號測試**: 中文標點符號正確處理
- **字體渲染測試**: 不同字體大小下的清晰度

#### 5. UI/UX 改進測試

- **視圖切換測試**: 編輯/預覽模式切換流暢性
- **佈局響應測試**: 不同視窗大小下的適應性
- **工具列測試**: 側邊欄和頂部工具列功能完整性
- **預覽功能測試**: 獨立預覽視窗的功能驗證
- **快捷鍵測試**: 所有鍵盤快捷鍵的響應性
- **動畫效果測試**: 視圖切換動畫的流暢度
- **記憶功能測試**: 視圖狀態和佈局偏好的持久化

### 自動化測試 (Automated Tests)

#### 1. CI/CD 整合

```yaml
# GitHub Actions 工作流程
name: Test Suite
on: [push, pull_request]
jobs:
  test:
    runs-on: macos-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v3
      - run: go test -v -race -coverprofile=coverage.out ./...
      - run: go tool cover -html=coverage.out -o coverage.html
```

#### 2. 測試環境

- **本地開發環境**: 快速反饋循環
- **CI 環境**: 自動化回歸測試
- **預發布環境**: 完整功能驗證
- **效能測試環境**: 專用硬體基準測試

#### 3. 測試資料管理

- **測試資料生成**: 自動生成測試筆記和資料夾
- **測試環境隔離**: 獨立的測試資料目錄
- **測試清理**: 自動清理測試產生的檔案
- **快照測試**: UI 元件視覺回歸測試

### 品質保證 (Quality Assurance)

#### 1. 程式碼品質

- **靜態分析**: golint, go vet, staticcheck
- **程式碼覆蓋率**: 目標 85%+
- **複雜度分析**: 循環複雜度監控
- **依賴性分析**: 第三方套件安全性掃描

#### 2. 文件測試

- **API 文件**: 自動生成和驗證
- **使用者手冊**: 功能完整性檢查
- **程式碼註解**: 文件覆蓋率檢查
- **範例程式碼**: 可執行性驗證

#### 3. 發布前檢查

- **功能完整性**: 所有需求功能驗證
- **效能基準**: 與前版本效能比較
- **安全掃描**: 漏洞和安全問題檢查
- **相容性測試**: 不同 macOS 版本相容性
