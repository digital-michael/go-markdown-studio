package ui

import (
	"log"
	"os"

	"fyne.io/fyne/v2/widget"
)

// Editor encapsulates the markdown editor widget and its state.
type Editor struct {
	Widget          *widget.Entry
	SaveButton      *widget.Button
	CurrentFilePath string
	OriginalContent string
	IsDirty         bool
}

// NewEditor creates and returns a new Editor struct.
func NewEditor() *Editor {
	editor := &Editor{}

	editor.Widget = widget.NewMultiLineEntry()
	editor.Widget.SetPlaceHolder("Select a markdown file to edit...")

	editor.SaveButton = widget.NewButton("💾 Save", func() {
		editor.Save()
	})
	editor.SaveButton.Disable() // start disabled

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
			e.SaveButton.Disable()
		}
	}
}

// OnContentChanged handles changes in the editor content.
func (e *Editor) OnContentChanged(content string) {
	if content != e.OriginalContent {
		e.IsDirty = true
		e.SaveButton.Enable()
	} else {
		e.IsDirty = false
		e.SaveButton.Disable()
	}
}

// SetFile loads a file into the editor and updates state.
func (e *Editor) SetFile(path string, content string) {
	e.CurrentFilePath = path
	e.OriginalContent = content
	e.Widget.SetText(content)
	e.IsDirty = false
	e.SaveButton.Disable()
}
