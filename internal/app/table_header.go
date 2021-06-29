package app

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"math"
	"server-stat/pkg/pad"
	"strings"
	"unicode/utf8"
)

// UpdateTableHeader renders the table header
func (app *App) UpdateTableHeader() error {
	// headers := []string{}
	var headers []string
	for i, col := range tableHeaders {
		width := app.GetTableColumnWidth(col)

		leftAlign := app.GetTableColumnAlignLeft(col)
		padfn := pad.Left
		padLeft := 1
		if i == 0 {
			padLeft = 0
		}
		if leftAlign {
			padfn = pad.Right
		}

		colStr := fmt.Sprintf(
			"%s%s%s",
			strings.Repeat(" ", padLeft),
			padfn(col, width+(1-padLeft), " "),
			strings.Repeat(" ", 1),
		)
		headers = append(headers, colStr)
	}

	app.Gui.Update(func(gui *gocui.Gui) error {
		v, err := app.Gui.View(TableHeaderViewName)
		if err != nil {
			return err
		}
		v.Clear()
		fmt.Fprintln(v, strings.Join(headers, ""))
		return nil
	})
	return nil
}

// SetTableColumnAlignLeft sets the column alignment direction for header
func (app *App) SetTableColumnAlignLeft(header string, alignLeft bool) {
	app.State.tableColumnAlignLeft.Store(header, alignLeft)
}

// GetTableColumnAlignLeft gets the column alignment direction for header
func (app *App) GetTableColumnAlignLeft(header string) bool {
	ifc, ok := app.State.tableColumnAlignLeft.Load(header)
	if ok {
		return ifc.(bool)
	}
	return false
}

// SetTableColumnWidth sets the column width for header
func (app *App) SetTableColumnWidth(header string, width int) {
	prevIfc, ok := app.State.tableColumnWidths.Load(header)
	var prev int
	if ok {
		prev = prevIfc.(int)
	} else {
		prev = utf8.RuneCountInString(header)
	}

	app.State.tableColumnWidths.Store(header, int(math.Max(float64(width), float64(prev))))
}

// SetTableColumnWidthFromString sets the column width for header given size of string
func (app *App) SetTableColumnWidthFromString(header string, text string) {
	app.SetTableColumnWidth(header, utf8.RuneCountInString(text))
}

// GetTableColumnWidth gets the column width for header
func (app *App) GetTableColumnWidth(header string) int {
	ifc, ok := app.State.tableColumnWidths.Load(header)
	if ok {
		return ifc.(int)
	}
	return 0
}
