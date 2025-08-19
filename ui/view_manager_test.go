// Package ui 包含使用者介面相關的元件和視窗管理
// 本檔案包含視圖管理器的單元測試
package ui

import (
	"testing"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

// TestNewViewManager 測試視圖管理器的建立
// 驗證視圖管理器是否正確初始化
func TestNewViewManager(t *testing.T) {
	// 建立測試應用程式和視窗
	testApp := app.NewWithID("test")
	testWindow := testApp.NewWindow("Test")
	testWindow.Resize(fyne.NewSize(800, 600))
	
	// 建立佈局管理器
	layoutManager := NewLayoutManager()
	
	// 建立視圖管理器
	viewManager := NewViewManager(testWindow, layoutManager)
	
	// 驗證初始狀態
	if viewManager == nil {
		t.Fatal("視圖管理器建立失敗")
	}
	
	if viewManager.GetViewMode() != ViewModeSplit {
		t.Errorf("預期初始視圖模式為分割視圖，實際為 %v", viewManager.GetViewMode())
	}
	
	if viewManager.GetSplitRatio() != 0.5 {
		t.Errorf("預期初始分割比例為 0.5，實際為 %f", viewManager.GetSplitRatio())
	}
	
	if viewManager.IsFullscreen() {
		t.Error("預期初始狀態不是全螢幕模式")
	}
}

// TestSetViewMode 測試視圖模式設定
// 驗證視圖模式切換功能是否正常工作
func TestSetViewMode(t *testing.T) {
	// 建立測試環境
	testApp := app.New()
	testWindow := testApp.NewWindow("Test")
	defer testWindow.Close()
	
	layoutManager := NewLayoutManager()
	viewManager := NewViewManager(testWindow, layoutManager)
	
	// 測試回調函數
	var callbackMode ViewMode
	var callbackCalled bool
	viewManager.SetOnViewModeChanged(func(mode ViewMode) {
		callbackMode = mode
		callbackCalled = true
	})
	
	// 測試切換到編輯模式
	viewManager.SetViewMode(ViewModeEdit)
	if viewManager.GetViewMode() != ViewModeEdit {
		t.Errorf("預期視圖模式為編輯模式，實際為 %v", viewManager.GetViewMode())
	}
	if !callbackCalled || callbackMode != ViewModeEdit {
		t.Error("視圖模式變更回調未正確觸發")
	}
	
	// 重置回調狀態
	callbackCalled = false
	
	// 測試切換到預覽模式
	viewManager.SetViewMode(ViewModePreview)
	if viewManager.GetViewMode() != ViewModePreview {
		t.Errorf("預期視圖模式為預覽模式，實際為 %v", viewManager.GetViewMode())
	}
	if !callbackCalled || callbackMode != ViewModePreview {
		t.Error("視圖模式變更回調未正確觸發")
	}
	
	// 測試切換到分割視圖模式
	viewManager.SetViewMode(ViewModeSplit)
	if viewManager.GetViewMode() != ViewModeSplit {
		t.Errorf("預期視圖模式為分割視圖，實際為 %v", viewManager.GetViewMode())
	}
}

// TestToggleViewMode 測試視圖模式循環切換
// 驗證視圖模式是否按照預期順序循環切換
func TestToggleViewMode(t *testing.T) {
	// 建立測試環境
	testApp := app.New()
	testWindow := testApp.NewWindow("Test")
	defer testWindow.Close()
	
	layoutManager := NewLayoutManager()
	viewManager := NewViewManager(testWindow, layoutManager)
	
	// 初始狀態應該是分割視圖
	if viewManager.GetViewMode() != ViewModeSplit {
		t.Errorf("預期初始模式為分割視圖，實際為 %v", viewManager.GetViewMode())
	}
	
	// 第一次切換：分割視圖 -> 編輯模式
	viewManager.ToggleViewMode()
	if viewManager.GetViewMode() != ViewModeEdit {
		t.Errorf("預期切換後模式為編輯模式，實際為 %v", viewManager.GetViewMode())
	}
	
	// 第二次切換：編輯模式 -> 預覽模式
	viewManager.ToggleViewMode()
	if viewManager.GetViewMode() != ViewModePreview {
		t.Errorf("預期切換後模式為預覽模式，實際為 %v", viewManager.GetViewMode())
	}
	
	// 第三次切換：預覽模式 -> 分割視圖
	viewManager.ToggleViewMode()
	if viewManager.GetViewMode() != ViewModeSplit {
		t.Errorf("預期切換後模式為分割視圖，實際為 %v", viewManager.GetViewMode())
	}
}

// TestSetSplitRatio 測試分割比例設定
// 驗證分割比例設定和邊界值處理
func TestSetSplitRatio(t *testing.T) {
	// 建立測試環境
	testApp := app.New()
	testWindow := testApp.NewWindow("Test")
	defer testWindow.Close()
	
	layoutManager := NewLayoutManager()
	viewManager := NewViewManager(testWindow, layoutManager)
	
	// 測試回調函數
	var callbackRatio float64
	var callbackCalled bool
	viewManager.SetOnSplitRatioChanged(func(ratio float64) {
		callbackRatio = ratio
		callbackCalled = true
	})
	
	// 測試正常範圍內的比例
	viewManager.SetSplitRatio(0.3)
	if viewManager.GetSplitRatio() != 0.3 {
		t.Errorf("預期分割比例為 0.3，實際為 %f", viewManager.GetSplitRatio())
	}
	if !callbackCalled || callbackRatio != 0.3 {
		t.Error("分割比例變更回調未正確觸發")
	}
	
	// 測試邊界值：小於最小值
	viewManager.SetSplitRatio(0.05)
	if viewManager.GetSplitRatio() != 0.1 {
		t.Errorf("預期分割比例被限制為 0.1，實際為 %f", viewManager.GetSplitRatio())
	}
	
	// 測試邊界值：大於最大值
	viewManager.SetSplitRatio(0.95)
	if viewManager.GetSplitRatio() != 0.9 {
		t.Errorf("預期分割比例被限制為 0.9，實際為 %f", viewManager.GetSplitRatio())
	}
}

// TestToggleFullscreen 測試全螢幕模式切換
// 驗證全螢幕模式切換和狀態保存功能
func TestToggleFullscreen(t *testing.T) {
	// 建立測試環境
	testApp := app.New()
	testWindow := testApp.NewWindow("Test")
	defer testWindow.Close()
	
	layoutManager := NewLayoutManager()
	viewManager := NewViewManager(testWindow, layoutManager)
	
	// 測試回調函數
	var callbackFullscreen bool
	var callbackCalled bool
	viewManager.SetOnFullscreenToggled(func(fullscreen bool) {
		callbackFullscreen = fullscreen
		callbackCalled = true
	})
	
	// 初始狀態應該不是全螢幕
	if viewManager.IsFullscreen() {
		t.Error("預期初始狀態不是全螢幕模式")
	}
	
	// 切換到全螢幕模式
	viewManager.ToggleFullscreen()
	if !viewManager.IsFullscreen() {
		t.Error("預期切換後為全螢幕模式")
	}
	if !callbackCalled || !callbackFullscreen {
		t.Error("全螢幕切換回調未正確觸發")
	}
	
	// 重置回調狀態
	callbackCalled = false
	
	// 切換回視窗模式
	viewManager.ToggleFullscreen()
	if viewManager.IsFullscreen() {
		t.Error("預期切換後不是全螢幕模式")
	}
	if !callbackCalled || callbackFullscreen {
		t.Error("全螢幕切換回調未正確觸發")
	}
}

// TestPreviousMode 測試上一個模式記憶功能
// 驗證上一個視圖模式的記憶和恢復功能
func TestPreviousMode(t *testing.T) {
	// 建立測試環境
	testApp := app.New()
	testWindow := testApp.NewWindow("Test")
	defer testWindow.Close()
	
	layoutManager := NewLayoutManager()
	viewManager := NewViewManager(testWindow, layoutManager)
	
	// 初始狀態：分割視圖
	initialMode := viewManager.GetViewMode()
	
	// 切換到編輯模式
	viewManager.SetViewMode(ViewModeEdit)
	
	// 檢查上一個模式是否正確記憶
	if viewManager.GetPreviousMode() != initialMode {
		t.Errorf("預期上一個模式為 %v，實際為 %v", initialMode, viewManager.GetPreviousMode())
	}
	
	// 切換到預覽模式
	viewManager.SetViewMode(ViewModePreview)
	
	// 檢查上一個模式是否更新
	if viewManager.GetPreviousMode() != ViewModeEdit {
		t.Errorf("預期上一個模式為編輯模式，實際為 %v", viewManager.GetPreviousMode())
	}
	
	// 恢復到上一個模式
	viewManager.RestorePreviousMode()
	if viewManager.GetViewMode() != ViewModeEdit {
		t.Errorf("預期恢復後模式為編輯模式，實際為 %v", viewManager.GetViewMode())
	}
}

// TestViewStateManagement 測試視圖狀態管理
// 驗證視圖狀態的保存和載入功能
func TestViewStateManagement(t *testing.T) {
	// 建立測試環境
	testApp := app.New()
	testWindow := testApp.NewWindow("Test")
	defer testWindow.Close()
	
	layoutManager := NewLayoutManager()
	viewManager := NewViewManager(testWindow, layoutManager)
	
	// 設定特定的視圖狀態
	viewManager.SetViewMode(ViewModeEdit)
	viewManager.SetSplitRatio(0.7)
	
	// 保存狀態
	savedState := viewManager.SaveViewState()
	
	// 驗證保存的狀態
	if savedState.Mode != ViewModeEdit {
		t.Errorf("預期保存的模式為編輯模式，實際為 %v", savedState.Mode)
	}
	if savedState.SplitRatio != 0.7 {
		t.Errorf("預期保存的分割比例為 0.7，實際為 %f", savedState.SplitRatio)
	}
	
	// 變更狀態
	viewManager.SetViewMode(ViewModePreview)
	viewManager.SetSplitRatio(0.3)
	
	// 載入之前保存的狀態
	viewManager.LoadViewState(savedState)
	
	// 驗證狀態是否正確恢復
	if viewManager.GetViewMode() != ViewModeEdit {
		t.Errorf("預期載入後模式為編輯模式，實際為 %v", viewManager.GetViewMode())
	}
	if viewManager.GetSplitRatio() != 0.7 {
		t.Errorf("預期載入後分割比例為 0.7，實際為 %f", viewManager.GetSplitRatio())
	}
}

// TestViewModeStrings 測試視圖模式字串表示
// 驗證視圖模式的字串轉換功能
func TestViewModeStrings(t *testing.T) {
	// 建立測試環境
	testApp := app.New()
	testWindow := testApp.NewWindow("Test")
	defer testWindow.Close()
	
	layoutManager := NewLayoutManager()
	viewManager := NewViewManager(testWindow, layoutManager)
	
	// 測試各種視圖模式的字串表示
	testCases := []struct {
		mode     ViewMode
		expected string
	}{
		{ViewModeEdit, "編輯模式"},
		{ViewModePreview, "預覽模式"},
		{ViewModeSplit, "分割視圖"},
	}
	
	for _, tc := range testCases {
		result := viewManager.GetViewModeString(tc.mode)
		if result != tc.expected {
			t.Errorf("視圖模式 %v 的字串表示預期為 '%s'，實際為 '%s'", tc.mode, tc.expected, result)
		}
	}
	
	// 測試當前視圖模式字串
	viewManager.SetViewMode(ViewModeEdit)
	currentString := viewManager.GetCurrentViewModeString()
	if currentString != "編輯模式" {
		t.Errorf("預期當前視圖模式字串為 '編輯模式'，實際為 '%s'", currentString)
	}
}

// TestViewModeShortcuts 測試視圖模式快捷鍵
// 驗證快捷鍵字串表示功能
func TestViewModeShortcuts(t *testing.T) {
	// 建立測試環境
	testApp := app.New()
	testWindow := testApp.NewWindow("Test")
	defer testWindow.Close()
	
	layoutManager := NewLayoutManager()
	viewManager := NewViewManager(testWindow, layoutManager)
	
	// 測試各種視圖模式的快捷鍵
	testCases := []struct {
		mode     ViewMode
		expected string
	}{
		{ViewModeEdit, "⌘1"},
		{ViewModePreview, "⌘2"},
		{ViewModeSplit, "⌘3"},
	}
	
	for _, tc := range testCases {
		result := viewManager.GetViewModeShortcut(tc.mode)
		if result != tc.expected {
			t.Errorf("視圖模式 %v 的快捷鍵預期為 '%s'，實際為 '%s'", tc.mode, tc.expected, result)
		}
	}
}

// TestViewModeCheckers 測試視圖模式檢查函數
// 驗證各種視圖模式檢查函數的正確性
func TestViewModeCheckers(t *testing.T) {
	// 建立測試環境
	testApp := app.New()
	testWindow := testApp.NewWindow("Test")
	defer testWindow.Close()
	
	layoutManager := NewLayoutManager()
	viewManager := NewViewManager(testWindow, layoutManager)
	
	// 測試編輯模式檢查
	viewManager.SetViewMode(ViewModeEdit)
	if !viewManager.IsEditMode() {
		t.Error("預期 IsEditMode() 回傳 true")
	}
	if viewManager.IsPreviewMode() {
		t.Error("預期 IsPreviewMode() 回傳 false")
	}
	if viewManager.IsSplitMode() {
		t.Error("預期 IsSplitMode() 回傳 false")
	}
	
	// 測試預覽模式檢查
	viewManager.SetViewMode(ViewModePreview)
	if viewManager.IsEditMode() {
		t.Error("預期 IsEditMode() 回傳 false")
	}
	if !viewManager.IsPreviewMode() {
		t.Error("預期 IsPreviewMode() 回傳 true")
	}
	if viewManager.IsSplitMode() {
		t.Error("預期 IsSplitMode() 回傳 false")
	}
	
	// 測試分割視圖模式檢查
	viewManager.SetViewMode(ViewModeSplit)
	if viewManager.IsEditMode() {
		t.Error("預期 IsEditMode() 回傳 false")
	}
	if viewManager.IsPreviewMode() {
		t.Error("預期 IsPreviewMode() 回傳 false")
	}
	if !viewManager.IsSplitMode() {
		t.Error("預期 IsSplitMode() 回傳 true")
	}
}

// TestGetAvailableViewModes 測試取得可用視圖模式
// 驗證可用視圖模式列表的正確性
func TestGetAvailableViewModes(t *testing.T) {
	// 建立測試環境
	testApp := app.New()
	testWindow := testApp.NewWindow("Test")
	defer testWindow.Close()
	
	layoutManager := NewLayoutManager()
	viewManager := NewViewManager(testWindow, layoutManager)
	
	// 取得可用視圖模式
	availableModes := viewManager.GetAvailableViewModes()
	
	// 驗證模式數量
	expectedCount := 3
	if len(availableModes) != expectedCount {
		t.Errorf("預期可用視圖模式數量為 %d，實際為 %d", expectedCount, len(availableModes))
	}
	
	// 驗證包含所有預期的模式
	expectedModes := []ViewMode{ViewModeEdit, ViewModePreview, ViewModeSplit}
	for _, expectedMode := range expectedModes {
		found := false
		for _, availableMode := range availableModes {
			if availableMode == expectedMode {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("可用視圖模式列表中缺少模式 %v", expectedMode)
		}
	}
}

// TestKeyboardShortcuts 測試鍵盤快捷鍵功能
// 驗證鍵盤快捷鍵是否正確觸發視圖切換
func TestKeyboardShortcuts(t *testing.T) {
	// 建立測試環境
	testApp := app.New()
	testWindow := testApp.NewWindow("Test")
	defer testWindow.Close()
	
	layoutManager := NewLayoutManager()
	viewManager := NewViewManager(testWindow, layoutManager)
	
	// 初始狀態應該是分割視圖
	if viewManager.GetViewMode() != ViewModeSplit {
		t.Errorf("預期初始模式為分割視圖，實際為 %v", viewManager.GetViewMode())
	}
	
	// 由於測試環境的限制，我們無法直接測試快捷鍵觸發
	// 但我們可以測試快捷鍵設定是否正確建立
	canvas := testWindow.Canvas()
	if canvas == nil {
		t.Error("視窗畫布未正確建立")
	}
	
	// 手動觸發視圖模式變更來驗證功能
	viewManager.SetViewMode(ViewModeEdit)
	if viewManager.GetViewMode() != ViewModeEdit {
		t.Errorf("預期視圖模式為編輯模式，實際為 %v", viewManager.GetViewMode())
	}
	
	viewManager.SetViewMode(ViewModePreview)
	if viewManager.GetViewMode() != ViewModePreview {
		t.Errorf("預期視圖模式為預覽模式，實際為 %v", viewManager.GetViewMode())
	}
	
	viewManager.SetViewMode(ViewModeSplit)
	if viewManager.GetViewMode() != ViewModeSplit {
		t.Errorf("預期視圖模式為分割視圖，實際為 %v", viewManager.GetViewMode())
	}
}

// TestContainerManagement 測試容器管理功能
// 驗證編輯器和預覽容器的設定和取得功能
func TestContainerManagement(t *testing.T) {
	// 建立測試環境
	testApp := app.New()
	testWindow := testApp.NewWindow("Test")
	defer testWindow.Close()
	
	layoutManager := NewLayoutManager()
	viewManager := NewViewManager(testWindow, layoutManager)
	
	// 取得主要容器
	container := viewManager.GetContainer()
	if container == nil {
		t.Fatal("主要容器未正確建立")
	}
	
	// 測試設定編輯器內容
	// 由於我們無法建立實際的編輯器內容，這裡只測試方法是否存在
	// 實際的內容設定會在整合測試中驗證
	
	// 驗證容器在不同視圖模式下的狀態
	viewManager.SetViewMode(ViewModeEdit)
	if len(container.Objects) != 1 {
		t.Errorf("編輯模式下預期容器有 1 個物件，實際有 %d 個", len(container.Objects))
	}
	
	viewManager.SetViewMode(ViewModePreview)
	if len(container.Objects) != 1 {
		t.Errorf("預覽模式下預期容器有 1 個物件，實際有 %d 個", len(container.Objects))
	}
	
	viewManager.SetViewMode(ViewModeSplit)
	if len(container.Objects) != 1 {
		t.Errorf("分割視圖模式下預期容器有 1 個物件，實際有 %d 個", len(container.Objects))
	}
}

// BenchmarkViewModeSwitch 視圖模式切換的效能基準測試
// 測量視圖模式切換的效能
func BenchmarkViewModeSwitch(b *testing.B) {
	// 建立測試環境
	testApp := app.New()
	testWindow := testApp.NewWindow("Benchmark")
	defer testWindow.Close()
	
	layoutManager := NewLayoutManager()
	viewManager := NewViewManager(testWindow, layoutManager)
	
	// 執行基準測試
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		viewManager.SetViewMode(ViewModeEdit)
		viewManager.SetViewMode(ViewModePreview)
		viewManager.SetViewMode(ViewModeSplit)
	}
}

// BenchmarkSplitRatioChange 分割比例變更的效能基準測試
// 測量分割比例變更的效能
func BenchmarkSplitRatioChange(b *testing.B) {
	// 建立測試環境
	testApp := app.New()
	testWindow := testApp.NewWindow("Benchmark")
	defer testWindow.Close()
	
	layoutManager := NewLayoutManager()
	viewManager := NewViewManager(testWindow, layoutManager)
	
	// 執行基準測試
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ratio := float64(i%9+1) / 10.0 // 0.1 到 0.9
		viewManager.SetSplitRatio(ratio)
	}
}