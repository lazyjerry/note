// Package ui 包含使用者介面相關的元件和視窗管理
// 本檔案包含增強版工具欄的單元測試
package ui

import (
	"testing"
)

// TestNewEnhancedToolbar 測試增強版工具欄的建立
// 驗證工具欄是否正確初始化所有區段和按鈕
func TestNewEnhancedToolbar(t *testing.T) {
	// 建立增強版工具欄
	toolbar := NewEnhancedToolbar()
	
	// 驗證工具欄不為 nil
	if toolbar == nil {
		t.Fatal("增強版工具欄不應該為 nil")
	}
	
	// 驗證主要容器已建立
	if toolbar.container == nil {
		t.Error("工具欄容器不應該為 nil")
	}
	
	// 驗證各區段已建立
	if toolbar.fileSection == nil {
		t.Error("檔案區段不應該為 nil")
	}
	
	if toolbar.editSection == nil {
		t.Error("編輯區段不應該為 nil")
	}
	
	if toolbar.formatSection == nil {
		t.Error("格式化區段不應該為 nil")
	}
	
	if toolbar.insertSection == nil {
		t.Error("插入區段不應該為 nil")
	}
	
	if toolbar.viewSection == nil {
		t.Error("視圖區段不應該為 nil")
	}
	
	if toolbar.toolsSection == nil {
		t.Error("工具區段不應該為 nil")
	}
	
	// 驗證按鈕映射已初始化
	if toolbar.buttons == nil {
		t.Error("按鈕映射不應該為 nil")
	}
	
	// 驗證區段可見性設定已初始化
	if toolbar.sectionsVisible == nil {
		t.Error("區段可見性設定不應該為 nil")
	}
	
	// 驗證預設狀態
	if toolbar.compactMode {
		t.Error("預設不應該是緊湊模式")
	}
}

// TestEnhancedToolbarSectionVisibility 測試區段可見性控制
// 驗證區段的顯示/隱藏功能是否正常工作
func TestEnhancedToolbarSectionVisibility(t *testing.T) {
	toolbar := NewEnhancedToolbar()
	
	// 測試所有區段預設都可見
	sections := toolbar.GetAvailableSections()
	for _, section := range sections {
		if !toolbar.IsSectionVisible(section) {
			t.Errorf("區段 %s 預設應該可見", section)
		}
	}
	
	// 隱藏檔案區段
	toolbar.SetSectionVisible("file", false)
	
	// 驗證檔案區段已隱藏
	if toolbar.IsSectionVisible("file") {
		t.Error("檔案區段應該已隱藏")
	}
	
	// 顯示檔案區段
	toolbar.SetSectionVisible("file", true)
	
	// 驗證檔案區段已顯示
	if !toolbar.IsSectionVisible("file") {
		t.Error("檔案區段應該已顯示")
	}
}

// TestEnhancedToolbarCompactMode 測試緊湊模式
// 驗證緊湊模式的切換和按鈕大小調整
func TestEnhancedToolbarCompactMode(t *testing.T) {
	toolbar := NewEnhancedToolbar()
	
	// 驗證初始狀態不是緊湊模式
	if toolbar.IsCompactMode() {
		t.Error("初始狀態不應該是緊湊模式")
	}
	
	// 啟用緊湊模式
	toolbar.SetCompactMode(true)
	
	// 驗證緊湊模式已啟用
	if !toolbar.IsCompactMode() {
		t.Error("緊湊模式應該已啟用")
	}
	
	// 停用緊湊模式
	toolbar.SetCompactMode(false)
	
	// 驗證緊湊模式已停用
	if toolbar.IsCompactMode() {
		t.Error("緊湊模式應該已停用")
	}
}

// TestEnhancedToolbarButtonOperations 測試按鈕操作
// 驗證按鈕的啟用、停用和文字設定功能
func TestEnhancedToolbarButtonOperations(t *testing.T) {
	toolbar := NewEnhancedToolbar()
	
	// 測試按鈕存在性
	buttonId := "new_note"
	button := toolbar.GetButton(buttonId)
	if button == nil {
		t.Errorf("按鈕 %s 應該存在", buttonId)
	}
	
	// 測試停用按鈕
	toolbar.DisableButton(buttonId)
	if !button.Disabled() {
		t.Errorf("按鈕 %s 應該已停用", buttonId)
	}
	
	// 測試啟用按鈕
	toolbar.EnableButton(buttonId)
	if button.Disabled() {
		t.Errorf("按鈕 %s 應該已啟用", buttonId)
	}
	
	// 測試設定按鈕文字
	newText := "測試文字"
	toolbar.SetButtonText(buttonId, newText)
	if button.Text != newText {
		t.Errorf("按鈕文字應該為 '%s'，實際為 '%s'", newText, button.Text)
	}
}

// TestEnhancedToolbarGetSectionButtons 測試取得區段按鈕
// 驗證各區段的按鈕列表是否正確
func TestEnhancedToolbarGetSectionButtons(t *testing.T) {
	toolbar := NewEnhancedToolbar()
	
	// 測試檔案區段按鈕
	fileButtons := toolbar.GetSectionButtons("file")
	expectedFileButtons := []string{"new_note", "new_folder", "open_file", "save_file", "save_as", "import", "export"}
	
	if len(fileButtons) != len(expectedFileButtons) {
		t.Errorf("檔案區段按鈕數量應該為 %d，實際為 %d", len(expectedFileButtons), len(fileButtons))
	}
	
	for i, expected := range expectedFileButtons {
		if i < len(fileButtons) && fileButtons[i] != expected {
			t.Errorf("檔案區段按鈕 %d 應該為 '%s'，實際為 '%s'", i, expected, fileButtons[i])
		}
	}
	
	// 測試編輯區段按鈕
	editButtons := toolbar.GetSectionButtons("edit")
	expectedEditButtons := []string{"undo", "redo", "cut", "copy", "paste", "find", "replace"}
	
	if len(editButtons) != len(expectedEditButtons) {
		t.Errorf("編輯區段按鈕數量應該為 %d，實際為 %d", len(expectedEditButtons), len(editButtons))
	}
	
	// 測試格式化區段按鈕
	formatButtons := toolbar.GetSectionButtons("format")
	expectedFormatButtons := []string{"format_bold", "format_italic", "format_underline", "format_strikethrough", "heading_1", "heading_2", "heading_3"}
	
	if len(formatButtons) != len(expectedFormatButtons) {
		t.Errorf("格式化區段按鈕數量應該為 %d，實際為 %d", len(expectedFormatButtons), len(formatButtons))
	}
	
	// 測試插入區段按鈕
	insertButtons := toolbar.GetSectionButtons("insert")
	expectedInsertButtons := []string{"insert_link", "insert_image", "insert_table", "insert_code", "list_bullet", "list_numbered", "list_todo"}
	
	if len(insertButtons) != len(expectedInsertButtons) {
		t.Errorf("插入區段按鈕數量應該為 %d，實際為 %d", len(expectedInsertButtons), len(insertButtons))
	}
	
	// 測試視圖區段按鈕
	viewButtons := toolbar.GetSectionButtons("view")
	expectedViewButtons := []string{"toggle_preview", "edit_mode", "preview_mode", "split_view", "fullscreen", "zoom_in", "zoom_out", "toggle_theme"}
	
	if len(viewButtons) != len(expectedViewButtons) {
		t.Errorf("視圖區段按鈕數量應該為 %d，實際為 %d", len(expectedViewButtons), len(viewButtons))
	}
	
	// 測試工具區段按鈕
	toolsButtons := toolbar.GetSectionButtons("tools")
	expectedToolsButtons := []string{"toggle_encryption", "toggle_favorite", "manage_tags", "show_stats", "open_settings", "show_help"}
	
	if len(toolsButtons) != len(expectedToolsButtons) {
		t.Errorf("工具區段按鈕數量應該為 %d，實際為 %d", len(expectedToolsButtons), len(toolsButtons))
	}
	
	// 測試無效區段
	invalidButtons := toolbar.GetSectionButtons("invalid")
	if len(invalidButtons) != 0 {
		t.Error("無效區段應該回傳空的按鈕列表")
	}
}

// TestEnhancedToolbarActionCallback 測試動作回調
// 驗證動作觸發回調是否正確工作
func TestEnhancedToolbarActionCallback(t *testing.T) {
	toolbar := NewEnhancedToolbar()
	
	// 設定回調函數
	var actionTriggered bool
	var lastAction string
	var lastParams map[string]interface{}
	
	toolbar.SetOnActionTriggered(func(action string, params map[string]interface{}) {
		actionTriggered = true
		lastAction = action
		lastParams = params
	})
	
	// 觸發動作
	testAction := "test_action"
	testParams := map[string]interface{}{"key": "value"}
	toolbar.triggerAction(testAction, testParams)
	
	// 驗證回調已觸發
	if !actionTriggered {
		t.Error("動作回調應該已觸發")
	}
	
	if lastAction != testAction {
		t.Errorf("動作名稱應該為 '%s'，實際為 '%s'", testAction, lastAction)
	}
	
	if lastParams == nil {
		t.Error("動作參數不應該為 nil")
	}
	
	if lastParams["key"] != "value" {
		t.Error("動作參數應該包含正確的值")
	}
}

// TestEnhancedToolbarGetAvailableSections 測試取得可用區段
// 驗證可用區段列表是否正確
func TestEnhancedToolbarGetAvailableSections(t *testing.T) {
	toolbar := NewEnhancedToolbar()
	
	sections := toolbar.GetAvailableSections()
	expectedSections := []string{"file", "edit", "format", "insert", "view", "tools"}
	
	if len(sections) != len(expectedSections) {
		t.Errorf("可用區段數量應該為 %d，實際為 %d", len(expectedSections), len(sections))
	}
	
	for i, expected := range expectedSections {
		if i < len(sections) && sections[i] != expected {
			t.Errorf("區段 %d 應該為 '%s'，實際為 '%s'", i, expected, sections[i])
		}
	}
}

// TestEnhancedToolbarButtonCreation 測試按鈕建立
// 驗證所有預期的按鈕是否都已建立
func TestEnhancedToolbarButtonCreation(t *testing.T) {
	toolbar := NewEnhancedToolbar()
	
	// 取得所有區段的所有按鈕
	allExpectedButtons := []string{}
	sections := toolbar.GetAvailableSections()
	
	for _, section := range sections {
		sectionButtons := toolbar.GetSectionButtons(section)
		allExpectedButtons = append(allExpectedButtons, sectionButtons...)
	}
	
	// 驗證所有按鈕都已建立
	for _, buttonId := range allExpectedButtons {
		button := toolbar.GetButton(buttonId)
		if button == nil {
			t.Errorf("按鈕 '%s' 應該已建立", buttonId)
		}
	}
	
	// 驗證按鈕映射的大小
	if len(toolbar.buttons) != len(allExpectedButtons) {
		t.Errorf("按鈕映射大小應該為 %d，實際為 %d", len(allExpectedButtons), len(toolbar.buttons))
	}
}

// TestEnhancedToolbarContainer 測試工具欄容器
// 驗證工具欄容器是否正確組合
func TestEnhancedToolbarContainer(t *testing.T) {
	toolbar := NewEnhancedToolbar()
	
	// 取得容器
	container := toolbar.GetContainer()
	
	// 驗證容器不為 nil
	if container == nil {
		t.Error("工具欄容器不應該為 nil")
	}
	
	// 驗證容器包含內容
	if len(container.Objects) == 0 {
		t.Error("工具欄容器應該包含內容")
	}
}

// TestEnhancedToolbarSectionToggle 測試區段切換對容器的影響
// 驗證區段可見性變更是否正確更新容器內容
func TestEnhancedToolbarSectionToggle(t *testing.T) {
	toolbar := NewEnhancedToolbar()
	
	// 記錄初始容器物件數量
	initialObjectCount := len(toolbar.GetContainer().Objects)
	
	// 隱藏一個區段
	toolbar.SetSectionVisible("file", false)
	
	// 驗證容器物件數量減少
	newObjectCount := len(toolbar.GetContainer().Objects)
	if newObjectCount >= initialObjectCount {
		t.Error("隱藏區段後容器物件數量應該減少")
	}
	
	// 重新顯示區段
	toolbar.SetSectionVisible("file", true)
	
	// 驗證容器物件數量恢復
	finalObjectCount := len(toolbar.GetContainer().Objects)
	if finalObjectCount != initialObjectCount {
		t.Error("重新顯示區段後容器物件數量應該恢復")
	}
}