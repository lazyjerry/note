// Package ui 包含拖拽功能相關的元件和管理
// 提供檔案拖拽支援、視覺回饋和檔案管理服務整合
package ui

import (
	"fmt"                                    // Go 標準庫，用於格式化字串
	"path/filepath"                          // Go 標準庫，用於檔案路徑操作
	"strings"                                // Go 標準庫，用於字串處理
	"fyne.io/fyne/v2"                        // Fyne GUI 框架核心套件
	"fyne.io/fyne/v2/container"              // Fyne 容器佈局套件
	"fyne.io/fyne/v2/widget"                 // Fyne UI 元件套件
	"fyne.io/fyne/v2/theme"                  // Fyne 主題套件
	"mac-notebook-app/internal/services"     // 內部服務套件
)

// DragDropManager 拖拽管理器
// 負責管理檔案拖拽操作、視覺回饋和與檔案管理服務的整合
type DragDropManager struct {
	fileManager    services.FileManagerService // 檔案管理服務
	parent         fyne.Window                 // 父視窗
	dropZones      map[string]*DropZone        // 拖拽區域映射
	dragFeedback   *DragFeedback               // 拖拽視覺回饋
	onFileDropped  func(sourcePath, targetPath string) error // 檔案拖拽完成回調
	onFileMoved    func(oldPath, newPath string)             // 檔案移動完成回調
	onError        func(error)                               // 錯誤處理回調
}

// DropZone 拖拽區域
// 定義可以接受拖拽操作的 UI 區域
type DropZone struct {
	widget       fyne.CanvasObject // 關聯的 UI 元件
	targetPath   string            // 目標路徑
	acceptTypes  []string          // 接受的檔案類型
	isActive     bool              // 是否啟用拖拽
	onDragEnter  func()            // 拖拽進入回調
	onDragLeave  func()            // 拖拽離開回調
	onDrop       func(sourcePath string) error // 拖拽放下回調
}

// DragFeedback 拖拽視覺回饋
// 提供拖拽操作過程中的視覺指示和狀態顯示
type DragFeedback struct {
	overlay      *fyne.Container // 覆蓋層容器
	indicator    *widget.Label   // 拖拽指示器
	isVisible    bool            // 是否可見
	currentZone  string          // 目前拖拽區域
}

// DragOperation 拖拽操作類型
type DragOperation int

const (
	DragOperationMove DragOperation = iota // 移動操作
	DragOperationCopy                      // 複製操作
	DragOperationLink                      // 連結操作
)

// NewDragDropManager 建立新的拖拽管理器
// 參數：fileManager（檔案管理服務）, parent（父視窗）
// 回傳：指向新建立的 DragDropManager 的指標
//
// 執行流程：
// 1. 建立 DragDropManager 結構體實例
// 2. 設定檔案管理服務和父視窗
// 3. 初始化拖拽區域映射和視覺回饋
// 4. 回傳管理器實例
func NewDragDropManager(fileManager services.FileManagerService, parent fyne.Window) *DragDropManager {
	return &DragDropManager{
		fileManager:  fileManager,
		parent:       parent,
		dropZones:    make(map[string]*DropZone),
		dragFeedback: NewDragFeedback(),
	}
}

// RegisterDropZone 註冊拖拽區域
// 參數：zoneID（區域識別碼）, widget（UI 元件）, targetPath（目標路徑）, acceptTypes（接受的檔案類型）
//
// 執行流程：
// 1. 建立新的 DropZone 實例
// 2. 設定區域屬性和回調函數
// 3. 將區域註冊到管理器中
// 4. 啟用區域的拖拽功能
func (ddm *DragDropManager) RegisterDropZone(zoneID string, widget fyne.CanvasObject, targetPath string, acceptTypes []string) {
	// 建立拖拽區域
	zone := &DropZone{
		widget:      widget,
		targetPath:  targetPath,
		acceptTypes: acceptTypes,
		isActive:    true,
		onDragEnter: func() {
			ddm.handleDragEnter(zoneID)
		},
		onDragLeave: func() {
			ddm.handleDragLeave(zoneID)
		},
		onDrop: func(sourcePath string) error {
			return ddm.handleDrop(sourcePath, targetPath)
		},
	}
	
	// 註冊到管理器
	ddm.dropZones[zoneID] = zone
	
	// 啟用拖拽功能（注意：Fyne 目前對拖拽支援有限，這裡提供架構）
	ddm.enableDragDrop(widget, zone)
}

// UnregisterDropZone 取消註冊拖拽區域
// 參數：zoneID（區域識別碼）
//
// 執行流程：
// 1. 檢查區域是否存在
// 2. 停用區域的拖拽功能
// 3. 從管理器中移除區域
func (ddm *DragDropManager) UnregisterDropZone(zoneID string) {
	if zone, exists := ddm.dropZones[zoneID]; exists {
		// 停用拖拽功能
		ddm.disableDragDrop(zone.widget)
		
		// 從映射中移除
		delete(ddm.dropZones, zoneID)
	}
}

// SetCallbacks 設定回調函數
// 參數：onFileDropped（檔案拖拽完成回調）, onFileMoved（檔案移動完成回調）, onError（錯誤處理回調）
//
// 執行流程：
// 1. 設定檔案拖拽完成回調函數
// 2. 設定檔案移動完成回調函數
// 3. 設定錯誤處理回調函數
func (ddm *DragDropManager) SetCallbacks(
	onFileDropped func(sourcePath, targetPath string) error,
	onFileMoved func(oldPath, newPath string),
	onError func(error),
) {
	ddm.onFileDropped = onFileDropped
	ddm.onFileMoved = onFileMoved
	ddm.onError = onError
}

// EnableZone 啟用拖拽區域
// 參數：zoneID（區域識別碼）
//
// 執行流程：
// 1. 檢查區域是否存在
// 2. 設定區域為啟用狀態
// 3. 更新 UI 元件的拖拽功能
func (ddm *DragDropManager) EnableZone(zoneID string) {
	if zone, exists := ddm.dropZones[zoneID]; exists {
		zone.isActive = true
		ddm.enableDragDrop(zone.widget, zone)
	}
}

// DisableZone 停用拖拽區域
// 參數：zoneID（區域識別碼）
//
// 執行流程：
// 1. 檢查區域是否存在
// 2. 設定區域為停用狀態
// 3. 停用 UI 元件的拖拽功能
func (ddm *DragDropManager) DisableZone(zoneID string) {
	if zone, exists := ddm.dropZones[zoneID]; exists {
		zone.isActive = false
		ddm.disableDragDrop(zone.widget)
	}
}

// IsValidFileType 檢查檔案類型是否被接受
// 參數：filePath（檔案路徑）, acceptTypes（接受的檔案類型列表）
// 回傳：是否為接受的檔案類型
//
// 執行流程：
// 1. 提取檔案副檔名
// 2. 轉換為小寫進行比較
// 3. 檢查是否在接受的類型列表中
// 4. 回傳檢查結果
func (ddm *DragDropManager) IsValidFileType(filePath string, acceptTypes []string) bool {
	// 如果沒有指定接受類型，則接受所有檔案
	if len(acceptTypes) == 0 {
		return true
	}
	
	// 提取檔案副檔名
	ext := strings.ToLower(filepath.Ext(filePath))
	
	// 檢查是否在接受的類型列表中
	for _, acceptType := range acceptTypes {
		if strings.ToLower(acceptType) == ext {
			return true
		}
	}
	
	return false
}

// handleDragEnter 處理拖拽進入事件
// 參數：zoneID（區域識別碼）
//
// 執行流程：
// 1. 更新視覺回饋顯示拖拽進入狀態
// 2. 高亮顯示目標區域
// 3. 顯示拖拽指示器
func (ddm *DragDropManager) handleDragEnter(zoneID string) {
	// 更新視覺回饋
	ddm.dragFeedback.ShowEnterFeedback(zoneID)
	
	// 高亮顯示目標區域
	if zone, exists := ddm.dropZones[zoneID]; exists {
		ddm.highlightDropZone(zone, true)
	}
}

// handleDragLeave 處理拖拽離開事件
// 參數：zoneID（區域識別碼）
//
// 執行流程：
// 1. 更新視覺回饋隱藏拖拽狀態
// 2. 取消高亮顯示目標區域
// 3. 隱藏拖拽指示器
func (ddm *DragDropManager) handleDragLeave(zoneID string) {
	// 更新視覺回饋
	ddm.dragFeedback.ShowLeaveFeedback()
	
	// 取消高亮顯示目標區域
	if zone, exists := ddm.dropZones[zoneID]; exists {
		ddm.highlightDropZone(zone, false)
	}
}

// handleDrop 處理拖拽放下事件
// 參數：sourcePath（來源路徑）, targetPath（目標路徑）
// 回傳：處理結果錯誤
//
// 執行流程：
// 1. 驗證來源和目標路徑
// 2. 檢查檔案類型是否被接受
// 3. 執行檔案移動或複製操作
// 4. 更新視覺回饋和呼叫回調函數
// 5. 處理錯誤情況
func (ddm *DragDropManager) handleDrop(sourcePath, targetPath string) error {
	// 隱藏視覺回饋
	ddm.dragFeedback.Hide()
	
	// 驗證路徑
	if sourcePath == "" || targetPath == "" {
		err := fmt.Errorf("無效的拖拽路徑：來源 '%s'，目標 '%s'", sourcePath, targetPath)
		if ddm.onError != nil {
			ddm.onError(err)
		}
		return err
	}
	
	// 檢查來源檔案是否存在
	if !ddm.fileExists(sourcePath) {
		err := fmt.Errorf("來源檔案不存在：%s", sourcePath)
		if ddm.onError != nil {
			ddm.onError(err)
		}
		return err
	}
	
	// 執行檔案操作
	var err error
	var newPath string
	
	// 判斷目標是檔案還是目錄
	if ddm.isDirectory(targetPath) {
		// 目標是目錄，移動檔案到目錄中
		fileName := filepath.Base(sourcePath)
		newPath = filepath.Join(targetPath, fileName)
		err = ddm.fileManager.MoveFile(sourcePath, newPath)
	} else {
		// 目標是檔案，直接移動
		newPath = targetPath
		err = ddm.fileManager.MoveFile(sourcePath, newPath)
	}
	
	// 處理結果
	if err != nil {
		if ddm.onError != nil {
			ddm.onError(err)
		}
		return err
	}
	
	// 呼叫成功回調
	if ddm.onFileDropped != nil {
		if dropErr := ddm.onFileDropped(sourcePath, newPath); dropErr != nil {
			if ddm.onError != nil {
				ddm.onError(dropErr)
			}
			return dropErr
		}
	}
	
	if ddm.onFileMoved != nil {
		ddm.onFileMoved(sourcePath, newPath)
	}
	
	return nil
}

// enableDragDrop 啟用 UI 元件的拖拽功能
// 參數：widget（UI 元件）, zone（拖拽區域）
//
// 執行流程：
// 1. 檢查元件類型並設定拖拽屬性
// 2. 註冊拖拽事件處理器
// 3. 設定視覺回饋
//
// 注意：Fyne 目前對拖拽支援有限，這裡提供基礎架構
func (ddm *DragDropManager) enableDragDrop(widget fyne.CanvasObject, zone *DropZone) {
	// 注意：Fyne v2 對原生拖拽支援有限
	// 這裡提供架構，實際實作可能需要使用 Fyne 的擴展或自訂實作
	
	// 為不同類型的元件設定拖拽功能
	switch w := widget.(type) {
	// Note: widget.Tree is not available in current Fyne version
	// case *widget.Tree:
	//	ddm.enableTreeDragDrop(w, zone)
	case *fyne.Container:
		// 為容器設定拖拽功能
		ddm.enableContainerDragDrop(w, zone)
	default:
		// 為一般元件設定基本拖拽功能
		ddm.enableBasicDragDrop(widget, zone)
	}
}

// disableDragDrop 停用 UI 元件的拖拽功能
// 參數：widget（UI 元件）
//
// 執行流程：
// 1. 移除拖拽事件處理器
// 2. 重置元件的拖拽屬性
// 3. 清理視覺回饋
func (ddm *DragDropManager) disableDragDrop(widget fyne.CanvasObject) {
	// 根據元件類型停用拖拽功能
	// 實際實作取決於 Fyne 的拖拽 API
}

// enableTreeDragDrop 為樹狀元件啟用拖拽功能
// 參數：tree（樹狀元件）, zone（拖拽區域）
//
// 執行流程：
// 1. 設定樹狀節點的拖拽處理
// 2. 註冊節點拖拽事件
// 3. 實作節點間的拖拽移動
//
// Note: This is a placeholder for future Tree widget support
func (ddm *DragDropManager) enableTreeDragDrop(tree interface{}, zone *DropZone) {
	// 為樹狀元件實作拖拽功能
	// 這裡可以擴展樹狀元件的拖拽行為
	// 當 Fyne 支援 Tree widget 時，可以實作具體功能
}



// enableContainerDragDrop 為容器啟用拖拽功能
// 參數：container（容器）, zone（拖拽區域）
//
// 執行流程：
// 1. 設定容器的拖拽接受區域
// 2. 註冊拖拽進入和離開事件
// 3. 處理拖拽放下操作
func (ddm *DragDropManager) enableContainerDragDrop(container *fyne.Container, zone *DropZone) {
	// 為容器實作拖拽功能
	// 這裡可以設定容器作為拖拽目標
}

// enableBasicDragDrop 為一般元件啟用基本拖拽功能
// 參數：widget（UI 元件）, zone（拖拽區域）
//
// 執行流程：
// 1. 設定元件的基本拖拽屬性
// 2. 註冊基本拖拽事件處理
// 3. 提供預設的拖拽行為
func (ddm *DragDropManager) enableBasicDragDrop(widget fyne.CanvasObject, zone *DropZone) {
	// 為一般元件實作基本拖拽功能
	// 這裡提供預設的拖拽行為
}

// startDragOperation 開始拖拽操作
// 參數：sourcePath（來源路徑）, zone（拖拽區域）
//
// 執行流程：
// 1. 初始化拖拽操作狀態
// 2. 顯示拖拽視覺回饋
// 3. 準備拖拽資料
func (ddm *DragDropManager) startDragOperation(sourcePath string, zone *DropZone) {
	// 顯示拖拽開始的視覺回饋
	ddm.dragFeedback.ShowDragStart(sourcePath)
	
	// 這裡可以設定拖拽游標和視覺效果
}

// highlightDropZone 高亮顯示拖拽區域
// 參數：zone（拖拽區域）, highlight（是否高亮）
//
// 執行流程：
// 1. 根據高亮狀態設定區域樣式
// 2. 更新區域的視覺外觀
// 3. 提供視覺回饋給使用者
func (ddm *DragDropManager) highlightDropZone(zone *DropZone, highlight bool) {
	// 根據元件類型設定高亮效果
	switch zone.widget.(type) {
	case *widget.Card:
		// 為卡片元件設定高亮邊框
		if highlight {
			// 設定高亮樣式（概念性實作）
		} else {
			// 恢復正常樣式
		}
	case *fyne.Container:
		// 為容器設定高亮背景
		if highlight {
			// 設定高亮背景色
		} else {
			// 恢復正常背景色
		}
	}
}

// fileExists 檢查檔案是否存在
// 參數：path（檔案路徑）
// 回傳：檔案是否存在
func (ddm *DragDropManager) fileExists(path string) bool {
	if ddm.fileManager == nil {
		return false
	}
	
	// 使用檔案管理服務檢查檔案存在性
	files, err := ddm.fileManager.ListFiles(filepath.Dir(path))
	if err != nil {
		return false
	}
	
	fileName := filepath.Base(path)
	for _, file := range files {
		if file.Name == fileName {
			return true
		}
	}
	
	return false
}

// isDirectory 檢查路徑是否為目錄
// 參數：path（檔案路徑）
// 回傳：是否為目錄
func (ddm *DragDropManager) isDirectory(path string) bool {
	if ddm.fileManager == nil {
		return false
	}
	
	// 使用檔案管理服務檢查是否為目錄
	files, err := ddm.fileManager.ListFiles(filepath.Dir(path))
	if err != nil {
		return false
	}
	
	fileName := filepath.Base(path)
	for _, file := range files {
		if file.Name == fileName {
			return file.IsDirectory
		}
	}
	
	return false
}

// NewDragFeedback 建立新的拖拽視覺回饋
// 回傳：指向新建立的 DragFeedback 的指標
//
// 執行流程：
// 1. 建立拖拽指示器標籤
// 2. 建立覆蓋層容器
// 3. 設定初始狀態為隱藏
// 4. 回傳視覺回饋實例
func NewDragFeedback() *DragFeedback {
	// 建立拖拽指示器
	indicator := widget.NewLabel("")
	indicator.TextStyle = fyne.TextStyle{Bold: true}
	indicator.Alignment = fyne.TextAlignCenter
	
	// 建立覆蓋層容器
	overlay := container.NewWithoutLayout(indicator)
	
	return &DragFeedback{
		overlay:   overlay,
		indicator: indicator,
		isVisible: false,
	}
}

// ShowDragStart 顯示拖拽開始的視覺回饋
// 參數：sourcePath（來源路徑）
//
// 執行流程：
// 1. 設定拖拽指示器文字
// 2. 顯示覆蓋層
// 3. 更新可見狀態
func (df *DragFeedback) ShowDragStart(sourcePath string) {
	fileName := filepath.Base(sourcePath)
	df.indicator.SetText(fmt.Sprintf("拖拽中: %s", fileName))
	df.indicator.Importance = widget.MediumImportance
	
	df.isVisible = true
	df.overlay.Show()
}

// ShowEnterFeedback 顯示拖拽進入區域的視覺回饋
// 參數：zoneID（區域識別碼）
//
// 執行流程：
// 1. 更新指示器文字顯示目標區域
// 2. 設定進入狀態的樣式
// 3. 更新目前區域記錄
func (df *DragFeedback) ShowEnterFeedback(zoneID string) {
	df.indicator.SetText(fmt.Sprintf("放下到: %s", zoneID))
	df.indicator.Importance = widget.SuccessImportance
	df.currentZone = zoneID
	
	if !df.isVisible {
		df.isVisible = true
		df.overlay.Show()
	}
}

// ShowLeaveFeedback 顯示拖拽離開區域的視覺回饋
//
// 執行流程：
// 1. 重置指示器文字
// 2. 設定離開狀態的樣式
// 3. 清空目前區域記錄
func (df *DragFeedback) ShowLeaveFeedback() {
	df.indicator.SetText("拖拽中...")
	df.indicator.Importance = widget.MediumImportance
	df.currentZone = ""
}

// Hide 隱藏拖拽視覺回饋
//
// 執行流程：
// 1. 隱藏覆蓋層
// 2. 重置可見狀態
// 3. 清空指示器文字和目前區域
func (df *DragFeedback) Hide() {
	df.overlay.Hide()
	df.isVisible = false
	df.indicator.SetText("")
	df.currentZone = ""
}

// IsVisible 檢查視覺回饋是否可見
// 回傳：是否可見
func (df *DragFeedback) IsVisible() bool {
	return df.isVisible
}

// GetCurrentZone 取得目前拖拽區域
// 回傳：目前區域識別碼
func (df *DragFeedback) GetCurrentZone() string {
	return df.currentZone
}

// GetOverlay 取得覆蓋層容器
// 回傳：覆蓋層容器，用於添加到主視窗中
func (df *DragFeedback) GetOverlay() *fyne.Container {
	return df.overlay
}

// DragDropHelper 拖拽輔助工具
// 提供拖拽操作的輔助功能和工具方法
type DragDropHelper struct{}

// NewDragDropHelper 建立新的拖拽輔助工具
// 回傳：指向新建立的 DragDropHelper 的指標
func NewDragDropHelper() *DragDropHelper {
	return &DragDropHelper{}
}

// CreateFileDropZone 建立檔案拖拽區域
// 參數：title（區域標題）, targetPath（目標路徑）, acceptTypes（接受的檔案類型）
// 回傳：拖拽區域的 UI 元件
//
// 執行流程：
// 1. 建立拖拽區域的視覺元件
// 2. 設定區域標題和說明文字
// 3. 添加拖拽指示圖示
// 4. 回傳完整的拖拽區域元件
func (ddh *DragDropHelper) CreateFileDropZone(title, targetPath string, acceptTypes []string) *fyne.Container {
	// 建立標題標籤
	titleLabel := widget.NewLabel(title)
	titleLabel.TextStyle = fyne.TextStyle{Bold: true}
	titleLabel.Alignment = fyne.TextAlignCenter
	
	// 建立說明文字
	var acceptText string
	if len(acceptTypes) > 0 {
		acceptText = fmt.Sprintf("接受檔案類型: %s", strings.Join(acceptTypes, ", "))
	} else {
		acceptText = "接受所有檔案類型"
	}
	
	descLabel := widget.NewLabel(acceptText)
	descLabel.Alignment = fyne.TextAlignCenter
	
	// 建立拖拽圖示
	dragIcon := widget.NewIcon(theme.MoveDownIcon())
	
	// 建立路徑顯示
	pathLabel := widget.NewLabel(fmt.Sprintf("目標: %s", targetPath))
	pathLabel.Alignment = fyne.TextAlignCenter
	
	// 組合拖拽區域
	dropZone := container.NewVBox(
		titleLabel,
		dragIcon,
		descLabel,
		pathLabel,
	)
	
	// 設定拖拽區域樣式（概念性實作）
	// 實際樣式設定取決於 Fyne 的主題系統
	
	return dropZone
}

// CreateDragHandle 建立拖拽控制項
// 參數：sourcePath（來源路徑）
// 回傳：拖拽控制項的 UI 元件
//
// 執行流程：
// 1. 建立拖拽圖示按鈕
// 2. 設定拖拽開始事件
// 3. 添加視覺回饋
// 4. 回傳拖拽控制項
func (ddh *DragDropHelper) CreateDragHandle(sourcePath string) *widget.Button {
	// 建立拖拽按鈕
	dragButton := widget.NewButtonWithIcon("", theme.MoveUpIcon(), func() {
		// 開始拖拽操作（概念性實作）
		// 實際實作需要整合拖拽管理器
	})
	
	// 設定按鈕樣式
	dragButton.Importance = widget.LowImportance
	
	return dragButton
}

// ValidateDropOperation 驗證拖拽操作是否有效
// 參數：sourcePath（來源路徑）, targetPath（目標路徑）, operation（操作類型）
// 回傳：是否有效和錯誤訊息
//
// 執行流程：
// 1. 檢查來源和目標路徑的有效性
// 2. 驗證操作類型是否支援
// 3. 檢查檔案權限和衝突
// 4. 回傳驗證結果
func (ddh *DragDropHelper) ValidateDropOperation(sourcePath, targetPath string, operation DragOperation) (bool, string) {
	// 檢查路徑有效性
	if sourcePath == "" {
		return false, "來源路徑不能為空"
	}
	
	if targetPath == "" {
		return false, "目標路徑不能為空"
	}
	
	// 檢查是否為相同路徑
	if sourcePath == targetPath {
		return false, "來源和目標路徑不能相同"
	}
	
	// 檢查是否為子目錄移動（避免循環）
	if strings.HasPrefix(targetPath, sourcePath+string(filepath.Separator)) {
		return false, "不能將目錄移動到其子目錄中"
	}
	
	// 根據操作類型進行額外驗證
	switch operation {
	case DragOperationMove:
		// 移動操作的特殊驗證
		return true, ""
	case DragOperationCopy:
		// 複製操作的特殊驗證
		return true, ""
	case DragOperationLink:
		// 連結操作的特殊驗證
		return true, ""
	default:
		return false, "不支援的拖拽操作類型"
	}
}