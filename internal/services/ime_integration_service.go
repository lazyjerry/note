// Package services 包含應用程式的業務邏輯服務
// 本檔案實作輸入法整合服務，提供與系統輸入法的深度整合功能
package services

import (
	"fmt"                           // 格式化輸出
	"strings"                       // 字串處理
	"time"                          // 時間處理
)

// IMEIntegrationService 代表輸入法整合服務介面
// 提供與系統輸入法的深度整合功能，特別針對繁體中文輸入法優化
type IMEIntegrationService interface {
	// 輸入法狀態管理
	IsIMEActive() bool
	GetCurrentIME() string
	SwitchToChineseIME() error
	SwitchToEnglishIME() error
	
	// 組合文字處理
	StartComposition(text string) error
	UpdateComposition(text string) error
	EndComposition() error
	CancelComposition() error
	GetCompositionText() string
	
	// 候選字管理
	GetCandidates() []string
	SelectCandidate(index int) error
	NavigateCandidates(direction int) error
	
	// 輸入法設定
	GetIMESettings() IMESettings
	UpdateIMESettings(settings IMESettings) error
	
	// 事件處理
	RegisterCompositionHandler(handler CompositionHandler)
	RegisterCandidateHandler(handler CandidateHandler)
	UnregisterHandlers()
}

// IMESettings 代表輸入法設定
// 包含輸入法的各種配置選項
type IMESettings struct {
	PreferredIME        string            // 偏好的輸入法
	AutoSwitchEnabled   bool              // 是否啟用自動切換
	CandidateWindowSize int               // 候選字視窗大小
	CompositionFont     string            // 組合文字字型
	CompositionFontSize float32           // 組合文字字型大小
	ShowPinyin          bool              // 是否顯示拼音
	ShowZhuyin          bool              // 是否顯示注音
	CustomKeyBindings   map[string]string // 自訂按鍵綁定
}

// CompositionHandler 代表組合文字處理器
// 處理輸入法組合文字的各種事件
type CompositionHandler interface {
	OnCompositionStart(text string)
	OnCompositionUpdate(text string)
	OnCompositionEnd(text string)
	OnCompositionCancel()
}

// CandidateHandler 代表候選字處理器
// 處理輸入法候選字的各種事件
type CandidateHandler interface {
	OnCandidatesUpdate(candidates []string)
	OnCandidateSelect(index int, candidate string)
	OnCandidatePageChange(page int)
}

// IMEState 代表輸入法狀態
// 記錄輸入法的當前狀態資訊
type IMEState struct {
	IsActive           bool              // 是否啟用
	CurrentIME         string            // 當前輸入法
	IsComposing        bool              // 是否正在組合
	CompositionText    string            // 組合文字
	CandidateWords     []string          // 候選字列表
	SelectedCandidate  int               // 選中的候選字索引
	LastUpdateTime     time.Time         // 最後更新時間
}

// imeIntegrationServiceImpl 實作 IMEIntegrationService 介面
// 提供輸入法整合的具體實作
type imeIntegrationServiceImpl struct {
	// 輸入法狀態
	state              IMEState          // 當前狀態
	settings           IMESettings       // 輸入法設定
	
	// 事件處理器
	compositionHandler CompositionHandler // 組合文字處理器
	candidateHandler   CandidateHandler   // 候選字處理器
	
	// 內部狀態
	isInitialized      bool              // 是否已初始化
	supportedIMEs      []string          // 支援的輸入法列表
}

// NewIMEIntegrationService 建立新的輸入法整合服務實例
// 回傳：IMEIntegrationService 介面的實作實例
//
// 執行流程：
// 1. 建立服務實例並設定預設配置
// 2. 初始化輸入法狀態
// 3. 載入系統輸入法資訊
// 4. 設定預設的事件處理器
// 5. 回傳配置完成的服務實例
func NewIMEIntegrationService() IMEIntegrationService {
	service := &imeIntegrationServiceImpl{
		state: IMEState{
			IsActive:          false,
			CurrentIME:        "",
			IsComposing:       false,
			CompositionText:   "",
			CandidateWords:    []string{},
			SelectedCandidate: -1,
			LastUpdateTime:    time.Now(),
		},
		settings: IMESettings{
			PreferredIME:        "注音",
			AutoSwitchEnabled:   true,
			CandidateWindowSize: 10,
			CompositionFont:     "PingFang TC",
			CompositionFontSize: 14.0,
			ShowPinyin:          false,
			ShowZhuyin:          true,
			CustomKeyBindings:   make(map[string]string),
		},
		isInitialized: false,
		supportedIMEs: []string{"注音", "倉頡", "速成", "拼音", "英數"},
	}
	
	// 初始化服務
	service.initialize()
	
	return service
}

// initialize 初始化輸入法整合服務
// 設定系統輸入法整合和事件監聽
//
// 執行流程：
// 1. 檢測系統支援的輸入法
// 2. 設定預設輸入法
// 3. 初始化事件監聽
// 4. 載入用戶設定
func (imes *imeIntegrationServiceImpl) initialize() {
	// 檢測當前系統輸入法
	imes.detectCurrentIME()
	
	// 設定預設的按鍵綁定
	imes.setupDefaultKeyBindings()
	
	// 標記為已初始化
	imes.isInitialized = true
}

// detectCurrentIME 檢測當前系統輸入法
// 識別系統當前使用的輸入法
//
// 執行流程：
// 1. 查詢系統輸入法狀態
// 2. 識別當前啟用的輸入法
// 3. 更新內部狀態
func (imes *imeIntegrationServiceImpl) detectCurrentIME() {
	// 由於這是跨平台的實作，這裡使用模擬的檢測邏輯
	// 實際應用中需要調用系統 API 來檢測輸入法狀態
	
	// 模擬檢測結果
	imes.state.CurrentIME = "注音"
	imes.state.IsActive = true
	imes.state.LastUpdateTime = time.Now()
}

// setupDefaultKeyBindings 設定預設按鍵綁定
// 建立常用的中文輸入法按鍵綁定
//
// 執行流程：
// 1. 設定候選字選擇按鍵
// 2. 設定輸入法切換按鍵
// 3. 設定組合文字控制按鍵
func (imes *imeIntegrationServiceImpl) setupDefaultKeyBindings() {
	// 候選字選擇按鍵
	imes.settings.CustomKeyBindings["1"] = "select_candidate_1"
	imes.settings.CustomKeyBindings["2"] = "select_candidate_2"
	imes.settings.CustomKeyBindings["3"] = "select_candidate_3"
	imes.settings.CustomKeyBindings["4"] = "select_candidate_4"
	imes.settings.CustomKeyBindings["5"] = "select_candidate_5"
	
	// 候選字導航按鍵
	imes.settings.CustomKeyBindings["Up"] = "previous_candidate"
	imes.settings.CustomKeyBindings["Down"] = "next_candidate"
	imes.settings.CustomKeyBindings["Left"] = "previous_page"
	imes.settings.CustomKeyBindings["Right"] = "next_page"
	
	// 組合文字控制按鍵
	imes.settings.CustomKeyBindings["Enter"] = "commit_composition"
	imes.settings.CustomKeyBindings["Escape"] = "cancel_composition"
	imes.settings.CustomKeyBindings["Backspace"] = "delete_composition"
	
	// 輸入法切換按鍵
	imes.settings.CustomKeyBindings["Ctrl+Space"] = "toggle_ime"
	imes.settings.CustomKeyBindings["Shift"] = "switch_to_english"
}

// IsIMEActive 檢查輸入法是否啟用
// 回傳：輸入法是否處於啟用狀態
func (imes *imeIntegrationServiceImpl) IsIMEActive() bool {
	return imes.state.IsActive
}

// GetCurrentIME 取得當前輸入法
// 回傳：當前使用的輸入法名稱
func (imes *imeIntegrationServiceImpl) GetCurrentIME() string {
	return imes.state.CurrentIME
}

// SwitchToChineseIME 切換到中文輸入法
// 回傳：切換操作的錯誤狀態
//
// 執行流程：
// 1. 檢查是否已經是中文輸入法
// 2. 執行輸入法切換
// 3. 更新內部狀態
// 4. 觸發相關事件
func (imes *imeIntegrationServiceImpl) SwitchToChineseIME() error {
	// 如果已經是中文輸入法，直接回傳
	if imes.isChineseIME(imes.state.CurrentIME) {
		return nil
	}
	
	// 切換到偏好的中文輸入法
	targetIME := imes.settings.PreferredIME
	if targetIME == "" {
		targetIME = "注音" // 預設使用注音輸入法
	}
	
	// 執行切換（模擬）
	err := imes.switchIME(targetIME)
	if err != nil {
		return fmt.Errorf("切換到中文輸入法失敗: %v", err)
	}
	
	// 更新狀態
	imes.state.CurrentIME = targetIME
	imes.state.IsActive = true
	imes.state.LastUpdateTime = time.Now()
	
	return nil
}

// SwitchToEnglishIME 切換到英文輸入法
// 回傳：切換操作的錯誤狀態
//
// 執行流程：
// 1. 檢查是否已經是英文輸入法
// 2. 執行輸入法切換
// 3. 更新內部狀態
// 4. 觸發相關事件
func (imes *imeIntegrationServiceImpl) SwitchToEnglishIME() error {
	// 如果已經是英文輸入法，直接回傳
	if imes.state.CurrentIME == "英數" {
		return nil
	}
	
	// 執行切換到英文輸入法（模擬）
	err := imes.switchIME("英數")
	if err != nil {
		return fmt.Errorf("切換到英文輸入法失敗: %v", err)
	}
	
	// 更新狀態
	imes.state.CurrentIME = "英數"
	imes.state.IsActive = false
	imes.state.LastUpdateTime = time.Now()
	
	return nil
}

// isChineseIME 檢查是否為中文輸入法
// 參數：imeName（輸入法名稱）
// 回傳：是否為中文輸入法
func (imes *imeIntegrationServiceImpl) isChineseIME(imeName string) bool {
	chineseIMEs := []string{"注音", "倉頡", "速成", "拼音", "大易", "行列"}
	for _, chineseIME := range chineseIMEs {
		if imeName == chineseIME {
			return true
		}
	}
	return false
}

// switchIME 執行輸入法切換（內部方法）
// 參數：targetIME（目標輸入法）
// 回傳：切換操作的錯誤狀態
func (imes *imeIntegrationServiceImpl) switchIME(targetIME string) error {
	// 檢查目標輸入法是否支援
	if !imes.isSupportedIME(targetIME) {
		return fmt.Errorf("不支援的輸入法: %s", targetIME)
	}
	
	// 模擬輸入法切換
	// 實際應用中需要調用系統 API 來執行切換
	
	return nil
}

// isSupportedIME 檢查是否為支援的輸入法
// 參數：imeName（輸入法名稱）
// 回傳：是否支援該輸入法
func (imes *imeIntegrationServiceImpl) isSupportedIME(imeName string) bool {
	for _, supportedIME := range imes.supportedIMEs {
		if imeName == supportedIME {
			return true
		}
	}
	return false
}

// StartComposition 開始組合輸入
// 參數：text（初始組合文字）
// 回傳：操作錯誤狀態
//
// 執行流程：
// 1. 設定組合狀態
// 2. 初始化組合文字
// 3. 觸發組合開始事件
// 4. 更新候選字列表
func (imes *imeIntegrationServiceImpl) StartComposition(text string) error {
	// 設定組合狀態
	imes.state.IsComposing = true
	imes.state.CompositionText = text
	imes.state.LastUpdateTime = time.Now()
	
	// 觸發組合開始事件
	if imes.compositionHandler != nil {
		imes.compositionHandler.OnCompositionStart(text)
	}
	
	// 更新候選字
	imes.updateCandidatesForComposition(text)
	
	return nil
}

// UpdateComposition 更新組合文字
// 參數：text（新的組合文字）
// 回傳：操作錯誤狀態
//
// 執行流程：
// 1. 更新組合文字
// 2. 重新生成候選字
// 3. 觸發組合更新事件
func (imes *imeIntegrationServiceImpl) UpdateComposition(text string) error {
	if !imes.state.IsComposing {
		return fmt.Errorf("目前沒有進行組合輸入")
	}
	
	// 更新組合文字
	imes.state.CompositionText = text
	imes.state.LastUpdateTime = time.Now()
	
	// 觸發組合更新事件
	if imes.compositionHandler != nil {
		imes.compositionHandler.OnCompositionUpdate(text)
	}
	
	// 更新候選字
	imes.updateCandidatesForComposition(text)
	
	return nil
}

// EndComposition 結束組合輸入
// 回傳：操作錯誤狀態
//
// 執行流程：
// 1. 確認當前組合文字
// 2. 清理組合狀態
// 3. 觸發組合結束事件
func (imes *imeIntegrationServiceImpl) EndComposition() error {
	if !imes.state.IsComposing {
		return fmt.Errorf("目前沒有進行組合輸入")
	}
	
	// 保存組合文字
	compositionText := imes.state.CompositionText
	
	// 清理組合狀態
	imes.clearCompositionState()
	
	// 觸發組合結束事件
	if imes.compositionHandler != nil {
		imes.compositionHandler.OnCompositionEnd(compositionText)
	}
	
	return nil
}

// CancelComposition 取消組合輸入
// 回傳：操作錯誤狀態
//
// 執行流程：
// 1. 清理組合狀態
// 2. 觸發組合取消事件
func (imes *imeIntegrationServiceImpl) CancelComposition() error {
	if !imes.state.IsComposing {
		return fmt.Errorf("目前沒有進行組合輸入")
	}
	
	// 清理組合狀態
	imes.clearCompositionState()
	
	// 觸發組合取消事件
	if imes.compositionHandler != nil {
		imes.compositionHandler.OnCompositionCancel()
	}
	
	return nil
}

// clearCompositionState 清理組合狀態（內部方法）
// 重置所有與組合輸入相關的狀態
func (imes *imeIntegrationServiceImpl) clearCompositionState() {
	imes.state.IsComposing = false
	imes.state.CompositionText = ""
	imes.state.CandidateWords = []string{}
	imes.state.SelectedCandidate = -1
	imes.state.LastUpdateTime = time.Now()
}

// GetCompositionText 取得當前組合文字
// 回傳：當前的組合文字
func (imes *imeIntegrationServiceImpl) GetCompositionText() string {
	return imes.state.CompositionText
}

// GetCandidates 取得候選字列表
// 回傳：當前的候選字列表
func (imes *imeIntegrationServiceImpl) GetCandidates() []string {
	return imes.state.CandidateWords
}

// SelectCandidate 選擇候選字
// 參數：index（候選字索引）
// 回傳：操作錯誤狀態
//
// 執行流程：
// 1. 驗證索引有效性
// 2. 選擇指定的候選字
// 3. 觸發候選字選擇事件
// 4. 結束組合輸入
func (imes *imeIntegrationServiceImpl) SelectCandidate(index int) error {
	if index < 0 || index >= len(imes.state.CandidateWords) {
		return fmt.Errorf("無效的候選字索引: %d", index)
	}
	
	// 取得選擇的候選字
	selectedCandidate := imes.state.CandidateWords[index]
	
	// 更新選擇狀態
	imes.state.SelectedCandidate = index
	
	// 觸發候選字選擇事件
	if imes.candidateHandler != nil {
		imes.candidateHandler.OnCandidateSelect(index, selectedCandidate)
	}
	
	// 結束組合輸入
	return imes.EndComposition()
}

// NavigateCandidates 導航候選字列表
// 參數：direction（導航方向，-1為上，1為下）
// 回傳：操作錯誤狀態
//
// 執行流程：
// 1. 計算新的選擇索引
// 2. 處理邊界條件
// 3. 更新選擇狀態
func (imes *imeIntegrationServiceImpl) NavigateCandidates(direction int) error {
	if len(imes.state.CandidateWords) == 0 {
		return fmt.Errorf("沒有可用的候選字")
	}
	
	// 計算新的選擇索引
	newIndex := imes.state.SelectedCandidate + direction
	
	// 處理邊界條件
	if newIndex < 0 {
		newIndex = len(imes.state.CandidateWords) - 1
	} else if newIndex >= len(imes.state.CandidateWords) {
		newIndex = 0
	}
	
	// 更新選擇狀態
	imes.state.SelectedCandidate = newIndex
	
	return nil
}

// updateCandidatesForComposition 為組合文字更新候選字
// 參數：compositionText（組合文字）
//
// 執行流程：
// 1. 根據組合文字生成候選字
// 2. 更新候選字列表
// 3. 觸發候選字更新事件
func (imes *imeIntegrationServiceImpl) updateCandidatesForComposition(compositionText string) {
	// 根據當前輸入法生成候選字
	candidates := imes.generateCandidates(compositionText)
	
	// 更新候選字列表
	imes.state.CandidateWords = candidates
	imes.state.SelectedCandidate = 0 // 預設選擇第一個候選字
	
	// 觸發候選字更新事件
	if imes.candidateHandler != nil {
		imes.candidateHandler.OnCandidatesUpdate(candidates)
	}
}

// generateCandidates 生成候選字列表
// 參數：input（輸入文字）
// 回傳：候選字列表
//
// 執行流程：
// 1. 根據當前輸入法類型
// 2. 分析輸入內容
// 3. 生成對應的候選字
func (imes *imeIntegrationServiceImpl) generateCandidates(input string) []string {
	if input == "" {
		return []string{}
	}
	
	// 根據當前輸入法生成候選字
	switch imes.state.CurrentIME {
	case "注音":
		return imes.generateZhuyinCandidates(input)
	case "倉頡":
		return imes.generateCangjieCandidate(input)
	case "拼音":
		return imes.generatePinyinCandidates(input)
	default:
		return []string{}
	}
}

// generateZhuyinCandidates 生成注音候選字
// 參數：zhuyin（注音輸入）
// 回傳：候選字列表
func (imes *imeIntegrationServiceImpl) generateZhuyinCandidates(zhuyin string) []string {
	// 簡化的注音候選字對照表
	zhuyinMap := map[string][]string{
		"ㄋㄧˇ":     {"你", "妳", "尼", "泥"},
		"ㄏㄠˇ":     {"好", "號", "豪", "毫"},
		"ㄕˋ":      {"是", "事", "世", "勢"},
		"ㄐㄧㄝˋ":   {"界", "借", "戒", "介"},
		"ㄓㄨㄥ":    {"中", "鐘", "忠", "終"},
		"ㄨㄣˊ":     {"文", "聞", "溫", "紋"},
		"ㄧㄡˇ":     {"有", "友", "又", "右"},
		"ㄧˊ":      {"一", "以", "已", "意"},
		"ㄍㄜˋ":     {"個", "各", "格", "隔"},
	}
	
	if candidates, exists := zhuyinMap[zhuyin]; exists {
		return candidates
	}
	
	// 如果沒有完全匹配，嘗試部分匹配
	for key, candidates := range zhuyinMap {
		if strings.HasPrefix(key, zhuyin) && len(zhuyin) >= 3 {
			return candidates
		}
	}
	
	return []string{}
}

// generateCangjieCandidate 生成倉頡候選字
// 參數：cangjie（倉頡輸入）
// 回傳：候選字列表
func (imes *imeIntegrationServiceImpl) generateCangjieCandidate(cangjie string) []string {
	// 簡化的倉頡候選字對照表
	cangjieMap := map[string][]string{
		"人":   {"人", "入"},
		"人心": {"你", "妳"},
		"口":   {"口", "古"},
		"口口": {"呂", "品"},
		"手":   {"手", "才"},
		"木":   {"木", "本"},
	}
	
	if candidates, exists := cangjieMap[cangjie]; exists {
		return candidates
	}
	
	return []string{}
}

// generatePinyinCandidates 生成拼音候選字
// 參數：pinyin（拼音輸入）
// 回傳：候選字列表
func (imes *imeIntegrationServiceImpl) generatePinyinCandidates(pinyin string) []string {
	// 簡化的拼音候選字對照表
	pinyinMap := map[string][]string{
		"ni":   {"你", "妳", "尼", "泥"},
		"hao":  {"好", "號", "豪", "毫"},
		"shi":  {"是", "事", "世", "勢"},
		"jie":  {"界", "借", "戒", "介"},
		"zhong": {"中", "鐘", "忠", "終"},
		"wen":  {"文", "聞", "溫", "紋"},
	}
	
	if candidates, exists := pinyinMap[pinyin]; exists {
		return candidates
	}
	
	return []string{}
}

// GetIMESettings 取得輸入法設定
// 回傳：當前的輸入法設定
func (imes *imeIntegrationServiceImpl) GetIMESettings() IMESettings {
	return imes.settings
}

// UpdateIMESettings 更新輸入法設定
// 參數：settings（新的輸入法設定）
// 回傳：更新操作的錯誤狀態
//
// 執行流程：
// 1. 驗證設定的有效性
// 2. 更新內部設定
// 3. 應用新的設定
func (imes *imeIntegrationServiceImpl) UpdateIMESettings(settings IMESettings) error {
	// 驗證設定
	if err := imes.validateSettings(settings); err != nil {
		return fmt.Errorf("無效的輸入法設定: %v", err)
	}
	
	// 更新設定
	imes.settings = settings
	
	// 應用新設定
	imes.applySettings()
	
	return nil
}

// validateSettings 驗證輸入法設定
// 參數：settings（要驗證的設定）
// 回傳：驗證錯誤
func (imes *imeIntegrationServiceImpl) validateSettings(settings IMESettings) error {
	// 檢查偏好輸入法是否支援
	if settings.PreferredIME != "" && !imes.isSupportedIME(settings.PreferredIME) {
		return fmt.Errorf("不支援的偏好輸入法: %s", settings.PreferredIME)
	}
	
	// 檢查候選字視窗大小
	if settings.CandidateWindowSize < 1 || settings.CandidateWindowSize > 20 {
		return fmt.Errorf("候選字視窗大小必須在 1-20 之間")
	}
	
	// 檢查字型大小
	if settings.CompositionFontSize < 8.0 || settings.CompositionFontSize > 72.0 {
		return fmt.Errorf("組合文字字型大小必須在 8-72 之間")
	}
	
	return nil
}

// applySettings 應用輸入法設定
// 將新的設定應用到輸入法系統
func (imes *imeIntegrationServiceImpl) applySettings() {
	// 如果啟用自動切換且有偏好輸入法，嘗試切換
	if imes.settings.AutoSwitchEnabled && imes.settings.PreferredIME != "" {
		if imes.settings.PreferredIME != imes.state.CurrentIME {
			imes.switchIME(imes.settings.PreferredIME)
		}
	}
}

// RegisterCompositionHandler 註冊組合文字處理器
// 參數：handler（組合文字處理器）
func (imes *imeIntegrationServiceImpl) RegisterCompositionHandler(handler CompositionHandler) {
	imes.compositionHandler = handler
}

// RegisterCandidateHandler 註冊候選字處理器
// 參數：handler（候選字處理器）
func (imes *imeIntegrationServiceImpl) RegisterCandidateHandler(handler CandidateHandler) {
	imes.candidateHandler = handler
}

// UnregisterHandlers 取消註冊所有處理器
// 清理所有已註冊的事件處理器
func (imes *imeIntegrationServiceImpl) UnregisterHandlers() {
	imes.compositionHandler = nil
	imes.candidateHandler = nil
}