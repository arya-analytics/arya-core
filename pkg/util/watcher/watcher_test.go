package watcher_test

import (
	"github.com/arya-analytics/aryacore/pkg/util/watcher"
	"github.com/fsnotify/fsnotify"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"log"
	"os"
	"time"
)

var _ = Describe("Watcher", func() {
	Describe("NewWatcher", func() {
		It("Should create a new watcher", func() {
			cfg := watcher.WatcherConfig{}
			_, err := watcher.NewWatcher(cfg)
			Expect(err).To(BeNil())
		})
	})
	Describe("Start", func() {
		It("Should start the watcher and listen for file changes", func() {
			triggered := false
			cfg := watcher.WatcherConfig{
				Dirs:      []string{tmpDir},
				Recursive: false,
				Triggers:  []fsnotify.Op{fsnotify.Create},
				Action: func(event fsnotify.Event) {
					triggered = true
				},
			}
			w, err := watcher.NewWatcher(cfg)
			if err != nil {
				log.Fatalln(err)
			}
			go w.Start()
			t := time.NewTicker(200 * time.Millisecond)
			var tmpFiles []*os.File
			defer func() {
				for _, f := range tmpFiles {
					err := f.Close()
					err = os.Remove(f.Name())
					if err != nil {
						log.Fatalln(err)
					}
				}

			}()
			for !triggered {
				select {
				case <-t.C:
					f, err := os.CreateTemp(tmpDir, "randomtmpfile")
					tmpFiles = append(tmpFiles, f)
					if err != nil {
						log.Fatalln(err)
					}
				}

			}
			Expect(triggered).To(BeTrue())
		})
	})
})
