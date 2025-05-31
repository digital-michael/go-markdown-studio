package ui

import (
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"

	"go-markdown-studio/internal/config"
)

// FileList encapsulates the file list widget and its state.
type FileList struct {
	List      *widget.List
	FileNames []string
	MDFiles   []string
	Cfg       *config.AppConfig
}

// NewFileList creates and returns a new FileList struct.
func NewFileList(cfg *config.AppConfig) *FileList {
	fl := &FileList{Cfg: cfg}
	fl.RefreshFiles()

	fl.List = widget.NewList(
		func() int {
			return len(fl.FileNames)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(fl.FileNames[i])
		},
	)

	return fl
}

// RefreshFiles updates the list of markdown files and their names.
func (fl *FileList) RefreshFiles() {
	fl.MDFiles = ScanMarkdownFilesFromConfig(fl.Cfg)
	fl.FileNames = make([]string, len(fl.MDFiles))
	for i, path := range fl.MDFiles {
		fl.FileNames[i] = filepath.Base(path)
	}
}

// UpdateList refreshes the file list widget and its data.
func (fl *FileList) UpdateList() {
	fl.RefreshFiles()
	fl.List.Refresh()
}

// OnSelected sets the callback for when a file is selected.
func (fl *FileList) OnSelected(callback func(int)) {
	fl.List.OnSelected = callback
}

// Widget returns the underlying fyne.CanvasObject for the file list.
func (fl *FileList) Widget() fyne.CanvasObject {
	return fl.List
}
