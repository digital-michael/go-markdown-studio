package ui

import (
	// "io/fs"
	"log"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

func WatchMarkdownDir(dir string, updateFiles func()) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Println("Failed to initialize watcher:", err)
		return
	}

	go func() {
		defer watcher.Close()

		debounce := time.NewTimer(0)
		<-debounce.C // discard initial tick

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if isRelevantEvent(event) {
					debounce.Reset(500 * time.Millisecond)
				}

			case <-debounce.C:
				updateFiles() // refresh list in UI
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("Watcher error:", err)
			}
		}
	}()

	err = watcher.Add(dir)
	if err != nil {
		log.Println("Failed to add watch directory:", err)
	}
}

func isRelevantEvent(event fsnotify.Event) bool {
	name := event.Name
	return strings.HasSuffix(name, ".md") && (event.Op&fsnotify.Create == fsnotify.Create || event.Op&fsnotify.Remove == fsnotify.Remove)
}
