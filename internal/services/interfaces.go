// Package services 定義了應用程式的業務邏輯服務介面
// 這些介面將業務邏輯與具體實作分離，提高程式碼的可測試性和可維護性
package services

import (
	"mac-notebook-app/internal/models" // 引入資料模型
	"time"                             // 時間處理套件
)

// EditorService 定義筆記編輯操作的介面
// 負責處理筆記的建立、開啟、保存、更新和預覽等核心功能
type EditorService interface {
	// CreateNote 建立新的筆記
	// 參數：title（標題）、content（內容）
	// 回傳：筆記實例和可能的錯誤
	CreateNote(title, content string) (*models.Note, error)
	
	// OpenNote 從檔案路徑開啟筆記
	// 參數：filePath（檔案路徑）
	// 回傳：筆記實例和可能的錯誤
	OpenNote(filePath string) (*models.Note, error)
	
	// SaveNote 保存筆記到檔案系統
	// 參數：note（要保存的筆記）
	// 回傳：可能的錯誤
	SaveNote(note *models.Note) error
	
	// UpdateContent 更新指定筆記的內容
	// 參數：noteID（筆記 ID）、content（新內容）
	// 回傳：可能的錯誤
	UpdateContent(noteID, content string) error
	
	// PreviewMarkdown 將 Markdown 內容轉換為 HTML 預覽
	// 參數：content（Markdown 內容）
	// 回傳：HTML 字串
	PreviewMarkdown(content string) string
	
	// DecryptWithPassword 使用密碼解密筆記內容
	// 參數：noteID（筆記 ID）、password（解密密碼）
	// 回傳：解密後的內容和可能的錯誤
	DecryptWithPassword(noteID, password string) (string, error)
	
	// GetActiveNotes 取得所有活躍筆記
	// 回傳：活躍筆記的映射表
	GetActiveNotes() map[string]*models.Note
	
	// CloseNote 關閉指定的筆記
	// 參數：noteID（筆記 ID）
	CloseNote(noteID string)
	
	// GetActiveNote 取得指定的活躍筆記
	// 參數：noteID（筆記 ID）
	// 回傳：筆記實例和是否存在
	GetActiveNote(noteID string) (*models.Note, bool)
	
	// 智慧編輯功能
	// GetAutoCompleteSuggestions 取得自動完成建議
	// 參數：content（當前內容）、cursorPosition（游標位置）
	// 回傳：自動完成建議陣列
	GetAutoCompleteSuggestions(content string, cursorPosition int) []AutoCompleteSuggestion
	
	// FormatTableContent 格式化表格內容
	// 參數：tableContent（表格內容）
	// 回傳：格式化後的表格字串和可能的錯誤
	FormatTableContent(tableContent string) (string, error)
	
	// InsertLinkMarkdown 插入 Markdown 連結
	// 參數：text（連結文字）、url（連結網址）
	// 回傳：格式化的 Markdown 連結字串
	InsertLinkMarkdown(text, url string) string
	
	// InsertImageMarkdown 插入 Markdown 圖片
	// 參數：altText（替代文字）、imagePath（圖片路徑）
	// 回傳：格式化的 Markdown 圖片字串
	InsertImageMarkdown(altText, imagePath string) string
	
	// GetSupportedCodeLanguages 取得支援的程式語言列表
	// 回傳：支援的程式語言陣列
	GetSupportedCodeLanguages() []string
	
	// FormatCodeBlockMarkdown 格式化程式碼區塊
	// 參數：code（程式碼內容）、language（程式語言）
	// 回傳：格式化的 Markdown 程式碼區塊
	FormatCodeBlockMarkdown(code, language string) string
	
	// FormatMathExpressionMarkdown 格式化數學公式
	// 參數：expression（數學表達式）、isInline（是否為行內公式）
	// 回傳：格式化的 LaTeX 數學公式字串
	FormatMathExpressionMarkdown(expression string, isInline bool) string
	
	// ValidateMarkdownContent 驗證 Markdown 內容的語法正確性
	// 參數：content（要驗證的 Markdown 內容）
	// 回傳：驗證結果和可能的錯誤列表
	ValidateMarkdownContent(content string) (bool, []string)
	
	// GenerateTableTemplateMarkdown 生成表格模板
	// 參數：rows（行數）、cols（列數）
	// 回傳：表格模板字串
	GenerateTableTemplateMarkdown(rows, cols int) string
	
	// PreviewMarkdownWithHighlight 預覽 Markdown 內容並包含程式碼高亮
	// 參數：content（Markdown 格式的內容）
	// 回傳：轉換後的 HTML 字串（包含語法高亮）
	PreviewMarkdownWithHighlight(content string) string
	
	// GetSmartEditingService 取得智慧編輯服務實例
	// 回傳：SmartEditingService 介面實例
	GetSmartEditingService() SmartEditingService
	
	// SetSmartEditingService 設定智慧編輯服務實例
	// 參數：smartEditSvc（智慧編輯服務實例）
	SetSmartEditingService(smartEditSvc SmartEditingService)
}

// FileManagerService 定義檔案系統操作的介面
// 負責處理檔案和目錄的管理操作，包含列表、建立、刪除、重新命名和移動
type FileManagerService interface {
	// ListFiles 列出指定目錄中的檔案和子目錄
	// 參數：directory（目錄路徑）
	// 回傳：檔案資訊陣列和可能的錯誤
	ListFiles(directory string) ([]*models.FileInfo, error)
	
	// CreateDirectory 建立新目錄
	// 參數：path（目錄路徑）
	// 回傳：可能的錯誤
	CreateDirectory(path string) error
	
	// DeleteFile 刪除檔案或目錄
	// 參數：path（檔案或目錄路徑）
	// 回傳：可能的錯誤
	DeleteFile(path string) error
	
	// RenameFile 重新命名檔案或目錄
	// 參數：oldPath（舊路徑）、newPath（新路徑）
	// 回傳：可能的錯誤
	RenameFile(oldPath, newPath string) error
	
	// MoveFile 移動檔案或目錄到新位置
	// 參數：sourcePath（來源路徑）、destPath（目標路徑）
	// 回傳：可能的錯誤
	MoveFile(sourcePath, destPath string) error
	
	// CopyFile 複製檔案或目錄
	// 參數：sourcePath（來源路徑）、destPath（目標路徑）
	// 回傳：可能的錯誤
	CopyFile(sourcePath, destPath string) error
	
	// SearchFiles 搜尋檔案
	// 參數：searchPath（搜尋路徑）、pattern（搜尋模式）、includeSubdirs（是否包含子目錄）
	// 回傳：符合條件的檔案資訊陣列和可能的錯誤
	SearchFiles(searchPath, pattern string, includeSubdirs bool) ([]*models.FileInfo, error)
}

// EncryptionService 定義加密操作的介面
// 負責處理內容加密、解密、生物識別驗證和密碼驗證等安全功能
type EncryptionService interface {
	// EncryptContent 使用指定演算法和密碼加密內容
	// 參數：content（要加密的內容）、password（密碼）、algorithm（加密演算法）
	// 回傳：加密後的位元組陣列和可能的錯誤
	EncryptContent(content, password string, algorithm string) ([]byte, error)
	
	// DecryptContent 使用指定演算法和密碼解密內容
	// 參數：encryptedData（加密的資料）、password（密碼）、algorithm（加密演算法）
	// 回傳：解密後的內容字串和可能的錯誤
	DecryptContent(encryptedData []byte, password string, algorithm string) (string, error)
	
	// SetupBiometricAuth 為指定筆記設定生物識別驗證
	// 參數：noteID（筆記 ID）
	// 回傳：可能的錯誤
	SetupBiometricAuth(noteID string) error
	
	// AuthenticateWithBiometric 使用生物識別進行驗證
	// 參數：noteID（筆記 ID）
	// 回傳：驗證結果（成功/失敗）和可能的錯誤
	AuthenticateWithBiometric(noteID string) (bool, error)
	
	// ValidatePassword 驗證密碼強度是否符合要求
	// 參數：password（要驗證的密碼）
	// 回傳：密碼是否有效
	ValidatePassword(password string) bool
}

// AutoSaveService 定義自動保存操作的介面
// 負責管理筆記的自動保存功能，包含啟動、停止、立即保存和狀態查詢
type AutoSaveService interface {
	// StartAutoSave 為指定筆記啟動自動保存
	// 參數：note（要自動保存的筆記）、interval（保存間隔）
	StartAutoSave(note *models.Note, interval time.Duration)
	
	// StopAutoSave 停止指定筆記的自動保存
	// 參數：noteID（筆記 ID）
	StopAutoSave(noteID string)
	
	// SaveNow 立即保存指定筆記
	// 參數：noteID（筆記 ID）
	// 回傳：可能的錯誤
	SaveNow(noteID string) error
	
	// GetSaveStatus 取得指定筆記的保存狀態
	// 參數：noteID（筆記 ID）
	// 回傳：保存狀態資訊
	GetSaveStatus(noteID string) SaveStatus
}

// SaveStatus 代表保存操作的狀態資訊
// 包含保存進度、時間戳、錯誤資訊和統計資料
type SaveStatus struct {
	NoteID      string    `json:"note_id"`              // 筆記 ID
	IsSaving    bool      `json:"is_saving"`            // 是否正在保存中
	LastSaved   time.Time `json:"last_saved"`           // 最後保存時間
	LastError   error     `json:"last_error,omitempty"` // 最後發生的錯誤（如果有）
	SaveCount   int       `json:"save_count"`           // 累計保存次數
}

// NotificationType 定義通知類型的列舉
type NotificationType int

const (
	// NotificationInfo 資訊通知（藍色）
	NotificationInfo NotificationType = iota
	// NotificationSuccess 成功通知（綠色）
	NotificationSuccess
	// NotificationWarning 警告通知（橙色）
	NotificationWarning
	// NotificationError 錯誤通知（紅色）
	NotificationError
)

// Notification 代表一個通知訊息
// 包含通知的所有必要資訊和顯示屬性
type Notification struct {
	ID          string           `json:"id"`          // 通知的唯一識別碼
	Type        NotificationType `json:"type"`        // 通知類型
	Title       string           `json:"title"`       // 通知標題
	Message     string           `json:"message"`     // 通知內容
	Duration    time.Duration    `json:"duration"`    // 顯示持續時間
	CreatedAt   time.Time        `json:"created_at"`  // 建立時間
	IsRead      bool             `json:"is_read"`     // 是否已讀
	IsPersistent bool            `json:"is_persistent"` // 是否持久顯示（不自動消失）
}

// SaveStatusInfo 代表保存操作的狀態資訊
// 用於顯示檔案保存的即時狀態
type SaveStatusInfo struct {
	NoteID       string    `json:"note_id"`       // 筆記 ID
	FileName     string    `json:"file_name"`     // 檔案名稱
	IsSaving     bool      `json:"is_saving"`     // 是否正在保存中
	LastSaved    time.Time `json:"last_saved"`    // 最後保存時間
	SaveProgress float64   `json:"save_progress"` // 保存進度（0.0 - 1.0）
	HasChanges   bool      `json:"has_changes"`   // 是否有未保存的變更
}

// SettingsService 定義設定管理的介面
// 負責處理應用程式設定的載入、保存和預設值管理
type SettingsService interface {
	// LoadSettings 從儲存位置載入應用程式設定
	// 回傳：設定實例和可能的錯誤
	LoadSettings() (*models.Settings, error)
	
	// SaveSettings 將設定保存到儲存位置
	// 參數：settings（要保存的設定）
	// 回傳：可能的錯誤
	SaveSettings(settings *models.Settings) error
	
	// GetDefaultSettings 取得預設的應用程式設定
	// 回傳：預設設定實例
	GetDefaultSettings() *models.Settings
}

// ErrorService 定義錯誤處理的介面
// 負責統一的錯誤處理、本地化和日誌記錄功能
type ErrorService interface {
	// LogError 記錄錯誤到日誌檔案
	// 參數：err（要記錄的錯誤）、context（錯誤發生的上下文資訊）
	// 回傳：記錄過程中可能發生的錯誤
	LogError(err error, context string) error
	
	// LocalizeError 將錯誤訊息本地化為繁體中文
	// 參數：err（要本地化的錯誤）
	// 回傳：本地化後的錯誤訊息
	LocalizeError(err error) string
	
	// WrapError 包裝錯誤並添加上下文資訊
	// 參數：err（原始錯誤）、context（上下文資訊）
	// 回傳：包裝後的錯誤
	WrapError(err error, context string) error
	
	// HandleError 統一處理錯誤（記錄日誌並本地化）
	// 參數：err（要處理的錯誤）、context（錯誤發生的上下文）
	// 回傳：本地化後的錯誤訊息
	HandleError(err error, context string) string
	
	// CreateAppError 建立應用程式特定錯誤
	// 參數：code（錯誤代碼）、message（錯誤訊息）、details（詳細資訊）
	// 回傳：AppError 實例
	CreateAppError(code, message, details string) *models.AppError
	
	// IsRetryableError 判斷錯誤是否可重試
	// 參數：err（要檢查的錯誤）
	// 回傳：是否可重試
	IsRetryableError(err error) bool
}

// NotificationService 定義通知系統的介面
// 負責用戶通知顯示、保存狀態指示和操作回饋功能
type NotificationService interface {
	// ShowNotification 顯示通知訊息
	// 參數：notificationType（通知類型）、title（標題）、message（內容）、duration（持續時間）
	// 回傳：通知 ID
	ShowNotification(notificationType NotificationType, title, message string, duration time.Duration) string
	
	// ShowSuccess 顯示成功通知
	// 參數：title（標題）、message（內容）
	// 回傳：通知 ID
	ShowSuccess(title, message string) string
	
	// ShowError 顯示錯誤通知
	// 參數：title（標題）、message（內容）
	// 回傳：通知 ID
	ShowError(title, message string) string
	
	// ShowWarning 顯示警告通知
	// 參數：title（標題）、message（內容）
	// 回傳：通知 ID
	ShowWarning(title, message string) string
	
	// ShowInfo 顯示資訊通知
	// 參數：title（標題）、message（內容）
	// 回傳：通知 ID
	ShowInfo(title, message string) string
	
	// DismissNotification 關閉指定的通知
	// 參數：notificationID（通知 ID）
	// 回傳：是否成功關閉
	DismissNotification(notificationID string) bool
	
	// DismissAllNotifications 關閉所有通知
	DismissAllNotifications()
	
	// GetActiveNotifications 取得所有活躍的通知
	// 回傳：活躍通知的陣列
	GetActiveNotifications() []*Notification
	
	// UpdateSaveStatus 更新保存狀態指示器
	// 參數：noteID（筆記 ID）、fileName（檔案名稱）、status（保存狀態資訊）
	UpdateSaveStatus(noteID, fileName string, status SaveStatusInfo)
	
	// GetSaveStatus 取得指定筆記的保存狀態
	// 參數：noteID（筆記 ID）
	// 回傳：保存狀態資訊
	GetSaveStatus(noteID string) *SaveStatusInfo
	
	// ClearSaveStatus 清除指定筆記的保存狀態
	// 參數：noteID（筆記 ID）
	ClearSaveStatus(noteID string)
	
	// SetNotificationCallback 設定通知回調函數（用於 UI 更新）
	// 參數：callback（通知更新時的回調函數）
	SetNotificationCallback(callback func(*Notification))
	
	// SetSaveStatusCallback 設定保存狀態回調函數（用於 UI 更新）
	// 參數：callback（保存狀態更新時的回調函數）
	SetSaveStatusCallback(callback func(string, *SaveStatusInfo))
}

// SmartEditingService 定義智慧編輯功能的介面
// 負責處理 Markdown 語法自動完成、表格編輯、連結插入等進階編輯功能
type SmartEditingService interface {
	// AutoCompleteMarkdown 提供 Markdown 語法自動完成建議
	// 參數：content（當前內容）、cursorPosition（游標位置）
	// 回傳：自動完成建議陣列
	AutoCompleteMarkdown(content string, cursorPosition int) []AutoCompleteSuggestion
	
	// FormatTable 格式化表格內容
	// 參數：tableContent（表格內容）
	// 回傳：格式化後的表格字串和可能的錯誤
	FormatTable(tableContent string) (string, error)
	
	// InsertLink 插入連結
	// 參數：text（連結文字）、url（連結網址）
	// 回傳：格式化的 Markdown 連結字串
	InsertLink(text, url string) string
	
	// InsertImage 插入圖片
	// 參數：altText（替代文字）、imagePath（圖片路徑）
	// 回傳：格式化的 Markdown 圖片字串
	InsertImage(altText, imagePath string) string
	
	// HighlightCodeBlock 為程式碼區塊添加語法高亮
	// 參數：code（程式碼內容）、language（程式語言）
	// 回傳：帶有語法高亮的 HTML 字串
	HighlightCodeBlock(code, language string) string
	
	// FormatMathExpression 格式化數學公式
	// 參數：expression（數學表達式）、isInline（是否為行內公式）
	// 回傳：格式化的 LaTeX 數學公式字串
	FormatMathExpression(expression string, isInline bool) string
	
	// GetSupportedLanguages 取得支援的程式語言列表
	// 回傳：支援的程式語言陣列
	GetSupportedLanguages() []string
	
	// ValidateMarkdownSyntax 驗證 Markdown 語法的正確性
	// 參數：content（要驗證的 Markdown 內容）
	// 回傳：驗證結果和可能的錯誤列表
	ValidateMarkdownSyntax(content string) (bool, []string)
	
	// GenerateTableTemplate 生成表格模板
	// 參數：rows（行數）、cols（列數）
	// 回傳：表格模板字串
	GenerateTableTemplate(rows, cols int) string
	
	// FormatCodeBlock 格式化程式碼區塊
	// 參數：code（程式碼內容）、language（程式語言）
	// 回傳：格式化的 Markdown 程式碼區塊
	FormatCodeBlock(code, language string) string
}

// ExportService 定義匯出和分享功能的介面
// 負責處理筆記的各種格式匯出、批量匯出和分享功能
type ExportService interface {
	// ExportToPDF 將筆記匯出為 PDF 格式
	// 參數：note（要匯出的筆記）、outputPath（輸出檔案路徑）、options（匯出選項）
	// 回傳：可能的錯誤
	ExportToPDF(note *models.Note, outputPath string, options *ExportOptions) error
	
	// ExportToHTML 將筆記匯出為 HTML 格式
	// 參數：note（要匯出的筆記）、outputPath（輸出檔案路徑）、options（匯出選項）
	// 回傳：可能的錯誤
	ExportToHTML(note *models.Note, outputPath string, options *ExportOptions) error
	
	// ExportToWord 將筆記匯出為 Word 文件格式
	// 參數：note（要匯出的筆記）、outputPath（輸出檔案路徑）、options（匯出選項）
	// 回傳：可能的錯誤
	ExportToWord(note *models.Note, outputPath string, options *ExportOptions) error
	
	// BatchExport 批量匯出多個筆記
	// 參數：notes（要匯出的筆記陣列）、outputDir（輸出目錄）、format（匯出格式）、options（匯出選項）
	// 回傳：匯出結果和可能的錯誤
	BatchExport(notes []*models.Note, outputDir string, format ExportFormat, options *ExportOptions) (*BatchExportResult, error)
	
	// ShareNote 分享筆記
	// 參數：note（要分享的筆記）、shareOptions（分享選項）
	// 回傳：分享結果和可能的錯誤
	ShareNote(note *models.Note, shareOptions *ShareOptions) (*ShareResult, error)
	
	// GetSupportedFormats 取得支援的匯出格式列表
	// 回傳：支援的匯出格式陣列
	GetSupportedFormats() []ExportFormat
	
	// ValidateExportPath 驗證匯出路徑的有效性
	// 參數：path（要驗證的路徑）、format（匯出格式）
	// 回傳：路徑是否有效和可能的錯誤訊息
	ValidateExportPath(path string, format ExportFormat) (bool, string)
	
	// GetExportProgress 取得匯出進度
	// 參數：exportID（匯出任務 ID）
	// 回傳：匯出進度資訊
	GetExportProgress(exportID string) *ExportProgress
	
	// CancelExport 取消匯出任務
	// 參數：exportID（匯出任務 ID）
	// 回傳：是否成功取消
	CancelExport(exportID string) bool
}

// ExportFormat 定義匯出格式的列舉
type ExportFormat int

const (
	// ExportFormatPDF PDF 格式
	ExportFormatPDF ExportFormat = iota
	// ExportFormatHTML HTML 格式
	ExportFormatHTML
	// ExportFormatWord Word 文件格式
	ExportFormatWord
	// ExportFormatMarkdown Markdown 格式
	ExportFormatMarkdown
)

// String 回傳匯出格式的字串表示
func (f ExportFormat) String() string {
	switch f {
	case ExportFormatPDF:
		return "PDF"
	case ExportFormatHTML:
		return "HTML"
	case ExportFormatWord:
		return "Word"
	case ExportFormatMarkdown:
		return "Markdown"
	default:
		return "Unknown"
	}
}

// ExportOptions 定義匯出選項
type ExportOptions struct {
	IncludeMetadata    bool   `json:"include_metadata"`    // 是否包含元資料
	IncludeTableOfContents bool `json:"include_toc"`      // 是否包含目錄
	Theme              string `json:"theme"`               // 主題樣式
	FontSize           int    `json:"font_size"`           // 字體大小
	PageSize           string `json:"page_size"`           // 頁面大小（A4, Letter 等）
	Margins            string `json:"margins"`             // 頁面邊距
	IncludeImages      bool   `json:"include_images"`      // 是否包含圖片
	ImageQuality       int    `json:"image_quality"`       // 圖片品質（1-100）
	WatermarkText      string `json:"watermark_text"`      // 浮水印文字
	HeaderText         string `json:"header_text"`         // 頁首文字
	FooterText         string `json:"footer_text"`         // 頁尾文字
}

// BatchExportResult 代表批量匯出的結果
type BatchExportResult struct {
	TotalFiles    int      `json:"total_files"`    // 總檔案數
	SuccessCount  int      `json:"success_count"`  // 成功匯出數量
	FailureCount  int      `json:"failure_count"`  // 失敗匯出數量
	FailedFiles   []string `json:"failed_files"`   // 失敗的檔案列表
	OutputPath    string   `json:"output_path"`    // 輸出路徑
	ElapsedTime   time.Duration `json:"elapsed_time"` // 耗費時間
}

// ShareOptions 定義分享選項
type ShareOptions struct {
	ShareType     ShareType `json:"share_type"`     // 分享類型
	ExpiryTime    time.Time `json:"expiry_time"`    // 過期時間
	Password      string    `json:"password"`       // 分享密碼
	AllowDownload bool      `json:"allow_download"` // 是否允許下載
	AllowEdit     bool      `json:"allow_edit"`     // 是否允許編輯
	Recipients    []string  `json:"recipients"`     // 收件人列表
}

// ShareType 定義分享類型的列舉
type ShareType int

const (
	// ShareTypeLink 連結分享
	ShareTypeLink ShareType = iota
	// ShareTypeEmail 電子郵件分享
	ShareTypeEmail
	// ShareTypeAirDrop AirDrop 分享
	ShareTypeAirDrop
	// ShareTypeClipboard 複製到剪貼簿
	ShareTypeClipboard
)

// ShareResult 代表分享操作的結果
type ShareResult struct {
	ShareID   string    `json:"share_id"`   // 分享 ID
	ShareURL  string    `json:"share_url"`  // 分享連結
	ExpiryTime time.Time `json:"expiry_time"` // 過期時間
	Success   bool      `json:"success"`    // 是否成功
	Message   string    `json:"message"`    // 結果訊息
}

// ExportProgress 代表匯出進度資訊
type ExportProgress struct {
	ExportID    string        `json:"export_id"`    // 匯出任務 ID
	Progress    float64       `json:"progress"`     // 進度百分比（0.0 - 1.0）
	Status      ExportStatus  `json:"status"`       // 匯出狀態
	CurrentFile string        `json:"current_file"` // 當前處理的檔案
	ElapsedTime time.Duration `json:"elapsed_time"` // 已耗費時間
	EstimatedTime time.Duration `json:"estimated_time"` // 預估剩餘時間
	Error       error         `json:"error,omitempty"` // 錯誤資訊（如果有）
}

// ExportStatus 定義匯出狀態的列舉
type ExportStatus int

const (
	// ExportStatusPending 等待中
	ExportStatusPending ExportStatus = iota
	// ExportStatusInProgress 進行中
	ExportStatusInProgress
	// ExportStatusCompleted 已完成
	ExportStatusCompleted
	// ExportStatusFailed 失敗
	ExportStatusFailed
	// ExportStatusCancelled 已取消
	ExportStatusCancelled
)

// AutoCompleteSuggestion 代表自動完成建議項目
type AutoCompleteSuggestion struct {
	Text        string `json:"text"`        // 建議的文字內容
	Description string `json:"description"` // 建議的描述
	Type        string `json:"type"`        // 建議類型（header, list, link, etc.）
	InsertText  string `json:"insert_text"` // 要插入的實際文字
}