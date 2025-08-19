// Package ui 包含使用者介面相關的元件和視窗管理
// 本檔案實作視圖管理器，負責管理編輯器和預覽面板的視圖模式切換
package ui

import (
	"fyne.io/fyne/v2"               // Fyne GUI 框架核心套件
	"fyne.io/fyne/v2/container"     // Fyne 容器佈局套件
	"fyne.io/fyne/v2/driver/desktop" // Fyne 桌面驅動套件
)

// ViewMode 代表視圖模式的類型
type ViewMode int

const (
	// ViewModeEdit 純編輯模式 - 只顯示編輯器
	ViewModeEdit ViewMode = iota
	// ViewModePreview 純預覽模式 - 只顯示預覽面板
	ViewModePreview
	// ViewModeSplit 分割視圖模式 - 並排顯示編輯器和預覽面板
	ViewModeSplit
)

// ViewState 代表視圖狀態資訊
// 用於保存和恢復視圖設定
type ViewState struct {
	Mode        ViewMode `json:"mode"`         // 當前視圖模式
	SplitRatio  float64  `json:"split_ratio"`  // 分割比例（0.0-1.0）
	IsFullscreen bool    `json:"is_fullscreen"` // 是否全螢幕模式
	SidebarHidden bool   `json:"sidebar_hidden"` // 側邊欄是否隱藏
	NoteListHidden bool  `json:"notelist_hidden"` // 筆記列表是否隱藏
}

// ViewManager 負責管理編輯器和預覽面板的視圖模式
// 提供視圖切換、全螢幕模式和狀態記憶功能
// 支援鍵盤快捷鍵和響應式佈局調整
type ViewManager struct {
	// 主要容器和元件
	container       *fyne.Container      // 主要容器
	editorContainer *fyne.Container      // 編輯器容器
	previewContainer *fyne.Container     // 預覽容器
	splitContainer  *container.Split     // 分割容器
	
	// 視圖狀態
	currentMode     ViewMode             // 當前視圖模式
	previousMode    ViewMode             // 上一個視圖模式
	splitRatio      float64              // 分割比例
	isFullscreen    bool                 // 是否全螢幕模式
	
	// 佈局狀態記憶
	savedSidebarVisible  bool            // 保存的側邊欄可見性
	savedNoteListVisible bool            // 保存的筆記列表可見性
	
	// 外部依賴
	layoutManager   *LayoutManager       // 佈局管理器
	window          fyne.Window          // 主視窗
	
	// 回調函數
	onViewModeChanged func(mode ViewMode) // 視圖模式變更回調
	onFullscreenToggled func(fullscreen bool) // 全螢幕切換回調
	onSplitRatioChanged func(ratio float64) // 分割比例變更回調
}

// NewViewManager 建立新的視圖管理器實例
// 參數：
//   - window: 主視窗實例
//   - layoutManager: 佈局管理器實例
// 回傳：指向新建立的 ViewManager 的指標
//
// 執行流程：
// 1. 建立 ViewManager 結構體實例
// 2. 設定預設的視圖狀態和參數
// 3. 初始化容器和佈局元件
// 4. 設定鍵盤快捷鍵監聽
// 5. 回傳配置完成的視圖管理器實例
func NewViewManager(window fyne.Window, layoutManager *LayoutManager) *ViewManager {
	vm := &ViewManager{
		window:        window,
		layoutManager: layoutManager,
		currentMode:   ViewModeSplit, // 預設為分割視圖模式
		previousMode:  ViewModeSplit,
		splitRatio:    0.5,           // 預設 50/50 分割
		isFullscreen:  false,
	}
	
	// 初始化容器
	vm.initializeContainers()
	
	// 設定鍵盤快捷鍵
	vm.setupKeyboardShortcuts()
	
	return vm
}

// initializeContainers 初始化所有容器元件
// 建立編輯器容器、預覽容器和分割容器
//
// 執行流程：
// 1. 建立編輯器和預覽容器
// 2. 建立分割容器並設定初始比例
// 3. 建立主要容器並設定預設佈局
func (vm *ViewManager) initializeContainers() {
	// 建立編輯器和預覽容器
	vm.editorContainer = container.NewVBox()
	vm.previewContainer = container.NewVBox()
	
	// 建立分割容器
	vm.splitContainer = container.NewHSplit(
		vm.editorContainer,
		vm.previewContainer,
	)
	vm.splitContainer.Offset = vm.splitRatio
	
	// 建立主要容器，預設使用分割視圖
	vm.container = container.NewVBox(vm.splitContainer)
}

// setupKeyboardShortcuts 設定鍵盤快捷鍵
// 註冊視圖切換的快捷鍵監聽器
//
// 執行流程：
// 1. 註冊 ⌘1 快捷鍵 - 切換到編輯模式
// 2. 註冊 ⌘2 快捷鍵 - 切換到預覽模式
// 3. 註冊 ⌘3 快捷鍵 - 切換到分割視圖模式
// 4. 註冊 ⌃⌘F 快捷鍵 - 切換全螢幕模式
func (vm *ViewManager) setupKeyboardShortcuts() {
	// 註冊快捷鍵到視窗的畫布
	canvas := vm.window.Canvas()
	
	// ⌘1 - 編輯模式
	canvas.AddShortcut(&desktop.CustomShortcut{
		Modifier: fyne.KeyModifierSuper,
		KeyName:  fyne.Key1,
	}, func(shortcut fyne.Shortcut) {
		vm.SetViewMode(ViewModeEdit)
	})
	
	// ⌘2 - 預覽模式
	canvas.AddShortcut(&desktop.CustomShortcut{
		Modifier: fyne.KeyModifierSuper,
		KeyName:  fyne.Key2,
	}, func(shortcut fyne.Shortcut) {
		vm.SetViewMode(ViewModePreview)
	})
	
	// ⌘3 - 分割視圖模式
	canvas.AddShortcut(&desktop.CustomShortcut{
		Modifier: fyne.KeyModifierSuper,
		KeyName:  fyne.Key3,
	}, func(shortcut fyne.Shortcut) {
		vm.SetViewMode(ViewModeSplit)
	})
	
	// ⌃⌘F - 切換全螢幕模式
	canvas.AddShortcut(&desktop.CustomShortcut{
		Modifier: fyne.KeyModifierSuper | fyne.KeyModifierControl,
		KeyName:  fyne.KeyF,
	}, func(shortcut fyne.Shortcut) {
		vm.ToggleFullscreen()
	})
}

// GetContainer 取得視圖管理器的主要容器
// 回傳：視圖管理器的 fyne.Container 實例
// 用於將視圖管理器嵌入到其他 UI 佈局中
func (vm *ViewManager) GetContainer() *fyne.Container {
	return vm.container
}

// SetEditorContent 設定編輯器內容
// 參數：content（編輯器內容容器）
//
// 執行流程：
// 1. 清空編輯器容器的現有內容
// 2. 添加新的編輯器內容
// 3. 刷新編輯器容器顯示
func (vm *ViewManager) SetEditorContent(content *fyne.Container) {
	vm.editorContainer.Objects = []fyne.CanvasObject{content}
	vm.editorContainer.Refresh()
}

// SetPreviewContent 設定預覽內容
// 參數：content（預覽內容容器）
//
// 執行流程：
// 1. 清空預覽容器的現有內容
// 2. 添加新的預覽內容
// 3. 刷新預覽容器顯示
func (vm *ViewManager) SetPreviewContent(content *fyne.Container) {
	vm.previewContainer.Objects = []fyne.CanvasObject{content}
	vm.previewContainer.Refresh()
}

// SetViewMode 設定視圖模式
// 參數：mode（要設定的視圖模式）
//
// 執行流程：
// 1. 檢查是否需要變更模式
// 2. 保存當前模式為上一個模式
// 3. 更新當前模式
// 4. 更新佈局以反映新模式
// 5. 觸發視圖模式變更回調
func (vm *ViewManager) SetViewMode(mode ViewMode) {
	if vm.currentMode == mode {
		return // 模式相同，無需變更
	}
	
	// 保存上一個模式
	vm.previousMode = vm.currentMode
	vm.currentMode = mode
	
	// 更新佈局
	vm.updateLayout()
	
	// 觸發回調
	if vm.onViewModeChanged != nil {
		vm.onViewModeChanged(mode)
	}
}

// GetViewMode 取得當前視圖模式
// 回傳：當前的視圖模式
func (vm *ViewManager) GetViewMode() ViewMode {
	return vm.currentMode
}

// ToggleViewMode 在視圖模式之間循環切換
// 按照 編輯 -> 預覽 -> 分割 -> 編輯 的順序循環
//
// 執行流程：
// 1. 根據當前模式決定下一個模式
// 2. 設定新的視圖模式
func (vm *ViewManager) ToggleViewMode() {
	var nextMode ViewMode
	
	switch vm.currentMode {
	case ViewModeEdit:
		nextMode = ViewModePreview
	case ViewModePreview:
		nextMode = ViewModeSplit
	case ViewModeSplit:
		nextMode = ViewModeEdit
	default:
		nextMode = ViewModeSplit
	}
	
	vm.SetViewMode(nextMode)
}

// SetSplitRatio 設定分割比例
// 參數：ratio（分割比例，0.0-1.0，左側/上方的比例）
//
// 執行流程：
// 1. 驗證比例的有效範圍
// 2. 更新分割比例
// 3. 應用到分割容器
// 4. 觸發分割比例變更回調
func (vm *ViewManager) SetSplitRatio(ratio float64) {
	if ratio < 0.1 {
		ratio = 0.1
	} else if ratio > 0.9 {
		ratio = 0.9
	}
	
	vm.splitRatio = ratio
	if vm.splitContainer != nil {
		vm.splitContainer.Offset = ratio
	}
	
	if vm.onSplitRatioChanged != nil {
		vm.onSplitRatioChanged(ratio)
	}
}

// GetSplitRatio 取得當前分割比例
// 回傳：當前分割比例（0.0-1.0）
func (vm *ViewManager) GetSplitRatio() float64 {
	return vm.splitRatio
}

// ToggleFullscreen 切換全螢幕模式
// 在全螢幕和視窗模式之間切換
//
// 執行流程：
// 1. 切換全螢幕狀態
// 2. 如果進入全螢幕，保存當前佈局狀態並隱藏側邊欄
// 3. 如果退出全螢幕，恢復之前的佈局狀態
// 4. 設定視窗的全螢幕狀態
// 5. 觸發全螢幕切換回調
func (vm *ViewManager) ToggleFullscreen() {
	vm.isFullscreen = !vm.isFullscreen
	
	if vm.isFullscreen {
		// 進入全螢幕模式
		vm.enterFullscreenMode()
	} else {
		// 退出全螢幕模式
		vm.exitFullscreenMode()
	}
	
	// 設定視窗全螢幕狀態
	vm.window.SetFullScreen(vm.isFullscreen)
	
	// 觸發回調
	if vm.onFullscreenToggled != nil {
		vm.onFullscreenToggled(vm.isFullscreen)
	}
}

// enterFullscreenMode 進入全螢幕模式
// 保存當前佈局狀態並隱藏不必要的 UI 元件
//
// 執行流程：
// 1. 保存當前的側邊欄和筆記列表可見性
// 2. 隱藏側邊欄和筆記列表以最大化編輯/預覽空間
// 3. 更新佈局以反映全螢幕狀態
func (vm *ViewManager) enterFullscreenMode() {
	// 保存當前佈局狀態
	vm.savedSidebarVisible = vm.layoutManager.IsSidebarVisible()
	vm.savedNoteListVisible = vm.layoutManager.IsNoteListVisible()
	
	// 隱藏側邊欄和筆記列表以最大化編輯空間
	if vm.savedSidebarVisible {
		vm.layoutManager.ToggleSidebar()
	}
	if vm.savedNoteListVisible {
		vm.layoutManager.ToggleNoteList()
	}
}

// exitFullscreenMode 退出全螢幕模式
// 恢復之前保存的佈局狀態
//
// 執行流程：
// 1. 根據保存的狀態恢復側邊欄可見性
// 2. 根據保存的狀態恢復筆記列表可見性
// 3. 更新佈局以反映正常視窗狀態
func (vm *ViewManager) exitFullscreenMode() {
	// 恢復之前的佈局狀態
	if vm.savedSidebarVisible && !vm.layoutManager.IsSidebarVisible() {
		vm.layoutManager.ToggleSidebar()
	}
	if vm.savedNoteListVisible && !vm.layoutManager.IsNoteListVisible() {
		vm.layoutManager.ToggleNoteList()
	}
}

// IsFullscreen 檢查是否為全螢幕模式
// 回傳：是否為全螢幕模式的布林值
func (vm *ViewManager) IsFullscreen() bool {
	return vm.isFullscreen
}

// updateLayout 更新佈局以反映當前視圖模式
// 根據當前視圖模式調整容器內容
//
// 執行流程：
// 1. 根據當前視圖模式選擇適當的佈局
// 2. 更新主要容器的內容
// 3. 刷新容器顯示
func (vm *ViewManager) updateLayout() {
	switch vm.currentMode {
	case ViewModeEdit:
		// 純編輯模式 - 只顯示編輯器
		vm.container.Objects = []fyne.CanvasObject{vm.editorContainer}
		
	case ViewModePreview:
		// 純預覽模式 - 只顯示預覽面板
		vm.container.Objects = []fyne.CanvasObject{vm.previewContainer}
		
	case ViewModeSplit:
		// 分割視圖模式 - 並排顯示編輯器和預覽面板
		vm.splitContainer.Offset = vm.splitRatio
		vm.container.Objects = []fyne.CanvasObject{vm.splitContainer}
		
	default:
		// 預設使用分割視圖模式
		vm.container.Objects = []fyne.CanvasObject{vm.splitContainer}
	}
	
	vm.container.Refresh()
}

// GetPreviousMode 取得上一個視圖模式
// 回傳：上一個視圖模式
// 用於實作「返回上一個視圖」功能
func (vm *ViewManager) GetPreviousMode() ViewMode {
	return vm.previousMode
}

// RestorePreviousMode 恢復到上一個視圖模式
// 切換回上一個使用的視圖模式
//
// 執行流程：
// 1. 取得上一個視圖模式
// 2. 設定為當前視圖模式
func (vm *ViewManager) RestorePreviousMode() {
	previousMode := vm.previousMode
	vm.SetViewMode(previousMode)
}

// SaveViewState 保存當前視圖狀態
// 回傳：包含當前視圖狀態的 ViewState 結構體
// 用於保存使用者的視圖偏好設定
func (vm *ViewManager) SaveViewState() ViewState {
	return ViewState{
		Mode:           vm.currentMode,
		SplitRatio:     vm.splitRatio,
		IsFullscreen:   vm.isFullscreen,
		SidebarHidden:  !vm.layoutManager.IsSidebarVisible(),
		NoteListHidden: !vm.layoutManager.IsNoteListVisible(),
	}
}

// LoadViewState 載入視圖狀態
// 參數：state（要載入的視圖狀態）
//
// 執行流程：
// 1. 從狀態結構體中讀取各項設定
// 2. 應用視圖模式設定
// 3. 應用分割比例設定
// 4. 應用全螢幕狀態設定
// 5. 更新佈局以反映載入的狀態
func (vm *ViewManager) LoadViewState(state ViewState) {
	// 載入視圖模式
	vm.SetViewMode(state.Mode)
	
	// 載入分割比例
	vm.SetSplitRatio(state.SplitRatio)
	
	// 載入全螢幕狀態
	if state.IsFullscreen != vm.isFullscreen {
		vm.ToggleFullscreen()
	}
	
	// 載入側邊欄和筆記列表狀態
	if state.SidebarHidden && vm.layoutManager.IsSidebarVisible() {
		vm.layoutManager.ToggleSidebar()
	} else if !state.SidebarHidden && !vm.layoutManager.IsSidebarVisible() {
		vm.layoutManager.ToggleSidebar()
	}
	
	if state.NoteListHidden && vm.layoutManager.IsNoteListVisible() {
		vm.layoutManager.ToggleNoteList()
	} else if !state.NoteListHidden && !vm.layoutManager.IsNoteListVisible() {
		vm.layoutManager.ToggleNoteList()
	}
}

// GetViewModeString 取得視圖模式的字串表示
// 參數：mode（視圖模式）
// 回傳：視圖模式的中文字串描述
func (vm *ViewManager) GetViewModeString(mode ViewMode) string {
	switch mode {
	case ViewModeEdit:
		return "編輯模式"
	case ViewModePreview:
		return "預覽模式"
	case ViewModeSplit:
		return "分割視圖"
	default:
		return "未知模式"
	}
}

// GetCurrentViewModeString 取得當前視圖模式的字串表示
// 回傳：當前視圖模式的中文字串描述
func (vm *ViewManager) GetCurrentViewModeString() string {
	return vm.GetViewModeString(vm.currentMode)
}

// SetOnViewModeChanged 設定視圖模式變更回調函數
// 參數：callback（視圖模式變更時的回調函數）
func (vm *ViewManager) SetOnViewModeChanged(callback func(mode ViewMode)) {
	vm.onViewModeChanged = callback
}

// SetOnFullscreenToggled 設定全螢幕切換回調函數
// 參數：callback（全螢幕切換時的回調函數）
func (vm *ViewManager) SetOnFullscreenToggled(callback func(fullscreen bool)) {
	vm.onFullscreenToggled = callback
}

// SetOnSplitRatioChanged 設定分割比例變更回調函數
// 參數：callback（分割比例變更時的回調函數）
func (vm *ViewManager) SetOnSplitRatioChanged(callback func(ratio float64)) {
	vm.onSplitRatioChanged = callback
}

// IsEditMode 檢查是否為編輯模式
// 回傳：是否為編輯模式的布林值
func (vm *ViewManager) IsEditMode() bool {
	return vm.currentMode == ViewModeEdit
}

// IsPreviewMode 檢查是否為預覽模式
// 回傳：是否為預覽模式的布林值
func (vm *ViewManager) IsPreviewMode() bool {
	return vm.currentMode == ViewModePreview
}

// IsSplitMode 檢查是否為分割視圖模式
// 回傳：是否為分割視圖模式的布林值
func (vm *ViewManager) IsSplitMode() bool {
	return vm.currentMode == ViewModeSplit
}

// GetAvailableViewModes 取得所有可用的視圖模式
// 回傳：視圖模式的切片
func (vm *ViewManager) GetAvailableViewModes() []ViewMode {
	return []ViewMode{ViewModeEdit, ViewModePreview, ViewModeSplit}
}

// GetViewModeShortcut 取得視圖模式的快捷鍵描述
// 參數：mode（視圖模式）
// 回傳：快捷鍵的字串描述
func (vm *ViewManager) GetViewModeShortcut(mode ViewMode) string {
	switch mode {
	case ViewModeEdit:
		return "⌘1"
	case ViewModePreview:
		return "⌘2"
	case ViewModeSplit:
		return "⌘3"
	default:
		return ""
	}
}