package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

// AuthMethod 驗證方法類型
// 定義不同的驗證方法
type AuthMethod int

const (
	// AuthMethodPassword 密碼驗證
	AuthMethodPassword AuthMethod = iota
	// AuthMethodBiometric 生物驗證
	AuthMethodBiometric
	// AuthMethodBoth 兩種驗證方法都支援
	AuthMethodBoth
)

// AuthResult 驗證結果
// 包含驗證的完整結果資訊
type AuthResult struct {
	Success      bool       // 驗證是否成功
	Password     string     // 密碼（如果使用密碼驗證）
	Method       AuthMethod // 使用的驗證方法
	ErrorMessage string     // 錯誤訊息（如果有）
	Cancelled    bool       // 用戶是否取消了驗證
}

// AuthCallback 驗證回調函數類型
// 當驗證完成時調用
type AuthCallback func(result AuthResult)

// AuthDialogManager 驗證對話框管理器
// 統一管理密碼和生物驗證對話框
type AuthDialogManager struct {
	parent              fyne.Window           // 父視窗
	biometricAvailable  bool                  // 生物驗證是否可用
	preferredMethod     AuthMethod            // 首選驗證方法
	allowFallback       bool                  // 是否允許回退到其他驗證方法
	maxPasswordAttempts int                   // 最大密碼嘗試次數
}

// NewAuthDialogManager 創建新的驗證對話框管理器
// 參數：
//   - parent: 父視窗
//   - biometricAvailable: 生物驗證是否可用
//   - preferredMethod: 首選驗證方法
//   - allowFallback: 是否允許回退到其他驗證方法
//   - maxPasswordAttempts: 最大密碼嘗試次數
// 回傳：驗證對話框管理器實例
//
// 執行流程：
// 1. 初始化管理器參數
// 2. 驗證參數的有效性
// 3. 設置預設值
func NewAuthDialogManager(parent fyne.Window, biometricAvailable bool, preferredMethod AuthMethod, allowFallback bool, maxPasswordAttempts int) *AuthDialogManager {
	if maxPasswordAttempts <= 0 {
		maxPasswordAttempts = 3 // 預設最大嘗試次數
	}

	return &AuthDialogManager{
		parent:              parent,
		biometricAvailable:  biometricAvailable,
		preferredMethod:     preferredMethod,
		allowFallback:       allowFallback,
		maxPasswordAttempts: maxPasswordAttempts,
	}
}

// ShowAuthDialog 顯示驗證對話框
// 參數：
//   - title: 對話框標題
//   - message: 驗證提示訊息
//   - callback: 完成時的回調函數
//
// 執行流程：
// 1. 根據首選方法和可用性決定顯示哪種對話框
// 2. 如果首選方法不可用，根據回退設置選擇替代方法
// 3. 顯示相應的驗證對話框
// 4. 處理驗證結果和回退邏輯
func (m *AuthDialogManager) ShowAuthDialog(title, message string, callback AuthCallback) {
	// 根據首選方法和可用性決定驗證方式
	switch m.preferredMethod {
	case AuthMethodBiometric:
		if m.biometricAvailable {
			m.showBiometricDialog(title, message, callback)
		} else if m.allowFallback {
			m.showPasswordDialog(title, "生物驗證不可用，請輸入密碼", callback)
		} else {
			m.handleAuthError("生物驗證不可用且未啟用密碼回退", callback)
		}

	case AuthMethodPassword:
		m.showPasswordDialog(title, message, callback)

	case AuthMethodBoth:
		if m.biometricAvailable {
			m.showBiometricDialog(title, message, callback)
		} else {
			m.showPasswordDialog(title, message, callback)
		}

	default:
		m.handleAuthError("未知的驗證方法", callback)
	}
}

// ShowPasswordSetupDialog 顯示密碼設定對話框
// 參數：
//   - title: 對話框標題
//   - callback: 完成時的回調函數
//
// 執行流程：
// 1. 創建密碼設定對話框
// 2. 處理設定結果
// 3. 調用回調函數
func (m *AuthDialogManager) ShowPasswordSetupDialog(title string, callback AuthCallback) {
	ShowPasswordSetupDialog(m.parent, title, func(result PasswordDialogResult) {
		authResult := AuthResult{
			Success:      result.Confirmed,
			Password:     result.Password,
			Method:       AuthMethodPassword,
			ErrorMessage: "",
			Cancelled:    !result.Confirmed,
		}
		callback(authResult)
	})
}

// ShowBiometricSetupDialog 顯示生物驗證設定對話框
// 參數：
//   - title: 對話框標題
//   - callback: 完成時的回調函數
//
// 執行流程：
// 1. 創建生物驗證設定對話框
// 2. 處理設定結果
// 3. 調用回調函數
func (m *AuthDialogManager) ShowBiometricSetupDialog(title string, callback AuthCallback) {
	ShowBiometricSetupDialog(m.parent, title, m.biometricAvailable, func(result BiometricAuthResult) {
		authResult := AuthResult{
			Success:      result.Success,
			Password:     "",
			Method:       AuthMethodBiometric,
			ErrorMessage: result.ErrorMessage,
			Cancelled:    result.UserCancelled,
		}
		callback(authResult)
	})
}

// showBiometricDialog 顯示生物驗證對話框
// 參數：
//   - title: 對話框標題
//   - message: 驗證提示訊息
//   - callback: 完成時的回調函數
//
// 執行流程：
// 1. 創建生物驗證對話框
// 2. 處理驗證結果
// 3. 如果失敗且允許回退，顯示密碼對話框
// 4. 調用回調函數
func (m *AuthDialogManager) showBiometricDialog(title, message string, callback AuthCallback) {
	ShowBiometricAuthDialog(m.parent, title, message, m.allowFallback, func(result BiometricAuthResult) {
		if result.Success {
			// 生物驗證成功
			callback(AuthResult{
				Success:      true,
				Password:     "",
				Method:       AuthMethodBiometric,
				ErrorMessage: "",
				Cancelled:    false,
			})
		} else if result.FallbackUsed && m.allowFallback {
			// 用戶選擇使用密碼回退
			m.showPasswordDialog(title, "請輸入密碼", callback)
		} else if result.UserCancelled {
			// 用戶取消了驗證
			callback(AuthResult{
				Success:      false,
				Password:     "",
				Method:       AuthMethodBiometric,
				ErrorMessage: "",
				Cancelled:    true,
			})
		} else {
			// 生物驗證失敗
			callback(AuthResult{
				Success:      false,
				Password:     "",
				Method:       AuthMethodBiometric,
				ErrorMessage: result.ErrorMessage,
				Cancelled:    false,
			})
		}
	})
}

// showPasswordDialog 顯示密碼驗證對話框
// 參數：
//   - title: 對話框標題
//   - message: 驗證提示訊息
//   - callback: 完成時的回調函數
//
// 執行流程：
// 1. 創建密碼驗證對話框
// 2. 處理驗證結果
// 3. 調用回調函數
func (m *AuthDialogManager) showPasswordDialog(title, message string, callback AuthCallback) {
	ShowPasswordVerifyDialog(m.parent, title, m.maxPasswordAttempts, func(result PasswordDialogResult) {
		callback(AuthResult{
			Success:      result.Confirmed,
			Password:     result.Password,
			Method:       AuthMethodPassword,
			ErrorMessage: "",
			Cancelled:    !result.Confirmed,
		})
	})
}

// handleAuthError 處理驗證錯誤
// 參數：
//   - errorMessage: 錯誤訊息
//   - callback: 完成時的回調函數
//
// 執行流程：
// 1. 顯示錯誤對話框
// 2. 調用回調函數報告錯誤
func (m *AuthDialogManager) handleAuthError(errorMessage string, callback AuthCallback) {
	dialog.ShowError(fmt.Errorf(errorMessage), m.parent)
	callback(AuthResult{
		Success:      false,
		Password:     "",
		Method:       AuthMethodPassword, // 預設方法
		ErrorMessage: errorMessage,
		Cancelled:    false,
	})
}

// SetBiometricAvailable 設置生物驗證可用性
// 參數：
//   - available: 生物驗證是否可用
//
// 用於動態更新生物驗證的可用性狀態
func (m *AuthDialogManager) SetBiometricAvailable(available bool) {
	m.biometricAvailable = available
}

// SetPreferredMethod 設置首選驗證方法
// 參數：
//   - method: 首選驗證方法
//
// 用於動態更改首選的驗證方法
func (m *AuthDialogManager) SetPreferredMethod(method AuthMethod) {
	m.preferredMethod = method
}

// SetAllowFallback 設置是否允許回退
// 參數：
//   - allow: 是否允許回退到其他驗證方法
//
// 用於動態控制是否允許驗證方法之間的回退
func (m *AuthDialogManager) SetAllowFallback(allow bool) {
	m.allowFallback = allow
}

// SetMaxPasswordAttempts 設置最大密碼嘗試次數
// 參數：
//   - maxAttempts: 最大嘗試次數
//
// 用於動態調整密碼驗證的最大嘗試次數
func (m *AuthDialogManager) SetMaxPasswordAttempts(maxAttempts int) {
	if maxAttempts > 0 {
		m.maxPasswordAttempts = maxAttempts
	}
}

// GetBiometricAvailable 獲取生物驗證可用性
// 回傳：生物驗證是否可用
func (m *AuthDialogManager) GetBiometricAvailable() bool {
	return m.biometricAvailable
}

// GetPreferredMethod 獲取首選驗證方法
// 回傳：首選驗證方法
func (m *AuthDialogManager) GetPreferredMethod() AuthMethod {
	return m.preferredMethod
}

// GetAllowFallback 獲取是否允許回退
// 回傳：是否允許回退到其他驗證方法
func (m *AuthDialogManager) GetAllowFallback() bool {
	return m.allowFallback
}

// GetMaxPasswordAttempts 獲取最大密碼嘗試次數
// 回傳：最大密碼嘗試次數
func (m *AuthDialogManager) GetMaxPasswordAttempts() int {
	return m.maxPasswordAttempts
}

// ShowQuickAuthDialog 顯示快速驗證對話框的便利函數
// 參數：
//   - parent: 父視窗
//   - title: 對話框標題
//   - message: 驗證提示訊息
//   - biometricAvailable: 生物驗證是否可用
//   - callback: 完成時的回調函數
//
// 執行流程：
// 1. 創建臨時的驗證對話框管理器
// 2. 使用預設設置顯示驗證對話框
func ShowQuickAuthDialog(parent fyne.Window, title, message string, biometricAvailable bool, callback AuthCallback) {
	manager := NewAuthDialogManager(parent, biometricAvailable, AuthMethodBoth, true, 3)
	manager.ShowAuthDialog(title, message, callback)
}

// ShowQuickPasswordSetup 顯示快速密碼設定的便利函數
// 參數：
//   - parent: 父視窗
//   - title: 對話框標題
//   - callback: 完成時的回調函數
//
// 執行流程：
// 1. 直接顯示密碼設定對話框
func ShowQuickPasswordSetup(parent fyne.Window, title string, callback AuthCallback) {
	ShowPasswordSetupDialog(parent, title, func(result PasswordDialogResult) {
		callback(AuthResult{
			Success:      result.Confirmed,
			Password:     result.Password,
			Method:       AuthMethodPassword,
			ErrorMessage: "",
			Cancelled:    !result.Confirmed,
		})
	})
}

// ShowQuickBiometricSetup 顯示快速生物驗證設定的便利函數
// 參數：
//   - parent: 父視窗
//   - title: 對話框標題
//   - biometricAvailable: 生物驗證是否可用
//   - callback: 完成時的回調函數
//
// 執行流程：
// 1. 直接顯示生物驗證設定對話框
func ShowQuickBiometricSetup(parent fyne.Window, title string, biometricAvailable bool, callback AuthCallback) {
	ShowBiometricSetupDialog(parent, title, biometricAvailable, func(result BiometricAuthResult) {
		callback(AuthResult{
			Success:      result.Success,
			Password:     "",
			Method:       AuthMethodBiometric,
			ErrorMessage: result.ErrorMessage,
			Cancelled:    result.UserCancelled,
		})
	})
}