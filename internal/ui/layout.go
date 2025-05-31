package ui

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/fsnotify/fsnotify"
	"go-markdown-studio/internal/config"
)

func BuildMainUI(a fyne.App, w fyne.Window, cfg *config.AppConfig) fyne.CanvasObject {
	// Markdown editor
	editor := widget.NewMultiLineEntry()
	editor.SetPlaceHolder("Select a markdown file to edit...")

	// // Register Ctrl+S shortcut on the editor
	// ctrlS := &desktop.CustomShortcut{
	// 	KeyName:  fyne.KeyS,
	// 	Modifier: fyne.KeyModifierControl,
	// }

	var currentFilePath string
	var originalContent string
	var isDirty bool

	// Get markdown files
	cwd, err := os.Getwd()
	if err != nil {
		log.Println("Could not get working directory:", err)
		cwd = "."
	}

	// Ensure config directories are not empty
	if len(cfg.Directories) == 0 {
		cwd, err := os.Getwd()
		if err != nil {
			log.Println("Could not get working directory:", err)
			cwd = "." // fallback
		}

		// Add current working directory as non-recursive default
		cfg.Directories = []config.DirectoryEntry{
			{
				Path:      cwd,
				Recursive: false,
			},
		}

		log.Printf("No directories in config; using current directory: %s", cwd)
		log.Printf("Loaded config: %+v", *cfg)
		config.SaveConfig(*cfg) // Save config now, not later
	}

	mdFiles := scanMarkdownFilesFromConfig(cfg)

	fileNames := make([]string, len(mdFiles))
	for i, path := range mdFiles {
		fileNames[i] = filepath.Base(path)
	}

	selectedPathLabel := widget.NewLabel("")
	selectedPathLabel.Wrapping = fyne.TextWrapBreak

	// UI list widget
	list := widget.NewList(
		func() int {
			return len(fileNames)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(fileNames[i])
		},
	)

	updateListsAll := func() {
		mdFiles = scanMarkdownFilesFromConfig(cfg)
		fileNames = make([]string, len(mdFiles))
		for i, path := range mdFiles {
			fileNames[i] = filepath.Base(path)
		}
		list.Refresh()
	}

	// updateList := func(cwd string) {
	// 	fyne.CurrentApp().Driver().RunOnMainThread(func() {
	// 		mdFiles = scanMarkdownFilesFlat(cwd)
	// 		fileNames = make([]string, len(mdFiles))
	// 		for i, path := range mdFiles {
	// 			fileNames[i] = filepath.Base(path)
	// 		}
	// 		list.Refresh()
	// 	})
	// }

	updateChan := make(chan string) // send updated cwd values

	go func() {
		for cwd := range updateChan {
			md := scanMarkdownFilesFlat(cwd)
			names := make([]string, len(md))
			for i, path := range md {
				names[i] = filepath.Base(path)
			}

			fyne.App().Driver().RunOnMainThread(func() {
				log.Printf("Updating file list for directory: %s ... ", cwd)
				list.Refresh()
				log.Print( "Done")
			})
		}
	}()

	// Initial call
	updateListsAll()

	// Start watching the directory
	WatchMarkdownDir(cwd, func() {
		log.Printf("Directory change detected, updating file list...%s\n", cwd)
		fyne.Do(func() { updateChan <- cwd } )
	})

	// Define saveBtn as a variable so it can be referenced in closures
	var saveBtn *widget.Button
	saveBtn = widget.NewButton("💾 Save", func() {

		if currentFilePath != "" {
			err := os.WriteFile(currentFilePath, []byte(editor.Text), 0644)
			if err != nil {
				log.Printf("Failed to save: %v", err)
			} else {
				log.Printf("Saved: %s", currentFilePath)
				originalContent = editor.Text
				isDirty = false
				saveBtn.Disable()
			}
		}
	})
	// fixedSize := fyne.NewSize(100, 40)
	buttonContainer := container.NewHBox(layout.NewSpacer(), saveBtn)

	saveBtn.Disable() // start disabled

	loadFile := func(i int) {
		content, err := os.ReadFile(mdFiles[i])
		if err != nil {
			log.Printf("Could not read file %s: %v", mdFiles[i], err)
			editor.SetText("Error loading file.")
			return
		}

		currentFilePath = mdFiles[i]
		originalContent = string(content)
		editor.SetText(originalContent)
		selectedPathLabel.SetText(currentFilePath)
		w.SetTitle("Markdown Studio - " + fileNames[i])
		isDirty = false
		saveBtn.Disable()
		log.Printf("Loaded file: %s", currentFilePath)
	}

	// Load file into editor on select
	list.OnSelected = func(i widget.ListItemID) {
		// prompt if dirty
		if isDirty {
			confirm := dialog.NewConfirm("Unsaved Changes",
				"You have unsaved changes. Do you want to save them before switching files?",
				func(save bool) {
					if save && currentFilePath != "" {
						err := os.WriteFile(currentFilePath, []byte(editor.Text), 0644)
						if err != nil {
							log.Printf("Failed to save: %v", err)
							return
						}
						log.Printf("Saved: %s", currentFilePath)
						isDirty = false
						saveBtn.Disable()
					}

					loadFile(i)
				}, w)
			confirm.Show()
		} else {
			loadFile(i)
		}
	}

	w.SetCloseIntercept(func() {
		if isDirty {
			confirm := dialog.NewConfirm("Unsaved Changes",
				"You have unsaved changes. Do you want to save before exiting?",
				func(save bool) {
					if save && currentFilePath != "" {
						err := os.WriteFile(currentFilePath, []byte(editor.Text), 0644)
						if err != nil {
							log.Printf("Failed to save: %v", err)
							return
						}
					}
					w.Close() // proceed with closing the window
				}, w)
			confirm.Show()
		} else {
			w.Close() // no unsaved changes, just close
		}
	})

	w.Canvas().AddShortcut(
		&desktop.CustomShortcut{KeyName: fyne.KeyS, Modifier: fyne.KeyModifierControl},
		func(shortcut fyne.Shortcut) {
			log.Printf("Ctrl+S pressed, isDirty: %v, currentFilePath: %s", isDirty, currentFilePath)
			if isDirty && currentFilePath != "" {
				err := os.WriteFile(currentFilePath, []byte(editor.Text), 0644)
				if err != nil {
					log.Printf("Failed to save: %v", err)
				} else {
					originalContent = editor.Text
					isDirty = false
					saveBtn.Disable()
					log.Printf("Saved via Ctrl+S: %s", currentFilePath)
				}
			}
		},
	)

	editor.OnChanged = func(content string) {
		if content != originalContent {
			isDirty = true
			saveBtn.Enable()
		} else {
			isDirty = false
			saveBtn.Disable()
		}
	}

	// Create the main layout
	leftPanel := container.NewMax(list)
	rightPanel := container.NewBorder(buttonContainer, nil, nil, nil, editor)
	split := container.NewHSplit(leftPanel, rightPanel)
	split.Offset = 0.25

	main := container.NewBorder(nil, selectedPathLabel, nil, nil, split)
	return main
}

// scanMarkdownFiles walks a directory and returns all .md/.MD files
func scanMarkdownFilesFromConfig(cfg *config.AppConfig) []string {
	var allFiles []string
	for _, entry := range cfg.Directories {
		if entry.Recursive {
			allFiles = append(allFiles, scanMarkdownFilesRecursive(entry.Path)...)
		} else {
			allFiles = append(allFiles, scanMarkdownFilesFlat(entry.Path)...)
		}
	}
	return allFiles
}

// Walk all subdirectories
func scanMarkdownFilesRecursive(dir string) []string {
	var files []string
	filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if !d.IsDir() && strings.HasSuffix(strings.ToLower(d.Name()), ".md") {
			files = append(files, path)
		}
		return nil
	})
	return files
}

// Scan only one level
func scanMarkdownFilesFlat(dir string) []string {
	var files []string
	entries, err := os.ReadDir(dir)
	if err != nil {
		return files
	}
	for _, d := range entries {
		if !d.IsDir() && strings.HasSuffix(strings.ToLower(d.Name()), ".md") {
			files = append(files, filepath.Join(dir, d.Name()))
		}
	}
	return files
}

// WatchMarkdownDir watches the specified directory and updates the file list on changes.
func WatchMarkdownDir(
	dir string,
	updateFiles func(), // function to call when files change
) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Println("Failed to initialize watcher:", err)
		return
	}

	go func() {
		defer watcher.Close()

		debounce := time.NewTimer(0)
		<-debounce.C // discard initial tick

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if isRelevantEvent(event) {
					debounce.Reset(500 * time.Millisecond)
				}

			case <-debounce.C:
				updateFiles() // refresh list in UI
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("Watcher error:", err)
			}
		}
	}()

	err = watcher.Add(dir)
	if err != nil {
		log.Println("Failed to add watch directory:", err)
	}
}

func isRelevantEvent(event fsnotify.Event) bool {
	name := strings.ToLower(event.Name)
	return strings.HasSuffix(name, ".md") && (event.Op&fsnotify.Create == fsnotify.Create || event.Op&fsnotify.Remove == fsnotify.Remove)
}
