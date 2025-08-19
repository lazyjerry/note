// Package services 測試智慧編輯服務的功能
// 本檔案包含 SmartEditingService 的完整測試案例
package services

import (
	"strings"
	"testing"
)

// TestNewSmartEditingService 測試智慧編輯服務的建立
// 驗證服務實例是否正確初始化
func TestNewSmartEditingService(t *testing.T) {
	// 建立智慧編輯服務實例
	service := NewSmartEditingService()
	
	// 驗證服務不為空
	if service == nil {
		t.Fatal("智慧編輯服務建立失敗：服務實例為空")
	}
	
	// 驗證服務實作了正確的介面
	_, ok := service.(SmartEditingService)
	if !ok {
		t.Fatal("智慧編輯服務建立失敗：未實作 SmartEditingService 介面")
	}
}

// TestAutoCompleteMarkdown 測試 Markdown 自動完成功能
// 驗證各種 Markdown 語法的自動完成建議
func TestAutoCompleteMarkdown(t *testing.T) {
	service := NewSmartEditingService()
	
	testCases := []struct {
		name           string
		content        string
		cursorPosition int
		expectedTypes  []string
		description    string
	}{
		{
			name:           "標題自動完成",
			content:        "#",
			cursorPosition: 1,
			expectedTypes:  []string{"header"},
			description:    "當輸入 # 時應該提供標題建議",
		},
		{
			name:           "列表自動完成",
			content:        "- ",
			cursorPosition: 2,
			expectedTypes:  []string{"list"},
			description:    "當輸入 - 時應該提供列表建議",
		},
		{
			name:           "程式碼區塊自動完成",
			content:        "```",
			cursorPosition: 3,
			expectedTypes:  []string{"codeblock"},
			description:    "當輸入 ``` 時應該提供程式語言建議",
		},
		{
			name:           "連結自動完成",
			content:        "[",
			cursorPosition: 1,
			expectedTypes:  []string{"link"},
			description:    "當輸入 [ 時應該提供連結建議",
		},
		{
			name:           "表格自動完成",
			content:        "|",
			cursorPosition: 1,
			expectedTypes:  []string{"table"},
			description:    "當輸入 | 時應該提供表格建議",
		},
		{
			name:           "數學公式自動完成",
			content:        "$",
			cursorPosition: 1,
			expectedTypes:  []string{"math"},
			description:    "當輸入 $ 時應該提供數學公式建議",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 執行自動完成
			suggestions := service.AutoCompleteMarkdown(tc.content, tc.cursorPosition)
			
			// 驗證是否有建議
			if len(suggestions) == 0 {
				t.Errorf("%s：沒有取得任何自動完成建議", tc.description)
				return
			}
			
			// 驗證建議類型
			foundExpectedType := false
			for _, suggestion := range suggestions {
				for _, expectedType := range tc.expectedTypes {
					if suggestion.Type == expectedType {
						foundExpectedType = true
						break
					}
				}
				if foundExpectedType {
					break
				}
			}
			
			if !foundExpectedType {
				t.Errorf("%s：未找到預期的建議類型 %v", tc.description, tc.expectedTypes)
			}
		})
	}
}

// TestFormatTable 測試表格格式化功能
// 驗證表格內容的正確格式化
func TestFormatTable(t *testing.T) {
	service := NewSmartEditingService()
	
	testCases := []struct {
		name        string
		input       string
		expectError bool
		description string
	}{
		{
			name: "基本表格格式化",
			input: `| 姓名 | 年齡 | 城市 |
|------|------|------|
| 張三 | 25 | 台北 |
| 李四 | 30 | 高雄 |`,
			expectError: false,
			description: "基本表格應該能正確格式化",
		},
		{
			name: "不對齊的表格",
			input: `|姓名|年齡|城市|
|-|-|-|
|張三|25|台北|
|李四|30|高雄|`,
			expectError: false,
			description: "不對齊的表格應該能格式化為對齊的表格",
		},
		{
			name:        "空表格",
			input:       "",
			expectError: true,
			description: "空表格應該回傳錯誤",
		},
		{
			name:        "只有一行的表格",
			input:       "| 姓名 | 年齡 |",
			expectError: true,
			description: "只有一行的表格應該回傳錯誤",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 執行表格格式化
			result, err := service.FormatTable(tc.input)
			
			// 檢查錯誤狀態
			if tc.expectError {
				if err == nil {
					t.Errorf("%s：預期會發生錯誤，但沒有錯誤發生", tc.description)
				}
				return
			}
			
			if err != nil {
				t.Errorf("%s：不應該發生錯誤，但發生了錯誤：%v", tc.description, err)
				return
			}
			
			// 驗證格式化結果
			if result == "" {
				t.Errorf("%s：格式化結果不應該為空", tc.description)
				return
			}
			
			// 驗證結果包含表格標記
			if !strings.Contains(result, "|") {
				t.Errorf("%s：格式化結果應該包含表格標記 |", tc.description)
			}
			
			// 驗證每行都以 | 開頭和結尾
			lines := strings.Split(strings.TrimSpace(result), "\n")
			for i, line := range lines {
				if !strings.HasPrefix(line, "|") || !strings.HasSuffix(line, "|") {
					t.Errorf("%s：第 %d 行格式不正確：%s", tc.description, i+1, line)
				}
			}
		})
	}
}

// TestInsertLink 測試連結插入功能
// 驗證各種連結格式的正確生成
func TestInsertLink(t *testing.T) {
	service := NewSmartEditingService()
	
	testCases := []struct {
		name        string
		text        string
		url         string
		expected    string
		description string
	}{
		{
			name:        "完整連結",
			text:        "Google",
			url:         "https://www.google.com",
			expected:    "[Google](https://www.google.com)",
			description: "完整的連結文字和網址應該正確格式化",
		},
		{
			name:        "沒有協議的網址",
			text:        "Google",
			url:         "www.google.com",
			expected:    "[Google](https://www.google.com)",
			description: "沒有協議的網址應該自動添加 https://",
		},
		{
			name:        "電子郵件連結",
			text:        "聯絡我們",
			url:         "contact@example.com",
			expected:    "[聯絡我們](mailto:contact@example.com)",
			description: "電子郵件地址應該自動添加 mailto: 協議",
		},
		{
			name:        "空文字",
			text:        "",
			url:         "https://www.google.com",
			expected:    "[https://www.google.com](https://www.google.com)",
			description: "空文字時應該使用網址作為文字",
		},
		{
			name:        "空網址",
			text:        "連結文字",
			url:         "",
			expected:    "[連結文字]()",
			description: "空網址時應該建立空的連結模板",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 執行連結插入
			result := service.InsertLink(tc.text, tc.url)
			
			// 驗證結果
			if result != tc.expected {
				t.Errorf("%s：預期 %s，實際得到 %s", tc.description, tc.expected, result)
			}
		})
	}
}

// TestInsertImage 測試圖片插入功能
// 驗證各種圖片格式的正確生成
func TestInsertImage(t *testing.T) {
	service := NewSmartEditingService()
	
	testCases := []struct {
		name        string
		altText     string
		imagePath   string
		expected    string
		description string
	}{
		{
			name:        "完整圖片",
			altText:     "示例圖片",
			imagePath:   "/images/example.png",
			expected:    "![示例圖片](/images/example.png)",
			description: "完整的替代文字和圖片路徑應該正確格式化",
		},
		{
			name:        "空替代文字",
			altText:     "",
			imagePath:   "/images/example.png",
			expected:    "![圖片](/images/example.png)",
			description: "空替代文字時應該使用預設文字",
		},
		{
			name:        "空圖片路徑",
			altText:     "示例圖片",
			imagePath:   "",
			expected:    "![示例圖片]()",
			description: "空圖片路徑時應該建立空的圖片模板",
		},
		{
			name:        "網路圖片",
			altText:     "網路圖片",
			imagePath:   "https://example.com/image.jpg",
			expected:    "![網路圖片](https://example.com/image.jpg)",
			description: "網路圖片路徑應該正確處理",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 執行圖片插入
			result := service.InsertImage(tc.altText, tc.imagePath)
			
			// 驗證結果
			if result != tc.expected {
				t.Errorf("%s：預期 %s，實際得到 %s", tc.description, tc.expected, result)
			}
		})
	}
}

// TestHighlightCodeBlock 測試程式碼區塊語法高亮功能
// 驗證各種程式語言的語法高亮
func TestHighlightCodeBlock(t *testing.T) {
	service := NewSmartEditingService()
	
	testCases := []struct {
		name        string
		code        string
		language    string
		description string
	}{
		{
			name:        "Go 程式碼高亮",
			code:        "func main() {\n    fmt.Println(\"Hello, World!\")\n}",
			language:    "go",
			description: "Go 程式碼應該正確高亮",
		},
		{
			name:        "JavaScript 程式碼高亮",
			code:        "function hello() {\n    console.log('Hello, World!');\n}",
			language:    "javascript",
			description: "JavaScript 程式碼應該正確高亮",
		},
		{
			name:        "Python 程式碼高亮",
			code:        "def hello():\n    print('Hello, World!')",
			language:    "python",
			description: "Python 程式碼應該正確高亮",
		},
		{
			name:        "不支援的語言",
			code:        "some code",
			language:    "unknown",
			description: "不支援的語言應該使用 text 格式",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 執行程式碼高亮
			result := service.HighlightCodeBlock(tc.code, tc.language)
			
			// 驗證結果包含 HTML 標籤
			if !strings.Contains(result, "<pre>") || !strings.Contains(result, "<code") {
				t.Errorf("%s：結果應該包含 HTML 標籤", tc.description)
			}
			
			// 驗證語言類別
			expectedClass := tc.language
			if tc.language == "unknown" {
				expectedClass = "text"
			}
			expectedClassAttr := "class=\"language-" + expectedClass + "\""
			if !strings.Contains(result, expectedClassAttr) {
				t.Errorf("%s：結果應該包含正確的語言類別 %s", tc.description, expectedClassAttr)
			}
			
			// 驗證程式碼內容被正確轉義
			if strings.Contains(tc.code, "<") && !strings.Contains(result, "&lt;") {
				t.Errorf("%s：HTML 字元應該被正確轉義", tc.description)
			}
		})
	}
}

// TestFormatMathExpression 測試數學公式格式化功能
// 驗證各種數學公式的正確格式化
func TestFormatMathExpression(t *testing.T) {
	service := NewSmartEditingService()
	
	testCases := []struct {
		name        string
		expression  string
		isInline    bool
		expected    string
		description string
	}{
		{
			name:        "行內數學公式",
			expression:  "x^2 + y^2 = z^2",
			isInline:    true,
			expected:    "$x^{2} + y^{2} = z^{2}$",
			description: "行內數學公式應該用單個 $ 包圍",
		},
		{
			name:        "獨立數學公式",
			expression:  "E = mc^2",
			isInline:    false,
			expected:    "$$\nE = mc^{2}\n$$",
			description: "獨立數學公式應該用雙 $$ 包圍並換行",
		},
		{
			name:        "空表達式行內",
			expression:  "",
			isInline:    true,
			expected:    "$公式$",
			description: "空表達式應該提供模板",
		},
		{
			name:        "空表達式獨立",
			expression:  "",
			isInline:    false,
			expected:    "$$\n公式\n$$",
			description: "空表達式應該提供獨立公式模板",
		},
		{
			name:        "包含 $ 符號的表達式",
			expression:  "$x + y$",
			isInline:    true,
			expected:    "$x + y$",
			description: "已包含 $ 符號的表達式應該正確處理",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 執行數學公式格式化
			result := service.FormatMathExpression(tc.expression, tc.isInline)
			
			// 驗證結果
			if result != tc.expected {
				t.Errorf("%s：預期 %s，實際得到 %s", tc.description, tc.expected, result)
			}
		})
	}
}

// TestGetSupportedLanguages 測試取得支援語言列表功能
// 驗證支援的程式語言列表是否正確
func TestGetSupportedLanguages(t *testing.T) {
	service := NewSmartEditingService()
	
	// 取得支援的語言列表
	languages := service.GetSupportedLanguages()
	
	// 驗證列表不為空
	if len(languages) == 0 {
		t.Fatal("支援的語言列表不應該為空")
	}
	
	// 驗證包含常見的程式語言
	expectedLanguages := []string{"go", "javascript", "python", "html", "css"}
	for _, expected := range expectedLanguages {
		found := false
		for _, lang := range languages {
			if lang == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("支援的語言列表應該包含 %s", expected)
		}
	}
}

// TestValidateMarkdownSyntax 測試 Markdown 語法驗證功能
// 驗證各種 Markdown 語法的正確性檢查
func TestValidateMarkdownSyntax(t *testing.T) {
	service := NewSmartEditingService()
	
	testCases := []struct {
		name        string
		content     string
		expectValid bool
		description string
	}{
		{
			name: "正確的 Markdown 語法",
			content: `# 標題
這是一段文字。

## 子標題
- 列表項目 1
- 列表項目 2

[連結](https://example.com)

![圖片](image.png)

` + "```go\nfunc main() {}\n```",
			expectValid: true,
			description: "正確的 Markdown 語法應該通過驗證",
		},
		{
			name:        "不完整的連結語法",
			content:     "這是一個 [不完整的連結",
			expectValid: false,
			description: "不完整的連結語法應該被檢測出來",
		},
		{
			name:        "不完整的圖片語法",
			content:     "這是一個 ![不完整的圖片",
			expectValid: false,
			description: "不完整的圖片語法應該被檢測出來",
		},
		{
			name: "未關閉的程式碼區塊",
			content: `這是一段文字。

` + "```go\nfunc main() {}",
			expectValid: false,
			description: "未關閉的程式碼區塊應該被檢測出來",
		},
		{
			name:        "不正確的表格語法",
			content:     "| 欄位1 | 欄位2",
			expectValid: false,
			description: "不正確的表格語法應該被檢測出來",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 執行語法驗證
			isValid, errors := service.ValidateMarkdownSyntax(tc.content)
			
			// 驗證結果
			if isValid != tc.expectValid {
				t.Errorf("%s：預期驗證結果為 %v，實際得到 %v", tc.description, tc.expectValid, isValid)
				if len(errors) > 0 {
					t.Logf("錯誤詳情：%v", errors)
				}
			}
			
			// 如果預期無效，應該有錯誤訊息
			if !tc.expectValid && len(errors) == 0 {
				t.Errorf("%s：預期會有錯誤訊息，但沒有收到任何錯誤", tc.description)
			}
		})
	}
}

// TestGenerateTableTemplate 測試表格模板生成功能
// 驗證不同大小表格模板的正確生成
func TestGenerateTableTemplate(t *testing.T) {
	service := NewSmartEditingService()
	
	testCases := []struct {
		name        string
		rows        int
		cols        int
		description string
	}{
		{
			name:        "3x3 表格",
			rows:        3,
			cols:        3,
			description: "3x3 表格模板應該正確生成",
		},
		{
			name:        "2x4 表格",
			rows:        2,
			cols:        4,
			description: "2x4 表格模板應該正確生成",
		},
		{
			name:        "無效參數（使用預設值）",
			rows:        0,
			cols:        0,
			description: "無效參數應該使用預設的 3x3 表格",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 執行表格模板生成
			result := service.GenerateTableTemplate(tc.rows, tc.cols)
			
			// 驗證結果不為空
			if result == "" {
				t.Errorf("%s：表格模板不應該為空", tc.description)
				return
			}
			
			// 驗證包含表格標記
			if !strings.Contains(result, "|") {
				t.Errorf("%s：表格模板應該包含表格標記 |", tc.description)
			}
			
			// 驗證包含分隔行
			if !strings.Contains(result, "-------") {
				t.Errorf("%s：表格模板應該包含分隔行", tc.description)
			}
			
			// 計算實際的行數和列數
			lines := strings.Split(strings.TrimSpace(result), "\n")
			if len(lines) < 2 {
				t.Errorf("%s：表格模板至少應該有 2 行（標題行和分隔行）", tc.description)
			}
			
			// 驗證每行的列數
			for i, line := range lines {
				colCount := strings.Count(line, "|") - 1 // 減去行首的 |
				expectedCols := tc.cols
				if tc.rows == 0 && tc.cols == 0 {
					expectedCols = 3 // 預設值
				}
				
				if colCount != expectedCols {
					t.Errorf("%s：第 %d 行的列數不正確，預期 %d，實際 %d", tc.description, i+1, expectedCols, colCount)
				}
			}
		})
	}
}

// TestFormatCodeBlock 測試程式碼區塊格式化功能
// 驗證程式碼區塊的正確格式化
func TestFormatCodeBlock(t *testing.T) {
	service := NewSmartEditingService()
	
	testCases := []struct {
		name        string
		code        string
		language    string
		expected    string
		description string
	}{
		{
			name:        "Go 程式碼區塊",
			code:        "func main() {\n    fmt.Println(\"Hello\")\n}",
			language:    "go",
			expected:    "```go\nfunc main() {\n    fmt.Println(\"Hello\")\n}\n```",
			description: "Go 程式碼區塊應該正確格式化",
		},
		{
			name:        "沒有語言的程式碼區塊",
			code:        "some code",
			language:    "",
			expected:    "```text\nsome code\n```",
			description: "沒有指定語言時應該使用 text",
		},
		{
			name:        "包含 ``` 的程式碼",
			code:        "```\nsome code\n```",
			language:    "text",
			expected:    "```text\n\nsome code\n\n```",
			description: "程式碼中的 ``` 應該被移除",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 執行程式碼區塊格式化
			result := service.FormatCodeBlock(tc.code, tc.language)
			
			// 驗證結果
			if result != tc.expected {
				t.Errorf("%s：預期 %s，實際得到 %s", tc.description, tc.expected, result)
			}
		})
	}
}

// BenchmarkAutoCompleteMarkdown 效能測試：自動完成功能
// 測試自動完成功能在大量內容下的效能
func BenchmarkAutoCompleteMarkdown(b *testing.B) {
	service := NewSmartEditingService()
	content := strings.Repeat("# 標題\n這是一段很長的內容。\n", 1000)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.AutoCompleteMarkdown(content, len(content)/2)
	}
}

// BenchmarkFormatTable 效能測試：表格格式化功能
// 測試表格格式化功能在大表格下的效能
func BenchmarkFormatTable(b *testing.B) {
	service := NewSmartEditingService()
	
	// 建立一個大表格
	var tableBuilder strings.Builder
	tableBuilder.WriteString("| 欄位1 | 欄位2 | 欄位3 | 欄位4 | 欄位5 |\n")
	tableBuilder.WriteString("|-------|-------|-------|-------|-------|\n")
	for i := 0; i < 100; i++ {
		tableBuilder.WriteString("| 資料1 | 資料2 | 資料3 | 資料4 | 資料5 |\n")
	}
	
	tableContent := tableBuilder.String()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.FormatTable(tableContent)
	}
}

// BenchmarkHighlightCodeBlock 效能測試：程式碼高亮功能
// 測試程式碼高亮功能在大程式碼下的效能
func BenchmarkHighlightCodeBlock(b *testing.B) {
	service := NewSmartEditingService()
	
	// 建立一個大的程式碼區塊
	code := strings.Repeat(`func example() {
    fmt.Println("Hello, World!")
    for i := 0; i < 100; i++ {
        fmt.Printf("Number: %d\n", i)
    }
}
`, 100)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.HighlightCodeBlock(code, "go")
	}
}