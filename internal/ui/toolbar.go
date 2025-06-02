package ui

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"

	"go-markdown-studio/internal/config"
	"go-markdown-studio/internal/ui/actions"
)

// Toolbar holds and manages a set of actions as a UI component.
type Toolbar struct {
	Name    string
	Actions []actions.Action
	Box     *fyne.Container
}

// NewToolbar creates a toolbar based on config, action registry, and context.
func NewToolbar(name string, cfg *config.AppConfig, ctx actions.ActionContext) *Toolbar {
	var toolbarCfg *config.ToolbarConfig
	for _, tb := range cfg.Toolbars {
		if tb.Name == name {
			toolbarCfg = &tb
			break
		}
	}
	if toolbarCfg == nil {
		log.Printf("Toolbar config '%s' not found", name)
		return nil
	}

	var objs []fyne.CanvasObject
	var acts []actions.Action

	for _, actName := range toolbarCfg.Actions {
		if actName == "separator" {
			objs = append(objs, container.NewVBox()) // Simple separator, can be improved
			continue
		}
		constructor, ok := actions.ActionRegistry[actName]
		if !ok {
			log.Printf("Unknown action: %s", actName)
			continue
		}
		act := constructor(ctx)
		acts = append(acts, act)
		objs = append(objs, act.CanvasObject())
	}

	var box *fyne.Container
	if toolbarCfg.Orientation == "vertical" {
		box = container.NewVBox(objs...)
	} else {
		box = container.NewHBox(objs...)
	}

	return &Toolbar{
		Name:    name,
		Actions: acts,
		Box:     box,
	}
}

// Widget returns the fyne.CanvasObject for the toolbar.
func (tb *Toolbar) Widget() fyne.CanvasObject {
	return tb.Box
}

// UpdateState propagates editor state to all actions.
func (tb *Toolbar) UpdateState(state actions.EditorState) {
	for _, act := range tb.Actions {
		act.UpdateState(state)
	}
}
