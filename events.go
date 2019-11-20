package main

import (
	"os"
	"strings"

	"golang.org/x/sys/unix"

	"github.com/s3rj1k/go-fanotify/fanotify"
)

// Process single fanotify event.
func event(notify *fanotify.NotifyFD, db *ContainerDB) {
	var (
		err   error
		event *fanotify.EventMetadata

		path string
		name string
		pid  int

		fdInfo fanotify.FdInfo

		ok bool
	)

	event, err = notify.GetEvent(os.Getpid())
	if err != nil {
		Error.Printf("%v\n", err)
	}

	if event == nil {
		return
	}

	defer func(event *fanotify.EventMetadata) {
		if err = event.Close(); err != nil {
			Error.Printf("%v\n", err)
		}
	}(event)

	if !event.MatchMask(unix.FAN_CLOSE_WRITE) {
		return
	}

	path, err = event.GetPath()
	if err != nil {
		Error.Printf("%v\n", err)

		return
	}

	for i := range WhitelistedPaths {
		if strings.HasPrefix(path, WhitelistedPaths[i]) {
			ok = true

			break
		}
	}

	if !ok {
		return
	}

	fdInfo, err = event.GetFdInfo()
	if err != nil {
		Error.Printf("%v\n", err)

		return
	}

	if name, ok = db.GetContainerNameFromMountID(fdInfo.MountID); !ok {
		Error.Printf("cnotifyd: container with rootfs mount ID '%d', is not registered\n", fdInfo.MountID)

		return
	}

	if pid, ok = db.GetPIDFromMountID(fdInfo.MountID); ok {
		Info.Printf("CT:'%s' PID:'%d' Path:'%s'\n", name, pid, path)

		return
	}

	Info.Printf("CT:'%s' Path:'%s'\n", name, path)
}
