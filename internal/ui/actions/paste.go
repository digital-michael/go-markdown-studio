package actions

import (
	"fyne.io/fyne/v2/widget"
)

type PasteAction struct {
	BaseAction
	ctx ActionContext
}

func NewPasteAction(ctx ActionContext) Action {
	a := &PasteAction{ctx: ctx}
	a.name = "paste"
	a.button = widget.NewButton("📋 Paste", func() {
		if a.ctx.EventBus != nil {
			a.ctx.EventBus.Publish("editor.paste", nil)
		}
	})
	a.SetEnabled(true) // Paste is usually enabled
	return a
}

func (a *PasteAction) UpdateState(state EditorState) {
	// Optionally, check clipboard state
	a.SetEnabled(true)
}

func init() {
	RegisterAction("paste", NewPasteAction)
}
