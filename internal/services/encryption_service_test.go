// Package services 提供加密服務的單元測試
// 測試 AES-256 和 ChaCha20-Poly1305 加密演算法的正確性和安全性
package services

import (
	"encoding/json"
	"strings"
	"testing"
)

// TestNewEncryptionService 測試加密服務的建立
// 驗證服務實例是否正確建立並實作了 EncryptionService 介面
func TestNewEncryptionService(t *testing.T) {
	service := NewEncryptionService()
	
	if service == nil {
		t.Fatal("NewEncryptionService() 回傳 nil")
	}
	
	// 驗證是否實作了 EncryptionService 介面
	_, ok := service.(EncryptionService)
	if !ok {
		t.Fatal("回傳的服務未實作 EncryptionService 介面")
	}
}

// TestEncryptContent_AES256 測試 AES-256 加密功能
// 驗證加密過程的正確性和輸出格式
func TestEncryptContent_AES256(t *testing.T) {
	service := NewEncryptionService()
	
	testCases := []struct {
		name     string
		content  string
		password string
		wantErr  bool
	}{
		{
			name:     "正常加密測試",
			content:  "這是一個測試內容",
			password: "TestPass123!",
			wantErr:  false,
		},
		{
			name:     "長內容加密測試",
			content:  strings.Repeat("長內容測試", 1000),
			password: "LongContent456@",
			wantErr:  false,
		},
		{
			name:     "空內容測試",
			content:  "",
			password: "EmptyContent789#",
			wantErr:  true,
		},
		{
			name:     "空密碼測試",
			content:  "測試內容",
			password: "",
			wantErr:  true,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			encrypted, err := service.EncryptContent(tc.content, tc.password, AlgorithmAES256)
			
			if tc.wantErr {
				if err == nil {
					t.Errorf("期望發生錯誤，但沒有錯誤發生")
				}
				return
			}
			
			if err != nil {
				t.Fatalf("加密失敗: %v", err)
			}
			
			if len(encrypted) == 0 {
				t.Fatal("加密結果為空")
			}
			
			// 驗證加密結果是有效的 JSON
			var encData EncryptedData
			if err := json.Unmarshal(encrypted, &encData); err != nil {
				t.Fatalf("加密結果不是有效的 JSON: %v", err)
			}
			
			// 驗證加密資料結構
			if encData.Version != "1.0" {
				t.Errorf("版本不正確: 期望 1.0，實際 %s", encData.Version)
			}
			
			if encData.Algorithm != AlgorithmAES256 {
				t.Errorf("演算法不正確: 期望 %s，實際 %s", AlgorithmAES256, encData.Algorithm)
			}
			
			if encData.Salt == "" || encData.Nonce == "" || encData.Data == "" || encData.Checksum == "" {
				t.Error("加密資料結構不完整")
			}
		})
	}
}

// TestEncryptContent_ChaCha20 測試 ChaCha20-Poly1305 加密功能
// 驗證加密過程的正確性和輸出格式
func TestEncryptContent_ChaCha20(t *testing.T) {
	service := NewEncryptionService()
	
	testCases := []struct {
		name     string
		content  string
		password string
		wantErr  bool
	}{
		{
			name:     "正常加密測試",
			content:  "這是 ChaCha20 測試內容",
			password: "ChaCha20Pass123!",
			wantErr:  false,
		},
		{
			name:     "Unicode 內容測試",
			content:  "測試中文內容 🔒 加密功能",
			password: "Unicode456@",
			wantErr:  false,
		},
		{
			name:     "不支援的演算法測試",
			content:  "測試內容",
			password: "TestPass123!",
			wantErr:  true,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			algorithm := AlgorithmChaCha20
			if tc.name == "不支援的演算法測試" {
				algorithm = "unsupported"
			}
			
			encrypted, err := service.EncryptContent(tc.content, tc.password, algorithm)
			
			if tc.wantErr {
				if err == nil {
					t.Errorf("期望發生錯誤，但沒有錯誤發生")
				}
				return
			}
			
			if err != nil {
				t.Fatalf("加密失敗: %v", err)
			}
			
			// 驗證加密結果格式
			var encData EncryptedData
			if err := json.Unmarshal(encrypted, &encData); err != nil {
				t.Fatalf("加密結果不是有效的 JSON: %v", err)
			}
			
			if encData.Algorithm != AlgorithmChaCha20 {
				t.Errorf("演算法不正確: 期望 %s，實際 %s", AlgorithmChaCha20, encData.Algorithm)
			}
		})
	}
}

// TestDecryptContent_AES256 測試 AES-256 解密功能
// 驗證解密過程的正確性和錯誤處理
func TestDecryptContent_AES256(t *testing.T) {
	service := NewEncryptionService()
	
	testCases := []struct {
		name     string
		content  string
		password string
	}{
		{
			name:     "正常解密測試",
			content:  "這是解密測試內容",
			password: "DecryptTest123!",
		},
		{
			name:     "特殊字元解密測試",
			content:  "特殊字元: !@#$%^&*()_+-=[]{}|;:,.<>?",
			password: "SpecialChars456@",
		},
		{
			name:     "多行內容解密測試",
			content:  "第一行\n第二行\n第三行",
			password: "MultiLine789#",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 先加密
			encrypted, err := service.EncryptContent(tc.content, tc.password, AlgorithmAES256)
			if err != nil {
				t.Fatalf("加密失敗: %v", err)
			}
			
			// 再解密
			decrypted, err := service.DecryptContent(encrypted, tc.password, AlgorithmAES256)
			if err != nil {
				t.Fatalf("解密失敗: %v", err)
			}
			
			// 驗證解密結果
			if decrypted != tc.content {
				t.Errorf("解密結果不匹配:\n期望: %s\n實際: %s", tc.content, decrypted)
			}
		})
	}
}

// TestDecryptContent_ChaCha20 測試 ChaCha20-Poly1305 解密功能
// 驗證解密過程的正確性和錯誤處理
func TestDecryptContent_ChaCha20(t *testing.T) {
	service := NewEncryptionService()
	
	testCases := []struct {
		name     string
		content  string
		password string
	}{
		{
			name:     "ChaCha20 正常解密測試",
			content:  "ChaCha20 解密測試內容",
			password: "ChaCha20Decrypt123!",
		},
		{
			name:     "長內容解密測試",
			content:  strings.Repeat("ChaCha20 長內容測試 ", 500),
			password: "LongChaCha20Content456@",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 先加密
			encrypted, err := service.EncryptContent(tc.content, tc.password, AlgorithmChaCha20)
			if err != nil {
				t.Fatalf("加密失敗: %v", err)
			}
			
			// 再解密
			decrypted, err := service.DecryptContent(encrypted, tc.password, AlgorithmChaCha20)
			if err != nil {
				t.Fatalf("解密失敗: %v", err)
			}
			
			// 驗證解密結果
			if decrypted != tc.content {
				t.Errorf("解密結果不匹配:\n期望: %s\n實際: %s", tc.content, decrypted)
			}
		})
	}
}

// TestDecryptContent_WrongPassword 測試錯誤密碼的解密處理
// 驗證使用錯誤密碼時是否正確拒絕解密
func TestDecryptContent_WrongPassword(t *testing.T) {
	service := NewEncryptionService()
	
	content := "測試錯誤密碼處理"
	correctPassword := "CorrectPass123!"
	wrongPassword := "WrongPass456@"
	
	// 使用正確密碼加密
	encrypted, err := service.EncryptContent(content, correctPassword, AlgorithmAES256)
	if err != nil {
		t.Fatalf("加密失敗: %v", err)
	}
	
	// 使用錯誤密碼解密
	_, err = service.DecryptContent(encrypted, wrongPassword, AlgorithmAES256)
	if err == nil {
		t.Error("使用錯誤密碼解密應該失敗，但成功了")
	}
	
	// 驗證錯誤訊息包含解密失敗的資訊
	if !strings.Contains(err.Error(), "解密失敗") {
		t.Errorf("錯誤訊息不正確: %v", err)
	}
}

// TestDecryptContent_InvalidData 測試無效資料的解密處理
// 驗證對無效或損壞資料的錯誤處理
func TestDecryptContent_InvalidData(t *testing.T) {
	service := NewEncryptionService()
	
	testCases := []struct {
		name        string
		data        []byte
		password    string
		expectError string
	}{
		{
			name:        "空資料測試",
			data:        []byte{},
			password:    "TestPass123!",
			expectError: "加密資料不能為空",
		},
		{
			name:        "無效 JSON 測試",
			data:        []byte("invalid json"),
			password:    "TestPass123!",
			expectError: "解析加密資料失敗",
		},
		{
			name:        "空密碼測試",
			data:        []byte(`{"version":"1.0","algorithm":"aes256"}`),
			password:    "",
			expectError: "密碼不能為空",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := service.DecryptContent(tc.data, tc.password, AlgorithmAES256)
			
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

// TestValidatePassword 測試密碼強度驗證功能
// 驗證各種密碼強度要求的檢查
func TestValidatePassword(t *testing.T) {
	service := NewEncryptionService()
	
	testCases := []struct {
		name     string
		password string
		expected bool
		reason   string
	}{
		{
			name:     "有效密碼測試",
			password: "ValidPass123!",
			expected: true,
			reason:   "包含大小寫字母、數字和特殊字元",
		},
		{
			name:     "太短密碼測試",
			password: "Short1!",
			expected: false,
			reason:   "密碼長度不足 8 字元",
		},
		{
			name:     "太長密碼測試",
			password: strings.Repeat("VeryLongPassword123!", 10),
			expected: false,
			reason:   "密碼長度超過 128 字元",
		},
		{
			name:     "缺少大寫字母測試",
			password: "lowercase123!",
			expected: false,
			reason:   "缺少大寫字母",
		},
		{
			name:     "缺少小寫字母測試",
			password: "UPPERCASE123!",
			expected: false,
			reason:   "缺少小寫字母",
		},
		{
			name:     "缺少數字測試",
			password: "NoNumbers!",
			expected: false,
			reason:   "缺少數字",
		},
		{
			name:     "缺少特殊字元測試",
			password: "NoSpecialChars123",
			expected: false,
			reason:   "缺少特殊字元",
		},
		{
			name:     "中文密碼測試",
			password: "中文密碼123!Aa",
			expected: true,
			reason:   "包含中文字元但滿足其他要求",
		},
		{
			name:     "邊界長度測試（8字元）",
			password: "Valid1!a",
			expected: true,
			reason:   "剛好 8 字元且滿足所有要求",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := service.ValidatePassword(tc.password)
			
			if result != tc.expected {
				t.Errorf("密碼驗證結果不正確:\n密碼: %s\n期望: %v\n實際: %v\n原因: %s",
					tc.password, tc.expected, result, tc.reason)
			}
		})
	}
}

// TestEncryptDecrypt_CrossAlgorithm 測試跨演算法的加密解密
// 驗證不同演算法之間的相容性處理
func TestEncryptDecrypt_CrossAlgorithm(t *testing.T) {
	service := NewEncryptionService()
	
	content := "跨演算法測試內容"
	password := "CrossAlgo123!"
	
	// 使用 AES256 加密
	aesEncrypted, err := service.EncryptContent(content, password, AlgorithmAES256)
	if err != nil {
		t.Fatalf("AES256 加密失敗: %v", err)
	}
	
	// 嘗試使用 ChaCha20 解密（應該失敗）
	_, err = service.DecryptContent(aesEncrypted, password, AlgorithmChaCha20)
	if err == nil {
		t.Error("跨演算法解密應該失敗，但成功了")
	}
	
	// 使用正確的演算法解密（應該成功）
	decrypted, err := service.DecryptContent(aesEncrypted, password, AlgorithmAES256)
	if err != nil {
		t.Fatalf("正確演算法解密失敗: %v", err)
	}
	
	if decrypted != content {
		t.Errorf("解密結果不匹配: 期望 %s，實際 %s", content, decrypted)
	}
}

// TestBiometricAuth_Integration 測試生物識別驗證的整合
// 驗證生物識別功能的實際行為
func TestBiometricAuth_Integration(t *testing.T) {
	service := NewEncryptionService()
	noteID := "test-note-biometric"
	
	// 測試設定生物識別驗證
	err := service.SetupBiometricAuth(noteID)
	// 在 macOS 上可能成功也可能失敗（取決於硬體和用戶操作）
	// 在其他平台上應該失敗
	t.Logf("SetupBiometricAuth 結果: %v", err)
	
	// 測試生物識別驗證
	success, err := service.AuthenticateWithBiometric(noteID)
	// 結果取決於平台和設定狀態
	t.Logf("AuthenticateWithBiometric 結果: success=%v, error=%v", success, err)
	
	// 基本驗證：如果有錯誤，success 應該為 false
	if err != nil && success {
		t.Error("有錯誤時 success 不應該為 true")
	}
}

// BenchmarkEncryptContent_AES256 效能測試：AES-256 加密
// 測試 AES-256 加密的效能表現
func BenchmarkEncryptContent_AES256(b *testing.B) {
	service := NewEncryptionService()
	content := "這是效能測試內容"
	password := "BenchmarkPass123!"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.EncryptContent(content, password, AlgorithmAES256)
		if err != nil {
			b.Fatalf("加密失敗: %v", err)
		}
	}
}

// BenchmarkEncryptContent_ChaCha20 效能測試：ChaCha20-Poly1305 加密
// 測試 ChaCha20-Poly1305 加密的效能表現
func BenchmarkEncryptContent_ChaCha20(b *testing.B) {
	service := NewEncryptionService()
	content := "這是效能測試內容"
	password := "BenchmarkPass123!"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.EncryptContent(content, password, AlgorithmChaCha20)
		if err != nil {
			b.Fatalf("加密失敗: %v", err)
		}
	}
}

// BenchmarkDecryptContent_AES256 效能測試：AES-256 解密
// 測試 AES-256 解密的效能表現
func BenchmarkDecryptContent_AES256(b *testing.B) {
	service := NewEncryptionService()
	content := "這是效能測試內容"
	password := "BenchmarkPass123!"
	
	// 預先加密資料
	encrypted, err := service.EncryptContent(content, password, AlgorithmAES256)
	if err != nil {
		b.Fatalf("預先加密失敗: %v", err)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.DecryptContent(encrypted, password, AlgorithmAES256)
		if err != nil {
			b.Fatalf("解密失敗: %v", err)
		}
	}
}

// BenchmarkValidatePassword 效能測試：密碼驗證
// 測試密碼強度驗證的效能表現
func BenchmarkValidatePassword(b *testing.B) {
	service := NewEncryptionService()
	password := "BenchmarkPassword123!"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.ValidatePassword(password)
	}
}