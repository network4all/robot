package main

import (
	"runtime"
	"strings"

	ui "github.com/gizak/termui"
	"github.com/gorilla/websocket"
)

func handleMessage(msg Message, devicename string, log *ui.Par, command *ui.Par, c *websocket.Conn) {

	if msg.MessageType == 1 {

		// echo
		uilog("#"+strings.ToUpper(msg.Source)+":"+msg.Message, log, command)

		// command
		output := doCommand(msg, devicename, c)
		if output != "" {
			uilog("#"+strings.ToUpper(msg.Source)+":"+output, log, command)
		}
		// stop
		if output == "stop" {
			ui.StopLoop()
		}
	}

	if msg.MessageType == 2 {
		// photo/file
		photo := "c:\\temp\\" + msg.Message
		uilog("#"+strings.ToUpper(msg.Source)+" Copy "+photo+" to c:\\temp!", log, command)

		if len(msg.Data) > 0 {
			decode(photo, msg.Data)
		}
		if runtime.GOOS == "windows" {
			//exec.Command("mspaint", photo).Output()
		}
	}

}
