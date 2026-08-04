package watcher

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/yokkesan/apiguard/internal/analyzer"
)

func Watch(path string) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	err = watcher.Add(path)
	if err != nil {
		return err
	}

	fmt.Println("watching:", path)

	for {
		select {
		case event := <-watcher.Events:
			fmt.Println("changed:", event.Name)

			issues := analyzer.Analyze(event.Name)

			for _, issue := range issues {
				analyzer.Report(issue)
			}

		case err := <-watcher.Errors:
			fmt.Println("watch error:", err)
		}
	}
}
