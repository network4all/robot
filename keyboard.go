package main

import (
	ui "github.com/gizak/termui"
	"github.com/gorilla/websocket"
)

func mykeyboard(commandline string, device string, keypressed string, log *ui.Par, c *websocket.Conn) string {

	// ks := []string{"<Insert>", "<Delete>", "<Home>", "<End>", "<Previous>", "<Next>", "<Up>", "<Down>", "<Left>", "<Right>"}
	switch keypressed {
	case "<Escape>":
		ui.StopLoop()
	case "<C-c>":
		ui.StopLoop()
	case "<Space>":
		commandline = commandline + " "
	case "<Enter>":
		sendMessage(commandline, 1, device, c)
		//log.Text = fmt.Sprintf("%s\n%s", commandline, log.Text)
		commandline = ""
	case "<Backspace>":
	case "<\b>":
		if len(commandline) > 0 {
			runes := []rune(commandline)
			commandline = string(runes[0 : len(commandline)-1])
		}
	default:
		commandline = commandline + keypressed
	}
	return commandline
}

