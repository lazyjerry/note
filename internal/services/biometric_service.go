// Package services 提供 macOS 生物識別驗證的具體實作
// 使用 CGO 調用 LocalAuthentication API 實現 Touch ID/Face ID 驗證
// +build darwin

package services

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework LocalAuthentication -framework Foundation
#import <LocalAuthentication/LocalAuthentication.h>
#import <Foundation/Foundation.h>

// 檢查生物識別驗證是否可用
int checkBiometricAvailability() {
    LAContext *context = [[LAContext alloc] init];
    NSError *error = nil;
    
    BOOL canEvaluate = [context canEvaluatePolicy:LAPolicyDeviceOwnerAuthenticationWithBiometrics error:&error];
    
    if (canEvaluate) {
        return 1; // 可用
    } else {
        return 0; // 不可用
    }
}

// 取得生物識別類型
int getBiometricType() {
    LAContext *context = [[LAContext alloc] init];
    NSError *error = nil;
    
    if ([context canEvaluatePolicy:LAPolicyDeviceOwnerAuthenticationWithBiometrics error:&error]) {
        if (@available(macOS 10.13.2, *)) {
            switch (context.biometryType) {
                case LABiometryTypeTouchID:
                    return 1; // Touch ID
                case LABiometryTypeFaceID:
                    return 2; // Face ID
                default:
                    return 0; // 未知或不支援
            }
        } else {
            return 1; // 舊版本預設為 Touch ID
        }
    }
    
    return 0; // 不可用
}

// 執行生物識別驗證
// 回傳值：1=成功, 0=失敗, -1=用戶取消, -2=系統錯誤
int performBiometricAuthentication(const char* reason) {
    LAContext *context = [[LAContext alloc] init];
    NSString *reasonString = [NSString stringWithUTF8String:reason];
    
    // 設定回退標題
    context.localizedFallbackTitle = @"使用密碼";
    
    __block int result = -2; // 預設為系統錯誤
    __block BOOL finished = NO;
    
    [context evaluatePolicy:LAPolicyDeviceOwnerAuthenticationWithBiometrics
            localizedReason:reasonString
                      reply:^(BOOL success, NSError *error) {
        if (success) {
            result = 1; // 成功
        } else {
            if (error) {
                switch (error.code) {
                    case LAErrorUserCancel:
                    case LAErrorAppCancel:
                    case LAErrorSystemCancel:
                        result = -1; // 用戶取消
                        break;
                    case LAErrorUserFallback:
                        result = 0; // 用戶選擇回退到密碼
                        break;
                    case LAErrorBiometryNotAvailable:
                    case LAErrorBiometryNotEnrolled:
                    case LAErrorBiometryLockout:
                        result = 0; // 生物識別不可用
                        break;
                    default:
                        result = -2; // 其他系統錯誤
                        break;
                }
            } else {
                result = 0; // 驗證失敗
            }
        }
        finished = YES;
    }];
    
    // 等待驗證完成（最多等待 30 秒）
    NSDate *timeout = [NSDate dateWithTimeIntervalSinceNow:30.0];
    while (!finished && [[NSDate date] compare:timeout] == NSOrderedAscending) {
        [[NSRunLoop currentRunLoop] runMode:NSDefaultRunLoopMode beforeDate:[NSDate dateWithTimeIntervalSinceNow:0.1]];
    }
    
    if (!finished) {
        result = -2; // 超時
    }
    
    return result;
}

// 取得錯誤描述
const char* getBiometricErrorDescription(int errorCode) {
    switch (errorCode) {
        case -1:
            return "用戶取消驗證";
        case -2:
            return "系統錯誤或超時";
        case 0:
            return "驗證失敗或生物識別不可用";
        default:
            return "未知錯誤";
    }
}
*/
import "C"

import (
	"errors"
	"fmt"
	"sync"
	"time"
	"unsafe"
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

// biometricService 實作 BiometricService 介面
// 提供完整的 macOS 生物識別驗證功能
type biometricService struct {
	enabledNotes map[string]bool // 啟用生物識別的筆記映射表
	mutex        sync.RWMutex    // 讀寫鎖保護並發存取
}

// NewBiometricService 建立新的生物識別服務實例
// 回傳：BiometricService 介面實例
//
// 執行流程：
// 1. 建立 biometricService 結構體實例
// 2. 初始化啟用筆記映射表
// 3. 回傳服務介面
func NewBiometricService() BiometricService {
	return &biometricService{
		enabledNotes: make(map[string]bool),
		mutex:        sync.RWMutex{},
	}
}

// IsAvailable 檢查生物識別驗證是否可用
// 回傳：是否可用和生物識別類型
//
// 執行流程：
// 1. 調用 C 函數檢查可用性
// 2. 調用 C 函數取得生物識別類型
// 3. 回傳結果
func (bs *biometricService) IsAvailable() (bool, BiometricType) {
	// 檢查生物識別是否可用
	available := C.checkBiometricAvailability()
	if available == 0 {
		return false, BiometricTypeNone
	}
	
	// 取得生物識別類型
	biometricType := C.getBiometricType()
	switch biometricType {
	case 1:
		return true, BiometricTypeTouchID
	case 2:
		return true, BiometricTypeFaceID
	default:
		return false, BiometricTypeNone
	}
}

// Authenticate 執行生物識別驗證
// 參數：reason（驗證原因，顯示給用戶）
// 回傳：驗證結果
//
// 執行流程：
// 1. 檢查生物識別是否可用
// 2. 準備驗證參數
// 3. 調用 C 函數執行驗證
// 4. 解析驗證結果
// 5. 回傳結果結構
func (bs *biometricService) Authenticate(reason string) *BiometricResult {
	startTime := time.Now()
	
	// 檢查生物識別是否可用
	available, biometricType := bs.IsAvailable()
	if !available {
		return &BiometricResult{
			Success:   false,
			Cancelled: false,
			Error:     errors.New("生物識別驗證不可用"),
			Duration:  time.Since(startTime),
		}
	}
	
	// 準備驗證原因字串
	if reason == "" {
		switch biometricType {
		case BiometricTypeTouchID:
			reason = "請使用 Touch ID 驗證身份"
		case BiometricTypeFaceID:
			reason = "請使用 Face ID 驗證身份"
		default:
			reason = "請使用生物識別驗證身份"
		}
	}
	
	// 轉換為 C 字串
	cReason := C.CString(reason)
	defer C.free(unsafe.Pointer(cReason))
	
	// 執行生物識別驗證
	result := C.performBiometricAuthentication(cReason)
	
	// 解析驗證結果
	duration := time.Since(startTime)
	switch result {
	case 1:
		// 驗證成功
		return &BiometricResult{
			Success:   true,
			Cancelled: false,
			Error:     nil,
			Duration:  duration,
		}
	case -1:
		// 用戶取消
		return &BiometricResult{
			Success:   false,
			Cancelled: true,
			Error:     errors.New("用戶取消驗證"),
			Duration:  duration,
		}
	default:
		// 驗證失敗或系統錯誤
		errorDesc := C.getBiometricErrorDescription(result)
		errorMsg := C.GoString(errorDesc)
		return &BiometricResult{
			Success:   false,
			Cancelled: false,
			Error:     fmt.Errorf("生物識別驗證失敗: %s", errorMsg),
			Duration:  duration,
		}
	}
}

// AuthenticateForNote 為特定筆記執行生物識別驗證
// 參數：noteID（筆記 ID）、reason（驗證原因）
// 回傳：驗證結果
//
// 執行流程：
// 1. 檢查筆記是否啟用生物識別
// 2. 準備驗證原因
// 3. 執行生物識別驗證
// 4. 回傳驗證結果
func (bs *biometricService) AuthenticateForNote(noteID, reason string) *BiometricResult {
	if noteID == "" {
		return &BiometricResult{
			Success:   false,
			Cancelled: false,
			Error:     errors.New("筆記 ID 不能為空"),
			Duration:  0,
		}
	}
	
	// 檢查筆記是否啟用生物識別
	if !bs.IsEnabledForNote(noteID) {
		return &BiometricResult{
			Success:   false,
			Cancelled: false,
			Error:     errors.New("此筆記未啟用生物識別驗證"),
			Duration:  0,
		}
	}
	
	// 準備驗證原因
	if reason == "" {
		reason = "請驗證身份以開啟加密筆記"
	}
	
	// 執行生物識別驗證
	return bs.Authenticate(reason)
}

// SetupForNote 為特定筆記設定生物識別驗證
// 參數：noteID（筆記 ID）
// 回傳：可能的錯誤
//
// 執行流程：
// 1. 驗證筆記 ID
// 2. 檢查生物識別是否可用
// 3. 執行一次驗證確認用戶身份
// 4. 將筆記加入啟用清單
func (bs *biometricService) SetupForNote(noteID string) error {
	if noteID == "" {
		return errors.New("筆記 ID 不能為空")
	}
	
	// 檢查生物識別是否可用
	available, biometricType := bs.IsAvailable()
	if !available {
		return errors.New("生物識別驗證不可用，無法設定")
	}
	
	// 執行一次驗證確認用戶身份
	reason := fmt.Sprintf("請使用 %s 確認身份以啟用生物識別保護", biometricType.String())
	result := bs.Authenticate(reason)
	
	if !result.Success {
		if result.Cancelled {
			return errors.New("用戶取消設定生物識別驗證")
		}
		return fmt.Errorf("設定生物識別驗證失敗: %v", result.Error)
	}
	
	// 將筆記加入啟用清單
	bs.mutex.Lock()
	defer bs.mutex.Unlock()
	bs.enabledNotes[noteID] = true
	
	return nil
}

// RemoveForNote 移除特定筆記的生物識別驗證設定
// 參數：noteID（筆記 ID）
// 回傳：可能的錯誤
//
// 執行流程：
// 1. 驗證筆記 ID
// 2. 從啟用清單中移除筆記
func (bs *biometricService) RemoveForNote(noteID string) error {
	if noteID == "" {
		return errors.New("筆記 ID 不能為空")
	}
	
	bs.mutex.Lock()
	defer bs.mutex.Unlock()
	delete(bs.enabledNotes, noteID)
	
	return nil
}

// IsEnabledForNote 檢查特定筆記是否啟用生物識別驗證
// 參數：noteID（筆記 ID）
// 回傳：是否啟用
//
// 執行流程：
// 1. 檢查筆記 ID 是否在啟用清單中
// 2. 回傳啟用狀態
func (bs *biometricService) IsEnabledForNote(noteID string) bool {
	if noteID == "" {
		return false
	}
	
	bs.mutex.RLock()
	defer bs.mutex.RUnlock()
	return bs.enabledNotes[noteID]
}

