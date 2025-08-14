// Package ui 包含使用者介面相關的元件和視窗管理測試
// 測試 EditorWithPreview 的建立、初始化和整合功能
package ui

import (
	"testing"                           // Go 標準測試套件
	"mac-notebook-app/internal/models"  // 引入資料模型
)

// TestNewEditorWithPreview 測試整合編輯器和預覽面板的建立和初始化
// 驗證複合元件是否正確建立並包含所有必要的子元件
//
// 測試項目：
// 1. 複合元件實例是否正確建立
// 2. 編輯器和預覽面板子元件是否正確初始化
// 3. 服務依賴是否正確設定
// 4. 初始狀態是否正確
func TestNewEditorWithPreview(t *testing.T) {
	// 建立模擬編輯器服務
	mockService := newMockEditorServiceForPreview()
	
	// 建立整合編輯器和預覽面板實例
	ewp := NewEditorWithPreview(mockService)
	
	// 驗證複合元件實例不為 nil
	if ewp == nil {
		t.Fatal("NewEditorWithPreview 應該回傳有效的 EditorWithPreview 實例")
	}
	
	// 驗證主要容器
	if ewp.container == nil {
		t.Error("複合元件的主要容器不應該為 nil")
	}
	
	// 驗證分割容器
	if ewp.splitContainer == nil {
		t.Error("複合元件的分割容器不應該為 nil")
	}
	
	// 驗證編輯器子元件
	if ewp.editor == nil {
		t.Error("編輯器子元件不應該為 nil")
	}
	
	// 驗證預覽面板子元件
	if ewp.preview == nil {
		t.Error("預覽面板子元件不應該為 nil")
	}
	
	// 驗證服務依賴
	if ewp.editorService == nil {
		t.Error("編輯器服務不應該為 nil")
	}
	
	// 驗證初始狀態
	if !ewp.previewVisible {
		t.Error("初始狀態預覽面板應該為可見")
	}
	
	if ewp.splitRatio != 0.5 {
		t.Errorf("初始分割比例應該是 0.5，但得到 %f", ewp.splitRatio)
	}
}

// TestEditorWithPreviewCreateNewNote 測試建立新筆記功能
// 驗證複合元件是否能正確建立和載入新筆記
//
// 測試項目：
// 1. 新筆記是否正確建立
// 2. 編輯器和預覽面板是否正確更新
// 3. 狀態是否正確設定
func TestEditorWithPreviewCreateNewNote(t *testing.T) {
	// 建立模擬編輯器服務
	mockService := newMockEditorServiceForPreview()
	
	// 建立整合編輯器和預覽面板實例
	ewp := NewEditorWithPreview(mockService)
	
	// 建立新筆記
	testTitle := "測試筆記標題"
	err := ewp.CreateNewNote(testTitle)
	
	// 驗證沒有錯誤
	if err != nil {
		t.Errorf("建立新筆記不應該產生錯誤，但得到: %s", err.Error())
	}
	
	// 驗證編輯器狀態
	if ewp.editor.GetCurrentNote() == nil {
		t.Error("建立新筆記後編輯器應該有當前筆記")
	}
	
	if ewp.editor.GetTitle() != testTitle {
		t.Errorf("筆記標題應該是 '%s'，但得到 '%s'", testTitle, ewp.editor.GetTitle())
	}
	
	// 驗證複合元件狀態
	if ewp.GetTitle() != testTitle {
		t.Errorf("複合元件標題應該是 '%s'，但得到 '%s'", testTitle, ewp.GetTitle())
	}
}

// TestEditorWithPreviewLoadNote 測試載入筆記功能
// 驗證複合元件是否能正確載入現有筆記
//
// 測試項目：
// 1. 筆記是否正確載入到編輯器
// 2. 預覽面板是否正確更新
// 3. 內容同步是否正確
func TestEditorWithPreviewLoadNote(t *testing.T) {
	// 建立模擬編輯器服務
	mockService := newMockEditorServiceForPreview()
	
	// 建立整合編輯器和預覽面板實例
	ewp := NewEditorWithPreview(mockService)
	
	// 建立測試筆記
	testNote := &models.Note{
		ID:      "test-id",
		Title:   "測試筆記",
		Content: "# 測試標題\n\n這是測試內容。",
	}
	
	// 載入筆記
	ewp.LoadNote(testNote)
	
	// 驗證編輯器內容
	if ewp.GetContent() != testNote.Content {
		t.Errorf("編輯器內容應該是 '%s'，但得到 '%s'", testNote.Content, ewp.GetContent())
	}
	
	// 驗證預覽面板內容
	if ewp.preview.GetCurrentContent() != testNote.Content {
		t.Errorf("預覽面板內容應該是 '%s'，但得到 '%s'", testNote.Content, ewp.preview.GetCurrentContent())
	}
	
	// 驗證標題
	if ewp.GetTitle() != testNote.Title {
		t.Errorf("標題應該是 '%s'，但得到 '%s'", testNote.Title, ewp.GetTitle())
	}
}

// TestEditorWithPreviewContentSync 測試內容同步功能
// 驗證編輯器和預覽面板之間的內容同步
//
// 測試項目：
// 1. 編輯器內容變更是否自動更新預覽
// 2. 自動刷新功能是否正確工作
// 3. 內容同步的即時性
func TestEditorWithPreviewContentSync(t *testing.T) {
	// 建立模擬編輯器服務
	mockService := newMockEditorServiceForPreview()
	
	// 建立整合編輯器和預覽面板實例
	ewp := NewEditorWithPreview(mockService)
	
	// 建立新筆記
	err := ewp.CreateNewNote("測試筆記")
	if err != nil {
		t.Fatalf("建立新筆記失敗: %s", err.Error())
	}
	
	// 設定內容
	testContent := "# 新標題\n\n這是新內容。"
	ewp.SetContent(testContent)
	
	// 驗證編輯器內容
	if ewp.GetContent() != testContent {
		t.Errorf("編輯器內容應該是 '%s'，但得到 '%s'", testContent, ewp.GetContent())
	}
	
	// 驗證預覽面板內容（應該自動同步）
	if ewp.preview.GetCurrentContent() != testContent {
		t.Errorf("預覽面板內容應該自動同步為 '%s'，但得到 '%s'", testContent, ewp.preview.GetCurrentContent())
	}
}

// TestEditorWithPreviewSaveNote 測試保存筆記功能
// 驗證複合元件的筆記保存功能
//
// 測試項目：
// 1. 筆記保存是否正確
// 2. 保存狀態是否正確更新
// 3. 錯誤處理是否正確
func TestEditorWithPreviewSaveNote(t *testing.T) {
	// 建立模擬編輯器服務
	mockService := newMockEditorServiceForPreview()
	
	// 建立整合編輯器和預覽面板實例
	ewp := NewEditorWithPreview(mockService)
	
	// 建立新筆記
	err := ewp.CreateNewNote("測試筆記")
	if err != nil {
		t.Fatalf("建立新筆記失敗: %s", err.Error())
	}
	
	// 修改內容
	testContent := "修改後的內容"
	ewp.SetContent(testContent)
	
	// 手動設定修改狀態（因為 SetContent 會重置修改狀態）
	ewp.editor.isModified = true
	
	// 保存筆記
	err = ewp.SaveNote()
	
	// 驗證沒有錯誤
	if err != nil {
		t.Errorf("保存筆記不應該產生錯誤，但得到: %s", err.Error())
	}
	
	// 驗證修改狀態
	if ewp.IsModified() {
		t.Error("保存後修改狀態應該為 false")
	}
}

// TestEditorWithPreviewPreviewToggle 測試預覽切換功能
// 驗證預覽面板的顯示和隱藏功能
//
// 測試項目：
// 1. 預覽切換是否正確
// 2. 佈局更新是否正確
// 3. 狀態管理是否正確
func TestEditorWithPreviewPreviewToggle(t *testing.T) {
	// 建立模擬編輯器服務
	mockService := newMockEditorServiceForPreview()
	
	// 建立整合編輯器和預覽面板實例
	ewp := NewEditorWithPreview(mockService)
	
	// 驗證初始狀態
	if !ewp.IsPreviewVisible() {
		t.Error("初始狀態預覽應該為可見")
	}
	
	// 切換預覽（隱藏）
	ewp.TogglePreview()
	
	if ewp.IsPreviewVisible() {
		t.Error("切換後預覽應該為隱藏")
	}
	
	// 切換預覽（顯示）
	ewp.TogglePreview()
	
	if !ewp.IsPreviewVisible() {
		t.Error("再次切換後預覽應該為可見")
	}
	
	// 測試直接設定可見性
	ewp.SetPreviewVisible(false)
	
	if ewp.IsPreviewVisible() {
		t.Error("設定後預覽應該為隱藏")
	}
	
	ewp.SetPreviewVisible(true)
	
	if !ewp.IsPreviewVisible() {
		t.Error("設定後預覽應該為可見")
	}
}

// TestEditorWithPreviewSplitRatio 測試分割比例功能
// 驗證編輯器和預覽面板的分割比例設定
//
// 測試項目：
// 1. 分割比例設定是否正確
// 2. 分割比例查詢是否正確
// 3. 邊界值處理是否正確
func TestEditorWithPreviewSplitRatio(t *testing.T) {
	// 建立模擬編輯器服務
	mockService := newMockEditorServiceForPreview()
	
	// 建立整合編輯器和預覽面板實例
	ewp := NewEditorWithPreview(mockService)
	
	// 驗證初始分割比例
	if ewp.GetSplitRatio() != 0.5 {
		t.Errorf("初始分割比例應該是 0.5，但得到 %f", ewp.GetSplitRatio())
	}
	
	// 測試設定分割比例
	testRatio := 0.7
	ewp.SetSplitRatio(testRatio)
	
	if ewp.GetSplitRatio() != testRatio {
		t.Errorf("分割比例應該是 %f，但得到 %f", testRatio, ewp.GetSplitRatio())
	}
	
	// 測試邊界值處理
	ewp.SetSplitRatio(-0.1) // 小於 0
	if ewp.GetSplitRatio() != 0.0 {
		t.Errorf("小於 0 的比例應該被設為 0.0，但得到 %f", ewp.GetSplitRatio())
	}
	
	ewp.SetSplitRatio(1.1) // 大於 1
	if ewp.GetSplitRatio() != 1.0 {
		t.Errorf("大於 1 的比例應該被設為 1.0，但得到 %f", ewp.GetSplitRatio())
	}
}

// TestEditorWithPreviewAutoRefresh 測試自動刷新功能
// 驗證預覽面板的自動刷新功能
//
// 測試項目：
// 1. 自動刷新設定是否正確
// 2. 自動刷新狀態查詢是否正確
// 3. 手動刷新功能是否正確
func TestEditorWithPreviewAutoRefresh(t *testing.T) {
	// 建立模擬編輯器服務
	mockService := newMockEditorServiceForPreview()
	
	// 建立整合編輯器和預覽面板實例
	ewp := NewEditorWithPreview(mockService)
	
	// 驗證初始自動刷新狀態
	if !ewp.IsAutoRefreshEnabled() {
		t.Error("初始狀態應該啟用自動刷新")
	}
	
	// 測試停用自動刷新
	ewp.SetAutoRefresh(false)
	
	if ewp.IsAutoRefreshEnabled() {
		t.Error("設定後自動刷新應該被停用")
	}
	
	// 測試啟用自動刷新
	ewp.SetAutoRefresh(true)
	
	if !ewp.IsAutoRefreshEnabled() {
		t.Error("設定後自動刷新應該被啟用")
	}
	
	// 測試手動刷新
	testContent := "# 測試內容"
	ewp.SetContent(testContent)
	ewp.RefreshPreview()
	
	// 驗證預覽內容已更新
	if ewp.preview.GetCurrentContent() != testContent {
		t.Errorf("手動刷新後預覽內容應該是 '%s'，但得到 '%s'", testContent, ewp.preview.GetCurrentContent())
	}
}

// TestEditorWithPreviewClearContent 測試清空內容功能
// 驗證複合元件的內容清空功能
//
// 測試項目：
// 1. 編輯器內容是否正確清空
// 2. 預覽面板內容是否正確清空
// 3. 狀態是否正確重置
func TestEditorWithPreviewClearContent(t *testing.T) {
	// 建立模擬編輯器服務
	mockService := newMockEditorServiceForPreview()
	
	// 建立整合編輯器和預覽面板實例
	ewp := NewEditorWithPreview(mockService)
	
	// 建立新筆記並添加內容
	err := ewp.CreateNewNote("測試筆記")
	if err != nil {
		t.Fatalf("建立新筆記失敗: %s", err.Error())
	}
	
	testContent := "# 測試內容"
	ewp.SetContent(testContent)
	
	// 驗證內容已設定
	if ewp.GetContent() != testContent {
		t.Error("內容設定失敗")
	}
	
	// 清空內容
	ewp.Clear()
	
	// 驗證編輯器內容已清空
	if ewp.GetContent() != "" {
		t.Error("清空後編輯器內容應該為空")
	}
	
	// 驗證預覽面板內容已清空
	if ewp.preview.GetCurrentContent() != "" {
		t.Error("清空後預覽面板內容應該為空")
	}
	
	// 驗證當前筆記已清空
	if ewp.GetCurrentNote() != nil {
		t.Error("清空後當前筆記應該為 nil")
	}
}

// TestEditorWithPreviewCallbacks 測試回調函數功能
// 驗證複合元件的回調函數設定和觸發
//
// 測試項目：
// 1. 內容變更回調是否正確觸發
// 2. 保存請求回調是否正確觸發
// 3. 字數變更回調是否正確觸發
// 4. 預覽切換回調是否正確觸發
func TestEditorWithPreviewCallbacks(t *testing.T) {
	// 建立模擬編輯器服務
	mockService := newMockEditorServiceForPreview()
	
	// 建立整合編輯器和預覽面板實例
	ewp := NewEditorWithPreview(mockService)
	
	// 設定回調函數
	var wordCountChangedCalled bool
	var previewToggledCalled bool
	var lastWordCount int
	var lastPreviewVisible bool
	
	ewp.SetOnWordCountChanged(func(count int) {
		wordCountChangedCalled = true
		lastWordCount = count
	})
	
	ewp.SetOnPreviewToggled(func(visible bool) {
		previewToggledCalled = true
		lastPreviewVisible = visible
	})
	
	// 建立新筆記以觸發回調
	err := ewp.CreateNewNote("測試筆記")
	if err != nil {
		t.Fatalf("建立新筆記失敗: %s", err.Error())
	}
	
	// 驗證字數變更回調
	if !wordCountChangedCalled {
		t.Error("建立新筆記應該觸發字數變更回調")
	}
	
	if lastWordCount != 0 {
		t.Errorf("新筆記的字數應該是 0，但得到 %d", lastWordCount)
	}
	
	// 測試預覽切換回調
	ewp.TogglePreview()
	
	if !previewToggledCalled {
		t.Error("預覽切換應該觸發預覽切換回調")
	}
	
	if lastPreviewVisible {
		t.Error("切換後預覽應該為隱藏")
	}
}

// TestEditorWithPreviewGetContainer 測試容器取得功能
// 驗證複合元件是否能正確回傳主要容器
//
// 測試項目：
// 1. GetContainer 是否回傳正確的容器實例
// 2. 回傳的容器是否與內部容器相同
func TestEditorWithPreviewGetContainer(t *testing.T) {
	// 建立模擬編輯器服務
	mockService := newMockEditorServiceForPreview()
	
	// 建立整合編輯器和預覽面板實例
	ewp := NewEditorWithPreview(mockService)
	
	// 取得容器
	container := ewp.GetContainer()
	
	// 驗證容器不為 nil
	if container == nil {
		t.Error("GetContainer 應該回傳有效的容器實例")
	}
	
	// 驗證回傳的容器與內部容器相同
	if container != ewp.container {
		t.Error("GetContainer 應該回傳與內部容器相同的實例")
	}
}

// TestEditorWithPreviewSubComponents 測試子元件存取功能
// 驗證複合元件是否能正確提供子元件的存取
//
// 測試項目：
// 1. GetEditor 是否回傳正確的編輯器實例
// 2. GetPreview 是否回傳正確的預覽面板實例
func TestEditorWithPreviewSubComponents(t *testing.T) {
	// 建立模擬編輯器服務
	mockService := newMockEditorServiceForPreview()
	
	// 建立整合編輯器和預覽面板實例
	ewp := NewEditorWithPreview(mockService)
	
	// 取得編輯器子元件
	editor := ewp.GetEditor()
	
	// 驗證編輯器不為 nil
	if editor == nil {
		t.Error("GetEditor 應該回傳有效的編輯器實例")
	}
	
	// 驗證回傳的編輯器與內部編輯器相同
	if editor != ewp.editor {
		t.Error("GetEditor 應該回傳與內部編輯器相同的實例")
	}
	
	// 取得預覽面板子元件
	preview := ewp.GetPreview()
	
	// 驗證預覽面板不為 nil
	if preview == nil {
		t.Error("GetPreview 應該回傳有效的預覽面板實例")
	}
	
	// 驗證回傳的預覽面板與內部預覽面板相同
	if preview != ewp.preview {
		t.Error("GetPreview 應該回傳與內部預覽面板相同的實例")
	}
}