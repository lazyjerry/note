// Package models 定義了應用程式的核心資料模型
// 包含筆記、設定、檔案資訊等結構體定義
package models

import "time"

// Note 代表一個具有加密功能的 Markdown 筆記
// 這個結構體包含了筆記的所有基本資訊和加密相關的元資料
type Note struct {
	ID             string    `json:"id"`              // 筆記的唯一識別符
	Title          string    `json:"title"`           // 筆記標題
	Content        string    `json:"content"`         // 筆記的 Markdown 內容
	FilePath       string    `json:"file_path"`       // 筆記檔案在檔案系統中的路徑
	IsEncrypted    bool      `json:"is_encrypted"`    // 標示筆記是否已加密
	EncryptionType string    `json:"encryption_type"` // 加密類型："password"（密碼）、"biometric"（生物識別）、"both"（兩者皆有）
	CreatedAt      time.Time `json:"created_at"`      // 筆記建立時間
	UpdatedAt      time.Time `json:"updated_at"`      // 筆記最後修改時間
	LastSaved      time.Time `json:"last_saved"`      // 筆記最後保存時間
}

// NewNote 建立一個新的筆記實例並設定預設值
// 參數：
//   - title: 筆記標題
//   - content: 筆記的 Markdown 內容
//   - filePath: 筆記檔案的儲存路徑
// 回傳：指向新建立筆記的指標
//
// 執行流程：
// 1. 取得當前時間作為建立和修改時間
// 2. 產生唯一的筆記 ID
// 3. 設定預設的非加密狀態
// 4. 將 LastSaved 設為零值表示尚未保存
func NewNote(title, content, filePath string) *Note {
	now := time.Now()
	return &Note{
		ID:             generateID(),    // 產生唯一識別符
		Title:          title,           // 設定筆記標題
		Content:        content,         // 設定筆記內容
		FilePath:       filePath,        // 設定檔案路徑
		IsEncrypted:    false,           // 預設為未加密狀態
		EncryptionType: "",              // 預設無加密類型
		CreatedAt:      now,             // 設定建立時間
		UpdatedAt:      now,             // 設定修改時間
		LastSaved:      time.Time{},     // 零時間表示尚未保存
	}
}

// UpdateContent 更新筆記內容並更新修改時間戳
// 參數：
//   - content: 新的筆記內容
//
// 執行流程：
// 1. 更新筆記的內容
// 2. 將修改時間設為當前時間
// 這個方法會自動標記筆記為已修改狀態
func (n *Note) UpdateContent(content string) {
	n.Content = content      // 更新筆記內容
	n.UpdatedAt = time.Now() // 更新修改時間戳
}

// MarkSaved 標記筆記為已保存狀態並更新保存時間戳
// 這個方法在筆記成功保存到檔案系統後呼叫
//
// 執行流程：
// 1. 將最後保存時間設為當前時間
// 2. 這會使 IsModified() 方法回傳 false
func (n *Note) MarkSaved() {
	n.LastSaved = time.Now() // 更新最後保存時間
}

// IsModified 檢查筆記自上次保存後是否已被修改
// 回傳：如果筆記已修改則回傳 true，否則回傳 false
//
// 判斷邏輯：
// - 比較 UpdatedAt（最後修改時間）和 LastSaved（最後保存時間）
// - 如果修改時間晚於保存時間，表示筆記已被修改
// - 如果 LastSaved 為零值，表示從未保存過，視為已修改
func (n *Note) IsModified() bool {
	return n.UpdatedAt.After(n.LastSaved)
}

// generateID 為筆記產生唯一的識別符
// 回傳：格式為 "YYYYMMDDHHMMSS-XXXXXXXX" 的唯一 ID 字串
//
// ID 組成：
// - 前半部：當前時間戳（年月日時分秒）
// - 後半部：8 位隨機字串
// 這種格式確保 ID 的唯一性和時間順序性
func generateID() string {
	// 使用時間戳格式 "20060102150405" (Go 的參考時間格式)
	timeStamp := time.Now().Format("20060102150405")
	// 組合時間戳和隨機字串
	return timeStamp + "-" + randomString(8)
}

// randomString 產生指定長度的隨機字串
// 參數：
//   - length: 要產生的隨機字串長度
// 回傳：包含字母和數字的隨機字串
//
// 執行流程：
// 1. 定義字符集（包含大小寫字母和數字）
// 2. 建立指定長度的位元組陣列
// 3. 使用當前時間的奈秒值作為隨機種子
// 4. 從字符集中隨機選擇字符填入陣列
// 5. 轉換為字串並回傳
func randomString(length int) string {
	// 定義可用的字符集
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	
	// 為每個位置隨機選擇字符
	for i := range b {
		// 使用當前時間的奈秒值作為隨機種子
		randomIndex := time.Now().UnixNano() % int64(len(charset))
		b[i] = charset[randomIndex]
	}
	
	return string(b)
}