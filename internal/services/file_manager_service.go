package services

import (
	"fmt"      // 格式化輸出套件
	"os"       // 作業系統介面套件
	"path/filepath" // 檔案路徑處理套件
	"strings"  // 字串處理套件
	
	"mac-notebook-app/internal/models"      // 引入資料模型
	"mac-notebook-app/internal/repositories" // 引入儲存庫介面
)

// LocalFileManagerService 實作 FileManagerService 介面
// 提供本地檔案系統的檔案和目錄管理功能，包含 CRUD 操作和檔案樹遍歷
type LocalFileManagerService struct {
	// fileRepo 檔案儲存庫，用於執行底層檔案操作
	fileRepo repositories.FileRepository
	
	// baseDir 基礎工作目錄，所有操作都相對於此目錄
	baseDir string
}

// NewLocalFileManagerService 建立新的本地檔案管理服務實例
// 參數：
//   - fileRepo: 檔案儲存庫介面實例
//   - baseDir: 基礎工作目錄路徑
// 回傳：指向新建立的 LocalFileManagerService 的指標和可能的錯誤
//
// 執行流程：
// 1. 驗證輸入參數的有效性
// 2. 檢查基礎目錄是否存在且可存取
// 3. 建立並回傳服務實例
func NewLocalFileManagerService(fileRepo repositories.FileRepository, baseDir string) (*LocalFileManagerService, error) {
	// 驗證檔案儲存庫不能為 nil
	if fileRepo == nil {
		return nil, models.NewAppError(
			models.ErrValidationFailed,
			"檔案儲存庫不能為空",
			"請提供有效的檔案儲存庫實例",
		)
	}
	
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
	
	// 檢查目錄是否存在
	if _, err := os.Stat(cleanPath); os.IsNotExist(err) {
		return nil, models.NewAppError(
			models.ErrFileNotFound,
			"指定的基礎目錄不存在",
			fmt.Sprintf("目錄路徑：%s", cleanPath),
		)
	}
	
	return &LocalFileManagerService{
		fileRepo: fileRepo,
		baseDir:  cleanPath,
	}, nil
}

// ListFiles 列出指定目錄中的檔案和子目錄
// 參數：directory（目錄路徑，相對於基礎目錄）
// 回傳：檔案資訊陣列和可能的錯誤
//
// 執行流程：
// 1. 驗證目錄路徑的有效性
// 2. 使用檔案儲存庫列出目錄內容
// 3. 對結果進行排序和過濾
// 4. 回傳檔案資訊陣列
func (s *LocalFileManagerService) ListFiles(directory string) ([]*models.FileInfo, error) {
	// 驗證目錄路徑
	if err := s.validatePath(directory); err != nil {
		return nil, err
	}
	
	// 使用檔案儲存庫列出目錄內容
	fileInfos, err := s.fileRepo.ListDirectory(directory)
	if err != nil {
		return nil, err
	}
	
	// 對檔案資訊進行排序：目錄優先，然後按名稱排序
	s.sortFileInfos(fileInfos)
	
	return fileInfos, nil
}

// CreateDirectory 建立新目錄
// 參數：path（目錄路徑，相對於基礎目錄）
// 回傳：可能的錯誤
//
// 執行流程：
// 1. 驗證目錄路徑的有效性
// 2. 檢查目錄是否已存在
// 3. 使用檔案儲存庫建立目錄
func (s *LocalFileManagerService) CreateDirectory(path string) error {
	// 驗證路徑
	if err := s.validatePath(path); err != nil {
		return err
	}
	
	// 檢查目錄是否已存在
	if s.fileRepo.FileExists(path) {
		return models.NewAppError(
			models.ErrValidationFailed,
			"目錄已存在",
			fmt.Sprintf("目錄路徑：%s", path),
		)
	}
	
	// 建立目錄
	return s.fileRepo.CreateDirectory(path)
}

// DeleteFile 刪除檔案或目錄
// 參數：path（檔案或目錄路徑，相對於基礎目錄）
// 回傳：可能的錯誤
//
// 執行流程：
// 1. 驗證路徑的有效性
// 2. 檢查檔案或目錄是否存在
// 3. 如果是目錄，檢查是否為空目錄
// 4. 執行刪除操作
func (s *LocalFileManagerService) DeleteFile(path string) error {
	// 驗證路徑
	if err := s.validatePath(path); err != nil {
		return err
	}
	
	// 檢查檔案或目錄是否存在
	if !s.fileRepo.FileExists(path) {
		return models.NewAppError(
			models.ErrFileNotFound,
			"找不到要刪除的檔案或目錄",
			fmt.Sprintf("路徑：%s", path),
		)
	}
	
	// 建構完整路徑以檢查是否為目錄
	fullPath := filepath.Join(s.baseDir, path)
	fileInfo, err := os.Stat(fullPath)
	if err != nil {
		return models.NewAppError(
			models.ErrPermissionDenied,
			"無法存取檔案或目錄",
			fmt.Sprintf("路徑：%s，錯誤：%v", path, err),
		)
	}
	
	// 如果是目錄，檢查是否為空
	if fileInfo.IsDir() {
		isEmpty, err := s.isDirectoryEmpty(path)
		if err != nil {
			return err
		}
		
		if !isEmpty {
			return models.NewAppError(
				models.ErrValidationFailed,
				"無法刪除非空目錄",
				fmt.Sprintf("目錄路徑：%s 包含檔案或子目錄", path),
			)
		}
		
		// 刪除空目錄
		if err := os.Remove(fullPath); err != nil {
			return models.NewAppError(
				models.ErrPermissionDenied,
				"無法刪除目錄",
				fmt.Sprintf("目錄路徑：%s，錯誤：%v", path, err),
			)
		}
		
		return nil
	}
	
	// 刪除檔案
	return s.fileRepo.DeleteFile(path)
}

// RenameFile 重新命名檔案或目錄
// 參數：
//   - oldPath: 舊路徑（相對於基礎目錄）
//   - newPath: 新路徑（相對於基礎目錄）
// 回傳：可能的錯誤
//
// 執行流程：
// 1. 驗證舊路徑和新路徑的有效性
// 2. 檢查舊檔案是否存在
// 3. 檢查新路徑是否已被佔用
// 4. 執行重新命名操作
func (s *LocalFileManagerService) RenameFile(oldPath, newPath string) error {
	// 驗證路徑
	if err := s.validatePath(oldPath); err != nil {
		return fmt.Errorf("舊路徑無效：%w", err)
	}
	
	if err := s.validatePath(newPath); err != nil {
		return fmt.Errorf("新路徑無效：%w", err)
	}
	
	// 檢查舊檔案是否存在
	if !s.fileRepo.FileExists(oldPath) {
		return models.NewAppError(
			models.ErrFileNotFound,
			"找不到要重新命名的檔案或目錄",
			fmt.Sprintf("路徑：%s", oldPath),
		)
	}
	
	// 檢查新路徑是否已被佔用
	if s.fileRepo.FileExists(newPath) {
		return models.NewAppError(
			models.ErrValidationFailed,
			"目標路徑已存在",
			fmt.Sprintf("新路徑：%s", newPath),
		)
	}
	
	// 建構完整路徑
	fullOldPath := filepath.Join(s.baseDir, oldPath)
	fullNewPath := filepath.Join(s.baseDir, newPath)
	
	// 確保新路徑的父目錄存在
	newParentDir := filepath.Dir(fullNewPath)
	if err := os.MkdirAll(newParentDir, 0755); err != nil {
		return models.NewAppError(
			models.ErrPermissionDenied,
			"無法建立目標目錄",
			fmt.Sprintf("目錄路徑：%s，錯誤：%v", newParentDir, err),
		)
	}
	
	// 執行重新命名操作
	if err := os.Rename(fullOldPath, fullNewPath); err != nil {
		return models.NewAppError(
			models.ErrPermissionDenied,
			"無法重新命名檔案或目錄",
			fmt.Sprintf("從 %s 到 %s，錯誤：%v", oldPath, newPath, err),
		)
	}
	
	return nil
}

// MoveFile 移動檔案或目錄到新位置
// 參數：
//   - sourcePath: 來源路徑（相對於基礎目錄）
//   - destPath: 目標路徑（相對於基礎目錄）
// 回傳：可能的錯誤
//
// 執行流程：
// 1. 驗證來源路徑和目標路徑的有效性
// 2. 檢查來源檔案是否存在
// 3. 處理目標路徑（如果是目錄，則移動到該目錄內）
// 4. 執行移動操作
func (s *LocalFileManagerService) MoveFile(sourcePath, destPath string) error {
	// 驗證路徑
	if err := s.validatePath(sourcePath); err != nil {
		return fmt.Errorf("來源路徑無效：%w", err)
	}
	
	if err := s.validatePath(destPath); err != nil {
		return fmt.Errorf("目標路徑無效：%w", err)
	}
	
	// 檢查來源檔案是否存在
	if !s.fileRepo.FileExists(sourcePath) {
		return models.NewAppError(
			models.ErrFileNotFound,
			"找不到要移動的檔案或目錄",
			fmt.Sprintf("來源路徑：%s", sourcePath),
		)
	}
	
	// 建構完整路徑
	fullSourcePath := filepath.Join(s.baseDir, sourcePath)
	fullDestPath := filepath.Join(s.baseDir, destPath)
	
	// 檢查目標路徑是否為目錄
	if s.fileRepo.FileExists(destPath) {
		destInfo, err := os.Stat(fullDestPath)
		if err != nil {
			return models.NewAppError(
				models.ErrPermissionDenied,
				"無法存取目標路徑",
				fmt.Sprintf("目標路徑：%s，錯誤：%v", destPath, err),
			)
		}
		
		// 如果目標是目錄，將檔案移動到該目錄內
		if destInfo.IsDir() {
			fileName := filepath.Base(sourcePath)
			fullDestPath = filepath.Join(fullDestPath, fileName)
			destPath = filepath.Join(destPath, fileName)
		}
	}
	
	// 檢查最終目標路徑是否已存在
	if _, err := os.Stat(fullDestPath); !os.IsNotExist(err) {
		return models.NewAppError(
			models.ErrValidationFailed,
			"目標位置已存在檔案或目錄",
			fmt.Sprintf("目標路徑：%s", destPath),
		)
	}
	
	// 確保目標目錄存在
	destDir := filepath.Dir(fullDestPath)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return models.NewAppError(
			models.ErrPermissionDenied,
			"無法建立目標目錄",
			fmt.Sprintf("目錄路徑：%s，錯誤：%v", destDir, err),
		)
	}
	
	// 執行移動操作
	if err := os.Rename(fullSourcePath, fullDestPath); err != nil {
		return models.NewAppError(
			models.ErrPermissionDenied,
			"無法移動檔案或目錄",
			fmt.Sprintf("從 %s 到 %s，錯誤：%v", sourcePath, destPath, err),
		)
	}
	
	return nil
}

// GetFileTree 取得完整的檔案樹狀結構
// 參數：rootPath（根目錄路徑，相對於基礎目錄）
// 回傳：檔案樹節點和可能的錯誤
//
// 執行流程：
// 1. 驗證根目錄路徑
// 2. 遞迴遍歷目錄結構
// 3. 建立樹狀結構的節點
// 4. 回傳完整的檔案樹
func (s *LocalFileManagerService) GetFileTree(rootPath string) (*FileTreeNode, error) {
	// 驗證路徑
	if err := s.validatePath(rootPath); err != nil {
		return nil, err
	}
	
	// 檢查根目錄是否存在
	if !s.fileRepo.FileExists(rootPath) {
		return nil, models.NewAppError(
			models.ErrFileNotFound,
			"找不到指定的根目錄",
			fmt.Sprintf("路徑：%s", rootPath),
		)
	}
	
	// 建立根節點
	rootNode, err := s.createFileTreeNode(rootPath)
	if err != nil {
		return nil, err
	}
	
	// 如果根節點是目錄，遞迴建立子節點
	if rootNode.IsDirectory {
		children, err := s.buildFileTreeChildren(rootPath)
		if err != nil {
			return nil, err
		}
		rootNode.Children = children
	}
	
	return rootNode, nil
}

// SearchFiles 在指定目錄中搜尋檔案
// 參數：
//   - searchPath: 搜尋路徑（相對於基礎目錄）
//   - pattern: 搜尋模式（支援萬用字元）
//   - includeSubdirs: 是否包含子目錄
// 回傳：符合條件的檔案資訊陣列和可能的錯誤
//
// 執行流程：
// 1. 驗證搜尋路徑和模式
// 2. 遍歷指定目錄（可選擇是否包含子目錄）
// 3. 使用模式匹配過濾檔案
// 4. 回傳符合條件的檔案列表
func (s *LocalFileManagerService) SearchFiles(searchPath, pattern string, includeSubdirs bool) ([]*models.FileInfo, error) {
	// 驗證搜尋路徑
	if err := s.validatePath(searchPath); err != nil {
		return nil, err
	}
	
	// 驗證搜尋模式
	if pattern == "" {
		return nil, models.NewAppError(
			models.ErrValidationFailed,
			"搜尋模式不能為空",
			"請提供有效的搜尋模式",
		)
	}
	
	var matchedFiles []*models.FileInfo
	
	if includeSubdirs {
		// 遞迴搜尋子目錄
		err := s.fileRepo.WalkDirectory(searchPath, func(info *models.FileInfo) error {
			if !info.IsDirectory {
				matched, err := filepath.Match(pattern, info.Name)
				if err != nil {
					return err
				}
				if matched {
					matchedFiles = append(matchedFiles, info)
				}
			}
			return nil
		})
		
		if err != nil {
			return nil, models.NewAppError(
				models.ErrPermissionDenied,
				"搜尋檔案時發生錯誤",
				fmt.Sprintf("搜尋路徑：%s，錯誤：%v", searchPath, err),
			)
		}
	} else {
		// 只搜尋當前目錄
		fileInfos, err := s.fileRepo.ListDirectory(searchPath)
		if err != nil {
			return nil, err
		}
		
		for _, info := range fileInfos {
			if !info.IsDirectory {
				matched, err := filepath.Match(pattern, info.Name)
				if err != nil {
					return nil, models.NewAppError(
						models.ErrValidationFailed,
						"搜尋模式無效",
						fmt.Sprintf("模式：%s，錯誤：%v", pattern, err),
					)
				}
				if matched {
					matchedFiles = append(matchedFiles, info)
				}
			}
		}
	}
	
	return matchedFiles, nil
}

// CopyFile 複製檔案或目錄
// 參數：
//   - sourcePath: 來源路徑（相對於基礎目錄）
//   - destPath: 目標路徑（相對於基礎目錄）
// 回傳：可能的錯誤
//
// 執行流程：
// 1. 驗證來源和目標路徑
// 2. 檢查來源檔案是否存在
// 3. 根據檔案類型執行複製操作
// 4. 處理目錄的遞迴複製
func (s *LocalFileManagerService) CopyFile(sourcePath, destPath string) error {
	// 驗證路徑
	if err := s.validatePath(sourcePath); err != nil {
		return fmt.Errorf("來源路徑無效：%w", err)
	}
	
	if err := s.validatePath(destPath); err != nil {
		return fmt.Errorf("目標路徑無效：%w", err)
	}
	
	// 檢查來源檔案是否存在
	if !s.fileRepo.FileExists(sourcePath) {
		return models.NewAppError(
			models.ErrFileNotFound,
			"找不到要複製的檔案或目錄",
			fmt.Sprintf("來源路徑：%s", sourcePath),
		)
	}
	
	// 建構完整路徑
	fullSourcePath := filepath.Join(s.baseDir, sourcePath)
	fullDestPath := filepath.Join(s.baseDir, destPath)
	
	// 取得來源檔案資訊
	sourceInfo, err := os.Stat(fullSourcePath)
	if err != nil {
		return models.NewAppError(
			models.ErrPermissionDenied,
			"無法存取來源檔案",
			fmt.Sprintf("來源路徑：%s，錯誤：%v", sourcePath, err),
		)
	}
	
	// 檢查目標路徑是否已存在
	if _, err := os.Stat(fullDestPath); !os.IsNotExist(err) {
		return models.NewAppError(
			models.ErrValidationFailed,
			"目標路徑已存在",
			fmt.Sprintf("目標路徑：%s", destPath),
		)
	}
	
	if sourceInfo.IsDir() {
		// 複製目錄
		return s.copyDirectory(sourcePath, destPath)
	} else {
		// 複製檔案
		return s.copyFile(sourcePath, destPath)
	}
}

// GetDirectorySize 計算目錄的總大小
// 參數：dirPath（目錄路徑，相對於基礎目錄）
// 回傳：目錄大小（位元組）和可能的錯誤
//
// 執行流程：
// 1. 驗證目錄路徑
// 2. 遞迴遍歷目錄中的所有檔案
// 3. 累計所有檔案的大小
// 4. 回傳總大小
func (s *LocalFileManagerService) GetDirectorySize(dirPath string) (int64, error) {
	// 驗證路徑
	if err := s.validatePath(dirPath); err != nil {
		return 0, err
	}
	
	// 檢查目錄是否存在
	if !s.fileRepo.FileExists(dirPath) {
		return 0, models.NewAppError(
			models.ErrFileNotFound,
			"找不到指定的目錄",
			fmt.Sprintf("目錄路徑：%s", dirPath),
		)
	}
	
	var totalSize int64
	
	// 遞迴遍歷目錄計算大小
	err := s.fileRepo.WalkDirectory(dirPath, func(info *models.FileInfo) error {
		if !info.IsDirectory {
			totalSize += info.Size
		}
		return nil
	})
	
	if err != nil {
		return 0, models.NewAppError(
			models.ErrPermissionDenied,
			"計算目錄大小時發生錯誤",
			fmt.Sprintf("目錄路徑：%s，錯誤：%v", dirPath, err),
		)
	}
	
	return totalSize, nil
}

// validatePath 驗證檔案路徑的安全性和有效性
// 參數：path（要驗證的路徑）
// 回傳：可能的錯誤
//
// 安全性檢查：
// 1. 路徑不能為空
// 2. 路徑不能包含 ".." 以防止目錄遍歷攻擊
// 3. 路徑不能是絕對路徑
// 4. 路徑不能包含危險字符
func (s *LocalFileManagerService) validatePath(path string) error {
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
	fullPath := filepath.Join(s.baseDir, cleanPath)
	if !strings.HasPrefix(fullPath, s.baseDir) {
		return models.NewAppError(
			models.ErrValidationFailed,
			"路徑超出允許的範圍",
			fmt.Sprintf("路徑：%s", path),
		)
	}
	
	return nil
}

// isDirectoryEmpty 檢查目錄是否為空
// 參數：dirPath（目錄路徑，相對於基礎目錄）
// 回傳：目錄是否為空和可能的錯誤
func (s *LocalFileManagerService) isDirectoryEmpty(dirPath string) (bool, error) {
	fileInfos, err := s.fileRepo.ListDirectory(dirPath)
	if err != nil {
		return false, err
	}
	
	return len(fileInfos) == 0, nil
}

// sortFileInfos 對檔案資訊陣列進行排序
// 排序規則：目錄優先，然後按名稱字母順序排序
// 參數：fileInfos（要排序的檔案資訊陣列）
func (s *LocalFileManagerService) sortFileInfos(fileInfos []*models.FileInfo) {
	// 使用簡單的冒泡排序演算法
	n := len(fileInfos)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			// 目錄優先排序
			if !fileInfos[j].IsDirectory && fileInfos[j+1].IsDirectory {
				fileInfos[j], fileInfos[j+1] = fileInfos[j+1], fileInfos[j]
			} else if fileInfos[j].IsDirectory == fileInfos[j+1].IsDirectory {
				// 同類型按名稱排序
				if strings.ToLower(fileInfos[j].Name) > strings.ToLower(fileInfos[j+1].Name) {
					fileInfos[j], fileInfos[j+1] = fileInfos[j+1], fileInfos[j]
				}
			}
		}
	}
}

// createFileTreeNode 為指定路徑建立檔案樹節點
// 參數：path（檔案或目錄路徑）
// 回傳：檔案樹節點和可能的錯誤
func (s *LocalFileManagerService) createFileTreeNode(path string) (*FileTreeNode, error) {
	// 建構完整路徑
	fullPath := filepath.Join(s.baseDir, path)
	
	// 取得檔案資訊
	fileInfo, err := os.Stat(fullPath)
	if err != nil {
		return nil, models.NewAppError(
			models.ErrPermissionDenied,
			"無法存取檔案或目錄",
			fmt.Sprintf("路徑：%s，錯誤：%v", path, err),
		)
	}
	
	// 建立檔案資訊實例
	modelFileInfo := models.NewFileInfo(filepath.Base(path), path, fileInfo)
	
	// 建立檔案樹節點
	return &FileTreeNode{
		FileInfo:    modelFileInfo,
		IsDirectory: fileInfo.IsDir(),
		Children:    nil, // 子節點將在需要時建立
	}, nil
}

// buildFileTreeChildren 為目錄建立子節點
// 參數：dirPath（目錄路徑）
// 回傳：子節點陣列和可能的錯誤
func (s *LocalFileManagerService) buildFileTreeChildren(dirPath string) ([]*FileTreeNode, error) {
	// 列出目錄內容
	fileInfos, err := s.fileRepo.ListDirectory(dirPath)
	if err != nil {
		return nil, err
	}
	
	var children []*FileTreeNode
	
	// 為每個項目建立子節點
	for _, info := range fileInfos {
		childNode := &FileTreeNode{
			FileInfo:    info,
			IsDirectory: info.IsDirectory,
			Children:    nil,
		}
		
		// 如果是目錄，遞迴建立子節點
		if info.IsDirectory {
			grandChildren, err := s.buildFileTreeChildren(info.Path)
			if err != nil {
				return nil, err
			}
			childNode.Children = grandChildren
		}
		
		children = append(children, childNode)
	}
	
	return children, nil
}

// copyFile 複製單個檔案
// 參數：sourcePath、destPath（來源和目標路徑）
// 回傳：可能的錯誤
func (s *LocalFileManagerService) copyFile(sourcePath, destPath string) error {
	// 讀取來源檔案內容
	data, err := s.fileRepo.ReadFile(sourcePath)
	if err != nil {
		return err
	}
	
	// 寫入目標檔案
	return s.fileRepo.WriteFile(destPath, data)
}

// copyDirectory 遞迴複製目錄
// 參數：sourcePath、destPath（來源和目標路徑）
// 回傳：可能的錯誤
func (s *LocalFileManagerService) copyDirectory(sourcePath, destPath string) error {
	// 建立目標目錄
	if err := s.fileRepo.CreateDirectory(destPath); err != nil {
		return err
	}
	
	// 列出來源目錄內容
	fileInfos, err := s.fileRepo.ListDirectory(sourcePath)
	if err != nil {
		return err
	}
	
	// 遞迴複製每個項目
	for _, info := range fileInfos {
		sourceItemPath := info.Path
		destItemPath := filepath.Join(destPath, info.Name)
		
		if info.IsDirectory {
			// 遞迴複製子目錄
			if err := s.copyDirectory(sourceItemPath, destItemPath); err != nil {
				return err
			}
		} else {
			// 複製檔案
			if err := s.copyFile(sourceItemPath, destItemPath); err != nil {
				return err
			}
		}
	}
	
	return nil
}

// FileTreeNode 代表檔案樹中的一個節點
// 用於表示檔案系統的樹狀結構
type FileTreeNode struct {
	FileInfo    *models.FileInfo `json:"file_info"`    // 檔案資訊
	IsDirectory bool             `json:"is_directory"` // 是否為目錄
	Children    []*FileTreeNode  `json:"children"`     // 子節點陣列
}

// GetPath 取得節點的完整路徑
// 回傳：節點的路徑字串
func (n *FileTreeNode) GetPath() string {
	return n.FileInfo.Path
}

// GetName 取得節點的名稱
// 回傳：節點的名稱字串
func (n *FileTreeNode) GetName() string {
	return n.FileInfo.Name
}

// HasChildren 檢查節點是否有子節點
// 回傳：是否有子節點
func (n *FileTreeNode) HasChildren() bool {
	return len(n.Children) > 0
}

// GetChildCount 取得子節點數量
// 回傳：子節點的數量
func (n *FileTreeNode) GetChildCount() int {
	return len(n.Children)
}