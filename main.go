package main

import (
	"os"
	"os/signal"

	"golang.org/x/sys/unix"

	"github.com/s3rj1k/go-fanotify/fanotify"
)

func main() {
	// initialize fanotify
	notify, err := fanotify.Initialize(
		unix.FAN_CLOEXEC|
			unix.FAN_CLASS_NOTIF|
			unix.FAN_UNLIMITED_QUEUE|
			unix.FAN_UNLIMITED_MARKS,
		os.O_RDONLY|
			unix.O_LARGEFILE|
			unix.O_CLOEXEC,
	)
	if err != nil {
		Error.Fatalf("%v\n", err)
	}

	// process SIGTERM signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, unix.SIGTERM)

	// do cleanup on exit
	go func(notify *fanotify.NotifyFD, c chan os.Signal) {
		sig := <-c
		Exit(notify, sig)
		os.Exit(0)
	}(notify, sigChan)

	// initialize global container dictionary
	db := NewContainerDB()

	// Serv API
	go ServAPI(notify, db)

	// get events
	for {
		event(notify, db)
	}
}
