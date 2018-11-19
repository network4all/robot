package main

import (
   "time"
   "fmt"
   ui "github.com/gizak/termui"
   "github.com/gorilla/websocket"
   "log"
   "os" 
   "os/signal"
   "net/url"
   "encoding/json"
   "runtime"
   "os/exec"
   "strings"
   // "bytes"
   d "./conf"
)

// Define our message object
type Message struct {
        MessageId   string `json:"messageid"`   // timestamp+node
        MessageType int    `json:"messagetype"` // ping, sendobject, ...
        Source      string `json:"source"`      // node, serial mac
        Destination string `json:"destination"` // broadcast, serial mac
        Message     string `json:"message"`     // json object data
        Ack         bool   `json:"ack"`         // read ack (tcp/udp) (true/false)
}

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

	ui.Render(p,g)
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
                    // echo
                    uilog("#" + strings.ToUpper(msg.Source) + ":" + msg.Message, p, g)

                    // command
                    output := doCommand(msg, devicename, c)
                    if (output != "") {
                    	uilog("#" + strings.ToUpper(msg.Source) + ":" + output, p, g)
                    }

                    if output == "stop" {
                    	ui.StopLoop()
                    }

            }
    }()

    // welcome
    sendMessage(fmt.Sprintf("client %s (%s) is connected to %s", devicename, runtime.GOOS, u.String()), devicename, c)

    // key   
	ui.Handle("<Keyboard>" , func(e ui.Event) {
		cmd = mykeyboard(cmd, devicename, e.ID, p, c) 
		g.Text = fmt.Sprintf("> %s", cmd)
		ui.Render(p,g)
	})


    // keepalive
    //keepaliveTicker := time.NewTicker(time.Second*15)

    // update loop
    drawTicker := time.NewTicker(time.Second)
	go func() {
        for {
			//p.Text = fmt.Sprintf("%d\n%s", timer, p.Text)
			ui.Render(p,g)
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

func mykeyboard (commandline string, device string, keypressed string, log *ui.Par, c *websocket.Conn) string {

 	// ks := []string{"<Insert>", "<Delete>", "<Home>", "<End>", "<Previous>", "<Next>", "<Up>", "<Down>", "<Left>", "<Right>"}    
    switch (keypressed) {
    	case "<Escape>":
    		ui.StopLoop()
    	case "<C-c>":
    		ui.StopLoop()
    	case "<Space>":
    		commandline = commandline + " "
    	case "<Enter>":
    		sendMessage(commandline, device, c)
    		//log.Text = fmt.Sprintf("%s\n%s", commandline, log.Text)
    		commandline = ""
    	case "<Backspace>":
    		if len(commandline)>0 {
    			runes := []rune(commandline)
    			commandline = string(runes[0:len(commandline)-1])
    		}
    	default: commandline = commandline + keypressed
    }
	return commandline
}

func sendMessage (message string, device string, c *websocket.Conn) {
	t := time.Now()
	
    var msg Message
    msg.MessageId   = fmt.Sprintf("%s %s", device, t.Format(time.StampMilli))
    msg.MessageType = 1
    msg.Source      = device
    msg.Destination = ""
    msg.Message     = message
    msg.Ack         = false

    err := c.WriteJSON(msg)
    if err != nil {
            // todo: reconnect!
            log.Println("write:", err)
            return
    }
}

func uilog (line string, log *ui.Par, command *ui.Par) {
	log.Text = fmt.Sprintf("%s\n%s", line, log.Text)
	ui.Render(log, command)
}

func executeShell(cmd string) string {
	if cmdout, err := exec.Command("sh","-c", cmd).Output(); err != nil {
		return err.Error()
	} else {
		return fmt.Sprint("\n" + string(cmdout))
	}
}

func doCommand(msg Message, devicename string, c *websocket.Conn, ) string {

    whoami := "@" + devicename + " "

    if strings.HasPrefix(msg.Message, whoami) {
    	// for me
		command := strings.Replace(msg.Message, whoami, "", -1)
		
		// restart
		if (command == "restart") {
			sendMessage ("terminating console!", devicename, c)
			return "stop"
		}

		// execute command
		if (strings.HasPrefix(command, "shell")) {
			shell := strings.Replace(command, "shell ", "", -1)
			sendMessage ("execute command: '" + shell + "'.", devicename, c)
			out := executeShell(shell)
			sendMessage (out, devicename, c)
			return ""
		}
		sendMessage ("received command '" + command + "'", devicename, c)

	} else {
		// for all
		if strings.HasPrefix(msg.Message, "hi") {
	    	answer := fmt.Sprintf("device #%s says hi\n", devicename)
	    	sendMessage (answer, devicename, c)
    		return ""
    	}
	}

    // nothing to do
	return ""
}

func ping(ws *websocket.Conn, done chan struct{}) {
    ticker := time.NewTicker(time.Second * 15)
    defer ticker.Stop()
    for {
        select {
        case <-ticker.C:
            if err := ws.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(writeWait)); err != nil {
                log.Println("ping:", err)
            }
        case <-done:
            return
        }
    }
}