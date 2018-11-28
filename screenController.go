package main

import (
	"fmt"

	ui "github.com/gizak/termui"
)

func initscreen(p *ui.Par, label string) {
	p.Height = 25
	p.Width = 120
	p.TextFgColor = ui.ColorWhite
	p.BorderLabel = label
	p.BorderFg = ui.ColorCyan
}

func initconsole(g *ui.Par) {
	g.Height = 5
	g.Width = 120
	g.Y = 25
	g.TextFgColor = ui.ColorWhite
	g.BorderLabel = "Commands"
	g.BorderFg = ui.ColorGreen
}

func uilog(line string, log *ui.Par, command *ui.Par) {
	log.Text = fmt.Sprintf("%s\n%s", line, log.Text)
	ui.Render(log, command)
}
