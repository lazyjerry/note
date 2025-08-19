// Package ui 包含使用者介面相關的元件和視窗管理
// 本檔案包含佈局管理器的單元測試
package ui

import (
	"testing"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// TestNewLayoutManager 測試佈局管理器的建立
// 驗證佈局管理器是否正確初始化所有元件和預設狀態
func TestNewLayoutManager(t *testing.T) {
	// 建立佈局管理器
	layoutManager := NewLayoutManager()
	
	// 驗證佈局管理器不為 nil
	if layoutManager == nil {
		t.Fatal("佈局管理器不應該為 nil")
	}
	
	// 驗證主要容器已建立
	if layoutManager.mainContainer == nil {
		t.Error("主要容器不應該為 nil")
	}
	
	// 驗證工具欄已建立
	if layoutManager.quickToolbar == nil {
		t.Error("快速工具欄不應該為 nil")
	}
	
	// 驗證面板容器已建立
	if layoutManager.sidebarPanel == nil {
		t.Error("側邊欄面板不應該為 nil")
	}
	
	if layoutManager.noteListPanel == nil {
		t.Error("筆記列表面板不應該為 nil")
	}
	
	if layoutManager.editorPanel == nil {
		t.Error("編輯器面板不應該為 nil")
	}
	
	// 驗證預設狀態
	if !layoutManager.sidebarVisible {
		t.Error("側邊欄預設應該可見")
	}
	
	if !layoutManager.noteListVisible {
		t.Error("筆記列表預設應該可見")
	}
	
	// 驗證預設寬度比例
	expectedSidebarWidth := 0.2
	if layoutManager.sidebarWidth != expectedSidebarWidth {
		t.Errorf("側邊欄寬度預設應該為 %f，實際為 %f", expectedSidebarWidth, layoutManager.sidebarWidth)
	}
	
	expectedNoteListWidth := 0.25
	if layoutManager.noteListWidth != expectedNoteListWidth {
		t.Errorf("筆記列表寬度預設應該為 %f，實際為 %f", expectedNoteListWidth, layoutManager.noteListWidth)
	}
}

// TestLayoutManagerToggleSidebar 測試側邊欄切換功能
// 驗證側邊欄的顯示/隱藏切換是否正常工作
func TestLayoutManagerToggleSidebar(t *testing.T) {
	layoutManager := NewLayoutManager()
	
	// 記錄初始狀態
	initialVisible := layoutManager.sidebarVisible
	
	// 切換側邊欄
	layoutManager.ToggleSidebar()
	
	// 驗證狀態已變更
	if layoutManager.sidebarVisible == initialVisible {
		t.Error("側邊欄可見性狀態應該已變更")
	}
	
	// 再次切換
	layoutManager.ToggleSidebar()
	
	// 驗證狀態已恢復
	if layoutManager.sidebarVisible != initialVisible {
		t.Error("側邊欄可見性狀態應該已恢復到初始狀態")
	}
}

// TestLayoutManagerToggleNoteList 測試筆記列表切換功能
// 驗證筆記列表的顯示/隱藏切換是否正常工作
func TestLayoutManagerToggleNoteList(t *testing.T) {
	layoutManager := NewLayoutManager()
	
	// 記錄初始狀態
	initialVisible := layoutManager.noteListVisible
	
	// 切換筆記列表
	layoutManager.ToggleNoteList()
	
	// 驗證狀態已變更
	if layoutManager.noteListVisible == initialVisible {
		t.Error("筆記列表可見性狀態應該已變更")
	}
	
	// 再次切換
	layoutManager.ToggleNoteList()
	
	// 驗證狀態已恢復
	if layoutManager.noteListVisible != initialVisible {
		t.Error("筆記列表可見性狀態應該已恢復到初始狀態")
	}
}

// TestLayoutManagerSetSidebarWidth 測試設定側邊欄寬度
// 驗證側邊欄寬度設定是否正確處理邊界值和有效範圍
func TestLayoutManagerSetSidebarWidth(t *testing.T) {
	layoutManager := NewLayoutManager()
	
	// 測試有效寬度
	validWidth := 0.3
	layoutManager.SetSidebarWidth(validWidth)
	
	if layoutManager.sidebarWidth != validWidth {
		t.Errorf("側邊欄寬度應該為 %f，實際為 %f", validWidth, layoutManager.sidebarWidth)
	}
	
	// 測試最小邊界值
	layoutManager.SetSidebarWidth(0.05) // 小於最小值 0.1
	if layoutManager.sidebarWidth != 0.1 {
		t.Errorf("側邊欄寬度應該被限制為最小值 0.1，實際為 %f", layoutManager.sidebarWidth)
	}
	
	// 測試最大邊界值
	layoutManager.SetSidebarWidth(0.6) // 大於最大值 0.5
	if layoutManager.sidebarWidth != 0.5 {
		t.Errorf("側邊欄寬度應該被限制為最大值 0.5，實際為 %f", layoutManager.sidebarWidth)
	}
}

// TestLayoutManagerSetNoteListWidth 測試設定筆記列表寬度
// 驗證筆記列表寬度設定是否正確處理邊界值和有效範圍
func TestLayoutManagerSetNoteListWidth(t *testing.T) {
	layoutManager := NewLayoutManager()
	
	// 測試有效寬度
	validWidth := 0.4
	layoutManager.SetNoteListWidth(validWidth)
	
	if layoutManager.noteListWidth != validWidth {
		t.Errorf("筆記列表寬度應該為 %f，實際為 %f", validWidth, layoutManager.noteListWidth)
	}
	
	// 測試最小邊界值
	layoutManager.SetNoteListWidth(0.05) // 小於最小值 0.1
	if layoutManager.noteListWidth != 0.1 {
		t.Errorf("筆記列表寬度應該被限制為最小值 0.1，實際為 %f", layoutManager.noteListWidth)
	}
	
	// 測試最大邊界值
	layoutManager.SetNoteListWidth(0.9) // 大於最大值 0.8
	if layoutManager.noteListWidth != 0.8 {
		t.Errorf("筆記列表寬度應該被限制為最大值 0.8，實際為 %f", layoutManager.noteListWidth)
	}
}

// TestLayoutManagerSetContent 測試設定面板內容
// 驗證各面板內容設定是否正確更新
func TestLayoutManagerSetContent(t *testing.T) {
	layoutManager := NewLayoutManager()
	
	// 建立測試內容
	testSidebarContent := container.NewVBox(widget.NewLabel("測試側邊欄"))
	testNoteListContent := container.NewVBox(widget.NewLabel("測試筆記列表"))
	testEditorContent := container.NewVBox(widget.NewLabel("測試編輯器"))
	testStatusContent := container.NewHBox(widget.NewLabel("測試狀態欄"))
	
	// 設定各面板內容
	layoutManager.SetSidebarContent(testSidebarContent)
	layoutManager.SetNoteListContent(testNoteListContent)
	layoutManager.SetEditorContent(testEditorContent)
	layoutManager.SetStatusBarContent(testStatusContent)
	
	// 驗證內容已設定（檢查面板是否包含測試內容）
	if len(layoutManager.sidebarPanel.Objects) == 0 {
		t.Error("側邊欄面板應該包含內容")
	}
	
	if len(layoutManager.noteListPanel.Objects) == 0 {
		t.Error("筆記列表面板應該包含內容")
	}
	
	if len(layoutManager.editorPanel.Objects) == 0 {
		t.Error("編輯器面板應該包含內容")
	}
	
	if len(layoutManager.bottomBar.Objects) == 0 {
		t.Error("狀態欄應該包含內容")
	}
}

// TestLayoutManagerCompactMode 測試緊湊模式
// 驗證緊湊模式的切換和相關設定調整
func TestLayoutManagerCompactMode(t *testing.T) {
	layoutManager := NewLayoutManager()
	
	// 驗證初始狀態不是緊湊模式
	if layoutManager.IsCompactMode() {
		t.Error("初始狀態不應該是緊湊模式")
	}
	
	// 記錄初始寬度
	initialSidebarWidth := layoutManager.sidebarWidth
	initialNoteListWidth := layoutManager.noteListWidth
	
	// 啟用緊湊模式
	layoutManager.SetCompactMode(true)
	
	// 驗證緊湊模式已啟用
	if !layoutManager.IsCompactMode() {
		t.Error("緊湊模式應該已啟用")
	}
	
	// 驗證寬度已調整
	if layoutManager.sidebarWidth >= initialSidebarWidth {
		t.Error("緊湊模式下側邊欄寬度應該減少")
	}
	
	if layoutManager.noteListWidth >= initialNoteListWidth {
		t.Error("緊湊模式下筆記列表寬度應該減少")
	}
	
	// 停用緊湊模式
	layoutManager.SetCompactMode(false)
	
	// 驗證緊湊模式已停用
	if layoutManager.IsCompactMode() {
		t.Error("緊湊模式應該已停用")
	}
}

// TestLayoutManagerResizeToWindow 測試視窗大小調整
// 驗證佈局管理器是否正確響應視窗大小變更
func TestLayoutManagerResizeToWindow(t *testing.T) {
	layoutManager := NewLayoutManager()
	
	// 測試小視窗大小
	smallSize := fyne.NewSize(800, 600)
	layoutManager.ResizeToWindow(smallSize)
	
	// 驗證是否啟用緊湊模式
	if !layoutManager.IsCompactMode() {
		t.Error("小視窗應該啟用緊湊模式")
	}
	
	// 測試大視窗大小
	largeSize := fyne.NewSize(1600, 1200)
	layoutManager.ResizeToWindow(largeSize)
	
	// 驗證是否停用緊湊模式
	if layoutManager.IsCompactMode() {
		t.Error("大視窗應該停用緊湊模式")
	}
}

// TestLayoutManagerSaveLoadState 測試佈局狀態保存和載入
// 驗證佈局狀態的保存和載入功能是否正確
func TestLayoutManagerSaveLoadState(t *testing.T) {
	layoutManager := NewLayoutManager()
	
	// 修改一些設定
	layoutManager.SetSidebarWidth(0.3)
	layoutManager.SetNoteListWidth(0.35)
	layoutManager.ToggleSidebar()
	layoutManager.SetCompactMode(true)
	
	// 保存狀態
	state := layoutManager.SaveLayoutState()
	
	// 驗證狀態包含預期的鍵
	expectedKeys := []string{"sidebar_visible", "notelist_visible", "sidebar_width", "notelist_width", "compact_mode"}
	for _, key := range expectedKeys {
		if _, exists := state[key]; !exists {
			t.Errorf("保存的狀態應該包含鍵: %s", key)
		}
	}
	
	// 建立新的佈局管理器
	newLayoutManager := NewLayoutManager()
	
	// 載入狀態
	newLayoutManager.LoadLayoutState(state)
	
	// 驗證狀態已正確載入
	if newLayoutManager.sidebarWidth != 0.3 {
		t.Errorf("載入的側邊欄寬度應該為 0.3，實際為 %f", newLayoutManager.sidebarWidth)
	}
	
	if newLayoutManager.noteListWidth != 0.35 {
		t.Errorf("載入的筆記列表寬度應該為 0.35，實際為 %f", newLayoutManager.noteListWidth)
	}
	
	if newLayoutManager.sidebarVisible {
		t.Error("載入的側邊欄可見性應該為 false")
	}
	
	if !newLayoutManager.compactMode {
		t.Error("載入的緊湊模式應該為 true")
	}
}

// TestLayoutManagerCallbacks 測試回調函數
// 驗證佈局變更和面板大小變更回調是否正確觸發
func TestLayoutManagerCallbacks(t *testing.T) {
	layoutManager := NewLayoutManager()
	
	// 設定回調函數
	var layoutChangedCalled bool
	var panelResizedCalled bool
	var lastLayoutChange string
	var lastPanelChange string
	var lastPanelSize float64
	
	layoutManager.SetOnLayoutChanged(func(layout string) {
		layoutChangedCalled = true
		lastLayoutChange = layout
	})
	
	layoutManager.SetOnPanelResized(func(panel string, size float64) {
		panelResizedCalled = true
		lastPanelChange = panel
		lastPanelSize = size
	})
	
	// 觸發佈局變更
	layoutManager.ToggleSidebar()
	
	// 驗證佈局變更回調已觸發
	if !layoutChangedCalled {
		t.Error("佈局變更回調應該已觸發")
	}
	
	if lastLayoutChange != "sidebar_hidden" {
		t.Errorf("佈局變更類型應該為 'sidebar_hidden'，實際為 '%s'", lastLayoutChange)
	}
	
	// 觸發面板大小變更
	layoutManager.SetSidebarWidth(0.4)
	
	// 驗證面板大小變更回調已觸發
	if !panelResizedCalled {
		t.Error("面板大小變更回調應該已觸發")
	}
	
	if lastPanelChange != "sidebar" {
		t.Errorf("面板變更類型應該為 'sidebar'，實際為 '%s'", lastPanelChange)
	}
	
	if lastPanelSize != 0.4 {
		t.Errorf("面板大小應該為 0.4，實際為 %f", lastPanelSize)
	}
}