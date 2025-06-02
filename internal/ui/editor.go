package ui

import (
	"go-markdown-studio/internal/config"
	"go-markdown-studio/internal/ui/actions"
	"log"
	"os"

	"fyne.io/fyne/v2/widget"
)

// Editor encapsulates the markdown editor widget and its state.
type Editor struct {
	Widget          *widget.Entry
	Toolbar         *Toolbar
	CurrentFilePath string
	OriginalContent string
	IsDirty         bool
	EventBus        actions.EventBus
}

// NewEditor creates and returns a new Editor struct.
func NewEditor(cfg *config.AppConfig, eventBus actions.EventBus) *Editor {
	editor := &Editor{
		EventBus: eventBus,
	}

	editor.Widget = widget.NewMultiLineEntry()
	editor.Widget.SetPlaceHolder("Select a markdown file to edit...")

	// Example: use toolbar named "editorMain"
	ctx := actions.ActionContext{
		EventBus: eventBus,
	}
	editor.Toolbar = NewToolbar("editorMain", cfg, ctx)

	editor.Widget.OnChanged = func(content string) {
		editor.OnContentChanged(content)
	}

	return editor
}

// Save writes the editor content to the current file path.
func (e *Editor) Save() {
	if e.CurrentFilePath != "" {
		err := os.WriteFile(e.CurrentFilePath, []byte(e.Widget.Text), 0644)
		if err != nil {
			log.Printf("Failed to save: %v", err)
		} else {
			log.Printf("Saved: %s", e.CurrentFilePath)
			e.OriginalContent = e.Widget.Text
			e.IsDirty = false
			e.UpdateToolbarState()
		}
	}
}

// OnContentChanged handles changes in the editor content.
func (e *Editor) OnContentChanged(content string) {
	if content != e.OriginalContent {
		e.IsDirty = true
	} else {
		e.IsDirty = false
	}
	e.UpdateToolbarState()
}

// SetFile loads a file into the editor and updates state.
func (e *Editor) SetFile(path string, content string) {
	e.CurrentFilePath = path
	e.OriginalContent = content
	e.Widget.SetText(content)
	e.IsDirty = false
	e.UpdateToolbarState()
}

// Call this when editor state changes
func (e *Editor) UpdateToolbarState() {
	state := actions.EditorState{
		CurrentFilePath: e.CurrentFilePath,
		IsDirty:         e.IsDirty,
		HasSelection:    false, // Fyne Entry does not support selection info directly
	}
	if e.Toolbar != nil {
		e.Toolbar.UpdateState(state)
	}
}
