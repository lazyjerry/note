package models

// Settings 代表應用程式的組態設定
// 包含加密、自動保存、主題等各種使用者偏好設定
type Settings struct {
	DefaultEncryption   string `json:"default_encryption"`   // 預設加密演算法："aes256" 或 "chacha20"
	AutoSaveInterval    int    `json:"auto_save_interval"`    // 自動保存間隔（分鐘）
	DefaultSaveLocation string `json:"default_save_location"` // 預設筆記保存位置
	BiometricEnabled    bool   `json:"biometric_enabled"`     // 是否啟用生物識別驗證
	Theme              string `json:"theme"`                 // 主題設定："light"（淺色）、"dark"（深色）、"auto"（自動）
}

// NewDefaultSettings 建立具有預設值的設定實例
// 回傳：指向新建立設定的指標
//
// 預設值說明：
// - 加密演算法：AES-256（較為通用且安全）
// - 自動保存間隔：5 分鐘（平衡效能和資料安全）
// - 預設保存位置：使用者文件夾下的 NotebookApp/notes 目錄
// - 生物識別：預設關閉（需要使用者手動啟用）
// - 主題：自動（跟隨系統設定）
func NewDefaultSettings() *Settings {
	return &Settings{
		DefaultEncryption:   "aes256",                        // 使用 AES-256 作為預設加密演算法
		AutoSaveInterval:    5,                               // 每 5 分鐘自動保存一次
		DefaultSaveLocation: "~/Documents/NotebookApp/notes", // 預設保存到文件夾
		BiometricEnabled:    false,                           // 預設不啟用生物識別
		Theme:              "auto",                           // 自動跟隨系統主題
	}
}

// Validate 驗證設定值是否有效
// 回傳：如果設定無效則回傳對應的錯誤，否則回傳 nil
//
// 驗證規則：
// 1. 自動保存間隔必須在 1-60 分鐘之間
// 2. 加密演算法必須是支援的類型（aes256 或 chacha20）
// 3. 主題設定必須是有效的選項（light、dark 或 auto）
//
// 執行流程：
// 1. 檢查自動保存間隔的有效範圍
// 2. 驗證加密演算法是否受支援
// 3. 確認主題設定是否有效
// 4. 如果所有驗證都通過，回傳 nil
func (s *Settings) Validate() error {
	// 驗證自動保存間隔（1-60 分鐘）
	if s.AutoSaveInterval < 1 || s.AutoSaveInterval > 60 {
		return ErrInvalidAutoSaveInterval
	}
	
	// 驗證加密演算法是否受支援
	if s.DefaultEncryption != "aes256" && s.DefaultEncryption != "chacha20" {
		return ErrInvalidEncryptionAlgorithm
	}
	
	// 驗證主題設定是否有效
	if s.Theme != "light" && s.Theme != "dark" && s.Theme != "auto" {
		return ErrInvalidTheme
	}
	
	// 所有驗證都通過
	return nil
}