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

type AppConfig struct {
	Theme       string           `json:"theme"`
	Directories []DirectoryEntry `json:"directories"`
}

const configFile = "config.json"

func LoadConfig() AppConfig {
	cfg := AppConfig{
		Theme:       "system",
		Directories: []DirectoryEntry{},
	}
	file, err := os.Open(configFile)
	if err != nil {
		log.Println("No config file found, using default.")
		return cfg
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cfg)
	if err != nil {
		log.Println("Error decoding config file:", err)
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
