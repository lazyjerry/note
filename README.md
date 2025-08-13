# Mac Notebook App

一個使用 Go 和 Fyne 建立的安全 Markdown 筆記應用程式，專為 macOS 設計。

## 功能特色

- Markdown 編輯與即時預覽
- 檔案和資料夾管理
- 密碼和生物識別加密
- 自動保存功能
- 跨平台相容性

## 專案結構

```
mac-notebook-app/
├── main.go                     # 應用程式入口點
├── go.mod                      # Go 模組定義
├── internal/                   # 內部應用程式程式碼
│   ├── models/                 # 資料模型和結構體
│   │   ├── note.go            # 筆記模型
│   │   ├── settings.go        # 設定模型
│   │   ├── file_info.go       # 檔案資訊模型
│   │   └── errors.go          # 錯誤定義
│   ├── services/              # 業務邏輯服務
│   │   └── interfaces.go      # 服務介面
│   └── repositories/          # 資料存取層
│       └── interfaces.go      # 儲存庫介面
└── ui/                        # 使用者介面元件
    └── main_window.go         # 主視窗實作
```

## 系統需求

- Go 1.21 或更新版本
- macOS（用於生物識別驗證功能）

## 相依套件

- Fyne v2.4+ - 跨平台 GUI 框架
- goldmark - Markdown 解析器
- fsnotify - 檔案系統通知
- golang.org/x/crypto - 加密函數

## 開發狀態

此專案目前正在開發中。基本的專案結構和核心介面已經建立完成。

## 後續步驟

1. 實作資料模型和驗證
2. 建立檔案系統操作
3. 添加加密功能
4. 建立使用者介面元件
5. 整合所有服務

## 授權條款

詳見 LICENSE 檔案。
