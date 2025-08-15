package services

import (
	"sync"
	"testing"
	"time"
)

// TestNewNotificationService 測試通知服務的建立
// 驗證服務能夠正確初始化
func TestNewNotificationService(t *testing.T) {
	service := NewNotificationService()
	
	if service == nil {
		t.Error("通知服務不應該為空")
	}

	// 驗證初始狀態
	notifications := service.GetActiveNotifications()
	if len(notifications) != 0 {
		t.Errorf("初始通知數量應該為 0，得到 %d", len(notifications))
	}
}

// TestShowNotification 測試顯示通知功能
// 驗證各種類型的通知能夠正確顯示
func TestShowNotification(t *testing.T) {
	service := NewNotificationService()

	tests := []struct {
		name         string
		notifyType   NotificationType
		title        string
		message      string
		duration     time.Duration
		expectedType NotificationType
	}{
		{
			name:         "資訊通知",
			notifyType:   NotificationInfo,
			title:        "資訊標題",
			message:      "資訊內容",
			duration:     3 * time.Second,
			expectedType: NotificationInfo,
		},
		{
			name:         "成功通知",
			notifyType:   NotificationSuccess,
			title:        "成功標題",
			message:      "操作成功",
			duration:     3 * time.Second,
			expectedType: NotificationSuccess,
		},
		{
			name:         "警告通知",
			notifyType:   NotificationWarning,
			title:        "警告標題",
			message:      "警告內容",
			duration:     4 * time.Second,
			expectedType: NotificationWarning,
		},
		{
			name:         "錯誤通知",
			notifyType:   NotificationError,
			title:        "錯誤標題",
			message:      "錯誤內容",
			duration:     5 * time.Second,
			expectedType: NotificationError,
		},
		{
			name:         "持久通知",
			notifyType:   NotificationInfo,
			title:        "持久標題",
			message:      "持久內容",
			duration:     0, // 持久通知
			expectedType: NotificationInfo,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			notificationID := service.ShowNotification(tt.notifyType, tt.title, tt.message, tt.duration)
			
			if notificationID == "" {
				t.Error("通知 ID 不應該為空")
			}

			// 驗證通知是否被正確加入
			notifications := service.GetActiveNotifications()
			found := false
			for _, notification := range notifications {
				if notification.ID == notificationID {
					found = true
					if notification.Type != tt.expectedType {
						t.Errorf("通知類型 = %v, 期望 %v", notification.Type, tt.expectedType)
					}
					if notification.Title != tt.title {
						t.Errorf("通知標題 = %v, 期望 %v", notification.Title, tt.title)
					}
					if notification.Message != tt.message {
						t.Errorf("通知內容 = %v, 期望 %v", notification.Message, tt.message)
					}
					if notification.Duration != tt.duration {
						t.Errorf("通知持續時間 = %v, 期望 %v", notification.Duration, tt.duration)
					}
					if notification.IsPersistent != (tt.duration <= 0) {
						t.Errorf("持久標誌 = %v, 期望 %v", notification.IsPersistent, tt.duration <= 0)
					}
					break
				}
			}
			
			if !found {
				t.Error("通知應該被加入到活躍通知列表中")
			}
		})
	}
}

// TestConvenienceMethods 測試便利方法
// 驗證 ShowSuccess、ShowError、ShowWarning、ShowInfo 方法
func TestConvenienceMethods(t *testing.T) {
	service := NewNotificationService()

	tests := []struct {
		name         string
		method       func(string, string) string
		expectedType NotificationType
	}{
		{
			name:         "ShowSuccess",
			method:       service.ShowSuccess,
			expectedType: NotificationSuccess,
		},
		{
			name:         "ShowError",
			method:       service.ShowError,
			expectedType: NotificationError,
		},
		{
			name:         "ShowWarning",
			method:       service.ShowWarning,
			expectedType: NotificationWarning,
		},
		{
			name:         "ShowInfo",
			method:       service.ShowInfo,
			expectedType: NotificationInfo,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			title := tt.name + " 標題"
			message := tt.name + " 內容"
			
			notificationID := tt.method(title, message)
			
			if notificationID == "" {
				t.Error("通知 ID 不應該為空")
			}

			// 驗證通知類型
			notifications := service.GetActiveNotifications()
			found := false
			for _, notification := range notifications {
				if notification.ID == notificationID {
					found = true
					if notification.Type != tt.expectedType {
						t.Errorf("通知類型 = %v, 期望 %v", notification.Type, tt.expectedType)
					}
					break
				}
			}
			
			if !found {
				t.Error("通知應該被加入到活躍通知列表中")
			}
		})
	}
}

// TestDismissNotification 測試關閉通知功能
// 驗證能夠正確關閉指定的通知
func TestDismissNotification(t *testing.T) {
	service := NewNotificationService()

	// 建立測試通知
	notificationID := service.ShowInfo("測試標題", "測試內容")
	
	// 驗證通知存在
	notifications := service.GetActiveNotifications()
	if len(notifications) != 1 {
		t.Errorf("應該有 1 個通知，得到 %d", len(notifications))
	}

	// 關閉通知
	success := service.DismissNotification(notificationID)
	if !success {
		t.Error("關閉通知應該成功")
	}

	// 驗證通知已被移除
	notifications = service.GetActiveNotifications()
	if len(notifications) != 0 {
		t.Errorf("通知應該被移除，但還有 %d 個", len(notifications))
	}

	// 嘗試關閉不存在的通知
	success = service.DismissNotification("不存在的ID")
	if success {
		t.Error("關閉不存在的通知應該失敗")
	}
}

// TestDismissAllNotifications 測試關閉所有通知功能
// 驗證能夠一次關閉所有活躍通知
func TestDismissAllNotifications(t *testing.T) {
	service := NewNotificationService()

	// 建立多個測試通知
	service.ShowInfo("通知 1", "內容 1")
	service.ShowSuccess("通知 2", "內容 2")
	service.ShowWarning("通知 3", "內容 3")
	service.ShowError("通知 4", "內容 4")

	// 驗證通知存在
	notifications := service.GetActiveNotifications()
	if len(notifications) != 4 {
		t.Errorf("應該有 4 個通知，得到 %d", len(notifications))
	}

	// 關閉所有通知
	service.DismissAllNotifications()

	// 驗證所有通知已被移除
	notifications = service.GetActiveNotifications()
	if len(notifications) != 0 {
		t.Errorf("所有通知應該被移除，但還有 %d 個", len(notifications))
	}
}

// TestGetActiveNotifications 測試取得活躍通知功能
// 驗證通知列表的排序和內容正確性
func TestGetActiveNotifications(t *testing.T) {
	service := NewNotificationService()

	// 建立測試通知（有時間間隔）
	id1 := service.ShowInfo("通知 1", "內容 1")
	time.Sleep(10 * time.Millisecond)
	id2 := service.ShowSuccess("通知 2", "內容 2")
	time.Sleep(10 * time.Millisecond)
	id3 := service.ShowWarning("通知 3", "內容 3")

	notifications := service.GetActiveNotifications()
	
	// 驗證通知數量
	if len(notifications) != 3 {
		t.Errorf("應該有 3 個通知，得到 %d", len(notifications))
	}

	// 驗證排序（最新的在前）
	if notifications[0].ID != id3 {
		t.Error("最新的通知應該在第一位")
	}
	if notifications[1].ID != id2 {
		t.Error("第二新的通知應該在第二位")
	}
	if notifications[2].ID != id1 {
		t.Error("最舊的通知應該在最後一位")
	}
}

// TestSaveStatus 測試保存狀態功能
// 驗證保存狀態的更新、取得和清除
func TestSaveStatus(t *testing.T) {
	service := NewNotificationService()

	noteID := "test-note-123"
	fileName := "test.md"
	
	// 測試初始狀態
	status := service.GetSaveStatus(noteID)
	if status != nil {
		t.Error("初始狀態應該為空")
	}

	// 更新保存狀態
	saveStatus := SaveStatusInfo{
		NoteID:       noteID,
		FileName:     fileName,
		IsSaving:     true,
		LastSaved:    time.Now(),
		SaveProgress: 0.5,
		HasChanges:   true,
	}
	
	service.UpdateSaveStatus(noteID, fileName, saveStatus)

	// 驗證狀態更新
	status = service.GetSaveStatus(noteID)
	if status == nil {
		t.Error("保存狀態不應該為空")
	}
	if status.NoteID != noteID {
		t.Errorf("筆記 ID = %v, 期望 %v", status.NoteID, noteID)
	}
	if status.FileName != fileName {
		t.Errorf("檔案名稱 = %v, 期望 %v", status.FileName, fileName)
	}
	if status.IsSaving != true {
		t.Errorf("保存狀態 = %v, 期望 %v", status.IsSaving, true)
	}
	if status.SaveProgress != 0.5 {
		t.Errorf("保存進度 = %v, 期望 %v", status.SaveProgress, 0.5)
	}
	if status.HasChanges != true {
		t.Errorf("變更狀態 = %v, 期望 %v", status.HasChanges, true)
	}

	// 清除保存狀態
	service.ClearSaveStatus(noteID)
	status = service.GetSaveStatus(noteID)
	if status != nil {
		t.Error("清除後的狀態應該為空")
	}
}

// TestCallbacks 測試回調函數功能
// 驗證通知和保存狀態的回調機制
func TestCallbacks(t *testing.T) {
	service := NewNotificationService()

	// 測試通知回調
	var callbackNotification *Notification
	var callbackMutex sync.Mutex
	
	service.SetNotificationCallback(func(notification *Notification) {
		callbackMutex.Lock()
		defer callbackMutex.Unlock()
		callbackNotification = notification
	})

	// 顯示通知並驗證回調
	notificationID := service.ShowInfo("回調測試", "回調內容")
	
	// 等待回調執行
	time.Sleep(10 * time.Millisecond)
	
	callbackMutex.Lock()
	if callbackNotification == nil {
		t.Error("通知回調應該被觸發")
	} else if callbackNotification.ID != notificationID {
		t.Error("回調通知 ID 不匹配")
	}
	callbackMutex.Unlock()

	// 測試保存狀態回調
	var callbackNoteID string
	var callbackSaveStatus *SaveStatusInfo
	var saveCallbackMutex sync.Mutex
	
	service.SetSaveStatusCallback(func(noteID string, status *SaveStatusInfo) {
		saveCallbackMutex.Lock()
		defer saveCallbackMutex.Unlock()
		callbackNoteID = noteID
		callbackSaveStatus = status
	})

	// 更新保存狀態並驗證回調
	testNoteID := "callback-test-note"
	testFileName := "callback-test.md"
	testStatus := SaveStatusInfo{
		IsSaving:     true,
		SaveProgress: 0.8,
		HasChanges:   false,
	}
	
	service.UpdateSaveStatus(testNoteID, testFileName, testStatus)
	
	// 等待回調執行
	time.Sleep(10 * time.Millisecond)
	
	saveCallbackMutex.Lock()
	if callbackNoteID != testNoteID {
		t.Errorf("回調筆記 ID = %v, 期望 %v", callbackNoteID, testNoteID)
	}
	if callbackSaveStatus == nil {
		t.Error("保存狀態回調應該被觸發")
	} else if callbackSaveStatus.SaveProgress != 0.8 {
		t.Errorf("回調保存進度 = %v, 期望 %v", callbackSaveStatus.SaveProgress, 0.8)
	}
	saveCallbackMutex.Unlock()
}

// TestAutoNotificationDismissal 測試自動通知消失功能
// 驗證非持久通知能夠在指定時間後自動消失
func TestAutoNotificationDismissal(t *testing.T) {
	service := NewNotificationService()

	// 建立短時間的通知
	notificationID := service.ShowNotification(NotificationInfo, "自動消失測試", "內容", 100*time.Millisecond)
	
	// 驗證通知存在
	notifications := service.GetActiveNotifications()
	if len(notifications) != 1 {
		t.Errorf("應該有 1 個通知，得到 %d", len(notifications))
	}
	
	// 驗證通知 ID 正確
	if len(notifications) > 0 && notifications[0].ID != notificationID {
		t.Error("通知 ID 不匹配")
	}

	// 等待通知自動消失
	time.Sleep(150 * time.Millisecond)
	
	// 驗證通知已自動消失
	notifications = service.GetActiveNotifications()
	if len(notifications) != 0 {
		t.Errorf("通知應該自動消失，但還有 %d 個", len(notifications))
	}

	// 測試持久通知不會自動消失
	persistentID := service.ShowNotification(NotificationInfo, "持久測試", "內容", 0)
	
	// 等待一段時間
	time.Sleep(100 * time.Millisecond)
	
	// 驗證持久通知仍然存在
	notifications = service.GetActiveNotifications()
	if len(notifications) != 1 {
		t.Errorf("持久通知應該仍然存在，得到 %d 個", len(notifications))
	}
	
	// 驗證持久通知的 ID 正確
	if len(notifications) > 0 && notifications[0].ID != persistentID {
		t.Error("持久通知 ID 不匹配")
	}
	
	// 手動清理持久通知
	service.DismissNotification(persistentID)
}

// TestNotificationConcurrentAccess 測試並發存取安全性
// 驗證多個 goroutine 同時操作通知服務的安全性
func TestNotificationConcurrentAccess(t *testing.T) {
	service := NewNotificationService()
	
	var wg sync.WaitGroup
	numGoroutines := 5
	notificationsPerGoroutine := 2

	// 並發建立通知
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()
			for j := 0; j < notificationsPerGoroutine; j++ {
				service.ShowNotification(NotificationInfo, "並發測試", "測試內容", 0) // 使用持久通知避免自動消失
			}
		}(i)
	}

	wg.Wait()

	// 驗證通知被正確建立（至少有一些通知）
	notifications := service.GetActiveNotifications()
	if len(notifications) == 0 {
		t.Error("應該有通知被建立")
	}

	// 測試並發關閉不會造成 panic
	var dismissWg sync.WaitGroup
	for i := 0; i < 3; i++ {
		dismissWg.Add(1)
		go func() {
			defer dismissWg.Done()
			service.DismissAllNotifications()
		}()
	}

	dismissWg.Wait()

	// 驗證最終狀態一致
	finalNotifications := service.GetActiveNotifications()
	if len(finalNotifications) != 0 {
		t.Errorf("所有通知應該被清除，但還有 %d 個", len(finalNotifications))
	}
}

// TestUtilityFunctions 測試工具函數
// 驗證通知類型字串和顏色函數
func TestUtilityFunctions(t *testing.T) {
	tests := []struct {
		notificationType NotificationType
		expectedString   string
		expectedColor    string
	}{
		{NotificationInfo, "資訊", "#2196F3"},
		{NotificationSuccess, "成功", "#4CAF50"},
		{NotificationWarning, "警告", "#FF9800"},
		{NotificationError, "錯誤", "#F44336"},
		{NotificationType(999), "未知", "#757575"}, // 無效類型
	}

	for _, tt := range tests {
		t.Run(tt.expectedString, func(t *testing.T) {
			typeString := GetNotificationTypeString(tt.notificationType)
			if typeString != tt.expectedString {
				t.Errorf("GetNotificationTypeString() = %v, 期望 %v", typeString, tt.expectedString)
			}

			typeColor := GetNotificationTypeColor(tt.notificationType)
			if typeColor != tt.expectedColor {
				t.Errorf("GetNotificationTypeColor() = %v, 期望 %v", typeColor, tt.expectedColor)
			}
		})
	}
}

// BenchmarkShowNotification 效能測試 - 顯示通知
// 測試顯示通知功能的效能表現
func BenchmarkShowNotification(b *testing.B) {
	service := NewNotificationService()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.ShowInfo("效能測試", "測試內容")
	}
}

// BenchmarkGetActiveNotifications 效能測試 - 取得活躍通知
// 測試取得通知列表功能的效能表現
func BenchmarkGetActiveNotifications(b *testing.B) {
	service := NewNotificationService()
	
	// 預先建立一些通知
	for i := 0; i < 100; i++ {
		service.ShowInfo("測試通知", "測試內容")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.GetActiveNotifications()
	}
}