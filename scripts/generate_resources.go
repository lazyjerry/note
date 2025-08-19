package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// generateResources 生成應用程式資源檔案
// 此函數負責將應用程式的圖示和其他資源檔案打包成 Go 程式碼
//
// 執行流程：
// 1. 檢查 fyne bundle 工具是否可用
// 2. 生成圖示資源檔案
// 3. 生成字體資源檔案
// 4. 創建資源包檔案
func generateResources() error {
	log.Println("開始生成應用程式資源檔案...")

	// 檢查 fyne bundle 工具
	if err := checkFyneBundle(); err != nil {
		return err
	}

	// 生成圖示資源
	if err := generateIconResource(); err != nil {
		return err
	}

	// 生成字體資源
	if err := generateFontResources(); err != nil {
		return err
	}

	log.Println("資源檔案生成完成！")
	return nil
}

// checkFyneBundle 檢查 fyne bundle 工具是否可用
// 回傳：如果工具不可用則回傳錯誤
func checkFyneBundle() error {
	log.Println("檢查 fyne bundle 工具...")
	
	cmd := exec.Command("fyne", "bundle", "--help")
	if err := cmd.Run(); err != nil {
		log.Println("fyne bundle 工具不可用，嘗試安裝...")
		installCmd := exec.Command("go", "install", "fyne.io/fyne/v2/cmd/fyne@latest")
		if err := installCmd.Run(); err != nil {
			return err
		}
		log.Println("fyne 工具安裝完成")
	}
	
	return nil
}

// generateIconResource 生成應用程式圖示資源
// 回傳：生成過程中的錯誤（如果有）
//
// 執行流程：
// 1. 檢查圖示檔案是否存在
// 2. 使用 fyne bundle 生成 Go 資源檔案
// 3. 將生成的檔案保存到指定位置
func generateIconResource() error {
	log.Println("生成應用程式圖示資源...")
	
	iconPath := "assets/Icon.png"
	outputPath := "internal/resources/icon.go"
	
	// 檢查圖示檔案是否存在
	if _, err := os.Stat(iconPath); os.IsNotExist(err) {
		log.Printf("警告：圖示檔案 %s 不存在，跳過圖示資源生成", iconPath)
		return nil
	}
	
	// 確保輸出目錄存在
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return err
	}
	
	// 生成圖示資源
	cmd := exec.Command("fyne", "bundle", "--name", "AppIcon", "-o", outputPath, iconPath)
	if err := cmd.Run(); err != nil {
		return err
	}
	
	log.Printf("圖示資源已生成：%s", outputPath)
	return nil
}

// generateFontResources 生成字體資源檔案
// 回傳：生成過程中的錯誤（如果有）
//
// 執行流程：
// 1. 掃描字體目錄
// 2. 為每個字體檔案生成資源
// 3. 創建字體資源包
func generateFontResources() error {
	log.Println("生成字體資源...")
	
	fontDir := "assets/font"
	outputDir := "internal/resources"
	
	// 檢查字體目錄是否存在
	if _, err := os.Stat(fontDir); os.IsNotExist(err) {
		log.Printf("警告：字體目錄 %s 不存在，跳過字體資源生成", fontDir)
		return nil
	}
	
	// 確保輸出目錄存在
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}
	
	// 掃描字體檔案
	fontFiles, err := filepath.Glob(filepath.Join(fontDir, "*.ttf"))
	if err != nil {
		return err
	}
	
	// 為每個字體檔案生成資源
	for _, fontFile := range fontFiles {
		fontName := filepath.Base(fontFile)
		varName := "Font" + filepath.Base(fontName[:len(fontName)-4]) // 移除 .ttf 副檔名
		outputPath := filepath.Join(outputDir, "font_"+fontName[:len(fontName)-4]+".go")
		
		cmd := exec.Command("fyne", "bundle", "--name", varName, "-o", outputPath, fontFile)
		if err := cmd.Run(); err != nil {
			log.Printf("警告：無法生成字體資源 %s: %v", fontFile, err)
			continue
		}
		
		log.Printf("字體資源已生成：%s", outputPath)
	}
	
	return nil
}

func main() {
	if err := generateResources(); err != nil {
		log.Fatalf("資源生成失敗：%v", err)
	}
}