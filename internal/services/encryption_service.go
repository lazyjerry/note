// Package services 提供加密服務的具體實作
// 支援 AES-256 和 ChaCha20-Poly1305 兩種加密演算法
// 提供完整的加密、解密、密碼驗證和生物識別驗證功能
package services

import (
	"crypto/aes"           // AES 加密演算法
	"crypto/cipher"        // 加密模式介面
	"crypto/rand"          // 安全隨機數產生器
	"crypto/sha256"        // SHA-256 雜湊演算法
	"encoding/base64"      // Base64 編碼
	"encoding/json"        // JSON 序列化
	"errors"               // 錯誤處理
	"fmt"                  // 格式化輸出
	"golang.org/x/crypto/chacha20poly1305" // ChaCha20-Poly1305 加密演算法
	"golang.org/x/crypto/pbkdf2"           // PBKDF2 金鑰衍生函數
	"io"                   // 輸入輸出介面
	"unicode"              // Unicode 字元處理
)

// 支援的加密演算法常數
const (
	AlgorithmAES256    = "aes256"    // AES-256-GCM 加密演算法
	AlgorithmChaCha20  = "chacha20"  // ChaCha20-Poly1305 加密演算法
)

// 加密參數常數
const (
	SaltSize     = 32  // 鹽值大小（位元組）
	NonceSize    = 12  // Nonce 大小（位元組）
	KeySize      = 32  // 金鑰大小（位元組）
	PBKDF2Rounds = 100000 // PBKDF2 迭代次數
)

// 密碼強度要求常數
const (
	MinPasswordLength = 8   // 最小密碼長度
	MaxPasswordLength = 128 // 最大密碼長度
)

// EncryptedData 代表加密後的資料結構
// 包含加密演算法、鹽值、隨機數和加密內容等資訊
type EncryptedData struct {
	Version   string `json:"version"`   // 加密格式版本
	Algorithm string `json:"algorithm"` // 使用的加密演算法
	Salt      string `json:"salt"`      // Base64 編碼的鹽值
	Nonce     string `json:"nonce"`     // Base64 編碼的隨機數
	Data      string `json:"data"`      // Base64 編碼的加密內容
	Checksum  string `json:"checksum"`  // SHA-256 校驗和
}

// encryptionService 實作 EncryptionService 介面
// 提供完整的加密解密功能和密碼管理
type encryptionService struct {
	// 可以在這裡添加配置選項或依賴注入
}

// NewEncryptionService 建立新的加密服務實例
// 回傳：EncryptionService 介面實例
//
// 執行流程：
// 1. 建立 encryptionService 結構體實例
// 2. 初始化必要的配置（如果需要）
// 3. 回傳服務介面
func NewEncryptionService() EncryptionService {
	return &encryptionService{}
}

// EncryptContent 使用指定演算法和密碼加密內容
// 參數：
//   - content: 要加密的明文內容
//   - password: 用於加密的密碼
//   - algorithm: 加密演算法（"aes256" 或 "chacha20"）
// 回傳：加密後的位元組陣列和可能的錯誤
//
// 執行流程：
// 1. 驗證輸入參數的有效性
// 2. 產生隨機鹽值和隨機數
// 3. 使用 PBKDF2 從密碼衍生金鑰
// 4. 根據指定演算法進行加密
// 5. 建立加密資料結構並序列化為 JSON
// 6. 回傳序列化後的位元組陣列
func (s *encryptionService) EncryptContent(content, password string, algorithm string) ([]byte, error) {
	// 驗證輸入參數
	if content == "" {
		return nil, errors.New("加密內容不能為空")
	}
	if password == "" {
		return nil, errors.New("密碼不能為空")
	}
	if !s.isValidAlgorithm(algorithm) {
		return nil, fmt.Errorf("不支援的加密演算法: %s", algorithm)
	}

	// 產生隨機鹽值
	salt := make([]byte, SaltSize)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, fmt.Errorf("產生鹽值失敗: %w", err)
	}

	// 使用 PBKDF2 從密碼衍生金鑰
	key := pbkdf2.Key([]byte(password), salt, PBKDF2Rounds, KeySize, sha256.New)

	var encryptedContent []byte
	var nonce []byte
	var err error

	// 根據演算法進行加密
	switch algorithm {
	case AlgorithmAES256:
		encryptedContent, nonce, err = s.encryptWithAES([]byte(content), key)
	case AlgorithmChaCha20:
		encryptedContent, nonce, err = s.encryptWithChaCha20([]byte(content), key)
	default:
		return nil, fmt.Errorf("不支援的加密演算法: %s", algorithm)
	}

	if err != nil {
		return nil, fmt.Errorf("加密失敗: %w", err)
	}

	// 計算校驗和
	checksum := sha256.Sum256(encryptedContent)

	// 建立加密資料結構
	encData := EncryptedData{
		Version:   "1.0",
		Algorithm: algorithm,
		Salt:      base64.StdEncoding.EncodeToString(salt),
		Nonce:     base64.StdEncoding.EncodeToString(nonce),
		Data:      base64.StdEncoding.EncodeToString(encryptedContent),
		Checksum:  base64.StdEncoding.EncodeToString(checksum[:]),
	}

	// 序列化為 JSON
	return json.Marshal(encData)
}

// DecryptContent 使用指定演算法和密碼解密內容
// 參數：
//   - encryptedData: 加密的資料位元組陣列
//   - password: 用於解密的密碼
//   - algorithm: 加密演算法（"aes256" 或 "chacha20"）
// 回傳：解密後的內容字串和可能的錯誤
//
// 執行流程：
// 1. 反序列化加密資料結構
// 2. 驗證加密格式和演算法
// 3. 解碼 Base64 編碼的資料
// 4. 使用 PBKDF2 從密碼衍生金鑰
// 5. 根據演算法進行解密
// 6. 驗證校驗和確保資料完整性
// 7. 回傳解密後的明文內容
func (s *encryptionService) DecryptContent(encryptedData []byte, password string, algorithm string) (string, error) {
	// 驗證輸入參數
	if len(encryptedData) == 0 {
		return "", errors.New("加密資料不能為空")
	}
	if password == "" {
		return "", errors.New("密碼不能為空")
	}

	// 反序列化加密資料
	var encData EncryptedData
	if err := json.Unmarshal(encryptedData, &encData); err != nil {
		return "", fmt.Errorf("解析加密資料失敗: %w", err)
	}

	// 驗證加密格式版本
	if encData.Version != "1.0" {
		return "", fmt.Errorf("不支援的加密格式版本: %s", encData.Version)
	}

	// 驗證演算法一致性
	if algorithm != "" && encData.Algorithm != algorithm {
		return "", fmt.Errorf("演算法不匹配: 期望 %s，實際 %s", algorithm, encData.Algorithm)
	}

	// 解碼 Base64 資料
	salt, err := base64.StdEncoding.DecodeString(encData.Salt)
	if err != nil {
		return "", fmt.Errorf("解碼鹽值失敗: %w", err)
	}

	nonce, err := base64.StdEncoding.DecodeString(encData.Nonce)
	if err != nil {
		return "", fmt.Errorf("解碼隨機數失敗: %w", err)
	}

	ciphertext, err := base64.StdEncoding.DecodeString(encData.Data)
	if err != nil {
		return "", fmt.Errorf("解碼加密內容失敗: %w", err)
	}

	expectedChecksum, err := base64.StdEncoding.DecodeString(encData.Checksum)
	if err != nil {
		return "", fmt.Errorf("解碼校驗和失敗: %w", err)
	}

	// 驗證校驗和
	actualChecksum := sha256.Sum256(ciphertext)
	if string(expectedChecksum) != string(actualChecksum[:]) {
		return "", errors.New("資料校驗和不匹配，可能已被篡改")
	}

	// 使用 PBKDF2 從密碼衍生金鑰
	key := pbkdf2.Key([]byte(password), salt, PBKDF2Rounds, KeySize, sha256.New)

	var plaintext []byte

	// 根據演算法進行解密
	switch encData.Algorithm {
	case AlgorithmAES256:
		plaintext, err = s.decryptWithAES(ciphertext, key, nonce)
	case AlgorithmChaCha20:
		plaintext, err = s.decryptWithChaCha20(ciphertext, key, nonce)
	default:
		return "", fmt.Errorf("不支援的加密演算法: %s", encData.Algorithm)
	}

	if err != nil {
		return "", fmt.Errorf("解密失敗: %w", err)
	}

	return string(plaintext), nil
}

// encryptWithAES 使用 AES-256-GCM 模式加密資料
// 參數：
//   - plaintext: 要加密的明文資料
//   - key: 32 位元組的加密金鑰
// 回傳：加密後的資料、隨機數和可能的錯誤
//
// 執行流程：
// 1. 建立 AES 加密器
// 2. 建立 GCM 模式包裝器
// 3. 產生隨機 nonce
// 4. 使用 GCM 模式加密資料
// 5. 回傳加密結果和 nonce
func (s *encryptionService) encryptWithAES(plaintext, key []byte) ([]byte, []byte, error) {
	// 建立 AES 加密器
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, fmt.Errorf("建立 AES 加密器失敗: %w", err)
	}

	// 建立 GCM 模式
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, fmt.Errorf("建立 GCM 模式失敗: %w", err)
	}

	// 產生隨機 nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, fmt.Errorf("產生 nonce 失敗: %w", err)
	}

	// 加密資料
	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)

	return ciphertext, nonce, nil
}

// decryptWithAES 使用 AES-256-GCM 模式解密資料
// 參數：
//   - ciphertext: 要解密的加密資料
//   - key: 32 位元組的解密金鑰
//   - nonce: 用於解密的隨機數
// 回傳：解密後的明文資料和可能的錯誤
//
// 執行流程：
// 1. 建立 AES 解密器
// 2. 建立 GCM 模式包裝器
// 3. 使用 GCM 模式解密資料
// 4. 回傳解密結果
func (s *encryptionService) decryptWithAES(ciphertext, key, nonce []byte) ([]byte, error) {
	// 建立 AES 解密器
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("建立 AES 解密器失敗: %w", err)
	}

	// 建立 GCM 模式
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("建立 GCM 模式失敗: %w", err)
	}

	// 解密資料
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("AES 解密失敗: %w", err)
	}

	return plaintext, nil
}

// encryptWithChaCha20 使用 ChaCha20-Poly1305 加密資料
// 參數：
//   - plaintext: 要加密的明文資料
//   - key: 32 位元組的加密金鑰
// 回傳：加密後的資料、隨機數和可能的錯誤
//
// 執行流程：
// 1. 建立 ChaCha20-Poly1305 加密器
// 2. 產生隨機 nonce
// 3. 使用 ChaCha20-Poly1305 加密資料
// 4. 回傳加密結果和 nonce
func (s *encryptionService) encryptWithChaCha20(plaintext, key []byte) ([]byte, []byte, error) {
	// 建立 ChaCha20-Poly1305 加密器
	aead, err := chacha20poly1305.New(key)
	if err != nil {
		return nil, nil, fmt.Errorf("建立 ChaCha20-Poly1305 加密器失敗: %w", err)
	}

	// 產生隨機 nonce
	nonce := make([]byte, aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, fmt.Errorf("產生 nonce 失敗: %w", err)
	}

	// 加密資料
	ciphertext := aead.Seal(nil, nonce, plaintext, nil)

	return ciphertext, nonce, nil
}

// decryptWithChaCha20 使用 ChaCha20-Poly1305 解密資料
// 參數：
//   - ciphertext: 要解密的加密資料
//   - key: 32 位元組的解密金鑰
//   - nonce: 用於解密的隨機數
// 回傳：解密後的明文資料和可能的錯誤
//
// 執行流程：
// 1. 建立 ChaCha20-Poly1305 解密器
// 2. 使用 ChaCha20-Poly1305 解密資料
// 3. 回傳解密結果
func (s *encryptionService) decryptWithChaCha20(ciphertext, key, nonce []byte) ([]byte, error) {
	// 建立 ChaCha20-Poly1305 解密器
	aead, err := chacha20poly1305.New(key)
	if err != nil {
		return nil, fmt.Errorf("建立 ChaCha20-Poly1305 解密器失敗: %w", err)
	}

	// 解密資料
	plaintext, err := aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("ChaCha20-Poly1305 解密失敗: %w", err)
	}

	return plaintext, nil
}

// ValidatePassword 驗證密碼強度是否符合要求
// 參數：password（要驗證的密碼）
// 回傳：密碼是否有效
//
// 密碼強度要求：
// 1. 長度在 8-128 字元之間
// 2. 至少包含一個大寫字母
// 3. 至少包含一個小寫字母
// 4. 至少包含一個數字
// 5. 至少包含一個特殊字元
//
// 執行流程：
// 1. 檢查密碼長度
// 2. 檢查是否包含大寫字母
// 3. 檢查是否包含小寫字母
// 4. 檢查是否包含數字
// 5. 檢查是否包含特殊字元
// 6. 回傳驗證結果
func (s *encryptionService) ValidatePassword(password string) bool {
	// 檢查密碼長度
	if len(password) < MinPasswordLength || len(password) > MaxPasswordLength {
		return false
	}

	var (
		hasUpper   = false // 是否包含大寫字母
		hasLower   = false // 是否包含小寫字母
		hasNumber  = false // 是否包含數字
		hasSpecial = false // 是否包含特殊字元
	)

	// 遍歷密碼中的每個字元
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	// 檢查是否滿足所有強度要求
	return hasUpper && hasLower && hasNumber && hasSpecial
}

// isValidAlgorithm 檢查加密演算法是否受支援
// 參數：algorithm（演算法名稱）
// 回傳：是否為有效的演算法
func (s *encryptionService) isValidAlgorithm(algorithm string) bool {
	return algorithm == AlgorithmAES256 || algorithm == AlgorithmChaCha20
}

// SetupBiometricAuth 為指定筆記設定生物識別驗證
// 參數：noteID（筆記 ID）
// 回傳：可能的錯誤
//
// 執行流程：
// 1. 建立生物識別服務實例
// 2. 調用設定方法
// 3. 回傳結果
func (s *encryptionService) SetupBiometricAuth(noteID string) error {
	biometricService := NewBiometricService()
	return biometricService.SetupForNote(noteID)
}

// AuthenticateWithBiometric 使用生物識別進行驗證
// 參數：noteID（筆記 ID）
// 回傳：驗證結果（成功/失敗）和可能的錯誤
//
// 執行流程：
// 1. 建立生物識別服務實例
// 2. 執行驗證
// 3. 回傳驗證結果
func (s *encryptionService) AuthenticateWithBiometric(noteID string) (bool, error) {
	biometricService := NewBiometricService()
	result := biometricService.AuthenticateForNote(noteID, "")
	
	if result.Error != nil && !result.Cancelled {
		return false, result.Error
	}
	
	return result.Success, nil
}