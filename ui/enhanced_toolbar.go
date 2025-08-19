// Package ui 包含使用者介面相關的元件和視窗管理
// 本檔案實作增強版工具欄系統，提供更豐富的編輯功能和更好的使用者體驗
package ui

import (
	"fyne.io/fyne/v2"               // Fyne GUI 框架核心套件
	"fyne.io/fyne/v2/container"     // Fyne 容器佈局套件
	"fyne.io/fyne/v2/widget"        // Fyne UI 元件套件
	"fyne.io/fyne/v2/theme"         // Fyne 主題套件
)

// EnhancedToolbar 代表增強版工具欄系統
// 提供分類的工具按鈕、快速存取功能和自訂工具欄配置
// 支援響應式佈局和工具欄的動態調整
type EnhancedToolbar struct {
	// 主要容器
	container       *fyne.Container      // 主要工具欄容器
	
	// 工具欄區段
	fileSection     *fyne.Container      // 檔案操作區段
	editSection     *fyne.Container      // 編輯功能區段
	formatSection   *fyne.Container      // 格式化區段
	insertSection   *fyne.Container      // 插入功能區段
	viewSection     *fyne.Container      // 視圖控制區段
	toolsSection    *fyne.Container      // 工具和設定區段
	
	// 工具按鈕
	buttons         map[string]*widget.Button // 所有工具按鈕的映射
	
	// 狀態和設定
	compactMode     bool                 // 緊湊模式
	sectionsVisible map[string]bool      // 各區段的可見性
	
	// 回調函數
	onActionTriggered func(action string, params map[string]interface{}) // 動作觸發回調
}

// NewEnhancedToolbar 建立新的增強版工具欄實例
// 回傳：指向新建立的 EnhancedToolbar 的指標
//
// 執行流程：
// 1. 建立 EnhancedToolbar 結構體實例
// 2. 初始化按鈕映射和區段可見性設定
// 3. 建立所有工具欄區段
// 4. 組合完整的工具欄佈局
// 5. 回傳配置完成的工具欄實例
func NewEnhancedToolbar() *EnhancedToolbar {
	et := &EnhancedToolbar{
		buttons:         make(map[string]*widget.Button),
		sectionsVisible: make(map[string]bool),
		compactMode:     false,
	}
	
	// 設定預設的區段可見性
	et.sectionsVisible["file"] = true
	et.sectionsVisible["edit"] = true
	et.sectionsVisible["format"] = true
	et.sectionsVisible["insert"] = true
	et.sectionsVisible["view"] = true
	et.sectionsVisible["tools"] = true
	
	// 建立所有工具欄區段
	et.createSections()
	
	// 組合工具欄佈局
	et.setupLayout()
	
	return et
}

// createSections 建立所有工具欄區段
// 建立並配置各個功能區段的按鈕和佈局
//
// 執行流程：
// 1. 建立檔案操作區段
// 2. 建立編輯功能區段
// 3. 建立格式化區段
// 4. 建立插入功能區段
// 5. 建立視圖控制區段
// 6. 建立工具和設定區段
func (et *EnhancedToolbar) createSections() {
	et.createFileSection()
	et.createEditSection()
	et.createFormatSection()
	et.createInsertSection()
	et.createViewSection()
	et.createToolsSection()
}

// createFileSection 建立檔案操作區段
// 包含新增、開啟、保存、匯入、匯出等檔案相關功能
//
// 執行流程：
// 1. 建立檔案操作按鈕
// 2. 設定按鈕的圖示和回調函數
// 3. 組合按鈕到檔案區段容器
func (et *EnhancedToolbar) createFileSection() {
	// 新增筆記按鈕
	newNoteBtn := et.createButton("new_note", theme.DocumentCreateIcon(), "新增筆記 (⌘N)", func() {
		et.triggerAction("new_note", nil)
	})
	
	// 新增資料夾按鈕
	newFolderBtn := et.createButton("new_folder", theme.FolderNewIcon(), "新增資料夾 (⇧⌘N)", func() {
		et.triggerAction("new_folder", nil)
	})
	
	// 開啟檔案按鈕
	openBtn := et.createButton("open_file", theme.FolderOpenIcon(), "開啟檔案 (⌘O)", func() {
		et.triggerAction("open_file", nil)
	})
	
	// 保存按鈕
	saveBtn := et.createButton("save_file", theme.DocumentSaveIcon(), "保存 (⌘S)", func() {
		et.triggerAction("save_file", nil)
	})
	
	// 另存新檔按鈕
	saveAsBtn := et.createButton("save_as", theme.DocumentSaveIcon(), "另存新檔 (⇧⌘S)", func() {
		et.triggerAction("save_as", nil)
	})
	
	// 匯入按鈕
	importBtn := et.createButton("import", theme.MailAttachmentIcon(), "匯入檔案 (⌘I)", func() {
		et.triggerAction("import_file", nil)
	})
	
	// 匯出按鈕
	exportBtn := et.createButton("export", theme.MailSendIcon(), "匯出檔案 (⌘E)", func() {
		et.triggerAction("export_file", nil)
	})
	
	// 組合檔案區段
	et.fileSection = container.NewHBox(
		newNoteBtn,
		newFolderBtn,
		widget.NewSeparator(),
		openBtn,
		saveBtn,
		saveAsBtn,
		widget.NewSeparator(),
		importBtn,
		exportBtn,
	)
}

// createEditSection 建立編輯功能區段
// 包含復原、重做、剪下、複製、貼上、尋找、取代等編輯功能
//
// 執行流程：
// 1. 建立編輯操作按鈕
// 2. 設定按鈕的圖示和快捷鍵
// 3. 組合按鈕到編輯區段容器
func (et *EnhancedToolbar) createEditSection() {
	// 復原按鈕
	undoBtn := et.createButton("undo", theme.NavigateBackIcon(), "復原 (⌘Z)", func() {
		et.triggerAction("undo", nil)
	})
	
	// 重做按鈕
	redoBtn := et.createButton("redo", theme.NavigateNextIcon(), "重做 (⇧⌘Z)", func() {
		et.triggerAction("redo", nil)
	})
	
	// 剪下按鈕
	cutBtn := et.createButton("cut", theme.ContentCutIcon(), "剪下 (⌘X)", func() {
		et.triggerAction("cut", nil)
	})
	
	// 複製按鈕
	copyBtn := et.createButton("copy", theme.ContentCopyIcon(), "複製 (⌘C)", func() {
		et.triggerAction("copy", nil)
	})
	
	// 貼上按鈕
	pasteBtn := et.createButton("paste", theme.ContentPasteIcon(), "貼上 (⌘V)", func() {
		et.triggerAction("paste", nil)
	})
	
	// 尋找按鈕
	findBtn := et.createButton("find", theme.SearchIcon(), "尋找 (⌘F)", func() {
		et.triggerAction("find", nil)
	})
	
	// 取代按鈕
	replaceBtn := et.createButton("replace", theme.SearchReplaceIcon(), "取代 (⌥⌘F)", func() {
		et.triggerAction("replace", nil)
	})
	
	// 組合編輯區段
	et.editSection = container.NewHBox(
		undoBtn,
		redoBtn,
		widget.NewSeparator(),
		cutBtn,
		copyBtn,
		pasteBtn,
		widget.NewSeparator(),
		findBtn,
		replaceBtn,
	)
}

// createFormatSection 建立格式化區段
// 包含粗體、斜體、底線、刪除線、標題等格式化功能
//
// 執行流程：
// 1. 建立文字格式化按鈕
// 2. 建立標題格式化按鈕
// 3. 組合按鈕到格式化區段容器
func (et *EnhancedToolbar) createFormatSection() {
	// 粗體按鈕
	boldBtn := et.createButton("format_bold", theme.ContentCopyIcon(), "粗體 (⌘B)", func() {
		et.triggerAction("format_bold", nil)
	})
	boldBtn.Text = "𝐁"
	
	// 斜體按鈕
	italicBtn := et.createButton("format_italic", theme.ContentCopyIcon(), "斜體 (⌘I)", func() {
		et.triggerAction("format_italic", nil)
	})
	italicBtn.Text = "𝐼"
	
	// 底線按鈕
	underlineBtn := et.createButton("format_underline", theme.ContentCopyIcon(), "底線 (⌘U)", func() {
		et.triggerAction("format_underline", nil)
	})
	underlineBtn.Text = "𝐔"
	
	// 刪除線按鈕
	strikeBtn := et.createButton("format_strikethrough", theme.ContentCopyIcon(), "刪除線", func() {
		et.triggerAction("format_strikethrough", nil)
	})
	strikeBtn.Text = "~~"
	
	// 標題 1 按鈕
	h1Btn := et.createButton("heading_1", theme.DocumentIcon(), "標題 1", func() {
		et.triggerAction("heading_1", nil)
	})
	h1Btn.Text = "H1"
	
	// 標題 2 按鈕
	h2Btn := et.createButton("heading_2", theme.DocumentIcon(), "標題 2", func() {
		et.triggerAction("heading_2", nil)
	})
	h2Btn.Text = "H2"
	
	// 標題 3 按鈕
	h3Btn := et.createButton("heading_3", theme.DocumentIcon(), "標題 3", func() {
		et.triggerAction("heading_3", nil)
	})
	h3Btn.Text = "H3"
	
	// 組合格式化區段
	et.formatSection = container.NewHBox(
		boldBtn,
		italicBtn,
		underlineBtn,
		strikeBtn,
		widget.NewSeparator(),
		h1Btn,
		h2Btn,
		h3Btn,
	)
}

// createInsertSection 建立插入功能區段
// 包含連結、圖片、表格、程式碼、列表等插入功能
//
// 執行流程：
// 1. 建立插入元素按鈕
// 2. 建立列表和程式碼按鈕
// 3. 組合按鈕到插入區段容器
func (et *EnhancedToolbar) createInsertSection() {
	// 連結按鈕
	linkBtn := et.createButton("insert_link", theme.ContentCopyIcon(), "插入連結 (⌘K)", func() {
		et.triggerAction("insert_link", nil)
	})
	linkBtn.Text = "🔗"
	
	// 圖片按鈕
	imageBtn := et.createButton("insert_image", theme.ContentCopyIcon(), "插入圖片", func() {
		et.triggerAction("insert_image", nil)
	})
	imageBtn.Text = "🖼"
	
	// 表格按鈕
	tableBtn := et.createButton("insert_table", theme.ContentCopyIcon(), "插入表格", func() {
		et.triggerAction("insert_table", nil)
	})
	tableBtn.Text = "📊"
	
	// 程式碼按鈕
	codeBtn := et.createButton("insert_code", theme.ContentCopyIcon(), "插入程式碼", func() {
		et.triggerAction("insert_code", nil)
	})
	codeBtn.Text = "💻"
	
	// 項目符號列表按鈕
	bulletBtn := et.createButton("list_bullet", theme.ListIcon(), "項目符號列表", func() {
		et.triggerAction("list_bullet", nil)
	})
	
	// 編號列表按鈕
	numberedBtn := et.createButton("list_numbered", theme.ListIcon(), "編號列表", func() {
		et.triggerAction("list_numbered", nil)
	})
	
	// 待辦事項按鈕
	todoBtn := et.createButton("list_todo", theme.ContentCopyIcon(), "待辦事項", func() {
		et.triggerAction("list_todo", nil)
	})
	todoBtn.Text = "☑"
	
	// 組合插入區段
	et.insertSection = container.NewHBox(
		linkBtn,
		imageBtn,
		tableBtn,
		codeBtn,
		widget.NewSeparator(),
		bulletBtn,
		numberedBtn,
		todoBtn,
	)
}

// createViewSection 建立視圖控制區段
// 包含預覽切換、全螢幕、縮放、主題切換等視圖功能
//
// 執行流程：
// 1. 建立視圖模式切換按鈕
// 2. 建立縮放和主題按鈕
// 3. 組合按鈕到視圖區段容器
func (et *EnhancedToolbar) createViewSection() {
	// 預覽切換按鈕
	previewBtn := et.createButton("toggle_preview", theme.VisibilityIcon(), "切換預覽 (⌘3)", func() {
		et.triggerAction("toggle_preview", nil)
	})
	
	// 編輯模式按鈕
	editModeBtn := et.createButton("edit_mode", theme.DocumentCreateIcon(), "編輯模式 (⌘1)", func() {
		et.triggerAction("edit_mode", nil)
	})
	
	// 預覽模式按鈕
	previewModeBtn := et.createButton("preview_mode", theme.VisibilityIcon(), "預覽模式 (⌘2)", func() {
		et.triggerAction("preview_mode", nil)
	})
	
	// 分割視圖按鈕
	splitViewBtn := et.createButton("split_view", theme.ViewRefreshIcon(), "分割視圖", func() {
		et.triggerAction("split_view", nil)
	})
	
	// 全螢幕按鈕
	fullscreenBtn := et.createButton("fullscreen", theme.ViewFullScreenIcon(), "全螢幕 (⌃⌘F)", func() {
		et.triggerAction("fullscreen", nil)
	})
	
	// 縮放放大按鈕
	zoomInBtn := et.createButton("zoom_in", theme.ContentAddIcon(), "放大", func() {
		et.triggerAction("zoom_in", nil)
	})
	
	// 縮放縮小按鈕
	zoomOutBtn := et.createButton("zoom_out", theme.ContentRemoveIcon(), "縮小", func() {
		et.triggerAction("zoom_out", nil)
	})
	
	// 主題切換按鈕
	themeBtn := et.createButton("toggle_theme", theme.ColorPaletteIcon(), "切換主題 (⌘D)", func() {
		et.triggerAction("toggle_theme", nil)
	})
	
	// 組合視圖區段
	et.viewSection = container.NewHBox(
		editModeBtn,
		previewModeBtn,
		splitViewBtn,
		previewBtn,
		widget.NewSeparator(),
		fullscreenBtn,
		zoomInBtn,
		zoomOutBtn,
		widget.NewSeparator(),
		themeBtn,
	)
}

// createToolsSection 建立工具和設定區段
// 包含加密、標籤、統計、設定等工具功能
//
// 執行流程：
// 1. 建立安全和標籤按鈕
// 2. 建立統計和設定按鈕
// 3. 組合按鈕到工具區段容器
func (et *EnhancedToolbar) createToolsSection() {
	// 加密切換按鈕
	encryptBtn := et.createButton("toggle_encryption", theme.VisibilityOffIcon(), "切換加密", func() {
		et.triggerAction("toggle_encryption", nil)
	})
	
	// 最愛切換按鈕
	favoriteBtn := et.createButton("toggle_favorite", theme.ContentCopyIcon(), "切換最愛", func() {
		et.triggerAction("toggle_favorite", nil)
	})
	favoriteBtn.Text = "⭐"
	
	// 標籤管理按鈕
	tagsBtn := et.createButton("manage_tags", theme.ContentCopyIcon(), "管理標籤", func() {
		et.triggerAction("manage_tags", nil)
	})
	tagsBtn.Text = "🏷"
	
	// 統計資訊按鈕
	statsBtn := et.createButton("show_stats", theme.InfoIcon(), "統計資訊", func() {
		et.triggerAction("show_stats", nil)
	})
	
	// 設定按鈕
	settingsBtn := et.createButton("open_settings", theme.SettingsIcon(), "設定 (⌘,)", func() {
		et.triggerAction("open_settings", nil)
	})
	
	// 說明按鈕
	helpBtn := et.createButton("show_help", theme.HelpIcon(), "說明", func() {
		et.triggerAction("show_help", nil)
	})
	
	// 組合工具區段
	et.toolsSection = container.NewHBox(
		encryptBtn,
		favoriteBtn,
		tagsBtn,
		widget.NewSeparator(),
		statsBtn,
		settingsBtn,
		helpBtn,
	)
}

// createButton 建立工具欄按鈕
// 參數：id（按鈕 ID）、icon（圖示）、tooltip（提示文字）、callback（回調函數）
// 回傳：配置完成的按鈕元件
//
// 執行流程：
// 1. 建立按鈕元件並設定圖示
// 2. 設定按鈕的回調函數
// 3. 將按鈕添加到按鈕映射中
// 4. 回傳配置完成的按鈕
func (et *EnhancedToolbar) createButton(id string, icon fyne.Resource, tooltip string, callback func()) *widget.Button {
	button := widget.NewButtonWithIcon("", icon, callback)
	
	// 設定按鈕樣式
	if et.compactMode {
		button.Resize(fyne.NewSize(32, 32))
	} else {
		button.Resize(fyne.NewSize(40, 40))
	}
	
	// 將按鈕添加到映射中
	et.buttons[id] = button
	
	return button
}

// setupLayout 設定工具欄佈局
// 組合所有區段到主要工具欄容器中
//
// 執行流程：
// 1. 建立主要工具欄容器
// 2. 根據區段可見性添加區段
// 3. 在區段之間添加分隔線
func (et *EnhancedToolbar) setupLayout() {
	var sections []*fyne.Container
	
	// 根據可見性添加區段
	if et.sectionsVisible["file"] {
		sections = append(sections, et.fileSection)
	}
	
	if et.sectionsVisible["edit"] {
		sections = append(sections, et.editSection)
	}
	
	if et.sectionsVisible["format"] {
		sections = append(sections, et.formatSection)
	}
	
	if et.sectionsVisible["insert"] {
		sections = append(sections, et.insertSection)
	}
	
	if et.sectionsVisible["view"] {
		sections = append(sections, et.viewSection)
	}
	
	if et.sectionsVisible["tools"] {
		sections = append(sections, et.toolsSection)
	}
	
	// 建立主要容器並添加區段
	var objects []fyne.CanvasObject
	for i, section := range sections {
		if i > 0 {
			// 在區段之間添加分隔線
			objects = append(objects, widget.NewSeparator())
		}
		objects = append(objects, section)
	}
	
	et.container = container.NewHBox(objects...)
}

// GetContainer 取得工具欄容器
// 回傳：工具欄的 fyne.Container 實例
func (et *EnhancedToolbar) GetContainer() *fyne.Container {
	return et.container
}

// SetSectionVisible 設定區段可見性
// 參數：section（區段名稱）、visible（是否可見）
//
// 執行流程：
// 1. 更新區段可見性設定
// 2. 重新建立佈局以反映變更
func (et *EnhancedToolbar) SetSectionVisible(section string, visible bool) {
	et.sectionsVisible[section] = visible
	et.setupLayout()
}

// IsSectionVisible 檢查區段是否可見
// 參數：section（區段名稱）
// 回傳：區段是否可見的布林值
func (et *EnhancedToolbar) IsSectionVisible(section string) bool {
	return et.sectionsVisible[section]
}

// SetCompactMode 設定緊湊模式
// 參數：compact（是否啟用緊湊模式）
//
// 執行流程：
// 1. 更新緊湊模式狀態
// 2. 調整所有按鈕的大小
// 3. 重新建立佈局
func (et *EnhancedToolbar) SetCompactMode(compact bool) {
	et.compactMode = compact
	
	// 調整按鈕大小
	buttonSize := fyne.NewSize(40, 40)
	if compact {
		buttonSize = fyne.NewSize(32, 32)
	}
	
	for _, button := range et.buttons {
		button.Resize(buttonSize)
	}
	
	et.setupLayout()
}

// IsCompactMode 檢查是否為緊湊模式
// 回傳：是否為緊湊模式的布林值
func (et *EnhancedToolbar) IsCompactMode() bool {
	return et.compactMode
}

// EnableButton 啟用按鈕
// 參數：buttonId（按鈕 ID）
func (et *EnhancedToolbar) EnableButton(buttonId string) {
	if button, exists := et.buttons[buttonId]; exists {
		button.Enable()
	}
}

// DisableButton 停用按鈕
// 參數：buttonId（按鈕 ID）
func (et *EnhancedToolbar) DisableButton(buttonId string) {
	if button, exists := et.buttons[buttonId]; exists {
		button.Disable()
	}
}

// SetButtonText 設定按鈕文字
// 參數：buttonId（按鈕 ID）、text（按鈕文字）
func (et *EnhancedToolbar) SetButtonText(buttonId, text string) {
	if button, exists := et.buttons[buttonId]; exists {
		button.SetText(text)
	}
}

// GetButton 取得按鈕實例
// 參數：buttonId（按鈕 ID）
// 回傳：按鈕實例，如果不存在則回傳 nil
func (et *EnhancedToolbar) GetButton(buttonId string) *widget.Button {
	return et.buttons[buttonId]
}

// SetOnActionTriggered 設定動作觸發回調函數
// 參數：callback（動作觸發時的回調函數）
func (et *EnhancedToolbar) SetOnActionTriggered(callback func(action string, params map[string]interface{})) {
	et.onActionTriggered = callback
}

// triggerAction 觸發動作事件
// 參數：action（動作名稱）、params（動作參數）
//
// 執行流程：
// 1. 檢查是否有設定回調函數
// 2. 調用回調函數並傳遞動作和參數
func (et *EnhancedToolbar) triggerAction(action string, params map[string]interface{}) {
	if et.onActionTriggered != nil {
		et.onActionTriggered(action, params)
	}
}

// GetAvailableSections 取得所有可用的區段名稱
// 回傳：區段名稱的切片
func (et *EnhancedToolbar) GetAvailableSections() []string {
	return []string{"file", "edit", "format", "insert", "view", "tools"}
}

// GetSectionButtons 取得指定區段的所有按鈕 ID
// 參數：section（區段名稱）
// 回傳：按鈕 ID 的切片
func (et *EnhancedToolbar) GetSectionButtons(section string) []string {
	switch section {
	case "file":
		return []string{"new_note", "new_folder", "open_file", "save_file", "save_as", "import", "export"}
	case "edit":
		return []string{"undo", "redo", "cut", "copy", "paste", "find", "replace"}
	case "format":
		return []string{"format_bold", "format_italic", "format_underline", "format_strikethrough", "heading_1", "heading_2", "heading_3"}
	case "insert":
		return []string{"insert_link", "insert_image", "insert_table", "insert_code", "list_bullet", "list_numbered", "list_todo"}
	case "view":
		return []string{"toggle_preview", "edit_mode", "preview_mode", "split_view", "fullscreen", "zoom_in", "zoom_out", "toggle_theme"}
	case "tools":
		return []string{"toggle_encryption", "toggle_favorite", "manage_tags", "show_stats", "open_settings", "show_help"}
	default:
		return []string{}
	}
}