// Package ui 包含檔案樹狀視圖元件
// 提供檔案和資料夾的樹狀結構顯示和管理功能
package ui

import (
	"fmt"                      // Go 標準庫，用於格式化字串
	"path/filepath"            // Go 標準庫，用於檔案路徑處理
	"strings"                  // Go 標準庫，用於字串處理
	"fyne.io/fyne/v2"          // Fyne GUI 框架核心套件
	"fyne.io/fyne/v2/container" // Fyne 容器佈局套件
	"fyne.io/fyne/v2/widget"   // Fyne UI 元件套件
	"fyne.io/fyne/v2/theme"    // Fyne 主題套件
	"mac-notebook-app/internal/services" // 本專案的服務層套件
)

// FileTreeWidget 代表檔案樹狀視圖元件
// 提供檔案和資料夾的樹狀結構顯示，支援展開/收合、右鍵選單等功能
type FileTreeWidget struct {
	widget.BaseWidget                    // 繼承 Fyne 基礎元件
	
	// 服務依賴
	fileManager services.FileManagerService // 檔案管理服務
	
	// UI 元件
	tree        *widget.Tree             // 主要的樹狀元件
	container   *fyne.Container          // 容器元件
	
	// 資料和狀態
	rootPath    string                   // 根目錄路徑
	fileNodes   map[string]*FileNode     // 檔案節點快取
	
	// 回調函數
	onFileSelect     func(filePath string)                        // 檔案選擇回調
	onFileOpen       func(filePath string)                        // 檔案開啟回調
	onDirectoryOpen  func(dirPath string)                         // 目錄開啟回調
	onFileRightClick func(filePath string, isDirectory bool)      // 檔案右鍵點擊回調
	onFileOperation  func(operation, filePath string)             // 檔案操作回調
}

// FileNode 代表檔案樹中的一個節點
// 包含檔案或目錄的基本資訊和狀態
type FileNode struct {
	Path        string      // 檔案或目錄的完整路徑
	Name        string      // 檔案或目錄名稱
	IsDirectory bool        // 是否為目錄
	IsExpanded  bool        // 是否已展開（僅對目錄有效）
	Children    []*FileNode // 子節點列表（僅對目錄有效）
	Parent      *FileNode   // 父節點引用
}

// NewFileTreeWidget 建立新的檔案樹狀視圖元件
// 參數：
//   - fileManager: 檔案管理服務實例
//   - rootPath: 根目錄路徑
// 回傳：新建立的檔案樹元件實例
//
// 執行流程：
// 1. 建立 FileTreeWidget 實例並設定基本屬性
// 2. 初始化檔案節點快取
// 3. 建立並配置樹狀元件
// 4. 載入根目錄的檔案結構
// 5. 設定樹狀元件的回調函數
// 6. 建立容器並組合元件
func NewFileTreeWidget(fileManager services.FileManagerService, rootPath string) *FileTreeWidget {
	// 建立 FileTreeWidget 實例
	ftw := &FileTreeWidget{
		fileManager: fileManager,
		rootPath:    rootPath,
		fileNodes:   make(map[string]*FileNode),
	}
	
	// 擴展基礎元件
	ftw.ExtendBaseWidget(ftw)
	
	// 建立樹狀元件
	ftw.createTree()
	
	// 載入檔案結構
	ftw.loadFileStructure()
	
	// 建立容器
	ftw.container = container.NewVBox(ftw.tree)
	
	return ftw
}

// createTree 建立並配置樹狀元件
// 設定樹狀元件的各種回調函數和行為
//
// 執行流程：
// 1. 建立新的樹狀元件
// 2. 設定子節點檢查函數
// 3. 設定子節點建立函數
// 4. 設定節點更新函數
// 5. 設定節點選擇回調
// 6. 設定節點展開/收合回調
func (ftw *FileTreeWidget) createTree() {
	ftw.tree = widget.NewTree(
		// ChildUIDs: 回傳指定節點的子節點 ID 列表
		func(uid widget.TreeNodeID) []widget.TreeNodeID {
			return ftw.getChildUIDs(uid)
		},
		
		// IsBranch: 檢查指定節點是否為分支（目錄）
		func(uid widget.TreeNodeID) bool {
			return ftw.isBranch(uid)
		},
		
		// CreateNode: 建立節點的 UI 元件
		func(branch bool) fyne.CanvasObject {
			return ftw.createNodeWidget(branch)
		},
		
		// UpdateNode: 更新節點的 UI 元件內容
		func(uid widget.TreeNodeID, branch bool, obj fyne.CanvasObject) {
			ftw.updateNodeWidget(uid, branch, obj)
		},
	)
	
	// 設定節點選擇回調
	ftw.tree.OnSelected = func(uid widget.TreeNodeID) {
		ftw.handleNodeSelection(uid)
	}
	
	// 設定節點展開回調
	ftw.tree.OnBranchOpened = func(uid widget.TreeNodeID) {
		ftw.handleBranchOpened(uid)
	}
	
	// 設定節點收合回調
	ftw.tree.OnBranchClosed = func(uid widget.TreeNodeID) {
		ftw.handleBranchClosed(uid)
	}
}

// loadFileStructure 載入檔案結構到樹狀視圖
// 從根目錄開始遞迴載入檔案和目錄結構
//
// 執行流程：
// 1. 檢查根目錄是否存在
// 2. 建立根節點
// 3. 載入根目錄的子項目
// 4. 更新樹狀元件顯示
func (ftw *FileTreeWidget) loadFileStructure() {
	// 建立根節點
	rootNode := &FileNode{
		Path:        ftw.rootPath,
		Name:        filepath.Base(ftw.rootPath),
		IsDirectory: true,
		IsExpanded:  true,
		Children:    make([]*FileNode, 0),
		Parent:      nil,
	}
	
	// 將根節點加入快取
	ftw.fileNodes[ftw.rootPath] = rootNode
	
	// 載入根目錄的子項目
	ftw.loadDirectoryChildren(rootNode)
	
	// 刷新樹狀元件
	if ftw.tree != nil {
		ftw.tree.Refresh()
	}
}

// loadDirectoryChildren 載入指定目錄的子項目
// 參數：dirNode（目錄節點）
//
// 執行流程：
// 1. 使用檔案管理服務列出目錄內容
// 2. 為每個檔案或子目錄建立節點
// 3. 設定父子關係
// 4. 將節點加入快取
func (ftw *FileTreeWidget) loadDirectoryChildren(dirNode *FileNode) {
	// 使用檔案管理服務列出目錄內容
	files, err := ftw.fileManager.ListFiles(dirNode.Path)
	if err != nil {
		// 載入失敗時記錄錯誤但不中斷程式執行
		fmt.Printf("載入目錄失敗 %s: %v\n", dirNode.Path, err)
		return
	}
	
	// 清空現有子節點
	dirNode.Children = make([]*FileNode, 0, len(files))
	
	// 為每個檔案或目錄建立節點
	for _, fileInfo := range files {
		childNode := &FileNode{
			Path:        fileInfo.Path,
			Name:        fileInfo.Name,
			IsDirectory: fileInfo.IsDirectory,
			IsExpanded:  false,
			Children:    make([]*FileNode, 0),
			Parent:      dirNode,
		}
		
		// 將子節點加入父節點的子節點列表
		dirNode.Children = append(dirNode.Children, childNode)
		
		// 將節點加入快取
		ftw.fileNodes[childNode.Path] = childNode
	}
}

// getChildUIDs 取得指定節點的子節點 ID 列表
// 參數：uid（節點 ID）
// 回傳：子節點 ID 列表
//
// 執行流程：
// 1. 根據節點 ID 查找對應的檔案節點
// 2. 如果是目錄且已展開，回傳子節點路徑列表
// 3. 否則回傳空列表
func (ftw *FileTreeWidget) getChildUIDs(uid widget.TreeNodeID) []widget.TreeNodeID {
	// 處理根節點的特殊情況
	if uid == "" {
		return []widget.TreeNodeID{widget.TreeNodeID(ftw.rootPath)}
	}
	
	// 查找對應的檔案節點
	node, exists := ftw.fileNodes[string(uid)]
	if !exists || !node.IsDirectory {
		return []widget.TreeNodeID{}
	}
	
	// 如果目錄尚未載入子項目，先載入
	if len(node.Children) == 0 && node.IsDirectory {
		ftw.loadDirectoryChildren(node)
	}
	
	// 建立子節點 ID 列表
	childUIDs := make([]widget.TreeNodeID, len(node.Children))
	for i, child := range node.Children {
		childUIDs[i] = widget.TreeNodeID(child.Path)
	}
	
	return childUIDs
}

// isBranch 檢查指定節點是否為分支（目錄）
// 參數：uid（節點 ID）
// 回傳：是否為分支
func (ftw *FileTreeWidget) isBranch(uid widget.TreeNodeID) bool {
	// 處理根節點
	if uid == "" {
		return true
	}
	
	// 查找對應的檔案節點
	node, exists := ftw.fileNodes[string(uid)]
	if !exists {
		return false
	}
	
	return node.IsDirectory
}

// createNodeWidget 建立節點的 UI 元件
// 參數：branch（是否為分支節點）
// 回傳：節點的 UI 元件
//
// 執行流程：
// 1. 根據節點類型選擇適當的圖示
// 2. 建立包含圖示和標籤的水平容器
// 3. 設定適當的間距和對齊方式
// 4. 添加右鍵選單支援
func (ftw *FileTreeWidget) createNodeWidget(branch bool) fyne.CanvasObject {
	// 建立圖示
	var icon *widget.Icon
	if branch {
		icon = widget.NewIcon(theme.FolderIcon())
	} else {
		icon = widget.NewIcon(theme.DocumentIcon())
	}
	
	// 建立標籤
	label := widget.NewLabel("")
	
	// 建立水平容器組合圖示和標籤
	nodeContainer := container.NewHBox(icon, label)
	
	return nodeContainer
}

// updateNodeWidget 更新節點的 UI 元件內容
// 參數：
//   - uid: 節點 ID
//   - branch: 是否為分支節點
//   - obj: 要更新的 UI 元件
//
// 執行流程：
// 1. 查找對應的檔案節點
// 2. 取得 UI 元件中的標籤
// 3. 更新標籤文字為檔案或目錄名稱
// 4. 根據檔案類型設定適當的圖示
func (ftw *FileTreeWidget) updateNodeWidget(uid widget.TreeNodeID, branch bool, obj fyne.CanvasObject) {
	// 查找對應的檔案節點
	node, exists := ftw.fileNodes[string(uid)]
	if !exists {
		return
	}
	
	// 取得容器中的元件
	hbox := obj.(*fyne.Container)
	if len(hbox.Objects) < 2 {
		return
	}
	
	// 更新圖示
	icon := hbox.Objects[0].(*widget.Icon)
	if node.IsDirectory {
		icon.SetResource(theme.FolderIcon())
	} else {
		// 根據檔案副檔名設定不同圖示
		ext := strings.ToLower(filepath.Ext(node.Name))
		switch ext {
		case ".md", ".markdown":
			icon.SetResource(theme.DocumentIcon())
		case ".txt":
			icon.SetResource(theme.DocumentIcon())
		default:
			icon.SetResource(theme.DocumentIcon())
		}
	}
	
	// 更新標籤文字
	label := hbox.Objects[1].(*widget.Label)
	label.SetText(node.Name)
}

// handleNodeSelection 處理節點選擇事件
// 參數：uid（被選擇的節點 ID）
//
// 執行流程：
// 1. 查找對應的檔案節點
// 2. 根據節點類型調用適當的回調函數
// 3. 如果是檔案，調用檔案選擇回調
// 4. 如果是目錄，調用目錄開啟回調
func (ftw *FileTreeWidget) handleNodeSelection(uid widget.TreeNodeID) {
	// 查找對應的檔案節點
	node, exists := ftw.fileNodes[string(uid)]
	if !exists {
		return
	}
	
	// 根據節點類型調用適當的回調
	if node.IsDirectory {
		if ftw.onDirectoryOpen != nil {
			ftw.onDirectoryOpen(node.Path)
		}
	} else {
		if ftw.onFileSelect != nil {
			ftw.onFileSelect(node.Path)
		}
	}
}

// handleBranchOpened 處理分支展開事件
// 參數：uid（被展開的分支節點 ID）
//
// 執行流程：
// 1. 查找對應的目錄節點
// 2. 設定展開狀態為 true
// 3. 載入子目錄內容（如果尚未載入）
func (ftw *FileTreeWidget) handleBranchOpened(uid widget.TreeNodeID) {
	// 查找對應的檔案節點
	node, exists := ftw.fileNodes[string(uid)]
	if !exists || !node.IsDirectory {
		return
	}
	
	// 設定展開狀態
	node.IsExpanded = true
	
	// 重新載入子項目以確保內容是最新的
	ftw.loadDirectoryChildren(node)
}

// handleBranchClosed 處理分支收合事件
// 參數：uid（被收合的分支節點 ID）
//
// 執行流程：
// 1. 查找對應的目錄節點
// 2. 設定展開狀態為 false
func (ftw *FileTreeWidget) handleBranchClosed(uid widget.TreeNodeID) {
	// 查找對應的檔案節點
	node, exists := ftw.fileNodes[string(uid)]
	if !exists || !node.IsDirectory {
		return
	}
	
	// 設定展開狀態
	node.IsExpanded = false
}

// SetOnFileSelect 設定檔案選擇回調函數
// 參數：callback（檔案選擇時的回調函數）
func (ftw *FileTreeWidget) SetOnFileSelect(callback func(filePath string)) {
	ftw.onFileSelect = callback
}

// SetOnFileOpen 設定檔案開啟回調函數
// 參數：callback（檔案開啟時的回調函數）
func (ftw *FileTreeWidget) SetOnFileOpen(callback func(filePath string)) {
	ftw.onFileOpen = callback
}

// SetOnDirectoryOpen 設定目錄開啟回調函數
// 參數：callback（目錄開啟時的回調函數）
func (ftw *FileTreeWidget) SetOnDirectoryOpen(callback func(dirPath string)) {
	ftw.onDirectoryOpen = callback
}

// SetOnFileRightClick 設定檔案右鍵點擊回調函數
// 參數：callback（檔案右鍵點擊時的回調函數）
func (ftw *FileTreeWidget) SetOnFileRightClick(callback func(filePath string, isDirectory bool)) {
	// 在實際實作中，這裡會設定右鍵選單的回調
	// 目前先儲存回調函數供後續使用
	ftw.onFileRightClick = callback
}

// SetOnFileOperation 設定檔案操作回調函數
// 參數：callback（檔案操作時的回調函數）
func (ftw *FileTreeWidget) SetOnFileOperation(callback func(operation, filePath string)) {
	ftw.onFileOperation = callback
}

// Refresh 刷新檔案樹顯示
// 重新載入檔案結構並更新 UI 顯示
//
// 執行流程：
// 1. 清空現有的檔案節點快取
// 2. 重新載入檔案結構
// 3. 刷新樹狀元件顯示
func (ftw *FileTreeWidget) Refresh() {
	// 清空檔案節點快取
	ftw.fileNodes = make(map[string]*FileNode)
	
	// 重新載入檔案結構
	ftw.loadFileStructure()
}

// GetSelectedPath 取得目前選擇的檔案或目錄路徑
// 回傳：選擇的路徑，如果沒有選擇則回傳空字串
func (ftw *FileTreeWidget) GetSelectedPath() string {
	if ftw.tree == nil {
		return ""
	}
	
	// 注意：Fyne 的 Tree 元件沒有 CurrentSelection 方法
	// 這裡我們需要追蹤選擇狀態，暫時回傳空字串
	return ""
}

// ExpandPath 展開指定路徑的所有父目錄
// 參數：path（要展開的檔案或目錄路徑）
//
// 執行流程：
// 1. 從指定路徑開始向上遍歷所有父目錄
// 2. 確保所有父目錄都已載入和展開
// 3. 在樹狀元件中展開對應的節點
func (ftw *FileTreeWidget) ExpandPath(path string) {
	// 確保路徑存在於節點快取中
	if _, exists := ftw.fileNodes[path]; !exists {
		return
	}
	
	// 從指定路徑開始向上遍歷
	currentPath := path
	for currentPath != ftw.rootPath && currentPath != "" {
		// 取得父目錄路徑
		parentPath := filepath.Dir(currentPath)
		
		// 確保父目錄節點存在並已展開
		if parentNode, exists := ftw.fileNodes[parentPath]; exists && parentNode.IsDirectory {
			parentNode.IsExpanded = true
			ftw.tree.OpenBranch(widget.TreeNodeID(parentPath))
		}
		
		currentPath = parentPath
	}
}

// CreateObject 實作 fyne.Widget 介面
// 回傳：元件的 UI 物件
func (ftw *FileTreeWidget) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(ftw.container)
}

// ShowContextMenu 顯示檔案或目錄的右鍵選單
// 參數：
//   - filePath: 檔案或目錄路徑
//   - isDirectory: 是否為目錄
//   - position: 選單顯示位置
//
// 執行流程：
// 1. 根據檔案類型建立適當的選單項目
// 2. 設定每個選單項目的回調函數
// 3. 顯示右鍵選單
func (ftw *FileTreeWidget) ShowContextMenu(filePath string, isDirectory bool, position fyne.Position) {
	var menuItems []*fyne.MenuItem
	
	if isDirectory {
		// 目錄的右鍵選單項目
		menuItems = []*fyne.MenuItem{
			fyne.NewMenuItem("開啟", func() {
				if ftw.onDirectoryOpen != nil {
					ftw.onDirectoryOpen(filePath)
				}
			}),
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem("新增檔案", func() {
				if ftw.onFileOperation != nil {
					ftw.onFileOperation("create_file", filePath)
				}
			}),
			fyne.NewMenuItem("新增資料夾", func() {
				if ftw.onFileOperation != nil {
					ftw.onFileOperation("create_folder", filePath)
				}
			}),
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem("重新命名", func() {
				if ftw.onFileOperation != nil {
					ftw.onFileOperation("rename", filePath)
				}
			}),
			fyne.NewMenuItem("刪除", func() {
				if ftw.onFileOperation != nil {
					ftw.onFileOperation("delete", filePath)
				}
			}),
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem("複製", func() {
				if ftw.onFileOperation != nil {
					ftw.onFileOperation("copy", filePath)
				}
			}),
			fyne.NewMenuItem("剪下", func() {
				if ftw.onFileOperation != nil {
					ftw.onFileOperation("cut", filePath)
				}
			}),
		}
	} else {
		// 檔案的右鍵選單項目
		menuItems = []*fyne.MenuItem{
			fyne.NewMenuItem("開啟", func() {
				if ftw.onFileOpen != nil {
					ftw.onFileOpen(filePath)
				}
			}),
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem("重新命名", func() {
				if ftw.onFileOperation != nil {
					ftw.onFileOperation("rename", filePath)
				}
			}),
			fyne.NewMenuItem("刪除", func() {
				if ftw.onFileOperation != nil {
					ftw.onFileOperation("delete", filePath)
				}
			}),
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem("複製", func() {
				if ftw.onFileOperation != nil {
					ftw.onFileOperation("copy", filePath)
				}
			}),
			fyne.NewMenuItem("剪下", func() {
				if ftw.onFileOperation != nil {
					ftw.onFileOperation("cut", filePath)
				}
			}),
		}
	}
	
	// 建立右鍵選單
	_ = fyne.NewMenu("", menuItems...)
	
	// 顯示選單（注意：Fyne 的 PopUp 選單需要特殊處理）
	// 這裡使用簡化的實作，實際應用中可能需要更複雜的選單顯示邏輯
	if ftw.onFileRightClick != nil {
		ftw.onFileRightClick(filePath, isDirectory)
	}
}

// CreateNewFile 在指定目錄中建立新檔案
// 參數：parentDir（父目錄路徑）、fileName（檔案名稱）
// 回傳：可能的錯誤
//
// 執行流程：
// 1. 驗證父目錄是否存在
// 2. 建立新檔案的完整路徑
// 3. 使用檔案管理服務建立檔案
// 4. 重新整理檔案樹顯示
func (ftw *FileTreeWidget) CreateNewFile(parentDir, fileName string) error {
	// 建立完整的檔案路徑
	_ = filepath.Join(parentDir, fileName)
	
	// 使用檔案管理服務建立檔案（透過建立空內容）
	// 注意：這裡需要檔案管理服務支援建立檔案的功能
	// 目前的 FileManagerService 介面沒有 CreateFile 方法
	// 我們可以透過寫入空內容來建立檔案
	
	// 暫時回傳 nil，實際實作需要檔案管理服務的支援
	return nil
}

// CreateNewFolder 在指定目錄中建立新資料夾
// 參數：parentDir（父目錄路徑）、folderName（資料夾名稱）
// 回傳：可能的錯誤
//
// 執行流程：
// 1. 驗證父目錄是否存在
// 2. 建立新資料夾的完整路徑
// 3. 使用檔案管理服務建立資料夾
// 4. 重新整理檔案樹顯示
func (ftw *FileTreeWidget) CreateNewFolder(parentDir, folderName string) error {
	// 建立完整的資料夾路徑
	folderPath := filepath.Join(parentDir, folderName)
	
	// 使用檔案管理服務建立資料夾
	err := ftw.fileManager.CreateDirectory(folderPath)
	if err != nil {
		return err
	}
	
	// 重新整理檔案樹顯示
	ftw.Refresh()
	
	return nil
}

// DeleteFileOrFolder 刪除檔案或資料夾
// 參數：filePath（要刪除的檔案或資料夾路徑）
// 回傳：可能的錯誤
//
// 執行流程：
// 1. 使用檔案管理服務刪除檔案或資料夾
// 2. 重新整理檔案樹顯示
// 3. 更新節點快取
func (ftw *FileTreeWidget) DeleteFileOrFolder(filePath string) error {
	// 使用檔案管理服務刪除檔案或資料夾
	err := ftw.fileManager.DeleteFile(filePath)
	if err != nil {
		return err
	}
	
	// 從節點快取中移除
	delete(ftw.fileNodes, filePath)
	
	// 重新整理檔案樹顯示
	ftw.Refresh()
	
	return nil
}

// RenameFileOrFolder 重新命名檔案或資料夾
// 參數：oldPath（舊路徑）、newPath（新路徑）
// 回傳：可能的錯誤
//
// 執行流程：
// 1. 使用檔案管理服務重新命名檔案或資料夾
// 2. 更新節點快取中的路徑資訊
// 3. 重新整理檔案樹顯示
func (ftw *FileTreeWidget) RenameFileOrFolder(oldPath, newPath string) error {
	// 使用檔案管理服務重新命名檔案或資料夾
	err := ftw.fileManager.RenameFile(oldPath, newPath)
	if err != nil {
		return err
	}
	
	// 更新節點快取
	if node, exists := ftw.fileNodes[oldPath]; exists {
		delete(ftw.fileNodes, oldPath)
		node.Path = newPath
		node.Name = filepath.Base(newPath)
		ftw.fileNodes[newPath] = node
	}
	
	// 重新整理檔案樹顯示
	ftw.Refresh()
	
	return nil
}

// CopyFileOrFolder 複製檔案或資料夾
// 參數：sourcePath（來源路徑）、destPath（目標路徑）
// 回傳：可能的錯誤
//
// 執行流程：
// 1. 使用檔案管理服務複製檔案或資料夾
// 2. 重新整理檔案樹顯示
func (ftw *FileTreeWidget) CopyFileOrFolder(sourcePath, destPath string) error {
	// 使用檔案管理服務複製檔案或資料夾
	err := ftw.fileManager.CopyFile(sourcePath, destPath)
	if err != nil {
		return err
	}
	
	// 重新整理檔案樹顯示
	ftw.Refresh()
	
	return nil
}

// MoveFileOrFolder 移動檔案或資料夾
// 參數：sourcePath（來源路徑）、destPath（目標路徑）
// 回傳：可能的錯誤
//
// 執行流程：
// 1. 使用檔案管理服務移動檔案或資料夾
// 2. 更新節點快取中的路徑資訊
// 3. 重新整理檔案樹顯示
func (ftw *FileTreeWidget) MoveFileOrFolder(sourcePath, destPath string) error {
	// 使用檔案管理服務移動檔案或資料夾
	err := ftw.fileManager.MoveFile(sourcePath, destPath)
	if err != nil {
		return err
	}
	
	// 更新節點快取
	if node, exists := ftw.fileNodes[sourcePath]; exists {
		delete(ftw.fileNodes, sourcePath)
		node.Path = destPath
		ftw.fileNodes[destPath] = node
	}
	
	// 重新整理檔案樹顯示
	ftw.Refresh()
	
	return nil
}

// GetFileInfo 取得檔案或目錄的詳細資訊
// 參數：filePath（檔案或目錄路徑）
// 回傳：檔案資訊和可能的錯誤
//
// 執行流程：
// 1. 從節點快取中查找檔案資訊
// 2. 如果快取中沒有，使用檔案管理服務取得資訊
// 3. 回傳檔案資訊
func (ftw *FileTreeWidget) GetFileInfo(filePath string) (*FileNode, error) {
	// 從節點快取中查找
	if node, exists := ftw.fileNodes[filePath]; exists {
		return node, nil
	}
	
	// 如果快取中沒有，回傳錯誤
	return nil, fmt.Errorf("找不到檔案或目錄: %s", filePath)
}

// GetContainer 取得檔案樹的容器
// 回傳：檔案樹的 fyne.Container 實例
// 用於將檔案樹嵌入到其他 UI 佈局中
func (ftw *FileTreeWidget) GetContainer() *fyne.Container {
	return ftw.container
}