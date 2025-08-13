package models

import "errors"

// AppError 代表應用程式特定的錯誤類型
// 提供結構化的錯誤資訊，包含錯誤代碼、訊息和詳細資訊
type AppError struct {
	Code    string `json:"code"`              // 錯誤代碼，用於程式化處理
	Message string `json:"message"`           // 使用者友善的錯誤訊息
	Details string `json:"details,omitempty"` // 可選的詳細錯誤資訊
}

// Error 實作 error 介面，回傳錯誤訊息
// 這使得 AppError 可以作為標準的 Go error 使用
func (e *AppError) Error() string {
	return e.Message
}

// ValidationError 代表資料驗證錯誤
// 用於標識特定欄位的驗證失敗
type ValidationError struct {
	Field   string `json:"field"`   // 驗證失敗的欄位名稱
	Message string `json:"message"` // 驗證錯誤的詳細訊息
}

// Error 實作 error 介面，回傳驗證錯誤訊息
func (e *ValidationError) Error() string {
	return e.Message
}

// 錯誤代碼常數定義
// 這些常數用於標識不同類型的應用程式錯誤
const (
	ErrFileNotFound     = "FILE_NOT_FOUND"     // 檔案未找到
	ErrInvalidPassword  = "INVALID_PASSWORD"   // 密碼無效
	ErrEncryptionFailed = "ENCRYPTION_FAILED"  // 加密失敗
	ErrBiometricFailed  = "BIOMETRIC_FAILED"   // 生物識別驗證失敗
	ErrSaveFailed       = "SAVE_FAILED"        // 保存失敗
	ErrPermissionDenied = "PERMISSION_DENIED"  // 權限被拒絕
	ErrValidationFailed = "VALIDATION_FAILED"  // 資料驗證失敗
)

// 預定義的錯誤實例
// 這些錯誤用於設定驗證和其他常見的驗證場景
var (
	// 自動保存間隔驗證錯誤
	ErrInvalidAutoSaveInterval = errors.New("自動保存間隔必須在 1 到 60 分鐘之間")
	
	// 加密演算法驗證錯誤
	ErrInvalidEncryptionAlgorithm = errors.New("加密演算法必須是 'aes256' 或 'chacha20'")
	
	// 主題設定驗證錯誤
	ErrInvalidTheme = errors.New("主題必須是 'light'、'dark' 或 'auto'")
)

// NewAppError 建立一個新的應用程式錯誤實例
// 參數：
//   - code: 錯誤代碼（用於程式化處理）
//   - message: 使用者友善的錯誤訊息
//   - details: 可選的詳細錯誤資訊
// 回傳：指向新建立的 AppError 的指標
//
// 使用範例：
//   err := NewAppError(ErrFileNotFound, "找不到指定的筆記檔案", "檔案路徑：/path/to/note.md")
func NewAppError(code, message, details string) *AppError {
	return &AppError{
		Code:    code,    // 設定錯誤代碼
		Message: message, // 設定錯誤訊息
		Details: details, // 設定詳細資訊
	}
}

// NewValidationError 建立一個新的驗證錯誤實例
// 參數：
//   - field: 驗證失敗的欄位名稱
//   - message: 驗證錯誤的詳細訊息
// 回傳：指向新建立的 ValidationError 的指標
//
// 使用範例：
//   err := NewValidationError("Title", "筆記標題不能為空")
func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{
		Field:   field,   // 設定欄位名稱
		Message: message, // 設定錯誤訊息
	}
}