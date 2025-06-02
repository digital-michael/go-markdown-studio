package config

import (
	"encoding/json"
	"log"
	"os"
)

type DirectoryEntry struct {
	Path      string `json:"path"`
	Recursive bool   `json:"recursive"`
}

type ToolbarConfig struct {
	Name        string   `json:"name"`
	Orientation string   `json:"orientation"` // "horizontal" or "vertical"
	Actions     []string `json:"actions"`
}

type AppConfig struct {
	Theme       string           `json:"theme"`
	Directories []DirectoryEntry `json:"directories"`
	Toolbars    []ToolbarConfig  `json:"toolbars"`
}

const configFile = "config.json"

func defaultConfig() AppConfig {
	return AppConfig{
		Theme: "system",
		Directories: []DirectoryEntry{
			{
				Path:      ".",
				Recursive: false,
			},
		},
		Toolbars: []ToolbarConfig{
			{
				Name:        "editorMain",
				Orientation: "horizontal",
				Actions: []string{
					"newfile",
					"separator",
					"save",
					"copy",
					"cut",
					"paste",
					"separator",
					"undo",
					"redo",
					"separator",
					"deletefile",
					"movefile",
				},
			},
		},
	}
}

func LoadConfig() AppConfig {
	cfg := defaultConfig()
	file, err := os.Open(configFile)
	if err != nil {
		log.Println("No config file found, creating default config.")
		SaveConfig(cfg)
		return cfg
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cfg)
	if err != nil {
		log.Println("Error decoding config file, using default:", err)
		cfg = defaultConfig()
		SaveConfig(cfg)
		return cfg
	}

	return cfg
}

func SaveConfig(cfg AppConfig) {
	file, err := os.Create(configFile)
	if err != nil {
		log.Println("Error saving config file:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(cfg)
	if err != nil {
		log.Println("Error encoding config file:", err)
	}
	log.Println("Config saved successfully.")
}
