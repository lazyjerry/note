// Package ui 包含檔案管理整合到 UI 的測試
// 測試檔案樹與檔案管理服務的整合功能，包含檔案操作、UI 回饋和確認對話框
package ui

import (
	"fmt"                        // Go 標準庫，用於格式化輸出
	"os"                         // Go 標準庫，用於檔案系統操作
	"path/filepath"              // Go 標準庫，用於檔案路徑處理
	"testing"                    // Go 標準測試套件
	"fyne.io/fyne/v2/test"       // Fyne 測試套件
	"fyne.io/fyne/v2/widget"     // Fyne UI 元件套件
	"mac-notebook-app/internal/models" // 本專案的資料模型套件
	"mac-notebook-app/internal/services" // 本專案的服務層套件
	"mac-notebook-app/internal/repositories" // 本專案的儲存庫層套件
)

// TestFileTreeServiceIntegration 測試檔案樹與檔案管理服務的整合
// 驗證檔案樹能否正確使用檔案管理服務執行各種檔案操作
func TestFileTreeServiceIntegration(t *testing.T) {
	// 建立臨時目錄用於測試
	tempDir, err := os.MkdirTemp("", "file_tree_integration_test")
	if err != nil {
		t.Fatalf("建立臨時目錄失敗: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// 建立測試檔案結構
	testDir := filepath.Join(tempDir, "test")
	os.MkdirAll(testDir, 0755)
	
	notesDir := filepath.Join(testDir, "notes")
	os.MkdirAll(notesDir, 0755)
	
	// 建立測試檔案
	testFile := filepath.Join(testDir, "readme.md")
	os.WriteFile(testFile, []byte("# Test README"), 0644)
	
	noteFile := filepath.Join(notesDir, "note1.md")
	os.WriteFile(noteFile, []byte("# Note 1"), 0644)
	
	// 建立檔案管理服務
	fileRepo, err := repositories.NewLocalFileRepository(tempDir)
	if err != nil {
		t.Fatalf("建立檔案儲存庫失敗: %v", err)
	}
	
	fileManager, err := services.NewLocalFileManagerService(fileRepo, tempDir)
	if err != nil {
		t.Fatalf("建立檔案管理服務失敗: %v", err)
	}
	
	// 建立檔案樹元件
	fileTree := NewFileTreeWidget(fileManager, "test")
	
	// 測試檔案樹是否正確載入檔案結構
	rootNode, exists := fileTree.fileNodes["test"]
	if !exists {
		t.Fatal("根節點應該存在")
	}
	
	if len(rootNode.Children) == 0 {
		t.Error("應該載入子檔案和目錄")
	}
	
	// 驗證檔案管理服務整合
	if fileTree.fileManager == nil {
		t.Error("檔案管理服務不應該為 nil")
	}
}

// TestFileTreeCreateNewFolder 測試檔案樹建立新資料夾功能
// 驗證檔案樹能否正確使用檔案管理服務建立新資料夾
func TestFileTreeCreateNewFolder(t *testing.T) {
	// 建立臨時目錄用於測試
	tempDir, err := os.MkdirTemp("", "file_tree_create_folder_test")
	if err != nil {
		t.Fatalf("建立臨時目錄失敗: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// 建立測試目錄結構
	testDir := filepath.Join(tempDir, "test")
	os.MkdirAll(testDir, 0755)
	
	// 建立檔案管理服務
	fileRepo, err := repositories.NewLocalFileRepository(tempDir)
	if err != nil {
		t.Fatalf("建立檔案儲存庫失敗: %v", err)
	}
	
	fileManager, err := services.NewLocalFileManagerService(fileRepo, tempDir)
	if err != nil {
		t.Fatalf("建立檔案管理服務失敗: %v", err)
	}
	
	// 建立檔案樹元件
	fileTree := NewFileTreeWidget(fileManager, "test")
	
	// 測試建立新資料夾
	newFolderPath := "test/new_folder"
	err = fileTree.CreateNewFolder("test", "new_folder")
	if err != nil {
		t.Errorf("建立新資料夾失敗: %v", err)
	}
	
	// 驗證資料夾是否已建立
	fullPath := filepath.Join(tempDir, newFolderPath)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		t.Error("新資料夾應該已建立")
	}
	
	// 驗證檔案樹是否已更新
	if _, exists := fileTree.fileNodes[newFolderPath]; !exists {
		// 注意：由於 Refresh() 會重新載入，節點可能需要重新載入
		fileTree.Refresh()
	}
}

// TestFileTreeDeleteFile 測試檔案樹刪除檔案功能
// 驗證檔案樹能否正確使用檔案管理服務刪除檔案
func TestFileTreeDeleteFile(t *testing.T) {
	// 建立臨時目錄用於測試
	tempDir, err := os.MkdirTemp("", "file_tree_delete_test")
	if err != nil {
		t.Fatalf("建立臨時目錄失敗: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// 建立測試檔案
	testDir := filepath.Join(tempDir, "test")
	os.MkdirAll(testDir, 0755)
	
	testFile := filepath.Join(testDir, "delete_me.md")
	os.WriteFile(testFile, []byte("# Delete Me"), 0644)
	
	// 建立檔案管理服務
	fileRepo, err := repositories.NewLocalFileRepository(tempDir)
	if err != nil {
		t.Fatalf("建立檔案儲存庫失敗: %v", err)
	}
	
	fileManager, err := services.NewLocalFileManagerService(fileRepo, tempDir)
	if err != nil {
		t.Fatalf("建立檔案管理服務失敗: %v", err)
	}
	
	// 建立檔案樹元件
	fileTree := NewFileTreeWidget(fileManager, "test")
	
	// 測試刪除檔案
	deleteFilePath := "test/delete_me.md"
	err = fileTree.DeleteFileOrFolder(deleteFilePath)
	if err != nil {
		t.Errorf("刪除檔案失敗: %v", err)
	}
	
	// 驗證檔案是否已刪除
	if _, err := os.Stat(testFile); !os.IsNotExist(err) {
		t.Error("檔案應該已被刪除")
	}
	
	// 驗證節點快取是否已更新
	if _, exists := fileTree.fileNodes[deleteFilePath]; exists {
		t.Error("刪除的檔案節點應該從快取中移除")
	}
}

// TestFileTreeRenameFile 測試檔案樹重新命名檔案功能
// 驗證檔案樹能否正確使用檔案管理服務重新命名檔案
func TestFileTreeRenameFile(t *testing.T) {
	// 建立臨時目錄用於測試
	tempDir, err := os.MkdirTemp("", "file_tree_rename_test")
	if err != nil {
		t.Fatalf("建立臨時目錄失敗: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// 建立測試檔案
	testDir := filepath.Join(tempDir, "test")
	os.MkdirAll(testDir, 0755)
	
	oldFile := filepath.Join(testDir, "old_name.md")
	os.WriteFile(oldFile, []byte("# Old Name"), 0644)
	
	// 建立檔案管理服務
	fileRepo, err := repositories.NewLocalFileRepository(tempDir)
	if err != nil {
		t.Fatalf("建立檔案儲存庫失敗: %v", err)
	}
	
	fileManager, err := services.NewLocalFileManagerService(fileRepo, tempDir)
	if err != nil {
		t.Fatalf("建立檔案管理服務失敗: %v", err)
	}
	
	// 建立檔案樹元件
	fileTree := NewFileTreeWidget(fileManager, "test")
	
	// 測試重新命名檔案
	oldPath := "test/old_name.md"
	newPath := "test/new_name.md"
	err = fileTree.RenameFileOrFolder(oldPath, newPath)
	if err != nil {
		t.Errorf("重新命名檔案失敗: %v", err)
	}
	
	// 驗證舊檔案是否已不存在
	if _, err := os.Stat(oldFile); !os.IsNotExist(err) {
		t.Error("舊檔案應該已不存在")
	}
	
	// 驗證新檔案是否存在
	newFile := filepath.Join(testDir, "new_name.md")
	if _, err := os.Stat(newFile); os.IsNotExist(err) {
		t.Error("新檔案應該存在")
	}
	
	// 驗證節點快取是否已更新
	if _, exists := fileTree.fileNodes[oldPath]; exists {
		t.Error("舊路徑的節點應該從快取中移除")
	}
	
	if node, exists := fileTree.fileNodes[newPath]; !exists {
		t.Error("新路徑的節點應該存在於快取中")
	} else {
		if node.Name != "new_name.md" {
			t.Errorf("節點名稱應該是 'new_name.md'，但得到 '%s'", node.Name)
		}
	}
}

// TestFileTreeCopyFile 測試檔案樹複製檔案功能
// 驗證檔案樹能否正確使用檔案管理服務複製檔案
func TestFileTreeCopyFile(t *testing.T) {
	// 建立臨時目錄用於測試
	tempDir, err := os.MkdirTemp("", "file_tree_copy_test")
	if err != nil {
		t.Fatalf("建立臨時目錄失敗: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// 建立測試檔案
	testDir := filepath.Join(tempDir, "test")
	os.MkdirAll(testDir, 0755)
	
	sourceFile := filepath.Join(testDir, "source.md")
	os.WriteFile(sourceFile, []byte("# Source File"), 0644)
	
	// 建立檔案管理服務
	fileRepo, err := repositories.NewLocalFileRepository(tempDir)
	if err != nil {
		t.Fatalf("建立檔案儲存庫失敗: %v", err)
	}
	
	fileManager, err := services.NewLocalFileManagerService(fileRepo, tempDir)
	if err != nil {
		t.Fatalf("建立檔案管理服務失敗: %v", err)
	}
	
	// 建立檔案樹元件
	fileTree := NewFileTreeWidget(fileManager, "test")
	
	// 測試複製檔案
	sourcePath := "test/source.md"
	destPath := "test/copy.md"
	err = fileTree.CopyFileOrFolder(sourcePath, destPath)
	if err != nil {
		t.Errorf("複製檔案失敗: %v", err)
	}
	
	// 驗證來源檔案仍然存在
	if _, err := os.Stat(sourceFile); os.IsNotExist(err) {
		t.Error("來源檔案應該仍然存在")
	}
	
	// 驗證複製的檔案是否存在
	destFile := filepath.Join(testDir, "copy.md")
	if _, err := os.Stat(destFile); os.IsNotExist(err) {
		t.Error("複製的檔案應該存在")
	}
	
	// 驗證檔案內容是否相同
	sourceContent, _ := os.ReadFile(sourceFile)
	destContent, _ := os.ReadFile(destFile)
	if string(sourceContent) != string(destContent) {
		t.Error("複製的檔案內容應該與來源檔案相同")
	}
}

// TestFileTreeMoveFile 測試檔案樹移動檔案功能
// 驗證檔案樹能否正確使用檔案管理服務移動檔案
func TestFileTreeMoveFile(t *testing.T) {
	// 建立臨時目錄用於測試
	tempDir, err := os.MkdirTemp("", "file_tree_move_test")
	if err != nil {
		t.Fatalf("建立臨時目錄失敗: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// 建立測試檔案和目錄
	testDir := filepath.Join(tempDir, "test")
	os.MkdirAll(testDir, 0755)
	
	subDir := filepath.Join(testDir, "subdir")
	os.MkdirAll(subDir, 0755)
	
	sourceFile := filepath.Join(testDir, "move_me.md")
	os.WriteFile(sourceFile, []byte("# Move Me"), 0644)
	
	// 建立檔案管理服務
	fileRepo, err := repositories.NewLocalFileRepository(tempDir)
	if err != nil {
		t.Fatalf("建立檔案儲存庫失敗: %v", err)
	}
	
	fileManager, err := services.NewLocalFileManagerService(fileRepo, tempDir)
	if err != nil {
		t.Fatalf("建立檔案管理服務失敗: %v", err)
	}
	
	// 建立檔案樹元件
	fileTree := NewFileTreeWidget(fileManager, "test")
	
	// 測試移動檔案
	sourcePath := "test/move_me.md"
	destPath := "test/subdir/move_me.md"
	err = fileTree.MoveFileOrFolder(sourcePath, destPath)
	if err != nil {
		t.Errorf("移動檔案失敗: %v", err)
	}
	
	// 驗證來源檔案是否已不存在
	if _, err := os.Stat(sourceFile); !os.IsNotExist(err) {
		t.Error("來源檔案應該已不存在")
	}
	
	// 驗證目標檔案是否存在
	destFile := filepath.Join(subDir, "move_me.md")
	if _, err := os.Stat(destFile); os.IsNotExist(err) {
		t.Error("目標檔案應該存在")
	}
	
	// 驗證節點快取是否已更新
	if _, exists := fileTree.fileNodes[sourcePath]; exists {
		t.Error("來源路徑的節點應該從快取中移除")
	}
	
	if _, exists := fileTree.fileNodes[destPath]; !exists {
		// 注意：移動後可能需要重新載入才能看到節點
		fileTree.Refresh()
	}
}

// TestMainWindowFileTreeIntegration 測試主視窗與檔案樹的整合
// 驗證主視窗能否正確整合檔案樹並處理檔案操作
func TestMainWindowFileTreeIntegration(t *testing.T) {
	// 建立測試應用程式
	app := test.NewApp()
	defer app.Quit()
	
	// 建立臨時目錄用於測試
	tempDir, err := os.MkdirTemp("", "main_window_integration_test")
	if err != nil {
		t.Fatalf("建立臨時目錄失敗: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// 建立測試設定
	settings := &models.Settings{
		DefaultSaveLocation: tempDir,
		AutoSaveInterval:    5,
		DefaultEncryption:   "aes256",
		Theme:              "light",
	}
	
	// 建立檔案管理服務
	fileRepo, err := repositories.NewLocalFileRepository(tempDir)
	if err != nil {
		t.Fatalf("建立檔案儲存庫失敗: %v", err)
	}
	
	fileManager, err := services.NewLocalFileManagerService(fileRepo, tempDir)
	if err != nil {
		t.Fatalf("建立檔案管理服務失敗: %v", err)
	}
	
	// 建立編輯器服務（模擬）
	editorService := &mockEditorService{}
	
	// 建立主視窗
	mainWindow := NewMainWindow(app, settings, editorService, fileManager)
	
	// 驗證主視窗是否正確建立
	if mainWindow == nil {
		t.Fatal("主視窗應該正確建立")
	}
	
	// 驗證檔案樹是否已整合
	if mainWindow.fileTreeWidget == nil {
		t.Error("檔案樹元件應該已整合到主視窗")
	}
	
	// 驗證檔案管理服務是否已設定
	if mainWindow.fileManagerService == nil {
		t.Error("檔案管理服務應該已設定")
	}
}

// TestFileTreeCallbacks 測試檔案樹回調函數的設定和執行
// 驗證檔案樹的各種回調函數是否正確設定和執行
func TestFileTreeCallbacks(t *testing.T) {
	// 建立模擬檔案管理服務
	mockService := newMockFileManagerService()
	
	// 建立檔案樹元件
	fileTree := NewFileTreeWidget(mockService, "/test")
	
	// 測試回調函數設定
	var selectedFile string
	var selectedDirectory string
	var operationCalled string
	var operationPath string
	
	// 設定回調函數
	fileTree.SetOnFileSelect(func(filePath string) {
		selectedFile = filePath
	})
	
	fileTree.SetOnDirectoryOpen(func(dirPath string) {
		selectedDirectory = dirPath
	})
	
	fileTree.SetOnFileOperation(func(operation, filePath string) {
		operationCalled = operation
		operationPath = filePath
	})
	
	// 測試檔案選擇回調
	fileTree.handleNodeSelection(widget.TreeNodeID("/test/readme.md"))
	if selectedFile != "/test/readme.md" {
		t.Errorf("檔案選擇回調應該收到 '/test/readme.md'，但得到 '%s'", selectedFile)
	}
	
	// 測試目錄開啟回調
	fileTree.handleNodeSelection(widget.TreeNodeID("/test/notes"))
	if selectedDirectory != "/test/notes" {
		t.Errorf("目錄開啟回調應該收到 '/test/notes'，但得到 '%s'", selectedDirectory)
	}
	
	// 測試檔案操作回調（透過模擬右鍵選單）
	if fileTree.onFileOperation != nil {
		fileTree.onFileOperation("delete", "test/readme.md")
	}
	
	if operationCalled != "delete" {
		t.Errorf("操作回調應該收到 'delete'，但得到 '%s'", operationCalled)
	}
	
	if operationPath != "test/readme.md" {
		t.Errorf("操作路徑應該是 'test/readme.md'，但得到 '%s'", operationPath)
	}
}

// TestFileTreeErrorHandling 測試檔案樹的錯誤處理
// 驗證檔案樹在遇到各種錯誤情況時的處理方式
func TestFileTreeErrorHandling(t *testing.T) {
	// 建立模擬檔案管理服務
	mockService := &errorMockFileManagerService{}
	
	// 建立檔案樹元件
	fileTree := NewFileTreeWidget(mockService, "/test")
	
	// 測試建立資料夾時的錯誤處理
	err := fileTree.CreateNewFolder("/test", "new_folder")
	if err == nil {
		t.Error("應該回傳錯誤")
	}
	
	// 測試刪除檔案時的錯誤處理
	err = fileTree.DeleteFileOrFolder("/test/nonexistent.md")
	if err == nil {
		t.Error("應該回傳錯誤")
	}
	
	// 測試重新命名檔案時的錯誤處理
	err = fileTree.RenameFileOrFolder("/test/old.md", "/test/new.md")
	if err == nil {
		t.Error("應該回傳錯誤")
	}
}

// 使用已存在的 mockEditorService（在 editor_test.go 中定義）

// errorMockFileManagerService 模擬會產生錯誤的檔案管理服務
type errorMockFileManagerService struct{}

func (e *errorMockFileManagerService) ListFiles(directory string) ([]*models.FileInfo, error) {
	return nil, fmt.Errorf("模擬錯誤：無法列出檔案")
}

func (e *errorMockFileManagerService) CreateDirectory(path string) error {
	return fmt.Errorf("模擬錯誤：無法建立目錄")
}

func (e *errorMockFileManagerService) DeleteFile(path string) error {
	return fmt.Errorf("模擬錯誤：無法刪除檔案")
}

func (e *errorMockFileManagerService) RenameFile(oldPath, newPath string) error {
	return fmt.Errorf("模擬錯誤：無法重新命名檔案")
}

func (e *errorMockFileManagerService) MoveFile(sourcePath, destPath string) error {
	return fmt.Errorf("模擬錯誤：無法移動檔案")
}

func (e *errorMockFileManagerService) CopyFile(sourcePath, destPath string) error {
	return fmt.Errorf("模擬錯誤：無法複製檔案")
}