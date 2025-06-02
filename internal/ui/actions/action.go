package actions

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// EditorState is a lightweight struct for editor state propagation.
type EditorState struct {
	CurrentFilePath string
	IsDirty         bool
	HasSelection    bool
	// Add more as needed
}

// Action is the interface all toolbar actions must implement.
type Action interface {
	Name() string
	CanvasObject() fyne.CanvasObject
	SetEnabled(enabled bool)
	UpdateState(state EditorState)
	OnEvent(event string, payload any)
}

// BaseAction provides common fields and methods for actions.
type BaseAction struct {
	name    string
	button  fyne.CanvasObject
	enabled bool
}

func (a *BaseAction) Name() string {
	return a.name
}

func (a *BaseAction) SetEnabled(enabled bool) {
	a.enabled = enabled
	if btn, ok := a.button.(*widget.Button); ok {
		btn.Disable()
		if enabled {
			btn.Enable()
		}
	}
}

func (a *BaseAction) CanvasObject() fyne.CanvasObject {
	return a.button
}

func (a *BaseAction) UpdateState(state EditorState) {
	// Default: do nothing
}

func (a *BaseAction) OnEvent(event string, payload any) {
	// Default: do nothing
}

// ActionRegistry holds constructors for all available actions.
var ActionRegistry = map[string]func(ctx ActionContext) Action{}

// ActionContext provides context for action construction.
type ActionContext struct {
	EditorState EditorState
	EventBus    EventBus
	// Add more as needed (e.g., app, window)
}

// RegisterAction registers an action constructor by name.
func RegisterAction(name string, constructor func(ctx ActionContext) Action) {
	ActionRegistry[name] = constructor
}
