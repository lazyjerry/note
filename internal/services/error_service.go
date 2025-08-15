package services

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"mac-notebook-app/internal/models"
)



// errorService 錯誤處理服務的具體實作
type errorService struct {
	logFile    *os.File           // 日誌檔案
	logger     *log.Logger        // 日誌記錄器
	errorMsgs  map[string]string  // 錯誤訊息本地化對照表
}

// NewErrorService 建立新的錯誤處理服務實例
// 參數：
//   - logDir: 日誌檔案目錄路徑
// 回傳：ErrorService 介面實例和可能的錯誤
//
// 執行流程：
// 1. 建立日誌目錄（如果不存在）
// 2. 開啟或建立日誌檔案
// 3. 初始化日誌記錄器
// 4. 載入錯誤訊息本地化對照表
func NewErrorService(logDir string) (ErrorService, error) {
	// 確保日誌目錄存在
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("無法建立日誌目錄: %w", err)
	}

	// 建立日誌檔案路徑（按日期命名）
	logFileName := fmt.Sprintf("app_%s.log", time.Now().Format("2006-01-02"))
	logPath := filepath.Join(logDir, logFileName)

	// 開啟日誌檔案（追加模式）
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("無法開啟日誌檔案: %w", err)
	}

	// 建立日誌記錄器
	logger := log.New(logFile, "", log.LstdFlags|log.Lshortfile)

	// 初始化錯誤訊息本地化對照表
	errorMsgs := initializeErrorMessages()

	return &errorService{
		logFile:   logFile,
		logger:    logger,
		errorMsgs: errorMsgs,
	}, nil
}

// LogError 記錄錯誤到日誌檔案
// 參數：
//   - err: 要記錄的錯誤
//   - context: 錯誤發生的上下文資訊
// 回傳：記錄過程中可能發生的錯誤
//
// 執行流程：
// 1. 檢查錯誤是否為空
// 2. 格式化錯誤訊息
// 3. 寫入日誌檔案
func (s *errorService) LogError(err error, context string) error {
	if err == nil {
		return nil
	}

	// 格式化日誌訊息
	logMsg := fmt.Sprintf("[ERROR] Context: %s | Error: %s", context, err.Error())
	
	// 如果是 AppError，添加額外資訊
	if appErr, ok := err.(*models.AppError); ok {
		logMsg += fmt.Sprintf(" | Code: %s", appErr.Code)
		if appErr.Details != "" {
			logMsg += fmt.Sprintf(" | Details: %s", appErr.Details)
		}
	}

	// 寫入日誌
	s.logger.Println(logMsg)
	return nil
}

// LocalizeError 將錯誤訊息本地化為繁體中文
// 參數：
//   - err: 要本地化的錯誤
// 回傳：本地化後的錯誤訊息
//
// 執行流程：
// 1. 檢查錯誤類型
// 2. 查找本地化訊息
// 3. 回傳適當的繁體中文訊息
func (s *errorService) LocalizeError(err error) string {
	if err == nil {
		return ""
	}

	// 處理 AppError 類型
	if appErr, ok := err.(*models.AppError); ok {
		if localizedMsg, exists := s.errorMsgs[appErr.Code]; exists {
			return localizedMsg
		}
		return appErr.Message
	}

	// 處理 ValidationError 類型
	if validationErr, ok := err.(*models.ValidationError); ok {
		return fmt.Sprintf("欄位 '%s' 驗證失敗：%s", validationErr.Field, validationErr.Message)
	}

	// 處理標準錯誤
	errorMsg := err.Error()
	for code, localizedMsg := range s.errorMsgs {
		if errorMsg == code {
			return localizedMsg
		}
	}

	// 如果沒有找到對應的本地化訊息，回傳原始錯誤訊息
	return errorMsg
}

// WrapError 包裝錯誤並添加上下文資訊
// 參數：
//   - err: 原始錯誤
//   - context: 上下文資訊
// 回傳：包裝後的錯誤
func (s *errorService) WrapError(err error, context string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", context, err)
}

// HandleError 統一處理錯誤（記錄日誌並本地化）
// 參數：
//   - err: 要處理的錯誤
//   - context: 錯誤發生的上下文
// 回傳：本地化後的錯誤訊息
//
// 執行流程：
// 1. 記錄錯誤到日誌
// 2. 本地化錯誤訊息
// 3. 回傳處理後的訊息
func (s *errorService) HandleError(err error, context string) string {
	if err == nil {
		return ""
	}

	// 記錄錯誤到日誌
	s.LogError(err, context)

	// 本地化錯誤訊息
	return s.LocalizeError(err)
}

// CreateAppError 建立應用程式特定錯誤
// 參數：
//   - code: 錯誤代碼
//   - message: 錯誤訊息
//   - details: 詳細資訊
// 回傳：AppError 實例
func (s *errorService) CreateAppError(code, message, details string) *models.AppError {
	return models.NewAppError(code, message, details)
}

// IsRetryableError 判斷錯誤是否可重試
// 參數：
//   - err: 要檢查的錯誤
// 回傳：是否可重試
//
// 執行流程：
// 1. 檢查錯誤類型
// 2. 根據錯誤代碼判斷是否可重試
func (s *errorService) IsRetryableError(err error) bool {
	if err == nil {
		return false
	}

	// 檢查 AppError 類型
	if appErr, ok := err.(*models.AppError); ok {
		switch appErr.Code {
		case models.ErrSaveFailed, models.ErrFileNotFound:
			return true // 這些錯誤可能是暫時性的，可以重試
		case models.ErrInvalidPassword, models.ErrPermissionDenied:
			return false // 這些錯誤不應該重試
		default:
			return false
		}
	}

	return false
}

// Close 關閉錯誤服務並釋放資源
// 執行流程：
// 1. 關閉日誌檔案
// 2. 清理資源
func (s *errorService) Close() error {
	if s.logFile != nil {
		return s.logFile.Close()
	}
	return nil
}

// initializeErrorMessages 初始化錯誤訊息本地化對照表
// 回傳：錯誤代碼到繁體中文訊息的對照表
//
// 執行流程：
// 1. 建立對照表
// 2. 添加所有錯誤代碼的繁體中文翻譯
func initializeErrorMessages() map[string]string {
	return map[string]string{
		models.ErrFileNotFound:     "找不到指定的檔案",
		models.ErrInvalidPassword:  "密碼不正確",
		models.ErrEncryptionFailed: "檔案加密失敗",
		models.ErrBiometricFailed:  "生物識別驗證失敗",
		models.ErrSaveFailed:       "檔案保存失敗",
		models.ErrPermissionDenied: "存取權限被拒絕",
		models.ErrValidationFailed: "資料驗證失敗",
		
		// 檔案操作相關錯誤
		"FILE_READ_ERROR":    "檔案讀取失敗",
		"FILE_WRITE_ERROR":   "檔案寫入失敗",
		"FILE_DELETE_ERROR":  "檔案刪除失敗",
		"DIRECTORY_ERROR":    "目錄操作失敗",
		
		// 加密相關錯誤
		"ENCRYPTION_KEY_ERROR": "加密金鑰錯誤",
		"DECRYPTION_ERROR":     "檔案解密失敗",
		"KEY_GENERATION_ERROR": "金鑰生成失敗",
		
		// 網路和系統錯誤
		"NETWORK_ERROR":    "網路連線錯誤",
		"SYSTEM_ERROR":     "系統錯誤",
		"MEMORY_ERROR":     "記憶體不足",
		"DISK_SPACE_ERROR": "磁碟空間不足",
		
		// 用戶介面錯誤
		"UI_RENDER_ERROR":   "介面渲染失敗",
		"UI_EVENT_ERROR":    "介面事件處理失敗",
		"DIALOG_ERROR":      "對話框顯示失敗",
		
		// 自動保存相關錯誤
		"AUTO_SAVE_ERROR":     "自動保存失敗",
		"SAVE_CONFLICT_ERROR": "保存衝突",
		"BACKUP_ERROR":        "備份建立失敗",
	}
}