package test

import (
	"os"
	"path/filepath"
	"testing"

	"go-markdown-studio/internal/config"
)

func TestConfigSaveLoad(t *testing.T) {
	tmpFile := filepath.Join(os.TempDir(), "test_config.json")
	defer os.Remove(tmpFile)

	cfg := config.AppConfig{
		Theme: "dark",
		Directories: []config.DirectoryEntry{
			{Path: "/tmp/docs", Recursive: false},
		},
	}
	config.SaveConfig(cfg)

	loaded := config.LoadConfig()

	if loaded.Theme != "dark" {
		t.Errorf("Expected theme 'dark', got %s", loaded.Theme)
	}
	if len(loaded.Directories) != 1 || loaded.Directories[0].Path != "/tmp/docs" {
		t.Errorf("Expected directory '/tmp/docs', got %v", loaded.Directories)
	}
}
