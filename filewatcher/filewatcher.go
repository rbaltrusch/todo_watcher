package filewatcher

import (
	"fmt"
	"log"

	"github.com/fsnotify/fsnotify"
)

func determineOperation(op *fsnotify.Event) string {
	if op.Has(fsnotify.Write) {
		return "WRITE"
	} else if op.Has(fsnotify.Create) {
		return "CREATE"
	} else if op.Has(fsnotify.Remove) {
		return "DELETE"
	} else if op.Has(fsnotify.Rename) {
		return "MOVE"
	}
	return ""
}

func HandleFileEvents(watcher *fsnotify.Watcher, broadcast chan string) {
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			log.Println("event:", event)
			operation := determineOperation(&event)
			message := fmt.Sprintf("%s; %s", operation, event.Name)
			broadcast <- message
			log.Println("Broadcast filewatcher message:", message)
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("Filewatcher error:", err)
		}
	}
}

func CreateWatcher(path string) *fsnotify.Watcher {
	// Create new watcher.
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	// Add a path to watch.
	err = watcher.Add(path)
	if err != nil {
		log.Fatal(err)
	}
	return watcher
}
