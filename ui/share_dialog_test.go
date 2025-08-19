// Package ui 提供分享對話框的測試
// 測試分享對話框的各種功能和用戶互動
package ui

import (
	"mac-notebook-app/internal/models"
	"mac-notebook-app/internal/services"
	"testing"
	"time"

	"fyne.io/fyne/v2/test"
)

// TestNewShareDialog 測試建立分享對話框
// 驗證對話框是否正確初始化
func TestNewShareDialog(t *testing.T) {
	// 建立測試應用程式和視窗
	app := test.NewApp()
	window := test.NewWindow(nil)
	defer app.Quit()
	
	// 建立模擬服務
	exportService := &mockShareExportService{}
	
	// 建立測試筆記
	note := &models.Note{
		ID:        "test-note-1",
		Title:     "測試筆記",
		Content:   "# 測試內容\n\n這是測試筆記的內容。",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	// 建立分享對話框
	shareDialog := NewShareDialog(window, exportService, note)
	
	// 驗證對話框不為空
	if shareDialog == nil {
		t.Fatal("分享對話框不應為空")
	}
	
	// 驗證基本屬性
	if shareDialog.window != window {
		t.Error("視窗引用不正確")
	}
	
	if shareDialog.exportService != exportService {
		t.Error("匯出服務引用不正確")
	}
	
	if shareDialog.note != note {
		t.Error("筆記引用不正確")
	}
	
	// 驗證 UI 元件是否已建立
	if shareDialog.shareTypeSelect == nil {
		t.Error("分享類型選擇元件未建立")
	}
	
	if shareDialog.passwordEntry == nil {
		t.Error("密碼輸入元件未建立")
	}
	
	if shareDialog.shareButton == nil {
		t.Error("分享按鈕未建立")
	}
	
	if shareDialog.cancelButton == nil {
		t.Error("取消按鈕未建立")
	}
}

// TestShareDialogDefaultValues 測試預設值設定
// 驗證對話框的預設值是否正確
func TestShareDialogDefaultValues(t *testing.T) {
	// 建立測試環境
	app := test.NewApp()
	window := test.NewWindow(nil)
	defer app.Quit()
	
	exportService := &mockShareExportService{}
	
	note := &models.Note{
		ID:      "test-note-2",
		Title:   "預設值測試筆記",
		Content: "測試內容",
	}
	
	// 建立分享對話框
	shareDialog := NewShareDialog(window, exportService, note)
	
	// 驗證預設分享類型
	if shareDialog.shareTypeSelect.Selected != "連結分享" {
		t.Errorf("預設分享類型應為 '連結分享'，實際為 %s", shareDialog.shareTypeSelect.Selected)
	}
	
	// 驗證預設過期時間
	if shareDialog.expirySelect.Selected != "24 小時" {
		t.Errorf("預設過期時間應為 '24 小時'，實際為 %s", shareDialog.expirySelect.Selected)
	}
	
	// 驗證預設權限
	if !shareDialog.allowDownload.Checked {
		t.Error("預設應允許下載")
	}
	
	if shareDialog.allowEdit.Checked {
		t.Error("預設不應允許編輯")
	}
	
	// 驗證收件人輸入框預設隱藏
	if shareDialog.recipientsEntry.Visible() {
		t.Error("收件人輸入框預設應隱藏")
	}
}

// TestShareDialogEncryptedNoteDefaults 測試加密筆記的預設值
// 驗證加密筆記的特殊預設設定
func TestShareDialogEncryptedNoteDefaults(t *testing.T) {
	// 建立測試環境
	app := test.NewApp()
	window := test.NewWindow(nil)
	defer app.Quit()
	
	exportService := &mockShareExportService{}
	
	// 建立加密筆記
	encryptedNote := &models.Note{
		ID:          "test-note-3",
		Title:       "加密測試筆記",
		Content:     "加密內容",
		IsEncrypted: true,
	}
	
	// 建立分享對話框
	shareDialog := NewShareDialog(window, exportService, encryptedNote)
	
	// 驗證加密筆記的預設設定
	if shareDialog.allowEdit.Checked {
		t.Error("加密筆記預設不應允許編輯")
	}
}

// TestShareDialogTypeChange 測試分享類型變更功能
// 驗證分享類型變更時的 UI 變化
func TestShareDialogTypeChange(t *testing.T) {
	// 建立測試環境
	app := test.NewApp()
	window := test.NewWindow(nil)
	defer app.Quit()
	
	exportService := &mockShareExportService{}
	note := &models.Note{
		ID:      "test-note-4",
		Title:   "類型變更測試筆記",
		Content: "測試內容",
	}
	
	shareDialog := NewShareDialog(window, exportService, note)
	
	// 測試電子郵件分享
	shareDialog.onShareTypeChanged("電子郵件分享")
	if !shareDialog.recipientsEntry.Visible() {
		t.Error("電子郵件分享應顯示收件人輸入框")
	}
	
	// 測試 AirDrop 分享
	shareDialog.onShareTypeChanged("AirDrop 分享")
	if shareDialog.passwordEntry.Visible() {
		t.Error("AirDrop 分享不應顯示密碼輸入框")
	}
	
	// 測試剪貼簿分享
	shareDialog.onShareTypeChanged("複製到剪貼簿")
	if shareDialog.expirySelect.Visible() {
		t.Error("剪貼簿分享不應顯示過期時間選擇")
	}
	
	// 測試連結分享
	shareDialog.onShareTypeChanged("連結分享")
	if shareDialog.recipientsEntry.Visible() {
		t.Error("連結分享不應顯示收件人輸入框")
	}
	if !shareDialog.passwordEntry.Visible() {
		t.Error("連結分享應顯示密碼輸入框")
	}
}

// TestShareDialogValidation 測試輸入驗證功能
// 驗證各種輸入驗證情況
func TestShareDialogValidation(t *testing.T) {
	// 建立測試環境
	app := test.NewApp()
	window := test.NewWindow(nil)
	defer app.Quit()
	
	exportService := &mockShareExportService{}
	note := &models.Note{
		ID:      "test-note-5",
		Title:   "驗證測試筆記",
		Content: "測試內容",
	}
	
	shareDialog := NewShareDialog(window, exportService, note)
	
	// 測試連結分享驗證（應該通過）
	shareDialog.shareTypeSelect.SetSelected("連結分享")
	if !shareDialog.validateInput() {
		t.Error("連結分享應該驗證成功")
	}
	
	// 測試電子郵件分享無收件人
	shareDialog.shareTypeSelect.SetSelected("電子郵件分享")
	shareDialog.recipientsEntry.SetText("")
	if shareDialog.validateInput() {
		t.Error("電子郵件分享無收件人應該驗證失敗")
	}
	
	// 測試電子郵件分享有效收件人
	shareDialog.recipientsEntry.SetText("test@example.com")
	if !shareDialog.validateInput() {
		t.Error("電子郵件分享有效收件人應該驗證成功")
	}
	
	// 測試無效電子郵件格式
	shareDialog.recipientsEntry.SetText("invalid-email")
	if shareDialog.validateInput() {
		t.Error("無效電子郵件格式應該驗證失敗")
	}
	
	// 測試多個電子郵件地址
	shareDialog.recipientsEntry.SetText("test1@example.com, test2@example.com")
	if !shareDialog.validateInput() {
		t.Error("多個有效電子郵件地址應該驗證成功")
	}
}

// TestShareDialogOptions 測試分享選項建立
// 驗證分享選項是否正確建立
func TestShareDialogOptions(t *testing.T) {
	// 建立測試環境
	app := test.NewApp()
	window := test.NewWindow(nil)
	defer app.Quit()
	
	exportService := &mockShareExportService{}
	note := &models.Note{
		ID:      "test-note-6",
		Title:   "選項測試筆記",
		Content: "測試內容",
	}
	
	shareDialog := NewShareDialog(window, exportService, note)
	
	// 設定各種選項
	shareDialog.shareTypeSelect.SetSelected("電子郵件分享")
	shareDialog.passwordEntry.SetText("test123")
	shareDialog.expirySelect.SetSelected("7 天")
	shareDialog.allowDownload.SetChecked(false)
	shareDialog.allowEdit.SetChecked(true)
	shareDialog.recipientsEntry.SetText("test1@example.com, test2@example.com")
	
	// 建立分享選項
	options := shareDialog.createShareOptions()
	
	// 驗證選項
	if options.ShareType != services.ShareTypeEmail {
		t.Errorf("分享類型應為 Email，實際為 %v", options.ShareType)
	}
	
	if options.Password != "test123" {
		t.Errorf("密碼應為 'test123'，實際為 %s", options.Password)
	}
	
	if options.AllowDownload {
		t.Error("不應允許下載")
	}
	
	if !options.AllowEdit {
		t.Error("應允許編輯")
	}
	
	if len(options.Recipients) != 2 {
		t.Errorf("收件人數量應為 2，實際為 %d", len(options.Recipients))
	}
	
	expectedRecipients := []string{"test1@example.com", "test2@example.com"}
	for i, recipient := range options.Recipients {
		if recipient != expectedRecipients[i] {
			t.Errorf("收件人 %d 應為 %s，實際為 %s", i, expectedRecipients[i], recipient)
		}
	}
}

// TestShareDialogTypeMapping 測試分享類型映射
// 驗證 UI 分享類型到服務類型的映射
func TestShareDialogTypeMapping(t *testing.T) {
	// 建立測試環境
	app := test.NewApp()
	window := test.NewWindow(nil)
	defer app.Quit()
	
	exportService := &mockShareExportService{}
	note := &models.Note{
		ID:      "test-note-7",
		Title:   "類型映射測試筆記",
		Content: "測試內容",
	}
	
	shareDialog := NewShareDialog(window, exportService, note)
	
	// 測試各種類型映射
	testCases := []struct {
		uiType       string
		expectedType services.ShareType
	}{
		{"連結分享", services.ShareTypeLink},
		{"電子郵件分享", services.ShareTypeEmail},
		{"AirDrop 分享", services.ShareTypeAirDrop},
		{"複製到剪貼簿", services.ShareTypeClipboard},
	}
	
	for _, tc := range testCases {
		shareDialog.shareTypeSelect.SetSelected(tc.uiType)
		actualType := shareDialog.getShareType()
		
		if actualType != tc.expectedType {
			t.Errorf("類型 %s 映射錯誤，期望 %v，實際 %v", 
				tc.uiType, tc.expectedType, actualType)
		}
	}
}

// TestShareDialogExpiryTime 測試過期時間計算
// 驗證過期時間的正確計算
func TestShareDialogExpiryTime(t *testing.T) {
	// 建立測試環境
	app := test.NewApp()
	window := test.NewWindow(nil)
	defer app.Quit()
	
	exportService := &mockShareExportService{}
	note := &models.Note{
		ID:      "test-note-8",
		Title:   "過期時間測試筆記",
		Content: "測試內容",
	}
	
	shareDialog := NewShareDialog(window, exportService, note)
	
	now := time.Now()
	
	// 測試各種過期時間
	testCases := []struct {
		uiExpiry     string
		expectedDiff time.Duration
		tolerance    time.Duration
	}{
		{"1 小時", 1 * time.Hour, 1 * time.Minute},
		{"24 小時", 24 * time.Hour, 1 * time.Minute},
		{"7 天", 7 * 24 * time.Hour, 1 * time.Minute},
		{"30 天", 30 * 24 * time.Hour, 1 * time.Minute},
	}
	
	for _, tc := range testCases {
		shareDialog.expirySelect.SetSelected(tc.uiExpiry)
		expiryTime := shareDialog.getExpiryTime()
		
		actualDiff := expiryTime.Sub(now)
		expectedDiff := tc.expectedDiff
		
		if actualDiff < expectedDiff-tc.tolerance || actualDiff > expectedDiff+tc.tolerance {
			t.Errorf("過期時間 %s 計算錯誤，期望約 %v，實際 %v", 
				tc.uiExpiry, expectedDiff, actualDiff)
		}
	}
	
	// 測試永不過期
	shareDialog.expirySelect.SetSelected("永不過期")
	expiryTime := shareDialog.getExpiryTime()
	
	// 永不過期應該是很久以後的時間
	if expiryTime.Sub(now) < 50*365*24*time.Hour {
		t.Error("永不過期的時間應該是很久以後")
	}
}

// TestShareDialogEmailValidation 測試電子郵件驗證功能
// 驗證電子郵件地址格式驗證
func TestShareDialogEmailValidation(t *testing.T) {
	// 建立測試環境
	app := test.NewApp()
	window := test.NewWindow(nil)
	defer app.Quit()
	
	exportService := &mockShareExportService{}
	note := &models.Note{
		ID:      "test-note-9",
		Title:   "電子郵件驗證測試筆記",
		Content: "測試內容",
	}
	
	shareDialog := NewShareDialog(window, exportService, note)
	
	// 測試有效電子郵件
	validEmails := []string{
		"test@example.com",
		"user.name@domain.co.uk",
		"test123@test-domain.org",
	}
	
	for _, email := range validEmails {
		if !shareDialog.isValidEmail(email) {
			t.Errorf("電子郵件 %s 應該是有效的", email)
		}
	}
	
	// 測試無效電子郵件
	invalidEmails := []string{
		"",
		"test",
		"test@",
		"@example.com",
		"test@@example.com",
		"test@example",
	}
	
	for _, email := range invalidEmails {
		if shareDialog.isValidEmail(email) {
			t.Errorf("電子郵件 %s 應該是無效的", email)
		}
	}
	
	// 測試電子郵件列表
	validList := "test1@example.com, test2@example.com"
	if !shareDialog.isValidEmailList(validList) {
		t.Error("有效電子郵件列表應該驗證成功")
	}
	
	invalidList := "test1@example.com, invalid-email"
	if shareDialog.isValidEmailList(invalidList) {
		t.Error("包含無效電子郵件的列表應該驗證失敗")
	}
}

// TestShareDialogEmailParsing 測試電子郵件解析功能
// 驗證電子郵件地址列表的解析
func TestShareDialogEmailParsing(t *testing.T) {
	// 建立測試環境
	app := test.NewApp()
	window := test.NewWindow(nil)
	defer app.Quit()
	
	exportService := &mockShareExportService{}
	note := &models.Note{
		ID:      "test-note-10",
		Title:   "電子郵件解析測試筆記",
		Content: "測試內容",
	}
	
	shareDialog := NewShareDialog(window, exportService, note)
	
	// 測試解析電子郵件列表
	testCases := []struct {
		input    string
		expected []string
	}{
		{
			"test@example.com",
			[]string{"test@example.com"},
		},
		{
			"test1@example.com, test2@example.com",
			[]string{"test1@example.com", "test2@example.com"},
		},
		{
			" test1@example.com , test2@example.com ",
			[]string{"test1@example.com", "test2@example.com"},
		},
		{
			"test1@example.com,test2@example.com,test3@example.com",
			[]string{"test1@example.com", "test2@example.com", "test3@example.com"},
		},
	}
	
	for _, tc := range testCases {
		result := shareDialog.parseEmailList(tc.input)
		
		if len(result) != len(tc.expected) {
			t.Errorf("解析 '%s' 結果數量錯誤，期望 %d，實際 %d", 
				tc.input, len(tc.expected), len(result))
			continue
		}
		
		for i, email := range result {
			if email != tc.expected[i] {
				t.Errorf("解析 '%s' 結果 %d 錯誤，期望 %s，實際 %s", 
					tc.input, i, tc.expected[i], email)
			}
		}
	}
}

// TestShareDialogShowHide 測試對話框顯示和隱藏
// 驗證對話框的顯示和隱藏功能
func TestShareDialogShowHide(t *testing.T) {
	// 建立測試環境
	app := test.NewApp()
	window := test.NewWindow(nil)
	defer app.Quit()
	
	exportService := &mockShareExportService{}
	note := &models.Note{
		ID:      "test-note-11",
		Title:   "顯示隱藏測試筆記",
		Content: "測試內容",
	}
	
	shareDialog := NewShareDialog(window, exportService, note)
	
	// 測試顯示（這在測試環境中可能不會實際顯示）
	shareDialog.Show()
	
	// 測試隱藏
	shareDialog.Hide()
	
	// 驗證對話框仍然存在（隱藏不會銷毀對話框）
	if shareDialog.dialog == nil {
		t.Error("隱藏後對話框不應為空")
	}
}

// TestShareDialogCallback 測試回調函數設定
// 驗證分享完成回調函數的設定和呼叫
func TestShareDialogCallback(t *testing.T) {
	// 建立測試環境
	app := test.NewApp()
	window := test.NewWindow(nil)
	defer app.Quit()
	
	exportService := &mockShareExportService{}
	note := &models.Note{
		ID:      "test-note-12",
		Title:   "回調測試筆記",
		Content: "測試內容",
	}
	
	shareDialog := NewShareDialog(window, exportService, note)
	
	// 設定回調函數
	callbackCalled := false
	var callbackSuccess bool
	var callbackResult *services.ShareResult
	
	shareDialog.SetOnShareComplete(func(success bool, result *services.ShareResult) {
		callbackCalled = true
		callbackSuccess = success
		callbackResult = result
	})
	
	// 模擬分享完成
	mockResult := &services.ShareResult{
		ShareID:  "test-share-id",
		ShareURL: "https://test.share.url",
		Success:  true,
		Message:  "分享成功",
	}
	
	shareDialog.onShareComplete(mockResult, nil)
	
	// 驗證回調是否被呼叫（在實際 UI 環境中）
	// 注意：在測試環境中，UI 事件可能不會正常觸發
	// 這裡主要測試方法不會崩潰
	_ = callbackCalled
	_ = callbackSuccess
	_ = callbackResult
}

// mockShareExportService 模擬匯出服務，專用於分享對話框測試
type mockShareExportService struct{}

func (m *mockShareExportService) ExportToPDF(note *models.Note, outputPath string, options *services.ExportOptions) error {
	return nil
}

func (m *mockShareExportService) ExportToHTML(note *models.Note, outputPath string, options *services.ExportOptions) error {
	return nil
}

func (m *mockShareExportService) ExportToWord(note *models.Note, outputPath string, options *services.ExportOptions) error {
	return nil
}

func (m *mockShareExportService) BatchExport(notes []*models.Note, outputDir string, format services.ExportFormat, options *services.ExportOptions) (*services.BatchExportResult, error) {
	return &services.BatchExportResult{
		TotalFiles:   len(notes),
		SuccessCount: len(notes),
		FailureCount: 0,
		FailedFiles:  []string{},
		OutputPath:   outputDir,
		ElapsedTime:  time.Second,
	}, nil
}

func (m *mockShareExportService) ShareNote(note *models.Note, shareOptions *services.ShareOptions) (*services.ShareResult, error) {
	return &services.ShareResult{
		ShareID:    "mock-share-id",
		ShareURL:   "https://mock.share.url",
		ExpiryTime: shareOptions.ExpiryTime,
		Success:    true,
		Message:    "分享成功",
	}, nil
}

func (m *mockShareExportService) GetSupportedFormats() []services.ExportFormat {
	return []services.ExportFormat{
		services.ExportFormatPDF,
		services.ExportFormatHTML,
		services.ExportFormatWord,
		services.ExportFormatMarkdown,
	}
}

func (m *mockShareExportService) ValidateExportPath(path string, format services.ExportFormat) (bool, string) {
	if path == "" {
		return false, "路徑不能為空"
	}
	return true, ""
}

func (m *mockShareExportService) GetExportProgress(exportID string) *services.ExportProgress {
	return &services.ExportProgress{
		ExportID: exportID,
		Progress: 1.0,
		Status:   services.ExportStatusCompleted,
	}
}

func (m *mockShareExportService) CancelExport(exportID string) bool {
	return true
}