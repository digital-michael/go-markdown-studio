package ui

import (
	"log"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"go-markdown-studio/internal/config"
	"go-markdown-studio/internal/ui/actions"
)

// BuildMainUI constructs the main UI layout using Editor and FileList components.
func BuildMainUI(a fyne.App, w fyne.Window, cfg *config.AppConfig) fyne.CanvasObject {
	// Ensure config directories are not empty
	ensureConfigDirectories(cfg)

	// Create the event bus
	eventBus := actions.NewSimpleEventBus()

	// Markdown editor
	editor := NewEditor(cfg, eventBus)

	// File list
	fileList := NewFileList(cfg)

	selectedPathLabel := widget.NewLabel("")
	selectedPathLabel.Wrapping = fyne.TextWrapBreak

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

	// ...existing eventBus.Subscribe handlers...

	eventBus.Subscribe("app.newfile", func(_ any) {
		// TODO: Implement new file logic (e.g., show dialog, create file, refresh list)
		log.Println("New file requested")
	})

	eventBus.Subscribe("app.deletefile", func(_ any) {
		// TODO: Implement delete file logic (e.g., confirm, delete, refresh list)
		log.Println("Delete file requested")
	})

	eventBus.Subscribe("app.movefile", func(_ any) {
		// TODO: Implement move file logic (e.g., show dialog, move, refresh list)
		log.Println("Move file requested")
	})

	eventBus.Subscribe("editor.undo", func(_ any) {
		// TODO: Implement undo logic (requires undo stack in editor)
		log.Println("Undo requested")
	})

	eventBus.Subscribe("editor.redo", func(_ any) {
		// TODO: Implement redo logic (requires redo stack in editor)
		log.Println("Redo requested")
	})

	// Wire up event bus handlers for editor actions
	eventBus.Subscribe("editor.save", func(_ any) {
		editor.Save()
	})
	eventBus.Subscribe("editor.copy", func(_ any) {
		editor.Widget.TypedShortcut(&fyne.ShortcutCopy{Clipboard: fyne.CurrentApp().Driver().AllWindows()[0].Clipboard()})
	})
	eventBus.Subscribe("editor.cut", func(_ any) {
		editor.Widget.TypedShortcut(&fyne.ShortcutCut{Clipboard: fyne.CurrentApp().Driver().AllWindows()[0].Clipboard()})
	})
	eventBus.Subscribe("editor.paste", func(_ any) {
		editor.Widget.TypedShortcut(&fyne.ShortcutPaste{Clipboard: fyne.CurrentApp().Driver().AllWindows()[0].Clipboard()})
	})
	// Add more event handlers for newfile, deletefile, movefile, undo, redo as needed

	// Directory/file watcher logic
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

	// Layout
	leftPanel := container.NewStack(fileList.Widget())
	rightPanel := container.NewBorder(editor.Toolbar.Widget(), nil, nil, nil, editor.Widget)
	split := container.NewHSplit(leftPanel, rightPanel)
	split.Offset = 0.25

	main := container.NewBorder(nil, selectedPathLabel, nil, nil, split)
	return main
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
