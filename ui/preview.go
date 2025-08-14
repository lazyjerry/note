// Package ui 包含使用者介面相關的元件和視窗管理
// 本檔案實作 Markdown 預覽面板 UI 元件，提供即時預覽和同步滾動功能
package ui

import (
	"fmt"                           // Go 標準庫，用於格式化字串
	"strings"                       // 字串處理
	"mac-notebook-app/internal/services" // 引入服務層

	"fyne.io/fyne/v2"               // Fyne GUI 框架核心套件
	"fyne.io/fyne/v2/container"     // Fyne 容器佈局套件
	"fyne.io/fyne/v2/widget"        // Fyne UI 元件套件
	"fyne.io/fyne/v2/theme"         // Fyne 主題套件
)

// MarkdownPreview 代表 Markdown 預覽面板 UI 元件
// 提供即時 HTML 預覽、同步滾動和預覽控制功能
// 整合編輯器服務以提供 Markdown 到 HTML 的轉換
type MarkdownPreview struct {
	container     *fyne.Container      // 主要容器
	toolbar       *widget.Toolbar      // 預覽工具欄
	previewArea   *widget.RichText     // HTML 預覽顯示區域
	statusLabel   *widget.Label        // 狀態標籤
	
	// 服務依賴
	editorService services.EditorService // 編輯器服務
	
	// 當前狀態
	currentContent string              // 當前預覽的 Markdown 內容
	isVisible     bool                 // 預覽面板是否可見
	autoRefresh   bool                 // 是否自動刷新預覽
	
	// 回調函數
	onVisibilityChanged func(visible bool) // 可見性變更回調
	onRefreshRequested  func()             // 刷新請求回調
}

// NewMarkdownPreview 建立新的 Markdown 預覽面板實例
// 參數：editorService（編輯器服務介面）
// 回傳：指向新建立的 MarkdownPreview 的指標
//
// 執行流程：
// 1. 建立 MarkdownPreview 結構體實例
// 2. 初始化編輯器服務依賴
// 3. 建立並配置所有 UI 元件
// 4. 設定預設狀態和屬性
// 5. 組合完整的預覽面板佈局
// 6. 回傳配置完成的預覽實例
func NewMarkdownPreview(editorService services.EditorService) *MarkdownPreview {
	// 建立 MarkdownPreview 實例
	preview := &MarkdownPreview{
		editorService: editorService,
		isVisible:     true,  // 預設可見
		autoRefresh:   true,  // 預設自動刷新
	}
	
	// 初始化 UI 元件
	preview.setupUI()
	
	return preview
}

// setupUI 初始化預覽面板的使用者介面元件
// 這個方法負責建立和配置預覽面板的所有 UI 元件
//
// 執行流程：
// 1. 建立預覽工具欄和控制按鈕
// 2. 建立 HTML 預覽顯示區域
// 3. 建立狀態標籤顯示預覽狀態
// 4. 設定事件處理和回調函數
// 5. 組合所有元件到主容器中
func (mp *MarkdownPreview) setupUI() {
	// 建立預覽工具欄
	mp.createToolbar()
	
	// 建立預覽顯示區域
	mp.createPreviewArea()
	
	// 建立狀態標籤
	mp.createStatusLabel()
	
	// 組合預覽面板佈局
	mp.assembleLayout()
}

// createToolbar 建立預覽面板工具欄
// 包含刷新、同步滾動、匯出等預覽控制功能
//
// 執行流程：
// 1. 建立手動刷新按鈕
// 2. 建立自動刷新切換按鈕
// 3. 建立同步滾動切換按鈕
// 4. 建立匯出 HTML 按鈕
// 5. 建立預覽隱藏/顯示按鈕
// 6. 組合所有按鈕到工具欄
func (mp *MarkdownPreview) createToolbar() {
	mp.toolbar = widget.NewToolbar(
		// 手動刷新按鈕
		widget.NewToolbarAction(theme.ViewRefreshIcon(), func() {
			mp.refreshPreview()
		}),
		
		// 分隔線
		widget.NewToolbarSeparator(),
		
		// 自動刷新切換按鈕
		widget.NewToolbarAction(theme.MediaPlayIcon(), func() {
			mp.toggleAutoRefresh()
		}),
		
		// 同步滾動按鈕
		widget.NewToolbarAction(theme.MediaSkipNextIcon(), func() {
			mp.toggleSyncScroll()
		}),
		
		// 分隔線
		widget.NewToolbarSeparator(),
		
		// 匯出 HTML 按鈕
		widget.NewToolbarAction(theme.DocumentSaveIcon(), func() {
			mp.exportHTML()
		}),
		
		// 複製 HTML 按鈕
		widget.NewToolbarAction(theme.ContentCopyIcon(), func() {
			mp.copyHTML()
		}),
		
		// 分隔線
		widget.NewToolbarSeparator(),
		
		// 全螢幕預覽按鈕
		widget.NewToolbarAction(theme.ViewFullScreenIcon(), func() {
			mp.toggleFullscreen()
		}),
		
		// 隱藏預覽按鈕
		widget.NewToolbarAction(theme.VisibilityOffIcon(), func() {
			mp.toggleVisibility()
		}),
	)
}

// createPreviewArea 建立 HTML 預覽顯示區域
// 配置富文本顯示、滾動和樣式設定
//
// 執行流程：
// 1. 建立富文本顯示元件
// 2. 設定預覽區域屬性（滾動、換行等）
// 3. 設定初始內容和樣式
// 4. 配置滾動事件處理
func (mp *MarkdownPreview) createPreviewArea() {
	// 建立富文本預覽區域
	mp.previewArea = widget.NewRichTextFromMarkdown("")
	
	// 設定預覽區域屬性
	mp.previewArea.Wrapping = fyne.TextWrapWord  // 自動換行
	mp.previewArea.Scroll = container.ScrollBoth // 雙向滾動
	
	// 設定初始內容
	mp.previewArea.ParseMarkdown("# 預覽面板\n\n在此顯示 Markdown 內容的即時預覽。\n\n開始編輯以查看預覽效果。")
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

// assembleLayout 組合預覽面板的完整佈局
// 將工具欄、預覽區域和狀態欄組合成完整的預覽介面
//
// 執行流程：
// 1. 建立垂直容器作為主要佈局
// 2. 依序添加工具欄、預覽區域和狀態欄
// 3. 設定容器屬性和佈局比例
func (mp *MarkdownPreview) assembleLayout() {
	// 建立主要容器，使用垂直佈局
	mp.container = container.NewVBox(
		mp.toolbar,                    // 工具欄在頂部
		widget.NewSeparator(),         // 分隔線
		mp.previewArea,                // 預覽區域在中間（主要區域）
		widget.NewSeparator(),         // 分隔線
		mp.statusLabel,                // 狀態欄在底部
	)
}

// GetContainer 取得預覽面板的主要容器
// 回傳：預覽面板的 fyne.Container 實例
// 用於將預覽面板嵌入到其他 UI 佈局中
func (mp *MarkdownPreview) GetContainer() *fyne.Container {
	return mp.container
}

// UpdatePreview 更新預覽內容
// 參數：content（要預覽的 Markdown 內容）
//
// 執行流程：
// 1. 檢查內容是否有變更
// 2. 使用編輯器服務轉換 Markdown 為 HTML
// 3. 更新預覽區域顯示
// 4. 更新狀態標籤
// 5. 觸發相關回調
func (mp *MarkdownPreview) UpdatePreview(content string) {
	// 檢查內容是否有變更
	if mp.currentContent == content {
		return // 內容未變更，無需更新
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
	
	// 更新狀態
	wordCount := len(strings.Fields(content))
	mp.updateStatus(fmt.Sprintf("已更新預覽 - 字數: %d", wordCount))
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