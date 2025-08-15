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