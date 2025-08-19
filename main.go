// Package main 是 Mac Notebook App 的主要入口點
// 這是一個安全的 Markdown 筆記應用程式，支援密碼和生物識別加密
package main

import (
	"fyne.io/fyne/v2"          // Fyne GUI 框架核心套件，提供基礎介面類型和功能
	"fyne.io/fyne/v2/app"      // 提供應用程式生命週期管理和視窗創建功能
	"fyne.io/fyne/v2/theme"    // 提供主題相關功能，用於自訂 UI 外觀樣式
	_ "embed"                  // Go 1.16+ 嵌入式檔案支援，用於嵌入字型資源
	"image/color"              // Go 標準庫，提供顏色定義和處理功能
	"log"                      // Go 標準庫，用於錯誤記錄
	"os"                       // Go 標準庫，用於作業系統介面
	"path/filepath"            // Go 標準庫，用於檔案路徑處理
	"mac-notebook-app/internal/models" // 本專案的資料模型套件
	"mac-notebook-app/internal/services" // 本專案的服務層套件
	"mac-notebook-app/internal/repositories" // 本專案的儲存庫層套件
	"mac-notebook-app/ui"      // 本專案的 UI 套件，包含主視窗和其他 UI 元件
)

// main 函數是應用程式的主要入口點
// 執行流程：
// 1. 建立 Fyne 應用程式實例
// 2. 設定應用程式基本屬性和主題
// 3. 建立主視窗實例
// 4. 顯示主視窗並啟動應用程式主迴圈
func main() {
	// 建立新的 Fyne 應用程式實例，使用唯一 ID
	// 這將初始化 GUI 框架並準備建立視窗
	myApp := app.NewWithID("com.macnotebook.app")
	
	// 設定應用程式的基本屬性
	myApp.SetIcon(nil) // 暫時不設定圖示，後續會添加自訂圖示

	// 設定應用程式主題為支援中日韓字型的深色主題
	myApp.Settings().SetTheme(&cjkTheme{base: theme.DarkTheme()})

	// 載入應用程式設定
	settings, err := models.LoadDefault()
	if err != nil {
		// 如果載入設定失敗，使用預設設定
		settings = models.NewDefaultSettings()
	}

	// 建立必要的服務實例
	// 1. 建立檔案儲存庫
	baseDir := settings.DefaultSaveLocation
	if baseDir == "" {
		// 如果沒有設定預設位置，使用使用者主目錄下的 Documents/NotebookApp
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Printf("無法取得使用者主目錄: %v", err)
			baseDir = "."
		} else {
			baseDir = filepath.Join(homeDir, "Documents", "NotebookApp")
		}
		
		// 確保目錄存在
		if err := os.MkdirAll(baseDir, 0755); err != nil {
			log.Printf("無法建立基礎目錄 %s: %v", baseDir, err)
			baseDir = "."
		}
	}
	
	fileRepo, err := repositories.NewLocalFileRepository(baseDir)
	if err != nil {
		log.Fatalf("建立檔案儲存庫失敗: %v", err)
	}

	// 2. 建立檔案管理服務
	fileManagerService, err := services.NewLocalFileManagerService(fileRepo, baseDir)
	if err != nil {
		log.Fatalf("建立檔案管理服務失敗: %v", err)
	}

	// 3. 建立各種服務
	encryptionService := services.NewEncryptionService()
	passwordService := services.NewPasswordService()
	biometricService := services.NewBiometricService()
	performanceService := services.NewPerformanceService(nil)
	smartEditingService := services.NewSmartEditingService()
	
	// 4. 建立編輯器服務
	editorService := services.NewEditorService(fileRepo, encryptionService, passwordService, biometricService, performanceService, smartEditingService)

	// 建立主視窗實例
	// 使用新的 MainWindow 結構，包含完整的 UI 佈局和服務整合
	mainWindow := ui.NewMainWindow(myApp, settings, editorService, fileManagerService)

	// 顯示主視窗並啟動應用程式的主事件迴圈
	// 這個函數會阻塞直到使用者關閉應用程式
	mainWindow.ShowAndRun()
}



//go:embed assets/font/GoogleSansCode-Regular.ttf
var fontRegular []byte

//go:embed assets/font/GoogleSansCode-Bold.ttf
var fontBold []byte

// 若沒有可留空，程式會回退 Regular
//go:embed assets/font/GoogleSansCode-Italic.ttf
var fontItalic []byte

//go:embed assets/font/GoogleSansCode-Regular.ttf
var fontMono []byte

type cjkTheme struct{ base fyne.Theme }

func (t cjkTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	return t.base.Color(n, v)
}

func (t cjkTheme) Icon(n fyne.ThemeIconName) fyne.Resource {
	return t.base.Icon(n)
}
func (t cjkTheme) Size(n fyne.ThemeSizeName) float32 {
	return t.base.Size(n)
}

func (t cjkTheme) Font(s fyne.TextStyle) fyne.Resource {
	// 依樣式回傳對應字型，沒有就回退 Regular
	switch {
	case s.Monospace && len(fontMono) > 0:
		return fyne.NewStaticResource("GoogleSansCode-Regular.ttf", fontMono)
	case s.Bold && len(fontBold) > 0:
		return fyne.NewStaticResource("GoogleSansCode-Bold.ttf", fontBold)
	case s.Italic && len(fontItalic) > 0:
		return fyne.NewStaticResource("GoogleSansCode-Italic.ttf", fontItalic)
	default:
		return fyne.NewStaticResource("GoogleSansCode-Regular.ttf", fontRegular)
	}
}