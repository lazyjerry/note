package repositories

import (
	"fmt"      // 格式化輸出套件
	"io/fs"    // 檔案系統介面套件
	"os"       // 作業系統介面套件
	"path/filepath" // 檔案路徑處理套件
	"strings"  // 字串處理套件
	
	"mac-notebook-app/internal/models" // 引入資料模型
)

// LocalFileRepository 實作 FileRepository 介面
// 提供本地檔案系統的檔案操作功能，包含 Markdown 檔案的特殊處理
type LocalFileRepository struct {
	// 基礎目錄路徑，所有檔案操作都相對於此目錄
	baseDir string
}

// NewLocalFileRepository 建立新的本地檔案儲存庫實例
// 參數：
//   - baseDir: 基礎目錄路徑，用作所有檔案操作的根目錄
// 回傳：指向新建立的 LocalFileRepository 的指標和可能的錯誤
//
// 執行流程：
// 1. 驗證基礎目錄路徑是否有效
// 2. 如果目錄不存在，嘗試建立目錄
// 3. 建立並回傳 LocalFileRepository 實例
func NewLocalFileRepository(baseDir string) (*LocalFileRepository, error) {
	// 驗證基礎目錄路徑不能為空
	if baseDir == "" {
		return nil, models.NewAppError(
			models.ErrValidationFailed,
			"基礎目錄路徑不能為空",
			"請提供有效的目錄路徑",
		)
	}
	
	// 清理並標準化路徑
	cleanPath := filepath.Clean(baseDir)
	
	// 檢查目錄是否存在，如果不存在則建立
	if _, err := os.Stat(cleanPath); os.IsNotExist(err) {
		// 建立目錄，包含所有必要的父目錄
		if err := os.MkdirAll(cleanPath, 0755); err != nil {
			return nil, models.NewAppError(
				models.ErrPermissionDenied,
				"無法建立基礎目錄",
				fmt.Sprintf("目錄路徑：%s，錯誤：%v", cleanPath, err),
			)
		}
	}
	
	return &LocalFileRepository{
		baseDir: cleanPath,
	}, nil
}

// ReadFile 讀取指定路徑的檔案內容
// 參數：path（檔案路徑，相對於基礎目錄）
// 回傳：檔案內容的位元組陣列和可能的錯誤
//
// 執行流程：
// 1. 驗證檔案路徑的安全性
// 2. 建構完整的檔案路徑
// 3. 檢查檔案是否存在
// 4. 讀取檔案內容並回傳
func (r *LocalFileRepository) ReadFile(path string) ([]byte, error) {
	// 驗證路徑安全性
	if err := r.validatePath(path); err != nil {
		return nil, err
	}
	
	// 建構完整的檔案路徑
	fullPath := filepath.Join(r.baseDir, path)
	
	// 檢查檔案是否存在
	if !r.FileExists(path) {
		return nil, models.NewAppError(
			models.ErrFileNotFound,
			"找不到指定的檔案",
			fmt.Sprintf("檔案路徑：%s", fullPath),
		)
	}
	
	// 讀取檔案內容
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, models.NewAppError(
			models.ErrPermissionDenied,
			"無法讀取檔案",
			fmt.Sprintf("檔案路徑：%s，錯誤：%v", fullPath, err),
		)
	}
	
	return data, nil
}

// WriteFile 將資料寫入指定路徑的檔案
// 參數：
//   - path: 檔案路徑（相對於基礎目錄）
//   - data: 要寫入的資料
// 回傳：可能的錯誤
//
// 執行流程：
// 1. 驗證檔案路徑的安全性
// 2. 建構完整的檔案路徑
// 3. 確保父目錄存在
// 4. 寫入檔案內容
func (r *LocalFileRepository) WriteFile(path string, data []byte) error {
	// 驗證路徑安全性
	if err := r.validatePath(path); err != nil {
		return err
	}
	
	// 建構完整的檔案路徑
	fullPath := filepath.Join(r.baseDir, path)
	
	// 確保父目錄存在
	parentDir := filepath.Dir(fullPath)
	if err := os.MkdirAll(parentDir, 0755); err != nil {
		return models.NewAppError(
			models.ErrPermissionDenied,
			"無法建立父目錄",
			fmt.Sprintf("目錄路徑：%s，錯誤：%v", parentDir, err),
		)
	}
	
	// 寫入檔案內容
	if err := os.WriteFile(fullPath, data, 0644); err != nil {
		return models.NewAppError(
			models.ErrSaveFailed,
			"無法寫入檔案",
			fmt.Sprintf("檔案路徑：%s，錯誤：%v", fullPath, err),
		)
	}
	
	return nil
}

// FileExists 檢查指定路徑的檔案是否存在
// 參數：path（檔案路徑，相對於基礎目錄）
// 回傳：檔案是否存在
//
// 執行流程：
// 1. 建構完整的檔案路徑
// 2. 使用 os.Stat 檢查檔案狀態
// 3. 回傳檔案是否存在的結果
func (r *LocalFileRepository) FileExists(path string) bool {
	// 建構完整的檔案路徑
	fullPath := filepath.Join(r.baseDir, path)
	
	// 檢查檔案是否存在
	_, err := os.Stat(fullPath)
	return !os.IsNotExist(err)
}

// DeleteFile 刪除指定路徑的檔案
// 參數：path（檔案路徑，相對於基礎目錄）
// 回傳：可能的錯誤
//
// 執行流程：
// 1. 驗證檔案路徑的安全性
// 2. 檢查檔案是否存在
// 3. 刪除檔案
func (r *LocalFileRepository) DeleteFile(path string) error {
	// 驗證路徑安全性
	if err := r.validatePath(path); err != nil {
		return err
	}
	
	// 建構完整的檔案路徑
	fullPath := filepath.Join(r.baseDir, path)
	
	// 檢查檔案是否存在
	if !r.FileExists(path) {
		return models.NewAppError(
			models.ErrFileNotFound,
			"找不到要刪除的檔案",
			fmt.Sprintf("檔案路徑：%s", fullPath),
		)
	}
	
	// 刪除檔案
	if err := os.Remove(fullPath); err != nil {
		return models.NewAppError(
			models.ErrPermissionDenied,
			"無法刪除檔案",
			fmt.Sprintf("檔案路徑：%s，錯誤：%v", fullPath, err),
		)
	}
	
	return nil
}

// CreateDirectory 建立指定路徑的目錄
// 參數：path（目錄路徑，相對於基礎目錄）
// 回傳：可能的錯誤
//
// 執行流程：
// 1. 驗證目錄路徑的安全性
// 2. 建構完整的目錄路徑
// 3. 建立目錄（包含所有必要的父目錄）
func (r *LocalFileRepository) CreateDirectory(path string) error {
	// 驗證路徑安全性
	if err := r.validatePath(path); err != nil {
		return err
	}
	
	// 建構完整的目錄路徑
	fullPath := filepath.Join(r.baseDir, path)
	
	// 建立目錄，包含所有必要的父目錄
	if err := os.MkdirAll(fullPath, 0755); err != nil {
		return models.NewAppError(
			models.ErrPermissionDenied,
			"無法建立目錄",
			fmt.Sprintf("目錄路徑：%s，錯誤：%v", fullPath, err),
		)
	}
	
	return nil
}// ListD
irectory 列出指定目錄中的檔案和子目錄
// 參數：path（目錄路徑，相對於基礎目錄）
// 回傳：檔案資訊陣列和可能的錯誤
//
// 執行流程：
// 1. 驗證目錄路徑的安全性
// 2. 建構完整的目錄路徑
// 3. 檢查目錄是否存在
// 4. 讀取目錄內容
// 5. 為每個項目建立 FileInfo 實例
func (r *LocalFileRepository) ListDirectory(path string) ([]*models.FileInfo, error) {
	// 驗證路徑安全性
	if err := r.validatePath(path); err != nil {
		return nil, err
	}
	
	// 建構完整的目錄路徑
	fullPath := filepath.Join(r.baseDir, path)
	
	// 檢查目錄是否存在
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return nil, models.NewAppError(
			models.ErrFileNotFound,
			"找不到指定的目錄",
			fmt.Sprintf("目錄路徑：%s", fullPath),
		)
	}
	
	// 讀取目錄內容
	entries, err := os.ReadDir(fullPath)
	if err != nil {
		return nil, models.NewAppError(
			models.ErrPermissionDenied,
			"無法讀取目錄內容",
			fmt.Sprintf("目錄路徑：%s，錯誤：%v", fullPath, err),
		)
	}
	
	// 建立檔案資訊陣列
	var fileInfos []*models.FileInfo
	
	// 遍歷目錄項目，為每個項目建立 FileInfo
	for _, entry := range entries {
		// 取得詳細的檔案資訊
		info, err := entry.Info()
		if err != nil {
			// 如果無法取得檔案資訊，跳過此項目並記錄錯誤
			continue
		}
		
		// 建構項目的完整路徑
		itemPath := filepath.Join(path, entry.Name())
		
		// 建立 FileInfo 實例並加入陣列
		fileInfo := models.NewFileInfo(entry.Name(), itemPath, info)
		fileInfos = append(fileInfos, fileInfo)
	}
	
	return fileInfos, nil
}

// validatePath 驗證檔案路徑的安全性
// 參數：path（要驗證的路徑）
// 回傳：可能的錯誤
//
// 安全性檢查：
// 1. 路徑不能為空
// 2. 路徑不能包含 ".." 以防止目錄遍歷攻擊
// 3. 路徑不能是絕對路徑
// 4. 路徑不能包含危險字符
func (r *LocalFileRepository) validatePath(path string) error {
	// 檢查路徑是否為空
	if path == "" {
		return models.NewAppError(
			models.ErrValidationFailed,
			"檔案路徑不能為空",
			"請提供有效的檔案路徑",
		)
	}
	
	// 清理路徑
	cleanPath := filepath.Clean(path)
	
	// 檢查是否為絕對路徑
	if filepath.IsAbs(cleanPath) {
		return models.NewAppError(
			models.ErrValidationFailed,
			"不允許使用絕對路徑",
			fmt.Sprintf("路徑：%s", path),
		)
	}
	
	// 檢查路徑是否包含 ".." 以防止目錄遍歷攻擊
	if strings.Contains(cleanPath, "..") {
		return models.NewAppError(
			models.ErrValidationFailed,
			"路徑不能包含 '..' 字符",
			fmt.Sprintf("路徑：%s", path),
		)
	}
	
	// 檢查路徑是否嘗試跳出基礎目錄
	fullPath := filepath.Join(r.baseDir, cleanPath)
	if !strings.HasPrefix(fullPath, r.baseDir) {
		return models.NewAppError(
			models.ErrValidationFailed,
			"路徑超出允許的範圍",
			fmt.Sprintf("路徑：%s", path),
		)
	}
	
	return nil
}

// IsMarkdownFile 檢查指定路徑的檔案是否為 Markdown 檔案
// 參數：path（檔案路徑，相對於基礎目錄）
// 回傳：是否為 Markdown 檔案和可能的錯誤
//
// 執行流程：
// 1. 驗證檔案路徑
// 2. 檢查檔案是否存在
// 3. 根據檔案副檔名判斷是否為 Markdown 檔案
func (r *LocalFileRepository) IsMarkdownFile(path string) (bool, error) {
	// 驗證路徑安全性
	if err := r.validatePath(path); err != nil {
		return false, err
	}
	
	// 檢查檔案是否存在
	if !r.FileExists(path) {
		return false, models.NewAppError(
			models.ErrFileNotFound,
			"找不到指定的檔案",
			fmt.Sprintf("檔案路徑：%s", path),
		)
	}
	
	// 取得檔案名稱
	fileName := filepath.Base(path)
	
	// 檢查是否為加密的 Markdown 檔案
	if strings.HasSuffix(fileName, ".md.enc") {
		return true, nil
	}
	
	// 檢查是否為一般的 Markdown 檔案
	if strings.HasSuffix(fileName, ".md") {
		return true, nil
	}
	
	return false, nil
}

// ReadMarkdownFile 讀取 Markdown 檔案內容並回傳字串格式
// 參數：path（檔案路徑，相對於基礎目錄）
// 回傳：檔案內容字串和可能的錯誤
//
// 執行流程：
// 1. 驗證檔案是否為 Markdown 檔案
// 2. 讀取檔案內容
// 3. 將位元組陣列轉換為字串並回傳
func (r *LocalFileRepository) ReadMarkdownFile(path string) (string, error) {
	// 檢查是否為 Markdown 檔案
	isMarkdown, err := r.IsMarkdownFile(path)
	if err != nil {
		return "", err
	}
	
	if !isMarkdown {
		return "", models.NewAppError(
			models.ErrValidationFailed,
			"指定的檔案不是 Markdown 檔案",
			fmt.Sprintf("檔案路徑：%s", path),
		)
	}
	
	// 讀取檔案內容
	data, err := r.ReadFile(path)
	if err != nil {
		return "", err
	}
	
	// 將位元組陣列轉換為字串
	return string(data), nil
}

// WriteMarkdownFile 將字串內容寫入 Markdown 檔案
// 參數：
//   - path: 檔案路徑（相對於基礎目錄）
//   - content: 要寫入的 Markdown 內容
// 回傳：可能的錯誤
//
// 執行流程：
// 1. 驗證檔案路徑
// 2. 確保檔案具有正確的 Markdown 副檔名
// 3. 將字串內容轉換為位元組陣列並寫入檔案
func (r *LocalFileRepository) WriteMarkdownFile(path, content string) error {
	// 驗證路徑安全性
	if err := r.validatePath(path); err != nil {
		return err
	}
	
	// 確保檔案具有 .md 副檔名（如果不是加密檔案）
	if !strings.HasSuffix(path, ".md") && !strings.HasSuffix(path, ".md.enc") {
		return models.NewAppError(
			models.ErrValidationFailed,
			"Markdown 檔案必須具有 .md 或 .md.enc 副檔名",
			fmt.Sprintf("檔案路徑：%s", path),
		)
	}
	
	// 將字串內容轉換為位元組陣列並寫入檔案
	return r.WriteFile(path, []byte(content))
}

// GetBaseDirectory 取得基礎目錄路徑
// 回傳：基礎目錄的完整路徑
//
// 此方法用於取得儲存庫的基礎目錄路徑，
// 可用於除錯或顯示目前工作目錄資訊
func (r *LocalFileRepository) GetBaseDirectory() string {
	return r.baseDir
}

// WalkDirectory 遞迴遍歷目錄樹
// 參數：
//   - startPath: 開始遍歷的路徑（相對於基礎目錄）
//   - walkFunc: 對每個檔案/目錄執行的函數
// 回傳：可能的錯誤
//
// 執行流程：
// 1. 驗證起始路徑
// 2. 使用 filepath.WalkDir 遞迴遍歷目錄
// 3. 對每個項目執行指定的函數
func (r *LocalFileRepository) WalkDirectory(startPath string, walkFunc func(*models.FileInfo) error) error {
	// 驗證路徑安全性
	if err := r.validatePath(startPath); err != nil {
		return err
	}
	
	// 建構完整的起始路徑
	fullStartPath := filepath.Join(r.baseDir, startPath)
	
	// 遞迴遍歷目錄
	return filepath.WalkDir(fullStartPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		
		// 取得相對於基礎目錄的路徑
		relPath, err := filepath.Rel(r.baseDir, path)
		if err != nil {
			return err
		}
		
		// 取得檔案資訊
		info, err := d.Info()
		if err != nil {
			return err
		}
		
		// 建立 FileInfo 實例
		fileInfo := models.NewFileInfo(d.Name(), relPath, info)
		
		// 執行使用者提供的函數
		return walkFunc(fileInfo)
	})
}