// Package repositories 定義了資料存取層的介面
// 這些介面抽象化了資料儲存和檢索的操作，支援不同的儲存後端
package repositories

import "mac-notebook-app/internal/models" // 引入資料模型

// FileRepository 定義檔案操作的介面
// 負責處理檔案系統的基本操作，包含讀取、寫入、刪除和目錄管理
type FileRepository interface {
	// ReadFile 讀取指定路徑的檔案內容
	// 參數：path（檔案路徑）
	// 回傳：檔案內容的位元組陣列和可能的錯誤
	ReadFile(path string) ([]byte, error)
	
	// WriteFile 將資料寫入指定路徑的檔案
	// 參數：path（檔案路徑）、data（要寫入的資料）
	// 回傳：可能的錯誤
	WriteFile(path string, data []byte) error
	
	// FileExists 檢查指定路徑的檔案是否存在
	// 參數：path（檔案路徑）
	// 回傳：檔案是否存在
	FileExists(path string) bool
	
	// DeleteFile 刪除指定路徑的檔案
	// 參數：path（檔案路徑）
	// 回傳：可能的錯誤
	DeleteFile(path string) error
	
	// CreateDirectory 建立指定路徑的目錄
	// 參數：path（目錄路徑）
	// 回傳：可能的錯誤
	CreateDirectory(path string) error
	
	// ListDirectory 列出指定目錄中的檔案和子目錄
	// 參數：path（目錄路徑）
	// 回傳：檔案資訊陣列和可能的錯誤
	ListDirectory(path string) ([]*models.FileInfo, error)
}

// SettingsRepository 定義設定持久化的介面
// 負責處理應用程式設定的儲存和載入操作
type SettingsRepository interface {
	// LoadSettings 從儲存位置載入設定
	// 回傳：設定實例和可能的錯誤
	LoadSettings() (*models.Settings, error)
	
	// SaveSettings 將設定保存到儲存位置
	// 參數：settings（要保存的設定）
	// 回傳：可能的錯誤
	SaveSettings(settings *models.Settings) error
	
	// SettingsExist 檢查設定檔案是否存在
	// 回傳：設定檔案是否存在
	SettingsExist() bool
}

// EncryptionRepository 定義加密金鑰管理的介面
// 負責處理密碼雜湊和生物識別金鑰的安全儲存和檢索
type EncryptionRepository interface {
	// StorePasswordHash 儲存指定筆記的密碼雜湊
	// 參數：noteID（筆記 ID）、hash（密碼雜湊）
	// 回傳：可能的錯誤
	StorePasswordHash(noteID, hash string) error
	
	// GetPasswordHash 取得指定筆記的密碼雜湊
	// 參數：noteID（筆記 ID）
	// 回傳：密碼雜湊和可能的錯誤
	GetPasswordHash(noteID string) (string, error)
	
	// StoreBiometricKey 儲存指定筆記的生物識別金鑰
	// 參數：noteID（筆記 ID）、keyData（金鑰資料）
	// 回傳：可能的錯誤
	StoreBiometricKey(noteID string, keyData []byte) error
	
	// GetBiometricKey 取得指定筆記的生物識別金鑰
	// 參數：noteID（筆記 ID）
	// 回傳：金鑰資料和可能的錯誤
	GetBiometricKey(noteID string) ([]byte, error)
	
	// DeleteKeys 刪除指定筆記的所有金鑰資料
	// 參數：noteID（筆記 ID）
	// 回傳：可能的錯誤
	DeleteKeys(noteID string) error
}