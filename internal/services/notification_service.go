package services

import (
	"sync"
	"time"
)



// notificationService 通知服務的具體實作
type notificationService struct {
	notifications        map[string]*Notification    // 活躍通知的映射表
	saveStatuses         map[string]*SaveStatusInfo  // 保存狀態的映射表
	notificationCallback func(*Notification)         // 通知更新回調
	saveStatusCallback   func(string, *SaveStatusInfo) // 保存狀態更新回調
	mutex               sync.RWMutex                 // 讀寫鎖，保護並發存取
	idCounter           int                          // 通知 ID 計數器
}

// NewNotificationService 建立新的通知服務實例
// 回傳：NotificationService 介面實例
//
// 執行流程：
// 1. 初始化通知和狀態映射表
// 2. 設定並發安全的讀寫鎖
// 3. 初始化 ID 計數器
func NewNotificationService() NotificationService {
	return &notificationService{
		notifications: make(map[string]*Notification),
		saveStatuses:  make(map[string]*SaveStatusInfo),
		idCounter:     0,
	}
}

// ShowNotification 顯示通知訊息
// 參數：
//   - notificationType: 通知類型
//   - title: 通知標題
//   - message: 通知內容
//   - duration: 顯示持續時間
// 回傳：通知 ID
//
// 執行流程：
// 1. 生成唯一的通知 ID
// 2. 建立通知實例
// 3. 加入活躍通知列表
// 4. 觸發 UI 回調
// 5. 設定自動消失計時器（如果不是持久通知）
func (s *notificationService) ShowNotification(notificationType NotificationType, title, message string, duration time.Duration) string {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// 生成通知 ID
	s.idCounter++
	notificationID := generateNotificationID(s.idCounter)

	// 建立通知實例
	notification := &Notification{
		ID:           notificationID,
		Type:         notificationType,
		Title:        title,
		Message:      message,
		Duration:     duration,
		CreatedAt:    time.Now(),
		IsRead:       false,
		IsPersistent: duration <= 0, // 持續時間為 0 或負數表示持久通知
	}

	// 加入活躍通知列表
	s.notifications[notificationID] = notification

	// 觸發 UI 回調
	if s.notificationCallback != nil {
		s.notificationCallback(notification)
	}

	// 設定自動消失計時器（非持久通知）
	if !notification.IsPersistent {
		go s.scheduleNotificationDismissal(notificationID, duration)
	}

	return notificationID
}

// ShowSuccess 顯示成功通知
// 參數：
//   - title: 通知標題
//   - message: 通知內容
// 回傳：通知 ID
func (s *notificationService) ShowSuccess(title, message string) string {
	return s.ShowNotification(NotificationSuccess, title, message, 3*time.Second)
}

// ShowError 顯示錯誤通知
// 參數：
//   - title: 通知標題
//   - message: 通知內容
// 回傳：通知 ID
func (s *notificationService) ShowError(title, message string) string {
	return s.ShowNotification(NotificationError, title, message, 5*time.Second)
}

// ShowWarning 顯示警告通知
// 參數：
//   - title: 通知標題
//   - message: 通知內容
// 回傳：通知 ID
func (s *notificationService) ShowWarning(title, message string) string {
	return s.ShowNotification(NotificationWarning, title, message, 4*time.Second)
}

// ShowInfo 顯示資訊通知
// 參數：
//   - title: 通知標題
//   - message: 通知內容
// 回傳：通知 ID
func (s *notificationService) ShowInfo(title, message string) string {
	return s.ShowNotification(NotificationInfo, title, message, 3*time.Second)
}

// DismissNotification 關閉指定的通知
// 參數：
//   - notificationID: 要關閉的通知 ID
// 回傳：是否成功關閉
//
// 執行流程：
// 1. 檢查通知是否存在
// 2. 從活躍通知列表中移除
// 3. 觸發 UI 回調（傳入 nil 表示移除）
func (s *notificationService) DismissNotification(notificationID string) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if notification, exists := s.notifications[notificationID]; exists {
		delete(s.notifications, notificationID)
		
		// 觸發 UI 回調表示通知被移除
		if s.notificationCallback != nil {
			notification.IsRead = true
			s.notificationCallback(notification)
		}
		
		return true
	}
	return false
}

// DismissAllNotifications 關閉所有通知
// 執行流程：
// 1. 遍歷所有活躍通知
// 2. 逐一觸發移除回調
// 3. 清空通知列表
func (s *notificationService) DismissAllNotifications() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// 觸發所有通知的移除回調
	if s.notificationCallback != nil {
		for _, notification := range s.notifications {
			notification.IsRead = true
			s.notificationCallback(notification)
		}
	}

	// 清空通知列表
	s.notifications = make(map[string]*Notification)
}

// GetActiveNotifications 取得所有活躍的通知
// 回傳：活躍通知的陣列
//
// 執行流程：
// 1. 建立通知陣列
// 2. 複製所有活躍通知
// 3. 按建立時間排序（最新的在前）
func (s *notificationService) GetActiveNotifications() []*Notification {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	notifications := make([]*Notification, 0, len(s.notifications))
	for _, notification := range s.notifications {
		notifications = append(notifications, notification)
	}

	// 按建立時間排序（最新的在前）
	for i := 0; i < len(notifications)-1; i++ {
		for j := i + 1; j < len(notifications); j++ {
			if notifications[i].CreatedAt.Before(notifications[j].CreatedAt) {
				notifications[i], notifications[j] = notifications[j], notifications[i]
			}
		}
	}

	return notifications
}

// UpdateSaveStatus 更新保存狀態指示器
// 參數：
//   - noteID: 筆記 ID
//   - fileName: 檔案名稱
//   - status: 保存狀態資訊
//
// 執行流程：
// 1. 更新保存狀態映射表
// 2. 觸發狀態更新回調
func (s *notificationService) UpdateSaveStatus(noteID, fileName string, status SaveStatusInfo) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// 確保 NoteID 和 FileName 正確設定
	status.NoteID = noteID
	status.FileName = fileName

	// 更新狀態
	s.saveStatuses[noteID] = &status

	// 觸發 UI 回調
	if s.saveStatusCallback != nil {
		s.saveStatusCallback(noteID, &status)
	}
}

// GetSaveStatus 取得指定筆記的保存狀態
// 參數：
//   - noteID: 筆記 ID
// 回傳：保存狀態資訊，如果不存在則回傳 nil
func (s *notificationService) GetSaveStatus(noteID string) *SaveStatusInfo {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if status, exists := s.saveStatuses[noteID]; exists {
		// 回傳狀態的副本，避免外部修改
		statusCopy := *status
		return &statusCopy
	}
	return nil
}

// ClearSaveStatus 清除指定筆記的保存狀態
// 參數：
//   - noteID: 筆記 ID
//
// 執行流程：
// 1. 從狀態映射表中移除
// 2. 觸發清除回調（傳入 nil）
func (s *notificationService) ClearSaveStatus(noteID string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.saveStatuses[noteID]; exists {
		delete(s.saveStatuses, noteID)
		
		// 觸發 UI 回調表示狀態被清除
		if s.saveStatusCallback != nil {
			s.saveStatusCallback(noteID, nil)
		}
	}
}

// SetNotificationCallback 設定通知回調函數（用於 UI 更新）
// 參數：
//   - callback: 通知更新時的回調函數
func (s *notificationService) SetNotificationCallback(callback func(*Notification)) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.notificationCallback = callback
}

// SetSaveStatusCallback 設定保存狀態回調函數（用於 UI 更新）
// 參數：
//   - callback: 保存狀態更新時的回調函數
func (s *notificationService) SetSaveStatusCallback(callback func(string, *SaveStatusInfo)) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.saveStatusCallback = callback
}

// scheduleNotificationDismissal 安排通知的自動消失
// 參數：
//   - notificationID: 通知 ID
//   - duration: 延遲時間
//
// 執行流程：
// 1. 等待指定的持續時間
// 2. 自動關閉通知
func (s *notificationService) scheduleNotificationDismissal(notificationID string, duration time.Duration) {
	time.Sleep(duration)
	s.DismissNotification(notificationID)
}

// generateNotificationID 生成通知 ID
// 參數：
//   - counter: 計數器值
// 回傳：格式化的通知 ID
func generateNotificationID(counter int) string {
	return "notification_" + time.Now().Format("20060102_150405") + "_" + string(rune('0'+counter%10))
}

// GetNotificationTypeString 取得通知類型的字串表示
// 參數：
//   - notificationType: 通知類型
// 回傳：類型的繁體中文字串
func GetNotificationTypeString(notificationType NotificationType) string {
	switch notificationType {
	case NotificationInfo:
		return "資訊"
	case NotificationSuccess:
		return "成功"
	case NotificationWarning:
		return "警告"
	case NotificationError:
		return "錯誤"
	default:
		return "未知"
	}
}

// GetNotificationTypeColor 取得通知類型對應的顏色
// 參數：
//   - notificationType: 通知類型
// 回傳：顏色的十六進位字串
func GetNotificationTypeColor(notificationType NotificationType) string {
	switch notificationType {
	case NotificationInfo:
		return "#2196F3" // 藍色
	case NotificationSuccess:
		return "#4CAF50" // 綠色
	case NotificationWarning:
		return "#FF9800" // 橙色
	case NotificationError:
		return "#F44336" // 紅色
	default:
		return "#757575" // 灰色
	}
}