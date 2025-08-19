// Package services 實作應用程式的業務邏輯服務
// 本檔案包含 SmartEditingService 的具體實作，負責處理智慧編輯功能
package services

import (
	"fmt"
	"regexp"
	"strings"
)



// smartEditingService 實作 SmartEditingService 介面
type smartEditingService struct {
	// 常用的 Markdown 語法模式
	markdownPatterns map[string]*regexp.Regexp
	// 支援的程式語言列表
	supportedLanguages []string
}

// NewSmartEditingService 建立新的智慧編輯服務實例
// 回傳：SmartEditingService 介面實例
func NewSmartEditingService() SmartEditingService {
	service := &smartEditingService{
		markdownPatterns: make(map[string]*regexp.Regexp),
		supportedLanguages: []string{
			"go", "javascript", "typescript", "python", "java", "c", "cpp",
			"html", "css", "json", "xml", "yaml", "markdown", "bash", "sql",
		},
	}
	
	// 初始化 Markdown 語法模式
	service.initializePatterns()
	
	return service
}

// initializePatterns 初始化 Markdown 語法識別模式
// 執行流程：
// 1. 建立各種 Markdown 語法的正規表達式
// 2. 儲存到模式映射表中供後續使用
func (s *smartEditingService) initializePatterns() {
	// 標題模式（# ## ### 等）
	s.markdownPatterns["header"] = regexp.MustCompile(`^(#{1,6})\s`)
	
	// 列表模式（- * + 或數字列表）
	s.markdownPatterns["list"] = regexp.MustCompile(`^(\s*)([-*+]|\d+\.)\s`)
	
	// 程式碼區塊模式（```）
	s.markdownPatterns["codeblock"] = regexp.MustCompile("^```(\\w+)?")
	
	// 連結模式（[text](url)）
	s.markdownPatterns["link"] = regexp.MustCompile(`\[([^\]]*)\]\(([^)]*)\)`)
	
	// 圖片模式（![alt](src)）
	s.markdownPatterns["image"] = regexp.MustCompile(`!\[([^\]]*)\]\(([^)]*)\)`)
	
	// 表格模式（| col1 | col2 |）
	s.markdownPatterns["table"] = regexp.MustCompile(`^\s*\|.*\|\s*$`)
	
	// 數學公式模式（$...$ 或 $$...$$）
	s.markdownPatterns["math_inline"] = regexp.MustCompile(`\$([^$]+)\$`)
	s.markdownPatterns["math_block"] = regexp.MustCompile(`\$\$([^$]+)\$\$`)
}

// AutoCompleteMarkdown 提供 Markdown 語法自動完成建議
// 參數：content（當前內容）、cursorPosition（游標位置）
// 回傳：自動完成建議陣列
//
// 執行流程：
// 1. 分析游標位置的上下文
// 2. 識別當前輸入的模式
// 3. 根據模式提供相應的自動完成建議
// 4. 回傳建議列表
func (s *smartEditingService) AutoCompleteMarkdown(content string, cursorPosition int) []AutoCompleteSuggestion {
	var suggestions []AutoCompleteSuggestion
	
	// 確保游標位置有效
	if cursorPosition < 0 || cursorPosition > len(content) {
		return suggestions
	}
	
	// 取得游標所在行的內容
	lines := strings.Split(content[:cursorPosition], "\n")
	currentLine := ""
	if len(lines) > 0 {
		currentLine = lines[len(lines)-1]
	}
	
	// 根據當前行內容提供建議
	suggestions = append(suggestions, s.getHeaderSuggestions(currentLine)...)
	suggestions = append(suggestions, s.getListSuggestions(currentLine)...)
	suggestions = append(suggestions, s.getCodeBlockSuggestions(currentLine)...)
	suggestions = append(suggestions, s.getLinkSuggestions(currentLine)...)
	suggestions = append(suggestions, s.getTableSuggestions(currentLine)...)
	suggestions = append(suggestions, s.getMathSuggestions(currentLine)...)
	
	return suggestions
}

// getHeaderSuggestions 取得標題相關的自動完成建議
// 參數：currentLine（當前行內容）
// 回傳：標題建議陣列
func (s *smartEditingService) getHeaderSuggestions(currentLine string) []AutoCompleteSuggestion {
	var suggestions []AutoCompleteSuggestion
	
	// 如果行開頭是 #，提供標題建議
	if strings.HasPrefix(strings.TrimSpace(currentLine), "#") {
		suggestions = append(suggestions, AutoCompleteSuggestion{
			Text:        "# 主標題",
			Description: "一級標題（最大）",
			Type:        "header",
			InsertText:  "# ",
		})
		suggestions = append(suggestions, AutoCompleteSuggestion{
			Text:        "## 次標題",
			Description: "二級標題",
			Type:        "header",
			InsertText:  "## ",
		})
		suggestions = append(suggestions, AutoCompleteSuggestion{
			Text:        "### 小標題",
			Description: "三級標題",
			Type:        "header",
			InsertText:  "### ",
		})
	}
	
	return suggestions
}

// getListSuggestions 取得列表相關的自動完成建議
// 參數：currentLine（當前行內容）
// 回傳：列表建議陣列
func (s *smartEditingService) getListSuggestions(currentLine string) []AutoCompleteSuggestion {
	var suggestions []AutoCompleteSuggestion
	
	trimmed := strings.TrimSpace(currentLine)
	
	// 如果行開頭是列表符號，提供列表建議
	if strings.HasPrefix(trimmed, "-") || strings.HasPrefix(trimmed, "*") || strings.HasPrefix(trimmed, "+") {
		suggestions = append(suggestions, AutoCompleteSuggestion{
			Text:        "- 項目",
			Description: "無序列表項目",
			Type:        "list",
			InsertText:  "- ",
		})
		suggestions = append(suggestions, AutoCompleteSuggestion{
			Text:        "  - 子項目",
			Description: "縮排的子列表項目",
			Type:        "list",
			InsertText:  "  - ",
		})
	}
	
	// 數字列表建議
	if regexp.MustCompile(`^\d+\.`).MatchString(trimmed) {
		suggestions = append(suggestions, AutoCompleteSuggestion{
			Text:        "1. 項目",
			Description: "有序列表項目",
			Type:        "list",
			InsertText:  "1. ",
		})
	}
	
	return suggestions
}

// getCodeBlockSuggestions 取得程式碼區塊相關的自動完成建議
// 參數：currentLine（當前行內容）
// 回傳：程式碼區塊建議陣列
func (s *smartEditingService) getCodeBlockSuggestions(currentLine string) []AutoCompleteSuggestion {
	var suggestions []AutoCompleteSuggestion
	
	// 如果行開頭是 ```，提供程式語言建議
	if strings.HasPrefix(strings.TrimSpace(currentLine), "```") {
		for _, lang := range s.supportedLanguages {
			suggestions = append(suggestions, AutoCompleteSuggestion{
				Text:        fmt.Sprintf("```%s", lang),
				Description: fmt.Sprintf("%s 程式碼區塊", strings.ToUpper(lang)),
				Type:        "codeblock",
				InsertText:  fmt.Sprintf("```%s\n\n```", lang),
			})
		}
	}
	
	return suggestions
}

// getLinkSuggestions 取得連結相關的自動完成建議
// 參數：currentLine（當前行內容）
// 回傳：連結建議陣列
func (s *smartEditingService) getLinkSuggestions(currentLine string) []AutoCompleteSuggestion {
	var suggestions []AutoCompleteSuggestion
	
	// 如果包含 [ 但沒有完整的連結格式，提供連結建議
	if strings.Contains(currentLine, "[") && !s.markdownPatterns["link"].MatchString(currentLine) {
		suggestions = append(suggestions, AutoCompleteSuggestion{
			Text:        "[連結文字](網址)",
			Description: "插入連結",
			Type:        "link",
			InsertText:  "[]()",
		})
		suggestions = append(suggestions, AutoCompleteSuggestion{
			Text:        "[連結文字](網址 \"標題\")",
			Description: "插入帶標題的連結",
			Type:        "link",
			InsertText:  "[](\"\")",
		})
	}
	
	return suggestions
}

// getTableSuggestions 取得表格相關的自動完成建議
// 參數：currentLine（當前行內容）
// 回傳：表格建議陣列
func (s *smartEditingService) getTableSuggestions(currentLine string) []AutoCompleteSuggestion {
	var suggestions []AutoCompleteSuggestion
	
	// 如果行包含 |，提供表格建議
	if strings.Contains(currentLine, "|") {
		suggestions = append(suggestions, AutoCompleteSuggestion{
			Text:        "| 欄位1 | 欄位2 | 欄位3 |",
			Description: "表格標題行",
			Type:        "table",
			InsertText:  "| 欄位1 | 欄位2 | 欄位3 |\n|-------|-------|-------|\n| 內容1 | 內容2 | 內容3 |",
		})
		suggestions = append(suggestions, AutoCompleteSuggestion{
			Text:        "|-------|-------|",
			Description: "表格分隔線",
			Type:        "table",
			InsertText:  "|-------|-------|",
		})
	}
	
	return suggestions
}

// getMathSuggestions 取得數學公式相關的自動完成建議
// 參數：currentLine（當前行內容）
// 回傳：數學公式建議陣列
func (s *smartEditingService) getMathSuggestions(currentLine string) []AutoCompleteSuggestion {
	var suggestions []AutoCompleteSuggestion
	
	// 如果包含 $，提供數學公式建議
	if strings.Contains(currentLine, "$") {
		suggestions = append(suggestions, AutoCompleteSuggestion{
			Text:        "$公式$",
			Description: "行內數學公式",
			Type:        "math",
			InsertText:  "$  $",
		})
		suggestions = append(suggestions, AutoCompleteSuggestion{
			Text:        "$$公式$$",
			Description: "獨立數學公式",
			Type:        "math",
			InsertText:  "$$\n  \n$$",
		})
	}
	
	return suggestions
}

// FormatTable 格式化表格內容
// 參數：tableContent（表格內容）
// 回傳：格式化後的表格字串和可能的錯誤
//
// 執行流程：
// 1. 解析表格內容，分離標題和資料行
// 2. 計算每欄的最大寬度
// 3. 重新格式化表格，確保對齊
// 4. 回傳格式化後的表格
func (s *smartEditingService) FormatTable(tableContent string) (string, error) {
	lines := strings.Split(strings.TrimSpace(tableContent), "\n")
	if len(lines) < 2 {
		return "", fmt.Errorf("表格至少需要標題行和分隔行")
	}
	
	var rows [][]string
	var maxWidths []int
	
	// 解析每一行
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		// 移除行首尾的 |
		line = strings.Trim(line, "|")
		
		// 分割欄位
		cells := strings.Split(line, "|")
		for i, cell := range cells {
			cells[i] = strings.TrimSpace(cell)
		}
		
		rows = append(rows, cells)
		
		// 更新最大寬度
		for i, cell := range cells {
			if i >= len(maxWidths) {
				maxWidths = append(maxWidths, 0)
			}
			if len(cell) > maxWidths[i] {
				maxWidths[i] = len(cell)
			}
		}
	}
	
	// 格式化表格
	var result strings.Builder
	
	for rowIndex, row := range rows {
		result.WriteString("|")
		
		for colIndex, cell := range row {
			width := maxWidths[colIndex]
			if width < 3 {
				width = 3 // 最小寬度
			}
			
			// 如果是分隔行（通常是第二行），使用 - 填充
			if rowIndex == 1 && (cell == "" || strings.Contains(cell, "-")) {
				result.WriteString(fmt.Sprintf(" %s |", strings.Repeat("-", width)))
			} else {
				result.WriteString(fmt.Sprintf(" %-*s |", width, cell))
			}
		}
		
		result.WriteString("\n")
	}
	
	return result.String(), nil
}

// InsertLink 插入連結
// 參數：text（連結文字）、url（連結網址）
// 回傳：格式化的 Markdown 連結字串
//
// 執行流程：
// 1. 驗證輸入參數
// 2. 清理連結文字和網址
// 3. 格式化為 Markdown 連結格式
// 4. 回傳格式化的連結字串
func (s *smartEditingService) InsertLink(text, url string) string {
	// 清理輸入
	text = strings.TrimSpace(text)
	url = strings.TrimSpace(url)
	
	// 如果沒有提供文字，使用網址作為文字
	if text == "" {
		text = url
	}
	
	// 如果沒有提供網址，建立空的連結模板
	if url == "" {
		return fmt.Sprintf("[%s]()", text)
	}
	
	// 檢查網址是否需要添加協議
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") && !strings.HasPrefix(url, "mailto:") {
		// 如果看起來像電子郵件地址
		if strings.Contains(url, "@") && strings.Contains(url, ".") {
			url = "mailto:" + url
		} else {
			url = "https://" + url
		}
	}
	
	return fmt.Sprintf("[%s](%s)", text, url)
}

// InsertImage 插入圖片
// 參數：altText（替代文字）、imagePath（圖片路徑）
// 回傳：格式化的 Markdown 圖片字串
//
// 執行流程：
// 1. 驗證輸入參數
// 2. 清理替代文字和圖片路徑
// 3. 格式化為 Markdown 圖片格式
// 4. 回傳格式化的圖片字串
func (s *smartEditingService) InsertImage(altText, imagePath string) string {
	// 清理輸入
	altText = strings.TrimSpace(altText)
	imagePath = strings.TrimSpace(imagePath)
	
	// 如果沒有提供替代文字，使用預設文字
	if altText == "" {
		altText = "圖片"
	}
	
	// 如果沒有提供路徑，建立空的圖片模板
	if imagePath == "" {
		return fmt.Sprintf("![%s]()", altText)
	}
	
	return fmt.Sprintf("![%s](%s)", altText, imagePath)
}

// HighlightCodeBlock 為程式碼區塊添加語法高亮
// 參數：code（程式碼內容）、language（程式語言）
// 回傳：帶有語法高亮的 HTML 字串
//
// 執行流程：
// 1. 驗證程式語言是否支援
// 2. 清理程式碼內容
// 3. 應用基本的語法高亮規則
// 4. 回傳 HTML 格式的高亮程式碼
func (s *smartEditingService) HighlightCodeBlock(code, language string) string {
	// 清理輸入
	code = strings.TrimSpace(code)
	language = strings.ToLower(strings.TrimSpace(language))
	
	// 檢查是否支援該語言
	supported := false
	for _, lang := range s.supportedLanguages {
		if lang == language {
			supported = true
			break
		}
	}
	
	if !supported {
		language = "text"
	}
	
	// 基本的 HTML 轉義
	code = strings.ReplaceAll(code, "&", "&amp;")
	code = strings.ReplaceAll(code, "<", "&lt;")
	code = strings.ReplaceAll(code, ">", "&gt;")
	code = strings.ReplaceAll(code, "\"", "&quot;")
	
	// 應用基本的語法高亮（簡化版本）
	highlightedCode := s.applyBasicHighlighting(code, language)
	
	return fmt.Sprintf(`<pre><code class="language-%s">%s</code></pre>`, language, highlightedCode)
}

// applyBasicHighlighting 應用基本的語法高亮
// 參數：code（程式碼）、language（程式語言）
// 回傳：高亮後的 HTML 程式碼
func (s *smartEditingService) applyBasicHighlighting(code, language string) string {
	switch language {
	case "go":
		return s.highlightGo(code)
	case "javascript", "js":
		return s.highlightJavaScript(code)
	case "python":
		return s.highlightPython(code)
	case "html":
		return s.highlightHTML(code)
	case "css":
		return s.highlightCSS(code)
	default:
		return code
	}
}

// highlightGo 高亮 Go 語言程式碼
func (s *smartEditingService) highlightGo(code string) string {
	result := code
	
	// 簡化的語法高亮，避免嵌套問題
	// 只高亮關鍵字，不處理字串和註解以避免複雜的衝突
	keywords := []string{
		"package", "import", "func", "var", "const", "type", "struct", "interface",
		"if", "else", "for", "range", "switch", "case", "default", "break", "continue",
		"return", "go", "defer", "chan", "select", "map", "make", "new", "len", "cap",
	}
	
	for _, keyword := range keywords {
		pattern := regexp.MustCompile(`\b` + keyword + `\b`)
		result = pattern.ReplaceAllString(result, `<span class="keyword">`+keyword+`</span>`)
	}
	
	return result
}

// highlightJavaScript 高亮 JavaScript 程式碼
func (s *smartEditingService) highlightJavaScript(code string) string {
	keywords := []string{
		"function", "var", "let", "const", "if", "else", "for", "while", "do",
		"switch", "case", "default", "break", "continue", "return", "try", "catch",
		"finally", "throw", "new", "this", "typeof", "instanceof", "in", "of",
	}
	
	result := code
	for _, keyword := range keywords {
		pattern := regexp.MustCompile(`\b` + keyword + `\b`)
		result = pattern.ReplaceAllString(result, `<span class="keyword">`+keyword+`</span>`)
	}
	
	return result
}

// highlightPython 高亮 Python 程式碼
func (s *smartEditingService) highlightPython(code string) string {
	keywords := []string{
		"def", "class", "if", "elif", "else", "for", "while", "try", "except",
		"finally", "with", "as", "import", "from", "return", "yield", "lambda",
		"and", "or", "not", "in", "is", "True", "False", "None",
	}
	
	result := code
	for _, keyword := range keywords {
		pattern := regexp.MustCompile(`\b` + keyword + `\b`)
		result = pattern.ReplaceAllString(result, `<span class="keyword">`+keyword+`</span>`)
	}
	
	return result
}

// highlightHTML 高亮 HTML 程式碼
func (s *smartEditingService) highlightHTML(code string) string {
	// HTML 標籤高亮
	tagPattern := regexp.MustCompile(`&lt;/?[^&gt;]+&gt;`)
	result := tagPattern.ReplaceAllString(code, `<span class="tag">$0</span>`)
	
	return result
}

// highlightCSS 高亮 CSS 程式碼
func (s *smartEditingService) highlightCSS(code string) string {
	// CSS 選擇器高亮
	selectorPattern := regexp.MustCompile(`^[^{]+(?=\s*{)`)
	result := selectorPattern.ReplaceAllString(code, `<span class="selector">$0</span>`)
	
	// CSS 屬性高亮
	propertyPattern := regexp.MustCompile(`([a-zA-Z-]+)\s*:`)
	result = propertyPattern.ReplaceAllString(result, `<span class="property">$1</span>:`)
	
	return result
}

// FormatMathExpression 格式化數學公式
// 參數：expression（數學表達式）、isInline（是否為行內公式）
// 回傳：格式化的 LaTeX 數學公式字串
//
// 執行流程：
// 1. 清理數學表達式
// 2. 驗證基本的 LaTeX 語法
// 3. 根據是否為行內公式選擇格式
// 4. 回傳格式化的數學公式
func (s *smartEditingService) FormatMathExpression(expression string, isInline bool) string {
	// 清理輸入
	expression = strings.TrimSpace(expression)
	
	// 如果表達式為空，提供模板
	if expression == "" {
		if isInline {
			return "$公式$"
		}
		return "$$\n公式\n$$"
	}
	
	// 移除現有的 $ 符號（如果有）
	expression = strings.Trim(expression, "$")
	expression = strings.TrimSpace(expression)
	
	// 基本的 LaTeX 語法檢查和修正
	expression = s.fixCommonMathSyntax(expression)
	
	// 根據是否為行內公式格式化
	if isInline {
		return fmt.Sprintf("$%s$", expression)
	}
	
	// 獨立公式（區塊公式）
	return fmt.Sprintf("$$\n%s\n$$", expression)
}

// fixCommonMathSyntax 修正常見的數學語法錯誤
// 參數：expression（數學表達式）
// 回傳：修正後的表達式
func (s *smartEditingService) fixCommonMathSyntax(expression string) string {
	// 常見的數學符號替換
	replacements := map[string]string{
		"*":     "\\cdot ",     // 乘號
		"+-":    "\\pm ",       // 正負號
		"-+":    "\\mp ",       // 負正號
		"<=":    "\\leq ",      // 小於等於
		">=":    "\\geq ",      // 大於等於
		"!=":    "\\neq ",      // 不等於
		"~=":    "\\approx ",   // 約等於
		"alpha": "\\alpha ",    // 希臘字母
		"beta":  "\\beta ",
		"gamma": "\\gamma ",
		"delta": "\\delta ",
		"pi":    "\\pi ",
		"theta": "\\theta ",
		"lambda": "\\lambda ",
		"mu":    "\\mu ",
		"sigma": "\\sigma ",
		"phi":   "\\phi ",
		"omega": "\\omega ",
	}
	
	result := expression
	for old, new := range replacements {
		result = strings.ReplaceAll(result, old, new)
	}
	
	// 處理分數（簡單的 a/b 格式）
	fractionPattern := regexp.MustCompile(`(\w+)/(\w+)`)
	result = fractionPattern.ReplaceAllString(result, `\frac{$1}{$2}`)
	
	// 處理上標（x^2 格式）
	superscriptPattern := regexp.MustCompile(`(\w+)\^(\w+)`)
	result = superscriptPattern.ReplaceAllString(result, `$1^{$2}`)
	
	// 處理下標（x_1 格式）
	subscriptPattern := regexp.MustCompile(`(\w+)_(\w+)`)
	result = subscriptPattern.ReplaceAllString(result, `$1_{$2}`)
	
	return result
}

// GetSupportedLanguages 取得支援的程式語言列表
// 回傳：支援的程式語言陣列
func (s *smartEditingService) GetSupportedLanguages() []string {
	return s.supportedLanguages
}

// ValidateMarkdownSyntax 驗證 Markdown 語法的正確性
// 參數：content（要驗證的 Markdown 內容）
// 回傳：驗證結果和可能的錯誤列表
func (s *smartEditingService) ValidateMarkdownSyntax(content string) (bool, []string) {
	var errors []string
	lines := strings.Split(content, "\n")
	
	inCodeBlock := false
	codeBlockStart := -1
	
	for i, line := range lines {
		lineNum := i + 1
		
		// 檢查程式碼區塊
		if strings.HasPrefix(strings.TrimSpace(line), "```") {
			if inCodeBlock {
				inCodeBlock = false
			} else {
				inCodeBlock = true
				codeBlockStart = lineNum
			}
			continue
		}
		
		// 如果在程式碼區塊內，跳過語法檢查
		if inCodeBlock {
			continue
		}
		
		// 檢查連結語法
		if strings.Contains(line, "[") || strings.Contains(line, "]") {
			if !s.validateLinkSyntax(line) {
				errors = append(errors, fmt.Sprintf("第 %d 行：連結語法不完整", lineNum))
			}
		}
		
		// 檢查圖片語法
		if strings.Contains(line, "![") {
			if !s.validateImageSyntax(line) {
				errors = append(errors, fmt.Sprintf("第 %d 行：圖片語法不完整", lineNum))
			}
		}
		
		// 檢查表格語法
		if strings.Contains(line, "|") {
			if !s.validateTableSyntax(line) {
				errors = append(errors, fmt.Sprintf("第 %d 行：表格語法不正確", lineNum))
			}
		}
	}
	
	// 檢查是否有未關閉的程式碼區塊
	if inCodeBlock {
		errors = append(errors, fmt.Sprintf("第 %d 行：程式碼區塊未正確關閉", codeBlockStart))
	}
	
	return len(errors) == 0, errors
}

// validateLinkSyntax 驗證連結語法
func (s *smartEditingService) validateLinkSyntax(line string) bool {
	// 簡單的連結語法檢查
	openBrackets := strings.Count(line, "[")
	closeBrackets := strings.Count(line, "]")
	openParens := strings.Count(line, "(")
	closeParens := strings.Count(line, ")")
	
	return openBrackets == closeBrackets && openParens == closeParens
}

// validateImageSyntax 驗證圖片語法
func (s *smartEditingService) validateImageSyntax(line string) bool {
	// 檢查圖片語法 ![alt](src)
	return s.markdownPatterns["image"].MatchString(line) || !strings.Contains(line, "![")
}

// validateTableSyntax 驗證表格語法
func (s *smartEditingService) validateTableSyntax(line string) bool {
	// 簡單的表格語法檢查：確保 | 符號成對出現
	trimmed := strings.TrimSpace(line)
	if !strings.HasPrefix(trimmed, "|") || !strings.HasSuffix(trimmed, "|") {
		return false
	}
	
	// 檢查是否至少有兩個 | 符號
	return strings.Count(line, "|") >= 2
}

// GenerateTableTemplate 生成表格模板
// 參數：rows（行數）、cols（列數）
// 回傳：表格模板字串
func (s *smartEditingService) GenerateTableTemplate(rows, cols int) string {
	if rows < 2 || cols < 1 {
		rows, cols = 3, 3 // 預設 3x3 表格
	}
	
	var result strings.Builder
	
	// 生成標題行
	result.WriteString("|")
	for i := 0; i < cols; i++ {
		result.WriteString(fmt.Sprintf(" 欄位%d |", i+1))
	}
	result.WriteString("\n")
	
	// 生成分隔行
	result.WriteString("|")
	for i := 0; i < cols; i++ {
		result.WriteString("-------|")
	}
	result.WriteString("\n")
	
	// 生成資料行
	for r := 0; r < rows-1; r++ {
		result.WriteString("|")
		for c := 0; c < cols; c++ {
			result.WriteString(fmt.Sprintf(" 內容%d-%d |", r+1, c+1))
		}
		result.WriteString("\n")
	}
	
	return result.String()
}

// FormatCodeBlock 格式化程式碼區塊
// 參數：code（程式碼內容）、language（程式語言）
// 回傳：格式化的 Markdown 程式碼區塊
func (s *smartEditingService) FormatCodeBlock(code, language string) string {
	// 清理輸入
	code = strings.TrimSpace(code)
	language = strings.TrimSpace(language)
	
	// 如果沒有指定語言，使用 text
	if language == "" {
		language = "text"
	}
	
	// 確保程式碼不包含 ``` 標記
	code = strings.ReplaceAll(code, "```", "")
	
	return fmt.Sprintf("```%s\n%s\n```", language, code)
}