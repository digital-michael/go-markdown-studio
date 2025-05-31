package ui

import (
	"log"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"go-markdown-studio/internal/config"
)

// BuildMainUI constructs the main UI layout using Editor and FileList components.
func BuildMainUI(a fyne.App, w fyne.Window, cfg *config.AppConfig) fyne.CanvasObject {
	// Ensure config directories are not empty
	ensureConfigDirectories(cfg)

	// Markdown editor
	editor := NewEditor()

	// File list
	fileList := NewFileList(cfg)

	selectedPathLabel := widget.NewLabel("")
	selectedPathLabel.Wrapping = fyne.TextWrapBreak

	// add selectors
	// Handle file selection: load file into editor and update label
	fileList.OnSelected(func(i int) {
		if i >= 0 && i < len(fileList.MDFiles) {
			path := fileList.MDFiles[i]
			content, err := os.ReadFile(path)
			if err == nil {
				editor.SetFile(path, string(content))
				selectedPathLabel.SetText(path)
			} else {
				selectedPathLabel.SetText("Failed to load: " + path)
			}
		}
	})

	// handle the event handlers
	addEventHandlers(a, w, *editor, *fileList, cfg)

	// Layout
	leftPanel := container.NewStack(fileList.Widget())
	rightPanel := container.NewBorder(nil, editor.SaveButton, nil, nil, editor.Widget)
	split := container.NewHSplit(leftPanel, rightPanel)
	split.Offset = 0.25

	main := container.NewBorder(nil, selectedPathLabel, nil, nil, split)
	return main
}

func addEventHandlers(a fyne.App, w fyne.Window, editor Editor, fileList FileList, cfg *config.AppConfig) {
	// Add any additional event handlers here if needed
	// For example, you could handle window close events or other actions
	// w.SetCloseIntercept(func() {
	// 	log.Println("Window close intercepted. You can add cleanup code here.")
	// 	// Optionally save state or perform cleanup
	// })

	updateChan := make(chan string)

	cwd, err := os.Getwd()
	if err != nil {
		log.Println("Could not get working directory:", err)
		cwd = "."
	}

	go func() {
		for cwd := range updateChan {
			fileList.RefreshFiles()
			fyne.DoAndWait(func() {
				log.Printf("Updating file list for directory: %s ... ", cwd)
				fileList.UpdateList()
				log.Print("Done")
			})
		}
	}()

	WatchMarkdownDir(cwd, func() {
		log.Printf("Directory change detected, updating file list...%s\n", cwd)
		fyne.Do(func() { updateChan <- cwd })
	})
}

// ensureConfigDirectories guarantees at least one directory is present in config.
func ensureConfigDirectories(cfg *config.AppConfig) {
	if len(cfg.Directories) == 0 {
		cwd, err := os.Getwd()
		if err != nil {
			log.Println("Could not get working directory:", err)
			cwd = "." // fallback
		}
		cfg.Directories = []config.DirectoryEntry{
			{
				Path:      cwd,
				Recursive: false,
			},
		}
		log.Printf("No directories in config; using current directory: %s", cwd)
		log.Printf("Loaded config: %+v", *cfg)
		config.SaveConfig(*cfg)
	}
}

// getFileNames returns the base names of the given file paths.
func getFileNames(mdFiles []string) []string {
	fileNames := make([]string, len(mdFiles))
	for i, path := range mdFiles {
		fileNames[i] = filepath.Base(path)
	}
	return fileNames
}
