// Package ui 包含使用者介面相關的元件和視窗管理
// 本檔案實作繁體中文輸入法優化元件，改善注音輸入法的顯示效果和用戶體驗
package ui

import (
	"fmt"                           // 格式化輸出
	"unicode/utf8"                  // UTF-8 編碼處理

	"fyne.io/fyne/v2"               // Fyne GUI 框架核心套件
	"fyne.io/fyne/v2/container"     // Fyne 容器佈局套件
	"fyne.io/fyne/v2/widget"        // Fyne UI 元件套件
)

// ChineseInputEnhancer 代表繁體中文輸入法優化器
// 提供注音輸入法顯示優化、候選字視窗和中文字符渲染改善功能
// 整合到文字編輯器中以提供更好的中文輸入體驗
type ChineseInputEnhancer struct {
	// 主要元件
	container         *fyne.Container      // 主要容器
	textEntry         *widget.Entry       // 增強的文字輸入元件
	candidateWindow   *fyne.Container     // 候選字視窗容器
	candidateList     *widget.List        // 候選字列表
	compositionLabel  *widget.Label       // 注音組合顯示標籤
	
	// 輸入狀態
	isComposing       bool                // 是否正在組合輸入
	compositionText   string              // 當前組合的注音文字
	candidateWords    []string            // 候選字詞列表
	selectedCandidate int                 // 選中的候選字索引
	
	// 設定選項
	showCandidates    bool                // 是否顯示候選字視窗
	autoComplete      bool                // 是否啟用自動完成
	fontName          string              // 中文字型名稱
	fontSize          float32             // 字型大小
	
	// 回調函數
	onTextChanged     func(text string)   // 文字變更回調
	onCompositionChanged func(text string) // 組合文字變更回調
	onCandidateSelected func(word string)  // 候選字選擇回調
}

// NewChineseInputEnhancer 建立新的繁體中文輸入法優化器實例
// 回傳：指向新建立的 ChineseInputEnhancer 的指標
//
// 執行流程：
// 1. 建立 ChineseInputEnhancer 結構體實例
// 2. 設定預設的配置選項
// 3. 初始化所有 UI 元件
// 4. 設定事件處理和回調函數
// 5. 組合完整的輸入增強器佈局
// 6. 回傳配置完成的輸入增強器實例
func NewChineseInputEnhancer() *ChineseInputEnhancer {
	enhancer := &ChineseInputEnhancer{
		// 預設設定
		showCandidates: true,
		autoComplete:   true,
		fontName:       "PingFang TC",  // macOS 繁體中文預設字型
		fontSize:       14.0,
		selectedCandidate: -1,
	}
	
	// 初始化 UI 元件
	enhancer.setupUI()
	
	// 設定事件處理
	enhancer.setupEventHandlers()
	
	return enhancer
}

// setupUI 初始化中文輸入增強器的使用者介面元件
// 建立文字輸入元件、候選字視窗和組合文字顯示
//
// 執行流程：
// 1. 建立增強的文字輸入元件
// 2. 建立候選字視窗和列表
// 3. 建立注音組合顯示標籤
// 4. 設定中文字型和樣式
// 5. 組合所有元件到主容器
func (cie *ChineseInputEnhancer) setupUI() {
	// 建立增強的文字輸入元件
	cie.createEnhancedTextEntry()
	
	// 建立候選字視窗
	cie.createCandidateWindow()
	
	// 建立注音組合顯示
	cie.createCompositionDisplay()
	
	// 組合完整佈局
	cie.assembleLayout()
}

// createEnhancedTextEntry 建立增強的文字輸入元件
// 配置支援中文輸入的多行文字編輯器
//
// 執行流程：
// 1. 建立多行文字輸入元件
// 2. 設定中文字型和渲染屬性
// 3. 配置文字選取和編輯行為
// 4. 設定輸入法相關屬性
func (cie *ChineseInputEnhancer) createEnhancedTextEntry() {
	// 建立多行文字輸入元件
	cie.textEntry = widget.NewMultiLineEntry()
	
	// 設定基本屬性
	cie.textEntry.Wrapping = fyne.TextWrapWord
	cie.textEntry.Scroll = container.ScrollBoth
	cie.textEntry.SetPlaceHolder("請輸入繁體中文內容...")
	
	// 設定中文字型
	cie.applyChineseFont()
	
	// 設定文字選取行為，優化中文字符的選取體驗
	cie.optimizeTextSelection()
}

// applyChineseFont 應用中文字型設定
// 設定適合繁體中文顯示的字型和渲染屬性
//
// 執行流程：
// 1. 建立中文字型資源
// 2. 設定字型大小和樣式
// 3. 應用到文字輸入元件
// 4. 設定字符間距和行高
func (cie *ChineseInputEnhancer) applyChineseFont() {
	// 建立自訂字型樣式，針對中文優化
	customFont := &fyne.TextStyle{
		Bold:      false,
		Italic:    false,
		Monospace: false,
	}
	
	// 設定文字樣式
	cie.textEntry.TextStyle = *customFont
	
	// 設定文字對齊方式，改善中文顯示
	cie.textEntry.MultiLine = true
	cie.textEntry.Wrapping = fyne.TextWrapWord
	
	// 針對中文字符優化的設定
	cie.optimizeChineseRendering()
}

// optimizeChineseRendering 優化中文字符渲染
// 針對中文字符的特殊渲染需求進行優化
//
// 執行流程：
// 1. 設定適合中文的行高
// 2. 優化字符間距
// 3. 改善標點符號顯示
// 4. 設定中文輸入法相關屬性
func (cie *ChineseInputEnhancer) optimizeChineseRendering() {
	// 由於 Fyne 的限制，我們主要透過設定來優化中文顯示
	// 實際的字型渲染優化需要在系統層面處理
	
	// 設定文字輸入的基本屬性
	cie.textEntry.Validator = cie.createChineseInputValidator()
	
	// 設定輸入提示，使用中文
	if cie.textEntry.PlaceHolder == "" {
		cie.textEntry.SetPlaceHolder("請輸入繁體中文內容...")
	}
}

// createChineseInputValidator 建立中文輸入驗證器
// 回傳：適用於中文輸入的驗證器
//
// 執行流程：
// 1. 建立自訂驗證器
// 2. 設定中文字符驗證規則
// 3. 處理特殊字符和標點符號
// 4. 回傳驗證器實例
func (cie *ChineseInputEnhancer) createChineseInputValidator() fyne.StringValidator {
	return func(text string) error {
		// 允許所有輸入，但可以在這裡添加特殊的中文輸入驗證邏輯
		// 例如：檢查是否包含不支援的字符、驗證輸入格式等
		
		// 目前允許所有輸入，包括中文、英文、數字和標點符號
		return nil
	}
}

// optimizeTextSelection 優化文字選取行為
// 改善中文字符的選取和編輯體驗
//
// 執行流程：
// 1. 設定字符邊界檢測
// 2. 優化雙擊選取行為
// 3. 改善游標定位精度
// 4. 設定中文標點符號處理
func (cie *ChineseInputEnhancer) optimizeTextSelection() {
	// 設定文字變更處理，用於檢測中文輸入
	cie.textEntry.OnChanged = func(text string) {
		cie.handleTextChanged(text)
	}
	
	// 設定游標位置變更處理
	cie.textEntry.OnCursorChanged = func() {
		cie.handleCursorChanged()
	}
	
	// 設定鍵盤事件處理，改善中文編輯體驗
	cie.setupChineseEditingKeyHandlers()
}

// handleCursorChanged 處理游標位置變更
// 優化中文字符的游標定位體驗
//
// 執行流程：
// 1. 檢測游標周圍的字符類型
// 2. 調整游標位置以適應中文字符
// 3. 更新選取狀態
func (cie *ChineseInputEnhancer) handleCursorChanged() {
	// 取得當前游標位置
	cursorPos := cie.textEntry.CursorColumn
	text := cie.textEntry.Text
	
	// 檢測游標周圍的中文字符
	cie.analyzeCursorContext(text, cursorPos)
}

// analyzeCursorContext 分析游標上下文
// 參數：text（文字內容）、cursorPos（游標位置）
//
// 執行流程：
// 1. 分析游標前後的字符
// 2. 檢測中文詞彙邊界
// 3. 提供上下文相關的輔助功能
func (cie *ChineseInputEnhancer) analyzeCursorContext(text string, cursorPos int) {
	if text == "" {
		return
	}
	
	runes := []rune(text)
	if cursorPos < 0 || cursorPos > len(runes) {
		return
	}
	
	// 分析游標前的字符
	var prevChar rune
	if cursorPos > 0 {
		prevChar = runes[cursorPos-1]
	}
	
	// 分析游標後的字符
	var nextChar rune
	if cursorPos < len(runes) {
		nextChar = runes[cursorPos]
	}
	
	// 檢測是否在中文詞彙中間
	if cie.isChineseCharacter(prevChar) || cie.isChineseCharacter(nextChar) {
		cie.handleChineseWordContext(text, cursorPos)
	}
}

// handleChineseWordContext 處理中文詞彙上下文
// 參數：text（文字內容）、cursorPos（游標位置）
//
// 執行流程：
// 1. 識別當前中文詞彙
// 2. 提供詞彙相關的輔助功能
// 3. 優化選取行為
func (cie *ChineseInputEnhancer) handleChineseWordContext(text string, cursorPos int) {
	// 找到當前詞彙的邊界
	wordStart, wordEnd := cie.findChineseWordBoundary(text, cursorPos)
	
	if wordStart != wordEnd {
		currentWord := string([]rune(text)[wordStart:wordEnd])
		
		// 如果當前詞彙是中文，可以提供相關建議
		if cie.containsChineseCharacters(currentWord) {
			cie.handleCurrentChineseWord(currentWord)
		}
	}
}

// findChineseWordBoundary 找到中文詞彙邊界
// 參數：text（文字內容）、cursorPos（游標位置）
// 回傳：詞彙開始位置和結束位置
func (cie *ChineseInputEnhancer) findChineseWordBoundary(text string, cursorPos int) (int, int) {
	runes := []rune(text)
	if len(runes) == 0 {
		return 0, 0
	}
	
	// 向前找詞彙開始
	start := cursorPos
	for start > 0 && cie.isChineseCharacter(runes[start-1]) {
		start--
	}
	
	// 向後找詞彙結束
	end := cursorPos
	for end < len(runes) && cie.isChineseCharacter(runes[end]) {
		end++
	}
	
	return start, end
}

// handleCurrentChineseWord 處理當前中文詞彙
// 參數：word（當前詞彙）
//
// 執行流程：
// 1. 分析詞彙特性
// 2. 提供相關建議
// 3. 更新輔助資訊
func (cie *ChineseInputEnhancer) handleCurrentChineseWord(word string) {
	// 這裡可以實作各種中文詞彙相關的輔助功能
	// 例如：同義詞建議、詞彙解釋、拼音顯示等
	
	// 暫時只觸發組合文字變更回調
	if cie.onCompositionChanged != nil {
		cie.onCompositionChanged(fmt.Sprintf("當前詞彙: %s", word))
	}
}

// setupChineseEditingKeyHandlers 設定中文編輯鍵盤處理器
// 改善中文編輯的鍵盤操作體驗
//
// 執行流程：
// 1. 設定中文特殊按鍵處理
// 2. 優化詞彙選取快捷鍵
// 3. 設定中文標點符號快捷鍵
func (cie *ChineseInputEnhancer) setupChineseEditingKeyHandlers() {
	// 由於 Fyne 的鍵盤事件處理限制，這裡主要設定基本的按鍵處理
	// 實際的中文編輯優化需要在更高層級實作
	
	// 設定提交事件處理
	cie.textEntry.OnSubmitted = func(text string) {
		cie.handleTextSubmitted(text)
	}
}

// createCandidateWindow 建立候選字視窗
// 顯示輸入法的候選字詞列表
//
// 執行流程：
// 1. 建立候選字列表元件
// 2. 設定列表項目的顯示格式
// 3. 設定選擇事件處理
// 4. 建立候選字視窗容器
// 5. 設定視窗的顯示和隱藏邏輯
func (cie *ChineseInputEnhancer) createCandidateWindow() {
	// 建立候選字列表
	cie.candidateList = widget.NewList(
		// 取得列表項目數量
		func() int {
			return len(cie.candidateWords)
		},
		// 建立列表項目 UI
		func() fyne.CanvasObject {
			// 建立候選字項目容器
			label := widget.NewLabel("")
			label.TextStyle = fyne.TextStyle{Bold: false}
			
			// 設定中文字型樣式
			cie.applyCandidateItemStyle(label)
			
			return label
		},
		// 更新列表項目內容
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if id < len(cie.candidateWords) {
				label := obj.(*widget.Label)
				
				// 顯示候選字和編號，使用更好的格式
				candidateText := cie.formatCandidateItem(id, cie.candidateWords[id])
				label.SetText(candidateText)
				
				// 高亮選中的候選字
				if id == cie.selectedCandidate {
					label.TextStyle = fyne.TextStyle{Bold: true}
					cie.applyCandidateItemStyle(label)
				} else {
					label.TextStyle = fyne.TextStyle{Bold: false}
					cie.applyCandidateItemStyle(label)
				}
			}
		},
	)
	
	// 設定候選字選擇事件
	cie.candidateList.OnSelected = func(id widget.ListItemID) {
		cie.selectCandidate(id)
	}
	
	// 建立候選字視窗標題
	candidateTitle := widget.NewLabel("🔤 候選字")
	candidateTitle.TextStyle = fyne.TextStyle{Bold: true}
	
	// 建立候選字統計資訊
	candidateInfo := widget.NewLabel("")
	candidateInfo.TextStyle = fyne.TextStyle{Italic: true}
	
	// 建立候選字視窗容器，使用更好的佈局
	cie.candidateWindow = container.NewVBox(
		candidateTitle,
		candidateInfo,
		widget.NewSeparator(),
		container.NewScroll(cie.candidateList),
		cie.createCandidateWindowControls(),
	)
	
	// 設定候選字視窗的大小限制
	cie.candidateWindow.Resize(fyne.NewSize(300, 200))
	
	// 預設隱藏候選字視窗
	cie.candidateWindow.Hide()
}

// applyCandidateItemStyle 應用候選字項目樣式
// 參數：label（候選字標籤元件）
//
// 執行流程：
// 1. 設定中文字型
// 2. 設定文字對齊
// 3. 設定顏色和樣式
func (cie *ChineseInputEnhancer) applyCandidateItemStyle(label *widget.Label) {
	// 設定文字對齊方式
	label.Alignment = fyne.TextAlignLeading
	
	// 設定文字換行
	label.Wrapping = fyne.TextWrapOff
}

// formatCandidateItem 格式化候選字項目顯示
// 參數：index（候選字索引）、word（候選字詞）
// 回傳：格式化後的顯示文字
//
// 執行流程：
// 1. 建立編號顯示
// 2. 添加候選字詞
// 3. 添加額外資訊（如頻率、拼音等）
// 4. 回傳格式化結果
func (cie *ChineseInputEnhancer) formatCandidateItem(index int, word string) string {
	// 基本格式：編號 + 候選字
	baseFormat := fmt.Sprintf("%d. %s", index+1, word)
	
	// 如果候選字是中文，可以添加額外資訊
	if cie.containsChineseCharacters(word) {
		// 添加字符數量資訊
		charCount := len([]rune(word))
		if charCount > 1 {
			baseFormat += fmt.Sprintf(" (%d字)", charCount)
		}
	}
	
	return baseFormat
}

// createCandidateWindowControls 建立候選字視窗控制項
// 回傳：候選字視窗控制項容器
//
// 執行流程：
// 1. 建立上一頁/下一頁按鈕
// 2. 建立關閉按鈕
// 3. 建立快捷鍵提示
// 4. 組合控制項佈局
func (cie *ChineseInputEnhancer) createCandidateWindowControls() *fyne.Container {
	// 建立控制按鈕
	prevButton := widget.NewButton("◀", func() {
		cie.navigateCandidates(-1)
	})
	prevButton.Resize(fyne.NewSize(30, 30))
	
	nextButton := widget.NewButton("▶", func() {
		cie.navigateCandidates(1)
	})
	nextButton.Resize(fyne.NewSize(30, 30))
	
	closeButton := widget.NewButton("✕", func() {
		cie.hideCandidateWindow()
	})
	closeButton.Resize(fyne.NewSize(30, 30))
	
	// 建立快捷鍵提示
	shortcutHint := widget.NewLabel("↑↓選擇 Enter確認 Esc關閉")
	shortcutHint.TextStyle = fyne.TextStyle{Italic: true}
	
	// 組合控制項
	controls := container.NewHBox(
		prevButton,
		nextButton,
		widget.NewSeparator(),
		closeButton,
	)
	
	return container.NewVBox(
		widget.NewSeparator(),
		controls,
		shortcutHint,
	)
}

// navigateCandidates 導航候選字列表
// 參數：direction（導航方向，-1為上，1為下）
//
// 執行流程：
// 1. 計算新的選擇索引
// 2. 檢查邊界條件
// 3. 更新選擇狀態
// 4. 刷新顯示
func (cie *ChineseInputEnhancer) navigateCandidates(direction int) {
	if len(cie.candidateWords) == 0 {
		return
	}
	
	// 計算新的選擇索引
	newIndex := cie.selectedCandidate + direction
	
	// 處理邊界條件
	if newIndex < 0 {
		newIndex = len(cie.candidateWords) - 1 // 循環到最後一個
	} else if newIndex >= len(cie.candidateWords) {
		newIndex = 0 // 循環到第一個
	}
	
	// 更新選擇狀態
	cie.selectedCandidate = newIndex
	
	// 刷新候選字列表顯示
	cie.candidateList.Refresh()
	
	// 確保選中的項目可見
	cie.candidateList.ScrollTo(newIndex)
}

// createCompositionDisplay 建立注音組合顯示
// 顯示當前正在組合的注音符號
//
// 執行流程：
// 1. 建立注音顯示標籤
// 2. 設定標籤樣式和字型
// 3. 設定顯示位置和對齊方式
// 4. 預設隱藏組合顯示
func (cie *ChineseInputEnhancer) createCompositionDisplay() {
	// 建立注音組合顯示標籤
	cie.compositionLabel = widget.NewLabel("")
	cie.compositionLabel.TextStyle = fyne.TextStyle{
		Bold:   true,
		Italic: false,
	}
	
	// 設定標籤對齊方式
	cie.compositionLabel.Alignment = fyne.TextAlignLeading
	
	// 預設隱藏組合顯示
	cie.compositionLabel.Hide()
}

// assembleLayout 組合中文輸入增強器的完整佈局
// 將所有元件組合成完整的輸入增強器介面
//
// 執行流程：
// 1. 建立主要垂直佈局容器
// 2. 添加注音組合顯示（頂部）
// 3. 添加文字輸入元件（中間主要區域）
// 4. 添加候選字視窗（底部，可隱藏）
// 5. 設定佈局比例和間距
func (cie *ChineseInputEnhancer) assembleLayout() {
	// 建立輸入區域容器（注音顯示 + 文字輸入）
	inputArea := container.NewVBox(
		cie.compositionLabel,  // 注音組合顯示
		cie.textEntry,         // 文字輸入元件
	)
	
	// 建立主要容器
	cie.container = container.NewVBox(
		inputArea,             // 輸入區域
		cie.candidateWindow,   // 候選字視窗
	)
}

// setupEventHandlers 設定事件處理器
// 配置鍵盤事件、滑鼠事件和輸入法事件的處理
//
// 執行流程：
// 1. 設定鍵盤快捷鍵處理
// 2. 設定輸入法組合事件處理
// 3. 設定候選字選擇事件處理
// 4. 設定文字變更事件處理
func (cie *ChineseInputEnhancer) setupEventHandlers() {
	// 設定鍵盤事件處理
	cie.setupKeyboardHandlers()
	
	// 設定輸入法事件處理
	cie.setupIMEHandlers()
}

// setupKeyboardHandlers 設定鍵盤事件處理器
// 處理方向鍵、Enter 鍵和數字鍵的候選字選擇
//
// 執行流程：
// 1. 設定方向鍵處理（上下選擇候選字）
// 2. 設定 Enter 鍵處理（確認選擇）
// 3. 設定數字鍵處理（快速選擇候選字）
// 4. 設定 Escape 鍵處理（取消輸入）
func (cie *ChineseInputEnhancer) setupKeyboardHandlers() {
	// 由於 Fyne 的鍵盤事件處理限制，這裡主要設定基本的按鍵處理
	// 實際的候選字選擇會透過滑鼠點擊或其他方式實作
	
	// 設定文字提交事件（Enter 鍵）
	cie.textEntry.OnSubmitted = func(text string) {
		cie.handleTextSubmitted(text)
	}
}

// setupIMEHandlers 設定輸入法事件處理器
// 處理輸入法的組合開始、更新和結束事件
//
// 執行流程：
// 1. 監聽輸入法組合開始事件
// 2. 處理組合文字更新事件
// 3. 處理組合結束和文字確認事件
// 4. 更新候選字列表和顯示
func (cie *ChineseInputEnhancer) setupIMEHandlers() {
	// 由於 Fyne 對輸入法事件的支援有限，這裡主要透過文字變更事件
	// 來檢測和處理中文輸入的狀態變化
	
	// 實際的輸入法處理會在 handleTextChanged 中實作
}

// handleTextChanged 處理文字變更事件
// 檢測中文輸入狀態並更新相關顯示
//
// 執行流程：
// 1. 檢測是否為中文輸入
// 2. 分析輸入的字符類型
// 3. 更新候選字列表
// 4. 觸發相關回調函數
func (cie *ChineseInputEnhancer) handleTextChanged(text string) {
	// 檢測是否包含中文字符
	hasChineseChars := cie.containsChineseCharacters(text)
	
	// 分析最後輸入的字符
	lastChar := cie.getLastCharacter(text)
	
	// 檢測是否可能是注音輸入
	if cie.isPossibleZhuyinInput(text) {
		cie.handleZhuyinInput(text)
	}
	
	// 如果是中文字符，更新中文輸入狀態
	if hasChineseChars {
		cie.updateChineseInputStatus(text)
	}
	
	// 處理自動完成
	if cie.autoComplete && text != "" {
		cie.handleAutoComplete(text, lastChar)
	}
	
	// 觸發文字變更回調
	if cie.onTextChanged != nil {
		cie.onTextChanged(text)
	}
}

// isPossibleZhuyinInput 檢測是否可能是注音輸入
// 參數：text（輸入文字）
// 回傳：是否可能是注音輸入
//
// 執行流程：
// 1. 檢查是否包含注音符號
// 2. 分析注音組合模式
// 3. 判斷是否為有效的注音輸入
func (cie *ChineseInputEnhancer) isPossibleZhuyinInput(text string) bool {
	if text == "" {
		return false
	}
	
	// 檢查最後幾個字符是否包含注音符號
	lastPart := cie.getLastInputPart(text, 10) // 取最後10個字符
	
	// 注音符號 Unicode 範圍 (U+3105-U+312F)
	for _, r := range lastPart {
		if r >= 0x3105 && r <= 0x312F {
			return true
		}
	}
	
	return false
}

// getLastInputPart 取得最後輸入的部分
// 參數：text（完整文字）、maxLength（最大長度）
// 回傳：最後輸入的部分
func (cie *ChineseInputEnhancer) getLastInputPart(text string, maxLength int) string {
	runes := []rune(text)
	if len(runes) <= maxLength {
		return text
	}
	return string(runes[len(runes)-maxLength:])
}

// handleZhuyinInput 處理注音輸入
// 參數：text（包含注音的文字）
//
// 執行流程：
// 1. 提取注音部分
// 2. 分析注音組合
// 3. 生成候選字
// 4. 更新候選字視窗
func (cie *ChineseInputEnhancer) handleZhuyinInput(text string) {
	// 提取可能的注音組合
	zhuyinPart := cie.extractZhuyinPart(text)
	
	if zhuyinPart != "" {
		// 設定組合狀態
		cie.isComposing = true
		cie.compositionText = zhuyinPart
		
		// 顯示組合文字
		cie.showCompositionText(zhuyinPart)
		
		// 生成候選字（這裡是簡化的實作）
		candidates := cie.generateZhuyinCandidates(zhuyinPart)
		if len(candidates) > 0 {
			cie.candidateWords = candidates
			cie.showCandidateWindow()
		}
	}
}

// extractZhuyinPart 提取注音部分
// 參數：text（完整文字）
// 回傳：注音部分
func (cie *ChineseInputEnhancer) extractZhuyinPart(text string) string {
	// 簡化的注音提取邏輯
	// 實際應用中需要更複雜的注音分析
	
	var zhuyinPart []rune
	runes := []rune(text)
	
	// 從後往前找注音符號
	for i := len(runes) - 1; i >= 0; i-- {
		r := runes[i]
		if r >= 0x3105 && r <= 0x312F { // 注音符號範圍
			zhuyinPart = append([]rune{r}, zhuyinPart...)
		} else if len(zhuyinPart) > 0 {
			// 遇到非注音符號且已有注音，停止
			break
		}
	}
	
	return string(zhuyinPart)
}

// generateZhuyinCandidates 生成注音候選字
// 參數：zhuyin（注音組合）
// 回傳：候選字列表
func (cie *ChineseInputEnhancer) generateZhuyinCandidates(zhuyin string) []string {
	// 簡化的注音候選字生成
	// 實際應用中需要完整的注音字典
	
	zhuyinMap := map[string][]string{
		"ㄋㄧˇ":   {"你", "妳", "尼"},
		"ㄏㄠˇ":   {"好", "號", "豪"},
		"ㄕˋ":    {"是", "事", "世"},
		"ㄐㄧㄝˋ": {"界", "借", "戒"},
		"ㄓㄨㄥ":  {"中", "鐘", "忠"},
		"ㄨㄣˊ":   {"文", "聞", "溫"},
	}
	
	if candidates, exists := zhuyinMap[zhuyin]; exists {
		return candidates
	}
	
	// 如果沒有完全匹配，嘗試部分匹配
	for key, candidates := range zhuyinMap {
		if len(zhuyin) >= 2 && len(key) >= 2 && key[:6] == zhuyin[:6] { // 比較前兩個字符
			return candidates
		}
	}
	
	return []string{}
}

// showCompositionText 顯示組合文字
// 參數：compositionText（組合文字）
//
// 執行流程：
// 1. 更新組合標籤內容
// 2. 設定標籤可見性
// 3. 調整標籤位置
func (cie *ChineseInputEnhancer) showCompositionText(compositionText string) {
	cie.compositionLabel.SetText("組合中: " + compositionText)
	cie.compositionLabel.Show()
	
	// 觸發組合文字變更回調
	if cie.onCompositionChanged != nil {
		cie.onCompositionChanged(compositionText)
	}
}

// handleAutoComplete 處理自動完成
// 參數：text（當前文字）、lastChar（最後字符）
//
// 執行流程：
// 1. 分析輸入上下文
// 2. 生成自動完成建議
// 3. 更新候選字列表
// 4. 顯示自動完成提示
func (cie *ChineseInputEnhancer) handleAutoComplete(text string, lastChar rune) {
	// 如果正在組合輸入，不進行自動完成
	if cie.isComposing {
		return
	}
	
	// 取得當前詞彙的前綴
	prefix := cie.getCurrentWordPrefix(text)
	
	if len(prefix) >= 1 && cie.isChineseCharacter(lastChar) {
		// 生成自動完成建議
		suggestions := cie.generateAutoCompleteSuggestions(prefix)
		
		if len(suggestions) > 0 {
			cie.candidateWords = suggestions
			cie.showCandidateWindow()
		}
	}
}

// getCurrentWordPrefix 取得當前詞彙的前綴
// 參數：text（完整文字）
// 回傳：當前詞彙前綴
func (cie *ChineseInputEnhancer) getCurrentWordPrefix(text string) string {
	if text == "" {
		return ""
	}
	
	runes := []rune(text)
	var prefix []rune
	
	// 從後往前找，直到遇到空格或標點符號
	for i := len(runes) - 1; i >= 0; i-- {
		r := runes[i]
		if cie.isChineseCharacter(r) || (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
			prefix = append([]rune{r}, prefix...)
		} else {
			break
		}
	}
	
	return string(prefix)
}

// generateAutoCompleteSuggestions 生成自動完成建議
// 參數：prefix（前綴）
// 回傳：自動完成建議列表
func (cie *ChineseInputEnhancer) generateAutoCompleteSuggestions(prefix string) []string {
	// 常用詞彙自動完成
	commonCompletions := map[string][]string{
		"你": {"你好", "你們", "你的"},
		"我": {"我們", "我的", "我是"},
		"這": {"這個", "這些", "這樣"},
		"那": {"那個", "那些", "那樣"},
		"什": {"什麼", "什麼時候"},
		"怎": {"怎麼", "怎樣", "怎麼辦"},
		"為": {"為什麼", "為了", "為何"},
		"可": {"可以", "可能", "可是"},
		"應": {"應該", "應當", "應用"},
		"需": {"需要", "需求", "需求"},
	}
	
	if suggestions, exists := commonCompletions[prefix]; exists {
		return suggestions
	}
	
	// 如果沒有預定義的建議，回傳空列表
	return []string{}
}

// containsChineseCharacters 檢查文字是否包含中文字符
// 參數：text（要檢查的文字）
// 回傳：是否包含中文字符的布林值
//
// 執行流程：
// 1. 遍歷文字中的每個字符
// 2. 檢查字符是否在中文 Unicode 範圍內
// 3. 回傳檢查結果
func (cie *ChineseInputEnhancer) containsChineseCharacters(text string) bool {
	for _, r := range text {
		// 檢查是否為中文字符（包含繁體、簡體和擴展字符）
		if cie.isChineseCharacter(r) {
			return true
		}
	}
	return false
}

// isChineseCharacter 檢查單個字符是否為中文字符
// 參數：r（要檢查的字符）
// 回傳：是否為中文字符的布林值
//
// 執行流程：
// 1. 檢查字符是否在 CJK 統一漢字範圍
// 2. 檢查字符是否在 CJK 擴展範圍
// 3. 檢查字符是否在中文標點符號範圍
// 4. 回傳檢查結果
func (cie *ChineseInputEnhancer) isChineseCharacter(r rune) bool {
	// CJK 統一漢字基本區塊 (U+4E00-U+9FFF)
	if r >= 0x4E00 && r <= 0x9FFF {
		return true
	}
	
	// CJK 統一漢字擴展 A 區塊 (U+3400-U+4DBF)
	if r >= 0x3400 && r <= 0x4DBF {
		return true
	}
	
	// CJK 統一漢字擴展 B 區塊 (U+20000-U+2A6DF)
	if r >= 0x20000 && r <= 0x2A6DF {
		return true
	}
	
	// CJK 相容漢字 (U+F900-U+FAFF)
	if r >= 0xF900 && r <= 0xFAFF {
		return true
	}
	
	// 中文標點符號 (U+3000-U+303F)
	if r >= 0x3000 && r <= 0x303F {
		return true
	}
	
	return false
}

// getLastCharacter 取得文字中的最後一個字符
// 參數：text（文字內容）
// 回傳：最後一個字符，如果文字為空則回傳空字符
//
// 執行流程：
// 1. 檢查文字是否為空
// 2. 使用 UTF-8 解碼取得最後一個字符
// 3. 回傳字符結果
func (cie *ChineseInputEnhancer) getLastCharacter(text string) rune {
	if text == "" {
		return 0
	}
	
	// 使用 UTF-8 解碼取得最後一個字符
	lastRune, _ := utf8.DecodeLastRuneInString(text)
	return lastRune
}

// updateChineseInputStatus 更新中文輸入狀態
// 參數：text（當前文字內容）
//
// 執行流程：
// 1. 分析中文輸入的狀態
// 2. 更新候選字列表（如果適用）
// 3. 更新組合文字顯示
// 4. 調整 UI 元件的可見性
func (cie *ChineseInputEnhancer) updateChineseInputStatus(text string) {
	// 統計中文字符數量
	chineseCharCount := cie.countChineseCharacters(text)
	
	// 如果有中文字符，可以提供一些輔助功能
	if chineseCharCount > 0 {
		// 這裡可以實作一些中文輸入的輔助功能
		// 例如：字數統計、常用詞建議等
		cie.updateChineseInputHelpers(text, chineseCharCount)
	}
}

// countChineseCharacters 統計文字中的中文字符數量
// 參數：text（要統計的文字）
// 回傳：中文字符的數量
//
// 執行流程：
// 1. 遍歷文字中的每個字符
// 2. 檢查並統計中文字符
// 3. 回傳統計結果
func (cie *ChineseInputEnhancer) countChineseCharacters(text string) int {
	count := 0
	for _, r := range text {
		if cie.isChineseCharacter(r) {
			count++
		}
	}
	return count
}

// updateChineseInputHelpers 更新中文輸入輔助功能
// 參數：text（當前文字）、chineseCharCount（中文字符數量）
//
// 執行流程：
// 1. 更新字數統計顯示
// 2. 提供常用詞建議（如果啟用）
// 3. 更新輸入法狀態指示
func (cie *ChineseInputEnhancer) updateChineseInputHelpers(text string, chineseCharCount int) {
	// 這裡可以實作各種中文輸入輔助功能
	// 例如：
	// - 顯示中文字數統計
	// - 提供常用詞彙建議
	// - 檢測並提示可能的錯字
	
	// 暫時只更新一些基本的狀態資訊
	if cie.onCompositionChanged != nil {
		cie.onCompositionChanged(fmt.Sprintf("中文字符: %d", chineseCharCount))
	}
}

// handleTextSubmitted 處理文字提交事件
// 參數：text（提交的文字內容）
//
// 執行流程：
// 1. 處理當前的組合狀態
// 2. 確認候選字選擇
// 3. 清理輸入狀態
// 4. 觸發相關回調
func (cie *ChineseInputEnhancer) handleTextSubmitted(text string) {
	// 如果正在組合輸入，結束組合
	if cie.isComposing {
		cie.finishComposition()
	}
	
	// 隱藏候選字視窗
	cie.hideCandidateWindow()
	
	// 觸發文字變更回調
	if cie.onTextChanged != nil {
		cie.onTextChanged(text)
	}
}

// showCandidateWindow 顯示候選字視窗
// 顯示輸入法的候選字詞列表
//
// 執行流程：
// 1. 更新候選字列表內容
// 2. 設定視窗可見性
// 3. 調整視窗位置和大小
func (cie *ChineseInputEnhancer) showCandidateWindow() {
	if cie.showCandidates && len(cie.candidateWords) > 0 {
		cie.candidateWindow.Show()
		cie.candidateList.Refresh()
	}
}

// hideCandidateWindow 隱藏候選字視窗
// 隱藏輸入法的候選字詞列表
//
// 執行流程：
// 1. 設定視窗不可見
// 2. 清理選擇狀態
func (cie *ChineseInputEnhancer) hideCandidateWindow() {
	cie.candidateWindow.Hide()
	cie.selectedCandidate = -1
}

// selectCandidate 選擇候選字
// 參數：index（候選字的索引）
//
// 執行流程：
// 1. 驗證索引的有效性
// 2. 取得選擇的候選字
// 3. 插入候選字到文字中
// 4. 清理輸入狀態
// 5. 觸發選擇回調
func (cie *ChineseInputEnhancer) selectCandidate(index int) {
	if index < 0 || index >= len(cie.candidateWords) {
		return
	}
	
	// 取得選擇的候選字
	selectedWord := cie.candidateWords[index]
	
	// 插入候選字到文字中
	cie.insertCandidateWord(selectedWord)
	
	// 清理輸入狀態
	cie.clearComposition()
	
	// 隱藏候選字視窗
	cie.hideCandidateWindow()
	
	// 觸發候選字選擇回調
	if cie.onCandidateSelected != nil {
		cie.onCandidateSelected(selectedWord)
	}
}

// insertCandidateWord 插入候選字到文字中
// 參數：word（要插入的候選字詞）
//
// 執行流程：
// 1. 取得當前游標位置
// 2. 插入候選字到適當位置
// 3. 更新游標位置
// 4. 刷新文字顯示
func (cie *ChineseInputEnhancer) insertCandidateWord(word string) {
	// 取得當前文字內容
	currentText := cie.textEntry.Text
	
	// 簡單的插入實作：在文字末尾添加候選字
	// 實際應用中可能需要更複雜的游標位置處理
	newText := currentText + word
	
	// 更新文字內容
	cie.textEntry.SetText(newText)
}

// clearComposition 清理組合輸入狀態
// 重置所有與組合輸入相關的狀態
//
// 執行流程：
// 1. 重置組合狀態標誌
// 2. 清空組合文字
// 3. 清空候選字列表
// 4. 隱藏相關 UI 元件
func (cie *ChineseInputEnhancer) clearComposition() {
	cie.isComposing = false
	cie.compositionText = ""
	cie.candidateWords = nil
	cie.selectedCandidate = -1
	
	// 隱藏組合顯示
	cie.compositionLabel.Hide()
}

// finishComposition 完成組合輸入
// 結束當前的組合輸入過程
//
// 執行流程：
// 1. 確認當前組合的文字
// 2. 清理組合狀態
// 3. 更新文字內容
func (cie *ChineseInputEnhancer) finishComposition() {
	if cie.isComposing && cie.compositionText != "" {
		// 將組合文字添加到主文字中
		currentText := cie.textEntry.Text
		newText := currentText + cie.compositionText
		cie.textEntry.SetText(newText)
	}
	
	// 清理組合狀態
	cie.clearComposition()
}

// GetContainer 取得中文輸入增強器的主要容器
// 回傳：增強器的 fyne.Container 實例
// 用於將增強器嵌入到其他 UI 佈局中
func (cie *ChineseInputEnhancer) GetContainer() *fyne.Container {
	return cie.container
}

// GetTextEntry 取得文字輸入元件
// 回傳：增強的文字輸入元件實例
// 用於直接存取文字輸入功能
func (cie *ChineseInputEnhancer) GetTextEntry() *widget.Entry {
	return cie.textEntry
}

// SetText 設定文字內容
// 參數：text（要設定的文字內容）
//
// 執行流程：
// 1. 設定文字輸入元件的內容
// 2. 觸發文字變更處理
func (cie *ChineseInputEnhancer) SetText(text string) {
	cie.textEntry.SetText(text)
}

// GetText 取得文字內容
// 回傳：當前的文字內容
func (cie *ChineseInputEnhancer) GetText() string {
	return cie.textEntry.Text
}

// SetPlaceHolder 設定佔位符文字
// 參數：placeholder（佔位符文字）
func (cie *ChineseInputEnhancer) SetPlaceHolder(placeholder string) {
	cie.textEntry.SetPlaceHolder(placeholder)
}

// SetShowCandidates 設定是否顯示候選字視窗
// 參數：show（是否顯示候選字視窗）
func (cie *ChineseInputEnhancer) SetShowCandidates(show bool) {
	cie.showCandidates = show
	if !show {
		cie.hideCandidateWindow()
	}
}

// SetAutoComplete 設定是否啟用自動完成
// 參數：enable（是否啟用自動完成）
func (cie *ChineseInputEnhancer) SetAutoComplete(enable bool) {
	cie.autoComplete = enable
}

// SetFontName 設定中文字型名稱
// 參數：fontName（字型名稱）
func (cie *ChineseInputEnhancer) SetFontName(fontName string) {
	cie.fontName = fontName
	cie.applyChineseFont()
}

// SetFontSize 設定字型大小
// 參數：fontSize（字型大小）
func (cie *ChineseInputEnhancer) SetFontSize(fontSize float32) {
	cie.fontSize = fontSize
	cie.applyChineseFont()
}

// SetOnTextChanged 設定文字變更回調函數
// 參數：callback（文字變更時的回調函數）
func (cie *ChineseInputEnhancer) SetOnTextChanged(callback func(text string)) {
	cie.onTextChanged = callback
}

// SetOnCompositionChanged 設定組合文字變更回調函數
// 參數：callback（組合文字變更時的回調函數）
func (cie *ChineseInputEnhancer) SetOnCompositionChanged(callback func(text string)) {
	cie.onCompositionChanged = callback
}

// SetOnCandidateSelected 設定候選字選擇回調函數
// 參數：callback（候選字選擇時的回調函數）
func (cie *ChineseInputEnhancer) SetOnCandidateSelected(callback func(word string)) {
	cie.onCandidateSelected = callback
}

// Focus 設定輸入焦點
// 讓文字輸入元件獲得輸入焦點
func (cie *ChineseInputEnhancer) Focus() {
	cie.textEntry.FocusGained()
}

// IsComposing 檢查是否正在組合輸入
// 回傳：是否正在組合輸入的布林值
func (cie *ChineseInputEnhancer) IsComposing() bool {
	return cie.isComposing
}

// GetCompositionText 取得當前組合文字
// 回傳：當前組合的文字內容
func (cie *ChineseInputEnhancer) GetCompositionText() string {
	return cie.compositionText
}

// GetCandidateWords 取得候選字列表
// 回傳：當前的候選字詞列表
func (cie *ChineseInputEnhancer) GetCandidateWords() []string {
	return cie.candidateWords
}

// GetSelectedCandidate 取得選中的候選字索引
// 回傳：選中的候選字索引，-1 表示沒有選中
func (cie *ChineseInputEnhancer) GetSelectedCandidate() int {
	return cie.selectedCandidate
}