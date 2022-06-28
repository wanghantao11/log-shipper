package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"

	"github.com/wanghantao11/log-shipper/config"
	logsvc "github.com/wanghantao11/log-shipper/internal/pkg/log"
)

const serviceName = "log-shipper"

var watcher *fsnotify.Watcher

// main
func main() {
	config.Init(serviceName)

	// Initialize log service
	logService := logsvc.New()

	// creates a new file watcher
	watcher, _ = fsnotify.NewWatcher()
	defer watcher.Close()

	// starting at the root of the project, walk each file/directory searching for
	// directories
	directory := config.Get(config.Path)
	if err := filepath.Walk(directory, watchDir); err != nil {
		fmt.Println("ERROR", err)
	}

	done := make(chan bool)

	go func() {
		for {
			select {
			// watch for events
			case event := <-watcher.Events:
				log.Println("event:", event)
				if event.Op&fsnotify.Create == fsnotify.Create {
					// Parse file contents
					fileData, err := logService.ParseFile(event.Name)
					if err != nil {
						fmt.Println("Fail to parse file", err)
					}

					log.Println("fileData:", fileData)

					// Send API call to log-receiver to create logs
					err = logService.AddLog(fileData)
					if err != nil {
						fmt.Println("Fail to add log via api call", err)
					}
				}

				// watch for errors
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error", err)
			}
		}
	}()

	<-done
}

func watchDir(path string, fi os.FileInfo, err error) error {
	// since fsnotify can watch all the files in a directory, watchers only need
	// to be added to each nested directory
	if fi.Mode().IsDir() {
		return watcher.Add(path)
	}

	return nil
}
