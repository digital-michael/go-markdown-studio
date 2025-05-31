package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
	// "fyne.io/fyne/v2/widget"

	// "fyne.io/fyne/v2/container"

	"fyne.io/fyne/v2"
	// "log"

	"go-markdown-studio/internal/config"
	"go-markdown-studio/internal/ui"
)

func main() {
	a := app.NewWithID("com.example.markdownstudio")
	if a == nil {
		panic("Failed to create fyne app instance")
	}
	cfg := config.LoadConfig()
	// TODO: don't save here as this cfg and the one in ui.BuildMainUI are different
	//       Let the UI handle saving the config when needed
	// defer config.SaveConfig(cfg) // Save config on exit

	switch cfg.Theme {
	case "dark":
		a.Settings().SetTheme(theme.DarkTheme())
	case "light":
		a.Settings().SetTheme(theme.LightTheme())
	default:
		// Use system default or fallback
	}

	w := a.NewWindow("Markdown Studio")
	w.Resize(fyne.NewSize(1000, 600))

	mainUI := ui.BuildMainUI(a, w, &cfg)
	w.SetContent(mainUI)

	w.ShowAndRun()

}
