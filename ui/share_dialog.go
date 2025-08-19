// Package ui 提供分享對話框的 UI 元件
// 負責處理筆記分享的用戶介面，包含分享類型選擇、選項設定和分享執行
package ui

import (
	"fmt"
	"mac-notebook-app/internal/models"
	"mac-notebook-app/internal/services"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// ShareDialog 分享對話框結構
// 提供完整的筆記分享功能介面
type ShareDialog struct {
	// UI 元件
	window       fyne.Window                // 父視窗
	dialog       *dialog.CustomDialog       // 自訂對話框
	content      *fyne.Container            // 對話框內容容器
	
	// 分享選項 UI 元件
	shareTypeSelect *widget.Select          // 分享類型選擇
	passwordEntry   *widget.Entry           // 分享密碼輸入
	expirySelect    *widget.Select          // 過期時間選擇
	allowDownload   *widget.Check           // 允許下載選項
	allowEdit       *widget.Check           // 允許編輯選項
	recipientsEntry *widget.Entry           // 收件人輸入（電子郵件分享）
	
	// 分享結果顯示
	resultLabel     *widget.Label           // 結果標籤
	shareURLEntry   *widget.Entry           // 分享連結顯示
	copyButton      *widget.Button          // 複製連結按鈕
	
	// 按鈕
	shareButton     *widget.Button          // 分享按鈕
	cancelButton    *widget.Button          // 取消按鈕
	
	// 服務和資料
	exportService   services.ExportService  // 匯出服務（包含分享功能）
	note           *models.Note             // 要分享的筆記
	
	// 回調函數
	onShareCompleteCallback func(success bool, shareResult *services.ShareResult) // 分享完成回調
}

// NewShareDialog 建立新的分享對話框
// 參數：window（父視窗）、exportService（匯出服務）、note（要分享的筆記）
// 回傳：ShareDialog 實例
//
// 執行流程：
// 1. 初始化對話框結構和基本屬性
// 2. 建立所有 UI 元件
// 3. 設定預設值和事件處理
// 4. 組裝對話框佈局
func NewShareDialog(window fyne.Window, exportService services.ExportService, note *models.Note) *ShareDialog {
	dialog := &ShareDialog{
		window:        window,
		exportService: exportService,
		note:         note,
	}
	
	// 建立 UI 元件
	dialog.createUIComponents()
	
	// 設定預設值
	dialog.setDefaultValues()
	
	// 建立對話框佈局
	dialog.createLayout()
	
	return dialog
}

// Show 顯示分享對話框
// 將對話框顯示給用戶並等待操作
func (d *ShareDialog) Show() {
	d.dialog.Show()
}

// Hide 隱藏分享對話框
// 關閉對話框並清理資源
func (d *ShareDialog) Hide() {
	if d.dialog != nil {
		d.dialog.Hide()
	}
}

// SetOnShareComplete 設定分享完成回調函數
// 參數：callback（分享完成時的回調函數）
func (d *ShareDialog) SetOnShareComplete(callback func(success bool, shareResult *services.ShareResult)) {
	d.onShareCompleteCallback = callback
}

// createUIComponents 建立所有 UI 元件
// 初始化對話框中的所有控制項和輸入元件
func (d *ShareDialog) createUIComponents() {
	// 分享密碼
	d.passwordEntry = widget.NewPasswordEntry()
	d.passwordEntry.SetPlaceHolder("設定分享密碼（可選）")
	
	// 過期時間選擇
	d.expirySelect = widget.NewSelect([]string{
		"1 小時",
		"24 小時",
		"7 天",
		"30 天",
		"永不過期",
	}, nil)
	d.expirySelect.SetSelected("24 小時")
	
	// 權限選項
	d.allowDownload = widget.NewCheck("允許下載", nil)
	d.allowDownload.SetChecked(true)
	
	d.allowEdit = widget.NewCheck("允許編輯", nil)
	d.allowEdit.SetChecked(false)
	
	// 收件人輸入（電子郵件分享用）
	d.recipientsEntry = widget.NewEntry()
	d.recipientsEntry.SetPlaceHolder("輸入收件人電子郵件地址，多個地址用逗號分隔")
	d.recipientsEntry.Hide() // 預設隱藏
	
	// 分享類型選擇（放在最後，避免在其他組件創建前觸發回調）
	d.shareTypeSelect = widget.NewSelect([]string{
		"連結分享",
		"電子郵件分享", 
		"AirDrop 分享",
		"複製到剪貼簿",
	}, d.onShareTypeChanged)
	d.shareTypeSelect.SetSelected("連結分享")
	
	// 分享結果顯示
	d.resultLabel = widget.NewLabel("")
	d.resultLabel.Hide()
	
	d.shareURLEntry = widget.NewEntry()
	d.shareURLEntry.SetPlaceHolder("分享連結將顯示在這裡")
	d.shareURLEntry.Disable()
	d.shareURLEntry.Hide()
	
	d.copyButton = widget.NewButton("複製連結", d.onCopyClicked)
	d.copyButton.Hide()
	
	// 按鈕
	d.shareButton = widget.NewButton("分享", d.onShareClicked)
	d.shareButton.Importance = widget.HighImportance
	
	d.cancelButton = widget.NewButton("取消", d.onCancelClicked)
}

// setDefaultValues 設定預設值
// 根據筆記資訊和用戶偏好設定預設的分享選項
func (d *ShareDialog) setDefaultValues() {
	// 根據筆記類型設定預設權限
	if d.note != nil && d.note.IsEncrypted {
		// 加密筆記預設不允許編輯
		d.allowEdit.SetChecked(false)
	}
}

// createLayout 建立對話框佈局
// 組織所有 UI 元件到適當的佈局容器中
func (d *ShareDialog) createLayout() {
	// 分享設定區域
	shareSection := container.NewVBox(
		widget.NewLabel("分享設定"),
		widget.NewSeparator(),
		container.NewGridWithColumns(2,
			widget.NewLabel("分享方式:"),
			d.shareTypeSelect,
		),
		d.passwordEntry,
		container.NewGridWithColumns(2,
			widget.NewLabel("過期時間:"),
			d.expirySelect,
		),
		container.NewGridWithColumns(2,
			d.allowDownload,
			d.allowEdit,
		),
		d.recipientsEntry,
	)
	
	// 分享結果區域
	resultSection := container.NewVBox(
		d.resultLabel,
		d.shareURLEntry,
		d.copyButton,
	)
	
	// 按鈕區域
	buttonSection := container.NewBorder(nil, nil, nil, 
		container.NewHBox(d.shareButton, d.cancelButton))
	
	// 主要內容
	d.content = container.NewVBox(
		shareSection,
		widget.NewSeparator(),
		resultSection,
		buttonSection,
	)
	
	// 建立自訂對話框
	title := "分享筆記"
	if d.note != nil && d.note.Title != "" {
		title = fmt.Sprintf("分享筆記: %s", d.note.Title)
	}
	
	d.dialog = dialog.NewCustom(title, "關閉", d.content, d.window)
	d.dialog.Resize(fyne.NewSize(450, 400))
}

// 事件處理方法

// onShareTypeChanged 處理分享類型變更事件
// 參數：shareType（選擇的分享類型）
func (d *ShareDialog) onShareTypeChanged(shareType string) {
	// 根據分享類型顯示/隱藏相關選項
	switch shareType {
	case "電子郵件分享":
		d.recipientsEntry.Show()
		d.passwordEntry.Show()
		d.expirySelect.Show()
		d.allowDownload.Show()
		d.allowEdit.Show()
		
	case "連結分享":
		d.recipientsEntry.Hide()
		d.passwordEntry.Show()
		d.expirySelect.Show()
		d.allowDownload.Show()
		d.allowEdit.Show()
		
	case "AirDrop 分享":
		d.recipientsEntry.Hide()
		d.passwordEntry.Hide()
		d.expirySelect.Hide()
		d.allowDownload.Hide()
		d.allowEdit.Hide()
		
	case "複製到剪貼簿":
		d.recipientsEntry.Hide()
		d.passwordEntry.Hide()
		d.expirySelect.Hide()
		d.allowDownload.Hide()
		d.allowEdit.Hide()
	}
	
	// 重新整理佈局
	if d.content != nil {
		d.content.Refresh()
	}
}

// onShareClicked 處理分享按鈕點擊事件
func (d *ShareDialog) onShareClicked() {
	// 驗證輸入
	if !d.validateInput() {
		return
	}
	
	// 禁用分享按鈕
	d.shareButton.SetText("分享中...")
	d.shareButton.Disable()
	
	// 建立分享選項
	shareOptions := d.createShareOptions()
	
	// 在背景執行分享
	go d.performShare(shareOptions)
}

// onCancelClicked 處理取消按鈕點擊事件
func (d *ShareDialog) onCancelClicked() {
	d.Hide()
}

// onCopyClicked 處理複製連結按鈕點擊事件
func (d *ShareDialog) onCopyClicked() {
	// 複製分享連結到剪貼簿
	shareURL := d.shareURLEntry.Text
	if shareURL != "" {
		d.window.Clipboard().SetContent(shareURL)
		d.showSuccess("分享連結已複製到剪貼簿")
	}
}

// 輔助方法

// validateInput 驗證用戶輸入
// 回傳：輸入是否有效
func (d *ShareDialog) validateInput() bool {
	shareType := d.shareTypeSelect.Selected
	
	// 電子郵件分享需要收件人
	if shareType == "電子郵件分享" {
		recipients := d.recipientsEntry.Text
		if recipients == "" {
			d.showError("請輸入收件人電子郵件地址")
			return false
		}
		
		// 簡單的電子郵件格式驗證
		if !d.isValidEmailList(recipients) {
			d.showError("請輸入有效的電子郵件地址")
			return false
		}
	}
	
	return true
}

// createShareOptions 建立分享選項
// 回傳：分享選項結構
func (d *ShareDialog) createShareOptions() *services.ShareOptions {
	shareType := d.getShareType()
	expiryTime := d.getExpiryTime()
	
	options := &services.ShareOptions{
		ShareType:     shareType,
		ExpiryTime:    expiryTime,
		Password:      d.passwordEntry.Text,
		AllowDownload: d.allowDownload.Checked,
		AllowEdit:     d.allowEdit.Checked,
	}
	
	// 電子郵件分享需要收件人列表
	if shareType == services.ShareTypeEmail {
		recipients := d.parseEmailList(d.recipientsEntry.Text)
		options.Recipients = recipients
	}
	
	return options
}

// getShareType 取得選擇的分享類型
// 回傳：分享類型列舉值
func (d *ShareDialog) getShareType() services.ShareType {
	switch d.shareTypeSelect.Selected {
	case "連結分享":
		return services.ShareTypeLink
	case "電子郵件分享":
		return services.ShareTypeEmail
	case "AirDrop 分享":
		return services.ShareTypeAirDrop
	case "複製到剪貼簿":
		return services.ShareTypeClipboard
	default:
		return services.ShareTypeLink
	}
}

// getExpiryTime 取得過期時間
// 回傳：過期時間
func (d *ShareDialog) getExpiryTime() time.Time {
	now := time.Now()
	
	switch d.expirySelect.Selected {
	case "1 小時":
		return now.Add(1 * time.Hour)
	case "24 小時":
		return now.Add(24 * time.Hour)
	case "7 天":
		return now.Add(7 * 24 * time.Hour)
	case "30 天":
		return now.Add(30 * 24 * time.Hour)
	case "永不過期":
		return now.Add(100 * 365 * 24 * time.Hour) // 100 年後
	default:
		return now.Add(24 * time.Hour)
	}
}

// performShare 執行分享操作
// 參數：shareOptions（分享選項）
func (d *ShareDialog) performShare(shareOptions *services.ShareOptions) {
	// 執行分享
	result, err := d.exportService.ShareNote(d.note, shareOptions)
	
	// 更新 UI（在主執行緒中）
	go func() {
		time.Sleep(100 * time.Millisecond) // 短暫延遲確保分享完成
		d.onShareComplete(result, err)
	}()
}

// onShareComplete 處理分享完成事件
// 參數：result（分享結果）、err（錯誤資訊）
func (d *ShareDialog) onShareComplete(result *services.ShareResult, err error) {
	// 重置按鈕狀態
	d.shareButton.SetText("分享")
	d.shareButton.Enable()
	
	if err != nil || result == nil || !result.Success {
		// 分享失敗
		errorMsg := "分享失敗"
		if err != nil {
			errorMsg = err.Error()
		} else if result != nil {
			errorMsg = result.Message
		}
		
		d.resultLabel.SetText(errorMsg)
		d.resultLabel.Show()
		d.showError(errorMsg)
		
		// 呼叫回調函數
		if d.onShareCompleteCallback != nil {
			d.onShareCompleteCallback(false, result)
		}
		
	} else {
		// 分享成功
		d.resultLabel.SetText("分享成功！")
		d.resultLabel.Show()
		
		// 如果有分享連結，顯示它
		if result.ShareURL != "" {
			d.shareURLEntry.SetText(result.ShareURL)
			d.shareURLEntry.Show()
			d.copyButton.Show()
		}
		
		d.showSuccess(result.Message)
		
		// 呼叫回調函數
		if d.onShareCompleteCallback != nil {
			d.onShareCompleteCallback(true, result)
		}
	}
	
	// 重新整理佈局
	if d.content != nil {
		d.content.Refresh()
	}
}

// isValidEmailList 驗證電子郵件地址列表
// 參數：emailList（電子郵件地址列表字串）
// 回傳：是否有效
func (d *ShareDialog) isValidEmailList(emailList string) bool {
	emails := d.parseEmailList(emailList)
	
	for _, email := range emails {
		if !d.isValidEmail(email) {
			return false
		}
	}
	
	return len(emails) > 0
}

// isValidEmail 驗證單個電子郵件地址
// 參數：email（電子郵件地址）
// 回傳：是否有效
func (d *ShareDialog) isValidEmail(email string) bool {
	// 簡單的電子郵件格式驗證
	if len(email) < 5 {
		return false
	}
	
	atIndex := -1
	dotIndex := -1
	
	for i, char := range email {
		if char == '@' {
			if atIndex != -1 {
				return false // 多個 @ 符號
			}
			atIndex = i
		} else if char == '.' && atIndex != -1 {
			dotIndex = i
		}
	}
	
	return atIndex > 0 && dotIndex > atIndex+1 && dotIndex < len(email)-1
}

// parseEmailList 解析電子郵件地址列表
// 參數：emailList（電子郵件地址列表字串）
// 回傳：電子郵件地址陣列
func (d *ShareDialog) parseEmailList(emailList string) []string {
	var emails []string
	
	// 按逗號分割
	parts := splitString(emailList, ",")
	
	for _, part := range parts {
		email := trimString(part)
		if email != "" {
			emails = append(emails, email)
		}
	}
	
	return emails
}

// showError 顯示錯誤訊息
// 參數：message（錯誤訊息）
func (d *ShareDialog) showError(message string) {
	dialog.ShowError(fmt.Errorf("%s", message), d.window)
}

// showSuccess 顯示成功訊息
// 參數：message（成功訊息）
func (d *ShareDialog) showSuccess(message string) {
	dialog.ShowInformation("分享成功", message, d.window)
}

// 輔助函數

// splitString 分割字串
// 參數：s（要分割的字串）、sep（分隔符）
// 回傳：分割後的字串陣列
func splitString(s, sep string) []string {
	var result []string
	var current string
	
	for _, char := range s {
		if string(char) == sep {
			if current != "" {
				result = append(result, current)
				current = ""
			}
		} else {
			current += string(char)
		}
	}
	
	if current != "" {
		result = append(result, current)
	}
	
	return result
}

// trimString 去除字串前後空白
// 參數：s（要處理的字串）
// 回傳：處理後的字串
func trimString(s string) string {
	// 去除前面的空白
	start := 0
	for start < len(s) && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n') {
		start++
	}
	
	// 去除後面的空白
	end := len(s)
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n') {
		end--
	}
	
	return s[start:end]
}