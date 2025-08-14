package ui

import (
	"testing"
	"time"

	"fyne.io/fyne/v2/test"
)

// TestBiometricAuthStatus 測試生物驗證狀態枚舉
func TestBiometricAuthStatus(t *testing.T) {
	// 驗證狀態枚舉值
	if BiometricStatusIdle != 0 {
		t.Errorf("BiometricStatusIdle 應該等於 0，實際為 %d", BiometricStatusIdle)
	}

	if BiometricStatusWaiting != 1 {
		t.Errorf("BiometricStatusWaiting 應該等於 1，實際為 %d", BiometricStatusWaiting)
	}

	if BiometricStatusSuccess != 2 {
		t.Errorf("BiometricStatusSuccess 應該等於 2，實際為 %d", BiometricStatusSuccess)
	}

	if BiometricStatusFailed != 3 {
		t.Errorf("BiometricStatusFailed 應該等於 3，實際為 %d", BiometricStatusFailed)
	}

	if BiometricStatusUnavailable != 4 {
		t.Errorf("BiometricStatusUnavailable 應該等於 4，實際為 %d", BiometricStatusUnavailable)
	}
}

// TestBiometricAuthResult 測試生物驗證結果結構
func TestBiometricAuthResult(t *testing.T) {
	// 測試成功結果
	successResult := BiometricAuthResult{
		Success:       true,
		ErrorMessage:  "",
		FallbackUsed:  false,
		UserCancelled: false,
	}

	if !successResult.Success {
		t.Error("成功結果的 Success 應該為 true")
	}

	if successResult.ErrorMessage != "" {
		t.Error("成功結果的 ErrorMessage 應該為空")
	}

	// 測試失敗結果
	failureResult := BiometricAuthResult{
		Success:       false,
		ErrorMessage:  "驗證失敗",
		FallbackUsed:  false,
		UserCancelled: false,
	}

	if failureResult.Success {
		t.Error("失敗結果的 Success 應該為 false")
	}

	if failureResult.ErrorMessage == "" {
		t.Error("失敗結果應該包含錯誤訊息")
	}

	// 測試回退結果
	fallbackResult := BiometricAuthResult{
		Success:       false,
		ErrorMessage:  "",
		FallbackUsed:  true,
		UserCancelled: false,
	}

	if !fallbackResult.FallbackUsed {
		t.Error("回退結果的 FallbackUsed 應該為 true")
	}

	// 測試取消結果
	cancelResult := BiometricAuthResult{
		Success:       false,
		ErrorMessage:  "用戶取消了驗證",
		FallbackUsed:  false,
		UserCancelled: true,
	}

	if !cancelResult.UserCancelled {
		t.Error("取消結果的 UserCancelled 應該為 true")
	}
}

// TestBiometricAuthDialog 測試生物驗證對話框的創建和基本功能
func TestBiometricAuthDialog(t *testing.T) {
	// 創建測試視窗
	testWindow := test.NewWindow(nil)
	defer testWindow.Close()

	// 測試回調函數
	var callbackResult BiometricAuthResult
	var callbackCalled bool
	callback := func(result BiometricAuthResult) {
		callbackResult = result
		callbackCalled = true
	}

	// 創建生物驗證對話框
	dialog := NewBiometricAuthDialog(testWindow, "測試驗證", "請進行生物驗證", true, callback)

	// 驗證對話框是否正確創建
	if dialog == nil {
		t.Fatal("生物驗證對話框創建失敗")
	}

	// 驗證對話框組件是否正確初始化
	if dialog.statusLabel == nil {
		t.Error("狀態標籤未正確初始化")
	}

	if dialog.statusIcon == nil {
		t.Error("狀態圖示未正確初始化")
	}

	if dialog.progressBar == nil {
		t.Error("進度指示器未正確初始化")
	}

	if dialog.fallbackButton == nil {
		t.Error("回退按鈕未正確初始化")
	}

	if dialog.cancelButton == nil {
		t.Error("取消按鈕未正確初始化")
	}

	// 驗證初始狀態
	if dialog.status != BiometricStatusIdle {
		t.Errorf("初始狀態應該為 BiometricStatusIdle，實際為 %d", dialog.status)
	}

	if !dialog.fallbackEnabled {
		t.Error("回退功能應該已啟用")
	}

	// 測試取消操作
	callbackCalled = false
	dialog.handleCancel()

	if !callbackCalled {
		t.Error("取消操作應該觸發回調函數")
	}

	if !callbackResult.UserCancelled {
		t.Error("取消操作的結果應該標記為用戶取消")
	}

	// 測試回退操作
	callbackCalled = false
	dialog.handleFallback()

	if !callbackCalled {
		t.Error("回退操作應該觸發回調函數")
	}

	if !callbackResult.FallbackUsed {
		t.Error("回退操作的結果應該標記為使用回退")
	}
}

// TestBiometricAuthDialogWithoutFallback 測試不啟用回退的生物驗證對話框
func TestBiometricAuthDialogWithoutFallback(t *testing.T) {
	// 創建測試視窗
	testWindow := test.NewWindow(nil)
	defer testWindow.Close()

	// 測試回調函數
	callback := func(result BiometricAuthResult) {
		// 回調函數用於測試
	}

	// 創建不啟用回退的生物驗證對話框
	dialog := NewBiometricAuthDialog(testWindow, "測試驗證", "請進行生物驗證", false, callback)

	// 驗證回退功能被禁用
	if dialog.fallbackEnabled {
		t.Error("回退功能應該被禁用")
	}

	// 驗證回退按鈕被隱藏
	if dialog.fallbackButton.Visible() {
		t.Error("回退按鈕應該被隱藏")
	}
}

// TestBiometricAuthDialogStatusChanges 測試生物驗證對話框的狀態變更
func TestBiometricAuthDialogStatusChanges(t *testing.T) {
	// 創建測試視窗
	testWindow := test.NewWindow(nil)
	defer testWindow.Close()

	// 測試回調函數
	callback := func(result BiometricAuthResult) {
		// 回調函數用於測試
	}

	// 創建生物驗證對話框
	dialog := NewBiometricAuthDialog(testWindow, "測試驗證", "請進行生物驗證", true, callback)

	// 測試設置等待狀態
	dialog.SetStatus(BiometricStatusWaiting, "正在等待驗證...")
	if dialog.status != BiometricStatusWaiting {
		t.Error("狀態應該更新為 BiometricStatusWaiting")
	}

	if dialog.statusLabel.Text != "正在等待驗證..." {
		t.Error("狀態標籤文字應該更新")
	}

	if !dialog.progressBar.Visible() {
		t.Error("等待狀態時進度指示器應該可見")
	}

	// 測試設置成功狀態
	dialog.SetStatus(BiometricStatusSuccess, "驗證成功！")
	if dialog.status != BiometricStatusSuccess {
		t.Error("狀態應該更新為 BiometricStatusSuccess")
	}

	if dialog.progressBar.Visible() {
		t.Error("成功狀態時進度指示器應該隱藏")
	}

	// 測試設置失敗狀態
	dialog.SetStatus(BiometricStatusFailed, "驗證失敗")
	if dialog.status != BiometricStatusFailed {
		t.Error("狀態應該更新為 BiometricStatusFailed")
	}

	// 測試設置不可用狀態
	dialog.SetStatus(BiometricStatusUnavailable, "生物驗證不可用")
	if dialog.status != BiometricStatusUnavailable {
		t.Error("狀態應該更新為 BiometricStatusUnavailable")
	}
}

// TestBiometricAuthDialogNotifications 測試生物驗證對話框的通知方法
func TestBiometricAuthDialogNotifications(t *testing.T) {
	// 創建測試視窗
	testWindow := test.NewWindow(nil)
	defer testWindow.Close()

	// 測試回調函數
	var callbackResult BiometricAuthResult
	var callbackCalled bool
	callback := func(result BiometricAuthResult) {
		callbackResult = result
		callbackCalled = true
	}

	// 創建生物驗證對話框
	dialog := NewBiometricAuthDialog(testWindow, "測試驗證", "請進行生物驗證", true, callback)

	// 測試成功通知
	callbackCalled = false
	dialog.NotifySuccess()

	if !callbackCalled {
		t.Error("成功通知應該觸發回調函數")
	}

	if !callbackResult.Success {
		t.Error("成功通知的結果應該標記為成功")
	}

	if dialog.status != BiometricStatusSuccess {
		t.Error("成功通知應該更新狀態為 BiometricStatusSuccess")
	}

	// 測試失敗通知
	callbackCalled = false
	dialog.NotifyFailure("測試錯誤")

	if !callbackCalled {
		t.Error("失敗通知應該觸發回調函數")
	}

	if callbackResult.Success {
		t.Error("失敗通知的結果應該標記為失敗")
	}

	if callbackResult.ErrorMessage != "測試錯誤" {
		t.Error("失敗通知應該包含錯誤訊息")
	}

	if dialog.status != BiometricStatusFailed {
		t.Error("失敗通知應該更新狀態為 BiometricStatusFailed")
	}

	// 測試不可用通知
	dialog.NotifyUnavailable("設備不支援")

	if dialog.status != BiometricStatusUnavailable {
		t.Error("不可用通知應該更新狀態為 BiometricStatusUnavailable")
	}
}

// TestBiometricSetupDialog 測試生物驗證設置對話框的創建和基本功能
func TestBiometricSetupDialog(t *testing.T) {
	// 創建測試視窗
	testWindow := test.NewWindow(nil)
	defer testWindow.Close()

	// 測試回調函數
	callback := func(result BiometricAuthResult) {
		// 回調函數用於測試
	}

	// 測試生物驗證可用的情況
	t.Run("生物驗證可用", func(t *testing.T) {
		dialog := NewBiometricSetupDialog(testWindow, "設置生物驗證", true, callback)

		// 驗證對話框是否正確創建
		if dialog == nil {
			t.Fatal("生物驗證設置對話框創建失敗")
		}

		// 驗證對話框組件是否正確初始化
		if dialog.enableCheckbox == nil {
			t.Error("啟用複選框未正確初始化")
		}

		if dialog.fallbackCheckbox == nil {
			t.Error("回退複選框未正確初始化")
		}

		if dialog.statusLabel == nil {
			t.Error("狀態標籤未正確初始化")
		}

		if dialog.testButton == nil {
			t.Error("測試按鈕未正確初始化")
		}

		// 驗證生物驗證可用性
		if !dialog.biometricAvailable {
			t.Error("生物驗證應該標記為可用")
		}

		// 驗證啟用複選框是否可用
		if dialog.enableCheckbox.Disabled() {
			t.Error("生物驗證可用時，啟用複選框應該可用")
		}
	})

	// 測試生物驗證不可用的情況
	t.Run("生物驗證不可用", func(t *testing.T) {
		dialog := NewBiometricSetupDialog(testWindow, "設置生物驗證", false, callback)

		// 驗證生物驗證不可用性
		if dialog.biometricAvailable {
			t.Error("生物驗證應該標記為不可用")
		}

		// 驗證啟用複選框是否被禁用
		if !dialog.enableCheckbox.Disabled() {
			t.Error("生物驗證不可用時，啟用複選框應該被禁用")
		}

		// 驗證測試按鈕是否被禁用
		if !dialog.testButton.Disabled() {
			t.Error("生物驗證不可用時，測試按鈕應該被禁用")
		}
	})
}

// TestBiometricSetupDialogInteractions 測試生物驗證設置對話框的互動
func TestBiometricSetupDialogInteractions(t *testing.T) {
	// 創建測試視窗
	testWindow := test.NewWindow(nil)
	defer testWindow.Close()

	// 測試回調函數
	var callbackResult BiometricAuthResult
	var callbackCalled bool
	callback := func(result BiometricAuthResult) {
		callbackResult = result
		callbackCalled = true
	}

	// 創建生物驗證設置對話框
	dialog := NewBiometricSetupDialog(testWindow, "設置生物驗證", true, callback)

	// 測試啟用生物驗證
	dialog.enableCheckbox.SetChecked(true)

	// 驗證測試按鈕是否啟用
	if dialog.testButton.Disabled() {
		t.Error("啟用生物驗證後，測試按鈕應該可用")
	}

	// 驗證回退複選框是否啟用
	if dialog.fallbackCheckbox.Disabled() {
		t.Error("啟用生物驗證後，回退複選框應該可用")
	}

	// 測試禁用生物驗證
	dialog.enableCheckbox.SetChecked(false)

	// 驗證測試按鈕是否禁用
	if !dialog.testButton.Disabled() {
		t.Error("禁用生物驗證後，測試按鈕應該被禁用")
	}

	// 測試保存操作
	callbackCalled = false
	dialog.handleSave()

	if !callbackCalled {
		t.Error("保存操作應該觸發回調函數")
	}

	if !callbackResult.Success {
		t.Error("保存操作的結果應該標記為成功")
	}

	// 測試取消操作
	callbackCalled = false
	dialog.handleCancel()

	if !callbackCalled {
		t.Error("取消操作應該觸發回調函數")
	}

	if !callbackResult.UserCancelled {
		t.Error("取消操作的結果應該標記為用戶取消")
	}
}

// TestBiometricDialogConvenienceFunctions 測試便利函數
func TestBiometricDialogConvenienceFunctions(t *testing.T) {
	// 創建測試視窗
	testWindow := test.NewWindow(nil)
	defer testWindow.Close()

	// 測試回調函數
	callback := func(result BiometricAuthResult) {
		// 回調函數用於測試
	}

	// 測試 ShowBiometricAuthDialog 便利函數
	authDialog := ShowBiometricAuthDialog(testWindow, "測試驗證", "請進行驗證", true, callback)

	if authDialog == nil {
		t.Error("ShowBiometricAuthDialog 應該返回對話框實例")
	}

	if authDialog.status != BiometricStatusWaiting {
		t.Error("便利函數創建的對話框應該處於等待狀態")
	}

	// 測試 ShowBiometricSetupDialog 便利函數
	ShowBiometricSetupDialog(testWindow, "設置驗證", true, callback)
	// 便利函數不返回實例，只要不崩潰就算成功
}

// TestBiometricDialogProgressAnimation 測試進度動畫
func TestBiometricDialogProgressAnimation(t *testing.T) {
	// 創建測試視窗
	testWindow := test.NewWindow(nil)
	defer testWindow.Close()

	// 測試回調函數
	callback := func(result BiometricAuthResult) {
		// 回調函數用於測試
	}

	// 創建生物驗證對話框
	dialog := NewBiometricAuthDialog(testWindow, "測試驗證", "請進行生物驗證", true, callback)

	// 設置等待狀態以啟動進度動畫
	dialog.SetStatus(BiometricStatusWaiting, "正在等待驗證...")

	// 等待一小段時間讓動畫開始
	time.Sleep(200 * time.Millisecond)

	// 驗證進度條是否可見
	if !dialog.progressBar.Visible() {
		t.Error("等待狀態時進度條應該可見")
	}

	// 停止動畫
	dialog.SetStatus(BiometricStatusSuccess, "驗證成功")

	// 等待一小段時間讓動畫停止
	time.Sleep(100 * time.Millisecond)

	// 驗證進度條是否隱藏
	if dialog.progressBar.Visible() {
		t.Error("成功狀態時進度條應該隱藏")
	}
}

// BenchmarkBiometricDialogCreation 生物驗證對話框創建的效能基準測試
func BenchmarkBiometricDialogCreation(b *testing.B) {
	// 創建測試視窗
	testWindow := test.NewWindow(nil)
	defer testWindow.Close()

	// 測試回調函數
	callback := func(result BiometricAuthResult) {
		// 回調函數用於測試
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dialog := NewBiometricAuthDialog(testWindow, "基準測試", "測試訊息", true, callback)
		dialog.Hide()
	}
}

// TestBiometricDialogMemoryUsage 測試生物驗證對話框的記憶體使用
func TestBiometricDialogMemoryUsage(t *testing.T) {
	// 創建測試視窗
	testWindow := test.NewWindow(nil)
	defer testWindow.Close()

	// 創建多個對話框實例以測試記憶體使用
	dialogs := make([]*BiometricAuthDialog, 10)
	for i := 0; i < 10; i++ {
		dialogs[i] = NewBiometricAuthDialog(testWindow, "記憶體測試", "測試訊息", true, nil)
	}

	// 驗證所有對話框都正確創建
	for i, dialog := range dialogs {
		if dialog == nil {
			t.Errorf("第 %d 個對話框創建失敗", i)
		}
	}

	// 清理資源
	for _, dialog := range dialogs {
		if dialog.dialog != nil {
			dialog.Hide()
		}
	}
}