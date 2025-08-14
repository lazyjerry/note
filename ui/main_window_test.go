// Package ui 包含使用者介面相關的元件和視窗管理測試
// 測試 MainWindow 的建立、初始化和基本功能
package ui

import (
	"testing"                    // Go 標準測試套件
	"fyne.io/fyne/v2/test"       // Fyne 測試工具套件
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
	
	// 建立主視窗實例
	mainWindow := NewMainWindow(testApp)
	
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
}

// TestMainWindowUIComponents 測試主視窗的 UI 元件初始化
// 驗證選單欄、工具欄、狀態欄等元件是否正確建立
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
	
	// 建立主視窗實例
	mainWindow := NewMainWindow(testApp)
	
	// 驗證選單欄
	if mainWindow.menuBar == nil {
		t.Error("選單欄不應該為 nil")
	}
	
	// 驗證工具欄
	if mainWindow.toolBar == nil {
		t.Error("工具欄不應該為 nil")
	}
	
	// 驗證狀態欄容器
	if mainWindow.statusBar == nil {
		t.Error("狀態欄容器不應該為 nil")
	}
	
	// 驗證狀態欄元件
	if mainWindow.saveStatus == nil {
		t.Error("保存狀態指示器不應該為 nil")
	}
	
	if mainWindow.encStatus == nil {
		t.Error("加密狀態指示器不應該為 nil")
	}
	
	if mainWindow.wordCount == nil {
		t.Error("字數統計顯示不應該為 nil")
	}
	
	// 驗證主要內容區域
	if mainWindow.leftPanel == nil {
		t.Error("左側面板不應該為 nil")
	}
	
	if mainWindow.rightPanel == nil {
		t.Error("右側面板不應該為 nil")
	}
	
	if mainWindow.mainSplit == nil {
		t.Error("主要分割容器不應該為 nil")
	}
	
	if mainWindow.content == nil {
		t.Error("主要內容容器不應該為 nil")
	}
}

// TestMainWindowStatusUpdates 測試主視窗狀態更新功能
// 驗證保存狀態、加密狀態和字數統計的更新功能
//
// 測試項目：
// 1. 保存狀態更新功能
// 2. 加密狀態更新功能
// 3. 字數統計更新功能
func TestMainWindowStatusUpdates(t *testing.T) {
	// 建立測試用的 Fyne 應用程式
	testApp := test.NewApp()
	defer testApp.Quit()
	
	// 建立主視窗實例
	mainWindow := NewMainWindow(testApp)
	
	// 測試保存狀態更新
	testSaveStatus := "正在儲存..."
	mainWindow.UpdateSaveStatus(testSaveStatus)
	if mainWindow.saveStatus.Text != testSaveStatus {
		t.Errorf("保存狀態應該是 '%s'，但得到 '%s'", testSaveStatus, mainWindow.saveStatus.Text)
	}
	
	// 測試加密狀態更新（已加密）
	mainWindow.UpdateEncryptionStatus(true, "AES-256")
	expectedEncStatus := "已加密 (AES-256)"
	if mainWindow.encStatus.Text != expectedEncStatus {
		t.Errorf("加密狀態應該是 '%s'，但得到 '%s'", expectedEncStatus, mainWindow.encStatus.Text)
	}
	
	// 測試加密狀態更新（未加密）
	mainWindow.UpdateEncryptionStatus(false, "")
	expectedEncStatus = "未加密"
	if mainWindow.encStatus.Text != expectedEncStatus {
		t.Errorf("加密狀態應該是 '%s'，但得到 '%s'", expectedEncStatus, mainWindow.encStatus.Text)
	}
	
	// 測試字數統計更新
	testWordCount := 150
	mainWindow.UpdateWordCount(testWordCount)
	expectedWordCount := "字數: 150"
	if mainWindow.wordCount.Text != expectedWordCount {
		t.Errorf("字數統計應該是 '%s'，但得到 '%s'", expectedWordCount, mainWindow.wordCount.Text)
	}
}

// TestMainWindowGetWindow 測試取得視窗實例功能
// 驗證 GetWindow 方法是否正確回傳視窗實例
//
// 測試項目：
// 1. GetWindow 方法是否回傳正確的視窗實例
// 2. 回傳的視窗實例是否與原始視窗相同
func TestMainWindowGetWindow(t *testing.T) {
	// 建立測試用的 Fyne 應用程式
	testApp := test.NewApp()
	defer testApp.Quit()
	
	// 建立主視窗實例
	mainWindow := NewMainWindow(testApp)
	
	// 取得視窗實例
	window := mainWindow.GetWindow()
	
	// 驗證回傳的視窗實例不為 nil
	if window == nil {
		t.Error("GetWindow 應該回傳有效的視窗實例")
	}
	
	// 驗證回傳的視窗實例與原始視窗相同
	if window != mainWindow.window {
		t.Error("GetWindow 應該回傳與原始視窗相同的實例")
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
	
	// 建立主視窗實例
	mainWindow := NewMainWindow(testApp)
	
	// 驗證分割比例
	expectedOffset := 0.3
	if mainWindow.mainSplit.Offset != expectedOffset {
		t.Errorf("分割容器的偏移量應該是 %f，但得到 %f", expectedOffset, mainWindow.mainSplit.Offset)
	}
}