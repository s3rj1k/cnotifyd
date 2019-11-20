package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// GetRootFsParentMountID returns parent mount ID for rootfs taken from '/proc/PID/mountinfo'.
func GetRootFsParentMountID(pid int) (int, error) {
	var (
		err  error
		file *os.File

		mountInfo string
		fields    []string

		id int
	)

	if file, err = os.Open(fmt.Sprintf("/proc/%d/mountinfo", pid)); err != nil {
		return 0, fmt.Errorf("cnotifyd: procfs error, %w", err)
	}

	defer file.Close()

	r := bufio.NewReader(file)

	if mountInfo, err = r.ReadString('\n'); err != nil {
		return 0, fmt.Errorf("cnotifyd: procfs error, %w", err)
	}

	if fields = strings.Fields(mountInfo); len(fields) < 1 {
		return 0, fmt.Errorf("cnotifyd: wrong '/proc/PID/mountinfo' file format")
	}

	if id, err = strconv.Atoi(fields[0]); err != nil {
		return 0, fmt.Errorf("cnotifyd: wrong '/proc/PID/mountinfo' file format, %w", err)
	}

	return id, nil
}

// IsValidContainer checks that requested container name matches value inside '/proc/PID/cpuset'.
func IsValidContainer(pid int, name string) bool {
	content, err := os.ReadFile(fmt.Sprintf("/proc/%d/cpuset", pid))
	if err != nil {
		Debug.Printf("cnotifyd: procfs error, %v", err)

		return false
	}

	content = bytes.Trim(content, "\x00")
	content = bytes.Trim(content, "\n")
	content = bytes.TrimSpace(content)

	return name == filepath.Base(string(content))
}
