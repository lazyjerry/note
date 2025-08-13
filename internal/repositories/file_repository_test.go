package repositories

import (
	"os"        // 作業系統介面套件
	"path/filepath" // 檔案路徑處理套件
	"testing"   // Go 測試套件
	"time"      // 時間處理套件
	
	"mac-notebook-app/internal/models" // 引入資料模型
)

// TestNewLocalFileRepository 測試 LocalFileRepository 的建立
// 驗證建構函數是否正確處理各種輸入情況
func TestNewLocalFileRepository(t *testing.T) {
	// 測試案例：成功建立儲存庫
	t.Run("成功建立儲存庫", func(t *testing.T) {
		// 建立臨時目錄用於測試
		tempDir := t.TempDir()
		
		// 建立儲存庫實例
		repo, err := NewLocalFileRepository(tempDir)
		
		// 驗證結果
		if err != nil {
			t.Fatalf("建立儲存庫時發生錯誤：%v", err)
		}
		
		if repo == nil {
			t.Fatal("儲存庫實例不應為 nil")
		}
		
		if repo.GetBaseDirectory() != tempDir {
			t.Errorf("基礎目錄不符合預期，期望：%s，實際：%s", tempDir, repo.GetBaseDirectory())
		}
	})
	
	// 測試案例：空路徑應該回傳錯誤
	t.Run("空路徑應該回傳錯誤", func(t *testing.T) {
		repo, err := NewLocalFileRepository("")
		
		if err == nil {
			t.Fatal("空路徑應該回傳錯誤")
		}
		
		if repo != nil {
			t.Fatal("錯誤情況下儲存庫實例應為 nil")
		}
		
		// 檢查錯誤類型
		if appErr, ok := err.(*models.AppError); ok {
			if appErr.Code != models.ErrValidationFailed {
				t.Errorf("錯誤代碼不符合預期，期望：%s，實際：%s", models.ErrValidationFailed, appErr.Code)
			}
		} else {
			t.Error("應該回傳 AppError 類型的錯誤")
		}
	})
	
	// 測試案例：自動建立不存在的目錄
	t.Run("自動建立不存在的目錄", func(t *testing.T) {
		// 建立臨時目錄
		tempDir := t.TempDir()
		
		// 建立不存在的子目錄路徑
		newDir := filepath.Join(tempDir, "new_directory")
		
		// 建立儲存庫實例
		repo, err := NewLocalFileRepository(newDir)
		
		// 驗證結果
		if err != nil {
			t.Fatalf("建立儲存庫時發生錯誤：%v", err)
		}
		
		// 檢查目錄是否被建立
		if _, err := os.Stat(newDir); os.IsNotExist(err) {
			t.Error("目錄應該被自動建立")
		}
		
		if repo.GetBaseDirectory() != newDir {
			t.Errorf("基礎目錄不符合預期，期望：%s，實際：%s", newDir, repo.GetBaseDirectory())
		}
	})
}

// TestFileOperations 測試基本檔案操作功能
func TestFileOperations(t *testing.T) {
	// 建立測試用的儲存庫
	tempDir := t.TempDir()
	repo, err := NewLocalFileRepository(tempDir)
	if err != nil {
		t.Fatalf("建立測試儲存庫失敗：%v", err)
	}
	
	testFilePath := "test.md"
	testContent := "# 測試標題\n\n這是測試內容。"
	
	// 測試案例：寫入檔案
	t.Run("寫入檔案", func(t *testing.T) {
		err := repo.WriteFile(testFilePath, []byte(testContent))
		if err != nil {
			t.Fatalf("寫入檔案失敗：%v", err)
		}
		
		// 檢查檔案是否存在
		if !repo.FileExists(testFilePath) {
			t.Error("檔案應該存在")
		}
	})
	
	// 測試案例：讀取檔案
	t.Run("讀取檔案", func(t *testing.T) {
		data, err := repo.ReadFile(testFilePath)
		if err != nil {
			t.Fatalf("讀取檔案失敗：%v", err)
		}
		
		if string(data) != testContent {
			t.Errorf("檔案內容不符合預期，期望：%s，實際：%s", testContent, string(data))
		}
	})
	
	// 測試案例：檔案存在性檢查
	t.Run("檔案存在性檢查", func(t *testing.T) {
		// 檢查存在的檔案
		if !repo.FileExists(testFilePath) {
			t.Error("檔案應該存在")
		}
		
		// 檢查不存在的檔案
		if repo.FileExists("nonexistent.md") {
			t.Error("不存在的檔案應該回傳 false")
		}
	})
	
	// 測試案例：刪除檔案
	t.Run("刪除檔案", func(t *testing.T) {
		err := repo.DeleteFile(testFilePath)
		if err != nil {
			t.Fatalf("刪除檔案失敗：%v", err)
		}
		
		// 檢查檔案是否被刪除
		if repo.FileExists(testFilePath) {
			t.Error("檔案應該被刪除")
		}
	})
	
	// 測試案例：讀取不存在的檔案應該回傳錯誤
	t.Run("讀取不存在的檔案應該回傳錯誤", func(t *testing.T) {
		_, err := repo.ReadFile("nonexistent.md")
		if err == nil {
			t.Fatal("讀取不存在的檔案應該回傳錯誤")
		}
		
		// 檢查錯誤類型
		if appErr, ok := err.(*models.AppError); ok {
			if appErr.Code != models.ErrFileNotFound {
				t.Errorf("錯誤代碼不符合預期，期望：%s，實際：%s", models.ErrFileNotFound, appErr.Code)
			}
		}
	})
	
	// 測試案例：刪除不存在的檔案應該回傳錯誤
	t.Run("刪除不存在的檔案應該回傳錯誤", func(t *testing.T) {
		err := repo.DeleteFile("nonexistent.md")
		if err == nil {
			t.Fatal("刪除不存在的檔案應該回傳錯誤")
		}
		
		// 檢查錯誤類型
		if appErr, ok := err.(*models.AppError); ok {
			if appErr.Code != models.ErrFileNotFound {
				t.Errorf("錯誤代碼不符合預期，期望：%s，實際：%s", models.ErrFileNotFound, appErr.Code)
			}
		}
	})
}

// TestDirectoryOperations 測試目錄操作功能
func TestDirectoryOperations(t *testing.T) {
	// 建立測試用的儲存庫
	tempDir := t.TempDir()
	repo, err := NewLocalFileRepository(tempDir)
	if err != nil {
		t.Fatalf("建立測試儲存庫失敗：%v", err)
	}
	
	// 測試案例：建立目錄
	t.Run("建立目錄", func(t *testing.T) {
		testDirPath := "test_directory"
		
		err := repo.CreateDirectory(testDirPath)
		if err != nil {
			t.Fatalf("建立目錄失敗：%v", err)
		}
		
		// 檢查目錄是否存在
		fullPath := filepath.Join(tempDir, testDirPath)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			t.Error("目錄應該存在")
		}
	})
	
	// 測試案例：建立巢狀目錄
	t.Run("建立巢狀目錄", func(t *testing.T) {
		nestedPath := "parent/child/grandchild"
		
		err := repo.CreateDirectory(nestedPath)
		if err != nil {
			t.Fatalf("建立巢狀目錄失敗：%v", err)
		}
		
		// 檢查目錄是否存在
		fullPath := filepath.Join(tempDir, nestedPath)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			t.Error("巢狀目錄應該存在")
		}
	})
	
	// 測試案例：列出目錄內容
	t.Run("列出目錄內容", func(t *testing.T) {
		// 建立測試檔案和目錄
		testDir := "list_test"
		repo.CreateDirectory(testDir)
		
		// 建立測試檔案
		testFiles := []string{
			filepath.Join(testDir, "file1.md"),
			filepath.Join(testDir, "file2.txt"),
			filepath.Join(testDir, "encrypted.md.enc"),
		}
		
		for _, filePath := range testFiles {
			repo.WriteFile(filePath, []byte("測試內容"))
		}
		
		// 建立子目錄
		subDir := filepath.Join(testDir, "subdirectory")
		repo.CreateDirectory(subDir)
		
		// 列出目錄內容
		fileInfos, err := repo.ListDirectory(testDir)
		if err != nil {
			t.Fatalf("列出目錄內容失敗：%v", err)
		}
		
		// 驗證結果
		if len(fileInfos) != 4 { // 3 個檔案 + 1 個子目錄
			t.Errorf("目錄項目數量不符合預期，期望：4，實際：%d", len(fileInfos))
		}
		
		// 檢查是否包含預期的項目
		foundItems := make(map[string]bool)
		for _, info := range fileInfos {
			foundItems[info.Name] = true
			
			// 檢查 Markdown 檔案識別
			if info.Name == "file1.md" && !info.IsMarkdownFile() {
				t.Error("file1.md 應該被識別為 Markdown 檔案")
			}
			
			if info.Name == "encrypted.md.enc" && !info.IsMarkdownFile() {
				t.Error("encrypted.md.enc 應該被識別為 Markdown 檔案")
			}
			
			if info.Name == "file2.txt" && info.IsMarkdownFile() {
				t.Error("file2.txt 不應該被識別為 Markdown 檔案")
			}
			
			if info.Name == "subdirectory" && !info.IsDirectory {
				t.Error("subdirectory 應該被識別為目錄")
			}
		}
		
		expectedItems := []string{"file1.md", "file2.txt", "encrypted.md.enc", "subdirectory"}
		for _, item := range expectedItems {
			if !foundItems[item] {
				t.Errorf("找不到預期的項目：%s", item)
			}
		}
	})
	
	// 測試案例：列出不存在的目錄應該回傳錯誤
	t.Run("列出不存在的目錄應該回傳錯誤", func(t *testing.T) {
		_, err := repo.ListDirectory("nonexistent_directory")
		if err == nil {
			t.Fatal("列出不存在的目錄應該回傳錯誤")
		}
		
		// 檢查錯誤類型
		if appErr, ok := err.(*models.AppError); ok {
			if appErr.Code != models.ErrFileNotFound {
				t.Errorf("錯誤代碼不符合預期，期望：%s，實際：%s", models.ErrFileNotFound, appErr.Code)
			}
		}
	})
}// 
TestMarkdownOperations 測試 Markdown 檔案特殊操作
func TestMarkdownOperations(t *testing.T) {
	// 建立測試用的儲存庫
	tempDir := t.TempDir()
	repo, err := NewLocalFileRepository(tempDir)
	if err != nil {
		t.Fatalf("建立測試儲存庫失敗：%v", err)
	}
	
	// 測試案例：Markdown 檔案識別
	t.Run("Markdown 檔案識別", func(t *testing.T) {
		testCases := []struct {
			fileName   string
			isMarkdown bool
		}{
			{"test.md", true},
			{"encrypted.md.enc", true},
			{"document.txt", false},
			{"readme.MD", false}, // 大小寫敏感
			{"note.markdown", false}, // 不支援 .markdown 副檔名
		}
		
		// 建立測試檔案
		for _, tc := range testCases {
			repo.WriteFile(tc.fileName, []byte("測試內容"))
		}
		
		// 測試檔案識別
		for _, tc := range testCases {
			isMarkdown, err := repo.IsMarkdownFile(tc.fileName)
			if err != nil {
				t.Fatalf("檢查 Markdown 檔案時發生錯誤：%v", err)
			}
			
			if isMarkdown != tc.isMarkdown {
				t.Errorf("檔案 %s 的 Markdown 識別結果不符合預期，期望：%v，實際：%v", 
					tc.fileName, tc.isMarkdown, isMarkdown)
			}
		}
	})
	
	// 測試案例：讀取 Markdown 檔案
	t.Run("讀取 Markdown 檔案", func(t *testing.T) {
		testContent := "# 測試標題\n\n這是 **粗體** 文字。"
		testFile := "markdown_test.md"
		
		// 寫入 Markdown 檔案
		err := repo.WriteMarkdownFile(testFile, testContent)
		if err != nil {
			t.Fatalf("寫入 Markdown 檔案失敗：%v", err)
		}
		
		// 讀取 Markdown 檔案
		content, err := repo.ReadMarkdownFile(testFile)
		if err != nil {
			t.Fatalf("讀取 Markdown 檔案失敗：%v", err)
		}
		
		if content != testContent {
			t.Errorf("Markdown 檔案內容不符合預期，期望：%s，實際：%s", testContent, content)
		}
	})
	
	// 測試案例：寫入 Markdown 檔案
	t.Run("寫入 Markdown 檔案", func(t *testing.T) {
		testContent := "## 子標題\n\n- 項目 1\n- 項目 2"
		testFile := "new_markdown.md"
		
		err := repo.WriteMarkdownFile(testFile, testContent)
		if err != nil {
			t.Fatalf("寫入 Markdown 檔案失敗：%v", err)
		}
		
		// 驗證檔案是否存在
		if !repo.FileExists(testFile) {
			t.Error("Markdown 檔案應該存在")
		}
		
		// 驗證內容
		content, err := repo.ReadMarkdownFile(testFile)
		if err != nil {
			t.Fatalf("讀取 Markdown 檔案失敗：%v", err)
		}
		
		if content != testContent {
			t.Errorf("Markdown 檔案內容不符合預期，期望：%s，實際：%s", testContent, content)
		}
	})
	
	// 測試案例：寫入沒有正確副檔名的檔案應該回傳錯誤
	t.Run("寫入沒有正確副檔名的檔案應該回傳錯誤", func(t *testing.T) {
		err := repo.WriteMarkdownFile("invalid.txt", "內容")
		if err == nil {
			t.Fatal("寫入沒有正確副檔名的檔案應該回傳錯誤")
		}
		
		// 檢查錯誤類型
		if appErr, ok := err.(*models.AppError); ok {
			if appErr.Code != models.ErrValidationFailed {
				t.Errorf("錯誤代碼不符合預期，期望：%s，實際：%s", models.ErrValidationFailed, appErr.Code)
			}
		}
	})
	
	// 測試案例：讀取非 Markdown 檔案應該回傳錯誤
	t.Run("讀取非 Markdown 檔案應該回傳錯誤", func(t *testing.T) {
		// 建立非 Markdown 檔案
		textFile := "document.txt"
		repo.WriteFile(textFile, []byte("這不是 Markdown 檔案"))
		
		_, err := repo.ReadMarkdownFile(textFile)
		if err == nil {
			t.Fatal("讀取非 Markdown 檔案應該回傳錯誤")
		}
		
		// 檢查錯誤類型
		if appErr, ok := err.(*models.AppError); ok {
			if appErr.Code != models.ErrValidationFailed {
				t.Errorf("錯誤代碼不符合預期，期望：%s，實際：%s", models.ErrValidationFailed, appErr.Code)
			}
		}
	})
}

// TestPathValidation 測試路徑驗證功能
func TestPathValidation(t *testing.T) {
	// 建立測試用的儲存庫
	tempDir := t.TempDir()
	repo, err := NewLocalFileRepository(tempDir)
	if err != nil {
		t.Fatalf("建立測試儲存庫失敗：%v", err)
	}
	
	// 測試案例：有效路徑
	t.Run("有效路徑", func(t *testing.T) {
		validPaths := []string{
			"file.md",
			"folder/file.md",
			"deep/nested/path/file.md",
			"file with spaces.md",
			"中文檔案.md",
		}
		
		for _, path := range validPaths {
			err := repo.WriteFile(path, []byte("測試內容"))
			if err != nil {
				t.Errorf("有效路徑 %s 應該被接受，但發生錯誤：%v", path, err)
			}
		}
	})
	
	// 測試案例：無效路徑
	t.Run("無效路徑", func(t *testing.T) {
		invalidPaths := []string{
			"",                    // 空路徑
			"../outside.md",       // 包含 ..
			"/absolute/path.md",   // 絕對路徑
			"folder/../outside.md", // 目錄遍歷攻擊
		}
		
		for _, path := range invalidPaths {
			err := repo.WriteFile(path, []byte("測試內容"))
			if err == nil {
				t.Errorf("無效路徑 %s 應該被拒絕", path)
			}
			
			// 檢查錯誤類型
			if appErr, ok := err.(*models.AppError); ok {
				if appErr.Code != models.ErrValidationFailed {
					t.Errorf("路徑 %s 的錯誤代碼不符合預期，期望：%s，實際：%s", 
						path, models.ErrValidationFailed, appErr.Code)
				}
			}
		}
	})
}

// TestWalkDirectory 測試目錄遍歷功能
func TestWalkDirectory(t *testing.T) {
	// 建立測試用的儲存庫
	tempDir := t.TempDir()
	repo, err := NewLocalFileRepository(tempDir)
	if err != nil {
		t.Fatalf("建立測試儲存庫失敗：%v", err)
	}
	
	// 建立測試目錄結構
	testStructure := []string{
		"root.md",
		"folder1/file1.md",
		"folder1/file2.txt",
		"folder1/subfolder/deep.md",
		"folder2/encrypted.md.enc",
	}
	
	for _, path := range testStructure {
		if err := repo.WriteFile(path, []byte("測試內容")); err != nil {
			t.Fatalf("建立測試檔案 %s 失敗：%v", path, err)
		}
	}
	
	// 測試案例：遍歷整個目錄樹
	t.Run("遍歷整個目錄樹", func(t *testing.T) {
		var foundFiles []string
		
		err := repo.WalkDirectory(".", func(info *models.FileInfo) error {
			if !info.IsDirectory {
				foundFiles = append(foundFiles, info.Path)
			}
			return nil
		})
		
		if err != nil {
			t.Fatalf("遍歷目錄失敗：%v", err)
		}
		
		// 檢查是否找到所有檔案
		if len(foundFiles) != len(testStructure) {
			t.Errorf("找到的檔案數量不符合預期，期望：%d，實際：%d", len(testStructure), len(foundFiles))
		}
		
		// 檢查每個檔案是否都被找到
		foundMap := make(map[string]bool)
		for _, file := range foundFiles {
			foundMap[file] = true
		}
		
		for _, expectedFile := range testStructure {
			if !foundMap[expectedFile] {
				t.Errorf("找不到預期的檔案：%s", expectedFile)
			}
		}
	})
	
	// 測試案例：遍歷特定子目錄
	t.Run("遍歷特定子目錄", func(t *testing.T) {
		var foundFiles []string
		
		err := repo.WalkDirectory("folder1", func(info *models.FileInfo) error {
			if !info.IsDirectory {
				foundFiles = append(foundFiles, info.Name)
			}
			return nil
		})
		
		if err != nil {
			t.Fatalf("遍歷子目錄失敗：%v", err)
		}
		
		expectedFiles := []string{"file1.md", "file2.txt", "deep.md"}
		if len(foundFiles) != len(expectedFiles) {
			t.Errorf("找到的檔案數量不符合預期，期望：%d，實際：%d", len(expectedFiles), len(foundFiles))
		}
	})
	
	// 測試案例：遍歷過程中發生錯誤
	t.Run("遍歷過程中發生錯誤", func(t *testing.T) {
		expectedError := "測試錯誤"
		
		err := repo.WalkDirectory(".", func(info *models.FileInfo) error {
			if info.Name == "file1.md" {
				return fmt.Errorf(expectedError)
			}
			return nil
		})
		
		if err == nil {
			t.Fatal("應該回傳錯誤")
		}
		
		if err.Error() != expectedError {
			t.Errorf("錯誤訊息不符合預期，期望：%s，實際：%s", expectedError, err.Error())
		}
	})
}

// TestConcurrentAccess 測試並發存取
func TestConcurrentAccess(t *testing.T) {
	// 建立測試用的儲存庫
	tempDir := t.TempDir()
	repo, err := NewLocalFileRepository(tempDir)
	if err != nil {
		t.Fatalf("建立測試儲存庫失敗：%v", err)
	}
	
	// 測試案例：並發寫入不同檔案
	t.Run("並發寫入不同檔案", func(t *testing.T) {
		const numGoroutines = 10
		done := make(chan bool, numGoroutines)
		
		// 啟動多個 goroutine 同時寫入不同檔案
		for i := 0; i < numGoroutines; i++ {
			go func(index int) {
				defer func() { done <- true }()
				
				fileName := fmt.Sprintf("concurrent_%d.md", index)
				content := fmt.Sprintf("# 檔案 %d\n\n這是第 %d 個檔案的內容。", index, index)
				
				if err := repo.WriteFile(fileName, []byte(content)); err != nil {
					t.Errorf("並發寫入檔案 %s 失敗：%v", fileName, err)
				}
			}(i)
		}
		
		// 等待所有 goroutine 完成
		for i := 0; i < numGoroutines; i++ {
			<-done
		}
		
		// 驗證所有檔案都被正確建立
		for i := 0; i < numGoroutines; i++ {
			fileName := fmt.Sprintf("concurrent_%d.md", i)
			if !repo.FileExists(fileName) {
				t.Errorf("並發建立的檔案 %s 不存在", fileName)
			}
		}
	})
}

// BenchmarkFileOperations 效能測試
func BenchmarkFileOperations(b *testing.B) {
	// 建立測試用的儲存庫
	tempDir := b.TempDir()
	repo, err := NewLocalFileRepository(tempDir)
	if err != nil {
		b.Fatalf("建立測試儲存庫失敗：%v", err)
	}
	
	testContent := []byte("# 效能測試\n\n這是用於效能測試的內容。")
	
	// 測試寫入效能
	b.Run("WriteFile", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			fileName := fmt.Sprintf("bench_%d.md", i)
			if err := repo.WriteFile(fileName, testContent); err != nil {
				b.Fatalf("寫入檔案失敗：%v", err)
			}
		}
	})
	
	// 測試讀取效能
	b.Run("ReadFile", func(b *testing.B) {
		// 先建立測試檔案
		testFile := "bench_read.md"
		repo.WriteFile(testFile, testContent)
		
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if _, err := repo.ReadFile(testFile); err != nil {
				b.Fatalf("讀取檔案失敗：%v", err)
			}
		}
	})
	
	// 測試檔案存在性檢查效能
	b.Run("FileExists", func(b *testing.B) {
		testFile := "bench_exists.md"
		repo.WriteFile(testFile, testContent)
		
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			repo.FileExists(testFile)
		}
	})
}