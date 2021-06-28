package app

import (
	"github.com/jroimartin/gocui"
	"os"
)

func Keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, exit); err != nil {
		return err
	}
	return nil
}

func exit(g *gocui.Gui, v *gocui.View) error {
	os.Exit(0)
	return nil
}
