// Package ui 提供使用者介面元件
// 本檔案實作設定對話框相關功能
package ui

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"

	"mac-notebook-app/internal/models"
)

// SettingsDialog 代表設定對話框的結構體
// 包含所有設定相關的 UI 元件和狀態管理
type SettingsDialog struct {
	window   fyne.Window                // 父視窗
	settings *models.Settings           // 當前設定實例
	dialog   *dialog.CustomDialog       // 自訂對話框
	
	// UI 元件
	encryptionSelect    *widget.Select   // 加密演算法選擇器
	autoSaveEntry      *widget.Entry     // 自動保存間隔輸入框
	saveLocationEntry  *widget.Entry     // 預設保存位置輸入框
	biometricCheck     *widget.Check     // 生物識別啟用勾選框
	themeSelect        *widget.Select    // 主題選擇器
	
	// 回調函數
	onSettingsChanged func(*models.Settings) // 設定變更時的回調函數
}

// NewSettingsDialog 建立新的設定對話框實例
// 參數：
//   - parent: 父視窗
//   - currentSettings: 當前的設定實例
//   - onChanged: 設定變更時的回調函數
// 回傳：新建立的設定對話框實例
//
// 執行流程：
// 1. 建立設定對話框結構體
// 2. 初始化所有 UI 元件
// 3. 設定當前值到各個元件
// 4. 建立對話框佈局
// 5. 回傳完整的對話框實例
func NewSettingsDialog(parent fyne.Window, currentSettings *models.Settings, onChanged func(*models.Settings)) *SettingsDialog {
	// 建立設定對話框實例
	sd := &SettingsDialog{
		window:            parent,
		settings:          currentSettings.Clone(), // 使用設定的複製避免直接修改
		onSettingsChanged: onChanged,
	}
	
	// 初始化 UI 元件
	sd.initializeComponents()
	
	// 建立對話框內容
	content := sd.createContent()
	
	// 建立自訂對話框
	sd.dialog = dialog.NewCustom("應用程式設定", "關閉", content, parent)
	sd.dialog.Resize(fyne.NewSize(500, 400))
	
	return sd
}

// initializeComponents 初始化所有 UI 元件
// 執行流程：
// 1. 建立加密演算法選擇器
// 2. 建立自動保存間隔輸入框
// 3. 建立預設保存位置輸入框和瀏覽按鈕
// 4. 建立生物識別勾選框
// 5. 建立主題選擇器
// 6. 設定各元件的當前值
func (sd *SettingsDialog) initializeComponents() {
	// 建立加密演算法選擇器
	sd.encryptionSelect = widget.NewSelect(
		sd.settings.GetSupportedEncryptionAlgorithms(),
		func(value string) {
			// 當選擇變更時更新設定
			if err := sd.settings.UpdateEncryption(value); err == nil {
				sd.notifySettingsChanged()
			}
		},
	)
	sd.encryptionSelect.SetSelected(sd.settings.DefaultEncryption)
	
	// 建立自動保存間隔輸入框
	sd.autoSaveEntry = widget.NewEntry()
	sd.autoSaveEntry.SetText(strconv.Itoa(sd.settings.AutoSaveInterval))
	sd.autoSaveEntry.OnChanged = func(text string) {
		// 驗證並更新自動保存間隔
		if interval, err := strconv.Atoi(text); err == nil {
			if err := sd.settings.UpdateAutoSaveInterval(interval); err == nil {
				sd.notifySettingsChanged()
			}
		}
	}
	
	// 建立預設保存位置輸入框
	sd.saveLocationEntry = widget.NewEntry()
	sd.saveLocationEntry.SetText(sd.settings.DefaultSaveLocation)
	sd.saveLocationEntry.OnChanged = func(text string) {
		// 更新預設保存位置
		sd.settings.UpdateDefaultSaveLocation(text)
		sd.notifySettingsChanged()
	}
	
	// 建立生物識別勾選框
	sd.biometricCheck = widget.NewCheck("啟用生物識別驗證 (Touch ID/Face ID)", func(checked bool) {
		// 更新生物識別設定
		sd.settings.SetBiometric(checked)
		sd.notifySettingsChanged()
	})
	sd.biometricCheck.SetChecked(sd.settings.BiometricEnabled)
	
	// 建立主題選擇器
	sd.themeSelect = widget.NewSelect(
		sd.settings.GetSupportedThemes(),
		func(value string) {
			// 當主題選擇變更時更新設定
			if err := sd.settings.UpdateTheme(value); err == nil {
				sd.notifySettingsChanged()
			}
		},
	)
	sd.themeSelect.SetSelected(sd.settings.Theme)
}

// createContent 建立對話框的內容佈局
// 回傳：包含所有設定元件的容器
//
// 執行流程：
// 1. 建立加密設定區塊
// 2. 建立檔案管理設定區塊
// 3. 建立外觀設定區塊
// 4. 建立操作按鈕區塊
// 5. 組合所有區塊成為完整佈局
func (sd *SettingsDialog) createContent() *fyne.Container {
	// 建立加密設定區塊
	encryptionSection := sd.createEncryptionSection()
	
	// 建立檔案管理設定區塊
	fileSection := sd.createFileSection()
	
	// 建立外觀設定區塊
	appearanceSection := sd.createAppearanceSection()
	
	// 建立操作按鈕區塊
	buttonSection := sd.createButtonSection()
	
	// 組合所有區塊
	content := container.NewVBox(
		encryptionSection,
		widget.NewSeparator(),
		fileSection,
		widget.NewSeparator(),
		appearanceSection,
		widget.NewSeparator(),
		buttonSection,
	)
	
	return content
}

// createEncryptionSection 建立加密設定區塊
// 回傳：包含加密相關設定的容器
//
// 執行流程：
// 1. 建立區塊標題
// 2. 建立加密演算法選擇器佈局
// 3. 建立生物識別設定佈局
// 4. 組合成完整的加密設定區塊
func (sd *SettingsDialog) createEncryptionSection() *fyne.Container {
	// 區塊標題
	title := widget.NewRichTextFromMarkdown("## 🔐 加密設定")
	
	// 加密演算法選擇
	encryptionLabel := widget.NewLabel("預設加密演算法：")
	encryptionRow := container.NewBorder(nil, nil, encryptionLabel, nil, sd.encryptionSelect)
	
	// 加密演算法說明
	encryptionHelp := widget.NewRichTextFromMarkdown(`
**AES-256**: 廣泛使用的標準加密演算法，相容性佳
**ChaCha20**: 現代化加密演算法，效能較佳`)
	encryptionHelp.Wrapping = fyne.TextWrapWord
	
	// 生物識別設定
	biometricRow := container.NewHBox(sd.biometricCheck)
	
	// 組合加密設定區塊
	section := container.NewVBox(
		title,
		encryptionRow,
		encryptionHelp,
		biometricRow,
	)
	
	return section
}

// createFileSection 建立檔案管理設定區塊
// 回傳：包含檔案管理相關設定的容器
//
// 執行流程：
// 1. 建立區塊標題
// 2. 建立自動保存間隔設定佈局
// 3. 建立預設保存位置設定佈局
// 4. 組合成完整的檔案管理設定區塊
func (sd *SettingsDialog) createFileSection() *fyne.Container {
	// 區塊標題
	title := widget.NewRichTextFromMarkdown("## 📁 檔案管理")
	
	// 自動保存間隔設定
	autoSaveLabel := widget.NewLabel("自動保存間隔（分鐘）：")
	autoSaveHelp := widget.NewLabel("範圍：1-60 分鐘")
	autoSaveRow := container.NewBorder(nil, nil, autoSaveLabel, autoSaveHelp, sd.autoSaveEntry)
	
	// 預設保存位置設定
	saveLocationLabel := widget.NewLabel("預設保存位置：")
	browseButton := widget.NewButton("瀏覽...", sd.onBrowseLocation)
	saveLocationRow := container.NewBorder(nil, nil, saveLocationLabel, browseButton, sd.saveLocationEntry)
	
	// 組合檔案管理設定區塊
	section := container.NewVBox(
		title,
		autoSaveRow,
		saveLocationRow,
	)
	
	return section
}

// createAppearanceSection 建立外觀設定區塊
// 回傳：包含外觀相關設定的容器
//
// 執行流程：
// 1. 建立區塊標題
// 2. 建立主題選擇器佈局
// 3. 建立主題說明
// 4. 組合成完整的外觀設定區塊
func (sd *SettingsDialog) createAppearanceSection() *fyne.Container {
	// 區塊標題
	title := widget.NewRichTextFromMarkdown("## 🎨 外觀設定")
	
	// 主題選擇
	themeLabel := widget.NewLabel("應用程式主題：")
	themeRow := container.NewBorder(nil, nil, themeLabel, nil, sd.themeSelect)
	
	// 主題說明
	themeHelp := widget.NewRichTextFromMarkdown(`
**淺色 (Light)**: 使用淺色主題
**深色 (Dark)**: 使用深色主題  
**自動 (Auto)**: 跟隨系統主題設定`)
	themeHelp.Wrapping = fyne.TextWrapWord
	
	// 組合外觀設定區塊
	section := container.NewVBox(
		title,
		themeRow,
		themeHelp,
	)
	
	return section
}

// createButtonSection 建立操作按鈕區塊
// 回傳：包含操作按鈕的容器
//
// 執行流程：
// 1. 建立重設為預設值按鈕
// 2. 建立儲存設定按鈕
// 3. 組合成按鈕列佈局
func (sd *SettingsDialog) createButtonSection() *fyne.Container {
	// 重設為預設值按鈕
	resetButton := widget.NewButton("重設為預設值", sd.onResetToDefaults)
	
	// 儲存設定按鈕
	saveButton := widget.NewButton("儲存設定", sd.onSaveSettings)
	saveButton.Importance = widget.HighImportance
	
	// 組合按鈕列
	buttonRow := container.NewHBox(
		resetButton,
		widget.NewSeparator(),
		saveButton,
	)
	
	return buttonRow
}

// onBrowseLocation 處理瀏覽保存位置按鈕點擊事件
// 執行流程：
// 1. 開啟資料夾選擇對話框
// 2. 當使用者選擇資料夾時更新設定
// 3. 更新 UI 顯示的路徑
func (sd *SettingsDialog) onBrowseLocation() {
	// 建立資料夾選擇對話框
	folderDialog := dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
		if err != nil {
			return
		}
		if uri != nil {
			// 更新預設保存位置
			path := uri.Path()
			sd.saveLocationEntry.SetText(path)
			sd.settings.UpdateDefaultSaveLocation(path)
			sd.notifySettingsChanged()
		}
	}, sd.window)
	
	// 設定對話框標題和初始位置
	if uri := storage.NewFileURI(sd.settings.DefaultSaveLocation); uri != nil {
		if listableURI, ok := uri.(fyne.ListableURI); ok {
			folderDialog.SetLocation(listableURI)
		}
	}
	folderDialog.Show()
}

// onResetToDefaults 處理重設為預設值按鈕點擊事件
// 執行流程：
// 1. 顯示確認對話框
// 2. 如果使用者確認，重設所有設定為預設值
// 3. 更新所有 UI 元件顯示
// 4. 通知設定變更
func (sd *SettingsDialog) onResetToDefaults() {
	// 建立確認對話框
	confirmDialog := dialog.NewConfirm(
		"確認重設",
		"確定要將所有設定重設為預設值嗎？此操作無法復原。",
		func(confirmed bool) {
			if confirmed {
				// 重設設定為預設值
				sd.settings = models.NewDefaultSettings()
				
				// 更新所有 UI 元件
				sd.updateUIFromSettings()
				
				// 通知設定變更
				sd.notifySettingsChanged()
			}
		},
		sd.window,
	)
	
	confirmDialog.Show()
}

// onSaveSettings 處理儲存設定按鈕點擊事件
// 執行流程：
// 1. 驗證當前設定是否有效
// 2. 嘗試保存設定到檔案
// 3. 顯示保存結果訊息
// 4. 如果保存成功，關閉對話框
func (sd *SettingsDialog) onSaveSettings() {
	// 驗證設定
	if err := sd.settings.Validate(); err != nil {
		// 顯示驗證錯誤訊息
		dialog.ShowError(fmt.Errorf("設定驗證失敗：%v", err), sd.window)
		return
	}
	
	// 保存設定到檔案
	if err := sd.settings.SaveDefault(); err != nil {
		// 顯示保存錯誤訊息
		dialog.ShowError(fmt.Errorf("保存設定失敗：%v", err), sd.window)
		return
	}
	
	// 顯示保存成功訊息
	dialog.ShowInformation("設定已保存", "設定已成功保存並套用。", sd.window)
	
	// 關閉對話框
	sd.dialog.Hide()
}

// updateUIFromSettings 根據當前設定更新所有 UI 元件
// 執行流程：
// 1. 更新加密演算法選擇器
// 2. 更新自動保存間隔輸入框
// 3. 更新預設保存位置輸入框
// 4. 更新生物識別勾選框
// 5. 更新主題選擇器
func (sd *SettingsDialog) updateUIFromSettings() {
	sd.encryptionSelect.SetSelected(sd.settings.DefaultEncryption)
	sd.autoSaveEntry.SetText(strconv.Itoa(sd.settings.AutoSaveInterval))
	sd.saveLocationEntry.SetText(sd.settings.DefaultSaveLocation)
	sd.biometricCheck.SetChecked(sd.settings.BiometricEnabled)
	sd.themeSelect.SetSelected(sd.settings.Theme)
}

// notifySettingsChanged 通知設定變更
// 執行流程：
// 1. 檢查是否有設定變更回調函數
// 2. 如果有，呼叫回調函數並傳遞當前設定
func (sd *SettingsDialog) notifySettingsChanged() {
	if sd.onSettingsChanged != nil {
		sd.onSettingsChanged(sd.settings)
	}
}

// Show 顯示設定對話框
// 執行流程：
// 1. 顯示對話框
func (sd *SettingsDialog) Show() {
	sd.dialog.Show()
}

// Hide 隱藏設定對話框
// 執行流程：
// 1. 隱藏對話框
func (sd *SettingsDialog) Hide() {
	sd.dialog.Hide()
}

// GetSettings 取得當前的設定實例
// 回傳：當前的設定實例
//
// 執行流程：
// 1. 回傳當前設定的複製
func (sd *SettingsDialog) GetSettings() *models.Settings {
	return sd.settings.Clone()
}

// SetSettings 設定新的設定值並更新 UI
// 參數：
//   - newSettings: 新的設定實例
//
// 執行流程：
// 1. 更新內部設定實例
// 2. 更新所有 UI 元件顯示
func (sd *SettingsDialog) SetSettings(newSettings *models.Settings) {
	sd.settings = newSettings.Clone()
	sd.updateUIFromSettings()
}