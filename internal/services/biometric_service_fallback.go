// Package services 提供非 macOS 系統的生物識別驗證回退實作
// 在不支援的平台上提供基本的介面實作
// +build !darwin

package services

import (
	"errors"
	"sync"
	"time"
)

// BiometricType 代表生物識別類型
type BiometricType int

const (
	BiometricTypeNone    BiometricType = iota // 不支援生物識別
	BiometricTypeTouchID                      // Touch ID
	BiometricTypeFaceID                       // Face ID
)

// String 回傳生物識別類型的字串表示
func (bt BiometricType) String() string {
	switch bt {
	case BiometricTypeNone:
		return "無"
	case BiometricTypeTouchID:
		return "Touch ID"
	case BiometricTypeFaceID:
		return "Face ID"
	default:
		return "未知"
	}
}

// BiometricResult 代表生物識別驗證結果
type BiometricResult struct {
	Success   bool          `json:"success"`    // 驗證是否成功
	Cancelled bool          `json:"cancelled"`  // 用戶是否取消
	Error     error         `json:"error"`      // 錯誤資訊（如果有）
	Duration  time.Duration `json:"duration"`   // 驗證耗時
}

// BiometricService 定義生物識別驗證服務的介面
// 提供 macOS Touch ID/Face ID 驗證功能
type BiometricService interface {
	// IsAvailable 檢查生物識別驗證是否可用
	// 回傳：是否可用和生物識別類型
	IsAvailable() (bool, BiometricType)
	
	// Authenticate 執行生物識別驗證
	// 參數：reason（驗證原因，顯示給用戶）
	// 回傳：驗證結果
	Authenticate(reason string) *BiometricResult
	
	// AuthenticateForNote 為特定筆記執行生物識別驗證
	// 參數：noteID（筆記 ID）、reason（驗證原因）
	// 回傳：驗證結果
	AuthenticateForNote(noteID, reason string) *BiometricResult
	
	// SetupForNote 為特定筆記設定生物識別驗證
	// 參數：noteID（筆記 ID）
	// 回傳：可能的錯誤
	SetupForNote(noteID string) error
	
	// RemoveForNote 移除特定筆記的生物識別驗證設定
	// 參數：noteID（筆記 ID）
	// 回傳：可能的錯誤
	RemoveForNote(noteID string) error
	
	// IsEnabledForNote 檢查特定筆記是否啟用生物識別驗證
	// 參數：noteID（筆記 ID）
	// 回傳：是否啟用
	IsEnabledForNote(noteID string) bool
}

// biometricService 實作 BiometricService 介面的回退版本
// 在非 macOS 系統上提供基本功能
type biometricService struct {
	enabledNotes map[string]bool // 啟用生物識別的筆記映射表
	mutex        sync.RWMutex    // 讀寫鎖保護並發存取
}

// NewBiometricService 建立新的生物識別服務實例
// 回傳：BiometricService 介面實例
//
// 注意：在非 macOS 系統上，此服務將始終回傳不可用狀態
func NewBiometricService() BiometricService {
	return &biometricService{
		enabledNotes: make(map[string]bool),
		mutex:        sync.RWMutex{},
	}
}

// IsAvailable 檢查生物識別驗證是否可用
// 回傳：是否可用和生物識別類型
//
// 注意：在非 macOS 系統上始終回傳不可用
func (bs *biometricService) IsAvailable() (bool, BiometricType) {
	return false, BiometricTypeNone
}

// Authenticate 執行生物識別驗證
// 參數：reason（驗證原因，顯示給用戶）
// 回傳：驗證結果
//
// 注意：在非 macOS 系統上始終回傳不支援錯誤
func (bs *biometricService) Authenticate(reason string) *BiometricResult {
	return &BiometricResult{
		Success:   false,
		Cancelled: false,
		Error:     errors.New("此平台不支援生物識別驗證"),
		Duration:  0,
	}
}

// AuthenticateForNote 為特定筆記執行生物識別驗證
// 參數：noteID（筆記 ID）、reason（驗證原因）
// 回傳：驗證結果
//
// 注意：在非 macOS 系統上始終回傳不支援錯誤
func (bs *biometricService) AuthenticateForNote(noteID, reason string) *BiometricResult {
	return &BiometricResult{
		Success:   false,
		Cancelled: false,
		Error:     errors.New("此平台不支援生物識別驗證"),
		Duration:  0,
	}
}

// SetupForNote 為特定筆記設定生物識別驗證
// 參數：noteID（筆記 ID）
// 回傳：可能的錯誤
//
// 注意：在非 macOS 系統上始終回傳不支援錯誤
func (bs *biometricService) SetupForNote(noteID string) error {
	return errors.New("此平台不支援生物識別驗證")
}

// RemoveForNote 移除特定筆記的生物識別驗證設定
// 參數：noteID（筆記 ID）
// 回傳：可能的錯誤
//
// 注意：在非 macOS 系統上始終回傳不支援錯誤
func (bs *biometricService) RemoveForNote(noteID string) error {
	return errors.New("此平台不支援生物識別驗證")
}

// IsEnabledForNote 檢查特定筆記是否啟用生物識別驗證
// 參數：noteID（筆記 ID）
// 回傳：是否啟用
//
// 注意：在非 macOS 系統上始終回傳 false
func (bs *biometricService) IsEnabledForNote(noteID string) bool {
	return false
}

