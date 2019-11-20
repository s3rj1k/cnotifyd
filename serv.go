package main

import (
	"encoding/json"
	"net"
	"os"
	"strings"

	"golang.org/x/sys/unix"

	"github.com/s3rj1k/go-fanotify/fanotify"
)

const unixSocketPath = "/var/run/cnotifyd.socket"

// ServAPI starts Unix-Socket API.
func ServAPI(notify *fanotify.NotifyFD, db *ContainerDB) {
	listener, err := net.Listen("unix", unixSocketPath)
	if err != nil {
		Error.Fatalf("cnotifyd: socket error, %v\n", err)
	}

	defer listener.Close()

	if err := os.Chmod(unixSocketPath, 0777); err != nil {
		Error.Fatalf("cnotifyd: socket error, %v\n", err)
	}

	f := func(c net.Conn) {
		defer c.Close()

		data := new(NotifyMarkRequest)

		if err := json.NewDecoder(c).Decode(data); err == nil {
			switch {
			case strings.EqualFold(data.Action, "ADD"):
				AddMark(notify, db, data)
			case strings.EqualFold(data.Action, "REMOVE"):
				RemoveMark(notify, db, data)
			default:
				Debug.Printf("cnotifyd: invalid action\n")
			}
		} else {
			Debug.Printf("cnotifyd: invalid request, %v\n", err)
		}
	}

	for {
		connection, err := listener.Accept()
		if err != nil {
			Error.Fatalf("cnotifyd: socket error, %v\n", err)
		}

		go f(connection)
	}
}

// AddMark handles unix.FAN_MARK_ADD.
func AddMark(notify *fanotify.NotifyFD, db *ContainerDB, req *NotifyMarkRequest) {
	if req == nil {
		Debug.Printf("cnotifyd: request object is nil\n")

		return
	}

	if db == nil {
		Debug.Printf("cnotifyd: dictionary object is nil\n")

		return
	}

	mountID, err := GetRootFsParentMountID(req.PID)
	if err != nil {
		Error.Printf("cnotifyd: failed to add mark for '%s': %v\n", req.GetRootFsPath(), err)

		return
	}

	if !IsValidContainer(req.PID, req.Name) {
		Error.Printf("cnotifyd: invalid container data, name does not match cgroup\n")

		return
	}

	// add mark for container rootfs
	if err := notify.Mark(
		unix.FAN_MARK_ADD|
			unix.FAN_MARK_MOUNT,
		unix.FAN_CLOSE_WRITE,
		unix.AT_FDCWD,
		req.GetRootFsPath(),
	); err != nil {
		Error.Printf("cnotifyd: failed to add mark for '%s', CT:'%s': %v\n", req.GetRootFsPath(), req.Name, err)

		return
	}

	db.BindMountIDToContainerName(mountID, req.Name)
	db.BindPIDToContainerName(req.PID, req.Name)

	Info.Printf("cnotifyd: add mark for '%s', CT:'%s'\n", req.GetRootFsPath(), req.Name)
}

// RemoveMark handles unix.FAN_MARK_REMOVE.
func RemoveMark(notify *fanotify.NotifyFD, db *ContainerDB, req *NotifyMarkRequest) {
	if req == nil {
		Debug.Printf("cnotifyd: request object is nil\n")

		return
	}

	if db == nil {
		Debug.Printf("cnotifyd: dictionary object is nil\n")

		return
	}

	mountID, ok := db.GetMountIDFromName(req.Name)
	if !ok {
		Error.Printf("cnotifyd: container with name '%s', is not registered\n", req.Name)

		return
	}

	db.UnBind(mountID, req.Name, req.PID)

	if _, err := os.Stat(req.GetRootFsPath()); err != nil {
		// remove mark for container rootfs
		_ = notify.Mark(
			unix.FAN_MARK_REMOVE|
				unix.FAN_MARK_MOUNT,
			unix.FAN_CLOSE_WRITE,
			unix.AT_FDCWD,
			req.GetRootFsPath(),
		) // we intentionally ignore errors here
	}

	Info.Printf("cnotifyd: remove mark for '%s', CT:'%s'\n", req.GetRootFsPath(), req.Name)
}
