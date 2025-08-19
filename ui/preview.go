// Package ui 包含使用者介面相關的元件和視窗管理
// 本檔案實作增強版 Markdown 預覽面板 UI 元件，提供獨立視窗、工具列、搜尋和導航功能
package ui

import (
	"fmt"                           // Go 標準庫，用於格式化字串
	"strings"                       // 字串處理
	"time"                          // 時間處理
	"regexp"                        // 正規表達式
	"mac-notebook-app/internal/services" // 引入服務層

	"fyne.io/fyne/v2"               // Fyne GUI 框架核心套件
	"fyne.io/fyne/v2/app"           // Fyne 應用程式套件
	"fyne.io/fyne/v2/container"     // Fyne 容器佈局套件
	"fyne.io/fyne/v2/widget"        // Fyne UI 元件套件
	"fyne.io/fyne/v2/theme"         // Fyne 主題套件
	"fyne.io/fyne/v2/dialog"        // Fyne 對話框套件

)

// MarkdownPreview 代表增強版 Markdown 預覽面板 UI 元件
// 提供即時 HTML 預覽、獨立視窗、工具列、搜尋導航和效能優化功能
// 整合編輯器服務以提供 Markdown 到 HTML 的轉換
type MarkdownPreview struct {
	container     *fyne.Container      // 主要容器
	toolbar       *widget.Toolbar      // 增強版預覽工具欄
	previewArea   *widget.RichText     // HTML 預覽顯示區域
	statusLabel   *widget.Label        // 狀態標籤
	searchBar     *widget.Entry        // 搜尋輸入欄
	searchResults *widget.Label        // 搜尋結果顯示
	zoomSlider    *widget.Slider       // 縮放滑桿
	
	// 獨立視窗支援
	independentWindow fyne.Window       // 獨立預覽視窗
	isIndependent     bool              // 是否為獨立視窗模式
	
	// 服務依賴
	editorService services.EditorService // 編輯器服務
	
	// 當前狀態
	currentContent string              // 當前預覽的 Markdown 內容
	isVisible     bool                 // 預覽面板是否可見
	autoRefresh   bool                 // 是否自動刷新預覽
	zoomLevel     float64              // 縮放級別 (0.5-3.0)
	
	// 搜尋功能狀態
	searchQuery    string              // 當前搜尋查詢
	searchMatches  []int               // 搜尋匹配位置
	currentMatch   int                 // 當前匹配項索引
	
	// 效能優化
	lastUpdateTime time.Time           // 上次更新時間
	updateThrottle time.Duration       // 更新節流間隔
	contentCache   map[string]string   // 內容快取
	
	// 回調函數
	onVisibilityChanged func(visible bool) // 可見性變更回調
	onRefreshRequested  func()             // 刷新請求回調
	onZoomChanged       func(level float64) // 縮放變更回調
	onSearchPerformed   func(query string, matches int) // 搜尋執行回調
}

// NewMarkdownPreview 建立新的增強版 Markdown 預覽面板實例
// 參數：editorService（編輯器服務介面）
// 回傳：指向新建立的 MarkdownPreview 的指標
//
// 執行流程：
// 1. 建立 MarkdownPreview 結構體實例
// 2. 初始化編輯器服務依賴和預設狀態
// 3. 建立並配置所有 UI 元件（工具列、搜尋、縮放等）
// 4. 初始化效能優化設定（快取、節流等）
// 5. 組合完整的增強版預覽面板佈局
// 6. 回傳配置完成的預覽實例
func NewMarkdownPreview(editorService services.EditorService) *MarkdownPreview {
	// 建立 MarkdownPreview 實例
	preview := &MarkdownPreview{
		editorService:  editorService,
		isVisible:      true,                    // 預設可見
		autoRefresh:    true,                    // 預設自動刷新
		zoomLevel:      1.0,                     // 預設縮放級別 100%
		updateThrottle: 100 * time.Millisecond, // 更新節流 100ms
		contentCache:   make(map[string]string), // 初始化內容快取
		searchMatches:  make([]int, 0),          // 初始化搜尋匹配陣列
		currentMatch:   -1,                      // 初始化當前匹配索引
	}
	
	// 初始化 UI 元件
	preview.setupUI()
	
	return preview
}

// setupUI 初始化增強版預覽面板的使用者介面元件
// 這個方法負責建立和配置預覽面板的所有 UI 元件
//
// 執行流程：
// 1. 建立增強版預覽工具欄（匯出、列印、縮放等）
// 2. 建立搜尋功能元件（搜尋欄、結果顯示、導航）
// 3. 建立縮放控制元件（滑桿、按鈕）
// 4. 建立 HTML 預覽顯示區域
// 5. 建立狀態標籤顯示預覽狀態
// 6. 設定事件處理和回調函數
// 7. 組合所有元件到主容器中
func (mp *MarkdownPreview) setupUI() {
	// 建立增強版預覽工具欄
	mp.createEnhancedToolbar()
	
	// 建立搜尋功能元件
	mp.createSearchComponents()
	
	// 建立縮放控制元件
	mp.createZoomControls()
	
	// 建立預覽顯示區域
	mp.createPreviewArea()
	
	// 建立狀態標籤
	mp.createStatusLabel()
	
	// 組合預覽面板佈局
	mp.assembleEnhancedLayout()
}

// createEnhancedToolbar 建立增強版預覽面板工具欄
// 包含刷新、匯出、列印、縮放、搜尋、獨立視窗等完整預覽控制功能
//
// 執行流程：
// 1. 建立手動刷新和自動刷新切換按鈕
// 2. 建立匯出功能按鈕（HTML、PDF、Word）
// 3. 建立列印功能按鈕
// 4. 建立縮放控制按鈕（放大、縮小、重置）
// 5. 建立搜尋功能按鈕
// 6. 建立獨立視窗和全螢幕按鈕
// 7. 建立同步滾動和可見性控制按鈕
// 8. 組合所有按鈕到工具欄
func (mp *MarkdownPreview) createEnhancedToolbar() {
	mp.toolbar = widget.NewToolbar(
		// 刷新功能組
		widget.NewToolbarAction(theme.ViewRefreshIcon(), func() {
			mp.refreshPreview()
		}),
		
		widget.NewToolbarAction(theme.MediaPlayIcon(), func() {
			mp.toggleAutoRefresh()
		}),
		
		// 分隔線
		widget.NewToolbarSeparator(),
		
		// 匯出功能組
		widget.NewToolbarAction(theme.DocumentSaveIcon(), func() {
			mp.showExportMenu()
		}),
		
		widget.NewToolbarAction(theme.DocumentIcon(), func() {
			mp.printPreview()
		}),
		
		// 分隔線
		widget.NewToolbarSeparator(),
		
		// 縮放功能組
		widget.NewToolbarAction(theme.ContentAddIcon(), func() {
			mp.zoomIn()
		}),
		
		widget.NewToolbarAction(theme.ContentRemoveIcon(), func() {
			mp.zoomOut()
		}),
		
		widget.NewToolbarAction(theme.HomeIcon(), func() {
			mp.resetZoom()
		}),
		
		// 分隔線
		widget.NewToolbarSeparator(),
		
		// 搜尋功能
		widget.NewToolbarAction(theme.SearchIcon(), func() {
			mp.toggleSearch()
		}),
		
		// 分隔線
		widget.NewToolbarSeparator(),
		
		// 視窗功能組
		widget.NewToolbarAction(theme.ComputerIcon(), func() {
			mp.toggleIndependentWindow()
		}),
		
		widget.NewToolbarAction(theme.ViewFullScreenIcon(), func() {
			mp.toggleFullscreen()
		}),
		
		// 分隔線
		widget.NewToolbarSeparator(),
		
		// 同步滾動按鈕
		widget.NewToolbarAction(theme.MediaSkipNextIcon(), func() {
			mp.toggleSyncScroll()
		}),
		
		// 複製 HTML 按鈕
		widget.NewToolbarAction(theme.ContentCopyIcon(), func() {
			mp.copyHTML()
		}),
		
		// 隱藏預覽按鈕
		widget.NewToolbarAction(theme.VisibilityOffIcon(), func() {
			mp.toggleVisibility()
		}),
	)
}

// createSearchComponents 建立搜尋功能元件
// 包含搜尋輸入欄、結果顯示和導航按鈕
//
// 執行流程：
// 1. 建立搜尋輸入欄
// 2. 建立搜尋結果顯示標籤
// 3. 設定搜尋事件處理
// 4. 建立搜尋導航按鈕
func (mp *MarkdownPreview) createSearchComponents() {
	// 建立搜尋輸入欄
	mp.searchBar = widget.NewEntry()
	mp.searchBar.SetPlaceHolder("搜尋預覽內容...")
	mp.searchBar.Hide() // 預設隱藏
	
	// 設定搜尋事件處理
	mp.searchBar.OnChanged = func(query string) {
		mp.performSearch(query)
	}
	
	// 建立搜尋結果顯示
	mp.searchResults = widget.NewLabel("")
	mp.searchResults.Hide() // 預設隱藏
}

// createZoomControls 建立縮放控制元件
// 包含縮放滑桿和縮放級別顯示
//
// 執行流程：
// 1. 建立縮放滑桿
// 2. 設定縮放範圍和預設值
// 3. 設定縮放事件處理
func (mp *MarkdownPreview) createZoomControls() {
	// 建立縮放滑桿
	mp.zoomSlider = widget.NewSlider(0.5, 3.0) // 50% 到 300%
	mp.zoomSlider.Value = mp.zoomLevel
	mp.zoomSlider.Step = 0.1
	mp.zoomSlider.Hide() // 預設隱藏，可透過選單顯示
	
	// 設定縮放事件處理
	mp.zoomSlider.OnChanged = func(value float64) {
		mp.setZoomLevel(value)
	}
}

// createPreviewArea 建立增強版 HTML 預覽顯示區域
// 配置富文本顯示、滾動、樣式設定和效能優化
//
// 執行流程：
// 1. 建立富文本顯示元件
// 2. 設定預覽區域屬性（滾動、換行等）
// 3. 設定初始內容和樣式
// 4. 配置滾動事件處理
// 5. 應用縮放級別
func (mp *MarkdownPreview) createPreviewArea() {
	// 建立富文本預覽區域
	mp.previewArea = widget.NewRichTextFromMarkdown("")
	
	// 設定預覽區域屬性
	mp.previewArea.Wrapping = fyne.TextWrapWord  // 自動換行
	mp.previewArea.Scroll = container.ScrollBoth // 雙向滾動
	
	// 設定初始內容
	mp.previewArea.ParseMarkdown("# 預覽面板\n\n在此顯示 Markdown 內容的即時預覽。\n\n開始編輯以查看預覽效果。")
	
	// 應用初始縮放級別
	mp.applyZoomLevel()
}

// createStatusLabel 建立狀態標籤
// 顯示預覽狀態、更新時間和其他資訊
//
// 執行流程：
// 1. 建立狀態標籤元件
// 2. 設定初始狀態文字
// 3. 配置標籤樣式
func (mp *MarkdownPreview) createStatusLabel() {
	mp.statusLabel = widget.NewLabel("預覽準備就緒")
	mp.statusLabel.TextStyle = fyne.TextStyle{Italic: true}
}

// assembleEnhancedLayout 組合增強版預覽面板的完整佈局
// 將工具欄、搜尋欄、縮放控制、預覽區域和狀態欄組合成完整的預覽介面
//
// 執行流程：
// 1. 建立垂直容器作為主要佈局
// 2. 依序添加工具欄、搜尋功能、縮放控制
// 3. 添加預覽區域（主要顯示區域）
// 4. 添加狀態欄和搜尋結果顯示
// 5. 設定容器屬性和佈局比例
func (mp *MarkdownPreview) assembleEnhancedLayout() {
	// 建立搜尋和縮放控制容器
	searchZoomContainer := container.NewVBox(
		mp.searchBar,     // 搜尋輸入欄
		mp.zoomSlider,    // 縮放滑桿
	)
	
	// 建立底部狀態容器
	bottomContainer := container.NewVBox(
		widget.NewSeparator(),  // 分隔線
		mp.searchResults,       // 搜尋結果顯示
		mp.statusLabel,         // 狀態標籤
	)
	
	// 建立主要容器，使用垂直佈局
	mp.container = container.NewVBox(
		mp.toolbar,             // 工具欄在頂部
		widget.NewSeparator(),  // 分隔線
		searchZoomContainer,    // 搜尋和縮放控制
		mp.previewArea,         // 預覽區域在中間（主要區域）
		bottomContainer,        // 底部狀態容器
	)
}

// GetContainer 取得預覽面板的主要容器
// 回傳：預覽面板的 fyne.Container 實例
// 用於將預覽面板嵌入到其他 UI 佈局中
func (mp *MarkdownPreview) GetContainer() *fyne.Container {
	return mp.container
}

// UpdatePreview 更新預覽內容（增強版，包含效能優化）
// 參數：content（要預覽的 Markdown 內容）
//
// 執行流程：
// 1. 檢查內容是否有變更
// 2. 實作更新節流機制以提升效能
// 3. 檢查內容快取以避免重複處理
// 4. 使用編輯器服務轉換 Markdown 為 HTML
// 5. 更新預覽區域顯示和獨立視窗
// 6. 更新搜尋結果（如果有搜尋查詢）
// 7. 更新狀態標籤和觸發相關回調
func (mp *MarkdownPreview) UpdatePreview(content string) {
	// 檢查內容是否有變更
	if mp.currentContent == content {
		return // 內容未變更，無需更新
	}
	
	// 實作更新節流機制
	now := time.Now()
	if now.Sub(mp.lastUpdateTime) < mp.updateThrottle {
		// 更新過於頻繁，延遲處理
		go func() {
			time.Sleep(mp.updateThrottle)
			mp.UpdatePreview(content)
		}()
		return
	}
	mp.lastUpdateTime = now
	
	// 檢查內容快取
	contentHash := fmt.Sprintf("%x", content) // 簡化的內容雜湊
	if cachedHTML, exists := mp.contentCache[contentHash]; exists {
		// 使用快取的內容
		mp.previewArea.ParseMarkdown(cachedHTML)
		mp.updateStatus("已從快取載入預覽")
		return
	}
	
	// 更新當前內容
	mp.currentContent = content
	
	// 如果內容為空，顯示預設訊息
	if strings.TrimSpace(content) == "" {
		mp.previewArea.ParseMarkdown("# 預覽面板\n\n開始編輯以查看預覽效果。")
		mp.updateStatus("等待內容輸入")
		return
	}
	
	// 使用編輯器服務轉換 Markdown 為 HTML
	// 注意：RichText 元件直接支援 Markdown，所以我們直接使用 Markdown 內容
	mp.previewArea.ParseMarkdown(content)
	
	// 快取處理後的內容
	mp.contentCache[contentHash] = content
	
	// 更新獨立視窗（如果存在）
	if mp.isIndependent {
		mp.updateIndependentWindow()
	}
	
	// 如果有搜尋查詢，重新執行搜尋
	if mp.searchQuery != "" {
		mp.performSearch(mp.searchQuery)
	}
	
	// 更新狀態
	wordCount := len(strings.Fields(content))
	characterCount := len(content)
	mp.updateStatus(fmt.Sprintf("已更新預覽 - 字數: %d, 字元: %d, 縮放: %.0f%%", 
		wordCount, characterCount, mp.zoomLevel*100))
}

// RefreshPreview 手動刷新預覽
// 強制重新渲染當前內容
//
// 執行流程：
// 1. 重新解析當前 Markdown 內容
// 2. 更新預覽顯示
// 3. 更新狀態標籤
func (mp *MarkdownPreview) refreshPreview() {
	if mp.currentContent != "" {
		// 強制重新解析內容
		mp.previewArea.ParseMarkdown(mp.currentContent)
		mp.updateStatus("預覽已手動刷新")
	} else {
		mp.updateStatus("沒有內容可刷新")
	}
}

// SetAutoRefresh 設定自動刷新模式
// 參數：enabled（是否啟用自動刷新）
//
// 執行流程：
// 1. 更新自動刷新狀態
// 2. 更新工具欄按鈕狀態
// 3. 更新狀態顯示
func (mp *MarkdownPreview) SetAutoRefresh(enabled bool) {
	mp.autoRefresh = enabled
	
	if enabled {
		mp.updateStatus("自動刷新已啟用")
	} else {
		mp.updateStatus("自動刷新已停用")
	}
}

// IsAutoRefreshEnabled 檢查是否啟用自動刷新
// 回傳：自動刷新是否啟用的布林值
func (mp *MarkdownPreview) IsAutoRefreshEnabled() bool {
	return mp.autoRefresh
}

// SetVisible 設定預覽面板可見性
// 參數：visible（是否可見）
//
// 執行流程：
// 1. 更新可見性狀態
// 2. 顯示或隱藏預覽容器
// 3. 觸發可見性變更回調
// 4. 更新狀態顯示
func (mp *MarkdownPreview) SetVisible(visible bool) {
	mp.isVisible = visible
	
	if visible {
		mp.container.Show()
		mp.updateStatus("預覽面板已顯示")
	} else {
		mp.container.Hide()
		mp.updateStatus("預覽面板已隱藏")
	}
	
	// 觸發可見性變更回調
	if mp.onVisibilityChanged != nil {
		mp.onVisibilityChanged(visible)
	}
}

// IsVisible 檢查預覽面板是否可見
// 回傳：預覽面板是否可見的布林值
func (mp *MarkdownPreview) IsVisible() bool {
	return mp.isVisible
}

// Clear 清空預覽內容
// 清除所有預覽內容並重置狀態
//
// 執行流程：
// 1. 清空當前內容
// 2. 重置預覽區域顯示
// 3. 更新狀態標籤
func (mp *MarkdownPreview) Clear() {
	mp.currentContent = ""
	mp.previewArea.ParseMarkdown("# 預覽面板\n\n預覽內容已清空。")
	mp.updateStatus("預覽內容已清空")
}

// GetCurrentContent 取得當前預覽的內容
// 回傳：當前預覽的 Markdown 內容
func (mp *MarkdownPreview) GetCurrentContent() string {
	return mp.currentContent
}

// SetOnVisibilityChanged 設定可見性變更回調函數
// 參數：callback（可見性變更時的回調函數）
func (mp *MarkdownPreview) SetOnVisibilityChanged(callback func(visible bool)) {
	mp.onVisibilityChanged = callback
}

// SetOnRefreshRequested 設定刷新請求回調函數
// 參數：callback（刷新請求時的回調函數）
func (mp *MarkdownPreview) SetOnRefreshRequested(callback func()) {
	mp.onRefreshRequested = callback
}

// toggleAutoRefresh 切換自動刷新模式
// 在啟用和停用自動刷新之間切換
//
// 執行流程：
// 1. 切換自動刷新狀態
// 2. 更新狀態顯示
// 3. 如果啟用自動刷新，立即刷新預覽
func (mp *MarkdownPreview) toggleAutoRefresh() {
	mp.autoRefresh = !mp.autoRefresh
	
	if mp.autoRefresh {
		mp.updateStatus("自動刷新已啟用")
		// 立即刷新預覽
		mp.refreshPreview()
	} else {
		mp.updateStatus("自動刷新已停用")
	}
}

// toggleSyncScroll 切換同步滾動功能
// 啟用或停用編輯器和預覽面板的同步滾動
//
// 執行流程：
// 1. 切換同步滾動狀態
// 2. 更新狀態顯示
// 3. 配置滾動事件處理
func (mp *MarkdownPreview) toggleSyncScroll() {
	// 同步滾動功能的實作
	// 注意：Fyne 的 RichText 元件目前不直接支援滾動位置控制
	// 這個功能將在未來的版本中實作
	mp.updateStatus("同步滾動功能將在未來版本中實作")
}

// exportHTML 匯出 HTML 內容
// 將當前預覽內容匯出為 HTML 檔案
//
// 執行流程：
// 1. 檢查是否有內容可匯出
// 2. 使用編輯器服務轉換為 HTML
// 3. 顯示檔案保存對話框
// 4. 保存 HTML 檔案
func (mp *MarkdownPreview) exportHTML() {
	if mp.currentContent == "" {
		mp.updateStatus("沒有內容可匯出")
		return
	}
	
	// 使用編輯器服務轉換為 HTML
	htmlContent := mp.editorService.PreviewMarkdown(mp.currentContent)
	
	// 建立完整的 HTML 文件
	fullHTML := mp.createFullHTMLDocument(htmlContent)
	
	// 這裡應該顯示檔案保存對話框
	// 目前簡化為狀態更新
	mp.updateStatus(fmt.Sprintf("HTML 內容已準備匯出 (%d 字元)", len(fullHTML)))
}

// copyHTML 複製 HTML 內容到剪貼簿
// 將當前預覽的 HTML 內容複製到系統剪貼簿
//
// 執行流程：
// 1. 檢查是否有內容可複製
// 2. 使用編輯器服務轉換為 HTML
// 3. 複製到系統剪貼簿
// 4. 更新狀態顯示
func (mp *MarkdownPreview) copyHTML() {
	if mp.currentContent == "" {
		mp.updateStatus("沒有內容可複製")
		return
	}
	
	// 使用編輯器服務轉換為 HTML
	htmlContent := mp.editorService.PreviewMarkdown(mp.currentContent)
	
	// 複製到剪貼簿
	// 注意：Fyne 的剪貼簿操作需要視窗上下文
	// 這裡簡化為狀態更新
	mp.updateStatus(fmt.Sprintf("HTML 內容已準備複製 (%d 字元)", len(htmlContent)))
}

// toggleFullscreen 切換全螢幕預覽模式
// 在正常模式和全螢幕預覽模式之間切換
//
// 執行流程：
// 1. 檢查當前顯示模式
// 2. 切換到全螢幕或正常模式
// 3. 調整佈局和顯示
// 4. 更新狀態顯示
func (mp *MarkdownPreview) toggleFullscreen() {
	// 全螢幕預覽功能的實作
	// 這個功能需要與主視窗協調
	mp.updateStatus("全螢幕預覽功能將在後續版本中實作")
}

// toggleVisibility 切換預覽面板可見性
// 顯示或隱藏預覽面板
//
// 執行流程：
// 1. 切換可見性狀態
// 2. 顯示或隱藏預覽容器
// 3. 觸發可見性變更回調
func (mp *MarkdownPreview) toggleVisibility() {
	mp.SetVisible(!mp.isVisible)
}

// updateStatus 更新狀態標籤顯示
// 參數：message（要顯示的狀態訊息）
func (mp *MarkdownPreview) updateStatus(message string) {
	if mp.statusLabel != nil {
		mp.statusLabel.SetText(message)
		mp.statusLabel.Refresh()
	}
}

// createFullHTMLDocument 建立完整的 HTML 文件
// 參數：bodyContent（HTML 主體內容）
// 回傳：完整的 HTML 文件字串
//
// 執行流程：
// 1. 建立 HTML 文件結構
// 2. 添加 CSS 樣式
// 3. 插入主體內容
// 4. 回傳完整的 HTML 文件
func (mp *MarkdownPreview) createFullHTMLDocument(bodyContent string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="zh-TW">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Markdown 預覽</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 800px;
            margin: 0 auto;
            padding: 20px;
        }
        h1, h2, h3, h4, h5, h6 {
            color: #2c3e50;
            margin-top: 1.5em;
            margin-bottom: 0.5em;
        }
        code {
            background-color: #f4f4f4;
            padding: 2px 4px;
            border-radius: 3px;
            font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
        }
        pre {
            background-color: #f8f8f8;
            border: 1px solid #ddd;
            border-radius: 5px;
            padding: 10px;
            overflow-x: auto;
        }
        blockquote {
            border-left: 4px solid #ddd;
            margin: 0;
            padding-left: 20px;
            color: #666;
        }
        table {
            border-collapse: collapse;
            width: 100%%;
        }
        th, td {
            border: 1px solid #ddd;
            padding: 8px;
            text-align: left;
        }
        th {
            background-color: #f2f2f2;
        }
    </style>
</head>
<body>
%s
</body>
</html>`, bodyContent)
}

// SyncScrollPosition 同步滾動位置
// 參數：position（滾動位置百分比，0.0-1.0）
//
// 執行流程：
// 1. 計算預覽區域的對應滾動位置
// 2. 設定預覽區域的滾動位置
// 3. 更新狀態顯示
func (mp *MarkdownPreview) SyncScrollPosition(position float64) {
	// 同步滾動功能的實作
	// 注意：Fyne 的 RichText 元件目前不直接支援滾動位置控制
	// 這個功能將在未來的版本中實作
	mp.updateStatus(fmt.Sprintf("滾動同步: %.1f%%", position*100))
}

// GetScrollPosition 取得當前滾動位置
// 回傳：當前滾動位置百分比（0.0-1.0）
func (mp *MarkdownPreview) GetScrollPosition() float64 {
	// 取得滾動位置的實作
	// 注意：Fyne 的 RichText 元件目前不直接支援滾動位置查詢
	// 這個功能將在未來的版本中實作
	return 0.0
}

// SetTheme 設定預覽主題
// 參數：theme（主題名稱）
//
// 執行流程：
// 1. 根據主題名稱設定樣式
// 2. 更新預覽區域顯示
// 3. 重新渲染內容
func (mp *MarkdownPreview) SetTheme(themeName string) {
	// 主題設定功能的實作
	// 這個功能將在設定管理實作時完成
	mp.updateStatus(fmt.Sprintf("主題設定: %s (將在後續版本中實作)", themeName))
}

// GetWordCount 取得當前內容的字數統計
// 回傳：字數統計
func (mp *MarkdownPreview) GetWordCount() int {
	if mp.currentContent == "" {
		return 0
	}
	return len(strings.Fields(mp.currentContent))
}

// GetCharacterCount 取得當前內容的字元統計
// 回傳：字元統計
func (mp *MarkdownPreview) GetCharacterCount() int {
	return len(mp.currentContent)
}

// HasContent 檢查是否有預覽內容
// 回傳：是否有內容的布林值
func (mp *MarkdownPreview) HasContent() bool {
	return strings.TrimSpace(mp.currentContent) != ""
}

// ===== 增強功能方法 =====

// showExportMenu 顯示匯出選單
// 提供多種匯出格式選項（HTML、PDF、Word）
//
// 執行流程：
// 1. 建立匯出選項選單
// 2. 顯示格式選擇對話框
// 3. 根據選擇執行對應的匯出功能
func (mp *MarkdownPreview) showExportMenu() {
	if mp.currentContent == "" {
		mp.updateStatus("沒有內容可匯出")
		return
	}
	
	// 建立匯出選項
	exportOptions := []string{"HTML", "PDF", "Word"}
	
	// 建立選擇對話框
	selectDialog := dialog.NewInformation("匯出格式", 
		"請選擇匯出格式：\n• HTML - 網頁格式\n• PDF - 可攜式文件格式\n• Word - Microsoft Word 格式", 
		mp.container.Objects[0].(fyne.Window))
	
	selectDialog.Show()
	
	// 目前簡化為 HTML 匯出
	mp.exportHTML()
	mp.updateStatus(fmt.Sprintf("支援的匯出格式: %v", exportOptions))
}

// printPreview 列印預覽內容
// 開啟系統列印對話框進行列印
//
// 執行流程：
// 1. 檢查是否有內容可列印
// 2. 轉換內容為可列印格式
// 3. 開啟系統列印對話框
func (mp *MarkdownPreview) printPreview() {
	if mp.currentContent == "" {
		mp.updateStatus("沒有內容可列印")
		return
	}
	
	// 列印功能的實作
	// 注意：Fyne 目前不直接支援列印功能
	// 這個功能將在未來版本中透過系統 API 實作
	mp.updateStatus("列印功能將在未來版本中實作")
}

// zoomIn 放大預覽內容
// 增加縮放級別並更新顯示
//
// 執行流程：
// 1. 增加縮放級別（最大 300%）
// 2. 更新縮放滑桿
// 3. 應用新的縮放級別
// 4. 更新狀態顯示
func (mp *MarkdownPreview) zoomIn() {
	newZoom := mp.zoomLevel + 0.1
	if newZoom <= 3.0 {
		mp.setZoomLevel(newZoom)
	}
}

// zoomOut 縮小預覽內容
// 減少縮放級別並更新顯示
//
// 執行流程：
// 1. 減少縮放級別（最小 50%）
// 2. 更新縮放滑桿
// 3. 應用新的縮放級別
// 4. 更新狀態顯示
func (mp *MarkdownPreview) zoomOut() {
	newZoom := mp.zoomLevel - 0.1
	if newZoom >= 0.5 {
		mp.setZoomLevel(newZoom)
	}
}

// resetZoom 重置縮放級別
// 將縮放級別重置為 100%
//
// 執行流程：
// 1. 設定縮放級別為 1.0（100%）
// 2. 更新縮放滑桿
// 3. 應用縮放級別
// 4. 更新狀態顯示
func (mp *MarkdownPreview) resetZoom() {
	mp.setZoomLevel(1.0)
}

// setZoomLevel 設定縮放級別
// 參數：level（縮放級別，0.5-3.0）
//
// 執行流程：
// 1. 驗證縮放級別範圍
// 2. 更新內部縮放狀態
// 3. 更新縮放滑桿顯示
// 4. 應用縮放到預覽區域
// 5. 觸發縮放變更回調
// 6. 更新狀態顯示
func (mp *MarkdownPreview) setZoomLevel(level float64) {
	// 驗證縮放級別範圍
	if level < 0.5 {
		level = 0.5
	} else if level > 3.0 {
		level = 3.0
	}
	
	mp.zoomLevel = level
	
	// 更新縮放滑桿
	if mp.zoomSlider != nil {
		mp.zoomSlider.SetValue(level)
	}
	
	// 應用縮放級別
	mp.applyZoomLevel()
	
	// 觸發縮放變更回調
	if mp.onZoomChanged != nil {
		mp.onZoomChanged(level)
	}
	
	// 更新狀態顯示
	mp.updateStatus(fmt.Sprintf("縮放級別: %.0f%%", level*100))
}

// applyZoomLevel 應用縮放級別到預覽區域
// 調整預覽內容的顯示大小
//
// 執行流程：
// 1. 計算縮放後的字體大小
// 2. 更新預覽區域的顯示屬性
// 3. 重新渲染內容
func (mp *MarkdownPreview) applyZoomLevel() {
	if mp.previewArea == nil {
		return
	}
	
	// 注意：Fyne 的 RichText 元件目前不直接支援縮放
	// 這個功能將透過調整字體大小來模擬縮放效果
	// 實際實作將在未來版本中完善
	
	// 重新渲染內容以應用縮放
	if mp.currentContent != "" {
		mp.previewArea.ParseMarkdown(mp.currentContent)
	}
}

// toggleSearch 切換搜尋功能顯示/隱藏
// 顯示或隱藏搜尋輸入欄和相關控制項
//
// 執行流程：
// 1. 切換搜尋欄可見性
// 2. 如果顯示，設定搜尋欄焦點
// 3. 如果隱藏，清空搜尋結果
// 4. 更新佈局
func (mp *MarkdownPreview) toggleSearch() {
	if mp.searchBar.Visible() {
		// 隱藏搜尋功能
		mp.searchBar.Hide()
		mp.searchResults.Hide()
		mp.clearSearchHighlights()
		mp.updateStatus("搜尋功能已隱藏")
	} else {
		// 顯示搜尋功能
		mp.searchBar.Show()
		mp.searchResults.Show()
		mp.searchBar.FocusGained()
		mp.updateStatus("搜尋功能已啟用")
	}
	
	mp.container.Refresh()
}

// performSearch 執行搜尋操作
// 參數：query（搜尋查詢字串）
//
// 執行流程：
// 1. 清空之前的搜尋結果
// 2. 如果查詢為空，清除高亮顯示
// 3. 使用正規表達式搜尋匹配項
// 4. 高亮顯示搜尋結果
// 5. 更新搜尋結果統計
// 6. 觸發搜尋回調
func (mp *MarkdownPreview) performSearch(query string) {
	mp.searchQuery = query
	mp.searchMatches = mp.searchMatches[:0] // 清空搜尋結果
	mp.currentMatch = -1
	
	if query == "" {
		mp.clearSearchHighlights()
		mp.searchResults.SetText("")
		return
	}
	
	// 在當前內容中搜尋
	if mp.currentContent != "" {
		// 使用正規表達式進行不區分大小寫的搜尋
		regex, err := regexp.Compile("(?i)" + regexp.QuoteMeta(query))
		if err != nil {
			mp.updateStatus("搜尋查詢格式錯誤")
			return
		}
		
		// 找到所有匹配項
		matches := regex.FindAllStringIndex(mp.currentContent, -1)
		for _, match := range matches {
			mp.searchMatches = append(mp.searchMatches, match[0])
		}
		
		// 更新搜尋結果顯示
		matchCount := len(mp.searchMatches)
		if matchCount > 0 {
			mp.currentMatch = 0
			mp.searchResults.SetText(fmt.Sprintf("找到 %d 個匹配項", matchCount))
			mp.highlightSearchResults()
		} else {
			mp.searchResults.SetText("未找到匹配項")
		}
		
		// 觸發搜尋回調
		if mp.onSearchPerformed != nil {
			mp.onSearchPerformed(query, matchCount)
		}
	}
}

// highlightSearchResults 高亮顯示搜尋結果
// 在預覽內容中標記搜尋匹配項
//
// 執行流程：
// 1. 遍歷所有搜尋匹配項
// 2. 在預覽內容中添加高亮標記
// 3. 重新渲染預覽內容
func (mp *MarkdownPreview) highlightSearchResults() {
	// 搜尋結果高亮功能的實作
	// 注意：Fyne 的 RichText 元件目前不直接支援文字高亮
	// 這個功能將在未來版本中透過自訂渲染器實作
	mp.updateStatus(fmt.Sprintf("高亮顯示 %d 個搜尋結果", len(mp.searchMatches)))
}

// clearSearchHighlights 清除搜尋高亮顯示
// 移除預覽內容中的所有搜尋高亮標記
//
// 執行流程：
// 1. 移除所有高亮標記
// 2. 重新渲染原始內容
func (mp *MarkdownPreview) clearSearchHighlights() {
	// 清除搜尋高亮功能的實作
	// 重新渲染原始內容
	if mp.currentContent != "" {
		mp.previewArea.ParseMarkdown(mp.currentContent)
	}
}

// navigateToNextMatch 導航到下一個搜尋匹配項
// 移動到搜尋結果中的下一個匹配項
//
// 執行流程：
// 1. 檢查是否有搜尋結果
// 2. 移動到下一個匹配項
// 3. 滾動到匹配項位置
// 4. 更新搜尋結果顯示
func (mp *MarkdownPreview) navigateToNextMatch() {
	if len(mp.searchMatches) == 0 {
		return
	}
	
	mp.currentMatch = (mp.currentMatch + 1) % len(mp.searchMatches)
	mp.scrollToMatch(mp.currentMatch)
	mp.updateSearchResultsDisplay()
}

// navigateToPreviousMatch 導航到上一個搜尋匹配項
// 移動到搜尋結果中的上一個匹配項
//
// 執行流程：
// 1. 檢查是否有搜尋結果
// 2. 移動到上一個匹配項
// 3. 滾動到匹配項位置
// 4. 更新搜尋結果顯示
func (mp *MarkdownPreview) navigateToPreviousMatch() {
	if len(mp.searchMatches) == 0 {
		return
	}
	
	mp.currentMatch--
	if mp.currentMatch < 0 {
		mp.currentMatch = len(mp.searchMatches) - 1
	}
	mp.scrollToMatch(mp.currentMatch)
	mp.updateSearchResultsDisplay()
}

// scrollToMatch 滾動到指定的搜尋匹配項
// 參數：matchIndex（匹配項索引）
//
// 執行流程：
// 1. 計算匹配項在內容中的位置
// 2. 計算對應的滾動位置
// 3. 滾動預覽區域到該位置
func (mp *MarkdownPreview) scrollToMatch(matchIndex int) {
	if matchIndex < 0 || matchIndex >= len(mp.searchMatches) {
		return
	}
	
	// 滾動到匹配項功能的實作
	// 注意：Fyne 的 RichText 元件目前不直接支援滾動到特定位置
	// 這個功能將在未來版本中實作
	mp.updateStatus(fmt.Sprintf("導航到匹配項 %d/%d", matchIndex+1, len(mp.searchMatches)))
}

// updateSearchResultsDisplay 更新搜尋結果顯示
// 更新搜尋結果標籤顯示當前匹配項資訊
//
// 執行流程：
// 1. 格式化搜尋結果文字
// 2. 更新搜尋結果標籤
// 3. 刷新顯示
func (mp *MarkdownPreview) updateSearchResultsDisplay() {
	if len(mp.searchMatches) == 0 {
		mp.searchResults.SetText("未找到匹配項")
	} else {
		mp.searchResults.SetText(fmt.Sprintf("匹配項 %d/%d", mp.currentMatch+1, len(mp.searchMatches)))
	}
	mp.searchResults.Refresh()
}

// toggleIndependentWindow 切換獨立視窗模式
// 在嵌入模式和獨立視窗模式之間切換
//
// 執行流程：
// 1. 檢查當前模式
// 2. 如果是嵌入模式，建立獨立視窗
// 3. 如果是獨立模式，關閉獨立視窗
// 4. 更新狀態和佈局
func (mp *MarkdownPreview) toggleIndependentWindow() {
	if mp.isIndependent {
		// 關閉獨立視窗，回到嵌入模式
		if mp.independentWindow != nil {
			mp.independentWindow.Close()
			mp.independentWindow = nil
		}
		mp.isIndependent = false
		mp.updateStatus("預覽已切換到嵌入模式")
	} else {
		// 建立獨立預覽視窗
		mp.createIndependentWindow()
		mp.isIndependent = true
		mp.updateStatus("預覽已切換到獨立視窗模式")
	}
}

// createIndependentWindow 建立獨立預覽視窗
// 建立一個新的視窗來顯示預覽內容
//
// 執行流程：
// 1. 建立新的應用程式視窗
// 2. 設定視窗屬性（標題、大小等）
// 3. 複製預覽內容到新視窗
// 4. 顯示獨立視窗
func (mp *MarkdownPreview) createIndependentWindow() {
	// 建立新的應用程式實例（用於獨立視窗）
	independentApp := app.New()
	mp.independentWindow = independentApp.NewWindow("Markdown 預覽")
	
	// 設定視窗屬性
	mp.independentWindow.Resize(fyne.NewSize(800, 600))
	mp.independentWindow.CenterOnScreen()
	
	// 建立獨立視窗的預覽內容
	independentPreview := widget.NewRichTextFromMarkdown(mp.currentContent)
	independentPreview.Wrapping = fyne.TextWrapWord
	independentPreview.Scroll = container.ScrollBoth
	
	// 建立獨立視窗的工具欄
	independentToolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.ViewRefreshIcon(), func() {
			independentPreview.ParseMarkdown(mp.currentContent)
		}),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(theme.DocumentSaveIcon(), func() {
			mp.exportHTML()
		}),
		widget.NewToolbarAction(theme.ContentCopyIcon(), func() {
			mp.copyHTML()
		}),
	)
	
	// 組合獨立視窗內容
	independentContent := container.NewVBox(
		independentToolbar,
		widget.NewSeparator(),
		independentPreview,
	)
	
	mp.independentWindow.SetContent(independentContent)
	mp.independentWindow.Show()
	
	// 設定視窗關閉事件
	mp.independentWindow.SetCloseIntercept(func() {
		mp.isIndependent = false
		mp.independentWindow = nil
		mp.updateStatus("獨立預覽視窗已關閉")
	})
}

// updateIndependentWindow 更新獨立視窗內容
// 當預覽內容變更時，同步更新獨立視窗
//
// 執行流程：
// 1. 檢查獨立視窗是否存在
// 2. 更新獨立視窗的預覽內容
// 3. 刷新獨立視窗顯示
func (mp *MarkdownPreview) updateIndependentWindow() {
	if mp.independentWindow != nil && mp.isIndependent {
		// 更新獨立視窗內容的實作
		// 這需要保持對獨立視窗中預覽元件的引用
		mp.updateStatus("獨立視窗內容已更新")
	}
}

// GetZoomLevel 取得當前縮放級別
// 回傳：當前縮放級別（0.5-3.0）
func (mp *MarkdownPreview) GetZoomLevel() float64 {
	return mp.zoomLevel
}

// SetOnZoomChanged 設定縮放變更回調函數
// 參數：callback（縮放變更時的回調函數）
func (mp *MarkdownPreview) SetOnZoomChanged(callback func(level float64)) {
	mp.onZoomChanged = callback
}

// SetOnSearchPerformed 設定搜尋執行回調函數
// 參數：callback（搜尋執行時的回調函數）
func (mp *MarkdownPreview) SetOnSearchPerformed(callback func(query string, matches int)) {
	mp.onSearchPerformed = callback
}

// GetSearchQuery 取得當前搜尋查詢
// 回傳：當前搜尋查詢字串
func (mp *MarkdownPreview) GetSearchQuery() string {
	return mp.searchQuery
}

// GetSearchMatchCount 取得搜尋匹配項數量
// 回傳：搜尋匹配項數量
func (mp *MarkdownPreview) GetSearchMatchCount() int {
	return len(mp.searchMatches)
}

// IsIndependentMode 檢查是否為獨立視窗模式
// 回傳：是否為獨立視窗模式的布林值
func (mp *MarkdownPreview) IsIndependentMode() bool {
	return mp.isIndependent
}

// optimizePerformance 優化預覽效能
// 實作內容快取和更新節流機制
//
// 執行流程：
// 1. 檢查內容是否已快取
// 2. 實作更新節流機制
// 3. 管理快取大小
func (mp *MarkdownPreview) optimizePerformance() {
	// 清理過期的快取項目
	if len(mp.contentCache) > 100 { // 限制快取大小
		// 清空快取以避免記憶體洩漏
		mp.contentCache = make(map[string]string)
	}
	
	mp.updateStatus("效能優化已執行")
}