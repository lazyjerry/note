// Package ui 包含使用者介面相關的元件和視窗管理
// 本檔案包含繁體中文輸入法優化元件的單元測試
package ui

import (
	"testing"
)

// TestNewChineseInputEnhancer 測試中文輸入增強器的建立
// 驗證中文輸入增強器是否正確初始化
func TestNewChineseInputEnhancer(t *testing.T) {
	// 建立中文輸入增強器
	enhancer := NewChineseInputEnhancer()
	
	// 驗證基本屬性
	if enhancer == nil {
		t.Fatal("中文輸入增強器建立失敗")
	}
	
	if enhancer.GetContainer() == nil {
		t.Error("主要容器未正確建立")
	}
	
	if enhancer.GetTextEntry() == nil {
		t.Error("文字輸入元件未正確建立")
	}
	
	// 驗證預設設定
	if !enhancer.showCandidates {
		t.Error("預期預設顯示候選字視窗")
	}
	
	if !enhancer.autoComplete {
		t.Error("預期預設啟用自動完成")
	}
	
	if enhancer.fontName != "PingFang TC" {
		t.Errorf("預期預設字型為 'PingFang TC'，實際為 '%s'", enhancer.fontName)
	}
	
	if enhancer.fontSize != 14.0 {
		t.Errorf("預期預設字型大小為 14.0，實際為 %f", enhancer.fontSize)
	}
}

// TestChineseCharacterDetection 測試中文字符檢測功能
// 驗證中文字符的識別和統計功能
func TestChineseCharacterDetection(t *testing.T) {
	enhancer := NewChineseInputEnhancer()
	
	// 測試中文字符檢測
	testCases := []struct {
		text     string
		expected bool
		desc     string
	}{
		{"Hello World", false, "純英文文字"},
		{"你好世界", true, "純中文文字"},
		{"Hello 你好", true, "中英混合文字"},
		{"123456", false, "純數字"},
		{"你好123", true, "中文數字混合"},
		{"", false, "空字串"},
		{"！？。，", true, "中文標點符號"},
		{"Hello, World!", false, "英文標點符號"},
	}
	
	for _, tc := range testCases {
		result := enhancer.containsChineseCharacters(tc.text)
		if result != tc.expected {
			t.Errorf("%s: 預期 %v，實際 %v", tc.desc, tc.expected, result)
		}
	}
}

// TestIndividualChineseCharacterCheck 測試單個中文字符檢查
// 驗證單個字符的中文識別功能
func TestIndividualChineseCharacterCheck(t *testing.T) {
	enhancer := NewChineseInputEnhancer()
	
	// 測試各種字符類型
	testCases := []struct {
		char     rune
		expected bool
		desc     string
	}{
		{'你', true, "常用中文字符"},
		{'好', true, "常用中文字符"},
		{'A', false, "英文字母"},
		{'1', false, "數字字符"},
		{'！', true, "中文標點符號"},
		{'!', false, "英文標點符號"},
		{'　', true, "中文空格"},
		{' ', false, "英文空格"},
		{'龍', true, "繁體中文字符"},
		{'龙', true, "簡體中文字符"},
	}
	
	for _, tc := range testCases {
		result := enhancer.isChineseCharacter(tc.char)
		if result != tc.expected {
			t.Errorf("%s (%c): 預期 %v，實際 %v", tc.desc, tc.char, tc.expected, result)
		}
	}
}

// TestChineseCharacterCounting 測試中文字符統計功能
// 驗證中文字符數量的統計準確性
func TestChineseCharacterCounting(t *testing.T) {
	enhancer := NewChineseInputEnhancer()
	
	// 測試字符統計
	testCases := []struct {
		text     string
		expected int
		desc     string
	}{
		{"", 0, "空字串"},
		{"Hello", 0, "純英文"},
		{"你好", 2, "兩個中文字符"},
		{"Hello你好World", 2, "中英混合"},
		{"你好世界！", 5, "中文字符加標點"},
		{"123你好456", 2, "數字中文混合"},
		{"你好\n世界", 4, "包含換行符"},
	}
	
	for _, tc := range testCases {
		result := enhancer.countChineseCharacters(tc.text)
		if result != tc.expected {
			t.Errorf("%s: 預期統計 %d 個中文字符，實際統計 %d 個", tc.desc, tc.expected, result)
		}
	}
}

// TestTextOperations 測試文字操作功能
// 驗證文字設定、取得和佔位符功能
func TestTextOperations(t *testing.T) {
	enhancer := NewChineseInputEnhancer()
	
	// 測試文字設定和取得
	testText := "測試中文輸入"
	enhancer.SetText(testText)
	
	if enhancer.GetText() != testText {
		t.Errorf("文字設定失敗，預期 '%s'，實際 '%s'", testText, enhancer.GetText())
	}
	
	// 測試佔位符設定
	placeholder := "請輸入中文內容"
	enhancer.SetPlaceHolder(placeholder)
	
	// 由於 Fyne 的限制，我們無法直接驗證佔位符是否設定成功
	// 但可以確保方法調用不會出錯
}

// TestConfigurationOptions 測試配置選項功能
// 驗證各種配置選項的設定和取得
func TestConfigurationOptions(t *testing.T) {
	enhancer := NewChineseInputEnhancer()
	
	// 測試候選字視窗顯示設定
	enhancer.SetShowCandidates(false)
	if enhancer.showCandidates {
		t.Error("候選字視窗顯示設定失敗")
	}
	
	enhancer.SetShowCandidates(true)
	if !enhancer.showCandidates {
		t.Error("候選字視窗顯示設定失敗")
	}
	
	// 測試自動完成設定
	enhancer.SetAutoComplete(false)
	if enhancer.autoComplete {
		t.Error("自動完成設定失敗")
	}
	
	enhancer.SetAutoComplete(true)
	if !enhancer.autoComplete {
		t.Error("自動完成設定失敗")
	}
	
	// 測試字型設定
	testFontName := "Heiti TC"
	enhancer.SetFontName(testFontName)
	if enhancer.fontName != testFontName {
		t.Errorf("字型名稱設定失敗，預期 '%s'，實際 '%s'", testFontName, enhancer.fontName)
	}
	
	// 測試字型大小設定
	testFontSize := float32(16.0)
	enhancer.SetFontSize(testFontSize)
	if enhancer.fontSize != testFontSize {
		t.Errorf("字型大小設定失敗，預期 %f，實際 %f", testFontSize, enhancer.fontSize)
	}
}

// TestCompositionState 測試組合輸入狀態管理
// 驗證組合輸入的狀態管理功能
func TestCompositionState(t *testing.T) {
	enhancer := NewChineseInputEnhancer()
	
	// 初始狀態應該不是組合輸入
	if enhancer.IsComposing() {
		t.Error("初始狀態不應該是組合輸入")
	}
	
	if enhancer.GetCompositionText() != "" {
		t.Error("初始組合文字應該為空")
	}
	
	// 測試清理組合狀態
	enhancer.clearComposition()
	
	if enhancer.IsComposing() {
		t.Error("清理後不應該是組合輸入狀態")
	}
	
	if enhancer.GetCompositionText() != "" {
		t.Error("清理後組合文字應該為空")
	}
}

// TestCandidateManagement 測試候選字管理功能
// 驗證候選字列表的管理和選擇功能
func TestCandidateManagement(t *testing.T) {
	enhancer := NewChineseInputEnhancer()
	
	// 初始狀態應該沒有候選字
	if len(enhancer.GetCandidateWords()) != 0 {
		t.Error("初始狀態不應該有候選字")
	}
	
	if enhancer.GetSelectedCandidate() != -1 {
		t.Error("初始狀態不應該有選中的候選字")
	}
	
	// 測試候選字設定
	testCandidates := []string{"你", "妳", "尼"}
	enhancer.candidateWords = testCandidates
	
	if len(enhancer.GetCandidateWords()) != len(testCandidates) {
		t.Errorf("候選字數量不正確，預期 %d，實際 %d", len(testCandidates), len(enhancer.GetCandidateWords()))
	}
	
	// 測試候選字選擇
	enhancer.selectedCandidate = 1
	if enhancer.GetSelectedCandidate() != 1 {
		t.Errorf("候選字選擇不正確，預期 1，實際 %d", enhancer.GetSelectedCandidate())
	}
	
	// 測試清理候選字
	enhancer.clearComposition()
	if len(enhancer.GetCandidateWords()) != 0 {
		t.Error("清理後不應該有候選字")
	}
	
	if enhancer.GetSelectedCandidate() != -1 {
		t.Error("清理後不應該有選中的候選字")
	}
}

// TestChineseInputCallbackFunctions 測試中文輸入回調函數功能
// 驗證各種回調函數的設定和觸發
func TestChineseInputCallbackFunctions(t *testing.T) {
	enhancer := NewChineseInputEnhancer()
	
	// 測試文字變更回調
	var textChangedCalled bool
	var textChangedContent string
	enhancer.SetOnTextChanged(func(text string) {
		textChangedCalled = true
		textChangedContent = text
	})
	
	// 模擬文字變更
	testText := "測試文字"
	enhancer.handleTextChanged(testText)
	
	if !textChangedCalled {
		t.Error("文字變更回調未被觸發")
	}
	
	if textChangedContent != testText {
		t.Errorf("文字變更回調內容不正確，預期 '%s'，實際 '%s'", testText, textChangedContent)
	}
	
	// 測試組合文字變更回調
	var compositionChangedCalled bool
	enhancer.SetOnCompositionChanged(func(text string) {
		compositionChangedCalled = true
	})
	
	// 模擬組合文字變更
	enhancer.updateChineseInputHelpers("你好", 2)
	
	if !compositionChangedCalled {
		t.Error("組合文字變更回調未被觸發")
	}
	
	// 測試候選字選擇回調
	var candidateSelectedCalled bool
	var candidateSelectedWord string
	enhancer.SetOnCandidateSelected(func(word string) {
		candidateSelectedCalled = true
		candidateSelectedWord = word
	})
	
	// 模擬候選字選擇
	enhancer.candidateWords = []string{"你", "妳", "尼"}
	enhancer.selectCandidate(0)
	
	if !candidateSelectedCalled {
		t.Error("候選字選擇回調未被觸發")
	}
	
	if candidateSelectedWord != "你" {
		t.Errorf("候選字選擇回調內容不正確，預期 '你'，實際 '%s'", candidateSelectedWord)
	}
}

// TestLastCharacterExtraction 測試最後字符提取功能
// 驗證從文字中提取最後一個字符的功能
func TestLastCharacterExtraction(t *testing.T) {
	enhancer := NewChineseInputEnhancer()
	
	// 測試各種情況的最後字符提取
	testCases := []struct {
		text     string
		expected rune
		desc     string
	}{
		{"", 0, "空字串"},
		{"A", 'A', "單個英文字符"},
		{"你", '你', "單個中文字符"},
		{"Hello", 'o', "英文單詞"},
		{"你好", '好', "中文詞語"},
		{"Hello你好", '好', "中英混合"},
		{"測試123", '3', "中文數字混合"},
	}
	
	for _, tc := range testCases {
		result := enhancer.getLastCharacter(tc.text)
		if result != tc.expected {
			t.Errorf("%s: 預期最後字符 '%c'，實際 '%c'", tc.desc, tc.expected, result)
		}
	}
}

// TestCandidateSelection 測試候選字選擇功能
// 驗證候選字選擇的邊界條件和錯誤處理
func TestCandidateSelection(t *testing.T) {
	enhancer := NewChineseInputEnhancer()
	
	// 設定測試候選字
	testCandidates := []string{"你", "妳", "尼"}
	enhancer.candidateWords = testCandidates
	
	// 測試有效的候選字選擇
	enhancer.selectCandidate(1)
	// 由於選擇後會清理狀態，我們主要測試不會出錯
	
	// 測試無效的候選字選擇（負數索引）
	enhancer.candidateWords = testCandidates
	enhancer.selectCandidate(-1)
	// 應該不會出錯，但也不會有任何效果
	
	// 測試無效的候選字選擇（超出範圍）
	enhancer.candidateWords = testCandidates
	enhancer.selectCandidate(10)
	// 應該不會出錯，但也不會有任何效果
	
	// 測試空候選字列表的選擇
	enhancer.candidateWords = []string{}
	enhancer.selectCandidate(0)
	// 應該不會出錯
}

// TestTextSubmission 測試文字提交功能
// 驗證文字提交時的狀態處理
func TestTextSubmission(t *testing.T) {
	enhancer := NewChineseInputEnhancer()
	
	// 設定組合狀態
	enhancer.isComposing = true
	enhancer.compositionText = "測試"
	enhancer.candidateWords = []string{"測試", "測驗"}
	
	// 測試文字提交
	enhancer.handleTextSubmitted("測試內容")
	
	// 驗證組合狀態已清理
	if enhancer.IsComposing() {
		t.Error("文字提交後組合狀態應該被清理")
	}
	
	if enhancer.GetCompositionText() != "" {
		t.Error("文字提交後組合文字應該被清理")
	}
	
	if len(enhancer.GetCandidateWords()) != 0 {
		t.Error("文字提交後候選字應該被清理")
	}
}

// TestUIComponentAccess 測試 UI 元件存取功能
// 驗證各種 UI 元件的存取方法
func TestUIComponentAccess(t *testing.T) {
	enhancer := NewChineseInputEnhancer()
	
	// 測試容器存取
	container := enhancer.GetContainer()
	if container == nil {
		t.Error("無法取得主要容器")
	}
	
	// 測試文字輸入元件存取
	textEntry := enhancer.GetTextEntry()
	if textEntry == nil {
		t.Error("無法取得文字輸入元件")
	}
	
	// 測試焦點設定
	// 由於測試環境的限制，我們主要確保方法調用不會出錯
	enhancer.Focus()
}

// BenchmarkChineseCharacterDetection 中文字符檢測的效能基準測試
// 測量中文字符檢測的效能
func BenchmarkChineseCharacterDetection(b *testing.B) {
	enhancer := NewChineseInputEnhancer()
	testText := "這是一段包含中文和English的混合文字，用於測試中文字符檢測的效能。"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		enhancer.containsChineseCharacters(testText)
	}
}

// BenchmarkChineseCharacterCounting 中文字符統計的效能基準測試
// 測量中文字符統計的效能
func BenchmarkChineseCharacterCounting(b *testing.B) {
	enhancer := NewChineseInputEnhancer()
	testText := "這是一段很長的中文文字，包含許多中文字符，用於測試統計功能的效能。重複的內容可以讓測試更準確。"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		enhancer.countChineseCharacters(testText)
	}
}

// BenchmarkTextOperations 文字操作的效能基準測試
// 測量文字設定和取得操作的效能
func BenchmarkTextOperations(b *testing.B) {
	enhancer := NewChineseInputEnhancer()
	testText := "測試文字內容"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		enhancer.SetText(testText)
		_ = enhancer.GetText()
	}
}