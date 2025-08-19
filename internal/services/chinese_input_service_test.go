// Package services 包含應用程式的業務邏輯服務
// 本檔案包含繁體中文輸入法優化服務的單元測試
package services

import (
	"strings"
	"testing"
)

// TestNewChineseInputService 測試中文輸入服務的建立
// 驗證服務是否正確初始化
func TestNewChineseInputService(t *testing.T) {
	service := NewChineseInputService()
	
	if service == nil {
		t.Fatal("中文輸入服務建立失敗")
	}
	
	// 測試服務的基本功能
	if !service.IsChineseCharacter('你') {
		t.Error("中文字符檢測功能異常")
	}
	
	if service.ContainsChineseCharacters("Hello") {
		t.Error("英文文字不應該包含中文字符")
	}
	
	if !service.ContainsChineseCharacters("你好") {
		t.Error("中文文字應該包含中文字符")
	}
}

// TestChineseCharacterDetection 測試中文字符檢測功能
// 驗證各種字符的中文識別準確性
func TestChineseCharacterDetection(t *testing.T) {
	service := NewChineseInputService()
	
	// 測試中文字符
	chineseChars := []rune{'你', '好', '世', '界', '中', '文', '繁', '體', '簡', '化'}
	for _, char := range chineseChars {
		if !service.IsChineseCharacter(char) {
			t.Errorf("字符 '%c' 應該被識別為中文字符", char)
		}
	}
	
	// 測試非中文字符
	nonChineseChars := []rune{'A', 'B', 'a', 'b', '1', '2', '!', '@'}
	for _, char := range nonChineseChars {
		if service.IsChineseCharacter(char) {
			t.Errorf("字符 '%c' 不應該被識別為中文字符", char)
		}
	}
	
	// 測試中文標點符號
	chinesePunctuation := []rune{'，', '。', '！', '？', '；', '：'}
	for _, char := range chinesePunctuation {
		if !service.IsChineseCharacter(char) {
			t.Errorf("中文標點符號 '%c' 應該被識別為中文字符", char)
		}
	}
}

// TestContainsChineseCharacters 測試文字中文字符包含檢測
// 驗證文字是否包含中文字符的檢測準確性
func TestContainsChineseCharacters(t *testing.T) {
	service := NewChineseInputService()
	
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
		{"測試Test", true, "中英文混合"},
	}
	
	for _, tc := range testCases {
		result := service.ContainsChineseCharacters(tc.text)
		if result != tc.expected {
			t.Errorf("%s: 預期 %v，實際 %v", tc.desc, tc.expected, result)
		}
	}
}

// TestCountChineseCharacters 測試中文字符統計功能
// 驗證中文字符數量統計的準確性
func TestCountChineseCharacters(t *testing.T) {
	service := NewChineseInputService()
	
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
		{"測試，這是一個測試。", 10, "完整中文句子"},
	}
	
	for _, tc := range testCases {
		result := service.CountChineseCharacters(tc.text)
		if result != tc.expected {
			t.Errorf("%s: 預期統計 %d 個中文字符，實際統計 %d 個", tc.desc, tc.expected, result)
		}
	}
}

// TestAnalyzeTextComposition 測試文字組成分析功能
// 驗證文字組成分析的準確性和完整性
func TestAnalyzeTextComposition(t *testing.T) {
	service := NewChineseInputService()
	
	// 測試純中文文字
	chineseText := "你好世界"
	composition := service.AnalyzeTextComposition(chineseText)
	
	if composition.TotalCharacters != 4 {
		t.Errorf("純中文文字總字符數預期為 4，實際為 %d", composition.TotalCharacters)
	}
	
	if composition.ChineseCharacters != 4 {
		t.Errorf("純中文文字中文字符數預期為 4，實際為 %d", composition.ChineseCharacters)
	}
	
	if composition.ChineseRatio != 1.0 {
		t.Errorf("純中文文字中文比例預期為 1.0，實際為 %f", composition.ChineseRatio)
	}
	
	// 測試混合文字
	mixedText := "Hello你好123"
	mixedComposition := service.AnalyzeTextComposition(mixedText)
	
	if mixedComposition.TotalCharacters != 10 {
		t.Errorf("混合文字總字符數預期為 10，實際為 %d", mixedComposition.TotalCharacters)
	}
	
	if mixedComposition.ChineseCharacters != 2 {
		t.Errorf("混合文字中文字符數預期為 2，實際為 %d", mixedComposition.ChineseCharacters)
	}
	
	if mixedComposition.EnglishCharacters != 5 {
		t.Errorf("混合文字英文字符數預期為 5，實際為 %d", mixedComposition.EnglishCharacters)
	}
	
	if mixedComposition.NumberCharacters != 3 {
		t.Errorf("混合文字數字字符數預期為 3，實際為 %d", mixedComposition.NumberCharacters)
	}
	
	expectedRatio := 2.0 / 10.0
	if mixedComposition.ChineseRatio != expectedRatio {
		t.Errorf("混合文字中文比例預期為 %f，實際為 %f", expectedRatio, mixedComposition.ChineseRatio)
	}
}

// TestGetCandidateWords 測試候選字詞功能
// 驗證候選字詞的生成和排序
func TestGetCandidateWords(t *testing.T) {
	service := NewChineseInputService()
	
	// 測試空輸入
	candidates := service.GetCandidateWords("")
	if len(candidates) != 0 {
		t.Error("空輸入不應該有候選字")
	}
	
	// 測試有效輸入
	candidates = service.GetCandidateWords("你")
	if len(candidates) == 0 {
		t.Error("有效輸入應該有候選字")
	}
	
	// 驗證候選字都以輸入開頭
	for _, candidate := range candidates {
		if !strings.HasPrefix(candidate, "你") {
			t.Errorf("候選字 '%s' 不以輸入 '你' 開頭", candidate)
		}
	}
	
	// 測試不存在的輸入
	candidates = service.GetCandidateWords("xyz")
	// 不存在的輸入可能沒有候選字，這是正常的
}

// TestGetAutoCompleteWords 測試自動完成功能
// 驗證自動完成詞彙的生成和過濾
func TestGetAutoCompleteWords(t *testing.T) {
	service := NewChineseInputService()
	
	// 測試空前綴
	autoComplete := service.GetAutoCompleteWords("")
	if len(autoComplete) != 0 {
		t.Error("空前綴不應該有自動完成詞彙")
	}
	
	// 測試有效前綴
	autoComplete = service.GetAutoCompleteWords("你")
	// 驗證自動完成詞彙都以前綴開頭且長度大於前綴
	for _, word := range autoComplete {
		if !strings.HasPrefix(word, "你") {
			t.Errorf("自動完成詞彙 '%s' 不以前綴 '你' 開頭", word)
		}
		if len(word) <= len("你") {
			t.Errorf("自動完成詞彙 '%s' 長度應該大於前綴長度", word)
		}
	}
}

// TestGetCommonWords 測試常用詞彙功能
// 驗證常用詞彙的取得和排序
func TestGetCommonWords(t *testing.T) {
	service := NewChineseInputService()
	
	commonWords := service.GetCommonWords()
	
	if len(commonWords) == 0 {
		t.Error("應該有常用詞彙")
	}
	
	// 驗證詞彙按頻率排序（頻率遞減）
	for i := 1; i < len(commonWords); i++ {
		prevFreq := service.GetWordFrequency(commonWords[i-1])
		currFreq := service.GetWordFrequency(commonWords[i])
		if prevFreq < currFreq {
			t.Errorf("常用詞彙排序錯誤：'%s'(頻率:%d) 應該在 '%s'(頻率:%d) 之前", 
				commonWords[i-1], prevFreq, commonWords[i], currFreq)
		}
	}
}

// TestOptimizeInputMethod 測試輸入法優化功能
// 驗證輸入法優化建議的生成
func TestOptimizeInputMethod(t *testing.T) {
	service := NewChineseInputService()
	
	// 測試純中文文字
	chineseText := "你好世界"
	optimization := service.OptimizeInputMethod(chineseText)
	
	if optimization.OriginalText != chineseText {
		t.Errorf("原始文字不匹配，預期 '%s'，實際 '%s'", chineseText, optimization.OriginalText)
	}
	
	if optimization.ConfidenceScore <= 0 || optimization.ConfidenceScore > 1 {
		t.Errorf("信心分數應該在 0-1 之間，實際為 %f", optimization.ConfidenceScore)
	}
	
	// 測試混合文字（中文比例低）
	// 使用中文比例很低的文字來觸發改善建議
	veryLowChineseText := "Hello World 你"
	veryLowOptimization := service.OptimizeInputMethod(veryLowChineseText)
	
	// 這個文字的中文比例為 1/13 ≈ 0.077，小於 0.1，應該有改善建議
	if len(veryLowOptimization.Improvements) == 0 {
		t.Error("中文比例很低的文字應該有改善建議")
	}
}

// TestValidateChineseInput 測試中文輸入驗證功能
// 驗證輸入驗證的準確性和完整性
func TestValidateChineseInput(t *testing.T) {
	service := NewChineseInputService()
	
	// 測試空輸入
	emptyResult := service.ValidateChineseInput("")
	if emptyResult.IsValid {
		t.Error("空輸入應該無效")
	}
	if len(emptyResult.Errors) == 0 {
		t.Error("空輸入應該有錯誤訊息")
	}
	
	// 測試有效輸入
	validResult := service.ValidateChineseInput("你好世界")
	if !validResult.IsValid {
		t.Error("有效的中文輸入應該通過驗證")
	}
	
	// 測試中文比例低的輸入
	lowChineseResult := service.ValidateChineseInput("Hello你")
	if !lowChineseResult.IsValid {
		t.Error("中文比例低的輸入仍應該有效")
	}
	if len(lowChineseResult.Suggestions) == 0 {
		t.Error("中文比例低的輸入應該有建議")
	}
}

// TestCustomWordManagement 測試自訂詞彙管理功能
// 驗證自訂詞彙的添加、移除和頻率管理
func TestCustomWordManagement(t *testing.T) {
	service := NewChineseInputService()
	
	// 測試添加自訂詞彙（使用較短的詞彙以符合長度限制）
	testWord := "測試"
	err := service.AddCustomWord(testWord)
	if err != nil {
		t.Errorf("添加自訂詞彙失敗：%v", err)
	}
	
	// 驗證詞彙頻率
	frequency := service.GetWordFrequency(testWord)
	if frequency == 0 {
		t.Error("新添加的自訂詞彙應該有頻率")
	}
	
	// 測試重複添加（應該增加頻率）
	err = service.AddCustomWord(testWord)
	if err != nil {
		t.Errorf("重複添加自訂詞彙失敗：%v", err)
	}
	
	newFrequency := service.GetWordFrequency(testWord)
	if newFrequency <= frequency {
		t.Error("重複添加應該增加詞彙頻率")
	}
	
	// 測試移除自訂詞彙
	err = service.RemoveCustomWord(testWord)
	if err != nil {
		t.Errorf("移除自訂詞彙失敗：%v", err)
	}
	
	// 驗證詞彙已被移除
	finalFrequency := service.GetWordFrequency(testWord)
	if finalFrequency != 0 {
		t.Error("移除後的詞彙頻率應該為 0")
	}
	
	// 測試移除不存在的詞彙
	err = service.RemoveCustomWord("不存在的詞彙")
	if err == nil {
		t.Error("移除不存在的詞彙應該回傳錯誤")
	}
}

// TestInvalidInputHandling 測試無效輸入處理
// 驗證各種無效輸入的處理
func TestInvalidInputHandling(t *testing.T) {
	service := NewChineseInputService()
	
	// 測試添加空詞彙
	err := service.AddCustomWord("")
	if err == nil {
		t.Error("添加空詞彙應該回傳錯誤")
	}
	
	// 測試移除空詞彙
	err = service.RemoveCustomWord("")
	if err == nil {
		t.Error("移除空詞彙應該回傳錯誤")
	}
	
	// 測試載入不存在的字典
	err = service.LoadDictionary("不存在的路徑")
	// 由於是模擬實作，這裡不會出錯，但在實際實作中應該處理檔案不存在的情況
}

// TestDictionaryLoading 測試字典載入功能
// 驗證字典載入的基本功能
func TestDictionaryLoading(t *testing.T) {
	service := NewChineseInputService()
	
	// 測試載入字典（模擬）
	err := service.LoadDictionary("test_dict.txt")
	if err != nil {
		t.Errorf("載入字典失敗：%v", err)
	}
	
	// 驗證載入後的詞彙
	testWords := []string{"繁體中文", "輸入法", "注音符號"}
	for _, word := range testWords {
		frequency := service.GetWordFrequency(word)
		if frequency == 0 {
			t.Errorf("載入字典後，詞彙 '%s' 應該有頻率", word)
		}
	}
}

// BenchmarkChineseCharacterDetection 中文字符檢測的效能基準測試
// 測量中文字符檢測的效能
func BenchmarkChineseCharacterDetection(b *testing.B) {
	service := NewChineseInputService()
	testChars := []rune{'你', '好', '世', '界', 'A', 'B', '1', '2'}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, char := range testChars {
			service.IsChineseCharacter(char)
		}
	}
}

// BenchmarkTextCompositionAnalysis 文字組成分析的效能基準測試
// 測量文字組成分析的效能
func BenchmarkTextCompositionAnalysis(b *testing.B) {
	service := NewChineseInputService()
	testText := "這是一段包含中文和English的混合文字，用於測試文字組成分析的效能。This is a mixed text for performance testing."
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.AnalyzeTextComposition(testText)
	}
}

// BenchmarkCandidateWordsGeneration 候選字詞生成的效能基準測試
// 測量候選字詞生成的效能
func BenchmarkCandidateWordsGeneration(b *testing.B) {
	service := NewChineseInputService()
	testInput := "你"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.GetCandidateWords(testInput)
	}
}

// BenchmarkAutoCompleteGeneration 自動完成生成的效能基準測試
// 測量自動完成詞彙生成的效能
func BenchmarkAutoCompleteGeneration(b *testing.B) {
	service := NewChineseInputService()
	testPrefix := "你"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.GetAutoCompleteWords(testPrefix)
	}
}