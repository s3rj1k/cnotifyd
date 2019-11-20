package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

// LXDHook sends container state data to a server API via Unix-Socket.
// https://linuxcontainers.org/ru/lxc/manpages/man5/lxc.container.conf.5.html
func LXDHook() {
	if !strings.EqualFold(os.Getenv("LXC_HOOK_SECTION"), "lxc") {
		fmt.Printf("HOOK: wrong LXC section\n")

		return
	}

	if _, err := os.Stat(unixSocketPath); err != nil {
		fmt.Printf("HOOK: unix socket does not exist\n")

		return
	}

	c, err := net.Dial("unix", unixSocketPath)
	if err != nil {
		fmt.Printf("HOOK: unix socket connection error\n")

		return
	}
	defer c.Close()

	data := new(NotifyMarkRequest)

	data.Name = os.Getenv("LXC_NAME")

	data.PID, err = strconv.Atoi(os.Getenv("LXC_PID"))
	if err != nil {
		fmt.Printf("HOOK: container PID is not a number\n")

		return
	}

	switch {
	case strings.EqualFold(os.Getenv("LXC_HOOK_TYPE"), "start-host"):
		data.Action = "Add"
	case strings.EqualFold(os.Getenv("LXC_HOOK_TYPE"), "post-stop"):
		data.Action = "Remove"
	}

	if err = json.NewEncoder(c).Encode(data); err != nil {
		fmt.Printf("HOOK: writing to unix socket failed\n")

		return
	}
}
