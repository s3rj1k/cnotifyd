package main

import (
	"os"

	"golang.org/x/sys/unix"

	"github.com/s3rj1k/go-fanotify/fanotify"
)

// Exit is a helper to do cleanup on exit.
func Exit(notify *fanotify.NotifyFD, sig os.Signal) {
	Info.Printf("cnotifyd: caught signal '%s': shutting down\n", sig)

	_ = os.Remove(unixSocketPath)

	if notify != nil {
		if err := notify.Mark(
			unix.FAN_MARK_FLUSH|
				unix.FAN_MARK_MOUNT,
			0, -1, "",
		); err != nil {
			Error.Printf("cnotifyd: failed to flush marks: %v\n", err)
		}
	}
}
