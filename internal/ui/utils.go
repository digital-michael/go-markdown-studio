package ui

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"go-markdown-studio/internal/config"
)

// ScanMarkdownFilesFromConfig scans markdown files based on the provided configuration.
func ScanMarkdownFilesFromConfig(cfg *config.AppConfig) []string {
	var allFiles []string
	for _, entry := range cfg.Directories {
		if entry.Recursive {
			allFiles = append(allFiles, scanMarkdownFilesRecursive(entry.Path)...)
		} else {
			allFiles = append(allFiles, scanMarkdownFilesFlat(entry.Path)...)
		}
	}
	return allFiles
}

// scanMarkdownFilesRecursive walks a directory and returns all .md/.MD files.
func scanMarkdownFilesRecursive(dir string) []string {
	var files []string
	filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if !d.IsDir() && strings.HasSuffix(strings.ToLower(d.Name()), ".md") {
			files = append(files, path)
		}
		return nil
	})
	return files
}

// scanMarkdownFilesFlat scans only one level of the specified directory for markdown files.
func scanMarkdownFilesFlat(dir string) []string {
	var files []string
	entries, err := os.ReadDir(dir)
	if err != nil {
		return files
	}
	for _, d := range entries {
		if !d.IsDir() && strings.HasSuffix(strings.ToLower(d.Name()), ".md") {
			files = append(files, filepath.Join(dir, d.Name()))
		}
	}
	return files
}
