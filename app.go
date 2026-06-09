package main

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"Booklet/pkg/booklet"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// ShowSuccessDialog shows a dialog with "Open Folder" and "Close" options
func (a *App) ShowSuccessDialog(message string) string {
	resp, err := wailsRuntime.MessageDialog(a.ctx, wailsRuntime.MessageDialogOptions{
		Type:          wailsRuntime.QuestionDialog,
		Title:         "생성 완료",
		Message:       message,
		Buttons:       []string{"폴더 열기", "닫기"},
		DefaultButton: "폴더 열기",
	})
	if err != nil {
		return "닫기"
	}
	return resp
}

// SelectFile opens a system file dialog to select a PDF file
func (a *App) SelectFile() (string, error) {
	selection, err := wailsRuntime.OpenFileDialog(a.ctx, wailsRuntime.OpenDialogOptions{
		Title: "PDF 파일 선택",
		Filters: []wailsRuntime.FileFilter{
			{DisplayName: "PDF Files (*.pdf)", Pattern: "*.pdf"},
		},
	})
	return selection, err
}

// SelectSaveFile opens a dialog to select the output file path
func (a *App) SelectSaveFile(inputPath string) (string, error) {
	var defaultDir, defaultFilename string
	if inputPath != "" {
		defaultDir = filepath.Dir(inputPath)
		ext := filepath.Ext(inputPath)
		base := strings.TrimSuffix(filepath.Base(inputPath), ext)
		defaultFilename = base + "_booklet.pdf"
	} else {
		defaultFilename = "booklet_output.pdf"
	}

	selection, err := wailsRuntime.SaveFileDialog(a.ctx, wailsRuntime.SaveDialogOptions{
		Title:            "결과 파일 저장",
		DefaultDirectory: defaultDir,
		DefaultFilename:  defaultFilename,
		Filters: []wailsRuntime.FileFilter{
			{DisplayName: "PDF Files (*.pdf)", Pattern: "*.pdf"},
		},
	})
	return selection, err
}

// ProcessBooklet processes the input PDF and creates a booklet
func (a *App) ProcessBooklet(opts booklet.Options) string {
	err := booklet.Process(opts)
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}
	return "Success"
}

// OpenFolder opens the specified path in the default file browser
func (a *App) OpenFolder(path string) {
	dir := filepath.Dir(path)
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("explorer", dir)
	case "darwin":
		cmd = exec.Command("open", dir)
	case "linux":
		cmd = exec.Command("xdg-open", dir)
	default:
		return
	}

	_ = cmd.Run()
}

// GetPageCount returns the page count of the PDF file at the specified path
func (a *App) GetPageCount(path string) (int, error) {
	ctx, err := api.ReadContextFile(path)
	if err != nil {
		return 0, fmt.Errorf("PDF 파일 분석 실패: %v", err)
	}
	return ctx.PageCount, nil
}
