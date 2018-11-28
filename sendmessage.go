package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

// send message to destination
func sendMessageTo(destination string, message string, msgtype int, data string, device string, c *websocket.Conn) error {
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
		//Todo: reconnect
		return fmt.Errorf("could not decode string :%v", err)
	}
	return nil
}

// send message to all
func sendMessage(message string, msgtype int, device string, c *websocket.Conn) {
	sendMessageTo("", message, msgtype, "", device, c)
}

// send last 4 photo's
func sendAllPhotos(destination string, device string, c *websocket.Conn) (int, error) {
	photopath := "/root/scripts/photo/"
	fis, err := ioutil.ReadDir(photopath)
	if err != nil {
		return 0, fmt.Errorf("could not read dir : %v", err)
	}

	// sort date
	sort.Slice(fis, func(i, j int) bool {
		return fis[i].ModTime().Unix() > fis[j].ModTime().Unix()
	})

	// send photo's
	cnt := 0
	for _, fi := range fis {
		cnt++
		if cnt > 4 {
			break
		}

		name := strings.ToLower(fi.Name())
		if fi.IsDir() || filepath.Ext(name) != ".jpeg" {
			continue
		}

		sendMessage("sending: "+name, 1, destination, c)
		encoded, err := encode(photopath + name)
		if err != nil {
			return 0, fmt.Errorf("sendphoto failed :%s :%v", name, err)
		}
		sendMessageTo(destination, name, 2, encoded, device, c)

	}
	return 1, nil
}

// send 1 photo
func sendPhoto(destination string, device string, c *websocket.Conn) (int, error) {
	photo := "/root/scripts/photo/201811230854.jpeg"
	encoded, err := encode(photo)
	if err != nil {
		sendMessageTo(destination, "201811230854.jpeg", 2, encoded, device, c)
		return 0, fmt.Errorf("sendphoto failed :%s :%v", photo, err)
	}
	return len(encoded), nil
}
