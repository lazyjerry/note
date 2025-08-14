// Package ui 包含檔案樹狀視圖元件的測試
// 測試檔案樹的建立、載入、選擇和各種互動功能
package ui

import (
	"os"                         // Go 標準庫，用於檔案系統操作
	"path/filepath"              // Go 標準庫，用於檔案路徑處理
	"testing"                    // Go 標準測試套件
	"fyne.io/fyne/v2/widget"     // Fyne UI 元件套件
	"mac-notebook-app/internal/models" // 本專案的資料模型套件
	"mac-notebook-app/internal/services" // 本專案的服務層套件
	"mac-notebook-app/internal/repositories" // 本專案的儲存庫層套件
)

// mockFileManagerService 模擬檔案管理服務
// 用於測試時提供可控制的檔案系統行為
type mockFileManagerService struct {
	files map[string][]*models.FileInfo // 模擬的檔案結構
}

// newMockFileManagerService 建立新的模擬檔案管理服務
// 回傳：配置好測試資料的模擬服務實例
func newMockFileManagerService() *mockFileManagerService {
	service := &mockFileManagerService{
		files: make(map[string][]*models.FileInfo),
	}
	
	// 設定測試用的檔案結構
	service.setupTestFileStructure()
	
	return service
}

// setupTestFileStructure 設定測試用的檔案結構
// 建立一個包含檔案和目錄的模擬檔案系統
func (m *mockFileManagerService) setupTestFileStructure() {
	// 根目錄檔案
	m.files["/test"] = []*models.FileInfo{
		{Path: "/test/notes", Name: "notes", IsDirectory: true},
		{Path: "/test/readme.md", Name: "readme.md", IsDirectory: false},
	}
	
	// notes 目錄檔案
	m.files["/test/notes"] = []*models.FileInfo{
		{Path: "/test/notes/work", Name: "work", IsDirectory: true},
		{Path: "/test/notes/personal", Name: "personal", IsDirectory: true},
		{Path: "/test/notes/todo.md", Name: "todo.md", IsDirectory: false},
	}
	
	// work 目錄檔案
	m.files["/test/notes/work"] = []*models.FileInfo{
		{Path: "/test/notes/work/project1.md", Name: "project1.md", IsDirectory: false},
		{Path: "/test/notes/work/meeting.md", Name: "meeting.md", IsDirectory: false},
	}
	
	// personal 目錄檔案
	m.files["/test/notes/personal"] = []*models.FileInfo{
		{Path: "/test/notes/personal/diary.md", Name: "diary.md", IsDirectory: false},
	}
}

// ListFiles 實作 FileManagerService 介面的 ListFiles 方法
func (m *mockFileManagerService) ListFiles(directory string) ([]*models.FileInfo, error) {
	files, exists := m.files[directory]
	if !exists {
		return []*models.FileInfo{}, nil
	}
	return files, nil
}

// 實作其他 FileManagerService 介面方法（測試中不使用）
func (m *mockFileManagerService) CreateDirectory(path string) error { return nil }
func (m *mockFileManagerService) DeleteFile(path string) error { return nil }
func (m *mockFileManagerService) RenameFile(oldPath, newPath string) error { return nil }
func (m *mockFileManagerService) MoveFile(sourcePath, destPath string) error { return nil }

// TestNewFileTreeWidget 測試檔案樹元件的建立
// 驗證檔案樹元件是否正確建立並包含所有必要的屬性
func TestNewFileTreeWidget(t *testing.T) {
	// 建立模擬檔案管理服務
	mockService := newMockFileManagerService()
	
	// 建立檔案樹元件
	fileTree := NewFileTreeWidget(mockService, "/test")
	
	// 驗證檔案樹元件不為 nil
	if fileTree == nil {
		t.Fatal("NewFileTreeWidget 應該回傳有效的 FileTreeWidget 實例")
	}
	
	// 驗證基本屬性
	if fileTree.rootPath != "/test" {
		t.Errorf("根目錄路徑應該是 '/test'，但得到 '%s'", fileTree.rootPath)
	}
	
	if fileTree.fileManager == nil {
		t.Error("檔案管理服務不應該為 nil")
	}
	
	if fileTree.tree == nil {
		t.Error("樹狀元件不應該為 nil")
	}
	
	if fileTree.container == nil {
		t.Error("容器元件不應該為 nil")
	}
	
	if fileTree.fileNodes == nil {
		t.Error("檔案節點快取不應該為 nil")
	}
}

// TestFileTreeLoadFileStructure 測試檔案結構載入
// 驗證檔案樹是否正確載入檔案和目錄結構
func TestFileTreeLoadFileStructure(t *testing.T) {
	// 建立模擬檔案管理服務
	mockService := newMockFileManagerService()
	
	// 建立檔案樹元件
	fileTree := NewFileTreeWidget(mockService, "/test")
	
	// 驗證根節點是否正確載入
	rootNode, exists := fileTree.fileNodes["/test"]
	if !exists {
		t.Fatal("根節點應該存在於檔案節點快取中")
	}
	
	if rootNode.Name != "test" {
		t.Errorf("根節點名稱應該是 'test'，但得到 '%s'", rootNode.Name)
	}
	
	if !rootNode.IsDirectory {
		t.Error("根節點應該是目錄")
	}
	
	if !rootNode.IsExpanded {
		t.Error("根節點應該是展開狀態")
	}
	
	// 驗證子節點是否正確載入
	if len(rootNode.Children) != 2 {
		t.Errorf("根節點應該有 2 個子節點，但得到 %d 個", len(rootNode.Children))
	}
	
	// 檢查子節點內容
	childNames := make(map[string]bool)
	for _, child := range rootNode.Children {
		childNames[child.Name] = true
	}
	
	if !childNames["notes"] {
		t.Error("應該包含 'notes' 目錄")
	}
	
	if !childNames["readme.md"] {
		t.Error("應該包含 'readme.md' 檔案")
	}
}

// TestFileTreeGetChildUIDs 測試子節點 ID 取得功能
// 驗證樹狀元件能否正確取得子節點的 ID 列表
func TestFileTreeGetChildUIDs(t *testing.T) {
	// 建立模擬檔案管理服務
	mockService := newMockFileManagerService()
	
	// 建立檔案樹元件
	fileTree := NewFileTreeWidget(mockService, "/test")
	
	// 測試根節點的子節點
	rootChildUIDs := fileTree.getChildUIDs("")
	if len(rootChildUIDs) != 1 {
		t.Errorf("根節點應該有 1 個子節點 ID，但得到 %d 個", len(rootChildUIDs))
	}
	
	if string(rootChildUIDs[0]) != "/test" {
		t.Errorf("根節點的子節點 ID 應該是 '/test'，但得到 '%s'", string(rootChildUIDs[0]))
	}
	
	// 測試 /test 節點的子節點
	testChildUIDs := fileTree.getChildUIDs(widget.TreeNodeID("/test"))
	if len(testChildUIDs) != 2 {
		t.Errorf("/test 節點應該有 2 個子節點 ID，但得到 %d 個", len(testChildUIDs))
	}
	
	// 測試檔案節點（應該沒有子節點）
	fileChildUIDs := fileTree.getChildUIDs(widget.TreeNodeID("/test/readme.md"))
	if len(fileChildUIDs) != 0 {
		t.Errorf("檔案節點不應該有子節點，但得到 %d 個", len(fileChildUIDs))
	}
}

// TestFileTreeIsBranch 測試分支檢查功能
// 驗證樹狀元件能否正確識別目錄和檔案
func TestFileTreeIsBranch(t *testing.T) {
	// 建立模擬檔案管理服務
	mockService := newMockFileManagerService()
	
	// 建立檔案樹元件
	fileTree := NewFileTreeWidget(mockService, "/test")
	
	// 測試根節點
	if !fileTree.isBranch("") {
		t.Error("根節點應該是分支")
	}
	
	// 測試目錄節點
	if !fileTree.isBranch(widget.TreeNodeID("/test")) {
		t.Error("/test 節點應該是分支")
	}
	
	if !fileTree.isBranch(widget.TreeNodeID("/test/notes")) {
		t.Error("/test/notes 節點應該是分支")
	}
	
	// 測試檔案節點
	if fileTree.isBranch(widget.TreeNodeID("/test/readme.md")) {
		t.Error("/test/readme.md 節點不應該是分支")
	}
	
	// 測試不存在的節點
	if fileTree.isBranch(widget.TreeNodeID("/nonexistent")) {
		t.Error("不存在的節點不應該是分支")
	}
}

// TestFileTreeNodeSelection 測試節點選擇功能
// 驗證檔案樹的節點選擇回調是否正確工作
func TestFileTreeNodeSelection(t *testing.T) {
	// 建立模擬檔案管理服務
	mockService := newMockFileManagerService()
	
	// 建立檔案樹元件
	fileTree := NewFileTreeWidget(mockService, "/test")
	
	// 設定回調函數來捕獲選擇事件
	var selectedFile string
	var selectedDirectory string
	
	fileTree.SetOnFileSelect(func(filePath string) {
		selectedFile = filePath
	})
	
	fileTree.SetOnDirectoryOpen(func(dirPath string) {
		selectedDirectory = dirPath
	})
	
	// 模擬選擇檔案
	fileTree.handleNodeSelection(widget.TreeNodeID("/test/readme.md"))
	if selectedFile != "/test/readme.md" {
		t.Errorf("檔案選擇回調應該收到 '/test/readme.md'，但得到 '%s'", selectedFile)
	}
	
	// 模擬選擇目錄
	fileTree.handleNodeSelection(widget.TreeNodeID("/test/notes"))
	if selectedDirectory != "/test/notes" {
		t.Errorf("目錄選擇回調應該收到 '/test/notes'，但得到 '%s'", selectedDirectory)
	}
}

// TestFileTreeBranchOperations 測試分支展開和收合操作
// 驗證目錄節點的展開和收合狀態管理
func TestFileTreeBranchOperations(t *testing.T) {
	// 建立模擬檔案管理服務
	mockService := newMockFileManagerService()
	
	// 建立檔案樹元件
	fileTree := NewFileTreeWidget(mockService, "/test")
	
	// 取得 notes 目錄節點
	notesNode, exists := fileTree.fileNodes["/test/notes"]
	if !exists {
		t.Fatal("notes 節點應該存在")
	}
	
	// 初始狀態應該是未展開
	if notesNode.IsExpanded {
		t.Error("notes 節點初始狀態應該是未展開")
	}
	
	// 模擬展開分支
	fileTree.handleBranchOpened(widget.TreeNodeID("/test/notes"))
	if !notesNode.IsExpanded {
		t.Error("展開後 notes 節點應該是展開狀態")
	}
	
	// 驗證子節點是否已載入
	if len(notesNode.Children) == 0 {
		t.Error("展開後應該載入子節點")
	}
	
	// 模擬收合分支
	fileTree.handleBranchClosed(widget.TreeNodeID("/test/notes"))
	if notesNode.IsExpanded {
		t.Error("收合後 notes 節點應該是收合狀態")
	}
}

// TestFileTreeRefresh 測試檔案樹刷新功能
// 驗證檔案樹能否正確重新載入檔案結構
func TestFileTreeRefresh(t *testing.T) {
	// 建立模擬檔案管理服務
	mockService := newMockFileManagerService()
	
	// 建立檔案樹元件
	fileTree := NewFileTreeWidget(mockService, "/test")
	
	// 記錄刷新前的節點數量（用於驗證刷新效果）
	_ = len(fileTree.fileNodes)
	
	// 修改模擬服務的檔案結構
	mockService.files["/test"] = append(mockService.files["/test"], &models.FileInfo{
		Path: "/test/new_file.md", Name: "new_file.md", IsDirectory: false,
	})
	
	// 執行刷新
	fileTree.Refresh()
	
	// 驗證節點快取是否重新建立
	if len(fileTree.fileNodes) == 0 {
		t.Error("刷新後應該重新載入節點")
	}
	
	// 驗證根節點是否存在
	rootNode, exists := fileTree.fileNodes["/test"]
	if !exists {
		t.Error("刷新後根節點應該存在")
	}
	
	// 驗證新檔案是否已載入
	if len(rootNode.Children) != 3 {
		t.Errorf("刷新後根節點應該有 3 個子節點，但得到 %d 個", len(rootNode.Children))
	}
}

// TestFileTreeGetSelectedPath 測試取得選擇路徑功能
// 驗證檔案樹能否正確回傳目前選擇的路徑
func TestFileTreeGetSelectedPath(t *testing.T) {
	// 建立模擬檔案管理服務
	mockService := newMockFileManagerService()
	
	// 建立檔案樹元件
	fileTree := NewFileTreeWidget(mockService, "/test")
	
	// 初始狀態應該沒有選擇
	selectedPath := fileTree.GetSelectedPath()
	if selectedPath != "" {
		t.Errorf("初始狀態應該沒有選擇，但得到 '%s'", selectedPath)
	}
}

// TestFileTreeExpandPath 測試路徑展開功能
// 驗證檔案樹能否正確展開指定路徑的所有父目錄
func TestFileTreeExpandPath(t *testing.T) {
	// 建立模擬檔案管理服務
	mockService := newMockFileManagerService()
	
	// 建立檔案樹元件
	fileTree := NewFileTreeWidget(mockService, "/test")
	
	// 先載入 notes 目錄的子項目
	fileTree.handleBranchOpened(widget.TreeNodeID("/test/notes"))
	
	// 展開深層路徑
	fileTree.ExpandPath("/test/notes/work/project1.md")
	
	// 驗證所有父目錄都已展開
	testNode := fileTree.fileNodes["/test"]
	if !testNode.IsExpanded {
		t.Error("/test 節點應該已展開")
	}
	
	notesNode := fileTree.fileNodes["/test/notes"]
	if notesNode == nil {
		t.Error("/test/notes 節點應該存在")
	} else if !notesNode.IsExpanded {
		t.Error("/test/notes 節點應該已展開")
	}
}

// TestFileTreeWithRealFileSystem 測試與真實檔案系統的整合
// 使用臨時目錄測試檔案樹與真實檔案系統的互動
func TestFileTreeWithRealFileSystem(t *testing.T) {
	// 建立臨時目錄
	tempDir, err := os.MkdirTemp("", "fileTree_test")
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
	os.WriteFile(testFile, []byte("# Test"), 0644)
	
	noteFile := filepath.Join(notesDir, "note1.md")
	os.WriteFile(noteFile, []byte("# Note 1"), 0644)
	
	// 建立真實的檔案管理服務
	fileRepo, err := repositories.NewLocalFileRepository(tempDir)
	if err != nil {
		t.Fatalf("建立檔案儲存庫失敗: %v", err)
	}
	
	fileManager, err := services.NewLocalFileManagerService(fileRepo, tempDir)
	if err != nil {
		t.Fatalf("建立檔案管理服務失敗: %v", err)
	}
	
	// 建立檔案樹元件（使用相對路徑）
	fileTree := NewFileTreeWidget(fileManager, "test")
	
	// 驗證檔案樹是否正確載入真實檔案結構
	rootNode, exists := fileTree.fileNodes["test"]
	if !exists {
		t.Fatal("根節點應該存在")
	}
	
	if len(rootNode.Children) == 0 {
		t.Error("應該載入子檔案和目錄")
	}
	
	// 驗證檔案和目錄是否正確識別
	foundReadme := false
	foundNotes := false
	
	for _, child := range rootNode.Children {
		if child.Name == "readme.md" && !child.IsDirectory {
			foundReadme = true
		}
		if child.Name == "notes" && child.IsDirectory {
			foundNotes = true
		}
	}
	
	if !foundReadme {
		t.Error("應該找到 readme.md 檔案")
	}
	
	if !foundNotes {
		t.Error("應該找到 notes 目錄")
	}
}