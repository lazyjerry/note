package services

import (
	"os"        // 作業系統介面套件
	"path/filepath" // 檔案路徑處理套件
	"testing"   // Go 測試套件
	
	"mac-notebook-app/internal/models"      // 引入資料模型
	"mac-notebook-app/internal/repositories" // 引入儲存庫介面
)

// TestNewLocalFileManagerService 測試 LocalFileManagerService 的建立
func TestNewLocalFileManagerService(t *testing.T) {
	// 測試案例：成功建立服務
	t.Run("成功建立服務", func(t *testing.T) {
		// 建立臨時目錄和檔案儲存庫
		tempDir := t.TempDir()
		fileRepo, err := repositories.NewLocalFileRepository(tempDir)
		if err != nil {
			t.Fatalf("建立檔案儲存庫失敗：%v", err)
		}
		
		// 建立檔案管理服務
		service, err := NewLocalFileManagerService(fileRepo, tempDir)
		
		// 驗證結果
		if err != nil {
			t.Fatalf("建立檔案管理服務時發生錯誤：%v", err)
		}
		
		if service == nil {
			t.Fatal("服務實例不應為 nil")
		}
		
		if service.baseDir != tempDir {
			t.Errorf("基礎目錄不符合預期，期望：%s，實際：%s", tempDir, service.baseDir)
		}
	})
	
	// 測試案例：nil 檔案儲存庫應該回傳錯誤
	t.Run("nil 檔案儲存庫應該回傳錯誤", func(t *testing.T) {
		tempDir := t.TempDir()
		
		service, err := NewLocalFileManagerService(nil, tempDir)
		
		if err == nil {
			t.Fatal("nil 檔案儲存庫應該回傳錯誤")
		}
		
		if service != nil {
			t.Fatal("錯誤情況下服務實例應為 nil")
		}
		
		// 檢查錯誤類型
		if appErr, ok := err.(*models.AppError); ok {
			if appErr.Code != models.ErrValidationFailed {
				t.Errorf("錯誤代碼不符合預期，期望：%s，實際：%s", models.ErrValidationFailed, appErr.Code)
			}
		}
	})
	
	// 測試案例：空基礎目錄應該回傳錯誤
	t.Run("空基礎目錄應該回傳錯誤", func(t *testing.T) {
		tempDir := t.TempDir()
		fileRepo, _ := repositories.NewLocalFileRepository(tempDir)
		
		service, err := NewLocalFileManagerService(fileRepo, "")
		
		if err == nil {
			t.Fatal("空基礎目錄應該回傳錯誤")
		}
		
		if service != nil {
			t.Fatal("錯誤情況下服務實例應為 nil")
		}
	})
	
	// 測試案例：不存在的基礎目錄應該回傳錯誤
	t.Run("不存在的基礎目錄應該回傳錯誤", func(t *testing.T) {
		tempDir := t.TempDir()
		fileRepo, _ := repositories.NewLocalFileRepository(tempDir)
		
		nonExistentDir := filepath.Join(tempDir, "nonexistent")
		
		service, err := NewLocalFileManagerService(fileRepo, nonExistentDir)
		
		if err == nil {
			t.Fatal("不存在的基礎目錄應該回傳錯誤")
		}
		
		if service != nil {
			t.Fatal("錯誤情況下服務實例應為 nil")
		}
		
		// 檢查錯誤類型
		if appErr, ok := err.(*models.AppError); ok {
			if appErr.Code != models.ErrFileNotFound {
				t.Errorf("錯誤代碼不符合預期，期望：%s，實際：%s", models.ErrFileNotFound, appErr.Code)
			}
		}
	})
}

// TestListFiles 測試檔案列表功能
func TestListFiles(t *testing.T) {
	// 建立測試環境
	tempDir := t.TempDir()
	fileRepo, _ := repositories.NewLocalFileRepository(tempDir)
	service, _ := NewLocalFileManagerService(fileRepo, tempDir)
	
	// 建立測試檔案和目錄
	testFiles := []string{
		"file1.md",
		"file2.txt",
		"encrypted.md.enc",
	}
	
	testDirs := []string{
		"folder1",
		"folder2",
	}
	
	// 建立測試檔案
	for _, file := range testFiles {
		fileRepo.WriteFile(file, []byte("測試內容"))
	}
	
	// 建立測試目錄
	for _, dir := range testDirs {
		fileRepo.CreateDirectory(dir)
	}
	
	// 測試案例：列出根目錄內容
	t.Run("列出根目錄內容", func(t *testing.T) {
		fileInfos, err := service.ListFiles(".")
		if err != nil {
			t.Fatalf("列出檔案失敗：%v", err)
		}
		
		// 驗證檔案數量
		expectedCount := len(testFiles) + len(testDirs)
		if len(fileInfos) != expectedCount {
			t.Errorf("檔案數量不符合預期，期望：%d，實際：%d", expectedCount, len(fileInfos))
		}
		
		// 驗證排序（目錄應該在前面）
		dirCount := 0
		for i, info := range fileInfos {
			if info.IsDirectory {
				dirCount++
			} else {
				// 檔案應該在所有目錄之後
				if i < len(testDirs) {
					t.Error("檔案排序不正確，目錄應該在檔案之前")
				}
			}
		}
		
		if dirCount != len(testDirs) {
			t.Errorf("目錄數量不符合預期，期望：%d，實際：%d", len(testDirs), dirCount)
		}
	})
	
	// 測試案例：列出不存在的目錄應該回傳錯誤
	t.Run("列出不存在的目錄應該回傳錯誤", func(t *testing.T) {
		_, err := service.ListFiles("nonexistent")
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
}

// TestCreateDirectory 測試目錄建立功能
func TestCreateDirectory(t *testing.T) {
	// 建立測試環境
	tempDir := t.TempDir()
	fileRepo, _ := repositories.NewLocalFileRepository(tempDir)
	service, _ := NewLocalFileManagerService(fileRepo, tempDir)
	
	// 測試案例：建立新目錄
	t.Run("建立新目錄", func(t *testing.T) {
		testDir := "new_directory"
		
		err := service.CreateDirectory(testDir)
		if err != nil {
			t.Fatalf("建立目錄失敗：%v", err)
		}
		
		// 驗證目錄是否存在
		if !fileRepo.FileExists(testDir) {
			t.Error("目錄應該存在")
		}
		
		// 驗證是否為目錄
		fullPath := filepath.Join(tempDir, testDir)
		info, err := os.Stat(fullPath)
		if err != nil {
			t.Fatalf("無法取得目錄資訊：%v", err)
		}
		
		if !info.IsDir() {
			t.Error("建立的項目應該是目錄")
		}
	})
	
	// 測試案例：建立巢狀目錄
	t.Run("建立巢狀目錄", func(t *testing.T) {
		nestedDir := "parent/child/grandchild"
		
		err := service.CreateDirectory(nestedDir)
		if err != nil {
			t.Fatalf("建立巢狀目錄失敗：%v", err)
		}
		
		// 驗證目錄是否存在
		if !fileRepo.FileExists(nestedDir) {
			t.Error("巢狀目錄應該存在")
		}
	})
	
	// 測試案例：建立已存在的目錄應該回傳錯誤
	t.Run("建立已存在的目錄應該回傳錯誤", func(t *testing.T) {
		existingDir := "existing_directory"
		
		// 先建立目錄
		service.CreateDirectory(existingDir)
		
		// 嘗試再次建立相同目錄
		err := service.CreateDirectory(existingDir)
		if err == nil {
			t.Fatal("建立已存在的目錄應該回傳錯誤")
		}
		
		// 檢查錯誤類型
		if appErr, ok := err.(*models.AppError); ok {
			if appErr.Code != models.ErrValidationFailed {
				t.Errorf("錯誤代碼不符合預期，期望：%s，實際：%s", models.ErrValidationFailed, appErr.Code)
			}
		}
	})
}

// TestDeleteFile 測試檔案刪除功能
func TestDeleteFile(t *testing.T) {
	// 建立測試環境
	tempDir := t.TempDir()
	fileRepo, _ := repositories.NewLocalFileRepository(tempDir)
	service, _ := NewLocalFileManagerService(fileRepo, tempDir)
	
	// 測試案例：刪除檔案
	t.Run("刪除檔案", func(t *testing.T) {
		testFile := "test_file.md"
		
		// 建立測試檔案
		fileRepo.WriteFile(testFile, []byte("測試內容"))
		
		// 刪除檔案
		err := service.DeleteFile(testFile)
		if err != nil {
			t.Fatalf("刪除檔案失敗：%v", err)
		}
		
		// 驗證檔案是否被刪除
		if fileRepo.FileExists(testFile) {
			t.Error("檔案應該被刪除")
		}
	})
	
	// 測試案例：刪除空目錄
	t.Run("刪除空目錄", func(t *testing.T) {
		testDir := "empty_directory"
		
		// 建立空目錄
		service.CreateDirectory(testDir)
		
		// 刪除目錄
		err := service.DeleteFile(testDir)
		if err != nil {
			t.Fatalf("刪除空目錄失敗：%v", err)
		}
		
		// 驗證目錄是否被刪除
		if fileRepo.FileExists(testDir) {
			t.Error("目錄應該被刪除")
		}
	})
	
	// 測試案例：刪除非空目錄應該回傳錯誤
	t.Run("刪除非空目錄應該回傳錯誤", func(t *testing.T) {
		testDir := "non_empty_directory"
		testFile := filepath.Join(testDir, "file.txt")
		
		// 建立目錄和檔案
		service.CreateDirectory(testDir)
		fileRepo.WriteFile(testFile, []byte("內容"))
		
		// 嘗試刪除非空目錄
		err := service.DeleteFile(testDir)
		if err == nil {
			t.Fatal("刪除非空目錄應該回傳錯誤")
		}
		
		// 檢查錯誤類型
		if appErr, ok := err.(*models.AppError); ok {
			if appErr.Code != models.ErrValidationFailed {
				t.Errorf("錯誤代碼不符合預期，期望：%s，實際：%s", models.ErrValidationFailed, appErr.Code)
			}
		}
	})
	
	// 測試案例：刪除不存在的檔案應該回傳錯誤
	t.Run("刪除不存在的檔案應該回傳錯誤", func(t *testing.T) {
		err := service.DeleteFile("nonexistent.txt")
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

// TestRenameFile 測試檔案重新命名功能
func TestRenameFile(t *testing.T) {
	// 建立測試環境
	tempDir := t.TempDir()
	fileRepo, _ := repositories.NewLocalFileRepository(tempDir)
	service, _ := NewLocalFileManagerService(fileRepo, tempDir)
	
	// 測試案例：重新命名檔案
	t.Run("重新命名檔案", func(t *testing.T) {
		oldName := "old_file.md"
		newName := "new_file.md"
		testContent := "測試內容"
		
		// 建立測試檔案
		fileRepo.WriteFile(oldName, []byte(testContent))
		
		// 重新命名檔案
		err := service.RenameFile(oldName, newName)
		if err != nil {
			t.Fatalf("重新命名檔案失敗：%v", err)
		}
		
		// 驗證舊檔案不存在
		if fileRepo.FileExists(oldName) {
			t.Error("舊檔案應該不存在")
		}
		
		// 驗證新檔案存在
		if !fileRepo.FileExists(newName) {
			t.Error("新檔案應該存在")
		}
		
		// 驗證檔案內容
		data, err := fileRepo.ReadFile(newName)
		if err != nil {
			t.Fatalf("讀取重新命名的檔案失敗：%v", err)
		}
		
		if string(data) != testContent {
			t.Errorf("檔案內容不符合預期，期望：%s，實際：%s", testContent, string(data))
		}
	})
	
	// 測試案例：重新命名目錄
	t.Run("重新命名目錄", func(t *testing.T) {
		oldDirName := "old_directory"
		newDirName := "new_directory"
		
		// 建立測試目錄
		service.CreateDirectory(oldDirName)
		
		// 重新命名目錄
		err := service.RenameFile(oldDirName, newDirName)
		if err != nil {
			t.Fatalf("重新命名目錄失敗：%v", err)
		}
		
		// 驗證舊目錄不存在
		if fileRepo.FileExists(oldDirName) {
			t.Error("舊目錄應該不存在")
		}
		
		// 驗證新目錄存在
		if !fileRepo.FileExists(newDirName) {
			t.Error("新目錄應該存在")
		}
	})
	
	// 測試案例：重新命名到已存在的路徑應該回傳錯誤
	t.Run("重新命名到已存在的路徑應該回傳錯誤", func(t *testing.T) {
		file1 := "file1.txt"
		file2 := "file2.txt"
		
		// 建立兩個檔案
		fileRepo.WriteFile(file1, []byte("內容1"))
		fileRepo.WriteFile(file2, []byte("內容2"))
		
		// 嘗試將 file1 重新命名為 file2
		err := service.RenameFile(file1, file2)
		if err == nil {
			t.Fatal("重新命名到已存在的路徑應該回傳錯誤")
		}
		
		// 檢查錯誤類型
		if appErr, ok := err.(*models.AppError); ok {
			if appErr.Code != models.ErrValidationFailed {
				t.Errorf("錯誤代碼不符合預期，期望：%s，實際：%s", models.ErrValidationFailed, appErr.Code)
			}
		}
	})
	
	// 測試案例：重新命名不存在的檔案應該回傳錯誤
	t.Run("重新命名不存在的檔案應該回傳錯誤", func(t *testing.T) {
		err := service.RenameFile("nonexistent.txt", "new_name.txt")
		if err == nil {
			t.Fatal("重新命名不存在的檔案應該回傳錯誤")
		}
		
		// 檢查錯誤類型
		if appErr, ok := err.(*models.AppError); ok {
			if appErr.Code != models.ErrFileNotFound {
				t.Errorf("錯誤代碼不符合預期，期望：%s，實際：%s", models.ErrFileNotFound, appErr.Code)
			}
		}
	})
}

// TestMoveFile 測試檔案移動功能
func TestMoveFile(t *testing.T) {
	// 建立測試環境
	tempDir := t.TempDir()
	fileRepo, _ := repositories.NewLocalFileRepository(tempDir)
	service, _ := NewLocalFileManagerService(fileRepo, tempDir)
	
	// 測試案例：移動檔案到目錄
	t.Run("移動檔案到目錄", func(t *testing.T) {
		sourceFile := "source.md"
		targetDir := "target_directory"
		testContent := "測試內容"
		
		// 建立來源檔案和目標目錄
		fileRepo.WriteFile(sourceFile, []byte(testContent))
		service.CreateDirectory(targetDir)
		
		// 移動檔案
		err := service.MoveFile(sourceFile, targetDir)
		if err != nil {
			t.Fatalf("移動檔案失敗：%v", err)
		}
		
		// 驗證來源檔案不存在
		if fileRepo.FileExists(sourceFile) {
			t.Error("來源檔案應該不存在")
		}
		
		// 驗證目標檔案存在
		targetFile := filepath.Join(targetDir, "source.md")
		if !fileRepo.FileExists(targetFile) {
			t.Error("目標檔案應該存在")
		}
		
		// 驗證檔案內容
		data, err := fileRepo.ReadFile(targetFile)
		if err != nil {
			t.Fatalf("讀取移動後的檔案失敗：%v", err)
		}
		
		if string(data) != testContent {
			t.Errorf("檔案內容不符合預期，期望：%s，實際：%s", testContent, string(data))
		}
	})
	
	// 測試案例：移動檔案到新位置
	t.Run("移動檔案到新位置", func(t *testing.T) {
		sourceFile := "move_source.txt"
		targetFile := "moved_file.txt"
		testContent := "移動測試內容"
		
		// 建立來源檔案
		fileRepo.WriteFile(sourceFile, []byte(testContent))
		
		// 移動檔案
		err := service.MoveFile(sourceFile, targetFile)
		if err != nil {
			t.Fatalf("移動檔案失敗：%v", err)
		}
		
		// 驗證來源檔案不存在
		if fileRepo.FileExists(sourceFile) {
			t.Error("來源檔案應該不存在")
		}
		
		// 驗證目標檔案存在
		if !fileRepo.FileExists(targetFile) {
			t.Error("目標檔案應該存在")
		}
		
		// 驗證檔案內容
		data, err := fileRepo.ReadFile(targetFile)
		if err != nil {
			t.Fatalf("讀取移動後的檔案失敗：%v", err)
		}
		
		if string(data) != testContent {
			t.Errorf("檔案內容不符合預期，期望：%s，實際：%s", testContent, string(data))
		}
	})
	
	// 測試案例：移動目錄
	t.Run("移動目錄", func(t *testing.T) {
		sourceDir := "source_dir"
		targetDir := "target_dir"
		testFile := filepath.Join(sourceDir, "test.txt")
		testContent := "目錄移動測試"
		
		// 建立來源目錄和檔案
		service.CreateDirectory(sourceDir)
		fileRepo.WriteFile(testFile, []byte(testContent))
		
		// 移動目錄
		err := service.MoveFile(sourceDir, targetDir)
		if err != nil {
			t.Fatalf("移動目錄失敗：%v", err)
		}
		
		// 驗證來源目錄不存在
		if fileRepo.FileExists(sourceDir) {
			t.Error("來源目錄應該不存在")
		}
		
		// 驗證目標目錄存在
		if !fileRepo.FileExists(targetDir) {
			t.Error("目標目錄應該存在")
		}
		
		// 驗證目錄內的檔案也被移動
		movedFile := filepath.Join(targetDir, "test.txt")
		if !fileRepo.FileExists(movedFile) {
			t.Error("目錄內的檔案應該被移動")
		}
		
		// 驗證檔案內容
		data, err := fileRepo.ReadFile(movedFile)
		if err != nil {
			t.Fatalf("讀取移動後的檔案失敗：%v", err)
		}
		
		if string(data) != testContent {
			t.Errorf("檔案內容不符合預期，期望：%s，實際：%s", testContent, string(data))
		}
	})
	
	// 測試案例：移動不存在的檔案應該回傳錯誤
	t.Run("移動不存在的檔案應該回傳錯誤", func(t *testing.T) {
		err := service.MoveFile("nonexistent.txt", "target.txt")
		if err == nil {
			t.Fatal("移動不存在的檔案應該回傳錯誤")
		}
		
		// 檢查錯誤類型
		if appErr, ok := err.(*models.AppError); ok {
			if appErr.Code != models.ErrFileNotFound {
				t.Errorf("錯誤代碼不符合預期，期望：%s，實際：%s", models.ErrFileNotFound, appErr.Code)
			}
		}
	})
	
	// 測試案例：移動到已存在的檔案應該回傳錯誤
	t.Run("移動到已存在的檔案應該回傳錯誤", func(t *testing.T) {
		sourceFile := "move_source2.txt"
		targetFile := "move_target2.txt"
		
		// 建立來源和目標檔案
		fileRepo.WriteFile(sourceFile, []byte("來源內容"))
		fileRepo.WriteFile(targetFile, []byte("目標內容"))
		
		// 嘗試移動到已存在的檔案
		err := service.MoveFile(sourceFile, targetFile)
		if err == nil {
			t.Fatal("移動到已存在的檔案應該回傳錯誤")
		}
		
		// 檢查錯誤類型
		if appErr, ok := err.(*models.AppError); ok {
			if appErr.Code != models.ErrValidationFailed {
				t.Errorf("錯誤代碼不符合預期，期望：%s，實際：%s", models.ErrValidationFailed, appErr.Code)
			}
		}
	})
}

// TestSearchFiles 測試檔案搜尋功能
func TestSearchFiles(t *testing.T) {
	// 建立測試環境
	tempDir := t.TempDir()
	fileRepo, _ := repositories.NewLocalFileRepository(tempDir)
	service, _ := NewLocalFileManagerService(fileRepo, tempDir)
	
	// 建立測試檔案結構
	testFiles := []string{
		"note1.md",
		"note2.md",
		"document.txt",
		"readme.MD",
		"folder1/nested_note.md",
		"folder1/other.txt",
		"folder2/another_note.md",
	}
	
	for _, file := range testFiles {
		fileRepo.WriteFile(file, []byte("測試內容"))
	}
	
	// 建立目錄
	service.CreateDirectory("folder1")
	service.CreateDirectory("folder2")
	
	// 測試案例：搜尋 Markdown 檔案（不包含子目錄）
	t.Run("搜尋 Markdown 檔案（不包含子目錄）", func(t *testing.T) {
		results, err := service.SearchFiles(".", "*.md", false)
		if err != nil {
			t.Fatalf("搜尋檔案失敗：%v", err)
		}
		
		// 應該找到根目錄下的 .md 檔案
		expectedCount := 2 // note1.md, note2.md
		if len(results) != expectedCount {
			t.Errorf("搜尋結果數量不符合預期，期望：%d，實際：%d", expectedCount, len(results))
		}
		
		// 檢查結果是否包含預期的檔案
		foundFiles := make(map[string]bool)
		for _, result := range results {
			foundFiles[result.Name] = true
		}
		
		expectedFiles := []string{"note1.md", "note2.md"}
		for _, expectedFile := range expectedFiles {
			if !foundFiles[expectedFile] {
				t.Errorf("搜尋結果中找不到預期的檔案：%s", expectedFile)
			}
		}
	})
	
	// 測試案例：搜尋 Markdown 檔案（包含子目錄）
	t.Run("搜尋 Markdown 檔案（包含子目錄）", func(t *testing.T) {
		results, err := service.SearchFiles(".", "*.md", true)
		if err != nil {
			t.Fatalf("搜尋檔案失敗：%v", err)
		}
		
		// 應該找到所有 .md 檔案
		expectedCount := 4 // note1.md, note2.md, nested_note.md, another_note.md
		if len(results) != expectedCount {
			t.Errorf("搜尋結果數量不符合預期，期望：%d，實際：%d", expectedCount, len(results))
		}
	})
	
	// 測試案例：搜尋特定目錄
	t.Run("搜尋特定目錄", func(t *testing.T) {
		results, err := service.SearchFiles("folder1", "*", false)
		if err != nil {
			t.Fatalf("搜尋特定目錄失敗：%v", err)
		}
		
		// 應該找到 folder1 中的所有檔案
		expectedCount := 2 // nested_note.md, other.txt
		if len(results) != expectedCount {
			t.Errorf("搜尋結果數量不符合預期，期望：%d，實際：%d", expectedCount, len(results))
		}
	})
	
	// 測試案例：空搜尋模式應該回傳錯誤
	t.Run("空搜尋模式應該回傳錯誤", func(t *testing.T) {
		_, err := service.SearchFiles(".", "", false)
		if err == nil {
			t.Fatal("空搜尋模式應該回傳錯誤")
		}
		
		// 檢查錯誤類型
		if appErr, ok := err.(*models.AppError); ok {
			if appErr.Code != models.ErrValidationFailed {
				t.Errorf("錯誤代碼不符合預期，期望：%s，實際：%s", models.ErrValidationFailed, appErr.Code)
			}
		}
	})
	
	// 測試案例：搜尋不存在的目錄應該回傳錯誤
	t.Run("搜尋不存在的目錄應該回傳錯誤", func(t *testing.T) {
		_, err := service.SearchFiles("nonexistent", "*.md", false)
		if err == nil {
			t.Fatal("搜尋不存在的目錄應該回傳錯誤")
		}
	})
}

// TestCopyFile 測試檔案複製功能
func TestCopyFile(t *testing.T) {
	// 建立測試環境
	tempDir := t.TempDir()
	fileRepo, _ := repositories.NewLocalFileRepository(tempDir)
	service, _ := NewLocalFileManagerService(fileRepo, tempDir)
	
	// 測試案例：複製檔案
	t.Run("複製檔案", func(t *testing.T) {
		sourceFile := "source_copy.md"
		targetFile := "target_copy.md"
		testContent := "複製測試內容"
		
		// 建立來源檔案
		fileRepo.WriteFile(sourceFile, []byte(testContent))
		
		// 複製檔案
		err := service.CopyFile(sourceFile, targetFile)
		if err != nil {
			t.Fatalf("複製檔案失敗：%v", err)
		}
		
		// 驗證來源檔案仍然存在
		if !fileRepo.FileExists(sourceFile) {
			t.Error("來源檔案應該仍然存在")
		}
		
		// 驗證目標檔案存在
		if !fileRepo.FileExists(targetFile) {
			t.Error("目標檔案應該存在")
		}
		
		// 驗證兩個檔案的內容相同
		sourceData, _ := fileRepo.ReadFile(sourceFile)
		targetData, _ := fileRepo.ReadFile(targetFile)
		
		if string(sourceData) != string(targetData) {
			t.Error("來源檔案和目標檔案的內容應該相同")
		}
		
		if string(targetData) != testContent {
			t.Errorf("目標檔案內容不符合預期，期望：%s，實際：%s", testContent, string(targetData))
		}
	})
	
	// 測試案例：複製目錄
	t.Run("複製目錄", func(t *testing.T) {
		sourceDir := "source_copy_dir"
		targetDir := "target_copy_dir"
		testFile := filepath.Join(sourceDir, "test.txt")
		testContent := "目錄複製測試"
		
		// 建立來源目錄和檔案
		service.CreateDirectory(sourceDir)
		fileRepo.WriteFile(testFile, []byte(testContent))
		
		// 複製目錄
		err := service.CopyFile(sourceDir, targetDir)
		if err != nil {
			t.Fatalf("複製目錄失敗：%v", err)
		}
		
		// 驗證來源目錄仍然存在
		if !fileRepo.FileExists(sourceDir) {
			t.Error("來源目錄應該仍然存在")
		}
		
		// 驗證目標目錄存在
		if !fileRepo.FileExists(targetDir) {
			t.Error("目標目錄應該存在")
		}
		
		// 驗證目錄內的檔案也被複製
		copiedFile := filepath.Join(targetDir, "test.txt")
		if !fileRepo.FileExists(copiedFile) {
			t.Error("目錄內的檔案應該被複製")
		}
		
		// 驗證檔案內容
		data, err := fileRepo.ReadFile(copiedFile)
		if err != nil {
			t.Fatalf("讀取複製後的檔案失敗：%v", err)
		}
		
		if string(data) != testContent {
			t.Errorf("檔案內容不符合預期，期望：%s，實際：%s", testContent, string(data))
		}
	})
	
	// 測試案例：複製不存在的檔案應該回傳錯誤
	t.Run("複製不存在的檔案應該回傳錯誤", func(t *testing.T) {
		err := service.CopyFile("nonexistent.txt", "target.txt")
		if err == nil {
			t.Fatal("複製不存在的檔案應該回傳錯誤")
		}
		
		// 檢查錯誤類型
		if appErr, ok := err.(*models.AppError); ok {
			if appErr.Code != models.ErrFileNotFound {
				t.Errorf("錯誤代碼不符合預期，期望：%s，實際：%s", models.ErrFileNotFound, appErr.Code)
			}
		}
	})
	
	// 測試案例：複製到已存在的路徑應該回傳錯誤
	t.Run("複製到已存在的路徑應該回傳錯誤", func(t *testing.T) {
		sourceFile := "copy_source3.txt"
		targetFile := "copy_target3.txt"
		
		// 建立來源和目標檔案
		fileRepo.WriteFile(sourceFile, []byte("來源內容"))
		fileRepo.WriteFile(targetFile, []byte("目標內容"))
		
		// 嘗試複製到已存在的檔案
		err := service.CopyFile(sourceFile, targetFile)
		if err == nil {
			t.Fatal("複製到已存在的路徑應該回傳錯誤")
		}
		
		// 檢查錯誤類型
		if appErr, ok := err.(*models.AppError); ok {
			if appErr.Code != models.ErrValidationFailed {
				t.Errorf("錯誤代碼不符合預期，期望：%s，實際：%s", models.ErrValidationFailed, appErr.Code)
			}
		}
	})
}

// TestGetDirectorySize 測試目錄大小計算功能
func TestGetDirectorySize(t *testing.T) {
	// 建立測試環境
	tempDir := t.TempDir()
	fileRepo, _ := repositories.NewLocalFileRepository(tempDir)
	service, _ := NewLocalFileManagerService(fileRepo, tempDir)
	
	// 測試案例：計算目錄大小
	t.Run("計算目錄大小", func(t *testing.T) {
		testDir := "size_test_dir"
		
		// 建立測試目錄
		service.CreateDirectory(testDir)
		
		// 建立測試檔案
		testFiles := map[string]string{
			filepath.Join(testDir, "file1.txt"): "內容1",
			filepath.Join(testDir, "file2.txt"): "內容22",
			filepath.Join(testDir, "file3.txt"): "內容333",
		}
		
		var expectedSize int64
		for filePath, content := range testFiles {
			fileRepo.WriteFile(filePath, []byte(content))
			expectedSize += int64(len(content))
		}
		
		// 計算目錄大小
		actualSize, err := service.GetDirectorySize(testDir)
		if err != nil {
			t.Fatalf("計算目錄大小失敗：%v", err)
		}
		
		if actualSize != expectedSize {
			t.Errorf("目錄大小不符合預期，期望：%d，實際：%d", expectedSize, actualSize)
		}
	})
	
	// 測試案例：計算巢狀目錄大小
	t.Run("計算巢狀目錄大小", func(t *testing.T) {
		rootDir := "nested_size_test"
		subDir := filepath.Join(rootDir, "subdir")
		
		// 建立目錄結構
		service.CreateDirectory(rootDir)
		service.CreateDirectory(subDir)
		
		// 建立測試檔案
		testFiles := map[string]string{
			filepath.Join(rootDir, "root.txt"):    "根目錄檔案",
			filepath.Join(subDir, "nested.txt"):   "子目錄檔案",
		}
		
		var expectedSize int64
		for filePath, content := range testFiles {
			fileRepo.WriteFile(filePath, []byte(content))
			expectedSize += int64(len(content))
		}
		
		// 計算根目錄大小（應該包含子目錄）
		actualSize, err := service.GetDirectorySize(rootDir)
		if err != nil {
			t.Fatalf("計算巢狀目錄大小失敗：%v", err)
		}
		
		if actualSize != expectedSize {
			t.Errorf("巢狀目錄大小不符合預期，期望：%d，實際：%d", expectedSize, actualSize)
		}
	})
	
	// 測試案例：計算不存在目錄的大小應該回傳錯誤
	t.Run("計算不存在目錄的大小應該回傳錯誤", func(t *testing.T) {
		_, err := service.GetDirectorySize("nonexistent_dir")
		if err == nil {
			t.Fatal("計算不存在目錄的大小應該回傳錯誤")
		}
		
		// 檢查錯誤類型
		if appErr, ok := err.(*models.AppError); ok {
			if appErr.Code != models.ErrFileNotFound {
				t.Errorf("錯誤代碼不符合預期，期望：%s，實際：%s", models.ErrFileNotFound, appErr.Code)
			}
		}
	})
}

// TestGetFileTree 測試檔案樹功能
func TestGetFileTree(t *testing.T) {
	// 建立測試環境
	tempDir := t.TempDir()
	fileRepo, _ := repositories.NewLocalFileRepository(tempDir)
	service, _ := NewLocalFileManagerService(fileRepo, tempDir)
	
	// 建立測試檔案結構
	testStructure := []string{
		"root.md",
		"folder1/file1.txt",
		"folder1/subfolder/deep.md",
		"folder2/file2.txt",
	}
	
	for _, path := range testStructure {
		fileRepo.WriteFile(path, []byte("測試內容"))
	}
	
	// 測試案例：取得檔案樹
	t.Run("取得檔案樹", func(t *testing.T) {
		tree, err := service.GetFileTree(".")
		if err != nil {
			t.Fatalf("取得檔案樹失敗：%v", err)
		}
		
		if tree == nil {
			t.Fatal("檔案樹不應為 nil")
		}
		
		// 驗證根節點是目錄
		if !tree.IsDirectory {
			t.Error("根節點應該是目錄")
		}
		
		// 驗證有子節點
		if !tree.HasChildren() {
			t.Error("根節點應該有子節點")
		}
		
		// 計算預期的子節點數量（1個檔案 + 2個目錄）
		expectedChildCount := 3
		if tree.GetChildCount() != expectedChildCount {
			t.Errorf("子節點數量不符合預期，期望：%d，實際：%d", expectedChildCount, tree.GetChildCount())
		}
	})
	
	// 測試案例：取得不存在目錄的檔案樹應該回傳錯誤
	t.Run("取得不存在目錄的檔案樹應該回傳錯誤", func(t *testing.T) {
		_, err := service.GetFileTree("nonexistent")
		if err == nil {
			t.Fatal("取得不存在目錄的檔案樹應該回傳錯誤")
		}
		
		// 檢查錯誤類型
		if appErr, ok := err.(*models.AppError); ok {
			if appErr.Code != models.ErrFileNotFound {
				t.Errorf("錯誤代碼不符合預期，期望：%s，實際：%s", models.ErrFileNotFound, appErr.Code)
			}
		}
	})
}

// TestPathValidation 測試路徑驗證功能
func TestPathValidation(t *testing.T) {
	// 建立測試環境
	tempDir := t.TempDir()
	fileRepo, _ := repositories.NewLocalFileRepository(tempDir)
	service, _ := NewLocalFileManagerService(fileRepo, tempDir)
	
	// 測試案例：有效路徑
	t.Run("有效路徑", func(t *testing.T) {
		validPaths := []string{
			".",
			"file.md",
			"folder/file.md",
			"deep/nested/path/file.md",
		}
		
		for _, path := range validPaths {
			err := service.validatePath(path)
			if err != nil {
				t.Errorf("有效路徑 %s 應該被接受，但發生錯誤：%v", path, err)
			}
		}
	})
	
	// 測試案例：無效路徑（通過 ListFiles 間接測試）
	t.Run("無效路徑", func(t *testing.T) {
		invalidPaths := []string{
			"",                    // 空路徑
			"../outside.md",       // 包含 ..
			"/absolute/path.md",   // 絕對路徑
			"../../outside.md",    // 真正的目錄遍歷攻擊
		}
		
		for _, path := range invalidPaths {
			// 使用 ListFiles 來間接測試路徑驗證
			_, err := service.ListFiles(path)
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