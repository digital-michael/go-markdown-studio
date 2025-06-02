package actions

import (
	"fyne.io/fyne/v2/widget"
)

type DeleteFileAction struct {
	BaseAction
	ctx ActionContext
}

func NewDeleteFileAction(ctx ActionContext) Action {
	a := &DeleteFileAction{ctx: ctx}
	a.name = "deletefile"
	a.button = widget.NewButton("🗑️ Delete", func() {
		if a.ctx.EventBus != nil {
			a.ctx.EventBus.Publish("app.deletefile", nil)
		}
	})
	a.SetEnabled(false)
	return a
}

func (a *DeleteFileAction) UpdateState(state EditorState) {
	a.SetEnabled(state.CurrentFilePath != "")
}

func init() {
	RegisterAction("deletefile", NewDeleteFileAction)
}
