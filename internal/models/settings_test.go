// Package models 的測試檔案
// 包含 Settings 資料模型的完整單元測試
package models

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// TestNewDefaultSettings 測試預設設定建立功能
// 驗證預設設定是否具有正確的預設值
func TestNewDefaultSettings(t *testing.T) {
	settings := NewDefaultSettings()

	// 驗證預設值
	if settings.DefaultEncryption != "aes256" {
		t.Errorf("期望預設加密演算法為 'aes256'，實際得到 '%s'", settings.DefaultEncryption)
	}
	if settings.AutoSaveInterval != 5 {
		t.Errorf("期望預設自動保存間隔為 5，實際得到 %d", settings.AutoSaveInterval)
	}
	if settings.DefaultSaveLocation != "~/Documents/NotebookApp/notes" {
		t.Errorf("期望預設保存位置為 '~/Documents/NotebookApp/notes'，實際得到 '%s'", settings.DefaultSaveLocation)
	}
	if settings.BiometricEnabled {
		t.Error("預設應該不啟用生物識別驗證")
	}
	if settings.Theme != "auto" {
		t.Errorf("期望預設主題為 'auto'，實際得到 '%s'", settings.Theme)
	}

	// 驗證預設設定是有效的
	if err := settings.Validate(); err != nil {
		t.Errorf("預設設定應該是有效的，但得到錯誤：%v", err)
	}
}

// TestSettings_Validate 測試設定驗證功能
// 驗證各種無效設定情況下的驗證錯誤
func TestSettings_Validate(t *testing.T) {
	// 測試有效設定
	validSettings := NewDefaultSettings()
	if err := validSettings.Validate(); err != nil {
		t.Errorf("有效設定不應該產生驗證錯誤：%v", err)
	}

	// 測試無效的自動保存間隔（小於 1）
	invalidSettings := NewDefaultSettings()
	invalidSettings.AutoSaveInterval = 0
	if err := invalidSettings.Validate(); err == nil {
		t.Error("自動保存間隔為 0 應該產生驗證錯誤")
	}

	// 測試無效的自動保存間隔（大於 60）
	invalidSettings = NewDefaultSettings()
	invalidSettings.AutoSaveInterval = 61
	if err := invalidSettings.Validate(); err == nil {
		t.Error("自動保存間隔為 61 應該產生驗證錯誤")
	}

	// 測試無效的加密演算法
	invalidSettings = NewDefaultSettings()
	invalidSettings.DefaultEncryption = "invalid_algorithm"
	if err := invalidSettings.Validate(); err == nil {
		t.Error("無效的加密演算法應該產生驗證錯誤")
	}

	// 測試無效的主題設定
	invalidSettings = NewDefaultSettings()
	invalidSettings.Theme = "invalid_theme"
	if err := invalidSettings.Validate(); err == nil {
		t.Error("無效的主題設定應該產生驗證錯誤")
	}

	// 測試邊界值（有效範圍的邊界）
	validSettings = NewDefaultSettings()
	validSettings.AutoSaveInterval = 1
	if err := validSettings.Validate(); err != nil {
		t.Errorf("自動保存間隔為 1 應該是有效的：%v", err)
	}

	validSettings.AutoSaveInterval = 60
	if err := validSettings.Validate(); err != nil {
		t.Errorf("自動保存間隔為 60 應該是有效的：%v", err)
	}
}

// TestSettings_UpdateEncryption 測試加密演算法更新功能
// 驗證有效和無效的加密演算法設定
func TestSettings_UpdateEncryption(t *testing.T) {
	settings := NewDefaultSettings()

	// 測試有效的加密演算法
	validAlgorithms := []string{"aes256", "chacha20"}
	for _, algorithm := range validAlgorithms {
		err := settings.UpdateEncryption(algorithm)
		if err != nil {
			t.Errorf("更新為有效加密演算法 '%s' 不應該產生錯誤：%v", algorithm, err)
		}
		if settings.DefaultEncryption != algorithm {
			t.Errorf("期望加密演算法為 '%s'，實際得到 '%s'", algorithm, settings.DefaultEncryption)
		}
	}

	// 測試無效的加密演算法
	err := settings.UpdateEncryption("invalid_algorithm")
	if err == nil {
		t.Error("設定無效加密演算法應該產生錯誤")
	}
}

// TestSettings_UpdateAutoSaveInterval 測試自動保存間隔更新功能
// 驗證有效和無效的間隔值設定
func TestSettings_UpdateAutoSaveInterval(t *testing.T) {
	settings := NewDefaultSettings()

	// 測試有效的間隔值
	validIntervals := []int{1, 5, 30, 60}
	for _, interval := range validIntervals {
		err := settings.UpdateAutoSaveInterval(interval)
		if err != nil {
			t.Errorf("更新為有效間隔 %d 不應該產生錯誤：%v", interval, err)
		}
		if settings.AutoSaveInterval != interval {
			t.Errorf("期望自動保存間隔為 %d，實際得到 %d", interval, settings.AutoSaveInterval)
		}
	}

	// 測試無效的間隔值
	invalidIntervals := []int{0, -1, 61, 100}
	for _, interval := range invalidIntervals {
		err := settings.UpdateAutoSaveInterval(interval)
		if err == nil {
			t.Errorf("設定無效間隔 %d 應該產生錯誤", interval)
		}
	}
}

// TestSettings_UpdateTheme 測試主題更新功能
// 驗證有效和無效的主題設定
func TestSettings_UpdateTheme(t *testing.T) {
	settings := NewDefaultSettings()

	// 測試有效的主題
	validThemes := []string{"light", "dark", "auto"}
	for _, theme := range validThemes {
		err := settings.UpdateTheme(theme)
		if err != nil {
			t.Errorf("更新為有效主題 '%s' 不應該產生錯誤：%v", theme, err)
		}
		if settings.Theme != theme {
			t.Errorf("期望主題為 '%s'，實際得到 '%s'", theme, settings.Theme)
		}
	}

	// 測試無效的主題
	err := settings.UpdateTheme("invalid_theme")
	if err == nil {
		t.Error("設定無效主題應該產生錯誤")
	}
}

// TestSettings_UpdateDefaultSaveLocation 測試預設保存位置更新功能
// 驗證保存位置的設定是否正確
func TestSettings_UpdateDefaultSaveLocation(t *testing.T) {
	settings := NewDefaultSettings()
	newLocation := "/new/path/to/notes"

	settings.UpdateDefaultSaveLocation(newLocation)

	if settings.DefaultSaveLocation != newLocation {
		t.Errorf("期望預設保存位置為 '%s'，實際得到 '%s'", newLocation, settings.DefaultSaveLocation)
	}
}

// TestSettings_ToggleBiometric 測試生物識別切換功能
// 驗證生物識別狀態的切換是否正確
func TestSettings_ToggleBiometric(t *testing.T) {
	settings := NewDefaultSettings()
	originalState := settings.BiometricEnabled

	// 第一次切換
	newState := settings.ToggleBiometric()
	if newState == originalState {
		t.Error("切換後的狀態應該與原始狀態不同")
	}
	if settings.BiometricEnabled != newState {
		t.Error("回傳的狀態應該與設定中的狀態一致")
	}

	// 第二次切換（應該回到原始狀態）
	finalState := settings.ToggleBiometric()
	if finalState != originalState {
		t.Error("兩次切換後應該回到原始狀態")
	}
}

// TestSettings_SetBiometric 測試生物識別設定功能
// 驗證生物識別狀態的直接設定是否正確
func TestSettings_SetBiometric(t *testing.T) {
	settings := NewDefaultSettings()

	// 設定為啟用
	settings.SetBiometric(true)
	if !settings.BiometricEnabled {
		t.Error("設定為啟用後，生物識別應該被啟用")
	}

	// 設定為停用
	settings.SetBiometric(false)
	if settings.BiometricEnabled {
		t.Error("設定為停用後，生物識別應該被停用")
	}
}

// TestSettings_Clone 測試設定複製功能
// 驗證複製的設定是否包含相同資料但不同記憶體位址
func TestSettings_Clone(t *testing.T) {
	// 建立原始設定
	original := NewDefaultSettings()
	original.UpdateEncryption("chacha20")
	original.UpdateAutoSaveInterval(10)
	original.UpdateTheme("dark")
	original.SetBiometric(true)

	// 複製設定
	cloned := original.Clone()

	// 驗證資料相同
	if cloned.DefaultEncryption != original.DefaultEncryption {
		t.Error("複製的設定加密演算法應該相同")
	}
	if cloned.AutoSaveInterval != original.AutoSaveInterval {
		t.Error("複製的設定自動保存間隔應該相同")
	}
	if cloned.DefaultSaveLocation != original.DefaultSaveLocation {
		t.Error("複製的設定預設保存位置應該相同")
	}
	if cloned.BiometricEnabled != original.BiometricEnabled {
		t.Error("複製的設定生物識別狀態應該相同")
	}
	if cloned.Theme != original.Theme {
		t.Error("複製的設定主題應該相同")
	}

	// 驗證記憶體位址不同
	if cloned == original {
		t.Error("複製的設定應該有不同的記憶體位址")
	}

	// 驗證修改複製品不會影響原始設定
	cloned.UpdateEncryption("aes256")
	if original.DefaultEncryption == cloned.DefaultEncryption {
		t.Error("修改複製品不應該影響原始設定")
	}
}

// TestSettings_IsDefault 測試預設狀態檢查功能
// 驗證設定是否與預設設定相同的判斷
func TestSettings_IsDefault(t *testing.T) {
	// 預設設定應該被識別為預設狀態
	defaultSettings := NewDefaultSettings()
	if !defaultSettings.IsDefault() {
		t.Error("預設設定應該被識別為預設狀態")
	}

	// 修改後的設定不應該被識別為預設狀態
	modifiedSettings := NewDefaultSettings()
	modifiedSettings.UpdateEncryption("chacha20")
	if modifiedSettings.IsDefault() {
		t.Error("修改後的設定不應該被識別為預設狀態")
	}

	// 測試各個欄位的修改
	testCases := []struct {
		name   string
		modify func(*Settings)
	}{
		{"修改加密演算法", func(s *Settings) { s.UpdateEncryption("chacha20") }},
		{"修改自動保存間隔", func(s *Settings) { s.UpdateAutoSaveInterval(10) }},
		{"修改預設保存位置", func(s *Settings) { s.UpdateDefaultSaveLocation("/new/path") }},
		{"啟用生物識別", func(s *Settings) { s.SetBiometric(true) }},
		{"修改主題", func(s *Settings) { s.UpdateTheme("dark") }},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			settings := NewDefaultSettings()
			tc.modify(settings)
			if settings.IsDefault() {
				t.Errorf("%s 後不應該被識別為預設狀態", tc.name)
			}
		})
	}
}

// TestSettings_GetSupportedEncryptionAlgorithms 測試支援的加密演算法清單
// 驗證回傳的演算法清單是否正確
func TestSettings_GetSupportedEncryptionAlgorithms(t *testing.T) {
	settings := NewDefaultSettings()
	algorithms := settings.GetSupportedEncryptionAlgorithms()

	expectedAlgorithms := []string{"aes256", "chacha20"}
	if len(algorithms) != len(expectedAlgorithms) {
		t.Errorf("期望 %d 個演算法，實際得到 %d 個", len(expectedAlgorithms), len(algorithms))
	}

	for i, expected := range expectedAlgorithms {
		if i >= len(algorithms) || algorithms[i] != expected {
			t.Errorf("期望演算法 %d 為 '%s'，實際得到 '%s'", i, expected, algorithms[i])
		}
	}
}

// TestSettings_GetSupportedThemes 測試支援的主題清單
// 驗證回傳的主題清單是否正確
func TestSettings_GetSupportedThemes(t *testing.T) {
	settings := NewDefaultSettings()
	themes := settings.GetSupportedThemes()

	expectedThemes := []string{"light", "dark", "auto"}
	if len(themes) != len(expectedThemes) {
		t.Errorf("期望 %d 個主題，實際得到 %d 個", len(expectedThemes), len(themes))
	}

	for i, expected := range expectedThemes {
		if i >= len(themes) || themes[i] != expected {
			t.Errorf("期望主題 %d 為 '%s'，實際得到 '%s'", i, expected, themes[i])
		}
	}
}

// TestSettings_SaveToFile_LoadFromFile 測試設定的保存和載入功能
// 驗證設定能否正確保存到檔案並從檔案載入
func TestSettings_SaveToFile_LoadFromFile(t *testing.T) {
	// 建立測試目錄
	tempDir := t.TempDir()
	testFilePath := filepath.Join(tempDir, "test_settings.json")

	// 建立測試設定
	originalSettings := NewDefaultSettings()
	originalSettings.UpdateEncryption("chacha20")
	originalSettings.UpdateAutoSaveInterval(15)
	originalSettings.UpdateTheme("dark")
	originalSettings.SetBiometric(true)
	originalSettings.UpdateDefaultSaveLocation("/custom/path")

	// 保存設定到檔案
	err := originalSettings.SaveToFile(testFilePath)
	if err != nil {
		t.Fatalf("保存設定到檔案失敗：%v", err)
	}

	// 驗證檔案是否存在
	if _, err := os.Stat(testFilePath); os.IsNotExist(err) {
		t.Fatal("設定檔案應該存在")
	}

	// 從檔案載入設定
	loadedSettings, err := LoadFromFile(testFilePath)
	if err != nil {
		t.Fatalf("從檔案載入設定失敗：%v", err)
	}

	// 驗證載入的設定與原始設定相同
	if loadedSettings.DefaultEncryption != originalSettings.DefaultEncryption {
		t.Error("載入的加密演算法與原始設定不同")
	}
	if loadedSettings.AutoSaveInterval != originalSettings.AutoSaveInterval {
		t.Error("載入的自動保存間隔與原始設定不同")
	}
	if loadedSettings.DefaultSaveLocation != originalSettings.DefaultSaveLocation {
		t.Error("載入的預設保存位置與原始設定不同")
	}
	if loadedSettings.BiometricEnabled != originalSettings.BiometricEnabled {
		t.Error("載入的生物識別狀態與原始設定不同")
	}
	if loadedSettings.Theme != originalSettings.Theme {
		t.Error("載入的主題與原始設定不同")
	}
}

// TestLoadFromFile_NonExistentFile 測試載入不存在的設定檔案
// 驗證當檔案不存在時是否回傳預設設定
func TestLoadFromFile_NonExistentFile(t *testing.T) {
	nonExistentPath := "/path/that/does/not/exist/settings.json"

	settings, err := LoadFromFile(nonExistentPath)
	if err != nil {
		t.Fatalf("載入不存在的檔案不應該產生錯誤：%v", err)
	}

	// 應該回傳預設設定
	if !settings.IsDefault() {
		t.Error("載入不存在的檔案應該回傳預設設定")
	}
}

// TestLoadFromFile_InvalidJSON 測試載入無效 JSON 格式的設定檔案
// 驗證無效 JSON 是否產生適當的錯誤
func TestLoadFromFile_InvalidJSON(t *testing.T) {
	// 建立包含無效 JSON 的測試檔案
	tempDir := t.TempDir()
	testFilePath := filepath.Join(tempDir, "invalid_settings.json")

	invalidJSON := `{"invalid": json content}`
	err := os.WriteFile(testFilePath, []byte(invalidJSON), 0644)
	if err != nil {
		t.Fatalf("建立測試檔案失敗：%v", err)
	}

	// 嘗試載入無效的設定檔案
	_, err = LoadFromFile(testFilePath)
	if err == nil {
		t.Error("載入無效 JSON 應該產生錯誤")
	}
}

// TestLoadFromFile_InvalidSettings 測試載入包含無效設定值的檔案
// 驗證無效設定值是否產生驗證錯誤
func TestLoadFromFile_InvalidSettings(t *testing.T) {
	// 建立包含無效設定的測試檔案
	tempDir := t.TempDir()
	testFilePath := filepath.Join(tempDir, "invalid_values_settings.json")

	invalidSettings := Settings{
		DefaultEncryption:   "invalid_algorithm",
		AutoSaveInterval:    100, // 超出有效範圍
		DefaultSaveLocation: "/some/path",
		BiometricEnabled:    false,
		Theme:              "invalid_theme",
	}

	data, _ := json.MarshalIndent(invalidSettings, "", "  ")
	err := os.WriteFile(testFilePath, data, 0644)
	if err != nil {
		t.Fatalf("建立測試檔案失敗：%v", err)
	}

	// 嘗試載入無效的設定檔案
	_, err = LoadFromFile(testFilePath)
	if err == nil {
		t.Error("載入包含無效設定值的檔案應該產生錯誤")
	}
}

// TestGetDefaultSettingsPath 測試預設設定檔案路徑取得功能
// 驗證回傳的路徑格式是否正確
func TestGetDefaultSettingsPath(t *testing.T) {
	path := GetDefaultSettingsPath()

	// 路徑不應該為空
	if path == "" {
		t.Error("預設設定檔案路徑不應該為空")
	}

	// 路徑應該以 settings.json 結尾
	if filepath.Base(path) != "settings.json" {
		t.Error("預設設定檔案路徑應該以 settings.json 結尾")
	}
}

// TestLoadDefault_SaveDefault 測試預設位置的載入和保存功能
// 驗證預設位置的設定操作是否正常
func TestLoadDefault_SaveDefault(t *testing.T) {
	// 注意：這個測試會在實際的預設位置操作，需要小心處理

	// 建立測試設定
	testSettings := NewDefaultSettings()
	testSettings.UpdateEncryption("chacha20")

	// 保存到預設位置（實際上會建立目錄和檔案）
	err := testSettings.SaveDefault()
	if err != nil {
		t.Logf("保存到預設位置失敗（可能是權限問題）：%v", err)
		return // 跳過這個測試，因為可能沒有寫入權限
	}

	// 從預設位置載入
	loadedSettings, err := LoadDefault()
	if err != nil {
		t.Fatalf("從預設位置載入設定失敗：%v", err)
	}

	// 驗證載入的設定
	if loadedSettings.DefaultEncryption != "chacha20" {
		t.Error("從預設位置載入的設定不正確")
	}

	// 清理：嘗試刪除測試建立的檔案
	defaultPath := GetDefaultSettingsPath()
	os.Remove(defaultPath)
	os.Remove(filepath.Dir(defaultPath)) // 嘗試刪除目錄（如果為空）
}

// BenchmarkNewDefaultSettings 效能測試：預設設定建立
func BenchmarkNewDefaultSettings(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewDefaultSettings()
	}
}

// BenchmarkSettings_Validate 效能測試：設定驗證
func BenchmarkSettings_Validate(b *testing.B) {
	settings := NewDefaultSettings()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		settings.Validate()
	}
}

// BenchmarkSettings_Clone 效能測試：設定複製
func BenchmarkSettings_Clone(b *testing.B) {
	settings := NewDefaultSettings()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		settings.Clone()
	}
}

// BenchmarkSettings_SaveToFile 效能測試：設定保存
func BenchmarkSettings_SaveToFile(b *testing.B) {
	settings := NewDefaultSettings()
	tempDir := b.TempDir()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		testFilePath := filepath.Join(tempDir, "bench_settings.json")
		settings.SaveToFile(testFilePath)
		os.Remove(testFilePath) // 清理檔案
	}
}

// BenchmarkLoadFromFile 效能測試：設定載入
func BenchmarkLoadFromFile(b *testing.B) {
	// 準備測試檔案
	settings := NewDefaultSettings()
	tempDir := b.TempDir()
	testFilePath := filepath.Join(tempDir, "bench_settings.json")
	settings.SaveToFile(testFilePath)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		LoadFromFile(testFilePath)
	}
}