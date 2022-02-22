package main

import (
	"phoenixbuilder_3rd_gui/gui"
)

func main() {
	dataFolder := "./data"
	gui := gui.NewGUI(dataFolder)
	gui.Run()
}
