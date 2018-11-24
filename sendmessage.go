package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

func sendMessageTo(destination string, message string, msgtype int, device string, c *websocket.Conn) {
	t := time.Now()

	var msg Message
	msg.MessageID = fmt.Sprintf("%s %s", device, t.Format(time.StampMilli))
	msg.MessageType = msgtype
	msg.Source = device
	msg.Destination = destination
	msg.Message = message
	msg.Data = ""
	msg.Ack = false

	err := c.WriteJSON(msg)
	if err != nil {
		// todo: reconnect!
		log.Println("write:", err)
		return
	}
}

func sendMessage(message string, msgtype int, device string, c *websocket.Conn) {
	sendMessageTo("", message, msgtype, device, c)
}

func sendAllPhotos(destination string, device string, c *websocket.Conn) (int, error) {
	photopath := "/root/scripts/photo/"
	fis, err := ioutil.ReadDir(photopath)
	if err != nil {
		// return nil, 
	}

	for_, fi := range (fis) {
   		if fi.IsDir || filepath.Ext(name) := ".jpeg" {
   			continue
   		}
        sendMessage("sending: "+ fi.Name, 1, destination, c)
   	}
}


func sendPhoto(destination string, device string, c *websocket.Conn) int {

	photo := "/root/scripts/photo/201811230854.jpeg"
	encoded := encode(photo)
	sendMessageTo(destination, encoded, 2, device, c)
	return len(encoded)
}
