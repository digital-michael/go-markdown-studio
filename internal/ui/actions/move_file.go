package actions

import (
	"fyne.io/fyne/v2/widget"
)

type MoveFileAction struct {
	BaseAction
	ctx ActionContext
}

func NewMoveFileAction(ctx ActionContext) Action {
	a := &MoveFileAction{ctx: ctx}
	a.name = "movefile"
	a.button = widget.NewButton("📂 Move", func() {
		if a.ctx.EventBus != nil {
			a.ctx.EventBus.Publish("app.movefile", nil)
		}
	})
	a.SetEnabled(false)
	return a
}

func (a *MoveFileAction) UpdateState(state EditorState) {
	a.SetEnabled(state.CurrentFilePath != "")
}

func init() {
	RegisterAction("movefile", NewMoveFileAction)
}
