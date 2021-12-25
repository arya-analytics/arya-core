package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.WithFields(log.Fields{
		"animal": "walrus",
		"size": 10,
	}).Error("info")
	fmt.Println("Starting Dummy Arya Core")
	t := time.NewTicker(5 * time.Second)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := true
	for done {
		select {
		case <-t.C:
			fmt.Println("This is such a fantastic concept")
		case <-sigs:
			fmt.Println("Terminating")
			done = false
		}
	}
}
