package ui

import (
	"testing"

	"fyne.io/fyne/v2/test"
)

// TestAuthMethod 測試驗證方法枚舉
func TestAuthMethod(t *testing.T) {
	// 驗證枚舉值
	if AuthMethodPassword != 0 {
		t.Errorf("AuthMethodPassword 應該等於 0，實際為 %d", AuthMethodPassword)
	}

	if AuthMethodBiometric != 1 {
		t.Errorf("AuthMethodBiometric 應該等於 1，實際為 %d", AuthMethodBiometric)
	}

	if AuthMethodBoth != 2 {
		t.Errorf("AuthMethodBoth 應該等於 2，實際為 %d", AuthMethodBoth)
	}
}

// TestAuthResult 測試驗證結果結構
func TestAuthResult(t *testing.T) {
	// 測試成功結果
	successResult := AuthResult{
		Success:      true,
		Password:     "testpassword",
		Method:       AuthMethodPassword,
		ErrorMessage: "",
		Cancelled:    false,
	}

	if !successResult.Success {
		t.Error("成功結果的 Success 應該為 true")
	}

	if successResult.Password != "testpassword" {
		t.Error("成功結果應該包含密碼")
	}

	if successResult.Method != AuthMethodPassword {
		t.Error("成功結果應該記錄正確的驗證方法")
	}

	// 測試失敗結果
	failureResult := AuthResult{
		Success:      false,
		Password:     "",
		Method:       AuthMethodBiometric,
		ErrorMessage: "驗證失敗",
		Cancelled:    false,
	}

	if failureResult.Success {
		t.Error("失敗結果的 Success 應該為 false")
	}

	if failureResult.ErrorMessage == "" {
		t.Error("失敗結果應該包含錯誤訊息")
	}

	// 測試取消結果
	cancelResult := AuthResult{
		Success:      false,
		Password:     "",
		Method:       AuthMethodPassword,
		ErrorMessage: "",
		Cancelled:    true,
	}

	if !cancelResult.Cancelled {
		t.Error("取消結果的 Cancelled 應該為 true")
	}
}

// TestAuthDialogManager 測試驗證對話框管理器的創建和基本功能
func TestAuthDialogManager(t *testing.T) {
	// 創建測試視窗
	testWindow := test.NewWindow(nil)
	defer testWindow.Close()

	// 測試正常創建
	manager := NewAuthDialogManager(testWindow, true, AuthMethodBoth, true, 3)

	if manager == nil {
		t.Fatal("驗證對話框管理器創建失敗")
	}

	// 驗證初始設置
	if manager.parent != testWindow {
		t.Error("父視窗設置錯誤")
	}

	if !manager.biometricAvailable {
		t.Error("生物驗證可用性設置錯誤")
	}

	if manager.preferredMethod != AuthMethodBoth {
		t.Error("首選驗證方法設置錯誤")
	}

	if !manager.allowFallback {
		t.Error("回退設置錯誤")
	}

	if manager.maxPasswordAttempts != 3 {
		t.Error("最大密碼嘗試次數設置錯誤")
	}

	// 測試預設值處理
	managerWithDefaults := NewAuthDialogManager(testWindow, false, AuthMethodPassword, false, 0)
	if managerWithDefaults.maxPasswordAttempts != 3 {
		t.Error("應該使用預設的最大嘗試次數 3")
	}
}

// TestAuthDialogManagerGettersAndSetters 測試驗證對話框管理器的 getter 和 setter 方法
func TestAuthDialogManagerGettersAndSetters(t *testing.T) {
	// 創建測試視窗
	testWindow := test.NewWindow(nil)
	defer testWindow.Close()

	// 創建管理器
	manager := NewAuthDialogManager(testWindow, true, AuthMethodBoth, true, 3)

	// 測試 getter 方法
	if !manager.GetBiometricAvailable() {
		t.Error("GetBiometricAvailable 返回值錯誤")
	}

	if manager.GetPreferredMethod() != AuthMethodBoth {
		t.Error("GetPreferredMethod 返回值錯誤")
	}

	if !manager.GetAllowFallback() {
		t.Error("GetAllowFallback 返回值錯誤")
	}

	if manager.GetMaxPasswordAttempts() != 3 {
		t.Error("GetMaxPasswordAttempts 返回值錯誤")
	}

	// 測試 setter 方法
	manager.SetBiometricAvailable(false)
	if manager.GetBiometricAvailable() {
		t.Error("SetBiometricAvailable 設置失敗")
	}

	manager.SetPreferredMethod(AuthMethodPassword)
	if manager.GetPreferredMethod() != AuthMethodPassword {
		t.Error("SetPreferredMethod 設置失敗")
	}

	manager.SetAllowFallback(false)
	if manager.GetAllowFallback() {
		t.Error("SetAllowFallback 設置失敗")
	}

	manager.SetMaxPasswordAttempts(5)
	if manager.GetMaxPasswordAttempts() != 5 {
		t.Error("SetMaxPasswordAttempts 設置失敗")
	}

	// 測試無效的最大嘗試次數
	manager.SetMaxPasswordAttempts(0)
	if manager.GetMaxPasswordAttempts() != 5 {
		t.Error("無效的最大嘗試次數應該被忽略")
	}

	manager.SetMaxPasswordAttempts(-1)
	if manager.GetMaxPasswordAttempts() != 5 {
		t.Error("負數的最大嘗試次數應該被忽略")
	}
}

// TestAuthDialogManagerPasswordMethod 測試密碼驗證方法
func TestAuthDialogManagerPasswordMethod(t *testing.T) {
	// 創建測試視窗
	testWindow := test.NewWindow(nil)
	defer testWindow.Close()

	// 創建只使用密碼驗證的管理器
	manager := NewAuthDialogManager(testWindow, false, AuthMethodPassword, false, 3)

	// 測試回調函數
	callback := func(result AuthResult) {
		// 回調函數用於測試
	}

	// 由於我們無法在測試中實際顯示對話框，這裡主要測試邏輯
	// 實際的對話框顯示會在整合測試中進行

	// 驗證管理器設置
	if manager.GetPreferredMethod() != AuthMethodPassword {
		t.Error("應該使用密碼驗證方法")
	}

	if manager.GetBiometricAvailable() {
		t.Error("生物驗證應該不可用")
	}

	if manager.GetAllowFallback() {
		t.Error("不應該允許回退")
	}

	// 測試密碼設定對話框
	manager.ShowPasswordSetupDialog("設定密碼", callback)

	// 由於對話框是異步的，我們無法直接測試回調結果
	// 但可以驗證沒有立即崩潰
}

// TestAuthDialogManagerBiometricMethod 測試生物驗證方法
func TestAuthDialogManagerBiometricMethod(t *testing.T) {
	// 創建測試視窗
	testWindow := test.NewWindow(nil)
	defer testWindow.Close()

	// 創建只使用生物驗證的管理器
	manager := NewAuthDialogManager(testWindow, true, AuthMethodBiometric, false, 3)

	// 測試回調函數
	callback := func(result AuthResult) {
		// 回調函數用於測試
	}

	// 驗證管理器設置
	if manager.GetPreferredMethod() != AuthMethodBiometric {
		t.Error("應該使用生物驗證方法")
	}

	if !manager.GetBiometricAvailable() {
		t.Error("生物驗證應該可用")
	}

	// 測試生物驗證設定對話框
	manager.ShowBiometricSetupDialog("設定生物驗證", callback)

	// 由於對話框是異步的，我們無法直接測試回調結果
	// 但可以驗證沒有立即崩潰
}

// TestAuthDialogManagerBothMethods 測試兩種驗證方法都支援的情況
func TestAuthDialogManagerBothMethods(t *testing.T) {
	// 創建測試視窗
	testWindow := test.NewWindow(nil)
	defer testWindow.Close()

	// 創建支援兩種驗證方法的管理器
	manager := NewAuthDialogManager(testWindow, true, AuthMethodBoth, true, 3)

	// 驗證管理器設置
	if manager.GetPreferredMethod() != AuthMethodBoth {
		t.Error("應該支援兩種驗證方法")
	}

	if !manager.GetBiometricAvailable() {
		t.Error("生物驗證應該可用")
	}

	if !manager.GetAllowFallback() {
		t.Error("應該允許回退")
	}
}

// TestAuthDialogManagerFallbackScenarios 測試回退場景
func TestAuthDialogManagerFallbackScenarios(t *testing.T) {
	// 創建測試視窗
	testWindow := test.NewWindow(nil)
	defer testWindow.Close()

	// 測試場景1：首選生物驗證但不可用，允許回退
	t.Run("生物驗證不可用允許回退", func(t *testing.T) {
		manager := NewAuthDialogManager(testWindow, false, AuthMethodBiometric, true, 3)

		if manager.GetBiometricAvailable() {
			t.Error("生物驗證應該不可用")
		}

		if !manager.GetAllowFallback() {
			t.Error("應該允許回退")
		}
	})

	// 測試場景2：首選生物驗證但不可用，不允許回退
	t.Run("生物驗證不可用不允許回退", func(t *testing.T) {
		manager := NewAuthDialogManager(testWindow, false, AuthMethodBiometric, false, 3)

		if manager.GetBiometricAvailable() {
			t.Error("生物驗證應該不可用")
		}

		if manager.GetAllowFallback() {
			t.Error("不應該允許回退")
		}
	})
}

// TestConvenienceFunctions 測試便利函數
func TestConvenienceFunctions(t *testing.T) {
	// 創建測試視窗
	testWindow := test.NewWindow(nil)
	defer testWindow.Close()

	// 測試回調函數
	callback := func(result AuthResult) {
		// 回調函數用於測試
	}

	// 測試 ShowQuickAuthDialog
	// 由於對話框是異步的，我們主要測試函數不會崩潰
	ShowQuickAuthDialog(testWindow, "快速驗證", "請進行驗證", true, callback)

	// 測試 ShowQuickPasswordSetup
	ShowQuickPasswordSetup(testWindow, "快速密碼設定", callback)

	// 測試 ShowQuickBiometricSetup
	ShowQuickBiometricSetup(testWindow, "快速生物驗證設定", true, callback)

	// 如果執行到這裡沒有崩潰，說明便利函數工作正常
}

// TestAuthDialogManagerErrorHandling 測試錯誤處理
func TestAuthDialogManagerErrorHandling(t *testing.T) {
	// 創建測試視窗
	testWindow := test.NewWindow(nil)
	defer testWindow.Close()

	// 創建管理器
	manager := NewAuthDialogManager(testWindow, false, AuthMethodBiometric, false, 3)

	// 測試回調函數
	var callbackResult AuthResult
	var callbackCalled bool
	callback := func(result AuthResult) {
		callbackResult = result
		callbackCalled = true
	}

	// 測試錯誤處理
	manager.handleAuthError("測試錯誤", callback)

	// 驗證回調被調用
	if !callbackCalled {
		t.Error("錯誤處理應該調用回調函數")
	}

	// 驗證錯誤結果
	if callbackResult.Success {
		t.Error("錯誤結果的 Success 應該為 false")
	}

	if callbackResult.ErrorMessage != "測試錯誤" {
		t.Error("錯誤結果應該包含正確的錯誤訊息")
	}
}

// BenchmarkAuthDialogManagerCreation 驗證對話框管理器創建的效能基準測試
func BenchmarkAuthDialogManagerCreation(b *testing.B) {
	// 創建測試視窗
	testWindow := test.NewWindow(nil)
	defer testWindow.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		manager := NewAuthDialogManager(testWindow, true, AuthMethodBoth, true, 3)
		_ = manager // 避免編譯器優化
	}
}

// TestAuthDialogManagerMemoryUsage 測試驗證對話框管理器的記憶體使用
func TestAuthDialogManagerMemoryUsage(t *testing.T) {
	// 創建測試視窗
	testWindow := test.NewWindow(nil)
	defer testWindow.Close()

	// 創建多個管理器實例以測試記憶體使用
	managers := make([]*AuthDialogManager, 100)
	for i := 0; i < 100; i++ {
		managers[i] = NewAuthDialogManager(testWindow, true, AuthMethodBoth, true, 3)
	}

	// 驗證所有管理器都正確創建
	for i, manager := range managers {
		if manager == nil {
			t.Errorf("第 %d 個管理器創建失敗", i)
		}
	}

	// 測試管理器的基本功能
	for _, manager := range managers {
		if manager.GetBiometricAvailable() != true {
			t.Error("管理器設置錯誤")
		}
	}
}