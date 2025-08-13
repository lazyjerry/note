// Package main 是 Mac Notebook App 的主要入口點
// 這是一個安全的 Markdown 筆記應用程式，支援密碼和生物識別加密
package main

import (
	"fyne.io/fyne/v2"          // Fyne GUI 框架核心套件
	"fyne.io/fyne/v2/app"      // Fyne 應用程式建立套件
	"fyne.io/fyne/v2/container" // Fyne 容器佈局套件
	"fyne.io/fyne/v2/widget"   // Fyne UI 元件套件
	"fyne.io/fyne/v2/theme"
	_ "embed"
	"image/color"
)

// main 函數是應用程式的主要入口點
// 執行流程：
// 1. 建立 Fyne 應用程式實例
// 2. 設定應用程式基本屬性
// 3. 建立主視窗並設定大小
// 4. 建立基本的 UI 佈局
// 5. 顯示視窗並啟動應用程式主迴圈
func main() {
	// 建立新的 Fyne 應用程式實例
	// 這將初始化 GUI 框架並準備建立視窗
	myApp := app.New()
	
	// 設定應用程式的基本屬性
	// 在 Fyne v2 中，應用程式 ID 通過不同的方式設定
	myApp.SetIcon(nil) // 暫時不設定圖示，後續會添加

	myApp.Settings().SetTheme(&cjkTheme{base: theme.LightTheme()})

	// 建立主視窗
	// 這將建立應用程式的主要使用者介面視窗
	myWindow := myApp.NewWindow("Mac Notebook App")
	
	// 設定視窗的初始大小為 1200x800 像素
	// 這個大小適合筆記編輯和檔案管理的雙面板佈局
	myWindow.Resize(fyne.NewSize(1200, 800))
	
	// 設定視窗居中顯示
	myWindow.CenterOnScreen()

	// 建立基本的 UI 佈局
	// 這是一個暫時的佈局，展示應用程式的基本結構
	content := createBasicLayout()
	myWindow.SetContent(content)

	// 顯示視窗並啟動應用程式的主事件迴圈
	// 這個函數會阻塞直到使用者關閉應用程式
	myWindow.ShowAndRun()
}

// createBasicLayout 建立基本的應用程式佈局
// 回傳：包含基本 UI 元素的容器
//
// 執行流程：
// 1. 建立標題標籤
// 2. 建立狀態資訊標籤
// 3. 建立功能說明文字
// 4. 使用垂直佈局組合所有元素
func createBasicLayout() fyne.CanvasObject {
	// 建立應用程式標題
	title := widget.NewLabel("Mac Notebook App")
	title.TextStyle = fyne.TextStyle{Bold: true}
	
	// 建立版本資訊
	version := widget.NewLabel("版本 0.5.0 - 編輯器服務完成")
	
	// 建立功能狀態說明
	status := widget.NewLabel("✅ 已完成功能：")
	
	// 建立功能清單
	features := widget.NewRichTextFromMarkdown(`
**已完成的核心功能：**

• 資料模型和驗證 (Note, Settings)
• 檔案系統操作 (FileRepository, FileManagerService)  
• 加密功能 (AES-256, ChaCha20, 密碼驗證, 生物識別)
• 編輯器服務 (Markdown 解析, 即時預覽, 加密整合)
• 完整的單元測試覆蓋

**🚧 進行中：** 準備實作 UI 介面

**📋 下一步：** 實作自動保存系統和 Fyne UI 元件`)
	
	// 建立開發資訊
	devInfo := widget.NewLabel("所有後端服務已完成，準備開始 UI 開發階段")
	devInfo.TextStyle = fyne.TextStyle{Italic: true}
	
	// 使用垂直佈局組合所有元素
	// 添加適當的間距讓介面更美觀
	content := container.NewVBox(
		widget.NewSeparator(),
		title,
		widget.NewSeparator(),
		version,
		widget.NewSeparator(),
		status,
		features,
		widget.NewSeparator(),
		devInfo,
		widget.NewSeparator(),
	)
	
	return content
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