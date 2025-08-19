// Package ui 包含使用者介面相關的元件和視窗管理
// 本檔案實作佈局管理器，負責管理主視窗的佈局結構和響應式設計
package ui

import (
	"fyne.io/fyne/v2"               // Fyne GUI 框架核心套件
	"fyne.io/fyne/v2/container"     // Fyne 容器佈局套件
	"fyne.io/fyne/v2/widget"        // Fyne UI 元件套件
	"fyne.io/fyne/v2/theme"         // Fyne 主題套件
)

// LayoutManager 負責管理主視窗的佈局結構
// 提供響應式佈局、面板大小調整和視圖模式切換功能
// 支援三欄式佈局的動態調整和優化
type LayoutManager struct {
	// 主要容器
	mainContainer    *fyne.Container      // 主要容器
	topBar          *fyne.Container      // 頂部工具欄容器
	contentArea     *fyne.Container      // 內容區域容器
	bottomBar       *fyne.Container      // 底部狀態欄容器
	
	// 佈局分割容器
	mainSplit       *container.Split     // 主要水平分割（側邊欄 | 內容）
	contentSplit    *container.Split     // 內容區域分割（筆記列表 | 編輯器）
	
	// 面板容器
	sidebarPanel    *fyne.Container      // 側邊欄面板
	noteListPanel   *fyne.Container      // 筆記列表面板
	editorPanel     *fyne.Container      // 編輯器面板
	
	// 工具欄
	quickToolbar    *widget.Toolbar      // 快速存取工具欄
	sideToolbar     *fyne.Container      // 側邊工具欄
	
	// 佈局狀態
	sidebarVisible  bool                 // 側邊欄是否可見
	noteListVisible bool                 // 筆記列表是否可見
	sidebarWidth    float64              // 側邊欄寬度比例
	noteListWidth   float64              // 筆記列表寬度比例
	
	// 響應式設定
	minWindowWidth  float32              // 最小視窗寬度
	minWindowHeight float32              // 最小視窗高度
	compactMode     bool                 // 緊湊模式（小螢幕）
	
	// 回調函數
	onLayoutChanged func(layout string)  // 佈局變更回調
	onPanelResized  func(panel string, size float64) // 面板大小變更回調
}

// NewLayoutManager 建立新的佈局管理器實例
// 回傳：指向新建立的 LayoutManager 的指標
//
// 執行流程：
// 1. 建立 LayoutManager 結構體實例
// 2. 設定預設的佈局參數和狀態
// 3. 初始化所有容器和工具欄
// 4. 設定響應式佈局參數
// 5. 建立完整的佈局結構
// 6. 回傳配置完成的佈局管理器實例
func NewLayoutManager() *LayoutManager {
	lm := &LayoutManager{
		// 預設佈局狀態
		sidebarVisible:  true,
		noteListVisible: true,
		sidebarWidth:    0.2,   // 側邊欄佔 20%
		noteListWidth:   0.25,  // 筆記列表佔 25%
		
		// 響應式設定
		minWindowWidth:  800,
		minWindowHeight: 600,
		compactMode:     false,
	}
	
	// 初始化佈局元件
	lm.initializeComponents()
	
	// 建立佈局結構
	lm.setupLayout()
	
	return lm
}

// initializeComponents 初始化所有佈局元件
// 建立所有必要的容器、工具欄和面板
//
// 執行流程：
// 1. 建立頂部工具欄和快速存取工具欄
// 2. 建立側邊工具欄和面板容器
// 3. 建立內容區域的各個面板
// 4. 建立底部狀態欄容器
func (lm *LayoutManager) initializeComponents() {
	// 建立頂部工具欄容器
	lm.topBar = container.NewVBox()
	
	// 建立快速存取工具欄
	lm.createQuickToolbar()
	
	// 建立側邊工具欄
	lm.createSideToolbar()
	
	// 建立面板容器
	lm.sidebarPanel = container.NewVBox()
	lm.noteListPanel = container.NewVBox()
	lm.editorPanel = container.NewVBox()
	
	// 建立底部狀態欄容器
	lm.bottomBar = container.NewHBox()
}

// createQuickToolbar 建立快速存取工具欄
// 包含最常用的功能按鈕，位於頂部
//
// 執行流程：
// 1. 建立檔案操作按鈕組
// 2. 建立視圖切換按鈕組
// 3. 建立搜尋和設定按鈕組
// 4. 組合所有按鈕到快速工具欄
func (lm *LayoutManager) createQuickToolbar() {
	lm.quickToolbar = widget.NewToolbar(
		// 檔案操作組
		widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {
			// 新增筆記
			lm.triggerAction("new_note")
		}),
		widget.NewToolbarAction(theme.FolderOpenIcon(), func() {
			// 開啟檔案
			lm.triggerAction("open_file")
		}),
		widget.NewToolbarAction(theme.DocumentSaveIcon(), func() {
			// 保存檔案
			lm.triggerAction("save_file")
		}),
		
		widget.NewToolbarSeparator(),
		
		// 視圖切換組
		widget.NewToolbarAction(theme.ListIcon(), func() {
			// 切換側邊欄
			lm.ToggleSidebar()
		}),
		widget.NewToolbarAction(theme.ViewRefreshIcon(), func() {
			// 切換筆記列表
			lm.ToggleNoteList()
		}),
		widget.NewToolbarAction(theme.VisibilityIcon(), func() {
			// 切換預覽
			lm.triggerAction("toggle_preview")
		}),
		
		widget.NewToolbarSeparator(),
		
		// 搜尋和設定組
		widget.NewToolbarAction(theme.SearchIcon(), func() {
			// 開啟搜尋
			lm.triggerAction("open_search")
		}),
		widget.NewToolbarAction(theme.SettingsIcon(), func() {
			// 開啟設定
			lm.triggerAction("open_settings")
		}),
	)
	
	// 將快速工具欄添加到頂部容器
	lm.topBar.Add(lm.quickToolbar)
}

// createSideToolbar 建立側邊工具欄
// 包含格式化和編輯功能按鈕，位於側邊
//
// 執行流程：
// 1. 建立格式化按鈕組
// 2. 建立插入功能按鈕組
// 3. 建立進階功能按鈕組
// 4. 使用垂直佈局組合側邊工具欄
func (lm *LayoutManager) createSideToolbar() {
	// 格式化按鈕組
	formatButtons := container.NewVBox(
		widget.NewLabel("格式"),
		widget.NewSeparator(),
		lm.createToolbarButton("𝐁", "粗體", "format_bold"),
		lm.createToolbarButton("𝐼", "斜體", "format_italic"),
		lm.createToolbarButton("𝐔", "底線", "format_underline"),
		lm.createToolbarButton("~~", "刪除線", "format_strikethrough"),
	)
	
	// 標題按鈕組
	headingButtons := container.NewVBox(
		widget.NewLabel("標題"),
		widget.NewSeparator(),
		lm.createToolbarButton("H1", "標題 1", "heading_1"),
		lm.createToolbarButton("H2", "標題 2", "heading_2"),
		lm.createToolbarButton("H3", "標題 3", "heading_3"),
	)
	
	// 列表按鈕組
	listButtons := container.NewVBox(
		widget.NewLabel("列表"),
		widget.NewSeparator(),
		lm.createToolbarButton("•", "項目符號", "list_bullet"),
		lm.createToolbarButton("1.", "編號列表", "list_numbered"),
		lm.createToolbarButton("☑", "待辦事項", "list_todo"),
	)
	
	// 插入按鈕組
	insertButtons := container.NewVBox(
		widget.NewLabel("插入"),
		widget.NewSeparator(),
		lm.createToolbarButton("🔗", "連結", "insert_link"),
		lm.createToolbarButton("🖼", "圖片", "insert_image"),
		lm.createToolbarButton("📊", "表格", "insert_table"),
		lm.createToolbarButton("💻", "程式碼", "insert_code"),
	)
	
	// 進階功能按鈕組
	advancedButtons := container.NewVBox(
		widget.NewLabel("進階"),
		widget.NewSeparator(),
		lm.createToolbarButton("🔒", "加密", "toggle_encryption"),
		lm.createToolbarButton("⭐", "最愛", "toggle_favorite"),
		lm.createToolbarButton("🏷", "標籤", "manage_tags"),
	)
	
	// 組合側邊工具欄
	lm.sideToolbar = container.NewVBox(
		formatButtons,
		widget.NewSeparator(),
		headingButtons,
		widget.NewSeparator(),
		listButtons,
		widget.NewSeparator(),
		insertButtons,
		widget.NewSeparator(),
		advancedButtons,
	)
}

// createToolbarButton 建立工具欄按鈕
// 參數：text（按鈕文字）、tooltip（提示文字）、action（動作名稱）
// 回傳：配置完成的按鈕元件
//
// 執行流程：
// 1. 建立按鈕元件並設定文字
// 2. 設定按鈕的點擊回調函數
// 3. 設定按鈕的樣式和提示文字
// 4. 回傳配置完成的按鈕
func (lm *LayoutManager) createToolbarButton(text, tooltip, action string) *widget.Button {
	button := widget.NewButton(text, func() {
		lm.triggerAction(action)
	})
	
	// 設定按鈕樣式為緊湊模式
	button.Resize(fyne.NewSize(40, 30))
	
	return button
}

// setupLayout 設定完整的佈局結構
// 組合所有容器和面板到主要佈局中
//
// 執行流程：
// 1. 建立內容區域的分割佈局
// 2. 建立主要的水平分割佈局
// 3. 組合頂部、內容和底部區域
// 4. 設定分割容器的初始比例
func (lm *LayoutManager) setupLayout() {
	// 建立內容區域分割（筆記列表 | 編輯器）
	lm.contentSplit = container.NewHSplit(
		lm.noteListPanel,
		lm.editorPanel,
	)
	lm.contentSplit.Offset = lm.noteListWidth / (1.0 - lm.sidebarWidth)
	
	// 建立主要分割（側邊欄 | 內容區域）
	// 側邊欄包含檔案樹和側邊工具欄
	sidebarContent := container.NewHSplit(
		lm.sidebarPanel,
		lm.sideToolbar,
	)
	sidebarContent.Offset = 0.7 // 檔案樹佔側邊欄的 70%
	
	lm.mainSplit = container.NewHSplit(
		sidebarContent,
		lm.contentSplit,
	)
	lm.mainSplit.Offset = lm.sidebarWidth
	
	// 建立內容區域容器
	lm.contentArea = container.NewVBox(lm.mainSplit)
	
	// 組合完整佈局
	lm.mainContainer = container.NewVBox(
		lm.topBar,      // 頂部工具欄
		lm.contentArea, // 主要內容區域
		lm.bottomBar,   // 底部狀態欄
	)
	
	// 更新佈局以反映當前狀態
	lm.updateLayout()
}

// GetContainer 取得主要佈局容器
// 回傳：主要佈局的 fyne.Container 實例
// 用於將佈局嵌入到主視窗中
func (lm *LayoutManager) GetContainer() *fyne.Container {
	return lm.mainContainer
}

// SetSidebarContent 設定側邊欄內容
// 參數：content（要設定的內容容器）
//
// 執行流程：
// 1. 清空側邊欄面板的現有內容
// 2. 添加新的內容到側邊欄面板
// 3. 刷新側邊欄面板顯示
func (lm *LayoutManager) SetSidebarContent(content *fyne.Container) {
	lm.sidebarPanel.Objects = []fyne.CanvasObject{content}
	lm.sidebarPanel.Refresh()
}

// SetNoteListContent 設定筆記列表內容
// 參數：content（要設定的內容容器）
//
// 執行流程：
// 1. 清空筆記列表面板的現有內容
// 2. 添加新的內容到筆記列表面板
// 3. 刷新筆記列表面板顯示
func (lm *LayoutManager) SetNoteListContent(content *fyne.Container) {
	lm.noteListPanel.Objects = []fyne.CanvasObject{content}
	lm.noteListPanel.Refresh()
}

// SetEditorContent 設定編輯器內容
// 參數：content（要設定的內容容器）
//
// 執行流程：
// 1. 清空編輯器面板的現有內容
// 2. 添加新的內容到編輯器面板
// 3. 刷新編輯器面板顯示
func (lm *LayoutManager) SetEditorContent(content *fyne.Container) {
	lm.editorPanel.Objects = []fyne.CanvasObject{content}
	lm.editorPanel.Refresh()
}

// SetStatusBarContent 設定狀態欄內容
// 參數：content（要設定的內容容器）
//
// 執行流程：
// 1. 清空底部狀態欄的現有內容
// 2. 添加新的內容到狀態欄
// 3. 刷新狀態欄顯示
func (lm *LayoutManager) SetStatusBarContent(content *fyne.Container) {
	lm.bottomBar.Objects = []fyne.CanvasObject{content}
	lm.bottomBar.Refresh()
}

// ToggleSidebar 切換側邊欄的顯示/隱藏
// 在顯示和隱藏側邊欄之間切換
//
// 執行流程：
// 1. 切換側邊欄可見性狀態
// 2. 更新佈局以反映變更
// 3. 觸發佈局變更回調
func (lm *LayoutManager) ToggleSidebar() {
	lm.sidebarVisible = !lm.sidebarVisible
	lm.updateLayout()
	
	if lm.onLayoutChanged != nil {
		if lm.sidebarVisible {
			lm.onLayoutChanged("sidebar_shown")
		} else {
			lm.onLayoutChanged("sidebar_hidden")
		}
	}
}

// ToggleNoteList 切換筆記列表的顯示/隱藏
// 在顯示和隱藏筆記列表之間切換
//
// 執行流程：
// 1. 切換筆記列表可見性狀態
// 2. 更新佈局以反映變更
// 3. 觸發佈局變更回調
func (lm *LayoutManager) ToggleNoteList() {
	lm.noteListVisible = !lm.noteListVisible
	lm.updateLayout()
	
	if lm.onLayoutChanged != nil {
		if lm.noteListVisible {
			lm.onLayoutChanged("notelist_shown")
		} else {
			lm.onLayoutChanged("notelist_hidden")
		}
	}
}

// SetSidebarWidth 設定側邊欄寬度比例
// 參數：width（寬度比例，0.0-1.0）
//
// 執行流程：
// 1. 驗證寬度比例的有效範圍
// 2. 更新側邊欄寬度比例
// 3. 更新分割容器的比例
// 4. 觸發面板大小變更回調
func (lm *LayoutManager) SetSidebarWidth(width float64) {
	if width < 0.1 {
		width = 0.1
	} else if width > 0.5 {
		width = 0.5
	}
	
	lm.sidebarWidth = width
	if lm.mainSplit != nil {
		lm.mainSplit.Offset = width
	}
	
	if lm.onPanelResized != nil {
		lm.onPanelResized("sidebar", width)
	}
}

// SetNoteListWidth 設定筆記列表寬度比例
// 參數：width（寬度比例，相對於內容區域）
//
// 執行流程：
// 1. 驗證寬度比例的有效範圍
// 2. 更新筆記列表寬度比例
// 3. 更新內容分割容器的比例
// 4. 觸發面板大小變更回調
func (lm *LayoutManager) SetNoteListWidth(width float64) {
	if width < 0.1 {
		width = 0.1
	} else if width > 0.8 {
		width = 0.8
	}
	
	lm.noteListWidth = width
	if lm.contentSplit != nil {
		// 計算相對於內容區域的比例
		contentWidth := 1.0 - lm.sidebarWidth
		lm.contentSplit.Offset = width / contentWidth
	}
	
	if lm.onPanelResized != nil {
		lm.onPanelResized("notelist", width)
	}
}

// updateLayout 更新佈局以反映當前狀態
// 根據面板可見性和大小設定調整佈局
//
// 執行流程：
// 1. 根據側邊欄可見性調整主分割容器
// 2. 根據筆記列表可見性調整內容分割容器
// 3. 更新分割比例
// 4. 刷新所有容器顯示
func (lm *LayoutManager) updateLayout() {
	if !lm.sidebarVisible {
		// 隱藏側邊欄，只顯示內容區域
		lm.mainSplit.Leading = container.NewWithoutLayout()
		lm.mainSplit.Offset = 0.0
	} else {
		// 顯示側邊欄
		sidebarContent := container.NewHSplit(
			lm.sidebarPanel,
			lm.sideToolbar,
		)
		sidebarContent.Offset = 0.7
		lm.mainSplit.Leading = sidebarContent
		lm.mainSplit.Offset = lm.sidebarWidth
	}
	
	if !lm.noteListVisible {
		// 隱藏筆記列表，只顯示編輯器
		lm.contentSplit.Leading = container.NewWithoutLayout()
		lm.contentSplit.Offset = 0.0
	} else {
		// 顯示筆記列表
		lm.contentSplit.Leading = lm.noteListPanel
		contentWidth := 1.0 - lm.sidebarWidth
		lm.contentSplit.Offset = lm.noteListWidth / contentWidth
	}
	
	// 刷新所有容器
	lm.mainContainer.Refresh()
}

// SetCompactMode 設定緊湊模式
// 參數：compact（是否啟用緊湊模式）
//
// 執行流程：
// 1. 更新緊湊模式狀態
// 2. 根據模式調整佈局參數
// 3. 更新佈局以反映變更
func (lm *LayoutManager) SetCompactMode(compact bool) {
	lm.compactMode = compact
	
	if compact {
		// 緊湊模式：調整面板大小和間距
		lm.SetSidebarWidth(0.15)
		lm.SetNoteListWidth(0.2)
	} else {
		// 正常模式：恢復預設大小
		lm.SetSidebarWidth(0.2)
		lm.SetNoteListWidth(0.25)
	}
}

// IsCompactMode 檢查是否為緊湊模式
// 回傳：是否為緊湊模式的布林值
func (lm *LayoutManager) IsCompactMode() bool {
	return lm.compactMode
}

// GetSidebarWidth 取得側邊欄寬度比例
// 回傳：側邊欄寬度比例（0.0-1.0）
func (lm *LayoutManager) GetSidebarWidth() float64 {
	return lm.sidebarWidth
}

// GetNoteListWidth 取得筆記列表寬度比例
// 回傳：筆記列表寬度比例
func (lm *LayoutManager) GetNoteListWidth() float64 {
	return lm.noteListWidth
}

// IsSidebarVisible 檢查側邊欄是否可見
// 回傳：側邊欄是否可見的布林值
func (lm *LayoutManager) IsSidebarVisible() bool {
	return lm.sidebarVisible
}

// IsNoteListVisible 檢查筆記列表是否可見
// 回傳：筆記列表是否可見的布林值
func (lm *LayoutManager) IsNoteListVisible() bool {
	return lm.noteListVisible
}

// SetOnLayoutChanged 設定佈局變更回調函數
// 參數：callback（佈局變更時的回調函數）
func (lm *LayoutManager) SetOnLayoutChanged(callback func(layout string)) {
	lm.onLayoutChanged = callback
}

// SetOnPanelResized 設定面板大小變更回調函數
// 參數：callback（面板大小變更時的回調函數）
func (lm *LayoutManager) SetOnPanelResized(callback func(panel string, size float64)) {
	lm.onPanelResized = callback
}

// triggerAction 觸發動作事件
// 參數：action（動作名稱）
//
// 執行流程：
// 1. 根據動作名稱執行相應的操作
// 2. 觸發相關的回調函數
// 3. 更新 UI 狀態（如果需要）
func (lm *LayoutManager) triggerAction(action string) {
	// 這裡可以實作動作分發邏輯
	// 實際的動作處理會由主視窗或其他元件負責
	if lm.onLayoutChanged != nil {
		lm.onLayoutChanged("action:" + action)
	}
}

// ResizeToWindow 根據視窗大小調整佈局
// 參數：windowSize（視窗大小）
//
// 執行流程：
// 1. 檢查視窗大小是否需要緊湊模式
// 2. 根據視窗大小調整面板比例
// 3. 更新佈局以適應新的視窗大小
func (lm *LayoutManager) ResizeToWindow(windowSize fyne.Size) {
	// 檢查是否需要緊湊模式
	needCompact := windowSize.Width < lm.minWindowWidth || windowSize.Height < lm.minWindowHeight
	
	if needCompact != lm.compactMode {
		lm.SetCompactMode(needCompact)
	}
	
	// 根據視窗寬度調整面板比例
	if windowSize.Width < 1000 {
		// 小視窗：減少側邊欄寬度
		lm.SetSidebarWidth(0.15)
		lm.SetNoteListWidth(0.2)
	} else if windowSize.Width > 1400 {
		// 大視窗：增加側邊欄寬度
		lm.SetSidebarWidth(0.25)
		lm.SetNoteListWidth(0.3)
	}
}

// SaveLayoutState 保存佈局狀態
// 回傳：佈局狀態的 map
// 用於保存使用者的佈局偏好設定
func (lm *LayoutManager) SaveLayoutState() map[string]interface{} {
	return map[string]interface{}{
		"sidebar_visible":   lm.sidebarVisible,
		"notelist_visible":  lm.noteListVisible,
		"sidebar_width":     lm.sidebarWidth,
		"notelist_width":    lm.noteListWidth,
		"compact_mode":      lm.compactMode,
	}
}

// LoadLayoutState 載入佈局狀態
// 參數：state（佈局狀態的 map）
//
// 執行流程：
// 1. 從狀態 map 中讀取各項設定
// 2. 應用設定到佈局管理器
// 3. 更新佈局以反映載入的狀態
func (lm *LayoutManager) LoadLayoutState(state map[string]interface{}) {
	if visible, ok := state["sidebar_visible"].(bool); ok {
		lm.sidebarVisible = visible
	}
	
	if visible, ok := state["notelist_visible"].(bool); ok {
		lm.noteListVisible = visible
	}
	
	if width, ok := state["sidebar_width"].(float64); ok {
		lm.sidebarWidth = width
	}
	
	if width, ok := state["notelist_width"].(float64); ok {
		lm.noteListWidth = width
	}
	
	if compact, ok := state["compact_mode"].(bool); ok {
		lm.compactMode = compact
	}
	
	// 更新佈局以反映載入的狀態
	lm.updateLayout()
}