package watcher

import (
	"io/fs"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

const debounceDuration = 300 * time.Millisecond

type Watcher struct {
	fsWatcher *fsnotify.Watcher
	timers    map[string]*time.Timer
	mu        sync.Mutex
}

func New() (*Watcher, error) {
	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	return &Watcher{
		fsWatcher: fsWatcher,
		timers:    make(map[string]*time.Timer),
	}, nil
}

func (w *Watcher) Close() error {
	w.mu.Lock()

	for _, timer := range w.timers {
		timer.Stop()
	}

	w.timers = make(map[string]*time.Timer)

	w.mu.Unlock()

	return w.fsWatcher.Close()
}

func (w *Watcher) AddRecursive(root string) error {
	return filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !entry.IsDir() {
			return nil
		}

		if shouldSkipDir(entry.Name()) {
			return filepath.SkipDir
		}

		return w.fsWatcher.Add(path)
	})
}

func (w *Watcher) Run(onChange func(string)) error {
	for {
		select {
		case event, ok := <-w.fsWatcher.Events:
			if !ok {
				return nil
			}

			if event.Op&fsnotify.Write == 0 {
				continue
			}

			if filepath.Ext(event.Name) != ".go" {
				continue
			}

			w.debounce(event.Name, onChange)

		case err, ok := <-w.fsWatcher.Errors:
			if !ok {
				return nil
			}

			return err
		}
	}
}

func (w *Watcher) debounce(path string, onChange func(string)) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if timer, ok := w.timers[path]; ok {
		timer.Stop()
	}

	w.timers[path] = time.AfterFunc(debounceDuration, func() {
		onChange(path)

		w.mu.Lock()
		delete(w.timers, path)
		w.mu.Unlock()
	})
}

func shouldSkipDir(name string) bool {
	switch name {
	case ".git", "vendor":
		return true
	default:
		return false
	}
}
