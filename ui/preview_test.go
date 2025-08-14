// Package ui 包含使用者介面相關的元件和視窗管理測試
// 測試 MarkdownPreview 的建立、初始化和預覽功能
package ui

import (
	"testing"                           // Go 標準測試套件
	"mac-notebook-app/internal/models"  // 引入資料模型
	"mac-notebook-app/internal/services" // 引入服務層
)

// mockEditorServiceForPreview 模擬編輯器服務，用於預覽測試
// 實作 EditorService 介面的 PreviewMarkdown 方法
type mockEditorServiceForPreview struct {
	mockEditorService // 繼承基本的模擬服務
}

// PreviewMarkdown 模擬 Markdown 預覽功能
// 參數：content（Markdown 內容）
// 回傳：HTML 字串
func (m *mockEditorServiceForPreview) PreviewMarkdown(content string) string {
	// 簡單的 Markdown 到 HTML 轉換模擬
	if content == "" {
		return ""
	}
	
	// 模擬一些基本的 Markdown 轉換
	html := content
	if content == "# 測試標題" {
		html = "<h1>測試標題</h1>"
	} else if content == "**粗體文字**" {
		html = "<p><strong>粗體文字</strong></p>"
	} else {
		html = "<p>" + content + "</p>"
	}
	
	return html
}

// newMockEditorServiceForPreview 建立模擬編輯器服務實例
// 回傳：模擬編輯器服務實例
func newMockEditorServiceForPreview() services.EditorService {
	return &mockEditorServiceForPreview{
		mockEditorService: mockEditorService{
			notes: make(map[string]*models.Note),
		},
	}
}

// TestNewMarkdownPreview 測試 Markdown 預覽面板的建立和初始化
// 驗證預覽面板是否正確建立並包含所有必要的 UI 元件
//
// 測試項目：
// 1. 預覽面板實例是否正確建立
// 2. 所有 UI 元件是否正確初始化
// 3. 服務依賴是否正確設定
// 4. 初始狀態是否正確
func TestNewMarkdownPreview(t *testing.T) {
	// 建立模擬編輯器服務
	mockService := newMockEditorServiceForPreview()
	
	// 建立 Markdown 預覽面板實例
	preview := NewMarkdownPreview(mockService)
	
	// 驗證預覽面板實例不為 nil
	if preview == nil {
		t.Fatal("NewMarkdownPreview 應該回傳有效的 MarkdownPreview 實例")
	}
	
	// 驗證主要容器
	if preview.container == nil {
		t.Error("預覽面板的主要容器不應該為 nil")
	}
	
	// 驗證工具欄
	if preview.toolbar == nil {
		t.Error("預覽面板的工具欄不應該為 nil")
	}
	
	// 驗證預覽區域
	if preview.previewArea == nil {
		t.Error("預覽面板的預覽區域不應該為 nil")
	}
	
	// 驗證狀態標籤
	if preview.statusLabel == nil {
		t.Error("預覽面板的狀態標籤不應該為 nil")
	}
	
	// 驗證服務依賴
	if preview.editorService == nil {
		t.Error("編輯器服務不應該為 nil")
	}
	
	// 驗證初始狀態
	if !preview.isVisible {
		t.Error("預覽面板初始狀態應該為可見")
	}
	
	if !preview.autoRefresh {
		t.Error("預覽面板初始狀態應該啟用自動刷新")
	}
	
	if preview.currentContent != "" {
		t.Error("預覽面板初始內容應該為空")
	}
}

// TestMarkdownPreviewUpdatePreview 測試預覽內容更新功能
// 驗證預覽面板是否能正確更新和顯示 Markdown 內容
//
// 測試項目：
// 1. 內容更新是否正確
// 2. 空內容處理是否正確
// 3. 重複內容是否正確處理
// 4. 狀態更新是否正確
func TestMarkdownPreviewUpdatePreview(t *testing.T) {
	// 建立模擬編輯器服務
	mockService := newMockEditorServiceForPreview()
	
	// 建立 Markdown 預覽面板實例
	preview := NewMarkdownPreview(mockService)
	
	// 測試內容更新
	testContent := "# 測試標題\n\n這是測試內容。"
	preview.UpdatePreview(testContent)
	
	// 驗證內容是否正確更新
	if preview.currentContent != testContent {
		t.Errorf("當前內容應該是 '%s'，但得到 '%s'", testContent, preview.currentContent)
	}
	
	// 測試空內容處理
	preview.UpdatePreview("")
	
	if preview.currentContent != "" {
		t.Error("空內容更新後當前內容應該為空")
	}
	
	// 測試重複內容（不應該重複處理）
	preview.UpdatePreview(testContent)
	firstUpdate := preview.currentContent
	preview.UpdatePreview(testContent) // 相同內容
	
	if preview.currentContent != firstUpdate {
		t.Error("重複內容不應該改變當前內容")
	}
}

// TestMarkdownPreviewRefreshPreview 測試手動刷新功能
// 驗證預覽面板的手動刷新功能是否正確
//
// 測試項目：
// 1. 有內容時的刷新是否正確
// 2. 無內容時的刷新處理
// 3. 刷新後狀態更新
func TestMarkdownPreviewRefreshPreview(t *testing.T) {
	// 建立模擬編輯器服務
	mockService := newMockEditorServiceForPreview()
	
	// 建立 Markdown 預覽面板實例
	preview := NewMarkdownPreview(mockService)
	
	// 測試無內容時的刷新
	preview.refreshPreview()
	
	// 驗證狀態標籤是否正確更新
	expectedStatus := "沒有內容可刷新"
	if preview.statusLabel.Text != expectedStatus {
		t.Errorf("無內容刷新狀態應該是 '%s'，但得到 '%s'", expectedStatus, preview.statusLabel.Text)
	}
	
	// 添加內容後測試刷新
	testContent := "# 測試內容"
	preview.UpdatePreview(testContent)
	preview.refreshPreview()
	
	// 驗證刷新後狀態
	if preview.statusLabel.Text != "預覽已手動刷新" {
		t.Error("有內容刷新後狀態應該顯示已刷新")
	}
}

// TestMarkdownPreviewAutoRefresh 測試自動刷新功能
// 驗證自動刷新的啟用、停用和狀態管理
//
// 測試項目：
// 1. 自動刷新啟用和停用
// 2. 狀態查詢功能
// 3. 狀態切換功能
func TestMarkdownPreviewAutoRefresh(t *testing.T) {
	// 建立模擬編輯器服務
	mockService := newMockEditorServiceForPreview()
	
	// 建立 Markdown 預覽面板實例
	preview := NewMarkdownPreview(mockService)
	
	// 驗證初始自動刷新狀態
	if !preview.IsAutoRefreshEnabled() {
		t.Error("初始狀態應該啟用自動刷新")
	}
	
	// 測試停用自動刷新
	preview.SetAutoRefresh(false)
	
	if preview.IsAutoRefreshEnabled() {
		t.Error("設定後自動刷新應該被停用")
	}
	
	// 測試啟用自動刷新
	preview.SetAutoRefresh(true)
	
	if !preview.IsAutoRefreshEnabled() {
		t.Error("設定後自動刷新應該被啟用")
	}
	
	// 測試切換功能
	initialState := preview.IsAutoRefreshEnabled()
	preview.toggleAutoRefresh()
	
	if preview.IsAutoRefreshEnabled() == initialState {
		t.Error("切換後自動刷新狀態應該改變")
	}
}

// TestMarkdownPreviewVisibility 測試可見性控制功能
// 驗證預覽面板的顯示和隱藏功能
//
// 測試項目：
// 1. 可見性設定和查詢
// 2. 可見性切換功能
// 3. 可見性變更回調
func TestMarkdownPreviewVisibility(t *testing.T) {
	// 建立模擬編輯器服務
	mockService := newMockEditorServiceForPreview()
	
	// 建立 Markdown 預覽面板實例
	preview := NewMarkdownPreview(mockService)
	
	// 驗證初始可見性
	if !preview.IsVisible() {
		t.Error("初始狀態應該為可見")
	}
	
	// 測試隱藏預覽面板
	preview.SetVisible(false)
	
	if preview.IsVisible() {
		t.Error("設定後預覽面板應該為隱藏")
	}
	
	// 測試顯示預覽面板
	preview.SetVisible(true)
	
	if !preview.IsVisible() {
		t.Error("設定後預覽面板應該為可見")
	}
	
	// 測試切換功能
	initialVisibility := preview.IsVisible()
	preview.toggleVisibility()
	
	if preview.IsVisible() == initialVisibility {
		t.Error("切換後可見性狀態應該改變")
	}
}

// TestMarkdownPreviewVisibilityCallback 測試可見性變更回調
// 驗證可見性變更時回調函數是否正確觸發
//
// 測試項目：
// 1. 回調函數設定
// 2. 回調函數觸發
// 3. 回調參數正確性
func TestMarkdownPreviewVisibilityCallback(t *testing.T) {
	// 建立模擬編輯器服務
	mockService := newMockEditorServiceForPreview()
	
	// 建立 Markdown 預覽面板實例
	preview := NewMarkdownPreview(mockService)
	
	// 設定回調函數
	var callbackCalled bool
	var callbackVisible bool
	
	preview.SetOnVisibilityChanged(func(visible bool) {
		callbackCalled = true
		callbackVisible = visible
	})
	
	// 測試隱藏觸發回調
	preview.SetVisible(false)
	
	if !callbackCalled {
		t.Error("可見性變更應該觸發回調函數")
	}
	
	if callbackVisible {
		t.Error("回調參數應該為 false（隱藏）")
	}
	
	// 重置回調狀態
	callbackCalled = false
	
	// 測試顯示觸發回調
	preview.SetVisible(true)
	
	if !callbackCalled {
		t.Error("可見性變更應該觸發回調函數")
	}
	
	if !callbackVisible {
		t.Error("回調參數應該為 true（顯示）")
	}
}

// TestMarkdownPreviewContentOperations 測試內容操作功能
// 驗證預覽面板的內容管理功能
//
// 測試項目：
// 1. 內容取得功能
// 2. 內容清空功能
// 3. 內容檢查功能
func TestMarkdownPreviewContentOperations(t *testing.T) {
	// 建立模擬編輯器服務
	mockService := newMockEditorServiceForPreview()
	
	// 建立 Markdown 預覽面板實例
	preview := NewMarkdownPreview(mockService)
	
	// 測試初始狀態
	if preview.HasContent() {
		t.Error("初始狀態不應該有內容")
	}
	
	if preview.GetCurrentContent() != "" {
		t.Error("初始內容應該為空")
	}
	
	// 添加內容
	testContent := "# 測試內容\n\n這是測試。"
	preview.UpdatePreview(testContent)
	
	// 驗證內容操作
	if !preview.HasContent() {
		t.Error("添加內容後應該有內容")
	}
	
	if preview.GetCurrentContent() != testContent {
		t.Errorf("當前內容應該是 '%s'，但得到 '%s'", testContent, preview.GetCurrentContent())
	}
	
	// 測試清空功能
	preview.Clear()
	
	if preview.HasContent() {
		t.Error("清空後不應該有內容")
	}
	
	if preview.GetCurrentContent() != "" {
		t.Error("清空後內容應該為空")
	}
}

// TestMarkdownPreviewWordCount 測試字數統計功能
// 驗證預覽面板的字數和字元統計功能
//
// 測試項目：
// 1. 字數統計準確性
// 2. 字元統計準確性
// 3. 空內容統計
func TestMarkdownPreviewWordCount(t *testing.T) {
	// 建立模擬編輯器服務
	mockService := newMockEditorServiceForPreview()
	
	// 建立 Markdown 預覽面板實例
	preview := NewMarkdownPreview(mockService)
	
	// 測試空內容統計
	if preview.GetWordCount() != 0 {
		t.Error("空內容的字數應該為 0")
	}
	
	if preview.GetCharacterCount() != 0 {
		t.Error("空內容的字元數應該為 0")
	}
	
	// 添加測試內容
	testContent := "這是 測試 內容"
	preview.UpdatePreview(testContent)
	
	// 驗證字數統計
	expectedWordCount := 3 // "這是", "測試", "內容"
	if preview.GetWordCount() != expectedWordCount {
		t.Errorf("字數應該是 %d，但得到 %d", expectedWordCount, preview.GetWordCount())
	}
	
	// 驗證字元統計
	expectedCharCount := len(testContent)
	if preview.GetCharacterCount() != expectedCharCount {
		t.Errorf("字元數應該是 %d，但得到 %d", expectedCharCount, preview.GetCharacterCount())
	}
}

// TestMarkdownPreviewHTMLExport 測試 HTML 匯出功能
// 驗證預覽面板的 HTML 匯出和複製功能
//
// 測試項目：
// 1. 有內容時的匯出處理
// 2. 無內容時的匯出處理
// 3. HTML 文件結構生成
func TestMarkdownPreviewHTMLExport(t *testing.T) {
	// 建立模擬編輯器服務
	mockService := newMockEditorServiceForPreview()
	
	// 建立 Markdown 預覽面板實例
	preview := NewMarkdownPreview(mockService)
	
	// 測試無內容時的匯出
	preview.exportHTML()
	
	expectedStatus := "沒有內容可匯出"
	if preview.statusLabel.Text != expectedStatus {
		t.Errorf("無內容匯出狀態應該是 '%s'，但得到 '%s'", expectedStatus, preview.statusLabel.Text)
	}
	
	// 添加內容後測試匯出
	testContent := "# 測試標題"
	preview.UpdatePreview(testContent)
	preview.exportHTML()
	
	// 驗證匯出狀態（應該顯示準備匯出的訊息）
	if !contains(preview.statusLabel.Text, "HTML 內容已準備匯出") {
		t.Error("有內容匯出後應該顯示準備匯出的狀態")
	}
	
	// 測試 HTML 文件生成
	htmlContent := "<h1>測試標題</h1>"
	fullHTML := preview.createFullHTMLDocument(htmlContent)
	
	// 驗證 HTML 文件結構
	if !contains(fullHTML, "<!DOCTYPE html>") {
		t.Error("完整 HTML 文件應該包含 DOCTYPE 聲明")
	}
	
	if !contains(fullHTML, htmlContent) {
		t.Error("完整 HTML 文件應該包含內容")
	}
	
	if !contains(fullHTML, "<title>Markdown 預覽</title>") {
		t.Error("完整 HTML 文件應該包含標題")
	}
}

// TestMarkdownPreviewCopyHTML 測試 HTML 複製功能
// 驗證預覽面板的 HTML 複製功能
//
// 測試項目：
// 1. 有內容時的複製處理
// 2. 無內容時的複製處理
func TestMarkdownPreviewCopyHTML(t *testing.T) {
	// 建立模擬編輯器服務
	mockService := newMockEditorServiceForPreview()
	
	// 建立 Markdown 預覽面板實例
	preview := NewMarkdownPreview(mockService)
	
	// 測試無內容時的複製
	preview.copyHTML()
	
	expectedStatus := "沒有內容可複製"
	if preview.statusLabel.Text != expectedStatus {
		t.Errorf("無內容複製狀態應該是 '%s'，但得到 '%s'", expectedStatus, preview.statusLabel.Text)
	}
	
	// 添加內容後測試複製
	testContent := "**粗體文字**"
	preview.UpdatePreview(testContent)
	preview.copyHTML()
	
	// 驗證複製狀態（應該顯示準備複製的訊息）
	if !contains(preview.statusLabel.Text, "HTML 內容已準備複製") {
		t.Error("有內容複製後應該顯示準備複製的狀態")
	}
}

// TestMarkdownPreviewGetContainer 測試容器取得功能
// 驗證預覽面板是否能正確回傳主要容器
//
// 測試項目：
// 1. GetContainer 是否回傳正確的容器實例
// 2. 回傳的容器是否與內部容器相同
func TestMarkdownPreviewGetContainer(t *testing.T) {
	// 建立模擬編輯器服務
	mockService := newMockEditorServiceForPreview()
	
	// 建立 Markdown 預覽面板實例
	preview := NewMarkdownPreview(mockService)
	
	// 取得容器
	container := preview.GetContainer()
	
	// 驗證容器不為 nil
	if container == nil {
		t.Error("GetContainer 應該回傳有效的容器實例")
	}
	
	// 驗證回傳的容器與內部容器相同
	if container != preview.container {
		t.Error("GetContainer 應該回傳與內部容器相同的實例")
	}
}

// TestMarkdownPreviewScrollSync 測試滾動同步功能
// 驗證預覽面板的滾動同步功能（目前為佔位實作）
//
// 測試項目：
// 1. 滾動位置設定
// 2. 滾動位置取得
// 3. 同步滾動切換
func TestMarkdownPreviewScrollSync(t *testing.T) {
	// 建立模擬編輯器服務
	mockService := newMockEditorServiceForPreview()
	
	// 建立 Markdown 預覽面板實例
	preview := NewMarkdownPreview(mockService)
	
	// 測試滾動位置同步
	testPosition := 0.5
	preview.SyncScrollPosition(testPosition)
	
	// 驗證狀態更新（目前為佔位實作）
	expectedStatus := "滾動同步: 50.0%"
	if preview.statusLabel.Text != expectedStatus {
		t.Errorf("滾動同步狀態應該是 '%s'，但得到 '%s'", expectedStatus, preview.statusLabel.Text)
	}
	
	// 測試滾動位置取得（目前回傳 0.0）
	position := preview.GetScrollPosition()
	if position != 0.0 {
		t.Errorf("滾動位置應該是 0.0，但得到 %f", position)
	}
	
	// 測試同步滾動切換
	preview.toggleSyncScroll()
	
	// 驗證狀態更新
	if !contains(preview.statusLabel.Text, "同步滾動功能將在未來版本中實作") {
		t.Error("同步滾動切換應該顯示未來實作的訊息")
	}
}

// TestMarkdownPreviewThemeAndSettings 測試主題和設定功能
// 驗證預覽面板的主題設定功能（目前為佔位實作）
//
// 測試項目：
// 1. 主題設定功能
// 2. 全螢幕切換功能
func TestMarkdownPreviewThemeAndSettings(t *testing.T) {
	// 建立模擬編輯器服務
	mockService := newMockEditorServiceForPreview()
	
	// 建立 Markdown 預覽面板實例
	preview := NewMarkdownPreview(mockService)
	
	// 測試主題設定
	testTheme := "dark"
	preview.SetTheme(testTheme)
	
	// 驗證狀態更新
	expectedStatus := "主題設定: dark (將在後續版本中實作)"
	if preview.statusLabel.Text != expectedStatus {
		t.Errorf("主題設定狀態應該是 '%s'，但得到 '%s'", expectedStatus, preview.statusLabel.Text)
	}
	
	// 測試全螢幕切換
	preview.toggleFullscreen()
	
	// 驗證狀態更新
	if !contains(preview.statusLabel.Text, "全螢幕預覽功能將在後續版本中實作") {
		t.Error("全螢幕切換應該顯示未來實作的訊息")
	}
}

// contains 檢查字串是否包含子字串的輔助函數
// 參數：s（主字串）、substr（子字串）
// 回傳：是否包含的布林值
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || 
		(len(s) > len(substr) && 
		 (s[:len(substr)] == substr || 
		  s[len(s)-len(substr):] == substr || 
		  containsInMiddle(s, substr))))
}

// containsInMiddle 檢查字串中間是否包含子字串的輔助函數
func containsInMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}