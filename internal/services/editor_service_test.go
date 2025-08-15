// Package services 包含編輯器服務的單元測試
// 測試 EditorService 的所有核心功能，包含筆記建立、開啟、保存、更新和 Markdown 預覽
package services

import (
	"fmt"
	"mac-notebook-app/internal/models"
	"strings"
	"testing"
	"time"
)

// mockFileRepository 模擬檔案存取介面，用於測試
// 實作 FileRepository 介面的所有方法，提供可控制的測試環境
type mockFileRepository struct {
	files       map[string][]byte // 模擬檔案系統的檔案內容
	fileExists  map[string]bool   // 模擬檔案存在狀態
	writeError  error             // 模擬寫入錯誤
	readError   error             // 模擬讀取錯誤
}

// newMockFileRepository 建立新的模擬檔案存取實例
// 回傳：初始化的模擬檔案存取介面
func newMockFileRepository() *mockFileRepository {
	return &mockFileRepository{
		files:      make(map[string][]byte),
		fileExists: make(map[string]bool),
	}
}

// mockEncryptionService 模擬加密服務介面，用於測試
type mockEncryptionService struct{}

func (m *mockEncryptionService) EncryptContent(content, password string, algorithm string) ([]byte, error) {
	return []byte("encrypted:" + content), nil
}

func (m *mockEncryptionService) DecryptContent(encryptedData []byte, password string, algorithm string) (string, error) {
	data := string(encryptedData)
	if strings.HasPrefix(data, "encrypted:") {
		return strings.TrimPrefix(data, "encrypted:"), nil
	}
	return "", fmt.Errorf("invalid encrypted data")
}

func (m *mockEncryptionService) SetupBiometricAuth(noteID string) error {
	return nil
}

func (m *mockEncryptionService) AuthenticateWithBiometric(noteID string) (bool, error) {
	return false, fmt.Errorf("biometric not available in test")
}

func (m *mockEncryptionService) ValidatePassword(password string) bool {
	return len(password) >= 8
}

// mockPasswordService 模擬密碼服務介面，用於測試
type mockPasswordService struct{}

func (m *mockPasswordService) HashPassword(password string) (*PasswordHash, error) {
	return &PasswordHash{
		Salt:      "mock-salt",
		Hash:      "hashed:" + password,
		Algorithm: "pbkdf2-sha256",
		Rounds:    100000,
		CreatedAt: time.Now(),
	}, nil
}

func (m *mockPasswordService) VerifyPassword(password string, hash *PasswordHash) (bool, error) {
	return hash.Hash == "hashed:"+password, nil
}

func (m *mockPasswordService) CheckPasswordStrength(password string) (PasswordStrength, []string) {
	if len(password) >= 8 {
		return PasswordStrong, []string{}
	}
	return PasswordWeak, []string{"密碼太短"}
}

func (m *mockPasswordService) RecordFailedAttempt(identifier string) error {
	return nil
}

func (m *mockPasswordService) IsLocked(identifier string) (bool, time.Duration) {
	return false, 0
}

func (m *mockPasswordService) ResetRetryCount(identifier string) {
}

func (m *mockPasswordService) GetRetryInfo(identifier string) *RetryInfo {
	return &RetryInfo{}
}

// mockBiometricService 模擬生物識別服務介面，用於測試
type mockBiometricService struct{}

func (m *mockBiometricService) IsAvailable() (bool, BiometricType) {
	return false, BiometricTypeNone
}

func (m *mockBiometricService) Authenticate(reason string) *BiometricResult {
	return &BiometricResult{
		Success: false,
		Error:   fmt.Errorf("not available in test"),
	}
}

func (m *mockBiometricService) SetupForNote(noteID string) error {
	return nil
}

func (m *mockBiometricService) AuthenticateForNote(noteID, reason string) *BiometricResult {
	return &BiometricResult{
		Success: false,
		Error:   fmt.Errorf("not available in test"),
	}
}

func (m *mockBiometricService) RemoveForNote(noteID string) error {
	return nil
}

func (m *mockBiometricService) IsEnabledForNote(noteID string) bool {
	return false
}

// createTestEditorService 建立用於測試的編輯器服務實例
// 回傳：編輯器服務實例和模擬檔案存取介面
func createTestEditorService() (EditorService, *mockFileRepository) {
	mockRepo := newMockFileRepository()
	mockEncryption := &mockEncryptionService{}
	mockPassword := &mockPasswordService{}
	mockBiometric := &mockBiometricService{}
	
	service := NewEditorService(mockRepo, mockEncryption, mockPassword, mockBiometric, nil)
	return service, mockRepo
}

// WriteFile 模擬檔案寫入操作
func (m *mockFileRepository) WriteFile(path string, data []byte) error {
	if m.writeError != nil {
		return m.writeError
	}
	m.files[path] = data
	m.fileExists[path] = true
	return nil
}

// ReadFile 模擬檔案讀取操作
func (m *mockFileRepository) ReadFile(path string) ([]byte, error) {
	if m.readError != nil {
		return nil, m.readError
	}
	if data, exists := m.files[path]; exists {
		return data, nil
	}
	return nil, fmt.Errorf("file not found")
}

// FileExists 模擬檔案存在性檢查
func (m *mockFileRepository) FileExists(path string) bool {
	return m.fileExists[path]
}

// DeleteFile 模擬檔案刪除操作
func (m *mockFileRepository) DeleteFile(path string) error {
	delete(m.files, path)
	delete(m.fileExists, path)
	return nil
}

// CreateDirectory 模擬目錄建立操作
func (m *mockFileRepository) CreateDirectory(path string) error {
	return nil
}

// ListDirectory 模擬目錄列表操作
func (m *mockFileRepository) ListDirectory(path string) ([]*models.FileInfo, error) {
	return nil, nil
}

// WalkDirectory 模擬目錄遍歷操作
func (m *mockFileRepository) WalkDirectory(path string, walkFunc func(*models.FileInfo) error) error {
	return nil
}

// TestNewEditorService 測試編輯器服務的建立
// 驗證服務實例是否正確初始化
func TestNewEditorService(t *testing.T) {
	// 建立測試編輯器服務
	service, _ := createTestEditorService()
	
	// 驗證服務不為空
	if service == nil {
		t.Fatal("編輯器服務不應該為空")
	}
	
	// 驗證服務實例類型
	editorSvc, ok := service.(*editorService)
	if !ok {
		t.Fatal("服務應該是 editorService 類型")
	}
	
	// 驗證內部狀態初始化
	if editorSvc.fileRepo == nil {
		t.Error("檔案存取介面不應該為空")
	}
	
	if editorSvc.encryptionSvc == nil {
		t.Error("加密服務介面不應該為空")
	}
	
	if editorSvc.passwordSvc == nil {
		t.Error("密碼服務介面不應該為空")
	}
	
	if editorSvc.biometricSvc == nil {
		t.Error("生物識別服務介面不應該為空")
	}
	
	if editorSvc.markdown == nil {
		t.Error("Markdown 解析器不應該為空")
	}
	
	if editorSvc.activeNotes == nil {
		t.Error("活躍筆記快取不應該為空")
	}
}

// TestCreateNote 測試筆記建立功能
// 驗證新筆記的建立和屬性設定
func TestCreateNote(t *testing.T) {
	// 建立測試環境
	service, _ := createTestEditorService()
	
	// 測試資料
	title := "測試筆記"
	content := "這是測試內容"
	
	// 執行建立筆記操作
	note, err := service.CreateNote(title, content)
	
	// 驗證結果
	if err != nil {
		t.Fatalf("建立筆記失敗: %v", err)
	}
	
	if note == nil {
		t.Fatal("筆記不應該為空")
	}
	
	// 驗證筆記屬性
	if note.Title != title {
		t.Errorf("筆記標題不正確，期望: %s，實際: %s", title, note.Title)
	}
	
	if note.Content != content {
		t.Errorf("筆記內容不正確，期望: %s，實際: %s", content, note.Content)
	}
	
	if note.ID == "" {
		t.Error("筆記 ID 不應該為空")
	}
	
	if note.IsEncrypted {
		t.Error("新筆記預設不應該加密")
	}
	
	if note.CreatedAt.IsZero() {
		t.Error("建立時間不應該為零值")
	}
	
	if note.UpdatedAt.IsZero() {
		t.Error("更新時間不應該為零值")
	}
	
	// 驗證筆記是否加入活躍快取
	editorSvc := service.(*editorService)
	if _, exists := editorSvc.activeNotes[note.ID]; !exists {
		t.Error("筆記應該加入活躍快取")
	}
}

// TestOpenNote 測試筆記開啟功能
// 驗證從檔案系統開啟現有筆記
func TestOpenNote(t *testing.T) {
	// 建立測試環境
	service, mockRepo := createTestEditorService()
	
	// 準備測試檔案
	filePath := "test-note.md"
	fileContent := "# 測試標題\n\n這是測試內容"
	mockRepo.files[filePath] = []byte(fileContent)
	mockRepo.fileExists[filePath] = true
	
	// 執行開啟筆記操作
	note, err := service.OpenNote(filePath)
	
	// 驗證結果
	if err != nil {
		t.Fatalf("開啟筆記失敗: %v", err)
	}
	
	if note == nil {
		t.Fatal("筆記不應該為空")
	}
	
	// 驗證筆記屬性
	if note.Title != "test-note" {
		t.Errorf("筆記標題不正確，期望: test-note，實際: %s", note.Title)
	}
	
	if note.Content != fileContent {
		t.Errorf("筆記內容不正確，期望: %s，實際: %s", fileContent, note.Content)
	}
	
	if note.FilePath != filePath {
		t.Errorf("檔案路徑不正確，期望: %s，實際: %s", filePath, note.FilePath)
	}
	
	if note.IsEncrypted {
		t.Error("普通 Markdown 檔案不應該標記為加密")
	}
}

// TestOpenEncryptedNote 測試開啟加密筆記
// 驗證加密檔案的正確識別
func TestOpenEncryptedNote(t *testing.T) {
	// 建立測試環境
	service, mockRepo := createTestEditorService()
	
	// 準備加密測試檔案（使用模擬加密格式）
	filePath := "encrypted-note.md.enc"
	fileContent := "encrypted:加密的筆記內容" // 使用模擬加密格式
	mockRepo.files[filePath] = []byte(fileContent)
	mockRepo.fileExists[filePath] = true
	
	// 由於加密檔案需要密碼驗證，我們期望會收到錯誤
	_, err := service.OpenNote(filePath)
	
	// 驗證應該收到需要密碼驗證的錯誤
	if err == nil {
		t.Fatal("開啟加密檔案應該要求密碼驗證")
	}
	
	if !strings.Contains(err.Error(), "需要密碼驗證") {
		t.Errorf("錯誤訊息應該包含密碼驗證要求，實際: %s", err.Error())
	}
}

// TestOpenNonExistentNote 測試開啟不存在的筆記
// 驗證錯誤處理機制
func TestOpenNonExistentNote(t *testing.T) {
	// 建立測試環境
	service, _ := createTestEditorService()
	
	// 嘗試開啟不存在的檔案
	_, err := service.OpenNote("non-existent.md")
	
	// 驗證應該回傳錯誤
	if err == nil {
		t.Fatal("開啟不存在的檔案應該回傳錯誤")
	}
	
	if !strings.Contains(err.Error(), "檔案不存在") {
		t.Errorf("錯誤訊息應該包含檔案不存在，實際: %s", err.Error())
	}
}

// TestSaveNote 測試筆記保存功能
// 驗證筆記保存到檔案系統
func TestSaveNote(t *testing.T) {
	// 建立測試環境
	service, mockRepo := createTestEditorService()
	
	// 建立測試筆記
	note, err := service.CreateNote("測試筆記", "測試內容")
	if err != nil {
		t.Fatalf("建立筆記失敗: %v", err)
	}
	
	// 記錄保存前的時間
	beforeSave := time.Now()
	
	// 執行保存操作
	err = service.SaveNote(note)
	if err != nil {
		t.Fatalf("保存筆記失敗: %v", err)
	}
	
	// 驗證檔案是否寫入
	expectedPath := "測試筆記.md"
	if data, exists := mockRepo.files[expectedPath]; !exists {
		t.Error("筆記檔案應該被寫入")
	} else if string(data) != note.Content {
		t.Errorf("檔案內容不正確，期望: %s，實際: %s", note.Content, string(data))
	}
	
	// 驗證時間戳更新
	if note.LastSaved.Before(beforeSave) {
		t.Error("最後保存時間應該更新")
	}
	
	if note.UpdatedAt.Before(beforeSave) {
		t.Error("更新時間應該更新")
	}
}

// TestSaveNoteWithCustomPath 測試使用自訂路徑保存筆記
func TestSaveNoteWithCustomPath(t *testing.T) {
	// 建立測試環境
	service, mockRepo := createTestEditorService()
	
	// 建立測試筆記並設定自訂路徑
	note, err := service.CreateNote("測試筆記", "測試內容")
	if err != nil {
		t.Fatalf("建立筆記失敗: %v", err)
	}
	
	customPath := "custom/path/my-note.md"
	note.FilePath = customPath
	
	// 執行保存操作
	err = service.SaveNote(note)
	if err != nil {
		t.Fatalf("保存筆記失敗: %v", err)
	}
	
	// 驗證檔案是否寫入到正確路徑
	if _, exists := mockRepo.files[customPath]; !exists {
		t.Errorf("筆記應該保存到自訂路徑: %s", customPath)
	}
}

// TestSaveNullNote 測試保存空筆記的錯誤處理
func TestSaveNullNote(t *testing.T) {
	// 建立測試環境
	service, _ := createTestEditorService()
	
	// 嘗試保存空筆記
	err := service.SaveNote(nil)
	
	// 驗證應該回傳錯誤
	if err == nil {
		t.Fatal("保存空筆記應該回傳錯誤")
	}
	
	if !strings.Contains(err.Error(), "筆記實例不能為空") {
		t.Errorf("錯誤訊息不正確，實際: %s", err.Error())
	}
}

// TestUpdateContent 測試筆記內容更新功能
func TestUpdateContent(t *testing.T) {
	// 建立測試環境
	service, _ := createTestEditorService()
	
	// 建立測試筆記
	note, err := service.CreateNote("測試筆記", "原始內容")
	if err != nil {
		t.Fatalf("建立筆記失敗: %v", err)
	}
	
	// 記錄更新前的時間
	beforeUpdate := time.Now()
	
	// 執行內容更新
	newContent := "更新後的內容"
	err = service.UpdateContent(note.ID, newContent)
	if err != nil {
		t.Fatalf("更新內容失敗: %v", err)
	}
	
	// 驗證內容是否更新
	if note.Content != newContent {
		t.Errorf("內容未正確更新，期望: %s，實際: %s", newContent, note.Content)
	}
	
	// 驗證更新時間
	if note.UpdatedAt.Before(beforeUpdate) {
		t.Error("更新時間應該更新")
	}
}

// TestUpdateNonExistentNote 測試更新不存在筆記的錯誤處理
func TestUpdateNonExistentNote(t *testing.T) {
	// 建立測試環境
	service, _ := createTestEditorService()
	
	// 嘗試更新不存在的筆記
	err := service.UpdateContent("non-existent-id", "新內容")
	
	// 驗證應該回傳錯誤
	if err == nil {
		t.Fatal("更新不存在的筆記應該回傳錯誤")
	}
	
	if !strings.Contains(err.Error(), "找不到指定的筆記") {
		t.Errorf("錯誤訊息不正確，實際: %s", err.Error())
	}
}

// TestPreviewMarkdown 測試 Markdown 預覽功能
// 驗證各種 Markdown 語法的正確轉換
func TestPreviewMarkdown(t *testing.T) {
	// 建立測試環境
	service, _ := createTestEditorService()
	
	// 測試案例
	testCases := []struct {
		name     string
		markdown string
		expected string
	}{
		{
			name:     "標題轉換",
			markdown: "# 主標題\n## 副標題",
			expected: "<h1",
		},
		{
			name:     "粗體和斜體",
			markdown: "**粗體** 和 *斜體*",
			expected: "<strong>粗體</strong> 和 <em>斜體</em>",
		},
		{
			name:     "程式碼區塊",
			markdown: "```go\nfunc main() {}\n```",
			expected: "<pre><code class=\"language-go\">func main() {}\n</code></pre>",
		},
		{
			name:     "連結",
			markdown: "[Google](https://google.com)",
			expected: "<a href=\"https://google.com\">Google</a>",
		},
		{
			name:     "列表",
			markdown: "- 項目 1\n- 項目 2",
			expected: "<ul>\n<li>項目 1</li>\n<li>項目 2</li>\n</ul>",
		},
	}
	
	// 執行測試案例
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := service.PreviewMarkdown(tc.markdown)
			
			// 驗證結果包含期望的 HTML 元素
			if !strings.Contains(result, tc.expected) {
				t.Errorf("Markdown 轉換結果不正確\n期望包含: %s\n實際結果: %s", tc.expected, result)
			}
			
			// 對於標題轉換測試，額外檢查是否包含標題內容
			if tc.name == "標題轉換" {
				if !strings.Contains(result, "主標題") || !strings.Contains(result, "副標題") {
					t.Errorf("標題內容應該被正確轉換")
				}
			}
		})
	}
}

// TestPreviewMarkdownTable 測試表格 Markdown 預覽
func TestPreviewMarkdownTable(t *testing.T) {
	// 建立測試環境
	service, _ := createTestEditorService()
	
	// 測試表格 Markdown
	tableMarkdown := `| 欄位1 | 欄位2 |
|-------|-------|
| 值1   | 值2   |`
	
	result := service.PreviewMarkdown(tableMarkdown)
	
	// 驗證表格元素
	expectedElements := []string{"<table>", "<thead>", "<tbody>", "<tr>", "<th>", "<td>"}
	for _, element := range expectedElements {
		if !strings.Contains(result, element) {
			t.Errorf("表格 HTML 應該包含 %s，實際結果: %s", element, result)
		}
	}
}

// TestPreviewMarkdownTaskList 測試任務列表 Markdown 預覽
func TestPreviewMarkdownTaskList(t *testing.T) {
	// 建立測試環境
	service, _ := createTestEditorService()
	
	// 測試任務列表 Markdown
	taskListMarkdown := `- [x] 已完成任務
- [ ] 未完成任務`
	
	result := service.PreviewMarkdown(taskListMarkdown)
	
	// 驗證任務列表元素
	if !strings.Contains(result, "checkbox") {
		t.Errorf("任務列表應該包含 checkbox，實際結果: %s", result)
	}
}

// TestSanitizeFileName 測試檔案名稱清理功能
func TestSanitizeFileName(t *testing.T) {
	// 建立測試環境
	service, _ := createTestEditorService()
	editorSvc := service.(*editorService)
	
	// 測試案例
	testCases := []struct {
		input    string
		expected string
	}{
		{"正常檔案名", "正常檔案名"},
		{"包含/斜線", "包含_斜線"},
		{"包含\\反斜線", "包含_反斜線"},
		{"包含:冒號", "包含_冒號"},
		{"包含*星號", "包含_星號"},
		{"包含?問號", "包含_問號"},
		{"包含\"引號", "包含_引號"},
		{"包含<小於", "包含_小於"},
		{"包含>大於", "包含_大於"},
		{"包含|管道", "包含_管道"},
		{"  前後空白  ", "前後空白"},
	}
	
	// 執行測試案例
	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := editorSvc.sanitizeFileName(tc.input)
			if result != tc.expected {
				t.Errorf("檔案名稱清理不正確，輸入: %s，期望: %s，實際: %s", tc.input, tc.expected, result)
			}
		})
	}
}

// TestGetActiveNote 測試取得活躍筆記功能
func TestGetActiveNote(t *testing.T) {
	// 建立測試環境
	service, _ := createTestEditorService()
	editorSvc := service.(*editorService)
	
	// 建立測試筆記
	note, err := service.CreateNote("測試筆記", "測試內容")
	if err != nil {
		t.Fatalf("建立筆記失敗: %v", err)
	}
	
	// 測試取得存在的筆記
	retrievedNote, exists := editorSvc.GetActiveNote(note.ID)
	if !exists {
		t.Error("應該能找到活躍筆記")
	}
	if retrievedNote.ID != note.ID {
		t.Error("取得的筆記 ID 不正確")
	}
	
	// 測試取得不存在的筆記
	_, exists = editorSvc.GetActiveNote("non-existent-id")
	if exists {
		t.Error("不應該找到不存在的筆記")
	}
}

// TestCloseNote 測試關閉筆記功能
func TestCloseNote(t *testing.T) {
	// 建立測試環境
	service, _ := createTestEditorService()
	editorSvc := service.(*editorService)
	
	// 建立測試筆記
	note, err := service.CreateNote("測試筆記", "測試內容")
	if err != nil {
		t.Fatalf("建立筆記失敗: %v", err)
	}
	
	// 驗證筆記在活躍快取中
	if _, exists := editorSvc.activeNotes[note.ID]; !exists {
		t.Error("筆記應該在活躍快取中")
	}
	
	// 關閉筆記
	editorSvc.CloseNote(note.ID)
	
	// 驗證筆記已從活躍快取中移除
	if _, exists := editorSvc.activeNotes[note.ID]; exists {
		t.Error("筆記應該從活躍快取中移除")
	}
}

// TestGetActiveNotes 測試取得所有活躍筆記功能
func TestGetActiveNotes(t *testing.T) {
	// 建立測試環境
	service, _ := createTestEditorService()
	editorSvc := service.(*editorService)
	
	// 建立多個測試筆記
	note1, _ := service.CreateNote("筆記1", "內容1")
	note2, _ := service.CreateNote("筆記2", "內容2")
	
	// 取得所有活躍筆記
	activeNotes := editorSvc.GetActiveNotes()
	
	// 驗證筆記數量
	if len(activeNotes) != 2 {
		t.Errorf("活躍筆記數量不正確，期望: 2，實際: %d", len(activeNotes))
	}
	
	// 驗證筆記內容
	if _, exists := activeNotes[note1.ID]; !exists {
		t.Error("筆記1 應該在活躍筆記列表中")
	}
	
	if _, exists := activeNotes[note2.ID]; !exists {
		t.Error("筆記2 應該在活躍筆記列表中")
	}
	
	// 驗證回傳的是副本（修改不會影響原始快取）
	delete(activeNotes, note1.ID)
	if _, exists := editorSvc.activeNotes[note1.ID]; !exists {
		t.Error("修改回傳的活躍筆記列表不應該影響原始快取")
	}
}

// TestEnableEncryption 測試啟用加密功能
func TestEnableEncryption(t *testing.T) {
	// 建立測試環境
	service, _ := createTestEditorService()
	
	// 建立測試筆記
	note, err := service.CreateNote("測試筆記", "測試內容")
	if err != nil {
		t.Fatalf("建立筆記失敗: %v", err)
	}
	
	// 啟用加密
	err = service.(*editorService).EnableEncryption(note.ID, "TestPassword123!", "aes256", false)
	if err != nil {
		t.Fatalf("啟用加密失敗: %v", err)
	}
	
	// 驗證加密狀態
	if !note.IsEncrypted {
		t.Error("筆記應該標記為已加密")
	}
	
	if note.EncryptionType != "aes256" {
		t.Errorf("加密類型不正確，期望: aes256，實際: %s", note.EncryptionType)
	}
}

// TestDisableEncryption 測試停用加密功能
func TestDisableEncryption(t *testing.T) {
	// 建立測試環境
	service, _ := createTestEditorService()
	editorSvc := service.(*editorService)
	
	// 建立測試筆記並啟用加密
	note, err := service.CreateNote("測試筆記", "測試內容")
	if err != nil {
		t.Fatalf("建立筆記失敗: %v", err)
	}
	
	// 先啟用加密
	err = editorSvc.EnableEncryption(note.ID, "TestPassword123!", "aes256", false)
	if err != nil {
		t.Fatalf("啟用加密失敗: %v", err)
	}
	
	// 停用加密
	err = editorSvc.DisableEncryption(note.ID)
	if err != nil {
		t.Fatalf("停用加密失敗: %v", err)
	}
	
	// 驗證加密狀態
	if note.IsEncrypted {
		t.Error("筆記不應該標記為已加密")
	}
	
	if note.EncryptionType != "" {
		t.Errorf("加密類型應該為空，實際: %s", note.EncryptionType)
	}
}

// TestIsEncrypted 測試檢查加密狀態功能
func TestIsEncrypted(t *testing.T) {
	// 建立測試環境
	service, _ := createTestEditorService()
	editorSvc := service.(*editorService)
	
	// 建立測試筆記
	note, err := service.CreateNote("測試筆記", "測試內容")
	if err != nil {
		t.Fatalf("建立筆記失敗: %v", err)
	}
	
	// 測試未加密狀態
	isEncrypted, exists := editorSvc.IsEncrypted(note.ID)
	if !exists {
		t.Error("筆記應該存在")
	}
	if isEncrypted {
		t.Error("新筆記不應該加密")
	}
	
	// 啟用加密
	err = editorSvc.EnableEncryption(note.ID, "TestPassword123!", "aes256", false)
	if err != nil {
		t.Fatalf("啟用加密失敗: %v", err)
	}
	
	// 測試已加密狀態
	isEncrypted, exists = editorSvc.IsEncrypted(note.ID)
	if !exists {
		t.Error("筆記應該存在")
	}
	if !isEncrypted {
		t.Error("筆記應該已加密")
	}
	
	// 測試不存在的筆記
	isEncrypted, exists = editorSvc.IsEncrypted("non-existent-id")
	if exists {
		t.Error("不存在的筆記不應該回傳存在狀態")
	}
}

// TestGetEncryptionType 測試取得加密類型功能
func TestGetEncryptionType(t *testing.T) {
	// 建立測試環境
	service, _ := createTestEditorService()
	editorSvc := service.(*editorService)
	
	// 建立測試筆記
	note, err := service.CreateNote("測試筆記", "測試內容")
	if err != nil {
		t.Fatalf("建立筆記失敗: %v", err)
	}
	
	// 測試未加密狀態
	encType, exists := editorSvc.GetEncryptionType(note.ID)
	if !exists {
		t.Error("筆記應該存在")
	}
	if encType != "" {
		t.Errorf("未加密筆記的加密類型應該為空，實際: %s", encType)
	}
	
	// 啟用加密
	err = editorSvc.EnableEncryption(note.ID, "TestPassword123!", "chacha20", false)
	if err != nil {
		t.Fatalf("啟用加密失敗: %v", err)
	}
	
	// 測試已加密狀態
	encType, exists = editorSvc.GetEncryptionType(note.ID)
	if !exists {
		t.Error("筆記應該存在")
	}
	if encType != "chacha20" {
		t.Errorf("加密類型不正確，期望: chacha20，實際: %s", encType)
	}
}