package test

import (
	"go-markdown-studio/internal/config"
	"testing"
)

func TestConfigSaveLoad(t *testing.T) {
	cfg := config.AppConfig{
		Theme:       "dark",
		Directories: []string{"/tmp/docs"},
	}
	config.SaveConfig(cfg)

	loaded := config.LoadConfig()

	if loaded.Theme != "dark" {
		t.Errorf("Expected theme 'dark', got %s", loaded.Theme)
	}
	if len(loaded.Directories) != 1 || loaded.Directories[0] != "/tmp/docs" {
		t.Errorf("Expected directory '/tmp/docs', got %v", loaded.Directories)
	}
}
