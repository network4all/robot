package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"

	//"os/exec"
	"os/signal"
	"runtime"
	"time"

	ui "github.com/gizak/termui"
	"github.com/gorilla/websocket"

	// "bytes"
	d "./conf"
)

func main() {
	devicename := d.Id()
	// ui
	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	// build screen
	p := ui.NewPar("")
	g := ui.NewPar(">")
	initscreen(p, "Websocket console of "+devicename)
	initconsole(g)
	ui.Render(p, g)
	cmd := ""

	// websocket
	log.SetFlags(0)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	u := url.URL{Scheme: "wss", Host: d.Url(), Path: d.Path()}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	// message reader
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			// read
			_, messagejson, err := c.ReadMessage()
			if err != nil {
				uilog(err.Error(), p, g)
				return
			}
			// decode msg
			var msg Message
			err = json.Unmarshal(messagejson, &msg)
			if err != nil {
				uilog(err.Error(), p, g)
				return
			}
			// handle
			handleMessage(msg, devicename, p, g, c)
		}
	}()

	// welcome
	sendMessage(fmt.Sprintf("client %s (%s) is connected to %s", devicename, runtime.GOOS, u.String()), 1, devicename, c)

	// key
	ui.Handle("<Keyboard>", func(e ui.Event) {
		cmd = mykeyboard(cmd, devicename, e.ID, p, c)
		g.Text = fmt.Sprintf("> %s", cmd)
		ui.Render(p, g)
	})

	// keepalive
	go ping(c)

	// update loop
	drawTicker := time.NewTicker(time.Second)
	go func() {
		for {
			ui.Render(p, g)
			<-drawTicker.C
		}
	}()
	ui.Loop()

	// close connection
	err = c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Println("write close:", err)
		return
	}
	select {
	case <-done:
	case <-time.After(time.Second):
	}
}
