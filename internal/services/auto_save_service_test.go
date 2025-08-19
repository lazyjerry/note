// Package services 包含自動保存服務的單元測試
// 測試自動保存功能的各種場景，包含正常操作、錯誤處理和邊界條件
package services

import (
	"mac-notebook-app/internal/models" // 引入資料模型
	"testing"                          // Go 標準測試套件
	"time"                             // 時間處理套件
)

// MockEditorService 模擬編輯器服務，用於測試自動保存功能
// 實作 EditorService 介面，提供可控制的保存行為
type MockEditorService struct {
	saveError    error                    // 模擬保存時的錯誤
	saveCallLog  []*models.Note           // 記錄所有保存呼叫的筆記
	shouldFail   bool                     // 控制是否模擬保存失敗
	saveDelay    time.Duration            // 模擬保存操作的延遲
	notes        map[string]*models.Note  // 儲存的筆記
}

// CreateNote 模擬建立筆記功能
// 參數：title（標題）、content（內容）
// 回傳：新建立的筆記和可能的錯誤
func (m *MockEditorService) CreateNote(title, content string) (*models.Note, error) {
	return models.NewNote(title, content, "/test/path/"+title+".md"), nil
}

// OpenNote 模擬開啟筆記功能
// 參數：filePath（檔案路徑）
// 回傳：筆記實例和可能的錯誤
func (m *MockEditorService) OpenNote(filePath string) (*models.Note, error) {
	return models.NewNote("Test Note", "Test Content", filePath), nil
}

// SaveNote 模擬保存筆記功能，記錄保存呼叫並可模擬錯誤
// 參數：note（要保存的筆記）
// 回傳：可能的錯誤
//
// 執行流程：
// 1. 記錄保存呼叫到日誌中
// 2. 如果設定了延遲，則等待指定時間
// 3. 根據 shouldFail 標誌決定是否回傳錯誤
// 4. 如果保存成功，標記筆記為已保存
func (m *MockEditorService) SaveNote(note *models.Note) error {
	// 記錄保存呼叫
	m.saveCallLog = append(m.saveCallLog, note)
	
	// 模擬保存延遲
	if m.saveDelay > 0 {
		time.Sleep(m.saveDelay)
	}
	
	// 模擬保存失敗
	if m.shouldFail {
		return m.saveError
	}
	
	// 保存成功，標記筆記為已保存
	note.MarkSaved()
	return nil
}

// UpdateContent 模擬更新筆記內容功能
// 參數：noteID（筆記 ID）、content（新內容）
// 回傳：可能的錯誤
func (m *MockEditorService) UpdateContent(noteID, content string) error {
	return nil
}

// PreviewMarkdown 模擬 Markdown 預覽功能
// 參數：content（Markdown 內容）
// 回傳：HTML 字串
func (m *MockEditorService) PreviewMarkdown(content string) string {
	return "<p>" + content + "</p>"
}

// DecryptWithPassword 模擬使用密碼解密筆記內容功能
// 參數：noteID（筆記 ID）、password（解密密碼）
// 回傳：解密後的內容和可能的錯誤
func (m *MockEditorService) DecryptWithPassword(noteID, password string) (string, error) {
	return "解密後的內容", nil
}

// GetActiveNotes 模擬取得所有活躍筆記功能
// 回傳：活躍筆記的映射表
func (m *MockEditorService) GetActiveNotes() map[string]*models.Note {
	return make(map[string]*models.Note)
}

// CloseNote 模擬關閉筆記功能
// 參數：noteID（筆記 ID）
func (m *MockEditorService) CloseNote(noteID string) {
	// 模擬關閉筆記操作
}

// GetActiveNote 模擬取得活躍筆記功能
// 參數：noteID（筆記 ID）
// 回傳：筆記實例和是否存在
func (m *MockEditorService) GetActiveNote(noteID string) (*models.Note, bool) {
	if m.notes == nil {
		m.notes = make(map[string]*models.Note)
	}
	note, exists := m.notes[noteID]
	return note, exists
}

// ========== 智慧編輯功能模擬實作 ==========

// GetAutoCompleteSuggestions 模擬自動完成建議功能
func (m *MockEditorService) GetAutoCompleteSuggestions(content string, cursorPosition int) []AutoCompleteSuggestion {
	return []AutoCompleteSuggestion{}
}

// FormatTableContent 模擬表格格式化功能
func (m *MockEditorService) FormatTableContent(tableContent string) (string, error) {
	return tableContent, nil
}

// InsertLinkMarkdown 模擬連結插入功能
func (m *MockEditorService) InsertLinkMarkdown(text, url string) string {
	return "[" + text + "](" + url + ")"
}

// InsertImageMarkdown 模擬圖片插入功能
func (m *MockEditorService) InsertImageMarkdown(altText, imagePath string) string {
	return "![" + altText + "](" + imagePath + ")"
}

// GetSupportedCodeLanguages 模擬取得支援語言列表功能
func (m *MockEditorService) GetSupportedCodeLanguages() []string {
	return []string{"go", "javascript", "python"}
}

// FormatCodeBlockMarkdown 模擬程式碼區塊格式化功能
func (m *MockEditorService) FormatCodeBlockMarkdown(code, language string) string {
	return "```" + language + "\n" + code + "\n```"
}

// FormatMathExpressionMarkdown 模擬數學公式格式化功能
func (m *MockEditorService) FormatMathExpressionMarkdown(expression string, isInline bool) string {
	if isInline {
		return "$" + expression + "$"
	}
	return "$$\n" + expression + "\n$$"
}

// ValidateMarkdownContent 模擬 Markdown 內容驗證功能
func (m *MockEditorService) ValidateMarkdownContent(content string) (bool, []string) {
	return true, []string{}
}

// GenerateTableTemplateMarkdown 模擬表格模板生成功能
func (m *MockEditorService) GenerateTableTemplateMarkdown(rows, cols int) string {
	return "| 欄位1 | 欄位2 |\n|-------|-------|\n| 內容1 | 內容2 |"
}

// PreviewMarkdownWithHighlight 模擬帶語法高亮的預覽功能
func (m *MockEditorService) PreviewMarkdownWithHighlight(content string) string {
	return "<p>" + content + "</p>"
}

// GetSmartEditingService 模擬取得智慧編輯服務功能
func (m *MockEditorService) GetSmartEditingService() SmartEditingService {
	return NewSmartEditingService()
}

// SetSmartEditingService 模擬設定智慧編輯服務功能
func (m *MockEditorService) SetSmartEditingService(smartEditSvc SmartEditingService) {
	// 模擬設定操作
}

// NewMockEditorService 建立新的模擬編輯器服務
func NewMockEditorService() *MockEditorService {
	return &MockEditorService{
		notes: make(map[string]*models.Note),
	}
}

// GetSaveCallCount 取得保存呼叫的次數
// 回傳：保存呼叫次數
func (m *MockEditorService) GetSaveCallCount() int {
	return len(m.saveCallLog)
}

// GetLastSavedNote 取得最後一次保存的筆記
// 回傳：最後保存的筆記實例，如果沒有則回傳 nil
func (m *MockEditorService) GetLastSavedNote() *models.Note {
	if len(m.saveCallLog) == 0 {
		return nil
	}
	return m.saveCallLog[len(m.saveCallLog)-1]
}

// ClearSaveCallLog 清空保存呼叫日誌
func (m *MockEditorService) ClearSaveCallLog() {
	m.saveCallLog = []*models.Note{}
}

// SetShouldFail 設定是否模擬保存失敗
// 參數：shouldFail（是否失敗）、err（錯誤實例）
func (m *MockEditorService) SetShouldFail(shouldFail bool, err error) {
	m.shouldFail = shouldFail
	m.saveError = err
}

// SetSaveDelay 設定保存操作的模擬延遲
// 參數：delay（延遲時間）
func (m *MockEditorService) SetSaveDelay(delay time.Duration) {
	m.saveDelay = delay
}

// TestNewAutoSaveService 測試自動保存服務的建立
// 驗證服務實例是否正確初始化
func TestNewAutoSaveService(t *testing.T) {
	mockEditor := NewMockEditorService()
	service := NewAutoSaveServiceWithDefaults(mockEditor)

	// 驗證服務實例不為 nil
	if service == nil {
		t.Fatal("自動保存服務實例不應該為 nil")
	}

	// 驗證編輯器服務依賴已正確設定
	if service.editorService != mockEditor {
		t.Error("編輯器服務依賴設定不正確")
	}

	// 驗證內部映射表已初始化
	if service.timers == nil {
		t.Error("定時器映射表未初始化")
	}
	if service.saveStatus == nil {
		t.Error("保存狀態映射表未初始化")
	}
	if service.notes == nil {
		t.Error("筆記快取映射表未初始化")
	}
}

// TestStartAutoSave 測試啟動自動保存功能
// 驗證自動保存是否正確啟動並設定相關狀態
func TestStartAutoSave(t *testing.T) {
	mockEditor := &MockEditorService{}
	service := NewAutoSaveServiceWithDefaults(mockEditor)

	// 建立測試筆記
	note := models.NewNote("測試筆記", "測試內容", "/test/note.md")
	interval := 100 * time.Millisecond

	// 啟動自動保存
	service.StartAutoSave(note, interval)

	// 驗證筆記已加入快取
	if _, exists := service.notes[note.ID]; !exists {
		t.Error("筆記未正確加入快取")
	}

	// 驗證保存狀態已初始化
	status := service.GetSaveStatus(note.ID)
	if status.NoteID != note.ID {
		t.Error("保存狀態的筆記 ID 不正確")
	}
	if status.IsSaving {
		t.Error("初始保存狀態不應該為保存中")
	}

	// 驗證自動保存已啟用
	if !service.IsAutoSaveActive(note.ID) {
		t.Error("自動保存未正確啟用")
	}

	// 清理資源
	service.StopAutoSave(note.ID)
}

// TestStopAutoSave 測試停止自動保存功能
// 驗證自動保存是否正確停止並清理相關資源
func TestStopAutoSave(t *testing.T) {
	mockEditor := &MockEditorService{}
	service := NewAutoSaveServiceWithDefaults(mockEditor)

	// 建立測試筆記並啟動自動保存
	note := models.NewNote("測試筆記", "測試內容", "/test/note.md")
	interval := 100 * time.Millisecond
	service.StartAutoSave(note, interval)

	// 驗證自動保存已啟用
	if !service.IsAutoSaveActive(note.ID) {
		t.Error("自動保存應該已啟用")
	}

	// 停止自動保存
	service.StopAutoSave(note.ID)

	// 驗證自動保存已停止
	if service.IsAutoSaveActive(note.ID) {
		t.Error("自動保存應該已停止")
	}

	// 驗證狀態已清理
	status := service.GetSaveStatus(note.ID)
	if status.LastError == nil {
		t.Error("停止後的狀態查詢應該回傳錯誤")
	}
}

// TestSaveNow 測試立即保存功能
// 驗證立即保存是否正確執行並更新狀態
func TestSaveNow(t *testing.T) {
	mockEditor := &MockEditorService{}
	service := NewAutoSaveServiceWithDefaults(mockEditor)

	// 建立測試筆記並啟動自動保存
	note := models.NewNote("測試筆記", "測試內容", "/test/note.md")
	note.UpdateContent("修改後的內容") // 標記為已修改
	interval := 1 * time.Hour        // 設定長間隔避免自動觸發
	service.StartAutoSave(note, interval)

	// 執行立即保存
	err := service.SaveNow(note.ID)
	if err != nil {
		t.Fatalf("立即保存失敗: %v", err)
	}

	// 驗證編輯器服務被呼叫
	if mockEditor.GetSaveCallCount() != 1 {
		t.Errorf("預期保存呼叫次數為 1，實際為 %d", mockEditor.GetSaveCallCount())
	}

	// 驗證保存狀態已更新
	status := service.GetSaveStatus(note.ID)
	if status.SaveCount != 1 {
		t.Errorf("預期保存計數為 1，實際為 %d", status.SaveCount)
	}
	if status.LastError != nil {
		t.Errorf("保存後不應該有錯誤: %v", status.LastError)
	}

	// 清理資源
	service.StopAutoSave(note.ID)
}

// TestSaveNowWithNonExistentNote 測試對不存在筆記的立即保存
// 驗證錯誤處理是否正確
func TestSaveNowWithNonExistentNote(t *testing.T) {
	mockEditor := &MockEditorService{}
	service := NewAutoSaveServiceWithDefaults(mockEditor)

	// 嘗試保存不存在的筆記
	err := service.SaveNow("不存在的筆記ID")
	if err == nil {
		t.Error("對不存在筆記的保存應該回傳錯誤")
	}

	// 驗證錯誤類型
	if appErr, ok := err.(*models.AppError); ok {
		if appErr.Code != "NOTE_NOT_FOUND" {
			t.Errorf("預期錯誤代碼為 NOTE_NOT_FOUND，實際為 %s", appErr.Code)
		}
	} else {
		t.Error("預期回傳 AppError 類型的錯誤")
	}
}

// TestSaveNowWithSaveInProgress 測試在保存進行中時的立即保存
// 驗證並發保存的防護機制
func TestSaveNowWithSaveInProgress(t *testing.T) {
	mockEditor := &MockEditorService{}
	mockEditor.SetSaveDelay(200 * time.Millisecond) // 設定保存延遲
	service := NewAutoSaveServiceWithDefaults(mockEditor)

	// 建立測試筆記並啟動自動保存
	note := models.NewNote("測試筆記", "測試內容", "/test/note.md")
	note.UpdateContent("修改後的內容")
	service.StartAutoSave(note, 1*time.Hour)

	// 啟動第一個保存操作（會有延遲）
	go func() {
		service.SaveNow(note.ID)
	}()

	// 等待一小段時間確保第一個保存已開始
	time.Sleep(50 * time.Millisecond)

	// 嘗試第二個立即保存
	err := service.SaveNow(note.ID)
	if err == nil {
		t.Error("在保存進行中時應該回傳錯誤")
	}

	// 驗證錯誤類型
	if appErr, ok := err.(*models.AppError); ok {
		if appErr.Code != "SAVE_IN_PROGRESS" {
			t.Errorf("預期錯誤代碼為 SAVE_IN_PROGRESS，實際為 %s", appErr.Code)
		}
	} else {
		t.Error("預期回傳 AppError 類型的錯誤")
	}

	// 等待第一個保存完成
	time.Sleep(200 * time.Millisecond)

	// 清理資源
	service.StopAutoSave(note.ID)
}

// TestAutoSaveWithModifiedNote 測試自動保存修改過的筆記
// 驗證自動保存是否在筆記修改後正確觸發
func TestAutoSaveWithModifiedNote(t *testing.T) {
	mockEditor := &MockEditorService{}
	service := NewAutoSaveServiceWithDefaults(mockEditor)

	// 建立測試筆記
	note := models.NewNote("測試筆記", "測試內容", "/test/note.md")
	note.MarkSaved() // 先標記為已保存
	interval := 50 * time.Millisecond

	// 啟動自動保存
	service.StartAutoSave(note, interval)

	// 修改筆記內容
	note.UpdateContent("修改後的內容")

	// 等待自動保存觸發
	time.Sleep(100 * time.Millisecond)

	// 驗證自動保存被觸發
	if mockEditor.GetSaveCallCount() == 0 {
		t.Error("自動保存應該被觸發")
	}

	// 驗證保存狀態已更新
	status := service.GetSaveStatus(note.ID)
	if status.SaveCount == 0 {
		t.Error("保存計數應該大於 0")
	}

	// 清理資源
	service.StopAutoSave(note.ID)
}

// TestAutoSaveWithUnmodifiedNote 測試自動保存未修改的筆記
// 驗證自動保存是否正確跳過未修改的筆記
func TestAutoSaveWithUnmodifiedNote(t *testing.T) {
	mockEditor := &MockEditorService{}
	service := NewAutoSaveServiceWithDefaults(mockEditor)

	// 建立測試筆記並標記為已保存
	note := models.NewNote("測試筆記", "測試內容", "/test/note.md")
	note.MarkSaved()
	interval := 50 * time.Millisecond

	// 啟動自動保存
	service.StartAutoSave(note, interval)

	// 等待自動保存觸發時間
	time.Sleep(100 * time.Millisecond)

	// 驗證自動保存未被觸發（因為筆記未修改）
	if mockEditor.GetSaveCallCount() > 0 {
		t.Error("未修改的筆記不應該觸發自動保存")
	}

	// 清理資源
	service.StopAutoSave(note.ID)
}

// TestGetSaveStatus 測試取得保存狀態功能
// 驗證狀態資訊是否正確回傳
func TestGetSaveStatus(t *testing.T) {
	mockEditor := &MockEditorService{}
	service := NewAutoSaveServiceWithDefaults(mockEditor)

	// 建立測試筆記並啟動自動保存
	note := models.NewNote("測試筆記", "測試內容", "/test/note.md")
	service.StartAutoSave(note, 1*time.Hour)

	// 取得保存狀態
	status := service.GetSaveStatus(note.ID)

	// 驗證狀態資訊
	if status.NoteID != note.ID {
		t.Errorf("預期筆記 ID 為 %s，實際為 %s", note.ID, status.NoteID)
	}
	if status.IsSaving {
		t.Error("初始狀態不應該為保存中")
	}
	if status.SaveCount != 0 {
		t.Errorf("初始保存計數應該為 0，實際為 %d", status.SaveCount)
	}

	// 清理資源
	service.StopAutoSave(note.ID)
}

// TestGetAllSaveStatuses 測試取得所有保存狀態功能
// 驗證是否正確回傳所有筆記的狀態
func TestGetAllSaveStatuses(t *testing.T) {
	mockEditor := &MockEditorService{}
	service := NewAutoSaveServiceWithDefaults(mockEditor)

	// 建立多個測試筆記
	note1 := models.NewNote("筆記1", "內容1", "/test/note1.md")
	note2 := models.NewNote("筆記2", "內容2", "/test/note2.md")

	// 啟動自動保存
	service.StartAutoSave(note1, 1*time.Hour)
	service.StartAutoSave(note2, 1*time.Hour)

	// 取得所有保存狀態
	allStatuses := service.GetAllSaveStatuses()

	// 驗證狀態數量
	if len(allStatuses) != 2 {
		t.Errorf("預期狀態數量為 2，實際為 %d", len(allStatuses))
	}

	// 驗證每個狀態
	if _, exists := allStatuses[note1.ID]; !exists {
		t.Error("筆記1的狀態未找到")
	}
	if _, exists := allStatuses[note2.ID]; !exists {
		t.Error("筆記2的狀態未找到")
	}

	// 清理資源
	service.StopAutoSave(note1.ID)
	service.StopAutoSave(note2.ID)
}

// TestShutdown 測試服務關閉功能
// 驗證所有資源是否正確清理
func TestShutdown(t *testing.T) {
	mockEditor := &MockEditorService{}
	service := NewAutoSaveServiceWithDefaults(mockEditor)

	// 建立測試筆記並啟動自動保存
	note := models.NewNote("測試筆記", "測試內容", "/test/note.md")
	service.StartAutoSave(note, 1*time.Hour)

	// 驗證自動保存已啟用
	if !service.IsAutoSaveActive(note.ID) {
		t.Error("自動保存應該已啟用")
	}

	// 關閉服務
	service.Shutdown()

	// 驗證所有資源已清理
	if service.IsAutoSaveActive(note.ID) {
		t.Error("關閉後自動保存應該已停止")
	}

	// 驗證狀態已清理
	allStatuses := service.GetAllSaveStatuses()
	if len(allStatuses) != 0 {
		t.Errorf("關閉後狀態應該已清理，實際還有 %d 個狀態", len(allStatuses))
	}
}

// TestSaveErrorHandling 測試保存錯誤處理
// 驗證保存失敗時的錯誤處理和狀態更新
func TestSaveErrorHandling(t *testing.T) {
	mockEditor := &MockEditorService{}
	service := NewAutoSaveServiceWithDefaults(mockEditor)

	// 設定模擬保存失敗
	saveError := models.NewAppError("SAVE_FAILED", "保存失敗", "磁碟空間不足")
	mockEditor.SetShouldFail(true, saveError)

	// 建立測試筆記並啟動自動保存
	note := models.NewNote("測試筆記", "測試內容", "/test/note.md")
	note.UpdateContent("修改後的內容")
	service.StartAutoSave(note, 50*time.Millisecond)

	// 等待自動保存觸發
	time.Sleep(100 * time.Millisecond)

	// 驗證保存狀態包含錯誤資訊
	status := service.GetSaveStatus(note.ID)
	if status.LastError == nil {
		t.Error("保存失敗後狀態應該包含錯誤資訊")
	}

	// 驗證錯誤類型和內容
	if appErr, ok := status.LastError.(*models.AppError); ok {
		if appErr.Code != "SAVE_FAILED" {
			t.Errorf("預期錯誤代碼為 SAVE_FAILED，實際為 %s", appErr.Code)
		}
	} else {
		t.Error("預期錯誤類型為 AppError")
	}

	// 清理資源
	service.StopAutoSave(note.ID)
}

// TestConcurrentAutoSave 測試並發自動保存
// 驗證多個筆記同時進行自動保存時的執行緒安全性
func TestConcurrentAutoSave(t *testing.T) {
	mockEditor := &MockEditorService{}
	service := NewAutoSaveServiceWithDefaults(mockEditor)

	// 建立多個測試筆記
	notes := make([]*models.Note, 5)
	for i := 0; i < 5; i++ {
		notes[i] = models.NewNote(
			"筆記"+string(rune('1'+i)),
			"內容"+string(rune('1'+i)),
			"/test/note"+string(rune('1'+i))+".md",
		)
		notes[i].UpdateContent("修改後的內容" + string(rune('1'+i)))
	}

	// 同時啟動多個自動保存
	for _, note := range notes {
		service.StartAutoSave(note, 50*time.Millisecond)
	}

	// 等待所有自動保存觸發
	time.Sleep(200 * time.Millisecond)

	// 驗證至少有一些筆記被保存（並發環境下可能不是全部）
	if mockEditor.GetSaveCallCount() == 0 {
		t.Error("應該至少有一些筆記被自動保存")
	}

	// 驗證所有狀態都正確更新
	allStatuses := service.GetAllSaveStatuses()
	if len(allStatuses) != len(notes) {
		t.Errorf("預期狀態數量為 %d，實際為 %d", len(notes), len(allStatuses))
	}

	// 清理資源
	for _, note := range notes {
		service.StopAutoSave(note.ID)
	}
}

// TestRescheduleTimer 測試定時器重新排程
// 驗證自動保存定時器是否正確重新設定
func TestRescheduleTimer(t *testing.T) {
	mockEditor := &MockEditorService{}
	service := NewAutoSaveServiceWithDefaults(mockEditor)

	// 建立測試筆記
	note := models.NewNote("測試筆記", "測試內容", "/test/note.md")
	note.MarkSaved() // 標記為已保存，第一次不會觸發保存
	interval := 50 * time.Millisecond

	// 啟動自動保存
	service.StartAutoSave(note, interval)

	// 等待第一次定時器觸發（不會保存因為未修改）
	time.Sleep(100 * time.Millisecond)

	// 修改筆記內容
	note.UpdateContent("修改後的內容")

	// 等待第二次定時器觸發（會保存因為已修改）
	time.Sleep(100 * time.Millisecond)

	// 驗證保存被觸發
	if mockEditor.GetSaveCallCount() == 0 {
		t.Error("修改後的筆記應該被自動保存")
	}

	// 清理資源
	service.StopAutoSave(note.ID)
}
// MockSettingsService 模擬設定服務，用於測試自動保存配置功能
type MockSettingsService struct {
	settings *models.Settings
	loadError error
	saveError error
}

// LoadSettings 模擬載入設定功能
func (m *MockSettingsService) LoadSettings() (*models.Settings, error) {
	if m.loadError != nil {
		return nil, m.loadError
	}
	if m.settings == nil {
		return models.NewDefaultSettings(), nil
	}
	return m.settings, nil
}

// SaveSettings 模擬保存設定功能
func (m *MockSettingsService) SaveSettings(settings *models.Settings) error {
	if m.saveError != nil {
		return m.saveError
	}
	m.settings = settings
	return nil
}

// GetDefaultSettings 模擬取得預設設定功能
func (m *MockSettingsService) GetDefaultSettings() *models.Settings {
	return models.NewDefaultSettings()
}

// SetSettings 設定模擬的設定值
func (m *MockSettingsService) SetSettings(settings *models.Settings) {
	m.settings = settings
}

// SetLoadError 設定載入設定時的錯誤
func (m *MockSettingsService) SetLoadError(err error) {
	m.loadError = err
}

// SetSaveError 設定保存設定時的錯誤
func (m *MockSettingsService) SetSaveError(err error) {
	m.saveError = err
}

// TestNewAutoSaveServiceWithSettings 測試使用設定服務建立自動保存服務
func TestNewAutoSaveServiceWithSettings(t *testing.T) {
	mockEditor := &MockEditorService{}
	mockSettings := &MockSettingsService{}
	service := NewAutoSaveService(mockEditor, mockSettings)

	// 驗證服務實例不為 nil
	if service == nil {
		t.Fatal("自動保存服務實例不應該為 nil")
	}

	// 驗證設定服務依賴已正確設定
	if service.settingsService != mockSettings {
		t.Error("設定服務依賴設定不正確")
	}
}

// TestStartAutoSaveWithSettings 測試使用設定服務的自動保存間隔
func TestStartAutoSaveWithSettings(t *testing.T) {
	mockEditor := &MockEditorService{}
	mockSettings := &MockSettingsService{}
	
	// 設定自訂的自動保存間隔
	customSettings := models.NewDefaultSettings()
	customSettings.AutoSaveInterval = 10 // 10 分鐘
	mockSettings.SetSettings(customSettings)
	
	service := NewAutoSaveService(mockEditor, mockSettings)

	// 建立測試筆記
	note := models.NewNote("測試筆記", "測試內容", "/test/note.md")
	
	// 使用設定服務的間隔啟動自動保存
	service.StartAutoSaveWithSettings(note)

	// 驗證自動保存已啟用
	if !service.IsAutoSaveActive(note.ID) {
		t.Error("自動保存應該已啟用")
	}

	// 清理資源
	service.StopAutoSave(note.ID)
}

// TestEncryptedFileAutoSave 測試加密檔案的自動保存
func TestEncryptedFileAutoSave(t *testing.T) {
	mockEditor := &MockEditorService{}
	mockSettings := &MockSettingsService{}
	service := NewAutoSaveService(mockEditor, mockSettings)

	// 建立加密測試筆記
	note := models.NewNote("加密筆記", "機密內容", "/test/encrypted.md.enc")
	note.SetEncryption("password") // 設定為密碼加密
	note.UpdateContent("修改後的機密內容")
	
	interval := 50 * time.Millisecond
	service.StartAutoSave(note, interval)

	// 等待自動保存觸發
	time.Sleep(100 * time.Millisecond)

	// 驗證加密檔案被保存
	if mockEditor.GetSaveCallCount() == 0 {
		t.Error("加密檔案應該被自動保存")
	}

	// 驗證保存狀態
	status := service.GetSaveStatus(note.ID)
	if status.SaveCount == 0 {
		t.Error("加密檔案的保存計數應該大於 0")
	}

	// 清理資源
	service.StopAutoSave(note.ID)
}

// TestEncryptedFileRetryMechanism 測試加密檔案的重試機制
func TestEncryptedFileRetryMechanism(t *testing.T) {
	mockEditor := &MockEditorService{}
	mockSettings := &MockSettingsService{}
	service := NewAutoSaveService(mockEditor, mockSettings)

	// 設定模擬保存失敗
	saveError := models.NewAppError("ENCRYPTION_FAILED", "加密失敗", "密鑰錯誤")
	mockEditor.SetShouldFail(true, saveError)

	// 建立加密測試筆記
	note := models.NewNote("加密筆記", "機密內容", "/test/encrypted.md.enc")
	note.SetEncryption("password")
	note.UpdateContent("修改後的內容")

	// 執行立即保存（會觸發重試機制）
	service.StartAutoSave(note, 1*time.Hour) // 長間隔避免自動觸發
	err := service.SaveNow(note.ID)

	// 驗證重試機制被觸發
	if err == nil {
		t.Error("加密檔案保存失敗應該回傳錯誤")
	}

	// 驗證錯誤類型
	if appErr, ok := err.(*models.AppError); ok {
		if appErr.Code != "ENCRYPTED_SAVE_RETRY_FAILED" {
			t.Errorf("預期錯誤代碼為 ENCRYPTED_SAVE_RETRY_FAILED，實際為 %s", appErr.Code)
		}
	} else {
		t.Error("預期回傳 AppError 類型的錯誤")
	}

	// 驗證多次保存嘗試
	if mockEditor.GetSaveCallCount() < 3 {
		t.Errorf("預期至少重試 3 次，實際重試 %d 次", mockEditor.GetSaveCallCount())
	}

	// 清理資源
	service.StopAutoSave(note.ID)
}

// TestGetEncryptedFileCount 測試取得加密檔案數量
func TestGetEncryptedFileCount(t *testing.T) {
	mockEditor := &MockEditorService{}
	service := NewAutoSaveServiceWithDefaults(mockEditor)

	// 建立混合的測試筆記
	note1 := models.NewNote("一般筆記", "內容1", "/test/note1.md")
	note2 := models.NewNote("加密筆記1", "機密內容1", "/test/encrypted1.md.enc")
	note2.SetEncryption("password")
	note3 := models.NewNote("加密筆記2", "機密內容2", "/test/encrypted2.md.enc")
	note3.SetEncryption("biometric")

	// 啟動自動保存
	service.StartAutoSave(note1, 1*time.Hour)
	service.StartAutoSave(note2, 1*time.Hour)
	service.StartAutoSave(note3, 1*time.Hour)

	// 驗證加密檔案數量
	encryptedCount := service.GetEncryptedFileCount()
	if encryptedCount != 2 {
		t.Errorf("預期加密檔案數量為 2，實際為 %d", encryptedCount)
	}

	// 清理資源
	service.StopAutoSave(note1.ID)
	service.StopAutoSave(note2.ID)
	service.StopAutoSave(note3.ID)
}

// TestUpdateAutoSaveInterval 測試更新自動保存間隔
func TestUpdateAutoSaveInterval(t *testing.T) {
	mockEditor := &MockEditorService{}
	service := NewAutoSaveServiceWithDefaults(mockEditor)

	// 建立測試筆記並啟動自動保存
	note := models.NewNote("測試筆記", "測試內容", "/test/note.md")
	service.StartAutoSave(note, 1*time.Hour)

	// 更新自動保存間隔
	newInterval := 30 * time.Second
	err := service.UpdateAutoSaveInterval(note.ID, newInterval)
	if err != nil {
		t.Fatalf("更新自動保存間隔失敗: %v", err)
	}

	// 驗證自動保存仍然活躍
	if !service.IsAutoSaveActive(note.ID) {
		t.Error("更新間隔後自動保存應該仍然活躍")
	}

	// 清理資源
	service.StopAutoSave(note.ID)
}

// TestUpdateAutoSaveIntervalWithNonExistentNote 測試更新不存在筆記的間隔
func TestUpdateAutoSaveIntervalWithNonExistentNote(t *testing.T) {
	mockEditor := &MockEditorService{}
	service := NewAutoSaveServiceWithDefaults(mockEditor)

	// 嘗試更新不存在筆記的間隔
	err := service.UpdateAutoSaveInterval("不存在的筆記ID", 30*time.Second)
	if err == nil {
		t.Error("更新不存在筆記的間隔應該回傳錯誤")
	}

	// 驗證錯誤類型
	if appErr, ok := err.(*models.AppError); ok {
		if appErr.Code != "NOTE_NOT_FOUND" {
			t.Errorf("預期錯誤代碼為 NOTE_NOT_FOUND，實際為 %s", appErr.Code)
		}
	} else {
		t.Error("預期回傳 AppError 類型的錯誤")
	}
}

// TestSetEncryptedBackoff 測試設定加密檔案延遲
func TestSetEncryptedBackoff(t *testing.T) {
	mockEditor := &MockEditorService{}
	service := NewAutoSaveServiceWithDefaults(mockEditor)

	// 設定新的加密檔案延遲
	newBackoff := 1 * time.Minute
	service.SetEncryptedBackoff(newBackoff)

	// 建立加密筆記並測試間隔
	note := models.NewNote("加密筆記", "機密內容", "/test/encrypted.md.enc")
	note.SetEncryption("password")

	// 使用設定啟動自動保存
	service.StartAutoSaveWithSettings(note)

	// 驗證自動保存已啟用
	if !service.IsAutoSaveActive(note.ID) {
		t.Error("加密筆記的自動保存應該已啟用")
	}

	// 清理資源
	service.StopAutoSave(note.ID)
}

// TestAutoSaveWithSettingsLoadError 測試設定載入錯誤的處理
func TestAutoSaveWithSettingsLoadError(t *testing.T) {
	mockEditor := &MockEditorService{}
	mockSettings := &MockSettingsService{}
	
	// 設定載入錯誤
	mockSettings.SetLoadError(models.NewAppError("SETTINGS_LOAD_FAILED", "設定載入失敗", "檔案不存在"))
	
	service := NewAutoSaveService(mockEditor, mockSettings)

	// 建立測試筆記
	note := models.NewNote("測試筆記", "測試內容", "/test/note.md")
	
	// 使用設定服務啟動自動保存（應該回退到預設值）
	service.StartAutoSaveWithSettings(note)

	// 驗證自動保存仍然啟用（使用預設間隔）
	if !service.IsAutoSaveActive(note.ID) {
		t.Error("即使設定載入失敗，自動保存也應該使用預設間隔啟用")
	}

	// 清理資源
	service.StopAutoSave(note.ID)
}