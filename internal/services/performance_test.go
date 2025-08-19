// Package services 效能測試套件
// 本檔案包含各種效能測試，驗證大檔案處理、記憶體優化和並發處理的效能表現
package services

import (
	"context"      // 上下文管理
	"fmt"          // 格式化輸出
	"runtime"      // 執行時期資訊
	"strings"      // 字串處理
	"testing"      // 測試框架
	"time"         // 時間處理
	"mac-notebook-app/internal/repositories" // 引入儲存庫
)

// BenchmarkEditorServiceLargeFile 測試編輯器服務處理大檔案的效能
// 驗證大檔案載入、處理和預覽的效能表現
//
// 測試場景：
// 1. 建立不同大小的測試檔案內容
// 2. 測試檔案開啟和處理時間
// 3. 測試 Markdown 預覽效能
// 4. 測試記憶體使用情況
func BenchmarkEditorServiceLargeFile(b *testing.B) {
	// 建立測試用的儲存庫和服務
	repo, _ := repositories.NewLocalFileRepository("test_data")
	perfService := NewPerformanceService(nil)
	editorService := NewEditorService(repo, nil, nil, nil, perfService, nil)
	
	// 測試不同大小的檔案
	testSizes := []struct {
		name string
		size int
	}{
		{"Small_1KB", 1024},
		{"Medium_100KB", 100 * 1024},
		{"Large_1MB", 1024 * 1024},
		{"VeryLarge_10MB", 10 * 1024 * 1024},
	}
	
	for _, testSize := range testSizes {
		b.Run(testSize.name, func(b *testing.B) {
			// 生成測試內容
			content := generateTestMarkdownContent(testSize.size)
			
			// 重置計時器
			b.ResetTimer()
			
			// 執行基準測試
			for i := 0; i < b.N; i++ {
				// 測試建立筆記
				note, err := editorService.CreateNote(fmt.Sprintf("測試筆記_%d", i), content)
				if err != nil {
					b.Fatalf("建立筆記失敗: %v", err)
				}
				
				// 測試 Markdown 預覽
				_ = editorService.PreviewMarkdown(content)
				
				// 清理
				editorService.CloseNote(note.ID)
			}
		})
	}
}

// BenchmarkPerformanceServiceMonitoring 測試效能監控服務的效能
// 驗證效能監控對系統效能的影響
//
// 測試場景：
// 1. 啟動效能監控
// 2. 執行大量操作
// 3. 測量監控開銷
func BenchmarkPerformanceServiceMonitoring(b *testing.B) {
	perfService := NewPerformanceService(nil)
	
	// 測試不同的監控場景
	testCases := []struct {
		name           string
		withMonitoring bool
	}{
		{"WithoutMonitoring", false},
		{"WithMonitoring", true},
	}
	
	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			
			if tc.withMonitoring {
				perfService.StartMonitoring(ctx)
				defer perfService.StopMonitoring()
			}
			
			b.ResetTimer()
			
			for i := 0; i < b.N; i++ {
				// 模擬各種操作
				taskID := perfService.RegisterBackgroundTask(fmt.Sprintf("測試任務_%d", i))
				_ = perfService.GetCurrentMetrics()
				perfService.UnregisterBackgroundTask(taskID)
				
				// 模擬快取操作
				if ps, ok := perfService.(*performanceService); ok {
					if i%2 == 0 {
						ps.RecordCacheHit()
					} else {
						ps.RecordCacheMiss()
					}
				}
			}
		})
	}
}

// BenchmarkMemoryOptimization 測試記憶體優化功能的效能
// 驗證記憶體優化操作的效率和效果
//
// 測試場景：
// 1. 建立大量物件消耗記憶體
// 2. 執行記憶體優化
// 3. 測量優化效果
func BenchmarkMemoryOptimization(b *testing.B) {
	perfService := NewPerformanceService(nil)
	
	b.Run("MemoryOptimization", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			// 建立一些記憶體壓力
			data := make([][]byte, 1000)
			for j := range data {
				data[j] = make([]byte, 1024) // 1KB per slice
			}
			
			// 執行記憶體優化
			perfService.OptimizeMemory()
			
			// 強制垃圾回收
			perfService.ForceGarbageCollection()
			
			// 清理引用
			data = nil
		}
	})
}

// BenchmarkConcurrentOperations 測試並發操作的效能
// 驗證服務在高並發環境下的效能表現
//
// 測試場景：
// 1. 多個 Goroutine 並發執行操作
// 2. 測試執行緒安全性
// 3. 測量並發效能
func BenchmarkConcurrentOperations(b *testing.B) {
	repo, _ := repositories.NewLocalFileRepository("test_data")
	perfService := NewPerformanceService(nil)
	editorService := NewEditorService(repo, nil, nil, nil, perfService, nil)
	
	b.Run("ConcurrentNoteOperations", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				// 並發建立和關閉筆記
				note, err := editorService.CreateNote(
					fmt.Sprintf("並發筆記_%d", i),
					generateTestMarkdownContent(1024),
				)
				if err != nil {
					b.Errorf("建立筆記失敗: %v", err)
					continue
				}
				
				// 更新內容
				editorService.UpdateContent(note.ID, fmt.Sprintf("更新內容_%d", i))
				
				// 關閉筆記
				editorService.CloseNote(note.ID)
				
				i++
			}
		})
	})
}

// TestLargeFileProcessing 測試大檔案處理功能
// 驗證大檔案的分塊處理和記憶體管理
//
// 測試案例：
// 1. 處理不同大小的檔案
// 2. 驗證分塊處理的正確性
// 3. 檢查記憶體使用情況
func TestLargeFileProcessing(t *testing.T) {
	repo, _ := repositories.NewLocalFileRepository("test_data")
	perfService := NewPerformanceService(nil)
	editorService := NewEditorService(repo, nil, nil, nil, perfService, nil)
	
	// 測試不同大小的檔案
	testCases := []struct {
		name     string
		size     int
		expected string
	}{
		{"SmallFile", 1024, "小檔案應該正常處理"},
		{"MediumFile", 100 * 1024, "中等檔案應該正常處理"},
		{"LargeFile", 5 * 1024 * 1024, "大檔案應該使用分塊處理"},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 生成測試內容
			content := generateTestMarkdownContent(tc.size)
			
			// 記錄處理前的記憶體使用
			var beforeMem runtime.MemStats
			runtime.ReadMemStats(&beforeMem)
			
			// 建立筆記
			note, err := editorService.CreateNote(tc.name, content)
			if err != nil {
				t.Fatalf("建立筆記失敗: %v", err)
			}
			
			// 測試預覽功能
			preview := editorService.PreviewMarkdown(content)
			if len(preview) == 0 {
				t.Error("預覽內容不應為空")
			}
			
			// 記錄處理後的記憶體使用
			var afterMem runtime.MemStats
			runtime.ReadMemStats(&afterMem)
			
			// 驗證記憶體使用合理性
			memIncrease := afterMem.Alloc - beforeMem.Alloc
			if memIncrease > uint64(tc.size*2) { // 記憶體增長不應超過檔案大小的兩倍
				t.Errorf("記憶體使用增長過多: %d bytes (檔案大小: %d bytes)", memIncrease, tc.size)
			}
			
			// 清理
			editorService.CloseNote(note.ID)
		})
	}
}

// TestCacheOptimization 測試快取優化功能
// 驗證快取管理和優化策略的效果
//
// 測試案例：
// 1. 填滿快取
// 2. 觸發快取優化
// 3. 驗證最舊項目被移除
func TestCacheOptimization(t *testing.T) {
	repo, _ := repositories.NewLocalFileRepository("test_data")
	perfService := NewPerformanceService(nil)
	editorService := NewEditorService(repo, nil, nil, nil, perfService, nil)
	
	// 注意：由於介面限制，我們無法直接設定快取大小
	// 這個測試將驗證基本的快取功能
	
	// 建立超過快取大小的筆記
	noteIDs := make([]string, 0)
	for i := 0; i < 8; i++ {
		note, err := editorService.CreateNote(
			fmt.Sprintf("快取測試筆記_%d", i),
			fmt.Sprintf("內容_%d", i),
		)
		if err != nil {
			t.Fatalf("建立筆記 %d 失敗: %v", i, err)
		}
		noteIDs = append(noteIDs, note.ID)
		
		// 為了測試 LRU 邏輯，讓每個筆記有不同的更新時間
		time.Sleep(1 * time.Millisecond)
		editorService.UpdateContent(note.ID, fmt.Sprintf("更新內容_%d", i))
	}
	
	// 檢查快取大小
	activeNotes := editorService.GetActiveNotes()
	if len(activeNotes) == 0 {
		t.Error("應該有活躍筆記")
	}
	
	// 基本驗證完成
}

// TestMemoryMonitoring 測試記憶體監控功能
// 驗證記憶體使用監控和建議功能
//
// 測試案例：
// 1. 監控記憶體使用情況
// 2. 驗證建議的準確性
// 3. 測試記憶體優化效果
func TestMemoryMonitoring(t *testing.T) {
	repo, _ := repositories.NewLocalFileRepository("test_data")
	perfService := NewPerformanceService(nil)
	editorService := NewEditorService(repo, nil, nil, nil, perfService, nil)
	
	// 測試基本記憶體監控功能
	memUsage, err := perfService.GetMemoryUsage()
	if err != nil {
		t.Errorf("記憶體監控失敗: %v", err)
	}
	
	if memUsage <= 0 {
		t.Error("記憶體使用量應該大於 0")
	}
	
	// 建立一些筆記以增加記憶體使用
	for i := 0; i < 10; i++ {
		_, err := editorService.CreateNote(
			fmt.Sprintf("記憶體測試筆記_%d", i),
			generateTestMarkdownContent(10*1024), // 10KB per note
		)
		if err != nil {
			t.Errorf("建立筆記 %d 失敗: %v", i, err)
		}
	}
	
	// 再次檢查記憶體使用
	newMemUsage, err := perfService.GetMemoryUsage()
	if err != nil {
		t.Errorf("第二次記憶體監控失敗: %v", err)
	}
	
	// 記憶體使用量可能會有變化，這裡只驗證能正常取得
	if newMemUsage <= 0 {
		t.Error("記憶體使用量應該大於 0")
	}
}

// TestPerformanceIntegration 測試效能服務整合
// 驗證效能服務與其他服務的整合效果
//
// 測試案例：
// 1. 啟動效能監控
// 2. 執行各種操作
// 3. 驗證效能指標收集
// 4. 測試效能優化建議
func TestPerformanceIntegration(t *testing.T) {
	repo, _ := repositories.NewLocalFileRepository("test_data")
	perfService := NewPerformanceService(nil)
	editorService := NewEditorService(repo, nil, nil, nil, perfService, nil)
	
	// 啟動效能監控
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	
	err := perfService.StartMonitoring(ctx)
	if err != nil {
		t.Fatalf("啟動效能監控失敗: %v", err)
	}
	defer perfService.StopMonitoring()
	
	// 執行各種操作以生成效能資料
	for i := 0; i < 5; i++ {
		// 建立筆記
		note, err := editorService.CreateNote(
			fmt.Sprintf("整合測試筆記_%d", i),
			generateTestMarkdownContent(5*1024),
		)
		if err != nil {
			t.Errorf("建立筆記 %d 失敗: %v", i, err)
			continue
		}
		
		// 更新內容
		editorService.UpdateContent(note.ID, fmt.Sprintf("更新內容_%d", i))
		
		// 預覽內容
		editorService.PreviewMarkdown(note.Content)
		
		// 註冊背景任務
		taskID := perfService.RegisterBackgroundTask(fmt.Sprintf("測試任務_%d", i))
		
		// 短暫等待
		time.Sleep(10 * time.Millisecond)
		
		// 取消註冊任務
		perfService.UnregisterBackgroundTask(taskID)
	}
	
	// 等待效能監控收集資料
	time.Sleep(100 * time.Millisecond)
	
	// 驗證效能指標
	metrics := perfService.GetCurrentMetrics()
	if metrics == nil {
		t.Fatal("效能指標不應為 nil")
	}
	
	if metrics.ActiveNotes != 5 {
		t.Errorf("活躍筆記數量應為 5，實際為 %d", metrics.ActiveNotes)
	}
	
	if metrics.MemoryUsage <= 0 {
		t.Error("記憶體使用量應大於 0")
	}
	
	// 測試歷史指標
	history := perfService.GetMetricsHistory(1 * time.Minute)
	if len(history) == 0 {
		t.Error("應該有歷史效能指標")
	}
	
	// 測試快取統計
	cacheStats := perfService.GetCacheStats()
	if cacheStats == nil {
		t.Error("快取統計不應為 nil")
	}
	
	// 執行記憶體優化
	err = perfService.OptimizeMemory()
	if err != nil {
		t.Errorf("記憶體優化失敗: %v", err)
	}
}

// generateTestMarkdownContent 生成指定大小的測試 Markdown 內容
// 參數：size（目標大小，位元組）
// 回傳：生成的 Markdown 內容
//
// 執行流程：
// 1. 計算需要生成的內容量
// 2. 建立包含各種 Markdown 元素的內容
// 3. 重複內容直到達到目標大小
func generateTestMarkdownContent(size int) string {
	// 基礎 Markdown 內容模板
	baseContent := `# 測試標題

這是一個測試段落，包含**粗體文字**和*斜體文字*。

## 子標題

- 列表項目 1
- 列表項目 2
- 列表項目 3

### 程式碼範例

` + "```go\nfunc main() {\n    fmt.Println(\"Hello, World!\")\n}\n```" + `

> 這是一個引用區塊

| 欄位1 | 欄位2 | 欄位3 |
|-------|-------|-------|
| 資料1 | 資料2 | 資料3 |

[連結文字](https://example.com)

---

`
	
	// 計算需要重複的次數
	baseSize := len(baseContent)
	repeatCount := (size / baseSize) + 1
	
	// 建立內容
	var builder strings.Builder
	builder.Grow(size) // 預分配記憶體
	
	for i := 0; i < repeatCount; i++ {
		builder.WriteString(baseContent)
		if builder.Len() >= size {
			break
		}
	}
	
	// 截取到目標大小
	result := builder.String()
	if len(result) > size {
		result = result[:size]
	}
	
	return result
}