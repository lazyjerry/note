// Package ui 包含使用者介面相關的元件和視窗管理
// 使用 Fyne 框架建立跨平台的圖形使用者介面
package ui

import (
	"fyne.io/fyne/v2"          // Fyne GUI 框架核心套件
	"fyne.io/fyne/v2/container" // Fyne 容器佈局套件
	"fyne.io/fyne/v2/widget"   // Fyne UI 元件套件
)

// MainWindow 代表應用程式的主視窗
// 包含所有主要的 UI 元件，如選單欄、工具欄、內容區域和狀態欄
type MainWindow struct {
	window    fyne.Window      // 主視窗實例
	content   *fyne.Container  // 主要內容容器
	menuBar   *fyne.MainMenu   // 選單欄
	toolBar   *widget.Toolbar  // 工具欄
	statusBar *widget.Label    // 狀態欄
}

// NewMainWindow 建立新的主視窗實例
// 參數：app（Fyne 應用程式實例）
// 回傳：指向新建立的 MainWindow 的指標
//
// 執行流程：
// 1. 建立新的視窗並設定標題
// 2. 設定視窗的初始大小
// 3. 建立 MainWindow 結構體實例
// 4. 初始化使用者介面元件
// 5. 回傳主視窗實例
func NewMainWindow(app fyne.App) *MainWindow {
	// 建立新視窗並設定標題
	window := app.NewWindow("Mac Notebook App")
	
	// 設定視窗初始大小為 1200x800 像素
	// 這個大小適合筆記編輯和檔案管理的雙面板佈局
	window.Resize(fyne.NewSize(1200, 800))
	
	// 建立 MainWindow 實例
	mw := &MainWindow{
		window: window, // 設定視窗實例
	}
	
	// 初始化使用者介面
	mw.setupUI()
	
	return mw
}

// setupUI 初始化使用者介面元件
// 這個方法負責建立和配置主視窗的所有 UI 元件
//
// 執行流程：
// 1. 建立暫時的佔位元件
// 2. 設定主要內容容器
// 3. 將內容設定到視窗中
//
// 注意：目前只有佔位內容，完整的 UI 元件將在後續任務中實作
func (mw *MainWindow) setupUI() {
	// 建立暫時的佔位標籤
	// 這個標籤將在後續任務中被完整的 UI 元件替換
	placeholder := widget.NewLabel("UI components will be implemented in subsequent tasks")
	
	// 建立垂直容器作為主要內容區域
	mw.content = container.NewVBox(placeholder)
	
	// 將內容設定到主視窗
	mw.window.SetContent(mw.content)
}

// Show 顯示主視窗
// 這個方法會顯示視窗但不會阻塞程式執行
// 適用於需要在背景繼續執行其他操作的情況
func (mw *MainWindow) Show() {
	mw.window.Show()
}

// ShowAndRun 顯示主視窗並啟動應用程式主迴圈
// 這個方法會顯示視窗並阻塞程式執行，直到使用者關閉應用程式
//
// 執行流程：
// 1. 顯示主視窗
// 2. 啟動 Fyne 的事件迴圈
// 3. 處理使用者互動事件
// 4. 當視窗關閉時結束應用程式
func (mw *MainWindow) ShowAndRun() {
	mw.window.ShowAndRun()
}