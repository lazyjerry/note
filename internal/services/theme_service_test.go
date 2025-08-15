// Package services 提供應用程式的核心業務邏輯服務的測試
// 本檔案測試主題管理相關功能
package services

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"

	"mac-notebook-app/internal/models"
)

// mockThemeListener 模擬主題監聽器，用於測試
type mockThemeListener struct {
	lastTheme string
	callCount int
}

// OnThemeChanged 實作 ThemeListener 介面
func (mtl *mockThemeListener) OnThemeChanged(themeName string) {
	mtl.lastTheme = themeName
	mtl.callCount++
}

// TestNewThemeService 測試建立新的主題服務
// 驗證：
// 1. 主題服務能正確建立
// 2. 初始主題設定正確載入
// 3. 系統主題偵測功能正常
func TestNewThemeService(t *testing.T) {
	// 建立測試應用程式和設定
	testApp := test.NewApp()
	testSettings := models.NewDefaultSettings()
	testSettings.Theme = "dark"

	// 建立主題服務
	themeService := NewThemeService(testApp, testSettings)

	// 驗證主題服務已建立
	if themeService == nil {
		t.Fatal("主題服務建立失敗")
	}

	// 驗證當前主題設定
	if themeService.GetCurrentTheme() != "dark" {
		t.Errorf("當前主題設定不正確，期望 'dark'，實際 '%s'", themeService.GetCurrentTheme())
	}

	// 驗證系統主題偵測
	systemTheme := themeService.GetSystemTheme()
	if systemTheme != "light" && systemTheme != "dark" {
		t.Errorf("系統主題偵測結果無效，實際 '%s'", systemTheme)
	}
}

// TestThemeService_SetTheme 測試主題設定功能
// 驗證：
// 1. 有效主題能正確設定
// 2. 無效主題會回傳錯誤
// 3. 主題變更會通知監聽器
func TestThemeService_SetTheme(t *testing.T) {
	// 建立測試環境
	testApp := test.NewApp()
	testSettings := models.NewDefaultSettings()
	themeService := NewThemeService(testApp, testSettings)

	// 建立模擬監聽器
	mockListener := &mockThemeListener{}
	themeService.AddThemeListener(mockListener)

	// 測試設定有效主題
	err := themeService.SetTheme("dark")
	if err != nil {
		t.Errorf("設定有效主題失敗：%v", err)
	}

	// 驗證主題已更新
	if themeService.GetCurrentTheme() != "dark" {
		t.Errorf("主題未正確更新，期望 'dark'，實際 '%s'", themeService.GetCurrentTheme())
	}

	// 驗證監聽器被通知
	if mockListener.lastTheme != "dark" {
		t.Errorf("監聽器未收到正確的主題通知，期望 'dark'，實際 '%s'", mockListener.lastTheme)
	}
	if mockListener.callCount != 1 {
		t.Errorf("監聽器呼叫次數不正確，期望 1，實際 %d", mockListener.callCount)
	}

	// 測試設定無效主題
	err = themeService.SetTheme("invalid_theme")
	if err == nil {
		t.Error("設定無效主題應該回傳錯誤")
	}

	// 驗證主題未變更
	if themeService.GetCurrentTheme() != "dark" {
		t.Error("無效主題設定後，當前主題不應該變更")
	}
}

// TestThemeService_AutoTheme 測試自動主題功能
// 驗證：
// 1. 自動主題能正確設定
// 2. 自動主題會根據系統主題決定實際主題
func TestThemeService_AutoTheme(t *testing.T) {
	// 建立測試環境
	testApp := test.NewApp()
	testSettings := models.NewDefaultSettings()
	themeService := NewThemeService(testApp, testSettings)

	// 設定為自動主題
	err := themeService.SetTheme("auto")
	if err != nil {
		t.Errorf("設定自動主題失敗：%v", err)
	}

	// 驗證主題設定為自動
	if themeService.GetCurrentTheme() != "auto" {
		t.Errorf("自動主題設定不正確，期望 'auto'，實際 '%s'", themeService.GetCurrentTheme())
	}

	// 驗證系統主題偵測
	systemTheme := themeService.GetSystemTheme()
	if systemTheme != "light" && systemTheme != "dark" {
		t.Errorf("系統主題偵測結果無效，實際 '%s'", systemTheme)
	}
}

// TestThemeService_ThemeListeners 測試主題監聽器管理
// 驗證：
// 1. 監聽器能正確新增和移除
// 2. 多個監聽器都能收到通知
// 3. 移除的監聽器不會收到通知
func TestThemeService_ThemeListeners(t *testing.T) {
	// 建立測試環境
	testApp := test.NewApp()
	testSettings := models.NewDefaultSettings()
	themeService := NewThemeService(testApp, testSettings)

	// 建立多個模擬監聽器
	listener1 := &mockThemeListener{}
	listener2 := &mockThemeListener{}
	listener3 := &mockThemeListener{}

	// 新增監聽器
	themeService.AddThemeListener(listener1)
	themeService.AddThemeListener(listener2)
	themeService.AddThemeListener(listener3)

	// 變更主題
	themeService.SetTheme("dark")

	// 驗證所有監聽器都收到通知
	if listener1.callCount != 1 || listener1.lastTheme != "dark" {
		t.Error("監聽器1未正確收到主題變更通知")
	}
	if listener2.callCount != 1 || listener2.lastTheme != "dark" {
		t.Error("監聽器2未正確收到主題變更通知")
	}
	if listener3.callCount != 1 || listener3.lastTheme != "dark" {
		t.Error("監聽器3未正確收到主題變更通知")
	}

	// 移除一個監聽器
	themeService.RemoveThemeListener(listener2)

	// 重設監聽器狀態
	listener1.callCount = 0
	listener2.callCount = 0
	listener3.callCount = 0

	// 再次變更主題
	themeService.SetTheme("light")

	// 驗證只有未移除的監聽器收到通知
	if listener1.callCount != 1 || listener1.lastTheme != "light" {
		t.Error("監聽器1未正確收到第二次主題變更通知")
	}
	if listener2.callCount != 0 {
		t.Error("已移除的監聽器2不應該收到通知")
	}
	if listener3.callCount != 1 || listener3.lastTheme != "light" {
		t.Error("監聽器3未正確收到第二次主題變更通知")
	}
}

// TestThemeService_GetAvailableThemes 測試取得可用主題列表
// 驗證：
// 1. 回傳的主題列表包含所有預期的主題
// 2. 主題列表不為空
func TestThemeService_GetAvailableThemes(t *testing.T) {
	// 建立測試環境
	testApp := test.NewApp()
	testSettings := models.NewDefaultSettings()
	themeService := NewThemeService(testApp, testSettings)

	// 取得可用主題列表
	themes := themeService.GetAvailableThemes()

	// 驗證主題列表不為空
	if len(themes) == 0 {
		t.Error("可用主題列表不應該為空")
	}

	// 驗證包含預期的主題
	expectedThemes := []string{"light", "dark", "auto"}
	if len(themes) != len(expectedThemes) {
		t.Errorf("主題列表長度不正確，期望 %d，實際 %d", len(expectedThemes), len(themes))
	}

	for i, expected := range expectedThemes {
		if i >= len(themes) || themes[i] != expected {
			t.Errorf("主題列表內容不正確，位置 %d 期望 '%s'，實際 '%s'", i, expected, themes[i])
		}
	}
}

// TestThemeService_GetThemeDisplayName 測試取得主題顯示名稱
// 驗證：
// 1. 所有有效主題都有對應的中文顯示名稱
// 2. 無效主題回傳適當的錯誤訊息
func TestThemeService_GetThemeDisplayName(t *testing.T) {
	// 建立測試環境
	testApp := test.NewApp()
	testSettings := models.NewDefaultSettings()
	themeService := NewThemeService(testApp, testSettings)

	// 測試有效主題的顯示名稱
	testCases := map[string]string{
		"light": "淺色主題",
		"dark":  "深色主題",
		"auto":  "自動（跟隨系統）",
	}

	for theme, expectedName := range testCases {
		displayName := themeService.GetThemeDisplayName(theme)
		if displayName != expectedName {
			t.Errorf("主題 '%s' 的顯示名稱不正確，期望 '%s'，實際 '%s'", theme, expectedName, displayName)
		}
	}

	// 測試無效主題
	invalidDisplayName := themeService.GetThemeDisplayName("invalid")
	if invalidDisplayName != "未知主題" {
		t.Errorf("無效主題的顯示名稱不正確，期望 '未知主題'，實際 '%s'", invalidDisplayName)
	}
}

// TestThemeService_IsSystemThemeSupported 測試系統主題支援檢查
// 驗證：
// 1. 系統主題支援檢查回傳合理的結果
func TestThemeService_IsSystemThemeSupported(t *testing.T) {
	// 建立測試環境
	testApp := test.NewApp()
	testSettings := models.NewDefaultSettings()
	themeService := NewThemeService(testApp, testSettings)

	// 檢查系統主題支援
	supported := themeService.IsSystemThemeSupported()

	// 在測試環境中，我們期望支援系統主題偵測
	if !supported {
		t.Error("在 macOS 環境中應該支援系統主題偵測")
	}
}

// TestThemeService_RefreshSystemTheme 測試重新整理系統主題
// 驗證：
// 1. 重新整理功能不會產生錯誤
// 2. 自動模式下會重新套用主題
func TestThemeService_RefreshSystemTheme(t *testing.T) {
	// 建立測試環境
	testApp := test.NewApp()
	testSettings := models.NewDefaultSettings()
	themeService := NewThemeService(testApp, testSettings)

	// 建立模擬監聽器
	mockListener := &mockThemeListener{}
	themeService.AddThemeListener(mockListener)

	// 設定為自動主題
	themeService.SetTheme("auto")
	
	// 重設監聽器狀態
	mockListener.callCount = 0

	// 重新整理系統主題
	themeService.RefreshSystemTheme()

	// 驗證自動模式下監聽器被通知
	if mockListener.callCount != 1 {
		t.Errorf("自動模式下重新整理系統主題應該通知監聽器，期望呼叫 1 次，實際 %d 次", mockListener.callCount)
	}
	if mockListener.lastTheme != "auto" {
		t.Errorf("重新整理後監聽器收到的主題不正確，期望 'auto'，實際 '%s'", mockListener.lastTheme)
	}
}

// TestCustomTheme 測試自訂主題實作
// 驗證：
// 1. 自訂主題能正確實作 fyne.Theme 介面
// 2. 主題變體能正確套用
func TestCustomTheme(t *testing.T) {
	// 建立自訂主題實例
	lightTheme := &customTheme{variant: theme.VariantLight}
	darkTheme := &customTheme{variant: theme.VariantDark}

	// 測試顏色方法
	lightColor := lightTheme.Color(theme.ColorNameBackground, 0) // 0 表示預設變體
	darkColor := darkTheme.Color(theme.ColorNameBackground, 0)

	// 驗證顏色不為 nil
	if lightColor == nil {
		t.Error("淺色主題背景顏色不應該為 nil")
	}
	if darkColor == nil {
		t.Error("深色主題背景顏色不應該為 nil")
	}

	// 測試字體方法
	font := lightTheme.Font(fyne.TextStyle{})
	if font == nil {
		t.Error("主題字體不應該為 nil")
	}

	// 測試圖示方法
	icon := lightTheme.Icon(theme.IconNameHome)
	if icon == nil {
		t.Error("主題圖示不應該為 nil")
	}

	// 測試尺寸方法
	size := lightTheme.Size(theme.SizeNameText)
	if size <= 0 {
		t.Error("主題文字尺寸應該大於 0")
	}
}