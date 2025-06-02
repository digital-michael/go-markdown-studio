package actions

import (
	"fyne.io/fyne/v2/widget"
)

// SaveAction implements the Action interface for saving the current file.
type SaveAction struct {
	BaseAction
	ctx     ActionContext
	onClick func()
}

func NewSaveAction(ctx ActionContext) Action {
	a := &SaveAction{
		ctx: ctx,
	}
	a.name = "save"
	a.button = widget.NewButton("💾 Save", func() {
		if a.onClick != nil {
			a.onClick()
		}
		// Publish a save event
		if a.ctx.EventBus != nil {
			a.ctx.EventBus.Publish("editor.save", nil)
		}
	})
	a.SetEnabled(false) // Disabled by default

	// Example: Enable when editor is dirty
	a.onClick = func() {
		// Actual save logic should be handled by listening to "editor.save" event elsewhere
	}

	return a
}

func (a *SaveAction) UpdateState(state EditorState) {
	a.SetEnabled(state.IsDirty)
}

// Register SaveAction in the registry
func init() {
	RegisterAction("save", NewSaveAction)
}
