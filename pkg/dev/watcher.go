package dev

import (
	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"io/fs"
	"os"
	"strings"
)

type WatcherConfig struct {
	// Slice of absolute paths to track
	Dirs       []string
	// Whether to do a recursive watch on those paths
	Recursive  bool
	// List of directories to ignore
	IgnoreDirs []string
	// List of triggers to operate. fsnotify.Write
	// would trigger the action on a file write.
	Triggers   []fsnotify.Op
	// Action to trigger when a file is modified
	Action     func(event fsnotify.Event)
}

type Watcher struct {
	cfg       WatcherConfig
	fsWatcher *fsnotify.Watcher
}

// Start starts the watcher and looks for file changes
func (w *Watcher) Start() {
	var err error
	w.fsWatcher, err = fsnotify.NewWatcher()
	defer func() {
		if err := w.fsWatcher.Close(); err != nil {
			log.Fatalln(err)
		}
	}()
	if err != nil {
		log.Fatalln(err)
	}
	w.bindPaths()
	w.listen()
}

func (w *Watcher) bindPaths() {
	for _, dirPath := range w.cfg.Dirs {
		if !w.cfg.Recursive {
			w.addDir(dirPath)
			return
		}
		dir := os.DirFS(dirPath)
		if err := fs.WalkDir(dir, ".", func(path string, d fs.DirEntry,
			err error) error {
			if err != nil {
				log.Fatalln(err)
			}
			if !w.ignored(path) && d.IsDir() {
				w.addDir(path)
			}
			return nil
		}); err != nil {
			log.Fatalln(err)
		}
	}
}

func (w *Watcher) addDir(path string) {
	if err := w.fsWatcher.Add(path); err != nil {
		log.Fatalln(err)
	}

}

func (w *Watcher) ignored(path string) bool {
	for _, v := range w.cfg.IgnoreDirs {
		if strings.Contains(path, v) {
			return true
		}
	}
	return false
}

func (w *Watcher) listen() {
	for {
		select {
		case event, ok := <-w.fsWatcher.Events:
			if !ok {
				log.Fatalln(event)
			}
			if w.triggerAction(event) {
				w.cfg.Action(event)
			}
		}
	}
}

func (w *Watcher) triggerAction(event fsnotify.Event) bool {
	for _, t := range w.cfg.Triggers {
		if event.Op&t == t {
			return true
		}
	}
	return false
}