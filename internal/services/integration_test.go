// Package services 整合測試套件
// 本檔案包含端到端的整合測試，驗證各個服務之間的協作和完整的使用場景
package services

import (
	"context"      // 上下文管理
	"fmt"          // 格式化輸出
	"os"           // 作業系統介面
	"path/filepath" // 檔案路徑處理
	"strings"      // 字串處理
	"testing"      // 測試框架
	"time"         // 時間處理
	"mac-notebook-app/internal/models"      // 引入資料模型
	"mac-notebook-app/internal/repositories" // 引入儲存庫
)

// IntegrationTestSuite 整合測試套件
// 包含所有必要的服務實例和測試資料
type IntegrationTestSuite struct {
	// 服務實例
	fileRepo        repositories.FileRepository
	fileManager     FileManagerService
	editorService   EditorService
	encryptionSvc   EncryptionService
	passwordSvc     PasswordService
	biometricSvc    BiometricService
	autoSaveService AutoSaveService
	perfService     PerformanceService
	errorService    ErrorService
	notificationSvc NotificationService
	
	// 測試環境
	testDir         string
	testNotes       []*models.Note
	cleanup         func()
}

// setupIntegrationTest 設定整合測試環境
// 回傳：配置完成的測試套件和清理函數
//
// 執行流程：
// 1. 建立臨時測試目錄
// 2. 初始化所有服務實例
// 3. 設定服務之間的依賴關係
// 4. 準備測試資料
// 5. 回傳測試套件和清理函數
func setupIntegrationTest(t *testing.T) *IntegrationTestSuite {
	// 建立臨時測試目錄
	testDir, err := os.MkdirTemp("", "notebook_integration_test_*")
	if err != nil {
		t.Fatalf("建立測試目錄失敗: %v", err)
	}
	
	// 初始化檔案儲存庫
	fileRepo, err := repositories.NewLocalFileRepository(testDir)
	if err != nil {
		os.RemoveAll(testDir)
		t.Fatalf("建立檔案儲存庫失敗: %v", err)
	}
	
	// 初始化各種服務
	fileManager, err := NewLocalFileManagerService(fileRepo, testDir)
	if err != nil {
		os.RemoveAll(testDir)
		t.Fatalf("建立檔案管理服務失敗: %v", err)
	}
	
	encryptionSvc := NewEncryptionService()
	passwordSvc := NewPasswordService()
	biometricSvc := NewBiometricService()
	errorService, _ := NewErrorService("test.log")
	notificationSvc := NewNotificationService()
	perfService := NewPerformanceService(nil) // 暫時不設定編輯器服務
	
	// 建立編輯器服務（需要在效能服務之後建立以避免循環依賴）
	editorService := NewEditorService(fileRepo, encryptionSvc, passwordSvc, biometricSvc, perfService)
	
	// 更新效能服務的編輯器服務引用
	if ps, ok := perfService.(*performanceService); ok {
		ps.editorService = editorService
	}
	
	// 建立自動保存服務
	autoSaveService := NewAutoSaveServiceWithDefaults(editorService)
	
	// 建立清理函數
	cleanup := func() {
		// 停止所有背景服務
		if perfService != nil {
			perfService.StopMonitoring()
		}
		
		// 清理測試目錄
		os.RemoveAll(testDir)
	}
	
	return &IntegrationTestSuite{
		fileRepo:        fileRepo,
		fileManager:     fileManager,
		editorService:   editorService,
		encryptionSvc:   encryptionSvc,
		passwordSvc:     passwordSvc,
		biometricSvc:    biometricSvc,
		autoSaveService: autoSaveService,
		perfService:     perfService,
		errorService:    errorService,
		notificationSvc: notificationSvc,
		testDir:         testDir,
		testNotes:       make([]*models.Note, 0),
		cleanup:         cleanup,
	}
}

// TestEndToEndNoteCreationAndEditing 測試端到端的筆記建立和編輯流程
// 驗證完整的筆記生命週期，從建立到保存的所有步驟
//
// 測試場景：
// 1. 建立新筆記
// 2. 編輯筆記內容
// 3. 自動保存功能
// 4. 檔案系統操作
// 5. 效能監控
func TestEndToEndNoteCreationAndEditing(t *testing.T) {
	suite := setupIntegrationTest(t)
	defer suite.cleanup()
	
	// 啟動效能監控
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	err := suite.perfService.StartMonitoring(ctx)
	if err != nil {
		t.Errorf("啟動效能監控失敗: %v", err)
	}
	defer suite.perfService.StopMonitoring()
	
	// 步驟 1: 建立新筆記
	t.Log("步驟 1: 建立新筆記")
	note, err := suite.editorService.CreateNote("整合測試筆記", "# 測試標題\n\n這是測試內容。")
	if err != nil {
		t.Fatalf("建立筆記失敗: %v", err)
	}
	suite.testNotes = append(suite.testNotes, note)
	
	// 驗證筆記建立
	if note.ID == "" {
		t.Error("筆記 ID 不應為空")
	}
	if note.Title != "整合測試筆記" {
		t.Errorf("筆記標題不正確，期望：整合測試筆記，實際：%s", note.Title)
	}
	
	// 步驟 2: 設定筆記檔案路徑並保存
	t.Log("步驟 2: 保存筆記到檔案系統")
	note.FilePath = "integration_test_note.md"
	err = suite.editorService.SaveNote(note)
	if err != nil {
		t.Fatalf("保存筆記失敗: %v", err)
	}
	
	// 驗證檔案是否存在
	fullPath := filepath.Join(suite.testDir, note.FilePath)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		t.Error("筆記檔案應該存在於檔案系統中")
	}
	
	// 步驟 3: 編輯筆記內容
	t.Log("步驟 3: 編輯筆記內容")
	newContent := "# 更新的標題\n\n這是更新後的內容。\n\n- 項目 1\n- 項目 2"
	err = suite.editorService.UpdateContent(note.ID, newContent)
	if err != nil {
		t.Fatalf("更新筆記內容失敗: %v", err)
	}
	
	// 驗證內容更新
	updatedNote, exists := suite.editorService.GetActiveNote(note.ID)
	if !exists {
		t.Fatal("更新後的筆記應該存在於活躍筆記中")
	}
	if updatedNote.Content != newContent {
		t.Error("筆記內容未正確更新")
	}
	
	// 步驟 4: 測試 Markdown 預覽
	t.Log("步驟 4: 測試 Markdown 預覽")
	preview := suite.editorService.PreviewMarkdown(newContent)
	if !strings.Contains(preview, "<h1>") {
		t.Error("Markdown 預覽應該包含 HTML 標題標籤")
	}
	if !strings.Contains(preview, "<li>") {
		t.Error("Markdown 預覽應該包含 HTML 列表標籤")
	}
	
	// 步驟 5: 啟動自動保存
	t.Log("步驟 5: 測試自動保存功能")
	suite.autoSaveService.StartAutoSave(updatedNote, 100*time.Millisecond)
	
	// 等待自動保存執行
	time.Sleep(200 * time.Millisecond)
	
	// 檢查保存狀態
	saveStatus := suite.autoSaveService.GetSaveStatus(note.ID)
	if saveStatus.NoteID != note.ID {
		t.Error("保存狀態的筆記 ID 不正確")
	}
	
	// 停止自動保存
	suite.autoSaveService.StopAutoSave(note.ID)
	
	// 步驟 6: 驗證效能指標
	t.Log("步驟 6: 驗證效能指標")
	metrics := suite.perfService.GetCurrentMetrics()
	if metrics == nil {
		t.Error("效能指標不應為 nil")
	}
	if metrics.MemoryUsage <= 0 {
		t.Error("記憶體使用量應該大於 0")
	}
	
	// 步驟 7: 測試檔案管理操作
	t.Log("步驟 7: 測試檔案管理操作")
	files, err := suite.fileManager.ListFiles(".")
	if err != nil {
		t.Errorf("列出檔案失敗: %v", err)
	}
	
	// 驗證筆記檔案在列表中
	found := false
	for _, file := range files {
		if file.Name == "integration_test_note.md" {
			found = true
			break
		}
	}
	if !found {
		t.Error("筆記檔案應該出現在檔案列表中")
	}
	
	t.Log("端到端測試完成")
}

// TestEncryptionIntegration 測試加密功能的整合
// 驗證加密筆記的完整工作流程
//
// 測試場景：
// 1. 建立普通筆記
// 2. 啟用加密
// 3. 保存加密筆記
// 4. 重新開啟加密筆記
// 5. 解密和編輯
func TestEncryptionIntegration(t *testing.T) {
	suite := setupIntegrationTest(t)
	defer suite.cleanup()
	
	// 步驟 1: 建立普通筆記
	t.Log("步驟 1: 建立普通筆記")
	note, err := suite.editorService.CreateNote("加密測試筆記", "這是需要加密的敏感內容。")
	if err != nil {
		t.Fatalf("建立筆記失敗: %v", err)
	}
	
	// 步驟 2: 啟用加密
	t.Log("步驟 2: 啟用加密")
	testPassword := "TestPassword123!"
	
	// 注意：密碼服務的 ValidatePassword 方法尚未實作，跳過驗證
	
	// 為筆記啟用加密（透過編輯器服務的具體實作）
	if es, ok := suite.editorService.(*editorService); ok {
		err = es.EnableEncryption(note.ID, testPassword, "aes256", false)
		if err != nil {
			t.Fatalf("啟用加密失敗: %v", err)
		}
		
		// 驗證加密狀態
		isEncrypted, exists := es.IsEncrypted(note.ID)
		if !exists {
			t.Fatal("筆記應該存在")
		}
		if !isEncrypted {
			t.Error("筆記應該已加密")
		}
		
		// 步驟 3: 保存加密筆記
		t.Log("步驟 3: 保存加密筆記")
		note.FilePath = "encrypted_test_note.md.enc"
		
		// 使用密碼加密並保存
		err = es.EncryptWithPassword(note.ID, testPassword)
		if err != nil {
			t.Errorf("加密筆記失敗: %v", err)
		}
		
		err = suite.editorService.SaveNote(note)
		if err != nil {
			t.Errorf("保存加密筆記失敗: %v", err)
		}
		
		// 步驟 4: 解密筆記
		t.Log("步驟 4: 解密筆記")
		decryptedContent, err := es.DecryptWithPassword(note.ID, testPassword)
		if err != nil {
			t.Errorf("解密筆記失敗: %v", err)
		}
		
		if decryptedContent != note.Content {
			t.Error("解密後的內容應該與原始內容相同")
		}
		
		// 步驟 5: 測試錯誤密碼
		t.Log("步驟 5: 測試錯誤密碼")
		_, err = es.DecryptWithPassword(note.ID, "WrongPassword")
		if err == nil {
			t.Error("使用錯誤密碼解密應該失敗")
		}
	} else {
		t.Skip("無法存取編輯器服務的具體實作，跳過加密測試")
	}
	
	t.Log("加密整合測試完成")
}

// TestFileManagementIntegration 測試檔案管理的整合
// 驗證檔案和目錄操作的完整工作流程
//
// 測試場景：
// 1. 建立目錄結構
// 2. 建立和保存多個筆記
// 3. 檔案操作（複製、移動、重新命名）
// 4. 目錄遍歷和搜尋
func TestFileManagementIntegration(t *testing.T) {
	suite := setupIntegrationTest(t)
	defer suite.cleanup()
	
	// 步驟 1: 建立目錄結構
	t.Log("步驟 1: 建立目錄結構")
	directories := []string{"projects", "personal", "archive"}
	for _, dir := range directories {
		err := suite.fileManager.CreateDirectory(dir)
		if err != nil {
			t.Errorf("建立目錄 %s 失敗: %v", dir, err)
		}
	}
	
	// 步驟 2: 在不同目錄中建立筆記
	t.Log("步驟 2: 建立多個筆記")
	testNotes := []struct {
		title    string
		content  string
		filePath string
	}{
		{"專案筆記", "# 專案計劃\n\n這是專案相關的筆記。", "projects/project_note.md"},
		{"個人筆記", "# 個人想法\n\n這是個人筆記。", "personal/personal_note.md"},
		{"封存筆記", "# 舊資料\n\n這是封存的筆記。", "archive/old_note.md"},
	}
	
	for _, noteData := range testNotes {
		note, err := suite.editorService.CreateNote(noteData.title, noteData.content)
		if err != nil {
			t.Errorf("建立筆記 %s 失敗: %v", noteData.title, err)
			continue
		}
		
		note.FilePath = noteData.filePath
		err = suite.editorService.SaveNote(note)
		if err != nil {
			t.Errorf("保存筆記 %s 失敗: %v", noteData.title, err)
		}
		
		suite.testNotes = append(suite.testNotes, note)
	}
	
	// 步驟 3: 驗證目錄內容
	t.Log("步驟 3: 驗證目錄內容")
	for _, dir := range directories {
		files, err := suite.fileManager.ListFiles(dir)
		if err != nil {
			t.Errorf("列出目錄 %s 的檔案失敗: %v", dir, err)
			continue
		}
		
		if len(files) != 1 {
			t.Errorf("目錄 %s 應該包含 1 個檔案，實際有 %d 個", dir, len(files))
		}
	}
	
	// 步驟 4: 測試檔案複製
	t.Log("步驟 4: 測試檔案複製")
	err := suite.fileManager.CopyFile("projects/project_note.md", "projects/project_note_backup.md")
	if err != nil {
		t.Errorf("複製檔案失敗: %v", err)
	}
	
	// 驗證複製結果
	files, err := suite.fileManager.ListFiles("projects")
	if err != nil {
		t.Errorf("列出 projects 目錄失敗: %v", err)
	} else if len(files) != 2 {
		t.Errorf("複製後 projects 目錄應該有 2 個檔案，實際有 %d 個", len(files))
	}
	
	// 步驟 5: 測試檔案移動
	t.Log("步驟 5: 測試檔案移動")
	err = suite.fileManager.MoveFile("archive/old_note.md", "personal/moved_note.md")
	if err != nil {
		t.Errorf("移動檔案失敗: %v", err)
	}
	
	// 驗證移動結果
	archiveFiles, err := suite.fileManager.ListFiles("archive")
	if err != nil {
		t.Errorf("列出 archive 目錄失敗: %v", err)
	} else if len(archiveFiles) != 0 {
		t.Errorf("移動後 archive 目錄應該為空，實際有 %d 個檔案", len(archiveFiles))
	}
	
	personalFiles, err := suite.fileManager.ListFiles("personal")
	if err != nil {
		t.Errorf("列出 personal 目錄失敗: %v", err)
	} else if len(personalFiles) != 2 {
		t.Errorf("移動後 personal 目錄應該有 2 個檔案，實際有 %d 個", len(personalFiles))
	}
	
	// 步驟 6: 測試檔案重新命名
	t.Log("步驟 6: 測試檔案重新命名")
	err = suite.fileManager.RenameFile("personal/moved_note.md", "personal/renamed_note.md")
	if err != nil {
		t.Errorf("重新命名檔案失敗: %v", err)
	}
	
	// 步驟 7: 測試檔案搜尋
	t.Log("步驟 7: 測試檔案搜尋")
	searchResults, err := suite.fileManager.SearchFiles(".", "*.md", true)
	if err != nil {
		t.Errorf("搜尋檔案失敗: %v", err)
	}
	
	if len(searchResults) < 3 {
		t.Errorf("搜尋結果應該至少有 3 個 .md 檔案，實際找到 %d 個", len(searchResults))
	}
	
	t.Log("檔案管理整合測試完成")
}

// TestPerformanceUnderLoad 測試系統在負載下的效能表現
// 驗證系統在高負載情況下的穩定性和效能
//
// 測試場景：
// 1. 建立大量筆記
// 2. 並發操作
// 3. 記憶體使用監控
// 4. 效能指標收集
func TestPerformanceUnderLoad(t *testing.T) {
	suite := setupIntegrationTest(t)
	defer suite.cleanup()
	
	// 啟動效能監控
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	
	err := suite.perfService.StartMonitoring(ctx)
	if err != nil {
		t.Errorf("啟動效能監控失敗: %v", err)
	}
	defer suite.perfService.StopMonitoring()
	
	// 步驟 1: 建立大量筆記
	t.Log("步驟 1: 建立大量筆記")
	noteCount := 50
	notes := make([]*models.Note, 0, noteCount)
	
	startTime := time.Now()
	for i := 0; i < noteCount; i++ {
		title := fmt.Sprintf("負載測試筆記_%d", i)
		content := fmt.Sprintf("# 筆記 %d\n\n這是負載測試的內容。\n\n", i)
		
		// 為一些筆記添加更多內容以模擬不同大小
		if i%10 == 0 {
			for j := 0; j < 100; j++ {
				content += fmt.Sprintf("這是第 %d 行額外內容。\n", j)
			}
		}
		
		note, err := suite.editorService.CreateNote(title, content)
		if err != nil {
			t.Errorf("建立筆記 %d 失敗: %v", i, err)
			continue
		}
		
		note.FilePath = fmt.Sprintf("load_test_note_%d.md", i)
		notes = append(notes, note)
	}
	creationTime := time.Since(startTime)
	t.Logf("建立 %d 個筆記耗時: %v", noteCount, creationTime)
	
	// 步驟 2: 批量保存筆記
	t.Log("步驟 2: 批量保存筆記")
	startTime = time.Now()
	for _, note := range notes {
		err := suite.editorService.SaveNote(note)
		if err != nil {
			t.Errorf("保存筆記 %s 失敗: %v", note.Title, err)
		}
	}
	saveTime := time.Since(startTime)
	t.Logf("保存 %d 個筆記耗時: %v", len(notes), saveTime)
	
	// 步驟 3: 並發編輯操作
	t.Log("步驟 3: 並發編輯操作")
	startTime = time.Now()
	
	// 使用通道來同步並發操作
	done := make(chan bool, len(notes))
	
	for i, note := range notes {
		go func(index int, n *models.Note) {
			defer func() { done <- true }()
			
			// 更新內容
			newContent := fmt.Sprintf("# 更新的筆記 %d\n\n這是並發更新的內容。", index)
			err := suite.editorService.UpdateContent(n.ID, newContent)
			if err != nil {
				t.Errorf("並發更新筆記 %d 失敗: %v", index, err)
			}
			
			// 生成預覽
			_ = suite.editorService.PreviewMarkdown(newContent)
		}(i, note)
	}
	
	// 等待所有並發操作完成
	for i := 0; i < len(notes); i++ {
		<-done
	}
	concurrentTime := time.Since(startTime)
	t.Logf("並發編輯 %d 個筆記耗時: %v", len(notes), concurrentTime)
	
	// 步驟 4: 檢查效能指標
	t.Log("步驟 4: 檢查效能指標")
	metrics := suite.perfService.GetCurrentMetrics()
	if metrics == nil {
		t.Error("效能指標不應為 nil")
	} else {
		t.Logf("記憶體使用量: %d bytes", metrics.MemoryUsage)
		t.Logf("Goroutine 數量: %d", metrics.Goroutines)
		t.Logf("垃圾回收次數: %d", metrics.GCCount)
		
		// 驗證記憶體使用合理性
		if metrics.MemoryUsage > 100*1024*1024 { // 100MB
			t.Logf("警告：記憶體使用量較高: %d bytes", metrics.MemoryUsage)
		}
	}
	
	// 步驟 5: 執行記憶體優化
	t.Log("步驟 5: 執行記憶體優化")
	err = suite.perfService.OptimizeMemory()
	if err != nil {
		t.Errorf("記憶體優化失敗: %v", err)
	}
	
	// 檢查優化後的記憶體使用
	optimizedMetrics := suite.perfService.GetCurrentMetrics()
	if optimizedMetrics != nil {
		t.Logf("優化後記憶體使用量: %d bytes", optimizedMetrics.MemoryUsage)
	}
	
	// 步驟 6: 清理筆記
	t.Log("步驟 6: 清理筆記")
	for _, note := range notes {
		suite.editorService.CloseNote(note.ID)
	}
	
	t.Log("負載測試完成")
}

// TestLongRunningStability 測試長時間運行的穩定性
// 驗證系統在長時間運行下的穩定性和資源管理
//
// 測試場景：
// 1. 長時間效能監控
// 2. 週期性操作
// 3. 記憶體洩漏檢測
// 4. 資源清理驗證
func TestLongRunningStability(t *testing.T) {
	if testing.Short() {
		t.Skip("跳過長時間運行測試（使用 -short 標誌）")
	}
	
	suite := setupIntegrationTest(t)
	defer suite.cleanup()
	
	// 啟動效能監控
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	err := suite.perfService.StartMonitoring(ctx)
	if err != nil {
		t.Errorf("啟動效能監控失敗: %v", err)
	}
	defer suite.perfService.StopMonitoring()
	
	// 記錄初始記憶體使用
	initialMetrics := suite.perfService.GetCurrentMetrics()
	if initialMetrics == nil {
		t.Fatal("無法取得初始效能指標")
	}
	initialMemory := initialMetrics.MemoryUsage
	
	t.Logf("初始記憶體使用量: %d bytes", initialMemory)
	
	// 執行週期性操作
	cycles := 20
	for i := 0; i < cycles; i++ {
		t.Logf("執行週期 %d/%d", i+1, cycles)
		
		// 建立筆記
		note, err := suite.editorService.CreateNote(
			fmt.Sprintf("穩定性測試筆記_%d", i),
			fmt.Sprintf("這是第 %d 個穩定性測試筆記。", i),
		)
		if err != nil {
			t.Errorf("週期 %d: 建立筆記失敗: %v", i, err)
			continue
		}
		
		// 編輯內容
		newContent := fmt.Sprintf("# 更新的筆記 %d\n\n更新時間: %s", i, time.Now().Format(time.RFC3339))
		err = suite.editorService.UpdateContent(note.ID, newContent)
		if err != nil {
			t.Errorf("週期 %d: 更新內容失敗: %v", i, err)
		}
		
		// 生成預覽
		_ = suite.editorService.PreviewMarkdown(newContent)
		
		// 保存筆記
		note.FilePath = fmt.Sprintf("stability_test_%d.md", i)
		err = suite.editorService.SaveNote(note)
		if err != nil {
			t.Errorf("週期 %d: 保存筆記失敗: %v", i, err)
		}
		
		// 關閉筆記以釋放記憶體
		suite.editorService.CloseNote(note.ID)
		
		// 每 5 個週期執行一次記憶體優化
		if (i+1)%5 == 0 {
			err = suite.perfService.OptimizeMemory()
			if err != nil {
				t.Errorf("週期 %d: 記憶體優化失敗: %v", i, err)
			}
			
			// 檢查記憶體使用
			currentMetrics := suite.perfService.GetCurrentMetrics()
			if currentMetrics != nil {
				t.Logf("週期 %d 記憶體使用量: %d bytes", i+1, currentMetrics.MemoryUsage)
				
				// 檢查記憶體洩漏
				memoryIncrease := currentMetrics.MemoryUsage - initialMemory
				if memoryIncrease > 50*1024*1024 { // 50MB
					t.Logf("警告：記憶體使用量增長較大: %d bytes", memoryIncrease)
				}
			}
		}
		
		// 短暫等待
		time.Sleep(100 * time.Millisecond)
	}
	
	// 最終記憶體檢查
	finalMetrics := suite.perfService.GetCurrentMetrics()
	if finalMetrics != nil {
		t.Logf("最終記憶體使用量: %d bytes", finalMetrics.MemoryUsage)
		
		memoryIncrease := finalMetrics.MemoryUsage - initialMemory
		t.Logf("記憶體增長: %d bytes", memoryIncrease)
		
		// 驗證記憶體增長在合理範圍內
		if memoryIncrease > 100*1024*1024 { // 100MB
			t.Errorf("記憶體增長過大，可能存在記憶體洩漏: %d bytes", memoryIncrease)
		}
	}
	
	// 檢查效能歷史
	history := suite.perfService.GetMetricsHistory(10 * time.Second)
	if len(history) == 0 {
		t.Error("應該有效能歷史記錄")
	} else {
		t.Logf("收集到 %d 筆效能歷史記錄", len(history))
	}
	
	t.Log("長時間穩定性測試完成")
}

// TestErrorHandlingIntegration 測試錯誤處理的整合
// 驗證各種錯誤情況下系統的行為和恢復能力
//
// 測試場景：
// 1. 檔案系統錯誤
// 2. 加密錯誤
// 3. 記憶體不足模擬
// 4. 並發錯誤處理
func TestErrorHandlingIntegration(t *testing.T) {
	suite := setupIntegrationTest(t)
	defer suite.cleanup()
	
	// 步驟 1: 測試檔案系統錯誤
	t.Log("步驟 1: 測試檔案系統錯誤")
	
	// 嘗試在不存在的目錄中建立檔案
	note, err := suite.editorService.CreateNote("錯誤測試筆記", "測試內容")
	if err != nil {
		t.Errorf("建立筆記失敗: %v", err)
	} else {
		note.FilePath = "nonexistent/directory/test.md"
		err = suite.editorService.SaveNote(note)
		if err == nil {
			t.Error("在不存在的目錄中保存檔案應該失敗")
		} else {
			t.Logf("預期的檔案系統錯誤: %v", err)
		}
	}
	
	// 步驟 2: 測試無效路徑
	t.Log("步驟 2: 測試無效路徑")
	_, err = suite.fileManager.ListFiles("../../../etc")
	if err == nil {
		t.Error("存取系統目錄應該失敗")
	} else {
		t.Logf("預期的路徑驗證錯誤: %v", err)
	}
	
	// 步驟 3: 測試加密錯誤
	t.Log("步驟 3: 測試加密錯誤")
	if es, ok := suite.editorService.(*editorService); ok {
		// 嘗試使用弱密碼
		weakPassword := "123"
		err = es.EnableEncryption(note.ID, weakPassword, "aes256", false)
		if err == nil {
			t.Error("使用弱密碼啟用加密應該失敗")
		} else {
			t.Logf("預期的密碼強度錯誤: %v", err)
		}
	}
	
	// 步驟 4: 測試並發錯誤處理
	t.Log("步驟 4: 測試並發錯誤處理")
	errorCount := 0
	done := make(chan bool, 10)
	
	for i := 0; i < 10; i++ {
		go func(index int) {
			defer func() { done <- true }()
			
			// 嘗試操作不存在的筆記
			err := suite.editorService.UpdateContent("nonexistent_note_id", "test content")
			if err != nil {
				errorCount++
			}
		}(i)
	}
	
	// 等待所有並發操作完成
	for i := 0; i < 10; i++ {
		<-done
	}
	
	if errorCount != 10 {
		t.Errorf("應該有 10 個錯誤，實際有 %d 個", errorCount)
	}
	
	// 步驟 5: 測試錯誤恢復
	t.Log("步驟 5: 測試錯誤恢復")
	
	// 在錯誤後系統應該仍能正常工作
	recoveryNote, err := suite.editorService.CreateNote("恢復測試筆記", "系統恢復測試")
	if err != nil {
		t.Errorf("錯誤後建立筆記失敗: %v", err)
	} else {
		recoveryNote.FilePath = "recovery_test.md"
		err = suite.editorService.SaveNote(recoveryNote)
		if err != nil {
			t.Errorf("錯誤後保存筆記失敗: %v", err)
		}
	}
	
	t.Log("錯誤處理整合測試完成")
}

// TestCompleteWorkflow 測試完整的工作流程
// 驗證從應用程式啟動到關閉的完整使用場景
//
// 測試場景：
// 1. 系統初始化
// 2. 用戶工作流程模擬
// 3. 多種功能組合使用
// 4. 系統關閉和清理
func TestCompleteWorkflow(t *testing.T) {
	suite := setupIntegrationTest(t)
	defer suite.cleanup()
	
	// 步驟 1: 系統初始化
	t.Log("步驟 1: 系統初始化")
	
	// 啟動效能監控
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	err := suite.perfService.StartMonitoring(ctx)
	if err != nil {
		t.Errorf("啟動效能監控失敗: %v", err)
	}
	defer suite.perfService.StopMonitoring()
	
	// 建立工作目錄
	workDirs := []string{"work", "personal", "archive"}
	for _, dir := range workDirs {
		err := suite.fileManager.CreateDirectory(dir)
		if err != nil {
			t.Errorf("建立工作目錄 %s 失敗: %v", dir, err)
		}
	}
	
	// 步驟 2: 模擬用戶工作流程
	t.Log("步驟 2: 模擬用戶工作流程")
	
	// 場景 A: 建立工作筆記
	workNote, err := suite.editorService.CreateNote("每日工作計劃", "# 今日任務\n\n- [ ] 完成專案文檔\n- [ ] 程式碼審查\n- [ ] 會議準備")
	if err != nil {
		t.Fatalf("建立工作筆記失敗: %v", err)
	}
	
	workNote.FilePath = "work/daily_plan.md"
	err = suite.editorService.SaveNote(workNote)
	if err != nil {
		t.Errorf("保存工作筆記失敗: %v", err)
	}
	
	// 啟動自動保存
	suite.autoSaveService.StartAutoSave(workNote, 1*time.Second)
	
	// 場景 B: 建立個人筆記並加密
	personalNote, err := suite.editorService.CreateNote("個人日記", "今天的心情很好，完成了很多工作。")
	if err != nil {
		t.Fatalf("建立個人筆記失敗: %v", err)
	}
	
	personalNote.FilePath = "personal/diary.md.enc"
	
	// 為個人筆記啟用加密
	if es, ok := suite.editorService.(*editorService); ok {
		password := "MySecretPassword123!"
		err = es.EnableEncryption(personalNote.ID, password, "aes256", false)
		if err != nil {
			t.Errorf("啟用加密失敗: %v", err)
		} else {
			err = es.EncryptWithPassword(personalNote.ID, password)
			if err != nil {
				t.Errorf("加密筆記失敗: %v", err)
			}
		}
	}
	
	err = suite.editorService.SaveNote(personalNote)
	if err != nil {
		t.Errorf("保存個人筆記失敗: %v", err)
	}
	
	// 場景 C: 編輯和更新筆記
	t.Log("場景 C: 編輯和更新筆記")
	updatedContent := "# 今日任務\n\n- [x] 完成專案文檔\n- [ ] 程式碼審查\n- [ ] 會議準備\n\n## 新增任務\n- [ ] 回覆郵件"
	err = suite.editorService.UpdateContent(workNote.ID, updatedContent)
	if err != nil {
		t.Errorf("更新工作筆記失敗: %v", err)
	}
	
	// 等待自動保存
	time.Sleep(1500 * time.Millisecond)
	
	// 場景 D: 檔案管理操作
	t.Log("場景 D: 檔案管理操作")
	
	// 建立備份目錄
	err = suite.fileManager.CreateDirectory("backup")
	if err != nil {
		t.Errorf("建立備份目錄失敗: %v", err)
	}
	
	// 複製重要檔案到備份目錄
	err = suite.fileManager.CopyFile("work/daily_plan.md", "backup/daily_plan_backup.md")
	if err != nil {
		t.Errorf("備份檔案失敗: %v", err)
	}
	
	// 場景 E: 搜尋和瀏覽
	t.Log("場景 E: 搜尋和瀏覽")
	
	// 搜尋所有 Markdown 檔案
	allFiles, err := suite.fileManager.SearchFiles(".", "*.md*", true)
	if err != nil {
		t.Errorf("搜尋檔案失敗: %v", err)
	} else {
		t.Logf("找到 %d 個檔案", len(allFiles))
	}
	
	// 列出所有目錄
	rootFiles, err := suite.fileManager.ListFiles(".")
	if err != nil {
		t.Errorf("列出根目錄失敗: %v", err)
	} else {
		t.Logf("根目錄包含 %d 個項目", len(rootFiles))
	}
	
	// 步驟 3: 效能檢查
	t.Log("步驟 3: 效能檢查")
	
	metrics := suite.perfService.GetCurrentMetrics()
	if metrics != nil {
		t.Logf("工作流程完成後的效能指標:")
		t.Logf("- 記憶體使用: %d bytes", metrics.MemoryUsage)
		t.Logf("- Goroutines: %d", metrics.Goroutines)
		t.Logf("- 垃圾回收: %d 次", metrics.GCCount)
	}
	
	// 步驟 4: 清理和關閉
	t.Log("步驟 4: 清理和關閉")
	
	// 停止自動保存
	suite.autoSaveService.StopAutoSave(workNote.ID)
	
	// 關閉所有筆記
	suite.editorService.CloseNote(workNote.ID)
	suite.editorService.CloseNote(personalNote.ID)
	
	// 執行最終的記憶體優化
	err = suite.perfService.OptimizeMemory()
	if err != nil {
		t.Errorf("最終記憶體優化失敗: %v", err)
	}
	
	// 驗證清理效果
	finalMetrics := suite.perfService.GetCurrentMetrics()
	if finalMetrics != nil {
		t.Logf("清理後的記憶體使用: %d bytes", finalMetrics.MemoryUsage)
	}
	
	t.Log("完整工作流程測試完成")
}