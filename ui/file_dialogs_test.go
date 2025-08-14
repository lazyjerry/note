// Package ui 包含檔案對話框功能的測試
package ui

import (
	"testing"                  // Go 標準測試套件
)

// TestNewFileDialogManager 測試檔案對話框管理器的建立
// 驗證管理器實例是否正確建立並設定父視窗
func TestNewFileDialogManager(t *testing.T) {
	// 建立檔案對話框管理器（使用 nil 作為測試）
	manager := NewFileDialogManager(nil)
	
	// 驗證管理器不為 nil
	if manager == nil {
		t.Fatal("檔案對話框管理器建立失敗，回傳 nil")
	}
	
	// 驗證父視窗設定正確
	if manager.parent != nil {
		t.Error("父視窗設定不正確")
	}
}

// TestFileDialogManager_isValidFileType 測試檔案類型驗證功能
// 驗證各種檔案副檔名的驗證結果是否正確
func TestFileDialogManager_isValidFileType(t *testing.T) {
	// 建立測試管理器（不需要實際的視窗）
	manager := NewFileDialogManager(nil)
	
	// 測試案例：檔案路徑和預期的驗證結果
	testCases := []struct {
		name     string // 測試案例名稱
		filePath string // 檔案路徑
		expected bool   // 預期的驗證結果
	}{
		{
			name:     "Markdown 檔案 (.md)",
			filePath: "/path/to/file.md",
			expected: true,
		},
		{
			name:     "Markdown 檔案 (.markdown)",
			filePath: "/path/to/file.markdown",
			expected: true,
		},
		{
			name:     "文字檔案 (.txt)",
			filePath: "/path/to/file.txt",
			expected: true,
		},
		{
			name:     "大寫副檔名 (.MD)",
			filePath: "/path/to/file.MD",
			expected: true,
		},
		{
			name:     "混合大小寫副檔名 (.Md)",
			filePath: "/path/to/file.Md",
			expected: true,
		},
		{
			name:     "不支援的檔案類型 (.doc)",
			filePath: "/path/to/file.doc",
			expected: false,
		},
		{
			name:     "不支援的檔案類型 (.pdf)",
			filePath: "/path/to/file.pdf",
			expected: false,
		},
		{
			name:     "沒有副檔名的檔案",
			filePath: "/path/to/file",
			expected: false,
		},
		{
			name:     "空檔案路徑",
			filePath: "",
			expected: false,
		},
	}
	
	// 執行所有測試案例
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := manager.isValidFileType(tc.filePath)
			if result != tc.expected {
				t.Errorf("檔案類型驗證失敗：檔案 %s，預期 %v，實際 %v", 
					tc.filePath, tc.expected, result)
			}
		})
	}
}

// TestFileDialogManager_isValidFileTypeCustom 測試自訂檔案類型驗證功能
// 驗證自訂檔案類型列表的驗證結果是否正確
func TestFileDialogManager_isValidFileTypeCustom(t *testing.T) {
	// 建立測試管理器（不需要實際的視窗）
	manager := NewFileDialogManager(nil)
	
	// 測試案例：檔案路徑、支援類型列表和預期結果
	testCases := []struct {
		name           string   // 測試案例名稱
		filePath       string   // 檔案路徑
		supportedTypes []string // 支援的檔案類型列表
		expected       bool     // 預期的驗證結果
	}{
		{
			name:           "支援的 Markdown 檔案",
			filePath:       "/path/to/file.md",
			supportedTypes: []string{".md", ".txt"},
			expected:       true,
		},
		{
			name:           "支援的文字檔案",
			filePath:       "/path/to/file.txt",
			supportedTypes: []string{".md", ".txt"},
			expected:       true,
		},
		{
			name:           "不在支援列表中的檔案",
			filePath:       "/path/to/file.doc",
			supportedTypes: []string{".md", ".txt"},
			expected:       false,
		},
		{
			name:           "空的支援類型列表",
			filePath:       "/path/to/file.md",
			supportedTypes: []string{},
			expected:       false,
		},
		{
			name:           "大小寫不敏感測試",
			filePath:       "/path/to/file.MD",
			supportedTypes: []string{".md", ".txt"},
			expected:       true,
		},
		{
			name:           "支援類型列表包含大寫",
			filePath:       "/path/to/file.md",
			supportedTypes: []string{".MD", ".TXT"},
			expected:       true,
		},
	}
	
	// 執行所有測試案例
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := manager.isValidFileTypeCustom(tc.filePath, tc.supportedTypes)
			if result != tc.expected {
				t.Errorf("自訂檔案類型驗證失敗：檔案 %s，支援類型 %v，預期 %v，實際 %v", 
					tc.filePath, tc.supportedTypes, tc.expected, result)
			}
		})
	}
}

// TestFileTypeError 測試檔案類型錯誤結構
// 驗證錯誤結構的建立和方法是否正常工作
func TestFileTypeError(t *testing.T) {
	// 建立檔案類型錯誤實例
	testPath := "/path/to/invalid/file.doc"
	testMessage := "不支援的檔案類型"
	
	err := &FileTypeError{
		Path:    testPath,
		Message: testMessage,
	}
	
	// 測試 Error() 方法
	if err.Error() != testMessage {
		t.Errorf("Error() 方法回傳錯誤：預期 %s，實際 %s", testMessage, err.Error())
	}
	
	// 測試 GetPath() 方法
	if err.GetPath() != testPath {
		t.Errorf("GetPath() 方法回傳錯誤：預期 %s，實際 %s", testPath, err.GetPath())
	}
}

// TestFileDialogConfig 測試檔案對話框配置結構
// 驗證配置結構的各個欄位是否正確設定
func TestFileDialogConfig(t *testing.T) {
	// 建立測試配置
	config := FileDialogConfig{
		Title:           "測試對話框",
		DefaultName:     "test.md",
		DefaultLocation: "/home/user/documents",
		FileTypes:       []string{".md", ".txt"},
		AllowMultiple:   true,
	}
	
	// 驗證各個欄位的設定
	if config.Title != "測試對話框" {
		t.Errorf("標題設定錯誤：預期 '測試對話框'，實際 '%s'", config.Title)
	}
	
	if config.DefaultName != "test.md" {
		t.Errorf("預設名稱設定錯誤：預期 'test.md'，實際 '%s'", config.DefaultName)
	}
	
	if config.DefaultLocation != "/home/user/documents" {
		t.Errorf("預設位置設定錯誤：預期 '/home/user/documents'，實際 '%s'", config.DefaultLocation)
	}
	
	if len(config.FileTypes) != 2 {
		t.Errorf("檔案類型數量錯誤：預期 2，實際 %d", len(config.FileTypes))
	}
	
	if !config.AllowMultiple {
		t.Error("多選設定錯誤：預期 true，實際 false")
	}
}

// Note: Integration tests that require actual Fyne windows are omitted
// to avoid threading issues in test environment. These would be tested
// manually or in a separate integration test suite.

// BenchmarkFileTypeValidation 效能測試：檔案類型驗證
// 測試檔案類型驗證功能的效能表現
func BenchmarkFileTypeValidation(b *testing.B) {
	// 建立測試管理器（不需要實際的視窗）
	manager := NewFileDialogManager(nil)
	
	// 測試檔案路徑
	testPath := "/path/to/test/file.md"
	
	// 執行效能測試
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		manager.isValidFileType(testPath)
	}
}

// BenchmarkCustomFileTypeValidation 效能測試：自訂檔案類型驗證
// 測試自訂檔案類型驗證功能的效能表現
func BenchmarkCustomFileTypeValidation(b *testing.B) {
	// 建立測試管理器（不需要實際的視窗）
	manager := NewFileDialogManager(nil)
	
	// 測試檔案路徑和支援類型
	testPath := "/path/to/test/file.md"
	supportedTypes := []string{".md", ".txt", ".markdown"}
	
	// 執行效能測試
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		manager.isValidFileTypeCustom(testPath, supportedTypes)
	}
}