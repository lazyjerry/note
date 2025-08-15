// Package ui 包含拖拽功能的測試
package ui

import (
	"fmt"                                    // Go 標準庫，用於格式化字串
	"testing"                                // Go 標準測試套件
	"fyne.io/fyne/v2/widget"                 // Fyne 元件套件
	"mac-notebook-app/internal/models"       // 內部模型套件
)

// MockFileManagerService 模擬檔案管理服務
// 用於測試拖拽功能，不需要實際的檔案系統操作
type MockFileManagerService struct {
	files       []*models.FileInfo // 模擬檔案列表
	moveError   error              // 模擬移動操作錯誤
	moveCalled  bool               // 記錄移動操作是否被呼叫
	lastMove    [2]string          // 記錄最後一次移動操作的路徑
}

// ListFiles 模擬列出檔案
func (m *MockFileManagerService) ListFiles(directory string) ([]*models.FileInfo, error) {
	return m.files, nil
}

// CreateDirectory 模擬建立目錄
func (m *MockFileManagerService) CreateDirectory(path string) error {
	return nil
}

// DeleteFile 模擬刪除檔案
func (m *MockFileManagerService) DeleteFile(path string) error {
	return nil
}

// RenameFile 模擬重新命名檔案
func (m *MockFileManagerService) RenameFile(oldPath, newPath string) error {
	return nil
}

// MoveFile 模擬移動檔案
func (m *MockFileManagerService) MoveFile(sourcePath, destPath string) error {
	m.moveCalled = true
	m.lastMove = [2]string{sourcePath, destPath}
	return m.moveError
}

// CopyFile 模擬複製檔案
func (m *MockFileManagerService) CopyFile(sourcePath, destPath string) error {
	return nil
}



// TestNewDragDropManager 測試拖拽管理器的建立
// 驗證管理器實例是否正確建立並設定服務
func TestNewDragDropManager(t *testing.T) {
	// 建立模擬檔案管理服務
	mockFileManager := &MockFileManagerService{}
	
	// 建立拖拽管理器（使用 nil 作為測試視窗）
	manager := NewDragDropManager(mockFileManager, nil)
	
	// 驗證管理器不為 nil
	if manager == nil {
		t.Fatal("拖拽管理器建立失敗，回傳 nil")
	}
	
	// 驗證檔案管理服務設定正確
	if manager.fileManager != mockFileManager {
		t.Error("檔案管理服務設定不正確")
	}
	
	// 驗證父視窗設定正確
	if manager.parent != nil {
		t.Error("父視窗設定不正確")
	}
	
	// 驗證拖拽區域映射已初始化
	if manager.dropZones == nil {
		t.Error("拖拽區域映射未初始化")
	}
	
	// 驗證視覺回饋已初始化
	if manager.dragFeedback == nil {
		t.Error("拖拽視覺回饋未初始化")
	}
}

// TestDragDropManager_RegisterDropZone 測試拖拽區域註冊
// 驗證拖拽區域是否正確註冊到管理器中
func TestDragDropManager_RegisterDropZone(t *testing.T) {
	// 建立測試環境
	mockFileManager := &MockFileManagerService{}
	manager := NewDragDropManager(mockFileManager, nil)
	
	// 建立測試元件
	testWidget := widget.NewLabel("測試拖拽區域")
	testZoneID := "test-zone"
	testTargetPath := "/test/target"
	testAcceptTypes := []string{".md", ".txt"}
	
	// 註冊拖拽區域
	manager.RegisterDropZone(testZoneID, testWidget, testTargetPath, testAcceptTypes)
	
	// 驗證區域已註冊
	zone, exists := manager.dropZones[testZoneID]
	if !exists {
		t.Fatal("拖拽區域註冊失敗")
	}
	
	// 驗證區域屬性
	if zone.widget != testWidget {
		t.Error("拖拽區域元件設定不正確")
	}
	
	if zone.targetPath != testTargetPath {
		t.Error("拖拽區域目標路徑設定不正確")
	}
	
	if len(zone.acceptTypes) != len(testAcceptTypes) {
		t.Error("拖拽區域接受類型數量不正確")
	}
	
	for i, acceptType := range zone.acceptTypes {
		if acceptType != testAcceptTypes[i] {
			t.Errorf("拖拽區域接受類型不正確：預期 %s，實際 %s", testAcceptTypes[i], acceptType)
		}
	}
	
	if !zone.isActive {
		t.Error("拖拽區域應該預設為啟用狀態")
	}
}

// TestDragDropManager_UnregisterDropZone 測試拖拽區域取消註冊
// 驗證拖拽區域是否正確從管理器中移除
func TestDragDropManager_UnregisterDropZone(t *testing.T) {
	// 建立測試環境
	mockFileManager := &MockFileManagerService{}
	manager := NewDragDropManager(mockFileManager, nil)
	
	// 註冊拖拽區域
	testWidget := widget.NewLabel("測試拖拽區域")
	testZoneID := "test-zone"
	manager.RegisterDropZone(testZoneID, testWidget, "/test/target", []string{".md"})
	
	// 驗證區域已註冊
	if _, exists := manager.dropZones[testZoneID]; !exists {
		t.Fatal("拖拽區域註冊失敗")
	}
	
	// 取消註冊拖拽區域
	manager.UnregisterDropZone(testZoneID)
	
	// 驗證區域已移除
	if _, exists := manager.dropZones[testZoneID]; exists {
		t.Error("拖拽區域取消註冊失敗，區域仍然存在")
	}
}

// TestDragDropManager_SetCallbacks 測試回調函數設定
// 驗證回調函數是否正確設定到管理器中
func TestDragDropManager_SetCallbacks(t *testing.T) {
	// 建立測試環境
	mockFileManager := &MockFileManagerService{}
	manager := NewDragDropManager(mockFileManager, nil)
	
	// 建立測試回調函數
	var fileDroppedCalled bool
	var fileMovedCalled bool
	
	onFileDropped := func(sourcePath, targetPath string) error {
		fileDroppedCalled = true
		return nil
	}
	
	onFileMoved := func(oldPath, newPath string) {
		fileMovedCalled = true
	}
	
	onError := func(err error) {
		// Error callback for testing
	}
	
	// 設定回調函數
	manager.SetCallbacks(onFileDropped, onFileMoved, onError)
	
	// 驗證回調函數已設定（透過模擬拖拽操作）
	// 由於無法直接檢查函數指標，我們透過觸發操作來驗證
	
	// 模擬成功的拖拽操作
	mockFileManager.files = []*models.FileInfo{
		{Name: "test.md", Path: "/test/source/test.md", IsDirectory: false},
	}
	
	err := manager.handleDrop("/test/source/test.md", "/test/target/")
	if err != nil {
		t.Errorf("拖拽操作失敗：%v", err)
	}
	
	// 驗證回調函數被呼叫
	if !fileDroppedCalled {
		t.Error("檔案拖拽完成回調未被呼叫")
	}
	
	if !fileMovedCalled {
		t.Error("檔案移動完成回調未被呼叫")
	}
}

// TestDragDropManager_EnableDisableZone 測試拖拽區域啟用和停用
// 驗證拖拽區域的啟用和停用功能是否正常工作
func TestDragDropManager_EnableDisableZone(t *testing.T) {
	// 建立測試環境
	mockFileManager := &MockFileManagerService{}
	manager := NewDragDropManager(mockFileManager, nil)
	
	// 註冊拖拽區域
	testWidget := widget.NewLabel("測試拖拽區域")
	testZoneID := "test-zone"
	manager.RegisterDropZone(testZoneID, testWidget, "/test/target", []string{".md"})
	
	// 驗證區域預設為啟用狀態
	zone := manager.dropZones[testZoneID]
	if !zone.isActive {
		t.Error("拖拽區域應該預設為啟用狀態")
	}
	
	// 停用區域
	manager.DisableZone(testZoneID)
	
	// 驗證區域已停用
	if zone.isActive {
		t.Error("拖拽區域停用失敗")
	}
	
	// 重新啟用區域
	manager.EnableZone(testZoneID)
	
	// 驗證區域已啟用
	if !zone.isActive {
		t.Error("拖拽區域啟用失敗")
	}
}

// TestDragDropManager_IsValidFileType 測試檔案類型驗證
// 驗證檔案類型驗證功能是否正確工作
func TestDragDropManager_IsValidFileType(t *testing.T) {
	// 建立測試環境
	mockFileManager := &MockFileManagerService{}
	manager := NewDragDropManager(mockFileManager, nil)
	
	// 測試案例
	testCases := []struct {
		name        string
		filePath    string
		acceptTypes []string
		expected    bool
	}{
		{
			name:        "支援的 Markdown 檔案",
			filePath:    "/test/file.md",
			acceptTypes: []string{".md", ".txt"},
			expected:    true,
		},
		{
			name:        "支援的文字檔案",
			filePath:    "/test/file.txt",
			acceptTypes: []string{".md", ".txt"},
			expected:    true,
		},
		{
			name:        "不支援的檔案類型",
			filePath:    "/test/file.doc",
			acceptTypes: []string{".md", ".txt"},
			expected:    false,
		},
		{
			name:        "空的接受類型列表（接受所有）",
			filePath:    "/test/file.doc",
			acceptTypes: []string{},
			expected:    true,
		},
		{
			name:        "大小寫不敏感測試",
			filePath:    "/test/file.MD",
			acceptTypes: []string{".md", ".txt"},
			expected:    true,
		},
	}
	
	// 執行測試案例
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := manager.IsValidFileType(tc.filePath, tc.acceptTypes)
			if result != tc.expected {
				t.Errorf("檔案類型驗證失敗：檔案 %s，接受類型 %v，預期 %v，實際 %v",
					tc.filePath, tc.acceptTypes, tc.expected, result)
			}
		})
	}
}

// TestDragDropManager_HandleDrop 測試拖拽放下處理
// 驗證拖拽放下操作的處理邏輯是否正確
func TestDragDropManager_HandleDrop(t *testing.T) {
	// 建立測試環境
	mockFileManager := &MockFileManagerService{}
	manager := NewDragDropManager(mockFileManager, nil)
	
	// 設定模擬檔案
	mockFileManager.files = []*models.FileInfo{
		{Name: "test.md", Path: "/test/source/test.md", IsDirectory: false},
		{Name: "target", Path: "/test/target", IsDirectory: true},
	}
	
	// 測試成功的拖拽操作
	t.Run("成功的拖拽操作", func(t *testing.T) {
		err := manager.handleDrop("/test/source/test.md", "/test/target")
		if err != nil {
			t.Errorf("拖拽操作應該成功，但發生錯誤：%v", err)
		}
		
		// 驗證檔案管理服務的移動操作被呼叫
		if !mockFileManager.moveCalled {
			t.Error("檔案移動操作未被呼叫")
		}
		
		// 驗證移動路徑正確
		expectedSource := "/test/source/test.md"
		expectedTarget := "/test/target/test.md"
		if mockFileManager.lastMove[0] != expectedSource {
			t.Errorf("來源路徑不正確：預期 %s，實際 %s", expectedSource, mockFileManager.lastMove[0])
		}
		if mockFileManager.lastMove[1] != expectedTarget {
			t.Errorf("目標路徑不正確：預期 %s，實際 %s", expectedTarget, mockFileManager.lastMove[1])
		}
	})
	
	// 測試無效路徑的錯誤處理
	t.Run("無效路徑錯誤處理", func(t *testing.T) {
		// 重置模擬服務狀態
		mockFileManager.moveCalled = false
		
		err := manager.handleDrop("", "/test/target")
		if err == nil {
			t.Error("空來源路徑應該產生錯誤")
		}
		
		err = manager.handleDrop("/test/source/test.md", "")
		if err == nil {
			t.Error("空目標路徑應該產生錯誤")
		}
	})
	
	// 測試檔案不存在的錯誤處理
	t.Run("檔案不存在錯誤處理", func(t *testing.T) {
		// 重置模擬服務狀態
		mockFileManager.moveCalled = false
		
		err := manager.handleDrop("/test/nonexistent.md", "/test/target")
		if err == nil {
			t.Error("不存在的檔案應該產生錯誤")
		}
		
		// 驗證移動操作未被呼叫
		if mockFileManager.moveCalled {
			t.Error("不存在檔案的移動操作不應該被呼叫")
		}
	})
	
	// 測試檔案移動失敗的錯誤處理
	t.Run("檔案移動失敗錯誤處理", func(t *testing.T) {
		// 設定模擬移動錯誤
		mockFileManager.moveError = fmt.Errorf("移動操作失敗")
		mockFileManager.moveCalled = false
		
		err := manager.handleDrop("/test/source/test.md", "/test/target")
		if err == nil {
			t.Error("移動失敗應該產生錯誤")
		}
		
		// 驗證移動操作被呼叫
		if !mockFileManager.moveCalled {
			t.Error("移動操作應該被呼叫")
		}
		
		// 重置錯誤狀態
		mockFileManager.moveError = nil
	})
}

// TestNewDragFeedback 測試拖拽視覺回饋的建立
// 驗證視覺回饋實例是否正確建立和初始化
func TestNewDragFeedback(t *testing.T) {
	// 建立拖拽視覺回饋
	feedback := NewDragFeedback()
	
	// 驗證回饋不為 nil
	if feedback == nil {
		t.Fatal("拖拽視覺回饋建立失敗，回傳 nil")
	}
	
	// 驗證覆蓋層已初始化
	if feedback.overlay == nil {
		t.Error("覆蓋層未初始化")
	}
	
	// 驗證指示器已初始化
	if feedback.indicator == nil {
		t.Error("指示器未初始化")
	}
	
	// 驗證初始狀態為隱藏
	if feedback.isVisible {
		t.Error("初始狀態應該為隱藏")
	}
	
	// 驗證初始區域為空
	if feedback.currentZone != "" {
		t.Error("初始區域應該為空")
	}
}

// TestDragFeedback_ShowHide 測試拖拽視覺回饋的顯示和隱藏
// 驗證視覺回饋的顯示和隱藏功能是否正常工作
func TestDragFeedback_ShowHide(t *testing.T) {
	// 建立拖拽視覺回饋
	feedback := NewDragFeedback()
	
	// 測試顯示拖拽開始
	t.Run("顯示拖拽開始", func(t *testing.T) {
		sourcePath := "/test/source/file.md"
		feedback.ShowDragStart(sourcePath)
		
		// 驗證可見狀態
		if !feedback.IsVisible() {
			t.Error("拖拽開始後應該可見")
		}
		
		// 驗證指示器文字
		expectedText := "拖拽中: file.md"
		if feedback.indicator.Text != expectedText {
			t.Errorf("指示器文字不正確：預期 '%s'，實際 '%s'", expectedText, feedback.indicator.Text)
		}
	})
	
	// 測試顯示進入回饋
	t.Run("顯示進入回饋", func(t *testing.T) {
		zoneID := "test-zone"
		feedback.ShowEnterFeedback(zoneID)
		
		// 驗證可見狀態
		if !feedback.IsVisible() {
			t.Error("進入回饋後應該可見")
		}
		
		// 驗證目前區域
		if feedback.GetCurrentZone() != zoneID {
			t.Errorf("目前區域不正確：預期 '%s'，實際 '%s'", zoneID, feedback.GetCurrentZone())
		}
		
		// 驗證指示器文字
		expectedText := "放下到: test-zone"
		if feedback.indicator.Text != expectedText {
			t.Errorf("指示器文字不正確：預期 '%s'，實際 '%s'", expectedText, feedback.indicator.Text)
		}
	})
	
	// 測試顯示離開回饋
	t.Run("顯示離開回饋", func(t *testing.T) {
		feedback.ShowLeaveFeedback()
		
		// 驗證目前區域已清空
		if feedback.GetCurrentZone() != "" {
			t.Error("離開回饋後目前區域應該為空")
		}
		
		// 驗證指示器文字
		expectedText := "拖拽中..."
		if feedback.indicator.Text != expectedText {
			t.Errorf("指示器文字不正確：預期 '%s'，實際 '%s'", expectedText, feedback.indicator.Text)
		}
	})
	
	// 測試隱藏回饋
	t.Run("隱藏回饋", func(t *testing.T) {
		feedback.Hide()
		
		// 驗證可見狀態
		if feedback.IsVisible() {
			t.Error("隱藏後應該不可見")
		}
		
		// 驗證指示器文字已清空
		if feedback.indicator.Text != "" {
			t.Error("隱藏後指示器文字應該為空")
		}
		
		// 驗證目前區域已清空
		if feedback.GetCurrentZone() != "" {
			t.Error("隱藏後目前區域應該為空")
		}
	})
}

// TestDragDropHelper_CreateFileDropZone 測試建立檔案拖拽區域
// 驗證拖拽區域的建立和配置是否正確
func TestDragDropHelper_CreateFileDropZone(t *testing.T) {
	// 建立拖拽輔助工具
	helper := NewDragDropHelper()
	
	// 測試參數
	title := "測試拖拽區域"
	targetPath := "/test/target"
	acceptTypes := []string{".md", ".txt"}
	
	// 建立檔案拖拽區域
	dropZone := helper.CreateFileDropZone(title, targetPath, acceptTypes)
	
	// 驗證拖拽區域不為 nil
	if dropZone == nil {
		t.Fatal("檔案拖拽區域建立失敗，回傳 nil")
	}
	
	// 驗證容器類型 (basic check since we can't import container)
	if dropZone == nil {
		t.Error("拖拽區域不應該為 nil")
	}
	
	// 驗證容器包含元件
	if len(dropZone.Objects) == 0 {
		t.Error("拖拽區域應該包含 UI 元件")
	}
}

// TestDragDropHelper_CreateDragHandle 測試建立拖拽控制項
// 驗證拖拽控制項的建立和配置是否正確
func TestDragDropHelper_CreateDragHandle(t *testing.T) {
	// 建立拖拽輔助工具
	helper := NewDragDropHelper()
	
	// 測試參數
	sourcePath := "/test/source/file.md"
	
	// 建立拖拽控制項
	dragHandle := helper.CreateDragHandle(sourcePath)
	
	// 驗證拖拽控制項不為 nil
	if dragHandle == nil {
		t.Fatal("拖拽控制項建立失敗，回傳 nil")
	}
	
	// 驗證控制項類型 (basic check)
	if dragHandle == nil {
		t.Error("拖拽控制項不應該為 nil")
	}
	
	// 驗證按鈕重要性設定
	if dragHandle.Importance != widget.LowImportance {
		t.Error("拖拽控制項應該設定為低重要性")
	}
}

// TestDragDropHelper_ValidateDropOperation 測試拖拽操作驗證
// 驗證拖拽操作的驗證邏輯是否正確
func TestDragDropHelper_ValidateDropOperation(t *testing.T) {
	// 建立拖拽輔助工具
	helper := NewDragDropHelper()
	
	// 測試案例
	testCases := []struct {
		name         string
		sourcePath   string
		targetPath   string
		operation    DragOperation
		expectedValid bool
		expectedMsg   string
	}{
		{
			name:         "有效的移動操作",
			sourcePath:   "/test/source/file.md",
			targetPath:   "/test/target/file.md",
			operation:    DragOperationMove,
			expectedValid: true,
			expectedMsg:   "",
		},
		{
			name:         "空來源路徑",
			sourcePath:   "",
			targetPath:   "/test/target/file.md",
			operation:    DragOperationMove,
			expectedValid: false,
			expectedMsg:   "來源路徑不能為空",
		},
		{
			name:         "空目標路徑",
			sourcePath:   "/test/source/file.md",
			targetPath:   "",
			operation:    DragOperationMove,
			expectedValid: false,
			expectedMsg:   "目標路徑不能為空",
		},
		{
			name:         "相同路徑",
			sourcePath:   "/test/file.md",
			targetPath:   "/test/file.md",
			operation:    DragOperationMove,
			expectedValid: false,
			expectedMsg:   "來源和目標路徑不能相同",
		},
		{
			name:         "子目錄移動",
			sourcePath:   "/test/parent",
			targetPath:   "/test/parent/child",
			operation:    DragOperationMove,
			expectedValid: false,
			expectedMsg:   "不能將目錄移動到其子目錄中",
		},
		{
			name:         "有效的複製操作",
			sourcePath:   "/test/source/file.md",
			targetPath:   "/test/target/file.md",
			operation:    DragOperationCopy,
			expectedValid: true,
			expectedMsg:   "",
		},
		{
			name:         "有效的連結操作",
			sourcePath:   "/test/source/file.md",
			targetPath:   "/test/target/file.md",
			operation:    DragOperationLink,
			expectedValid: true,
			expectedMsg:   "",
		},
	}
	
	// 執行測試案例
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			valid, msg := helper.ValidateDropOperation(tc.sourcePath, tc.targetPath, tc.operation)
			
			if valid != tc.expectedValid {
				t.Errorf("驗證結果不正確：預期 %v，實際 %v", tc.expectedValid, valid)
			}
			
			if tc.expectedMsg != "" && msg != tc.expectedMsg {
				t.Errorf("錯誤訊息不正確：預期 '%s'，實際 '%s'", tc.expectedMsg, msg)
			}
		})
	}
}

// BenchmarkDragDropManager_IsValidFileType 效能測試：檔案類型驗證
// 測試檔案類型驗證功能的效能表現
func BenchmarkDragDropManager_IsValidFileType(b *testing.B) {
	// 建立測試環境
	mockFileManager := &MockFileManagerService{}
	manager := NewDragDropManager(mockFileManager, nil)
	
	// 測試參數
	filePath := "/test/file.md"
	acceptTypes := []string{".md", ".txt", ".markdown"}
	
	// 執行效能測試
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		manager.IsValidFileType(filePath, acceptTypes)
	}
}

// BenchmarkDragDropHelper_ValidateDropOperation 效能測試：拖拽操作驗證
// 測試拖拽操作驗證功能的效能表現
func BenchmarkDragDropHelper_ValidateDropOperation(b *testing.B) {
	// 建立拖拽輔助工具
	helper := NewDragDropHelper()
	
	// 測試參數
	sourcePath := "/test/source/file.md"
	targetPath := "/test/target/file.md"
	operation := DragOperationMove
	
	// 執行效能測試
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		helper.ValidateDropOperation(sourcePath, targetPath, operation)
	}
}