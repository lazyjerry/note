// Package main 是 Mac Notebook App 的主要入口點
// 這是一個安全的 Markdown 筆記應用程式，支援密碼和生物識別加密
package main

import (
	"fyne.io/fyne/v2"        // Fyne GUI 框架核心套件
	"fyne.io/fyne/v2/app"    // Fyne 應用程式建立套件
	"fyne.io/fyne/v2/widget" // Fyne UI 元件套件
)

// main 函數是應用程式的主要入口點
// 執行流程：
// 1. 建立 Fyne 應用程式實例
// 2. 設定應用程式元資料（ID 和名稱）
// 3. 建立主視窗並設定大小
// 4. 設定暫時的佔位內容
// 5. 顯示視窗並啟動應用程式主迴圈
func main() {
	// 建立新的 Fyne 應用程式實例
	// 這將初始化 GUI 框架並準備建立視窗
	myApp := app.New()
	
	// 設定應用程式的元資料
	// ID: 用於系統識別應用程式的唯一標識符
	// Name: 顯示在系統中的應用程式名稱
	myApp.SetMetadata(&app.Metadata{
		ID:   "com.notebook.mac-notebook-app",
		Name: "Mac Notebook App",
	})

	// 建立主視窗
	// 這將建立應用程式的主要使用者介面視窗
	myWindow := myApp.NewWindow("Mac Notebook App")
	
	// 設定視窗的初始大小為 1200x800 像素
	// 這個大小適合筆記編輯和檔案管理的雙面板佈局
	myWindow.Resize(fyne.NewSize(1200, 800))

	// 建立暫時的佔位內容
	// 這個標籤將在後續任務中被完整的 UI 元件替換
	content := widget.NewLabel("Mac Notebook App - Ready for implementation")
	myWindow.SetContent(content)

	// 顯示視窗並啟動應用程式的主事件迴圈
	// 這個函數會阻塞直到使用者關閉應用程式
	myWindow.ShowAndRun()
}