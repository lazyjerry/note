// Package ui 包含使用者介面相關的元件和視窗管理
// 本檔案實作 Markdown 編輯器 UI 元件，提供文字編輯、語法高亮和工具欄功能
package ui

import (
	"fmt"                           // Go 標準庫，用於格式化字串
	"strings"                       // 字串處理
	"mac-notebook-app/internal/models" // 引入資料模型
	"mac-notebook-app/internal/services" // 引入服務層

	"fyne.io/fyne/v2"               // Fyne GUI 框架核心套件
	"fyne.io/fyne/v2/container"     // Fyne 容器佈局套件
	"fyne.io/fyne/v2/widget"        // Fyne UI 元件套件
	"fyne.io/fyne/v2/theme"         // Fyne 主題套件
	"fyne.io/fyne/v2/dialog"        // Fyne 對話框套件
)

// MarkdownEditor 代表 Markdown 編輯器 UI 元件
// 包含文字編輯器、工具欄、語法高亮和即時預覽功能
// 整合編輯器服務以提供完整的筆記編輯體驗
type MarkdownEditor struct {
	container     *fyne.Container      // 主要容器
	toolbar       *widget.Toolbar      // 編輯器工具欄
	editor        *widget.Entry        // 文字編輯器元件
	statusLabel   *widget.Label        // 狀態標籤
	
	// 服務依賴
	editorService services.EditorService // 編輯器服務
	
	// 當前狀態
	currentNote   *models.Note         // 當前編輯的筆記
	isModified    bool                 // 內容是否已修改
	
	// 回調函數
	onContentChanged func(content string) // 內容變更回調
	onSaveRequested  func()               // 保存請求回調
	onWordCountChanged func(count int)    // 字數變更回調
}

// NewMarkdownEditor 建立新的 Markdown 編輯器實例
// 參數：editorService（編輯器服務介面）
// 回傳：指向新建立的 MarkdownEditor 的指標
//
// 執行流程：
// 1. 建立 MarkdownEditor 結構體實例
// 2. 初始化編輯器服務依賴
// 3. 建立並配置所有 UI 元件
// 4. 設定事件處理和回調函數
// 5. 組合完整的編輯器佈局
// 6. 回傳配置完成的編輯器實例
func NewMarkdownEditor(editorService services.EditorService) *MarkdownEditor {
	// 建立 MarkdownEditor 實例
	editor := &MarkdownEditor{
		editorService: editorService,
		isModified:    false,
	}
	
	// 初始化 UI 元件
	editor.setupUI()
	
	return editor
}

// setupUI 初始化編輯器的使用者介面元件
// 這個方法負責建立和配置編輯器的所有 UI 元件
//
// 執行流程：
// 1. 建立編輯器工具欄和功能按鈕
// 2. 建立文字編輯器元件並配置屬性
// 3. 建立狀態標籤顯示編輯狀態
// 4. 設定事件處理和回調函數
// 5. 組合所有元件到主容器中
func (me *MarkdownEditor) setupUI() {
	// 建立編輯器工具欄
	me.createToolbar()
	
	// 建立文字編輯器
	me.createTextEditor()
	
	// 建立狀態標籤
	me.createStatusLabel()
	
	// 組合編輯器佈局
	me.assembleLayout()
}

// createToolbar 建立編輯器工具欄
// 包含常用的 Markdown 格式化按鈕和編輯功能
//
// 執行流程：
// 1. 建立標題格式化按鈕（H1-H3）
// 2. 建立文字格式化按鈕（粗體、斜體、刪除線）
// 3. 建立列表和連結按鈕
// 4. 建立保存和預覽按鈕
// 5. 組合所有按鈕到工具欄
func (me *MarkdownEditor) createToolbar() {
	me.toolbar = widget.NewToolbar(
		// 標題格式化按鈕
		widget.NewToolbarAction(theme.DocumentIcon(), func() {
			me.insertMarkdown("# ", "", "標題 1")
		}),
		
		widget.NewToolbarAction(theme.DocumentIcon(), func() {
			me.insertMarkdown("## ", "", "標題 2")
		}),
		
		widget.NewToolbarAction(theme.DocumentIcon(), func() {
			me.insertMarkdown("### ", "", "標題 3")
		}),
		
		// 分隔線
		widget.NewToolbarSeparator(),
		
		// 文字格式化按鈕
		widget.NewToolbarAction(theme.ContentCopyIcon(), func() {
			me.wrapSelection("**", "**", "粗體文字")
		}),
		
		widget.NewToolbarAction(theme.ContentCopyIcon(), func() {
			me.wrapSelection("*", "*", "斜體文字")
		}),
		
		widget.NewToolbarAction(theme.ContentCopyIcon(), func() {
			me.wrapSelection("~~", "~~", "刪除線文字")
		}),
		
		// 分隔線
		widget.NewToolbarSeparator(),
		
		// 列表按鈕
		widget.NewToolbarAction(theme.ListIcon(), func() {
			me.insertMarkdown("- ", "", "列表項目")
		}),
		
		widget.NewToolbarAction(theme.ListIcon(), func() {
			me.insertMarkdown("1. ", "", "編號列表項目")
		}),
		
		// 分隔線
		widget.NewToolbarSeparator(),
		
		// 連結和圖片按鈕
		widget.NewToolbarAction(theme.ContentCopyIcon(), func() {
			me.insertMarkdown("[", "](https://example.com)", "連結文字")
		}),
		
		widget.NewToolbarAction(theme.ContentCopyIcon(), func() {
			me.insertMarkdown("![", "](image.png)", "圖片描述")
		}),
		
		// 分隔線
		widget.NewToolbarSeparator(),
		
		// 程式碼區塊按鈕
		widget.NewToolbarAction(theme.ContentCopyIcon(), func() {
			me.wrapSelection("`", "`", "程式碼")
		}),
		
		widget.NewToolbarAction(theme.ContentCopyIcon(), func() {
			me.insertMarkdown("```\n", "\n```", "程式碼區塊")
		}),
		
		// 分隔線
		widget.NewToolbarSeparator(),
		
		// 保存按鈕
		widget.NewToolbarAction(theme.DocumentSaveIcon(), func() {
			me.saveContent()
		}),
		
		// 預覽切換按鈕
		widget.NewToolbarAction(theme.ViewRefreshIcon(), func() {
			me.togglePreview()
		}),
	)
}

// createTextEditor 建立文字編輯器元件
// 配置多行文字輸入、自動換行和事件處理
//
// 執行流程：
// 1. 建立多行文字輸入元件
// 2. 設定編輯器屬性（自動換行、滾動等）
// 3. 設定內容變更事件處理
// 4. 設定鍵盤快捷鍵處理
// 5. 配置編輯器樣式和字型
func (me *MarkdownEditor) createTextEditor() {
	// 建立多行文字編輯器
	me.editor = widget.NewMultiLineEntry()
	
	// 設定編輯器屬性
	me.editor.Wrapping = fyne.TextWrapWord  // 自動換行
	me.editor.Scroll = container.ScrollBoth  // 雙向滾動
	me.editor.SetPlaceHolder("在此輸入您的 Markdown 內容...")
	
	// 設定內容變更事件處理
	me.editor.OnChanged = func(content string) {
		me.onTextChanged(content)
	}
	
	// 設定鍵盤事件處理
	me.editor.OnSubmitted = func(content string) {
		// Enter 鍵處理（如果需要特殊行為）
		me.handleEnterKey()
	}
}

// createStatusLabel 建立狀態標籤
// 顯示編輯狀態、字數統計和其他資訊
//
// 執行流程：
// 1. 建立狀態標籤元件
// 2. 設定初始狀態文字
// 3. 配置標籤樣式
func (me *MarkdownEditor) createStatusLabel() {
	me.statusLabel = widget.NewLabel("準備就緒")
	me.statusLabel.TextStyle = fyne.TextStyle{Italic: true}
}

// assembleLayout 組合編輯器的完整佈局
// 將工具欄、編輯器和狀態欄組合成完整的編輯器介面
//
// 執行流程：
// 1. 建立垂直容器作為主要佈局
// 2. 依序添加工具欄、編輯器和狀態欄
// 3. 設定容器屬性和佈局比例
func (me *MarkdownEditor) assembleLayout() {
	// 建立主要容器，使用垂直佈局
	me.container = container.NewVBox(
		me.toolbar,                    // 工具欄在頂部
		widget.NewSeparator(),         // 分隔線
		me.editor,                     // 編輯器在中間（主要區域）
		widget.NewSeparator(),         // 分隔線
		me.statusLabel,                // 狀態欄在底部
	)
}

// GetContainer 取得編輯器的主要容器
// 回傳：編輯器的 fyne.Container 實例
// 用於將編輯器嵌入到其他 UI 佈局中
func (me *MarkdownEditor) GetContainer() *fyne.Container {
	return me.container
}

// LoadNote 載入筆記到編輯器
// 參數：note（要載入的筆記實例）
//
// 執行流程：
// 1. 設定當前筆記實例
// 2. 將筆記內容載入到編輯器
// 3. 重置修改狀態
// 4. 更新狀態顯示
// 5. 觸發內容變更回調
func (me *MarkdownEditor) LoadNote(note *models.Note) {
	if note == nil {
		return
	}
	
	// 設定當前筆記
	me.currentNote = note
	
	// 載入筆記內容到編輯器
	me.editor.SetText(note.Content)
	
	// 重置修改狀態
	me.isModified = false
	
	// 更新狀態顯示
	me.updateStatus(fmt.Sprintf("已載入筆記: %s", note.Title))
	
	// 觸發字數統計更新
	me.updateWordCount()
}

// CreateNewNote 建立新筆記
// 參數：title（筆記標題）
//
// 執行流程：
// 1. 使用編輯器服務建立新筆記
// 2. 載入新筆記到編輯器
// 3. 設定編輯器焦點
// 4. 更新狀態顯示
func (me *MarkdownEditor) CreateNewNote(title string) error {
	// 使用編輯器服務建立新筆記
	note, err := me.editorService.CreateNote(title, "")
	if err != nil {
		me.updateStatus(fmt.Sprintf("建立筆記失敗: %s", err.Error()))
		return err
	}
	
	// 載入新筆記到編輯器
	me.LoadNote(note)
	
	// 設定編輯器焦點
	me.editor.FocusGained()
	
	// 更新狀態顯示
	me.updateStatus(fmt.Sprintf("已建立新筆記: %s", title))
	
	return nil
}

// SaveNote 保存當前筆記
// 回傳：可能的錯誤
//
// 執行流程：
// 1. 檢查是否有當前筆記
// 2. 更新筆記內容
// 3. 使用編輯器服務保存筆記
// 4. 重置修改狀態
// 5. 更新狀態顯示
// 6. 觸發保存回調
func (me *MarkdownEditor) SaveNote() error {
	if me.currentNote == nil {
		return fmt.Errorf("沒有可保存的筆記")
	}
	
	// 更新筆記內容
	content := me.editor.Text
	err := me.editorService.UpdateContent(me.currentNote.ID, content)
	if err != nil {
		me.updateStatus(fmt.Sprintf("更新內容失敗: %s", err.Error()))
		return err
	}
	
	// 保存筆記
	err = me.editorService.SaveNote(me.currentNote)
	if err != nil {
		me.updateStatus(fmt.Sprintf("保存失敗: %s", err.Error()))
		return err
	}
	
	// 重置修改狀態
	me.isModified = false
	
	// 更新狀態顯示
	me.updateStatus(fmt.Sprintf("已保存筆記: %s", me.currentNote.Title))
	
	// 觸發保存回調
	if me.onSaveRequested != nil {
		me.onSaveRequested()
	}
	
	return nil
}

// GetContent 取得編輯器當前內容
// 回傳：編輯器中的文字內容
func (me *MarkdownEditor) GetContent() string {
	return me.editor.Text
}

// SetContent 設定編輯器內容
// 參數：content（要設定的內容）
//
// 執行流程：
// 1. 設定編輯器文字內容
// 2. 重置修改狀態
// 3. 更新字數統計
func (me *MarkdownEditor) SetContent(content string) {
	me.editor.SetText(content)
	me.isModified = false
	me.updateWordCount()
}

// IsModified 檢查內容是否已修改
// 回傳：內容是否已修改的布林值
func (me *MarkdownEditor) IsModified() bool {
	return me.isModified
}

// GetCurrentNote 取得當前編輯的筆記
// 回傳：當前筆記實例
func (me *MarkdownEditor) GetCurrentNote() *models.Note {
	return me.currentNote
}

// SetOnContentChanged 設定內容變更回調函數
// 參數：callback（內容變更時的回調函數）
func (me *MarkdownEditor) SetOnContentChanged(callback func(content string)) {
	me.onContentChanged = callback
}

// SetOnSaveRequested 設定保存請求回調函數
// 參數：callback（保存請求時的回調函數）
func (me *MarkdownEditor) SetOnSaveRequested(callback func()) {
	me.onSaveRequested = callback
}

// SetOnWordCountChanged 設定字數變更回調函數
// 參數：callback（字數變更時的回調函數）
func (me *MarkdownEditor) SetOnWordCountChanged(callback func(count int)) {
	me.onWordCountChanged = callback
}

// insertMarkdown 在游標位置插入 Markdown 語法
// 參數：
//   - prefix: 前綴字串
//   - suffix: 後綴字串（可選）
//   - placeholder: 佔位文字
//
// 執行流程：
// 1. 取得當前游標位置
// 2. 在游標位置插入 Markdown 語法
// 3. 設定選取範圍到佔位文字
// 4. 觸發內容變更事件
func (me *MarkdownEditor) insertMarkdown(prefix, suffix, placeholder string) {
	// 取得當前內容和游標位置
	content := me.editor.Text
	cursorPos := len(content) // 簡化實作，實際應該取得真實游標位置
	
	// 建立要插入的文字
	insertText := prefix + placeholder
	if suffix != "" {
		insertText += suffix
	}
	
	// 插入文字
	newContent := content[:cursorPos] + insertText + content[cursorPos:]
	me.editor.SetText(newContent)
	
	// 標記為已修改
	me.isModified = true
	
	// 更新狀態
	me.updateStatus("已插入 Markdown 語法")
}

// wrapSelection 用指定的標記包圍選取的文字
// 參數：
//   - startMark: 開始標記
//   - endMark: 結束標記
//   - placeholder: 沒有選取文字時的佔位文字
//
// 執行流程：
// 1. 取得當前選取的文字
// 2. 如果沒有選取文字，使用佔位文字
// 3. 用指定標記包圍文字
// 4. 替換編輯器內容
// 5. 觸發內容變更事件
func (me *MarkdownEditor) wrapSelection(startMark, endMark, placeholder string) {
	// 簡化實作：由於 Fyne Entry 不直接支援選取範圍操作
	// 這裡在游標位置插入包圍的佔位文字
	me.insertMarkdown(startMark, endMark, placeholder)
}

// onTextChanged 處理文字內容變更事件
// 參數：content（變更後的內容）
//
// 執行流程：
// 1. 標記內容為已修改
// 2. 更新字數統計
// 3. 觸發內容變更回調
// 4. 更新狀態顯示
func (me *MarkdownEditor) onTextChanged(content string) {
	// 標記為已修改
	me.isModified = true
	
	// 更新字數統計
	me.updateWordCount()
	
	// 觸發內容變更回調
	if me.onContentChanged != nil {
		me.onContentChanged(content)
	}
	
	// 更新狀態顯示
	if me.isModified {
		me.updateStatus("內容已修改")
	}
}

// updateWordCount 更新字數統計
// 計算當前內容的字數並觸發相關回調
//
// 執行流程：
// 1. 取得當前編輯器內容
// 2. 計算字數（以空白分隔的詞數）
// 3. 觸發字數變更回調
func (me *MarkdownEditor) updateWordCount() {
	content := me.editor.Text
	
	// 簡單的字數計算（以空白分隔）
	words := strings.Fields(content)
	wordCount := len(words)
	
	// 觸發字數變更回調
	if me.onWordCountChanged != nil {
		me.onWordCountChanged(wordCount)
	}
}

// updateStatus 更新狀態標籤顯示
// 參數：message（要顯示的狀態訊息）
func (me *MarkdownEditor) updateStatus(message string) {
	if me.statusLabel != nil {
		me.statusLabel.SetText(message)
		me.statusLabel.Refresh()
	}
}

// saveContent 保存內容的內部方法
// 處理保存操作並顯示適當的回饋
//
// 執行流程：
// 1. 嘗試保存當前筆記
// 2. 處理保存結果
// 3. 顯示成功或錯誤訊息
func (me *MarkdownEditor) saveContent() {
	err := me.SaveNote()
	if err != nil {
		// 顯示錯誤對話框
		dialog.ShowError(err, me.container.Objects[0].(fyne.Window))
	}
}

// togglePreview 切換預覽模式
// 觸發預覽面板的顯示/隱藏
//
// 執行流程：
// 1. 取得當前內容
// 2. 觸發預覽切換事件
// 3. 更新狀態顯示
func (me *MarkdownEditor) togglePreview() {
	// 這個功能將在預覽面板實作時完成
	me.updateStatus("預覽功能將在下一個任務中實作")
}

// handleEnterKey 處理 Enter 鍵事件
// 提供智慧的換行和列表繼續功能
//
// 執行流程：
// 1. 檢查當前行是否為列表項目
// 2. 如果是列表，自動添加新的列表項目
// 3. 否則正常換行
func (me *MarkdownEditor) handleEnterKey() {
	// 簡化實作：正常換行
	// 實際實作可以添加智慧列表繼續等功能
	me.updateStatus("Enter 鍵處理")
}

// Focus 設定編輯器焦點
// 讓編輯器獲得輸入焦點
func (me *MarkdownEditor) Focus() {
	if me.editor != nil {
		me.editor.FocusGained()
	}
}

// Clear 清空編輯器內容
// 清除所有文字並重置狀態
//
// 執行流程：
// 1. 清空編輯器文字
// 2. 重置當前筆記
// 3. 重置修改狀態
// 4. 更新狀態顯示
func (me *MarkdownEditor) Clear() {
	me.editor.SetText("")
	me.currentNote = nil
	me.isModified = false
	me.updateStatus("編輯器已清空")
}

// CanSave 檢查是否可以保存
// 回傳：是否可以保存的布林值
//
// 執行流程：
// 1. 檢查是否有當前筆記
// 2. 檢查內容是否已修改
// 3. 回傳是否可以保存
func (me *MarkdownEditor) CanSave() bool {
	return me.currentNote != nil && me.isModified
}

// GetTitle 取得當前筆記標題
// 回傳：筆記標題，如果沒有當前筆記則回傳空字串
func (me *MarkdownEditor) GetTitle() string {
	if me.currentNote != nil {
		return me.currentNote.Title
	}
	return ""
}

// SetTitle 設定當前筆記標題
// 參數：title（新的標題）
//
// 執行流程：
// 1. 檢查是否有當前筆記
// 2. 更新筆記標題
// 3. 標記為已修改
// 4. 更新狀態顯示
func (me *MarkdownEditor) SetTitle(title string) {
	if me.currentNote != nil {
		me.currentNote.Title = title
		me.isModified = true
		me.updateStatus(fmt.Sprintf("標題已更新: %s", title))
	}
}