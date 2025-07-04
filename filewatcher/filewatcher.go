package main

import (
	"fmt"
	"log"

	"github.com/fsnotify/fsnotify"
)

func handleFileEvents(watcher *fsnotify.Watcher) {
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			log.Println("event:", event)
			if event.Has(fsnotify.Write) {
				log.Println("modified file:", event.Name)
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("error:", err)
		}
	}
}

func createWatcher(path string) *fsnotify.Watcher {
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

func main() {
	fmt.Println("Starting file watcher...")
	path := "C:/Users/richa/Desktop/todo"
	watcher := createWatcher(path)
	defer watcher.Close()
	go handleFileEvents(watcher)
	select {} // Block main goroutine forever.
}
