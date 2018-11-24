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
	"strings"
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

	p := ui.NewPar("")
	p.Height = 25
	p.Width = 120
	p.TextFgColor = ui.ColorWhite
	p.BorderLabel = fmt.Sprintf("Websocket console of %s", devicename)
	p.BorderFg = ui.ColorCyan

	g := ui.NewPar("> ")
	g.Height = 5
	g.Width = 120
	g.Y = 25
	g.TextFgColor = ui.ColorWhite
	g.BorderLabel = "Commands"
	g.BorderFg = ui.ColorGreen

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

	//ui.Render(p,g)

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

			if msg.MessageType == 2 {
				// photo/file
				uilog("#"+strings.ToUpper(msg.Source)+" receiver photo!", p, g)
				photo := "c:\\temp\\" + msg.Message
				decode(photo, msg.Data)

				if runtime.GOOS == "windows" {
					//exec.Command("mspaint", photo).Output()
				}
			}

			if msg.MessageType == 1 {

				// echo
				uilog("#"+strings.ToUpper(msg.Source)+":"+msg.Message, p, g)

				// command
				output := doCommand(msg, devicename, c)
				if output != "" {
					uilog("#"+strings.ToUpper(msg.Source)+":"+output, p, g)
				}
				// stop
				if output == "stop" {
					ui.StopLoop()
				}
			}
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
			//p.Text = fmt.Sprintf("%d\n%s", timer, p.Text)
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
