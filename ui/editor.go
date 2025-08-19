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
)

// MarkdownEditor 代表 Markdown 編輯器 UI 元件
// 包含文字編輯器、工具欄、語法高亮和即時預覽功能
// 整合編輯器服務以提供完整的筆記編輯體驗
// 支援繁體中文輸入法優化和中文字符處理
type MarkdownEditor struct {
	container     *fyne.Container      // 主要容器
	toolbar       *fyne.Container      // 編輯器工具欄容器（包含兩行工具欄和標籤）
	editor        *widget.Entry        // 文字編輯器元件
	statusLabel   *widget.Label        // 狀態標籤
	
	// 中文輸入增強
	chineseInputEnhancer *ChineseInputEnhancer // 中文輸入增強器
	enableChineseInput   bool                  // 是否啟用中文輸入增強
	
	// 服務依賴
	editorService        services.EditorService        // 編輯器服務
	chineseInputService  services.ChineseInputService  // 中文輸入服務
	
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
// 3. 建立中文輸入增強器和服務
// 4. 建立並配置所有 UI 元件
// 5. 設定事件處理和回調函數
// 6. 組合完整的編輯器佈局
// 7. 回傳配置完成的編輯器實例
func NewMarkdownEditor(editorService services.EditorService) *MarkdownEditor {
	// 建立 MarkdownEditor 實例
	editor := &MarkdownEditor{
		editorService:       editorService,
		chineseInputService: services.NewChineseInputService(),
		enableChineseInput:  true, // 預設啟用中文輸入增強
		isModified:          false,
	}
	
	// 建立中文輸入增強器
	if editor.enableChineseInput {
		editor.chineseInputEnhancer = NewChineseInputEnhancer()
		editor.setupChineseInputIntegration()
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
// 包含常用的 Markdown 格式化按鈕和編輯功能，採用兩行佈局並加上文字說明
//
// 執行流程：
// 1. 建立第一行工具欄：標題格式化和文字格式化按鈕
// 2. 建立第二行工具欄：列表、連結、程式碼和操作按鈕
// 3. 為每個按鈕添加文字說明，避免圖示混淆
// 4. 組合兩行工具欄到垂直容器中
func (me *MarkdownEditor) createToolbar() {
	// 第一行工具欄：標題和文字格式化
	firstRow := widget.NewToolbar(
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
	)
	
	// 第二行工具欄：列表、連結、程式碼和操作
	secondRow := widget.NewToolbar(
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
	
	// 建立文字說明標籤
	firstRowLabels := container.NewHBox(
		widget.NewLabel("H1"),
		widget.NewLabel("H2"), 
		widget.NewLabel("H3"),
		widget.NewLabel(""),  // 分隔線對應的空白
		widget.NewLabel("粗體"),
		widget.NewLabel("斜體"),
		widget.NewLabel("刪除線"),
	)
	
	secondRowLabels := container.NewHBox(
		widget.NewLabel("無序列表"),
		widget.NewLabel("有序列表"),
		widget.NewLabel(""),  // 分隔線對應的空白
		widget.NewLabel("連結"),
		widget.NewLabel("圖片"),
		widget.NewLabel(""),  // 分隔線對應的空白
		widget.NewLabel("行內程式碼"),
		widget.NewLabel("程式碼區塊"),
		widget.NewLabel(""),  // 分隔線對應的空白
		widget.NewLabel("保存"),
		widget.NewLabel("預覽"),
	)
	
	// 設定標籤樣式
	for _, label := range []*widget.Label{
		firstRowLabels.Objects[0].(*widget.Label),
		firstRowLabels.Objects[1].(*widget.Label),
		firstRowLabels.Objects[2].(*widget.Label),
		firstRowLabels.Objects[4].(*widget.Label),
		firstRowLabels.Objects[5].(*widget.Label),
		firstRowLabels.Objects[6].(*widget.Label),
	} {
		label.TextStyle = fyne.TextStyle{Italic: true}
		label.Alignment = fyne.TextAlignCenter
	}
	
	for _, label := range []*widget.Label{
		secondRowLabels.Objects[0].(*widget.Label),
		secondRowLabels.Objects[1].(*widget.Label),
		secondRowLabels.Objects[3].(*widget.Label),
		secondRowLabels.Objects[4].(*widget.Label),
		secondRowLabels.Objects[6].(*widget.Label),
		secondRowLabels.Objects[7].(*widget.Label),
		secondRowLabels.Objects[9].(*widget.Label),
		secondRowLabels.Objects[10].(*widget.Label),
	} {
		label.TextStyle = fyne.TextStyle{Italic: true}
		label.Alignment = fyne.TextAlignCenter
	}
	
	// 組合工具欄和標籤
	me.toolbar = container.NewVBox(
		firstRow,
		firstRowLabels,
		widget.NewSeparator(),
		secondRow,
		secondRowLabels,
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
	
	// 如果啟用中文輸入增強，使用增強器的文字輸入元件
	if me.enableChineseInput && me.chineseInputEnhancer != nil {
		// 將標準編輯器替換為中文增強編輯器
		me.replaceWithChineseEnhancedEditor()
	}
}

// setupChineseInputIntegration 設定中文輸入整合
// 配置中文輸入增強器與編輯器的整合
//
// 執行流程：
// 1. 設定中文輸入增強器的回調函數
// 2. 配置中文輸入相關的設定
// 3. 整合中文輸入服務
func (me *MarkdownEditor) setupChineseInputIntegration() {
	if me.chineseInputEnhancer == nil {
		return
	}
	
	// 設定文字變更回調
	me.chineseInputEnhancer.SetOnTextChanged(func(text string) {
		me.handleChineseTextChanged(text)
	})
	
	// 設定組合文字變更回調
	me.chineseInputEnhancer.SetOnCompositionChanged(func(text string) {
		me.handleCompositionChanged(text)
	})
	
	// 設定候選字選擇回調
	me.chineseInputEnhancer.SetOnCandidateSelected(func(word string) {
		me.handleCandidateSelected(word)
	})
	
	// 配置中文輸入設定
	me.chineseInputEnhancer.SetShowCandidates(true)
	me.chineseInputEnhancer.SetAutoComplete(true)
	me.chineseInputEnhancer.SetFontName("PingFang TC")
	me.chineseInputEnhancer.SetFontSize(14.0)
}

// replaceWithChineseEnhancedEditor 將標準編輯器替換為中文增強編輯器
// 使用中文輸入增強器的文字輸入元件替換標準編輯器
//
// 執行流程：
// 1. 保存標準編輯器的設定
// 2. 將設定應用到中文增強編輯器
// 3. 替換編輯器元件
func (me *MarkdownEditor) replaceWithChineseEnhancedEditor() {
	if me.chineseInputEnhancer == nil {
		return
	}
	
	// 保存標準編輯器的設定
	placeholder := me.editor.PlaceHolder
	wrapping := me.editor.Wrapping
	
	// 取得中文增強編輯器的文字輸入元件
	enhancedEditor := me.chineseInputEnhancer.GetTextEntry()
	
	// 應用設定到增強編輯器
	enhancedEditor.SetPlaceHolder(placeholder)
	enhancedEditor.Wrapping = wrapping
	enhancedEditor.Scroll = container.ScrollBoth
	
	// 設定事件處理
	enhancedEditor.OnChanged = func(content string) {
		me.onTextChanged(content)
	}
	
	enhancedEditor.OnSubmitted = func(content string) {
		me.handleEnterKey()
	}
	
	// 替換編輯器元件
	me.editor = enhancedEditor
}

// handleChineseTextChanged 處理中文文字變更事件
// 參數：text（變更後的文字內容）
//
// 執行流程：
// 1. 分析文字的中文內容
// 2. 更新中文輸入狀態
// 3. 觸發標準的文字變更處理
func (me *MarkdownEditor) handleChineseTextChanged(text string) {
	// 分析文字組成
	if me.chineseInputService != nil {
		composition := me.chineseInputService.AnalyzeTextComposition(text)
		
		// 更新狀態顯示，包含中文字符統計
		me.updateChineseInputStatus(composition)
	}
	
	// 觸發標準的文字變更處理
	me.onTextChanged(text)
}

// handleCompositionChanged 處理組合文字變更事件
// 參數：text（組合文字內容）
//
// 執行流程：
// 1. 更新組合文字顯示
// 2. 提供候選字建議（如果適用）
func (me *MarkdownEditor) handleCompositionChanged(text string) {
	// 更新狀態顯示
	me.updateStatus(fmt.Sprintf("組合輸入: %s", text))
}

// handleCandidateSelected 處理候選字選擇事件
// 參數：word（選擇的候選字詞）
//
// 執行流程：
// 1. 記錄候選字選擇
// 2. 更新詞彙頻率（如果適用）
// 3. 觸發文字變更處理
func (me *MarkdownEditor) handleCandidateSelected(word string) {
	// 更新狀態顯示
	me.updateStatus(fmt.Sprintf("已選擇: %s", word))
	
	// 如果有中文輸入服務，可以記錄詞彙使用
	if me.chineseInputService != nil {
		// 這裡可以實作詞彙使用統計
		me.chineseInputService.AddCustomWord(word)
	}
}

// updateChineseInputStatus 更新中文輸入狀態顯示
// 參數：composition（文字組成分析結果）
//
// 執行流程：
// 1. 格式化中文輸入狀態資訊
// 2. 更新狀態標籤顯示
func (me *MarkdownEditor) updateChineseInputStatus(composition services.TextComposition) {
	statusText := fmt.Sprintf("字符: %d | 中文: %d (%.1f%%) | 英文: %d | 數字: %d",
		composition.TotalCharacters,
		composition.ChineseCharacters,
		composition.ChineseRatio*100,
		composition.EnglishCharacters,
		composition.NumberCharacters)
	
	me.updateStatus(statusText)
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
		me.toolbar,                    // 工具欄容器在頂部（包含兩行工具欄和標籤）
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
// 處理保存操作並觸發適當的回調
//
// 執行流程：
// 1. 嘗試保存當前筆記
// 2. 觸發保存請求回調，讓上層處理結果
func (me *MarkdownEditor) saveContent() {
	err := me.SaveNote()
	if err != nil {
		// 錯誤處理由上層負責，這裡只觸發保存請求回調
		if me.onSaveRequested != nil {
			me.onSaveRequested()
		}
	} else {
		// 保存成功，觸發保存請求回調
		if me.onSaveRequested != nil {
			me.onSaveRequested()
		}
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
	if me.enableChineseInput && me.chineseInputEnhancer != nil {
		me.chineseInputEnhancer.Focus()
	} else if me.editor != nil {
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

// ApplyFormat 應用格式化到編輯器內容
// 參數：prefix（前綴標記）、suffix（後綴標記）、placeholder（佔位文字）
//
// 執行流程：
// 1. 檢查是否有選取的文字
// 2. 如果有選取，包圍選取的文字
// 3. 如果沒有選取，插入佔位文字
// 4. 觸發內容變更事件
func (me *MarkdownEditor) ApplyFormat(prefix, suffix, placeholder string) {
	if suffix != "" {
		// 使用包圍格式（如粗體、斜體）
		me.wrapSelection(prefix, suffix, placeholder)
	} else {
		// 使用插入格式（如標題）
		me.insertMarkdown(prefix, "", placeholder)
	}
}

// InsertText 在編輯器中插入文字
// 參數：text（要插入的文字）
//
// 執行流程：
// 1. 取得當前游標位置
// 2. 在游標位置插入文字
// 3. 更新編輯器內容
// 4. 觸發內容變更事件
func (me *MarkdownEditor) InsertText(text string) {
	// 取得當前內容
	content := me.editor.Text
	
	// 簡化實作：在內容末尾添加文字
	// 實際實作中可以取得游標位置並在該位置插入
	newContent := content + "\n" + text
	
	// 更新編輯器內容
	me.editor.SetText(newContent)
	
	// 觸發內容變更事件
	me.onTextChanged(newContent)
}

// SetEnableChineseInput 設定是否啟用中文輸入增強
// 參數：enable（是否啟用中文輸入增強）
//
// 執行流程：
// 1. 更新中文輸入啟用狀態
// 2. 如果啟用且尚未建立增強器，則建立
// 3. 重新配置編輯器
func (me *MarkdownEditor) SetEnableChineseInput(enable bool) {
	me.enableChineseInput = enable
	
	if enable && me.chineseInputEnhancer == nil {
		me.chineseInputEnhancer = NewChineseInputEnhancer()
		me.setupChineseInputIntegration()
		me.replaceWithChineseEnhancedEditor()
	}
}

// IsChineseInputEnabled 檢查是否啟用中文輸入增強
// 回傳：是否啟用中文輸入增強的布林值
func (me *MarkdownEditor) IsChineseInputEnabled() bool {
	return me.enableChineseInput
}

// GetChineseInputEnhancer 取得中文輸入增強器
// 回傳：中文輸入增強器實例，如果未啟用則回傳 nil
func (me *MarkdownEditor) GetChineseInputEnhancer() *ChineseInputEnhancer {
	return me.chineseInputEnhancer
}

// GetChineseInputService 取得中文輸入服務
// 回傳：中文輸入服務實例
func (me *MarkdownEditor) GetChineseInputService() services.ChineseInputService {
	return me.chineseInputService
}

// AnalyzeCurrentText 分析當前文字的中文內容
// 回傳：文字組成分析結果
//
// 執行流程：
// 1. 取得當前文字內容
// 2. 使用中文輸入服務分析文字組成
// 3. 回傳分析結果
func (me *MarkdownEditor) AnalyzeCurrentText() services.TextComposition {
	if me.chineseInputService == nil {
		return services.TextComposition{}
	}
	
	currentText := me.GetContent()
	return me.chineseInputService.AnalyzeTextComposition(currentText)
}

// GetChineseCharacterCount 取得當前文字的中文字符數量
// 回傳：中文字符數量
func (me *MarkdownEditor) GetChineseCharacterCount() int {
	if me.chineseInputService == nil {
		return 0
	}
	
	currentText := me.GetContent()
	return me.chineseInputService.CountChineseCharacters(currentText)
}

// OptimizeChineseInput 優化當前的中文輸入
// 回傳：輸入法優化建議
//
// 執行流程：
// 1. 取得當前文字內容
// 2. 使用中文輸入服務分析和優化
// 3. 回傳優化建議
func (me *MarkdownEditor) OptimizeChineseInput() services.InputOptimization {
	if me.chineseInputService == nil {
		return services.InputOptimization{}
	}
	
	currentText := me.GetContent()
	return me.chineseInputService.OptimizeInputMethod(currentText)
}

// ValidateChineseInput 驗證當前的中文輸入
// 回傳：輸入驗證結果
//
// 執行流程：
// 1. 取得當前文字內容
// 2. 使用中文輸入服務驗證輸入
// 3. 回傳驗證結果
func (me *MarkdownEditor) ValidateChineseInput() services.ValidationResult {
	if me.chineseInputService == nil {
		return services.ValidationResult{IsValid: true}
	}
	
	currentText := me.GetContent()
	return me.chineseInputService.ValidateChineseInput(currentText)
}