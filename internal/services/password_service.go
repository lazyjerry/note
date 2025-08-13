// Package services 提供密碼驗證系統的具體實作
// 包含密碼雜湊、驗證、強度檢查和重試機制
package services

import (
	"crypto/rand"     // 安全隨機數產生器
	"crypto/sha256"   // SHA-256 雜湊演算法
	"crypto/subtle"   // 安全比較函數
	"encoding/base64" // Base64 編碼
	"errors"          // 錯誤處理
	"fmt"             // 格式化輸出
	"golang.org/x/crypto/pbkdf2" // PBKDF2 金鑰衍生函數
	"strings"         // 字串處理
	"sync"            // 同步原語
	"time"            // 時間處理
)

// 密碼驗證相關常數
const (
	PasswordSaltSize    = 32     // 密碼鹽值大小（位元組）
	PasswordHashSize    = 64     // 密碼雜湊大小（位元組）
	PasswordPBKDF2Rounds = 100000 // PBKDF2 迭代次數
	MaxRetryAttempts    = 3      // 最大重試次數
	RetryLockoutDuration = 5 * time.Minute // 重試鎖定時間
)

// PasswordHash 代表密碼雜湊資料結構
// 包含鹽值、雜湊值和相關元資料
type PasswordHash struct {
	Salt      string    `json:"salt"`       // Base64 編碼的鹽值
	Hash      string    `json:"hash"`       // Base64 編碼的密碼雜湊
	Algorithm string    `json:"algorithm"`  // 雜湊演算法（目前為 "pbkdf2-sha256"）
	Rounds    int       `json:"rounds"`     // PBKDF2 迭代次數
	CreatedAt time.Time `json:"created_at"` // 建立時間
}

// RetryInfo 代表密碼重試資訊
// 追蹤重試次數和鎖定狀態
type RetryInfo struct {
	Attempts    int       `json:"attempts"`     // 當前重試次數
	LastAttempt time.Time `json:"last_attempt"` // 最後嘗試時間
	LockedUntil time.Time `json:"locked_until"` // 鎖定到期時間
}

// PasswordService 定義密碼驗證服務的介面
// 提供密碼雜湊、驗證、重試管理等功能
type PasswordService interface {
	// HashPassword 將明文密碼轉換為安全的雜湊值
	// 參數：password（明文密碼）
	// 回傳：密碼雜湊結構和可能的錯誤
	HashPassword(password string) (*PasswordHash, error)
	
	// VerifyPassword 驗證明文密碼是否與雜湊值匹配
	// 參數：password（明文密碼）、hash（密碼雜湊結構）
	// 回傳：驗證結果和可能的錯誤
	VerifyPassword(password string, hash *PasswordHash) (bool, error)
	
	// CheckPasswordStrength 檢查密碼強度
	// 參數：password（要檢查的密碼）
	// 回傳：強度等級和建議
	CheckPasswordStrength(password string) (PasswordStrength, []string)
	
	// RecordFailedAttempt 記錄失敗的密碼嘗試
	// 參數：identifier（識別符，如用戶 ID 或檔案路徑）
	// 回傳：可能的錯誤
	RecordFailedAttempt(identifier string) error
	
	// IsLocked 檢查是否因重試次數過多而被鎖定
	// 參數：identifier（識別符）
	// 回傳：是否被鎖定和剩餘鎖定時間
	IsLocked(identifier string) (bool, time.Duration)
	
	// ResetRetryCount 重置重試計數（成功驗證後調用）
	// 參數：identifier（識別符）
	ResetRetryCount(identifier string)
	
	// GetRetryInfo 取得重試資訊
	// 參數：identifier（識別符）
	// 回傳：重試資訊
	GetRetryInfo(identifier string) *RetryInfo
}

// PasswordStrength 代表密碼強度等級
type PasswordStrength int

const (
	PasswordWeak   PasswordStrength = iota // 弱密碼
	PasswordFair                           // 一般密碼
	PasswordGood                           // 良好密碼
	PasswordStrong                         // 強密碼
)

// String 回傳密碼強度的字串表示
func (ps PasswordStrength) String() string {
	switch ps {
	case PasswordWeak:
		return "弱"
	case PasswordFair:
		return "一般"
	case PasswordGood:
		return "良好"
	case PasswordStrong:
		return "強"
	default:
		return "未知"
	}
}

// passwordService 實作 PasswordService 介面
// 提供完整的密碼管理功能
type passwordService struct {
	retryMap map[string]*RetryInfo // 重試資訊映射表
	mutex    sync.RWMutex          // 讀寫鎖保護並發存取
}

// NewPasswordService 建立新的密碼服務實例
// 回傳：PasswordService 介面實例
//
// 執行流程：
// 1. 建立 passwordService 結構體實例
// 2. 初始化重試資訊映射表
// 3. 回傳服務介面
func NewPasswordService() PasswordService {
	return &passwordService{
		retryMap: make(map[string]*RetryInfo),
		mutex:    sync.RWMutex{},
	}
}

// HashPassword 將明文密碼轉換為安全的雜湊值
// 參數：password（明文密碼）
// 回傳：密碼雜湊結構和可能的錯誤
//
// 執行流程：
// 1. 驗證密碼不為空
// 2. 產生隨機鹽值
// 3. 使用 PBKDF2-SHA256 計算密碼雜湊
// 4. 建立並回傳密碼雜湊結構
func (ps *passwordService) HashPassword(password string) (*PasswordHash, error) {
	if password == "" {
		return nil, errors.New("密碼不能為空")
	}
	
	// 產生隨機鹽值
	salt := make([]byte, PasswordSaltSize)
	if _, err := rand.Read(salt); err != nil {
		return nil, fmt.Errorf("產生鹽值失敗: %w", err)
	}
	
	// 使用 PBKDF2-SHA256 計算密碼雜湊
	hash := pbkdf2.Key([]byte(password), salt, PasswordPBKDF2Rounds, PasswordHashSize, sha256.New)
	
	// 建立密碼雜湊結構
	passwordHash := &PasswordHash{
		Salt:      base64.StdEncoding.EncodeToString(salt),
		Hash:      base64.StdEncoding.EncodeToString(hash),
		Algorithm: "pbkdf2-sha256",
		Rounds:    PasswordPBKDF2Rounds,
		CreatedAt: time.Now(),
	}
	
	return passwordHash, nil
}

// VerifyPassword 驗證明文密碼是否與雜湊值匹配
// 參數：password（明文密碼）、hash（密碼雜湊結構）
// 回傳：驗證結果和可能的錯誤
//
// 執行流程：
// 1. 驗證輸入參數
// 2. 解碼鹽值和雜湊值
// 3. 使用相同參數重新計算雜湊
// 4. 使用安全比較函數比較雜湊值
// 5. 回傳比較結果
func (ps *passwordService) VerifyPassword(password string, hash *PasswordHash) (bool, error) {
	if password == "" {
		return false, errors.New("密碼不能為空")
	}
	
	if hash == nil {
		return false, errors.New("密碼雜湊不能為空")
	}
	
	// 驗證演算法支援
	if hash.Algorithm != "pbkdf2-sha256" {
		return false, fmt.Errorf("不支援的雜湊演算法: %s", hash.Algorithm)
	}
	
	// 解碼鹽值
	salt, err := base64.StdEncoding.DecodeString(hash.Salt)
	if err != nil {
		return false, fmt.Errorf("解碼鹽值失敗: %w", err)
	}
	
	// 解碼預期雜湊值
	expectedHash, err := base64.StdEncoding.DecodeString(hash.Hash)
	if err != nil {
		return false, fmt.Errorf("解碼雜湊值失敗: %w", err)
	}
	
	// 使用相同參數重新計算雜湊
	computedHash := pbkdf2.Key([]byte(password), salt, hash.Rounds, len(expectedHash), sha256.New)
	
	// 使用安全比較函數比較雜湊值（防止時序攻擊）
	return subtle.ConstantTimeCompare(expectedHash, computedHash) == 1, nil
}

// CheckPasswordStrength 檢查密碼強度
// 參數：password（要檢查的密碼）
// 回傳：強度等級和改進建議
//
// 強度評估標準：
// - 弱：不滿足基本要求
// - 一般：滿足基本要求但缺少某些特徵
// - 良好：滿足大部分要求
// - 強：滿足所有要求且長度充足
//
// 執行流程：
// 1. 檢查密碼長度
// 2. 分析字元類型分佈
// 3. 檢查常見弱密碼模式
// 4. 計算強度分數
// 5. 生成改進建議
func (ps *passwordService) CheckPasswordStrength(password string) (PasswordStrength, []string) {
	var suggestions []string
	score := 0
	
	// 檢查密碼長度
	length := len(password)
	if length < MinPasswordLength {
		suggestions = append(suggestions, fmt.Sprintf("密碼長度至少需要 %d 個字元", MinPasswordLength))
		return PasswordWeak, suggestions
	}
	
	if length >= 8 {
		score += 1
	}
	if length >= 12 {
		score += 1
	}
	if length >= 16 {
		score += 1
	}
	
	// 檢查字元類型
	var (
		hasLower   = false
		hasUpper   = false
		hasNumber  = false
		hasSpecial = false
	)
	
	for _, char := range password {
		switch {
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= '0' && char <= '9':
			hasNumber = true
		default:
			hasSpecial = true
		}
	}
	
	// 根據字元類型給分和建議
	if hasLower {
		score += 1
	} else {
		suggestions = append(suggestions, "建議包含小寫字母")
	}
	
	if hasUpper {
		score += 1
	} else {
		suggestions = append(suggestions, "建議包含大寫字母")
	}
	
	if hasNumber {
		score += 1
	} else {
		suggestions = append(suggestions, "建議包含數字")
	}
	
	if hasSpecial {
		score += 1
	} else {
		suggestions = append(suggestions, "建議包含特殊字元 (!@#$%^&* 等)")
	}
	
	// 檢查常見弱密碼模式
	if ps.isCommonWeakPassword(password) {
		suggestions = append(suggestions, "避免使用常見的弱密碼")
		score -= 2
	}
	
	// 檢查重複字元
	if ps.hasRepeatingChars(password) {
		suggestions = append(suggestions, "避免使用重複字元")
		score -= 1
	}
	
	// 檢查是否滿足所有基本要求
	hasAllRequiredTypes := hasLower && hasUpper && hasNumber && hasSpecial
	
	// 根據分數和基本要求判定強度等級
	switch {
	case score <= 2:
		return PasswordWeak, suggestions
	case score <= 4 || !hasAllRequiredTypes:
		return PasswordFair, suggestions
	case score <= 6:
		return PasswordGood, suggestions
	default:
		if len(suggestions) == 0 {
			suggestions = append(suggestions, "密碼強度良好")
		}
		return PasswordStrong, suggestions
	}
}

// RecordFailedAttempt 記錄失敗的密碼嘗試
// 參數：identifier（識別符，如用戶 ID 或檔案路徑）
// 回傳：可能的錯誤
//
// 執行流程：
// 1. 取得或建立重試資訊
// 2. 增加失敗次數
// 3. 更新最後嘗試時間
// 4. 檢查是否需要鎖定
// 5. 更新鎖定狀態
func (ps *passwordService) RecordFailedAttempt(identifier string) error {
	if identifier == "" {
		return errors.New("識別符不能為空")
	}
	
	ps.mutex.Lock()
	defer ps.mutex.Unlock()
	
	// 取得或建立重試資訊
	retryInfo, exists := ps.retryMap[identifier]
	if !exists {
		retryInfo = &RetryInfo{
			Attempts:    0,
			LastAttempt: time.Time{},
			LockedUntil: time.Time{},
		}
		ps.retryMap[identifier] = retryInfo
	}
	
	// 檢查是否仍在鎖定期間
	now := time.Now()
	if now.Before(retryInfo.LockedUntil) {
		return fmt.Errorf("帳戶已鎖定，請在 %v 後重試", retryInfo.LockedUntil.Sub(now).Round(time.Second))
	}
	
	// 增加失敗次數
	retryInfo.Attempts++
	retryInfo.LastAttempt = now
	
	// 檢查是否達到最大重試次數
	if retryInfo.Attempts >= MaxRetryAttempts {
		retryInfo.LockedUntil = now.Add(RetryLockoutDuration)
		return fmt.Errorf("密碼錯誤次數過多，帳戶已鎖定 %v", RetryLockoutDuration)
	}
	
	return nil
}

// IsLocked 檢查是否因重試次數過多而被鎖定
// 參數：identifier（識別符）
// 回傳：是否被鎖定和剩餘鎖定時間
//
// 執行流程：
// 1. 查找重試資訊
// 2. 檢查鎖定到期時間
// 3. 計算剩餘鎖定時間
// 4. 回傳鎖定狀態
func (ps *passwordService) IsLocked(identifier string) (bool, time.Duration) {
	ps.mutex.RLock()
	defer ps.mutex.RUnlock()
	
	retryInfo, exists := ps.retryMap[identifier]
	if !exists {
		return false, 0
	}
	
	now := time.Now()
	if now.Before(retryInfo.LockedUntil) {
		return true, retryInfo.LockedUntil.Sub(now)
	}
	
	return false, 0
}

// ResetRetryCount 重置重試計數（成功驗證後調用）
// 參數：identifier（識別符）
//
// 執行流程：
// 1. 查找重試資訊
// 2. 重置失敗次數
// 3. 清除鎖定狀態
func (ps *passwordService) ResetRetryCount(identifier string) {
	ps.mutex.Lock()
	defer ps.mutex.Unlock()
	
	retryInfo, exists := ps.retryMap[identifier]
	if exists {
		retryInfo.Attempts = 0
		retryInfo.LockedUntil = time.Time{}
	}
}

// GetRetryInfo 取得重試資訊
// 參數：identifier（識別符）
// 回傳：重試資訊（如果不存在則回傳 nil）
func (ps *passwordService) GetRetryInfo(identifier string) *RetryInfo {
	ps.mutex.RLock()
	defer ps.mutex.RUnlock()
	
	retryInfo, exists := ps.retryMap[identifier]
	if !exists {
		return nil
	}
	
	// 回傳副本以避免外部修改
	return &RetryInfo{
		Attempts:    retryInfo.Attempts,
		LastAttempt: retryInfo.LastAttempt,
		LockedUntil: retryInfo.LockedUntil,
	}
}

// isCommonWeakPassword 檢查是否為常見的弱密碼
// 參數：password（要檢查的密碼）
// 回傳：是否為常見弱密碼
func (ps *passwordService) isCommonWeakPassword(password string) bool {
	// 常見弱密碼清單（可以擴展）
	commonWeakPasswords := []string{
		"password", "123456", "123456789", "qwerty", "abc123",
		"password123", "admin", "root", "user", "guest",
		"12345678", "1234567890", "qwerty123", "password1",
		"123123", "111111", "000000", "1qaz2wsx",
	}
	
	lowerPassword := strings.ToLower(password)
	for _, weak := range commonWeakPasswords {
		if lowerPassword == weak {
			return true
		}
	}
	
	return false
}

// hasRepeatingChars 檢查是否有過多重複字元
// 參數：password（要檢查的密碼）
// 回傳：是否有重複字元問題
func (ps *passwordService) hasRepeatingChars(password string) bool {
	if len(password) < 3 {
		return false
	}
	
	// 檢查連續重複字元（如 "aaa", "111"）
	repeatCount := 1
	for i := 1; i < len(password); i++ {
		if password[i] == password[i-1] {
			repeatCount++
			if repeatCount >= 3 {
				return true
			}
		} else {
			repeatCount = 1
		}
	}
	
	return false
}