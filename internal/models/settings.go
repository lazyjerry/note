// Package models 定義了應用程式的核心資料模型
// 包含設定管理相關的結構體和方法
package models

import (
	"encoding/json"
	"os"
	"path/filepath"
)

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

// UpdateEncryption 更新預設加密演算法設定
// 參數：
//   - algorithm: 新的加密演算法（"aes256" 或 "chacha20"）
// 回傳：如果演算法無效則回傳錯誤，否則回傳 nil
//
// 執行流程：
// 1. 驗證新的加密演算法是否有效
// 2. 如果有效，更新設定值
// 3. 回傳操作結果
func (s *Settings) UpdateEncryption(algorithm string) error {
	if algorithm != "aes256" && algorithm != "chacha20" {
		return ErrInvalidEncryptionAlgorithm
	}
	s.DefaultEncryption = algorithm
	return nil
}

// UpdateAutoSaveInterval 更新自動保存間隔設定
// 參數：
//   - interval: 新的自動保存間隔（分鐘，範圍 1-60）
// 回傳：如果間隔無效則回傳錯誤，否則回傳 nil
//
// 執行流程：
// 1. 驗證新的間隔值是否在有效範圍內
// 2. 如果有效，更新設定值
// 3. 回傳操作結果
func (s *Settings) UpdateAutoSaveInterval(interval int) error {
	if interval < 1 || interval > 60 {
		return ErrInvalidAutoSaveInterval
	}
	s.AutoSaveInterval = interval
	return nil
}

// UpdateTheme 更新主題設定
// 參數：
//   - theme: 新的主題設定（"light"、"dark" 或 "auto"）
// 回傳：如果主題無效則回傳錯誤，否則回傳 nil
//
// 執行流程：
// 1. 驗證新的主題設定是否有效
// 2. 如果有效，更新設定值
// 3. 回傳操作結果
func (s *Settings) UpdateTheme(theme string) error {
	if theme != "light" && theme != "dark" && theme != "auto" {
		return ErrInvalidTheme
	}
	s.Theme = theme
	return nil
}

// UpdateDefaultSaveLocation 更新預設保存位置
// 參數：
//   - location: 新的預設保存位置路徑
//
// 執行流程：
// 1. 直接更新預設保存位置
// 2. 不進行路徑驗證，因為路徑可能在設定時尚不存在
func (s *Settings) UpdateDefaultSaveLocation(location string) {
	s.DefaultSaveLocation = location
}

// ToggleBiometric 切換生物識別驗證的啟用狀態
// 回傳：切換後的生物識別啟用狀態
//
// 執行流程：
// 1. 將當前的生物識別狀態取反
// 2. 回傳新的狀態值
func (s *Settings) ToggleBiometric() bool {
	s.BiometricEnabled = !s.BiometricEnabled
	return s.BiometricEnabled
}

// SetBiometric 設定生物識別驗證的啟用狀態
// 參數：
//   - enabled: 是否啟用生物識別驗證
//
// 執行流程：
// 1. 直接設定生物識別的啟用狀態
func (s *Settings) SetBiometric(enabled bool) {
	s.BiometricEnabled = enabled
}

// Clone 建立設定的深度複製
// 回傳：新的設定實例，包含相同的資料但不同的記憶體位址
//
// 執行流程：
// 1. 建立新的 Settings 結構體
// 2. 複製所有欄位的值
// 3. 回傳新的設定實例
func (s *Settings) Clone() *Settings {
	return &Settings{
		DefaultEncryption:   s.DefaultEncryption,
		AutoSaveInterval:    s.AutoSaveInterval,
		DefaultSaveLocation: s.DefaultSaveLocation,
		BiometricEnabled:    s.BiometricEnabled,
		Theme:              s.Theme,
	}
}

// IsDefault 檢查當前設定是否與預設設定相同
// 回傳：如果與預設設定相同則回傳 true，否則回傳 false
//
// 執行流程：
// 1. 建立預設設定實例
// 2. 逐一比較各個欄位
// 3. 如果所有欄位都相同則回傳 true
func (s *Settings) IsDefault() bool {
	defaultSettings := NewDefaultSettings()
	return s.DefaultEncryption == defaultSettings.DefaultEncryption &&
		s.AutoSaveInterval == defaultSettings.AutoSaveInterval &&
		s.DefaultSaveLocation == defaultSettings.DefaultSaveLocation &&
		s.BiometricEnabled == defaultSettings.BiometricEnabled &&
		s.Theme == defaultSettings.Theme
}

// GetSupportedEncryptionAlgorithms 取得支援的加密演算法清單
// 回傳：包含所有支援的加密演算法名稱的字串切片
//
// 執行流程：
// 1. 回傳預定義的支援演算法清單
func (s *Settings) GetSupportedEncryptionAlgorithms() []string {
	return []string{"aes256", "chacha20"}
}

// GetSupportedThemes 取得支援的主題清單
// 回傳：包含所有支援的主題名稱的字串切片
//
// 執行流程：
// 1. 回傳預定義的支援主題清單
func (s *Settings) GetSupportedThemes() []string {
	return []string{"light", "dark", "auto"}
}

// LoadFromFile 從指定的檔案載入設定
// 參數：
//   - filePath: 設定檔案的完整路徑
// 回傳：載入的設定實例和可能的錯誤
//
// 執行流程：
// 1. 檢查檔案是否存在
// 2. 讀取檔案內容
// 3. 解析 JSON 格式的設定資料
// 4. 驗證載入的設定是否有效
// 5. 回傳設定實例或錯誤
func LoadFromFile(filePath string) (*Settings, error) {
	// 檢查檔案是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// 檔案不存在時回傳預設設定
		return NewDefaultSettings(), nil
	}

	// 讀取檔案內容
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, NewAppError(ErrFileNotFound, "無法讀取設定檔案", err.Error())
	}

	// 解析 JSON 資料
	var settings Settings
	if err := json.Unmarshal(data, &settings); err != nil {
		return nil, NewAppError(ErrValidationFailed, "設定檔案格式無效", err.Error())
	}

	// 驗證載入的設定
	if err := settings.Validate(); err != nil {
		return nil, err
	}

	return &settings, nil
}

// SaveToFile 將設定保存到指定的檔案
// 參數：
//   - filePath: 要保存的檔案完整路徑
// 回傳：保存操作的錯誤（如果有）
//
// 執行流程：
// 1. 驗證當前設定是否有效
// 2. 確保目標目錄存在
// 3. 將設定序列化為 JSON 格式
// 4. 寫入檔案
// 5. 回傳操作結果
func (s *Settings) SaveToFile(filePath string) error {
	// 驗證設定有效性
	if err := s.Validate(); err != nil {
		return err
	}

	// 確保目標目錄存在
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return NewAppError(ErrPermissionDenied, "無法建立設定目錄", err.Error())
	}

	// 序列化設定為 JSON
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return NewAppError(ErrSaveFailed, "無法序列化設定資料", err.Error())
	}

	// 寫入檔案
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return NewAppError(ErrSaveFailed, "無法寫入設定檔案", err.Error())
	}

	return nil
}

// GetDefaultSettingsPath 取得預設的設定檔案路徑
// 回傳：預設設定檔案的完整路徑
//
// 執行流程：
// 1. 取得使用者的主目錄
// 2. 組合預設的設定檔案路徑
// 3. 回傳完整路徑
func GetDefaultSettingsPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// 如果無法取得主目錄，使用當前目錄
		return ".notebook/settings.json"
	}
	return filepath.Join(homeDir, "Documents", "NotebookApp", ".notebook", "settings.json")
}

// LoadDefault 載入預設位置的設定檔案
// 回傳：載入的設定實例和可能的錯誤
//
// 執行流程：
// 1. 取得預設設定檔案路徑
// 2. 呼叫 LoadFromFile 載入設定
// 3. 回傳載入結果
func LoadDefault() (*Settings, error) {
	return LoadFromFile(GetDefaultSettingsPath())
}

// SaveDefault 將設定保存到預設位置
// 回傳：保存操作的錯誤（如果有）
//
// 執行流程：
// 1. 取得預設設定檔案路徑
// 2. 呼叫 SaveToFile 保存設定
// 3. 回傳保存結果
func (s *Settings) SaveDefault() error {
	return s.SaveToFile(GetDefaultSettingsPath())
}