package ui

import (
	"fmt"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"mac-notebook-app/internal/services"
)

// NotificationWidget 通知顯示元件
// 負責在 UI 中顯示通知訊息和保存狀態指示器
type NotificationWidget struct {
	widget.BaseWidget
	
	// 服務依賴
	notificationService services.NotificationService
	
	// UI 元件
	container           *fyne.Container
	notificationList    *fyne.Container
	saveStatusContainer *fyne.Container
	saveStatusLabel     *widget.Label
	saveProgressBar     *widget.ProgressBar
	
	// 狀態管理
	activeNotifications map[string]*NotificationItem
	currentSaveStatus   *services.SaveStatusInfo
}

// NotificationItem 單個通知項目
type NotificationItem struct {
	container     *fyne.Container
	titleLabel    *widget.Label
	messageLabel  *widget.Label
	closeButton   *widget.Button
	notification  *services.Notification
}

// NewNotificationWidget 建立新的通知元件
// 參數：
//   - notificationService: 通知服務實例
// 回傳：NotificationWidget 指標
//
// 執行流程：
// 1. 初始化元件結構
// 2. 建立 UI 佈局
// 3. 設定服務回調
// 4. 回傳元件實例
func NewNotificationWidget(notificationService services.NotificationService) *NotificationWidget {
	widget := &NotificationWidget{
		notificationService: notificationService,
		activeNotifications: make(map[string]*NotificationItem),
	}
	
	widget.ExtendBaseWidget(widget)
	widget.createUI()
	widget.setupCallbacks()
	
	return widget
}

// CreateRenderer 建立元件的渲染器
// 回傳：fyne.WidgetRenderer 介面實例
func (w *NotificationWidget) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(w.container)
}

// createUI 建立通知元件的 UI 佈局
// 執行流程：
// 1. 建立通知列表容器
// 2. 建立保存狀態指示器
// 3. 組合主容器佈局
func (w *NotificationWidget) createUI() {
	// 建立通知列表容器
	w.notificationList = container.NewVBox()
	
	// 建立保存狀態指示器
	w.saveStatusLabel = widget.NewLabel("就緒")
	w.saveStatusLabel.TextStyle = fyne.TextStyle{Italic: true}
	
	w.saveProgressBar = widget.NewProgressBar()
	w.saveProgressBar.Hide() // 初始隱藏
	
	w.saveStatusContainer = container.NewHBox(
		widget.NewIcon(theme.DocumentSaveIcon()),
		w.saveStatusLabel,
		w.saveProgressBar,
	)
	
	// 組合主容器
	w.container = container.NewVBox(
		w.notificationList,
		widget.NewSeparator(),
		w.saveStatusContainer,
	)
}

// setupCallbacks 設定通知服務的回調函數
// 執行流程：
// 1. 設定通知更新回調
// 2. 設定保存狀態更新回調
func (w *NotificationWidget) setupCallbacks() {
	// 設定通知回調
	w.notificationService.SetNotificationCallback(w.onNotificationUpdate)
	
	// 設定保存狀態回調
	w.notificationService.SetSaveStatusCallback(w.onSaveStatusUpdate)
}

// onNotificationUpdate 處理通知更新回調
// 參數：
//   - notification: 更新的通知實例
//
// 執行流程：
// 1. 檢查通知是否已讀（表示要移除）
// 2. 如果是新通知，建立通知項目
// 3. 如果是移除通知，從 UI 中移除
func (w *NotificationWidget) onNotificationUpdate(notification *services.Notification) {
	if notification == nil {
		return
	}
	
	if notification.IsRead {
		// 移除通知
		w.removeNotificationItem(notification.ID)
	} else {
		// 添加或更新通知
		w.addNotificationItem(notification)
	}
}

// addNotificationItem 添加通知項目到 UI
// 參數：
//   - notification: 要添加的通知
//
// 執行流程：
// 1. 檢查通知是否已存在
// 2. 建立通知項目 UI
// 3. 添加到通知列表
// 4. 刷新 UI
func (w *NotificationWidget) addNotificationItem(notification *services.Notification) {
	// 檢查是否已存在
	if _, exists := w.activeNotifications[notification.ID]; exists {
		return
	}
	
	// 建立通知項目
	item := w.createNotificationItem(notification)
	w.activeNotifications[notification.ID] = item
	
	// 添加到列表頂部（最新的在上面）
	w.notificationList.Add(item.container)
	w.notificationList.Refresh()
}

// removeNotificationItem 從 UI 移除通知項目
// 參數：
//   - notificationID: 要移除的通知 ID
//
// 執行流程：
// 1. 查找通知項目
// 2. 從容器中移除
// 3. 從映射表中刪除
// 4. 刷新 UI
func (w *NotificationWidget) removeNotificationItem(notificationID string) {
	item, exists := w.activeNotifications[notificationID]
	if !exists {
		return
	}
	
	// 從容器中移除
	w.notificationList.Remove(item.container)
	
	// 從映射表中刪除
	delete(w.activeNotifications, notificationID)
	
	// 刷新 UI
	w.notificationList.Refresh()
}

// createNotificationItem 建立單個通知項目的 UI
// 參數：
//   - notification: 通知資料
// 回傳：NotificationItem 指標
//
// 執行流程：
// 1. 建立標題和內容標籤
// 2. 建立關閉按鈕
// 3. 設定通知樣式（顏色）
// 4. 組合佈局
func (w *NotificationWidget) createNotificationItem(notification *services.Notification) *NotificationItem {
	// 建立標題標籤
	titleLabel := widget.NewLabel(notification.Title)
	titleLabel.TextStyle = fyne.TextStyle{Bold: true}
	
	// 建立內容標籤
	messageLabel := widget.NewLabel(notification.Message)
	messageLabel.Wrapping = fyne.TextWrapWord
	
	// 建立關閉按鈕
	closeButton := widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
		w.notificationService.DismissNotification(notification.ID)
	})
	closeButton.Importance = widget.LowImportance
	
	// 建立內容容器
	contentContainer := container.NewVBox(titleLabel, messageLabel)
	
	// 建立主容器
	mainContainer := container.NewBorder(
		nil, nil, nil, closeButton,
		contentContainer,
	)
	
	// 設定通知樣式
	w.applyNotificationStyle(mainContainer, notification.Type)
	
	return &NotificationItem{
		container:    mainContainer,
		titleLabel:   titleLabel,
		messageLabel: messageLabel,
		closeButton:  closeButton,
		notification: notification,
	}
}

// applyNotificationStyle 應用通知類型對應的樣式
// 參數：
//   - container: 要應用樣式的容器
//   - notificationType: 通知類型
//
// 執行流程：
// 1. 根據通知類型選擇顏色
// 2. 設定容器背景色
// 3. 添加邊框效果
func (w *NotificationWidget) applyNotificationStyle(container *fyne.Container, notificationType services.NotificationType) {
	// 根據通知類型設定顏色
	var bgColor color.Color
	switch notificationType {
	case services.NotificationInfo:
		bgColor = color.RGBA{33, 150, 243, 50} // 淺藍色
	case services.NotificationSuccess:
		bgColor = color.RGBA{76, 175, 80, 50} // 淺綠色
	case services.NotificationWarning:
		bgColor = color.RGBA{255, 152, 0, 50} // 淺橙色
	case services.NotificationError:
		bgColor = color.RGBA{244, 67, 54, 50} // 淺紅色
	default:
		bgColor = color.RGBA{117, 117, 117, 50} // 淺灰色
	}
	
	// 建立背景矩形
	background := widget.NewCard("", "", container)
	background.SetContent(container)
	
	// 注意：Fyne 的 Card 元件會自動提供背景色和邊框
	// 這裡我們使用 Card 來提供視覺效果
	// bgColor 變數保留供未來擴展使用
	_ = bgColor
}

// onSaveStatusUpdate 處理保存狀態更新回調
// 參數：
//   - noteID: 筆記 ID
//   - status: 保存狀態資訊（nil 表示清除狀態）
//
// 執行流程：
// 1. 更新當前保存狀態
// 2. 更新狀態標籤文字
// 3. 更新進度條顯示
// 4. 刷新 UI
func (w *NotificationWidget) onSaveStatusUpdate(noteID string, status *services.SaveStatusInfo) {
	w.currentSaveStatus = status
	
	if status == nil {
		// 清除狀態
		w.saveStatusLabel.SetText("就緒")
		w.saveProgressBar.Hide()
	} else {
		// 更新狀態
		w.updateSaveStatusDisplay(status)
	}
	
	w.saveStatusContainer.Refresh()
}

// updateSaveStatusDisplay 更新保存狀態顯示
// 參數：
//   - status: 保存狀態資訊
//
// 執行流程：
// 1. 根據保存狀態設定標籤文字
// 2. 更新進度條值和可見性
// 3. 設定狀態指示顏色
func (w *NotificationWidget) updateSaveStatusDisplay(status *services.SaveStatusInfo) {
	if status.IsSaving {
		// 正在保存
		w.saveStatusLabel.SetText("正在保存 " + status.FileName + "...")
		w.saveProgressBar.SetValue(status.SaveProgress)
		w.saveProgressBar.Show()
	} else if status.HasChanges {
		// 有未保存的變更
		w.saveStatusLabel.SetText(status.FileName + " (有未保存的變更)")
		w.saveProgressBar.Hide()
	} else {
		// 已保存
		lastSavedText := "已保存"
		if !status.LastSaved.IsZero() {
			lastSavedText += " (" + formatTimeAgo(status.LastSaved) + ")"
		}
		w.saveStatusLabel.SetText(lastSavedText)
		w.saveProgressBar.Hide()
	}
}

// ShowTestNotifications 顯示測試通知（用於開發和測試）
// 執行流程：
// 1. 顯示各種類型的測試通知
// 2. 模擬保存狀態更新
func (w *NotificationWidget) ShowTestNotifications() {
	// 顯示各種類型的通知
	w.notificationService.ShowInfo("資訊通知", "這是一個資訊通知的範例")
	w.notificationService.ShowSuccess("操作成功", "檔案已成功保存")
	w.notificationService.ShowWarning("注意事項", "檔案可能包含未保存的變更")
	w.notificationService.ShowError("錯誤發生", "無法連接到伺服器")
	
	// 模擬保存狀態
	go func() {
		time.Sleep(1 * time.Second)
		w.notificationService.UpdateSaveStatus("test-note", "test.md", services.SaveStatusInfo{
			IsSaving:     true,
			SaveProgress: 0.3,
			HasChanges:   true,
		})
		
		time.Sleep(2 * time.Second)
		w.notificationService.UpdateSaveStatus("test-note", "test.md", services.SaveStatusInfo{
			IsSaving:     true,
			SaveProgress: 0.8,
			HasChanges:   true,
		})
		
		time.Sleep(1 * time.Second)
		w.notificationService.UpdateSaveStatus("test-note", "test.md", services.SaveStatusInfo{
			IsSaving:     false,
			SaveProgress: 1.0,
			HasChanges:   false,
			LastSaved:    time.Now(),
		})
	}()
}

// ClearAllNotifications 清除所有通知
func (w *NotificationWidget) ClearAllNotifications() {
	w.notificationService.DismissAllNotifications()
}

// GetNotificationCount 取得當前通知數量
// 回傳：活躍通知的數量
func (w *NotificationWidget) GetNotificationCount() int {
	return len(w.activeNotifications)
}

// formatTimeAgo 格式化時間為「多久之前」的形式
// 參數：
//   - t: 要格式化的時間
// 回傳：格式化後的時間字串
func formatTimeAgo(t time.Time) string {
	duration := time.Since(t)
	
	if duration < time.Minute {
		return "剛剛"
	} else if duration < time.Hour {
		minutes := int(duration.Minutes())
		return fmt.Sprintf("%d 分鐘前", minutes)
	} else if duration < 24*time.Hour {
		hours := int(duration.Hours())
		return fmt.Sprintf("%d 小時前", hours)
	} else {
		days := int(duration.Hours() / 24)
		return fmt.Sprintf("%d 天前", days)
	}
}

// SetVisible 設定通知元件的可見性
// 參數：
//   - visible: 是否可見
func (w *NotificationWidget) SetVisible(visible bool) {
	if visible {
		w.container.Show()
	} else {
		w.container.Hide()
	}
}