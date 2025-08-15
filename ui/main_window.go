// Package ui 包含使用者介面相關的元件和視窗管理
// 使用 Fyne 框架建立跨平台的圖形使用者介面
package ui

import (
	"fmt"                      // Go 標準庫，用於格式化字串
	"path/filepath"            // 檔案路徑處理
	"strings"                  // 字串處理
	"fyne.io/fyne/v2"          // Fyne GUI 框架核心套件
	"fyne.io/fyne/v2/container" // Fyne 容器佈局套件
	"fyne.io/fyne/v2/widget"   // Fyne UI 元件套件
	"fyne.io/fyne/v2/theme"    // Fyne 主題套件
	"fyne.io/fyne/v2/dialog"   // Fyne 對話框套件

	"mac-notebook-app/internal/models"
	"mac-notebook-app/internal/services"
)

// FileTree 代表檔案樹 UI 元件（暫時實作）
// 在 Task 12.2 中將實作完整的檔案管理整合
type FileTree struct {
	container *fyne.Container
}

// GetContainer 取得檔案樹的容器
func (ft *FileTree) GetContainer() *fyne.Container {
	return ft.container
}

// Refresh 重新整理檔案樹
func (ft *FileTree) Refresh() error {
	// 暫時實作，在 Task 12.2 中將實作完整功能
	return nil
}

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
	
	// UI 元件
	fileTree       *FileTree        // 檔案樹元件（舊版，保留相容性）
	fileTreeWidget *FileTreeWidget  // 新的檔案樹元件
	editor         *MarkdownEditor  // Markdown 編輯器元件

	// 服務和設定
	app              fyne.App                         // Fyne 應用程式實例
	settings         *models.Settings                 // 應用程式設定
	themeService     *services.ThemeService           // 主題管理服務
	editorService    services.EditorService           // 編輯器服務
	fileManagerService services.FileManagerService   // 檔案管理服務
}

// NewMainWindow 建立新的主視窗實例
// 參數：
//   - app: Fyne 應用程式實例
//   - settings: 應用程式設定
//   - editorService: 編輯器服務實例
//   - fileManagerService: 檔案管理服務實例
// 回傳：指向新建立的 MainWindow 的指標
//
// 執行流程：
// 1. 建立新的視窗並設定標題和基本屬性
// 2. 設定視窗的初始大小和位置
// 3. 建立 MainWindow 結構體實例
// 4. 初始化主題服務和業務服務
// 5. 初始化所有 UI 元件（選單、工具欄、狀態欄）
// 6. 設定主要佈局結構和服務整合
// 7. 回傳完整配置的主視窗實例
func NewMainWindow(app fyne.App, settings *models.Settings, editorService services.EditorService, fileManagerService services.FileManagerService) *MainWindow {
	// 建立新視窗並設定標題
	window := app.NewWindow("Mac Notebook App - 安全筆記編輯器")
	
	// 設定視窗初始大小為 1200x800 像素
	// 這個大小適合筆記編輯和檔案管理的雙面板佈局
	window.Resize(fyne.NewSize(1200, 800))
	
	// 設定視窗居中顯示
	window.CenterOnScreen()
	
	// 建立 MainWindow 實例
	mw := &MainWindow{
		window:             window,             // 設定視窗實例
		app:                app,                // 設定應用程式實例
		settings:           settings,           // 設定應用程式設定
		editorService:      editorService,      // 設定編輯器服務
		fileManagerService: fileManagerService, // 設定檔案管理服務
	}

	// 初始化主題服務
	mw.themeService = services.NewThemeService(app, settings)
	
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
// 6. 設定主題監聽器
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

	// 設定主題監聽器
	mw.SetupThemeListener()
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
			mw.createNewNote()
		}),
		fyne.NewMenuItem("開啟檔案", func() {
			mw.openFile()
		}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("儲存", func() {
			mw.saveCurrentNote()
		}),
		fyne.NewMenuItem("另存新檔", func() {
			mw.saveAsNewFile()
		}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("設定", func() {
			mw.showSettingsDialog()
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
			mw.createNewNote()
		}),
		
		// 開啟檔案按鈕
		widget.NewToolbarAction(theme.FolderOpenIcon(), func() {
			mw.openFile()
		}),
		
		// 儲存按鈕
		widget.NewToolbarAction(theme.DocumentSaveIcon(), func() {
			mw.saveCurrentNote()
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
			mw.showSettingsDialog()
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
// 2. 建立右側面板包含 Markdown 編輯器
// 3. 整合編輯器服務到 UI 元件
// 4. 使用水平分割容器組合左右面板
func (mw *MainWindow) createContentArea() {
	// 建立左側面板包含檔案樹
	mw.createLeftPanel()
	
	// 建立右側面板包含 Markdown 編輯器
	mw.createRightPanel()
	
	// 使用水平分割容器組合左右面板
	// 左側面板佔 30%，右側面板佔 70%
	mw.mainSplit = container.NewHSplit(mw.leftPanel, mw.rightPanel)
	mw.mainSplit.Offset = 0.3 // 設定分割比例
}

// createLeftPanel 建立左側面板
// 包含檔案樹和相關控制元件，整合檔案管理服務
//
// 執行流程：
// 1. 建立檔案樹元件並整合檔案管理服務
// 2. 設定檔案樹的回調函數和事件處理
// 3. 建立面板標題和控制按鈕
// 4. 組合所有元件到左側面板
func (mw *MainWindow) createLeftPanel() {
	// 建立面板標題
	titleLabel := widget.NewLabel("檔案瀏覽器")
	titleLabel.TextStyle = fyne.TextStyle{Bold: true}
	
	// 建立檔案樹元件並整合檔案管理服務
	// 使用設定中的預設保存位置或當前目錄
	rootPath := mw.settings.DefaultSaveLocation
	if rootPath == "" {
		rootPath = "."
	}
	
	// 建立真正的檔案樹元件並整合檔案管理服務
	mw.fileTreeWidget = NewFileTreeWidget(mw.fileManagerService, rootPath)
	
	// 設定檔案樹的回調函數
	mw.setupFileTreeCallbacks()
	
	// 建立控制按鈕
	refreshButton := widget.NewButton("重新整理", func() {
		mw.refreshFileTree()
	})
	
	newFolderButton := widget.NewButton("新增資料夾", func() {
		mw.createNewFolder()
	})
	
	newFileButton := widget.NewButton("新增檔案", func() {
		mw.createNewFileInCurrentDir()
	})
	
	// 建立按鈕容器
	buttonContainer := container.NewHBox(refreshButton, newFolderButton, newFileButton)
	
	// 組合左側面板
	mw.leftPanel = container.NewVBox(
		titleLabel,
		widget.NewSeparator(),
		mw.fileTreeWidget,
		widget.NewSeparator(),
		buttonContainer,
	)
}

// createRightPanel 建立右側面板
// 包含 Markdown 編輯器和預覽面板，整合編輯器服務
//
// 執行流程：
// 1. 建立 Markdown 編輯器元件並整合編輯器服務
// 2. 設定編輯器的回調函數和事件處理
// 3. 建立預覽面板（如果需要）
// 4. 組合編輯器和預覽面板到右側面板
func (mw *MainWindow) createRightPanel() {
	// 建立 Markdown 編輯器元件並整合編輯器服務
	mw.editor = NewMarkdownEditor(mw.editorService)
	
	// 設定編輯器的回調函數
	mw.setupEditorCallbacks()
	
	// 建立右側面板容器
	mw.rightPanel = container.NewVBox(
		mw.editor.GetContainer(),
	)
}

// setupEditorCallbacks 設定編輯器的回調函數
// 整合編輯器事件到主視窗的狀態管理
//
// 執行流程：
// 1. 設定內容變更回調，更新狀態欄
// 2. 設定保存請求回調，處理保存操作
// 3. 設定字數變更回調，更新字數統計
func (mw *MainWindow) setupEditorCallbacks() {
	// 設定內容變更回調
	mw.editor.SetOnContentChanged(func(content string) {
		// 更新保存狀態為未保存
		mw.UpdateSaveStatus("未保存")
		
		// 檢查是否為加密筆記並更新加密狀態
		if currentNote := mw.editor.GetCurrentNote(); currentNote != nil {
			mw.UpdateEncryptionStatus(currentNote.IsEncrypted, currentNote.EncryptionType)
		}
	})
	
	// 設定保存請求回調
	mw.editor.SetOnSaveRequested(func() {
		// 更新保存狀態
		mw.UpdateSaveStatus("已保存")
	})
	
	// 設定字數變更回調
	mw.editor.SetOnWordCountChanged(func(count int) {
		mw.UpdateWordCount(count)
	})
}

// setupFileTreeCallbacks 設定檔案樹的回調函數
// 整合檔案樹事件到主視窗的檔案管理功能
//
// 執行流程：
// 1. 設定檔案選擇回調，載入選擇的檔案到編輯器
// 2. 設定檔案開啟回調，處理檔案開啟操作
// 3. 設定目錄開啟回調，展開目錄結構
// 4. 設定檔案操作回調，處理各種檔案操作
// 5. 設定右鍵點擊回調，顯示操作選單
func (mw *MainWindow) setupFileTreeCallbacks() {
	if mw.fileTreeWidget == nil {
		return
	}
	
	// 設定檔案選擇回調
	mw.fileTreeWidget.SetOnFileSelect(func(filePath string) {
		// 載入選擇的檔案到編輯器
		mw.openFileFromPath(filePath)
	})
	
	// 設定檔案開啟回調
	mw.fileTreeWidget.SetOnFileOpen(func(filePath string) {
		// 開啟檔案到編輯器
		mw.openFileFromPath(filePath)
	})
	
	// 設定目錄開啟回調
	mw.fileTreeWidget.SetOnDirectoryOpen(func(dirPath string) {
		// 展開目錄（檔案樹會自動處理）
		fmt.Printf("開啟目錄: %s\n", dirPath)
	})
	
	// 設定檔案操作回調
	mw.fileTreeWidget.SetOnFileOperation(func(operation, filePath string) {
		mw.handleFileTreeOperation(operation, filePath)
	})
	
	// 設定右鍵點擊回調
	mw.fileTreeWidget.SetOnFileRightClick(func(filePath string, isDirectory bool) {
		mw.showFileContextMenu(filePath, isDirectory)
	})
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

// showSettingsDialog 顯示設定對話框
// 執行流程：
// 1. 建立設定對話框實例
// 2. 設定主題變更回調函數
// 3. 顯示對話框
func (mw *MainWindow) showSettingsDialog() {
	// 建立設定對話框
	settingsDialog := NewSettingsDialog(
		mw.window,
		mw.settings,
		func(newSettings *models.Settings) {
			// 當設定變更時的回調函數
			mw.onSettingsChanged(newSettings)
		},
	)
	
	// 顯示設定對話框
	settingsDialog.Show()
}

// onSettingsChanged 處理設定變更事件
// 參數：
//   - newSettings: 新的設定實例
//
// 執行流程：
// 1. 更新內部設定實例
// 2. 套用主題變更
// 3. 更新其他相關的 UI 元件
func (mw *MainWindow) onSettingsChanged(newSettings *models.Settings) {
	// 更新內部設定
	mw.settings = newSettings
	
	// 如果主題有變更，套用新主題
	if mw.themeService.GetCurrentTheme() != newSettings.Theme {
		mw.themeService.SetTheme(newSettings.Theme)
	}
	
	// 更新狀態欄顯示（如果需要）
	mw.updateUIFromSettings()
}

// updateUIFromSettings 根據設定更新 UI 元件
// 執行流程：
// 1. 更新加密狀態顯示
// 2. 更新其他相關的狀態指示器
func (mw *MainWindow) updateUIFromSettings() {
	// 更新加密狀態顯示
	if mw.encStatus != nil {
		encryptionText := fmt.Sprintf("加密: %s", mw.settings.DefaultEncryption)
		mw.encStatus.SetText(encryptionText)
	}
	
	// 更新保存狀態（顯示自動保存間隔）
	if mw.saveStatus != nil {
		saveText := fmt.Sprintf("自動保存: %d分鐘", mw.settings.AutoSaveInterval)
		mw.saveStatus.SetText(saveText)
	}
}

// OnThemeChanged 實作 ThemeListener 介面
// 參數：
//   - themeName: 新的主題名稱
//
// 執行流程：
// 1. 更新內部主題狀態
// 2. 重新整理 UI 元件以反映主題變更
func (mw *MainWindow) OnThemeChanged(themeName string) {
	// 主題變更時的處理邏輯
	// Fyne 會自動處理大部分的主題變更
	// 這裡可以添加自訂的主題變更邏輯
	
	// 更新狀態欄或其他需要手動更新的元件
	mw.updateUIFromSettings()
}

// SetupThemeListener 設定主題監聽器
// 執行流程：
// 1. 將主視窗註冊為主題變更監聽器
// 2. 初始化 UI 狀態
func (mw *MainWindow) SetupThemeListener() {
	// 註冊主題監聽器
	mw.themeService.AddThemeListener(mw)
	
	// 初始化 UI 狀態
	mw.updateUIFromSettings()
}

// createNewNote 建立新筆記
// 整合編輯器服務建立新筆記並載入到編輯器
//
// 執行流程：
// 1. 提示使用者輸入筆記標題
// 2. 使用編輯器服務建立新筆記
// 3. 載入新筆記到編輯器
// 4. 更新狀態顯示
func (mw *MainWindow) createNewNote() {
	// 建立標題輸入對話框
	titleEntry := widget.NewEntry()
	titleEntry.SetPlaceHolder("請輸入筆記標題...")
	
	// 建立對話框內容
	content := container.NewVBox(
		widget.NewLabel("新增筆記"),
		titleEntry,
	)
	
	// 建立對話框
	dialog := dialog.NewCustomConfirm("新增筆記", "建立", "取消", content, func(confirmed bool) {
		if confirmed {
			title := titleEntry.Text
			if title == "" {
				title = "未命名筆記"
			}
			
			// 使用編輯器建立新筆記
			err := mw.editor.CreateNewNote(title)
			if err != nil {
				dialog.ShowError(err, mw.window)
				return
			}
			
			// 更新狀態顯示
			mw.UpdateSaveStatus("新筆記")
			mw.UpdateEncryptionStatus(false, "")
		}
	}, mw.window)
	
	// 顯示對話框
	dialog.Show()
	
	// 設定焦點到標題輸入框
	titleEntry.FocusGained()
}

// openFile 開啟檔案
// 顯示檔案選擇對話框並開啟選擇的檔案
//
// 執行流程：
// 1. 顯示檔案選擇對話框
// 2. 使用編輯器服務開啟選擇的檔案
// 3. 載入檔案內容到編輯器
// 4. 更新狀態顯示
func (mw *MainWindow) openFile() {
	// 建立檔案開啟對話框
	fileDialog := NewFileOpenDialog(mw.window, func(filePath string) {
		mw.openFileFromPath(filePath)
	})
	
	// 顯示對話框
	fileDialog.Show()
}

// openFileFromPath 從指定路徑開啟檔案
// 參數：filePath（檔案路徑）
//
// 執行流程：
// 1. 使用編輯器服務開啟檔案
// 2. 處理加密檔案的密碼驗證
// 3. 載入檔案內容到編輯器
// 4. 更新狀態顯示
func (mw *MainWindow) openFileFromPath(filePath string) {
	// 使用編輯器服務開啟檔案
	note, err := mw.editorService.OpenNote(filePath)
	if err != nil {
		// 檢查是否為加密檔案需要密碼
		if strings.Contains(err.Error(), "需要密碼驗證") {
			mw.handleEncryptedFileOpen(filePath)
			return
		}
		
		dialog.ShowError(err, mw.window)
		return
	}
	
	// 載入筆記到編輯器
	mw.editor.LoadNote(note)
	
	// 更新狀態顯示
	mw.UpdateSaveStatus("已載入")
	mw.UpdateEncryptionStatus(note.IsEncrypted, note.EncryptionType)
	
	// 重新整理檔案樹以反映變更
	mw.refreshFileTree()
}

// handleEncryptedFileOpen 處理加密檔案的開啟
// 參數：filePath（加密檔案路徑）
//
// 執行流程：
// 1. 顯示密碼輸入對話框
// 2. 使用密碼解密檔案
// 3. 載入解密後的內容到編輯器
func (mw *MainWindow) handleEncryptedFileOpen(filePath string) {
	// 建立密碼輸入對話框
	passwordDialog := NewPasswordDialog(mw.window, "開啟加密檔案", func(password string) {
		// 先開啟檔案取得筆記 ID
		note, err := mw.editorService.OpenNote(filePath)
		if err != nil {
			// 如果還是失敗，嘗試使用密碼解密
			// 這裡需要特殊處理，因為我們需要筆記 ID 來解密
			dialog.ShowError(fmt.Errorf("無法開啟加密檔案: %w", err), mw.window)
			return
		}
		
		// 使用密碼解密內容
		decryptedContent, err := mw.editorService.DecryptWithPassword(note.ID, password)
		if err != nil {
			dialog.ShowError(fmt.Errorf("密碼錯誤或解密失敗: %w", err), mw.window)
			return
		}
		
		// 更新筆記內容
		note.Content = decryptedContent
		
		// 載入筆記到編輯器
		mw.editor.LoadNote(note)
		
		// 更新狀態顯示
		mw.UpdateSaveStatus("已載入")
		mw.UpdateEncryptionStatus(true, note.EncryptionType)
	})
	
	// 顯示密碼對話框
	passwordDialog.Show()
}

// saveCurrentNote 保存當前筆記
// 使用編輯器服務保存當前編輯的筆記
//
// 執行流程：
// 1. 檢查是否有當前筆記
// 2. 使用編輯器保存筆記
// 3. 處理保存結果和錯誤
// 4. 更新狀態顯示
func (mw *MainWindow) saveCurrentNote() {
	// 檢查編輯器是否可以保存
	if !mw.editor.CanSave() {
		mw.UpdateSaveStatus("無需保存")
		return
	}
	
	// 使用編輯器保存筆記
	err := mw.editor.SaveNote()
	if err != nil {
		dialog.ShowError(err, mw.window)
		mw.UpdateSaveStatus("保存失敗")
		return
	}
	
	// 更新狀態顯示
	mw.UpdateSaveStatus("已保存")
	
	// 重新整理檔案樹以反映變更
	mw.refreshFileTree()
}

// saveAsNewFile 另存新檔
// 顯示檔案保存對話框並將當前筆記保存為新檔案
//
// 執行流程：
// 1. 顯示檔案保存對話框
// 2. 設定新的檔案路徑
// 3. 保存筆記到新位置
// 4. 更新狀態顯示
func (mw *MainWindow) saveAsNewFile() {
	// 檢查是否有當前筆記
	currentNote := mw.editor.GetCurrentNote()
	if currentNote == nil {
		dialog.ShowInformation("提示", "沒有可保存的筆記", mw.window)
		return
	}
	
	// 建立檔案保存對話框
	fileDialog := NewFileSaveDialog(mw.window, func(filePath string) {
		// 更新筆記的檔案路徑
		currentNote.FilePath = filePath
		
		// 保存筆記
		err := mw.editorService.SaveNote(currentNote)
		if err != nil {
			dialog.ShowError(err, mw.window)
			mw.UpdateSaveStatus("保存失敗")
			return
		}
		
		// 更新狀態顯示
		mw.UpdateSaveStatus("已保存")
		
		// 重新整理檔案樹以反映變更
		mw.refreshFileTree()
	})
	
	// 設定預設檔案名稱
	if currentNote.Title != "" {
		fileDialog.SetFileName(currentNote.Title + ".md")
	}
	
	// 顯示對話框
	fileDialog.Show()
}

// refreshFileTree 重新整理檔案樹
// 重新載入檔案樹的內容以反映檔案系統的變更
//
// 執行流程：
// 1. 重新載入檔案樹內容
// 2. 更新 UI 顯示
// 3. 處理重新載入錯誤
func (mw *MainWindow) refreshFileTree() {
	// 使用新的檔案樹元件
	if mw.fileTreeWidget != nil {
		mw.fileTreeWidget.Refresh()
		return
	}
	
	// 保留舊版相容性
	if mw.fileTree != nil {
		err := mw.fileTree.Refresh()
		if err != nil {
			fmt.Printf("重新整理檔案樹失敗: %v\n", err)
		}
	}
}

// createNewFolder 建立新資料夾
// 顯示資料夾名稱輸入對話框並建立新資料夾
//
// 執行流程：
// 1. 提示使用者輸入資料夾名稱
// 2. 使用檔案管理服務建立資料夾
// 3. 重新整理檔案樹
// 4. 處理建立結果
func (mw *MainWindow) createNewFolder() {
	// 建立資料夾名稱輸入對話框
	folderEntry := widget.NewEntry()
	folderEntry.SetPlaceHolder("請輸入資料夾名稱...")
	
	// 建立對話框內容
	content := container.NewVBox(
		widget.NewLabel("新增資料夾"),
		folderEntry,
	)
	
	// 建立對話框
	dialog := dialog.NewCustomConfirm("新增資料夾", "建立", "取消", content, func(confirmed bool) {
		if confirmed {
			folderName := folderEntry.Text
			if folderName == "" {
				dialog.ShowError(fmt.Errorf("資料夾名稱不能為空"), mw.window)
				return
			}
			
			// 使用檔案管理服務建立資料夾
			err := mw.fileManagerService.CreateDirectory(folderName)
			if err != nil {
				dialog.ShowError(err, mw.window)
				return
			}
			
			// 重新整理檔案樹
			mw.refreshFileTree()
		}
	}, mw.window)
	
	// 顯示對話框
	dialog.Show()
	
	// 設定焦點到資料夾名稱輸入框
	folderEntry.FocusGained()
}

// handleFileOperation 處理檔案操作
// 參數：operation（操作類型）、path（檔案路徑）
//
// 執行流程：
// 1. 根據操作類型執行相應的檔案操作
// 2. 使用檔案管理服務執行操作
// 3. 更新 UI 狀態
// 4. 處理操作結果
func (mw *MainWindow) handleFileOperation(operation, path string) {
	switch operation {
	case "delete":
		mw.deleteFile(path)
	case "rename":
		mw.renameFile(path)
	case "copy":
		mw.copyFile(path)
	case "move":
		mw.moveFile(path)
	default:
		fmt.Printf("未知的檔案操作: %s\n", operation)
	}
}

// deleteFile 刪除檔案
// 參數：filePath（要刪除的檔案路徑）
//
// 執行流程：
// 1. 顯示確認對話框
// 2. 使用檔案管理服務刪除檔案
// 3. 重新整理檔案樹
// 4. 處理刪除結果
func (mw *MainWindow) deleteFile(filePath string) {
	// 顯示確認對話框
	dialog.ShowConfirm("確認刪除", 
		fmt.Sprintf("確定要刪除檔案 '%s' 嗎？", filepath.Base(filePath)), 
		func(confirmed bool) {
			if confirmed {
				// 使用檔案管理服務刪除檔案
				err := mw.fileManagerService.DeleteFile(filePath)
				if err != nil {
					dialog.ShowError(err, mw.window)
					return
				}
				
				// 重新整理檔案樹
				mw.refreshFileTree()
			}
		}, mw.window)
}

// renameFile 重新命名檔案
// 參數：filePath（要重新命名的檔案路徑）
//
// 執行流程：
// 1. 顯示新名稱輸入對話框
// 2. 使用檔案管理服務重新命名檔案
// 3. 重新整理檔案樹
// 4. 處理重新命名結果
func (mw *MainWindow) renameFile(filePath string) {
	// 取得當前檔案名稱
	currentName := filepath.Base(filePath)
	
	// 建立新名稱輸入對話框
	nameEntry := widget.NewEntry()
	nameEntry.SetText(currentName)
	
	// 建立對話框內容
	content := container.NewVBox(
		widget.NewLabel("重新命名"),
		nameEntry,
	)
	
	// 建立對話框
	dialog := dialog.NewCustomConfirm("重新命名", "確定", "取消", content, func(confirmed bool) {
		if confirmed {
			newName := nameEntry.Text
			if newName == "" || newName == currentName {
				return
			}
			
			// 建立新路徑
			newPath := filepath.Join(filepath.Dir(filePath), newName)
			
			// 使用檔案管理服務重新命名檔案
			err := mw.fileManagerService.RenameFile(filePath, newPath)
			if err != nil {
				dialog.ShowError(err, mw.window)
				return
			}
			
			// 重新整理檔案樹
			mw.refreshFileTree()
		}
	}, mw.window)
	
	// 顯示對話框
	dialog.Show()
	
	// 設定焦點到名稱輸入框
	nameEntry.FocusGained()
}

// copyFile 複製檔案
// 參數：filePath（要複製的檔案路徑）
//
// 執行流程：
// 1. 顯示目標路徑選擇對話框
// 2. 使用檔案管理服務複製檔案
// 3. 重新整理檔案樹
// 4. 處理複製結果
func (mw *MainWindow) copyFile(filePath string) {
	// 建立目標路徑輸入對話框
	targetEntry := widget.NewEntry()
	targetEntry.SetPlaceHolder("請輸入目標路徑...")
	
	// 建立對話框內容
	content := container.NewVBox(
		widget.NewLabel("複製檔案"),
		widget.NewLabel(fmt.Sprintf("來源: %s", filePath)),
		targetEntry,
	)
	
	// 建立對話框
	dialog := dialog.NewCustomConfirm("複製檔案", "複製", "取消", content, func(confirmed bool) {
		if confirmed {
			targetPath := targetEntry.Text
			if targetPath == "" {
				dialog.ShowError(fmt.Errorf("目標路徑不能為空"), mw.window)
				return
			}
			
			// 使用檔案管理服務複製檔案
			err := mw.fileManagerService.CopyFile(filePath, targetPath)
			if err != nil {
				dialog.ShowError(err, mw.window)
				return
			}
			
			// 重新整理檔案樹
			mw.refreshFileTree()
		}
	}, mw.window)
	
	// 顯示對話框
	dialog.Show()
	
	// 設定焦點到目標路徑輸入框
	targetEntry.FocusGained()
}

// moveFile 移動檔案
// 參數：filePath（要移動的檔案路徑）
//
// 執行流程：
// 1. 顯示目標路徑選擇對話框
// 2. 使用檔案管理服務移動檔案
// 3. 重新整理檔案樹
// 4. 處理移動結果
func (mw *MainWindow) moveFile(filePath string) {
	// 建立目標路徑輸入對話框
	targetEntry := widget.NewEntry()
	targetEntry.SetPlaceHolder("請輸入目標路徑...")
	
	// 建立對話框內容
	content := container.NewVBox(
		widget.NewLabel("移動檔案"),
		widget.NewLabel(fmt.Sprintf("來源: %s", filePath)),
		targetEntry,
	)
	
	// 建立對話框
	dialog := dialog.NewCustomConfirm("移動檔案", "移動", "取消", content, func(confirmed bool) {
		if confirmed {
			targetPath := targetEntry.Text
			if targetPath == "" {
				dialog.ShowError(fmt.Errorf("目標路徑不能為空"), mw.window)
				return
			}
			
			// 使用檔案管理服務移動檔案
			err := mw.fileManagerService.MoveFile(filePath, targetPath)
			if err != nil {
				dialog.ShowError(err, mw.window)
				return
			}
			
			// 重新整理檔案樹
			mw.refreshFileTree()
		}
	}, mw.window)
	
	// 顯示對話框
	dialog.Show()
	
	// 設定焦點到目標路徑輸入框
	targetEntry.FocusGained()
}

// handleFileTreeOperation 處理檔案樹的檔案操作
// 參數：operation（操作類型）、filePath（檔案路徑）
//
// 執行流程：
// 1. 根據操作類型執行相應的檔案操作
// 2. 顯示操作確認對話框（如果需要）
// 3. 使用檔案管理服務執行操作
// 4. 更新 UI 狀態和顯示操作回饋
func (mw *MainWindow) handleFileTreeOperation(operation, filePath string) {
	switch operation {
	case "create_file":
		mw.createNewFileInDirectory(filePath)
	case "create_folder":
		mw.createNewFolderInDirectory(filePath)
	case "rename":
		mw.renameFileWithDialog(filePath)
	case "delete":
		mw.deleteFileWithConfirmation(filePath)
	case "copy":
		mw.copyFileWithDialog(filePath)
	case "cut":
		mw.cutFileWithDialog(filePath)
	default:
		fmt.Printf("未知的檔案操作: %s\n", operation)
	}
}

// showFileContextMenu 顯示檔案或目錄的右鍵選單
// 參數：filePath（檔案路徑）、isDirectory（是否為目錄）
//
// 執行流程：
// 1. 根據檔案類型建立適當的選單項目
// 2. 顯示右鍵選單
// 3. 處理選單項目的點擊事件
func (mw *MainWindow) showFileContextMenu(filePath string, isDirectory bool) {
	// 建立右鍵選單項目
	var menuItems []*fyne.MenuItem
	
	if isDirectory {
		// 目錄的右鍵選單
		menuItems = []*fyne.MenuItem{
			fyne.NewMenuItem("開啟", func() {
				fmt.Printf("開啟目錄: %s\n", filePath)
			}),
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem("新增檔案", func() {
				mw.createNewFileInDirectory(filePath)
			}),
			fyne.NewMenuItem("新增資料夾", func() {
				mw.createNewFolderInDirectory(filePath)
			}),
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem("重新命名", func() {
				mw.renameFileWithDialog(filePath)
			}),
			fyne.NewMenuItem("刪除", func() {
				mw.deleteFileWithConfirmation(filePath)
			}),
		}
	} else {
		// 檔案的右鍵選單
		menuItems = []*fyne.MenuItem{
			fyne.NewMenuItem("開啟", func() {
				mw.openFileFromPath(filePath)
			}),
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem("重新命名", func() {
				mw.renameFileWithDialog(filePath)
			}),
			fyne.NewMenuItem("刪除", func() {
				mw.deleteFileWithConfirmation(filePath)
			}),
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem("複製", func() {
				mw.copyFileWithDialog(filePath)
			}),
		}
	}
	
	// 建立並顯示右鍵選單
	_ = fyne.NewMenu("", menuItems...)
	
	// 注意：Fyne 的右鍵選單需要特殊處理
	// 這裡使用簡化的實作，實際應用中可能需要更複雜的選單顯示邏輯
	fmt.Printf("顯示 %s 的右鍵選單\n", filePath)
}

// createNewFileInDirectory 在指定目錄中建立新檔案
// 參數：dirPath（目錄路徑）
//
// 執行流程：
// 1. 顯示檔案名稱輸入對話框
// 2. 驗證檔案名稱的有效性
// 3. 在指定目錄中建立新檔案
// 4. 重新整理檔案樹並顯示操作結果
func (mw *MainWindow) createNewFileInDirectory(dirPath string) {
	// 建立檔案名稱輸入對話框
	fileNameEntry := widget.NewEntry()
	fileNameEntry.SetPlaceHolder("請輸入檔案名稱（例如：note.md）...")
	
	// 建立對話框內容
	content := container.NewVBox(
		widget.NewLabel("在目錄中新增檔案"),
		widget.NewLabel(fmt.Sprintf("目錄: %s", dirPath)),
		fileNameEntry,
	)
	
	// 建立對話框
	dialog := dialog.NewCustomConfirm("新增檔案", "建立", "取消", content, func(confirmed bool) {
		if confirmed {
			fileName := fileNameEntry.Text
			if fileName == "" {
				dialog.ShowError(fmt.Errorf("檔案名稱不能為空"), mw.window)
				return
			}
			
			// 建立完整的檔案路徑
			_ = filepath.Join(dirPath, fileName)
			
			// 建立新檔案（透過編輯器服務）
			_, err := mw.editorService.CreateNote("", "")
			if err != nil {
				dialog.ShowError(fmt.Errorf("建立檔案失敗: %w", err), mw.window)
				return
			}
			
			// 重新整理檔案樹
			mw.refreshFileTree()
			
			// 顯示成功訊息
			dialog.ShowInformation("成功", fmt.Sprintf("檔案 '%s' 已建立", fileName), mw.window)
		}
	}, mw.window)
	
	// 顯示對話框
	dialog.Show()
	
	// 設定焦點到檔案名稱輸入框
	fileNameEntry.FocusGained()
}

// createNewFolderInDirectory 在指定目錄中建立新資料夾
// 參數：dirPath（目錄路徑）
//
// 執行流程：
// 1. 顯示資料夾名稱輸入對話框
// 2. 驗證資料夾名稱的有效性
// 3. 在指定目錄中建立新資料夾
// 4. 重新整理檔案樹並顯示操作結果
func (mw *MainWindow) createNewFolderInDirectory(dirPath string) {
	// 建立資料夾名稱輸入對話框
	folderNameEntry := widget.NewEntry()
	folderNameEntry.SetPlaceHolder("請輸入資料夾名稱...")
	
	// 建立對話框內容
	content := container.NewVBox(
		widget.NewLabel("在目錄中新增資料夾"),
		widget.NewLabel(fmt.Sprintf("目錄: %s", dirPath)),
		folderNameEntry,
	)
	
	// 建立對話框
	dialog := dialog.NewCustomConfirm("新增資料夾", "建立", "取消", content, func(confirmed bool) {
		if confirmed {
			folderName := folderNameEntry.Text
			if folderName == "" {
				dialog.ShowError(fmt.Errorf("資料夾名稱不能為空"), mw.window)
				return
			}
			
			// 建立完整的資料夾路徑
			folderPath := filepath.Join(dirPath, folderName)
			
			// 使用檔案管理服務建立資料夾
			err := mw.fileManagerService.CreateDirectory(folderPath)
			if err != nil {
				dialog.ShowError(fmt.Errorf("建立資料夾失敗: %w", err), mw.window)
				return
			}
			
			// 重新整理檔案樹
			mw.refreshFileTree()
			
			// 顯示成功訊息
			dialog.ShowInformation("成功", fmt.Sprintf("資料夾 '%s' 已建立", folderName), mw.window)
		}
	}, mw.window)
	
	// 顯示對話框
	dialog.Show()
	
	// 設定焦點到資料夾名稱輸入框
	folderNameEntry.FocusGained()
}

// createNewFileInCurrentDir 在當前目錄中建立新檔案
// 這是工具欄按鈕的回調函數
//
// 執行流程：
// 1. 取得當前選擇的目錄或使用根目錄
// 2. 調用 createNewFileInDirectory 方法
func (mw *MainWindow) createNewFileInCurrentDir() {
	// 使用根目錄作為預設位置
	rootPath := mw.settings.DefaultSaveLocation
	if rootPath == "" {
		rootPath = "."
	}
	
	mw.createNewFileInDirectory(rootPath)
}

// renameFileWithDialog 顯示重新命名對話框並執行重新命名操作
// 參數：filePath（要重新命名的檔案路徑）
//
// 執行流程：
// 1. 顯示新名稱輸入對話框
// 2. 驗證新名稱的有效性
// 3. 使用檔案管理服務執行重新命名
// 4. 重新整理檔案樹並顯示操作結果
func (mw *MainWindow) renameFileWithDialog(filePath string) {
	// 取得當前檔案名稱
	currentName := filepath.Base(filePath)
	
	// 建立新名稱輸入對話框
	nameEntry := widget.NewEntry()
	nameEntry.SetText(currentName)
	
	// 建立對話框內容
	content := container.NewVBox(
		widget.NewLabel("重新命名"),
		widget.NewLabel(fmt.Sprintf("當前名稱: %s", currentName)),
		nameEntry,
	)
	
	// 建立對話框
	dialog := dialog.NewCustomConfirm("重新命名", "確定", "取消", content, func(confirmed bool) {
		if confirmed {
			newName := nameEntry.Text
			if newName == "" || newName == currentName {
				return
			}
			
			// 建立新路徑
			newPath := filepath.Join(filepath.Dir(filePath), newName)
			
			// 使用檔案管理服務重新命名檔案
			err := mw.fileManagerService.RenameFile(filePath, newPath)
			if err != nil {
				dialog.ShowError(fmt.Errorf("重新命名失敗: %w", err), mw.window)
				return
			}
			
			// 重新整理檔案樹
			mw.refreshFileTree()
			
			// 顯示成功訊息
			dialog.ShowInformation("成功", fmt.Sprintf("已重新命名為 '%s'", newName), mw.window)
		}
	}, mw.window)
	
	// 顯示對話框
	dialog.Show()
	
	// 設定焦點到名稱輸入框並選擇文字
	nameEntry.FocusGained()
}

// deleteFileWithConfirmation 顯示刪除確認對話框並執行刪除操作
// 參數：filePath（要刪除的檔案路徑）
//
// 執行流程：
// 1. 顯示刪除確認對話框
// 2. 如果用戶確認，使用檔案管理服務執行刪除
// 3. 重新整理檔案樹並顯示操作結果
func (mw *MainWindow) deleteFileWithConfirmation(filePath string) {
	fileName := filepath.Base(filePath)
	
	// 顯示確認對話框
	dialog.ShowConfirm("確認刪除", 
		fmt.Sprintf("確定要刪除 '%s' 嗎？\n\n此操作無法復原。", fileName), 
		func(confirmed bool) {
			if confirmed {
				// 使用檔案管理服務刪除檔案
				err := mw.fileManagerService.DeleteFile(filePath)
				if err != nil {
					dialog.ShowError(fmt.Errorf("刪除失敗: %w", err), mw.window)
					return
				}
				
				// 重新整理檔案樹
				mw.refreshFileTree()
				
				// 顯示成功訊息
				dialog.ShowInformation("成功", fmt.Sprintf("'%s' 已刪除", fileName), mw.window)
			}
		}, mw.window)
}

// copyFileWithDialog 顯示複製對話框並執行複製操作
// 參數：filePath（要複製的檔案路徑）
//
// 執行流程：
// 1. 顯示目標路徑輸入對話框
// 2. 驗證目標路徑的有效性
// 3. 使用檔案管理服務執行複製
// 4. 重新整理檔案樹並顯示操作結果
func (mw *MainWindow) copyFileWithDialog(filePath string) {
	fileName := filepath.Base(filePath)
	
	// 建立目標路徑輸入對話框
	targetEntry := widget.NewEntry()
	targetEntry.SetPlaceHolder("請輸入目標路徑...")
	
	// 建立對話框內容
	content := container.NewVBox(
		widget.NewLabel("複製檔案"),
		widget.NewLabel(fmt.Sprintf("來源: %s", fileName)),
		targetEntry,
	)
	
	// 建立對話框
	dialog := dialog.NewCustomConfirm("複製檔案", "複製", "取消", content, func(confirmed bool) {
		if confirmed {
			targetPath := targetEntry.Text
			if targetPath == "" {
				dialog.ShowError(fmt.Errorf("目標路徑不能為空"), mw.window)
				return
			}
			
			// 使用檔案管理服務複製檔案
			err := mw.fileManagerService.CopyFile(filePath, targetPath)
			if err != nil {
				dialog.ShowError(fmt.Errorf("複製失敗: %w", err), mw.window)
				return
			}
			
			// 重新整理檔案樹
			mw.refreshFileTree()
			
			// 顯示成功訊息
			dialog.ShowInformation("成功", fmt.Sprintf("'%s' 已複製到 '%s'", fileName, targetPath), mw.window)
		}
	}, mw.window)
	
	// 顯示對話框
	dialog.Show()
	
	// 設定焦點到目標路徑輸入框
	targetEntry.FocusGained()
}

// cutFileWithDialog 顯示剪下對話框並執行移動操作
// 參數：filePath（要剪下的檔案路徑）
//
// 執行流程：
// 1. 顯示目標路徑輸入對話框
// 2. 驗證目標路徑的有效性
// 3. 使用檔案管理服務執行移動
// 4. 重新整理檔案樹並顯示操作結果
func (mw *MainWindow) cutFileWithDialog(filePath string) {
	fileName := filepath.Base(filePath)
	
	// 建立目標路徑輸入對話框
	targetEntry := widget.NewEntry()
	targetEntry.SetPlaceHolder("請輸入目標路徑...")
	
	// 建立對話框內容
	content := container.NewVBox(
		widget.NewLabel("剪下檔案"),
		widget.NewLabel(fmt.Sprintf("來源: %s", fileName)),
		targetEntry,
	)
	
	// 建立對話框
	dialog := dialog.NewCustomConfirm("剪下檔案", "移動", "取消", content, func(confirmed bool) {
		if confirmed {
			targetPath := targetEntry.Text
			if targetPath == "" {
				dialog.ShowError(fmt.Errorf("目標路徑不能為空"), mw.window)
				return
			}
			
			// 使用檔案管理服務移動檔案
			err := mw.fileManagerService.MoveFile(filePath, targetPath)
			if err != nil {
				dialog.ShowError(fmt.Errorf("移動失敗: %w", err), mw.window)
				return
			}
			
			// 重新整理檔案樹
			mw.refreshFileTree()
			
			// 顯示成功訊息
			dialog.ShowInformation("成功", fmt.Sprintf("'%s' 已移動到 '%s'", fileName, targetPath), mw.window)
		}
	}, mw.window)
	
	// 顯示對話框
	dialog.Show()
	
	// 設定焦點到目標路徑輸入框
	targetEntry.FocusGained()
}

// FileOpenDialog 代表檔案開啟對話框的包裝器
type FileOpenDialog struct {
	dialog *dialog.FileDialog
}

// NewFileOpenDialog 建立檔案開啟對話框（暫時實作）
// 在實際應用中應該使用 Fyne 的檔案對話框
func NewFileOpenDialog(parent fyne.Window, callback func(string)) *FileOpenDialog {
	// 暫時實作，回傳一個基本的檔案對話框
	fileDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err == nil && reader != nil {
			callback(reader.URI().Path())
			reader.Close()
		}
	}, parent)
	
	return &FileOpenDialog{dialog: fileDialog}
}

// Show 顯示檔案開啟對話框
func (fod *FileOpenDialog) Show() {
	fod.dialog.Show()
}

// FileSaveDialog 代表檔案保存對話框的包裝器
type FileSaveDialog struct {
	dialog *dialog.FileDialog
}

// NewFileSaveDialog 建立檔案保存對話框（暫時實作）
// 在實際應用中應該使用 Fyne 的檔案對話框
func NewFileSaveDialog(parent fyne.Window, callback func(string)) *FileSaveDialog {
	// 暫時實作，回傳一個基本的檔案對話框
	fileDialog := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err == nil && writer != nil {
			callback(writer.URI().Path())
			writer.Close()
		}
	}, parent)
	
	return &FileSaveDialog{dialog: fileDialog}
}

// Show 顯示檔案保存對話框
func (fsd *FileSaveDialog) Show() {
	fsd.dialog.Show()
}

// SetFileName 設定檔案對話框的預設檔案名稱（暫時實作）
func (fsd *FileSaveDialog) SetFileName(name string) {
	// 暫時實作，Fyne 的 FileDialog 沒有直接的 SetFileName 方法
	// 在實際應用中可能需要使用其他方式設定預設名稱
}

// PasswordDialog 代表密碼輸入對話框的包裝器
type PasswordDialog struct {
	dialog *dialog.ConfirmDialog
}

// NewPasswordDialog 建立密碼輸入對話框（暫時實作）
func NewPasswordDialog(parent fyne.Window, title string, callback func(string)) *PasswordDialog {
	// 建立密碼輸入框
	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("請輸入密碼...")
	
	// 建立對話框內容
	content := container.NewVBox(
		widget.NewLabel(title),
		passwordEntry,
	)
	
	// 建立對話框
	customDialog := dialog.NewCustomConfirm("密碼驗證", "確定", "取消", content, func(confirmed bool) {
		if confirmed {
			callback(passwordEntry.Text)
		}
	}, parent)
	
	return &PasswordDialog{dialog: customDialog}
}

// Show 顯示密碼輸入對話框
func (pd *PasswordDialog) Show() {
	pd.dialog.Show()
}