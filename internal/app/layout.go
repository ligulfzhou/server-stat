package app

import (
	"github.com/jroimartin/gocui"
	"log"
)

func (app *App) Layout(gui *gocui.Gui) error {
	maxX, maxY := app.Gui.Size()

	thv, thvErr := app.Gui.SetView(TableHeaderViewName, 0, 0, maxX, 2)
	tv, tvErr := app.Gui.SetView(TableViewName, 0, 2, maxX, maxY)
	if tvErr != nil && tvErr != gocui.ErrUnknownView {
		log.Fatalf("set table_header failed with %s", tvErr.Error())
	}
	if thvErr != nil && thvErr != gocui.ErrUnknownView {
		log.Fatalf("set table_header failed with %s", thvErr.Error())
	}
	tv.Frame = false
	thv.Frame = false
	thv.FgColor = gocui.ColorBlack
	thv.BgColor = gocui.ColorGreen
	go app.RefreshAll()

	return nil
}
