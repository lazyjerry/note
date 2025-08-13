// Package services 提供生物識別驗證服務的整合測試
// 測試 macOS 生物識別驗證功能的正確性和錯誤處理
package services

import (
	"errors"
	"runtime"
	"strings"
	"testing"
	"time"
)

// TestNewBiometricService 測試生物識別服務的建立
// 驗證服務實例是否正確建立並實作了 BiometricService 介面
func TestNewBiometricService(t *testing.T) {
	service := NewBiometricService()
	
	if service == nil {
		t.Fatal("NewBiometricService() 回傳 nil")
	}
	
	// 驗證是否實作了 BiometricService 介面
	_, ok := service.(BiometricService)
	if !ok {
		t.Fatal("回傳的服務未實作 BiometricService 介面")
	}
}

// TestIsAvailable 測試生物識別可用性檢查
// 驗證在不同平台上的可用性檢查結果
func TestIsAvailable(t *testing.T) {
	service := NewBiometricService()
	
	available, biometricType := service.IsAvailable()
	
	// 在 macOS 上可能可用，在其他平台上應該不可用
	if runtime.GOOS == "darwin" {
		// macOS 系統：可能可用也可能不可用（取決於硬體和設定）
		t.Logf("macOS 系統生物識別可用性: %v, 類型: %s", available, biometricType.String())
		
		if available {
			// 如果可用，類型應該是 Touch ID 或 Face ID
			if biometricType != BiometricTypeTouchID && biometricType != BiometricTypeFaceID {
				t.Errorf("可用時生物識別類型應該是 Touch ID 或 Face ID，實際: %s", biometricType.String())
			}
		} else {
			// 如果不可用，類型應該是 None
			if biometricType != BiometricTypeNone {
				t.Errorf("不可用時生物識別類型應該是 None，實際: %s", biometricType.String())
			}
		}
	} else {
		// 非 macOS 系統：應該不可用
		if available {
			t.Error("非 macOS 系統不應該支援生物識別驗證")
		}
		
		if biometricType != BiometricTypeNone {
			t.Errorf("非 macOS 系統生物識別類型應該是 None，實際: %s", biometricType.String())
		}
	}
}

// TestBiometricType_String 測試生物識別類型的字串表示
// 驗證各種生物識別類型的字串轉換
func TestBiometricType_String(t *testing.T) {
	testCases := []struct {
		biometricType BiometricType
		expected      string
	}{
		{BiometricTypeNone, "無"},
		{BiometricTypeTouchID, "Touch ID"},
		{BiometricTypeFaceID, "Face ID"},
		{BiometricType(999), "未知"},
	}
	
	for _, tc := range testCases {
		t.Run(tc.expected, func(t *testing.T) {
			result := tc.biometricType.String()
			if result != tc.expected {
				t.Errorf("字串表示不正確: 期望 %s，實際 %s", tc.expected, result)
			}
		})
	}
}

// TestAuthenticate_NotAvailable 測試生物識別不可用時的驗證行為
// 驗證在生物識別不可用時的錯誤處理
func TestAuthenticate_NotAvailable(t *testing.T) {
	service := NewBiometricService()
	
	// 檢查生物識別是否可用
	available, _ := service.IsAvailable()
	
	if !available {
		// 如果生物識別不可用，測試驗證行為
		result := service.Authenticate("測試驗證")
		
		if result == nil {
			t.Fatal("驗證結果不應該為 nil")
		}
		
		if result.Success {
			t.Error("生物識別不可用時驗證不應該成功")
		}
		
		if result.Cancelled {
			t.Error("生物識別不可用時不應該是取消狀態")
		}
		
		if result.Error == nil {
			t.Error("生物識別不可用時應該有錯誤")
		}
		
		if result.Duration < 0 {
			t.Error("驗證耗時不應該為負數")
		}
		
		// 檢查錯誤訊息
		expectedMessages := []string{"不可用", "不支援"}
		errorMsg := result.Error.Error()
		hasExpectedMessage := false
		for _, msg := range expectedMessages {
			if strings.Contains(errorMsg, msg) {
				hasExpectedMessage = true
				break
			}
		}
		
		if !hasExpectedMessage {
			t.Errorf("錯誤訊息應該包含不可用或不支援的資訊: %s", errorMsg)
		}
	} else {
		t.Skip("生物識別可用，跳過不可用測試")
	}
}

// TestAuthenticateForNote_EmptyNoteID 測試空筆記 ID 的錯誤處理
// 驗證對無效輸入的錯誤處理
func TestAuthenticateForNote_EmptyNoteID(t *testing.T) {
	service := NewBiometricService()
	
	result := service.AuthenticateForNote("", "測試原因")
	
	if result == nil {
		t.Fatal("驗證結果不應該為 nil")
	}
	
	if result.Success {
		t.Error("空筆記 ID 驗證不應該成功")
	}
	
	if result.Cancelled {
		t.Error("空筆記 ID 不應該是取消狀態")
	}
	
	if result.Error == nil {
		t.Error("空筆記 ID 應該有錯誤")
	}
	
	if !strings.Contains(result.Error.Error(), "筆記 ID 不能為空") {
		t.Errorf("錯誤訊息不正確: %v", result.Error)
	}
	
	if result.Duration != 0 {
		t.Error("空筆記 ID 驗證耗時應該為 0")
	}
}

// TestAuthenticateForNote_NotEnabled 測試未啟用生物識別的筆記驗證
// 驗證對未啟用筆記的錯誤處理
func TestAuthenticateForNote_NotEnabled(t *testing.T) {
	service := NewBiometricService()
	noteID := "test-note-not-enabled"
	
	// 確保筆記未啟用生物識別
	if service.IsEnabledForNote(noteID) {
		t.Fatal("測試筆記不應該已啟用生物識別")
	}
	
	result := service.AuthenticateForNote(noteID, "測試原因")
	
	if result == nil {
		t.Fatal("驗證結果不應該為 nil")
	}
	
	if result.Success {
		t.Error("未啟用筆記驗證不應該成功")
	}
	
	if result.Cancelled {
		t.Error("未啟用筆記不應該是取消狀態")
	}
	
	if result.Error == nil {
		t.Error("未啟用筆記應該有錯誤")
	}
	
	if !strings.Contains(result.Error.Error(), "未啟用生物識別驗證") {
		t.Errorf("錯誤訊息不正確: %v", result.Error)
	}
	
	if result.Duration != 0 {
		t.Error("未啟用筆記驗證耗時應該為 0")
	}
}

// TestSetupForNote_EmptyNoteID 測試空筆記 ID 的設定錯誤處理
// 驗證對無效輸入的錯誤處理
func TestSetupForNote_EmptyNoteID(t *testing.T) {
	service := NewBiometricService()
	
	err := service.SetupForNote("")
	
	if err == nil {
		t.Error("空筆記 ID 設定應該回傳錯誤")
	}
	
	if !strings.Contains(err.Error(), "筆記 ID 不能為空") {
		t.Errorf("錯誤訊息不正確: %v", err)
	}
}

// TestSetupForNote_NotAvailable 測試生物識別不可用時的設定行為
// 驗證在生物識別不可用時的錯誤處理
func TestSetupForNote_NotAvailable(t *testing.T) {
	service := NewBiometricService()
	noteID := "test-note-setup"
	
	// 檢查生物識別是否可用
	available, _ := service.IsAvailable()
	
	if !available {
		// 如果生物識別不可用，測試設定行為
		err := service.SetupForNote(noteID)
		
		if err == nil {
			t.Error("生物識別不可用時設定應該回傳錯誤")
		}
		
		// 檢查錯誤訊息
		expectedMessages := []string{"不可用", "不支援"}
		errorMsg := err.Error()
		hasExpectedMessage := false
		for _, msg := range expectedMessages {
			if strings.Contains(errorMsg, msg) {
				hasExpectedMessage = true
				break
			}
		}
		
		if !hasExpectedMessage {
			t.Errorf("錯誤訊息應該包含不可用或不支援的資訊: %s", errorMsg)
		}
		
		// 確保筆記未被啟用
		if service.IsEnabledForNote(noteID) {
			t.Error("設定失敗時筆記不應該被啟用")
		}
	} else {
		t.Skip("生物識別可用，跳過不可用測試")
	}
}

// TestRemoveForNote 測試移除筆記生物識別設定
// 驗證移除功能的正確性
func TestRemoveForNote(t *testing.T) {
	service := NewBiometricService()
	noteID := "test-note-remove"
	
	// 測試移除不存在的筆記（應該成功）
	err := service.RemoveForNote(noteID)
	if err != nil {
		t.Errorf("移除不存在的筆記設定不應該回傳錯誤: %v", err)
	}
	
	// 測試空筆記 ID
	err = service.RemoveForNote("")
	if err == nil {
		t.Error("空筆記 ID 移除應該回傳錯誤")
	}
	
	if !strings.Contains(err.Error(), "筆記 ID 不能為空") {
		t.Errorf("錯誤訊息不正確: %v", err)
	}
}

// TestIsEnabledForNote 測試筆記生物識別啟用狀態檢查
// 驗證啟用狀態檢查的正確性
func TestIsEnabledForNote(t *testing.T) {
	service := NewBiometricService()
	
	// 測試空筆記 ID
	if service.IsEnabledForNote("") {
		t.Error("空筆記 ID 不應該啟用生物識別")
	}
	
	// 測試不存在的筆記
	noteID := "test-note-enabled-check"
	if service.IsEnabledForNote(noteID) {
		t.Error("不存在的筆記不應該啟用生物識別")
	}
	
	// 在非 macOS 系統上，所有筆記都應該回傳 false
	if runtime.GOOS != "darwin" {
		if service.IsEnabledForNote("any-note") {
			t.Error("非 macOS 系統上所有筆記都不應該啟用生物識別")
		}
	}
}

// TestEncryptionService_BiometricIntegration 測試加密服務的生物識別整合
// 驗證 EncryptionService 與生物識別功能的整合
func TestEncryptionService_BiometricIntegration(t *testing.T) {
	encService := NewEncryptionService()
	noteID := "test-encryption-biometric"
	
	// 測試設定生物識別驗證
	err := encService.SetupBiometricAuth(noteID)
	
	if runtime.GOOS == "darwin" {
		// macOS 系統：可能成功也可能失敗（取決於硬體和用戶操作）
		t.Logf("macOS 系統生物識別設定結果: %v", err)
	} else {
		// 非 macOS 系統：應該失敗
		if err == nil {
			t.Error("非 macOS 系統設定生物識別應該失敗")
		}
		
		if !strings.Contains(err.Error(), "不支援") {
			t.Errorf("錯誤訊息應該包含不支援資訊: %v", err)
		}
	}
	
	// 測試生物識別驗證
	success, err := encService.AuthenticateWithBiometric(noteID)
	
	if runtime.GOOS == "darwin" {
		// macOS 系統：結果取決於設定和用戶操作
		t.Logf("macOS 系統生物識別驗證結果: success=%v, error=%v", success, err)
	} else {
		// 非 macOS 系統：應該失敗
		if success {
			t.Error("非 macOS 系統生物識別驗證不應該成功")
		}
		
		if err == nil {
			t.Error("非 macOS 系統生物識別驗證應該回傳錯誤")
		}
		
		if !strings.Contains(err.Error(), "不支援") {
			t.Errorf("錯誤訊息應該包含不支援資訊: %v", err)
		}
	}
}

// TestBiometricResult_Structure 測試生物識別結果結構
// 驗證結果結構的完整性和正確性
func TestBiometricResult_Structure(t *testing.T) {
	// 建立測試結果
	result := &BiometricResult{
		Success:   true,
		Cancelled: false,
		Error:     nil,
		Duration:  100 * time.Millisecond,
	}
	
	// 驗證結構欄位
	if !result.Success {
		t.Error("Success 欄位設定不正確")
	}
	
	if result.Cancelled {
		t.Error("Cancelled 欄位設定不正確")
	}
	
	if result.Error != nil {
		t.Error("Error 欄位設定不正確")
	}
	
	if result.Duration != 100*time.Millisecond {
		t.Error("Duration 欄位設定不正確")
	}
	
	// 測試失敗結果
	failResult := &BiometricResult{
		Success:   false,
		Cancelled: true,
		Error:     errors.New("測試錯誤"),
		Duration:  50 * time.Millisecond,
	}
	
	if failResult.Success {
		t.Error("失敗結果 Success 應該為 false")
	}
	
	if !failResult.Cancelled {
		t.Error("取消結果 Cancelled 應該為 true")
	}
	
	if failResult.Error == nil {
		t.Error("錯誤結果 Error 不應該為 nil")
	}
	
	if failResult.Duration != 50*time.Millisecond {
		t.Error("失敗結果 Duration 設定不正確")
	}
}

// TestConcurrentAccess_BiometricService 測試生物識別服務的並發存取
// 驗證多個 goroutine 同時存取時的執行緒安全性
func TestConcurrentAccess_BiometricService(t *testing.T) {
	service := NewBiometricService()
	noteID := "concurrent-biometric-test"
	
	// 啟動多個 goroutine 同時檢查啟用狀態
	done := make(chan bool, 10)
	
	for i := 0; i < 10; i++ {
		go func() {
			defer func() { done <- true }()
			
			// 檢查啟用狀態
			service.IsEnabledForNote(noteID)
			
			// 嘗試移除設定
			service.RemoveForNote(noteID)
			
			// 再次檢查啟用狀態
			service.IsEnabledForNote(noteID)
		}()
	}
	
	// 等待所有 goroutine 完成
	for i := 0; i < 10; i++ {
		<-done
	}
	
	// 檢查最終狀態（應該是未啟用）
	if service.IsEnabledForNote(noteID) {
		t.Error("並發存取後筆記不應該啟用生物識別")
	}
}

// BenchmarkIsAvailable 效能測試：生物識別可用性檢查
// 測試可用性檢查的效能表現
func BenchmarkIsAvailable(b *testing.B) {
	service := NewBiometricService()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.IsAvailable()
	}
}

// BenchmarkIsEnabledForNote 效能測試：筆記啟用狀態檢查
// 測試啟用狀態檢查的效能表現
func BenchmarkIsEnabledForNote(b *testing.B) {
	service := NewBiometricService()
	noteID := "benchmark-note"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.IsEnabledForNote(noteID)
	}
}