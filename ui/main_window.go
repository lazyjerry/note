// Package ui 包含使用者介面相關的元件和視窗管理
// 使用 Fyne 框架建立跨平台的圖形使用者介面
package ui

import (
	"fmt"                      // Go 標準庫，用於格式化字串
	"fyne.io/fyne/v2"          // Fyne GUI 框架核心套件
	"fyne.io/fyne/v2/container" // Fyne 容器佈局套件
	"fyne.io/fyne/v2/widget"   // Fyne UI 元件套件
	"fyne.io/fyne/v2/theme"    // Fyne 主題套件
)

// MainWindow 代表應用程式的主視窗
// 包含所有主要的 UI 元件，如選單欄、工具欄、內容區域和狀態欄
// 採用標準的桌面應用程式佈局結構
type MainWindow struct {
	window       fyne.Window      // 主視窗實例
	content      *fyne.Container  // 主要內容容器
	menuBar      *fyne.MainMenu   // 選單欄
	toolBar      *widget.Toolbar  // 工具欄
	statusBar    *fyne.Container  // 狀態欄容器
	
	// 狀態欄元件
	saveStatus   *widget.Label    // 保存狀態指示器
	encStatus    *widget.Label    // 加密狀態指示器
	wordCount    *widget.Label    // 字數統計顯示
	
	// 主要內容區域
	leftPanel    *fyne.Container  // 左側面板（檔案樹和筆記列表）
	rightPanel   *fyne.Container  // 右側面板（編輯器和預覽）
	mainSplit    *container.Split // 主要分割容器
}

// NewMainWindow 建立新的主視窗實例
// 參數：app（Fyne 應用程式實例）
// 回傳：指向新建立的 MainWindow 的指標
//
// 執行流程：
// 1. 建立新的視窗並設定標題和基本屬性
// 2. 設定視窗的初始大小和位置
// 3. 建立 MainWindow 結構體實例
// 4. 初始化所有 UI 元件（選單、工具欄、狀態欄）
// 5. 設定主要佈局結構
// 6. 回傳完整配置的主視窗實例
func NewMainWindow(app fyne.App) *MainWindow {
	// 建立新視窗並設定標題
	window := app.NewWindow("Mac Notebook App - 安全筆記編輯器")
	
	// 設定視窗初始大小為 1200x800 像素
	// 這個大小適合筆記編輯和檔案管理的雙面板佈局
	window.Resize(fyne.NewSize(1200, 800))
	
	// 設定視窗居中顯示
	window.CenterOnScreen()
	
	// 建立 MainWindow 實例
	mw := &MainWindow{
		window: window, // 設定視窗實例
	}
	
	// 初始化使用者介面元件
	mw.setupUI()
	
	// 設定視窗關閉時的清理工作
	window.SetCloseIntercept(func() {
		// 在這裡可以添加保存未儲存的工作等清理邏輯
		window.Close()
	})
	
	return mw
}

// setupUI 初始化使用者介面元件
// 這個方法負責建立和配置主視窗的所有 UI 元件
//
// 執行流程：
// 1. 建立選單欄和所有選單項目
// 2. 建立工具欄和常用功能按鈕
// 3. 建立狀態欄和狀態指示器
// 4. 建立主要內容區域的佈局結構
// 5. 組合所有元件到主視窗中
func (mw *MainWindow) setupUI() {
	// 建立選單欄
	mw.createMenuBar()
	
	// 建立工具欄
	mw.createToolBar()
	
	// 建立狀態欄
	mw.createStatusBar()
	
	// 建立主要內容區域
	mw.createContentArea()
	
	// 組合所有元件到主視窗
	mw.assembleMainLayout()
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

// createMenuBar 建立應用程式的選單欄
// 包含檔案、編輯、檢視等主要選單項目
//
// 執行流程：
// 1. 建立檔案選單（新增、開啟、儲存、設定等）
// 2. 建立編輯選單（復原、重做、尋找等）
// 3. 建立檢視選單（主題、預覽等）
// 4. 組合所有選單到主選單欄
func (mw *MainWindow) createMenuBar() {
	// 建立檔案選單項目
	fileMenu := fyne.NewMenu("檔案",
		fyne.NewMenuItem("新增筆記", func() {
			// TODO: 實作新增筆記功能
			fmt.Println("新增筆記功能將在後續任務中實作")
		}),
		fyne.NewMenuItem("開啟檔案", func() {
			// TODO: 實作開啟檔案功能
			fmt.Println("開啟檔案功能將在後續任務中實作")
		}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("儲存", func() {
			// TODO: 實作儲存功能
			fmt.Println("儲存功能將在後續任務中實作")
		}),
		fyne.NewMenuItem("另存新檔", func() {
			// TODO: 實作另存新檔功能
			fmt.Println("另存新檔功能將在後續任務中實作")
		}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("設定", func() {
			// TODO: 實作設定對話框
			fmt.Println("設定對話框將在後續任務中實作")
		}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("結束", func() {
			mw.window.Close()
		}),
	)
	
	// 建立編輯選單項目
	editMenu := fyne.NewMenu("編輯",
		fyne.NewMenuItem("復原", func() {
			// TODO: 實作復原功能
			fmt.Println("復原功能將在後續任務中實作")
		}),
		fyne.NewMenuItem("重做", func() {
			// TODO: 實作重做功能
			fmt.Println("重做功能將在後續任務中實作")
		}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("尋找", func() {
			// TODO: 實作尋找功能
			fmt.Println("尋找功能將在後續任務中實作")
		}),
		fyne.NewMenuItem("取代", func() {
			// TODO: 實作取代功能
			fmt.Println("取代功能將在後續任務中實作")
		}),
	)
	
	// 建立檢視選單項目
	viewMenu := fyne.NewMenu("檢視",
		fyne.NewMenuItem("切換主題", func() {
			// TODO: 實作主題切換功能
			fmt.Println("主題切換功能將在後續任務中實作")
		}),
		fyne.NewMenuItem("切換預覽", func() {
			// TODO: 實作預覽切換功能
			fmt.Println("預覽切換功能將在後續任務中實作")
		}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("全螢幕", func() {
			mw.window.SetFullScreen(!mw.window.FullScreen())
		}),
	)
	
	// 建立說明選單項目
	helpMenu := fyne.NewMenu("說明",
		fyne.NewMenuItem("關於", func() {
			// TODO: 實作關於對話框
			fmt.Println("關於對話框將在後續任務中實作")
		}),
	)
	
	// 組合主選單欄
	mw.menuBar = fyne.NewMainMenu(fileMenu, editMenu, viewMenu, helpMenu)
	mw.window.SetMainMenu(mw.menuBar)
}

// createToolBar 建立應用程式的工具欄
// 包含常用功能的快速存取按鈕
//
// 執行流程：
// 1. 建立新增筆記按鈕
// 2. 建立儲存按鈕
// 3. 建立加密切換按鈕
// 4. 建立預覽切換按鈕
// 5. 組合所有按鈕到工具欄
func (mw *MainWindow) createToolBar() {
	mw.toolBar = widget.NewToolbar(
		// 新增筆記按鈕
		widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {
			// TODO: 實作新增筆記功能
			fmt.Println("新增筆記功能將在後續任務中實作")
		}),
		
		// 開啟檔案按鈕
		widget.NewToolbarAction(theme.FolderOpenIcon(), func() {
			// TODO: 實作開啟檔案功能
			fmt.Println("開啟檔案功能將在後續任務中實作")
		}),
		
		// 儲存按鈕
		widget.NewToolbarAction(theme.DocumentSaveIcon(), func() {
			// TODO: 實作儲存功能
			fmt.Println("儲存功能將在後續任務中實作")
		}),
		
		// 分隔線
		widget.NewToolbarSeparator(),
		
		// 加密切換按鈕
		widget.NewToolbarAction(theme.VisibilityOffIcon(), func() {
			// TODO: 實作加密切換功能
			fmt.Println("加密切換功能將在後續任務中實作")
		}),
		
		// 預覽切換按鈕
		widget.NewToolbarAction(theme.ViewRefreshIcon(), func() {
			// TODO: 實作預覽切換功能
			fmt.Println("預覽切換功能將在後續任務中實作")
		}),
		
		// 分隔線
		widget.NewToolbarSeparator(),
		
		// 設定按鈕
		widget.NewToolbarAction(theme.SettingsIcon(), func() {
			// TODO: 實作設定對話框
			fmt.Println("設定對話框將在後續任務中實作")
		}),
	)
}

// createStatusBar 建立應用程式的狀態欄
// 顯示保存狀態、加密狀態和字數統計等資訊
//
// 執行流程：
// 1. 建立保存狀態指示器
// 2. 建立加密狀態指示器
// 3. 建立字數統計顯示
// 4. 使用水平佈局組合狀態欄元件
func (mw *MainWindow) createStatusBar() {
	// 建立保存狀態指示器
	mw.saveStatus = widget.NewLabel("已儲存")
	mw.saveStatus.TextStyle = fyne.TextStyle{Italic: true}
	
	// 建立加密狀態指示器
	mw.encStatus = widget.NewLabel("未加密")
	mw.encStatus.TextStyle = fyne.TextStyle{Italic: true}
	
	// 建立字數統計顯示
	mw.wordCount = widget.NewLabel("字數: 0")
	mw.wordCount.TextStyle = fyne.TextStyle{Italic: true}
	
	// 建立分隔線
	separator1 := widget.NewSeparator()
	separator2 := widget.NewSeparator()
	
	// 使用水平佈局組合狀態欄
	// 左側顯示保存和加密狀態，右側顯示字數統計
	mw.statusBar = container.NewHBox(
		mw.saveStatus,
		separator1,
		mw.encStatus,
		widget.NewLabel(""), // 彈性空間
		separator2,
		mw.wordCount,
	)
}

// createContentArea 建立主要內容區域
// 包含左側面板（檔案樹和筆記列表）和右側面板（編輯器和預覽）
//
// 執行流程：
// 1. 建立左側面板包含檔案樹
// 2. 建立右側面板的佔位內容
// 3. 使用水平分割容器組合左右面板
func (mw *MainWindow) createContentArea() {
	// 建立左側面板包含檔案樹
	mw.createLeftPanel()
	
	// 建立右側面板佔位內容
	// 這將在後續任務中被編輯器和預覽面板替換
	rightPlaceholder := widget.NewLabel("Markdown 編輯器和預覽面板\n將在後續任務中實作")
	rightPlaceholder.Alignment = fyne.TextAlignCenter
	mw.rightPanel = container.NewVBox(rightPlaceholder)
	
	// 使用水平分割容器組合左右面板
	// 左側面板佔 30%，右側面板佔 70%
	mw.mainSplit = container.NewHSplit(mw.leftPanel, mw.rightPanel)
	mw.mainSplit.Offset = 0.3 // 設定分割比例
}

// createLeftPanel 建立左側面板
// 包含檔案樹和相關控制元件
//
// 執行流程：
// 1. 建立檔案樹元件
// 2. 設定檔案樹的回調函數
// 3. 建立面板標題和控制按鈕
// 4. 組合所有元件到左側面板
func (mw *MainWindow) createLeftPanel() {
	// 建立面板標題
	titleLabel := widget.NewLabel("檔案瀏覽器")
	titleLabel.TextStyle = fyne.TextStyle{Bold: true}
	
	// 建立檔案樹元件（暫時使用當前目錄作為根目錄）
	// 在實際應用中，這應該從設定中讀取或讓使用者選擇
	fileTree := mw.createFileTree(".")
	
	// 建立控制按鈕
	refreshButton := widget.NewButton("重新整理", func() {
		// TODO: 實作檔案樹重新整理功能
		fmt.Println("檔案樹重新整理功能將在後續實作")
	})
	
	newFolderButton := widget.NewButton("新增資料夾", func() {
		// TODO: 實作新增資料夾功能
		fmt.Println("新增資料夾功能將在後續實作")
	})
	
	// 建立按鈕容器
	buttonContainer := container.NewHBox(refreshButton, newFolderButton)
	
	// 組合左側面板
	mw.leftPanel = container.NewVBox(
		titleLabel,
		widget.NewSeparator(),
		fileTree,
		widget.NewSeparator(),
		buttonContainer,
	)
}

// createFileTree 建立檔案樹元件
// 參數：rootPath（根目錄路徑）
// 回傳：檔案樹元件
//
// 執行流程：
// 1. 建立模擬的檔案管理服務（暫時使用）
// 2. 建立檔案樹元件
// 3. 設定檔案選擇和目錄開啟回調
// 4. 回傳檔案樹元件
func (mw *MainWindow) createFileTree(rootPath string) fyne.CanvasObject {
	// 暫時建立一個簡單的檔案樹佔位元件
	// 在後續任務中將整合真實的檔案管理服務
	treeLabel := widget.NewLabel("檔案樹元件\n（整合檔案管理服務）")
	treeLabel.Alignment = fyne.TextAlignCenter
	
	// 建立一個簡單的樹狀結構示例
	tree := widget.NewTree(
		func(uid widget.TreeNodeID) []widget.TreeNodeID {
			// 根節點
			if uid == "" {
				return []widget.TreeNodeID{"root"}
			}
			// 根節點的子項目
			if uid == "root" {
				return []widget.TreeNodeID{"notes", "docs", "readme.md"}
			}
			// notes 目錄的子項目
			if uid == "notes" {
				return []widget.TreeNodeID{"work", "personal"}
			}
			return []widget.TreeNodeID{}
		},
		func(uid widget.TreeNodeID) bool {
			// 目錄節點
			return uid == "root" || uid == "notes"
		},
		func(branch bool) fyne.CanvasObject {
			// 建立節點 UI
			var icon *widget.Icon
			if branch {
				icon = widget.NewIcon(theme.FolderIcon())
			} else {
				icon = widget.NewIcon(theme.DocumentIcon())
			}
			label := widget.NewLabel("")
			return container.NewHBox(icon, label)
		},
		func(uid widget.TreeNodeID, branch bool, obj fyne.CanvasObject) {
			// 更新節點 UI
			hbox := obj.(*fyne.Container)
			if len(hbox.Objects) >= 2 {
				label := hbox.Objects[1].(*widget.Label)
				switch uid {
				case "root":
					label.SetText("專案根目錄")
				case "notes":
					label.SetText("筆記")
				case "docs":
					label.SetText("文件")
				case "work":
					label.SetText("工作")
				case "personal":
					label.SetText("個人")
				case "readme.md":
					label.SetText("README.md")
				default:
					label.SetText(string(uid))
				}
			}
		},
	)
	
	// 設定節點選擇回調
	tree.OnSelected = func(uid widget.TreeNodeID) {
		fmt.Printf("選擇了節點: %s\n", uid)
		// TODO: 在後續任務中實作檔案開啟功能
	}
	
	return tree
}

// assembleMainLayout 組合主視窗的完整佈局
// 將工具欄、內容區域和狀態欄組合成完整的視窗佈局
//
// 執行流程：
// 1. 建立垂直容器作為主要佈局
// 2. 依序添加工具欄、內容區域和狀態欄
// 3. 將完整佈局設定到主視窗
func (mw *MainWindow) assembleMainLayout() {
	// 建立主要內容容器，使用垂直佈局
	mw.content = container.NewVBox(
		mw.toolBar,                    // 工具欄在頂部
		mw.mainSplit,                  // 主要內容區域在中間
		widget.NewSeparator(),         // 分隔線
		mw.statusBar,                  // 狀態欄在底部
	)
	
	// 將完整佈局設定到主視窗
	mw.window.SetContent(mw.content)
}

// UpdateSaveStatus 更新保存狀態顯示
// 參數：status（保存狀態文字）
//
// 執行流程：
// 1. 更新保存狀態標籤的文字
// 2. 根據狀態設定適當的顏色或樣式
func (mw *MainWindow) UpdateSaveStatus(status string) {
	if mw.saveStatus != nil {
		mw.saveStatus.SetText(status)
		mw.saveStatus.Refresh()
	}
}

// UpdateEncryptionStatus 更新加密狀態顯示
// 參數：isEncrypted（是否已加密）, encType（加密類型）
//
// 執行流程：
// 1. 根據加密狀態設定適當的顯示文字
// 2. 更新加密狀態標籤
func (mw *MainWindow) UpdateEncryptionStatus(isEncrypted bool, encType string) {
	if mw.encStatus != nil {
		if isEncrypted {
			mw.encStatus.SetText(fmt.Sprintf("已加密 (%s)", encType))
		} else {
			mw.encStatus.SetText("未加密")
		}
		mw.encStatus.Refresh()
	}
}

// UpdateWordCount 更新字數統計顯示
// 參數：count（字數）
//
// 執行流程：
// 1. 格式化字數顯示文字
// 2. 更新字數統計標籤
func (mw *MainWindow) UpdateWordCount(count int) {
	if mw.wordCount != nil {
		mw.wordCount.SetText(fmt.Sprintf("字數: %d", count))
		mw.wordCount.Refresh()
	}
}

// GetWindow 取得主視窗實例
// 回傳：主視窗的 fyne.Window 介面
// 用於其他元件需要存取視窗功能時使用
func (mw *MainWindow) GetWindow() fyne.Window {
	return mw.window
}