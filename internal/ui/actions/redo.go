package actions

import (
	"fyne.io/fyne/v2/widget"
)

type RedoAction struct {
	BaseAction
	ctx ActionContext
}

func NewRedoAction(ctx ActionContext) Action {
	a := &RedoAction{ctx: ctx}
	a.name = "redo"
	a.button = widget.NewButton("↪️ Redo", func() {
		if a.ctx.EventBus != nil {
			a.ctx.EventBus.Publish("editor.redo", nil)
		}
	})
	a.SetEnabled(false)
	return a
}

func (a *RedoAction) UpdateState(state EditorState) {
	// Enable if redo is available (expand EditorState as needed)
	a.SetEnabled(true)
}

func init() {
	RegisterAction("redo", NewRedoAction)
}
