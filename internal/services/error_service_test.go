package services

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"mac-notebook-app/internal/models"
)

// TestNewErrorService 測試錯誤服務的建立
// 驗證服務能夠正確初始化並建立日誌檔案
func TestNewErrorService(t *testing.T) {
	// 建立臨時目錄用於測試
	tempDir := t.TempDir()
	
	// 建立錯誤服務
	service, err := NewErrorService(tempDir)
	if err != nil {
		t.Fatalf("建立錯誤服務失敗: %v", err)
	}
	defer service.(*errorService).Close()

	// 驗證服務不為空
	if service == nil {
		t.Error("錯誤服務不應該為空")
	}

	// 驗證日誌檔案是否建立
	logFileName := "app_" + time.Now().Format("2006-01-02") + ".log"
	logPath := filepath.Join(tempDir, logFileName)
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		t.Error("日誌檔案應該被建立")
	}
}

// TestNewErrorService_InvalidDirectory 測試無效目錄的處理
// 驗證當無法建立日誌目錄時的錯誤處理
func TestNewErrorService_InvalidDirectory(t *testing.T) {
	// 使用無效的目錄路徑（在只讀檔案系統中）
	invalidDir := "/invalid/readonly/path"
	
	service, err := NewErrorService(invalidDir)
	
	// 應該回傳錯誤
	if err == nil {
		t.Error("應該回傳錯誤當目錄無法建立時")
		if service != nil {
			service.(*errorService).Close()
		}
	}
}

// TestLogError 測試錯誤日誌記錄功能
// 驗證錯誤能夠正確記錄到日誌檔案
func TestLogError(t *testing.T) {
	tempDir := t.TempDir()
	service, err := NewErrorService(tempDir)
	if err != nil {
		t.Fatalf("建立錯誤服務失敗: %v", err)
	}
	defer service.(*errorService).Close()

	// 測試記錄標準錯誤
	testErr := errors.New("測試錯誤")
	context := "單元測試"
	
	err = service.LogError(testErr, context)
	if err != nil {
		t.Errorf("記錄錯誤失敗: %v", err)
	}

	// 測試記錄 AppError
	appErr := models.NewAppError(models.ErrFileNotFound, "檔案未找到", "詳細資訊")
	err = service.LogError(appErr, context)
	if err != nil {
		t.Errorf("記錄 AppError 失敗: %v", err)
	}

	// 測試記錄空錯誤
	err = service.LogError(nil, context)
	if err != nil {
		t.Errorf("記錄空錯誤應該成功: %v", err)
	}
}

// TestLocalizeError 測試錯誤訊息本地化功能
// 驗證錯誤訊息能夠正確轉換為繁體中文
func TestLocalizeError(t *testing.T) {
	tempDir := t.TempDir()
	service, err := NewErrorService(tempDir)
	if err != nil {
		t.Fatalf("建立錯誤服務失敗: %v", err)
	}
	defer service.(*errorService).Close()

	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{
			name:     "AppError 本地化",
			err:      models.NewAppError(models.ErrFileNotFound, "File not found", ""),
			expected: "找不到指定的檔案",
		},
		{
			name:     "ValidationError 本地化",
			err:      models.NewValidationError("Title", "標題不能為空"),
			expected: "欄位 'Title' 驗證失敗：標題不能為空",
		},
		{
			name:     "標準錯誤",
			err:      errors.New("未知錯誤"),
			expected: "未知錯誤",
		},
		{
			name:     "空錯誤",
			err:      nil,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.LocalizeError(tt.err)
			if result != tt.expected {
				t.Errorf("LocalizeError() = %v, 期望 %v", result, tt.expected)
			}
		})
	}
}

// TestWrapError 測試錯誤包裝功能
// 驗證錯誤能夠正確包裝並添加上下文資訊
func TestWrapError(t *testing.T) {
	tempDir := t.TempDir()
	service, err := NewErrorService(tempDir)
	if err != nil {
		t.Fatalf("建立錯誤服務失敗: %v", err)
	}
	defer service.(*errorService).Close()

	// 測試包裝標準錯誤
	originalErr := errors.New("原始錯誤")
	context := "測試上下文"
	
	wrappedErr := service.WrapError(originalErr, context)
	if wrappedErr == nil {
		t.Error("包裝後的錯誤不應該為空")
	}

	errMsg := wrappedErr.Error()
	if !strings.Contains(errMsg, context) {
		t.Errorf("包裝後的錯誤應該包含上下文: %s", errMsg)
	}
	if !strings.Contains(errMsg, "原始錯誤") {
		t.Errorf("包裝後的錯誤應該包含原始錯誤: %s", errMsg)
	}

	// 測試包裝空錯誤
	wrappedNil := service.WrapError(nil, context)
	if wrappedNil != nil {
		t.Error("包裝空錯誤應該回傳 nil")
	}
}

// TestHandleError 測試統一錯誤處理功能
// 驗證錯誤能夠被正確處理（記錄和本地化）
func TestHandleError(t *testing.T) {
	tempDir := t.TempDir()
	service, err := NewErrorService(tempDir)
	if err != nil {
		t.Fatalf("建立錯誤服務失敗: %v", err)
	}
	defer service.(*errorService).Close()

	// 測試處理 AppError
	appErr := models.NewAppError(models.ErrInvalidPassword, "Invalid password", "")
	context := "用戶登入"
	
	result := service.HandleError(appErr, context)
	expected := "密碼不正確"
	if result != expected {
		t.Errorf("HandleError() = %v, 期望 %v", result, expected)
	}

	// 測試處理空錯誤
	result = service.HandleError(nil, context)
	if result != "" {
		t.Errorf("處理空錯誤應該回傳空字串，得到: %s", result)
	}
}

// TestCreateAppError 測試應用程式錯誤建立功能
// 驗證能夠正確建立 AppError 實例
func TestCreateAppError(t *testing.T) {
	tempDir := t.TempDir()
	service, err := NewErrorService(tempDir)
	if err != nil {
		t.Fatalf("建立錯誤服務失敗: %v", err)
	}
	defer service.(*errorService).Close()

	code := models.ErrSaveFailed
	message := "保存失敗"
	details := "磁碟空間不足"

	appErr := service.CreateAppError(code, message, details)

	if appErr == nil {
		t.Error("建立的 AppError 不應該為空")
	}
	if appErr.Code != code {
		t.Errorf("錯誤代碼 = %v, 期望 %v", appErr.Code, code)
	}
	if appErr.Message != message {
		t.Errorf("錯誤訊息 = %v, 期望 %v", appErr.Message, message)
	}
	if appErr.Details != details {
		t.Errorf("錯誤詳情 = %v, 期望 %v", appErr.Details, details)
	}
}

// TestIsRetryableError 測試錯誤重試判斷功能
// 驗證能夠正確判斷錯誤是否可重試
func TestIsRetryableError(t *testing.T) {
	tempDir := t.TempDir()
	service, err := NewErrorService(tempDir)
	if err != nil {
		t.Fatalf("建立錯誤服務失敗: %v", err)
	}
	defer service.(*errorService).Close()

	tests := []struct {
		name      string
		err       error
		retryable bool
	}{
		{
			name:      "可重試錯誤 - 保存失敗",
			err:       models.NewAppError(models.ErrSaveFailed, "保存失敗", ""),
			retryable: true,
		},
		{
			name:      "可重試錯誤 - 檔案未找到",
			err:       models.NewAppError(models.ErrFileNotFound, "檔案未找到", ""),
			retryable: true,
		},
		{
			name:      "不可重試錯誤 - 密碼錯誤",
			err:       models.NewAppError(models.ErrInvalidPassword, "密碼錯誤", ""),
			retryable: false,
		},
		{
			name:      "不可重試錯誤 - 權限拒絕",
			err:       models.NewAppError(models.ErrPermissionDenied, "權限拒絕", ""),
			retryable: false,
		},
		{
			name:      "標準錯誤",
			err:       errors.New("標準錯誤"),
			retryable: false,
		},
		{
			name:      "空錯誤",
			err:       nil,
			retryable: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.IsRetryableError(tt.err)
			if result != tt.retryable {
				t.Errorf("IsRetryableError() = %v, 期望 %v", result, tt.retryable)
			}
		})
	}
}

// TestErrorServiceIntegration 測試錯誤服務的整合功能
// 驗證完整的錯誤處理流程
func TestErrorServiceIntegration(t *testing.T) {
	tempDir := t.TempDir()
	service, err := NewErrorService(tempDir)
	if err != nil {
		t.Fatalf("建立錯誤服務失敗: %v", err)
	}
	defer service.(*errorService).Close()

	// 建立一個複雜的錯誤場景
	originalErr := errors.New("底層檔案系統錯誤")
	wrappedErr := service.WrapError(originalErr, "檔案保存操作")
	
	// 記錄錯誤
	logErr := service.LogError(wrappedErr, "用戶保存筆記")
	if logErr != nil {
		t.Errorf("記錄錯誤失敗: %v", logErr)
	}

	// 本地化錯誤
	localizedMsg := service.LocalizeError(wrappedErr)
	if localizedMsg == "" {
		t.Error("本地化訊息不應該為空")
	}

	// 統一處理錯誤
	handledMsg := service.HandleError(wrappedErr, "完整錯誤處理測試")
	if handledMsg == "" {
		t.Error("處理後的訊息不應該為空")
	}
}

// BenchmarkLogError 效能測試 - 錯誤記錄
// 測試錯誤記錄功能的效能表現
func BenchmarkLogError(b *testing.B) {
	tempDir := b.TempDir()
	service, err := NewErrorService(tempDir)
	if err != nil {
		b.Fatalf("建立錯誤服務失敗: %v", err)
	}
	defer service.(*errorService).Close()

	testErr := errors.New("效能測試錯誤")
	context := "效能測試"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.LogError(testErr, context)
	}
}

// BenchmarkLocalizeError 效能測試 - 錯誤本地化
// 測試錯誤本地化功能的效能表現
func BenchmarkLocalizeError(b *testing.B) {
	tempDir := b.TempDir()
	service, err := NewErrorService(tempDir)
	if err != nil {
		b.Fatalf("建立錯誤服務失敗: %v", err)
	}
	defer service.(*errorService).Close()

	appErr := models.NewAppError(models.ErrFileNotFound, "File not found", "")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.LocalizeError(appErr)
	}
}