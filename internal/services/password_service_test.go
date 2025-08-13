// Package services 提供密碼驗證系統的單元測試
// 測試密碼雜湊、驗證、強度檢查和重試機制的正確性
package services

import (
	"strings"
	"testing"
	"time"
)

// TestNewPasswordService 測試密碼服務的建立
// 驗證服務實例是否正確建立並實作了 PasswordService 介面
func TestNewPasswordService(t *testing.T) {
	service := NewPasswordService()
	
	if service == nil {
		t.Fatal("NewPasswordService() 回傳 nil")
	}
	
	// 驗證是否實作了 PasswordService 介面
	_, ok := service.(PasswordService)
	if !ok {
		t.Fatal("回傳的服務未實作 PasswordService 介面")
	}
}

// TestHashPassword 測試密碼雜湊功能
// 驗證密碼雜湊的正確性和安全性
func TestHashPassword(t *testing.T) {
	service := NewPasswordService()
	
	testCases := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "正常密碼雜湊測試",
			password: "TestPassword123!",
			wantErr:  false,
		},
		{
			name:     "長密碼雜湊測試",
			password: "VeryLongPasswordWithManyCharacters123!@#$%^&*()",
			wantErr:  false,
		},
		{
			name:     "中文密碼雜湊測試",
			password: "中文密碼123!",
			wantErr:  false,
		},
		{
			name:     "空密碼測試",
			password: "",
			wantErr:  true,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hash, err := service.HashPassword(tc.password)
			
			if tc.wantErr {
				if err == nil {
					t.Errorf("期望發生錯誤，但沒有錯誤發生")
				}
				return
			}
			
			if err != nil {
				t.Fatalf("雜湊密碼失敗: %v", err)
			}
			
			if hash == nil {
				t.Fatal("雜湊結果為 nil")
			}
			
			// 驗證雜湊結構
			if hash.Salt == "" {
				t.Error("鹽值不能為空")
			}
			
			if hash.Hash == "" {
				t.Error("雜湊值不能為空")
			}
			
			if hash.Algorithm != "pbkdf2-sha256" {
				t.Errorf("演算法不正確: 期望 pbkdf2-sha256，實際 %s", hash.Algorithm)
			}
			
			if hash.Rounds != PasswordPBKDF2Rounds {
				t.Errorf("迭代次數不正確: 期望 %d，實際 %d", PasswordPBKDF2Rounds, hash.Rounds)
			}
			
			if hash.CreatedAt.IsZero() {
				t.Error("建立時間不能為零值")
			}
		})
	}
}

// TestVerifyPassword 測試密碼驗證功能
// 驗證密碼驗證的正確性和錯誤處理
func TestVerifyPassword(t *testing.T) {
	service := NewPasswordService()
	
	testCases := []struct {
		name     string
		password string
	}{
		{
			name:     "正常密碼驗證測試",
			password: "VerifyTest123!",
		},
		{
			name:     "特殊字元密碼驗證測試",
			password: "Special!@#$%^&*()_+-=[]{}|;:,.<>?",
		},
		{
			name:     "中文密碼驗證測試",
			password: "中文驗證密碼123!",
		},
		{
			name:     "長密碼驗證測試",
			password: strings.Repeat("LongPassword123!", 5),
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 先雜湊密碼
			hash, err := service.HashPassword(tc.password)
			if err != nil {
				t.Fatalf("雜湊密碼失敗: %v", err)
			}
			
			// 驗證正確密碼
			valid, err := service.VerifyPassword(tc.password, hash)
			if err != nil {
				t.Fatalf("驗證密碼失敗: %v", err)
			}
			
			if !valid {
				t.Error("正確密碼驗證失敗")
			}
			
			// 驗證錯誤密碼
			wrongPassword := tc.password + "wrong"
			valid, err = service.VerifyPassword(wrongPassword, hash)
			if err != nil {
				t.Fatalf("驗證錯誤密碼時發生錯誤: %v", err)
			}
			
			if valid {
				t.Error("錯誤密碼驗證應該失敗")
			}
		})
	}
}

// TestVerifyPassword_ErrorCases 測試密碼驗證的錯誤情況
// 驗證各種無效輸入的錯誤處理
func TestVerifyPassword_ErrorCases(t *testing.T) {
	service := NewPasswordService()
	
	// 建立有效的密碼雜湊
	validHash, err := service.HashPassword("ValidPassword123!")
	if err != nil {
		t.Fatalf("建立有效雜湊失敗: %v", err)
	}
	
	testCases := []struct {
		name        string
		password    string
		hash        *PasswordHash
		expectError string
	}{
		{
			name:        "空密碼測試",
			password:    "",
			hash:        validHash,
			expectError: "密碼不能為空",
		},
		{
			name:        "空雜湊測試",
			password:    "ValidPassword123!",
			hash:        nil,
			expectError: "密碼雜湊不能為空",
		},
		{
			name:     "不支援的演算法測試",
			password: "ValidPassword123!",
			hash: &PasswordHash{
				Salt:      validHash.Salt,
				Hash:      validHash.Hash,
				Algorithm: "unsupported",
				Rounds:    validHash.Rounds,
				CreatedAt: validHash.CreatedAt,
			},
			expectError: "不支援的雜湊演算法",
		},
		{
			name:     "無效鹽值測試",
			password: "ValidPassword123!",
			hash: &PasswordHash{
				Salt:      "invalid-base64!",
				Hash:      validHash.Hash,
				Algorithm: "pbkdf2-sha256",
				Rounds:    validHash.Rounds,
				CreatedAt: validHash.CreatedAt,
			},
			expectError: "解碼鹽值失敗",
		},
		{
			name:     "無效雜湊值測試",
			password: "ValidPassword123!",
			hash: &PasswordHash{
				Salt:      validHash.Salt,
				Hash:      "invalid-base64!",
				Algorithm: "pbkdf2-sha256",
				Rounds:    validHash.Rounds,
				CreatedAt: validHash.CreatedAt,
			},
			expectError: "解碼雜湊值失敗",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := service.VerifyPassword(tc.password, tc.hash)
			
			if err == nil {
				t.Error("期望發生錯誤，但沒有錯誤發生")
				return
			}
			
			if !strings.Contains(err.Error(), tc.expectError) {
				t.Errorf("錯誤訊息不匹配:\n期望包含: %s\n實際: %s", tc.expectError, err.Error())
			}
		})
	}
}

// TestCheckPasswordStrength 測試密碼強度檢查功能
// 驗證各種密碼強度的正確評估
func TestCheckPasswordStrength(t *testing.T) {
	service := NewPasswordService()
	
	testCases := []struct {
		name             string
		password         string
		expectedStrength PasswordStrength
		expectSuggestions bool
	}{
		{
			name:             "強密碼測試",
			password:         "StrongPassword123!@#",
			expectedStrength: PasswordStrong,
			expectSuggestions: false,
		},
		{
			name:             "良好密碼測試",
			password:         "GoodPass123!",
			expectedStrength: PasswordGood,
			expectSuggestions: false,
		},
		{
			name:             "一般密碼測試",
			password:         "fairpass123",
			expectedStrength: PasswordFair,
			expectSuggestions: true,
		},
		{
			name:             "弱密碼測試（太短）",
			password:         "weak",
			expectedStrength: PasswordWeak,
			expectSuggestions: true,
		},
		{
			name:             "弱密碼測試（常見密碼）",
			password:         "password123",
			expectedStrength: PasswordWeak,
			expectSuggestions: true,
		},
		{
			name:             "缺少大寫字母",
			password:         "lowercase123!",
			expectedStrength: PasswordFair,
			expectSuggestions: true,
		},
		{
			name:             "缺少小寫字母",
			password:         "UPPERCASE123!",
			expectedStrength: PasswordFair,
			expectSuggestions: true,
		},
		{
			name:             "缺少數字",
			password:         "NoNumbers!@#",
			expectedStrength: PasswordFair,
			expectSuggestions: true,
		},
		{
			name:             "缺少特殊字元",
			password:         "NoSpecialChars123",
			expectedStrength: PasswordFair,
			expectSuggestions: true,
		},
		{
			name:             "重複字元密碼",
			password:         "Password111!",
			expectedStrength: PasswordGood,
			expectSuggestions: true,
		},
		{
			name:             "中文密碼測試",
			password:         "中文密碼123!Aa",
			expectedStrength: PasswordStrong,
			expectSuggestions: false,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			strength, suggestions := service.CheckPasswordStrength(tc.password)
			
			if strength != tc.expectedStrength {
				t.Errorf("密碼強度不正確:\n密碼: %s\n期望: %s\n實際: %s",
					tc.password, tc.expectedStrength.String(), strength.String())
			}
			
			if tc.expectSuggestions && len(suggestions) == 0 {
				t.Error("期望有改進建議，但沒有建議")
			}
			
			if !tc.expectSuggestions && len(suggestions) > 1 {
				// 強密碼可能有一個 "密碼強度良好" 的建議
				hasOnlyGoodMessage := len(suggestions) == 1 && strings.Contains(suggestions[0], "密碼強度良好")
				if !hasOnlyGoodMessage {
					t.Errorf("不期望有改進建議，但有建議: %v", suggestions)
				}
			}
		})
	}
}

// TestRetryMechanism 測試密碼重試機制
// 驗證重試次數限制和鎖定功能
func TestRetryMechanism(t *testing.T) {
	service := NewPasswordService()
	identifier := "test-user"
	
	// 初始狀態應該沒有鎖定
	locked, duration := service.IsLocked(identifier)
	if locked {
		t.Error("初始狀態不應該被鎖定")
	}
	if duration != time.Duration(0) {
		t.Error("初始狀態鎖定時間應該為 0")
	}
	
	// 記錄失敗嘗試（未達到最大次數）
	for i := 0; i < MaxRetryAttempts-1; i++ {
		err := service.RecordFailedAttempt(identifier)
		if err != nil {
			t.Fatalf("記錄失敗嘗試 %d 時發生錯誤: %v", i+1, err)
		}
		
		// 檢查未被鎖定
		locked, _ := service.IsLocked(identifier)
		if locked {
			t.Errorf("嘗試 %d 次後不應該被鎖定", i+1)
		}
	}
	
	// 記錄最後一次失敗嘗試（應該觸發鎖定）
	err := service.RecordFailedAttempt(identifier)
	if err == nil {
		t.Error("達到最大重試次數時應該回傳錯誤")
	}
	
	if !strings.Contains(err.Error(), "帳戶已鎖定") {
		t.Errorf("錯誤訊息應該包含鎖定資訊: %v", err)
	}
	
	// 檢查已被鎖定
	locked, duration = service.IsLocked(identifier)
	if !locked {
		t.Error("達到最大重試次數後應該被鎖定")
	}
	
	if duration <= time.Duration(0) {
		t.Error("鎖定時間應該大於 0")
	}
	
	// 在鎖定期間嘗試記錄失敗應該回傳錯誤
	err = service.RecordFailedAttempt(identifier)
	if err == nil {
		t.Error("在鎖定期間記錄失敗嘗試應該回傳錯誤")
	}
	
	// 重置重試計數
	service.ResetRetryCount(identifier)
	
	// 檢查已解除鎖定
	locked, _ = service.IsLocked(identifier)
	if locked {
		t.Error("重置後不應該被鎖定")
	}
}

// TestGetRetryInfo 測試重試資訊取得功能
// 驗證重試資訊的正確性和完整性
func TestGetRetryInfo(t *testing.T) {
	service := NewPasswordService()
	identifier := "test-retry-info"
	
	// 初始狀態應該沒有重試資訊
	info := service.GetRetryInfo(identifier)
	if info != nil {
		t.Error("初始狀態不應該有重試資訊")
	}
	
	// 記錄一次失敗嘗試
	err := service.RecordFailedAttempt(identifier)
	if err != nil {
		t.Fatalf("記錄失敗嘗試失敗: %v", err)
	}
	
	// 取得重試資訊
	info = service.GetRetryInfo(identifier)
	if info == nil {
		t.Fatal("應該有重試資訊")
	}
	
	if info.Attempts != 1 {
		t.Errorf("重試次數不正確: 期望 1，實際 %d", info.Attempts)
	}
	
	if info.LastAttempt.IsZero() {
		t.Error("最後嘗試時間不應該為零值")
	}
	
	if !info.LockedUntil.IsZero() {
		t.Error("未達到最大重試次數時不應該有鎖定時間")
	}
	
	// 記錄更多失敗嘗試直到鎖定
	for i := 1; i < MaxRetryAttempts; i++ {
		service.RecordFailedAttempt(identifier)
	}
	
	// 檢查鎖定狀態的重試資訊
	info = service.GetRetryInfo(identifier)
	if info.Attempts != MaxRetryAttempts {
		t.Errorf("重試次數不正確: 期望 %d，實際 %d", MaxRetryAttempts, info.Attempts)
	}
	
	if info.LockedUntil.IsZero() {
		t.Error("達到最大重試次數時應該有鎖定時間")
	}
}

// TestRecordFailedAttempt_EmptyIdentifier 測試空識別符的錯誤處理
// 驗證對無效輸入的錯誤處理
func TestRecordFailedAttempt_EmptyIdentifier(t *testing.T) {
	service := NewPasswordService()
	
	err := service.RecordFailedAttempt("")
	if err == nil {
		t.Error("空識別符應該回傳錯誤")
	}
	
	if !strings.Contains(err.Error(), "識別符不能為空") {
		t.Errorf("錯誤訊息不正確: %v", err)
	}
}

// TestPasswordStrength_String 測試密碼強度字串表示
// 驗證密碼強度等級的字串轉換
func TestPasswordStrength_String(t *testing.T) {
	testCases := []struct {
		strength PasswordStrength
		expected string
	}{
		{PasswordWeak, "弱"},
		{PasswordFair, "一般"},
		{PasswordGood, "良好"},
		{PasswordStrong, "強"},
		{PasswordStrength(999), "未知"},
	}
	
	for _, tc := range testCases {
		t.Run(tc.expected, func(t *testing.T) {
			result := tc.strength.String()
			if result != tc.expected {
				t.Errorf("字串表示不正確: 期望 %s，實際 %s", tc.expected, result)
			}
		})
	}
}

// TestConcurrentAccess 測試並發存取的安全性
// 驗證多個 goroutine 同時存取時的執行緒安全性
func TestConcurrentAccess(t *testing.T) {
	service := NewPasswordService()
	identifier := "concurrent-test"
	
	// 啟動多個 goroutine 同時記錄失敗嘗試
	done := make(chan bool, 10)
	
	for i := 0; i < 10; i++ {
		go func() {
			defer func() { done <- true }()
			service.RecordFailedAttempt(identifier)
		}()
	}
	
	// 等待所有 goroutine 完成
	for i := 0; i < 10; i++ {
		<-done
	}
	
	// 檢查最終狀態
	info := service.GetRetryInfo(identifier)
	if info == nil {
		t.Fatal("應該有重試資訊")
	}
	
	// 由於並發存取，重試次數應該至少為最大值
	if info.Attempts < MaxRetryAttempts {
		t.Errorf("並發存取後重試次數不正確: 期望至少 %d，實際 %d", MaxRetryAttempts, info.Attempts)
	}
}

// BenchmarkHashPassword 效能測試：密碼雜湊
// 測試密碼雜湊的效能表現
func BenchmarkHashPassword(b *testing.B) {
	service := NewPasswordService()
	password := "BenchmarkPassword123!"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.HashPassword(password)
		if err != nil {
			b.Fatalf("雜湊密碼失敗: %v", err)
		}
	}
}

// BenchmarkVerifyPassword 效能測試：密碼驗證
// 測試密碼驗證的效能表現
func BenchmarkVerifyPassword(b *testing.B) {
	service := NewPasswordService()
	password := "BenchmarkPassword123!"
	
	// 預先雜湊密碼
	hash, err := service.HashPassword(password)
	if err != nil {
		b.Fatalf("預先雜湊密碼失敗: %v", err)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.VerifyPassword(password, hash)
		if err != nil {
			b.Fatalf("驗證密碼失敗: %v", err)
		}
	}
}

// BenchmarkCheckPasswordStrength 效能測試：密碼強度檢查
// 測試密碼強度檢查的效能表現
func BenchmarkCheckPasswordStrength(b *testing.B) {
	service := NewPasswordService()
	password := "BenchmarkPassword123!"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.CheckPasswordStrength(password)
	}
}