// Package services 提供應用程式的核心業務邏輯服務
// 本檔案實作主題管理相關功能
package services

import (
	"image/color"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"

	"mac-notebook-app/internal/models"
)

// ThemeService 主題管理服務
// 負責處理應用程式主題的切換、偵測和管理
type ThemeService struct {
	app              fyne.App           // Fyne 應用程式實例
	currentTheme     string             // 當前主題設定
	systemTheme      string             // 系統主題狀態
	themeListeners   []ThemeListener    // 主題變更監聽器列表
	mutex            sync.RWMutex       // 讀寫鎖，保護並發存取
	settings         *models.Settings   // 應用程式設定
}

// ThemeListener 主題變更監聽器介面
// 實作此介面的元件可以接收主題變更通知
type ThemeListener interface {
	OnThemeChanged(themeName string)
}

// NewThemeService 建立新的主題管理服務實例
// 參數：
//   - app: Fyne 應用程式實例
//   - settings: 應用程式設定
// 回傳：新建立的主題服務實例
//
// 執行流程：
// 1. 建立主題服務結構體
// 2. 初始化主題監聽器列表
// 3. 設定當前主題為設定中的主題
// 4. 偵測系統主題狀態
// 5. 套用初始主題
func NewThemeService(app fyne.App, settings *models.Settings) *ThemeService {
	service := &ThemeService{
		app:            app,
		settings:       settings,
		currentTheme:   settings.Theme,
		themeListeners: make([]ThemeListener, 0),
	}
	
	// 偵測系統主題
	service.detectSystemTheme()
	
	// 套用初始主題
	service.applyTheme(settings.Theme)
	
	return service
}

// detectSystemTheme 偵測系統主題設定
// 執行流程：
// 1. 檢查系統是否支援深色模式偵測
// 2. 根據系統設定判斷當前系統主題
// 3. 更新內部系統主題狀態
func (ts *ThemeService) detectSystemTheme() {
	ts.mutex.Lock()
	defer ts.mutex.Unlock()
	
	// 檢查當前系統主題
	// 注意：Fyne 會自動偵測系統主題，我們可以透過比較當前主題來判斷
	currentVariant := ts.app.Settings().ThemeVariant()
	
	switch currentVariant {
	case theme.VariantDark:
		ts.systemTheme = "dark"
	case theme.VariantLight:
		ts.systemTheme = "light"
	default:
		ts.systemTheme = "light" // 預設為淺色主題
	}
}

// GetCurrentTheme 取得當前主題設定
// 回傳：當前主題名稱（"light"、"dark" 或 "auto"）
//
// 執行流程：
// 1. 使用讀鎖保護並發存取
// 2. 回傳當前主題設定
func (ts *ThemeService) GetCurrentTheme() string {
	ts.mutex.RLock()
	defer ts.mutex.RUnlock()
	return ts.currentTheme
}

// GetSystemTheme 取得系統主題狀態
// 回傳：系統主題名稱（"light" 或 "dark"）
//
// 執行流程：
// 1. 重新偵測系統主題（確保最新狀態）
// 2. 使用讀鎖保護並發存取
// 3. 回傳系統主題狀態
func (ts *ThemeService) GetSystemTheme() string {
	ts.detectSystemTheme()
	ts.mutex.RLock()
	defer ts.mutex.RUnlock()
	return ts.systemTheme
}

// SetTheme 設定應用程式主題
// 參數：
//   - themeName: 主題名稱（"light"、"dark" 或 "auto"）
// 回傳：設定操作的錯誤（如果有）
//
// 執行流程：
// 1. 驗證主題名稱是否有效
// 2. 更新設定中的主題值
// 3. 套用新主題到應用程式
// 4. 通知所有主題監聽器
// 5. 保存設定到檔案
func (ts *ThemeService) SetTheme(themeName string) error {
	// 驗證主題名稱
	if err := ts.settings.UpdateTheme(themeName); err != nil {
		return err
	}
	
	ts.mutex.Lock()
	ts.currentTheme = themeName
	ts.mutex.Unlock()
	
	// 套用主題
	ts.applyTheme(themeName)
	
	// 通知監聽器
	ts.notifyThemeListeners(themeName)
	
	// 保存設定
	return ts.settings.SaveDefault()
}

// applyTheme 套用指定主題到應用程式
// 參數：
//   - themeName: 要套用的主題名稱
//
// 執行流程：
// 1. 根據主題名稱決定實際要套用的主題
// 2. 如果是 "auto"，則根據系統主題決定
// 3. 設定 Fyne 應用程式的主題變體
func (ts *ThemeService) applyTheme(themeName string) {
	var variant fyne.ThemeVariant
	
	switch themeName {
	case "light":
		variant = theme.VariantLight
	case "dark":
		variant = theme.VariantDark
	case "auto":
		// 自動模式：跟隨系統主題
		ts.detectSystemTheme()
		if ts.systemTheme == "dark" {
			variant = theme.VariantDark
		} else {
			variant = theme.VariantLight
		}
	default:
		variant = theme.VariantLight // 預設為淺色主題
	}
	
	// 套用主題變體到應用程式
	ts.app.Settings().SetTheme(&customTheme{variant: variant})
}

// AddThemeListener 新增主題變更監聽器
// 參數：
//   - listener: 實作 ThemeListener 介面的監聽器
//
// 執行流程：
// 1. 使用寫鎖保護並發存取
// 2. 將監聽器加入監聽器列表
func (ts *ThemeService) AddThemeListener(listener ThemeListener) {
	ts.mutex.Lock()
	defer ts.mutex.Unlock()
	ts.themeListeners = append(ts.themeListeners, listener)
}

// RemoveThemeListener 移除主題變更監聽器
// 參數：
//   - listener: 要移除的監聽器
//
// 執行流程：
// 1. 使用寫鎖保護並發存取
// 2. 在監聽器列表中尋找指定監聽器
// 3. 如果找到，從列表中移除
func (ts *ThemeService) RemoveThemeListener(listener ThemeListener) {
	ts.mutex.Lock()
	defer ts.mutex.Unlock()
	
	for i, l := range ts.themeListeners {
		if l == listener {
			// 移除監聽器（保持切片順序）
			ts.themeListeners = append(ts.themeListeners[:i], ts.themeListeners[i+1:]...)
			break
		}
	}
}

// notifyThemeListeners 通知所有主題監聽器主題已變更
// 參數：
//   - themeName: 新的主題名稱
//
// 執行流程：
// 1. 使用讀鎖保護並發存取
// 2. 遍歷所有監聽器並呼叫其 OnThemeChanged 方法
func (ts *ThemeService) notifyThemeListeners(themeName string) {
	ts.mutex.RLock()
	listeners := make([]ThemeListener, len(ts.themeListeners))
	copy(listeners, ts.themeListeners)
	ts.mutex.RUnlock()
	
	// 在鎖外通知監聽器，避免死鎖
	for _, listener := range listeners {
		listener.OnThemeChanged(themeName)
	}
}

// GetAvailableThemes 取得可用的主題列表
// 回傳：包含所有可用主題名稱的字串切片
//
// 執行流程：
// 1. 回傳預定義的主題列表
func (ts *ThemeService) GetAvailableThemes() []string {
	return []string{"light", "dark", "auto"}
}

// GetThemeDisplayName 取得主題的顯示名稱
// 參數：
//   - themeName: 主題名稱
// 回傳：主題的中文顯示名稱
//
// 執行流程：
// 1. 根據主題名稱回傳對應的中文顯示名稱
func (ts *ThemeService) GetThemeDisplayName(themeName string) string {
	switch themeName {
	case "light":
		return "淺色主題"
	case "dark":
		return "深色主題"
	case "auto":
		return "自動（跟隨系統）"
	default:
		return "未知主題"
	}
}

// IsSystemThemeSupported 檢查是否支援系統主題偵測
// 回傳：如果支援系統主題偵測則回傳 true
//
// 執行流程：
// 1. 檢查當前平台是否支援系統主題偵測
// 2. 在 macOS 上通常支援，其他平台可能不支援
func (ts *ThemeService) IsSystemThemeSupported() bool {
	// 在 macOS 上，Fyne 支援系統主題偵測
	// 其他平台的支援情況可能不同
	return true
}

// RefreshSystemTheme 重新整理系統主題狀態
// 執行流程：
// 1. 重新偵測系統主題
// 2. 如果當前設定為自動模式，重新套用主題
// 3. 通知監聽器主題可能已變更
func (ts *ThemeService) RefreshSystemTheme() {
	ts.detectSystemTheme()
	
	ts.mutex.RLock()
	currentTheme := ts.currentTheme
	ts.mutex.RUnlock()
	
	if currentTheme == "auto" {
		ts.applyTheme("auto")
		ts.notifyThemeListeners("auto")
	}
}

// customTheme 自訂主題結構體
// 實作 fyne.Theme 介面，提供主題變體控制
type customTheme struct {
	variant fyne.ThemeVariant
}

// Color 實作 fyne.Theme 介面的 Color 方法
// 參數：
//   - name: 顏色名稱
//   - variant: 主題變體（可選）
// 回傳：對應的顏色值
func (ct *customTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	// 使用指定的變體，如果沒有指定則使用預設變體
	if variant == 0 { // 0 表示預設變體
		variant = ct.variant
	}
	return theme.DefaultTheme().Color(name, variant)
}

// Font 實作 fyne.Theme 介面的 Font 方法
// 參數：
//   - style: 字體樣式
// 回傳：對應的字體資源
func (ct *customTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

// Icon 實作 fyne.Theme 介面的 Icon 方法
// 參數：
//   - name: 圖示名稱
// 回傳：對應的圖示資源
func (ct *customTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

// Size 實作 fyne.Theme 介面的 Size 方法
// 參數：
//   - name: 尺寸名稱
// 回傳：對應的尺寸值
func (ct *customTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}