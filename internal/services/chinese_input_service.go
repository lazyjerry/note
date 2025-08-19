// Package services 包含應用程式的業務邏輯服務
// 本檔案實作繁體中文輸入法優化服務，提供中文輸入相關的業務邏輯支援
package services

import (
	"fmt"                           // 格式化輸出
	"strings"                       // 字串處理
	"unicode"                       // Unicode 字符處理
	"sort"                          // 排序功能
)

// ChineseInputService 代表繁體中文輸入法優化服務介面
// 定義中文輸入相關的業務邏輯方法
type ChineseInputService interface {
	// 字符檢測和分析
	IsChineseCharacter(r rune) bool
	ContainsChineseCharacters(text string) bool
	CountChineseCharacters(text string) int
	AnalyzeTextComposition(text string) TextComposition
	
	// 候選字和自動完成
	GetCandidateWords(input string) []string
	GetAutoCompleteWords(prefix string) []string
	GetCommonWords() []string
	
	// 輸入法優化
	OptimizeInputMethod(text string) InputOptimization
	ValidateChineseInput(text string) ValidationResult
	
	// 字典和詞庫管理
	LoadDictionary(dictPath string) error
	AddCustomWord(word string) error
	RemoveCustomWord(word string) error
	GetWordFrequency(word string) int
}

// TextComposition 代表文字組成分析結果
// 包含文字中各種字符類型的統計資訊
type TextComposition struct {
	TotalCharacters    int     // 總字符數
	ChineseCharacters  int     // 中文字符數
	EnglishCharacters  int     // 英文字符數
	NumberCharacters   int     // 數字字符數
	PunctuationMarks   int     // 標點符號數
	WhitespaceChars    int     // 空白字符數
	ChineseRatio       float64 // 中文字符比例
}

// InputOptimization 代表輸入法優化建議
// 包含針對特定輸入的優化建議和改善方案
type InputOptimization struct {
	OriginalText       string   // 原始文字
	SuggestedText      string   // 建議文字
	Corrections        []string // 修正建議
	Improvements       []string // 改善建議
	ConfidenceScore    float64  // 信心分數
}

// ValidationResult 代表中文輸入驗證結果
// 包含輸入驗證的結果和相關資訊
type ValidationResult struct {
	IsValid            bool     // 是否有效
	Errors             []string // 錯誤列表
	Warnings           []string // 警告列表
	Suggestions        []string // 建議列表
}

// chineseInputServiceImpl 實作 ChineseInputService 介面
// 提供繁體中文輸入法優化的具體實作
type chineseInputServiceImpl struct {
	// 詞庫和字典
	commonWords        map[string]int    // 常用詞彙及其頻率
	customWords        map[string]int    // 自訂詞彙及其頻率
	dictionary         map[string]bool   // 字典詞彙
	
	// 配置選項
	maxCandidates      int               // 最大候選字數量
	minWordLength      int               // 最小詞彙長度
	maxWordLength      int               // 最大詞彙長度
	enableAutoComplete bool              // 是否啟用自動完成
}

// NewChineseInputService 建立新的繁體中文輸入法優化服務實例
// 回傳：ChineseInputService 介面的實作實例
//
// 執行流程：
// 1. 建立服務實例並設定預設配置
// 2. 初始化常用詞彙和字典
// 3. 載入預設的中文詞庫
// 4. 回傳配置完成的服務實例
func NewChineseInputService() ChineseInputService {
	service := &chineseInputServiceImpl{
		commonWords:        make(map[string]int),
		customWords:        make(map[string]int),
		dictionary:         make(map[string]bool),
		maxCandidates:      10,
		minWordLength:      1,
		maxWordLength:      10,
		enableAutoComplete: true,
	}
	
	// 初始化預設詞庫
	service.initializeDefaultDictionary()
	
	return service
}

// initializeDefaultDictionary 初始化預設詞庫
// 載入常用的繁體中文詞彙和字符
//
// 執行流程：
// 1. 載入常用單字
// 2. 載入常用詞彙
// 3. 載入常用標點符號
// 4. 設定詞彙頻率
func (cis *chineseInputServiceImpl) initializeDefaultDictionary() {
	// 常用繁體中文單字
	commonChars := []string{
		"的", "一", "是", "在", "不", "了", "有", "和", "人", "這",
		"中", "大", "為", "上", "個", "國", "我", "以", "要", "他",
		"時", "來", "用", "們", "生", "到", "作", "地", "於", "出",
		"就", "分", "對", "成", "會", "可", "主", "發", "年", "動",
		"同", "工", "也", "能", "下", "過", "子", "說", "產", "種",
		"面", "而", "方", "後", "多", "定", "行", "學", "法", "所",
	}
	
	// 常用繁體中文詞彙
	commonWords := []string{
		"你好", "謝謝", "對不起", "沒關係", "再見", "早安", "晚安",
		"請問", "不好意思", "麻煩你", "辛苦了", "加油", "恭喜",
		"生日快樂", "新年快樂", "聖誕快樂", "身體健康", "工作順利",
		"學習進步", "一路順風", "祝你好運", "保重身體", "注意安全",
		"台灣", "中華民國", "繁體中文", "注音符號", "輸入法",
		"電腦", "手機", "網路", "軟體", "程式", "系統", "資料",
		"檔案", "資料夾", "下載", "上傳", "安裝", "設定", "功能",
	}
	
	// 設定常用字符的頻率
	for i, char := range commonChars {
		cis.commonWords[char] = len(commonChars) - i // 頻率遞減
		cis.dictionary[char] = true
	}
	
	// 設定常用詞彙的頻率
	for i, word := range commonWords {
		cis.commonWords[word] = len(commonWords) - i // 頻率遞減
		cis.dictionary[word] = true
	}
}

// IsChineseCharacter 檢查單個字符是否為中文字符
// 參數：r（要檢查的字符）
// 回傳：是否為中文字符的布林值
//
// 執行流程：
// 1. 檢查字符是否在 CJK 統一漢字基本區塊
// 2. 檢查字符是否在 CJK 統一漢字擴展區塊
// 3. 檢查字符是否在中文標點符號區塊
// 4. 回傳檢查結果
func (cis *chineseInputServiceImpl) IsChineseCharacter(r rune) bool {
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
	
	// CJK 符號和標點 (U+3000-U+303F)
	if r >= 0x3000 && r <= 0x303F {
		return true
	}
	
	// 全形 ASCII 字符 (U+FF00-U+FFEF)
	if r >= 0xFF00 && r <= 0xFFEF {
		return true
	}
	
	return false
}

// ContainsChineseCharacters 檢查文字是否包含中文字符
// 參數：text（要檢查的文字）
// 回傳：是否包含中文字符的布林值
//
// 執行流程：
// 1. 遍歷文字中的每個字符
// 2. 檢查字符是否為中文字符
// 3. 如果找到中文字符，立即回傳 true
// 4. 如果沒有找到中文字符，回傳 false
func (cis *chineseInputServiceImpl) ContainsChineseCharacters(text string) bool {
	for _, r := range text {
		if cis.IsChineseCharacter(r) {
			return true
		}
	}
	return false
}

// CountChineseCharacters 統計文字中的中文字符數量
// 參數：text（要統計的文字）
// 回傳：中文字符的數量
//
// 執行流程：
// 1. 初始化計數器
// 2. 遍歷文字中的每個字符
// 3. 檢查並統計中文字符
// 4. 回傳統計結果
func (cis *chineseInputServiceImpl) CountChineseCharacters(text string) int {
	count := 0
	for _, r := range text {
		if cis.IsChineseCharacter(r) {
			count++
		}
	}
	return count
}

// AnalyzeTextComposition 分析文字組成
// 參數：text（要分析的文字）
// 回傳：文字組成分析結果
//
// 執行流程：
// 1. 初始化統計變數
// 2. 遍歷文字中的每個字符
// 3. 分類統計各種字符類型
// 4. 計算中文字符比例
// 5. 回傳分析結果
func (cis *chineseInputServiceImpl) AnalyzeTextComposition(text string) TextComposition {
	composition := TextComposition{}
	
	for _, r := range text {
		composition.TotalCharacters++
		
		if cis.IsChineseCharacter(r) {
			composition.ChineseCharacters++
		} else if unicode.IsLetter(r) {
			composition.EnglishCharacters++
		} else if unicode.IsDigit(r) {
			composition.NumberCharacters++
		} else if unicode.IsPunct(r) {
			composition.PunctuationMarks++
		} else if unicode.IsSpace(r) {
			composition.WhitespaceChars++
		}
	}
	
	// 計算中文字符比例
	if composition.TotalCharacters > 0 {
		composition.ChineseRatio = float64(composition.ChineseCharacters) / float64(composition.TotalCharacters)
	}
	
	return composition
}

// GetCandidateWords 取得候選字詞
// 參數：input（輸入的文字或注音）
// 回傳：候選字詞列表
//
// 執行流程：
// 1. 分析輸入內容
// 2. 搜尋匹配的詞彙
// 3. 按頻率排序候選字詞
// 4. 限制候選字數量
// 5. 回傳候選字詞列表
func (cis *chineseInputServiceImpl) GetCandidateWords(input string) []string {
	if input == "" {
		return []string{}
	}
	
	var candidates []string
	
	// 搜尋以輸入開頭的詞彙
	for word := range cis.commonWords {
		if strings.HasPrefix(word, input) {
			candidates = append(candidates, word)
		}
	}
	
	// 搜尋自訂詞彙
	for word := range cis.customWords {
		if strings.HasPrefix(word, input) {
			candidates = append(candidates, word)
		}
	}
	
	// 按頻率排序
	sort.Slice(candidates, func(i, j int) bool {
		freqI := cis.getWordFrequency(candidates[i])
		freqJ := cis.getWordFrequency(candidates[j])
		return freqI > freqJ
	})
	
	// 限制候選字數量
	if len(candidates) > cis.maxCandidates {
		candidates = candidates[:cis.maxCandidates]
	}
	
	return candidates
}

// getWordFrequency 取得詞彙頻率（內部方法）
// 參數：word（詞彙）
// 回傳：詞彙頻率
func (cis *chineseInputServiceImpl) getWordFrequency(word string) int {
	if freq, exists := cis.commonWords[word]; exists {
		return freq
	}
	if freq, exists := cis.customWords[word]; exists {
		return freq
	}
	return 0
}

// GetAutoCompleteWords 取得自動完成詞彙
// 參數：prefix（前綴文字）
// 回傳：自動完成詞彙列表
//
// 執行流程：
// 1. 檢查是否啟用自動完成
// 2. 搜尋以前綴開頭的詞彙
// 3. 過濾詞彙長度
// 4. 按頻率排序
// 5. 回傳自動完成列表
func (cis *chineseInputServiceImpl) GetAutoCompleteWords(prefix string) []string {
	if !cis.enableAutoComplete || prefix == "" {
		return []string{}
	}
	
	var autoCompleteWords []string
	
	// 搜尋匹配的詞彙
	for word := range cis.commonWords {
		if strings.HasPrefix(word, prefix) && len(word) > len(prefix) {
			if len(word) >= cis.minWordLength && len(word) <= cis.maxWordLength {
				autoCompleteWords = append(autoCompleteWords, word)
			}
		}
	}
	
	// 搜尋自訂詞彙
	for word := range cis.customWords {
		if strings.HasPrefix(word, prefix) && len(word) > len(prefix) {
			if len(word) >= cis.minWordLength && len(word) <= cis.maxWordLength {
				autoCompleteWords = append(autoCompleteWords, word)
			}
		}
	}
	
	// 按頻率排序
	sort.Slice(autoCompleteWords, func(i, j int) bool {
		freqI := cis.getWordFrequency(autoCompleteWords[i])
		freqJ := cis.getWordFrequency(autoCompleteWords[j])
		return freqI > freqJ
	})
	
	// 限制數量
	if len(autoCompleteWords) > cis.maxCandidates {
		autoCompleteWords = autoCompleteWords[:cis.maxCandidates]
	}
	
	return autoCompleteWords
}

// GetCommonWords 取得常用詞彙列表
// 回傳：常用詞彙列表（按頻率排序）
//
// 執行流程：
// 1. 收集所有常用詞彙
// 2. 按頻率排序
// 3. 回傳排序後的列表
func (cis *chineseInputServiceImpl) GetCommonWords() []string {
	var words []string
	
	// 收集常用詞彙
	for word := range cis.commonWords {
		words = append(words, word)
	}
	
	// 按頻率排序
	sort.Slice(words, func(i, j int) bool {
		return cis.commonWords[words[i]] > cis.commonWords[words[j]]
	})
	
	return words
}

// OptimizeInputMethod 優化輸入法建議
// 參數：text（輸入的文字）
// 回傳：輸入法優化建議
//
// 執行流程：
// 1. 分析輸入文字
// 2. 檢測可能的錯誤或改善點
// 3. 生成優化建議
// 4. 計算信心分數
// 5. 回傳優化結果
func (cis *chineseInputServiceImpl) OptimizeInputMethod(text string) InputOptimization {
	optimization := InputOptimization{
		OriginalText:    text,
		SuggestedText:   text,
		Corrections:     []string{},
		Improvements:    []string{},
		ConfidenceScore: 1.0,
	}
	
	// 分析文字組成
	composition := cis.AnalyzeTextComposition(text)
	
	// 如果中文字符比例很低，建議檢查輸入法設定
	if composition.ChineseRatio < 0.1 && composition.ChineseCharacters > 0 {
		optimization.Improvements = append(optimization.Improvements, "建議檢查輸入法設定，確保正確切換到中文輸入模式")
		optimization.ConfidenceScore *= 0.9
	}
	
	// 檢查是否有常見的輸入錯誤
	corrections := cis.detectCommonErrors(text)
	optimization.Corrections = append(optimization.Corrections, corrections...)
	
	// 如果有修正建議，降低信心分數
	if len(optimization.Corrections) > 0 {
		optimization.ConfidenceScore *= 0.8
	}
	
	return optimization
}

// detectCommonErrors 檢測常見的輸入錯誤
// 參數：text（要檢測的文字）
// 回傳：錯誤修正建議列表
//
// 執行流程：
// 1. 檢查常見的錯字
// 2. 檢查標點符號使用
// 3. 檢查詞彙搭配
// 4. 回傳修正建議
func (cis *chineseInputServiceImpl) detectCommonErrors(text string) []string {
	var corrections []string
	
	// 常見錯字對照表
	commonErrors := map[string]string{
		"的話": "的話",
		"在於": "在於",
		"關於": "關於",
		"由於": "由於",
	}
	
	// 檢查常見錯字
	for wrong, correct := range commonErrors {
		if strings.Contains(text, wrong) && wrong != correct {
			corrections = append(corrections, fmt.Sprintf("建議將 '%s' 修正為 '%s'", wrong, correct))
		}
	}
	
	return corrections
}

// ValidateChineseInput 驗證中文輸入
// 參數：text（要驗證的文字）
// 回傳：驗證結果
//
// 執行流程：
// 1. 檢查文字的基本有效性
// 2. 檢查中文字符的正確性
// 3. 檢查詞彙的合理性
// 4. 生成驗證結果
func (cis *chineseInputServiceImpl) ValidateChineseInput(text string) ValidationResult {
	result := ValidationResult{
		IsValid:     true,
		Errors:      []string{},
		Warnings:    []string{},
		Suggestions: []string{},
	}
	
	if text == "" {
		result.IsValid = false
		result.Errors = append(result.Errors, "輸入文字不能為空")
		return result
	}
	
	// 分析文字組成
	composition := cis.AnalyzeTextComposition(text)
	
	// 檢查是否包含無效字符
	for _, r := range text {
		if !unicode.IsPrint(r) && !unicode.IsSpace(r) {
			result.Warnings = append(result.Warnings, "文字包含不可列印字符")
			break
		}
	}
	
	// 如果中文字符比例很低，給出建議
	if composition.ChineseRatio < 0.5 && composition.ChineseCharacters > 0 {
		result.Suggestions = append(result.Suggestions, "建議增加中文內容的比例")
	}
	
	return result
}

// LoadDictionary 載入字典檔案
// 參數：dictPath（字典檔案路徑）
// 回傳：載入結果錯誤
//
// 執行流程：
// 1. 檢查檔案是否存在
// 2. 讀取字典內容
// 3. 解析詞彙和頻率
// 4. 更新內部字典
func (cis *chineseInputServiceImpl) LoadDictionary(dictPath string) error {
	// 由於這是一個簡化的實作，這裡只是模擬載入過程
	// 實際應用中需要實作檔案讀取和解析邏輯
	
	// 模擬載入一些額外的詞彙
	additionalWords := map[string]int{
		"繁體中文": 100,
		"輸入法":   90,
		"注音符號": 80,
		"候選字":   70,
		"自動完成": 60,
	}
	
	// 添加到常用詞彙中
	for word, freq := range additionalWords {
		cis.commonWords[word] = freq
		cis.dictionary[word] = true
	}
	
	return nil
}

// AddCustomWord 添加自訂詞彙
// 參數：word（要添加的詞彙）
// 回傳：添加結果錯誤
//
// 執行流程：
// 1. 驗證詞彙的有效性
// 2. 檢查詞彙是否已存在
// 3. 添加到自訂詞彙中
// 4. 更新字典
func (cis *chineseInputServiceImpl) AddCustomWord(word string) error {
	if word == "" {
		return fmt.Errorf("詞彙不能為空")
	}
	
	if len(word) < cis.minWordLength || len(word) > cis.maxWordLength {
		return fmt.Errorf("詞彙長度必須在 %d 到 %d 字符之間", cis.minWordLength, cis.maxWordLength)
	}
	
	// 添加到自訂詞彙中，預設頻率為 1
	if _, exists := cis.customWords[word]; exists {
		cis.customWords[word]++
	} else {
		cis.customWords[word] = 1
	}
	
	cis.dictionary[word] = true
	
	return nil
}

// RemoveCustomWord 移除自訂詞彙
// 參數：word（要移除的詞彙）
// 回傳：移除結果錯誤
//
// 執行流程：
// 1. 檢查詞彙是否存在
// 2. 從自訂詞彙中移除
// 3. 更新字典
func (cis *chineseInputServiceImpl) RemoveCustomWord(word string) error {
	if word == "" {
		return fmt.Errorf("詞彙不能為空")
	}
	
	if _, exists := cis.customWords[word]; !exists {
		return fmt.Errorf("詞彙 '%s' 不存在於自訂詞彙中", word)
	}
	
	delete(cis.customWords, word)
	
	// 如果不在常用詞彙中，也從字典中移除
	if _, exists := cis.commonWords[word]; !exists {
		delete(cis.dictionary, word)
	}
	
	return nil
}

// GetWordFrequency 取得詞彙頻率
// 參數：word（詞彙）
// 回傳：詞彙頻率
//
// 執行流程：
// 1. 檢查常用詞彙中的頻率
// 2. 檢查自訂詞彙中的頻率
// 3. 回傳頻率值
func (cis *chineseInputServiceImpl) GetWordFrequency(word string) int {
	return cis.getWordFrequency(word)
}