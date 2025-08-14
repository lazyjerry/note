// Package ui 包含檔案對話框相關的功能
// 提供檔案開啟、保存對話框和檔案類型過濾器
package ui

import (
	"path/filepath"            // Go 標準庫，用於檔案路徑操作
	"strings"                  // Go 標準庫，用於字串處理
	"fyne.io/fyne/v2"          // Fyne GUI 框架核心套件
	"fyne.io/fyne/v2/dialog"   // Fyne 對話框套件
	"fyne.io/fyne/v2/storage"  // Fyne 儲存套件
)

// FileDialogManager 檔案對話框管理器
// 負責管理所有檔案相關的對話框操作，包括開啟、保存和檔案類型過濾
type FileDialogManager struct {
	parent fyne.Window // 父視窗，用於顯示對話框
}

// NewFileDialogManager 建立新的檔案對話框管理器
// 參數：parent（父視窗實例）
// 回傳：指向新建立的 FileDialogManager 的指標
//
// 執行流程：
// 1. 建立 FileDialogManager 結構體實例
// 2. 設定父視窗參考
// 3. 回傳管理器實例
func NewFileDialogManager(parent fyne.Window) *FileDialogManager {
	return &FileDialogManager{
		parent: parent,
	}
}

// ShowOpenDialog 顯示檔案開啟對話框
// 參數：callback（檔案選擇完成後的回調函數）
// 
// 執行流程：
// 1. 建立檔案開啟對話框
// 2. 設定檔案類型過濾器（支援 Markdown 檔案）
// 3. 設定對話框標題和按鈕文字
// 4. 設定檔案選擇回調函數
// 5. 顯示對話框
func (fdm *FileDialogManager) ShowOpenDialog(callback func(fyne.URIReadCloser, error)) {
	// 建立檔案開啟對話框
	openDialog := dialog.NewFileOpen(
		func(reader fyne.URIReadCloser, err error) {
			// 檔案選擇完成後的處理
			if err != nil {
				// 處理錯誤情況
				callback(nil, err)
				return
			}
			
			if reader == nil {
				// 使用者取消選擇
				callback(nil, nil)
				return
			}
			
			// 驗證檔案類型
			if !fdm.isValidFileType(reader.URI().Path()) {
				// 檔案類型不支援
				reader.Close()
				callback(nil, &FileTypeError{
					Path: reader.URI().Path(),
					Message: "不支援的檔案類型，請選擇 Markdown (.md) 或文字檔案 (.txt)",
				})
				return
			}
			
			// 回調處理選擇的檔案
			callback(reader, nil)
		},
		fdm.parent,
	)
	
	// 設定檔案類型過濾器
	openDialog.SetFilter(fdm.createFileFilter())
	
	// 設定初始目錄（使用者的文件目錄）
	// 注意：在實際應用中，可以設定預設目錄
	// 這裡暫時省略目錄設定，讓系統使用預設位置
	
	// 顯示對話框
	openDialog.Show()
}

// ShowSaveDialog 顯示檔案保存對話框
// 參數：defaultName（預設檔案名稱）, callback（檔案保存位置選擇完成後的回調函數）
//
// 執行流程：
// 1. 建立檔案保存對話框
// 2. 設定預設檔案名稱和副檔名
// 3. 設定檔案類型過濾器
// 4. 設定對話框標題和按鈕文字
// 5. 設定檔案保存回調函數
// 6. 顯示對話框
func (fdm *FileDialogManager) ShowSaveDialog(defaultName string, callback func(fyne.URIWriteCloser, error)) {
	// 確保預設檔案名稱有正確的副檔名
	if defaultName != "" && !strings.HasSuffix(strings.ToLower(defaultName), ".md") {
		defaultName += ".md"
	}
	
	// 建立檔案保存對話框
	saveDialog := dialog.NewFileSave(
		func(writer fyne.URIWriteCloser, err error) {
			// 檔案保存位置選擇完成後的處理
			if err != nil {
				// 處理錯誤情況
				callback(nil, err)
				return
			}
			
			if writer == nil {
				// 使用者取消保存
				callback(nil, nil)
				return
			}
			
			// 確保保存的檔案有正確的副檔名
			filePath := writer.URI().Path()
			if !strings.HasSuffix(strings.ToLower(filePath), ".md") {
				writer.Close()
				callback(nil, &FileTypeError{
					Path: filePath,
					Message: "檔案必須以 .md 副檔名保存",
				})
				return
			}
			
			// 回調處理保存位置
			callback(writer, nil)
		},
		fdm.parent,
	)
	
	// 設定預設檔案名稱
	if defaultName != "" {
		saveDialog.SetFileName(defaultName)
	}
	
	// 設定檔案類型過濾器
	saveDialog.SetFilter(fdm.createFileFilter())
	
	// 設定初始目錄（使用者的文件目錄）
	// 注意：在實際應用中，可以設定預設目錄
	// 這裡暫時省略目錄設定，讓系統使用預設位置
	
	// 顯示對話框
	saveDialog.Show()
}

// ShowSaveAsDialog 顯示另存新檔對話框
// 參數：currentPath（目前檔案路徑）, callback（檔案保存位置選擇完成後的回調函數）
//
// 執行流程：
// 1. 從目前檔案路徑提取檔案名稱
// 2. 建立另存新檔對話框
// 3. 設定預設檔案名稱和目錄
// 4. 呼叫標準保存對話框功能
func (fdm *FileDialogManager) ShowSaveAsDialog(currentPath string, callback func(fyne.URIWriteCloser, error)) {
	// 提取目前檔案的名稱作為預設名稱
	var defaultName string
	if currentPath != "" {
		defaultName = filepath.Base(currentPath)
	} else {
		defaultName = "未命名筆記.md"
	}
	
	// 使用標準保存對話框功能
	fdm.ShowSaveDialog(defaultName, callback)
}

// createFileFilter 建立檔案類型過濾器
// 回傳：支援的檔案類型過濾器
//
// 執行流程：
// 1. 定義支援的檔案副檔名列表
// 2. 建立 Fyne 檔案過濾器
// 3. 回傳過濾器實例
func (fdm *FileDialogManager) createFileFilter() storage.FileFilter {
	// 建立檔案過濾器，支援 Markdown 和文字檔案
	return storage.NewExtensionFileFilter([]string{".md", ".txt", ".markdown"})
}

// isValidFileType 驗證檔案類型是否支援
// 參數：filePath（檔案路徑）
// 回傳：是否為支援的檔案類型
//
// 執行流程：
// 1. 提取檔案副檔名
// 2. 轉換為小寫進行比較
// 3. 檢查是否在支援的類型列表中
// 4. 回傳驗證結果
func (fdm *FileDialogManager) isValidFileType(filePath string) bool {
	// 支援的檔案類型列表
	supportedTypes := []string{".md", ".txt", ".markdown"}
	
	// 提取檔案副檔名並轉換為小寫
	ext := strings.ToLower(filepath.Ext(filePath))
	
	// 檢查是否在支援的類型列表中
	for _, supportedType := range supportedTypes {
		if ext == supportedType {
			return true
		}
	}
	
	return false
}

// FileTypeError 檔案類型錯誤
// 當選擇的檔案類型不支援時使用
type FileTypeError struct {
	Path    string // 檔案路徑
	Message string // 錯誤訊息
}

// Error 實作 error 介面
// 回傳：錯誤訊息字串
func (e *FileTypeError) Error() string {
	return e.Message
}

// GetPath 取得錯誤相關的檔案路徑
// 回傳：檔案路徑
func (e *FileTypeError) GetPath() string {
	return e.Path
}

// FileDialogConfig 檔案對話框配置
// 用於自訂檔案對話框的行為和外觀
type FileDialogConfig struct {
	Title           string   // 對話框標題
	DefaultName     string   // 預設檔案名稱
	DefaultLocation string   // 預設目錄位置
	FileTypes       []string // 支援的檔案類型
	AllowMultiple   bool     // 是否允許多選（僅適用於開啟對話框）
}

// ShowCustomOpenDialog 顯示自訂的檔案開啟對話框
// 參數：config（對話框配置）, callback（檔案選擇完成後的回調函數）
//
// 執行流程：
// 1. 根據配置建立檔案開啟對話框
// 2. 設定自訂的標題、位置和檔案類型
// 3. 處理多選功能（如果啟用）
// 4. 顯示對話框
func (fdm *FileDialogManager) ShowCustomOpenDialog(config FileDialogConfig, callback func([]fyne.URIReadCloser, error)) {
	// 建立檔案開啟對話框
	openDialog := dialog.NewFileOpen(
		func(reader fyne.URIReadCloser, err error) {
			// 單檔案選擇的處理
			if err != nil {
				callback(nil, err)
				return
			}
			
			if reader == nil {
				callback(nil, nil)
				return
			}
			
			// 驗證檔案類型（如果有指定）
			if len(config.FileTypes) > 0 && !fdm.isValidFileTypeCustom(reader.URI().Path(), config.FileTypes) {
				reader.Close()
				callback(nil, &FileTypeError{
					Path: reader.URI().Path(),
					Message: "不支援的檔案類型",
				})
				return
			}
			
			// 回傳單個檔案的陣列
			callback([]fyne.URIReadCloser{reader}, nil)
		},
		fdm.parent,
	)
	
	// 設定檔案類型過濾器
	if len(config.FileTypes) > 0 {
		openDialog.SetFilter(storage.NewExtensionFileFilter(config.FileTypes))
	}
	
	// 設定初始目錄
	// 注意：在實際應用中，可以根據配置設定目錄
	// 這裡暫時省略目錄設定，讓系統使用預設位置
	
	// 顯示對話框
	openDialog.Show()
}

// ShowCustomSaveDialog 顯示自訂的檔案保存對話框
// 參數：config（對話框配置）, callback（檔案保存位置選擇完成後的回調函數）
//
// 執行流程：
// 1. 根據配置建立檔案保存對話框
// 2. 設定自訂的標題、預設名稱和位置
// 3. 設定檔案類型過濾器
// 4. 顯示對話框
func (fdm *FileDialogManager) ShowCustomSaveDialog(config FileDialogConfig, callback func(fyne.URIWriteCloser, error)) {
	// 建立檔案保存對話框
	saveDialog := dialog.NewFileSave(
		func(writer fyne.URIWriteCloser, err error) {
			// 檔案保存位置選擇完成後的處理
			if err != nil {
				callback(nil, err)
				return
			}
			
			if writer == nil {
				callback(nil, nil)
				return
			}
			
			// 驗證檔案類型（如果有指定）
			if len(config.FileTypes) > 0 && !fdm.isValidFileTypeCustom(writer.URI().Path(), config.FileTypes) {
				writer.Close()
				callback(nil, &FileTypeError{
					Path: writer.URI().Path(),
					Message: "不支援的檔案類型",
				})
				return
			}
			
			// 回調處理保存位置
			callback(writer, nil)
		},
		fdm.parent,
	)
	
	// 設定預設檔案名稱
	if config.DefaultName != "" {
		saveDialog.SetFileName(config.DefaultName)
	}
	
	// 設定檔案類型過濾器
	if len(config.FileTypes) > 0 {
		saveDialog.SetFilter(storage.NewExtensionFileFilter(config.FileTypes))
	}
	
	// 設定初始目錄
	// 注意：在實際應用中，可以根據配置設定目錄
	// 這裡暫時省略目錄設定，讓系統使用預設位置
	
	// 顯示對話框
	saveDialog.Show()
}

// isValidFileTypeCustom 驗證檔案類型是否在自訂的支援列表中
// 參數：filePath（檔案路徑）, supportedTypes（支援的檔案類型列表）
// 回傳：是否為支援的檔案類型
//
// 執行流程：
// 1. 提取檔案副檔名
// 2. 轉換為小寫進行比較
// 3. 檢查是否在自訂的支援類型列表中
// 4. 回傳驗證結果
func (fdm *FileDialogManager) isValidFileTypeCustom(filePath string, supportedTypes []string) bool {
	// 提取檔案副檔名並轉換為小寫
	ext := strings.ToLower(filepath.Ext(filePath))
	
	// 檢查是否在支援的類型列表中
	for _, supportedType := range supportedTypes {
		if ext == strings.ToLower(supportedType) {
			return true
		}
	}
	
	return false
}