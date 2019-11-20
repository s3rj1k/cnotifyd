package main

import (
	"fmt"
)

// NotifyMarkRequest contains PID, that is needed to compose NS rootfs path what we want to mark for fanotify.
type NotifyMarkRequest struct {
	Action string `json:"Action"`
	PID    int    `json:"PID"`
	Name   string `json:"Name"`
}

// GetRootFsPath returns NS rootfs path inside procfs.
func (r NotifyMarkRequest) GetRootFsPath() string {
	return fmt.Sprintf("/proc/%d/root", r.PID)
}
