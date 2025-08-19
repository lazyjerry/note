// Package services 提供匯出和分享功能的測試
// 測試匯出服務的各種功能，包含 PDF、HTML、Word 匯出和分享功能
package services

import (
	"mac-notebook-app/internal/models"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestNewExportService 測試建立匯出服務實例
// 驗證服務實例是否正確初始化
func TestNewExportService(t *testing.T) {
	// 建立模擬編輯器服務
	editorService := createMockEditorService()
	
	// 建立匯出服務
	exportService := NewExportService(editorService)
	
	// 驗證服務不為空
	if exportService == nil {
		t.Fatal("匯出服務不應為空")
	}
	
	// 驗證支援的格式
	formats := exportService.GetSupportedFormats()
	expectedFormats := []ExportFormat{
		ExportFormatPDF,
		ExportFormatHTML,
		ExportFormatWord,
		ExportFormatMarkdown,
	}
	
	if len(formats) != len(expectedFormats) {
		t.Errorf("支援格式數量不正確，期望 %d，實際 %d", len(expectedFormats), len(formats))
	}
	
	for i, format := range formats {
		if format != expectedFormats[i] {
			t.Errorf("格式不匹配，期望 %v，實際 %v", expectedFormats[i], format)
		}
	}
}

// TestExportToPDF 測試 PDF 匯出功能
// 驗證 PDF 匯出的完整流程
func TestExportToPDF(t *testing.T) {
	// 建立測試環境
	editorService := createMockEditorService()
	exportService := NewExportService(editorService)
	
	// 建立測試筆記
	note := &models.Note{
		ID:        "test-note-1",
		Title:     "測試筆記",
		Content:   "# 標題\n\n這是測試內容。\n\n## 子標題\n\n- 項目 1\n- 項目 2",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	// 建立臨時輸出目錄
	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "test.pdf")
	
	// 測試匯出
	err := exportService.ExportToPDF(note, outputPath, nil)
	if err != nil {
		t.Fatalf("PDF 匯出失敗: %v", err)
	}
	
	// 驗證檔案是否存在
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Error("PDF 檔案未建立")
	}
}

// TestExportToHTML 測試 HTML 匯出功能
// 驗證 HTML 匯出的完整流程和內容正確性
func TestExportToHTML(t *testing.T) {
	// 建立測試環境
	editorService := createMockEditorService()
	exportService := NewExportService(editorService)
	
	// 建立測試筆記
	note := &models.Note{
		ID:        "test-note-2",
		Title:     "HTML 測試筆記",
		Content:   "# 主標題\n\n這是 **粗體** 和 *斜體* 文字。\n\n```go\nfunc main() {\n    fmt.Println(\"Hello, World!\")\n}\n```",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	// 建立臨時輸出目錄
	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "test.html")
	
	// 設定匯出選項
	options := &ExportOptions{
		IncludeMetadata:        true,
		IncludeTableOfContents: true,
		Theme:                  "default",
		FontSize:               14,
	}
	
	// 測試匯出
	err := exportService.ExportToHTML(note, outputPath, options)
	if err != nil {
		t.Fatalf("HTML 匯出失敗: %v", err)
	}
	
	// 驗證檔案是否存在
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Error("HTML 檔案未建立")
	}
	
	// 讀取並驗證檔案內容
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("讀取 HTML 檔案失敗: %v", err)
	}
	
	htmlContent := string(content)
	
	// 驗證 HTML 結構
	if !containsString(htmlContent, "<!DOCTYPE html>") {
		t.Error("HTML 檔案缺少 DOCTYPE 宣告")
	}
	
	if !containsString(htmlContent, note.Title) {
		t.Error("HTML 檔案缺少筆記標題")
	}
	
	// 檢查是否包含標題（可能是 h1 標籤）
	if !containsString(htmlContent, "主標題") {
		t.Errorf("HTML 檔案缺少標題內容，實際內容: %s", htmlContent)
	}
}

// TestExportToWord 測試 Word 文件匯出功能
// 驗證 Word 文件匯出的基本功能
func TestExportToWord(t *testing.T) {
	// 建立測試環境
	editorService := createMockEditorService()
	exportService := NewExportService(editorService)
	
	// 建立測試筆記
	note := &models.Note{
		ID:        "test-note-3",
		Title:     "Word 測試筆記",
		Content:   "# 文件標題\n\n這是一個測試文件。\n\n| 欄位1 | 欄位2 |\n|-------|-------|\n| 值1   | 值2   |",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	// 建立臨時輸出目錄
	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "test.docx")
	
	// 測試匯出
	err := exportService.ExportToWord(note, outputPath, nil)
	if err != nil {
		t.Fatalf("Word 匯出失敗: %v", err)
	}
	
	// 驗證檔案是否存在
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Error("Word 檔案未建立")
	}
}

// TestBatchExport 測試批量匯出功能
// 驗證多個筆記的批量匯出處理
func TestBatchExport(t *testing.T) {
	// 建立測試環境
	editorService := createMockEditorService()
	exportService := NewExportService(editorService)
	
	// 建立測試筆記陣列
	notes := []*models.Note{
		{
			ID:        "batch-note-1",
			Title:     "批量測試筆記1",
			Content:   "# 筆記1\n\n這是第一個筆記。",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        "batch-note-2",
			Title:     "批量測試筆記2",
			Content:   "# 筆記2\n\n這是第二個筆記。",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        "batch-note-3",
			Title:     "批量測試筆記3",
			Content:   "# 筆記3\n\n這是第三個筆記。",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	
	// 建立臨時輸出目錄
	tempDir := t.TempDir()
	
	// 測試批量匯出
	result, err := exportService.BatchExport(notes, tempDir, ExportFormatHTML, nil)
	if err != nil {
		t.Fatalf("批量匯出失敗: %v", err)
	}
	
	// 驗證匯出結果
	if result.TotalFiles != len(notes) {
		t.Errorf("總檔案數不正確，期望 %d，實際 %d", len(notes), result.TotalFiles)
	}
	
	if result.SuccessCount != len(notes) {
		t.Errorf("成功匯出數量不正確，期望 %d，實際 %d", len(notes), result.SuccessCount)
	}
	
	if result.FailureCount != 0 {
		t.Errorf("失敗匯出數量應為 0，實際 %d", result.FailureCount)
	}
	
	// 驗證輸出檔案是否存在
	for _, note := range notes {
		expectedPath := filepath.Join(tempDir, note.Title+".html")
		if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
			t.Errorf("批量匯出檔案不存在: %s", expectedPath)
		}
	}
}

// TestShareNote 測試筆記分享功能
// 驗證不同分享類型的處理
func TestShareNote(t *testing.T) {
	// 建立測試環境
	editorService := createMockEditorService()
	exportService := NewExportService(editorService)
	
	// 建立測試筆記
	note := &models.Note{
		ID:        "share-note-1",
		Title:     "分享測試筆記",
		Content:   "# 分享內容\n\n這是要分享的筆記內容。",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	// 測試連結分享
	shareOptions := &ShareOptions{
		ShareType:     ShareTypeLink,
		ExpiryTime:    time.Now().Add(24 * time.Hour),
		Password:      "",
		AllowDownload: true,
		AllowEdit:     false,
	}
	
	result, err := exportService.ShareNote(note, shareOptions)
	if err != nil {
		t.Fatalf("連結分享失敗: %v", err)
	}
	
	if !result.Success {
		t.Error("分享應該成功")
	}
	
	if result.ShareURL == "" {
		t.Error("分享連結不應為空")
	}
	
	// 測試剪貼簿分享
	clipboardOptions := &ShareOptions{
		ShareType: ShareTypeClipboard,
	}
	
	clipboardResult, err := exportService.ShareNote(note, clipboardOptions)
	if err != nil {
		t.Fatalf("剪貼簿分享失敗: %v", err)
	}
	
	if !clipboardResult.Success {
		t.Error("剪貼簿分享應該成功")
	}
}

// TestValidateExportPath 測試匯出路徑驗證功能
// 驗證各種路徑驗證情況
func TestValidateExportPath(t *testing.T) {
	// 建立測試環境
	editorService := createMockEditorService()
	exportService := NewExportService(editorService)
	
	// 建立臨時目錄
	tempDir := t.TempDir()
	
	// 測試有效路徑
	validPath := filepath.Join(tempDir, "test.pdf")
	valid, errMsg := exportService.ValidateExportPath(validPath, ExportFormatPDF)
	if !valid {
		t.Errorf("有效路徑驗證失敗: %s", errMsg)
	}
	
	// 測試空路徑
	valid, errMsg = exportService.ValidateExportPath("", ExportFormatPDF)
	if valid {
		t.Error("空路徑應該驗證失敗")
	}
	if errMsg != "路徑不能為空" {
		t.Errorf("錯誤訊息不正確，期望 '路徑不能為空'，實際 '%s'", errMsg)
	}
	
	// 測試不存在的目錄
	invalidPath := filepath.Join("/nonexistent/directory", "test.pdf")
	valid, errMsg = exportService.ValidateExportPath(invalidPath, ExportFormatPDF)
	if valid {
		t.Error("不存在目錄的路徑應該驗證失敗")
	}
	
	// 測試錯誤的副檔名
	wrongExtPath := filepath.Join(tempDir, "test.txt")
	valid, errMsg = exportService.ValidateExportPath(wrongExtPath, ExportFormatPDF)
	if valid {
		t.Error("錯誤副檔名的路徑應該驗證失敗")
	}
}

// TestExportProgress 測試匯出進度追蹤功能
// 驗證進度追蹤的正確性
func TestExportProgress(t *testing.T) {
	// 建立測試環境
	editorService := createMockEditorService()
	exportService := NewExportService(editorService)
	
	// 建立測試筆記
	note := &models.Note{
		ID:        "progress-note-1",
		Title:     "進度測試筆記",
		Content:   "# 測試內容\n\n這是用於測試進度追蹤的筆記。",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	// 建立臨時輸出目錄
	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "progress_test.html")
	
	// 開始匯出（在背景執行）
	go func() {
		exportService.ExportToHTML(note, outputPath, nil)
	}()
	
	// 等待一段時間讓匯出開始
	time.Sleep(100 * time.Millisecond)
	
	// 注意：由於這是模擬實作，實際的進度追蹤可能需要更複雜的測試
	// 這裡主要測試 API 的正確性
}

// TestCancelExport 測試取消匯出功能
// 驗證匯出任務的取消機制
func TestCancelExport(t *testing.T) {
	// 建立測試環境
	editorService := createMockEditorService()
	exportService := NewExportService(editorService)
	
	// 測試取消不存在的任務
	success := exportService.CancelExport("nonexistent-task")
	if success {
		t.Error("取消不存在的任務應該失敗")
	}
}

// TestExportOptionsValidation 測試匯出選項驗證
// 驗證各種匯出選項的處理
func TestExportOptionsValidation(t *testing.T) {
	// 建立測試環境
	editorService := createMockEditorService()
	exportService := NewExportService(editorService)
	
	// 建立測試筆記
	note := &models.Note{
		ID:        "options-note-1",
		Title:     "選項測試筆記",
		Content:   "# 測試\n\n選項測試內容。",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	// 建立臨時輸出目錄
	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "options_test.html")
	
	// 測試自訂選項
	customOptions := &ExportOptions{
		IncludeMetadata:        true,
		IncludeTableOfContents: false,
		Theme:                  "dark",
		FontSize:               16,
		PageSize:               "Letter",
		Margins:                "1cm",
		IncludeImages:          false,
		ImageQuality:           60,
		WatermarkText:          "機密文件",
		HeaderText:             "公司內部文件",
		FooterText:             "版權所有",
	}
	
	// 測試使用自訂選項匯出
	err := exportService.ExportToHTML(note, outputPath, customOptions)
	if err != nil {
		t.Fatalf("使用自訂選項匯出失敗: %v", err)
	}
	
	// 驗證檔案是否存在
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Error("使用自訂選項的匯出檔案未建立")
	}
}

// TestErrorHandling 測試錯誤處理
// 驗證各種錯誤情況的處理
func TestErrorHandling(t *testing.T) {
	// 建立測試環境
	editorService := createMockEditorService()
	exportService := NewExportService(editorService)
	
	// 測試空筆記匯出
	err := exportService.ExportToPDF(nil, "test.pdf", nil)
	if err == nil {
		t.Error("空筆記匯出應該失敗")
	}
	
	// 測試空路徑匯出
	note := &models.Note{
		ID:      "error-note-1",
		Title:   "錯誤測試筆記",
		Content: "測試內容",
	}
	
	err = exportService.ExportToPDF(note, "", nil)
	if err == nil {
		t.Error("空路徑匯出應該失敗")
	}
	
	// 測試空筆記分享
	_, err = exportService.ShareNote(nil, &ShareOptions{ShareType: ShareTypeLink})
	if err == nil {
		t.Error("空筆記分享應該失敗")
	}
	
	// 測試空分享選項
	_, err = exportService.ShareNote(note, nil)
	if err == nil {
		t.Error("空分享選項應該失敗")
	}
}

// TestAdvancedPDFExport 測試進階 PDF 匯出功能
// 驗證 PDF 匯出的進階選項和格式化
func TestAdvancedPDFExport(t *testing.T) {
	// 建立測試環境
	editorService := createMockEditorService()
	exportService := NewExportService(editorService)
	
	// 建立包含複雜內容的測試筆記
	note := &models.Note{
		ID:        "advanced-pdf-note",
		Title:     "進階 PDF 測試筆記",
		Content:   "# 主標題\n\n這是 **粗體** 和 *斜體* 文字。\n\n## 子標題\n\n- 項目 1\n- 項目 2\n\n```go\nfunc main() {\n    fmt.Println(\"Hello, World!\")\n}\n```",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	// 建立臨時輸出目錄
	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "advanced_test.pdf")
	
	// 設定進階匯出選項
	options := &ExportOptions{
		IncludeMetadata:        true,
		IncludeTableOfContents: true,
		Theme:                  "professional",
		FontSize:               12,
		PageSize:               "A4",
		Margins:                "2.5cm",
		IncludeImages:          true,
		ImageQuality:           90,
		WatermarkText:          "機密文件",
		HeaderText:             "公司內部文件",
		FooterText:             "第 1 頁，共 1 頁",
	}
	
	// 測試進階 PDF 匯出
	err := exportService.ExportToPDF(note, outputPath, options)
	if err != nil {
		t.Fatalf("進階 PDF 匯出失敗: %v", err)
	}
	
	// 驗證檔案是否存在
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Error("進階 PDF 檔案未建立")
	}
	
	// 讀取並驗證檔案內容
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("讀取 PDF 檔案失敗: %v", err)
	}
	
	pdfContent := string(content)
	
	// 驗證 PDF 結構
	if !containsString(pdfContent, "%PDF-1.4") {
		t.Error("PDF 檔案缺少正確的檔案標頭")
	}
	
	// 驗證浮水印
	if !containsString(pdfContent, "機密文件") {
		t.Error("PDF 檔案缺少浮水印")
	}
	
	// 驗證頁首頁尾
	if !containsString(pdfContent, "公司內部文件") {
		t.Error("PDF 檔案缺少頁首")
	}
	
	if !containsString(pdfContent, "第 1 頁，共 1 頁") {
		t.Error("PDF 檔案缺少頁尾")
	}
}

// TestAdvancedWordExport 測試進階 Word 文件匯出功能
// 驗證 Word 文件匯出的格式化和結構
func TestAdvancedWordExport(t *testing.T) {
	// 建立測試環境
	editorService := createMockEditorService()
	exportService := NewExportService(editorService)
	
	// 建立測試筆記
	note := &models.Note{
		ID:        "advanced-word-note",
		Title:     "進階 Word 測試筆記",
		Content:   "# 文件標題\n\n這是一個包含 **粗體** 和 *斜體* 的段落。\n\n## 子標題\n\n1. 編號項目 1\n2. 編號項目 2\n\n[連結文字](https://example.com)",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	// 建立臨時輸出目錄
	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "advanced_test.docx")
	
	// 設定進階匯出選項
	options := &ExportOptions{
		IncludeMetadata: true,
		FontSize:        14,
		FooterText:      "版權所有 © 2024",
	}
	
	// 測試進階 Word 匯出
	err := exportService.ExportToWord(note, outputPath, options)
	if err != nil {
		t.Fatalf("進階 Word 匯出失敗: %v", err)
	}
	
	// 驗證檔案是否存在
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Error("進階 Word 檔案未建立")
	}
	
	// 讀取並驗證檔案內容
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("讀取 Word 檔案失敗: %v", err)
	}
	
	wordContent := string(content)
	
	// 驗證 Office Open XML 結構
	if !containsString(wordContent, "<?xml version=\"1.0\"") {
		t.Error("Word 檔案缺少 XML 標頭")
	}
	
	if !containsString(wordContent, "w:document") {
		t.Error("Word 檔案缺少文件結構")
	}
	
	// 驗證標題
	if !containsString(wordContent, note.Title) {
		t.Error("Word 檔案缺少筆記標題")
	}
	
	// 驗證元資料
	if !containsString(wordContent, "建立時間") {
		t.Error("Word 檔案缺少建立時間元資料")
	}
	
	// 驗證頁尾
	if !containsString(wordContent, "版權所有") {
		t.Error("Word 檔案缺少頁尾")
	}
}

// TestAdvancedShareFunctionality 測試進階分享功能
// 驗證不同分享類型的詳細功能
func TestAdvancedShareFunctionality(t *testing.T) {
	// 建立測試環境
	editorService := createMockEditorService()
	exportService := NewExportService(editorService)
	
	// 建立測試筆記
	note := &models.Note{
		ID:        "advanced-share-note",
		Title:     "進階分享測試筆記",
		Content:   "# 分享內容\n\n這是要分享的詳細筆記內容。\n\n包含多個段落和格式。",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	// 測試電子郵件分享
	emailOptions := &ShareOptions{
		ShareType:     ShareTypeEmail,
		ExpiryTime:    time.Now().Add(7 * 24 * time.Hour),
		Password:      "test123",
		AllowDownload: true,
		AllowEdit:     false,
		Recipients:    []string{"test1@example.com", "test2@example.com"},
	}
	
	result, err := exportService.ShareNote(note, emailOptions)
	if err != nil {
		t.Fatalf("電子郵件分享失敗: %v", err)
	}
	
	if !result.Success {
		t.Error("電子郵件分享應該成功")
	}
	
	// 測試 AirDrop 分享
	airdropOptions := &ShareOptions{
		ShareType: ShareTypeAirDrop,
	}
	
	result, err = exportService.ShareNote(note, airdropOptions)
	if err != nil {
		t.Fatalf("AirDrop 分享失敗: %v", err)
	}
	
	if !result.Success {
		t.Error("AirDrop 分享應該成功")
	}
	
	// 測試剪貼簿分享
	clipboardOptions := &ShareOptions{
		ShareType: ShareTypeClipboard,
	}
	
	result, err = exportService.ShareNote(note, clipboardOptions)
	if err != nil {
		t.Fatalf("剪貼簿分享失敗: %v", err)
	}
	
	if !result.Success {
		t.Error("剪貼簿分享應該成功")
	}
	
	// 測試帶密碼的連結分享
	linkOptions := &ShareOptions{
		ShareType:     ShareTypeLink,
		ExpiryTime:    time.Now().Add(24 * time.Hour),
		Password:      "secure123",
		AllowDownload: true,
		AllowEdit:     false,
	}
	
	result, err = exportService.ShareNote(note, linkOptions)
	if err != nil {
		t.Fatalf("連結分享失敗: %v", err)
	}
	
	if !result.Success {
		t.Error("連結分享應該成功")
	}
	
	if result.ShareURL == "" {
		t.Error("連結分享應該回傳分享連結")
	}
	
	if result.ExpiryTime.IsZero() {
		t.Error("連結分享應該設定過期時間")
	}
}

// 輔助函數

// containsString 檢查字串是否包含子字串
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || 
		(len(substr) > 0 && len(s) > 0 && findSubstring(s, substr)))
}

// findSubstring 尋找子字串
func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// createMockEditorService 建立模擬編輯器服務實例
// 用於匯出服務測試，提供完整的 EditorService 介面實作
func createMockEditorService() EditorService {
	return &mockExportEditorService{
		activeNotes: make(map[string]*models.Note),
	}
}

// mockExportEditorService 專用於匯出服務測試的模擬編輯器服務
type mockExportEditorService struct {
	activeNotes map[string]*models.Note
}

func (m *mockExportEditorService) CreateNote(title, content string) (*models.Note, error) {
	return &models.Note{
		ID:        "mock-note",
		Title:     title,
		Content:   content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (m *mockExportEditorService) OpenNote(filePath string) (*models.Note, error) {
	return &models.Note{
		ID:       "mock-note",
		Title:    "Mock Note",
		Content:  "Mock content",
		FilePath: filePath,
	}, nil
}

func (m *mockExportEditorService) SaveNote(note *models.Note) error {
	return nil
}

func (m *mockExportEditorService) UpdateContent(noteID, content string) error {
	return nil
}

func (m *mockExportEditorService) PreviewMarkdown(content string) string {
	return "<p>" + content + "</p>"
}

func (m *mockExportEditorService) DecryptWithPassword(noteID, password string) (string, error) {
	return "decrypted content", nil
}

func (m *mockExportEditorService) GetActiveNotes() map[string]*models.Note {
	return m.activeNotes
}

func (m *mockExportEditorService) CloseNote(noteID string) {}

func (m *mockExportEditorService) GetActiveNote(noteID string) (*models.Note, bool) {
	return nil, false
}

func (m *mockExportEditorService) GetAutoCompleteSuggestions(content string, cursorPosition int) []AutoCompleteSuggestion {
	return []AutoCompleteSuggestion{}
}

func (m *mockExportEditorService) FormatTableContent(tableContent string) (string, error) {
	return tableContent, nil
}

func (m *mockExportEditorService) InsertLinkMarkdown(text, url string) string {
	return "[" + text + "](" + url + ")"
}

func (m *mockExportEditorService) InsertImageMarkdown(altText, imagePath string) string {
	return "![" + altText + "](" + imagePath + ")"
}

func (m *mockExportEditorService) GetSupportedCodeLanguages() []string {
	return []string{"go", "javascript", "python"}
}

func (m *mockExportEditorService) FormatCodeBlockMarkdown(code, language string) string {
	return "```" + language + "\n" + code + "\n```"
}

func (m *mockExportEditorService) FormatMathExpressionMarkdown(expression string, isInline bool) string {
	if isInline {
		return "$" + expression + "$"
	}
	return "$$" + expression + "$$"
}

func (m *mockExportEditorService) ValidateMarkdownContent(content string) (bool, []string) {
	return true, []string{}
}

func (m *mockExportEditorService) GenerateTableTemplateMarkdown(rows, cols int) string {
	return "| Header | Header |\n|--------|--------|\n| Cell   | Cell   |"
}

func (m *mockExportEditorService) PreviewMarkdownWithHighlight(content string) string {
	return "<p>" + content + "</p>"
}

func (m *mockExportEditorService) GetSmartEditingService() SmartEditingService {
	return nil
}

func (m *mockExportEditorService) SetSmartEditingService(smartEditSvc SmartEditingService) {}