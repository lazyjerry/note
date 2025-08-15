package ui

import (
	"testing"
	"time"

	"fyne.io/fyne/v2/test"

	"mac-notebook-app/internal/services"
)

// TestNewNotificationWidget 測試通知元件的建立
// 驗證元件能夠正確初始化
func TestNewNotificationWidget(t *testing.T) {
	// 建立通知服務
	notificationService := services.NewNotificationService()
	
	// 建立通知元件
	widget := NewNotificationWidget(notificationService)
	
	if widget == nil {
		t.Error("通知元件不應該為空")
	}
	
	if widget.notificationService == nil {
		t.Error("通知服務不應該為空")
	}
	
	if widget.container == nil {
		t.Error("容器不應該為空")
	}
	
	if widget.activeNotifications == nil {
		t.Error("活躍通知映射表不應該為空")
	}
	
	// 驗證初始狀態
	if widget.GetNotificationCount() != 0 {
		t.Errorf("初始通知數量應該為 0，得到 %d", widget.GetNotificationCount())
	}
}

// TestNotificationDisplay 測試通知顯示功能
// 驗證通知能夠正確顯示在 UI 中
func TestNotificationDisplay(t *testing.T) {
	// 建立測試應用
	app := test.NewApp()
	defer app.Quit()
	
	// 建立通知服務和元件
	notificationService := services.NewNotificationService()
	widget := NewNotificationWidget(notificationService)
	
	// 顯示測試通知
	notificationID := notificationService.ShowInfo("測試標題", "測試內容")
	
	// 等待 UI 更新
	time.Sleep(10 * time.Millisecond)
	
	// 驗證通知被添加到 UI
	if widget.GetNotificationCount() != 1 {
		t.Errorf("應該有 1 個通知，得到 %d", widget.GetNotificationCount())
	}
	
	// 驗證通知項目存在
	if _, exists := widget.activeNotifications[notificationID]; !exists {
		t.Error("通知項目應該存在於活躍通知映射表中")
	}
	
	// 關閉通知
	notificationService.DismissNotification(notificationID)
	
	// 等待 UI 更新
	time.Sleep(10 * time.Millisecond)
	
	// 驗證通知被移除
	if widget.GetNotificationCount() != 0 {
		t.Errorf("通知應該被移除，但還有 %d 個", widget.GetNotificationCount())
	}
}

// TestMultipleNotifications 測試多個通知的顯示
// 驗證能夠同時顯示多個不同類型的通知
func TestMultipleNotifications(t *testing.T) {
	// 建立測試應用
	app := test.NewApp()
	defer app.Quit()
	
	// 建立通知服務和元件
	notificationService := services.NewNotificationService()
	widget := NewNotificationWidget(notificationService)
	
	// 顯示多個不同類型的通知
	id1 := notificationService.ShowInfo("資訊", "資訊內容")
	id2 := notificationService.ShowSuccess("成功", "成功內容")
	id3 := notificationService.ShowWarning("警告", "警告內容")
	id4 := notificationService.ShowError("錯誤", "錯誤內容")
	
	// 等待 UI 更新
	time.Sleep(10 * time.Millisecond)
	
	// 驗證所有通知都被顯示
	if widget.GetNotificationCount() != 4 {
		t.Errorf("應該有 4 個通知，得到 %d", widget.GetNotificationCount())
	}
	
	// 驗證每個通知項目都存在
	expectedIDs := []string{id1, id2, id3, id4}
	for _, id := range expectedIDs {
		if _, exists := widget.activeNotifications[id]; !exists {
			t.Errorf("通知 %s 應該存在於活躍通知映射表中", id)
		}
	}
	
	// 清除所有通知
	widget.ClearAllNotifications()
	
	// 等待 UI 更新
	time.Sleep(10 * time.Millisecond)
	
	// 驗證所有通知都被清除
	if widget.GetNotificationCount() != 0 {
		t.Errorf("所有通知應該被清除，但還有 %d 個", widget.GetNotificationCount())
	}
}

// TestSaveStatusDisplay 測試保存狀態顯示功能
// 驗證保存狀態能夠正確顯示和更新
func TestSaveStatusDisplay(t *testing.T) {
	// 建立測試應用
	app := test.NewApp()
	defer app.Quit()
	
	// 建立通知服務和元件
	notificationService := services.NewNotificationService()
	widget := NewNotificationWidget(notificationService)
	
	noteID := "test-note"
	fileName := "test.md"
	
	// 測試保存中狀態
	savingStatus := services.SaveStatusInfo{
		NoteID:       noteID,
		FileName:     fileName,
		IsSaving:     true,
		SaveProgress: 0.5,
		HasChanges:   true,
	}
	
	notificationService.UpdateSaveStatus(noteID, fileName, savingStatus)
	
	// 等待 UI 更新
	time.Sleep(10 * time.Millisecond)
	
	// 驗證保存狀態被更新
	if widget.currentSaveStatus == nil {
		t.Error("當前保存狀態不應該為空")
	}
	
	if widget.currentSaveStatus.IsSaving != true {
		t.Error("保存狀態應該為正在保存中")
	}
	
	if widget.currentSaveStatus.SaveProgress != 0.5 {
		t.Errorf("保存進度應該為 0.5，得到 %f", widget.currentSaveStatus.SaveProgress)
	}
	
	// 測試保存完成狀態
	completedStatus := services.SaveStatusInfo{
		NoteID:       noteID,
		FileName:     fileName,
		IsSaving:     false,
		SaveProgress: 1.0,
		HasChanges:   false,
		LastSaved:    time.Now(),
	}
	
	notificationService.UpdateSaveStatus(noteID, fileName, completedStatus)
	
	// 等待 UI 更新
	time.Sleep(10 * time.Millisecond)
	
	// 驗證保存狀態被更新
	if widget.currentSaveStatus.IsSaving != false {
		t.Error("保存狀態應該為已完成")
	}
	
	if widget.currentSaveStatus.HasChanges != false {
		t.Error("應該沒有未保存的變更")
	}
	
	// 測試清除狀態
	notificationService.ClearSaveStatus(noteID)
	
	// 等待 UI 更新
	time.Sleep(10 * time.Millisecond)
	
	// 驗證狀態被清除
	if widget.currentSaveStatus != nil {
		t.Error("保存狀態應該被清除")
	}
}

// TestNotificationVisibility 測試通知元件的可見性控制
// 驗證能夠正確控制元件的顯示和隱藏
func TestNotificationVisibility(t *testing.T) {
	// 建立測試應用
	app := test.NewApp()
	defer app.Quit()
	
	// 建立通知服務和元件
	notificationService := services.NewNotificationService()
	widget := NewNotificationWidget(notificationService)
	
	// 測試初始可見性（應該是可見的）
	if !widget.container.Visible() {
		t.Error("通知元件初始應該是可見的")
	}
	
	// 隱藏元件
	widget.SetVisible(false)
	if widget.container.Visible() {
		t.Error("通知元件應該被隱藏")
	}
	
	// 顯示元件
	widget.SetVisible(true)
	if !widget.container.Visible() {
		t.Error("通知元件應該被顯示")
	}
}

// TestShowTestNotifications 測試顯示測試通知功能
// 驗證測試通知功能能夠正常工作
func TestShowTestNotifications(t *testing.T) {
	// 建立測試應用
	app := test.NewApp()
	defer app.Quit()
	
	// 建立通知服務和元件
	notificationService := services.NewNotificationService()
	widget := NewNotificationWidget(notificationService)
	
	// 顯示測試通知
	widget.ShowTestNotifications()
	
	// 等待通知顯示
	time.Sleep(50 * time.Millisecond)
	
	// 驗證測試通知被顯示
	if widget.GetNotificationCount() == 0 {
		t.Error("應該有測試通知被顯示")
	}
	
	// 驗證至少有 4 個不同類型的通知
	if widget.GetNotificationCount() < 4 {
		t.Errorf("應該至少有 4 個測試通知，得到 %d", widget.GetNotificationCount())
	}
}

// TestFormatTimeAgo 測試時間格式化功能
// 驗證時間能夠正確格式化為「多久之前」的形式
func TestFormatTimeAgo(t *testing.T) {
	now := time.Now()
	
	tests := []struct {
		name     string
		time     time.Time
		expected string
	}{
		{
			name:     "剛剛",
			time:     now.Add(-30 * time.Second),
			expected: "剛剛",
		},
		{
			name:     "5分鐘前",
			time:     now.Add(-5 * time.Minute),
			expected: "5 分鐘前",
		},
		{
			name:     "2小時前",
			time:     now.Add(-2 * time.Hour),
			expected: "2 小時前",
		},
		{
			name:     "3天前",
			time:     now.Add(-3 * 24 * time.Hour),
			expected: "3 天前",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatTimeAgo(tt.time)
			if result != tt.expected {
				t.Errorf("formatTimeAgo() = %v, 期望 %v", result, tt.expected)
			}
		})
	}
}

// TestNotificationItemCreation 測試通知項目建立功能
// 驗證通知項目能夠正確建立和配置
func TestNotificationItemCreation(t *testing.T) {
	// 建立測試應用
	app := test.NewApp()
	defer app.Quit()
	
	// 建立通知服務和元件
	notificationService := services.NewNotificationService()
	widget := NewNotificationWidget(notificationService)
	
	// 建立測試通知
	notification := &services.Notification{
		ID:       "test-notification",
		Type:     services.NotificationInfo,
		Title:    "測試標題",
		Message:  "測試內容",
		Duration: 3 * time.Second,
	}
	
	// 建立通知項目
	item := widget.createNotificationItem(notification)
	
	if item == nil {
		t.Error("通知項目不應該為空")
	}
	
	if item.container == nil {
		t.Error("通知項目容器不應該為空")
	}
	
	if item.titleLabel == nil {
		t.Error("標題標籤不應該為空")
	}
	
	if item.messageLabel == nil {
		t.Error("內容標籤不應該為空")
	}
	
	if item.closeButton == nil {
		t.Error("關閉按鈕不應該為空")
	}
	
	if item.notification != notification {
		t.Error("通知項目應該保存通知引用")
	}
	
	// 驗證標籤內容
	if item.titleLabel.Text != "測試標題" {
		t.Errorf("標題標籤內容 = %v, 期望 %v", item.titleLabel.Text, "測試標題")
	}
	
	if item.messageLabel.Text != "測試內容" {
		t.Errorf("內容標籤內容 = %v, 期望 %v", item.messageLabel.Text, "測試內容")
	}
}

// TestNotificationCallback 測試通知回調功能
// 驗證通知回調能夠正確觸發 UI 更新
func TestNotificationCallback(t *testing.T) {
	// 建立測試應用
	app := test.NewApp()
	defer app.Quit()
	
	// 建立通知服務和元件
	notificationService := services.NewNotificationService()
	widget := NewNotificationWidget(notificationService)
	
	// 直接觸發通知回調
	notification := &services.Notification{
		ID:      "callback-test",
		Type:    services.NotificationSuccess,
		Title:   "回調測試",
		Message: "回調內容",
		IsRead:  false,
	}
	
	// 觸發添加通知回調
	widget.onNotificationUpdate(notification)
	
	// 驗證通知被添加
	if widget.GetNotificationCount() != 1 {
		t.Errorf("應該有 1 個通知，得到 %d", widget.GetNotificationCount())
	}
	
	// 觸發移除通知回調
	notification.IsRead = true
	widget.onNotificationUpdate(notification)
	
	// 驗證通知被移除
	if widget.GetNotificationCount() != 0 {
		t.Errorf("通知應該被移除，但還有 %d 個", widget.GetNotificationCount())
	}
	
	// 測試空通知回調
	widget.onNotificationUpdate(nil)
	
	// 應該不會造成錯誤或變更
	if widget.GetNotificationCount() != 0 {
		t.Error("空通知回調不應該影響通知數量")
	}
}

// BenchmarkNotificationDisplay 效能測試 - 通知顯示
// 測試通知顯示功能的效能表現
func BenchmarkNotificationDisplay(b *testing.B) {
	// 建立測試應用
	app := test.NewApp()
	defer app.Quit()
	
	// 建立通知服務和元件
	notificationService := services.NewNotificationService()
	notificationWidget := NewNotificationWidget(notificationService)
	
	// 確保元件被使用
	_ = notificationWidget
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		notificationService.ShowInfo("效能測試", "測試內容")
	}
}

// BenchmarkSaveStatusUpdate 效能測試 - 保存狀態更新
// 測試保存狀態更新功能的效能表現
func BenchmarkSaveStatusUpdate(b *testing.B) {
	// 建立測試應用
	app := test.NewApp()
	defer app.Quit()
	
	// 建立通知服務和元件
	notificationService := services.NewNotificationService()
	notificationWidget := NewNotificationWidget(notificationService)
	
	// 確保元件被使用
	_ = notificationWidget
	
	status := services.SaveStatusInfo{
		IsSaving:     true,
		SaveProgress: 0.5,
		HasChanges:   true,
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		notificationService.UpdateSaveStatus("test-note", "test.md", status)
	}
}