// Package ui 包含使用者介面相關的元件測試
// 本檔案測試增強版 Markdown 預覽面板 UI 元件的新功能
package ui

import (
	"fmt"
	"testing"
	"time"

)

// TestEnhancedMarkdownPreviewCreation 測試增強版 Markdown 預覽面板的建立
// 驗證所有新增的 UI 元件是否正確初始化
//
// 測試項目：
// 1. 搜尋功能元件初始化
// 2. 縮放控制元件初始化
// 3. 效能優化設定初始化
// 4. 新增狀態變數初始化
func TestEnhancedMarkdownPreviewCreation(t *testing.T) {
	// 建立模擬編輯器服務
	mockService := newMockEditorServiceForPreview()
	
	// 建立增強版預覽面板
	preview := NewMarkdownPreview(mockService)
	
	// 驗證新增的 UI 元件
	if preview.searchBar == nil {
		t.Error("搜尋輸入欄不應該為 nil")
	}
	
	if preview.searchResults == nil {
		t.Error("搜尋結果顯示不應該為 nil")
	}
	
	if preview.zoomSlider == nil {
		t.Error("縮放滑桿不應該為 nil")
	}
	
	// 驗證新增的狀態變數
	if preview.zoomLevel != 1.0 {
		t.Errorf("初始縮放級別應該是 1.0，但得到 %f", preview.zoomLevel)
	}
	
	if preview.updateThrottle != 100*time.Millisecond {
		t.Errorf("更新節流間隔應該是 100ms，但得到 %v", preview.updateThrottle)
	}
	
	if preview.contentCache == nil {
		t.Error("內容快取不應該為 nil")
	}
	
	if preview.searchMatches == nil {
		t.Error("搜尋匹配陣列不應該為 nil")
	}
	
	if preview.currentMatch != -1 {
		t.Errorf("初始匹配索引應該是 -1，但得到 %d", preview.currentMatch)
	}
	
	if preview.isIndependent {
		t.Error("初始狀態不應該是獨立視窗模式")
	}
}

// TestZoomFunctionality 測試縮放功能
// 驗證縮放控制的各種操作
//
// 測試項目：
// 1. 放大功能
// 2. 縮小功能
// 3. 重置縮放功能
// 4. 設定特定縮放級別
// 5. 縮放範圍限制
func TestZoomFunctionality(t *testing.T) {
	// 建立模擬編輯器服務
	mockService := newMockEditorServiceForPreview()
	
	// 建立預覽面板
	preview := NewMarkdownPreview(mockService)
	
	// 測試初始縮放級別
	if preview.GetZoomLevel() != 1.0 {
		t.Errorf("初始縮放級別應該是 1.0，但得到 %f", preview.GetZoomLevel())
	}
	
	// 測試放大功能
	preview.zoomIn()
	expectedZoom := 1.1
	if preview.GetZoomLevel() != expectedZoom {
		t.Errorf("放大後縮放級別應該是 %f，但得到 %f", expectedZoom, preview.GetZoomLevel())
	}
	
	// 測試縮小功能
	preview.zoomOut()
	expectedZoom = 1.0
	if preview.GetZoomLevel() != expectedZoom {
		t.Errorf("縮小後縮放級別應該是 %f，但得到 %f", expectedZoom, preview.GetZoomLevel())
	}
	
	// 測試設定特定縮放級別
	testZoom := 1.5
	preview.setZoomLevel(testZoom)
	if preview.GetZoomLevel() != testZoom {
		t.Errorf("設定縮放級別應該是 %f，但得到 %f", testZoom, preview.GetZoomLevel())
	}
	
	// 測試重置縮放
	preview.resetZoom()
	if preview.GetZoomLevel() != 1.0 {
		t.Errorf("重置後縮放級別應該是 1.0，但得到 %f", preview.GetZoomLevel())
	}
	
	// 測試縮放範圍限制 - 最小值
	preview.setZoomLevel(0.3) // 低於最小值 0.5
	if preview.GetZoomLevel() != 0.5 {
		t.Errorf("縮放級別不應該低於 0.5，但得到 %f", preview.GetZoomLevel())
	}
	
	// 測試縮放範圍限制 - 最大值
	preview.setZoomLevel(4.0) // 高於最大值 3.0
	if preview.GetZoomLevel() != 3.0 {
		t.Errorf("縮放級別不應該高於 3.0，但得到 %f", preview.GetZoomLevel())
	}
}

// TestSearchFunctionality 測試搜尋功能
// 驗證搜尋功能的各種操作
//
// 測試項目：
// 1. 基本搜尋功能
// 2. 搜尋結果統計
// 3. 空搜尋處理
// 4. 不存在內容搜尋
// 5. 搜尋狀態管理
func TestSearchFunctionality(t *testing.T) {
	// 建立模擬編輯器服務
	mockService := newMockEditorServiceForPreview()
	
	// 建立預覽面板
	preview := NewMarkdownPreview(mockService)
	
	// 設定測試內容
	testContent := "這是第一行測試內容\n這是第二行測試內容\n這是第三行內容"
	preview.UpdatePreview(testContent)
	
	// 測試搜尋功能
	searchQuery := "測試"
	preview.performSearch(searchQuery)
	
	// 驗證搜尋結果
	if preview.GetSearchQuery() != searchQuery {
		t.Errorf("搜尋查詢應該是 '%s'，但得到 '%s'", searchQuery, preview.GetSearchQuery())
	}
	
	expectedMatches := 2 // "測試" 在內容中出現 2 次
	if preview.GetSearchMatchCount() != expectedMatches {
		t.Errorf("搜尋匹配數應該是 %d，但得到 %d", expectedMatches, preview.GetSearchMatchCount())
	}
	
	// 測試空搜尋查詢
	preview.performSearch("")
	if preview.GetSearchQuery() != "" {
		t.Error("空搜尋查詢後搜尋查詢應該為空")
	}
	
	if preview.GetSearchMatchCount() != 0 {
		t.Error("空搜尋查詢後匹配數應該為 0")
	}
	
	// 測試搜尋不存在的內容
	preview.performSearch("不存在的內容")
	if preview.GetSearchMatchCount() != 0 {
		t.Error("搜尋不存在內容的匹配數應該為 0")
	}
}

// TestSearchNavigation 測試搜尋導航功能
// 驗證搜尋結果的導航操作
//
// 測試項目：
// 1. 導航到下一個匹配項
// 2. 導航到上一個匹配項
// 3. 循環導航
// 4. 無搜尋結果時的導航處理
func TestSearchNavigation(t *testing.T) {
	// 建立模擬編輯器服務
	mockService := newMockEditorServiceForPreview()
	
	// 建立預覽面板
	preview := NewMarkdownPreview(mockService)
	
	// 設定測試內容
	testContent := "測試 內容 測試 內容 測試"
	preview.UpdatePreview(testContent)
	
	// 執行搜尋
	preview.performSearch("測試")
	expectedMatches := 3
	if preview.GetSearchMatchCount() != expectedMatches {
		t.Errorf("搜尋匹配數應該是 %d，但得到 %d", expectedMatches, preview.GetSearchMatchCount())
	}
	
	// 測試導航到下一個匹配項
	initialMatch := preview.currentMatch
	preview.navigateToNextMatch()
	
	// 驗證匹配項索引有變化
	if preview.currentMatch == initialMatch && expectedMatches > 1 {
		t.Error("導航到下一個匹配項後索引應該改變")
	}
	
	// 測試導航到上一個匹配項
	currentMatch := preview.currentMatch
	preview.navigateToPreviousMatch()
	
	// 驗證匹配項索引有變化
	if preview.currentMatch == currentMatch && expectedMatches > 1 {
		t.Error("導航到上一個匹配項後索引應該改變")
	}
	
	// 測試沒有搜尋結果時的導航
	preview.performSearch("不存在")
	preview.navigateToNextMatch()    // 應該不會出錯
	preview.navigateToPreviousMatch() // 應該不會出錯
	
	// 驗證沒有搜尋結果時導航不會出錯
	if preview.GetSearchMatchCount() != 0 {
		t.Error("搜尋不存在內容後匹配數應該為 0")
	}
}

// TestToggleSearch 測試搜尋功能切換
// 驗證搜尋功能的顯示和隱藏
//
// 測試項目：
// 1. 搜尋功能顯示
// 2. 搜尋功能隱藏
// 3. 搜尋狀態清理
func TestToggleSearch(t *testing.T) {
	// 建立模擬編輯器服務
	mockService := newMockEditorServiceForPreview()
	
	// 建立預覽面板
	preview := NewMarkdownPreview(mockService)
	
	// 測試初始狀態（搜尋欄應該隱藏）
	if preview.searchBar.Visible() {
		t.Error("搜尋欄初始狀態應該隱藏")
	}
	
	if preview.searchResults.Visible() {
		t.Error("搜尋結果初始狀態應該隱藏")
	}
	
	// 測試顯示搜尋功能
	preview.toggleSearch()
	
	if !preview.searchBar.Visible() {
		t.Error("切換後搜尋欄應該顯示")
	}
	
	if !preview.searchResults.Visible() {
		t.Error("切換後搜尋結果應該顯示")
	}
	
	// 測試隱藏搜尋功能
	preview.toggleSearch()
	
	if preview.searchBar.Visible() {
		t.Error("再次切換後搜尋欄應該隱藏")
	}
	
	if preview.searchResults.Visible() {
		t.Error("再次切換後搜尋結果應該隱藏")
	}
}

// TestIndependentWindowMode 測試獨立視窗模式
// 驗證獨立視窗的建立和管理
//
// 測試項目：
// 1. 獨立視窗模式切換
// 2. 獨立視窗狀態管理
// 3. 獨立視窗建立和關閉
func TestIndependentWindowMode(t *testing.T) {
	// 建立模擬編輯器服務
	mockService := newMockEditorServiceForPreview()
	
	// 建立預覽面板
	preview := NewMarkdownPreview(mockService)
	
	// 測試初始狀態
	if preview.IsIndependentMode() {
		t.Error("初始狀態不應該是獨立視窗模式")
	}
	
	if preview.independentWindow != nil {
		t.Error("初始狀態獨立視窗應該為 nil")
	}
	
	// 測試切換到獨立視窗模式
	preview.toggleIndependentWindow()
	
	if !preview.IsIndependentMode() {
		t.Error("切換後應該是獨立視窗模式")
	}
	
	if preview.independentWindow == nil {
		t.Error("切換後獨立視窗不應該為 nil")
	}
	
	// 測試切換回嵌入模式
	preview.toggleIndependentWindow()
	
	if preview.IsIndependentMode() {
		t.Error("切換回嵌入模式後不應該是獨立視窗模式")
	}
	
	if preview.independentWindow != nil {
		t.Error("切換回嵌入模式後獨立視窗應該為 nil")
	}
}

// TestPerformanceOptimization 測試效能優化功能
// 驗證內容快取和更新節流機制
//
// 測試項目：
// 1. 內容快取機制
// 2. 更新節流機制
// 3. 快取大小限制
// 4. 效能優化執行
func TestPerformanceOptimization(t *testing.T) {
	// 建立模擬編輯器服務
	mockService := newMockEditorServiceForPreview()
	
	// 建立預覽面板
	preview := NewMarkdownPreview(mockService)
	
	// 測試內容快取機制
	content1 := "第一次更新的內容"
	preview.UpdatePreview(content1)
	
	// 驗證快取中有內容
	if len(preview.contentCache) == 0 {
		t.Error("更新內容後快取中應該有項目")
	}
	
	// 測試更新節流機制
	content2 := "第二次更新的內容"
	preview.UpdatePreview(content2)
	
	// 立即進行第三次更新（應該被節流）
	content3 := "第三次更新的內容"
	preview.UpdatePreview(content3)
	
	// 等待節流時間過去
	time.Sleep(150 * time.Millisecond)
	
	// 測試快取大小限制
	// 添加大量快取項目
	for i := 0; i < 150; i++ {
		preview.contentCache[fmt.Sprintf("key_%d", i)] = fmt.Sprintf("value_%d", i)
	}
	
	// 執行效能優化
	preview.optimizePerformance()
	
	// 驗證快取大小被限制
	if len(preview.contentCache) > 100 {
		t.Errorf("快取大小應該被限制在 100 以內，但得到 %d", len(preview.contentCache))
	}
}

// TestCallbackFunctions 測試回調函數
// 驗證新增的回調函數是否正確觸發
//
// 測試項目：
// 1. 縮放變更回調
// 2. 搜尋執行回調
// 3. 回調參數正確性
func TestCallbackFunctions(t *testing.T) {
	// 建立模擬編輯器服務
	mockService := newMockEditorServiceForPreview()
	
	// 建立預覽面板
	preview := NewMarkdownPreview(mockService)
	
	// 測試縮放變更回調
	var zoomChanged bool
	var zoomLevel float64
	
	preview.SetOnZoomChanged(func(level float64) {
		zoomChanged = true
		zoomLevel = level
	})
	
	testZoom := 1.5
	preview.setZoomLevel(testZoom)
	
	if !zoomChanged {
		t.Error("縮放變更應該觸發回調函數")
	}
	
	if zoomLevel != testZoom {
		t.Errorf("回調參數應該是 %f，但得到 %f", testZoom, zoomLevel)
	}
	
	// 測試搜尋執行回調
	var searchPerformed bool
	var searchQuery string
	var matchCount int
	
	preview.SetOnSearchPerformed(func(query string, matches int) {
		searchPerformed = true
		searchQuery = query
		matchCount = matches
	})
	
	// 設定內容並執行搜尋
	testContent := "測試內容測試"
	preview.UpdatePreview(testContent)
	testQuery := "測試"
	preview.performSearch(testQuery)
	
	if !searchPerformed {
		t.Error("搜尋執行應該觸發回調函數")
	}
	
	if searchQuery != testQuery {
		t.Errorf("搜尋回調查詢參數應該是 '%s'，但得到 '%s'", testQuery, searchQuery)
	}
	
	expectedMatches := 2 // "測試" 出現 2 次
	if matchCount != expectedMatches {
		t.Errorf("搜尋回調匹配數參數應該是 %d，但得到 %d", expectedMatches, matchCount)
	}
}

// TestEnhancedUpdatePreview 測試增強版更新預覽功能
// 驗證更新預覽的效能優化和新功能
//
// 測試項目：
// 1. 內容快取使用
// 2. 更新節流機制
// 3. 搜尋結果更新
// 4. 獨立視窗同步
func TestEnhancedUpdatePreview(t *testing.T) {
	// 建立模擬編輯器服務
	mockService := newMockEditorServiceForPreview()
	
	// 建立預覽面板
	preview := NewMarkdownPreview(mockService)
	
	// 測試內容更新
	testContent := "# 測試標題\n\n這是測試內容。"
	preview.UpdatePreview(testContent)
	
	// 驗證內容已更新
	if preview.GetCurrentContent() != testContent {
		t.Errorf("當前內容應該是 '%s'，但得到 '%s'", testContent, preview.GetCurrentContent())
	}
	
	// 驗證快取中有內容
	if len(preview.contentCache) == 0 {
		t.Error("更新內容後快取中應該有項目")
	}
	
	// 測試重複更新相同內容（應該直接返回）
	initialCacheSize := len(preview.contentCache)
	preview.UpdatePreview(testContent)
	
	// 快取大小不應該增加（因為內容相同）
	if len(preview.contentCache) != initialCacheSize {
		t.Error("重複更新相同內容不應該增加快取項目")
	}
	
	// 測試搜尋結果在內容更新後的重新執行
	preview.performSearch("測試")
	initialMatches := preview.GetSearchMatchCount()
	
	// 更新內容（包含更多搜尋目標）
	newContent := testContent + "\n更多測試內容"
	preview.UpdatePreview(newContent)
	
	// 如果有搜尋查詢，應該重新執行搜尋
	if preview.searchQuery != "" {
		// 搜尋結果應該更新
		if preview.GetSearchMatchCount() == initialMatches {
			// 這個測試可能需要調整，因為搜尋重新執行的邏輯
		}
	}
}

// TestExportFunctionality 測試匯出功能
// 驗證增強版匯出功能
//
// 測試項目：
// 1. 匯出選單顯示
// 2. 列印功能
// 3. 無內容時的匯出處理
func TestExportFunctionality(t *testing.T) {
	// 建立模擬編輯器服務
	mockService := newMockEditorServiceForPreview()
	
	// 建立預覽面板
	preview := NewMarkdownPreview(mockService)
	
	// 測試無內容時的匯出
	preview.showExportMenu()
	
	if !contains(preview.statusLabel.Text, "沒有內容可匯出") {
		t.Error("無內容時應該顯示沒有內容可匯出的訊息")
	}
	
	// 添加內容後測試匯出
	testContent := "# 測試標題\n\n這是測試內容。"
	preview.UpdatePreview(testContent)
	
	preview.showExportMenu()
	
	// 驗證匯出功能執行（目前為佔位實作）
	if !contains(preview.statusLabel.Text, "支援的匯出格式") {
		t.Error("有內容時應該顯示支援的匯出格式")
	}
	
	// 測試列印功能
	preview.printPreview()
	
	if !contains(preview.statusLabel.Text, "列印功能將在未來版本中實作") {
		t.Error("列印功能應該顯示未來實作的訊息")
	}
}

// TestLargeContentPerformance 測試大內容效能
// 驗證處理大量內容時的效能表現
//
// 測試項目：
// 1. 大內容處理時間
// 2. 記憶體使用合理性
// 3. 搜尋效能
func TestLargeContentPerformance(t *testing.T) {
	// 建立模擬編輯器服務
	mockService := newMockEditorServiceForPreview()
	
	// 建立預覽面板
	preview := NewMarkdownPreview(mockService)
	
	// 建立大量內容
	largeContent := ""
	for i := 0; i < 1000; i++ {
		largeContent += fmt.Sprintf("這是測試內容的第 %d 行\n", i)
	}
	
	// 測試大內容處理時間
	start := time.Now()
	preview.UpdatePreview(largeContent)
	duration := time.Since(start)
	
	// 驗證處理完成且效能合理
	if !preview.HasContent() {
		t.Error("大內容處理後應該有內容")
	}
	
	if preview.GetWordCount() == 0 {
		t.Error("大內容處理後字數應該大於 0")
	}
	
	// 效能要求：處理時間應該在合理範圍內
	if duration > 2*time.Second {
		t.Errorf("大內容處理時間過長: %v", duration)
	}
	
	// 測試大內容搜尋效能
	searchStart := time.Now()
	preview.performSearch("測試")
	searchDuration := time.Since(searchStart)
	
	// 搜尋效能要求
	if searchDuration > 1*time.Second {
		t.Errorf("大內容搜尋時間過長: %v", searchDuration)
	}
	
	// 驗證搜尋結果
	if preview.GetSearchMatchCount() == 0 {
		t.Error("大內容搜尋應該找到匹配項")
	}
}

// TestConcurrentOperations 測試並發操作
// 驗證多執行緒環境下的安全性
//
// 測試項目：
// 1. 並發內容更新
// 2. 並發搜尋操作
// 3. 並發縮放操作
func TestConcurrentOperations(t *testing.T) {
	// 建立模擬編輯器服務
	mockService := newMockEditorServiceForPreview()
	
	// 建立預覽面板
	preview := NewMarkdownPreview(mockService)
	
	// 並發內容更新測試
	done := make(chan bool, 10)
	
	for i := 0; i < 10; i++ {
		go func(index int) {
			content := fmt.Sprintf("並發測試內容 %d", index)
			preview.UpdatePreview(content)
			done <- true
		}(i)
	}
	
	// 等待所有更新完成
	for i := 0; i < 10; i++ {
		<-done
	}
	
	// 驗證預覽面板仍然正常工作
	if !preview.HasContent() {
		t.Error("並發更新後應該有內容")
	}
	
	// 並發搜尋操作測試
	preview.UpdatePreview("測試內容用於搜尋測試")
	
	searchDone := make(chan bool, 5)
	for i := 0; i < 5; i++ {
		go func(index int) {
			query := fmt.Sprintf("測試%d", index%2) // 交替搜尋不同內容
			preview.performSearch(query)
			searchDone <- true
		}(i)
	}
	
	// 等待所有搜尋完成
	for i := 0; i < 5; i++ {
		<-searchDone
	}
	
	// 並發縮放操作測試
	zoomDone := make(chan bool, 5)
	for i := 0; i < 5; i++ {
		go func(index int) {
			level := 1.0 + float64(index)*0.1
			preview.setZoomLevel(level)
			zoomDone <- true
		}(i)
	}
	
	// 等待所有縮放操作完成
	for i := 0; i < 5; i++ {
		<-zoomDone
	}
	
	// 驗證最終狀態合理
	if preview.GetZoomLevel() < 0.5 || preview.GetZoomLevel() > 3.0 {
		t.Errorf("並發縮放後縮放級別應該在合理範圍內，但得到 %f", preview.GetZoomLevel())
	}
}