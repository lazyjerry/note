// Package services 提供匯出和分享功能的服務實作
// 負責處理筆記的各種格式匯出、批量匯出和分享功能
package services

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"html/template"
	"io"
	"mac-notebook-app/internal/models"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

// exportServiceImpl 實作 ExportService 介面
// 提供完整的匯出和分享功能實作
type exportServiceImpl struct {
	// 依賴的服務
	editorService EditorService // 編輯器服務，用於處理筆記內容
	
	// 匯出任務管理
	exportTasks map[string]*ExportProgress // 匯出任務進度追蹤
	tasksMutex  sync.RWMutex               // 任務存取的讀寫鎖
	
	// Markdown 處理器
	markdownProcessor goldmark.Markdown // Goldmark Markdown 處理器
	
	// HTML 模板
	htmlTemplate *template.Template // HTML 匯出模板
}

// NewExportService 建立新的匯出服務實例
// 參數：editorService（編輯器服務實例）
// 回傳：ExportService 介面實例
//
// 執行流程：
// 1. 建立服務實例並初始化基本屬性
// 2. 設定 Markdown 處理器和擴展功能
// 3. 載入 HTML 匯出模板
// 4. 初始化匯出任務管理
func NewExportService(editorService EditorService) ExportService {
	service := &exportServiceImpl{
		editorService: editorService,
		exportTasks:   make(map[string]*ExportProgress),
	}
	
	// 設定 Markdown 處理器，啟用各種擴展功能
	service.markdownProcessor = goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,        // GitHub Flavored Markdown
			extension.Table,      // 表格支援
			extension.Strikethrough, // 刪除線
			extension.TaskList,   // 任務列表
			extension.Linkify,    // 自動連結
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(), // 自動生成標題 ID
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(), // 硬換行
			html.WithXHTML(),     // XHTML 相容
		),
	)
	
	// 載入 HTML 匯出模板
	service.loadHTMLTemplate()
	
	return service
}

// ExportToPDF 將筆記匯出為 PDF 格式
// 參數：note（要匯出的筆記）、outputPath（輸出檔案路徑）、options（匯出選項）
// 回傳：可能的錯誤
//
// 執行流程：
// 1. 驗證輸入參數和匯出路徑
// 2. 建立匯出任務並開始進度追蹤
// 3. 將 Markdown 內容轉換為 HTML
// 4. 使用 HTML 到 PDF 轉換器生成 PDF
// 5. 應用匯出選項（頁面設定、浮水印等）
// 6. 保存 PDF 檔案並更新進度
func (s *exportServiceImpl) ExportToPDF(note *models.Note, outputPath string, options *ExportOptions) error {
	// 驗證輸入參數
	if note == nil {
		return fmt.Errorf("筆記不能為空")
	}
	if outputPath == "" {
		return fmt.Errorf("輸出路徑不能為空")
	}
	
	// 驗證匯出路徑
	if valid, errMsg := s.ValidateExportPath(outputPath, ExportFormatPDF); !valid {
		return fmt.Errorf("無效的匯出路徑: %s", errMsg)
	}
	
	// 建立匯出任務
	exportID := s.generateExportID()
	progress := &ExportProgress{
		ExportID:    exportID,
		Progress:    0.0,
		Status:      ExportStatusInProgress,
		CurrentFile: note.Title,
	}
	
	s.tasksMutex.Lock()
	s.exportTasks[exportID] = progress
	s.tasksMutex.Unlock()
	
	// 設定預設選項
	if options == nil {
		options = s.getDefaultExportOptions()
	}
	
	// 更新進度：開始轉換 Markdown
	s.updateProgress(exportID, 0.2, "轉換 Markdown 內容...")
	
	// 將 Markdown 轉換為 HTML
	htmlContent, err := s.convertMarkdownToHTML(note.Content, options)
	if err != nil {
		s.updateProgressError(exportID, fmt.Errorf("Markdown 轉換失敗: %v", err))
		return err
	}
	
	// 更新進度：生成 PDF
	s.updateProgress(exportID, 0.6, "生成 PDF 檔案...")
	
	// 生成 PDF（這裡使用模擬實作，實際應用中需要整合 PDF 生成庫）
	err = s.generatePDFFromHTML(htmlContent, outputPath, options)
	if err != nil {
		s.updateProgressError(exportID, fmt.Errorf("PDF 生成失敗: %v", err))
		return err
	}
	
	// 更新進度：完成
	s.updateProgress(exportID, 1.0, "匯出完成")
	s.completeExport(exportID)
	
	return nil
}

// ExportToHTML 將筆記匯出為 HTML 格式
// 參數：note（要匯出的筆記）、outputPath（輸出檔案路徑）、options（匯出選項）
// 回傳：可能的錯誤
//
// 執行流程：
// 1. 驗證輸入參數和匯出路徑
// 2. 建立匯出任務並開始進度追蹤
// 3. 將 Markdown 內容轉換為 HTML
// 4. 應用 CSS 樣式和主題
// 5. 處理圖片和附件（如果需要）
// 6. 保存 HTML 檔案並更新進度
func (s *exportServiceImpl) ExportToHTML(note *models.Note, outputPath string, options *ExportOptions) error {
	// 驗證輸入參數
	if note == nil {
		return fmt.Errorf("筆記不能為空")
	}
	if outputPath == "" {
		return fmt.Errorf("輸出路徑不能為空")
	}
	
	// 驗證匯出路徑
	if valid, errMsg := s.ValidateExportPath(outputPath, ExportFormatHTML); !valid {
		return fmt.Errorf("無效的匯出路徑: %s", errMsg)
	}
	
	// 建立匯出任務
	exportID := s.generateExportID()
	progress := &ExportProgress{
		ExportID:    exportID,
		Progress:    0.0,
		Status:      ExportStatusInProgress,
		CurrentFile: note.Title,
	}
	
	s.tasksMutex.Lock()
	s.exportTasks[exportID] = progress
	s.tasksMutex.Unlock()
	
	// 設定預設選項
	if options == nil {
		options = s.getDefaultExportOptions()
	}
	
	// 更新進度：開始轉換
	s.updateProgress(exportID, 0.3, "轉換 Markdown 內容...")
	
	// 將 Markdown 轉換為 HTML
	htmlContent, err := s.convertMarkdownToHTML(note.Content, options)
	if err != nil {
		s.updateProgressError(exportID, fmt.Errorf("Markdown 轉換失敗: %v", err))
		return err
	}
	
	// 更新進度：應用樣式
	s.updateProgress(exportID, 0.6, "應用樣式和主題...")
	
	// 生成完整的 HTML 文件
	fullHTML, err := s.generateFullHTML(note, htmlContent, options)
	if err != nil {
		s.updateProgressError(exportID, fmt.Errorf("HTML 生成失敗: %v", err))
		return err
	}
	
	// 更新進度：保存檔案
	s.updateProgress(exportID, 0.9, "保存 HTML 檔案...")
	
	// 保存 HTML 檔案
	err = s.saveHTMLFile(fullHTML, outputPath)
	if err != nil {
		s.updateProgressError(exportID, fmt.Errorf("檔案保存失敗: %v", err))
		return err
	}
	
	// 更新進度：完成
	s.updateProgress(exportID, 1.0, "匯出完成")
	s.completeExport(exportID)
	
	return nil
}

// ExportToWord 將筆記匯出為 Word 文件格式
// 參數：note（要匯出的筆記）、outputPath（輸出檔案路徑）、options（匯出選項）
// 回傳：可能的錯誤
//
// 執行流程：
// 1. 驗證輸入參數和匯出路徑
// 2. 建立匯出任務並開始進度追蹤
// 3. 解析 Markdown 內容結構
// 4. 建立 Word 文件並設定格式
// 5. 轉換內容到 Word 格式
// 6. 保存 Word 檔案並更新進度
func (s *exportServiceImpl) ExportToWord(note *models.Note, outputPath string, options *ExportOptions) error {
	// 驗證輸入參數
	if note == nil {
		return fmt.Errorf("筆記不能為空")
	}
	if outputPath == "" {
		return fmt.Errorf("輸出路徑不能為空")
	}
	
	// 驗證匯出路徑
	if valid, errMsg := s.ValidateExportPath(outputPath, ExportFormatWord); !valid {
		return fmt.Errorf("無效的匯出路徑: %s", errMsg)
	}
	
	// 建立匯出任務
	exportID := s.generateExportID()
	progress := &ExportProgress{
		ExportID:    exportID,
		Progress:    0.0,
		Status:      ExportStatusInProgress,
		CurrentFile: note.Title,
	}
	
	s.tasksMutex.Lock()
	s.exportTasks[exportID] = progress
	s.tasksMutex.Unlock()
	
	// 設定預設選項
	if options == nil {
		options = s.getDefaultExportOptions()
	}
	
	// 更新進度：解析內容
	s.updateProgress(exportID, 0.2, "解析 Markdown 內容...")
	
	// 解析 Markdown 內容（這裡使用模擬實作）
	parsedContent, err := s.parseMarkdownForWord(note.Content)
	if err != nil {
		s.updateProgressError(exportID, fmt.Errorf("內容解析失敗: %v", err))
		return err
	}
	
	// 更新進度：建立文件
	s.updateProgress(exportID, 0.5, "建立 Word 文件...")
	
	// 生成 Word 文件（這裡使用模擬實作，實際應用中需要整合 Word 生成庫）
	err = s.generateWordDocument(note, parsedContent, outputPath, options)
	if err != nil {
		s.updateProgressError(exportID, fmt.Errorf("Word 文件生成失敗: %v", err))
		return err
	}
	
	// 更新進度：完成
	s.updateProgress(exportID, 1.0, "匯出完成")
	s.completeExport(exportID)
	
	return nil
}

// BatchExport 批量匯出多個筆記
// 參數：notes（要匯出的筆記陣列）、outputDir（輸出目錄）、format（匯出格式）、options（匯出選項）
// 回傳：匯出結果和可能的錯誤
//
// 執行流程：
// 1. 驗證輸入參數和輸出目錄
// 2. 建立批量匯出任務
// 3. 並行處理多個筆記匯出
// 4. 追蹤每個檔案的匯出進度
// 5. 統計匯出結果並回傳
func (s *exportServiceImpl) BatchExport(notes []*models.Note, outputDir string, format ExportFormat, options *ExportOptions) (*BatchExportResult, error) {
	startTime := time.Now()
	
	// 驗證輸入參數
	if len(notes) == 0 {
		return nil, fmt.Errorf("沒有要匯出的筆記")
	}
	if outputDir == "" {
		return nil, fmt.Errorf("輸出目錄不能為空")
	}
	
	// 確保輸出目錄存在
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, fmt.Errorf("建立輸出目錄失敗: %v", err)
	}
	
	// 建立批量匯出任務
	exportID := s.generateExportID()
	progress := &ExportProgress{
		ExportID:    exportID,
		Progress:    0.0,
		Status:      ExportStatusInProgress,
		CurrentFile: fmt.Sprintf("批量匯出 %d 個檔案", len(notes)),
	}
	
	s.tasksMutex.Lock()
	s.exportTasks[exportID] = progress
	s.tasksMutex.Unlock()
	
	// 設定預設選項
	if options == nil {
		options = s.getDefaultExportOptions()
	}
	
	// 初始化結果
	result := &BatchExportResult{
		TotalFiles:   len(notes),
		SuccessCount: 0,
		FailureCount: 0,
		FailedFiles:  make([]string, 0),
		OutputPath:   outputDir,
	}
	
	// 並行匯出處理
	const maxWorkers = 4 // 最大並行工作者數量
	semaphore := make(chan struct{}, maxWorkers)
	var wg sync.WaitGroup
	var resultMutex sync.Mutex
	
	for i, note := range notes {
		wg.Add(1)
		go func(index int, n *models.Note) {
			defer wg.Done()
			semaphore <- struct{}{} // 取得工作許可
			defer func() { <-semaphore }() // 釋放工作許可
			
			// 更新進度
			currentProgress := float64(index) / float64(len(notes))
			s.updateProgress(exportID, currentProgress, fmt.Sprintf("匯出: %s", n.Title))
			
			// 生成輸出檔案路徑
			outputPath := s.generateOutputPath(outputDir, n.Title, format)
			
			// 根據格式執行匯出
			var err error
			switch format {
			case ExportFormatPDF:
				err = s.ExportToPDF(n, outputPath, options)
			case ExportFormatHTML:
				err = s.ExportToHTML(n, outputPath, options)
			case ExportFormatWord:
				err = s.ExportToWord(n, outputPath, options)
			case ExportFormatMarkdown:
				err = s.exportToMarkdown(n, outputPath, options)
			default:
				err = fmt.Errorf("不支援的匯出格式: %s", format.String())
			}
			
			// 更新結果
			resultMutex.Lock()
			if err != nil {
				result.FailureCount++
				result.FailedFiles = append(result.FailedFiles, n.Title)
			} else {
				result.SuccessCount++
			}
			resultMutex.Unlock()
		}(i, note)
	}
	
	// 等待所有匯出完成
	wg.Wait()
	
	// 計算耗費時間
	result.ElapsedTime = time.Since(startTime)
	
	// 更新最終進度
	s.updateProgress(exportID, 1.0, "批量匯出完成")
	s.completeExport(exportID)
	
	return result, nil
}

// ShareNote 分享筆記
// 參數：note（要分享的筆記）、shareOptions（分享選項）
// 回傳：分享結果和可能的錯誤
//
// 執行流程：
// 1. 驗證輸入參數和分享選項
// 2. 根據分享類型執行不同的分享邏輯
// 3. 生成分享連結或執行分享動作
// 4. 回傳分享結果
func (s *exportServiceImpl) ShareNote(note *models.Note, shareOptions *ShareOptions) (*ShareResult, error) {
	// 驗證輸入參數
	if note == nil {
		return nil, fmt.Errorf("筆記不能為空")
	}
	if shareOptions == nil {
		return nil, fmt.Errorf("分享選項不能為空")
	}
	
	// 生成分享 ID
	shareID := s.generateShareID()
	
	result := &ShareResult{
		ShareID: shareID,
		Success: false,
	}
	
	// 根據分享類型執行不同邏輯
	switch shareOptions.ShareType {
	case ShareTypeLink:
		// 生成分享連結
		shareURL, err := s.generateShareLink(note, shareOptions)
		if err != nil {
			result.Message = fmt.Sprintf("生成分享連結失敗: %v", err)
			return result, err
		}
		result.ShareURL = shareURL
		result.ExpiryTime = shareOptions.ExpiryTime
		result.Success = true
		result.Message = "分享連結已生成"
		
	case ShareTypeEmail:
		// 電子郵件分享
		err := s.shareViaEmail(note, shareOptions)
		if err != nil {
			result.Message = fmt.Sprintf("電子郵件分享失敗: %v", err)
			return result, err
		}
		result.Success = true
		result.Message = "已透過電子郵件分享"
		
	case ShareTypeAirDrop:
		// AirDrop 分享
		err := s.shareViaAirDrop(note, shareOptions)
		if err != nil {
			result.Message = fmt.Sprintf("AirDrop 分享失敗: %v", err)
			return result, err
		}
		result.Success = true
		result.Message = "已透過 AirDrop 分享"
		
	case ShareTypeClipboard:
		// 複製到剪貼簿
		err := s.shareToClipboard(note, shareOptions)
		if err != nil {
			result.Message = fmt.Sprintf("複製到剪貼簿失敗: %v", err)
			return result, err
		}
		result.Success = true
		result.Message = "內容已複製到剪貼簿"
		
	default:
		return nil, fmt.Errorf("不支援的分享類型")
	}
	
	return result, nil
}

// GetSupportedFormats 取得支援的匯出格式列表
// 回傳：支援的匯出格式陣列
func (s *exportServiceImpl) GetSupportedFormats() []ExportFormat {
	return []ExportFormat{
		ExportFormatPDF,
		ExportFormatHTML,
		ExportFormatWord,
		ExportFormatMarkdown,
	}
}

// ValidateExportPath 驗證匯出路徑的有效性
// 參數：path（要驗證的路徑）、format（匯出格式）
// 回傳：路徑是否有效和可能的錯誤訊息
func (s *exportServiceImpl) ValidateExportPath(path string, format ExportFormat) (bool, string) {
	// 檢查路徑是否為空
	if path == "" {
		return false, "路徑不能為空"
	}
	
	// 檢查目錄是否存在
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return false, "目錄不存在"
	}
	
	// 檢查副檔名是否正確
	expectedExt := s.getFileExtension(format)
	actualExt := strings.ToLower(filepath.Ext(path))
	if actualExt != expectedExt {
		return false, fmt.Sprintf("副檔名應為 %s，實際為 %s", expectedExt, actualExt)
	}
	
	// 檢查是否有寫入權限
	testFile := filepath.Join(dir, ".write_test")
	if file, err := os.Create(testFile); err != nil {
		return false, "沒有寫入權限"
	} else {
		file.Close()
		os.Remove(testFile)
	}
	
	return true, ""
}

// GetExportProgress 取得匯出進度
// 參數：exportID（匯出任務 ID）
// 回傳：匯出進度資訊
func (s *exportServiceImpl) GetExportProgress(exportID string) *ExportProgress {
	s.tasksMutex.RLock()
	defer s.tasksMutex.RUnlock()
	
	if progress, exists := s.exportTasks[exportID]; exists {
		// 回傳進度的副本，避免併發修改
		return &ExportProgress{
			ExportID:      progress.ExportID,
			Progress:      progress.Progress,
			Status:        progress.Status,
			CurrentFile:   progress.CurrentFile,
			ElapsedTime:   progress.ElapsedTime,
			EstimatedTime: progress.EstimatedTime,
			Error:         progress.Error,
		}
	}
	
	return nil
}

// CancelExport 取消匯出任務
// 參數：exportID（匯出任務 ID）
// 回傳：是否成功取消
func (s *exportServiceImpl) CancelExport(exportID string) bool {
	s.tasksMutex.Lock()
	defer s.tasksMutex.Unlock()
	
	if progress, exists := s.exportTasks[exportID]; exists {
		if progress.Status == ExportStatusInProgress {
			progress.Status = ExportStatusCancelled
			return true
		}
	}
	
	return false
}

// 私有輔助方法

// generateExportID 生成唯一的匯出任務 ID
func (s *exportServiceImpl) generateExportID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// generateShareID 生成唯一的分享 ID
func (s *exportServiceImpl) generateShareID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// updateProgress 更新匯出進度
func (s *exportServiceImpl) updateProgress(exportID string, progress float64, currentFile string) {
	s.tasksMutex.Lock()
	defer s.tasksMutex.Unlock()
	
	if task, exists := s.exportTasks[exportID]; exists {
		task.Progress = progress
		task.CurrentFile = currentFile
		task.ElapsedTime = time.Since(time.Now().Add(-task.ElapsedTime))
	}
}

// updateProgressError 更新匯出進度錯誤
func (s *exportServiceImpl) updateProgressError(exportID string, err error) {
	s.tasksMutex.Lock()
	defer s.tasksMutex.Unlock()
	
	if task, exists := s.exportTasks[exportID]; exists {
		task.Status = ExportStatusFailed
		task.Error = err
	}
}

// completeExport 完成匯出任務
func (s *exportServiceImpl) completeExport(exportID string) {
	s.tasksMutex.Lock()
	defer s.tasksMutex.Unlock()
	
	if task, exists := s.exportTasks[exportID]; exists {
		task.Status = ExportStatusCompleted
	}
}

// getDefaultExportOptions 取得預設匯出選項
func (s *exportServiceImpl) getDefaultExportOptions() *ExportOptions {
	return &ExportOptions{
		IncludeMetadata:        false,
		IncludeTableOfContents: true,
		Theme:                  "default",
		FontSize:               13,
		PageSize:               "A4",
		Margins:                "2cm",
		IncludeImages:          true,
		ImageQuality:           80,
		WatermarkText:          "",
		HeaderText:             "",
		FooterText:             "",
	}
}

// getFileExtension 根據匯出格式取得檔案副檔名
func (s *exportServiceImpl) getFileExtension(format ExportFormat) string {
	switch format {
	case ExportFormatPDF:
		return ".pdf"
	case ExportFormatHTML:
		return ".html"
	case ExportFormatWord:
		return ".docx"
	case ExportFormatMarkdown:
		return ".md"
	default:
		return ".txt"
	}
}

// convertMarkdownToHTML 將 Markdown 內容轉換為 HTML
func (s *exportServiceImpl) convertMarkdownToHTML(content string, options *ExportOptions) (string, error) {
	var buf bytes.Buffer
	
	// 使用 Goldmark 轉換 Markdown 到 HTML
	if err := s.markdownProcessor.Convert([]byte(content), &buf); err != nil {
		return "", fmt.Errorf("Markdown 轉換失敗: %v", err)
	}
	
	return buf.String(), nil
}

// loadHTMLTemplate 載入 HTML 匯出模板
func (s *exportServiceImpl) loadHTMLTemplate() {
	templateContent := `<!DOCTYPE html>
<html lang="zh-TW">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}}</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'SF Pro Text', sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 800px;
            margin: 0 auto;
            padding: 20px;
        }
        h1, h2, h3, h4, h5, h6 {
            color: #2c3e50;
            margin-top: 2em;
            margin-bottom: 1em;
        }
        code {
            background-color: #f8f9fa;
            padding: 2px 4px;
            border-radius: 3px;
            font-family: 'SF Mono', Monaco, monospace;
        }
        pre {
            background-color: #f8f9fa;
            padding: 1em;
            border-radius: 5px;
            overflow-x: auto;
        }
        blockquote {
            border-left: 4px solid #3498db;
            margin: 0;
            padding-left: 1em;
            color: #7f8c8d;
        }
        table {
            border-collapse: collapse;
            width: 100%;
            margin: 1em 0;
        }
        th, td {
            border: 1px solid #ddd;
            padding: 8px;
            text-align: left;
        }
        th {
            background-color: #f2f2f2;
        }
    </style>
</head>
<body>
    {{if .IncludeMetadata}}
    <div class="metadata">
        <p><strong>建立時間：</strong>{{.CreatedAt}}</p>
        <p><strong>修改時間：</strong>{{.UpdatedAt}}</p>
    </div>
    {{end}}
    
    <h1>{{.Title}}</h1>
    
    {{.Content}}
    
    {{if .FooterText}}
    <footer>
        <p>{{.FooterText}}</p>
    </footer>
    {{end}}
</body>
</html>`
	
	var err error
	s.htmlTemplate, err = template.New("html_export").Parse(templateContent)
	if err != nil {
		// 如果模板載入失敗，使用簡單的預設模板
		s.htmlTemplate = template.Must(template.New("simple").Parse("<html><body><h1>{{.Title}}</h1>{{.Content}}</body></html>"))
	}
}

// 其他輔助方法的模擬實作（實際應用中需要完整實作）

// generatePDFFromHTML 從 HTML 內容生成 PDF 檔案
// 參數：htmlContent（HTML 內容）、outputPath（輸出路徑）、options（匯出選項）
// 回傳：可能的錯誤
//
// 執行流程：
// 1. 建立 PDF 文件結構
// 2. 設定頁面屬性（大小、邊距等）
// 3. 轉換 HTML 內容為 PDF 格式
// 4. 應用浮水印和頁首頁尾
// 5. 保存 PDF 檔案
func (s *exportServiceImpl) generatePDFFromHTML(htmlContent, outputPath string, options *ExportOptions) error {
	// 建立 PDF 內容結構
	pdfContent := s.buildPDFContent(htmlContent, options)
	
	// 這裡應該整合真正的 PDF 生成庫，如：
	// - github.com/jung-kurt/gofpdf (純 Go PDF 生成)
	// - github.com/chromedp/chromedp (使用 Chrome 引擎)
	// - wkhtmltopdf 命令行工具
	
	// 目前使用增強的模擬實作，包含更多 PDF 特性
	return s.writeToFile(outputPath, pdfContent)
}

// buildPDFContent 建立 PDF 內容結構
// 參數：htmlContent（HTML 內容）、options（匯出選項）
// 回傳：PDF 內容字串
func (s *exportServiceImpl) buildPDFContent(htmlContent string, options *ExportOptions) string {
	var content strings.Builder
	
	// PDF 檔案標頭
	content.WriteString("%PDF-1.4\n")
	content.WriteString("% 由 Mac 筆記本應用程式生成\n\n")
	
	// 文件屬性
	content.WriteString("1 0 obj\n")
	content.WriteString("<<\n")
	content.WriteString("/Type /Catalog\n")
	content.WriteString("/Pages 2 0 R\n")
	content.WriteString(">>\n")
	content.WriteString("endobj\n\n")
	
	// 頁面設定
	content.WriteString("2 0 obj\n")
	content.WriteString("<<\n")
	content.WriteString("/Type /Pages\n")
	content.WriteString("/Kids [3 0 R]\n")
	content.WriteString("/Count 1\n")
	content.WriteString(">>\n")
	content.WriteString("endobj\n\n")
	
	// 頁面內容
	content.WriteString("3 0 obj\n")
	content.WriteString("<<\n")
	content.WriteString("/Type /Page\n")
	content.WriteString("/Parent 2 0 R\n")
	
	// 設定頁面大小
	pageSize := s.getPageSizeForPDF(options.PageSize)
	content.WriteString(fmt.Sprintf("/MediaBox [0 0 %s]\n", pageSize))
	
	content.WriteString("/Contents 4 0 R\n")
	content.WriteString(">>\n")
	content.WriteString("endobj\n\n")
	
	// 內容流
	content.WriteString("4 0 obj\n")
	content.WriteString("<<\n")
	content.WriteString(fmt.Sprintf("/Length %d\n", len(htmlContent)))
	content.WriteString(">>\n")
	content.WriteString("stream\n")
	
	// 添加頁首
	if options.HeaderText != "" {
		content.WriteString(fmt.Sprintf("頁首: %s\n", options.HeaderText))
	}
	
	// 添加浮水印
	if options.WatermarkText != "" {
		content.WriteString(fmt.Sprintf("浮水印: %s\n", options.WatermarkText))
	}
	
	// 主要內容（簡化的 HTML 到文字轉換）
	textContent := s.htmlToText(htmlContent)
	content.WriteString(textContent)
	
	// 添加頁尾
	if options.FooterText != "" {
		content.WriteString(fmt.Sprintf("\n頁尾: %s", options.FooterText))
	}
	
	content.WriteString("\nendstream\n")
	content.WriteString("endobj\n\n")
	
	// PDF 結尾
	content.WriteString("xref\n")
	content.WriteString("0 5\n")
	content.WriteString("0000000000 65535 f \n")
	content.WriteString("0000000010 00000 n \n")
	content.WriteString("0000000079 00000 n \n")
	content.WriteString("0000000173 00000 n \n")
	content.WriteString("0000000301 00000 n \n")
	content.WriteString("trailer\n")
	content.WriteString("<<\n")
	content.WriteString("/Size 5\n")
	content.WriteString("/Root 1 0 R\n")
	content.WriteString(">>\n")
	content.WriteString("startxref\n")
	content.WriteString("492\n")
	content.WriteString("%%EOF\n")
	
	return content.String()
}

// getPageSizeForPDF 取得 PDF 頁面大小設定
// 參數：pageSize（頁面大小字串）
// 回傳：PDF 頁面大小規格
func (s *exportServiceImpl) getPageSizeForPDF(pageSize string) string {
	switch pageSize {
	case "A4":
		return "595 842"
	case "Letter":
		return "612 792"
	case "A3":
		return "842 1191"
	case "Legal":
		return "612 1008"
	default:
		return "595 842" // 預設 A4
	}
}

// htmlToText 將 HTML 內容轉換為純文字
// 參數：htmlContent（HTML 內容）
// 回傳：純文字內容
func (s *exportServiceImpl) htmlToText(htmlContent string) string {
	// 簡單的 HTML 標籤移除（實際應用中應使用專門的 HTML 解析器）
	text := htmlContent
	
	// 移除常見的 HTML 標籤
	htmlTags := []string{
		"<p>", "</p>", "<div>", "</div>", "<span>", "</span>",
		"<h1>", "</h1>", "<h2>", "</h2>", "<h3>", "</h3>",
		"<h4>", "</h4>", "<h5>", "</h5>", "<h6>", "</h6>",
		"<strong>", "</strong>", "<b>", "</b>", "<em>", "</em>",
		"<i>", "</i>", "<u>", "</u>", "<br>", "<br/>",
	}
	
	for _, tag := range htmlTags {
		text = strings.ReplaceAll(text, tag, "")
	}
	
	// 處理特殊字符
	text = strings.ReplaceAll(text, "&nbsp;", " ")
	text = strings.ReplaceAll(text, "&amp;", "&")
	text = strings.ReplaceAll(text, "&lt;", "<")
	text = strings.ReplaceAll(text, "&gt;", ">")
	text = strings.ReplaceAll(text, "&quot;", "\"")
	
	return text
}

func (s *exportServiceImpl) generateFullHTML(note *models.Note, htmlContent string, options *ExportOptions) (string, error) {
	data := struct {
		Title           string
		Content         string
		IncludeMetadata bool
		CreatedAt       string
		UpdatedAt       string
		FooterText      string
	}{
		Title:           note.Title,
		Content:         htmlContent,
		IncludeMetadata: options.IncludeMetadata,
		CreatedAt:       note.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:       note.UpdatedAt.Format("2006-01-02 15:04:05"),
		FooterText:      options.FooterText,
	}
	
	var buf bytes.Buffer
	if err := s.htmlTemplate.Execute(&buf, data); err != nil {
		return "", err
	}
	
	return buf.String(), nil
}

func (s *exportServiceImpl) saveHTMLFile(content, outputPath string) error {
	return s.writeToFile(outputPath, content)
}

func (s *exportServiceImpl) parseMarkdownForWord(content string) (interface{}, error) {
	// 這裡應該解析 Markdown 結構用於 Word 文件生成
	return content, nil
}

// generateWordDocument 生成 Word 文件
// 參數：note（筆記）、parsedContent（解析後的內容）、outputPath（輸出路徑）、options（匯出選項）
// 回傳：可能的錯誤
//
// 執行流程：
// 1. 建立 Word 文件結構
// 2. 設定文件屬性和樣式
// 3. 轉換 Markdown 內容為 Word 格式
// 4. 應用格式化和樣式
// 5. 保存 Word 檔案
func (s *exportServiceImpl) generateWordDocument(note *models.Note, parsedContent interface{}, outputPath string, options *ExportOptions) error {
	// 建立 Word 文件內容
	wordContent := s.buildWordContent(note, parsedContent, options)
	
	// 這裡應該整合真正的 Word 文件生成庫，如：
	// - github.com/unidoc/unioffice (商業授權)
	// - github.com/lukasjarosch/go-docx (開源)
	// - 或使用 Office Open XML 格式直接生成
	
	// 目前使用增強的模擬實作，生成 Office Open XML 結構
	return s.writeToFile(outputPath, wordContent)
}

// buildWordContent 建立 Word 文件內容
// 參數：note（筆記）、parsedContent（解析後的內容）、options（匯出選項）
// 回傳：Word 文件內容字串
func (s *exportServiceImpl) buildWordContent(note *models.Note, parsedContent interface{}, options *ExportOptions) string {
	var content strings.Builder
	
	// Office Open XML 文件結構開始
	content.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\" standalone=\"yes\"?>\n")
	content.WriteString("<w:document xmlns:w=\"http://schemas.openxmlformats.org/wordprocessingml/2006/main\">\n")
	content.WriteString("  <w:body>\n")
	
	// 文件標題
	if note.Title != "" {
		content.WriteString("    <w:p>\n")
		content.WriteString("      <w:pPr>\n")
		content.WriteString("        <w:pStyle w:val=\"Title\"/>\n")
		content.WriteString("      </w:pPr>\n")
		content.WriteString("      <w:r>\n")
		content.WriteString("        <w:rPr>\n")
		content.WriteString(fmt.Sprintf("          <w:sz w:val=\"%d\"/>\n", options.FontSize*2)) // Word 使用半點
		content.WriteString("          <w:b/>\n")
		content.WriteString("        </w:rPr>\n")
		content.WriteString(fmt.Sprintf("        <w:t>%s</w:t>\n", s.escapeXML(note.Title)))
		content.WriteString("      </w:r>\n")
		content.WriteString("    </w:p>\n")
	}
	
	// 元資料（如果啟用）
	if options.IncludeMetadata {
		content.WriteString("    <w:p>\n")
		content.WriteString("      <w:r>\n")
		content.WriteString("        <w:rPr><w:i/></w:rPr>\n")
		content.WriteString(fmt.Sprintf("        <w:t>建立時間：%s</w:t>\n", note.CreatedAt.Format("2006-01-02 15:04:05")))
		content.WriteString("      </w:r>\n")
		content.WriteString("    </w:p>\n")
		content.WriteString("    <w:p>\n")
		content.WriteString("      <w:r>\n")
		content.WriteString("        <w:rPr><w:i/></w:rPr>\n")
		content.WriteString(fmt.Sprintf("        <w:t>修改時間：%s</w:t>\n", note.UpdatedAt.Format("2006-01-02 15:04:05")))
		content.WriteString("      </w:r>\n")
		content.WriteString("    </w:p>\n")
	}
	
	// 主要內容
	content.WriteString("    <w:p>\n")
	content.WriteString("      <w:r>\n")
	content.WriteString("        <w:rPr>\n")
	content.WriteString(fmt.Sprintf("          <w:sz w:val=\"%d\"/>\n", options.FontSize*2))
	content.WriteString("        </w:rPr>\n")
	
	// 轉換 Markdown 內容為 Word 格式
	wordText := s.markdownToWordText(note.Content)
	content.WriteString(fmt.Sprintf("        <w:t xml:space=\"preserve\">%s</w:t>\n", s.escapeXML(wordText)))
	content.WriteString("      </w:r>\n")
	content.WriteString("    </w:p>\n")
	
	// 頁尾（如果有）
	if options.FooterText != "" {
		content.WriteString("    <w:p>\n")
		content.WriteString("      <w:pPr>\n")
		content.WriteString("        <w:jc w:val=\"center\"/>\n")
		content.WriteString("      </w:pPr>\n")
		content.WriteString("      <w:r>\n")
		content.WriteString("        <w:rPr><w:i/></w:rPr>\n")
		content.WriteString(fmt.Sprintf("        <w:t>%s</w:t>\n", s.escapeXML(options.FooterText)))
		content.WriteString("      </w:r>\n")
		content.WriteString("    </w:p>\n")
	}
	
	// Office Open XML 文件結構結束
	content.WriteString("  </w:body>\n")
	content.WriteString("</w:document>\n")
	
	return content.String()
}

// markdownToWordText 將 Markdown 內容轉換為適合 Word 的文字格式
// 參數：markdown（Markdown 內容）
// 回傳：轉換後的文字
func (s *exportServiceImpl) markdownToWordText(markdown string) string {
	text := markdown
	
	// 處理 Markdown 語法
	// 標題
	text = strings.ReplaceAll(text, "# ", "")
	text = strings.ReplaceAll(text, "## ", "")
	text = strings.ReplaceAll(text, "### ", "")
	
	// 粗體和斜體（簡化處理）
	text = strings.ReplaceAll(text, "**", "")
	text = strings.ReplaceAll(text, "*", "")
	
	// 程式碼區塊
	text = strings.ReplaceAll(text, "```", "")
	text = strings.ReplaceAll(text, "`", "")
	
	// 連結（簡化處理）
	// [文字](連結) -> 文字
	for {
		start := strings.Index(text, "[")
		if start == -1 {
			break
		}
		end := strings.Index(text[start:], ")")
		if end == -1 {
			break
		}
		
		linkText := text[start:][:end+1]
		if closeIdx := strings.Index(linkText, "]("); closeIdx != -1 {
			displayText := linkText[1:closeIdx]
			text = strings.Replace(text, linkText, displayText, 1)
		} else {
			break
		}
	}
	
	return text
}

// escapeXML 轉義 XML 特殊字符
// 參數：text（要轉義的文字）
// 回傳：轉義後的文字
func (s *exportServiceImpl) escapeXML(text string) string {
	text = strings.ReplaceAll(text, "&", "&amp;")
	text = strings.ReplaceAll(text, "<", "&lt;")
	text = strings.ReplaceAll(text, ">", "&gt;")
	text = strings.ReplaceAll(text, "\"", "&quot;")
	text = strings.ReplaceAll(text, "'", "&apos;")
	return text
}

func (s *exportServiceImpl) exportToMarkdown(note *models.Note, outputPath string, options *ExportOptions) error {
	content := note.Content
	if options.IncludeMetadata {
		metadata := fmt.Sprintf("---\ntitle: %s\ncreated: %s\nupdated: %s\n---\n\n",
			note.Title,
			note.CreatedAt.Format("2006-01-02 15:04:05"),
			note.UpdatedAt.Format("2006-01-02 15:04:05"))
		content = metadata + content
	}
	return s.writeToFile(outputPath, content)
}

func (s *exportServiceImpl) generateOutputPath(outputDir, title string, format ExportFormat) string {
	// 清理檔案名稱中的無效字符
	cleanTitle := strings.ReplaceAll(title, "/", "_")
	cleanTitle = strings.ReplaceAll(cleanTitle, "\\", "_")
	cleanTitle = strings.ReplaceAll(cleanTitle, ":", "_")
	
	ext := s.getFileExtension(format)
	return filepath.Join(outputDir, cleanTitle+ext)
}

func (s *exportServiceImpl) generateShareLink(note *models.Note, options *ShareOptions) (string, error) {
	// 這裡應該實作分享連結生成邏輯
	return fmt.Sprintf("https://share.notebook.app/%s", s.generateShareID()), nil
}

// shareViaEmail 透過電子郵件分享筆記
// 參數：note（要分享的筆記）、options（分享選項）
// 回傳：可能的錯誤
//
// 執行流程：
// 1. 建立電子郵件內容
// 2. 設定收件人和主旨
// 3. 附加筆記內容或檔案
// 4. 發送電子郵件
func (s *exportServiceImpl) shareViaEmail(note *models.Note, options *ShareOptions) error {
	// 建立電子郵件內容
	emailContent := s.buildEmailContent(note, options)
	
	// 這裡應該整合電子郵件發送功能，如：
	// - net/smtp 套件
	// - 第三方服務 API（SendGrid, Mailgun 等）
	// - macOS 系統郵件應用程式整合
	
	// 模擬發送電子郵件
	fmt.Printf("發送電子郵件到: %v\n", options.Recipients)
	fmt.Printf("主旨: 分享筆記 - %s\n", note.Title)
	fmt.Printf("內容: %s\n", emailContent)
	
	return nil
}

// shareViaAirDrop 透過 AirDrop 分享筆記
// 參數：note（要分享的筆記）、options（分享選項）
// 回傳：可能的錯誤
//
// 執行流程：
// 1. 建立臨時檔案
// 2. 調用 macOS AirDrop API
// 3. 顯示 AirDrop 選擇器
// 4. 清理臨時檔案
func (s *exportServiceImpl) shareViaAirDrop(note *models.Note, options *ShareOptions) error {
	// 建立臨時檔案
	tempFile, err := s.createTempFileForShare(note)
	if err != nil {
		return fmt.Errorf("建立臨時檔案失敗: %v", err)
	}
	defer os.Remove(tempFile) // 清理臨時檔案
	
	// 這裡應該整合 macOS AirDrop 功能，如：
	// - 使用 NSWorkspace 的 openFile:withApplication: 方法
	// - 調用系統分享服務
	// - 使用 CGO 調用 Objective-C 代碼
	
	// 模擬 AirDrop 分享
	fmt.Printf("透過 AirDrop 分享檔案: %s\n", tempFile)
	fmt.Printf("筆記標題: %s\n", note.Title)
	
	return nil
}

// shareToClipboard 將筆記內容複製到剪貼簿
// 參數：note（要分享的筆記）、options（分享選項）
// 回傳：可能的錯誤
//
// 執行流程：
// 1. 格式化筆記內容
// 2. 複製到系統剪貼簿
// 3. 提供用戶回饋
func (s *exportServiceImpl) shareToClipboard(note *models.Note, options *ShareOptions) error {
	// 格式化筆記內容
	clipboardContent := s.formatContentForClipboard(note)
	
	// 這裡應該整合剪貼簿操作，如：
	// - github.com/atotto/clipboard 套件
	// - 系統原生剪貼簿 API
	// - Fyne 的剪貼簿功能
	
	// 模擬複製到剪貼簿
	fmt.Printf("已複製到剪貼簿:\n%s\n", clipboardContent)
	
	return nil
}

// buildEmailContent 建立電子郵件內容
// 參數：note（筆記）、options（分享選項）
// 回傳：電子郵件內容字串
func (s *exportServiceImpl) buildEmailContent(note *models.Note, options *ShareOptions) string {
	var content strings.Builder
	
	content.WriteString(fmt.Sprintf("親愛的朋友，\n\n"))
	content.WriteString(fmt.Sprintf("我想與您分享一篇筆記：%s\n\n", note.Title))
	
	// 如果允許查看內容
	if options.AllowDownload {
		content.WriteString("筆記內容：\n")
		content.WriteString("=" + strings.Repeat("=", len(note.Title)) + "\n")
		content.WriteString(note.Content)
		content.WriteString("\n\n")
	}
	
	// 添加分享連結（如果有）
	if shareURL := s.generateShareURL(note, options); shareURL != "" {
		content.WriteString(fmt.Sprintf("您也可以透過以下連結查看：%s\n", shareURL))
		
		if !options.ExpiryTime.IsZero() {
			content.WriteString(fmt.Sprintf("連結將於 %s 過期\n", options.ExpiryTime.Format("2006-01-02 15:04:05")))
		}
	}
	
	content.WriteString("\n此郵件由 Mac 筆記本應用程式自動發送。")
	
	return content.String()
}

// createTempFileForShare 為分享建立臨時檔案
// 參數：note（筆記）
// 回傳：臨時檔案路徑和可能的錯誤
func (s *exportServiceImpl) createTempFileForShare(note *models.Note) (string, error) {
	// 建立臨時目錄
	tempDir := os.TempDir()
	
	// 生成安全的檔案名稱
	safeTitle := s.sanitizeFileName(note.Title)
	if safeTitle == "" {
		safeTitle = "筆記"
	}
	
	tempFile := filepath.Join(tempDir, safeTitle+".md")
	
	// 建立檔案內容
	content := fmt.Sprintf("# %s\n\n%s", note.Title, note.Content)
	
	// 寫入檔案
	err := s.writeToFile(tempFile, content)
	if err != nil {
		return "", err
	}
	
	return tempFile, nil
}

// formatContentForClipboard 格式化內容用於剪貼簿
// 參數：note（筆記）
// 回傳：格式化後的內容
func (s *exportServiceImpl) formatContentForClipboard(note *models.Note) string {
	var content strings.Builder
	
	// 添加標題
	if note.Title != "" {
		content.WriteString(note.Title)
		content.WriteString("\n")
		content.WriteString(strings.Repeat("=", len(note.Title)))
		content.WriteString("\n\n")
	}
	
	// 添加內容
	content.WriteString(note.Content)
	
	// 添加時間戳記
	content.WriteString("\n\n---\n")
	content.WriteString(fmt.Sprintf("分享時間: %s", time.Now().Format("2006-01-02 15:04:05")))
	
	return content.String()
}

// generateShareURL 生成分享連結
// 參數：note（筆記）、options（分享選項）
// 回傳：分享連結
func (s *exportServiceImpl) generateShareURL(note *models.Note, options *ShareOptions) string {
	// 這裡應該實作真正的分享連結生成邏輯
	// 包含安全性驗證、過期時間管理等
	baseURL := "https://share.notebook.app"
	shareID := s.generateShareID()
	
	return fmt.Sprintf("%s/note/%s", baseURL, shareID)
}

// sanitizeFileName 清理檔案名稱
// 參數：fileName（原始檔案名稱）
// 回傳：清理後的檔案名稱
func (s *exportServiceImpl) sanitizeFileName(fileName string) string {
	// 移除或替換無效字符
	invalidChars := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	result := fileName
	
	for _, char := range invalidChars {
		result = strings.ReplaceAll(result, char, "_")
	}
	
	// 限制長度
	if len(result) > 100 {
		result = result[:100]
	}
	
	return result
}

func (s *exportServiceImpl) writeToFile(path, content string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	
	_, err = io.WriteString(file, content)
	return err
}