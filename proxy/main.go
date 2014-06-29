package main

import (
	"flag"
	"log"
	"syscall"

	"github.com/apokalyptik/gopid"
)

var pidFile = ""
var uid = 0
var gid = 0

func init() {
	flag.StringVar(&listenHTTP, "http", listenHTTP, "http address and port number to listen on")
	flag.StringVar(&pidFile, "pid", pidFile, "file to lock and write our PID to (empty string disables)")
	flag.DurationVar(&cacheDuration, "cache", cacheDuration, "keep items in the local cache for this long")
	flag.IntVar(&uid, "uid", uid, "set UID (0 disables)")
	flag.IntVar(&gid, "gid", gid, "set GID (0 disables)")
}

func main() {
	flag.Parse()
	if pidFile != "" {
		if _, err := pid.Do(pidFile); err != nil {
			log.Fatalf("error locking pid file: %s", err.Error())
		}
	}
	if gid != 0 {
		if err := syscall.Setgid(gid); err != nil {
			log.Fatal(err)
		}
	}
	if uid != 0 {
		if err := syscall.Setuid(uid); err != nil {
			log.Fatal(err)
		}
	}
	go mindEngine()
	mindHTTP()
}
