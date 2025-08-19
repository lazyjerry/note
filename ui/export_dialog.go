// Package ui 提供匯出對話框的 UI 元件
// 負責處理筆記匯出的用戶介面，包含格式選擇、選項設定和進度顯示
package ui

import (
	"fmt"
	"mac-notebook-app/internal/models"
	"mac-notebook-app/internal/services"
	"path/filepath"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

// ExportDialog 匯出對話框結構
// 提供完整的筆記匯出功能介面
type ExportDialog struct {
	// UI 元件
	window       fyne.Window                // 父視窗
	dialog       *dialog.CustomDialog       // 自訂對話框
	content      *fyne.Container            // 對話框內容容器
	
	// 匯出選項 UI 元件
	formatSelect     *widget.Select         // 匯出格式選擇
	pathEntry       *widget.Entry          // 輸出路徑輸入框
	browseButton    *widget.Button         // 瀏覽按鈕
	
	// 進階選項
	includeMetadata *widget.Check          // 包含元資料選項
	includeTOC      *widget.Check          // 包含目錄選項
	themeSelect     *widget.Select         // 主題選擇
	fontSizeEntry   *widget.Entry          // 字體大小輸入
	pageSizeSelect  *widget.Select         // 頁面大小選擇
	marginsEntry    *widget.Entry          // 邊距設定
	includeImages   *widget.Check          // 包含圖片選項
	imageQuality    *widget.Slider         // 圖片品質滑桿
	watermarkEntry  *widget.Entry          // 浮水印文字
	headerEntry     *widget.Entry          // 頁首文字
	footerEntry     *widget.Entry          // 頁尾文字
	
	// 進度顯示
	progressBar     *widget.ProgressBar    // 進度條
	statusLabel     *widget.Label          // 狀態標籤
	
	// 按鈕
	exportButton    *widget.Button         // 匯出按鈕
	cancelButton    *widget.Button         // 取消按鈕
	
	// 服務和資料
	exportService   services.ExportService // 匯出服務
	note           *models.Note            // 要匯出的筆記
	exportID       string                  // 當前匯出任務 ID
	
	// 回調函數
	onExportCompleteCallback func(success bool, outputPath string) // 匯出完成回調
}

// NewExportDialog 建立新的匯出對話框
// 參數：window（父視窗）、exportService（匯出服務）、note（要匯出的筆記）
// 回傳：ExportDialog 實例
//
// 執行流程：
// 1. 初始化對話框結構和基本屬性
// 2. 建立所有 UI 元件
// 3. 設定預設值和事件處理
// 4. 組裝對話框佈局
func NewExportDialog(window fyne.Window, exportService services.ExportService, note *models.Note) *ExportDialog {
	dialog := &ExportDialog{
		window:        window,
		exportService: exportService,
		note:         note,
	}
	
	// 建立 UI 元件
	dialog.createUIComponents()
	
	// 設定預設值
	dialog.setDefaultValues()
	
	// 建立對話框佈局
	dialog.createLayout()
	
	return dialog
}

// Show 顯示匯出對話框
// 將對話框顯示給用戶並等待操作
func (d *ExportDialog) Show() {
	d.dialog.Show()
}

// Hide 隱藏匯出對話框
// 關閉對話框並清理資源
func (d *ExportDialog) Hide() {
	if d.dialog != nil {
		d.dialog.Hide()
	}
	
	// 如果有進行中的匯出任務，取消它
	if d.exportID != "" {
		d.exportService.CancelExport(d.exportID)
	}
}

// SetOnExportComplete 設定匯出完成回調函數
// 參數：callback（匯出完成時的回調函數）
func (d *ExportDialog) SetOnExportComplete(callback func(success bool, outputPath string)) {
	d.onExportCompleteCallback = callback
}

// createUIComponents 建立所有 UI 元件
// 初始化對話框中的所有控制項和輸入元件
func (d *ExportDialog) createUIComponents() {
	// 匯出格式選擇
	d.formatSelect = widget.NewSelect([]string{"PDF", "HTML", "Word", "Markdown"}, nil)
	d.formatSelect.SetSelected("PDF")
	d.formatSelect.OnChanged = d.onFormatChanged
	
	// 輸出路徑
	d.pathEntry = widget.NewEntry()
	d.pathEntry.SetPlaceHolder("選擇匯出檔案路徑...")
	
	// 瀏覽按鈕
	d.browseButton = widget.NewButton("瀏覽...", d.onBrowseClicked)
	
	// 進階選項
	d.includeMetadata = widget.NewCheck("包含元資料", nil)
	d.includeTOC = widget.NewCheck("包含目錄", nil)
	d.includeTOC.SetChecked(true)
	
	d.themeSelect = widget.NewSelect([]string{"預設", "淺色", "深色", "專業"}, nil)
	d.themeSelect.SetSelected("預設")
	
	d.fontSizeEntry = widget.NewEntry()
	d.fontSizeEntry.SetText("13")
	
	d.pageSizeSelect = widget.NewSelect([]string{"A4", "Letter", "A3", "Legal"}, nil)
	d.pageSizeSelect.SetSelected("A4")
	
	d.marginsEntry = widget.NewEntry()
	d.marginsEntry.SetText("2cm")
	
	d.includeImages = widget.NewCheck("包含圖片", nil)
	d.includeImages.SetChecked(true)
	
	d.imageQuality = widget.NewSlider(1, 100)
	d.imageQuality.SetValue(80)
	
	d.watermarkEntry = widget.NewEntry()
	d.watermarkEntry.SetPlaceHolder("浮水印文字（可選）")
	
	d.headerEntry = widget.NewEntry()
	d.headerEntry.SetPlaceHolder("頁首文字（可選）")
	
	d.footerEntry = widget.NewEntry()
	d.footerEntry.SetPlaceHolder("頁尾文字（可選）")
	
	// 進度顯示
	d.progressBar = widget.NewProgressBar()
	d.progressBar.Hide()
	
	d.statusLabel = widget.NewLabel("")
	d.statusLabel.Hide()
	
	// 按鈕
	d.exportButton = widget.NewButton("匯出", d.onExportClicked)
	d.exportButton.Importance = widget.HighImportance
	
	d.cancelButton = widget.NewButton("取消", d.onCancelClicked)
}

// setDefaultValues 設定預設值
// 根據筆記資訊和用戶偏好設定預設的匯出選項
func (d *ExportDialog) setDefaultValues() {
	// 根據筆記標題生成預設檔案名
	if d.note != nil {
		defaultName := d.note.Title
		if defaultName == "" {
			defaultName = "未命名筆記"
		}
		
		// 清理檔案名稱中的無效字符
		defaultName = d.sanitizeFileName(defaultName)
		
		// 設定預設路徑（用戶文件目錄）
		if homeURI := fyne.CurrentApp().Storage().RootURI(); homeURI != nil {
			defaultPath := filepath.Join(homeURI.Path(), "Documents", defaultName+".pdf")
			d.pathEntry.SetText(defaultPath)
		}
	}
}

// createLayout 建立對話框佈局
// 組織所有 UI 元件到適當的佈局容器中
func (d *ExportDialog) createLayout() {
	// 基本設定區域
	basicSection := container.NewVBox(
		widget.NewLabel("基本設定"),
		widget.NewSeparator(),
		container.NewGridWithColumns(2,
			widget.NewLabel("匯出格式:"),
			d.formatSelect,
		),
		container.NewBorder(nil, nil, nil, d.browseButton, d.pathEntry),
	)
	
	// 進階選項區域
	advancedSection := container.NewVBox(
		widget.NewLabel("進階選項"),
		widget.NewSeparator(),
		container.NewGridWithColumns(2,
			d.includeMetadata,
			d.includeTOC,
		),
		container.NewGridWithColumns(2,
			widget.NewLabel("主題:"),
			d.themeSelect,
			widget.NewLabel("字體大小:"),
			d.fontSizeEntry,
			widget.NewLabel("頁面大小:"),
			d.pageSizeSelect,
			widget.NewLabel("邊距:"),
			d.marginsEntry,
		),
		d.includeImages,
		container.NewBorder(nil, nil, widget.NewLabel("圖片品質:"), 
			widget.NewLabel(fmt.Sprintf("%.0f%%", d.imageQuality.Value)), d.imageQuality),
		d.watermarkEntry,
		d.headerEntry,
		d.footerEntry,
	)
	
	// 進度區域
	progressSection := container.NewVBox(
		d.progressBar,
		d.statusLabel,
	)
	
	// 按鈕區域
	buttonSection := container.NewBorder(nil, nil, nil, 
		container.NewHBox(d.exportButton, d.cancelButton))
	
	// 主要內容
	d.content = container.NewVBox(
		basicSection,
		widget.NewSeparator(),
		advancedSection,
		widget.NewSeparator(),
		progressSection,
		buttonSection,
	)
	
	// 建立自訂對話框
	d.dialog = dialog.NewCustom("匯出筆記", "關閉", d.content, d.window)
	d.dialog.Resize(fyne.NewSize(500, 600))
}

// 事件處理方法

// onFormatChanged 處理匯出格式變更事件
// 參數：format（選擇的格式）
func (d *ExportDialog) onFormatChanged(format string) {
	// 更新檔案路徑的副檔名
	currentPath := d.pathEntry.Text
	if currentPath != "" {
		dir := filepath.Dir(currentPath)
		baseName := d.getFileNameWithoutExt(filepath.Base(currentPath))
		
		var newExt string
		switch format {
		case "PDF":
			newExt = ".pdf"
		case "HTML":
			newExt = ".html"
		case "Word":
			newExt = ".docx"
		case "Markdown":
			newExt = ".md"
		default:
			newExt = ".pdf"
		}
		
		newPath := filepath.Join(dir, baseName+newExt)
		d.pathEntry.SetText(newPath)
	}
	
	// 根據格式啟用/禁用相關選項
	d.updateOptionsForFormat(format)
}

// onBrowseClicked 處理瀏覽按鈕點擊事件
func (d *ExportDialog) onBrowseClicked() {
	// 根據選擇的格式設定檔案過濾器
	format := d.formatSelect.Selected
	var fileFilter storage.FileFilter
	
	switch format {
	case "PDF":
		fileFilter = storage.NewExtensionFileFilter([]string{".pdf"})
	case "HTML":
		fileFilter = storage.NewExtensionFileFilter([]string{".html", ".htm"})
	case "Word":
		fileFilter = storage.NewExtensionFileFilter([]string{".docx"})
	case "Markdown":
		fileFilter = storage.NewExtensionFileFilter([]string{".md"})
	default:
		fileFilter = nil
	}
	
	// 顯示檔案保存對話框
	saveDialog := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil || writer == nil {
			return
		}
		
		// 設定選擇的路徑
		d.pathEntry.SetText(writer.URI().Path())
		writer.Close()
	}, d.window)
	
	if fileFilter != nil {
		saveDialog.SetFilter(fileFilter)
	}
	
	// 設定預設檔案名
	if d.note != nil && d.note.Title != "" {
		defaultName := d.sanitizeFileName(d.note.Title)
		saveDialog.SetFileName(defaultName)
	}
	
	saveDialog.Show()
}

// onExportClicked 處理匯出按鈕點擊事件
func (d *ExportDialog) onExportClicked() {
	// 驗證輸入
	if !d.validateInput() {
		return
	}
	
	// 禁用匯出按鈕，顯示進度
	d.exportButton.SetText("匯出中...")
	d.exportButton.Disable()
	d.progressBar.Show()
	d.statusLabel.Show()
	d.statusLabel.SetText("準備匯出...")
	
	// 建立匯出選項
	options := d.createExportOptions()
	
	// 取得匯出格式
	format := d.getExportFormat()
	outputPath := d.pathEntry.Text
	
	// 在背景執行匯出
	go d.performExport(format, outputPath, options)
}

// onCancelClicked 處理取消按鈕點擊事件
func (d *ExportDialog) onCancelClicked() {
	// 如果有進行中的匯出，取消它
	if d.exportID != "" {
		d.exportService.CancelExport(d.exportID)
	}
	
	d.Hide()
}

// 輔助方法

// validateInput 驗證用戶輸入
// 回傳：輸入是否有效
func (d *ExportDialog) validateInput() bool {
	// 檢查輸出路徑
	outputPath := d.pathEntry.Text
	if outputPath == "" {
		d.showError("請選擇輸出檔案路徑")
		return false
	}
	
	// 驗證路徑
	format := d.getExportFormat()
	if valid, errMsg := d.exportService.ValidateExportPath(outputPath, format); !valid {
		d.showError(fmt.Sprintf("無效的輸出路徑: %s", errMsg))
		return false
	}
	
	// 驗證字體大小
	if fontSize := d.fontSizeEntry.Text; fontSize != "" {
		if size, err := strconv.Atoi(fontSize); err != nil || size < 8 || size > 72 {
			d.showError("字體大小必須是 8-72 之間的數字")
			return false
		}
	}
	
	return true
}

// createExportOptions 建立匯出選項
// 回傳：匯出選項結構
func (d *ExportDialog) createExportOptions() *services.ExportOptions {
	fontSize, _ := strconv.Atoi(d.fontSizeEntry.Text)
	if fontSize == 0 {
		fontSize = 13
	}
	
	return &services.ExportOptions{
		IncludeMetadata:        d.includeMetadata.Checked,
		IncludeTableOfContents: d.includeTOC.Checked,
		Theme:                  d.themeSelect.Selected,
		FontSize:               fontSize,
		PageSize:               d.pageSizeSelect.Selected,
		Margins:                d.marginsEntry.Text,
		IncludeImages:          d.includeImages.Checked,
		ImageQuality:           int(d.imageQuality.Value),
		WatermarkText:          d.watermarkEntry.Text,
		HeaderText:             d.headerEntry.Text,
		FooterText:             d.footerEntry.Text,
	}
}

// getExportFormat 取得選擇的匯出格式
// 回傳：匯出格式列舉值
func (d *ExportDialog) getExportFormat() services.ExportFormat {
	switch d.formatSelect.Selected {
	case "PDF":
		return services.ExportFormatPDF
	case "HTML":
		return services.ExportFormatHTML
	case "Word":
		return services.ExportFormatWord
	case "Markdown":
		return services.ExportFormatMarkdown
	default:
		return services.ExportFormatPDF
	}
}

// performExport 執行匯出操作
// 參數：format（匯出格式）、outputPath（輸出路徑）、options（匯出選項）
func (d *ExportDialog) performExport(format services.ExportFormat, outputPath string, options *services.ExportOptions) {
	var err error
	
	// 根據格式執行匯出
	switch format {
	case services.ExportFormatPDF:
		err = d.exportService.ExportToPDF(d.note, outputPath, options)
	case services.ExportFormatHTML:
		err = d.exportService.ExportToHTML(d.note, outputPath, options)
	case services.ExportFormatWord:
		err = d.exportService.ExportToWord(d.note, outputPath, options)
	case services.ExportFormatMarkdown:
		// Markdown 匯出使用批量匯出功能
		notes := []*models.Note{d.note}
		result, batchErr := d.exportService.BatchExport(notes, filepath.Dir(outputPath), format, options)
		if batchErr != nil {
			err = batchErr
		} else if result.FailureCount > 0 {
			err = fmt.Errorf("匯出失敗")
		}
	}
	
	// 更新 UI（在主執行緒中）
	go func() {
		time.Sleep(100 * time.Millisecond) // 短暫延遲確保匯出完成
		d.onExportComplete(err == nil, outputPath, err)
	}()
}

// onExportComplete 處理匯出完成事件
// 參數：success（是否成功）、outputPath（輸出路徑）、err（錯誤資訊）
func (d *ExportDialog) onExportComplete(success bool, outputPath string, err error) {
	// 重置 UI 狀態
	d.exportButton.SetText("匯出")
	d.exportButton.Enable()
	d.progressBar.Hide()
	
	if success {
		d.statusLabel.SetText("匯出完成！")
		d.showSuccess(fmt.Sprintf("檔案已成功匯出到: %s", outputPath))
		
		// 呼叫回調函數
		if d.onExportCompleteCallback != nil {
			d.onExportCompleteCallback(true, outputPath)
		}
		
		// 延遲關閉對話框
		time.AfterFunc(2*time.Second, func() {
			d.Hide()
		})
	} else {
		d.statusLabel.SetText("匯出失敗")
		errorMsg := "匯出過程中發生錯誤"
		if err != nil {
			errorMsg = err.Error()
		}
		d.showError(errorMsg)
		
		// 呼叫回調函數
		if d.onExportCompleteCallback != nil {
			d.onExportCompleteCallback(false, "")
		}
	}
}

// updateOptionsForFormat 根據格式更新選項可用性
// 參數：format（選擇的格式）
func (d *ExportDialog) updateOptionsForFormat(format string) {
	// PDF 和 Word 支援所有選項
	isPDFOrWord := format == "PDF" || format == "Word"
	
	// 頁面相關選項只對 PDF 和 Word 有效
	d.pageSizeSelect.Enable()
	d.marginsEntry.Enable()
	d.watermarkEntry.Enable()
	d.headerEntry.Enable()
	d.footerEntry.Enable()
	
	if !isPDFOrWord {
		d.pageSizeSelect.Disable()
		d.marginsEntry.Disable()
		d.watermarkEntry.Disable()
		d.headerEntry.Disable()
		d.footerEntry.Disable()
	}
}

// sanitizeFileName 清理檔案名稱中的無效字符
// 參數：fileName（原始檔案名稱）
// 回傳：清理後的檔案名稱
func (d *ExportDialog) sanitizeFileName(fileName string) string {
	// 替換無效字符
	invalidChars := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	result := fileName
	
	for _, invalidChar := range invalidChars {
		for i := 0; i < len(result); i++ {
			if string(result[i]) == invalidChar {
				result = result[:i] + "_" + result[i+1:]
			}
		}
	}
	
	// 限制長度
	if len(result) > 100 {
		result = result[:100]
	}
	
	return result
}

// getFileNameWithoutExt 取得不含副檔名的檔案名稱
// 參數：fileName（完整檔案名稱）
// 回傳：不含副檔名的檔案名稱
func (d *ExportDialog) getFileNameWithoutExt(fileName string) string {
	ext := filepath.Ext(fileName)
	if ext != "" {
		return fileName[:len(fileName)-len(ext)]
	}
	return fileName
}

// showError 顯示錯誤訊息
// 參數：message（錯誤訊息）
func (d *ExportDialog) showError(message string) {
	dialog.ShowError(fmt.Errorf("%s", message), d.window)
}

// showSuccess 顯示成功訊息
// 參數：message（成功訊息）
func (d *ExportDialog) showSuccess(message string) {
	dialog.ShowInformation("匯出成功", message, d.window)
}