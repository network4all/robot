package main

import (
	"fmt"

	ui "github.com/gizak/termui"
)

func uilog(line string, log *ui.Par, command *ui.Par) {
	log.Text = fmt.Sprintf("%s\n%s", line, log.Text)
	ui.Render(log, command)
}
