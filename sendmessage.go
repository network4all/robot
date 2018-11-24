package main

import (
	"fmt"
	"log"
	"time"
	"io/ioutil"
	"path/filepath"
	"github.com/gorilla/websocket"
	"strings"
)

func sendMessageTo(destination string, message string, msgtype int, data string, device string, c *websocket.Conn) {
	t := time.Now()

	var msg Message
	msg.MessageID = fmt.Sprintf("%s %s", device, t.Format(time.StampMilli))
	msg.MessageType = msgtype
	msg.Source = device
	msg.Destination = destination
	msg.Message = message
	msg.Data = data
	msg.Ack = false

	err := c.WriteJSON(msg)
	if err != nil {
		// todo: reconnect!
		log.Println("write:", err)
		return
	}
}

func sendMessage(message string, msgtype int, device string, c *websocket.Conn) {
	sendMessageTo("", message, msgtype, "", device, c)
}

func sendAllPhotos(destination string, device string, c *websocket.Conn) (int, error) {
	photopath := "/root/scripts/photo/"
	fis, err := ioutil.ReadDir(photopath)
	if err != nil {
		return 0, fmt.Errorf("could not read dir : %v", err)

	}
	for _, fi := range (fis) {
		name := strings.ToLower(fi.Name())
   		if fi.IsDir() || filepath.Ext(name) != ".jpeg" {
   			continue
   		}
        sendMessage("sending: "+ name, 1, destination, c)
        encoded := encode(photopath + name)
        if len (encoded) >0 {
           sendMessageTo(destination, fmt.Sprintf("%s", fi.Name()), 2, encoded, device, c)
        }
   	}
   	return 1, nil
}


func sendPhoto(destination string, device string, c *websocket.Conn) int {
	photo := "/root/scripts/photo/201811230854.jpeg"
	encoded := encode(photo)
	sendMessageTo(destination, "201811230854.jpeg", 2, encoded, device, c)
	return len(encoded)
}
