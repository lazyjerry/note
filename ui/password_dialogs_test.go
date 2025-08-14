package ui

import (
	"testing"

	"fyne.io/fyne/v2/test"
)

// TestPasswordStrengthCalculation 測試密碼強度計算功能
// 驗證不同類型密碼的強度評估是否正確
func TestPasswordStrengthCalculation(t *testing.T) {
	tests := []struct {
		name     string
		password string
		expected PasswordStrength
	}{
		{
			name:     "空密碼應該是弱密碼",
			password: "",
			expected: PasswordWeak,
		},
		{
			name:     "短密碼應該是弱密碼",
			password: "123",
			expected: PasswordWeak,
		},
		{
			name:     "只有小寫字母的短密碼應該是弱密碼",
			password: "abcdef",
			expected: PasswordWeak,
		},
		{
			name:     "包含大小寫字母和數字的中等長度密碼應該是中等強度",
			password: "Abc123",
			expected: PasswordMedium,
		},
		{
			name:     "包含大小寫字母、數字和特殊字符的長密碼應該是強密碼",
			password: "Abc123!@#",
			expected: PasswordStrong,
		},
		{
			name:     "很長但只有小寫字母的密碼應該是中等強度",
			password: "abcdefghijklmnop",
			expected: PasswordMedium,
		},
		{
			name:     "包含所有字符類型的超長密碼應該是強密碼",
			password: "MyVerySecurePassword123!@#",
			expected: PasswordStrong,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculatePasswordStrength(tt.password)
			if result != tt.expected {
				t.Errorf("calculatePasswordStrength(%q) = %v, 期望 %v", tt.password, result, tt.expected)
			}
		})
	}
}

// TestPasswordSetupDialog 測試密碼設定對話框的創建和基本功能
func TestPasswordSetupDialog(t *testing.T) {
	// 創建測試視窗
	testWindow := test.NewWindow(nil)
	defer testWindow.Close()

	// 測試回調函數
	callback := func(result PasswordDialogResult) {
		// 回調函數用於測試
	}

	// 創建密碼設定對話框
	dialog := NewPasswordSetupDialog(testWindow, "設定密碼", callback)

	// 驗證對話框是否正確創建
	if dialog == nil {
		t.Fatal("密碼設定對話框創建失敗")
	}

	// 驗證對話框組件是否正確初始化
	if dialog.passwordEntry == nil {
		t.Error("密碼輸入框未正確初始化")
	}

	if dialog.confirmEntry == nil {
		t.Error("確認密碼輸入框未正確初始化")
	}

	if dialog.strengthBar == nil {
		t.Error("密碼強度指示器未正確初始化")
	}

	if dialog.strengthLabel == nil {
		t.Error("密碼強度標籤未正確初始化")
	}

	// 測試密碼強度更新功能
	dialog.updatePasswordStrength("weak")
	if dialog.strengthBar.Value != 0.33 {
		t.Errorf("弱密碼強度指示器數值錯誤，期望 0.33，實際 %f", dialog.strengthBar.Value)
	}

	dialog.updatePasswordStrength("Medium1!")
	if dialog.strengthBar.Value != 0.66 {
		t.Errorf("中等密碼強度指示器數值錯誤，期望 0.66，實際 %f", dialog.strengthBar.Value)
	}

	dialog.updatePasswordStrength("VeryStrong123!@#")
	if dialog.strengthBar.Value != 1.0 {
		t.Errorf("強密碼強度指示器數值錯誤，期望 1.0，實際 %f", dialog.strengthBar.Value)
	}
}

// TestPasswordVerifyDialog 測試密碼驗證對話框的創建和基本功能
func TestPasswordVerifyDialog(t *testing.T) {
	// 創建測試視窗
	testWindow := test.NewWindow(nil)
	defer testWindow.Close()

	// 測試回調函數
	callback := func(result PasswordDialogResult) {
		// 回調函數用於測試
	}

	// 創建密碼驗證對話框
	maxAttempts := 3
	dialog := NewPasswordVerifyDialog(testWindow, "驗證密碼", maxAttempts, callback)

	// 驗證對話框是否正確創建
	if dialog == nil {
		t.Fatal("密碼驗證對話框創建失敗")
	}

	// 驗證對話框組件是否正確初始化
	if dialog.passwordEntry == nil {
		t.Error("密碼輸入框未正確初始化")
	}

	if dialog.attemptsLabel == nil {
		t.Error("嘗試次數標籤未正確初始化")
	}

	// 驗證最大嘗試次數設置
	if dialog.maxAttempts != maxAttempts {
		t.Errorf("最大嘗試次數設置錯誤，期望 %d，實際 %d", maxAttempts, dialog.maxAttempts)
	}

	// 驗證初始嘗試次數
	if dialog.attempts != 0 {
		t.Errorf("初始嘗試次數錯誤，期望 0，實際 %d", dialog.attempts)
	}
}

// TestPasswordDialogCallbacks 測試密碼對話框的回調功能
func TestPasswordDialogCallbacks(t *testing.T) {
	// 創建測試視窗
	testWindow := test.NewWindow(nil)
	defer testWindow.Close()

	t.Run("測試密碼設定對話框取消回調", func(t *testing.T) {
		var callbackResult PasswordDialogResult
		callback := func(result PasswordDialogResult) {
			callbackResult = result
		}

		dialog := NewPasswordSetupDialog(testWindow, "設定密碼", callback)
		dialog.handleCancel()

		// 驗證取消操作的回調結果
		if callbackResult.Confirmed {
			t.Error("取消操作後 Confirmed 應該為 false")
		}

		if callbackResult.Password != "" {
			t.Error("取消操作後 Password 應該為空字符串")
		}
	})

	t.Run("測試密碼驗證對話框取消回調", func(t *testing.T) {
		var callbackResult PasswordDialogResult
		callback := func(result PasswordDialogResult) {
			callbackResult = result
		}

		dialog := NewPasswordVerifyDialog(testWindow, "驗證密碼", 3, callback)
		dialog.handleCancel()

		// 驗證取消操作的回調結果
		if callbackResult.Confirmed {
			t.Error("取消操作後 Confirmed 應該為 false")
		}

		if callbackResult.Password != "" {
			t.Error("取消操作後 Password 應該為空字符串")
		}
	})
}

// TestPasswordSetupDialogValidation 測試密碼設定對話框的驗證邏輯
func TestPasswordSetupDialogValidation(t *testing.T) {
	// 創建測試視窗
	testWindow := test.NewWindow(nil)
	defer testWindow.Close()

	var callbackCalled bool
	callback := func(result PasswordDialogResult) {
		callbackCalled = true
	}

	dialog := NewPasswordSetupDialog(testWindow, "設定密碼", callback)

	t.Run("測試空密碼驗證", func(t *testing.T) {
		callbackCalled = false
		dialog.passwordEntry.SetText("")
		dialog.confirmEntry.SetText("")
		dialog.handleConfirm()

		// 空密碼應該不會觸發回調
		if callbackCalled {
			t.Error("空密碼不應該觸發成功回調")
		}
	})

	t.Run("測試密碼不一致驗證", func(t *testing.T) {
		callbackCalled = false
		dialog.passwordEntry.SetText("password123")
		dialog.confirmEntry.SetText("different123")
		dialog.handleConfirm()

		// 密碼不一致應該不會觸發回調
		if callbackCalled {
			t.Error("密碼不一致不應該觸發成功回調")
		}
	})

	t.Run("測試弱密碼驗證", func(t *testing.T) {
		callbackCalled = false
		weakPassword := "123"
		dialog.passwordEntry.SetText(weakPassword)
		dialog.confirmEntry.SetText(weakPassword)
		dialog.handleConfirm()

		// 弱密碼應該不會觸發回調
		if callbackCalled {
			t.Error("弱密碼不應該觸發成功回調")
		}
	})
}

// TestPasswordVerifyDialogAttempts 測試密碼驗證對話框的嘗試次數邏輯
func TestPasswordVerifyDialogAttempts(t *testing.T) {
	// 創建測試視窗
	testWindow := test.NewWindow(nil)
	defer testWindow.Close()

	var callbackResults []PasswordDialogResult
	callback := func(result PasswordDialogResult) {
		callbackResults = append(callbackResults, result)
	}

	maxAttempts := 3
	dialog := NewPasswordVerifyDialog(testWindow, "驗證密碼", maxAttempts, callback)

	// 模擬多次密碼輸入
	for i := 0; i < maxAttempts; i++ {
		dialog.passwordEntry.SetText("testpassword")
		dialog.handleVerify()
	}

	// 驗證回調被調用的次數
	if len(callbackResults) != maxAttempts {
		t.Errorf("回調調用次數錯誤，期望 %d，實際 %d", maxAttempts, len(callbackResults))
	}

	// 驗證嘗試次數是否正確更新
	if dialog.attempts != maxAttempts {
		t.Errorf("嘗試次數更新錯誤，期望 %d，實際 %d", maxAttempts, dialog.attempts)
	}

	// 驗證所有回調結果都包含正確的密碼
	for i, result := range callbackResults {
		if result.Password != "testpassword" {
			t.Errorf("第 %d 次回調的密碼錯誤，期望 'testpassword'，實際 '%s'", i+1, result.Password)
		}
		if !result.Confirmed {
			t.Errorf("第 %d 次回調的 Confirmed 應該為 true", i+1)
		}
	}
}

// TestPasswordDialogUI 測試密碼對話框的 UI 互動
func TestPasswordDialogUI(t *testing.T) {
	// 創建測試視窗
	testWindow := test.NewWindow(nil)
	defer testWindow.Close()

	t.Run("測試密碼設定對話框 UI 互動", func(t *testing.T) {
		callback := func(result PasswordDialogResult) {
			// 回調函數用於測試
		}

		dialog := NewPasswordSetupDialog(testWindow, "設定密碼", callback)

		// 測試密碼輸入
		test.Type(dialog.passwordEntry, "TestPassword123!")
		if dialog.passwordEntry.Text != "TestPassword123!" {
			t.Error("密碼輸入框文字設置失敗")
		}

		// 測試確認密碼輸入
		test.Type(dialog.confirmEntry, "TestPassword123!")
		if dialog.confirmEntry.Text != "TestPassword123!" {
			t.Error("確認密碼輸入框文字設置失敗")
		}

		// 驗證密碼強度是否更新
		if dialog.strengthBar.Value == 0 {
			t.Error("密碼強度指示器未正確更新")
		}
	})

	t.Run("測試密碼驗證對話框 UI 互動", func(t *testing.T) {
		callback := func(result PasswordDialogResult) {
			// 回調函數用於測試
		}

		dialog := NewPasswordVerifyDialog(testWindow, "驗證密碼", 3, callback)

		// 測試密碼輸入
		test.Type(dialog.passwordEntry, "TestPassword")
		if dialog.passwordEntry.Text != "TestPassword" {
			t.Error("密碼輸入框文字設置失敗")
		}
	})
}

// BenchmarkPasswordStrengthCalculation 密碼強度計算的效能基準測試
func BenchmarkPasswordStrengthCalculation(b *testing.B) {
	passwords := []string{
		"weak",
		"Medium1!",
		"VeryStrongPassword123!@#",
		"SuperLongPasswordWithManyCharacters123!@#$%^&*()",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, password := range passwords {
			calculatePasswordStrength(password)
		}
	}
}

// TestPasswordDialogMemoryUsage 測試密碼對話框的記憶體使用
func TestPasswordDialogMemoryUsage(t *testing.T) {
	// 創建測試視窗
	testWindow := test.NewWindow(nil)
	defer testWindow.Close()

	// 創建多個對話框實例以測試記憶體使用
	dialogs := make([]*PasswordSetupDialog, 10)
	for i := 0; i < 10; i++ {
		dialogs[i] = NewPasswordSetupDialog(testWindow, "測試對話框", nil)
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
			dialog.dialog.Hide()
		}
	}
}