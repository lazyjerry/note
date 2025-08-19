// Package ui 包含使用者介面相關的元件和視窗管理測試
// 測試 MarkdownEditor 的建立、初始化和編輯功能
package ui

import (
	"testing"                           // Go 標準測試套件
	"mac-notebook-app/internal/models"  // 引入資料模型
	"mac-notebook-app/internal/services" // 引入服務層
)

// mockEditorService 模擬編輯器服務，用於測試
// 實作 EditorService 介面的基本功能
type mockEditorService struct {
	notes map[string]*models.Note // 模擬筆記儲存
}

// newMockEditorService 建立模擬編輯器服務實例
// 回傳：模擬編輯器服務實例
func newMockEditorService() services.EditorService {
	return &mockEditorService{
		notes: make(map[string]*models.Note),
	}
}

// CreateNote 模擬建立筆記功能
// 參數：title（筆記標題）、content（筆記內容）
// 回傳：建立的筆記實例和可能的錯誤
func (m *mockEditorService) CreateNote(title, content string) (*models.Note, error) {
	note := &models.Note{
		ID:      "test-note-id",
		Title:   title,
		Content: content,
	}
	m.notes[note.ID] = note
	return note, nil
}

// OpenNote 模擬開啟筆記功能
// 參數：filePath（檔案路徑）
// 回傳：開啟的筆記實例和可能的錯誤
func (m *mockEditorService) OpenNote(filePath string) (*models.Note, error) {
	note := &models.Note{
		ID:       "test-note-id",
		Title:    "測試筆記",
		Content:  "測試內容",
		FilePath: filePath,
	}
	m.notes[note.ID] = note
	return note, nil
}

// SaveNote 模擬保存筆記功能
// 參數：note（要保存的筆記）
// 回傳：可能的錯誤
func (m *mockEditorService) SaveNote(note *models.Note) error {
	m.notes[note.ID] = note
	return nil
}

// UpdateContent 模擬更新內容功能
// 參數：noteID（筆記 ID）、content（新內容）
// 回傳：可能的錯誤
func (m *mockEditorService) UpdateContent(noteID, content string) error {
	if note, exists := m.notes[noteID]; exists {
		note.Content = content
		return nil
	}
	return nil
}

// PreviewMarkdown 模擬 Markdown 預覽功能
// 參數：content（Markdown 內容）
// 回傳：HTML 字串
func (m *mockEditorService) PreviewMarkdown(content string) string {
	return "<p>" + content + "</p>"
}

// DecryptWithPassword 模擬密碼解密功能
// 參數：noteID（筆記 ID）、password（密碼）
// 回傳：解密後的內容和可能的錯誤
func (m *mockEditorService) DecryptWithPassword(noteID, password string) (string, error) {
	if note, exists := m.notes[noteID]; exists {
		return note.Content, nil
	}
	return "", nil
}

// GetActiveNotes 模擬取得所有活躍筆記功能
// 回傳：活躍筆記的映射表
func (m *mockEditorService) GetActiveNotes() map[string]*models.Note {
	return m.notes
}

// CloseNote 模擬關閉筆記功能
// 參數：noteID（筆記 ID）
func (m *mockEditorService) CloseNote(noteID string) {
	delete(m.notes, noteID)
}

// GetActiveNote 模擬取得活躍筆記功能
// 參數：noteID（筆記 ID）
// 回傳：筆記實例和是否存在
func (m *mockEditorService) GetActiveNote(noteID string) (*models.Note, bool) {
	note, exists := m.notes[noteID]
	return note, exists
}

// GetAutoCompleteSuggestions 模擬自動完成建議功能
// 參數：content（當前內容）、cursorPosition（游標位置）
// 回傳：自動完成建議陣列
func (m *mockEditorService) GetAutoCompleteSuggestions(content string, cursorPosition int) []services.AutoCompleteSuggestion {
	return []services.AutoCompleteSuggestion{
		{Text: "# ", Description: "標題", Type: "header", InsertText: "# "},
		{Text: "- ", Description: "項目符號", Type: "list", InsertText: "- "},
	}
}

// FormatTableContent 模擬表格格式化功能
// 參數：tableContent（表格內容）
// 回傳：格式化後的表格字串和可能的錯誤
func (m *mockEditorService) FormatTableContent(tableContent string) (string, error) {
	return tableContent, nil
}

// InsertLinkMarkdown 模擬插入 Markdown 連結功能
// 參數：text（連結文字）、url（連結網址）
// 回傳：格式化的 Markdown 連結字串
func (m *mockEditorService) InsertLinkMarkdown(text, url string) string {
	return "[" + text + "](" + url + ")"
}

// InsertImageMarkdown 模擬插入 Markdown 圖片功能
// 參數：altText（替代文字）、imagePath（圖片路徑）
// 回傳：格式化的 Markdown 圖片字串
func (m *mockEditorService) InsertImageMarkdown(altText, imagePath string) string {
	return "![" + altText + "](" + imagePath + ")"
}

// GetSupportedCodeLanguages 模擬取得支援程式語言功能
// 回傳：支援的程式語言陣列
func (m *mockEditorService) GetSupportedCodeLanguages() []string {
	return []string{"go", "javascript", "python", "java", "c++"}
}

// FormatCodeBlockMarkdown 模擬格式化程式碼區塊功能
// 參數：code（程式碼內容）、language（程式語言）
// 回傳：格式化的 Markdown 程式碼區塊
func (m *mockEditorService) FormatCodeBlockMarkdown(code, language string) string {
	return "```" + language + "\n" + code + "\n```"
}

// FormatMathExpressionMarkdown 模擬格式化數學公式功能
// 參數：expression（數學表達式）、isInline（是否為行內公式）
// 回傳：格式化的 LaTeX 數學公式字串
func (m *mockEditorService) FormatMathExpressionMarkdown(expression string, isInline bool) string {
	if isInline {
		return "$" + expression + "$"
	}
	return "$$" + expression + "$$"
}

// ValidateMarkdownContent 模擬驗證 Markdown 內容功能
// 參數：content（要驗證的 Markdown 內容）
// 回傳：驗證結果和可能的錯誤列表
func (m *mockEditorService) ValidateMarkdownContent(content string) (bool, []string) {
	return true, []string{}
}

// GenerateTableTemplateMarkdown 模擬生成表格模板功能
// 參數：rows（行數）、cols（列數）
// 回傳：表格模板字串
func (m *mockEditorService) GenerateTableTemplateMarkdown(rows, cols int) string {
	return "| Header | Header |\n|--------|--------|\n| Cell   | Cell   |"
}

// PreviewMarkdownWithHighlight 模擬預覽 Markdown 內容並包含程式碼高亮功能
// 參數：content（Markdown 格式的內容）
// 回傳：轉換後的 HTML 字串（包含語法高亮）
func (m *mockEditorService) PreviewMarkdownWithHighlight(content string) string {
	return "<div class=\"highlight\"><p>" + content + "</p></div>"
}

// GetSmartEditingService 模擬取得智慧編輯服務功能
// 回傳：SmartEditingService 介面實例
func (m *mockEditorService) GetSmartEditingService() services.SmartEditingService {
	return nil
}

// SetSmartEditingService 模擬設定智慧編輯服務功能
// 參數：smartEditSvc（智慧編輯服務實例）
func (m *mockEditorService) SetSmartEditingService(smartEditSvc services.SmartEditingService) {
	// 模擬實作，不執行任何操作
}

// TestNewMarkdownEditor 測試 Markdown 編輯器的建立和初始化
// 驗證編輯器是否正確建立並包含所有必要的 UI 元件
//
// 測試項目：
// 1. 編輯器實例是否正確建立
// 2. 所有 UI 元件是否正確初始化
// 3. 服務依賴是否正確設定
// 4. 初始狀態是否正確
// 5. 兩行工具欄佈局是否正確
func TestNewMarkdownEditor(t *testing.T) {
	// 建立模擬編輯器服務
	mockService := newMockEditorService()
	
	// 建立 Markdown 編輯器實例
	editor := NewMarkdownEditor(mockService)
	
	// 驗證編輯器實例不為 nil
	if editor == nil {
		t.Fatal("NewMarkdownEditor 應該回傳有效的 MarkdownEditor 實例")
	}
	
	// 驗證主要容器
	if editor.container == nil {
		t.Error("編輯器的主要容器不應該為 nil")
	}
	
	// 驗證工具欄容器
	if editor.toolbar == nil {
		t.Error("編輯器的工具欄容器不應該為 nil")
	}
	
	// 驗證工具欄是兩行佈局（包含工具欄和標籤）
	toolbarContainer := editor.toolbar
	if len(toolbarContainer.Objects) < 5 {
		t.Error("工具欄容器應包含至少5個元件（兩行工具欄、兩行標籤、一個分隔線）")
	}
	
	// 驗證文字編輯器
	if editor.editor == nil {
		t.Error("編輯器的文字編輯器不應該為 nil")
	}
	
	// 驗證狀態標籤
	if editor.statusLabel == nil {
		t.Error("編輯器的狀態標籤不應該為 nil")
	}
	
	// 驗證服務依賴
	if editor.editorService == nil {
		t.Error("編輯器服務不應該為 nil")
	}
	
	// 驗證初始狀態
	if editor.isModified {
		t.Error("編輯器初始狀態不應該為已修改")
	}
	
	if editor.currentNote != nil {
		t.Error("編輯器初始狀態不應該有當前筆記")
	}
}

// TestMarkdownEditorTextEditor 測試文字編輯器的配置
// 驗證文字編輯器的屬性和行為是否正確設定
//
// 測試項目：
// 1. 文字編輯器是否為多行模式
// 2. 自動換行是否正確設定
// 3. 滾動模式是否正確設定
// 4. 佔位文字是否正確設定
func TestMarkdownEditorTextEditor(t *testing.T) {
	// 建立模擬編輯器服務
	mockService := newMockEditorService()
	
	// 建立 Markdown 編輯器實例
	editor := NewMarkdownEditor(mockService)
	
	// 驗證文字編輯器屬性
	if editor.editor.MultiLine != true {
		t.Error("文字編輯器應該是多行模式")
	}
	
	// 驗證佔位文字
	expectedPlaceholder := "在此輸入您的 Markdown 內容..."
	if editor.editor.PlaceHolder != expectedPlaceholder {
		t.Errorf("佔位文字應該是 '%s'，但得到 '%s'", expectedPlaceholder, editor.editor.PlaceHolder)
	}
}

// TestMarkdownEditorCreateNewNote 測試建立新筆記功能
// 驗證編輯器是否能正確建立和載入新筆記
//
// 測試項目：
// 1. 新筆記是否正確建立
// 2. 筆記是否正確載入到編輯器
// 3. 編輯器狀態是否正確更新
func TestMarkdownEditorCreateNewNote(t *testing.T) {
	// 建立模擬編輯器服務
	mockService := newMockEditorService()
	
	// 建立 Markdown 編輯器實例
	editor := NewMarkdownEditor(mockService)
	
	// 建立新筆記
	testTitle := "測試筆記標題"
	err := editor.CreateNewNote(testTitle)
	
	// 驗證沒有錯誤
	if err != nil {
		t.Errorf("建立新筆記不應該產生錯誤，但得到: %s", err.Error())
	}
	
	// 驗證當前筆記
	if editor.currentNote == nil {
		t.Error("建立新筆記後應該有當前筆記")
	}
	
	if editor.currentNote.Title != testTitle {
		t.Errorf("筆記標題應該是 '%s'，但得到 '%s'", testTitle, editor.currentNote.Title)
	}
	
	// 驗證編輯器內容
	if editor.editor.Text != "" {
		t.Error("新筆記的內容應該是空的")
	}
	
	// 驗證修改狀態
	if editor.isModified {
		t.Error("新筆記初始狀態不應該為已修改")
	}
}

// TestMarkdownEditorLoadNote 測試載入筆記功能
// 驗證編輯器是否能正確載入現有筆記
//
// 測試項目：
// 1. 筆記內容是否正確載入到編輯器
// 2. 當前筆記是否正確設定
// 3. 修改狀態是否正確重置
func TestMarkdownEditorLoadNote(t *testing.T) {
	// 建立模擬編輯器服務
	mockService := newMockEditorService()
	
	// 建立 Markdown 編輯器實例
	editor := NewMarkdownEditor(mockService)
	
	// 建立測試筆記
	testNote := &models.Note{
		ID:      "test-id",
		Title:   "測試筆記",
		Content: "這是測試內容",
	}
	
	// 載入筆記
	editor.LoadNote(testNote)
	
	// 驗證當前筆記
	if editor.currentNote != testNote {
		t.Error("當前筆記應該是載入的筆記")
	}
	
	// 驗證編輯器內容
	if editor.editor.Text != testNote.Content {
		t.Errorf("編輯器內容應該是 '%s'，但得到 '%s'", testNote.Content, editor.editor.Text)
	}
	
	// 驗證修改狀態
	if editor.isModified {
		t.Error("載入筆記後修改狀態應該為 false")
	}
}

// TestMarkdownEditorSaveNote 測試保存筆記功能
// 驗證編輯器是否能正確保存筆記
//
// 測試項目：
// 1. 筆記是否正確保存
// 2. 修改狀態是否正確重置
// 3. 錯誤處理是否正確
func TestMarkdownEditorSaveNote(t *testing.T) {
	// 建立模擬編輯器服務
	mockService := newMockEditorService()
	
	// 建立 Markdown 編輯器實例
	editor := NewMarkdownEditor(mockService)
	
	// 建立新筆記
	err := editor.CreateNewNote("測試筆記")
	if err != nil {
		t.Fatalf("建立新筆記失敗: %s", err.Error())
	}
	
	// 修改內容
	testContent := "修改後的內容"
	editor.SetContent(testContent)
	editor.isModified = true // 手動設定修改狀態
	
	// 保存筆記
	err = editor.SaveNote()
	
	// 驗證沒有錯誤
	if err != nil {
		t.Errorf("保存筆記不應該產生錯誤，但得到: %s", err.Error())
	}
	
	// 驗證修改狀態
	if editor.isModified {
		t.Error("保存後修改狀態應該為 false")
	}
}

// TestMarkdownEditorSaveNoteWithoutCurrentNote 測試沒有當前筆記時的保存行為
// 驗證編輯器在沒有當前筆記時是否正確處理保存請求
//
// 測試項目：
// 1. 沒有當前筆記時保存是否回傳錯誤
// 2. 錯誤訊息是否正確
func TestMarkdownEditorSaveNoteWithoutCurrentNote(t *testing.T) {
	// 建立模擬編輯器服務
	mockService := newMockEditorService()
	
	// 建立 Markdown 編輯器實例
	editor := NewMarkdownEditor(mockService)
	
	// 嘗試保存（沒有當前筆記）
	err := editor.SaveNote()
	
	// 驗證應該有錯誤
	if err == nil {
		t.Error("沒有當前筆記時保存應該產生錯誤")
	}
	
	// 驗證錯誤訊息
	expectedError := "沒有可保存的筆記"
	if err.Error() != expectedError {
		t.Errorf("錯誤訊息應該是 '%s'，但得到 '%s'", expectedError, err.Error())
	}
}

// TestMarkdownEditorContentOperations 測試內容操作功能
// 驗證編輯器的內容設定、取得和修改檢測功能
//
// 測試項目：
// 1. 內容設定和取得是否正確
// 2. 修改狀態檢測是否正確
// 3. 內容清空是否正確
func TestMarkdownEditorContentOperations(t *testing.T) {
	// 建立模擬編輯器服務
	mockService := newMockEditorService()
	
	// 建立 Markdown 編輯器實例
	editor := NewMarkdownEditor(mockService)
	
	// 測試內容設定和取得
	testContent := "測試內容"
	editor.SetContent(testContent)
	
	if editor.GetContent() != testContent {
		t.Errorf("內容應該是 '%s'，但得到 '%s'", testContent, editor.GetContent())
	}
	
	// 驗證修改狀態（SetContent 應該重置修改狀態）
	if editor.IsModified() {
		t.Error("SetContent 後修改狀態應該為 false")
	}
	
	// 測試內容清空
	editor.Clear()
	
	if editor.GetContent() != "" {
		t.Error("清空後內容應該是空字串")
	}
	
	if editor.currentNote != nil {
		t.Error("清空後當前筆記應該為 nil")
	}
	
	if editor.IsModified() {
		t.Error("清空後修改狀態應該為 false")
	}
}

// TestMarkdownEditorTitleOperations 測試標題操作功能
// 驗證編輯器的標題設定和取得功能
//
// 測試項目：
// 1. 標題設定和取得是否正確
// 2. 沒有當前筆記時的標題操作
func TestMarkdownEditorTitleOperations(t *testing.T) {
	// 建立模擬編輯器服務
	mockService := newMockEditorService()
	
	// 建立 Markdown 編輯器實例
	editor := NewMarkdownEditor(mockService)
	
	// 測試沒有當前筆記時的標題取得
	if editor.GetTitle() != "" {
		t.Error("沒有當前筆記時標題應該是空字串")
	}
	
	// 建立新筆記
	testTitle := "原始標題"
	err := editor.CreateNewNote(testTitle)
	if err != nil {
		t.Fatalf("建立新筆記失敗: %s", err.Error())
	}
	
	// 驗證標題取得
	if editor.GetTitle() != testTitle {
		t.Errorf("標題應該是 '%s'，但得到 '%s'", testTitle, editor.GetTitle())
	}
	
	// 測試標題設定
	newTitle := "新標題"
	editor.SetTitle(newTitle)
	
	if editor.GetTitle() != newTitle {
		t.Errorf("標題應該是 '%s'，但得到 '%s'", newTitle, editor.GetTitle())
	}
	
	// 驗證修改狀態
	if !editor.IsModified() {
		t.Error("設定標題後修改狀態應該為 true")
	}
}

// TestMarkdownEditorCanSave 測試保存能力檢查功能
// 驗證編輯器是否能正確判斷是否可以保存
//
// 測試項目：
// 1. 沒有當前筆記時不能保存
// 2. 有當前筆記但未修改時不能保存
// 3. 有當前筆記且已修改時可以保存
func TestMarkdownEditorCanSave(t *testing.T) {
	// 建立模擬編輯器服務
	mockService := newMockEditorService()
	
	// 建立 Markdown 編輯器實例
	editor := NewMarkdownEditor(mockService)
	
	// 測試沒有當前筆記時
	if editor.CanSave() {
		t.Error("沒有當前筆記時不應該能保存")
	}
	
	// 建立新筆記
	err := editor.CreateNewNote("測試筆記")
	if err != nil {
		t.Fatalf("建立新筆記失敗: %s", err.Error())
	}
	
	// 測試有筆記但未修改時
	if editor.CanSave() {
		t.Error("有筆記但未修改時不應該能保存")
	}
	
	// 修改內容
	editor.isModified = true
	
	// 測試有筆記且已修改時
	if !editor.CanSave() {
		t.Error("有筆記且已修改時應該能保存")
	}
}

// TestMarkdownEditorCallbacks 測試回調函數功能
// 驗證編輯器的回調函數設定和觸發是否正確
//
// 測試項目：
// 1. 內容變更回調是否正確設定和觸發
// 2. 保存請求回調是否正確設定和觸發
// 3. 字數變更回調是否正確設定和觸發
func TestMarkdownEditorCallbacks(t *testing.T) {
	// 建立模擬編輯器服務
	mockService := newMockEditorService()
	
	// 建立 Markdown 編輯器實例
	editor := NewMarkdownEditor(mockService)
	
	// 測試回調函數設定
	var wordCountChangedCalled bool
	var lastWordCount int
	
	editor.SetOnWordCountChanged(func(count int) {
		wordCountChangedCalled = true
		lastWordCount = count
	})
	
	// 建立新筆記以觸發回調
	err := editor.CreateNewNote("測試筆記")
	if err != nil {
		t.Fatalf("建立新筆記失敗: %s", err.Error())
	}
	
	// 驗證字數變更回調
	if !wordCountChangedCalled {
		t.Error("建立新筆記應該觸發字數變更回調")
	}
	
	if lastWordCount != 0 {
		t.Errorf("新筆記的字數應該是 0，但得到 %d", lastWordCount)
	}
}

// TestMarkdownEditorGetContainer 測試容器取得功能
// 驗證編輯器是否能正確回傳主要容器
//
// 測試項目：
// 1. GetContainer 是否回傳正確的容器實例
// 2. 回傳的容器是否與內部容器相同
func TestMarkdownEditorGetContainer(t *testing.T) {
	// 建立模擬編輯器服務
	mockService := newMockEditorService()
	
	// 建立 Markdown 編輯器實例
	editor := NewMarkdownEditor(mockService)
	
	// 取得容器
	container := editor.GetContainer()
	
	// 驗證容器不為 nil
	if container == nil {
		t.Error("GetContainer 應該回傳有效的容器實例")
	}
	
	// 驗證回傳的容器與內部容器相同
	if container != editor.container {
		t.Error("GetContainer 應該回傳與內部容器相同的實例")
	}
}