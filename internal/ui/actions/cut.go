package actions

import (
	"fyne.io/fyne/v2/widget"
)

type CutAction struct {
	BaseAction
	ctx ActionContext
}

func NewCutAction(ctx ActionContext) Action {
	a := &CutAction{ctx: ctx}
	a.name = "cut"
	a.button = widget.NewButton("✂️ Cut", func() {
		if a.ctx.EventBus != nil {
			a.ctx.EventBus.Publish("editor.cut", nil)
		}
	})
	a.SetEnabled(false)
	return a
}

func (a *CutAction) UpdateState(state EditorState) {
	a.SetEnabled(state.HasSelection)
}

func init() {
	RegisterAction("cut", NewCutAction)
}
