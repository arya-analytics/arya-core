package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	fmt.Println("Starting Dummy Arya Core")
	t := time.NewTicker(5 * time.Second)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := true
	for done {
		select {
		case <-t.C:
			fmt.Println("Ticker X")
		case <-sigs:
			fmt.Println("Terminating")
			done = false
		}
	}
}
