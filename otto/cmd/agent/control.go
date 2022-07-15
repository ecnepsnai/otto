package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"github.com/ecnepsnai/logtic"
	"github.com/ecnepsnai/snapshot"
)

var controlPath = ".control"
var controlFd net.Listener

func setupControl() {
	if fileExists(controlPath) {
		if err := os.Remove(controlPath); err != nil && !os.IsNotExist(err) {
			panic("cant remove agent pid: " + err.Error())
		}
	}

	l, err := net.Listen("unix", controlPath)
	if err != nil {
		panic("err listen agent pid: " + err.Error())
	}
	controlFd = l
}

func controlMain() {
	for {
		conn, err := controlFd.Accept()
		if err != nil {
			panic("pid accept err: " + err.Error())
		}
		go controlAccept(conn)
	}
}

func controlAccept(conn net.Conn) {
	defer conn.Close()

	conn.Write([]byte(fmt.Sprintf("Otto %s. Type command + newline, then ^D to submit.\n\n# ", Version)))

	data, _ := io.ReadAll(conn)
	command := string(data[0 : len(data)-1])
	switch command {
	case "stat":
		data, err := formatJSON(Stats)
		if err != nil {
			conn.Write([]byte(fmt.Sprintf("err: %s", err.Error())))
		} else {
			conn.Write(data)
		}
	case "dump":
		name := fmt.Sprintf("dump_%d_%s.zip", os.Getpid(), time.Now().Format("20060102150405"))
		if err := snapshot.Full(name); err != nil {
			conn.Write([]byte(fmt.Sprintf("err: %s", err.Error())))
		} else {
			conn.Write([]byte(fmt.Sprintf("dump saved to %s", name)))
		}
	case "config":
		data, err := formatJSON(*config)
		if err != nil {
			conn.Write([]byte(fmt.Sprintf("err: %s", err.Error())))
		} else {
			conn.Write(data)
		}
	case "reload":
		mustLoadConfig()
		mustLoadIdentity()
		conn.Write([]byte("config & identity reloaded"))
	case "debug":
		logtic.Log.Level = logtic.LevelDebug
		conn.Write([]byte("debug logging enabled"))
	case "nodebug":
		logtic.Log.Level = logtic.LevelError
		conn.Write([]byte("debug logging disabled"))
	case "help":
		conn.Write([]byte("valid commands are: stat, dump, config, reload, debug, nodebug, help"))
	default:
		conn.Write([]byte(fmt.Sprintf("unknown command '%s'", command)))
	}

	conn.Write([]byte("\n"))
}

func stopControl() {
	controlFd.Close()
}
