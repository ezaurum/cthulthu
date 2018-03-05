package cthulthu

import (
	"github.com/fsnotify/fsnotify"
	"log"
)

func WatchDir(targetDir string, run func(watcher *fsnotify.Watcher)) *fsnotify.Watcher {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	go run(watcher)

	err = watcher.Add(targetDir)
	if err != nil {
		log.Fatal(err)
	}
	return watcher
}
