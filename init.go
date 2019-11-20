package main

import (
	"log"
	"os"
	"strings"

	"golang.org/x/sys/unix"
)

var (
	// Logging levels.
	Debug *log.Logger
	Info  *log.Logger
	Error *log.Logger

	// Whitelisted paths inside container.
	WhitelistedPaths = []string{"/var/www/"}
)

func init() {
	// run LXD hook code and exit
	if len(os.Args) == 2 {
		if strings.EqualFold(os.Args[1], "-hook") {
			LXDHook()

			os.Exit(0)
		}
	}

	// initialize loggers
	Debug = log.New(
		os.Stdout,
		"DEBUG: ",
		log.Ldate|log.Ltime|log.Lshortfile,
	)
	Info = log.New(
		os.Stdout,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile,
	)
	Error = log.New(
		os.Stderr,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile,
	)

	if os.Geteuid() != 0 {
		Error.Fatalf("cnotifyd: must be run under root\n")
	}

	if _, err := os.Stat(unixSocketPath); err == nil {
		if err := unix.Unlink(unixSocketPath); err != nil {
			Error.Fatalf("cnotifyd: unlink '%s', %v\n", unixSocketPath, err)
		}
	}
}
