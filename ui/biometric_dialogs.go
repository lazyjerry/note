package ui

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// BiometricAuthStatus 生物驗證狀態
// 定義生物驗證過程中的不同狀態
type BiometricAuthStatus int

const (
	// BiometricStatusIdle 閒置狀態 - 尚未開始驗證
	BiometricStatusIdle BiometricAuthStatus = iota
	// BiometricStatusWaiting 等待中 - 正在等待用戶進行生物驗證
	BiometricStatusWaiting
	// BiometricStatusSuccess 成功 - 生物驗證成功
	BiometricStatusSuccess
	// BiometricStatusFailed 失敗 - 生物驗證失敗
	BiometricStatusFailed
	// BiometricStatusUnavailable 不可用 - 生物驗證功能不可用
	BiometricStatusUnavailable
)

// BiometricAuthResult 生物驗證結果
// 包含驗證結果和相關資訊
type BiometricAuthResult struct {
	Success       bool   // 驗證是否成功
	ErrorMessage  string // 錯誤訊息（如果有）
	FallbackUsed  bool   // 是否使用了密碼回退
	UserCancelled bool   // 用戶是否取消了驗證
}

// BiometricAuthCallback 生物驗證回調函數類型
// 當生物驗證完成時調用
type BiometricAuthCallback func(result BiometricAuthResult)

// BiometricAuthDialog 生物驗證對話框結構
// 用於顯示生物驗證提示和狀態的對話框
type BiometricAuthDialog struct {
	dialog          dialog.Dialog           // Fyne 對話框實例
	statusLabel     *widget.Label           // 狀態顯示標籤
	statusIcon      *widget.Icon            // 狀態圖示
	progressBar     *widget.ProgressBar     // 進度指示器
	fallbackButton  *widget.Button          // 回退到密碼驗證按鈕
	cancelButton    *widget.Button          // 取消按鈕
	callback        BiometricAuthCallback   // 完成時的回調函數
	status          BiometricAuthStatus     // 當前驗證狀態
	fallbackEnabled bool                    // 是否啟用密碼回退
}

// NewBiometricAuthDialog 創建新的生物驗證對話框
// 參數：
//   - parent: 父視窗，用於模態顯示
//   - title: 對話框標題
//   - message: 驗證提示訊息
//   - enableFallback: 是否啟用密碼回退
//   - callback: 完成時的回調函數
// 回傳：生物驗證對話框實例
//
// 執行流程：
// 1. 創建狀態顯示元件
// 2. 創建進度指示器
// 3. 創建操作按鈕
// 4. 組裝對話框佈局
// 5. 設置初始狀態
func NewBiometricAuthDialog(parent fyne.Window, title, message string, enableFallback bool, callback BiometricAuthCallback) *BiometricAuthDialog {
	d := &BiometricAuthDialog{
		callback:        callback,
		status:          BiometricStatusIdle,
		fallbackEnabled: enableFallback,
	}

	// 創建狀態標籤
	d.statusLabel = widget.NewLabel(message)
	d.statusLabel.Alignment = fyne.TextAlignCenter

	// 創建狀態圖示（初始為空）
	d.statusIcon = widget.NewIcon(nil)

	// 創建進度指示器
	d.progressBar = widget.NewProgressBar()
	d.progressBar.Hide() // 初始隱藏

	// 創建回退按鈕
	d.fallbackButton = widget.NewButton("使用密碼", func() {
		d.handleFallback()
	})
	if !enableFallback {
		d.fallbackButton.Hide()
	}

	// 創建取消按鈕
	d.cancelButton = widget.NewButton("取消", func() {
		d.handleCancel()
	})

	// 組裝按鈕容器
	var buttonContainer *fyne.Container
	if enableFallback {
		buttonContainer = container.NewHBox(
			d.fallbackButton,
			d.cancelButton,
		)
	} else {
		buttonContainer = container.NewHBox(
			d.cancelButton,
		)
	}

	// 組裝對話框內容
	content := container.NewVBox(
		container.NewCenter(d.statusIcon),
		widget.NewSeparator(),
		d.statusLabel,
		d.progressBar,
		widget.NewSeparator(),
		container.NewCenter(buttonContainer),
	)

	// 創建對話框
	d.dialog = dialog.NewCustom(title, "", content, parent)
	d.dialog.Resize(fyne.NewSize(400, 250))

	return d
}

// Show 顯示生物驗證對話框
// 將對話框以模態方式顯示給用戶
func (d *BiometricAuthDialog) Show() {
	d.dialog.Show()
}

// Hide 隱藏生物驗證對話框
// 關閉對話框
func (d *BiometricAuthDialog) Hide() {
	d.dialog.Hide()
}

// SetStatus 設置驗證狀態
// 參數：
//   - status: 新的驗證狀態
//   - message: 狀態訊息
//
// 執行流程：
// 1. 更新內部狀態
// 2. 更新狀態標籤文字
// 3. 更新狀態圖示
// 4. 更新進度指示器
// 5. 更新按鈕狀態
func (d *BiometricAuthDialog) SetStatus(status BiometricAuthStatus, message string) {
	d.status = status
	d.statusLabel.SetText(message)

	switch status {
	case BiometricStatusIdle:
		d.statusIcon.SetResource(nil)
		d.progressBar.Hide()
		d.enableButtons(true)

	case BiometricStatusWaiting:
		// 設置等待圖示（可以是指紋或 Face ID 圖示）
		d.statusIcon.SetResource(nil) // 這裡可以設置相應的圖示資源
		d.progressBar.Show()
		d.startProgressAnimation()
		d.enableButtons(true)

	case BiometricStatusSuccess:
		// 設置成功圖示
		d.statusIcon.SetResource(nil) // 這裡可以設置成功圖示
		d.progressBar.Hide()
		d.enableButtons(false)
		// 延遲關閉對話框
		time.AfterFunc(1*time.Second, func() {
			d.Hide()
		})

	case BiometricStatusFailed:
		// 設置失敗圖示
		d.statusIcon.SetResource(nil) // 這裡可以設置失敗圖示
		d.progressBar.Hide()
		d.enableButtons(true)

	case BiometricStatusUnavailable:
		// 設置不可用圖示
		d.statusIcon.SetResource(nil) // 這裡可以設置不可用圖示
		d.progressBar.Hide()
		d.enableButtons(true)
		// 如果生物驗證不可用且啟用了回退，自動顯示密碼選項
		if d.fallbackEnabled {
			d.fallbackButton.SetText("輸入密碼")
		}
	}
}

// startProgressAnimation 開始進度動畫
// 在等待生物驗證時顯示動畫效果
func (d *BiometricAuthDialog) startProgressAnimation() {
	// 創建一個簡單的進度動畫
	go func() {
		for d.status == BiometricStatusWaiting {
			for i := 0.0; i <= 1.0 && d.status == BiometricStatusWaiting; i += 0.1 {
				d.progressBar.SetValue(i)
				time.Sleep(100 * time.Millisecond)
			}
			for i := 1.0; i >= 0.0 && d.status == BiometricStatusWaiting; i -= 0.1 {
				d.progressBar.SetValue(i)
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
}

// enableButtons 啟用或禁用按鈕
// 參數：
//   - enabled: 是否啟用按鈕
func (d *BiometricAuthDialog) enableButtons(enabled bool) {
	if enabled {
		d.cancelButton.Enable()
		if d.fallbackEnabled {
			d.fallbackButton.Enable()
		}
	} else {
		d.cancelButton.Disable()
		if d.fallbackEnabled {
			d.fallbackButton.Disable()
		}
	}
}

// handleFallback 處理回退到密碼驗證
// 執行流程：
// 1. 調用回調函數通知使用密碼回退
// 2. 關閉對話框
func (d *BiometricAuthDialog) handleFallback() {
	if d.callback != nil {
		d.callback(BiometricAuthResult{
			Success:       false,
			ErrorMessage:  "",
			FallbackUsed:  true,
			UserCancelled: false,
		})
	}
	d.Hide()
}

// handleCancel 處理取消操作
// 執行流程：
// 1. 調用回調函數通知用戶取消
// 2. 關閉對話框
func (d *BiometricAuthDialog) handleCancel() {
	if d.callback != nil {
		d.callback(BiometricAuthResult{
			Success:       false,
			ErrorMessage:  "用戶取消了驗證",
			FallbackUsed:  false,
			UserCancelled: true,
		})
	}
	d.Hide()
}

// NotifySuccess 通知驗證成功
// 更新對話框狀態並調用回調函數
func (d *BiometricAuthDialog) NotifySuccess() {
	d.SetStatus(BiometricStatusSuccess, "驗證成功！")
	if d.callback != nil {
		d.callback(BiometricAuthResult{
			Success:       true,
			ErrorMessage:  "",
			FallbackUsed:  false,
			UserCancelled: false,
		})
	}
}

// NotifyFailure 通知驗證失敗
// 參數：
//   - errorMessage: 失敗原因
// 更新對話框狀態並調用回調函數
func (d *BiometricAuthDialog) NotifyFailure(errorMessage string) {
	d.SetStatus(BiometricStatusFailed, fmt.Sprintf("驗證失敗：%s", errorMessage))
	if d.callback != nil {
		d.callback(BiometricAuthResult{
			Success:       false,
			ErrorMessage:  errorMessage,
			FallbackUsed:  false,
			UserCancelled: false,
		})
	}
}

// NotifyUnavailable 通知生物驗證不可用
// 參數：
//   - reason: 不可用的原因
// 更新對話框狀態
func (d *BiometricAuthDialog) NotifyUnavailable(reason string) {
	d.SetStatus(BiometricStatusUnavailable, fmt.Sprintf("生物驗證不可用：%s", reason))
}

// BiometricSetupDialog 生物驗證設置對話框結構
// 用於設置和配置生物驗證功能
type BiometricSetupDialog struct {
	dialog              dialog.Dialog           // Fyne 對話框實例
	enableCheckbox      *widget.Check           // 啟用生物驗證複選框
	fallbackCheckbox    *widget.Check           // 啟用密碼回退複選框
	statusLabel         *widget.Label           // 狀態顯示標籤
	testButton          *widget.Button          // 測試生物驗證按鈕
	saveButton          *widget.Button          // 保存設置按鈕
	cancelButton        *widget.Button          // 取消按鈕
	callback            BiometricAuthCallback   // 完成時的回調函數
	biometricAvailable  bool                    // 生物驗證是否可用
}

// NewBiometricSetupDialog 創建新的生物驗證設置對話框
// 參數：
//   - parent: 父視窗，用於模態顯示
//   - title: 對話框標題
//   - biometricAvailable: 生物驗證是否可用
//   - callback: 完成時的回調函數
// 回傳：生物驗證設置對話框實例
//
// 執行流程：
// 1. 創建設置選項控制項
// 2. 創建狀態顯示
// 3. 創建操作按鈕
// 4. 組裝對話框佈局
// 5. 設置初始狀態
func NewBiometricSetupDialog(parent fyne.Window, title string, biometricAvailable bool, callback BiometricAuthCallback) *BiometricSetupDialog {
	d := &BiometricSetupDialog{
		callback:           callback,
		biometricAvailable: biometricAvailable,
	}

	// 創建啟用生物驗證複選框
	d.enableCheckbox = widget.NewCheck("啟用生物驗證 (Touch ID/Face ID)", func(checked bool) {
		d.updateUI()
	})
	d.enableCheckbox.SetChecked(false)

	// 創建密碼回退複選框
	d.fallbackCheckbox = widget.NewCheck("允許密碼回退", func(checked bool) {
		// 回退選項變更處理
	})
	d.fallbackCheckbox.SetChecked(true)

	// 創建狀態標籤
	d.statusLabel = widget.NewLabel("")
	d.updateStatusLabel()

	// 創建測試按鈕
	d.testButton = widget.NewButton("測試生物驗證", func() {
		d.handleTest()
	})

	// 創建保存按鈕
	d.saveButton = widget.NewButton("保存", func() {
		d.handleSave()
	})

	// 創建取消按鈕
	d.cancelButton = widget.NewButton("取消", func() {
		d.handleCancel()
	})

	// 組裝對話框內容
	content := container.NewVBox(
		widget.NewLabel("生物驗證設置："),
		widget.NewSeparator(),
		d.enableCheckbox,
		d.fallbackCheckbox,
		widget.NewSeparator(),
		d.statusLabel,
		d.testButton,
		widget.NewSeparator(),
		container.NewHBox(
			d.saveButton,
			d.cancelButton,
		),
	)

	// 創建對話框
	d.dialog = dialog.NewCustom(title, "", content, parent)
	d.dialog.Resize(fyne.NewSize(400, 300))

	// 初始化 UI 狀態
	d.updateUI()

	return d
}

// Show 顯示生物驗證設置對話框
// 將對話框以模態方式顯示給用戶
func (d *BiometricSetupDialog) Show() {
	d.dialog.Show()
}

// Hide 隱藏生物驗證設置對話框
// 關閉對話框
func (d *BiometricSetupDialog) Hide() {
	d.dialog.Hide()
}

// updateUI 更新 UI 狀態
// 根據生物驗證可用性和用戶選擇更新界面
func (d *BiometricSetupDialog) updateUI() {
	if !d.biometricAvailable {
		d.enableCheckbox.SetChecked(false)
		d.enableCheckbox.Disable()
		d.testButton.Disable()
		d.fallbackCheckbox.Disable()
	} else {
		d.enableCheckbox.Enable()
		if d.enableCheckbox.Checked {
			d.testButton.Enable()
			d.fallbackCheckbox.Enable()
		} else {
			d.testButton.Disable()
			d.fallbackCheckbox.Disable()
		}
	}
	d.updateStatusLabel()
}

// updateStatusLabel 更新狀態標籤
// 根據當前設置顯示相應的狀態訊息
func (d *BiometricSetupDialog) updateStatusLabel() {
	if !d.biometricAvailable {
		d.statusLabel.SetText("⚠️ 此設備不支援生物驗證或未設置 Touch ID/Face ID")
	} else if d.enableCheckbox.Checked {
		d.statusLabel.SetText("✅ 生物驗證已啟用")
	} else {
		d.statusLabel.SetText("ℹ️ 生物驗證已禁用")
	}
}

// handleTest 處理測試生物驗證
// 執行流程：
// 1. 創建測試用的生物驗證對話框
// 2. 模擬生物驗證過程
// 3. 顯示測試結果
func (d *BiometricSetupDialog) handleTest() {
	// 創建測試對話框
	testDialog := NewBiometricAuthDialog(
		d.dialog.(fyne.Window),
		"測試生物驗證",
		"請使用 Touch ID 或 Face ID 進行驗證",
		true,
		func(result BiometricAuthResult) {
			if result.Success {
				dialog.ShowInformation("測試成功", "生物驗證測試成功！", d.dialog.(fyne.Window))
			} else if result.FallbackUsed {
				dialog.ShowInformation("回退測試", "密碼回退功能正常", d.dialog.(fyne.Window))
			} else if result.UserCancelled {
				dialog.ShowInformation("測試取消", "用戶取消了測試", d.dialog.(fyne.Window))
			} else {
				dialog.ShowError(fmt.Errorf("測試失敗：%s", result.ErrorMessage), d.dialog.(fyne.Window))
			}
		},
	)

	// 顯示測試對話框
	testDialog.Show()

	// 模擬生物驗證過程
	testDialog.SetStatus(BiometricStatusWaiting, "正在等待生物驗證...")

	// 模擬驗證結果（在實際實作中，這裡會調用真正的生物驗證 API）
	go func() {
		time.Sleep(2 * time.Second)
		// 這裡可以隨機模擬成功或失敗
		testDialog.NotifySuccess()
	}()
}

// handleSave 處理保存設置
// 執行流程：
// 1. 收集用戶設置
// 2. 調用回調函數保存設置
// 3. 關閉對話框
func (d *BiometricSetupDialog) handleSave() {
	// 這裡可以收集設置並通過回調函數傳遞
	if d.callback != nil {
		d.callback(BiometricAuthResult{
			Success:       true,
			ErrorMessage:  "",
			FallbackUsed:  false,
			UserCancelled: false,
		})
	}
	d.Hide()
}

// handleCancel 處理取消操作
// 執行流程：
// 1. 調用回調函數通知取消
// 2. 關閉對話框
func (d *BiometricSetupDialog) handleCancel() {
	if d.callback != nil {
		d.callback(BiometricAuthResult{
			Success:       false,
			ErrorMessage:  "",
			FallbackUsed:  false,
			UserCancelled: true,
		})
	}
	d.Hide()
}

// ShowBiometricAuthDialog 顯示生物驗證對話框的便利函數
// 參數：
//   - parent: 父視窗
//   - title: 對話框標題
//   - message: 驗證提示訊息
//   - enableFallback: 是否啟用密碼回退
//   - callback: 完成時的回調函數
//
// 執行流程：
// 1. 創建生物驗證對話框
// 2. 顯示對話框
// 3. 開始驗證過程
func ShowBiometricAuthDialog(parent fyne.Window, title, message string, enableFallback bool, callback BiometricAuthCallback) *BiometricAuthDialog {
	dialog := NewBiometricAuthDialog(parent, title, message, enableFallback, callback)
	dialog.Show()
	dialog.SetStatus(BiometricStatusWaiting, message)
	return dialog
}

// ShowBiometricSetupDialog 顯示生物驗證設置對話框的便利函數
// 參數：
//   - parent: 父視窗
//   - title: 對話框標題
//   - biometricAvailable: 生物驗證是否可用
//   - callback: 完成時的回調函數
//
// 執行流程：
// 1. 創建生物驗證設置對話框
// 2. 顯示對話框
func ShowBiometricSetupDialog(parent fyne.Window, title string, biometricAvailable bool, callback BiometricAuthCallback) {
	dialog := NewBiometricSetupDialog(parent, title, biometricAvailable, callback)
	dialog.Show()
}