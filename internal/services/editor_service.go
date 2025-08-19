// Package services 實作應用程式的業務邏輯服務
// 本檔案包含 EditorService 的具體實作，負責處理筆記的編輯、保存和 Markdown 預覽功能
package services

import (
	"bytes"                           // 位元組緩衝區處理
	"fmt"                            // 格式化輸出
	"mac-notebook-app/internal/models" // 引入資料模型
	"mac-notebook-app/internal/repositories" // 引入資料存取層
	"path/filepath"                  // 檔案路徑處理
	"strings"                        // 字串處理
	"time"                          // 時間處理
	"github.com/google/uuid"        // UUID 生成

	"github.com/yuin/goldmark"       // Markdown 解析器
	"github.com/yuin/goldmark/extension" // Markdown 擴展功能
	"github.com/yuin/goldmark/parser"    // Markdown 解析器配置
	"github.com/yuin/goldmark/renderer/html" // HTML 渲染器
)

// editorService 實作 EditorService 介面
// 負責處理筆記的核心編輯功能，包含建立、開啟、保存、更新和 Markdown 預覽
// 整合加密功能，支援加密檔案的開啟、保存和管理
// 新增效能優化功能，支援大檔案處理和記憶體管理
// 整合智慧編輯功能，提供進階的編輯輔助工具
type editorService struct {
	fileRepo      repositories.FileRepository // 檔案存取介面
	encryptionSvc EncryptionService           // 加密服務介面
	passwordSvc   PasswordService             // 密碼服務介面
	biometricSvc  BiometricService            // 生物識別服務介面
	smartEditSvc  SmartEditingService         // 智慧編輯服務介面
	markdown      goldmark.Markdown           // Markdown 解析器實例
	activeNotes   map[string]*models.Note     // 當前開啟的筆記快取
	perfService   PerformanceService          // 效能服務介面
	
	// 效能優化相關欄位
	maxCacheSize     int                      // 最大快取大小
	largeFileThreshold int64                  // 大檔案閾值（位元組）
	chunkSize        int64                    // 分塊處理大小
}

// NewEditorService 建立新的編輯器服務實例
// 參數：
//   - fileRepo: 檔案存取介面
//   - encryptionSvc: 加密服務介面
//   - passwordSvc: 密碼服務介面
//   - biometricSvc: 生物識別服務介面
//   - perfService: 效能服務介面（可選，用於效能監控和優化）
//   - smartEditSvc: 智慧編輯服務介面（可選，用於進階編輯功能）
// 回傳：EditorService 介面實例
//
// 執行流程：
// 1. 初始化 goldmark Markdown 解析器，啟用常用擴展功能
// 2. 建立筆記快取映射表
// 3. 整合加密相關服務
// 4. 整合智慧編輯服務
// 5. 設定效能優化參數
// 6. 回傳配置完成的編輯器服務實例
func NewEditorService(fileRepo repositories.FileRepository, encryptionSvc EncryptionService, passwordSvc PasswordService, biometricSvc BiometricService, perfService PerformanceService, smartEditSvc SmartEditingService) EditorService {
	// 配置 Markdown 解析器，啟用表格、刪除線、任務列表等擴展功能
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,        // GitHub Flavored Markdown 支援
			extension.Table,      // 表格支援
			extension.Strikethrough, // 刪除線支援
			extension.TaskList,   // 任務列表支援
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(), // 自動生成標題 ID
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(), // 硬換行支援
			html.WithXHTML(),     // XHTML 相容性
		),
	)

	// 如果沒有提供智慧編輯服務，建立預設實例
	if smartEditSvc == nil {
		smartEditSvc = NewSmartEditingService()
	}

	return &editorService{
		fileRepo:           fileRepo,
		encryptionSvc:      encryptionSvc,
		passwordSvc:        passwordSvc,
		biometricSvc:       biometricSvc,
		smartEditSvc:       smartEditSvc,
		markdown:           md,
		activeNotes:        make(map[string]*models.Note),
		perfService:        perfService,
		maxCacheSize:       100,              // 最多快取 100 個筆記
		largeFileThreshold: 5 * 1024 * 1024,  // 5MB 以上視為大檔案
		chunkSize:          1024 * 1024,      // 1MB 分塊大小
	}
}

// CreateNote 建立新的筆記
// 參數：title（筆記標題）、content（筆記內容）
// 回傳：建立的筆記實例和可能的錯誤
//
// 執行流程：
// 1. 生成唯一的筆記 ID
// 2. 建立筆記實例並設定基本屬性
// 3. 將筆記加入活躍筆記快取
// 4. 回傳建立的筆記實例
func (e *editorService) CreateNote(title, content string) (*models.Note, error) {
	// 生成唯一的筆記 ID
	noteID := uuid.New().String()
	
	// 建立新筆記實例
	note := &models.Note{
		ID:          noteID,
		Title:       title,
		Content:     content,
		IsEncrypted: false, // 預設不加密
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// 將新筆記加入活躍筆記快取
	e.activeNotes[noteID] = note

	return note, nil
}

// OpenNote 從檔案路徑開啟筆記
// 參數：filePath（筆記檔案路徑）
// 回傳：開啟的筆記實例和可能的錯誤
//
// 執行流程：
// 1. 檢查檔案是否存在
// 2. 讀取檔案內容
// 3. 檢查是否為加密檔案並進行解密
// 4. 解析檔案資訊（標題、加密狀態等）
// 5. 建立筆記實例
// 6. 將筆記加入活躍筆記快取
// 7. 回傳筆記實例
func (e *editorService) OpenNote(filePath string) (*models.Note, error) {
	// 檢查檔案是否存在
	if !e.fileRepo.FileExists(filePath) {
		return nil, fmt.Errorf("檔案不存在: %s", filePath)
	}

	// 讀取檔案內容
	rawContent, err := e.fileRepo.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("讀取檔案失敗: %w", err)
	}

	// 從檔案路徑提取標題（去除副檔名）
	fileName := filepath.Base(filePath)
	title := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	
	// 檢查是否為加密檔案（副檔名為 .enc）
	isEncrypted := strings.HasSuffix(filePath, ".enc")
	
	// 處理檔案內容（解密或直接使用）
	var content string
	if isEncrypted {
		// 加密檔案需要解密
		content, err = e.decryptFileContent(rawContent, title)
		if err != nil {
			return nil, fmt.Errorf("解密檔案失敗: %w", err)
		}
		
		// 移除 .enc 副檔名以取得正確的標題
		if strings.HasSuffix(title, ".md") {
			title = strings.TrimSuffix(title, ".md")
		}
	} else {
		// 非加密檔案直接使用
		content = string(rawContent)
	}

	// 生成筆記 ID
	noteID := uuid.New().String()

	// 建立筆記實例
	note := &models.Note{
		ID:          noteID,
		Title:       title,
		Content:     content,
		FilePath:    filePath,
		IsEncrypted: isEncrypted,
		CreatedAt:   time.Now(), // 實際應用中可能需要從檔案屬性讀取
		UpdatedAt:   time.Now(),
	}

	// 將筆記加入活躍筆記快取
	e.activeNotes[noteID] = note

	return note, nil
}

// SaveNote 保存筆記到檔案系統
// 參數：note（要保存的筆記實例）
// 回傳：可能的錯誤
//
// 執行流程：
// 1. 驗證筆記實例的有效性
// 2. 確定保存路徑（如果未設定則生成預設路徑）
// 3. 處理筆記內容（加密或直接使用）
// 4. 將處理後的內容寫入檔案
// 5. 更新筆記的最後保存時間
// 6. 更新活躍筆記快取
func (e *editorService) SaveNote(note *models.Note) error {
	if note == nil {
		return fmt.Errorf("筆記實例不能為空")
	}

	// 如果筆記沒有設定檔案路徑，生成預設路徑
	if note.FilePath == "" {
		// 清理標題作為檔案名稱
		fileName := e.sanitizeFileName(note.Title)
		if fileName == "" {
			fileName = "untitled"
		}
		
		// 根據加密狀態決定副檔名
		extension := ".md"
		if note.IsEncrypted {
			extension = ".md.enc"
		}
		
		note.FilePath = fileName + extension
	}

	// 處理筆記內容（加密或直接使用）
	var contentToSave []byte
	var err error
	
	if note.IsEncrypted {
		// 加密筆記內容
		contentToSave, err = e.encryptFileContent(note.Content, note.ID)
		if err != nil {
			return fmt.Errorf("加密筆記內容失敗: %w", err)
		}
	} else {
		// 非加密筆記直接轉換為位元組
		contentToSave = []byte(note.Content)
	}

	// 將處理後的內容寫入檔案
	err = e.fileRepo.WriteFile(note.FilePath, contentToSave)
	if err != nil {
		return fmt.Errorf("保存筆記失敗: %w", err)
	}

	// 更新筆記的時間戳
	note.UpdatedAt = time.Now()
	note.LastSaved = time.Now()

	// 更新活躍筆記快取
	e.activeNotes[note.ID] = note

	return nil
}

// UpdateContent 更新指定筆記的內容
// 參數：noteID（筆記 ID）、content（新的內容）
// 回傳：可能的錯誤
//
// 執行流程：
// 1. 從活躍筆記快取中查找指定筆記
// 2. 更新筆記內容
// 3. 更新最後修改時間
// 4. 更新快取中的筆記實例
func (e *editorService) UpdateContent(noteID, content string) error {
	// 從活躍筆記快取中查找筆記
	note, exists := e.activeNotes[noteID]
	if !exists {
		return fmt.Errorf("找不到指定的筆記: %s", noteID)
	}

	// 更新筆記內容和時間戳
	note.Content = content
	note.UpdatedAt = time.Now()

	// 更新快取中的筆記實例
	e.activeNotes[noteID] = note

	return nil
}

// PreviewMarkdown 將 Markdown 內容轉換為 HTML 預覽
// 參數：content（Markdown 格式的內容）
// 回傳：轉換後的 HTML 字串
//
// 執行流程：
// 1. 建立輸出緩衝區
// 2. 使用 goldmark 解析器將 Markdown 轉換為 HTML
// 3. 處理轉換錯誤（如果有）
// 4. 回傳 HTML 字串
func (e *editorService) PreviewMarkdown(content string) string {
	// 建立輸出緩衝區
	var buf bytes.Buffer
	
	// 使用 goldmark 將 Markdown 轉換為 HTML
	err := e.markdown.Convert([]byte(content), &buf)
	if err != nil {
		// 如果轉換失敗，回傳錯誤訊息的 HTML
		return fmt.Sprintf("<p>Markdown 轉換錯誤: %s</p>", err.Error())
	}

	return buf.String()
}

// sanitizeFileName 清理檔案名稱，移除不合法的字元
// 參數：fileName（原始檔案名稱）
// 回傳：清理後的檔案名稱
//
// 執行流程：
// 1. 定義不合法字元列表
// 2. 逐一替換不合法字元為底線
// 3. 移除多餘的空白字元
// 4. 回傳清理後的檔案名稱
func (e *editorService) sanitizeFileName(fileName string) string {
	// 定義不合法的檔案名稱字元
	invalidChars := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	
	// 替換不合法字元為底線
	result := fileName
	for _, char := range invalidChars {
		result = strings.ReplaceAll(result, char, "_")
	}
	
	// 移除前後空白並限制長度
	result = strings.TrimSpace(result)
	if len(result) > 100 {
		result = result[:100]
	}
	
	return result
}

// GetActiveNote 取得指定 ID 的活躍筆記
// 參數：noteID（筆記 ID）
// 回傳：筆記實例和是否存在的布林值
//
// 執行流程：
// 1. 從活躍筆記快取中查找指定筆記
// 2. 回傳筆記實例和存在狀態
func (e *editorService) GetActiveNote(noteID string) (*models.Note, bool) {
	note, exists := e.activeNotes[noteID]
	return note, exists
}

// CloseNote 關閉指定的筆記，從活躍快取中移除
// 參數：noteID（筆記 ID）
//
// 執行流程：
// 1. 從活躍筆記快取中移除指定筆記
func (e *editorService) CloseNote(noteID string) {
	delete(e.activeNotes, noteID)
}

// GetActiveNotes 取得所有活躍筆記的列表
// 回傳：活躍筆記的映射表
func (e *editorService) GetActiveNotes() map[string]*models.Note {
	// 建立副本以避免外部修改
	result := make(map[string]*models.Note)
	for id, note := range e.activeNotes {
		result[id] = note
	}
	return result
}

// EnableEncryption 為指定筆記啟用加密
// 參數：
//   - noteID: 筆記 ID
//   - password: 加密密碼
//   - algorithm: 加密演算法
//   - useBiometric: 是否啟用生物識別驗證
// 回傳：可能的錯誤
//
// 執行流程：
// 1. 驗證筆記是否存在於活躍快取中
// 2. 驗證密碼強度
// 3. 設定筆記的加密狀態
// 4. 如果啟用生物識別，設定生物識別驗證
// 5. 更新檔案路徑以包含 .enc 副檔名
// 6. 更新活躍筆記快取
func (e *editorService) EnableEncryption(noteID, password, algorithm string, useBiometric bool) error {
	// 檢查筆記是否存在
	note, exists := e.activeNotes[noteID]
	if !exists {
		return fmt.Errorf("找不到指定的筆記: %s", noteID)
	}

	// 驗證密碼強度
	if !e.encryptionSvc.ValidatePassword(password) {
		return fmt.Errorf("密碼不符合安全要求")
	}

	// 設定筆記的加密狀態
	note.IsEncrypted = true
	if algorithm != "" {
		note.EncryptionType = algorithm
	} else {
		note.EncryptionType = "aes256" // 預設使用 AES-256
	}

	// 如果啟用生物識別驗證
	available, _ := e.biometricSvc.IsAvailable()
	if useBiometric && available {
		err := e.biometricSvc.SetupForNote(noteID)
		if err != nil {
			return fmt.Errorf("設定生物識別驗證失敗: %w", err)
		}
	}

	// 更新檔案路徑以包含 .enc 副檔名
	if note.FilePath != "" && !strings.HasSuffix(note.FilePath, ".enc") {
		note.FilePath = note.FilePath + ".enc"
	}

	// 更新時間戳
	note.UpdatedAt = time.Now()

	// 更新活躍筆記快取
	e.activeNotes[noteID] = note

	return nil
}

// DisableEncryption 為指定筆記停用加密
// 參數：noteID（筆記 ID）
// 回傳：可能的錯誤
//
// 執行流程：
// 1. 驗證筆記是否存在於活躍快取中
// 2. 移除筆記的加密狀態
// 3. 移除生物識別驗證設定
// 4. 更新檔案路徑移除 .enc 副檔名
// 5. 更新活躍筆記快取
func (e *editorService) DisableEncryption(noteID string) error {
	// 檢查筆記是否存在
	note, exists := e.activeNotes[noteID]
	if !exists {
		return fmt.Errorf("找不到指定的筆記: %s", noteID)
	}

	// 移除筆記的加密狀態
	note.IsEncrypted = false
	note.EncryptionType = ""

	// 移除生物識別驗證設定
	e.biometricSvc.RemoveForNote(noteID)

	// 更新檔案路徑移除 .enc 副檔名
	if note.FilePath != "" && strings.HasSuffix(note.FilePath, ".enc") {
		note.FilePath = strings.TrimSuffix(note.FilePath, ".enc")
	}

	// 更新時間戳
	note.UpdatedAt = time.Now()

	// 更新活躍筆記快取
	e.activeNotes[noteID] = note

	return nil
}

// decryptFileContent 解密檔案內容
// 參數：
//   - encryptedData: 加密的檔案內容
//   - noteID: 筆記 ID（用於生物識別驗證）
// 回傳：解密後的內容和可能的錯誤
//
// 執行流程：
// 1. 嘗試使用生物識別驗證（如果可用）
// 2. 如果生物識別失敗或不可用，提示輸入密碼
// 3. 使用密碼解密內容
// 4. 回傳解密後的內容
func (e *editorService) decryptFileContent(encryptedData []byte, noteID string) (string, error) {
	// 首先嘗試生物識別驗證
	if e.biometricSvc.IsEnabledForNote(noteID) {
		result := e.biometricSvc.AuthenticateForNote(noteID, "開啟加密筆記")
		if result.Error == nil && result.Success {
			// 生物識別成功，使用儲存的密碼解密
			// 注意：實際實作中需要安全地儲存和檢索密碼
			// 這裡簡化處理，實際應用中應該使用 Keychain 或類似的安全儲存
			return e.decryptWithStoredCredentials(encryptedData, noteID)
		}
	}

	// 生物識別失敗或不可用，需要密碼驗證
	// 注意：在實際的 UI 實作中，這裡應該彈出密碼輸入對話框
	// 目前返回錯誤，要求上層處理密碼輸入
	return "", fmt.Errorf("需要密碼驗證才能開啟加密檔案")
}

// DecryptWithPassword 使用密碼解密筆記內容
// 參數：
//   - noteID: 筆記 ID
//   - password: 解密密碼
// 回傳：解密後的內容和可能的錯誤
//
// 執行流程：
// 1. 檢查筆記是否存在且已加密
// 2. 讀取加密檔案內容
// 3. 使用密碼解密內容
// 4. 更新筆記內容
// 5. 回傳解密後的內容
func (e *editorService) DecryptWithPassword(noteID, password string) (string, error) {
	// 檢查筆記是否存在
	note, exists := e.activeNotes[noteID]
	if !exists {
		return "", fmt.Errorf("找不到指定的筆記: %s", noteID)
	}

	if !note.IsEncrypted {
		return note.Content, nil // 非加密筆記直接回傳內容
	}

	// 讀取加密檔案內容
	encryptedData, err := e.fileRepo.ReadFile(note.FilePath)
	if err != nil {
		return "", fmt.Errorf("讀取加密檔案失敗: %w", err)
	}

	// 使用密碼解密內容
	algorithm := note.EncryptionType
	if algorithm == "" {
		algorithm = "aes256" // 預設演算法
	}

	decryptedContent, err := e.encryptionSvc.DecryptContent(encryptedData, password, algorithm)
	if err != nil {
		return "", fmt.Errorf("解密失敗: %w", err)
	}

	// 更新筆記內容
	note.Content = decryptedContent
	note.UpdatedAt = time.Now()
	e.activeNotes[noteID] = note

	return decryptedContent, nil
}

// encryptFileContent 加密檔案內容
// 參數：
//   - content: 要加密的內容
//   - noteID: 筆記 ID
// 回傳：加密後的位元組陣列和可能的錯誤
//
// 執行流程：
// 1. 取得筆記的加密設定
// 2. 取得加密密碼（從安全儲存或提示輸入）
// 3. 使用指定演算法加密內容
// 4. 回傳加密後的資料
func (e *editorService) encryptFileContent(content, noteID string) ([]byte, error) {
	// 取得筆記實例
	note, exists := e.activeNotes[noteID]
	if !exists {
		return nil, fmt.Errorf("找不到指定的筆記: %s", noteID)
	}

	// 取得加密演算法
	algorithm := note.EncryptionType
	if algorithm == "" {
		algorithm = "aes256" // 預設演算法
	}

	// 注意：在實際實作中，這裡需要安全地取得密碼
	// 可能從 Keychain、用戶輸入或其他安全儲存中取得
	// 目前返回錯誤，要求上層處理密碼取得
	return nil, fmt.Errorf("需要密碼才能加密檔案內容")
}

// EncryptWithPassword 使用密碼加密筆記內容
// 參數：
//   - noteID: 筆記 ID
//   - password: 加密密碼
// 回傳：可能的錯誤
//
// 執行流程：
// 1. 檢查筆記是否存在
// 2. 取得加密演算法
// 3. 使用密碼加密內容
// 4. 更新筆記狀態
func (e *editorService) EncryptWithPassword(noteID, password string) error {
	// 檢查筆記是否存在
	note, exists := e.activeNotes[noteID]
	if !exists {
		return fmt.Errorf("找不到指定的筆記: %s", noteID)
	}

	// 取得加密演算法
	algorithm := note.EncryptionType
	if algorithm == "" {
		algorithm = "aes256" // 預設演算法
	}

	// 加密內容
	_, err := e.encryptionSvc.EncryptContent(note.Content, password, algorithm)
	if err != nil {
		return fmt.Errorf("加密失敗: %w", err)
	}

	// 更新筆記狀態
	note.UpdatedAt = time.Now()
	e.activeNotes[noteID] = note

	return nil
}

// decryptWithStoredCredentials 使用儲存的憑證解密內容
// 參數：
//   - encryptedData: 加密的資料
//   - noteID: 筆記 ID
// 回傳：解密後的內容和可能的錯誤
//
// 注意：這是一個簡化的實作，實際應用中需要安全的憑證管理
func (e *editorService) decryptWithStoredCredentials(encryptedData []byte, noteID string) (string, error) {
	// 實際實作中應該從安全儲存（如 Keychain）中取得密碼
	// 目前返回錯誤，表示需要進一步實作
	return "", fmt.Errorf("安全憑證管理功能尚未實作")
}

// IsEncrypted 檢查指定筆記是否已加密
// 參數：noteID（筆記 ID）
// 回傳：是否已加密和筆記是否存在
func (e *editorService) IsEncrypted(noteID string) (bool, bool) {
	note, exists := e.activeNotes[noteID]
	if !exists {
		return false, false
	}
	return note.IsEncrypted, true
}

// GetEncryptionType 取得指定筆記的加密類型
// 參數：noteID（筆記 ID）
// 回傳：加密類型和筆記是否存在
func (e *editorService) GetEncryptionType(noteID string) (string, bool) {
	note, exists := e.activeNotes[noteID]
	if !exists {
		return "", false
	}
	return note.EncryptionType, true
}

// OptimizeForLargeFile 為大檔案處理進行優化
// 參數：
//   - filePath: 檔案路徑
//   - fileSize: 檔案大小
// 回傳：可能的錯誤
//
// 執行流程：
// 1. 檢查檔案大小是否超過大檔案閾值
// 2. 如果是大檔案，執行記憶體優化
// 3. 調整快取策略
// 4. 通知效能服務進行優化
func (e *editorService) OptimizeForLargeFile(filePath string, fileSize int64) error {
	// 檢查是否為大檔案
	if fileSize > e.largeFileThreshold {
		// 清理快取以釋放記憶體
		if err := e.optimizeCache(); err != nil {
			return fmt.Errorf("快取優化失敗: %w", err)
		}
		
		// 通知效能服務進行優化
		if e.perfService != nil {
			if err := e.perfService.OptimizeForLargeFile(filePath, fileSize); err != nil {
				return fmt.Errorf("效能服務優化失敗: %w", err)
			}
		}
	}
	
	return nil
}

// ProcessLargeFileInChunks 分塊處理大檔案內容
// 參數：
//   - content: 檔案內容
//   - processor: 處理函數
// 回傳：處理結果和可能的錯誤
//
// 執行流程：
// 1. 將內容分割成指定大小的塊
// 2. 逐塊處理內容
// 3. 在處理過程中監控記憶體使用
// 4. 必要時執行垃圾回收
func (e *editorService) ProcessLargeFileInChunks(content string, processor func(string) (string, error)) (string, error) {
	contentBytes := []byte(content)
	contentSize := int64(len(contentBytes))
	
	// 如果內容不大，直接處理
	if contentSize <= e.largeFileThreshold {
		return processor(content)
	}
	
	// 分塊處理大內容
	var result strings.Builder
	result.Grow(len(content)) // 預分配記憶體
	
	for i := int64(0); i < contentSize; i += e.chunkSize {
		end := i + e.chunkSize
		if end > contentSize {
			end = contentSize
		}
		
		chunk := string(contentBytes[i:end])
		processedChunk, err := processor(chunk)
		if err != nil {
			return "", fmt.Errorf("處理第 %d 塊時發生錯誤: %w", i/e.chunkSize+1, err)
		}
		
		result.WriteString(processedChunk)
		
		// 每處理 10 個塊後檢查記憶體使用
		if (i/e.chunkSize+1)%10 == 0 {
			if e.perfService != nil {
				memUsage, _ := e.perfService.GetMemoryUsage()
				// 如果記憶體使用超過 100MB，執行垃圾回收
				if memUsage > 100*1024*1024 {
					e.perfService.ForceGarbageCollection()
				}
			}
		}
	}
	
	return result.String(), nil
}

// optimizeCache 優化筆記快取
// 回傳：可能的錯誤
//
// 執行流程：
// 1. 檢查快取大小是否超過限制
// 2. 如果超過，移除最舊的筆記
// 3. 記錄快取優化統計
func (e *editorService) optimizeCache() error {
	if len(e.activeNotes) <= e.maxCacheSize {
		return nil // 快取大小在限制內
	}
	
	// 找出最舊的筆記並移除
	var oldestID string
	var oldestTime time.Time = time.Now()
	
	for id, note := range e.activeNotes {
		if note.UpdatedAt.Before(oldestTime) {
			oldestTime = note.UpdatedAt
			oldestID = id
		}
	}
	
	// 移除最舊的筆記
	if oldestID != "" {
		delete(e.activeNotes, oldestID)
		
		// 記錄快取未命中（因為移除了快取項目）
		if e.perfService != nil {
			if ps, ok := e.perfService.(*performanceService); ok {
				ps.RecordCacheMiss()
			}
		}
	}
	
	return nil
}

// PreviewMarkdownOptimized 優化的 Markdown 預覽功能
// 參數：content（Markdown 格式的內容）
// 回傳：轉換後的 HTML 字串
//
// 執行流程：
// 1. 檢查內容大小
// 2. 如果是大內容，使用分塊處理
// 3. 對每個塊進行 Markdown 轉換
// 4. 合併結果並回傳
func (e *editorService) PreviewMarkdownOptimized(content string) string {
	contentSize := int64(len(content))
	
	// 小內容直接處理
	if contentSize <= e.largeFileThreshold {
		return e.PreviewMarkdown(content)
	}
	
	// 大內容分塊處理
	result, err := e.ProcessLargeFileInChunks(content, func(chunk string) (string, error) {
		var buf bytes.Buffer
		err := e.markdown.Convert([]byte(chunk), &buf)
		if err != nil {
			return "", err
		}
		return buf.String(), nil
	})
	
	if err != nil {
		return fmt.Sprintf("<p>Markdown 轉換錯誤: %s</p>", err.Error())
	}
	
	return result
}

// GetCacheStats 取得快取統計資訊
// 回傳：快取統計資料
//
// 執行流程：
// 1. 計算當前快取使用情況
// 2. 計算快取使用率
// 3. 回傳統計資料
func (e *editorService) GetCacheStats() map[string]interface{} {
	cacheUsage := float64(len(e.activeNotes)) / float64(e.maxCacheSize)
	
	return map[string]interface{}{
		"active_notes":    len(e.activeNotes),
		"max_cache_size":  e.maxCacheSize,
		"cache_usage":     cacheUsage,
		"large_file_threshold": e.largeFileThreshold,
		"chunk_size":      e.chunkSize,
	}
}

// ClearCache 清空筆記快取
// 回傳：可能的錯誤
//
// 執行流程：
// 1. 清空所有活躍筆記
// 2. 通知效能服務記錄快取清理
func (e *editorService) ClearCache() error {
	// 記錄清理前的快取大小
	cacheSize := len(e.activeNotes)
	
	// 清空快取
	e.activeNotes = make(map[string]*models.Note)
	
	// 通知效能服務
	if e.perfService != nil {
		if ps, ok := e.perfService.(*performanceService); ok {
			// 記錄多次快取未命中以反映清理操作
			for i := 0; i < cacheSize; i++ {
				ps.RecordCacheMiss()
			}
		}
	}
	
	return nil
}

// MonitorMemoryUsage 監控記憶體使用情況
// 回傳：當前記憶體使用量和建議
//
// 執行流程：
// 1. 取得當前記憶體使用量
// 2. 分析快取使用情況
// 3. 提供優化建議
func (e *editorService) MonitorMemoryUsage() (int64, string, error) {
	if e.perfService == nil {
		return 0, "效能服務未啟用", nil
	}
	
	memUsage, err := e.perfService.GetMemoryUsage()
	if err != nil {
		return 0, "", fmt.Errorf("取得記憶體使用量失敗: %w", err)
	}
	
	var suggestion string
	cacheSize := len(e.activeNotes)
	
	// 根據記憶體使用情況提供建議
	if memUsage > 200*1024*1024 { // 超過 200MB
		suggestion = "記憶體使用量較高，建議清理快取或關閉不需要的筆記"
	} else if cacheSize > e.maxCacheSize*8/10 { // 快取使用率超過 80%
		suggestion = "快取使用率較高，建議關閉一些不常用的筆記"
	} else {
		suggestion = "記憶體使用正常"
	}
	
	return memUsage, suggestion, nil
}

// ========== 智慧編輯功能實作 ==========

// GetAutoCompleteSuggestions 取得自動完成建議
// 參數：content（當前內容）、cursorPosition（游標位置）
// 回傳：自動完成建議陣列
//
// 執行流程：
// 1. 委託給智慧編輯服務處理
// 2. 回傳自動完成建議列表
func (e *editorService) GetAutoCompleteSuggestions(content string, cursorPosition int) []AutoCompleteSuggestion {
	if e.smartEditSvc == nil {
		return []AutoCompleteSuggestion{}
	}
	return e.smartEditSvc.AutoCompleteMarkdown(content, cursorPosition)
}

// FormatTableContent 格式化表格內容
// 參數：tableContent（表格內容）
// 回傳：格式化後的表格字串和可能的錯誤
//
// 執行流程：
// 1. 委託給智慧編輯服務處理表格格式化
// 2. 回傳格式化結果
func (e *editorService) FormatTableContent(tableContent string) (string, error) {
	if e.smartEditSvc == nil {
		return tableContent, fmt.Errorf("智慧編輯服務未啟用")
	}
	return e.smartEditSvc.FormatTable(tableContent)
}

// InsertLinkMarkdown 插入 Markdown 連結
// 參數：text（連結文字）、url（連結網址）
// 回傳：格式化的 Markdown 連結字串
//
// 執行流程：
// 1. 委託給智慧編輯服務處理連結插入
// 2. 回傳格式化的連結字串
func (e *editorService) InsertLinkMarkdown(text, url string) string {
	if e.smartEditSvc == nil {
		return fmt.Sprintf("[%s](%s)", text, url)
	}
	return e.smartEditSvc.InsertLink(text, url)
}

// InsertImageMarkdown 插入 Markdown 圖片
// 參數：altText（替代文字）、imagePath（圖片路徑）
// 回傳：格式化的 Markdown 圖片字串
//
// 執行流程：
// 1. 委託給智慧編輯服務處理圖片插入
// 2. 回傳格式化的圖片字串
func (e *editorService) InsertImageMarkdown(altText, imagePath string) string {
	if e.smartEditSvc == nil {
		return fmt.Sprintf("![%s](%s)", altText, imagePath)
	}
	return e.smartEditSvc.InsertImage(altText, imagePath)
}

// GetSupportedCodeLanguages 取得支援的程式語言列表
// 回傳：支援的程式語言陣列
//
// 執行流程：
// 1. 委託給智慧編輯服務取得支援的語言列表
// 2. 回傳語言陣列
func (e *editorService) GetSupportedCodeLanguages() []string {
	if e.smartEditSvc == nil {
		return []string{"text", "go", "javascript", "python", "html", "css"}
	}
	return e.smartEditSvc.GetSupportedLanguages()
}

// FormatCodeBlockMarkdown 格式化程式碼區塊
// 參數：code（程式碼內容）、language（程式語言）
// 回傳：格式化的 Markdown 程式碼區塊
//
// 執行流程：
// 1. 委託給智慧編輯服務處理程式碼區塊格式化
// 2. 回傳格式化的程式碼區塊
func (e *editorService) FormatCodeBlockMarkdown(code, language string) string {
	if e.smartEditSvc == nil {
		return fmt.Sprintf("```%s\n%s\n```", language, code)
	}
	return e.smartEditSvc.FormatCodeBlock(code, language)
}

// FormatMathExpressionMarkdown 格式化數學公式
// 參數：expression（數學表達式）、isInline（是否為行內公式）
// 回傳：格式化的 LaTeX 數學公式字串
//
// 執行流程：
// 1. 委託給智慧編輯服務處理數學公式格式化
// 2. 回傳格式化的數學公式
func (e *editorService) FormatMathExpressionMarkdown(expression string, isInline bool) string {
	if e.smartEditSvc == nil {
		if isInline {
			return fmt.Sprintf("$%s$", expression)
		}
		return fmt.Sprintf("$$\n%s\n$$", expression)
	}
	return e.smartEditSvc.FormatMathExpression(expression, isInline)
}

// ValidateMarkdownContent 驗證 Markdown 內容的語法正確性
// 參數：content（要驗證的 Markdown 內容）
// 回傳：驗證結果和可能的錯誤列表
//
// 執行流程：
// 1. 委託給智慧編輯服務進行語法驗證
// 2. 回傳驗證結果和錯誤列表
func (e *editorService) ValidateMarkdownContent(content string) (bool, []string) {
	if e.smartEditSvc == nil {
		return true, []string{} // 如果沒有智慧編輯服務，預設為有效
	}
	return e.smartEditSvc.ValidateMarkdownSyntax(content)
}

// GenerateTableTemplateMarkdown 生成表格模板
// 參數：rows（行數）、cols（列數）
// 回傳：表格模板字串
//
// 執行流程：
// 1. 委託給智慧編輯服務生成表格模板
// 2. 回傳表格模板字串
func (e *editorService) GenerateTableTemplateMarkdown(rows, cols int) string {
	if e.smartEditSvc == nil {
		// 簡單的預設表格模板
		return "| 欄位1 | 欄位2 | 欄位3 |\n|-------|-------|-------|\n| 內容1 | 內容2 | 內容3 |"
	}
	return e.smartEditSvc.GenerateTableTemplate(rows, cols)
}

// PreviewMarkdownWithHighlight 預覽 Markdown 內容並包含程式碼高亮
// 參數：content（Markdown 格式的內容）
// 回傳：轉換後的 HTML 字串（包含語法高亮）
//
// 執行流程：
// 1. 使用 goldmark 進行基本的 Markdown 轉換
// 2. 如果有智慧編輯服務，對程式碼區塊進行語法高亮處理
// 3. 回傳增強的 HTML 內容
func (e *editorService) PreviewMarkdownWithHighlight(content string) string {
	// 先進行基本的 Markdown 轉換
	var buf bytes.Buffer
	err := e.markdown.Convert([]byte(content), &buf)
	if err != nil {
		return fmt.Sprintf("<p>Markdown 轉換錯誤: %s</p>", err.Error())
	}
	
	htmlContent := buf.String()
	
	// 如果有智慧編輯服務，進行程式碼高亮處理
	if e.smartEditSvc != nil {
		htmlContent = e.enhanceCodeBlocks(htmlContent)
	}
	
	return htmlContent
}

// enhanceCodeBlocks 增強程式碼區塊的語法高亮
// 參數：htmlContent（HTML 內容）
// 回傳：增強後的 HTML 內容
//
// 執行流程：
// 1. 尋找 HTML 中的程式碼區塊
// 2. 提取程式碼內容和語言資訊
// 3. 使用智慧編輯服務進行語法高亮
// 4. 替換原始的程式碼區塊
func (e *editorService) enhanceCodeBlocks(htmlContent string) string {
	// 這是一個簡化的實作，實際應用中可能需要更複雜的 HTML 解析
	// 目前直接回傳原始內容，未來可以擴展
	return htmlContent
}

// GetSmartEditingService 取得智慧編輯服務實例
// 回傳：SmartEditingService 介面實例
//
// 執行流程：
// 1. 回傳內部的智慧編輯服務實例
// 2. 供外部直接存取智慧編輯功能
func (e *editorService) GetSmartEditingService() SmartEditingService {
	return e.smartEditSvc
}

// SetSmartEditingService 設定智慧編輯服務實例
// 參數：smartEditSvc（智慧編輯服務實例）
//
// 執行流程：
// 1. 更新內部的智慧編輯服務實例
// 2. 允許動態替換智慧編輯服務
func (e *editorService) SetSmartEditingService(smartEditSvc SmartEditingService) {
	e.smartEditSvc = smartEditSvc
}