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
type editorService struct {
	fileRepo      repositories.FileRepository // 檔案存取介面
	encryptionSvc EncryptionService           // 加密服務介面
	passwordSvc   PasswordService             // 密碼服務介面
	biometricSvc  BiometricService            // 生物識別服務介面
	markdown      goldmark.Markdown           // Markdown 解析器實例
	activeNotes   map[string]*models.Note     // 當前開啟的筆記快取
}

// NewEditorService 建立新的編輯器服務實例
// 參數：
//   - fileRepo: 檔案存取介面
//   - encryptionSvc: 加密服務介面
//   - passwordSvc: 密碼服務介面
//   - biometricSvc: 生物識別服務介面
// 回傳：EditorService 介面實例
//
// 執行流程：
// 1. 初始化 goldmark Markdown 解析器，啟用常用擴展功能
// 2. 建立筆記快取映射表
// 3. 整合加密相關服務
// 4. 回傳配置完成的編輯器服務實例
func NewEditorService(fileRepo repositories.FileRepository, encryptionSvc EncryptionService, passwordSvc PasswordService, biometricSvc BiometricService) EditorService {
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

	return &editorService{
		fileRepo:      fileRepo,
		encryptionSvc: encryptionSvc,
		passwordSvc:   passwordSvc,
		biometricSvc:  biometricSvc,
		markdown:      md,
		activeNotes:   make(map[string]*models.Note),
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