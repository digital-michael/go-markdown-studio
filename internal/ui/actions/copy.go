package actions

import (
	"fyne.io/fyne/v2/widget"
)

type CopyAction struct {
	BaseAction
	ctx ActionContext
}

func NewCopyAction(ctx ActionContext) Action {
	a := &CopyAction{ctx: ctx}
	a.name = "copy"
	a.button = widget.NewButton("📋 Copy", func() {
		if a.ctx.EventBus != nil {
			a.ctx.EventBus.Publish("editor.copy", nil)
		}
	})
	a.SetEnabled(false)
	return a
}

func (a *CopyAction) UpdateState(state EditorState) {
	a.SetEnabled(state.HasSelection)
}

func init() {
	RegisterAction("copy", NewCopyAction)
}
