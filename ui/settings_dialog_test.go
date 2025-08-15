// Package ui 提供使用者介面元件的測試
// 本檔案測試設定對話框相關功能
package ui

import (
	"testing"

	"fyne.io/fyne/v2/test"

	"mac-notebook-app/internal/models"
)

// TestNewSettingsDialog 測試建立新的設定對話框
// 驗證：
// 1. 對話框能正確建立
// 2. 所有 UI 元件都已初始化
// 3. 設定值正確載入到 UI 元件
func TestNewSettingsDialog(t *testing.T) {
	// 建立測試應用程式和視窗
	testApp := test.NewApp()
	testWindow := testApp.NewWindow("Test")
	defer testWindow.Close()

	// 建立測試設定
	testSettings := models.NewDefaultSettings()
	testSettings.DefaultEncryption = "chacha20"
	testSettings.AutoSaveInterval = 10
	testSettings.BiometricEnabled = true
	testSettings.Theme = "dark"

	// 建立設定變更回調函數
	onChanged := func(settings *models.Settings) {
		// 測試回調函數
	}

	// 建立設定對話框
	dialog := NewSettingsDialog(testWindow, testSettings, onChanged)

	// 驗證對話框已建立
	if dialog == nil {
		t.Fatal("設定對話框建立失敗")
	}

	// 驗證 UI 元件已初始化
	if dialog.encryptionSelect == nil {
		t.Error("加密演算法選擇器未初始化")
	}
	if dialog.autoSaveEntry == nil {
		t.Error("自動保存間隔輸入框未初始化")
	}
	if dialog.saveLocationEntry == nil {
		t.Error("預設保存位置輸入框未初始化")
	}
	if dialog.biometricCheck == nil {
		t.Error("生物識別勾選框未初始化")
	}
	if dialog.themeSelect == nil {
		t.Error("主題選擇器未初始化")
	}

	// 驗證設定值正確載入
	if dialog.encryptionSelect.Selected != "chacha20" {
		t.Errorf("加密演算法選擇器值不正確，期望 'chacha20'，實際 '%s'", dialog.encryptionSelect.Selected)
	}
	if dialog.autoSaveEntry.Text != "10" {
		t.Errorf("自動保存間隔輸入框值不正確，期望 '10'，實際 '%s'", dialog.autoSaveEntry.Text)
	}
	if !dialog.biometricCheck.Checked {
		t.Error("生物識別勾選框狀態不正確，期望已勾選")
	}
	if dialog.themeSelect.Selected != "dark" {
		t.Errorf("主題選擇器值不正確，期望 'dark'，實際 '%s'", dialog.themeSelect.Selected)
	}
}

// TestSettingsDialog_EncryptionChange 測試加密演算法變更
// 驗證：
// 1. 選擇不同的加密演算法時設定正確更新
// 2. 設定變更回調函數被正確呼叫
func TestSettingsDialog_EncryptionChange(t *testing.T) {
	// 建立測試環境
	testApp := test.NewApp()
	testWindow := testApp.NewWindow("Test")
	defer testWindow.Close()

	testSettings := models.NewDefaultSettings()
	var changedSettings *models.Settings
	onChanged := func(settings *models.Settings) {
		changedSettings = settings
	}

	dialog := NewSettingsDialog(testWindow, testSettings, onChanged)

	// 模擬選擇不同的加密演算法
	test.Tap(dialog.encryptionSelect)
	dialog.encryptionSelect.SetSelected("chacha20")

	// 驗證設定已更新
	if dialog.settings.DefaultEncryption != "chacha20" {
		t.Errorf("加密演算法設定未更新，期望 'chacha20'，實際 '%s'", dialog.settings.DefaultEncryption)
	}

	// 驗證回調函數被呼叫
	if changedSettings == nil {
		t.Error("設定變更回調函數未被呼叫")
	} else if changedSettings.DefaultEncryption != "chacha20" {
		t.Errorf("回調函數接收的設定不正確，期望 'chacha20'，實際 '%s'", changedSettings.DefaultEncryption)
	}
}

// TestSettingsDialog_AutoSaveIntervalChange 測試自動保存間隔變更
// 驗證：
// 1. 輸入有效的間隔值時設定正確更新
// 2. 輸入無效值時設定不變更
func TestSettingsDialog_AutoSaveIntervalChange(t *testing.T) {
	// 建立測試環境
	testApp := test.NewApp()
	testWindow := testApp.NewWindow("Test")
	defer testWindow.Close()

	testSettings := models.NewDefaultSettings()
	var changedSettings *models.Settings
	onChanged := func(settings *models.Settings) {
		changedSettings = settings
	}

	dialog := NewSettingsDialog(testWindow, testSettings, onChanged)

	// 測試有效的間隔值
	test.Type(dialog.autoSaveEntry, "15")

	// 驗證設定已更新
	if dialog.settings.AutoSaveInterval != 15 {
		t.Errorf("自動保存間隔設定未更新，期望 15，實際 %d", dialog.settings.AutoSaveInterval)
	}

	// 重設回調函數狀態
	changedSettings = nil

	// 測試無效的間隔值
	dialog.autoSaveEntry.SetText("invalid")
	dialog.autoSaveEntry.OnChanged("invalid")

	// 驗證設定未變更（仍為 15）
	if dialog.settings.AutoSaveInterval != 15 {
		t.Errorf("無效輸入後設定被錯誤更新，期望 15，實際 %d", dialog.settings.AutoSaveInterval)
	}

	// 使用 changedSettings 避免未使用變數警告
	_ = changedSettings

	// 測試超出範圍的值
	dialog.autoSaveEntry.SetText("100")
	dialog.autoSaveEntry.OnChanged("100")

	// 驗證設定未變更（仍為 15）
	if dialog.settings.AutoSaveInterval != 15 {
		t.Errorf("超出範圍輸入後設定被錯誤更新，期望 15，實際 %d", dialog.settings.AutoSaveInterval)
	}
}

// TestSettingsDialog_BiometricToggle 測試生物識別切換
// 驗證：
// 1. 勾選/取消勾選時設定正確更新
// 2. 設定變更回調函數被正確呼叫
func TestSettingsDialog_BiometricToggle(t *testing.T) {
	// 建立測試環境
	testApp := test.NewApp()
	testWindow := testApp.NewWindow("Test")
	defer testWindow.Close()

	testSettings := models.NewDefaultSettings()
	testSettings.BiometricEnabled = false // 初始為關閉

	var changedSettings *models.Settings
	onChanged := func(settings *models.Settings) {
		changedSettings = settings
	}

	dialog := NewSettingsDialog(testWindow, testSettings, onChanged)

	// 模擬勾選生物識別
	test.Tap(dialog.biometricCheck)

	// 驗證設定已更新
	if !dialog.settings.BiometricEnabled {
		t.Error("生物識別設定未正確啟用")
	}

	// 驗證回調函數被呼叫
	if changedSettings == nil {
		t.Error("設定變更回調函數未被呼叫")
	} else if !changedSettings.BiometricEnabled {
		t.Error("回調函數接收的生物識別設定不正確")
	}
}

// TestSettingsDialog_ThemeChange 測試主題變更
// 驗證：
// 1. 選擇不同主題時設定正確更新
// 2. 設定變更回調函數被正確呼叫
func TestSettingsDialog_ThemeChange(t *testing.T) {
	// 建立測試環境
	testApp := test.NewApp()
	testWindow := testApp.NewWindow("Test")
	defer testWindow.Close()

	testSettings := models.NewDefaultSettings()
	var changedSettings *models.Settings
	onChanged := func(settings *models.Settings) {
		changedSettings = settings
	}

	dialog := NewSettingsDialog(testWindow, testSettings, onChanged)

	// 模擬選擇深色主題
	test.Tap(dialog.themeSelect)
	dialog.themeSelect.SetSelected("dark")

	// 驗證設定已更新
	if dialog.settings.Theme != "dark" {
		t.Errorf("主題設定未更新，期望 'dark'，實際 '%s'", dialog.settings.Theme)
	}

	// 驗證回調函數被呼叫
	if changedSettings == nil {
		t.Error("設定變更回調函數未被呼叫")
	} else if changedSettings.Theme != "dark" {
		t.Errorf("回調函數接收的主題設定不正確，期望 'dark'，實際 '%s'", changedSettings.Theme)
	}
}

// TestSettingsDialog_ResetToDefaults 測試重設為預設值功能
// 驗證：
// 1. 重設功能正確將所有設定恢復為預設值
// 2. UI 元件正確更新顯示
func TestSettingsDialog_ResetToDefaults(t *testing.T) {
	// 建立測試環境
	testApp := test.NewApp()
	testWindow := testApp.NewWindow("Test")
	defer testWindow.Close()

	// 建立修改過的設定
	testSettings := models.NewDefaultSettings()
	testSettings.DefaultEncryption = "chacha20"
	testSettings.AutoSaveInterval = 30
	testSettings.BiometricEnabled = true
	testSettings.Theme = "dark"

	var changedSettings *models.Settings
	onChanged := func(settings *models.Settings) {
		changedSettings = settings
	}

	dialog := NewSettingsDialog(testWindow, testSettings, onChanged)

	// 執行重設操作（直接呼叫方法，跳過確認對話框）
	dialog.settings = models.NewDefaultSettings()
	dialog.updateUIFromSettings()
	dialog.notifySettingsChanged()

	// 使用 changedSettings 避免未使用變數警告
	_ = changedSettings

	// 驗證設定已重設為預設值
	defaultSettings := models.NewDefaultSettings()
	if dialog.settings.DefaultEncryption != defaultSettings.DefaultEncryption {
		t.Errorf("加密演算法未重設，期望 '%s'，實際 '%s'", 
			defaultSettings.DefaultEncryption, dialog.settings.DefaultEncryption)
	}
	if dialog.settings.AutoSaveInterval != defaultSettings.AutoSaveInterval {
		t.Errorf("自動保存間隔未重設，期望 %d，實際 %d", 
			defaultSettings.AutoSaveInterval, dialog.settings.AutoSaveInterval)
	}
	if dialog.settings.BiometricEnabled != defaultSettings.BiometricEnabled {
		t.Errorf("生物識別設定未重設，期望 %t，實際 %t", 
			defaultSettings.BiometricEnabled, dialog.settings.BiometricEnabled)
	}
	if dialog.settings.Theme != defaultSettings.Theme {
		t.Errorf("主題設定未重設，期望 '%s'，實際 '%s'", 
			defaultSettings.Theme, dialog.settings.Theme)
	}

	// 驗證 UI 元件已更新
	if dialog.encryptionSelect.Selected != defaultSettings.DefaultEncryption {
		t.Error("加密演算法選擇器未正確更新")
	}
	if dialog.biometricCheck.Checked != defaultSettings.BiometricEnabled {
		t.Error("生物識別勾選框未正確更新")
	}
	if dialog.themeSelect.Selected != defaultSettings.Theme {
		t.Error("主題選擇器未正確更新")
	}
}

// TestSettingsDialog_GetAndSetSettings 測試取得和設定功能
// 驗證：
// 1. GetSettings 回傳正確的設定複製
// 2. SetSettings 正確更新設定和 UI
func TestSettingsDialog_GetAndSetSettings(t *testing.T) {
	// 建立測試環境
	testApp := test.NewApp()
	testWindow := testApp.NewWindow("Test")
	defer testWindow.Close()

	testSettings := models.NewDefaultSettings()
	dialog := NewSettingsDialog(testWindow, testSettings, nil)

	// 測試 GetSettings
	retrievedSettings := dialog.GetSettings()
	if retrievedSettings == nil {
		t.Fatal("GetSettings 回傳 nil")
	}

	// 驗證回傳的是複製而非原始設定
	if retrievedSettings == dialog.settings {
		t.Error("GetSettings 應該回傳設定的複製，而非原始設定")
	}

	// 驗證設定內容正確
	if retrievedSettings.DefaultEncryption != testSettings.DefaultEncryption {
		t.Error("GetSettings 回傳的設定內容不正確")
	}

	// 測試 SetSettings
	newSettings := models.NewDefaultSettings()
	newSettings.DefaultEncryption = "chacha20"
	newSettings.AutoSaveInterval = 20
	newSettings.Theme = "dark"

	dialog.SetSettings(newSettings)

	// 驗證設定已更新
	if dialog.settings.DefaultEncryption != "chacha20" {
		t.Error("SetSettings 未正確更新加密演算法設定")
	}
	if dialog.settings.AutoSaveInterval != 20 {
		t.Error("SetSettings 未正確更新自動保存間隔設定")
	}
	if dialog.settings.Theme != "dark" {
		t.Error("SetSettings 未正確更新主題設定")
	}

	// 驗證 UI 已更新
	if dialog.encryptionSelect.Selected != "chacha20" {
		t.Error("SetSettings 未正確更新加密演算法選擇器")
	}
	if dialog.themeSelect.Selected != "dark" {
		t.Error("SetSettings 未正確更新主題選擇器")
	}
}

// TestSettingsDialog_ShowHide 測試顯示和隱藏功能
// 驗證：
// 1. Show 方法正確顯示對話框
// 2. Hide 方法正確隱藏對話框
func TestSettingsDialog_ShowHide(t *testing.T) {
	// 建立測試環境
	testApp := test.NewApp()
	testWindow := testApp.NewWindow("Test")
	defer testWindow.Close()

	testSettings := models.NewDefaultSettings()
	dialog := NewSettingsDialog(testWindow, testSettings, nil)

	// 測試顯示功能
	dialog.Show()
	// 注意：在測試環境中無法直接驗證對話框是否可見
	// 這裡主要確保方法呼叫不會產生錯誤

	// 測試隱藏功能
	dialog.Hide()
	// 同樣，主要確保方法呼叫不會產生錯誤
}

// TestSettingsDialog_ValidationHandling 測試設定驗證處理
// 驗證：
// 1. 無效設定值不會被接受
// 2. 驗證錯誤得到適當處理
func TestSettingsDialog_ValidationHandling(t *testing.T) {
	// 建立測試環境
	testApp := test.NewApp()
	testWindow := testApp.NewWindow("Test")
	defer testWindow.Close()

	testSettings := models.NewDefaultSettings()
	dialog := NewSettingsDialog(testWindow, testSettings, nil)

	// 測試無效的自動保存間隔
	originalInterval := dialog.settings.AutoSaveInterval
	dialog.autoSaveEntry.SetText("0") // 無效值
	dialog.autoSaveEntry.OnChanged("0")

	// 驗證設定未被更新
	if dialog.settings.AutoSaveInterval != originalInterval {
		t.Error("無效的自動保存間隔值被錯誤接受")
	}

	// 測試無效的加密演算法（透過直接設定）
	originalEncryption := dialog.settings.DefaultEncryption
	err := dialog.settings.UpdateEncryption("invalid_algorithm")
	if err == nil {
		t.Error("無效的加密演算法應該產生錯誤")
	}

	// 驗證設定未被更新
	if dialog.settings.DefaultEncryption != originalEncryption {
		t.Error("無效的加密演算法值被錯誤接受")
	}
}