// Package models 的測試檔案
// 包含 Note 資料模型的完整單元測試
package models

import (
	"strings"
	"testing"
	"time"
)

// TestNewNote 測試 NewNote 函數的基本功能
// 驗證新建立的筆記是否具有正確的預設值和屬性
func TestNewNote(t *testing.T) {
	// 測試資料
	title := "測試筆記標題"
	content := "這是測試筆記的內容"
	filePath := "/path/to/test.md"

	// 執行測試
	note := NewNote(title, content, filePath)

	// 驗證基本屬性
	if note.Title != title {
		t.Errorf("期望標題為 '%s'，實際得到 '%s'", title, note.Title)
	}
	if note.Content != content {
		t.Errorf("期望內容為 '%s'，實際得到 '%s'", content, note.Content)
	}
	if note.FilePath != filePath {
		t.Errorf("期望檔案路徑為 '%s'，實際得到 '%s'", filePath, note.FilePath)
	}

	// 驗證預設值
	if note.ID == "" {
		t.Error("筆記 ID 不應該為空")
	}
	if note.IsEncrypted {
		t.Error("新筆記預設應該是未加密狀態")
	}
	if note.EncryptionType != "" {
		t.Error("新筆記的加密類型應該為空字串")
	}

	// 驗證時間戳
	if note.CreatedAt.IsZero() {
		t.Error("建立時間不應該為零值")
	}
	if note.UpdatedAt.IsZero() {
		t.Error("修改時間不應該為零值")
	}
	if !note.LastSaved.IsZero() {
		t.Error("新筆記的最後保存時間應該為零值")
	}

	// 驗證 ID 格式（應該包含時間戳和隨機字串）
	if len(note.ID) < 15 { // 時間戳14位 + 連字符1位 + 至少8位隨機字串
		t.Error("筆記 ID 格式不正確")
	}
	if !strings.Contains(note.ID, "-") {
		t.Error("筆記 ID 應該包含連字符")
	}
}

// TestNote_UpdateContent 測試筆記內容更新功能
// 驗證內容更新後修改時間戳是否正確更新
func TestNote_UpdateContent(t *testing.T) {
	// 建立測試筆記
	note := NewNote("測試標題", "原始內容", "/test.md")
	originalUpdateTime := note.UpdatedAt

	// 等待一小段時間確保時間戳不同
	time.Sleep(time.Millisecond * 10)

	// 更新內容
	newContent := "更新後的內容"
	note.UpdateContent(newContent)

	// 驗證內容已更新
	if note.Content != newContent {
		t.Errorf("期望內容為 '%s'，實際得到 '%s'", newContent, note.Content)
	}

	// 驗證修改時間已更新
	if !note.UpdatedAt.After(originalUpdateTime) {
		t.Error("修改時間應該在內容更新後被更新")
	}
}

// TestNote_MarkSaved 測試筆記保存標記功能
// 驗證保存標記後最後保存時間是否正確設定
func TestNote_MarkSaved(t *testing.T) {
	// 建立測試筆記
	note := NewNote("測試標題", "測試內容", "/test.md")

	// 驗證初始狀態
	if !note.LastSaved.IsZero() {
		t.Error("新筆記的最後保存時間應該為零值")
	}

	// 標記為已保存
	note.MarkSaved()

	// 驗證最後保存時間已設定
	if note.LastSaved.IsZero() {
		t.Error("標記保存後，最後保存時間不應該為零值")
	}

	// 驗證最後保存時間是最近的時間
	timeDiff := time.Since(note.LastSaved)
	if timeDiff > time.Second {
		t.Error("最後保存時間應該是最近的時間")
	}
}

// TestNote_IsModified 測試筆記修改狀態檢查功能
// 驗證各種情況下修改狀態的判斷是否正確
func TestNote_IsModified(t *testing.T) {
	// 建立測試筆記
	note := NewNote("測試標題", "測試內容", "/test.md")

	// 新筆記應該被視為已修改（因為從未保存）
	if !note.IsModified() {
		t.Error("新筆記應該被視為已修改")
	}

	// 標記為已保存後應該不被視為已修改
	note.MarkSaved()
	if note.IsModified() {
		t.Error("剛保存的筆記不應該被視為已修改")
	}

	// 更新內容後應該被視為已修改
	time.Sleep(time.Millisecond * 10) // 確保時間戳不同
	note.UpdateContent("新的內容")
	if !note.IsModified() {
		t.Error("更新內容後的筆記應該被視為已修改")
	}
}

// TestNote_Validate 測試筆記資料驗證功能
// 驗證各種無效資料情況下的驗證錯誤
func TestNote_Validate(t *testing.T) {
	// 測試有效的筆記
	validNote := NewNote("有效標題", "有效內容", "/valid/path.md")
	if err := validNote.Validate(); err != nil {
		t.Errorf("有效筆記不應該產生驗證錯誤：%v", err)
	}

	// 測試空 ID
	invalidNote := NewNote("標題", "內容", "/path.md")
	invalidNote.ID = ""
	if err := invalidNote.Validate(); err == nil {
		t.Error("空 ID 應該產生驗證錯誤")
	}

	// 測試空標題
	invalidNote = NewNote("", "內容", "/path.md")
	if err := invalidNote.Validate(); err == nil {
		t.Error("空標題應該產生驗證錯誤")
	}

	// 測試過長標題
	longTitle := strings.Repeat("長", 201) // 201個字符
	invalidNote = NewNote(longTitle, "內容", "/path.md")
	if err := invalidNote.Validate(); err == nil {
		t.Error("過長標題應該產生驗證錯誤")
	}

	// 測試空檔案路徑
	invalidNote = NewNote("標題", "內容", "")
	if err := invalidNote.Validate(); err == nil {
		t.Error("空檔案路徑應該產生驗證錯誤")
	}

	// 測試無效的加密類型
	invalidNote = NewNote("標題", "內容", "/path.md")
	invalidNote.IsEncrypted = true
	invalidNote.EncryptionType = "invalid_type"
	if err := invalidNote.Validate(); err == nil {
		t.Error("無效的加密類型應該產生驗證錯誤")
	}

	// 測試零值時間戳
	invalidNote = NewNote("標題", "內容", "/path.md")
	invalidNote.CreatedAt = time.Time{}
	if err := invalidNote.Validate(); err == nil {
		t.Error("零值建立時間應該產生驗證錯誤")
	}
}

// TestNote_SetEncryption 測試筆記加密設定功能
// 驗證各種加密類型的設定是否正確
func TestNote_SetEncryption(t *testing.T) {
	note := NewNote("測試標題", "測試內容", "/test.md")
	originalUpdateTime := note.UpdatedAt

	// 等待一小段時間確保時間戳不同
	time.Sleep(time.Millisecond * 10)

	// 測試有效的加密類型
	validTypes := []string{"password", "biometric", "both"}
	for _, encType := range validTypes {
		testNote := NewNote("標題", "內容", "/path.md")
		err := testNote.SetEncryption(encType)
		if err != nil {
			t.Errorf("設定有效加密類型 '%s' 不應該產生錯誤：%v", encType, err)
		}
		if !testNote.IsEncrypted {
			t.Errorf("設定加密類型 '%s' 後，筆記應該被標記為已加密", encType)
		}
		if testNote.EncryptionType != encType {
			t.Errorf("期望加密類型為 '%s'，實際得到 '%s'", encType, testNote.EncryptionType)
		}
	}

	// 測試無效的加密類型
	err := note.SetEncryption("invalid_type")
	if err == nil {
		t.Error("設定無效加密類型應該產生錯誤")
	}

	// 測試修改時間是否更新
	note.SetEncryption("password")
	if !note.UpdatedAt.After(originalUpdateTime) {
		t.Error("設定加密後修改時間應該被更新")
	}
}

// TestNote_RemoveEncryption 測試移除筆記加密功能
// 驗證加密移除後狀態是否正確重置
func TestNote_RemoveEncryption(t *testing.T) {
	note := NewNote("測試標題", "測試內容", "/test.md")
	
	// 先設定加密
	note.SetEncryption("password")
	originalUpdateTime := note.UpdatedAt

	// 等待一小段時間確保時間戳不同
	time.Sleep(time.Millisecond * 10)

	// 移除加密
	note.RemoveEncryption()

	// 驗證加密狀態已重置
	if note.IsEncrypted {
		t.Error("移除加密後，筆記不應該被標記為已加密")
	}
	if note.EncryptionType != "" {
		t.Error("移除加密後，加密類型應該為空字串")
	}

	// 驗證修改時間已更新
	if !note.UpdatedAt.After(originalUpdateTime) {
		t.Error("移除加密後修改時間應該被更新")
	}
}

// TestNote_GetWordCount 測試筆記字數統計功能
// 驗證各種內容情況下的字數統計是否正確
func TestNote_GetWordCount(t *testing.T) {
	// 測試空內容
	note := NewNote("標題", "", "/test.md")
	if count := note.GetWordCount(); count != 0 {
		t.Errorf("空內容的字數應該為 0，實際得到 %d", count)
	}

	// 測試單個單詞
	note.Content = "Hello"
	if count := note.GetWordCount(); count != 1 {
		t.Errorf("單個單詞的字數應該為 1，實際得到 %d", count)
	}

	// 測試多個單詞
	note.Content = "Hello World Test"
	if count := note.GetWordCount(); count != 3 {
		t.Errorf("三個單詞的字數應該為 3，實際得到 %d", count)
	}

	// 測試包含多個空白字符的內容
	note.Content = "  Hello   World  \n\t Test  "
	if count := note.GetWordCount(); count != 3 {
		t.Errorf("包含多個空白字符的內容字數應該為 3，實際得到 %d", count)
	}

	// 測試中文內容
	note.Content = "這是 一個 測試"
	if count := note.GetWordCount(); count != 3 {
		t.Errorf("中文內容的字數應該為 3，實際得到 %d", count)
	}
}

// TestNote_Clone 測試筆記複製功能
// 驗證複製的筆記是否包含相同資料但不同記憶體位址
func TestNote_Clone(t *testing.T) {
	// 建立原始筆記
	original := NewNote("原始標題", "原始內容", "/original.md")
	original.SetEncryption("password")
	original.MarkSaved()

	// 複製筆記
	cloned := original.Clone()

	// 驗證資料相同
	if cloned.ID != original.ID {
		t.Error("複製的筆記 ID 應該相同")
	}
	if cloned.Title != original.Title {
		t.Error("複製的筆記標題應該相同")
	}
	if cloned.Content != original.Content {
		t.Error("複製的筆記內容應該相同")
	}
	if cloned.FilePath != original.FilePath {
		t.Error("複製的筆記檔案路徑應該相同")
	}
	if cloned.IsEncrypted != original.IsEncrypted {
		t.Error("複製的筆記加密狀態應該相同")
	}
	if cloned.EncryptionType != original.EncryptionType {
		t.Error("複製的筆記加密類型應該相同")
	}
	if !cloned.CreatedAt.Equal(original.CreatedAt) {
		t.Error("複製的筆記建立時間應該相同")
	}
	if !cloned.UpdatedAt.Equal(original.UpdatedAt) {
		t.Error("複製的筆記修改時間應該相同")
	}
	if !cloned.LastSaved.Equal(original.LastSaved) {
		t.Error("複製的筆記最後保存時間應該相同")
	}

	// 驗證記憶體位址不同
	if cloned == original {
		t.Error("複製的筆記應該有不同的記憶體位址")
	}

	// 驗證修改複製品不會影響原始筆記
	cloned.UpdateContent("修改後的內容")
	if original.Content == cloned.Content {
		t.Error("修改複製品不應該影響原始筆記")
	}
}

// TestGenerateID 測試 ID 生成功能
// 驗證生成的 ID 是否符合預期格式和唯一性
func TestGenerateID(t *testing.T) {
	// 生成多個 ID 測試唯一性
	ids := make(map[string]bool)
	for i := 0; i < 100; i++ {
		id := generateID()
		
		// 檢查格式
		if len(id) < 15 {
			t.Errorf("ID 長度不足：%s", id)
		}
		if !strings.Contains(id, "-") {
			t.Errorf("ID 應該包含連字符：%s", id)
		}
		
		// 檢查唯一性
		if ids[id] {
			t.Errorf("發現重複的 ID：%s", id)
		}
		ids[id] = true
	}
}

// TestRandomString 測試隨機字串生成功能
// 驗證生成的隨機字串長度和字符集是否正確
func TestRandomString(t *testing.T) {
	// 測試不同長度
	lengths := []int{1, 5, 10, 20}
	for _, length := range lengths {
		str := randomString(length)
		if len(str) != length {
			t.Errorf("期望長度 %d，實際得到 %d：%s", length, len(str), str)
		}
		
		// 檢查字符集（應該只包含字母和數字）
		for _, char := range str {
			if !((char >= 'a' && char <= 'z') || 
				 (char >= 'A' && char <= 'Z') || 
				 (char >= '0' && char <= '9')) {
				t.Errorf("隨機字串包含無效字符：%c in %s", char, str)
			}
		}
	}
	
	// 測試零長度
	str := randomString(0)
	if len(str) != 0 {
		t.Errorf("零長度應該回傳空字串，實際得到：%s", str)
	}
}

// BenchmarkNewNote 效能測試：筆記建立
func BenchmarkNewNote(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewNote("測試標題", "測試內容", "/test.md")
	}
}

// BenchmarkNote_Validate 效能測試：筆記驗證
func BenchmarkNote_Validate(b *testing.B) {
	note := NewNote("測試標題", "測試內容", "/test.md")
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		note.Validate()
	}
}

// BenchmarkNote_Clone 效能測試：筆記複製
func BenchmarkNote_Clone(b *testing.B) {
	note := NewNote("測試標題", "測試內容", "/test.md")
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		note.Clone()
	}
}