// Package ui 包含使用者介面相關的元件和視窗管理
// 本檔案實作整合編輯器和預覽面板的複合 UI 元件
package ui

import (
	"mac-notebook-app/internal/models"  // 引入資料模型
	"mac-notebook-app/internal/services" // 引入服務層

	"fyne.io/fyne/v2"               // Fyne GUI 框架核心套件
	"fyne.io/fyne/v2/container"     // Fyne 容器佈局套件
)

// EditorWithPreview 代表整合編輯器和預覽面板的複合 UI 元件
// 提供完整的 Markdown 編輯和即時預覽體驗
// 支援同步滾動、預覽切換和統一的狀態管理
type EditorWithPreview struct {
	container     *fyne.Container      // 主要容器
	splitContainer *container.Split    // 分割容器
	
	// 子元件
	editor        *MarkdownEditor      // Markdown 編輯器
	preview       *MarkdownPreview     // Markdown 預覽面板
	
	// 服務依賴
	editorService services.EditorService // 編輯器服務
	
	// 當前狀態
	previewVisible bool                 // 預覽面板是否可見
	splitRatio     float64              // 分割比例
	
	// 回調函數
	onContentChanged   func(content string) // 內容變更回調
	onSaveRequested    func()               // 保存請求回調
	onWordCountChanged func(count int)      // 字數變更回調
	onPreviewToggled   func(visible bool)   // 預覽切換回調
}

// NewEditorWithPreview 建立新的整合編輯器和預覽面板實例
// 參數：editorService（編輯器服務介面）
// 回傳：指向新建立的 EditorWithPreview 的指標
//
// 執行流程：
// 1. 建立 EditorWithPreview 結構體實例
// 2. 初始化編輯器服務依賴
// 3. 建立編輯器和預覽面板子元件
// 4. 設定元件間的事件連接
// 5. 組合完整的複合元件佈局
// 6. 回傳配置完成的複合元件實例
func NewEditorWithPreview(editorService services.EditorService) *EditorWithPreview {
	// 建立 EditorWithPreview 實例
	ewp := &EditorWithPreview{
		editorService:  editorService,
		previewVisible: true,  // 預設顯示預覽
		splitRatio:     0.5,   // 預設 50/50 分割
	}
	
	// 建立子元件
	ewp.createComponents()
	
	// 設定元件間的連接
	ewp.setupConnections()
	
	// 組合佈局
	ewp.setupLayout()
	
	return ewp
}

// createComponents 建立編輯器和預覽面板子元件
// 初始化所有必要的子元件實例
//
// 執行流程：
// 1. 建立 Markdown 編輯器實例
// 2. 建立 Markdown 預覽面板實例
// 3. 設定子元件的初始狀態
func (ewp *EditorWithPreview) createComponents() {
	// 建立 Markdown 編輯器
	ewp.editor = NewMarkdownEditor(ewp.editorService)
	
	// 建立 Markdown 預覽面板
	ewp.preview = NewMarkdownPreview(ewp.editorService)
}

// setupConnections 設定元件間的事件連接
// 建立編輯器和預覽面板之間的事件通信
//
// 執行流程：
// 1. 設定編輯器的內容變更回調，自動更新預覽
// 2. 設定編輯器的保存請求回調，傳遞到外部
// 3. 設定編輯器的字數變更回調，傳遞到外部
// 4. 設定預覽面板的可見性變更回調
// 5. 設定預覽面板的刷新請求回調
func (ewp *EditorWithPreview) setupConnections() {
	// 設定編輯器內容變更回調 - 自動更新預覽
	ewp.editor.SetOnContentChanged(func(content string) {
		// 如果預覽面板可見且啟用自動刷新，更新預覽
		if ewp.previewVisible && ewp.preview.IsAutoRefreshEnabled() {
			ewp.preview.UpdatePreview(content)
		}
		
		// 觸發外部內容變更回調
		if ewp.onContentChanged != nil {
			ewp.onContentChanged(content)
		}
	})
	
	// 設定編輯器保存請求回調
	ewp.editor.SetOnSaveRequested(func() {
		if ewp.onSaveRequested != nil {
			ewp.onSaveRequested()
		}
	})
	
	// 設定編輯器字數變更回調
	ewp.editor.SetOnWordCountChanged(func(count int) {
		if ewp.onWordCountChanged != nil {
			ewp.onWordCountChanged(count)
		}
	})
	
	// 設定預覽面板可見性變更回調
	ewp.preview.SetOnVisibilityChanged(func(visible bool) {
		ewp.previewVisible = visible
		ewp.updateLayout()
		
		if ewp.onPreviewToggled != nil {
			ewp.onPreviewToggled(visible)
		}
	})
	
	// 設定預覽面板刷新請求回調
	ewp.preview.SetOnRefreshRequested(func() {
		// 從編輯器取得當前內容並更新預覽
		content := ewp.editor.GetContent()
		ewp.preview.UpdatePreview(content)
	})
}

// setupLayout 設定複合元件的佈局
// 組合編輯器和預覽面板到分割容器中
//
// 執行流程：
// 1. 建立水平分割容器
// 2. 設定分割比例
// 3. 添加編輯器和預覽面板到分割容器
// 4. 建立主要容器包含分割容器
func (ewp *EditorWithPreview) setupLayout() {
	// 建立水平分割容器
	ewp.splitContainer = container.NewHSplit(
		ewp.editor.GetContainer(),
		ewp.preview.GetContainer(),
	)
	
	// 設定分割比例
	ewp.splitContainer.Offset = ewp.splitRatio
	
	// 建立主要容器
	ewp.container = container.NewVBox(ewp.splitContainer)
	
	// 根據預覽可見性更新佈局
	ewp.updateLayout()
}

// updateLayout 更新佈局以反映當前狀態
// 根據預覽面板的可見性調整佈局
//
// 執行流程：
// 1. 檢查預覽面板可見性
// 2. 如果可見，使用分割佈局
// 3. 如果隱藏，只顯示編輯器
// 4. 更新主容器內容
func (ewp *EditorWithPreview) updateLayout() {
	if ewp.previewVisible {
		// 顯示分割佈局（編輯器 + 預覽）
		ewp.container.Objects = []fyne.CanvasObject{ewp.splitContainer}
	} else {
		// 只顯示編輯器
		ewp.container.Objects = []fyne.CanvasObject{ewp.editor.GetContainer()}
	}
	
	ewp.container.Refresh()
}

// GetContainer 取得複合元件的主要容器
// 回傳：複合元件的 fyne.Container 實例
// 用於將複合元件嵌入到其他 UI 佈局中
func (ewp *EditorWithPreview) GetContainer() *fyne.Container {
	return ewp.container
}

// LoadNote 載入筆記到編輯器和預覽面板
// 參數：note（要載入的筆記實例）
//
// 執行流程：
// 1. 載入筆記到編輯器
// 2. 如果預覽可見，更新預覽內容
// 3. 同步滾動位置（如果啟用）
func (ewp *EditorWithPreview) LoadNote(note *models.Note) {
	// 載入筆記到編輯器
	ewp.editor.LoadNote(note)
	
	// 如果預覽可見，更新預覽
	if ewp.previewVisible {
		ewp.preview.UpdatePreview(note.Content)
	}
}

// CreateNewNote 建立新筆記
// 參數：title（筆記標題）
//
// 執行流程：
// 1. 在編輯器中建立新筆記
// 2. 清空預覽面板
// 3. 設定編輯器焦點
func (ewp *EditorWithPreview) CreateNewNote(title string) error {
	// 在編輯器中建立新筆記
	err := ewp.editor.CreateNewNote(title)
	if err != nil {
		return err
	}
	
	// 清空預覽面板
	ewp.preview.Clear()
	
	// 設定編輯器焦點
	ewp.editor.Focus()
	
	return nil
}

// SaveNote 保存當前筆記
// 回傳：可能的錯誤
//
// 執行流程：
// 1. 保存編輯器中的筆記
// 2. 更新預覽內容（如果需要）
func (ewp *EditorWithPreview) SaveNote() error {
	return ewp.editor.SaveNote()
}

// GetContent 取得編輯器當前內容
// 回傳：編輯器中的文字內容
func (ewp *EditorWithPreview) GetContent() string {
	return ewp.editor.GetContent()
}

// SetContent 設定編輯器內容
// 參數：content（要設定的內容）
//
// 執行流程：
// 1. 設定編輯器內容
// 2. 如果預覽可見，更新預覽
func (ewp *EditorWithPreview) SetContent(content string) {
	ewp.editor.SetContent(content)
	
	if ewp.previewVisible {
		ewp.preview.UpdatePreview(content)
	}
}

// IsModified 檢查內容是否已修改
// 回傳：內容是否已修改的布林值
func (ewp *EditorWithPreview) IsModified() bool {
	return ewp.editor.IsModified()
}

// GetCurrentNote 取得當前編輯的筆記
// 回傳：當前筆記實例
func (ewp *EditorWithPreview) GetCurrentNote() *models.Note {
	return ewp.editor.GetCurrentNote()
}

// CanSave 檢查是否可以保存
// 回傳：是否可以保存的布林值
func (ewp *EditorWithPreview) CanSave() bool {
	return ewp.editor.CanSave()
}

// Clear 清空編輯器和預覽內容
// 清除所有內容並重置狀態
//
// 執行流程：
// 1. 清空編輯器內容
// 2. 清空預覽面板內容
func (ewp *EditorWithPreview) Clear() {
	ewp.editor.Clear()
	ewp.preview.Clear()
}

// TogglePreview 切換預覽面板的顯示/隱藏
// 在顯示和隱藏預覽面板之間切換
//
// 執行流程：
// 1. 切換預覽可見性狀態
// 2. 更新預覽面板可見性
// 3. 更新佈局
// 4. 如果顯示預覽，更新預覽內容
func (ewp *EditorWithPreview) TogglePreview() {
	ewp.previewVisible = !ewp.previewVisible
	ewp.preview.SetVisible(ewp.previewVisible)
	
	// 如果顯示預覽，更新預覽內容
	if ewp.previewVisible {
		content := ewp.editor.GetContent()
		ewp.preview.UpdatePreview(content)
	}
}

// IsPreviewVisible 檢查預覽面板是否可見
// 回傳：預覽面板是否可見的布林值
func (ewp *EditorWithPreview) IsPreviewVisible() bool {
	return ewp.previewVisible
}

// SetPreviewVisible 設定預覽面板可見性
// 參數：visible（是否可見）
//
// 執行流程：
// 1. 設定預覽可見性狀態
// 2. 更新預覽面板可見性
// 3. 如果顯示預覽，更新預覽內容
func (ewp *EditorWithPreview) SetPreviewVisible(visible bool) {
	ewp.previewVisible = visible
	ewp.preview.SetVisible(visible)
	
	if visible {
		content := ewp.editor.GetContent()
		ewp.preview.UpdatePreview(content)
	}
}

// SetSplitRatio 設定分割比例
// 參數：ratio（分割比例，0.0-1.0）
//
// 執行流程：
// 1. 更新分割比例
// 2. 應用到分割容器
func (ewp *EditorWithPreview) SetSplitRatio(ratio float64) {
	if ratio < 0.0 {
		ratio = 0.0
	} else if ratio > 1.0 {
		ratio = 1.0
	}
	
	ewp.splitRatio = ratio
	if ewp.splitContainer != nil {
		ewp.splitContainer.Offset = ratio
	}
}

// GetSplitRatio 取得當前分割比例
// 回傳：當前分割比例（0.0-1.0）
func (ewp *EditorWithPreview) GetSplitRatio() float64 {
	return ewp.splitRatio
}

// RefreshPreview 手動刷新預覽
// 強制更新預覽面板內容
//
// 執行流程：
// 1. 取得編輯器當前內容
// 2. 更新預覽面板
func (ewp *EditorWithPreview) RefreshPreview() {
	content := ewp.editor.GetContent()
	ewp.preview.UpdatePreview(content)
}

// SetAutoRefresh 設定預覽自動刷新
// 參數：enabled（是否啟用自動刷新）
func (ewp *EditorWithPreview) SetAutoRefresh(enabled bool) {
	ewp.preview.SetAutoRefresh(enabled)
}

// IsAutoRefreshEnabled 檢查是否啟用自動刷新
// 回傳：自動刷新是否啟用的布林值
func (ewp *EditorWithPreview) IsAutoRefreshEnabled() bool {
	return ewp.preview.IsAutoRefreshEnabled()
}

// GetTitle 取得當前筆記標題
// 回傳：筆記標題，如果沒有當前筆記則回傳空字串
func (ewp *EditorWithPreview) GetTitle() string {
	return ewp.editor.GetTitle()
}

// SetTitle 設定當前筆記標題
// 參數：title（新的標題）
func (ewp *EditorWithPreview) SetTitle(title string) {
	ewp.editor.SetTitle(title)
}

// Focus 設定編輯器焦點
// 讓編輯器獲得輸入焦點
func (ewp *EditorWithPreview) Focus() {
	ewp.editor.Focus()
}

// GetEditor 取得編輯器元件
// 回傳：Markdown 編輯器實例
// 用於直接存取編輯器的特定功能
func (ewp *EditorWithPreview) GetEditor() *MarkdownEditor {
	return ewp.editor
}

// GetPreview 取得預覽面板元件
// 回傳：Markdown 預覽面板實例
// 用於直接存取預覽面板的特定功能
func (ewp *EditorWithPreview) GetPreview() *MarkdownPreview {
	return ewp.preview
}

// SetOnContentChanged 設定內容變更回調函數
// 參數：callback（內容變更時的回調函數）
func (ewp *EditorWithPreview) SetOnContentChanged(callback func(content string)) {
	ewp.onContentChanged = callback
}

// SetOnSaveRequested 設定保存請求回調函數
// 參數：callback（保存請求時的回調函數）
func (ewp *EditorWithPreview) SetOnSaveRequested(callback func()) {
	ewp.onSaveRequested = callback
}

// SetOnWordCountChanged 設定字數變更回調函數
// 參數：callback（字數變更時的回調函數）
func (ewp *EditorWithPreview) SetOnWordCountChanged(callback func(count int)) {
	ewp.onWordCountChanged = callback
}

// SetOnPreviewToggled 設定預覽切換回調函數
// 參數：callback（預覽切換時的回調函數）
func (ewp *EditorWithPreview) SetOnPreviewToggled(callback func(visible bool)) {
	ewp.onPreviewToggled = callback
}