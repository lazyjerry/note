package ui

import (
	"fmt"
	"unicode"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// PasswordStrength 密碼強度等級
// 定義密碼強度的不同等級，用於密碼強度指示器
type PasswordStrength int

const (
	// PasswordWeak 弱密碼 - 不符合基本安全要求
	PasswordWeak PasswordStrength = iota
	// PasswordMedium 中等密碼 - 符合基本要求但可以更強
	PasswordMedium
	// PasswordStrong 強密碼 - 符合所有安全要求
	PasswordStrong
)

// PasswordDialogResult 密碼對話框結果
// 包含用戶輸入的密碼和操作結果
type PasswordDialogResult struct {
	Password  string // 用戶輸入的密碼
	Confirmed bool   // 用戶是否確認操作
}

// PasswordDialogCallback 密碼對話框回調函數類型
// 當用戶完成密碼輸入操作時調用
type PasswordDialogCallback func(result PasswordDialogResult)

// PasswordSetupDialog 密碼設定對話框結構
// 用於創建新密碼或修改現有密碼的對話框
type PasswordSetupDialog struct {
	dialog          dialog.Dialog        // Fyne 對話框實例
	passwordEntry   *widget.Entry        // 密碼輸入框
	confirmEntry    *widget.Entry        // 確認密碼輸入框
	strengthBar     *widget.ProgressBar  // 密碼強度指示器
	strengthLabel   *widget.Label        // 密碼強度文字說明
	callback        PasswordDialogCallback // 完成時的回調函數
}

// NewPasswordSetupDialog 創建新的密碼設定對話框
// 參數：
//   - parent: 父視窗，用於模態顯示
//   - title: 對話框標題
//   - callback: 完成時的回調函數
// 回傳：密碼設定對話框實例
//
// 執行流程：
// 1. 創建密碼輸入框和確認輸入框
// 2. 創建密碼強度指示器
// 3. 設置輸入框的變更監聽器
// 4. 創建確認和取消按鈕
// 5. 組裝對話框佈局
func NewPasswordSetupDialog(parent fyne.Window, title string, callback PasswordDialogCallback) *PasswordSetupDialog {
	d := &PasswordSetupDialog{
		callback: callback,
	}

	// 創建密碼輸入框
	d.passwordEntry = widget.NewPasswordEntry()
	d.passwordEntry.SetPlaceHolder("請輸入密碼")
	
	// 創建確認密碼輸入框
	d.confirmEntry = widget.NewPasswordEntry()
	d.confirmEntry.SetPlaceHolder("請再次輸入密碼")

	// 創建密碼強度指示器
	d.strengthBar = widget.NewProgressBar()
	d.strengthBar.SetValue(0)
	
	// 創建密碼強度文字說明
	d.strengthLabel = widget.NewLabel("密碼強度：無")

	// 設置密碼輸入框的變更監聽器
	d.passwordEntry.OnChanged = func(text string) {
		d.updatePasswordStrength(text)
	}

	// 創建確認按鈕
	confirmButton := widget.NewButton("確認", func() {
		d.handleConfirm()
	})

	// 創建取消按鈕
	cancelButton := widget.NewButton("取消", func() {
		d.handleCancel()
	})

	// 組裝對話框內容
	content := container.NewVBox(
		widget.NewLabel("設定密碼："),
		d.passwordEntry,
		widget.NewLabel("確認密碼："),
		d.confirmEntry,
		widget.NewSeparator(),
		widget.NewLabel("密碼強度："),
		d.strengthBar,
		d.strengthLabel,
		widget.NewSeparator(),
		container.NewHBox(
			confirmButton,
			cancelButton,
		),
	)

	// 創建對話框
	d.dialog = dialog.NewCustom(title, "", content, parent)
	d.dialog.Resize(fyne.NewSize(400, 300))

	return d
}

// Show 顯示密碼設定對話框
// 將對話框以模態方式顯示給用戶
func (d *PasswordSetupDialog) Show() {
	d.dialog.Show()
}

// updatePasswordStrength 更新密碼強度指示器
// 參數：
//   - password: 要檢查的密碼
//
// 執行流程：
// 1. 計算密碼強度等級
// 2. 更新進度條數值
// 3. 更新強度文字說明
// 4. 根據強度設置不同的顏色主題
func (d *PasswordSetupDialog) updatePasswordStrength(password string) {
	strength := calculatePasswordStrength(password)
	
	// 根據強度等級設置進度條數值
	switch strength {
	case PasswordWeak:
		d.strengthBar.SetValue(0.33)
		d.strengthLabel.SetText("密碼強度：弱")
	case PasswordMedium:
		d.strengthBar.SetValue(0.66)
		d.strengthLabel.SetText("密碼強度：中等")
	case PasswordStrong:
		d.strengthBar.SetValue(1.0)
		d.strengthLabel.SetText("密碼強度：強")
	}
}

// handleConfirm 處理確認按鈕點擊事件
// 執行流程：
// 1. 驗證密碼輸入
// 2. 檢查密碼確認是否一致
// 3. 驗證密碼強度
// 4. 調用回調函數並關閉對話框
func (d *PasswordSetupDialog) handleConfirm() {
	password := d.passwordEntry.Text
	confirm := d.confirmEntry.Text

	// 檢查密碼是否為空
	if password == "" {
		dialog.ShowError(fmt.Errorf("密碼不能為空"), d.dialog.(fyne.Window))
		return
	}

	// 檢查密碼確認是否一致
	if password != confirm {
		dialog.ShowError(fmt.Errorf("兩次輸入的密碼不一致"), d.dialog.(fyne.Window))
		return
	}

	// 檢查密碼強度
	if calculatePasswordStrength(password) == PasswordWeak {
		dialog.ShowError(fmt.Errorf("密碼強度太弱，請使用更複雜的密碼"), d.dialog.(fyne.Window))
		return
	}

	// 調用回調函數
	if d.callback != nil {
		d.callback(PasswordDialogResult{
			Password:  password,
			Confirmed: true,
		})
	}

	// 關閉對話框
	d.dialog.Hide()
}

// handleCancel 處理取消按鈕點擊事件
// 調用回調函數並關閉對話框
func (d *PasswordSetupDialog) handleCancel() {
	// 調用回調函數
	if d.callback != nil {
		d.callback(PasswordDialogResult{
			Password:  "",
			Confirmed: false,
		})
	}

	// 關閉對話框
	d.dialog.Hide()
}

// PasswordVerifyDialog 密碼驗證對話框結構
// 用於驗證現有密碼的對話框
type PasswordVerifyDialog struct {
	dialog        dialog.Dialog        // Fyne 對話框實例
	passwordEntry *widget.Entry        // 密碼輸入框
	attemptsLabel *widget.Label        // 剩餘嘗試次數標籤
	callback      PasswordDialogCallback // 完成時的回調函數
	maxAttempts   int                  // 最大嘗試次數
	attempts      int                  // 當前嘗試次數
}

// NewPasswordVerifyDialog 創建新的密碼驗證對話框
// 參數：
//   - parent: 父視窗，用於模態顯示
//   - title: 對話框標題
//   - maxAttempts: 最大嘗試次數
//   - callback: 完成時的回調函數
// 回傳：密碼驗證對話框實例
//
// 執行流程：
// 1. 創建密碼輸入框
// 2. 創建嘗試次數標籤
// 3. 創建確認和取消按鈕
// 4. 組裝對話框佈局
func NewPasswordVerifyDialog(parent fyne.Window, title string, maxAttempts int, callback PasswordDialogCallback) *PasswordVerifyDialog {
	d := &PasswordVerifyDialog{
		callback:    callback,
		maxAttempts: maxAttempts,
		attempts:    0,
	}

	// 創建密碼輸入框
	d.passwordEntry = widget.NewPasswordEntry()
	d.passwordEntry.SetPlaceHolder("請輸入密碼")

	// 創建嘗試次數標籤
	d.attemptsLabel = widget.NewLabel(fmt.Sprintf("剩餘嘗試次數：%d", maxAttempts))

	// 設置 Enter 鍵提交
	d.passwordEntry.OnSubmitted = func(text string) {
		d.handleVerify()
	}

	// 創建確認按鈕
	verifyButton := widget.NewButton("驗證", func() {
		d.handleVerify()
	})

	// 創建取消按鈕
	cancelButton := widget.NewButton("取消", func() {
		d.handleCancel()
	})

	// 組裝對話框內容
	content := container.NewVBox(
		widget.NewLabel("請輸入密碼："),
		d.passwordEntry,
		d.attemptsLabel,
		widget.NewSeparator(),
		container.NewHBox(
			verifyButton,
			cancelButton,
		),
	)

	// 創建對話框
	d.dialog = dialog.NewCustom(title, "", content, parent)
	d.dialog.Resize(fyne.NewSize(350, 200))

	return d
}

// Show 顯示密碼驗證對話框
// 將對話框以模態方式顯示給用戶，並聚焦到密碼輸入框
func (d *PasswordVerifyDialog) Show() {
	d.dialog.Show()
	// 聚焦到密碼輸入框
	d.passwordEntry.FocusGained()
}

// handleVerify 處理密碼驗證
// 執行流程：
// 1. 獲取用戶輸入的密碼
// 2. 增加嘗試次數
// 3. 調用回調函數進行驗證
// 4. 根據驗證結果決定後續操作
func (d *PasswordVerifyDialog) handleVerify() {
	password := d.passwordEntry.Text
	d.attempts++

	// 檢查密碼是否為空
	if password == "" {
		dialog.ShowError(fmt.Errorf("密碼不能為空"), d.dialog.(fyne.Window))
		return
	}

	// 調用回調函數進行驗證
	if d.callback != nil {
		d.callback(PasswordDialogResult{
			Password:  password,
			Confirmed: true,
		})
	}

	// 更新嘗試次數顯示
	remaining := d.maxAttempts - d.attempts
	d.attemptsLabel.SetText(fmt.Sprintf("剩餘嘗試次數：%d", remaining))

	// 如果達到最大嘗試次數，關閉對話框
	if d.attempts >= d.maxAttempts {
		d.dialog.Hide()
	} else {
		// 清空密碼輸入框，準備下次輸入
		d.passwordEntry.SetText("")
	}
}

// handleCancel 處理取消按鈕點擊事件
// 調用回調函數並關閉對話框
func (d *PasswordVerifyDialog) handleCancel() {
	// 調用回調函數
	if d.callback != nil {
		d.callback(PasswordDialogResult{
			Password:  "",
			Confirmed: false,
		})
	}

	// 關閉對話框
	d.dialog.Hide()
}

// calculatePasswordStrength 計算密碼強度
// 參數：
//   - password: 要檢查的密碼
// 回傳：密碼強度等級
//
// 執行流程：
// 1. 檢查密碼長度
// 2. 檢查是否包含大寫字母
// 3. 檢查是否包含小寫字母
// 4. 檢查是否包含數字
// 5. 檢查是否包含特殊字符
// 6. 根據滿足的條件數量判斷強度等級
func calculatePasswordStrength(password string) PasswordStrength {
	if len(password) < 6 {
		return PasswordWeak
	}

	score := 0

	// 檢查密碼長度
	if len(password) >= 8 {
		score++
	}
	if len(password) >= 12 {
		score++
	}

	// 檢查字符類型
	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	// 根據字符類型增加分數
	if hasUpper {
		score++
	}
	if hasLower {
		score++
	}
	if hasDigit {
		score++
	}
	if hasSpecial {
		score++
	}

	// 根據分數判斷強度
	switch {
	case score >= 5:
		return PasswordStrong
	case score >= 3:
		return PasswordMedium
	default:
		return PasswordWeak
	}
}

// ShowPasswordSetupDialog 顯示密碼設定對話框的便利函數
// 參數：
//   - parent: 父視窗
//   - title: 對話框標題
//   - callback: 完成時的回調函數
//
// 執行流程：
// 1. 創建密碼設定對話框
// 2. 顯示對話框
func ShowPasswordSetupDialog(parent fyne.Window, title string, callback PasswordDialogCallback) {
	dialog := NewPasswordSetupDialog(parent, title, callback)
	dialog.Show()
}

// ShowPasswordVerifyDialog 顯示密碼驗證對話框的便利函數
// 參數：
//   - parent: 父視窗
//   - title: 對話框標題
//   - maxAttempts: 最大嘗試次數
//   - callback: 完成時的回調函數
//
// 執行流程：
// 1. 創建密碼驗證對話框
// 2. 顯示對話框
func ShowPasswordVerifyDialog(parent fyne.Window, title string, maxAttempts int, callback PasswordDialogCallback) {
	dialog := NewPasswordVerifyDialog(parent, title, maxAttempts, callback)
	dialog.Show()
}