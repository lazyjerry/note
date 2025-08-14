// Package services 實作應用程式的業務邏輯服務
// 本檔案包含自動保存服務的實作，負責管理筆記的定時自動保存功能
package services

import (
	"fmt"                              // 格式化輸出
	"mac-notebook-app/internal/models" // 引入資料模型
	"sync"                             // 同步原語套件
	"time"                             // 時間處理套件
)

// AutoSaveServiceImpl 實作 AutoSaveService 介面
// 負責管理多個筆記的自動保存功能，包含定時器管理和狀態追蹤
// 支援加密檔案的自動保存，並提供可配置的保存間隔
type AutoSaveServiceImpl struct {
	editorService    EditorService                // 編輯器服務，用於執行實際的保存操作
	settingsService  SettingsService              // 設定服務，用於取得自動保存配置
	timers           map[string]*time.Timer       // 儲存每個筆記的定時器
	saveStatus       map[string]*SaveStatus       // 儲存每個筆記的保存狀態
	mutex            sync.RWMutex                 // 讀寫鎖，保護並發存取
	notes            map[string]*models.Note      // 儲存筆記實例的快取
	defaultInterval  time.Duration                // 預設自動保存間隔
	encryptedBackoff time.Duration                // 加密檔案的額外延遲（避免頻繁加密操作）
}

// NewAutoSaveService 建立新的自動保存服務實例
// 參數：
//   - editorService: 編輯器服務實例，用於執行保存操作
//   - settingsService: 設定服務實例，用於取得自動保存配置
// 回傳：自動保存服務實例
//
// 執行流程：
// 1. 初始化所有內部映射表
// 2. 設定編輯器服務和設定服務依賴
// 3. 設定預設的自動保存間隔和加密檔案延遲
// 4. 初始化讀寫鎖
func NewAutoSaveService(editorService EditorService, settingsService SettingsService) *AutoSaveServiceImpl {
	return &AutoSaveServiceImpl{
		editorService:    editorService,
		settingsService:  settingsService,
		timers:           make(map[string]*time.Timer),
		saveStatus:       make(map[string]*SaveStatus),
		notes:            make(map[string]*models.Note),
		mutex:            sync.RWMutex{},
		defaultInterval:  5 * time.Minute,  // 預設 5 分鐘間隔
		encryptedBackoff: 30 * time.Second, // 加密檔案額外延遲 30 秒
	}
}

// NewAutoSaveServiceWithDefaults 建立使用預設設定的自動保存服務實例
// 參數：
//   - editorService: 編輯器服務實例，用於執行保存操作
// 回傳：自動保存服務實例
//
// 這個方法用於向後相容，當沒有設定服務時使用預設配置
func NewAutoSaveServiceWithDefaults(editorService EditorService) *AutoSaveServiceImpl {
	return &AutoSaveServiceImpl{
		editorService:    editorService,
		settingsService:  nil, // 沒有設定服務，使用預設值
		timers:           make(map[string]*time.Timer),
		saveStatus:       make(map[string]*SaveStatus),
		notes:            make(map[string]*models.Note),
		mutex:            sync.RWMutex{},
		defaultInterval:  5 * time.Minute,  // 預設 5 分鐘間隔
		encryptedBackoff: 30 * time.Second, // 加密檔案額外延遲 30 秒
	}
}

// StartAutoSave 為指定筆記啟動自動保存功能
// 參數：
//   - note: 要自動保存的筆記實例
//   - interval: 自動保存的時間間隔
//
// 執行流程：
// 1. 取得寫入鎖以確保執行緒安全
// 2. 停止該筆記現有的自動保存定時器（如果存在）
// 3. 儲存筆記實例到快取中
// 4. 初始化該筆記的保存狀態
// 5. 建立新的定時器並設定回調函數
// 6. 將定時器儲存到映射表中
func (a *AutoSaveServiceImpl) StartAutoSave(note *models.Note, interval time.Duration) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	// 如果該筆記已有定時器，先停止它
	if timer, exists := a.timers[note.ID]; exists {
		timer.Stop()
	}

	// 儲存筆記實例
	a.notes[note.ID] = note

	// 初始化保存狀態
	a.saveStatus[note.ID] = &SaveStatus{
		NoteID:    note.ID,
		IsSaving:  false,
		LastSaved: note.LastSaved,
		LastError: nil,
		SaveCount: 0,
	}

	// 建立定時器，定期執行自動保存
	timer := time.AfterFunc(interval, func() {
		a.performAutoSave(note.ID, interval)
	})

	// 儲存定時器
	a.timers[note.ID] = timer
}

// StopAutoSave 停止指定筆記的自動保存功能
// 參數：
//   - noteID: 要停止自動保存的筆記 ID
//
// 執行流程：
// 1. 取得寫入鎖以確保執行緒安全
// 2. 檢查該筆記是否有活躍的定時器
// 3. 停止定時器並從映射表中移除
// 4. 清理相關的狀態資訊和快取
func (a *AutoSaveServiceImpl) StopAutoSave(noteID string) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	// 停止定時器
	if timer, exists := a.timers[noteID]; exists {
		timer.Stop()
		delete(a.timers, noteID)
	}

	// 清理狀態和快取
	delete(a.saveStatus, noteID)
	delete(a.notes, noteID)
}

// SaveNow 立即保存指定筆記，不等待定時器觸發
// 參數：
//   - noteID: 要立即保存的筆記 ID
// 回傳：保存操作的錯誤（如果有）
//
// 執行流程：
// 1. 取得讀取鎖檢查筆記是否存在
// 2. 檢查筆記是否正在保存中，避免重複保存
// 3. 呼叫內部保存方法執行實際保存
func (a *AutoSaveServiceImpl) SaveNow(noteID string) error {
	a.mutex.RLock()
	note, noteExists := a.notes[noteID]
	status, statusExists := a.saveStatus[noteID]
	a.mutex.RUnlock()

	if !noteExists || !statusExists {
		return models.NewAppError("NOTE_NOT_FOUND", "找不到指定的筆記", "筆記 ID: "+noteID)
	}

	// 檢查是否正在保存中
	if status.IsSaving {
		return models.NewAppError("SAVE_IN_PROGRESS", "筆記正在保存中，請稍後再試", "筆記 ID: "+noteID)
	}

	return a.saveNoteWithRetry(note)
}

// GetSaveStatus 取得指定筆記的保存狀態資訊
// 參數：
//   - noteID: 要查詢狀態的筆記 ID
// 回傳：保存狀態資訊
//
// 執行流程：
// 1. 取得讀取鎖以安全存取狀態資訊
// 2. 檢查筆記是否存在於狀態映射表中
// 3. 回傳狀態資訊的副本，避免外部修改
func (a *AutoSaveServiceImpl) GetSaveStatus(noteID string) SaveStatus {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	if status, exists := a.saveStatus[noteID]; exists {
		// 回傳狀態的副本，避免外部修改
		return SaveStatus{
			NoteID:    status.NoteID,
			IsSaving:  status.IsSaving,
			LastSaved: status.LastSaved,
			LastError: status.LastError,
			SaveCount: status.SaveCount,
		}
	}

	// 如果找不到狀態，回傳預設狀態
	return SaveStatus{
		NoteID:    noteID,
		IsSaving:  false,
		LastSaved: time.Time{},
		LastError: models.NewAppError("STATUS_NOT_FOUND", "找不到筆記的保存狀態", "筆記 ID: "+noteID),
		SaveCount: 0,
	}
}

// performAutoSave 執行自動保存操作的內部方法
// 參數：
//   - noteID: 要保存的筆記 ID
//   - interval: 自動保存間隔，用於重新設定下次保存
//
// 執行流程：
// 1. 取得筆記實例和狀態資訊
// 2. 檢查筆記是否需要保存（是否已修改）
// 3. 執行保存操作並更新狀態
// 4. 重新設定下次自動保存的定時器
// 5. 處理保存過程中的任何錯誤
func (a *AutoSaveServiceImpl) performAutoSave(noteID string, interval time.Duration) {
	a.mutex.Lock()
	note, noteExists := a.notes[noteID]
	_, statusExists := a.saveStatus[noteID]
	a.mutex.Unlock()

	if !noteExists || !statusExists {
		return // 筆記已被移除，不需要保存
	}

	// 檢查筆記是否需要保存
	if !note.IsModified() {
		// 筆記未修改，重新設定定時器並返回
		a.rescheduleTimer(noteID, interval)
		return
	}

	// 執行保存操作，對加密檔案進行特殊處理
	err := a.saveNoteWithRetry(note)

	// 更新狀態
	a.mutex.Lock()
	if status, exists := a.saveStatus[noteID]; exists {
		if err != nil {
			status.LastError = err
			// 如果是加密檔案保存失敗，記錄特殊錯誤類型
			if note.IsEncrypted {
				status.LastError = models.NewAppError("ENCRYPTED_SAVE_FAILED", 
					"加密檔案自動保存失敗", err.Error())
			}
		} else {
			status.LastError = nil
			status.LastSaved = time.Now()
			status.SaveCount++
		}
	}
	a.mutex.Unlock()

	// 重新設定下次自動保存
	a.rescheduleTimer(noteID, interval)
}

// saveNoteWithRetry 執行帶重試機制的筆記保存操作
// 參數：
//   - note: 要保存的筆記實例
// 回傳：保存操作的錯誤（如果有）
//
// 執行流程：
// 1. 對於加密檔案，實作重試機制以處理可能的加密失敗
// 2. 對於一般檔案，直接呼叫 saveNote
// 3. 記錄重試次數和失敗原因
func (a *AutoSaveServiceImpl) saveNoteWithRetry(note *models.Note) error {
	const maxRetries = 3
	const retryDelay = 1 * time.Second

	var lastErr error
	
	for attempt := 1; attempt <= maxRetries; attempt++ {
		err := a.saveNote(note)
		if err == nil {
			// 保存成功
			return nil
		}
		
		lastErr = err
		
		// 如果不是加密檔案，或者已經是最後一次嘗試，不再重試
		if !note.IsEncrypted || attempt == maxRetries {
			break
		}
		
		// 等待一段時間後重試（只對加密檔案）
		time.Sleep(retryDelay * time.Duration(attempt))
	}
	
	// 所有重試都失敗，回傳最後的錯誤
	if note.IsEncrypted {
		return models.NewAppError("ENCRYPTED_SAVE_RETRY_FAILED", 
			"加密檔案保存重試失敗", 
			fmt.Sprintf("嘗試 %d 次後仍然失敗: %v", maxRetries, lastErr))
	}
	
	return lastErr
}

// saveNote 執行實際的筆記保存操作
// 參數：
//   - note: 要保存的筆記實例
// 回傳：保存操作的錯誤（如果有）
//
// 執行流程：
// 1. 更新保存狀態為進行中
// 2. 呼叫編輯器服務執行保存
// 3. 根據保存結果更新筆記的保存時間戳
// 4. 更新保存狀態為完成
func (a *AutoSaveServiceImpl) saveNote(note *models.Note) error {
	// 設定保存狀態為進行中
	a.mutex.Lock()
	if status, exists := a.saveStatus[note.ID]; exists {
		status.IsSaving = true
	}
	a.mutex.Unlock()

	// 執行保存操作
	err := a.editorService.SaveNote(note)

	// 更新保存狀態
	a.mutex.Lock()
	if status, exists := a.saveStatus[note.ID]; exists {
		status.IsSaving = false
		if err == nil {
			// 保存成功，更新筆記的保存時間戳和狀態
			note.MarkSaved()
			status.LastSaved = time.Now()
			status.SaveCount++
			status.LastError = nil
		} else {
			status.LastError = err
		}
	}
	a.mutex.Unlock()

	return err
}

// rescheduleTimer 重新設定指定筆記的自動保存定時器
// 參數：
//   - noteID: 筆記 ID
//   - interval: 自動保存間隔
//
// 執行流程：
// 1. 取得寫入鎖以安全操作定時器
// 2. 停止現有的定時器
// 3. 建立新的定時器並設定回調函數
// 4. 將新定時器儲存到映射表中
func (a *AutoSaveServiceImpl) rescheduleTimer(noteID string, interval time.Duration) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	// 停止現有定時器
	if timer, exists := a.timers[noteID]; exists {
		timer.Stop()
	}

	// 建立新定時器
	timer := time.AfterFunc(interval, func() {
		a.performAutoSave(noteID, interval)
	})

	// 儲存新定時器
	a.timers[noteID] = timer
}

// Shutdown 關閉自動保存服務，停止所有定時器並清理資源
// 這個方法應該在應用程式關閉時呼叫
//
// 執行流程：
// 1. 取得寫入鎖以確保執行緒安全
// 2. 停止所有活躍的定時器
// 3. 清理所有內部映射表
// 4. 釋放所有資源
func (a *AutoSaveServiceImpl) Shutdown() {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	// 停止所有定時器
	for _, timer := range a.timers {
		timer.Stop()
	}

	// 清理所有映射表
	a.timers = make(map[string]*time.Timer)
	a.saveStatus = make(map[string]*SaveStatus)
	a.notes = make(map[string]*models.Note)
}

// GetAllSaveStatuses 取得所有筆記的保存狀態
// 回傳：包含所有筆記保存狀態的映射表
//
// 執行流程：
// 1. 取得讀取鎖以安全存取狀態資訊
// 2. 建立狀態副本的映射表
// 3. 回傳狀態副本，避免外部修改原始資料
func (a *AutoSaveServiceImpl) GetAllSaveStatuses() map[string]SaveStatus {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	result := make(map[string]SaveStatus)
	for noteID, status := range a.saveStatus {
		result[noteID] = SaveStatus{
			NoteID:    status.NoteID,
			IsSaving:  status.IsSaving,
			LastSaved: status.LastSaved,
			LastError: status.LastError,
			SaveCount: status.SaveCount,
		}
	}

	return result
}

// IsAutoSaveActive 檢查指定筆記是否啟用了自動保存
// 參數：
//   - noteID: 要檢查的筆記 ID
// 回傳：如果自動保存已啟用則回傳 true，否則回傳 false
//
// 執行流程：
// 1. 取得讀取鎖以安全存取定時器資訊
// 2. 檢查該筆記是否存在活躍的定時器
func (a *AutoSaveServiceImpl) IsAutoSaveActive(noteID string) bool {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	_, exists := a.timers[noteID]
	return exists
}

// StartAutoSaveWithSettings 使用設定服務的間隔啟動自動保存
// 參數：
//   - note: 要自動保存的筆記實例
//
// 執行流程：
// 1. 從設定服務取得自動保存間隔
// 2. 根據筆記是否加密調整間隔
// 3. 呼叫 StartAutoSave 啟動自動保存
func (a *AutoSaveServiceImpl) StartAutoSaveWithSettings(note *models.Note) {
	interval := a.getAutoSaveInterval(note)
	a.StartAutoSave(note, interval)
}

// getAutoSaveInterval 取得指定筆記的自動保存間隔
// 參數：
//   - note: 筆記實例
// 回傳：自動保存間隔
//
// 執行流程：
// 1. 嘗試從設定服務取得使用者配置的間隔
// 2. 如果筆記是加密的，增加額外延遲以減少加密操作頻率
// 3. 如果無法取得設定，使用預設間隔
func (a *AutoSaveServiceImpl) getAutoSaveInterval(note *models.Note) time.Duration {
	var baseInterval time.Duration = a.defaultInterval

	// 嘗試從設定服務取得使用者配置
	if a.settingsService != nil {
		if settings, err := a.settingsService.LoadSettings(); err == nil {
			if settings.AutoSaveInterval > 0 {
				baseInterval = time.Duration(settings.AutoSaveInterval) * time.Minute
			}
		}
	}

	// 如果是加密檔案，增加額外延遲以減少加密操作的頻率
	if note.IsEncrypted {
		return baseInterval + a.encryptedBackoff
	}

	return baseInterval
}

// UpdateAutoSaveInterval 更新指定筆記的自動保存間隔
// 參數：
//   - noteID: 筆記 ID
//   - newInterval: 新的自動保存間隔
// 回傳：操作錯誤（如果有）
//
// 執行流程：
// 1. 檢查筆記是否存在且啟用了自動保存
// 2. 停止現有的定時器
// 3. 使用新間隔重新啟動自動保存
func (a *AutoSaveServiceImpl) UpdateAutoSaveInterval(noteID string, newInterval time.Duration) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	_, noteExists := a.notes[noteID]
	if !noteExists {
		return models.NewAppError("NOTE_NOT_FOUND", "找不到指定的筆記", "筆記 ID: "+noteID)
	}

	// 停止現有定時器
	if timer, exists := a.timers[noteID]; exists {
		timer.Stop()
		delete(a.timers, noteID)
	}

	// 使用新間隔重新啟動
	timer := time.AfterFunc(newInterval, func() {
		a.performAutoSave(noteID, newInterval)
	})

	a.timers[noteID] = timer
	return nil
}

// GetEncryptedFileCount 取得目前正在自動保存的加密檔案數量
// 回傳：加密檔案數量
//
// 執行流程：
// 1. 遍歷所有快取的筆記
// 2. 計算其中加密檔案的數量
func (a *AutoSaveServiceImpl) GetEncryptedFileCount() int {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	count := 0
	for _, note := range a.notes {
		if note.IsEncrypted {
			count++
		}
	}
	return count
}

// SetEncryptedBackoff 設定加密檔案的額外延遲時間
// 參數：
//   - backoff: 額外延遲時間
//
// 這個方法允許動態調整加密檔案的保存頻率，
// 在系統負載較高時可以增加延遲以減少加密操作的影響
func (a *AutoSaveServiceImpl) SetEncryptedBackoff(backoff time.Duration) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	a.encryptedBackoff = backoff
}