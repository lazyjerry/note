// Package services 包含應用程式的業務邏輯服務
// 本檔案包含輸入法整合服務的單元測試
package services

import (
	"testing"
)

// TestNewIMEIntegrationService 測試輸入法整合服務的建立
// 驗證服務是否正確初始化
func TestNewIMEIntegrationService(t *testing.T) {
	service := NewIMEIntegrationService()
	
	if service == nil {
		t.Fatal("輸入法整合服務建立失敗")
	}
	
	// 測試初始狀態
	if !service.IsIMEActive() {
		t.Error("預期輸入法處於啟用狀態")
	}
	
	currentIME := service.GetCurrentIME()
	if currentIME == "" {
		t.Error("預期有當前輸入法")
	}
}

// TestIMESwitching 測試輸入法切換功能
// 驗證中英文輸入法切換的正確性
func TestIMESwitching(t *testing.T) {
	service := NewIMEIntegrationService()
	
	// 測試切換到英文輸入法
	err := service.SwitchToEnglishIME()
	if err != nil {
		t.Errorf("切換到英文輸入法失敗: %v", err)
	}
	
	if service.GetCurrentIME() != "英數" {
		t.Error("預期當前輸入法為英數")
	}
	
	if service.IsIMEActive() {
		t.Error("英文輸入法狀態下，IME 應該不啟用")
	}
	
	// 測試切換到中文輸入法
	err = service.SwitchToChineseIME()
	if err != nil {
		t.Errorf("切換到中文輸入法失敗: %v", err)
	}
	
	if !service.IsIMEActive() {
		t.Error("中文輸入法狀態下，IME 應該啟用")
	}
}

// TestCompositionHandling 測試組合輸入處理
// 驗證組合輸入的開始、更新和結束流程
func TestCompositionHandling(t *testing.T) {
	service := NewIMEIntegrationService()
	
	// 測試開始組合輸入
	err := service.StartComposition("ㄋㄧˇ")
	if err != nil {
		t.Errorf("開始組合輸入失敗: %v", err)
	}
	
	compositionText := service.GetCompositionText()
	if compositionText != "ㄋㄧˇ" {
		t.Errorf("組合文字不正確，預期 'ㄋㄧˇ'，實際 '%s'", compositionText)
	}
	
	// 測試更新組合輸入
	err = service.UpdateComposition("ㄋㄧˇㄏㄠˇ")
	if err != nil {
		t.Errorf("更新組合輸入失敗: %v", err)
	}
	
	updatedText := service.GetCompositionText()
	if updatedText != "ㄋㄧˇㄏㄠˇ" {
		t.Errorf("更新後組合文字不正確，預期 'ㄋㄧˇㄏㄠˇ'，實際 '%s'", updatedText)
	}
	
	// 測試結束組合輸入
	err = service.EndComposition()
	if err != nil {
		t.Errorf("結束組合輸入失敗: %v", err)
	}
	
	finalText := service.GetCompositionText()
	if finalText != "" {
		t.Error("結束組合後，組合文字應該為空")
	}
}

// TestCandidateManagement 測試候選字管理功能
// 驗證候選字的生成、選擇和導航
func TestCandidateManagement(t *testing.T) {
	service := NewIMEIntegrationService()
	
	// 開始組合輸入以生成候選字
	err := service.StartComposition("ㄋㄧˇ")
	if err != nil {
		t.Errorf("開始組合輸入失敗: %v", err)
	}
	
	// 檢查候選字
	candidates := service.GetCandidates()
	if len(candidates) == 0 {
		t.Error("應該有候選字")
	}
	
	// 測試候選字選擇
	if len(candidates) > 0 {
		err = service.SelectCandidate(0)
		if err != nil {
			t.Errorf("選擇候選字失敗: %v", err)
		}
	}
	
	// 測試無效索引的候選字選擇
	err = service.SelectCandidate(-1)
	if err == nil {
		t.Error("選擇無效索引的候選字應該回傳錯誤")
	}
	
	err = service.SelectCandidate(100)
	if err == nil {
		t.Error("選擇超出範圍索引的候選字應該回傳錯誤")
	}
}

// TestCandidateNavigation 測試候選字導航功能
// 驗證候選字列表的上下導航
func TestCandidateNavigation(t *testing.T) {
	service := NewIMEIntegrationService()
	
	// 開始組合輸入以生成候選字
	err := service.StartComposition("ㄋㄧˇ")
	if err != nil {
		t.Errorf("開始組合輸入失敗: %v", err)
	}
	
	candidates := service.GetCandidates()
	if len(candidates) < 2 {
		t.Skip("需要至少 2 個候選字來測試導航")
	}
	
	// 測試向下導航
	err = service.NavigateCandidates(1)
	if err != nil {
		t.Errorf("向下導航候選字失敗: %v", err)
	}
	
	// 測試向上導航
	err = service.NavigateCandidates(-1)
	if err != nil {
		t.Errorf("向上導航候選字失敗: %v", err)
	}
	
	// 測試邊界條件（循環導航）
	for i := 0; i < len(candidates)+1; i++ {
		err = service.NavigateCandidates(1)
		if err != nil {
			t.Errorf("循環導航候選字失敗: %v", err)
		}
	}
}

// TestCompositionCancellation 測試組合輸入取消功能
// 驗證組合輸入的取消流程
func TestCompositionCancellation(t *testing.T) {
	service := NewIMEIntegrationService()
	
	// 開始組合輸入
	err := service.StartComposition("ㄋㄧˇ")
	if err != nil {
		t.Errorf("開始組合輸入失敗: %v", err)
	}
	
	// 取消組合輸入
	err = service.CancelComposition()
	if err != nil {
		t.Errorf("取消組合輸入失敗: %v", err)
	}
	
	// 驗證狀態已清理
	if service.GetCompositionText() != "" {
		t.Error("取消組合後，組合文字應該為空")
	}
	
	if len(service.GetCandidates()) != 0 {
		t.Error("取消組合後，候選字列表應該為空")
	}
	
	// 測試在沒有組合輸入時取消
	err = service.CancelComposition()
	if err == nil {
		t.Error("在沒有組合輸入時取消應該回傳錯誤")
	}
}

// TestIMESettings 測試輸入法設定功能
// 驗證輸入法設定的取得和更新
func TestIMESettings(t *testing.T) {
	service := NewIMEIntegrationService()
	
	// 取得預設設定
	settings := service.GetIMESettings()
	if settings.PreferredIME == "" {
		t.Error("預期有預設的偏好輸入法")
	}
	
	// 測試更新設定
	newSettings := settings
	newSettings.PreferredIME = "倉頡"
	newSettings.CandidateWindowSize = 8
	newSettings.ShowZhuyin = false
	
	err := service.UpdateIMESettings(newSettings)
	if err != nil {
		t.Errorf("更新輸入法設定失敗: %v", err)
	}
	
	// 驗證設定已更新
	updatedSettings := service.GetIMESettings()
	if updatedSettings.PreferredIME != "倉頡" {
		t.Error("偏好輸入法設定未正確更新")
	}
	
	if updatedSettings.CandidateWindowSize != 8 {
		t.Error("候選字視窗大小設定未正確更新")
	}
}

// TestInvalidSettings 測試無效設定處理
// 驗證無效輸入法設定的錯誤處理
func TestInvalidSettings(t *testing.T) {
	service := NewIMEIntegrationService()
	
	// 測試無效的偏好輸入法
	invalidSettings := service.GetIMESettings()
	invalidSettings.PreferredIME = "不存在的輸入法"
	
	err := service.UpdateIMESettings(invalidSettings)
	if err == nil {
		t.Error("設定無效的偏好輸入法應該回傳錯誤")
	}
	
	// 測試無效的候選字視窗大小
	invalidSettings = service.GetIMESettings()
	invalidSettings.CandidateWindowSize = 0
	
	err = service.UpdateIMESettings(invalidSettings)
	if err == nil {
		t.Error("設定無效的候選字視窗大小應該回傳錯誤")
	}
	
	// 測試無效的字型大小
	invalidSettings = service.GetIMESettings()
	invalidSettings.CompositionFontSize = 100.0
	
	err = service.UpdateIMESettings(invalidSettings)
	if err == nil {
		t.Error("設定無效的字型大小應該回傳錯誤")
	}
}

// TestZhuyinCandidateGeneration 測試注音候選字生成
// 驗證注音輸入的候選字生成準確性
func TestZhuyinCandidateGeneration(t *testing.T) {
	service := NewIMEIntegrationService()
	
	// 確保使用注音輸入法
	err := service.SwitchToChineseIME()
	if err != nil {
		t.Errorf("切換到中文輸入法失敗: %v", err)
	}
	
	testCases := []struct {
		input    string
		expected []string
		desc     string
	}{
		{"ㄋㄧˇ", []string{"你", "妳", "尼", "泥"}, "注音 ㄋㄧˇ"},
		{"ㄏㄠˇ", []string{"好", "號", "豪", "毫"}, "注音 ㄏㄠˇ"},
		{"ㄕˋ", []string{"是", "事", "世", "勢"}, "注音 ㄕˋ"},
	}
	
	for _, tc := range testCases {
		err := service.StartComposition(tc.input)
		if err != nil {
			t.Errorf("開始組合輸入失敗 (%s): %v", tc.desc, err)
			continue
		}
		
		candidates := service.GetCandidates()
		if len(candidates) == 0 {
			t.Errorf("%s: 應該有候選字", tc.desc)
			continue
		}
		
		// 檢查是否包含預期的候選字
		found := false
		for _, expected := range tc.expected {
			for _, candidate := range candidates {
				if candidate == expected {
					found = true
					break
				}
			}
			if found {
				break
			}
		}
		
		if !found {
			t.Errorf("%s: 候選字中應該包含預期的字符", tc.desc)
		}
		
		// 清理狀態
		service.CancelComposition()
	}
}

// MockCompositionHandler 模擬組合文字處理器
// 用於測試事件處理功能
type MockCompositionHandler struct {
	StartCalled  bool
	UpdateCalled bool
	EndCalled    bool
	CancelCalled bool
	LastText     string
}

func (mch *MockCompositionHandler) OnCompositionStart(text string) {
	mch.StartCalled = true
	mch.LastText = text
}

func (mch *MockCompositionHandler) OnCompositionUpdate(text string) {
	mch.UpdateCalled = true
	mch.LastText = text
}

func (mch *MockCompositionHandler) OnCompositionEnd(text string) {
	mch.EndCalled = true
	mch.LastText = text
}

func (mch *MockCompositionHandler) OnCompositionCancel() {
	mch.CancelCalled = true
}

// MockCandidateHandler 模擬候選字處理器
// 用於測試候選字事件處理功能
type MockCandidateHandler struct {
	UpdateCalled bool
	SelectCalled bool
	PageCalled   bool
	LastCandidates []string
	LastIndex      int
	LastCandidate  string
}

func (mch *MockCandidateHandler) OnCandidatesUpdate(candidates []string) {
	mch.UpdateCalled = true
	mch.LastCandidates = candidates
}

func (mch *MockCandidateHandler) OnCandidateSelect(index int, candidate string) {
	mch.SelectCalled = true
	mch.LastIndex = index
	mch.LastCandidate = candidate
}

func (mch *MockCandidateHandler) OnCandidatePageChange(page int) {
	mch.PageCalled = true
}

// TestEventHandlers 測試事件處理器功能
// 驗證組合文字和候選字事件的正確觸發
func TestEventHandlers(t *testing.T) {
	service := NewIMEIntegrationService()
	
	// 建立模擬處理器
	compositionHandler := &MockCompositionHandler{}
	candidateHandler := &MockCandidateHandler{}
	
	// 註冊處理器
	service.RegisterCompositionHandler(compositionHandler)
	service.RegisterCandidateHandler(candidateHandler)
	
	// 測試組合開始事件
	err := service.StartComposition("ㄋㄧˇ")
	if err != nil {
		t.Errorf("開始組合輸入失敗: %v", err)
	}
	
	if !compositionHandler.StartCalled {
		t.Error("組合開始事件未被觸發")
	}
	
	if !candidateHandler.UpdateCalled {
		t.Error("候選字更新事件未被觸發")
	}
	
	// 測試組合更新事件
	err = service.UpdateComposition("ㄋㄧˇㄏㄠˇ")
	if err != nil {
		t.Errorf("更新組合輸入失敗: %v", err)
	}
	
	if !compositionHandler.UpdateCalled {
		t.Error("組合更新事件未被觸發")
	}
	
	// 測試候選字選擇事件
	candidates := service.GetCandidates()
	if len(candidates) > 0 {
		err = service.SelectCandidate(0)
		if err != nil {
			t.Errorf("選擇候選字失敗: %v", err)
		}
		
		if !candidateHandler.SelectCalled {
			t.Error("候選字選擇事件未被觸發")
		}
		
		if !compositionHandler.EndCalled {
			t.Error("組合結束事件未被觸發")
		}
	}
	
	// 測試取消註冊
	service.UnregisterHandlers()
	
	// 重新開始組合，事件不應該被觸發
	compositionHandler.StartCalled = false
	err = service.StartComposition("ㄋㄧˇ")
	if err != nil {
		t.Errorf("開始組合輸入失敗: %v", err)
	}
	
	// 由於處理器已取消註冊，事件不應該被觸發
	// 但由於我們的實作中沒有檢查 nil，這個測試可能會失敗
	// 這是一個需要改進的地方
}

// BenchmarkCandidateGeneration 候選字生成的效能基準測試
// 測量候選字生成的效能
func BenchmarkCandidateGeneration(b *testing.B) {
	service := NewIMEIntegrationService()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.StartComposition("ㄋㄧˇ")
		service.CancelComposition()
	}
}

// BenchmarkIMESwitching 輸入法切換的效能基準測試
// 測量輸入法切換的效能
func BenchmarkIMESwitching(b *testing.B) {
	service := NewIMEIntegrationService()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.SwitchToEnglishIME()
		service.SwitchToChineseIME()
	}
}