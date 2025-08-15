// Package ui 包含使用者介面相關的元件和視窗管理測試
// 測試 MainWindow 的建立、初始化和基本功能
package ui

import (
	"testing"                    // Go 標準測試套件
	"fyne.io/fyne/v2/test"       // Fyne 測試工具套件

	"mac-notebook-app/internal/models"
	"mac-notebook-app/internal/services"
	"mac-notebook-app/internal/repositories"
)

// TestNewMainWindow 測試主視窗的建立和初始化
// 驗證主視窗是否正確建立並包含所有必要的 UI 元件
//
// 測試項目：
// 1. 主視窗實例是否正確建立
// 2. 視窗標題是否正確設定
// 3. 視窗大小是否符合預期
// 4. 所有 UI 元件是否正確初始化
func TestNewMainWindow(t *testing.T) {
	// 建立測試用的 Fyne 應用程式
	testApp := test.NewApp()
	defer testApp.Quit()

	// 建立測試設定
	testSettings := models.NewDefaultSettings()
	
	// 建立測試服務
	fileRepo, _ := repositories.NewLocalFileRepository("./test_notes")
	encryptionSvc := services.NewEncryptionService()
	passwordSvc := services.NewPasswordService()
	biometricSvc := services.NewBiometricService()
	editorService := services.NewEditorService(fileRepo, encryptionSvc, passwordSvc, biometricSvc, services.NewPerformanceService(nil))
	fileManagerService, _ := services.NewLocalFileManagerService(fileRepo, "./test_notes")
	
	// 建立主視窗實例
	mainWindow := NewMainWindow(testApp, testSettings, editorService, fileManagerService)
	
	// 驗證主視窗實例不為 nil
	if mainWindow == nil {
		t.Fatal("NewMainWindow 應該回傳有效的 MainWindow 實例")
	}
	
	// 驗證視窗實例不為 nil
	if mainWindow.window == nil {
		t.Error("主視窗的 window 欄位不應該為 nil")
	}
	
	// 驗證視窗標題
	expectedTitle := "Mac Notebook App - 安全筆記編輯器"
	if mainWindow.window.Title() != expectedTitle {
		t.Errorf("視窗標題應該是 '%s'，但得到 '%s'", expectedTitle, mainWindow.window.Title())
	}
	
	// 驗證視窗大小
	size := mainWindow.window.Canvas().Size()
	expectedWidth := float32(1200)
	expectedHeight := float32(800)
	
	if size.Width != expectedWidth {
		t.Errorf("視窗寬度應該是 %f，但得到 %f", expectedWidth, size.Width)
	}
	
	if size.Height != expectedHeight {
		t.Errorf("視窗高度應該是 %f，但得到 %f", expectedHeight, size.Height)
	}

	// 驗證設定和服務已初始化
	if mainWindow.settings == nil {
		t.Error("主視窗的 settings 欄位不應該為 nil")
	}
	
	if mainWindow.themeService == nil {
		t.Error("主視窗的 themeService 欄位不應該為 nil")
	}
	
	if mainWindow.editorService == nil {
		t.Error("主視窗的 editorService 欄位不應該為 nil")
	}
	
	if mainWindow.fileManagerService == nil {
		t.Error("主視窗的 fileManagerService 欄位不應該為 nil")
	}
}

// TestMainWindowUIComponents 測試主視窗 UI 元件的初始化
// 驗證所有主要 UI 元件是否正確建立和配置
//
// 測試項目：
// 1. 選單欄是否正確建立
// 2. 工具欄是否正確建立
// 3. 狀態欄元件是否正確初始化
// 4. 主要內容區域是否正確設定
func TestMainWindowUIComponents(t *testing.T) {
	// 建立測試用的 Fyne 應用程式
	testApp := test.NewApp()
	defer testApp.Quit()

	// 建立測試設定
	testSettings := models.NewDefaultSettings()
	
	// 建立測試服務
	fileRepo, _ := repositories.NewLocalFileRepository("./test_notes")
	encryptionSvc := services.NewEncryptionService()
	passwordSvc := services.NewPasswordService()
	biometricSvc := services.NewBiometricService()
	editorService := services.NewEditorService(fileRepo, encryptionSvc, passwordSvc, biometricSvc, services.NewPerformanceService(nil))
	fileManagerService, _ := services.NewLocalFileManagerService(fileRepo, "./test_notes")
	
	// 建立主視窗實例
	mainWindow := NewMainWindow(testApp, testSettings, editorService, fileManagerService)
	
	// 驗證選單欄已建立
	if mainWindow.menuBar == nil {
		t.Error("主視窗的選單欄不應該為 nil")
	}
	
	// 驗證工具欄已建立
	if mainWindow.toolBar == nil {
		t.Error("主視窗的工具欄不應該為 nil")
	}
	
	// 驗證狀態欄已建立
	if mainWindow.statusBar == nil {
		t.Error("主視窗的狀態欄不應該為 nil")
	}
	
	// 驗證狀態欄元件已初始化
	if mainWindow.saveStatus == nil {
		t.Error("保存狀態標籤不應該為 nil")
	}
	
	if mainWindow.encStatus == nil {
		t.Error("加密狀態標籤不應該為 nil")
	}
	
	if mainWindow.wordCount == nil {
		t.Error("字數統計標籤不應該為 nil")
	}
	
	// 驗證主要內容區域已建立
	if mainWindow.content == nil {
		t.Error("主要內容容器不應該為 nil")
	}
	
	if mainWindow.leftPanel == nil {
		t.Error("左側面板不應該為 nil")
	}
	
	if mainWindow.rightPanel == nil {
		t.Error("右側面板不應該為 nil")
	}
	
	if mainWindow.mainSplit == nil {
		t.Error("主要分割容器不應該為 nil")
	}
	
	// 驗證編輯器和檔案樹元件已建立
	if mainWindow.editor == nil {
		t.Error("Markdown 編輯器不應該為 nil")
	}
	
	if mainWindow.fileTree == nil {
		t.Error("檔案樹元件不應該為 nil")
	}
}

// TestMainWindowStatusUpdates 測試主視窗狀態更新功能
// 驗證狀態欄的各種更新功能是否正常工作
//
// 測試項目：
// 1. 保存狀態更新功能
// 2. 加密狀態更新功能
// 3. 字數統計更新功能
func TestMainWindowStatusUpdates(t *testing.T) {
	// 建立測試用的 Fyne 應用程式
	testApp := test.NewApp()
	defer testApp.Quit()

	// 建立測試設定
	testSettings := models.NewDefaultSettings()
	
	// 建立測試服務
	fileRepo, _ := repositories.NewLocalFileRepository("./test_notes")
	encryptionSvc := services.NewEncryptionService()
	passwordSvc := services.NewPasswordService()
	biometricSvc := services.NewBiometricService()
	editorService := services.NewEditorService(fileRepo, encryptionSvc, passwordSvc, biometricSvc, services.NewPerformanceService(nil))
	fileManagerService, _ := services.NewLocalFileManagerService(fileRepo, "./test_notes")
	
	// 建立主視窗實例
	mainWindow := NewMainWindow(testApp, testSettings, editorService, fileManagerService)
	
	// 測試保存狀態更新
	testSaveStatus := "正在儲存..."
	mainWindow.UpdateSaveStatus(testSaveStatus)
	if mainWindow.saveStatus.Text != testSaveStatus {
		t.Errorf("保存狀態更新失敗，期望 '%s'，得到 '%s'", testSaveStatus, mainWindow.saveStatus.Text)
	}
	
	// 測試加密狀態更新
	mainWindow.UpdateEncryptionStatus(true, "AES-256")
	expectedEncStatus := "已加密 (AES-256)"
	if mainWindow.encStatus.Text != expectedEncStatus {
		t.Errorf("加密狀態更新失敗，期望 '%s'，得到 '%s'", expectedEncStatus, mainWindow.encStatus.Text)
	}
	
	// 測試未加密狀態
	mainWindow.UpdateEncryptionStatus(false, "")
	expectedUnencStatus := "未加密"
	if mainWindow.encStatus.Text != expectedUnencStatus {
		t.Errorf("未加密狀態更新失敗，期望 '%s'，得到 '%s'", expectedUnencStatus, mainWindow.encStatus.Text)
	}
	
	// 測試字數統計更新
	testWordCount := 150
	mainWindow.UpdateWordCount(testWordCount)
	expectedWordCountText := "字數: 150"
	if mainWindow.wordCount.Text != expectedWordCountText {
		t.Errorf("字數統計更新失敗，期望 '%s'，得到 '%s'", expectedWordCountText, mainWindow.wordCount.Text)
	}
}

// TestMainWindowGetWindow 測試 GetWindow 方法
// 驗證 GetWindow 方法是否回傳正確的視窗實例
//
// 測試項目：
// 1. GetWindow 方法是否回傳正確的視窗實例
// 2. 回傳的視窗實例是否與原始視窗相同
func TestMainWindowGetWindow(t *testing.T) {
	// 建立測試用的 Fyne 應用程式
	testApp := test.NewApp()
	defer testApp.Quit()

	// 建立測試設定
	testSettings := models.NewDefaultSettings()
	
	// 建立測試服務
	fileRepo, _ := repositories.NewLocalFileRepository("./test_notes")
	encryptionSvc := services.NewEncryptionService()
	passwordSvc := services.NewPasswordService()
	biometricSvc := services.NewBiometricService()
	editorService := services.NewEditorService(fileRepo, encryptionSvc, passwordSvc, biometricSvc, services.NewPerformanceService(nil))
	fileManagerService, _ := services.NewLocalFileManagerService(fileRepo, "./test_notes")
	
	// 建立主視窗實例
	mainWindow := NewMainWindow(testApp, testSettings, editorService, fileManagerService)
	
	// 測試 GetWindow 方法
	window := mainWindow.GetWindow()
	
	// 驗證回傳的視窗不為 nil
	if window == nil {
		t.Error("GetWindow 不應該回傳 nil")
	}
	
	// 驗證回傳的視窗與原始視窗相同
	if window != mainWindow.window {
		t.Error("GetWindow 應該回傳與內部 window 欄位相同的實例")
	}
}

// TestMainWindowSplitRatio 測試主視窗分割比例設定
// 驗證左右面板的分割比例是否正確設定
//
// 測試項目：
// 1. 分割容器的偏移量是否正確設定為 0.3
func TestMainWindowSplitRatio(t *testing.T) {
	// 建立測試用的 Fyne 應用程式
	testApp := test.NewApp()
	defer testApp.Quit()

	// 建立測試設定
	testSettings := models.NewDefaultSettings()
	
	// 建立測試服務
	fileRepo, _ := repositories.NewLocalFileRepository("./test_notes")
	encryptionSvc := services.NewEncryptionService()
	passwordSvc := services.NewPasswordService()
	biometricSvc := services.NewBiometricService()
	editorService := services.NewEditorService(fileRepo, encryptionSvc, passwordSvc, biometricSvc, services.NewPerformanceService(nil))
	fileManagerService, _ := services.NewLocalFileManagerService(fileRepo, "./test_notes")
	
	// 建立主視窗實例
	mainWindow := NewMainWindow(testApp, testSettings, editorService, fileManagerService)
	
	// 驗證分割容器的偏移量
	expectedOffset := 0.3
	actualOffset := mainWindow.mainSplit.Offset
	
	if actualOffset != expectedOffset {
		t.Errorf("分割容器偏移量應該是 %f，但得到 %f", expectedOffset, actualOffset)
	}
}

// TestMainWindowThemeIntegration 測試主題整合功能
// 驗證主題服務是否正確整合到主視窗中
//
// 測試項目：
// 1. 主題服務是否正確初始化
// 2. 主題監聽器是否正確設定
// 3. 設定變更是否正確處理
func TestMainWindowThemeIntegration(t *testing.T) {
	// 建立測試用的 Fyne 應用程式
	testApp := test.NewApp()
	defer testApp.Quit()

	// 建立測試設定
	testSettings := models.NewDefaultSettings()
	testSettings.Theme = "dark"
	
	// 建立測試服務
	fileRepo, _ := repositories.NewLocalFileRepository("./test_notes")
	encryptionSvc := services.NewEncryptionService()
	passwordSvc := services.NewPasswordService()
	biometricSvc := services.NewBiometricService()
	editorService := services.NewEditorService(fileRepo, encryptionSvc, passwordSvc, biometricSvc, services.NewPerformanceService(nil))
	fileManagerService, _ := services.NewLocalFileManagerService(fileRepo, "./test_notes")
	
	// 建立主視窗實例
	mainWindow := NewMainWindow(testApp, testSettings, editorService, fileManagerService)
	
	// 驗證主題服務已初始化
	if mainWindow.themeService == nil {
		t.Fatal("主題服務不應該為 nil")
	}
	
	// 驗證當前主題設定
	currentTheme := mainWindow.themeService.GetCurrentTheme()
	if currentTheme != "dark" {
		t.Errorf("當前主題應該是 'dark'，但得到 '%s'", currentTheme)
	}
	
	// 測試設定變更處理
	newSettings := models.NewDefaultSettings()
	newSettings.Theme = "light"
	newSettings.DefaultEncryption = "chacha20"
	newSettings.AutoSaveInterval = 10
	
	mainWindow.onSettingsChanged(newSettings)
	
	// 驗證設定已更新
	if mainWindow.settings.Theme != "light" {
		t.Error("設定變更後主題應該更新為 'light'")
	}
	
	if mainWindow.settings.DefaultEncryption != "chacha20" {
		t.Error("設定變更後加密演算法應該更新為 'chacha20'")
	}
}

// TestMainWindowOnThemeChanged 測試主題變更監聽器
// 驗證主題變更監聽器是否正確實作
//
// 測試項目：
// 1. OnThemeChanged 方法是否正確實作
// 2. 主題變更時 UI 是否正確更新
func TestMainWindowOnThemeChanged(t *testing.T) {
	// 建立測試用的 Fyne 應用程式
	testApp := test.NewApp()
	defer testApp.Quit()

	// 建立測試設定
	testSettings := models.NewDefaultSettings()
	
	// 建立測試服務
	fileRepo, _ := repositories.NewLocalFileRepository("./test_notes")
	encryptionSvc := services.NewEncryptionService()
	passwordSvc := services.NewPasswordService()
	biometricSvc := services.NewBiometricService()
	editorService := services.NewEditorService(fileRepo, encryptionSvc, passwordSvc, biometricSvc, services.NewPerformanceService(nil))
	fileManagerService, _ := services.NewLocalFileManagerService(fileRepo, "./test_notes")
	
	// 建立主視窗實例
	mainWindow := NewMainWindow(testApp, testSettings, editorService, fileManagerService)
	
	// 測試主題變更監聽器
	mainWindow.OnThemeChanged("dark")
	
	// 驗證方法執行不會產生錯誤
	// 實際的主題變更效果由 Fyne 框架處理
	// 這裡主要確保方法能正常執行
}

// TestMainWindowEditorServiceIntegration 測試編輯器服務整合
// 驗證編輯器服務是否正確整合到主視窗中
//
// 測試項目：
// 1. 編輯器元件是否正確建立
// 2. 編輯器服務是否正確注入
// 3. 編輯器回調函數是否正確設定
func TestMainWindowEditorServiceIntegration(t *testing.T) {
	// 建立測試用的 Fyne 應用程式
	testApp := test.NewApp()
	defer testApp.Quit()

	// 建立測試設定
	testSettings := models.NewDefaultSettings()
	
	// 建立測試服務
	fileRepo, _ := repositories.NewLocalFileRepository("./test_notes")
	encryptionSvc := services.NewEncryptionService()
	passwordSvc := services.NewPasswordService()
	biometricSvc := services.NewBiometricService()
	editorService := services.NewEditorService(fileRepo, encryptionSvc, passwordSvc, biometricSvc, services.NewPerformanceService(nil))
	fileManagerService, _ := services.NewLocalFileManagerService(fileRepo, "./test_notes")
	
	// 建立主視窗實例
	mainWindow := NewMainWindow(testApp, testSettings, editorService, fileManagerService)
	
	// 驗證編輯器服務已正確注入
	if mainWindow.editorService == nil {
		t.Fatal("編輯器服務不應該為 nil")
	}
	
	// 驗證編輯器元件已建立
	if mainWindow.editor == nil {
		t.Fatal("編輯器元件不應該為 nil")
	}
	
	// 驗證編輯器容器已正確嵌入到右側面板
	if mainWindow.rightPanel == nil {
		t.Error("右側面板不應該為 nil")
	}
	
	editorContainer := mainWindow.editor.GetContainer()
	if editorContainer == nil {
		t.Error("編輯器容器不應該為 nil")
	}
}

// TestMainWindowFileManagerServiceIntegration 測試檔案管理服務整合
// 驗證檔案管理服務是否正確整合到主視窗中
//
// 測試項目：
// 1. 檔案樹元件是否正確建立
// 2. 檔案管理服務是否正確注入
// 3. 檔案樹回調函數是否正確設定
func TestMainWindowFileManagerServiceIntegration(t *testing.T) {
	// 建立測試用的 Fyne 應用程式
	testApp := test.NewApp()
	defer testApp.Quit()

	// 建立測試設定
	testSettings := models.NewDefaultSettings()
	
	// 建立測試服務
	fileRepo, _ := repositories.NewLocalFileRepository("./test_notes")
	encryptionSvc := services.NewEncryptionService()
	passwordSvc := services.NewPasswordService()
	biometricSvc := services.NewBiometricService()
	editorService := services.NewEditorService(fileRepo, encryptionSvc, passwordSvc, biometricSvc, services.NewPerformanceService(nil))
	fileManagerService, _ := services.NewLocalFileManagerService(fileRepo, "./test_notes")
	
	// 建立主視窗實例
	mainWindow := NewMainWindow(testApp, testSettings, editorService, fileManagerService)
	
	// 驗證檔案管理服務已正確注入
	if mainWindow.fileManagerService == nil {
		t.Fatal("檔案管理服務不應該為 nil")
	}
	
	// 驗證檔案樹元件已建立
	if mainWindow.fileTree == nil {
		t.Fatal("檔案樹元件不應該為 nil")
	}
	
	// 驗證檔案樹容器已正確嵌入到左側面板
	if mainWindow.leftPanel == nil {
		t.Error("左側面板不應該為 nil")
	}
	
	fileTreeContainer := mainWindow.fileTree.GetContainer()
	if fileTreeContainer == nil {
		t.Error("檔案樹容器不應該為 nil")
	}
}