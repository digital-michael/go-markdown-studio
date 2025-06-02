package actions

import (
	"fyne.io/fyne/v2/widget"
)

type NewFileAction struct {
	BaseAction
	ctx ActionContext
}

func NewNewFileAction(ctx ActionContext) Action {
	a := &NewFileAction{ctx: ctx}
	a.name = "newfile"
	a.button = widget.NewButton("🆕 New", func() {
		if a.ctx.EventBus != nil {
			a.ctx.EventBus.Publish("app.newfile", nil)
		}
	})
	a.SetEnabled(true)
	return a
}

func (a *NewFileAction) UpdateState(state EditorState) {
	a.SetEnabled(true)
}

func init() {
	RegisterAction("newfile", NewNewFileAction)
}
