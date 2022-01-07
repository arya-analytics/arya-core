package main

import (
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.Info("Starting Dummy Arya Core")
	t := time.NewTicker(5 * time.Second)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := true
	for done {
		select {
		case <-t.C:
			log.Info("This is such a fantastic concept")
		case <-sigs:
			log.Info("Terminating")
			done = false
		}
	}
}
