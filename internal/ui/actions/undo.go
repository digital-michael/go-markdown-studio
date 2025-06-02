package actions

import (
	"fyne.io/fyne/v2/widget"
)

type UndoAction struct {
	BaseAction
	ctx ActionContext
}

func NewUndoAction(ctx ActionContext) Action {
	a := &UndoAction{ctx: ctx}
	a.name = "undo"
	a.button = widget.NewButton("↩️ Undo", func() {
		if a.ctx.EventBus != nil {
			a.ctx.EventBus.Publish("editor.undo", nil)
		}
	})
	a.SetEnabled(false)
	return a
}

func (a *UndoAction) UpdateState(state EditorState) {
	// Enable if undo is available (expand EditorState as needed)
	a.SetEnabled(true)
}

func init() {
	RegisterAction("undo", NewUndoAction)
}
