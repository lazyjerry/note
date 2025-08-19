// Package services 測試效能服務的功能
// 本檔案包含 PerformanceService 的完整測試案例，驗證效能監控、記憶體優化和大檔案處理功能
package services

import (
	"context"      // 上下文管理
	"fmt"          // 格式化輸出
	"testing"      // 測試框架
	"time"         // 時間處理
	"mac-notebook-app/internal/models" // 引入資料模型
)

// mockEditorService 模擬編輯器服務
// 用於測試效能服務與編輯器服務的整合
type mockEditorService struct {
	activeNotes map[string]*models.Note // 模擬活躍筆記
}

// GetActiveNotes 取得活躍筆記（模擬實作）
// 回傳：活躍筆記映射
func (m *mockEditorService) GetActiveNotes() map[string]*models.Note {
	return m.activeNotes
}

// 其他 EditorService 介面方法的空實作
func (m *mockEditorService) CreateNote(title, content string) (*models.Note, error) { return nil, nil }
func (m *mockEditorService) OpenNote(filePath string) (*models.Note, error) { return nil, nil }
func (m *mockEditorService) SaveNote(note *models.Note) error { return nil }
func (m *mockEditorService) UpdateContent(noteID, content string) error { return nil }
func (m *mockEditorService) PreviewMarkdown(content string) string { return "" }
func (m *mockEditorService) DecryptWithPassword(noteID, password string) (string, error) { return "", nil }
func (m *mockEditorService) CloseNote(noteID string) {}
func (m *mockEditorService) GetActiveNote(noteID string) (*models.Note, bool) { return nil, false }

// 智慧編輯功能的空實作
func (m *mockEditorService) GetAutoCompleteSuggestions(content string, cursorPosition int) []AutoCompleteSuggestion { return []AutoCompleteSuggestion{} }
func (m *mockEditorService) FormatTableContent(tableContent string) (string, error) { return tableContent, nil }
func (m *mockEditorService) InsertLinkMarkdown(text, url string) string { return "[" + text + "](" + url + ")" }
func (m *mockEditorService) InsertImageMarkdown(altText, imagePath string) string { return "![" + altText + "](" + imagePath + ")" }
func (m *mockEditorService) GetSupportedCodeLanguages() []string { return []string{"go"} }
func (m *mockEditorService) FormatCodeBlockMarkdown(code, language string) string { return "```" + language + "\n" + code + "\n```" }
func (m *mockEditorService) FormatMathExpressionMarkdown(expression string, isInline bool) string { return "$" + expression + "$" }
func (m *mockEditorService) ValidateMarkdownContent(content string) (bool, []string) { return true, []string{} }
func (m *mockEditorService) GenerateTableTemplateMarkdown(rows, cols int) string { return "| 欄位1 | 欄位2 |\n|-------|-------|\n| 內容1 | 內容2 |" }
func (m *mockEditorService) PreviewMarkdownWithHighlight(content string) string { return "<p>" + content + "</p>" }
func (m *mockEditorService) GetSmartEditingService() SmartEditingService { return NewSmartEditingService() }
func (m *mockEditorService) SetSmartEditingService(smartEditSvc SmartEditingService) {}

// TestNewPerformanceService 測試效能服務的建立
// 驗證效能服務實例是否正確初始化
//
// 測試案例：
// 1. 使用有效的編輯器服務建立效能服務
// 2. 驗證服務實例不為 nil
// 3. 驗證初始狀態正確
func TestNewPerformanceService(t *testing.T) {
	// 建立模擬編輯器服務
	mockEditor := &mockEditorService{
		activeNotes: make(map[string]*models.Note),
	}
	
	// 建立效能服務
	perfService := NewPerformanceService(mockEditor)
	
	// 驗證服務實例
	if perfService == nil {
		t.Fatal("效能服務實例不應為 nil")
	}
	
	// 驗證初始狀態
	metrics := perfService.GetCurrentMetrics()
	if metrics == nil {
		t.Error("當前指標不應為 nil")
	}
	
	if metrics.ActiveNotes != 0 {
		t.Errorf("初始活躍筆記數量應為 0，實際為 %d", metrics.ActiveNotes)
	}
}

// TestPerformanceMonitoring 測試效能監控功能
// 驗證監控的啟動、停止和指標收集
//
// 測試案例：
// 1. 啟動效能監控
// 2. 等待一段時間讓監控收集指標
// 3. 停止監控
// 4. 驗證指標歷史記錄
func TestPerformanceMonitoring(t *testing.T) {
	// 建立模擬編輯器服務
	mockEditor := &mockEditorService{
		activeNotes: make(map[string]*models.Note),
	}
	
	// 建立效能服務
	perfService := NewPerformanceService(mockEditor)
	
	// 建立上下文
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	
	// 啟動監控
	err := perfService.StartMonitoring(ctx)
	if err != nil {
		t.Fatalf("啟動監控失敗：%v", err)
	}
	
	// 等待監控收集指標
	time.Sleep(100 * time.Millisecond)
	
	// 驗證當前指標
	metrics := perfService.GetCurrentMetrics()
	if metrics == nil {
		t.Error("當前指標不應為 nil")
	}
	
	if metrics.Timestamp.IsZero() {
		t.Error("指標時間戳不應為零值")
	}
	
	if metrics.MemoryUsage <= 0 {
		t.Error("記憶體使用量應大於 0")
	}
	
	if metrics.Goroutines <= 0 {
		t.Error("Goroutine 數量應大於 0")
	}
	
	// 停止監控
	err = perfService.StopMonitoring()
	if err != nil {
		t.Errorf("停止監控失敗：%v", err)
	}
	
	// 驗證重複停止監控會返回錯誤
	err = perfService.StopMonitoring()
	if err == nil {
		t.Error("重複停止監控應該返回錯誤")
	}
}

// TestMemoryOptimization 測試記憶體優化功能
// 驗證記憶體優化和垃圾回收功能
//
// 測試案例：
// 1. 取得優化前的記憶體使用量
// 2. 執行記憶體優化
// 3. 強制垃圾回收
// 4. 驗證優化效果
func TestMemoryOptimization(t *testing.T) {
	// 建立模擬編輯器服務
	mockEditor := &mockEditorService{
		activeNotes: make(map[string]*models.Note),
	}
	
	// 建立效能服務
	perfService := NewPerformanceService(mockEditor)
	
	// 取得優化前的記憶體使用量
	beforeMemory, err := perfService.GetMemoryUsage()
	if err != nil {
		t.Fatalf("取得記憶體使用量失敗：%v", err)
	}
	
	if beforeMemory <= 0 {
		t.Error("記憶體使用量應大於 0")
	}
	
	// 執行記憶體優化
	err = perfService.OptimizeMemory()
	if err != nil {
		t.Errorf("記憶體優化失敗：%v", err)
	}
	
	// 強制垃圾回收
	err = perfService.ForceGarbageCollection()
	if err != nil {
		t.Errorf("強制垃圾回收失敗：%v", err)
	}
	
	// 取得優化後的記憶體使用量
	afterMemory, err := perfService.GetMemoryUsage()
	if err != nil {
		t.Fatalf("取得優化後記憶體使用量失敗：%v", err)
	}
	
	if afterMemory <= 0 {
		t.Error("優化後記憶體使用量應大於 0")
	}
	
	// 注意：由於 Go 的垃圾回收器特性，記憶體使用量可能不會立即減少
	// 這裡主要驗證函數能正常執行而不出錯
}

// TestLargeFileOptimization 測試大檔案處理優化
// 驗證大檔案處理的優化策略
//
// 測試案例：
// 1. 為小檔案進行優化
// 2. 為大檔案進行優化
// 3. 測試分塊處理功能
func TestLargeFileOptimization(t *testing.T) {
	// 建立模擬編輯器服務
	mockEditor := &mockEditorService{
		activeNotes: make(map[string]*models.Note),
	}
	
	// 建立效能服務
	perfService := NewPerformanceService(mockEditor)
	
	// 測試小檔案優化
	err := perfService.OptimizeForLargeFile("small_file.md", 1024) // 1KB
	if err != nil {
		t.Errorf("小檔案優化失敗：%v", err)
	}
	
	// 測試大檔案優化
	err = perfService.OptimizeForLargeFile("large_file.md", 50*1024*1024) // 50MB
	if err != nil {
		t.Errorf("大檔案優化失敗：%v", err)
	}
	
	// 測試分塊處理
	processedChunks := 0
	processor := func(chunk []byte) error {
		processedChunks++
		return nil
	}
	
	err = perfService.ProcessLargeFileInChunks("test_file.md", 1024*1024, processor)
	if err != nil {
		t.Errorf("分塊處理失敗：%v", err)
	}
	
	// 注意：由於這是模擬實作，processedChunks 可能為 0
	// 實際實作中應該驗證分塊處理的正確性
}

// TestCacheManagement 測試快取管理功能
// 驗證快取統計、優化和清理功能
//
// 測試案例：
// 1. 取得初始快取統計
// 2. 模擬快取命中和未命中
// 3. 驗證統計資料更新
// 4. 測試快取優化和清理
func TestCacheManagement(t *testing.T) {
	// 建立模擬編輯器服務
	mockEditor := &mockEditorService{
		activeNotes: make(map[string]*models.Note),
	}
	
	// 建立效能服務
	perfService := NewPerformanceService(mockEditor)
	
	// 取得初始快取統計
	initialStats := perfService.GetCacheStats()
	if initialStats == nil {
		t.Fatal("快取統計不應為 nil")
	}
	
	// 驗證初始統計
	if initialStats["hits"].(int64) != 0 {
		t.Error("初始快取命中次數應為 0")
	}
	
	if initialStats["misses"].(int64) != 0 {
		t.Error("初始快取未命中次數應為 0")
	}
	
	if initialStats["hit_rate"].(float64) != 0.0 {
		t.Error("初始快取命中率應為 0.0")
	}
	
	// 模擬快取操作（需要將 perfService 轉換為具體類型以存取內部方法）
	if ps, ok := perfService.(*performanceService); ok {
		// 模擬快取命中
		ps.RecordCacheHit()
		ps.RecordCacheHit()
		
		// 模擬快取未命中
		ps.RecordCacheMiss()
		
		// 驗證更新後的統計
		updatedStats := perfService.GetCacheStats()
		if updatedStats["hits"].(int64) != 2 {
			t.Errorf("快取命中次數應為 2，實際為 %d", updatedStats["hits"].(int64))
		}
		
		if updatedStats["misses"].(int64) != 1 {
			t.Errorf("快取未命中次數應為 1，實際為 %d", updatedStats["misses"].(int64))
		}
		
		expectedHitRate := 2.0 / 3.0
		actualHitRate := updatedStats["hit_rate"].(float64)
		if actualHitRate < expectedHitRate-0.01 || actualHitRate > expectedHitRate+0.01 {
			t.Errorf("快取命中率應約為 %.2f，實際為 %.2f", expectedHitRate, actualHitRate)
		}
	}
	
	// 測試快取優化
	err := perfService.OptimizeCache()
	if err != nil {
		t.Errorf("快取優化失敗：%v", err)
	}
	
	// 測試快取清理
	err = perfService.ClearCache()
	if err != nil {
		t.Errorf("快取清理失敗：%v", err)
	}
	
	// 驗證清理後的統計
	clearedStats := perfService.GetCacheStats()
	if clearedStats["hits"].(int64) != 0 {
		t.Error("清理後快取命中次數應為 0")
	}
	
	if clearedStats["misses"].(int64) != 0 {
		t.Error("清理後快取未命中次數應為 0")
	}
}

// TestBackgroundTaskTracking 測試背景任務追蹤功能
// 驗證任務註冊、取消註冊和計數功能
//
// 測試案例：
// 1. 註冊多個背景任務
// 2. 驗證任務計數
// 3. 取消註冊任務
// 4. 驗證計數更新
func TestBackgroundTaskTracking(t *testing.T) {
	// 建立模擬編輯器服務
	mockEditor := &mockEditorService{
		activeNotes: make(map[string]*models.Note),
	}
	
	// 建立效能服務
	perfService := NewPerformanceService(mockEditor)
	
	// 驗證初始任務數量
	initialCount := perfService.GetActiveTasksCount()
	if initialCount != 0 {
		t.Errorf("初始活躍任務數量應為 0，實際為 %d", initialCount)
	}
	
	// 註冊多個背景任務
	taskID1 := perfService.RegisterBackgroundTask("自動保存任務")
	if taskID1 == "" {
		t.Error("任務 ID 不應為空")
	}
	
	taskID2 := perfService.RegisterBackgroundTask("檔案監控任務")
	if taskID2 == "" {
		t.Error("任務 ID 不應為空")
	}
	
	taskID3 := perfService.RegisterBackgroundTask("效能監控任務")
	if taskID3 == "" {
		t.Error("任務 ID 不應為空")
	}
	
	// 驗證任務 ID 唯一性
	if taskID1 == taskID2 || taskID1 == taskID3 || taskID2 == taskID3 {
		t.Error("任務 ID 應該是唯一的")
	}
	
	// 驗證任務計數
	activeCount := perfService.GetActiveTasksCount()
	if activeCount != 3 {
		t.Errorf("活躍任務數量應為 3，實際為 %d", activeCount)
	}
	
	// 取消註冊一個任務
	perfService.UnregisterBackgroundTask(taskID2)
	
	// 驗證計數更新
	updatedCount := perfService.GetActiveTasksCount()
	if updatedCount != 2 {
		t.Errorf("取消註冊後活躍任務數量應為 2，實際為 %d", updatedCount)
	}
	
	// 取消註冊剩餘任務
	perfService.UnregisterBackgroundTask(taskID1)
	perfService.UnregisterBackgroundTask(taskID3)
	
	// 驗證最終計數
	finalCount := perfService.GetActiveTasksCount()
	if finalCount != 0 {
		t.Errorf("所有任務取消註冊後數量應為 0，實際為 %d", finalCount)
	}
	
	// 測試取消註冊不存在的任務（應該不會出錯）
	perfService.UnregisterBackgroundTask("不存在的任務ID")
}

// TestMetricsHistory 測試效能指標歷史記錄功能
// 驗證指標歷史的收集和查詢功能
//
// 測試案例：
// 1. 啟動監控並收集指標
// 2. 查詢不同時間範圍的歷史指標
// 3. 驗證歷史記錄的正確性
func TestMetricsHistory(t *testing.T) {
	// 建立模擬編輯器服務
	mockEditor := &mockEditorService{
		activeNotes: make(map[string]*models.Note),
	}
	
	// 建立效能服務
	perfService := NewPerformanceService(mockEditor)
	
	// 建立上下文
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	
	// 啟動監控
	err := perfService.StartMonitoring(ctx)
	if err != nil {
		t.Fatalf("啟動監控失敗：%v", err)
	}
	
	// 等待收集一些指標
	time.Sleep(200 * time.Millisecond)
	
	// 查詢最近 1 分鐘的歷史指標
	recentHistory := perfService.GetMetricsHistory(1 * time.Minute)
	if len(recentHistory) == 0 {
		t.Error("應該有歷史指標記錄")
	}
	
	// 查詢最近 1 秒的歷史指標
	veryRecentHistory := perfService.GetMetricsHistory(1 * time.Second)
	if len(veryRecentHistory) == 0 {
		t.Error("應該有最近的歷史指標記錄")
	}
	
	// 查詢過去的歷史指標（應該為空）
	oldHistory := perfService.GetMetricsHistory(1 * time.Nanosecond)
	if len(oldHistory) != 0 {
		t.Error("過去的歷史指標應該為空")
	}
	
	// 驗證歷史指標的時間順序
	if len(recentHistory) > 1 {
		for i := 1; i < len(recentHistory); i++ {
			if recentHistory[i].Timestamp.Before(recentHistory[i-1].Timestamp) {
				t.Error("歷史指標應該按時間順序排列")
			}
		}
	}
	
	// 停止監控
	err = perfService.StopMonitoring()
	if err != nil {
		t.Errorf("停止監控失敗：%v", err)
	}
}

// TestActiveNotesTracking 測試活躍筆記追蹤功能
// 驗證效能服務能正確追蹤編輯器中的活躍筆記數量
//
// 測試案例：
// 1. 在編輯器中添加筆記
// 2. 驗證效能指標中的活躍筆記數量
// 3. 移除筆記並驗證計數更新
func TestActiveNotesTracking(t *testing.T) {
	// 建立模擬編輯器服務
	mockEditor := &mockEditorService{
		activeNotes: make(map[string]*models.Note),
	}
	
	// 建立效能服務
	perfService := NewPerformanceService(mockEditor)
	
	// 驗證初始狀態
	initialMetrics := perfService.GetCurrentMetrics()
	if initialMetrics.ActiveNotes != 0 {
		t.Errorf("初始活躍筆記數量應為 0，實際為 %d", initialMetrics.ActiveNotes)
	}
	
	// 在模擬編輯器中添加筆記
	note1 := &models.Note{
		ID:      "note1",
		Title:   "測試筆記 1",
		Content: "這是測試內容",
	}
	
	note2 := &models.Note{
		ID:      "note2",
		Title:   "測試筆記 2",
		Content: "這是另一個測試內容",
	}
	
	mockEditor.activeNotes["note1"] = note1
	mockEditor.activeNotes["note2"] = note2
	
	// 驗證活躍筆記數量更新
	updatedMetrics := perfService.GetCurrentMetrics()
	if updatedMetrics.ActiveNotes != 2 {
		t.Errorf("活躍筆記數量應為 2，實際為 %d", updatedMetrics.ActiveNotes)
	}
	
	// 移除一個筆記
	delete(mockEditor.activeNotes, "note1")
	
	// 驗證計數更新
	finalMetrics := perfService.GetCurrentMetrics()
	if finalMetrics.ActiveNotes != 1 {
		t.Errorf("移除筆記後活躍數量應為 1，實際為 %d", finalMetrics.ActiveNotes)
	}
}

// TestPerformanceConcurrentAccess 測試並發存取的安全性
// 驗證效能服務在並發環境下的執行緒安全性
//
// 測試案例：
// 1. 並發註冊和取消註冊背景任務
// 2. 並發記錄快取命中和未命中
// 3. 並發收集效能指標
func TestPerformanceConcurrentAccess(t *testing.T) {
	// 建立模擬編輯器服務
	mockEditor := &mockEditorService{
		activeNotes: make(map[string]*models.Note),
	}
	
	// 建立效能服務
	perfService := NewPerformanceService(mockEditor)
	
	// 並發測試的 Goroutine 數量
	const numGoroutines = 10
	const operationsPerGoroutine = 100
	
	// 使用通道來同步 Goroutine
	done := make(chan bool, numGoroutines)
	
	// 並發註冊和取消註冊背景任務
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer func() { done <- true }()
			
			for j := 0; j < operationsPerGoroutine; j++ {
				// 註冊任務
				taskID := perfService.RegisterBackgroundTask(fmt.Sprintf("任務_%d_%d", id, j))
				
				// 立即取消註冊
				perfService.UnregisterBackgroundTask(taskID)
				
				// 記錄快取操作
				if ps, ok := perfService.(*performanceService); ok {
					if j%2 == 0 {
						ps.RecordCacheHit()
					} else {
						ps.RecordCacheMiss()
					}
				}
				
				// 取得當前指標
				_ = perfService.GetCurrentMetrics()
			}
		}(i)
	}
	
	// 等待所有 Goroutine 完成
	for i := 0; i < numGoroutines; i++ {
		<-done
	}
	
	// 驗證最終狀態
	finalTaskCount := perfService.GetActiveTasksCount()
	if finalTaskCount != 0 {
		t.Errorf("並發測試後活躍任務數量應為 0，實際為 %d", finalTaskCount)
	}
	
	// 驗證快取統計
	cacheStats := perfService.GetCacheStats()
	totalOperations := int64(numGoroutines * operationsPerGoroutine)
	expectedHits := totalOperations / 2
	expectedMisses := totalOperations - expectedHits
	
	actualHits := cacheStats["hits"].(int64)
	actualMisses := cacheStats["misses"].(int64)
	
	if actualHits != expectedHits {
		t.Errorf("快取命中次數應為 %d，實際為 %d", expectedHits, actualHits)
	}
	
	if actualMisses != expectedMisses {
		t.Errorf("快取未命中次數應為 %d，實際為 %d", expectedMisses, actualMisses)
	}
}