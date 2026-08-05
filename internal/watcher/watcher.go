package watcher

import (
	"fmt"
	"time"

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

	var timer *time.Timer
	var changedFile string

	for {
		select {
		case event := <-watcher.Events:

			// 書き込み以外は無視
			if event.Op&fsnotify.Write == 0 {
				continue
			}

			changedFile = event.Name

			if timer != nil {
				timer.Stop()
			}

			timer = time.AfterFunc(300*time.Millisecond, func() {
				fmt.Println("changed:", changedFile)

				issues := analyzer.Analyze(changedFile)

				for _, issue := range issues {
					analyzer.Report(issue)
				}
			})

		case err := <-watcher.Errors:
			fmt.Println("watch error:", err)
		}
	}
}
