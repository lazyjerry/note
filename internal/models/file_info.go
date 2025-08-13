package models

import (
	"os"   // 作業系統介面，用於檔案系統操作
	"time" // 時間處理套件
)

// FileInfo 代表檔案系統資訊
// 封裝了檔案或目錄的基本屬性，包含加密狀態等擴充資訊
type FileInfo struct {
	Name         string      `json:"name"`         // 檔案或目錄名稱
	Path         string      `json:"path"`         // 完整的檔案路徑
	IsDirectory  bool        `json:"is_directory"` // 是否為目錄
	Size         int64       `json:"size"`         // 檔案大小（位元組）
	ModTime      time.Time   `json:"mod_time"`     // 最後修改時間
	IsEncrypted  bool        `json:"is_encrypted"` // 是否為加密檔案
	Permissions  os.FileMode `json:"permissions"`  // 檔案權限
}

// NewFileInfo 從 os.FileInfo 建立 FileInfo 實例
// 參數：
//   - name: 檔案或目錄名稱
//   - path: 完整的檔案路徑
//   - info: 系統提供的檔案資訊
// 回傳：指向新建立的 FileInfo 的指標
//
// 執行流程：
// 1. 檢查檔案是否為加密檔案（以 .enc 結尾）
// 2. 從 os.FileInfo 提取基本檔案屬性
// 3. 建立並回傳 FileInfo 實例
func NewFileInfo(name, path string, info os.FileInfo) *FileInfo {
	// 判斷是否為加密檔案
	// 加密檔案的判斷條件：非目錄且檔名以 .enc 結尾
	isEncrypted := false
	if !info.IsDir() && len(name) > 4 && name[len(name)-4:] == ".enc" {
		isEncrypted = true
	}
	
	return &FileInfo{
		Name:        name,           // 設定檔案名稱
		Path:        path,           // 設定檔案路徑
		IsDirectory: info.IsDir(),   // 設定是否為目錄
		Size:        info.Size(),    // 設定檔案大小
		ModTime:     info.ModTime(), // 設定修改時間
		IsEncrypted: isEncrypted,    // 設定加密狀態
		Permissions: info.Mode(),    // 設定檔案權限
	}
}

// IsMarkdownFile 檢查檔案是否為 Markdown 檔案
// 回傳：如果是 Markdown 檔案則回傳 true，否則回傳 false
//
// 判斷邏輯：
// 1. 如果是目錄，直接回傳 false
// 2. 如果是加密檔案，先移除 .enc 副檔名
// 3. 檢查檔名是否以 .md 結尾
//
// 支援的檔案格式：
// - 一般 Markdown 檔案：*.md
// - 加密 Markdown 檔案：*.md.enc
func (f *FileInfo) IsMarkdownFile() bool {
	// 目錄不是 Markdown 檔案
	if f.IsDirectory {
		return false
	}
	
	name := f.Name
	
	// 如果是加密檔案，移除 .enc 副檔名以檢查原始檔案類型
	if f.IsEncrypted && len(name) > 4 {
		name = name[:len(name)-4] // 移除 ".enc" 副檔名
	}
	
	// 檢查檔名是否以 .md 結尾
	return len(name) > 3 && name[len(name)-3:] == ".md"
}